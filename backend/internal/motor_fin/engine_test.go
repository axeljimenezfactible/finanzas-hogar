package motor_fin

import (
	"encoding/json"
	"fmt"
	"math/big"
	"path/filepath"
	"sort"
	"testing"
)

type M = map[string]any

func clone(t testing.TB, m M) M { return genericJSON(t, mustJSON(t, m)).(M) }

func push(m M, key string, v ...any) { m[key] = append(m[key].([]any), v...) }

func cashExp(id int, cat string, amt int64, date, kind string) M {
	return M{"id": id, "kind": kind, "category": cat, "amount_cents": amt, "expense_date": date, "payment_method": "cash",
		"paid_by": nil, "shares": []any{M{"person": "axel", "amount_cents": amt}}}
}

// ratHalfUp es un redondeo independiente con big.Rat (verificación, no el del motor).
func ratHalfUp(n, d int64) int64 {
	r := new(big.Rat).Add(big.NewRat(n, d), big.NewRat(1, 2))
	return new(big.Int).Quo(r.Num(), r.Denom()).Int64()
}

// ratOf convierte una hoja numérica genérica (json.Number) en un racional exacto.
func ratOf(v any) (*big.Rat, bool) {
	n, ok := v.(json.Number)
	if !ok {
		return nil, false
	}
	return new(big.Rat).SetString(string(n))
}

func sumOf(m map[string]any) (*big.Rat, bool) {
	s := new(big.Rat)
	for _, v := range m {
		r, ok := ratOf(v)
		if !ok {
			return nil, false
		}
		s.Add(s, r)
	}
	return s, true
}

// partsInvariants recorre una salida genérica y aplica las invariantes de S2-05 donde aparecen sus
// claves: Σ by_person = total_cents; payments_cents + cash_expenses_cents = total_cents;
// Σ shares = amount_cents; y en las vistas del oráculo, axel + karen = total (o monto).
func partsInvariants(v any, path string, f func(kind, path string, got, want *big.Rat)) {
	switch x := v.(type) {
	case []any:
		for i, e := range x {
			partsInvariants(e, fmt.Sprintf("%s[%d]", path, i), f)
		}
	case map[string]any:
		if tot, ok := ratOf(x["total_cents"]); ok {
			if bp, ok := x["by_person"].(map[string]any); ok {
				if s, ok := sumOf(bp); ok {
					f("Σ by_person = total_cents", path, s, tot)
				}
			}
			pc, ok1 := ratOf(x["payments_cents"])
			ce, ok2 := ratOf(x["cash_expenses_cents"])
			if ok1 && ok2 {
				f("payments_cents + cash_expenses_cents = total_cents", path, new(big.Rat).Add(pc, ce), tot)
			}
		}
		if amt, ok := ratOf(x["amount_cents"]); ok {
			if sh, ok := x["shares"].(map[string]any); ok {
				if s, ok := sumOf(sh); ok {
					f("Σ shares = amount_cents", path, s, amt)
				}
			}
		}
		ax, okA := ratOf(x["axel"])
		ka, okK := ratOf(x["karen"])
		if okA && okK {
			for _, key := range []string{"total", "monto"} {
				if w, ok := ratOf(x[key]); ok {
					f("axel + karen = "+key, path, new(big.Rat).Add(ax, ka), w)
					break
				}
			}
		}
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			partsInvariants(x[k], path+"."+k, f)
		}
	}
}

