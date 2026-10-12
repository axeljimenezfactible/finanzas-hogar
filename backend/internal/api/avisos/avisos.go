// Package avisos expone las 10 operaciones de avisos de C1 (Ajustes «Notificaciones» y avisos enviados)
// y registra los jobs de resúmenes y alertas en el planificador de base.
package avisos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	motor "finanzas-hogar/internal/reportes/avisos"
	"finanzas-hogar/internal/telegram"
)

// Options cambia el envío y las cifras (pruebas: servidor falso). Vacío = producción.
type Options struct {
	Enviador motor.Enviador
	Fuente   motor.Fuente
}

// Register es la entrada de areas.go: cliente real de Telegram con el token de ~/.openclaw/openclaw.json (K30).
func Register(reg *api.Registry) { RegisterWith(Options{})(reg) }

// RegisterWith arma el área con opciones.
func RegisterWith(o Options) func(*api.Registry) {
	return func(reg *api.Registry) {
		d := reg.Deps()
		if o.Enviador == nil {
			o.Enviador = telegram.New(telegram.FileToken(telegram.DefaultConfigPath()))
		}
		if o.Fuente == nil {
			o.Fuente = motor.FuenteBD{DB: d.DB}
		}
		h := &area{reg: reg, m: &motor.Motor{DB: d.DB, Fuente: o.Fuente, Enviador: o.Enviador, Log: d.Log}}
		if d.Jobs != nil {
			for _, j := range h.m.Jobs() {
				if err := d.Jobs.Register(j); err != nil {
					panic(err)
				}
			}
		}
		reg.Handle("listNotifications", h.listNotifications)
		reg.Handle("updateNotification", h.updateNotification)
		reg.Handle("setNotificationRecipients", h.setNotificationRecipients)
		reg.Handle("listNotificationRecipients", h.listNotificationRecipients)
		reg.Handle("updateNotificationRecipient", h.updateNotificationRecipient)
		reg.Handle("sendTestNotification", h.sendTestNotification)
		reg.Handle("listTelegramChatCandidates", h.listTelegramChatCandidates)
		reg.Handle("confirmTelegramChatCandidate", h.confirmTelegramChatCandidate)
		reg.Handle("dismissTelegramChatCandidate", h.dismissTelegramChatCandidate)
		reg.Handle("listNotificationDeliveries", h.listNotificationDeliveries)
	}
}

type area struct {
	reg *api.Registry
	m   *motor.Motor
}

func (a *area) db() *sql.DB { return a.reg.Deps().DB }

const notifCols = `n.id, n.code, n.name, n.is_enabled, n.params_json, n.created_at, n.updated_at,
 (SELECT j.id FROM job j WHERE j.notification_id = n.id OR j.code = CASE WHEN n.code IN ('payment_due', 'budget_threshold', 'stale_prices', 'meal_plan_missing') THEN 'alerts' ELSE n.code END ORDER BY j.id LIMIT 1)`

func scanNotif(ctx context.Context, q dominio.Querier, sc interface{ Scan(...any) error }) (api.Notification, error) {
	var n api.Notification
	var raw, ca, ua string
	var job sql.NullInt64
	if err := sc.Scan(&n.ID, &n.Code, &n.Name, &n.IsEnabled, &raw, &ca, &ua, &job); err != nil {
		return n, err
	}
	n.ParamsJSON = map[string]any{}
	_ = json.Unmarshal([]byte(raw), &n.ParamsJSON)
	n.CreatedAt, n.UpdatedAt = api.APITime(ca), api.APITime(ua)
	if job.Valid {
		n.JobID = api.Ptr(job.Int64)
	}
	ids, err := subscriptions(ctx, q, n.ID)
	n.RecipientIDs = ids
	return n, err
}

