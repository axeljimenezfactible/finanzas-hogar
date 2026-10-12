package comprasa_test

import (
	"fmt"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comprasa"
	"finanzas-hogar/internal/testutil"
)

const w42 = "?week=2026-W42" // lunes 12-oct-2026, el día del reloj de prueba

type env struct {
	t                                       *testing.T
	s                                       *testutil.Server
	c                                       *testutil.Client
	n                                       int
	chain, chain2, market, grocery, veggies int64
	rice, rice1kg, rice500, tomato, tom1kg  int64
}

func (e *env) post(path string, body any) *testutil.Response {
	e.n++
	return e.c.Do("POST", path, body, "Idempotency-Key", fmt.Sprintf("prueba-compras-a-%06d", e.n))
}

func (e *env) id(res *testutil.Response, want int) int64 {
	e.t.Helper()
	if res.StatusCode != want {
		e.t.Fatalf("esperaba %d y llegó %d: %s", want, res.StatusCode, res.Body)
	}
	var x struct {
		ID int64 `json:"id"`
	}
	res.JSON(e.t, &x)
	return x.ID
}

func (e *env) code(res *testutil.Response, want int) {
	e.t.Helper()
	if res.StatusCode != want {
		e.t.Fatalf("esperaba %d y llegó %d: %s", want, res.StatusCode, res.Body)
	}
}

func (e *env) exec(q string, args ...any) int64 {
	e.t.Helper()
	r, err := e.s.DB.Exec(q, args...)
	if err != nil {
		e.t.Fatal(err)
	}
	id, _ := r.LastInsertId()
	return id
}

// newEnv arma fixtures mínimas: dos cadenas, un tianguis (sábado) que surte verduras, arroz y jitomate con precios.
func newEnv(t *testing.T) *env {
	s := testutil.NewServer(t, comprasa.Register)
	e := &env{t: t, s: s, c: s.SessionFor(t, testutil.AxelUser)}
	e.grocery = e.exec("INSERT INTO product_group (name, purchase_channel, sort_order) VALUES ('Prueba abarrotes', 'chain', 1)")
	e.veggies = e.exec("INSERT INTO product_group (name, purchase_channel, sort_order) VALUES ('Prueba verduras', 'market', 2)")
	e.chain = e.id(e.post("/api/v1/stores", map[string]any{"kind": "chain", "name": "Walmart prueba"}), 201)
	e.chain2 = e.id(e.post("/api/v1/stores", map[string]any{"kind": "chain", "name": "Bodega prueba"}), 201)
	e.market = e.id(e.post("/api/v1/stores", map[string]any{"kind": "market", "name": "Tiánguis prueba", "visit_weekday": 6}), 201)
	e.code(e.c.Do("PUT", fmt.Sprintf("/api/v1/stores/%d/product-groups", e.market), map[string]any{"ids": []int64{e.veggies}}), 204)
	var pd api.ProductDetail
	res := e.post("/api/v1/products", map[string]any{"name": "Arroz prueba", "base_unit": "g", "product_group_id": e.grocery,
		"presentations": []map[string]any{{"name": "1 kg", "size_milli": 1000000}, {"name": "500 g", "size_milli": 500000}}})
	e.code(res, 201)
	res.JSON(t, &pd)
	e.rice, e.rice1kg, e.rice500 = pd.ID, pd.Presentations[0].ID, pd.Presentations[1].ID
	res = e.post("/api/v1/products", map[string]any{"name": "Jitomate prueba", "base_unit": "g", "product_group_id": e.veggies,
		"presentations": []map[string]any{{"name": "1 kg", "size_milli": 1000000}}})
	e.code(res, 201)
	res.JSON(t, &pd)
	e.tomato, e.tom1kg = pd.ID, pd.Presentations[0].ID
	e.exec(`INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source) VALUES
		(?, ?, 2800, '2026-10-10', 'manual'), (?, ?, 2600, '2026-10-10', 'manual'), (?, ?, 1500, '2026-10-10', 'manual'), (?, ?, 3500, '2026-10-11', 'manual')`,
		e.rice1kg, e.chain, e.rice1kg, e.chain2, e.rice500, e.chain, e.tom1kg, e.market)
	return e
}

