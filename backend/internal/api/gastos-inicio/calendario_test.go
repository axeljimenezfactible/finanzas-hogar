package gastos_test

import (
	"net/url"
	"testing"

	"finanzas-hogar/internal/api"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/testutil"
)

// credFx: línea revolvente con una compra de 1,000.00 en 2 pagos de 500.00 (10-oct vencido y 30-oct);
// cada pago es 300.00 del hogar (Axel) y 200.00 de Karen.
func credFx(t *testing.T, s *testutil.Server) (cr int64, sp []int64) {
	t.Helper()
	cr = ins(t, s.DB, `INSERT INTO credit (name, kind, holder_person_id, credit_limit_cents, payment_frequency) VALUES ('Línea de prueba', 'revolving_line', ?, 1000000, 'fortnightly')`, s.Fx.AxelPersonID)
	pu := ins(t, s.DB, `INSERT INTO credit_purchase (credit_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents, term_count, payment_frequency, first_payment_date, interest_source)
 VALUES (?, 'installment_plan', 'Despensa', '2026-09-28', 90000, 10000, 100000, 2, 'fortnightly', '2026-10-10', 'quoted')`, cr)
	ins(t, s.DB, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 54000, 6000, 60000)`, pu, s.Fx.AxelPersonID)
	ins(t, s.DB, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 36000, 4000, 40000)`, pu, s.Fx.KarenPersonID)
	for i, d := range []string{"2026-10-10", "2026-10-30"} {
		id := ins(t, s.DB, `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, ?, ?, 50000, 45000, 5000)`, pu, i+1, d)
		ins(t, s.DB, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 30000)`, id, s.Fx.AxelPersonID)
		ins(t, s.DB, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 20000)`, id, s.Fx.KarenPersonID)
		sp = append(sp, id)
	}
	return cr, sp
}