// TestSumaPartes es la prueba de propiedad de S2-05: la semilla mes por mes y toda salida de las
// consultas sin error de cada fixture de oraculo/fixtures/s2.
func TestSumaPartes(t *testing.T) {
	t.Run("fixtures", func(t *testing.T) {
		names, _ := filepath.Glob(filepath.Join(repoRoot(t), fixDir, "s2-*.json"))
		if len(names) == 0 {
			t.Fatalf("falta %s/s2-*.json (autor: agente Python de S2)", fixDir)
		}
		nq, n, kinds := 0, 0, map[string]int{}
		for _, f := range names {
			fx := readJSON(t, filepath.Join(fixDir, filepath.Base(f))).(map[string]any)
			ent := fx["entradas"].(map[string]any)
			datos := mustJSON(t, resolveData(t, ent["datos"]))
			for _, c := range ent["consultas"].([]any) {
				qq := c.(map[string]any)
				out, e := Eval(datos, qq["funcion"].(string), mustJSON(t, qq["argumentos"]))
				if e != nil {
					continue // consulta de error: no tiene partes (TestOraculoFixtures compara el código)
				}
				nq++
				where := fx["caso"].(string) + " " + qq["id"].(string)
				partsInvariants(genericJSON(t, mustJSON(t, out)), "", func(kind, path string, got, want *big.Rat) {
					n++
					kinds[kind]++
					if got.Cmp(want) != 0 {
						t.Errorf("%s%s: %s da %s, esperado %s", where, path, kind, got.RatString(), want.RatString())
					}
				})
			}
		}
		if n == 0 {
			t.Fatal("ninguna invariante evaluada")
		}
		t.Logf("S2-05: %d fixtures, %d consultas sin error, %d invariantes con diferencia 0 %v", len(names), nq, n, kinds)
	})
	t.Run("semilla-12-meses", testSumaSemilla)
}

func testSumaSemilla(t *testing.T) {
	seed := seedFromModel(t)
	n := 0
	for i := 0; i < 12; i++ {
		m := (Month{2026, 10}).Add(i).String()
		p := leaf(t, q(t, seed, "periodos", `{"months":["`+m+`"]}`), "[0]")
		for _, b := range []string{"accrued", "cash"} {
			eq(t, m+" "+b+" Σ personas", qc(t, p, b+".by_person.axel")+qc(t, p, b+".by_person.karen"), qc(t, p, b+".total_cents"))
			n++
		}
		eq(t, m+" pagos + contado", qc(t, p, "cash.payments_cents")+qc(t, p, "cash.cash_expenses_cents"), qc(t, p, "cash.total_cents"))
		for j := range leaf(t, p, "cash.payments").([]any) {
			x := fmt.Sprintf("cash.payments[%d].", j)
			eq(t, m+" partes del pago", qc(t, p, x+"shares.axel")+qc(t, p, x+"shares.karen"), qc(t, p, x+"amount_cents"))
			n++
		}
	}
	t.Logf("S2-05: %d casos de Σ partes = total con diferencia 0", n)
}

func TestQuincenaFinDeMes(t *testing.T) {
	seed := seedFromModel(t)
	cal := func(t *testing.T, s M, d, f, a, b, pay string) {
		c := q(t, s, "calendario", `{"date":"`+d+`"}`)
		for k, v := range map[string]string{"fortnight": f, "fortnight_start": a, "fortnight_end": b, "pay_date": pay} {
			eq(t, d+" "+k, leaf(t, c, k), v)
		}
	}
	t.Run("feb-2027", func(t *testing.T) {
		cal(t, seed, "2027-02-20", "2027-02-q2", "2027-02-16", "2027-02-28", ms(t, "despensas[2].pagos[0].fecha"))
	})
	t.Run("mes-30-dias", func(t *testing.T) {
		cal(t, seed, "2026-11-20", "2026-11-q2", "2026-11-16", "2026-11-30", ms(t, "despensas[0].pagos[2].fecha"))
	})
	t.Run("mes-31-dias", func(t *testing.T) {
		s := clone(t, seed)
		s["params"].(M)["today"] = "2026-10-31"
		push(s, "expenses", cashExp(7, "despensa", 50000, "2026-10-31", "variable"))
		cal(t, s, "2026-10-31", "2026-10-q2", "2026-10-16", "2026-10-31", ms(t, "despensas[0].pagos[0].fecha"))
		f := q(t, s, "flujo_quincenal", `{"from":"2026-10-q1","count":2}`)
		eq(t, "q1 variable", qc(t, f, "[0].variable_spent_cents"), 0)
		eq(t, "q1 gastado", qc(t, f, "[0].spent_cents"), 0)
		eq(t, "q2 variable (gasto del 31)", qc(t, f, "[1].variable_spent_cents"), 50000)
		eq(t, "q2 gastado", qc(t, f, "[1].spent_cents"), 50000)
		eq(t, "q2 pendiente del hogar", qc(t, f, "[1].credit_payments.household_pending_cents"), mc(t, "despensas[0].pagos[0].axel"))
		p := q(t, s, "periodos", `{"months":["2026-10"]}`)
		eq(t, "contado oct", qc(t, p, "[0].cash.cash_expenses_cents"), 50000)
		eq(t, "efectivo oct", qc(t, p, "[0].cash.total_cents"), mc(t, "periodos[0].efectivo.total")+50000)
	})
	t.Run("semana", func(t *testing.T) {
		c := q(t, seed, "calendario", `{"date":"2026-10-14"}`)
		eq(t, "semana", sp(leaf(t, c, "week"), leaf(t, c, "week_start"), leaf(t, c, "week_end")), "2026-W42 2026-10-12 2026-10-18")
	})
}

