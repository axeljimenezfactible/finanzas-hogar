// Paquete recetas: motor de recetas, comidas y salud (S4). Núcleo de funciones puras:
// sin base de datos, sin reloj, sin red y sin flotantes; «hoy» lo pasa el llamador
// (docs/contratos/motores/S4-interfaz.md §1).
package recetas

import (
	"errors"
	"fmt"
)

// codeError es un error del paquete con su código en snake_case (S4-interfaz §12.2 y §14).
type codeError struct{ code string }

func (e *codeError) Error() string { return e.code }

// Errores del paquete (§14) y códigos genéricos de las fixtures (§12.2).
var (
	ErrPersonWithoutPlan     error = &codeError{"person_without_plan"}
	ErrEatingOutConflict     error = &codeError{"eating_out_conflict"}
	ErrUnknownProduct        error = &codeError{"unknown_product"}
	ErrUnknownPresentation   error = &codeError{"unknown_presentation"}
	ErrUnitMismatch          error = &codeError{"unit_mismatch"}
	ErrProductWithoutChannel error = &codeError{"product_without_channel"}
	ErrKcalMismatch          error = &codeError{"kcal_mismatch"}
	ErrInvalidServings       error = &codeError{"invalid_servings"}
	ErrNotMonday             error = &codeError{"not_monday"}
	ErrOverflow              error = &codeError{"overflow"}
	ErrUnknownField          error = &codeError{"unknown_field"}
	ErrMissingField          error = &codeError{"missing_field"}
	ErrInvalidDate           error = &codeError{"invalid_date"}
	ErrInvalidArgument       error = &codeError{"invalid_argument"}
	// ErrUnconfirmed es el motivo de rechazo de un renglón «propuesto» de tabla.csv en modo production (§9.1).
	ErrUnconfirmed error = &codeError{"unconfirmed"}
)

// errf envuelve un error base con un detalle en español.
func errf(base error, format string, a ...any) error {
	return fmt.Errorf("%w: %s", base, fmt.Sprintf(format, a...))
}

// Code devuelve el código de un error del paquete; cualquier otro error es invalid_argument.
func Code(err error) string {
	var ce *codeError
	if errors.As(err, &ce) {
		return ce.code
	}
	return "invalid_argument"
}

// overflowPanic se lanza dentro del paquete cuando un resultado no cabe en int64;
// las funciones exportadas lo convierten en ErrOverflow con catch.
type overflowPanic struct{}

func catch(err *error) {
	if r := recover(); r != nil {
		if _, ok := r.(overflowPanic); ok {
			*err = ErrOverflow
			return
		}
		panic(r)
	}
}
