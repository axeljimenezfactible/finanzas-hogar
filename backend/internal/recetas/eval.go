package recetas

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"sync"
)

// Evaluador de fixtures (S4-interfaz §12): lee entradas.datos (§12.1) con campos estrictos y resuelve
// cada consulta (§12.2). Lo usa el subcomando «s4 dump» de backend/cmd/motores.

type inGroup struct {
	Code            *string `json:"code" s4:"req"`
	Name            string  `json:"name"`
	PurchaseChannel *string `json:"purchase_channel" s4:"req"`
}

type inProduct struct {
	Code              *string `json:"code" s4:"req"`
	Name              string  `json:"name"`
	BaseUnit          *string `json:"base_unit" s4:"req"`
	Group             *string `json:"group"`
	GramsPerUnitMilli *int64  `json:"grams_per_unit_milli"`
}

type inPresentation struct {
	Code      *string `json:"code" s4:"req"`
	Product   *string `json:"product" s4:"req"`
	Name      string  `json:"name"`
	SizeMilli *int64  `json:"size_milli" s4:"req"`
}

type inCatalog struct {
	Groups        []inGroup        `json:"groups"`
	Products      []inProduct      `json:"products"`
	Presentations []inPresentation `json:"presentations"`
}

type inIngredient struct {
	Product       *string `json:"product" s4:"req"`
	Presentation  *string `json:"presentation" s4:"req"`
	QuantityMilli *int64  `json:"quantity_milli" s4:"req"`
	Optional      bool    `json:"optional"`
}

type inRecipe struct {
	Code         *string        `json:"code" s4:"req"`
	Name         string         `json:"name"`
	Slots        []SlotCode     `json:"slots"`
	BaseServings *int           `json:"base_servings" s4:"req"`
	PrepMinutes  *int           `json:"prep_minutes"`
	CookMinutes  *int           `json:"cook_minutes"`
	Origin       *string        `json:"origin"`
	IsActive     *bool          `json:"is_active"`
	Ingredients  []inIngredient `json:"ingredients" s4:"req"`
}

type inNutrition struct {
	Product          *string `json:"product" s4:"req"`
	Presentation     *string `json:"presentation"`
	KcalMilliPer100g *int64  `json:"kcal_milli_per_100g" s4:"req"`
	Source           *string `json:"source" s4:"req"`
	SourceRef        *string `json:"source_ref" s4:"req"`
	Status           *string `json:"status" s4:"req"`
}

type inPerson struct {
	Code        *string `json:"code" s4:"req"`
	HasMealPlan *bool   `json:"has_meal_plan" s4:"req"`
}

type inStore struct {
	Code      *string  `json:"code" s4:"req"`
	Kind      *string  `json:"kind" s4:"req"`
	Groups    []string `json:"groups"`
	NextVisit *Date    `json:"next_visit"`
}

type inPlan struct {
	Date          *Date     `json:"date" s4:"req"`
	Slot          *SlotCode `json:"slot" s4:"req"`
	Person        *string   `json:"person" s4:"req"`
	Recipe        *string   `json:"recipe"`
	ServingsMilli *int64    `json:"servings_milli"`
	EatingOut     bool      `json:"eating_out"`
}

type inInventory struct {
	Code              *string `json:"code" s4:"req"`
	Product           *string `json:"product"`
	Recipe            *string `json:"recipe"`
	Mode              *string `json:"mode" s4:"req"`
	QuantityMilli     *int64  `json:"quantity_milli" s4:"req"`
	CountedDate       *Date   `json:"counted_date"`
	RateQuantityMilli *int64  `json:"rate_quantity_milli"`
	RatePeriodDays    *int64  `json:"rate_period_days"`
	ReorderLevelMilli *int64  `json:"reorder_level_milli"`
	ExpiresDate       *Date   `json:"expires_date"`
}

type inChoice struct {
	Product       *string `json:"product" s4:"req"`
	Kind          *string `json:"kind" s4:"req"`
	Store         *string `json:"store"`
	Presentation  *string `json:"presentation"`
	QuantityMilli *int64  `json:"quantity_milli"`
	IsBought      bool    `json:"is_bought"`
}

type inQuote struct {
	Product      *string `json:"product" s4:"req"`
	Store        *string `json:"store" s4:"req"`
	Presentation *string `json:"presentation" s4:"req"`
	PriceCents   *int64  `json:"price_cents" s4:"req"`
	SizeMilli    *int64  `json:"size_milli" s4:"req"`
	ObservedDate *Date   `json:"observed_date" s4:"req"`
	IsStale      bool    `json:"is_stale"`
	Recommended  *bool   `json:"recommended" s4:"req"`
}

