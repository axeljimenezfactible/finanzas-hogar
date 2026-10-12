package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/ajustes"
	"finanzas-hogar/internal/api/avisos"
	"finanzas-hogar/internal/api/comidas"
	"finanzas-hogar/internal/api/comprasa"
	"finanzas-hogar/internal/api/comprasb"
	"finanzas-hogar/internal/api/creditos"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/api/presupuesto"
	reportesapi "finanzas-hogar/internal/api/reportes"
	"finanzas-hogar/internal/api/tablero"
	"finanzas-hogar/internal/auth"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/telegram"
	"finanzas-hogar/internal/testutil"
)

func init() { auth.HashParams = testutil.FastArgon2 }

func creds(user, pw string) map[string]any {
	return map[string]any{"username": user, "password": pw, "device_label": "Teléfono de prueba"}
}

func sessionCookie(t *testing.T, res *testutil.Response) *http.Cookie {
	t.Helper()
	for _, c := range res.Cookies() {
		if c.Name == api.SessionCookie {
			return c
		}
	}
	t.Fatalf("sin cookie %s: %v", api.SessionCookie, res.Header)
	return nil
}

// login inicia sesión por la API (valida la respuesta contra openapi.yaml) y devuelve el cliente con su cookie.
func login(t *testing.T, s *testutil.Server, user, pw string) (*testutil.Client, api.CurrentSession, *testutil.Response) {
	t.Helper()
	res := s.Anonymous(t).Do("POST", "/api/v1/auth/login", creds(user, pw))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login %s: %d %s", user, res.StatusCode, res.Body)
	}
	var cs api.CurrentSession
	res.JSON(t, &cs)
	c := s.Anonymous(t)
	c.Token = sessionCookie(t, res).Value
	return c, cs, res
}

func code(t *testing.T, res *testutil.Response) string {
	t.Helper()
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	res.JSON(t, &e)
	return e.Code
}

func TestOp_login(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	_, cs, res := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	raw := strings.Join(res.Header.Values("Set-Cookie"), " | ")
	for _, want := range []string{"HttpOnly", "Secure", "SameSite=Strict", "Max-Age=7776000", "Path=/"} {
		if !strings.Contains(raw, want) {
			t.Errorf("Set-Cookie sin %s: %s", want, raw)
		}
	}
	if cs.Username != "axel" || cs.Participation != "holder" || cs.DeviceSessionID == 0 {
		t.Errorf("sesión: %+v", cs)
	}
	bad := s.Anonymous(t).Do("POST", "/api/v1/auth/login", creds(testutil.AxelUser, "otra-contrasena-mala"))
	if bad.StatusCode != 401 || code(t, bad) != "invalid_credentials" || len(bad.Cookies()) != 0 {
		t.Errorf("credenciales malas: %d %s", bad.StatusCode, bad.Body)
	}
	if !strings.Contains(string(bad.Body), "Usuario o contraseña incorrectos") {
		t.Errorf("mensaje genérico: %s", bad.Body)
	}
	if r := s.Anonymous(t).Do("POST", "/api/v1/auth/login", map[string]any{"username": "", "password": ""}); r.StatusCode != 422 {
		t.Errorf("vacíos: %d", r.StatusCode)
	}
	if r := s.Anonymous(t).Do("POST", "/api/v1/auth/login", `{"username": 1}`); r.StatusCode != 400 {
		t.Errorf("formato: %d", r.StatusCode)
	}
	rows, err := s.DB.Query(`SELECT substr(password_hash, 1, 10) FROM user_account`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p string
		if rows.Scan(&p); p != "$argon2id$" {
			t.Errorf("hash sin argon2id: %q", p)
		}
	}
}

func TestOp_logout(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c, _, _ := login(t, s, testutil.XimenaUser, testutil.FakeCredential)
	res := c.Do("POST", "/api/v1/auth/logout", nil)
	if res.StatusCode != 204 {
		t.Fatalf("logout: %d %s", res.StatusCode, res.Body)
	}
	if ck := sessionCookie(t, res); ck.MaxAge >= 0 || ck.Value != "" {
		t.Errorf("la cookie no se borró: %+v", ck)
	}
	if r := c.Do("GET", "/api/v1/auth/me", nil); r.StatusCode != 401 {
		t.Errorf("tras logout: %d", r.StatusCode)
	}
	if r := s.Anonymous(t).Do("POST", "/api/v1/auth/logout", nil); r.StatusCode != 401 {
		t.Errorf("sin sesión: %d", r.StatusCode)
	}
}

