package semilla_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
	"finanzas-hogar/internal/semilla"
	"finanzas-hogar/internal/testutil"
)

// Pruebas de la semilla de DESARROLLO (desarrollo.go): idempotencia, montos = fixtures del oráculo S2 y cotejo de las
// cifras de Inicio y de P1–P6 con oraculo/model.json (C3-09).

func devOpts(t *testing.T, fixed bool) semilla.Options {
	o := opts(t)
	o.DevFixedExpense = fixed
	return o
}

// sembrada devuelve una base con la semilla real y el overlay de desarrollo.
func sembrada(t *testing.T, fixed bool) *sql.DB {
	t.Helper()
	db := emptyDB(t)
	if _, err := semilla.RunWith(context.Background(), db, opts(t)); err != nil {
		t.Fatal(err)
	}
	if _, err := semilla.RunDevWith(context.Background(), db, devOpts(t, fixed)); err != nil {
		t.Fatal(err)
	}
	return db
}

func devCounts(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	m := map[string]int{}
	for _, tb := range []string{"category", "budget", "planned_item", "income", "savings_goal", "savings_movement", "price", "meal_plan_entry", "audit_log"} {
		m[tb] = one[int](t, db, "SELECT count(*) FROM "+tb)
	}
	return m
}

func TestDesarrolloIdempotente(t *testing.T) {
	db := sembrada(t, true)
	want := map[string]int{"category": 4, "budget": 1, "planned_item": 2, "income": 0, "savings_goal": 1, "savings_movement": 1, "price": 72, "meal_plan_entry": 12}
	before := devCounts(t, db)
	for k, w := range want {
		if before[k] != w {
			t.Errorf("%s: %d, se esperaban %d", k, before[k], w)
		}
	}
	rep, err := semilla.RunDevWith(context.Background(), db, devOpts(t, true))
	if err != nil {
		t.Fatal(err)
	}
	if rep.TotalCreated() != 0 {
		t.Errorf("la segunda carga creó %d renglones: %s", rep.TotalCreated(), rep)
	}
	after := devCounts(t, db)
	for k := range before {
		if before[k] != after[k] {
			t.Errorf("%s cambió al sembrar dos veces: %d → %d", k, before[k], after[k])
		}
	}
	// Todo lo del overlay firma audit_log con su propio request_id y deja una marca de desarrollo en el nombre o las notas.
	if n := one[int](t, db, `SELECT count(*) FROM audit_log WHERE request_id = ?`, semilla.RequestIDDev); n == 0 {
		t.Error("el overlay no escribió audit_log con request_id semilla-desarrollo")
	}
	if n := one[int](t, db, `SELECT count(*) FROM planned_item WHERE description NOT LIKE '%(desarrollo)'`); n != 0 {
		t.Errorf("%d partidas sin la marca (desarrollo)", n)
	}
	// La semilla real no se tocó: el overlay no agrega gastos (la real trae 6: la despensa del 10-oct y los 5 proyectados).
	if n := one[int](t, db, `SELECT count(*) FROM expense`); n != 6 {
		t.Errorf("expense: %d (la semilla real trae 6)", n)
	}
}

func TestDesarrolloExigeLaSemillaReal(t *testing.T) {
	_, err := semilla.RunDevWith(context.Background(), emptyDB(t), devOpts(t, false))
	if err == nil || !strings.Contains(err.Error(), "corre primero") {
		t.Fatalf("sin semilla real debe pedir «finanzas semilla»: %v", err)
	}
}

func TestDesarrolloSinFixturesNoInventa(t *testing.T) {
	db := emptyDB(t)
	if _, err := semilla.RunWith(context.Background(), db, opts(t)); err != nil {
		t.Fatal(err)
	}
	o := devOpts(t, false)
	o.Dir = miniSeed(t) // datos/tmp/<x>: sin oraculo/ alrededor
	_, err := semilla.RunDevWith(context.Background(), db, o)
	if err == nil || !strings.Contains(err.Error(), "oraculo/fixtures/s2") {
		t.Fatalf("sin fixtures no debe cargar nada: %v", err)
	}
	if n := one[int](t, db, `SELECT count(*) FROM budget`); n != 0 {
		t.Errorf("budget: %d", n)
	}
}

