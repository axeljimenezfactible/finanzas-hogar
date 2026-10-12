package ajustes

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

var reCode = regexp.MustCompile(`^[a-z0-9_]{1,40}$`)

func categoryOut(c dominio.CategoryRow) api.Category {
	return api.Category{
		ID: c.ID, ParentID: pInt(c.ParentID), Code: pStr(c.Code), Name: c.Name, Kind: pStr(c.Kind),
		DefaultCostShareRuleID: pInt(c.DefaultCostShareRuleID), SortOrder: c.SortOrder, IsActive: c.IsActive,
		CreatedAt: api.APITime(c.CreatedAt), UpdatedAt: api.APITime(c.UpdatedAt),
		ExpenseCount: api.Ptr(c.ExpenseCount), RecentUsageCount: api.Ptr(c.RecentUsageCount),
	}
}

// listCategories: orden por defecto sort_order; N-050 pide sort=-recent_usage_count,sort_order&limit=6.
func listCategories(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"sort_order", "name", "recent_usage_count", "expense_count"}, "sort_order")
		if err != nil {
			return err
		}
		f := dominio.CategoryFilter{Status: r.URL.Query().Get("status"), Today: today(reg)}
		if f.Status == "" {
			f.Status = "active"
		}
		if !oneOf(f.Status, "active", "archived") {
			return badEnum("status")
		}
		db := reg.Deps().DB
		if q := fold(r.URL.Query().Get("q")); q != "" {
			all, err := dominio.ListCategoryRows(r.Context(), db, dominio.CategoryFilter{Today: f.Today}, "")
			if err != nil {
				return err
			}
			f.FilterIDs, f.IDs = true, []int64{}
			for _, c := range all {
				if strings.Contains(fold(c.Name), q) || (c.Code.Valid && strings.Contains(c.Code.String, q)) {
					f.IDs = append(f.IDs, c.ID)
				}
			}
		}
		order := p.OrderBy(map[string]string{"sort_order": "c.sort_order", "name": "c.name", "id": "c.id"}) + " " + p.LimitOffset()
		cs, err := dominio.ListCategoryRows(r.Context(), db, f, order)
		if err != nil {
			return err
		}
		items := make([]api.Category, 0, len(cs))
		for _, c := range cs {
			items = append(items, categoryOut(c))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// checkCategory valida las reglas de negocio comunes al alta y al cambio.
func checkCategory(r *http.Request, reg *api.Registry, c dominio.CategoryRow) error {
	var v api.Validator
	v.Required(strings.TrimSpace(c.Name) != "", "name")
	v.Check(!c.Code.Valid || reCode.MatchString(c.Code.String), "code", "invalid_format",
		"La clave va en minúsculas sin acentos, con dígitos o guion bajo (hasta 40).")
	if err := refCheck(r, reg, &v, "category", c.ParentID, "parent_id", "La categoría padre no existe."); err != nil {
		return err
	}
	if err := refCheck(r, reg, &v, "cost_share_rule", c.DefaultCostShareRuleID, "default_cost_share_rule_id", "La regla de reparto no existe."); err != nil {
		return err
	}
	if c.ID > 0 && c.ParentID.Valid {
		cyc, err := dominio.CategoryInSubtree(r.Context(), reg.Deps().DB, c.ID, c.ParentID.Int64)
		if err != nil {
			return err
		}
		v.Check(!cyc, "parent_id", "conflict", "Una categoría no puede quedar dentro de sí misma ni de una de sus subcategorías.")
	}
	return v.Err()
}

func categoryCols(c dominio.CategoryRow) map[string]any {
	set := map[string]any{"name": strings.TrimSpace(c.Name), "sort_order": c.SortOrder, "is_active": c.IsActive,
		"parent_id": nil, "code": nil, "kind": nil, "default_cost_share_rule_id": nil}
	if c.ParentID.Valid {
		set["parent_id"] = c.ParentID.Int64
	}
	if c.Code.Valid {
		set["code"] = c.Code.String
	}
	if c.Kind.Valid {
		set["kind"] = c.Kind.String
	}
	if c.DefaultCostShareRuleID.Valid {
		set["default_cost_share_rule_id"] = c.DefaultCostShareRuleID.Int64
	}
	return set
}

const dupCategory = "Ya existe una categoría con esa clave."

func createCategory(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.CategoryCreate
		if _, err := decode(reg, r, &in); err != nil {
			return err
		}
		if in.Kind != nil && !oneOf(*in.Kind, "fixed", "variable") {
			return badEnum("kind")
		}
		c := dominio.CategoryRow{ParentID: nInt(in.ParentID), Code: nStr(in.Code), Name: in.Name, Kind: nStr(in.Kind),
			DefaultCostShareRuleID: nInt(in.DefaultCostShareRuleID), IsActive: in.IsActive == nil || *in.IsActive}
		if err := checkCategory(r, reg, c); err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		if in.SortOrder != nil {
			c.SortOrder = *in.SortOrder
		} else {
			n, err := dominio.NextCategorySortOrder(ctx, db, in.ParentID)
			if err != nil {
				return err
			}
			c.SortOrder = n
		}
		var id int64
		err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			var err error
			id, err = dominio.InsertOwned(ctx, tx, "category", categoryCols(c), api.Actor(r), api.RequestID(r))
			return err
		})
		if err := writeErr(err, dupCategory); err != nil {
			return err
		}
		out, err := dominio.CategoryRowByID(ctx, db, id, today(reg))
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, categoryOut(out))
	}
}

