package precios

import (
	"errors"
	"testing"
	"time"
)

const repo = "../../.."

func d(s string) time.Time {
	t, err := ParseDate(s)
	if err != nil {
		panic(err)
	}
	return t
}

func i64(v int64) *int64 { return &v }

func pr(pres, store string, cents int64, date, src string) PriceObs {
	return PriceObs{Presentation: pres, Store: store, PriceCents: cents, ObservedDate: d(date), Source: src}
}

// testCatalog es un catálogo ILUSTRATIVO para los casos de §14.1, §7.3, §8 y §9.
func testCatalog(t *testing.T) *Catalog {
	t.Helper()
	c, err := NewCatalog(
		[]Group{{"abarrotes", "Abarrotes", "chain"}, {"frutas-verduras", "Frutas y verduras", "market"}},
		[]Product{{"arroz", "Arroz", "g", "abarrotes"}, {"frijol", "Frijol", "g", "abarrotes"},
			{"aceite", "Aceite", "ml", "abarrotes"}, {"jitomate", "Jitomate", "g", "frutas-verduras"},
			{"huevo", "Huevo", "piece", "abarrotes"}},
		[]Presentation{{"arroz-1kg", "arroz", "Arroz 1 kg", 1000000}, {"arroz-750", "arroz", "Arroz 750 g", 750000},
			{"arroz-500", "arroz", "Arroz 500 g", 500000}, {"frijol-1kg", "frijol", "Frijol 1 kg", 1000000},
			{"frijol-500", "frijol", "Frijol 500 g", 500000}, {"aceite-946", "aceite", "Aceite 946 ml", 946000},
			{"jitomate-1kg", "jitomate", "Jitomate 1 kg", 1000000}, {"huevo-12", "huevo", "Huevo 12 piezas", 12000}})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

var (
	stWalmart  = Store{Code: "walmart", Kind: KindChain, Name: "Walmart"}
	stBodega   = Store{Code: "bodega", Kind: KindChain, Name: "Bodega Aurrera"}
	stSams     = Store{Code: "sams", Kind: KindChain, Name: "Sam's Club", RequiresMembership: true, MembershipFeeCents: i64(100000)}
	stTianguis = Store{Code: "tianguis", Kind: KindMarket, Name: "Tianguis", VisitWeekday: 2, Groups: []string{"frutas-verduras"}}
	stMercado  = Store{Code: "mercado", Kind: KindMarket, Name: "Mercado", VisitWeekday: 6}
)

func newBook(t *testing.T, stores []Store, prices []PriceObs, today string) *Book {
	t.Helper()
	b, err := NewBook(testCatalog(t), stores, prices, DefaultSettings(), d(today))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// canon arma la compra de Bodega contra el carrito de Walmart desde data.json (sin copias).
func canon(t *testing.T) (*CartComparison, *Report) {
	t.Helper()
	cat, err := LoadCatalogFile(repo + "/datos/semilla/productos.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := LoadReport(repo+"/datos/semilla/data.json", cat)
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewBook(cat, nil, nil, DefaultSettings(), d("2026-10-10"))
	if err != nil {
		t.Fatal(err)
	}
	cmp, err := b.CompareCarts(r.Bodega, []Cart{r.Walmart})
	if err != nil {
		t.Fatal(err)
	}
	return cmp, r
}

func eq[T comparable](t *testing.T, what string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, se esperaba %v", what, got, want)
	}
}

func TestComparativaCanon(t *testing.T) {
	cmp, r := canon(t)
	bo, wa := cmp.Stores[0], cmp.Stores[1]
	eq(t, "tienda más barata", bo.Store, "bodega")
	eq(t, "K-01 Bodega", bo.TotalCents, 309450)
	eq(t, "K-02 Walmart", wa.TotalCents, 323300)
	eq(t, "K-03 Walmart − Bodega", dif(t, wa), 13850)
	eq(t, "K-04 idénticos Bodega", bo.IdenticalTotalCents, 277350)
	eq(t, "K-05 idénticos Walmart", wa.IdenticalTotalCents, 287400)
	eq(t, "K-06 idénticos Walmart − Bodega", wa.IdenticalTotalCents-bo.IdenticalTotalCents, 10050)
	// Contra data.json: totales.* leídos del archivo.
	eq(t, "totales.bodega", bo.TotalCents, r.Totals["bodega"])
	eq(t, "totales.walmart", wa.TotalCents, r.Totals["walmart"])
	eq(t, "totales.bodega_identicos", bo.IdenticalTotalCents, r.Totals["bodega_identicos"])
	eq(t, "totales.walmart_identicos", wa.IdenticalTotalCents, r.Totals["walmart_identicos"])
	// Descomposición de +138.50 = +100.50 + 59.00 − 21.00 (§7.2).
	sub := map[string]int64{}
	only := map[string]int64{}
	for _, row := range cmp.Rows {
		for _, l := range row.Lines {
			switch row.Kind {
			case "substitute":
				sub[l.Store] += l.TotalCents
			case "only_one":
				only[l.Store] += l.TotalCents
			}
		}
	}
	eq(t, "sustitutos Bodega", sub["bodega"], 30000)
	eq(t, "sustitutos Walmart", sub["walmart"], 35900)
	eq(t, "solo Bodega (sexta lata)", only["bodega"], 2100)
	eq(t, "solo Walmart", only["walmart"], 0)
	eq(t, "descomposición", (wa.IdenticalTotalCents-bo.IdenticalTotalCents)+(sub["walmart"]-sub["bodega"])-only["bodega"], 13850)
	t.Logf("S5-01: Bodega %d, Walmart %d (+%d); idénticos %d / %d (+%d)", bo.TotalCents, wa.TotalCents,
		dif(t, wa), bo.IdenticalTotalCents, wa.IdenticalTotalCents, wa.IdenticalTotalCents-bo.IdenticalTotalCents)
}

func TestComparativaGanadores(t *testing.T) {
	cmp, r := canon(t)
	w := map[string]int{}
	for _, s := range cmp.Winners {
		w[s.Store] = s.Count
	}
	eq(t, "Bodega gana", w["bodega"], 22)
	eq(t, "Walmart gana", w["walmart"], 3)
	eq(t, "empates", cmp.Ties, 8)
	eq(t, "bodega_gana (data.json)", int64(w["bodega"]), r.Winners["bodega_gana"])
	eq(t, "walmart_gana (data.json)", int64(w["walmart"]), r.Winners["walmart_gana"])
	eq(t, "empates (data.json)", int64(cmp.Ties), r.Winners["empates"])
}

func TestComparativaPiezas(t *testing.T) {
	cmp, r := canon(t)
	bo, wa := cmp.Stores[0], cmp.Stores[1]
	eq(t, "K-07 piezas Bodega", bo.PieceCount, 61)
	eq(t, "K-07 piezas Walmart", wa.PieceCount, 60)
	eq(t, "K-07 piezas idénticas", cmp.IdenticalPieces, 55)
	eq(t, "articulos_bodega", bo.PieceCount, r.Totals["articulos_bodega"])
	eq(t, "articulos_walmart", wa.PieceCount, r.Totals["articulos_walmart"])
	eq(t, "articulos_identicos", cmp.IdenticalPieces, r.Totals["articulos_identicos"])
	eq(t, "K-08 renglones idénticos", bo.IdenticalCount, 33)
	eq(t, "K-08 renglones sustitutos", wa.SubstituteCount, 3)
	eq(t, "len(items_identicos)", bo.IdenticalCount, r.Identical)
	eq(t, "len(sustitutos)", bo.SubstituteCount, r.Substitutes)
	eq(t, "faltantes", bo.MissingCount+wa.MissingCount, 0)
}

func TestUnidadSustitutos(t *testing.T) {
	cmp, _ := canon(t)
	want := [][3]any{{int64(2933), int64(2800), "walmart"}, {int64(12692), int64(16406), "bodega"}, {int64(3375), int64(2833), "walmart"}}
	k := 0
	for _, row := range cmp.Rows {
		if row.Kind != "substitute" {
			continue
		}
		if k >= len(want) {
			t.Fatalf("más de 3 sustitutos")
		}
		bu := map[string]int64{}
		for _, l := range row.Lines {
			bu[l.Store] = l.BaseUnitPriceCents
		}
		eq(t, row.Product+" Bodega por unidad base", bu["bodega"], want[k][0].(int64))
		eq(t, row.Product+" Walmart por unidad base", bu["walmart"], want[k][1].(int64))
		if row.WinnerStore == nil {
			t.Fatalf("%s sin ganador", row.Product)
		}
		eq(t, row.Product+" más barato", *row.WinnerStore, want[k][2].(string))
		t.Logf("S5-02: %s %d vs %d → %s", row.Product, bu["bodega"], bu["walmart"], *row.WinnerStore)
		k++
	}
	eq(t, "sustitutos", k, 3)
}

func TestUnidadRedondeo(t *testing.T) {
	for _, c := range []struct {
		price, size int64
		unit        string
		want        int64
		per         string
	}{
		{6600, 2250000, "g", 2933, "kg"}, {8400, 3000000, "g", 2800, "kg"}, {9900, 780000, "g", 12692, "kg"},
		{10500, 640000, "g", 16406, "kg"}, {13500, 4000000, "ml", 3375, "l"}, {17000, 6000000, "ml", 2833, "l"},
		{1, 2000000, "g", 1, "kg"}, {4000, 12000, "piece", 333, "piece"}, {1500, 3000, "piece", 500, "piece"},
	} {
		got, per, exact, err := BaseUnitPrice(c.price, c.size, c.unit)
		if err != nil || got != c.want || per != c.per {
			t.Errorf("BaseUnitPrice(%d, %d, %s) = %d %s %v, se esperaba %d %s", c.price, c.size, c.unit, got, per, err, c.want, c.per)
		}
		if exact.Den <= 0 {
			t.Errorf("fracción inválida %v", exact)
		}
	}
	if _, _, _, err := BaseUnitPrice(-1, 1000, "g"); err == nil {
		t.Error("precio negativo aceptado")
	}
	if _, _, _, err := BaseUnitPrice(1, 1000, "oz"); err == nil {
		t.Error("unidad desconocida aceptada")
	}
	if _, _, _, err := BaseUnitPrice(1<<62, 1000, "g"); !errors.Is(err, ErrOverflow) {
		t.Errorf("desbordamiento: %v", err)
	}
}

// Caso s5-cotizacion-recomendada (§14.1, §7.3).
func cotizacionBook(t *testing.T) *Book {
	return newBook(t, []Store{stWalmart, stBodega}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-750", "walmart", 2200, "2026-10-10", SourceManual),
		pr("arroz-750", "bodega", 2200, "2026-10-10", SourceManual),
		pr("arroz-500", "bodega", 1500, "2026-10-10", SourceManual),
		pr("frijol-1kg", "walmart", 3000, "2026-10-10", SourceManual),
		pr("frijol-500", "bodega", 1500, "2026-10-10", SourceManual),
	}, "2026-10-12")
}

func quote(t *testing.T, b *Book, prod, store, pres string) Quote {
	t.Helper()
	q, ok, err := b.QuoteFor(prod, store, pres, b.Today)
	if err != nil || !ok {
		t.Fatalf("QuoteFor(%s, %s, %s): ok=%v err=%v", prod, store, pres, ok, err)
	}
	return q
}

func TestMejorTiendaCotizacion(t *testing.T) {
	b := cotizacionBook(t)
	q := quote(t, b, "arroz", "", "")
	eq(t, "arroz sin elección", q.Store+"/"+q.Presentation, "walmart/arroz-1kg")
	eq(t, "recommended", q.Recommended, true)
	q = quote(t, b, "arroz", "bodega", "")
	eq(t, "arroz en bodega", q.Presentation, "arroz-750")
	eq(t, "recommended con tienda", q.Recommended, false)
	q = quote(t, b, "arroz", "", "arroz-750")
	eq(t, "arroz 750 g (gana el orden de stores)", q.Store, "walmart")
	q = quote(t, b, "frijol", "", "")
	eq(t, "frijol (gana el menor precio de paquete)", q.Store+"/"+q.Presentation, "bodega/frijol-500")
	q = quote(t, b, "arroz", "bodega", "arroz-500")
	eq(t, "par elegido", q.PriceCents, 1500)
	if _, ok, err := b.QuoteFor("aceite", "", "", b.Today); ok || err != nil {
		t.Errorf("aceite sin precio: ok=%v err=%v", ok, err)
	}
	for _, c := range [][3]string{{"nada", "", ""}, {"arroz", "nada", ""}, {"arroz", "", "frijol-500"}} {
		if _, _, err := b.QuoteFor(c[0], c[1], c[2], b.Today); err == nil {
			t.Errorf("QuoteFor%v sin error", c)
		}
	}
	cmp, err := b.CompareBasket([]ListItem{{Product: "arroz", NeededMilli: i64(500000), RecommendedPresentation: "arroz-1kg"}})
	if err != nil {
		t.Fatal(err)
	}
	it := cmp.Items[0]
	eq(t, "recommended_store", *it.RecommendedStore, "walmart")
	eq(t, "recommended_presentation", *it.RecommendedPresentation, "arroz-1kg")
	eq(t, "estimated_price_cents", *it.EstimatedPriceCents, 2800)
	eq(t, "quantity_milli", it.QuantityMilli, 500000)
	eq(t, "línea Walmart", it.Prices[0].Presentation, "arroz-750")
	eq(t, "línea Walmart sustituto", it.Prices[0].IsSubstitute, true)
	eq(t, "línea Bodega", it.Prices[1].Presentation, "arroz-500")
	eq(t, "cheapest_store", *cmp.CheapestStore, "bodega")
	eq(t, "Bodega", cmp.Stores[0].TotalCents, 1500)
	eq(t, "Walmart", cmp.Stores[1].TotalCents, 2200)
	eq(t, "diferencia", dif(t, cmp.Stores[1]), 700)
}

func TestMejorTiendaCanastaReglaUnica(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-500", "walmart", 1450, "2026-10-10", SourceManual),
		pr("arroz-500", "bodega", 1460, "2026-10-10", SourceManual),
		pr("aceite-946", "walmart", 3990, "2026-10-10", SourceManual),
		pr("aceite-946", "bodega", 4000, "2026-10-10", SourceManual),
	}, "2026-10-12")
	cmp, err := b.CompareBasket([]ListItem{
		{Product: "arroz", NeededMilli: i64(500000), RecommendedPresentation: "arroz-1kg"},
		{Product: "aceite", QuantityMilli: i64(1892000), RecommendedPresentation: "aceite-946"},
	})
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "arroz recomendado", *cmp.Items[0].RecommendedPresentation, "arroz-1kg")
	eq(t, "arroz estimado", *cmp.Items[0].EstimatedPriceCents, 2800)
	eq(t, "aceite recomendado", *cmp.Items[1].RecommendedStore, "walmart")
	eq(t, "aceite estimado (2 paquetes)", *cmp.Items[1].EstimatedPriceCents, 7980)
	eq(t, "aceite cantidad a comprar", cmp.Items[1].QuantityMilli, 1892000)
	eq(t, "aceite idéntico", cmp.Items[1].Prices[1].IsSubstitute, false)
	eq(t, "cheapest_store", *cmp.CheapestStore, "walmart")
	eq(t, "Walmart", cmp.Stores[0].TotalCents, 9430)
	eq(t, "Bodega", cmp.Stores[1].TotalCents, 9460)
	eq(t, "diferencia", dif(t, cmp.Stores[1]), 30)
	eq(t, "idénticos Walmart", cmp.Stores[0].IdenticalCount, 1)
	eq(t, "sustitutos Walmart", cmp.Stores[0].SubstituteCount, 1)
	if _, err := b.CompareBasket([]ListItem{{Product: "arroz"}}); err == nil {
		t.Error("cantidades nulas aceptadas")
	}
	if _, err := b.CompareBasket([]ListItem{{Product: "x", NeededMilli: i64(1)}}); err == nil {
		t.Error("producto desconocido aceptado")
	}
	if _, err := b.CompareBasket([]ListItem{{Product: "arroz", NeededMilli: i64(1), RecommendedPresentation: "aceite-946"}}); err == nil {
		t.Error("presentación ajena aceptada")
	}
}

