package dominio_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func gxIns(t *testing.T, db *sql.DB, q string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(q+" RETURNING id", args...).Scan(&id); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return id
}

func gxSave(t *testing.T, db *sql.DB, in dominio.ExpenseIn) int64 {
	t.Helper()
	var id int64
	err := dominio.WithTx(context.Background(), db, func(tx *sql.Tx) error {
		var err error
		id, err = dominio.SaveExpense(context.Background(), tx, 0, in, true)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestGastosRepositorio(t *testing.T) {
	ctx := context.Background()
	db, fx := testutil.NewDB(t)
	cat := gxIns(t, db, `INSERT INTO category (name) VALUES ('Comida fuera')`)
	store := gxIns(t, db, `INSERT INTO store (kind, name) VALUES ('chain', 'Mercado Jamaica')`)
	half := int64(5000)
	fixed := int64(300)
	desc := "Taquería"
	key := "clave-de-prueba-0001"
	// Reparto propio del gasto y renglón «mixto»: crea dos reglas de un solo uso.
	id := gxSave(t, db, dominio.ExpenseIn{Kind: "variable", Description: &desc, CategoryID: cat, AmountCents: 1001, ExpenseDate: "2026-10-05", AccrualMonths: 1,
		PaymentMethod: "cash", StoreID: &store, PersonID: &fx.XimenaPersonID, IdempotencyKey: &key, ItemsSet: true,
		Parts: []dominio.ExpenseSharePart{{PersonID: fx.KarenPersonID, PctBP: &half}},
		Items: []dominio.ExpenseItemIn{{LineNo: 1, Description: "Tacos", QuantityMilli: 1000, LineTotalCents: 900,
			Parts: []dominio.ExpenseSharePart{{PersonID: fx.KarenPersonID, FixedCents: &fixed}}}}})
	x, err := dominio.GetExpense(ctx, db, id)
	if err != nil || x.AmountCents != 1001 || x.HouseholdCents+x.ExternalCents != 1001 || len(x.Shares) != 2 {
		t.Fatalf("%+v %v", x, err)
	}
	// Grupo del renglón 900: Karen 300 fijo; diferencia 101 con el reparto propio: Karen HALF_UP(50.5) = 51.
	if x.ExternalCents != 351 {
		t.Fatalf("reparto por grupos RN-03: %d", x.ExternalCents)
	}
	if got, ok, err := dominio.ExpenseIDByIdempotencyKey(ctx, db, key); !ok || got != id || err != nil {
		t.Fatal("K9")
	}
	if _, ok, _ := dominio.ExpenseIDByIdempotencyKey(ctx, db, "otra-clave-000000000"); ok {
		t.Fatal("clave inexistente")
	}
	no := false
	for _, f := range []dominio.ExpenseFilter{
		{Q: "TAQUERIA"}, {Q: "jamaica"}, {Q: "comida"}, {Family: "joint"}, {Family: "household"}, {Family: "external"},
		{Kinds: []string{"variable"}, PaymentMethods: []string{"cash"}}, {PersonIDs: []int64{fx.XimenaPersonID}, StoreIDs: []int64{store}, CategoryIDs: []int64{cat}},
		{HasAttachment: &no},
	} {
		xs, err := dominio.ListExpenses(ctx, db, f, "ORDER BY e.id", "LIMIT 10 OFFSET 0")
		if err != nil || len(xs) != 1 {
			t.Errorf("filtro %+v: %d %v", f, len(xs), err)
		}
	}
	if xs, _ := dominio.ListExpenses(ctx, db, dominio.ExpenseFilter{CreditIDs: []int64{1}}, "", ""); len(xs) != 0 {
		t.Error("filtro de crédito")
	}
	tot, err := dominio.SumExpenses(ctx, db, dominio.ExpenseFilter{})
	if err != nil || tot.Count != 1 || tot.ExternalCents != 351 {
		t.Fatalf("%+v %v", tot, err)
	}
	if c, cents, ok, err := dominio.TopHouseholdCategory(ctx, db, dominio.ExpenseFilter{}); !ok || c != cat || cents != 650 || err != nil {
		t.Fatalf("top %d %d %v", c, cents, err)
	}
	if _, _, ok, _ := dominio.TopHouseholdCategory(ctx, db, dominio.ExpenseFilter{Q: "nada"}); ok {
		t.Fatal("sin gastos no hay categoría principal")
	}
	in, err := dominio.ExpenseInOf(ctx, db, x)
	if err != nil || len(in.Items) != 1 || in.Items[0].CostShareRuleID == nil {
		t.Fatalf("%+v %v", in, err)
	}
	in.ItemsSet = true
	dup := gxSave(t, db, in)
	if err := dominio.CopyExpenseShares(ctx, db, id, dup); err != nil {
		t.Fatal(err)
	}
	if sum, maxLine, n, err := dominio.ExpenseItemsSum(ctx, db, dup); sum != 900 || maxLine != 1 || n != 1 || err != nil {
		t.Fatal("suma de renglones")
	}
	if its, err := dominio.ListExpenseItems(ctx, db, dup, "ORDER BY i.line_no", ""); err != nil || len(its) != 1 {
		t.Fatal("renglones")
	}
	var de *dominio.Error
	if err := dominio.DeleteExpense(ctx, db, dominio.ExpenseRow{ID: id, CreditPurchaseID: sql.NullInt64{Int64: 1, Valid: true}}); !errors.As(err, &de) || de.Code != "has_credit_purchase" {
		t.Fatalf("409 has_credit_purchase: %v", err)
	}
	if err := dominio.DeleteExpense(ctx, db, x); err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.GetExpense(ctx, db, id); !errors.As(err, &de) {
		t.Fatal("404 tras borrar")
	}
	// Errores de reparto y de partes.
	big := int64(2000)
	if _, err := dominio.ExpenseSplitByRule(1000, []dominio.ExpenseSharePart{{PersonID: fx.KarenPersonID, FixedCents: &big}}, fx.AxelPersonID); err == nil {
		t.Fatal("fijo mayor que el monto")
	}
	bad := int64(10001)
	neg := int64(-1)
	for _, ps := range [][]dominio.ExpenseSharePart{{{PersonID: 999, PctBP: &half}}, {{PersonID: fx.KarenPersonID, PctBP: &half}, {PersonID: fx.KarenPersonID, PctBP: &half}},
		{{PersonID: fx.KarenPersonID}}, {{PersonID: fx.KarenPersonID, PctBP: &bad}}, {{PersonID: fx.KarenPersonID, FixedCents: &neg}}} {
		if err := dominio.ExpenseCheckParts(ctx, db, ps, "cost_share_rule_parts"); err == nil {
			t.Errorf("partes inválidas aceptadas: %+v", ps)
		}
	}
	missing := int64(999)
	if _, err := dominio.ExpenseComputeShares(ctx, db, 100, &missing, cat, nil); err == nil {
		t.Fatal("regla inexistente")
	}
	if _, err := dominio.ExpenseComputeShares(ctx, db, 100, nil, cat, []dominio.ExpenseItemIn{{LineTotalCents: 200}}); err == nil {
		t.Fatal("renglones mayores que el importe")
	}
	if sh, err := dominio.ExpenseComputeShares(ctx, db, 0, nil, cat, nil); err != nil || len(sh) != 1 || sh[0].PersonID != fx.AxelPersonID {
		t.Fatalf("importe 0: %+v %v", sh, err)
	}
	if k, err := dominio.ExpensePlannedItemKind(ctx, db, nil); k != "" || err != nil {
		t.Fatal("sin partida")
	}
	if k, err := dominio.ExpensePlannedItemKind(ctx, db, &missing); k != "" || err != nil {
		t.Fatal("partida inexistente")
	}
	if dominio.GastosFoldText(" ÑANDÚ ") != "nandu" || dominio.GastosLikePattern("50%_") != `%50\%\_%` {
		t.Fatal("búsqueda sin acentos")
	}
}

func TestTicketsRepositorio(t *testing.T) {
	ctx := context.Background()
	db, _ := testutil.NewDB(t)
	id := gxIns(t, db, `INSERT INTO attachment (kind, file_path, mime_type, size_bytes, sha256, ocr_status, ocr_text) VALUES ('ticket_photo', 'a/1.jpg', 'image/jpeg', 10, ?, 'done', 'PAN 10.00')`, "a"+string(make([]byte, 0))+"0000000000000000000000000000000000000000000000000000000000000000"[:63])
	conf := gxIns(t, db, `INSERT INTO attachment (kind, file_path, mime_type, size_bytes, sha256, ocr_status) VALUES ('ticket_photo', 'a/2.jpg', 'image/jpeg', 10, ?, 'confirmed')`, "b"+"000000000000000000000000000000000000000000000000000000000000000")
	if ts, err := dominio.ListTicketScans(ctx, db, "ORDER BY a.id", ""); err != nil || len(ts) != 1 || ts[0].ID != id {
		t.Fatalf("cola: %+v %v", ts, err)
	}
	if err := dominio.SaveTicketScanDraft(ctx, db, id, `{"lines":[]}`, "2026-10-12T18:00:00.000Z"); err != nil {
		t.Fatal(err)
	}
	if tk, err := dominio.GetTicketScan(ctx, db, id); err != nil || !tk.DraftJSON.Valid {
		t.Fatal("borrador K3")
	}
	var de *dominio.Error
	if err := dominio.SaveTicketScanDraft(ctx, db, conf, `{"lines":[]}`, "2026-10-12T18:00:00.000Z"); !errors.As(err, &de) || de.Code != "conflict" {
		t.Fatalf("confirmado: %v", err)
	}
	if _, err := dominio.GetTicketScan(ctx, db, 999); !errors.As(err, &de) {
		t.Fatal("404")
	}
}
