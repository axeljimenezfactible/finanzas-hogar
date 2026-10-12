package dominio_test

import (
	"context"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func TestGastosInicioRepos(t *testing.T) {
	db, fx := testutil.NewDB(t)
	ctx := context.Background()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	exec(`INSERT INTO credit (id, name, kind, holder_person_id, credit_limit_cents, payment_frequency) VALUES (1, 'Línea', 'revolving_line', ?, 100000, 'fortnightly')`, fx.AxelPersonID)
	exec(`INSERT INTO credit_purchase (id, credit_id, kind, purchase_date, principal_cents, total_cents, term_count, payment_frequency, first_payment_date, interest_source)
 VALUES (1, 1, 'installment_plan', '2026-10-01', 2000, 2000, 2, 'fortnightly', '2026-10-15', 'quoted')`)
	exec(`INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (1, ?, 2000, 0, 2000)`, fx.KarenPersonID)
	exec(`INSERT INTO scheduled_payment (id, credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (1, 1, 1, '2026-10-15', 1000, 1000, 0), (2, 1, 2, '2026-10-30', 1000, 1000, 0)`)
	exec(`INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (1, ?, 1000), (2, ?, 1000)`, fx.KarenPersonID, fx.KarenPersonID)
	exec(`INSERT INTO credit_payment (credit_id, scheduled_payment_id, payer_person_id, paid_date, amount_cents) VALUES (1, 1, ?, '2026-10-10', 400)`, fx.KarenPersonID)
	exec(`INSERT INTO savings_goal (name, kind, target_cents, target_date, start_date) VALUES ('Meta', 'goal', 5000, '2027-01-31', '2026-10-01')`)

	if err := dominio.GastosIsExternalPerson(ctx, db, fx.AxelPersonID); err == nil {
		t.Error("Axel no es externa")
	}
	if err := dominio.GastosIsExternalPerson(ctx, db, fx.KarenPersonID); err != nil {
		t.Error(err)
	}
	sh, err := dominio.GastosPersonShares(ctx, db, fx.KarenPersonID)
	if err != nil || len(sh) != 2 || sh[0].Paid() || sh[0].PendingCents() != 600 || sh[1].PendingCents() != 1000 {
		t.Fatalf("%v %+v", err, sh)
	}
	ps, err := dominio.GastosParticularPurchases(ctx, db, fx.KarenPersonID)
	if err != nil || len(ps) != 1 || ps[0].TotalCents != 2000 {
		t.Fatalf("%v %+v", err, ps)
	}
	if g, ok, err := dominio.GastosNextSavingsGoal(ctx, db); err != nil || !ok || g.Name != "Meta" {
		t.Fatalf("%v %v %+v", err, ok, g)
	}
	if ids, err := dominio.GastosExpenseIDs(ctx, db, dominio.ExpenseFilter{}); err != nil || len(ids) != 0 {
		t.Fatalf("%v %v", err, ids)
	}
	for typ, n := range map[string]int{"credit": 1, "credit_purchase": 1, "savings_goal": 1, "settings_section": 1, "store": 0} {
		hs, err := dominio.GastosSearch(ctx, db, typ, map[string]string{"settings_section": "categorías"}[typ]+map[bool]string{true: "", false: "LÍN"}[typ == "settings_section"])
		if typ == "credit_purchase" || typ == "savings_goal" {
			hs, err = dominio.GastosSearch(ctx, db, typ, map[string]string{"credit_purchase": "", "savings_goal": "meta"}[typ])
		}
		if err != nil || len(hs) != n {
			t.Errorf("%s: %v %+v", typ, err, hs)
		}
	}
	if _, err := dominio.GastosSearch(ctx, db, "otro", "x"); err == nil {
		t.Error("tipo desconocido")
	}
}
