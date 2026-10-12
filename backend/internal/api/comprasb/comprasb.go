// Package comprasb implementa las operaciones de api-compras-b (CONVENCIONES-backend §14): precios con frescura
// (K22, K40, K44), nutrición, inventario por unidad y por ritmo, ajustes de precios y recetario, PROFECO y los lotes
// del agente externo (K53). Los cálculos son de los motores precios, profeco y recetas: aquí solo se llaman.
package comprasb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
)

// Register registra las 23 operaciones del área.
func Register(reg *api.Registry) {
	h := &handlers{reg: reg}
	for op, f := range map[string]api.HandlerFunc{
		"listProductPrices": h.listProductPrices, "getPriceSummary": h.getPriceSummary, "createPrice": h.createPrice,
		"updatePrice": h.updatePrice, "deletePrice": h.deletePrice, "getPriceSources": h.getPriceSources,
		"getPriceSettings": h.getPriceSettings, "updatePriceSettings": h.updatePriceSettings, "importProfeco": h.importProfeco,
		"listNutritionFacts": h.listNutritionFacts, "createNutritionFact": h.createNutritionFact,
		"getRecipeSettings": h.getRecipeSettings, "reloadNutritionBase": h.reloadNutritionBase,
		"listInventoryItems": h.listInventoryItems, "createInventoryItem": h.createInventoryItem,
		"getInventorySummary": h.getInventorySummary, "listInventoryMovements": h.listInventoryMovements,
		"adjustInventoryItem": h.adjustInventoryItem, "updateInventoryItem": h.updateInventoryItem,
		"getExternalAgentReference": h.getExternalAgentReference, "importExternalCatalog": h.importExternalCatalog,
		"importExternalPrices": h.importExternalPrices, "importExternalRecipes": h.importExternalRecipes,
	} {
		reg.Handle(op, f)
	}
}

type handlers struct{ reg *api.Registry }

func (h *handlers) db() *sql.DB { return h.reg.Deps().DB }

func (h *handlers) today() dominio.Date { return dominio.Today(h.reg.Deps().Clock) }

// todayTime es «hoy» (America/Mexico_City) en la forma que piden los motores.
func (h *handlers) todayTime() time.Time {
	t, _ := precios.ParseDate(h.today().String())
	return t
}

// Claves de app_setting de la sección «Tiendas y precios» (K44).
var priceParamKeys = []string{"price.stale_after_days", "price.zero_weight_after_days", "price.profeco_import_max_gb"}

const assistedKey = "price.assisted_reading_enabled"

// priceSettings lee los parámetros de frescura vigentes (nunca constantes en código).
func priceSettings(ctx context.Context, q dominio.Querier) (precios.Settings, error) {
	st, err := dominio.SettingInt(ctx, q, "price.stale_after_days")
	if err != nil {
		return precios.Settings{}, err
	}
	zw, err := dominio.SettingInt(ctx, q, "price.zero_weight_after_days")
	if err != nil {
		return precios.Settings{}, err
	}
	ar, err := dominio.SettingBool(ctx, q, assistedKey)
	return precios.Settings{StaleAfterDays: int(st), ZeroWeightAfterDays: int(zw), AssistedReadingEnabled: ar}, err
}

// priceOut arma el Price de C1; precio por unidad base (RN-20, K7) y frescura (RN-12, K22, K40) los da el motor.
func priceOut(p dominio.CBPriceRow, s precios.Settings, today time.Time) api.Price {
	out := api.Price{ID: p.ID, PresentationID: p.PresentationID, StoreID: p.StoreID, PriceCents: p.PriceCents,
		ObservedDate: p.ObservedDate, Source: p.Source, IsPromotional: p.IsPromotional,
		CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt)}
	if p.SourceRef.Valid {
		out.SourceRef = api.Ptr(p.SourceRef.String)
	}
	if p.ExpenseItemID.Valid {
		out.ExpenseItemID = api.Ptr(p.ExpenseItemID.Int64)
	}
	if p.AttachmentID.Valid {
		out.AttachmentID = api.Ptr(p.AttachmentID.Int64)
	}
	if c, _, _, err := precios.BaseUnitPrice(p.PriceCents, p.SizeMilli, p.BaseUnit); err == nil {
		out.BaseUnitPriceCents = api.Ptr(c)
	}
	if obs, err := precios.ParseDate(p.ObservedDate); err == nil {
		if w, err := precios.Freshness(obs, today, s, p.Source); err == nil {
			out.AgeDays, out.IsStale = api.Ptr(int64(w.AgeDays)), api.Ptr(w.IsStale)
		}
	}
	return out
}

