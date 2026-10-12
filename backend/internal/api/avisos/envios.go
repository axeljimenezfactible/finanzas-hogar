package avisos

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/telegram"
)

func (a *area) sendTestNotification(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "notification_id")
	if err != nil {
		return err
	}
	var in api.TestNotificationRequest
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	n, err := loadNotif(r.Context(), a.db(), id)
	if err != nil {
		return err
	}
	var v api.Validator
	v.Required(in.PersonID > 0, "person_id")
	if err := v.Err(); err != nil {
		return err
	}
	if err := telegram.PersonaDelHogar(r.Context(), a.db(), in.PersonID, "person_id"); err != nil {
		return err
	}
	did, err := a.m.Probar(r.Context(), n.Code, in.PersonID, a.reg.Deps().Clock.Now())
	if err != nil {
		return err
	}
	d, err := loadDelivery(r, a.db(), did)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusAccepted, d)
}

const delivCols = `d.id, d.notification_id, d.recipient_id, d.dedupe_key, d.status, d.message_text, d.http_status, d.error, d.sent_at, d.created_at, d.updated_at, n.code`

func scanDelivery(sc interface{ Scan(...any) error }) (api.NotificationDelivery, error) {
	var x api.NotificationDelivery
	var text, errTxt, sent sql.NullString
	var hs sql.NullInt64
	var ca, ua, code string
	if err := sc.Scan(&x.ID, &x.NotificationID, &x.RecipientID, &x.DedupeKey, &x.Status, &text, &hs, &errTxt, &sent, &ca, &ua, &code); err != nil {
		return x, err
	}
	if text.Valid {
		x.MessageText = api.Ptr(text.String)
	}
	if hs.Valid {
		x.HTTPStatus = api.Ptr(hs.Int64)
	}
	if errTxt.Valid {
		x.Error = api.Ptr(errTxt.String)
	}
	if sent.Valid {
		x.SentAt = api.Ptr(api.APITime(sent.String))
	}
	x.CreatedAt, x.UpdatedAt, x.NotificationCode = api.APITime(ca), api.APITime(ua), api.Ptr(code)
	return x, nil
}

func loadDelivery(r *http.Request, q dominio.Querier, id int64) (api.NotificationDelivery, error) {
	return scanDelivery(q.QueryRowContext(r.Context(), `SELECT `+delivCols+` FROM notification_delivery d JOIN notification n ON n.id = d.notification_id WHERE d.id = ?`, id))
}

func (a *area) listNotificationDeliveries(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"created_at", "sent_at", "id"}, "-created_at")
	if err != nil {
		return err
	}
	per, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	ids, err := api.QueryIDs(r, "notification_id")
	if err != nil {
		return err
	}
	where, args := []string{"1 = 1"}, []any{}
	if per != nil {
		where = append(where, "d.created_at >= ? AND d.created_at < ?")
		args = append(args, dominio.FormatDBTime(per.Start.StartOfDay()), dominio.FormatDBTime(per.End.AddDays(1).StartOfDay()))
	}
	if len(ids) > 0 {
		where = append(where, "d.notification_id IN (?"+strings.Repeat(", ?", len(ids)-1)+")")
		for _, id := range ids {
			args = append(args, id)
		}
	}
	rows, err := a.db().QueryContext(r.Context(), `SELECT `+delivCols+` FROM notification_delivery d JOIN notification n ON n.id = d.notification_id WHERE `+
		strings.Join(where, " AND ")+" "+p.OrderBy(map[string]string{"created_at": "d.created_at", "sent_at": "d.sent_at", "id": "d.id"})+" "+p.LimitOffset(), args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []api.NotificationDelivery{}
	for rows.Next() {
		x, err := scanDelivery(rows)
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

func (a *area) listTelegramChatCandidates(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"detected_at", "id"}, "-detected_at")
	if err != nil {
		return err
	}
	rows, err := a.db().QueryContext(r.Context(), `SELECT id, chat_id, telegram_name, detected_at, created_at, updated_at FROM telegram_chat_candidate
 WHERE confirmed_at IS NULL `+p.OrderBy(map[string]string{"detected_at": "detected_at", "id": "id"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	defer rows.Close()
	items := []api.TelegramChatCandidate{}
	for rows.Next() {
		var x api.TelegramChatCandidate
		var chat, det, ca, ua string
		if err := rows.Scan(&x.ID, &chat, &x.TelegramName, &det, &ca, &ua); err != nil {
			return err
		}
		x.DetectedAt, x.CreatedAt, x.UpdatedAt = api.APITime(det), api.APITime(ca), api.APITime(ua)
		x.ChatIDMasked = api.Ptr(telegram.MaskChatID(chat)) // el chat_id completo nunca sale (K4)
		items = append(items, x)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

func (a *area) confirmTelegramChatCandidate(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "candidate_id")
	if err != nil {
		return err
	}
	var in api.TelegramChatCandidateConfirm
	if err := a.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	pid := api.Session(r).PersonID
	if in.PersonID != nil {
		pid = *in.PersonID
	}
	var rid int64
	err = dominio.WithTx(r.Context(), a.db(), func(tx *sql.Tx) (err error) {
		rid, err = telegram.ConfirmarCandidato(r.Context(), tx, id, pid, a.reg.Deps().Clock.Now(), api.Actor(r), api.RequestID(r))
		return err
	})
	if err != nil {
		return err
	}
	x, err := loadRecip(r.Context(), a.db(), rid)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, x)
}

func (a *area) dismissTelegramChatCandidate(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "candidate_id")
	if err != nil {
		return err
	}
	err = dominio.WithTx(r.Context(), a.db(), func(tx *sql.Tx) error {
		var n int
		err := tx.QueryRowContext(r.Context(), `SELECT count(*) FROM telegram_chat_candidate WHERE id = ? AND confirmed_at IS NULL`, id).Scan(&n)
		if err != nil {
			return err
		}
		if n == 0 {
			return dominio.NotFound("Ese chat detectado no existe o ya está confirmado.")
		}
		before, _ := dominio.RowJSON(r.Context(), tx, "telegram_chat_candidate", id)
		if _, err := tx.ExecContext(r.Context(), `DELETE FROM telegram_chat_candidate WHERE id = ?`, id); err != nil {
			return err
		}
		return dominio.Audit(r.Context(), tx, api.Actor(r), "delete", "telegram_chat_candidate", id, before, nil, api.RequestID(r))
	})
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.NotFound("Ese chat detectado no existe.")
	}
	if err != nil {
		return err
	}
	return api.NoContent(w)
}
