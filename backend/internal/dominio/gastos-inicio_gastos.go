package dominio

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
)

// Repositorio de gastos (expense, expense_item, expense_share) del agente api-gastos-inicio.
// El reparto sigue RN-03: regla del renglón → regla del gasto → regla de la categoría → regla por defecto
// («Hogar», 100 % al titular); porcentajes con ROUND_HALF_UP, fijos tal cual y la persona del resto absorbe la diferencia.

// ExpenseSharePart es una parte de una regla de reparto: porcentaje en puntos base o monto fijo (cost_share_rule_part).
type ExpenseSharePart struct {
	PersonID   int64
	PctBP      *int64
	FixedCents *int64
}

// ExpenseShareAmount es la parte congelada de una persona (expense_share).
type ExpenseShareAmount struct {
	PersonID    int64
	AmountCents int64
}

// ExpenseItemIn es un renglón de ticket por guardar.
type ExpenseItemIn struct {
	LineNo          int64
	Description     string
	PresentationID  *int64
	QuantityMilli   int64
	UnitPriceCents  *int64
	LineTotalCents  int64
	CostShareRuleID *int64
	Parts           []ExpenseSharePart // reparto «mixto»: crea una regla de un solo uso
}

// ExpenseIn es un gasto completo por guardar (alta o cambio ya mezclado).
type ExpenseIn struct {
	Kind            string
	Description     *string
	CategoryID      int64
	AmountCents     int64
	ExpenseDate     string
	AccrualMonths   int64
	PaymentMethod   string
	PaidByPersonID  *int64
	PersonID        *int64
	CostShareRuleID *int64
	StoreID         *int64
	PlannedItemID   *int64
	Notes           *string
	CapturedAt      *string // UTC de la base (FormatDBTime)
	IdempotencyKey  *string
	Parts           []ExpenseSharePart // reparto propio del gasto: crea una regla de un solo uso
	Items           []ExpenseItemIn    // renglones; con ItemsSet reemplazan los guardados
	ItemsSet        bool
}

// ExpenseRow es un gasto leído con sus datos calculados.
type ExpenseRow struct {
	ID                                                                int64
	Kind                                                              string
	Description                                                       sql.NullString
	CategoryID, AmountCents                                           int64
	ExpenseDate                                                       string
	AccrualMonths                                                     int64
	PaymentMethod                                                     string
	PaidByPersonID, PersonID, CostShareRuleID, StoreID, PlannedItemID sql.NullInt64
	Notes                                                             sql.NullString
	CreatedAt, UpdatedAt                                              string
	CapturedAt                                                        sql.NullString
	HouseholdCents, ExternalCents                                     int64
	HasAttachment                                                     bool
	CreditPurchaseID                                                  sql.NullInt64
	Shares                                                            []ExpenseShareAmount
}

// ExpenseFilter son los filtros de listExpenses, getExpenseTotals y getExpenseSummary (C1 §6).
type ExpenseFilter struct {
	PersonIDs, CategoryIDs, StoreIDs, CreditIDs []int64
	Family                                      string // household, external, joint (incluyentes, bloqueo C1-filtro-familia opción A)
	Period                                      *Period
	Q                                           string
	Kinds, PaymentMethods                       []string
	HasAttachment                               *bool
}

const expenseCols = `e.id, e.kind, e.description, e.category_id, e.amount_cents, e.expense_date, e.accrual_months, e.payment_method,
 e.paid_by_person_id, e.person_id, e.cost_share_rule_id, e.store_id, e.planned_item_id, e.notes, e.created_at, e.updated_at, e.captured_at,
 (SELECT coalesce(sum(s.amount_cents), 0) FROM expense_share s JOIN person p ON p.id = s.person_id WHERE s.expense_id = e.id AND p.participation = 'holder'),
 EXISTS (SELECT 1 FROM attachment a WHERE a.expense_id = e.id),
 (SELECT cp.id FROM credit_purchase cp WHERE cp.expense_id = e.id)`

