package presupuesto

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
)

// ---------- Ingresos ----------

func incomeOut(x dominio.IncomeRecord) api.Income {
	return api.Income{ID: x.ID, PersonID: x.PersonID, Kind: x.Kind, PlannedItemID: optID(x.PlannedItemID), Description: x.Description,
		AmountCents: x.AmountCents, IncomeDate: x.IncomeDate, Notes: optStr(x.Notes), CreatedAt: api.APITime(x.CreatedAt),
		UpdatedAt: api.APITime(x.UpdatedAt), CreatedByUserAccountID: optID(x.CreatedBy)}
}

func (h *handlers) listIncomes(w http.ResponseWriter, r *http.Request) error {
	// Supuesto: C1 no fija el orden por defecto de ingresos; se usa -income_date (los más recientes primero).
	pg, err := api.ParsePage(r, []string{"income_date", "amount_cents"}, "-income_date")
	if err != nil {
		return err
	}
	p, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	kind := r.URL.Query().Get("kind")
	if kind != "" && kind != "fixed" && kind != "extra" {
		return api.BadRequest("El filtro kind debe ser fixed o extra.")
	}
	from, to := "", ""
	if p != nil {
		from, to = p.Start.String(), p.End.String()
	}
	recs, err := dominio.ListIncomes(r.Context(), h.db(), from, to, kind,
		pg.OrderBy(map[string]string{"income_date": "i.income_date", "amount_cents": "i.amount_cents", "id": "i.id"})+" "+pg.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.Income, 0, len(recs))
	for _, x := range recs {
		items = append(items, incomeOut(x))
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
}

func (h *handlers) createIncome(w http.ResponseWriter, r *http.Request) error {
	var in api.IncomeCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	var v api.Validator
	if _, err := v.Date(in.IncomeDate, "income_date"); err != nil {
		return err
	}
	v.Check(in.Kind == "fixed" || in.Kind == "extra", "kind", "invalid_format", "El tipo de ingreso es fixed (quincena) o extra.")
	v.Check(strings.TrimSpace(in.Description) != "", "description", "required", "Escribe el concepto del ingreso.")
	v.PositiveCents(in.AmountCents, "amount_cents")
	if in.PlannedItemID != nil {
		ok, err := dominio.RowExists(ctx, h.db(), "planned_item", *in.PlannedItemID, "kind = 'income'")
		if err != nil {
			return err
		}
		v.Check(ok, "planned_item_id", "not_found", "No encontramos esa partida de ingreso planificada.")
	}
	if err := v.Err(); err != nil {
		return err
	}
	var out api.Income
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		holder, err := dominio.HolderPersonID(ctx, tx)
		if err != nil {
			return err
		}
		id, err := dominio.BudgetAreaInsert(ctx, tx, "income", []string{"person_id", "kind", "planned_item_id", "description", "amount_cents", "income_date", "notes"},
			[]any{holder, in.Kind, in.PlannedItemID, strings.TrimSpace(in.Description), in.AmountCents, in.IncomeDate, in.Notes}, api.Actor(r), api.RequestID(r))
		if err != nil {
			return constraintErr(err)
		}
		x, err := dominio.GetIncome(ctx, tx, id)
		out = incomeOut(x)
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func (h *handlers) updateIncome(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "income_id")
	if err != nil {
		return err
	}
	var in api.IncomeUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	var v api.Validator
	var cols []string
	var vals []any
	if in.Description != nil {
		v.Check(strings.TrimSpace(*in.Description) != "", "description", "required", "Escribe el concepto del ingreso.")
		cols, vals = append(cols, "description"), append(vals, strings.TrimSpace(*in.Description))
	}
	if in.AmountCents != nil {
		v.PositiveCents(*in.AmountCents, "amount_cents")
		cols, vals = append(cols, "amount_cents"), append(vals, *in.AmountCents)
	}
	if in.IncomeDate != nil {
		if _, err := v.Date(*in.IncomeDate, "income_date"); err != nil {
			return err
		}
		cols, vals = append(cols, "income_date"), append(vals, *in.IncomeDate)
	}
	if in.Notes.Set {
		cols, vals = append(cols, "notes"), append(vals, sql.NullString{String: in.Notes.Value, Valid: in.Notes.Valid})
	}
	v.Check(len(cols) > 0, "body", "required", "Manda al menos un cambio.")
	if err := v.Err(); err != nil {
		return err
	}
	ctx := r.Context()
	var out api.Income
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetIncome(ctx, tx, id); err != nil {
			return err
		}
		if err := dominio.BudgetAreaUpdate(ctx, tx, "income", id, cols, vals, h.now(), api.Actor(r), api.RequestID(r)); err != nil {
			return constraintErr(err)
		}
		x, err := dominio.GetIncome(ctx, tx, id)
		out = incomeOut(x)
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) deleteIncome(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "income_id")
	if err != nil {
		return err
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetIncome(ctx, tx, id); err != nil {
			return err
		}
		return dominio.BudgetAreaDelete(ctx, tx, "income", id, api.Actor(r), api.RequestID(r))
	}); err != nil {
		return err
	}
	return api.NoContent(w)
}

