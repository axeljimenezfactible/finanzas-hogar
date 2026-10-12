package telegram

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/dominio"
)

// Destino es un destinatario activo suscrito a un aviso.
type Destino struct {
	RecipientID, PersonID int64
	Nombre, ChatID        string
}

// Omitido es una persona del hogar que no recibe el aviso y por qué (S7-09: se reporta, no falla el resto).
type Omitido struct {
	PersonID      int64
	Nombre, Razon string
}

// Destinatarios da a quién va un aviso: suscripciones activas de personas del hogar (titular y co-titular;
// una persona externa como Karen nunca, decisión A) y, aparte, las personas del hogar sin chat registrado.
func Destinatarios(ctx context.Context, q dominio.Querier, notificationID int64) ([]Destino, []Omitido, error) {
	rows, err := q.QueryContext(ctx, `SELECT r.id, r.person_id, p.display_name, r.chat_id
 FROM notification_subscription s JOIN notification_recipient r ON r.id = s.recipient_id JOIN person p ON p.id = r.person_id
 WHERE s.notification_id = ? AND r.is_active = 1 AND p.is_active = 1 AND p.participation IN ('holder', 'co_holder')
 ORDER BY p.id, r.id`, notificationID)
	if err != nil {
		return nil, nil, err
	}
	var ds []Destino
	for rows.Next() {
		var d Destino
		if err := rows.Scan(&d.RecipientID, &d.PersonID, &d.Nombre, &d.ChatID); err != nil {
			rows.Close()
			return nil, nil, err
		}
		ds = append(ds, d)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	rows, err = q.QueryContext(ctx, `SELECT p.id, p.display_name FROM person p
 WHERE p.is_active = 1 AND p.participation IN ('holder', 'co_holder')
   AND NOT EXISTS (SELECT 1 FROM notification_recipient r WHERE r.person_id = p.id AND r.is_active = 1)
 ORDER BY p.id`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var om []Omitido
	for rows.Next() {
		o := Omitido{Razon: "sin chat de Telegram registrado"}
		if err := rows.Scan(&o.PersonID, &o.Nombre); err != nil {
			return nil, nil, err
		}
		om = append(om, o)
	}
	return ds, om, rows.Err()
}

// PersonaDelHogar dice si la persona existe y puede recibir avisos (titular o co-titular activa).
// Karen (external) da un error 422 not_allowed en field.
func PersonaDelHogar(ctx context.Context, q dominio.Querier, personID int64, field string) error {
	var part string
	var active bool
	err := q.QueryRowContext(ctx, `SELECT participation, is_active FROM person WHERE id = ?`, personID).Scan(&part, &active)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Invalid(dominio.FieldIssue{Field: field, Code: "not_found", Message: "Esa persona no existe."})
	}
	if err != nil {
		return err
	}
	if part == "external" || !active {
		return dominio.Invalid(dominio.FieldIssue{Field: field, Code: "not_allowed",
			Message: "Solo Axel y Ximena pueden recibir avisos; una persona externa no puede ser destinataria."})
	}
	return nil
}

// ValidarDestinatarios revisa una lista completa de ids de notification_recipient (setNotificationRecipients).
func ValidarDestinatarios(ctx context.Context, q dominio.Querier, ids []int64) error {
	seen := map[int64]bool{}
	for i, id := range ids {
		f := "ids[" + strconv.Itoa(i) + "]"
		if seen[id] {
			return dominio.Invalid(dominio.FieldIssue{Field: f, Code: "conflict", Message: "El destinatario está repetido."})
		}
		seen[id] = true
		var pid int64
		err := q.QueryRowContext(ctx, `SELECT person_id FROM notification_recipient WHERE id = ?`, id).Scan(&pid)
		if errors.Is(err, sql.ErrNoRows) {
			return dominio.Invalid(dominio.FieldIssue{Field: f, Code: "not_found", Message: "Ese destinatario no existe."})
		}
		if err != nil {
			return err
		}
		if err := PersonaDelHogar(ctx, q, pid, f); err != nil {
			return err
		}
	}
	return nil
}

// MaskChatID muestra «•••» y los últimos 4 caracteres (K4, C1-18).
func MaskChatID(chatID string) string {
	r := []rune(chatID)
	if len(r) > 4 {
		r = r[len(r)-4:]
	}
	return "•••" + string(r)
}

// Update es la parte que se usa de un objeto Update de la API de Telegram (core.telegram.org/bots/api#update).
type Update struct {
	UpdateID int64 `json:"update_id"`
	Message  *struct {
		Chat struct {
			ID        int64  `json:"id"`
			Title     string `json:"title"`
			Username  string `json:"username"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"chat"`
	} `json:"message"`
}

// RegistrarActualizacion anota como candidato pendiente el chat de un mensaje entrante (K4): guarda solo
// chat_id y el nombre que da Telegram. Un chat ya registrado o ya detectado no cambia (nuevo = false).
// Quién entrega las actualizaciones sin chocar con MultiX lo decide el bloqueo S7-convivencia-multix.
func RegistrarActualizacion(ctx context.Context, q dominio.Querier, raw []byte, now time.Time) (id int64, nuevo bool, err error) {
	var u Update
	if err := json.Unmarshal(raw, &u); err != nil || u.Message == nil || u.Message.Chat.ID == 0 {
		return 0, false, errors.New("telegram: la actualización no trae un mensaje con chat")
	}
	c := u.Message.Chat
	chatID := strconv.FormatInt(c.ID, 10)
	name := strings.TrimSpace(c.FirstName + " " + c.LastName)
	if name == "" {
		name = strings.TrimSpace(c.Title)
	}
	if name == "" && c.Username != "" {
		name = "@" + c.Username
	}
	if name == "" {
		name = "Sin nombre"
	}
	var n int
	if err := q.QueryRowContext(ctx, `SELECT count(*) FROM notification_recipient WHERE channel = 'telegram' AND chat_id = ?`, chatID).Scan(&n); err != nil {
		return 0, false, err
	}
	if n > 0 {
		return 0, false, nil
	}
	res, err := q.ExecContext(ctx, `INSERT INTO telegram_chat_candidate (chat_id, telegram_name, detected_at) VALUES (?, ?, ?)
 ON CONFLICT (chat_id) DO NOTHING`, chatID, name, dominio.FormatDBTime(now))
	if err != nil {
		return 0, false, err
	}
	if k, _ := res.RowsAffected(); k == 0 {
		err := q.QueryRowContext(ctx, `SELECT id FROM telegram_chat_candidate WHERE chat_id = ?`, chatID).Scan(&id)
		return id, false, err
	}
	id, err = res.LastInsertId()
	return id, true, err
}

// ConfirmarCandidato liga un candidato pendiente a una persona del hogar y crea su notification_recipient
// en la misma transacción (K4). Ya confirmado → 409 conflict; inexistente → 404.
func ConfirmarCandidato(ctx context.Context, tx *sql.Tx, candidateID, personID int64, now time.Time, actor *int64, requestID string) (int64, error) {
	var chatID, name string
	var confirmed sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT chat_id, telegram_name, confirmed_at FROM telegram_chat_candidate WHERE id = ?`, candidateID).Scan(&chatID, &name, &confirmed)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, dominio.NotFound("Ese chat detectado no existe.")
	}
	if err != nil {
		return 0, err
	}
	if confirmed.Valid {
		return 0, dominio.Conflict("conflict", "Ese chat ya está confirmado.")
	}
	if err := PersonaDelHogar(ctx, tx, personID, "person_id"); err != nil {
		return 0, err
	}
	before, _ := dominio.RowJSON(ctx, tx, "telegram_chat_candidate", candidateID)
	ts := dominio.FormatDBTime(now)
	if _, err := tx.ExecContext(ctx, `UPDATE telegram_chat_candidate SET person_id = ?, confirmed_at = ? WHERE id = ?`, personID, ts, candidateID); err != nil {
		return 0, err
	}
	after, _ := dominio.RowJSON(ctx, tx, "telegram_chat_candidate", candidateID)
	if err := dominio.Audit(ctx, tx, actor, "update", "telegram_chat_candidate", candidateID, before, after, requestID); err != nil {
		return 0, err
	}
	res, err := tx.ExecContext(ctx, `INSERT INTO notification_recipient (person_id, channel, chat_id, label) VALUES (?, 'telegram', ?, ?)`, personID, chatID, name)
	if dominio.IsUnique(err) {
		return 0, dominio.Conflict("conflict", "Ese chat ya está registrado como destinatario.")
	}
	if err != nil {
		return 0, fmt.Errorf("telegram: alta de destinatario: %w", err)
	}
	rid, _ := res.LastInsertId()
	row, _ := dominio.RowJSON(ctx, tx, "notification_recipient", rid)
	return rid, dominio.Audit(ctx, tx, actor, "insert", "notification_recipient", rid, nil, row, requestID)
}
