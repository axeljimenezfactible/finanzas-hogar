package dominio

import (
	"context"
	"database/sql"
	"strings"
)

// Repositorios de créditos (api-creditos): credit, credit_term_rate, promotion, credit_purchase,
// credit_purchase_share, scheduled_payment, scheduled_payment_share, credit_payment y
// credit_statement. Solo leen y escriben renglones; los cálculos son del motor S3 (motor_cred).

// CreditRow es un renglón de credit (NULL = puntero nil).
type CreditRow struct {
	ID                    int64
	Name, Kind            string
	Lender                *string
	HolderPersonID        int64
	CreditLimitCents      *int64
	PaymentFrequency      string
	FirstPaymentMinDays   *int64
	AvailabilityReading   string
	AvailabilityConfirmed bool
	ReadingAAnchorDate    *string
	StatementDay          *int64
	PaymentDueDay         *int64
	AnnualRateNum         *int64
	AnnualRateDen         *int64
	VatOnInterestBp       *int64
	AnnualFeeCents        *int64
	LateFeeCents          *int64
	OpeningFeeCents       *int64
	IsActive              bool
	Notes                 *string
	CreatedAt, UpdatedAt  string
	RateBase              *string
	MinPaymentPctBp       *int64
	MinPaymentFixedCents  *int64
	AnnualRateDecimals    *int64
	VatIncluded           *bool
}

const creditCols = `c.id, c.name, c.kind, c.lender, c.holder_person_id, c.credit_limit_cents, c.payment_frequency,
 c.first_payment_min_days, c.availability_reading, c.availability_confirmed, c.reading_a_anchor_date, c.statement_day,
 c.payment_due_day, c.annual_rate_num, c.annual_rate_den, c.vat_on_interest_bp, c.annual_fee_cents, c.late_fee_cents,
 c.opening_fee_cents, c.is_active, c.notes, c.created_at, c.updated_at, c.rate_base, c.min_payment_pct_bp,
 c.min_payment_fixed_cents, c.annual_rate_decimals, c.vat_included`

func scanCredit(s interface{ Scan(...any) error }) (CreditRow, error) {
	var c CreditRow
	err := s.Scan(&c.ID, &c.Name, &c.Kind, &c.Lender, &c.HolderPersonID, &c.CreditLimitCents, &c.PaymentFrequency,
		&c.FirstPaymentMinDays, &c.AvailabilityReading, &c.AvailabilityConfirmed, &c.ReadingAAnchorDate, &c.StatementDay,
		&c.PaymentDueDay, &c.AnnualRateNum, &c.AnnualRateDen, &c.VatOnInterestBp, &c.AnnualFeeCents, &c.LateFeeCents,
		&c.OpeningFeeCents, &c.IsActive, &c.Notes, &c.CreatedAt, &c.UpdatedAt, &c.RateBase, &c.MinPaymentPctBp,
		&c.MinPaymentFixedCents, &c.AnnualRateDecimals, &c.VatIncluded)
	return c, err
}

// GetCredit lee un crédito; sql.ErrNoRows si no existe (la API lo vuelve 404).
func GetCredit(ctx context.Context, q Querier, id int64) (CreditRow, error) {
	return scanCredit(q.QueryRowContext(ctx, `SELECT `+creditCols+` FROM credit c WHERE c.id = ?`, id))
}

