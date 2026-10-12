package dominio

// Repositorios del área de Ajustes (agente api-ajustes): category, person, cost_share_rule y sus partes,
// planned_item y savings_goal. Cada escritura deja su renglón en audit_log (CONVENCIONES-backend §8).

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ajOwnedTables son las tablas que escriben InsertOwned, UpdateOwned y DeleteOwned (nombres fijos en código).
var ajOwnedTables = map[string]bool{
	"category": true, "person": true, "cost_share_rule": true, "cost_share_rule_part": true,
	"planned_item": true, "savings_goal": true,
}

func ajOwnedTable(table string) error {
	if !ajOwnedTables[table] {
		return fmt.Errorf("tabla %q fuera del área de ajustes", table)
	}
	return nil
}

func ajSortedCols(set map[string]any) []string {
	cols := make([]string, 0, len(set))
	for c := range set {
		cols = append(cols, c)
	}
	sort.Strings(cols)
	return cols
}

// InsertOwned inserta un renglón (columnas → valor; nil = NULL) y registra el alta en audit_log.
func InsertOwned(ctx context.Context, tx *sql.Tx, table string, set map[string]any, actor *int64, requestID string) (int64, error) {
	if err := ajOwnedTable(table); err != nil {
		return 0, err
	}
	cols := ajSortedCols(set)
	args := make([]any, len(cols))
	for i, c := range cols {
		args[i] = set[c]
	}
	q := `INSERT INTO ` + table + ` (` + strings.Join(cols, ", ") + `) VALUES (` + strings.TrimSuffix(strings.Repeat("?, ", len(cols)), ", ") + `)`
	res, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	after, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return 0, err
	}
	return id, Audit(ctx, tx, actor, "insert", table, id, nil, after, requestID)
}

// UpdateOwned cambia columnas de un renglón y registra antes y después en audit_log. Sin columnas no hace nada.
func UpdateOwned(ctx context.Context, tx *sql.Tx, table string, id int64, set map[string]any, actor *int64, requestID string) error {
	if err := ajOwnedTable(table); err != nil {
		return err
	}
	if len(set) == 0 {
		return nil
	}
	before, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return err
	}
	cols := ajSortedCols(set)
	args := make([]any, 0, len(cols)+1)
	parts := make([]string, len(cols))
	for i, c := range cols {
		parts[i] = c + " = ?"
		args = append(args, set[c])
	}
	args = append(args, id)
	if _, err := tx.ExecContext(ctx, `UPDATE `+table+` SET `+strings.Join(parts, ", ")+` WHERE id = ?`, args...); err != nil {
		return err
	}
	after, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return err
	}
	return Audit(ctx, tx, actor, "update", table, id, before, after, requestID)
}

