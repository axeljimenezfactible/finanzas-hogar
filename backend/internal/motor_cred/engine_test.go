package motor_cred_test

// Pruebas sin canon: datos en línea con las cifras de S3-interfaz.md (capturas de fixture, no de
// Axel). Las que usan cifras canon las leen de model.json al correr (canontest).

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"testing"

	mc "finanzas-hogar/internal/motor_cred"
)

func ptr[T any](v T) *T { return &v }

func day(t *testing.T, s string) mc.Date {
	t.Helper()
	d, err := mc.ParseDate(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func wantErr(t *testing.T, label string, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("%s: error %v, se esperaba %v", label, err, target)
	}
}

func eqList(t *testing.T, label string, got, want []string) {
	t.Helper()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("%s: %v, se esperaba %v", label, got, want)
	}
}

func holder(p mc.Cents) []mc.ShareInput { return []mc.ShareInput{{Code: "axel", Principal: p}} }

// checkInvariants: invariantes de §5 (TestReparto).
func checkInvariants(t *testing.T, label string, p mc.Purchase) {
	t.Helper()
	if p.Total != p.Principal+p.Interest {
		t.Errorf("%s: total %d ≠ capital + intereses", label, p.Total)
	}
	var sa, sp mc.Cents
	per := map[string]mc.Cents{}
	for _, s := range p.Payments {
		sa, sp = sa+s.Amount, sp+s.Principal
		if s.Amount != s.Principal+s.Interest {
			t.Errorf("%s: pago %d no cuadra", label, s.Seq)
		}
		var ss mc.Cents
		for _, x := range s.Shares {
			ss += x.Amount
			per[x.Code] += x.Amount
		}
		if len(s.Shares) > 0 && ss != s.Amount {
			t.Errorf("%s: partes del pago %d suman %d de %d", label, s.Seq, ss, s.Amount)
		}
	}
	if sa != p.Total || sp != p.Principal {
		t.Errorf("%s: Σ montos %d (total %d), Σ capital %d (capital %d)", label, sa, p.Total, sp, p.Principal)
	}
	for _, s := range p.Shares {
		if per[s.Code] != s.Total {
			t.Errorf("%s: pagos de %s suman %d de %d", label, s.Code, per[s.Code], s.Total)
		}
	}
	if !sort.StringsAreSorted(p.Provisional) {
		t.Errorf("%s: provisional sin orden %v", label, p.Provisional)
	}
}

func cardCredit(k18 bool) mc.Credit {
	c := mc.Credit{ID: 2, Name: "Tarjeta", Kind: mc.KindCard, Holder: "axel", Limit: ptr(mc.Cents(2000000)),
		Frequency: "monthly", StatementDay: ptr(5), PaymentDueDay: ptr(25),
		Promotions: []mc.Promotion{
			{ID: 1, Merchant: "Tienda X", TermCount: 3, Rate: mc.Rate{Num: 0, Den: 1, Decimals: 0}, ValidFrom: mc.Date{Year: 2026, Month: 1, Day: 1}, ValidTo: mc.Date{Year: 2027, Month: 12, Day: 31}, Active: true},
			{ID: 2, Merchant: "Tienda X", TermCount: 6, Rate: mc.Rate{Num: 10, Den: 100, Decimals: 2}, ValidFrom: mc.Date{Year: 2026, Month: 1, Day: 1}, ValidTo: mc.Date{Year: 2027, Month: 12, Day: 31}, Active: true},
		}}
	if k18 {
		c.AnnualRate, c.RateBase = &mc.Rate{Num: 455, Den: 1000, Decimals: 1}, "daily"
		c.VATOnInterestBP, c.MinPaymentPctBP, c.MinPaymentFixed = ptr(1600), ptr(500), ptr(mc.Cents(20000))
	}
	return c
}

