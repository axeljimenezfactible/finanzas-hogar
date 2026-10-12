package comidas_test

import (
	"fmt"
	"strings"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comidas"
	"finanzas-hogar/internal/testutil"
)

// Catálogo mínimo ficticio: arroz 1 kg a $30.00 y frijol 1 kg a $40.00 (precio manual del 10-oct-2026),
// 360 y 340 kcal por 100 g, y la receta 1 (2 porciones: 200 g de arroz y 100 g de frijol).
// Costo por porción esperado = (200/1000 × 3000 + 100/1000 × 4000) / 2 = 500 centavos; kcal = (720 + 340) / 2 = 530.
const fixture = `
INSERT INTO product_group (id, name, purchase_channel, sort_order) VALUES (1, 'Abarrotes de prueba', 'chain', 1);
INSERT INTO product (id, name, product_group_id, base_unit) VALUES (1, 'Arroz de prueba', 1, 'g'), (2, 'Frijol de prueba', 1, 'g');
INSERT INTO presentation (id, product_id, name, size_milli) VALUES (1, 1, 'Arroz 1 kg', 1000000), (2, 2, 'Frijol 1 kg', 1000000);
INSERT INTO store (id, kind, name) VALUES (1, 'chain', 'Tienda de prueba');
INSERT INTO store_product_group (store_id, product_group_id) VALUES (1, 1);
INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source) VALUES (1, 1, 3000, '2026-10-10', 'manual'), (2, 1, 4000, '2026-10-10', 'manual');
INSERT INTO nutrition_fact (product_id, kcal_milli_per_100g, source, source_ref, is_preferred) VALUES (1, 360000, 'usda_fdc', 'FDC-PRUEBA-1', 1), (2, 340000, 'usda_fdc', 'FDC-PRUEBA-2', 1);
INSERT INTO recipe (id, name, instructions, base_servings, origin) VALUES (1, 'Arroz con frijoles', 'Cocer y servir.', 2, 'household');
INSERT INTO recipe_meal_slot (recipe_id, meal_slot_id) VALUES (1, 2), (1, 3);
INSERT INTO recipe_ingredient (recipe_id, product_id, presentation_id, quantity_milli, sort_order) VALUES (1, 1, 1, 200000, 1), (1, 2, 2, 100000, 2);
`

func setup(t *testing.T) (*testutil.Server, *testutil.Client) {
	t.Helper()
	s := testutil.NewServer(t, comidas.Register)
	if _, err := s.DB.Exec(fixture); err != nil {
		t.Fatal(err)
	}
	return s, s.SessionFor(t, testutil.AxelUser)
}

var keyN int

func key() []string {
	keyN++
	return []string{"Idempotency-Key", fmt.Sprintf("prueba_comidas_%06d", keyN)}
}

func want(t *testing.T, res *testutil.Response, status int) *testutil.Response {
	t.Helper()
	if res.StatusCode != status {
		t.Fatalf("estado %d, se esperaba %d: %s", res.StatusCode, status, res.Body)
	}
	return res
}

func wantCode(t *testing.T, res *testutil.Response, status int, code string) {
	t.Helper()
	want(t, res, status)
	if !strings.Contains(string(res.Body), `"`+code+`"`) {
		t.Fatalf("se esperaba el código %s: %s", code, res.Body)
	}
}

func plan(t *testing.T, c *testutil.Client, date string, persons []int64, body map[string]any) api.MealPlanEntryList {
	t.Helper()
	b := map[string]any{"plan_dates": []string{date}, "meal_slot_id": 2, "person_ids": persons, "is_eating_out": false, "recipe_id": 1}
	for k, v := range body {
		b[k] = v
	}
	var out api.MealPlanEntryList
	want(t, c.Do("POST", "/api/v1/meal-plan-entries", b, key()...), 201).JSON(t, &out)
	return out
}

func TestOp_listMealSlots(t *testing.T) {
	s, c := setup(t)
	var pg api.MealSlotPage
	want(t, c.Do("GET", "/api/v1/meal-slots", nil), 200).JSON(t, &pg)
	if len(pg.Items) != 4 || pg.Items[0].Code != "breakfast" || !pg.Items[3].IsOptional {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/meal-slots?limit=201", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/meal-slots", nil), 401)
}