func TestOp_changeOwnPassword(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	r := c.Do("POST", "/api/v1/auth/password", map[string]any{"current_password": "no-es-la-actual", "new_password": "nueva-de-doce-x"})
	if r.StatusCode != 422 || code(t, r) != "invalid_current_password" {
		t.Errorf("actual mala: %d %s", r.StatusCode, r.Body)
	}
	r = c.Do("POST", "/api/v1/auth/password", map[string]any{"current_password": testutil.FakeCredential, "new_password": "once-carac1"}) // 11
	if r.StatusCode != 422 || code(t, r) != "validation_failed" || !strings.Contains(string(r.Body), "12 caracteres") {
		t.Errorf("11 caracteres: %d %s", r.StatusCode, r.Body)
	}
	nueva := "doce-caract1" // 12
	if r = c.Do("POST", "/api/v1/auth/password", map[string]any{"current_password": testutil.FakeCredential, "new_password": nueva}); r.StatusCode != 204 {
		t.Fatalf("12 caracteres: %d %s", r.StatusCode, r.Body)
	}
	login(t, s, testutil.AxelUser, nueva)
	var h string
	if err := s.DB.QueryRow(`SELECT password_hash FROM user_account WHERE username = 'axel'`).Scan(&h); err != nil || !strings.HasPrefix(h, "$argon2id$") {
		t.Errorf("hash nuevo: %q %v", h[:10], err)
	}
	var n int
	s.DB.QueryRow(`SELECT COUNT(*) FROM audit_log WHERE table_name = 'user_account' AND (coalesce(before_json, '') || coalesce(after_json, '')) LIKE '%argon2id%'`).Scan(&n)
	if n != 0 {
		t.Errorf("audit_log guarda el hash")
	}
	if r := s.Anonymous(t).Do("POST", "/api/v1/auth/password", map[string]any{"current_password": "a", "new_password": "b"}); r.StatusCode != 401 {
		t.Errorf("sin sesión: %d", r.StatusCode)
	}
}

