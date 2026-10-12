package evaluar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mc "finanzas-hogar/internal/motor_cred"
	ct "finanzas-hogar/internal/motor_cred/canontest"
)

func call(t *testing.T, input []byte, fn, args string) (any, string) {
	t.Helper()
	out, code := Evaluate(input, fn, json.RawMessage(args))
	if code != "" {
		return nil, code
	}
	return ct.JSON(t, out), ""
}

func TestEvaluar(t *testing.T) {
	m, d := ct.Load(t, "oraculo/model.json"), ct.Load(t, "datos/semilla/data.json")
	seed := ct.Seed(t, m, d)
	input, _ := json.Marshal(seed)
	si, ipm := ct.Sensitivity(t, m)
	sij, _ := json.Marshal(si)
	d1, _ := json.Marshal(seed.Credits[0].Purchases[0].Input)
	from, to := ct.S(t, m, "parametros.horizonte[0]"), ct.S(t, m, "parametros.horizonte[1]")
	win := fmt.Sprintf(`{"from": %q, "to": %q}`, from, to)
	rate, _ := json.Marshal(si.Rate)
	amounts, _ := json.Marshal(ipm)
	ok := map[string]string{
		"cotizar_compra":          `{"credit_id": 1, "purchase": ` + string(d1) + `}`,
		"promociones_aplicables":  fmt.Sprintf(`{"credit_id": 1, "date": %q, "principal_cents": 100, "term_count": 4, "store_id": null, "merchant": ""}`, from),
		"linea_tiempo":            fmt.Sprintf(`{"credit_id": 1, "today": %q, "from": %q, "to": %q, "reading": ""}`, from, from, to),
		"calendario_unificado":    win,
		"efectivo_ventana":        win,
		"deuda_viva":              fmt.Sprintf(`{"as_of": %q}`, to),
		"compras_vivas":           fmt.Sprintf(`{"credit_id": 1, "at": %q}`, to),
		"cotizar_monto":           fmt.Sprintf(`{"amount_cents": %d, "rate": %s, "term_count": 4}`, ipm[0], rate),
		"sensibilidad":            string(sij),
		"amortizar":               fmt.Sprintf(`{"principal_cents": 100000, "annual_rate": {"num": 24, "den": 100, "decimales": 0}, "vat_bp": 1600, "term_count": 7, "frequency": "monthly", "first_date": %q}`, to),
		"vista_parametros":        `{"credit_id": 1}`,
		"vista_despensa":          `{"credit_id": 1, "n": 1}`,
		"vista_disposicion":       `{"credit_id": 1, "n": 2}`,
		"vista_credito":           fmt.Sprintf(`{"credit_id": 1, "today": %q, "from": %q, "to": %q}`, from, from, to),
		"vista_sensibilidad":      string(sij),
		"vista_interes_por_monto": fmt.Sprintf(`{"amounts_cents": %s, "rate": %s, "term_count": 4}`, amounts, rate),
	}
	outs := map[string]any{}
	for fn, args := range ok {
		out, code := call(t, input, fn, args)
		if code != "" {
			t.Errorf("%s: código %q", fn, code)
		}
		outs[fn] = out
	}
	chk := func(fn, path string, want mc.Cents) {
		x, err := ct.Get(outs[fn], path)
		if err != nil || fmt.Sprint(x) != fmt.Sprint(int64(want)) {
			t.Errorf("%s %s: %v, se esperaba %d (%v)", fn, path, x, want, err)
		}
	}
	chk("cotizar_compra", "interest_cents", ct.C(t, m, "despensas[0].intereses"))
	chk("efectivo_ventana", "total_cents", ct.C(t, m, "horizonte.total_efectivo"))
	chk("deuda_viva", "total_cents", ct.C(t, m, "horizonte.deuda_viva_10_abr"))
	chk("vista_credito", "minimo_B", ct.C(t, m, "credito.minimo_B"))
	chk("vista_despensa", "pagos[3].monto", ct.C(t, m, "despensas[0].pagos[3].monto"))
	chk("vista_disposicion", "karen.total", ct.C(t, m, "disposicion_por_despensa[1].karen.total"))
	chk("sensibilidad", "rows[0].interest_cents", ct.C(t, m, "sensibilidad[0].intereses"))
	chk("cotizar_monto", "reference_payment_cents", ct.C(t, m, "interes_por_monto[0].pago_quincenal"))
	if r, _ := ct.Get(outs["linea_tiempo"], "reading"); r != "B" {
		t.Errorf("lectura omitida: %v", r)
	}
	// vista_dos_quincenas contra data.json leído al correr (R-05, R-06).
	pin := seed.Credits[0].Purchases[0].Input
	pin.TermCount, pin.QuotedInterest, pin.RateDecimals = 2, nil, nil
	pj, _ := json.Marshal(pin)
	v, code := call(t, input, "vista_dos_quincenas", `{"credit_id": 1, "purchase": `+string(pj)+`}`)
	if code != "" {
		t.Fatalf("vista_dos_quincenas: %s", code)
	}
	for _, k := range []string{"total", "intereses", "axel.intereses", "axel.total", "axel.por_quincena", "karen.intereses", "karen.total", "karen.por_quincena"} {
		x, _ := ct.Get(v, k)
		if want := ct.C(t, d, "financiamiento.dos_quincenas."+k); fmt.Sprint(x) != fmt.Sprint(int64(want)) {
			t.Errorf("dos quincenas %s: %v, data.json %d", k, x, want)
		}
	}
	card := `{"credits": [{"credit": {"id": 2, "kind": "bank_card", "holder": "axel", "limit_cents": 2000000, "statement_day": 5, "payment_due_day": 25,
		"annual_rate": {"num": 455, "den": 1000, "decimales": 1}, "rate_base": "monthly", "vat_on_interest_bp": 1600, "min_payment_pct_bp": 500, "min_payment_fixed_cents": 20000}}]}`
	st := `"statement": {"statement_date": "2026-11-05", "due_date": "2026-11-25", "balance_cents": 100000, "daily_balances_cents": [], "source": "bank"}`
	if v, code := call(t, []byte(card), "interes_tarjeta", `{"credit_id": 2, `+st+`, "paid_cents": 0, "paid_date": "2026-11-26"}`); code != "" || fmt.Sprint(v) == "" {
		t.Errorf("interes_tarjeta: %s", code)
	}
	if _, code := call(t, []byte(card), "minimo_tarjeta", `{"credit_id": 2, `+st+`}`); code != "" {
		t.Errorf("minimo_tarjeta: %s", code)
	}
	bad := []struct{ in, fn, args, code string }{
		{string(input), "no_existe", `{}`, "invalid"},
		{string(input), "vista_parametros", `{"credit_id": 1, "otro": 1}`, "unknown_field"},
		{string(input), "vista_parametros", `{"credit_id": 9}`, "invalid"},
		{`{"otro": 1}`, "vista_parametros", `{}`, "unknown_field"},
		{string(input), "comparar", `{}`, "invalid"},
		{string(input), "comparar", `{"timeout_seconds": 5, "request": {}, "otro": 1}`, "unknown_field"},
		{string(input), "comparar", `{"timeout_seconds": -1, "request": {}}`, "invalid"},
		{string(input), "cotizar_monto", `{"amounts_cents": [1], "term_count": 4}`, "invalid"},
		{card, "minimo_tarjeta", `{"credit_id": 2, "statement": {"balance_cents": -1}}`, "invalid"},
		{card, "interes_tarjeta", `{"credit_id": 2, "statement": {"balance_cents": -1}}`, "invalid"},
		{string(input), "compras_vivas", `{"credit_id": 2}`, "invalid"},
		{string(input), "linea_tiempo", `{"credit_id": 1}`, "invalid"},
		{string(input), "sensibilidad", `{"rate": {"num": 1, "den": 0}}`, "invalid"},
		{string(input), "vista_sensibilidad", `{"term_count": 0}`, "invalid"},
		{string(input), "vista_interes_por_monto", `{"term_count": 0}`, "invalid"},
		{string(input), "vista_despensa", `{"credit_id": 1, "n": 9}`, "invalid"},
		{string(input), "vista_credito", `{"credit_id": 1}`, "invalid"},
		{string(input), "vista_dos_quincenas", `{"credit_id": 1, "purchase": {}}`, "invalid"},
		{string(input), "calendario_unificado", `{"from": "x"}`, "invalid"},
		{string(input), "efectivo_ventana", `{"from": "2027-01-01", "to": "2026-01-01"}`, "invalid"},
		{string(input), "deuda_viva", `{"as_of": 3}`, "invalid"},
		{string(input), "amortizar", `null`, "invalid"},
		{``, "vista_parametros", `{"credit_id": 1}`, "invalid"},
	}
	for _, b := range bad {
		if _, code := call(t, []byte(b.in), b.fn, b.args); code != b.code {
			t.Errorf("%s %s: código %q, se esperaba %q", b.fn, b.args, code, b.code)
		}
	}
	if _, code := call(t, nil, "cotizar_monto", `{"amount_cents": 100, "rate": {"num": 1, "den": 4, "decimales": 0}, "term_count": 4}`); code != "" {
		t.Errorf("sin datos: %s", code)
	}
}

