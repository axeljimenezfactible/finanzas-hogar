package tablero_test

import (
	"database/sql"
	"fmt"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/creditos"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/api/tablero"
	"finanzas-hogar/internal/testutil"
)

// Pruebas con datos (agente api-tablero, WF-3c). Las cifras esperadas son las de los motores: se comprueban contra
// las operaciones vigentes que usan el mismo motor (getCreditTotals de S3, getDashboardMonthVsPlan de S2) y, donde la
// cuenta es del tablero (RN-30, RN-31), contra el cálculo a mano con ROUND_HALF_UP. Datos falsos de prueba, ninguno
// de la casa. Reloj de testutil: 12-oct-2026 12:00.

type tabFx struct {
	s                              *testutil.Server
	c                              *testutil.Client
	groceries, transport, misc, hl int64
}

func ins(t *testing.T, db *sql.DB, q string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow(q+" RETURNING id", args...).Scan(&id); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return id
}

var nKey int

func expense(t *testing.T, c *testutil.Client, cat, cents int64, date string) {
	t.Helper()
	nKey++
	res := c.Do("POST", "/api/v1/expenses", map[string]any{"category_id": cat, "amount_cents": cents, "expense_date": date,
		"payment_method": "cash", "accrual_months": 1}, "Idempotency-Key", fmt.Sprintf("tablero-prueba-%010d", nKey))
	if res.StatusCode != 201 {
		t.Fatalf("gasto %s: %d %s", date, res.StatusCode, res.Body)
	}
}

// fixture: sueldo de Axel de 10,000.00 por quincena desde jul-2026 (capturado jul–sep, esperado en oct), un extra de
// 3,000.00 el 5-oct, una compra a 3 pagos de 3,750.00 (15-oct, 30-oct y 15-nov; por pago 3,250.00 del hogar y 500.00
// de Karen), planes de 1,000.00 desde jul-2026 en Despensa, Transporte y Varios (Salud sin plan) y gastos de contado.
func fixture(t *testing.T) tabFx {
	t.Helper()
	s := testutil.NewServer(t, tablero.Register, gastos.Register, creditos.Register)
	f := tabFx{s: s, c: s.SessionFor(t, testutil.AxelUser)}
	db, axel, karen := s.DB, s.Fx.AxelPersonID, s.Fx.KarenPersonID
	ins(t, db, `INSERT INTO cost_share_rule (name, remainder_person_id, is_default) VALUES ('Hogar', ?, 1)`, axel)
	f.groceries = ins(t, db, `INSERT INTO category (code, name, kind) VALUES ('groceries', 'Despensa', 'variable')`)
	f.transport = ins(t, db, `INSERT INTO category (code, name, kind) VALUES ('transport', 'Transporte', 'variable')`)
	f.misc = ins(t, db, `INSERT INTO category (code, name, kind) VALUES ('misc', 'Varios', 'variable')`)
	f.hl = ins(t, db, `INSERT INTO category (code, name, kind) VALUES ('health', 'Salud', 'variable')`)
	for _, cat := range []int64{f.groceries, f.transport, f.misc} {
		ins(t, db, `INSERT INTO budget (category_id, amount_cents, valid_from_date) VALUES (?, 100000, '2026-07-01')`, cat)
	}
	pi := ins(t, db, `INSERT INTO planned_item (kind, description, amount_cents, frequency, start_date, person_id) VALUES ('income', 'Sueldo', 1000000, 'fortnightly', '2026-07-15', ?)`, axel)
	for _, d := range []string{"2026-07-15", "2026-07-30", "2026-08-15", "2026-08-30", "2026-09-15", "2026-09-30"} {
		ins(t, db, `INSERT INTO income (person_id, kind, planned_item_id, description, amount_cents, income_date) VALUES (?, 'fixed', ?, 'Sueldo', 1000000, ?)`, axel, pi, d)
	}
	ins(t, db, `INSERT INTO income (person_id, kind, description, amount_cents, income_date) VALUES (?, 'extra', 'Bono', 300000, '2026-10-05')`, axel)
	cr := ins(t, db, `INSERT INTO credit (name, kind, holder_person_id, credit_limit_cents, payment_frequency) VALUES ('Línea de prueba', 'revolving_line', ?, 5000000, 'fortnightly')`, axel)
	pu := ins(t, db, `INSERT INTO credit_purchase (credit_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents, term_count, payment_frequency, first_payment_date, interest_source)
 VALUES (?, 'installment_plan', 'Refrigerador', '2026-10-01', 1125000, 0, 1125000, 3, 'fortnightly', '2026-10-15', 'quoted')`, cr)
	ins(t, db, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 975000, 0, 975000)`, pu, axel)
	ins(t, db, `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, 150000, 0, 150000)`, pu, karen)
	for i, d := range []string{"2026-10-15", "2026-10-30", "2026-11-15"} {
		sp := ins(t, db, `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, ?, ?, 375000, 375000, 0)`, pu, i+1, d)
		ins(t, db, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 325000)`, sp, axel)
		ins(t, db, `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, 50000)`, sp, karen)
	}
	// Gastos de contado (Hogar = todo de Axel). % usado: Despensa 115/120/111 (110.5 sube con ROUND_HALF_UP),
	// Transporte 50/89/0, Varios 150/110/150 (110.4 baja a 110: no pasa de 110).
	for _, g := range []struct {
		cat   int64
		cents int64
		date  string
	}{
		{f.groceries, 115000, "2026-07-10"}, {f.groceries, 120000, "2026-08-10"}, {f.groceries, 110500, "2026-09-10"},
		{f.transport, 50000, "2026-07-11"}, {f.transport, 89000, "2026-08-11"},
		{f.misc, 150000, "2026-07-12"}, {f.misc, 110400, "2026-08-12"}, {f.misc, 150000, "2026-09-12"},
		{f.hl, 20000, "2026-09-13"}, {f.groceries, 50000, "2026-10-05"},
	} {
		expense(t, f.c, g.cat, g.cents, g.date)
	}
	return f
}

