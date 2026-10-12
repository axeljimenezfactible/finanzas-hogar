package comparador

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	mc "finanzas-hogar/internal/motor_cred"
	ct "finanzas-hogar/internal/motor_cred/canontest"
)

type fixture struct {
	Entradas struct {
		Datos     json.RawMessage `json:"datos"`
		Consultas []struct {
			ID         string          `json:"id"`
			Argumentos json.RawMessage `json:"argumentos"`
		} `json:"consultas"`
	} `json:"entradas"`
	Esperado map[string]json.RawMessage `json:"esperado"`
}

// fixtureDir es oraculo/fixtures/s3, la única fuente de la tabla a mano (S3-interfaz.md §10.0;
// corrige REVISION-S3-r1 H3: ya no hay copias en testdata/).
func fixtureDir(t *testing.T) string {
	t.Helper()
	root, err := ct.Root()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(root, "oraculo", "fixtures", "s3")
}

// readFixture lee oraculo/fixtures/s3/<file> y resuelve entradas.datos {"$ref": …} (sin parche:
// las fixtures del comparador no lo usan; si apareciera, falla en vez de ignorarlo).
func readFixture(t *testing.T, file string) fixture {
	t.Helper()
	dir := fixtureDir(t)
	b, err := os.ReadFile(filepath.Join(dir, file))
	if err != nil {
		t.Fatalf("falta oraculo/fixtures/s3/%s (autor: optimizador S3): %v", file, err)
	}
	var f fixture
	if err := json.Unmarshal(b, &f); err != nil {
		t.Fatalf("%s: %v", file, err)
	}
	var ref struct {
		Ref    *string         `json:"$ref"`
		Parche json.RawMessage `json:"parche"`
	}
	if json.Unmarshal(f.Entradas.Datos, &ref) == nil && ref.Ref != nil {
		if len(ref.Parche) > 0 {
			t.Fatalf("%s: parche no soportado en las fixtures del comparador", file)
		}
		if f.Entradas.Datos, err = os.ReadFile(filepath.Join(dir, filepath.FromSlash(*ref.Ref))); err != nil {
			t.Fatalf("falta oraculo/fixtures/s3/%s (autor: optimizador S3): %v", *ref.Ref, err)
		}
	}
	return f
}

// loadFixture lee oraculo/fixtures/s3/s3-comparador-<name>.json: datos, argumentos de la
// consulta y esperado genérico.
func loadFixture(t *testing.T, name string) (mc.Input, json.RawMessage, any) {
	t.Helper()
	f := readFixture(t, "s3-comparador-"+name+".json")
	in, err := mc.DecodeInput(bytes.NewReader(f.Entradas.Datos))
	if err != nil {
		t.Fatalf("%s: datos: %v", name, err)
	}
	q := f.Entradas.Consultas[0]
	return in, q.Argumentos, generic(t, f.Esperado[q.ID])
}

// TestOraculoFixtures (§11, «ambos»): las 6 fixtures del comparador existen y cada consulta
// reproduce su esperado (subconjunto de claves, listas completas) por la vía de las fixtures.
func TestOraculoFixtures(t *testing.T) {
	names := []string{"caso-1", "caso-2", "caso-3", "sin-restricciones", "tiempo", "carga"}
	for _, name := range names {
		f := readFixture(t, "s3-comparador-"+name+".json")
		in, err := mc.DecodeInput(bytes.NewReader(f.Entradas.Datos))
		if err != nil {
			t.Fatalf("%s: datos: %v", name, err)
		}
		for _, q := range f.Entradas.Consultas {
			out, code := EvalComparar(in, q.Argumentos)
			var got any = map[string]any{"error": code}
			if code == "" {
				got = generic(t, out)
			}
			var diffs []string
			subset(name+":"+q.ID, generic(t, f.Esperado[q.ID]), got, &diffs)
			for _, d := range diffs {
				t.Errorf("esperado ≠ motor: %s", d)
			}
		}
	}
	t.Logf("%d fixtures del comparador corridas", len(names))
}

