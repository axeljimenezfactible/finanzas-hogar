// Package ajustes implementa las operaciones de Ajustes de C1 v1.4 (agente api-ajustes): categorías, personas,
// reglas de reparto, partidas planificadas e importación de gastos_fijos_ingresos.csv, metas de ahorro y los
// ajustes de proyección y de avisos (app_setting).
package ajustes

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"io"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

// Register reclama las 19 operaciones del área (CONVENCIONES-backend §14).
func Register(reg *api.Registry) {
	reg.Handle("listCategories", listCategories(reg))
	reg.Handle("createCategory", createCategory(reg))
	reg.Handle("updateCategory", updateCategory(reg))
	reg.Handle("reorderCategories", reorderCategories(reg))
	reg.Handle("listPersons", listPersons(reg))
	reg.Handle("updatePerson", updatePerson(reg))
	reg.Handle("listCostShareRules", listCostShareRules(reg))
	reg.Handle("updateCostShareRule", updateCostShareRule(reg))
	reg.Handle("listPlannedItems", listPlannedItems(reg))
	reg.Handle("createPlannedItem", createPlannedItem(reg))
	reg.Handle("updatePlannedItem", updatePlannedItem(reg))
	reg.Handle("previewPlannedItemsImport", previewPlannedItemsImport(reg))
	reg.Handle("importPlannedItems", importPlannedItems(reg))
	reg.Handle("getProjectionSettings", getProjectionSettings(reg))
	reg.Handle("updateProjectionSettings", updateProjectionSettings(reg))
	reg.Handle("createSavingsGoal", createSavingsGoal(reg))
	reg.Handle("updateSavingsGoal", updateSavingsGoal(reg))
	reg.Handle("getAlertSettings", getAlertSettings(reg))
	reg.Handle("updateAlertSettings", updateAlertSettings(reg))
}

// decode lee el cuerpo JSON en dst (reglas de reg.DecodeJSON: 400/413) y devuelve además las claves presentes,
// para distinguir «ausente», «null» y «vacío» (requeridos en alta y JSON Merge Patch sin cambios).
func decode(reg *api.Registry, r *http.Request, dst any) (map[string]json.RawMessage, error) {
	limit := reg.Deps().Config.MaxJSONBytes
	b, err := io.ReadAll(io.LimitReader(r.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El cuerpo de la petición es demasiado grande."}
	}
	r.Body = io.NopCloser(bytes.NewReader(b))
	if err := reg.DecodeJSON(r, dst); err != nil {
		return nil, err
	}
	keys := map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &keys); err != nil {
		return nil, api.BadRequest("El cuerpo debe ser un objeto JSON.")
	}
	return keys, nil
}

// noChanges es el 400 de un PATCH vacío (minProperties: 1 en C1).
func noChanges() error {
	return api.BadRequest("No hay cambios que guardar: manda al menos un campo.")
}

// badEnum es el 400 de un valor fuera de la lista de C1 (formato, no regla de negocio).
func badEnum(field string) error {
	return &api.Problem{Status: http.StatusBadRequest, Code: "bad_request", Message: "El valor de " + field + " no es uno de los permitidos.",
		Fields: []api.FieldError{{Field: field, Code: "invalid_format", Message: "Valor no permitido."}}}
}

func oneOf(v string, allowed ...string) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}

func today(reg *api.Registry) dominio.Date { return dominio.Today(reg.Deps().Clock) }

func pInt(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return api.Ptr(n.Int64)
}

func pStr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return api.Ptr(s.String)
}

func nInt(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

func nStr(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

// mergeInt aplica un campo Nullable de un PATCH sobre la columna.
func mergeInt(cur sql.NullInt64, v api.Nullable[int64]) sql.NullInt64 {
	if !v.Set {
		return cur
	}
	return sql.NullInt64{Int64: v.Value, Valid: v.Valid}
}

func mergeStr(cur sql.NullString, v api.Nullable[string]) sql.NullString {
	if !v.Set {
		return cur
	}
	return sql.NullString{String: v.Value, Valid: v.Valid}
}

// writeErr traduce violaciones de la base: UNIQUE → 409 already_exists; CHECK o disparador → 422.
func writeErr(err error, unique string) error {
	switch {
	case err == nil:
		return nil
	case dominio.IsUnique(err):
		return dominio.Conflict("already_exists", unique)
	case dominio.IsConstraint(err):
		return dominio.Invalid(dominio.FieldIssue{Field: "body", Code: "conflict", Message: "Los datos no cumplen una regla del registro. Revisa los campos e intenta de nuevo."})
	}
	return err
}

// refCheck valida que una referencia exista (422 not_found en el campo).
func refCheck(r *http.Request, reg *api.Registry, v *api.Validator, table string, id sql.NullInt64, field, msg string) error {
	if !id.Valid {
		return nil
	}
	ok, err := dominio.RefExists(r.Context(), reg.Deps().DB, table, id.Int64)
	if err != nil {
		return err
	}
	v.Check(ok, field, "not_found", msg)
	return nil
}

// fold quita mayúsculas y acentos para la búsqueda q («sin distinguir mayúsculas ni acentos», C1).
func fold(s string) string {
	return strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n").Replace(strings.ToLower(strings.TrimSpace(s)))
}
