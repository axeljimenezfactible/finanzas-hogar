package comprasa

import (
	"context"
	"database/sql"
	"net/http"
	"sort"
	"strings"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
	"finanzas-hogar/internal/recetas"
)

// engineErr: el motor no pudo calcular con los datos actuales (p. ej. un producto sin grupo de compra).
func engineErr(status int, err error) error {
	code := "validation_failed"
	if status == http.StatusConflict {
		code = "conflict"
	}
	return &api.Problem{Status: status, Code: code, Message: "No se pudo calcular con el catálogo actual: " + err.Error() + ". Revisa productos, tiendas y precios en Ajustes."}
}

func hasChoice(it dominio.ShopListItem) bool {
	return it.StoreID.Valid || it.PresentationID.Valid || it.QuantityMilli.Valid || it.IsBought
}

// choicesOf convierte los artículos guardados en elecciones de S4 (list_choices) y dice de qué artículo sale cada una.
// allCalc=false deja fuera los artículos calculated sin elección de una persona (generar los recalcula desde el plan).
func choicesOf(items []dominio.ShopListItem, allCalc bool) ([]recetas.ListChoice, []int) {
	var out []recetas.ListChoice
	var idx []int
	seen := map[int64]bool{}
	for i, it := range items {
		if !it.ProductID.Valid {
			continue
		}
		if it.Origin == recetas.KindCalculated {
			if seen[it.ProductID.Int64] || (!allCalc && !hasChoice(it)) {
				continue
			}
			seen[it.ProductID.Int64] = true
		}
		out = append(out, recetas.ListChoice{Product: code(it.ProductID.Int64), Kind: it.Origin, Store: codeOf(it.StoreID),
			Presentation: codeOf(it.PresentationID), QuantityMilli: pi(it.QuantityMilli), IsBought: it.IsBought})
		idx = append(idx, i)
	}
	return out, idx
}

func baseItem(it dominio.ShopListItem) api.ShoppingListItem {
	out := api.ShoppingListItem{ID: it.ID, ShoppingListID: it.ListID, Origin: it.Origin, ProductID: pi(it.ProductID), FreeText: ps(it.FreeText),
		FreeTextUnit: ps(it.FreeTextUnit), PresentationID: pi(it.PresentationID), NeededMilli: pi(it.NeededMilli), QuantityMilli: pi(it.QuantityMilli),
		StoreID: pi(it.StoreID), IsBought: it.IsBought, Notes: ps(it.Notes), CreatedAt: api.APITime(it.CreatedAt), UpdatedAt: api.APITime(it.UpdatedAt),
		Reasons: []string{}, RecipeIDs: []int64{}}
	if it.BoughtAt.Valid {
		out.BoughtAt = api.Ptr(api.APITime(it.BoughtAt.String))
	}
	if it.Origin == recetas.KindManual {
		out.Reasons = []string{"manual"}
	}
	return out
}

// built es la lista de la semana guardada más lo que calculan S4 y S5 al consultar.
type built struct {
	list  *dominio.ShopList
	rows  []dominio.ShopListItem
	out   api.ShoppingList
	items map[int64]api.ShoppingListItem
	world *world
}