func TestDesarrolloCandado(t *testing.T) {
	db := emptyDB(t) // semilla.db: sin «dev» en el nombre
	t.Setenv("FINANZAS_ENTORNO", "")
	t.Setenv("FINANZAS_SEMILLA_DEV", "")
	if err := semilla.RunDev(context.Background(), db, seedDir(t), false); err == nil || !strings.Contains(err.Error(), "base de desarrollo") {
		t.Errorf("una base sin «dev» debe negarse: %v", err)
	}
	t.Setenv("FINANZAS_ENTORNO", "produccion")
	t.Setenv("FINANZAS_SEMILLA_DEV", "1")
	if err := semilla.RunDev(context.Background(), db, seedDir(t), false); err == nil || !strings.Contains(err.Error(), "producción") {
		t.Errorf("FINANZAS_ENTORNO=produccion siempre se niega: %v", err)
	}
	// Una base con «dev» en el nombre pasa el candado (y falla después solo porque no hay semilla real).
	t.Setenv("FINANZAS_ENTORNO", "")
	t.Setenv("FINANZAS_SEMILLA_DEV", "")
	dev, err := dominio.Open(filepath.Join(testutil.TempDir(t), "finanzas-dev.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer dev.Close()
	if _, err := dominio.Migrate(context.Background(), dev); err != nil {
		t.Fatal(err)
	}
	if err := semilla.RunDev(context.Background(), dev, seedDir(t), false); err == nil || !strings.Contains(err.Error(), "corre primero") {
		t.Errorf("finanzas-dev.db debe pasar el candado y pedir la semilla real: %v", err)
	}
}

// fixture lee una fixture del oráculo como JSON genérico (independiente de los tipos de desarrollo.go).
func fixture(t *testing.T, name string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(testutil.Root(t), "oraculo", "fixtures", "s2", name))
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func parche(t *testing.T, name, key string, id float64, idKey string) map[string]any {
	t.Helper()
	p := fixture(t, name)["entradas"].(map[string]any)["datos"].(map[string]any)["parche"].(map[string]any)
	for _, x := range p[key].([]any) {
		if m := x.(map[string]any); m[idKey] == id {
			return m
		}
	}
	t.Fatalf("%s: %s %v no existe", name, key, id)
	return nil
}

func TestDesarrolloMontosSalenDeLasFixtures(t *testing.T) {
	db := sembrada(t, true)
	sueldo := parche(t, "s2-proyeccion-csv.json", "planned_items", 1, "id")
	renta := parche(t, "s2-proyeccion-csv.json", "planned_items", 2, "id")
	for _, c := range []struct {
		name string
		fx   map[string]any
		desc string
	}{{"sueldo", sueldo, "Sueldo (desarrollo)"}, {"renta", renta, "Renta (desarrollo)"}} {
		var amt float64
		var start, freq string
		var second sql.NullFloat64
		if err := db.QueryRow(`SELECT amount_cents, second_half_amount_cents, start_date, frequency FROM planned_item WHERE description = ?`, c.desc).Scan(&amt, &second, &start, &freq); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if f := c.fx["amount_cents"].(float64); amt != f {
			t.Errorf("%s: monto %v, la fixture trae %v", c.name, amt, f)
		}
		if f, ok := c.fx["second_half_amount_cents"].(float64); ok != second.Valid || (ok && second.Float64 != f) {
			t.Errorf("%s: segunda quincena %v, la fixture trae %v", c.name, second, c.fx["second_half_amount_cents"])
		}
		if start != c.fx["start_date"] || freq != c.fx["frequency"] {
			t.Errorf("%s: %s/%s contra %v/%v", c.name, start, freq, c.fx["start_date"], c.fx["frequency"])
		}
	}
	bud := parche(t, "s2-presupuesto-scope.json", "budgets", 1, "id")
	if got := one[float64](t, db, `SELECT amount_cents FROM budget b JOIN category c ON c.id = b.category_id WHERE c.code = 'groceries' AND b.valid_from_date = ? AND b.valid_to_date IS NULL`, bud["valid_from"]); got != bud["amount_cents"].(float64) {
		t.Errorf("presupuesto Despensa: %v, la fixture trae %v", got, bud["amount_cents"])
	}
	g := parche(t, "s2-metas.json", "savings_goals", 2, "id")
	var target, monthly float64
	var tdate string
	if err := db.QueryRow(`SELECT target_cents, planned_monthly_cents, target_date FROM savings_goal`).Scan(&target, &monthly, &tdate); err != nil {
		t.Fatal(err)
	}
	if target != g["target_cents"].(float64) || monthly != g["planned_monthly_cents"].(float64) || tdate != g["target_date"] {
		t.Errorf("meta: %v/%v/%s contra %v", target, monthly, tdate, g)
	}
	mov := parche(t, "s2-metas.json", "savings_movements", 2, "goal_id")
	if got := one[float64](t, db, `SELECT amount_cents FROM savings_movement WHERE kind = 'deposit'`); got != mov["amount_cents"].(float64) {
		t.Errorf("aportación: %v, la fixture trae %v", got, mov["amount_cents"])
	}
}

// modelCents lee un número de pesos de model.json (json.Number) y lo vuelve centavos exactos.
func modelCents(t *testing.T, v any) int64 {
	t.Helper()
	r, ok := new(big.Rat).SetString(string(v.(json.Number)))
	if !ok {
		t.Fatalf("número inválido %v", v)
	}
	r.Mul(r, big.NewRat(100, 1))
	if !r.IsInt() {
		t.Fatalf("%v no son centavos exactos", v)
	}
	return r.Num().Int64()
}

func readModelJSON(t *testing.T) map[string]any {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(testutil.Root(t), "oraculo", "model.json"))
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		t.Fatal(err)
	}
	return m
}

