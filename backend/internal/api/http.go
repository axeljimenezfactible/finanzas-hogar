package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/dominio"
)

// Problem es un error HTTP con el esquema único de C1 §3 (para lo que no es de negocio: 400, 401, 413, 415, 429, 501…).
type Problem struct {
	Status  int
	Code    string
	Message string
	Fields  []FieldError
	RetryAt time.Time // solo 429
}

func (p *Problem) Error() string { return fmt.Sprintf("%d %s: %s", p.Status, p.Code, p.Message) }

// BadRequest es un 400 bad_request (formato incorrecto) con mensaje en español de México.
func BadRequest(msg string) *Problem {
	return &Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: msg}
}

var kindStatus = map[dominio.Kind]int{
	dominio.KindValidation: http.StatusUnprocessableEntity, dominio.KindNotFound: http.StatusNotFound,
	dominio.KindConflict: http.StatusConflict, dominio.KindForbidden: http.StatusForbidden,
	dominio.KindUnavailable: http.StatusServiceUnavailable, dominio.KindBadRequest: http.StatusBadRequest,
}

// Fail escribe err con el esquema único de error: *Problem tal cual; *dominio.Error según su clase; sql.ErrNoRows
// como 404; cualquier otro como 500 internal_error (el detalle va al registro, nunca a la respuesta).
func Fail(w http.ResponseWriter, r *http.Request, err error) {
	body := ErrorBody{RequestID: RequestID(r)}
	status := http.StatusInternalServerError
	var p *Problem
	if errors.As(err, &p) {
		status, body.Code, body.Message, body.Fields = p.Status, p.Code, p.Message, p.Fields
		if !p.RetryAt.IsZero() {
			body.RetryAt = dominio.FormatDateTime(p.RetryAt)
		}
	} else if e, ok := dominio.AsError(err); ok {
		status, body.Code, body.Message = kindStatus[e.Kind], e.Code, e.Message
		for _, f := range e.Fields {
			body.Fields = append(body.Fields, FieldError{Field: f.Field, Code: f.Code, Message: f.Message})
		}
	} else {
		body.Code, body.Message = "internal_error", "Ocurrió un error inesperado. Intenta de nuevo en un momento."
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteJSON responde v como JSON con el código dado.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// NoContent responde 204.
func NoContent(w http.ResponseWriter) error { w.WriteHeader(http.StatusNoContent); return nil }

// DecodeJSON lee el cuerpo JSON en dst con las reglas de C1 §3: JSON inválido, tipo incorrecto (973.75 o "973.75" en un
// *_cents), campo desconocido o datos de sobra → 400 bad_request; cuerpo mayor que Config.MaxJSONBytes → 413.
func (g *Registry) DecodeJSON(r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return BadRequest("El cuerpo debe ser JSON (Content-Type: application/json).")
	}
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, g.deps.Config.MaxJSONBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return jsonProblem(err)
	}
	if dec.More() {
		return BadRequest("El cuerpo trae datos de sobra después del JSON.")
	}
	return nil
}

func jsonProblem(err error) *Problem {
	var te *json.UnmarshalTypeError
	var mb *http.MaxBytesError
	switch {
	case errors.As(err, &mb):
		return &Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El cuerpo de la petición es demasiado grande."}
	case errors.As(err, &te):
		f := te.Field
		return &Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "El campo " + f + " tiene un tipo incorrecto.",
			Fields: []FieldError{{Field: f, Code: "invalid_format", Message: "Tipo incorrecto: los importes son enteros en centavos y las fechas, texto AAAA-MM-DD."}}}
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		f := strings.Trim(strings.TrimPrefix(err.Error(), "json: unknown field "), `"`)
		return &Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "El campo " + f + " no existe en esta operación.",
			Fields: []FieldError{{Field: f, Code: "not_allowed", Message: "Campo no admitido."}}}
	}
	return BadRequest("El cuerpo no es JSON válido.")
}

// Validator junta errores por campo y los devuelve como 422 validation_failed.
type Validator struct{ issues []dominio.FieldIssue }

// Add agrega un error de campo (code: required, invalid_format, out_of_range, not_found, not_allowed, conflict).
func (v *Validator) Add(field, code, msg string) {
	v.issues = append(v.issues, dominio.FieldIssue{Field: field, Code: code, Message: msg})
}

