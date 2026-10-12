package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/ia"
	"finanzas-hogar/internal/jobs"
)

// SessionCookie es el nombre de la cookie de sesión (securitySchemes.cookieAuth de C1).
const SessionCookie = "fh_session"

// ServiceAuth verifica la credencial de servicio del agente externo (Authorization: Bearer, K53). PROVISIONAL: dónde
// vive la credencial está abierto (bloqueo S5-agente-externo, pregunta e); sin verificador, esas 4 operaciones dan 401.
type ServiceAuth interface {
	VerifyServiceToken(ctx context.Context, token string) (bool, error)
}

// Config son rutas y límites del servidor. Los límites marcados PROVISIONAL se aprueban en la puerta 4 (K31).
type Config struct {
	DataDir       string // carpeta de datos (adjuntos, reportes)
	BackupDir     string // carpeta de respaldos (solo local)
	TmpDir        string // temporales (restauración de prueba, cargas grandes)
	OpenAPIPath   string // docs/contratos/openapi.yaml que sirve getOpenApiDocument
	FrontendDir   string // frontend/dist (PWA, S6-12); vacío = no sirve interfaz
	MaxJSONBytes  int64  // PROVISIONAL: 8 MiB por cuerpo JSON
	IdemMemBytes  int64  // PROVISIONAL: 32 MiB en memoria por cuerpo idempotente; lo demás va a un temporal en TmpDir
	SecureHeaders bool   // PROVISIONAL (K31): cabeceras de seguridad; true en serve
}

// Deps son las dependencias que reciben las áreas por medio del Registry.
type Deps struct {
	DB      *sql.DB
	Clock   dominio.Clock
	Log     *slog.Logger
	IA      *ia.Registry    // proveedores de IA, apagados por defecto (S6-13)
	Jobs    *jobs.Scheduler // planificador; RunNow encola una corrida
	Service ServiceAuth     // nil = el agente externo no puede entrar (PROVISIONAL)
	Config  Config
}

// HandlerFunc es la forma de un manejador de área: escribe la respuesta y devuelve nil, o devuelve un error
// (*dominio.Error, *Problem o cualquier otro) que la API convierte en el esquema único de error de C1.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// Registry junta los manejadores de todas las áreas; cada área exporta Register(reg *api.Registry).
type Registry struct {
	deps  *Deps
	impl  map[string]HandlerFunc
	known map[string]OperationMeta
}

var operationIndex = func() map[string]OperationMeta {
	m := make(map[string]OperationMeta, len(Operations))
	for _, o := range Operations {
		m[o.OperationID] = o
	}
	return m
}()

// Meta devuelve los metadatos del contrato de una operación.
func Meta(operationID string) (OperationMeta, bool) {
	m, ok := operationIndex[operationID]
	return m, ok
}

