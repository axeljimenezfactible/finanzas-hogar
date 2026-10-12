package gastos

import (
	"context"
	"fmt"
	"net/http"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// Tarjeta «Próximo pago y disponible» (U3, N-022, M-03 y M-04). El próximo pago sale de los pagos programados
// guardados y sus partes; el disponible A/B, del motor S3 (motor_cred.AvailabilityTimeline).

// credInput es lo de créditos leído de la base para armar la entrada del motor S3. Es el mismo armado que
// internal/api/creditos/motor.go (snapshot.state), que no se exporta: PROPUESTA del bloqueo
// S6-gastos-inicio §13 (exportarlo para que haya un solo armado).
type credInput struct {
	today      dominio.Date
	credits    []dominio.CreditRow
	persons    map[int64]dominio.PersonRow
	rates      map[int64][]dominio.TermRateRow
	promos     []dominio.PromotionRow
	purchases  []dominio.PurchaseRow
	shares     map[int64][]dominio.PurchaseShareRow
	sched      []dominio.ScheduledRow
	payments   []dominio.PaymentRow
	statements []dominio.StatementRow
}

func loadCredInput(ctx context.Context, q dominio.Querier, today dominio.Date) (*credInput, error) {
	s := &credInput{today: today, rates: map[int64][]dominio.TermRateRow{}}
	var err error
	if s.credits, err = dominio.QueryCredits(ctx, q, "", "ORDER BY c.id"); err != nil {
		return nil, err
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
	s.statements, err = dominio.QueryStatements(ctx, q)
	return s, err
}

func cdate(s string) mc.Date {
	d, _ := mc.ParseDate(s)
	return d
}

func cdateOf(d dominio.Date) mc.Date { return mc.Date{Year: d.Year, Month: d.Month, Day: d.Day} }

func cdatePtr(s *string) *mc.Date {
	if s == nil {
		return nil
	}
	d := cdate(*s)
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

// credit arma el crédito para el motor S3 (tasas por plazo y promociones incluidas); lo usan state y confirmTicketScan.
func (s *credInput) credit(c dominio.CreditRow) mc.Credit {
	cr := mc.Credit{ID: c.ID, Name: c.Name, Kind: mc.Kind(c.Kind), Holder: s.persons[c.HolderPersonID].Code,
		HolderPersonID: c.HolderPersonID, Limit: centsPtr(c.CreditLimitCents), Frequency: c.PaymentFrequency,
		FirstPaymentMinDays: intPtr(c.FirstPaymentMinDays), Reading: mc.Reading(c.AvailabilityReading),
		ReadingConfirmed: c.AvailabilityConfirmed, ReadingAAnchor: cdatePtr(c.ReadingAAnchorDate),
		StatementDay: intPtr(c.StatementDay), PaymentDueDay: intPtr(c.PaymentDueDay),
		MinPaymentPctBP: intPtr(c.MinPaymentPctBp), MinPaymentFixed: centsPtr(c.MinPaymentFixedCents)}
	if c.AnnualRateNum != nil && c.AnnualRateDen != nil {
		cr.AnnualRate = &mc.Rate{Num: *c.AnnualRateNum, Den: *c.AnnualRateDen, Decimals: decimals(c.AnnualRateDecimals)}
	}
	if c.RateBase != nil {
		cr.RateBase = *c.RateBase
	}
	if c.VatIncluded == nil || !*c.VatIncluded {
		cr.VATOnInterestBP = intPtr(c.VatOnInterestBp)
	}
	for _, r := range s.rates[c.ID] {
		cr.TermRates = append(cr.TermRates, mc.TermRate{TermCount: int(r.TermCount),
			Rate: mc.Rate{Num: r.RateNum, Den: r.RateDen, Decimals: decimals(r.RateDecimals)}, Source: r.Source, QuoteDate: cdatePtr(r.QuoteDate)})
	}
	for _, p := range s.promos {
		if p.CreditID != nil && *p.CreditID != c.ID {
			continue
		}
		cr.Promotions = append(cr.Promotions, mc.Promotion{ID: p.ID, CreditID: p.CreditID, Merchant: p.MerchantName,
			StoreID: p.StoreID, TermCount: int(p.TermCount), Rate: mc.Rate{Num: p.RateNum, Den: p.RateDen, Decimals: decimals(p.RateDecimals)},
			MinPurchase: mc.Cents(p.MinPurchaseCents), ValidFrom: cdate(p.ValidFrom), ValidTo: cdate(p.ValidTo), Active: p.IsActive})
	}
	return cr
}

func (s *credInput) state(c dominio.CreditRow) mc.CreditState {
	cr := s.credit(c)
	st := mc.CreditState{Credit: cr}
	ref := func(id int64) string { return fmt.Sprintf("p%09d", id) }
	for _, p := range s.purchases {
		if p.CreditID != c.ID {
			continue
		}
		fp := cdate(p.FirstPaymentDate)
		in := mc.PurchaseInput{Credit: cr, Date: cdate(p.PurchaseDate), Principal: mc.Cents(p.PrincipalCents),
			TermCount: int(p.TermCount), RateDecimals: intPtr(p.RateDecimals), FirstPaymentDate: &fp}
		if p.InterestSource != "amortization" {
			qi := mc.Cents(p.InterestCents)
			in.QuotedInterest = &qi
		}
		for _, sh := range s.shares[p.ID] {
			pr := s.persons[sh.PersonID]
			in.Shares = append(in.Shares, mc.ShareInput{PersonID: sh.PersonID, Code: pr.Code, External: pr.Participation == "external", Principal: mc.Cents(sh.PrincipalCents)})
		}
		st.Purchases = append(st.Purchases, mc.PurchaseState{Ref: ref(p.ID), Status: p.Status, Input: in})
	}
	for _, m := range s.payments {
		if m.CreditID != c.ID || m.PurchaseID == 0 {
			continue
		}
		st.Payments = append(st.Payments, mc.PaymentRecord{PurchaseRef: ref(m.PurchaseID), Seq: int(m.Seq), Date: cdate(m.PaidDate),
			Payer: s.persons[m.PayerPersonID].Code, PayerPersonID: m.PayerPersonID, Amount: mc.Cents(m.AmountCents)})
	}
	for _, t := range s.statements {
		if t.CreditID == c.ID {
			st.Statements = append(st.Statements, mc.StatementInput{StatementDate: cdate(t.StatementDate), DueDate: cdate(t.DueDate),
				Balance: mc.Cents(t.BalanceCents), Source: t.Source})
		}
	}
	return st
}

// availability: una fila por línea revolvente activa con el disponible de hoy en su lectura activa (motor S3).
func (s *credInput) availability(reg *api.Registry, ctx context.Context) []api.DashboardNextPaymentAvailabilityItem {
	out := []api.DashboardNextPaymentAvailabilityItem{}
	for _, c := range s.credits {
		if !c.IsActive || c.Kind != "revolving_line" {
			continue
		}
		h := s.today.String()
		for _, sp := range s.sched {
			if sp.CreditID == c.ID && sp.DueDate > h {
				h = sp.DueDate
			}
		}
		t := cdateOf(s.today)
		tl, err := mc.AvailabilityTimeline(mc.TimelineInput{State: s.state(c), Today: t, From: t, To: cdate(h)})
		if err != nil || len(tl.Events) == 0 {
			// Sin dato del motor la línea no se muestra (PROVISIONAL, bloqueo S6-gastos-inicio §13).
			reg.Deps().Log.WarnContext(ctx, "disponible sin dato del motor", "credit_id", c.ID, "err", err)
			continue
		}
		v := tl.Events[0].AvailableB
		if tl.Reading == "A" {
			v = tl.Events[0].AvailableA
		}
		out = append(out, api.DashboardNextPaymentAvailabilityItem{CreditID: c.ID, AvailabilityReading: string(tl.Reading),
			AvailableCents: int64(v), AvailabilityConfirmed: c.AvailabilityConfirmed})
	}
	return out
}

// payShares parte un pago programado en hogar (titulares) y externa (Karen) y dice si queda algo pendiente.
// Sin partes guardadas, todo es del hogar y cuenta como pendiente (PROVISIONAL, §13).
func payShares(sp dominio.ScheduledRow, shares []dominio.ScheduledShareRow, persons map[int64]dominio.PersonRow) (hh, ext int64, pending bool) {
	if len(shares) == 0 {
		return sp.AmountCents, 0, true
	}
	for _, sh := range shares {
		if persons[sh.PersonID].Participation == "external" {
			ext += sh.AmountCents
		} else {
			hh += sh.AmountCents
		}
		pending = pending || sh.PaidCents < sh.AmountCents
	}
	return hh, ext, pending
}

func getDashboardNextPayment(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx, db := r.Context(), reg.Deps().DB
		t := today(reg)
		s, err := loadCredInput(ctx, db, t)
		if err != nil {
			return err
		}
		spShares, err := dominio.ScheduledShares(ctx, db, 0)
		if err != nil {
			return err
		}
		cancelled := map[int64]bool{}
		for _, p := range s.purchases {
			cancelled[p.ID] = p.Status == "cancelled"
		}
		// N-022 a: el pendiente vencido más antiguo; si no hay, el pendiente más cercano ≥ hoy. s.sched va por fecha.
		var pend []dominio.ScheduledRow
		for _, sp := range s.sched {
			if _, _, p := payShares(sp, spShares[sp.ID], s.persons); p && !cancelled[sp.CreditPurchaseID] {
				pend = append(pend, sp)
			}
		}
		out := api.DashboardNextPayment{Availability: s.availability(reg, ctx)}
		if len(pend) > 0 {
			day := pend[0].DueDate // vencido más antiguo o, sin vencidos, el más cercano ≥ hoy
			np := &api.DashboardNextPaymentNextPayment{CreditID: pend[0].CreditID, IsOverdue: day < t.String()}
			for _, sp := range pend {
				if sp.DueDate != day {
					continue
				}
				hh, ext, _ := payShares(sp, spShares[sp.ID], s.persons)
				np.AmountCents += sp.AmountCents
				np.HouseholdShareCents += hh
				np.ExternalShareCents += ext
				np.PaymentCount++
			}
			d, err := dominio.ParseDate(day)
			if err != nil {
				return err
			}
			np.DueDate, np.DaysLeft = d.String(), int64(t.DaysUntil(d))
			out.NextPayment = np
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}