func (h *handlers) buildList(ctx context.Context, q dominio.Querier, monday dominio.Date, status int) (*built, error) {
	w, err := h.loadWorld(ctx, q)
	if err != nil {
		return nil, err
	}
	b := &built{world: w, items: map[int64]api.ShoppingListItem{},
		out: api.ShoppingList{StartDate: monday.String(), EndDate: monday.AddDays(6).String(), Groups: []api.ShoppingListGroupsItem{}}}
	l, err := dominio.ShopA.ListByStart(ctx, q, monday.String())
	if isNoRows(err) {
		return b, nil
	} else if err != nil {
		return nil, err
	}
	b.list, b.out.ShoppingListID = &l, api.Ptr(l.ID)
	if l.GeneratedAt.Valid {
		b.out.GeneratedAt = api.Ptr(api.APITime(l.GeneratedAt.String))
	}
	if b.rows, err = dominio.ShopA.Items(ctx, q, l.ID); err != nil {
		return nil, err
	}
	choices, idx := choicesOf(b.rows, true)
	e, err := w.engine(ctx, q, monday, choices)
	if err != nil {
		return nil, engineErr(status, err)
	}
	sl, err := e.ShoppingListFor(rdate(monday.String()))
	if err != nil {
		return nil, engineErr(status, err)
	}
	calc, man := map[string]int{}, []int{}
	for k, ch := range choices {
		if ch.Kind == recetas.KindCalculated {
			calc[ch.Product] = idx[k]
		} else {
			man = append(man, idx[k])
		}
	}
	shown := make([]*api.ShoppingListItem, len(b.rows))
	m := 0
	for _, si := range sl.Items {
		i, ok := calc[si.Product]
		if si.Kind == recetas.KindManual {
			ok = m < len(man)
			if ok {
				i, m = man[m], m+1
			}
		}
		if !ok {
			continue // necesidad nueva del plan sin artículo guardado: aparece al generar (N-127)
		}
		it := baseItem(b.rows[i])
		it.PackagesToBuy, it.EstimatedPriceCents, it.IsStalePrice = si.PackagesToBuy, si.EstimatedCents, si.IsStalePrice
		it.AssignedStoreID, it.RecommendedStoreID = idPtr(si.Store), idPtr(si.RecommendedStore)
		it.RecommendedPresentationID = idPtr(si.RecommendedPresentation)
		if si.PriceObservedDate != nil {
			it.PriceObservedDate = api.Ptr(si.PriceObservedDate.String())
		}
		it.Reasons = append([]string{}, si.Reasons...)
		for _, rc := range si.Recipes {
			it.RecipeIDs = append(it.RecipeIDs, idOf(rc))
		}
		shown[i] = &it
	}
	type gkey struct {
		ch    string
		store int64
	}
	groups := map[gkey]*api.ShoppingListGroupsItem{}
	var order []gkey
	for i, row := range b.rows {
		it := shown[i]
		if it == nil { // texto libre (S4 no lo calcula) o renglón repetido: solo lo guardado
			x := baseItem(row)
			x.AssignedStoreID = pi(row.StoreID)
			it = &x
		}
		ch := "chain"
		if row.ProductID.Valid {
			ch = w.channelOf(row.ProductID.Int64)
		} else if s, ok := w.store[row.StoreID.Int64]; ok && s.Kind == "market" {
			ch = "market" // PROVISIONAL: texto libre sigue a la tienda elegida (bloqueo S6-compras-a-motores)
		}
		k := gkey{ch: ch}
		if it.AssignedStoreID != nil {
			k.store = *it.AssignedStoreID
		}
		g := groups[k]
		if g == nil {
			g = &api.ShoppingListGroupsItem{PurchaseChannel: ch, Items: []api.ShoppingListItem{}}
			if k.store != 0 {
				g.StoreID = api.Ptr(k.store)
				if s := w.store[k.store]; s.Kind == "market" && w.visit[k.store] != "" {
					g.VisitDate, g.VisitWeekday = api.Ptr(w.visit[k.store]), pi(s.VisitWeekday)
				}
			}
			groups[k] = g
			order = append(order, k)
		}
		g.Items = append(g.Items, *it)
		if it.EstimatedPriceCents != nil {
			g.TotalCents += *it.EstimatedPriceCents
		}
		b.items[row.ID] = *it
	}
	sort.SliceStable(order, func(i, j int) bool { // Q2-42 A: mercado primero; visita ascendente (sin fecha al final); id de tienda; sin tienda al final
		a, c := groups[order[i]], groups[order[j]]
		if a.PurchaseChannel != c.PurchaseChannel {
			return a.PurchaseChannel == "market"
		}
		if (a.VisitDate == nil) != (c.VisitDate == nil) {
			return a.VisitDate != nil
		}
		if a.VisitDate != nil && *a.VisitDate != *c.VisitDate {
			return *a.VisitDate < *c.VisitDate
		}
		if (order[i].store == 0) != (order[j].store == 0) {
			return order[j].store == 0
		}
		return order[i].store < order[j].store
	})
	rowOf := map[int64]dominio.ShopListItem{}
	for _, row := range b.rows {
		rowOf[row.ID] = row
	}
	for _, k := range order {
		g := groups[k]
		sort.SliceStable(g.Items, func(i, j int) bool { // categoría de producto, luego orden alfabético
			ri, rj := rowOf[g.Items[i].ID], rowOf[g.Items[j].ID]
			ci, ni := w.categoryOrder(ri.ProductID)
			cj, nj := w.categoryOrder(rj.ProductID)
			if ci != cj {
				return ci < cj
			}
			ni, nj = ni+ri.FreeText.String, nj+rj.FreeText.String
			if !strings.EqualFold(ni, nj) {
				return strings.ToLower(ni) < strings.ToLower(nj)
			}
			return g.Items[i].ID < g.Items[j].ID
		})
		b.out.Groups = append(b.out.Groups, *g)
		b.out.TotalCents += g.TotalCents
	}
	return b, nil
}

