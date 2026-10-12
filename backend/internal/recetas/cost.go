package recetas

import "math/big"

// Escalado exacto (§4, RN-14) y costo y kcal de una receta (§8, §9.3; S4-02, S4-06, S4-08).

// Scale = quantity_milli × servings_milli / (base_servings × 1000), fracción exacta sin redondeo (§4).
func Scale(qtyMilli int64, baseServings int, servingsMilli int64) (*big.Rat, error) {
	if servingsMilli <= 0 || baseServings <= 0 {
		return nil, errf(ErrInvalidServings, "porciones %d (base %d)", servingsMilli, baseServings)
	}
	return quoR(mulR(r64(qtyMilli), r64(servingsMilli)), r64(int64(baseServings)*1000)), nil
}

// ScaledIngredient es un ingrediente escalado: misma unidad, cantidad exacta.
type ScaledIngredient struct {
	Product  string `json:"product"`
	Unit     Unit   `json:"unit"`
	Quantity Frac   `json:"quantity"`
}

// ScaledRecipe es la salida de scale_recipe (§12.2).
type ScaledRecipe struct {
	Recipe        string             `json:"recipe"`
	ServingsMilli int64              `json:"servings_milli"`
	Ingredients   []ScaledIngredient `json:"ingredients"`
}

// ScaleRecipe escala todos los ingredientes de una receta (opcionales incluidos: el escalado no los filtra).
func (e *Engine) ScaleRecipe(code string, servingsMilli int64) (out *ScaledRecipe, err error) {
	defer catch(&err)
	r, err := e.recipe(code)
	if err != nil {
		return nil, err
	}
	out = &ScaledRecipe{Recipe: code, ServingsMilli: servingsMilli, Ingredients: []ScaledIngredient{}}
	for _, in := range r.Ingredients {
		q, err := Scale(in.QuantityMilli, r.BaseServings, servingsMilli)
		if err != nil {
			return nil, err
		}
		out.Ingredients = append(out.Ingredients, ScaledIngredient{in.Product, e.Catalog.Products[in.Product].BaseUnit, fracOf(q)})
	}
	return out, nil
}

// quotes son las cotizaciones del producto; se valida aquí, al usarlas: exactamente una recomendada y a lo
// más una por (tienda, presentación) (§12.1). Un producto sin cotizaciones queda «sin precio».
func (e *Engine) quotes(product string) ([]Quote, error) {
	if e.Prices == nil {
		return nil, nil
	}
	qs := e.Prices.Quotes(product, e.Today)
	rec, seen := 0, map[string]bool{}
	for _, q := range qs {
		k := q.Store + "|" + q.Presentation
		if seen[k] || q.SizeMilli <= 0 || q.PriceCents < 0 {
			return nil, errf(ErrInvalidArgument, "cotización repetida o inválida de %q", product)
		}
		seen[k] = true
		if q.Recommended {
			rec++
		}
	}
	if len(qs) > 0 && rec != 1 {
		return nil, errf(ErrInvalidArgument, "%q trae %d cotizaciones recomendadas; debe ser exactamente una", product, rec)
	}
	return qs, nil
}

// recommended es la cotización recomendada del producto (nil si no hay).
func (e *Engine) recommended(product string) (*Quote, error) {
	qs, err := e.quotes(product)
	for i := range qs {
		if qs[i].Recommended {
			return &qs[i], nil
		}
	}
	return nil, err
}

// governing es la cotización que rige (§8.1): con elección de tienda o presentación, la que casa (si casan
// varias, la recomendada; si no, invalid_argument; si ninguna, «sin precio» sin caer a la recomendada);
// sin elección, la recomendada.
func (e *Engine) governing(product string, ch *ListChoice) (*Quote, error) {
	if ch == nil || (ch.Store == "" && ch.Presentation == "") {
		return e.recommended(product)
	}
	qs, err := e.quotes(product)
	if err != nil {
		return nil, err
	}
	var match []int
	for i, q := range qs {
		if (ch.Store == "" || q.Store == ch.Store) && (ch.Presentation == "" || q.Presentation == ch.Presentation) {
			match = append(match, i)
		}
	}
	if len(match) == 1 {
		return &qs[match[0]], nil
	}
	for _, i := range match {
		if qs[i].Recommended {
			return &qs[i], nil
		}
	}
	if len(match) > 1 {
		return nil, errf(ErrInvalidArgument, "varias cotizaciones casan con la elección de %q", product)
	}
	return nil, nil
}

// calcChoice es la elección del artículo calculated del producto; rige costos de recetas, plan y lista.
func (e *Engine) calcChoice(product string) *ListChoice {
	for i := range e.Choices {
		if e.Choices[i].Kind == KindCalculated && e.Choices[i].Product == product {
			return &e.Choices[i]
		}
	}
	return nil
}

