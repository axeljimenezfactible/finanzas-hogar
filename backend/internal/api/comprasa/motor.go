package comprasa

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
	"finanzas-hogar/internal/recetas"
)

// world reúne lo que leen los motores S5 (precios.Book) y S4 (recetas.Engine). Los motores trabajan con códigos de
// texto: aquí el código de cada registro es su id en decimal (code/idOf).
type world struct {
	today    dominio.Date
	now      time.Time // «hoy» para S5 (S5 §14.1, America/Mexico_City)
	settings precios.Settings
	stores   []dominio.ShopStore
	store    map[int64]dominio.ShopStore
	pstore   map[string]precios.Store
	groups   []dominio.ShopGroup
	group    map[int64]dominio.ShopGroup
	products []dominio.ShopProduct
	product  map[int64]dominio.ShopProduct
	pres     []dominio.ShopPresentation
	presByID map[int64]dominio.ShopPresentation
	prices   []dominio.ShopPrice
	book     *precios.Book
	visit    map[int64]string // próxima visita de cada mercado (S4 next_visit), la llena engine
}

func code(id int64) string { return strconv.FormatInt(id, 10) }

func idOf(c string) int64 { n, _ := strconv.ParseInt(c, 10, 64); return n }

func codeOf(n sql.NullInt64) string {
	if !n.Valid {
		return ""
	}
	return code(n.Int64)
}

func idPtr(c *string) *int64 {
	if c == nil || *c == "" {
		return nil
	}
	return api.Ptr(idOf(*c))
}

func rdate(s string) recetas.Date { d, _ := recetas.ParseDate(s); return d }

func rdatePtr(n sql.NullString) *recetas.Date {
	if !n.Valid {
		return nil
	}
	d := rdate(n.String)
	return &d
}

// loadWorld lee catálogo, tiendas activas, precios y parámetros price.* (K44) y arma el libro de precios de S5.
func (h *handlers) loadWorld(ctx context.Context, q dominio.Querier) (*world, error) {
	w := &world{today: h.today(), store: map[int64]dominio.ShopStore{}, pstore: map[string]precios.Store{}, group: map[int64]dominio.ShopGroup{},
		product: map[int64]dominio.ShopProduct{}, presByID: map[int64]dominio.ShopPresentation{}, visit: map[int64]string{}}
	var err error
	if w.now, err = precios.ParseDate(w.today.String()); err != nil {
		return nil, err
	}
	stale, err := dominio.SettingInt(ctx, q, "price.stale_after_days")
	if err != nil {
		return nil, err
	}
	zero, err := dominio.SettingInt(ctx, q, "price.zero_weight_after_days")
	if err != nil {
		return nil, err
	}
	w.settings = precios.Settings{StaleAfterDays: int(stale), ZeroWeightAfterDays: int(zero)}
	if w.stores, err = dominio.ShopA.Stores(ctx, q, "", "ORDER BY id"); err != nil {
		return nil, err
	}
	if w.groups, err = dominio.ShopA.Groups(ctx, q, "ORDER BY sort_order, id"); err != nil {
		return nil, err
	}
	if w.products, err = dominio.ShopA.Products(ctx, q, 0, "ORDER BY id"); err != nil {
		return nil, err
	}
	if w.pres, err = dominio.ShopA.Presentations(ctx, q, 0); err != nil {
		return nil, err
	}
	if w.prices, err = dominio.ShopA.Prices(ctx, q); err != nil {
		return nil, err
	}
	var pg []precios.Group
	for _, g := range w.groups {
		w.group[g.ID] = g
		pg = append(pg, precios.Group{Code: code(g.ID), Name: g.Name, PurchaseChannel: g.PurchaseChannel})
	}
	var pp []precios.Product
	for _, p := range w.products {
		w.product[p.ID] = p
		pp = append(pp, precios.Product{Code: code(p.ID), Name: p.Name, BaseUnit: p.BaseUnit, Group: codeOf(p.GroupID)})
	}
	var ppr []precios.Presentation
	for _, p := range w.pres {
		w.presByID[p.ID] = p
		ppr = append(ppr, precios.Presentation{Code: code(p.ID), Product: code(p.ProductID), Name: p.Name, SizeMilli: p.SizeMilli})
	}
	cat, err := precios.NewCatalog(pg, pp, ppr)
	if err != nil {
		return nil, err
	}
	active := map[int64]bool{}
	for _, s := range w.stores {
		w.store[s.ID] = s
		active[s.ID] = s.IsActive
	}
	var pst []precios.Store
	for _, s := range w.stores {
		if !s.IsActive {
			continue // PROVISIONAL: una tienda dada de baja no compite (bloqueo S6-compras-a-motores)
		}
		st := precios.Store{Code: code(s.ID), Kind: s.Kind, Name: s.Name, ProfecoName: s.ProfecoName.String,
			Municipality: s.Municipality.String, VisitWeekday: int(s.VisitWeekday.Int64), RequiresMembership: s.RequiresMembership,
			HasMembership: s.HasMembership, AssistedReading: s.AssistedReadingEnabled, MembershipFeeCents: pi(s.MembershipAnnualFeeCents)}
		if s.ParentID.Valid && active[s.ParentID.Int64] {
			st.Parent = code(s.ParentID.Int64)
		}
		for _, g := range s.GroupIDs {
			st.Groups = append(st.Groups, code(g))
		}
		w.pstore[st.Code] = st
		pst = append(pst, st)
	}
	var obs []precios.PriceObs
	for _, p := range w.prices {
		if !active[p.StoreID] {
			continue
		}
		d, err := precios.ParseDate(p.ObservedDate)
		if err != nil {
			return nil, err
		}
		created, _ := time.Parse(time.RFC3339Nano, p.CreatedAt)
		obs = append(obs, precios.PriceObs{Presentation: code(p.PresentationID), Store: code(p.StoreID), PriceCents: p.PriceCents,
			ObservedDate: d, Source: p.Source, IsPromotional: p.IsPromotional, SourceRef: p.SourceRef.String, CreatedAt: created})
	}
	w.book, err = precios.NewBook(cat, pst, obs, w.settings, w.now)
	return w, err
}