func scanExpense(sc interface{ Scan(...any) error }) (ExpenseRow, error) {
	var x ExpenseRow
	err := sc.Scan(&x.ID, &x.Kind, &x.Description, &x.CategoryID, &x.AmountCents, &x.ExpenseDate, &x.AccrualMonths, &x.PaymentMethod,
		&x.PaidByPersonID, &x.PersonID, &x.CostShareRuleID, &x.StoreID, &x.PlannedItemID, &x.Notes, &x.CreatedAt, &x.UpdatedAt, &x.CapturedAt,
		&x.HouseholdCents, &x.HasAttachment, &x.CreditPurchaseID)
	x.ExternalCents = x.AmountCents - x.HouseholdCents
	return x, err
}

// GastosFoldSQL envuelve una expresión SQL para comparar sin mayúsculas ni acentos (lower de SQLite solo cubre ASCII).
func GastosFoldSQL(expr string) string {
	s := "lower(" + expr + ")"
	for _, p := range [][2]string{{"á", "a"}, {"é", "e"}, {"í", "i"}, {"ó", "o"}, {"ú", "u"}, {"ü", "u"}, {"ñ", "n"},
		{"Á", "a"}, {"É", "e"}, {"Í", "i"}, {"Ó", "o"}, {"Ú", "u"}, {"Ü", "u"}, {"Ñ", "n"}} {
		s = "replace(" + s + ", '" + p[0] + "', '" + p[1] + "')"
	}
	return s
}

// GastosFoldText es GastosFoldSQL del lado de Go (para el texto buscado).
func GastosFoldText(s string) string {
	return strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n").Replace(strings.ToLower(strings.TrimSpace(s)))
}

// GastosLikePattern arma %texto% con los comodines escapados (ESCAPE '\').
func GastosLikePattern(s string) string {
	return "%" + strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(GastosFoldText(s)) + "%"
}

func gxInList(col string, n int) string {
	return col + " IN (" + strings.TrimSuffix(strings.Repeat("?,", n), ",") + ")"
}

// Where arma la cláusula WHERE (con alias e) y sus argumentos.
func (f ExpenseFilter) Where() (string, []any) {
	var w []string
	var args []any
	addIDs := func(expr string, ids []int64) {
		if len(ids) == 0 {
			return
		}
		w = append(w, gxInList(expr, len(ids)))
		for _, id := range ids {
			args = append(args, id)
		}
	}
	addIDs("e.person_id", f.PersonIDs)
	addIDs("e.category_id", f.CategoryIDs)
	addIDs("e.store_id", f.StoreIDs)
	if len(f.CreditIDs) > 0 {
		w = append(w, "EXISTS (SELECT 1 FROM credit_purchase cp WHERE cp.expense_id = e.id AND "+gxInList("cp.credit_id", len(f.CreditIDs))+")")
		for _, id := range f.CreditIDs {
			args = append(args, id)
		}
	}
	for _, x := range []struct {
		col  string
		vals []string
	}{{"e.kind", f.Kinds}, {"e.payment_method", f.PaymentMethods}} {
		if len(x.vals) > 0 {
			w = append(w, gxInList(x.col, len(x.vals)))
			for _, v := range x.vals {
				args = append(args, v)
			}
		}
	}
	const part = "EXISTS (SELECT 1 FROM expense_share s JOIN person p ON p.id = s.person_id WHERE s.expense_id = e.id AND s.amount_cents > 0 AND p.participation = '%s')"
	hh, ext := strings.Replace(part, "%s", "holder", 1), strings.Replace(part, "%s", "external", 1)
	switch f.Family {
	case "household":
		w = append(w, hh)
	case "external":
		w = append(w, ext)
	case "joint":
		w = append(w, hh, ext)
	}
	if f.Period != nil {
		w = append(w, "e.expense_date BETWEEN ? AND ?")
		args = append(args, f.Period.Start.String(), f.Period.End.String())
	}
	if f.HasAttachment != nil {
		not := ""
		if !*f.HasAttachment {
			not = "NOT "
		}
		w = append(w, not+"EXISTS (SELECT 1 FROM attachment a WHERE a.expense_id = e.id)")
	}
	if strings.TrimSpace(f.Q) != "" {
		pat := GastosLikePattern(f.Q)
		w = append(w, "("+GastosFoldSQL("coalesce(e.description, '')")+` LIKE ? ESCAPE '\' OR `+GastosFoldSQL("coalesce(e.notes, '')")+` LIKE ? ESCAPE '\'`+
			` OR EXISTS (SELECT 1 FROM category c WHERE c.id = e.category_id AND `+GastosFoldSQL("c.name")+` LIKE ? ESCAPE '\')`+
			` OR EXISTS (SELECT 1 FROM store st WHERE st.id = e.store_id AND `+GastosFoldSQL("st.name")+` LIKE ? ESCAPE '\'))`)
		args = append(args, pat, pat, pat, pat)
	}
	if len(w) == 0 {
		return "1 = 1", args
	}
	return strings.Join(w, " AND "), args
}

