// Package presupuesto implementa las operaciones de presupuesto-a.yaml y presupuesto-b.yaml (C1 v1.4):
// periodos, presupuesto con alcance K2, resúmenes del mes, flujo y proyección, ingresos y ahorro.
// Los cálculos de dinero los hace el motor S2 (internal/motor_fin); aquí solo se arman entradas y respuestas.
package presupuesto

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
)

type handlers struct{ reg *api.Registry }

// Register da de alta las 28 operaciones del área.
func Register(reg *api.Registry) {
	h := &handlers{reg: reg}
	ops := []struct {
		id string
		fn api.HandlerFunc
	}{
		{"listPeriods", h.listPeriods}, {"listBudgets", h.listBudgets}, {"createBudget", h.createBudget},
		{"updateBudget", h.updateBudget}, {"deleteBudget", h.deleteBudget}, {"getBudgetMonth", h.getBudgetMonth},
		{"getBudgetByCategory", h.getBudgetByCategory}, {"getBudgetFortnights", h.getBudgetFortnights},
		{"getAccruedGroceries", h.getAccruedGroceries}, {"getBudgetAlerts", h.getBudgetAlerts},
		{"getBudgetMonthlyTrend", h.getBudgetMonthlyTrend}, {"getCashFlowFortnights", h.getCashFlowFortnights},
		{"getAccrualVsCash", h.getAccrualVsCash}, {"getCashFlowAssumptions", h.getCashFlowAssumptions},
		{"getLiveDebt", h.getLiveDebt}, {"listIncomes", h.listIncomes}, {"createIncome", h.createIncome},
		{"getIncomeMonth", h.getIncomeMonth}, {"getIncomeMonthlyTrend", h.getIncomeMonthlyTrend},
		{"updateIncome", h.updateIncome}, {"deleteIncome", h.deleteIncome}, {"getSavingsSummary", h.getSavingsSummary},
		{"getSavingsMonthlyTrend", h.getSavingsMonthlyTrend}, {"listSavingsGoals", h.listSavingsGoals},
		{"listSavingsMovements", h.listSavingsMovements}, {"createSavingsMovement", h.createSavingsMovement},
		{"updateSavingsMovement", h.updateSavingsMovement}, {"deleteSavingsMovement", h.deleteSavingsMovement},
	}
	for _, o := range ops {
		reg.Handle(o.id, o.fn)
	}
}

func (h *handlers) db() *sql.DB         { return h.reg.Deps().DB }
func (h *handlers) today() dominio.Date { return dominio.Today(h.reg.Deps().Clock) }
func (h *handlers) now() string         { return dominio.FormatDBTime(h.reg.Deps().Clock.Now()) }

func invalid(field, code, msg string) error {
	return dominio.Invalid(dominio.FieldIssue{Field: field, Code: code, Message: msg})
}

// engine arma la entrada del motor S2 desde la base y la prepara con «hoy».
func (h *handlers) engine(ctx context.Context) (*dominio.FinEngine, mf.Input, error) {
	fe, err := dominio.LoadFinEngine(ctx, h.db())
	if err != nil {
		return nil, mf.Input{}, err
	}
	in, err := fe.At(h.today(), nil)
	if err != nil {
		return nil, in, fmt.Errorf("motor_fin: entrada inválida: %w", err)
	}
	return fe, in, nil
}

func projectionMode(fe *dominio.FinEngine) string {
	if fe.RepeatCadence {
		return "with_repeat_cadence"
	}
	return "contracted_only"
}

// monthParam lee period como mes (por defecto el mes en curso) y devuelve su día 1.
func (h *handlers) monthParam(r *http.Request) (dominio.Date, error) {
	p, err := api.QueryPeriod(r)
	if err != nil {
		return dominio.Date{}, err
	}
	if p == nil {
		t := h.today()
		return dominio.MonthStart(t.Year, t.Month), nil
	}
	if p.Kind != dominio.PeriodMonth {
		return dominio.Date{}, invalid("period", "not_allowed", "Esta vista es por mes: usa un periodo como 2026-10.")
	}
	return p.Start, nil
}

func cents(c mf.Cents) int64 { return int64(c) }

