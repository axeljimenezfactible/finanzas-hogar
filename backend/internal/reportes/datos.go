package reportes

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
)

// Monthly es el reporte del mes: la respuesta de getMonthlyReport más el renglón del mes de getAccrualVsCash
// (devengado, efectivo y diferencia acumulada, R-45…R-50), calculado con las mismas llamadas al motor S2.
type Monthly struct {
	Report     api.MonthlyReport
	Accrual    api.AccrualVsCashMonthsItem
	Categories map[int64]string
	Persons    map[int64]string
	Externals  []int64 // personas externas en el orden de external_participation
}

func c64(c mf.Cents) int64 { return int64(c) }

func cptr(c *mf.Cents) *int64 {
	if c == nil {
		return nil
	}
	return api.Ptr(int64(*c))
}

func names(ctx context.Context, q dominio.Querier, query string) (map[int64]string, error) {
	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]string{}
	for rows.Next() {
		var id int64
		var n string
		if err := rows.Scan(&id, &n); err != nil {
			return nil, err
		}
		out[id] = n
	}
	return out, rows.Err()
}

func externalPersons(ctx context.Context, q dominio.Querier) ([]int64, error) {
	rows, err := q.QueryContext(ctx, `SELECT id FROM person WHERE participation = 'external' ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

// BuildMonthly arma el reporte del mes que contiene month con la fecha de hoy today.
func BuildMonthly(ctx context.Context, db *sql.DB, today, month dominio.Date) (*Monthly, error) {
	m := dominio.MonthStart(month.Year, month.Month)
	end := dominio.MonthEnd(m)
	em := dominio.EngineMonth(m)
	fe, err := dominio.LoadFinEngine(ctx, db)
	if err != nil {
		return nil, err
	}
	in, err := fe.At(today, nil)
	if err != nil {
		return nil, fmt.Errorf("motor_fin: entrada inválida: %w", err)
	}
	out := &Monthly{}
	r := &out.Report
	r.Month = em.String()
	// Presupuesto del mes (las mismas llamadas que getBudgetMonth, getBudgetByCategory y getBudgetFortnights).
	bm, err := mf.BudgetMonth(in, em)
	if err != nil {
		return nil, err
	}
	var income int64
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount_cents), 0) FROM income WHERE income_date BETWEEN ? AND ?`,
		m.String(), end.String()).Scan(&income); err != nil {
		return nil, err
	}
	ref := today
	if dominio.EngineMonth(ref) != em {
		ref = m
	}
	sin, err := fe.At(ref, nil)
	if err != nil {
		return nil, err
	}
	sv, err := mf.SavingsGoals(sin)
	if err != nil {
		return nil, err
	}
	r.Budget = api.BudgetMonth{Month: em.String(), IncomeCents: income, PlannedSpendingCents: c64(bm.Total.PlanCents),
		PlannedSavingsCents: c64(sv.MonthSummary.PlannedThisMonthCents), ActualCents: c64(bm.Total.ActualCents),
		RemainingCents: c64(bm.Total.RemainingCents), DaysLeft: int64(bm.DaysLeft)}
	r.Budget.UnassignedCents = r.Budget.IncomeCents - r.Budget.PlannedSpendingCents - r.Budget.PlannedSavingsCents
	if bm.Total.UsedPct != nil {
		r.Budget.UsedPct, r.Budget.Status = api.Ptr(*bm.Total.UsedPct), api.Ptr(bm.Total.State)
	}
	counts, err := dominio.ExpenseCountByCategory(ctx, db, m.String(), end.String())
	if err != nil {
		return nil, err
	}
	r.Categories = api.BudgetCategoryRows{Month: em.String(), Rows: []api.BudgetCategoryRowsRowsItem{}}
	for _, l := range bm.Categories {
		id, ok := fe.CategoryID[l.Category]
		if !ok {
			continue
		}
		it := api.BudgetCategoryRowsRowsItem{CategoryID: id, PlannedCents: cptr(l.PlanCents), ActualCents: c64(l.ActualCents),
			DifferenceCents: cptr(l.RemainingCents), ExpenseCount: counts[id]}
		if l.UsedPct != nil {
			it.UsedPct, it.Status = api.Ptr(*l.UsedPct), api.Ptr(l.State)
		}
		r.Categories.Rows = append(r.Categories.Rows, it)
	}
	q1 := dominio.FortnightOf(m)
	q2 := q1.Next()
	// El motor pide ventanas sin huecos desde el primer costo (§2.12 reglas 2 y 3): si hay gastos antes de la Q1 del
	// mes, se antepone una ventana «antes» que no se muestra.
	wins := []mf.Window{{Label: "q1", Start: dominio.EngineDate(q1.Start()), End: dominio.EngineDate(q1.End())},
		{Label: "q2", Start: dominio.EngineDate(q2.Start()), End: dominio.EngineDate(q2.End())}}
	first := wins[0].Start
	for _, x := range fe.Expenses() {
		if x.ExpenseDate.Before(first) {
			first = x.ExpenseDate
		}
	}
	if first.Before(wins[0].Start) {
		wins = append([]mf.Window{{Label: "antes", Start: first, End: dominio.EngineDate(q1.Start().AddDays(-1))}}, wins...)
	}
	ps, err := mf.Periods(in, mf.PeriodSpec{Windows: wins})
	if err != nil {
		return nil, err
	}
	ps = ps[len(ps)-2:]
	// Igual que getBudgetFortnights (PROVISIONAL del área de presupuesto, bloqueo S6-presupuesto-quincenas).
	plan := c64(bm.Total.PlanCents)
	p1 := dominio.MulDivRound(plan, 15, int64(end.Day))
	r.Fortnights = api.BudgetFortnights{Month: em.String()}
	for i, f := range []dominio.Fortnight{q1, q2} {
		var actual int64
		if i < len(ps) {
			actual = c64(ps[i].Accrued.ByPerson.Get(fe.Holder))
		}
		pl := p1
		if i == 1 {
			pl = plan - p1
		}
		r.Fortnights.Fortnights = append(r.Fortnights.Fortnights, api.BudgetFortnightsFortnightsItem{Half: int64(f.Half),
			StartDate: f.Start().String(), EndDate: f.End().String(), PayDate: f.PayDay().String(), ActualCents: actual, PlannedCents: pl})
	}
	// Devengado contra efectivo del mes (motor S2): como getAccrualVsCash.
	mp, err := mf.Periods(in, mf.PeriodSpec{Months: []mf.Month{em}})
	if err != nil {
		return nil, err
	}
	ld, err := mf.LiveDebt(in, em.Last())
	if err != nil {
		return nil, err
	}
	if len(mp) == 1 {
		p := mp[0]
		a := api.AccrualVsCashMonthsItem{Month: em.String(), AccruedCents: c64(p.Accrued.TotalCents), CashCents: c64(p.Cash.TotalCents),
			AccruedHouseholdCents: c64(p.Accrued.ByPerson.Get(fe.Holder)), CashHouseholdCents: c64(p.Cash.ByPerson.Get(fe.Holder)),
			CumulativeAccruedCents: c64(ld.AccruedToDateCents), CumulativeDifferenceCents: c64(p.CumulativeDiffCents),
			IsCanonical: em.String() >= "2026-10" && em.String() <= "2027-03"}
		a.CumulativeCashCents = a.CumulativeAccruedCents - a.CumulativeDifferenceCents
		for _, x := range fe.Externals {
			a.AccruedExternalCents += c64(p.Accrued.ByPerson.Get(x))
			a.CashExternalCents += c64(p.Cash.ByPerson.Get(x))
		}
		out.Accrual = a
	}
	r.LiveDebtEndCents = c64(ld.LiveDebtCents)
	if r.Credits, err = creditTotals(ctx, db, m, end); err != nil {
		return nil, err
	}
	r.Credits.LiveDebtCents = r.LiveDebtEndCents
	// PROVISIONAL (bloqueo S7-reporte-mensual-secciones): el plan de comidas y la canasta los calcula su área; mientras no
	// haya un cálculo compartido, el reporte no los inventa: semanas vacías y sin canasta.
	r.Meals = api.MealPlanTrends{Weeks: []api.MealPlanTrendsWeeksItem{}}
	exts, err := externalPersons(ctx, db)
	if err != nil {
		return nil, err
	}
	out.Externals = exts
	r.ExternalParticipation = []api.PersonParticipation{}
	for _, pid := range exts {
		pp := api.PersonParticipation{PersonID: pid, Month: em.String(), Purchases: []api.PersonParticipationPurchasesItem{},
			PaymentShares: []api.PersonParticipationPaymentSharesItem{}}
		if len(mp) == 1 {
			pp.AccruedGroceriesCents = c64(mp[0].Accrued.ByPerson.Get(fe.PersonCode[pid]))
		}
		shares, err := personShares(ctx, db, pid)
		if err != nil {
			return nil, err
		}
		for _, s := range shares {
			if s.Due.Compare(m) >= 0 && s.Due.Compare(end) <= 0 {
				pp.PaymentShares = append(pp.PaymentShares, api.PersonParticipationPaymentSharesItem{ScheduledPaymentID: s.ID,
					DueDate: s.Due.String(), AmountCents: s.Amount, Status: s.status(today)})
			}
			if s.Paid < s.Amount {
				pp.LiveDebtCents += s.Amount - s.Paid
			}
		}
		if pp.Purchases, err = personPurchases(ctx, db, pid, end, shares); err != nil {
			return nil, err
		}
		r.ExternalParticipation = append(r.ExternalParticipation, pp)
	}
	if out.Categories, err = names(ctx, db, `SELECT id, name FROM category`); err != nil {
		return nil, err
	}
	if out.Persons, err = names(ctx, db, `SELECT id, display_name FROM person`); err != nil {
		return nil, err
	}
	return out, nil
}

