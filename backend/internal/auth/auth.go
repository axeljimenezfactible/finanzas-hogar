// Package auth: inicio de sesión de las 2 cuentas co-titulares (Axel y Ximena; Karen no tiene cuenta, decisión A),
// sesiones de 90 días por dispositivo que se renuevan con cada uso (K11), límite de intentos fallidos en ventana
// deslizante por usuario y por IP (S6-06) y cambio de la contraseña propia (≥ auth.password_min_length, argon2id).
//
// Register también registra las operaciones de sistema (internal/api/sistema): en cmd/finanzas/areas.go basta la
// línea auth.Register (CONVENCIONES-backend.md §2).
package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/sistema"
	"finanzas-hogar/internal/dominio"
)

// HashParams son los parámetros de argon2id de las contraseñas nuevas (las pruebas pueden abaratarlos).
var HashParams = dominio.DefaultArgon2

// dummyHash se verifica cuando la cuenta no existe, para que el tiempo de respuesta no revele qué usuarios hay.
var dummyHash = sync.OnceValue(func() string {
	h, _ := dominio.HashPassword("no-es-una-cuenta-de-la-casa", HashParams)
	return h
})

// Register registra las 8 operaciones de autenticación y las 19 de sistema (27 de auth-sistema).
func Register(reg *api.Registry) {
	reg.Handle("login", login(reg))
	reg.Handle("logout", logout(reg))
	reg.Handle("changeOwnPassword", changeOwnPassword(reg))
	reg.Handle("listUsers", listUsers(reg))
	reg.Handle("listSessions", listSessions(reg))
	reg.Handle("revokeSession", revokeSession(reg))
	reg.Handle("revokeSessions", revokeSessions(reg))
	reg.Handle("getCurrentSession", getCurrentSession(reg))
	sistema.Register(reg)
}

func invalidCredentials() error {
	return &api.Problem{Status: http.StatusUnauthorized, Code: "invalid_credentials", Message: "Usuario o contraseña incorrectos."}
}

// Lockout lee los valores vigentes (auth.max_failed_attempts, auth.lockout_window_minutes; K44) y aplica la ventana
// deslizante de dominio.CheckLockout: se cuentan los fallos por usuario y por IP, lo que llegue primero al máximo.
func Lockout(ctx context.Context, q dominio.Querier, username, ip string, now time.Time) (dominio.Lockout, error) {
	maxFail, err := dominio.SettingInt(ctx, q, "auth.max_failed_attempts")
	if err != nil {
		return dominio.Lockout{}, err
	}
	win, err := dominio.SettingInt(ctx, q, "auth.lockout_window_minutes")
	if err != nil {
		return dominio.Lockout{}, err
	}
	return dominio.CheckLockout(ctx, q, username, ip, now, int(maxFail), time.Duration(win)*time.Minute)
}

