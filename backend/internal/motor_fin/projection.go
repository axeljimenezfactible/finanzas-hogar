package motor_fin

import "sort"

// ProjectionMonth es un mes de la proyección (§3.12).
type ProjectionMonth struct {
	Index               int      `json:"index"`
	Month               Month    `json:"month"`
	Accrued             Accrued  `json:"accrued"`
	Cash                Cash     `json:"cash"`
	IncomeCents         Cents    `json:"income_cents"`
	PlannedExpenseCents Cents    `json:"planned_expense_cents"`
	NetCashCents        Cents    `json:"net_cash_cents"`
	CumulativeDiffCents Cents    `json:"cumulative_diff_cents"`
	Provisional         []string `json:"provisional"`
}

// ProjectionResult §3.12.
type ProjectionResult struct {
	RepeatCadenceEnabled bool              `json:"repeat_cadence_enabled"`
	Months               []ProjectionMonth `json:"months"`
}

// Projection proyecta 12 meses en devengado y efectivo; los meses 7–12 suman la cadencia
// solo con el interruptor encendido (K21, RN-28; cadencia PROVISIONAL, S2-proyeccion-interruptor).
func Projection(in Input, from Month) (ProjectionResult, error) {
	var res ProjectionResult
	p, err := prepared(in)
	if err != nil {
		return res, err
	}
	if from.IsZero() {
		from = p.today.MonthOf()
	}
	m7, end := from.Add(6).First(), from.Add(11).Last()
	for _, c := range p.Cadence.Purchases {
		if c.PurchaseDate.Before(m7) {
			return res, errf(CodeCadenceBefore7, "compra proyectada %d antes del mes 7", c.ID)
		}
	}
	for _, e := range p.Cadence.Expenses {
		if e.ExpenseDate.Before(m7) {
			return res, errf(CodeCadenceBefore7, "gasto proyectado %d antes del mes 7", e.ID)
		}
	}
	var xc []costItem
	var xk []cashEvent
	incOcc, planned, linkedM := map[int]Cents{}, map[int]Cents{}, map[int]bool{}
	for _, o := range p.allOccurrences(from.First(), end) {
		mi := o.date.MonthOf().idx()
		if p.linked(o) {
			linkedM[mi] = true
			continue
		}
		if o.item.Kind == "income" {
			incOcc[mi] += o.amount
			continue
		}
		sh, err := splitRule(o.amount, o.item.ShareRule)
		if err != nil {
			return res, errf(CodeInvalidArgument, "partida %d: %v", o.item.ID, err)
		}
		cat := ""
		if o.item.Category != nil {
			cat = *o.item.Category
		}
		c := p.newCost(o.date, count1(o.item.AccrualMonths), o.amount, sh, cat)
		c.paperCash = cat == p.Params.PaperCategory
		xc = append(xc, c)
		xk = append(xk, cashEvent{date: o.date, total: o.amount, shares: []Share{{p.holder, o.amount}}, category: cat, paperCash: c.paperCash})
		planned[mi] += o.amount
	}
	exps, purs := p.Expenses, p.CreditPurchases
	if p.Params.RepeatCadenceEnabled {
		exps = append(append([]Expense(nil), exps...), p.Cadence.Expenses...)
		purs = append(append([]CreditPurchase(nil), purs...), p.Cadence.Purchases...)
	}
	l := buildLedger(p, exps, purs, xc, xk)
	ms := l.monthly()
	res.RepeatCadenceEnabled = p.Params.RepeatCadenceEnabled
	for i := 0; i < 12; i++ {
		m := from.Add(i)
		s, e := MonthBounds(m)
		pm := ProjectionMonth{Index: i + 1, Month: m, Accrued: l.accruedBlock(ms[m.idx()]), Cash: l.cashBlock(s, e)}
		for _, x := range p.Incomes {
			if x.IncomeDate.in(s, e) {
				pm.IncomeCents += x.AmountCents
			}
		}
		pm.IncomeCents += incOcc[m.idx()]
		pm.PlannedExpenseCents = planned[m.idx()]
		pm.NetCashCents = pm.IncomeCents - pm.Cash.ByPerson.Get(p.holder)
		pm.CumulativeDiffCents = accruedUpTo(ms, m.idx()) - l.cashUpTo(e)
		if linkedM[m.idx()] {
			pm.Provisional = append(pm.Provisional, blockFlow)
		}
		if i >= 6 && p.Params.RepeatCadenceEnabled {
			pm.Provisional = append(pm.Provisional, blockProjection)
		}
		res.Months = append(res.Months, pm)
	}
	return res, nil
}

