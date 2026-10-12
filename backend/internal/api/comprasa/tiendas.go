package comprasa

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

func storeOut(s dominio.ShopStore) api.Store {
	out := api.Store{ID: s.ID, Kind: s.Kind, ParentID: pi(s.ParentID), Name: s.Name, ProfecoName: ps(s.ProfecoName),
		WebsiteURL: ps(s.WebsiteURL), AssistedReadingEnabled: s.AssistedReadingEnabled, Address: ps(s.Address),
		Municipality: ps(s.Municipality), State: ps(s.State), VisitWeekday: pi(s.VisitWeekday),
		RequiresMembership: s.RequiresMembership, HasMembership: s.HasMembership,
		MembershipAnnualFeeCents: pi(s.MembershipAnnualFeeCents), IsActive: s.IsActive, Notes: ps(s.Notes),
		CreatedAt: api.APITime(s.CreatedAt), UpdatedAt: api.APITime(s.UpdatedAt), ProductGroupIDs: s.GroupIDs}
	if out.ProductGroupIDs == nil {
		out.ProductGroupIDs = []int64{}
	}
	return out
}

// checkStore valida la tienda completa (alta o resultado del PATCH) con las reglas de la tabla store.
func (h *handlers) checkStore(ctx context.Context, q dominio.Querier, s dominio.ShopStore) error {
	var v api.Validator
	v.Check(oneOf(s.Kind, "chain", "branch", "market"), "kind", "invalid_format", "El tipo debe ser chain (cadena), branch (sucursal) o market (mercado o tianguis).")
	v.Check(strings.TrimSpace(s.Name) != "", "name", "required", "Escribe el nombre de la tienda.")
	if s.Kind == "branch" {
		v.Check(s.ParentID.Valid, "parent_id", "required", "Una sucursal necesita su cadena (parent_id).")
		if s.ParentID.Valid {
			p, err := dominio.ShopA.Store(ctx, q, s.ParentID.Int64)
			switch {
			case isNoRows(err):
				v.Add("parent_id", "not_found", "No encontramos esa cadena.")
			case err != nil:
				return err
			default:
				v.Check(p.Kind == "chain" && p.ID != s.ID, "parent_id", "not_allowed", "La tienda padre de una sucursal debe ser una cadena.")
			}
		}
	} else {
		v.Check(!s.ParentID.Valid, "parent_id", "not_allowed", "Solo una sucursal lleva cadena (parent_id).")
	}
	if s.Kind == "market" {
		v.Check(s.VisitWeekday.Valid, "visit_weekday", "required", "Un mercado necesita su día de visita (1 = lunes … 7 = domingo).")
	} else {
		v.Check(!s.VisitWeekday.Valid, "visit_weekday", "not_allowed", "Solo un mercado o tianguis lleva día de visita.")
	}
	v.Check(!s.VisitWeekday.Valid || (s.VisitWeekday.Int64 >= 1 && s.VisitWeekday.Int64 <= 7), "visit_weekday", "out_of_range", "El día de visita va de 1 (lunes) a 7 (domingo).")
	if s.Kind != "chain" {
		v.Check(!s.RequiresMembership, "requires_membership", "not_allowed", "La membresía se captura en la cadena.")
		v.Check(!s.WebsiteURL.Valid, "website_url", "not_allowed", "El sitio se captura en la cadena.")
	}
	v.Check(s.RequiresMembership || !s.HasMembership, "has_membership", "not_allowed", "Primero marca que la tienda pide membresía (requires_membership).")
	v.Check(s.RequiresMembership || !s.MembershipAnnualFeeCents.Valid, "membership_annual_fee_cents", "not_allowed", "El costo de la membresía va solo en una tienda que la pide.")
	v.Check(!s.MembershipAnnualFeeCents.Valid || s.MembershipAnnualFeeCents.Int64 >= 0, "membership_annual_fee_cents", "out_of_range", "El costo de la membresía no puede ser negativo.")
	return v.Err()
}

func (h *handlers) listStores(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"name", "kind"}, "name")
	if err != nil {
		return err
	}
	kind := r.URL.Query().Get("kind")
	if kind != "" && !oneOf(kind, "chain", "branch", "market") {
		return dominio.BadRequest("El filtro kind debe ser chain, branch o market.")
	}
	all, err := dominio.ShopA.Stores(r.Context(), h.db(), kind, p.OrderBy(map[string]string{"name": "name COLLATE NOCASE"}))
	if err != nil {
		return err
	}
	items := []api.Store{}
	for _, s := range all {
		if matchQ(r.URL.Query().Get("q"), s.Name) {
			items = append(items, storeOut(s))
		}
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(pageSlice(items, p), p))
}

var storeInsertCols = []string{"kind", "parent_id", "name", "profeco_name", "website_url", "address", "municipality", "state",
	"visit_weekday", "requires_membership", "has_membership", "membership_annual_fee_cents", "is_active", "notes"}