func generic(t *testing.T, v any) any {
	t.Helper()
	b, ok := v.(json.RawMessage)
	if !ok {
		var err error
		if b, err = json.Marshal(v); err != nil {
			t.Fatal(err)
		}
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var out any
	if err := dec.Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

// subset: cada clave del esperado está en la salida con el mismo valor; listas completas.
func subset(path string, exp, got any, diffs *[]string) {
	switch e := exp.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("%s: esperado objeto, salida %v", path, got))
			return
		}
		for k, v := range e {
			gv, ok := g[k]
			if !ok {
				*diffs = append(*diffs, path+"."+k+": falta en la salida")
				continue
			}
			subset(path+"."+k, v, gv, diffs)
		}
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(e) {
			*diffs = append(*diffs, fmt.Sprintf("%s: lista de %d, salida %v", path, len(e), got))
			return
		}
		for i := range e {
			subset(fmt.Sprintf("%s[%d]", path, i), e[i], g[i], diffs)
		}
	default:
		if !reflect.DeepEqual(exp, got) {
			*diffs = append(*diffs, fmt.Sprintf("%s: esperado %v, salida %v", path, exp, got))
		}
	}
}

func requestOf(t *testing.T, in mc.Input, args json.RawMessage) (Request, int64) {
	t.Helper()
	var a struct {
		TimeoutSeconds int64   `json:"timeout_seconds"`
		Request        Request `json:"request"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		t.Fatal(err)
	}
	a.Request.Today, a.Request.Credits = in.Today, in.Credits
	return a.Request, a.TimeoutSeconds
}

func compareT(t *testing.T, req Request) Result {
	t.Helper()
	res, err := Compare(context.Background(), req)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	return res
}

var cases = []string{"caso-1", "caso-2", "caso-3"}

func keyOf(s Scenario) string {
	f := func(p *int64) string {
		if p == nil {
			return "-"
		}
		return strconv.FormatInt(*p, 10)
	}
	return fmt.Sprintf("%d/%dq/%s/t%s/p%s", s.CreditID, s.TermCount, s.PurchaseDate, f(s.StoreID), f(s.PromotionID))
}

// checkConstraints: 100 % de los válidos cumplen las 3 restricciones y cada descartado viola la
// que dice y ninguna anterior en el orden de C1.
func checkConstraints(t *testing.T, label string, res Result) {
	t.Helper()
	u := res.ConstraintsUsed
	for _, s := range res.Valid {
		if s.MinAvailable < *u.MinAvailableBuffer || s.MaxLivePurchases > *u.MaxLivePurchases || s.ZeroDebtDate.Compare(*u.ZeroDebtBy) > 0 {
			t.Errorf("%s: válido %s viola una restricción", label, keyOf(s))
		}
		if s.ScoreBP == nil || *s.ScoreBP < 0 || *s.ScoreBP > 10000 {
			t.Errorf("%s: puntaje fuera de [0, 10000] en %s", label, keyOf(s))
		}
	}
	for _, r := range res.Rejected {
		if v := violation(r.Scenario, u); v != r.ViolatedConstraint || v == "" || r.Scenario.ScoreBP != nil {
			t.Errorf("%s: descartado %s con %q (recalculado %q)", label, keyOf(r.Scenario), r.ViolatedConstraint, v)
		}
	}
	if res.EvaluatedCount != len(res.Valid)+len(res.Rejected) {
		t.Errorf("%s: evaluated_count %d", label, res.EvaluatedCount)
	}
}

func TestRestricciones(t *testing.T) {
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			in, args, exp := loadFixture(t, name)
			out, code := EvalComparar(in, args)
			if code != "" {
				t.Fatalf("código %q", code)
			}
			var diffs []string
			subset(name, exp, generic(t, out), &diffs)
			for i, d := range diffs {
				if i == 15 {
					t.Errorf("… %d diferencias en total", len(diffs))
					break
				}
				t.Errorf("tabla a mano ≠ motor: %s", d)
			}
			res := out.(Result)
			checkConstraints(t, name, res)
			seen := map[string]bool{}
			for _, r := range res.Rejected {
				seen[r.ViolatedConstraint] = true
			}
			if len(seen) != 3 {
				t.Errorf("descartes por restricción: %v (se esperan las 3)", seen)
			}
			t.Logf("S3-12 %s: %d evaluados, %d válidos, %d descartados, 0 válidos con violación", name, res.EvaluatedCount, len(res.Valid), len(res.Rejected))
		})
	}
	// el cargo normal de la tarjeta siempre se genera (§8.2): caso 2, cada fecha × tienda
	in, args, _ := loadFixture(t, "caso-2")
	req, _ := requestOf(t, in, args)
	res := compareT(t, req)
	normal := map[string]bool{}
	for _, s := range append(res.Valid, rejectedScenarios(res)...) {
		if s.CreditID == 3 && s.TermCount == 1 && s.PromotionID == nil {
			normal[fmt.Sprintf("%s/%d", s.PurchaseDate, *s.StoreID)] = true
		}
	}
	if len(normal) != 3*len(req.Baskets) {
		t.Errorf("cargos normales generados: %d, se esperan %d", len(normal), 3*len(req.Baskets))
	}
}

func rejectedScenarios(res Result) []Scenario {
	out := []Scenario{}
	for _, r := range res.Rejected {
		out = append(out, r.Scenario)
	}
	return out
}

// TestRestriccionesPropiedad: 300 combinaciones al azar (semilla fija) de las 4 restricciones
// sobre los datos de los 3 casos; en todas, 100 % de los válidos cumplen.
func TestRestriccionesPropiedad(t *testing.T) {
	rng := rand.New(rand.NewSource(20261011))
	total := 0
	for i := 0; i < 300; i++ {
		in, args, _ := loadFixture(t, cases[i%3])
		req, _ := requestOf(t, in, args)
		buf := mc.Cents(rng.Int63n(15)) * 100000
		live, w := rng.Intn(5), rng.Intn(4)
		zero := mc.Date{Year: 2026, Month: 12, Day: 1}.AddDays(rng.Intn(120))
		req.Override = Constraints{&buf, &live, &zero, &w}
		res := compareT(t, req)
		checkConstraints(t, fmt.Sprintf("iteración %d", i), res)
		if (res.Status == StatusOK) != (len(res.Valid) > 0) || (res.Optima == nil) == (res.Status == StatusOK) {
			t.Errorf("iteración %d: status %s con %d válidos", i, res.Status, len(res.Valid))
		}
		total += res.EvaluatedCount
	}
	t.Logf("S3-12 propiedad: 300 comparaciones, %d candidatos, 0 válidos con violación", total)
}

// TestRestriccionesVivas: el conteo propio de compras vivas coincide con LivePurchases del núcleo.
func TestRestriccionesVivas(t *testing.T) {
	in, args, _ := loadFixture(t, "caso-1")
	req, _ := requestOf(t, in, args)
	res := compareT(t, req)
	for _, s := range append(res.Valid, rejectedScenarios(res)...) {
		var cs mc.CreditState
		for _, c := range req.Credits {
			if c.Credit.ID == s.CreditID {
				cs = c
			}
		}
		pin := mc.PurchaseInput{Credit: cs.Credit, Date: s.PurchaseDate, Principal: req.Amount, TermCount: s.TermCount,
			PromotionID: s.PromotionID, StoreID: s.StoreID, Shares: req.Shares}
		q, err := mc.QuotePurchase(pin)
		if err != nil {
			t.Fatal(err)
		}
		st := withScenario(cs, pin)
		top := 0
		for d := s.PurchaseDate; d.Compare(q.Payments[len(q.Payments)-1].Due) <= 0; d = d.AddDays(1) {
			refs, err := mc.LivePurchases(st, d)
			if err != nil {
				t.Fatal(err)
			}
			top = max(top, len(refs))
		}
		if top != s.MaxLivePurchases {
			t.Errorf("%s: vivas %d, núcleo %d", keyOf(s), s.MaxLivePurchases, top)
		}
	}
}

func TestRestriccionesFaltan(t *testing.T) {
	in, args, exp := loadFixture(t, "sin-restricciones")
	out, code := EvalComparar(in, args)
	if code != "" {
		t.Fatalf("código %q", code)
	}
	var diffs []string
	subset("sin-restricciones", exp, generic(t, out), &diffs)
	for _, d := range diffs {
		t.Error(d)
	}
	res := out.(Result)
	if res.Status != StatusMissing || res.Recommendation != nil || res.Optima != nil || res.EvaluatedCount != 0 || len(res.MissingConstraints) != 4 {
		t.Errorf("sin restricciones: %+v", res)
	}
	// una sola faltante; la de la comparación completa la global (K16)
	in1, args1, _ := loadFixture(t, "caso-1")
	req, _ := requestOf(t, in1, args1)
	req.Global.DateWindowDays = nil
	if r := compareT(t, req); r.Status != StatusMissing || !reflect.DeepEqual(r.MissingConstraints, []string{"date_window_days"}) {
		t.Errorf("faltante única: %s %v", r.Status, r.MissingConstraints)
	}
	w := 1
	req.Override.DateWindowDays = &w
	if r := compareT(t, req); r.Status != StatusOK || *r.ConstraintsUsed.DateWindowDays != 1 {
		t.Errorf("override no completó la global: %s", r.Status)
	}
}

func TestOptimos(t *testing.T) {
	for _, name := range cases {
		in, args, exp := loadFixture(t, name)
		req, _ := requestOf(t, in, args)
		res := compareT(t, req)
		o := res.Optima
		keys := map[string]bool{keyOf(o.LowestInterest): true, keyOf(o.EvenLowPayment): true, keyOf(o.MaxAvailable): true, keyOf(o.MostValuePerPeso): true}
		if len(keys) != 4 {
			t.Errorf("%s: óptimos repetidos %v", name, keys)
		}
		if o.LowestInterest.PromotionID == nil {
			t.Errorf("%s: el óptimo de intereses no es la promoción", name)
		}
		var diffs []string
		subset(name+".optima", exp.(map[string]any)["optima"], generic(t, o), &diffs)
		for _, d := range diffs {
			t.Error(d)
		}
		t.Logf("S3-13 %s: intereses %s · parejo %s · disponible %s · por peso %s · recomendación %s (%d pb)", name,
			keyOf(o.LowestInterest), keyOf(o.EvenLowPayment), keyOf(o.MaxAvailable), keyOf(o.MostValuePerPeso),
			keyOf(*res.Recommendation), *res.Recommendation.ScoreBP)
	}
}

func TestObjetivoParejo(t *testing.T) {
	if got := stdDevMilli([]mc.Cents{0, 100, 200}); got != 81650 {
		t.Errorf("0/100/200 → %d, se espera 81,650", got)
	}
	if got := stdDevMilli([]mc.Cents{973_75, 973_75, 973_75, 973_75}); got != 0 {
		t.Errorf("pagos iguales → %d", got)
	}
	d := func(y int, m time.Month, day int) mc.Date { return mc.Date{Year: y, Month: m, Day: day} }
	entries := []mc.CalendarEntry{
		{Date: d(2026, 12, 15), Amount: 100, Shares: []mc.PaymentShare{{Code: "axel", Amount: 100}}},
		{Date: d(2027, 1, 15), Amount: 200, Shares: []mc.PaymentShare{{Code: "axel", Amount: 150}, {Code: "karen", Amount: 50}}, Provisional: []string{"S3-formula-tarjeta"}},
		{Date: d(2027, 3, 15), Amount: 999},
	}
	tot, byPerson, prov := series(entries, d(2026, 12, 1), d(2027, 1, 15))
	if !reflect.DeepEqual(tot, []mc.Cents{100, 0, 200}) || !reflect.DeepEqual(prov, []string{"S3-formula-tarjeta"}) {
		t.Errorf("serie con quincena en 0: %v %v", tot, prov)
	}
	if oc := overCap([]PersonCap{{Person: "karen", Cap: 40}, {Person: "axel", Cap: 200}}, byPerson); len(oc) != 1 || oc[0].Person != "karen" || oc[0].Peak != 50 {
		t.Errorf("over_cap: %+v", oc)
	}
	if stdDevMilli(tot) != 81650 {
		t.Errorf("la quincena en 0 debe contar en n")
	}
	// empate en desviación → gana el de menos personas sobre el tope, aunque su puntaje sea menor
	one, two := int64(1), int64(2)
	v := []Scenario{
		{CreditID: 1, TermCount: 2, StdDevMilliCents: 5, Interest: 10, OverCap: []PersonPeak{{"axel", 3, 2}}, StoreID: &one},
		{CreditID: 2, TermCount: 2, StdDevMilliCents: 5, Interest: 50, StoreID: &two},
	}
	exact := scores(v, Weights{1, 1, 1, 1})
	if best(v, exact, 1) != 1 || best(v, exact, 0) != 0 {
		t.Errorf("desempate del objetivo 2 por over_cap")
	}
	if cmp := exact[0].Cmp(exact[1]); cmp <= 0 {
		t.Errorf("el de over_cap debía tener más puntaje")
	}
}

func TestPuntaje(t *testing.T) {
	v := []Scenario{{Interest: 100}, {Interest: 200}, {Interest: 400}}
	n := normalize(v, objectives[0])
	want := []*big.Rat{big.NewRat(1, 1), big.NewRat(2, 3), big.NewRat(0, 1)}
	for i := range want {
		if n[i].Cmp(want[i]) != 0 {
			t.Errorf("100/200/400 → %v", n)
		}
	}
	for _, x := range normalize(v, objectives[3]) { // todos nulos → 1
		if x.Cmp(big.NewRat(1, 1)) != 0 {
			t.Errorf("nulos deben valer 1")
		}
	}
	eq := []Scenario{{MinAvailable: 7}, {MinAvailable: 7}}
	if x := normalize(eq, objectives[2]); x[0].Cmp(big.NewRat(1, 1)) != 0 || x[1].Cmp(big.NewRat(1, 1)) != 0 {
		t.Errorf("empate debe valer 1")
	}
	ws := weightsOf(Weights{1, 1, 1, 1})
	sum := new(big.Rat)
	for _, w := range ws {
		sum.Add(sum, big.NewRat(w, 4))
	}
	if sum.Cmp(big.NewRat(1, 1)) != 0 {
		t.Errorf("pesos por defecto Σ = %v", sum)
	}
	if exact := scores(v, Weights{1, 1, 1, 1}); *v[0].ScoreBP != 10000 || *v[2].ScoreBP != 7500 || exact[1].Cmp(big.NewRat(11, 12)) != 0 {
		t.Errorf("puntajes %d %d %d", *v[0].ScoreBP, *v[1].ScoreBP, *v[2].ScoreBP)
	}
	// cambiar pesos cambia la recomendación; un descartado extremo no mueve la escala
	in, args, _ := loadFixture(t, "caso-1")
	req, _ := requestOf(t, in, args)
	base := compareT(t, req)
	req.Weights = Weights{0, 1, 0, 0}
	alt := compareT(t, req)
	if keyOf(*alt.Recommendation) == keyOf(*base.Recommendation) || keyOf(*alt.Recommendation) != keyOf(base.Optima.EvenLowPayment) {
		t.Errorf("pesos editados: %s contra %s", keyOf(*alt.Recommendation), keyOf(*base.Recommendation))
	}
	var hiValid, hiRejected mc.Cents
	for _, s := range base.Valid {
		hiValid = max(hiValid, s.Interest)
	}
	for _, s := range rejectedScenarios(base) {
		hiRejected = max(hiRejected, s.Interest)
	}
	low := big.NewRat(1, 1)
	for _, x := range normalize(base.Valid, objectives[0]) {
		if x.Cmp(low) < 0 {
			low = x
		}
	}
	if hiRejected <= hiValid || low.Sign() != 0 {
		t.Errorf("la escala debe salir solo de los válidos: máx válido %d, máx descartado %d, mínimo normalizado %v", hiValid, hiRejected, low)
	}
}

func ptr[T any](v T) *T { return &v }

func TestRestriccionesTiempo(t *testing.T) {
	in, args, _ := loadFixture(t, "tiempo")
	if out, code := EvalComparar(in, args); code != "timeout" || out != nil {
		t.Errorf("timeout_seconds 0: código %q", code)
	}
	req, _ := requestOf(t, in, args)
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	_, err := Compare(ctx, req)
	if st, code, ok := mc.ResponseFor(err); !errors.Is(err, mc.ErrTimeout) || st != 503 || code != "unavailable" || !ok {
		t.Errorf("vencido: %v → %d %q", err, st, code)
	}
	limit, _ := TimeoutFromSetting(5)
	for _, name := range append(append([]string{}, cases...), "carga") {
		// la rejilla de carga es la fixture s3-comparador-carga (6 créditos × 7 plazos × 15 fechas × 6
		// tiendas + 1 promoción por crédito; cuenta a mano 3,870)
		in, args, _ := loadFixture(t, name)
		r, _ := requestOf(t, in, args)
		ctx, cancel := context.WithTimeout(context.Background(), limit)
		start := time.Now()
		res, err := Compare(ctx, r)
		ms := time.Since(start).Milliseconds()
		cancel()
		if err != nil {
			t.Fatalf("%s con 5 s: %v", name, err)
		}
		if name == "carga" && (res.EvaluatedCount != 3870 || res.Status != StatusOK) {
			t.Errorf("carga: %s, %d evaluados (cuenta a mano 3,870)", res.Status, res.EvaluatedCount)
		}
		t.Logf("S3-TIEMPO %s %d ms (%d candidatos)", name, ms, res.EvaluatedCount)
	}
}

func TestTiempoClaveSembrada(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join("..", "..", "..", "db", "migrations", "00003_*.sql"))
	if len(files) != 1 {
		t.Fatalf("migración 00003: %v", files)
	}
	b, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`\('comparator\.timeout_seconds',\s*'(\d+)'`).FindSubmatch(b)
	if m == nil {
		t.Fatal("no se encontró el valor sembrado de comparator.timeout_seconds")
	}
	v, _ := strconv.ParseInt(string(m[1]), 10, 64)
	if d, err := TimeoutFromSetting(v); err != nil || d != 5*time.Second {
		t.Errorf("valor sembrado %d → %v %v", v, d, err)
	}
	if _, err := TimeoutFromSetting(-1); !errors.Is(err, mc.ErrInvalid) {
		t.Errorf("-1 debe ser ErrInvalid")
	}
	if d, _ := TimeoutFromSetting(0); d != 0 {
		t.Errorf("0 vence de inmediato")
	}
}