func centsPtr(c *mf.Cents) *int64 {
	if c == nil {
		return nil
	}
	return api.Ptr(int64(*c))
}

// ---------- Periodos ----------

func (h *handlers) listPeriods(w http.ResponseWriter, r *http.Request) error {
	pg, err := api.ParsePage(r, []string{"start_date", "end_date", "kind"}, "start_date")
	if err != nil {
		return err
	}
	p, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	t := h.today()
	// PROVISIONAL (bloqueo S6-presupuesto-periodos): sin period, del mes en curso − 11 al mes en curso + 12.
	from, to := dominio.MonthStart(t.Year, t.Month-11), dominio.MonthEnd(dominio.MonthStart(t.Year, t.Month+12))
	if p != nil {
		from, to = p.Start, p.End
	}
	if dominio.MonthStart(from.Year, from.Month+36).Before(to) {
		return invalid("period", "out_of_range", "El selector de periodo abarca como máximo 36 meses.")
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		return dominio.EnsurePeriods(ctx, tx, from, to, api.Actor(r), api.RequestID(r))
	}); err != nil {
		return err
	}
	rows, err := dominio.ListPeriods(ctx, h.db(), from.String(), to.String(),
		pg.OrderBy(map[string]string{"start_date": "start_date", "end_date": "end_date", "kind": "kind", "id": "id"})+" "+pg.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.Period, 0, len(rows))
	for _, x := range rows {
		it := api.Period{ID: x.ID, Kind: x.Kind, Year: x.Year, Month: x.Month, StartDate: x.StartDate, EndDate: x.EndDate,
			CreatedAt: api.APITime(x.CreatedAt), UpdatedAt: api.APITime(x.UpdatedAt)}
		label := x.StartDate[:7]
		if x.ParentID.Valid {
			it.ParentID = api.Ptr(x.ParentID.Int64)
		}
		if x.Half.Valid {
			it.Half = api.Ptr(x.Half.Int64)
			label = fmt.Sprintf("%s-q%d", label, x.Half.Int64)
		}
		if x.PayDate.Valid {
			it.PayDate = api.Ptr(x.PayDate.String)
		}
		it.Label = api.Ptr(label)
		items = append(items, it)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
}

// ---------- Presupuesto (budget) ----------

func (h *handlers) budgetOut(ctx context.Context, q dominio.Querier, b dominio.BudgetRecord, followDay dominio.Date) (api.Budget, error) {
	out := api.Budget{ID: b.ID, AmountCents: b.AmountCents, ValidFromDate: b.ValidFrom,
		CreatedAt: api.APITime(b.CreatedAt), UpdatedAt: api.APITime(b.UpdatedAt)}
	if b.CategoryID.Valid {
		out.CategoryID = api.Ptr(b.CategoryID.Int64)
	}
	if b.ValidTo.Valid {
		out.ValidToDate = api.Ptr(b.ValidTo.String)
	}
	if b.Notes.Valid {
		out.Notes = api.Ptr(b.Notes.String)
	}
	single := b.ValidTo.Valid && b.ValidTo.String[:7] == b.ValidFrom[:7]
	out.IsSingleMonth = &single
	amt, err := dominio.BudgetAmountAt(ctx, q, b.CategoryID, followDay.String())
	out.FollowingMonthAmountCents = amt
	return out, err
}

