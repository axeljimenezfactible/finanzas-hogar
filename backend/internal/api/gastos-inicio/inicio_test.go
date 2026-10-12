package gastos_test

import (
	"fmt"
	"strconv"
	"testing"

	"finanzas-hogar/internal/api"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/testutil"
)

// inicioFx: Varios 500.00 de contado (Hogar, todo de Axel) y Despensa 1,000.00 de contado (Conjunto 50/50),
// ambos devengados en un mes, en oct-2026; el reloj de testutil está en 12-oct-2026.
func inicioFx(t *testing.T) (*testutil.Server, *testutil.Client, fx) {
	t.Helper()
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	// dominio.LoadFinEngine (api-presupuesto) aún no acepta categorías con code NULL: bloqueo S6-gastos-inicio §11.
	if _, err := s.DB.Exec(`UPDATE category SET code = 'misc' WHERE id = ?`, f.varios); err != nil {
		t.Fatal(err)
	}
	c := s.SessionFor(t, testutil.AxelUser)
	newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 100000, "expense_date": "2026-10-06", "payment_method": "cash", "accrual_months": 1}, 501)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 50000, "expense_date": "2026-10-05", "payment_method": "cash", "accrual_months": 1}, 502)
	return s, c, f
}

func TestOp_getExpenseMonthlyTrend(t *testing.T) {
	s, c, f := inicioFx(t)
	var tr api.ExpenseMonthlyTrend
	res := c.Do("GET", "/api/v1/expenses/monthly-trend", nil)
	want(t, res, 200)
	res.JSON(t, &tr)
	// Del primer mes con datos (oct) al mes en curso; empate 500.00 y 500.00 → por id (Despensa primero); sin «Otras».
	if len(tr.Months) != 1 || tr.Months[0].Month != "2026-10" || len(tr.Months[0].Categories) != 2 {
		t.Fatalf("%+v", tr)
	}
	cs := tr.Months[0].Categories
	if *cs[0].CategoryID != f.despensa || cs[0].HouseholdShareCents != 50000 || *cs[1].CategoryID != f.varios || cs[1].HouseholdShareCents != 50000 {
		t.Fatalf("%+v", cs)
	}
	res = c.Do("GET", "/api/v1/expenses/monthly-trend?family=external", nil)
	want(t, res, 200)
	res.JSON(t, &tr)
	if len(tr.Months) != 1 || len(tr.Months[0].Categories) != 1 || tr.Months[0].Categories[0].HouseholdShareCents != 50000 {
		t.Fatalf("familia external: %+v", tr)
	}
	want(t, c.Do("GET", "/api/v1/expenses/monthly-trend?family=otra", nil), 400)
	want(t, c.Do("GET", "/api/v1/expenses/monthly-trend?person_id=x", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/expenses/monthly-trend", nil), 401)
}

func TestOp_getDashboardMonthVsPlan(t *testing.T) {
	s, c, f := inicioFx(t)
	var mv api.DashboardMonthVsPlan
	res := c.Do("GET", "/api/v1/dashboard/month-vs-plan", nil)
	want(t, res, 200)
	res.JSON(t, &mv)
	// Sin plan: used_pct y status ausentes; 12-oct contando hoy → 20 días.
	if mv.Month != "2026-10" || mv.ActualCents != 100000 || mv.PlannedCents != nil || mv.UsedPct != nil || mv.Status != nil || mv.DaysLeft != 20 {
		t.Fatalf("sin plan: %+v", mv)
	}
	ins(t, s.DB, `INSERT INTO budget (category_id, amount_cents, valid_from_date) VALUES (?, 120000, '2026-10-01')`, f.despensa)
	res = c.Do("GET", "/api/v1/dashboard/month-vs-plan", nil)
	want(t, res, 200)
	res.JSON(t, &mv)
	// 1,000.00 ÷ 1,200.00 = 83.3 % → 83, warning (≥ 80); restante 200.00.
	if mv.PlannedCents == nil || *mv.PlannedCents != 120000 || *mv.UsedPct != 83 || *mv.Status != "warning" || *mv.RemainingCents != 20000 {
		t.Fatalf("con plan: %+v", mv)
	}
	res = c.Do("GET", "/api/v1/dashboard/month-vs-plan?period=2026-09", nil)
	want(t, res, 200)
	res.JSON(t, &mv)
	if mv.Month != "2026-09" || mv.DaysLeft != 0 || mv.ActualCents != 0 {
		t.Fatalf("mes pasado: %+v", mv)
	}
	want(t, c.Do("GET", "/api/v1/dashboard/month-vs-plan?period=octubre", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/dashboard/month-vs-plan", nil), 401)
}

func TestOp_getDashboardFortnightFlow(t *testing.T) {
	s, c, _ := inicioFx(t)
	var ff api.DashboardFortnightFlow
	res := c.Do("GET", "/api/v1/dashboard/fortnight-flow", nil)
	want(t, res, 200)
	res.JSON(t, &ff)
	if ff.StartDate != "2026-10-01" || ff.EndDate != "2026-10-15" || ff.SpentCents <= 0 ||
		ff.ClosingBalanceCents != ff.IncomeCents-ff.SpentCents-ff.CommittedCents {
		t.Fatalf("%+v", ff)
	}
	res = c.Do("GET", "/api/v1/dashboard/fortnight-flow?period=2026-10-q2", nil)
	want(t, res, 200)
	res.JSON(t, &ff)
	if ff.StartDate != "2026-10-16" || ff.EndDate != "2026-10-31" {
		t.Fatalf("q2: %+v", ff)
	}
	want(t, c.Do("GET", "/api/v1/dashboard/fortnight-flow?period=2026-13", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/dashboard/fortnight-flow", nil), 401)
}

func TestOp_getDashboardSavings(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	c := s.SessionFor(t, testutil.XimenaUser)
	var sv api.DashboardSavings
	res := c.Do("GET", "/api/v1/dashboard/savings", nil)
	want(t, res, 200)
	res.JSON(t, &sv)
	if sv.SavedTotalCents != 0 || sv.NextGoal != nil {
		t.Fatalf("vacío: %+v", sv)
	}
	ins(t, s.DB, `INSERT INTO savings_goal (name, kind, start_date, planned_monthly_cents) VALUES ('Fondo', 'fund', '2026-10-01', 50000)`)
	far := ins(t, s.DB, `INSERT INTO savings_goal (name, kind, target_cents, target_date, planned_monthly_cents, start_date) VALUES ('Viaje', 'goal', 2000000, '2027-12-31', 100000, '2026-10-01')`)
	near := ins(t, s.DB, `INSERT INTO savings_goal (name, kind, target_cents, target_date, planned_monthly_cents, start_date) VALUES ('Lentes', 'goal', 1000000, '2027-06-30', 100000, '2026-10-01')`)
	ins(t, s.DB, `INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (?, '2026-10-05', 'deposit', 150000)`, near)
	ins(t, s.DB, `INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (?, '2026-09-05', 'deposit', 30000)`, far)
	res = c.Do("GET", "/api/v1/dashboard/savings", nil)
	want(t, res, 200)
	res.JSON(t, &sv)
	if sv.SavedTotalCents != 180000 || sv.ContributedMonthCents != 150000 || sv.NextGoal == nil || sv.NextGoal.ID != near ||
		*sv.NextGoal.SavedCents != 150000 || *sv.NextGoal.ProgressPct != 15 {
		t.Fatalf("%+v %+v", sv, sv.NextGoal)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/dashboard/savings", nil), 401)
}

func TestOp_getPersonParticipation(t *testing.T) {
	s, c, _ := inicioFx(t)
	k := s.Fx.KarenPersonID
	cr := ins(t, s.DB, `INSERT INTO credit (name, kind, holder_person_id, credit_limit_cents, payment_frequency) VALUES ('Línea de prueba', 'revolving_line', ?, 1000000, 'fortnightly')`, s.Fx.AxelPersonID)
	pu := ins(t, s.DB, `INSERT INTO credit_purchase (credit_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents, term_count, payment_frequency, first_payment_date, interest_source)
 VALUES (?, 'installment_plan', 'Tenis', '2026-10-01', 90000, 10000, 100000, 2, 'fortnightly', '2026-10-15', 'quoted')`, cr)
	ins(t, s.DB, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 90000, 10000, 100000)`, pu, k)
	sp1 := ins(t, s.DB, `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, 1, '2026-10-15', 50000, 45000, 5000)`, pu)
	sp2 := ins(t, s.DB, `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, 2, '2026-10-30', 50000, 45000, 5000)`, pu)
	ins(t, s.DB, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 50000)`, sp1, k)
	ins(t, s.DB, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 50000)`, sp2, k)
	ins(t, s.DB, `INSERT INTO credit_payment (credit_id, scheduled_payment_id, payer_person_id, paid_date, amount_cents) VALUES (?, ?, ?, '2026-10-10', 50000)`, cr, sp1, k)

	var pp api.PersonParticipation
	res := c.Do("GET", "/api/v1/persons/"+itoa(k)+"/participation", nil)
	want(t, res, 200)
	res.JSON(t, &pp)
	if pp.PersonID != k || pp.Month != "2026-10" || pp.AccruedGroceriesCents != 50000 || pp.LiveDebtCents != 50000 || len(pp.PaymentShares) != 2 || len(pp.Purchases) != 1 {
		t.Fatalf("%+v", pp)
	}
	if pp.PaymentShares[0].Status != "paid" || pp.PaymentShares[1].Status != "pending" {
		t.Fatalf("estados: %+v", pp.PaymentShares)
	}
	if p := pp.Purchases[0]; p.CreditPurchaseID != pu || *p.Description != "Tenis" || p.PaymentsMade != 1 || p.PaymentsTotal != 2 || p.BalanceCents != 50000 || p.TotalCents != 100000 {
		t.Fatalf("compra: %+v", p)
	}
	want(t, c.Do("GET", "/api/v1/persons/"+itoa(s.Fx.AxelPersonID)+"/participation", nil), 404)
	want(t, c.Do("GET", "/api/v1/persons/9999/participation", nil), 404)
	if c.Raw("GET", "/api/v1/persons/abc/participation", nil).StatusCode != 400 {
		t.Error("person_id no numérico debe ser 400")
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/persons/"+itoa(k)+"/participation", nil), 401)
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }

func TestOp_search(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.XimenaUser)
	var sr api.SearchResult
	get := func(q string, code int) {
		t.Helper()
		res := c.Do("GET", "/api/v1/search?"+q, nil)
		want(t, res, code)
		if code == 200 {
			sr = api.SearchResult{}
			res.JSON(t, &sr)
		}
	}
	get("q=", 200)
	if len(sr.Groups) != 0 {
		t.Fatalf("texto vacío: %+v", sr)
	}
	for i := range 6 {
		newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 1000 + i, "expense_date": fmt.Sprintf("2026-10-0%d", i+1),
			"payment_method": "cash", "description": "Café de olla", "store_id": f.tienda}, 600+i)
	}
	// Sin mayúsculas ni acentos; máximo 5 por grupo con el total; fecha descendente.
	get("q=CAFE", 200)
	if len(sr.Groups) != 1 || sr.Groups[0].Type != "expense" || sr.Groups[0].Total != 6 || len(sr.Groups[0].Items) != 5 {
		t.Fatalf("%+v", sr)
	}
	if it := sr.Groups[0].Items[0]; it.Title != "Café de olla" || *it.Subtitle != "Varios" || *it.Date != "2026-10-06" {
		t.Fatalf("%+v", it)
	}
	get("q=aurrera", 200)
	if len(sr.Groups) != 2 || sr.Groups[0].Type != "expense" || sr.Groups[1].Type != "store" || sr.Groups[1].Items[0].ID != itoa(f.tienda) {
		t.Fatalf("aurrera: %+v", sr)
	}
	get("q=aurrera&type=store", 200)
	if len(sr.Groups) != 1 || sr.Groups[0].Items[0].Title != "Bodega Aurrerá" {
		t.Fatalf("type=store: %+v", sr)
	}
	get("q=categorias&type=settings_section,expense", 200)
	if len(sr.Groups) != 1 || sr.Groups[0].Items[0].ID != "P-26.04" {
		t.Fatalf("ajustes: %+v", sr)
	}
	get("q=nada-que-coincida", 200)
	if len(sr.Groups) != 0 {
		t.Fatalf("sin coincidencias: %+v", sr)
	}
	get("q=cafe&type=otro", 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/search?q=cafe", nil), 401)
}
