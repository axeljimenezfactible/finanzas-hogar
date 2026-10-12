package dominio_test

import (
	"context"
	"database/sql"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

// Repositorios de api-compras-b: catálogo, precios por llave natural, recetas en borrador, inventario y nutrición.
func TestComprasBRepositorios(t *testing.T) {
	ctx := context.Background()
	db, _ := testutil.NewDB(t)
	for _, q := range []string{
		`INSERT INTO store (id, kind, name) VALUES (1, 'chain', 'Cadena')`,
		`INSERT INTO store (id, kind, parent_id, name, municipality) VALUES (2, 'branch', 1, 'Sucursal', 'Cuauhtémoc')`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	tx := func(f func(tx *sql.Tx) error) error { return dominio.WithTx(ctx, db, f) }
	item := dominio.CBCatalogItem{Name: "Aceite", BaseUnit: "ml", Presentations: []dominio.CBCatalogPresentation{{Name: "500 ml", SizeMilli: 500000}}}
	var st string
	var pid int64
	var pres []int64
	for i, wantSt := range []string{dominio.CBRowCreated, dominio.CBRowUnchanged} {
		if err := tx(func(tx *sql.Tx) (err error) {
			st, pid, pres, err = dominio.CBUpsertCatalogItem(ctx, tx, item, nil, "req")
			return err
		}); err != nil || st != wantSt || len(pres) != 1 {
			t.Fatalf("vuelta %d: %s %d %v %v", i, st, pid, pres, err)
		}
	}
	g := int64(1000)
	item.GramsPerUnitMilli = &g
	item.Presentations = append(item.Presentations, dominio.CBCatalogPresentation{Name: "1 l", SizeMilli: 1000000, Barcode: dominio.CBStrPtr("75010000000")})
	if err := tx(func(tx *sql.Tx) (err error) {
		st, _, pres, err = dominio.CBUpsertCatalogItem(ctx, tx, item, nil, "req")
		return err
	}); err != nil || st != dominio.CBRowUpdated || len(pres) != 2 {
		t.Fatalf("%s %v %v", st, pres, err)
	}
	item.Presentations[0].SizeMilli = 450000
	if err := tx(func(tx *sql.Tx) (err error) {
		st, _, _, err = dominio.CBUpsertCatalogItem(ctx, tx, item, nil, "req")
		return err
	}); err != nil || st != dominio.CBRowUpdated {
		t.Fatalf("%s %v", st, err)
	}
	for _, bad := range []dominio.CBCatalogItem{
		{Name: "aceite", BaseUnit: "g"},
		{Name: "Sal", BaseUnit: "g", Presentations: []dominio.CBCatalogPresentation{{Name: "Bolsa", SizeMilli: 1, Barcode: dominio.CBStrPtr("75010000000")}}},
	} {
		err := tx(func(tx *sql.Tx) error {
			_, _, _, err := dominio.CBUpsertCatalogItem(ctx, tx, bad, nil, "req")
			return err
		})
		if e, ok := dominio.AsError(err); !ok || e.Code != "validation_failed" {
			t.Fatalf("%s: %v", bad.Name, err)
		}
	}
	if owner, err := dominio.CBPresentationOf(ctx, db, pres[0]); err != nil || owner != pid {
		t.Fatal(owner, err)
	}
	p := dominio.CBPriceRow{PresentationID: pres[0], StoreID: 2, PriceCents: 100, ObservedDate: "2026-10-01", Source: "site_assisted"}
	for i, want := range []string{dominio.CBRowCreated, dominio.CBRowUnchanged, dominio.CBRowUpdated} {
		if i == 2 {
			p.PriceCents = 120
		}
		if err := tx(func(tx *sql.Tx) (err error) {
			st, _, err = dominio.CBUpsertSourcedPrice(ctx, tx, p, nil, "req")
			return err
		}); err != nil || st != want {
			t.Fatalf("precio %d: %s %v", i, st, err)
		}
	}
	p.Source = "profeco"
	if _, err := dominio.CBInsertPrice(ctx, db, p); err != nil {
		t.Fatal(err)
	}
	cov, total, last, err := dominio.CBProfecoStats(ctx, db)
	if err != nil || total != 1 || len(cov) != 1 || cov[0].StoreID != 1 || !last.Valid {
		t.Fatalf("%+v %d %v", cov, total, err)
	}
	ps, err := dominio.CBQueryPrices(ctx, db, "WHERE p.store_id = ?", []any{2}, "ORDER BY p.id")
	if err != nil || len(ps) != 2 || ps[0].BaseUnit != "ml" || ps[0].SizeMilli != 450000 {
		t.Fatalf("%+v %v", ps, err)
	}
	if _, err := dominio.CBGetPrice(ctx, db, 999); err == nil {
		t.Fatal("precio inexistente")
	}
	if ok, err := dominio.CBExists(ctx, db, "store", 2); !ok || err != nil {
		t.Fatal(ok, err)
	}
	rec := dominio.CBDraftRecipe{Name: "Ensalada", BaseServings: 2, MealSlotIDs: []int64{1},
		Ingredients: []dominio.CBDraftIngredient{{ProductID: pid, QuantityMilli: 1000, PresentationID: &pres[0]}}}
	for i, wantErr := range []bool{false, true} {
		err := tx(func(tx *sql.Tx) error { _, err := dominio.CBInsertDraftRecipe(ctx, tx, rec, nil, "req"); return err })
		if (err != nil) != wantErr {
			t.Fatalf("receta %d: %v", i, err)
		}
	}
	if _, err := db.Exec(`INSERT INTO nutrition_fact (product_id, kcal_milli_per_100g, source) VALUES (?, 884000, 'manual')`, pid); err != nil {
		t.Fatal(err)
	}
	if ns, err := dominio.CBQueryNutrition(ctx, db, "WHERE product_id = ?", []any{pid}, ""); err != nil || len(ns) != 1 || ns[0].KcalMilli != 884000 {
		t.Fatalf("%+v %v", ns, err)
	}
	if _, err := db.Exec(`INSERT INTO inventory_item (id, product_id, tracking_mode, quantity_milli) VALUES (5, ?, 'stock', 1000)`, pid); err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.CBInsertMovement(ctx, db, 5, "2026-10-12", "adjustment", dominio.CBStrPtr("count"), 1000, nil); err != nil {
		t.Fatal(err)
	}
	if it, err := dominio.CBGetInventory(ctx, db, 5); err != nil || it.Name != "Aceite" {
		t.Fatalf("%+v %v", it, err)
	}
	if _, err := dominio.CBGetInventory(ctx, db, 99); err == nil {
		t.Fatal("inventario inexistente")
	}
	if ms, err := dominio.CBQueryMovements(ctx, db, 5, "ORDER BY m.id"); err != nil || len(ms) != 1 || ms[0].Delta != 1000 {
		t.Fatalf("%+v %v", ms, err)
	}
}

// Un error de la base (contexto cancelado) se propaga sin tragarse en cada repositorio del área.
func TestComprasBRepositoriosErrores(t *testing.T) {
	db, _ := testutil.NewDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	checks := map[string]error{}
	_, checks["prices"] = dominio.CBQueryPrices(ctx, db, "", nil, "")
	_, checks["price"] = dominio.CBGetPrice(ctx, db, 1)
	_, checks["insertPrice"] = dominio.CBInsertPrice(ctx, db, dominio.CBPriceRow{})
	_, _, _, checks["profeco"] = dominio.CBProfecoStats(ctx, db)
	_, checks["nutrition"] = dominio.CBQueryNutrition(ctx, db, "", nil, "")
	_, checks["inventory"] = dominio.CBQueryInventory(ctx, db, "")
	_, checks["inventoryOne"] = dominio.CBGetInventory(ctx, db, 1)
	_, checks["movements"] = dominio.CBQueryMovements(ctx, db, 1, "")
	_, checks["insertMovement"] = dominio.CBInsertMovement(ctx, db, 1, "2026-10-12", "adjustment", nil, 1, nil)
	_, checks["presentation"] = dominio.CBPresentationOf(ctx, db, 1)
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	_, _, _, checks["catalog"] = dominio.CBUpsertCatalogItem(ctx, tx, dominio.CBCatalogItem{Name: "x"}, nil, "req")
	_, _, checks["sourced"] = dominio.CBUpsertSourcedPrice(ctx, tx, dominio.CBPriceRow{}, nil, "req")
	_, checks["recipe"] = dominio.CBInsertDraftRecipe(ctx, tx, dominio.CBDraftRecipe{Name: "x"}, nil, "req")
	checks["audit"] = dominio.CBAuditWrite(ctx, tx, nil, "insert", "price", 1, nil, "req")
	for name, err := range checks {
		if err == nil {
			t.Errorf("%s: se esperaba el error de contexto cancelado", name)
		}
	}
}
