package dominio_test

import (
	"context"
	"database/sql"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

// Repositorios de comidas y salud (agente api-comidas-salud) con datos ficticios.
func TestComidasSaludRepositorios(t *testing.T) {
	d, fx := testutil.NewDB(t)
	ctx := context.Background()
	if _, err := d.Exec(`
INSERT INTO product_group (id, name, purchase_channel, sort_order) VALUES (1, 'Grupo de prueba', 'chain', 1);
INSERT INTO product (id, name, product_group_id, base_unit) VALUES (1, 'Producto de prueba', 1, 'g');
INSERT INTO presentation (id, product_id, name, size_milli) VALUES (1, 1, 'Bolsa 1 kg', 1000000);
INSERT INTO store (id, kind, name) VALUES (1, 'chain', 'Tienda de prueba');
INSERT INTO store_product_group (store_id, product_group_id) VALUES (1, 1);
INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source) VALUES (1, 1, 2500, '2026-10-01', 'manual'), (1, 1, 2600, '2026-11-01', 'manual');
INSERT INTO nutrition_fact (product_id, kcal_milli_per_100g, source, is_preferred) VALUES (1, 100000, 'manual', 1);`); err != nil {
		t.Fatal(err)
	}
	slots, err := dominio.ListMealSlots(ctx, d, "ORDER BY ms.sort_order")
	if err != nil || len(slots) != 4 || slots[3].Code != "snack" {
		t.Fatalf("%v %+v", err, slots)
	}
	ps, err := dominio.ListMealPersons(ctx, d)
	if err != nil || len(ps) != 3 || !ps[0].HasMealPlan() || ps[2].HasMealPlan() {
		t.Fatalf("%v %+v", err, ps)
	}
	// Recetas
	r := dominio.RecipeRow{Name: "Receta de prueba", BaseServings: 2, Origin: "household", IsActive: true, SlotIDs: []int64{2, 1},
		Ingredients: []dominio.RecipeIngredientRow{{ProductID: 1, PresentationID: sql.NullInt64{Int64: 1, Valid: true}, QuantityMilli: 100000, SortOrder: 1}}}
	if err := dominio.SaveRecipe(ctx, d, &r, true, true); err != nil || r.ID == 0 || r.Ingredients[0].ID == 0 {
		t.Fatalf("%v %+v", err, r)
	}
	r.Name, r.SlotIDs = "Receta de prueba 2", []int64{3}
	if err := dominio.SaveRecipe(ctx, d, &r, true, false); err != nil {
		t.Fatal(err)
	}
	got, err := dominio.GetRecipeRow(ctx, d, r.ID)
	if err != nil || got.Name != "Receta de prueba 2" || len(got.SlotIDs) != 1 || len(got.Ingredients) != 1 {
		t.Fatalf("%v %+v", err, got)
	}
	if _, err := dominio.GetRecipeRow(ctx, d, 999); err == nil {
		t.Fatal("receta inexistente debe ser 404")
	}
	if all, err := dominio.ListRecipeRows(ctx, d); err != nil || len(all) != 1 {
		t.Fatalf("%v %+v", err, all)
	}
	// Plan
	e := dominio.MealPlanEntryRow{PlanDate: "2026-10-12", MealSlotID: 3, PersonID: fx.AxelPersonID, RecipeID: sql.NullInt64{Int64: r.ID, Valid: true}, ServingsMilli: sql.NullInt64{Int64: 1000, Valid: true}}
	if err := dominio.SaveMealPlanEntry(ctx, d, &e); err != nil || e.ID == 0 {
		t.Fatal(err)
	}
	e.ServingsMilli.Int64 = 2000
	if err := dominio.SaveMealPlanEntry(ctx, d, &e); err != nil {
		t.Fatal(err)
	}
	if x, err := dominio.GetMealPlanEntry(ctx, d, e.ID); err != nil || x.ServingsMilli.Int64 != 2000 {
		t.Fatalf("%v %+v", err, x)
	}
	if _, err := dominio.GetMealPlanEntry(ctx, d, 999); err == nil {
		t.Fatal("entrada inexistente debe ser 404")
	}
	if es, err := dominio.ListMealPlanEntries(ctx, d, "2026-10-12", "2026-10-18"); err != nil || len(es) != 1 {
		t.Fatalf("%v %+v", err, es)
	}
	if cell, err := dominio.MealPlanCell(ctx, d, "2026-10-12", 3, fx.AxelPersonID); err != nil || len(cell) != 1 {
		t.Fatalf("%v %+v", err, cell)
	}
	if ds, err := dominio.MealPlanDates(ctx, d, "2026-12-31"); err != nil || len(ds) != 1 {
		t.Fatalf("%v %+v", err, ds)
	}
	if n, err := dominio.RecipeUseCounts(ctx, d, "2026-10-01", "2026-10-31"); err != nil || n[r.ID] != 1 {
		t.Fatalf("%v %+v", err, n)
	}
	used, next, err := dominio.RecipePlanInfo(ctx, d, r.ID, "2026-10-12")
	if err != nil || !used || next.String != "2026-10-12" {
		t.Fatalf("%v %v %+v", err, used, next)
	}
	if err := dominio.DeleteMealPlanEntry(ctx, d, e.ID); err != nil {
		t.Fatal(err)
	}
	if err := dominio.ArchiveRecipe(ctx, d, r.ID); err != nil {
		t.Fatal(err)
	}
	if err := dominio.DeleteRecipe(ctx, d, r.ID); err != nil {
		t.Fatal(err)
	}
	// Salud
	h := dominio.HealthProfileRow{PersonID: fx.XimenaPersonID, HeightMM: sql.NullInt64{Int64: 1600, Valid: true}}
	if err := dominio.SaveHealthProfile(ctx, d, &h); err != nil || h.ID == 0 {
		t.Fatal(err)
	}
	h.Goal = sql.NullString{String: "gain", Valid: true}
	if err := dominio.SaveHealthProfile(ctx, d, &h); err != nil {
		t.Fatal(err)
	}
	if x, err := dominio.GetHealthProfileRow(ctx, d, h.ID); err != nil || x.Goal.String != "gain" {
		t.Fatalf("%v %+v", err, x)
	}
	if _, err := dominio.GetHealthProfileRow(ctx, d, 999); err == nil {
		t.Fatal("perfil inexistente debe ser 404")
	}
	if hs, err := dominio.ListHealthProfileRows(ctx, d, ""); err != nil || len(hs) != 1 {
		t.Fatalf("%v %+v", err, hs)
	}
	w := dominio.WeightEntryRow{PersonID: fx.XimenaPersonID, MeasuredDate: "2026-10-11", WeightG: 60000, IdempotencyKey: sql.NullString{String: "clave_de_prueba_0001", Valid: true}}
	if err := dominio.SaveWeightEntry(ctx, d, &w); err != nil || w.ID == 0 {
		t.Fatal(err)
	}
	w.WeightG = 59000
	if err := dominio.SaveWeightEntry(ctx, d, &w); err != nil {
		t.Fatal(err)
	}
	if x, ok, err := dominio.WeightEntryByKey(ctx, d, "clave_de_prueba_0001"); err != nil || !ok || x.WeightG != 59000 {
		t.Fatalf("%v %v %+v", err, ok, x)
	}
	if _, ok, err := dominio.WeightEntryByKey(ctx, d, "otra_clave_de_prueba"); err != nil || ok {
		t.Fatalf("%v %v", err, ok)
	}
	if ws, err := dominio.ListWeightEntryRows(ctx, d, dominio.WeightFilter{PersonIDs: []int64{fx.XimenaPersonID, fx.AxelPersonID}, From: "2026-10-01", To: "2026-10-31"}, ""); err != nil || len(ws) != 1 {
		t.Fatalf("%v %+v", err, ws)
	}
	if _, err := dominio.GetWeightEntryRow(ctx, d, 999); err == nil {
		t.Fatal("peso inexistente debe ser 404")
	}
	if err := dominio.DeleteWeightEntry(ctx, d, w.ID); err != nil {
		t.Fatal(err)
	}
	// Datos de los motores: el precio de noviembre (futuro respecto a hoy) no entra.
	fd, err := dominio.LoadMealFoodData(ctx, d, "2026-10-12")
	if err != nil || len(fd.Groups) != 1 || len(fd.Products) != 1 || len(fd.Presentations) != 1 || len(fd.Nutrition) != 1 ||
		len(fd.Stores) != 1 || len(fd.Stores[0].GroupIDs) != 1 || len(fd.Prices) != 1 || len(fd.Leftovers) != 0 {
		t.Fatalf("%v %+v", err, fd)
	}
}
