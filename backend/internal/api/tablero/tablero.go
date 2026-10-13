// Package tablero registra las cinco operaciones del tablero de control de Inicio que agrega C1 v1.5 (K54,
// U3 2.0): getDashboardDebtCapacity (B-1, M-20, RN-30), getDashboardMonthlyOverview (G-1 y G-6, M-21),
// getBudgetSignals (G-3, M-22, RN-31) y los parámetros del tablero de Ajustes (getDashboardSettings y
// updateDashboardSettings, K44). Todas las cifras salen de los motores aprobados: S2 (motor_fin, con la entrada de
// dominio.LoadFinEngine) para ingreso, devengado, efectivo y plan, y S3 (motor_cred, por creditos.PaymentsByMonth)
// para los pagos de crédito. Aquí solo se agrupa, se divide con ROUND_HALF_UP (dominio.PercentOf, MulDivRound) y se
// compara contra los umbrales vigentes de app_setting. Agente correccion-especificacion-r2 (hallazgo A-1 de
// REVISION-especificacion-r2; H-1 reabierto); pruebas con datos, escala compartida y regla «sin datos» de api-tablero.
package tablero

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/creditos"
	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
)

// Register registra las cinco operaciones del tablero (C1 v1.5).
func Register(reg *api.Registry) {
	reg.Handle("getDashboardDebtCapacity", getDashboardDebtCapacity(reg))
	reg.Handle("getDashboardMonthlyOverview", getDashboardMonthlyOverview(reg))
	reg.Handle("getBudgetSignals", getBudgetSignals(reg))
	reg.Handle("getDashboardSettings", getDashboardSettings(reg))
	reg.Handle("updateDashboardSettings", updateDashboardSettings(reg))
}

// Claves de app_setting del tablero (migración 00005, K54). La sección de Ajustes es PROVISIONAL (bloqueo C1-tablero).
const (
	keyWarning   = "debt.capacity_warning_pct"
	keyLimit     = "debt.capacity_limit_pct"
	keyReminder  = "dashboard.payment_reminder_count"
	keySignalN   = "budget.signal_months"
	keySignalPct = "budget.signal_margin_pct"
)

var dashboardKeys = []string{keyWarning, keyLimit, keyReminder, keySignalN, keySignalPct}

// maxOverviewMonths es el máximo de la ventana de getDashboardMonthlyOverview (C1 v1.5, OverviewFrom).
const maxOverviewMonths = 24

func today(reg *api.Registry) dominio.Date { return dominio.Today(reg.Deps().Clock) }

// engine carga la entrada del motor S2 con «hoy» = el reloj de la API.
func engine(reg *api.Registry, ctx context.Context) (*dominio.FinEngine, mf.Input, error) {
	fe, err := dominio.LoadFinEngine(ctx, reg.Deps().DB)
	if err != nil {
		return nil, mf.Input{}, err
	}
	in, err := fe.At(today(reg), nil)
	if err != nil {
		return nil, mf.Input{}, fmt.Errorf("motor_fin: entrada inválida: %w", err)
	}
	return fe, in, nil
}

// monthParam lee period: solo mes calendario; por omisión el mes en curso. Una quincena, una semana o un año
// responden 422 (PROVISIONAL, bloqueo C1-tablero, punto 5).
func monthParam(reg *api.Registry, r *http.Request) (mf.Month, error) {
	p, err := api.QueryPeriod(r)
	if err != nil { // 400 propio: aquí solo vale un mes, así que el ejemplo no propone una quincena (B-6 de REVISION-tablero-backend-r3)
		return mf.Month{}, api.BadRequest(fmt.Sprintf("El periodo «%s» no tiene un formato válido. Aquí va un mes con el formato AAAA-MM, por ejemplo 2026-10.", r.URL.Query().Get("period")))
	}
	if p == nil {
		return dominio.EngineMonth(today(reg)), nil
	}
	if p.Kind != dominio.PeriodMonth {
		return mf.Month{}, dominio.Invalid(dominio.FieldIssue{Field: "period", Code: "not_allowed",
			Message: "Aquí solo se puede pedir un mes completo, por ejemplo 2026-10."})
	}
	return dominio.EngineMonth(p.Start), nil
}

func monthIndex(m mf.Month) int { return m.Year*12 + int(m.Mon) - 1 }

// thresholds son los umbrales vigentes del barómetro (K54 decisión 2, K44).
func thresholds(ctx context.Context, q dominio.Querier) (warning, limit int64, err error) {
	if warning, err = dominio.SettingInt(ctx, q, keyWarning); err != nil {
		return 0, 0, err
	}
	limit, err = dominio.SettingInt(ctx, q, keyLimit)
	return warning, limit, err
}

