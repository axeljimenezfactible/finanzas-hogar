package dominio

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
)

var reIdent = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Audit escribe un renglón de audit_log (las migraciones no traen disparadores de bitácora: la API la escribe).
// action: insert (solo after), update (before y after) o delete (solo before). actor nil = proceso del sistema.
func Audit(ctx context.Context, q Querier, actor *int64, action, table string, rowID int64, before, after json.RawMessage, requestID string) error {
	var b, a, rid any
	if before != nil {
		b = string(before)
	}
	if after != nil {
		a = string(after)
	}
	if requestID != "" {
		rid = requestID
	}
	_, err := q.ExecContext(ctx, `INSERT INTO audit_log (user_account_id, action, table_name, row_id, before_json, after_json, request_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		actor, action, table, rowID, b, a, rid)
	return err
}

// RowJSON lee un renglón completo como objeto JSON (para before_json y after_json de audit_log).
func RowJSON(ctx context.Context, q Querier, table string, id int64) (json.RawMessage, error) {
	if !reIdent.MatchString(table) {
		return nil, fmt.Errorf("tabla inválida %q", table)
	}
	rows, err := q.QueryContext(ctx, `SELECT * FROM `+table+` WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, ErrNotFound
	}
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return nil, err
	}
	m := make(map[string]any, len(cols))
	for i, c := range cols {
		if b, ok := vals[i].([]byte); ok {
			vals[i] = string(b)
		}
		m[c] = vals[i]
	}
	return json.Marshal(m)
}
