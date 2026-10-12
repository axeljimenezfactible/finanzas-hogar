package precios

import (
	"fmt"
	"slices"
	"time"
)

// option es un par (tienda, presentación) con precio vigente que compite en una fecha.
type option struct {
	store       Store
	order       int // posición de la tienda en stores (la API la entrega por id ascendente)
	pres        Presentation
	price       PriceObs
	weight      Weight
	packs, cost int64
	substitute  bool
}

// adjCmp compara el precio por unidad base ajustado = (precio / contenido) ÷ peso como
// fracción exacta (§7.1, PROVISIONAL bloqueo S5-reglas-comparativa, pregunta 1, opción A).
// Solo se llama con pesos > 0 (candidatos que compiten).
func adjCmp(a, b option) int {
	return cmpRatBig(a.price.PriceCents, a.weight.Frac.Den, a.pres.SizeMilli, a.weight.Frac.Num,
		b.price.PriceCents, b.weight.Frac.Den, b.pres.SizeMilli, b.weight.Frac.Num)
}

// eligible: cadenas y sucursales compiten por todo; un mercado solo por los grupos que
// declara en store_product_group (RN-16, S5-11). Sam's compite como si tuvieran membresía.
func (b *Book) eligible(st Store, prod Product) bool {
	return st.Kind != KindMarket || slices.Contains(st.Groups, prod.Group)
}

// candidates devuelve los pares (tienda, presentación) del producto con precio vigente que
// compite en at, en el orden de stores y luego del catálogo (§14.1). store y pres vacíos =
// sin elección.
func (b *Book) candidates(prod Product, store, pres string, at time.Time) []option {
	var out []option
	for i, st := range b.Stores {
		if (store != "" && st.Code != store) || !b.eligible(st, prod) {
			continue
		}
		for _, p := range b.Catalog.Presentations {
			if p.Product != prod.Code || (pres != "" && p.Code != pres) {
				continue
			}
			cur := b.CurrentPrice(p.Code, st.Code, at)
			if cur.Price == nil || !cur.Weight.Competes {
				continue
			}
			out = append(out, option{store: st, order: i, pres: p, price: *cur.Price, weight: cur.Weight})
		}
	}
	return out
}

// quoteLess es el orden de §14.1: (1) menor precio por unidad base ajustado; (2) menor precio
// del paquete; (3) la tienda que va antes en stores; (4) la presentación que va antes en el
// catálogo. Los pasos 2 a 4 son PROVISIONALES (bloqueo S5-reglas-comparativa, pregunta 2).
func (b *Book) quoteLess(x, y option) bool {
	if c := adjCmp(x, y); c != 0 {
		return c < 0
	}
	if x.price.PriceCents != y.price.PriceCents {
		return x.price.PriceCents < y.price.PriceCents
	}
	if x.order != y.order {
		return x.order < y.order
	}
	return b.Catalog.PresOrder(x.pres.Code) < b.Catalog.PresOrder(y.pres.Code)
}

// bestQuote aplica la regla única de §14.1 a los candidatos.
func (b *Book) bestQuote(prod Product, store, pres string, at time.Time) (option, bool) {
	var best option
	found := false
	for _, o := range b.candidates(prod, store, pres, at) {
		if !found || b.quoteLess(o, best) {
			best, found = o, true
		}
	}
	return best, found
}

// Quote es la cotización que consume S4 (PriceSource; claves de quotes en S4-interfaz §12.1).
type Quote struct {
	Product, Store, Presentation string
	PriceCents, SizeMilli        int64
	ObservedDate                 time.Time
	IsStale, Recommended         bool
	Weight                       Weight
}

// QuoteFor cotiza un producto en la fecha at con la regla única de §14.1: sin tienda ni
// presentación, el par recomendado (Recommended = true); con tienda, con presentación o con
// las dos, el mismo criterio dentro de esa elección. ok = false es «sin precio». La API pasa
// at = hoy (America/Mexico_City); solo StalePricesForPlan usa la fecha de despensa (§10).
func (b *Book) QuoteFor(product, store, pres string, at time.Time) (Quote, bool, error) {
	prod, ok := b.Catalog.Product(product)
	if !ok {
		return Quote{}, false, fmt.Errorf("%w: %s", ErrUnknownProduct, product)
	}
	if store != "" {
		if _, ok := b.Store(store); !ok {
			return Quote{}, false, fmt.Errorf("%w: %s", ErrUnknownStore, store)
		}
	}
	if pres != "" {
		if p, ok := b.Catalog.Presentation(pres); !ok || p.Product != prod.Code {
			return Quote{}, false, fmt.Errorf("%w: %s", ErrUnknownPresentation, pres)
		}
	}
	o, found := b.bestQuote(prod, store, pres, at)
	if !found {
		return Quote{}, false, nil
	}
	return Quote{prod.Code, o.store.Code, o.pres.Code, o.price.PriceCents, o.pres.SizeMilli,
		o.price.ObservedDate, o.weight.IsStale, store == "" && pres == "", o.weight}, true, nil
}

