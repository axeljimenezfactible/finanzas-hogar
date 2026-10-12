// Package telegram envía avisos por el bot MultiX (decisión F, K30): texto plano sin parse_mode,
// solo al host de la API de Telegram, con el token leído en ejecución y nunca impreso.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf16"
)

// APIBase es el único host al que habla el paquete (S7-14).
const APIBase = "https://api.telegram.org"

// MaxText es el límite de sendMessage para text: «1-4096 characters after entities parsing»
// (core.telegram.org/bots/api#sendmessage, consultado el 12-oct-2026). Se cuenta en unidades
// UTF-16, como Telegram, para que un emoji o un acento no hagan pasar un trozo del límite.
const MaxText = 4096

// MaxDocumentBytes es el límite de sendDocument con multipart (50 MB, core.telegram.org/bots/api#senddocument);
// lo usa el reporte mensual de s7-pdf (S7-10).
const MaxDocumentBytes = 50 * 1024 * 1024

// Result es el resultado de un envío (uno o varios trozos).
type Result struct {
	HTTPStatus int
	Parts      int
}

// APIError es una respuesta distinta de 200 de la API de Telegram (sin token).
type APIError struct {
	HTTPStatus  int
	Description string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Telegram respondió %d: %s", e.HTTPStatus, e.Description)
}

// Client manda mensajes. Base solo cambia en pruebas (servidor httptest).
type Client struct {
	Base  string
	Token TokenSource
	HTTP  *http.Client
}

// New arma el cliente real contra APIBase.
func New(tok TokenSource) *Client {
	return &Client{Base: APIBase, Token: tok, HTTP: &http.Client{Timeout: 20 * time.Second}}
}

// Mensaje arma un texto plano con el encabezado en mayúsculas y una línea por renglón.
func Mensaje(encabezado string, lineas ...string) string {
	return strings.Join(append([]string{strings.ToUpper(strings.TrimSpace(encabezado))}, lineas...), "\n")
}

// SendMessage envía text a chatID como texto plano (sin parse_mode); si pasa de MaxText lo divide.
func (c *Client) SendMessage(ctx context.Context, chatID, text string) (Result, error) {
	var res Result
	tok, err := c.Token()
	if err != nil {
		return res, err
	}
	for _, p := range Split(text, MaxText) {
		body, _ := json.Marshal(map[string]string{"chat_id": chatID, "text": p})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Base+"/bot"+tok+"/sendMessage", bytes.NewReader(body))
		if err != nil {
			return res, errors.New("telegram: no se pudo armar la petición")
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return res, errors.New(redact("telegram: no se pudo conectar con la API de Telegram: "+inner(err), tok))
		}
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		resp.Body.Close()
		res.HTTPStatus = resp.StatusCode
		if resp.StatusCode != http.StatusOK {
			var e struct {
				Description string `json:"description"`
			}
			_ = json.Unmarshal(raw, &e)
			if e.Description == "" {
				e.Description = http.StatusText(resp.StatusCode)
			}
			return res, &APIError{HTTPStatus: resp.StatusCode, Description: redact(e.Description, tok)}
		}
		res.Parts++
	}
	return res, nil
}

// inner quita la URL (que lleva el token) de un error de transporte.
func inner(err error) string {
	var ue *url.Error
	if errors.As(err, &ue) {
		err = ue.Err
	}
	return err.Error()
}

func redact(s, tok string) string {
	if tok == "" {
		return s
	}
	return strings.ReplaceAll(s, tok, "•••")
}

// Split divide s en trozos de a lo más max unidades UTF-16, cortando en el último salto de línea
// cuando lo hay; ningún trozo queda vacío.
func Split(s string, max int) []string {
	if len(utf16.Encode([]rune(s))) <= max {
		return []string{s}
	}
	var out []string
	rs := []rune(s)
	for len(rs) > 0 {
		n, units, lastNL := 0, 0, -1
		for n < len(rs) {
			u := utf16.RuneLen(rs[n])
			if u < 1 {
				u = 1
			}
			if units+u > max {
				break
			}
			units += u
			if rs[n] == '\n' {
				lastNL = n
			}
			n++
		}
		cut := n
		if n < len(rs) && lastNL > 0 {
			cut = lastNL + 1
		}
		if p := strings.TrimRight(string(rs[:cut]), "\n"); strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
		rs = rs[cut:]
	}
	return out
}
