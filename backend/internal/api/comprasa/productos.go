package comprasa

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

func presOut(p dominio.ShopPresentation) api.Presentation {
	return api.Presentation{ID: p.ID, ProductID: p.ProductID, Name: p.Name, Brand: ps(p.Brand), SizeMilli: p.SizeMilli, IsBulk: p.IsBulk,
		Barcode: ps(p.Barcode), IsActive: p.IsActive, Notes: ps(p.Notes), CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt)}
}

func productOut(p dominio.ShopProduct, w *world) api.Product {
	out := api.Product{ID: p.ID, Name: p.Name, ProductGroupID: pi(p.GroupID), BaseUnit: p.BaseUnit, GramsPerUnitMilli: pi(p.GramsPerUnitMilli),
		IsActive: p.IsActive, Notes: ps(p.Notes), CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt)}
	var fresh string
	out.BestPrice, out.LastObservedDate, fresh, _ = w.priceInfo(p)
	out.Freshness = api.Ptr(fresh)
	return out
}

func (h *handlers) productDetail(ctx context.Context, q dominio.Querier, id int64) (*api.ProductDetail, error) {
	p, err := dominio.ShopA.Product(ctx, q, id)
	if isNoRows(err) {
		return nil, dominio.NotFound("No encontramos ese producto.")
	} else if err != nil {
		return nil, err
	}
	w, err := h.loadWorld(ctx, q)
	if err != nil {
		return nil, err
	}
	b := productOut(p, w)
	out := &api.ProductDetail{ID: b.ID, Name: b.Name, ProductGroupID: b.ProductGroupID, BaseUnit: b.BaseUnit, GramsPerUnitMilli: b.GramsPerUnitMilli,
		IsActive: b.IsActive, Notes: b.Notes, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt, BestPrice: b.BestPrice,
		LastObservedDate: b.LastObservedDate, Freshness: b.Freshness, Presentations: []api.Presentation{}, NutritionFacts: []api.NutritionFact{}}
	pres, err := dominio.ShopA.Presentations(ctx, q, id)
	if err != nil {
		return nil, err
	}
	for _, x := range pres {
		out.Presentations = append(out.Presentations, presOut(x))
	}
	nf, err := dominio.ShopA.Nutrition(ctx, q, id)
	if err != nil {
		return nil, err
	}
	for _, n := range nf {
		out.NutritionFacts = append(out.NutritionFacts, api.NutritionFact{ID: n.ID, ProductID: n.ProductID, PresentationID: pi(n.PresentationID),
			KcalMilliPer100g: n.KcalMilliPer100g, ProteinMgPer100g: pi(n.Protein), FatMgPer100g: pi(n.Fat), CarbsMgPer100g: pi(n.Carbs),
			FiberMgPer100g: pi(n.Fib), Source: n.Source, SourceRef: ps(n.SourceRef), SourceLicense: ps(n.SourceLicense), IsPreferred: n.IsPreferred,
			CreatedAt: api.APITime(n.CreatedAt), UpdatedAt: api.APITime(n.UpdatedAt)})
	}
	return out, nil
}

func (h *handlers) checkProduct(ctx context.Context, q dominio.Querier, v *api.Validator, p dominio.ShopProduct) error {
	v.Check(strings.TrimSpace(p.Name) != "", "name", "required", "Escribe el nombre del producto.")
	v.Check(oneOf(p.BaseUnit, "g", "ml", "piece"), "base_unit", "invalid_format", "La unidad base debe ser g, ml o piece.")
	v.Check(p.BaseUnit != "g" || !p.GramsPerUnitMilli.Valid, "grams_per_unit_milli", "not_allowed", "Un producto en gramos no lleva gramos por unidad.")
	v.Check(!p.GramsPerUnitMilli.Valid || p.GramsPerUnitMilli.Int64 > 0, "grams_per_unit_milli", "out_of_range", "Los gramos por unidad deben ser mayores que cero.")
	return h.mustExist(ctx, q, v, "product_group", p.GroupID, "product_group_id", "No encontramos ese grupo de producto.")
}

