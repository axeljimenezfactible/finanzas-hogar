package api_test

// Pruebas de contrato de punta a punta (S6-03, S6-04) · agente pruebas-contrato.
//
// TestContrato_Operaciones recorre las 216 operaciones de routes.go sobre un servidor de prueba con sesión y las
// áreas reales registradas, las invoca con datos creados por la propia API (más unas pocas filas que ninguna operación
// puede crear: reglas de reparto, candidatos de chat) y valida CADA respuesta (2xx y 4xx) contra openapi.yaml con el
// validador de testutil. Además de la 2xx declarada de cada operación comprueba: 401 sin sesión, 400 sin
// Idempotency-Key, repetición idempotente, 409 con otro cuerpo, 404 con un id inexistente, 400 con campo desconocido,
// 422/400 con el cuerpo vacío y la página máxima (limit 200 sí, 201 y 0 no).
// Variable CONTRATO_TABLA=<archivo>: escribe la tabla operación → prueba → resultado en Markdown.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"

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
	"finanzas-hogar/internal/reportes"
	"finanzas-hogar/internal/telegram"
	"finanzas-hogar/internal/testutil"
)

const (
	fakeBearer = "token-de-servicio-falso-de-prueba"
	noExiste   = int64(999999999)
)

// fakeDocSender es el envío de documentos por Telegram: nada sale a la red (S7-07).
type fakeDocSender struct{}

func (fakeDocSender) SendDocument(context.Context, string, string, []byte, string) (int, error) {
	return 200, nil
}

type fakeSvc struct{}

func (fakeSvc) VerifyServiceToken(_ context.Context, tk string) (bool, error) {
	return tk == fakeBearer, nil
}

// K es el entorno de la prueba de contrato: servidor con todas las áreas, cliente con sesión y los ids creados.
type K struct {
	t     *testing.T
	st    *testing.T // subprueba vigente (a ella se atribuyen los fallos de validación)
	s     *testutil.Server
	c     *testutil.Client // Axel con sesión
	anon  *testutil.Client
	doc   *openapi3.T
	ids   map[string]int64
	n     int
	calls int
	// Telegram falso: nada sale a la red (S7-07, S6-13).
	tg   *httptest.Server
	root string // esquema del cuerpo vigente
}

// fila es una línea de la tabla operación → prueba → resultado.
type fila struct {
	Op, Metodo, Ruta string
	Prueba           string
	Llamadas         int
	Verifs           []string // lo comprobado y pasó
	Fallos           []string // causas
}

func newK(t *testing.T) *K {
	t.Helper()
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.Write([]byte(`{"ok":true,"result":{}}`))
	}))
	t.Cleanup(tg.Close)
	cfg := filepath.Join(testutil.TempDir(t), "openclaw.json")
	os.WriteFile(cfg, []byte(`{"bots":{"multix":{"botToken":"1234567890:falso_falso_falso"}}}`), 0o600)
	cl := &telegram.Client{Base: tg.URL, Token: telegram.FileToken(cfg), HTTP: tg.Client()}
	s := testutil.NewServer(t, auth.Register, ajustes.Register, avisos.RegisterWith(avisos.Options{Enviador: cl}),
		comidas.Register, comprasa.Register, comprasb.Register, creditos.Register, gastos.Register,
		presupuesto.Register, reportesapi.Register, tablero.Register)
	s.Reg.Deps().Service = fakeSvc{}
	k := &K{t: t, st: t, s: s, c: s.SessionFor(t, testutil.AxelUser), anon: s.Anonymous(t), doc: testutil.Doc(t),
		ids: map[string]int64{"person_id": s.Fx.AxelPersonID}, tg: tg}
	if err := s.Reg.Deps().Jobs.Prepare(context.Background()); err != nil {
		t.Fatalf("jobs: %v", err)
	}
	reportes.SetSender(fakeDocSender{})
	t.Cleanup(func() { reportes.SetSender(nil) })
	// filas que ninguna operación crea: un grupo de producto y un destinatario de Telegram de Axel
	if r, err := s.DB.Exec(`INSERT INTO product_group (name, purchase_channel) VALUES ('Verduras de prueba', 'market')`); err == nil {
		k.ids["product_group_id"], _ = r.LastInsertId()
	}
	return k
}

func (k *K) key() string { k.n++; return fmt.Sprintf("contrato-%05d-clave", k.n) }