type inProfile struct {
	Person              *string `json:"person" s4:"req"`
	HeightMM            *int64  `json:"height_mm"`
	BirthDate           *Date   `json:"birth_date"`
	Sex                 *string `json:"sex"`
	Activity            *string `json:"activity"`
	Goal                *string `json:"goal"`
	AdjustLoseKcal      *int64  `json:"adjust_lose_kcal"`
	AdjustGainKcal      *int64  `json:"adjust_gain_kcal"`
	DailyKcalTarget     *int64  `json:"daily_kcal_target"`
	DietaryRestrictions *string `json:"dietary_restrictions"`
}

type inWeight struct {
	Person  *string `json:"person" s4:"req"`
	Date    *Date   `json:"date" s4:"req"`
	WeightG *int64  `json:"weight_g" s4:"req"`
}

type inSettings struct {
	HouseAdjustLoseKcal *int64 `json:"house_adjust_lose_kcal"`
	HouseAdjustGainKcal *int64 `json:"house_adjust_gain_kcal"`
}

type input struct {
	Today          *Date         `json:"today" s4:"req"`
	CatalogFile    *string       `json:"catalog_file"`
	Catalog        *inCatalog    `json:"catalog"`
	RecipeBookFile *string       `json:"recipe_book_file"`
	Recipes        []inRecipe    `json:"recipes"`
	NutritionFile  *string       `json:"nutrition_file"`
	Nutrition      []inNutrition `json:"nutrition"`
	Persons        []inPerson    `json:"persons" s4:"req"`
	Stores         []inStore     `json:"stores"`
	MealPlan       []inPlan      `json:"meal_plan"`
	Inventory      []inInventory `json:"inventory"`
	ListChoices    []inChoice    `json:"list_choices"`
	Quotes         []inQuote     `json:"quotes"`
	HealthProfiles []inProfile   `json:"health_profiles"`
	Weights        []inWeight    `json:"weights"`
	Settings       *inSettings   `json:"settings"`
}

// decodeStrict decodifica sin campos desconocidos (unknown_field) y exige los obligatorios (missing_field).
func decodeStrict(raw []byte, v any) error {
	if len(bytes.TrimSpace(raw)) == 0 || string(bytes.TrimSpace(raw)) == "null" {
		raw = []byte("{}")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		var ce *codeError
		switch {
		case errors.As(err, &ce):
			return err
		case strings.Contains(err.Error(), "unknown field"):
			return errf(ErrUnknownField, "%v", err)
		}
		return errf(ErrInvalidArgument, "%v", err)
	}
	return checkRequired(reflect.ValueOf(v))
}

