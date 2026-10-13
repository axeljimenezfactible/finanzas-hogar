package presupuesto_test

import (
	"fmt"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/presupuesto"
	"finanzas-hogar/internal/testutil"
)

type env struct {
	t        *testing.T
	s        *testutil.Server
	c        *testutil.Client
	cat      int64
	goal     int64
	planned  int64
	keyCount int
}

// setup arma el servidor con fixtures mínimas: una categoría (despensa), una meta de ahorro y la quincena planificada de Axel.
func setup(t *testing.T) *env {
	t.Helper()
	s := testutil.NewServer(t, presupuesto.Register)
	e := &env{t: t, s: s, c: s.SessionFor(t, testutil.AxelUser)}
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(s.DB.QueryRow(`INSERT INTO category (code, name) VALUES ('groceries', 'Despensa') RETURNING id`).Scan(&e.cat))
	must(s.DB.QueryRow(`INSERT INTO savings_goal (name, kind, target_cents, target_date, planned_monthly_cents, start_date)
		VALUES ('Fondo de emergencia', 'goal', 629344, '2027-04-10', 50000, '2026-10-01') RETURNING id`).Scan(&e.goal))
	must(s.DB.QueryRow(`INSERT INTO planned_item (kind, description, amount_cents, second_half_amount_cents, frequency, interval_count, start_date, accrual_months, person_id)
		VALUES ('income', 'Sueldo', 1000000, 1100000, 'fortnightly', 1, '2026-01-15', 1, ?) RETURNING id`, s.Fx.AxelPersonID).Scan(&e.planned))
	return e
}

func (e *env) key() string { e.keyCount++; return fmt.Sprintf("prueba-presupuesto-%04d", e.keyCount) }

func (e *env) want(res *testutil.Response, code int) *testutil.Response {
	e.t.Helper()
	if res.StatusCode != code {
		e.t.Fatalf("esperaba %d y llegó %d: %s", code, res.StatusCode, res.Body)
	}
	return res
}

func (e *env) get(path string, code int) *testutil.Response {
	return e.want(e.c.Do("GET", path, nil), code)
}

func (e *env) anon(path string) {
	e.t.Helper()
	e.want(e.s.Anonymous(e.t).Do("GET", path, nil), 401)
}

func (e *env) budget(amount int64, from string) api.Budget {
	e.t.Helper()
	var b api.Budget
	e.want(e.c.Do("POST", "/api/v1/budgets", map[string]any{"category_id": e.cat, "amount_cents": amount, "valid_from_date": from}, "Idempotency-Key", e.key()), 201).JSON(e.t, &b)
	return b
}

func (e *env) income(kind string, amount int64, date string) api.Income {
	e.t.Helper()
	var x api.Income
	body := map[string]any{"kind": kind, "description": "Quincena", "amount_cents": amount, "income_date": date}
	if kind == "fixed" {
		body["planned_item_id"] = e.planned
	}
	e.want(e.c.Do("POST", "/api/v1/incomes", body, "Idempotency-Key", e.key()), 201).JSON(e.t, &x)
	return x
}

func (e *env) movement(kind string, amount int64) api.SavingsMovement {
	e.t.Helper()
	var m api.SavingsMovement
	e.want(e.c.Do("POST", "/api/v1/savings-movements", map[string]any{"savings_goal_id": e.goal, "movement_date": "2026-10-05", "kind": kind, "amount_cents": amount},
		"Idempotency-Key", e.key()), 201).JSON(e.t, &m)
	return m
}

func TestOp_listPeriods(t *testing.T) {
	e := setup(t)
	var pg api.PeriodPage
	e.get("/api/v1/periods?period=2026-10", 200).JSON(t, &pg)
	if len(pg.Items) != 3 || *pg.Items[2].PayDate != "2026-10-30" || *pg.Items[0].Label != "2026-10" {
		t.Fatalf("%+v", pg.Items)
	}
	e.get("/api/v1/periods?period=2026-13", 400)
	e.get("/api/v1/periods?period=2020-01-01..2026-01-31", 422)
}