// creditTotals suma los pagos programados del mes por papel de la persona. PROVISIONAL: la «quincena» de un mes
// cerrado es su segunda quincena (16 al último día); los intereses pagados son los de pagos ya cubiertos.
func creditTotals(ctx context.Context, db *sql.DB, m, end dominio.Date) (api.CreditTotals, error) {
	var t api.CreditTotals
	rows, err := db.QueryContext(ctx, `SELECT sp.id, sp.due_date, sp.amount_cents, sp.interest_cents, sps.amount_cents, p.participation,
		COALESCE((SELECT SUM(cp.amount_cents) FROM credit_payment cp WHERE cp.scheduled_payment_id = sp.id), 0)
		FROM scheduled_payment sp JOIN scheduled_payment_share sps ON sps.scheduled_payment_id = sp.id
		JOIN person p ON p.id = sps.person_id JOIN credit_purchase pu ON pu.id = sp.credit_purchase_id
		WHERE pu.status <> 'cancelled' AND sp.due_date BETWEEN ? AND ? ORDER BY sp.due_date, sp.id, sps.person_id`,
		fmt.Sprintf("%04d-01-01", m.Year), end.String())
	if err != nil {
		return t, err
	}
	defer rows.Close()
	seen := map[int64]bool{}
	q2 := dominio.NewDate(m.Year, m.Month, 16).String()
	for rows.Next() {
		var id, amount, interest, share, paid int64
		var due, part string
		if err := rows.Scan(&id, &due, &amount, &interest, &share, &part, &paid); err != nil {
			return t, err
		}
		inMonth := due >= m.String()
		if !seen[id] {
			seen[id] = true
			if paid >= amount {
				t.InterestPaidYearCents += interest
				if inMonth {
					t.InterestPaidMonthCents += interest
				}
			}
		}
		if !inMonth {
			continue
		}
		ext := part == "external"
		t.MonthTotalCents += share
		if ext {
			t.MonthExternalCents += share
		} else {
			t.MonthHouseholdCents += share
		}
		if due >= q2 {
			t.FortnightTotalCents += share
			if ext {
				t.FortnightExternalCents += share
			} else {
				t.FortnightHouseholdCents += share
			}
		}
	}
	return t, rows.Err()
}

