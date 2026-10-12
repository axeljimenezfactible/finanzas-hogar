package main

import (
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/testutil"
)

// Hallazgo H4 de REVISION-backend-core-r1 (S6-03 c, S6-05, C3-03): el registro que arma el binario con areas cubre
// las operaciones de C1 (ninguna responde 501) y el inicio de sesión funciona con él. Solo se consulta el registro:
// no se llama a operaciones que envían avisos.
func TestAreas_SinOperacionesEn501(t *testing.T) {
	s := testutil.NewServer(t, areas...)
	var faltan []string
	for _, rt := range api.Routes {
		if !s.Reg.Has(rt.OperationID) {
			faltan = append(faltan, rt.OperationID)
		}
	}
	if len(faltan) > 0 {
		t.Fatalf("%d de %d operaciones sin registrar (responden 501): %v", len(faltan), len(api.Routes), faltan)
	}
	t.Logf("operaciones registradas: %d de %d", len(s.Reg.Implemented()), len(api.Routes))
	login := map[string]any{"username": testutil.AxelUser, "password": testutil.FakeCredential, "device_label": "Prueba del binario"}
	if res := s.Anonymous(t).Do("POST", "/api/v1/auth/login", login); res.StatusCode != 200 {
		t.Fatalf("login con el registro del binario: %d %s", res.StatusCode, res.Body)
	}
	if res := s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/budget/month", nil); res.StatusCode != 200 {
		t.Fatalf("una operación de área con el registro del binario: %d %s", res.StatusCode, res.Body)
	}
}
