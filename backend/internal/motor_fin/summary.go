package motor_fin

import "sort"

// PurchaseInHorizon §3.3.
type PurchaseInHorizon struct {
	PurchaseID     int64 `json:"purchase_id"`
	PurchaseSeq    *int  `json:"purchase_seq"`
	InHorizonCents Cents `json:"in_horizon_cents"`
	Complete       bool  `json:"complete"`
}

// HorizonSpec §3.3.
type HorizonSpec struct {
	Start, End Date
	Months     []Month
}

// HorizonSummary §3.3.
type HorizonSummary struct {
	PaymentsByPurchase     []PurchaseInHorizon `json:"payments_by_purchase"`
	PaperCashCents         Cents               `json:"paper_cash_cents"`
	CashTotalCents         Cents               `json:"cash_total_cents"`
	SumOfPeriodsCashCents  Cents               `json:"sum_of_periods_cash_cents"`
	AccruedTotalCents      Cents               `json:"accrued_total_cents"`
	LiveDebt               LiveDebtResult      `json:"live_debt"`
	MonthlyAvgCashCents    Cents               `json:"monthly_avg_cash_cents"`
	MonthlyAvgAccruedCents Cents               `json:"monthly_avg_accrued_cents"`
	FinalDiffCents         Cents               `json:"final_diff_cents"`
}

// PendingPurchase es una compra con pagos después del corte (§3.4).
type PendingPurchase struct {
	PurchaseID   int64 `json:"purchase_id"`
	PurchaseSeq  *int  `json:"purchase_seq"`
	PendingCents Cents `json:"pending_cents"`
	Started      bool  `json:"started"`
}

// LiveDebtResult §3.4 (K19, RN-21).
type LiveDebtResult struct {
	Cut                        Date              `json:"cut"`
	LiveDebtCents              Cents             `json:"live_debt_cents"`
	ByPurchase                 []PendingPurchase `json:"by_purchase"`
	AccruedToDateCents         Cents             `json:"accrued_to_date_cents"`
	CashToDateCents            Cents             `json:"cash_to_date_cents"`
	CumulativeDiffCents        Cents             `json:"cumulative_diff_cents"`
	UnaccruedBalanceCents      Cents             `json:"unaccrued_balance_cents"`
	IdentityHolds              bool              `json:"identity_holds"`
	IdentityWithoutBalanceHold bool              `json:"identity_without_balance_holds"`
}

func (l *ledger) sortedPurchases() []*CreditPurchase {
	ps := append([]*CreditPurchase(nil), l.purchs...)
	sort.SliceStable(ps, func(a, b int) bool {
		if ps[a].PurchaseDate != ps[b].PurchaseDate {
			return ps[a].PurchaseDate.Before(ps[b].PurchaseDate)
		}
		return ps[a].ID < ps[b].ID
	})
	return ps
}

func (l *ledger) liveDebt(cut Date) LiveDebtResult {
	r := LiveDebtResult{Cut: cut, ByPurchase: []PendingPurchase{}}
	for _, p := range l.sortedPurchases() {
		if p.PurchaseDate.After(cut) {
			continue
		}
		var pend Cents
		started := false
		for _, pay := range p.Payments {
			if pay.DueDate.After(cut) {
				pend += pay.AmountCents
			} else {
				started = true
			}
		}
		if pend > 0 {
			r.ByPurchase = append(r.ByPurchase, PendingPurchase{p.ID, seqPtr(l.seqOf[p.ID]), pend, started})
			r.LiveDebtCents += pend
		}
	}
	ms := l.monthly()
	for i, a := range ms {
		if !monthFromIdx(i).Last().After(cut) {
			r.AccruedToDateCents += a.total
		}
	}
	r.CashToDateCents = l.cashUpTo(cut)
	r.CumulativeDiffCents = r.AccruedToDateCents - r.CashToDateCents
	var cost Cents
	for _, c := range l.costs {
		if !c.date.After(cut) {
			cost += c.total
		}
	}
	r.UnaccruedBalanceCents = cost - r.AccruedToDateCents
	r.IdentityHolds = r.LiveDebtCents == r.CumulativeDiffCents+r.UnaccruedBalanceCents
	r.IdentityWithoutBalanceHold = r.LiveDebtCents == r.CumulativeDiffCents
	return r
}

