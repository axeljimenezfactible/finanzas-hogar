package dominio

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Job es un renglón de la tabla job (S6-08): estado y última ejecución visibles.
type Job struct {
	ID             int64
	Code           string
	Name           string
	Schedule       string // cron de 5 campos en hora de America/Mexico_City
	IsEnabled      bool
	NotificationID sql.NullInt64
	LastRunAt      sql.NullString
	LastStatus     sql.NullString // running, ok, error, skipped
	LastError      sql.NullString
	NextRunAt      sql.NullString
	CreatedAt      string
	UpdatedAt      string
}

const jobCols = `id, code, name, schedule, is_enabled, notification_id, last_run_at, last_status, last_error, next_run_at, created_at, updated_at`

func scanJob(r interface{ Scan(...any) error }) (Job, error) {
	var j Job
	err := r.Scan(&j.ID, &j.Code, &j.Name, &j.Schedule, &j.IsEnabled, &j.NotificationID, &j.LastRunAt, &j.LastStatus, &j.LastError, &j.NextRunAt, &j.CreatedAt, &j.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return j, ErrNotFound
	}
	return j, err
}

// ListJobs devuelve los jobs ordenados por código.
func ListJobs(ctx context.Context, q Querier) ([]Job, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+jobCols+` FROM job ORDER BY code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Job
	for rows.Next() {
		j, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

// GetJob lee un job por id.
func GetJob(ctx context.Context, q Querier, id int64) (Job, error) {
	return scanJob(q.QueryRowContext(ctx, `SELECT `+jobCols+` FROM job WHERE id = ?`, id))
}

// GetJobByCode lee un job por código.
func GetJobByCode(ctx context.Context, q Querier, code string) (Job, error) {
	return scanJob(q.QueryRowContext(ctx, `SELECT `+jobCols+` FROM job WHERE code = ?`, code))
}

// EnsureJob crea el renglón del job si no existe (no cambia uno existente: Ajustes manda sobre el valor inicial).
func EnsureJob(ctx context.Context, q Querier, code, name, schedule string) error {
	_, err := q.ExecContext(ctx, `INSERT INTO job (code, name, schedule) VALUES (?, ?, ?) ON CONFLICT (code) DO NOTHING`, code, name, schedule)
	return err
}

// SetJobNextRun fija la próxima corrida sin tocar el resto.
func SetJobNextRun(ctx context.Context, q Querier, id int64, next time.Time) error {
	_, err := q.ExecContext(ctx, `UPDATE job SET next_run_at = ? WHERE id = ?`, FormatDBTime(next), id)
	return err
}

// ClaimJob toma un disparo de forma atómica: solo una llamada gana para el mismo next_run_at (idempotencia por disparo).
// Marca running, anota last_run_at = now y mueve next_run_at a next. Devuelve false si otro ya lo tomó.
func ClaimJob(ctx context.Context, q Querier, id int64, expectedNext string, now, next time.Time) (bool, error) {
	res, err := q.ExecContext(ctx, `UPDATE job SET last_status = 'running', last_run_at = ?, last_error = NULL, next_run_at = ?
		WHERE id = ? AND is_enabled = 1 AND next_run_at = ? AND (last_status IS NULL OR last_status <> 'running')`,
		FormatDBTime(now), FormatDBTime(next), id, expectedNext)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}

// FinishJob anota el resultado de la corrida (ok, error o skipped) y el error resumido.
func FinishJob(ctx context.Context, q Querier, id int64, status, errMsg string) error {
	_, err := q.ExecContext(ctx, `UPDATE job SET last_status = ?, last_error = ? WHERE id = ?`, status, nullable(errMsg), id)
	return err
}

// ResetInterruptedJobs pasa a error los jobs que quedaron en running por un reinicio (al arrancar el planificador).
func ResetInterruptedJobs(ctx context.Context, q Querier) error {
	_, err := q.ExecContext(ctx, `UPDATE job SET last_status = 'error', last_error = 'Interrumpido por un reinicio del servicio.' WHERE last_status = 'running'`)
	return err
}