func updateCategory(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "category_id")
		if err != nil {
			return err
		}
		var in api.CategoryUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return noChanges()
		}
		if in.Kind.Valid && !oneOf(in.Kind.Value, "fixed", "variable") {
			return badEnum("kind")
		}
		ctx, db := r.Context(), reg.Deps().DB
		c, err := dominio.CategoryRowByID(ctx, db, id, today(reg))
		if err != nil {
			return err
		}
		c.ParentID = mergeInt(c.ParentID, in.ParentID)
		c.Code = mergeStr(c.Code, in.Code)
		c.Kind = mergeStr(c.Kind, in.Kind)
		c.DefaultCostShareRuleID = mergeInt(c.DefaultCostShareRuleID, in.DefaultCostShareRuleID)
		if in.Name != nil {
			c.Name = *in.Name
		}
		if in.SortOrder != nil {
			c.SortOrder = *in.SortOrder
		}
		if in.IsActive != nil {
			c.IsActive = *in.IsActive
		}
		if err := checkCategory(r, reg, c); err != nil {
			return err
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			return dominio.UpdateOwned(ctx, tx, "category", id, categoryCols(c), api.Actor(r), api.RequestID(r))
		})
		if err := writeErr(err, dupCategory); err != nil {
			return err
		}
		out, err := dominio.CategoryRowByID(ctx, db, id, today(reg))
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, categoryOut(out))
	}
}

// reorderCategories: «lista completa de ids en el orden nuevo». PROVISIONAL (bloqueo S6-ajustes-api): deben venir
// todas las categorías activas, una vez cada una; las archivadas pueden venir o no (si no vienen, conservan su orden).
func reorderCategories(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.IdList
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if _, ok := keys["ids"]; !ok || in.IDs == nil {
			return dominio.Invalid(dominio.FieldIssue{Field: "ids", Code: "required", Message: "Manda la lista de categorías en el orden nuevo."})
		}
		ctx, db := r.Context(), reg.Deps().DB
		all, err := dominio.ListCategoryRows(ctx, db, dominio.CategoryFilter{Today: today(reg)}, "ORDER BY c.id")
		if err != nil {
			return err
		}
		byID := map[int64]dominio.CategoryRow{}
		for _, c := range all {
			byID[c.ID] = c
		}
		var v api.Validator
		seen := map[int64]bool{}
		for i, id := range in.IDs {
			f := "ids[" + strconv.Itoa(i) + "]"
			_, ok := byID[id]
			v.Check(ok, f, "not_found", "Esa categoría no existe.")
			v.Check(!seen[id], f, "conflict", "La categoría viene repetida.")
			seen[id] = true
		}
		for _, c := range all {
			v.Check(!c.IsActive || seen[c.ID], "ids", "required", "Falta la categoría «"+c.Name+"»: manda la lista completa.")
		}
		if err := v.Err(); err != nil {
			return err
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			for i, id := range in.IDs {
				if byID[id].SortOrder == int64(i+1) {
					continue
				}
				if err := dominio.UpdateOwned(ctx, tx, "category", id, map[string]any{"sort_order": i + 1}, api.Actor(r), api.RequestID(r)); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return api.NoContent(w)
	}
}

func personOut(p dominio.PersonRecord) api.Person {
	return api.Person{ID: p.ID, Code: p.Code, DisplayName: p.DisplayName, Participation: p.Participation,
		PaysShareDirectly: p.PaysShareDirectly, FortnightPaymentCapCents: pInt(p.FortnightPaymentCapCents),
		IsActive: p.IsActive, Notes: pStr(p.Notes), CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt),
		HasAccount: api.Ptr(p.HasAccount)}
}

