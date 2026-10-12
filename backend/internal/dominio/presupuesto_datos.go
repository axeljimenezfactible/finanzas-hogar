package dominio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Tablas del área de presupuesto que escriben BudgetAreaInsert/Update/Delete (lista cerrada: los nombres
// de tabla y columna salen del código, nunca de la petición).
var budgetAreaTables = map[string]bool{"budget": true, "income": true, "savings_movement": true, "period": true}

// BudgetAreaInsert inserta un renglón y escribe audit_log (insert).
func BudgetAreaInsert(ctx context.Context, tx *sql.Tx, table string, cols []string, vals []any, actor *int64, requestID string) (int64, error) {
	if !budgetAreaTables[table] {
		return 0, fmt.Errorf("tabla %q fuera del área de presupuesto", table)
	}
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", table, strings.Join(cols, ", "), strings.TrimSuffix(strings.Repeat("?, ", len(cols)), ", "))
	var id int64
	if err := tx.QueryRowContext(ctx, q, vals...).Scan(&id); err != nil {
		return 0, err
	}
	after, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return 0, err
	}
	return id, Audit(ctx, tx, actor, "insert", table, id, nil, after, requestID)
}

// BudgetAreaUpdate cambia columnas de un renglón, pone updated_at = now y escribe audit_log (update).
func BudgetAreaUpdate(ctx context.Context, tx *sql.Tx, table string, id int64, cols []string, vals []any, now string, actor *int64, requestID string) error {
	if !budgetAreaTables[table] {
		return fmt.Errorf("tabla %q fuera del área de presupuesto", table)
	}
	before, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return err
	}
	sets := make([]string, 0, len(cols)+1)
	for _, c := range cols {
		sets = append(sets, c+" = ?")
	}
	sets = append(sets, "updated_at = ?")
	args := append(append(append([]any{}, vals...), now), id)
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", table, strings.Join(sets, ", ")), args...); err != nil {
		return err
	}
	after, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return err
	}
	return Audit(ctx, tx, actor, "update", table, id, before, after, requestID)
}

// BudgetAreaDelete borra un renglón y escribe audit_log (delete). Sin renglón → ErrNotFound.
func BudgetAreaDelete(ctx context.Context, tx *sql.Tx, table string, id int64, actor *int64, requestID string) error {
	if !budgetAreaTables[table] {
		return fmt.Errorf("tabla %q fuera del área de presupuesto", table)
	}
	before, err := RowJSON(ctx, tx, table, id)
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, "DELETE FROM "+table+" WHERE id = ?", id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return Audit(ctx, tx, actor, "delete", table, id, before, nil, requestID)
}

// BudgetRecord es un renglón de budget.
type BudgetRecord struct {
	ID          int64
	CategoryID  sql.NullInt64
	AmountCents int64
	ValidFrom   string
	ValidTo     sql.NullString
	Notes       sql.NullString
	CreatedAt   string
	UpdatedAt   string
}

const budgetCols = `id, category_id, amount_cents, valid_from_date, valid_to_date, notes, created_at, updated_at`

func scanBudgets(ctx context.Context, q Querier, query string, args ...any) ([]BudgetRecord, error) {
	var out []BudgetRecord
	err := presEach(ctx, q, query, func(s presScanner) error {
		var b BudgetRecord
		if err := s.Scan(&b.ID, &b.CategoryID, &b.AmountCents, &b.ValidFrom, &b.ValidTo, &b.Notes, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return err
		}
		out = append(out, b)
		return nil
	}, args...)
	return out, err
}

// GetBudget lee un renglón de budget (404 si no existe).
func GetBudget(ctx context.Context, q Querier, id int64) (BudgetRecord, error) {
	bs, err := scanBudgets(ctx, q, "SELECT "+budgetCols+" FROM budget WHERE id = ?", id)
	if err != nil {
		return BudgetRecord{}, err
	}
	if len(bs) == 0 {
		return BudgetRecord{}, NotFound("No encontramos ese monto de presupuesto.")
	}
	return bs[0], nil
}

