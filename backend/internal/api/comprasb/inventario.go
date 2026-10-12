package comprasb

import (
	"context"
	"database/sql"
	"net/http"
	"sort"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/recetas"
)

// invStatus calcula existencia estimada, estado, consumo semanal y días restantes con el motor de recetas (§7, RN-15, K7).
func (h *handlers) invStatus(rows []dominio.CBInventoryRow) (map[int64]recetas.ItemStatus, recetas.InventorySummary, error) {
	today, err := recetas.ParseDate(h.today().String())
	if err != nil {
		return nil, recetas.InventorySummary{}, err
	}
	e := &recetas.Engine{Today: today}
	opt := func(v sql.NullInt64) *int64 {
		if v.Valid {
			return &v.Int64
		}
		return nil
	}
	for _, it := range rows {
		mi := recetas.InventoryItem{Code: code(it.ID), Mode: it.TrackingMode, QuantityMilli: it.QuantityMilli,
			RateQuantityMilli: opt(it.RateQuantityMilli), RatePeriodDays: opt(it.RatePeriodDays), ReorderLevelMilli: opt(it.ReorderLevelMilli)}
		if it.CountedDate.Valid {
			d, err := recetas.ParseDate(it.CountedDate.String)
			if err != nil {
				return nil, recetas.InventorySummary{}, err
			}
			mi.CountedDate = &d
		}
		e.Inventory = append(e.Inventory, mi)
	}
	rep, err := e.InventoryStatus()
	if err != nil {
		return nil, recetas.InventorySummary{}, err
	}
	out := make(map[int64]recetas.ItemStatus, len(rows))
	for i, s := range rep.Items {
		out[rows[i].ID] = s
	}
	return out, rep.Summary, nil
}

func invOut(it dominio.CBInventoryRow, s recetas.ItemStatus) api.InventoryItem {
	out := api.InventoryItem{ID: it.ID, TrackingMode: it.TrackingMode, QuantityMilli: it.QuantityMilli,
		CreatedAt: api.APITime(it.CreatedAt), UpdatedAt: api.APITime(it.UpdatedAt),
		Status: api.Ptr(s.Status), EstimatedQuantityMilli: api.Ptr(s.EstimatedQuantityMilli),
		WeeklyConsumptionMilli: s.WeeklyConsumptionMilli, DaysLeft: s.DaysLeft}
	i64 := func(v sql.NullInt64) *int64 {
		if v.Valid {
			return api.Ptr(v.Int64)
		}
		return nil
	}
	str := func(v sql.NullString) *string {
		if v.Valid {
			return api.Ptr(v.String)
		}
		return nil
	}
	out.ProductID, out.RecipeID, out.RateQuantityMilli, out.RatePeriodDays, out.ReorderLevelMilli =
		i64(it.ProductID), i64(it.RecipeID), i64(it.RateQuantityMilli), i64(it.RatePeriodDays), i64(it.ReorderLevelMilli)
	out.CountedDate, out.ExpiresDate, out.Notes = str(it.CountedDate), str(it.ExpiresDate), str(it.Notes)
	return out
}

func (h *handlers) writeInventory(w http.ResponseWriter, r *http.Request, status int, id int64) error {
	it, err := dominio.CBGetInventory(r.Context(), h.db(), id)
	if err != nil {
		return err
	}
	st, _, err := h.invStatus([]dominio.CBInventoryRow{it})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, status, invOut(it, st[it.ID]))
}

var statusRank = map[string]int{recetas.StatusOut: 0, recetas.StatusReorder: 1, recetas.StatusSufficient: 2}

