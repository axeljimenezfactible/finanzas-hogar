package recetas

import "sort"

// Plan semanal por persona (§5; RN-14, RN-19, K8, K38, K43; S4-05).

// weekEntries valida y ordena las entradas de la semana (fecha, tiempo, persona, orden de captura).
func (e *Engine) weekEntries(monday Date) ([]PlanEntry, error) {
	pidx := map[string]int{}
	for i, p := range e.Persons {
		pidx[p.Code] = i
	}
	var out []PlanEntry
	cells := map[string][2]bool{} // [hay receta, hay «fuera»]
	for _, en := range e.Plan {
		if !en.Date.inRange(monday, monday.AddDays(6)) {
			continue
		}
		p, err := e.person(en.Person)
		if err != nil {
			return nil, err
		}
		if !p.HasMealPlan {
			return nil, errf(ErrPersonWithoutPlan, "%q no tiene plan de comidas", en.Person)
		}
		if _, ok := slotOrder[en.Slot]; !ok {
			return nil, errf(ErrInvalidArgument, "tiempo %q", en.Slot)
		}
		k := en.Date.String() + "|" + string(en.Slot) + "|" + en.Person
		c := cells[k]
		if en.EatingOut {
			if en.Recipe != "" || en.ServingsMilli != nil || c[0] {
				return nil, errf(ErrEatingOutConflict, "%s", k)
			}
			c[1] = true
		} else {
			if c[1] {
				return nil, errf(ErrEatingOutConflict, "%s", k)
			}
			if _, err := e.recipe(en.Recipe); err != nil {
				return nil, err
			}
			if en.ServingsMilli == nil || *en.ServingsMilli <= 0 {
				return nil, errf(ErrInvalidServings, "%s", k)
			}
			c[0] = true
		}
		cells[k] = c
		out = append(out, en)
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Date != b.Date {
			return a.Date.Before(b.Date)
		}
		if a.Slot != b.Slot {
			return slotOrder[a.Slot] < slotOrder[b.Slot]
		}
		return pidx[a.Person] < pidx[b.Person]
	})
	return out, nil
}

// coverLeftovers: cada sobrante registrado (con receta y porciones > 0) cubre primero las entradas de su
// receta en orden de fecha y tiempo, si no caducó antes de la entrada (§6.4, FIRME). Devuelve las porciones
// cubiertas por entrada y los sobrantes que el plan no alcanza a consumir (FIRME: Q2-44 = A del bloqueo S4-reglas-lista, RESUELTO).
func (e *Engine) coverLeftovers(entries []PlanEntry) ([]int64, []Leftover) {
	covered := make([]int64, len(entries))
	rest := []Leftover{}
	for _, it := range e.Inventory {
		if it.Recipe == "" || it.QuantityMilli <= 0 {
			continue
		}
		left := it.QuantityMilli
		for i, en := range entries {
			if left == 0 {
				break
			}
			if en.EatingOut || en.Recipe != it.Recipe || (it.ExpiresDate != nil && it.ExpiresDate.Before(en.Date)) {
				continue
			}
			take := min(left, *en.ServingsMilli-covered[i])
			covered[i] += take
			left -= take
		}
		if left > 0 {
			rest = append(rest, Leftover{it.Recipe, left})
		}
	}
	return covered, rest
}

// UnplannedSlots cuenta pares día × tiempo obligatorio en que alguna persona con plan no tiene receta ni
// «fuera»; máximo 21 (K43, §5.2).
func (e *Engine) UnplannedSlots(monday Date) (int, error) {
	entries, err := e.weekEntries(monday)
	if err != nil {
		return 0, err
	}
	has := map[string]bool{}
	for _, en := range entries {
		has[en.Date.String()+"|"+string(en.Slot)+"|"+en.Person] = true
	}
	n := 0
	for d := 0; d < 7; d++ {
		day := monday.AddDays(d).String()
		for _, s := range mandatorySlots {
			for _, p := range e.Persons {
				if p.HasMealPlan && !has[day+"|"+string(s)+"|"+p.Code] {
					n++
					break
				}
			}
		}
	}
	return n, nil
}

// PlanCell es una celda del plan (C1 MealPlanWeek.cells).
type PlanCell struct {
	Date              Date     `json:"date"`
	Slot              SlotCode `json:"slot"`
	Person            string   `json:"person"`
	Recipe            *string  `json:"recipe"`
	ServingsMilli     *int64   `json:"servings_milli"`
	EatingOut         bool     `json:"eating_out"`
	CostCents         *int64   `json:"cost_cents"`
	Kcal              *int64   `json:"kcal"`
	PriceStore        *string  `json:"price_store"`
	PriceObservedDate *Date    `json:"price_observed_date"`
}