// Fixtures de §10.0: archivo y agente autor. Si falta alguno, la prueba falla con su nombre y su autor.
var required = func() [][2]string {
	reg := []string{"_datos/kueski-semilla.json", "s3-kueski-despensa-1", "s3-kueski-despensa-2", "s3-kueski-despensa-3", "s3-kueski-despensa-4",
		"s3-feb-2027", "s3-series-a-b", "s3-series-desde-otra-fecha", "s3-reposicion-separada", "s3-deuda-viva-10-abr-2027", "s3-sensibilidad",
		"s3-interes-por-monto", "s3-tasa-plazo", "s3-promocion-0", "s3-tarjeta-sin-intereses", "s3-tarjeta-minimo", "s3-prestamo-frances",
		"s3-tienda-abonos", "s3-tienda-abonos-interes", "s3-unificado-sin-corte", "s3-unificado", "s3-unificado-duplicado"}
	opt := []string{"_datos/comparador-caso-1.json", "_datos/comparador-caso-2.json", "_datos/comparador-caso-3.json", "s3-comparador-caso-1",
		"s3-comparador-caso-2", "s3-comparador-caso-3", "s3-comparador-sin-restricciones", "s3-comparador-tiempo", "s3-comparador-carga"}
	var out [][2]string
	for _, f := range reg {
		out = append(out, [2]string{f, "regresión S3"})
	}
	for _, f := range opt {
		out = append(out, [2]string{f, "optimizador S3"})
	}
	return out
}()

