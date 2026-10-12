// Paquete comidas: operaciones de C1 de comidas (plan semanal, recetas) y salud (perfiles, peso, kcal y
// ajuste de la casa) del agente api-comidas-salud. Los cálculos de costo, kcal, escalado, IMC y meta
// sugerida los hace el motor S4 (internal/recetas) con los precios de S5 (internal/precios); aquí solo se
// leen datos, se validan entradas y se convierten tipos.
package comidas

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
	"finanzas-hogar/internal/recetas"
)

// Register registra las 26 operaciones de comidas y salud.
func Register(reg *api.Registry) {
	a := &area{reg: reg}
	for op, h := range map[string]api.HandlerFunc{
		"listMealSlots": a.listMealSlots, "getMealPlanWeek": a.getMealPlanWeek, "getMealPlanMonth": a.getMealPlanMonth,
		"getMealPlanTrends": a.getMealPlanTrends, "createMealPlanEntries": a.createMealPlanEntries,
		"updateMealPlanEntry": a.updateMealPlanEntry, "deleteMealPlanEntry": a.deleteMealPlanEntry,
		"copyPreviousWeek": a.copyPreviousWeek, "listRecipes": a.listRecipes, "createRecipe": a.createRecipe,
		"suggestRecipeWithAi": a.suggestRecipeWithAi, "getRecipe": a.getRecipe, "updateRecipe": a.updateRecipe,
		"deleteRecipe": a.deleteRecipe, "getRecipeUsage": a.getRecipeUsage, "duplicateRecipe": a.duplicateRecipe,
		"listHealthProfiles": a.listHealthProfiles, "getTodayKcal": a.getTodayKcal, "listWeightEntries": a.listWeightEntries,
		"createWeightEntry": a.createWeightEntry, "updateWeightEntry": a.updateWeightEntry, "deleteWeightEntry": a.deleteWeightEntry,
		"createHealthProfile": a.createHealthProfile, "updateHealthProfile": a.updateHealthProfile,
		"getHealthSettings": a.getHealthSettings, "updateHealthSettings": a.updateHealthSettings,
	} {
		reg.Handle(op, h)
	}
}

type area struct{ reg *api.Registry }

func (a *area) db() *sql.DB         { return a.reg.Deps().DB }
func (a *area) today() dominio.Date { return dominio.Today(a.reg.Deps().Clock) }

// audit escribe audit_log con el renglón posterior (salvo en delete).
func (a *area) audit(ctx context.Context, tx *sql.Tx, r *http.Request, action, table string, id int64, before json.RawMessage) error {
	var after json.RawMessage
	if action != "delete" {
		var err error
		if after, err = dominio.RowJSON(ctx, tx, table, id); err != nil {
			return err
		}
	}
	return dominio.Audit(ctx, tx, api.Actor(r), action, table, id, before, after, api.RequestID(r))
}

// ---- Conversión de tipos ----

func code(prefix string, id int64) string { return prefix + strconv.FormatInt(id, 10) }

func idOf(prefix, c string) int64 {
	n, _ := strconv.ParseInt(strings.TrimPrefix(c, prefix), 10, 64)
	return n
}

func rdate(s string) recetas.Date { d, _ := recetas.ParseDate(s); return d }

func ddate(d recetas.Date) dominio.Date { x, _ := dominio.ParseDate(d.String()); return x }

func mondayOf(d dominio.Date) dominio.Date { return d.AddDays(-((int(d.Weekday()) + 6) % 7)) }

func nint(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}

func nstr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	v := n.String
	return &v
}

