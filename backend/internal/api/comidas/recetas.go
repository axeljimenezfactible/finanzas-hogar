package comidas

import (
	"cmp"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/ia"
	"finanzas-hogar/internal/recetas"
)

// uses8w cuenta las entradas del plan por receta en las últimas 8 semanas: de hoy − 55 días a hoy
// (56 días; PROVISIONAL, bloqueo S6-comidas-calculos).
func (a *area) uses8w(ctx context.Context, q dominio.Querier) (map[int64]int64, error) {
	t := a.today()
	return dominio.RecipeUseCounts(ctx, q, t.AddDays(-55).String(), t.String())
}

// apiRecipe convierte la receta y pide al motor S4 costo y kcal por porción a sus porciones base.
func (a *area) apiRecipe(m *motor, r dominio.RecipeRow, uses map[int64]int64) api.Recipe {
	out := api.Recipe{ID: r.ID, Name: r.Name, Description: nstr(r.Description), Instructions: nstr(r.Instructions), BaseServings: r.BaseServings,
		PrepMinutes: nint(r.PrepMinutes), CookMinutes: nint(r.CookMinutes), Origin: r.Origin, AuthorPersonID: nint(r.AuthorPersonID),
		IsActive: r.IsActive, Notes: nstr(r.Notes), CreatedAt: api.APITime(r.CreatedAt), UpdatedAt: api.APITime(r.UpdatedAt), MealSlotIDs: r.SlotIDs}
	n := uses[r.ID]
	out.TimesUsed8w = &n
	if d, err := m.e.RecipeDetailFor(code("r", r.ID), r.BaseServings*1000); err == nil {
		out.CostPerServingCents, out.KcalPerServing = d.CostPerServingCents, d.KcalPerServing
	} else { // dato del recetario que el motor rechaza: los calculados quedan ausentes (nunca 0)
		a.reg.Deps().Log.Warn("comidas: el motor S4 no calculó la receta", "recipe_id", r.ID, "code", recetas.Code(err))
	}
	return out
}

func (a *area) recipeResponse(r *http.Request, id int64) (api.Recipe, error) {
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{food: true})
	if err != nil {
		return api.Recipe{}, err
	}
	uses, err := a.uses8w(ctx, db)
	if err != nil {
		return api.Recipe{}, err
	}
	rr, ok := m.recipes[id]
	if !ok {
		return api.Recipe{}, dominio.NotFound("No encontramos esa receta.")
	}
	return a.apiRecipe(m, rr, uses), nil
}

func cmpOpt(x, y *int64) int {
	switch {
	case x == nil && y == nil:
		return 0
	case x == nil:
		return 1
	case y == nil:
		return -1
	}
	return cmp.Compare(*x, *y)
}

func (a *area) listRecipes(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"name", "cost_per_serving_cents", "kcal_per_serving", "times_used_8w"}, "name")
	if err != nil {
		return err
	}
	qv := r.URL.Query()
	text := strings.ToLower(strings.TrimSpace(qv.Get("q")))
	var slot int64
	if s := qv.Get("meal_slot_id"); s != "" {
		if slot, err = strconv.ParseInt(s, 10, 64); err != nil {
			return api.BadRequest("meal_slot_id debe ser un número entero.")
		}
	}
	origin := qv.Get("origin")
	if origin != "" && origin != recetas.OriginSeed && origin != recetas.OriginHousehold {
		return api.BadRequest("origin debe ser claude_seed o household.")
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{food: true})
	if err != nil {
		return err
	}
	uses, err := a.uses8w(ctx, db)
	if err != nil {
		return err
	}
	items := []api.Recipe{}
	for _, x := range m.e.Recipes {
		rr := m.recipes[x.ID]
		if (text != "" && !strings.Contains(strings.ToLower(rr.Name), text)) || (slot != 0 && !slices.Contains(rr.SlotIDs, slot)) || (origin != "" && rr.Origin != origin) {
			continue
		}
		items = append(items, a.apiRecipe(m, rr, uses))
	}
	slices.SortStableFunc(items, func(x, y api.Recipe) int {
		for _, f := range p.Sort {
			c := 0
			switch f.Name {
			case "name":
				c = strings.Compare(strings.ToLower(x.Name), strings.ToLower(y.Name))
			case "cost_per_serving_cents":
				c = cmpOpt(x.CostPerServingCents, y.CostPerServingCents)
			case "kcal_per_serving":
				c = cmpOpt(x.KcalPerServing, y.KcalPerServing)
			case "times_used_8w":
				c = cmpOpt(x.TimesUsed8w, y.TimesUsed8w)
			case "id":
				c = cmp.Compare(x.ID, y.ID)
			}
			if f.Desc {
				c = -c
			}
			if c != 0 {
				return c
			}
		}
		return cmp.Compare(x.ID, y.ID)
	})
	start := min(p.Offset, len(items))
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items[start:min(p.Offset+p.Limit+1, len(items))], p))
}

