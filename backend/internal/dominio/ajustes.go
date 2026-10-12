package dominio

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// ParamDef describe un parámetro operativo de app_setting (K44) tal como lo anotan las migraciones:
//
//	-- @param <clave> | <tipo> | <unidad> | <valor inicial JSON> | <decisión> | <sección de Ajustes> | <migración>
//	-- @param_rango <clave> | <rango aprobado> | <decisión>
//
// No hay una tabla escrita a mano: la fuente es la misma que usa el generador de C1 (§11.7).
type ParamDef struct {
	Key       string
	Type      string // entero, lista de enteros, hora HH:MM
	Unit      string
	Initial   json.RawMessage
	Decision  string
	Section   string
	Migration string
	Range     string // texto de @param_rango, si lo hay (el rango lo hacen cumplir los disparadores de S1)
}

// ValueType es el value_type de C1: integer, integer_list o time.
func (p ParamDef) ValueType() string {
	switch {
	case strings.HasPrefix(p.Type, "lista"):
		return "integer_list"
	case strings.HasPrefix(p.Type, "hora"):
		return "time"
	}
	return "integer"
}

var (
	paramOnce sync.Once
	paramDefs map[string]ParamDef
	paramErr  error
	reHHMM    = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)
)

// ParamDefs lee las anotaciones @param de las migraciones embebidas (una vez por proceso).
func ParamDefs() (map[string]ParamDef, error) {
	paramOnce.Do(func() { paramDefs, paramErr = parseParamDefs(MigrationsFS()) })
	return paramDefs, paramErr
}

func parseParamDefs(fsys fs.FS) (map[string]ParamDef, error) {
	names, err := fs.Glob(fsys, "*.sql")
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	defs := map[string]ParamDef{}
	ranges := map[string]string{}
	for _, n := range names {
		b, err := fs.ReadFile(fsys, n)
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(bytes.NewReader(b))
		sc.Buffer(make([]byte, 1<<20), 1<<20)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			switch {
			case strings.HasPrefix(line, "-- @param "):
				f := splitBar(strings.TrimPrefix(line, "-- @param "))
				if len(f) != 7 || strings.HasPrefix(f[0], "<") {
					continue // la línea que documenta el formato
				}
				if !json.Valid([]byte(f[3])) {
					return nil, fmt.Errorf("%s: valor inicial de %s no es JSON: %s", n, f[0], f[3])
				}
				defs[f[0]] = ParamDef{Key: f[0], Type: f[1], Unit: f[2], Initial: json.RawMessage(f[3]), Decision: f[4], Section: f[5], Migration: f[6]}
			case strings.HasPrefix(line, "-- @param_rango "):
				f := splitBar(strings.TrimPrefix(line, "-- @param_rango "))
				if len(f) >= 2 && !strings.HasPrefix(f[0], "<") {
					ranges[f[0]] = f[1]
				}
			}
		}
		if err := sc.Err(); err != nil {
			return nil, err
		}
	}
	for k, r := range ranges {
		if d, ok := defs[k]; ok {
			d.Range = r
			defs[k] = d
		}
	}
	return defs, nil
}

