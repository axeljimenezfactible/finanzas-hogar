// Package testutil ayuda a las pruebas del backend: base temporal migrada y sembrada con fixtures mínimas, reloj
// simulable, servidor httptest (TLS) con sesión y validación de cada respuesta contra docs/contratos/openapi.yaml
// (kin-openapi, solo en pruebas). Los temporales van a datos/tmp/ (S6) y se borran al terminar la prueba.
// Las contraseñas de las fixtures son falsas y solo existen aquí.
package testutil

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"finanzas-hogar/internal/dominio"
)

// Usuarios y credencial FALSOS de las fixtures (nunca son los reales; FakeCredential es la contraseña de prueba).
const (
	AxelUser       = "axel"
	XimenaUser     = "ximena"
	FakeCredential = "contrasena-falsa-de-prueba"
)

// FastArgon2 abarata el hash de las fixtures (VerifyPassword lee los parámetros del propio hash).
var FastArgon2 = dominio.Argon2Params{Time: 1, Memory: 8 * 1024, Threads: 1, SaltLen: 16, KeyLen: 32}

// Fixtures son los ids de lo que siembra Seed.
type Fixtures struct {
	AxelPersonID, XimenaPersonID, KarenPersonID int64
	AxelUserID, XimenaUserID                    int64
}

// Root devuelve la raíz del repositorio (la carpeta que contiene backend/go.mod).
func Root(t testing.TB) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "backend", "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("testutil: no encontré la raíz del repositorio")
		}
		dir = parent
	}
}

// TempDir crea una carpeta temporal en datos/tmp/ que se borra al terminar la prueba.
func TempDir(t testing.TB) string {
	t.Helper()
	base := filepath.Join(Root(t), "datos", "tmp")
	if err := os.MkdirAll(base, 0o700); err != nil {
		t.Fatal(err)
	}
	d, err := os.MkdirTemp(base, "gotest-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

// Clock es un reloj simulable y seguro entre gorrutinas.
type Clock struct {
	mu sync.Mutex
	t  time.Time
}

// NewClock crea el reloj en t.
func NewClock(t time.Time) *Clock { return &Clock{t: t} }

// Now devuelve la hora simulada.
func (c *Clock) Now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.t }

// Set fija la hora.
func (c *Clock) Set(t time.Time) { c.mu.Lock(); c.t = t; c.mu.Unlock() }

// Advance avanza la hora.
func (c *Clock) Advance(d time.Duration) { c.mu.Lock(); c.t = c.t.Add(d); c.mu.Unlock() }

// MX arma un instante en America/Mexico_City.
func MX(y int, m time.Month, d, hh, mm int) time.Time {
	return time.Date(y, m, d, hh, mm, 0, 0, dominio.MX)
}

// NewDB abre una base temporal, aplica todas las migraciones (goose como biblioteca) y siembra las fixtures mínimas.
func NewDB(t testing.TB) (*sql.DB, Fixtures) {
	t.Helper()
	d, err := dominio.Open(filepath.Join(TempDir(t), "prueba.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	ctx := context.Background()
	if _, err := dominio.Migrate(ctx, d); err != nil {
		t.Fatal(err)
	}
	fx, err := Seed(ctx, d)
	if err != nil {
		t.Fatal(err)
	}
	return d, fx
}

// Seed siembra las 3 personas (Axel titular, Ximena co-titular, Karen externa) y las 2 cuentas; Karen no tiene cuenta.
func Seed(ctx context.Context, d *sql.DB) (Fixtures, error) {
	var fx Fixtures
	err := dominio.WithTx(ctx, d, func(tx *sql.Tx) error {
		for _, p := range []struct {
			code, name, part string
			id               *int64
		}{{"axel", "Axel", "holder", &fx.AxelPersonID}, {"ximena", "Ximena", "co_holder", &fx.XimenaPersonID}, {"karen", "Karen", "external", &fx.KarenPersonID}} {
			if err := tx.QueryRowContext(ctx, `INSERT INTO person (code, display_name, participation) VALUES (?, ?, ?) RETURNING id`, p.code, p.name, p.part).Scan(p.id); err != nil {
				return err
			}
		}
		for _, u := range []struct {
			user     string
			personID int64
			id       *int64
		}{{AxelUser, fx.AxelPersonID, &fx.AxelUserID}, {XimenaUser, fx.XimenaPersonID, &fx.XimenaUserID}} {
			h, err := dominio.HashPassword(FakeCredential, FastArgon2)
			if err != nil {
				return err
			}
			if err := tx.QueryRowContext(ctx, `INSERT INTO user_account (person_id, username, password_hash) VALUES (?, ?, ?) RETURNING id`, u.personID, u.user, h).Scan(u.id); err != nil {
				return err
			}
		}
		return nil
	})
	return fx, err
}
