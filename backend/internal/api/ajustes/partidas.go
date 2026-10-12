package ajustes

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

// plannedOut arma la respuesta. next_date (x-calculado, opcional) se omite: PROVISIONAL, bloqueo S6-ajustes-api
// (motor_fin no exporta las ocurrencias de una partida y aquí no se recalcula a mano).
func plannedOut(p dominio.PlannedItemRow) api.PlannedItem {
	return api.PlannedItem{ID: p.ID, Kind: p.Kind, Description: p.Description, AmountCents: p.AmountCents,
		SecondHalfAmountCents: pInt(p.SecondHalfAmountCents), Frequency: p.Frequency, IntervalCount: p.IntervalCount,
		StartDate: p.StartDate, EndDate: pStr(p.EndDate), CategoryID: pInt(p.CategoryID), PersonID: pInt(p.PersonID),
		CostShareRuleID: pInt(p.CostShareRuleID), PaymentMethod: pStr(p.PaymentMethod), CreditID: pInt(p.CreditID),
		TermCount: pInt(p.TermCount), AccrualMonths: p.AccrualMonths, IsActive: p.IsActive, Notes: pStr(p.Notes),
		CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt)}
}

// plannedEnums: valores fuera de la lista de C1 son 400 (formato).
func plannedEnums(p dominio.PlannedItemRow) error {
	switch {
	case !oneOf(p.Kind, "income", "fixed_expense", "upcoming_expense"):
		return badEnum("kind")
	case !oneOf(p.Frequency, "weekly", "fortnightly", "monthly", "yearly", "once"):
		return badEnum("frequency")
	case p.PaymentMethod.Valid && !oneOf(p.PaymentMethod.String, "cash", "debit", "transfer", "credit"):
		return badEnum("payment_method")
	}
	return nil
}

// plannedIssues son las reglas de negocio de planned_item (S1: CHECK y disparador del titular en ingresos) como 422
// por campo, más las referencias que deben existir.
func plannedIssues(ctx context.Context, db dominio.Querier, p dominio.PlannedItemRow) ([]dominio.FieldIssue, error) {
	var v api.Validator
	v.Required(strings.TrimSpace(p.Description) != "", "description")
	v.Check(p.AmountCents >= 0, "amount_cents", "out_of_range", "El monto no puede ser negativo.")
	if p.SecondHalfAmountCents.Valid {
		v.Check(p.SecondHalfAmountCents.Int64 >= 0, "second_half_amount_cents", "out_of_range", "El monto no puede ser negativo.")
		v.Check(p.Frequency == "fortnightly", "second_half_amount_cents", "not_allowed", "El monto de la quincena del 30 solo va en partidas quincenales.")
	}
	v.Check(p.IntervalCount >= 1 && p.IntervalCount <= 24, "interval_count", "out_of_range", "La repetición va de 1 a 24.")
	v.Check(p.AccrualMonths >= 1 && p.AccrualMonths <= 24, "accrual_months", "out_of_range", "Los meses de devengo van de 1 a 24.")
	if p.EndDate.Valid {
		v.Check(p.EndDate.String >= p.StartDate, "end_date", "out_of_range", "La fecha final no puede ser anterior a la de inicio.")
	}
	if p.Frequency == "once" {
		v.Check(p.IntervalCount == 1, "interval_count", "not_allowed", "Una partida única no se repite.")
		v.Check(!p.EndDate.Valid || p.EndDate.String == p.StartDate, "end_date", "not_allowed", "En una partida única la fecha final es la misma de inicio.")
	}
	if p.Kind == "income" {
		v.Check(!p.CategoryID.Valid, "category_id", "not_allowed", "Un ingreso no lleva categoría.")
		v.Check(!p.CostShareRuleID.Valid, "cost_share_rule_id", "not_allowed", "Un ingreso no lleva reparto.")
		v.Check(!p.PaymentMethod.Valid, "payment_method", "not_allowed", "Un ingreso no lleva forma de pago.")
		v.Check(!p.CreditID.Valid, "credit_id", "not_allowed", "Un ingreso no lleva crédito.")
		v.Check(p.AccrualMonths == 1, "accrual_months", "not_allowed", "Un ingreso se devenga en un solo mes.")
	} else {
		v.Check(p.CategoryID.Valid, "category_id", "required", "Elige la categoría del gasto.")
	}
	v.Check((p.PaymentMethod.Valid && p.PaymentMethod.String == "credit") == p.CreditID.Valid, "credit_id", "required",
		"Con forma de pago a crédito elige la línea de crédito (y solo entonces).")
	v.Check(p.CreditID.Valid == p.TermCount.Valid, "term_count", "required", "El plazo va junto con el crédito.")
	if p.TermCount.Valid {
		v.Check(p.TermCount.Int64 >= 1 && p.TermCount.Int64 <= 120, "term_count", "out_of_range", "El plazo va de 1 a 120 pagos.")
	}
	refs := []struct {
		table, field, msg string
		id                sql.NullInt64
	}{
		{"category", "category_id", "La categoría no existe.", p.CategoryID},
		{"person", "person_id", "La persona no existe.", p.PersonID},
		{"cost_share_rule", "cost_share_rule_id", "La regla de reparto no existe.", p.CostShareRuleID},
		{"credit", "credit_id", "El crédito no existe.", p.CreditID},
	}
	for _, x := range refs {
		if !x.id.Valid {
			continue
		}
		ok, err := dominio.RefExists(ctx, db, x.table, x.id.Int64)
		if err != nil {
			return nil, err
		}
		v.Check(ok, x.field, "not_found", x.msg)
	}
	if p.Kind == "income" && p.PersonID.Valid {
		pr, err := dominio.PersonRecordByID(ctx, db, p.PersonID.Int64)
		if err == nil {
			v.Check(pr.Participation == "holder", "person_id", "not_allowed", "Los ingresos van solo a nombre del titular.")
		}
	}
	if err, ok := dominio.AsError(v.Err()); ok {
		return err.Fields, nil
	}
	return nil, nil
}

