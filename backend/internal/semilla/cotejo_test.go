package semilla_test

import (
	"context"
	"io"
	"log/slog"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comprasa"
	"finanzas-hogar/internal/api/presupuesto"
	"finanzas-hogar/internal/jobs"
	"finanzas-hogar/internal/semilla"
	"finanzas-hogar/internal/testutil"
)

// Cotejo de la semilla REAL con el oráculo a través de la API (WF-3b, agente semilla-c3; C3-09):
// getAccrualVsCash P1–P6 = model.json (R-45…R-50), deuda viva al 10-abr-2027 = 6,293.44 y la comparativa de canasta
// Bodega contra Walmart de data.json (+138.50 en la canasta completa y +100.50 en los artículos idénticos).
// Decisiones provisionales del orquestador: bloqueo WF3a-bis-semilla-dev, puntos 4 y 7.

// servidorSembrado es la API (presupuesto y compras) sobre una base con SOLO la semilla real, con el reloj en 12-oct-2026.
func servidorSembrado(t *testing.T) (*testutil.Server, *testutil.Client) {
	t.Helper()
	db := soloReal(t)
	clk := testutil.NewClock(testutil.MX(2026, time.October, 12, 12, 0))
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	tmp := testutil.TempDir(t)
	deps := &api.Deps{DB: db, Clock: clk, Log: log, Jobs: jobs.New(db, clk, log), Config: api.Config{
		DataDir: filepath.Join(tmp, "datos"), BackupDir: filepath.Join(tmp, "respaldos"), TmpDir: tmp,
		OpenAPIPath: filepath.Join(testutil.Root(t), "docs", "contratos", "openapi.yaml"), SecureHeaders: true}}
	reg := api.NewRegistry(deps)
	presupuesto.Register(reg)
	comprasa.Register(reg)
	api.RegisterBase(reg)
	ts := httptest.NewTLSServer(reg.Handler())
	t.Cleanup(ts.Close)
	s := &testutil.Server{Server: ts, DB: db, Clock: clk, Reg: reg}
	return s, s.SessionFor(t, "axel")
}

func TestCotejoAccrualVsCashP1aP6ConModelJSON(t *testing.T) {
	_, c := servidorSembrado(t)
	var got api.AccrualVsCash
	res := c.Do("GET", "/api/v1/cash-flow/accrual-vs-cash", nil)
	if res.StatusCode != 200 {
		t.Fatalf("getAccrualVsCash: %d %s", res.StatusCode, res.Body)
	}
	res.JSON(t, &got)
	want := modelPeriodos(t)
	meses := []string{"2026-10", "2026-11", "2026-12", "2027-01", "2027-02", "2027-03"}
	if len(got.Months) < 6 {
		t.Fatalf("getAccrualVsCash trae %d meses", len(got.Months))
	}
	for i := 0; i < 6; i++ {
		m := got.Months[i]
		if string(m.Month) != meses[i] {
			t.Errorf("mes %d = %s, se esperaba %s", i+1, m.Month, meses[i])
		}
		if int64(m.AccruedCents) != want[i][0] || int64(m.CashCents) != want[i][1] {
			t.Errorf("P%d (R-%d) %s: devengado %d y efectivo %d; model.json %d y %d", i+1, 45+i, m.Month, m.AccruedCents, m.CashCents, want[i][0], want[i][1])
		}
		if !m.IsCanonical {
			t.Errorf("P%d: is_canonical = false", i+1)
		}
	}
	// La diferencia acumulada de P6 es la de model.json (1,258.70).
	if got.Months[5].CumulativeDifferenceCents != 125870 {
		t.Errorf("diferencia acumulada de P6: %d, model.json 125870", got.Months[5].CumulativeDifferenceCents)
	}
}

