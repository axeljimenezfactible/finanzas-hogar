package precios

import (
	"fmt"
	"math/big"
	"sort"
	"time"
)

// CartLine es un renglón de canasta: presentaciones × 1000 y precio unitario en centavos.
type CartLine struct {
	Presentation   string
	QuantityMilli  int64
	UnitPriceCents int64
}

// Cart es una canasta de una tienda: compra registrada (expense) o carrito no comprado (cart).
type Cart struct {
	Code, Store, Kind string
	ObservedDate      time.Time
	Source            string
	Lines             []CartLine
}

// StoreTotals es el resumen por tienda de una comparativa (C1 BasketComparison.stores[]).
type StoreTotals struct {
	Store                     string `json:"store"`
	TotalCents                int64  `json:"total_cents"`
	IdenticalCount            int    `json:"identical_count"`
	IdenticalTotalCents       int64  `json:"identical_total_cents"`
	SubstituteCount           int    `json:"substitute_count"`
	PieceCount                int64  `json:"piece_count"`
	MissingCount              int    `json:"missing_count"`
	StaleCount                int    `json:"stale_count"`
	DifferenceVsCheapestCents *int64 `json:"difference_vs_cheapest_cents"` // null en incompletas (§7.4)
	AssumesMembership         bool   `json:"assumes_membership"`
	order, assigned           int
}

// RowLine es la parte de una tienda en un renglón de la comparativa.
type RowLine struct {
	Store              string `json:"store"`
	Presentation       string `json:"presentation"`
	QuantityMilli      int64  `json:"quantity_milli"`
	TotalCents         int64  `json:"total_cents"`
	BaseUnitPriceCents int64  `json:"base_unit_price_cents"`
	cart               int
}

// Row es un renglón: identical, substitute u only_one.
type Row struct {
	Kind        string    `json:"kind"`
	Product     string    `json:"product"`
	Lines       []RowLine `json:"lines"`
	WinnerStore *string   `json:"winner_store"`
}

// StoreCount cuenta renglones idénticos ganados por tienda.
type StoreCount struct {
	Store string `json:"store"`
	Count int    `json:"count"`
}

// CartComparison es la salida de CompareCarts (§7.2, S5-01).
type CartComparison struct {
	Stores          []StoreTotals `json:"stores"`
	IdenticalPieces int64         `json:"identical_pieces"`
	Winners         []StoreCount  `json:"winners"`
	Ties            int           `json:"ties"`
	Rows            []Row         `json:"rows"`
}

// agg agrupa los renglones de una presentación en una canasta y los consume en orden de
// captura: la parte emparejada sale de los primeros renglones y el excedente de los últimos
// (la sexta lata de atún de Bodega es el último renglón, a 21.00). PROVISIONAL (bloqueo
// S5-motor-oraculo-reglas pregunta 2; §7.2 pasos 1 y 2 no dicen qué precio toma la parte emparejada si la
// misma presentación trae precios distintos): en uso, el orden de captura; el oráculo usa el promedio.
type agg struct {
	pres  Presentation
	qty   int64 // milésimas de presentación sin consumir
	lines []CartLine
}

// take consume qty milésimas y devuelve su importe exacto en milésimas de centavo.
func (a *agg) take(qty int64) *big.Int {
	sum := new(big.Int)
	a.qty -= qty
	for qty > 0 && len(a.lines) > 0 {
		l := &a.lines[0]
		q := min(qty, l.QuantityMilli)
		sum.Add(sum, new(big.Int).Mul(big.NewInt(q), big.NewInt(l.UnitPriceCents)))
		l.QuantityMilli -= q
		qty -= q
		if l.QuantityMilli == 0 {
			a.lines = a.lines[1:]
		}
	}
	return sum
}

func pieces(qty int64) int64 {
	if qty%1000 == 0 {
		return qty / 1000
	}
	// Renglón a granel = 1 pieza (supuesto S-11). PROVISIONAL (bloqueo S5-motor-oraculo-reglas pregunta 1;
	// §7.2 paso 6 no dice qué es «renglón»): en uso, la presentación agrupada del paso 1; el
	// oráculo cuenta cada línea de rows[] (emparejada y excedente por separado).
	return 1
}