func rquote(q precios.Quote) recetas.Quote {
	return recetas.Quote{Product: q.Product, Store: q.Store, Presentation: q.Presentation, PriceCents: q.PriceCents, SizeMilli: q.SizeMilli,
		ObservedDate: rdate(q.ObservedDate.Format(time.DateOnly)), IsStale: q.IsStale, Recommended: q.Recommended}
}

// quotes arma la entrada PriceSource de S4 (S4-interfaz §8.1): por producto la recomendada de S5 (QuoteFor sin tienda
// ni presentación) y la de cada elección distinta de la lista (QuoteFor con esa elección), sin repetir tienda y
// presentación (el motor S4 vigente las indexa por ese par).
func (w *world) quotes(choices []recetas.ListChoice) (recetas.QuoteList, error) {
	ql := recetas.QuoteList{}
	for _, p := range w.products {
		pc := code(p.ID)
		rec, ok, err := w.book.QuoteFor(pc, "", "", w.now)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		seen := map[string]bool{rec.Store + "|" + rec.Presentation: true}
		ql = append(ql, rquote(rec))
		for _, ch := range choices {
			if ch.Product != pc || (ch.Store == "" && ch.Presentation == "") {
				continue
			}
			if _, ok := w.pstore[ch.Store]; ch.Store != "" && !ok {
				continue
			}
			if pr, ok := w.presByID[idOf(ch.Presentation)]; ch.Presentation != "" && (!ok || pr.ProductID != p.ID) {
				continue
			}
			q, ok, err := w.book.QuoteFor(pc, ch.Store, ch.Presentation, w.now)
			if err != nil {
				return nil, err
			}
			if ok && !seen[q.Store+"|"+q.Presentation] {
				seen[q.Store+"|"+q.Presentation] = true
				ql = append(ql, rquote(q))
			}
		}
	}
	return ql, nil
}

