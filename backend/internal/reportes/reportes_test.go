package reportes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/ajustes"
	"finanzas-hogar/internal/api/presupuesto"
	reportesapi "finanzas-hogar/internal/api/reportes"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/reportes"
	"finanzas-hogar/internal/testutil"
)

var keyN atomic.Int64

func key() string { return fmt.Sprintf("s7-pdf-prueba-%08d", keyN.Add(1)) }

type env struct {
	s                        *testutil.Server
	c                        *testutil.Client
	axel, ximena, karen, sp1 int64
}

func ex(t *testing.T, s *testutil.Server, q string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := s.DB.QueryRow(q+" RETURNING id", args...).Scan(&id); err != nil {
		t.Fatalf("%v\n%s", err, q)
	}
	return id
}

// setup: despensa 1 canónica (semilla del 10-oct-2026): 3,094.50 a crédito Kueski en 4 quincenas, 800.50 de intereses,
// repartida Axel 1,903.00 / Karen 1,191.50 (credit_purchase_share y scheduled_payment_share de C1-datos), más un
// ingreso quincenal de Axel que el estado de cuenta de Karen no debe mostrar.
func setup(t *testing.T, areas ...func(*api.Registry)) *env {
	areas = append(areas, reportesapi.Register)
	s := testutil.NewServer(t, areas...)
	e := &env{s: s, c: s.SessionFor(t, testutil.AxelUser), axel: s.Fx.AxelPersonID, ximena: s.Fx.XimenaPersonID, karen: s.Fx.KarenPersonID}
	cat := ex(t, s, `INSERT INTO category (code, name) VALUES ('groceries', 'Despensa')`)
	cr := ex(t, s, `INSERT INTO credit (name, kind, holder_person_id, credit_limit_cents, payment_frequency) VALUES ('Kueski Pay', 'revolving_line', ?, 1000000, 'fortnightly')`, e.axel)
	exp := ex(t, s, `INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, accrual_months, payment_method)
		VALUES ('variable', 'Despensa 1', ?, 309450, '2026-10-10', 2, 'credit')`, cat)
	ex(t, s, `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, 190300)`, exp, e.axel)
	ex(t, s, `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, 119150)`, exp, e.karen)
	pu := ex(t, s, `INSERT INTO credit_purchase (credit_id, expense_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents,
		term_count, payment_frequency, first_payment_date, interest_source, rate_num, rate_den, rate_basis)
		VALUES (?, ?, 'installment_plan', 'Despensa 1', '2026-10-10', 309450, 80050, 389500, 4, 'fortnightly', '2026-10-30', 'quoted', 80050, 309450, 'term_total')`, cr, exp)
	ex(t, s, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 190300, 49228, 239528)`, pu, e.axel)
	ex(t, s, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 119150, 30822, 149972)`, pu, e.karen)
	for i, d := range []string{"2026-10-30", "2026-11-15", "2026-11-30", "2026-12-15"} {
		pr, in := int64(77363), int64(20012)
		if i == 3 {
			pr, in = 77361, 20014
		}
		sp := ex(t, s, `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, ?, ?, 97375, ?, ?)`, pu, i+1, d, pr, in)
		if i == 0 {
			e.sp1 = sp
		}
		ex(t, s, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 59882)`, sp, e.axel)
		ex(t, s, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 37493)`, sp, e.karen)
	}
	var ddl string
	if err := s.DB.QueryRow(`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = 'income'`).Scan(&ddl); err != nil {
		t.Fatal(err)
	}
	kind := regexp.MustCompile(`kind TEXT NOT NULL CHECK \(kind IN \('([a-z_]+)'`).FindStringSubmatch(ddl)
	if kind == nil {
		t.Fatalf("no encontré el enumerado de income.kind")
	}
	ex(t, s, `INSERT INTO income (person_id, kind, description, amount_cents, income_date) VALUES (?, ?, 'Sueldo', 1234567, '2026-10-15')`, e.axel, kind[1])
	return e
}

func (e *env) json(t *testing.T, path string) map[string]any {
	t.Helper()
	res := e.c.Do("GET", path, nil)
	if res.StatusCode != 200 {
		t.Fatalf("GET %s = %d %s\n%s", path, res.StatusCode, res.Body, tail(e.s.Logs.String()))
	}
	var m map[string]any
	d := json.NewDecoder(bytes.NewReader(res.Body))
	d.UseNumber()
	if err := d.Decode(&m); err != nil {
		t.Fatal(err)
	}
	return m
}