func statement(t *testing.T, balance mc.Cents) mc.StatementInput {
	daily := make([]mc.Cents, 30)
	for i := range daily {
		daily[i] = balance
	}
	return mc.StatementInput{StatementDate: day(t, "2026-11-05"), DueDate: day(t, "2026-11-25"), Balance: balance, DailyBalances: daily, Source: "bank"}
}

func TestTarjetaPagoTotal(t *testing.T) {
	st := statement(t, 100000)
	for _, k18 := range []bool{true, false} {
		c := cardCredit(k18)
		for _, paid := range []string{"2026-11-24", "2026-11-25"} {
			v, prov, err := c.CardInterest(st, 100000, day(t, paid))
			if err != nil || v != 0 || len(prov) != 0 {
				t.Errorf("k18=%v pagado %s: %d %v %v", k18, paid, v, prov, err)
			}
		}
		p, err := mc.QuotePurchase(mc.PurchaseInput{Credit: c, Date: day(t, "2026-10-10"), Principal: 50000, TermCount: 1, Shares: holder(50000)})
		if err != nil || p.Interest != 0 || p.Total != 50000 || len(p.Payments) != 1 || p.Payments[0].Due != day(t, "2026-11-25") {
			t.Fatalf("cargo normal: %+v %v", p, err)
		}
		if p.Source != "quoted" || p.PurchaseKind != "single_charge" || p.RateBasis != "" || p.Rate != (mc.Rate{Num: 0, Den: 1, Decimals: 0}) {
			t.Errorf("cargo normal: %s %s %q %v", p.Source, p.PurchaseKind, p.RateBasis, p.Rate)
		}
		eqList(t, "cargo normal", p.Provisional, []string{"S3-formula-tarjeta", "S3-origen-cargo-normal"})
		checkInvariants(t, "cargo normal", p)
	}
}

func TestTarjetaMSI(t *testing.T) {
	c := cardCredit(false)
	in := mc.PurchaseInput{Credit: c, Date: day(t, "2026-10-10"), Principal: 100000, TermCount: 3, PromotionID: ptr(int64(1)), Merchant: "tienda x", Shares: holder(100000)}
	p, err := mc.QuotePurchase(in)
	if err != nil {
		t.Fatal(err)
	}
	if p.Interest != 0 || p.Payments[0].Amount != 33333 || p.Payments[1].Amount != 33333 || p.Payments[2].Amount != 33334 {
		t.Errorf("MSI: %+v", p.Payments)
	}
	if p.Payments[2].Due != day(t, "2027-01-25") || p.Source != "promotion" {
		t.Errorf("MSI: fechas o source %v %s", p.Payments[2].Due, p.Source)
	}
	eqList(t, "MSI", p.Provisional, []string{"S3-formula-tarjeta", "S3-serie-sin-intereses"})
	checkInvariants(t, "MSI", p)
	bad := in
	bad.TermCount, bad.PromotionID = 6, ptr(int64(2))
	_, err = mc.QuotePurchase(bad)
	wantErr(t, "MSI con intereses", err, mc.ErrUnsupported)
	bad.PromotionID = nil
	_, err = mc.QuotePurchase(bad)
	wantErr(t, "plazo sin promoción", err, mc.ErrNoRate)
	bad.QuotedInterest = ptr(mc.Cents(1))
	_, err = mc.QuotePurchase(bad)
	wantErr(t, "tarjeta con cotización", err, mc.ErrInvalid)
	bad = in
	bad.Merchant = "Otra"
	_, err = mc.QuotePurchase(bad)
	wantErr(t, "promoción de otra tienda", err, mc.ErrPromotionNotApplicable)
}

func TestTarjetaSinCaptura(t *testing.T) {
	c := cardCredit(false)
	_, _, err := c.CardInterest(statement(t, 100000), 100000, day(t, "2026-11-26"))
	wantErr(t, "interés sin datos", err, mc.ErrMissingCapture)
	_, err = c.CardMinimumPayment(statement(t, 100000))
	wantErr(t, "mínimo sin datos", err, mc.ErrMissingCapture)
	_, _, err = c.CardMinimum(statement(t, 100000))
	wantErr(t, "mínimo sin datos", err, mc.ErrMissingCapture)
}