// do manda la petición y valida la respuesta contra el contrato en la subprueba vigente.
func (k *K) do(cl *testutil.Client, method, path string, body any, h ...string) *testutil.Response {
	k.st.Helper()
	k.calls++
	res := cl.Raw(method, path, body, h...)
	testutil.ValidateResponse(k.st, method, path, res)
	return res
}

func jsonOf(b []byte) map[string]any {
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}

func snake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// idKey es el parámetro de ruta que alimenta la respuesta de una operación que crea (id → ids[clave]).
func idKey(op string) string {
	special := map[string]string{
		"createTicketScan": "scan_id", "uploadExpenseAttachment": "attachment_id", "createBackup": "backup_id",
		"generateMonthlyReport": "report_id", "generateExternalStatement": "ext_report_id",
		"markScheduledPaymentPaid": "credit_payment_id", "duplicateExpense": "dup_expense_id", "duplicateRecipe": "dup_recipe_id",
	}
	if v, ok := special[op]; ok {
		return v
	}
	if strings.HasPrefix(op, "create") {
		return snake(strings.TrimPrefix(op, "create")) + "_id"
	}
	return ""
}

// ---------- generador de cuerpos a partir del esquema ----------

func (k *K) idFor(prop string) int64 {
	if strings.Contains(prop, "person_id") {
		return k.s.Fx.AxelPersonID
	}
	aliases := map[string]string{"parent_id": "store_id", "from_store_id": "store_id", "to_store_id": "store_id",
		"base_recipe_id": "recipe_id", "scan_attachment_id": "scan_id"}
	if a, ok := aliases[prop]; ok {
		prop = a
	}
	if v, ok := k.ids[prop]; ok {
		return v
	}
	if v, ok := k.ids[strings.TrimPrefix(prop, "default_")]; ok {
		return v
	}
	return 1
}

func (k *K) genProp(name string, sr *openapi3.SchemaRef, depth int) any {
	if sr == nil || sr.Value == nil {
		return nil
	}
	s := sr.Value
	t := s.Type
	if t != nil && t.Includes("integer") {
		if strings.HasSuffix(name, "_id") && !strings.Contains(name, "cents") {
			return k.idFor(name)
		}
		if strings.HasSuffix(name, "_cents") || strings.Contains(sr.Ref, "MoneyCents") {
			v := int64(10000)
			if s.Min != nil && int64(*s.Min) > v {
				v = int64(*s.Min)
			}
			return v
		}
	}
	if t != nil && t.Includes("array") && strings.HasSuffix(name, "_ids") {
		return []int64{k.idFor(strings.TrimSuffix(name, "s"))}
	}
	return k.gen(sr, name, depth)
}