func mergePatch(target, patch any) any {
	p, ok := patch.(map[string]any)
	if !ok {
		return patch
	}
	t, ok := target.(map[string]any)
	if !ok {
		t = map[string]any{}
	}
	for k, v := range p {
		if v == nil {
			delete(t, k)
		} else {
			t[k] = mergePatch(t[k], v)
		}
	}
	return t
}

// subset compara esperado (subconjunto de claves; listas completas) con la salida del motor.
func subset(path string, want, got any) []string {
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return []string{path + ": __falta__"}
		}
		var out []string
		for k, wv := range w {
			gv, ok := g[k]
			if !ok {
				out = append(out, path+"."+k+": __falta__")
				continue
			}
			out = append(out, subset(path+"."+k, wv, gv)...)
		}
		return out
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return []string{path + ".#: largo distinto"}
		}
		var out []string
		for i := range w {
			out = append(out, subset(fmt.Sprintf("%s[%d]", path, i), w[i], g[i])...)
		}
		return out
	}
	if fmt.Sprint(want) != fmt.Sprint(got) {
		return []string{fmt.Sprintf("%s: motor %v, esperado %v", path, got, want)}
	}
	return nil
}

func runFixture(t *testing.T, dir, path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	f, err := ct.Decode(b)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	datos, _ := ct.Get(f, "entradas.datos")
	if dm, ok := datos.(map[string]any); ok {
		if ref, ok := dm["$ref"].(string); ok {
			rb, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(ref)))
			if err != nil {
				t.Fatalf("%s: %v", ref, err)
			}
			base, _ := ct.Decode(rb)
			if p, ok := dm["parche"]; ok {
				base = mergePatch(base, p)
			}
			datos = base
		}
	}
	input, _ := json.Marshal(datos)
	qs, _ := ct.Get(f, "entradas.consultas")
	esperado, _ := ct.Get(f, "esperado")
	em, _ := esperado.(map[string]any)
	list, _ := qs.([]any)
	for _, q := range list {
		qm, _ := q.(map[string]any)
		id, fn := fmt.Sprint(qm["id"]), fmt.Sprint(qm["funcion"])
		args, _ := json.Marshal(qm["argumentos"])
		want, ok := em[id]
		if !ok {
			continue
		}
		out, code := Evaluate(input, fn, args)
		var got any = map[string]any{"error": code}
		if code == "" {
			got = ct.JSON(t, out)
		}
		for _, diff := range subset(filepath.Base(path)+":"+id, want, got) {
			t.Error(diff)
		}
	}
}

func TestOraculoFixtures(t *testing.T) {
	root, err := ct.Root()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "oraculo", "fixtures", "s3")
	for _, r := range required {
		name := r[0]
		if !strings.HasSuffix(name, ".json") {
			name += ".json"
		}
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Errorf("falta oraculo/fixtures/s3/%s (autor: agente de %s)", name, r[1])
		}
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	for _, p := range files {
		runFixture(t, dir, p)
	}
	t.Logf("%d fixtures corridas", len(files))
}
