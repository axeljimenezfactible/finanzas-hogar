// Package comprasa implementa la parte a de «Compras y tiendas» de C1 (agente api-compras-a, CONVENCIONES-backend.md
// §14): lista de compras de la semana, canastas y comparativa (motores S4 recetas y S5 precios, solo llamadas),
// tiendas, grupos de producto, productos y presentaciones.
package comprasa

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
	"finanzas-hogar/internal/recetas"
)

// Register registra las 18 operaciones del área.
func Register(reg *api.Registry) {
	h := &handlers{reg: reg}
	for op, f := range map[string]api.HandlerFunc{
		"getShoppingList": h.getShoppingList, "generateShoppingList": h.generateShoppingList,
		"createShoppingListItem": h.createShoppingListItem, "updateShoppingListItem": h.updateShoppingListItem,
		"assignShoppingListStore": h.assignShoppingListStore, "listBaskets": h.listBaskets,
		"getBasketComparison": h.getBasketComparison, "listStores": h.listStores, "createStore": h.createStore,
		"updateStore": h.updateStore, "setStoreProductGroups": h.setStoreProductGroups,
		"listProductGroups": h.listProductGroups, "listProducts": h.listProducts, "createProduct": h.createProduct,
		"getProduct": h.getProduct, "updateProduct": h.updateProduct, "createPresentation": h.createPresentation,
		"updatePresentation": h.updatePresentation,
	} {
		reg.Handle(op, f)
	}
}

type handlers struct {
	reg *api.Registry
	day *dominio.Date // «hoy» fijo cuando se calcula sin Registry (BestBasket); nil = reloj de la API
}

func (h *handlers) db() *sql.DB { return h.reg.Deps().DB }

func (h *handlers) today() dominio.Date {
	if h.day != nil {
		return *h.day
	}
	return dominio.Today(h.reg.Deps().Clock)
}

// WeekParam lee la llave obligatoria week (R3-01) con la regla del área: falta → 422 required; otro formato que
// 2026-W42 → 400. Devuelve el lunes de esa semana. La usa la tarjeta «Mejor canasta» de Inicio (api-gastos-inicio).
func WeekParam(r *http.Request) (dominio.Date, error) { return weekParam(r) }

// BestBasket arma la tarjeta «Mejor canasta» de Inicio (N-024, M-06; bloqueo S6-gastos-inicio §15) con la misma regla y
// el mismo motor S5 que getBasketComparison, sin cifras nuevas: la canasta es la lista de la semana que empieza en
// monday; si esa lista no tiene artículos con producto y cantidad por comprar (PROVISIONAL, §15), es la última compra
// con renglones (basket = expense, regla de compareExpense). today es «hoy» para la frescura de los precios (K22).
func BestBasket(ctx context.Context, q dominio.Querier, today, monday dominio.Date) (api.DashboardBestBasket, error) {
	h := &handlers{day: &today}
	b, err := h.buildList(ctx, q, monday, http.StatusUnprocessableEntity)
	if err != nil {
		return api.DashboardBestBasket{}, err
	}
	out := api.DashboardBestBasket{Basket: "shopping_list", ItemCount: int64(len(b.rows))} // item_count como listBaskets
	if list := comparableList(b); len(list) > 0 {
		cmp, err := b.world.book.CompareBasket(list)
		if err != nil {
			return api.DashboardBestBasket{}, engineErr(http.StatusUnprocessableEntity, err)
		}
		return bestOf(out, cmp.Stores, cmp.CheapestStore), nil
	}
	exps, err := dominio.ShopA.ExpenseBaskets(ctx, q, "ORDER BY e.expense_date DESC, e.id DESC")
	if err != nil || len(exps) == 0 {
		return out, err // sin lista ni compras con renglones: tarjeta vacía (sin store_id)
	}
	out = api.DashboardBestBasket{Basket: "expense", ExpenseID: api.Ptr(exps[0].ID), ItemCount: exps[0].ItemCount}
	e, lines, err := dominio.ShopA.ExpenseBasket(ctx, q, exps[0].ID)
	if err != nil {
		return api.DashboardBestBasket{}, err
	}
	cart := expenseCart(e, lines)
	if _, ok := b.world.pstore[cart.Store]; !ok || len(cart.Lines) == 0 {
		return out, nil // la compra no se puede comparar (getBasketComparison daría 422): tarjeta vacía, PROVISIONAL §15
	}
	cmp, err := b.world.book.CompareExpense(cart)
	if errors.Is(err, precios.ErrInvalidArgument) {
		// p. ej. ninguna otra tienda tiene precios que compitan («sin canastas que comparar»): getBasketComparison da 422;
		// la tarjeta queda vacía («Aún no hay precios para comparar», U3 tarjeta 5). PROVISIONAL §15.
		return out, nil
	} else if err != nil {
		return api.DashboardBestBasket{}, engineErr(http.StatusUnprocessableEntity, err)
	}
	return bestOf(out, cmp.Stores, firstComplete(cmp.Stores)), nil
}

