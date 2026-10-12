package ajustes

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/motor_fin"
)

// goalFigures pide al motor S2 (motor_fin.SavingsGoals, §3.13) el ahorrado, el avance, la aportación necesaria y el
// estado. El motor solo informa metas activas: una meta desactivada sale sin campos calculados.
func goalFigures(ctx context.Context, reg *api.Registry, g dominio.SavingsGoalRow) (*motor_fin.GoalStatus, error) {
	db := reg.Deps().DB
	ps, err := dominio.ListPersonRecords(ctx, db, 0, "ORDER BY p.id")
	if err != nil {
		return nil, err
	}
	td, err := motor_fin.ParseDate(today(reg).String())
	if err != nil {
		return nil, err
	}
	// groceries_category y paper_category son obligatorios en Params pero SavingsGoals no los usa.
	in := motor_fin.Input{Params: motor_fin.Params{Today: &td, Zone: motor_fin.DefaultZone, GroceriesCategory: "groceries", PaperCategory: "paper"}}
	for _, p := range ps {
		in.Persons = append(in.Persons, motor_fin.Person{Code: p.Code, DisplayName: p.DisplayName, Participation: p.Participation, PaysShareDirectly: p.PaysShareDirectly})
		if p.Participation == "holder" {
			in.Params.Holder = p.Code
		}
	}
	mg := motor_fin.SavingsGoal{ID: g.ID, Name: g.Name, Kind: g.Kind, IsActive: api.Ptr(g.IsActive)}
	if mg.StartDate, err = motor_fin.ParseDate(g.StartDate); err != nil {
		return nil, err
	}
	if g.TargetCents.Valid {
		mg.TargetCents = api.Ptr(motor_fin.Cents(g.TargetCents.Int64))
	}
	if g.PlannedMonthlyCents.Valid {
		mg.PlannedMonthlyCents = api.Ptr(motor_fin.Cents(g.PlannedMonthlyCents.Int64))
	}
	if g.TargetDate.Valid {
		d, err := motor_fin.ParseDate(g.TargetDate.String)
		if err != nil {
			return nil, err
		}
		mg.TargetDate = &d
	}
	in.SavingsGoals = []motor_fin.SavingsGoal{mg}
	mv, err := dominio.SavingsMovementsOf(ctx, db, g.ID)
	if err != nil {
		return nil, err
	}
	for _, m := range mv {
		d, err := motor_fin.ParseDate(m.MovementDate)
		if err != nil {
			return nil, err
		}
		in.SavingsMovements = append(in.SavingsMovements, motor_fin.SavingsMovement{GoalID: g.ID, MovementDate: d, Kind: m.Kind, AmountCents: motor_fin.Cents(m.AmountCents)})
	}
	res, err := motor_fin.SavingsGoals(in)
	if err != nil {
		return nil, err
	}
	for i := range res.Goals {
		if res.Goals[i].GoalID == g.ID {
			return &res.Goals[i], nil
		}
	}
	return nil, nil
}

func goalOut(ctx context.Context, reg *api.Registry, g dominio.SavingsGoalRow) (api.SavingsGoal, error) {
	out := api.SavingsGoal{ID: g.ID, Name: g.Name, Kind: g.Kind, TargetCents: pInt(g.TargetCents), TargetDate: pStr(g.TargetDate),
		PlannedMonthlyCents: pInt(g.PlannedMonthlyCents), StartDate: g.StartDate, IsActive: g.IsActive, Notes: pStr(g.Notes),
		CreatedAt: api.APITime(g.CreatedAt), UpdatedAt: api.APITime(g.UpdatedAt)}
	st, err := goalFigures(ctx, reg, g)
	if err != nil || st == nil {
		return out, err
	}
	out.SavedCents = api.Ptr(int64(st.SavedCents))
	out.ProgressTenthsPct, out.ProgressPct = st.ProgressTenthsPct, st.ProgressPct
	if st.NeededMonthlyCents != nil {
		out.RequiredMonthlyCents = api.Ptr(int64(*st.NeededMonthlyCents))
	}
	out.Status = api.Ptr(st.State)
	return out, nil
}

