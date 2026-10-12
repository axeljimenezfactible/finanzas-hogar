package precios

import (
	"fmt"
	"math/big"
	"sort"
	"time"
)

// Catalog es el catálogo de productos y presentaciones (product, presentation, product_group).
type Catalog struct {
	Groups        []Group
	Products      []Product
	Presentations []Presentation
	prod, pres    map[string]int
}

// NewCatalog indexa el catálogo y valida códigos únicos, producto existente y contenido > 0.
func NewCatalog(groups []Group, products []Product, pres []Presentation) (*Catalog, error) {
	c := &Catalog{Groups: groups, Products: products, Presentations: pres,
		prod: map[string]int{}, pres: map[string]int{}}
	for i, p := range products {
		if _, dup := c.prod[p.Code]; dup || p.Code == "" {
			return nil, fmt.Errorf("%w: producto %q repetido o vacío", ErrInvalidArgument, p.Code)
		}
		if p.BaseUnit != "g" && p.BaseUnit != "ml" && p.BaseUnit != "piece" {
			return nil, fmt.Errorf("%w: unidad base %q", ErrInvalidArgument, p.BaseUnit)
		}
		c.prod[p.Code] = i
	}
	for i, p := range pres {
		if _, dup := c.pres[p.Code]; dup || p.Code == "" {
			return nil, fmt.Errorf("%w: presentación %q repetida o vacía", ErrInvalidArgument, p.Code)
		}
		if _, ok := c.prod[p.Product]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownProduct, p.Product)
		}
		if p.SizeMilli <= 0 {
			return nil, fmt.Errorf("%w: contenido de %s", ErrInvalidArgument, p.Code)
		}
		c.pres[p.Code] = i
	}
	return c, nil
}

// Product devuelve el producto por código.
func (c *Catalog) Product(code string) (Product, bool) {
	i, ok := c.prod[code]
	if !ok {
		return Product{}, false
	}
	return c.Products[i], true
}

// Presentation devuelve la presentación por código.
func (c *Catalog) Presentation(code string) (Presentation, bool) {
	i, ok := c.pres[code]
	if !ok {
		return Presentation{}, false
	}
	return c.Presentations[i], true
}

// PresOrder es la posición de la presentación en el catálogo (sustituye al «menor id»).
func (c *Catalog) PresOrder(code string) int { return c.pres[code] }

// unitOf devuelve la unidad base del producto de una presentación.
func (c *Catalog) unitOf(p Presentation) string {
	pr, _ := c.Product(p.Product)
	return pr.BaseUnit
}

// BaseUnitPrice (RN-20, K7) devuelve el precio por 1 kg, 1 l o 1 pieza con HALF_UP, la
// unidad publicada y la fracción exacta en centavos por esa unidad (sin redondear).
func BaseUnitPrice(priceCents, sizeMilli int64, baseUnit string) (cents int64, per string, exact Frac, err error) {
	if priceCents < 0 || sizeMilli <= 0 {
		return 0, "", Frac{}, ErrInvalidArgument
	}
	factor, per := int64(1000000), "kg"
	switch baseUnit {
	case "ml":
		per = "l"
	case "piece":
		factor, per = 1000, "piece"
	case "g":
	default:
		return 0, "", Frac{}, ErrInvalidArgument
	}
	n, err := mul(priceCents, factor)
	if err != nil {
		return 0, "", Frac{}, err
	}
	cents, err = halfUp(big.NewInt(n), big.NewInt(sizeMilli))
	return cents, per, NewFrac(n, sizeMilli), err
}

// Weight es el peso de frescura de un precio (RN-12, K22, K40, K44).
type Weight struct {
	Frac     Frac
	Competes bool
	IsStale  bool
	AgeDays  int
}

// Display muestra el peso con 4 decimales HALF_UP (solo para pantalla).
func (w Weight) Display() string {
	v, _ := halfUp(big.NewInt(w.Frac.Num*10000), big.NewInt(w.Frac.Den))
	return fmt.Sprintf("%d.%04d", v/10000, v%10000)
}

// Freshness calcula edad, «viejo» y peso: 1 si edad < Dv; (D0 − edad)/(D0 − Dv + 1) si
// Dv ≤ edad ≤ D0; 0 y sin competir si edad > D0. Compite si la fuente compite y peso > 0.
func Freshness(observed, today time.Time, s Settings, source string) (Weight, error) {
	comp, ok := Competes(source)
	if !ok {
		return Weight{}, fmt.Errorf("%w: fuente %q", ErrInvalidArgument, source)
	}
	if err := s.Validate(); err != nil {
		return Weight{}, err
	}
	age := AgeDays(observed, today)
	if age < 0 {
		return Weight{}, ErrFutureDate
	}
	w := Weight{AgeDays: age, IsStale: age >= s.StaleAfterDays}
	dv, d0 := int64(s.StaleAfterDays), int64(s.ZeroWeightAfterDays)
	switch {
	case age < s.StaleAfterDays:
		w.Frac = Frac{1, 1}
	case age <= s.ZeroWeightAfterDays:
		w.Frac = NewFrac(d0-int64(age), d0-dv+1)
	default:
		w.Frac = Frac{0, 1}
	}
	w.Competes = comp && w.Frac.Num > 0
	return w, nil
}

