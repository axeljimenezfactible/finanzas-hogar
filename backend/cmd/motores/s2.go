package main

import (
	"encoding/json"

	"finanzas-hogar/internal/motor_fin"
)

// Registro del motor financiero (S2): evalúa cada consulta con motor_fin.Eval.
func init() {
	register("s2", func(input json.RawMessage, q Query) (any, string) {
		out, err := motor_fin.Eval(input, q.Func, q.Args)
		if err != nil {
			return nil, err.Code
		}
		return out, ""
	})
}
