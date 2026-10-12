package motor_cred_test

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	mc "finanzas-hogar/internal/motor_cred"
	ct "finanzas-hogar/internal/motor_cred/canontest"
)

func timeline(t *testing.T, c canon, st mc.CreditState, today, from, to mc.Date) mc.Timeline {
	t.Helper()
	tl, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, Cash: c.in.Cash, Today: today, From: from, To: to})
	if err != nil {
		t.Fatal(err)
	}
	return tl
}

func TestSeriesAB(t *testing.T) {
	c := loadCanon(t)
	st := c.in.Credits[0]
	tl := timeline(t, c, st, c.in.Today, c.from, c.to)
	n := 0
	for _, ab := range []struct {
		got []mc.Cents
		key string
	}{{tl.SeriesA, "esperado_A"}, {tl.SeriesB, "esperado_B"}} {
		for i := 0; i < ct.Len(t, c.m, "credito.comparacion_series."+ab.key); i++ {
			eqC(t, ab.key, ab.got[i], ct.C(t, c.m, fmt.Sprintf("credito.comparacion_series.%s[%d]", ab.key, i)))
			n++
		}
	}
	if n != 2*len(tl.SeriesA) || len(tl.SeriesA) != len(tl.SeriesB) {
		t.Errorf("%d valores comparados, series de %d y %d", n, len(tl.SeriesA), len(tl.SeriesB))
	}
	st.Credit.Reading = ""
	om := timeline(t, c, st, c.in.Today, c.from, c.to)
	if om.Reading != "B" || !om.PendingConfirmation || fmt.Sprint(om.SeriesB) != fmt.Sprint(tl.SeriesB) {
		t.Errorf("lectura omitida: %s %v", om.Reading, om.PendingConfirmation)
	}
	st.Credit.Reading = "C"
	_, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, From: c.from, To: c.to})
	wantErr(t, "lectura C", err, mc.ErrInvalid)
	t.Logf("S3-03: %d de %d valores", n, n)
}

func TestSeriesABDesdeOtraFecha(t *testing.T) {
	c := loadCanon(t)
	st := c.in.Credits[0]
	base := timeline(t, c, st, c.in.Today, c.from, c.to)
	checked := 0
	for i := 1; i < len(base.Events); i++ {
		from := base.Events[i].Date
		last := i
		for last+1 < len(base.Events) && base.Events[last+1].Date.Compare(from) <= 0 {
			last++
		}
		tl := timeline(t, c, st, c.in.Today, from, c.to)
		if tl.Events[0].AvailableA != base.Events[last].AvailableA || tl.Events[0].AvailableB != base.Events[last].AvailableB {
			t.Errorf("desde %s: inicio %d/%d, canon %d/%d", from, tl.Events[0].AvailableA, tl.Events[0].AvailableB, base.Events[last].AvailableA, base.Events[last].AvailableB)
		}
		rest := base.Events[last+1:]
		if len(tl.Events)-1 != len(rest) {
			t.Fatalf("desde %s: %d eventos, canon %d", from, len(tl.Events)-1, len(rest))
		}
		for k, e := range rest {
			g := tl.Events[k+1]
			g.Description, e.Description = "", ""
			if g != e {
				t.Errorf("desde %s, evento %d: %+v ≠ %+v", from, k, g, e)
			}
		}
		checked++
	}
	_, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, From: c.from.AddDays(-1), To: c.to})
	wantErr(t, "from antes del ancla", err, mc.ErrInvalid)
	st.Credit.ReadingAAnchor, st.Credit.Reading = nil, "A"
	_, err = mc.AvailabilityTimeline(mc.TimelineInput{State: st, From: c.from, To: c.to})
	wantErr(t, "lectura A sin ancla", err, mc.ErrInvalid)
	st.Credit.Reading = "B"
	if tl := timeline(t, c, st, c.in.Today, c.from, c.to); !tl.AAnchorDefaulted || tl.AAnchor != c.from {
		t.Errorf("ancla por omisión")
	}
	_, err = mc.AvailabilityTimeline(mc.TimelineInput{State: st, From: c.to, To: c.from})
	wantErr(t, "ventana invertida", err, mc.ErrInvalid)
	t.Logf("%d fechas de inicio probadas", checked)
}

