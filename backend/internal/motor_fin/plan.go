package motor_fin

import "sort"

// ---------- Presupuesto (§3.8–§3.10) ----------

// BudgetLine es una fila de categoría o el total del presupuesto del mes.
type BudgetLine struct {
	Category       string `json:"category,omitempty"`
	PlanCents      *Cents `json:"plan_cents"`
	ActualCents    Cents  `json:"actual_cents"`
	RemainingCents *Cents `json:"remaining_cents"`
	UsedPct        *int64 `json:"used_pct"`
	Warning        bool   `json:"warning"`
	Limit          bool   `json:"limit"`
	State          string `json:"state"`
}

// BudgetMonthResult §3.8.
type BudgetMonthResult struct {
	Month                 Month        `json:"month"`
	DaysInMonth           int          `json:"days_in_month"`
	DaysLeft              int          `json:"days_left"`
	Categories            []BudgetLine `json:"categories"`
	Total                 BudgetTotal  `json:"total"`
	AccruedHouseholdCents Cents        `json:"accrued_household_cents"`
	CashHouseholdCents    Cents        `json:"cash_household_cents"`
}

// BudgetTotal es el renglón total (declared_total_row_cents PROVISIONAL, S2-presupuesto-total).
type BudgetTotal struct {
	PlanCents             Cents  `json:"plan_cents"`
	ActualCents           Cents  `json:"actual_cents"`
	RemainingCents        Cents  `json:"remaining_cents"`
	UsedPct               *int64 `json:"used_pct"`
	Warning               bool   `json:"warning"`
	Limit                 bool   `json:"limit"`
	State                 string `json:"state"`
	DeclaredTotalRowCents *Cents `json:"declared_total_row_cents"`
}

func (in *Input) flags(plan, actual Cents) (*int64, bool, bool, string) {
	if plan == 0 { // plan 0: PROVISIONAL (S2-presupuesto-total)
		return nil, false, false, "no_plan"
	}
	u := halfUp(int64(actual)*100, int64(plan))
	w := int64(actual)*100 >= in.warnPct*int64(plan)
	l := int64(actual)*100 >= in.limitPct*int64(plan)
	st := "normal"
	if l {
		st = "exceeded"
	} else if w {
		st = "warning"
	}
	return &u, w, l, st
}

func rowFor(rows []BudgetRow, cat *string, m Month) *BudgetRow {
	for i := range rows {
		r := &rows[i]
		if (r.Category == nil) != (cat == nil) || (cat != nil && *r.Category != *cat) {
			continue
		}
		if !r.ValidFrom.After(m.First()) && (r.ValidTo == nil || !r.ValidTo.Before(m.Last())) {
			return r
		}
	}
	return nil
}

func (in *Input) budgetMonth(ms map[int]*agg, l *ledger, m Month) BudgetMonthResult {
	r := BudgetMonthResult{Month: m, DaysInMonth: m.Days(), Categories: []BudgetLine{}}
	switch {
	case in.today.After(m.Last()):
	case in.today.Before(m.First()):
		r.DaysLeft = m.Days()
	default:
		r.DaysLeft = daysIncl(in.today, m.Last())
	}
	a := ms[m.idx()]
	if a == nil {
		a = newAgg()
	}
	cats := append([]Category(nil), in.Categories...)
	sort.SliceStable(cats, func(x, y int) bool { return cats[x].Code < cats[y].Code })
	for _, c := range cats {
		if !active(c.IsActive) {
			continue
		}
		code := c.Code
		line := BudgetLine{Category: code, ActualCents: a.catPers[code][in.holder], State: "no_plan"}
		if row := rowFor(in.Budgets, &code, m); row != nil {
			pl := row.AmountCents
			rem := pl - line.ActualCents
			line.PlanCents, line.RemainingCents = &pl, &rem
			line.UsedPct, line.Warning, line.Limit, line.State = in.flags(pl, line.ActualCents)
			r.Total.PlanCents += pl
		}
		r.Total.ActualCents += line.ActualCents
		r.Categories = append(r.Categories, line)
	}
	t := &r.Total
	t.RemainingCents = t.PlanCents - t.ActualCents
	t.UsedPct, t.Warning, t.Limit, t.State = in.flags(t.PlanCents, t.ActualCents)
	if row := rowFor(in.Budgets, nil, m); row != nil {
		v := row.AmountCents
		t.DeclaredTotalRowCents = &v
	}
	r.AccruedHouseholdCents = a.byPerson[in.holder]
	r.CashHouseholdCents = l.cashBlock(m.First(), m.Last()).ByPerson.Get(in.holder)
	return r
}