// NewRegistry crea un registro vacío con las dependencias dadas (con valores por omisión donde falten).
func NewRegistry(d *Deps) *Registry {
	if d.Clock == nil {
		d.Clock = dominio.SystemClock{}
	}
	if d.Log == nil {
		d.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if d.IA == nil {
		d.IA = ia.NewRegistry(d.DB)
	}
	if d.Config.MaxJSONBytes == 0 {
		d.Config.MaxJSONBytes = 8 << 20
	}
	if d.Config.IdemMemBytes == 0 {
		d.Config.IdemMemBytes = 32 << 20
	}
	if d.Config.TmpDir == "" {
		d.Config.TmpDir = os.TempDir()
	}
	return &Registry{deps: d, impl: map[string]HandlerFunc{}, known: operationIndex}
}

// Deps devuelve las dependencias compartidas.
func (g *Registry) Deps() *Deps { return g.deps }

// Handle registra el manejador de una operación de C1. Falla (panic) si la operación no está en el contrato o ya
// tiene manejador: dos áreas no pueden reclamar la misma operación.
func (g *Registry) Handle(operationID string, h HandlerFunc) {
	if _, ok := g.known[operationID]; !ok {
		panic("api: la operación " + operationID + " no está en docs/contratos/openapi.yaml")
	}
	if _, dup := g.impl[operationID]; dup {
		panic("api: la operación " + operationID + " ya tiene manejador")
	}
	g.impl[operationID] = h
}

// Has dice si la operación ya tiene manejador.
func (g *Registry) Has(operationID string) bool { _, ok := g.impl[operationID]; return ok }

// Implemented lista las operaciones con manejador.
func (g *Registry) Implemented() []string {
	out := make([]string, 0, len(g.impl))
	for k := range g.impl {
		out = append(out, k)
	}
	return out
}

// Handler arma el http.Handler completo: las 216 rutas de routes.go (las que no tienen manejador responden 501),
// con sesión, idempotencia, errores únicos, 404/405 de C1, cabeceras y la PWA en las rutas que no son de la API.
func (g *Registry) Handler() http.Handler {
	mux := http.NewServeMux()
	impl := make(map[string]http.Handler, len(Routes))
	for _, rt := range Routes {
		meta := g.known[rt.OperationID]
		h, ok := g.impl[rt.OperationID]
		if !ok {
			h = func(w http.ResponseWriter, r *http.Request) error {
				return &Problem{Status: http.StatusNotImplemented, Code: "not_implemented", Message: "Esta operación todavía no está implementada."}
			}
		}
		impl[rt.OperationID] = g.wrap(meta, h)
	}
	Register(mux, impl)
	if g.deps.Config.FrontendDir != "" {
		mux.Handle("GET /", staticHandler(g.deps.Config.FrontendDir))
	}
	return g.outer(mux)
}

type ctxKey int

const (
	ctxSession ctxKey = iota + 1
	ctxIdemKey
	ctxOperation
)

// SessionInfo es la sesión de la petición (cookie fh_session).
type SessionInfo struct {
	SessionID int64
	UserID    int64
	PersonID  int64
	Username  string
	Token     string // token en claro de la cookie; solo para renovar la cookie, nunca se registra
}

// Session devuelve la sesión de la petición (nil en operaciones públicas o del agente externo).
func Session(r *http.Request) *SessionInfo {
	s, _ := r.Context().Value(ctxSession).(*SessionInfo)
	return s
}

// Actor devuelve el user_account_id para audit_log (nil = sistema o agente externo).
func Actor(r *http.Request) *int64 {
	if s := Session(r); s != nil {
		return &s.UserID
	}
	return nil
}

// OperationID devuelve el operationId de la petición.
func OperationID(r *http.Request) string { s, _ := r.Context().Value(ctxOperation).(string); return s }

// IdempotencyKey devuelve la Idempotency-Key validada de la petición ("" si la operación no la lleva). Gastos y pesos
// la guardan además en su renglón (K9).
func IdempotencyKey(r *http.Request) string { s, _ := r.Context().Value(ctxIdemKey).(string); return s }

// MarkReplayed pone Idempotent-Replayed: true (p. ej. createExpense al encontrar la clave en el renglón pasadas 24 h).
func MarkReplayed(w http.ResponseWriter) { w.Header().Set("Idempotent-Replayed", "true") }

// SetSessionCookie fija la cookie de sesión: HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age (C1 §2).
func SetSessionCookie(w http.ResponseWriter, token string, maxAge time.Duration) {
	http.SetCookie(w, &http.Cookie{Name: SessionCookie, Value: token, Path: "/", MaxAge: int(maxAge / time.Second), HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})
}

// ClearSessionCookie borra la cookie (logout).
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: SessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})
}

// SessionMaxAge lee session.max_age_days (90 al sembrar, K11).
func SessionMaxAge(ctx context.Context, q dominio.Querier) (time.Duration, error) {
	n, err := dominio.SettingInt(ctx, q, "session.max_age_days")
	return time.Duration(n) * 24 * time.Hour, err
}

// ClientIP es la IP de origen (sin puerto).
func ClientIP(r *http.Request) string {
	h, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return h
}

func (g *Registry) wrap(meta OperationMeta, h HandlerFunc) http.Handler {
	d := g.deps
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), ctxOperation, meta.OperationID))
		defer func() {
			if v := recover(); v != nil {
				d.Log.Error("pánico en manejador", "operation", meta.OperationID, "request_id", RequestID(r), "panic", fmt.Sprint(v), "stack", string(debug.Stack()))
				Fail(w, r, errors.New("pánico"))
			}
		}()
		var err error
		switch meta.Auth {
		case "session":
			r, err = g.authSession(w, r)
		case "service":
			err = g.authService(r)
		}
		if err == nil && meta.Idempotent {
			err = g.idempotent(w, r, meta, h)
		} else if err == nil {
			err = h(w, r)
		}
		if err != nil {
			if _, ok := dominio.AsError(err); !ok {
				var p *Problem
				if !errors.As(err, &p) {
					d.Log.Error("error interno", "operation", meta.OperationID, "request_id", RequestID(r), "error", err.Error())
				}
			}
			Fail(w, r, err)
		}
	})
}