func TestOp_createStore(t *testing.T) {
	e := newEnv(t)
	e.id(e.post("/api/v1/stores", map[string]any{"kind": "branch", "name": "Walmart Coyoacán", "parent_id": e.chain}), 201)
	e.code(e.post("/api/v1/stores", map[string]any{"kind": "market", "name": "Mercado sin día"}), 422)
	e.code(e.post("/api/v1/stores", map[string]any{"kind": "chain", "name": 5}), 400)
	e.code(e.c.Do("POST", "/api/v1/stores", map[string]any{"kind": "chain", "name": "Sin clave"}), 400)
}

func TestOp_updateStore(t *testing.T) {
	e := newEnv(t)
	var st api.Store
	res := e.c.Do("PATCH", fmt.Sprintf("/api/v1/stores/%d", e.chain), map[string]any{"name": "Walmart Taxqueña", "requires_membership": true, "membership_annual_fee_cents": 50000})
	e.code(res, 200)
	if res.JSON(t, &st); st.Name != "Walmart Taxqueña" || st.MembershipAnnualFeeCents == nil {
		t.Fatalf("%+v", st)
	}
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/stores/%d", e.chain), map[string]any{"assisted_reading_enabled": true}), 409)
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/stores/%d", e.chain), map[string]any{"visit_weekday": 3}), 422)
	e.code(e.c.Do("PATCH", "/api/v1/stores/9999", map[string]any{"name": "X"}), 404)
}

func TestOp_setStoreProductGroups(t *testing.T) {
	e := newEnv(t)
	e.code(e.c.Do("PUT", fmt.Sprintf("/api/v1/stores/%d/product-groups", e.market), map[string]any{"ids": []int64{e.veggies, e.grocery}}), 204)
	e.code(e.c.Do("PUT", "/api/v1/stores/9999/product-groups", map[string]any{"ids": []int64{e.veggies}}), 404)
	e.code(e.c.Do("PUT", fmt.Sprintf("/api/v1/stores/%d/product-groups", e.market), map[string]any{"ids": []int64{99999}}), 422)
	// ids es obligatorio (IdList): {} y {"ids": null} no vacían la tienda.
	for _, body := range []any{map[string]any{}, map[string]any{"ids": nil}} {
		e.code(e.c.Do("PUT", fmt.Sprintf("/api/v1/stores/%d/product-groups", e.market), body), 422)
	}
	groups := func() []int64 {
		var pg api.StorePage
		res := e.c.Do("GET", "/api/v1/stores?limit=200", nil)
		e.code(res, 200)
		res.JSON(t, &pg)
		for _, st := range pg.Items {
			if st.ID == e.market {
				return st.ProductGroupIDs
			}
		}
		t.Fatalf("no encontré la tienda %d", e.market)
		return nil
	}
	if g := groups(); len(g) != 2 {
		t.Fatalf("un cuerpo sin ids no debe tocar los grupos: %v", g)
	}
	// Vaciar es explícito: {"ids": []}.
	e.code(e.c.Do("PUT", fmt.Sprintf("/api/v1/stores/%d/product-groups", e.market), map[string]any{"ids": []int64{}}), 204)
	if g := groups(); len(g) != 0 {
		t.Fatalf("{\"ids\": []} debe vaciar los grupos: %v", g)
	}
}

func TestOp_listStores(t *testing.T) {
	e := newEnv(t)
	var pg api.StorePage
	res := e.c.Do("GET", "/api/v1/stores?q=TIANGUIS&kind=market", nil)
	e.code(res, 200)
	if res.JSON(t, &pg); len(pg.Items) != 1 || pg.Items[0].ID != e.market || len(pg.Items[0].ProductGroupIDs) != 1 {
		t.Fatalf("%+v", pg)
	}
	res = e.c.Do("GET", "/api/v1/stores?limit=1", nil)
	if res.JSON(t, &pg); len(pg.Items) != 1 || !pg.NextCursor.Valid {
		t.Fatalf("paginación: %+v", pg)
	}
	e.code(e.c.Do("GET", "/api/v1/stores?kind=otro", nil), 400)
	e.code(e.s.Anonymous(t).Do("GET", "/api/v1/stores", nil), 401)
}

func TestOp_listProductGroups(t *testing.T) {
	e := newEnv(t)
	var pg api.ProductGroupPage
	res := e.c.Do("GET", "/api/v1/product-groups?limit=200", nil)
	e.code(res, 200)
	if res.JSON(t, &pg); len(pg.Items) < 2 {
		t.Fatalf("%+v", pg)
	}
	e.code(e.c.Do("GET", "/api/v1/product-groups?limit=500", nil), 400)
}