// listInventoryItems: el estado lo calcula el motor, así que filtro por estado, orden y página se hacen sobre el
// resultado. Orden por defecto (C1): agotado, reponer, suficiente y luego nombre.
func (h *handlers) listInventoryItems(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"status"}, "status")
	if err != nil {
		return err
	}
	qs := r.URL.Query()
	where, args := "WHERE 1 = 1", []any{}
	if q := strings.TrimSpace(qs.Get("q")); q != "" {
		where += " AND coalesce(pd.name, r.name) LIKE ?"
		args = append(args, "%"+q+"%")
	}
	if m := qs.Get("tracking_mode"); m != "" {
		if !oneOf(m, recetas.ModeStock, recetas.ModeRate) {
			return api.BadRequest("El modo debe ser stock o consumption_rate.")
		}
		where += " AND i.tracking_mode = ?"
		args = append(args, m)
	}
	status := qs.Get("status")
	if status != "" && !oneOf(status, recetas.StatusSufficient, recetas.StatusReorder, recetas.StatusOut) {
		return api.BadRequest("El estado debe ser sufficient, reorder u out.")
	}
	rows, err := dominio.CBQueryInventory(r.Context(), h.db(), where, args...)
	if err != nil {
		return err
	}
	st, _, err := h.invStatus(rows)
	if err != nil {
		return err
	}
	items := make([]api.InventoryItem, 0, len(rows))
	names := map[int64]string{}
	for _, it := range rows {
		if status == "" || st[it.ID].Status == status {
			items = append(items, invOut(it, st[it.ID]))
			names[it.ID] = strings.ToLower(it.Name)
		}
	}
	desc := len(p.Sort) > 0 && p.Sort[0].Name == "status" && p.Sort[0].Desc
	sort.SliceStable(items, func(i, j int) bool {
		a, b := statusRank[*items[i].Status], statusRank[*items[j].Status]
		if a != b {
			return (a < b) != desc
		}
		if names[items[i].ID] != names[items[j].ID] {
			return names[items[i].ID] < names[items[j].ID]
		}
		return items[i].ID < items[j].ID
	})
	lo := min(p.Offset, len(items))
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items[lo:min(lo+p.Limit+1, len(items))], p))
}

func (h *handlers) getInventorySummary(w http.ResponseWriter, r *http.Request) error {
	rows, err := dominio.CBQueryInventory(r.Context(), h.db(), "")
	if err != nil {
		return err
	}
	_, sum, err := h.invStatus(rows)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, api.InventorySummary{SufficientCount: int64(sum.Sufficient), ReorderCount: int64(sum.Reorder), OutCount: int64(sum.Out)})
}

func checkInvConstraint(err error) error {
	if dominio.IsUnique(err) {
		return dominio.Conflict("already_exists", "Ese producto ya está en el inventario; ajusta su existencia en lugar de darlo de alta otra vez.")
	}
	if dominio.IsConstraint(err) {
		return dominio.Invalid(dominio.FieldIssue{Field: "tracking_mode", Code: "invalid_format",
			Message: "La combinación no es válida: el ritmo de consumo necesita cantidad, días y fecha de conteo; la existencia por unidad no lleva ritmo; un sobrante va por unidad."})
	}
	return err
}