var priceSources = []string{precios.SourceTicketManual, precios.SourceTicketOCR, precios.SourceSiteAssisted, precios.SourceProfeco, precios.SourceManual}

func oneOf(v string, set ...string) bool {
	for _, s := range set {
		if v == s {
			return true
		}
	}
	return false
}

func placeholders(n int) string { return strings.TrimSuffix(strings.Repeat("?,", n), ",") }

func (h *handlers) listProductPrices(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	pid, err := api.PathID(r, "product_id")
	if err != nil {
		return err
	}
	p, err := api.ParsePage(r, []string{"observed_date", "price_cents", "created_at"}, "-observed_date")
	if err != nil {
		return err
	}
	if ok, err := dominio.CBExists(ctx, h.db(), "product", pid); err != nil {
		return err
	} else if !ok {
		return dominio.NotFound("No encontramos ese producto.")
	}
	where, args := "WHERE pr.product_id = ?", []any{pid}
	stores, err := api.QueryIDs(r, "store_id")
	if err != nil {
		return err
	}
	if len(stores) > 0 {
		where += " AND p.store_id IN (" + placeholders(len(stores)) + ")"
		for _, s := range stores {
			args = append(args, s)
		}
	}
	per, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	if per != nil {
		where += " AND p.observed_date BETWEEN ? AND ?"
		args = append(args, per.Start.String(), per.End.String())
	}
	if src := r.URL.Query().Get("source"); src != "" {
		if !oneOf(src, priceSources...) {
			return api.BadRequest("La fuente no es válida: usa ticket_manual, ticket_ocr, site_assisted, profeco o manual.")
		}
		where += " AND p.source = ?"
		args = append(args, src)
	}
	tail := p.OrderBy(map[string]string{"observed_date": "p.observed_date", "price_cents": "p.price_cents", "created_at": "p.created_at", "id": "p.id"}) + " " + p.LimitOffset()
	rows, err := dominio.CBQueryPrices(ctx, h.db(), where, args, tail)
	if err != nil {
		return err
	}
	s, err := priceSettings(ctx, h.db())
	if err != nil {
		return err
	}
	today := h.todayTime()
	items := make([]api.Price, 0, len(rows))
	for _, row := range rows {
		items = append(items, priceOut(row, s, today))
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

// getPriceSummary arma el libro del motor de precios con el catálogo activo y los precios que compiten, y usa
// Book.Summary (fresh, stale, no_price y fecha más reciente por tienda). Las fuentes que no compiten (PROFECO y,
// en el motor aprobado, site_assisted) no cambian el resultado del motor, así que no se cargan.
func (h *handlers) getPriceSummary(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	db := h.db()
	s, err := priceSettings(ctx, db)
	if err != nil {
		return err
	}
	var prods []precios.Product
	var pres []precios.Presentation
	var stores []precios.Store
	if err := eachRow(ctx, db, `SELECT id, name, base_unit FROM product WHERE is_active = 1 ORDER BY id`, func(sc scanner) error {
		var id int64
		var p precios.Product
		err := sc.Scan(&id, &p.Name, &p.BaseUnit)
		p.Code = code(id)
		prods = append(prods, p)
		return err
	}); err != nil {
		return err
	}
	if err := eachRow(ctx, db, `SELECT pr.id, pr.product_id, pr.name, pr.size_milli FROM presentation pr JOIN product pd ON pd.id = pr.product_id
 WHERE pd.is_active = 1 ORDER BY pr.id`, func(sc scanner) error {
		var id, pid int64
		var p precios.Presentation
		err := sc.Scan(&id, &pid, &p.Name, &p.SizeMilli)
		p.Code, p.Product = code(id), code(pid)
		pres = append(pres, p)
		return err
	}); err != nil {
		return err
	}
	var storeIDs []int64
	if err := eachRow(ctx, db, `SELECT id, kind, coalesce(parent_id, 0), name, coalesce(visit_weekday, 0) FROM store WHERE is_active = 1 ORDER BY id`, func(sc scanner) error {
		var id, parent int64
		var st precios.Store
		err := sc.Scan(&id, &st.Kind, &parent, &st.Name, &st.VisitWeekday)
		st.Code = code(id)
		if parent > 0 {
			st.Parent = code(parent)
		}
		stores, storeIDs = append(stores, st), append(storeIDs, id)
		return err
	}); err != nil {
		return err
	}
	cat, err := precios.NewCatalog(nil, prods, pres)
	if err != nil {
		return err
	}
	today := h.todayTime()
	book, err := precios.NewBook(cat, stores, nil, s, today)
	if err != nil {
		return err
	}
	if err := eachRow(ctx, db, `SELECT presentation_id, store_id, price_cents, observed_date, source, is_promotional, coalesce(source_ref, ''), created_at FROM price ORDER BY id`, func(sc scanner) error {
		var presID, storeID int64
		var o precios.PriceObs
		var obs, created string
		if err := sc.Scan(&presID, &storeID, &o.PriceCents, &obs, &o.Source, &o.IsPromotional, &o.SourceRef, &created); err != nil {
			return err
		}
		if c, _ := precios.Competes(o.Source); !c {
			return nil
		}
		o.Presentation, o.Store = code(presID), code(storeID)
		if _, ok := cat.Presentation(o.Presentation); !ok {
			return nil // producto inactivo
		}
		if _, ok := book.Store(o.Store); !ok {
			return nil // tienda inactiva
		}
		o.ObservedDate, _ = precios.ParseDate(obs)
		if t, err := dominio.ParseDBTime(created); err == nil {
			o.CreatedAt = t
		}
		if err := book.AddPrice(o); err != nil && !errors.Is(err, precios.ErrFutureDate) {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	ps, sl := book.Summary()
	out := api.PriceSummary{Stores: make([]api.PriceSummaryStoresItem, 0, len(sl))}
	for _, p := range ps {
		switch p.Status {
		case "fresh":
			out.FreshCount++
		case "stale":
			out.StaleCount++
		default:
			out.NoPriceCount++
		}
	}
	for i, st := range sl {
		it := api.PriceSummaryStoresItem{StoreID: storeIDs[i]}
		if st.LatestDate != nil {
			it.LatestObservedDate = api.Ptr(*st.LatestDate)
		}
		out.Stores = append(out.Stores, it)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func code(id int64) string { return strconv.FormatInt(id, 10) }

type scanner interface{ Scan(...any) error }

func eachRow(ctx context.Context, q dominio.Querier, query string, f func(scanner) error, args ...any) error {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := f(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

// checkPriceRefs valida tienda, presentación, importe y fecha (no futura: el motor de precios rechaza fechas futuras).
func (h *handlers) checkPriceRefs(ctx context.Context, q dominio.Querier, v *api.Validator, presID, storeID *int64, cents *int64, date *string) error {
	if presID != nil {
		if _, err := dominio.CBPresentationOf(ctx, q, *presID); errors.Is(err, sql.ErrNoRows) {
			v.Add("presentation_id", "not_found", "No existe esa presentación.")
		} else if err != nil {
			return err
		}
	}
	if storeID != nil {
		ok, err := dominio.CBExists(ctx, q, "store", *storeID)
		if err != nil {
			return err
		}
		v.Check(ok, "store_id", "not_found", "No existe esa tienda.")
	}
	if cents != nil {
		v.Check(*cents >= 0, "price_cents", "out_of_range", "El precio no puede ser negativo.")
	}
	if date != nil {
		d, err := v.Date(*date, "observed_date")
		if err != nil {
			return err
		}
		v.Check(!d.After(h.today()), "observed_date", "out_of_range", "La fecha del precio no puede ser futura.")
	}
	return nil
}

func (h *handlers) writePrice(w http.ResponseWriter, r *http.Request, status int, id int64) error {
	p, err := dominio.CBGetPrice(r.Context(), h.db(), id)
	if err != nil {
		return err
	}
	s, err := priceSettings(r.Context(), h.db())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, status, priceOut(p, s, h.todayTime()))
}

func duplicatePrice(err error) error {
	if dominio.IsUnique(err) {
		return dominio.Conflict("already_exists", "Ya hay un precio de esa presentación en esa tienda, fecha y fuente. Corrige el que ya existe.")
	}
	return err
}

func (h *handlers) createPrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var in api.PriceCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	var v api.Validator
	if err := h.checkPriceRefs(ctx, h.db(), &v, &in.PresentationID, &in.StoreID, &in.PriceCents, &in.ObservedDate); err != nil {
		return err
	}
	if err := v.Err(); err != nil {
		return err
	}
	row := dominio.CBPriceRow{PresentationID: in.PresentationID, StoreID: in.StoreID, PriceCents: in.PriceCents, ObservedDate: in.ObservedDate,
		Source: precios.SourceManual, IsPromotional: in.IsPromotional != nil && *in.IsPromotional}
	if in.SourceRef != nil {
		row.SourceRef = sql.NullString{String: *in.SourceRef, Valid: true}
	}
	var id int64
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		var err error
		if id, err = dominio.CBInsertPrice(ctx, tx, row); err != nil {
			return duplicatePrice(err)
		}
		return dominio.CBAuditWrite(ctx, tx, api.Actor(r), "insert", "price", id, nil, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	return h.writePrice(w, r, http.StatusCreated, id)
}

func readOnlyProfeco(p dominio.CBPriceRow) error {
	if p.Source == precios.SourceProfeco {
		return dominio.Forbidden("read_only_source", "Los precios de PROFECO son de referencia y no se editan ni se borran.")
	}
	return nil
}

func (h *handlers) updatePrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := api.PathID(r, "price_id")
	if err != nil {
		return err
	}
	var in api.PriceUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	cur, err := dominio.CBGetPrice(ctx, h.db(), id)
	if err != nil {
		return err
	}
	if err := readOnlyProfeco(cur); err != nil {
		return err
	}
	if in.StoreID == nil && in.PriceCents == nil && in.ObservedDate == nil {
		return api.BadRequest("Manda al menos un campo para corregir el precio.")
	}
	var v api.Validator
	if err := h.checkPriceRefs(ctx, h.db(), &v, nil, in.StoreID, in.PriceCents, in.ObservedDate); err != nil {
		return err
	}
	if err := v.Err(); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "price", id)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE price SET store_id = coalesce(?, store_id), price_cents = coalesce(?, price_cents),
 observed_date = coalesce(?, observed_date) WHERE id = ?`, in.StoreID, in.PriceCents, in.ObservedDate, id); err != nil {
			return duplicatePrice(err)
		}
		return dominio.CBAuditWrite(ctx, tx, api.Actor(r), "update", "price", id, before, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	return h.writePrice(w, r, http.StatusOK, id)
}

func (h *handlers) deletePrice(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := api.PathID(r, "price_id")
	if err != nil {
		return err
	}
	cur, err := dominio.CBGetPrice(ctx, h.db(), id)
	if err != nil {
		return err
	}
	if err := readOnlyProfeco(cur); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "price", id)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM price WHERE id = ?`, id); err != nil {
			return err
		}
		return dominio.Audit(ctx, tx, api.Actor(r), "delete", "price", id, before, nil, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	return api.NoContent(w)
}

// paramOut convierte un parámetro de app_setting (K44, C1 §11.7) al esquema SettingParameterInteger.
func paramOut(p dominio.Parameter) api.SettingParameterInteger {
	out := api.SettingParameterInteger{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, IsInitial: p.IsInitial,
		Value: api.Null[int64](), InitialValue: api.Null[int64]()}
	var n int64
	if json.Unmarshal(p.Value, &n) == nil && len(p.Value) > 0 && string(p.Value) != "null" {
		out.Value = api.Some(n)
	}
	if json.Unmarshal(p.InitialValue, &n) == nil && len(p.InitialValue) > 0 && string(p.InitialValue) != "null" {
		out.InitialValue = api.Some(n)
	}
	return out
}

func (h *handlers) priceSettingsOut(ctx context.Context, q dominio.Querier) (api.PriceSettings, error) {
	var out api.PriceSettings
	var err error
	if out.AssistedReadingEnabled, err = dominio.SettingBool(ctx, q, assistedKey); err != nil {
		return out, err
	}
	ps := make([]api.SettingParameterInteger, len(priceParamKeys))
	for i, k := range priceParamKeys {
		p, err := dominio.ReadParameter(ctx, q, k)
		if err != nil {
			return out, err
		}
		ps[i] = paramOut(p)
	}
	out.Parameters = api.PriceSettingsParameters{PriceStaleAfterDays: ps[0], PriceZeroWeightAfterDays: ps[1], PriceProfecoImportMaxGB: ps[2]}
	return out, nil
}

func (h *handlers) getPriceSettings(w http.ResponseWriter, r *http.Request) error {
	out, err := h.priceSettingsOut(r.Context(), h.db())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// priceSettingsPatch es PriceSettingsUpdate con los parámetros crudos, para que un null o un tipo incorrecto
// lleguen a ApplyParams (422) en lugar de perderse al decodificar.
type priceSettingsPatch struct {
	AssistedReadingEnabled *bool                      `json:"assisted_reading_enabled,omitzero"`
	Parameters             map[string]json.RawMessage `json:"parameters,omitzero"`
	ResetParameters        []string                   `json:"reset_parameters,omitzero"`
}

func (h *handlers) updatePriceSettings(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var in priceSettingsPatch
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if in.AssistedReadingEnabled == nil && in.Parameters == nil && in.ResetParameters == nil {
		return api.BadRequest("Manda al menos un ajuste para cambiar.")
	}
	if in.AssistedReadingEnabled != nil && *in.AssistedReadingEnabled {
		return dominio.Conflict("assisted_reading_not_allowed", "La aplicación ya no lee sitios de tiendas (K53): los precios de sitios los mete el agente externo. El interruptor se queda apagado.")
	}
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if in.AssistedReadingEnabled != nil {
			if err := dominio.SetSetting(ctx, tx, assistedKey, json.RawMessage("false"), api.Actor(r), api.RequestID(r)); err != nil {
				return err
			}
		}
		if in.Parameters != nil || in.ResetParameters != nil {
			ch := dominio.ParamChange{Set: in.Parameters, Reset: in.ResetParameters}
			if err := dominio.ApplyParams(ctx, tx, ch, priceParamKeys, api.Actor(r), api.RequestID(r)); err != nil {
				return err
			}
		}
		s, err := priceSettings(ctx, tx)
		if err != nil {
			return err
		}
		if s.Validate() != nil { // regla entre claves de Q2-77: la valida el motor de precios
			return dominio.Invalid(dominio.FieldIssue{Field: "parameters.price.zero_weight_after_days", Code: "out_of_range",
				Message: "Los días para que un precio deje de pesar tienen que ser iguales o mayores que los días para que sea viejo."})
		}
		return nil
	})
	if err != nil {
		return err
	}
	return h.getPriceSettings(w, r)
}