func toNI(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

func toNS(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func mergeI(n api.Nullable[int64], cur sql.NullInt64) sql.NullInt64 {
	if !n.Set {
		return cur
	}
	return sql.NullInt64{Int64: n.Value, Valid: n.Valid}
}

func mergeS(n api.Nullable[string], cur sql.NullString) sql.NullString {
	if !n.Set {
		return cur
	}
	return sql.NullString{String: n.Value, Valid: n.Valid}
}

func isNotFound(err error) bool {
	e, ok := dominio.AsError(err)
	return ok && e.Kind == dominio.KindNotFound
}

func issue(field, c, msg string) dominio.FieldIssue {
	return dominio.FieldIssue{Field: field, Code: c, Message: msg}
}

func externalPerson(field string) error {
	return &api.Problem{Status: http.StatusUnprocessableEntity, Code: "external_person",
		Message: "Solo Axel y Ximena tienen plan de comidas y perfil de salud.",
		Fields:  []api.FieldError{{Field: field, Code: "not_allowed", Message: "Persona externa."}}}
}

// ---- Entrada de los motores ----

type loadOpts struct {
	food     bool         // catálogo, nutrición, precios, recetas y sobrantes
	from, to dominio.Date // plan entre estas fechas; cero = sin plan
}

type motor struct {
	e        *recetas.Engine
	persons  []dominio.MealPersonRow
	personID map[string]int64
	slots    []dominio.MealSlotRow
	slotCode map[int64]recetas.SlotCode
	slotID   map[recetas.SlotCode]int64
	recipes  map[int64]dominio.RecipeRow
	nutSrc   map[string]string
	entries  []dominio.MealPlanEntryRow
	profiles []dominio.HealthProfileRow
}

func (m *motor) person(id int64) (dominio.MealPersonRow, bool) {
	for _, p := range m.persons {
		if p.ID == id {
			return p, true
		}
	}
	return dominio.MealPersonRow{}, false
}

func (m *motor) personCode(id int64) string { p, _ := m.person(id); return p.Code }

// loadMotor lee de SQLite lo que necesitan S4 y S5 (C-05: «hoy» en America/Mexico_City lo pone la API).
func (a *area) loadMotor(ctx context.Context, q dominio.Querier, o loadOpts) (*motor, error) {
	today := a.today()
	m := &motor{e: &recetas.Engine{Today: rdate(today.String())}, personID: map[string]int64{}, slotCode: map[int64]recetas.SlotCode{},
		slotID: map[recetas.SlotCode]int64{}, recipes: map[int64]dominio.RecipeRow{}, nutSrc: map[string]string{}}
	var err error
	if m.persons, err = dominio.ListMealPersons(ctx, q); err != nil {
		return nil, err
	}
	for _, p := range m.persons {
		m.personID[p.Code] = p.ID
		m.e.Persons = append(m.e.Persons, recetas.Person{Code: p.Code, HasMealPlan: p.HasMealPlan()})
	}
	if m.slots, err = dominio.ListMealSlots(ctx, q, "ORDER BY ms.sort_order"); err != nil {
		return nil, err
	}
	for _, s := range m.slots {
		m.slotCode[s.ID], m.slotID[recetas.SlotCode(s.Code)] = recetas.SlotCode(s.Code), s.ID
	}
	if err := m.loadHealth(ctx, q); err != nil {
		return nil, err
	}
	if o.food {
		if err := m.loadFood(ctx, q, today); err != nil {
			return nil, err
		}
	}
	if !o.from.IsZero() {
		if m.entries, err = dominio.ListMealPlanEntries(ctx, q, o.from.String(), o.to.String()); err != nil {
			return nil, err
		}
		for _, en := range m.entries {
			pe := recetas.PlanEntry{Date: rdate(en.PlanDate), Slot: m.slotCode[en.MealSlotID], Person: m.personCode(en.PersonID),
				EatingOut: en.IsEatingOut, ServingsMilli: nint(en.ServingsMilli)}
			if en.RecipeID.Valid {
				pe.Recipe = code("r", en.RecipeID.Int64)
			}
			m.e.Plan = append(m.e.Plan, pe)
		}
	}
	return m, nil
}

func (m *motor) loadHealth(ctx context.Context, q dominio.Querier) error {
	var err error
	if m.profiles, err = dominio.ListHealthProfileRows(ctx, q, ""); err != nil {
		return err
	}
	for _, h := range m.profiles {
		m.e.Profiles = append(m.e.Profiles, toProfile(m.personCode(h.PersonID), h))
	}
	ws, err := dominio.ListWeightEntryRows(ctx, q, dominio.WeightFilter{}, "")
	if err != nil {
		return err
	}
	for _, w := range ws {
		m.e.Weights = append(m.e.Weights, recetas.WeightEntry{Person: m.personCode(w.PersonID), Date: rdate(w.MeasuredDate), WeightG: w.WeightG})
	}
	lose, err := dominio.SettingInt(ctx, q, "health.goal_adjust_lose_kcal")
	if err != nil {
		return err
	}
	gain, err := dominio.SettingInt(ctx, q, "health.goal_adjust_gain_kcal")
	m.e.Settings = recetas.Settings{HouseAdjustLoseKcal: lose, HouseAdjustGainKcal: gain}
	return err
}

func toProfile(person string, h dominio.HealthProfileRow) recetas.HealthProfile {
	p := recetas.HealthProfile{Person: person, HeightMM: nint(h.HeightMM), Sex: nstr(h.Sex), Activity: nstr(h.ActivityLevel), Goal: nstr(h.Goal),
		AdjustLoseKcal: nint(h.GoalAdjustLoseKcal), AdjustGainKcal: nint(h.GoalAdjustGainKcal), DailyKcalTarget: nint(h.DailyKcalTarget),
		DietaryRestrictions: nstr(h.DietaryRestrictions)}
	if h.BirthDate.Valid {
		d := rdate(h.BirthDate.String)
		p.BirthDate = &d
	}
	return p
}

func (m *motor) loadFood(ctx context.Context, q dominio.Querier, today dominio.Date) error {
	d, err := dominio.LoadMealFoodData(ctx, q, today.String())
	if err != nil {
		return err
	}
	var (
		groups   []recetas.Group
		products []recetas.Product
		pres     []recetas.Presentation
		pgroups  []precios.Group
		pprods   []precios.Product
		ppres    []precios.Presentation
	)
	for _, g := range d.Groups {
		groups = append(groups, recetas.Group{Code: code("g", g.ID), Name: g.Name, Channel: recetas.Channel(g.Channel), SortOrder: int(g.SortOrder)})
		pgroups = append(pgroups, precios.Group{Code: code("g", g.ID), Name: g.Name, PurchaseChannel: g.Channel})
	}
	for _, p := range d.Products {
		grp := ""
		if p.GroupID.Valid {
			grp = code("g", p.GroupID.Int64)
		}
		products = append(products, recetas.Product{ID: p.ID, Code: code("p", p.ID), Name: p.Name, BaseUnit: recetas.Unit(p.BaseUnit), GramsPerUnitMilli: nint(p.GramsPerUnitMilli), Group: grp})
		pprods = append(pprods, precios.Product{Code: code("p", p.ID), Name: p.Name, BaseUnit: p.BaseUnit, Group: grp})
	}
	for _, x := range d.Presentations {
		pres = append(pres, recetas.Presentation{ID: x.ID, Code: code("pr", x.ID), Product: code("p", x.ProductID), Name: x.Name, SizeMilli: x.SizeMilli})
		ppres = append(ppres, precios.Presentation{Code: code("pr", x.ID), Product: code("p", x.ProductID), Name: x.Name, SizeMilli: x.SizeMilli})
	}
	if m.e.Catalog, err = recetas.NewCatalog(groups, products, pres); err != nil {
		return err
	}
	var facts []recetas.NutritionFact
	for _, n := range d.Nutrition {
		f := recetas.NutritionFact{Product: code("p", n.ProductID), KcalMilliPer100g: n.KcalMilliPer100g, ProteinMg: nint(n.ProteinMg), FatMg: nint(n.FatMg),
			CarbsMg: nint(n.CarbsMg), Source: n.Source, SourceRef: n.SourceRef.String, License: n.License.String, ObservedDate: rdate(n.CreatedAt[:10])}
		if n.PresentationID.Valid {
			f.Presentation = code("pr", n.PresentationID.Int64)
		}
		if f.SourceRef == "" { // fuente manual sin cita (el esquema la permite): la cita es la fuente misma
			f.SourceRef = n.Source
		}
		m.nutSrc[f.Product+"|"+f.Presentation] = n.Source
		facts = append(facts, f)
	}
	if m.e.Nutrition, err = recetas.NewNutrition(facts, recetas.Production); err != nil {
		return err
	}
	rs, err := dominio.ListRecipeRows(ctx, q)
	if err != nil {
		return err
	}
	for _, r := range rs {
		m.recipes[r.ID] = r
		m.e.Recipes = append(m.e.Recipes, m.toRecipe(r))
	}
	for _, l := range d.Leftovers {
		it := recetas.InventoryItem{Code: code("i", l.ID), Recipe: code("r", l.RecipeID), Mode: recetas.ModeStock, QuantityMilli: l.QuantityMilli}
		if l.ExpiresDate.Valid {
			x := rdate(l.ExpiresDate.String)
			it.ExpiresDate = &x
		}
		m.e.Inventory = append(m.e.Inventory, it)
	}
	return m.loadPrices(ctx, q, d, pgroups, pprods, ppres, today)
}

func (m *motor) loadPrices(ctx context.Context, q dominio.Querier, d dominio.MealFoodData, g []precios.Group, p []precios.Product, x []precios.Presentation, today dominio.Date) error {
	cat, err := precios.NewCatalog(g, p, x)
	if err != nil {
		return err
	}
	var stores []precios.Store
	for _, s := range d.Stores {
		st := precios.Store{Code: code("s", s.ID), Kind: s.Kind, Name: s.Name, ProfecoName: s.ProfecoName.String, Municipality: s.Municipality.String,
			VisitWeekday: int(s.VisitWeekday.Int64), RequiresMembership: s.RequiresMembership, HasMembership: s.HasMembership,
			AssistedReading: s.AssistedReading, MembershipFeeCents: nint(s.MembershipFeeCents)}
		if s.ParentID.Valid {
			st.Parent = code("s", s.ParentID.Int64)
		}
		for _, gid := range s.GroupIDs {
			st.Groups = append(st.Groups, code("g", gid))
		}
		stores = append(stores, st)
	}
	var obs []precios.PriceObs
	for _, pr := range d.Prices {
		od, err := precios.ParseDate(pr.ObservedDate)
		if err != nil {
			return err
		}
		ca, _ := dominio.ParseDBTime(pr.CreatedAt)
		obs = append(obs, precios.PriceObs{Presentation: code("pr", pr.PresentationID), Store: code("s", pr.StoreID), PriceCents: pr.PriceCents,
			ObservedDate: od, Source: pr.Source, IsPromotional: pr.IsPromotional, SourceRef: pr.SourceRef.String, CreatedAt: ca})
	}
	var set precios.Settings
	stale, err := dominio.SettingInt(ctx, q, "price.stale_after_days")
	if err != nil {
		return err
	}
	zero, err := dominio.SettingInt(ctx, q, "price.zero_weight_after_days")
	if err != nil {
		return err
	}
	set.StaleAfterDays, set.ZeroWeightAfterDays = int(stale), int(zero)
	if set.AssistedReadingEnabled, err = dominio.SettingBool(ctx, q, "price.assisted_reading_enabled"); err != nil {
		return err
	}
	at, err := precios.ParseDate(today.String())
	if err != nil {
		return err
	}
	book, err := precios.NewBook(cat, stores, obs, set, at)
	if err != nil {
		return err
	}
	m.e.Prices = priceSource{book: book, at: at}
	return nil
}

func (m *motor) toRecipe(r dominio.RecipeRow) *recetas.Recipe {
	x := &recetas.Recipe{ID: r.ID, Code: code("r", r.ID), Name: r.Name, Description: r.Description.String, BaseServings: int(r.BaseServings),
		PrepMinutes: int(r.PrepMinutes.Int64), Origin: r.Origin, IsActive: r.IsActive}
	if r.CookMinutes.Valid {
		v := int(r.CookMinutes.Int64)
		x.CookMinutes = &v
	}
	if r.Instructions.Valid {
		x.Steps = []string{r.Instructions.String}
	}
	for _, s := range r.SlotIDs {
		x.Slots = append(x.Slots, m.slotCode[s])
	}
	for _, in := range r.Ingredients {
		ing := recetas.Ingredient{Product: code("p", in.ProductID), QuantityMilli: in.QuantityMilli, DisplayText: in.DisplayText.String, Optional: in.IsOptional, SortOrder: int(in.SortOrder)}
		if in.PresentationID.Valid {
			ing.Presentation = code("pr", in.PresentationID.Int64)
		}
		x.Ingredients = append(x.Ingredients, ing)
	}
	return x
}

// priceSource conecta S5 con S4 (S4-interfaz §8.1): la cotización recomendada de QuoteFor a la fecha de hoy.
// Las elecciones de tienda o presentación de la lista no se leen aquí (PROVISIONAL, bloqueo S6-comidas-calculos).
type priceSource struct {
	book *precios.Book
	at   time.Time
}

func (s priceSource) Quotes(product string, _ recetas.Date) []recetas.Quote {
	q, ok, err := s.book.QuoteFor(product, "", "", s.at)
	if err != nil || !ok {
		return nil
	}
	return []recetas.Quote{{Product: q.Product, Store: q.Store, Presentation: q.Presentation, PriceCents: q.PriceCents, SizeMilli: q.SizeMilli,
		ObservedDate: rdate(q.ObservedDate.Format(time.DateOnly)), IsStale: q.IsStale, Recommended: q.Recommended}}
}