func TestOp_listBudgets(t *testing.T) {
	e := setup(t)
	e.budget(200000, "2026-10-01")
	var pg api.BudgetPage
	e.get("/api/v1/budgets?period=2026-10", 200).JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].FollowingMonthAmountCents != 200000 || *pg.Items[0].IsSingleMonth {
		t.Fatalf("%+v", pg.Items)
	}
	e.get("/api/v1/budgets?limit=500", 400)
}

func TestOp_createBudget(t *testing.T) {
	e := setup(t)
	b := e.budget(194750, "2026-10-01")
	if b.AmountCents != 194750 || b.ValidToDate != nil {
		t.Fatalf("%+v", b)
	}
	e.want(e.c.Do("POST", "/api/v1/budgets", map[string]any{"category_id": e.cat, "amount_cents": 1, "valid_from_date": "2026-12-01"}, "Idempotency-Key", e.key()), 409)
	e.want(e.c.Do("POST", "/api/v1/budgets", map[string]any{"amount_cents": 1, "valid_from_date": "2026-12-05"}, "Idempotency-Key", e.key()), 422)
	e.want(e.c.Do("POST", "/api/v1/budgets", map[string]any{"amount_cents": "1", "valid_from_date": "2026-12-01"}, "Idempotency-Key", e.key()), 400)
}

// Ejemplo de C1-datos §7.1: 194750 desde oct; cambio a 251737 en dic «solo este mes» y luego en nov «en adelante».
func TestOp_updateBudget(t *testing.T) {
	e := setup(t)
	b := e.budget(194750, "2026-10-01")
	var r1 api.Budget
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/budgets/%d?month=2026-12&scope=this_month", b.ID), map[string]any{"amount_cents": 251737}), 200).JSON(t, &r1)
	if !*r1.IsSingleMonth || *r1.FollowingMonthAmountCents != 194750 || r1.ValidFromDate != "2026-12-01" {
		t.Fatalf("this_month: %+v", r1)
	}
	var r2 api.Budget
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/budgets/%d?month=2026-11&scope=onward", b.ID), map[string]any{"amount_cents": 251737}), 200).JSON(t, &r2)
	if r2.ValidFromDate != "2026-11-01" || r2.ValidToDate != nil || *r2.FollowingMonthAmountCents != 251737 {
		t.Fatalf("onward: %+v", r2)
	}
	var n int
	if err := e.s.DB.QueryRow(`SELECT count(*) FROM budget`).Scan(&n); err != nil || n != 2 {
		t.Fatalf("renglones = %d, %v", n, err)
	}
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/budgets/%d?month=2027-01", b.ID), map[string]any{"amount_cents": 1}), 422)
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/budgets/%d?scope=siempre", b.ID), map[string]any{"amount_cents": 1}), 400)
	e.want(e.c.Do("PATCH", "/api/v1/budgets/9999", map[string]any{"notes": "x"}), 404)
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/budgets/%d", b.ID), map[string]any{"notes": "Despensa del mes"}), 200)
}

func TestOp_deleteBudget(t *testing.T) {
	e := setup(t)
	b := e.budget(100000, "2026-10-01")
	e.want(e.c.Do("DELETE", fmt.Sprintf("/api/v1/budgets/%d", b.ID), nil), 204)
	e.want(e.c.Do("DELETE", fmt.Sprintf("/api/v1/budgets/%d", b.ID), nil), 404)
}

func TestOp_getBudgetMonth(t *testing.T) {
	e := setup(t)
	e.budget(200000, "2026-10-01")
	e.income("fixed", 1000000, "2026-10-15")
	e.movement("deposit", 1000)
	var m api.BudgetMonth
	e.get("/api/v1/budget/month?period=2026-10", 200).JSON(t, &m)
	if m.Month != "2026-10" || m.PlannedSpendingCents != 200000 || m.IncomeCents != 2100000 || m.PlannedSavingsCents != 50000 ||
		m.UnassignedCents != 2100000-200000-50000 || m.DaysLeft != 20 {
		t.Fatalf("%+v", m)
	}
	e.get("/api/v1/budget/month?period=2026-10-q1", 422)
}