func (h *handlers) listBudgets(w http.ResponseWriter, r *http.Request) error {
	pg, err := api.ParsePage(r, []string{"category_id", "valid_from_date", "amount_cents"}, "category_id,valid_from_date")
	if err != nil {
		return err
	}
	p, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	cats, err := api.QueryIDs(r, "category_id")
	if err != nil {
		return err
	}
	period := dominio.MonthPeriod(h.today())
	if p != nil {
		period = *p
	}
	ctx := r.Context()
	rows, err := dominio.ListBudgets(ctx, h.db(), period.Start.String(), period.End.String(), cats,
		pg.OrderBy(map[string]string{"category_id": "coalesce(category_id, 0)", "valid_from_date": "valid_from_date", "amount_cents": "amount_cents", "id": "id"})+" "+pg.LimitOffset())
	if err != nil {
		return err
	}
	follow := dominio.MonthStart(period.End.Year, period.End.Month+1)
	items := make([]api.Budget, 0, len(rows))
	for _, b := range rows {
		it, err := h.budgetOut(ctx, h.db(), b, follow)
		if err != nil {
			return err
		}
		items = append(items, it)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
}

func overlapErr() error {
	return dominio.Conflict("budget_overlap", "Ya hay un monto para esa categoría en esos meses. Cambia el monto existente con su alcance o elige otra vigencia.")
}

func (h *handlers) checkBudgetDates(v *api.Validator, from string, to *string) {
	f, err := dominio.ParseDate(from)
	v.Check(err == nil && f.Day == 1, "valid_from_date", "invalid_format", "La vigencia empieza el día 1 de un mes.")
	if to != nil {
		t, err := dominio.ParseDate(*to)
		v.Check(err == nil && t == dominio.MonthEnd(t), "valid_to_date", "invalid_format", "La vigencia termina el último día de un mes.")
		v.Check(*to >= from, "valid_to_date", "out_of_range", "La vigencia no puede terminar antes de empezar.")
	}
}

func (h *handlers) createBudget(w http.ResponseWriter, r *http.Request) error {
	var in api.BudgetCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx := r.Context()
	var v api.Validator
	v.Check(in.AmountCents >= 0, "amount_cents", "out_of_range", "El monto no puede ser negativo.")
	h.checkBudgetDates(&v, in.ValidFromDate, in.ValidToDate)
	cat := sql.NullInt64{}
	if in.CategoryID != nil {
		cat = sql.NullInt64{Int64: *in.CategoryID, Valid: true}
		ok, err := dominio.RowExists(ctx, h.db(), "category", *in.CategoryID, "")
		if err != nil {
			return err
		}
		v.Check(ok, "category_id", "not_found", "No encontramos esa categoría.")
	}
	if err := v.Err(); err != nil {
		return err
	}
	to := sql.NullString{}
	if in.ValidToDate != nil {
		to = sql.NullString{String: *in.ValidToDate, Valid: true}
	}
	var out api.Budget
	err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		over, err := dominio.BudgetOverlaps(ctx, tx, cat, in.ValidFromDate, to, 0)
		if err != nil {
			return err
		}
		if over {
			return overlapErr()
		}
		id, err := dominio.BudgetAreaInsert(ctx, tx, "budget", []string{"category_id", "amount_cents", "valid_from_date", "valid_to_date", "notes"},
			[]any{cat, in.AmountCents, in.ValidFromDate, to, in.Notes}, api.Actor(r), api.RequestID(r))
		if err != nil {
			return constraintErr(err)
		}
		b, err := dominio.GetBudget(ctx, tx, id)
		if err != nil {
			return err
		}
		f, _ := dominio.ParseDate(in.ValidFromDate)
		out, err = h.budgetOut(ctx, tx, b, dominio.MonthStart(f.Year, f.Month+1))
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func constraintErr(err error) error {
	if dominio.IsUnique(err) {
		return dominio.Conflict("already_exists", "Ese registro ya existe.")
	}
	if dominio.IsConstraint(err) {
		return invalid("body", "out_of_range", "Algún dato no cumple las reglas del registro. Revisa montos y fechas.")
	}
	return err
}

func (h *handlers) updateBudget(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "budget_id")
	if err != nil {
		return err
	}
	scope := r.URL.Query().Get("scope")
	if scope == "" {
		scope = "onward"
	}
	if scope != "onward" && scope != "this_month" {
		return api.BadRequest("El alcance (scope) debe ser onward o this_month.")
	}
	month := dominio.MonthStart(h.today().Year, h.today().Month)
	if s := r.URL.Query().Get("month"); s != "" {
		p, err := dominio.ParsePeriod(s)
		if err != nil || p.Kind != dominio.PeriodMonth {
			return api.BadRequest("El mes (month) debe tener el formato AAAA-MM.")
		}
		month = p.Start
	}
	var in api.BudgetUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if in.AmountCents == nil && !in.ValidToDate.Set && !in.Notes.Set {
		return invalid("body", "required", "Manda al menos un cambio: monto, fin de la vigencia o notas.")
	}
	ctx := r.Context()
	var out api.Budget
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		b, err := dominio.GetBudget(ctx, tx, id)
		if err != nil {
			return err
		}
		var cols []string
		var vals []any
		if in.ValidToDate.Set {
			to := sql.NullString{String: in.ValidToDate.Value, Valid: in.ValidToDate.Valid}
			var v api.Validator
			var tp *string
			if to.Valid {
				tp = &to.String
			}
			h.checkBudgetDates(&v, b.ValidFrom, tp)
			if err := v.Err(); err != nil {
				return err
			}
			over, err := dominio.BudgetOverlaps(ctx, tx, b.CategoryID, b.ValidFrom, to, b.ID)
			if err != nil {
				return err
			}
			if over {
				return overlapErr()
			}
			cols, vals = append(cols, "valid_to_date"), append(vals, to)
		}
		if in.Notes.Set {
			cols, vals = append(cols, "notes"), append(vals, sql.NullString{String: in.Notes.Value, Valid: in.Notes.Valid})
		}
		if len(cols) > 0 {
			if err := dominio.BudgetAreaUpdate(ctx, tx, "budget", id, cols, vals, h.now(), api.Actor(r), api.RequestID(r)); err != nil {
				return constraintErr(err)
			}
		}
		effective := id
		if in.AmountCents != nil {
			if *in.AmountCents < 0 {
				return invalid("amount_cents", "out_of_range", "El monto no puede ser negativo.")
			}
			if effective, err = h.applyScope(ctx, tx, r, id, month, *in.AmountCents, scope); err != nil {
				return err
			}
		}
		nb, err := dominio.GetBudget(ctx, tx, effective)
		if err != nil {
			return err
		}
		out, err = h.budgetOut(ctx, tx, nb, dominio.MonthStart(month.Year, month.Month+1))
		return err
	})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// applyScope aplica el cambio de monto con el motor S2 (ApplyBudgetChange, K2) y guarda la diferencia de renglones.
