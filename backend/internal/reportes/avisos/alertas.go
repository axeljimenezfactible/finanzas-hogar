package avisos

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/telegram"
)

// TickAlertas evalúa las cuatro alertas de la decisión F a partir de alerts.send_time.
// Cada alerta sale una vez por disparador y periodo (dedupe_key); un reintento no duplica.
func (m *Motor) TickAlertas(ctx context.Context, now time.Time) error {
	a, err := LeerAjustes(ctx, m.DB)
	if err != nil || hhmm(now) < a.AlertTime {
		return err
	}
	hoy := dominio.DateOf(now)
	var errs []error
	pend := func(name string, err error) {
		if errors.Is(err, ErrSinFuente) {
			m.log().Info("avisos: alerta sin fuente de cifras (PROVISIONAL)", "alerta", name)
			return
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}
	pend("payment_due", m.alertaPagos(ctx, a, hoy, now))
	pend("budget_threshold", m.alertaPresupuesto(ctx, a, hoy, now))
	pend("stale_prices", m.alertaPrecios(ctx, a, hoy, now))
	pend("meal_plan_missing", m.alertaPlan(ctx, a, hoy, now))
	return errors.Join(errs...)
}

func (m *Motor) una(ctx context.Context, code, key string, now time.Time, texto func() string, soloAnotar bool) error {
	if done, err := m.yaEntregado(ctx, code, key); err != nil || done {
		return err
	}
	_, err := m.Entregar(ctx, code, key, texto(), soloAnotar, now)
	return err
}

// (1) Pago próximo: alerts.payment_due_days_before ([1, 0]: el día anterior y el mismo día).
func (m *Motor) alertaPagos(ctx context.Context, a Ajustes, hoy dominio.Date, now time.Time) error {
	for _, d := range a.DaysBefore {
		ps, err := m.Fuente.Pagos(ctx, hoy.AddDays(int(d)), hoy.AddDays(int(d)))
		if err != nil {
			return err
		}
		for _, p := range ps {
			key := "payment_due:" + strconv.FormatInt(p.ID, 10) + ":" + hoy.String()
			if err := m.una(ctx, "payment_due", key, now, func() string { return TextoPago(p, int(d)) }, false); err != nil {
				return err
			}
		}
	}
	return nil
}

// TextoPago: «PAGO PRÓXIMO» con cuándo, cuánto y de qué crédito.
func TextoPago(p Pago, dias int) string {
	cuando := "Hoy"
	switch {
	case dias == 1:
		cuando = "Mañana"
	case dias > 1:
		cuando = "En " + strconv.Itoa(dias) + " días"
	}
	return telegram.Mensaje("Pago próximo", cuando+", "+fechaLarga(p.Fecha)+", vence un pago de "+pesos(p.MontoCents)+" ("+p.Credito+").")
}

// (2) Presupuesto: alerts.budget_warning_pct y alerts.budget_limit_pct por categoría y total (scopes).
// Si se rebasan los dos a la vez sale solo el mayor y el menor queda anotado (skipped) para no repetir.
func (m *Motor) alertaPresupuesto(ctx context.Context, a Ajustes, hoy dominio.Date, now time.Time) error {
	mesP := dominio.MonthPeriod(hoy)
	us, err := m.Fuente.Presupuesto(ctx, mesP)
	if err != nil {
		return err
	}
	sc, err := m.scopes(ctx)
	if err != nil {
		return err
	}
	for _, u := range us {
		if u.PlanCents <= 0 || !sc[scopeOf(u.Clave)] {
			continue
		}
		pct := u.UsadoCents * 100 / u.PlanCents
		enviado := false
		for _, th := range []int64{a.Limit, a.Warn} {
			if pct < th {
				continue
			}
			key := fmt.Sprintf("budget_threshold:%s:%s:%d", hoy.MonthKey(), u.Clave, th)
			if err := m.una(ctx, "budget_threshold", key, now, func() string { return TextoPresupuesto(u, th, hoy) }, enviado); err != nil {
				return err
			}
			enviado = true
		}
	}
	return nil
}

// TextoPresupuesto: «PRESUPUESTO AL 80 %».
func TextoPresupuesto(u Uso, th int64, hoy dominio.Date) string {
	return telegram.Mensaje(fmt.Sprintf("Presupuesto al %d %%", th),
		fmt.Sprintf("%s lleva %s de %s en %s.", u.Nombre, pesos(u.UsadoCents), pesos(u.PlanCents), mes(hoy)))
}

// (3) Precios viejos: alerts.stale_prices_advance_days antes del día de despensa.
func (m *Motor) alertaPrecios(ctx context.Context, a Ajustes, hoy dominio.Date, now time.Time) error {
	g, ok, err := m.Fuente.Despensa(ctx, hoy)
	if err != nil || !ok || hoy != g.AddDays(-a.StaleAdvance) {
		return err
	}
	n, err := m.Fuente.PreciosViejos(ctx, g)
	if err != nil || n == 0 {
		return err
	}
	return m.una(ctx, "stale_prices", "stale_prices:"+g.String(), now, func() string { return TextoPrecios(g, n) }, false)
}

// TextoPrecios: «PRECIOS VIEJOS».
func TextoPrecios(g dominio.Date, n int) string {
	return telegram.Mensaje("Precios viejos",
		fmt.Sprintf("Para la despensa del %s hay %d precio(s) viejo(s). Actualízalos en Compras antes de ir.", fechaLarga(g), n))
}

// (4) Plan de comidas sin armar: el día alerts.meal_plan_missing_weekday, sobre la semana que empieza el
// lunes siguiente (PROVISIONAL: qué semana revisa, bloqueo S7-fuentes-avisos punto 3).
func (m *Motor) alertaPlan(ctx context.Context, a Ajustes, hoy dominio.Date, now time.Time) error {
	if isoWeekday(hoy) != a.MealDay {
		return nil
	}
	lunes := hoy.AddDays(8 - isoWeekday(hoy))
	n, total, err := m.Fuente.Casillas(ctx, lunes)
	if err != nil || n == 0 {
		return err
	}
	return m.una(ctx, "meal_plan_missing", "meal_plan_missing:"+lunes.String(), now, func() string { return TextoPlan(lunes, n, total) }, false)
}

// TextoPlan: «PLAN DE COMIDAS SIN ARMAR».
func TextoPlan(lunes dominio.Date, n, total int) string {
	return telegram.Mensaje("Plan de comidas sin armar",
		fmt.Sprintf("La semana del %s al %s tiene %d de %d casillas sin armar.", fechaCorta(lunes), fechaCorta(lunes.AddDays(6)), n, total))
}

// Probar manda ahora el aviso code a una persona (sendTestNotification, disparo manual de S7-07).
// Devuelve el id de notification_delivery; si sale bien pone verified_at al destinatario.
func (m *Motor) Probar(ctx context.Context, code string, personID int64, now time.Time) (int64, error) {
	var nid, rid int64
	var chat string
	if err := m.DB.QueryRowContext(ctx, `SELECT id FROM notification WHERE code = ?`, code).Scan(&nid); err != nil {
		return 0, err
	}
	err := m.DB.QueryRowContext(ctx, `SELECT id, chat_id FROM notification_recipient WHERE person_id = ? AND is_active = 1 ORDER BY id LIMIT 1`, personID).Scan(&rid, &chat)
	if err != nil {
		return 0, dominio.Conflict("conflict", "Esta persona todavía no tiene su chat de Telegram confirmado. Confírmalo en Ajustes, Notificaciones.")
	}
	text, err := m.TextoPrueba(ctx, code, dominio.DateOf(now))
	if err != nil {
		return 0, err
	}
	id, _, err := m.anotar(ctx, nid, rid, "test:"+code+":"+strconv.FormatInt(now.UnixNano(), 10), text)
	if err != nil {
		return 0, err
	}
	if m.enviar(ctx, id, chat, text, now) {
		_, err = m.DB.ExecContext(ctx, `UPDATE notification_recipient SET verified_at = ? WHERE id = ?`, dominio.FormatDBTime(now), rid)
	}
	return id, err
}

// TextoPrueba arma el texto real del aviso con las cifras de hoy, más una línea que dice que es una prueba.
func (m *Motor) TextoPrueba(ctx context.Context, code string, hoy dominio.Date) (string, error) {
	const pie = "\nMensaje de prueba enviado desde Ajustes."
	switch code {
	case "daily_summary":
		t, err := m.TextoDiario(ctx, hoy)
		return t + pie, err
	case "weekly_summary":
		t, err := m.TextoSemanal(ctx, hoy)
		return t + pie, err
	case "payment_due":
		ps, err := m.Fuente.Pagos(ctx, hoy, hoy.AddDays(31))
		if err != nil {
			return "", err
		}
		if len(ps) == 0 {
			return telegram.Mensaje("Pago próximo", "No hay pagos programados en los próximos 31 días.") + pie, nil
		}
		return TextoPago(ps[0], hoy.DaysUntil(ps[0].Fecha)) + pie, nil
	case "meal_plan_missing":
		lunes := hoy.AddDays(8 - isoWeekday(hoy))
		n, total, err := m.Fuente.Casillas(ctx, lunes)
		return TextoPlan(lunes, n, total) + pie, err
	}
	var name string
	if err := m.DB.QueryRowContext(ctx, `SELECT name FROM notification WHERE code = ?`, code).Scan(&name); err != nil {
		return "", err
	}
	return telegram.Mensaje("Aviso de prueba", "Este chat recibirá el aviso «"+name+"».") + pie, nil
}