// packChoice elige, en una tienda, la presentación que cubre buy con el menor costo de
// paquetes completos = techo(buy / contenido) × precio; empate: menor precio por unidad base,
// luego la que va antes en el catálogo (§7.3). Es la regla del total de cada tienda.
func (b *Book) packChoice(st Store, order int, prod Product, buy int64, at time.Time) (option, bool, error) {
	var best option
	found := false
	for _, o := range b.candidates(prod, st.Code, "", at) {
		o.order = order
		o.packs = (buy + o.pres.SizeMilli - 1) / o.pres.SizeMilli
		var err error
		if o.cost, err = mul(o.packs, o.price.PriceCents); err != nil {
			return option{}, false, err
		}
		if !found || o.cost < best.cost || (o.cost == best.cost && b.unitLess(o, best)) {
			best, found = o, true
		}
	}
	return best, found, nil
}

func (b *Book) unitLess(x, y option) bool {
	if c := cmpRat(x.price.PriceCents, x.pres.SizeMilli, y.price.PriceCents, y.pres.SizeMilli); c != 0 {
		return c < 0
	}
	return b.Catalog.PresOrder(x.pres.Code) < b.Catalog.PresOrder(y.pres.Code)
}

// ListItem es un renglón de la lista de compras de S4 (canasta por necesidad, §7.3, §15.1).
type ListItem struct {
	Product                 string
	NeededMilli             *int64 // nil: artículo manual
	QuantityMilli           *int64 // nil: sin cantidad fijada por la persona
	RecommendedPresentation string // "" = sin presentación que rija en la lista
}

// BuyMilli es la cantidad a comprar: quantity_milli si no es nula; si no, needed_milli
// (S4-interfaz §6.6). Las dos nulas o una cantidad ≤ 0 es argumento inválido.
func (it ListItem) BuyMilli() (int64, error) {
	q := it.QuantityMilli
	if q == nil {
		q = it.NeededMilli
	}
	if q == nil || *q <= 0 {
		return 0, fmt.Errorf("%w: cantidad a comprar de %s", ErrInvalidArgument, it.Product)
	}
	return *q, nil
}

// ItemPrice es la cotización de una tienda para un renglón.
type ItemPrice struct {
	Store              string `json:"store"`
	Presentation       string `json:"presentation"`
	PriceCents         int64  `json:"price_cents"`
	BaseUnitPriceCents int64  `json:"base_unit_price_cents"`
	ObservedDate       string `json:"observed_date"`
	IsSubstitute       bool   `json:"is_substitute"`
}

// BasketItem es la recomendación por renglón (C1 ShoppingListItem).
type BasketItem struct {
	Product                 string      `json:"product"`
	QuantityMilli           int64       `json:"quantity_milli"`
	RecommendedStore        *string     `json:"recommended_store"`
	RecommendedPresentation *string     `json:"recommended_presentation"`
	EstimatedPriceCents     *int64      `json:"estimated_price_cents"`
	IsStalePrice            bool        `json:"is_stale_price"`
	Prices                  []ItemPrice `json:"prices"`
}

// BasketComparison es la salida de CompareBasket (§7.3, §7.4; basket = shopping_list).
type BasketComparison struct {
	Stores        []StoreTotals `json:"stores"`
	Items         []BasketItem  `json:"items"`
	CheapestStore *string       `json:"cheapest_store"`
	opts          [][]option    // por renglón, la elección de cada tienda (para la membresía)
}