func TestTarjetaMinimo(t *testing.T) {
	c := cardCredit(true)
	st := statement(t, 100000)
	// Diaria: Σ saldos 3,000,000 × 0.455 / 360 = 3,791.67 → 3,792; IVA 16 % = 606.72 → 607; total 4,399.
	v, prov, err := c.CardInterest(st, 0, day(t, "2026-11-26"))
	if err != nil || v != 4399 {
		t.Errorf("interés diario: %d %v", v, err)
	}
	eqList(t, "interés", prov, []string{"S3-formula-tarjeta", "U4-iva-credito"})
	if v, err := c.CardPeriodInterest(st, 100000, day(t, "2026-11-26")); err != nil || v != 4399 {
		t.Errorf("pagado tarde: %d %v", v, err)
	}
	c.RateBase, c.VATOnInterestBP = "monthly", nil // mensual sin IVA: 100,000 × 0.455 / 12 = 3,791.67 → 3,792
	if v, _, err := c.CardInterest(st, 0, day(t, "2026-11-26")); err != nil || v != 3792 {
		t.Errorf("interés mensual: %d %v", v, err)
	}
	m, prov, err := c.CardMinimum(st) // máx(5 % = 5,000, fijo 20,000) = 20,000
	if err != nil || m != 20000 {
		t.Errorf("mínimo: %d %v", m, err)
	}
	eqList(t, "mínimo", prov, []string{"S3-formula-tarjeta", "U4-iva-credito"})
	if m, _ := c.CardMinimumPayment(statement(t, 15000)); m != 15000 {
		t.Errorf("mínimo topado al saldo: %d", m)
	}
	c.RateBase = "daily"
	st.DailyBalances = nil
	_, _, err = c.CardInterest(st, 0, day(t, "2026-11-26"))
	wantErr(t, "diaria sin saldos", err, mc.ErrInvalid)
	c.RateBase = "anual"
	_, _, err = c.CardInterest(st, 0, day(t, "2026-11-26"))
	wantErr(t, "base inválida", err, mc.ErrInvalid)
	_, err = c.CardMinimumPayment(mc.StatementInput{Balance: -1})
	wantErr(t, "saldo negativo", err, mc.ErrInvalid)
	_, _, err = c.CardInterest(mc.StatementInput{Balance: -1}, 0, day(t, "2026-11-26"))
	wantErr(t, "saldo negativo", err, mc.ErrInvalid)
	c.Kind = mc.KindLoan
	_, _, err = c.CardInterest(st, 0, day(t, "2026-11-26"))
	wantErr(t, "no es tarjeta", err, mc.ErrInvalid)
	_, err = c.CardMinimumPayment(st)
	wantErr(t, "no es tarjeta", err, mc.ErrInvalid)
}

func checkLoan(t *testing.T, label string, p mc.Purchase, principal mc.Cents) {
	t.Helper()
	checkInvariants(t, label, p)
	var bal = principal
	for i, s := range p.Payments {
		bal -= s.Principal
		if i < len(p.Payments)-1 && s.Amount != p.ReferencePayment {
			t.Errorf("%s: pago %d = %d, fijo %d", label, s.Seq, s.Amount, p.ReferencePayment)
		}
	}
	if bal != 0 {
		t.Errorf("%s: saldo final %d", label, bal)
	}
}