// engine arma la entrada del motor S2 con «hoy» = 12-oct-2026, igual que las tarjetas de Inicio.
func engine(t *testing.T, db *sql.DB) (*dominio.FinEngine, mf.Input) {
	t.Helper()
	fe, err := dominio.LoadFinEngine(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	today, err := dominio.ParseDate("2026-10-12")
	if err != nil {
		t.Fatal(err)
	}
	in, err := fe.At(today, nil)
	if err != nil {
		t.Fatal(err)
	}
	return fe, in
}

// meses devuelve devengado y efectivo (centavos) de los 12 meses de la proyección desde oct-2026, como getAccrualVsCash.
func meses(t *testing.T, in mf.Input) [12][2]int64 {
	t.Helper()
	pr, err := mf.Projection(in, mf.NewDate(2026, 10, 1).MonthOf())
	if err != nil {
		t.Fatal(err)
	}
	var out [12][2]int64
	for i := range out {
		out[i] = [2]int64{int64(pr.Months[i].Accrued.TotalCents), int64(pr.Months[i].Cash.TotalCents)}
	}
	return out
}

// modelPeriodos son devengado y efectivo de model.json periodos[0..5] (P1–P6).
func modelPeriodos(t *testing.T) (want [6][2]int64) {
	t.Helper()
	for i, x := range readModelJSON(t)["periodos"].([]any) {
		p := x.(map[string]any)
		want[i] = [2]int64{modelCents(t, p["devengado"].(map[string]any)["total"]), modelCents(t, p["efectivo"].(map[string]any)["total"])}
	}
	return want
}

// soloReal es la base con la semilla real y SIN el overlay (la referencia de «no se movió nada»).
func soloReal(t *testing.T) *sql.DB {
	t.Helper()
	db := emptyDB(t)
	if _, err := semilla.RunWith(context.Background(), db, opts(t)); err != nil {
		t.Fatal(err)
	}
	return db
}

// C3-09: con el overlay por omisión, P1–P6 y la deuda viva siguen siendo las cifras de model.json y las tarjetas de
// Inicio traen las cifras de las fixtures.
func TestDesarrolloCotejoConElOraculo(t *testing.T) {
	db := sembrada(t, false)
	fe, in := engine(t, db)
	_, real := engine(t, soloReal(t))
	got, base, want := meses(t, in), meses(t, real), modelPeriodos(t)
	// 1. El overlay por omisión no mueve ni un centavo de devengado y efectivo de los 12 meses (getAccrualVsCash).
	if got != base {
		t.Errorf("el overlay movió getAccrualVsCash:\n con overlay %v\n solo real   %v", got, base)
	}
	// 2. P1 a P6 son las cifras de model.json (R-45…R-50): desde WF-3b la semilla real trae los gastos proyectados 2 a 6
	//    (bloqueo WF3a-bis-semilla-dev, punto 4, opción a), así que ya no hay hallazgo en P3 a P6.
	for i := 0; i < 6; i++ {
		if got[i] != want[i] {
			t.Errorf("P%d = %v; model.json = %v", i+1, got[i], want[i])
		}
	}
	h := readModelJSON(t)["horizonte"].(map[string]any)
	ld, err := mf.LiveDebt(in, mf.NewDate(2027, 4, 10))
	if err != nil {
		t.Fatal(err)
	}
	if w := modelCents(t, h["deuda_viva_10_abr"]); int64(ld.LiveDebtCents) != w || w != 629344 {
		t.Errorf("deuda viva al 10-abr-2027: %d; model.json %d", ld.LiveDebtCents, w)
	}
	// Tarjeta «Mes contra plan»: Despensa 2,000.00 (fixture) contra el real del hogar de octubre (119,764 = 59.9 %).
	rows, err := mf.PlanVsActual(in, mf.NewDate(2026, 10, 1).MonthOf(), 1, "")
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].PlanCents == nil || int64(*rows[0].PlanCents) != 200000 || int64(rows[0].ActualCents) != 119764 {
		t.Errorf("mes contra plan: %+v", rows[0])
	}
	if pct := dominio.PercentOf(int64(rows[0].ActualCents), 200000); pct >= fe.WarningPct {
		t.Errorf("la tarjeta debería estar en «normal» (%d %% >= %d %%)", pct, fe.WarningPct)
	}
	// Tarjeta «Flujo de la quincena»: sueldo esperado 12,000.00 (quincena del 15) y 12,500.00 (la del 30), de la fixture.
	for _, c := range []struct {
		day  mf.Date
		want int64
	}{{mf.NewDate(2026, 10, 12), 1200000}, {mf.NewDate(2026, 10, 20), 1250000}} {
		ff, err := mf.FortnightFlow(in, mf.FortnightOf(c.day), 1)
		if err != nil {
			t.Fatal(err)
		}
		if int64(ff[0].Income.TotalCents) != c.want || !ff[0].Income.IncomeIsExpected {
			t.Errorf("quincena de %v: ingreso %d (esperado=%v), la fixture trae %d", c.day, ff[0].Income.TotalCents, ff[0].Income.IncomeIsExpected, c.want)
		}
	}
	// Tarjeta «Ahorro»: 314,672 de 629,344 = 50.0 % (RN-18), plan mensual 449.53.
	sv, err := mf.SavingsGoals(in)
	if err != nil {
		t.Fatal(err)
	}
	if int64(sv.MonthSummary.SavedCents) != 314672 || int64(sv.MonthSummary.ContributedThisMonthCents) != 314672 || int64(sv.MonthSummary.PlannedThisMonthCents) != 44953 ||
		len(sv.Goals) != 1 || sv.Goals[0].ProgressTenthsPct == nil || *sv.Goals[0].ProgressTenthsPct != 500 {
		t.Errorf("ahorro: %+v", sv)
	}
}

