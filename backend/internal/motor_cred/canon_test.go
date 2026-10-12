package motor_cred_test

// Pruebas canon (S3-01, S3-02): todo valor esperado se lee de oraculo/model.json y de
// datos/semilla/data.json al correr; ninguna cifra canon está escrita en este archivo.

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	mc "finanzas-hogar/internal/motor_cred"
	ct "finanzas-hogar/internal/motor_cred/canontest"
)

type canon struct {
	m, d     any
	in       mc.Input
	si       mc.SensitivityInput
	ipm      []mc.Cents
	ps       []mc.Purchase
	from, to mc.Date
}

func loadCanon(t *testing.T) canon {
	t.Helper()
	c := canon{m: ct.Load(t, "oraculo/model.json"), d: ct.Load(t, "datos/semilla/data.json")}
	c.in = ct.Seed(t, c.m, c.d)
	c.si, c.ipm = ct.Sensitivity(t, c.m)
	var err error
	if c.ps, err = c.in.Credits[0].Quote(); err != nil {
		t.Fatal(err)
	}
	c.from, c.to = ct.D(t, c.m, "parametros.horizonte[0]"), ct.D(t, c.m, "parametros.horizonte[1]")
	checkSeedFiles(t, c)
	return c
}

// checkSeedFiles: si existen las fixtures del agente de regresión, deben ser iguales (como JSON)
// al Input de seedFromModel y al insumo de sensitivityFromModel (§10.0); si faltan, no falla aquí.
func checkSeedFiles(t *testing.T, c canon) {
	t.Helper()
	root, _ := ct.Root()
	dir := filepath.Join(root, "oraculo", "fixtures", "s3")
	if b, err := os.ReadFile(filepath.Join(dir, "_datos", "kueski-semilla.json")); err == nil {
		f, err := ct.Decode(b)
		if err != nil || !reflect.DeepEqual(f, ct.JSON(t, c.in)) {
			t.Errorf("hallazgo: _datos/kueski-semilla.json difiere del Input de seedFromModel (%v)", err)
		}
	}
	want := map[string]any{
		"s3-sensibilidad":      ct.JSON(t, c.si),
		"s3-interes-por-monto": ct.JSON(t, map[string]any{"amounts_cents": c.ipm, "rate": c.si.Rate, "term_count": c.si.TermCount}),
	}
	fn := map[string]string{"s3-sensibilidad": "vista_sensibilidad", "s3-interes-por-monto": "vista_interes_por_monto"}
	for caso, w := range want {
		b, err := os.ReadFile(filepath.Join(dir, caso+".json"))
		if err != nil {
			continue
		}
		f, err := ct.Decode(b)
		if err != nil {
			t.Errorf("%s: %v", caso, err)
			continue
		}
		qs, _ := ct.Get(f, "entradas.consultas")
		list, _ := qs.([]any)
		for _, q := range list {
			if qm, _ := q.(map[string]any); qm["funcion"] == fn[caso] && !reflect.DeepEqual(qm["argumentos"], w) {
				t.Errorf("hallazgo: argumentos de %s difieren de sensitivityFromModel", caso)
			}
		}
	}
}

func eqC(t *testing.T, label string, got, want mc.Cents) {
	t.Helper()
	if got != want {
		t.Errorf("%s: motor %d, canon %d", label, got, want)
	}
}

func share(p mc.ScheduledPayment, code string) mc.Cents {
	for _, s := range p.Shares {
		if s.Code == code {
			return s.Amount
		}
	}
	return -1
}

func person(p mc.Purchase, code string) mc.PersonShare {
	for _, s := range p.Shares {
		if s.Code == code {
			return s
		}
	}
	return mc.PersonShare{}
}

func num(t *testing.T, v any, path string) mc.Cents {
	t.Helper()
	x, err := ct.Get(v, path)
	if err != nil {
		t.Fatal(err)
	}
	n, err := strconv.ParseInt(fmt.Sprint(x), 10, 64)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return mc.Cents(n)
}