// ListExpenses lee una página de gastos con sus partes (orderBy y limitOffset vienen de api.PageParams).
func ListExpenses(ctx context.Context, q Querier, f ExpenseFilter, orderBy, limitOffset string) ([]ExpenseRow, error) {
	where, args := f.Where()
	rows, err := q.QueryContext(ctx, "SELECT "+expenseCols+" FROM expense e WHERE "+where+" "+orderBy+" "+limitOffset, args...)
	if err != nil {
		return nil, err
	}
	var out []ExpenseRow
	for rows.Next() {
		x, err := scanExpense(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, x)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		if out[i].Shares, err = expenseShares(ctx, q, out[i].ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// ExpenseTotals son las sumas sobre todo el conjunto filtrado.
type ExpenseTotals struct{ Count, AmountCents, HouseholdCents, ExternalCents int64 }

// SumExpenses suma el conjunto filtrado (todas las páginas). La parte externa es la de personas external.
func SumExpenses(ctx context.Context, q Querier, f ExpenseFilter) (ExpenseTotals, error) {
	where, args := f.Where()
	var t ExpenseTotals
	err := q.QueryRowContext(ctx, `SELECT count(*), coalesce(sum(e.amount_cents), 0),
 coalesce(sum((SELECT coalesce(sum(s.amount_cents), 0) FROM expense_share s JOIN person p ON p.id = s.person_id WHERE s.expense_id = e.id AND p.participation = 'holder')), 0),
 coalesce(sum((SELECT coalesce(sum(s.amount_cents), 0) FROM expense_share s JOIN person p ON p.id = s.person_id WHERE s.expense_id = e.id AND p.participation = 'external')), 0)
 FROM expense e WHERE `+where, args...).Scan(&t.Count, &t.AmountCents, &t.HouseholdCents, &t.ExternalCents)
	return t, err
}

// TopHouseholdCategory da la categoría con mayor parte del hogar del conjunto filtrado (ok = false si no hay gastos).
func TopHouseholdCategory(ctx context.Context, q Querier, f ExpenseFilter) (id, cents int64, ok bool, err error) {
	where, args := f.Where()
	err = q.QueryRowContext(ctx, `SELECT e.category_id, sum((SELECT coalesce(sum(s.amount_cents), 0) FROM expense_share s JOIN person p ON p.id = s.person_id
 WHERE s.expense_id = e.id AND p.participation = 'holder')) AS h FROM expense e WHERE `+where+` GROUP BY e.category_id ORDER BY h DESC, e.category_id LIMIT 1`, args...).Scan(&id, &cents)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, false, nil
	}
	return id, cents, err == nil, err
}

func expenseShares(ctx context.Context, q Querier, id int64) ([]ExpenseShareAmount, error) {
	rows, err := q.QueryContext(ctx, `SELECT person_id, amount_cents FROM expense_share WHERE expense_id = ? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ExpenseShareAmount
	for rows.Next() {
		var s ExpenseShareAmount
		if err := rows.Scan(&s.PersonID, &s.AmountCents); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// GetExpense lee un gasto (ErrNotFound si no existe).
func GetExpense(ctx context.Context, q Querier, id int64) (ExpenseRow, error) {
	x, err := scanExpense(q.QueryRowContext(ctx, "SELECT "+expenseCols+" FROM expense e WHERE e.id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return x, NotFound("No encontramos ese gasto. Puede que alguien lo haya eliminado.")
	}
	if err != nil {
		return x, err
	}
	x.Shares, err = expenseShares(ctx, q, id)
	return x, err
}

// ExpenseIDByIdempotencyKey busca la captura ya registrada con esa clave (K9).
func ExpenseIDByIdempotencyKey(ctx context.Context, q Querier, key string) (int64, bool, error) {
	var id int64
	err := q.QueryRowContext(ctx, `SELECT id FROM expense WHERE idempotency_key = ?`, key).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return id, err == nil, err
}

// ExpenseRefExists dice si hay un renglón con ese id en la tabla (nombre fijo del código, nunca del usuario).
func ExpenseRefExists(ctx context.Context, q Querier, table string, id int64) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, "SELECT count(*) FROM "+table+" WHERE id = ?", id).Scan(&n)
	return n > 0, err
}

// GastosHolderID es la persona titular (person.participation = holder).
func GastosHolderID(ctx context.Context, q Querier) (int64, error) {
	var id int64
	err := q.QueryRowContext(ctx, `SELECT id FROM person WHERE participation = 'holder' ORDER BY id LIMIT 1`).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, Conflict("conflict", "Falta registrar a la persona titular. Revisa Ajustes «Personas y reparto».")
	}
	return id, err
}

// ExpenseSplitByRule reparte un monto (RN-03): porcentaje = HALF_UP(monto × pct_bp / 10000), fijo tal cual, el resto a rest.
func ExpenseSplitByRule(amount int64, parts []ExpenseSharePart, rest int64) ([]ExpenseShareAmount, error) {
	left := amount
	out := make([]ExpenseShareAmount, 0, len(parts)+1)
	for _, p := range parts {
		var v int64
		switch {
		case p.PctBP != nil:
			v = ApplyBasisPoints(amount, *p.PctBP)
		case p.FixedCents != nil:
			v = *p.FixedCents
		}
		out = append(out, ExpenseShareAmount{p.PersonID, v})
		left -= v
	}
	if left < 0 {
		return nil, Invalid(FieldIssue{"amount_cents", "out_of_range", "Las partes del reparto suman más que el importe. Revisa los montos fijos."})
	}
	return append(out, ExpenseShareAmount{rest, left}), nil
}

func gxLoadRule(ctx context.Context, q Querier, ruleID int64) ([]ExpenseSharePart, int64, error) {
	var rest int64
	if err := q.QueryRowContext(ctx, `SELECT remainder_person_id FROM cost_share_rule WHERE id = ?`, ruleID).Scan(&rest); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, Invalid(FieldIssue{"cost_share_rule_id", "not_found", "No existe esa regla de reparto."})
		}
		return nil, 0, err
	}
	rows, err := q.QueryContext(ctx, `SELECT person_id, pct_bp, fixed_cents FROM cost_share_rule_part WHERE rule_id = ? ORDER BY id`, ruleID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var parts []ExpenseSharePart
	for rows.Next() {
		var p ExpenseSharePart
		var pct, fix sql.NullInt64
		if err := rows.Scan(&p.PersonID, &pct, &fix); err != nil {
			return nil, 0, err
		}
		if pct.Valid {
			p.PctBP = &pct.Int64
		}
		if fix.Valid {
			p.FixedCents = &fix.Int64
		}
		parts = append(parts, p)
	}
	return parts, rest, rows.Err()
}

// effectiveRule: regla del gasto → de la categoría → por defecto activa → 0 («Hogar» implícita, 100 % al titular, RN-03).
func gxEffectiveRule(ctx context.Context, q Querier, ruleID *int64, categoryID int64) (int64, error) {
	if ruleID != nil {
		return *ruleID, nil
	}
	var id sql.NullInt64
	err := q.QueryRowContext(ctx, `SELECT coalesce(
 (SELECT default_cost_share_rule_id FROM category WHERE id = ?),
 (SELECT id FROM cost_share_rule WHERE is_default = 1 AND is_active = 1 ORDER BY id LIMIT 1))`, categoryID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id.Int64, nil
}

// ExpenseCheckParts valida las partes de un reparto propio (campo field en los errores).
func ExpenseCheckParts(ctx context.Context, q Querier, parts []ExpenseSharePart, field string) error {
	seen := map[int64]bool{}
	for _, p := range parts {
		ok, err := ExpenseRefExists(ctx, q, "person", p.PersonID)
		if err != nil {
			return err
		}
		switch {
		case !ok:
			return Invalid(FieldIssue{field, "not_found", "Una de las personas del reparto no existe."})
		case seen[p.PersonID]:
			return Invalid(FieldIssue{field, "conflict", "Cada persona puede aparecer una sola vez en el reparto."})
		case (p.PctBP == nil) == (p.FixedCents == nil):
			return Invalid(FieldIssue{field, "invalid_format", "Cada parte lleva un porcentaje o un monto fijo, no los dos."})
		case p.PctBP != nil && (*p.PctBP < 0 || *p.PctBP > 10000):
			return Invalid(FieldIssue{field, "out_of_range", "El porcentaje va de 0 a 100 %."})
		case p.FixedCents != nil && *p.FixedCents < 0:
			return Invalid(FieldIssue{field, "out_of_range", "El monto fijo no puede ser negativo."})
		}
		seen[p.PersonID] = true
	}
	return nil
}

// oneOffRule crea una regla de un solo uso (is_template = 0) con el resto al titular.
func gxOneOffRule(ctx context.Context, q Querier, name string, parts []ExpenseSharePart) (int64, error) {
	holder, err := GastosHolderID(ctx, q)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := q.QueryRowContext(ctx, `INSERT INTO cost_share_rule (name, remainder_person_id, is_template) VALUES (?, ?, 0) RETURNING id`, name, holder).Scan(&id); err != nil {
		return 0, err
	}
	for _, p := range parts {
		if _, err := q.ExecContext(ctx, `INSERT INTO cost_share_rule_part (rule_id, person_id, pct_bp, fixed_cents) VALUES (?, ?, ?, ?)`, id, p.PersonID, p.PctBP, p.FixedCents); err != nil {
			return 0, err
		}
	}
	return id, nil
}

// ExpenseComputeShares calcula las partes de un gasto con sus renglones (RN-03): agrupa por regla, reparte cada grupo y suma por persona.
// La diferencia entre el importe y la suma de renglones va con la regla del gasto.
func ExpenseComputeShares(ctx context.Context, q Querier, amount int64, ruleID *int64, categoryID int64, items []ExpenseItemIn) ([]ExpenseShareAmount, error) {
	base, err := gxEffectiveRule(ctx, q, ruleID, categoryID)
	if err != nil {
		return nil, err
	}
	groups := map[int64]int64{}
	sum := int64(0)
	for _, it := range items {
		r := base
		if it.CostShareRuleID != nil {
			r = *it.CostShareRuleID
		}
		groups[r] += it.LineTotalCents
		sum += it.LineTotalCents
	}
	if sum > amount {
		return nil, Invalid(FieldIssue{"items", "out_of_range", "La suma de los renglones es mayor que el importe del gasto."})
	}
	groups[base] += amount - sum
	keys := make([]int64, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	holder, err := GastosHolderID(ctx, q)
	if err != nil {
		return nil, err
	}
	total := map[int64]int64{}
	var order []int64
	for _, k := range keys {
		parts, rest := []ExpenseSharePart(nil), holder
		if k != 0 {
			if parts, rest, err = gxLoadRule(ctx, q, k); err != nil {
				return nil, err
			}
		}
		sh, err := ExpenseSplitByRule(groups[k], parts, rest)
		if err != nil {
			return nil, err
		}
		for _, s := range sh {
			if _, ok := total[s.PersonID]; !ok {
				order = append(order, s.PersonID)
			}
			total[s.PersonID] += s.AmountCents
		}
	}
	var out []ExpenseShareAmount
	for _, p := range order {
		if total[p] != 0 {
			out = append(out, ExpenseShareAmount{p, total[p]})
		}
	}
	if len(out) == 0 {
		out = []ExpenseShareAmount{{holder, 0}}
	}
	return out, nil
}

func expenseItemsIn(ctx context.Context, q Querier, id int64) ([]ExpenseItemIn, error) {
	rows, err := q.QueryContext(ctx, `SELECT line_no, description, presentation_id, quantity_milli, unit_price_cents, line_total_cents, cost_share_rule_id
 FROM expense_item WHERE expense_id = ? ORDER BY line_no`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ExpenseItemIn
	for rows.Next() {
		var it ExpenseItemIn
		var pres, unit, rule sql.NullInt64
		if err := rows.Scan(&it.LineNo, &it.Description, &pres, &it.QuantityMilli, &unit, &it.LineTotalCents, &rule); err != nil {
			return nil, err
		}
		it.PresentationID, it.UnitPriceCents, it.CostShareRuleID = gxNullPtr(pres), gxNullPtr(unit), gxNullPtr(rule)
		out = append(out, it)
	}
	return out, rows.Err()
}

func gxNullPtr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}

// SaveExpense da de alta (id = 0) o reescribe un gasto con sus renglones y partes, en la transacción tx.
// recompute = false conserva las partes congeladas (un cambio que no toca importe, categoría, regla ni renglones).
func SaveExpense(ctx context.Context, tx *sql.Tx, id int64, in ExpenseIn, recompute bool) (int64, error) {
	if len(in.Parts) > 0 {
		if err := ExpenseCheckParts(ctx, tx, in.Parts, "cost_share_rule_parts"); err != nil {
			return 0, err
		}
		r, err := gxOneOffRule(ctx, tx, "Reparto propio del gasto", in.Parts)
		if err != nil {
			return 0, err
		}
		in.CostShareRuleID = &r
	}
	args := []any{in.Kind, in.Description, in.CategoryID, in.AmountCents, in.ExpenseDate, in.AccrualMonths, in.PaymentMethod,
		in.PaidByPersonID, in.PersonID, in.CostShareRuleID, in.StoreID, in.PlannedItemID, in.Notes}
	var err error
	if id == 0 {
		err = tx.QueryRowContext(ctx, `INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, accrual_months, payment_method,
 paid_by_person_id, person_id, cost_share_rule_id, store_id, planned_item_id, notes, captured_at, idempotency_key)
 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`, append(args, in.CapturedAt, in.IdempotencyKey)...).Scan(&id)
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE expense SET kind = ?, description = ?, category_id = ?, amount_cents = ?, expense_date = ?, accrual_months = ?,
 payment_method = ?, paid_by_person_id = ?, person_id = ?, cost_share_rule_id = ?, store_id = ?, planned_item_id = ?, notes = ? WHERE id = ?`, append(args, id)...)
	}
	if err != nil {
		return 0, err
	}
	items := in.Items
	if in.ItemsSet {
		if _, err := tx.ExecContext(ctx, `DELETE FROM expense_item WHERE expense_id = ?`, id); err != nil {
			return 0, err
		}
		for i := range items {
			it := &items[i]
			if len(it.Parts) > 0 {
				if err := ExpenseCheckParts(ctx, tx, it.Parts, "items"); err != nil {
					return 0, err
				}
				r, err := gxOneOffRule(ctx, tx, "Reparto mixto del renglón "+strings.TrimSpace(it.Description), it.Parts)
				if err != nil {
					return 0, err
				}
				it.CostShareRuleID = &r
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO expense_item (expense_id, line_no, description, presentation_id, quantity_milli, unit_price_cents,
 line_total_cents, cost_share_rule_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, id, it.LineNo, it.Description, it.PresentationID, it.QuantityMilli,
				it.UnitPriceCents, it.LineTotalCents, it.CostShareRuleID); err != nil {
				return 0, err
			}
		}
	} else if recompute {
		if items, err = expenseItemsIn(ctx, tx, id); err != nil {
			return 0, err
		}
	}
	if !recompute {
		return id, nil
	}
	shares, err := ExpenseComputeShares(ctx, tx, in.AmountCents, in.CostShareRuleID, in.CategoryID, items)
	if err != nil {
		return 0, err
	}
	return id, ReplaceExpenseShares(ctx, tx, id, shares)
}

// ReplaceExpenseShares reescribe las partes congeladas de un gasto.
func ReplaceExpenseShares(ctx context.Context, q Querier, id int64, shares []ExpenseShareAmount) error {
	if _, err := q.ExecContext(ctx, `DELETE FROM expense_share WHERE expense_id = ?`, id); err != nil {
		return err
	}
	for _, s := range shares {
		if _, err := q.ExecContext(ctx, `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, ?)`, id, s.PersonID, s.AmountCents); err != nil {
			return err
		}
	}
	return nil
}

// ExpenseInOf convierte un gasto guardado en entrada (para mezclar un cambio o duplicar).
func ExpenseInOf(ctx context.Context, q Querier, x ExpenseRow) (ExpenseIn, error) {
	str := func(n sql.NullString) *string {
		if !n.Valid {
			return nil
		}
		v := n.String
		return &v
	}
	in := ExpenseIn{Kind: x.Kind, Description: str(x.Description), CategoryID: x.CategoryID, AmountCents: x.AmountCents, ExpenseDate: x.ExpenseDate,
		AccrualMonths: x.AccrualMonths, PaymentMethod: x.PaymentMethod, PaidByPersonID: gxNullPtr(x.PaidByPersonID), PersonID: gxNullPtr(x.PersonID),
		CostShareRuleID: gxNullPtr(x.CostShareRuleID), StoreID: gxNullPtr(x.StoreID), PlannedItemID: gxNullPtr(x.PlannedItemID), Notes: str(x.Notes)}
	var err error
	in.Items, err = expenseItemsIn(ctx, q, x.ID)
	return in, err
}

// CopyExpenseShares copia las partes congeladas de un gasto a otro (duplicar).
func CopyExpenseShares(ctx context.Context, q Querier, from, to int64) error {
	sh, err := expenseShares(ctx, q, from)
	if err != nil {
		return err
	}
	return ReplaceExpenseShares(ctx, q, to, sh)
}

// DeleteExpense borra un gasto; con compra a crédito vinculada responde 409 has_credit_purchase.
func DeleteExpense(ctx context.Context, q Querier, x ExpenseRow) error {
	if x.CreditPurchaseID.Valid {
		return Conflict("has_credit_purchase", "Este gasto viene de una compra a crédito. Primero elimina la compra en Créditos.")
	}
	_, err := q.ExecContext(ctx, `DELETE FROM expense WHERE id = ?`, x.ID)
	return err
}

// ExpensePlannedItemKind da el tipo de la partida planificada ("" si no hay).
func ExpensePlannedItemKind(ctx context.Context, q Querier, id *int64) (string, error) {
	if id == nil {
		return "", nil
	}
	var k string
	err := q.QueryRowContext(ctx, `SELECT kind FROM planned_item WHERE id = ?`, *id).Scan(&k)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return k, err
}

// ExpenseItemRow es un renglón leído.
type ExpenseItemRow struct {
	ID, ExpenseID, LineNo                     int64
	Description                               string
	PresentationID, UnitPriceCents, CostShare sql.NullInt64
	QuantityMilli, LineTotalCents             int64
	CreatedAt, UpdatedAt                      string
}

// ListExpenseItems lee los renglones de un gasto en una página.
func ListExpenseItems(ctx context.Context, q Querier, id int64, orderBy, limitOffset string) ([]ExpenseItemRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT i.id, i.expense_id, i.line_no, i.description, i.presentation_id, i.unit_price_cents, i.cost_share_rule_id,
 i.quantity_milli, i.line_total_cents, i.created_at, i.updated_at FROM expense_item i WHERE i.expense_id = ? `+orderBy+" "+limitOffset, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ExpenseItemRow
	for rows.Next() {
		var r ExpenseItemRow
		if err := rows.Scan(&r.ID, &r.ExpenseID, &r.LineNo, &r.Description, &r.PresentationID, &r.UnitPriceCents, &r.CostShare,
			&r.QuantityMilli, &r.LineTotalCents, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// ExpenseItemsSum da la suma de renglones y el mayor line_no de un gasto.
func ExpenseItemsSum(ctx context.Context, q Querier, id int64) (sum, maxLine, n int64, err error) {
	err = q.QueryRowContext(ctx, `SELECT coalesce(sum(line_total_cents), 0), coalesce(max(line_no), 0), count(*) FROM expense_item WHERE expense_id = ?`, id).Scan(&sum, &maxLine, &n)
	return
}