// engine arma el motor S4 de la semana con plan, recetas, inventario, elecciones de la lista y cotizaciones de S5.
func (w *world) engine(ctx context.Context, q dominio.Querier, monday dominio.Date, choices []recetas.ListChoice) (*recetas.Engine, error) {
	var rg []recetas.Group
	for _, g := range w.groups {
		rg = append(rg, recetas.Group{Code: code(g.ID), Name: g.Name, Channel: recetas.Channel(g.PurchaseChannel), SortOrder: int(g.SortOrder)})
	}
	var rp []recetas.Product
	for _, p := range w.products {
		rp = append(rp, recetas.Product{ID: p.ID, Code: code(p.ID), Name: p.Name, BaseUnit: recetas.Unit(p.BaseUnit),
			GramsPerUnitMilli: pi(p.GramsPerUnitMilli), Group: codeOf(p.GroupID)})
	}
	var rpr []recetas.Presentation
	for _, p := range w.pres {
		rpr = append(rpr, recetas.Presentation{ID: p.ID, Code: code(p.ID), Product: code(p.ProductID), Name: p.Name, SizeMilli: p.SizeMilli})
	}
	cat, err := recetas.NewCatalog(rg, rp, rpr)
	if err != nil {
		return nil, err
	}
	e := &recetas.Engine{Today: rdate(w.today.String()), Catalog: cat, Choices: choices}
	recs, err := dominio.ShopA.Recipes(ctx, q)
	if err != nil {
		return nil, err
	}
	for _, r := range recs {
		rr := &recetas.Recipe{ID: r.ID, Code: code(r.ID), Name: r.Name, BaseServings: int(r.BaseServings), Origin: r.Origin, IsActive: r.IsActive}
		for _, g := range r.Ingredients {
			rr.Ingredients = append(rr.Ingredients, recetas.Ingredient{Product: code(g.ProductID), Presentation: codeOf(g.PresentationID),
				QuantityMilli: g.QuantityMilli, DisplayText: g.DisplayText.String, Optional: g.IsOptional, SortOrder: int(g.SortOrder)})
		}
		e.Recipes = append(e.Recipes, rr)
	}
	people, err := dominio.ShopA.Persons(ctx, q)
	if err != nil {
		return nil, err
	}
	for _, p := range people {
		// PROVISIONAL: el esquema 1.4 no guarda quién tiene plan; se toma «co-titulares sí, externa no» (bloqueo S6-compras-a-motores).
		e.Persons = append(e.Persons, recetas.Person{Code: code(p.ID), HasMealPlan: p.Participation != "external"})
	}
	plan, err := dominio.ShopA.PlanEntries(ctx, q, monday.String(), monday.AddDays(6).String())
	if err != nil {
		return nil, err
	}
	for _, p := range plan {
		e.Plan = append(e.Plan, recetas.PlanEntry{Date: rdate(p.Date), Slot: recetas.SlotCode(p.Slot), Person: code(p.PersonID),
			Recipe: codeOf(p.RecipeID), ServingsMilli: pi(p.ServingsMilli), EatingOut: p.IsEatingOut})
	}
	inv, err := dominio.ShopA.Inventory(ctx, q)
	if err != nil {
		return nil, err
	}
	for _, it := range inv {
		e.Inventory = append(e.Inventory, recetas.InventoryItem{Code: code(it.ID), Product: codeOf(it.ProductID), Recipe: codeOf(it.RecipeID),
			Mode: it.Mode, QuantityMilli: it.QuantityMilli, CountedDate: rdatePtr(it.CountedDate), RateQuantityMilli: pi(it.RateQuantityMilli),
			RatePeriodDays: pi(it.RatePeriodDays), ReorderLevelMilli: pi(it.ReorderLevel), ExpiresDate: rdatePtr(it.ExpiresDate)})
	}
	from := monday // Q2-42 A: la próxima visita se cuenta desde la fecha más tardía entre hoy y el lunes de la lista
	if w.today.After(monday) {
		from = w.today
	}
	fromT, _ := precios.ParseDate(from.String())
	for _, s := range w.stores {
		st, ok := w.pstore[code(s.ID)]
		if !ok {
			continue
		}
		rs := recetas.Store{Code: st.Code, Kind: st.Kind, Groups: st.Groups}
		if s.Kind == "market" {
			nv, err := precios.NextVisit(st, fromT)
			if err != nil {
				return nil, err
			}
			d := rdate(nv.Format(time.DateOnly))
			rs.NextVisit, w.visit[s.ID] = &d, d.String()
		}
		e.Stores = append(e.Stores, rs)
	}
	if e.Prices, err = w.quotes(choices); err != nil {
		return nil, err
	}
	return e, nil
}

// priceInfo da best_price (QuoteFor de S5), la fecha del último precio, su frescura (precios.Freshness, K22) y su fuente.
func (w *world) priceInfo(p dominio.ShopProduct) (best *api.PriceAt, last *string, fresh string, source string) {
	fresh = "no_price"
	if q, ok, err := w.book.QuoteFor(code(p.ID), "", "", w.now); err == nil && ok {
		base, _, _, err := precios.BaseUnitPrice(q.PriceCents, q.SizeMilli, p.BaseUnit)
		if err == nil {
			best = &api.PriceAt{StoreID: idOf(q.Store), PresentationID: idOf(q.Presentation), PriceCents: q.PriceCents, BaseUnitPriceCents: base,
				ObservedDate: q.ObservedDate.Format(time.DateOnly), Source: w.sourceOf(q.Store, q.Presentation, precios.SourceManual), IsStale: q.IsStale}
		}
	}
	var lp *dominio.ShopPrice
	for i := range w.prices {
		x := &w.prices[i]
		if w.presByID[x.PresentationID].ProductID != p.ID {
			continue
		}
		if lp == nil || x.ObservedDate > lp.ObservedDate || (x.ObservedDate == lp.ObservedDate && x.ID > lp.ID) {
			lp = x
		}
	}
	if lp == nil {
		return best, nil, fresh, ""
	}
	fresh = "fresh"
	if d, err := precios.ParseDate(lp.ObservedDate); err == nil {
		if wt, err := precios.Freshness(d, w.now, w.settings, lp.Source); err == nil && wt.IsStale {
			fresh = "stale"
		}
	}
	return best, api.Ptr(lp.ObservedDate), fresh, lp.Source
}

// sourceOf es la fuente del precio vigente de S5 en esa tienda y presentación (o def si no hay).
func (w *world) sourceOf(store, pres, def string) string {
	if cur := w.book.CurrentPrice(pres, store, w.now); cur.Price != nil {
		return cur.Price.Source
	}
	return def
}

func (w *world) channelOf(productID int64) string {
	if g, ok := w.group[w.product[productID].GroupID.Int64]; ok {
		return g.PurchaseChannel
	}
	return "chain"
}

func (w *world) categoryOrder(productID sql.NullInt64) (int64, string) {
	if !productID.Valid {
		return 1 << 40, ""
	}
	p := w.product[productID.Int64]
	return w.group[p.GroupID.Int64].SortOrder, p.Name
}