func checkGoal(g dominio.SavingsGoalRow) error {
	if !oneOf(g.Kind, "fund", "goal") {
		return badEnum("kind")
	}
	var v api.Validator
	for _, d := range []struct {
		s     string
		field string
		ok    bool
	}{{g.StartDate, "start_date", true}, {g.TargetDate.String, "target_date", g.TargetDate.Valid}} {
		if d.ok {
			if _, err := v.Date(d.s, d.field); err != nil {
				return err
			}
		}
	}
	v.Required(strings.TrimSpace(g.Name) != "", "name")
	if g.Kind == "goal" {
		v.Check(g.TargetCents.Valid, "target_cents", "required", "Una meta lleva monto objetivo.")
		v.Check(g.TargetDate.Valid, "target_date", "required", "Una meta lleva fecha objetivo.")
	}
	v.Check(!g.TargetCents.Valid || g.TargetCents.Int64 > 0, "target_cents", "out_of_range", "El monto objetivo debe ser mayor que cero.")
	v.Check(!g.PlannedMonthlyCents.Valid || g.PlannedMonthlyCents.Int64 >= 0, "planned_monthly_cents", "out_of_range", "La aportación no puede ser negativa.")
	v.Check(!g.TargetDate.Valid || g.TargetDate.String > g.StartDate, "target_date", "out_of_range", "La fecha objetivo debe ser posterior al inicio.")
	return v.Err()
}

const dupGoal = "Ya existe una meta o fondo con ese nombre."

func createSavingsGoal(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.SavingsGoalCreate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		var v api.Validator
		for _, k := range []string{"name", "kind", "start_date"} {
			_, ok := keys[k]
			v.Required(ok, k)
		}
		if err := v.Err(); err != nil {
			return err
		}
		g := dominio.SavingsGoalRow{Name: strings.TrimSpace(in.Name), Kind: in.Kind, TargetCents: nInt(in.TargetCents),
			TargetDate: nStr(in.TargetDate), PlannedMonthlyCents: nInt(in.PlannedMonthlyCents), StartDate: in.StartDate,
			IsActive: in.IsActive == nil || *in.IsActive, Notes: nStr(in.Notes)}
		if err := checkGoal(g); err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		var id int64
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			id, err = dominio.InsertOwned(ctx, tx, "savings_goal", g.Columns(), api.Actor(r), api.RequestID(r))
			return err
		})
		if err := writeErr(err, dupGoal); err != nil {
			return err
		}
		saved, err := dominio.SavingsGoalRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		out, err := goalOut(ctx, reg, saved)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, out)
	}
}

func updateSavingsGoal(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "savings_goal_id")
		if err != nil {
			return err
		}
		var in api.SavingsGoalUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return noChanges()
		}
		ctx, db := r.Context(), reg.Deps().DB
		g, err := dominio.SavingsGoalRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		if in.Name != nil {
			g.Name = strings.TrimSpace(*in.Name)
		}
		if in.Kind != nil {
			g.Kind = *in.Kind
		}
		if in.StartDate != nil {
			g.StartDate = *in.StartDate
		}
		if in.IsActive != nil {
			g.IsActive = *in.IsActive
		}
		g.TargetCents = mergeInt(g.TargetCents, in.TargetCents)
		g.TargetDate = mergeStr(g.TargetDate, in.TargetDate)
		g.PlannedMonthlyCents = mergeInt(g.PlannedMonthlyCents, in.PlannedMonthlyCents)
		g.Notes = mergeStr(g.Notes, in.Notes)
		if err := checkGoal(g); err != nil {
			return err
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			return dominio.UpdateOwned(ctx, tx, "savings_goal", id, g.Columns(), api.Actor(r), api.RequestID(r))
		})
		if err := writeErr(err, dupGoal); err != nil {
			return err
		}
		saved, err := dominio.SavingsGoalRowByID(ctx, db, id)
		if err != nil {
			return err
		}
		out, err := goalOut(ctx, reg, saved)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

const projectionKey = "projection.repeat_cadence_enabled"

func getProjectionSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		on, err := dominio.SettingBool(r.Context(), reg.Deps().DB, projectionKey)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, api.ProjectionSettings{RepeatCadenceEnabled: on})
	}
}

// updateProjectionSettings cambia el interruptor K21 (clave sin @param: dominio.SetSetting con audit_log).
func updateProjectionSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.ProjectionSettingsUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		if raw, ok := keys["repeat_cadence_enabled"]; !ok || string(raw) == "null" {
			return dominio.Invalid(dominio.FieldIssue{Field: "repeat_cadence_enabled", Code: "required", Message: "Indica si se enciende o se apaga."})
		}
		ctx, db := r.Context(), reg.Deps().DB
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			return dominio.SetSetting(ctx, tx, projectionKey, dominio.JSONBool(in.RepeatCadenceEnabled), api.Actor(r), api.RequestID(r))
		})
		if err != nil {
			return err
		}
		return getProjectionSettings(reg)(w, r)
	}
}