// QueryCredits lee créditos con un WHERE y una cola (ORDER BY … LIMIT …) armados por quien llama.
func QueryCredits(ctx context.Context, q Querier, where, tail string, args ...any) ([]CreditRow, error) {
	if where == "" {
		where = "1 = 1"
	}
	rows, err := q.QueryContext(ctx, `SELECT `+creditCols+` FROM credit c WHERE `+where+` `+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CreditRow
	for rows.Next() {
		c, err := scanCredit(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func creditValues(c CreditRow) []any {
	return []any{c.Name, c.Kind, c.Lender, c.HolderPersonID, c.CreditLimitCents, c.PaymentFrequency, c.FirstPaymentMinDays,
		c.AvailabilityReading, c.AvailabilityConfirmed, c.ReadingAAnchorDate, c.StatementDay, c.PaymentDueDay,
		c.AnnualRateNum, c.AnnualRateDen, c.VatOnInterestBp, c.AnnualFeeCents, c.LateFeeCents, c.OpeningFeeCents,
		c.IsActive, c.Notes, c.RateBase, c.MinPaymentPctBp, c.MinPaymentFixedCents, c.AnnualRateDecimals, c.VatIncluded}
}

const creditWriteCols = `name, kind, lender, holder_person_id, credit_limit_cents, payment_frequency, first_payment_min_days,
 availability_reading, availability_confirmed, reading_a_anchor_date, statement_day, payment_due_day, annual_rate_num,
 annual_rate_den, vat_on_interest_bp, annual_fee_cents, late_fee_cents, opening_fee_cents, is_active, notes, rate_base,
 min_payment_pct_bp, min_payment_fixed_cents, annual_rate_decimals, vat_included`

// InsertCredit da de alta un crédito y devuelve su id.
func InsertCredit(ctx context.Context, q Querier, c CreditRow) (int64, error) {
	return insertRow(ctx, q, "credit", creditWriteCols, creditValues(c))
}

// UpdateCredit reescribe todas las columnas editables de un crédito.
func UpdateCredit(ctx context.Context, q Querier, c CreditRow) error {
	return updateRow(ctx, q, "credit", creditWriteCols, creditValues(c), c.ID)
}

func insertRow(ctx context.Context, q Querier, table, cols string, vals []any) (int64, error) {
	ph := strings.TrimSuffix(strings.Repeat("?, ", len(vals)), ", ")
	res, err := q.ExecContext(ctx, `INSERT INTO `+table+` (`+cols+`) VALUES (`+ph+`)`, vals...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func updateRow(ctx context.Context, q Querier, table, cols string, vals []any, id int64) error {
	parts := strings.Split(cols, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i]) + " = ?"
	}
	res, err := q.ExecContext(ctx, `UPDATE `+table+` SET `+strings.Join(parts, ", ")+` WHERE id = ?`, append(vals, id)...)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteRow borra un renglón por id; sql.ErrNoRows si no existía.
func DeleteRow(ctx context.Context, q Querier, table string, id int64) error {
	res, err := q.ExecContext(ctx, `DELETE FROM `+table+` WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Exists dice si hay un renglón con ese id en la tabla.
func Exists(ctx context.Context, q Querier, table string, id int64) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, `SELECT count(*) FROM `+table+` WHERE id = ?`, id).Scan(&n)
	return n > 0, err
}

// TermRateRow es un renglón de credit_term_rate.
type TermRateRow struct {
	ID, CreditID, TermCount, RateNum, RateDen int64
	Source                                    string
	QuoteDate, Notes                          *string
	CreatedAt, UpdatedAt                      string
	RateDecimals                              *int64
}

const termRateCols = `id, credit_id, term_count, rate_num, rate_den, source, quote_date, notes, created_at, updated_at, rate_decimals`
const termRateWrite = `credit_id, term_count, rate_num, rate_den, source, quote_date, notes, rate_decimals`

// QueryTermRates lee tasas por plazo con un WHERE y una cola armados por quien llama.
func QueryTermRates(ctx context.Context, q Querier, where, tail string, args ...any) ([]TermRateRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+termRateCols+` FROM credit_term_rate WHERE `+where+` `+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TermRateRow
	for rows.Next() {
		var t TermRateRow
		if err := rows.Scan(&t.ID, &t.CreditID, &t.TermCount, &t.RateNum, &t.RateDen, &t.Source, &t.QuoteDate, &t.Notes,
			&t.CreatedAt, &t.UpdatedAt, &t.RateDecimals); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// GetTermRate lee una tasa por plazo; sql.ErrNoRows si no existe.
func GetTermRate(ctx context.Context, q Querier, id int64) (TermRateRow, error) {
	ts, err := QueryTermRates(ctx, q, "id = ?", "", id)
	if err == nil && len(ts) == 0 {
		err = sql.ErrNoRows
	}
	if err != nil {
		return TermRateRow{}, err
	}
	return ts[0], nil
}

func termRateValues(t TermRateRow) []any {
	return []any{t.CreditID, t.TermCount, t.RateNum, t.RateDen, t.Source, t.QuoteDate, t.Notes, t.RateDecimals}
}

// InsertTermRate da de alta una tasa por plazo.
func InsertTermRate(ctx context.Context, q Querier, t TermRateRow) (int64, error) {
	return insertRow(ctx, q, "credit_term_rate", termRateWrite, termRateValues(t))
}

// UpdateTermRate reescribe una tasa por plazo.
func UpdateTermRate(ctx context.Context, q Querier, t TermRateRow) error {
	return updateRow(ctx, q, "credit_term_rate", termRateWrite, termRateValues(t), t.ID)
}

// PromotionRow es un renglón de promotion.
type PromotionRow struct {
	ID                   int64
	CreditID             *int64
	MerchantName         string
	StoreID              *int64
	TermCount            int64
	RateNum, RateDen     int64
	MinPurchaseCents     int64
	ValidFrom, ValidTo   string
	SourceNote           *string
	IsActive             bool
	CreatedAt, UpdatedAt string
	RateDecimals         *int64
}

const promoCols = `id, credit_id, merchant_name, store_id, term_count, rate_num, rate_den, min_purchase_cents, valid_from,
 valid_to, source_note, is_active, created_at, updated_at, rate_decimals`
const promoWrite = `credit_id, merchant_name, store_id, term_count, rate_num, rate_den, min_purchase_cents, valid_from,
 valid_to, source_note, is_active, rate_decimals`

// QueryPromotions lee promociones con un WHERE y una cola armados por quien llama.
func QueryPromotions(ctx context.Context, q Querier, where, tail string, args ...any) ([]PromotionRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+promoCols+` FROM promotion WHERE `+where+` `+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PromotionRow
	for rows.Next() {
		var p PromotionRow
		if err := rows.Scan(&p.ID, &p.CreditID, &p.MerchantName, &p.StoreID, &p.TermCount, &p.RateNum, &p.RateDen,
			&p.MinPurchaseCents, &p.ValidFrom, &p.ValidTo, &p.SourceNote, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
			&p.RateDecimals); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GetPromotion lee una promoción; sql.ErrNoRows si no existe.
func GetPromotion(ctx context.Context, q Querier, id int64) (PromotionRow, error) {
	ps, err := QueryPromotions(ctx, q, "id = ?", "", id)
	if err == nil && len(ps) == 0 {
		err = sql.ErrNoRows
	}
	if err != nil {
		return PromotionRow{}, err
	}
	return ps[0], nil
}

func promoValues(p PromotionRow) []any {
	return []any{p.CreditID, p.MerchantName, p.StoreID, p.TermCount, p.RateNum, p.RateDen, p.MinPurchaseCents,
		p.ValidFrom, p.ValidTo, p.SourceNote, p.IsActive, p.RateDecimals}
}

// InsertPromotion da de alta una promoción.
func InsertPromotion(ctx context.Context, q Querier, p PromotionRow) (int64, error) {
	return insertRow(ctx, q, "promotion", promoWrite, promoValues(p))
}

// UpdatePromotion reescribe una promoción.
func UpdatePromotion(ctx context.Context, q Querier, p PromotionRow) error {
	return updateRow(ctx, q, "promotion", promoWrite, promoValues(p), p.ID)
}

// CountWhere cuenta renglones de una tabla con un WHERE armado por quien llama.
func CountWhere(ctx context.Context, q Querier, table, where string, args ...any) (int64, error) {
	var n int64
	err := q.QueryRowContext(ctx, `SELECT count(*) FROM `+table+` WHERE `+where, args...).Scan(&n)
	return n, err
}
