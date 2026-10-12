package dominio

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

// UserAccount es una de las 2 cuentas co-titulares (Axel y Ximena; decisión A). No hay función para crear cuentas:
// Karen no tiene cuenta (S6-05); las cuentas las crea la semilla.
type UserAccount struct {
	ID           int64
	PersonID     int64
	Username     string
	DisplayName  string // person.display_name
	PasswordHash string
	IsActive     bool
	LastLoginAt  sql.NullString
	CreatedAt    string
	UpdatedAt    string
}

const userCols = `u.id, u.person_id, u.username, p.display_name, u.password_hash, u.is_active, u.last_login_at, u.created_at, u.updated_at`

func scanUser(r interface{ Scan(...any) error }) (UserAccount, error) {
	var u UserAccount
	err := r.Scan(&u.ID, &u.PersonID, &u.Username, &u.DisplayName, &u.PasswordHash, &u.IsActive, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

// UserByUsername busca la cuenta por usuario (sin distinguir mayúsculas).
func UserByUsername(ctx context.Context, q Querier, username string) (UserAccount, error) {
	return scanUser(q.QueryRowContext(ctx, `SELECT `+userCols+` FROM user_account u JOIN person p ON p.id = u.person_id WHERE u.username = ?`, username))
}

// UserByID busca la cuenta por id.
func UserByID(ctx context.Context, q Querier, id int64) (UserAccount, error) {
	return scanUser(q.QueryRowContext(ctx, `SELECT `+userCols+` FROM user_account u JOIN person p ON p.id = u.person_id WHERE u.id = ?`, id))
}

// ListUsers devuelve las cuentas ordenadas por id.
func ListUsers(ctx context.Context, q Querier) ([]UserAccount, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+userCols+` FROM user_account u JOIN person p ON p.id = u.person_id ORDER BY u.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserAccount
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// SetPasswordHash guarda un hash argon2id nuevo (la contraseña en claro nunca llega aquí).
func SetPasswordHash(ctx context.Context, q Querier, userID int64, hash string, at time.Time) error {
	_, err := q.ExecContext(ctx, `UPDATE user_account SET password_hash = ?, password_changed_at = ? WHERE id = ?`, hash, FormatDBTime(at), userID)
	return err
}

// TouchLogin anota el último inicio de sesión correcto.
func TouchLogin(ctx context.Context, q Querier, userID int64, at time.Time) error {
	_, err := q.ExecContext(ctx, `UPDATE user_account SET last_login_at = ? WHERE id = ?`, FormatDBTime(at), userID)
	return err
}

// Argon2Params son los parámetros de argon2id. PROVISIONAL (K31, puerta 4): DefaultArgon2 sigue la segunda opción
// recomendada de RFC 9106 §4 (t=3, m=64 MiB, p=4, sal de 16 bytes, clave de 32 bytes).
type Argon2Params struct {
	Time    uint32
	Memory  uint32 // KiB
	Threads uint8
	SaltLen int
	KeyLen  uint32
}

// DefaultArgon2 es el valor propuesto (PROVISIONAL).
var DefaultArgon2 = Argon2Params{Time: 3, Memory: 64 * 1024, Threads: 4, SaltLen: 16, KeyLen: 32}

// HashPassword devuelve el hash argon2id en formato PHC: $argon2id$v=19$m=…,t=…,p=…$sal$hash.
func HashPassword(password string, p Argon2Params) (string, error) {
	salt := make([]byte, p.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, p.KeyLen)
	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Time, p.Threads, b64.EncodeToString(salt), b64.EncodeToString(key)), nil
}

// VerifyPassword compara en tiempo constante una contraseña con un hash PHC argon2id.
func VerifyPassword(password, phc string) (bool, error) {
	parts := strings.Split(phc, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("hash con formato desconocido")
	}
	var v int
	var p Argon2Params
	if _, err := fmt.Sscanf(parts[2], "v=%d", &v); err != nil || v != argon2.Version {
		return false, errors.New("versión de argon2 no admitida")
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Time, &p.Threads); err != nil {
		return false, errors.New("parámetros de argon2 ilegibles")
	}
	b64 := base64.RawStdEncoding
	salt, err1 := b64.DecodeString(parts[4])
	want, err2 := b64.DecodeString(parts[5])
	if err1 != nil || err2 != nil {
		return false, errors.New("sal o hash ilegibles")
	}
	got := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1, nil
}

// NewSessionToken genera un token de sesión aleatorio (32 bytes, base64url) y su SHA-256 en hexadecimal.
// El token en claro solo viaja en la cookie fh_session; la base guarda el hash (C1 §2).
func NewSessionToken() (token, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token = base64.RawURLEncoding.EncodeToString(b)
	return token, HashToken(token), nil
}

// HashToken da el SHA-256 en hexadecimal (64 caracteres) de un token.
func HashToken(token string) string {
	s := sha256.Sum256([]byte(token))
	return hex.EncodeToString(s[:])
}

// DeviceSession es una sesión larga por dispositivo.
type DeviceSession struct {
	ID            int64
	UserAccountID int64
	DeviceLabel   sql.NullString
	UserAgent     sql.NullString
	IPAddress     sql.NullString
	LastSeenAt    string
	ExpiresAt     string
	RevokedAt     sql.NullString
	CreatedAt     string
}

const sessCols = `id, user_account_id, device_label, user_agent, ip_address, last_seen_at, expires_at, revoked_at, created_at`

func scanSession(r interface{ Scan(...any) error }) (DeviceSession, error) {
	var s DeviceSession
	err := r.Scan(&s.ID, &s.UserAccountID, &s.DeviceLabel, &s.UserAgent, &s.IPAddress, &s.LastSeenAt, &s.ExpiresAt, &s.RevokedAt, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return s, ErrNotFound
	}
	return s, err
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// CreateSession guarda una sesión nueva con vencimiento now + maxAge.
func CreateSession(ctx context.Context, q Querier, userID int64, tokenHash, label, userAgent, ip string, now time.Time, maxAge time.Duration) (int64, error) {
	t := FormatDBTime(now)
	res, err := q.ExecContext(ctx, `INSERT INTO device_session (user_account_id, token_hash, device_label, user_agent, ip_address, last_seen_at, expires_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, tokenHash, nullable(label), nullable(userAgent), nullable(ip), t, FormatDBTime(now.Add(maxAge)), t, t)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ActiveSession busca la sesión vigente por hash del token: no cerrada, no vencida y con cuenta activa. Si no, ErrNotFound.
func ActiveSession(ctx context.Context, q Querier, tokenHash string, now time.Time) (DeviceSession, error) {
	return scanSession(q.QueryRowContext(ctx, `SELECT s.`+strings.ReplaceAll(sessCols, ", ", ", s.")+` FROM device_session s JOIN user_account u ON u.id = s.user_account_id
		WHERE s.token_hash = ? AND s.revoked_at IS NULL AND s.expires_at > ? AND u.is_active = 1`, tokenHash, FormatDBTime(now)))
}

// RenewSession mueve el vencimiento a now + maxAge con cada uso (K11) y anota último uso e IP.
func RenewSession(ctx context.Context, q Querier, id int64, ip string, now time.Time, maxAge time.Duration) error {
	_, err := q.ExecContext(ctx, `UPDATE device_session SET last_seen_at = ?, expires_at = ?, ip_address = coalesce(?, ip_address) WHERE id = ?`,
		FormatDBTime(now), FormatDBTime(now.Add(maxAge)), nullable(ip), id)
	return err
}

// GetSession lee una sesión por id (cerrada o no).
func GetSession(ctx context.Context, q Querier, id int64) (DeviceSession, error) {
	return scanSession(q.QueryRowContext(ctx, `SELECT `+sessCols+` FROM device_session WHERE id = ?`, id))
}

// ListSessions lista las sesiones vigentes de ambas cuentas (C1 §2), la de uso más reciente primero.
func ListSessions(ctx context.Context, q Querier, now time.Time) ([]DeviceSession, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+sessCols+` FROM device_session WHERE revoked_at IS NULL AND expires_at > ? ORDER BY last_seen_at DESC, id DESC`, FormatDBTime(now))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DeviceSession
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// RevokeSession cierra una sesión. Devuelve ErrNotFound si no existe o ya estaba cerrada.
func RevokeSession(ctx context.Context, q Querier, id int64, now time.Time) error {
	res, err := q.ExecContext(ctx, `UPDATE device_session SET revoked_at = ? WHERE id = ? AND revoked_at IS NULL`, FormatDBTime(now), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// RevokeOtherSessions cierra todas las sesiones vigentes salvo exceptID (0 = ninguna). Devuelve cuántas cerró.
func RevokeOtherSessions(ctx context.Context, q Querier, exceptID int64, now time.Time) (int64, error) {
	res, err := q.ExecContext(ctx, `UPDATE device_session SET revoked_at = ? WHERE revoked_at IS NULL AND id <> ?`, FormatDBTime(now), exceptID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// RecordLoginAttempt anota un intento de inicio de sesión. Un intento rechazado por bloqueo (429) NO se anota como
// fallo: si se anotara, el bloqueo nunca se levantaría con la ventana deslizante (S6-06 caso b).
func RecordLoginAttempt(ctx context.Context, q Querier, username string, userID *int64, ip string, ok bool, at time.Time) error {
	t := FormatDBTime(at)
	_, err := q.ExecContext(ctx, `INSERT INTO login_attempt (username, user_account_id, ip_address, succeeded, attempted_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		username, userID, nullable(ip), ok, t, t, t)
	return err
}

// Lockout es el estado del límite de intentos.
type Lockout struct {
	Locked  bool
	RetryAt time.Time // cuándo vuelve a haber menos de max fallos en la ventana
}

// CheckLockout aplica la ventana deslizante (K11): fallos con attempted_at > now − window, contados por usuario y por IP
// (lo que llegue primero a max). Con n ≥ max fallos, RetryAt = el fallo (n − max + 1)-ésimo más antiguo + window, que con
// n = max es «el fallo más antiguo de la ventana + 15 minutos» de C1 §2.
func CheckLockout(ctx context.Context, q Querier, username, ip string, now time.Time, max int, window time.Duration) (Lockout, error) {
	since := FormatDBTime(now.Add(-window))
	var out Lockout
	check := func(where string, arg string) error {
		rows, err := q.QueryContext(ctx, `SELECT attempted_at FROM login_attempt WHERE succeeded = 0 AND attempted_at > ? AND `+where+` ORDER BY attempted_at`, since, arg)
		if err != nil {
			return err
		}
		defer rows.Close()
		var ts []time.Time
		for rows.Next() {
			var s string
			if err := rows.Scan(&s); err != nil {
				return err
			}
			t, err := ParseDBTime(s)
			if err != nil {
				return err
			}
			ts = append(ts, t)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if max > 0 && len(ts) >= max {
			r := ts[len(ts)-max].Add(window)
			if !out.Locked || r.After(out.RetryAt) {
				out.RetryAt = r
			}
			out.Locked = true
		}
		return nil
	}
	if err := check("username = ?", username); err != nil {
		return out, err
	}
	if ip != "" {
		if err := check("ip_address = ?", ip); err != nil {
			return out, err
		}
	}
	return out, nil
}