// alertKeys son las claves de «Notificaciones» (K44, K47, K50).
var alertKeys = []string{"alerts.budget_warning_pct", "alerts.budget_limit_pct", "alerts.payment_due_days_before",
	"alerts.send_time", "alerts.stale_prices_advance_days", "alerts.meal_plan_missing_weekday"}

func nullableOf[T any](raw json.RawMessage) (api.Nullable[T], error) {
	if len(raw) == 0 || string(raw) == "null" {
		return api.Null[T](), nil
	}
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return api.Nullable[T]{}, err
	}
	return api.Some(v), nil
}

func intParam(ctx context.Context, q dominio.Querier, key string) (api.SettingParameterInteger, error) {
	p, err := dominio.ReadParameter(ctx, q, key)
	if err != nil {
		return api.SettingParameterInteger{}, err
	}
	out := api.SettingParameterInteger{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, IsInitial: p.IsInitial}
	if out.Value, err = nullableOf[int64](p.Value); err != nil {
		return out, err
	}
	out.InitialValue, err = nullableOf[int64](p.InitialValue)
	return out, err
}

func readAlerts(ctx context.Context, q dominio.Querier) (api.AlertSettings, error) {
	var out api.AlertSettings
	ps := &out.Parameters
	var err error
	for key, dst := range map[string]*api.SettingParameterInteger{
		alertKeys[0]: &ps.AlertsBudgetWarningPct, alertKeys[1]: &ps.AlertsBudgetLimitPct,
		alertKeys[4]: &ps.AlertsStalePricesAdvanceDays, alertKeys[5]: &ps.AlertsMealPlanMissingWeekday,
	} {
		if *dst, err = intParam(ctx, q, key); err != nil {
			return out, err
		}
	}
	l, err := dominio.ReadParameter(ctx, q, alertKeys[2])
	if err != nil {
		return out, err
	}
	ps.AlertsPaymentDueDaysBefore = api.SettingParameterIntegerList{Key: l.Key, ValueType: l.ValueType, Unit: l.Unit, IsInitial: l.IsInitial}
	if ps.AlertsPaymentDueDaysBefore.Value, err = nullableOf[[]int64](l.Value); err != nil {
		return out, err
	}
	if ps.AlertsPaymentDueDaysBefore.InitialValue, err = nullableOf[[]int64](l.InitialValue); err != nil {
		return out, err
	}
	t, err := dominio.ReadParameter(ctx, q, alertKeys[3])
	if err != nil {
		return out, err
	}
	ps.AlertsSendTime = api.SettingParameterTime{Key: t.Key, ValueType: t.ValueType, Unit: t.Unit, IsInitial: t.IsInitial}
	if ps.AlertsSendTime.Value, err = nullableOf[string](t.Value); err != nil {
		return out, err
	}
	ps.AlertsSendTime.InitialValue, err = nullableOf[string](t.InitialValue)
	return out, err
}

func getAlertSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		out, err := readAlerts(r.Context(), reg.Deps().DB)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// updateAlertSettings: parameters (valor nuevo por clave) y reset_parameters (vuelven al inicial), con
// dominio.ApplyParams (tipo, clave de la sección, «cambiar y restablecer a la vez» y rango de los disparadores: 422).
func updateAlertSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.AlertSettingsUpdate
		keys, err := decode(reg, r, &in)
		if err != nil {
			return err
		}
		var raw struct {
			Parameters map[string]json.RawMessage `json:"parameters"`
		}
		_ = json.Unmarshal(mustJSON(keys), &raw)
		_, hasReset := keys["reset_parameters"]
		if len(keys) == 0 || (keys["parameters"] != nil && len(raw.Parameters) == 0) || (hasReset && len(in.ResetParameters) == 0) {
			return noChanges()
		}
		seen := map[string]bool{}
		for _, k := range in.ResetParameters {
			if seen[k] {
				return api.BadRequest("Una clave viene repetida en reset_parameters.")
			}
			seen[k] = true
		}
		ctx, db := r.Context(), reg.Deps().DB
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			return dominio.ApplyParams(ctx, tx, dominio.ParamChange{Set: raw.Parameters, Reset: in.ResetParameters}, alertKeys, api.Actor(r), api.RequestID(r))
		})
		if err != nil {
			return err
		}
		return getAlertSettings(reg)(w, r)
	}
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