func TestCalendarioKueski(t *testing.T) {
	c := loadCanon(t)
	n := 0
	for j, p := range c.ps {
		for k, s := range p.Payments {
			if want := ct.S(t, c.m, fmt.Sprintf("despensas[%d].pagos[%d].fecha", j, k)); s.Due.String() != want {
				t.Errorf("despensa %d pago %d: %s, canon %s", j+1, k+1, s.Due, want)
			}
			n++
		}
	}
	in := c.in.Credits[0].Purchases[2].Input
	in.Credit, in.Date = c.in.Credits[0].Credit, day(t, "2028-02-10")
	p, err := mc.QuotePurchase(in)
	if err != nil || p.Payments[0].Due != day(t, "2028-02-29") || p.Payments[3].Due != day(t, "2028-04-15") {
		t.Errorf("febrero bisiesto: %v", err)
	}
	t.Logf("S3-04: %d de %d fechas", n, n)
}

func TestReposicionCapital(t *testing.T) {
	c := loadCanon(t)
	tl := timeline(t, c, c.in.Credits[0], c.in.Today, c.from, c.to)
	got := map[string]mc.Cents{}
	for _, e := range tl.Events {
		got[e.PurchaseRef] += e.PrincipalRestored
		if e.Kind == "cash" && (e.AvailableA != 0 || e.PrincipalRestored != 0) && e.PrincipalRestored != 0 {
			t.Errorf("el contado movió la línea")
		}
	}
	for j, p := range c.ps {
		var s mc.Cents
		for _, x := range p.Payments {
			s += x.Principal
		}
		eqC(t, fmt.Sprintf("Σ capital repuesto despensa %d", j+1), s, ct.C(t, c.m, fmt.Sprintf("despensas[%d].capital", j)))
	}
	for k := range c.ps[:2] { // las dos primeras despensas se pagan completas dentro del horizonte
		ref := fmt.Sprintf("d%d", k+1)
		eqC(t, "repuesto en la línea "+ref, got[ref], c.ps[k].Principal)
	}
}

// karenAxel: crédito con solo la despensa 1 y abonos registrados.
func karenAxel(c canon, recs []mc.PaymentRecord) mc.CreditState {
	st := c.in.Credits[0]
	st.Purchases = st.Purchases[:1]
	st.Payments = recs
	return st
}

func TestReposicionParteSeparada(t *testing.T) {
	c := loadCanon(t)
	m := c.m
	kar, axe := ct.C(t, m, "despensas[0].pagos[0].karen"), ct.C(t, m, "despensas[0].pagos[0].axel")
	k1 := ct.C(t, m, "despensas[0].pagos[0].capital_repuesto")
	i1 := ct.C(t, m, "despensas[0].pagos[0].monto") - k1
	st := karenAxel(c, []mc.PaymentRecord{{PurchaseRef: "d1", Seq: 1, Date: day(t, "2026-10-31"), Payer: "axel", Amount: axe},
		{PurchaseRef: "d1", Seq: 1, Date: day(t, "2026-10-29"), Payer: "karen", Amount: kar}})
	tl := timeline(t, c, st, c.in.Today, c.from, day(t, "2026-11-15"))
	a0, b0 := tl.Events[0].AvailableA, tl.Events[0].AvailableB
	want := []mc.TimelineEvent{
		{Date: c.from, Kind: "start", AvailableA: a0, AvailableB: b0},
		{Date: day(t, "2026-10-29"), Kind: "payment", PurchaseRef: "d1", Seq: 1, Amount: kar, AvailableA: a0, AvailableB: b0},
		{Date: day(t, "2026-10-31"), Kind: "payment", PurchaseRef: "d1", Seq: 1, Amount: axe, PrincipalRestored: k1, InterestNotRestored: i1, AvailableA: a0 + k1, AvailableB: b0 + k1},
		{Date: day(t, "2026-11-15"), Kind: "payment", PurchaseRef: "d1", Seq: 2, Amount: ct.C(t, m, "despensas[0].pagos[1].monto"),
			PrincipalRestored: ct.C(t, m, "despensas[0].pagos[1].capital_repuesto"), InterestNotRestored: ct.C(t, m, "despensas[0].pagos[1].monto") - ct.C(t, m, "despensas[0].pagos[1].capital_repuesto"),
			AvailableA: a0 + k1 + ct.C(t, m, "despensas[0].pagos[1].capital_repuesto"), AvailableB: b0 + k1 + ct.C(t, m, "despensas[0].pagos[1].capital_repuesto")},
	}
	if len(tl.Events) != len(want) {
		t.Fatalf("%d eventos, se esperaban %d: %+v", len(tl.Events), len(want), tl.Events)
	}
	for i := range want {
		g := tl.Events[i]
		g.Description = ""
		if g != want[i] {
			t.Errorf("evento %d: %+v, se esperaba %+v", i, g, want[i])
		}
	}
	if len(tl.SeriesA) != 3 || tl.MinB != b0 || tl.Reading != "B" || !tl.PendingConfirmation {
		t.Errorf("series %v, mínimo B %d", tl.SeriesA, tl.MinB)
	}
}

