package gastos_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/testutil"
)

// fx son los catálogos mínimos: reglas Hogar (por defecto), Conjunto 50/50 y Karen; categorías Despensa (Conjunto) y Varios; una tienda.
type fx struct{ hogar, conjunto, karen, despensa, varios, tienda int64 }

func ins(t *testing.T, db *sql.DB, q string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(q+" RETURNING id", args...).Scan(&id); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return id
}

func seed(t *testing.T, s *testutil.Server) fx {
	var f fx
	f.hogar = ins(t, s.DB, `INSERT INTO cost_share_rule (name, remainder_person_id, is_default) VALUES ('Hogar', ?, 1)`, s.Fx.AxelPersonID)
	f.conjunto = ins(t, s.DB, `INSERT INTO cost_share_rule (name, remainder_person_id) VALUES ('Conjunto 50/50', ?)`, s.Fx.AxelPersonID)
	ins(t, s.DB, `INSERT INTO cost_share_rule_part (rule_id, person_id, pct_bp) VALUES (?, ?, 5000)`, f.conjunto, s.Fx.KarenPersonID)
	f.karen = ins(t, s.DB, `INSERT INTO cost_share_rule (name, remainder_person_id) VALUES ('Karen', ?)`, s.Fx.KarenPersonID)
	f.despensa = ins(t, s.DB, `INSERT INTO category (code, name, kind, default_cost_share_rule_id) VALUES ('groceries', 'Despensa', 'variable', ?)`, f.conjunto)
	f.varios = ins(t, s.DB, `INSERT INTO category (name) VALUES ('Varios')`)
	f.tienda = ins(t, s.DB, `INSERT INTO store (kind, name) VALUES ('chain', 'Bodega Aurrerá')`)
	return f
}

func key(n int) []string { return []string{"Idempotency-Key", fmt.Sprintf("gasto-prueba-%010d", n)} }

func newExpense(t *testing.T, c *testutil.Client, body map[string]any, n int) api.Expense {
	t.Helper()
	res := c.Do("POST", "/api/v1/expenses", body, key(n)...)
	if res.StatusCode != 201 {
		t.Fatalf("alta: %d %s", res.StatusCode, res.Body)
	}
	var e api.Expense
	res.JSON(t, &e)
	return e
}

func shareOf(e api.Expense, person int64) int64 {
	for _, s := range e.Shares {
		if s.PersonID == person {
			return s.AmountCents
		}
	}
	return 0
}

func want(t *testing.T, res *testutil.Response, code int) {
	t.Helper()
	if res.StatusCode != code {
		t.Fatalf("esperaba %d, llegó %d: %s", code, res.StatusCode, res.Body)
	}
}

// RN-03 con el ejemplo de la despensa 1: Axel 1,903.00 y Karen 1,191.50 sobre 3,094.50.
func TestOp_createExpense(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	body := map[string]any{"category_id": f.despensa, "amount_cents": 309450, "expense_date": "2026-10-10", "payment_method": "debit", "accrual_months": 2,
		"store_id": f.tienda, "captured_at": "2026-10-10T19:20:00.000Z", "items": []map[string]any{
			{"line_no": 1, "description": "Conjunto", "quantity_milli": 1000, "line_total_cents": 170200},
			{"line_no": 2, "description": "Karen", "quantity_milli": 1000, "line_total_cents": 34050, "cost_share_rule_id": f.karen},
			{"line_no": 3, "description": "Hogar", "quantity_milli": 1000, "line_total_cents": 105200, "cost_share_rule_id": f.hogar}}}
	e := newExpense(t, c, body, 1)
	if shareOf(e, s.Fx.AxelPersonID) != 190300 || shareOf(e, s.Fx.KarenPersonID) != 119150 || *e.HouseholdShareCents != 190300 || *e.ExternalShareCents != 119150 {
		t.Fatalf("reparto RN-03: %+v", e.Shares)
	}
	if e.Kind != "variable" || e.CapturedAt == nil || len(e.CashFlows) != 1 {
		t.Fatalf("tipo, captura o efectivo: %+v", e)
	}
	// K9: mismo cuerpo y clave → 200 con el mismo gasto; pasadas las 24 h también; otro cuerpo dentro de las 24 h → 409.
	res := c.Do("POST", "/api/v1/expenses", body, key(1)...)
	want(t, res, 200)
	if res.Header.Get("Idempotent-Replayed") != "true" {
		t.Fatal("falta Idempotent-Replayed")
	}
	want(t, c.Do("POST", "/api/v1/expenses", map[string]any{"category_id": f.varios, "amount_cents": 100, "expense_date": "2026-10-10", "payment_method": "cash"}, key(1)...), 409)
	s.Clock.Advance(25 * time.Hour)
	if _, err := s.DB.Exec(`DELETE FROM idempotency_key`); err != nil {
		t.Fatal(err)
	}
	res = c.Do("POST", "/api/v1/expenses", body, key(1)...)
	want(t, res, 200)
	var again api.Expense
	res.JSON(t, &again)
	if again.ID != e.ID || res.Header.Get("Idempotent-Replayed") != "true" {
		t.Fatalf("K9 pasadas 24 h: %d vs %d", again.ID, e.ID)
	}
	// Futuro → upcoming; reparto propio «Conjunto con %».
	e2 := newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 1001, "expense_date": "2026-11-01", "payment_method": "cash",
		"cost_share_rule_parts": []map[string]any{{"person_id": s.Fx.KarenPersonID, "pct_bp": 5000}}}, 2)
	if e2.Kind != "upcoming" || shareOf(e2, s.Fx.KarenPersonID) != 501 || shareOf(e2, s.Fx.AxelPersonID) != 500 {
		t.Fatalf("upcoming o HALF_UP: %+v", e2)
	}
	want(t, c.Do("POST", "/api/v1/expenses", map[string]any{"category_id": f.varios, "amount_cents": 100, "expense_date": "2026-10-10", "payment_method": "credit"}, key(3)...), 422)
	want(t, c.Do("POST", "/api/v1/expenses", map[string]any{"category_id": 999, "amount_cents": 0, "expense_date": "2026-10-10", "payment_method": "cash"}, key(4)...), 422)
	want(t, c.Do("POST", "/api/v1/expenses", map[string]any{"category_id": f.varios, "amount_cents": 100, "expense_date": "2026-02-30", "payment_method": "cash"}, key(5)...), 400)
	want(t, c.Do("POST", "/api/v1/expenses", map[string]any{"category_id": f.varios, "amount_cents": 100, "expense_date": "2026-10-10", "payment_method": "cash"}), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/expenses", nil), 401)
}