// Hallazgo H1 de REVISION-backend-core-r1: las categorías que crea el usuario tienen code NULL (00001_esquema.sql,
// tabla category) y el motor S2 debe aceptarlas; antes respondían 500 internal_error.
func TestOp_getBudgetMonth_CategoriaSinCodigo(t *testing.T) {
	e := setup(t)
	var mascotas int64
	if err := e.s.DB.QueryRow(`INSERT INTO category (name, kind) VALUES ('Mascotas', 'variable') RETURNING id`).Scan(&mascotas); err != nil {
		t.Fatal(err)
	}
	e.budget(200000, "2026-10-01")
	e.want(e.c.Do("POST", "/api/v1/budgets", map[string]any{"category_id": mascotas, "amount_cents": 80000, "valid_from_date": "2026-10-01"},
		"Idempotency-Key", e.key()), 201)
	var m api.BudgetMonth
	e.get("/api/v1/budget/month?period=2026-10", 200).JSON(t, &m)
	if m.PlannedSpendingCents != 280000 {
		t.Fatalf("plan con la categoría sin código: %+v", m)
	}
	var rows api.BudgetCategoryRows
	e.get("/api/v1/budget/categories?period=2026-10", 200).JSON(t, &rows)
	found := false
	for _, r := range rows.Rows {
		if r.CategoryID == mascotas && r.PlannedCents != nil && *r.PlannedCents == 80000 {
			found = true
		}
	}
	if !found {
		t.Fatalf("falta la fila de la categoría sin código: %+v", rows)
	}
	var tr api.BudgetMonthlyTrend
	e.get(fmt.Sprintf("/api/v1/budget/monthly-trend?category_id=%d", mascotas), 200).JSON(t, &tr)
	if tr.CategoryID == nil || *tr.CategoryID != mascotas {
		t.Fatalf("%+v", tr)
	}
	for _, path := range []string{"/api/v1/budget/fortnights?period=2026-10", "/api/v1/budget/alerts", "/api/v1/debt/live",
		"/api/v1/savings/summary", "/api/v1/cash-flow/fortnights"} {
		e.get(path, 200)
	}
}

func TestOp_getBudgetByCategory(t *testing.T) {
	e := setup(t)
	e.budget(200000, "2026-10-01")
	var rows api.BudgetCategoryRows
	e.get("/api/v1/budget/categories", 200).JSON(t, &rows)
	if rows.Month != "2026-10" || len(rows.Rows) != 1 || rows.Rows[0].CategoryID != e.cat || *rows.Rows[0].PlannedCents != 200000 {
		t.Fatalf("%+v", rows)
	}
	e.get("/api/v1/budget/categories?period=octubre", 400)
}

func TestOp_getBudgetFortnights(t *testing.T) {
	e := setup(t)
	e.budget(200000, "2027-02-01")
	var f api.BudgetFortnights
	e.get("/api/v1/budget/fortnights?period=2027-02", 200).JSON(t, &f)
	if len(f.Fortnights) != 2 || f.Fortnights[1].PayDate != "2027-02-28" || f.Fortnights[0].PlannedCents+f.Fortnights[1].PlannedCents != 200000 {
		t.Fatalf("%+v", f)
	}
	e.get("/api/v1/budget/fortnights?period=2027-W05", 422)
}

func TestOp_getAccruedGroceries(t *testing.T) {
	e := setup(t)
	var a api.AccruedGroceries
	e.get("/api/v1/budget/accrued-groceries?period=2026-11", 200).JSON(t, &a)
	if a.Month != "2026-11" || len(a.Rows) != 0 {
		t.Fatalf("%+v", a)
	}
	e.get("/api/v1/budget/accrued-groceries?period=2026-11-q3", 400)
}

func TestOp_getBudgetAlerts(t *testing.T) {
	e := setup(t)
	var a api.BudgetAlerts
	e.get("/api/v1/budget/alerts", 200).JSON(t, &a)
	if len(a.Rows) != 0 {
		t.Fatalf("%+v", a)
	}
	e.anon("/api/v1/budget/alerts")
}

