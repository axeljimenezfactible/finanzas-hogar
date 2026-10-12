package dominio

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// BackupResult es el resultado de un respaldo con su restauración de prueba (S6-09).
type BackupResult struct {
	ID            int64
	Path          string
	SizeBytes     int64
	SHA256        string
	Integrity     string           // PRAGMA integrity_check de la copia: ok o el error
	RestoreResult string           // ok o el motivo
	Counts        map[string]int64 // filas por tabla de la copia restaurada
}

// Backup hace un respaldo con VACUUM INTO en dir (solo local), comprueba integrity_check, lo restaura en una
// carpeta temporal dentro de tmpDir y compara las filas por tabla; anota todo en la tabla backup.
// kind: scheduled, manual o pre_migration. La hora y la retención son PROVISIONAL (K31, puerta 4).
func Backup(ctx context.Context, d *sql.DB, dir, tmpDir, kind string, now time.Time) (BackupResult, error) {
	var res BackupResult
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return res, err
	}
	t := FormatDBTime(now)
	r, err := d.ExecContext(ctx, `INSERT INTO backup (kind, status, started_at, created_at, updated_at) VALUES (?, 'running', ?, ?, ?)`, kind, t, t, t)
	if err != nil {
		return res, err
	}
	res.ID, _ = r.LastInsertId()
	res.Path = filepath.Join(dir, fmt.Sprintf("finanzas-%s-%s-%d.db", now.In(MX).Format("20060102-150405"), kind, res.ID))
	fail := func(e error) (BackupResult, error) {
		_, _ = d.ExecContext(ctx, `UPDATE backup SET status = 'failed', finished_at = ?, error = ? WHERE id = ?`, FormatDBTime(time.Now()), e.Error(), res.ID)
		return res, e
	}
	if _, err := d.ExecContext(ctx, `VACUUM INTO ?`, res.Path); err != nil {
		return fail(fmt.Errorf("VACUUM INTO: %w", err))
	}
	if res.SHA256, res.SizeBytes, err = fileSHA256(res.Path); err != nil {
		return fail(err)
	}
	if res.Integrity, err = IntegrityCheck(ctx, res.Path); err != nil {
		return fail(err)
	}
	if res.Integrity != "ok" {
		return fail(fmt.Errorf("integrity_check de la copia: %s", res.Integrity))
	}
	res.RestoreResult, res.Counts, err = RestoreTest(ctx, res.Path, tmpDir)
	if err != nil {
		return fail(err)
	}
	fin := FormatDBTime(now)
	_, err = d.ExecContext(ctx, `UPDATE backup SET status = 'ok', finished_at = ?, file_path = ?, size_bytes = ?, sha256 = ?, integrity_check = ?, restore_tested_at = ?, restore_test_result = ? WHERE id = ?`,
		fin, res.Path, res.SizeBytes, res.SHA256, res.Integrity, fin, res.RestoreResult, res.ID)
	return res, err
}

// IntegrityCheck abre el archivo y corre PRAGMA integrity_check.
func IntegrityCheck(ctx context.Context, path string) (string, error) {
	d, err := OpenReadOnly(path)
	if err != nil {
		return "", err
	}
	defer d.Close()
	var s string
	err = d.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&s)
	return s, err
}

// RestoreTest copia el respaldo a una ruta temporal, corre integrity_check y compara las filas por tabla con el respaldo.
func RestoreTest(ctx context.Context, backupPath, tmpDir string) (string, map[string]int64, error) {
	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		return "", nil, err
	}
	work, err := os.MkdirTemp(tmpDir, "restaura-")
	if err != nil {
		return "", nil, err
	}
	defer os.RemoveAll(work)
	dst := filepath.Join(work, "restaura.db")
	if err := copyFile(backupPath, dst); err != nil {
		return "", nil, err
	}
	ic, err := IntegrityCheck(ctx, dst)
	if err != nil {
		return "", nil, err
	}
	if ic != "ok" {
		return "integrity_check: " + ic, nil, nil
	}
	want, err := countsOf(ctx, backupPath)
	if err != nil {
		return "", nil, err
	}
	got, err := countsOf(ctx, dst)
	if err != nil {
		return "", nil, err
	}
	for k, v := range want {
		if got[k] != v {
			return fmt.Sprintf("la tabla %s tiene %d filas y se esperaban %d", k, got[k], v), got, nil
		}
	}
	return "ok", got, nil
}

func countsOf(ctx context.Context, path string) (map[string]int64, error) {
	d, err := OpenReadOnly(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	return TableCounts(ctx, d)
}

// TableCounts cuenta las filas de cada tabla de la base (sin las internas de SQLite).
func TableCounts(ctx context.Context, q Querier) (map[string]int64, error) {
	rows, err := q.QueryContext(ctx, `SELECT name FROM sqlite_schema WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			rows.Close()
			return nil, err
		}
		names = append(names, n)
	}
	rows.Close()
	sort.Strings(names)
	out := make(map[string]int64, len(names))
	for _, n := range names {
		var c int64
		if err := q.QueryRowContext(ctx, `SELECT COUNT(*) FROM "`+n+`"`).Scan(&c); err != nil {
			return nil, err
		}
		out[n] = c
	}
	return out, nil
}

// PruneBackups borra los archivos de respaldos correctos más viejos que los keep más recientes del mismo kind
// y anota pruned_at. La retención (keep) es PROVISIONAL (K31, puerta 4).
func PruneBackups(ctx context.Context, d *sql.DB, kind string, keep int, now time.Time) (int, error) {
	rows, err := d.QueryContext(ctx, `SELECT id, file_path FROM backup WHERE kind = ? AND status = 'ok' AND pruned_at IS NULL ORDER BY started_at DESC, id DESC`, kind)
	if err != nil {
		return 0, err
	}
	type item struct {
		id   int64
		path string
	}
	var all []item
	for rows.Next() {
		var it item
		if err := rows.Scan(&it.id, &it.path); err != nil {
			rows.Close()
			return 0, err
		}
		all = append(all, it)
	}
	rows.Close()
	n := 0
	for i, it := range all {
		if i < keep {
			continue
		}
		if err := os.Remove(it.path); err != nil && !os.IsNotExist(err) {
			return n, err
		}
		if _, err := d.ExecContext(ctx, `UPDATE backup SET pruned_at = ? WHERE id = ?`, FormatDBTime(now), it.id); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

func fileSHA256(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), n, err
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