// scaleMaxPct es el fin de la escala del dibujo del barómetro (scale_max_pct de M-20): el múltiplo de 10 igual o
// mayor que 1.5 × limit_pct, en enteros (60 con 40). PROVISIONAL, bloqueo U3-tablero-barometro. La comparten
// getDashboardDebtCapacity y getDashboardMonthlyOverview (C1 v1.5.1, hallazgo M-R3-1 de REVISION-especificacion-r3).
func scaleMaxPct(limit int64) int64 { return (3*limit + 19) / 20 * 10 }

// debtStatus aplica RN-30 al entero redondeado: fronteras inclusivas (30 es normal y 40 es warning).
func debtStatus(ratio, warning, limit int64) string {
	switch {
	case ratio <= warning:
		return "normal"
	case ratio <= limit:
		return "warning"
	}
	return "exceeded"
}

// monthIncome es el ingreso del mes según el flujo quincenal del motor S2 (las dos quincenas del mes, K8):
// quincenas = capturadas (fixed) + esperadas sin capturar; extras aparte.
type monthIncome struct {
	fortnights, extras int64
	expected           bool
}

func incomeOf(in mf.Input, m mf.Month) (monthIncome, error) {
	rows, err := mf.FortnightFlow(in, mf.FortnightOf(m.First()), 2)
	if err != nil {
		return monthIncome{}, err
	}
	var out monthIncome
	for _, x := range rows {
		out.fortnights += int64(x.Income.FixedCents) + int64(x.Income.ExpectedCents)
		out.extras += int64(x.Income.ExtraCents)
		out.expected = out.expected || x.Income.IncomeIsExpected
	}
	return out, nil
}