func debt(t *testing.T, c *testutil.Client, q string) api.DashboardDebtCapacity {
	t.Helper()
	res := c.Do("GET", "/api/v1/dashboard/debt-capacity"+q, nil)
	var out api.DashboardDebtCapacity
	res.JSON(t, &out)
	if res.StatusCode != 200 {
		t.Fatalf("debt-capacity%s: %d %s", q, res.StatusCode, res.Body)
	}
	return out
}

func val[T any](p *T) any {
	if p == nil {
		return nil
	}
	return *p
}

func TestOp_getDashboardDebtCapacity_ConDatos(t *testing.T) {
	f := fixture(t)
	d := debt(t, f.c, "")
	// Octubre: 2 quincenas esperadas = 20,000.00; pagos del hogar 6,500.00 → 32.5 % → 33 (ROUND_HALF_UP) → warning;
	// con Karen 7,500.00 → 37.5 % → 38; margen = 40 % × 20,000.00 − 6,500.00 = 1,500.00.
	got := []any{d.Month, val(d.IncomeCents), d.IncomeIsExpected, d.ExtrasCents, d.PaymentsCents, d.PaymentsTotalCents, d.PaymentCount,
		val(d.RatioPct), val(d.RatioTotalPct), val(d.Status), val(d.HeadroomCents), d.HasCredits, d.WarningPct, d.LimitPct, d.ScaleMaxPct}
	want := []any{"2026-10", int64(2000000), true, int64(300000), int64(650000), int64(750000), int64(2),
		int64(33), int64(38), "warning", int64(150000), true, int64(30), int64(40), int64(60)}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("octubre:\n got %v\nwant %v", got, want)
	}
	// Misma cifra que getCreditTotals del mes en curso (mismo motor S3).
	var ct api.CreditTotals
	res := f.c.Do("GET", "/api/v1/credits/totals", nil)
	res.JSON(t, &ct)
	if res.StatusCode != 200 || ct.MonthHouseholdCents != d.PaymentsCents || ct.MonthTotalCents != d.PaymentsTotalCents {
		t.Fatalf("getCreditTotals %d: hogar %d total %d; tablero %d %d", res.StatusCode, ct.MonthHouseholdCents, ct.MonthTotalCents, d.PaymentsCents, d.PaymentsTotalCents)
	}
	// Septiembre: quincenas capturadas, sin pagos de crédito: 0 % es cifra real (has_credits true), verde.
	if d := debt(t, f.c, "?period=2026-09"); val(d.IncomeCents) != int64(2000000) || d.IncomeIsExpected || d.PaymentsCents != 0 || val(d.RatioPct) != int64(0) || val(d.Status) != "normal" {
		t.Errorf("septiembre: %+v", d)
	}
	// Noviembre: un pago, 3,250.00 ÷ 20,000.00 = 16.25 % → 16.
	if d := debt(t, f.c, "?period=2026-11"); d.PaymentCount != 1 || d.PaymentsCents != 325000 || val(d.RatioPct) != int64(16) || val(d.Status) != "normal" {
		t.Errorf("noviembre: %+v", d)
	}
	// Junio: sin quincena capturada ni esperada → sin ingreso, porcentaje, estado ni margen (sin datos no es cero).
	if d := debt(t, f.c, "?period=2026-06"); d.IncomeCents != nil || d.RatioPct != nil || d.RatioTotalPct != nil || d.Status != nil || d.HeadroomCents != nil {
		t.Errorf("junio: %+v", d)
	}
	// El color sale de app_setting (K44): 35/45 → verde y escala 70; 20/32 → rojo con margen negativo.
	patch := func(w, l int) {
		t.Helper()
		if res := f.c.Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"debt.capacity_warning_pct": w, "debt.capacity_limit_pct": l}}); res.StatusCode != 200 {
			t.Fatalf("PATCH %d/%d: %d %s", w, l, res.StatusCode, res.Body)
		}
	}
	patch(35, 45)
	if d := debt(t, f.c, ""); val(d.Status) != "normal" || d.WarningPct != 35 || d.LimitPct != 45 || d.ScaleMaxPct != 70 || val(d.HeadroomCents) != int64(250000) {
		t.Errorf("35/45: %+v", d)
	}
	patch(20, 32)
	if d := debt(t, f.c, ""); val(d.Status) != "exceeded" || val(d.HeadroomCents) != int64(-10000) {
		t.Errorf("20/32: %+v", d)
	}
}

