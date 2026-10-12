package profeco

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

const migrationsDir = "../../../db/migrations"

// newSchemaDB crea una base en memoria con el esquema real (migraciones 00001 a
// 00003, parte «Up»). Si el esquema cambia, la prueba lo nota.
func newSchemaDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "app.db")+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	for _, name := range []string{"00001_esquema.sql", "00002_puerta1.sql", "00003_puerta1b.sql"} {
		b, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			t.Fatal(err)
		}
		src := string(b)
		if i := strings.Index(src, "-- +goose Down"); i >= 0 {
			src = src[:i]
		}
		if _, err := db.Exec(src); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
	return db
}

func seedCatalog(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, q := range []string{
		`INSERT INTO product (id, name, base_unit) VALUES (1, 'Atún', 'g'), (2, 'Arroz', 'g')`,
		`INSERT INTO presentation (id, product_id, name, brand, size_milli) VALUES (7, 1, 'Lata 140 g', 'Dolores', 140000), (8, 2, 'Bolsa 900 g', 'SOS', 900000)`,
		`INSERT INTO store (id, kind, name, profeco_name) VALUES (1, 'chain', 'Walmart', 'Wal-mart'), (2, 'chain', 'Bodega Aurrera', 'Bodega Aurrera')`,
		`INSERT INTO store (id, kind, parent_id, name, profeco_name, municipality, state) VALUES (3, 'branch', 1, 'Tepeyac', 'Walmart Sucursal Tepeyac', 'Gustavo a. Madero', 'Ciudad de México')`,
		`INSERT INTO store (id, kind, name, profeco_name, is_active) VALUES (4, 'chain', 'Inactiva', 'Chedraui', 0)`,
		`INSERT INTO store (id, kind, name, visit_weekday) VALUES (5, 'market', 'Sin PROFECO', 3)`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
}

func TestProfecoCargaEnBase(t *testing.T) {
	db := newSchemaDB(t)
	seedCatalog(t, db)
	ctx := context.Background()
	body := header +
		row("Atún", "Lata 140 g", "Dolores", "21.5", "2026/03/02", "Wal-mart", "Walmart Sucursal Tepeyac", "Calz", "Ciudad de México", "Gustavo a. Madero") +
		row("Atún", "Lata 140 g", "Dolores", "20", "2026/03/02", "Wal-mart", "Walmart Sucursal Tepeyac", "Calz", "Ciudad de México", "Gustavo a. Madero") + // duplicado del mismo día
		row("Atún", "Lata 140 g", "Dolores", "22", "2026/03/09", "Wal-mart", "Walmart Sucursal Tepeyac", "Calz", "Ciudad de México", "Gustavo a. Madero") +
		row("Arroz", "Bolsa 900 g", "SOS", "19", "2026/03/02", "Wal-mart", "Otra", "Calz", "Ciudad de México", "Gustavo a. Madero") + // propuesta: no cuenta en producción
		row("Atún", "Lata 140 g", "Dolores", "25", "2026/03/02", "Chedraui", "Chedraui", "Calz", "Ciudad de México", "Gustavo a. Madero") // cadena inactiva: no pasa
	p := writeTmp(t, "db.csv", body)
	links := mustLinks(t, "Atún,Lata 140 g,Dolores,7,confirmado,axel,2026-10-11\nArroz,Bolsa 900 g,SOS,8,propuesto,,\n")

	res, err := LoadIntoDB(ctx, db, Options{Corte: "2026", Links: links}, p)
	if err != nil {
		t.Fatal(err)
	}
	if res.Read != 5 || res.Loaded != 2 || res.Duplicates != 1 || res.Skipped != 1 || res.FilteredOut != 1 || !res.Balanced() {
		t.Fatalf("resultado: %+v", res)
	}
	var n int
	var src, ref string
	var promo int
	if err := db.QueryRow(`SELECT count(*) FROM price WHERE source = 'profeco'`).Scan(&n); err != nil || n != 2 {
		t.Fatalf("precios cargados: %d %v", n, err)
	}
	if err := db.QueryRow(`SELECT source, is_promotional, source_ref FROM price WHERE observed_date = '2026-03-02'`).Scan(&src, &promo, &ref); err != nil {
		t.Fatal(err)
	}
	if src != "profeco" || promo != 0 || ref != "QQP 2026 db.csv:1" {
		t.Errorf("fuente %q promo %d ref %q", src, promo, ref)
	}
	var cents int64
	var store int64
	if err := db.QueryRow(`SELECT price_cents, store_id FROM price WHERE observed_date = '2026-03-02'`).Scan(&cents, &store); err != nil || cents != 2150 || store != 3 {
		t.Errorf("centavos %d tienda %d %v", cents, store, err)
	}

	// Volver a cargar el mismo archivo no duplica nada: todo cuenta como duplicado.
	res, err = LoadIntoDB(ctx, db, Options{Corte: "2026", Links: links}, p)
	if err != nil || res.Loaded != 0 || res.Duplicates != 3 {
		t.Fatalf("segunda carga: %+v %v", res, err)
	}

	// Un archivo que no cabe en el límite de app_setting revierte todo.
	if _, err := db.Exec(`UPDATE app_setting SET value_json = '1' WHERE setting_key = 'price.profeco_import_max_gb'`); err != nil {
		t.Fatal(err)
	}
	if mb, err := MaxBytesFromSettings(ctx, db); err != nil || mb != 1<<30 {
		t.Errorf("límite: %d %v", mb, err)
	}
	if _, err := LoadIntoDB(ctx, db, Options{Links: links, MaxBytes: 5}, p); err == nil {
		t.Error("el límite de tamaño debía fallar")
	}
}

func TestProfecoLimiteDeAjustes(t *testing.T) {
	db := newSchemaDB(t)
	ctx := context.Background()
	if _, err := db.Exec(`UPDATE app_setting SET value_json = '2' WHERE setting_key = 'price.profeco_import_max_gb'`); err != nil {
		t.Fatal(err)
	}
	if mb, err := MaxBytesFromSettings(ctx, db); err != nil || mb != 2<<30 {
		t.Errorf("2 GB: %d %v", mb, err)
	}
	if _, err := db.Exec(`DELETE FROM app_setting WHERE setting_key = 'price.profeco_import_max_gb'`); err != nil {
		t.Fatal(err)
	}
	if _, err := MaxBytesFromSettings(ctx, db); err == nil {
		t.Error("sin la clave debía fallar (nunca una constante en código)")
	}
	// LoadIntoDB sin la clave y sin MaxBytes falla antes de leer nada.
	if _, err := LoadIntoDB(ctx, db, Options{}, "no-importa.csv"); err == nil {
		t.Error("LoadIntoDB sin límite debía fallar")
	}
}

func TestProfecoQQP(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "qqp.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	q, err := NewQQP(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	im, _ := NewImporter(Options{Corte: "2026", Records: q})
	res, err := im.ImportFile(ctx, muestraPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := q.Flush(); err != nil {
		t.Fatal(err)
	}
	var n int64
	var maxDate, maxPrecio string
	if err := db.QueryRow(`SELECT count(*), max(fecha), max(precio_cents) FROM profeco_record`).Scan(&n, &maxDate, &maxPrecio); err != nil {
		t.Fatal(err)
	}
	// El criterio S5-05: lo que hay en la base = renglones del archivo − encabezado.
	if n != res.Read-res.Rejected || n != q.Rows || n != res.Lines-1 {
		t.Errorf("base %d; leídos %d rechazados %d wc-l %d", n, res.Read, res.Rejected, res.Lines)
	}
	if maxDate != res.MaxDate {
		t.Errorf("max fecha %s != %s", maxDate, res.MaxDate)
	}
	// Un segundo corte con el mismo archivo choca con UNIQUE (archivo, renglón).
	im2, _ := NewImporter(Options{Records: q})
	if _, err := im2.ImportFile(ctx, muestraPath); err == nil {
		t.Error("cargar dos veces el mismo archivo debía chocar con la llave única")
	}
	_ = q.Flush()
	if err := (&QQP{}).Flush(); err != nil {
		t.Error("Flush sin lote debía ser inocuo")
	}
}