func TestInvarianteDeudaViva(t *testing.T) {
	seed := seedFromModel(t)
	cuts := []string{"2026-10-31", "2026-11-30", "2026-12-31", "2027-01-31", "2027-02-28", "2027-03-31", "2027-04-10"}
	for k, cut := range cuts {
		t.Run(cut, func(t *testing.T) {
			r := q(t, seed, "deuda_viva", `{"cut":"`+cut+`"}`)
			c, _ := ParseDate(cut)
			var live int64 // desde model.json: pagos posteriores al corte de compras hechas hasta el corte
			for j := 0; j < 4; j++ {
				if pd, _ := ParseDate(ms(t, fmt.Sprintf("despensas[%d].fecha", j))); pd.After(c) {
					continue
				}
				for i := 0; i < 4; i++ {
					if d, _ := ParseDate(ms(t, fmt.Sprintf("despensas[%d].pagos[%d].fecha", j, i))); d.After(c) {
						live += int64(mc(t, fmt.Sprintf("despensas[%d].pagos[%d].monto", j, i)))
					}
				}
			}
			eq(t, "deuda viva", qc(t, r, "live_debt_cents"), live)
			eq(t, "identidad K19", qc(t, r, "cumulative_diff_cents")+qc(t, r, "unaccrued_balance_cents"), live)
			eq(t, "identity_holds", leaf(t, r, "identity_holds"), true)
			if k < 6 {
				eq(t, "diferencia acumulada", qc(t, r, "cumulative_diff_cents"), mc(t, fmt.Sprintf("periodos[%d].diferencia_acumulada", k)))
				eq(t, "vale sin saldo", leaf(t, r, "identity_without_balance_holds"), k%2 == 1)
			} else {
				eq(t, "R-53", live, mc(t, "horizonte.deuda_viva_10_abr"))
				eq(t, "R-58", qc(t, r, "cumulative_diff_cents"), mc(t, "horizonte.diferencia_final"))
				eq(t, "R-55 saldo", qc(t, r, "unaccrued_balance_cents"), mc(t, "horizonte.deuda_viva_detalle.despensa_4"))
			}
			t.Logf("S2-07 %s: %d = %d + %d", cut, live, qc(t, r, "cumulative_diff_cents"), qc(t, r, "unaccrued_balance_cents"))
		})
	}
}