func ratioFloat(num, den int64, pct bool) float64 {
	r := big.NewRat(num, den)
	if pct {
		r.Mul(r, big.NewRat(100, 1))
	}
	f, _ := r.Float64() // solo para comparar con la hoja flotante de model.json (§1.1)
	return f
}

func modelFloat(t *testing.T, m any, path string) float64 {
	t.Helper()
	x, _ := ct.Get(m, path)
	f, err := strconv.ParseFloat(fmt.Sprint(x), 64)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func h2(c mc.Cents, n int64) mc.Cents { return mc.RoundHalfUp(big.NewInt(int64(c)), big.NewInt(n)) }

func TestCanon85(t *testing.T) {
	c := loadCanon(t)
	m, st := c.m, c.in.Credits[0]
	rows, err := mc.Sensitivity(c.si)
	if err != nil {
		t.Fatal(err)
	}
	var row mc.SensitivityRow
	for _, r := range rows {
		if r.Amount == ct.C(t, m, "despensas[1].capital") {
			row = r
		}
	}
	tl, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, Cash: c.in.Cash, Today: c.in.Today, From: c.from, To: c.to})
	if err != nil {
		t.Fatal(err)
	}
	par, err := mc.ViewParametros(c.in, 1)
	if err != nil {
		t.Fatal(err)
	}
	pj := ct.JSON(t, par)
	// Efectivo y devengado por ventana (§11.1).
	nw := ct.Len(t, m, "periodos")
	eff, dev := make([]mc.Cents, nw), make([]mc.Cents, nw)
	wins := make([][2]mc.Date, nw)
	for k := range wins {
		wins[k] = [2]mc.Date{ct.D(t, m, fmt.Sprintf("periodos[%d].inicio", k)), ct.D(t, m, fmt.Sprintf("periodos[%d].fin", k))}
		cal, err := mc.UnifiedCalendar(c.in.Credits, wins[k][0], wins[k][1])
		if err != nil {
			t.Fatal(err)
		}
		w := mc.CashFlow(cal, c.in.Cash, wins[k][0], wins[k][1])
		eff[k], dev[k] = w.Total, w.Cash
	}
	for j, p := range c.ps {
		halves := mc.SplitHalfUp(p.Total, 2)
		for k := range wins {
			d := st.Purchases[j].Input.Date
			if d.Compare(wins[k][0]) >= 0 && d.Compare(wins[k][1]) <= 0 {
				dev[k] += halves[0]
				if k+1 < nw {
					dev[k+1] += halves[1]
				}
			}
		}
	}
	var sumEff, sumDev mc.Cents
	for k := range wins {
		sumEff, sumDev = sumEff+eff[k], sumDev+dev[k]
	}
	hcal, err := mc.UnifiedCalendar(c.in.Credits, c.from, c.to)
	if err != nil {
		t.Fatal(err)
	}
	hw := mc.CashFlow(hcal, c.in.Cash, c.from, c.to)
	debt, err := mc.LiveDebt(c.in.Credits, c.to)
	if err != nil {
		t.Fatal(err)
	}
	debtOf := func(ref string) mc.Cents {
		for _, l := range debt.ByPurchase {
			if l.PurchaseRef == ref {
				return l.Pending
			}
		}
		return 0
	}

	cases := map[int]func(t *testing.T){}
	purchaseRow := func(j int) func(t *testing.T) {
		return func(t *testing.T) {
			p := fmt.Sprintf("despensas[%d]", j)
			eqC(t, "capital", c.ps[j].Principal, ct.C(t, m, p+".capital"))
			eqC(t, "intereses", c.ps[j].Interest, ct.C(t, m, p+".intereses"))
			eqC(t, "total", c.ps[j].Total, ct.C(t, m, p+".total"))
		}
	}
	payRow := func(j, k int) func(t *testing.T) {
		return func(t *testing.T) {
			p, pay := fmt.Sprintf("despensas[%d].pagos[%d]", j, k), c.ps[j].Payments[k]
			if pay.Due.String() != ct.S(t, m, p+".fecha") {
				t.Errorf("fecha %s, canon %s", pay.Due, ct.S(t, m, p+".fecha"))
			}
			eqC(t, "monto", pay.Amount, ct.C(t, m, p+".monto"))
			eqC(t, "axel", share(pay, "axel"), ct.C(t, m, p+".axel"))
			eqC(t, "karen", share(pay, "karen"), ct.C(t, m, p+".karen"))
			eqC(t, "capital_repuesto", pay.Principal, ct.C(t, m, p+".capital_repuesto"))
		}
	}
	cases[1], cases[2], cases[3] = purchaseRow(0), purchaseRow(0), purchaseRow(0)
	cases[17], cases[22], cases[27] = purchaseRow(1), purchaseRow(2), purchaseRow(3)
	for j, first := range []int{13, 18, 23, 28} {
		for k := 0; k < 4; k++ {
			cases[first+k] = payRow(j, k)
		}
	}
	cases[4] = func(t *testing.T) {
		if got, want := ratioFloat(c.ps[0].Rate.Num, c.ps[0].Rate.Den, true), modelFloat(t, m, "parametros.tasa_pct"); got != want {
			t.Errorf("tasa_pct %v, canon %v", got, want)
		}
	}
	dosQ := func(t *testing.T) any {
		pin := st.Purchases[0].Input
		pin.TermCount, pin.QuotedInterest, pin.RateDecimals = 2, nil, nil
		v, err := mc.ViewDosQuincenas(c.in, 1, pin)
		if err != nil {
			t.Fatal(err)
		}
		return ct.JSON(t, v)
	}
	cases[5] = func(t *testing.T) {
		v := dosQ(t)
		dec := num(t, v, "tasa_pct.decimales")
		scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(dec)), nil)
		n := new(big.Int).Mul(big.NewInt(int64(num(t, v, "tasa_pct.num"))*100), scale)
		got := mc.RoundHalfUp(n, big.NewInt(int64(num(t, v, "tasa_pct.den"))))
		x, _ := ct.Get(c.d, "financiamiento.dos_quincenas.tasa_pct")
		want, _ := new(big.Rat).SetString(fmt.Sprint(x))
		want.Mul(want, new(big.Rat).SetInt(scale))
		if !want.IsInt() || int64(got) != want.Num().Int64() {
			t.Errorf("tasa a 2 q: %d (× 10^%d), data.json %v", got, dec, x)
		}
	}
	cases[6] = func(t *testing.T) {
		v := dosQ(t)
		for _, k := range []string{"total", "intereses", "axel.intereses", "axel.total", "axel.por_quincena", "karen.intereses", "karen.total", "karen.por_quincena"} {
			eqC(t, k, num(t, v, k), ct.C(t, c.d, "financiamiento.dos_quincenas."+k))
		}
	}
	personRow := func(code, field string) func(t *testing.T) {
		return func(t *testing.T) {
			s := person(c.ps[0], code)
			got := map[string]mc.Cents{"total": s.Total, "pago_por_quincena": s.FirstPayment, "intereses": s.Interest}[field]
			eqC(t, code+"."+field, got, ct.C(t, m, "despensas[0]."+code+"."+field))
		}
	}
	for i, f := range []string{"total", "pago_por_quincena", "intereses"} {
		cases[7+i], cases[10+i] = personRow("axel", f), personRow("karen", f)
	}
	for i, k := range []string{"total", "pago", "capital", "intereses", "conjunto"} {
		k := k
		cases[32+i] = func(t *testing.T) { eqC(t, k, num(t, pj, "karen_fija."+k), ct.C(t, m, "parametros.karen_fija."+k)) }
	}
	ve := []struct {
		got  mc.Cents
		path string
	}{{row.MonthlyTotal, "mes.total"}, {row.HolderMonthly, "mes.axel"}, {row.ExternalMonthly, "mes.karen"},
		{h2(row.MonthlyTotal, 2), "quincena.total"}, {h2(row.HolderMonthly, 2), "quincena.axel"}, {h2(row.ExternalMonthly, 2), "quincena.karen"},
		{row.MonthlyPurchase, "componentes.mitad_despensa"}, {h2(num(t, pj, "papel"), 2), "componentes.papel_contado_mes"}}
	for i, v := range ve {
		v := v
		cases[37+i] = func(t *testing.T) { eqC(t, v.path, v.got, ct.C(t, m, "valor_estable."+v.path)) }
	}
	for k := 0; k < nw; k++ {
		k := k
		cases[45+k] = func(t *testing.T) {
			p := fmt.Sprintf("periodos[%d]", k)
			eqC(t, "efectivo", eff[k], ct.C(t, m, p+".efectivo.total"))
			eqC(t, "devengado", dev[k], ct.C(t, m, p+".devengado.total"))
			var acc mc.Cents
			for i := 0; i <= k; i++ {
				acc += dev[i] - eff[i]
			}
			eqC(t, "diferencia_acumulada", acc, ct.C(t, m, p+".diferencia_acumulada"))
		}
	}
	hz := func(n int, got mc.Cents, path string) {
		cases[n] = func(t *testing.T) { eqC(t, path, got, ct.C(t, m, "horizonte."+path)) }
	}
	hz(51, sumDev, "total_devengado")
	cases[52] = func(t *testing.T) {
		eqC(t, "total_efectivo", hw.Total, ct.C(t, m, "horizonte.total_efectivo"))
		eqC(t, "suma de ventanas", sumEff, ct.C(t, m, "horizonte.total_efectivo"))
	}
	hz(53, debt.Total, "deuda_viva_10_abr")
	hz(54, debtOf("d3"), "deuda_viva_detalle.despensa_3_pendiente")
	hz(55, debtOf("d4"), "deuda_viva_detalle.despensa_4")
	hz(56, h2(hw.Total, 6), "promedio_mensual_efectivo")
	hz(57, h2(sumDev, 6), "promedio_mensual_devengado")
	hz(58, sumDev-hw.Total, "diferencia_final")
	var d3 mc.Cents
	for _, e := range hcal.Entries {
		if e.PurchaseRef == "d3" {
			d3 += e.Amount
		}
	}
	hz(59, d3, "pagos_despensa_3_en_horizonte")
	hz(60, hw.Cash, "papel_contado")
	cases[61] = func(t *testing.T) { eqC(t, "minimo_A", tl.MinA, ct.C(t, m, "credito.minimo_A")) }
	cases[62] = func(t *testing.T) { eqC(t, "minimo_B", tl.MinB, ct.C(t, m, "credito.minimo_B")) }
	cases[63] = func(t *testing.T) {
		eqC(t, "limite", num(t, pj, "limite_credito"), ct.C(t, m, "parametros.limite_credito"))
	}
	cases[64] = func(t *testing.T) {
		eqC(t, "A", tl.Events[0].AvailableA, ct.C(t, m, "credito.eventos[0].disponible_A"))
		eqC(t, "B", tl.Events[0].AvailableB, ct.C(t, m, "credito.eventos[0].disponible_B"))
	}
	series := func(got []mc.Cents, ab string) func(t *testing.T) {
		return func(t *testing.T) {
			for _, key := range []string{"esperado_", "calculado_"} {
				p := "credito.comparacion_series." + key + ab
				if len(got) != ct.Len(t, m, p) {
					t.Fatalf("%s: %d valores, canon %d", p, len(got), ct.Len(t, m, p))
				}
				for i := range got {
					eqC(t, fmt.Sprintf("%s[%d]", p, i), got[i], ct.C(t, m, fmt.Sprintf("%s[%d]", p, i)))
				}
			}
		}
	}
	cases[65], cases[66] = series(tl.SeriesA, "A"), series(tl.SeriesB, "B")
	cases[67] = func(t *testing.T) {
		v, err := mc.ViewCredito(c.in, 1, c.in.Today, c.from, c.to)
		if err != nil {
			t.Fatal(err)
		}
		g := ct.JSON(t, v)
		if n := ct.Len(t, m, "credito.eventos"); len(tl.Events) != n {
			t.Fatalf("%d eventos, canon %d", len(tl.Events), n)
		}
		for i := range tl.Events {
			p := fmt.Sprintf("eventos[%d]", i)
			if ct.S(t, g, p+".fecha") != ct.S(t, m, "credito."+p+".fecha") || ct.S(t, g, p+".tipo") != ct.S(t, m, "credito."+p+".tipo") {
				t.Errorf("%s: fecha o tipo", p)
			}
			for _, k := range []string{"monto", "disponible_A", "disponible_B"} {
				eqC(t, p+"."+k, num(t, g, p+"."+k), ct.C(t, m, "credito."+p+"."+k))
			}
		}
	}
	type pkg struct {
		date        mc.Date
		who, how    string
		amount      mc.Cents
		order, rank int
	}
	var pkgs []pkg
	for j, p := range c.ps {
		for _, s := range p.Shares {
			if s.Composition != nil && s.Composition.PaperIn > 0 {
				pkgs = append(pkgs, pkg{st.Purchases[j].Input.Date, s.Code, fmt.Sprintf("dentro de despensa %d", j+1), s.Composition.PaperIn, len(pkgs), 0})
			}
		}
	}
	for _, ce := range c.in.Cash {
		pkgs = append(pkgs, pkg{ce.Date, ce.Person, "de contado", ce.Amount, len(pkgs), 1})
	}
	sort.SliceStable(pkgs, func(i, j int) bool { return pkgs[i].date.Compare(pkgs[j].date) < 0 })
	for k := 0; k < 6; k++ {
		k := k
		cases[68+k] = func(t *testing.T) {
			p := fmt.Sprintf("papel[%d]", k)
			if k >= len(pkgs) {
				t.Fatalf("faltan paquetes: %d", len(pkgs))
			}
			g := pkgs[k]
			if int(num(t, m, p+".paquete")) != k+1 || g.date.String() != ct.S(t, m, p+".compra") ||
				!strings.EqualFold(g.who, ct.S(t, m, p+".quien")) || g.how != ct.S(t, m, p+".forma") {
				t.Errorf("%s: %+v", p, g)
			}
			eqC(t, p+".monto", g.amount, ct.C(t, m, p+".monto"))
		}
	}
	for k := 0; k < 5; k++ {
		k := k
		cases[74+k] = func(t *testing.T) {
			p, r := fmt.Sprintf("sensibilidad[%d]", k), rows[k]
			eqC(t, "intereses", r.Interest, ct.C(t, m, p+".intereses"))
			eqC(t, "total", r.Total, ct.C(t, m, p+".total"))
			eqC(t, "pago_quincenal", r.ReferencePayment, ct.C(t, m, p+".pago_quincenal"))
			eqC(t, "axel_total", r.HolderTotal, ct.C(t, m, p+".axel_total"))
		}
	}
	for k := 0; k < 7; k++ {
		k := k
		cases[79+k] = func(t *testing.T) {
			p, q := fmt.Sprintf("interes_por_monto[%d]", k), mc.QuoteByAmount(c.ipm[k], c.si.Rate, c.si.TermCount)
			eqC(t, "intereses", q.Interest, ct.C(t, m, p+".intereses"))
			eqC(t, "total", q.Total, ct.C(t, m, p+".total"))
			eqC(t, "pago_quincenal", q.ReferencePayment, ct.C(t, m, p+".pago_quincenal"))
		}
	}
	for i := 1; i <= 85; i++ {
		f, ok := cases[i]
		if !ok {
			t.Errorf("R-%02d sin aserción", i)
			continue
		}
		t.Run(fmt.Sprintf("R-%02d", i), f)
	}
}

