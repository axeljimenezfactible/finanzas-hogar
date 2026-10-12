package api

import (
	"bytes"
	"encoding/json"
)

// Nullable representa un campo que admite null en C1 (tipo ["T","null"] o anyOf con null), con tres estados:
// ausente (Set=false), null (Set=true, Valid=false) o con valor (Valid=true). Sirve para JSON Merge Patch (C1 §1 punto 7:
// null borra un valor) y para respuestas: con la etiqueta omitzero un campo sin Set ni Valid se omite (NULL se omite, punto 6).
type Nullable[T any] struct {
	Value T
	Valid bool
	Set   bool
}

// Some devuelve un Nullable con valor.
func Some[T any](v T) Nullable[T] { return Nullable[T]{Value: v, Valid: true, Set: true} }

// Null devuelve un Nullable explícitamente null.
func Null[T any]() Nullable[T] { return Nullable[T]{Set: true} }

// IsZero hace que omitzero omita el campo cuando no se fijó.
func (n Nullable[T]) IsZero() bool { return !n.Set && !n.Valid }

// MarshalJSON escribe el valor o null.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

// UnmarshalJSON distingue null de un valor; la ausencia la deja el decodificador sin tocar (Set=false).
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true
	if bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
		var zero T
		n.Value, n.Valid = zero, false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(b, &n.Value)
}

// Ptr devuelve un puntero a v (para los campos opcionales no nulos de types_gen.go).
func Ptr[T any](v T) *T { return &v }