// Check agrega el error si ok es falso.
func (v *Validator) Check(ok bool, field, code, msg string) {
	if !ok {
		v.Add(field, code, msg)
	}
}

// Required exige un campo presente.
func (v *Validator) Required(present bool, field string) {
	v.Check(present, field, "required", "Este dato es obligatorio.")
}

// PositiveCents exige un importe mayor que cero.
func (v *Validator) PositiveCents(cents int64, field string) {
	v.Check(cents > 0, field, "out_of_range", "El monto debe ser mayor que cero.")
}

// Date valida una fecha AAAA-MM-DD; una fecha imposible es 400 según C1 §3, así que devuelve el *Problem para cortarse antes.
func (v *Validator) Date(s, field string) (dominio.Date, error) {
	d, err := dominio.ParseDate(s)
	if err != nil {
		return d, &Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "La fecha de " + field + " debe ser AAAA-MM-DD y existir.",
			Fields: []FieldError{{Field: field, Code: "invalid_format", Message: "Fecha inválida."}}}
	}
	return d, nil
}

// Err devuelve nil o el 422 con todos los campos.
func (v *Validator) Err() error {
	if len(v.issues) == 0 {
		return nil
	}
	return dominio.Invalid(v.issues...)
}

// PageParams son cursor, limit y sort de una colección (C1 §5). El cursor es opaco: lleva el desplazamiento y una
// huella de la consulta (orden y filtros), así que un cursor de otra consulta se rechaza con 400.
type PageParams struct {
	Limit  int
	Offset int
	Sort   []SortField
	print  string
}

// SortField es una propiedad del orden.
type SortField struct {
	Name string
	Desc bool
}

var reSort = regexp.MustCompile(`^-?[a-z][a-z0-9_]*(,-?[a-z][a-z0-9_]*)*$`)

// ParsePage lee limit (50 por defecto, de 1 a 200: fuera de rango → 400), sort (solo propiedades de allowed;
// sin valor = defaultSort) y cursor.
func ParsePage(r *http.Request, allowed []string, defaultSort string) (PageParams, error) {
	q := r.URL.Query()
	p := PageParams{Limit: 50}
	if s := q.Get("limit"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 200 {
			return p, BadRequest("El límite por página va de 1 a 200.")
		}
		p.Limit = n
	}
	sort := q.Get("sort")
	if sort == "" {
		sort = defaultSort
	}
	if sort != "" {
		if !reSort.MatchString(sort) {
			return p, BadRequest("El orden no tiene un formato válido.")
		}
		ok := map[string]bool{"id": true}
		for _, a := range allowed {
			ok[a] = true
		}
		for _, f := range strings.Split(sort, ",") {
			sf := SortField{Name: strings.TrimPrefix(f, "-"), Desc: strings.HasPrefix(f, "-")}
			if !ok[sf.Name] {
				return p, BadRequest("No se puede ordenar por " + sf.Name + ".")
			}
			p.Sort = append(p.Sort, sf)
		}
	}
	q.Del("cursor")
	q.Del("limit")
	q.Set("sort", sort)
	sum := sha256.Sum256([]byte(q.Encode()))
	p.print = hex.EncodeToString(sum[:8])
	if c := q.Get("cursor"); c == "" {
		if c = r.URL.Query().Get("cursor"); c != "" {
			var cur struct {
				O int    `json:"o"`
				F string `json:"f"`
			}
			b, err := base64.RawURLEncoding.DecodeString(c)
			if err != nil || len(c) > 512 || json.Unmarshal(b, &cur) != nil || cur.O < 0 || cur.F != p.print {
				return p, BadRequest("El cursor no es válido para esta consulta. Vuelve a la primera página.")
			}
			p.Offset = cur.O
		}
	}
	return p, nil
}

// OrderBy arma «ORDER BY …, id» con las columnas permitidas (propiedad → columna SQL); el id cierra el orden.
func (p PageParams) OrderBy(columns map[string]string) string {
	parts := []string{}
	hasID := false
	for _, f := range p.Sort {
		col := columns[f.Name]
		if col == "" {
			col = f.Name
		}
		if f.Name == "id" {
			hasID = true
		}
		dir := " ASC"
		if f.Desc {
			dir = " DESC"
		}
		parts = append(parts, col+dir)
	}
	if !hasID {
		idc := columns["id"]
		if idc == "" {
			idc = "id"
		}
		parts = append(parts, idc+" ASC")
	}
	return "ORDER BY " + strings.Join(parts, ", ")
}

