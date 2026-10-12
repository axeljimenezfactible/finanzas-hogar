package main

import (
	"encoding/json"

	"finanzas-hogar/internal/motor_cred/evaluar"
)

func init() { register("s3", evalS3) }

// evalS3 es el Evaluator de S3 (FORMATO.md 1.2 §4.1): input es entradas.datos ya resuelto.
func evalS3(input json.RawMessage, q Query) (out any, code string) {
	return evaluar.Evaluate(input, q.Func, q.Args)
}
