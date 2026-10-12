package comparador

import (
	"bytes"
	"context"
	"encoding/json"

	mc "finanzas-hogar/internal/motor_cred"
)

// EvalComparar atiende la consulta «comparar» de las fixtures y de «motores s3 dump»
// (S3-interfaz.md §9.2): argumentos {timeout_seconds, request}; today y credits salen de datos.
// Devuelve el Result o el código de error de motor_cred.Code (p. ej. «timeout»). El despacho
// (evaluar.Evaluate) lo llama para la función «comparar».
func EvalComparar(in mc.Input, args json.RawMessage) (any, string) {
	out, err := Query(in, args)
	if err != nil {
		return nil, mc.Code(err)
	}
	return out, ""
}

// Query es EvalComparar con el error tipado (ErrInvalid, ErrUnknownField, ErrTimeout, …) en vez
// del código; la usa el despacho de «comparar» de evaluar.Evaluate (S3-interfaz.md §9.2).
func Query(in mc.Input, args json.RawMessage) (Result, error) {
	var a struct {
		TimeoutSeconds *int64          `json:"timeout_seconds"`
		Request        json.RawMessage `json:"request"`
	}
	if err := mc.DecodeStrict(bytes.NewReader(args), &a); err != nil {
		return Result{}, err
	}
	if a.TimeoutSeconds == nil || len(a.Request) == 0 {
		return Result{}, invalid("faltan timeout_seconds o request")
	}
	var req Request
	if err := mc.DecodeStrict(bytes.NewReader(a.Request), &req); err != nil {
		return Result{}, err
	}
	req.Today, req.Credits = in.Today, in.Credits
	d, err := TimeoutFromSetting(*a.TimeoutSeconds)
	if err != nil {
		return Result{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return Compare(ctx, req)
}