// LiveDebt calcula la deuda viva al corte y la identidad de K19 (§3.4).
func LiveDebt(in Input, cut Date) (LiveDebtResult, error) {
	p, err := prepared(in)
	if err != nil {
		return LiveDebtResult{}, err
	}
	return records(p).liveDebt(cut), nil
}

// Horizon resume un horizonte de fechas y meses (§3.3).
func Horizon(in Input, spec HorizonSpec) (HorizonSummary, error) {
	var h HorizonSummary
	p, err := prepared(in)
	if err != nil {
		return h, err
	}
	if spec.End.Before(spec.Start) || len(spec.Months) == 0 {
		return h, errf(CodeInvalidArgument, "horizonte: start ≤ end y months no vacío")
	}
	l := records(p)
	h.PaymentsByPurchase = []PurchaseInHorizon{}
	for _, cp := range l.sortedPurchases() {
		var in Cents
		n, all := 0, 0
		for _, pay := range cp.Payments {
			all++
			if pay.DueDate.in(spec.Start, spec.End) {
				in += pay.AmountCents
				n++
			}
		}
		if n > 0 {
			h.PaymentsByPurchase = append(h.PaymentsByPurchase, PurchaseInHorizon{cp.ID, seqPtr(l.seqOf[cp.ID]), in, n == all})
		}
	}
	for _, e := range l.cash {
		if !e.payment && e.paperCash && e.date.in(spec.Start, spec.End) {
			h.PaperCashCents += e.total
		}
	}
	h.CashTotalCents = l.cashBetween(spec.Start, spec.End)
	ps, err := l.periods(PeriodSpec{Months: spec.Months})
	if err != nil {
		return h, err
	}
	for _, x := range ps {
		h.SumOfPeriodsCashCents += x.Cash.TotalCents
		h.AccruedTotalCents += x.Accrued.TotalCents
	}
	h.LiveDebt = l.liveDebt(spec.End)
	k := int64(len(spec.Months))
	h.MonthlyAvgCashCents = Cents(halfUp(int64(h.CashTotalCents), k))
	h.MonthlyAvgAccruedCents = Cents(halfUp(int64(h.AccruedTotalCents), k))
	h.FinalDiffCents = h.AccruedTotalCents - h.CashTotalCents
	return h, nil
}

// Average §3.5.
type Average struct {
	Days        int   `json:"days"`
	DailyCents  Cents `json:"daily_cents"`
	WeeklyCents Cents `json:"weekly_cents"`
}

// Averages: diario = HALF_UP(T/d); semanal = HALF_UP(7·T/d), un solo redondeo (K20, K39).
func Averages(total Cents, start, end Date) (Average, error) {
	if end.Before(start) {
		return Average{}, errf(CodeInvalidArgument, "promedios: start > end")
	}
	d := int64(daysIncl(start, end))
	return Average{int(d), Cents(halfUp(int64(total), d)), Cents(halfUp(7*int64(total), d))}, nil
}

// MonthlyAverage = HALF_UP(T/k) (R-56, R-57).
func MonthlyAverage(total Cents, months int) (Cents, error) {
	if months < 1 {
		return 0, errf(CodeInvalidArgument, "promedio_mensual: month_count ≥ 1")
	}
	return Cents(halfUp(int64(total), int64(months))), nil
}

// StableComponent §3.6.
type StableComponent struct {
	Label        string   `json:"label"`
	MonthlyCents Cents    `json:"monthly_cents"`
	ByPerson     CentsMap `json:"by_person"`
}

// StableBlock es el mes del valor estable.
type StableBlock struct {
	TotalCents Cents    `json:"total_cents"`
	ByPerson   CentsMap `json:"by_person"`
}

// StableFortnight es la quincena del valor estable.
type StableFortnight struct {
	TotalCents      Cents    `json:"total_cents"`
	ByPerson        CentsMap `json:"by_person"`
	SumPersonsCents Cents    `json:"sum_persons_cents"`
	CentDiffCents   Cents    `json:"cent_diff_cents"`
}

// StableValueResult §3.6.
type StableValueResult struct {
	Month      StableBlock       `json:"month"`
	Fortnight  StableFortnight   `json:"fortnight"`
	Components []StableComponent `json:"components"`
}

