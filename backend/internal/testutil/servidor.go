package testutil

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
)

// Server es un servidor de prueba TLS con la API completa (las áreas que se le pasen; el resto responde 501).
type Server struct {
	*httptest.Server
	DB    *sql.DB
	Fx    Fixtures
	Clock *Clock
	Reg   *api.Registry
	Logs  *bytes.Buffer // registro del servidor, para comprobar que no imprime contraseñas, cookies ni tokens
}

// NewServer arma el servidor con reloj simulado en 12-oct-2026 12:00 (America/Mexico_City).
func NewServer(t testing.TB, areas ...func(*api.Registry)) *Server {
	t.Helper()
	d, fx := NewDB(t)
	clk := NewClock(MX(2026, time.October, 12, 12, 0))
	logs := &bytes.Buffer{}
	log := slog.New(slog.NewTextHandler(&lockedWriter{w: logs}, nil))
	tmp := TempDir(t)
	deps := &api.Deps{DB: d, Clock: clk, Log: log, Jobs: jobs.New(d, clk, log), Config: api.Config{
		DataDir: filepath.Join(tmp, "datos"), BackupDir: filepath.Join(tmp, "respaldos"), TmpDir: tmp,
		OpenAPIPath: filepath.Join(Root(t), "docs", "contratos", "openapi.yaml"), SecureHeaders: true}}
	reg := api.NewRegistry(deps)
	for _, a := range areas {
		a(reg)
	}
	api.RegisterBase(reg)
	ts := httptest.NewTLSServer(reg.Handler())
	t.Cleanup(ts.Close)
	return &Server{Server: ts, DB: d, Fx: fx, Clock: clk, Reg: reg, Logs: logs}
}

type lockedWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (l *lockedWriter) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(b)
}

// Client hace peticiones con o sin sesión y valida cada respuesta contra el contrato.
type Client struct {
	t     testing.TB
	s     *Server
	hc    *http.Client
	Token string // token de sesión en claro ("" = sin sesión)
}

// Anonymous devuelve un cliente sin sesión.
func (s *Server) Anonymous(t testing.TB) *Client { return &Client{t: t, s: s, hc: s.Client()} }

// SessionFor crea una sesión de dispositivo para la cuenta (axel o ximena) directamente en la base, sin pasar por login.
func (s *Server) SessionFor(t testing.TB, username string) *Client {
	t.Helper()
	ctx := context.Background()
	u, err := dominio.UserByUsername(ctx, s.DB, username)
	if err != nil {
		t.Fatalf("testutil: cuenta %s: %v", username, err)
	}
	tok, hash, err := dominio.NewSessionToken()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.CreateSession(ctx, s.DB, u.ID, hash, "Prueba", "go-test", "127.0.0.1", s.Clock.Now(), 90*24*time.Hour); err != nil {
		t.Fatal(err)
	}
	return &Client{t: t, s: s, hc: s.Client(), Token: tok}
}

// Response es la respuesta con el cuerpo ya leído.
type Response struct {
	*http.Response
	Body []byte
}

// JSON decodifica el cuerpo.
func (r *Response) JSON(t testing.TB, v any) {
	t.Helper()
	if err := json.Unmarshal(r.Body, v); err != nil {
		t.Fatalf("respuesta no es JSON: %v: %s", err, r.Body)
	}
}

// Do manda la petición (body: nil, []byte, string o cualquier valor que se codifica como JSON), valida la respuesta
// contra openapi.yaml (código declarado, esquema y tipo de contenido) y la devuelve. headers: pares clave, valor.
func (c *Client) Do(method, path string, body any, headers ...string) *Response {
	c.t.Helper()
	resp := c.Raw(method, path, body, headers...)
	ValidateResponse(c.t, method, path, resp)
	return resp
}

// Raw manda la petición sin validar la respuesta.
func (c *Client) Raw(method, path string, body any, headers ...string) *Response {
	c.t.Helper()
	var rd io.Reader
	ctype := ""
	switch b := body.(type) {
	case nil:
	case []byte:
		rd = bytes.NewReader(b)
	case string:
		rd, ctype = strings.NewReader(b), "application/json"
	default:
		js, err := json.Marshal(b)
		if err != nil {
			c.t.Fatal(err)
		}
		rd, ctype = bytes.NewReader(js), "application/json"
	}
	req, err := http.NewRequest(method, c.s.URL+path, rd)
	if err != nil {
		c.t.Fatal(err)
	}
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	if c.Token != "" {
		req.AddCookie(&http.Cookie{Name: api.SessionCookie, Value: c.Token})
	}
	res, err := c.hc.Do(req)
	if err != nil {
		c.t.Fatal(err)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		c.t.Fatal(err)
	}
	return &Response{Response: res, Body: b}
}

var (
	specOnce   sync.Once
	specRouter routers.Router
	specDoc    *openapi3.T
	specErr    error
)

// Spec carga docs/contratos/openapi.yaml una vez y arma el enrutador de validación.
func Spec(t testing.TB) routers.Router {
	t.Helper()
	specOnce.Do(func() {
		loader := openapi3.NewLoader()
		doc, err := loader.LoadFromFile(filepath.Join(Root(t), "docs", "contratos", "openapi.yaml"))
		if err != nil {
			specErr = err
			return
		}
		doc.Servers = openapi3.Servers{{URL: "/"}}
		specDoc = doc
		for _, p := range doc.Paths.Map() {
			for _, op := range p.Operations() {
				for _, r := range op.Responses.Map() {
					if r.Value == nil {
						continue
					}
					for ct, mt := range r.Value.Content {
						text := mt.Schema != nil && mt.Schema.Value != nil && mt.Schema.Value.Type.Is("string")
						if openapi3filter.RegisteredBodyDecoder(ct) == nil || (text && ct != "application/json") {
							openapi3filter.RegisterBodyDecoder(ct, openapi3filter.FileBodyDecoder) // binario o texto (p. ej. application/yaml como string)
						}
					}
				}
			}
		}
		specRouter, specErr = gorillamux.NewRouter(doc)
	})
	if specErr != nil {
		t.Fatalf("testutil: openapi.yaml: %v", specErr)
	}
	return specRouter
}

// Doc devuelve el contrato cargado (para comparar rutas, S6-03).
func Doc(t testing.TB) *openapi3.T {
	t.Helper()
	Spec(t)
	return specDoc
}

// ValidateResponse comprueba que el código esté declarado (o caiga en default) y que el cuerpo cumpla el esquema.
func ValidateResponse(t testing.TB, method, path string, resp *Response) {
	t.Helper()
	req, _ := http.NewRequest(method, "http://finanzas.local"+path, nil)
	route, params, err := Spec(t).FindRoute(req)
	if err != nil {
		t.Fatalf("%s %s no está en el contrato: %v", method, path, err)
	}
	in := &openapi3filter.RequestValidationInput{Request: req, PathParams: params, Route: route,
		Options: &openapi3filter.Options{AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, IncludeResponseStatus: true}}
	out := &openapi3filter.ResponseValidationInput{RequestValidationInput: in, Status: resp.StatusCode, Header: resp.Header,
		Body: io.NopCloser(bytes.NewReader(resp.Body)), Options: in.Options}
	if err := openapi3filter.ValidateResponse(context.Background(), out); err != nil {
		t.Errorf("%s %s → %d no cumple el contrato: %v\ncuerpo: %.400s", method, path, resp.StatusCode, err, resp.Body)
	}
}