// ListBudgets lee los renglones cuya vigencia toca [from, to], con filtro de categorías; tail = ORDER BY … LIMIT … OFFSET ….
func ListBudgets(ctx context.Context, q Querier, from, to string, cats []int64, tail string) ([]BudgetRecord, error) {
	where := "valid_from_date <= ? AND coalesce(valid_to_date, '9999-12-31') >= ?"
	args := []any{to, from}
	if len(cats) > 0 {
		where += " AND category_id IN (" + strings.TrimSuffix(strings.Repeat("?,", len(cats)), ",") + ")"
		for _, c := range cats {
			args = append(args, c)
		}
	}
	return scanBudgets(ctx, q, "SELECT "+budgetCols+" FROM budget WHERE "+where+" "+tail, args...)
}

// BudgetsOfCategory lee todos los renglones de una categoría (cat inválido = el total del mes).
func BudgetsOfCategory(ctx context.Context, q Querier, cat sql.NullInt64) ([]BudgetRecord, error) {
	return scanBudgets(ctx, q, "SELECT "+budgetCols+" FROM budget WHERE coalesce(category_id, 0) = ? ORDER BY valid_from_date", cat.Int64)
}

// BudgetOverlaps dice si la vigencia [from, to] se traslapa con otro renglón de la misma categoría (o del total).
func BudgetOverlaps(ctx context.Context, q Querier, cat sql.NullInt64, from string, to sql.NullString, exceptID int64) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, `SELECT count(*) FROM budget WHERE coalesce(category_id, 0) = ? AND id <> ?
		AND valid_from_date <= coalesce(?, '9999-12-31') AND coalesce(valid_to_date, '9999-12-31') >= ?`, cat.Int64, exceptID, to, from).Scan(&n)
	return n > 0, err
}

