package avisos

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"finanzas-hogar/internal/api"
	apigastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
	"finanzas-hogar/internal/testutil"
)

func cambiar(t *testing.T, db *sql.DB, set map[string]string, reset []string) {
	t.Helper()
	ch := dominio.ParamChange{Set: map[string]json.RawMessage{}, Reset: reset}
	for k, v := range set {
		ch.Set[k] = json.RawMessage(v)
	}
	keys := []string{"reports.daily_send_time", "reports.weekly_send_time", "reports.weekly_weekday"}
	if err := dominio.WithTx(context.Background(), db, func(tx *sql.Tx) error { return dominio.ApplyParams(context.Background(), tx, ch, keys, nil, "prueba") }); err != nil {
		t.Fatal(err)
	}
}

// TestResumen (S7-04): diario y semanal por el planificador, una sola vez, a la hora y el día de app_setting,
// con las cifras de getExpenseTotals.
func TestResumen(t *testing.T) {
	ctx := context.Background()
	s := testutil.NewServer(t, apigastos.Register)
	destinos(t, s.DB, s.Fx.AxelPersonID, s.Fx.XimenaPersonID)
	res, err := s.DB.Exec(`INSERT INTO category (name) VALUES ('Prueba de avisos')`)
	if err != nil {
		t.Fatal(err)
	}
	cat, _ := res.LastInsertId()
	for _, e := range []struct {
		d string
		c int64
	}{{"2026-10-12", 12345}, {"2026-10-12", 67890}, {"2026-10-03", 100000}, {"2026-10-18", 5000}} {
		if _, err := s.DB.Exec(`INSERT INTO expense (kind, category_id, amount_cents, expense_date, payment_method) VALUES ('variable', ?, ?, ?, 'cash')`, cat, e.c, e.d); err != nil {
			t.Fatal(err)
		}
	}
	b := &buzon{}
	m := &Motor{DB: s.DB, Fuente: FuenteBD{s.DB}, Enviador: b}
	sch := jobs.New(s.DB, s.Clock, nil)
	for _, d := range m.Jobs() {
		if err := sch.Register(d); err != nil {
			t.Fatal(err)
		}
	}
	if err := sch.Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	at := func(dia string, hh, mm int) {
		s.Clock.Set(a9(fecha(t, dia), hh, mm))
		if _, err := sch.Tick(ctx); err != nil {
			t.Fatal(err)
		}
	}
	diario := func() int { return len(b.con("RESUMEN DEL DÍA")) }
	semanal := func() []enviado { return b.con("RESUMEN DE LA SEMANA") }
	c := s.SessionFor(t, testutil.AxelUser)
	tot := func(period string) string {
		var x api.ExpenseTotals
		c.Do("GET", "/api/v1/expenses/totals?period="+period, nil).JSON(t, &x)
		return gastosTxt(x.ExpenseCount) + " por " + pesos(x.AmountCents)
	}

	at("2026-10-12", 21, 0)
	at("2026-10-12", 21, 29)
	if diario() != 0 {
		t.Fatal("antes de reports.daily_send_time no sale")
	}
	at("2026-10-12", 21, 30)
	got := b.con("RESUMEN DEL DÍA")
	if len(got) != 2 {
		t.Fatalf("21:30: %d", len(got))
	}
	for _, w := range []string{"Gastos de hoy: " + tot("2026-10-12..2026-10-12"), "Gastos del mes (oct 2026): " + tot("2026-10")} {
		if !strings.Contains(got[0].text, w) {
			t.Fatalf("el texto no coincide con la API: falta %q en\n%s", w, got[0].text)
		}
	}
	at("2026-10-12", 21, 31)
	at("2026-10-12", 23, 59)
	if diario() != 0 {
		t.Fatal("el diario sale una sola vez")
	}
	t.Logf("S7-04 diario: 21:29 → 0; 21:30 → 2 (Axel y Ximena); 21:31 y 23:59 → 0; cifras = getExpenseTotals\n%s", got[0].text)

	cambiar(t, s.DB, map[string]string{"reports.daily_send_time": `"22:15"`}, nil)
	at("2026-10-13", 21, 30)
	if diario() != 0 {
		t.Fatal("con 22:15 no sale a las 21:30")
	}
	at("2026-10-13", 22, 15)
	if diario() != 2 {
		t.Fatal("con 22:15 sale a las 22:15")
	}
	cambiar(t, s.DB, nil, []string{"reports.daily_send_time"})
	at("2026-10-14", 21, 30)
	if diario() != 2 {
		t.Fatal("restablecido sale a las 21:30")
	}
	t.Logf("S7-04 diario con reports.daily_send_time = 22:15: 21:30 → 0, 22:15 → 2; restablecido → 21:30")

	at("2026-10-17", 19, 0) // sábado
	at("2026-10-18", 18, 59)
	if len(semanal()) != 0 {
		t.Fatal("el semanal solo el día y a la hora")
	}
	at("2026-10-18", 19, 0) // domingo
	sg := semanal()
	if len(sg) != 2 || !strings.Contains(sg[0].text, "Gastos de la semana: "+tot("2026-W42")) || !strings.Contains(sg[0].text, "Gastos del mes (oct 2026): "+tot("2026-10")) {
		t.Fatalf("semanal: %+v", sg)
	}
	at("2026-10-18", 19, 5)
	if len(semanal()) != 0 {
		t.Fatal("el semanal sale una vez")
	}
	t.Logf("S7-04 semanal: sábado → 0; domingo 18:59 → 0; 19:00 → 2; 19:05 → 0; cifras = getExpenseTotals (2026-W42)\n%s", sg[0].text)
	cambiar(t, s.DB, map[string]string{"reports.weekly_weekday": "3", "reports.weekly_send_time": `"20:00"`}, nil)
	at("2026-10-21", 19, 0)
	if len(semanal()) != 0 {
		t.Fatal("miércoles 19:00 con 20:00 no sale")
	}
	at("2026-10-21", 20, 0)
	at("2026-10-25", 20, 0)
	if sg := semanal(); len(sg) != 2 || !strings.Contains(sg[0].text, "19-oct-2026 al 25-oct-2026") {
		t.Fatalf("con día 3 sale el miércoles y no el domingo: %d", len(sg))
	}
	cambiar(t, s.DB, nil, []string{"reports.weekly_weekday", "reports.weekly_send_time"})
	at("2026-10-28", 19, 0)
	at("2026-11-01", 19, 0)
	if sg := semanal(); len(sg) != 2 || !strings.Contains(sg[0].text, "26-oct-2026 al 1-nov-2026") {
		t.Fatal("restablecido vuelve al domingo 19:00")
	}
	var pd, pw string
	s.DB.QueryRow(`SELECT params_json FROM notification WHERE code = 'daily_summary'`).Scan(&pd)
	s.DB.QueryRow(`SELECT params_json FROM notification WHERE code = 'weekly_summary'`).Scan(&pw)
	if pd != "{}" || pw != "{}" {
		t.Fatalf("params_json: %s %s", pd, pw)
	}
	var n int
	s.DB.QueryRow(`SELECT count(*) FROM notification_delivery d JOIN notification x ON x.id = d.notification_id WHERE x.code IN ('daily_summary', 'weekly_summary')`).Scan(&n)
	t.Logf("S7-04 semanal con weekday 3 y 20:00: miércoles 20:00 → 2, domingo → 0; restablecido (7, 19:00) → domingo; params_json = {}; %d envíos registrados sin repetir", n)
}

