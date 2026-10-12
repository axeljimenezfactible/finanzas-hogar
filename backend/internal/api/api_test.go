package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/ia"
	"finanzas-hogar/internal/testutil"
)

var reParam = regexp.MustCompile(`\{[a-z_]+\}`)

func samplePath(p string) string {
	return reParam.ReplaceAllStringFunc(p, func(s string) string {
		if s == "{provider}" {
			return "local_model"
		}
		if s == "{entity}" {
			return "expense"
		}
		return "1"
	})
}

// S6-03: las rutas registradas son idénticas a paths de openapi.yaml (0 sobrantes, 0 faltantes) y Operations las describe todas.
func TestRutasIgualesAlContrato(t *testing.T) {
	doc := testutil.Doc(t)
	var contrato, servidor []string
	for path, item := range doc.Paths.Map() {
		for m, op := range item.Operations() {
			contrato = append(contrato, m+" "+path+" "+op.OperationID)
		}
	}
	for _, r := range api.Routes {
		servidor = append(servidor, r.Method+" "+r.Pattern+" "+r.OperationID)
		if _, ok := api.Meta(r.OperationID); !ok {
			t.Errorf("%s sin metadatos en types_gen.go", r.OperationID)
		}
	}
	sort.Strings(contrato)
	sort.Strings(servidor)
	if strings.Join(contrato, "\n") != strings.Join(servidor, "\n") || len(api.Operations) != len(api.Routes) {
		t.Fatalf("rutas distintas: contrato %d, servidor %d, metadatos %d", len(contrato), len(servidor), len(api.Operations))
	}
	t.Logf("S6-03: %d operaciones en %d rutas, 0 diferencias", len(servidor), len(doc.Paths.Map()))
}

// Cada operación del contrato responde lo que el contrato declara sin sesión (401, o la pública) y con sesión
// (501 mientras su área no la implemente; 400 si falta Idempotency-Key). Valida el 100 % de respuestas.
func TestContratoSinSesionYEsqueleto(t *testing.T) {
	s := testutil.NewServer(t)
	anon, axel := s.Anonymous(t), s.SessionFor(t, testutil.AxelUser)
	for _, r := range api.Routes {
		meta, _ := api.Meta(r.OperationID)
		path := samplePath(r.Pattern)
		res := anon.Do(r.Method, path, nil)
		switch {
		case meta.Auth == "public" && res.StatusCode == http.StatusUnauthorized:
			t.Errorf("%s es pública y dio 401", r.OperationID)
		case meta.Auth != "public" && res.StatusCode != http.StatusUnauthorized:
			t.Errorf("%s sin sesión dio %d", r.OperationID, res.StatusCode)
		}
		if meta.Auth != "session" || s.Reg.Has(r.OperationID) {
			continue
		}
		res = axel.Do(r.Method, path, nil)
		want := http.StatusNotImplemented
		if meta.Idempotent {
			want = http.StatusBadRequest
		}
		if res.StatusCode != want {
			t.Errorf("%s con sesión: %d, se esperaba %d", r.OperationID, res.StatusCode, want)
		}
	}
	if res := axel.Do("GET", "/health", nil); res.StatusCode != 200 || res.Header.Get("X-Request-Id") == "" {
		t.Error("getHealth")
	}
	if res := axel.Do("GET", "/openapi.yaml", nil); res.StatusCode != 200 || !strings.Contains(string(res.Body), "openapi:") {
		t.Error("getOpenApiDocument")
	}
	for _, c := range []struct{ m, p, code string }{{"GET", "/api/v1/no-existe", "not_found"}, {"PUT", "/api/v1/users", "method_not_allowed"}} {
		res := axel.Raw(c.m, c.p, nil)
		var e api.ErrorBody
		res.JSON(t, &e)
		if e.Code != c.code || e.RequestID == "" || e.Message == "" {
			t.Errorf("%s %s: %+v", c.m, c.p, e)
		}
	}
}