func TestOp_listExpenses(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.XimenaUser)
	newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 20000, "expense_date": "2026-10-05", "payment_method": "cash", "description": "Súper semanal"}, 1)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 5000, "expense_date": "2026-10-06", "payment_method": "debit", "store_id": f.tienda}, 2)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 7000, "expense_date": "2026-09-06", "payment_method": "cash"}, 3)
	var pg api.ExpensePage
	c.Do("GET", "/api/v1/expenses", nil).JSON(t, &pg)
	if len(pg.Items) != 2 || pg.Items[0].ExpenseDate != "2026-10-06" || pg.NextCursor.Valid {
		t.Fatalf("mes en curso, orden por fecha: %+v", pg)
	}
	for q, n := range map[string]int{"?family=external": 1, "?family=joint": 1, "?q=SUPER": 1, "?q=aurrera": 1, "?period=2026-09": 1, "?kind=variable&payment_method=debit": 1,
		"?has_attachment=true": 0, "?category_id=" + fmt.Sprint(f.varios) + "&sort=amount_cents&limit=1": 1} {
		var p api.ExpensePage
		c.Do("GET", "/api/v1/expenses"+q, nil).JSON(t, &p)
		if len(p.Items) != n {
			t.Errorf("%s: %d", q, len(p.Items))
		}
	}
	for _, q := range []string{"?limit=201", "?kind=otro", "?family=x", "?has_attachment=si", "?sort=password", "?person_id=a"} {
		want(t, c.Do("GET", "/api/v1/expenses"+q, nil), 400)
	}
}

func TestOp_getExpenseTotals(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 20001, "expense_date": "2026-10-05", "payment_method": "cash"}, 1)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 5000, "expense_date": "2026-10-06", "payment_method": "cash"}, 2)
	var tot api.ExpenseTotals
	c.Do("GET", "/api/v1/expenses/totals?period=2026-10", nil).JSON(t, &tot)
	if tot.ExpenseCount != 2 || tot.AmountCents != 25001 || tot.ExternalShareCents != 10001 || tot.HouseholdShareCents != 15000 {
		t.Fatalf("%+v", tot)
	}
	want(t, c.Do("GET", "/api/v1/expenses/totals?family=otra", nil), 400)
}

func TestOp_getExpenseSummary(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 97375, "expense_date": "2026-10-05", "payment_method": "cash"}, 1)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 10000, "expense_date": "2026-09-05", "payment_method": "cash"}, 2)
	var sm api.ExpenseSummary
	c.Do("GET", "/api/v1/expenses/summary", nil).JSON(t, &sm)
	// RN-27 (K39): 973.75 en 31 días → diario 31.41 y semanal 219.88.
	if sm.PeriodStart != "2026-10-01" || sm.HouseholdShareCents != 97375 || sm.PreviousHouseholdShareCents != 10000 || sm.DifferenceCents != 87375 ||
		*sm.DailyAverageCents != 3141 || *sm.WeeklyAverageCents != 21988 || *sm.MonthlyAverageCents != 97375 || *sm.TopCategoryID != f.varios {
		t.Fatalf("%+v", sm)
	}
	for _, p := range []string{"2026-10-q2", "2026", "2026-W42", "2026-10-01..2026-10-10"} {
		want(t, c.Do("GET", "/api/v1/expenses/summary?period="+p, nil), 200)
	}
	if res := c.Do("GET", "/api/v1/expenses/summary?period=2026-13", nil); res.StatusCode < 400 || res.StatusCode > 422 {
		t.Fatalf("periodo inválido: %d", res.StatusCode)
	}
}

