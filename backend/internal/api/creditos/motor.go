package creditos

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// snapshot es todo lo de créditos leído de la base en una sola pasada (escala de un hogar). Con él se
// arma la entrada del motor S3 (motor_cred.CreditState), que hace todos los cálculos.
type snapshot struct {
	today      dominio.Date
	credits    []dominio.CreditRow
	byID       map[int64]dominio.CreditRow
	persons    map[int64]dominio.PersonRow
	rates      map[int64][]dominio.TermRateRow
	promos     []dominio.PromotionRow
	purchases  []dominio.PurchaseRow
	shares     map[int64][]dominio.PurchaseShareRow
	sched      []dominio.ScheduledRow
	paidBySP   map[int64]int64
	spShares   map[int64][]dominio.ScheduledShareRow
	payments   []dominio.PaymentRow
	statements []dominio.StatementRow
}

func load(ctx context.Context, q dominio.Querier, today dominio.Date) (*snapshot, error) {
	s := &snapshot{today: today, byID: map[int64]dominio.CreditRow{}, rates: map[int64][]dominio.TermRateRow{}, paidBySP: map[int64]int64{}}
	var err error
	if s.credits, err = dominio.QueryCredits(ctx, q, "", "ORDER BY c.id"); err != nil {
		return nil, err
	}
	for _, c := range s.credits {
		s.byID[c.ID] = c
	}
	if s.persons, err = dominio.CreditPersons(ctx, q); err != nil {
		return nil, err
	}
	rs, err := dominio.QueryTermRates(ctx, q, "1 = 1", "ORDER BY credit_id, term_count")
	if err != nil {
		return nil, err
	}
	for _, r := range rs {
		s.rates[r.CreditID] = append(s.rates[r.CreditID], r)
	}
	if s.promos, err = dominio.QueryPromotions(ctx, q, "1 = 1", "ORDER BY id"); err != nil {
		return nil, err
	}
	if s.purchases, err = dominio.QueryPurchases(ctx, q, "", "ORDER BY p.purchase_date, p.id"); err != nil {
		return nil, err
	}
	if s.shares, err = dominio.PurchaseShares(ctx, q, 0); err != nil {
		return nil, err
	}
	if s.sched, err = dominio.QueryScheduled(ctx, q, "", "ORDER BY sp.due_date, sp.id"); err != nil {
		return nil, err
	}
	if s.payments, err = dominio.QueryPayments(ctx, q, ""); err != nil {
		return nil, err
	}
	for _, m := range s.payments {
		if m.ScheduledPaymentID != nil {
			s.paidBySP[*m.ScheduledPaymentID] += m.AmountCents
		}
	}
	if s.spShares, err = dominio.ScheduledShares(ctx, q, 0); err != nil {
		return nil, err
	}
	s.statements, err = dominio.QueryStatements(ctx, q)
	return s, err
}

func mdate(s string) mc.Date {
	d, _ := mc.ParseDate(s)
	return d
}

func mdateOf(d dominio.Date) mc.Date { return mc.Date{Year: d.Year, Month: d.Month, Day: d.Day} }

func mdatePtr(s *string) *mc.Date {
	if s == nil {
		return nil
	}
	d := mdate(*s)
	return &d
}

func intPtr(v *int64) *int {
	if v == nil {
		return nil
	}
	n := int(*v)
	return &n
}

func centsPtr(v *int64) *mc.Cents {
	if v == nil {
		return nil
	}
	c := mc.Cents(*v)
	return &c
}

func decimals(v *int64) int {
	if v == nil {
		return -1
	}
	return int(*v)
}

func ref(purchaseID int64) string { return fmt.Sprintf("p%09d", purchaseID) }

// ratePct es la tasa como se capturó (K12): la redondea el motor (Rate.Display); sin decimales, ausente.
func ratePct(num, den int64, dec *int64) *string {
	if dec == nil || den <= 0 {
		return nil
	}
	s := strings.TrimSpace(strings.TrimSuffix(mc.Rate{Num: num, Den: den, Decimals: int(*dec)}.Display(), "%"))
	if s == "" {
		return nil
	}
	return &s
}

