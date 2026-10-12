package comidas

import (
	"database/sql"
	"fmt"
	"net/http"
	"slices"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/recetas"
)

// ---- Tiempos de comida ----

func (a *area) listMealSlots(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"sort_order", "name", "code"}, "sort_order")
	if err != nil {
		return err
	}
	rows, err := dominio.ListMealSlots(r.Context(), a.db(), p.OrderBy(map[string]string{"sort_order": "ms.sort_order", "name": "ms.name", "code": "ms.code", "id": "ms.id"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.MealSlot, 0, len(rows))
	for _, s := range rows {
		items = append(items, apiSlot(s))
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

func apiSlot(s dominio.MealSlotRow) api.MealSlot {
	return api.MealSlot{ID: s.ID, Code: s.Code, Name: s.Name, SortOrder: s.SortOrder, IsOptional: s.IsOptional,
		CreatedAt: api.APITime(s.CreatedAt), UpdatedAt: api.APITime(s.UpdatedAt)}
}

func apiEntry(e dominio.MealPlanEntryRow) api.MealPlanEntry {
	return api.MealPlanEntry{ID: e.ID, PlanDate: e.PlanDate, MealSlotID: e.MealSlotID, PersonID: e.PersonID, RecipeID: nint(e.RecipeID),
		ServingsMilli: nint(e.ServingsMilli), IsEatingOut: e.IsEatingOut, Notes: nstr(e.Notes), CreatedAt: api.APITime(e.CreatedAt), UpdatedAt: api.APITime(e.UpdatedAt)}
}

// ---- Semana, mes y tendencias ----

// weekParam: `period` de semana (K8) o, sin él, la semana por defecto de S4 (K38: los domingos, la siguiente).
// Con badRequest = true (operaciones sin 422 en C1) un periodo que no es semana responde 400.
func (a *area) weekParam(r *http.Request, badRequest bool) (dominio.Date, error) {
	p, err := api.QueryPeriod(r)
	if err != nil {
		return dominio.Date{}, err
	}
	if p == nil {
		return ddate(recetas.DefaultWeek(rdate(a.today().String()))), nil
	}
	if p.Kind != dominio.PeriodWeek && badRequest {
		return dominio.Date{}, api.BadRequest("El periodo debe ser una semana, por ejemplo 2026-W42.")
	}
	if p.Kind != dominio.PeriodWeek {
		return dominio.Date{}, dominio.Invalid(issue("period", "invalid_format", "Elige una semana, por ejemplo 2026-W42."))
	}
	return p.Start, nil
}

// weekFill: costo por día y kcal por persona y día del motor; cuando el motor los deja ausentes (falta un
// precio o un dato nutricional) y C1 los pide como entero, se usa la suma de las casillas que sí los tienen
// (PROVISIONAL, bloqueo S6-comidas-calculos).
type weekFill struct {
	dayCost map[string]int64
	kcal    map[string]int64
	week    int64
}

func fillWeek(wp *recetas.WeekPlan) weekFill {
	knownCost, knownKcal := map[string]int64{}, map[string]int64{}
	for _, c := range wp.Cells {
		if c.CostCents != nil {
			knownCost[c.Date.String()] += *c.CostCents
		}
		if c.Kcal != nil {
			knownKcal[c.Date.String()+"|"+c.Person] += *c.Kcal
		}
	}
	f := weekFill{dayCost: map[string]int64{}, kcal: map[string]int64{}}
	for _, d := range wp.Days {
		ds := d.Date.String()
		f.dayCost[ds] = knownCost[ds]
		if d.CostCents != nil {
			f.dayCost[ds] = *d.CostCents
		}
		f.week += f.dayCost[ds]
		for _, p := range d.Persons {
			k := ds + "|" + p.Person
			f.kcal[k] = knownKcal[k]
			if p.Kcal != nil {
				f.kcal[k] = *p.Kcal
			}
		}
	}
	if wp.WeekCostCents != nil {
		f.week = *wp.WeekCostCents
	}
	return f
}

func (m *motor) apiWeek(monday dominio.Date, wp *recetas.WeekPlan, only []int64) api.MealPlanWeek {
	keep := func(id int64) bool { return len(only) == 0 || slices.Contains(only, id) }
	entryID := map[string]int64{}
	for _, en := range m.entries {
		entryID[fmt.Sprint(en.PlanDate, en.MealSlotID, en.PersonID, en.RecipeID.Int64)] = en.ID
	}
	f := fillWeek(wp)
	out := api.MealPlanWeek{WeekStart: monday.String(), WeekEnd: monday.AddDays(6).String(), Cells: []api.MealPlanWeekCellsItem{},
		Days: []api.MealPlanWeekDaysItem{}, WeekCostCents: f.week, UnplannedSlots: int64(wp.UnplannedSlots), Leftovers: []api.MealPlanWeekLeftoversItem{}}
	for _, c := range wp.Cells {
		pid := m.personID[c.Person]
		if !keep(pid) {
			continue
		}
		var rid int64
		if c.Recipe != nil {
			rid = idOf("r", *c.Recipe)
		}
		it := api.MealPlanWeekCellsItem{PlanDate: c.Date.String(), MealSlotID: m.slotID[c.Slot], PersonID: pid, ServingsMilli: c.ServingsMilli,
			IsEatingOut: c.EatingOut, CostCents: c.CostCents, Kcal: c.Kcal}
		if id, ok := entryID[fmt.Sprint(c.Date.String(), it.MealSlotID, pid, rid)]; ok {
			it.MealPlanEntryID = &id
		}
		if rid != 0 {
			it.RecipeID = &rid
		}
		if c.PriceStore != nil {
			s := idOf("s", *c.PriceStore)
			it.PriceStoreID = &s
		}
		if c.PriceObservedDate != nil {
			s := c.PriceObservedDate.String()
			it.PriceObservedDate = &s
		}
		out.Cells = append(out.Cells, it)
	}
	for _, d := range wp.Days {
		day := api.MealPlanWeekDaysItem{PlanDate: d.Date.String(), CostCents: f.dayCost[d.Date.String()], Persons: []api.MealPlanWeekDaysItemPersonsItem{}}
		for _, p := range d.Persons {
			if pid := m.personID[p.Person]; keep(pid) {
				day.Persons = append(day.Persons, api.MealPlanWeekDaysItemPersonsItem{PersonID: pid, Kcal: f.kcal[d.Date.String()+"|"+p.Person],
					DailyKcalTarget: p.DailyKcalTarget, DifferenceKcal: p.DifferenceKcal})
			}
		}
		out.Days = append(out.Days, day)
	}
	for _, l := range wp.Leftovers {
		out.Leftovers = append(out.Leftovers, api.MealPlanWeekLeftoversItem{RecipeID: idOf("r", l.Recipe), ServingsMilli: l.ServingsMilli})
	}
	return out
}

func (a *area) getMealPlanWeek(w http.ResponseWriter, r *http.Request) error {
	monday, err := a.weekParam(r, false)
	if err != nil {
		return err
	}
	only, err := api.QueryIDs(r, "person_id")
	if err != nil {
		return err
	}
	m, err := a.loadMotor(r.Context(), a.db(), loadOpts{food: true, from: monday, to: monday.AddDays(6)})
	if err != nil {
		return err
	}
	wp, err := m.e.WeekPlanFor(rdate(monday.String()))
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, m.apiWeek(monday, wp, only))
}

func (a *area) getMealPlanMonth(w http.ResponseWriter, r *http.Request) error {
	p, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	if p == nil {
		mp := dominio.MonthPeriod(a.today())
		p = &mp
	}
	if p.Kind != dominio.PeriodMonth {
		return dominio.Invalid(issue("period", "invalid_format", "Elige un mes, por ejemplo 2026-10."))
	}
	ctx := r.Context()
	es, err := dominio.ListMealPlanEntries(ctx, a.db(), p.Start.String(), p.End.String())
	if err != nil {
		return err
	}
	m, err := a.loadMotor(ctx, a.db(), loadOpts{})
	if err != nil {
		return err
	}
	has, any := map[string]bool{}, map[string]bool{}
	for _, e := range es {
		has[fmt.Sprint(e.PlanDate, e.MealSlotID, e.PersonID)], any[e.PlanDate] = true, true
	}
	out := api.MealPlanMonth{Month: p.Key, Days: []api.MealPlanMonthDaysItem{}}
	for d := p.Start; !d.After(p.End); d = d.AddDays(1) {
		status := "empty"
		if any[d.String()] {
			status = "complete"
			for _, s := range m.slots {
				for _, per := range m.persons {
					if !s.IsOptional && per.HasMealPlan() && !has[fmt.Sprint(d.String(), s.ID, per.ID)] {
						status = "incomplete"
					}
				}
			}
		}
		out.Days = append(out.Days, api.MealPlanMonthDaysItem{PlanDate: d.String(), Status: status})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (a *area) getMealPlanTrends(w http.ResponseWriter, r *http.Request) error {
	only, err := api.QueryIDs(r, "person_id")
	if err != nil {
		return err
	}
	ctx := r.Context()
	last := ddate(recetas.DefaultWeek(rdate(a.today().String())))
	dates, err := dominio.MealPlanDates(ctx, a.db(), last.AddDays(6).String())
	if err != nil {
		return err
	}
	var mondays []dominio.Date
	for _, s := range dates {
		d, err := dominio.ParseDate(s)
		if err != nil {
			return err
		}
		if mo := mondayOf(d); len(mondays) == 0 || mondays[len(mondays)-1] != mo {
			mondays = append(mondays, mo)
		}
	}
	if len(mondays) > 12 {
		mondays = mondays[len(mondays)-12:]
	}
	out := api.MealPlanTrends{Weeks: []api.MealPlanTrendsWeeksItem{}}
	if len(mondays) == 0 {
		return api.WriteJSON(w, http.StatusOK, out)
	}
	m, err := a.loadMotor(ctx, a.db(), loadOpts{food: true, from: mondays[0], to: mondays[len(mondays)-1].AddDays(6)})
	if err != nil {
		return err
	}
	for _, mo := range mondays {
		wp, err := m.e.WeekPlanFor(rdate(mo.String()))
		if err != nil {
			return err
		}
		f := fillWeek(wp)
		wk := api.MealPlanTrendsWeeksItem{WeekStart: mo.String(), CostCents: f.week, Persons: []api.MealPlanTrendsWeeksItemPersonsItem{}}
		for _, p := range m.persons {
			if !p.HasMealPlan() || (len(only) > 0 && !slices.Contains(only, p.ID)) {
				continue
			}
			var sum int64
			it := api.MealPlanTrendsWeeksItemPersonsItem{PersonID: p.ID}
			for _, d := range wp.Days {
				sum += f.kcal[d.Date.String()+"|"+p.Code]
				for _, x := range d.Persons {
					if x.Person == p.Code && x.DailyKcalTarget != nil {
						it.DailyKcalTarget = x.DailyKcalTarget
					}
				}
			}
			avg := dominio.RoundHalfUp(sum, 7) // promedio por día de la semana de 7 días (PROVISIONAL, bloqueo S6-comidas-calculos)
			it.AvgKcalPerDay = &avg
			wk.Persons = append(wk.Persons, it)
		}
		out.Weeks = append(out.Weeks, wk)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// ---- Casillas del plan ----

func (a *area) checkPlanPerson(m *motor, v *api.Validator, id int64, field string) error {
	p, ok := m.person(id)
	if !ok {
		v.Add(field, "not_found", "Esa persona no existe.")
		return nil
	}
	if !p.HasMealPlan() {
		return externalPerson(field)
	}
	return nil
}

func (a *area) recipeExists(r *http.Request, q dominio.Querier, v *api.Validator, id int64, field string) error {
	if _, err := dominio.GetRecipeRow(r.Context(), q, id); isNotFound(err) {
		v.Add(field, "not_found", "Esa receta no existe.")
	} else if err != nil {
		return err
	}
	return nil
}

// cellConflict: misma casilla no mezcla «fuera» con recetas (409 slot_conflict) ni repite receta (409 already_exists).
func cellConflict(cell []dominio.MealPlanEntryRow, self int64, eatingOut bool, recipeID int64) error {
	for _, c := range cell {
		if c.ID == self {
			continue
		}
		if eatingOut || c.IsEatingOut {
			return dominio.Conflict("slot_conflict", "Esa casilla ya tiene «fuera» o una receta y no se pueden mezclar. Vacíala primero.")
		}
		if c.RecipeID.Int64 == recipeID {
			return dominio.Conflict("already_exists", "Esa receta ya está en esa casilla.")
		}
	}
	return nil
}

func (a *area) createMealPlanEntries(w http.ResponseWriter, r *http.Request) error {
	var in api.MealPlanEntryCreate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	var v api.Validator
	var dates []dominio.Date
	for i, s := range in.PlanDates {
		d, err := v.Date(s, fmt.Sprintf("plan_dates[%d]", i))
		if err != nil {
			return err
		}
		dates = append(dates, d)
	}
	m, err := a.loadMotor(ctx, db, loadOpts{})
	if err != nil {
		return err
	}
	v.Check(len(dates) > 0, "plan_dates", "required", "Elige al menos un día.")
	v.Check(len(in.PersonIDs) > 0, "person_ids", "required", "Elige a Axel, a Ximena o a los dos.")
	_, okSlot := m.slotCode[in.MealSlotID]
	v.Check(okSlot, "meal_slot_id", "not_found", "Ese tiempo de comida no existe.")
	for _, id := range in.PersonIDs {
		if err := a.checkPlanPerson(m, &v, id, "person_ids"); err != nil {
			return err
		}
	}
	servings := int64(1000)
	if in.ServingsMilli != nil {
		servings = *in.ServingsMilli
	}
	if in.IsEatingOut {
		v.Check(in.RecipeID == nil && in.ServingsMilli == nil, "recipe_id", "not_allowed", "«Fuera» no lleva receta ni porciones.")
	} else {
		v.Check(in.RecipeID != nil, "recipe_id", "required", "Elige una receta o marca «fuera».")
		v.Check(servings > 0, "servings_milli", "out_of_range", "Las porciones deben ser mayores que cero.")
		if in.RecipeID != nil {
			if err := a.recipeExists(r, db, &v, *in.RecipeID, "recipe_id"); err != nil {
				return err
			}
		}
	}
	if err := v.Err(); err != nil {
		return err
	}
	out := api.MealPlanEntryList{Entries: []api.MealPlanEntry{}}
	err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		for _, d := range dates {
			for _, pid := range in.PersonIDs {
				cell, err := dominio.MealPlanCell(ctx, tx, d.String(), in.MealSlotID, pid)
				if err != nil {
					return err
				}
				e := dominio.MealPlanEntryRow{PlanDate: d.String(), MealSlotID: in.MealSlotID, PersonID: pid, IsEatingOut: in.IsEatingOut, Notes: toNS(in.Notes)}
				if !in.IsEatingOut {
					e.RecipeID, e.ServingsMilli = toNI(in.RecipeID), sql.NullInt64{Int64: servings, Valid: true}
				}
				if err := cellConflict(cell, 0, e.IsEatingOut, e.RecipeID.Int64); err != nil {
					return err
				}
				if err := dominio.SaveMealPlanEntry(ctx, tx, &e); err != nil {
					return err
				}
				if err := a.audit(ctx, tx, r, "insert", "meal_plan_entry", e.ID, nil); err != nil {
					return err
				}
				row, err := dominio.GetMealPlanEntry(ctx, tx, e.ID)
				if err != nil {
					return err
				}
				out.Entries = append(out.Entries, apiEntry(row))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func (a *area) updateMealPlanEntry(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "meal_plan_entry_id")
	if err != nil {
		return err
	}
	var in api.MealPlanEntryUpdate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if !in.RecipeID.Set && !in.ServingsMilli.Set && in.IsEatingOut == nil && !in.Notes.Set {
		return dominio.Invalid(issue("body", "required", "Manda al menos un cambio."))
	}
	ctx, db := r.Context(), a.db()
	var row dominio.MealPlanEntryRow
	err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		e, err := dominio.GetMealPlanEntry(ctx, tx, id)
		if err != nil {
			return err
		}
		before, err := dominio.RowJSON(ctx, tx, "meal_plan_entry", id)
		if err != nil {
			return err
		}
		e.RecipeID, e.ServingsMilli, e.Notes = mergeI(in.RecipeID, e.RecipeID), mergeI(in.ServingsMilli, e.ServingsMilli), mergeS(in.Notes, e.Notes)
		if in.IsEatingOut != nil {
			e.IsEatingOut = *in.IsEatingOut
		}
		var v api.Validator
		if e.IsEatingOut {
			v.Check(!e.RecipeID.Valid && !e.ServingsMilli.Valid, "recipe_id", "not_allowed", "«Fuera» no lleva receta ni porciones: manda recipe_id y servings_milli en null.")
		} else {
			v.Check(e.RecipeID.Valid, "recipe_id", "required", "Elige una receta o marca «fuera».")
			v.Check(e.ServingsMilli.Valid && e.ServingsMilli.Int64 > 0, "servings_milli", "out_of_range", "Las porciones deben ser mayores que cero.")
			if in.RecipeID.Valid {
				if err := a.recipeExists(r, tx, &v, in.RecipeID.Value, "recipe_id"); err != nil {
					return err
				}
			}
		}
		if err := v.Err(); err != nil {
			return err
		}
		cell, err := dominio.MealPlanCell(ctx, tx, e.PlanDate, e.MealSlotID, e.PersonID)
		if err != nil {
			return err
		}
		if err := cellConflict(cell, e.ID, e.IsEatingOut, e.RecipeID.Int64); err != nil {
			return err
		}
		if err := dominio.SaveMealPlanEntry(ctx, tx, &e); err != nil {
			return err
		}
		if err := a.audit(ctx, tx, r, "update", "meal_plan_entry", id, before); err != nil {
			return err
		}
		row, err = dominio.GetMealPlanEntry(ctx, tx, id)
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, apiEntry(row))
}

func (a *area) deleteMealPlanEntry(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "meal_plan_entry_id")
	if err != nil {
		return err
	}
	ctx := r.Context()
	err = dominio.WithTx(ctx, a.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetMealPlanEntry(ctx, tx, id); err != nil {
			return err
		}
		before, err := dominio.RowJSON(ctx, tx, "meal_plan_entry", id)
		if err != nil {
			return err
		}
		if err := dominio.DeleteMealPlanEntry(ctx, tx, id); err != nil {
			return err
		}
		return a.audit(ctx, tx, r, "delete", "meal_plan_entry", id, before)
	})
	if err != nil {
		return err
	}
	return api.NoContent(w)
}

// copyPreviousWeek copia a la semana `period` (sin él, la semana por defecto, PROVISIONAL) las entradas de la
// semana previa, solo en casillas (día × tiempo × persona) vacías.
func (a *area) copyPreviousWeek(w http.ResponseWriter, r *http.Request) error {
	monday, err := a.weekParam(r, true)
	if err != nil {
		return err
	}
	ctx := r.Context()
	n := int64(0)
	err = dominio.WithTx(ctx, a.db(), func(tx *sql.Tx) error {
		prev, err := dominio.ListMealPlanEntries(ctx, tx, monday.AddDays(-7).String(), monday.AddDays(-1).String())
		if err != nil {
			return err
		}
		cur, err := dominio.ListMealPlanEntries(ctx, tx, monday.String(), monday.AddDays(6).String())
		if err != nil {
			return err
		}
		busy := map[string]bool{}
		for _, e := range cur {
			busy[fmt.Sprint(e.PlanDate, e.MealSlotID, e.PersonID)] = true
		}
		for _, e := range prev {
			d, err := dominio.ParseDate(e.PlanDate)
			if err != nil {
				return err
			}
			e.ID, e.PlanDate = 0, d.AddDays(7).String()
			if busy[fmt.Sprint(e.PlanDate, e.MealSlotID, e.PersonID)] {
				continue
			}
			if err := dominio.SaveMealPlanEntry(ctx, tx, &e); err != nil {
				return err
			}
			if err := a.audit(ctx, tx, r, "insert", "meal_plan_entry", e.ID, nil); err != nil {
				return err
			}
			n++
		}
		return nil
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, api.CopyWeekResult{CreatedCount: n})
}