func (e *env) pdf(t *testing.T, id int64) string {
	t.Helper()
	res := e.c.Do("GET", fmt.Sprintf("/api/v1/reports/%d/file", id), nil)
	if res.StatusCode != 200 || !bytes.HasPrefix(res.Body, []byte("%PDF-")) {
		t.Fatalf("descarga %d", res.StatusCode)
	}
	f := filepath.Join(testutil.TempDir(t), "r.pdf")
	if err := os.WriteFile(f, res.Body, 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("pdftotext", "-layout", f, "-").Output()
	if err != nil {
		t.Fatalf("pdftotext: %v", err)
	}
	return string(out)
}

var reToken = regexp.MustCompile(`-?\$-?[0-9]{1,3}(?:,[0-9]{3})*\.[0-9]{2}|(?:Dom|Lun|Mar|Mié|Jue|Vie|Sáb) [0-9]{1,2} (?:ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic) [0-9]{4}|[0-9]{1,2}-(?:ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)-[0-9]{4}`)

func resolve(t *testing.T, resp map[string]map[string]any, src string) any {
	t.Helper()
	name, path, _ := strings.Cut(src, "|")
	var cur any = resp[name]
	if cur == nil || resp[name] == nil {
		t.Fatalf("respuesta %q no cargada", name)
	}
	for _, p := range strings.Split(path, ".") {
		switch v := cur.(type) {
		case map[string]any:
			cur = v[p]
		case []any:
			i, _ := strconv.Atoi(p)
			if i >= len(v) {
				t.Fatalf("%s: índice fuera", src)
			}
			cur = v[i]
		default:
			t.Fatalf("%s: ruta inválida", src)
		}
	}
	return cur
}

// cruce compara cada cifra y fecha del PDF con el campo de la API que cita su celda; devuelve M y los campos.
func cruce(t *testing.T, text string, doc *reportes.Doc, resp map[string]map[string]any, only string) (int, []string) {
	t.Helper()
	toks := reToken.FindAllString(text, -1)
	cells := doc.Cells()
	if only == "" && len(toks) != len(cells) {
		t.Fatalf("el PDF trae %d cifras/fechas y el documento %d:\n%v\n%s", len(toks), len(cells), toks, text)
	}
	diffs, m := 0, 0
	var fields []string
	for i, c := range cells {
		if only != "" && !strings.HasPrefix(c.Src, only) {
			continue
		}
		v := resolve(t, resp, c.Src)
		var want string
		switch c.Kind {
		case "money":
			n, err := v.(json.Number).Int64()
			if err != nil {
				t.Fatalf("%s: %v", c.Src, err)
			}
			want = reportes.Pesos(n)
		default:
			s, _ := v.(string)
			w, err := reportes.FechaAPI(s, c.Kind == "date_long")
			if err != nil {
				t.Fatalf("%s: %v", c.Src, err)
			}
			want = w
		}
		m++
		fields = append(fields, c.Src)
		if i >= len(toks) || toks[i] != want {
			diffs++
			got := ""
			if i < len(toks) {
				got = toks[i]
			}
			t.Errorf("diferencia en %s: PDF %q, API %q", c.Src, got, want)
		}
	}
	t.Logf("%d de %d cifras y fechas iguales, %d diferencias", m-diffs, m, diffs)
	return m, fields
}

func accrualResp(t *testing.T, e *env) (map[string]map[string]any, []any) {
	resp := map[string]map[string]any{}
	av := e.json(t, "/api/v1/cash-flow/accrual-vs-cash")
	ms := av["months"].([]any)
	for _, x := range ms {
		mm := x.(map[string]any)
		resp["accrual:"+mm["month"].(string)] = mm
	}
	return resp, ms
}

func oracle(t *testing.T) []map[string]any {
	b, err := os.ReadFile(filepath.Join(testutil.Root(t), "oraculo", "model.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Periodos []map[string]any `json:"periodos"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m.Periodos
}

func cents(v any) int64 {
	switch x := v.(type) {
	case float64:
		return int64(math.Round(x * 100))
	case json.Number:
		n, _ := x.Int64()
		return n
	}
	return -1
}

func generate(t *testing.T, e *env, path string, body any) api.GeneratedReport {
	t.Helper()
	var r api.GeneratedReport
	res := e.c.Do("POST", path, body, "Idempotency-Key", key())
	if res.StatusCode != 202 {
		t.Fatalf("POST %s = %d %s\n%s", path, res.StatusCode, res.Body, tail(e.s.Logs.String()))
	}
	res.JSON(t, &r)
	return r
}

// userCategory crea una categoría por la API (createCategory), como las que Axel captura en Ajustes (U4 P-26.04): su
// category.code queda NULL (00001_esquema.sql, tabla category). Con ella el reporte mensual, la vista web y el job deben
// seguir respondiendo (hallazgo H-1 de REVISION-backend-areas-r1).
func userCategory(t *testing.T, e *env) int64 {
	t.Helper()
	var c api.Category
	res := e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "Mascotas", "kind": "variable"}, "Idempotency-Key", key())
	if res.StatusCode != 201 {
		t.Fatalf("createCategory = %d %s", res.StatusCode, res.Body)
	}
	res.JSON(t, &c)
	var null bool
	if err := e.s.DB.QueryRow(`SELECT code IS NULL FROM category WHERE id = ?`, c.ID).Scan(&null); err != nil || !null {
		t.Fatalf("la categoría creada por la API debe quedar sin code (NULL): %v %v", null, err)
	}
	return c.ID
}

// canonModel es lo de oraculo/model.json que describe las compras de P1–P6 (S2 §6, nota (e)).
type canonModel struct {
	Periodos  []map[string]any `json:"periodos"`
	Despensas []struct {
		N         int     `json:"n"`
		Fecha     string  `json:"fecha"`
		Capital   float64 `json:"capital"`
		Intereses float64 `json:"intereses"`
		Total     float64 `json:"total"`
		Pagos     []struct {
			Fecha           string  `json:"fecha"`
			Monto           float64 `json:"monto"`
			Axel            float64 `json:"axel"`
			Karen           float64 `json:"karen"`
			CapitalRepuesto float64 `json:"capital_repuesto"`
		} `json:"pagos"`
	} `json:"despensas"`
	Disposicion []struct {
		Despensa    int `json:"despensa"`
		Axel, Karen struct {
			Capital   float64 `json:"capital"`
			Intereses float64 `json:"intereses"`
			Total     float64 `json:"total"`
		}
	} `json:"disposicion_por_despensa"`
	Papel []struct {
		Paquete int     `json:"paquete"`
		Compra  string  `json:"compra"`
		Quien   string  `json:"quien"`
		Forma   string  `json:"forma"`
		Monto   float64 `json:"monto"`
	} `json:"papel"`
}

// pc pasa un importe de model.json (pesos con 2 decimales exactos) a centavos enteros.
func pc(x float64) int64 { return int64(math.Round(x * 100)) }

// canon completa la despensa 1 de setup con lo demás que cubre P1–P6 (oct-2026…mar-2027): las despensas 2, 3 y 4 a
// crédito Kueski y el papel alterno de contado. Ninguna cifra se escribe a mano: capital, intereses, partes y pagos salen
// de «despensas» y «disposicion_por_despensa», y el papel de «papel» (forma «de contado»; el de «dentro de despensa» ya
// va en la despensa). Las compras van como la despensa 1: gasto de Despensa con devengo de 2 meses y su compra a plazos.
func canon(t *testing.T, e *env) canonModel {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(testutil.Root(t), "oraculo", "model.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m canonModel
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if len(m.Periodos) < 6 || len(m.Despensas) < 3 || len(m.Disposicion) != len(m.Despensas) {
		t.Fatalf("model.json incompleto: %d periodos, %d despensas, %d disposiciones", len(m.Periodos), len(m.Despensas), len(m.Disposicion))
	}
	s := e.s
	var cat, cr int64
	if err := s.DB.QueryRow(`SELECT id FROM category WHERE code = 'groceries'`).Scan(&cat); err != nil {
		t.Fatal(err)
	}
	if err := s.DB.QueryRow(`SELECT id FROM credit WHERE name = 'Kueski Pay'`).Scan(&cr); err != nil {
		t.Fatal(err)
	}
	for i, d := range m.Despensas {
		if d.N == 1 {
			continue // ya la puso setup (semilla del 10-oct-2026)
		}
		dp := m.Disposicion[i]
		if dp.Despensa != d.N || pc(dp.Axel.Capital)+pc(dp.Karen.Capital) != pc(d.Capital) || pc(d.Capital)+pc(d.Intereses) != pc(d.Total) {
			t.Fatalf("despensa %d de model.json no cuadra consigo misma", d.N)
		}
		desc := fmt.Sprintf("Despensa %d", d.N)
		exp := ex(t, s, `INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, accrual_months, payment_method)
			VALUES ('variable', ?, ?, ?, ?, 2, 'credit')`, desc, cat, pc(d.Capital), d.Fecha)
		ex(t, s, `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, ?)`, exp, e.axel, pc(dp.Axel.Capital))
		ex(t, s, `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, ?)`, exp, e.karen, pc(dp.Karen.Capital))
		pu := ex(t, s, `INSERT INTO credit_purchase (credit_id, expense_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents,
			term_count, payment_frequency, first_payment_date, interest_source, rate_num, rate_den, rate_basis)
			VALUES (?, ?, 'installment_plan', ?, ?, ?, ?, ?, ?, 'fortnightly', ?, 'quoted', ?, ?, 'term_total')`,
			cr, exp, desc, d.Fecha, pc(d.Capital), pc(d.Intereses), pc(d.Total), len(d.Pagos), d.Pagos[0].Fecha, pc(d.Intereses), pc(d.Capital))
		for _, x := range []struct {
			p                 int64
			cap, inter, total float64
		}{{e.axel, dp.Axel.Capital, dp.Axel.Intereses, dp.Axel.Total}, {e.karen, dp.Karen.Capital, dp.Karen.Intereses, dp.Karen.Total}} {
			ex(t, s, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, ?, ?, ?)`,
				pu, x.p, pc(x.cap), pc(x.inter), pc(x.total))
		}
		for j, p := range d.Pagos {
			sp := ex(t, s, `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, ?, ?, ?, ?, ?)`,
				pu, j+1, p.Fecha, pc(p.Monto), pc(p.CapitalRepuesto), pc(p.Monto)-pc(p.CapitalRepuesto))
			ex(t, s, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, ?)`, sp, e.axel, pc(p.Axel))
			ex(t, s, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, ?)`, sp, e.karen, pc(p.Karen))
		}
	}
	if _, err := s.DB.Exec(`INSERT INTO category (code, name) VALUES ('paper', 'Papel') ON CONFLICT (code) DO NOTHING`); err != nil {
		t.Fatal(err)
	}
	var paper int64
	if err := s.DB.QueryRow(`SELECT id FROM category WHERE code = 'paper'`).Scan(&paper); err != nil {
		t.Fatal(err)
	}
	who := map[string]int64{"Axel": e.axel, "Karen": e.karen}
	for _, p := range m.Papel {
		if p.Forma != "de contado" {
			continue
		}
		person, ok := who[p.Quien]
		if !ok {
			t.Fatalf("papel %d: persona %q desconocida", p.Paquete, p.Quien)
		}
		exp := ex(t, s, `INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, accrual_months, payment_method, paid_by_person_id)
			VALUES ('variable', ?, ?, ?, ?, 1, 'cash', ?)`, fmt.Sprintf("Papel (paquete %d)", p.Paquete), paper, pc(p.Monto), p.Compra, person)
		ex(t, s, `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, ?)`, exp, person, pc(p.Monto))
	}
	return m
}

// S7-02: cifras de cada PDF mensual = API (getMonthlyReport y getAccrualVsCash del mismo mes) de oct-2026 a mar-2027, y
// P1–P6 = R-45…R-50 del oráculo (S2 §6, nota (e): oct-2026 ↔ P1 … mar-2027 ↔ P6), con una categoría creada por la API.
func TestCruceMensual(t *testing.T) {
	e := setup(t, presupuesto.Register, ajustes.Register)
	userCategory(t, e)
	m := canon(t, e)
	ctx := context.Background()
	// API = oráculo en los 6 meses, a la fecha de la prueba (12-oct-2026).
	_, ms := accrualResp(t, e)
	if len(ms) < 6 {
		t.Fatalf("getAccrualVsCash trae %d meses, se esperaban al menos 6", len(ms))
	}
	type fila struct {
		r, month      string
		dev, efe, dif int64
	}
	var filas []fila
	for i := 0; i < 6; i++ {
		mm := ms[i].(map[string]any)
		o := m.Periodos[i]
		dev, efe := o["devengado"].(map[string]any), o["efectivo"].(map[string]any)
		f := fila{r: fmt.Sprintf("R-%d", 45+i), month: time.Date(2026, time.October+time.Month(i), 1, 12, 0, 0, 0, time.UTC).Format("2006-01"),
			dev: cents(dev["total"]), efe: cents(efe["total"]), dif: cents(o["diferencia_acumulada"])}
		got := []int64{cents(mm["accrued_cents"]), cents(mm["cash_cents"]), cents(mm["cumulative_difference_cents"])}
		if mm["month"] != f.month || fmt.Sprint(got) != fmt.Sprint([]int64{f.dev, f.efe, f.dif}) {
			t.Errorf("%s (%s, %v): API %v, oráculo %v", f.r, f.month, mm["month"], got, []int64{f.dev, f.efe, f.dif})
		}
		t.Logf("%s (%s, %s): devengado %s, efectivo %s, diferencia acumulada %s = oráculo", f.r, f.month, o["periodo"],
			reportes.Pesos(got[0]), reportes.Pesos(got[1]), reportes.Pesos(got[2]))
		filas = append(filas, f)
	}
	// PDF = API mes por mes. El reporte solo se genera de un mes que ya empezó: el reloj se pone en el día 1 de cada mes
	// (a las 12:00), antes de cualquier evento de la semilla (S2-09: ninguno cae entre los días 1 y 9).
	total, nf := 0, 0
	for i, f := range filas {
		mo, _ := dominio.ParseDate(f.month + "-01")
		if i > 0 {
			e.s.Clock.Set(testutil.MX(mo.Year, mo.Month, 1, 12, 0))
		}
		rep := generate(t, e, "/api/v1/reports/monthly", map[string]any{"month": f.month, "send_telegram": false})
		text := e.pdf(t, rep.ID)
		resp, _ := accrualResp(t, e)
		resp["monthly:"+f.month] = e.json(t, "/api/v1/monthly-report?period="+f.month)
		mr, err := reportes.BuildMonthly(ctx, e.s.DB, dominio.Today(e.s.Clock), dominio.MonthStart(mo.Year, mo.Month))
		if err != nil {
			t.Fatal(err)
		}
		n, fields := cruce(t, text, reportes.MonthlyDoc(mr), resp, "")
		total += n
		nf = len(fields)
		for _, c := range []int64{f.dev, f.efe, f.dif} {
			if !strings.Contains(text, reportes.Pesos(c)) {
				t.Errorf("el PDF de %s no trae %s (%s)", f.month, reportes.Pesos(c), f.r)
			}
		}
		t.Logf("S7-02 %s (%s): M=%d, PDF = API y trae %s / %s / %s", f.month, f.r, n, reportes.Pesos(f.dev), reportes.Pesos(f.efe), reportes.Pesos(f.dif))
		if i == len(filas)-1 {
			t.Logf("S7-02 campos cruzados por PDF (%d): %s", nf, strings.Join(fields, ", "))
		}
	}
	t.Logf("S7-02 total: %d cifras y fechas de 6 PDF (oct-2026…mar-2027) iguales a la API; R-45…R-50 = oráculo", total)
}

// S7-11: la vista web mensual (getMonthlyReport) trae las mismas cifras que el PDF del mismo mes.
func TestVistaMensual(t *testing.T) {
	e := setup(t, presupuesto.Register, ajustes.Register)
	userCategory(t, e)
	rep := generate(t, e, "/api/v1/reports/monthly", map[string]any{"month": "2026-10", "send_telegram": false})
	text := e.pdf(t, rep.ID)
	view := e.json(t, "/api/v1/monthly-report?period=2026-10")
	if pdf, _ := view["pdf"].(map[string]any); pdf == nil || cents(pdf["id"]) != rep.ID {
		t.Fatalf("la vista no enlaza el PDF %d: %v", rep.ID, view["pdf"])
	}
	resp, _ := accrualResp(t, e)
	resp["monthly:2026-10"] = view
	mo, err := reportes.BuildMonthly(context.Background(), e.s.DB, dominio.Today(e.s.Clock), dominio.MonthStart(2026, 10))
	if err != nil {
		t.Fatal(err)
	}
	m, _ := cruce(t, text, reportes.MonthlyDoc(mo), resp, "monthly:")
	t.Logf("S7-11 vista web = PDF: %d cifras de getMonthlyReport, 0 diferencias", m)
	if d := e.json(t, "/api/v1/monthly-report"); d["month"] != "2026-09" {
		t.Errorf("por defecto debe ser el último mes cerrado: %v", d["month"])
	}
	if res := e.c.Do("GET", "/api/v1/monthly-report?period=2026-W41", nil); res.StatusCode != 422 && res.StatusCode != 400 {
		t.Errorf("periodo semanal: %d", res.StatusCode)
	}
}

// S7-03: estado de cuenta de Karen bajo demanda, solo su despensa en común y sus pagos; no se le envía.
func TestKaren(t *testing.T) {
	e := setup(t, presupuesto.Register)
	path := fmt.Sprintf("/api/v1/persons/%d/statement?period=2026-10", e.karen)
	resp := map[string]map[string]any{"statement": e.json(t, path)}
	for _, m := range []string{"2026-10", "2026-11", "2026-12"} {
		resp["monthly:"+m] = e.json(t, "/api/v1/monthly-report?period="+m)
	}
	rep := generate(t, e, "/api/v1/reports/external-statements", map[string]any{"person_id": e.karen, "period_start_date": "2026-10-01", "period_end_date": "2026-10-31"})
	if rep.Kind != "external_statement_pdf" || rep.PersonID == nil || *rep.PersonID != e.karen {
		t.Fatalf("%+v", rep)
	}
	text := e.pdf(t, rep.ID)
	x, err := reportes.BuildExternal(context.Background(), e.s.DB, dominio.Today(e.s.Clock), e.karen, dominio.MonthStart(2026, 10), dominio.NewDate(2026, 10, 31))
	if err != nil {
		t.Fatal(err)
	}
	m, fields := cruce(t, text, reportes.ExternalDoc(x), resp, "")
	t.Logf("S7-03 M=%d campos: %s", m, strings.Join(fields, ", "))
	for _, s := range []string{"$1,499.72", "$308.22", "$1,191.50", "Vie 30 oct 2026", "Dom 15 nov 2026", "Lun 30 nov 2026", "Mar 15 dic 2026"} {
		if !strings.Contains(text, s) {
			t.Errorf("falta %q", s)
		}
	}
	if n := strings.Count(text, "$374.93"); n != 4 {
		t.Errorf("$374.93 aparece %d veces, se esperaban 4", n)
	}
	for _, s := range []string{"$12,345.67", "Ingreso", "$1,903.00", "$598.82"} {
		if strings.Contains(text, s) {
			t.Errorf("el estado de cuenta de Karen no debe traer %q", s)
		}
	}
	res := e.c.Do("POST", fmt.Sprintf("/api/v1/reports/%d/send", rep.ID), nil, "Idempotency-Key", key())
	if res.StatusCode != 409 {
		t.Errorf("el estado de cuenta no se envía: %d", res.StatusCode)
	}
	if res := e.c.Do("POST", "/api/v1/reports/external-statements", map[string]any{"person_id": e.axel, "period_start_date": "2026-10-01",
		"period_end_date": "2026-10-31"}, "Idempotency-Key", key()); res.StatusCode != 422 {
		t.Errorf("Axel no es persona externa: %d", res.StatusCode)
	}
}

type fakeTelegram struct {
	mu    sync.Mutex
	calls []map[string]string
	srv   *httptest.Server
}

func newFakeTelegram(t *testing.T) *fakeTelegram {
	f := &fakeTelegram{}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		mr := multipart.NewReader(r.Body, params["boundary"])
		call := map[string]string{"method": filepath.Base(r.URL.Path)}
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}
			b, _ := io.ReadAll(p)
			if p.FormName() == "document" {
				call["document_bytes"] = strconv.Itoa(len(b))
				continue
			}
			call[p.FormName()] = string(b)
		}
		f.mu.Lock()
		f.calls = append(f.calls, call)
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1}}`))
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeTelegram) n() int { f.mu.Lock(); defer f.mu.Unlock(); return len(f.calls) }