func TestOp_createProduct(t *testing.T) {
	e := newEnv(t)
	e.id(e.post("/api/v1/products", map[string]any{"name": "Leche prueba", "base_unit": "ml", "grams_per_unit_milli": 1030}), 201)
	e.code(e.post("/api/v1/products", map[string]any{"name": "arroz PRUEBA", "base_unit": "g"}), 409)
	e.code(e.post("/api/v1/products", map[string]any{"name": "Huevo", "base_unit": "docena"}), 422)
}

func TestOp_getProduct(t *testing.T) {
	e := newEnv(t)
	var pd api.ProductDetail
	res := e.c.Do("GET", fmt.Sprintf("/api/v1/products/%d", e.rice), nil)
	e.code(res, 200)
	res.JSON(t, &pd)
	// S5: el mejor precio por unidad base es el de 500 g a 1500 en Walmart (3000 por kg contra 2600 y 2800 por kg: gana Bodega 1 kg).
	if pd.BestPrice == nil || pd.BestPrice.StoreID != e.chain2 || pd.BestPrice.BaseUnitPriceCents != 2600 || *pd.Freshness != "fresh" || len(pd.Presentations) != 2 {
		t.Fatalf("%+v %+v", pd, pd.BestPrice)
	}
	e.code(e.c.Do("GET", "/api/v1/products/9999", nil), 404)
	e.code(e.c.Do("GET", "/api/v1/products/abc", nil), 400)
}

func TestOp_updateProduct(t *testing.T) {
	e := newEnv(t)
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/products/%d", e.rice), map[string]any{"notes": "grano largo", "product_group_id": nil}), 200)
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/products/%d", e.rice), map[string]any{"grams_per_unit_milli": 5}), 422)
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/products/%d", e.rice), map[string]any{"name": "Jitomate prueba"}), 409)
	e.code(e.c.Do("PATCH", "/api/v1/products/9999", map[string]any{"notes": "x"}), 404)
}

func TestOp_listProducts(t *testing.T) {
	e := newEnv(t)
	var pg api.ProductPage
	res := e.c.Do("GET", fmt.Sprintf("/api/v1/products?freshness=fresh&store_id=%d&q=jitomate", e.market), nil)
	e.code(res, 200)
	if res.JSON(t, &pg); len(pg.Items) != 1 || pg.Items[0].ID != e.tomato || pg.Items[0].BestPrice == nil {
		t.Fatalf("%+v", pg)
	}
	res = e.c.Do("GET", "/api/v1/products?freshness=no_price", nil)
	if res.JSON(t, &pg); len(pg.Items) != 0 {
		t.Fatalf("sin precio: %+v", pg)
	}
	e.code(e.c.Do("GET", "/api/v1/products?freshness=vieja", nil), 400)
	e.code(e.c.Do("GET", "/api/v1/products?source=web", nil), 400)
}

func TestOp_createPresentation(t *testing.T) {
	e := newEnv(t)
	path := fmt.Sprintf("/api/v1/products/%d/presentations", e.rice)
	e.id(e.post(path, map[string]any{"name": "5 kg", "size_milli": 5000000, "barcode": "7501234567890"}), 201)
	e.code(e.post(path, map[string]any{"name": "1 kg", "size_milli": 1000000}), 409)
	e.code(e.post(path, map[string]any{"name": "Vacía", "size_milli": 0}), 422)
	e.code(e.post("/api/v1/products/9999/presentations", map[string]any{"name": "1 kg", "size_milli": 1}), 404)
}

func TestOp_updatePresentation(t *testing.T) {
	e := newEnv(t)
	var p api.Presentation
	res := e.c.Do("PATCH", fmt.Sprintf("/api/v1/presentations/%d", e.rice500), map[string]any{"brand": "Verde Valle", "is_active": false})
	e.code(res, 200)
	if res.JSON(t, &p); p.Brand == nil || p.IsActive {
		t.Fatalf("%+v", p)
	}
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/presentations/%d", e.rice500), map[string]any{"barcode": "12"}), 422)
	e.code(e.c.Do("PATCH", fmt.Sprintf("/api/v1/presentations/%d", e.rice500), map[string]any{"name": "1 kg"}), 409)
	e.code(e.c.Do("PATCH", "/api/v1/presentations/9999", map[string]any{"name": "x"}), 404)
}