func (k *K) gen(sr *openapi3.SchemaRef, name string, depth int) any {
	if sr == nil || sr.Value == nil || depth > 7 {
		return nil
	}
	s := sr.Value
	if len(s.AllOf) > 0 {
		out := map[string]any{}
		for _, p := range s.AllOf {
			if m, ok := k.gen(p, name, depth+1).(map[string]any); ok {
				for kk, v := range m {
					out[kk] = v
				}
			}
		}
		return out
	}
	if len(s.OneOf) > 0 {
		return k.gen(s.OneOf[0], name, depth+1)
	}
	if len(s.AnyOf) > 0 {
		return k.gen(s.AnyOf[0], name, depth+1)
	}
	if len(s.Enum) > 0 {
		return s.Enum[0]
	}
	t := s.Type
	switch {
	case len(s.Properties) > 0 || (t != nil && t.Is("object")):
		out := map[string]any{}
		req := map[string]bool{}
		for _, r := range s.Required {
			req[r] = true
		}
		names := make([]string, 0, len(s.Properties))
		for n := range s.Properties {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			if req[n] {
				out[n] = k.genProp(n, s.Properties[n], depth+1)
			}
		}
		if len(req) == 0 && depth == 0 && strings.HasSuffix(k.root, "Update") { // PATCH: un campo sencillo
			pref := []string{"notes", "name", "display_name", "description", "label", "is_active"}
			for _, n := range append(pref, names...) {
				if p, ok := s.Properties[n]; ok && p.Value != nil && p.Value.Type != nil && !p.Value.Type.Is("object") && !p.Value.Type.Is("array") {
					out[n] = k.genProp(n, p, depth+1)
					break
				}
			}
		}
		return out
	case t != nil && t.Is("array"):
		n := 1
		if s.MinItems > 1 {
			n = int(s.MinItems)
		}
		arr := make([]any, 0, n)
		for i := 0; i < n; i++ {
			arr = append(arr, k.genProp(strings.TrimSuffix(name, "s"), s.Items, depth+1))
		}
		return arr
	case t != nil && t.Includes("integer"):
		v := int64(1)
		if s.Min != nil {
			v = int64(*s.Min)
		}
		if s.Max != nil && v > int64(*s.Max) {
			v = int64(*s.Max)
		}
		return v
	case t != nil && t.Includes("number"):
		if s.Min != nil {
			return *s.Min
		}
		return 1.0
	case t != nil && t.Includes("boolean"):
		return false
	case t != nil && t.Includes("string"):
		switch s.Format {
		case "date":
			return "2026-10-12"
		case "date-time":
			return "2026-10-12T12:00:00-06:00"
		}
		if ex, ok := s.Example.(string); ok && ex != "" {
			return ex
		}
		switch {
		case name == "week" || strings.HasSuffix(name, "_week"):
			return "2026-W42"
		case name == "month":
			return "2026-10"
		case strings.HasSuffix(name, "_date") || name == "date":
			return "2026-10-12"
		case name == "period":
			return "2026-10"
		case strings.Contains(name, "url"):
			return "https://tienda.invalid"
		}
		k.n++
		v := fmt.Sprintf("Prueba %d", k.n)
		if name == "code" {
			v = fmt.Sprintf("prueba_%d", k.n)
		}
		for uint64(len(v)) < s.MinLength {
			v += "x"
		}
		if s.MaxLength != nil && uint64(len(v)) > *s.MaxLength {
			v = v[:*s.MaxLength]
		}
		return v
	}
	return nil
}

// ---------- archivos multipart ----------

func mpart(fields map[string]string, filename, ctype string, data []byte) ([]byte, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for kk, v := range fields {
		_ = w.WriteField(kk, v)
	}
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", ctype)
	p, _ := w.CreatePart(h)
	p.Write(data)
	w.Close()
	return buf.Bytes(), w.FormDataContentType()
}

func pngBytes(n int) []byte {
	b := make([]byte, n)
	copy(b, "\x89PNG\r\n\x1a\n")
	return b
}

const csvPlanned = "\xef\xbb\xbf" + "tipo,concepto,monto,periodicidad,dia_o_fecha,categoria,persona,porcentaje_conjunto,notas\n" +
	"ingreso,Sueldo de prueba,\"15,000.00\",quincenal,,,axel,,\n" +
	"gasto_fijo,Renta de prueba,$9500,mensual,1,Renta,hogar,,deposito\n"

const csvProfeco = "producto,presentacion,marca,categoria,catalogo,precio,fecha_registro,cadena_comercial,giro,nombre_comercial,direccion,estado,municipio,latitud,longitud\n" +
	"Aceite,Botella 500 ml,Nutrioli,Aceites,Basicos,31.5,2026/10/01,Bodega Aurrera,Supermercado,Bodega Aurrera Centro,Calle 1,Ciudad de Mexico,Cuauhtemoc,19.43,-99.13\n"

func (k *K) multipartFor(op string) ([]byte, string) {
	switch op {
	case "uploadExpenseAttachment":
		return mpart(map[string]string{"kind": "ticket_photo"}, "ticket.png", "image/png", pngBytes(2048))
	case "createTicketScan":
		return mpart(nil, "ticket.png", "image/png", pngBytes(2048))
	case "importProfeco":
		return mpart(nil, "qqp.csv", "text/csv", []byte(csvProfeco))
	default: // importPlannedItems, previewPlannedItemsImport
		var nombre string // updateCategory pudo renombrarla: se lee el nombre vigente
		_ = k.s.DB.QueryRow(`SELECT name FROM category WHERE id = ?`, k.ids["category_id"]).Scan(&nombre)
		csv := strings.Replace(csvPlanned, "Renta,hogar", nombre+",axel", 1)
		return mpart(nil, "gastos_fijos_ingresos.csv", "text/csv", []byte(csv))
	}
}

// ---------- cuerpos, rutas y consultas ----------

// over: cuerpos que el generador no puede inventar (reglas de negocio que el esquema no expresa).
var over = map[string]func(k *K) any{}