func (h *handlers) getShoppingList(w http.ResponseWriter, r *http.Request) error {
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	stores, err := api.QueryIDs(r, "store_id")
	if err != nil {
		return err
	}
	bought := r.URL.Query().Get("is_bought")
	if bought != "" && !oneOf(bought, "true", "false") {
		return dominio.BadRequest("El filtro is_bought debe ser true o false.")
	}
	b, err := h.buildList(r.Context(), h.db(), monday, http.StatusUnprocessableEntity)
	if err != nil {
		return err
	}
	if len(stores) == 0 && bought == "" {
		return api.WriteJSON(w, http.StatusOK, b.out)
	}
	out := b.out
	out.Groups, out.TotalCents = []api.ShoppingListGroupsItem{}, 0
	for _, g := range b.out.Groups {
		keep := len(stores) == 0
		for _, s := range stores {
			keep = keep || (g.StoreID != nil && *g.StoreID == s)
		}
		if !keep {
			continue
		}
		items := []api.ShoppingListItem{}
		for _, it := range g.Items {
			if bought == "" || it.IsBought == (bought == "true") {
				items = append(items, it)
			}
		}
		if len(items) > 0 {
			g.Items = items
			out.Groups, out.TotalCents = append(out.Groups, g), out.TotalCents+g.TotalCents
		}
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) generateShoppingList(w http.ResponseWriter, r *http.Request) error {
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	ctx := r.Context()
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		listID, created, err := dominio.ShopA.EnsureList(ctx, tx, monday.String(), monday.AddDays(6).String())
		if err != nil {
			return err
		}
		if created {
			if err := audit(ctx, tx, r, "insert", "shopping_list", listID, nil); err != nil {
				return err
			}
		}
		rows, err := dominio.ShopA.Items(ctx, tx, listID)
		if err != nil {
			return err
		}
		wd, err := h.loadWorld(ctx, tx)
		if err != nil {
			return err
		}
		choices, _ := choicesOf(rows, false)
		e, err := wd.engine(ctx, tx, monday, choices)
		if err != nil {
			return engineErr(http.StatusConflict, err)
		}
		sl, err := e.ShoppingListFor(rdate(monday.String()))
		if err != nil {
			return engineErr(http.StatusConflict, err)
		}
		calc := map[string]dominio.ShopListItem{}
		for _, row := range rows {
			if _, dup := calc[codeOf(row.ProductID)]; row.Origin == recetas.KindCalculated && !dup {
				calc[codeOf(row.ProductID)] = row
			}
		}
		keep := map[int64]bool{}
		for _, si := range sl.Items {
			if si.Kind != recetas.KindCalculated {
				continue
			}
			need := deref(si.NeededMilli, 0)
			row, ok := calc[si.Product]
			if !ok {
				id, err := dominio.ShopA.Insert(ctx, tx, "shopping_list_item", []string{"shopping_list_id", "origin", "product_id", "needed_milli"},
					listID, recetas.KindCalculated, idOf(si.Product), need)
				if err == nil {
					err = audit(ctx, tx, r, "insert", "shopping_list_item", id, nil)
				}
				if err != nil {
					return err
				}
				continue
			}
			keep[row.ID] = true
			if row.NeededMilli.Valid && row.NeededMilli.Int64 == need {
				continue
			}
			before, err := dominio.RowJSON(ctx, tx, "shopping_list_item", row.ID)
			if err != nil {
				return err
			}
			row.NeededMilli = sql.NullInt64{Int64: need, Valid: true}
			if err := dominio.ShopA.SaveItem(ctx, tx, row); err != nil {
				return err
			}
			if err := audit(ctx, tx, r, "update", "shopping_list_item", row.ID, before); err != nil {
				return err
			}
		}
		for _, row := range rows { // calculated sin elección que el plan ya no pide: se quita (lo elegido se conserva, inv. 13)
			if row.Origin != recetas.KindCalculated || keep[row.ID] || hasChoice(row) {
				continue
			}
			before, err := dominio.RowJSON(ctx, tx, "shopping_list_item", row.ID)
			if err == nil {
				err = dominio.ShopA.Delete(ctx, tx, "shopping_list_item", row.ID)
			}
			if err == nil {
				err = audit(ctx, tx, r, "delete", "shopping_list_item", row.ID, before)
			}
			if err != nil {
				return err
			}
		}
		before, err := dominio.RowJSON(ctx, tx, "shopping_list", listID)
		if err != nil {
			return err
		}
		now := dominio.FormatDBTime(h.reg.Deps().Clock.Now())
		if err := dominio.ShopA.Update(ctx, tx, "shopping_list", listID, []string{"generated_at"}, now); err != nil {
			return err
		}
		return audit(ctx, tx, r, "update", "shopping_list", listID, before)
	})
	if err != nil {
		return dbErr(err, "La lista cambió al mismo tiempo; intenta de nuevo.")
	}
	b, err := h.buildList(ctx, h.db(), monday, http.StatusConflict)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, b.out)
}