// S6-05 y S6-07: cookie HttpOnly; Secure; SameSite=Strict; 90 días renovables con el uso; sesión cerrada = 401.
func TestSesionCookieYRenovacion(t *testing.T) {
	s := testutil.NewServer(t)
	axel, xim := s.SessionFor(t, testutil.AxelUser), s.SessionFor(t, testutil.XimenaUser)
	t0 := s.Clock.Now()
	s.Clock.Advance(60 * 24 * time.Hour)
	res := axel.Do("GET", "/openapi.yaml", nil)
	sc := res.Header.Get("Set-Cookie")
	for _, want := range []string{"fh_session=", "Max-Age=7776000", "HttpOnly", "Secure", "SameSite=Strict", "Path=/"} {
		if !strings.Contains(sc, want) {
			t.Errorf("Set-Cookie sin %s: %s", want, sc)
		}
	}
	var exp string
	_ = s.DB.QueryRow(`SELECT expires_at FROM device_session WHERE token_hash = ?`, dominio.HashToken(axel.Token)).Scan(&exp)
	if got, _ := dominio.ParseDBTime(exp); !got.Equal(t0.Add(150 * 24 * time.Hour)) {
		t.Errorf("el uso del día 60 debe mover el vencimiento al día 150: %s", exp)
	}
	s.Clock.Advance(31 * 24 * time.Hour) // día 91: Ximena no usó su sesión
	if res := xim.Do("GET", "/openapi.yaml", nil); res.StatusCode != http.StatusUnauthorized {
		t.Errorf("día 91 sin uso: %d", res.StatusCode)
	}
	if res := axel.Do("GET", "/openapi.yaml", nil); res.StatusCode != 200 {
		t.Error("Axel sigue vigente")
	}
	var id int64
	_ = s.DB.QueryRow(`SELECT id FROM device_session WHERE token_hash = ?`, dominio.HashToken(axel.Token)).Scan(&id)
	_ = dominio.RevokeSession(context.Background(), s.DB, id, s.Clock.Now())
	if res := axel.Do("GET", "/openapi.yaml", nil); res.StatusCode != http.StatusUnauthorized {
		t.Error("sesión cerrada debe dar 401")
	}
	rec := httptest.NewRecorder()
	api.ClearSessionCookie(rec)
	if !strings.Contains(rec.Header().Get("Set-Cookie"), "Max-Age=0") {
		t.Error("ClearSessionCookie")
	}
}