func (s *snapshot) motorCredit(c dominio.CreditRow) mc.Credit {
	out := mc.Credit{ID: c.ID, Name: c.Name, Kind: mc.Kind(c.Kind), Holder: s.persons[c.HolderPersonID].Code,
		HolderPersonID: c.HolderPersonID, Limit: centsPtr(c.CreditLimitCents), Frequency: c.PaymentFrequency,
		FirstPaymentMinDays: intPtr(c.FirstPaymentMinDays), Reading: mc.Reading(c.AvailabilityReading),
		ReadingConfirmed: c.AvailabilityConfirmed, ReadingAAnchor: mdatePtr(c.ReadingAAnchorDate),
		StatementDay: intPtr(c.StatementDay), PaymentDueDay: intPtr(c.PaymentDueDay),
		MinPaymentPctBP: intPtr(c.MinPaymentPctBp), MinPaymentFixed: centsPtr(c.MinPaymentFixedCents)}
	if c.AnnualRateNum != nil && c.AnnualRateDen != nil {
		out.AnnualRate = &mc.Rate{Num: *c.AnnualRateNum, Den: *c.AnnualRateDen, Decimals: decimals(c.AnnualRateDecimals)}
	}
	if c.RateBase != nil {
		out.RateBase = *c.RateBase
	}
	// vat_included = 1: la tasa ya trae IVA (el motor no suma nada); si no, el porcentaje capturado (Q2-67).
	if c.VatIncluded == nil || !*c.VatIncluded {
		out.VATOnInterestBP = intPtr(c.VatOnInterestBp)
	}
	for _, r := range s.rates[c.ID] {
		out.TermRates = append(out.TermRates, mc.TermRate{TermCount: int(r.TermCount),
			Rate: mc.Rate{Num: r.RateNum, Den: r.RateDen, Decimals: decimals(r.RateDecimals)}, Source: r.Source, QuoteDate: mdatePtr(r.QuoteDate)})
	}
	for _, p := range s.promos {
		if p.CreditID != nil && *p.CreditID != c.ID {
			continue
		}
		out.Promotions = append(out.Promotions, mc.Promotion{ID: p.ID, CreditID: p.CreditID, Merchant: p.MerchantName,
			StoreID: p.StoreID, TermCount: int(p.TermCount), Rate: mc.Rate{Num: p.RateNum, Den: p.RateDen, Decimals: decimals(p.RateDecimals)},
			MinPurchase: mc.Cents(p.MinPurchaseCents), ValidFrom: mdate(p.ValidFrom), ValidTo: mdate(p.ValidTo), Active: p.IsActive})
	}
	return out
}

func (s *snapshot) shareInput(personID, principal int64) mc.ShareInput {
	p := s.persons[personID]
	return mc.ShareInput{PersonID: personID, Code: p.Code, External: p.Participation == "external", Principal: mc.Cents(principal)}
}

// state arma el CreditState de un crédito. Una compra guardada se vuelve a cotizar con sus intereses
// ya fijados (quoted) y su primer pago, así el motor reproduce la serie guardada sin depender de que la
// promoción o la tasa sigan vigentes (préstamo: la tabla de amortización con su tasa anual).
func (s *snapshot) state(c dominio.CreditRow) mc.CreditState {
	cr := s.motorCredit(c)
	st := mc.CreditState{Credit: cr}
	for _, p := range s.purchases {
		if p.CreditID != c.ID {
			continue
		}
		fp := mdate(p.FirstPaymentDate)
		in := mc.PurchaseInput{Credit: cr, Date: mdate(p.PurchaseDate), Principal: mc.Cents(p.PrincipalCents),
			TermCount: int(p.TermCount), RateDecimals: intPtr(p.RateDecimals), FirstPaymentDate: &fp}
		if p.InterestSource != "amortization" {
			qi := mc.Cents(p.InterestCents)
			in.QuotedInterest = &qi
		}
		for _, sh := range s.shares[p.ID] {
			in.Shares = append(in.Shares, s.shareInput(sh.PersonID, sh.PrincipalCents))
		}
		st.Purchases = append(st.Purchases, mc.PurchaseState{Ref: ref(p.ID), Status: p.Status, Input: in})
	}
	for _, m := range s.payments {
		if m.CreditID != c.ID || m.PurchaseID == 0 {
			continue // abonos a un corte de tarjeta: el motor los recibe por el corte (PROVISIONAL S3-formula-tarjeta)
		}
		st.Payments = append(st.Payments, mc.PaymentRecord{PurchaseRef: ref(m.PurchaseID), Seq: int(m.Seq), Date: mdate(m.PaidDate),
			Payer: s.persons[m.PayerPersonID].Code, PayerPersonID: m.PayerPersonID, Amount: mc.Cents(m.AmountCents)})
	}
	for _, t := range s.statements {
		if t.CreditID == c.ID {
			st.Statements = append(st.Statements, mc.StatementInput{StatementDate: mdate(t.StatementDate), DueDate: mdate(t.DueDate),
				Balance: mc.Cents(t.BalanceCents), Source: t.Source})
		}
	}
	return st
}