// BudgetMonth: plan, real (devengado de la parte del hogar) y banderas 80/100 % (§3.8).
func BudgetMonth(in Input, m Month) (BudgetMonthResult, error) {
	p, err := prepared(in)
	if err != nil {
		return BudgetMonthResult{}, err
	}
	l := records(p)
	return p.budgetMonth(l.monthly(), l, m), nil
}

// BudgetChangeResult §3.9.
type BudgetChangeResult struct {
	Rows           []BudgetRow `json:"rows"`
	EffectiveRowID int64       `json:"effective_row_id"`
	DeletedRowIDs  []int64     `json:"deleted_row_ids"`
}

// ApplyBudgetChange aplica un cambio con alcance onward o this_month (K2, RN-22).
func ApplyBudgetChange(rows []BudgetRow, ch BudgetChange) (BudgetChangeResult, error) {
	var res BudgetChangeResult
	rs := append([]BudgetRow(nil), rows...)
	idx := -1
	var maxID int64
	for i, r := range rs {
		if r.ID == ch.RowID {
			idx = i
		}
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	if idx < 0 {
		return res, errf(CodeInvalidArgument, "row_id %d no existe", ch.RowID)
	}
	scope := ch.Scope
	if scope == "" {
		scope = "onward"
	}
	if scope != "onward" && scope != "this_month" {
		return res, errf(CodeInvalidArgument, "scope %q", scope)
	}
	r := rs[idx]
	mf, ml := ch.Month.First(), ch.Month.Last()
	if r.ValidFrom.After(mf) || (r.ValidTo != nil && r.ValidTo.Before(ml)) {
		return res, errf(CodeOutOfRange, "el mes %s está fuera de la vigencia del renglón %d", ch.Month, r.ID)
	}
	newRow := func(from Date, to *Date, amt Cents) int64 {
		maxID++
		rs = append(rs, BudgetRow{ID: maxID, Category: r.Category, AmountCents: amt, ValidFrom: from, ValidTo: to})
		return maxID
	}
	prev := ch.Month.Add(-1).Last()
	origTo := r.ValidTo
	starts := r.ValidFrom == mf
	if scope == "onward" {
		if starts {
			rs[idx].AmountCents, rs[idx].ValidTo = ch.AmountCents, nil
			res.EffectiveRowID = r.ID
		} else {
			rs[idx].ValidTo = &prev
			res.EffectiveRowID = newRow(mf, nil, ch.AmountCents)
		}
		var kept []BudgetRow
		for _, x := range rs {
			same := (x.Category == nil) == (r.Category == nil) && (x.Category == nil || *x.Category == *r.Category)
			if same && x.ValidFrom.After(mf) {
				res.DeletedRowIDs = append(res.DeletedRowIDs, x.ID)
				continue
			}
			kept = append(kept, x)
		}
		rs = kept
	} else {
		end := ml
		if starts {
			rs[idx].AmountCents, rs[idx].ValidTo = ch.AmountCents, &end
			res.EffectiveRowID = r.ID
		} else {
			rs[idx].ValidTo = &prev
			res.EffectiveRowID = newRow(mf, &end, ch.AmountCents)
		}
		if origTo == nil || origTo.After(ml) {
			newRow(ch.Month.Add(1).First(), origTo, r.AmountCents)
		}
	}
	sort.SliceStable(rs, func(a, b int) bool {
		ca, cb := "", ""
		if rs[a].Category != nil {
			ca = "\x01" + *rs[a].Category
		}
		if rs[b].Category != nil {
			cb = "\x01" + *rs[b].Category
		}
		if ca != cb {
			return ca < cb
		}
		return rs[a].ValidFrom.Before(rs[b].ValidFrom)
	})
	sort.Slice(res.DeletedRowIDs, func(a, b int) bool { return res.DeletedRowIDs[a] < res.DeletedRowIDs[b] })
	if res.DeletedRowIDs == nil {
		res.DeletedRowIDs = []int64{}
	}
	res.Rows = rs
	return res, nil
}

// PlanActualRow §3.10.
type PlanActualRow struct {
	Month       Month  `json:"month"`
	PlanCents   *Cents `json:"plan_cents"`
	ActualCents Cents  `json:"actual_cents"`
}

// PlanVsActual da n meses de plan contra real de una categoría o del total (M-14).
func PlanVsActual(in Input, from Month, n int, category string) ([]PlanActualRow, error) {
	p, err := prepared(in)
	if err != nil {
		return nil, err
	}
	if n < 1 || n > 12 {
		return nil, errf(CodeInvalidArgument, "plan_contra_real: count fuera de 1–12")
	}
	if category != "" {
		ok := false
		for _, c := range p.Categories {
			ok = ok || (c.Code == category && active(c.IsActive))
		}
		if !ok {
			return nil, errf(CodeInvalidArgument, "categoría %q no está activa", category)
		}
	}
	l := records(p)
	ms := l.monthly()
	out := []PlanActualRow{}
	for i := 0; i < n; i++ {
		m := from.Add(i)
		b := p.budgetMonth(ms, l, m)
		row := PlanActualRow{Month: m}
		if category == "" {
			pl := b.Total.PlanCents
			row.PlanCents, row.ActualCents = &pl, b.Total.ActualCents
		} else {
			for _, c := range b.Categories {
				if c.Category == category {
					row.PlanCents, row.ActualCents = c.PlanCents, c.ActualCents
				}
			}
		}
		out = append(out, row)
	}
	return out, nil
}

// ---------- Ocurrencias de partidas planificadas (§2.7) ----------

type occurrence struct {
	item   *PlannedItem
	date   Date
	amount Cents
}

func clipDay(y int, m Month, d int) Date {
	if d > m.Days() {
		d = m.Days()
	}
	return Date{y, m.Mon, d}
}

// occurrencesIn da las ocurrencias de una partida activa sin crédito con fecha en [from, to].
func occurrencesIn(p *PlannedItem, from, to Date) []occurrence {
	if !active(p.IsActive) || p.CreditID != nil {
		return nil
	}
	var out []occurrence
	emit := func(d Date, amt Cents) bool {
		if d.After(to) || (p.EndDate != nil && d.After(*p.EndDate)) {
			return false
		}
		if !d.Before(from) && !d.Before(p.StartDate) {
			out = append(out, occurrence{p, d, amt})
		}
		return true
	}
	s, step := p.StartDate, count1(p.IntervalCount)
	switch p.Frequency {
	case "once":
		emit(s, p.AmountCents)
	case "weekly":
		for k := 0; emit(s.AddDays(7*step*k), p.AmountCents); k++ {
		}
	case "monthly":
		for k := 0; ; k += step {
			m := s.MonthOf().Add(k)
			if !emit(clipDay(m.Year, m, s.Day), p.AmountCents) {
				break
			}
		}
	case "yearly":
		for k := 0; ; k += step {
			m := Month{s.Year + k, s.Month}
			if !emit(clipDay(m.Year, m, s.Day), p.AmountCents) {
				break
			}
		}
	case "fortnightly":
		f := FortnightOf(s)
		if _, _, pay := FortnightBounds(f); pay.Before(s) {
			f = f.Add(1)
		}
		for ; ; f = f.Add(step) {
			_, _, pay := FortnightBounds(f)
			amt := p.AmountCents
			if f.Half == 2 && p.SecondHalfAmountCents != nil {
				amt = *p.SecondHalfAmountCents
			}
			if !emit(pay, amt) {
				break
			}
		}
	}
	return out
}

// linked: lectura de trabajo PROVISIONAL (docs/estado/bloqueos/S2-flujo-quincenal.md): una captura
// con ese planned_item_id en la misma quincena (semanal o quincenal) o en el mismo mes (mensual o mayor).
func (in *Input) linked(o occurrence) bool {
	same := func(d Date) bool {
		if o.item.Frequency == "weekly" || o.item.Frequency == "fortnightly" {
			return FortnightOf(d) == FortnightOf(o.date)
		}
		return d.MonthOf() == o.date.MonthOf()
	}
	for _, e := range in.Expenses {
		if e.PlannedItemID != nil && *e.PlannedItemID == o.item.ID && same(e.ExpenseDate) {
			return true
		}
	}
	for _, x := range in.Incomes {
		if x.PlannedItemID != nil && *x.PlannedItemID == o.item.ID && same(x.IncomeDate) {
			return true
		}
	}
	return false
}

func (in *Input) allOccurrences(from, to Date) []occurrence {
	var out []occurrence
	for i := range in.PlannedItems {
		out = append(out, occurrencesIn(&in.PlannedItems[i], from, to)...)
	}
	return out
}

// ---------- Flujo por quincena (§3.11) ----------

// FlowIncome es el bloque de ingresos de una quincena.
type FlowIncome struct {
	FixedCents       Cents `json:"fixed_cents"`
	ExtraCents       Cents `json:"extra_cents"`
	ExpectedCents    Cents `json:"expected_cents"`
	IncomeIsExpected bool  `json:"income_is_expected"`
	TotalCents       Cents `json:"total_cents"`
}

// FlowCredit es el bloque de pagos programados de una quincena.
type FlowCredit struct {
	ByPerson                 CentsMap `json:"by_person"`
	HouseholdCents           Cents    `json:"household_cents"`
	ExternalInformativeCents Cents    `json:"external_informative_cents"`
	HouseholdPaidCents       Cents    `json:"household_paid_cents"`
	HouseholdPendingCents    Cents    `json:"household_pending_cents"`
}

// FlowFixed es el bloque de gastos fijos y venideros de una quincena.
type FlowFixed struct {
	PlannedCents Cents `json:"planned_cents"`
	PaidCents    Cents `json:"paid_cents"`
	PendingCents Cents `json:"pending_cents"`
}

// FortnightRow §3.11.
type FortnightRow struct {
	Fortnight              Fortnight  `json:"fortnight"`
	Start                  Date       `json:"start"`
	End                    Date       `json:"end"`
	PayDate                Date       `json:"pay_date"`
	Income                 FlowIncome `json:"income"`
	CreditPayments         FlowCredit `json:"credit_payments"`
	FixedUpcoming          FlowFixed  `json:"fixed_upcoming"`
	VariableSpentCents     Cents      `json:"variable_spent_cents"`
	VariableEstimatedCents *Cents     `json:"variable_estimated_cents"`
	SavingsPlannedCents    *Cents     `json:"savings_planned_cents"`
	SpentCents             Cents      `json:"spent_cents"`
	CommittedCents         Cents      `json:"committed_cents"`
	BalanceCents           Cents      `json:"balance_cents"`
	CumulativeBalanceCents Cents      `json:"cumulative_balance_cents"`
	Provisional            []string   `json:"provisional"`
}

const (
	blockFlow       = "S2-flujo-quincenal"
	blockProjection = "S2-proyeccion-interruptor"
)

// FortnightFlow calcula el flujo del hogar por quincena (§3.11).
func FortnightFlow(in Input, from Fortnight, count int) ([]FortnightRow, error) {
	p, err := prepared(in)
	if err != nil {
		return nil, err
	}
	if count < 1 || count > 24 {
		return nil, errf(CodeInvalidArgument, "flujo_quincenal: count fuera de 1–24")
	}
	T, h := p.today, p.holder
	direct := map[string]bool{}
	for _, ps := range p.Persons {
		direct[ps.Code] = ps.PaysShareDirectly
	}
	l := records(p)
	out := []FortnightRow{}
	var cum Cents
	for i := 0; i < count; i++ {
		f := from.Add(i)
		s, e, pay := FortnightBounds(f)
		r := FortnightRow{Fortnight: f, Start: s, End: e, PayDate: pay}
		r.CreditPayments.ByPerson = newCentsMap(p.keys)
		for _, x := range p.Incomes {
			if x.IncomeDate.in(s, e) {
				if x.Kind == "fixed" {
					r.Income.FixedCents += x.AmountCents
				} else {
					r.Income.ExtraCents += x.AmountCents
				}
			}
		}
		for _, o := range p.allOccurrences(s, e) {
			lk := p.linked(o)
			if lk {
				r.Provisional = []string{blockFlow}
			}
			switch o.item.Kind {
			case "income":
				if !lk {
					r.Income.ExpectedCents += o.amount
				}
			default:
				r.FixedUpcoming.PlannedCents += o.amount
				if !lk {
					r.FixedUpcoming.PendingCents += o.amount
				}
			}
		}
		r.Income.IncomeIsExpected = r.Income.ExpectedCents > 0
		r.Income.TotalCents = r.Income.FixedCents + r.Income.ExtraCents + r.Income.ExpectedCents
		for _, cp := range l.purchs {
			for _, py := range cp.Payments {
				if !py.DueDate.in(s, e) {
					continue
				}
				var hs, hp Cents
				for _, sh := range py.Shares {
					r.CreditPayments.ByPerson.Vals[sh.Person] += sh.AmountCents
					if direct[sh.Person] {
						r.CreditPayments.ExternalInformativeCents += sh.AmountCents
					}
					if sh.Person == h {
						hs += sh.AmountCents
					}
				}
				for _, pd := range py.Paid {
					if pd.Person == h && !pd.PaidDate.After(T) {
						hp += pd.AmountCents
					}
				}
				if hp > hs {
					hp = hs
				}
				r.CreditPayments.HouseholdPaidCents += hp
			}
		}
		r.CreditPayments.HouseholdCents = r.CreditPayments.ByPerson.Get(h)
		r.CreditPayments.HouseholdPendingCents = r.CreditPayments.HouseholdCents - r.CreditPayments.HouseholdPaidCents
		for _, x := range p.Expenses {
			if x.PaymentMethod == "credit" || !x.ExpenseDate.in(s, e) || (x.PaidBy != nil && *x.PaidBy != h) {
				continue
			}
			past := !x.ExpenseDate.After(T)
			switch {
			case x.Kind == "variable" && past:
				r.VariableSpentCents += x.AmountCents
			case x.Kind != "variable" && past:
				r.FixedUpcoming.PaidCents += x.AmountCents
			case x.Kind != "variable":
				r.FixedUpcoming.PendingCents += x.AmountCents
			}
		}
		r.SpentCents = r.VariableSpentCents + r.FixedUpcoming.PaidCents + r.CreditPayments.HouseholdPaidCents
		r.CommittedCents = r.CreditPayments.HouseholdPendingCents + r.FixedUpcoming.PendingCents
		r.BalanceCents = r.Income.TotalCents - r.SpentCents - r.CommittedCents
		cum += r.BalanceCents
		r.CumulativeBalanceCents = cum
		out = append(out, r)
	}
	return out, nil
}