func TestReposicionAbonoParcial(t *testing.T) {
	c := loadCanon(t)
	m := c.m
	kar := ct.C(t, m, "despensas[0].pagos[1].karen")
	m2 := ct.C(t, m, "despensas[0].pagos[1].monto")
	st := karenAxel(c, []mc.PaymentRecord{{PurchaseRef: "d1", Seq: 2, Date: day(t, "2026-11-10"), Payer: "karen", Amount: kar},
		{PurchaseRef: "d1", Seq: 1, Date: day(t, "2026-10-20"), Payer: "karen", Amount: kar}})
	tl := timeline(t, c, st, c.in.Today, c.from, day(t, "2026-11-30"))
	var proj, abono2 *mc.TimelineEvent
	for i, e := range tl.Events {
		if e.Seq == 2 && e.Date == day(t, "2026-11-15") {
			proj = &tl.Events[i]
		}
		if e.Seq == 2 && e.Date == day(t, "2026-11-10") {
			abono2 = &tl.Events[i]
		}
	}
	if proj == nil || abono2 == nil || proj.Amount != m2-kar || proj.PrincipalRestored == 0 || abono2.PrincipalRestored != 0 {
		t.Errorf("resto proyectado: %+v %+v", proj, abono2)
	}
	// Con hoy después del vencimiento del pago 1 abonado en parte: no se proyecta ni repone.
	late := timeline(t, c, st, day(t, "2026-11-01"), c.from, day(t, "2026-11-30"))
	for _, e := range late.Events {
		if e.Seq == 1 && (e.PrincipalRestored != 0 || e.Date == day(t, "2026-10-30")) {
			t.Errorf("pago vencido sin completar repuso: %+v", e)
		}
	}
	st.Payments = append(st.Payments, mc.PaymentRecord{PurchaseRef: "d1", Seq: 9, Date: day(t, "2026-11-10"), Amount: 1})
	_, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, From: c.from, To: c.to})
	wantErr(t, "abono a un pago inexistente", err, mc.ErrInvalid)
}

func TestTasaPorCompra(t *testing.T) {
	c := loadCanon(t)
	d1 := c.ps[0]
	if d1.Rate != mc.QuotedRate(ct.C(t, c.m, "despensas[0].intereses"), ct.C(t, c.m, "despensas[0].capital"), 4) || d1.Source != "quoted" {
		t.Errorf("tasa cotizada: %+v", d1.Rate)
	}
	if want := strconv.FormatFloat(modelFloat(t, c.m, "parametros.tasa_pct"), 'f', 4, 64) + " %"; d1.RateDisplay != want {
		t.Errorf("pantalla %q, se esperaba %q", d1.RateDisplay, want)
	}
	in := c.in.Credits[0].Purchases[0].Input
	in.Credit, in.RateDecimals = c.in.Credits[0].Credit, nil
	p, err := mc.QuotePurchase(in)
	if err != nil || p.Interest != d1.Interest || p.Rate.Decimals != -1 || p.RateDisplay != "" {
		t.Errorf("sin rate_decimals: %+v %v", p.Rate, err)
	}
	if c.ps[1].Source != "term_table" || c.ps[1].Interest != ct.C(t, c.m, "despensas[1].intereses") {
		t.Errorf("tabla por plazo: %d", c.ps[1].Interest)
	}
}