func TestMejorTiendaCanasta(t *testing.T) {
	cmp, _ := canon(t)
	eq(t, "canasta del 10-oct: más barata", cmp.Stores[0].Store, "bodega")
	// Una tienda con faltantes no puede ser la más barata (§7.4).
	b := newBook(t, []Store{stWalmart, stSams}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("frijol-1kg", "walmart", 3000, "2026-10-10", SourceManual),
		pr("arroz-1kg", "sams", 100, "2026-10-10", SourceManual),
	}, "2026-10-12")
	c, err := b.CompareBasket([]ListItem{{Product: "arroz", NeededMilli: i64(1000000)}, {Product: "frijol", NeededMilli: i64(1000000)}})
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "más barata completa", *c.CheapestStore, "walmart")
	// §7.4 (H-r3-04, Q2-05 = A): primero las completas por total y después las incompletas
	// (Sam's, 1.00 sin frijol), que no pueden ser la más barata y van con diferencia null.
	eq(t, "Walmart primera", c.Stores[0].Store, "walmart")
	last := c.Stores[len(c.Stores)-1]
	eq(t, "Sam's al final aunque su total sea el menor", last.Store, "sams")
	eq(t, "Sam's faltantes", last.MissingCount, 1)
	eq(t, "Sam's con membresía supuesta", last.AssumesMembership, true)
	eq(t, "Sam's diferencia null", last.DifferenceVsCheapestCents == nil, true)
	eq(t, "Walmart diferencia 0", dif(t, storeRow(t, c.Stores, "walmart")), 0)
	// Canasta vacía: todas tienen missing_count 0, así que la más barata es la primera (§7.4 literal).
	empty, _ := b.CompareBasket(nil)
	if empty.CheapestStore == nil || *empty.CheapestStore != empty.Stores[0].Store {
		t.Error("canasta vacía: la más barata debe ser la primera tienda completa")
	}
}