func TestOp_createMealPlanEntries(t *testing.T) {
	s, c := setup(t)
	out := plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID, s.Fx.XimenaPersonID}, nil)
	if len(out.Entries) != 2 || *out.Entries[0].ServingsMilli != 1000 {
		t.Fatalf("%+v", out)
	}
	fuera := map[string]any{"plan_dates": []string{"2026-10-12"}, "meal_slot_id": 2, "person_ids": []int64{s.Fx.AxelPersonID}, "is_eating_out": true}
	wantCode(t, c.Do("POST", "/api/v1/meal-plan-entries", fuera, key()...), 409, "slot_conflict")
	karen := map[string]any{"plan_dates": []string{"2026-10-13"}, "meal_slot_id": 2, "person_ids": []int64{s.Fx.KarenPersonID}, "is_eating_out": true}
	wantCode(t, c.Do("POST", "/api/v1/meal-plan-entries", karen, key()...), 422, "external_person")
	sinReceta := map[string]any{"plan_dates": []string{"2026-10-13"}, "meal_slot_id": 2, "person_ids": []int64{s.Fx.AxelPersonID}, "is_eating_out": false}
	wantCode(t, c.Do("POST", "/api/v1/meal-plan-entries", sinReceta, key()...), 422, "validation_failed")
	want(t, c.Do("POST", "/api/v1/meal-plan-entries", fuera), 400) // sin Idempotency-Key
}

func TestOp_getMealPlanWeek(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	plan(t, c, "2026-10-13", []int64{s.Fx.XimenaPersonID}, map[string]any{"is_eating_out": true, "recipe_id": nil})
	var wk api.MealPlanWeek
	want(t, c.Do("GET", "/api/v1/meal-plan/week", nil), 200).JSON(t, &wk) // 12-oct-2026 es lunes: semana por defecto W42
	if wk.WeekStart != "2026-10-12" || wk.WeekEnd != "2026-10-18" || len(wk.Cells) != 2 || len(wk.Days) != 7 {
		t.Fatalf("%+v", wk)
	}
	cell := wk.Cells[0]
	if *cell.CostCents != 500 || *cell.Kcal != 530 || cell.MealPlanEntryID == nil || *cell.PriceStoreID != 1 || wk.WeekCostCents != 500 {
		t.Fatalf("celda %+v, semana %d", cell, wk.WeekCostCents)
	}
	if wk.Days[0].Persons[0].Kcal != 530 || wk.UnplannedSlots != 21 {
		t.Fatalf("%+v %d", wk.Days[0], wk.UnplannedSlots)
	}
	want(t, c.Do("GET", fmt.Sprintf("/api/v1/meal-plan/week?period=2026-W42&person_id=%d", s.Fx.XimenaPersonID), nil), 200).JSON(t, &wk)
	if len(wk.Cells) != 1 || !wk.Cells[0].IsEatingOut || len(wk.Days[0].Persons) != 1 {
		t.Fatalf("%+v", wk)
	}
	wantCode(t, c.Do("GET", "/api/v1/meal-plan/week?period=2026-10", nil), 422, "validation_failed")
	want(t, c.Do("GET", "/api/v1/meal-plan/week?period=semana", nil), 400)
}

func TestOp_getMealPlanMonth(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	var mo api.MealPlanMonth
	want(t, c.Do("GET", "/api/v1/meal-plan/month?period=2026-10", nil), 200).JSON(t, &mo)
	if mo.Month != "2026-10" || len(mo.Days) != 31 || mo.Days[11].Status != "incomplete" || mo.Days[0].Status != "empty" {
		t.Fatalf("%+v", mo)
	}
	want(t, c.Do("GET", "/api/v1/meal-plan/month", nil), 200)
	wantCode(t, c.Do("GET", "/api/v1/meal-plan/month?period=2026-W42", nil), 422, "validation_failed")
}

func TestOp_getMealPlanTrends(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	var tr api.MealPlanTrends
	want(t, c.Do("GET", "/api/v1/meal-plan/trends", nil), 200).JSON(t, &tr)
	if len(tr.Weeks) != 1 || tr.Weeks[0].CostCents != 500 || len(tr.Weeks[0].Persons) != 2 || *tr.Weeks[0].Persons[0].AvgKcalPerDay != 76 {
		t.Fatalf("%+v", tr) // 530 / 7 = 75.71 → 76 (ROUND_HALF_UP)
	}
	want(t, c.Do("GET", "/api/v1/meal-plan/trends?person_id=x", nil), 400)
}

