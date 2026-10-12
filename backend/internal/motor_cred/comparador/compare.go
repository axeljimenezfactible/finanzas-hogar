package comparador

import (
	"context"
	"sort"

	mc "finanzas-hogar/internal/motor_cred"
)

// Compare corre el comparador (S3-interfaz.md §8). Revisa el contexto antes de cada candidato y
// al terminar; si venció devuelve motor_cred.ErrTimeout sin resultados parciales (§8.7).
func Compare(ctx context.Context, req Request) (Result, error) {
	if err := req.validate(); err != nil {
		return Result{}, err
	}
	used, missing := merge(req.Global, req.Override)
	res := Result{MissingConstraints: missing, ConstraintsUsed: used, Valid: []Scenario{}, Rejected: []Rejected{},
		Excluded: []ExcludedCredit{}, Provisional: []string{}}
	if len(missing) > 0 { // K16: «faltan restricciones», sin escenarios ni recomendación
		res.Status = StatusMissing
		if err := checkCtx(ctx); err != nil {
			return Result{}, err
		}
		return res, nil
	}
	credits := append([]mc.CreditState(nil), req.Credits...)
	sort.SliceStable(credits, func(i, j int) bool { return credits[i].Credit.ID < credits[j].Credit.ID })
	for i := 1; i < len(credits); i++ {
		if credits[i].Credit.ID == credits[i-1].Credit.ID {
			return Result{}, invalid("crédito %d repetido", credits[i].Credit.ID)
		}
	}
	cands, excluded, err := generate(ctx, req, *used.DateWindowDays, credits)
	if err != nil {
		return Result{}, err
	}
	res.Excluded = excluded
	if len(cands) > 0 {
		from, to := fqOf(cands[0].in.Date).first(), cands[0].quote.Payments[0].Due
		for _, c := range cands {
			if c.in.Date.Compare(from) < 0 {
				from = fqOf(c.in.Date).first()
			}
			if d := c.quote.Payments[len(c.quote.Payments)-1].Due; d.Compare(to) > 0 {
				to = d
			}
		}
		to = fqOf(to).last()
		ctxs := make([]creditCtx, len(credits))
		for i, cs := range credits {
			if ctxs[i], err = prepare(cs, from, to); err != nil {
				return Result{}, err
			}
		}
		for _, c := range cands {
			if err := checkCtx(ctx); err != nil {
				return Result{}, err
			}
			s, err := evaluate(req, credits, ctxs, c)
			if err != nil {
				return Result{}, err
			}
			if v := violation(s, used); v != "" {
				res.Rejected = append(res.Rejected, Rejected{s, v})
			} else {
				res.Valid = append(res.Valid, s)
			}
		}
	}
	sort.SliceStable(res.Valid, func(i, j int) bool { return keyCmp(&res.Valid[i], &res.Valid[j]) < 0 })
	sort.SliceStable(res.Rejected, func(i, j int) bool { return keyCmp(&res.Rejected[i].Scenario, &res.Rejected[j].Scenario) < 0 })
	res.EvaluatedCount = len(res.Valid) + len(res.Rejected)
	provs := [][]string{{provObjectives}}
	for _, s := range res.Valid {
		provs = append(provs, s.Provisional)
	}
	for _, r := range res.Rejected {
		provs = append(provs, r.Scenario.Provisional)
	}
	res.Provisional = provList(provs...)
	res.Status = StatusNoValid
	if len(res.Valid) > 0 {
		res.Status = StatusOK
		exact := scores(res.Valid, req.Weights)
		v := res.Valid
		res.Optima = &Optima{v[best(v, exact, 0)], v[best(v, exact, 1)], v[best(v, exact, 2)], v[best(v, exact, 3)]}
		rec := v[recommend(v, exact)]
		res.Recommendation = &rec
	}
	if err := checkCtx(ctx); err != nil {
		return Result{}, err // sin resultados parciales
	}
	return res, nil
}

// evaluate calcula los valores de un candidato con el núcleo (§8.3, §8.4).
func evaluate(req Request, credits []mc.CreditState, ctxs []creditCtx, c candidate) (Scenario, error) {
	st := withScenario(credits[c.ci], c.in)
	pays := c.quote.Payments
	p, l := c.in.Date, pays[len(pays)-1].Due
	tl, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, Today: req.Today, From: req.Today, To: l})
	if err != nil {
		return Scenario{}, err
	}
	cc := ctxs[c.ci]
	zero := l
	if !cc.lastDue.IsZero() && cc.lastDue.Compare(zero) > 0 {
		zero = cc.lastDue
	}
	from, to := fqOf(p).first(), fqOf(l).last()
	cal, err := mc.UnifiedCalendar([]mc.CreditState{st}, from, to)
	if err != nil {
		return Scenario{}, err
	}
	entries := append([]mc.CalendarEntry(nil), cal.Entries...)
	for i, o := range ctxs {
		if i != c.ci {
			entries = append(entries, o.entries...)
		}
	}
	tot, byPerson, eprov := series(entries, p, l)
	hi, lo := maxMin(tot)
	return Scenario{
		CreditID: credits[c.ci].Credit.ID, TermCount: c.in.TermCount, PurchaseDate: p, StoreID: c.store.id,
		PromotionID: c.in.PromotionID, Interest: c.quote.Interest, MaxFortnightPayment: hi,
		FortnightPaymentSpread: hi - lo, StdDevMilliCents: stdDevMilli(tot), MinAvailable: minAvailable(tl, p, l),
		MaxLivePurchases: liveMax(append(append([]span(nil), cc.spans...), span{p, l}), p, l), ZeroDebtDate: zero,
		ValuePerPesoMilli: valuePerPeso(c.in.Principal, c.quote.Interest, c.store.basket),
		OverCap:           overCap(req.Caps, byPerson),
		Provisional:       provList([]string{provObjectives}, c.quote.Provisional, eprov),
	}, nil
}
