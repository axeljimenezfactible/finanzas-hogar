package recetas

import (
	"math/big"
	"sort"
)

// Lista de compras neteada con inventario híbrido (§6; RN-15, RN-16, K1; S4-03, S4-06, S4-12).

// Tipos de artículo de la lista (shopping_list_item.origin).
const (
	KindCalculated = "calculated"
	KindManual     = "manual"
)

// ShopItem es un artículo de la lista (C1 ShoppingListItem; traducción en §12.2).
type ShopItem struct {
	Product                 string   `json:"product"`
	Kind                    string   `json:"kind"`
	NeededMilli             *int64   `json:"needed_milli"`
	QuantityMilli           *int64   `json:"quantity_milli"`
	Presentation            *string  `json:"presentation"`
	PackagesToBuy           *int64   `json:"packages_to_buy"`
	LeftoverMilli           *int64   `json:"leftover_milli"`
	Reasons                 []string `json:"reasons"`
	Recipes                 []string `json:"recipes"`
	Store                   *string  `json:"store"`
	IsBought                bool     `json:"is_bought"`
	RecommendedStore        *string  `json:"recommended_store"`
	RecommendedPresentation *string  `json:"recommended_presentation"`
	PriceStore              *string  `json:"price_store"`
	PriceObservedDate       *Date    `json:"price_observed_date"`
	IsStalePrice            *bool    `json:"is_stale_price"`
	EstimatedCents          *int64   `json:"estimated_cents"`
	vertical                Channel
}

// ShopGroup es un grupo (vertical, tienda) de la lista.
type ShopGroup struct {
	Vertical            Channel `json:"vertical"`
	Store               *string `json:"store"`
	VisitDate           *Date   `json:"visit_date"`
	EstimatedTotalCents int64   `json:"estimated_total_cents"`
}

// ShoppingList es la salida de shopping_list (§6, §12.2).
type ShoppingList struct {
	Week                string      `json:"week"`
	Items               []ShopItem  `json:"items"`
	Groups              []ShopGroup `json:"groups"`
	EstimatedTotalCents int64       `json:"estimated_total_cents"`
	MissingPrices       []string    `json:"missing_prices"`
}

type need struct {
	recipe   *big.Rat
	repl     int64
	stock    int64
	recipes  []string
	pres     string
	fromPlan bool
}

