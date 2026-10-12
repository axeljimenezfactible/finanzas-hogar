package ajustes

// Importación de gastos_fijos_ingresos.csv (plan §8.2; N-304). Columnas de la plantilla: tipo, concepto, monto,
// periodicidad, dia_o_fecha, categoria, persona, porcentaje_conjunto (opcional), notas. La correspondencia de cada
// columna con planned_item no está escrita en C1 ni en el plan: es PROVISIONAL (bloqueo S6-importar-partidas-csv).

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

var csvColumns = []string{"tipo", "concepto", "monto", "periodicidad", "dia_o_fecha", "categoria", "persona", "porcentaje_conjunto", "notas"}

var csvRequired = map[string]bool{"tipo": true, "concepto": true, "monto": true, "periodicidad": true, "dia_o_fecha": true, "categoria": true, "persona": true}

// csvField da la columna del CSV para un campo de PlannedItemCreate (los errores salen con la columna, U4 P-26.11).
var csvField = map[string]string{"kind": "tipo", "description": "concepto", "amount_cents": "monto",
	"second_half_amount_cents": "monto", "frequency": "periodicidad", "interval_count": "periodicidad",
	"start_date": "dia_o_fecha", "end_date": "dia_o_fecha", "category_id": "categoria", "person_id": "persona",
	"cost_share_rule_id": "persona", "payment_method": "persona", "credit_id": "persona", "term_count": "persona",
	"accrual_months": "periodicidad", "notes": "notas"}

var reMonto = regexp.MustCompile(`^\$?\s*([0-9]{1,3}(,[0-9]{3})+|[0-9]+)(\.([0-9]{1,2}))?$`)

// parseMonto convierte pesos de texto («1,234.5», «$950») a centavos enteros sin punto flotante.
func parseMonto(s string) (int64, bool) {
	m := reMonto.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0, false
	}
	pesos, err := strconv.ParseInt(strings.ReplaceAll(m[1], ",", ""), 10, 64)
	if err != nil || pesos > 1e12 {
		return 0, false
	}
	frac := m[4]
	for len(frac) < 2 {
		frac += "0"
	}
	c, _ := strconv.ParseInt(frac, 10, 64)
	return pesos*100 + c, true
}

type csvRow struct {
	line   int64
	item   dominio.PlannedItemRow
	ok     bool // se pudo armar la partida (aunque tenga errores de negocio)
	errors []api.FieldError
}

func (c *csvRow) add(col, code, msg string) {
	c.errors = append(c.errors, api.FieldError{Field: col, Code: code, Message: msg})
}

// importCatalog son los catálogos para resolver categoria y persona.
type importCatalog struct {
	cats    []dominio.CategoryRow
	persons map[string]dominio.PersonRecord
	rules   []dominio.ShareRuleRow
	today   dominio.Date
}

func loadCatalog(ctx context.Context, db dominio.Querier, today dominio.Date) (*importCatalog, error) {
	cats, err := dominio.ListCategoryRows(ctx, db, dominio.CategoryFilter{Status: "active", Today: today}, "ORDER BY c.id")
	if err != nil {
		return nil, err
	}
	ps, err := dominio.ListPersonRecords(ctx, db, 0, "ORDER BY p.id")
	if err != nil {
		return nil, err
	}
	tpl := true
	rules, err := dominio.ListShareRuleRows(ctx, db, &tpl, 0, "ORDER BY r.id")
	if err != nil {
		return nil, err
	}
	cat := &importCatalog{cats: cats, persons: map[string]dominio.PersonRecord{}, rules: rules, today: today}
	for _, p := range ps {
		cat.persons[p.Code] = p
	}
	return cat, nil
}

// rule busca la regla reutilizable activa por nombre (PROVISIONAL): «karen» → la llamada Karen; «conjunto» → la que
// empieza con «Conjunto»; «hogar» → la regla por defecto.
func (c *importCatalog) rule(persona string) (dominio.ShareRuleRow, int) {
	var found []dominio.ShareRuleRow
	for _, r := range c.rules {
		n := fold(r.Name)
		if r.IsActive && ((persona == "hogar" && r.IsDefault) || (persona == "karen" && n == "karen") ||
			(persona == "conjunto" && strings.HasPrefix(n, "conjunto"))) {
			found = append(found, r)
		}
	}
	if len(found) != 1 {
		return dominio.ShareRuleRow{}, len(found)
	}
	return found[0], 1
}

