package dominio_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func TestAjustesRepos(t *testing.T) {
	ctx := context.Background()
	d, fx := testutil.NewDB(t)
	today, _ := dominio.ParseDate("2026-10-12")
	actor := &fx.AxelUserID
	var rule, cat, sub, item, goal int64
	err := dominio.WithTx(ctx, d, func(tx *sql.Tx) error {
		var err error
		if rule, err = dominio.InsertOwned(ctx, tx, "cost_share_rule", map[string]any{"name": "Conjunto", "remainder_person_id": fx.AxelPersonID}, actor, "r1"); err != nil {
			return err
		}
		if err = dominio.ReplaceShareRuleParts(ctx, tx, rule, []dominio.ShareRulePartInput{{PersonID: fx.KarenPersonID, PctBp: ajPtr(int64(5000))}}, actor, "r1"); err != nil {
			return err
		}
		if err = dominio.ReplaceShareRuleParts(ctx, tx, rule, []dominio.ShareRulePartInput{{PersonID: fx.KarenPersonID, FixedCents: ajPtr(int64(119150))}}, actor, "r1"); err != nil {
			return err
		}
		if cat, err = dominio.InsertOwned(ctx, tx, "category", map[string]any{"name": "Despensa", "code": "groceries", "default_cost_share_rule_id": rule}, actor, "r1"); err != nil {
			return err
		}
		if sub, err = dominio.InsertOwned(ctx, tx, "category", map[string]any{"name": "Frutas", "parent_id": cat}, actor, "r1"); err != nil {
			return err
		}
		p := dominio.PlannedItemRow{Kind: "fixed_expense", Description: "Renta", AmountCents: 950000, Frequency: "monthly", IntervalCount: 1,
			StartDate: "2026-11-01", CategoryID: sql.NullInt64{Int64: cat, Valid: true}, Notes: sql.NullString{String: "n", Valid: true}, AccrualMonths: 1, IsActive: true}
		if item, err = dominio.InsertOwned(ctx, tx, "planned_item", p.Columns(), actor, "r1"); err != nil {
			return err
		}
		g := dominio.SavingsGoalRow{Name: "Fondo", Kind: "fund", StartDate: "2026-10-01", IsActive: true, TargetCents: sql.NullInt64{Int64: 100, Valid: true}}
		if goal, err = dominio.InsertOwned(ctx, tx, "savings_goal", g.Columns(), actor, "r1"); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (?, '2026-10-02', 'deposit', 50)`, goal); err != nil {
			return err
		}
		if err = dominio.UpdateOwned(ctx, tx, "category", sub, map[string]any{"sort_order": 5}, actor, "r1"); err != nil {
			return err
		}
		if err = dominio.UpdateOwned(ctx, tx, "category", sub, nil, actor, "r1"); err != nil {
			return err
		}
		_, err = dominio.InsertOwned(ctx, tx, "expense", map[string]any{}, actor, "r1")
		if err == nil || dominio.UpdateOwned(ctx, tx, "expense", 1, map[string]any{"x": 1}, actor, "r1") == nil ||
			dominio.DeleteOwned(ctx, tx, "user_account", 1, actor, "r1") == nil {
			t.Error("las escrituras deben limitarse a las tablas del área")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	cs, err := dominio.ListCategoryRows(ctx, d, dominio.CategoryFilter{Status: "active", Today: today}, "ORDER BY c.id")
	if err != nil || len(cs) != 2 || cs[1].SortOrder != 5 {
		t.Fatalf("%v %+v", err, cs)
	}
	if cs, _ := dominio.ListCategoryRows(ctx, d, dominio.CategoryFilter{Status: "archived", Today: today}, ""); len(cs) != 0 {
		t.Fatal("archivadas")
	}
	if cs, _ := dominio.ListCategoryRows(ctx, d, dominio.CategoryFilter{FilterIDs: true, Today: today}, ""); len(cs) != 0 {
		t.Fatal("ids vacíos")
	}
	if _, err := dominio.CategoryRowByID(ctx, d, 999, today); !isNotFound(err) {
		t.Fatal(err)
	}
	if in, _ := dominio.CategoryInSubtree(ctx, d, cat, sub); !in {
		t.Fatal("Frutas está bajo Despensa")
	}
	if in, _ := dominio.CategoryInSubtree(ctx, d, sub, cat); in {
		t.Fatal("Despensa no está bajo Frutas")
	}
	if n, _ := dominio.NextCategorySortOrder(ctx, d, &cat); n != 6 {
		t.Fatalf("siguiente orden %d", n)
	}
	if n, _ := dominio.NextCategorySortOrder(ctx, d, nil); n != 1 {
		t.Fatalf("siguiente orden raíz %d", n)
	}
	ps, err := dominio.ListPersonRecords(ctx, d, 0, "ORDER BY "+dominio.PersonOrder)
	if err != nil || len(ps) != 3 || ps[0].Participation != "holder" || !ps[0].HasAccount || ps[2].HasAccount {
		t.Fatalf("%v %+v", err, ps)
	}
	if _, err := dominio.PersonRecordByID(ctx, d, 999); !isNotFound(err) {
		t.Fatal(err)
	}
	r, err := dominio.ShareRuleRowByID(ctx, d, rule)
	if err != nil || len(r.Parts) != 1 || r.Parts[0].FixedCents.Int64 != 119150 {
		t.Fatalf("%v %+v", err, r)
	}
	tpl := false
	if rs, _ := dominio.ListShareRuleRows(ctx, d, &tpl, 0, ""); len(rs) != 0 {
		t.Fatal("no hay reglas de un solo uso")
	}
	if _, err := dominio.ShareRuleRowByID(ctx, d, 999); !isNotFound(err) {
		t.Fatal(err)
	}
	pi, err := dominio.ListPlannedItemRows(ctx, d, dominio.PlannedItemFilter{Kind: "fixed_expense", CategoryIDs: []int64{cat}, CreditIDs: []int64{}}, "")
	if err != nil || len(pi) != 1 || pi[0].ID != item {
		t.Fatalf("%v %+v", err, pi)
	}
	if pi, _ := dominio.ListPlannedItemRows(ctx, d, dominio.PlannedItemFilter{CreditIDs: []int64{1}}, ""); len(pi) != 0 {
		t.Fatal("filtro crédito")
	}
	if _, err := dominio.PlannedItemRowByID(ctx, d, 999); !isNotFound(err) {
		t.Fatal(err)
	}
	g, err := dominio.SavingsGoalRowByID(ctx, d, goal)
	if err != nil || g.TargetCents.Int64 != 100 {
		t.Fatalf("%v %+v", err, g)
	}
	if _, err := dominio.SavingsGoalRowByID(ctx, d, 999); !isNotFound(err) {
		t.Fatal(err)
	}
	if mv, err := dominio.SavingsMovementsOf(ctx, d, goal); err != nil || len(mv) != 1 || mv[0].AmountCents != 50 {
		t.Fatalf("%v %+v", err, mv)
	}
	if ok, _ := dominio.RefExists(ctx, d, "category", cat); !ok {
		t.Fatal("existe")
	}
	if _, err := dominio.RefExists(ctx, d, "user_account", 1); err == nil {
		t.Fatal("tabla no admitida")
	}
	if string(dominio.JSONBool(true)) != "true" || string(dominio.JSONBool(false)) != "false" {
		t.Fatal("JSONBool")
	}
	var audits int
	_ = d.QueryRow(`SELECT count(*) FROM audit_log WHERE request_id = 'r1'`).Scan(&audits)
	if audits < 9 {
		t.Fatalf("audit_log: %d", audits)
	}
}

func isNotFound(err error) bool {
	e, ok := dominio.AsError(err)
	return ok && e.Kind == dominio.KindNotFound && !errors.Is(err, sql.ErrNoRows)
}

func ajPtr[T any](v T) *T { return &v }
