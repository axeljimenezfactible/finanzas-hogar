package dominio

import (
	"context"
	"database/sql"
)

// PurchaseRow es un renglón de credit_purchase.
type PurchaseRow struct {
	ID, CreditID              int64
	ExpenseID                 *int64
	Kind                      string
	Description               *string
	PurchaseDate              string
	PrincipalCents            int64
	InterestCents, TotalCents int64
	TermCount                 int64
	PaymentFrequency          string
	FirstPaymentDate          string
	InterestSource            string
	RateNum, RateDen          *int64
	RateBasis                 *string
	InterestIncludesVat       bool
	PromotionID               *int64
	Status                    string
	Notes                     *string
	CreatedAt, UpdatedAt      string
	RateDecimals              *int64
}

const purchaseCols = `p.id, p.credit_id, p.expense_id, p.kind, p.description, p.purchase_date, p.principal_cents,
 p.interest_cents, p.total_cents, p.term_count, p.payment_frequency, p.first_payment_date, p.interest_source, p.rate_num,
 p.rate_den, p.rate_basis, p.interest_includes_vat, p.promotion_id, p.status, p.notes, p.created_at, p.updated_at, p.rate_decimals`
const purchaseWrite = `credit_id, expense_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents,
 term_count, payment_frequency, first_payment_date, interest_source, rate_num, rate_den, rate_basis, interest_includes_vat,
 promotion_id, status, notes, rate_decimals`

