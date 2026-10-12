package dominio_test

import (
	"context"
	"database/sql"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

// Repositorios de confirmTicketScan (N-074): presentaciones, renglones, existencia, precio del día y cierre del ticket.
func TestGastosConfirmTicketRepos(t *testing.T) {
	db, _ := testutil.NewDB(t)
	ctx := context.Background()
	id := func(q string, args ...any) int64 {
		t.Helper()
		var n int64
		if err := db.QueryRow(q+" RETURNING id", args...).Scan(&n); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return n
	}
	prod := id(`INSERT INTO product (name, base_unit) VALUES ('Aceite', 'ml')`)
	pres := id(`INSERT INTO presentation (product_id, name, size_milli) VALUES (?, 'Aceite 946 ml', 946000)`, prod)
	store := id(`INSERT INTO store (kind, name) VALUES ('chain', 'Tienda')`)
	cat := id(`INSERT INTO category (name) VALUES ('Despensa')`)
	exp := id(`INSERT INTO expense (kind, category_id, amount_cents, expense_date, accrual_months, payment_method) VALUES ('variable', ?, 100, '2026-10-10', 1, 'cash')`, cat)
	line := id(`INSERT INTO expense_item (expense_id, line_no, description, presentation_id, quantity_milli, line_total_cents) VALUES (?, 7, 'ACEITE', ?, 1000, 100)`, exp, pres)
	att := id(`INSERT INTO attachment (kind, file_path, mime_type, size_bytes, sha256, ocr_status, ocr_draft_json, ocr_draft_saved_at)
 VALUES ('ticket_photo', 'adjuntos/r.jpg', 'image/jpeg', 1, ?, 'done', '{}', '2026-10-11T10:00:00Z')`, "0000000000000000000000000000000000000000000000000000000000000001")

	ps, err := dominio.GastosTicketPresentations(ctx, db, []int64{pres, 999})
	if err != nil || len(ps) != 1 || ps[pres].ProductID != prod || ps[pres].SizeMilli != 946000 || ps[pres].ProductName != "Aceite" {
		t.Fatalf("%v %+v", err, ps)
	}
	if empty, err := dominio.GastosTicketPresentations(ctx, db, nil); err != nil || len(empty) != 0 {
		t.Fatalf("%v %+v", err, empty)
	}
	lines, err := dominio.GastosExpenseItemIDs(ctx, db, exp)
	if err != nil || lines[7] != line {
		t.Fatalf("%v %+v", err, lines)
	}
	if _, ok, err := dominio.GastosStockItem(ctx, db, prod); ok || err != nil {
		t.Fatalf("sin inventario: %v %v", ok, err)
	}
	inv := id(`INSERT INTO inventory_item (product_id, tracking_mode, quantity_milli) VALUES (?, 'stock', 1000)`, prod)
	got, ok, err := dominio.GastosStockItem(ctx, db, prod)
	if !ok || err != nil || got != inv {
		t.Fatalf("%d %v %v", got, ok, err)
	}
	now := "2026-10-12T18:00:00.000Z"
	err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if err := dominio.GastosTicketStockIn(ctx, tx, inv, "2026-10-10", 946000, line, now, nil, "req"); err != nil {
			return err
		}
		p := dominio.TicketPrice{PresentationID: pres, StoreID: store, PriceCents: 4500, ExpenseItemID: line, AttachmentID: att, ObservedDate: "2026-10-10"}
		first, err := dominio.GastosSaveTicketPrice(ctx, tx, p, now, nil, "req")
		if err != nil {
			return err
		}
		p.PriceCents = 4400
		second, err := dominio.GastosSaveTicketPrice(ctx, tx, p, now, nil, "req")
		if err != nil || second != first {
			t.Errorf("mismo día y tienda: se actualiza (%d, %d, %v)", first, second, err)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	var qty, delta, cents, audits int64
	if err := db.QueryRow(`SELECT i.quantity_milli, m.quantity_delta_milli, (SELECT price_cents FROM price), (SELECT count(*) FROM audit_log)
 FROM inventory_item i JOIN inventory_movement m ON m.inventory_item_id = i.id AND m.expense_item_id = ? WHERE i.id = ?`, line, inv).Scan(&qty, &delta, &cents, &audits); err != nil {
		t.Fatal(err)
	}
	if qty != 947000 || delta != 946000 || cents != 4400 || audits != 4 {
		t.Fatalf("existencia %d, movimiento %d, precio %d, bitácora %d", qty, delta, cents, audits)
	}
	if ok, err := dominio.GastosConfirmTicket(ctx, db, att, exp, now); !ok || err != nil {
		t.Fatalf("confirmar: %v %v", ok, err)
	}
	tk, err := dominio.GetTicketScan(ctx, db, att)
	if err != nil || tk.OCRStatus != "confirmed" || tk.ExpenseID.Int64 != exp || tk.DraftJSON.Valid || tk.DraftSavedAt.Valid {
		t.Fatalf("%v %+v", err, tk)
	}
	if ok, err := dominio.GastosConfirmTicket(ctx, db, att, exp, now); ok || err != nil {
		t.Fatalf("una segunda confirmación no cambia nada: %v %v", ok, err)
	}
}