// startDate: AAAA-MM-DD tal cual; un día del mes (1 a 31) es la próxima fecha con ese día desde hoy (recortada al
// último día en meses cortos); en quincenal sin fecha, el día de pago de la quincena en curso (15 o último día).
func (c *importCatalog) startDate(s, freq string) (string, string) {
	s = strings.TrimSpace(s)
	if d, err := dominio.ParseDate(s); err == nil {
		return d.String(), ""
	}
	if s == "" {
		if freq == "fortnightly" {
			return dominio.FortnightOf(c.today).PayDay().String(), ""
		}
		return "", "Escribe el día del mes (1 a 31) o la fecha AAAA-MM-DD."
	}
	day, err := strconv.Atoi(s)
	if err != nil || day < 1 || day > 31 {
		return "", "Escribe el día del mes (1 a 31) o la fecha AAAA-MM-DD."
	}
	switch freq {
	case "once", "yearly":
		return "", "En una partida única o anual escribe la fecha completa AAAA-MM-DD."
	case "fortnightly":
		return dominio.FortnightOf(c.today).PayDay().String(), ""
	}
	y, m := c.today.Year, c.today.Month
	for i := 0; i < 2; i++ {
		d := min(day, dominio.DaysIn(y, m))
		cand, _ := dominio.ParseDate(fmt.Sprintf("%04d-%02d-%02d", y, int(m), d))
		if !cand.Before(c.today) {
			return cand.String(), ""
		}
		if m == 12 {
			y, m = y+1, 1
		} else {
			m++
		}
	}
	return "", "No se pudo calcular la fecha."
}

// mapRow arma la partida de un renglón; los errores de columna quedan en row.errors.
func (c *importCatalog) mapRow(row *csvRow, get func(string) string) {
	kinds := map[string]string{"ingreso": "income", "gasto_fijo": "fixed_expense", "gasto_venidero": "upcoming_expense"}
	freqs := map[string]string{"quincenal": "fortnightly", "mensual": "monthly", "anual": "yearly", "unica": "once"}
	p := dominio.PlannedItemRow{IntervalCount: 1, AccrualMonths: 1, IsActive: true}
	p.Kind = kinds[fold(get("tipo"))]
	if p.Kind == "" {
		row.add("tipo", "invalid_format", "El tipo es ingreso, gasto_fijo o gasto_venidero.")
	}
	p.Description = strings.TrimSpace(get("concepto"))
	if p.Description == "" {
		row.add("concepto", "required", "Escribe el concepto.")
	}
	amt, ok := parseMonto(get("monto"))
	if !ok || amt <= 0 {
		row.add("monto", "invalid_format", "El monto va en pesos con punto decimal y hasta dos decimales (p. ej. 1250.50).")
	}
	p.AmountCents = amt
	p.Frequency = freqs[fold(get("periodicidad"))]
	if p.Frequency == "" {
		row.add("periodicidad", "invalid_format", "La periodicidad es quincenal, mensual, anual o unica.")
	}
	if p.Frequency != "" {
		if d, msg := c.startDate(get("dia_o_fecha"), p.Frequency); msg != "" {
			row.add("dia_o_fecha", "invalid_format", msg)
		} else {
			p.StartDate = d
		}
	}
	catName := fold(get("categoria"))
	if p.Kind == "income" {
		if catName != "" {
			row.add("categoria", "not_allowed", "Un ingreso no lleva categoría: deja la columna vacía.")
		}
	} else if catName != "" {
		for _, x := range c.cats {
			if fold(x.Name) == catName || (x.Code.Valid && x.Code.String == catName) {
				p.CategoryID = sql.NullInt64{Int64: x.ID, Valid: true}
				break
			}
		}
		if !p.CategoryID.Valid {
			row.add("categoria", "not_found", "No hay una categoría activa con ese nombre.")
		}
	}
	persona := fold(get("persona"))
	pct := strings.TrimSpace(get("porcentaje_conjunto"))
	switch persona {
	case "axel", "ximena":
		if pr, ok := c.persons[persona]; ok {
			p.PersonID = sql.NullInt64{Int64: pr.ID, Valid: true}
		} else {
			row.add("persona", "not_found", "Esa persona no existe.")
		}
	case "hogar", "karen", "conjunto":
		if p.Kind == "income" {
			if persona != "hogar" {
				row.add("persona", "not_allowed", "Los ingresos van a nombre del titular o del hogar.")
			}
			break
		}
		r, n := c.rule(persona)
		if n != 1 {
			row.add("persona", "not_found", "No se encontró una sola regla de reparto para «"+persona+"» en Personas y reparto.")
			break
		}
		p.CostShareRuleID = sql.NullInt64{Int64: r.ID, Valid: true}
		if persona == "conjunto" && pct != "" && !pctMatches(pct, r, c.persons["karen"].ID) {
			row.add("porcentaje_conjunto", "not_allowed", "Por ahora el porcentaje debe ser el de la regla «"+r.Name+"»; cámbialo en Personas y reparto.")
		}
	default:
		row.add("persona", "invalid_format", "La persona es axel, ximena, hogar, karen o conjunto.")
	}
	if pct != "" && persona != "conjunto" {
		row.add("porcentaje_conjunto", "not_allowed", "El porcentaje solo va con persona «conjunto».")
	}
	if n := strings.TrimSpace(get("notas")); n != "" {
		p.Notes = sql.NullString{String: n, Valid: true}
	}
	row.item, row.ok = p, p.Kind != "" && p.Frequency != "" && p.StartDate != ""
}