// CompareCarts compara la canasta base contra otras, renglón contra renglón (§7.2): idénticos
// con la cantidad menor, sustitutos con su cantidad e importe y ganador por unidad base exacta,
// excedentes como «solo en una tienda». winners y ties cuentan solo idénticos (supuesto S-9).
func (b *Book) CompareCarts(base Cart, others []Cart) (*CartComparison, error) {
	if len(others) == 0 {
		return nil, fmt.Errorf("%w: sin canastas que comparar", ErrInvalidArgument)
	}
	carts := append([]Cart{base}, others...)
	aggs := make([]map[string]*agg, len(carts))
	byProd := make([]map[string][]string, len(carts))
	presOrder := make([][]string, len(carts)) // presentaciones por canasta, en orden de aparición
	totals := make([]StoreTotals, len(carts))
	exactTot := make([]*big.Int, len(carts))
	identTot := make([]*big.Int, len(carts))
	stale := make([]bool, len(carts))
	var prodOrder []string
	seenProd := map[string]bool{}
	for ci, c := range carts {
		if _, ok := Competes(c.Source); !ok {
			return nil, fmt.Errorf("%w: fuente %q", ErrInvalidArgument, c.Source)
		}
		st, ok := b.Store(c.Store)
		if !ok && len(b.Stores) > 0 {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStore, c.Store)
		}
		age := AgeDays(c.ObservedDate, b.Today)
		if age < 0 {
			return nil, ErrFutureDate
		}
		stale[ci] = age >= b.Settings.StaleAfterDays
		totals[ci] = StoreTotals{Store: c.Store, AssumesMembership: st.AssumesMembership(), order: ci}
		exactTot[ci], identTot[ci] = new(big.Int), new(big.Int)
		aggs[ci], byProd[ci] = map[string]*agg{}, map[string][]string{}
		for _, l := range c.Lines {
			p, ok := b.Catalog.Presentation(l.Presentation)
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownPresentation, l.Presentation)
			}
			if l.QuantityMilli <= 0 || l.UnitPriceCents < 0 {
				return nil, fmt.Errorf("%w: renglón de %s", ErrInvalidArgument, l.Presentation)
			}
			a := aggs[ci][p.Code]
			if a == nil {
				a = &agg{pres: p}
				aggs[ci][p.Code] = a
				byProd[ci][p.Product] = append(byProd[ci][p.Product], p.Code)
				presOrder[ci] = append(presOrder[ci], p.Code)
				if !seenProd[p.Product] {
					seenProd[p.Product] = true
					prodOrder = append(prodOrder, p.Product)
				}
			}
			a.qty += l.QuantityMilli
			a.lines = append(a.lines, l)
			// Total exacto de la canasta, también a granel, y un solo HALF_UP al publicar: PROVISIONAL, bloqueo
			// S5-carritos-granel, opción A en uso (§7.2 paso 8).
			exactTot[ci].Add(exactTot[ci], new(big.Int).Mul(big.NewInt(l.QuantityMilli), big.NewInt(l.UnitPriceCents)))
		}
		for _, a := range aggs[ci] {
			totals[ci].PieceCount += pieces(a.qty)
		}
	}
	out := &CartComparison{Winners: []StoreCount{}, Rows: []Row{}}
	wins := make([]int, len(carts))
	inAll := func(pc string) bool {
		for ci := range carts {
			if _, ok := aggs[ci][pc]; !ok {
				return false
			}
		}
		return true
	}
	// Orden de rows[] (convenciones 2 y 3 del oráculo S5; solicitud a §15.2): primero los
	// idénticos en el orden de la canasta base, luego los sustitutos por producto en orden de
	// aparición y al final lo que queda «solo en una tienda», por canasta y renglón.
	var identRows, subRows, onlyRows []Row
	// Idénticos: la misma presentación en todas las canastas, con la cantidad menor.
	for _, pc := range presOrder[0] {
		if !inAll(pc) {
			continue
		}
		m := aggs[0][pc].qty
		for ci := 1; ci < len(carts); ci++ {
			m = min(m, aggs[ci][pc].qty)
		}
		row := Row{Kind: "identical", Product: aggs[0][pc].pres.Product}
		var bestAmt *big.Int
		best, tie := -1, false
		for ci := range carts {
			a := aggs[ci][pc]
			amt := a.take(m)
			if err := b.addLine(&row, ci, carts[ci].Store, a.pres, m, amt); err != nil {
				return nil, err
			}
			identTot[ci].Add(identTot[ci], amt)
			// Misma cantidad en todas: comparar importes es comparar precio unitario exacto.
			switch {
			case best < 0 || amt.Cmp(bestAmt) < 0:
				best, bestAmt, tie = ci, amt, false
			case amt.Cmp(bestAmt) == 0:
				tie = true
			}
		}
		if tie {
			out.Ties++
		} else {
			wins[best]++
			row.WinnerStore = &carts[best].Store
		}
		out.IdenticalPieces += pieces(m)
		identRows = append(identRows, row)
	}
	// Sustitutos: un producto con presentaciones sin emparejar (que no están en todas las
	// canastas) en dos o más canastas; cada una conserva su cantidad y su importe.
	for _, prod := range prodOrder {
		var with []int
		for ci := range carts {
			for _, pc := range byProd[ci][prod] {
				if !inAll(pc) && aggs[ci][pc].qty > 0 {
					with = append(with, ci)
					break
				}
			}
		}
		if len(with) < 2 {
			continue
		}
		row := Row{Kind: "substitute", Product: prod}
		var bestR *big.Rat
		best, tie := -1, false
		for _, ci := range with {
			cost, content := new(big.Int), new(big.Int)
			for _, pc := range byProd[ci][prod] {
				a := aggs[ci][pc]
				q := a.qty
				if q == 0 || inAll(pc) {
					continue
				}
				amt := a.take(q)
				if err := b.addLine(&row, ci, carts[ci].Store, a.pres, q, amt); err != nil {
					return nil, err
				}
				cost.Add(cost, amt)
				content.Add(content, new(big.Int).Mul(big.NewInt(q), big.NewInt(a.pres.SizeMilli)))
			}
			// Precio por unidad base exacto (RN-20) de lo que esta tienda cobró por el producto:
			// importe total ÷ contenido total de sus presentaciones sin emparejar (PROVISIONAL,
			// bloqueo S5-carritos-detalle, opción A).
			r := new(big.Rat).SetFrac(cost, content)
			switch {
			case best < 0 || r.Cmp(bestR) < 0:
				best, bestR, tie = ci, r, false
			case r.Cmp(bestR) == 0:
				tie = true
			}
		}
		if !tie {
			row.WinnerStore = &carts[best].Store
		}
		subRows = append(subRows, row)
	}
	// Lo que queda solo está en una tienda: por canasta y en el orden de sus renglones.
	for ci := range carts {
		for _, pc := range presOrder[ci] {
			a := aggs[ci][pc]
			if q := a.qty; q > 0 {
				row := Row{Kind: "only_one", Product: a.pres.Product}
				if err := b.addLine(&row, ci, carts[ci].Store, a.pres, q, a.take(q)); err != nil {
					return nil, err
				}
				onlyRows = append(onlyRows, row)
			}
		}
	}
	// Faltantes y viejos de cada canasta (PROVISIONAL, bloqueo S5-carritos-detalle, opción A):
	// missing_count = productos de la unión de las canastas sin ninguna presentación en esta;
	// stale_count = renglones de la canasta agrupados por presentación (§7.2 paso 1) si su fecha
	// es vieja, cada presentación una sola vez aunque salga en un idéntico y en un «solo aquí».
	for ci := range carts {
		for _, prod := range prodOrder {
			if len(byProd[ci][prod]) == 0 {
				totals[ci].MissingCount++
			}
		}
		if stale[ci] {
			totals[ci].StaleCount = len(aggs[ci])
		}
	}
	out.Rows = append(append(append(out.Rows, identRows...), subRows...), onlyRows...)
	for _, r := range out.Rows {
		in := make([]bool, len(carts))
		for _, l := range r.Lines {
			in[l.cart] = true
		}
		for ci := range carts {
			switch {
			case r.Kind == "identical":
				totals[ci].IdenticalCount++
			case r.Kind == "substitute" && in[ci]:
				totals[ci].SubstituteCount++
			}
		}
	}
	for ci := range carts {
		var err error
		if totals[ci].TotalCents, err = halfUp(exactTot[ci], big.NewInt(1000)); err != nil {
			return nil, err
		}
		if totals[ci].IdenticalTotalCents, err = halfUp(identTot[ci], big.NewInt(1000)); err != nil {
			return nil, err
		}
	}
	out.Stores = rankCarts(totals)
	for _, st := range out.Stores {
		out.Winners = append(out.Winners, StoreCount{st.Store, wins[st.order]})
	}
	return out, nil
}