func TestCotejoDeudaVivaAl10Abril(t *testing.T) {
	_, c := servidorSembrado(t)
	res := c.Do("GET", "/api/v1/debt/live?date=2027-04-10", nil)
	if res.StatusCode != 200 {
		t.Fatalf("getLiveDebt: %d %s", res.StatusCode, res.Body)
	}
	if !strings.Contains(string(res.Body), "629344") {
		t.Errorf("la deuda viva al 10-abr-2027 debía ser 629344 (6,293.44): %s", res.Body)
	}
	// Y el cruce del motor de crédito sobre la misma base.
	cr, err := semilla.Verify(context.Background(), soloReal(t), seedDir(t), "")
	if err != nil || cr.DeudaViva != 629344 {
		t.Errorf("Verify: %+v %v", cr, err)
	}
}

func TestCotejoComparativaCanastaBodegaVsWalmart(t *testing.T) {
	s, c := servidorSembrado(t)
	var bodega, walmart, eid int64
	for _, q := range []struct {
		dst *int64
		sql string
	}{{&bodega, `SELECT s.id FROM store s JOIN store c ON c.id = s.parent_id WHERE s.kind = 'branch' AND c.name = 'Bodega Aurrera'`},
		{&walmart, `SELECT s.id FROM store s JOIN store c ON c.id = s.parent_id WHERE s.kind = 'branch' AND c.name = 'Walmart'`},
		{&eid, `SELECT id FROM expense WHERE expense_date = '2026-10-10'`}} {
		if err := s.DB.QueryRow(q.sql).Scan(q.dst); err != nil {
			t.Fatal(err)
		}
	}
	var bc api.BasketComparison
	res := c.Do("GET", "/api/v1/basket-comparison?basket=expense&expense_id="+strconv.FormatInt(eid, 10)+"&week=2026-10-12", nil)
	if res.StatusCode != 200 {
		t.Fatalf("getBasketComparison: %d %s", res.StatusCode, res.Body)
	}
	res.JSON(t, &bc)
	var b, w *api.BasketComparisonStoresItem
	for i := range bc.Stores {
		switch bc.Stores[i].StoreID {
		case bodega:
			b = &bc.Stores[i]
		case walmart:
			w = &bc.Stores[i]
		}
	}
	if b == nil || w == nil {
		t.Fatalf("faltan Bodega o Walmart en la comparativa: %+v", bc.Stores)
	}
	// Artículos idénticos (data.json: Bodega 2,773.50 y Walmart 2,874.00, +100.50). El motor S5 cuenta además la sexta lata de
	// atún (21.00 en las dos tiendas): 2,794.50 y 2,895.00; la diferencia, +100.50, es exactamente la de data.json.
	if int64(w.IdenticalTotalCents-b.IdenticalTotalCents) != 10050 || int64(b.IdenticalTotalCents) != 277350+2100 || int64(w.IdenticalTotalCents) != 287400+2100 {
		t.Errorf("idénticos: Bodega %d, Walmart %d (diferencia %d); data.json +2100 de la sexta lata: 279450, 289500 y 10050", b.IdenticalTotalCents, w.IdenticalTotalCents, w.IdenticalTotalCents-b.IdenticalTotalCents)
	}
	// Canasta completa: Bodega 3,094.50 = data.json. Walmart: data.json trae 3,233.00 (+138.50) con 60 artículos y el Coffee Mate
	// de 640 g como equivalente del pack de 780 g de Bodega. El motor S5 (§7.3) compra los paquetes que cubren la cantidad:
	// +21.00 por la sexta lata de atún que data.json no pone del lado de Walmart y +105.00 por el segundo Coffee Mate de 640 g.
	// 138.50 + 21.00 + 105.00 = 264.50. Es una diferencia de método, no de datos (hallazgo para Axel, ver el reporte).
	if int64(b.TotalCents) != 309450 || int64(w.TotalCents) != 323300+2100+10500 || int64(w.TotalCents-b.TotalCents) != 13850+2100+10500 {
		t.Errorf("canasta completa: Bodega %d, Walmart %d (diferencia %d); se esperaban 309450, 335900 y 26450 (data.json 323300 y 13850 + 2100 + 10500)", b.TotalCents, w.TotalCents, w.TotalCents-b.TotalCents)
	}
	if n := one[int](t, s.DB, `SELECT count(*) FROM price WHERE source = 'manual' AND store_id = ?`, walmart); n != 36 {
		t.Errorf("precios de Walmart sembrados: %d, data.json trae 36", n)
	}
}