func (h *handlers) createStore(w http.ResponseWriter, r *http.Request) error {
	var in api.StoreCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	s := dominio.ShopStore{Kind: in.Kind, Name: strings.TrimSpace(in.Name), ParentID: ni(in.ParentID), ProfecoName: ns(in.ProfecoName),
		WebsiteURL: ns(in.WebsiteURL), Address: ns(in.Address), Municipality: ns(in.Municipality), State: ns(in.State),
		VisitWeekday: ni(in.VisitWeekday), RequiresMembership: deref(in.RequiresMembership, false),
		HasMembership: deref(in.HasMembership, false), MembershipAnnualFeeCents: ni(in.MembershipAnnualFeeCents),
		IsActive: deref(in.IsActive, true), Notes: ns(in.Notes)}
	if err := h.checkStore(ctx, h.db(), s); err != nil {
		return err
	}
	var id int64
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		var err error
		id, err = dominio.ShopA.Insert(ctx, tx, "store", storeInsertCols, s.Kind, s.ParentID, s.Name, s.ProfecoName, s.WebsiteURL,
			s.Address, s.Municipality, s.State, s.VisitWeekday, s.RequiresMembership, s.HasMembership, s.MembershipAnnualFeeCents, s.IsActive, s.Notes)
		if err != nil {
			return err
		}
		return audit(ctx, tx, r, "insert", "store", id, nil)
	})
	if err != nil {
		return dbErr(err, "Ya existe esa tienda.")
	}
	out, err := dominio.ShopA.Store(ctx, h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, storeOut(out))
}

func (h *handlers) updateStore(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "store_id")
	if err != nil {
		return err
	}
	var in api.StoreUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	s, err := dominio.ShopA.Store(ctx, h.db(), id)
	if isNoRows(err) {
		return dominio.NotFound("No encontramos esa tienda.")
	} else if err != nil {
		return err
	}
	if in.AssistedReadingEnabled != nil && *in.AssistedReadingEnabled {
		return dominio.Conflict("assisted_reading_not_allowed", "La aplicación ya no lee sitios de tiendas (K53): la lectura asistida no se puede activar.")
	}
	if in.Kind != nil {
		s.Kind = *in.Kind
	}
	if in.Name != nil {
		s.Name = strings.TrimSpace(*in.Name)
	}
	s.ParentID, s.ProfecoName, s.WebsiteURL = patchInt(s.ParentID, in.ParentID), patchStr(s.ProfecoName, in.ProfecoName), patchStr(s.WebsiteURL, in.WebsiteURL)
	s.Address, s.Municipality, s.State = patchStr(s.Address, in.Address), patchStr(s.Municipality, in.Municipality), patchStr(s.State, in.State)
	s.VisitWeekday, s.Notes = patchInt(s.VisitWeekday, in.VisitWeekday), patchStr(s.Notes, in.Notes)
	s.MembershipAnnualFeeCents = patchInt(s.MembershipAnnualFeeCents, in.MembershipAnnualFeeCents)
	s.AssistedReadingEnabled = deref(in.AssistedReadingEnabled, s.AssistedReadingEnabled)
	s.RequiresMembership = deref(in.RequiresMembership, s.RequiresMembership)
	s.HasMembership, s.IsActive = deref(in.HasMembership, s.HasMembership), deref(in.IsActive, s.IsActive)
	if err := h.checkStore(ctx, h.db(), s); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "store", id)
		if err != nil {
			return err
		}
		if err := dominio.ShopA.SaveStore(ctx, tx, s); err != nil {
			return err
		}
		return audit(ctx, tx, r, "update", "store", id, before)
	})
	if err != nil {
		return dbErr(err, "Ya existe esa tienda.")
	}
	out, err := dominio.ShopA.Store(ctx, h.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, storeOut(out))
}

func (h *handlers) setStoreProductGroups(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "store_id")
	if err != nil {
		return err
	}
	var in api.IdList
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	s, err := dominio.ShopA.Store(ctx, h.db(), id)
	if isNoRows(err) {
		return dominio.NotFound("No encontramos esa tienda.")
	} else if err != nil {
		return err
	}
	var v api.Validator
	// ids es obligatorio en IdList (C1): {} o {"ids": null} → 422; para vaciar la tienda se manda {"ids": []}.
	v.Required(in.IDs != nil, "ids")
	seen := map[int64]bool{}
	for _, g := range in.IDs {
		v.Check(!seen[g], "ids", "conflict", "Un grupo de producto viene repetido.")
		seen[g] = true
		if err := h.mustExist(ctx, h.db(), &v, "product_group", sql.NullInt64{Int64: g, Valid: true}, "ids", "No encontramos uno de los grupos de producto."); err != nil {
			return err
		}
	}
	if err := v.Err(); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if err := dominio.ShopA.SetStoreGroups(ctx, tx, id, in.IDs); err != nil {
			return err
		}
		before, _ := json.Marshal(map[string]any{"product_group_ids": storeOut(s).ProductGroupIDs})
		after, _ := json.Marshal(map[string]any{"product_group_ids": in.IDs})
		return dominio.Audit(ctx, tx, api.Actor(r), "update", "store_product_group", id, before, after, api.RequestID(r))
	})
	if err != nil {
		return dbErr(err, "Un grupo de producto viene repetido.")
	}
	return api.NoContent(w)
}

func (h *handlers) listProductGroups(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"sort_order", "name"}, "sort_order")
	if err != nil {
		return err
	}
	gs, err := dominio.ShopA.Groups(r.Context(), h.db(), p.OrderBy(map[string]string{"name": "name COLLATE NOCASE"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.ProductGroup, 0, len(gs))
	for _, g := range gs {
		items = append(items, api.ProductGroup{ID: g.ID, Name: g.Name, PurchaseChannel: g.PurchaseChannel, SortOrder: g.SortOrder,
			IsActive: g.IsActive, CreatedAt: api.APITime(g.CreatedAt), UpdatedAt: api.APITime(g.UpdatedAt)})
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}