func TestTasaSinCaptura(t *testing.T) {
	c := loadCanon(t)
	in := c.in.Credits[0].Purchases[1].Input
	in.Credit = c.in.Credits[0].Credit
	for _, n := range []int{6, 8, 10, 12} {
		in.TermCount = n
		_, err := mc.QuotePurchase(in)
		wantErr(t, fmt.Sprintf("%d quincenas", n), err, mc.ErrNoRate)
	}
}

func TestTasaPorMonto(t *testing.T) {
	c := loadCanon(t)
	v, err := mc.ViewInteresPorMonto(c.ipm, c.si.Rate, c.si.TermCount)
	if err != nil {
		t.Fatal(err)
	}
	g := ct.JSON(t, v)
	for k := range c.ipm {
		for _, f := range []string{"monto", "intereses", "total", "pago_quincenal"} {
			p := fmt.Sprintf("[%d].%s", k, f)
			eqC(t, p, num(t, g, p), ct.C(t, c.m, "interes_por_monto"+p))
		}
	}
	if _, err := mc.ViewInteresPorMonto(c.ipm, mc.Rate{}, 4); err == nil {
		t.Errorf("tasa vacía aceptada")
	}
}

func TestTasaSensibilidad(t *testing.T) {
	c := loadCanon(t)
	v, err := mc.ViewSensibilidad(c.si)
	if err != nil {
		t.Fatal(err)
	}
	g := ct.JSON(t, v)
	for k := range c.si.Amounts {
		for _, f := range []string{"N", "intereses", "total", "pago_quincenal", "axel_total", "axel_capital", "costo_mensual_total", "axel_mensual", "karen_mensual"} {
			p := fmt.Sprintf("[%d].%s", k, f)
			eqC(t, p, num(t, g, p), ct.C(t, c.m, "sensibilidad"+p))
		}
		x, _ := ct.Get(g, fmt.Sprintf("[%d].axel_menor_que_actual", k))
		w, _ := ct.Get(c.m, fmt.Sprintf("sensibilidad[%d].axel_menor_que_actual", k))
		if x != w {
			t.Errorf("axel_menor_que_actual[%d]: %v, canon %v", k, x, w)
		}
	}
	bad := c.si
	bad.AccrualMonths = 0
	if _, err := mc.ViewSensibilidad(bad); err == nil {
		t.Errorf("meses 0 aceptado")
	}
}

func TestEfectivoVentana(t *testing.T) {
	c := loadCanon(t)
	for k := 0; k < ct.Len(t, c.m, "periodos"); k++ {
		p := fmt.Sprintf("periodos[%d]", k)
		from, to := ct.D(t, c.m, p+".inicio"), ct.D(t, c.m, p+".fin")
		cal, err := mc.UnifiedCalendar(c.in.Credits, from, to)
		if err != nil {
			t.Fatal(err)
		}
		w := mc.CashFlow(cal, c.in.Cash, from, to)
		eqC(t, p, w.Total, ct.C(t, c.m, p+".efectivo.total"))
		var s mc.Cents
		for _, x := range w.ByPerson {
			s += x.Amount
		}
		eqC(t, p+" por persona", s, w.Total)
		if w.PaymentCount != int(num(t, c.m, p+".efectivo.n_pagos")) {
			t.Errorf("%s: %d pagos", p, w.PaymentCount)
		}
	}
}

func TestDeudaViva(t *testing.T) {
	c := loadCanon(t)
	d, err := mc.LiveDebt(c.in.Credits, c.to)
	if err != nil {
		t.Fatal(err)
	}
	eqC(t, "deuda viva", d.Total, ct.C(t, c.m, "horizonte.deuda_viva_10_abr"))
	if len(d.ByCredit) != 1 || d.ByCredit[0].Pending != d.Total || len(d.ByPurchase) != 2 {
		t.Errorf("detalle: %+v", d)
	}
	st := c.in.Credits[0]
	st.Purchases = append([]mc.PurchaseState(nil), st.Purchases...)
	st.Purchases[3].Status = "cancelled"
	d2, _ := mc.LiveDebt([]mc.CreditState{st}, c.to)
	eqC(t, "sin la despensa cancelada", d2.Total, ct.C(t, c.m, "horizonte.deuda_viva_detalle.despensa_3_pendiente"))
	_, err = mc.LiveDebt([]mc.CreditState{st, st}, c.to)
	wantErr(t, "crédito repetido", err, mc.ErrInvalid)
}