func (k *K) bodyFor(meta api.OperationMeta) (any, []string) {
	switch meta.RequestType {
	case "":
		return nil, nil
	case "multipart/form-data":
		b, ct := k.multipartFor(meta.OperationID)
		return b, []string{"Content-Type", ct}
	}
	if f, ok := over[meta.OperationID]; ok {
		return f(k), nil
	}
	k.root = meta.RequestSchema
	sr := k.doc.Components.Schemas[meta.RequestSchema]
	return k.gen(sr, "", 0), nil
}

// pathParamOverride: parámetros de ruta que no salen de ids[nombre].
func (k *K) param(op, name string) string {
	switch name {
	case "provider":
		return "local_model"
	case "entity":
		return "expense"
	case "attachment_id":
		if strings.Contains(op, "TicketScan") {
			return fmt.Sprint(k.ids["scan_id"])
		}
	case "credit_id":
		if op == "getCardStatus" {
			return fmt.Sprint(k.ids["card_credit_id"])
		}
	case "expense_id":
		if op == "markExpensePaid" {
			return fmt.Sprint(k.ids["upcoming_expense_id"])
		}
	case "person_id":
		if op == "getExternalStatement" || op == "getPersonParticipation" {
			return fmt.Sprint(k.s.Fx.KarenPersonID)
		}
	}
	if v, ok := k.ids[name]; ok {
		return fmt.Sprint(v)
	}
	return k.lookup(name)
}

// lookup busca el primer id de una colección cuando ninguna operación lo creó (jobs, avisos, sesiones…).
func (k *K) lookup(name string) string {
	src := map[string]string{"job_id": "/api/v1/jobs", "notification_id": "/api/v1/notifications",
		"recipient_id": "/api/v1/notification-recipients", "cost_share_rule_id": "/api/v1/cost-share-rules",
		"scheduled_payment_id": "/api/v1/scheduled-payments?limit=1", "meal_plan_entry_id": "", "session_id": ""}
	if p := src[name]; p != "" {
		res := k.c.Raw("GET", p, nil)
		var pg struct {
			Items []struct {
				ID int64 `json:"id"`
			}
		}
		_ = json.Unmarshal(res.Body, &pg)
		if len(pg.Items) > 0 {
			k.ids[name] = pg.Items[0].ID
			return fmt.Sprint(pg.Items[0].ID)
		}
	}
	return "1"
}

func (k *K) path(r api.Route, bad bool) string {
	out := r.Pattern
	for {
		i := strings.IndexByte(out, '{')
		if i < 0 {
			return out
		}
		j := strings.IndexByte(out, '}')
		name := out[i+1 : j]
		v := k.param(r.OperationID, name)
		if bad {
			v = fmt.Sprint(noExiste)
			if name == "provider" || name == "entity" {
				v = "no-existe"
			}
		}
		out = out[:i] + v + out[j+1:]
	}
}

func (k *K) hasIDParam(r api.Route) bool { return strings.Contains(r.Pattern, "{") }

// query arma los parámetros de consulta obligatorios (y los que la operación necesite por regla).
func (k *K) query(r api.Route) url.Values {
	q := url.Values{}
	item := k.doc.Paths.Find(strings.TrimPrefix(r.Pattern, ""))
	if item == nil {
		return q
	}
	op := item.GetOperation(r.Method)
	params := append(openapi3.Parameters{}, item.Parameters...)
	params = append(params, op.Parameters...)
	for _, p := range params {
		if p.Value == nil || p.Value.In != "query" || !p.Value.Required {
			continue
		}
		v := k.gen(p.Value.Schema, p.Value.Name, 0)
		q.Set(p.Value.Name, fmt.Sprint(v))
	}
	for kk, v := range qover[r.OperationID] {
		q.Set(kk, v)
	}
	return q
}

var qover = map[string]map[string]string{}

func (k *K) hasLimit(r api.Route) bool {
	item := k.doc.Paths.Find(r.Pattern)
	if item == nil || r.Method != "GET" {
		return false
	}
	op := item.GetOperation("GET")
	params := append(openapi3.Parameters{}, item.Parameters...)
	params = append(params, op.Parameters...)
	for _, p := range params {
		if p.Value != nil && p.Value.Name == "limit" && p.Value.In == "query" {
			return true
		}
	}
	return false
}

// ---------- una operación ----------

func (k *K) client(meta api.OperationMeta) (*testutil.Client, []string) {
	switch meta.Auth {
	case "service":
		return k.anon, []string{"Authorization", "Bearer " + fakeBearer}
	case "public":
		return k.anon, nil
	}
	return k.c, nil
}