// ShoppingListFor calcula la lista de la semana (lunes): necesidad neta = max(0, N + R − E) hacia arriba
// a la milésima; paquetes de la presentación que rige con la cantidad a comprar; sobrante sin negativos.
func (e *Engine) ShoppingListFor(monday Date) (out *ShoppingList, err error) {
	defer catch(&err)
	if monday.ISOWeekday() != 1 {
		return nil, errf(ErrNotMonday, "%s", monday)
	}
	entries, err := e.weekEntries(monday)
	if err != nil {
		return nil, err
	}
	covered, _ := e.coverLeftovers(entries)
	needs := map[string]*need{}
	get := func(p string) *need {
		if needs[p] == nil {
			needs[p] = &need{recipe: r64(0)}
		}
		return needs[p]
	}
	for i, en := range entries {
		if en.EatingOut || *en.ServingsMilli-covered[i] <= 0 {
			continue
		}
		r, _ := e.recipe(en.Recipe)
		for _, in := range r.Ingredients {
			if in.Optional {
				continue
			}
			q, err := Scale(in.QuantityMilli, r.BaseServings, *en.ServingsMilli-covered[i])
			if err != nil {
				return nil, err
			}
			n := get(in.Product)
			n.recipe, n.fromPlan = addR(n.recipe, q), true
			if n.pres == "" {
				n.pres = in.Presentation
			}
			if len(n.recipes) == 0 || !contains(n.recipes, r.Code) {
				n.recipes = append(n.recipes, r.Code)
			}
		}
	}
	for _, it := range e.Inventory {
		if it.Product == "" {
			continue
		}
		n := get(it.Product)
		n.stock += estimateAt(it, monday)
		if it.Mode == ModeRate && it.ReorderLevelMilli != nil && estimateAt(it, monday.AddDays(6)) <= *it.ReorderLevelMilli {
			n.repl += weeklyConsumption(it) // FIRME (Q2-45 = A; bloqueo S4-reglas-lista, pregunta 2, RESUELTO)
		}
	}
	for _, ch := range e.Choices {
		if ch.Kind == KindCalculated {
			get(ch.Product)
		}
	}
	codes := make([]string, 0, len(needs))
	for p := range needs {
		codes = append(codes, p)
	}
	sort.Strings(codes)
	out = &ShoppingList{Week: WeekCode(monday), Items: []ShopItem{}, Groups: []ShopGroup{}, MissingPrices: []string{}}
	for _, p := range codes {
		n := needs[p]
		net := ceilRat(subR(addR(n.recipe, r64(n.repl)), r64(n.stock)))
		ch := e.calcChoice(p)
		if net <= 0 && ch == nil {
			continue
		}
		net = max(net, 0)
		pres := n.pres
		if pres == "" {
			pres = e.Catalog.fallbackPresentation(p)
		}
		item, err := e.buildItem(p, KindCalculated, &net, ch, pres)
		if err != nil {
			return nil, err
		}
		if n.recipe.Sign() > 0 {
			item.Reasons = append(item.Reasons, "recipe")
		}
		if n.repl > 0 {
			item.Reasons = append(item.Reasons, "replenishment")
		}
		if n.recipes != nil {
			sort.Strings(n.recipes) // orden por código (supuesto S-9: el contrato no fija el orden)
			item.Recipes = n.recipes
		}
		out.Items = append(out.Items, *item)
	}
	for i := range e.Choices {
		ch := &e.Choices[i]
		if ch.Kind != KindManual {
			continue
		}
		item, err := e.buildItem(ch.Product, KindManual, nil, ch, "")
		if err != nil {
			return nil, err
		}
		item.Reasons = append(item.Reasons, "manual")
		out.Items = append(out.Items, *item)
	}
	e.groupItems(out)
	return out, nil
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// buildItem fija presentación, paquetes, sobrante, tienda, precio y recomendación de un artículo (§6.2
// pasos 5 y 6, §6.5, §6.6). needed es nil en un artículo manual.
func (e *Engine) buildItem(p, kind string, needed *int64, ch *ListChoice, fallback string) (*ShopItem, error) {
	vertical, _, err := e.Catalog.channelOf(p)
	if err != nil {
		return nil, err
	}
	quote, err := e.governing(p, ch)
	if err != nil {
		return nil, err
	}
	rec, err := e.recommended(p)
	if err != nil {
		return nil, err
	}
	it := &ShopItem{Product: p, Kind: kind, NeededMilli: needed, Reasons: []string{}, Recipes: []string{}, vertical: vertical}
	chStore, chPres := "", ""
	if ch != nil {
		it.QuantityMilli, it.IsBought, chStore, chPres = ch.QuantityMilli, ch.IsBought, ch.Store, ch.Presentation
	}
	var size int64
	switch {
	case chPres != "":
		x, err := e.Catalog.presentationOf(p, chPres)
		if err != nil {
			return nil, err
		}
		it.Presentation, size = &chPres, x.SizeMilli
	case quote != nil:
		it.Presentation = ptr(quote.Presentation)
	case fallback != "":
		it.Presentation, size = &fallback, e.Catalog.Presentations[fallback].SizeMilli
	}
	if quote != nil {
		size = quote.SizeMilli
	}
	if it.Presentation != nil {
		qty := int64(0)
		if needed != nil {
			qty = *needed
		}
		if it.QuantityMilli != nil {
			qty = *it.QuantityMilli
		}
		pk := ceilRat(frac(qty, size))
		nm := int64(0)
		if needed != nil {
			nm = *needed
		}
		it.PackagesToBuy, it.LeftoverMilli = &pk, ptr(max(0, pk*size-nm))
		if quote != nil {
			it.PriceStore, it.PriceObservedDate = ptr(quote.Store), ptr(quote.ObservedDate)
			it.IsStalePrice, it.EstimatedCents = ptr(quote.IsStale), ptr(pk*quote.PriceCents)
		}
	}
	switch {
	case chStore != "":
		it.Store = &chStore
	case vertical == ChannelChain:
		if quote != nil {
			it.Store = ptr(quote.Store)
		}
	default:
		it.Store = e.marketFor(p) // FIRME (Q2-39 y Q2-40 = A; bloqueo S4-asignacion-mercado, RESUELTO)
	}
	if chStore == "" {
		it.RecommendedStore = it.PriceStore
	} else if rec != nil {
		it.RecommendedStore = ptr(rec.Store)
	}
	if chPres == "" {
		it.RecommendedPresentation = it.Presentation
	} else if rec != nil {
		it.RecommendedPresentation = ptr(rec.Presentation)
	}
	return it, nil
}

// marketFor: entre los mercados que surten el grupo del producto, el de next_visit más cercana; en empate,
// el primero en el orden de stores; sin ninguno, nil (§6.5 paso 3; FIRME, Q2-39 y Q2-40 = A).
func (e *Engine) marketFor(p string) *string {
	g := e.Catalog.Products[p].Group
	var best *Store
	for i := range e.Stores {
		s := &e.Stores[i]
		if s.Kind != "market" || !contains(s.Groups, g) {
			continue
		}
		if best == nil || s.NextVisit.Before(*best.NextVisit) {
			best = s
		}
	}
	if best == nil {
		return nil
	}
	return ptr(best.Code)
}

// groupItems arma groups[] (market primero; visit_date ascendente con nulo al final; orden de stores con
// tienda nula al final), sus totales con los artículos con precio, el total de la raíz y missing_prices[].
func (e *Engine) groupItems(out *ShoppingList) {
	sidx := map[string]int{}
	for i, s := range e.Stores {
		sidx[s.Code] = i
	}
	pos := map[string]int{}
	for _, it := range out.Items {
		k := string(it.vertical) + "|"
		if it.Store != nil {
			k += *it.Store
		}
		i, ok := pos[k]
		if !ok {
			g := ShopGroup{Vertical: it.vertical, Store: it.Store}
			if it.Store != nil {
				if j, ok := sidx[*it.Store]; ok && e.Stores[j].Kind == "market" {
					g.VisitDate = e.Stores[j].NextVisit
				}
			}
			i = len(out.Groups)
			pos[k] = i
			out.Groups = append(out.Groups, g)
		}
		if it.EstimatedCents != nil {
			out.Groups[i].EstimatedTotalCents += *it.EstimatedCents
			out.EstimatedTotalCents += *it.EstimatedCents
		} else if !contains(out.MissingPrices, it.Product) {
			out.MissingPrices = append(out.MissingPrices, it.Product)
		}
	}
	rank := func(g ShopGroup) (int, int) {
		if g.Store == nil {
			return 2, 0
		}
		if j, ok := sidx[*g.Store]; ok {
			return 0, j
		}
		return 1, 0
	}
	sort.SliceStable(out.Groups, func(i, j int) bool {
		a, b := out.Groups[i], out.Groups[j]
		if a.Vertical != b.Vertical {
			return a.Vertical == ChannelMarket
		}
		if (a.VisitDate == nil) != (b.VisitDate == nil) {
			return a.VisitDate != nil
		}
		if a.VisitDate != nil && *a.VisitDate != *b.VisitDate {
			return a.VisitDate.Before(*b.VisitDate)
		}
		ra, ia := rank(a)
		rb, ib := rank(b)
		if ra != rb {
			return ra < rb
		}
		return ia < ib
	})
}
