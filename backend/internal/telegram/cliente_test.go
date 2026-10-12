package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"unicode/utf16"

	"finanzas-hogar/internal/testutil"
)

// fakeToken tiene la forma de un token de bot pero es falso; se arma en partes para que el grep de S7-06
// no encuentre ni siquiera el falso en el repositorio.
var fakeToken = strings.Repeat("1", 9) + ":" + strings.Repeat("prueba", 6)

type captura struct {
	mu     sync.Mutex
	bodies []map[string]any
	paths  []string
	status int
}

func fakeTelegram(t *testing.T, status int) (*httptest.Server, *captura) {
	t.Helper()
	c := &captura{status: status}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
		c.mu.Lock()
		c.bodies = append(c.bodies, m)
		c.paths = append(c.paths, r.URL.Path)
		c.mu.Unlock()
		if c.status != 200 {
			w.WriteHeader(c.status)
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request: chat not found"}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1}}`))
	}))
	t.Cleanup(ts.Close)
	return ts, c
}

func fakeConfig(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(testutil.TempDir(t), "openclaw.json")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestTextoPlano(t *testing.T) { // S7-05
	ts, c := fakeTelegram(t, 200)
	cl := &Client{Base: ts.URL, Token: FileToken(fakeConfig(t, `{"bots":{"multix":{"botToken":"`+fakeToken+`"}}}`)), HTTP: ts.Client()}
	text := Mensaje("reporte mensual", "Archivo: reporte_mensual_2026_10.pdf", "Total: $1,234.56")
	res, err := cl.SendMessage(context.Background(), "123456789", text)
	if err != nil || res.HTTPStatus != 200 || res.Parts != 1 {
		t.Fatalf("%+v %v", res, err)
	}
	b := c.bodies[0]
	if _, ok := b["parse_mode"]; ok {
		t.Fatal("el cuerpo no debe llevar parse_mode")
	}
	if b["chat_id"] != "123456789" || b["text"] != text {
		t.Fatalf("cuerpo: %v", b)
	}
	if !strings.Contains(b["text"].(string), "reporte_mensual_2026_10.pdf") || !strings.HasPrefix(b["text"].(string), "REPORTE MENSUAL\n") {
		t.Fatalf("texto alterado: %q", b["text"])
	}
	if !strings.HasSuffix(c.paths[0], "/sendMessage") {
		t.Fatal(c.paths[0])
	}
	t.Logf("S7-05 cuerpo sin parse_mode, con chat_id y text; guiones bajos íntegros; encabezado en mayúsculas")
}

func TestMensajeLargo(t *testing.T) { // S7-05
	ts, c := fakeTelegram(t, 200)
	cl := &Client{Base: ts.URL, Token: FileToken(fakeConfig(t, `{"bots":{"multix":{"botToken":"`+fakeToken+`"}}}`)), HTTP: ts.Client()}
	var sb strings.Builder
	for i := 0; sb.Len() < 3*MaxText; i++ {
		sb.WriteString("Línea con acentos y ñ, gasto número ")
		sb.WriteString(strings.Repeat("9", i%7+1))
		sb.WriteString(" 😀\n")
	}
	res, err := cl.SendMessage(context.Background(), "-100200", sb.String())
	if err != nil || res.Parts < 3 || len(c.bodies) != res.Parts {
		t.Fatalf("%+v %v", res, err)
	}
	var joined []string
	for _, b := range c.bodies {
		s := b["text"].(string)
		if n := len(utf16.Encode([]rune(s))); n > MaxText || n == 0 {
			t.Fatalf("trozo de %d unidades", n)
		}
		joined = append(joined, s)
	}
	if strings.Join(joined, "\n") != strings.TrimRight(sb.String(), "\n") {
		t.Fatal("al juntar los trozos no sale el texto original")
	}
	if got := Split(strings.Repeat("a", MaxText+10), MaxText); len(got) != 2 || len(got[0]) != MaxText {
		t.Fatalf("corte sin saltos de línea: %d", len(got))
	}
	t.Logf("S7-05 mensaje de %d bytes dividido en %d trozos ≤ %d unidades UTF-16", sb.Len(), res.Parts, MaxText)
}

func TestToken(t *testing.T) { // S7-06
	var logs bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logs, nil))
	ts, _ := fakeTelegram(t, 400)
	cfg := fakeConfig(t, `{"channels":{"telegram":{"bots":{"multix":{"botToken":"`+fakeToken+`"}}}}}`)
	cl := &Client{Base: ts.URL, Token: FileToken(cfg), HTTP: ts.Client()}
	_, err := cl.SendMessage(context.Background(), "1", "hola")
	if err == nil || strings.Contains(err.Error(), fakeToken) {
		t.Fatalf("error esperado sin token: %v", err)
	}
	log.Error("envío", "error", err.Error())
	// Error de transporte: la URL (con token) no debe salir en el error.
	dead := &Client{Base: "http://127.0.0.1:1", Token: FileToken(cfg), HTTP: ts.Client()}
	_, err = dead.SendMessage(context.Background(), "1", "hola")
	if err == nil || strings.Contains(err.Error(), fakeToken) || strings.Contains(err.Error(), "/bot") {
		t.Fatalf("error de transporte con URL: %v", err)
	}
	log.Error("envío", "error", err.Error())
	if strings.Contains(logs.String(), fakeToken) {
		t.Fatal("el token salió en el registro")
	}
	for _, tc := range []struct{ name, content, want string }{
		{"sin clave", `{"bots":{"otro":{"botToken":"x"}}, "secreto":"` + fakeToken + `"}`, "no trae la clave bots.multix.botToken"},
		{"json roto", `{"bots": ` + fakeToken, "no es JSON válido"},
	} {
		_, err := FileToken(fakeConfig(t, tc.content))()
		if err == nil || !strings.Contains(err.Error(), tc.want) || strings.Contains(err.Error(), fakeToken) {
			t.Fatalf("%s: %v", tc.name, err)
		}
	}
	if _, err := FileToken(filepath.Join(testutil.TempDir(t), "no-existe.json"))(); err == nil || !strings.Contains(err.Error(), "no existe") {
		t.Fatal(err)
	}
	t.Logf("S7-06 token leído en ejecución del archivo falso; 0 apariciones en errores y registro; errores claros sin contenido")
}