// listPersons: orden holder, co_holder, external (C1).
func listPersons(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"participation", "display_name"}, "participation")
		if err != nil {
			return err
		}
		order := p.OrderBy(map[string]string{"participation": dominio.PersonOrder, "display_name": "p.display_name", "id": "p.id"})
		ps, err := dominio.ListPersonRecords(r.Context(), reg.Deps().DB, 0, order+" "+p.LimitOffset())
		if err != nil {
			return err
		}
		items := make([]api.Person, 0, len(ps))
		for _, x := range ps {
			items = append(items, personOut(x))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// updatePerson: nombre visible, tope por quincena (K13; null lo quita) y notas. No se crean personas ni cuentas.
func updatePerson(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "person_id")
		if err != nil {
			return err
		}
		var in api.PersonUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return noChanges()
		}
		ctx, db := r.Context(), reg.Deps().DB
		if _, err := dominio.PersonRecordByID(ctx, db, id); err != nil {
			return err
		}
		var v api.Validator
		set := map[string]any{}
		if in.DisplayName != nil {
			v.Required(strings.TrimSpace(*in.DisplayName) != "", "display_name")
			set["display_name"] = strings.TrimSpace(*in.DisplayName)
		}
		if in.FortnightPaymentCapCents.Set {
			set["fortnight_payment_cap_cents"] = nil
			if in.FortnightPaymentCapCents.Valid {
				v.Check(in.FortnightPaymentCapCents.Value >= 0, "fortnight_payment_cap_cents", "out_of_range", "El tope no puede ser negativo.")
				set["fortnight_payment_cap_cents"] = in.FortnightPaymentCapCents.Value
			}
		}
		if in.Notes.Set {
			set["notes"] = nil
			if in.Notes.Valid {
				set["notes"] = in.Notes.Value
			}
		}
		if err := v.Err(); err != nil {
			return err
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			return dominio.UpdateOwned(ctx, tx, "person", id, set, api.Actor(r), api.RequestID(r))
		})
		if err := writeErr(err, "Ya existe una persona con esos datos."); err != nil {
			return err
		}
		out, err := dominio.PersonRecordByID(ctx, db, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, personOut(out))
	}
}

func ruleOut(x dominio.ShareRuleRow) api.CostShareRule {
	out := api.CostShareRule{ID: x.ID, Name: x.Name, RemainderPersonID: x.RemainderPersonID, IsTemplate: x.IsTemplate,
		IsDefault: x.IsDefault, IsActive: x.IsActive, Notes: pStr(x.Notes), CreatedAt: api.APITime(x.CreatedAt),
		UpdatedAt: api.APITime(x.UpdatedAt), Parts: []api.CostShareRulePart{}}
	for _, p := range x.Parts {
		out.Parts = append(out.Parts, api.CostShareRulePart{ID: p.ID, RuleID: p.RuleID, PersonID: p.PersonID,
			PctBp: pInt(p.PctBp), FixedCents: pInt(p.FixedCents), CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt)})
	}
	return out
}