func TestPrestamoFrances(t *testing.T) {
	r := mc.Rate{Num: 24, Den: 100, Decimals: 0}
	p, err := mc.Amortize(1000000, r, ptr(1600), 12, "monthly", day(t, "2027-01-31"))
	if err != nil {
		t.Fatal(err)
	}
	checkLoan(t, "12 mensual", p, 1000000)
	eqList(t, "préstamo con IVA", p.Provisional, []string{"S3-convencion-prestamo"})
	if p.Payments[1].Due != day(t, "2027-02-28") || p.Source != "amortization" || p.PurchaseKind != "loan_disbursement" || p.RateBasis != "annual" {
		t.Errorf("fechas o etiquetas: %v %s", p.Payments[1].Due, p.Source)
	}
	last := p.Payments[11]
	if last.Amount == p.ReferencePayment {
		t.Logf("el último pago coincidió con el fijo")
	}
	p, err = mc.Amortize(100000, r, nil, 7, "fortnightly", day(t, "2026-10-15"))
	if err != nil {
		t.Fatal(err)
	}
	checkLoan(t, "7 quincenal", p, 100000)
	if p.Payments[6].Amount == p.ReferencePayment {
		t.Errorf("100,000 a 7 debía obligar ajuste en el último pago")
	}
	eqList(t, "préstamo sin IVA", p.Provisional, []string{"S3-convencion-prestamo", "U4-iva-credito"})
	if p.Payments[1].Due != day(t, "2026-10-30") {
		t.Errorf("quincena: %v", p.Payments[1].Due)
	}
	p, err = mc.Amortize(70000, mc.Rate{Num: 0, Den: 1}, nil, 3, "weekly", day(t, "2026-10-15"))
	if err != nil {
		t.Fatal(err)
	}
	checkLoan(t, "sin tasa semanal", p, 70000)
	if p.Payments[2].Due != day(t, "2026-10-29") {
		t.Errorf("semanal: %v", p.Payments[2].Due)
	}
	for _, f := range []string{"daily", ""} {
		_, err = mc.Amortize(100000, r, nil, 7, f, day(t, "2026-10-15"))
		wantErr(t, "frecuencia "+f, err, mc.ErrInvalid)
	}
	_, err = mc.Amortize(100000, r, ptr(-1), 7, "monthly", day(t, "2026-10-15"))
	wantErr(t, "IVA negativo", err, mc.ErrInvalid)
}

func storeCredit(rates bool, k18 bool) mc.Credit {
	c := mc.Credit{ID: 4, Name: "Tienda", Kind: mc.KindStore, Holder: "axel", Frequency: "monthly", Promotions: []mc.Promotion{}}
	if rates {
		c.TermRates = []mc.TermRate{{TermCount: 3, Rate: mc.Rate{Num: 0, Den: 1, Decimals: 0}, Source: "manual"},
			{TermCount: 6, Rate: mc.Rate{Num: 12, Den: 100, Decimals: 0}, Source: "manual"}}
	}
	if k18 {
		c.AnnualRate, c.RateBase = &mc.Rate{Num: 455, Den: 1000, Decimals: 1}, "monthly"
		c.MinPaymentPctBP, c.MinPaymentFixed = ptr(500), ptr(mc.Cents(20000))
	}
	return c
}

func TestTiendaAbonos(t *testing.T) {
	in := mc.PurchaseInput{Credit: storeCredit(true, false), Date: day(t, "2026-10-10"), Principal: 100000, TermCount: 3,
		FirstPaymentDate: ptr(day(t, "2026-11-20")), Shares: holder(100000)}
	p, err := mc.QuotePurchase(in)
	if err != nil {
		t.Fatal(err)
	}
	if p.Interest != 0 || p.Payments[0].Amount != 33333 || p.Payments[2].Amount != 33334 || p.Source != "term_table" || p.Payments[2].Principal != 33334 {
		t.Errorf("abonos sin interés: %+v", p.Payments)
	}
	eqList(t, "tienda", p.Provisional, []string{"S3-serie-sin-intereses"})
	checkInvariants(t, "tienda", p)
	in.TermCount, in.Principal, in.Shares = 6, 600000, holder(600000)
	if p, err = mc.QuotePurchase(in); err != nil || p.Interest != 72000 || len(p.Provisional) != 0 {
		t.Errorf("tienda con tasa por plazo: %d %v %v", p.Interest, p.Provisional, err)
	}
	checkInvariants(t, "tienda con tasa", p)
	in.Credit.Frequency, in.FirstPaymentDate, in.Credit.FirstPaymentMinDays = "fortnightly", nil, ptr(15)
	if p, err = mc.QuotePurchase(in); err != nil || p.Payments[0].Due != day(t, "2026-10-30") {
		t.Errorf("tienda quincenal: %v", err)
	}
}

