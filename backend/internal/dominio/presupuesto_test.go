package dominio_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
	"finanzas-hogar/internal/testutil"
)

func TestPresupuestoRepositorios(t *testing.T) {
	db, fx := testutil.NewDB(t)
	ctx := context.Background()
	var cat, rule, goal int64
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(db.QueryRow(`INSERT INTO category (code, name) VALUES ('groceries', 'Despensa') RETURNING id`).Scan(&cat))
	must(db.QueryRow(`INSERT INTO cost_share_rule (name, remainder_person_id, is_default) VALUES ('Mitad', ?, 1) RETURNING id`, fx.AxelPersonID).Scan(&rule))
	_, err := db.Exec(`INSERT INTO cost_share_rule_part (rule_id, person_id, pct_bp) VALUES (?, ?, 5000)`, rule, fx.KarenPersonID)
	must(err)
	_, err = db.Exec(`INSERT INTO planned_item (kind, description, amount_cents, frequency, interval_count, start_date, accrual_months, category_id)
		VALUES ('fixed_expense', 'Renta', 500000, 'monthly', 1, '2026-10-01', 1, ?)`, cat)
	must(err)
	_, err = db.Exec(`INSERT INTO planned_item (kind, description, amount_cents, frequency, interval_count, start_date, accrual_months, person_id)
		VALUES ('income', 'Sueldo', 1000000, 'fortnightly', 1, '2026-01-15', 1, ?)`, fx.AxelPersonID)
	must(err)
	must(db.QueryRow(`INSERT INTO savings_goal (name, kind, start_date, planned_monthly_cents) VALUES ('Fondo', 'fund', '2026-10-01', 1000) RETURNING id`).Scan(&goal))
	var exp int64
	if err := db.QueryRow(`INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, payment_method, paid_by_person_id)
		VALUES ('variable', 'Súper', ?, 10000, '2026-10-05', 'cash', ?) RETURNING id`, cat, fx.AxelPersonID).Scan(&exp); err == nil {
		if _, err := db.Exec(`INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, 10000)`, exp, fx.AxelPersonID); err != nil {
			t.Logf("sin parte del gasto: %v", err)
		}
	} else {
		t.Logf("sin gasto de prueba: %v", err)
	}

	actor := &fx.AxelUserID
	var b1 int64
	must(dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if _, err := dominio.BudgetAreaInsert(ctx, tx, "expense", nil, nil, actor, "r"); err == nil {
			t.Error("tabla fuera del área debe fallar")
		}
		var err error
		b1, err = dominio.BudgetAreaInsert(ctx, tx, "budget", []string{"category_id", "amount_cents", "valid_from_date"}, []any{cat, 1000, "2026-10-01"}, actor, "r")
		if err != nil {
			return err
		}
		if err := dominio.BudgetAreaUpdate(ctx, tx, "budget", b1, []string{"amount_cents"}, []any{2000}, "2026-10-12T18:00:00.000Z", actor, "r"); err != nil {
			return err
		}
		if err := dominio.BudgetAreaUpdate(ctx, tx, "person", 1, nil, nil, "", actor, "r"); err == nil {
			t.Error("update fuera del área debe fallar")
		}
		if err := dominio.BudgetAreaDelete(ctx, tx, "budget", 9999, actor, "r"); !errors.Is(err, dominio.ErrNotFound) && err == nil {
			t.Errorf("borrar inexistente: %v", err)
		}
		if err := dominio.BudgetAreaDelete(ctx, tx, "person", 1, actor, "r"); err == nil {
			t.Error("borrar fuera del área debe fallar")
		}
		return dominio.EnsurePeriods(ctx, tx, dominio.NewDate(2027, time.February, 1), dominio.NewDate(2027, time.March, 31), actor, "r")
	}))
	cn := sql.NullInt64{Int64: cat, Valid: true}
	if over, err := dominio.BudgetOverlaps(ctx, db, cn, "2026-12-01", sql.NullString{}, 0); err != nil || !over {
		t.Errorf("traslape: %v %v", over, err)
	}
	if a, err := dominio.BudgetAmountAt(ctx, db, cn, "2026-11-01"); err != nil || a == nil || *a != 2000 {
		t.Errorf("monto: %v %v", a, err)
	}
	if a, err := dominio.BudgetAmountAt(ctx, db, sql.NullInt64{}, "2026-11-01"); err != nil || a != nil {
		t.Errorf("total sin monto: %v %v", a, err)
	}
	if bs, err := dominio.BudgetsOfCategory(ctx, db, cn); err != nil || len(bs) != 1 {
		t.Errorf("renglones: %v %v", bs, err)
	}
	if bs, err := dominio.ListBudgets(ctx, db, "2026-10-01", "2027-10-31", nil, ""); err != nil || len(bs) != 1 {
		t.Errorf("lista: %v %v", bs, err)
	}
	if _, err := dominio.GetBudget(ctx, db, 9999); err == nil {
		t.Error("presupuesto inexistente")
	}
	if ps, err := dominio.ListPeriods(ctx, db, "2027-02-01", "2027-02-28", "ORDER BY start_date"); err != nil || len(ps) != 3 || ps[2].PayDate.String != "2027-02-28" {
		t.Errorf("periodos: %+v %v", ps, err)
	}
	if _, err := dominio.RowExists(ctx, db, "person", 1, ""); err == nil {
		t.Error("tabla no admitida")
	}
	if ok, err := dominio.RowExists(ctx, db, "category", cat, "is_active = 1"); err != nil || !ok {
		t.Errorf("categoría: %v %v", ok, err)
	}
	if n, err := dominio.ExpectedIncome(ctx, db, "2026-10-30", true); err != nil || n != 1000000 {
		t.Errorf("esperado: %d %v", n, err)
	}
	if _, err := dominio.GetIncome(ctx, db, 9999); err == nil {
		t.Error("ingreso inexistente")
	}
	if _, err := dominio.GetSavingsMovement(ctx, db, 9999); err == nil {
		t.Error("movimiento inexistente")
	}
	must(dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		for _, k := range []string{"fixed", "extra"} {
			id, err := dominio.BudgetAreaInsert(ctx, tx, "income", []string{"person_id", "kind", "description", "amount_cents", "income_date"},
				[]any{fx.AxelPersonID, k, "Ingreso", 1000, "2026-10-15"}, actor, "r")
			if err != nil {
				return err
			}
			if _, err := dominio.GetIncome(ctx, tx, id); err != nil {
				return err
			}
		}
		if _, err := dominio.BudgetAreaInsert(ctx, tx, "savings_movement", []string{"savings_goal_id", "movement_date", "kind", "amount_cents"},
			[]any{goal, "2026-10-05", "deposit", 500}, actor, "r"); err != nil {
			return err
		}
		id, err := dominio.BudgetAreaInsert(ctx, tx, "budget", []string{"amount_cents", "valid_from_date", "valid_to_date"}, []any{9, "2027-01-01", "2027-01-31"}, actor, "r")
		if err != nil {
			return err
		}
		if _, err := dominio.GetBudget(ctx, tx, id); err != nil {
			return err
		}
		return dominio.BudgetAreaDelete(ctx, tx, "budget", b1, actor, "r")
	}))
	if xs, err := dominio.ListIncomes(ctx, db, "2026-10-01", "2026-10-31", "fixed", "ORDER BY i.id"); err != nil || len(xs) != 1 || !xs[0].CreatedBy.Valid {
		t.Errorf("ingresos: %+v %v", xs, err)
	}
	if ms, err := dominio.ListSavingsMovements(ctx, db, "", "", goal, ""); err != nil || len(ms) != 1 {
		t.Errorf("movimientos: %v %v", ms, err)
	}
	_, _ = db.Exec(`UPDATE category SET kind = 'variable' WHERE id = ?`, cat)
	if m, err := dominio.IncomeByMonth(ctx, db, "2026-01-01", "2026-12-31"); err != nil || m["2026-10"] != [2]int64{1000, 1000} {
		t.Errorf("ingresos por mes: %v %v", m, err)
	}
	if _, err := dominio.ExpenseCountByCategory(ctx, db, "2026-10-01", "2026-10-31"); err != nil {
		t.Error(err)
	}
	if gs, err := dominio.ListSavingsGoals(ctx, db, true, "ORDER BY id"); err != nil || len(gs) != 1 {
		t.Errorf("metas: %v %v", gs, err)
	}
	if ms, err := dominio.ListSavingsMovements(ctx, db, "2026-10-01", "2026-10-31", goal, ""); err != nil || len(ms) != 1 {
		t.Errorf("movimientos: %v %v", ms, err)
	}
	if ps, err := dominio.ActivePlannedItems(ctx, db); err != nil || len(ps) != 2 {
		t.Errorf("partidas: %v %v", ps, err)
	}
	if id, err := dominio.HolderPersonID(ctx, db); err != nil || id != fx.AxelPersonID {
		t.Errorf("titular: %d %v", id, err)
	}

	// Una compra a crédito de 2 pagos quincenales, con partes y un abono, para recorrer compras, pagos y abonos.
	var credit, cp, sp1, sp2 int64
	must(db.QueryRow(`INSERT INTO credit (name, kind, holder_person_id, payment_frequency) VALUES ('Préstamo', 'term_loan', ?, 'fortnightly') RETURNING id`, fx.AxelPersonID).Scan(&credit))
	must(db.QueryRow(`INSERT INTO credit_purchase (credit_id, kind, purchase_date, principal_cents, interest_cents, total_cents, term_count, payment_frequency, first_payment_date, interest_source)
		VALUES (?, 'installment_plan', '2026-10-10', 10000, 1000, 11000, 2, 'fortnightly', '2026-10-30', 'quoted') RETURNING id`, credit).Scan(&cp))
	_, err = db.Exec(`INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 10000, 1000, 11000)`, cp, fx.AxelPersonID)
	must(err)
	must(db.QueryRow(`INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, 1, '2026-10-30', 5500, 5000, 500) RETURNING id`, cp).Scan(&sp1))
	must(db.QueryRow(`INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, 2, '2026-11-15', 5500, 5000, 500) RETURNING id`, cp).Scan(&sp2))
	_, err = db.Exec(`INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 5500), (?, ?, 5500)`, sp1, fx.AxelPersonID, sp2, fx.AxelPersonID)
	must(err)
	_, err = db.Exec(`INSERT INTO credit_payment (credit_id, scheduled_payment_id, payer_person_id, paid_date, amount_cents) VALUES (?, ?, ?, '2026-10-30', 5500)`, credit, sp1, fx.XimenaPersonID)
	must(err)

	fe, err := dominio.LoadFinEngine(ctx, db)
	must(err)
	today := dominio.NewDate(2026, time.October, 12)
	in, err := fe.At(today, nil)
	if err != nil {
		t.Fatalf("entrada del motor: %v", err)
	}
	if fe.Holder != "axel" || len(fe.Externals) != 1 || len(in.PlannedItems) != 2 || in.PlannedItems[0].ShareRule == nil {
		t.Fatalf("motor: %+v", fe)
	}
	if _, err := mf.BudgetMonth(in, dominio.EngineMonth(today)); err != nil {
		t.Fatal(err)
	}
	cut := dominio.NewDate(2026, time.October, 1)
	if early, err := fe.At(today, &cut); err != nil || len(early.Expenses) != 0 {
		t.Errorf("corte: %v %v", early.Expenses, err)
	}
	if len(fe.Expenses()) > 0 {
		one, err := fe.Only(today, fe.Expenses()[0].ID)
		if err != nil || len(one.Expenses) != 1 || len(one.PlannedItems) != 0 {
			t.Errorf("solo un gasto: %v", err)
		}
	}
	if c, ok := fe.Purchase(cp); !ok || len(c.Payments) != 2 || c.Payments[0].Paid[0].Person != "axel" || fe.PurchaseCredit[cp] != credit {
		t.Errorf("compra: %+v", c)
	}
	if _, ok := fe.Purchase(9999); ok {
		t.Error("compra inexistente")
	}
	// RN-21: deuda viva al 31-oct = pagos programados posteriores (el del 15-nov).
	if ld, err := mf.LiveDebt(in, mf.NewDate(2026, time.October, 31)); err != nil || ld.LiveDebtCents != 5500 {
		t.Errorf("deuda viva: %+v %v", ld.LiveDebtCents, err)
	}
	if dominio.FromEngineDate(dominio.EngineDate(today)) != today || dominio.MonthEnd(today).Day != 31 {
		t.Error("conversión de fechas")
	}
}