func checkPlanned(ctx context.Context, db dominio.Querier, p dominio.PlannedItemRow) error {
	if err := plannedEnums(p); err != nil {
		return err
	}
	issues, err := plannedIssues(ctx, db, p)
	if err != nil {
		return err
	}
	if len(issues) > 0 {
		return dominio.Invalid(issues...)
	}
	return nil
}

// listPlannedItems: filtros category_id, credit_id (ids separados por coma) y kind. Orden por defecto: start_date
// (C1 no lo fija; PROVISIONAL, bloqueo S6-ajustes-api).
func listPlannedItems(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"start_date", "description", "amount_cents", "kind"}, "start_date")
		if err != nil {
			return err
		}
		var f dominio.PlannedItemFilter
		if f.CategoryIDs, err = api.QueryIDs(r, "category_id"); err != nil {
			return err
		}
		if f.CreditIDs, err = api.QueryIDs(r, "credit_id"); err != nil {
			return err
		}
		f.Kind = r.URL.Query().Get("kind")
		if f.Kind != "" && !oneOf(f.Kind, "income", "fixed_expense", "upcoming_expense") {
			return badEnum("kind")
		}
		order := p.OrderBy(map[string]string{"start_date": "p.start_date", "description": "p.description",
			"amount_cents": "p.amount_cents", "kind": "p.kind", "id": "p.id"}) + " " + p.LimitOffset()
		ps, err := dominio.ListPlannedItemRows(r.Context(), reg.Deps().DB, f, order)
		if err != nil {
			return err
		}
		items := make([]api.PlannedItem, 0, len(ps))
		for _, x := range ps {
			items = append(items, plannedOut(x))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// plannedFromCreate aplica los valores por omisión de S1 (interval_count 1, accrual_months 1, is_active true).
func plannedFromCreate(in api.PlannedItemCreate) dominio.PlannedItemRow {
	p := dominio.PlannedItemRow{Kind: in.Kind, Description: strings.TrimSpace(in.Description), AmountCents: in.AmountCents,
		SecondHalfAmountCents: nInt(in.SecondHalfAmountCents), Frequency: in.Frequency, IntervalCount: 1,
		StartDate: in.StartDate, EndDate: nStr(in.EndDate), CategoryID: nInt(in.CategoryID), PersonID: nInt(in.PersonID),
		CostShareRuleID: nInt(in.CostShareRuleID), PaymentMethod: nStr(in.PaymentMethod), CreditID: nInt(in.CreditID),
		TermCount: nInt(in.TermCount), AccrualMonths: 1, IsActive: in.IsActive == nil || *in.IsActive, Notes: nStr(in.Notes)}
	if in.IntervalCount != nil {
		p.IntervalCount = *in.IntervalCount
	}
	if in.AccrualMonths != nil {
		p.AccrualMonths = *in.AccrualMonths
	}
	return p
}

func checkDates(p dominio.PlannedItemRow) error {
	var v api.Validator
	if _, err := v.Date(p.StartDate, "start_date"); err != nil {
		return err
	}
	if p.EndDate.Valid {
		if _, err := v.Date(p.EndDate.String, "end_date"); err != nil {
			return err
		}
	}
	return nil
}

func createPlannedItem(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.PlannedItemCreate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		var v api.Validator
		for _, k := range []string{"kind", "description", "amount_cents", "frequency", "start_date"} {
			_, ok := keys[k]
			v.Required(ok, k)
		}
		if err := v.Err(); err != nil {
			return err
		}
		p := plannedFromCreate(in)
		if err := checkDates(p); err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		if err := checkPlanned(ctx, db, p); err != nil {
			return err
		}
		var id int64
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			id, err = dominio.InsertOwned(ctx, tx, "planned_item", p.Columns(), api.Actor(r), api.RequestID(r))
			return err
		})
		if err := writeErr(err, "Ya existe esa partida."); err != nil {
			return err
		}
		out, err := dominio.PlannedItemRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, plannedOut(out))
	}
}