func checkPresentation(v *api.Validator, p dominio.ShopPresentation, prefix string) {
	v.Check(strings.TrimSpace(p.Name) != "", prefix+"name", "required", "Escribe el nombre de la presentación.")
	v.Check(p.SizeMilli > 0, prefix+"size_milli", "out_of_range", "El contenido de la presentación debe ser mayor que cero.")
	ok := !p.Barcode.Valid || (len(p.Barcode.String) >= 8 && len(p.Barcode.String) <= 14 && strings.Trim(p.Barcode.String, "0123456789") == "")
	v.Check(ok, prefix+"barcode", "invalid_format", "El código de barras lleva de 8 a 14 dígitos.")
}

var presInsertCols = []string{"product_id", "name", "brand", "size_milli", "is_bulk", "barcode", "is_active", "notes"}

func presFrom(productID int64, in api.PresentationCreate) dominio.ShopPresentation {
	return dominio.ShopPresentation{ProductID: productID, Name: strings.TrimSpace(in.Name), Brand: ns(in.Brand), SizeMilli: in.SizeMilli,
		IsBulk: deref(in.IsBulk, false), Barcode: ns(in.Barcode), IsActive: deref(in.IsActive, true), Notes: ns(in.Notes)}
}

func insertPres(ctx context.Context, tx *sql.Tx, r *http.Request, p dominio.ShopPresentation) (int64, error) {
	id, err := dominio.ShopA.Insert(ctx, tx, "presentation", presInsertCols, p.ProductID, p.Name, p.Brand, p.SizeMilli, p.IsBulk, p.Barcode, p.IsActive, p.Notes)
	if err != nil {
		return 0, err
	}
	return id, audit(ctx, tx, r, "insert", "presentation", id, nil)
}

func (h *handlers) listProducts(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"name"}, "name")
	if err != nil {
		return err
	}
	qv := r.URL.Query()
	var groupID int64
	if s := qv.Get("product_group_id"); s != "" {
		if groupID, err = strconv.ParseInt(s, 10, 64); err != nil || groupID < 1 {
			return dominio.BadRequest("El filtro product_group_id debe ser un identificador.")
		}
	}
	stores, err := api.QueryIDs(r, "store_id")
	if err != nil {
		return err
	}
	fresh, source := qv.Get("freshness"), qv.Get("source")
	if fresh != "" && !oneOf(fresh, "fresh", "stale", "no_price") {
		return dominio.BadRequest("El filtro freshness debe ser fresh, stale o no_price.")
	}
	if source != "" && !oneOf(source, "ticket_manual", "ticket_ocr", "site_assisted", "profeco", "manual") {
		return dominio.BadRequest("El filtro source debe ser ticket_manual, ticket_ocr, site_assisted, profeco o manual.")
	}
	ctx := r.Context()
	wd, err := h.loadWorld(ctx, h.db())
	if err != nil {
		return err
	}
	all, err := dominio.ShopA.Products(ctx, h.db(), groupID, p.OrderBy(map[string]string{"name": "name COLLATE NOCASE"}))
	if err != nil {
		return err
	}
	inStore := map[int64]bool{}
	for _, x := range wd.prices {
		for _, s := range stores {
			if x.StoreID == s {
				inStore[wd.presByID[x.PresentationID].ProductID] = true
			}
		}
	}
	items := []api.Product{}
	for _, x := range all {
		if !matchQ(qv.Get("q"), x.Name) || (len(stores) > 0 && !inStore[x.ID]) {
			continue
		}
		out := productOut(x, wd)
		_, _, _, src := wd.priceInfo(x)
		if (fresh != "" && *out.Freshness != fresh) || (source != "" && src != source) {
			continue
		}
		items = append(items, out)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(pageSlice(items, p), p))
}