// incomeMonth arma IncomeMonth: quincena del 15 y del 30 (esperada o capturada) y extras del mes.
func (h *handlers) incomeMonth(ctx context.Context, m dominio.Date) (api.IncomeMonth, error) {
	q := h.db()
	out := api.IncomeMonth{Month: dominio.EngineMonth(m).String(), Extras: []api.Income{}, Fortnights: []api.IncomeMonthFortnightsItem{}}
	q1 := dominio.FortnightOf(m)
	for i, f := range []dominio.Fortnight{q1, q1.Next()} {
		exp, err := dominio.ExpectedIncome(ctx, q, f.PayDay().String(), i == 1)
		if err != nil {
			return out, err
		}
		recs, err := dominio.ListIncomes(ctx, q, f.Start().String(), f.End().String(), "fixed", "ORDER BY i.income_date, i.id")
		if err != nil {
			return out, err
		}
		it := api.IncomeMonthFortnightsItem{PayDate: f.PayDay().String(), Status: "expected", ExpectedCents: exp}
		if len(recs) > 0 {
			var sum int64
			for _, x := range recs {
				sum += x.AmountCents
			}
			it.Status, it.IncomeID, it.AmountCents = "captured", api.Ptr(recs[0].ID), api.Ptr(sum)
			out.CapturedCents += sum
		} else {
			out.ExpectedPendingCents += exp
		}
		out.Fortnights = append(out.Fortnights, it)
	}
	extras, err := dominio.ListIncomes(ctx, q, m.String(), dominio.MonthEnd(m).String(), "extra", "ORDER BY i.income_date, i.id")
	if err != nil {
		return out, err
	}
	for _, x := range extras {
		out.Extras = append(out.Extras, incomeOut(x))
		out.ExtrasCents += x.AmountCents
	}
	out.TotalCents = out.CapturedCents + out.ExpectedPendingCents + out.ExtrasCents
	return out, nil
}