func TestMejorTiendaEmpate(t *testing.T) {
	// Empate 1: mismo total; gana la que tiene más renglones recomendados aunque vaya después.
	b := cotizacionBook(t)
	c, err := b.CompareBasket([]ListItem{{Product: "frijol", NeededMilli: i64(1000000), RecommendedPresentation: "frijol-500"}})
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "empate 1 totales", c.Stores[0].TotalCents, c.Stores[1].TotalCents)
	eq(t, "empate 1 primera", c.Stores[0].Store, "bodega")
	// Empate 2: mismo total y mismos recomendados; gana la de menos precios viejos.
	b = newBook(t, []Store{stBodega, stWalmart}, []PriceObs{
		pr("huevo-12", "walmart", 4000, "2026-10-10", SourceManual),
		pr("huevo-12", "bodega", 4000, "2026-09-01", SourceManual),
		pr("frijol-1kg", "walmart", 3000, "2026-10-10", SourceManual),
		pr("frijol-500", "bodega", 1500, "2026-10-10", SourceManual),
	}, "2026-10-12")
	c, err = b.CompareBasket([]ListItem{{Product: "huevo", NeededMilli: i64(12000)}, {Product: "frijol", NeededMilli: i64(1000000)}})
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "empate 2 totales", c.Stores[0].TotalCents, int64(7000))
	eq(t, "empate 2 primera (menos viejos)", c.Stores[0].Store, "walmart")
	eq(t, "empate 2 viejos Bodega", c.Stores[1].StaleCount, 1)
	eq(t, "huevo recomendado", *c.Items[0].RecommendedStore, "walmart")
	eq(t, "frijol recomendado", *c.Items[1].RecommendedStore, "bodega")
}