// DeleteOwned borra un renglón y registra la baja en audit_log.
func DeleteOwned(ctx context.Context, tx *sql.Tx, table string, id int64, actor *int64, requestID string) error {
	if err := ajOwnedTable(table); err != nil {
		return err
	}
	before, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM `+table+` WHERE id = ?`, id); err != nil {
		return err
	}
	return Audit(ctx, tx, actor, "delete", table, id, before, nil, requestID)
}

// RefExists dice si existe el renglón id en una tabla de catálogo (para validar referencias con 422 por campo).
func RefExists(ctx context.Context, q Querier, table string, id int64) (bool, error) {
	switch table {
	case "category", "person", "cost_share_rule", "credit", "planned_item", "savings_goal":
	default:
		return false, fmt.Errorf("tabla %q no admitida", table)
	}
	var n int
	err := q.QueryRowContext(ctx, `SELECT count(*) FROM `+table+` WHERE id = ?`, id).Scan(&n)
	return n > 0, err
}

func ajInInts(ids []int64) (string, []any) {
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return strings.TrimSuffix(strings.Repeat("?, ", len(ids)), ", "), args
}

// ---------- category ----------

// CategoryRow es un renglón de category con sus conteos calculados.
type CategoryRow struct {
	ID                     int64
	ParentID               sql.NullInt64
	Code                   sql.NullString
	Name                   string
	Kind                   sql.NullString
	DefaultCostShareRuleID sql.NullInt64
	SortOrder              int64
	IsActive               bool
	CreatedAt, UpdatedAt   string
	ExpenseCount           int64
	RecentUsageCount       int64
}

// CategoryFilter: Status "active", "archived" o "" (todas); IDs se aplica solo si FilterIDs; Today fija la ventana
// de 90 días de recent_usage_count (de Today − 89 a Today, N-050).
type CategoryFilter struct {
	Status    string
	FilterIDs bool
	IDs       []int64
	Today     Date
}

// ListCategoryRows lista categorías; orderLimit es el ORDER BY y LIMIT ya armados por la API (alias permitidos:
// c.*, expense_count, recent_usage_count).
func ListCategoryRows(ctx context.Context, q Querier, f CategoryFilter, orderLimit string) ([]CategoryRow, error) {
	args := []any{f.Today.AddDays(-89).String(), f.Today.String()}
	where := []string{"1 = 1"}
	switch f.Status {
	case "active":
		where = append(where, "c.is_active = 1")
	case "archived":
		where = append(where, "c.is_active = 0")
	}
	if f.FilterIDs {
		if len(f.IDs) == 0 {
			return []CategoryRow{}, nil
		}
		ph, a := ajInInts(f.IDs)
		where = append(where, "c.id IN ("+ph+")")
		args = append(args, a...)
	}
	rows, err := q.QueryContext(ctx, `SELECT c.id, c.parent_id, c.code, c.name, c.kind, c.default_cost_share_rule_id, c.sort_order,
	       c.is_active, c.created_at, c.updated_at,
	       (SELECT count(*) FROM expense e WHERE e.category_id = c.id) AS expense_count,
	       (SELECT count(*) FROM expense e WHERE e.category_id = c.id AND e.expense_date BETWEEN ?1 AND ?2) AS recent_usage_count
	  FROM category c WHERE `+strings.Join(where, " AND ")+` `+orderLimit, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []CategoryRow{}
	for rows.Next() {
		var c CategoryRow
		if err := rows.Scan(&c.ID, &c.ParentID, &c.Code, &c.Name, &c.Kind, &c.DefaultCostShareRuleID, &c.SortOrder,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.ExpenseCount, &c.RecentUsageCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// CategoryRowByID lee una categoría; si no existe, 404.
func CategoryRowByID(ctx context.Context, q Querier, id int64, today Date) (CategoryRow, error) {
	cs, err := ListCategoryRows(ctx, q, CategoryFilter{FilterIDs: true, IDs: []int64{id}, Today: today}, "")
	if err != nil {
		return CategoryRow{}, err
	}
	if len(cs) == 0 {
		return CategoryRow{}, NotFound("No encontramos esa categoría.")
	}
	return cs[0], nil
}

// CategoryInSubtree dice si candidate es root o una de sus descendientes (para impedir ciclos al cambiar parent_id).
func CategoryInSubtree(ctx context.Context, q Querier, root, candidate int64) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, `WITH RECURSIVE sub(id) AS (SELECT ?1 UNION SELECT c.id FROM category c JOIN sub ON c.parent_id = sub.id)
		SELECT count(*) FROM sub WHERE id = ?2`, root, candidate).Scan(&n)
	return n > 0, err
}

// NextCategorySortOrder es el siguiente sort_order entre las hermanas de parent (nil = raíz).
func NextCategorySortOrder(ctx context.Context, q Querier, parent *int64) (int64, error) {
	var n int64
	var p any
	if parent != nil {
		p = *parent
	}
	err := q.QueryRowContext(ctx, `SELECT coalesce(max(sort_order), 0) + 1 FROM category WHERE parent_id IS ?`, p).Scan(&n)
	return n, err
}

// ---------- person ----------

// PersonRecord es un renglón de person con has_account calculado.
type PersonRecord struct {
	ID                       int64
	Code                     string
	DisplayName              string
	Participation            string
	PaysShareDirectly        bool
	FortnightPaymentCapCents sql.NullInt64
	IsActive                 bool
	Notes                    sql.NullString
	CreatedAt, UpdatedAt     string
	HasAccount               bool
}

// PersonOrder ordena holder, co_holder, external (listPersons).
const PersonOrder = `CASE p.participation WHEN 'holder' THEN 1 WHEN 'co_holder' THEN 2 ELSE 3 END`

// ListPersonRecords lista personas; where filtra por id si id > 0.
func ListPersonRecords(ctx context.Context, q Querier, id int64, orderLimit string) ([]PersonRecord, error) {
	rows, err := q.QueryContext(ctx, `SELECT p.id, p.code, p.display_name, p.participation, p.pays_share_directly,
	       p.fortnight_payment_cap_cents, p.is_active, p.notes, p.created_at, p.updated_at,
	       EXISTS (SELECT 1 FROM user_account u WHERE u.person_id = p.id)
	  FROM person p WHERE (?1 = 0 OR p.id = ?1) `+orderLimit, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PersonRecord{}
	for rows.Next() {
		var p PersonRecord
		if err := rows.Scan(&p.ID, &p.Code, &p.DisplayName, &p.Participation, &p.PaysShareDirectly,
			&p.FortnightPaymentCapCents, &p.IsActive, &p.Notes, &p.CreatedAt, &p.UpdatedAt, &p.HasAccount); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// PersonRecordByID lee una persona; si no existe, 404.
func PersonRecordByID(ctx context.Context, q Querier, id int64) (PersonRecord, error) {
	ps, err := ListPersonRecords(ctx, q, id, "")
	if err != nil {
		return PersonRecord{}, err
	}
	if len(ps) == 0 {
		return PersonRecord{}, NotFound("No encontramos esa persona.")
	}
	return ps[0], nil
}

// ---------- cost_share_rule ----------

// ShareRulePartRow es un renglón de cost_share_rule_part.
type ShareRulePartRow struct {
	ID, RuleID, PersonID int64
	PctBp                sql.NullInt64
	FixedCents           sql.NullInt64
	CreatedAt, UpdatedAt string
}

// ShareRuleRow es un renglón de cost_share_rule con sus partes.
type ShareRuleRow struct {
	ID                   int64
	Name                 string
	RemainderPersonID    int64
	IsTemplate           bool
	IsDefault            bool
	IsActive             bool
	Notes                sql.NullString
	CreatedAt, UpdatedAt string
	Parts                []ShareRulePartRow
}

// ShareRulePartInput es una parte nueva: exactamente uno de PctBp o FixedCents.
type ShareRulePartInput struct {
	PersonID   int64
	PctBp      *int64
	FixedCents *int64
}

// ListShareRuleRows lista reglas (template: nil = todas) con sus partes; id > 0 filtra una.
func ListShareRuleRows(ctx context.Context, q Querier, template *bool, id int64, orderLimit string) ([]ShareRuleRow, error) {
	var t any
	if template != nil {
		t = *template
	}
	rows, err := q.QueryContext(ctx, `SELECT r.id, r.name, r.remainder_person_id, r.is_template, r.is_default, r.is_active, r.notes,
	       r.created_at, r.updated_at FROM cost_share_rule r
	 WHERE (?1 IS NULL OR r.is_template = ?1) AND (?2 = 0 OR r.id = ?2) `+orderLimit, t, id)
	if err != nil {
		return nil, err
	}
	out := []ShareRuleRow{}
	idx := map[int64]int{}
	for rows.Next() {
		var r ShareRuleRow
		if err := rows.Scan(&r.ID, &r.Name, &r.RemainderPersonID, &r.IsTemplate, &r.IsDefault, &r.IsActive, &r.Notes,
			&r.CreatedAt, &r.UpdatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		r.Parts = []ShareRulePartRow{}
		idx[r.ID] = len(out)
		out = append(out, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil || len(out) == 0 {
		return out, err
	}
	ids := make([]int64, len(out))
	for i, r := range out {
		ids[i] = r.ID
	}
	ph, args := ajInInts(ids)
	pr, err := q.QueryContext(ctx, `SELECT id, rule_id, person_id, pct_bp, fixed_cents, created_at, updated_at
	  FROM cost_share_rule_part WHERE rule_id IN (`+ph+`) ORDER BY rule_id, id`, args...)
	if err != nil {
		return nil, err
	}
	defer pr.Close()
	for pr.Next() {
		var p ShareRulePartRow
		if err := pr.Scan(&p.ID, &p.RuleID, &p.PersonID, &p.PctBp, &p.FixedCents, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		r := &out[idx[p.RuleID]]
		r.Parts = append(r.Parts, p)
	}
	return out, pr.Err()
}

// ShareRuleRowByID lee una regla con sus partes; si no existe, 404.
func ShareRuleRowByID(ctx context.Context, q Querier, id int64) (ShareRuleRow, error) {
	rs, err := ListShareRuleRows(ctx, q, nil, id, "")
	if err != nil {
		return ShareRuleRow{}, err
	}
	if len(rs) == 0 {
		return ShareRuleRow{}, NotFound("No encontramos esa regla de reparto.")
	}
	return rs[0], nil
}

// ReplaceShareRuleParts reemplaza las partes de una regla (baja y alta con audit_log por renglón).
func ReplaceShareRuleParts(ctx context.Context, tx *sql.Tx, ruleID int64, parts []ShareRulePartInput, actor *int64, requestID string) error {
	old, err := ShareRuleRowByID(ctx, tx, ruleID)
	if err != nil {
		return err
	}
	for _, p := range old.Parts {
		if err := DeleteOwned(ctx, tx, "cost_share_rule_part", p.ID, actor, requestID); err != nil {
			return err
		}
	}
	for _, p := range parts {
		set := map[string]any{"rule_id": ruleID, "person_id": p.PersonID, "pct_bp": nil, "fixed_cents": nil}
		if p.PctBp != nil {
			set["pct_bp"] = *p.PctBp
		}
		if p.FixedCents != nil {
			set["fixed_cents"] = *p.FixedCents
		}
		if _, err := InsertOwned(ctx, tx, "cost_share_rule_part", set, actor, requestID); err != nil {
			return err
		}
	}
	return nil
}

// ---------- planned_item ----------

// PlannedItemRow es un renglón de planned_item.
type PlannedItemRow struct {
	ID                    int64
	Kind                  string
	Description           string
	AmountCents           int64
	SecondHalfAmountCents sql.NullInt64
	Frequency             string
	IntervalCount         int64
	StartDate             string
	EndDate               sql.NullString
	CategoryID            sql.NullInt64
	PersonID              sql.NullInt64
	CostShareRuleID       sql.NullInt64
	PaymentMethod         sql.NullString
	CreditID              sql.NullInt64
	TermCount             sql.NullInt64
	AccrualMonths         int64
	IsActive              bool
	Notes                 sql.NullString
	CreatedAt, UpdatedAt  string
}

// PlannedItemFilter son los filtros de listPlannedItems.
type PlannedItemFilter struct {
	ID          int64
	CategoryIDs []int64
	CreditIDs   []int64
	Kind        string
}

// ListPlannedItemRows lista partidas con filtros; orderLimit lo arma la API.
func ListPlannedItemRows(ctx context.Context, q Querier, f PlannedItemFilter, orderLimit string) ([]PlannedItemRow, error) {
	where, args := []string{"1 = 1"}, []any{}
	if f.ID > 0 {
		where, args = append(where, "p.id = ?"), append(args, f.ID)
	}
	if f.Kind != "" {
		where, args = append(where, "p.kind = ?"), append(args, f.Kind)
	}
	if len(f.CategoryIDs) > 0 {
		ph, a := ajInInts(f.CategoryIDs)
		where, args = append(where, "p.category_id IN ("+ph+")"), append(args, a...)
	}
	if len(f.CreditIDs) > 0 {
		ph, a := ajInInts(f.CreditIDs)
		where, args = append(where, "p.credit_id IN ("+ph+")"), append(args, a...)
	}
	rows, err := q.QueryContext(ctx, `SELECT p.id, p.kind, p.description, p.amount_cents, p.second_half_amount_cents, p.frequency,
	       p.interval_count, p.start_date, p.end_date, p.category_id, p.person_id, p.cost_share_rule_id, p.payment_method,
	       p.credit_id, p.term_count, p.accrual_months, p.is_active, p.notes, p.created_at, p.updated_at
	  FROM planned_item p WHERE `+strings.Join(where, " AND ")+` `+orderLimit, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PlannedItemRow{}
	for rows.Next() {
		var p PlannedItemRow
		if err := rows.Scan(&p.ID, &p.Kind, &p.Description, &p.AmountCents, &p.SecondHalfAmountCents, &p.Frequency,
			&p.IntervalCount, &p.StartDate, &p.EndDate, &p.CategoryID, &p.PersonID, &p.CostShareRuleID, &p.PaymentMethod,
			&p.CreditID, &p.TermCount, &p.AccrualMonths, &p.IsActive, &p.Notes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// PlannedItemRowByID lee una partida; si no existe, 404.
func PlannedItemRowByID(ctx context.Context, q Querier, id int64) (PlannedItemRow, error) {
	ps, err := ListPlannedItemRows(ctx, q, PlannedItemFilter{ID: id}, "")
	if err != nil {
		return PlannedItemRow{}, err
	}
	if len(ps) == 0 {
		return PlannedItemRow{}, NotFound("No encontramos esa partida.")
	}
	return ps[0], nil
}

// Columns da las columnas escribibles de la partida (para InsertOwned/UpdateOwned).
func (p PlannedItemRow) Columns() map[string]any {
	return map[string]any{
		"kind": p.Kind, "description": p.Description, "amount_cents": p.AmountCents,
		"second_half_amount_cents": ajNullAny(p.SecondHalfAmountCents), "frequency": p.Frequency,
		"interval_count": p.IntervalCount, "start_date": p.StartDate, "end_date": ajNullAny(p.EndDate),
		"category_id": ajNullAny(p.CategoryID), "person_id": ajNullAny(p.PersonID),
		"cost_share_rule_id": ajNullAny(p.CostShareRuleID), "payment_method": ajNullAny(p.PaymentMethod),
		"credit_id": ajNullAny(p.CreditID), "term_count": ajNullAny(p.TermCount), "accrual_months": p.AccrualMonths,
		"is_active": p.IsActive, "notes": ajNullAny(p.Notes),
	}
}

func ajNullAny(v any) any {
	switch x := v.(type) {
	case sql.NullInt64:
		if x.Valid {
			return x.Int64
		}
	case sql.NullString:
		if x.Valid {
			return x.String
		}
	}
	return nil
}

// ---------- savings_goal ----------

// SavingsGoalRow es un renglón de savings_goal.
type SavingsGoalRow struct {
	ID                   int64
	Name                 string
	Kind                 string
	TargetCents          sql.NullInt64
	TargetDate           sql.NullString
	PlannedMonthlyCents  sql.NullInt64
	StartDate            string
	IsActive             bool
	Notes                sql.NullString
	CreatedAt, UpdatedAt string
}

// SavingsMovementRow es un movimiento (para el motor de metas).
type SavingsMovementRow struct {
	MovementDate string
	Kind         string
	AmountCents  int64
}

// SavingsGoalRowByID lee una meta; si no existe, 404.
func SavingsGoalRowByID(ctx context.Context, q Querier, id int64) (SavingsGoalRow, error) {
	var g SavingsGoalRow
	err := q.QueryRowContext(ctx, `SELECT id, name, kind, target_cents, target_date, planned_monthly_cents, start_date, is_active,
	       notes, created_at, updated_at FROM savings_goal WHERE id = ?`, id).Scan(&g.ID, &g.Name, &g.Kind, &g.TargetCents,
		&g.TargetDate, &g.PlannedMonthlyCents, &g.StartDate, &g.IsActive, &g.Notes, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return g, NotFound("No encontramos esa meta.")
	}
	return g, err
}

// Columns da las columnas escribibles de la meta.
func (g SavingsGoalRow) Columns() map[string]any {
	return map[string]any{
		"name": g.Name, "kind": g.Kind, "target_cents": ajNullAny(g.TargetCents), "target_date": ajNullAny(g.TargetDate),
		"planned_monthly_cents": ajNullAny(g.PlannedMonthlyCents), "start_date": g.StartDate, "is_active": g.IsActive,
		"notes": ajNullAny(g.Notes),
	}
}

// SavingsMovementsOf lista los movimientos de una meta en orden de fecha.
func SavingsMovementsOf(ctx context.Context, q Querier, goalID int64) ([]SavingsMovementRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT movement_date, kind, amount_cents FROM savings_movement WHERE savings_goal_id = ? ORDER BY movement_date, id`, goalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SavingsMovementRow{}
	for rows.Next() {
		var m SavingsMovementRow
		if err := rows.Scan(&m.MovementDate, &m.Kind, &m.AmountCents); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// JSONBool codifica un booleano como valor de app_setting.
func JSONBool(b bool) json.RawMessage {
	if b {
		return json.RawMessage("true")
	}
	return json.RawMessage("false")
}
