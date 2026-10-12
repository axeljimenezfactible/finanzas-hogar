package dominio_test

import (
	"context"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func TestGastosCalendarRepos(t *testing.T) {
	db, fx := testutil.NewDB(t)
	ctx := context.Background()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	exec(`INSERT INTO category (id, name) VALUES (90, 'Renta')`)
	exec(`INSERT INTO expense (id, kind, category_id, amount_cents, expense_date, payment_method) VALUES (1, 'fixed', 90, 1000, '2026-10-01', 'cash')`)
	exec(`INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (1, ?, 600), (1, ?, 400)`, fx.AxelPersonID, fx.KarenPersonID)
	exec(`INSERT INTO expense (id, kind, category_id, amount_cents, expense_date, payment_method) VALUES (2, 'variable', 90, 500, '2026-10-02', 'cash')`)
	exec(`INSERT INTO income (person_id, kind, description, amount_cents, income_date) VALUES (?, 'extra', 'Bono', 300, '2026-10-03')`, fx.AxelPersonID)
	exec(`INSERT INTO income (person_id, kind, description, amount_cents, income_date) VALUES (?, 'extra', 'Fuera', 300, '2026-11-03')`, fx.AxelPersonID)
	exec(`INSERT INTO savings_goal (id, name, kind, target_cents, target_date, start_date) VALUES (1, 'Fondo', 'goal', 5000, '2026-10-30', '2026-01-01')`)
	exec(`INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (1, '2026-10-04', 'withdrawal', 100)`)
	exec(`INSERT INTO store (kind, name, visit_weekday) VALUES ('market', 'Tianguis', 3)`)
	rows, err := dominio.GastosCalendarRows(ctx, db, "2026-10-01", "2026-10-31")
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]dominio.CalendarRow{}
	for _, r := range rows {
		got[r.Type] = r
	}
	if len(rows) != 4 || got["fixed_expense"].Title != "Renta" || *got["fixed_expense"].ExternalShareCents != 400 || *got["fixed_expense"].HouseholdShareCents != 600 ||
		got["income"].Title != "Bono" || got["savings_movement"].Title != "Retiro de Fondo" || *got["savings_goal"].AmountCents != 5000 {
		t.Fatalf("%+v", rows)
	}
	ms, err := dominio.GastosMarkets(ctx, db)
	if err != nil || len(ms) != 1 || ms[0].VisitWeekday != 3 {
		t.Fatalf("%+v %v", ms, err)
	}
	if id, err := dominio.GastosNotificationID(ctx, db, "no_existe"); err != nil || id != 0 {
		t.Fatalf("%d %v", id, err)
	}
}