func TestOp_getExpense(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	e := newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 100, "expense_date": "2026-10-05", "payment_method": "cash"}, 1)
	want(t, c.Do("GET", fmt.Sprintf("/api/v1/expenses/%d", e.ID), nil), 200)
	want(t, c.Do("GET", "/api/v1/expenses/999", nil), 404)
	want(t, c.Do("GET", "/api/v1/expenses/abc", nil), 400)
}

func TestOp_updateExpense(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	e := newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 1000, "expense_date": "2026-10-05", "payment_method": "cash", "description": "Tianguis"}, 1)
	u := fmt.Sprintf("/api/v1/expenses/%d", e.ID)
	var got api.Expense
	c.Do("PATCH", u, map[string]any{"amount_cents": 2001, "description": nil, "notes": "con nota"}).JSON(t, &got)
	if got.AmountCents != 2001 || got.Description != nil || shareOf(got, s.Fx.KarenPersonID) != 1001 || *got.Notes != "con nota" {
		t.Fatalf("recalcula partes y borra con null: %+v", got)
	}
	c.Do("PATCH", u, map[string]any{"cost_share_rule_id": f.hogar, "items": []map[string]any{{"line_no": 1, "description": "Todo", "quantity_milli": 1000, "line_total_cents": 2001}}}).JSON(t, &got)
	if shareOf(got, s.Fx.KarenPersonID) != 0 || *got.HouseholdShareCents != 2001 {
		t.Fatalf("regla Hogar: %+v", got.Shares)
	}
	want(t, c.Do("PATCH", u, map[string]any{"accrual_months": 30}), 422)
	want(t, c.Do("PATCH", u, map[string]any{"amount_cents": "20.01"}), 400)
	want(t, c.Do("PATCH", "/api/v1/expenses/999", map[string]any{"notes": "x"}), 404)
}