// addLine agrega la parte de una tienda: importe y precio por unidad base con HALF_UP.
func (b *Book) addLine(row *Row, ci int, store string, p Presentation, qty int64, amt *big.Int) error {
	tot, err := halfUp(amt, big.NewInt(1000))
	if err != nil {
		return err
	}
	factor := int64(1000000) // centavos por 1 kg o 1 l
	if b.Catalog.unitOf(p) == "piece" {
		factor = 1000
	}
	// amt en milésimas de centavo; contenido = qty × size / 1000 milésimas de unidad base.
	num := new(big.Int).Mul(amt, big.NewInt(factor))
	den := new(big.Int).Mul(big.NewInt(qty), big.NewInt(p.SizeMilli))
	bu, err := halfUp(num, den)
	if err != nil {
		return err
	}
	row.Lines = append(row.Lines, RowLine{store, p.Code, qty, tot, bu, ci})
	return nil
}

// rankCarts ordena las canastas de CompareCarts (§7.4, §15.2): completas por total ascendente y
// luego incompletas por total ascendente; a igual total, la base y luego others en su orden.
func rankCarts(t []StoreTotals) []StoreTotals {
	return rankComplete(t, func(a, b *StoreTotals) bool { return a.order < b.order })
}

// rankStores ordena las tiendas de la canasta por necesidad (§7.4, §15.2): completas por total
// ascendente y luego incompletas por total ascendente. Empate de totales (FIRME, puerta 2, Q2-60 = A):
// más renglones recomendados, luego menos precios viejos, luego el orden de stores.
func rankStores(t []StoreTotals) []StoreTotals {
	return rankComplete(t, func(a, b *StoreTotals) bool {
		if a.assigned != b.assigned {
			return a.assigned > b.assigned
		}
		if a.StaleCount != b.StaleCount {
			return a.StaleCount < b.StaleCount
		}
		return a.order < b.order
	})
}

// rankComplete aplica la regla única de §7.4 en los dos modos (H-r3-04, Q2-05 = A; M-r5-3,
// Q2-08 = A): una tienda está completa si su missing_count es 0; van primero las completas y
// después las incompletas, cada grupo por total ascendente y, a igual total, por tie. La
// diferencia es contra la primera completa y queda en null en las incompletas o si ninguna
// tienda está completa.
func rankComplete(t []StoreTotals, tie func(a, b *StoreTotals) bool) []StoreTotals {
	s := append([]StoreTotals(nil), t...)
	sort.SliceStable(s, func(i, j int) bool {
		if ci, cj := s[i].MissingCount == 0, s[j].MissingCount == 0; ci != cj {
			return ci
		}
		if s[i].TotalCents != s[j].TotalCents {
			return s[i].TotalCents < s[j].TotalCents
		}
		return tie(&s[i], &s[j])
	})
	for i := range s {
		s[i].DifferenceVsCheapestCents = nil
		if s[i].MissingCount == 0 {
			d := s[i].TotalCents - s[0].TotalCents // s[0] es completa si s[i] lo es
			s[i].DifferenceVsCheapestCents = &d
		}
	}
	return s
}