func TestScaleMaxPct(t *testing.T) {
	for limit, want := range map[int64]int64{40: 60, 45: 70, 32: 50, 1: 10, 100: 150, 20: 30} {
		if got := tablero.ScaleMaxPctForTest(limit); got != want {
			t.Errorf("limit %d: %d, se esperaba %d", limit, got, want)
		}
	}
}

func TestOp_getDashboardMonthlyOverview_Serie12(t *testing.T) {
	f := fixture(t)
	res := f.c.Do("GET", "/api/v1/dashboard/monthly-overview?from=2025-11&to=2026-10", nil)
	var ov api.DashboardMonthlyOverview
	res.JSON(t, &ov)
	if res.StatusCode != 200 || len(ov.Months) != 12 || ov.Months[0].Month != "2025-11" || ov.Months[11].Month != "2026-10" {
		t.Fatalf("%d: %d meses", res.StatusCode, len(ov.Months))
	}
	for _, m := range ov.Months[:8] { // nov-2025 a jun-2026: sin ingreso ni gasto → sin cifras (sin datos no es cero)
		if m.IncomeTotalCents != nil || m.AccruedCents != nil || m.CashCents != nil || m.DebtRatioPct != nil || m.IsProjection {
			t.Errorf("%s debería venir sin cifras: %+v", m.Month, m)
		}
	}
	jul := ov.Months[8]
	// Julio: 20,000.00 de quincenas capturadas; devengado 1,150.00 + 500.00 + 1,500.00 = 3,150.00; plan 3 × 1,000.00.
	if val(jul.IncomeFortnightsCents) != int64(2000000) || val(jul.IncomeExtrasCents) != int64(0) || jul.IncomeIsExpected ||
		val(jul.AccruedCents) != int64(315000) || val(jul.PlannedCents) != int64(300000) || val(jul.DebtRatioPct) != int64(0) || val(jul.DebtStatus) != "normal" {
		t.Errorf("julio: %+v", jul)
	}
	oct := ov.Months[11]
	if val(oct.IncomeTotalCents) != int64(2300000) || val(oct.IncomeFortnightsCents) != int64(2000000) || val(oct.IncomeExtrasCents) != int64(300000) ||
		!oct.IncomeIsExpected || oct.CreditPaymentsCents != 650000 || oct.CreditPaymentsTotalCents != 750000 || val(oct.DebtRatioPct) != int64(33) || val(oct.DebtStatus) != "warning" {
		t.Errorf("octubre: %+v", oct)
	}
	// Devengado, efectivo y plan del mes en curso = los de getDashboardMonthVsPlan (mismo motor S2).
	var mv api.DashboardMonthVsPlan
	res = f.c.Do("GET", "/api/v1/dashboard/month-vs-plan", nil)
	res.JSON(t, &mv)
	if res.StatusCode != 200 || val(oct.AccruedCents) != mv.ActualCents || fmt.Sprint(val(oct.CashCents)) != fmt.Sprint(val(mv.CashCents)) || fmt.Sprint(val(oct.PlannedCents)) != fmt.Sprint(val(mv.PlannedCents)) {
		t.Errorf("octubre contra month-vs-plan: tablero %v/%v/%v, M-01 %d/%v/%v", val(oct.AccruedCents), val(oct.CashCents), val(oct.PlannedCents), mv.ActualCents, val(mv.CashCents), val(mv.PlannedCents))
	}
	// Proyección: noviembre a la vista con el pago del 15-nov y sin devengado ni efectivo.
	res = f.c.Do("GET", "/api/v1/dashboard/monthly-overview?from=2026-11&to=2027-10", nil)
	ov = api.DashboardMonthlyOverview{}
	res.JSON(t, &ov)
	nov := ov.Months[0]
	if res.StatusCode != 200 || len(ov.Months) != 12 || !nov.IsProjection || nov.AccruedCents != nil || nov.CashCents != nil ||
		nov.CreditPaymentsCents != 325000 || val(nov.DebtRatioPct) != int64(16) || !nov.IncomeIsExpected {
		t.Errorf("noviembre (proyección): %d %+v", res.StatusCode, nov)
	}
	if last := ov.Months[11]; last.Month != "2027-10" || last.CreditPaymentsCents != 0 {
		t.Errorf("oct-2027: %+v", last)
	}
	// Los umbrales del resumen son los vigentes de app_setting.
	f.c.Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"debt.capacity_warning_pct": 35, "debt.capacity_limit_pct": 45}})
	res = f.c.Do("GET", "/api/v1/dashboard/monthly-overview?from=2026-10&to=2026-10", nil)
	ov = api.DashboardMonthlyOverview{}
	res.JSON(t, &ov)
	if ov.WarningPct != 35 || ov.LimitPct != 45 || val(ov.Months[0].DebtStatus) != "normal" {
		t.Errorf("umbrales 35/45: %+v", ov)
	}
	if res := f.c.Raw("GET", "/api/v1/dashboard/monthly-overview?from=2026-13&to=2026-12", nil); res.StatusCode != 400 {
		t.Errorf("mes inválido en from: %d", res.StatusCode)
	}
	if res := f.c.Raw("GET", "/api/v1/dashboard/monthly-overview?from=2026-10", nil); res.StatusCode != 400 {
		t.Errorf("sin to: %d", res.StatusCode)
	}
}