func (h *handlers) applyScope(ctx context.Context, tx *sql.Tx, r *http.Request, id int64, month dominio.Date, amount int64, scope string) (int64, error) {
	b, err := dominio.GetBudget(ctx, tx, id)
	if err != nil {
		return 0, err
	}
	recs, err := dominio.BudgetsOfCategory(ctx, tx, b.CategoryID)
	if err != nil {
		return 0, err
	}
	var rows []mf.BudgetRow
	orig := map[int64]mf.BudgetRow{}
	var maxID int64
	for _, x := range recs {
		from, err := mf.ParseDate(x.ValidFrom)
		if err != nil {
			return 0, err
		}
		row := mf.BudgetRow{ID: x.ID, AmountCents: mf.Cents(x.AmountCents), ValidFrom: from}
		if x.ValidTo.Valid {
			t, err := mf.ParseDate(x.ValidTo.String)
			if err != nil {
				return 0, err
			}
			row.ValidTo = &t
		}
		rows, orig[x.ID] = append(rows, row), row
		maxID = max(maxID, x.ID)
	}
	res, err := mf.ApplyBudgetChange(rows, mf.BudgetChange{RowID: id, Month: dominio.EngineMonth(month), AmountCents: mf.Cents(amount), Scope: scope})
	var me *mf.Error
	if errors.As(err, &me) && me.Code == mf.CodeOutOfRange {
		return 0, invalid("month", "out_of_range", "El mes elegido está fuera de la vigencia de ese monto.")
	}
	if err != nil {
		return 0, err
	}
	actor, rid := api.Actor(r), api.RequestID(r)
	for _, d := range res.DeletedRowIDs {
		if err := dominio.BudgetAreaDelete(ctx, tx, "budget", d, actor, rid); err != nil {
			return 0, err
		}
	}
	toSQL := func(t *mf.Date) sql.NullString {
		if t == nil {
			return sql.NullString{}
		}
		return sql.NullString{String: t.String(), Valid: true}
	}
	ids := map[int64]int64{}
	// Primero los cambios a renglones existentes (acortan vigencias) y luego los nuevos, para no chocar con el índice único.
	for _, x := range res.Rows {
		o, ok := orig[x.ID]
		if !ok || (o.AmountCents == x.AmountCents && toSQL(o.ValidTo) == toSQL(x.ValidTo)) {
			continue
		}
		if err := dominio.BudgetAreaUpdate(ctx, tx, "budget", x.ID, []string{"amount_cents", "valid_to_date"},
			[]any{int64(x.AmountCents), toSQL(x.ValidTo)}, h.now(), actor, rid); err != nil {
			return 0, constraintErr(err)
		}
	}
	for _, x := range res.Rows {
		if x.ID <= maxID {
			continue
		}
		nid, err := dominio.BudgetAreaInsert(ctx, tx, "budget", []string{"category_id", "amount_cents", "valid_from_date", "valid_to_date"},
			[]any{b.CategoryID, int64(x.AmountCents), x.ValidFrom.String(), toSQL(x.ValidTo)}, actor, rid)
		if err != nil {
			return 0, constraintErr(err)
		}
		ids[x.ID] = nid
	}
	if n, ok := ids[res.EffectiveRowID]; ok {
		return n, nil
	}
	return res.EffectiveRowID, nil
}