func (h *handlers) createProduct(w http.ResponseWriter, r *http.Request) error {
	var in api.ProductCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	p := dominio.ShopProduct{Name: strings.TrimSpace(in.Name), GroupID: ni(in.ProductGroupID), BaseUnit: in.BaseUnit,
		GramsPerUnitMilli: ni(in.GramsPerUnitMilli), IsActive: deref(in.IsActive, true), Notes: ns(in.Notes)}
	var v api.Validator
	if err := h.checkProduct(ctx, h.db(), &v, p); err != nil {
		return err
	}
	for i, pc := range in.Presentations {
		checkPresentation(&v, presFrom(0, pc), fmt.Sprintf("presentations[%d].", i))
	}
	if err := v.Err(); err != nil {
		return err
	}
	var id int64
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		var err error
		id, err = dominio.ShopA.Insert(ctx, tx, "product", []string{"name", "product_group_id", "base_unit", "grams_per_unit_milli", "is_active", "notes"},
			p.Name, p.GroupID, p.BaseUnit, p.GramsPerUnitMilli, p.IsActive, p.Notes)
		if err != nil {
			return err
		}
		if err := audit(ctx, tx, r, "insert", "product", id, nil); err != nil {
			return err
		}
		for _, pc := range in.Presentations {
			if _, err := insertPres(ctx, tx, r, presFrom(id, pc)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return dbErr(err, "Ya existe un producto con ese nombre, o dos presentaciones con el mismo nombre.")
	}
	out, err := h.productDetail(ctx, h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func (h *handlers) getProduct(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "product_id")
	if err != nil {
		return err
	}
	out, err := h.productDetail(r.Context(), h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) updateProduct(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "product_id")
	if err != nil {
		return err
	}
	var in api.ProductUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	p, err := dominio.ShopA.Product(ctx, h.db(), id)
	if isNoRows(err) {
		return dominio.NotFound("No encontramos ese producto.")
	} else if err != nil {
		return err
	}
	if in.Name != nil {
		p.Name = strings.TrimSpace(*in.Name)
	}
	p.BaseUnit, p.IsActive = deref(in.BaseUnit, p.BaseUnit), deref(in.IsActive, p.IsActive)
	p.GroupID, p.GramsPerUnitMilli, p.Notes = patchInt(p.GroupID, in.ProductGroupID), patchInt(p.GramsPerUnitMilli, in.GramsPerUnitMilli), patchStr(p.Notes, in.Notes)
	var v api.Validator
	if err := h.checkProduct(ctx, h.db(), &v, p); err != nil {
		return err
	}
	if err := v.Err(); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "product", id)
		if err != nil {
			return err
		}
		if err := dominio.ShopA.SaveProduct(ctx, tx, p); err != nil {
			return err
		}
		return audit(ctx, tx, r, "update", "product", id, before)
	})
	if err != nil {
		return dbErr(err, "Ya existe un producto con ese nombre.")
	}
	out, err := h.productDetail(ctx, h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) createPresentation(w http.ResponseWriter, r *http.Request) error {
	pid, err := api.PathID(r, "product_id")
	if err != nil {
		return err
	}
	var in api.PresentationCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	if ok, err := dominio.ShopA.Exists(ctx, h.db(), "product", pid); err != nil {
		return err
	} else if !ok {
		return dominio.NotFound("No encontramos ese producto.")
	}
	p := presFrom(pid, in)
	var v api.Validator
	if checkPresentation(&v, p, ""); v.Err() != nil {
		return v.Err()
	}
	var id int64
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		var err error
		id, err = insertPres(ctx, tx, r, p)
		return err
	})
	if err != nil {
		return dbErr(err, "Ese producto ya tiene una presentación con ese nombre.")
	}
	out, err := dominio.ShopA.Presentation(ctx, h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, presOut(out))
}

func (h *handlers) updatePresentation(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "presentation_id")
	if err != nil {
		return err
	}
	var in api.PresentationUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	p, err := dominio.ShopA.Presentation(ctx, h.db(), id)
	if isNoRows(err) {
		return dominio.NotFound("No encontramos esa presentación.")
	} else if err != nil {
		return err
	}
	if in.Name != nil {
		p.Name = strings.TrimSpace(*in.Name)
	}
	p.SizeMilli, p.IsBulk, p.IsActive = deref(in.SizeMilli, p.SizeMilli), deref(in.IsBulk, p.IsBulk), deref(in.IsActive, p.IsActive)
	p.Brand, p.Barcode, p.Notes = patchStr(p.Brand, in.Brand), patchStr(p.Barcode, in.Barcode), patchStr(p.Notes, in.Notes)
	var v api.Validator
	if checkPresentation(&v, p, ""); v.Err() != nil {
		return v.Err()
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "presentation", id)
		if err != nil {
			return err
		}
		if err := dominio.ShopA.SavePresentation(ctx, tx, p); err != nil {
			return err
		}
		return audit(ctx, tx, r, "update", "presentation", id, before)
	})
	if err != nil {
		return dbErr(err, "Ese producto ya tiene una presentación con ese nombre.")
	}
	out, err := dominio.ShopA.Presentation(ctx, h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, presOut(out))
}
