package gastos_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"finanzas-hogar/internal/api"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/testutil"
)

func scalar(t *testing.T, s *testutil.Server, q string, args ...any) int64 {
	t.Helper()
	var n int64
	if err := s.DB.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return n
}

// N-074: gasto con renglones, precios ticket_ocr, existencia por unidad con su movimiento y, a crédito, la compra.
func TestOp_confirmTicketScan(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	aceite := ins(t, s.DB, `INSERT INTO product (name, base_unit) VALUES ('Aceite de soya', 'ml')`)
	leche := ins(t, s.DB, `INSERT INTO product (name, base_unit) VALUES ('Leche entera', 'ml')`)
	jitomate := ins(t, s.DB, `INSERT INTO product (name, base_unit) VALUES ('Jitomate', 'g')`)
	pAceite := ins(t, s.DB, `INSERT INTO presentation (product_id, name, size_milli) VALUES (?, 'Aceite 946 ml', 946000)`, aceite)
	pLeche := ins(t, s.DB, `INSERT INTO presentation (product_id, name, size_milli) VALUES (?, 'Leche 1 L', 1000000)`, leche)
	pJitomate := ins(t, s.DB, `INSERT INTO presentation (product_id, name, size_milli, is_bulk) VALUES (?, 'Jitomate a granel', 1000000, 1)`, jitomate)
	invAceite := ins(t, s.DB, `INSERT INTO inventory_item (product_id, tracking_mode, quantity_milli) VALUES (?, 'stock', 0)`, aceite)
	invJitomate := ins(t, s.DB, `INSERT INTO inventory_item (product_id, tracking_mode, quantity_milli) VALUES (?, 'stock', 500000)`, jitomate)
	// Otro ticket del mismo día ya dejó el precio de la leche en esa tienda: la llave UNIQUE de price lo actualiza.
	oldPrice := ins(t, s.DB, `INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source) VALUES (?, ?, 2500, '2026-10-10', 'ticket_ocr')`, pLeche, f.tienda)
	id := ticket(t, s, 1, "done")
	if _, err := s.DB.Exec(`UPDATE attachment SET ocr_draft_json = '{"lines":[]}', ocr_draft_saved_at = '2026-10-11T10:00:00Z' WHERE id = ?`, id); err != nil {
		t.Fatal(err)
	}
	u := fmt.Sprintf("/api/v1/ticket-scans/%d/confirm", id)
	items := []map[string]any{
		{"line_no": 1, "description": "ACEITE SOYA 946ML", "presentation_id": pAceite, "quantity_milli": 3000, "unit_price_cents": 4500, "line_total_cents": 13500},
		{"line_no": 2, "description": "LECHE ENTERA 1L", "presentation_id": pLeche, "quantity_milli": 2000, "line_total_cents": 5800},
		{"line_no": 3, "description": "JITOMATE KG", "presentation_id": pJitomate, "quantity_milli": 1250, "line_total_cents": 3125, "cost_share_rule_id": f.karen},
		{"line_no": 4, "description": "LECHE ENTERA 1L", "presentation_id": pLeche, "quantity_milli": 1000, "line_total_cents": 2900},
		{"line_no": 5, "description": "BOLSA", "quantity_milli": 1000, "line_total_cents": 100},
	}
	expense := map[string]any{"category_id": f.despensa, "amount_cents": 25425, "expense_date": "2026-10-10", "payment_method": "cash",
		"store_id": f.tienda, "items": items}

	// 4xx antes de crear nada.
	want(t, s.Anonymous(t).Do("POST", u, map[string]any{"expense": expense}, key(9000)...), 401)
	want(t, c.Raw("POST", u, map[string]any{"expense": expense}), 400) // sin Idempotency-Key
	want(t, c.Raw("POST", u, map[string]any{"expense": expense, "extra": 1}, key(9001)...), 400)
	want(t, c.Do("POST", "/api/v1/ticket-scans/999/confirm", map[string]any{"expense": expense}, key(9002)...), 404)
	noStore := map[string]any{"category_id": f.despensa, "amount_cents": 25425, "expense_date": "2026-10-10", "payment_method": "cash", "items": items}
	want(t, c.Do("POST", u, map[string]any{"expense": noStore}, key(9003)...), 422)
	badSum := map[string]any{"category_id": f.despensa, "amount_cents": 99999, "expense_date": "2026-10-10", "payment_method": "cash", "store_id": f.tienda, "items": items}
	res := c.Do("POST", u, map[string]any{"expense": badSum}, key(9004)...)
	want(t, res, 422)
	if !strings.Contains(string(res.Body), `"field":"expense.items"`) {
		t.Fatalf("el campo debe llevar la ruta JSON: %s", res.Body)
	}
	ghost := map[string]any{"category_id": f.despensa, "amount_cents": 100, "expense_date": "2026-10-10", "payment_method": "cash", "store_id": f.tienda,
		"items": []map[string]any{{"line_no": 1, "description": "X", "presentation_id": 999, "quantity_milli": 1000, "line_total_cents": 100}}}
	want(t, c.Do("POST", u, map[string]any{"expense": ghost}, key(9005)...), 422)
	cashWithCredit := map[string]any{"expense": expense, "credit_purchase": map[string]any{"credit_id": 1, "kind": "installment_plan",
		"purchase_date": "2026-10-10", "principal_cents": 25425, "term_count": 4}}
	want(t, c.Do("POST", u, cashWithCredit, key(9006)...), 422)
	if n := scalar(t, s, `SELECT count(*) FROM expense`); n != 0 {
		t.Fatalf("un 4xx no debe dejar gastos: %d", n)
	}

	var out api.TicketConfirmResult
	res = c.Do("POST", u, map[string]any{"expense": expense}, key(9010)...)
	want(t, res, 201)
	res.JSON(t, &out)
	if out.ExpenseID == 0 || len(out.PriceIDs) != 3 || out.CreditPurchaseID != nil || len(out.InventoryItemIDs) != 2 ||
		out.InventoryItemIDs[0] != invAceite || out.InventoryItemIDs[1] != invJitomate {
		t.Fatalf("%+v", out)
	}
	// RN-15: 3 botellas de 946 ml = 2,838 ml; 1.250 kg a granel sobre 500 g = 1,750 g. La leche no lleva existencia.
	if q := scalar(t, s, `SELECT quantity_milli FROM inventory_item WHERE id = ?`, invAceite); q != 2838000 {
		t.Errorf("aceite %d", q)
	}
	if q := scalar(t, s, `SELECT quantity_milli FROM inventory_item WHERE id = ?`, invJitomate); q != 1750000 {
		t.Errorf("jitomate %d", q)
	}
	if n := scalar(t, s, `SELECT count(*) FROM inventory_movement m JOIN expense_item ei ON ei.id = m.expense_item_id
 WHERE m.kind = 'purchase' AND ei.expense_id = ?`, out.ExpenseID); n != 2 {
		t.Errorf("movimientos %d", n)
	}
	// Precios: el capturado, importe ÷ presentaciones (5800 / 2 = 2900) y 3125 / 1.25 = 2500; la leche reusa el renglón del día.
	for pres, cents := range map[int64]int64{pAceite: 4500, pLeche: 2900, pJitomate: 2500} {
		if got := scalar(t, s, `SELECT price_cents FROM price WHERE presentation_id = ? AND source = 'ticket_ocr' AND attachment_id = ?`, pres, id); got != cents {
			t.Errorf("precio %d: %d", pres, got)
		}
	}
	if !slices.Contains(out.PriceIDs, oldPrice) {
		t.Errorf("el precio del día se actualiza, no se duplica: %v", out.PriceIDs)
	}
	var ts api.TicketScan
	c.Do("GET", fmt.Sprintf("/api/v1/ticket-scans/%d", id), nil).JSON(t, &ts)
	if ts.OCRStatus != "confirmed" || ts.ExpenseID == nil || *ts.ExpenseID != out.ExpenseID || ts.HasDraft {
		t.Fatalf("ticket %+v", ts)
	}
	var e api.Expense
	c.Do("GET", fmt.Sprintf("/api/v1/expenses/%d", out.ExpenseID), nil).JSON(t, &e)
	if e.AmountCents != 25425 || shareOf(e, s.Fx.KarenPersonID) == 0 {
		t.Fatalf("gasto %+v", e)
	}
	if n := scalar(t, s, `SELECT count(*) FROM audit_log WHERE table_name IN ('expense', 'price', 'inventory_item', 'inventory_movement', 'attachment')`); n < 8 {
		t.Errorf("audit_log %d", n)
	}
	want(t, c.Do("POST", u, map[string]any{"expense": expense}, key(9011)...), 409)

	// A crédito: la despensa 1 con Kueski (canon S3: 309450 + 80050 = 389500 en 4 pagos de 97375).
	kueski := ins(t, s.DB, `INSERT INTO credit (name, kind, holder_person_id, credit_limit_cents, payment_frequency, first_payment_min_days)
 VALUES ('Kueski Pay', 'revolving_line', ?, 1000000, 'fortnightly', 15)`, s.Fx.AxelPersonID)
	id2 := ticket(t, s, 2, "failed")
	u2 := fmt.Sprintf("/api/v1/ticket-scans/%d/confirm", id2)
	credExp := map[string]any{"category_id": f.despensa, "amount_cents": 309450, "expense_date": "2026-10-10", "payment_method": "credit",
		"items": []map[string]any{{"line_no": 1, "description": "DESPENSA", "quantity_milli": 1000, "line_total_cents": 309450}}}
	purchase := func(principal, credit int64) map[string]any {
		return map[string]any{"credit_id": credit, "kind": "installment_plan", "purchase_date": "2026-10-10", "principal_cents": principal,
			"interest_cents": 80050, "term_count": 4, "shares": []map[string]any{{"person_id": s.Fx.AxelPersonID, "principal_cents": principal - 100000},
				{"person_id": s.Fx.KarenPersonID, "principal_cents": 100000}}}
	}
	want(t, c.Do("POST", u2, map[string]any{"expense": credExp}, key(9020)...), 422)
	want(t, c.Do("POST", u2, map[string]any{"expense": credExp, "credit_purchase": purchase(300000, kueski)}, key(9021)...), 422)
	want(t, c.Do("POST", u2, map[string]any{"expense": credExp, "credit_purchase": purchase(309450, 999)}, key(9022)...), 422)
	res = c.Do("POST", u2, map[string]any{"expense": credExp, "credit_purchase": purchase(309450, kueski)}, key(9023)...)
	want(t, res, 201)
	var out2 api.TicketConfirmResult
	res.JSON(t, &out2)
	if out2.CreditPurchaseID == nil || len(out2.PriceIDs) != 0 || len(out2.InventoryItemIDs) != 0 {
		t.Fatalf("%+v", out2)
	}
	if n, sum := scalar(t, s, `SELECT count(*) FROM scheduled_payment WHERE credit_purchase_id = ?`, *out2.CreditPurchaseID),
		scalar(t, s, `SELECT sum(amount_cents) FROM scheduled_payment WHERE credit_purchase_id = ?`, *out2.CreditPurchaseID); n != 4 || sum != 389500 {
		t.Fatalf("serie %d pagos, %d", n, sum)
	}
	if got := scalar(t, s, `SELECT expense_id FROM credit_purchase WHERE id = ?`, *out2.CreditPurchaseID); got != out2.ExpenseID {
		t.Fatalf("la compra va ligada al gasto del ticket: %d", got)
	}
	if got := scalar(t, s, `SELECT amount_cents FROM scheduled_payment WHERE credit_purchase_id = ? AND seq = 1`, *out2.CreditPurchaseID); got != 97375 {
		t.Fatalf("primer pago %d", got)
	}
}
