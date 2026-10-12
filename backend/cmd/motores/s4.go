package main

import (
	"encoding/json"

	"finanzas-hogar/internal/recetas"
)

// Registro del motor de recetas, comidas y salud (S4): evalúa cada consulta con recetas.Eval
// (docs/contratos/motores/S4-interfaz.md §12).
func init() {
	register("s4", func(input json.RawMessage, q Query) (any, string) {
		out, err := recetas.Eval(input, q.Func, q.Args)
		if err != nil {
			return nil, recetas.Code(err)
		}
		return out, ""
	})
}