func TestOp_getBudgetMonthlyTrend(t *testing.T) {
	e := setup(t)
	e.budget(200000, "2026-10-01")
	var m api.BudgetMonthlyTrend
	e.get(fmt.Sprintf("/api/v1/budget/monthly-trend?category_id=%d", e.cat), 200).JSON(t, &m)
	if len(m.Months) != 12 || m.Months[11].Month != "2026-10" || *m.CategoryID != e.cat {
		t.Fatalf("%+v", m)
	}
	e.get("/api/v1/budget/monthly-trend?category_id=1,2", 422)
}

func TestOp_getCashFlowFortnights(t *testing.T) {
	e := setup(t)
	var c api.CashFlowFortnights
	e.get("/api/v1/cash-flow/fortnights", 200).JSON(t, &c)
	if c.ProjectionMode != "contracted_only" || len(c.Rows) != 24 || c.Rows[0].StartDate != "2026-10-01" {
		t.Fatalf("%s %d", c.ProjectionMode, len(c.Rows))
	}
	e.anon("/api/v1/cash-flow/fortnights")
}

func TestOp_getAccrualVsCash(t *testing.T) {
	e := setup(t)
	var a api.AccrualVsCash
	e.get("/api/v1/cash-flow/accrual-vs-cash", 200).JSON(t, &a)
	if len(a.Months) != 12 || a.Months[0].Month != "2026-10" || !a.Months[0].IsCanonical || a.Months[6].IsCanonical {
		t.Fatalf("%+v", a.Months)
	}
	for _, m := range a.Months {
		if m.CumulativeDifferenceCents != m.CumulativeAccruedCents-m.CumulativeCashCents {
			t.Fatalf("identidad: %+v", m)
		}
	}
	e.anon("/api/v1/cash-flow/accrual-vs-cash")
}

func TestOp_getCashFlowAssumptions(t *testing.T) {
	e := setup(t)
	var a api.CashFlowAssumptions
	e.get("/api/v1/cash-flow/assumptions", 200).JSON(t, &a)
	if len(a.PlannedItems) != 1 || a.PlannedItems[0].ID != e.planned || a.VariableEstimateMethod == "" {
		t.Fatalf("%+v", a)
	}
	e.anon("/api/v1/cash-flow/assumptions")
}

func TestOp_getLiveDebt(t *testing.T) {
	e := setup(t)
	var d api.LiveDebt
	e.get("/api/v1/debt/live?date=2027-04-10", 200).JSON(t, &d)
	if d.Date != "2027-04-10" || d.TotalCents != 0 {
		t.Fatalf("%+v", d)
	}
	e.get("/api/v1/debt/live", 200).JSON(t, &d)
	if d.Date != "2027-09-30" {
		t.Fatalf("horizonte: %s", d.Date)
	}
	e.get("/api/v1/debt/live?date=2027-02-30", 400)
}

func TestOp_listIncomes(t *testing.T) {
	e := setup(t)
	e.income("fixed", 1000000, "2026-10-15")
	e.income("extra", 5000, "2026-09-20")
	var pg api.IncomePage
	e.get("/api/v1/incomes?period=2026-10&kind=fixed", 200).JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].PersonID != e.s.Fx.AxelPersonID {
		t.Fatalf("%+v", pg)
	}
	e.get("/api/v1/incomes?kind=bono", 400)
}

func TestOp_createIncome(t *testing.T) {
	e := setup(t)
	x := e.income("extra", 25000, "2026-10-10")
	if x.CreatedByUserAccountID == nil || *x.CreatedByUserAccountID != e.s.Fx.AxelUserID || x.PersonID != e.s.Fx.AxelPersonID {
		t.Fatalf("%+v", x)
	}
	e.want(e.c.Do("POST", "/api/v1/incomes", map[string]any{"kind": "extra", "description": "x", "amount_cents": 0, "income_date": "2026-10-10"}, "Idempotency-Key", e.key()), 422)
	e.want(e.c.Do("POST", "/api/v1/incomes", map[string]any{"kind": "extra", "description": "x", "amount_cents": 5}), 400)
}

