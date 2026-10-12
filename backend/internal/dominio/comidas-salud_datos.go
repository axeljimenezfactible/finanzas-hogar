package dominio

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// Comidas y salud (agente api-comidas-salud): repositorios de meal_slot, meal_plan_entry, recipe,
// recipe_ingredient, recipe_meal_slot, health_profile y weight_entry, y la lectura de los datos que la API
// convierte en la entrada del motor S4 (internal/recetas) y de S5 (internal/precios). Aquí no se calcula nada.

type csScanner interface{ Scan(dest ...any) error }

func csQuery[T any](ctx context.Context, q Querier, scan func(csScanner) (T, error), query string, args ...any) ([]T, error) {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func csOne[T any](ctx context.Context, q Querier, scan func(csScanner) (T, error), notFound, query string, args ...any) (T, error) {
	v, err := scan(q.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return v, NotFound(notFound)
	}
	return v, err
}

func csBool(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func csExec(ctx context.Context, q Querier, query string, args ...any) error {
	_, err := q.ExecContext(ctx, query, args...)
	return err
}

// ---- Tiempos de comida y personas ----

// MealSlotRow es un renglón de meal_slot.
type MealSlotRow struct {
	ID                   int64
	Code, Name           string
	SortOrder            int64
	IsOptional           bool
	CreatedAt, UpdatedAt string
}

func scanMealSlot(s csScanner) (r MealSlotRow, err error) {
	err = s.Scan(&r.ID, &r.Code, &r.Name, &r.SortOrder, &r.IsOptional, &r.CreatedAt, &r.UpdatedAt)
	return
}

// ListMealSlots lee los tiempos de comida; orderLimit es el ORDER BY … LIMIT … de la página.
func ListMealSlots(ctx context.Context, q Querier, orderLimit string) ([]MealSlotRow, error) {
	return csQuery(ctx, q, scanMealSlot, `SELECT id, code, name, sort_order, is_optional, created_at, updated_at FROM meal_slot ms `+orderLimit)
}

// MealPersonRow es una persona con su participación (las externas no tienen plan ni salud).
type MealPersonRow struct {
	ID                  int64
	Code, Participation string
}

// HasMealPlan dice si la persona entra al plan de comidas y a salud (titular y co-titular).
func (p MealPersonRow) HasMealPlan() bool { return p.Participation != "external" }

// ListMealPersons lee las personas en orden de id.
func ListMealPersons(ctx context.Context, q Querier) ([]MealPersonRow, error) {
	return csQuery(ctx, q, func(s csScanner) (p MealPersonRow, err error) {
		err = s.Scan(&p.ID, &p.Code, &p.Participation)
		return
	}, `SELECT id, code, participation FROM person ORDER BY id`)
}

// ---- Plan de comidas ----

// MealPlanEntryRow es un renglón de meal_plan_entry.
type MealPlanEntryRow struct {
	ID                      int64
	PlanDate                string
	MealSlotID, PersonID    int64
	RecipeID, ServingsMilli sql.NullInt64
	IsEatingOut             bool
	Notes                   sql.NullString
	CreatedAt, UpdatedAt    string
}

const mealPlanCols = `id, plan_date, meal_slot_id, person_id, recipe_id, servings_milli, is_eating_out, notes, created_at, updated_at`

func scanMealPlanEntry(s csScanner) (e MealPlanEntryRow, err error) {
	err = s.Scan(&e.ID, &e.PlanDate, &e.MealSlotID, &e.PersonID, &e.RecipeID, &e.ServingsMilli, &e.IsEatingOut, &e.Notes, &e.CreatedAt, &e.UpdatedAt)
	return
}

// GetMealPlanEntry lee una entrada del plan (404 si no existe).
func GetMealPlanEntry(ctx context.Context, q Querier, id int64) (MealPlanEntryRow, error) {
	return csOne(ctx, q, scanMealPlanEntry, "No encontramos esa casilla del plan.", `SELECT `+mealPlanCols+` FROM meal_plan_entry WHERE id = ?`, id)
}

// ListMealPlanEntries lee las entradas entre dos fechas (incluidas).
func ListMealPlanEntries(ctx context.Context, q Querier, from, to string) ([]MealPlanEntryRow, error) {
	return csQuery(ctx, q, scanMealPlanEntry, `SELECT `+mealPlanCols+` FROM meal_plan_entry WHERE plan_date BETWEEN ? AND ? ORDER BY plan_date, meal_slot_id, person_id, id`, from, to)
}

// MealPlanCell lee las entradas de una casilla (día × tiempo × persona).
func MealPlanCell(ctx context.Context, q Querier, date string, slotID, personID int64) ([]MealPlanEntryRow, error) {
	return csQuery(ctx, q, scanMealPlanEntry, `SELECT `+mealPlanCols+` FROM meal_plan_entry WHERE plan_date = ? AND meal_slot_id = ? AND person_id = ? ORDER BY id`, date, slotID, personID)
}

// SaveMealPlanEntry inserta (ID = 0) o actualiza la entrada.
func SaveMealPlanEntry(ctx context.Context, q Querier, e *MealPlanEntryRow) error {
	if e.ID == 0 {
		return q.QueryRowContext(ctx, `INSERT INTO meal_plan_entry (plan_date, meal_slot_id, person_id, recipe_id, servings_milli, is_eating_out, notes) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`,
			e.PlanDate, e.MealSlotID, e.PersonID, e.RecipeID, e.ServingsMilli, csBool(e.IsEatingOut), e.Notes).Scan(&e.ID)
	}
	return csExec(ctx, q, `UPDATE meal_plan_entry SET recipe_id = ?, servings_milli = ?, is_eating_out = ?, notes = ? WHERE id = ?`,
		e.RecipeID, e.ServingsMilli, csBool(e.IsEatingOut), e.Notes, e.ID)
}

// DeleteMealPlanEntry borra una entrada.
func DeleteMealPlanEntry(ctx context.Context, q Querier, id int64) error {
	return csExec(ctx, q, `DELETE FROM meal_plan_entry WHERE id = ?`, id)
}

// MealPlanDates son las fechas distintas con plan hasta until (incluida), en orden.
func MealPlanDates(ctx context.Context, q Querier, until string) ([]string, error) {
	return csQuery(ctx, q, func(s csScanner) (d string, err error) { err = s.Scan(&d); return },
		`SELECT DISTINCT plan_date FROM meal_plan_entry WHERE plan_date <= ? ORDER BY plan_date`, until)
}

// RecipeUseCounts cuenta las entradas del plan por receta entre dos fechas (incluidas).
func RecipeUseCounts(ctx context.Context, q Querier, from, to string) (map[int64]int64, error) {
	type pair struct{ id, n int64 }
	ps, err := csQuery(ctx, q, func(s csScanner) (p pair, err error) { err = s.Scan(&p.id, &p.n); return },
		`SELECT recipe_id, count(*) FROM meal_plan_entry WHERE recipe_id IS NOT NULL AND plan_date BETWEEN ? AND ? GROUP BY recipe_id`, from, to)
	out := map[int64]int64{}
	for _, p := range ps {
		out[p.id] = p.n
	}
	return out, err
}

// RecipePlanInfo: si la receta tiene alguna entrada en el plan y la primera fecha planeada desde from.
func RecipePlanInfo(ctx context.Context, q Querier, recipeID int64, from string) (anyPlan bool, next sql.NullString, err error) {
	err = q.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM meal_plan_entry WHERE recipe_id = ?), (SELECT min(plan_date) FROM meal_plan_entry WHERE recipe_id = ? AND plan_date >= ?)`,
		recipeID, recipeID, from).Scan(&anyPlan, &next)
	return
}

// ---- Recetas ----

// RecipeIngredientRow es un renglón de recipe_ingredient.
type RecipeIngredientRow struct {
	ID, ProductID  int64
	PresentationID sql.NullInt64
	QuantityMilli  int64
	DisplayText    sql.NullString
	IsOptional     bool
	SortOrder      int64
}

// RecipeRow es una receta con sus tiempos aptos y sus ingredientes en orden.
type RecipeRow struct {
	ID                        int64
	Name                      string
	Description, Instructions sql.NullString
	BaseServings              int64
	PrepMinutes, CookMinutes  sql.NullInt64
	Origin                    string
	AuthorPersonID            sql.NullInt64
	IsActive                  bool
	Notes                     sql.NullString
	CreatedAt, UpdatedAt      string
	SlotIDs                   []int64
	Ingredients               []RecipeIngredientRow
}

const recipeCols = `id, name, description, instructions, base_servings, prep_minutes, cook_minutes, origin, author_person_id, is_active, notes, created_at, updated_at`

func scanRecipe(s csScanner) (r RecipeRow, err error) {
	err = s.Scan(&r.ID, &r.Name, &r.Description, &r.Instructions, &r.BaseServings, &r.PrepMinutes, &r.CookMinutes, &r.Origin,
		&r.AuthorPersonID, &r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return
}

func recipeChildren(ctx context.Context, q Querier, rs []RecipeRow, id int64) error {
	type slot struct{ r, s int64 }
	ss, err := csQuery(ctx, q, func(s csScanner) (x slot, err error) { err = s.Scan(&x.r, &x.s); return },
		`SELECT rms.recipe_id, rms.meal_slot_id FROM recipe_meal_slot rms JOIN meal_slot ms ON ms.id = rms.meal_slot_id WHERE ? = 0 OR rms.recipe_id = ? ORDER BY rms.recipe_id, ms.sort_order`, id, id)
	if err != nil {
		return err
	}
	type ing struct {
		r int64
		RecipeIngredientRow
	}
	is, err := csQuery(ctx, q, func(s csScanner) (x ing, err error) {
		err = s.Scan(&x.r, &x.ID, &x.ProductID, &x.PresentationID, &x.QuantityMilli, &x.DisplayText, &x.IsOptional, &x.SortOrder)
		return
	}, `SELECT recipe_id, id, product_id, presentation_id, quantity_milli, display_text, is_optional, sort_order FROM recipe_ingredient WHERE ? = 0 OR recipe_id = ? ORDER BY recipe_id, sort_order, id`, id, id)
	if err != nil {
		return err
	}
	idx := map[int64]int{}
	for i := range rs {
		idx[rs[i].ID] = i
		rs[i].SlotIDs, rs[i].Ingredients = []int64{}, []RecipeIngredientRow{}
	}
	for _, x := range ss {
		if i, ok := idx[x.r]; ok {
			rs[i].SlotIDs = append(rs[i].SlotIDs, x.s)
		}
	}
	for _, x := range is {
		if i, ok := idx[x.r]; ok {
			rs[i].Ingredients = append(rs[i].Ingredients, x.RecipeIngredientRow)
		}
	}
	return nil
}

// ListRecipeRows lee todas las recetas (activas y archivadas) con tiempos e ingredientes.
func ListRecipeRows(ctx context.Context, q Querier) ([]RecipeRow, error) {
	rs, err := csQuery(ctx, q, scanRecipe, `SELECT `+recipeCols+` FROM recipe ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return rs, recipeChildren(ctx, q, rs, 0)
}

// GetRecipeRow lee una receta (404 si no existe).
func GetRecipeRow(ctx context.Context, q Querier, id int64) (RecipeRow, error) {
	r, err := csOne(ctx, q, scanRecipe, "No encontramos esa receta.", `SELECT `+recipeCols+` FROM recipe WHERE id = ?`, id)
	if err != nil {
		return r, err
	}
	rs := []RecipeRow{r}
	err = recipeChildren(ctx, q, rs, id)
	return rs[0], err
}

// SaveRecipe inserta (ID = 0) o actualiza la receta; con slots o ings reemplaza tiempos o ingredientes.
func SaveRecipe(ctx context.Context, q Querier, r *RecipeRow, slots, ings bool) error {
	var err error
	if r.ID == 0 {
		err = q.QueryRowContext(ctx, `INSERT INTO recipe (name, description, instructions, base_servings, prep_minutes, cook_minutes, origin, author_person_id, is_active, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
			strings.TrimSpace(r.Name), r.Description, r.Instructions, r.BaseServings, r.PrepMinutes, r.CookMinutes, r.Origin, r.AuthorPersonID, csBool(r.IsActive), r.Notes).Scan(&r.ID)
	} else {
		err = csExec(ctx, q, `UPDATE recipe SET name = ?, description = ?, instructions = ?, base_servings = ?, prep_minutes = ?, cook_minutes = ?, is_active = ?, notes = ? WHERE id = ?`,
			strings.TrimSpace(r.Name), r.Description, r.Instructions, r.BaseServings, r.PrepMinutes, r.CookMinutes, csBool(r.IsActive), r.Notes, r.ID)
	}
	if err != nil {
		return err
	}
	if slots {
		if err := csExec(ctx, q, `DELETE FROM recipe_meal_slot WHERE recipe_id = ?`, r.ID); err != nil {
			return err
		}
		for _, s := range r.SlotIDs {
			if err := csExec(ctx, q, `INSERT INTO recipe_meal_slot (recipe_id, meal_slot_id) VALUES (?, ?)`, r.ID, s); err != nil {
				return err
			}
		}
	}
	if ings {
		if err := csExec(ctx, q, `DELETE FROM recipe_ingredient WHERE recipe_id = ?`, r.ID); err != nil {
			return err
		}
		for i := range r.Ingredients {
			in := &r.Ingredients[i]
			if err := q.QueryRowContext(ctx, `INSERT INTO recipe_ingredient (recipe_id, product_id, presentation_id, quantity_milli, display_text, is_optional, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id`,
				r.ID, in.ProductID, in.PresentationID, in.QuantityMilli, in.DisplayText, csBool(in.IsOptional), in.SortOrder).Scan(&in.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// ArchiveRecipe deja la receta inactiva (el plan la usa: no se borra su historia).
func ArchiveRecipe(ctx context.Context, q Querier, id int64) error {
	return csExec(ctx, q, `UPDATE recipe SET is_active = 0 WHERE id = ?`, id)
}

// DeleteRecipe borra la receta (sus ingredientes y tiempos caen en cascada).
func DeleteRecipe(ctx context.Context, q Querier, id int64) error {
	return csExec(ctx, q, `DELETE FROM recipe WHERE id = ?`, id)
}

// ---- Perfil de salud y peso ----

// HealthProfileRow es un renglón de health_profile.
type HealthProfileRow struct {
	ID, PersonID                           int64
	HeightMM                               sql.NullInt64
	BirthDate, Sex, ActivityLevel          sql.NullString
	DailyKcalTarget                        sql.NullInt64
	DietaryRestrictions, Notes, Goal       sql.NullString
	GoalAdjustLoseKcal, GoalAdjustGainKcal sql.NullInt64
	CreatedAt, UpdatedAt                   string
}

const healthCols = `hp.id, hp.person_id, hp.height_mm, hp.birth_date, hp.sex, hp.activity_level, hp.daily_kcal_target, hp.dietary_restrictions, hp.notes, hp.goal, hp.goal_adjust_lose_kcal, hp.goal_adjust_gain_kcal, hp.created_at, hp.updated_at`

func scanHealth(s csScanner) (h HealthProfileRow, err error) {
	err = s.Scan(&h.ID, &h.PersonID, &h.HeightMM, &h.BirthDate, &h.Sex, &h.ActivityLevel, &h.DailyKcalTarget, &h.DietaryRestrictions,
		&h.Notes, &h.Goal, &h.GoalAdjustLoseKcal, &h.GoalAdjustGainKcal, &h.CreatedAt, &h.UpdatedAt)
	return
}

// ListHealthProfileRows lee los perfiles; orderLimit es el ORDER BY … LIMIT … de la página ("" = por id).
func ListHealthProfileRows(ctx context.Context, q Querier, orderLimit string) ([]HealthProfileRow, error) {
	if orderLimit == "" {
		orderLimit = "ORDER BY hp.id"
	}
	return csQuery(ctx, q, scanHealth, `SELECT `+healthCols+` FROM health_profile hp `+orderLimit)
}

// GetHealthProfileRow lee un perfil (404 si no existe).
func GetHealthProfileRow(ctx context.Context, q Querier, id int64) (HealthProfileRow, error) {
	return csOne(ctx, q, scanHealth, "No encontramos ese perfil de salud.", `SELECT `+healthCols+` FROM health_profile hp WHERE hp.id = ?`, id)
}

// SaveHealthProfile inserta (ID = 0) o actualiza el perfil.
func SaveHealthProfile(ctx context.Context, q Querier, h *HealthProfileRow) error {
	args := []any{h.HeightMM, h.BirthDate, h.Sex, h.ActivityLevel, h.DailyKcalTarget, h.DietaryRestrictions, h.Notes, h.Goal, h.GoalAdjustLoseKcal, h.GoalAdjustGainKcal}
	if h.ID == 0 {
		return q.QueryRowContext(ctx, `INSERT INTO health_profile (height_mm, birth_date, sex, activity_level, daily_kcal_target, dietary_restrictions, notes, goal, goal_adjust_lose_kcal, goal_adjust_gain_kcal, person_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
			append(args, h.PersonID)...).Scan(&h.ID)
	}
	return csExec(ctx, q, `UPDATE health_profile SET height_mm = ?, birth_date = ?, sex = ?, activity_level = ?, daily_kcal_target = ?, dietary_restrictions = ?, notes = ?, goal = ?, goal_adjust_lose_kcal = ?, goal_adjust_gain_kcal = ? WHERE id = ?`,
		append(args, h.ID)...)
}

// WeightEntryRow es un renglón de weight_entry.
type WeightEntryRow struct {
	ID, PersonID               int64
	MeasuredDate               string
	WeightG                    int64
	Notes                      sql.NullString
	CreatedAt, UpdatedAt       string
	CapturedAt, IdempotencyKey sql.NullString
}

// WeightFilter filtra el historial de peso (vacío = todo).
type WeightFilter struct {
	PersonIDs []int64
	From, To  string
}

const weightCols = `w.id, w.person_id, w.measured_date, w.weight_g, w.notes, w.created_at, w.updated_at, w.captured_at, w.idempotency_key`

func scanWeight(s csScanner) (w WeightEntryRow, err error) {
	err = s.Scan(&w.ID, &w.PersonID, &w.MeasuredDate, &w.WeightG, &w.Notes, &w.CreatedAt, &w.UpdatedAt, &w.CapturedAt, &w.IdempotencyKey)
	return
}

// ListWeightEntryRows lee el historial con filtro; orderLimit es el ORDER BY … LIMIT … ("" = por fecha y id).
func ListWeightEntryRows(ctx context.Context, q Querier, f WeightFilter, orderLimit string) ([]WeightEntryRow, error) {
	where, args := []string{"1 = 1"}, []any{}
	if len(f.PersonIDs) > 0 {
		where = append(where, "w.person_id IN (?"+strings.Repeat(", ?", len(f.PersonIDs)-1)+")")
		for _, id := range f.PersonIDs {
			args = append(args, id)
		}
	}
	if f.From != "" {
		where, args = append(where, "w.measured_date >= ?"), append(args, f.From)
	}
	if f.To != "" {
		where, args = append(where, "w.measured_date <= ?"), append(args, f.To)
	}
	if orderLimit == "" {
		orderLimit = "ORDER BY w.measured_date, w.id"
	}
	return csQuery(ctx, q, scanWeight, `SELECT `+weightCols+` FROM weight_entry w WHERE `+strings.Join(where, " AND ")+" "+orderLimit, args...)
}

// GetWeightEntryRow lee un peso (404 si no existe).
func GetWeightEntryRow(ctx context.Context, q Querier, id int64) (WeightEntryRow, error) {
	return csOne(ctx, q, scanWeight, "No encontramos ese registro de peso.", `SELECT `+weightCols+` FROM weight_entry w WHERE w.id = ?`, id)
}

// WeightEntryByKey busca el peso guardado con esa Idempotency-Key (K9; ok = false si no hay).
func WeightEntryByKey(ctx context.Context, q Querier, key string) (WeightEntryRow, bool, error) {
	w, err := scanWeight(q.QueryRowContext(ctx, `SELECT `+weightCols+` FROM weight_entry w WHERE w.idempotency_key = ?`, key))
	if errors.Is(err, sql.ErrNoRows) {
		return w, false, nil
	}
	return w, err == nil, err
}

// SaveWeightEntry inserta (ID = 0, con captured_at e idempotency_key) o actualiza fecha, peso y nota.
func SaveWeightEntry(ctx context.Context, q Querier, w *WeightEntryRow) error {
	if w.ID == 0 {
		return q.QueryRowContext(ctx, `INSERT INTO weight_entry (person_id, measured_date, weight_g, notes, captured_at, idempotency_key) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`,
			w.PersonID, w.MeasuredDate, w.WeightG, w.Notes, w.CapturedAt, w.IdempotencyKey).Scan(&w.ID)
	}
	return csExec(ctx, q, `UPDATE weight_entry SET measured_date = ?, weight_g = ?, notes = ? WHERE id = ?`, w.MeasuredDate, w.WeightG, w.Notes, w.ID)
}

// DeleteWeightEntry borra un peso.
func DeleteWeightEntry(ctx context.Context, q Querier, id int64) error {
	return csExec(ctx, q, `DELETE FROM weight_entry WHERE id = ?`, id)
}

// ---- Datos de catálogo, nutrición, tiendas, precios y sobrantes para los motores S4 y S5 ----

// MealGroupRow es un grupo de compra.
type MealGroupRow struct {
	ID            int64
	Name, Channel string
	SortOrder     int64
}

// MealProductRow es un producto genérico.
type MealProductRow struct {
	ID                         int64
	Name, BaseUnit             string
	GramsPerUnitMilli, GroupID sql.NullInt64
}

// MealPresentationRow es una presentación comprable.
type MealPresentationRow struct {
	ID, ProductID int64
	Name          string
	SizeMilli     int64
}

// MealNutritionRow es el dato nutricional preferido (is_preferred = 1).
type MealNutritionRow struct {
	ProductID                 int64
	PresentationID            sql.NullInt64
	KcalMilliPer100g          int64
	ProteinMg, FatMg, CarbsMg sql.NullInt64
	Source                    string
	SourceRef, License        sql.NullString
	CreatedAt                 string
}

// MealStoreRow es una tienda activa con sus grupos de compra.
type MealStoreRow struct {
	ID                                                 int64
	Kind, Name                                         string
	ParentID, VisitWeekday, MembershipFeeCents         sql.NullInt64
	ProfecoName, Municipality                          sql.NullString
	RequiresMembership, HasMembership, AssistedReading bool
	GroupIDs                                           []int64
}

// MealPriceRow es una observación de precio de una tienda activa, observada hasta hoy.
type MealPriceRow struct {
	PresentationID, StoreID, PriceCents int64
	ObservedDate, Source                string
	IsPromotional                       bool
	SourceRef                           sql.NullString
	CreatedAt                           string
}

// MealLeftoverRow es un sobrante de receta en inventario (porciones × 1000).
type MealLeftoverRow struct {
	ID, RecipeID, QuantityMilli int64
	ExpiresDate                 sql.NullString
}

// MealFoodData reúne lo que leen los motores para costo y kcal.
type MealFoodData struct {
	Groups        []MealGroupRow
	Products      []MealProductRow
	Presentations []MealPresentationRow
	Nutrition     []MealNutritionRow
	Stores        []MealStoreRow
	Prices        []MealPriceRow
	Leftovers     []MealLeftoverRow
}

// LoadMealFoodData lee catálogo, nutrición preferida, tiendas activas, precios hasta today y sobrantes.
func LoadMealFoodData(ctx context.Context, q Querier, today string) (d MealFoodData, err error) {
	if d.Groups, err = csQuery(ctx, q, func(s csScanner) (g MealGroupRow, err error) {
		err = s.Scan(&g.ID, &g.Name, &g.Channel, &g.SortOrder)
		return
	}, `SELECT id, name, purchase_channel, sort_order FROM product_group ORDER BY sort_order, id`); err != nil {
		return
	}
	if d.Products, err = csQuery(ctx, q, func(s csScanner) (p MealProductRow, err error) {
		err = s.Scan(&p.ID, &p.Name, &p.BaseUnit, &p.GramsPerUnitMilli, &p.GroupID)
		return
	}, `SELECT id, name, base_unit, grams_per_unit_milli, product_group_id FROM product ORDER BY id`); err != nil {
		return
	}
	if d.Presentations, err = csQuery(ctx, q, func(s csScanner) (p MealPresentationRow, err error) {
		err = s.Scan(&p.ID, &p.ProductID, &p.Name, &p.SizeMilli)
		return
	}, `SELECT id, product_id, name, size_milli FROM presentation ORDER BY id`); err != nil {
		return
	}
	if d.Nutrition, err = csQuery(ctx, q, func(s csScanner) (n MealNutritionRow, err error) {
		err = s.Scan(&n.ProductID, &n.PresentationID, &n.KcalMilliPer100g, &n.ProteinMg, &n.FatMg, &n.CarbsMg, &n.Source, &n.SourceRef, &n.License, &n.CreatedAt)
		return
	}, `SELECT product_id, presentation_id, kcal_milli_per_100g, protein_mg_per_100g, fat_mg_per_100g, carbs_mg_per_100g, source, source_ref, source_license, created_at FROM nutrition_fact WHERE is_preferred = 1 ORDER BY id`); err != nil {
		return
	}
	if d.Stores, err = csQuery(ctx, q, func(s csScanner) (x MealStoreRow, err error) {
		err = s.Scan(&x.ID, &x.Kind, &x.Name, &x.ParentID, &x.VisitWeekday, &x.MembershipFeeCents, &x.ProfecoName, &x.Municipality, &x.RequiresMembership, &x.HasMembership, &x.AssistedReading)
		return
	}, `SELECT id, kind, name, parent_id, visit_weekday, membership_annual_fee_cents, profeco_name, municipality, requires_membership, has_membership, assisted_reading_enabled FROM store WHERE is_active = 1 ORDER BY id`); err != nil {
		return
	}
	type sg struct{ s, g int64 }
	sgs, err := csQuery(ctx, q, func(s csScanner) (x sg, err error) { err = s.Scan(&x.s, &x.g); return },
		`SELECT store_id, product_group_id FROM store_product_group ORDER BY store_id, product_group_id`)
	if err != nil {
		return
	}
	for i := range d.Stores {
		d.Stores[i].GroupIDs = []int64{}
		for _, x := range sgs {
			if x.s == d.Stores[i].ID {
				d.Stores[i].GroupIDs = append(d.Stores[i].GroupIDs, x.g)
			}
		}
	}
	if d.Prices, err = csQuery(ctx, q, func(s csScanner) (p MealPriceRow, err error) {
		err = s.Scan(&p.PresentationID, &p.StoreID, &p.PriceCents, &p.ObservedDate, &p.Source, &p.IsPromotional, &p.SourceRef, &p.CreatedAt)
		return
	}, `SELECT p.presentation_id, p.store_id, p.price_cents, p.observed_date, p.source, p.is_promotional, p.source_ref, p.created_at FROM price p JOIN store s ON s.id = p.store_id AND s.is_active = 1 WHERE p.observed_date <= ? ORDER BY p.id`, today); err != nil {
		return
	}
	d.Leftovers, err = csQuery(ctx, q, func(s csScanner) (l MealLeftoverRow, err error) {
		err = s.Scan(&l.ID, &l.RecipeID, &l.QuantityMilli, &l.ExpiresDate)
		return
	}, `SELECT id, recipe_id, quantity_milli, expires_date FROM inventory_item WHERE recipe_id IS NOT NULL ORDER BY id`)
	return
}