// Con la renta (gasto fijo) el motor suma 800,000 al devengado y al efectivo de cada mes: lo que era canon (P1 y P2)
// deja de serlo. Esta prueba deja escrito por qué NO es la carga por omisión (bloqueo WF3a-bis-semilla-dev, punto 1).
func TestDesarrolloConGastoFijoMueveLosMeses(t *testing.T) {
	_, sin := engine(t, sembrada(t, false))
	_, con := engine(t, sembrada(t, true))
	base, got, want := meses(t, sin), meses(t, con), modelPeriodos(t)
	for i := 0; i < 6; i++ {
		if got[i] != [2]int64{base[i][0] + 800000, base[i][1] + 800000} {
			t.Errorf("mes %d con renta = %v; sin renta + 800000 = %v", i+1, got[i], [2]int64{base[i][0] + 800000, base[i][1] + 800000})
		}
	}
	for i := 0; i < 2; i++ {
		if got[i] == want[i] {
			t.Errorf("P%d con renta no debería coincidir con model.json", i+1)
		}
	}
}

func TestDesarrolloPlanDeComidas(t *testing.T) {
	db := sembrada(t, false)
	if n := one[int](t, db, `SELECT count(DISTINCT recipe_id) FROM meal_plan_entry`); n != 6 {
		t.Errorf("recetas distintas: %d, se esperaban 6", n)
	}
	if n := one[int](t, db, `SELECT count(*) FROM meal_plan_entry WHERE plan_date NOT BETWEEN '2026-10-12' AND '2026-10-18'`); n != 0 {
		t.Errorf("%d entradas fuera de la semana del 12 al 18-oct", n)
	}
	// Precios: solo de la compra del 10-oct (Bodega del ticket y Walmart del carrito comparado), nunca inventados.
	if n := one[int](t, db, `SELECT count(*) FROM price WHERE observed_date <> '2026-10-10'`); n != 0 {
		t.Errorf("%d precios con otra fecha", n)
	}
	if n := one[int](t, db, `SELECT count(*) FROM price p JOIN expense_item ei ON ei.id = p.expense_item_id WHERE p.price_cents <> ei.unit_price_cents`); n != 0 {
		t.Errorf("%d precios de Bodega distintos del renglón del ticket", n)
	}
}