func TestTiendaSinCaptura(t *testing.T) {
	in := mc.PurchaseInput{Credit: storeCredit(false, false), Date: day(t, "2026-10-10"), Principal: 100000, TermCount: 3,
		FirstPaymentDate: ptr(day(t, "2026-11-20")), Shares: holder(100000)}
	_, err := mc.QuotePurchase(in)
	wantErr(t, "tienda sin datos", err, mc.ErrMissingCapture)
	in.FirstPaymentDate = nil
	_, err = mc.QuotePurchase(in)
	wantErr(t, "tienda mensual sin primera fecha", err, mc.ErrInvalid)
	in.Credit.Frequency = "yearly"
	_, err = mc.QuotePurchase(in)
	wantErr(t, "tienda sin frecuencia", err, mc.ErrInvalid)
}

func TestTiendaConInteres(t *testing.T) {
	c := storeCredit(false, true)
	in := mc.PurchaseInput{Credit: c, Date: day(t, "2026-10-10"), Principal: 100000, TermCount: 3,
		FirstPaymentDate: ptr(day(t, "2026-11-20")), Shares: holder(100000)}
	_, err := mc.QuotePurchase(in)
	wantErr(t, "tienda con datos K18", err, mc.ErrUnsupported)
	in.PromotionID = ptr(int64(9))
	_, err = mc.QuotePurchase(in)
	wantErr(t, "promoción inexistente", err, mc.ErrPromotionNotApplicable)
	st := mc.StatementInput{StatementDate: day(t, "2026-11-05"), DueDate: day(t, "2026-11-25"), Balance: 100000, Source: "computed"}
	if v, prov, err := c.CardInterest(st, 0, day(t, "2026-11-26")); err != nil || v != 3792 || len(prov) != 2 {
		t.Errorf("interés de tienda: %d %v %v", v, prov, err)
	}
	if m, err := c.CardMinimumPayment(st); err != nil || m != 20000 {
		t.Errorf("mínimo de tienda: %d %v", m, err)
	}
	in.PromotionID, in.QuotedInterest = nil, ptr(mc.Cents(3000))
	if p, err := mc.QuotePurchase(in); err != nil || p.Interest != 3000 || p.Source != "quoted" {
		t.Errorf("tienda cotizada: %v", err)
	}
}