var countKeys = map[string]bool{"n": true, "n_pagos": true, "despensa": true, "despensa_que_cubre": true, "paquete": true}

// walkLeaves compara cada hoja numérica de model.json con la salida de la vista (FORMATO.md §3).
func walkLeaves(t *testing.T, path, key string, want, got any, count *int) {
	if !hasNumber(want) {
		return // textos, booleanos y listas de textos (parametros.horizonte) no se comparan
	}
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			t.Errorf("%s: la vista no trae objeto", path)
			return
		}
		for k, wv := range w {
			walkLeaves(t, path+"."+k, k, wv, g[k], count)
		}
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			t.Errorf("%s: largo distinto", path)
			return
		}
		for i := range w {
			walkLeaves(t, fmt.Sprintf("%s[%d]", path, i), key, w[i], g[i], count)
		}
	case json.Number:
		*count++
		switch {
		case key == "tasa" || key == "tasa_pct":
			g, _ := got.(map[string]any)
			n, _ := strconv.ParseInt(fmt.Sprint(g["num"]), 10, 64)
			d, _ := strconv.ParseInt(fmt.Sprint(g["den"]), 10, 64)
			wf, _ := strconv.ParseFloat(w.String(), 64)
			if d <= 0 || ratioFloat(n, d, key == "tasa_pct") != wf {
				t.Errorf("%s: razón %v, canon %v", path, got, w)
			}
		case countKeys[key]:
			if fmt.Sprint(got) != w.String() {
				t.Errorf("%s: %v, canon %v", path, got, w)
			}
		default:
			wc, err := ct.CentsOf(w)
			if err != nil || fmt.Sprint(got) != strconv.FormatInt(int64(wc), 10) {
				t.Errorf("%s: motor %v, canon %v centavos (%v)", path, got, wc, err)
			}
		}
	}
}

