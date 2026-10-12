package comidas

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/recetas"
)

var (
	sexes      = []string{"female", "male"}
	activities = []string{"sedentary", "light", "moderate", "active", "very_active"}
	goals      = []string{"lose", "maintain", "gain"}
	healthKeys = []string{"health.goal_adjust_lose_kcal", "health.goal_adjust_gain_kcal"}
)

func badEnum(field string) error {
	return &api.Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "El valor de " + field + " no es una de las opciones permitidas.",
		Fields: []api.FieldError{{Field: field, Code: "invalid_format", Message: "Opción no válida."}}}
}

// apiProfile convierte el perfil y pide al motor S4 edad, IMC, clase y meta sugerida (RN-17, K7, K23, K36).
func (m *motor) apiProfile(h dominio.HealthProfileRow) (api.HealthProfile, error) {
	lose, gain := m.e.Settings.HouseAdjustLoseKcal, m.e.Settings.HouseAdjustGainKcal
	out := api.HealthProfile{ID: h.ID, PersonID: h.PersonID, HeightMm: nint(h.HeightMM), BirthDate: nstr(h.BirthDate), Sex: nstr(h.Sex),
		ActivityLevel: nstr(h.ActivityLevel), DailyKcalTarget: nint(h.DailyKcalTarget), DietaryRestrictions: nstr(h.DietaryRestrictions),
		Notes: nstr(h.Notes), CreatedAt: api.APITime(h.CreatedAt), UpdatedAt: api.APITime(h.UpdatedAt), Goal: nstr(h.Goal),
		GoalAdjustLoseKcal: nint(h.GoalAdjustLoseKcal), GoalAdjustGainKcal: nint(h.GoalAdjustGainKcal),
		HouseholdGoalAdjustLoseKcal: &lose, HouseholdGoalAdjustGainKcal: &gain}
	pc := m.personCode(h.PersonID)
	var last *recetas.WeightEntry
	for i := range m.e.Weights {
		if x := &m.e.Weights[i]; x.Person == pc && (last == nil || !x.Date.Before(last.Date)) {
			last = x
		}
	}
	if last != nil {
		g, d := last.WeightG, last.Date.String()
		out.CurrentWeightG, out.CurrentWeightDate = &g, &d
	}
	hl, err := m.e.HealthFor(pc)
	if err != nil {
		return out, err
	}
	out.AgeYears, out.BMITenths, out.BMIClass, out.SuggestedKcalTarget = hl.Age, hl.BMITenths, hl.BMIClass, hl.SuggestedKcalTarget
	return out, nil
}