func TestDespachoCotizacion(t *testing.T) {
	base := mc.PurchaseInput{Credit: mc.Credit{ID: 1, Kind: mc.KindRevolving, Limit: ptr(mc.Cents(100000)), Frequency: "fortnightly",
		FirstPaymentMinDays: ptr(15), TermRates: []mc.TermRate{{TermCount: 2, Rate: mc.Rate{Num: 1, Den: 10, Decimals: 0}}}},
		Date: day(t, "2026-10-10"), Principal: 1000, TermCount: 2, Shares: holder(1000)}
	cases := []struct {
		name string
		mod  func(*mc.PurchaseInput)
		err  error
	}{
		{"tipo desconocido", func(p *mc.PurchaseInput) { p.Credit.Kind = "x" }, mc.ErrInvalid},
		{"capital 0", func(p *mc.PurchaseInput) { p.Principal = 0 }, mc.ErrInvalid},
		{"primer pago antes de la compra", func(p *mc.PurchaseInput) { p.FirstPaymentDate = ptr(day(t, "2026-10-01")) }, mc.ErrInvalid},
		{"dos titulares", func(p *mc.PurchaseInput) {
			p.Shares = []mc.ShareInput{{Code: "a", Principal: 500}, {Code: "b", Principal: 500}}
		}, mc.ErrInvalid},
		{"partes que no suman", func(p *mc.PurchaseInput) { p.Shares = holder(999) }, mc.ErrInvalid},
		{"cotización negativa", func(p *mc.PurchaseInput) { p.QuotedInterest = ptr(mc.Cents(-1)) }, mc.ErrInvalid},
		{"rate_decimals negativo", func(p *mc.PurchaseInput) { p.QuotedInterest, p.RateDecimals = ptr(mc.Cents(1)), ptr(-1) }, mc.ErrInvalid},
		{"titular negativo", func(p *mc.PurchaseInput) {
			p.QuotedInterest = ptr(mc.Cents(1))
			p.Shares = []mc.ShareInput{{Code: "axel", Principal: 0}, {Code: "karen", External: true, Principal: 1000}}
		}, mc.ErrInvalid},
		{"línea sin límite", func(p *mc.PurchaseInput) { p.Credit.Limit = nil }, mc.ErrInvalid},
		{"línea semanal", func(p *mc.PurchaseInput) { p.Credit.Frequency = "weekly" }, mc.ErrUnsupported},
		{"línea sin días mínimos", func(p *mc.PurchaseInput) { p.Credit.FirstPaymentMinDays = nil }, mc.ErrInvalid},
		{"promoción no aplicable", func(p *mc.PurchaseInput) { p.PromotionID = ptr(int64(1)) }, mc.ErrPromotionNotApplicable},
		{"línea sin tasa", func(p *mc.PurchaseInput) { p.TermCount = 6 }, mc.ErrNoRate},
		{"préstamo sin primera fecha", func(p *mc.PurchaseInput) { p.Credit.Kind = mc.KindLoan }, mc.ErrInvalid},
		{"préstamo sin tasa anual", func(p *mc.PurchaseInput) {
			p.Credit.Kind, p.FirstPaymentDate = mc.KindLoan, ptr(day(t, "2026-11-10"))
		}, mc.ErrMissingCapture},
		{"préstamo con cotización", func(p *mc.PurchaseInput) {
			p.Credit.Kind, p.FirstPaymentDate, p.QuotedInterest = mc.KindLoan, ptr(day(t, "2026-11-10")), ptr(mc.Cents(1))
		}, mc.ErrInvalid},
		{"tarjeta sin corte", func(p *mc.PurchaseInput) { p.Credit.Kind = mc.KindCard }, mc.ErrInvalid},
	}
	for _, c := range cases {
		in := base
		c.mod(&in)
		_, err := mc.QuotePurchase(in)
		wantErr(t, c.name, err, c.err)
	}
	p, err := mc.QuotePurchase(base)
	if err != nil || p.Source != "term_table" || p.PurchaseKind != "installment_plan" || p.RateBasis != "term_total" || p.Interest != 100 {
		t.Errorf("línea por tabla: %+v %v", p, err)
	}
	in := base
	in.FirstPaymentDate, in.Credit.Frequency = ptr(day(t, "2026-10-31")), "monthly"
	if p, err = mc.QuotePurchase(in); err != nil || p.Payments[1].Due != day(t, "2026-11-30") {
		t.Errorf("línea mensual con primera fecha: %v", err)
	}
	loan := base
	loan.Credit = mc.Credit{ID: 3, Kind: mc.KindLoan, Holder: "axel", Frequency: "monthly", AnnualRate: &mc.Rate{Num: 24, Den: 100}, VATOnInterestBP: ptr(1600)}
	loan.FirstPaymentDate, loan.TermCount, loan.Principal = ptr(day(t, "2026-11-10")), 6, 600000
	loan.Shares = []mc.ShareInput{{Code: "axel", Principal: 400000}, {Code: "karen", External: true, Principal: 200000}}
	p, err = mc.QuotePurchase(loan)
	if err != nil || len(p.Shares) != 2 || len(p.Payments[0].Shares) != 2 || p.PromotionID != nil {
		t.Fatalf("préstamo con partes: %v", err)
	}
	checkInvariants(t, "préstamo con partes", p)
}

