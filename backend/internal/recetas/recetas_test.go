package recetas

import (
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Datos de prueba FICTICIOS (testdata/, generados con testdata/generar_testdata.py): 5 recetas propias.

const seedDir = "../../../datos/semilla/"

func base() map[string]any {
	return map[string]any{"today": "2026-10-11", "catalog_file": "testdata/catalogo.json",
		"recipe_book_file": "testdata/recetario.json", "nutrition_file": "testdata/tabla.csv",
		"persons": []any{person("axel", true), person("ximena", true), person("karen", false)}}
}

func person(c string, p bool) map[string]any { return map[string]any{"code": c, "has_meal_plan": p} }

func entry(date, slot, who, recipe string, sm int64) map[string]any {
	if recipe == "" {
		return map[string]any{"date": date, "slot": slot, "person": who, "recipe": nil, "servings_milli": nil, "eating_out": true}
	}
	return map[string]any{"date": date, "slot": slot, "person": who, "recipe": recipe, "servings_milli": sm, "eating_out": false}
}

func quote(p, s, pres string, price, size int64, rec bool) map[string]any {
	return map[string]any{"product": p, "store": s, "presentation": pres, "price_cents": price, "size_milli": size,
		"observed_date": "2026-10-10", "is_stale": false, "recommended": rec}
}

func choice(p, kind string, store, pres any, qty any, bought bool) map[string]any {
	return map[string]any{"product": p, "kind": kind, "store": store, "presentation": pres, "quantity_milli": qty, "is_bought": bought}
}

func eval(in map[string]any, fn string, args any) (any, string) {
	raw, _ := json.Marshal(in)
	a, _ := json.Marshal(args)
	out, err := Eval(raw, fn, a)
	if err != nil {
		return nil, Code(err)
	}
	return out, ""
}

func must[T any](t *testing.T, in map[string]any, fn string, args any) T {
	t.Helper()
	out, code := eval(in, fn, args)
	if code != "" {
		t.Fatalf("%s: error %s", fn, code)
	}
	return out.(T)
}

func wantCode(t *testing.T, in map[string]any, fn string, args any, code string) {
	t.Helper()
	if _, got := eval(in, fn, args); got != code {
		t.Errorf("%s: código %q, se esperaba %q", fn, got, code)
	}
}

func eq(t *testing.T, name string, got, want any) {
	t.Helper()
	g, _ := json.Marshal(got)
	w, _ := json.Marshal(want)
	if string(g) != string(w) {
		t.Errorf("%s: obtuvo %s, se esperaba %s", name, g, w)
	}
}

func item(l *ShoppingList, p, kind string) *ShopItem {
	for i := range l.Items {
		if l.Items[i].Product == p && l.Items[i].Kind == kind {
			return &l.Items[i]
		}
	}
	return nil
}

func wk(w string) map[string]any { return map[string]any{"week": w} }
func rd(r string, s int64) map[string]any {
	return map[string]any{"recipe": r, "servings_milli": s}
}

func testCatalog(t *testing.T) *Catalog {
	c, err := LoadCatalogFile("testdata/catalogo.json")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// realBook devuelve el recetario oficial datos/semilla/recetario.json (S4-01); sin él la prueba falla.
func realBook(t *testing.T) []byte {
	b, err := os.ReadFile(seedDir + "recetario.json")
	if err != nil {
		t.Fatalf("S4-01: falta el recetario oficial: %v", err)
	}
	return b
}

func TestRecetarioValido(t *testing.T) {
	cat := testCatalog(t)
	nut, err := LoadNutritionCSV("testdata/tabla.csv", cat, Review)
	if err != nil {
		t.Fatal(err)
	}
	book, err := LoadRecipeBook("testdata/recetario.json", cat, nut, Review)
	if err != nil || len(book.Rejected) > 0 || len(book.Recipes) != 5 {
		t.Fatalf("testdata: %v %v", err, book.Rejected)
	}
	if book.SeedComplete() == nil {
		t.Error("5 recetas no completan la semilla de 60")
	}
	// Una llave fuera de §3.2 (como las etiquetas de los borradores) rechaza el archivo completo.
	if _, err := ParseRecipeBook([]byte(`{"version":"1.0","recetas":[{"id":"x","etiquetas":["rápida"]}]}`), cat, nut, Review); !errors.Is(err, ErrInvalidArgument) {
		t.Errorf("campo desconocido: %v", err)
	}
	raw := realBook(t)
	report := t.Errorf
	rcat, err := LoadCatalogFile(seedDir + "productos.json")
	if err != nil {
		t.Fatal(err)
	}
	var rnut *Nutrition
	if _, err := os.Stat(seedDir + "nutricion/tabla.csv"); err == nil {
		if rnut, err = LoadNutritionCSV(seedDir+"nutricion/tabla.csv", rcat, Review); err != nil {
			t.Fatalf("tabla.csv real: %v", err)
		}
	}
	rb, err := ParseRecipeBook(raw, rcat, rnut, Review)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range rb.Rejected {
		report("review: receta %s rechazada: %v", r.Recipe, r.Err)
	}
	pb, _ := ParseRecipeBook(raw, rcat, rnut, Production)
	for _, r := range pb.Rejected {
		report("production: receta %s rechazada: %v", r.Recipe, r.Err)
	}
	t.Logf("S4-01 recetario oficial: review %d válidas y %d rechazadas; production %d válidas y %d rechazadas; tabla.csv=%v",
		len(rb.Recipes), len(rb.Rejected), len(pb.Recipes), len(pb.Rejected), rnut != nil)
	if rnut != nil {
		c := NutritionCoverage(pb.Recipes, rnut)
		t.Logf("S4-08 cobertura real (review): únicos %d, con fuente %d (confirmados %d, propuestos %d), sin dato %d %v",
			c.Unique, c.WithSource, c.Confirmed, c.Proposed, len(c.Missing), c.Missing)
	}
	if err := rb.SeedComplete(); err != nil {
		report("%v", err)
	}
}

func TestRecetarioProduccion(t *testing.T) {
	cat := testCatalog(t)
	nut, err := LoadNutritionCSV("testdata/tabla.csv", cat, Production)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile("testdata/recetario.json")
	tampered := strings.Replace(string(b), `"kcal_por_porcion": 650`, `"kcal_por_porcion": 1`, 1)
	book, err := ParseRecipeBook([]byte(tampered), cat, nut, Production)
	if err != nil || len(book.Rejected) > 0 || len(book.Recipes) != 5 {
		t.Fatalf("production debe cargar las 5 sin cotejo: %v %v", err, book.Rejected)
	}
	e := &Engine{Catalog: cat, Recipes: book.Recipes, Nutrition: nut}
	d, err := e.RecipeDetailFor("huevos-a-la-mexicana", 2000)
	if err != nil || d.KcalPerServing != nil {
		t.Fatalf("receta con dato propuesto debe salir incompleta: %v", err)
	}
	eq(t, "missing_nutrition", d.MissingNutrition, []string{"tortilla"})
	if _, err := LoadRecipeBook("testdata/nada.json", cat, nut, Production); err == nil {
		t.Error("archivo inexistente")
	}
	if _, err := ParseRecipeBook([]byte("{"), cat, nut, Production); err == nil {
		t.Error("JSON inválido")
	}
}

func TestEscalado2a4(t *testing.T) {
	in := base()
	delete(in, "recipe_book_file")
	in["recipes"] = []any{map[string]any{"code": "r1", "name": "Prueba", "slots": []string{"lunch"}, "base_servings": 2,
		"prep_minutes": 1, "cook_minutes": 0, "origin": "household", "is_active": true, "ingredients": []any{
			map[string]any{"product": "espagueti", "presentation": "espagueti.la-moderna-250g", "quantity_milli": 250000, "optional": false},
			map[string]any{"product": "aceite-oliva", "presentation": "aceite-oliva.500ml", "quantity_milli": 7000, "optional": false},
			map[string]any{"product": "huevo", "presentation": "huevo.carton-12", "quantity_milli": 1000, "optional": false}}}}
	want := map[int64][]Frac{4000: {{500000, 1}, {14000, 1}, {2000, 1}}, 1000: {{125000, 1}, {3500, 1}, {500, 1}},
		3000: {{375000, 1}, {10500, 1}, {1500, 1}}, 6000: {{750000, 1}, {21000, 1}, {3000, 1}}}
	for s, w := range want {
		out := must[*ScaledRecipe](t, in, "scale_recipe", rd("r1", s))
		for i, x := range out.Ingredients {
			eq(t, x.Product, x.Quantity, w[i])
		}
	}
	eq(t, "unidad", must[*ScaledRecipe](t, in, "scale_recipe", rd("r1", 4000)).Ingredients[2].Unit, UnitPiece)
}

func TestEscaladoRecetario(t *testing.T) {
	var tab []struct {
		Recipe        string `json:"recipe"`
		ServingsMilli int64  `json:"servings_milli"`
		Product       string `json:"product"`
		Num, Den      int64
	}
	b, _ := os.ReadFile("testdata/escalado.json")
	if err := json.Unmarshal(b, &tab); err != nil || len(tab) == 0 {
		t.Fatal(err)
	}
	for _, x := range tab {
		out := must[*ScaledRecipe](t, base(), "scale_recipe", rd(x.Recipe, x.ServingsMilli))
		for _, in := range out.Ingredients {
			if in.Product == x.Product && (in.Quantity != Frac{x.Num, x.Den}) {
				t.Errorf("%s %s ×%d: %v", x.Recipe, x.Product, x.ServingsMilli, in.Quantity)
			}
		}
	}
	raw := realBook(t)
	rcat, _ := LoadCatalogFile(seedDir + "productos.json")
	rb, _ := ParseRecipeBook(raw, rcat, nil, Production)
	n := 0
	for _, r := range rb.Recipes {
		for _, in := range r.Ingredients {
			for _, s := range []int64{4000, 1000, 3000, 6000} {
				q, _ := Scale(in.QuantityMilli, r.BaseServings, s)
				// tabla independiente: q × base × 1000 = cantidad × porciones (multiplicación cruzada exacta)
				lhs := new(big.Rat).Mul(q, big.NewRat(int64(r.BaseServings)*1000, 1))
				if lhs.Cmp(big.NewRat(in.QuantityMilli*s, 1)) != 0 || (s == 4000 && q.Cmp(big.NewRat(2*in.QuantityMilli, 1)) != 0) {
					t.Errorf("%s %s ×%d", r.Code, in.Product, s)
				}
				n++
			}
		}
	}
	t.Logf("S4-02 recetario real: %d recetas, %d escalados, 0 diferencias si no hay errores arriba", len(rb.Recipes), n)
}

func TestEscaladoFracciones(t *testing.T) {
	q, _ := Scale(1, 2, 3000)
	eq(t, "1 × 3/2", fracOf(q), Frac{3, 2})
	wantCode(t, base(), "scale_recipe", rd("arroz-rojo", 0), "invalid_servings")
	wantCode(t, base(), "scale_recipe", rd("no-existe", 1000), "invalid_argument")
	wantCode(t, base(), "scale_recipe", map[string]any{"recipe": "arroz-rojo"}, "missing_field")
}

func planInput(entries ...any) map[string]any {
	in := base()
	in["meal_plan"] = entries
	return in
}

func threeCells() map[string]any {
	return planInput(entry("2026-10-12", "lunch", "axel", "espagueti-con-atun", 2000),
		entry("2026-10-13", "lunch", "axel", "arroz-rojo", 2000), entry("2026-10-14", "dinner", "ximena", "pollo-con-arroz", 4000))
}

func TestListaSinExistencia(t *testing.T) {
	l := must[*ShoppingList](t, threeCells(), "shopping_list", wk("2026-W42"))
	want := map[string][3]int64{"espagueti": {250000, 1, 0}, "atun": {140000, 1, 0}, "elote": {110000, 1, 110000},
		"aceite-oliva": {15000, 1, 485000}, "arroz": {500000, 1, 500000}, "jitomate": {100000, 1, 900000},
		"cebolla": {130000, 1, 870000}, "aceite-soya": {20000, 1, 926000}, "sal": {11000, 1, 989000}, "pollo": {1000000, 1, 0}}
	if len(l.Items) != len(want) {
		t.Fatalf("artículos %d", len(l.Items))
	}
	for p, w := range want {
		it := item(l, p, KindCalculated)
		eq(t, p, []int64{*it.NeededMilli, *it.PackagesToBuy, *it.LeftoverMilli}, w[:])
	}
	eq(t, "recetas del arroz", item(l, "arroz", KindCalculated).Recipes, []string{"arroz-rojo", "pollo-con-arroz"})
	eq(t, "grupos", l.Groups, []ShopGroup{{Vertical: ChannelMarket}, {Vertical: ChannelChain}})
	if len(l.MissingPrices) != 10 || l.EstimatedTotalCents != 0 {
		t.Error("sin cotizaciones todo va a missing_prices")
	}
	wantCode(t, threeCells(), "shopping_list", wk("2026-10-13"), "not_monday")
}

func TestListaExistenciaParcial(t *testing.T) {
	in := threeCells()
	in["inventory"] = []any{
		map[string]any{"code": "i1", "product": "arroz", "recipe": nil, "mode": "stock", "quantity_milli": 300000},
		map[string]any{"code": "i2", "product": "cebolla", "mode": "consumption_rate", "quantity_milli": 180000,
			"counted_date": "2026-10-05", "rate_quantity_milli": 70000, "rate_period_days": 7}}
	l := must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	eq(t, "arroz", *item(l, "arroz", KindCalculated).NeededMilli, 200000)
	eq(t, "cebolla (180000 − 70000 al lunes)", *item(l, "cebolla", KindCalculated).NeededMilli, 20000)
}

func TestListaExistenciaMayor(t *testing.T) {
	in := threeCells()
	in["inventory"] = []any{map[string]any{"code": "i1", "product": "arroz", "mode": "stock", "quantity_milli": 900000}}
	l := must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	if item(l, "arroz", KindCalculated) != nil {
		t.Error("existencia mayor: 0 a comprar, sin artículo")
	}
	for _, it := range l.Items {
		if *it.NeededMilli < 0 || *it.LeftoverMilli < 0 {
			t.Errorf("negativo en %s", it.Product)
		}
	}
}

func TestListaSobrantes(t *testing.T) {
	in := threeCells()
	in["inventory"] = []any{
		map[string]any{"code": "s1", "recipe": "espagueti-con-atun", "mode": "stock", "quantity_milli": 1000},
		map[string]any{"code": "s2", "recipe": "arroz-rojo", "mode": "stock", "quantity_milli": 1000, "expires_date": "2026-10-12"}}
	l := must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	eq(t, "espagueti cubierto a medias", *item(l, "espagueti", KindCalculated).NeededMilli, 125000)
	eq(t, "arroz sin cubrir (sobrante caducado)", *item(l, "arroz", KindCalculated).NeededMilli, 500000)
	p := must[*WeekPlan](t, in, "week_plan", wk("2026-W42"))
	eq(t, "leftovers", p.Leftovers, []Leftover{{"arroz-rojo", 1000}})
}

func TestListaConservaElecciones(t *testing.T) {
	in := threeCells()
	in["list_choices"] = []any{choice("sal", KindCalculated, "walmart", nil, 2000000, true),
		choice("servilletas", KindCalculated, nil, nil, nil, false)}
	l := must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	sal := item(l, "sal", KindCalculated)
	eq(t, "sal", []any{*sal.QuantityMilli, *sal.Store, sal.IsBought, *sal.PackagesToBuy, *sal.LeftoverMilli}, []any{2000000, "walmart", true, 2, 1989000})
	sv := item(l, "servilletas", KindCalculated)
	eq(t, "servilletas sin necesidad", []any{*sv.NeededMilli, sv.Reasons, *sv.PackagesToBuy}, []any{0, []string{}, 0})
	in["list_choices"] = []any{choice("sal", KindCalculated, nil, nil, nil, false), choice("sal", KindCalculated, nil, nil, nil, false)}
	wantCode(t, in, "shopping_list", wk("2026-W42"), "invalid_argument")
	in["list_choices"] = []any{choice("sal", KindManual, nil, nil, nil, false)}
	wantCode(t, in, "shopping_list", wk("2026-W42"), "invalid_argument")
}

func TestListaPrecioYRecomendacion(t *testing.T) {
	in := threeCells()
	in["quotes"] = []any{quote("arroz", "walmart", "arroz.1kg", 2800, 1000000, true), quote("arroz", "bodega", "arroz.750g", 2100, 750000, false)}
	in["stores"] = []any{map[string]any{"code": "walmart", "kind": "chain", "groups": []string{}, "next_visit": nil},
		map[string]any{"code": "bodega", "kind": "chain", "groups": []string{}, "next_visit": nil}}
	a := item(must[*ShoppingList](t, in, "shopping_list", wk("2026-W42")), "arroz", KindCalculated)
	eq(t, "sin elección", []any{*a.Presentation, *a.Store, *a.PriceStore, *a.RecommendedStore, *a.RecommendedPresentation, *a.EstimatedCents, *a.IsStalePrice, a.PriceObservedDate.String()},
		[]any{"arroz.1kg", "walmart", "walmart", "walmart", "arroz.1kg", 2800, false, "2026-10-10"})
	in["list_choices"] = []any{choice("arroz", KindCalculated, "bodega", nil, nil, false)}
	l := must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	a = item(l, "arroz", KindCalculated)
	eq(t, "solo tienda", []any{*a.Presentation, *a.Store, *a.PriceStore, *a.RecommendedStore, *a.RecommendedPresentation, *a.EstimatedCents},
		[]any{"arroz.750g", "bodega", "bodega", "walmart", "arroz.750g", 2100})
	eq(t, "total de la raíz", l.EstimatedTotalCents, 2100)
	in["list_choices"] = []any{choice("arroz", KindCalculated, nil, "arroz.750g", nil, false)}
	a = item(must[*ShoppingList](t, in, "shopping_list", wk("2026-W42")), "arroz", KindCalculated)
	eq(t, "solo presentación", []any{*a.Store, *a.RecommendedStore, *a.RecommendedPresentation, *a.EstimatedCents}, []any{"bodega", "bodega", "arroz.1kg", 2100})
}

func manualCase() map[string]any {
	in := base()
	delete(in, "recipe_book_file")
	mk := func(code, p, pres string, q int64) map[string]any {
		return map[string]any{"code": code, "name": code, "slots": []string{"lunch"}, "base_servings": 2, "origin": "household",
			"ingredients": []any{map[string]any{"product": p, "presentation": pres, "quantity_milli": q, "optional": false}}}
	}
	in["recipes"] = []any{mk("r-arroz", "arroz", "arroz.1kg", 500000), mk("r-frijol", "frijol", "frijol.1kg", 1200000)}
	in["meal_plan"] = []any{entry("2026-10-12", "lunch", "axel", "r-arroz", 2000), entry("2026-10-12", "lunch", "ximena", "r-frijol", 2000)}
	in["quotes"] = []any{quote("arroz", "walmart", "arroz.1kg", 2800, 1000000, true), quote("frijol", "walmart", "frijol.1kg", 3000, 1000000, true),
		quote("aceite-soya", "walmart", "aceite-soya.946ml", 3990, 946000, true)}
	in["list_choices"] = []any{choice("arroz", KindCalculated, nil, nil, 3000000, false), choice("frijol", KindCalculated, nil, nil, 500000, false),
		choice("aceite-soya", KindManual, nil, nil, 1892000, true), choice("servilletas", KindManual, nil, nil, 1000, false)}
	return in
}

func TestListaCantidadFijadaYManual(t *testing.T) {
	l := must[*ShoppingList](t, manualCase(), "shopping_list", wk("2026-W42"))
	row := func(it *ShopItem) []any {
		return []any{it.NeededMilli, it.QuantityMilli, it.PackagesToBuy, it.LeftoverMilli, it.EstimatedCents}
	}
	eq(t, "arroz", row(item(l, "arroz", KindCalculated)), []any{500000, 3000000, 3, 2500000, 8400})
	eq(t, "frijol", row(item(l, "frijol", KindCalculated)), []any{1200000, 500000, 1, 0, 3000})
	eq(t, "aceite manual", row(item(l, "aceite-soya", KindManual)), []any{nil, 1892000, 2, 1892000, 7980})
	sv := item(l, "servilletas", KindManual)
	eq(t, "servilletas manual", append(row(sv), sv.Presentation, sv.Reasons), []any{nil, 1000, nil, nil, nil, nil, []string{"manual"}})
	eq(t, "grupos", l.Groups, []ShopGroup{{ChannelChain, ptr("walmart"), nil, 19380}, {ChannelChain, nil, nil, 0}})
	eq(t, "raíz", []any{l.EstimatedTotalCents, l.MissingPrices}, []any{19380, []string{"servilletas"}})
}

func invItem(code, mode string, qty int64, extra map[string]any) map[string]any {
	m := map[string]any{"code": code, "product": "aceite-soya", "mode": mode, "quantity_milli": qty}
	for k, v := range extra {
		m[k] = v
	}
	return m
}

func TestInventarioModos(t *testing.T) {
	in := base()
	in["inventory"] = []any{invItem("st", "stock", 500000, nil),
		invItem("rt", "consumption_rate", 2838000, map[string]any{"counted_date": "2026-09-11", "rate_quantity_milli": 946000, "rate_period_days": 21}),
		map[string]any{"code": "lo", "product": nil, "recipe": "arroz-rojo", "mode": "stock", "quantity_milli": 1000}}
	r := must[*InventoryReport](t, in, "inventory_status", nil)
	eq(t, "stock", r.Items[0], ItemStatus{"st", 500000, StatusSufficient, nil, nil})
	eq(t, "ritmo (RN-15)", r.Items[1], ItemStatus{"rt", 1486571, StatusSufficient, ptr(int64(315333)), ptr(int64(32))})
	eq(t, "sobrante", r.Items[2].EstimatedQuantityMilli, 1000)
	in["inventory"] = []any{invItem("x", "consumption_rate", 1, nil)}
	wantCode(t, in, "inventory_status", nil, "invalid_argument")
	wantCode(t, base(), "inventory_status", map[string]any{"x": 1}, "unknown_field")
}

func TestInventarioEstado(t *testing.T) {
	in := base()
	in["inventory"] = []any{invItem("a", "stock", 200000, map[string]any{"reorder_level_milli": 200000}),
		invItem("b", "stock", 0, map[string]any{"reorder_level_milli": 200000}), invItem("c", "stock", 0, nil)}
	r := must[*InventoryReport](t, in, "inventory_status", nil)
	eq(t, "estados", []string{r.Items[0].Status, r.Items[1].Status, r.Items[2].Status}, []string{StatusReorder, StatusOut, StatusOut})
	eq(t, "resumen", r.Summary, InventorySummary{0, 1, 2})
}

func TestInventarioMovimientos(t *testing.T) {
	d, _ := ParseDate("2026-10-11")
	eq(t, "compra", PurchaseMovement("a", d, 3, 946000).DeltaMilli, 2838000)
	m, short := PlanConsumptionMovement("a", d, 100, 250)
	eq(t, "consumo", []int64{m.DeltaMilli, short}, []int64{-100, 150})
	eq(t, "sobrante", LeftoverMovement("s", d, 1000).Kind, MoveLeftover)
	a, err := AdjustmentMovement("a", d, ReasonWaste, 500, 300)
	eq(t, "ajuste", []any{a.DeltaMilli, a.Reason, err}, []any{-200, "waste", nil})
	if _, err := AdjustmentMovement("a", d, "robo", 1, 0); err == nil {
		t.Error("motivo inválido")
	}
}

func TestPlanCeldas(t *testing.T) {
	in := planInput(entry("2026-10-12", "breakfast", "axel", "huevos-a-la-mexicana", 2000), entry("2026-10-12", "lunch", "ximena", "", 0),
		entry("2026-10-12", "snack", "axel", "arroz-rojo", 1000), entry("2026-10-12", "breakfast", "ximena", "arroz-rojo", 2000))
	p := must[*WeekPlan](t, in, "week_plan", wk("2026-10-12"))
	eq(t, "orden de celdas", []string{p.Cells[0].Person, p.Cells[1].Person, string(p.Cells[2].Slot), string(p.Cells[3].Slot)}, []string{"axel", "ximena", "lunch", "snack"})
	eq(t, "fuera sin kcal ni costo", []any{p.Cells[2].Kcal, p.Cells[2].CostCents, p.Cells[2].Recipe, p.EatingOutSlots}, []any{nil, nil, nil, 1})
	eq(t, "colación llena (1 porción)", *p.Cells[3].Kcal, 461)
	eq(t, "kcal del día de ximena (fuera no suma)", *p.Days[0].Persons[1].Kcal, *p.Cells[1].Kcal)
	eq(t, "días y personas", []int{len(p.Days), len(p.Days[0].Persons)}, []int{7, 2})
	eq(t, "huecos", p.UnplannedSlots, 20)
}

func TestPlanFuera(t *testing.T) {
	in := planInput(entry("2026-10-12", "lunch", "axel", "", 0), entry("2026-10-12", "lunch", "axel", "arroz-rojo", 2000))
	wantCode(t, in, "week_plan", wk("2026-W42"), "eating_out_conflict")
	in = planInput(entry("2026-10-12", "lunch", "axel", "arroz-rojo", 2000), entry("2026-10-12", "lunch", "axel", "", 0))
	wantCode(t, in, "week_plan", wk("2026-W42"), "eating_out_conflict")
	bad := entry("2026-10-12", "lunch", "axel", "", 0)
	bad["recipe"] = "arroz-rojo"
	wantCode(t, planInput(bad), "week_plan", wk("2026-W42"), "eating_out_conflict")
	wantCode(t, planInput(entry("2026-10-12", "lunch", "axel", "arroz-rojo", -1)), "week_plan", wk("2026-W42"), "invalid_servings")
	wantCode(t, planInput(entry("2026-10-12", "brunch", "axel", "arroz-rojo", 1000)), "week_plan", wk("2026-W42"), "invalid_argument")
	wantCode(t, planInput(entry("2026-10-12", "lunch", "nadie", "arroz-rojo", 1000)), "week_plan", wk("2026-W42"), "invalid_argument")
}

func TestPlanKarenRechazada(t *testing.T) {
	wantCode(t, planInput(entry("2026-10-12", "lunch", "karen", "arroz-rojo", 2000)), "week_plan", wk("2026-W42"), "person_without_plan")
	wantCode(t, base(), "health_profile", map[string]any{"person": "karen"}, "person_without_plan")
}

func TestPlanCasillasSinArmar(t *testing.T) {
	n := must[map[string]any](t, base(), "unplanned_slots", wk("2026-W42"))
	eq(t, "semana vacía", n["unplanned_slots"], 21)
	var es []any
	for _, s := range []string{"breakfast", "lunch", "dinner"} {
		es = append(es, entry("2026-10-12", s, "axel", "arroz-rojo", 2000), entry("2026-10-12", s, "ximena", "", 0))
	}
	es = append(es, entry("2026-10-13", "lunch", "axel", "arroz-rojo", 2000))
	n = must[map[string]any](t, planInput(es...), "unplanned_slots", wk("2026-W42"))
	eq(t, "un día completo y una celda de una sola persona", n["unplanned_slots"], 18)
}

func TestPlanSemanaDomingo(t *testing.T) {
	w := must[map[string]any](t, base(), "default_week", map[string]any{"date": "2026-10-11"})
	eq(t, "domingo abre la semana siguiente", []any{w["week"], w["monday"]}, []any{"2026-W42", "2026-10-12"})
	w = must[map[string]any](t, base(), "default_week", map[string]any{"date": "2026-10-14"})
	eq(t, "miércoles", w["monday"], "2026-10-12")
	wantCode(t, base(), "default_week", map[string]any{"date": "14-oct-2026"}, "invalid_date")
	wantCode(t, base(), "week_plan", wk("2026-W99"), "invalid_date")
	wantCode(t, base(), "week_plan", wk("2026-10-13"), "not_monday")
	eq(t, "semana 1 de 2027", WeekCode(DefaultWeek(mustDate("2027-01-03"))), "2027-W01")
}

func mustDate(s string) Date { d, _ := ParseDate(s); return d }

func bodega() []any {
	return []any{quote("espagueti", "bodega", "espagueti.la-moderna-250g", 1200, 250000, true), quote("atun", "bodega", "atun.lata-140g", 2100, 140000, true),
		quote("elote", "bodega", "elote.lata-220g", 1350, 220000, true), quote("aceite-oliva", "bodega", "aceite-oliva.500ml", 15000, 500000, true)}
}

func TestCostoReceta(t *testing.T) {
	in := base()
	in["quotes"] = bodega()
	d := must[*RecipeDetail](t, in, "recipe_detail", rd("espagueti-con-atun", 2000))
	eq(t, "§8.3", []any{d.TotalCostCents, d.CostPerServingCents, d.Ingredients[2].CostCents}, []any{4425, 2213, 675})
	eq(t, "4 porciones", must[*RecipeDetail](t, in, "recipe_detail", rd("espagueti-con-atun", 4000)).TotalCostCents, 8850)
}

func TestCostoDiaSemana(t *testing.T) {
	in := planInput(entry("2026-10-12", "lunch", "axel", "espagueti-con-atun", 2000), entry("2026-10-12", "dinner", "ximena", "espagueti-con-atun", 1000))
	in["quotes"] = bodega()
	p := must[*WeekPlan](t, in, "week_plan", wk("2026-W42"))
	eq(t, "celdas", []any{p.Cells[0].CostCents, p.Cells[1].CostCents, *p.Cells[0].PriceStore, p.Cells[0].PriceObservedDate.String()}, []any{4425, 2213, "bodega", "2026-10-10"})
	eq(t, "día y semana", []any{p.Days[0].CostCents, p.Days[1].CostCents, p.WeekCostCents}, []any{6638, 0, 6638})
}

func TestCostoSinPrecio(t *testing.T) {
	in := planInput(entry("2026-10-12", "lunch", "axel", "espagueti-con-atun", 2000))
	in["quotes"] = bodega()[:2]
	d := must[*RecipeDetail](t, in, "recipe_detail", rd("espagueti-con-atun", 2000))
	eq(t, "sin precio", []any{d.TotalCostCents, d.CostPerServingCents, d.Ingredients[2].CostCents, d.MissingPrices}, []any{nil, nil, nil, []string{"elote", "aceite-oliva"}})
	p := must[*WeekPlan](t, in, "week_plan", wk("2026-W42"))
	eq(t, "celda, día y semana", []any{p.Cells[0].CostCents, p.Days[0].CostCents, p.WeekCostCents}, []any{nil, nil, nil})
}

func TestCostoCotizacionQueRige(t *testing.T) {
	in := base()
	in["quotes"] = append(bodega(), quote("espagueti", "walmart", "espagueti.la-moderna-250g", 1100, 250000, false))
	in["list_choices"] = []any{choice("espagueti", KindCalculated, "walmart", nil, nil, false)}
	eq(t, "tienda elegida", must[*RecipeDetail](t, in, "recipe_detail", rd("espagueti-con-atun", 2000)).TotalCostCents, 4325)
	in["list_choices"] = []any{choice("espagueti", KindCalculated, "soriana", nil, nil, false)}
	d := must[*RecipeDetail](t, in, "recipe_detail", rd("espagueti-con-atun", 2000))
	eq(t, "elegida sin cotización no cae a la recomendada", d.MissingPrices, []string{"espagueti"})
	in["list_choices"] = []any{choice("espagueti", KindCalculated, nil, "espagueti.la-moderna-250g", nil, false)}
	eq(t, "presentación elegida: casan dos, rige la recomendada", must[*RecipeDetail](t, in, "recipe_detail", rd("espagueti-con-atun", 2000)).TotalCostCents, 4425)
	in["quotes"] = append(bodega(), quote("espagueti", "walmart", "espagueti.la-moderna-250g", 1100, 250000, true))
	wantCode(t, in, "recipe_detail", rd("espagueti-con-atun", 2000), "invalid_argument")
	in["quotes"] = append(bodega(), quote("espagueti", "bodega", "espagueti.la-moderna-250g", 1, 250000, false))
	wantCode(t, in, "recipe_detail", rd("espagueti-con-atun", 2000), "invalid_argument")
}

func TestCostoCantidadFijada(t *testing.T) {
	in := manualCase()
	l := must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	eq(t, "con cantidad fijada", l.EstimatedTotalCents, 19380)
	in["list_choices"] = []any{}
	l = must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	eq(t, "sin fijar", []any{*item(l, "arroz", KindCalculated).EstimatedCents, *item(l, "arroz", KindCalculated).LeftoverMilli,
		*item(l, "frijol", KindCalculated).EstimatedCents, *item(l, "frijol", KindCalculated).LeftoverMilli}, []any{2800, 500000, 6000, 800000})
}

func healthInput(sex, birth, act, goal string, h, w int64, extra map[string]any) map[string]any {
	in := base()
	p := map[string]any{"person": "axel", "height_mm": h, "birth_date": birth, "sex": sex, "activity": act, "goal": goal}
	for k, v := range extra {
		p[k] = v
	}
	in["health_profiles"] = []any{p}
	in["weights"] = []any{map[string]any{"person": "axel", "date": "2026-01-01", "weight_g": 99000}, map[string]any{"person": "axel", "date": "2026-10-10", "weight_g": w}}
	return in
}

func hp(t *testing.T, in map[string]any) *Health {
	return must[*Health](t, in, "health_profile", map[string]any{"person": "axel"})
}

func TestSaludIMC(t *testing.T) {
	for _, c := range []struct {
		h, w, tenths int64
		class        string
	}{{2000, 74000, 185, "normal"}, {2000, 100000, 250, "overweight"}, {2000, 120000, 300, "obesity"}, {1750, 72500, 237, "normal"},
		{2000, 73990, 185, "normal"}, {2000, 60000, 150, "underweight"}} {
		h := hp(t, healthInput("male", "1996-01-01", "sedentary", "maintain", c.h, c.w, nil))
		eq(t, "IMC", []any{*h.BMITenths, *h.BMIClass, *h.Age}, []any{c.tenths, c.class, 30})
	}
}

func TestSaludMeta(t *testing.T) {
	h := hp(t, healthInput("male", "1996-01-01", "sedentary", "maintain", 1750, 80000, nil))
	eq(t, "hombre", []any{*h.SuggestedKcalExact, *h.SuggestedKcalTarget, h.TargetStatus}, []any{Frac{4197, 2}, 2099, TargetSuggested})
	h = hp(t, healthInput("female", "1998-01-01", "sedentary", "gain", 1650, 60000, nil))
	eq(t, "mujer, subir con 0", []any{*h.SuggestedKcalExact, *h.SuggestedKcalTarget}, []any{Frac{15963, 10}, 1596})
	h = hp(t, healthInput("male", "1991-01-01", "moderate", "maintain", 1750, 72500, nil))
	eq(t, "RN-17", []any{*h.SuggestedKcalExact, *h.SuggestedKcalTarget}, []any{Frac{40889, 16}, 2556})
	h = hp(t, healthInput("male", "1991-01-01", "moderate", "lose", 1750, 72500, map[string]any{"adjust_lose_kcal": 500}))
	eq(t, "RN-17 bajar 500 propios", *h.SuggestedKcalTarget, 2056)
	for act, want := range map[string]int64{"light": 2405, "active": 3017, "very_active": 3323} {
		eq(t, act, *hp(t, healthInput("male", "1996-01-01", act, "maintain", 1750, 80000, nil)).SuggestedKcalTarget, want)
	}
	if _, err := SuggestedKcal(1, 1, 1, "x", "light", "lose", 0, 0); err == nil {
		t.Error("sexo inválido")
	}
	if _, err := SuggestedKcal(1, 1, 1, "male", "x", "lose", 0, 0); err == nil {
		t.Error("actividad inválida")
	}
	wantCode(t, healthInput("male", "1996-01-01", "x", "maintain", 1750, 80000, nil), "health_profile", map[string]any{"person": "axel"}, "invalid_argument")
	wantCode(t, healthInput("otro", "1996-01-01", "light", "maintain", 1750, 80000, nil), "health_profile", map[string]any{"person": "axel"}, "invalid_argument")
}

func TestSaludAjustePersonaYCasa(t *testing.T) {
	house := func(in map[string]any) map[string]any {
		in["settings"] = map[string]any{"house_adjust_lose_kcal": 300, "house_adjust_gain_kcal": 250}
		return in
	}
	eq(t, "hombre bajar 300", *hp(t, house(healthInput("male", "1996-01-01", "sedentary", "lose", 1750, 80000, nil))).SuggestedKcalExact, Frac{3597, 2})
	eq(t, "mujer subir 250", *hp(t, house(healthInput("female", "1998-01-01", "sedentary", "gain", 1650, 60000, nil))).SuggestedKcalTarget, 1846)
	m := house(healthInput("male", "1996-01-01", "sedentary", "maintain", 1750, 80000, nil))
	m["settings"] = map[string]any{"house_adjust_lose_kcal": 300, "house_adjust_gain_kcal": 300}
	eq(t, "mantener no ajusta", *hp(t, m).SuggestedKcalTarget, 2099)
	eq(t, "propio gana a la casa", *hp(t, house(healthInput("male", "1996-01-01", "sedentary", "lose", 1750, 80000, map[string]any{"adjust_lose_kcal": 100}))).SuggestedKcalTarget, 1999)
	eq(t, "casa en 0 y 0", *hp(t, healthInput("male", "1996-01-01", "sedentary", "lose", 1750, 80000, nil)).SuggestedKcalTarget, 2099)
	wantCode(t, healthInput("male", "1996-01-01", "sedentary", "lose", 1750, 80000, map[string]any{"adjust_lose_kcal": -1}), "health_profile", map[string]any{"person": "axel"}, "invalid_argument")
	neg := healthInput("male", "1996-01-01", "sedentary", "lose", 1750, 80000, nil)
	neg["settings"] = map[string]any{"house_adjust_lose_kcal": -5}
	wantCode(t, neg, "health_profile", map[string]any{"person": "axel"}, "invalid_argument")
}

func TestSaludMetaConfirmada(t *testing.T) {
	h := hp(t, healthInput("male", "1996-01-01", "sedentary", "maintain", 1750, 80000, map[string]any{"daily_kcal_target": 1900}))
	eq(t, "escrita reemplaza a la sugerida", []any{*h.DailyKcalTarget, h.TargetStatus, *h.SuggestedKcalTarget}, []any{1900, TargetConfirmed, 2099})
	h = hp(t, base())
	eq(t, "sin perfil", []any{h.TargetStatus, h.MissingInputs, h.BMITenths}, []any{TargetNone, []string{"sex", "birth_date", "height_mm", "weight", "activity", "goal"}, nil})
	wantCode(t, healthInput("male", "1996-01-01", "sedentary", "maintain", 1750, 80000, map[string]any{"daily_kcal_target": 100}), "health_profile", map[string]any{"person": "axel"}, "invalid_argument")
	in := planInput(entry("2026-10-12", "lunch", "axel", "espagueti-con-atun", 2000))
	in["health_profiles"] = []any{map[string]any{"person": "axel", "daily_kcal_target": 2000}}
	p := must[*WeekPlan](t, in, "week_plan", wk("2026-W42"))
	eq(t, "kcal del día contra la meta", []any{*p.Days[0].Persons[0].Kcal, *p.Days[0].Persons[0].DailyKcalTarget, *p.Days[0].Persons[0].DifferenceKcal}, []any{1301, 2000, -699})
	eq(t, "% de la meta (RN-17)", *TargetPercentTenths(1850, 2000), 925)
	if TargetPercentTenths(1, 0) != nil {
		t.Error("meta 0")
	}
}

func TestKcalReceta(t *testing.T) {
	d := must[*RecipeDetail](t, base(), "recipe_detail", rd("espagueti-con-atun", 2000))
	eq(t, "kcal", []any{d.TotalKcal, d.KcalPerServing, d.Ingredients[0].Kcal, d.Ingredients[3].Kcal}, []any{1301, 650, 928, 121})
	eq(t, "4 porciones", must[*RecipeDetail](t, base(), "recipe_detail", rd("espagueti-con-atun", 4000)).KcalPerServing, 650)
}

func TestKcalSinDato(t *testing.T) {
	d := must[*RecipeDetail](t, base(), "recipe_detail", rd("frijoles-de-la-olla", 2000))
	eq(t, "sin dato", []any{d.TotalKcal, d.KcalPerServing, d.Ingredients[0].Kcal, d.MissingNutrition}, []any{nil, nil, nil, []string{"frijol"}})
	cat := testCatalog(t)
	if _, ok := gramsMilli(cat.Products["servilletas"], r64(1000)); ok {
		t.Error("pieza sin gramos por unidad es «sin dato»")
	}
	in := planInput(entry("2026-10-12", "lunch", "axel", "frijoles-de-la-olla", 2000))
	eq(t, "día sin dato", must[*WeekPlan](t, in, "week_plan", wk("2026-W42")).Days[0].Persons[0].Kcal, nil)
}

func TestKcalRecetasAutor(t *testing.T) {
	b, _ := os.ReadFile("testdata/recetario.json")
	var f bookFile
	_ = json.Unmarshal(b, &f)
	for _, r := range f.Recetas {
		d := must[*RecipeDetail](t, base(), "recipe_detail", rd(r.ID, 2000))
		eq(t, "S4-08 "+r.ID, d.KcalPerServing, r.Kcal)
	}
}

func TestKcalCargadorCSV(t *testing.T) {
	cat := testCatalog(t)
	hdr := NutritionHeader + "\n"
	ok := "arroz,365,,,,USDA,1,x,2026-10-11,confirmado,axel,2026-10-11\n"
	for name, body := range map[string]string{
		"encabezado":     "producto,kcal\n" + ok,
		"producto":       "nada,1,,,,USDA,1,x,2026-10-11,propuesto,,\n",
		"decimales":      "arroz,1.2345,,,,USDA,1,x,2026-10-11,propuesto,,\n",
		"macro":          "arroz,1,x,,,USDA,1,x,2026-10-11,propuesto,,\n",
		"fuente":         "arroz,1,,,,SMAE,1,x,2026-10-11,propuesto,,\n",
		"id vacío":       "arroz,1,,,,manual, ,x,2026-10-11,propuesto,,\n",
		"fdc":            "arroz,1,,,,USDA,abc,x,2026-10-11,propuesto,,\n",
		"barras":         "arroz,1,,,,OFF,123,x,2026-10-11,propuesto,,\n",
		"duplicado":      ok + ok,
		"estado":         "arroz,1,,,,USDA,1,x,2026-10-11,listo,,\n",
		"propuesto con":  "arroz,1,,,,USDA,1,x,2026-10-11,propuesto,axel,\n",
		"confirmo":       "arroz,1,,,,USDA,1,x,2026-10-11,confirmado,pedro,2026-10-11\n",
		"fecha anterior": "arroz,1,,,,USDA,1,x,2026-10-11,confirmado,axel,2026-10-01\n",
		"fecha inválida": "arroz,1,,,,USDA,1,x,11-oct,propuesto,,\n",
		"columnas":       "arroz,1\n",
		"vacía":          "",
	} {
		b := hdr + body
		if name == "encabezado" || name == "vacía" {
			b = body
		}
		if _, err := ParseNutritionCSV(strings.NewReader(b), cat, Review); err == nil {
			t.Errorf("%s: debió rechazarse", name)
		}
	}
	n, err := ParseNutritionCSV(strings.NewReader(hdr+"huevo,143,12.5,9.51,0.7,USDA,1,x,2026-10-11,confirmado,ximena,2026-10-12\n"), cat, Review)
	if err != nil || *n.fact("huevo", "").ProteinMg != 12500 || n.fact("huevo", "").License != "CC0" {
		t.Errorf("renglón válido: %v", err)
	}
	if _, err := LoadNutritionCSV("testdata/nada.csv", cat, Review); err == nil {
		t.Error("archivo inexistente")
	}
	c := must[Coverage](t, base(), "nutrition_coverage", nil)
	eq(t, "S4-08 cobertura", c, Coverage{13, 12, 11, 1, []string{"frijol"}})
	if _, err := NewNutrition([]NutritionFact{{Product: "a", SourceRef: ""}}, Review); err == nil {
		t.Error("sin fuente")
	}
}

func TestKcalCotejoRevision(t *testing.T) {
	cat := testCatalog(t)
	nut, _ := LoadNutritionCSV("testdata/tabla.csv", cat, Review)
	b, _ := os.ReadFile("testdata/recetario.json")
	for name, r := range map[string][2]string{
		"kcal":         {`"kcal_por_porcion": 650`, `"kcal_por_porcion": 651`},
		"null":         {`"kcal_por_porcion": 461`, `"kcal_por_porcion": null`},
		"fuentes":      {`"id_fuente": "100001"`, `"id_fuente": "100099"`},
		"fuente extra": {`"fuentes_kcal": [`, `"fuentes_kcal": [{"producto_id": "x", "fuente": "USDA", "id_fuente": "1"}, `},
	} {
		s := strings.Replace(string(b), r[0], r[1], 1)
		if s == string(b) {
			t.Fatalf("%s: no se aplicó el cambio", name)
		}
		book, _ := ParseRecipeBook([]byte(s), cat, nut, Review)
		if len(book.Rejected) != 1 || !errors.Is(book.Rejected[0].Err, ErrKcalMismatch) {
			t.Errorf("%s: se esperaba ErrKcalMismatch, %v", name, book.Rejected)
		}
		if pb, _ := ParseRecipeBook([]byte(s), cat, nut, Production); len(pb.Rejected) != 0 {
			t.Errorf("%s: production no coteja", name)
		}
	}
}

func TestNombresEnEspanol(t *testing.T) {
	cat := testCatalog(t)
	nut, _ := LoadNutritionCSV("testdata/tabla.csv", cat, Review)
	b, _ := os.ReadFile("testdata/recetario.json")
	s := strings.Replace(string(b), `"Arroz rojo"`, `"Onions, raw"`, 1)
	book, _ := ParseRecipeBook([]byte(s), cat, nut, Review)
	if len(book.Rejected) != 1 {
		t.Errorf("nombre en inglés aceptado: %v", book.Rejected)
	}
	cat.Products["sal"].Name = "Salt, table"
	book, _ = ParseRecipeBook(b, cat, nut, Review)
	if len(book.Rejected) != 3 {
		t.Errorf("producto con nombre en inglés: %d rechazadas", len(book.Rejected))
	}
}

func TestNutricionSinConfirmar(t *testing.T) {
	cat := testCatalog(t)
	n, err := LoadNutritionCSV("testdata/tabla.csv", cat, Production)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "rechazado", n.Rejected, []RowError{{10, "tortilla", "unconfirmed"}})
	book, _ := LoadRecipeBook("testdata/recetario.json", cat, n, Production)
	c := NutritionCoverage(book.Recipes, n)
	eq(t, "producción", []any{c.WithSource, c.Proposed, c.Missing}, []any{11, 0, []string{"frijol", "tortilla"}})
}

func marketInput() map[string]any {
	in := planInput(entry("2026-10-12", "lunch", "axel", "arroz-rojo", 2000), entry("2026-10-13", "lunch", "axel", "pollo-con-arroz", 2000))
	in["stores"] = []any{map[string]any{"code": "walmart", "kind": "chain", "groups": []string{}, "next_visit": nil},
		map[string]any{"code": "tianguis-mar", "kind": "market", "groups": []string{"frutas-verduras", "carnes"}, "next_visit": "2026-10-13"},
		map[string]any{"code": "mercado-lun", "kind": "market", "groups": []string{"frutas-verduras"}, "next_visit": "2026-10-12"},
		map[string]any{"code": "tianguis-lun", "kind": "market", "groups": []string{"frutas-verduras"}, "next_visit": "2026-10-12"}}
	in["quotes"] = []any{quote("arroz", "walmart", "arroz.1kg", 2800, 1000000, true), quote("jitomate", "bodega", "jitomate.kilo", 2500, 1000000, true)}
	return in
}

func TestCanalAgrupacion(t *testing.T) {
	l := must[*ShoppingList](t, marketInput(), "shopping_list", wk("2026-W42"))
	var gs []string
	for _, g := range l.Groups {
		gs = append(gs, string(g.Vertical)+":"+str(g.Store))
	}
	eq(t, "orden de grupos", gs, []string{"market:mercado-lun", "market:tianguis-mar", "chain:walmart", "chain:"})
	eq(t, "visita", l.Groups[0].VisitDate.String(), "2026-10-12")
	eq(t, "total = Σ grupos", []int64{l.EstimatedTotalCents, l.Groups[0].EstimatedTotalCents + l.Groups[2].EstimatedTotalCents}, []int64{5300, 5300})
	for _, it := range l.Items {
		if it.Store == nil && it.vertical == ChannelMarket {
			t.Errorf("%s sin mercado", it.Product)
		}
	}
}

func TestCanalSinGrupo(t *testing.T) {
	in := base()
	delete(in, "catalog_file")
	delete(in, "recipe_book_file")
	delete(in, "nutrition_file")
	in["catalog"] = map[string]any{"groups": []any{}, "products": []any{map[string]any{"code": "x", "name": "Equis", "base_unit": "g", "group": nil, "grams_per_unit_milli": nil}},
		"presentations": []any{map[string]any{"code": "x.1", "product": "x", "name": "X", "size_milli": 1000}}}
	in["recipes"] = []any{map[string]any{"code": "rx", "base_servings": 2, "ingredients": []any{map[string]any{"product": "x", "presentation": "x.1", "quantity_milli": 1}}}}
	in["meal_plan"] = []any{entry("2026-10-12", "lunch", "axel", "rx", 2000)}
	wantCode(t, in, "shopping_list", wk("2026-W42"), "product_without_channel")
	in["nutrition"] = []any{map[string]any{"product": "x", "presentation": "x.1", "kcal_milli_per_100g": 100000, "source": "store_label", "source_ref": "etiqueta", "status": "confirmed"}}
	eq(t, "dato por presentación", must[*RecipeDetail](t, in, "recipe_detail", rd("rx", 2000)).TotalKcal, 0)
	for _, s := range []string{"usda", "x"} {
		in["nutrition"] = []any{map[string]any{"product": "x", "kcal_milli_per_100g": 1, "source": "usda_fdc", "source_ref": "1", "status": s}}
		if s == "x" {
			in["nutrition"].([]any)[0].(map[string]any)["status"] = "confirmed"
			in["nutrition"].([]any)[0].(map[string]any)["source"] = "otra"
		}
		wantCode(t, in, "recipe_detail", rd("rx", 2000), "invalid_argument")
	}
	b, _ := os.ReadFile("testdata/catalogo.json")
	if _, err := ParseCatalog([]byte(strings.Replace(string(b), `"categoria": "abarrotes"`, `"categoria": "nada"`, 1))); err == nil {
		t.Error("grupo desconocido")
	}
}

func TestCanalMercados(t *testing.T) {
	l := must[*ShoppingList](t, marketInput(), "shopping_list", wk("2026-W42"))
	eq(t, "empate por orden de stores", *item(l, "cebolla", KindCalculated).Store, "mercado-lun")
	eq(t, "solo el tianguis surte carnes", *item(l, "pollo", KindCalculated).Store, "tianguis-mar")
	j := item(l, "jitomate", KindCalculated)
	eq(t, "recomendación de otra tienda no mueve el artículo", []any{*j.Store, *j.RecommendedStore, *j.EstimatedCents}, []any{"mercado-lun", "bodega", 2500})
	in := marketInput()
	in["list_choices"] = []any{choice("jitomate", KindCalculated, "walmart", nil, nil, false)}
	j = item(must[*ShoppingList](t, in, "shopping_list", wk("2026-W42")), "jitomate", KindCalculated)
	eq(t, "cadena elegida", []any{*j.Store, j.EstimatedCents, *j.RecommendedStore}, []any{"walmart", nil, "bodega"})
	in = marketInput()
	in["stores"] = []any{}
	l = must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))
	eq(t, "sin mercados", item(l, "cebolla", KindCalculated).Store, nil)
	in["stores"] = []any{map[string]any{"code": "m", "kind": "market", "groups": []string{}, "next_visit": nil}}
	wantCode(t, in, "shopping_list", wk("2026-W42"), "invalid_argument")
}