// checkItem valida un artículo completo (alta o resultado del PATCH) contra las reglas de shopping_list_item.
func (h *handlers) checkItem(ctx context.Context, q dominio.Querier, it dominio.ShopListItem) error {
	var v api.Validator
	v.Check(it.ProductID.Valid != it.FreeText.Valid, "product_id", "required", "Elige un producto del catálogo o escribe el artículo, pero no los dos.")
	if it.FreeText.Valid {
		v.Check(strings.TrimSpace(it.FreeText.String) != "", "free_text", "required", "Escribe el artículo.")
		v.Check(it.FreeTextUnit.Valid && oneOf(it.FreeTextUnit.String, "g", "ml", "piece"), "free_text_unit", "required", "Indica la unidad del artículo escrito: g, ml o piece.")
		v.Check(!it.PresentationID.Valid, "presentation_id", "not_allowed", "Un artículo escrito a mano no lleva presentación.")
	} else {
		v.Check(!it.FreeTextUnit.Valid, "free_text_unit", "not_allowed", "Un producto del catálogo usa su propia unidad.")
	}
	v.Check(it.Origin != recetas.KindManual || it.QuantityMilli.Valid, "quantity_milli", "required", "En un artículo agregado a mano la cantidad es obligatoria.")
	v.Check(!it.QuantityMilli.Valid || it.QuantityMilli.Int64 > 0, "quantity_milli", "out_of_range", "La cantidad debe ser mayor que cero.")
	if err := h.mustExist(ctx, q, &v, "store", it.StoreID, "store_id", "No encontramos esa tienda."); err != nil {
		return err
	}
	if it.Origin == recetas.KindManual {
		if err := h.mustExist(ctx, q, &v, "product", it.ProductID, "product_id", "No encontramos ese producto."); err != nil {
			return err
		}
	}
	if it.PresentationID.Valid && it.ProductID.Valid {
		p, err := dominio.ShopA.Presentation(ctx, q, it.PresentationID.Int64)
		if err != nil && !isNoRows(err) {
			return err
		}
		v.Check(err == nil && p.ProductID == it.ProductID.Int64, "presentation_id", "not_found", "Esa presentación no es de este producto.")
	}
	return v.Err()
}

func (h *handlers) createShoppingListItem(w http.ResponseWriter, r *http.Request) error {
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	var in api.ShoppingListItemCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	it := dominio.ShopListItem{Origin: recetas.KindManual, ProductID: ni(in.ProductID), FreeText: ns(in.FreeText), FreeTextUnit: ns(in.FreeTextUnit),
		PresentationID: ni(in.PresentationID), QuantityMilli: sql.NullInt64{Int64: in.QuantityMilli, Valid: true}, StoreID: ni(in.StoreID), Notes: ns(in.Notes)}
	if it.FreeText.Valid {
		it.FreeText.String = strings.TrimSpace(it.FreeText.String)
	}
	if err := h.checkItem(ctx, h.db(), it); err != nil {
		return err
	}
	var id int64
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		listID, created, err := dominio.ShopA.EnsureList(ctx, tx, monday.String(), monday.AddDays(6).String())
		if err == nil && created {
			err = audit(ctx, tx, r, "insert", "shopping_list", listID, nil)
		}
		if err != nil {
			return err
		}
		id, err = dominio.ShopA.Insert(ctx, tx, "shopping_list_item", []string{"shopping_list_id", "origin", "product_id", "free_text", "free_text_unit",
			"presentation_id", "quantity_milli", "store_id", "notes"}, listID, it.Origin, it.ProductID, it.FreeText, it.FreeTextUnit, it.PresentationID,
			it.QuantityMilli, it.StoreID, it.Notes)
		if err != nil {
			return err
		}
		return audit(ctx, tx, r, "insert", "shopping_list_item", id, nil)
	})
	if err != nil {
		return dbErr(err, "Ese artículo ya está en la lista.")
	}
	b, err := h.buildList(ctx, h.db(), monday, http.StatusConflict)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, b.items[id])
}