func login(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.LoginRequest
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var v api.Validator
		v.Required(strings.TrimSpace(in.Username) != "", "username")
		v.Required(in.Password != "", "password")
		if err := v.Err(); err != nil {
			return err
		}
		ctx, d := r.Context(), reg.Deps()
		now, ip := d.Clock.Now(), api.ClientIP(r)
		lk, err := Lockout(ctx, d.DB, in.Username, ip, now)
		if err != nil {
			return err
		}
		if lk.Locked { // no se anota como fallo: si no, la ventana deslizante nunca se libera (S6-06 b)
			return &api.Problem{Status: http.StatusTooManyRequests, Code: "too_many_attempts",
				Message: "Demasiados intentos. Vuelve a intentar a las " + lk.RetryAt.In(dominio.MX).Format("15:04") + ".", RetryAt: lk.RetryAt}
		}
		u, err := dominio.UserByUsername(ctx, d.DB, in.Username)
		found := err == nil && u.IsActive
		if err != nil && !errors.Is(err, dominio.ErrNotFound) && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		hash := dummyHash()
		if found {
			hash = u.PasswordHash
		}
		ok, err := dominio.VerifyPassword(in.Password, hash)
		if err != nil {
			return err
		}
		if !ok || !found {
			var uid *int64
			if err == nil && u.ID != 0 {
				uid = &u.ID
			}
			if err := dominio.RecordLoginAttempt(ctx, d.DB, in.Username, uid, ip, false, now); err != nil {
				return err
			}
			return invalidCredentials()
		}
		maxAge, err := api.SessionMaxAge(ctx, d.DB)
		if err != nil {
			return err
		}
		token, tokenHash, err := dominio.NewSessionToken()
		if err != nil {
			return err
		}
		label := ""
		if in.DeviceLabel != nil {
			label = strings.TrimSpace(*in.DeviceLabel)
		}
		ua := r.UserAgent()
		if len(ua) > 300 {
			ua = ua[:300]
		}
		var sid int64
		err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			if err := dominio.RecordLoginAttempt(ctx, tx, in.Username, &u.ID, ip, true, now); err != nil {
				return err
			}
			id, err := dominio.CreateSession(ctx, tx, u.ID, tokenHash, label, ua, ip, now, maxAge)
			if err != nil {
				return err
			}
			sid = id
			return dominio.TouchLogin(ctx, tx, u.ID, now)
		})
		if err != nil {
			return err
		}
		cs, err := currentSession(ctx, d.DB, sid)
		if err != nil {
			return err
		}
		api.SetSessionCookie(w, token, maxAge)
		return api.WriteJSON(w, http.StatusOK, cs)
	}
}

func currentSession(ctx context.Context, q dominio.Querier, sessionID int64) (api.CurrentSession, error) {
	var cs api.CurrentSession
	var exp, seen string
	err := q.QueryRowContext(ctx, `SELECT u.id, u.username, u.person_id, p.display_name, p.participation, s.id, s.expires_at, s.last_seen_at
		FROM device_session s JOIN user_account u ON u.id = s.user_account_id JOIN person p ON p.id = u.person_id WHERE s.id = ?`, sessionID).
		Scan(&cs.UserAccountID, &cs.Username, &cs.PersonID, &cs.DisplayName, &cs.Participation, &cs.DeviceSessionID, &exp, &seen)
	cs.ExpiresAt, cs.LastSeenAt = api.APITime(exp), api.APITime(seen)
	return cs, err
}

func getCurrentSession(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		cs, err := currentSession(r.Context(), reg.Deps().DB, api.Session(r).SessionID)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, cs)
	}
}

// clearCookie quita la cookie renovada por el middleware y la borra (la última Set-Cookie es la que vale).
func clearCookie(w http.ResponseWriter) {
	w.Header().Del("Set-Cookie")
	api.ClearSessionCookie(w)
}

func logout(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		d := reg.Deps()
		if err := dominio.RevokeSession(r.Context(), d.DB, api.Session(r).SessionID, d.Clock.Now()); err != nil {
			return err
		}
		clearCookie(w)
		return api.NoContent(w)
	}
}

func changeOwnPassword(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.PasswordChangeRequest
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var v api.Validator
		v.Required(in.CurrentPassword != "", "current_password")
		v.Required(in.NewPassword != "", "new_password")
		if err := v.Err(); err != nil {
			return err
		}
		ctx, d, s := r.Context(), reg.Deps(), api.Session(r)
		u, err := dominio.UserByID(ctx, d.DB, s.UserID)
		if err != nil {
			return err
		}
		ok, err := dominio.VerifyPassword(in.CurrentPassword, u.PasswordHash)
		if err != nil {
			return err
		}
		if !ok {
			return &dominio.Error{Kind: dominio.KindValidation, Code: "invalid_current_password", Message: "La contraseña actual no es correcta. Revísala e intenta de nuevo."}
		}
		minLen, err := dominio.SettingInt(ctx, d.DB, "auth.password_min_length")
		if err != nil {
			return err
		}
		if int64(utf8.RuneCountInString(in.NewPassword)) < minLen {
			return dominio.Invalid(dominio.FieldIssue{Field: "new_password", Code: "out_of_range",
				Message: fmt.Sprintf("La contraseña nueva debe tener al menos %d caracteres.", minLen)})
		}
		hash, err := dominio.HashPassword(in.NewPassword, HashParams)
		if err != nil {
			return err
		}
		now := d.Clock.Now()
		return withNoContent(w, dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			var prev sql.NullString
			if err := tx.QueryRowContext(ctx, `SELECT password_changed_at FROM user_account WHERE id = ?`, u.ID).Scan(&prev); err != nil {
				return err
			}
			if err := dominio.SetPasswordHash(ctx, tx, u.ID, hash, now); err != nil {
				return err
			}
			// la bitácora nunca guarda el hash: solo cuándo cambió la contraseña
			before, _ := json.Marshal(map[string]any{"password_changed_at": nullStr(prev)})
			after, _ := json.Marshal(map[string]any{"password_changed_at": dominio.FormatDBTime(now)})
			return dominio.Audit(ctx, tx, api.Actor(r), "update", "user_account", u.ID, before, after, api.RequestID(r))
		}))
	}
}