func TestMejorTiendaIngrediente(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega, stSams, stTianguis}, []PriceObs{
		pr("jitomate-1kg", "walmart", 3500, "2026-10-10", SourceManual),
		pr("jitomate-1kg", "tianguis", 2500, "2026-10-10", SourceManual),
		pr("arroz-1kg", "tianguis", 100, "2026-10-10", SourceManual),
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-1kg", "sams", 2600, "2026-10-10", SourceManual),
		pr("frijol-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("frijol-1kg", "bodega", 1000, "2026-08-13", SourceManual),
		pr("aceite-946", "walmart", 2800, "2026-10-10", SourceManual),
		pr("aceite-946", "bodega", 1500, "2026-08-13", SourceManual),
	}, "2026-10-12")
	eq(t, "jitomate en el tianguis", quote(t, b, "jitomate", "", "").Store, "tianguis")
	eq(t, "arroz: tianguis no surte abarrotes; Sam's con membresía", quote(t, b, "arroz", "", "").Store, "sams")
	// 60 días: peso 30/61. 1000 ÷ (30/61) = 2033.33 < 2800 gana; 1500 ÷ (30/61) = 3050 > 2800 pierde.
	eq(t, "frijol viejo barato", quote(t, b, "frijol", "", "").Store, "bodega")
	eq(t, "aceite viejo caro", quote(t, b, "aceite", "", "").Store, "walmart")
}

