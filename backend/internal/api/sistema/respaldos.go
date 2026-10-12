package sistema

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/respaldo"
)

func registerBackups(reg *api.Registry) {
	reg.Handle("listBackups", listBackups(reg))
	reg.Handle("createBackup", createBackup(reg))
	reg.Handle("restoreBackup", restoreBackup(reg))
}

const backupCols = `id, kind, status, started_at, finished_at, file_path, size_bytes, sha256, integrity_check, restore_tested_at, restore_test_result, pruned_at, error, created_at, updated_at`

func scanBackup(s interface{ Scan(...any) error }) (api.Backup, error) {
	var b api.Backup
	var fin, path, sha, ic, rt, rr, pr, er sql.NullString
	var size sql.NullInt64
	err := s.Scan(&b.ID, &b.Kind, &b.Status, &b.StartedAt, &fin, &path, &size, &sha, &ic, &rt, &rr, &pr, &er, &b.CreatedAt, &b.UpdatedAt)
	b.StartedAt, b.CreatedAt, b.UpdatedAt = api.APITime(b.StartedAt), api.APITime(b.CreatedAt), api.APITime(b.UpdatedAt)
	str := func(n sql.NullString, t bool) *string {
		if !n.Valid {
			return nil
		}
		if t {
			return api.Ptr(api.APITime(n.String))
		}
		return api.Ptr(n.String)
	}
	b.FinishedAt, b.FilePath, b.SHA256, b.IntegrityCheck = str(fin, true), str(path, false), str(sha, false), str(ic, false)
	b.RestoreTestedAt, b.RestoreTestResult, b.PrunedAt, b.Error = str(rt, true), str(rr, false), str(pr, true), str(er, false)
	if size.Valid {
		b.SizeBytes = api.Ptr(size.Int64)
	}
	return b, err
}

func listBackups(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"started_at", "kind", "status"}, "-started_at")
		if err != nil {
			return err
		}
		per, err := api.QueryPeriod(r)
		if err != nil {
			return err
		}
		where, args := "", []any{}
		if per != nil { // el periodo es de días de México; started_at está en UTC
			where = ` WHERE started_at >= ? AND started_at < ?`
			args = append(args, dominio.FormatDBTime(per.Start.StartOfDay()), dominio.FormatDBTime(per.End.AddDays(1).StartOfDay()))
		}
		rows, err := reg.Deps().DB.QueryContext(r.Context(), `SELECT `+backupCols+` FROM backup`+where+" "+
			p.OrderBy(map[string]string{"started_at": "started_at", "kind": "kind", "status": "status", "id": "id"})+" "+p.LimitOffset(), args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		items := []api.Backup{}
		for rows.Next() {
			b, err := scanBackup(rows)
			if err != nil {
				return err
			}
			items = append(items, b)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func getBackup(r *http.Request, q dominio.Querier, id int64) (api.Backup, error) {
	return scanBackup(q.QueryRowContext(r.Context(), `SELECT `+backupCols+` FROM backup WHERE id = ?`, id))
}

// createBackup hace el respaldo manual en el momento (VACUUM INTO + integrity_check + restauración de prueba) y
// responde 202 con el renglón; si falló, el renglón dice failed y el error resumido.
func createBackup(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx, d := r.Context(), reg.Deps()
		var running int
		if err := d.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM backup WHERE status = 'running'`).Scan(&running); err != nil {
			return err
		}
		if running > 0 {
			return dominio.Conflict("conflict", "Ya hay un respaldo en curso. Espera a que termine y vuelve a intentar.")
		}
		res, err := respaldo.Run(ctx, d.DB, d.Config.BackupDir, d.Config.TmpDir, "manual", d.Clock.Now())
		if res.ID == 0 {
			return err
		}
		if err != nil {
			d.Log.Error("respaldo manual falló", "backup_id", res.ID, "request_id", api.RequestID(r), "error", err.Error())
		}
		b, err := getBackup(r, d.DB, res.ID)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusAccepted, b)
	}
}

// restoreBackup: acción crítica; confirmation = RESTAURAR (la doble confirmación es de la pantalla).
func restoreBackup(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "backup_id")
		if err != nil {
			return err
		}
		var in api.BackupRestoreRequest
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		if in.Confirmation != "RESTAURAR" {
			return dominio.Invalid(dominio.FieldIssue{Field: "confirmation", Code: "invalid_format", Message: "Para restaurar escribe RESTAURAR en mayúsculas."})
		}
		ctx, d := r.Context(), reg.Deps()
		out, err := respaldo.Restore(ctx, d.DB, id, d.Config.BackupDir, d.Config.TmpDir, d.Clock.Now())
		if err != nil {
			return err
		}
		before, err := dominio.RowJSON(ctx, d.DB, "backup", id)
		if err != nil {
			return err
		}
		after, _ := json.Marshal(map[string]any{"restored_backup_id": out.BackupID, "pre_restore_backup_id": out.PreRestoreID, "integrity_check": out.Integrity})
		if err := dominio.Audit(ctx, d.DB, api.Actor(r), "update", "backup", id, before, after, api.RequestID(r)); err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusAccepted, api.RestoreResult{BackupID: out.BackupID, IntegrityCheck: out.Integrity, PreRestoreBackupID: out.PreRestoreID})
	}
}