// recipeIn son los campos de alta o cambio ya convertidos.
type recipeIn struct {
	ings    []api.RecipeIngredientInput
	ingsSet bool
}

// validateRecipe valida campos con mensajes por campo y luego las reglas del recetario del motor S4 (§3.3, S4-13).
func (a *area) validateRecipe(m *motor, rr *dominio.RecipeRow, in recipeIn) error {
	var v api.Validator
	n := utf8.RuneCountInString(strings.TrimSpace(rr.Name))
	v.Check(n >= 3 && n <= 80, "name", "out_of_range", "El nombre lleva de 3 a 80 caracteres.")
	v.Check(utf8.RuneCountInString(rr.Description.String) <= 200, "description", "out_of_range", "La descripción lleva hasta 200 caracteres.")
	v.Check(rr.BaseServings >= 1 && rr.BaseServings <= 50, "base_servings", "out_of_range", "Las porciones base van de 1 a 50.")
	v.Check(rr.Origin != recetas.OriginSeed || rr.BaseServings == 2, "base_servings", "out_of_range", "El recetario inicial está escrito para 2 porciones.")
	v.Check(!rr.PrepMinutes.Valid || rr.PrepMinutes.Int64 >= 0, "prep_minutes", "out_of_range", "Los minutos no pueden ser negativos.")
	v.Check(!rr.CookMinutes.Valid || rr.CookMinutes.Int64 >= 0, "cook_minutes", "out_of_range", "Los minutos no pueden ser negativos.")
	v.Check(strings.TrimSpace(rr.Instructions.String) != "", "instructions", "required", "Escribe los pasos de preparación.")
	v.Check(len(rr.SlotIDs) > 0, "meal_slot_ids", "required", "Elige para qué tiempos sirve la receta.")
	for _, s := range rr.SlotIDs {
		_, ok := m.slotCode[s]
		v.Check(ok, "meal_slot_ids", "not_found", "Ese tiempo de comida no existe.")
	}
	if in.ingsSet {
		rr.Ingredients = nil
		seen := map[int64]bool{}
		for i, x := range in.ings {
			f := fmt.Sprintf("ingredients[%d]", i)
			pc := code("p", x.ProductID)
			if m.e.Catalog.Products[pc] == nil {
				v.Add(f+".product_id", "not_found", "Ese producto no existe en el catálogo.")
				continue
			}
			v.Check(!seen[x.ProductID], f+".product_id", "conflict", "Ese producto ya está en la receta.")
			seen[x.ProductID] = true
			if x.PresentationID == nil { // el recetario de S4 pide la presentación comprable de cada ingrediente (§3.2)
				v.Add(f+".presentation_id", "required", "Elige la presentación comprable del ingrediente.")
			} else {
				pr := m.e.Catalog.Presentations[code("pr", *x.PresentationID)]
				v.Check(pr != nil && pr.Product == pc, f+".presentation_id", "not_found", "Esa presentación no es de ese producto.")
			}
			v.Check(x.QuantityMilli > 0, f+".quantity_milli", "out_of_range", "La cantidad debe ser mayor que cero.")
			so := int64(i + 1)
			if x.SortOrder != nil {
				so = *x.SortOrder
			}
			rr.Ingredients = append(rr.Ingredients, dominio.RecipeIngredientRow{ProductID: x.ProductID, PresentationID: toNI(x.PresentationID),
				QuantityMilli: x.QuantityMilli, DisplayText: toNS(x.DisplayText), IsOptional: x.IsOptional != nil && *x.IsOptional, SortOrder: so})
		}
	}
	v.Check(len(rr.Ingredients) > 0, "ingredients", "required", "Agrega al menos un ingrediente.")
	if err := v.Err(); err != nil {
		return err
	}
	slices.SortStableFunc(rr.Ingredients, func(x, y dominio.RecipeIngredientRow) int { return cmp.Compare(x.SortOrder, y.SortOrder) })
	if err := recetas.ValidateRecipe(m.toRecipe(*rr), m.e.Catalog); err != nil {
		return dominio.Invalid(issue("ingredients", "invalid_format", "La receta no cumple las reglas del recetario ("+recetas.Code(err)+"). Revisa unidades y presentaciones."))
	}
	return nil
}