func TestFrescuraFronteras(t *testing.T) {
	s := DefaultSettings()
	for _, c := range []struct {
		today    string
		age      int
		stale    bool
		num, den int64
		show     string
		competes bool
	}{
		{"2026-10-10", 0, false, 1, 1, "1.0000", true},
		{"2026-11-08", 29, false, 1, 1, "1.0000", true},
		{"2026-11-09", 30, true, 60, 61, "0.9836", true},
		{"2026-11-10", 31, true, 59, 61, "0.9672", true},
		{"2026-12-09", 60, true, 30, 61, "0.4918", true},
		{"2027-01-08", 90, true, 0, 1, "0.0000", false},
		{"2027-01-09", 91, true, 0, 1, "0.0000", false},
	} {
		w, err := Freshness(d("2026-10-10"), d(c.today), s, SourceTicketManual)
		if err != nil {
			t.Fatal(err)
		}
		if w.AgeDays != c.age || w.IsStale != c.stale || w.Frac != (Frac{c.num, c.den}) || w.Display() != c.show || w.Competes != c.competes {
			t.Errorf("%s: %+v %s", c.today, w, w.Display())
		}
		t.Logf("S5-04: %s días=%d viejo=%v peso=%d/%d (%s) compite=%v", c.today, w.AgeDays, w.IsStale, w.Frac.Num, w.Frac.Den, w.Display(), w.Competes)
	}
	if _, err := Freshness(d("2026-10-11"), d("2026-10-10"), s, SourceManual); !errors.Is(err, ErrFutureDate) {
		t.Errorf("fecha futura: %v", err)
	}
	if _, err := Freshness(d("2026-10-10"), d("2026-10-10"), s, "otra"); err == nil {
		t.Error("fuente desconocida aceptada")
	}
	if _, err := Freshness(d("2026-10-10"), d("2026-10-10"), Settings{StaleAfterDays: 40, ZeroWeightAfterDays: 30}, SourceManual); err == nil {
		t.Error("parámetros inválidos aceptados")
	}
	// Parámetros configurables (K44): Dv = 10 y D0 = 20 dan (20 − edad)/11.
	w, _ := Freshness(d("2026-10-01"), d("2026-10-11"), Settings{StaleAfterDays: 10, ZeroWeightAfterDays: 20}, SourceManual)
	eq(t, "peso con 10 y 20", w.Frac, Frac{10, 11})
}