// QueryPurchases lee compras a crédito (alias p) con un WHERE y una cola armados por quien llama.
func QueryPurchases(ctx context.Context, q Querier, where, tail string, args ...any) ([]PurchaseRow, error) {
	if where == "" {
		where = "1 = 1"
	}
	rows, err := q.QueryContext(ctx, `SELECT `+purchaseCols+` FROM credit_purchase p WHERE `+where+` `+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PurchaseRow
	for rows.Next() {
		var p PurchaseRow
		if err := rows.Scan(&p.ID, &p.CreditID, &p.ExpenseID, &p.Kind, &p.Description, &p.PurchaseDate, &p.PrincipalCents,
			&p.InterestCents, &p.TotalCents, &p.TermCount, &p.PaymentFrequency, &p.FirstPaymentDate, &p.InterestSource,
			&p.RateNum, &p.RateDen, &p.RateBasis, &p.InterestIncludesVat, &p.PromotionID, &p.Status, &p.Notes,
			&p.CreatedAt, &p.UpdatedAt, &p.RateDecimals); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetPurchase lee una compra a crédito; sql.ErrNoRows si no existe.
func GetPurchase(ctx context.Context, q Querier, id int64) (PurchaseRow, error) {
	ps, err := QueryPurchases(ctx, q, "p.id = ?", "", id)
	if err == nil && len(ps) == 0 {
		err = sql.ErrNoRows
	}
	if err != nil {
		return PurchaseRow{}, err
	}
	return ps[0], nil
}

// InsertPurchase da de alta una compra a crédito.
func InsertPurchase(ctx context.Context, q Querier, p PurchaseRow) (int64, error) {
	return insertRow(ctx, q, "credit_purchase", purchaseWrite, []any{p.CreditID, p.ExpenseID, p.Kind, p.Description,
		p.PurchaseDate, p.PrincipalCents, p.InterestCents, p.TotalCents, p.TermCount, p.PaymentFrequency,
		p.FirstPaymentDate, p.InterestSource, p.RateNum, p.RateDen, p.RateBasis, p.InterestIncludesVat, p.PromotionID,
		p.Status, p.Notes, p.RateDecimals})
}

// UpdatePurchaseText cambia descripción, notas y estado (los montos no se editan, N-219).
func UpdatePurchaseText(ctx context.Context, q Querier, p PurchaseRow) error {
	return updateRow(ctx, q, "credit_purchase", "description, notes, status", []any{p.Description, p.Notes, p.Status}, p.ID)
}

// PurchaseShareRow es un renglón de credit_purchase_share.
type PurchaseShareRow struct {
	ID, CreditPurchaseID, PersonID            int64
	PrincipalCents, InterestCents, TotalCents int64
	CreatedAt, UpdatedAt                      string
}

// InsertPurchaseShare da de alta la parte de una persona.
func InsertPurchaseShare(ctx context.Context, q Querier, s PurchaseShareRow) (int64, error) {
	return insertRow(ctx, q, "credit_purchase_share", "credit_purchase_id, person_id, principal_cents, interest_cents, total_cents",
		[]any{s.CreditPurchaseID, s.PersonID, s.PrincipalCents, s.InterestCents, s.TotalCents})
}

// ScheduledRow es un renglón de scheduled_payment con su crédito.
type ScheduledRow struct {
	ID, CreditPurchaseID, Seq                  int64
	DueDate                                    string
	AmountCents, PrincipalCents, InterestCents int64
	CreatedAt, UpdatedAt                       string
	CreditID                                   int64
}

// InsertScheduled da de alta un pago programado.
func InsertScheduled(ctx context.Context, q Querier, s ScheduledRow) (int64, error) {
	return insertRow(ctx, q, "scheduled_payment", "credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents",
		[]any{s.CreditPurchaseID, s.Seq, s.DueDate, s.AmountCents, s.PrincipalCents, s.InterestCents})
}

// InsertScheduledShare da de alta la parte de una persona en un pago programado.
func InsertScheduledShare(ctx context.Context, q Querier, spID, personID, amount int64) error {
	_, err := insertRow(ctx, q, "scheduled_payment_share", "scheduled_payment_id, person_id, amount_cents", []any{spID, personID, amount})
	return err
}

const scheduledCols = `sp.id, sp.credit_purchase_id, sp.seq, sp.due_date, sp.amount_cents, sp.principal_cents, sp.interest_cents,
 sp.created_at, sp.updated_at, p.credit_id`

// QueryScheduled lee pagos programados (alias sp, compra p) con un WHERE y una cola armados por quien llama.
func QueryScheduled(ctx context.Context, q Querier, where, tail string, args ...any) ([]ScheduledRow, error) {
	if where == "" {
		where = "1 = 1"
	}
	rows, err := q.QueryContext(ctx, `SELECT `+scheduledCols+` FROM scheduled_payment sp JOIN credit_purchase p ON p.id = sp.credit_purchase_id WHERE `+where+` `+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ScheduledRow
	for rows.Next() {
		var s ScheduledRow
		if err := rows.Scan(&s.ID, &s.CreditPurchaseID, &s.Seq, &s.DueDate, &s.AmountCents, &s.PrincipalCents,
			&s.InterestCents, &s.CreatedAt, &s.UpdatedAt, &s.CreditID); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// ScheduledShareRow es la parte de una persona en un pago programado con lo que lleva pagado.
type ScheduledShareRow struct {
	ScheduledPaymentID, PersonID, AmountCents, PaidCents int64
}

// ScheduledShares lee las partes (y lo pagado por cada persona) de los pagos programados de una compra o de todas (purchaseID 0).
func ScheduledShares(ctx context.Context, q Querier, purchaseID int64) (map[int64][]ScheduledShareRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT s.scheduled_payment_id, s.person_id, s.amount_cents,
 COALESCE((SELECT SUM(m.amount_cents) FROM credit_payment m WHERE m.scheduled_payment_id = s.scheduled_payment_id AND m.payer_person_id = s.person_id), 0)
 FROM scheduled_payment_share s JOIN scheduled_payment sp ON sp.id = s.scheduled_payment_id
 WHERE ? = 0 OR sp.credit_purchase_id = ? ORDER BY s.scheduled_payment_id, s.person_id`, purchaseID, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64][]ScheduledShareRow{}
	for rows.Next() {
		var s ScheduledShareRow
		if err := rows.Scan(&s.ScheduledPaymentID, &s.PersonID, &s.AmountCents, &s.PaidCents); err != nil {
			return nil, err
		}
		out[s.ScheduledPaymentID] = append(out[s.ScheduledPaymentID], s)
	}
	return out, rows.Err()
}

// PaymentRow es un renglón de credit_payment.
type PaymentRow struct {
	ID, CreditID         int64
	ScheduledPaymentID   *int64
	CreditStatementID    *int64
	PayerPersonID        int64
	PaidDate             string
	AmountCents          int64
	Notes                *string
	CreatedAt, UpdatedAt string
	PurchaseID, Seq      int64 // del pago programado (0 si es abono a un corte)
}

// QueryPayments lee abonos (alias m) con un WHERE armado por quien llama.
func QueryPayments(ctx context.Context, q Querier, where string, args ...any) ([]PaymentRow, error) {
	if where == "" {
		where = "1 = 1"
	}
	rows, err := q.QueryContext(ctx, `SELECT m.id, m.credit_id, m.scheduled_payment_id, m.credit_statement_id, m.payer_person_id,
 m.paid_date, m.amount_cents, m.notes, m.created_at, m.updated_at, COALESCE(sp.credit_purchase_id, 0), COALESCE(sp.seq, 0)
 FROM credit_payment m LEFT JOIN scheduled_payment sp ON sp.id = m.scheduled_payment_id WHERE `+where+` ORDER BY m.paid_date, m.id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PaymentRow
	for rows.Next() {
		var m PaymentRow
		if err := rows.Scan(&m.ID, &m.CreditID, &m.ScheduledPaymentID, &m.CreditStatementID, &m.PayerPersonID, &m.PaidDate,
			&m.AmountCents, &m.Notes, &m.CreatedAt, &m.UpdatedAt, &m.PurchaseID, &m.Seq); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// InsertPayment registra un abono.
func InsertPayment(ctx context.Context, q Querier, m PaymentRow) (int64, error) {
	return insertRow(ctx, q, "credit_payment", "credit_id, scheduled_payment_id, credit_statement_id, payer_person_id, paid_date, amount_cents, notes",
		[]any{m.CreditID, m.ScheduledPaymentID, m.CreditStatementID, m.PayerPersonID, m.PaidDate, m.AmountCents, m.Notes})
}

// PurchaseShares lee las partes de las compras (purchaseID 0 = todas), por compra.
func PurchaseShares(ctx context.Context, q Querier, purchaseID int64) (map[int64][]PurchaseShareRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, credit_purchase_id, person_id, principal_cents, interest_cents, total_cents, created_at, updated_at
 FROM credit_purchase_share WHERE ? = 0 OR credit_purchase_id = ? ORDER BY credit_purchase_id, id`, purchaseID, purchaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64][]PurchaseShareRow{}
	for rows.Next() {
		var s PurchaseShareRow
		if err := rows.Scan(&s.ID, &s.CreditPurchaseID, &s.PersonID, &s.PrincipalCents, &s.InterestCents, &s.TotalCents,
			&s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out[s.CreditPurchaseID] = append(out[s.CreditPurchaseID], s)
	}
	return out, rows.Err()
}

// StatementRow es un renglón de credit_statement (cortes de tarjeta).
type StatementRow struct {
	ID, CreditID                  int64
	StatementDate, DueDate        string
	BalanceCents                  int64
	NoInterestCents, MinimumCents *int64
	Source                        string
}

// QueryStatements lee los cortes de todos los créditos, en orden de fecha.
func QueryStatements(ctx context.Context, q Querier) ([]StatementRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, credit_id, statement_date, due_date, balance_cents, no_interest_payment_cents,
 minimum_payment_cents, source FROM credit_statement ORDER BY credit_id, statement_date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []StatementRow
	for rows.Next() {
		var s StatementRow
		if err := rows.Scan(&s.ID, &s.CreditID, &s.StatementDate, &s.DueDate, &s.BalanceCents, &s.NoInterestCents,
			&s.MinimumCents, &s.Source); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// PersonRow es lo que créditos usa de person: clave, participación y tope por quincena (K13).
type PersonRow struct {
	ID            int64
	Code          string
	Participation string
	CapCents      *int64
}

// CreditPersons lee las personas por id.
func CreditPersons(ctx context.Context, q Querier) (map[int64]PersonRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, code, participation, fortnight_payment_cap_cents FROM person ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]PersonRow{}
	for rows.Next() {
		var p PersonRow
		if err := rows.Scan(&p.ID, &p.Code, &p.Participation, &p.CapCents); err != nil {
			return nil, err
		}
		out[p.ID] = p
	}
	return out, rows.Err()
}