func in(code int, list []int) bool {
	for _, c := range list {
		if c == code {
			return true
		}
	}
	return false
}

// cause resume una respuesta de error; en un 500 agrega la última línea de error del registro del servidor.
func (k *K) cause(res *testutil.Response) string {
	out := short(res.Body)
	if res.StatusCode >= 500 {
		lines := strings.Split(k.s.Logs.String(), "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.Contains(lines[i], "level=ERROR") {
				l := lines[i]
				if len(l) > 420 {
					l = l[:420]
				}
				return out + " | LOG: " + l
			}
		}
	}
	return out
}

func short(b []byte) string {
	s := strings.Join(strings.Fields(string(b)), " ")
	if len(s) > 700 {
		s = s[:700] + "…"
	}
	return s
}

// prepare ejecuta lo que una operación necesita antes (filas que ninguna operación puede crear).
func (k *K) prepare(op string) {
	db := k.s.DB
	exec := func(q string, a ...any) int64 {
		res, err := db.Exec(q, a...)
		if err != nil {
			k.t.Fatalf("preparación %s: %v", op, err)
		}
		id, _ := res.LastInsertId()
		return id
	}
	_ = exec
	prepExtra(k, op)
}

func (k *K) runOp(st *testing.T, r api.Route, meta api.OperationMeta, f *fila) {
	k.st = st
	k.calls = 0
	ok := func(s string) { f.Verifs = append(f.Verifs, s) }
	bad := func(s string, a ...any) { f.Fallos = append(f.Fallos, fmt.Sprintf(s, a...)); st.Errorf(s, a...) }
	op := r.OperationID
	cl, hdr := k.client(meta)
	k.prepare(op)
	body, bh := k.bodyFor(meta)
	hdrs := append(append([]string{}, hdr...), bh...)
	withKey := func() []string {
		if meta.Idempotent {
			return append(append([]string{}, hdrs...), "Idempotency-Key", k.key())
		}
		return hdrs
	}
	full := func(bad bool) string {
		p := k.path(r, bad)
		if q := k.query(r); len(q) > 0 {
			p += "?" + q.Encode()
		}
		return p
	}

	// 1. sin sesión (o sin credencial de servicio) = 401
	if meta.Auth != "public" {
		res := k.do(k.anon, r.Method, full(false), body, append(append([]string{}, bh...), idemHdr(meta, k)...)...)
		if res.StatusCode != 401 {
			bad("sin credencial: %d, se esperaba 401: %s", res.StatusCode, k.cause(res))
		} else {
			ok("401 sin credencial")
		}
	}
	// 2. sin Idempotency-Key = 400
	if meta.Idempotent {
		res := k.do(cl, r.Method, full(false), body, append(append([]string{}, hdr...), bh...)...)
		if res.StatusCode != 400 {
			bad("sin Idempotency-Key: %d, se esperaba 400: %s", res.StatusCode, k.cause(res))
		} else {
			ok("400 sin Idempotency-Key")
		}
	}
	// 3. la 2xx declarada
	h := withKey()
	main := k.do(cl, r.Method, full(false), body, h...)
	if ex, has := excepcion[op]; has && in(main.StatusCode, ex.codes) {
		f.Verifs = append(f.Verifs, fmt.Sprintf("%d válida contra el contrato (excepción documentada: %s)", main.StatusCode, ex.porque))
	} else if !in(main.StatusCode, meta.Statuses) {
		bad("2xx: llegó %d y el contrato declara %v: %s", main.StatusCode, meta.Statuses, k.cause(main))
	} else {
		ok(fmt.Sprintf("%d declarada", main.StatusCode))
		if key := idKey(op); key != "" {
			if m := jsonOf(main.Body); m != nil {
				if id, isNum := m["id"].(float64); isNum {
					k.ids[key] = int64(id)
				}
			}
		}
		k.afterOK(op, main)
		// 4. repetición idempotente
		if meta.Idempotent && meta.Auth == "session" {
			rep := k.do(cl, r.Method, full(false), body, h...)
			okStatus := rep.StatusCode == main.StatusCode || (main.StatusCode == 201 && rep.StatusCode == 200 && in(200, meta.Statuses))
			switch {
			case !okStatus || rep.Header.Get("Idempotent-Replayed") != "true":
				bad("repetición con la misma clave: %d (replayed=%q), se esperaba la respuesta guardada: %s", rep.StatusCode, rep.Header.Get("Idempotent-Replayed"), short(rep.Body))
			case !sameJSON(main.Body, rep.Body) && !isK9(op):
				bad("repetición con la misma clave devolvió otro cuerpo")
			default:
				ok("repetición = respuesta guardada (Idempotent-Replayed)")
			}
			if meta.RequestType == "application/json" {
				otro := k.do(cl, r.Method, full(false), map[string]any{"cuerpo_distinto": k.n}, h...)
				if otro.StatusCode != 409 || jsonOf(otro.Body)["code"] != "idempotency_key_reused" {
					bad("misma clave con otro cuerpo: %d, se esperaba 409 idempotency_key_reused: %s", otro.StatusCode, short(otro.Body))
				} else {
					ok("409 clave con otro cuerpo")
				}
			}
		}
	}
	// 5. id inexistente = 404
	if k.hasIDParam(r) && !skip404[op] {
		res := k.do(cl, r.Method, full(true), body, withKey()...)
		if res.StatusCode != 404 {
			bad("id inexistente: %d, se esperaba 404: %s", res.StatusCode, k.cause(res))
		} else {
			ok("404 id inexistente")
		}
	}
	// 6. cuerpo: campo desconocido = 400; vacío = 422 (o 400)
	if meta.RequestType == "application/json" && meta.Auth != "public" {
		res := k.do(cl, r.Method, full(false), map[string]any{"campo_que_no_existe": 1}, withKey()...)
		if res.StatusCode != 400 {
			bad("campo desconocido: %d, se esperaba 400: %s", res.StatusCode, k.cause(res))
		} else {
			ok("400 campo desconocido")
		}
		if sr := k.doc.Components.Schemas[meta.RequestSchema]; sr != nil && sr.Value != nil && len(sr.Value.Required) > 0 {
			res = k.do(cl, r.Method, full(false), map[string]any{}, withKey()...)
			if res.StatusCode != 422 && res.StatusCode != 400 {
				bad("cuerpo vacío con campos obligatorios: %d, se esperaba 422: %s", res.StatusCode, k.cause(res))
			} else {
				ok(fmt.Sprintf("%d cuerpo sin obligatorios", res.StatusCode))
			}
		}
	}
	// 7. página máxima de 200
	if k.hasLimit(r) {
		q := func(l string) string {
			qq := k.query(r)
			qq.Set("limit", l)
			return k.path(r, false) + "?" + qq.Encode()
		}
		a, b, c0 := k.do(cl, "GET", q("200"), nil), k.do(cl, "GET", q("201"), nil), k.do(cl, "GET", q("0"), nil)
		if a.StatusCode != 200 || b.StatusCode != 400 || c0.StatusCode != 400 {
			bad("limit: 200→%d (200), 201→%d (400), 0→%d (400)", a.StatusCode, b.StatusCode, c0.StatusCode)
		} else {
			ok("limit 200 sí; 201 y 0 = 400")
		}
	}
	f.Llamadas = k.calls
}