func TestPromocionFronteras(t *testing.T) {
	c := mc.Credit{ID: 1, Kind: mc.KindRevolving, Holder: "axel", Limit: ptr(mc.Cents(1000000)), Frequency: "fortnightly", FirstPaymentMinDays: ptr(15),
		TermRates: []mc.TermRate{{TermCount: 4, Rate: mc.Rate{Num: 1, Den: 4, Decimals: 0}}},
		Promotions: []mc.Promotion{{ID: 1, Merchant: "Chedraui", TermCount: 4, Rate: mc.Rate{Num: 0, Den: 1, Decimals: 0}, MinPurchase: 200000,
			ValidFrom: day(t, "2026-10-28"), ValidTo: day(t, "2026-10-30"), Active: true}}}
	try := func(date string, principal mc.Cents) (mc.Purchase, error, int) {
		in := mc.PurchaseInput{Credit: c, Date: day(t, date), Principal: principal, TermCount: 4, PromotionID: ptr(int64(1)), Merchant: " CHEDRAUI ", Shares: holder(principal)}
		p, err := mc.QuotePurchase(in)
		return p, err, len(mc.ApplicablePromotions(c, in.Date, principal, 4, nil, in.Merchant))
	}
	for _, d := range []string{"2026-10-28", "2026-10-30"} {
		p, err, n := try(d, 300001)
		if err != nil || n != 1 || p.Interest != 0 || p.Payments[0].Amount != 75000 || p.Payments[3].Amount != 75001 {
			t.Errorf("promoción el %s: %v %d %+v", d, err, n, p.Payments)
		}
		eqList(t, "promoción", p.Provisional, []string{"S3-serie-sin-intereses"})
		checkInvariants(t, "promoción", p)
	}
	for _, x := range []struct {
		d string
		p mc.Cents
	}{{"2026-10-27", 300000}, {"2026-10-28", 199999}} {
		_, err, n := try(x.d, x.p)
		wantErr(t, "fuera de frontera "+x.d, err, mc.ErrPromotionNotApplicable)
		if n != 0 {
			t.Errorf("%s: %d promociones aplicables", x.d, n)
		}
		in := mc.PurchaseInput{Credit: c, Date: day(t, x.d), Principal: x.p, TermCount: 4, Shares: holder(x.p)}
		if p, err := mc.QuotePurchase(in); err != nil || p.Interest != mc.ApplyRate(x.p, c.TermRates[0].Rate) {
			t.Errorf("tasa normal sin promoción: %v", err)
		}
	}
	c.Promotions[0].StoreID = ptr(int64(7))
	if n := len(mc.ApplicablePromotions(c, day(t, "2026-10-28"), 300000, 4, ptr(int64(8)), "Chedraui")); n != 0 {
		t.Errorf("otra tienda por store_id: %d", n)
	}
}

