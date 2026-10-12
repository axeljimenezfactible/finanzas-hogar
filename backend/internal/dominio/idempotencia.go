package dominio

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// IdempotencyRecord es un renglón de idempotency_key (C1 §1 punto 8).
type IdempotencyRecord struct {
	ID         int64
	Method     string
	Path       string
	RequestSHA string
	Status     sql.NullInt64 // NULL = la petición original sigue en proceso
	Body       sql.NullString
	CreatedAt  time.Time
}

// FindIdempotency busca la clave de la cuenta creada después de notBefore (las más viejas cuentan como purgadas).
func FindIdempotency(ctx context.Context, q Querier, userID int64, key string, notBefore time.Time) (IdempotencyRecord, error) {
	var r IdempotencyRecord
	var created string
	err := q.QueryRowContext(ctx, `SELECT id, http_method, request_path, request_sha256, response_status, response_body, created_at FROM idempotency_key
		WHERE user_account_id = ? AND idempotency_key = ? AND created_at >= ?`, userID, key, FormatDBTime(notBefore)).
		Scan(&r.ID, &r.Method, &r.Path, &r.RequestSHA, &r.Status, &r.Body, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrNotFound
	}
	if err != nil {
		return r, err
	}
	r.CreatedAt, err = ParseDBTime(created)
	return r, err
}

// ReserveIdempotency aparta la clave (respuesta NULL) antes de ejecutar la operación; borra antes un renglón vencido
// de la misma clave. Si otra petición la apartó al mismo tiempo, devuelve Conflict.
func ReserveIdempotency(ctx context.Context, q Querier, userID int64, key, method, path, sha string, now, notBefore time.Time) (int64, error) {
	if _, err := q.ExecContext(ctx, `DELETE FROM idempotency_key WHERE user_account_id = ? AND idempotency_key = ? AND created_at < ?`, userID, key, FormatDBTime(notBefore)); err != nil {
		return 0, err
	}
	t := FormatDBTime(now)
	res, err := q.ExecContext(ctx, `INSERT INTO idempotency_key (user_account_id, idempotency_key, http_method, request_path, request_sha256, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, key, method, path, sha, t, t)
	if IsUnique(err) {
		return 0, Conflict("conflict", "Esta captura ya se está procesando. Espera un momento y vuelve a intentar.")
	}
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// CompleteIdempotency guarda la respuesta (código y cuerpo JSON) para repetirla igual.
func CompleteIdempotency(ctx context.Context, q Querier, id int64, status int, body []byte) error {
	var b any
	if len(body) > 0 {
		b = string(body)
	}
	_, err := q.ExecContext(ctx, `UPDATE idempotency_key SET response_status = ?, response_body = ? WHERE id = ?`, status, b, id)
	return err
}

// ReleaseIdempotency borra la reserva cuando la operación no se completó (respuesta no 2xx): la clave puede reintentarse.
func ReleaseIdempotency(ctx context.Context, q Querier, id int64) error {
	_, err := q.ExecContext(ctx, `DELETE FROM idempotency_key WHERE id = ?`, id)
	return err
}

// PurgeIdempotency borra las claves creadas antes de before (job idempotency_purge). Devuelve cuántas borró.
func PurgeIdempotency(ctx context.Context, q Querier, before time.Time) (int64, error) {
	res, err := q.ExecContext(ctx, `DELETE FROM idempotency_key WHERE created_at < ?`, FormatDBTime(before))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