// BudgetAmountAt es el monto que rige el día dado en la categoría (o el total); nil = sin monto (RN-22).
func BudgetAmountAt(ctx context.Context, q Querier, cat sql.NullInt64, day string) (*int64, error) {
	var a int64
	err := q.QueryRowContext(ctx, `SELECT amount_cents FROM budget WHERE coalesce(category_id, 0) = ? AND valid_from_date <= ?
		AND coalesce(valid_to_date, '9999-12-31') >= ? ORDER BY valid_from_date DESC LIMIT 1`, cat.Int64, day, day).Scan(&a)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// RowExists dice si existe el id en una tabla de catálogo (category, planned_item, savings_goal) con condición extra opcional.
func RowExists(ctx context.Context, q Querier, table string, id int64, cond string) (bool, error) {
	switch table {
	case "category", "planned_item", "savings_goal":
	default:
		return false, fmt.Errorf("tabla %q no admitida", table)
	}
	if cond != "" {
		cond = " AND " + cond
	}
	var n int
	err := q.QueryRowContext(ctx, "SELECT count(*) FROM "+table+" WHERE id = ?"+cond, id).Scan(&n)
	return n > 0, err
}

// HolderPersonID es la persona titular (decisión B: los ingresos son del titular).
func HolderPersonID(ctx context.Context, q Querier) (int64, error) {
	var id int64
	err := q.QueryRowContext(ctx, "SELECT id FROM person WHERE participation = 'holder' ORDER BY id LIMIT 1").Scan(&id)
	return id, err
}

// IncomeRecord es un renglón de income con quien lo capturó (audit_log).
type IncomeRecord struct {
	ID            int64
	PersonID      int64
	Kind          string
	PlannedItemID sql.NullInt64
	Description   string
	AmountCents   int64
	IncomeDate    string
	Notes         sql.NullString
	CreatedAt     string
	UpdatedAt     string
	CreatedBy     sql.NullInt64
}

const incomeSelect = `SELECT i.id, i.person_id, i.kind, i.planned_item_id, i.description, i.amount_cents, i.income_date, i.notes, i.created_at, i.updated_at,
	(SELECT a.user_account_id FROM audit_log a WHERE a.table_name = 'income' AND a.row_id = i.id AND a.action = 'insert' ORDER BY a.id LIMIT 1)
	FROM income i`

// ListIncomes lee ingresos con filtros opcionales (fechas inclusivas y kind); tail = ORDER BY … LIMIT ….
func ListIncomes(ctx context.Context, q Querier, from, to, kind string, tail string) ([]IncomeRecord, error) {
	var conds []string
	var args []any
	if from != "" {
		conds, args = append(conds, "i.income_date >= ?"), append(args, from)
	}
	if to != "" {
		conds, args = append(conds, "i.income_date <= ?"), append(args, to)
	}
	if kind != "" {
		conds, args = append(conds, "i.kind = ?"), append(args, kind)
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}
	var out []IncomeRecord
	err := presEach(ctx, q, incomeSelect+where+" "+tail, func(s presScanner) error {
		var x IncomeRecord
		if err := s.Scan(&x.ID, &x.PersonID, &x.Kind, &x.PlannedItemID, &x.Description, &x.AmountCents, &x.IncomeDate, &x.Notes, &x.CreatedAt, &x.UpdatedAt, &x.CreatedBy); err != nil {
			return err
		}
		out = append(out, x)
		return nil
	}, args...)
	return out, err
}

// GetIncome lee un ingreso (404 si no existe).
func GetIncome(ctx context.Context, q Querier, id int64) (IncomeRecord, error) {
	var x IncomeRecord
	err := q.QueryRowContext(ctx, incomeSelect+" WHERE i.id = ?", id).Scan(&x.ID, &x.PersonID, &x.Kind, &x.PlannedItemID, &x.Description,
		&x.AmountCents, &x.IncomeDate, &x.Notes, &x.CreatedAt, &x.UpdatedAt, &x.CreatedBy)
	if errors.Is(err, sql.ErrNoRows) {
		return x, NotFound("No encontramos ese ingreso.")
	}
	return x, err
}

// ExpectedIncome es el monto esperado de la quincena con ese día de pago: Σ partidas de ingreso quincenales activas que
// rigen ese día, amount_cents en la del 15 y second_half_amount_cents (o amount_cents) en la del 30 (C1 IncomeMonth).
func ExpectedIncome(ctx context.Context, q Querier, payDate string, secondHalf bool) (int64, error) {
	col := "amount_cents"
	if secondHalf {
		col = "coalesce(second_half_amount_cents, amount_cents)"
	}
	var n int64
	err := q.QueryRowContext(ctx, `SELECT coalesce(sum(`+col+`), 0) FROM planned_item WHERE kind = 'income' AND frequency = 'fortnightly'
		AND is_active = 1 AND start_date <= ? AND coalesce(end_date, '9999-12-31') >= ?`, payDate, payDate).Scan(&n)
	return n, err
}

// IncomeByMonth suma ingresos fijos y extra por mes (AAAA-MM) entre dos fechas.
func IncomeByMonth(ctx context.Context, q Querier, from, to string) (map[string][2]int64, error) {
	out := map[string][2]int64{}
	err := presEach(ctx, q, `SELECT substr(income_date, 1, 7), kind, sum(amount_cents) FROM income WHERE income_date BETWEEN ? AND ? GROUP BY 1, 2`, func(s presScanner) error {
		var m, k string
		var n int64
		if err := s.Scan(&m, &k, &n); err != nil {
			return err
		}
		v := out[m]
		if k == "fixed" {
			v[0] += n
		} else {
			v[1] += n
		}
		out[m] = v
		return nil
	}, from, to)
	return out, err
}

// ExpenseCountByCategory cuenta gastos por categoría entre dos fechas (fecha del gasto).
func ExpenseCountByCategory(ctx context.Context, q Querier, from, to string) (map[int64]int64, error) {
	out := map[int64]int64{}
	err := presEach(ctx, q, `SELECT category_id, count(*) FROM expense WHERE expense_date BETWEEN ? AND ? AND category_id IS NOT NULL GROUP BY category_id`, func(s presScanner) error {
		var c, n int64
		if err := s.Scan(&c, &n); err != nil {
			return err
		}
		out[c] = n
		return nil
	}, from, to)
	return out, err
}

// SavingsGoalRecord es un renglón de savings_goal.
type SavingsGoalRecord struct {
	ID                  int64
	Name, Kind          string
	TargetCents         sql.NullInt64
	TargetDate          sql.NullString
	PlannedMonthlyCents sql.NullInt64
	StartDate           string
	IsActive            bool
	Notes               sql.NullString
	CreatedAt           string
	UpdatedAt           string
}

// ListSavingsGoals lee fondos y metas activos o archivados; tail = ORDER BY … LIMIT ….
func ListSavingsGoals(ctx context.Context, q Querier, active bool, tail string) ([]SavingsGoalRecord, error) {
	var out []SavingsGoalRecord
	err := presEach(ctx, q, `SELECT id, name, kind, target_cents, target_date, planned_monthly_cents, start_date, is_active, notes, created_at, updated_at
		FROM savings_goal WHERE is_active = ? `+tail, func(s presScanner) error {
		var g SavingsGoalRecord
		if err := s.Scan(&g.ID, &g.Name, &g.Kind, &g.TargetCents, &g.TargetDate, &g.PlannedMonthlyCents, &g.StartDate, &g.IsActive, &g.Notes, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return err
		}
		out = append(out, g)
		return nil
	}, active)
	return out, err
}

// SavingsMovementRecord es un renglón de savings_movement.
type SavingsMovementRecord struct {
	ID            int64
	SavingsGoalID int64
	MovementDate  string
	Kind          string
	AmountCents   int64
	Notes         sql.NullString
	CreatedAt     string
	UpdatedAt     string
}

const movementSelect = `SELECT id, savings_goal_id, movement_date, kind, amount_cents, notes, created_at, updated_at FROM savings_movement`

// ListSavingsMovements lee movimientos con filtros opcionales; tail = ORDER BY … LIMIT ….
func ListSavingsMovements(ctx context.Context, q Querier, from, to string, goal int64, tail string) ([]SavingsMovementRecord, error) {
	var conds []string
	var args []any
	if from != "" {
		conds, args = append(conds, "movement_date >= ? AND movement_date <= ?"), append(args, from, to)
	}
	if goal > 0 {
		conds, args = append(conds, "savings_goal_id = ?"), append(args, goal)
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}
	var out []SavingsMovementRecord
	err := presEach(ctx, q, movementSelect+where+" "+tail, func(s presScanner) error {
		var m SavingsMovementRecord
		if err := s.Scan(&m.ID, &m.SavingsGoalID, &m.MovementDate, &m.Kind, &m.AmountCents, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return err
		}
		out = append(out, m)
		return nil
	}, args...)
	return out, err
}

// GetSavingsMovement lee un movimiento (404 si no existe).
func GetSavingsMovement(ctx context.Context, q Querier, id int64) (SavingsMovementRecord, error) {
	var m SavingsMovementRecord
	err := q.QueryRowContext(ctx, movementSelect+" WHERE id = ?", id).Scan(&m.ID, &m.SavingsGoalID, &m.MovementDate, &m.Kind, &m.AmountCents, &m.Notes, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return m, NotFound("No encontramos ese movimiento de ahorro.")
	}
	return m, err
}

// PeriodRecord es un renglón de period.
type PeriodRecord struct {
	ID        int64
	Kind      string
	ParentID  sql.NullInt64
	Year      int64
	Month     int64
	Half      sql.NullInt64
	StartDate string
	EndDate   string
	PayDate   sql.NullString
	CreatedAt string
	UpdatedAt string
}

// EnsurePeriods crea los meses y quincenas que falten entre dos meses (inclusive), con audit_log.
// La tabla period es calendario (K8, FIRME); se llena a pedido del selector de periodo.
func EnsurePeriods(ctx context.Context, tx *sql.Tx, from, to Date, actor *int64, requestID string) error {
	for m := MonthStart(from.Year, from.Month); !m.After(to); m = MonthStart(m.Year, m.Month+1) {
		var id int64
		err := tx.QueryRowContext(ctx, "SELECT id FROM period WHERE kind = 'month' AND year = ? AND month = ?", m.Year, int(m.Month)).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			id, err = BudgetAreaInsert(ctx, tx, "period", []string{"kind", "year", "month", "start_date", "end_date"},
				[]any{"month", m.Year, int(m.Month), m.String(), MonthEnd(m).String()}, actor, requestID)
		}
		if err != nil {
			return err
		}
		for half := 1; half <= 2; half++ {
			f := FortnightOf(NewDate(m.Year, m.Month, 1+15*(half-1)))
			var n int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM period WHERE kind = 'fortnight' AND parent_id = ? AND half = ?", id, half).Scan(&n); err != nil {
				return err
			}
			if n == 0 {
				if _, err := BudgetAreaInsert(ctx, tx, "period", []string{"kind", "parent_id", "year", "month", "half", "start_date", "end_date", "pay_date"},
					[]any{"fortnight", id, m.Year, int(m.Month), half, f.Start().String(), f.End().String(), f.PayDay().String()}, actor, requestID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ListPeriods lee los periodos que tocan [from, to]; tail = ORDER BY … LIMIT ….
func ListPeriods(ctx context.Context, q Querier, from, to string, tail string) ([]PeriodRecord, error) {
	var out []PeriodRecord
	err := presEach(ctx, q, `SELECT id, kind, parent_id, year, month, half, start_date, end_date, pay_date, created_at, updated_at FROM period
		WHERE start_date <= ? AND end_date >= ? `+tail, func(s presScanner) error {
		var p PeriodRecord
		if err := s.Scan(&p.ID, &p.Kind, &p.ParentID, &p.Year, &p.Month, &p.Half, &p.StartDate, &p.EndDate, &p.PayDate, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return err
		}
		out = append(out, p)
		return nil
	}, to, from)
	return out, err
}

// PlannedItemRecord es un renglón de planned_item (lectura para los supuestos de la proyección).
type PlannedItemRecord struct {
	ID                                          int64
	Kind, Description, Frequency, StartDate     string
	AmountCents, IntervalCount, AccrualMonths   int64
	SecondHalfAmountCents, CategoryID, PersonID sql.NullInt64
	CostShareRuleID, CreditID, TermCount        sql.NullInt64
	EndDate, PaymentMethod, Notes               sql.NullString
	IsActive                                    bool
	CreatedAt, UpdatedAt                        string
}

// ActivePlannedItems lee las partidas planificadas activas (ingreso esperado, gastos fijos y venideros) por fecha de inicio.
func ActivePlannedItems(ctx context.Context, q Querier) ([]PlannedItemRecord, error) {
	var out []PlannedItemRecord
	err := presEach(ctx, q, `SELECT id, kind, description, frequency, start_date, amount_cents, interval_count, accrual_months, second_half_amount_cents,
		category_id, person_id, cost_share_rule_id, credit_id, term_count, end_date, payment_method, notes, is_active, created_at, updated_at
		FROM planned_item WHERE is_active = 1 ORDER BY start_date, id`, func(s presScanner) error {
		var p PlannedItemRecord
		if err := s.Scan(&p.ID, &p.Kind, &p.Description, &p.Frequency, &p.StartDate, &p.AmountCents, &p.IntervalCount, &p.AccrualMonths,
			&p.SecondHalfAmountCents, &p.CategoryID, &p.PersonID, &p.CostShareRuleID, &p.CreditID, &p.TermCount, &p.EndDate,
			&p.PaymentMethod, &p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return err
		}
		out = append(out, p)
		return nil
	})
	return out, err
}