func TestPromedios(t *testing.T) {
	seed := seedFromModel(t)
	for _, b := range []string{"efectivo", "devengado"} {
		r := q(t, seed, "promedio_mensual", fmt.Sprintf(`{"total_cents":%d,"month_count":6}`, mc(t, "horizonte.total_"+b)))
		eq(t, "mensual "+b, qc(t, r, "monthly_cents"), mc(t, "horizonte.promedio_mensual_"+b))
	}
	for _, c := range []struct {
		month, basis, path string
		k                  int
	}{{"2026-10", "cash", "efectivo.total", 0}, {"2026-11", "cash", "efectivo.total", 1}, {"2027-02", "cash", "efectivo.total", 4}, {"2026-10", "accrued_groceries", "devengado.despensa_total", 0}} {
		r := q(t, seed, "promedios", `{"month":"`+c.month+`","basis":"`+c.basis+`"}`)
		T, d := int64(mc(t, fmt.Sprintf("periodos[%d].%s", c.k, c.path))), qc(t, r, "days")
		m, _ := ParseMonth(c.month)
		eq(t, c.month+" días reales", d, m.Days())
		eq(t, c.month+" diario", qc(t, r, "daily_cents"), ratHalfUp(T, d))
		eq(t, c.month+" semanal (un solo redondeo, K39)", qc(t, r, "weekly_cents"), ratHalfUp(7*T, d))
		t.Logf("S2-08 %s %s: T=%d d=%d diario=%d semanal=%d (diario redondeado ×7 = %d)", c.month, c.basis, T, d, qc(t, r, "daily_cents"), qc(t, r, "weekly_cents"), 7*ratHalfUp(T, d))
	}
	r := q(t, seed, "promedios", `{"total_cents":97375,"start":"2026-10-01","end":"2026-10-31"}`)
	eq(t, "forma total/start/end", qc(t, r, "weekly_cents"), ratHalfUp(7*97375, 31))
}

func incomePlan() M {
	return M{"id": 1, "kind": "income", "description": "Quincena", "amount_cents": 1000000, "second_half_amount_cents": 1200000,
		"frequency": "fortnightly", "interval_count": 1, "start_date": "2026-10-15", "end_date": nil, "person": "axel", "share_rule": nil, "credit_id": nil, "is_active": true}
}

func TestProyeccionDefecto(t *testing.T) {
	seed := seedFromModel(t)
	p := q(t, seed, "proyeccion", `{"from":"2026-10"}`)
	var acc712 int64
	var dates []any
	for i := 0; i < 12; i++ {
		x := fmt.Sprintf("months[%d].", i)
		eq(t, x+"provisional", leaf(t, p, x+"provisional"), nil)
		if i < 6 {
			P := fmt.Sprintf("periodos[%d].", i)
			eq(t, x+"devengado = P", qc(t, p, x+"accrued.total_cents"), mc(t, P+"devengado.total"))
			eq(t, x+"efectivo = P", qc(t, p, x+"cash.total_cents"), mc(t, P+"efectivo.total"))
			eq(t, x+"diferencia", qc(t, p, x+"cumulative_diff_cents"), mc(t, P+"diferencia_acumulada"))
			continue
		}
		acc712 += qc(t, p, x+"accrued.total_cents")
		for j := range leaf(t, p, x+"cash.payments").([]any) {
			dates = append(dates, leaf(t, p, fmt.Sprintf("%scash.payments[%d].due_date", x, j)))
		}
	}
	eq(t, "devengo 7–12 = despensa 4", acc712, mc(t, "despensas[3].total"))
	eq(t, "pagos 7–12", fmt.Sprint(dates), fmt.Sprint([]any{ms(t, "despensas[2].pagos[3].fecha"), ms(t, "despensas[3].pagos[0].fecha"),
		ms(t, "despensas[3].pagos[1].fecha"), ms(t, "despensas[3].pagos[2].fecha"), ms(t, "despensas[3].pagos[3].fecha")}))
	s := clone(t, seed)
	hogar := M{"parts": []any{}, "rest_person": "axel"}
	conj := M{"parts": []any{M{"person": "karen", "pct_bp": 5000}}, "rest_person": "axel"}
	push(s, "planned_items", incomePlan(),
		M{"id": 2, "kind": "fixed_expense", "description": "Renta", "amount_cents": 500000, "frequency": "monthly", "start_date": "2026-10-05", "category": "despensa", "share_rule": hogar},
		M{"id": 3, "kind": "upcoming_expense", "description": "Regalo", "amount_cents": 300000, "frequency": "once", "start_date": "2026-12-20", "category": "despensa", "accrual_months": 2, "share_rule": conj})
	c := q(t, s, "proyeccion", `{}`)
	eq(t, "ingreso oct (15 y 30)", qc(t, c, "months[0].income_cents"), 2200000)
	eq(t, "partidas dic", qc(t, c, "months[2].planned_expense_cents"), 800000)
	eq(t, "devengo dic", qc(t, c, "months[2].accrued.total_cents"), mc(t, "periodos[2].devengado.total")+500000+150000)
	eq(t, "Karen dic", qc(t, c, "months[2].accrued.by_person.karen"), mc(t, "periodos[2].devengado.karen")+75000)
	eq(t, "neto oct", qc(t, c, "months[0].net_cash_cents"), 2200000-int64(mc(t, "periodos[0].efectivo.axel"))-500000)
	eq(t, "csv provisional", leaf(t, c, "months[0].provisional"), nil)
}