// LimitOffset da «LIMIT limit+1 OFFSET offset»: se pide uno de más para saber si hay página siguiente.
func (p PageParams) LimitOffset() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit+1, p.Offset)
}

// PageOf recorta items a Limit y arma la envoltura Page de C1 con next_cursor (null si no hay más).
func PageOf[T any](items []T, p PageParams) struct {
	Items      []T              `json:"items"`
	NextCursor Nullable[string] `json:"next_cursor"`
} {
	out := struct {
		Items      []T              `json:"items"`
		NextCursor Nullable[string] `json:"next_cursor"`
	}{Items: items, NextCursor: Null[string]()}
	if out.Items == nil {
		out.Items = []T{}
	}
	if len(items) > p.Limit {
		out.Items = items[:p.Limit]
		b, _ := json.Marshal(map[string]any{"o": p.Offset + p.Limit, "f": p.print})
		out.NextCursor = Some(base64.RawURLEncoding.EncodeToString(b))
	}
	return out
}

// PathID lee un id de la ruta ({expense_id}…); no entero positivo → 400.
func PathID(r *http.Request, name string) (int64, error) {
	n, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	if err != nil || n < 1 {
		return 0, BadRequest("El identificador " + name + " no es válido.")
	}
	return n, nil
}

// QueryIDs lee una lista de ids separados por coma (filtros person_id, category_id, store_id, credit_id).
func QueryIDs(r *http.Request, name string) ([]int64, error) {
	s := r.URL.Query().Get(name)
	if s == "" {
		return nil, nil
	}
	var out []int64
	for _, x := range strings.Split(s, ",") {
		n, err := strconv.ParseInt(x, 10, 64)
		if err != nil || n < 1 {
			return nil, BadRequest("El filtro " + name + " debe ser una lista de identificadores separados por coma.")
		}
		out = append(out, n)
	}
	return out, nil
}

// QueryPeriod lee el filtro global period (nil si no viene).
func QueryPeriod(r *http.Request) (*dominio.Period, error) {
	s := r.URL.Query().Get("period")
	if s == "" {
		return nil, nil
	}
	p, err := dominio.ParsePeriod(s)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// QueryDate lee una fecha AAAA-MM-DD de la consulta (nil si no viene; formato incorrecto → 400).
func QueryDate(r *http.Request, name string) (*dominio.Date, error) {
	s := r.URL.Query().Get(name)
	if s == "" {
		return nil, nil
	}
	d, err := dominio.ParseDate(s)
	if err != nil {
		return nil, BadRequest("El parámetro " + name + " debe ser una fecha AAAA-MM-DD que exista.")
	}
	return &d, nil
}

// APITime convierte un instante de la base (UTC con Z) al formato DateTime de respuesta de C1 (-06:00).
func APITime(dbValue string) string { return dominio.DBTimeToAPI(dbValue) }

func fileSum(f *os.File) (string, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	_, err := f.Seek(0, io.SeekStart)
	return hex.EncodeToString(h.Sum(nil)), err
}

// staticHandler sirve frontend/dist (S6-12) y devuelve index.html en las rutas de la PWA que no son archivos.
func staticHandler(dir string) http.Handler {
	fsrv := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := filepath.Clean("/" + r.URL.Path)
		if st, err := os.Stat(filepath.Join(dir, clean)); err == nil && !st.IsDir() {
			fsrv.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(dir, "index.html"))
	})
}

// RegisterBase registra getHealth y getOpenApiDocument si ningún área lo hizo (son de auth-sistema; esto es el
// esqueleto de WF-0 para que el servidor responda mientras tanto).
func RegisterBase(g *Registry) {
	if !g.Has("getHealth") {
		g.Handle("getHealth", func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Cache-Control", "no-store")
			return WriteJSON(w, http.StatusOK, Health{Status: "ok", Time: dominio.FormatDateTime(g.deps.Clock.Now())})
		})
	}
	if !g.Has("getOpenApiDocument") {
		g.Handle("getOpenApiDocument", func(w http.ResponseWriter, r *http.Request) error {
			b, err := os.ReadFile(g.deps.Config.OpenAPIPath)
			if err != nil {
				return dominio.Unavailable("El contrato de la API no está disponible en este momento.")
			}
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
			_, err = w.Write(b)
			return err
		})
	}
}
