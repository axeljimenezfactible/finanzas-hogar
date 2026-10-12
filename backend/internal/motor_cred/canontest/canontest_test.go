package canontest

import (
	"encoding/json"
	"testing"
)

func TestCanontest(t *testing.T) {
	m, d := Load(t, "oraculo/model.json"), Load(t, "datos/semilla/data.json")
	in := Seed(t, m, d)
	if len(in.Credits) != 1 || len(in.Credits[0].Purchases) != Len(t, m, "despensas") || len(in.Cash) == 0 {
		t.Errorf("semilla incompleta: %+v", in)
	}
	si, ipm := Sensitivity(t, m)
	if len(si.Amounts) != Len(t, m, "sensibilidad") || len(ipm) != Len(t, m, "interes_por_monto") {
		t.Errorf("insumo incompleto")
	}
	for _, p := range []string{"x", "despensas[x]", "despensas[99]", "parametros.tasa.x", "despensas.x"} {
		if _, err := Get(m, p); err == nil {
			t.Errorf("ruta %q aceptada", p)
		}
	}
	for _, v := range []any{"texto", json.Number("1.005"), json.Number("x")} {
		if _, err := CentsOf(v); err == nil {
			t.Errorf("%v aceptado como dinero", v)
		}
	}
	if JSON(t, map[string]int{"a": 1}) == nil {
		t.Errorf("JSON vacío")
	}
}