func TestFrescuraRecta(t *testing.T) {
	s := DefaultSettings()
	obs := d("2026-10-10")
	prev, _ := Freshness(obs, obs.AddDate(0, 0, 29), s, SourceManual)
	for age := 30; age <= 90; age++ {
		w, _ := Freshness(obs, obs.AddDate(0, 0, age), s, SourceManual)
		// prev − w = 1/61 exacto: prev.Num·61·w.Den − w.Num·61·prev.Den = prev.Den·w.Den
		diff := prev.Frac.Num*61*w.Frac.Den - w.Frac.Num*61*prev.Frac.Den
		if diff != prev.Frac.Den*w.Frac.Den {
			t.Errorf("día %d: la recta no baja 1/61 (%v → %v)", age, prev.Frac, w.Frac)
		}
		if w.Frac.Num >= w.Frac.Den {
			t.Errorf("día %d: peso no menor que 1", age)
		}
		prev = w
	}
}

func TestFuenteMismoPeso(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega, stSams}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-01", SourceManual),
		pr("arroz-1kg", "bodega", 2800, "2026-10-01", SourceTicketManual),
		pr("arroz-1kg", "sams", 2800, "2026-10-01", SourceTicketOCR),
	}, "2026-11-05")
	var w0 Weight
	for i, st := range []string{"walmart", "bodega", "sams"} {
		cur := b.CurrentPrice("arroz-1kg", st, b.Today)
		if cur.Price == nil || !cur.Weight.Competes {
			t.Fatalf("%s no compite", st)
		}
		if i == 0 {
			w0 = cur.Weight
		}
		eq(t, st+" peso", cur.Weight.Frac, w0.Frac)
	}
	eq(t, "peso a 35 días", w0.Frac, Frac{55, 61})
	eq(t, "empate resuelto por el orden de stores", quote(t, b, "arroz", "", "").Store, "walmart")
}

func TestFuenteReferencia(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-1kg", "bodega", 100, "2026-10-10", SourceProfeco),
		pr("arroz-1kg", "bodega", 200, "2026-10-10", SourceSiteAssisted),
	}, "2026-10-12")
	eq(t, "la referencia nunca gana", quote(t, b, "arroz", "", "").Store, "walmart")
	cur := b.CurrentPrice("arroz-1kg", "bodega", b.Today)
	if cur.Price != nil || len(cur.References) != 2 {
		t.Fatalf("referencias: %+v", cur)
	}
	eq(t, "referencia 1", cur.References[0].Source, SourceProfeco)
	c, _ := b.CompareBasket([]ListItem{{Product: "arroz", NeededMilli: i64(1000000)}})
	eq(t, "Bodega sin precio que compita", storeRow(t, c.Stores, "bodega").MissingCount, 1)
	eq(t, "más barata", *c.CheapestStore, "walmart")
}

func TestFuenteTrazabilidad(t *testing.T) {
	b := newBook(t, []Store{stWalmart}, nil, "2026-10-12")
	srcs := []string{SourceTicketManual, SourceTicketOCR, SourceManual, SourceSiteAssisted, SourceProfeco}
	for i, s := range srcs {
		p := pr("arroz-1kg", "walmart", 2800+int64(i), "2026-10-10", s)
		p.SourceRef = "ref-" + s
		if err := b.AddPrice(p); err != nil {
			t.Fatal(err)
		}
		comp, ok := Competes(s)
		eq(t, s+" compite", comp, i < 3)
		eq(t, s+" existe", ok, true)
	}
	for i, p := range b.Prices {
		eq(t, "fuente guardada", p.Source, srcs[i])
		eq(t, "referencia guardada", p.SourceRef, "ref-"+srcs[i])
		eq(t, "tienda guardada", p.Store, "walmart")
	}
	cur := b.CurrentPrice("arroz-1kg", "walmart", b.Today)
	eq(t, "vigente = último que compite (misma fecha: el último capturado)", cur.Price.Source, SourceTicketManual)
	later := pr("arroz-1kg", "walmart", 2700, "2026-10-10", SourceManual)
	later.IsPromotional = false
	if err := b.AddPrice(later); !errors.Is(err, ErrDuplicatePrice) {
		t.Errorf("duplicado: %v", err)
	}
	for _, c := range []struct {
		p    PriceObs
		want error
	}{
		{pr("arroz-1kg", "walmart", 1, "2026-10-13", SourceManual), ErrFutureDate},
		{pr("nada", "walmart", 1, "2026-10-10", SourceManual), ErrUnknownPresentation},
		{pr("arroz-1kg", "nada", 1, "2026-10-10", SourceManual), ErrUnknownStore},
		{pr("arroz-1kg", "walmart", -1, "2026-10-10", SourceManual), ErrInvalidArgument},
		{pr("arroz-1kg", "walmart", 1, "2026-10-10", "otra"), ErrInvalidArgument},
	} {
		if err := b.AddPrice(c.p); !errors.Is(err, c.want) {
			t.Errorf("%+v: %v, se esperaba %v", c.p, err, c.want)
		}
	}
	// A igual fecha rige el de created_at más reciente (supuesto S-8).
	b2 := newBook(t, []Store{stWalmart}, nil, "2026-10-12")
	p1 := pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual)
	p1.CreatedAt = time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	p2 := pr("arroz-1kg", "walmart", 2700, "2026-10-10", SourceTicketManual)
	p2.CreatedAt = time.Date(2026, 10, 10, 9, 0, 0, 0, time.UTC)
	_ = b2.AddPrice(p1)
	_ = b2.AddPrice(p2)
	eq(t, "created_at más reciente", b2.CurrentPrice("arroz-1kg", "walmart", b2.Today).Price.PriceCents, 2800)
}

