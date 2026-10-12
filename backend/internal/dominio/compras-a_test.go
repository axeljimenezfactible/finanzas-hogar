package dominio_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

// TestShopA_Repositorios recorre los repositorios del área de compras (parte a) contra la base migrada.
func TestShopA_Repositorios(t *testing.T) {
	db, fx := testutil.NewDB(t)
	ctx, r := context.Background(), dominio.ShopA
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	ins := func(table string, cols []string, vals ...any) int64 {
		t.Helper()
		id, err := r.Insert(ctx, db, table, cols, vals...)
		must(err)
		return id
	}
	gid := ins("product_group", []string{"name", "purchase_channel"}, "Prueba grupo", "market")
	sid := ins("store", []string{"kind", "name", "visit_weekday"}, "market", "Mercado prueba", 3)
	must(r.SetStoreGroups(ctx, db, sid, []int64{gid}))
	st, err := r.Store(ctx, db, sid)
	must(err)
	st.Name = "Mercado dos"
	must(r.SaveStore(ctx, db, st))
	sts, err := r.Stores(ctx, db, "market", "ORDER BY id")
	must(err)
	if len(st.GroupIDs) != 1 || len(sts) == 0 || sts[len(sts)-1].Name != "Mercado dos" {
		t.Fatalf("tiendas %+v %+v", st, sts)
	}
	if gs, err := r.Groups(ctx, db, "ORDER BY id"); err != nil || len(gs) == 0 {
		t.Fatalf("grupos %v %v", gs, err)
	}
	pid := ins("product", []string{"name", "base_unit", "product_group_id"}, "Producto prueba", "g", gid)
	p, err := r.Product(ctx, db, pid)
	must(err)
	p.Notes = sql.NullString{String: "nota", Valid: true}
	must(r.SaveProduct(ctx, db, p))
	if ps, err := r.Products(ctx, db, gid, "ORDER BY id"); err != nil || len(ps) != 1 || ps[0].Notes.String != "nota" {
		t.Fatalf("productos %v %v", ps, err)
	}
	prid := ins("presentation", []string{"product_id", "name", "size_milli"}, pid, "1 kg", 1000000)
	pr, err := r.Presentation(ctx, db, prid)
	must(err)
	pr.IsBulk = true
	must(r.SavePresentation(ctx, db, pr))
	if prs, err := r.Presentations(ctx, db, pid); err != nil || len(prs) != 1 || !prs[0].IsBulk {
		t.Fatalf("presentaciones %v %v", prs, err)
	}
	ins("price", []string{"presentation_id", "store_id", "price_cents", "observed_date", "source"}, prid, sid, 2800, "2026-10-10", "manual")
	if prices, err := r.Prices(ctx, db); err != nil || len(prices) == 0 {
		t.Fatalf("precios %v %v", prices, err)
	}
	if nf, err := r.Nutrition(ctx, db, pid); err != nil || len(nf) != 0 {
		t.Fatalf("nutrición %v %v", nf, err)
	}
	lid, created, err := r.EnsureList(ctx, db, "2026-10-12", "2026-10-18")
	must(err)
	again, created2, err := r.EnsureList(ctx, db, "2026-10-12", "2026-10-18")
	must(err)
	if !created || created2 || again != lid {
		t.Fatalf("lista %d %v %d %v", lid, created, again, created2)
	}
	if l, err := r.ListByStart(ctx, db, "2026-10-12"); err != nil || l.ID != lid || l.EndDate != "2026-10-18" {
		t.Fatalf("lista %+v %v", l, err)
	}
	iid := ins("shopping_list_item", []string{"shopping_list_id", "origin", "product_id", "quantity_milli"}, lid, "manual", pid, 1000)
	it, err := r.Item(ctx, db, iid)
	must(err)
	it.IsBought, it.BoughtAt = true, sql.NullString{String: "2026-10-12T18:00:00.000Z", Valid: true}
	must(r.SaveItem(ctx, db, it))
	if items, err := r.Items(ctx, db, lid); err != nil || len(items) != 1 || !items[0].IsBought {
		t.Fatalf("artículos %v %v", items, err)
	}
	rid := ins("recipe", []string{"name", "origin"}, "Receta prueba", "household")
	ins("recipe_ingredient", []string{"recipe_id", "product_id", "quantity_milli"}, rid, pid, 200)
	recs, err := r.Recipes(ctx, db)
	must(err)
	if recs[len(recs)-1].ID != rid || len(recs[len(recs)-1].Ingredients) != 1 {
		t.Fatalf("recetas %+v", recs)
	}
	ins("meal_plan_entry", []string{"plan_date", "meal_slot_id", "person_id", "recipe_id", "servings_milli"}, "2026-10-13", 2, fx.AxelPersonID, rid, 2000)
	if plan, err := r.PlanEntries(ctx, db, "2026-10-12", "2026-10-18"); err != nil || len(plan) != 1 || plan[0].Slot != "lunch" {
		t.Fatalf("plan %v %v", plan, err)
	}
	if pers, err := r.Persons(ctx, db); err != nil || len(pers) < 3 {
		t.Fatalf("personas %v %v", pers, err)
	}
	ins("inventory_item", []string{"product_id", "tracking_mode", "quantity_milli"}, pid, "stock", 500)
	if inv, err := r.Inventory(ctx, db); err != nil || len(inv) == 0 {
		t.Fatalf("inventario %v %v", inv, err)
	}
	if _, err := r.ExpenseBaskets(ctx, db, "ORDER BY e.id"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := r.ExpenseBasket(ctx, db, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("compra inexistente: %v", err)
	}
	if ok, err := r.Exists(ctx, db, "store", sid); err != nil || !ok {
		t.Fatal("la tienda debe existir")
	}
	if ok, err := r.Exists(ctx, db, "store", 999); err != nil || ok {
		t.Fatal("la tienda 999 no existe")
	}
	if err := r.Update(ctx, db, "store", 999, []string{"name"}, "x"); !errors.Is(err, dominio.ErrNotFound) {
		t.Fatalf("Update sin renglón: %v", err)
	}
	must(r.Delete(ctx, db, "shopping_list_item", iid))
	if _, err := r.Item(ctx, db, iid); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("borrado: %v", err)
	}
	if _, err := r.Stores(ctx, db, "", "ORDER BY columna_que_no_existe"); err == nil {
		t.Fatal("un ORDER BY malo debe fallar")
	}
}