// pctMatches compara porcentaje_conjunto (en %, hasta dos decimales) con la parte de Karen en la regla.
func pctMatches(s string, r dominio.ShareRuleRow, karenID int64) bool {
	bp, ok := parseMonto(strings.TrimSuffix(strings.TrimSpace(s), "%"))
	if !ok {
		return false
	}
	for _, p := range r.Parts {
		if p.PersonID == karenID && p.PctBp.Valid {
			return p.PctBp.Int64 == bp
		}
	}
	return false
}

// readUpload lee el campo file (text/csv) de un multipart/form-data: 415 si no es CSV, 413 si pasa del tope.
func readUpload(reg *api.Registry, r *http.Request) ([]byte, error) {
	mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mt != "multipart/form-data" {
		return nil, &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Sube el archivo .csv como formulario (multipart/form-data)."}
	}
	limit := reg.Deps().Config.MaxJSONBytes // PROVISIONAL: mismo tope que un cuerpo JSON (bloqueo S6-importar-partidas-csv)
	r.Body = http.MaxBytesReader(nil, r.Body, limit)
	mr, err := r.MultipartReader()
	if err != nil {
		return nil, api.BadRequest("El formulario no es válido.")
	}
	for {
		part, err := mr.NextPart()
		var mb *http.MaxBytesError
		switch {
		case errors.As(err, &mb):
			return nil, &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El archivo es demasiado grande."}
		case err == io.EOF:
			return nil, &api.Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "Falta el archivo CSV.",
				Fields: []api.FieldError{{Field: "file", Code: "required", Message: "Elige el archivo."}}}
		case err != nil:
			return nil, api.BadRequest("El formulario no es válido.")
		}
		if part.FormName() != "file" {
			continue
		}
		ct, _, _ := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if ct != "text/csv" {
			return nil, &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Solo se admiten archivos .csv (text/csv)."}
		}
		b, err := io.ReadAll(part)
		if errors.As(err, &mb) {
			return nil, &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El archivo es demasiado grande."}
		}
		return b, err
	}
}