func (h *handlers) createInventoryItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var in api.InventoryItemCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	var v api.Validator
	v.Check(oneOf(in.TrackingMode, recetas.ModeStock, recetas.ModeRate), "tracking_mode", "invalid_format", "El modo debe ser stock o consumption_rate.")
	v.Check((in.ProductID == nil) != (in.RecipeID == nil), "product_id", "required", "Elige un producto o una receta (sobrante), no los dos.")
	v.Check(in.RecipeID == nil || in.TrackingMode == recetas.ModeStock, "tracking_mode", "not_allowed", "Un sobrante de receta se lleva por existencia (stock).")
	qty := int64(0)
	if in.QuantityMilli != nil {
		qty = *in.QuantityMilli
	}
	v.Check(qty >= 0, "quantity_milli", "out_of_range", "La existencia no puede ser negativa.")
	if in.TrackingMode == recetas.ModeRate {
		v.Required(in.RateQuantityMilli != nil, "rate_quantity_milli")
		v.Required(in.RatePeriodDays != nil, "rate_period_days")
		v.Required(in.CountedDate != nil, "counted_date")
	}
	for f, d := range map[string]*string{"counted_date": in.CountedDate, "expires_date": in.ExpiresDate} {
		if d != nil {
			if _, err := v.Date(*d, f); err != nil {
				return err
			}
		}
	}
	for t, id := range map[string]*int64{"product": in.ProductID, "recipe": in.RecipeID} {
		if id != nil {
			ok, err := dominio.CBExists(ctx, h.db(), t, *id)
			if err != nil {
				return err
			}
			v.Check(ok, t+"_id", "not_found", "No existe ese registro.")
		}
	}
	if err := v.Err(); err != nil {
		return err
	}
	var id int64
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `INSERT INTO inventory_item (product_id, recipe_id, tracking_mode, quantity_milli, counted_date, rate_quantity_milli,
 rate_period_days, reorder_level_milli, expires_date, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, in.ProductID, in.RecipeID, in.TrackingMode, qty,
			in.CountedDate, in.RateQuantityMilli, in.RatePeriodDays, in.ReorderLevelMilli, in.ExpiresDate, in.Notes)
		if err != nil {
			return checkInvConstraint(err)
		}
		if id, err = res.LastInsertId(); err != nil {
			return err
		}
		if err := dominio.CBAuditWrite(ctx, tx, api.Actor(r), "insert", "inventory_item", id, nil, api.RequestID(r)); err != nil {
			return err
		}
		if qty == 0 {
			return nil
		}
		// Invariante 12 de S1: la existencia inicial deja su movimiento (motor de recetas).
		day := h.today().String()
		if in.CountedDate != nil {
			day = *in.CountedDate
		}
		d, err := recetas.ParseDate(day)
		if err != nil {
			return err
		}
		var mv recetas.Movement
		if in.RecipeID != nil {
			mv = recetas.LeftoverMovement(code(id), d, qty)
		} else if mv, err = recetas.AdjustmentMovement(code(id), d, recetas.ReasonCount, 0, qty); err != nil {
			return err
		}
		return h.saveMovement(ctx, tx, r, id, mv, nil)
	})
	if err != nil {
		return err
	}
	return h.writeInventory(w, r, http.StatusCreated, id)
}

func (h *handlers) saveMovement(ctx context.Context, tx *sql.Tx, r *http.Request, itemID int64, mv recetas.Movement, notes *string) error {
	var reason *string
	if mv.Reason != "" {
		reason = &mv.Reason
	}
	mid, err := dominio.CBInsertMovement(ctx, tx, itemID, mv.Date.String(), mv.Kind, reason, mv.DeltaMilli, notes)
	if err != nil {
		return err
	}
	return dominio.CBAuditWrite(ctx, tx, api.Actor(r), "insert", "inventory_movement", mid, nil, api.RequestID(r))
}

func (h *handlers) listInventoryMovements(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "inventory_item_id")
	if err != nil {
		return err
	}
	p, err := api.ParsePage(r, []string{"movement_date"}, "-movement_date")
	if err != nil {
		return err
	}
	if _, err := dominio.CBGetInventory(r.Context(), h.db(), id); err != nil {
		return err
	}
	rows, err := dominio.CBQueryMovements(r.Context(), h.db(), id, p.OrderBy(map[string]string{"movement_date": "m.movement_date", "id": "m.id"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.InventoryMovement, 0, len(rows))
	for _, m := range rows {
		it := api.InventoryMovement{ID: m.ID, InventoryItemID: m.ItemID, MovementDate: m.Date, Kind: m.Kind, QuantityDeltaMilli: m.Delta,
			CreatedAt: api.APITime(m.CreatedAt), UpdatedAt: api.APITime(m.UpdatedAt)}
		if m.Reason.Valid {
			it.Reason = api.Ptr(m.Reason.String)
		}
		if m.Notes.Valid {
			it.Notes = api.Ptr(m.Notes.String)
		}
		for _, x := range []struct {
			v   sql.NullInt64
			dst **int64
		}{{m.ExpenseItemID, &it.ExpenseItemID}, {m.MealPlanEntryID, &it.MealPlanEntryID}, {m.ExpenseID, &it.ExpenseID}} {
			if x.v.Valid {
				*x.dst = api.Ptr(x.v.Int64)
			}
		}
		items = append(items, it)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

// adjustInventoryItem (N-183, RN-15): cantidad nueva o diferencia con motivo; el movimiento lo arma el motor de
// recetas (AdjustmentMovement) y se guarda en la misma transacción que la existencia.
func (h *handlers) adjustInventoryItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := api.PathID(r, "inventory_item_id")
	if err != nil {
		return err
	}
	var in api.InventoryAdjustmentRequest
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	cur, err := dominio.CBGetInventory(ctx, h.db(), id)
	if err != nil {
		return err
	}
	var v api.Validator
	v.Check(oneOf(in.Reason, recetas.ReasonCount, recetas.ReasonConsumption, recetas.ReasonWaste), "reason", "invalid_format", "El motivo debe ser count, consumption o waste.")
	v.Check((in.QuantityMilli == nil) != (in.QuantityDeltaMilli == nil), "quantity_milli", "required", "Manda la cantidad nueva o la diferencia, solo una de las dos.")
	day := h.today().String()
	if in.MovementDate != nil {
		if _, err := v.Date(*in.MovementDate, "movement_date"); err != nil {
			return err
		}
		day = *in.MovementDate
	}
	counted := day
	if in.CountedDate != nil {
		if _, err := v.Date(*in.CountedDate, "counted_date"); err != nil {
			return err
		}
		counted = *in.CountedDate
	}
	newQty := cur.QuantityMilli
	switch {
	case in.QuantityMilli != nil:
		newQty = *in.QuantityMilli
	case in.QuantityDeltaMilli != nil:
		newQty = cur.QuantityMilli + *in.QuantityDeltaMilli
	}
	v.Check(newQty >= 0, "quantity_milli", "out_of_range", "La existencia no puede quedar negativa.")
	if err := v.Err(); err != nil {
		return err
	}
	d, err := recetas.ParseDate(day)
	if err != nil {
		return err
	}
	mv, err := recetas.AdjustmentMovement(code(id), d, in.Reason, cur.QuantityMilli, newQty)
	if err != nil {
		return dominio.Invalid(dominio.FieldIssue{Field: "reason", Code: "invalid_format", Message: "El ajuste no es válido."})
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "inventory_item", id)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE inventory_item SET quantity_milli = ?, counted_date = ? WHERE id = ?`, newQty, counted, id); err != nil {
			return checkInvConstraint(err)
		}
		if err := dominio.CBAuditWrite(ctx, tx, api.Actor(r), "update", "inventory_item", id, before, api.RequestID(r)); err != nil {
			return err
		}
		return h.saveMovement(ctx, tx, r, id, mv, in.Notes)
	})
	if err != nil {
		return err
	}
	return h.writeInventory(w, r, http.StatusOK, id)
}