// Book reúne catálogo, tiendas, precios, parámetros y «hoy» (America/Mexico_City, C-05).
type Book struct {
	Catalog  *Catalog
	Stores   []Store
	Prices   []PriceObs
	Settings Settings
	Today    time.Time
	store    map[string]int
}

// NewBook valida tiendas y precios (S5-interfaz §3, validación común) y arma el libro.
func NewBook(c *Catalog, stores []Store, prices []PriceObs, s Settings, today time.Time) (*Book, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	b := &Book{Catalog: c, Settings: s, Today: today, store: map[string]int{}}
	for _, st := range stores {
		if _, dup := b.store[st.Code]; dup || st.Code == "" {
			return nil, fmt.Errorf("%w: tienda %q repetida o vacía", ErrInvalidArgument, st.Code)
		}
		if st.Kind != KindChain && st.Kind != KindBranch && st.Kind != KindMarket {
			return nil, fmt.Errorf("%w: tipo de tienda %q", ErrInvalidArgument, st.Kind)
		}
		if st.VisitWeekday < 0 || st.VisitWeekday > 7 || (st.Kind == KindMarket && st.VisitWeekday == 0) {
			return nil, fmt.Errorf("%w: día de visita de %s", ErrInvalidArgument, st.Code)
		}
		b.store[st.Code] = len(b.Stores)
		b.Stores = append(b.Stores, st)
	}
	for _, p := range prices {
		if err := b.AddPrice(p); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Store devuelve la tienda por código.
func (b *Book) Store(code string) (Store, bool) {
	i, ok := b.store[code]
	if !ok {
		return Store{}, false
	}
	return b.Stores[i], true
}

// AddPrice ingiere un precio de cualquier fuente con la validación común de §3: centavos ≥ 0,
// fecha no futura, presentación y tienda existentes y llave única
// (presentación, tienda, fecha, fuente, promoción).
func (b *Book) AddPrice(p PriceObs) error {
	if p.PriceCents < 0 {
		return fmt.Errorf("%w: precio negativo", ErrInvalidArgument)
	}
	if _, ok := Competes(p.Source); !ok {
		return fmt.Errorf("%w: fuente %q", ErrInvalidArgument, p.Source)
	}
	if _, ok := b.Catalog.Presentation(p.Presentation); !ok {
		return fmt.Errorf("%w: %s", ErrUnknownPresentation, p.Presentation)
	}
	if _, ok := b.Store(p.Store); !ok {
		return fmt.Errorf("%w: %s", ErrUnknownStore, p.Store)
	}
	if AgeDays(p.ObservedDate, b.Today) < 0 {
		return ErrFutureDate
	}
	for _, q := range b.Prices {
		if q.Presentation == p.Presentation && q.Store == p.Store && q.ObservedDate.Equal(p.ObservedDate) &&
			q.Source == p.Source && q.IsPromotional == p.IsPromotional {
			return ErrDuplicatePrice
		}
	}
	b.Prices = append(b.Prices, p)
	return nil
}

// later dice si a es posterior a b (fecha observada y, a igual fecha, created_at; supuesto S-8). Con el mismo created_at
// no es posterior, así que queda el primero de la entrada (PROVISIONAL, bloqueo S5-precio-vigente-captura pregunta 1).
func later(a, b PriceObs) bool {
	if !a.ObservedDate.Equal(b.ObservedDate) {
		return a.ObservedDate.After(b.ObservedDate)
	}
	return a.CreatedAt.After(b.CreatedAt)
}

// Current es el precio vigente de una presentación en una tienda (§6).
type Current struct {
	Price      *PriceObs // último precio regular de fuente que compite (o nil: «sin precio»)
	Weight     Weight
	References []PriceObs // último de cada fuente de referencia (profeco, site_assisted); PROVISIONAL, bloqueo S5-motor-oraculo-reglas pregunta 3 (§6)
	Promotion  *PriceObs  // promoción vigente: solo en su observed_date (FIRME, puerta 2, Q2-61 = A)
}

// CurrentPrice evalúa el precio vigente en la fecha at (hoy o la fecha de despensa).
func (b *Book) CurrentPrice(pres, store string, at time.Time) Current {
	var cur Current
	refs := map[string]PriceObs{}
	for i := range b.Prices {
		p := b.Prices[i]
		if p.Presentation != pres || p.Store != store || p.ObservedDate.After(at) {
			continue
		}
		comp, _ := Competes(p.Source)
		switch {
		case p.IsPromotional:
			if comp && AgeDays(p.ObservedDate, at) == 0 && (cur.Promotion == nil || later(p, *cur.Promotion)) {
				q := p
				cur.Promotion = &q
			}
		case comp:
			if cur.Price == nil || later(p, *cur.Price) {
				q := p
				cur.Price = &q
			}
		default:
			if r, ok := refs[p.Source]; !ok || later(p, r) {
				refs[p.Source] = p
			}
		}
	}
	if cur.Price != nil {
		cur.Weight, _ = Freshness(cur.Price.ObservedDate, at, b.Settings, cur.Price.Source)
	}
	for _, r := range refs {
		cur.References = append(cur.References, r)
	}
	// §15.2: por observed_date descendente y luego por source (orden de bytes).
	sort.Slice(cur.References, func(i, j int) bool {
		ri, rj := cur.References[i], cur.References[j]
		if !ri.ObservedDate.Equal(rj.ObservedDate) {
			return ri.ObservedDate.After(rj.ObservedDate)
		}
		return ri.Source < rj.Source
	})
	return cur
}