func TestCompraViva(t *testing.T) {
	c := loadCanon(t)
	st := c.in.Credits[0]
	at := func(s string) string {
		r, err := mc.LivePurchases(st, day(t, s))
		if err != nil {
			t.Fatal(err)
		}
		return fmt.Sprint(r)
	}
	// RN-11: con máximo 1, la despensa 2 el 10-dic-2026 viola (d1 sigue viva) y el 16-dic-2026 cumple.
	if at("2026-12-10") != "[d1 d2]" || at("2026-12-15") != "[d1 d2]" || at("2026-12-16") != "[d2]" || at(c.to.String()) != "[d3 d4]" {
		t.Errorf("compras vivas: %s %s %s", at("2026-12-10"), at("2026-12-16"), at(c.to.String()))
	}
}

func TestReparto(t *testing.T) {
	c := loadCanon(t)
	for j, p := range c.ps {
		checkInvariants(t, fmt.Sprintf("despensa %d", j+1), p)
		for _, s := range p.Shares {
			pre := fmt.Sprintf("despensas[%d].%s.", j, s.Code)
			eqC(t, pre+"propio", s.Composition.Own, ct.C(t, c.m, pre+"propio"))
			eqC(t, pre+"libre", s.Composition.Free, ct.C(t, c.m, pre+"libre"))
			eqC(t, pre+"intereses", s.Interest, ct.C(t, c.m, pre+"intereses"))
		}
		if p.Shares[0].External || !p.Shares[1].External {
			t.Errorf("orden de partes: titular primero")
		}
	}
}

func unifiedCredits(t *testing.T, c canon) []mc.CreditState {
	kue := karenAxel(c, []mc.PaymentRecord{})
	cardSt := mc.CreditState{Credit: cardCredit(false), Purchases: []mc.PurchaseState{
		{Ref: "c1", Status: "active", Input: mc.PurchaseInput{Date: day(t, "2026-10-10"), Principal: 50000, TermCount: 1, Shares: holder(50000)}},
		{Ref: "c2", Status: "active", Input: mc.PurchaseInput{Date: day(t, "2026-10-12"), Principal: 90000, TermCount: 3, PromotionID: ptr(int64(1)), Merchant: "Tienda X",
			Shares: []mc.ShareInput{{Code: "axel", Principal: 60000}, {Code: "karen", External: true, Principal: 30000}}}},
		{Ref: "c3", Status: "cancelled", Input: mc.PurchaseInput{Date: day(t, "2026-10-12"), Principal: 1000, TermCount: 1, Shares: holder(1000)}},
	}, Payments: []mc.PaymentRecord{{PurchaseRef: "c2", Seq: 2, Date: day(t, "2026-12-20"), Amount: 10000}}}
	loan := mc.CreditState{Credit: mc.Credit{ID: 3, Kind: mc.KindLoan, Holder: "axel", Frequency: "monthly", AnnualRate: &mc.Rate{Num: 24, Den: 100}, VATOnInterestBP: ptr(1600)},
		Purchases: []mc.PurchaseState{{Ref: "l1", Status: "active", Input: mc.PurchaseInput{Date: day(t, "2026-10-15"), Principal: 1000000, TermCount: 12,
			FirstPaymentDate: ptr(day(t, "2026-11-15")), Shares: holder(1000000)}}}}
	store := mc.CreditState{Credit: storeCredit(true, false), Purchases: []mc.PurchaseState{{Ref: "t1", Status: "active", Input: mc.PurchaseInput{
		Date: day(t, "2026-10-20"), Principal: 100000, TermCount: 3, FirstPaymentDate: ptr(day(t, "2026-11-20")), Shares: holder(100000)}}}}
	return []mc.CreditState{store, loan, cardSt, kue}
}