func subscriptions(ctx context.Context, q dominio.Querier, nid int64) ([]int64, error) {
	rows, err := q.QueryContext(ctx, `SELECT recipient_id FROM notification_subscription WHERE notification_id = ? ORDER BY recipient_id`, nid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (a *area) listNotifications(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"id", "code", "name"}, "id")
	if err != nil {
		return err
	}
	rows, err := a.db().QueryContext(r.Context(), `SELECT `+notifCols+` FROM notification n `+
		p.OrderBy(map[string]string{"id": "n.id", "code": "n.code", "name": "n.name"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	var raws []api.Notification
	for rows.Next() {
		var n api.Notification
		var raw, ca, ua string
		var job sql.NullInt64
		if err := rows.Scan(&n.ID, &n.Code, &n.Name, &n.IsEnabled, &raw, &ca, &ua, &job); err != nil {
			rows.Close()
			return err
		}
		n.ParamsJSON = map[string]any{}
		_ = json.Unmarshal([]byte(raw), &n.ParamsJSON)
		n.CreatedAt, n.UpdatedAt = api.APITime(ca), api.APITime(ua)
		if job.Valid {
			n.JobID = api.Ptr(job.Int64)
		}
		raws = append(raws, n)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	items := make([]api.Notification, 0, len(raws))
	for _, n := range raws {
		if n.RecipientIDs, err = subscriptions(r.Context(), a.db(), n.ID); err != nil {
			return err
		}
		items = append(items, n)
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

func loadNotif(ctx context.Context, q dominio.Querier, id int64) (api.Notification, error) {
	n, err := scanNotif(ctx, q, q.QueryRowContext(ctx, `SELECT `+notifCols+` FROM notification n WHERE n.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return n, dominio.NotFound("Ese aviso no existe.")
	}
	return n, err
}

// Claves que ya no viven en params_json (K47): se cambian en updateReportSettings y updateAlertSettings.
var fueraDeParams = map[string]bool{"send_time": true, "advance_days": true, "weekday": true, "day_of_month": true, "days_before": true, "thresholds_pct": true}

func validarParams(code string, p map[string]any, v *api.Validator) {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f := "params_json." + k
		switch {
		case fueraDeParams[k]:
			v.Add(f, "not_allowed", "La hora, el día, la anticipación, los umbrales y los días antes del pago se cambian en Ajustes de reportes o de alertas.")
		case code == "budget_threshold" && k == "scopes":
			arr, ok := p[k].([]any)
			seen := map[string]bool{}
			for _, x := range arr {
				s, _ := x.(string)
				ok = ok && (s == "category" || s == "total") && !seen[s]
				seen[s] = true
			}
			v.Check(ok && len(arr) > 0, f, "invalid_format", "Elige «category», «total» o los dos, sin repetir.")
		default:
			v.Add(f, "not_allowed", "Este aviso no tiene ese parámetro.")
		}
	}
}

func (a *area) updateNotification(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "notification_id")
	if err != nil {
		return err
	}
	var in api.NotificationUpdate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	cur, err := loadNotif(r.Context(), a.db(), id)
	if err != nil {
		return err
	}
	var v api.Validator
	var params any
	if in.ParamsJSON != nil {
		validarParams(cur.Code, in.ParamsJSON, &v)
		b, _ := json.Marshal(in.ParamsJSON)
		params = string(b)
	}
	if err := v.Err(); err != nil {
		return err
	}
	err = dominio.WithTx(r.Context(), a.db(), func(tx *sql.Tx) error {
		before, _ := dominio.RowJSON(r.Context(), tx, "notification", id)
		if _, err := tx.ExecContext(r.Context(), `UPDATE notification SET is_enabled = coalesce(?, is_enabled), params_json = coalesce(?, params_json) WHERE id = ?`,
			in.IsEnabled, params, id); err != nil {
			return err
		}
		after, _ := dominio.RowJSON(r.Context(), tx, "notification", id)
		return dominio.Audit(r.Context(), tx, api.Actor(r), "update", "notification", id, before, after, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	n, err := loadNotif(r.Context(), a.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, n)
}

func (a *area) setNotificationRecipients(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "notification_id")
	if err != nil {
		return err
	}
	var in api.IdList
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if _, err := loadNotif(r.Context(), a.db(), id); err != nil {
		return err
	}
	var v api.Validator
	v.Required(in.IDs != nil, "ids")
	if err := v.Err(); err != nil {
		return err
	}
	if err := telegram.ValidarDestinatarios(r.Context(), a.db(), in.IDs); err != nil {
		return err
	}
	err = dominio.WithTx(r.Context(), a.db(), func(tx *sql.Tx) error {
		prev, err := subscriptions(r.Context(), tx, id)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(r.Context(), `DELETE FROM notification_subscription WHERE notification_id = ?`, id); err != nil {
			return err
		}
		for _, rid := range in.IDs {
			if _, err := tx.ExecContext(r.Context(), `INSERT INTO notification_subscription (notification_id, recipient_id) VALUES (?, ?)`, id, rid); err != nil {
				return err
			}
		}
		before, _ := json.Marshal(map[string]any{"recipient_ids": prev})
		after, _ := json.Marshal(map[string]any{"recipient_ids": in.IDs})
		return dominio.Audit(r.Context(), tx, api.Actor(r), "update", "notification_subscription", id, before, after, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	return api.NoContent(w)
}

const recipCols = `id, person_id, channel, label, verified_at, is_active, created_at, updated_at`

func scanRecip(sc interface{ Scan(...any) error }) (api.NotificationRecipient, error) {
	var x api.NotificationRecipient
	var label, ver sql.NullString
	var ca, ua string
	if err := sc.Scan(&x.ID, &x.PersonID, &x.Channel, &label, &ver, &x.IsActive, &ca, &ua); err != nil {
		return x, err
	}
	if label.Valid {
		x.Label = api.Ptr(label.String)
	}
	if ver.Valid {
		x.VerifiedAt = api.Ptr(api.APITime(ver.String))
	}
	x.CreatedAt, x.UpdatedAt = api.APITime(ca), api.APITime(ua)
	x.IsRegistered = api.Ptr(ver.Valid && x.IsActive)
	return x, nil // nunca chat_id (C1: listNotificationRecipients)
}

func loadRecip(ctx context.Context, q dominio.Querier, id int64) (api.NotificationRecipient, error) {
	x, err := scanRecip(q.QueryRowContext(ctx, `SELECT `+recipCols+` FROM notification_recipient WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return x, dominio.NotFound("Ese destinatario no existe.")
	}
	return x, err
}

func (a *area) listNotificationRecipients(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"id", "person_id", "created_at"}, "id")
	if err != nil {
		return err
	}
	rows, err := a.db().QueryContext(r.Context(), `SELECT `+recipCols+` FROM notification_recipient `+
		p.OrderBy(map[string]string{"id": "id", "person_id": "person_id", "created_at": "created_at"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []api.NotificationRecipient{}
	for rows.Next() {
		x, err := scanRecip(rows)
		if err != nil {
			return err
		}
		items = append(items, x)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

func (a *area) updateNotificationRecipient(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "recipient_id")
	if err != nil {
		return err
	}
	var in api.NotificationRecipientUpdate
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	if _, err := loadRecip(r.Context(), a.db(), id); err != nil {
		return err
	}
	var v api.Validator
	label := strings.TrimSpace(in.Label.Value)
	v.Check(!in.Label.Valid || (label != "" && len([]rune(label)) <= 80), "label", "invalid_format", "Escribe una etiqueta de 1 a 80 caracteres o déjala vacía con null.")
	if err := v.Err(); err != nil {
		return err
	}
	err = dominio.WithTx(r.Context(), a.db(), func(tx *sql.Tx) error {
		before, _ := dominio.RowJSON(r.Context(), tx, "notification_recipient", id)
		if in.Label.Set {
			var val any
			if in.Label.Valid {
				val = label
			}
			if _, err := tx.ExecContext(r.Context(), `UPDATE notification_recipient SET label = ? WHERE id = ?`, val, id); err != nil {
				return err
			}
		}
		if in.IsActive != nil {
			if _, err := tx.ExecContext(r.Context(), `UPDATE notification_recipient SET is_active = ? WHERE id = ?`, *in.IsActive, id); err != nil {
				return err
			}
		}
		after, _ := dominio.RowJSON(r.Context(), tx, "notification_recipient", id)
		return dominio.Audit(r.Context(), tx, api.Actor(r), "update", "notification_recipient", id, before, after, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	x, err := loadRecip(r.Context(), a.db(), id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, x)
}
