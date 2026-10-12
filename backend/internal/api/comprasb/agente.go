package comprasb

// Agente externo (K53): la aplicación no lee sitios; una sesión de Claude Code de Axel mete productos, precios y
// recetas por lotes con la credencial de servicio. PROVISIONAL (bloqueo S5-agente-externo, pregunta e): audit_log
// queda sin usuario (api.Actor = nil, «proceso del sistema») y la Idempotency-Key no se guarda (CONVENCIONES §6);
// repetir un lote no duplica nada por las llaves naturales.

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
)

// Tamaños de lote de C1 1.4. PROVISIONAL (bloqueo S5-agente-externo, pregunta g).
const (
	maxCatalogBatch = 500
	maxPriceBatch   = 500
	maxRecipeBatch  = 100
)

func (h *handlers) getExternalAgentReference(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	out := api.ExternalAgentReference{Stores: []api.ExternalAgentReferenceStoresItem{}, ProductGroups: []api.ExternalAgentReferenceProductGroupsItem{},
		MealSlots: []api.ExternalAgentReferenceMealSlotsItem{}}
	if err := eachRow(ctx, h.db(), `SELECT id, kind, name, parent_id FROM store ORDER BY id`, func(sc scanner) error {
		var it api.ExternalAgentReferenceStoresItem
		var parent sql.NullInt64
		err := sc.Scan(&it.ID, &it.Kind, &it.Name, &parent)
		if parent.Valid {
			it.ParentID = api.Ptr(parent.Int64)
		}
		out.Stores = append(out.Stores, it)
		return err
	}); err != nil {
		return err
	}
	if err := eachRow(ctx, h.db(), `SELECT id, name, purchase_channel FROM product_group ORDER BY id`, func(sc scanner) error {
		var it api.ExternalAgentReferenceProductGroupsItem
		err := sc.Scan(&it.ID, &it.Name, &it.PurchaseChannel)
		out.ProductGroups = append(out.ProductGroups, it)
		return err
	}); err != nil {
		return err
	}
	if err := eachRow(ctx, h.db(), `SELECT id, code, name FROM meal_slot ORDER BY sort_order`, func(sc scanner) error {
		var it api.ExternalAgentReferenceMealSlotsItem
		err := sc.Scan(&it.ID, &it.Code, &it.Name)
		out.MealSlots = append(out.MealSlots, it)
		return err
	}); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func batchSize(n, max int) error {
	if n < 1 || n > max {
		return dominio.Invalid(dominio.FieldIssue{Field: "items", Code: "out_of_range", Message: "El lote debe traer de 1 a " + strconv.Itoa(max) + " renglones."})
	}
	return nil
}

// runBatch procesa cada renglón en su propia transacción: un renglón rechazado (error de negocio o regla de S1) no
// detiene a los demás. Otro error (base caída) corta el lote con 500.
func (h *handlers) runBatch(ctx context.Context, n int, row func(tx *sql.Tx, i int, it *api.ExternalImportResultItemsItem) error) (api.ExternalImportResult, error) {
	out := api.ExternalImportResult{Items: make([]api.ExternalImportResultItemsItem, 0, n)}
	for i := 0; i < n; i++ {
		it := api.ExternalImportResultItemsItem{Index: int64(i)}
		err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error { return row(tx, i, &it) })
		if err != nil {
			fe, ok := rowError(err)
			if !ok {
				return out, err
			}
			it = api.ExternalImportResultItemsItem{Index: int64(i), Status: "rejected", Error: &fe}
		}
		switch it.Status {
		case dominio.CBRowCreated:
			out.CreatedCount++
		case dominio.CBRowUpdated:
			out.UpdatedCount++
		case dominio.CBRowUnchanged:
			out.UnchangedCount++
		default:
			out.RejectedCount++
		}
		out.Items = append(out.Items, it)
	}
	return out, nil
}

func rowError(err error) (api.FieldError, bool) {
	if e, ok := dominio.AsError(err); ok {
		if len(e.Fields) > 0 {
			f := e.Fields[0]
			return api.FieldError{Field: f.Field, Code: f.Code, Message: f.Message}, true
		}
		return api.FieldError{Field: "items", Code: "conflict", Message: e.Message}, true
	}
	if dominio.IsConstraint(err) {
		return api.FieldError{Field: "items", Code: "invalid_format", Message: "El renglón no cumple las reglas del catálogo."}, true
	}
	return api.FieldError{}, false
}