// listCostShareRules: is_template true por defecto (Hogar, Karen, Conjunto 50/50).
func listCostShareRules(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"name"}, "id")
		if err != nil {
			return err
		}
		tpl := true
		if s := r.URL.Query().Get("is_template"); s != "" {
			b, err := strconv.ParseBool(s)
			if err != nil || (s != "true" && s != "false") {
				return badEnum("is_template")
			}
			tpl = b
		}
		order := p.OrderBy(map[string]string{"name": "r.name", "id": "r.id"}) + " " + p.LimitOffset()
		rs, err := dominio.ListShareRuleRows(r.Context(), reg.Deps().DB, &tpl, 0, order)
		if err != nil {
			return err
		}
		items := make([]api.CostShareRule, 0, len(rs))
		for _, x := range rs {
			items = append(items, ruleOut(x))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// updateCostShareRule: nombre, notas y partes con su alcance (Q2-68 C, C1-datos §7.2). scope = all queda en 409
// PROVISIONAL (bloqueo S6-ajustes-api): falta la decisión C1-reparto-alcance y el motor no expone el reparto.
func updateCostShareRule(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "cost_share_rule_id")
		if err != nil {
			return err
		}
		var in api.CostShareRuleUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return noChanges()
		}
		if in.Scope != nil && !oneOf(*in.Scope, "onward", "all") {
			return badEnum("scope")
		}
		ctx, db := r.Context(), reg.Deps().DB
		if _, err := dominio.ShareRuleRowByID(ctx, db, id); err != nil {
			return err
		}
		_, hasParts := keys["parts"]
		var v api.Validator
		if in.Name != nil {
			v.Required(strings.TrimSpace(*in.Name) != "", "name")
		}
		v.Check(!hasParts || in.Scope != nil, "scope", "required", "Elige el alcance del cambio: «Solo los nuevos» o «También los ya capturados».")
		v.Check(hasParts || in.Scope == nil, "scope", "not_allowed", "El alcance solo va cuando cambian las partes.")
		parts := make([]dominio.ShareRulePartInput, 0, len(in.Parts))
		seen := map[int64]bool{}
		var pctSum int64
		for i, p := range in.Parts {
			f := "parts[" + strconv.Itoa(i) + "]"
			ok, err := dominio.RefExists(ctx, db, "person", p.PersonID)
			if err != nil {
				return err
			}
			v.Check(ok, f+".person_id", "not_found", "Esa persona no existe.")
			v.Check(!seen[p.PersonID], f+".person_id", "conflict", "La persona viene dos veces en las partes.")
			seen[p.PersonID] = true
			v.Check((p.PctBp == nil) != (p.FixedCents == nil), f, "invalid_format", "Cada parte lleva un porcentaje o un monto fijo, no los dos.")
			if p.PctBp != nil {
				v.Check(*p.PctBp >= 0 && *p.PctBp <= 10000, f+".pct_bp", "out_of_range", "El porcentaje va de 0 a 100 %.")
				pctSum += *p.PctBp
			}
			if p.FixedCents != nil {
				v.Check(*p.FixedCents >= 0, f+".fixed_cents", "out_of_range", "El monto fijo no puede ser negativo.")
			}
			parts = append(parts, dominio.ShareRulePartInput{PersonID: p.PersonID, PctBp: p.PctBp, FixedCents: p.FixedCents})
		}
		v.Check(pctSum <= 10000, "parts", "out_of_range", "Los porcentajes suman más de 100 %.")
		if err := v.Err(); err != nil {
			return err
		}
		if in.Scope != nil && *in.Scope == "all" {
			return dominio.Conflict("conflict", "Por ahora el cambio no se puede aplicar a los gastos ya capturados. Elige «Solo los nuevos».")
		}
		set := map[string]any{}
		if in.Name != nil {
			set["name"] = strings.TrimSpace(*in.Name)
		}
		if in.Notes.Set {
			set["notes"] = nil
			if in.Notes.Valid {
				set["notes"] = in.Notes.Value
			}
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			if err := dominio.UpdateOwned(ctx, tx, "cost_share_rule", id, set, api.Actor(r), api.RequestID(r)); err != nil {
				return err
			}
			if hasParts {
				return dominio.ReplaceShareRuleParts(ctx, tx, id, parts, api.Actor(r), api.RequestID(r))
			}
			return nil
		})
		if err := writeErr(err, "Ya existe una regla con esos datos."); err != nil {
			return err
		}
		out, err := dominio.ShareRuleRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, ruleOut(out))
	}
}