// PersonDay son las kcal de una persona en un día contra su meta que rige.
type PersonDay struct {
	Person          string `json:"person"`
	Kcal            *int64 `json:"kcal"`
	DailyKcalTarget *int64 `json:"daily_kcal_target"`
	DifferenceKcal  *int64 `json:"difference_kcal"`
}

// PlanDay es un día de la semana.
type PlanDay struct {
	Date      Date        `json:"date"`
	CostCents *int64      `json:"cost_cents"`
	Persons   []PersonDay `json:"persons"`
}

// Leftover es un sobrante registrado que el plan no consume (receta y porciones × 1000).
type Leftover struct {
	Recipe        string `json:"recipe"`
	ServingsMilli int64  `json:"servings_milli"`
}

// WeekPlan es la salida de week_plan (§5.3).
type WeekPlan struct {
	Week           string     `json:"week"`
	Cells          []PlanCell `json:"cells"`
	Days           []PlanDay  `json:"days"`
	WeekCostCents  *int64     `json:"week_cost_cents"`
	UnplannedSlots int        `json:"unplanned_slots"`
	EatingOutSlots int        `json:"eating_out_slots"`
	Leftovers      []Leftover `json:"leftovers"`
}

func addOpt(a, b *int64) *int64 {
	if a == nil || b == nil {
		return nil
	}
	return ptr(*a + *b)
}

// WeekPlanFor arma la semana: celda = HALF_UP del costo exacto a sus porciones (§8.2) y de sus kcal (§9.3);
// «fuera» sin costo de recetas ni kcal y fuera de las sumas (supuesto S-8); día = Σ celdas; semana = Σ días.
func (e *Engine) WeekPlanFor(monday Date) (out *WeekPlan, err error) {
	defer catch(&err)
	entries, err := e.weekEntries(monday)
	if err != nil {
		return nil, err
	}
	unplanned, err := e.UnplannedSlots(monday)
	if err != nil {
		return nil, err
	}
	_, rest := e.coverLeftovers(entries)
	out = &WeekPlan{Week: WeekCode(monday), Cells: []PlanCell{}, Days: []PlanDay{}, UnplannedSlots: unplanned, Leftovers: rest}
	dayCost := map[string]*int64{}
	kcal := map[string]*int64{}
	for _, en := range entries {
		c := PlanCell{Date: en.Date, Slot: en.Slot, Person: en.Person, ServingsMilli: en.ServingsMilli, EatingOut: en.EatingOut}
		if en.EatingOut {
			out.EatingOutSlots++
			out.Cells = append(out.Cells, c)
			continue
		}
		r, _ := e.recipe(en.Recipe)
		c.Recipe = ptr(en.Recipe)
		ev, err := e.evalRecipe(r, *en.ServingsMilli)
		if err != nil {
			return nil, err
		}
		c.CostCents, c.Kcal = roundOpt(ev.cost, 1), roundOpt(ev.kcal, 1000)
		if ev.cost != nil && len(ev.ings) > 0 {
			store, oldest := ev.ings[0].quote.Store, ev.ings[0].quote.ObservedDate
			for _, ie := range ev.ings[1:] {
				if ie.quote.Store != store {
					store = ""
				}
				if ie.quote.ObservedDate.Before(oldest) {
					oldest = ie.quote.ObservedDate
				}
			}
			if store != "" {
				c.PriceStore = &store
			}
			c.PriceObservedDate = &oldest
		}
		d := en.Date.String()
		if _, ok := dayCost[d]; !ok {
			dayCost[d] = ptr(int64(0))
		}
		dayCost[d] = addOpt(dayCost[d], c.CostCents)
		k := d + "|" + en.Person
		if _, ok := kcal[k]; !ok {
			kcal[k] = ptr(int64(0))
		}
		kcal[k] = addOpt(kcal[k], c.Kcal)
		out.Cells = append(out.Cells, c)
	}
	week := ptr(int64(0))
	for i := 0; i < 7; i++ {
		day := monday.AddDays(i)
		pd := PlanDay{Date: day, CostCents: ptr(int64(0)), Persons: []PersonDay{}}
		if v, ok := dayCost[day.String()]; ok {
			pd.CostCents = v
		}
		week = addOpt(week, pd.CostCents)
		for _, p := range e.Persons {
			if !p.HasMealPlan {
				continue
			}
			x := PersonDay{Person: p.Code, Kcal: ptr(int64(0))}
			if v, ok := kcal[day.String()+"|"+p.Code]; ok {
				x.Kcal = v
			}
			if h, err := e.HealthFor(p.Code); err == nil {
				x.DailyKcalTarget = h.DailyKcalTarget
			}
			if x.Kcal != nil && x.DailyKcalTarget != nil {
				x.DifferenceKcal = ptr(*x.Kcal - *x.DailyKcalTarget)
			}
			pd.Persons = append(pd.Persons, x)
		}
		out.Days = append(out.Days, pd)
	}
	out.WeekCostCents = week
	return out, nil
}
