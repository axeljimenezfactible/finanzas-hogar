package gastos_test

import (
	"fmt"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comprasa"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/testutil"
)

// canastaIgual: la tarjeta copia de getBasketComparison (motor S5, misma regla; bloqueo S6-gastos-inicio §15) la ganadora,
// su total, su cobertura, sus precios viejos y la tienda siguiente con su diferencia.
func canastaIgual(t *testing.T, got api.DashboardBestBasket, bc api.BasketComparison) {
	t.Helper()
	if bc.CheapestStoreID == nil || got.StoreID == nil || *got.StoreID != *bc.CheapestStoreID || len(bc.Stores) < 2 {
		t.Fatalf("ganadora: tarjeta %+v, comparativa %+v", got, bc)
	}
	win, next := bc.Stores[0], bc.Stores[1] // S5 §7.4: la ganadora completa va primero
	if win.StoreID != *got.StoreID || *got.TotalCents != win.TotalCents || got.CoveredCount != win.IdenticalCount+win.SubstituteCount ||
		got.StalePriceCount != win.StaleCount || got.AssumesMembership != win.AssumesMembership {
		t.Fatalf("tarjeta %+v contra %+v", got, win)
	}
	if next.DifferenceVsCheapestCents == nil || got.SecondStoreID == nil || *got.SecondStoreID != next.StoreID ||
		*got.DifferenceVsSecondCents != *next.DifferenceVsCheapestCents {
		t.Fatalf("segunda tienda: tarjeta %+v contra %+v", got, next)
	}
}

func TestOp_getDashboardBestBasket(t *testing.T) {
	s := testutil.NewServer(t, gastos.Register, comprasa.Register)
	f := seed(t, s)
	c := s.SessionFor(t, testutil.AxelUser)
	const card = "/api/v1/dashboard/best-basket?week=2026-W42"
	if s.Anonymous(t).Do("GET", card, nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
	want(t, c.Do("GET", "/api/v1/dashboard/best-basket", nil), 422) // week es obligatoria (R3-01)
	want(t, c.Do("GET", "/api/v1/dashboard/best-basket?week=2026-42", nil), 400)

	var bb api.DashboardBestBasket
	res := c.Do("GET", card, nil) // valida el 200 contra openapi.yaml
	want(t, res, 200)
	if res.JSON(t, &bb); bb.Basket != "shopping_list" || bb.StoreID != nil || bb.ItemCount != 0 {
		t.Fatalf("sin lista ni compras la tarjeta va vacía: %+v", bb)
	}

	// Catálogo de prueba: arroz de 1 kg con precio en la tienda de seed() y en otra cadena.
	otra := ins(t, s.DB, `INSERT INTO store (kind, name) VALUES ('chain', 'Walmart prueba')`)
	grp := ins(t, s.DB, `INSERT INTO product_group (name, purchase_channel, sort_order) VALUES ('Prueba abarrotes', 'chain', 1)`)
	var pd api.ProductDetail
	res = c.Do("POST", "/api/v1/products", map[string]any{"name": "Arroz prueba", "base_unit": "g", "product_group_id": grp,
		"presentations": []map[string]any{{"name": "1 kg", "size_milli": 1000000}}}, key(900)...)
	want(t, res, 201)
	res.JSON(t, &pd)
	pres := pd.Presentations[0].ID
	if _, err := s.DB.Exec(`INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source) VALUES
		(?, ?, 2600, '2026-10-10', 'manual'), (?, ?, 2800, '2026-10-10', 'manual')`, pres, f.tienda, pres, otra); err != nil {
		t.Fatal(err)
	}

	// Lista vacía: la última compra con renglones (basket = expense), con la regla de compareExpense.
	exp := newExpense(t, c, map[string]any{"category_id": f.despensa, "amount_cents": 2600, "expense_date": "2026-10-10", "payment_method": "cash",
		"accrual_months": 1, "store_id": f.tienda, "items": []map[string]any{{"line_no": 1, "description": "ARROZ 1KG", "presentation_id": pres,
			"quantity_milli": 1000, "line_total_cents": 2600}}}, 901)
	res = c.Do("GET", card, nil)
	want(t, res, 200)
	var bc api.BasketComparison
	if res.JSON(t, &bb); bb.Basket != "expense" || bb.ExpenseID == nil || *bb.ExpenseID != exp.ID {
		t.Fatalf("con la lista vacía va la última compra: %+v", bb)
	}
	res = c.Do("GET", fmt.Sprintf("/api/v1/basket-comparison?basket=expense&expense_id=%d&week=2026-W42", exp.ID), nil)
	want(t, res, 200)
	res.JSON(t, &bc)
	canastaIgual(t, bb, bc)

	// Lista de la semana con un artículo: basket = shopping_list, misma regla que getBasketComparison.
	want(t, c.Do("POST", "/api/v1/shopping-list/items?week=2026-W42", map[string]any{"product_id": pd.ID, "quantity_milli": 1000000}, key(902)...), 201)
	res = c.Do("GET", card, nil)
	want(t, res, 200)
	bb = api.DashboardBestBasket{} // sin restos de la respuesta anterior (los campos ausentes no se tocan al decodificar)
	if res.JSON(t, &bb); bb.Basket != "shopping_list" || bb.ExpenseID != nil || bb.ItemCount != 1 {
		t.Fatalf("con lista va la lista: %+v", bb)
	}
	res = c.Do("GET", "/api/v1/basket-comparison?week=2026-W42", nil)
	want(t, res, 200)
	bc = api.BasketComparison{}
	res.JSON(t, &bc)
	canastaIgual(t, bb, bc)
	if s.SessionFor(t, testutil.XimenaUser).Do("GET", card, nil).StatusCode != 200 {
		t.Error("Ximena ve la misma tarjeta")
	}
}
