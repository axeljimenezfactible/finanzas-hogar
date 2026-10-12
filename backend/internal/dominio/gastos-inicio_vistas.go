package dominio

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// ViewSQL describe la consulta genérica de una vista de base de datos (queryView, C1 §8).
// Las expresiones usan el alias de la tabla (t por defecto); las vacías se leen como NULL.
type ViewSQL struct {
	Table      string
	Alias      string // alias de la tabla (t si va vacío)
	Title      string // título del registro
	Amount     string // importe principal (C1 §8), en centavos
	Date       string // campo de fecha del calendario (o del periodo)
	Group      string // clave del grupo del tablero; varias claves separadas por comas = el registro va en cada columna
	GroupLabel string // etiqueta del grupo (nombre del registro al que apunta un id)
	Image      string // adjunto con foto (galería)
	Extra      string // dato auxiliar para valores calculados por un motor (p. ej. último precio)
	Where      []string
	Args       []any
	OrderBy    string
}

// ViewRecord es un renglón de una vista con su registro completo y su proyección.
type ViewRecord struct {
	ID         int64
	Record     map[string]any
	Title      string
	Amount     sql.NullInt64
	Date       sql.NullString
	Group      sql.NullString
	GroupLabel sql.NullString
	Image      sql.NullInt64
	Extra      sql.NullString
}

// viewHidden son columnas que nunca salen en una vista (credenciales y datos internos).
var viewHidden = map[string]bool{"password_hash": true, "token_hash": true, "idempotency_key": true, "ocr_draft_json": true}

// viewBool son columnas 0/1 que el contrato publica como booleanas además de las is_*, has_* y requires_*.
var viewBool = map[string]bool{"pays_share_directly": true, "availability_confirmed": true, "assisted_reading_enabled": true,
	"interest_includes_vat": true, "vat_included": true}

func viewIsBool(col string) bool {
	return viewBool[col] || strings.HasPrefix(col, "is_") || strings.HasPrefix(col, "has_") || strings.HasPrefix(col, "requires_")
}

func viewExpr(e string) string {
	if e == "" {
		return "NULL"
	}
	return "(" + e + ")"
}

// GastosViewRows lee todos los renglones de la vista con los filtros SQL (la paginación y los filtros calculados van en la API).
func GastosViewRows(ctx context.Context, q Querier, s ViewSQL) ([]ViewRecord, error) {
	where := "1 = 1"
	if len(s.Where) > 0 {
		where = strings.Join(s.Where, " AND ")
	}
	a := s.Alias
	if a == "" {
		a = "t"
	}
	query := fmt.Sprintf("SELECT %s.*, %s, %s, %s, %s, %s, %s, %s FROM %s %s WHERE %s %s", a, viewExpr(s.Title), viewExpr(s.Amount),
		viewExpr(s.Date), viewExpr(s.Group), viewExpr(s.GroupLabel), viewExpr(s.Image), viewExpr(s.Extra), s.Table, a, where, s.OrderBy)
	rows, err := q.QueryContext(ctx, query, s.Args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	n := len(cols) - 7
	out := []ViewRecord{}
	for rows.Next() {
		vals := make([]any, n)
		dst := make([]any, 0, len(cols))
		for i := range vals {
			dst = append(dst, &vals[i])
		}
		var v ViewRecord
		var title sql.NullString
		dst = append(dst, &title, &v.Amount, &v.Date, &v.Group, &v.GroupLabel, &v.Image, &v.Extra)
		if err := rows.Scan(dst...); err != nil {
			return nil, err
		}
		v.Title = title.String
		v.Record = map[string]any{}
		for i, c := range cols[:n] {
			if viewHidden[c] || vals[i] == nil { // NULL → propiedad ausente (C1 §1 punto 6)
				continue
			}
			val := vals[i]
			switch x := val.(type) {
			case []byte:
				val = string(x)
			case time.Time:
				val = FormatDBTime(x)
			}
			if sv, ok := val.(string); ok && strings.HasSuffix(c, "_at") {
				val = DBTimeToAPI(sv)
			}
			if iv, ok := val.(int64); ok && viewIsBool(c) {
				val = iv != 0
			}
			if c == "id" {
				v.ID, _ = val.(int64)
			}
			v.Record[c] = val
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