func (g *Registry) authSession(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	unauth := &Problem{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Tu sesión terminó o se cerró. Inicia sesión de nuevo."}
	c, err := r.Cookie(SessionCookie)
	if err != nil || c.Value == "" {
		return r, unauth
	}
	ctx, d := r.Context(), g.deps
	now := d.Clock.Now()
	s, err := dominio.ActiveSession(ctx, d.DB, dominio.HashToken(c.Value), now)
	if errors.Is(err, dominio.ErrNotFound) {
		return r, unauth
	}
	if err != nil {
		return r, err
	}
	u, err := dominio.UserByID(ctx, d.DB, s.UserAccountID)
	if err != nil {
		return r, err
	}
	maxAge, err := SessionMaxAge(ctx, d.DB)
	if err != nil {
		return r, err
	}
	if err := dominio.RenewSession(ctx, d.DB, s.ID, ClientIP(r), now, maxAge); err != nil {
		return r, err
	}
	SetSessionCookie(w, c.Value, maxAge)
	info := &SessionInfo{SessionID: s.ID, UserID: u.ID, PersonID: u.PersonID, Username: u.Username, Token: c.Value}
	return r.WithContext(context.WithValue(ctx, ctxSession, info)), nil
}

var reBearer = regexp.MustCompile(`^Bearer ([A-Za-z0-9._~+/=-]{16,512})$`)

func (g *Registry) authService(r *http.Request) error {
	unauth := &Problem{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Esta operación necesita la credencial del agente externo."}
	m := reBearer.FindStringSubmatch(r.Header.Get("Authorization"))
	if m == nil || g.deps.Service == nil {
		return unauth
	}
	ok, err := g.deps.Service.VerifyServiceToken(r.Context(), m[1])
	if err != nil {
		return err
	}
	if !ok {
		return unauth
	}
	return nil
}

var reIdemKey = regexp.MustCompile(`^[A-Za-z0-9_-]{16,128}$`)

// idempotent aplica C1 §1 punto 8: misma clave y mismo cuerpo → la respuesta guardada con Idempotent-Replayed: true
// (201 se repite como 200 si la operación declara 200, K9); otro cuerpo → 409 idempotency_key_reused. Solo se guardan
// respuestas 2xx; una respuesta de error libera la clave para reintentar.
func (g *Registry) idempotent(w http.ResponseWriter, r *http.Request, meta OperationMeta, h HandlerFunc) error {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		return BadRequest("Falta el encabezado Idempotency-Key.")
	}
	if !reIdemKey.MatchString(key) {
		return BadRequest("La Idempotency-Key debe tener de 16 a 128 letras sin acento, dígitos, guiones o guiones bajos.")
	}
	r = r.WithContext(context.WithValue(r.Context(), ctxIdemKey, key))
	s := Session(r)
	if s == nil { // agente externo: el esquema 1.4 no tiene a nombre de quién guardar la clave (PROVISIONAL, S5-agente-externo)
		return h(w, r)
	}
	sum, cleanup, err := g.hashBody(r)
	defer cleanup()
	if err != nil {
		return err
	}
	ctx, d := r.Context(), g.deps
	now := d.Clock.Now()
	hours, err := dominio.SettingInt(ctx, d.DB, "idempotency.retention_hours")
	if err != nil {
		return err
	}
	notBefore := now.Add(-time.Duration(hours) * time.Hour)
	rec, err := dominio.FindIdempotency(ctx, d.DB, s.UserID, key, notBefore)
	switch {
	case err == nil && rec.RequestSHA != sum:
		return &Problem{Status: http.StatusConflict, Code: "idempotency_key_reused", Message: "Esta clave de captura ya se usó con otros datos. Vuelve a capturar el registro."}
	case err == nil && !rec.Status.Valid:
		return &Problem{Status: http.StatusConflict, Code: "conflict", Message: "Esta captura ya se está procesando. Espera un momento y vuelve a intentar."}
	case err == nil:
		status := int(rec.Status.Int64)
		if status == http.StatusCreated && containsInt(meta.Statuses, http.StatusOK) {
			status = http.StatusOK
		}
		MarkReplayed(w)
		if rec.Body.Valid {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, rec.Body.String)
		return nil
	case !errors.Is(err, dominio.ErrNotFound):
		return err
	}
	id, err := dominio.ReserveIdempotency(ctx, d.DB, s.UserID, key, r.Method, r.URL.Path, sum, now, notBefore)
	if err != nil {
		return err
	}
	cw := &captureWriter{ResponseWriter: w}
	done := false
	defer func() {
		if !done {
			_ = dominio.ReleaseIdempotency(context.WithoutCancel(ctx), d.DB, id)
		}
	}()
	if err := h(cw, r); err != nil {
		return err
	}
	if cw.status >= 200 && cw.status < 300 && (cw.buf.Len() == 0 || json.Valid(cw.buf.Bytes())) {
		done = dominio.CompleteIdempotency(context.WithoutCancel(ctx), d.DB, id, cw.status, cw.buf.Bytes()) == nil
	}
	return nil
}

// hashBody calcula el SHA-256 del cuerpo y lo deja releíble (en memoria hasta IdemMemBytes; el resto, en TmpDir).
func (g *Registry) hashBody(r *http.Request) (string, func(), error) {
	noop := func() {}
	if r.Body == nil {
		return dominio.HashToken(""), noop, nil
	}
	limit := g.deps.Config.IdemMemBytes
	var mem bytes.Buffer
	n, err := io.CopyN(&mem, r.Body, limit+1)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", noop, BadRequest("No se pudo leer el cuerpo de la petición.")
	}
	if n <= limit {
		r.Body = io.NopCloser(bytes.NewReader(mem.Bytes()))
		return dominio.HashToken(mem.String()), noop, nil
	}
	f, err := os.CreateTemp(g.deps.Config.TmpDir, "cuerpo-")
	if err != nil {
		return "", noop, err
	}
	cleanup := func() { f.Close(); os.Remove(f.Name()) }
	if _, err := io.Copy(f, io.MultiReader(&mem, r.Body)); err != nil {
		return "", cleanup, BadRequest("No se pudo leer el cuerpo de la petición.")
	}
	sum, err := fileSum(f)
	if err != nil {
		return "", cleanup, err
	}
	r.Body = f
	return sum, cleanup, nil
}