func checkCalendar(t *testing.T, label string, cal mc.Calendar) mc.Cents {
	t.Helper()
	var sum, fsum mc.Cents
	seen := map[string]bool{}
	for i, e := range cal.Entries {
		sum += e.Amount
		key := fmt.Sprint(e.CreditID, e.Source, e.PurchaseRef, e.Seq, e.StatementDate)
		if seen[key] {
			t.Errorf("%s: duplicado %s", label, key)
		}
		seen[key] = true
		if i > 0 && e.Date.Compare(cal.Entries[i-1].Date) < 0 {
			t.Errorf("%s: fuera de orden en %d", label, i)
		}
		if !sort.StringsAreSorted(e.Provisional) {
			t.Errorf("%s: provisional sin orden", label)
		}
	}
	for _, f := range cal.ByFortnight {
		fsum += f.Total
		var p mc.Cents
		for _, x := range f.ByPerson {
			p += x.Amount
		}
		if p != f.Total {
			t.Errorf("%s: quincena %s por persona %d de %d", label, f.PayDay, p, f.Total)
		}
	}
	if fsum != sum {
		t.Errorf("%s: Σ quincenas %d ≠ Σ entradas %d", label, fsum, sum)
	}
	return sum
}

func TestUnificado(t *testing.T) {
	c := loadCanon(t)
	credits := unifiedCredits(t, c)
	from, to := c.from, day(t, "2027-12-31")
	cal, err := mc.UnifiedCalendar(credits, from, to)
	if err != nil {
		t.Fatal(err)
	}
	sum := checkCalendar(t, "sin corte", cal)
	if len(cal.Entries) != 4+1+3+12+3 {
		t.Errorf("%d entradas", len(cal.Entries))
	}
	kinds := map[mc.Kind]bool{}
	for _, e := range cal.Entries {
		kinds[e.CreditKind] = true
		if e.PurchaseRef == "c2" && e.Seq == 2 && e.Status != "partial" {
			t.Errorf("estado del abono parcial: %s", e.Status)
		}
	}
	if len(kinds) != 4 {
		t.Errorf("tipos: %v", kinds)
	}
	credits[2].Statements = []mc.StatementInput{{StatementDate: day(t, "2026-11-05"), DueDate: day(t, "2026-11-25"), Balance: 81000, Source: "bank"}}
	cal2, err := mc.UnifiedCalendar(credits, from, to)
	if err != nil {
		t.Fatal(err)
	}
	sum2 := checkCalendar(t, "con corte", cal2)
	if len(cal2.Entries) != len(cal.Entries)-1 || sum2 != sum-80000+81000 {
		t.Errorf("el corte no sustituyó: %d entradas, Σ %d", len(cal2.Entries), sum2)
	}
	for _, e := range cal2.Entries {
		if e.Source == "statement" {
			if e.Principal != 80000 || e.Interest != 1000 || e.Shares[0].Amount != 71000 || e.Shares[1].Amount != 10000 || e.Status != "pending" {
				t.Errorf("corte: %+v", e)
			}
			eqList(t, "corte", e.Provisional, []string{"S3-formula-tarjeta"})
		}
	}
	credits[2].Statements[0].Balance = 70000
	_, err = mc.UnifiedCalendar(credits, from, to)
	wantErr(t, "corte menor que sus pagos", err, mc.ErrInvalid)
	credits[2].Statements = append(credits[2].Statements, credits[2].Statements[0])
	_, err = mc.UnifiedCalendar(credits, from, to)
	wantErr(t, "dos cortes de la misma fecha", err, mc.ErrInvalid)
	_, err = mc.UnifiedCalendar(credits, to, from)
	wantErr(t, "ventana invertida", err, mc.ErrInvalid)
}

func TestProvisionalListas(t *testing.T) {
	c := loadCanon(t)
	cal, err := mc.UnifiedCalendar(unifiedCredits(t, c), c.from, day(t, "2027-12-31"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{"d1": {}, "c1": {"S3-formula-tarjeta", "S3-origen-cargo-normal"},
		"c2": {"S3-formula-tarjeta", "S3-serie-sin-intereses"}, "l1": {"S3-convencion-prestamo"}, "t1": {"S3-serie-sin-intereses"}}
	for _, e := range cal.Entries {
		eqList(t, e.PurchaseRef, e.Provisional, want[e.PurchaseRef])
	}
	for _, p := range c.ps {
		eqList(t, "Kueski", p.Provisional, []string{})
	}
}