func TestOp_getIncomeMonth(t *testing.T) {
	e := setup(t)
	x := e.income("fixed", 990000, "2026-10-15")
	e.income("extra", 5000, "2026-10-20")
	var m api.IncomeMonth
	e.get("/api/v1/incomes/month?period=2026-10", 200).JSON(t, &m)
	if len(m.Fortnights) != 2 || m.Fortnights[0].Status != "captured" || *m.Fortnights[0].IncomeID != x.ID || m.Fortnights[1].ExpectedCents != 1100000 ||
		m.CapturedCents != 990000 || m.ExpectedPendingCents != 1100000 || m.ExtrasCents != 5000 || m.TotalCents != 2095000 {
		t.Fatalf("%+v", m)
	}
	e.get("/api/v1/incomes/month?period=2026", 422)
	// K-2 (M-18, C1 v1.5): mes anterior con la misma cuenta y diferencia de quincenas con signo.
	var sep api.IncomeMonth
	e.get("/api/v1/incomes/month?period=2026-09", 200).JSON(t, &sep)
	e.income("extra", 7000, "2026-09-05")
	e.get("/api/v1/incomes/month?period=2026-10", 200).JSON(t, &m)
	pf := sep.CapturedCents + sep.ExpectedPendingCents
	if m.PreviousMonth == nil || *m.PreviousMonth != "2026-09" || m.PreviousFortnightsCents == nil || *m.PreviousFortnightsCents != pf ||
		pf == 0 || *m.PreviousExtrasCents != 7000 || *m.FortnightsDifferenceCents != m.CapturedCents+m.ExpectedPendingCents-pf {
		t.Fatalf("mes anterior: %+v (septiembre %+v)", m, sep)
	}
	// Mes anterior sin quincenas ni extras: sin datos no es cero (S6-tablero-sin-datos): cifras ausentes, mes presente.
	var old api.IncomeMonth
	e.get("/api/v1/incomes/month?period=2000-01", 200).JSON(t, &old)
	if old.PreviousMonth == nil || *old.PreviousMonth != "1999-12" || old.PreviousFortnightsCents != nil || old.PreviousExtrasCents != nil ||
		old.FortnightsDifferenceCents != nil {
		t.Fatalf("mes anterior sin datos: %+v", old)
	}
}

func TestOp_getIncomeMonthlyTrend(t *testing.T) {
	e := setup(t)
	e.income("fixed", 1000000, "2026-10-15")
	e.income("extra", 5000, "2026-09-20")
	var m api.IncomeMonthlyTrend
	e.get("/api/v1/incomes/monthly-trend", 200).JSON(t, &m)
	if len(m.Months) != 12 || m.Months[11].FixedCents != 1000000 || m.Months[10].ExtraCents != 5000 || m.Months[0].Month != "2025-11" {
		t.Fatalf("%+v", m)
	}
	e.anon("/api/v1/incomes/monthly-trend")
}

func TestOp_updateIncome(t *testing.T) {
	e := setup(t)
	x := e.income("extra", 5000, "2026-10-20")
	var u api.Income
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/incomes/%d", x.ID), map[string]any{"amount_cents": 6000, "notes": "Bono"}), 200).JSON(t, &u)
	if u.AmountCents != 6000 || *u.Notes != "Bono" {
		t.Fatalf("%+v", u)
	}
	e.want(e.c.Do("PATCH", "/api/v1/incomes/9999", map[string]any{"amount_cents": 1}), 404)
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/incomes/%d", x.ID), map[string]any{"amount_cents": -5}), 422)
}