type ingEval struct {
	ing   Ingredient
	qty   *big.Rat
	cost  *big.Rat // nil = sin precio
	quote *Quote
	kcal  *big.Rat // milésimas de kcal; nil = sin dato
}

type recipeEval struct {
	ings                            []ingEval
	cost, kcal                      *big.Rat // nil si falta algún precio o dato
	missingPrices, missingNutrition []string
}

// evalRecipe calcula costo y kcal exactos para servings_milli; los opcionales no entran (§5.4, S-3).
func (e *Engine) evalRecipe(r *Recipe, sm int64) (*recipeEval, error) {
	ev := &recipeEval{cost: r64(0), kcal: r64(0), missingPrices: []string{}, missingNutrition: []string{}}
	for _, in := range r.Ingredients {
		if in.Optional {
			continue
		}
		q, err := Scale(in.QuantityMilli, r.BaseServings, sm)
		if err != nil {
			return nil, err
		}
		ie := ingEval{ing: in, qty: q}
		quote, err := e.governing(in.Product, e.calcChoice(in.Product))
		if err != nil {
			return nil, err
		}
		if quote != nil {
			ie.quote, ie.cost = quote, mulR(q, frac(quote.PriceCents, quote.SizeMilli))
			if ev.cost != nil {
				ev.cost = addR(ev.cost, ie.cost)
			}
		} else {
			ev.cost = nil
			ev.missingPrices = append(ev.missingPrices, in.Product)
		}
		if f := e.Nutrition.fact(in.Product, in.Presentation); f != nil {
			if k, ok := ingredientKcal(e.Catalog, f, in.Product, q); ok {
				ie.kcal = k
			}
		}
		if ie.kcal != nil {
			if ev.kcal != nil {
				ev.kcal = addR(ev.kcal, ie.kcal)
			}
		} else {
			ev.kcal = nil
			ev.missingNutrition = append(ev.missingNutrition, in.Product)
		}
		ev.ings = append(ev.ings, ie)
	}
	return ev, nil
}

func roundOpt(r *big.Rat, den int64) *int64 {
	if r == nil {
		return nil
	}
	return ptr(roundHalfUp(quoR(r, r64(den))))
}

// DetailIngredient es un renglón de recipe_detail: costo y kcal redondeados solo para mostrar.
type DetailIngredient struct {
	Product   string `json:"product"`
	Quantity  Frac   `json:"quantity"`
	CostCents *int64 `json:"cost_cents"`
	Kcal      *int64 `json:"kcal"`
}

// RecipeDetail es la salida de recipe_detail (§8.2, §9.3).
type RecipeDetail struct {
	Recipe              string             `json:"recipe"`
	ServingsMilli       int64              `json:"servings_milli"`
	Ingredients         []DetailIngredient `json:"ingredients"`
	TotalCostCents      *int64             `json:"total_cost_cents"`
	CostPerServingCents *int64             `json:"cost_per_serving_cents"`
	TotalKcal           *int64             `json:"total_kcal"`
	KcalPerServing      *int64             `json:"kcal_per_serving"`
	MissingPrices       []string           `json:"missing_prices"`
	MissingNutrition    []string           `json:"missing_nutrition"`
}

// RecipeDetailFor: total = HALF_UP(Σ exacto), un solo redondeo; por porción = HALF_UP(Σ ÷ porciones);
// sin precio o sin dato, el total queda ausente con la lista de productos (nunca 0).
func (e *Engine) RecipeDetailFor(code string, sm int64) (out *RecipeDetail, err error) {
	defer catch(&err)
	r, err := e.recipe(code)
	if err != nil {
		return nil, err
	}
	ev, err := e.evalRecipe(r, sm)
	if err != nil {
		return nil, err
	}
	out = &RecipeDetail{Recipe: code, ServingsMilli: sm, Ingredients: []DetailIngredient{},
		MissingPrices: ev.missingPrices, MissingNutrition: ev.missingNutrition}
	for _, ie := range ev.ings {
		out.Ingredients = append(out.Ingredients, DetailIngredient{ie.ing.Product, fracOf(ie.qty), roundOpt(ie.cost, 1), roundOpt(ie.kcal, 1000)})
	}
	out.TotalCostCents, out.TotalKcal = roundOpt(ev.cost, 1), roundOpt(ev.kcal, 1000)
	if ev.cost != nil {
		out.CostPerServingCents = roundOpt(mulR(ev.cost, r64(1000)), sm)
	}
	out.KcalPerServing = roundOpt(ev.kcal, sm)
	return out, nil
}
