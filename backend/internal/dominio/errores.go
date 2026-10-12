package dominio

import (
	"database/sql"
	"errors"
)

// Kind clasifica los errores de negocio; la API los traduce al esquema único de error de C1 (§3).
type Kind int

// Clases de error y el código HTTP que les da la API.
const (
	KindValidation  Kind = iota + 1 // 422 validation_failed (u otro código 422 de C1)
	KindNotFound                    // 404 not_found
	KindConflict                    // 409 conflict, in_use, already_exists, …
	KindForbidden                   // 403 forbidden, read_only_source
	KindUnavailable                 // 503 unavailable
	KindBadRequest                  // 400 bad_request (formato)
)

// FieldIssue es el detalle por campo (se vuelve FieldError en la respuesta).
type FieldIssue struct {
	Field   string
	Code    string // required, invalid_format, out_of_range, not_found, not_allowed, conflict
	Message string // español de México
}

// Error es un error de negocio con código estable en inglés y mensaje para la persona usuaria en español de México.
type Error struct {
	Kind    Kind
	Code    string
	Message string
	Fields  []FieldIssue
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + ": " + e.Err.Error()
	}
	return e.Code + ": " + e.Message
}

func (e *Error) Unwrap() error { return e.Err }

// ErrNotFound se devuelve cuando un registro no existe (equivale a sql.ErrNoRows).
var ErrNotFound = &Error{Kind: KindNotFound, Code: "not_found", Message: "No encontramos ese registro."}

// NotFound devuelve un 404 con un mensaje propio.
func NotFound(msg string) *Error { return &Error{Kind: KindNotFound, Code: "not_found", Message: msg} }

// Conflict devuelve un 409 con código estable (conflict, in_use, already_exists, …).
func Conflict(code, msg string) *Error { return &Error{Kind: KindConflict, Code: code, Message: msg} }

// Forbidden devuelve un 403 (forbidden o read_only_source).
func Forbidden(code, msg string) *Error { return &Error{Kind: KindForbidden, Code: code, Message: msg} }

// Unavailable devuelve un 503 unavailable.
func Unavailable(msg string) *Error {
	return &Error{Kind: KindUnavailable, Code: "unavailable", Message: msg}
}

// BadRequest devuelve un 400 bad_request (formato incorrecto).
func BadRequest(msg string) *Error {
	return &Error{Kind: KindBadRequest, Code: "bad_request", Message: msg}
}

// Invalid devuelve un 422 validation_failed con el detalle por campo.
func Invalid(fields ...FieldIssue) *Error {
	return &Error{Kind: KindValidation, Code: "validation_failed", Message: "Revisa los campos marcados.", Fields: fields}
}

// AsError convierte cualquier error en *Error cuando es de negocio; sql.ErrNoRows se vuelve ErrNotFound.
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound, true
	}
	return nil, false
}