// CompareBasket arma la canasta por necesidad (§7.3): el total de cada tienda, también el de
// la recomendada de un renglón, es el menor costo de paquetes completos para la cantidad a
// comprar; la recomendación por renglón es QuoteFor sin elección con la misma fecha (§14.1).
func (b *Book) CompareBasket(list []ListItem) (*BasketComparison, error) {
	out := &BasketComparison{Stores: []StoreTotals{}, Items: []BasketItem{}}
	totals := make([]StoreTotals, len(b.Stores))
	for i, st := range b.Stores {
		totals[i] = StoreTotals{Store: st.Code, AssumesMembership: st.AssumesMembership(), order: i}
	}
	for _, it := range list {
		prod, ok := b.Catalog.Product(it.Product)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownProduct, it.Product)
		}
		buy, err := it.BuyMilli()
		if err != nil {
			return nil, err
		}
		if it.RecommendedPresentation != "" {
			if p, ok := b.Catalog.Presentation(it.RecommendedPresentation); !ok || p.Product != prod.Code {
				return nil, fmt.Errorf("%w: %s", ErrUnknownPresentation, it.RecommendedPresentation)
			}
		}
		item := BasketItem{Product: prod.Code, QuantityMilli: buy, Prices: []ItemPrice{}}
		var opts []option
		for i, st := range b.Stores {
			o, ok, err := b.packChoice(st, i, prod, buy, b.Today)
			if err != nil {
				return nil, err
			}
			t := &totals[i]
			if !ok {
				t.MissingCount++
				continue
			}
			// Sin presentación que rija no hay contra qué marcar sustituto (supuesto S5-go-1).
			o.substitute = it.RecommendedPresentation != "" && o.pres.Code != it.RecommendedPresentation
			opts = append(opts, o)
			t.TotalCents += o.cost
			t.PieceCount += o.packs
			if o.substitute {
				t.SubstituteCount++
			} else {
				t.IdenticalCount++
				t.IdenticalTotalCents += o.cost
			}
			if o.weight.IsStale {
				t.StaleCount++
			}
			bu, _, _, err := BaseUnitPrice(o.price.PriceCents, o.pres.SizeMilli, prod.BaseUnit)
			if err != nil {
				return nil, err
			}
			item.Prices = append(item.Prices, ItemPrice{st.Code, o.pres.Code, o.price.PriceCents, bu,
				o.price.ObservedDate.Format(time.DateOnly), o.substitute})
		}
		if r, ok := b.bestQuote(prod, "", "", b.Today); ok {
			est, err := mul((buy+r.pres.SizeMilli-1)/r.pres.SizeMilli, r.price.PriceCents)
			if err != nil {
				return nil, err
			}
			item.RecommendedStore, item.RecommendedPresentation = &r.store.Code, &r.pres.Code
			item.EstimatedPriceCents, item.IsStalePrice = &est, r.weight.IsStale
			totals[r.order].assigned++
		}
		out.Items = append(out.Items, item)
		out.opts = append(out.opts, opts)
	}
	out.Stores = rankStores(totals)
	for i := range out.Stores {
		if len(list) > 0 && out.Stores[i].MissingCount == 0 {
			out.CheapestStore = &out.Stores[i].Store // la primera completa (§7.4)
			break
		}
	}
	return out, nil
}

// CompareExpense (§7.2, API basket = expense; PROVISIONAL, bloqueo S5-carrito-comparativa,
// pregunta 1, opción A): cotiza cada otra tienda con sus precios vigentes que compiten para los
// productos y cantidades de la compra base, con la presentación de §7.3, y llama a CompareCarts.
// La fecha y la fuente de cada carrito cotizado son las de su precio más antiguo; una tienda sin
// ningún precio que compita para la compra no se cotiza (supuesto S5-go-3).
func (b *Book) CompareExpense(base Cart) (*CartComparison, error) {
	need := map[string]int64{}
	var prods []string
	for _, l := range base.Lines {
		p, ok := b.Catalog.Presentation(l.Presentation)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownPresentation, l.Presentation)
		}
		c, err := mul(l.QuantityMilli, p.SizeMilli)
		if err != nil {
			return nil, err
		}
		if _, seen := need[p.Product]; !seen {
			prods = append(prods, p.Product)
		}
		need[p.Product] += c / 1000
	}
	var others []Cart
	for i, st := range b.Stores {
		if st.Code == base.Store {
			continue
		}
		c := Cart{Code: st.Code, Store: st.Code, Kind: "cart", ObservedDate: b.Today, Source: SourceManual}
		for _, pc := range prods {
			prod, _ := b.Catalog.Product(pc)
			if !b.eligible(st, prod) {
				continue
			}
			o, ok, err := b.packChoice(st, i, prod, need[pc], b.Today)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			c.Lines = append(c.Lines, CartLine{o.pres.Code, o.packs * 1000, o.price.PriceCents})
			if o.price.ObservedDate.Before(c.ObservedDate) {
				c.ObservedDate, c.Source = o.price.ObservedDate, o.price.Source
			}
		}
		if len(c.Lines) > 0 {
			others = append(others, c)
		}
	}
	return b.CompareCarts(base, others)
}