// states arma los CreditState de los créditos activos (ids vacío = todos).
func (s *snapshot) states(ids []int64) []mc.CreditState {
	want := map[int64]bool{}
	for _, id := range ids {
		want[id] = true
	}
	var out []mc.CreditState
	for _, c := range s.credits {
		if c.IsActive && (len(ids) == 0 || want[c.ID]) {
			out = append(out, s.state(c))
		}
	}
	return out
}

// horizon es el último vencimiento programado del crédito (o hoy): la ventana de la serie A/B.
func (s *snapshot) horizon(creditID int64) mc.Date {
	h := s.today.String()
	for _, sp := range s.sched {
		if sp.CreditID == creditID && sp.DueDate > h {
			h = sp.DueDate
		}
	}
	return mdate(h)
}

func (s *snapshot) timeline(c dominio.CreditRow) (mc.Timeline, error) {
	t := mdateOf(s.today)
	return mc.AvailabilityTimeline(mc.TimelineInput{State: s.state(c), Today: t, From: t, To: s.horizon(c.ID)})
}

func (s *snapshot) liveDebt(states []mc.CreditState) (mc.DebtSummary, error) {
	return mc.LiveDebt(states, mdateOf(s.today))
}

func (s *snapshot) isExternal(personID int64) bool {
	return s.persons[personID].Participation == "external"
}

// motorErr traduce un error del motor a la respuesta de C1 (422 con el campo, 503 por tiempo).
func motorErr(err error, field string) error {
	issue := func(code, msg string) error {
		return dominio.Invalid(dominio.FieldIssue{Field: field, Code: code, Message: msg})
	}
	switch {
	case errors.Is(err, mc.ErrTimeout), errors.Is(err, context.DeadlineExceeded):
		return dominio.Unavailable("El comparador tardó más del tiempo máximo. Prueba con menos días alrededor de la fecha o sube el tiempo máximo en Ajustes.")
	case errors.Is(err, mc.ErrNoRate):
		return dominio.Invalid(dominio.FieldIssue{Field: "term_count", Code: "not_found",
			Message: "Este crédito no tiene tasa capturada para ese plazo. Captúrala en Ajustes o manda los intereses cotizados."})
	case errors.Is(err, mc.ErrMissingCapture):
		return issue("required", "Faltan datos del estado de cuenta de este crédito (tasa, base, IVA o pago mínimo). Captúralos en Ajustes.")
	case errors.Is(err, mc.ErrPromotionNotApplicable):
		return dominio.Invalid(dominio.FieldIssue{Field: "promotion_id", Code: "not_allowed",
			Message: "La promoción no aplica a esta compra (comercio, compra mínima, plazo o vigencia)."})
	case errors.Is(err, mc.ErrUnsupported):
		return issue("not_allowed", "Este crédito todavía no tiene regla de intereses para este caso.")
	case errors.Is(err, mc.ErrInvalid):
		return issue("out_of_range", "Los datos no cuadran para el cálculo ("+err.Error()+"). Revisa montos, partes, plazo y fechas.")
	}
	return err
}

func nowDB(c dominio.Clock) string { return dominio.FormatDBTime(c.Now()) }