// S6-04 (K6, K9): misma clave y mismo cuerpo = 1 registro y repetición con Idempotent-Replayed; otro cuerpo = 409;
// pasadas 24 h la tabla ya no la tiene; una respuesta de error libera la clave.
func TestIdempotencia(t *testing.T) {
	var calls atomic.Int32
	area := func(reg *api.Registry) {
		reg.Handle("createExpense", func(w http.ResponseWriter, r *http.Request) error {
			var in struct {
				AmountCents int64 `json:"amount_cents"`
			}
			if err := reg.DecodeJSON(r, &in); err != nil {
				return err
			}
			if in.AmountCents <= 0 {
				var v api.Validator
				v.PositiveCents(in.AmountCents, "amount_cents")
				return v.Err()
			}
			n := calls.Add(1)
			return api.WriteJSON(w, http.StatusCreated, map[string]any{"id": n, "key": api.IdempotencyKey(r)})
		})
	}
	s := testutil.NewServer(t, area)
	c := s.SessionFor(t, testutil.AxelUser)
	key := "3f6c2a8e-7b1d-4c9a-9e2f-5a7d1b0c4e81"
	r1 := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 19900}`, "Idempotency-Key", key)
	r2 := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 19900}`, "Idempotency-Key", key)
	if r1.StatusCode != 201 || r2.StatusCode != 200 || r2.Header.Get("Idempotent-Replayed") != "true" || string(r1.Body) != string(r2.Body) || calls.Load() != 1 {
		t.Fatalf("repetición: %d %d %s %s (%d llamadas)", r1.StatusCode, r2.StatusCode, r1.Body, r2.Body, calls.Load())
	}
	if r := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 20000}`, "Idempotency-Key", key); r.StatusCode != 409 || !strings.Contains(string(r.Body), "idempotency_key_reused") {
		t.Errorf("otro cuerpo: %d %s", r.StatusCode, r.Body)
	}
	if r := s.SessionFor(t, testutil.XimenaUser).Raw("POST", "/api/v1/expenses", `{"amount_cents": 19900}`, "Idempotency-Key", key); r.StatusCode != 201 {
		t.Errorf("la clave es por cuenta: %d", r.StatusCode)
	}
	r := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 0}`, "Idempotency-Key", "clave-de-error-0000001")
	var e api.ErrorBody
	r.JSON(t, &e)
	if r.StatusCode != 422 || e.Code != "validation_failed" || len(e.Fields) != 1 || e.Fields[0].Field != "amount_cents" {
		t.Errorf("422: %d %+v", r.StatusCode, e)
	}
	if r := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 100}`, "Idempotency-Key", "clave-de-error-0000001"); r.StatusCode != 201 {
		t.Errorf("tras un error la clave se libera: %d", r.StatusCode)
	}
	s.Clock.Advance(25 * time.Hour)
	if r := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 19900}`, "Idempotency-Key", key); r.StatusCode != 201 {
		t.Errorf("pasadas 24 h la tabla ya no tiene la clave (gastos y pesos la buscan en su renglón, K9): %d", r.StatusCode)
	}
	for _, k := range []string{"", "corta", "con espacios no válida"} {
		if r := c.Raw("POST", "/api/v1/expenses", `{"amount_cents": 1}`, "Idempotency-Key", k); r.StatusCode != 400 {
			t.Errorf("clave %q: %d", k, r.StatusCode)
		}
	}
	for body, code := range map[string]int{`{"amount_cents": 973.75}`: 400, `{"amount_cents": "973.75"}`: 400, `{"otro": 1}`: 400, `{`: 400, `{"amount_cents": 1} {}`: 400} {
		if r := c.Raw("POST", "/api/v1/expenses", body, "Idempotency-Key", "k"+strings.Repeat("x", 20)+string(rune('a'+len(body)%26))); r.StatusCode != code {
			t.Errorf("cuerpo %s: %d", body, r.StatusCode)
		}
	}
	if r := c.Raw("POST", "/api/v1/expenses", []byte(`{"amount_cents": 1}`), "Idempotency-Key", "sin-content-type-0001"); r.StatusCode != 400 {
		t.Errorf("sin Content-Type JSON: %d", r.StatusCode)
	}
}