func (h *handlers) getIncomeMonth(w http.ResponseWriter, r *http.Request) error {
	m, err := h.monthParam(r)
	if err != nil {
		return err
	}
	out, err := h.incomeMonth(r.Context(), m)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) getIncomeMonthlyTrend(w http.ResponseWriter, r *http.Request) error {
	t := h.today()
	first := dominio.MonthStart(t.Year, t.Month-11)
	sums, err := dominio.IncomeByMonth(r.Context(), h.db(), first.String(), dominio.MonthEnd(t).String())
	if err != nil {
		return err
	}
	out := api.IncomeMonthlyTrend{Months: []api.IncomeMonthlyTrendMonthsItem{}}
	for i := 0; i < 12; i++ {
		m := dominio.MonthStart(first.Year, first.Month+timeMonth(i)).MonthKey()
		s := sums[m]
		out.Months = append(out.Months, api.IncomeMonthlyTrendMonthsItem{Month: m, FixedCents: s[0], ExtraCents: s[1], TotalCents: s[0] + s[1]})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// ---------- Ahorro ----------

func (h *handlers) savingsAt(fe *dominio.FinEngine, ref dominio.Date) (mf.SavingsResult, error) {
	in, err := fe.At(ref, nil)
	if err != nil {
		return mf.SavingsResult{}, err
	}
	return mf.SavingsGoals(in)
}

func (h *handlers) listSavingsGoals(w http.ResponseWriter, r *http.Request) error {
	pg, err := api.ParsePage(r, []string{"target_date", "name", "start_date"}, "target_date")
	if err != nil {
		return err
	}
	status := r.URL.Query().Get("status")
	if status != "" && status != "active" && status != "archived" {
		return api.BadRequest("El filtro status debe ser active o archived.")
	}
	ctx := r.Context()
	recs, err := dominio.ListSavingsGoals(ctx, h.db(), status != "archived",
		pg.OrderBy(map[string]string{"target_date": "target_date", "name": "name", "start_date": "start_date", "id": "id"})+" "+pg.LimitOffset())
	if err != nil {
		return err
	}
	fe, err := dominio.LoadFinEngine(ctx, h.db())
	if err != nil {
		return err
	}
	sv, err := h.savingsAt(fe, h.today())
	if err != nil {
		return err
	}
	st := map[int64]mf.GoalStatus{}
	for _, g := range sv.Goals {
		st[g.GoalID] = g
	}
	items := make([]api.SavingsGoal, 0, len(recs))
	for _, g := range recs {
		it := api.SavingsGoal{ID: g.ID, Name: g.Name, Kind: g.Kind, TargetCents: optID(g.TargetCents), TargetDate: optStr(g.TargetDate),
			PlannedMonthlyCents: optID(g.PlannedMonthlyCents), StartDate: g.StartDate, IsActive: g.IsActive, Notes: optStr(g.Notes),
			CreatedAt: api.APITime(g.CreatedAt), UpdatedAt: api.APITime(g.UpdatedAt)}
		if s, ok := st[g.ID]; ok {
			it.SavedCents, it.ProgressTenthsPct, it.ProgressPct = api.Ptr(cents(s.SavedCents)), s.ProgressTenthsPct, s.ProgressPct
			it.RequiredMonthlyCents, it.Status = centsPtr(s.NeededMonthlyCents), api.Ptr(s.State)
		}
		items = append(items, it)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
}

func (h *handlers) getSavingsSummary(w http.ResponseWriter, r *http.Request) error {
	fe, err := dominio.LoadFinEngine(r.Context(), h.db())
	if err != nil {
		return err
	}
	sv, err := h.savingsAt(fe, h.today())
	if err != nil {
		return err
	}
	out := api.SavingsSummary{SavedTotalCents: cents(sv.MonthSummary.SavedCents), Goals: []api.SavingsSummaryGoalsItem{},
		ContributedMonthCents: cents(sv.MonthSummary.ContributedThisMonthCents), PlannedMonthCents: cents(sv.MonthSummary.PlannedThisMonthCents)}
	for _, g := range sv.Goals {
		out.Goals = append(out.Goals, api.SavingsSummaryGoalsItem{SavingsGoalID: g.GoalID, SavedCents: cents(g.SavedCents)})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) getSavingsMonthlyTrend(w http.ResponseWriter, r *http.Request) error {
	fe, err := dominio.LoadFinEngine(r.Context(), h.db())
	if err != nil {
		return err
	}
	t := h.today()
	out := api.SavingsMonthlyTrend{Months: []api.SavingsMonthlyTrendMonthsItem{}}
	for i := 11; i >= 0; i-- {
		m := dominio.MonthStart(t.Year, t.Month-timeMonth(i))
		ref := dominio.MonthEnd(m)
		if i == 0 {
			ref = t
		}
		sv, err := h.savingsAt(fe, ref)
		if err != nil {
			return err
		}
		out.Months = append(out.Months, api.SavingsMonthlyTrendMonthsItem{Month: m.MonthKey(), SavedEndCents: cents(sv.MonthSummary.SavedCents),
			ContributedCents: cents(sv.MonthSummary.ContributedThisMonthCents), PlannedCents: cents(sv.MonthSummary.PlannedThisMonthCents)})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func movementOut(m dominio.SavingsMovementRecord) api.SavingsMovement {
	return api.SavingsMovement{ID: m.ID, SavingsGoalID: m.SavingsGoalID, MovementDate: m.MovementDate, Kind: m.Kind, AmountCents: m.AmountCents,
		Notes: optStr(m.Notes), CreatedAt: api.APITime(m.CreatedAt), UpdatedAt: api.APITime(m.UpdatedAt)}
}

func (h *handlers) listSavingsMovements(w http.ResponseWriter, r *http.Request) error {
	pg, err := api.ParsePage(r, []string{"movement_date", "amount_cents"}, "-movement_date")
	if err != nil {
		return err
	}
	p, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	goals, err := api.QueryIDs(r, "savings_goal_id")
	if err != nil {
		return err
	}
	if len(goals) > 1 {
		return api.BadRequest("El filtro savings_goal_id lleva un solo destino.")
	}
	from, to, goal := "", "", int64(0)
	if p != nil {
		from, to = p.Start.String(), p.End.String()
	}
	if len(goals) == 1 {
		goal = goals[0]
	}
	recs, err := dominio.ListSavingsMovements(r.Context(), h.db(), from, to, goal,
		pg.OrderBy(map[string]string{"movement_date": "movement_date", "amount_cents": "amount_cents", "id": "id"})+" "+pg.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.SavingsMovement, 0, len(recs))
	for _, m := range recs {
		items = append(items, movementOut(m))
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
}

func (h *handlers) goalCheck(ctx context.Context, v *api.Validator, id int64) error {
	ok, err := dominio.RowExists(ctx, h.db(), "savings_goal", id, "")
	if err != nil {
		return err
	}
	v.Check(ok, "savings_goal_id", "not_found", "No encontramos ese fondo o meta.")
	return nil
}

func (h *handlers) createSavingsMovement(w http.ResponseWriter, r *http.Request) error {
	var in api.SavingsMovementCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	var v api.Validator
	if _, err := v.Date(in.MovementDate, "movement_date"); err != nil {
		return err
	}
	v.Check(in.Kind == "deposit" || in.Kind == "withdrawal", "kind", "invalid_format", "El movimiento es deposit (aportación) o withdrawal (retiro).")
	v.PositiveCents(in.AmountCents, "amount_cents")
	if err := h.goalCheck(ctx, &v, in.SavingsGoalID); err != nil {
		return err
	}
	if err := v.Err(); err != nil {
		return err
	}
	var out api.SavingsMovement
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		id, err := dominio.BudgetAreaInsert(ctx, tx, "savings_movement", []string{"savings_goal_id", "movement_date", "kind", "amount_cents", "notes"},
			[]any{in.SavingsGoalID, in.MovementDate, in.Kind, in.AmountCents, in.Notes}, api.Actor(r), api.RequestID(r))
		if err != nil {
			return constraintErr(err)
		}
		m, err := dominio.GetSavingsMovement(ctx, tx, id)
		out = movementOut(m)
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func (h *handlers) updateSavingsMovement(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "savings_movement_id")
	if err != nil {
		return err
	}
	var in api.SavingsMovementUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	var v api.Validator
	var cols []string
	var vals []any
	if in.SavingsGoalID != nil {
		if err := h.goalCheck(ctx, &v, *in.SavingsGoalID); err != nil {
			return err
		}
		cols, vals = append(cols, "savings_goal_id"), append(vals, *in.SavingsGoalID)
	}
	if in.MovementDate != nil {
		if _, err := v.Date(*in.MovementDate, "movement_date"); err != nil {
			return err
		}
		cols, vals = append(cols, "movement_date"), append(vals, *in.MovementDate)
	}
	if in.AmountCents != nil {
		v.PositiveCents(*in.AmountCents, "amount_cents")
		cols, vals = append(cols, "amount_cents"), append(vals, *in.AmountCents)
	}
	if in.Notes.Set {
		cols, vals = append(cols, "notes"), append(vals, sql.NullString{String: in.Notes.Value, Valid: in.Notes.Valid})
	}
	v.Check(len(cols) > 0, "body", "required", "Manda al menos un cambio.")
	if err := v.Err(); err != nil {
		return err
	}
	var out api.SavingsMovement
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetSavingsMovement(ctx, tx, id); err != nil {
			return err
		}
		if err := dominio.BudgetAreaUpdate(ctx, tx, "savings_movement", id, cols, vals, h.now(), api.Actor(r), api.RequestID(r)); err != nil {
			return constraintErr(err)
		}
		m, err := dominio.GetSavingsMovement(ctx, tx, id)
		out = movementOut(m)
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) deleteSavingsMovement(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "savings_movement_id")
	if err != nil {
		return err
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetSavingsMovement(ctx, tx, id); err != nil {
			return err
		}
		return dominio.BudgetAreaDelete(ctx, tx, "savings_movement", id, api.Actor(r), api.RequestID(r))
	}); err != nil {
		return err
	}
	return api.NoContent(w)
}