type captureWriter struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (c *captureWriter) WriteHeader(code int) {
	if c.status == 0 {
		c.status = code
	}
	c.ResponseWriter.WriteHeader(code)
}

func (c *captureWriter) Write(b []byte) (int, error) {
	if c.status == 0 {
		c.status = http.StatusOK
	}
	c.buf.Write(b)
	return c.ResponseWriter.Write(b)
}

func containsInt(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// outer agrega X-Request-Id a todo, el registro de accesos (sin cookies, cuerpos ni tokens), las cabeceras
// PROVISIONALES y los 404/405 con el esquema único de error.
func (g *Registry) outer(mux *http.ServeMux) http.Handler {
	d := g.deps
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if id := r.Header.Get("X-Request-Id"); id == "" || len(id) > 64 {
			r.Header.Set("X-Request-Id", newRequestID())
		}
		w.Header().Set("X-Request-Id", r.Header.Get("X-Request-Id"))
		if d.Config.SecureHeaders {
			securityHeaders(w.Header(), r.URL.Path)
		}
		sw := &statusWriter{ResponseWriter: w}
		var out http.ResponseWriter = sw
		if _, pattern := mux.Handler(r); pattern == "" || (pattern == "GET /" && isAPIPath(r.URL.Path)) {
			out = &contractErrorWriter{ResponseWriter: sw, r: r}
			if pattern == "GET /" {
				Fail(out, r, &Problem{Status: http.StatusNotFound, Code: "not_found", Message: "La ruta no existe."})
				d.Log.Info("petición", "method", r.Method, "path", r.URL.Path, "status", sw.code(), "ms", time.Since(start).Milliseconds(), "request_id", r.Header.Get("X-Request-Id"))
				return
			}
		}
		mux.ServeHTTP(out, r)
		d.Log.Info("petición", "method", r.Method, "path", r.URL.Path, "operation", OperationID(r), "status", sw.code(), "ms", time.Since(start).Milliseconds(), "request_id", r.Header.Get("X-Request-Id"))
	})
}

func isAPIPath(p string) bool {
	return p == "/health" || p == "/openapi.yaml" || len(p) >= 5 && p[:5] == "/api/"
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(c int) {
	if s.status == 0 {
		s.status = c
	}
	s.ResponseWriter.WriteHeader(c)
}

func (s *statusWriter) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

func (s *statusWriter) code() int {
	if s.status == 0 {
		return http.StatusOK
	}
	return s.status
}

type contractErrorWriter struct {
	http.ResponseWriter
	r        *http.Request
	replaced bool
}

func (c *contractErrorWriter) WriteHeader(code int) {
	switch code {
	case http.StatusNotFound:
		c.replaced = true
		WriteError(c.ResponseWriter, c.r, code, "not_found", "La ruta no existe.")
	case http.StatusMethodNotAllowed:
		c.replaced = true
		WriteError(c.ResponseWriter, c.r, code, "method_not_allowed", "Método no permitido en esta ruta.")
	default:
		c.ResponseWriter.WriteHeader(code)
	}
}

func (c *contractErrorWriter) Write(b []byte) (int, error) {
	if c.replaced {
		return len(b), nil
	}
	return c.ResponseWriter.Write(b)
}

// securityHeaders son las cabeceras propuestas (PROVISIONAL, K31 punto «cabeceras», puerta 4).
func securityHeaders(h http.Header, path string) {
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "no-referrer")
	h.Set("X-Frame-Options", "DENY")
	if isAPIPath(path) {
		h.Set("Cache-Control", "no-store")
		return
	}
	h.Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; connect-src 'self'; font-src 'self'; worker-src 'self'; manifest-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
}