func (a *area) listHealthProfiles(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"person_id"}, "person_id")
	if err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	rows, err := dominio.ListHealthProfileRows(ctx, db, p.OrderBy(map[string]string{"person_id": "hp.person_id", "id": "hp.id"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	m, err := a.loadMotor(ctx, db, loadOpts{})
	if err != nil {
		return err
	}
	items := make([]api.HealthProfile, 0, len(rows))
	for _, h := range rows {
		x, err := m.apiProfile(h)
		if err != nil {
			return err
		}
		items = append(items, x)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

// checkHealth valida rangos (esquema 1.4) y opciones del perfil.
func checkHealth(v *api.Validator, h dominio.HealthProfileRow, today dominio.Date) error {
	for _, e := range []struct {
		val   sql.NullString
		opts  []string
		field string
	}{{h.Sex, sexes, "sex"}, {h.ActivityLevel, activities, "activity_level"}, {h.Goal, goals, "goal"}} {
		if e.val.Valid && !slices.Contains(e.opts, e.val.String) {
			return badEnum(e.field)
		}
	}
	if h.BirthDate.Valid {
		d, err := v.Date(h.BirthDate.String, "birth_date")
		if err != nil {
			return err
		}
		v.Check(!d.After(today), "birth_date", "out_of_range", "La fecha de nacimiento no puede ser futura.")
	}
	v.Check(!h.HeightMM.Valid || (h.HeightMM.Int64 >= 500 && h.HeightMM.Int64 <= 2500), "height_mm", "out_of_range", "La estatura va de 500 a 2500 mm (0.50 a 2.50 m).")
	v.Check(!h.DailyKcalTarget.Valid || (h.DailyKcalTarget.Int64 >= 500 && h.DailyKcalTarget.Int64 <= 10000), "daily_kcal_target", "out_of_range", "La meta va de 500 a 10000 kcal por día.")
	v.Check(!h.GoalAdjustLoseKcal.Valid || h.GoalAdjustLoseKcal.Int64 >= 0, "goal_adjust_lose_kcal", "out_of_range", "El ajuste es un entero de 0 o más.")
	v.Check(!h.GoalAdjustGainKcal.Valid || h.GoalAdjustGainKcal.Int64 >= 0, "goal_adjust_gain_kcal", "out_of_range", "El ajuste es un entero de 0 o más.")
	return v.Err()
}

func (a *area) profileResponse(ctx context.Context, id int64) (api.HealthProfile, error) {
	m, err := a.loadMotor(ctx, a.db(), loadOpts{})
	if err != nil {
		return api.HealthProfile{}, err
	}
	h, err := dominio.GetHealthProfileRow(ctx, a.db(), id)
	if err != nil {
		return api.HealthProfile{}, err
	}
	return m.apiProfile(h)
}

func (a *area) createHealthProfile(w http.ResponseWriter, r *http.Request) error {
	var in api.HealthProfileCreate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{})
	if err != nil {
		return err
	}
	var v api.Validator
	if p, ok := m.person(in.PersonID); !ok {
		v.Add("person_id", "not_found", "Esa persona no existe.")
	} else if !p.HasMealPlan() {
		return externalPerson("person_id")
	}
	h := dominio.HealthProfileRow{PersonID: in.PersonID, HeightMM: toNI(in.HeightMm), BirthDate: toNS(in.BirthDate), Sex: toNS(in.Sex),
		ActivityLevel: toNS(in.ActivityLevel), DailyKcalTarget: toNI(in.DailyKcalTarget), DietaryRestrictions: toNS(in.DietaryRestrictions),
		Notes: toNS(in.Notes), Goal: toNS(in.Goal), GoalAdjustLoseKcal: toNI(in.GoalAdjustLoseKcal), GoalAdjustGainKcal: toNI(in.GoalAdjustGainKcal)}
	if err := checkHealth(&v, h, a.today()); err != nil {
		return err
	}
	if slices.ContainsFunc(m.profiles, func(x dominio.HealthProfileRow) bool { return x.PersonID == in.PersonID }) {
		return dominio.Conflict("already_exists", "Esa persona ya tiene perfil de salud. Edítalo en lugar de crear otro.")
	}
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if err := dominio.SaveHealthProfile(ctx, tx, &h); err != nil {
			if dominio.IsUnique(err) {
				return dominio.Conflict("already_exists", "Esa persona ya tiene perfil de salud. Edítalo en lugar de crear otro.")
			}
			return err
		}
		return a.audit(ctx, tx, r, "insert", "health_profile", h.ID, nil)
	}); err != nil {
		return err
	}
	out, err := a.profileResponse(ctx, h.ID)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

func (a *area) updateHealthProfile(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "health_profile_id")
	if err != nil {
		return err
	}
	var in api.HealthProfileUpdate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if !in.HeightMm.Set && !in.BirthDate.Set && !in.Sex.Set && !in.ActivityLevel.Set && !in.DailyKcalTarget.Set && !in.DietaryRestrictions.Set &&
		!in.Notes.Set && !in.Goal.Set && !in.GoalAdjustLoseKcal.Set && !in.GoalAdjustGainKcal.Set {
		return dominio.Invalid(issue("body", "required", "Manda al menos un cambio."))
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, a.db(), func(tx *sql.Tx) error {
		h, err := dominio.GetHealthProfileRow(ctx, tx, id)
		if err != nil {
			return err
		}
		before, err := dominio.RowJSON(ctx, tx, "health_profile", id)
		if err != nil {
			return err
		}
		h.HeightMM, h.BirthDate, h.Sex = mergeI(in.HeightMm, h.HeightMM), mergeS(in.BirthDate, h.BirthDate), mergeS(in.Sex, h.Sex)
		h.ActivityLevel, h.DailyKcalTarget = mergeS(in.ActivityLevel, h.ActivityLevel), mergeI(in.DailyKcalTarget, h.DailyKcalTarget)
		h.DietaryRestrictions, h.Notes, h.Goal = mergeS(in.DietaryRestrictions, h.DietaryRestrictions), mergeS(in.Notes, h.Notes), mergeS(in.Goal, h.Goal)
		h.GoalAdjustLoseKcal, h.GoalAdjustGainKcal = mergeI(in.GoalAdjustLoseKcal, h.GoalAdjustLoseKcal), mergeI(in.GoalAdjustGainKcal, h.GoalAdjustGainKcal)
		var v api.Validator
		if err := checkHealth(&v, h, a.today()); err != nil {
			return err
		}
		if err := dominio.SaveHealthProfile(ctx, tx, &h); err != nil {
			return err
		}
		return a.audit(ctx, tx, r, "update", "health_profile", id, before)
	}); err != nil {
		return err
	}
	out, err := a.profileResponse(ctx, id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// getTodayKcal: kcal planeadas hoy por persona contra su meta (motor S4, semana que contiene hoy).
func (a *area) getTodayKcal(w http.ResponseWriter, r *http.Request) error {
	today := a.today()
	mo := mondayOf(today)
	m, err := a.loadMotor(r.Context(), a.db(), loadOpts{food: true, from: mo, to: mo.AddDays(6)})
	if err != nil {
		return err
	}
	wp, err := m.e.WeekPlanFor(rdate(mo.String()))
	if err != nil {
		return err
	}
	f := fillWeek(wp)
	out := api.TodayKcal{Date: today.String(), Persons: []api.TodayKcalPersonsItem{}}
	for _, d := range wp.Days {
		if d.Date.String() != today.String() {
			continue
		}
		for _, p := range d.Persons {
			it := api.TodayKcalPersonsItem{PersonID: m.personID[p.Person], PlannedKcal: f.kcal[d.Date.String()+"|"+p.Person],
				DailyKcalTarget: p.DailyKcalTarget, DifferenceKcal: p.DifferenceKcal}
			for _, c := range wp.Cells {
				if c.EatingOut && c.Person == p.Person && c.Date.String() == today.String() {
					it.EatingOutSlots++
				}
			}
			out.Persons = append(out.Persons, it)
		}
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// apiWeight convierte el peso; el IMC a esa fecha lo calcula el motor S4 con ese peso y el perfil de la persona.
func (m *motor) apiWeight(x dominio.WeightEntryRow) api.WeightEntry {
	out := api.WeightEntry{ID: x.ID, PersonID: x.PersonID, MeasuredDate: x.MeasuredDate, WeightG: x.WeightG, Notes: nstr(x.Notes),
		CreatedAt: api.APITime(x.CreatedAt), UpdatedAt: api.APITime(x.UpdatedAt)}
	if x.CapturedAt.Valid {
		c := api.APITime(x.CapturedAt.String)
		out.CapturedAt = &c
	}
	pc := m.personCode(x.PersonID)
	e := &recetas.Engine{Today: rdate(x.MeasuredDate), Persons: m.e.Persons, Profiles: m.e.Profiles, Settings: m.e.Settings,
		Weights: []recetas.WeightEntry{{Person: pc, Date: rdate(x.MeasuredDate), WeightG: x.WeightG}}}
	if h, err := e.HealthFor(pc); err == nil {
		out.BMITenths = h.BMITenths
	}
	return out
}

func (a *area) listWeightEntries(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"measured_date", "weight_g"}, "-measured_date")
	if err != nil {
		return err
	}
	ids, err := api.QueryIDs(r, "person_id")
	if err != nil {
		return err
	}
	per, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	f := dominio.WeightFilter{PersonIDs: ids}
	if per != nil {
		f.From, f.To = per.Start.String(), per.End.String()
	}
	ctx, db := r.Context(), a.db()
	rows, err := dominio.ListWeightEntryRows(ctx, db, f, p.OrderBy(map[string]string{"measured_date": "w.measured_date", "weight_g": "w.weight_g", "id": "w.id"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	m, err := a.loadMotor(ctx, db, loadOpts{})
	if err != nil {
		return err
	}
	items := make([]api.WeightEntry, 0, len(rows))
	for _, x := range rows {
		items = append(items, m.apiWeight(x))
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

func checkWeight(v *api.Validator, g int64) {
	v.Check(g >= 20000 && g <= 400000, "weight_g", "out_of_range", "El peso va de 20 a 400 kg (en gramos, de 20000 a 400000).")
}

func weightUnique(err error) error {
	if dominio.IsUnique(err) {
		return dominio.Conflict("already_exists", "Ya hay un peso de esa persona en ese día. Corrígelo en lugar de registrar otro.")
	}
	return err
}

// createWeightEntry: K9. La clave queda en weight_entry.idempotency_key; un reintento con la misma clave pasadas
// las 24 h de la tabla idempotency_key responde 200 con el mismo peso e Idempotent-Replayed: true.
func (a *area) createWeightEntry(w http.ResponseWriter, r *http.Request) error {
	var in api.WeightEntryCreate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ctx, db := r.Context(), a.db()
	m, err := a.loadMotor(ctx, db, loadOpts{})
	if err != nil {
		return err
	}
	key := api.IdempotencyKey(r)
	if key != "" {
		old, ok, err := dominio.WeightEntryByKey(ctx, db, key)
		if err != nil {
			return err
		}
		if ok {
			api.MarkReplayed(w)
			return api.WriteJSON(w, http.StatusOK, m.apiWeight(old))
		}
	}
	var v api.Validator
	d, err := v.Date(in.MeasuredDate, "measured_date")
	if err != nil {
		return err
	}
	row := dominio.WeightEntryRow{PersonID: in.PersonID, MeasuredDate: d.String(), WeightG: in.WeightG, Notes: toNS(in.Notes),
		IdempotencyKey: sql.NullString{String: key, Valid: key != ""}}
	if in.CapturedAt != nil {
		t, err := dominio.ParseDateTime(*in.CapturedAt)
		if err != nil {
			return &api.Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "captured_at debe ser fecha y hora ISO 8601.",
				Fields: []api.FieldError{{Field: "captured_at", Code: "invalid_format", Message: "Fecha y hora inválidas."}}}
		}
		row.CapturedAt = sql.NullString{String: dominio.FormatDBTime(t), Valid: true}
	}
	if p, ok := m.person(in.PersonID); !ok {
		v.Add("person_id", "not_found", "Esa persona no existe.")
	} else if !p.HasMealPlan() {
		return externalPerson("person_id")
	}
	checkWeight(&v, in.WeightG)
	if err := v.Err(); err != nil {
		return err
	}
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if err := dominio.SaveWeightEntry(ctx, tx, &row); err != nil {
			return weightUnique(err)
		}
		if err := a.audit(ctx, tx, r, "insert", "weight_entry", row.ID, nil); err != nil {
			return err
		}
		row, err = dominio.GetWeightEntryRow(ctx, tx, row.ID)
		return err
	}); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, m.apiWeight(row))
}

func (a *area) updateWeightEntry(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "weight_entry_id")
	if err != nil {
		return err
	}
	var in api.WeightEntryUpdate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if in.MeasuredDate == nil && in.WeightG == nil && !in.Notes.Set {
		return dominio.Invalid(issue("body", "required", "Manda al menos un cambio."))
	}
	ctx, db := r.Context(), a.db()
	var row dominio.WeightEntryRow
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		if row, err = dominio.GetWeightEntryRow(ctx, tx, id); err != nil {
			return err
		}
		before, err := dominio.RowJSON(ctx, tx, "weight_entry", id)
		if err != nil {
			return err
		}
		var v api.Validator
		if in.MeasuredDate != nil {
			d, err := v.Date(*in.MeasuredDate, "measured_date")
			if err != nil {
				return err
			}
			row.MeasuredDate = d.String()
		}
		if in.WeightG != nil {
			row.WeightG = *in.WeightG
		}
		row.Notes = mergeS(in.Notes, row.Notes)
		checkWeight(&v, row.WeightG)
		if err := v.Err(); err != nil {
			return err
		}
		if err := dominio.SaveWeightEntry(ctx, tx, &row); err != nil {
			return weightUnique(err)
		}
		if err := a.audit(ctx, tx, r, "update", "weight_entry", id, before); err != nil {
			return err
		}
		row, err = dominio.GetWeightEntryRow(ctx, tx, id)
		return err
	}); err != nil {
		return err
	}
	m, err := a.loadMotor(ctx, db, loadOpts{})
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, m.apiWeight(row))
}