func TestProyeccionCadencia(t *testing.T) {
	seed := seedFromModel(t)
	s := clone(t, seed)
	s["params"].(M)["repeat_cadence_enabled"] = true
	s["cadence"].(M)["expenses"] = []any{cashExp(20, "papel", 19900, "2027-05-10", "variable")}
	on, off := q(t, s, "proyeccion", `{"from":"2026-10"}`), q(t, seed, "proyeccion", `{"from":"2026-10"}`)
	eq(t, "papel alterno en may", qc(t, on, "months[7].cash.total_cents")-qc(t, off, "months[7].cash.total_cents"), 19900)
	eq(t, "sin marca mes 7 (FIRME, Q2-21 a Q2-24)", leaf(t, on, "months[6].provisional"), nil)
	eq(t, "sin marca mes 6", leaf(t, on, "months[5].provisional"), nil)
	s["cadence"].(M)["expenses"] = []any{cashExp(20, "papel", 19900, "2027-03-10", "variable")}
	if _, e := Eval(mustJSON(t, s), "proyeccion", []byte(`{"from":"2026-10"}`)); e == nil || e.Code != CodeCadenceBefore7 {
		t.Errorf("cadencia antes del mes 7: %v", e)
	}
	l := clone(t, seed) // s2-liga-provisional
	push(l, "planned_items", incomePlan())
	push(l, "incomes", M{"id": 1, "person": "axel", "kind": "fixed", "amount_cents": 990000, "income_date": "2026-10-14", "planned_item_id": 1})
	f := q(t, l, "flujo_quincenal", `{"from":"2026-10-q1","count":2}`)
	eq(t, "q1 sin marca (FIRME, Q2-14 = A)", leaf(t, f, "[0].provisional"), nil)
	eq(t, "q1 esperado", qc(t, f, "[0].income.expected_cents"), 0)
	eq(t, "q2 esperado", qc(t, f, "[1].income.expected_cents"), 1200000)
	eq(t, "q2 marca", leaf(t, f, "[1].provisional"), nil)
	p := q(t, l, "proyeccion", `{"from":"2026-10"}`)
	eq(t, "proy oct sin marca", leaf(t, p, "months[0].provisional"), nil)
	eq(t, "proy nov sin marca", leaf(t, p, "months[1].provisional"), nil)
}