// ---------- Metas y fondos (§3.13) ----------

// GoalPlanMonth es una aportación planificada.
type GoalPlanMonth struct {
	Month Month `json:"month"`
	Cents Cents `json:"cents"`
}

// GoalStatus es el avance de una meta o fondo (months_left, plan y state PROVISIONAL, S2-metas).
type GoalStatus struct {
	GoalID             int64           `json:"goal_id"`
	Kind               string          `json:"kind"`
	SavedCents         Cents           `json:"saved_cents"`
	TargetCents        *Cents          `json:"target_cents"`
	ProgressTenthsPct  *int64          `json:"progress_tenths_pct"`
	ProgressPct        *int64          `json:"progress_pct"`
	MonthsLeft         *int            `json:"months_left"`
	NeededMonthlyCents *Cents          `json:"needed_monthly_cents"`
	Plan               []GoalPlanMonth `json:"plan"`
	State              string          `json:"state"`
}

// MonthSummary es el resumen del mes en curso (PROVISIONAL, S2-metas-resumen).
type MonthSummary struct {
	SavedCents                Cents `json:"saved_cents"`
	ContributedThisMonthCents Cents `json:"contributed_this_month_cents"`
	PlannedThisMonthCents     Cents `json:"planned_this_month_cents"`
}

// SavingsResult §3.13.
type SavingsResult struct {
	MonthSummary MonthSummary `json:"month_summary"`
	Goals        []GoalStatus `json:"goals"`
}

// SavingsGoals calcula ahorrado, avance y plan de aportaciones (RN-18, S2-12).
func SavingsGoals(in Input) (SavingsResult, error) {
	res := SavingsResult{Goals: []GoalStatus{}}
	p, err := prepared(in)
	if err != nil {
		return res, err
	}
	T := p.today
	goals := append([]SavingsGoal(nil), p.SavingsGoals...)
	sort.SliceStable(goals, func(a, b int) bool { return goals[a].ID < goals[b].ID })
	for _, g := range goals {
		if !active(g.IsActive) {
			continue
		}
		st := GoalStatus{GoalID: g.ID, Kind: g.Kind, TargetCents: g.TargetCents, Plan: []GoalPlanMonth{}}
		for _, mv := range p.SavingsMovements {
			if mv.GoalID != g.ID || mv.MovementDate.After(T) {
				continue
			}
			if mv.Kind == "deposit" {
				st.SavedCents += mv.AmountCents
				if mv.MovementDate.MonthOf() == T.MonthOf() {
					res.MonthSummary.ContributedThisMonthCents += mv.AmountCents
				}
			} else {
				st.SavedCents -= mv.AmountCents
			}
		}
		st.State = "open"
		if g.TargetDate != nil {
			ml := 0
			if !g.TargetDate.Before(T) {
				ml = g.TargetDate.MonthOf().idx() - T.MonthOf().idx() + 1
			}
			st.MonthsLeft = &ml
		}
		if t := g.TargetCents; t != nil && *t > 0 {
			tp := halfUp(int64(st.SavedCents)*1000, int64(*t))
			pp := halfUp(int64(st.SavedCents)*100, int64(*t))
			st.ProgressTenthsPct, st.ProgressPct = &tp, &pp
		}
		if g.TargetCents != nil && g.TargetDate != nil {
			rem := *g.TargetCents - st.SavedCents
			var need Cents
			switch {
			case rem <= 0:
				st.State = "achieved"
			default:
				n := *st.MonthsLeft
				st.State = "in_progress"
				if n < 1 {
					n, st.State = 1, "overdue"
				}
				for k := 0; k < n; k++ {
					st.Plan = append(st.Plan, GoalPlanMonth{T.MonthOf().Add(k), series(rem, n, k)})
				}
				need = st.Plan[0].Cents
			}
			st.NeededMonthlyCents = &need
		}
		res.MonthSummary.SavedCents += st.SavedCents
		if g.PlannedMonthlyCents != nil {
			res.MonthSummary.PlannedThisMonthCents += *g.PlannedMonthlyCents
		}
		res.Goals = append(res.Goals, st)
	}
	return res, nil
}
