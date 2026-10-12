// Package respaldo: respaldo con VACUUM INTO en una carpeta local configurable (versionado por fecha, tipo e id),
// PRAGMA integrity_check, restauración de prueba en una ruta temporal con conteos por tabla comparados con el
// origen (S6-09) y restauración real de un respaldo sobre la base en uso (restoreBackup).
//
// Hora y retención: PROVISIONAL (K31, puerta 4); propuesta y justificación en docs/estado/bloqueos/S6-respaldos.md.
package respaldo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"modernc.org/sqlite"

	"finanzas-hogar/internal/dominio"
)

// Política PROVISIONAL (K31): 03:15 de México todos los días y 30 respaldos diarios (un mes para volver atrás); los
// manuales y los previos a una restauración no se borran solos. Son los mismos valores del job backup_daily.
const (
	DailySchedule = "15 3 * * *"
	KeepScheduled = 30
)

// PreservedTables no se toman del respaldo al restaurar: el catálogo de respaldos (para no perder la referencia a
// los archivos, incluido el respaldo previo) y las sesiones (una sesión cerrada después del respaldo no revive).
// PROVISIONAL: bloqueo S6-respaldos.
var PreservedTables = []string{"backup", "device_session"}

// Result es un respaldo hecho y comprobado contra la base de origen.
type Result struct {
	dominio.BackupResult
	Origin map[string]int64 // filas por tabla de la base de origen justo después del VACUUM INTO
	Diff   []string         // tablas cuyo conteo restaurado no es igual al de origen (vacío = iguales)
}

// Run hace el respaldo (dominio.Backup: VACUUM INTO, huella, integrity_check y restauración de prueba en tmpDir) y
// compara los conteos de la copia restaurada con los de la base de origen.
func Run(ctx context.Context, d *sql.DB, dir, tmpDir, kind string, now time.Time) (Result, error) {
	res, err := dominio.Backup(ctx, d, dir, tmpDir, kind, now)
	out := Result{BackupResult: res}
	if err != nil {
		return out, err
	}
	if out.Origin, err = dominio.TableCounts(ctx, d); err != nil {
		return out, err
	}
	out.Diff = DiffCounts(out.Origin, res.Counts)
	return out, nil
}

// DiffCounts lista las tablas cuyo conteo difiere (o falta) entre a y b.
func DiffCounts(a, b map[string]int64) []string {
	var out []string
	for k, v := range a {
		if w, ok := b[k]; !ok || w != v {
			out = append(out, fmt.Sprintf("%s: %d contra %d", k, v, w))
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			out = append(out, k+": sobra en la copia")
		}
	}
	sort.Strings(out)
	return out
}

// Restored es el resultado de una restauración real.
type Restored struct {
	BackupID     int64
	PreRestoreID int64
	Integrity    string
}

type backupRow struct {
	status, path string
	sha, pruned  sql.NullString
}

// Restore restaura el respaldo backupID sobre la base en uso: comprueba que el archivo exista, conserve su huella y
// pase integrity_check; toma antes un respaldo pre_migration; copia las páginas con la API de respaldo de SQLite;
// aplica las migraciones que falten; repone PreservedTables y corre integrity_check sobre la base restaurada.
func Restore(ctx context.Context, d *sql.DB, backupID int64, dir, tmpDir string, now time.Time) (Restored, error) {
	out := Restored{BackupID: backupID}
	var b backupRow
	err := d.QueryRowContext(ctx, `SELECT status, coalesce(file_path, ''), sha256, pruned_at FROM backup WHERE id = ?`, backupID).Scan(&b.status, &b.path, &b.sha, &b.pruned)
	if errors.Is(err, sql.ErrNoRows) {
		return out, dominio.NotFound("Ese respaldo no existe.")
	}
	if err != nil {
		return out, err
	}
	if b.status != "ok" || b.pruned.Valid || b.path == "" {
		return out, dominio.Conflict("conflict", "Ese respaldo no se puede restaurar: no terminó bien o su archivo ya se borró por la retención.")
	}
	sum, err := fileSHA256(b.path)
	if err != nil || !b.sha.Valid || sum != b.sha.String {
		return out, dominio.Conflict("conflict", "El archivo de ese respaldo falta o cambió desde que se hizo; no se restaura.")
	}
	if ic, err := dominio.IntegrityCheck(ctx, b.path); err != nil || ic != "ok" {
		return out, dominio.Conflict("conflict", "El archivo de ese respaldo no pasó la revisión de integridad; no se restaura.")
	}
	pre, err := dominio.Backup(ctx, d, dir, tmpDir, "pre_migration", now)
	if err != nil {
		return out, fmt.Errorf("respaldo previo a la restauración: %w", err)
	}
	out.PreRestoreID = pre.ID
	conn, err := d.Conn(ctx)
	if err != nil {
		return out, err
	}
	defer conn.Close()
	keep := make(map[string]snapshot, len(PreservedTables))
	for _, t := range PreservedTables {
		if keep[t], err = takeSnapshot(ctx, conn, t); err != nil {
			return out, err
		}
	}
	err = conn.Raw(func(dc any) error {
		rc, ok := dc.(interface {
			NewRestore(string) (*sqlite.Backup, error)
		})
		if !ok {
			return errors.New("respaldo: el controlador de SQLite no permite restaurar")
		}
		bk, err := rc.NewRestore(b.path)
		if err != nil {
			return err
		}
		for {
			more, err := bk.Step(-1)
			if err != nil {
				_ = bk.Finish()
				return err
			}
			if !more {
				break
			}
		}
		return bk.Finish()
	})
	if err != nil {
		return out, fmt.Errorf("restaurar: %w", err)
	}
	if _, err := dominio.Migrate(ctx, d); err != nil {
		return out, fmt.Errorf("migraciones tras restaurar: %w", err)
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return out, err
	}
	for _, t := range PreservedTables {
		if err := keep[t].put(ctx, tx); err != nil {
			_ = tx.Rollback()
			return out, err
		}
	}
	if err := tx.Commit(); err != nil {
		return out, err
	}
	err = conn.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&out.Integrity)
	return out, err
}

type snapshot struct {
	table string
	cols  []string
	rows  [][]any
}

func takeSnapshot(ctx context.Context, conn *sql.Conn, table string) (snapshot, error) {
	s := snapshot{table: table}
	rows, err := conn.QueryContext(ctx, `SELECT * FROM "`+table+`"`)
	if err != nil {
		return s, err
	}
	defer rows.Close()
	if s.cols, err = rows.Columns(); err != nil {
		return s, err
	}
	for rows.Next() {
		vals := make([]any, len(s.cols))
		ptrs := make([]any, len(s.cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return s, err
		}
		s.rows = append(s.rows, vals)
	}
	return s, rows.Err()
}

func (s snapshot) put(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM "`+s.table+`"`); err != nil {
		return err
	}
	cols, marks := "", ""
	for i, c := range s.cols {
		if i > 0 {
			cols, marks = cols+", ", marks+", "
		}
		cols, marks = cols+`"`+c+`"`, marks+"?"
	}
	for _, r := range s.rows {
		if _, err := tx.ExecContext(ctx, `INSERT INTO "`+s.table+`" (`+cols+`) VALUES (`+marks+`)`, r...); err != nil {
			return fmt.Errorf("reponer %s: %w", s.table, err)
		}
	}
	return nil
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