var reBarcode = regexp.MustCompile(`^[0-9]{8,14}$`)

func (h *handlers) importExternalCatalog(w http.ResponseWriter, r *http.Request) error {
	var in api.ExternalCatalogImport
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if err := batchSize(len(in.Items), maxCatalogBatch); err != nil {
		return err
	}
	ctx := r.Context()
	out, err := h.runBatch(ctx, len(in.Items), func(tx *sql.Tx, i int, res *api.ExternalImportResultItemsItem) error {
		x := in.Items[i]
		var v api.Validator
		v.Check(strings.TrimSpace(x.Name) != "", "name", "required", "Falta el nombre del producto.")
		v.Check(oneOf(x.BaseUnit, "g", "ml", "piece"), "base_unit", "invalid_format", "La unidad base debe ser g, ml o piece.")
		v.Check(x.GramsPerUnitMilli == nil || (x.BaseUnit != "g" && *x.GramsPerUnitMilli > 0), "grams_per_unit_milli", "not_allowed", "Los gramos por unidad van solo en ml o piece y mayores que 0.")
		v.Check(len(x.Presentations) > 0, "presentations", "required", "El producto necesita al menos una presentación.")
		if x.ProductGroupID != nil {
			ok, err := dominio.CBExists(ctx, tx, "product_group", *x.ProductGroupID)
			if err != nil {
				return err
			}
			v.Check(ok, "product_group_id", "not_found", "No existe ese grupo de producto.")
		}
		item := dominio.CBCatalogItem{Name: strings.TrimSpace(x.Name), BaseUnit: x.BaseUnit, GramsPerUnitMilli: x.GramsPerUnitMilli, ProductGroupID: x.ProductGroupID}
		for j, p := range x.Presentations {
			f := "presentations[" + strconv.Itoa(j) + "]."
			v.Check(strings.TrimSpace(p.Name) != "", f+"name", "required", "Falta el nombre de la presentación.")
			v.Check(p.SizeMilli > 0, f+"size_milli", "out_of_range", "El contenido debe ser mayor que 0.")
			v.Check(p.Barcode == nil || reBarcode.MatchString(*p.Barcode), f+"barcode", "invalid_format", "El código de barras debe tener de 8 a 14 dígitos.")
			item.Presentations = append(item.Presentations, dominio.CBCatalogPresentation{Name: strings.TrimSpace(p.Name), SizeMilli: p.SizeMilli,
				Brand: p.Brand, IsBulk: p.IsBulk != nil && *p.IsBulk, Barcode: p.Barcode})
		}
		if err := v.Err(); err != nil {
			return err
		}
		st, id, pids, err := dominio.CBUpsertCatalogItem(ctx, tx, item, api.Actor(r), api.RequestID(r))
		if err != nil {
			return err
		}
		res.Status, res.ID, res.PresentationIDs = st, api.Ptr(id), pids
		return nil
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// importExternalPrices: source = site_assisted (PROVISIONAL, bloqueo S5-agente-externo, preguntas a y d).
func (h *handlers) importExternalPrices(w http.ResponseWriter, r *http.Request) error {
	var in api.ExternalPriceImport
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if err := batchSize(len(in.Items), maxPriceBatch); err != nil {
		return err
	}
	ctx := r.Context()
	out, err := h.runBatch(ctx, len(in.Items), func(tx *sql.Tx, i int, res *api.ExternalImportResultItemsItem) error {
		x := in.Items[i]
		var v api.Validator
		if _, err := dominio.ParseDate(x.ObservedDate); err != nil {
			return dominio.Invalid(dominio.FieldIssue{Field: "observed_date", Code: "invalid_format", Message: "La fecha debe ser AAAA-MM-DD y existir."})
		}
		if err := h.checkPriceRefs(ctx, tx, &v, &x.PresentationID, &x.StoreID, &x.PriceCents, &x.ObservedDate); err != nil {
			return err
		}
		if err := v.Err(); err != nil {
			return err
		}
		row := dominio.CBPriceRow{PresentationID: x.PresentationID, StoreID: x.StoreID, PriceCents: x.PriceCents, ObservedDate: x.ObservedDate,
			Source: precios.SourceSiteAssisted, IsPromotional: x.IsPromotional != nil && *x.IsPromotional}
		if x.ExternalRef != nil {
			row.SourceRef = sql.NullString{String: *x.ExternalRef, Valid: true}
		}
		st, id, err := dominio.CBUpsertSourcedPrice(ctx, tx, row, api.Actor(r), api.RequestID(r))
		if err != nil {
			return err
		}
		res.Status, res.ID = st, api.Ptr(id)
		return nil
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// importExternalRecipes: cada receta entra como borrador (PROVISIONAL, opción f2 del bloqueo S5-agente-externo).
func (h *handlers) importExternalRecipes(w http.ResponseWriter, r *http.Request) error {
	var in api.ExternalRecipeImport
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if err := batchSize(len(in.Items), maxRecipeBatch); err != nil {
		return err
	}
	ctx := r.Context()
	out, err := h.runBatch(ctx, len(in.Items), func(tx *sql.Tx, i int, res *api.ExternalImportResultItemsItem) error {
		x := in.Items[i]
		var v api.Validator
		v.Check(strings.TrimSpace(x.Name) != "", "name", "required", "Falta el nombre de la receta.")
		rec := dominio.CBDraftRecipe{Name: strings.TrimSpace(x.Name), Description: x.Description, Instructions: x.Instructions, Notes: x.Notes,
			PrepMinutes: x.PrepMinutes, CookMinutes: x.CookMinutes, MealSlotIDs: x.MealSlotIDs}
		if x.BaseServings != nil {
			rec.BaseServings = *x.BaseServings
		} else if err := tx.QueryRowContext(ctx, `SELECT CAST(dflt_value AS INTEGER) FROM pragma_table_info('recipe') WHERE name = 'base_servings'`).Scan(&rec.BaseServings); err != nil {
			return err
		}
		v.Check(rec.BaseServings >= 1 && rec.BaseServings <= 50, "base_servings", "out_of_range", "Las porciones base van de 1 a 50.")
		for _, m := range []*int64{x.PrepMinutes, x.CookMinutes} {
			v.Check(m == nil || *m >= 0, "prep_minutes", "out_of_range", "Los minutos no pueden ser negativos.")
		}
		for _, s := range x.MealSlotIDs {
			ok, err := dominio.CBExists(ctx, tx, "meal_slot", s)
			if err != nil {
				return err
			}
			v.Check(ok, "meal_slot_ids", "not_found", "No existe ese tiempo de comida.")
		}
		seen := map[int64]bool{}
		for j, g := range x.Ingredients {
			f := "ingredients[" + strconv.Itoa(j) + "]."
			ok, err := dominio.CBExists(ctx, tx, "product", g.ProductID)
			if err != nil {
				return err
			}
			v.Check(ok, f+"product_id", "not_found", "No existe ese producto; impórtalo primero con el catálogo.")
			v.Check(!seen[g.ProductID], f+"product_id", "conflict", "El producto se repite en la receta.")
			seen[g.ProductID] = true
			v.Check(g.QuantityMilli > 0, f+"quantity_milli", "out_of_range", "La cantidad debe ser mayor que 0.")
			if g.PresentationID != nil {
				owner, err := dominio.CBPresentationOf(ctx, tx, *g.PresentationID)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return err
				}
				v.Check(err == nil && owner == g.ProductID, f+"presentation_id", "not_found", "Esa presentación no es de ese producto.")
			}
			ing := dominio.CBDraftIngredient{ProductID: g.ProductID, QuantityMilli: g.QuantityMilli, PresentationID: g.PresentationID,
				DisplayText: g.DisplayText, IsOptional: g.IsOptional != nil && *g.IsOptional, SortOrder: int64(j)}
			if g.SortOrder != nil {
				ing.SortOrder = *g.SortOrder
			}
			rec.Ingredients = append(rec.Ingredients, ing)
		}
		if err := v.Err(); err != nil {
			return err
		}
		id, err := dominio.CBInsertDraftRecipe(ctx, tx, rec, api.Actor(r), api.RequestID(r))
		if err != nil {
			return err
		}
		res.Status, res.ID = dominio.CBRowCreated, api.Ptr(id)
		return nil
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}