func TestOp_listUsers(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	res := s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/users", nil)
	var pg api.UserAccountPage
	res.JSON(t, &pg)
	if res.StatusCode != 200 || len(pg.Items) != 2 || pg.NextCursor.Valid || strings.Contains(string(res.Body), "password") {
		t.Fatalf("%d %s", res.StatusCode, res.Body)
	}
	if r := s.SessionFor(t, testutil.AxelUser).Do("GET", "/api/v1/users?sort=password_hash", nil); r.StatusCode != 400 {
		t.Errorf("orden no permitido: %d", r.StatusCode)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/users", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

func TestOp_listSessions(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	a := s.SessionFor(t, testutil.AxelUser)
	s.SessionFor(t, testutil.XimenaUser)
	var pg api.DeviceSessionPage
	res := a.Do("GET", "/api/v1/sessions", nil)
	res.JSON(t, &pg)
	cur := 0
	for _, it := range pg.Items {
		if it.IsCurrent != nil && *it.IsCurrent {
			cur++
		}
	}
	if res.StatusCode != 200 || len(pg.Items) != 2 || cur != 1 {
		t.Fatalf("%d %s", res.StatusCode, res.Body)
	}
	res = a.Do("GET", fmt.Sprintf("/api/v1/sessions?person_id=%d", s.Fx.XimenaPersonID), nil)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].Username != "ximena" {
		t.Errorf("filtro por persona: %s", res.Body)
	}
	if r := a.Do("GET", "/api/v1/sessions?limit=201", nil); r.StatusCode != 400 {
		t.Errorf("límite 201: %d", r.StatusCode)
	}
}

func TestOp_revokeSession(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	a, _, _ := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	_, x, _ := login(t, s, testutil.XimenaUser, testutil.FakeCredential)
	path := fmt.Sprintf("/api/v1/sessions/%d", x.DeviceSessionID)
	if r := a.Do("DELETE", path, nil); r.StatusCode != 204 { // co-titulares: Axel cierra la de Ximena
		t.Fatalf("cerrar: %d %s", r.StatusCode, r.Body)
	}
	if r := a.Do("DELETE", path, nil); r.StatusCode != 404 {
		t.Errorf("ya cerrada: %d", r.StatusCode)
	}
	if r := a.Do("DELETE", "/api/v1/sessions/999", nil); r.StatusCode != 404 {
		t.Errorf("no existe: %d", r.StatusCode)
	}
	if r := a.Do("DELETE", "/api/v1/sessions/abc", nil); r.StatusCode != 400 {
		t.Errorf("id malo: %d", r.StatusCode)
	}
}

func TestOp_revokeSessions(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	a, _, _ := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	b, _, _ := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	x, _, _ := login(t, s, testutil.XimenaUser, testutil.FakeCredential)
	var out api.SessionRevokeResult
	if r := a.Do("POST", "/api/v1/sessions/revoke", map[string]any{"scope": "user_except_current"}); r.StatusCode != 422 {
		t.Errorf("sin cuenta: %d", r.StatusCode)
	}
	if r := a.Do("POST", "/api/v1/sessions/revoke", map[string]any{"scope": "todas"}); r.StatusCode != 422 {
		t.Errorf("alcance malo: %d", r.StatusCode)
	}
	r := a.Do("POST", "/api/v1/sessions/revoke", map[string]any{"scope": "user_except_current", "user_account_id": s.Fx.AxelUserID})
	r.JSON(t, &out)
	if r.StatusCode != 200 || out.RevokedCount != 1 || b.Do("GET", "/api/v1/auth/me", nil).StatusCode != 401 || x.Do("GET", "/api/v1/auth/me", nil).StatusCode != 200 {
		t.Fatalf("user_except_current: %d %s", r.StatusCode, r.Body)
	}
	r = a.Do("POST", "/api/v1/sessions/revoke", map[string]any{"scope": "all_except_current"})
	r.JSON(t, &out)
	if out.RevokedCount != 1 || x.Do("GET", "/api/v1/auth/me", nil).StatusCode != 401 {
		t.Fatalf("all_except_current: %s", r.Body)
	}
	r = a.Do("POST", "/api/v1/sessions/revoke", map[string]any{"scope": "all"})
	r.JSON(t, &out)
	if out.RevokedCount != 1 || a.Do("GET", "/api/v1/auth/me", nil).StatusCode != 401 {
		t.Fatalf("all: %s", r.Body)
	}
}

func TestOp_getCurrentSession(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var cs api.CurrentSession
	res := s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/auth/me", nil)
	res.JSON(t, &cs)
	if res.StatusCode != 200 || cs.Username != "ximena" || cs.Participation != "co_holder" || cs.DisplayName != "Ximena" {
		t.Fatalf("%d %s", res.StatusCode, res.Body)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/auth/me", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

// S6-05: exactamente 2 cuentas; Karen no tiene cuenta ni forma de crearla.
func TestCuentas_DosYKarenSinCuenta(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var n int
	if s.DB.QueryRow(`SELECT COUNT(*) FROM user_account`).Scan(&n); n != 2 {
		t.Fatalf("cuentas: %d", n)
	}
	if r := s.Anonymous(t).Do("POST", "/api/v1/auth/login", creds("karen", testutil.FakeCredential)); r.StatusCode != 401 {
		t.Errorf("Karen inicia sesión: %d", r.StatusCode)
	}
	if r := s.SessionFor(t, testutil.AxelUser).Raw("POST", "/api/v1/users", map[string]any{"username": "karen"}); r.StatusCode != 405 {
		t.Errorf("crear cuenta por la API: %d", r.StatusCode)
	}
	h, _ := dominio.HashPassword(testutil.FakeCredential, testutil.FastArgon2)
	if _, err := s.DB.Exec(`INSERT INTO user_account (person_id, username, password_hash) VALUES (?, 'karen', ?)`, s.Fx.KarenPersonID, h); err == nil {
		t.Error("la base aceptó una cuenta para Karen")
	}
}

// todasLasAreas es la misma lista que backend/cmd/finanzas/areas.go (el binario), salvo el cliente de Telegram de
// s7-avisos, que aquí va a un servidor falso con un token falso: la prueba nunca llama a Telegram.
func todasLasAreas(t *testing.T) []func(*api.Registry) {
	t.Helper()
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true,"result":{}}`))
	}))
	t.Cleanup(tg.Close)
	cfg := filepath.Join(testutil.TempDir(t), "openclaw.json")
	if err := os.WriteFile(cfg, []byte(`{"bots":{"multix":{"botToken":"1234567890:falso_falso_falso"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cl := &telegram.Client{Base: tg.URL, Token: telegram.FileToken(cfg), HTTP: tg.Client()}
	return []func(*api.Registry){auth.Register, gastos.Register, presupuesto.Register, comidas.Register, comprasa.Register,
		comprasb.Register, creditos.Register, ajustes.Register, avisos.RegisterWith(avisos.Options{Enviador: cl}), reportesapi.Register, tablero.Register}
}

// S6-05: Axel y Ximena tienen los mismos permisos: con todas las áreas registradas (las del binario), cada operación
// de C1 da el mismo código a las dos cuentas (dos servidores iguales, mismo orden) y ninguna responde 403 ni 501
// (hallazgo M1 de REVISION-backend-core-r1). Las que cierran la propia sesión o restauran van aparte.
func TestPermisos_MismasOperacionesAmbasCuentas(t *testing.T) {
	skip := map[string]bool{"logout": true, "revokeSession": true, "revokeSessions": true, "restoreBackup": true}
	param := regexp.MustCompile(`\{[a-z_]+\}`)
	run := func(user string) map[string]int {
		s := testutil.NewServer(t, todasLasAreas(t)...)
		c := s.SessionFor(t, user)
		out := map[string]int{}
		for i, rt := range api.Routes {
			if skip[rt.OperationID] {
				continue
			}
			p := param.ReplaceAllStringFunc(rt.Pattern, func(m string) string {
				if m == "{provider}" {
					return "local_model"
				}
				return "1"
			})
			var body any
			if rt.Method == "POST" || rt.Method == "PATCH" || rt.Method == "PUT" {
				body = "{}"
			}
			out[rt.OperationID] = c.Raw(rt.Method, p, body, "Idempotency-Key", fmt.Sprintf("permisos-%04d-abcdefgh", i)).StatusCode
		}
		return out
	}
	ax, xi := run(testutil.AxelUser), run(testutil.XimenaUser)
	if len(ax) != len(api.Routes)-len(skip) {
		t.Fatalf("operaciones recorridas: %d", len(ax))
	}
	for op, c := range ax {
		if xi[op] != c || c == 403 || c == 501 || xi[op] == 403 || xi[op] == 501 {
			t.Errorf("%s: Axel %d, Ximena %d", op, c, xi[op])
		}
	}
	t.Logf("%d operaciones con el mismo resultado para las dos cuentas", len(ax))
}

// --- S6-07: sesión de 90 días por dispositivo, renovable con el uso ---

func TestSesion_DosDispositivosYCierre(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	d1, _, _ := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	d2, s2, _ := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	if r := d1.Do("DELETE", fmt.Sprintf("/api/v1/sessions/%d", s2.DeviceSessionID), nil); r.StatusCode != 204 {
		t.Fatalf("cerrar: %d", r.StatusCode)
	}
	if r := d2.Do("GET", "/api/v1/auth/me", nil); r.StatusCode != 401 {
		t.Errorf("la cookie del dispositivo cerrado: %d", r.StatusCode)
	}
	if r := d1.Do("GET", "/api/v1/auth/me", nil); r.StatusCode != 200 {
		t.Errorf("el otro dispositivo sigue: %d", r.StatusCode)
	}
}

func TestSesion_90DiasRenovableConElUso(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	t0 := s.Clock.Now()
	day := 24 * time.Hour
	a, sa, res := login(t, s, testutil.AxelUser, testutil.FakeCredential)
	b, sb, _ := login(t, s, testutil.XimenaUser, testutil.FakeCredential)
	if ck := sessionCookie(t, res); ck.MaxAge != 90*24*3600 {
		t.Errorf("Max-Age inicial: %d", ck.MaxAge)
	}
	exp := func(id int64) string {
		var e string
		s.DB.QueryRow(`SELECT expires_at FROM device_session WHERE id = ?`, id).Scan(&e)
		return e
	}
	if exp(sa.DeviceSessionID) != dominio.FormatDBTime(t0.Add(90*day)) || exp(sb.DeviceSessionID) != dominio.FormatDBTime(t0.Add(90*day)) {
		t.Errorf("vencimiento inicial: %s", exp(sa.DeviceSessionID))
	}
	s.Clock.Set(t0.Add(60 * day))
	r := a.Do("GET", "/api/v1/auth/me", nil)
	if r.StatusCode != 200 || exp(sa.DeviceSessionID) != dominio.FormatDBTime(t0.Add(150*day)) || sessionCookie(t, r).MaxAge != 90*24*3600 {
		t.Errorf("día 60: %d, vence %s", r.StatusCode, exp(sa.DeviceSessionID))
	}
	s.Clock.Set(t0.Add(91 * day))
	if r := b.Do("GET", "/api/v1/auth/me", nil); r.StatusCode != 401 {
		t.Errorf("sin uso, día 91: %d", r.StatusCode)
	}
	if r := a.Do("GET", "/api/v1/auth/me", nil); r.StatusCode != 200 {
		t.Errorf("con uso el día 60, día 91: %d", r.StatusCode)
	}
}

// --- S6-06: 5 fallos en una ventana deslizante de 15 minutos, por usuario y por IP ---

func intento(t *testing.T, s *testutil.Server, user, pw string) *testutil.Response {
	t.Helper()
	return s.Anonymous(t).Do("POST", "/api/v1/auth/login", creds(user, pw))
}

func TestIntentos_A_CincoFallosEnCuarentaSegundos(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	t0 := s.Clock.Now()
	for i := 0; i < 5; i++ {
		s.Clock.Set(t0.Add(time.Duration(i*10) * time.Second))
		if r := intento(t, s, "axel", "mala"); r.StatusCode != 401 {
			t.Fatalf("fallo %d: %d", i, r.StatusCode)
		}
	}
	s.Clock.Set(t0.Add(14 * time.Minute))
	r := intento(t, s, "axel", testutil.FakeCredential)
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		RetryAt string `json:"retry_at"`
	}
	r.JSON(t, &e)
	want := dominio.FormatDateTime(t0.Add(15 * time.Minute))
	if r.StatusCode != 429 || e.Code != "too_many_attempts" || !strings.Contains(e.Message, "a las 12:15") {
		t.Errorf("minuto 14 con la contraseña correcta: %d %s (retry_at esperado %s)", r.StatusCode, r.Body, want)
	}
	s.Clock.Set(t0.Add(16 * time.Minute))
	if r := intento(t, s, "axel", testutil.FakeCredential); r.StatusCode != 200 {
		t.Errorf("minuto 16: %d %s", r.StatusCode, r.Body)
	}
}

func TestIntentos_B_VentanaDeslizanteNoFija(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	t0 := s.Clock.Now()
	for _, m := range []float64{0, 1, 2, 3, 14} {
		s.Clock.Set(t0.Add(time.Duration(m * float64(time.Minute))))
		if r := intento(t, s, "ximena", "mala"); r.StatusCode != 401 {
			t.Fatalf("fallo en el minuto %v: %d", m, r.StatusCode)
		}
	}
	s.Clock.Set(t0.Add(14*time.Minute + 30*time.Second))
	if r := intento(t, s, "ximena", testutil.FakeCredential); r.StatusCode != 429 {
		t.Errorf("minuto 14.5 (5 fallos en la ventana): %d", r.StatusCode)
	}
	s.Clock.Set(t0.Add(15*time.Minute + 30*time.Second))
	if r := intento(t, s, "ximena", testutil.FakeCredential); r.StatusCode != 200 {
		t.Errorf("minuto 15.5 (quedan 4 fallos: 1, 2, 3 y 14; un bloqueo fijo desde el último fallo seguiría hasta el 29): %d %s", r.StatusCode, r.Body)
	}
}

// Caso (c), C1-datos.md §2: se cuentan los fallos por usuario y por IP, lo que llegue primero a 5.
func TestIntentos_C_UsuarioEIP(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	h := s.Reg.Handler()
	try := func(user, pw, ip string) int {
		req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(fmt.Sprintf(`{"username":%q,"password":%q}`, user, pw)))
		req.Header.Set("Content-Type", "application/json")
		req.RemoteAddr = ip + ":40000"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req.WithContext(context.Background()))
		return rec.Code
	}
	const ip1, ip2 = "192.168.50.21", "192.168.50.22"
	for i := 0; i < 5; i++ {
		if c := try("axel", "mala", ip1); c != 401 {
			t.Fatalf("fallo %d: %d", i, c)
		}
	}
	if c := try("axel", testutil.FakeCredential, ip2); c != 429 {
		t.Errorf("mismo usuario desde otra IP (cuenta por usuario): %d", c)
	}
	if c := try("ximena", testutil.FakeCredential, ip1); c != 429 {
		t.Errorf("otro usuario desde la misma IP (cuenta por IP): %d", c)
	}
	if c := try("ximena", testutil.FakeCredential, ip2); c != 200 {
		t.Errorf("otro usuario desde otra IP: %d", c)
	}
}
