// Package ia define los proveedores de IA opcionales (decisión H): modelo local (Ollama) y API de Claude, solo para
// OCR de tickets y sugerencias de recetas. Los dos están APAGADOS por defecto (app_setting ai.local_model_enabled y
// ai.claude_api_enabled = false) y aquí solo existe el adaptador «apagado»: con la IA apagada no hay ninguna llamada
// de salida (S6-13). Ninguna llave vive en el código, la configuración versionada ni la base.
package ia

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"finanzas-hogar/internal/dominio"
)

// Códigos de proveedor (parámetro {provider} de updateAiProvider en C1).
const (
	LocalModel = "local_model"
	ClaudeAPI  = "claude_api"
)

// ErrDisabled: la IA está apagada o no hay adaptador; la API responde 409 ai_disabled (C1 §3).
var ErrDisabled = errors.New("ia: apagada")

// TicketImage es la imagen de un ticket para OCR.
type TicketImage struct {
	Data     []byte
	MIMEType string
}

// TicketReading es la lectura de un ticket (el formato fino lo fija el área de gastos con el motor de OCR de S5).
type TicketReading struct {
	Text string
	Raw  json.RawMessage
}

// RecipeRequest pide una sugerencia de receta.
type RecipeRequest struct {
	Prompt   string
	Servings int
}

// RecipeSuggestion es la sugerencia devuelta (el área de comidas la convierte en RecipeCreate).
type RecipeSuggestion struct {
	Raw json.RawMessage
}

// Provider es un adaptador de IA.
type Provider interface {
	Code() string
	ReadTicket(ctx context.Context, img TicketImage) (TicketReading, error)
	SuggestRecipe(ctx context.Context, req RecipeRequest) (RecipeSuggestion, error)
}

// Off es el adaptador apagado: no llama a nada y siempre devuelve ErrDisabled.
type Off struct{ ProviderCode string }

// Code devuelve el código del proveedor.
func (o Off) Code() string { return o.ProviderCode }

// ReadTicket devuelve ErrDisabled.
func (Off) ReadTicket(context.Context, TicketImage) (TicketReading, error) {
	return TicketReading{}, ErrDisabled
}

// SuggestRecipe devuelve ErrDisabled.
func (Off) SuggestRecipe(context.Context, RecipeRequest) (RecipeSuggestion, error) {
	return RecipeSuggestion{}, ErrDisabled
}

// Registry elige el proveedor según los interruptores de app_setting.
type Registry struct {
	db        *sql.DB
	providers map[string]Provider
}

// NewRegistry crea el registro con los dos proveedores apagados.
func NewRegistry(db *sql.DB) *Registry {
	return &Registry{db: db, providers: map[string]Provider{LocalModel: Off{LocalModel}, ClaudeAPI: Off{ClaudeAPI}}}
}

// Set instala un adaptador real (fuera de alcance de WF-3a; ninguno existe todavía).
func (r *Registry) Set(p Provider) { r.providers[p.Code()] = p }

// SettingKey es la clave de app_setting del interruptor de un proveedor.
func SettingKey(code string) string { return "ai." + code + "_enabled" }

// Enabled lee el interruptor del proveedor en app_setting (false si falta).
func (r *Registry) Enabled(ctx context.Context, code string) (bool, error) {
	if _, ok := r.providers[code]; !ok {
		return false, errors.New("ia: proveedor desconocido " + code)
	}
	raw, ok, err := dominio.SettingRaw(ctx, r.db, SettingKey(code))
	if err != nil || !ok {
		return false, err
	}
	var on bool
	return on, json.Unmarshal(raw, &on)
}

// Available dice si hay un adaptador real instalado (no Off).
func (r *Registry) Available(code string) bool { _, off := r.providers[code].(Off); return !off }

// Active devuelve el primer proveedor encendido y disponible (local primero); si ninguno, ErrDisabled sin llamar a nada.
func (r *Registry) Active(ctx context.Context) (Provider, error) {
	for _, code := range []string{LocalModel, ClaudeAPI} {
		on, err := r.Enabled(ctx, code)
		if err != nil {
			return nil, err
		}
		if on && r.Available(code) {
			return r.providers[code], nil
		}
	}
	return nil, ErrDisabled
}