func TestOp_getDashboardNextPayment(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	c := s.SessionFor(t, testutil.XimenaUser)
	var np api.DashboardNextPayment
	res := c.Do("GET", "/api/v1/dashboard/next-payment", nil)
	want(t, res, 200)
	res.JSON(t, &np)
	if np.NextPayment != nil || len(np.Availability) != 0 {
		t.Fatalf("sin créditos: %+v", np)
	}
	cr, sp := credFx(t, s)
	res = c.Do("GET", "/api/v1/dashboard/next-payment", nil)
	want(t, res, 200)
	res.JSON(t, &np)
	n := np.NextPayment
	if n == nil || n.DueDate != "2026-10-10" || !n.IsOverdue || n.DaysLeft != -2 || n.AmountCents != 50000 || n.HouseholdShareCents != 30000 || n.ExternalShareCents != 20000 || n.PaymentCount != 1 || n.CreditID != cr {
		t.Fatalf("vencido: %+v", n)
	}
	if len(np.Availability) != 1 || np.Availability[0].CreditID != cr || np.Availability[0].AvailabilityReading != "B" || np.Availability[0].AvailabilityConfirmed {
		t.Fatalf("disponible: %+v", np.Availability)
	}
	// Karen pagó su parte del 10-oct y Axel la suya: el próximo es el 30-oct.
	for _, p := range []struct{ who, amt int64 }{{s.Fx.KarenPersonID, 20000}, {s.Fx.AxelPersonID, 30000}} {
		ins(t, s.DB, `INSERT INTO credit_payment (credit_id, scheduled_payment_id, payer_person_id, paid_date, amount_cents) VALUES (?, ?, ?, '2026-10-11', ?)`, cr, sp[0], p.who, p.amt)
	}
	res = c.Do("GET", "/api/v1/dashboard/next-payment", nil)
	want(t, res, 200)
	res.JSON(t, &np)
	if n = np.NextPayment; n == nil || n.DueDate != "2026-10-30" || n.IsOverdue || n.DaysLeft != 18 || n.AmountCents != 50000 {
		t.Fatalf("siguiente: %+v", n)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/dashboard/next-payment", nil), 401)
}

func TestOp_listCalendarEvents(t *testing.T) {
	s, c, f := inicioFx(t)
	cr, _ := credFx(t, s)
	newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 80000, "expense_date": "2026-10-20", "kind": "upcoming", "payment_method": "cash", "accrual_months": 1}, 601)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 120000, "expense_date": "2026-10-01", "kind": "fixed", "description": "Renta", "payment_method": "transfer", "accrual_months": 1}, 602)
	ins(t, s.DB, `INSERT INTO income (person_id, kind, description, amount_cents, income_date) VALUES (?, 'fixed', 'Sueldo', 2500000, '2026-10-15')`, s.Fx.AxelPersonID)
	ins(t, s.DB, `INSERT INTO store (kind, name, visit_weekday) VALUES ('market', 'Tianguis', 6)`)
	g := ins(t, s.DB, `INSERT INTO savings_goal (name, kind, target_cents, target_date, start_date) VALUES ('Viaje', 'goal', 1000000, '2026-10-31', '2026-01-01')`)
	ins(t, s.DB, `INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (?, '2026-10-05', 'deposit', 50000)`, g)
	ins(t, s.DB, `INSERT INTO credit_statement (credit_id, statement_date, due_date, balance_cents, source) VALUES (?, '2026-10-25', '2026-11-10', 100000, 'bank')`, cr)

	var pg api.CalendarEventPage
	res := c.Do("GET", "/api/v1/calendar/events", nil) // por defecto: oct-2026
	want(t, res, 200)
	res.JSON(t, &pg)
	got := map[string]int{}
	for _, e := range pg.Items {
		got[e.EventType]++
	}
	exp := map[string]int{"income": 1, "credit_payment": 2, "fixed_expense": 1, "upcoming_expense": 1, "card_statement": 1, "market_visit": 5, "savings_goal": 1, "savings_movement": 1, "monthly_report": 1}
	for k, v := range exp {
		if got[k] != v {
			t.Errorf("%s: %d eventos, se esperaban %d (%v)", k, got[k], v, got)
		}
	}
	// Orden: fecha y, dentro del día, el de N-029 (el 1-oct: gasto fijo antes que el reporte mensual).
	if e := pg.Items[0]; e.Date != "2026-10-01" || e.EventType != "fixed_expense" || *e.Status != "captured" || e.Title != "Renta" {
		t.Fatalf("primero: %+v", e)
	}
	if e := pg.Items[1]; e.EventType != "monthly_report" || e.RefTable != "notification" {
		t.Fatalf("segundo: %+v", e)
	}
	for _, e := range pg.Items {
		if e.EventType == "credit_payment" && e.Date == "2026-10-10" && (*e.Status != "overdue" || *e.ExternalShareCents != 20000 || *e.CreditID != cr || e.Title != "Pago de Línea de prueba") {
			t.Fatalf("pago vencido: %+v", e)
		}
	}
	res = c.Do("GET", "/api/v1/calendar/events?family=external&event_type=credit_payment,upcoming_expense", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 3 { // 2 pagos con parte de Karen + la despensa venidera (Conjunto 50/50)
		t.Fatalf("family=external: %+v", pg.Items)
	}
	res = c.Do("GET", "/api/v1/calendar/events?credit_id="+itoa(cr)+"&limit=2", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 2 || !pg.NextCursor.Valid {
		t.Fatalf("credit_id + limit: %+v", pg)
	}
	res = c.Do("GET", "/api/v1/calendar/events?credit_id="+itoa(cr)+"&limit=2&cursor="+url.QueryEscape(pg.NextCursor.Value), nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].Date != "2026-10-30" || pg.NextCursor.Valid { // 10-oct pago, 25-oct corte, 30-oct pago
		t.Fatalf("segunda página: %+v", pg)
	}
	res = c.Do("GET", "/api/v1/calendar/events?period=2026-11&event_type=monthly_report", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].Date != "2026-11-01" {
		t.Fatalf("noviembre: %+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/calendar/events?event_type=birthday", nil), 400)
	want(t, c.Do("GET", "/api/v1/calendar/events?family=friends", nil), 400)
	want(t, c.Do("GET", "/api/v1/calendar/events?period=octubre", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/calendar/events", nil), 401)
}