// updateInventoryItem (JSON Merge Patch): cambia modo, ritmo, nivel de reposición, caducidad y notas; la existencia
// solo cambia con adjustInventoryItem. Las reglas de S1 (CHECK) responden 422.
func (h *handlers) updateInventoryItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := api.PathID(r, "inventory_item_id")
	if err != nil {
		return err
	}
	var in api.InventoryItemUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if _, err := dominio.CBGetInventory(ctx, h.db(), id); err != nil {
		return err
	}
	sets, args := []string{}, []any{}
	var v api.Validator
	if in.TrackingMode != nil {
		v.Check(oneOf(*in.TrackingMode, recetas.ModeStock, recetas.ModeRate), "tracking_mode", "invalid_format", "El modo debe ser stock o consumption_rate.")
		sets, args = append(sets, "tracking_mode = ?"), append(args, *in.TrackingMode)
	}
	addInt := func(col string, n api.Nullable[int64]) {
		if n.Set {
			var x any
			if n.Valid {
				x = n.Value
			}
			sets, args = append(sets, col+" = ?"), append(args, x)
		}
	}
	addInt("rate_quantity_milli", in.RateQuantityMilli)
	addInt("rate_period_days", in.RatePeriodDays)
	addInt("reorder_level_milli", in.ReorderLevelMilli)
	if in.ExpiresDate.Set {
		var x any
		if in.ExpiresDate.Valid {
			if _, err := v.Date(in.ExpiresDate.Value, "expires_date"); err != nil {
				return err
			}
			x = in.ExpiresDate.Value
		}
		sets, args = append(sets, "expires_date = ?"), append(args, x)
	}
	if in.Notes.Set {
		var x any
		if in.Notes.Valid {
			x = in.Notes.Value
		}
		sets, args = append(sets, "notes = ?"), append(args, x)
	}
	if len(sets) == 0 {
		return api.BadRequest("Manda al menos un campo para cambiar.")
	}
	if err := v.Err(); err != nil {
		return err
	}
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		before, err := dominio.RowJSON(ctx, tx, "inventory_item", id)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE inventory_item SET `+strings.Join(sets, ", ")+` WHERE id = ?`, append(args, id)...); err != nil {
			return checkInvConstraint(err)
		}
		return dominio.CBAuditWrite(ctx, tx, api.Actor(r), "update", "inventory_item", id, before, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	return h.writeInventory(w, r, http.StatusOK, id)
}