func gastosTxt(n int64) string { return gastos(n) }

var (
	reDinero = regexp.MustCompile(`\$[-0-9.,]*`)
	reOkMXN  = regexp.MustCompile(`^\$-?[0-9]{1,3}(,[0-9]{3})*\.[0-9]{2}$`)
	reISO    = regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
)

// TestFormato (S7-12): dinero $1,234.56, fechas 10-oct-2026 / Vie 30 oct 2026, sin fechas ISO, encabezado en mayúsculas.
func TestFormato(t *testing.T) {
	ctx := context.Background()
	db, fx := testutil.NewDB(t)
	destinos(t, db, fx.AxelPersonID, fx.XimenaPersonID)
	f := &fuenteFija{FuenteBD: FuenteBD{db}, pagos: []Pago{{42, fecha(t, "2026-10-30"), 97375, "Kueski"}, {44, fecha(t, "2026-10-15"), 123456789, "Tarjeta"}}}
	m := &Motor{DB: db, Fuente: f, Enviador: &buzon{}}
	hoy := fecha(t, "2026-10-12")
	var textos []string
	d, _ := m.TextoDiario(ctx, hoy)
	w, _ := m.TextoSemanal(ctx, hoy)
	textos = append(textos, d, w, TextoPago(f.pagos[0], 1), TextoPago(f.pagos[1], 0),
		TextoPresupuesto(Uso{"total", "Total del mes", 500000, 400050}, 80, hoy), TextoPrecios(fecha(t, "2026-10-17"), 3), TextoPlan(fecha(t, "2026-10-19"), 1, 21))
	for _, c := range []string{"daily_summary", "weekly_summary", "monthly_report", "payment_due", "budget_threshold", "stale_prices", "meal_plan_missing"} {
		x, err := m.TextoPrueba(ctx, c, hoy)
		if err != nil {
			t.Fatal(err)
		}
		textos = append(textos, x)
	}
	cifras, lineas := 0, 0
	for _, x := range textos {
		ls := strings.Split(x, "\n")
		if ls[0] != strings.ToUpper(ls[0]) {
			t.Fatalf("encabezado no está en mayúsculas: %q", ls[0])
		}
		for _, l := range ls {
			lineas++
			if reISO.MatchString(l) {
				t.Fatalf("fecha ISO visible: %q", l)
			}
			for _, mny := range reDinero.FindAllString(l, -1) {
				cifras++
				if !reOkMXN.MatchString(strings.TrimRight(mny, ".,")) {
					t.Fatalf("dinero mal formado %q en %q", mny, l)
				}
			}
		}
	}
	if !strings.Contains(textos[2], "Vie 30 oct 2026") || !strings.Contains(textos[3], "$1,234,567.89") {
		t.Fatal(textos[2], textos[3])
	}
	t.Logf("S7-12 %d textos, %d líneas, %d cifras con ^\\$-?[0-9]{1,3}(,[0-9]{3})*\\.[0-9]{2}$, 0 fechas ISO. Muestra:\n%s\n%s\n%s\n%s",
		len(textos), lineas, cifras, textos[0], textos[2], textos[4], textos[6])
}