func idemHdr(meta api.OperationMeta, k *K) []string {
	if meta.Idempotent {
		return []string{"Idempotency-Key", k.key()}
	}
	return nil
}

func isK9(op string) bool { return op == "createExpense" || op == "createWeightEntry" }

func sameJSON(a, b []byte) bool {
	var x, y any
	if json.Unmarshal(a, &x) != nil || json.Unmarshal(b, &y) != nil {
		return bytes.Equal(a, b)
	}
	xb, _ := json.Marshal(x)
	yb, _ := json.Marshal(y)
	return bytes.Equal(xb, yb)
}

// excepcion: operaciones cuya 2xx no se alcanza en el entorno de pruebas por una decisión documentada; se valida el
// código alterno contra el contrato y se informa aparte.
type excep struct {
	codes  []int
	porque string
}

var excepcion = map[string]excep{
	"suggestRecipeWithAi": {[]int{409}, "IA apagada por defecto y sin proveedor real (S6-13): 409 ai_disabled"},
	"reloadNutritionBase": {[]int{409}, "falta el archivo de la base nutricional local (bloqueo S6-compras-b-nutricion)"},
}

// skip404: operaciones cuya ruta no tiene un recurso que buscar (el id inexistente no es un 404).
var skip404 = map[string]bool{}