func TestFlujoQuincenal(t *testing.T) {
	seed := seedFromModel(t)
	f := q(t, seed, "flujo_quincenal", `{"from":"2026-10-q2","count":16}`)
	plus, minus := 0, 0
	for i := 0; i < 16; i++ {
		k, j, x := i/4, i%4, fmt.Sprintf("[%d].credit_payments.", i)
		P := fmt.Sprintf("despensas[%d].pagos[%d].", k, j)
		eq(t, x+"hogar", qc(t, f, x+"household_cents"), mc(t, P+"axel"))
		eq(t, x+"karen", qc(t, f, x+"by_person.karen"), mc(t, P+"karen"))
		eq(t, x+"informativo", qc(t, f, x+"external_informative_cents"), mc(t, P+"karen"))
		if k > 0 {
			d := int64(mc(t, "valor_estable.quincena.axel")) - qc(t, f, x+"household_cents")
			plus, minus = plus+map[bool]int{true: 1}[d == 1], minus+map[bool]int{true: 1}[d == -1]
		}
	}
	for k := 1; k < 4; k++ {
		var h, kr int64
		for j := 0; j < 4; j++ {
			h += qc(t, f, fmt.Sprintf("[%d].credit_payments.household_cents", 4*k+j))
			kr += qc(t, f, fmt.Sprintf("[%d].credit_payments.by_person.karen", 4*k+j))
		}
		eq(t, fmt.Sprintf("ciclo despensa %d hogar", k+1), h, 2*int64(mc(t, "valor_estable.mes.axel")))
		eq(t, fmt.Sprintf("ciclo despensa %d karen", k+1), kr, 2*int64(mc(t, "valor_estable.componentes.karen_mitad_despensa")))
	}
	t.Logf("S2-10 contraste valor estable − pago: +1 en %d quincenas, −1 en %d (se informa, no se corrige)", plus, minus)
}

func TestPresupuestoUmbrales(t *testing.T) {
	seed := seedFromModel(t)
	for _, c := range []struct {
		amt   int64
		pct   int
		state string
		w, l  bool
	}{{79000, 79, "normal", false, false}, {79999, 80, "normal", false, false}, {80000, 80, "warning", true, false},
		{99000, 99, "warning", true, false}, {100000, 100, "exceeded", true, true}, {101000, 101, "exceeded", true, true}} {
		s := clone(t, seed)
		push(s, "categories", M{"code": "hogar", "name": "Hogar", "kind": "variable", "is_active": true})
		push(s, "budgets", M{"id": 1, "category": "hogar", "amount_cents": 100000, "valid_from": "2026-12-01", "valid_to": nil})
		push(s, "expenses", cashExp(7, "hogar", c.amt, "2026-12-05", "variable"))
		r := q(t, s, "presupuesto_mes", `{"month":"2026-12"}`)
		x := "categories[1]."
		eq(t, fmt.Sprint(c.amt, " estado"), sp(leaf(t, r, x+"used_pct"), leaf(t, r, x+"state"), leaf(t, r, x+"warning"), leaf(t, r, x+"limit")), sp(c.pct, c.state, c.w, c.l))
		eq(t, "total real", qc(t, r, "total.actual_cents"), int64(mc(t, "periodos[2].devengado.axel"))+c.amt)
	}
	s := clone(t, seed)
	push(s, "budgets", M{"id": 1, "category": "despensa", "amount_cents": 200000, "valid_from": "2026-10-01", "valid_to": nil})
	r := q(t, s, "presupuesto_mes", `{"month":"2026-12"}`)
	eq(t, "RN-22", sp(leaf(t, r, "categories[0].used_pct"), leaf(t, r, "categories[0].warning"), leaf(t, r, "categories[0].limit")), "88 true false")
	o := q(t, s, "presupuesto_mes", `{"month":"2026-10"}`)
	eq(t, "devengado hogar oct", qc(t, o, "accrued_household_cents"), mc(t, "periodos[0].devengado.axel"))
	eq(t, "efectivo hogar oct", qc(t, o, "cash_household_cents"), mc(t, "periodos[0].efectivo.axel"))
	eq(t, "días que faltan", qc(t, o, "days_left"), 21)
	pv := q(t, s, "plan_contra_real", `{"from":"2026-11","count":3,"category":"despensa"}`)
	eq(t, "plan contra real dic", qc(t, pv, "[1].actual_cents"), mc(t, "periodos[2].devengado.axel"))
	q(t, s, "plan_contra_real", `{"from":"2026-11","count":2}`)
}