func nullStr(n sql.NullString) any {
	if n.Valid {
		return n.String
	}
	return nil
}

func withNoContent(w http.ResponseWriter, err error) error {
	if err != nil {
		return err
	}
	return api.NoContent(w)
}

func listUsers(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"username", "created_at"}, "username")
		if err != nil {
			return err
		}
		rows, err := reg.Deps().DB.QueryContext(r.Context(), `SELECT u.id, u.person_id, u.username, u.is_active, u.last_login_at, u.created_at, u.updated_at, p.display_name, p.participation
			FROM user_account u JOIN person p ON p.id = u.person_id `+p.OrderBy(map[string]string{"username": "u.username", "created_at": "u.created_at", "id": "u.id"})+" "+p.LimitOffset())
		if err != nil {
			return err
		}
		defer rows.Close()
		items := []api.UserAccount{}
		for rows.Next() {
			var it api.UserAccount
			var last sql.NullString
			var name, part string
			if err := rows.Scan(&it.ID, &it.PersonID, &it.Username, &it.IsActive, &last, &it.CreatedAt, &it.UpdatedAt, &name, &part); err != nil {
				return err
			}
			it.CreatedAt, it.UpdatedAt = api.APITime(it.CreatedAt), api.APITime(it.UpdatedAt)
			if last.Valid {
				it.LastLoginAt = api.Ptr(api.APITime(last.String))
			}
			it.DisplayName, it.Participation = api.Ptr(name), api.Ptr(part)
			items = append(items, it) // nunca password_hash
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func listSessions(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"last_seen_at", "created_at", "expires_at"}, "-last_seen_at")
		if err != nil {
			return err
		}
		persons, err := api.QueryIDs(r, "person_id")
		if err != nil {
			return err
		}
		d := reg.Deps()
		where, args := "", []any{dominio.FormatDBTime(d.Clock.Now())}
		if len(persons) > 0 {
			where = " AND u.person_id IN (" + strings.TrimSuffix(strings.Repeat("?,", len(persons)), ",") + ")"
			for _, id := range persons {
				args = append(args, id)
			}
		}
		rows, err := d.DB.QueryContext(r.Context(), `SELECT s.id, s.user_account_id, s.device_label, s.user_agent, s.last_seen_at, s.expires_at, s.created_at, s.updated_at, u.username
			FROM device_session s JOIN user_account u ON u.id = s.user_account_id
			WHERE s.revoked_at IS NULL AND s.expires_at > ?`+where+" "+
			p.OrderBy(map[string]string{"last_seen_at": "s.last_seen_at", "created_at": "s.created_at", "expires_at": "s.expires_at", "id": "s.id"})+" "+p.LimitOffset(), args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		cur := api.Session(r).SessionID
		items := []api.DeviceSession{}
		for rows.Next() {
			var it api.DeviceSession
			var label, ua sql.NullString
			var user string
			if err := rows.Scan(&it.ID, &it.UserAccountID, &label, &ua, &it.LastSeenAt, &it.ExpiresAt, &it.CreatedAt, &it.UpdatedAt, &user); err != nil {
				return err
			}
			it.LastSeenAt, it.ExpiresAt = api.APITime(it.LastSeenAt), api.APITime(it.ExpiresAt)
			it.CreatedAt, it.UpdatedAt = api.APITime(it.CreatedAt), api.APITime(it.UpdatedAt)
			if label.Valid {
				it.DeviceLabel = api.Ptr(label.String)
			}
			if ua.Valid {
				it.UserAgent = api.Ptr(ua.String)
			}
			it.IsCurrent, it.Username = api.Ptr(it.ID == cur), api.Ptr(user)
			items = append(items, it)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// revokeOne cierra una sesión y la anota en audit_log.
func revokeOne(ctx context.Context, tx *sql.Tx, r *http.Request, id int64, now time.Time) error {
	before, err := dominio.RowJSON(ctx, tx, "device_session", id)
	if err != nil {
		return err
	}
	if err := dominio.RevokeSession(ctx, tx, id, now); err != nil {
		return err
	}
	after, err := dominio.RowJSON(ctx, tx, "device_session", id)
	if err != nil {
		return err
	}
	return dominio.Audit(ctx, tx, api.Actor(r), "update", "device_session", id, redactToken(before), redactToken(after), api.RequestID(r))
}

// redactToken quita token_hash de la copia que va a audit_log.
func redactToken(row json.RawMessage) json.RawMessage {
	var m map[string]any
	if json.Unmarshal(row, &m) != nil {
		return row
	}
	delete(m, "token_hash")
	b, _ := json.Marshal(m)
	return b
}

func revokeSession(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "session_id")
		if err != nil {
			return err
		}
		ctx, d := r.Context(), reg.Deps()
		now := d.Clock.Now()
		err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			s, err := dominio.GetSession(ctx, tx, id)
			if err != nil || s.RevokedAt.Valid {
				return dominio.NotFound("Esa sesión no existe o ya estaba cerrada.")
			}
			return revokeOne(ctx, tx, r, id, now)
		})
		if err != nil {
			return err
		}
		if id == api.Session(r).SessionID {
			clearCookie(w)
		}
		return api.NoContent(w)
	}
}