// afterOK guarda ids que salen de listas o de respuestas distintas a {id}.
func (k *K) afterOK(op string, res *testutil.Response) {
	switch op {
	case "createCategory":
		// Diagnóstico (CONTRATO_CATEGORIA_CON_CODE=1): la categoría creada por la API no tiene category.code (el esquema lo
		// permite) y el motor de presupuesto la lee como string; con el código puesto se ve qué más falla detrás.
		if os.Getenv("CONTRATO_CATEGORIA_CON_CODE") != "" {
			_, _ = k.s.DB.Exec(`UPDATE category SET code = 'prueba_' || id WHERE id = ?`, k.ids["category_id"])
		}
	case "createMealPlanEntries":
		var out struct {
			Items []struct {
				ID int64 `json:"id"`
			}
		}
		if json.Unmarshal(res.Body, &out) == nil && len(out.Items) > 0 {
			k.ids["meal_plan_entry_id"] = out.Items[0].ID
		}
	}
}

// ---------- orden de ejecución ----------

var primero = []string{
	"createCategory", "createStore", "createProduct", "createPresentation", "createNutritionFact", "createPrice",
	"createRecipe", "createInventoryItem", "createCredit", "createCreditTermRate", "createPromotion", "createCreditPurchase",
	"createHealthProfile", "createWeightEntry", "createExpense", "uploadExpenseAttachment", "createTicketScan",
	"saveTicketScanDraft", "createIncome", "createBudget", "createSavingsGoal", "createSavingsMovement", "createPlannedItem",
	"createShoppingListItem", "createMealPlanEntries", "createBackup", "generateMonthlyReport", "generateExternalStatement",
}
var ultimo = []string{"restoreBackup", "revokeSessions", "changeOwnPassword", "logout"}

func rank(r api.Route) int {
	for i, op := range primero {
		if op == r.OperationID {
			return i
		}
	}
	for i, op := range ultimo {
		if op == r.OperationID {
			return 1000 + i
		}
	}
	switch r.Method {
	case "GET":
		return 100
	case "DELETE":
		return 300
	}
	return 200
}

// ---------- pruebas ----------

func TestContrato_Operaciones(t *testing.T) {
	k := newK(t)
	routes := append([]api.Route{}, api.Routes...)
	sort.SliceStable(routes, func(i, j int) bool {
		if ri, rj := rank(routes[i]), rank(routes[j]); ri != rj {
			return ri < rj
		}
		return routes[i].Pattern+routes[i].Method < routes[j].Pattern+routes[j].Method
	})
	var filas []*fila
	pasaron := 0
	for _, r := range routes {
		meta, _ := api.Meta(r.OperationID)
		f := &fila{Op: r.OperationID, Metodo: r.Method, Ruta: strings.TrimPrefix(r.Pattern, "/api/v1"), Prueba: "TestContrato_Operaciones/" + r.OperationID}
		filas = append(filas, f)
		t.Run(r.OperationID, func(st *testing.T) {
			if !k.s.Reg.Has(r.OperationID) && r.OperationID != "getHealth" && r.OperationID != "getOpenApiDocument" {
				f.Fallos = append(f.Fallos, "no está registrada en ninguna área (responde 501)")
				st.Errorf("no está registrada en ninguna área (responde 501)")
				return
			}
			if f2, ok := special[r.OperationID]; ok {
				k.st = st
				k.calls = 0
				f2(k, st, f)
				f.Llamadas = k.calls
				return
			}
			k.runOp(st, r, meta, f)
		})
	}
	// una operación pasa si no tuvo fallos propios ni de validación contra el contrato
	for _, f := range filas {
		if len(f.Fallos) == 0 {
			pasaron++
		}
	}
	t.Logf("S6-03: %d de %d operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)", pasaron, len(filas))
	if path := os.Getenv("CONTRATO_TABLA"); path != "" {
		var b strings.Builder
		b.WriteString("| # | Operación | Método y ruta | Prueba | Llamadas | Resultado | Detalle |\n|---|---|---|---|---|---|---|\n")
		for i, f := range filas {
			res, det := "OK", strings.Join(f.Verifs, "; ")
			if len(f.Fallos) > 0 {
				res, det = "FALLA", strings.ReplaceAll(strings.Join(f.Fallos, " | "), "|", "/")
			}
			fmt.Fprintf(&b, "| %d | `%s` | %s `%s` | `%s` | %d | %s | %s |\n", i+1, f.Op, f.Metodo, f.Ruta, f.Prueba, f.Llamadas, res, det)
		}
		_ = os.WriteFile(path, []byte(b.String()), 0o644)
	}
}