// parseCSV lee el archivo y arma cada renglón (line_no cuenta el encabezado como 1).
func parseCSV(ctx context.Context, db dominio.Querier, today dominio.Date, b []byte) ([]*csvRow, error) {
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	rd := csv.NewReader(bytes.NewReader(b))
	if first, _, _ := bytes.Cut(b, []byte("\n")); bytes.Count(first, []byte(";")) > bytes.Count(first, []byte(",")) {
		rd.Comma = ';'
	}
	rd.FieldsPerRecord, rd.TrimLeadingSpace = -1, true
	recs, err := rd.ReadAll()
	if err != nil {
		var pe *csv.ParseError
		line := 0
		if errors.As(err, &pe) {
			line = pe.Line
		}
		return nil, dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "invalid_format", Message: fmt.Sprintf("El archivo no es un CSV válido (renglón %d).", line)})
	}
	if len(recs) == 0 {
		return nil, dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "required", Message: "El archivo está vacío."})
	}
	idx := map[string]int{}
	for i, h := range recs[0] {
		idx[fold(h)] = i
	}
	var missing []string
	for _, c := range csvColumns {
		if _, ok := idx[c]; !ok && csvRequired[c] {
			missing = append(missing, c)
		}
	}
	if len(missing) > 0 {
		return nil, dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "invalid_format", Message: "El archivo no tiene las columnas de la plantilla: falta " + strings.Join(missing, ", ") + "."})
	}
	if len(recs) == 1 {
		return nil, dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "required", Message: "El archivo no tiene renglones después del encabezado."})
	}
	cat, err := loadCatalog(ctx, db, today)
	if err != nil {
		return nil, err
	}
	rows := make([]*csvRow, 0, len(recs)-1)
	for n, rec := range recs[1:] {
		row := &csvRow{line: int64(n + 2), errors: []api.FieldError{}}
		get := func(col string) string {
			if i, ok := idx[col]; ok && i < len(rec) {
				return rec[i]
			}
			return ""
		}
		cat.mapRow(row, get)
		if row.ok {
			if err := plannedEnums(row.item); err != nil {
				return nil, err
			}
			issues, err := plannedIssues(ctx, db, row.item)
			if err != nil {
				return nil, err
			}
			for _, is := range issues {
				col := csvField[is.Field]
				if col == "" {
					col = is.Field
				}
				if !hasErr(row, col) {
					row.add(col, is.Code, is.Message)
				}
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func hasErr(row *csvRow, col string) bool {
	for _, e := range row.errors {
		if e.Field == col {
			return true
		}
	}
	return false
}

func plannedCreateOf(p dominio.PlannedItemRow) *api.PlannedItemCreate {
	return &api.PlannedItemCreate{Kind: p.Kind, Description: p.Description, AmountCents: p.AmountCents, Frequency: p.Frequency,
		IntervalCount: api.Ptr(p.IntervalCount), StartDate: p.StartDate, CategoryID: pInt(p.CategoryID), PersonID: pInt(p.PersonID),
		CostShareRuleID: pInt(p.CostShareRuleID), AccrualMonths: api.Ptr(p.AccrualMonths), IsActive: api.Ptr(p.IsActive), Notes: pStr(p.Notes)}
}

// previewPlannedItemsImport: vista previa sin escribir nada.
func previewPlannedItemsImport(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		b, err := readUpload(reg, r)
		if err != nil {
			return err
		}
		rows, err := parseCSV(r.Context(), reg.Deps().DB, today(reg), b)
		if err != nil {
			return err
		}
		out := api.ImportPreview{Rows: make([]api.ImportPreviewRowsItem, 0, len(rows))}
		for _, row := range rows {
			it := api.ImportPreviewRowsItem{LineNo: row.line, Errors: row.errors}
			if row.ok {
				it.PlannedItem = plannedCreateOf(row.item)
			}
			if len(row.errors) == 0 {
				out.ValidCount++
			} else {
				out.ErrorCount++
			}
			out.Rows = append(out.Rows, it)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// importPlannedItems: mismo archivo; con cualquier renglón con error no se importa nada (422).
func importPlannedItems(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		b, err := readUpload(reg, r)
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		rows, err := parseCSV(ctx, db, today(reg), b)
		if err != nil {
			return err
		}
		var issues []dominio.FieldIssue
		for _, row := range rows {
			for _, e := range row.errors {
				issues = append(issues, dominio.FieldIssue{Field: fmt.Sprintf("rows[%d].%s", row.line, e.Field), Code: e.Code, Message: e.Message})
			}
		}
		if len(issues) > 0 {
			return dominio.Invalid(issues...)
		}
		res := api.ImportResult{PlannedItemIDs: []int64{}}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			for _, row := range rows {
				id, err := dominio.InsertOwned(ctx, tx, "planned_item", row.item.Columns(), api.Actor(r), api.RequestID(r))
				if err != nil {
					return err
				}
				res.PlannedItemIDs = append(res.PlannedItemIDs, id)
			}
			return nil
		})
		if err := writeErr(err, "Ya existe esa partida."); err != nil {
			return err
		}
		res.CreatedCount = int64(len(res.PlannedItemIDs))
		return api.WriteJSON(w, http.StatusCreated, res)
	}
}
