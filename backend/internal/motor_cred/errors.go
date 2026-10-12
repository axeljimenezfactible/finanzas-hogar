// Package motor_cred es el núcleo del motor de créditos (S3): compras a crédito, series de
// pagos, disponible A/B, calendario unificado, deuda viva, tarjeta, préstamo y crédito de
// tienda. Es puro: no lee reloj, red ni base de datos (S3-interfaz.md §1).
package motor_cred

import (
	"errors"
	"fmt"
)

// Errores con tipo (S3-interfaz.md §1.7). Nunca se devuelve un número inventado en su lugar.
var (
	ErrInvalid                = errors.New("entrada inválida")
	ErrUnknownField           = fmt.Errorf("%w: campo desconocido", ErrInvalid)
	ErrNoRate                 = errors.New("sin tasa")
	ErrMissingCapture         = errors.New("falta capturar")
	ErrPromotionNotApplicable = errors.New("promoción no aplicable")
	ErrUnsupported            = errors.New("sin regla")
	ErrTimeout                = errors.New("tiempo máximo agotado")
)

// invalid arma un ErrInvalid con detalle.
func invalid(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, a...))
}

// ErrorResponse traduce un error del motor a la respuesta de C1 (dato puro, sin el paquete HTTP; §8.7).
type ErrorResponse struct {
	Err    error
	Status int
	Code   string
}

// ErrorResponses trae solo las dos respuestas que C1 fija para el motor.
var ErrorResponses = []ErrorResponse{
	{ErrInvalid, 422, "validation_failed"},
	{ErrTimeout, 503, "unavailable"},
}

// ResponseFor devuelve el primer renglón de ErrorResponses que casa con errors.Is.
func ResponseFor(err error) (status int, code string, ok bool) {
	for _, r := range ErrorResponses {
		if errors.Is(err, r.Err) {
			return r.Status, r.Code, true
		}
	}
	return 0, "", false
}

// Code es el código de error de las salidas de fixtures (S3-interfaz.md §9.1).
func Code(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, ErrUnknownField):
		return "unknown_field"
	case errors.Is(err, ErrInvalid):
		return "invalid"
	case errors.Is(err, ErrNoRate):
		return "no_rate"
	case errors.Is(err, ErrMissingCapture):
		return "missing_capture"
	case errors.Is(err, ErrPromotionNotApplicable):
		return "promotion_not_applicable"
	case errors.Is(err, ErrUnsupported):
		return "unsupported"
	case errors.Is(err, ErrTimeout):
		return "timeout"
	}
	return "invalid"
}