func TestOp_updateMealPlanEntry(t *testing.T) {
	s, c := setup(t)
	id := plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil).Entries[0].ID
	path := fmt.Sprintf("/api/v1/meal-plan-entries/%d", id)
	var e api.MealPlanEntry
	want(t, c.Do("PATCH", path, map[string]any{"servings_milli": 1500, "notes": "más arroz"}, nil...), 200).JSON(t, &e)
	if *e.ServingsMilli != 1500 || *e.Notes != "más arroz" {
		t.Fatalf("%+v", e)
	}
	e = api.MealPlanEntry{}
	want(t, c.Do("PATCH", path, map[string]any{"is_eating_out": true, "recipe_id": nil, "servings_milli": nil}), 200).JSON(t, &e)
	if !e.IsEatingOut || e.RecipeID != nil {
		t.Fatalf("%+v", e)
	}
	wantCode(t, c.Do("PATCH", path, map[string]any{"is_eating_out": false}), 422, "validation_failed")
	want(t, c.Do("PATCH", "/api/v1/meal-plan-entries/99999", map[string]any{"notes": "x"}), 404)
}

func TestOp_deleteMealPlanEntry(t *testing.T) {
	s, c := setup(t)
	id := plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil).Entries[0].ID
	path := fmt.Sprintf("/api/v1/meal-plan-entries/%d", id)
	want(t, c.Do("DELETE", path, nil), 204)
	want(t, c.Do("DELETE", path, nil), 404)
}

func TestOp_copyPreviousWeek(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-05", []int64{s.Fx.AxelPersonID, s.Fx.XimenaPersonID}, nil)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, map[string]any{"is_eating_out": true, "recipe_id": nil})
	var res api.CopyWeekResult
	want(t, c.Do("POST", "/api/v1/meal-plan/copy-previous-week?period=2026-W42", nil, key()...), 201).JSON(t, &res)
	if res.CreatedCount != 1 { // la casilla de Axel ya tenía «fuera»: solo se copia la de Ximena
		t.Fatalf("%+v", res)
	}
	want(t, c.Do("POST", "/api/v1/meal-plan/copy-previous-week?period=2026-W42", nil, key()...), 201).JSON(t, &res)
	if res.CreatedCount != 0 {
		t.Fatalf("%+v", res)
	}
	want(t, c.Do("POST", "/api/v1/meal-plan/copy-previous-week?period=2026-10", nil, key()...), 400)
}