func signals(t *testing.T, c *testutil.Client, q string) (api.BudgetSignals, map[int64]api.BudgetSignalsRowsItem) {
	t.Helper()
	res := c.Do("GET", "/api/v1/budget/signals"+q, nil)
	var out api.BudgetSignals
	res.JSON(t, &out)
	if res.StatusCode != 200 {
		t.Fatalf("signals%s: %d %s", q, res.StatusCode, res.Body)
	}
	by := map[int64]api.BudgetSignalsRowsItem{}
	for _, r := range out.Rows {
		by[r.CategoryID] = r
	}
	return out, by
}

func TestOp_getBudgetSignals_ConDatos(t *testing.T) {
	f := fixture(t)
	out, by := signals(t, f.c, "")
	if len(out.Rows) != 4 {
		t.Fatalf("se esperaban 4 categorías activas: %+v", out.Rows)
	}
	// RN-31 con 3 meses y margen 10: Despensa 115/120/111 → raise; Transporte 50/89/0 (con plan y sin gasto = 0 %) →
	// cut; Varios 150/110/150 → none (110 no pasa de 110); Salud sin plan → none y sin % usado.
	for cat, w := range map[int64]string{f.groceries: "raise", f.transport: "cut", f.misc: "none", f.hl: "none"} {
		if by[cat].Signal != w {
			t.Errorf("categoría %d: %s, se esperaba %s (%+v)", cat, by[cat].Signal, w, by[cat].Months)
		}
	}
	pct := func(cat int64) string {
		s := ""
		for _, m := range by[cat].Months {
			s += fmt.Sprintf("%s=%v/%d ", m.Month, val(m.UsedPct), m.ActualCents)
		}
		return s
	}
	if g := pct(f.groceries); g != "2026-07=115/115000 2026-08=120/120000 2026-09=111/110500 " {
		t.Errorf("Despensa: %s", g)
	}
	if g := pct(f.transport); g != "2026-07=50/50000 2026-08=89/89000 2026-09=0/0 " {
		t.Errorf("Transporte: %s", g)
	}
	if g := pct(f.hl); g != "2026-07=<nil>/0 2026-08=<nil>/0 2026-09=<nil>/20000 " || by[f.hl].Months[0].PlannedCents != nil {
		t.Errorf("Salud: %s", g)
	}
	// Margen 5 (app_setting): Varios 150/110/150 ya pasa de 105 → raise.
	if res := f.c.Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"budget.signal_margin_pct": 5}}); res.StatusCode != 200 {
		t.Fatalf("margen 5: %d %s", res.StatusCode, res.Body)
	}
	if out, by = signals(t, f.c, ""); out.MarginPct != 5 || by[f.misc].Signal != "raise" {
		t.Errorf("margen 5: %+v", by[f.misc])
	}
	// Racha de 2 meses (ago y sep) pidiendo octubre explícito.
	if res := f.c.Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"budget.signal_months": 2}}); res.StatusCode != 200 {
		t.Fatalf("2 meses: %d %s", res.StatusCode, res.Body)
	}
	if out, by = signals(t, f.c, "?period=2026-10"); out.MonthsConsidered != 2 || len(by[f.groceries].Months) != 2 || by[f.groceries].Months[0].Month != "2026-08" {
		t.Errorf("2 meses: %+v", out)
	}
	// Agosto pedido: mira may–jul; mayo y junio sin plan cortan la racha → none en todas.
	if _, by = signals(t, f.c, "?period=2026-08"); by[f.groceries].Signal != "none" || by[f.transport].Signal != "none" {
		t.Errorf("agosto: %+v %+v", by[f.groceries], by[f.transport])
	}
}