func TestFuentePromocionVencida(t *testing.T) {
	promo := func(cents int64, date string) PriceObs {
		p := pr("arroz-1kg", "bodega", cents, date, SourceManual)
		p.IsPromotional = true
		return p
	}
	b := newBook(t, []Store{stWalmart, stBodega}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-1kg", "bodega", 2900, "2026-10-10", SourceManual),
		promo(1000, "2026-10-11"), promo(1200, "2026-10-12"),
	}, "2026-10-12")
	cur := b.CurrentPrice("arroz-1kg", "bodega", b.Today)
	if cur.Promotion == nil || cur.Promotion.PriceCents != 1200 {
		t.Fatalf("promoción vigente: %+v", cur.Promotion)
	}
	eq(t, "precio regular", cur.Price.PriceCents, 2900)
	eq(t, "la promoción no compite", quote(t, b, "arroz", "", "").Store, "walmart")
	b.Today = d("2026-10-13")
	if b.CurrentPrice("arroz-1kg", "bodega", b.Today).Promotion != nil {
		t.Error("promoción vencida mostrada")
	}
}

func TestMembresiaEquilibrio(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega, stSams, stTianguis}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-1kg", "bodega", 2900, "2026-10-10", SourceManual),
		pr("arroz-1kg", "sams", 2500, "2026-10-10", SourceManual),
		pr("frijol-1kg", "walmart", 3000, "2026-10-10", SourceManual),
		pr("frijol-1kg", "sams", 3100, "2026-10-10", SourceManual),
	}, "2026-10-12")
	m, err := b.Membership("sams", i64(12500), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "S5-09: 100000 ÷ 12500", *m.BreakEvenPurchases, 8)
	eq(t, "supone membresía", m.AssumesMembership, true)
	m, _ = b.Membership("sams", i64(12000), nil, nil)
	eq(t, "8.33 compras → 9 (hacia arriba)", *m.BreakEvenPurchases, 9)
	hist := []Purchase{{d("2026-10-10"), "walmart"}, {d("2026-06-01"), "bodega"}, {d("2025-10-13"), "sams"},
		{d("2025-10-12"), "walmart"}, {d("2026-10-01"), "tianguis"}}
	list := []ListItem{{Product: "arroz", NeededMilli: i64(2000000)}, {Product: "frijol", NeededMilli: i64(1000000)}}
	m, err = b.Membership("sams", nil, list, hist)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "ahorro de la canasta (solo donde Sam's es más barato)", m.BasketSavingsCents, 600)
	eq(t, "compras al año (cadenas, 365 días)", m.PurchasesPerYear, 3)
	eq(t, "ahorro anual", m.AnnualSavingsCents, 1800)
	eq(t, "equilibrio", *m.BreakEvenPurchases, 167)
	eq(t, "veredicto", m.Verdict, "not_worth_it")
	hist8 := make([]Purchase, 8)
	for i := range hist8 {
		hist8[i] = Purchase{d("2026-10-01").AddDate(0, 0, -30*i), "walmart"}
	}
	m, _ = b.Membership("sams", i64(12500), nil, hist8)
	eq(t, "8 × 12500 cubre 100000", m.Verdict, "worth_it")
	if _, err := b.Membership("nada", nil, nil, nil); err == nil {
		t.Error("tienda desconocida aceptada")
	}
	if _, err := b.Membership("sams", i64(-1), nil, nil); err == nil {
		t.Error("ahorro negativo aceptado")
	}
	if _, err := b.Membership("sams", nil, nil, []Purchase{{d("2026-10-01"), "nada"}}); err == nil {
		t.Error("compra en tienda desconocida aceptada")
	}
}