func TestOp_deleteIncome(t *testing.T) {
	e := setup(t)
	x := e.income("fixed", 1000000, "2026-10-15")
	e.want(e.c.Do("DELETE", fmt.Sprintf("/api/v1/incomes/%d", x.ID), nil), 204)
	e.want(e.c.Do("DELETE", fmt.Sprintf("/api/v1/incomes/%d", x.ID), nil), 404)
	var m api.IncomeMonth
	e.get("/api/v1/incomes/month", 200).JSON(t, &m)
	if m.Fortnights[0].Status != "expected" {
		t.Fatal("una quincena eliminada vuelve a «esperada»")
	}
}

func TestOp_getSavingsSummary(t *testing.T) {
	e := setup(t)
	e.movement("deposit", 314672)
	var s api.SavingsSummary
	e.get("/api/v1/savings/summary", 200).JSON(t, &s)
	if s.SavedTotalCents != 314672 || s.ContributedMonthCents != 314672 || s.PlannedMonthCents != 50000 || len(s.Goals) != 1 {
		t.Fatalf("%+v", s)
	}
	e.anon("/api/v1/savings/summary")
}

func TestOp_getSavingsMonthlyTrend(t *testing.T) {
	e := setup(t)
	e.movement("deposit", 1000)
	var s api.SavingsMonthlyTrend
	e.get("/api/v1/savings/monthly-trend", 200).JSON(t, &s)
	if len(s.Months) != 12 || s.Months[11].SavedEndCents != 1000 || s.Months[10].SavedEndCents != 0 {
		t.Fatalf("%+v", s)
	}
	e.anon("/api/v1/savings/monthly-trend")
}

// RN-18: 3,146.72 aportados a una meta de 6,293.44 = 50.0 % (500 décimas).
func TestOp_listSavingsGoals(t *testing.T) {
	e := setup(t)
	e.movement("deposit", 314672)
	var pg api.SavingsGoalPage
	e.get("/api/v1/savings-goals", 200).JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].ProgressTenthsPct != 500 || *pg.Items[0].ProgressPct != 50 || *pg.Items[0].Status != "in_progress" {
		t.Fatalf("%+v", pg.Items)
	}
	e.get("/api/v1/savings-goals?status=todas", 400)
}

func TestOp_listSavingsMovements(t *testing.T) {
	e := setup(t)
	e.movement("deposit", 1000)
	e.movement("withdrawal", 300)
	var pg api.SavingsMovementPage
	e.get(fmt.Sprintf("/api/v1/savings-movements?period=2026-10&savings_goal_id=%d&limit=1", e.goal), 200).JSON(t, &pg)
	if len(pg.Items) != 1 || !pg.NextCursor.Valid {
		t.Fatalf("%+v", pg)
	}
	e.get("/api/v1/savings-movements?savings_goal_id=abc", 400)
}

func TestOp_createSavingsMovement(t *testing.T) {
	e := setup(t)
	m := e.movement("deposit", 1000)
	if m.SavingsGoalID != e.goal || m.Kind != "deposit" {
		t.Fatalf("%+v", m)
	}
	e.want(e.c.Do("POST", "/api/v1/savings-movements", map[string]any{"savings_goal_id": 9999, "movement_date": "2026-10-05", "kind": "deposit", "amount_cents": 1},
		"Idempotency-Key", e.key()), 422)
}

func TestOp_updateSavingsMovement(t *testing.T) {
	e := setup(t)
	m := e.movement("deposit", 1000)
	var u api.SavingsMovement
	e.want(e.c.Do("PATCH", fmt.Sprintf("/api/v1/savings-movements/%d", m.ID), map[string]any{"amount_cents": 1500, "notes": nil}), 200).JSON(t, &u)
	if u.AmountCents != 1500 {
		t.Fatalf("%+v", u)
	}
	e.want(e.c.Do("PATCH", "/api/v1/savings-movements/9999", map[string]any{"amount_cents": 1}), 404)
}

func TestOp_deleteSavingsMovement(t *testing.T) {
	e := setup(t)
	m := e.movement("deposit", 1000)
	e.want(e.c.Do("DELETE", fmt.Sprintf("/api/v1/savings-movements/%d", m.ID), nil), 204)
	e.want(e.c.Do("DELETE", fmt.Sprintf("/api/v1/savings-movements/%d", m.ID), nil), 404)
}