func saveRecipeErr(err error) error {
	if dominio.IsUnique(err) {
		return dominio.Conflict("already_exists", "Ya hay una receta con ese nombre. Usa otro.")
	}
	return err
}

func (a *area) createRecipe(w http.ResponseWriter, r *http.Request) error {
	var in api.RecipeCreate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{food: true})
	if err != nil {
		return err
	}
	rr := dominio.RecipeRow{ID: 0, Name: in.Name, Description: toNS(in.Description), Instructions: toNS(in.Instructions), BaseServings: 2,
		PrepMinutes: toNI(in.PrepMinutes), CookMinutes: toNI(in.CookMinutes), Origin: recetas.OriginHousehold, IsActive: true, Notes: toNS(in.Notes), SlotIDs: in.MealSlotIDs}
	if s := api.Session(r); s != nil {
		rr.AuthorPersonID = sql.NullInt64{Int64: s.PersonID, Valid: true}
	}
	if in.BaseServings != nil {
		rr.BaseServings = *in.BaseServings
	}
	if in.IsActive != nil {
		rr.IsActive = *in.IsActive
	}
	if err := a.validateRecipe(m, &rr, recipeIn{ings: in.Ingredients, ingsSet: true}); err != nil {
		return err
	}
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if err := dominio.SaveRecipe(ctx, tx, &rr, true, true); err != nil {
			return saveRecipeErr(err)
		}
		return a.audit(ctx, tx, r, "insert", "recipe", rr.ID, nil)
	}); err != nil {
		return err
	}
	out, err := a.recipeResponse(r, rr.ID)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func (a *area) suggestRecipeWithAi(w http.ResponseWriter, r *http.Request) error {
	var in api.AiRecipeRequest
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	var v api.Validator
	v.Check(strings.TrimSpace(in.Description) != "", "description", "required", "Escribe la idea o los ingredientes que tienes.")
	if in.MealSlotID != nil {
		slots, err := dominio.ListMealSlots(ctx, a.db(), "")
		if err != nil {
			return err
		}
		v.Check(slices.ContainsFunc(slots, func(s dominio.MealSlotRow) bool { return s.ID == *in.MealSlotID }), "meal_slot_id", "not_found", "Ese tiempo de comida no existe.")
	}
	if err := v.Err(); err != nil {
		return err
	}
	off := dominio.Conflict("ai_disabled", "La IA está apagada. Puedes encenderla en Ajustes.")
	reg := a.reg.Deps().IA
	if reg == nil {
		return off
	}
	p, err := reg.Active(ctx)
	if errors.Is(err, ia.ErrDisabled) {
		return off
	} else if err != nil {
		return err
	}
	sug, err := p.SuggestRecipe(ctx, ia.RecipeRequest{Prompt: in.Description, Servings: 2})
	if errors.Is(err, ia.ErrDisabled) {
		return off
	} else if err != nil {
		return dominio.Unavailable("La IA no respondió. Intenta de nuevo más tarde.")
	}
	var out api.RecipeCreate
	if json.Unmarshal(sug.Raw, &out) != nil || strings.TrimSpace(out.Name) == "" {
		return dominio.Unavailable("La IA respondió algo que no pudimos leer. Intenta de nuevo.")
	}
	if in.MealSlotID != nil && len(out.MealSlotIDs) == 0 {
		out.MealSlotIDs = []int64{*in.MealSlotID}
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (a *area) getRecipe(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "recipe_id")
	if err != nil {
		return err
	}
	var sm int64
	if s := r.URL.Query().Get("servings_milli"); s != "" {
		if sm, err = strconv.ParseInt(s, 10, 64); err != nil {
			return api.BadRequest("servings_milli debe ser un número entero (1 porción = 1000).")
		}
		if sm <= 0 {
			return dominio.Invalid(issue("servings_milli", "out_of_range", "Las porciones deben ser mayores que cero."))
		}
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{food: true})
	if err != nil {
		return err
	}
	rr, ok := m.recipes[id]
	if !ok {
		return dominio.NotFound("No encontramos esa receta.")
	}
	if sm == 0 {
		sm = rr.BaseServings * 1000
	}
	uses, err := a.uses8w(ctx, db)
	if err != nil {
		return err
	}
	d, err := m.e.RecipeDetailFor(code("r", id), sm)
	if err != nil {
		return err
	}
	out := api.RecipeDetail{Recipe: a.apiRecipe(m, rr, uses), ServingsMilli: sm, Ingredients: []api.RecipeDetailIngredientsItem{},
		TotalCostCents: d.TotalCostCents, CostPerServingCents: d.CostPerServingCents, TotalKcal: d.TotalKcal, KcalPerServing: d.KcalPerServing}
	for i, x := range d.Ingredients {
		in := rr.Ingredients[i]
		it := api.RecipeDetailIngredientsItem{RecipeIngredientID: in.ID, ProductID: in.ProductID, PresentationID: nint(in.PresentationID),
			QuantityMilli: dominio.RoundHalfUp(x.Quantity.Num, x.Quantity.Den), DisplayText: nstr(in.DisplayText), IsOptional: in.IsOptional,
			CostCents: x.CostCents, Kcal: x.Kcal}
		key := x.Product + "|"
		if in.PresentationID.Valid && m.nutSrc[key+code("pr", in.PresentationID.Int64)] != "" {
			key += code("pr", in.PresentationID.Int64)
		}
		if s, ok := m.nutSrc[key]; ok && x.Kcal != nil {
			it.NutritionSource = &s
		}
		out.Ingredients = append(out.Ingredients, it)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (a *area) updateRecipe(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "recipe_id")
	if err != nil {
		return err
	}
	var in api.RecipeUpdate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if in.Name == nil && !in.Description.Set && !in.Instructions.Set && in.BaseServings == nil && !in.PrepMinutes.Set && !in.CookMinutes.Set &&
		in.IsActive == nil && !in.Notes.Set && in.MealSlotIDs == nil && in.Ingredients == nil {
		return dominio.Invalid(issue("body", "required", "Manda al menos un cambio."))
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{food: true})
	if err != nil {
		return err
	}
	rr, ok := m.recipes[id]
	if !ok {
		return dominio.NotFound("No encontramos esa receta.")
	}
	if in.Name != nil {
		rr.Name = *in.Name
	}
	if in.BaseServings != nil {
		rr.BaseServings = *in.BaseServings
	}
	if in.IsActive != nil {
		rr.IsActive = *in.IsActive
	}
	if in.MealSlotIDs != nil {
		rr.SlotIDs = in.MealSlotIDs
	}
	rr.Description, rr.Instructions, rr.Notes = mergeS(in.Description, rr.Description), mergeS(in.Instructions, rr.Instructions), mergeS(in.Notes, rr.Notes)
	rr.PrepMinutes, rr.CookMinutes = mergeI(in.PrepMinutes, rr.PrepMinutes), mergeI(in.CookMinutes, rr.CookMinutes)
	if err := a.validateRecipe(m, &rr, recipeIn{ings: in.Ingredients, ingsSet: in.Ingredients != nil}); err != nil {
		return err
	}
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "recipe", id)
		if err != nil {
			return err
		}
		if err := dominio.SaveRecipe(ctx, tx, &rr, in.MealSlotIDs != nil, in.Ingredients != nil); err != nil {
			return saveRecipeErr(err)
		}
		return a.audit(ctx, tx, r, "update", "recipe", id, before)
	}); err != nil {
		return err
	}
	out, err := a.recipeResponse(r, id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// deleteRecipe: planeada de hoy en adelante sin confirm_planned → 409 recipe_planned; si el plan la usa, el motor
// decide archivarla (is_active = 0, C-06) y las casillas se conservan; si no, se borra.
func (a *area) deleteRecipe(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "recipe_id")
	if err != nil {
		return err
	}
	confirm := false
	if s := r.URL.Query().Get("confirm_planned"); s != "" {
		if confirm, err = strconv.ParseBool(s); err != nil {
			return api.BadRequest("confirm_planned debe ser true o false.")
		}
	}
	ctx := r.Context()
	err = dominio.WithTx(ctx, a.db(), func(tx *sql.Tx) error {
		rr, err := dominio.GetRecipeRow(ctx, tx, id)
		if err != nil {
			return err
		}
		used, next, err := dominio.RecipePlanInfo(ctx, tx, id, a.today().String())
		if err != nil {
			return err
		}
		if next.Valid && !confirm {
			return dominio.Conflict("recipe_planned", "Esa receta está en el plan desde el "+next.String+". Confirma para quitarla del recetario.")
		}
		before, err := dominio.RowJSON(ctx, tx, "recipe", id)
		if err != nil {
			return err
		}
		m := &motor{slotCode: map[int64]recetas.SlotCode{}}
		if recetas.RemoveRecipe(m.toRecipe(rr), used) == recetas.RemoveArchive {
			if err := dominio.ArchiveRecipe(ctx, tx, id); err != nil {
				return err
			}
			if err := a.audit(ctx, tx, r, "update", "recipe", id, before); err != nil {
				return err
			}
		} else {
			if err := dominio.DeleteRecipe(ctx, tx, id); err != nil {
				return err
			}
			if err := a.audit(ctx, tx, r, "delete", "recipe", id, before); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return api.NoContent(w)
}

func (a *area) getRecipeUsage(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "recipe_id")
	if err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	if _, err := dominio.GetRecipeRow(ctx, db, id); err != nil {
		return err
	}
	uses, err := a.uses8w(ctx, db)
	if err != nil {
		return err
	}
	_, next, err := dominio.RecipePlanInfo(ctx, db, id, a.today().String())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, api.RecipeUsage{TimesUsed8w: uses[id], NextPlanDate: nstr(next)})
}

func (a *area) duplicateRecipe(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "recipe_id")
	if err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{food: true})
	if err != nil {
		return err
	}
	rr, ok := m.recipes[id]
	if !ok {
		return dominio.NotFound("No encontramos esa receta.")
	}
	rr.ID, rr.Name, rr.Origin, rr.AuthorPersonID = 0, rr.Name+" (copia)", recetas.OriginHousehold, sql.NullInt64{}
	if s := api.Session(r); s != nil {
		rr.AuthorPersonID = sql.NullInt64{Int64: s.PersonID, Valid: true}
	}
	if err := a.validateRecipe(m, &rr, recipeIn{}); err != nil {
		return err
	}
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if err := dominio.SaveRecipe(ctx, tx, &rr, true, true); err != nil {
			return saveRecipeErr(err)
		}
		return a.audit(ctx, tx, r, "insert", "recipe", rr.ID, nil)
	}); err != nil {
		return err
	}
	out, err := a.recipeResponse(r, rr.ID)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}