// checkRequired recorre structs y listas: un campo con etiqueta s4:"req" nulo es missing_field.
func checkRequired(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			return checkRequired(v.Elem())
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := checkRequired(v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := v.Field(i)
			if !t.Field(i).IsExported() {
				continue
			}
			if t.Field(i).Tag.Get("s4") == "req" && (f.Kind() == reflect.Pointer || f.Kind() == reflect.Slice) && f.IsNil() {
				return errf(ErrMissingField, "%s", strings.Split(t.Field(i).Tag.Get("json"), ",")[0])
			}
			if err := checkRequired(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var sourceLabel = map[string]string{"usda_fdc": "USDA", "open_food_facts": "OFF", "manual": "manual", "mx_table": "mx_table", "store_label": "store_label"}
var unitEN = map[string]Unit{"g": UnitG, "ml": UnitML, "piece": UnitPiece}
var storeKinds = map[string]bool{"chain": true, "branch": true, "market": true}
var enumOK = map[string]map[string]bool{
	"sex": {"female": true, "male": true}, "goal": {"lose": true, "maintain": true, "gain": true},
	"status": {"proposed": true, "confirmed": true}, "kind": {KindCalculated: true, KindManual: true},
	"mode": {ModeStock: true, ModeRate: true},
}

func checkEnum(set string, v *string) error {
	if v != nil && !enumOK[set][*v] {
		return errf(ErrInvalidArgument, "%s %q", set, *v)
	}
	return nil
}

func buildCatalog(c *inCatalog) (*Catalog, error) {
	var gs []Group
	var ps []Product
	var xs []Presentation
	for _, g := range c.Groups {
		gs = append(gs, Group{Code: *g.Code, Name: g.Name, Channel: Channel(*g.PurchaseChannel)})
	}
	for _, p := range c.Products {
		u, ok := unitEN[*p.BaseUnit]
		if !ok {
			return nil, errf(ErrInvalidArgument, "base_unit %q", *p.BaseUnit)
		}
		ps = append(ps, Product{Code: *p.Code, Name: p.Name, BaseUnit: u, Group: str(p.Group), GramsPerUnitMilli: p.GramsPerUnitMilli})
	}
	for _, x := range c.Presentations {
		xs = append(xs, Presentation{Code: *x.Code, Product: *x.Product, Name: x.Name, SizeMilli: *x.SizeMilli})
	}
	return NewCatalog(gs, ps, xs)
}

// BuildEngine arma el motor desde entradas.datos de una fixture (§12.1).
func BuildEngine(raw json.RawMessage) (e *Engine, err error) {
	defer catch(&err)
	var in input
	if err := decodeStrict(raw, &in); err != nil {
		return nil, err
	}
	e = &Engine{Today: *in.Today}
	switch {
	case in.CatalogFile != nil && in.Catalog != nil:
		return nil, errf(ErrInvalidArgument, "catalog_file y catalog a la vez")
	case in.CatalogFile != nil:
		e.Catalog, err = LoadCatalogFile(*in.CatalogFile)
	case in.Catalog != nil:
		e.Catalog, err = buildCatalog(in.Catalog)
	default:
		return nil, errf(ErrMissingField, "catalog")
	}
	if err != nil {
		return nil, err
	}
	switch {
	case in.NutritionFile != nil && in.Nutrition != nil:
		return nil, errf(ErrInvalidArgument, "nutrition_file y nutrition a la vez")
	case in.NutritionFile != nil:
		e.Nutrition, err = LoadNutritionCSV(*in.NutritionFile, e.Catalog, Review)
	default:
		var fs []NutritionFact
		for _, n := range in.Nutrition {
			if _, ok := sourceLabel[*n.Source]; !ok {
				return nil, errf(ErrInvalidArgument, "source %q", *n.Source)
			}
			if err := checkEnum("status", n.Status); err != nil {
				return nil, err
			}
			if e.Catalog.Products[*n.Product] == nil {
				return nil, errf(ErrUnknownProduct, "%q", *n.Product)
			}
			fs = append(fs, NutritionFact{Product: *n.Product, Presentation: str(n.Presentation), KcalMilliPer100g: *n.KcalMilliPer100g,
				Source: *n.Source, SourceLabel: sourceLabel[*n.Source], SourceRef: *n.SourceRef, Proposed: *n.Status == "proposed"})
		}
		e.Nutrition, err = NewNutrition(fs, Review)
	}
	if err != nil {
		return nil, err
	}
	if in.RecipeBookFile == nil && in.Recipes == nil {
		return nil, errf(ErrMissingField, "recipes")
	}
	if in.RecipeBookFile != nil {
		book, err := LoadRecipeBook(*in.RecipeBookFile, e.Catalog, e.Nutrition, Review)
		if err != nil {
			return nil, err
		}
		if len(book.Rejected) > 0 {
			return nil, errf(book.Rejected[0].Err, "receta %q", book.Rejected[0].Recipe)
		}
		e.Recipes = book.Recipes
	}
	for _, r := range in.Recipes {
		if err := e.addInputRecipe(r); err != nil {
			return nil, err
		}
	}
	return e, e.fillRest(&in)
}

func (e *Engine) addInputRecipe(r inRecipe) error {
	if _, err := e.recipe(*r.Code); err == nil {
		return errf(ErrInvalidArgument, "receta repetida %q", *r.Code)
	}
	rc := &Recipe{Code: *r.Code, Name: r.Name, Slots: r.Slots, BaseServings: *r.BaseServings, CookMinutes: r.CookMinutes,
		Origin: str(r.Origin), IsActive: r.IsActive == nil || *r.IsActive}
	if r.PrepMinutes != nil {
		rc.PrepMinutes = *r.PrepMinutes
	}
	if rc.BaseServings < 1 {
		return errf(ErrInvalidServings, "base_servings de %q", rc.Code)
	}
	for _, s := range rc.Slots {
		if _, ok := slotOrder[s]; !ok {
			return errf(ErrInvalidArgument, "tiempo %q", s)
		}
	}
	for i, in := range r.Ingredients {
		if e.Catalog.Products[*in.Product] == nil {
			return errf(ErrUnknownProduct, "%q", *in.Product)
		}
		if _, err := e.Catalog.presentationOf(*in.Product, *in.Presentation); err != nil {
			return err
		}
		if *in.QuantityMilli <= 0 {
			return errf(ErrInvalidArgument, "cantidad de %q", *in.Product)
		}
		rc.Ingredients = append(rc.Ingredients, Ingredient{Product: *in.Product, Presentation: *in.Presentation,
			QuantityMilli: *in.QuantityMilli, Optional: in.Optional, SortOrder: i})
	}
	e.Recipes = append(e.Recipes, rc)
	return nil
}

func (e *Engine) fillRest(in *input) error {
	for _, p := range in.Persons {
		e.Persons = append(e.Persons, Person{Code: *p.Code, HasMealPlan: *p.HasMealPlan})
	}
	for _, s := range in.Stores {
		if !storeKinds[*s.Kind] || (*s.Kind == "market") != (s.NextVisit != nil) {
			return errf(ErrInvalidArgument, "tienda %q: kind y next_visit", *s.Code)
		}
		e.Stores = append(e.Stores, Store{Code: *s.Code, Kind: *s.Kind, Groups: s.Groups, NextVisit: s.NextVisit})
	}
	for _, m := range in.MealPlan {
		e.Plan = append(e.Plan, PlanEntry{Date: *m.Date, Slot: *m.Slot, Person: *m.Person, Recipe: str(m.Recipe), ServingsMilli: m.ServingsMilli, EatingOut: m.EatingOut})
	}
	for _, it := range in.Inventory {
		if err := checkEnum("mode", it.Mode); err != nil {
			return err
		}
		bad := (it.Product == nil) == (it.Recipe == nil) || *it.QuantityMilli < 0
		if it.Product != nil && e.Catalog.Products[*it.Product] == nil {
			return errf(ErrUnknownProduct, "%q", *it.Product)
		}
		if *it.Mode == ModeRate {
			bad = bad || it.CountedDate == nil || it.RateQuantityMilli == nil || it.RatePeriodDays == nil || *it.RateQuantityMilli <= 0 || *it.RatePeriodDays <= 0
		}
		if bad {
			return errf(ErrInvalidArgument, "artículo de inventario %q", *it.Code)
		}
		e.Inventory = append(e.Inventory, InventoryItem{Code: *it.Code, Product: str(it.Product), Recipe: str(it.Recipe), Mode: *it.Mode,
			QuantityMilli: *it.QuantityMilli, CountedDate: it.CountedDate, RateQuantityMilli: it.RateQuantityMilli,
			RatePeriodDays: it.RatePeriodDays, ReorderLevelMilli: it.ReorderLevelMilli, ExpiresDate: it.ExpiresDate})
	}
	calc := map[string]bool{}
	for _, c := range in.ListChoices {
		if err := checkEnum("kind", c.Kind); err != nil {
			return err
		}
		if e.Catalog.Products[*c.Product] == nil {
			return errf(ErrUnknownProduct, "%q", *c.Product)
		}
		if (*c.Kind == KindManual && c.QuantityMilli == nil) || (c.QuantityMilli != nil && *c.QuantityMilli < 0) || (*c.Kind == KindCalculated && calc[*c.Product]) {
			return errf(ErrInvalidArgument, "elección de %q", *c.Product)
		}
		calc[*c.Product] = calc[*c.Product] || *c.Kind == KindCalculated
		e.Choices = append(e.Choices, ListChoice{Product: *c.Product, Kind: *c.Kind, Store: str(c.Store), Presentation: str(c.Presentation), QuantityMilli: c.QuantityMilli, IsBought: c.IsBought})
	}
	var ql QuoteList
	for _, q := range in.Quotes {
		ql = append(ql, Quote{Product: *q.Product, Store: *q.Store, Presentation: *q.Presentation, PriceCents: *q.PriceCents,
			SizeMilli: *q.SizeMilli, ObservedDate: *q.ObservedDate, IsStale: q.IsStale, Recommended: *q.Recommended})
	}
	e.Prices = ql
	for _, h := range in.HealthProfiles {
		for _, v := range []*int64{h.AdjustLoseKcal, h.AdjustGainKcal} {
			if v != nil && *v < 0 {
				return errf(ErrInvalidArgument, "ajuste negativo de %q", *h.Person)
			}
		}
		if err := errors.Join(checkEnum("sex", h.Sex), checkEnum("goal", h.Goal)); err != nil {
			return err
		}
		if (h.DailyKcalTarget != nil && (*h.DailyKcalTarget < 500 || *h.DailyKcalTarget > 10000)) || (h.HeightMM != nil && *h.HeightMM <= 0) {
			return errf(ErrInvalidArgument, "perfil de %q", *h.Person)
		}
		e.Profiles = append(e.Profiles, HealthProfile{Person: *h.Person, HeightMM: h.HeightMM, BirthDate: h.BirthDate, Sex: h.Sex,
			Activity: h.Activity, Goal: h.Goal, AdjustLoseKcal: h.AdjustLoseKcal, AdjustGainKcal: h.AdjustGainKcal,
			DailyKcalTarget: h.DailyKcalTarget, DietaryRestrictions: h.DietaryRestrictions})
	}
	for _, w := range in.Weights {
		if *w.WeightG <= 0 {
			return errf(ErrInvalidArgument, "peso de %q", *w.Person)
		}
		e.Weights = append(e.Weights, WeightEntry{Person: *w.Person, Date: *w.Date, WeightG: *w.WeightG})
	}
	if s := in.Settings; s != nil {
		if s.HouseAdjustLoseKcal != nil {
			e.Settings.HouseAdjustLoseKcal = *s.HouseAdjustLoseKcal
		}
		if s.HouseAdjustGainKcal != nil {
			e.Settings.HouseAdjustGainKcal = *s.HouseAdjustGainKcal
		}
		if e.Settings.HouseAdjustLoseKcal < 0 || e.Settings.HouseAdjustGainKcal < 0 {
			return errf(ErrInvalidArgument, "ajustes de la casa negativos")
		}
	}
	return nil
}

type argRecipe struct {
	Recipe        *string `json:"recipe" s4:"req"`
	ServingsMilli *int64  `json:"servings_milli" s4:"req"`
}
type argDate struct {
	Date *string `json:"date" s4:"req"`
}
type argWeek struct {
	Week *string `json:"week" s4:"req"`
}
type argPerson struct {
	Person *string `json:"person" s4:"req"`
}
type argNone struct{}

var (
	cacheMu sync.Mutex
	cache   = map[string]*Engine{}
)

// Eval resuelve una consulta de fixture: funcion y argumentos de §12.2 sobre entradas.datos.
func Eval(raw json.RawMessage, fn string, args json.RawMessage) (out any, err error) {
	defer catch(&err)
	cacheMu.Lock()
	e, ok := cache[string(raw)]
	cacheMu.Unlock()
	if !ok {
		if e, err = BuildEngine(raw); err != nil {
			return nil, err
		}
		cacheMu.Lock()
		cache[string(raw)] = e
		cacheMu.Unlock()
	}
	week := func() (Date, error) {
		var a argWeek
		if err := decodeStrict(args, &a); err != nil {
			return Date{}, err
		}
		return ParseWeek(*a.Week)
	}
	switch fn {
	case "scale_recipe", "recipe_detail":
		var a argRecipe
		if err := decodeStrict(args, &a); err != nil {
			return nil, err
		}
		if fn == "scale_recipe" {
			return e.ScaleRecipe(*a.Recipe, *a.ServingsMilli)
		}
		return e.RecipeDetailFor(*a.Recipe, *a.ServingsMilli)
	case "default_week":
		var a argDate
		if err := decodeStrict(args, &a); err != nil {
			return nil, err
		}
		d, err := ParseDate(*a.Date)
		if err != nil {
			return nil, err
		}
		m := DefaultWeek(d)
		return map[string]any{"week": WeekCode(m), "monday": m}, nil
	case "unplanned_slots", "week_plan", "shopping_list":
		m, err := week()
		if err != nil {
			return nil, err
		}
		switch fn {
		case "unplanned_slots":
			n, err := e.UnplannedSlots(m)
			if err != nil {
				return nil, err
			}
			return map[string]any{"unplanned_slots": n}, nil
		case "week_plan":
			return e.WeekPlanFor(m)
		}
		return e.ShoppingListFor(m)
	case "inventory_status", "nutrition_coverage":
		if err := decodeStrict(args, &argNone{}); err != nil {
			return nil, err
		}
		if fn == "inventory_status" {
			return e.InventoryStatus()
		}
		return NutritionCoverage(e.Recipes, e.Nutrition), nil
	case "health_profile":
		var a argPerson
		if err := decodeStrict(args, &a); err != nil {
			return nil, err
		}
		return e.HealthFor(*a.Person)
	}
	return nil, errf(ErrInvalidArgument, "función desconocida %q", fn)
}
