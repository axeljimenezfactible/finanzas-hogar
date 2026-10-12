package ia

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

// Doer es el cliente HTTP de los adaptadores; las pruebas ponen uno que falla si se invoca (S6-13).
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// ErrNotConfigured: el adaptador existe pero le falta configuración (modelo, dirección o llave) o la llamada real
// aún no está decidida (PROVISIONAL, bloqueo S6-sistema).
var ErrNotConfigured = errors.New("ia: el adaptador no está configurado")

// Pinger es un adaptador que puede decir si responde.
type Pinger interface {
	Ping(ctx context.Context) error
}

// gate dice si el interruptor del proveedor (app_setting ai.*_enabled, apagado por defecto) está encendido.
type gate func(ctx context.Context) (bool, error)

func (g gate) check(ctx context.Context) error {
	on, err := g(ctx)
	if err != nil {
		return err
	}
	if !on {
		return ErrDisabled // apagada: ninguna llamada de salida
	}
	return nil
}

// Ollama es el adaptador del modelo local. Sin llave. BaseURL y Model vienen de la configuración del servidor
// (nunca del código); vacíos = no configurado.
type Ollama struct {
	BaseURL, Model string
	HTTP           Doer
	on             gate
}

// NewOllama crea el adaptador ligado al interruptor ai.local_model_enabled del registro.
func NewOllama(r *Registry, baseURL, model string, hc Doer) *Ollama {
	return &Ollama{BaseURL: strings.TrimRight(baseURL, "/"), Model: model, HTTP: hc, on: func(ctx context.Context) (bool, error) { return r.Enabled(ctx, LocalModel) }}
}

func (o *Ollama) Code() string { return LocalModel }

// Ping consulta GET {BaseURL}/api/tags solo si el interruptor está encendido.
func (o *Ollama) Ping(ctx context.Context) error {
	if err := o.on.check(ctx); err != nil {
		return err
	}
	if o.BaseURL == "" || o.HTTP == nil {
		return ErrNotConfigured
	}
	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(c, http.MethodGet, o.BaseURL+"/api/tags", nil)
	if err != nil {
		return err
	}
	res, err := o.HTTP.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("ia: el modelo local respondió " + res.Status)
	}
	return nil
}

// ReadTicket y SuggestRecipe: con el interruptor apagado devuelven ErrDisabled sin llamar a nada. La llamada real
// (modelo, instrucciones y formato) está PROVISIONAL: bloqueo S6-sistema.
func (o *Ollama) ReadTicket(ctx context.Context, _ TicketImage) (TicketReading, error) {
	if err := o.on.check(ctx); err != nil {
		return TicketReading{}, err
	}
	return TicketReading{}, ErrNotConfigured
}

func (o *Ollama) SuggestRecipe(ctx context.Context, _ RecipeRequest) (RecipeSuggestion, error) {
	if err := o.on.check(ctx); err != nil {
		return RecipeSuggestion{}, err
	}
	return RecipeSuggestion{}, ErrNotConfigured
}

// Claude es el adaptador de la API de Claude. La llave nunca va en el código ni en la configuración versionada:
// KeyFile es la ruta de un archivo fuera del repositorio que se lee al momento y nunca se imprime.
type Claude struct {
	KeyFile string
	HTTP    Doer
	on      gate
}

// NewClaude crea el adaptador ligado al interruptor ai.claude_api_enabled.
func NewClaude(r *Registry, keyFile string, hc Doer) *Claude {
	return &Claude{KeyFile: keyFile, HTTP: hc, on: func(ctx context.Context) (bool, error) { return r.Enabled(ctx, ClaudeAPI) }}
}

func (c *Claude) Code() string { return ClaudeAPI }

// Ping: sin llamada real decidida (PROVISIONAL), responde ErrNotConfigured; apagado, ErrDisabled.
func (c *Claude) Ping(ctx context.Context) error {
	if err := c.on.check(ctx); err != nil {
		return err
	}
	return ErrNotConfigured
}

func (c *Claude) ReadTicket(ctx context.Context, _ TicketImage) (TicketReading, error) {
	if err := c.on.check(ctx); err != nil {
		return TicketReading{}, err
	}
	return TicketReading{}, ErrNotConfigured
}

func (c *Claude) SuggestRecipe(ctx context.Context, _ RecipeRequest) (RecipeSuggestion, error) {
	if err := c.on.check(ctx); err != nil {
		return RecipeSuggestion{}, err
	}
	return RecipeSuggestion{}, ErrNotConfigured
}

// Check dice si el proveedor responde: apagado = ErrDisabled sin llamar a nada; sin adaptador = ErrNotConfigured.
func (r *Registry) Check(ctx context.Context, code string) error {
	on, err := r.Enabled(ctx, code)
	if err != nil {
		return err
	}
	if !on {
		return ErrDisabled
	}
	if p, ok := r.providers[code].(Pinger); ok {
		return p.Ping(ctx)
	}
	return ErrNotConfigured
}