func (h *handlers) deleteBudget(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "budget_id")
	if err != nil {
		return err
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetBudget(ctx, tx, id); err != nil {
			return err
		}
		return dominio.BudgetAreaDelete(ctx, tx, "budget", id, api.Actor(r), api.RequestID(r))
	}); err != nil {
		return err
	}
	return api.NoContent(w)
}

// ---------- Resúmenes del mes ----------

func (h *handlers) getBudgetMonth(w http.ResponseWriter, r *http.Request) error {
	m, err := h.monthParam(r)
	if err != nil {
		return err
	}
	ctx := r.Context()
	fe, in, err := h.engine(ctx)
	if err != nil {
		return err
	}
	bm, err := mf.BudgetMonth(in, dominio.EngineMonth(m))
	if err != nil {
		return err
	}
	inc, err := h.incomeMonth(ctx, m)
	if err != nil {
		return err
	}
	ref := h.today()
	if dominio.EngineMonth(ref) != dominio.EngineMonth(m) {
		ref = m
	}
	sin, err := fe.At(ref, nil)
	if err != nil {
		return err
	}
	sv, err := mf.SavingsGoals(sin)
	if err != nil {
		return err
	}
	out := api.BudgetMonth{Month: dominio.EngineMonth(m).String(), IncomeCents: inc.TotalCents, PlannedSpendingCents: cents(bm.Total.PlanCents),
		PlannedSavingsCents: cents(sv.MonthSummary.PlannedThisMonthCents), ActualCents: cents(bm.Total.ActualCents),
		RemainingCents: cents(bm.Total.RemainingCents), DaysLeft: int64(bm.DaysLeft)}
	out.UnassignedCents = out.IncomeCents - out.PlannedSpendingCents - out.PlannedSavingsCents
	if bm.Total.UsedPct != nil {
		out.UsedPct, out.Status = api.Ptr(*bm.Total.UsedPct), api.Ptr(bm.Total.State)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) getBudgetByCategory(w http.ResponseWriter, r *http.Request) error {
	m, err := h.monthParam(r)
	if err != nil {
		return err
	}
	ctx := r.Context()
	fe, in, err := h.engine(ctx)
	if err != nil {
		return err
	}
	bm, err := mf.BudgetMonth(in, dominio.EngineMonth(m))
	if err != nil {
		return err
	}
	counts, err := dominio.ExpenseCountByCategory(ctx, h.db(), m.String(), dominio.MonthEnd(m).String())
	if err != nil {
		return err
	}
	rows := []api.BudgetCategoryRowsRowsItem{}
	for _, l := range bm.Categories {
		id, ok := fe.CategoryID[l.Category]
		if !ok {
			continue
		}
		it := api.BudgetCategoryRowsRowsItem{CategoryID: id, PlannedCents: centsPtr(l.PlanCents), ActualCents: cents(l.ActualCents),
			DifferenceCents: centsPtr(l.RemainingCents), ExpenseCount: counts[id]}
		if l.UsedPct != nil {
			it.UsedPct, it.Status = api.Ptr(*l.UsedPct), api.Ptr(l.State)
		}
		rows = append(rows, it)
	}
	sort.SliceStable(rows, func(a, b int) bool {
		pa, pb := rows[a].UsedPct, rows[b].UsedPct
		if (pa == nil) != (pb == nil) {
			return pa != nil
		}
		if pa != nil && *pa != *pb {
			return *pa > *pb
		}
		return rows[a].CategoryID < rows[b].CategoryID
	})
	return api.WriteJSON(w, http.StatusOK, api.BudgetCategoryRows{Month: dominio.EngineMonth(m).String(), Rows: rows})
}

func (h *handlers) getBudgetFortnights(w http.ResponseWriter, r *http.Request) error {
	m, err := h.monthParam(r)
	if err != nil {
		return err
	}
	fe, in, err := h.engine(r.Context())
	if err != nil {
		return err
	}
	bm, err := mf.BudgetMonth(in, dominio.EngineMonth(m))
	if err != nil {
		return err
	}
	q1 := dominio.FortnightOf(m)
	q2 := q1.Next()
	wins := []mf.Window{{Label: "q1", Start: dominio.EngineDate(q1.Start()), End: dominio.EngineDate(q1.End())},
		{Label: "q2", Start: dominio.EngineDate(q2.Start()), End: dominio.EngineDate(q2.End())}}
	ps, err := mf.Periods(in, mf.PeriodSpec{Windows: wins})
	if err != nil {
		return err
	}
	// PROVISIONAL (bloqueo S6-presupuesto-quincenas): el motor S2 no reparte el plan mensual por quincena; se reparte
	// por días con ROUND_HALF_UP en la Q1 y la Q2 lleva el resto, para que las dos sumen el plan del mes.
	plan := cents(bm.Total.PlanCents)
	days := int64(dominio.MonthEnd(m).Day)
	p1 := dominio.MulDivRound(plan, 15, days)
	out := api.BudgetFortnights{Month: dominio.EngineMonth(m).String()}
	for i, f := range []dominio.Fortnight{q1, q2} {
		var actual int64
		if i < len(ps) {
			actual = cents(ps[i].Accrued.ByPerson.Get(fe.Holder))
		}
		pl := p1
		if i == 1 {
			pl = plan - p1
		}
		out.Fortnights = append(out.Fortnights, api.BudgetFortnightsFortnightsItem{Half: int64(f.Half), StartDate: f.Start().String(),
			EndDate: f.End().String(), PayDate: f.PayDay().String(), ActualCents: actual, PlannedCents: pl})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) getAccruedGroceries(w http.ResponseWriter, r *http.Request) error {
	m, err := h.monthParam(r)
	if err != nil {
		return err
	}
	fe, err := dominio.LoadFinEngine(r.Context(), h.db())
	if err != nil {
		return err
	}
	em := dominio.EngineMonth(m)
	from := dominio.EngineDate(dominio.MonthStart(m.Year, m.Month-24))
	rows := []api.AccruedGroceriesRowsItem{}
	for _, e := range fe.Expenses() {
		if e.Category != dominio.GroceriesCategoryCode || e.ExpenseDate.After(em.Last()) || e.ExpenseDate.Before(from) {
			continue
		}
		in, err := fe.Only(h.today(), e.ID)
		if err != nil {
			return fmt.Errorf("motor_fin: entrada inválida: %w", err)
		}
		ps, err := mf.Periods(in, mf.PeriodSpec{Months: []mf.Month{em}})
		if err != nil {
			return err
		}
		if len(ps) == 0 || ps[0].Accrued.TotalCents == 0 {
			continue
		}
		acc := ps[0].Accrued
		it := api.AccruedGroceriesRowsItem{ExpenseID: e.ID, CreditPurchaseID: e.CreditPurchaseID, ExpenseDate: e.ExpenseDate.String(),
			TotalCents: cents(e.AmountCents), AccruedCents: cents(acc.TotalCents), HouseholdShareCents: cents(acc.ByPerson.Get(fe.Holder))}
		if e.CreditPurchaseID != nil {
			if cp, ok := fe.Purchase(*e.CreditPurchaseID); ok {
				it.TotalCents = cents(cp.TotalCents)
			}
		}
		for _, x := range fe.Externals {
			it.ExternalShareCents += cents(acc.ByPerson.Get(x))
		}
		rows = append(rows, it)
	}
	return api.WriteJSON(w, http.StatusOK, api.AccruedGroceries{Month: em.String(), Rows: rows})
}

func (h *handlers) getBudgetAlerts(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	fe, in, err := h.engine(ctx)
	if err != nil {
		return err
	}
	today := h.today()
	em := dominio.EngineMonth(today)
	bm, err := mf.BudgetMonth(in, em)
	if err != nil {
		return err
	}
	type line struct {
		cat  *int64
		used int64
		thr  int64
		key  string
	}
	var active []line
	pick := func(key string, cat *int64, used *int64) {
		if used == nil || *used < fe.WarningPct {
			return
		}
		thr := fe.WarningPct
		if *used >= fe.LimitPct {
			thr = fe.LimitPct
		}
		active = append(active, line{cat, *used, thr, key})
	}
	for _, l := range bm.Categories {
		if id, ok := fe.CategoryID[l.Category]; ok {
			pick(l.Category, api.Ptr(id), l.UsedPct)
		}
	}
	pick("", nil, bm.Total.UsedPct)
	rows := []api.BudgetAlertsRowsItem{}
	// Fecha de cruce: el primer día del mes en que el motor, con los gastos hasta ese día, da % ≥ umbral.
	crossed := map[string]string{}
	for d := dominio.MonthStart(today.Year, today.Month); len(active) > 0 && !d.After(today) && len(crossed) < len(active); d = d.AddDays(1) {
		dd := d
		din, err := fe.At(dd, &dd)
		if err != nil {
			return err
		}
		dm, err := mf.BudgetMonth(din, em)
		if err != nil {
			return err
		}
		used := map[string]*int64{"": dm.Total.UsedPct}
		for _, l := range dm.Categories {
			used[l.Category] = l.UsedPct
		}
		for _, a := range active {
			if _, ok := crossed[a.key]; !ok && used[a.key] != nil && *used[a.key] >= a.thr {
				crossed[a.key] = d.String()
			}
		}
	}
	for _, a := range active {
		cd, ok := crossed[a.key]
		if !ok {
			cd = today.String()
		}
		rows = append(rows, api.BudgetAlertsRowsItem{CategoryID: a.cat, UsedPct: a.used, ThresholdPct: a.thr, CrossedDate: cd})
	}
	return api.WriteJSON(w, http.StatusOK, api.BudgetAlerts{Rows: rows})
}

func (h *handlers) getBudgetMonthlyTrend(w http.ResponseWriter, r *http.Request) error {
	cats, err := api.QueryIDs(r, "category_id")
	if err != nil {
		return err
	}
	if len(cats) > 1 {
		return invalid("category_id", "not_allowed", "Elige una sola categoría o ninguna para ver el total.")
	}
	fe, in, err := h.engine(r.Context())
	if err != nil {
		return err
	}
	out := api.BudgetMonthlyTrend{Months: []api.BudgetMonthlyTrendMonthsItem{}}
	code := ""
	if len(cats) == 1 {
		c, ok := fe.CategoryCode[cats[0]]
		if !ok {
			return invalid("category_id", "not_found", "No encontramos esa categoría.")
		}
		code, out.CategoryID = c, api.Ptr(cats[0])
	}
	// Supuesto: los 12 meses que terminan en el mes en curso (C1 dice «máximo 12 meses»).
	rows, err := mf.PlanVsActual(in, dominio.EngineMonth(h.today()).Add(-11), 12, code)
	var me *mf.Error
	if errors.As(err, &me) && me.Code == mf.CodeInvalidArgument && code != "" {
		return invalid("category_id", "not_allowed", "Esa categoría no está activa.")
	}
	if err != nil {
		return err
	}
	for _, x := range rows {
		out.Months = append(out.Months, api.BudgetMonthlyTrendMonthsItem{Month: x.Month.String(), PlannedCents: centsPtr(x.PlanCents), ActualCents: cents(x.ActualCents)})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func timeMonth(i int) time.Month { return time.Month(i) }
