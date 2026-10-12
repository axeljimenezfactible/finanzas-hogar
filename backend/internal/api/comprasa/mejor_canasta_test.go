package comprasa_test

import (
	"context"
	"path/filepath"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comprasa"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/semilla"
	"finanzas-hogar/internal/testutil"
)

// igualAComparativa exige que la tarjeta «Mejor canasta» copie de getBasketComparison (motor S5) la ganadora, su total,
// su cobertura, sus precios viejos y la tienda siguiente con su diferencia: ninguna cifra sale de otra cuenta.
func igualAComparativa(t *testing.T, got api.DashboardBestBasket, bc api.BasketComparison) {
	t.Helper()
	if bc.CheapestStoreID == nil || got.StoreID == nil || *got.StoreID != *bc.CheapestStoreID {
		t.Fatalf("ganadora distinta: tarjeta %+v, comparativa %+v", got, bc)
	}
	for i, st := range bc.Stores {
		if st.StoreID != *bc.CheapestStoreID {
			continue
		}
		if *got.TotalCents != st.TotalCents || got.CoveredCount != st.IdenticalCount+st.SubstituteCount ||
			got.StalePriceCount != st.StaleCount || got.AssumesMembership != st.AssumesMembership {
			t.Fatalf("tarjeta %+v contra tienda %+v", got, st)
		}
		if i+1 < len(bc.Stores) && bc.Stores[i+1].DifferenceVsCheapestCents != nil {
			if got.SecondStoreID == nil || *got.SecondStoreID != bc.Stores[i+1].StoreID || *got.DifferenceVsSecondCents != *bc.Stores[i+1].DifferenceVsCheapestCents {
				t.Fatalf("segunda tienda: tarjeta %+v, comparativa %+v", got, bc.Stores[i+1])
			}
		} else if got.SecondStoreID != nil || got.DifferenceVsSecondCents != nil {
			t.Fatalf("sin segunda tienda completa no hay diferencia: %+v", got)
		}
		return
	}
	t.Fatalf("la ganadora no está en stores: %+v", bc)
}

func TestBestBasket_MismaReglaQueLaComparativa(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	monday, _ := dominio.ParseDate("2026-10-12")
	today := dominio.Today(e.s.Clock)
	got, err := comprasa.BestBasket(ctx, e.s.DB, today, monday)
	if err != nil || got.Basket != "shopping_list" || got.StoreID != nil || got.ItemCount != 0 || got.ExpenseID != nil {
		t.Fatalf("sin lista ni compras la tarjeta va vacía: %+v %v", got, err)
	}
	e.manual(map[string]any{"product_id": e.rice, "quantity_milli": 1000000})
	if got, err = comprasa.BestBasket(ctx, e.s.DB, today, monday); err != nil || got.Basket != "shopping_list" {
		t.Fatalf("%+v %v", got, err)
	}
	var bc api.BasketComparison
	res := e.c.Do("GET", "/api/v1/basket-comparison"+w42, nil)
	e.code(res, 200)
	res.JSON(t, &bc)
	igualAComparativa(t, got, bc)
	if got.SecondStoreID == nil {
		t.Fatalf("las dos cadenas tienen arroz: debe haber segunda tienda: %+v", got)
	}
	var pg api.BasketPage
	res = e.c.Do("GET", "/api/v1/baskets"+w42, nil)
	e.code(res, 200)
	if res.JSON(t, &pg); got.ItemCount != pg.Items[0].ItemCount {
		t.Fatalf("item_count %d contra listBaskets %d", got.ItemCount, pg.Items[0].ItemCount)
	}
	// Otra semana sin lista ni compras: vacía.
	if got, err = comprasa.BestBasket(ctx, e.s.DB, today, monday.AddDays(7)); err != nil || got.StoreID != nil || got.ItemCount != 0 {
		t.Fatalf("semana sin lista: %+v %v", got, err)
	}
}

// Con la semilla real (datos/semilla) no hay lista de la semana y no hay precios cargados (los mete el agente externo,
// K53): la canasta es la compra del 10-oct-2026 y la tarjeta queda vacía sin error (PROVISIONAL, S6-gastos-inicio §15).
func TestBestBasket_SemillaSinPrecios(t *testing.T) {
	db, err := dominio.Open(filepath.Join(testutil.TempDir(t), "semilla.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	ctx := context.Background()
	if _, err := dominio.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	opt := semilla.Options{Dir: filepath.Join(testutil.Root(t), "datos", "semilla"), Argon2: &testutil.FastArgon2,
		Passwords: map[string]string{"axel": testutil.FakeCredential, "ximena": testutil.FakeCredential + "-x"}}
	if _, err := semilla.RunWith(ctx, db, opt); err != nil {
		t.Fatal(err)
	}
	var last int64
	if err := db.QueryRow(`SELECT e.id FROM expense e WHERE EXISTS (SELECT 1 FROM expense_item i WHERE i.expense_id = e.id)
		ORDER BY e.expense_date DESC, e.id DESC LIMIT 1`).Scan(&last); err != nil {
		t.Fatal(err)
	}
	monday, _ := dominio.ParseDate("2026-10-12")
	got, err := comprasa.BestBasket(ctx, db, monday, monday)
	if err != nil || got.Basket != "expense" || got.ExpenseID == nil || *got.ExpenseID != last || got.ItemCount == 0 || got.StoreID != nil {
		t.Fatalf("%+v %v", got, err)
	}
}