// Ajustes del tablero: el cambio queda en audit_log, un 422 no guarda nada y un cuerpo malformado es 400.
func TestOp_updateDashboardSettings_Bitacora(t *testing.T) {
	s := testutil.NewServer(t, tablero.Register)
	c := s.SessionFor(t, testutil.XimenaUser)
	if res := c.Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"dashboard.payment_reminder_count": 8}}); res.StatusCode != 200 {
		t.Fatalf("%d %s", res.StatusCode, res.Body)
	}
	var n int
	if err := s.DB.QueryRow(`SELECT count(*) FROM audit_log WHERE table_name = 'app_setting' AND action = 'update' AND user_account_id IS NOT NULL AND after_json LIKE '%dashboard.payment_reminder_count%'`).Scan(&n); err != nil || n == 0 {
		t.Errorf("audit_log del cambio: n=%d err=%v", n, err)
	}
	// 422 por amarillo ≥ tope: ni el tope ni el amarillo se guardan.
	if res := c.Do("PATCH", "/api/v1/settings/dashboard", map[string]any{"parameters": map[string]any{"debt.capacity_warning_pct": 60, "debt.capacity_limit_pct": 50}}); res.StatusCode != 422 {
		t.Fatalf("se esperaba 422: %d", res.StatusCode)
	}
	var out api.DashboardSettings
	c.Do("GET", "/api/v1/settings/dashboard", nil).JSON(t, &out)
	if out.Parameters.DebtCapacityWarningPct.Value.Value != 30 || out.Parameters.DebtCapacityLimitPct.Value.Value != 40 || out.Parameters.DashboardPaymentReminderCount.Value.Value != 8 {
		t.Errorf("tras el 422: %+v", out.Parameters)
	}
	for name, body := range map[string]any{
		"reset repetido":    map[string]any{"reset_parameters": []string{"budget.signal_months", "budget.signal_months"}},
		"clave desconocida": map[string]any{"parameters": map[string]any{"otra.clave": 1}},
		"tipo texto":        map[string]any{"parameters": map[string]any{"budget.signal_months": "4"}},
	} {
		// Formato (C1 §3 y CONVENCIONES-backend §4): repetido, clave fuera del esquema o tipo incorrecto → 400.
		if res := c.Raw("PATCH", "/api/v1/settings/dashboard", body); res.StatusCode != 400 {
			t.Errorf("%s: %d %s", name, res.StatusCode, res.Body)
		}
	}
}