// S6-04: limit > 200 se rechaza (400), 50 por defecto, cursor opaco ligado a la consulta.
func TestPaginacion(t *testing.T) {
	req := func(q string) *http.Request { return httptest.NewRequest("GET", "/api/v1/expenses?"+q, nil) }
	if _, err := api.ParsePage(req("limit=201"), nil, ""); err == nil {
		t.Error("limit=201 debe ser 400")
	}
	for _, q := range []string{"limit=0", "limit=x", "sort=Amount", "sort=amount_cents", "cursor=zzz"} {
		if _, err := api.ParsePage(req(q), []string{"expense_date"}, "-expense_date"); err == nil {
			t.Errorf("%s debió fallar", q)
		}
	}
	p, err := api.ParsePage(req("person_id=1"), []string{"expense_date"}, "-expense_date")
	if err != nil || p.Limit != 50 || p.OrderBy(map[string]string{"expense_date": "e.expense_date"}) != "ORDER BY e.expense_date DESC, id ASC" || p.LimitOffset() != "LIMIT 51 OFFSET 0" {
		t.Fatalf("por defecto: %+v %v %s", p, err, p.OrderBy(nil))
	}
	p.Limit = 2
	pg := api.PageOf([]int{1, 2, 3}, p)
	if len(pg.Items) != 2 || !pg.NextCursor.Valid {
		t.Fatal("PageOf")
	}
	p2, err := api.ParsePage(req("person_id=1&limit=2&cursor="+pg.NextCursor.Value), []string{"expense_date"}, "-expense_date")
	if err != nil || p2.Offset != 2 {
		t.Fatalf("cursor: %+v %v", p2, err)
	}
	if _, err := api.ParsePage(req("person_id=2&cursor="+pg.NextCursor.Value), []string{"expense_date"}, "-expense_date"); err == nil {
		t.Error("cursor de otra consulta")
	}
	if last := api.PageOf([]int(nil), p2); last.NextCursor.Valid || last.Items == nil {
		t.Error("última página: next_cursor null e items []")
	}
	b, _ := json.Marshal(api.PageOf([]int{}, p))
	if string(b) != `{"items":[],"next_cursor":null}` {
		t.Errorf("Page: %s", b)
	}
	if ids, err := api.QueryIDs(req("person_id=1,2"), "person_id"); err != nil || len(ids) != 2 {
		t.Error("QueryIDs")
	}
	if _, err := api.QueryIDs(req("person_id=1,x"), "person_id"); err == nil {
		t.Error("QueryIDs malo")
	}
	if pp, err := api.QueryPeriod(req("period=2026-10-q2")); err != nil || pp.Start.Day != 16 {
		t.Error("QueryPeriod")
	}
	if _, err := api.QueryPeriod(req("period=oct-2026")); err == nil {
		t.Error("period malo")
	}
	if dd, err := api.QueryDate(req("from=2026-10-10"), "from"); err != nil || dd.Day != 10 {
		t.Error("QueryDate")
	}
	if _, err := api.QueryDate(req("from=10-oct-2026"), "from"); err == nil {
		t.Error("QueryDate malo")
	}
	if api.APITime("2026-10-11T03:30:00.000Z") != "2026-10-10T21:30:00-06:00" {
		t.Error("APITime")
	}
}

type failTransport struct{ t *testing.T }

func (f failTransport) RoundTrip(*http.Request) (*http.Response, error) {
	f.t.Error("con la IA apagada no debe salir ninguna llamada")
	return nil, errors.New("bloqueado")
}

type netProvider struct {
	hc *http.Client
	ia.Off
}

func (p netProvider) SuggestRecipe(ctx context.Context, _ ia.RecipeRequest) (ia.RecipeSuggestion, error) {
	_, err := p.hc.Get("https://api.example.invalid/")
	return ia.RecipeSuggestion{}, err
}

// S6-13: IA apagada por defecto; con un adaptador que saldría a la red, Active no lo toca.
func TestIAApagadaSinLlamadas(t *testing.T) {
	s := testutil.NewServer(t)
	reg := s.Reg.Deps().IA
	reg.Set(netProvider{hc: &http.Client{Transport: failTransport{t}}, Off: ia.Off{ProviderCode: ia.ClaudeAPI}})
	ctx := context.Background()
	if _, err := reg.Active(ctx); !errors.Is(err, ia.ErrDisabled) {
		t.Fatalf("IA apagada: %v", err)
	}
	for _, c := range []string{ia.LocalModel, ia.ClaudeAPI} {
		if on, err := reg.Enabled(ctx, c); on || err != nil {
			t.Errorf("%s encendida por defecto", c)
		}
	}
	if _, err := reg.Enabled(ctx, "otro"); err == nil || reg.Available(ia.LocalModel) || !reg.Available(ia.ClaudeAPI) {
		t.Error("Enabled/Available")
	}
	off := ia.Off{ProviderCode: ia.LocalModel}
	if _, err := off.ReadTicket(ctx, ia.TicketImage{}); !errors.Is(err, ia.ErrDisabled) || off.Code() != ia.LocalModel {
		t.Error("Off.ReadTicket")
	}
	if _, err := off.SuggestRecipe(ctx, ia.RecipeRequest{}); !errors.Is(err, ia.ErrDisabled) {
		t.Error("Off.SuggestRecipe")
	}
	_, _ = s.DB.Exec(`UPDATE app_setting SET value_json = 'true' WHERE setting_key = 'ai.local_model_enabled'`)
	if _, err := reg.Active(ctx); !errors.Is(err, ia.ErrDisabled) {
		t.Error("encendida sin adaptador real sigue apagada")
	}
}