// ShareRow es la parte de una persona en un pago programado y lo que ella ha pagado (Σ credit_payment).
type ShareRow struct {
	ID       int64
	Purchase int64
	Due      dominio.Date
	Amount   int64
	Paid     int64
}

func (s ShareRow) status(today dominio.Date) string {
	switch {
	case s.Paid >= s.Amount:
		return "paid"
	case s.Due.Before(today):
		return "overdue"
	default:
		return "pending"
	}
}

func personShares(ctx context.Context, db *sql.DB, pid int64) ([]ShareRow, error) {
	rows, err := db.QueryContext(ctx, `SELECT sp.id, sp.credit_purchase_id, sp.due_date, sps.amount_cents,
		COALESCE((SELECT SUM(cp.amount_cents) FROM credit_payment cp WHERE cp.scheduled_payment_id = sp.id AND cp.payer_person_id = sps.person_id), 0)
		FROM scheduled_payment_share sps JOIN scheduled_payment sp ON sp.id = sps.scheduled_payment_id
		JOIN credit_purchase pu ON pu.id = sp.credit_purchase_id
		WHERE sps.person_id = ? AND pu.status <> 'cancelled' ORDER BY sp.due_date, sp.id`, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ShareRow
	for rows.Next() {
		var s ShareRow
		var due string
		if err := rows.Scan(&s.ID, &s.Purchase, &due, &s.Amount, &s.Paid); err != nil {
			return nil, err
		}
		if s.Due, err = dominio.ParseDate(due); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func personPurchases(ctx context.Context, db *sql.DB, pid int64, end dominio.Date, shares []ShareRow) ([]api.PersonParticipationPurchasesItem, error) {
	rows, err := db.QueryContext(ctx, `SELECT pu.id, pu.description, pu.purchase_date, s.total_cents, pu.term_count
		FROM credit_purchase_share s JOIN credit_purchase pu ON pu.id = s.credit_purchase_id
		WHERE s.person_id = ? AND pu.status <> 'cancelled' AND pu.purchase_date <= ? ORDER BY pu.purchase_date, pu.id`, pid, end.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []api.PersonParticipationPurchasesItem{}
	for rows.Next() {
		var it api.PersonParticipationPurchasesItem
		var desc sql.NullString
		if err := rows.Scan(&it.CreditPurchaseID, &desc, &it.PurchaseDate, &it.TotalCents, &it.PaymentsTotal); err != nil {
			return nil, err
		}
		if desc.Valid {
			it.Description = api.Ptr(desc.String)
		}
		it.BalanceCents = it.TotalCents
		for _, s := range shares {
			if s.Purchase != it.CreditPurchaseID {
				continue
			}
			if s.Paid >= s.Amount && !s.Due.After(end) {
				it.PaymentsMade++
			}
			it.BalanceCents -= min(s.Paid, s.Amount)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// External es el estado de cuenta de una persona externa (getExternalStatement) más lo que el PDF necesita para
// citar cada fecha de pago desde getMonthlyReport (C1 no trae due_date en payment_shares: bloqueo S7-estado-cuenta-fechas).
type External struct {
	Statement api.ExternalStatement
	Name      string
	DueRefs   []string // por payment_shares[i]: «monthly:AAAA-MM|external_participation.k.payment_shares.j.due_date»
	DueDates  []string
	Labels    []string // descripción de purchases[i]
}

// ErrNotExternal: el estado de cuenta es solo para personas externas.
var ErrNotExternal = dominio.Invalid(dominio.FieldIssue{Field: "person_id", Code: "not_allowed",
	Message: "El estado de cuenta es solo para personas externas, como Karen."})

// BuildExternal arma el estado de cuenta de pid entre start y end (inclusive).
func BuildExternal(ctx context.Context, db *sql.DB, today dominio.Date, pid int64, start, end dominio.Date) (*External, error) {
	var name, part string
	err := db.QueryRowContext(ctx, `SELECT display_name, participation FROM person WHERE id = ?`, pid).Scan(&name, &part)
	if err == sql.ErrNoRows {
		return nil, dominio.NotFound("No existe esa persona.")
	}
	if err != nil {
		return nil, err
	}
	if part != "external" {
		return nil, ErrNotExternal
	}
	exts, err := externalPersons(ctx, db)
	if err != nil {
		return nil, err
	}
	k := 0
	for i, id := range exts {
		if id == pid {
			k = i
		}
	}
	x := &External{Name: name}
	st := &x.Statement
	st.PersonID, st.PeriodStartDate, st.PeriodEndDate = pid, start.String(), end.String()
	st.Purchases, st.GroceryShares, st.PaymentShares = []api.CreditPurchase{}, []api.ExternalStatementGrocerySharesItem{}, []api.ExternalStatementPaymentSharesItem{}
	// Despensas en común: su parte en compras ligadas a un gasto de despensa que comparte con alguien más.
	rows, err := db.QueryContext(ctx, `SELECT pu.id, pu.purchase_date, s.principal_cents, s.interest_cents, s.total_cents
		FROM credit_purchase_share s JOIN credit_purchase pu ON pu.id = s.credit_purchase_id
		JOIN expense e ON e.id = pu.expense_id JOIN category c ON c.id = e.category_id
		WHERE s.person_id = ? AND c.code = ? AND pu.status <> 'cancelled' AND pu.purchase_date <= ?
		  AND (pu.status = 'active' OR pu.purchase_date >= ?)
		  AND EXISTS (SELECT 1 FROM credit_purchase_share o WHERE o.credit_purchase_id = pu.id AND o.person_id <> s.person_id)
		ORDER BY pu.purchase_date, pu.id`, pid, dominio.GroceriesCategoryCode, end.String(), start.String())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var g api.ExternalStatementGrocerySharesItem
		if err := rows.Scan(&g.CreditPurchaseID, &g.PurchaseDate, &g.PrincipalCents, &g.InterestCents, &g.TotalCents); err != nil {
			rows.Close()
			return nil, err
		}
		st.GroceryShares = append(st.GroceryShares, g)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Compras particulares: compras a crédito en las que solo participa ella.
	rows, err = db.QueryContext(ctx, `SELECT pu.id, pu.credit_id, pu.expense_id, pu.kind, pu.description, pu.purchase_date, pu.principal_cents,
		pu.interest_cents, pu.total_cents, pu.term_count, pu.payment_frequency, pu.first_payment_date, pu.interest_source, pu.rate_num,
		pu.rate_den, pu.rate_basis, pu.interest_includes_vat, pu.promotion_id, pu.status, pu.notes, pu.created_at, pu.updated_at
		FROM credit_purchase pu WHERE pu.status <> 'cancelled' AND pu.purchase_date <= ? AND (pu.status = 'active' OR pu.purchase_date >= ?)
		  AND EXISTS (SELECT 1 FROM credit_purchase_share s WHERE s.credit_purchase_id = pu.id AND s.person_id = ?)
		  AND NOT EXISTS (SELECT 1 FROM credit_purchase_share s WHERE s.credit_purchase_id = pu.id AND s.person_id <> ?)
		ORDER BY pu.purchase_date, pu.id`, end.String(), start.String(), pid, pid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p api.CreditPurchase
		var exp, rn, rd, promo sql.NullInt64
		var desc, rb, notes sql.NullString
		var vat int64
		if err := rows.Scan(&p.ID, &p.CreditID, &exp, &p.Kind, &desc, &p.PurchaseDate, &p.PrincipalCents, &p.InterestCents, &p.TotalCents,
			&p.TermCount, &p.PaymentFrequency, &p.FirstPaymentDate, &p.InterestSource, &rn, &rd, &rb, &vat, &promo, &p.Status, &notes,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		p.CreatedAt, p.UpdatedAt, p.InterestIncludesVat = api.APITime(p.CreatedAt), api.APITime(p.UpdatedAt), vat == 1
		if exp.Valid {
			p.ExpenseID = api.Ptr(exp.Int64)
		}
		if rn.Valid {
			p.RateNum, p.RateDen = api.Ptr(rn.Int64), api.Ptr(rd.Int64)
		}
		if rb.Valid {
			p.RateBasis = api.Ptr(rb.String)
		}
		if promo.Valid {
			p.PromotionID = api.Ptr(promo.Int64)
		}
		label := "Compra a crédito"
		if desc.Valid {
			p.Description, label = api.Ptr(desc.String), desc.String
		}
		if notes.Valid {
			p.Notes = api.Ptr(notes.String)
		}
		st.Purchases = append(st.Purchases, p)
		x.Labels = append(x.Labels, label)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	shares, err := personShares(ctx, db, pid)
	if err != nil {
		return nil, err
	}
	perMonth := map[string]int{}
	for _, s := range shares {
		mk := s.Due.MonthKey()
		j := perMonth[mk]
		perMonth[mk]++
		if s.Paid < s.Amount {
			st.PendingTotalCents += s.Amount - s.Paid
		}
		inPeriod := !s.Due.Before(start) && !s.Due.After(end)
		if !inPeriod && !(s.Due.After(end) && s.Paid < s.Amount) {
			continue
		}
		st.PaymentShares = append(st.PaymentShares, api.ExternalStatementPaymentSharesItem{PersonID: pid, AmountCents: s.Amount,
			PaidCents: s.Paid, Status: s.status(today)})
		x.DueDates = append(x.DueDates, s.Due.String())
		x.DueRefs = append(x.DueRefs, "monthly:"+mk+"|external_participation."+strconv.Itoa(k)+".payment_shares."+strconv.Itoa(j)+".due_date")
	}
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount_cents), 0) FROM credit_payment WHERE payer_person_id = ? AND paid_date BETWEEN ? AND ?`,
		pid, start.String(), end.String()).Scan(&st.PaidInPeriodCents); err != nil {
		return nil, err
	}
	return x, nil
}
