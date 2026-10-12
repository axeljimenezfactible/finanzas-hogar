// Package dominio reúne lo común del backend de Finanzas Hogar: apertura de SQLite (modernc.org/sqlite, puro Go),
// migraciones embebidas (goose como biblioteca), repositorios compartidos (user_account, device_session,
// login_attempt, app_setting, idempotency_key, job, audit_log y backup) y utilidades de dinero en centavos,
// fechas en America/Mexico_City, periodos y quincenas. Las áreas de la API lo usan; no reimplementan nada de esto.
package dominio

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/pressly/goose/v3"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"finanzas-hogar/db"
)

// Querier lo cumplen *sql.DB y *sql.Tx: todos los repositorios aceptan cualquiera de los dos.
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Clock da la hora actual. En producción es SystemClock; las pruebas usan testutil.Clock (reloj simulable).
type Clock interface{ Now() time.Time }

// SystemClock es el reloj del sistema.
type SystemClock struct{}

// Now devuelve la hora actual del sistema.
func (SystemClock) Now() time.Time { return time.Now() }

// Open abre (o crea) la base SQLite con llaves foráneas, WAL, espera de 10 s ante bloqueo y transacciones IMMEDIATE.
func Open(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_txlock=immediate"
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir la base %s: %w", path, err)
	}
	if err := d.Ping(); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("abrir la base %s: %w", path, err)
	}
	return d, nil
}

// OpenReadOnly abre una base solo para lectura y sin cambiar su modo de diario (respaldos y restauraciones de prueba:
// abrirlas con WAL cambiaría el archivo y su huella).
func OpenReadOnly(path string) (*sql.DB, error) {
	d, err := sql.Open("sqlite", "file:"+path+"?mode=ro&_pragma=busy_timeout(10000)")
	if err != nil {
		return nil, err
	}
	if err := d.Ping(); err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("abrir la base %s: %w", path, err)
	}
	return d, nil
}

// MigrationsFS devuelve las migraciones embebidas con los .sql en la raíz.
func MigrationsFS() fs.FS {
	sub, err := fs.Sub(db.Migrations, "migrations")
	if err != nil {
		panic(err) // la ruta es fija y se comprueba en las pruebas
	}
	return sub
}

// Migrate aplica todas las migraciones pendientes (goose como biblioteca) y devuelve la versión final.
func Migrate(ctx context.Context, d *sql.DB) (int64, error) {
	p, err := goose.NewProvider(goose.DialectSQLite3, d, MigrationsFS())
	if err != nil {
		return 0, fmt.Errorf("preparar migraciones: %w", err)
	}
	if _, err := p.Up(ctx); err != nil {
		return 0, fmt.Errorf("aplicar migraciones: %w", err)
	}
	return p.GetDBVersion(ctx)
}

// WithTx corre fn dentro de una transacción: confirma si fn no falla y revierte si falla.
func WithTx(ctx context.Context, d *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// IsConstraint dice si err es una violación de restricción de SQLite (CHECK, UNIQUE, FOREIGN KEY o RAISE de un disparador).
func IsConstraint(err error) bool {
	var se *sqlite.Error
	return errors.As(err, &se) && se.Code()&0xff == sqlite3.SQLITE_CONSTRAINT
}

// IsUnique dice si err es una violación de UNIQUE o de PRIMARY KEY.
func IsUnique(err error) bool {
	var se *sqlite.Error
	return errors.As(err, &se) && (se.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE || se.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY)
}