// S6-15: el registro no imprime el token de sesión, la cookie ni el cuerpo; errores con mensaje en español y code en inglés.
func TestLogSinSecretosYErrores(t *testing.T) {
	s := testutil.NewServer(t, func(reg *api.Registry) {
		reg.Handle("listUsers", func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("detalle interno con " + testutil.FakeCredential)
		})
		reg.Handle("listSessions", func(w http.ResponseWriter, r *http.Request) error {
			return &api.Problem{Status: 429, Code: "too_many_attempts", Message: "Demasiados intentos.", RetryAt: time.Date(2026, 10, 11, 3, 45, 0, 0, time.UTC)}
		})
		reg.Handle("listJobs", func(w http.ResponseWriter, r *http.Request) error { panic("x") })
		reg.Handle("getSystemStatus", func(w http.ResponseWriter, r *http.Request) error { return dominio.Conflict("in_use", "En uso.") })
		reg.Handle("getSeedStatus", func(w http.ResponseWriter, r *http.Request) error { return dominio.ErrNotFound })
	})
	c := s.SessionFor(t, testutil.AxelUser)
	for path, want := range map[string]string{"/api/v1/users": "internal_error", "/api/v1/sessions": "too_many_attempts", "/api/v1/jobs": "internal_error", "/api/v1/system/status": "in_use", "/api/v1/system/seed-status": "not_found"} {
		var e api.ErrorBody
		c.Raw("GET", path, nil).JSON(t, &e)
		if e.Code != want || strings.Contains(e.Message, testutil.FakeCredential) {
			t.Errorf("%s: %+v", path, e)
		}
		if want == "too_many_attempts" && e.RetryAt != "2026-10-10T21:45:00-06:00" {
			t.Errorf("retry_at: %s", e.RetryAt)
		}
	}
	logs := s.Logs.String()
	if strings.Contains(logs, c.Token) || !strings.Contains(logs, "status=500") {
		t.Error("el registro no debe tener el token de sesión")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("Handle de una operación que no está en C1 debe fallar")
			}
		}()
		s.Reg.Handle("noExiste", nil)
	}()
	func() {
		defer func() {
			if recover() == nil {
				t.Error("Handle duplicado debe fallar")
			}
		}()
		s.Reg.Handle("listUsers", nil)
	}()
	if len(s.Reg.Implemented()) != 7 {
		t.Errorf("Implemented: %v", s.Reg.Implemented())
	}
}

// S6-12 (parte del servidor): la PWA se sirve desde FrontendDir; las rutas de la API siguen con su 404 JSON.
func TestFrontendYCabeceras(t *testing.T) {
	dir := testutil.TempDir(t)
	_ = os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html><title>Finanzas</title>"), 0o600)
	_ = os.WriteFile(filepath.Join(dir, "manifest.webmanifest"), []byte(`{"name":"Finanzas"}`), 0o600)
	db, _ := testutil.NewDB(t)
	reg := api.NewRegistry(&api.Deps{DB: db, Config: api.Config{FrontendDir: dir, SecureHeaders: true, OpenAPIPath: "/no/existe"}})
	api.RegisterBase(reg)
	h := reg.Handler()
	for path, want := range map[string]string{"/gastos/12": "<!doctype", "/manifest.webmanifest": `{"name"`, "/api/v1/nada": `"not_found"`} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if !strings.Contains(rec.Body.String(), want) || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("%s: %d %s", path, rec.Code, rec.Body.String())
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/health", nil))
	if rec.Code != 200 || rec.Header().Get("Cache-Control") != "no-store" {
		t.Error("health")
	}
	if api.ClientIP(&http.Request{RemoteAddr: "sin-puerto"}) != "sin-puerto" {
		t.Error("ClientIP")
	}
}
