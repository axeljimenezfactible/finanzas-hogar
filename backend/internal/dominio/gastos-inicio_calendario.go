package dominio

import (
	"context"
	"database/sql"
	"fmt"
)

// CalendarRow es un evento del calendario del hogar (N-029) leído de una tabla, antes de armar el
// CalendarEvent de la API. Los montos son centavos tal como están guardados: aquí no se calcula nada.
type CalendarRow struct {
	Type                string
	Date                string
	EndDate             *string
	Title               string
	AmountCents         *int64
	HouseholdShareCents *int64
	ExternalShareCents  *int64
	Status              *string
	CreditID            *int64
	RefTable            string
	RefID               int64
}

// MarketStore es un mercado activo con su día de visita (RN-16; la fecha la da precios.NextVisit).
type MarketStore struct {
	ID, VisitWeekday int64
	Name             string
}

const calShare = `coalesce((SELECT sum(s.amount_cents) FROM expense_share s JOIN person p ON p.id = s.person_id WHERE s.expense_id = e.id AND p.participation = '%s'), 0)`

// GastosCalendarRows lee los eventos guardados del rango [from, to] (fechas AAAA-MM-DD): ingresos recibidos,
// gastos fijos y venideros, cortes de tarjeta, promociones vigentes, fechas meta y movimientos de ahorro.
// Los pagos de crédito, los mercados y el reporte mensual los arma la API con sus propias consultas.
func GastosCalendarRows(ctx context.Context, q Querier, from, to string) ([]CalendarRow, error) {
	var out []CalendarRow
	type src struct {
		typ, table, sql string
		args            []any
		scan            func(rs *sql.Rows, c *CalendarRow) error
	}
	plain := func(rs *sql.Rows, c *CalendarRow) error {
		var amt sql.NullInt64
		if err := rs.Scan(&c.RefID, &c.Date, &c.Title, &amt); err != nil {
			return err
		}
		if amt.Valid {
			c.AmountCents = &amt.Int64
		}
		return nil
	}
	withCredit := func(rs *sql.Rows, c *CalendarRow) error {
		var amt, cr sql.NullInt64
		var end sql.NullString
		if err := rs.Scan(&c.RefID, &c.Date, &end, &c.Title, &amt, &cr); err != nil {
			return err
		}
		if amt.Valid {
			c.AmountCents = &amt.Int64
		}
		if cr.Valid {
			c.CreditID = &cr.Int64
		}
		if end.Valid {
			c.EndDate = &end.String
		}
		return nil
	}
	hh := fmt.Sprintf(calShare, "holder")
	ext := fmt.Sprintf(calShare, "external")
	srcs := []src{
		{"income", "income", `SELECT id, income_date, description, amount_cents FROM income WHERE income_date BETWEEN ? AND ? ORDER BY income_date, id`, []any{from, to}, plain},
		{"expense", "expense", `SELECT e.id, e.kind, e.expense_date, coalesce(e.description, c.name), e.amount_cents, ` + hh + `, ` + ext + `
 FROM expense e JOIN category c ON c.id = e.category_id WHERE e.kind IN ('fixed', 'upcoming') AND e.expense_date BETWEEN ? AND ? ORDER BY e.expense_date, e.id`, []any{from, to},
			func(rs *sql.Rows, c *CalendarRow) error {
				var kind string
				var amt, h, x int64
				if err := rs.Scan(&c.RefID, &kind, &c.Date, &c.Title, &amt, &h, &x); err != nil {
					return err
				}
				c.Type = map[string]string{"fixed": "fixed_expense", "upcoming": "upcoming_expense"}[kind]
				c.AmountCents, c.HouseholdShareCents, c.ExternalShareCents = &amt, &h, &x
				return nil
			}},
		{"card_statement", "credit_statement", `SELECT s.id, s.statement_date, NULL, cr.name, s.balance_cents, s.credit_id FROM credit_statement s JOIN credit cr ON cr.id = s.credit_id
 WHERE s.statement_date BETWEEN ? AND ? ORDER BY s.statement_date, s.id`, []any{from, to}, withCredit},
		{"promotion", "promotion", `SELECT id, valid_from, valid_to, merchant_name || ' a ' || term_count || ' pagos', NULL, credit_id FROM promotion
 WHERE is_active = 1 AND valid_from <= ? AND valid_to >= ? ORDER BY valid_from, id`, []any{to, from}, withCredit},
		{"savings_goal", "savings_goal", `SELECT id, target_date, name, target_cents FROM savings_goal
 WHERE is_active = 1 AND kind = 'goal' AND target_date BETWEEN ? AND ? ORDER BY target_date, id`, []any{from, to}, plain},
		{"savings_movement", "savings_movement", `SELECT m.id, m.movement_date, CASE m.kind WHEN 'deposit' THEN 'Aportación a ' ELSE 'Retiro de ' END || g.name, m.amount_cents
 FROM savings_movement m JOIN savings_goal g ON g.id = m.savings_goal_id WHERE m.movement_date BETWEEN ? AND ? ORDER BY m.movement_date, m.id`, []any{from, to}, plain},
	}
	for _, s := range srcs {
		rs, err := q.QueryContext(ctx, s.sql, s.args...)
		if err != nil {
			return nil, err
		}
		for rs.Next() {
			c := CalendarRow{Type: s.typ, RefTable: s.table}
			if err := s.scan(rs, &c); err != nil {
				rs.Close()
				return nil, err
			}
			out = append(out, c)
		}
		if err := rs.Close(); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// GastosMarkets lee los mercados activos con día de visita.
func GastosMarkets(ctx context.Context, q Querier) ([]MarketStore, error) {
	rs, err := q.QueryContext(ctx, `SELECT id, name, visit_weekday FROM store WHERE kind = 'market' AND is_active = 1 AND visit_weekday IS NOT NULL ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rs.Close()
	var out []MarketStore
	for rs.Next() {
		var m MarketStore
		if err := rs.Scan(&m.ID, &m.Name, &m.VisitWeekday); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rs.Err()
}

// GastosNotificationID da el id de un aviso por su código (0 si no existe).
func GastosNotificationID(ctx context.Context, q Querier, code string) (int64, error) {
	var id int64
	err := q.QueryRowContext(ctx, `SELECT id FROM notification WHERE code = ?`, code).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}