func (a *area) deleteWeightEntry(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "weight_entry_id")
	if err != nil {
		return err
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, a.db(), func(tx *sql.Tx) error {
		if _, err := dominio.GetWeightEntryRow(ctx, tx, id); err != nil {
			return err
		}
		before, err := dominio.RowJSON(ctx, tx, "weight_entry", id)
		if err != nil {
			return err
		}
		if err := dominio.DeleteWeightEntry(ctx, tx, id); err != nil {
			return err
		}
		return a.audit(ctx, tx, r, "delete", "weight_entry", id, before)
	}); err != nil {
		return err
	}
	return api.NoContent(w)
}

// ---- Ajuste de kcal de la casa (K36, K44, C1 §11.7) ----

func rawInt(raw json.RawMessage) api.Nullable[int64] {
	var v *int64
	if len(raw) == 0 || json.Unmarshal(raw, &v) != nil || v == nil {
		return api.Null[int64]()
	}
	return api.Some(*v)
}

func (a *area) healthSettings(ctx context.Context, q dominio.Querier) (api.HealthSettings, error) {
	var out api.HealthSettings
	for i, dst := range []*api.SettingParameterInteger{&out.Parameters.HealthGoalAdjustLoseKcal, &out.Parameters.HealthGoalAdjustGainKcal} {
		p, err := dominio.ReadParameter(ctx, q, healthKeys[i])
		if err != nil {
			return out, err
		}
		*dst = api.SettingParameterInteger{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, Value: rawInt(p.Value), InitialValue: rawInt(p.InitialValue), IsInitial: p.IsInitial}
	}
	return out, nil
}

func (a *area) getHealthSettings(w http.ResponseWriter, r *http.Request) error {
	out, err := a.healthSettings(r.Context(), a.db())
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (a *area) updateHealthSettings(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Parameters      map[string]json.RawMessage `json:"parameters,omitempty"`
		ResetParameters []string                   `json:"reset_parameters,omitempty"`
	}
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if len(in.Parameters) == 0 && len(in.ResetParameters) == 0 {
		return dominio.Invalid(issue("parameters", "required", "Manda al menos un parámetro para cambiar o restablecer."))
	}
	ctx := r.Context()
	if err := dominio.WithTx(ctx, a.db(), func(tx *sql.Tx) error {
		return dominio.ApplyParams(ctx, tx, dominio.ParamChange{Set: in.Parameters, Reset: in.ResetParameters}, healthKeys, api.Actor(r), api.RequestID(r))
	}); err != nil {
		return err
	}
	return a.getHealthSettings(w, r)
}