func (h *handlers) updateShoppingListItem(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "shopping_list_item_id")
	if err != nil {
		return err
	}
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	var in api.ShoppingListItemUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if !in.PresentationID.Set && !in.QuantityMilli.Set && !in.StoreID.Set && !in.Notes.Set && in.IsBought == nil {
		return dominio.BadRequest("Manda al menos un cambio del artículo.")
	}
	ctx := r.Context()
	it, err := dominio.ShopA.Item(ctx, h.db(), id)
	if isNoRows(err) {
		return dominio.NotFound("No encontramos ese artículo de la lista.")
	} else if err != nil {
		return err
	}
	l, err := dominio.ShopA.ListByStart(ctx, h.db(), monday.String())
	if err != nil && !isNoRows(err) {
		return err
	}
	if err != nil || l.ID != it.ListID {
		return issue("week", "out_of_range", "Ese artículo no es de la lista de esa semana.")
	}
	it.PresentationID, it.QuantityMilli = patchInt(it.PresentationID, in.PresentationID), patchInt(it.QuantityMilli, in.QuantityMilli)
	it.StoreID, it.Notes = patchInt(it.StoreID, in.StoreID), patchStr(it.Notes, in.Notes)
	if in.IsBought != nil && *in.IsBought != it.IsBought {
		it.IsBought, it.BoughtAt = *in.IsBought, sql.NullString{}
		if it.IsBought {
			it.BoughtAt = sql.NullString{String: dominio.FormatDBTime(h.reg.Deps().Clock.Now()), Valid: true}
		}
	}
	if err := h.checkItem(ctx, h.db(), it); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "shopping_list_item", id)
		if err != nil {
			return err
		}
		if err := dominio.ShopA.SaveItem(ctx, tx, it); err != nil {
			return err
		}
		return audit(ctx, tx, r, "update", "shopping_list_item", id, before)
	})
	if err != nil {
		return dbErr(err, "Ese artículo ya está en la lista.")
	}
	b, err := h.buildList(ctx, h.db(), monday, http.StatusConflict)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, b.items[id])
}

func (h *handlers) assignShoppingListStore(w http.ResponseWriter, r *http.Request) error {
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	var in api.ShoppingListStoreAssignment
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if in.StoreID < 1 {
		return dominio.BadRequest("store_id debe ser un identificador.")
	}
	ctx := r.Context()
	var v api.Validator
	if err := h.mustExist(ctx, h.db(), &v, "store", sql.NullInt64{Int64: in.StoreID, Valid: true}, "store_id", "No encontramos esa tienda."); err != nil {
		return err
	}
	if err := v.Err(); err != nil {
		return err
	}
	l, err := dominio.ShopA.ListByStart(ctx, h.db(), monday.String())
	if isNoRows(err) {
		return dominio.NotFound("Esa semana todavía no tiene lista de compras.")
	} else if err != nil {
		return err
	}
	wd, err := h.loadWorld(ctx, h.db())
	if err != nil {
		return err
	}
	rows, err := dominio.ShopA.Items(ctx, h.db(), l.ID)
	if err != nil {
		return err
	}
	sc := code(in.StoreID)
	_, inBook := wd.pstore[sc]
	var res api.ShoppingListStoreAssignmentResult
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		for _, row := range rows {
			has := false
			if inBook && row.ProductID.Valid {
				_, ok, err := wd.book.QuoteFor(code(row.ProductID.Int64), sc, "", wd.now)
				if err != nil {
					return err
				}
				has = ok
			}
			if !has {
				res.WithoutPriceCount++
				continue
			}
			res.AssignedCount++
			if row.StoreID.Valid && row.StoreID.Int64 == in.StoreID {
				continue
			}
			before, err := dominio.RowJSON(ctx, tx, "shopping_list_item", row.ID)
			if err != nil {
				return err
			}
			row.StoreID = sql.NullInt64{Int64: in.StoreID, Valid: true}
			if err := dominio.ShopA.SaveItem(ctx, tx, row); err != nil {
				return err
			}
			if err := audit(ctx, tx, r, "update", "shopping_list_item", row.ID, before); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return dbErr(err, "La lista cambió al mismo tiempo; intenta de nuevo.")
	}
	return api.WriteJSON(w, http.StatusOK, res)
}