func splitBar(s string) []string {
	parts := strings.Split(s, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// SettingRaw lee el valor JSON vigente de una clave del hogar (user_account_id NULL). found=false si no existe.
func SettingRaw(ctx context.Context, q Querier, key string) (json.RawMessage, bool, error) {
	var v string
	err := q.QueryRowContext(ctx, `SELECT value_json FROM app_setting WHERE setting_key = ? AND user_account_id IS NULL`, key).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return json.RawMessage(v), true, nil
}

func settingInto(ctx context.Context, q Querier, key string, dst any) error {
	raw, ok, err := SettingRaw(ctx, q, key)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("app_setting %s no existe", key)
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("app_setting %s: %w", key, err)
	}
	return nil
}

// SettingInt lee una clave entera (auth.max_failed_attempts, session.max_age_days, attachment.max_size_mb, …).
func SettingInt(ctx context.Context, q Querier, key string) (int64, error) {
	var n int64
	return n, settingInto(ctx, q, key, &n)
}

// SettingIntList lee una clave de lista de enteros (alerts.payment_due_days_before).
func SettingIntList(ctx context.Context, q Querier, key string) ([]int64, error) {
	var n []int64
	return n, settingInto(ctx, q, key, &n)
}

// SettingTime lee una clave de hora HH:MM (reports.daily_send_time, alerts.send_time).
func SettingTime(ctx context.Context, q Querier, key string) (string, error) {
	var s string
	if err := settingInto(ctx, q, key, &s); err != nil {
		return "", err
	}
	if !reHHMM.MatchString(s) {
		return "", fmt.Errorf("app_setting %s: %q no es HH:MM", key, s)
	}
	return s, nil
}

// SettingBool lee una clave booleana (ai.local_model_enabled, ai.claude_api_enabled, price.assisted_reading_enabled).
func SettingBool(ctx context.Context, q Querier, key string) (bool, error) {
	var b bool
	return b, settingInto(ctx, q, key, &b)
}

// Parameter es la lectura de un parámetro para los GET de settings (SettingParameter* de C1 §11.7).
type Parameter struct {
	Key          string
	ValueType    string
	Unit         string
	Value        json.RawMessage
	InitialValue json.RawMessage // nil si no hay valor inicial anotado
	IsInitial    bool
}

// ReadParameter lee el valor vigente y el inicial de una clave con @param.
func ReadParameter(ctx context.Context, q Querier, key string) (Parameter, error) {
	defs, err := ParamDefs()
	if err != nil {
		return Parameter{}, err
	}
	def, ok := defs[key]
	if !ok {
		return Parameter{}, fmt.Errorf("app_setting %s no es un parámetro anotado con @param", key)
	}
	raw, found, err := SettingRaw(ctx, q, key)
	if err != nil {
		return Parameter{}, err
	}
	if !found {
		raw = def.Initial
	}
	return Parameter{Key: key, ValueType: def.ValueType(), Unit: def.Unit, Value: raw, InitialValue: def.Initial, IsInitial: jsonEqual(raw, def.Initial)}, nil
}

func jsonEqual(a, b json.RawMessage) bool {
	var x, y any
	if json.Unmarshal(a, &x) != nil || json.Unmarshal(b, &y) != nil {
		return false
	}
	xa, _ := json.Marshal(x)
	ya, _ := json.Marshal(y)
	return bytes.Equal(xa, ya)
}

// ParamChange es un cambio por PATCH de settings: parameters (valor nuevo por clave) y reset_parameters (vuelven al inicial).
type ParamChange struct {
	Set   map[string]json.RawMessage
	Reset []string
}

// ApplyParams valida y guarda un cambio de parámetros dentro de tx, con audit_log por clave (C1 §11.7).
// allowed limita las claves a las de la sección de la operación. Errores de valor: 422 validation_failed con
// fields «parameters.<clave>» (o «reset_parameters»). El rango lo hacen cumplir los disparadores de S1.
func ApplyParams(ctx context.Context, tx *sql.Tx, ch ParamChange, allowed []string, actor *int64, requestID string) error {
	defs, err := ParamDefs()
	if err != nil {
		return err
	}
	ok := map[string]bool{}
	for _, k := range allowed {
		ok[k] = true
	}
	var issues []FieldIssue
	final := map[string]json.RawMessage{}
	for k, v := range ch.Set {
		f := "parameters." + k
		def, known := defs[k]
		switch {
		case !known || !ok[k]:
			issues = append(issues, FieldIssue{f, "not_allowed", "Este parámetro no se puede cambiar aquí."})
		case !validParamValue(def.ValueType(), v):
			issues = append(issues, FieldIssue{f, "invalid_format", "El valor no tiene el tipo correcto para este parámetro."})
		default:
			final[k] = v
		}
	}
	for _, k := range ch.Reset {
		def, known := defs[k]
		switch {
		case !known || !ok[k]:
			issues = append(issues, FieldIssue{"reset_parameters", "not_allowed", "El parámetro " + k + " no se puede restablecer aquí."})
		case ch.Set[k] != nil:
			issues = append(issues, FieldIssue{"reset_parameters", "conflict", "El parámetro " + k + " no puede cambiarse y restablecerse a la vez."})
		default:
			final[k] = def.Initial
		}
	}
	if len(issues) > 0 {
		return Invalid(issues...)
	}
	keys := make([]string, 0, len(final))
	for k := range final {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := setSetting(ctx, tx, k, final[k], actor, requestID); err != nil {
			if IsConstraint(err) {
				return Invalid(FieldIssue{"parameters." + k, "out_of_range", "El valor está fuera del rango permitido para este parámetro."})
			}
			return err
		}
	}
	return nil
}

func validParamValue(vt string, v json.RawMessage) bool {
	switch vt {
	case "integer":
		var n int64
		return json.Unmarshal(v, &n) == nil && !bytes.Equal(bytes.TrimSpace(v), []byte("null"))
	case "integer_list":
		var l []int64
		return json.Unmarshal(v, &l) == nil && l != nil
	case "time":
		var s string
		return json.Unmarshal(v, &s) == nil && reHHMM.MatchString(s)
	}
	return false
}

// SetSetting guarda el valor JSON de una clave del hogar (la crea si no existe) y registra el cambio en audit_log.
// Para parámetros de Ajustes use ApplyParams; SetSetting es para claves sin @param (ai.*, optimizer.weights, …).
func SetSetting(ctx context.Context, tx *sql.Tx, key string, value json.RawMessage, actor *int64, requestID string) error {
	if !json.Valid(value) {
		return Invalid(FieldIssue{key, "invalid_format", "El valor no es válido."})
	}
	return setSetting(ctx, tx, key, value, actor, requestID)
}

func setSetting(ctx context.Context, tx *sql.Tx, key string, value json.RawMessage, actor *int64, requestID string) error {
	var id int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM app_setting WHERE setting_key = ? AND user_account_id IS NULL`, key).Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		res, err := tx.ExecContext(ctx, `INSERT INTO app_setting (setting_key, value_json) VALUES (?, ?)`, key, string(value))
		if err != nil {
			return err
		}
		id, _ = res.LastInsertId()
		after, err := RowJSON(ctx, tx, "app_setting", id)
		if err != nil {
			return err
		}
		return Audit(ctx, tx, actor, "insert", "app_setting", id, nil, after, requestID)
	case err != nil:
		return err
	}
	before, err := RowJSON(ctx, tx, "app_setting", id)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE app_setting SET value_json = ? WHERE id = ?`, string(value), id); err != nil {
		return err
	}
	after, err := RowJSON(ctx, tx, "app_setting", id)
	if err != nil {
		return err
	}
	return Audit(ctx, tx, actor, "update", "app_setting", id, before, after, requestID)
}