func TestOp_listRecipes(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	var pg api.RecipePage
	want(t, c.Do("GET", "/api/v1/recipes?sort=-times_used_8w&q=arroz&meal_slot_id=2&origin=household", nil), 200).JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].CostPerServingCents != 500 || *pg.Items[0].KcalPerServing != 530 || *pg.Items[0].TimesUsed8w != 1 {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/recipes?meal_slot_id=1", nil), 200).JSON(t, &pg)
	if len(pg.Items) != 0 {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/recipes?origin=otra", nil), 400)
	want(t, c.Do("GET", "/api/v1/recipes?sort=precio", nil), 400)
}

func newRecipe(name string) map[string]any {
	return map[string]any{"name": name, "instructions": "Mezclar.", "base_servings": 4, "meal_slot_ids": []int64{1},
		"ingredients": []map[string]any{{"product_id": 1, "presentation_id": 1, "quantity_milli": 400000}}}
}

func TestOp_createRecipe(t *testing.T) {
	s, c := setup(t)
	var r api.Recipe
	want(t, c.Do("POST", "/api/v1/recipes", newRecipe("Arroz blanco"), key()...), 201).JSON(t, &r)
	if r.Origin != "household" || *r.AuthorPersonID != s.Fx.AxelPersonID || *r.CostPerServingCents != 300 || r.MealSlotIDs[0] != 1 {
		t.Fatalf("%+v", r) // 400 g × 3000 / 1000 g = 1200 centavos para 4 porciones
	}
	wantCode(t, c.Do("POST", "/api/v1/recipes", newRecipe("arroz BLANCO"), key()...), 409, "already_exists")
	wantCode(t, c.Do("POST", "/api/v1/recipes", map[string]any{"name": "Sin ingredientes"}, key()...), 422, "validation_failed")
	want(t, c.Do("POST", "/api/v1/recipes", map[string]any{"name": 5}, key()...), 400)
}

func TestOp_suggestRecipeWithAi(t *testing.T) {
	_, c := setup(t)
	wantCode(t, c.Do("POST", "/api/v1/recipes/ai-suggestions", map[string]any{"description": "algo con arroz", "meal_slot_id": 2}), 409, "ai_disabled")
	wantCode(t, c.Do("POST", "/api/v1/recipes/ai-suggestions", map[string]any{"description": " "}), 422, "validation_failed")
}

func TestOp_getRecipe(t *testing.T) {
	_, c := setup(t)
	var d api.RecipeDetail
	want(t, c.Do("GET", "/api/v1/recipes/1?servings_milli=4000", nil), 200).JSON(t, &d)
	if *d.TotalCostCents != 2000 || *d.CostPerServingCents != 500 || *d.TotalKcal != 2120 || d.Ingredients[0].QuantityMilli != 400000 ||
		*d.Ingredients[1].CostCents != 800 || *d.Ingredients[0].NutritionSource != "usda_fdc" {
		t.Fatalf("%+v", d)
	}
	want(t, c.Do("GET", "/api/v1/recipes/1", nil), 200).JSON(t, &d)
	if d.ServingsMilli != 2000 || *d.TotalCostCents != 1000 {
		t.Fatalf("%+v", d)
	}
	wantCode(t, c.Do("GET", "/api/v1/recipes/1?servings_milli=0", nil), 422, "validation_failed")
	want(t, c.Do("GET", "/api/v1/recipes/1?servings_milli=1.5", nil), 400)
	want(t, c.Do("GET", "/api/v1/recipes/999", nil), 404)
}

func TestOp_updateRecipe(t *testing.T) {
	_, c := setup(t)
	var r api.Recipe
	want(t, c.Do("PATCH", "/api/v1/recipes/1", map[string]any{"name": "Arroz con frijol negro", "notes": "rinde bien", "meal_slot_ids": []int64{2}}), 200).JSON(t, &r)
	if r.Name != "Arroz con frijol negro" || *r.Notes != "rinde bien" || len(r.MealSlotIDs) != 1 {
		t.Fatalf("%+v", r)
	}
	want(t, c.Do("PATCH", "/api/v1/recipes/1", map[string]any{"ingredients": []map[string]any{{"product_id": 2, "presentation_id": 2, "quantity_milli": 100000}}, "description": nil}), 200).JSON(t, &r)
	if *r.CostPerServingCents != 200 {
		t.Fatalf("%+v", r)
	}
	wantCode(t, c.Do("PATCH", "/api/v1/recipes/1", map[string]any{"base_servings": 0}), 422, "validation_failed")
	wantCode(t, c.Do("PATCH", "/api/v1/recipes/1", map[string]any{"ingredients": []map[string]any{{"product_id": 2, "quantity_milli": 1000}}}), 422, "validation_failed")
	want(t, c.Do("PATCH", "/api/v1/recipes/999", map[string]any{"name": "Otra receta"}), 404)
}

func TestOp_deleteRecipe(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-14", []int64{s.Fx.AxelPersonID}, nil)
	wantCode(t, c.Do("DELETE", "/api/v1/recipes/1", nil), 409, "recipe_planned")
	want(t, c.Do("DELETE", "/api/v1/recipes/1?confirm_planned=true", nil), 204)
	var r api.Recipe
	want(t, c.Do("GET", "/api/v1/recipes/1", nil), 200).JSON(t, &api.RecipeDetail{Recipe: r})
	var id int64
	if err := s.DB.QueryRow(`SELECT is_active FROM recipe WHERE id = 1`).Scan(&id); err != nil || id != 0 {
		t.Fatalf("la receta planeada se archiva: %v %d", err, id)
	}
	want(t, c.Do("POST", "/api/v1/recipes", newRecipe("Arroz solo"), key()...), 201).JSON(t, &r)
	want(t, c.Do("DELETE", fmt.Sprintf("/api/v1/recipes/%d", r.ID), nil), 204)
	want(t, c.Do("DELETE", fmt.Sprintf("/api/v1/recipes/%d", r.ID), nil), 404)
}

func TestOp_getRecipeUsage(t *testing.T) {
	s, c := setup(t)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	plan(t, c, "2026-10-20", []int64{s.Fx.AxelPersonID}, nil)
	var u api.RecipeUsage
	want(t, c.Do("GET", "/api/v1/recipes/1/usage", nil), 200).JSON(t, &u)
	if u.TimesUsed8w != 1 || *u.NextPlanDate != "2026-10-12" {
		t.Fatalf("%+v", u)
	}
	want(t, c.Do("GET", "/api/v1/recipes/999/usage", nil), 404)
}

func TestOp_duplicateRecipe(t *testing.T) {
	_, c := setup(t)
	var r api.Recipe
	want(t, c.Do("POST", "/api/v1/recipes/1/duplicate", nil, key()...), 201).JSON(t, &r)
	if r.Name != "Arroz con frijoles (copia)" || r.Origin != "household" || len(r.MealSlotIDs) != 2 || *r.CostPerServingCents != 500 {
		t.Fatalf("%+v", r)
	}
	wantCode(t, c.Do("POST", "/api/v1/recipes/1/duplicate", nil, key()...), 409, "already_exists")
	want(t, c.Do("POST", "/api/v1/recipes/999/duplicate", nil, key()...), 404)
}

// Persona ficticia (nunca datos de Axel o Ximena): hombre de 80 kg, 1.75 m, nacido el 20-may-1990, actividad moderada.
// MSJ = 10 × 80 + 6.25 × 175 − 5 × 36 + 5 = 1718.75; × 1.55 = 2664.06 → 2664. IMC = 80 / 1.75² = 26.12 → 261.
func profile(t *testing.T, s *testutil.Server, c *testutil.Client) api.HealthProfile {
	t.Helper()
	var h api.HealthProfile
	body := map[string]any{"person_id": s.Fx.AxelPersonID, "height_mm": 1750, "birth_date": "1990-05-20", "sex": "male", "activity_level": "moderate", "goal": "maintain"}
	want(t, c.Do("POST", "/api/v1/health-profiles", body, key()...), 201).JSON(t, &h)
	want(t, c.Do("POST", "/api/v1/weight-entries", map[string]any{"person_id": s.Fx.AxelPersonID, "measured_date": "2026-10-11", "weight_g": 80000}, key()...), 201)
	return h
}

func TestOp_createHealthProfile(t *testing.T) {
	s, c := setup(t)
	h := profile(t, s, c)
	if *h.AgeYears != 36 || h.SuggestedKcalTarget != nil || *h.HouseholdGoalAdjustLoseKcal != 0 {
		t.Fatalf("%+v", h) // sin peso todavía: la meta sugerida queda ausente
	}
	wantCode(t, c.Do("POST", "/api/v1/health-profiles", map[string]any{"person_id": s.Fx.AxelPersonID}, key()...), 409, "already_exists")
	wantCode(t, c.Do("POST", "/api/v1/health-profiles", map[string]any{"person_id": s.Fx.KarenPersonID}, key()...), 422, "external_person")
	wantCode(t, c.Do("POST", "/api/v1/health-profiles", map[string]any{"person_id": s.Fx.XimenaPersonID, "height_mm": 100}, key()...), 422, "validation_failed")
}

func TestOp_listHealthProfiles(t *testing.T) {
	s, c := setup(t)
	profile(t, s, c)
	var pg api.HealthProfilePage
	want(t, c.Do("GET", "/api/v1/health-profiles", nil), 200).JSON(t, &pg)
	h := pg.Items[0]
	if len(pg.Items) != 1 || *h.SuggestedKcalTarget != 2664 || *h.BMITenths != 261 || *h.BMIClass != "overweight" || *h.CurrentWeightG != 80000 {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/health-profiles?sort=peso", nil), 400)
}

func TestOp_updateHealthProfile(t *testing.T) {
	s, c := setup(t)
	h := profile(t, s, c)
	path := fmt.Sprintf("/api/v1/health-profiles/%d", h.ID)
	want(t, c.Do("PATCH", path, map[string]any{"goal": "lose", "goal_adjust_lose_kcal": 500, "notes": nil}), 200).JSON(t, &h)
	if *h.SuggestedKcalTarget != 2164 || *h.GoalAdjustLoseKcal != 500 {
		t.Fatalf("%+v", h)
	}
	h = api.HealthProfile{}
	want(t, c.Do("PATCH", path, map[string]any{"goal_adjust_lose_kcal": nil, "daily_kcal_target": 2200}), 200).JSON(t, &h)
	if *h.SuggestedKcalTarget != 2664 || *h.DailyKcalTarget != 2200 || h.GoalAdjustLoseKcal != nil {
		t.Fatalf("%+v", h) // sin ajuste propio usa el de la casa (0)
	}
	wantCode(t, c.Do("PATCH", path, map[string]any{"daily_kcal_target": 100}), 422, "validation_failed")
	want(t, c.Do("PATCH", path, map[string]any{"sex": "otro"}), 400)
	want(t, c.Do("PATCH", "/api/v1/health-profiles/999", map[string]any{"notes": "x"}), 404)
}

func TestOp_getTodayKcal(t *testing.T) {
	s, c := setup(t)
	profile(t, s, c)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	plan(t, c, "2026-10-12", []int64{s.Fx.XimenaPersonID}, map[string]any{"is_eating_out": true, "recipe_id": nil})
	var k api.TodayKcal
	want(t, c.Do("GET", "/api/v1/health/today-kcal", nil), 200).JSON(t, &k)
	if k.Date != "2026-10-12" || len(k.Persons) != 2 || k.Persons[0].PlannedKcal != 530 || *k.Persons[0].DifferenceKcal != 530-2664 || k.Persons[1].EatingOutSlots != 1 {
		t.Fatalf("%+v", k)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/health/today-kcal", nil), 401)
}

func TestOp_listWeightEntries(t *testing.T) {
	s, c := setup(t)
	profile(t, s, c)
	var pg api.WeightEntryPage
	want(t, c.Do("GET", fmt.Sprintf("/api/v1/weight-entries?person_id=%d&period=2026-10&limit=1", s.Fx.AxelPersonID), nil), 200).JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].BMITenths != 261 {
		t.Fatalf("%+v", pg)
	}
	want(t, c.Do("GET", "/api/v1/weight-entries?limit=201", nil), 400)
}

func TestOp_createWeightEntry(t *testing.T) {
	s, c := setup(t)
	k := key()
	body := map[string]any{"person_id": s.Fx.XimenaPersonID, "measured_date": "2026-10-10", "weight_g": 60500, "captured_at": "2026-10-10T07:30:00-06:00"}
	var e1, e2 api.WeightEntry
	want(t, c.Do("POST", "/api/v1/weight-entries", body, k...), 201).JSON(t, &e1)
	if e1.CapturedAt == nil || e1.BMITenths != nil {
		t.Fatalf("%+v", e1) // sin estatura no hay IMC
	}
	res := want(t, c.Do("POST", "/api/v1/weight-entries", body, k...), 200)
	if res.Header.Get("Idempotent-Replayed") != "true" {
		t.Fatal("el reintento debe marcarse como repetido")
	}
	// Pasadas las 24 h (la tabla idempotency_key ya se purgó) la clave guardada en weight_entry evita el duplicado (K9).
	if _, err := s.DB.Exec(`DELETE FROM idempotency_key`); err != nil {
		t.Fatal(err)
	}
	res = want(t, c.Do("POST", "/api/v1/weight-entries", body, k...), 200)
	res.JSON(t, &e2)
	if e2.ID != e1.ID || res.Header.Get("Idempotent-Replayed") != "true" {
		t.Fatalf("%+v", e2)
	}
	wantCode(t, c.Do("POST", "/api/v1/weight-entries", body, key()...), 409, "already_exists")
	wantCode(t, c.Do("POST", "/api/v1/weight-entries", map[string]any{"person_id": s.Fx.KarenPersonID, "measured_date": "2026-10-10", "weight_g": 60000}, key()...), 422, "external_person")
	wantCode(t, c.Do("POST", "/api/v1/weight-entries", map[string]any{"person_id": s.Fx.AxelPersonID, "measured_date": "2026-10-10", "weight_g": 10}, key()...), 422, "validation_failed")
	want(t, c.Do("POST", "/api/v1/weight-entries", map[string]any{"person_id": s.Fx.AxelPersonID, "measured_date": "2026-02-30", "weight_g": 70000}, key()...), 400)
}

func TestOp_updateWeightEntry(t *testing.T) {
	s, c := setup(t)
	var e api.WeightEntry
	want(t, c.Do("POST", "/api/v1/weight-entries", map[string]any{"person_id": s.Fx.AxelPersonID, "measured_date": "2026-10-11", "weight_g": 80000}, key()...), 201).JSON(t, &e)
	path := fmt.Sprintf("/api/v1/weight-entries/%d", e.ID)
	want(t, c.Do("PATCH", path, map[string]any{"weight_g": 79500, "measured_date": "2026-10-12", "notes": "en ayunas"}), 200).JSON(t, &e)
	if e.WeightG != 79500 || e.MeasuredDate != "2026-10-12" || *e.Notes != "en ayunas" {
		t.Fatalf("%+v", e)
	}
	wantCode(t, c.Do("PATCH", path, map[string]any{"weight_g": 500000}), 422, "validation_failed")
	want(t, c.Do("PATCH", "/api/v1/weight-entries/999", map[string]any{"weight_g": 70000}), 404)
}

func TestOp_deleteWeightEntry(t *testing.T) {
	s, c := setup(t)
	var e api.WeightEntry
	want(t, c.Do("POST", "/api/v1/weight-entries", map[string]any{"person_id": s.Fx.AxelPersonID, "measured_date": "2026-10-11", "weight_g": 80000}, key()...), 201).JSON(t, &e)
	want(t, c.Do("DELETE", fmt.Sprintf("/api/v1/weight-entries/%d", e.ID), nil), 204)
	want(t, c.Do("DELETE", fmt.Sprintf("/api/v1/weight-entries/%d", e.ID), nil), 404)
}

func TestOp_getHealthSettings(t *testing.T) {
	s, c := setup(t)
	var hs api.HealthSettings
	want(t, c.Do("GET", "/api/v1/settings/health", nil), 200).JSON(t, &hs)
	p := hs.Parameters.HealthGoalAdjustLoseKcal
	if p.Key != "health.goal_adjust_lose_kcal" || !p.Value.Valid || p.Value.Value != 0 || !p.IsInitial {
		t.Fatalf("%+v", hs)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/settings/health", nil), 401)
}

func TestOp_updateHealthSettings(t *testing.T) {
	s, c := setup(t)
	h := profile(t, s, c)
	var hs api.HealthSettings
	want(t, c.Do("PATCH", "/api/v1/settings/health", map[string]any{"parameters": map[string]any{"health.goal_adjust_lose_kcal": 300}}), 200).JSON(t, &hs)
	if hs.Parameters.HealthGoalAdjustLoseKcal.Value.Value != 300 || hs.Parameters.HealthGoalAdjustLoseKcal.IsInitial {
		t.Fatalf("%+v", hs)
	}
	want(t, c.Do("PATCH", fmt.Sprintf("/api/v1/health-profiles/%d", h.ID), map[string]any{"goal": "lose"}), 200).JSON(t, &h)
	if *h.SuggestedKcalTarget != 2364 || *h.HouseholdGoalAdjustLoseKcal != 300 {
		t.Fatalf("%+v", h) // 2664.06 − 300 = 2364.06 → 2364 con el ajuste de la casa
	}
	both := map[string]any{"parameters": map[string]any{"health.goal_adjust_lose_kcal": 1}, "reset_parameters": []string{"health.goal_adjust_lose_kcal"}}
	wantCode(t, c.Do("PATCH", "/api/v1/settings/health", both), 422, "validation_failed")
	wantCode(t, c.Do("PATCH", "/api/v1/settings/health", map[string]any{"parameters": map[string]any{"health.goal_adjust_gain_kcal": -1}}), 422, "validation_failed")
	want(t, c.Do("PATCH", "/api/v1/settings/health", map[string]any{"reset_parameters": []string{"health.goal_adjust_lose_kcal"}}), 200).JSON(t, &hs)
	if !hs.Parameters.HealthGoalAdjustLoseKcal.IsInitial {
		t.Fatalf("%+v", hs)
	}
}