func (h *handlers) listBaskets(w http.ResponseWriter, r *http.Request) error {
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	p, err := api.ParsePage(r, []string{"expense_date"}, "-expense_date")
	if err != nil {
		return err
	}
	ctx := r.Context()
	b, err := h.buildList(ctx, h.db(), monday, http.StatusUnprocessableEntity)
	if err != nil {
		return err
	}
	all := []api.Basket{}
	if b.list != nil {
		all = append(all, api.Basket{Basket: "shopping_list", ItemCount: int64(len(b.rows)), TotalCents: api.Ptr(b.out.TotalCents)})
	}
	exps, err := dominio.ShopA.ExpenseBaskets(ctx, h.db(), p.OrderBy(map[string]string{"expense_date": "e.expense_date", "id": "e.id"}))
	if err != nil {
		return err
	}
	for _, e := range exps {
		all = append(all, api.Basket{Basket: "expense", ExpenseID: api.Ptr(e.ID), ExpenseDate: api.Ptr(e.Date), StoreID: pi(e.StoreID),
			TotalCents: api.Ptr(e.AmountCents), ItemCount: e.ItemCount})
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(pageSlice(all, p), p))
}

func storesOut(sts []precios.StoreTotals) []api.BasketComparisonStoresItem {
	out := []api.BasketComparisonStoresItem{}
	for _, s := range sts {
		out = append(out, api.BasketComparisonStoresItem{StoreID: idOf(s.Store), TotalCents: s.TotalCents, IdenticalCount: int64(s.IdenticalCount),
			IdenticalTotalCents: s.IdenticalTotalCents, SubstituteCount: int64(s.SubstituteCount), MissingCount: int64(s.MissingCount),
			StaleCount: int64(s.StaleCount), DifferenceVsCheapestCents: s.DifferenceVsCheapestCents, AssumesMembership: s.AssumesMembership})
	}
	return out
}

// priceLine es el precio de un artículo en una tienda con la fuente y la frescura del precio vigente de S5.
func (w *world) priceLine(store, pres string, price, base int64, observed string) api.BasketComparisonItemsItemPricesItem {
	out := api.BasketComparisonItemsItemPricesItem{StoreID: idOf(store), PresentationID: idOf(pres), PriceCents: price, BaseUnitPriceCents: base,
		ObservedDate: observed, Source: precios.SourceManual}
	if cur := w.book.CurrentPrice(pres, store, w.now); cur.Price != nil {
		out.Source, out.IsStale = cur.Price.Source, cur.Weight.IsStale
	}
	return out
}