// RN-31: solo cuentan los meses cerrados; el mes en curso (oct-2026 con el reloj de testutil) y los de proyección no
// se miran. PROVISIONAL, bloqueo S6-tablero-senal-periodo, opción (a); hallazgo M-1 de REVISION-tablero-backend-r1.
func TestOp_getBudgetSignals_PeriodoFuturo(t *testing.T) {
	f := fixture(t)
	months := func(r api.BudgetSignalsRowsItem) string {
		s := ""
		for _, m := range r.Months {
			s += m.Month + " "
		}
		return s
	}
	// Mes siguiente: la ventana es ago–oct, pero octubre está en curso. Sin el corte, Transporte (89 %, 0 % y el 0 % de
	// octubre, que aún no termina) daría cut; con 2 meses cerrados de 3 la racha no se completa → none en las 4 filas.
	out, by := signals(t, f.c, "?period=2026-11")
	if out.Month != "2026-11" || out.MonthsConsidered != 3 || len(out.Rows) != 4 {
		t.Fatalf("noviembre: %+v", out)
	}
	for cat, r := range by {
		if r.Signal != "none" || months(r) != "2026-08 2026-09 " {
			t.Errorf("noviembre, categoría %d: %s con meses %q (se esperaba none con 2026-08 y 2026-09)", cat, r.Signal, months(r))
		}
	}
	// Period futuro: ene-2027 mira oct (en curso), nov y dic (proyección); feb-2027 solo proyección. Antes daban cut con
	// gasto que todavía no ocurre; ahora ningún mes cuenta: filas sin meses y none.
	for _, p := range []string{"2027-01", "2027-02"} {
		out, by := signals(t, f.c, "?period="+p)
		if out.Month != p || len(out.Rows) != 4 {
			t.Fatalf("%s: %+v", p, out)
		}
		for cat, r := range by {
			if r.Signal != "none" || len(r.Months) != 0 {
				t.Errorf("%s, categoría %d: %s con meses %q (se esperaba none sin meses)", p, cat, r.Signal, months(r))
			}
		}
	}
	// El mes en curso explícito no cambia: mira jul–sep, todos cerrados (Despensa raise y Transporte cut).
	if _, by = signals(t, f.c, "?period=2026-10"); by[f.groceries].Signal != "raise" || by[f.transport].Signal != "cut" || months(by[f.transport]) != "2026-07 2026-08 2026-09 " {
		t.Errorf("octubre: %+v %+v", by[f.groceries], by[f.transport])
	}
}

