package avisos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
	"finanzas-hogar/internal/telegram"
)

// Enviador manda texto plano a un chat (telegram.Client en producción; un servidor falso en pruebas).
type Enviador interface {
	SendMessage(ctx context.Context, chatID, text string) (telegram.Result, error)
}

// Motor arma y entrega los avisos.
type Motor struct {
	DB       *sql.DB
	Fuente   Fuente
	Enviador Enviador
	Log      *slog.Logger
}

// Informe resume una entrega.
type Informe struct {
	Enviados, Fallidos, Repetidos, Omitidos int
	Notas                                   []string
}

func (m *Motor) log() *slog.Logger {
	if m.Log == nil {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return m.Log
}

// Cada revisión es por minuto: la hora real sale de app_setting en cada corrida (K47), así un cambio en
// Ajustes se respeta sin tocar el job. PROPUESTA (bloqueo S7-fuentes-avisos, punto 4).
const cadaMinuto = "* * * * *"

// Jobs son los trabajos del planificador de base (idempotentes por notification_delivery.dedupe_key).
func (m *Motor) Jobs() []jobs.Def {
	return []jobs.Def{
		{Code: "daily_summary", Name: "Resumen diario por Telegram", Schedule: cadaMinuto, Run: func(ctx context.Context, r jobs.Run) error { return m.TickDiario(ctx, r.Now) }},
		{Code: "weekly_summary", Name: "Resumen semanal por Telegram", Schedule: cadaMinuto, Run: func(ctx context.Context, r jobs.Run) error { return m.TickSemanal(ctx, r.Now) }},
		{Code: "alerts", Name: "Avisos: pago próximo, presupuesto, precios viejos y plan sin armar", Schedule: cadaMinuto, Run: func(ctx context.Context, r jobs.Run) error { return m.TickAlertas(ctx, r.Now) }},
	}
}

func (m *Motor) yaEntregado(ctx context.Context, code, dedupe string) (bool, error) {
	var n int
	err := m.DB.QueryRowContext(ctx, `SELECT count(*) FROM notification_delivery d JOIN notification n ON n.id = d.notification_id
 WHERE n.code = ? AND d.dedupe_key = ?`, code, dedupe).Scan(&n)
	return n > 0, err
}

// Entregar manda text a los destinatarios suscritos del aviso code, una sola vez por dedupe y destinatario
// (UNIQUE de notification_delivery). soloAnotar deja el renglón en skipped sin mandar (umbral menor ya rebasado).
func (m *Motor) Entregar(ctx context.Context, code, dedupe, text string, soloAnotar bool, now time.Time) (Informe, error) {
	var inf Informe
	var nid int64
	var enabled bool
	if err := m.DB.QueryRowContext(ctx, `SELECT id, is_enabled FROM notification WHERE code = ?`, code).Scan(&nid, &enabled); err != nil {
		return inf, fmt.Errorf("avisos: aviso %s: %w", code, err)
	}
	if !enabled {
		inf.Notas = append(inf.Notas, "aviso apagado en Ajustes")
		return inf, nil
	}
	ds, om, err := telegram.Destinatarios(ctx, m.DB, nid)
	if err != nil {
		return inf, err
	}
	for _, o := range om {
		inf.Omitidos++
		inf.Notas = append(inf.Notas, o.Nombre+": "+o.Razon)
		m.log().Warn("avisos: destinatario omitido", "aviso", code, "person_id", o.PersonID, "motivo", o.Razon)
	}
	for _, d := range ds {
		id, ok, err := m.anotar(ctx, nid, d.RecipientID, dedupe, text)
		if err != nil {
			return inf, err
		}
		if !ok {
			inf.Repetidos++
			continue
		}
		if soloAnotar {
			_, err = m.DB.ExecContext(ctx, `UPDATE notification_delivery SET status = 'skipped' WHERE id = ?`, id)
			if err != nil {
				return inf, err
			}
			continue
		}
		if m.enviar(ctx, id, d.ChatID, text, now) {
			inf.Enviados++
		} else {
			inf.Fallidos++
		}
	}
	if inf.Fallidos > 0 {
		return inf, fmt.Errorf("avisos: %s: %d envío(s) fallido(s); ver notification_delivery", code, inf.Fallidos)
	}
	return inf, nil
}

func (m *Motor) anotar(ctx context.Context, nid, rid int64, dedupe, text string) (int64, bool, error) {
	res, err := m.DB.ExecContext(ctx, `INSERT INTO notification_delivery (notification_id, recipient_id, dedupe_key, status, message_text)
 VALUES (?, ?, ?, 'pending', ?) ON CONFLICT (dedupe_key, recipient_id) DO NOTHING`, nid, rid, dedupe, text)
	if err != nil {
		return 0, false, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return 0, false, nil
	}
	id, err := res.LastInsertId()
	return id, true, err
}

// enviar manda y anota el resultado (HTTP y error resumido, sin token).
func (m *Motor) enviar(ctx context.Context, deliveryID int64, chatID, text string, now time.Time) bool {
	var res telegram.Result
	err := errors.New("el envío por Telegram no está configurado")
	if m.Enviador != nil {
		res, err = m.Enviador.SendMessage(ctx, chatID, text)
	}
	var httpStatus any
	if res.HTTPStatus != 0 {
		httpStatus = res.HTTPStatus
	}
	if err != nil {
		msg := err.Error()
		if len(msg) > 300 {
			msg = msg[:300]
		}
		_, _ = m.DB.ExecContext(context.WithoutCancel(ctx), `UPDATE notification_delivery SET status = 'failed', http_status = ?, error = ? WHERE id = ?`, httpStatus, msg, deliveryID)
		m.log().Error("avisos: envío fallido", "delivery_id", deliveryID, "error", msg)
		return false
	}
	_, _ = m.DB.ExecContext(context.WithoutCancel(ctx), `UPDATE notification_delivery SET status = 'sent', http_status = ?, sent_at = ? WHERE id = ?`, httpStatus, dominio.FormatDBTime(now), deliveryID)
	return true
}

// TickDiario manda el resumen del día cuando la hora local llega a reports.daily_send_time.
func (m *Motor) TickDiario(ctx context.Context, now time.Time) error {
	a, err := LeerAjustes(ctx, m.DB)
	if err != nil || hhmm(now) < a.DailyTime {
		return err
	}
	hoy := dominio.DateOf(now)
	key := "daily_summary:" + hoy.String()
	if done, err := m.yaEntregado(ctx, "daily_summary", key); err != nil || done {
		return err
	}
	text, err := m.TextoDiario(ctx, hoy)
	if err != nil {
		return err
	}
	_, err = m.Entregar(ctx, "daily_summary", key, text, false, now)
	return err
}

// TickSemanal manda el resumen de la semana el día reports.weekly_weekday a reports.weekly_send_time.
func (m *Motor) TickSemanal(ctx context.Context, now time.Time) error {
	a, err := LeerAjustes(ctx, m.DB)
	hoy := dominio.DateOf(now)
	if err != nil || isoWeekday(hoy) != a.WeeklyDay || hhmm(now) < a.WeeklyTime {
		return err
	}
	key := "weekly_summary:" + hoy.ISOWeek()
	if done, err := m.yaEntregado(ctx, "weekly_summary", key); err != nil || done {
		return err
	}
	text, err := m.TextoSemanal(ctx, hoy)
	if err != nil {
		return err
	}
	_, err = m.Entregar(ctx, "weekly_summary", key, text, false, now)
	return err
}

// PROPUESTA de contenido (K31, PROVISIONAL hasta la puerta 4): ventana de pagos del resumen diario.
const diasPagosDiario = 7

// TextoDiario: gastos de hoy y del mes (las cifras de getExpenseTotals) y pagos de los próximos días.
func (m *Motor) TextoDiario(ctx context.Context, hoy dominio.Date) (string, error) {
	dia := dominio.Period{Kind: dominio.PeriodRange, Start: hoy, End: hoy}
	mesP := dominio.MonthPeriod(hoy)
	g1, err := m.Fuente.Gastos(ctx, dia)
	if err != nil {
		return "", err
	}
	g2, err := m.Fuente.Gastos(ctx, mesP)
	if err != nil {
		return "", err
	}
	ps, err := m.Fuente.Pagos(ctx, hoy.AddDays(1), hoy.AddDays(diasPagosDiario))
	if err != nil {
		return "", err
	}
	ls := []string{
		fechaLarga(hoy),
		"Gastos de hoy: " + gastos(g1.Count) + " por " + pesos(g1.AmountCents),
		"Gastos del mes (" + mes(hoy) + "): " + gastos(g2.Count) + " por " + pesos(g2.AmountCents),
	}
	ls = append(ls, lineasPagos("Pagos de los próximos "+strconv.Itoa(diasPagosDiario)+" días", ps)...)
	return telegram.Mensaje("Resumen del día", ls...), nil
}

// TextoSemanal: gastos de la semana ISO y del mes, pagos de los 7 días siguientes y casillas sin armar
// de la semana que sigue.
func (m *Motor) TextoSemanal(ctx context.Context, hoy dominio.Date) (string, error) {
	lunes := hoy.AddDays(1 - isoWeekday(hoy))
	sem := dominio.Period{Kind: dominio.PeriodRange, Start: lunes, End: lunes.AddDays(6)}
	g1, err := m.Fuente.Gastos(ctx, sem)
	if err != nil {
		return "", err
	}
	g2, err := m.Fuente.Gastos(ctx, dominio.MonthPeriod(hoy))
	if err != nil {
		return "", err
	}
	ps, err := m.Fuente.Pagos(ctx, hoy.AddDays(1), hoy.AddDays(7))
	if err != nil {
		return "", err
	}
	prox := lunes.AddDays(7)
	n, total, err := m.Fuente.Casillas(ctx, prox)
	if err != nil {
		return "", err
	}
	ls := []string{
		"Del " + fechaCorta(lunes) + " al " + fechaCorta(lunes.AddDays(6)),
		"Gastos de la semana: " + gastos(g1.Count) + " por " + pesos(g1.AmountCents),
		"Gastos del mes (" + mes(hoy) + "): " + gastos(g2.Count) + " por " + pesos(g2.AmountCents),
	}
	ls = append(ls, lineasPagos("Pagos de los próximos 7 días", ps)...)
	ls = append(ls, fmt.Sprintf("Plan de comidas de la semana del %s: %d de %d casillas sin armar", fechaCorta(prox), n, total))
	return telegram.Mensaje("Resumen de la semana", ls...), nil
}

func lineasPagos(titulo string, ps []Pago) []string {
	if len(ps) == 0 {
		return []string{titulo + ": ninguno"}
	}
	out := []string{titulo + ":"}
	for _, p := range ps {
		out = append(out, "- "+fechaLarga(p.Fecha)+": "+pesos(p.MontoCents)+" ("+p.Credito+")")
	}
	return out
}

// scopes lee budget_threshold.params_json.scopes (único parámetro propio que queda; valor de la siembra).
func (m *Motor) scopes(ctx context.Context) (map[string]bool, error) {
	var raw string
	if err := m.DB.QueryRowContext(ctx, `SELECT params_json FROM notification WHERE code = 'budget_threshold'`).Scan(&raw); err != nil {
		return nil, err
	}
	var p struct {
		Scopes []string `json:"scopes"`
	}
	_ = json.Unmarshal([]byte(raw), &p)
	out := map[string]bool{}
	for _, s := range p.Scopes {
		out[s] = true
	}
	return out, nil
}

func scopeOf(clave string) string {
	if strings.HasPrefix(clave, "category:") {
		return "category"
	}
	return "total"
}