func (h *handlers) getBasketComparison(w http.ResponseWriter, r *http.Request) error {
	qv := r.URL.Query()
	basket := qv.Get("basket")
	if basket == "" {
		basket = "shopping_list"
	}
	if !oneOf(basket, "shopping_list", "expense") {
		return dominio.BadRequest("basket debe ser shopping_list o expense.")
	}
	var expID int64
	if s := qv.Get("expense_id"); s != "" {
		ids, err := api.QueryIDs(r, "expense_id")
		if err != nil || len(ids) != 1 {
			return dominio.BadRequest("expense_id debe ser un identificador.")
		}
		expID = ids[0]
	}
	var v api.Validator
	v.Check(qv.Get("week") != "", "week", "required", "Falta la semana (week), por ejemplo 2026-W42.")
	v.Check(basket != "expense" || expID > 0, "expense_id", "required", "Elige la compra que quieres comparar (expense_id).")
	if err := v.Err(); err != nil {
		return err
	}
	monday, err := weekParam(r)
	if err != nil {
		return err
	}
	if basket == "expense" {
		return h.compareExpense(w, r, expID)
	}
	ctx := r.Context()
	b, err := h.buildList(ctx, h.db(), monday, http.StatusUnprocessableEntity)
	if err != nil {
		return err
	}
	wd := b.world
	list := comparableList(b)
	cmp, err := wd.book.CompareBasket(list)
	if err != nil {
		return engineErr(http.StatusUnprocessableEntity, err)
	}
	out := api.BasketComparison{Basket: "shopping_list", Stores: storesOut(cmp.Stores), CheapestStoreID: idPtr(cmp.CheapestStore),
		Items: []api.BasketComparisonItemsItem{}, Winners: []api.BasketComparisonWinnersItem{}}
	for _, bi := range cmp.Items {
		it := api.BasketComparisonItemsItem{ProductID: idOf(bi.Product), QuantityMilli: bi.QuantityMilli, Prices: []api.BasketComparisonItemsItemPricesItem{}}
		for _, ip := range bi.Prices {
			it.Prices = append(it.Prices, wd.priceLine(ip.Store, ip.Presentation, ip.PriceCents, ip.BaseUnitPriceCents, ip.ObservedDate))
			it.IsSubstitute = it.IsSubstitute || ip.IsSubstitute
		}
		out.Items = append(out.Items, it)
	}
	if out.Membership, err = h.membership(ctx, wd, list); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// comparableList es la canasta de S5 (§7.3) de la lista armada: los artículos con producto y cantidad por comprar.
// La usan getBasketComparison y BestBasket (tarjeta «Mejor canasta» de Inicio) para que den las mismas cifras.
func comparableList(b *built) []precios.ListItem {
	var list []precios.ListItem
	for _, g := range b.out.Groups {
		for _, it := range g.Items {
			if it.ProductID == nil {
				continue
			}
			li := precios.ListItem{Product: code(*it.ProductID), NeededMilli: it.NeededMilli, QuantityMilli: it.QuantityMilli}
			if pr := it.PresentationID; pr != nil {
				li.RecommendedPresentation = code(*pr)
			} else if pr := it.RecommendedPresentationID; pr != nil {
				li.RecommendedPresentation = code(*pr)
			}
			if buy, err := li.BuyMilli(); err == nil && buy > 0 {
				list = append(list, li)
			}
		}
	}
	return list
}

// expenseCart es la compra base de S5 (§7.2): los renglones con presentación del catálogo, a su precio unitario.
// La usan compareExpense y BestBasket.
func expenseCart(e dominio.ShopExpense, lines []dominio.ShopExpenseLine) precios.Cart {
	sc := codeOf(e.StoreID)
	obs, _ := precios.ParseDate(e.Date)
	cart := precios.Cart{Code: sc, Store: sc, Kind: "expense", ObservedDate: obs, Source: "ticket_manual"}
	for _, l := range lines {
		if !l.PresentationID.Valid {
			continue
		}
		unit := l.UnitPriceCents.Int64
		if !l.UnitPriceCents.Valid { // PROVISIONAL: precio unitario = total del renglón ÷ cantidad, ROUND_HALF_UP
			unit = dominio.MulDivRound(l.LineTotalCents, 1000, l.QuantityMilli)
		}
		cart.Lines = append(cart.Lines, precios.CartLine{Presentation: code(l.PresentationID.Int64), QuantityMilli: l.QuantityMilli, UnitPriceCents: unit})
	}
	return cart
}

// firstComplete es la tienda más barata de S5 §7.4 en el modo compra: la primera completa del orden del motor
// (nil si ninguna lo está). La usan compareExpense y BestBasket.
func firstComplete(sts []precios.StoreTotals) *string {
	for i := range sts {
		if sts[i].MissingCount == 0 {
			return &sts[i].Store
		}
	}
	return nil
}

// membership: N-164 con S5 Membership para la primera tienda que pide membresía (historial = compras con renglones).
func (h *handlers) membership(ctx context.Context, wd *world, list []precios.ListItem) (*api.BasketComparisonMembership, error) {
	if len(list) == 0 {
		return nil, nil
	}
	exps, err := dominio.ShopA.ExpenseBaskets(ctx, h.db(), "ORDER BY e.expense_date, e.id")
	if err != nil {
		return nil, err
	}
	var hist []precios.Purchase
	for _, e := range exps {
		if _, ok := wd.pstore[codeOf(e.StoreID)]; ok {
			d, _ := precios.ParseDate(e.Date)
			hist = append(hist, precios.Purchase{Date: d, Store: codeOf(e.StoreID)})
		}
	}
	for _, st := range wd.book.Stores {
		if !st.RequiresMembership {
			continue
		}
		m, err := wd.book.Membership(st.Code, nil, list, hist)
		if err != nil {
			return nil, engineErr(http.StatusUnprocessableEntity, err)
		}
		return &api.BasketComparisonMembership{StoreID: idOf(m.Store), BasketSavingsCents: m.BasketSavingsCents, AnnualSavingsCents: api.Ptr(m.AnnualSavingsCents),
			AnnualFeeCents: m.AnnualFeeCents, BreakEvenPurchases: m.BreakEvenPurchases, Verdict: m.Verdict}, nil
	}
	return nil, nil
}

func (h *handlers) compareExpense(w http.ResponseWriter, r *http.Request, expID int64) error {
	ctx := r.Context()
	e, lines, err := dominio.ShopA.ExpenseBasket(ctx, h.db(), expID)
	if isNoRows(err) {
		return issue("expense_id", "not_found", "No encontramos esa compra o no tiene renglones.")
	} else if err != nil {
		return err
	}
	wd, err := h.loadWorld(ctx, h.db())
	if err != nil {
		return err
	}
	sc := codeOf(e.StoreID)
	if _, ok := wd.pstore[sc]; !ok {
		return issue("expense_id", "not_allowed", "La compra no tiene una tienda activa: no se puede comparar.")
	}
	cart := expenseCart(e, lines)
	if len(cart.Lines) == 0 {
		return issue("expense_id", "not_allowed", "La compra no tiene renglones con presentación del catálogo.")
	}
	cmp, err := wd.book.CompareExpense(cart)
	if err != nil {
		return engineErr(http.StatusUnprocessableEntity, err)
	}
	out := api.BasketComparison{Basket: "expense", ExpenseID: api.Ptr(expID), Stores: storesOut(cmp.Stores), Items: []api.BasketComparisonItemsItem{},
		Winners: []api.BasketComparisonWinnersItem{}, Ties: int64(cmp.Ties)}
	out.CheapestStoreID = idPtr(firstComplete(cmp.Stores)) // S5 §7.4: la primera tienda completa
	for _, c := range cmp.Winners {
		out.Winners = append(out.Winners, api.BasketComparisonWinnersItem{StoreID: idOf(c.Store), Wins: int64(c.Count)})
	}
	for _, row := range cmp.Rows {
		it := api.BasketComparisonItemsItem{ProductID: idOf(row.Product), IsSubstitute: row.Kind == "substitute", WinnerStoreID: idPtr(row.WinnerStore),
			Prices: []api.BasketComparisonItemsItemPricesItem{}}
		var lo, hi int64
		for i, l := range row.Lines {
			pl := wd.priceLine(l.Store, l.Presentation, l.TotalCents, l.BaseUnitPriceCents, wd.today.String())
			if l.Store == sc {
				pl.ObservedDate, pl.Source, pl.IsStale = e.Date, "ticket_manual", false
				it.QuantityMilli = l.QuantityMilli
			} else if cur := wd.book.CurrentPrice(l.Presentation, l.Store, wd.now); cur.Price != nil {
				pl.ObservedDate = cur.Price.ObservedDate.Format(time.DateOnly)
			}
			it.Prices = append(it.Prices, pl)
			if i == 0 || l.TotalCents < lo {
				lo = l.TotalCents
			}
			if i == 0 || l.TotalCents > hi {
				hi = l.TotalCents
			}
		}
		if len(row.Lines) > 1 {
			it.DifferenceCents = api.Ptr(hi - lo)
		}
		if it.QuantityMilli == 0 && len(row.Lines) > 0 {
			it.QuantityMilli = row.Lines[0].QuantityMilli
		}
		out.Items = append(out.Items, it)
	}
	sort.SliceStable(out.Items, func(i, j int) bool { // orden por diferencia descendente
		return deref(out.Items[i].DifferenceCents, 0) > deref(out.Items[j].DifferenceCents, 0)
	})
	return api.WriteJSON(w, http.StatusOK, out)
}
