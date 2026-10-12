package comidas

import (
	"context"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

// WeekMeals arma la tarjeta «Comidas de la semana» de Inicio (getDashboardWeekMeals, N-023, métrica M-05). La
// operación la registra api-gastos-inicio; el armado del motor S4 es de esta área y no se duplica (bloqueo
// S6-gastos-inicio §15). monday es el lunes de la semana (K8).
//
//   - cost_cents y unplanned_slots: los mismos que week_cost_cents y unplanned_slots de getMealPlanWeek (fillWeek).
//   - days_with_plan: días con al menos una casilla de la persona, con receta o «fuera».
//   - avg_kcal_per_day: Σ kcal de la semana ÷ days_with_plan, ROUND_HALF_UP; ausente si no hay días con plan.
//   - daily_kcal_target: la meta que rige el último día de la semana que la tiene (igual que getMealPlanTrends).
//
// PROVISIONAL (bloqueo S6-comidas-calculos, punto 14): si un día con solo «fuera» cuenta como día con plan.
func WeekMeals(ctx context.Context, reg *api.Registry, monday dominio.Date) (api.DashboardWeekMeals, error) {
	a := &area{reg: reg}
	m, err := a.loadMotor(ctx, a.db(), loadOpts{food: true, from: monday, to: monday.AddDays(6)})
	if err != nil {
		return api.DashboardWeekMeals{}, err
	}
	wp, err := m.e.WeekPlanFor(rdate(monday.String()))
	if err != nil {
		return api.DashboardWeekMeals{}, err
	}
	f := fillWeek(wp)
	out := api.DashboardWeekMeals{WeekStart: monday.String(), WeekEnd: monday.AddDays(6).String(), CostCents: f.week,
		Persons: []api.DashboardWeekMealsPersonsItem{}, UnplannedSlots: int64(wp.UnplannedSlots)}
	days := map[string]map[string]bool{}
	for _, c := range wp.Cells {
		if days[c.Person] == nil {
			days[c.Person] = map[string]bool{}
		}
		days[c.Person][c.Date.String()] = true
	}
	for _, p := range m.persons {
		if !p.HasMealPlan() {
			continue
		}
		it := api.DashboardWeekMealsPersonsItem{PersonID: p.ID, DaysWithPlan: int64(len(days[p.Code]))}
		var sum int64
		for _, d := range wp.Days {
			sum += f.kcal[d.Date.String()+"|"+p.Code]
			for _, x := range d.Persons {
				if x.Person == p.Code && x.DailyKcalTarget != nil {
					it.DailyKcalTarget = x.DailyKcalTarget
				}
			}
		}
		if it.DaysWithPlan > 0 {
			avg := dominio.RoundHalfUp(sum, it.DaysWithPlan)
			it.AvgKcalPerDay = &avg
		}
		out.Persons = append(out.Persons, it)
	}
	return out, nil
}