func TestCanonHojas437(t *testing.T) {
	c := loadCanon(t)
	var desp, disp []any
	for n := 1; n <= len(c.ps); n++ {
		a, err1 := mc.ViewDespensa(c.in, 1, n)
		b, err2 := mc.ViewDisposicion(c.in, 1, n)
		if err1 != nil || err2 != nil {
			t.Fatal(err1, err2)
		}
		desp, disp = append(desp, a), append(disp, b)
	}
	par, err1 := mc.ViewParametros(c.in, 1)
	cre, err2 := mc.ViewCredito(c.in, 1, c.in.Today, c.from, c.to)
	sen, err3 := mc.ViewSensibilidad(c.si)
	ipm, err4 := mc.ViewInteresPorMonto(c.ipm, c.si.Rate, c.si.TermCount)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		t.Fatal(err1, err2, err3, err4)
	}
	branches := []struct {
		name  string
		view  any
		count int // conteo de S3-interfaz.md §12
	}{{"parametros", par, 9}, {"despensas", desp, 144}, {"credito", cre, 143},
		{"sensibilidad", sen, 45}, {"disposicion_por_despensa", disp, 68}, {"interes_por_monto", ipm, 28}}
	total := 0
	for _, b := range branches {
		n := 0
		want, _ := ct.Get(c.m, b.name)
		walkLeaves(t, b.name, b.name, want, ct.JSON(t, b.view), &n)
		if n != b.count {
			t.Errorf("%s: %d hojas, se esperaban %d", b.name, n, b.count)
		}
		total += n
	}
	if total != 437 {
		t.Errorf("%d hojas numéricas, se esperaban 437", total)
	}
	t.Logf("model.json (ramas de s3): %d hojas comparadas", total)
}

func hasNumber(v any) bool {
	switch w := v.(type) {
	case json.Number:
		return true
	case map[string]any:
		for _, x := range w {
			if hasNumber(x) {
				return true
			}
		}
	case []any:
		for _, x := range w {
			if hasNumber(x) {
				return true
			}
		}
	}
	return false
}