// bestOf copia de S5 la tienda ganadora (cheapest, §7.4) y la siguiente del orden del motor, con su
// difference_vs_cheapest_cents; una tienda incompleta no tiene diferencia (null en S5), así que entonces no hay segunda.
func bestOf(out api.DashboardBestBasket, sts []precios.StoreTotals, cheapest *string) api.DashboardBestBasket {
	for i, st := range sts {
		if cheapest == nil || st.Store != *cheapest {
			continue
		}
		out.StoreID, out.TotalCents = api.Ptr(idOf(st.Store)), api.Ptr(st.TotalCents)
		out.CoveredCount = int64(st.IdenticalCount + st.SubstituteCount) // PROVISIONAL §15: un sustituto cuenta «con precio»
		out.StalePriceCount, out.AssumesMembership = int64(st.StaleCount), st.AssumesMembership
		if i+1 < len(sts) && sts[i+1].DifferenceVsCheapestCents != nil {
			out.SecondStoreID = api.Ptr(idOf(sts[i+1].Store))
			out.DifferenceVsSecondCents = api.Ptr(*sts[i+1].DifferenceVsCheapestCents)
		}
		break
	}
	return out
}

func ni(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

func ns(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func pi(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return api.Ptr(n.Int64)
}

func ps(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	return api.Ptr(n.String)
}

func deref[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

func patchInt(cur sql.NullInt64, n api.Nullable[int64]) sql.NullInt64 {
	if !n.Set {
		return cur
	}
	return sql.NullInt64{Int64: n.Value, Valid: n.Valid}
}

func patchStr(cur sql.NullString, n api.Nullable[string]) sql.NullString {
	if !n.Set {
		return cur
	}
	return sql.NullString{String: n.Value, Valid: n.Valid}
}

var accents = strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
	"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u", "Ü", "u", "Ñ", "n")

// matchQ aplica el filtro q: texto libre sin distinguir mayúsculas ni acentos.
func matchQ(q, name string) bool {
	f := func(s string) string { return strings.ToLower(accents.Replace(strings.TrimSpace(s))) }
	return q == "" || strings.Contains(f(name), f(q))
}

// pageSlice corta en memoria una colección ya ordenada (pide uno de más, como LimitOffset).
func pageSlice[T any](all []T, p api.PageParams) []T {
	if p.Offset >= len(all) {
		return nil
	}
	return all[p.Offset:min(len(all), p.Offset+p.Limit+1)]
}

func oneOf(s string, opts ...string) bool {
	for _, o := range opts {
		if s == o {
			return true
		}
	}
	return false
}

func issue(field, code, msg string) error {
	return dominio.Invalid(dominio.FieldIssue{Field: field, Code: code, Message: msg})
}

// audit escribe audit_log con el renglón antes y después (después vacío en un borrado).
func audit(ctx context.Context, tx *sql.Tx, r *http.Request, action, table string, id int64, before json.RawMessage) error {
	var after json.RawMessage
	if action != "delete" {
		var err error
		if after, err = dominio.RowJSON(ctx, tx, table, id); err != nil {
			return err
		}
	}
	return dominio.Audit(ctx, tx, api.Actor(r), action, table, id, before, after, api.RequestID(r))
}

// dbErr traduce choques con la base: UNIQUE → 409 already_exists; CHECK o llave → 409 conflict.
func dbErr(err error, dup string) error {
	switch {
	case err == nil:
		return nil
	case dominio.IsUnique(err):
		return dominio.Conflict("already_exists", dup)
	case dominio.IsConstraint(err):
		return dominio.Conflict("conflict", "El cambio choca con otro registro o con una regla de los datos. Revisa los valores e intenta de nuevo.")
	}
	return err
}

// weekParam lee la llave week (R3-01): falta → 422; formato distinto de 2026-W42 → 400. Devuelve el lunes.
func weekParam(r *http.Request) (dominio.Date, error) {
	s := r.URL.Query().Get("week")
	if s == "" {
		return dominio.Date{}, issue("week", "required", "Falta la semana (week), por ejemplo 2026-W42.")
	}
	rd, err := recetas.ParseWeek(s)
	if err != nil {
		return dominio.Date{}, dominio.BadRequest("La semana (week) debe tener el formato 2026-W42.")
	}
	return dominio.ParseDate(rd.String())
}

// refErr: un id de otra tabla que no existe → 422 not_found en ese campo.
func (h *handlers) mustExist(ctx context.Context, q dominio.Querier, v *api.Validator, table string, id sql.NullInt64, field, msg string) error {
	if !id.Valid {
		return nil
	}
	ok, err := dominio.ShopA.Exists(ctx, q, table, id.Int64)
	if err != nil {
		return err
	}
	v.Check(ok, field, "not_found", msg)
	return nil
}

var errNoRows = sql.ErrNoRows

func isNoRows(err error) bool { return errors.Is(err, errNoRows) }