func updatePlannedItem(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "planned_item_id")
		if err != nil {
			return err
		}
		var in api.PlannedItemUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return noChanges()
		}
		ctx, db := r.Context(), reg.Deps().DB
		p, err := dominio.PlannedItemRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		if in.Kind != nil {
			p.Kind = *in.Kind
		}
		if in.Description != nil {
			p.Description = strings.TrimSpace(*in.Description)
		}
		if in.AmountCents != nil {
			p.AmountCents = *in.AmountCents
		}
		if in.Frequency != nil {
			p.Frequency = *in.Frequency
		}
		if in.IntervalCount != nil {
			p.IntervalCount = *in.IntervalCount
		}
		if in.StartDate != nil {
			p.StartDate = *in.StartDate
		}
		if in.AccrualMonths != nil {
			p.AccrualMonths = *in.AccrualMonths
		}
		if in.IsActive != nil {
			p.IsActive = *in.IsActive
		}
		p.SecondHalfAmountCents = mergeInt(p.SecondHalfAmountCents, in.SecondHalfAmountCents)
		p.EndDate = mergeStr(p.EndDate, in.EndDate)
		p.CategoryID = mergeInt(p.CategoryID, in.CategoryID)
		p.PersonID = mergeInt(p.PersonID, in.PersonID)
		p.CostShareRuleID = mergeInt(p.CostShareRuleID, in.CostShareRuleID)
		p.PaymentMethod = mergeStr(p.PaymentMethod, in.PaymentMethod)
		p.CreditID = mergeInt(p.CreditID, in.CreditID)
		p.TermCount = mergeInt(p.TermCount, in.TermCount)
		p.Notes = mergeStr(p.Notes, in.Notes)
		if err := checkDates(p); err != nil {
			return err
		}
		if err := checkPlanned(ctx, db, p); err != nil {
			return err
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			return dominio.UpdateOwned(ctx, tx, "planned_item", id, p.Columns(), api.Actor(r), api.RequestID(r))
		})
		if err := writeErr(err, "Ya existe esa partida."); err != nil {
			return err
		}
		out, err := dominio.PlannedItemRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, plannedOut(out))
	}
}