// StableValue: cada partida aporta HALF_UP(total/ciclo) al mes, con personas como en el
// devengo; la quincena es HALF_UP(x/2) por renglón y la diferencia de centavo se informa.
func StableValue(items []StableItem, persons []Person) (StableValueResult, error) {
	var r StableValueResult
	var keys []string
	holder := ""
	for _, p := range persons {
		switch p.Participation {
		case "holder":
			holder = p.Code
			keys = append(keys, p.Code)
		case "external":
			keys = append(keys, p.Code)
		}
	}
	if holder == "" {
		return r, errf(CodeInvalidArgument, "valor_estable: falta el titular")
	}
	r.Month.ByPerson = newCentsMap(keys)
	r.Components = []StableComponent{}
	for _, it := range items {
		if it.CycleMonths < 1 {
			return r, errf(CodeInvalidArgument, "valor_estable: cycle_months < 1")
		}
		var sh []Share
		for _, s := range it.Shares {
			sh = append(sh, Share{s.Person, s.TotalCents})
		}
		c := (&Input{}).newCost(Date{}, it.CycleMonths, it.TotalCents, sh, it.Label)
		tot, pp := c.slot(0, holder)
		bp := newCentsMap(keys)
		for p, v := range pp {
			if _, ok := bp.Vals[p]; !ok {
				return r, errf(CodeInvalidArgument, "valor_estable: persona %q", p)
			}
			bp.Vals[p] += v
			r.Month.ByPerson.Vals[p] += v
		}
		r.Month.TotalCents += tot
		r.Components = append(r.Components, StableComponent{it.Label, tot, bp})
	}
	r.Fortnight.TotalCents = Cents(halfUp(int64(r.Month.TotalCents), 2))
	r.Fortnight.ByPerson = newCentsMap(keys)
	for _, k := range keys {
		v := Cents(halfUp(int64(r.Month.ByPerson.Vals[k]), 2))
		r.Fortnight.ByPerson.Vals[k] = v
		r.Fortnight.SumPersonsCents += v
	}
	r.Fortnight.CentDiffCents = r.Fortnight.TotalCents - r.Fortnight.SumPersonsCents
	return r, nil
}

// PaperPackage §3.7.
type PaperPackage struct {
	Package     int    `json:"package"`
	Date        Date   `json:"date"`
	Person      string `json:"person"`
	PersonName  string `json:"person_name"`
	Form        string `json:"form"`
	PurchaseSeq *int   `json:"purchase_seq"`
	AmountCents Cents  `json:"amount_cents"`
}

// PaperPlan lista los paquetes de papel: dentro de despensa y de contado (§3.7).
func PaperPlan(in Input) ([]PaperPackage, error) {
	p, err := prepared(in)
	if err != nil {
		return nil, err
	}
	l := records(p)
	type ent struct {
		pk      PaperPackage
		rank    int
		id, pos int64
	}
	var es []ent
	for _, cp := range l.purchs {
		for i, e := range cp.Embedded {
			if e.Category == p.Params.PaperCategory {
				es = append(es, ent{PaperPackage{Date: cp.PurchaseDate, Person: e.Person, Form: "in_groceries",
					PurchaseSeq: seqPtr(l.seqOf[cp.ID]), AmountCents: e.AmountCents}, 0, cp.ID, int64(i)})
			}
		}
	}
	for _, e := range p.Expenses {
		if e.PaymentMethod != "credit" && e.Category == p.Params.PaperCategory {
			who := p.holder
			if e.PaidBy != nil {
				who = *e.PaidBy
			}
			es = append(es, ent{PaperPackage{Date: e.ExpenseDate, Person: who, Form: "cash", AmountCents: e.AmountCents}, 1, e.ID, 0})
		}
	}
	sort.SliceStable(es, func(a, b int) bool {
		x, y := es[a], es[b]
		if x.pk.Date != y.pk.Date {
			return x.pk.Date.Before(y.pk.Date)
		}
		if x.rank != y.rank {
			return x.rank < y.rank
		}
		if x.id != y.id {
			return x.id < y.id
		}
		return x.pos < y.pos
	})
	out := []PaperPackage{}
	for i, e := range es {
		e.pk.Package = i + 1
		e.pk.PersonName = p.byCode[e.pk.Person].DisplayName
		out = append(out, e.pk)
	}
	return out, nil
}
