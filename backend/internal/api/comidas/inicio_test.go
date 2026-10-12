package comidas_test

import (
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comidas"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/testutil"
)

// La operación es de api-gastos-inicio (tarjeta 4 de Inicio, N-023) y el armado es comidas.WeekMeals; la prueba vive
// aquí para usar el catálogo ficticio de esta área (receta 1: 500 centavos y 530 kcal por porción).
func TestOp_getDashboardWeekMeals(t *testing.T) {
	s := testutil.NewServer(t, comidas.Register, gastos.Register)
	if _, err := s.DB.Exec(fixture); err != nil {
		t.Fatal(err)
	}
	c := s.SessionFor(t, testutil.AxelUser)
	plan(t, c, "2026-10-12", []int64{s.Fx.AxelPersonID}, nil)
	plan(t, c, "2026-10-13", []int64{s.Fx.AxelPersonID}, nil)
	plan(t, c, "2026-10-13", []int64{s.Fx.XimenaPersonID}, map[string]any{"is_eating_out": true, "recipe_id": nil})
	var wk api.DashboardWeekMeals
	want(t, c.Do("GET", "/api/v1/dashboard/week-meals", nil), 200).JSON(t, &wk) // 12-oct-2026 (lunes): semana W42
	var mp api.MealPlanWeek
	want(t, c.Do("GET", "/api/v1/meal-plan/week?period=2026-W42", nil), 200).JSON(t, &mp)
	// Costo: 2 casillas con receta × 500; «fuera» no suma. Sin armar: 21 pares menos el martes en la comida (las dos personas).
	if wk.WeekStart != "2026-10-12" || wk.WeekEnd != "2026-10-18" || wk.CostCents != 1000 || wk.CostCents != mp.WeekCostCents ||
		wk.UnplannedSlots != 20 || wk.UnplannedSlots != mp.UnplannedSlots || len(wk.Persons) != 2 {
		t.Fatalf("%+v (plan de la semana: costo %d, sin armar %d)", wk, mp.WeekCostCents, mp.UnplannedSlots)
	}
	axel, ximena := wk.Persons[0], wk.Persons[1]
	if axel.PersonID != s.Fx.AxelPersonID || axel.DaysWithPlan != 2 || axel.AvgKcalPerDay == nil || *axel.AvgKcalPerDay != 530 {
		t.Fatalf("Axel: %+v", axel)
	}
	// Ximena solo tiene «fuera»: cuenta como día con plan y 0 kcal (PROVISIONAL, S6-comidas-calculos punto 14).
	if ximena.PersonID != s.Fx.XimenaPersonID || ximena.DaysWithPlan != 1 || ximena.AvgKcalPerDay == nil || *ximena.AvgKcalPerDay != 0 {
		t.Fatalf("Ximena: %+v", ximena)
	}
	// Semana sin plan: sin promedio y 21 casillas sin armar.
	wk = api.DashboardWeekMeals{} // JSON reusa los elementos del slice: sin limpiar, un campo ausente conserva el valor anterior
	want(t, c.Do("GET", "/api/v1/dashboard/week-meals?period=2026-W41", nil), 200).JSON(t, &wk)
	if wk.WeekStart != "2026-10-05" || wk.CostCents != 0 || wk.UnplannedSlots != 21 || wk.Persons[0].DaysWithPlan != 0 || wk.Persons[0].AvgKcalPerDay != nil {
		t.Fatalf("semana vacía: %+v", wk)
	}
	// Domingo 18-oct-2026: la tarjeta sigue en la semana en curso (K8), no en la siguiente (K38 es solo de P-10).
	s.Clock.Set(testutil.MX(2026, time.October, 18, 9, 0))
	wk = api.DashboardWeekMeals{}
	want(t, c.Do("GET", "/api/v1/dashboard/week-meals", nil), 200).JSON(t, &wk)
	if wk.WeekStart != "2026-10-12" || wk.CostCents != 1000 {
		t.Fatalf("domingo: %+v", wk)
	}
	wantCode(t, c.Do("GET", "/api/v1/dashboard/week-meals?period=2026-10", nil), 422, "validation_failed")
	want(t, c.Do("GET", "/api/v1/dashboard/week-meals?period=semana", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/dashboard/week-meals", nil), 401)
}