func TestOp_deleteExpense(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	e := newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 100, "expense_date": "2026-10-05", "payment_method": "cash"}, 1)
	want(t, c.Do("DELETE", fmt.Sprintf("/api/v1/expenses/%d", e.ID), nil), 204)
	want(t, c.Do("DELETE", fmt.Sprintf("/api/v1/expenses/%d", e.ID), nil), 404)
	var n int
	if err := s.DB.QueryRow(`SELECT count(*) FROM audit_log WHERE table_name = 'expense' AND action = 'delete'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("copia en audit_log: %d %v", n, err)
	}
}

func TestOp_listExpenseItems(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	e := newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 3000, "expense_date": "2026-10-05", "payment_method": "cash", "items": []map[string]any{
		{"line_no": 1, "description": "Leche", "quantity_milli": 2000, "unit_price_cents": 1000, "line_total_cents": 2000},
		{"line_no": 2, "description": "Pan", "quantity_milli": 1000, "line_total_cents": 1000}}}, 1)
	c.Do("PATCH", fmt.Sprintf("/api/v1/expenses/%d", e.ID), map[string]any{"amount_cents": 3500})
	var pg api.ExpenseItemPage
	c.Do("GET", fmt.Sprintf("/api/v1/expenses/%d/items", e.ID), nil).JSON(t, &pg)
	if len(pg.Items) != 3 || pg.Items[2].Description != "Ajuste del ticket" || pg.Items[2].LineTotalCents != 500 {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/expenses/999/items", nil), 404)
	want(t, c.Do("GET", fmt.Sprintf("/api/v1/expenses/%d/items?limit=0", e.ID), nil), 400)
}

func TestOp_duplicateExpense(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	e := newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 1001, "expense_date": "2026-10-01", "payment_method": "cash"}, 1)
	res := c.Do("POST", fmt.Sprintf("/api/v1/expenses/%d/duplicate", e.ID), nil, key(2)...)
	want(t, res, 201)
	var d api.Expense
	res.JSON(t, &d)
	if d.ID == e.ID || d.ExpenseDate != "2026-10-12" || shareOf(d, s.Fx.KarenPersonID) != shareOf(e, s.Fx.KarenPersonID) {
		t.Fatalf("%+v", d)
	}
	want(t, c.Do("POST", "/api/v1/expenses/999/duplicate", nil, key(3)...), 404)
	want(t, c.Do("POST", fmt.Sprintf("/api/v1/expenses/%d/duplicate", e.ID), nil), 400)
}

func TestOp_markExpensePaid(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	e := newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 5000, "expense_date": "2026-10-20", "payment_method": "transfer"}, 1)
	u := fmt.Sprintf("/api/v1/expenses/%d/mark-paid", e.ID)
	var got api.Expense
	res := c.Do("POST", u, map[string]any{"amount_cents": 5200, "expense_date": "2026-10-12"}, key(2)...)
	want(t, res, 200)
	res.JSON(t, &got)
	if got.Kind != "variable" || got.AmountCents != 5200 || *got.HouseholdShareCents != 5200 {
		t.Fatalf("%+v", got)
	}
	want(t, c.Do("POST", u, map[string]any{"amount_cents": 5200, "expense_date": "2026-10-12"}, key(3)...), 409)
	want(t, c.Do("POST", u, map[string]any{"amount_cents": 0, "expense_date": "2026-10-12"}, key(4)...), 422)
	want(t, c.Do("POST", "/api/v1/expenses/999/mark-paid", map[string]any{"amount_cents": 1, "expense_date": "2026-10-12"}, key(5)...), 404)
}

func ticket(t *testing.T, s *testutil.Server, n int, status string) int64 {
	return ins(t, s.DB, `INSERT INTO attachment (kind, file_path, original_name, mime_type, size_bytes, sha256, ocr_status, ocr_engine, ocr_text)
 VALUES ('ticket_photo', ?, 'ticket.jpg', 'image/jpeg', 1234, ?, ?, 'tesseract 5.5.3 spa', ?)`, fmt.Sprintf("adjuntos/t%d.jpg", n),
		fmt.Sprintf("%064d", n), status, "2 LECHE ENTERA 1L 58.00\nPAN BLANCO 35.50\nTOTAL 93.50")
}

func TestOp_listTicketScans(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	ticket(t, s, 1, "done")
	ticket(t, s, 2, "confirmed")
	var pg api.TicketScanPage
	c.Do("GET", "/api/v1/ticket-scans", nil).JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].HasDraft {
		t.Fatalf("%+v", pg)
	}
	want(t, c.Do("GET", "/api/v1/ticket-scans?limit=0", nil), 400)
}

func TestOp_getTicketScan(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	id := ticket(t, s, 1, "done")
	var ts api.TicketScan
	c.Do("GET", fmt.Sprintf("/api/v1/ticket-scans/%d", id), nil).JSON(t, &ts)
	if len(ts.Lines) != 2 || ts.Lines[1].LineTotalCents == nil || *ts.Lines[1].LineTotalCents != 3550 {
		t.Fatalf("%+v", ts.Lines)
	}
	want(t, c.Do("GET", "/api/v1/ticket-scans/999", nil), 404)
}

func TestOp_saveTicketScanDraft(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	id := ticket(t, s, 1, "done")
	u := fmt.Sprintf("/api/v1/ticket-scans/%d/draft", id)
	var ts api.TicketScan
	c.Do("PUT", u, map[string]any{"store_id": f.tienda, "total_cents": 9350, "expense_date": "2026-10-10",
		"lines": []map[string]any{{"line_no": 1, "quantity_milli": 3000, "line_total_cents": 8700}}}).JSON(t, &ts)
	if !ts.HasDraft || ts.Draft == nil || *ts.DisplayTotalCents != 9350 || *ts.DisplayStoreID != f.tienda || ts.OCRDraftSavedAt == nil {
		t.Fatalf("%+v", ts)
	}
	want(t, c.Do("PUT", u, map[string]any{"payment_method": "vales", "lines": []any{}}), 422)
	// lines es obligatorio en TicketScanDraft (hallazgo H3): sin él, 422 validation_failed con el campo, y el borrador no cambia.
	sinLineas := c.Do("PUT", u, map[string]any{"total_cents": 1})
	want(t, sinLineas, 422)
	if b := string(sinLineas.Body); !bytes.Contains(sinLineas.Body, []byte(`"validation_failed"`)) || !bytes.Contains(sinLineas.Body, []byte(`"lines"`)) {
		t.Fatalf("sin lines: %s", b)
	}
	c.Do("GET", fmt.Sprintf("/api/v1/ticket-scans/%d", id), nil).JSON(t, &ts)
	if ts.DisplayTotalCents == nil || *ts.DisplayTotalCents != 9350 {
		t.Fatalf("el borrador cambió con un cuerpo inválido: %+v", ts)
	}
	want(t, c.Do("PUT", fmt.Sprintf("/api/v1/ticket-scans/%d/draft", ticket(t, s, 2, "confirmed")), map[string]any{"lines": []any{}}), 409)
	want(t, c.Do("PUT", "/api/v1/ticket-scans/999/draft", map[string]any{"lines": []any{}}), 404)
}