func TestMembresiaSinCosto(t *testing.T) {
	sams := stSams
	sams.MembershipFeeCents = nil
	b := newBook(t, []Store{stWalmart, sams}, nil, "2026-10-12")
	m, err := b.Membership("sams", i64(12500), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "sin costo capturado", m.Verdict, "missing_fee")
	eq(t, "S5-09: nota con la marca PROVISIONAL", m.Note, "sin costo capturado (PROVISIONAL, S5-09)")
	if m.AnnualFeeCents != nil || m.BreakEvenPurchases != nil {
		t.Error("costo o equilibrio inventados")
	}
	b = newBook(t, []Store{stWalmart, stSams}, nil, "2026-10-12")
	m, _ = b.Membership("sams", i64(0), nil, nil)
	eq(t, "ahorro 0", m.Verdict, "not_worth_it")
	if m.BreakEvenPurchases != nil {
		t.Error("equilibrio con ahorro 0")
	}
}

func TestMercadoVisita(t *testing.T) {
	n, err := NextVisit(stTianguis, d("2026-10-10"))
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "RN-16: sábado 10-oct + martes", n.Format(time.DateOnly), "2026-10-13")
	n, _ = NextVisit(stTianguis, d("2026-10-13"))
	eq(t, "hoy es día de visita (S-10)", n.Format(time.DateOnly), "2026-10-13")
	n, _ = NextVisit(Store{Code: "dom", VisitWeekday: 7}, d("2026-10-12"))
	eq(t, "domingo", n.Format(time.DateOnly), "2026-10-18")
	if _, err := NextVisit(stWalmart, d("2026-10-10")); err == nil {
		t.Error("tienda sin día de visita aceptada")
	}
	if _, err := NewBook(testCatalog(t), []Store{{Code: "m", Kind: KindMarket}}, nil, DefaultSettings(), d("2026-10-10")); err == nil {
		t.Error("mercado sin día de visita aceptado")
	}
}

func TestMercadoGrupo(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stMercado, stTianguis}, []PriceObs{
		pr("jitomate-1kg", "walmart", 3500, "2026-10-10", SourceManual),
		pr("jitomate-1kg", "mercado", 1000, "2026-10-10", SourceManual),
	}, "2026-10-12")
	eq(t, "mercado sin el grupo no compite", quote(t, b, "jitomate", "", "").Store, "walmart")
	if _, ok, _ := b.QuoteFor("jitomate", "mercado", "", b.Today); ok {
		t.Error("mercado sin el grupo cotizado")
	}
	c, _ := b.CompareBasket([]ListItem{{Product: "jitomate", NeededMilli: i64(1000000)}})
	eq(t, "mercado sin el grupo: faltante", storeRow(t, c.Stores, "mercado").MissingCount, 1)
}

// storeRow busca la tienda por código: el orden de stores[] es por total (§7.4).
// dif devuelve difference_vs_cheapest_cents de una tienda completa; falla si es null.
func dif(t *testing.T, s StoreTotals) int64 {
	t.Helper()
	if s.DifferenceVsCheapestCents == nil {
		t.Fatalf("%s: difference_vs_cheapest_cents null", s.Store)
		return 0
	}
	return *s.DifferenceVsCheapestCents
}

func storeRow(t *testing.T, s []StoreTotals, code string) StoreTotals {
	t.Helper()
	for _, x := range s {
		if x.Store == code {
			return x
		}
	}
	t.Fatalf("sin tienda %s en %+v", code, s)
	return StoreTotals{}
}

func TestPreciosViejosDespensa(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-1kg", "bodega", 2900, "2026-10-10", SourceManual),
		pr("huevo-12", "walmart", 4000, "2026-10-12", SourceManual),
		pr("frijol-1kg", "bodega", 3000, "2026-08-01", SourceManual),
	}, "2026-10-12")
	r, err := b.StalePricesForPlan([]string{"arroz", "huevo", "frijol", "jitomate"}, map[string]string{"huevo": "walmart"}, d("2026-11-09"))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Stale) != 1 || r.Stale[0].Store != "walmart" || r.Stale[0].AgeDays != 30 || r.Stale[0].Weight != (Frac{60, 61}) {
		t.Fatalf("viejos: %+v", r.Stale)
	}
	eq(t, "sin precio que compita", len(r.WithoutPrice), 2)
	eq(t, "frijol de 100 días", r.WithoutPrice[0], "frijol")
	r, _ = b.StalePricesForPlan([]string{"arroz"}, map[string]string{"arroz": "bodega"}, d("2026-11-09"))
	eq(t, "tienda asignada", r.Stale[0].Store, "bodega")
	if _, err := b.StalePricesForPlan([]string{"nada"}, nil, d("2026-11-09")); err == nil {
		t.Error("producto desconocido aceptado")
	}
}