// RN-30 «de todas las líneas»: una línea archivada sigue contando en los pagos de crédito (S3), igual que en el
// efectivo de S2 (dominio.LoadFinEngine carga todos los scheduled_payment). PROVISIONAL, bloqueo
// S6-tablero-lineas-archivadas, opción (a) y (a1); hallazgo M-2 de REVISION-tablero-backend-r1 y r2.
// Lo que vale con cualquier opción del bloqueo corre siempre. Lo propio de la opción (a) va en el subtest
// «opcion_a_PROVISIONAL», en SKIP mientras no se aplique la solicitud a backend/internal/api/creditos/tablero.go
// (fuera del alcance de api-tablero; ver el punto 4 del bloqueo).
func TestOp_Dashboard_LineaArchivada(t *testing.T) {
	f := fixture(t)
	oct := func() api.DashboardMonthlyOverviewMonthsItem {
		t.Helper()
		res := f.c.Do("GET", "/api/v1/dashboard/monthly-overview?from=2026-10&to=2026-10", nil)
		var ov api.DashboardMonthlyOverview
		res.JSON(t, &ov)
		if res.StatusCode != 200 || len(ov.Months) != 1 {
			t.Fatalf("monthly-overview: %d %s", res.StatusCode, res.Body)
		}
		return ov.Months[0]
	}
	before := oct()
	if _, err := f.s.DB.Exec(`UPDATE credit SET is_active = 0`); err != nil {
		t.Fatal(err)
	}
	after := oct()
	// Con cualquier opción. S2: el efectivo de octubre no cambia al archivar (7,000.00 = 6,500.00 de pagos del
	// hogar + 500.00 de contado del 5-oct).
	cash, ok := val(after.CashCents).(int64)
	if !ok || cash != 700000 || fmt.Sprint(val(before.CashCents)) != fmt.Sprint(cash) {
		t.Fatalf("efectivo S2: antes %v, después %v", val(before.CashCents), val(after.CashCents))
	}
	// Con cualquier opción. G-1/G-6 y el barómetro dan la misma cifra del mes (los dos salen de
	// creditos.PaymentsByMonth), y con has_credits false los pagos y el porcentaje son 0, cifra real (C1).
	byPeriod := map[string]api.DashboardDebtCapacity{}
	for _, q := range []string{"", "?period=2026-11", "?period=2026-12"} {
		d := debt(t, f.c, q)
		byPeriod[q] = d
		if r := val(d.RatioPct); !d.HasCredits && (d.PaymentsCents != 0 || d.PaymentsTotalCents != 0 || d.PaymentCount != 0 || (r != nil && r != int64(0))) {
			t.Errorf("debt-capacity%s: has_credits false con pagos o porcentaje: %+v", q, d)
		}
	}
	if d := byPeriod[""]; d.PaymentsCents != after.CreditPaymentsCents || d.PaymentsTotalCents != after.CreditPaymentsTotalCents {
		t.Errorf("octubre: barómetro %d · %d, serie %d · %d", d.PaymentsCents, d.PaymentsTotalCents, after.CreditPaymentsCents, after.CreditPaymentsTotalCents)
	}
	t.Run("opcion_a_PROVISIONAL", func(t *testing.T) {
		// S3: los pagos tampoco cambian, y son la parte de crédito del efectivo de S2 (efectivo − contado del mes).
		if after.CreditPaymentsCents != before.CreditPaymentsCents || after.CreditPaymentsCents != cash-50000 ||
			after.CreditPaymentsTotalCents != 750000 || val(after.DebtRatioPct) != int64(33) || val(after.DebtStatus) != "warning" {
			t.Errorf("octubre con la línea archivada: pagos %d (antes %d; parte de crédito del efectivo S2 %d), total %d, %v %%, %v",
				after.CreditPaymentsCents, before.CreditPaymentsCents, cash-50000, after.CreditPaymentsTotalCents, val(after.DebtRatioPct), val(after.DebtStatus))
		}
		// Barómetro: octubre y noviembre siguen con sus pagos y has_credits true (hay pagos de la línea archivada en el
		// mes, opción a1); diciembre, sin líneas activas ni pagos, es 0 % cifra real con has_credits false.
		if d := byPeriod[""]; d.PaymentsCents != 650000 || d.PaymentCount != 2 || val(d.RatioPct) != int64(33) || !d.HasCredits {
			t.Errorf("debt-capacity octubre: %+v", d)
		}
		if d := byPeriod["?period=2026-11"]; d.PaymentsCents != 325000 || val(d.RatioPct) != int64(16) || !d.HasCredits {
			t.Errorf("debt-capacity noviembre: %+v", d)
		}
		if d := byPeriod["?period=2026-12"]; d.PaymentsCents != 0 || val(d.RatioPct) != int64(0) || d.HasCredits {
			t.Errorf("debt-capacity diciembre: %+v", d)
		}
	})
}