func getDashboardDebtCapacity(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		m, err := monthParam(reg, r)
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		warning, limit, err := thresholds(ctx, db)
		if err != nil {
			return err
		}
		_, in, err := engine(reg, ctx)
		if err != nil {
			return err
		}
		inc, err := incomeOf(in, m)
		if err != nil {
			return err
		}
		// Qué líneas cuentan (RN-30 «de todas las líneas») y qué es has_credits: bloqueo S6-tablero-lineas-archivadas.
		pays, hasCredits, err := creditos.PaymentsByMonth(ctx, db, today(reg), dominio.FromEngineDate(m.First()), dominio.FromEngineDate(m.Last()))
		if err != nil {
			return err
		}
		p := pays[m.String()]
		out := api.DashboardDebtCapacity{Month: m.String(), IncomeIsExpected: inc.expected, ExtrasCents: inc.extras,
			PaymentsCents: p.HouseholdCents, PaymentsTotalCents: p.TotalCents, PaymentCount: p.Count,
			WarningPct: warning, LimitPct: limit, ScaleMaxPct: scaleMaxPct(limit), HasCredits: hasCredits}
		// Sin quincena capturada ni esperada no hay ingreso (sin datos no es cero): sin porcentaje, estado ni margen.
		if inc.fortnights > 0 {
			ratio := dominio.PercentOf(p.HouseholdCents, inc.fortnights)
			out.IncomeCents = api.Ptr(inc.fortnights)
			out.RatioPct = api.Ptr(ratio)
			out.RatioTotalPct = api.Ptr(dominio.PercentOf(p.TotalCents, inc.fortnights))
			out.Status = api.Ptr(debtStatus(ratio, warning, limit))
			out.HeadroomCents = api.Ptr(dominio.MulDivRound(inc.fortnights, limit, 100) - p.HouseholdCents)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// monthKey lee una llave AAAA-MM obligatoria (formato: 400; C1 §3).
func monthKey(r *http.Request, name string) (mf.Month, error) {
	s := r.URL.Query().Get(name)
	if s == "" {
		return mf.Month{}, api.BadRequest(fmt.Sprintf("Falta «%s»: un mes con el formato AAAA-MM, por ejemplo 2026-10.", name))
	}
	m, err := mf.ParseMonth(s)
	if err != nil {
		return mf.Month{}, api.BadRequest(fmt.Sprintf("«%s» no es un mes válido. Usa el formato AAAA-MM, por ejemplo 2026-10.", name))
	}
	return m, nil
}

func getDashboardMonthlyOverview(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		from, err := monthKey(r, "from")
		if err != nil {
			return err
		}
		to, err := monthKey(r, "to")
		if err != nil {
			return err
		}
		n := monthIndex(to) - monthIndex(from) + 1
		switch {
		case n < 1:
			return dominio.Invalid(dominio.FieldIssue{Field: "from", Code: "out_of_range",
				Message: "El primer mes tiene que ser igual o anterior al último. Revisa «from» y «to»."})
		case n > maxOverviewMonths:
			return dominio.Invalid(dominio.FieldIssue{Field: "to", Code: "out_of_range",
				Message: fmt.Sprintf("Puedes pedir como máximo %d meses a la vez. Acorta la ventana.", maxOverviewMonths)})
		}
		ctx, db := r.Context(), reg.Deps().DB
		warning, limit, err := thresholds(ctx, db)
		if err != nil {
			return err
		}
		_, in, err := engine(reg, ctx)
		if err != nil {
			return err
		}
		pays, _, err := creditos.PaymentsByMonth(ctx, db, today(reg), dominio.FromEngineDate(from.First()), dominio.FromEngineDate(to.Last()))
		if err != nil {
			return err
		}
		current := monthIndex(dominio.EngineMonth(today(reg)))
		out := api.DashboardMonthlyOverview{WarningPct: warning, LimitPct: limit, ScaleMaxPct: scaleMaxPct(limit), Months: make([]api.DashboardMonthlyOverviewMonthsItem, 0, n)}
		for i := 0; i < n; i++ {
			m := from.Add(i)
			p := pays[m.String()]
			it := api.DashboardMonthlyOverviewMonthsItem{Month: m.String(), IsProjection: monthIndex(m) > current,
				CreditPaymentsCents: p.HouseholdCents, CreditPaymentsTotalCents: p.TotalCents}
			inc, err := incomeOf(in, m)
			if err != nil {
				return err
			}
			it.IncomeIsExpected = inc.expected
			if inc.fortnights+inc.extras > 0 {
				it.IncomeTotalCents, it.IncomeFortnightsCents, it.IncomeExtrasCents = api.Ptr(inc.fortnights+inc.extras), api.Ptr(inc.fortnights), api.Ptr(inc.extras)
			}
			if inc.fortnights > 0 {
				ratio := dominio.PercentOf(p.HouseholdCents, inc.fortnights)
				it.DebtRatioPct, it.DebtStatus = api.Ptr(ratio), api.Ptr(debtStatus(ratio, warning, limit))
			}
			pva, err := mf.PlanVsActual(in, m, 1, "")
			if err != nil {
				return err
			}
			if pl := pva[0].PlanCents; pl != nil && *pl > 0 {
				it.PlannedCents = api.Ptr(int64(*pl))
			}
			if !it.IsProjection {
				bm, err := mf.BudgetMonth(in, m)
				if err != nil {
					return err
				}
				// Devengado = real de M-01 y efectivo = cash_household_cents de M-11 (motor S2). Mes sin datos: ausentes
				// (dominio.MonthHasData, la misma regla que el mes anterior de getDashboardMonthVsPlan; PROVISIONAL,
				// bloqueo S6-tablero-sin-datos).
				accrued, cash := int64(pva[0].ActualCents), int64(bm.CashHouseholdCents)
				if dominio.MonthHasData(accrued, cash) {
					it.AccruedCents, it.CashCents = api.Ptr(accrued), api.Ptr(cash)
				}
			}
			out.Months = append(out.Months, it)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func getBudgetSignals(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		m, err := monthParam(reg, r)
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		months, err := dominio.SettingInt(ctx, db, keySignalN)
		if err != nil {
			return err
		}
		margin, err := dominio.SettingInt(ctx, db, keySignalPct)
		if err != nil {
			return err
		}
		fe, in, err := engine(reg, ctx)
		if err != nil {
			return err
		}
		type cell struct {
			month string
			line  mf.BudgetLine
		}
		order := []string{}
		byCat := map[string][]cell{}
		current := monthIndex(dominio.EngineMonth(today(reg)))
		for i := int(months); i >= 1; i-- { // del más antiguo al más reciente; el mes pedido no cuenta
			mm := m.Add(-i)
			bm, err := mf.BudgetMonth(in, mm)
			if err != nil {
				return err
			}
			// Solo cuentan los meses cerrados (RN-31: «el mes en curso no cuenta»). Con un period posterior al mes en
			// curso, el mes en curso y los de proyección no se miran: quedan fuera de months y la racha no se completa,
			// así que la señal es none. PROVISIONAL, bloqueo S6-tablero-senal-periodo, opción (a).
			closed := monthIndex(mm) < current
			for _, l := range bm.Categories {
				if _, ok := fe.CategoryID[l.Category]; !ok {
					continue
				}
				if _, seen := byCat[l.Category]; !seen {
					order = append(order, l.Category)
					byCat[l.Category] = nil // la fila existe aunque no haya meses cerrados en la ventana
				}
				if closed {
					byCat[l.Category] = append(byCat[l.Category], cell{mm.String(), l})
				}
			}
		}
		out := api.BudgetSignals{Month: m.String(), MonthsConsidered: months, MarginPct: margin, Rows: []api.BudgetSignalsRowsItem{}}
		for _, code := range order {
			cells := byCat[code]
			row := api.BudgetSignalsRowsItem{CategoryID: fe.CategoryID[code], Months: make([]api.BudgetSignalsRowsItemMonthsItem, 0, len(cells))}
			// Racha completa solo con months_considered meses cerrados; con menos (period futuro) no hay sugerencia.
			raise, cut := len(cells) == int(months), len(cells) == int(months)
			for _, c := range cells {
				it := api.BudgetSignalsRowsItemMonthsItem{Month: c.month, ActualCents: int64(c.line.ActualCents)}
				if c.line.PlanCents != nil {
					it.PlannedCents = api.Ptr(int64(*c.line.PlanCents))
				}
				// % usado del motor (real ÷ plan × 100, ROUND_HALF_UP). Sin plan corta la racha (RN-31).
				if c.line.UsedPct == nil {
					raise, cut = false, false
				} else {
					u := *c.line.UsedPct
					it.UsedPct = api.Ptr(u)
					raise = raise && u > 100+margin
					cut = cut && u < 100-margin
				}
				row.Months = append(row.Months, it)
			}
			switch {
			case raise:
				row.Signal = "raise"
			case cut:
				row.Signal = "cut"
			default:
				row.Signal = "none"
			}
			out.Rows = append(out.Rows, row)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// ---------- Parámetros del tablero (Ajustes, K44) ----------

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

func readSettings(ctx context.Context, q dominio.Querier) (api.DashboardSettings, error) {
	var out api.DashboardSettings
	ps := &out.Parameters
	for key, dst := range map[string]*api.SettingParameterInteger{keyWarning: &ps.DebtCapacityWarningPct, keyLimit: &ps.DebtCapacityLimitPct,
		keyReminder: &ps.DashboardPaymentReminderCount, keySignalN: &ps.BudgetSignalMonths, keySignalPct: &ps.BudgetSignalMarginPct} {
		var err error
		if *dst, err = intParam(ctx, q, key); err != nil {
			return out, err
		}
	}
	return out, nil
}

func getDashboardSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		out, err := readSettings(r.Context(), reg.Deps().DB)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// updateDashboardSettings: parameters (valor nuevo por clave) y reset_parameters (vuelven al inicial) con
// dominio.ApplyParams (tipo, clave de la sección, «cambiar y restablecer a la vez» y rango de los disparadores de la
// 00005: 422). Después, con los valores que resultarían, el umbral amarillo tiene que quedar menor que el tope
// (RN-30, @param_rango de la 00005): si no, 422 y no se guarda nada.
func updateDashboardSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		limitBytes := reg.Deps().Config.MaxJSONBytes
		b, err := io.ReadAll(io.LimitReader(r.Body, limitBytes+1))
		if err != nil {
			return err
		}
		if int64(len(b)) > limitBytes {
			return &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El cuerpo de la petición es demasiado grande."}
		}
		r.Body = io.NopCloser(bytes.NewReader(b))
		var in api.DashboardSettingsUpdate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var raw struct {
			Parameters      map[string]json.RawMessage `json:"parameters"`
			ResetParameters *[]string                  `json:"reset_parameters"`
		}
		if err := json.Unmarshal(b, &raw); err != nil {
			return api.BadRequest("El cuerpo debe ser un objeto JSON.")
		}
		if len(raw.Parameters) == 0 && len(in.ResetParameters) == 0 {
			return api.BadRequest("No hay cambios que guardar: manda al menos un parámetro en parameters o en reset_parameters.")
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
			if err := dominio.ApplyParams(ctx, tx, dominio.ParamChange{Set: raw.Parameters, Reset: in.ResetParameters}, dashboardKeys, api.Actor(r), api.RequestID(r)); err != nil {
				return err
			}
			warning, limit, err := thresholds(ctx, tx)
			if err != nil {
				return err
			}
			if warning >= limit {
				return dominio.Invalid(dominio.FieldIssue{Field: "parameters." + keyWarning, Code: "out_of_range",
					Message: fmt.Sprintf("El umbral amarillo (%d %%) tiene que ser menor que el tope (%d %%). Baja el amarillo o sube el tope.", warning, limit)})
			}
			return nil
		})
		if err != nil {
			return err
		}
		return getDashboardSettings(reg)(w, r)
	}
}