// La semilla real marca los gastos proyectados, los liga a sus compras y es idempotente también al "actualizar" una
// base que ya tenía las compras a crédito sin su gasto (la forma anterior de la semilla).
func TestProyectadosSeLigan(t *testing.T) {
	db := soloReal(t)
	if n := one[int](t, db, `SELECT count(*) FROM credit_purchase WHERE expense_id IS NOT NULL`); n != 4 {
		t.Errorf("compras ligadas a su gasto: %d, se esperaban 4", n)
	}
	if n := one[int](t, db, `SELECT count(*) FROM expense WHERE expense_date > '2026-10-12' AND notes LIKE 'PROYECTADO%' AND description LIKE '%(proyectado)'`); n != 5 {
		t.Errorf("gastos proyectados marcados: %d, se esperaban 5", n)
	}
	for _, c := range []struct {
		date string
		cat  string
		amt  int
	}{{"2026-12-10", "groceries", 400000}, {"2027-02-10", "groceries", 400000}, {"2027-04-10", "groceries", 400000}, {"2027-01-10", "papel", 19900}, {"2027-03-10", "papel", 19900}} {
		if n := one[int](t, db, `SELECT count(*) FROM expense e JOIN category c ON c.id = e.category_id WHERE e.expense_date = ? AND c.code = ? AND e.amount_cents = ?`, c.date, c.cat, c.amt); n != 1 {
			t.Errorf("gasto %s %s %d: %d", c.date, c.cat, c.amt, n)
		}
	}
	// Los papeles de contado los paga Karen de su bolsillo y son 100 % de Karen.
	if n := one[int](t, db, `SELECT count(*) FROM expense e JOIN person p ON p.id = e.paid_by_person_id WHERE e.payment_method = 'cash' AND p.code = 'karen'`); n != 2 {
		t.Errorf("papeles de contado de Karen: %d", n)
	}
	// Base de la forma anterior: se desligan las compras 2 a 4 y se borran sus gastos; la semilla las vuelve a ligar.
	if _, err := db.Exec(`UPDATE credit_purchase SET expense_id = NULL WHERE purchase_date > '2026-10-12'`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`DELETE FROM expense WHERE expense_date > '2026-10-12'`); err != nil {
		t.Fatal(err)
	}
	rep, err := semilla.RunWith(context.Background(), db, opts(t))
	if err != nil {
		t.Fatal(err)
	}
	if rep.Get("gastos proyectados").Created != 5 || rep.Get("compras a crédito ligadas a su gasto").Created != 3 {
		t.Errorf("reparación: %+v", rep.Entries)
	}
	if n := one[int](t, db, `SELECT count(*) FROM credit_purchase WHERE expense_id IS NOT NULL`); n != 4 {
		t.Errorf("compras ligadas tras reparar: %d", n)
	}
	if n := one[int](t, db, `SELECT count(*) FROM audit_log WHERE action = 'update' AND table_name = 'credit_purchase' AND request_id = 'semilla'`); n != 6 { // 3 al sembrar de cero (la compra se inserta y luego se liga) y 3 al reparar
		t.Errorf("audit_log de las ligas: %d", n)
	}
	// Segunda carga: nada nuevo.
	rep, err = semilla.RunWith(context.Background(), db, opts(t))
	if err != nil || rep.TotalCreated() != 0 {
		t.Errorf("segunda carga: %v %+v", err, rep.Entries)
	}
}