func TestPresupuestoScope(t *testing.T) {
	s := seedFromModel(t)
	push(s, "budgets", M{"id": 1, "category": "despensa", "amount_cents": 200000, "valid_from": "2026-10-01", "valid_to": nil})
	for scope, want := range map[string]string{
		"onward":     "[1 2026-10-01 2026-11-30 200000] [2 2026-12-01 <nil> 250000]",
		"this_month": "[1 2026-10-01 2026-11-30 200000] [2 2026-12-01 2026-12-31 250000] [3 2027-01-01 <nil> 200000]"} {
		r := q(t, s, "aplicar_cambio_presupuesto", `{"row_id":1,"month":"2026-12","amount_cents":250000,"scope":"`+scope+`"}`)
		got := ""
		for i, row := range leaf(t, r, "rows").([]any) {
			x := row.(M)
			got += fmt.Sprint([]any{x["id"], x["valid_from"], x["valid_to"], x["amount_cents"]})
			if i < len(leaf(t, r, "rows").([]any))-1 {
				got += " "
			}
		}
		eq(t, scope, got, want)
	}
}

func TestMetasAhorro(t *testing.T) {
	s := seedFromModel(t)
	goal := func(id int, target int64, date string) M {
		return M{"id": id, "name": "Meta", "kind": "goal", "target_cents": target, "target_date": date, "planned_monthly_cents": 10000, "start_date": "2026-10-01"}
	}
	dv := int64(mc(t, "horizonte.deuda_viva_10_abr"))
	push(s, "savings_goals", goal(1, 600000, "2027-03-31"), goal(2, 100001, "2027-01-15"), goal(3, dv, "2027-04-10"), goal(4, dv, "2027-04-10"),
		M{"id": 5, "name": "Fondo", "kind": "fund", "start_date": "2026-10-01"})
	push(s, "savings_movements", M{"goal_id": 3, "movement_date": "2026-10-01", "kind": "deposit", "amount_cents": dv},
		M{"goal_id": 4, "movement_date": "2026-10-02", "kind": "deposit", "amount_cents": dv},
		M{"goal_id": 4, "movement_date": "2026-10-03", "kind": "withdrawal", "amount_cents": dv - dv/2})
	r := q(t, s, "metas", `{}`)
	for g := 0; g < 2; g++ {
		var sum int64
		for i := range leaf(t, r, fmt.Sprintf("goals[%d].plan", g)).([]any) {
			sum += qc(t, r, fmt.Sprintf("goals[%d].plan[%d].cents", g, i))
		}
		eq(t, fmt.Sprintf("meta %d Σ plan = objetivo", g+1), sum, qc(t, r, fmt.Sprintf("goals[%d].target_cents", g)))
	}
	eq(t, "residuo: la última ajusta", qc(t, r, "goals[1].plan[3].cents"), 25001)
	eq(t, "cumplida", sp(len(leaf(t, r, "goals[2].plan").([]any)), qc(t, r, "goals[2].progress_tenths_pct"), leaf(t, r, "goals[2].state")), "0 1000 achieved")
	eq(t, "RN-18 50.0 %", qc(t, r, "goals[3].progress_tenths_pct"), 500)
	eq(t, "fondo sin monto", leaf(t, r, "goals[4].progress_pct"), nil)
	eq(t, "aportado oct", qc(t, r, "month_summary.contributed_this_month_cents"), 2*dv)
}

func TestZonaHoy(t *testing.T) {
	s := seedFromModel(t)
	for now, want := range map[string]string{"2026-10-31T05:59:59Z": "2026-10-30", "2026-10-31T06:00:00Z": "2026-10-31"} {
		eq(t, now, leaf(t, q(t, s, "calendario", `{"now":"`+now+`"}`), "today"), want)
		x := clone(t, s)
		x["params"].(M)["today"], x["params"].(M)["now"] = nil, now
		in, err := DecodeInput(mustJSON(t, x))
		if err != nil || in.Today().String() != want {
			t.Errorf("params.now %s: %v %v", now, in.Today(), err)
		}
	}
}