func revokeSessions(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.SessionRevokeRequest
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var v api.Validator
		v.Check(in.Scope == "user_except_current" || in.Scope == "all_except_current" || in.Scope == "all", "scope", "invalid_format",
			"Elige qué sesiones cerrar: las de una cuenta salvo la actual, todas salvo la actual o todas.")
		if in.Scope == "user_except_current" {
			v.Required(in.UserAccountID != nil, "user_account_id")
		}
		if err := v.Err(); err != nil {
			return err
		}
		ctx, d := r.Context(), reg.Deps()
		now, cur := d.Clock.Now(), api.Session(r).SessionID
		q, args := `SELECT id FROM device_session WHERE revoked_at IS NULL AND expires_at > ?`, []any{dominio.FormatDBTime(now)}
		switch in.Scope {
		case "user_except_current":
			if _, err := dominio.UserByID(ctx, d.DB, *in.UserAccountID); err != nil {
				return dominio.Invalid(dominio.FieldIssue{Field: "user_account_id", Code: "not_found", Message: "Esa cuenta no existe."})
			}
			q, args = q+` AND user_account_id = ? AND id <> ?`, append(args, *in.UserAccountID, cur)
		case "all_except_current":
			q, args = q+` AND id <> ?`, append(args, cur)
		}
		var n int64
		err := dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			rows, err := tx.QueryContext(ctx, q, args...)
			if err != nil {
				return err
			}
			var ids []int64
			for rows.Next() {
				var id int64
				if err := rows.Scan(&id); err != nil {
					rows.Close()
					return err
				}
				ids = append(ids, id)
			}
			rows.Close()
			for _, id := range ids {
				if err := revokeOne(ctx, tx, r, id, now); err != nil {
					return err
				}
				n++
			}
			return nil
		})
		if err != nil {
			return err
		}
		if in.Scope == "all" {
			clearCookie(w)
		}
		return api.WriteJSON(w, http.StatusOK, api.SessionRevokeResult{RevokedCount: n})
	}
}