func household() *Recipe {
	return &Recipe{Code: "sopa-de-la-casa", Name: "Sopa de la casa", Slots: []SlotCode{Lunch}, BaseServings: 4, Origin: OriginHousehold,
		Steps: []string{"Hierve."}, Ingredients: []Ingredient{{Product: "jitomate", Presentation: "jitomate.kilo", QuantityMilli: 300000, Optional: true}}}
}

func TestRecetaEditable(t *testing.T) {
	cat := testCatalog(t)
	r := household()
	if err := ValidateRecipe(r, cat); err != nil {
		t.Fatal(err)
	}
	bad := []func(*Recipe){
		func(r *Recipe) { r.Code = "Sopa" }, func(r *Recipe) { r.Name = "So" }, func(r *Recipe) { r.Description = strings.Repeat("a", 201) },
		func(r *Recipe) { r.Slots = nil }, func(r *Recipe) { r.Slots = []SlotCode{Lunch, Lunch} }, func(r *Recipe) { r.PrepMinutes = -1 },
		func(r *Recipe) { r.Ingredients = nil }, func(r *Recipe) { r.BaseServings = 51 }, func(r *Recipe) { r.Origin = "otro" },
		func(r *Recipe) { r.Origin = OriginSeed }, func(r *Recipe) { r.Steps = nil }, func(r *Recipe) { r.Steps = []string{" "} },
		func(r *Recipe) { r.Ingredients[0].Product = "nada" }, func(r *Recipe) { r.Ingredients[0].Presentation = "arroz.1kg" },
		func(r *Recipe) { r.Ingredients[0].QuantityMilli = 0 }, func(r *Recipe) { r.Ingredients = append(r.Ingredients, r.Ingredients[0]) },
	}
	for i, f := range bad {
		x := household()
		f(x)
		if ValidateRecipe(x, cat) == nil {
			t.Errorf("caso inválido %d aceptado", i)
		}
	}
	edit := household()
	edit.Code, edit.Origin, edit.Name = "otro-codigo", OriginSeed, "Sopa de jitomate"
	got, err := ApplyRecipeEdit(r, edit, cat)
	if err != nil || got.Code != r.Code || got.Origin != OriginHousehold || got.Name != "Sopa de jitomate" {
		t.Errorf("edición: %v", err)
	}
	edit.Name = ""
	if _, err := ApplyRecipeEdit(r, edit, cat); err == nil {
		t.Error("edición inválida aceptada")
	}
	eq(t, "borrado", []string{RemoveRecipe(r, true), RemoveRecipe(r, false)}, []string{RemoveArchive, RemoveDelete})
	e := &Engine{Catalog: cat, Recipes: []*Recipe{r}}
	d, _ := e.RecipeDetailFor(r.Code, 4000)
	eq(t, "opcional fuera de costo y kcal", len(d.Ingredients), 0)
}