func (e *env) manual(body map[string]any) api.ShoppingListItem {
	e.t.Helper()
	var it api.ShoppingListItem
	res := e.post("/api/v1/shopping-list/items"+w42, body)
	e.code(res, 201)
	res.JSON(e.t, &it)
	return it
}

func TestOp_createShoppingListItem(t *testing.T) {
	e := newEnv(t)
	it := e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 1500000})
	// S4/S5: 1.5 kg con la recomendada (Bodega 1 kg a 2600) = 2 paquetes, 5200 estimado.
	if it.Origin != "manual" || it.PackagesToBuy == nil || *it.PackagesToBuy != 2 || it.EstimatedPriceCents == nil || *it.EstimatedPriceCents != 5200 ||
		it.AssignedStoreID == nil || *it.AssignedStoreID != e.chain2 {
		t.Fatalf("%+v", it)
	}
	free := e.manual(map[string]any{"free_text": "Servilletas", "free_text_unit": "piece", "quantity_milli": 2000})
	if free.FreeText == nil || free.ProductID != nil {
		t.Fatalf("%+v", free)
	}
	e.code(e.post("/api/v1/shopping-list/items"+w42, map[string]any{"product_id": e.rice, "free_text": "x", "free_text_unit": "g", "quantity_milli": 1}), 422)
	e.code(e.post("/api/v1/shopping-list/items", map[string]any{"product_id": e.rice, "quantity_milli": 1}), 422)
	e.code(e.post("/api/v1/shopping-list/items?week=2026-42", map[string]any{"product_id": e.rice, "quantity_milli": 1}), 400)
}

func TestOp_getShoppingList(t *testing.T) {
	e := newEnv(t)
	e.manual(map[string]any{"product_id": e.tomato, "quantity_milli": 1000000})
	e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 500000})
	var sl api.ShoppingList
	res := e.c.Do("GET", "/api/v1/shopping-list"+w42, nil)
	e.code(res, 200)
	res.JSON(t, &sl)
	// Q2-42 A: primero el tianguis (sábado 17-oct) y luego la cadena; total = 3500 + 2600.
	if len(sl.Groups) != 2 || sl.Groups[0].PurchaseChannel != "market" || sl.Groups[0].VisitDate == nil || *sl.Groups[0].VisitDate != "2026-10-17" ||
		sl.TotalCents != 6100 || sl.ShoppingListID == nil {
		t.Fatalf("%+v", sl)
	}
	res = e.c.Do("GET", fmt.Sprintf("/api/v1/shopping-list%s&store_id=%d&is_bought=false", w42, e.market), nil)
	if res.JSON(t, &sl); len(sl.Groups) != 1 || sl.TotalCents != 3500 {
		t.Fatalf("filtros: %+v", sl)
	}
	var empty api.ShoppingList
	res = e.c.Do("GET", "/api/v1/shopping-list?week=2026-W50", nil)
	if res.JSON(t, &empty); empty.ShoppingListID != nil || len(empty.Groups) != 0 {
		t.Fatalf("semana sin lista: %+v", empty)
	}
	e.code(e.c.Do("GET", "/api/v1/shopping-list", nil), 422)
	e.code(e.c.Do("GET", "/api/v1/shopping-list?week=hoy", nil), 400)
}

func TestOp_generateShoppingList(t *testing.T) {
	e := newEnv(t)
	rid := e.exec("INSERT INTO recipe (name, base_servings, origin) VALUES ('Arroz rojo prueba', 2, 'household')")
	e.exec("INSERT INTO recipe_ingredient (recipe_id, product_id, presentation_id, quantity_milli) VALUES (?, ?, ?, 250000)", rid, e.rice, e.rice1kg)
	e.exec("INSERT INTO meal_plan_entry (plan_date, meal_slot_id, person_id, recipe_id, servings_milli) VALUES ('2026-10-13', 2, ?, ?, 4000)", e.s.Fx.AxelPersonID, rid)
	var sl api.ShoppingList
	res := e.post("/api/v1/shopping-list/generate"+w42, nil)
	e.code(res, 200)
	res.JSON(t, &sl)
	if sl.GeneratedAt == nil || len(sl.Groups) != 1 || len(sl.Groups[0].Items) != 1 {
		t.Fatalf("%+v", sl)
	}
	it := sl.Groups[0].Items[0]
	// S4: 250 g para 2 porciones escalado a 4 = 500 g; paquete de 1 kg en Bodega a 2600.
	if it.Origin != "calculated" || it.NeededMilli == nil || *it.NeededMilli != 500000 || len(it.RecipeIDs) != 1 || it.RecipeIDs[0] != rid ||
		len(it.Reasons) != 1 || it.Reasons[0] != "recipe" {
		t.Fatalf("%+v", it)
	}
	e.code(e.post("/api/v1/shopping-list/generate"+w42, nil), 200) // regenerar conserva el artículo
	e.code(e.c.Do("POST", "/api/v1/shopping-list/generate"+w42, nil), 400)
	e.code(e.post("/api/v1/shopping-list/generate?week=W42", nil), 400)
}

