// Package avisos arma y manda por Telegram los resúmenes diario y semanal y las cuatro alertas de S7
// (decisión F). Horas, días, anticipación, umbrales y días antes del pago se leen de app_setting en cada
// revisión (K44, K47, K50): ninguno vive en el código ni en notification.params_json.
package avisos

import (
	"context"
	"time"

	"finanzas-hogar/internal/dominio"
)

// Ajustes son las claves reports.* y alerts.* que usan los avisos.
type Ajustes struct {
	DailyTime, WeeklyTime, AlertTime string
	WeeklyDay, MealDay, StaleAdvance int
	Warn, Limit                      int64
	DaysBefore                       []int64
}

// LeerAjustes lee las claves vigentes (valor actual o el inicial de la siembra 00003/00004).
func LeerAjustes(ctx context.Context, q dominio.Querier) (Ajustes, error) {
	var a Ajustes
	var err error
	hora := func(k string, dst *string) {
		if err == nil {
			*dst, err = dominio.SettingTime(ctx, q, k)
		}
	}
	entero := func(k string) int64 {
		if err != nil {
			return 0
		}
		var n int64
		n, err = dominio.SettingInt(ctx, q, k)
		return n
	}
	hora("reports.daily_send_time", &a.DailyTime)
	hora("reports.weekly_send_time", &a.WeeklyTime)
	hora("alerts.send_time", &a.AlertTime)
	a.WeeklyDay = int(entero("reports.weekly_weekday"))
	a.MealDay = int(entero("alerts.meal_plan_missing_weekday"))
	a.StaleAdvance = int(entero("alerts.stale_prices_advance_days"))
	a.Warn = entero("alerts.budget_warning_pct")
	a.Limit = entero("alerts.budget_limit_pct")
	if err == nil {
		a.DaysBefore, err = dominio.SettingIntList(ctx, q, "alerts.payment_due_days_before")
	}
	return a, err
}

// hhmm es la hora local (America/Mexico_City) en el formato de las claves.
func hhmm(t time.Time) string { return t.In(dominio.MX).Format("15:04") }

// isoWeekday da 1 = lunes … 7 = domingo (ISO 8601).
func isoWeekday(d dominio.Date) int {
	if w := int(d.Weekday()); w != 0 {
		return w
	}
	return 7
}
