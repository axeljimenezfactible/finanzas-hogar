package tablero_test

import (
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/tablero"
	"finanzas-hogar/internal/testutil"
)

// Una prueba por operación (CONVENCIONES-backend §12, S6-03): 2xx y los 4xx propios. Do valida cada respuesta
// contra openapi.yaml (C1 v1.5). El reloj de testutil está en 12-oct-2026 12:00.

func errField(t *testing.T, res *testutil.Response) string {
	t.Helper()
	var e api.ErrorBody
	res.JSON(t, &e)
	if len(e.Fields) == 0 {
		return e.Code
	}
	return e.Code + " " + e.Fields[0].Field
}

func TestOp_getDashboardDebtCapacity(t *testing.T) {
	s := testutil.NewServer(t, tablero.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	res := c.Do("GET", "/api/v1/dashboard/debt-capacity", nil)
	var out api.DashboardDebtCapacity
	res.JSON(t, &out)
	if res.StatusCode != 200 || out.Month != "2026-10" || out.WarningPct != 30 || out.LimitPct != 40 || out.ScaleMaxPct != 60 {
		t.Fatalf("%d %+v", res.StatusCode, out)
	}
	// Sin ingreso no hay porcentaje, estado ni margen (sin datos no es cero, RN-30).
	if out.IncomeCents == nil && (out.RatioPct != nil || out.Status != nil || out.HeadroomCents != nil) {
		t.Errorf("sin ingreso no debe haber porcentaje, estado ni margen: %+v", out)
	}
	if res := c.Do("GET", "/api/v1/dashboard/debt-capacity?period=2026-11", nil); res.StatusCode != 200 {
		t.Errorf("mes explícito: %d %s", res.StatusCode, res.Body)
	}
	if res := c.Do("GET", "/api/v1/dashboard/debt-capacity?period=2026-10-q1", nil); res.StatusCode != 422 || errField(t, res) != "validation_failed period" {
		t.Errorf("quincena: se esperaba 422 en period: %d %s", res.StatusCode, res.Body)
	}
	if res := c.Do("GET", "/api/v1/dashboard/debt-capacity?period=oct", nil); res.StatusCode != 400 {
		t.Errorf("periodo malformado: %d", res.StatusCode)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/dashboard/debt-capacity", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

func TestDebtStatus_Fronteras(t *testing.T) {
	// RN-30: fronteras inclusivas sobre el entero (30 verde, 31 amarillo, 40 amarillo, 41 rojo).
	for _, c := range []struct {
		ratio int64
		want  string
	}{{0, "normal"}, {30, "normal"}, {31, "warning"}, {40, "warning"}, {41, "exceeded"}} {
		if got := tablero.DebtStatusForTest(c.ratio, 30, 40); got != c.want {
			t.Errorf("%d: %s, se esperaba %s", c.ratio, got, c.want)
		}
	}
}

func TestOp_getDashboardMonthlyOverview(t *testing.T) {
	s := testutil.NewServer(t, tablero.Register)
	c := s.SessionFor(t, testutil.XimenaUser)
	res := c.Do("GET", "/api/v1/dashboard/monthly-overview?from=2026-09&to=2026-12", nil)
	var out api.DashboardMonthlyOverview
	res.JSON(t, &out)
	if res.StatusCode != 200 || len(out.Months) != 4 || out.WarningPct != 30 || out.LimitPct != 40 {
		t.Fatalf("%d %+v", res.StatusCode, out)
	}
	for i, want := range []struct {
		month string
		proj  bool
	}{{"2026-09", false}, {"2026-10", false}, {"2026-11", true}, {"2026-12", true}} {
		m := out.Months[i]
		if m.Month != want.month || m.IsProjection != want.proj {
			t.Errorf("mes %d: %s proyección=%v, se esperaba %s %v", i, m.Month, m.IsProjection, want.month, want.proj)
		}
		if m.IsProjection && (m.AccruedCents != nil || m.CashCents != nil) {
			t.Errorf("%s: un mes de proyección no trae devengado ni efectivo", m.Month)
		}
	}
	for q, field := range map[string]string{"from=2026-12&to=2026-09": "from", "from=2025-01&to=2027-01": "to"} {
		if res := c.Do("GET", "/api/v1/dashboard/monthly-overview?"+q, nil); res.StatusCode != 422 || errField(t, res) != "validation_failed "+field {
			t.Errorf("%s: se esperaba 422 en %s: %d %s", q, field, res.StatusCode, res.Body)
		}
	}
	if res := c.Do("GET", "/api/v1/dashboard/monthly-overview?from=2025-01&to=2026-12", nil); res.StatusCode != 200 {
		t.Errorf("24 meses es el máximo permitido: %d %s", res.StatusCode, res.Body)
	}
	if res := c.Raw("GET", "/api/v1/dashboard/monthly-overview?to=2026-12", nil); res.StatusCode != 400 {
		t.Errorf("sin from: %d", res.StatusCode)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/dashboard/monthly-overview?from=2026-10&to=2026-10", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

func TestOp_getBudgetSignals(t *testing.T) {
	s := testutil.NewServer(t, tablero.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	res := c.Do("GET", "/api/v1/budget/signals", nil)
	var out api.BudgetSignals
	res.JSON(t, &out)
	if res.StatusCode != 200 || out.Month != "2026-10" || out.MonthsConsidered != 3 || out.MarginPct != 10 {
		t.Fatalf("%d %+v", res.StatusCode, out)
	}
	for _, r := range out.Rows {
		if len(r.Months) != 3 || r.Months[0].Month != "2026-07" || r.Months[2].Month != "2026-09" {
			t.Errorf("categoría %d: meses mirados %+v (se esperaban jul a sep, sin el mes pedido)", r.CategoryID, r.Months)
		}
	}
	if res := c.Do("GET", "/api/v1/budget/signals?period=2026-W41", nil); res.StatusCode != 422 || errField(t, res) != "validation_failed period" {
		t.Errorf("semana: se esperaba 422: %d %s", res.StatusCode, res.Body)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/budget/signals", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

func TestOp_getDashboardSettings(t *testing.T) {
	s := testutil.NewServer(t, tablero.Register)
	res := s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/settings/dashboard", nil)
	var out api.DashboardSettings
	res.JSON(t, &out)
	p := out.Parameters
	got := []int64{p.DebtCapacityWarningPct.Value.Value, p.DebtCapacityLimitPct.Value.Value, p.DashboardPaymentReminderCount.Value.Value,
		p.BudgetSignalMonths.Value.Value, p.BudgetSignalMarginPct.Value.Value}
	if res.StatusCode != 200 || got[0] != 30 || got[1] != 40 || got[2] != 5 || got[3] != 3 || got[4] != 10 || !p.DebtCapacityLimitPct.IsInitial {
		t.Fatalf("%d %v %+v", res.StatusCode, got, p)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/settings/dashboard", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

func TestOp_updateDashboardSettings(t *testing.T) {
	s := testutil.NewServer(t, tablero.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	patch := func(body map[string]any) *testutil.Response {
		return c.Do("PATCH", "/api/v1/settings/dashboard", body)
	}
	res := patch(map[string]any{"parameters": map[string]any{"debt.capacity_limit_pct": 45, "budget.signal_months": 4}})
	var out api.DashboardSettings
	res.JSON(t, &out)
	if res.StatusCode != 200 || out.Parameters.DebtCapacityLimitPct.Value.Value != 45 || out.Parameters.BudgetSignalMonths.Value.Value != 4 || out.Parameters.DebtCapacityLimitPct.IsInitial {
		t.Fatalf("%d %+v", res.StatusCode, out)
	}
	// amarillo ≥ tope: 422 y no se guarda nada
	if res := patch(map[string]any{"parameters": map[string]any{"debt.capacity_warning_pct": 45}}); res.StatusCode != 422 || errField(t, res) != "validation_failed parameters.debt.capacity_warning_pct" {
		t.Errorf("amarillo igual al tope: %d %s", res.StatusCode, res.Body)
	}
	if res := patch(map[string]any{"parameters": map[string]any{"debt.capacity_warning_pct": 50, "debt.capacity_limit_pct": 40}}); res.StatusCode != 422 {
		t.Errorf("amarillo mayor que el tope: %d %s", res.StatusCode, res.Body)
	}
	// la misma clave en parameters y en reset_parameters: 422
	if res := patch(map[string]any{"parameters": map[string]any{"budget.signal_months": 5}, "reset_parameters": []string{"budget.signal_months"}}); res.StatusCode != 422 {
		t.Errorf("cambiar y restablecer a la vez: %d %s", res.StatusCode, res.Body)
	}
	// fuera de rango (disparadores de la 00005): 422
	if res := c.Raw("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"budget.signal_months": 13}}); res.StatusCode != 422 {
		t.Errorf("fuera de rango: %d %s", res.StatusCode, res.Body)
	}
	// restablecer vuelve al valor inicial
	res = patch(map[string]any{"reset_parameters": []string{"debt.capacity_limit_pct", "budget.signal_months"}})
	res.JSON(t, &out)
	if res.StatusCode != 200 || out.Parameters.DebtCapacityLimitPct.Value.Value != 40 || out.Parameters.BudgetSignalMonths.Value.Value != 3 {
		t.Errorf("restablecer: %d %+v", res.StatusCode, out)
	}
	if res := c.Raw("PATCH", "/api/v1/settings/dashboard", map[string]any{}); res.StatusCode != 400 {
		t.Errorf("sin cambios: %d", res.StatusCode)
	}
	if s.Anonymous(t).Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"reset_parameters": []string{"budget.signal_months"}}).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}