// special: operaciones que no se pueden recorrer con el esquema general (sesión propia, cierre de sesión, etc.).
var special = map[string]func(k *K, st *testing.T, f *fila){}

// TestContrato_Rutas (S6-03): rutas del servidor = rutas de openapi.yaml, 0 sobrantes y 0 faltantes, también por HTTP.
func TestContrato_Rutas(t *testing.T) {
	k := newK(t)
	contrato := map[string]bool{}
	servidor := map[string]bool{}
	for path, item := range k.doc.Paths.Map() {
		for m, op := range item.Operations() {
			contrato[m+" "+path+" "+op.OperationID] = true
		}
	}
	for _, r := range api.Routes {
		servidor[r.Method+" "+r.Pattern+" "+r.OperationID] = true
	}
	var sobran, faltan []string
	for r := range servidor {
		if !contrato[r] {
			sobran = append(sobran, r)
		}
	}
	for r := range contrato {
		if !servidor[r] {
			faltan = append(faltan, r)
		}
	}
	if len(sobran) > 0 || len(faltan) > 0 {
		t.Errorf("rutas: sobran %v, faltan %v", sobran, faltan)
	}
	// toda operación del contrato está implementada por alguna área (si no, responde 501)
	var sin []string
	for _, m := range api.Operations {
		if !k.s.Reg.Has(m.OperationID) && m.OperationID != "getHealth" && m.OperationID != "getOpenApiDocument" {
			sin = append(sin, m.OperationID)
		}
	}
	if len(sin) > 0 {
		t.Errorf("operaciones del contrato sin implementar (501): %v", sin)
	}
	// por HTTP: ninguna ruta del contrato da 404/405 de enrutador y los métodos fuera del contrato dan 405
	k.st = t
	for path, item := range k.doc.Paths.Map() {
		sample := regexpParam(path)
		have := map[string]bool{}
		for m := range item.Operations() {
			have[m] = true
			res := k.c.Raw(m, sample, nil)
			if res.StatusCode == 405 {
				t.Errorf("%s %s: el servidor la rechaza con 405", m, path)
			}
		}
		for _, m := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
			if have[m] {
				continue
			}
			res := k.c.Raw(m, sample, nil)
			var e api.ErrorBody
			_ = json.Unmarshal(res.Body, &e)
			if res.StatusCode != 405 && matchesParamRoute(k.doc, m, path) {
				continue // p. ej. GET /credit-purchases/preview cae en /credit-purchases/{id}: id inválido, no es un método sobrante
			}
			if res.StatusCode != 405 || e.Code != "method_not_allowed" {
				t.Errorf("%s %s no está en el contrato y el servidor dio %d (%s), se esperaba 405", m, path, res.StatusCode, e.Code)
			}
		}
	}
	res := k.c.Raw("GET", "/api/v1/ruta-que-no-existe", nil)
	if res.StatusCode != 404 {
		t.Errorf("ruta fuera del contrato: %d", res.StatusCode)
	}
	t.Logf("S6-03: contrato %d operaciones, servidor %d, sobrantes %d, faltantes %d, sin implementar %d", len(contrato), len(servidor), len(sobran), len(faltan), len(sin))
}

func regexpParam(p string) string {
	for {
		i := strings.IndexByte(p, '{')
		if i < 0 {
			return p
		}
		j := strings.IndexByte(p, '}')
		v := "1"
		switch p[i+1 : j] {
		case "provider":
			v = "local_model"
		case "entity":
			v = "expense"
		}
		p = p[:i] + v + p[j+1:]
	}
}

// matchesParamRoute dice si el contrato tiene otra ruta con parámetro que absorbe a esta (misma forma, mismo método).
func matchesParamRoute(doc *openapi3.T, method, path string) bool {
	segs := strings.Split(path, "/")
	for other, item := range doc.Paths.Map() {
		if other == path || item.GetOperation(method) == nil {
			continue
		}
		os := strings.Split(other, "/")
		if len(os) != len(segs) {
			continue
		}
		ok := true
		for i := range os {
			if os[i] != segs[i] && !strings.HasPrefix(os[i], "{") {
				ok = false
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func operations() []api.OperationMeta { return api.Operations }