func TestErroresYTipos(t *testing.T) {
	if _, err := mc.ParseDate("2027-02-30"); !errors.Is(err, mc.ErrInvalid) {
		t.Errorf("2027-02-30 aceptado")
	}
	if _, err := mc.ParseDate("2027-2-3"); !errors.Is(err, mc.ErrInvalid) {
		t.Errorf("fecha corta aceptada")
	}
	var d mc.Date
	if err := json.Unmarshal([]byte(`5`), &d); !errors.Is(err, mc.ErrInvalid) {
		t.Errorf("fecha numérica aceptada")
	}
	if b, _ := json.Marshal(mc.Date{}); string(b) != "null" {
		t.Errorf("fecha vacía: %s", b)
	}
	if f := mc.FortnightOf(day(t, "2027-02-16")); f.Half != 2 || f.PayDay != day(t, "2027-02-28") {
		t.Errorf("quincena de febrero: %+v", f)
	}
	if mc.DaysBetween(day(t, "2026-10-10"), day(t, "2026-10-30")) != 20 || day(t, "2026-12-31").AddDays(1) != day(t, "2027-01-01") {
		t.Errorf("aritmética de fechas")
	}
	for err, code := range map[error]string{nil: "", mc.ErrUnknownField: "unknown_field", mc.ErrInvalid: "invalid", mc.ErrNoRate: "no_rate",
		mc.ErrMissingCapture: "missing_capture", mc.ErrPromotionNotApplicable: "promotion_not_applicable",
		mc.ErrUnsupported: "unsupported", mc.ErrTimeout: "timeout", errors.New("otro"): "invalid"} {
		if got := mc.Code(err); got != code {
			t.Errorf("Code(%v) = %q", err, got)
		}
	}
	if s, c, ok := mc.ResponseFor(mc.ErrTimeout); s != 503 || c != "unavailable" || !ok {
		t.Errorf("503: %d %s", s, c)
	}
	if s, c, ok := mc.ResponseFor(mc.ErrInvalid); s != 422 || c != "validation_failed" || !ok {
		t.Errorf("422: %d %s", s, c)
	}
	if _, _, ok := mc.ResponseFor(mc.ErrNoRate); ok {
		t.Errorf("ErrNoRate no tiene respuesta fija")
	}
	if _, err := mc.NewRate(1, 0, 2); !errors.Is(err, mc.ErrInvalid) {
		t.Errorf("tasa con den 0")
	}
	r, err := mc.NewRate(53450, 309450, 2)
	if err != nil || r.Display() != "17.27 %" || (mc.Rate{Num: 1, Den: 3, Decimals: 0}).Display() != "33 %" ||
		(mc.Rate{Num: 1, Den: 3, Decimals: -1}).Display() != "" || (mc.Rate{Num: 1, Den: 1000, Decimals: 2}).Display() != "0.10 %" {
		t.Errorf("Display: %q", r.Display())
	}
	if mc.ApplyRate(100, mc.Rate{}) != 0 || mc.SplitTrunc(1, 0) != nil || mc.SplitHalfUp(1, 0) != nil || mc.ReferencePayment(nil, 1) != 0 {
		t.Errorf("casos vacíos")
	}
	_, err = mc.DecodeInput(strings.NewReader(`{"today": "2026-10-10", "otra": 1}`))
	if mc.Code(err) != "unknown_field" {
		t.Errorf("campo desconocido: %v", err)
	}
	for _, s := range []string{`{"today": "2026-02-30"}`, `{`, `{"credits": [{"credit": {"id": 1}}, {"credit": {"id": 1}}]}`,
		`{"credits": [{"credit": {"id": 1}, "payments": [{"purchase_ref": "x", "seq": 1}]}]}`,
		`{"credits": [{"credit": {"id": 1, "kind": "revolving_line"}, "statements": [{"balance_cents": 1}]}]}`,
		`{"credits": [{"credit": {"id": 1}, "purchases": [{"ref": "a", "status": "x"}]}]}`,
		`{"credits": [{"credit": {"id": 1}, "purchases": [{"ref": "a", "status": "active"}, {"ref": "a", "status": "active"}]}]}`} {
		if _, err := mc.DecodeInput(strings.NewReader(s)); mc.Code(err) != "invalid" {
			t.Errorf("%s: %v", s, err)
		}
	}
	in, err := mc.DecodeInput(strings.NewReader(`{"credits": [{"credit": {"id": 1}}]}`))
	if _, err2 := in.CreditByID(2); err != nil || !errors.Is(err2, mc.ErrInvalid) {
		t.Errorf("crédito inexistente")
	}
	p := mc.Purchase{TotalExact: nil}
	if b, _ := json.Marshal(p); !strings.Contains(string(b), `"total_exact":null`) {
		t.Errorf("total_exact nulo: %s", b)
	}
}