// SendDocument imita al cliente de internal/telegram contra el servidor falso (token falso de prueba).
func (f *fakeTelegram) SendDocument(ctx context.Context, chatID, name string, pdf []byte, caption string) (int, error) {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	_ = mw.WriteField("chat_id", chatID)
	_ = mw.WriteField("caption", caption)
	fw, _ := mw.CreateFormFile("document", name)
	_, _ = fw.Write(pdf)
	_ = mw.Close()
	req, _ := http.NewRequestWithContext(ctx, "POST", f.srv.URL+"/botTOKEN-FALSO-DE-PRUEBA/sendDocument", &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	res, err := f.srv.Client().Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

func jobRow(t *testing.T, e *env) (schedule, next, status string) {
	t.Helper()
	var n, st any
	if err := e.s.DB.QueryRow(`SELECT schedule, next_run_at, last_status FROM job WHERE code = ?`, reportes.MonthlyJobCode).Scan(&schedule, &n, &st); err != nil {
		t.Fatal(err)
	}
	return schedule, fmt.Sprint(n), fmt.Sprint(st)
}

func count(t *testing.T, e *env, q string, args ...any) int {
	var n int
	if err := e.s.DB.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// S7-10: el job del día configurado genera, guarda versionado y envía; con LaTeX roto, error y 0 envíos; la hora y el
// día salen de reports.monthly_send_time y reports.monthly_day_of_month (31 → 30-nov-2026).
func TestMensual(t *testing.T) {
	e := setup(t, presupuesto.Register, ajustes.Register)
	userCategory(t, e)
	ctx := context.Background()
	tg := newFakeTelegram(t)
	reportes.SetSender(tg)
	t.Cleanup(func() { reportes.SetSender(nil) })
	nid := ex(t, e.s, `INSERT INTO notification (code, name) VALUES ('monthly_report', 'Reporte mensual') ON CONFLICT (code) DO UPDATE SET is_enabled = 1`)
	for _, r := range []struct {
		p    int64
		chat string
	}{{e.axel, "1001"}, {e.ximena, "1002"}, {e.karen, "1003"}} {
		rid := ex(t, e.s, `INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, ?)`, r.p, r.chat)
		ex(t, e.s, `INSERT INTO notification_subscription (notification_id, recipient_id) VALUES (?, ?)`, nid, rid)
	}
	sch := e.s.Reg.Deps().Jobs
	e.s.Clock.Set(testutil.MX(2026, time.November, 1, 7, 59))
	if err := sch.Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := sch.Tick(ctx); err != nil {
		t.Fatal(err)
	}
	if s, next, _ := jobRow(t, e); s != "0 8 1 * *" || !strings.HasPrefix(next, "2026-11-01T14:00") {
		t.Fatalf("horario inicial %q %q", s, next)
	}
	e.s.Clock.Set(testutil.MX(2026, time.November, 1, 8, 0))
	if n, err := sch.Tick(ctx); err != nil || n == 0 {
		t.Fatalf("no corrió: %d %v", n, err)
	}
	if _, _, st := jobRow(t, e); st != "ok" {
		var msg any
		_ = e.s.DB.QueryRow(`SELECT last_error FROM job WHERE code = ?`, reportes.MonthlyJobCode).Scan(&msg)
		t.Fatalf("estado %s: %v", st, msg)
	}
	var path string
	var size int64
	if err := e.s.DB.QueryRow(`SELECT file_path, size_bytes FROM generated_report WHERE kind = 'monthly_pdf' AND period_start_date = '2026-10-01'`).Scan(&path, &size); err != nil {
		t.Fatal(err)
	}
	full := filepath.Join(e.s.Reg.Deps().Config.DataDir, path)
	if st, err := os.Stat(full); err != nil || st.Size() != size || size > reportes.MaxTelegramDocumentBytes {
		t.Fatalf("PDF guardado %s: %v", full, err)
	}
	t.Logf("S7-10 PDF guardado en %s (%d bytes; límite de sendDocument %d bytes)", path, size, reportes.MaxTelegramDocumentBytes)
	if tg.n() != 2 {
		t.Fatalf("envíos %d, se esperaban 2 (Axel y Ximena; Karen nunca)", tg.n())
	}
	for _, c := range tg.calls {
		if c["method"] != "sendDocument" || c["chat_id"] == "1003" || c["parse_mode"] != "" || c["document_bytes"] != strconv.FormatInt(size, 10) {
			t.Fatalf("envío %v", c)
		}
	}
	if n := count(t, e, `SELECT COUNT(*) FROM notification_delivery WHERE status = 'sent' AND http_status = 200`); n != 2 {
		t.Fatalf("entregas %d", n)
	}
	// Plantilla rota: el job queda en error, sin PDF nuevo y sin envíos.
	restore := reportes.SetTemplates(fstest.MapFS{
		"plantillas/comun.tex":   {Data: []byte(`<<define "preambulo">>\documentclass{article}<<end>><<define "cuerpo">>\begin{document}\comandoquenoexiste<<end>>`)},
		"plantillas/mensual.tex": {Data: []byte(`<<template "preambulo" .>><<template "cuerpo" .>>`)},
	})
	before := count(t, e, `SELECT COUNT(*) FROM generated_report`)
	if err := sch.RunNow(ctx, reportes.MonthlyJobCode); err != nil {
		t.Fatal(err)
	}
	if _, err := sch.Tick(ctx); err != nil {
		t.Fatal(err)
	}
	restore()
	if _, _, st := jobRow(t, e); st != "error" || count(t, e, `SELECT COUNT(*) FROM generated_report`) != before || tg.n() != 2 {
		t.Fatalf("con LaTeX roto: estado %s, reportes %d→%d, envíos %d", st, before, count(t, e, `SELECT COUNT(*) FROM generated_report`), tg.n())
	}
	t.Log("S7-10 plantilla rota: job en error, 0 PDF nuevos, 0 envíos")
	// Hora nueva desde Ajustes.
	res := e.c.Do("PATCH", "/api/v1/settings/reports", map[string]any{"parameters": map[string]any{"reports.monthly_send_time": "09:30"}})
	if res.StatusCode != 200 {
		t.Fatalf("updateReportSettings %d %s", res.StatusCode, res.Body)
	}
	if s, next, _ := jobRow(t, e); s != "30 9 1 * *" || !strings.HasPrefix(next, "2026-12-01T15:30") {
		t.Fatalf("hora nueva: %q %q", s, next)
	}
	var pj string
	if err := e.s.DB.QueryRow(`SELECT params_json FROM notification WHERE code = 'monthly_report'`).Scan(&pj); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"send_time", "day_of_month"} {
		if strings.Contains(pj, k) {
			t.Fatalf("params_json no debe traer %s: %s", k, pj)
		}
	}
	// Día 31: en noviembre sale el 30.
	res = e.c.Do("PATCH", "/api/v1/settings/reports", map[string]any{"parameters": map[string]any{"reports.monthly_day_of_month": 31}})
	if res.StatusCode != 200 {
		t.Fatalf("updateReportSettings %d %s", res.StatusCode, res.Body)
	}
	if s, _, _ := jobRow(t, e); s != "30 9 28-31 * *" {
		t.Fatalf("horario día 31: %q", s)
	}
	n0 := count(t, e, `SELECT COUNT(*) FROM generated_report WHERE kind = 'monthly_pdf'`)
	for _, d := range []int{28, 29} {
		e.s.Clock.Set(testutil.MX(2026, time.November, d, 9, 30))
		if _, err := sch.Tick(ctx); err != nil {
			t.Fatal(err)
		}
		if count(t, e, `SELECT COUNT(*) FROM generated_report WHERE kind = 'monthly_pdf'`) != n0 {
			t.Fatalf("salió el %d-nov", d)
		}
	}
	e.s.Clock.Set(testutil.MX(2026, time.November, 30, 9, 30))
	if _, err := sch.Tick(ctx); err != nil {
		t.Fatal(err)
	}
	if count(t, e, `SELECT COUNT(*) FROM generated_report WHERE kind = 'monthly_pdf' AND file_path LIKE '%-v2.pdf'`) != 1 {
		t.Fatal("con día 31 no salió el 30-nov-2026 (versión 2 de octubre)")
	}
	t.Log("S7-10 día 31 → 30-nov-2026; hora 09:30 respetada; versión 2 guardada sin pisar la 1")
}

// S7-12: dinero $1,234.56, fechas 10-oct-2026 / Vie 30 oct 2026, sin fechas ISO ni inglés en el texto visible.
func TestFormatoS712(t *testing.T) {
	e := setup(t, presupuesto.Register)
	r1 := generate(t, e, "/api/v1/reports/monthly", map[string]any{"month": "2026-10", "send_telegram": false})
	r2 := generate(t, e, "/api/v1/reports/external-statements", map[string]any{"person_id": e.karen, "period_start_date": "2026-10-01", "period_end_date": "2026-10-31"})
	texts := []string{e.pdf(t, r1.ID), e.pdf(t, r2.ID), "Reporte mensual de " + reportes.MesLargo(2026, time.October) + " (PDF)."}
	reMoney := regexp.MustCompile(`^\$-?[0-9]{1,3}(,[0-9]{3})*\.[0-9]{2}$`)
	reISO := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	reEN := regexp.MustCompile(`(?i)\b(page|total due|balance|monday|tuesday|wednesday|thursday|friday|saturday|sunday|january|february|october|november|december|report|statement|paid|pending|overdue)\b`)
	nm := 0
	for _, s := range texts {
		for _, m := range regexp.MustCompile(`-?\$\S+`).FindAllString(s, -1) {
			m = strings.TrimRight(m, ".,;")
			nm++
			if !reMoney.MatchString(strings.TrimPrefix(m, "-")) {
				t.Errorf("cifra mal formada %q", m)
			}
		}
		if f := reISO.FindAllString(s, -1); len(f) > 0 {
			t.Errorf("fechas ISO visibles: %v", f)
		}
		if f := reEN.FindAllString(s, -1); len(f) > 0 {
			t.Errorf("inglés visible: %v", f)
		}
	}
	lines := strings.Split(strings.TrimSpace(texts[0]), "\n")
	if len(lines) > 10 {
		lines = lines[:10]
	}
	t.Logf("S7-12 %d cifras con formato $1,234.56, 0 fechas ISO, 0 palabras en inglés. Primeras líneas:\n%s", nm, strings.Join(lines, "\n"))
}

// S7-13: cada \definecolor de las plantillas es un token de C2 y el texto tiene contraste ≥ 4.5:1 sobre su fondo.
func TestColoresC2(t *testing.T) {
	root := testutil.Root(t)
	b, err := os.ReadFile(filepath.Join(root, "docs", "contratos", "C2-tokens.json"))
	if err != nil {
		t.Fatal(err)
	}
	var tok map[string]any
	if err := json.Unmarshal(b, &tok); err != nil {
		t.Fatal(err)
	}
	tex, err := os.ReadFile(filepath.Join(root, "backend", "internal", "reportes", "plantillas", "comun.tex"))
	if err != nil {
		t.Fatal(err)
	}
	re := regexp.MustCompile(`\\definecolor\{(\w+)\}\{HTML\}\{([0-9A-F]{6})\}% ([a-z0-9.]+)`)
	hex := map[string]string{}
	for _, m := range re.FindAllStringSubmatch(string(tex), -1) {
		var cur any = tok
		for _, p := range strings.Split(m[3], ".") {
			cur = cur.(map[string]any)[p]
		}
		v, _ := cur.(map[string]any)["$value"].(string)
		if !strings.EqualFold(v, "#"+m[2]) {
			t.Errorf("%s = #%s pero %s = %s en C2", m[1], m[2], m[3], v)
		}
		hex[m[1]] = m[2]
		t.Logf("%s #%s = %s", m[1], m[2], m[3])
	}
	if len(hex) != 7 {
		t.Fatalf("se esperaban 7 colores, hay %d", len(hex))
	}
	lum := func(h string) float64 {
		var c [3]float64
		for i := 0; i < 3; i++ {
			n, _ := strconv.ParseUint(h[2*i:2*i+2], 16, 8)
			v := float64(n) / 255
			if v <= 0.04045 {
				c[i] = v / 12.92
			} else {
				c[i] = math.Pow((v+0.055)/1.055, 2.4)
			}
		}
		return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2]
	}
	for _, p := range [][2]string{{"fhText", "fhPage"}, {"fhMuted", "fhPage"}, {"fhBrand", "fhPage"}, {"fhText", "fhBrandSurface"},
		{"fhBrand", "fhBrandSurface"}, {"fhText", "fhHead"}} {
		a, b := lum(hex[p[0]]), lum(hex[p[1]])
		if a < b {
			a, b = b, a
		}
		r := (a + 0.05) / (b + 0.05)
		if r < 4.5 {
			t.Errorf("%s sobre %s: %.2f:1", p[0], p[1], r)
		}
		t.Logf("contraste %s sobre %s = %.2f:1", p[0], p[1], r)
	}
}

// S7-01: la plantilla compila con lualatex sin errores; con FH_S7_BUILD deja el .tex, el .log y el .pdf ahí.
func TestPlantillasCompilan(t *testing.T) {
	e := setup(t, presupuesto.Register)
	ctx := context.Background()
	c := &reportes.Compiler{TmpDir: testutil.TempDir(t), KeepDir: os.Getenv("FH_S7_BUILD")}
	mo, err := reportes.BuildMonthly(ctx, e.s.DB, dominio.Today(e.s.Clock), dominio.MonthStart(2026, 10))
	if err != nil {
		t.Fatal(err)
	}
	x, err := reportes.BuildExternal(ctx, e.s.DB, dominio.Today(e.s.Clock), e.karen, dominio.MonthStart(2026, 10), dominio.NewDate(2026, 10, 31))
	if err != nil {
		t.Fatal(err)
	}
	for name, d := range map[string]*reportes.Doc{"mensual": reportes.MonthlyDoc(mo), "estado-cuenta": reportes.ExternalDoc(x)} {
		pdf, err := c.Compile(ctx, name, d)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		t.Logf("S7-01 %s.pdf: %d bytes", name, len(pdf))
	}
	if !strings.Contains(reportes.Tex(`50% & $5 #1 _a ¿sí? ¡ya!`), `50\% \& \$5 \#1 \_a \textquestiondown{}sí? \textexclamdown{}ya!`) {
		t.Errorf("escape %q", reportes.Tex(`50% & $5 #1 _a ¿sí? ¡ya!`))
	}
}

func tail(s string) string {
	if len(s) > 700 {
		return s[len(s)-700:]
	}
	return s
}
