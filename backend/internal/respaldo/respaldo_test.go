package respaldo_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/respaldo"
	"finanzas-hogar/internal/testutil"
)

func copyFile(t *testing.T, src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		t.Fatal(err)
	}
}

// S6-09: VACUUM INTO en carpeta local, integrity_check = ok, restauración en una ruta temporal con integrity_check = ok
// y conteos por tabla iguales a los de la base de origen.
func TestRespaldo_VacuumIntoIntegridadYRestauracion(t *testing.T) {
	d, fx := testutil.NewDB(t)
	ctx := context.Background()
	var cat int64
	d.QueryRow(`INSERT INTO category (code, name) VALUES ('groceries', 'Despensa') RETURNING id`).Scan(&cat)
	for i := 0; i < 3; i++ {
		if _, err := d.Exec(`INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, payment_method, paid_by_person_id) VALUES ('variable', 'Súper', ?, 10000, '2026-10-05', 'cash', ?)`, cat, fx.AxelPersonID); err != nil {
			t.Fatal(err)
		}
	}
	dir := filepath.Join(testutil.TempDir(t), "respaldos")
	res, err := respaldo.Run(ctx, d, dir, testutil.TempDir(t), "scheduled", testutil.MX(2026, time.October, 13, 3, 15))
	if err != nil || res.Integrity != "ok" || res.RestoreResult != "ok" || len(res.Diff) != 0 {
		t.Fatalf("%+v %v", res, err)
	}
	if filepath.Dir(res.Path) != dir || !strings.HasPrefix(filepath.Base(res.Path), "finanzas-20261013-031500-scheduled-") || len(res.SHA256) != 64 {
		t.Errorf("archivo versionado: %s", res.Path)
	}
	dst := filepath.Join(testutil.TempDir(t), "restaura.db")
	copyFile(t, res.Path, dst)
	if ic, err := dominio.IntegrityCheck(ctx, dst); ic != "ok" || err != nil {
		t.Fatalf("integrity_check de la restauración: %q %v", ic, err)
	}
	rd, err := dominio.OpenReadOnly(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer rd.Close()
	got, _ := dominio.TableCounts(ctx, rd)
	want, _ := dominio.TableCounts(ctx, d)
	if diff := respaldo.DiffCounts(want, got); len(diff) != 0 || len(got) < 50 {
		t.Fatalf("conteos: %v (%d tablas)", diff, len(got))
	}
	t.Logf("restauración: integrity_check = ok; %d tablas con el mismo conteo que el origen (expense %d, person %d, user_account %d, app_setting %d, backup %d)",
		len(got), got["expense"], got["person"], got["user_account"], got["app_setting"], got["backup"])
}

func TestRespaldo_RestaurarConservaCatalogoYSesiones(t *testing.T) {
	d, fx := testutil.NewDB(t)
	ctx := context.Background()
	dir, tmp := filepath.Join(testutil.TempDir(t), "respaldos"), testutil.TempDir(t)
	now := testutil.MX(2026, time.October, 12, 12, 0)
	res, err := respaldo.Run(ctx, d, dir, tmp, "manual", now)
	if err != nil {
		t.Fatal(err)
	}
	d.Exec(`UPDATE person SET display_name = 'Cambiada' WHERE id = ?`, fx.KarenPersonID)
	_, h, _ := dominio.NewSessionToken()
	sid, err := dominio.CreateSession(ctx, d, fx.AxelUserID, h, "Teléfono", "go-test", "127.0.0.1", now, 90*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	out, err := respaldo.Restore(ctx, d, res.ID, dir, tmp, now.Add(time.Minute))
	if err != nil || out.Integrity != "ok" || out.PreRestoreID == 0 {
		t.Fatalf("%+v %v", out, err)
	}
	var name string
	d.QueryRow(`SELECT display_name FROM person WHERE id = ?`, fx.KarenPersonID).Scan(&name)
	if name != "Karen" {
		t.Errorf("no se restauró: %q", name)
	}
	if _, err := dominio.ActiveSession(ctx, d, h, now.Add(2*time.Minute)); err != nil {
		t.Errorf("la sesión %d abierta después del respaldo se perdió: %v", sid, err)
	}
	var n int
	d.QueryRow(`SELECT COUNT(*) FROM backup WHERE id IN (?, ?) AND status = 'ok'`, res.ID, out.PreRestoreID).Scan(&n)
	if n != 2 {
		t.Errorf("catálogo de respaldos: %d", n)
	}
	var de *dominio.Error
	if _, err := respaldo.Restore(ctx, d, 999, dir, tmp, now); !errors.As(err, &de) || de.Kind != dominio.KindNotFound {
		t.Errorf("no existe: %v", err)
	}
	if _, err := dominio.PruneBackups(ctx, d, "manual", 0, now); err != nil {
		t.Fatal(err)
	}
	if _, err := respaldo.Restore(ctx, d, res.ID, dir, tmp, now); !errors.As(err, &de) || de.Kind != dominio.KindConflict {
		t.Errorf("ya borrado por la retención: %v", err)
	}
}