func TestOp_updateShoppingListItem(t *testing.T) {
	e := newEnv(t)
	it := e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 1000000})
	path := fmt.Sprintf("/api/v1/shopping-list/items/%d", it.ID)
	var got api.ShoppingListItem
	res := e.c.Do("PATCH", path+w42, map[string]any{"is_bought": true, "store_id": e.chain})
	e.code(res, 200)
	if res.JSON(t, &got); !got.IsBought || got.BoughtAt == nil || got.EstimatedPriceCents == nil || *got.EstimatedPriceCents != 2800 {
		t.Fatalf("%+v", got)
	}
	e.code(e.c.Do("PATCH", path+"?week=2026-W43", map[string]any{"is_bought": false}), 422)
	e.code(e.c.Do("PATCH", path+w42, map[string]any{"quantity_milli": nil}), 422)
	e.code(e.c.Do("PATCH", path+w42, map[string]any{}), 400)
	e.code(e.c.Do("PATCH", "/api/v1/shopping-list/items/9999"+w42, map[string]any{"is_bought": true}), 404)
}

func TestOp_assignShoppingListStore(t *testing.T) {
	e := newEnv(t)
	e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 1000000})
	e.manual(map[string]any{"product_id": e.tomato, "quantity_milli": 1000000})
	var r api.ShoppingListStoreAssignmentResult
	res := e.post("/api/v1/shopping-list/assign-store"+w42, map[string]any{"store_id": e.chain})
	e.code(res, 200)
	if res.JSON(t, &r); r.AssignedCount != 1 || r.WithoutPriceCount != 1 {
		t.Fatalf("%+v", r)
	}
	e.code(e.post("/api/v1/shopping-list/assign-store?week=2026-W50", map[string]any{"store_id": e.chain}), 404)
	e.code(e.post("/api/v1/shopping-list/assign-store"+w42, map[string]any{"store_id": 9999}), 422)
}

func TestOp_listBaskets(t *testing.T) {
	e := newEnv(t)
	e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 1000000})
	var pg api.BasketPage
	res := e.c.Do("GET", "/api/v1/baskets"+w42, nil)
	e.code(res, 200)
	if res.JSON(t, &pg); len(pg.Items) != 1 || pg.Items[0].Basket != "shopping_list" || pg.Items[0].ItemCount != 1 {
		t.Fatalf("%+v", pg)
	}
	e.code(e.c.Do("GET", "/api/v1/baskets", nil), 422)
	e.code(e.c.Do("GET", "/api/v1/baskets"+w42+"&limit=0", nil), 400)
}

func TestOp_getBasketComparison(t *testing.T) {
	e := newEnv(t)
	e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 1000000})
	var bc api.BasketComparison
	res := e.c.Do("GET", "/api/v1/basket-comparison"+w42, nil)
	e.code(res, 200)
	if res.JSON(t, &bc); bc.Basket != "shopping_list" || len(bc.Stores) == 0 || len(bc.Items) != 1 || bc.CheapestStoreID == nil {
		t.Fatalf("%+v", bc)
	}
	e.code(e.c.Do("GET", "/api/v1/basket-comparison", nil), 422)
	e.code(e.c.Do("GET", "/api/v1/basket-comparison?basket=expense"+"&week=2026-W42", nil), 422)
	e.code(e.c.Do("GET", "/api/v1/basket-comparison?basket=expense&expense_id=9999&week=2026-W42", nil), 422)
	e.code(e.c.Do("GET", "/api/v1/basket-comparison?basket=otra&week=2026-W42", nil), 400)
}