func TestRestriccionesVacias(t *testing.T) {
	in := threeCells()
	in["quotes"] = bodega()
	ref := [2]any{must[*WeekPlan](t, in, "week_plan", wk("2026-W42")), must[*ShoppingList](t, in, "shopping_list", wk("2026-W42"))}
	for _, v := range []any{"", nil} {
		x := threeCells()
		x["quotes"] = bodega()
		x["health_profiles"] = []any{map[string]any{"person": "axel", "dietary_restrictions": v}, map[string]any{"person": "ximena", "dietary_restrictions": v}}
		eq(t, "plan", must[*WeekPlan](t, x, "week_plan", wk("2026-W42")), ref[0])
		eq(t, "lista", must[*ShoppingList](t, x, "shopping_list", wk("2026-W42")), ref[1])
	}
}

func TestSinIA(t *testing.T) {
	files, _ := filepath.Glob("*.go")
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, _ := os.ReadFile(f)
		for _, bad := range []string{"os.Getenv", "os.LookupEnv", "os.Environ", "http", "database/sql", "float32", "float64"} {
			if strings.Contains(string(b), bad) {
				t.Errorf("%s contiene %q", f, bad)
			}
		}
	}
}

func TestEvalErrores(t *testing.T) {
	cases := map[string]func(map[string]any){
		"unknown_field":    func(in map[string]any) { in["extra"] = 1 },
		"missing_field":    func(in map[string]any) { delete(in, "today") },
		"invalid_date":     func(in map[string]any) { in["today"] = "2026-13-01" },
		"invalid_argument": func(in map[string]any) { in["catalog"] = map[string]any{} },
	}
	for code, f := range cases {
		in := base()
		f(in)
		wantCode(t, in, "inventory_status", nil, code)
	}
	for _, f := range []func(map[string]any){
		func(in map[string]any) { delete(in, "catalog_file") },
		func(in map[string]any) { delete(in, "recipe_book_file") },
		func(in map[string]any) { in["nutrition"] = []any{} },
		func(in map[string]any) {
			in["recipes"] = []any{map[string]any{"code": "arroz-rojo", "base_servings": 2, "ingredients": []any{}}}
		},
		func(in map[string]any) { in["persons"] = []any{map[string]any{"code": "axel"}} },
		func(in map[string]any) { in["catalog_file"] = "testdata/nada.json" },
		func(in map[string]any) { in["persons"] = "x" },
	} {
		in := base()
		f(in)
		if _, code := eval(in, "inventory_status", nil); code == "" {
			t.Errorf("entrada inválida aceptada: %v", in)
		}
	}
	wantCode(t, base(), "otra_funcion", nil, "invalid_argument")
	wantCode(t, base(), "health_profile", nil, "missing_field")
	wantCode(t, base(), "health_profile", map[string]any{"person": "nadie"}, "invalid_argument")
	if Code(errors.New("x")) != "invalid_argument" || ErrOverflow.Error() != "overflow" {
		t.Error("códigos")
	}
	var err error
	func() { defer catch(&err); toInt64(new(big.Int).Lsh(big.NewInt(1), 70)) }()
	eq(t, "desbordamiento", Code(err), "overflow")
	if _, err := parseMilli("1.2.3"); err == nil {
		t.Error("decimal inválido")
	}
	if floorRat(big.NewRat(-3, 2)) != -2 || roundHalfUp(big.NewRat(-5, 2)) != -3 {
		t.Error("piso o redondeo negativo")
	}
}
