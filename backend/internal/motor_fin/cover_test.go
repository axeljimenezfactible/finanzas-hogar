package motor_fin

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func sp(a ...any) string {
	s := make([]string, len(a))
	for i, x := range a {
		s[i] = strings.TrimSpace(strings.ReplaceAll(sprint(x), "\n", ""))
	}
	return strings.Join(s, " ")
}

func sprint(x any) string { return strings.TrimSpace(strings.Trim(string(mustJSONq(x)), `"`)) }

func mustJSONq(x any) []byte { b, _ := json.Marshal(x); return b }

// TestErrores recorre los códigos de §3.15 y ramas de argumentos (S2-15, cobertura).
func TestErrores(t *testing.T) {
	seed := seedFromModel(t)
	set := func(f func(M)) M { s := clone(t, seed); f(s); return s }
	exp0 := func(s M) M { return s["expenses"].([]any)[0].(M) }
	pur0 := func(s M) M { return s["credit_purchases"].([]any)[0].(M) }
	cases := []struct {
		name, fn, args, code string
		in                   M
	}{
		{"campo", "metas", `{}`, CodeUnknownField, set(func(s M) { s["params"].(M)["x"] = 1 })},
		{"zona", "metas", `{}`, CodeInvalidZone, set(func(s M) { s["params"].(M)["zone"] = "UTC" })},
		{"fecha", "metas", `{}`, CodeInvalidDate, set(func(s M) { s["params"].(M)["today"] = "2026-13-01" })},
		{"hoy", "metas", `{}`, CodeInvalidArgument, set(func(s M) { s["params"].(M)["now"] = "2026-10-11T00:00:00Z" })},
		{"partes", "metas", `{}`, CodeSharesMismatch, set(func(s M) { exp0(s)["amount_cents"] = 1 })},
		{"compra", "metas", `{}`, CodeMissingPurchase, set(func(s M) { exp0(s)["credit_purchase_id"] = 99 })},
		{"pago", "metas", `{}`, CodeSumMismatch, set(func(s M) { pur0(s)["payments"].([]any)[0].(M)["interest_cents"] = 1 })},
		{"total", "metas", `{}`, CodeSumMismatch, set(func(s M) { pur0(s)["total_cents"] = 1 })},
		{"capital", "metas", `{}`, CodeSumMismatch, set(func(s M) {
			exp0(s)["amount_cents"], exp0(s)["shares"] = 5, []any{M{"person": "axel", "amount_cents": 5}}
		})},
		{"ingreso", "metas", `{}`, CodeIncomeNotHolder, set(func(s M) {
			push(s, "incomes", M{"id": 1, "person": "karen", "kind": "extra", "amount_cents": 5, "income_date": "2026-10-01"})
		})},
		{"traslape", "metas", `{}`, CodeBudgetOverlap, set(func(s M) {
			push(s, "budgets", M{"id": 1, "category": nil, "amount_cents": 5, "valid_from": "2026-10-01"}, M{"id": 2, "category": nil, "amount_cents": 5, "valid_from": "2026-11-01"})
		})},
		{"regla", "metas", `{}`, CodeInvalidArgument, set(func(s M) {
			push(s, "planned_items", M{"id": 9, "kind": "fixed_expense", "amount_cents": 5, "frequency": "monthly", "start_date": "2026-10-01"})
		})},
		{"persona", "metas", `{}`, CodeInvalidArgument, set(func(s M) { exp0(s)["paid_by"] = "ximena" })},
		{"titular", "metas", `{}`, CodeInvalidArgument, set(func(s M) { s["params"].(M)["holder"] = "karen" })},
		{"meta", "metas", `{}`, CodeInvalidArgument, set(func(s M) { push(s, "savings_goals", M{"id": 1, "kind": "goal", "start_date": "2026-10-01"}) })},
		{"funcion", "nada", `{}`, CodeInvalidArgument, seed},
		{"arg", "metas", `{"x":1}`, CodeUnknownField, seed},
		{"periodos", "periodos", `{"months":["2026-10"],"windows":[]}`, CodeInvalidArgument, seed},
		{"vacío", "periodos", `{"months":[]}`, CodeInvalidArgument, seed},
		// PROVISIONAL (docs/estado/bloqueos/S2-ventanas.md): costo o pago antes de la primera ventana o en un hueco.
		{"ventana posterior", "periodos", `{"windows":[{"label":"P2","start":"2026-11-10","end":"2026-12-09"}]}`, CodeInvalidArgument, seed},
		{"hueco con pagos", "periodos", `{"windows":[{"label":"P1","start":"2026-10-10","end":"2026-11-09"},{"label":"P3","start":"2026-12-10","end":"2027-01-09"}]}`, CodeInvalidArgument, seed},
		{"hueco con costo", "periodos", `{"windows":[{"label":"A","start":"2026-10-10","end":"2026-12-09"},{"label":"B","start":"2026-12-11","end":"2027-01-09"}]}`, CodeInvalidArgument, seed},
		{"mes", "presupuesto_mes", `{"month":"2026-1"}`, CodeInvalidArgument, seed},
		{"base", "promedios", `{"month":"2026-10","basis":"x"}`, CodeInvalidArgument, seed},
		{"promedios", "promedios", `{"month":"2026-10"}`, CodeInvalidArgument, seed},
		{"rango", "aplicar_cambio_presupuesto", `{"row_id":1,"month":"2026-09","amount_cents":5}`, CodeOutOfRange, set(func(s M) {
			push(s, "budgets", M{"id": 1, "category": "despensa", "amount_cents": 5, "valid_from": "2026-10-01"})
		})},
		{"renglón", "aplicar_cambio_presupuesto", `{"row_id":7,"month":"2026-10","amount_cents":5}`, CodeInvalidArgument, seed},
		{"forma", "vista_periodo", `{"month":"2026-10"}`, CodeOracleViewShape, set(func(s M) {
			push(s, "persons", M{"code": "otra", "display_name": "Otra", "participation": "external"})
		})},
		{"compra vista", "vista_interes_por_monto", `{"purchase_id":9}`, CodeInvalidArgument, seed},
		{"cuenta", "flujo_quincenal", `{"from":"2026-10-q1","count":0}`, CodeInvalidArgument, seed},
		{"quincena", "flujo_quincenal", `{"from":"2026-10-q3","count":1}`, CodeInvalidArgument, seed},
		{"categoría", "plan_contra_real", `{"from":"2026-10","count":1,"category":"x"}`, CodeInvalidArgument, seed},
		{"calendario", "calendario", `{}`, CodeInvalidArgument, seed},
		{"sensibilidad", "vista_sensibilidad", `{}`, CodeInvalidArgument, seed},
		{"horizonte", "horizonte", `{"start":"2026-10-10"}`, CodeInvalidArgument, seed},
		{"corte", "deuda_viva", `{}`, CodeInvalidArgument, seed},
		{"mensual", "promedio_mensual", `{"total_cents":5,"month_count":0}`, CodeInvalidArgument, seed},
	}
	for _, c := range cases {
		if _, e := Eval(mustJSON(t, c.in), c.fn, []byte(c.args)); e == nil || e.Code != c.code {
			t.Errorf("%s: se esperaba %s y salió %v", c.name, c.code, e)
		}
	}
	// Ramas con éxito: ventanas, genéricas, ocurrencias y abonos.
	s := clone(t, seed)
	push(s, "expenses", cashExp(7, "despensa", 1000, "2026-11-05", "fixed"), cashExp(8, "despensa", 700, "2026-10-20", "upcoming"))
	pur0(s)["payments"].([]any)[0].(M)["paid"] = []any{M{"person": "axel", "amount_cents": 99999, "paid_date": "2026-10-10"}}
	for i, f := range []string{"weekly", "yearly", "monthly", "fortnightly"} {
		push(s, "planned_items", M{"id": 10 + i, "kind": "fixed_expense", "amount_cents": 100, "frequency": f, "start_date": "2026-10-31",
			"end_date": "2027-12-31", "share_rule": M{"parts": []any{M{"person": "karen", "fixed_cents": 10}}, "rest_person": "axel"}})
	}
	w := q(t, s, "periodos", `{"windows":[{"label":"P1","start":"2026-10-10","end":"2026-11-09"},{"label":"P2","start":"2026-11-10","end":"2026-12-09"}]}`)
	eq(t, "gasto del día 5 en la ventana 1", qc(t, w, "[0].cash.cash_expenses_cents"), 1700)
	f := q(t, s, "flujo_quincenal", `{"from":"2026-10-q2","count":4}`)
	eq(t, "abono vigente topado a la parte", qc(t, f, "[0].credit_payments.household_paid_cents"), mc(t, "despensas[0].pagos[0].axel"))
	for _, fn := range []string{"horizonte", "valor_estable", "papel", "proyeccion"} {
		q(t, s, fn, map[string]string{"horizonte": `{"start":"2026-10-10","end":"2027-04-10","months":["2026-10"]}`}[fn])
	}
	push(s, "budgets", M{"id": 1, "category": "despensa", "amount_cents": 5, "valid_from": "2026-10-01", "valid_to": "2027-03-31"})
	for _, sc := range []string{"onward", "this_month"} {
		q(t, s, "aplicar_cambio_presupuesto", `{"row_id":1,"month":"2026-10","amount_cents":9,"scope":"`+sc+`"}`)
	}
	push(s, "savings_goals", M{"id": 1, "kind": "goal", "target_cents": 500, "target_date": "2026-09-30", "start_date": "2026-01-01"})
	eq(t, "meta vencida", leaf(t, q(t, s, "metas", `{}`), "goals[0].state"), "overdue")
}

// TestRangoConteos: accrual_months e interval_count distinguen ausente o null (vale 1) de un 0
// explícito; 0, -1 y 25 son invalid_argument (§2.4, §2.5, §2.7, rango 1–24), como en el oráculo.
func TestRangoConteos(t *testing.T) {
	seed := seedFromModel(t)
	gasto := func(s M, v any, set bool) {
		e := cashExp(7, "despensa", 30000, "2026-12-05", "variable")
		if set {
			e["accrual_months"] = v
		}
		push(s, "expenses", e)
	}
	compra := func(s M, v any, set bool) {
		p := s["credit_purchases"].([]any)[0].(M)
		delete(p, "accrual_months")
		if set {
			p["accrual_months"] = v
		}
	}
	partida := func(key string) func(M, any, bool) {
		return func(s M, v any, set bool) {
			m := M{"id": 9, "kind": "income", "description": "sueldo", "amount_cents": 500, "frequency": "monthly", "start_date": "2026-10-15", "person": "axel"}
			if set {
				m[key] = v
			}
			push(s, "planned_items", m)
		}
	}
	for _, tg := range []struct {
		name, fn, args string
		apply          func(M, any, bool)
	}{
		{"gasto.accrual_months", "periodos", `{"months":["2026-12"]}`, gasto},
		{"compra.accrual_months", "periodos", `{"months":["2026-12"]}`, compra},
		{"partida.interval_count", "proyeccion", `{"from":"2026-10"}`, partida("interval_count")},
		{"partida.accrual_months", "proyeccion", `{"from":"2026-10"}`, partida("accrual_months")},
	} {
		t.Run(tg.name, func(t *testing.T) {
			run := func(v any, set bool) (string, *Error) {
				s := clone(t, seed)
				tg.apply(s, v, set)
				out, e := Eval(mustJSON(t, s), tg.fn, []byte(tg.args))
				if e != nil {
					return "", e
				}
				return string(mustJSON(t, genericJSON(t, mustJSON(t, out)))), nil
			}
			absent, e := run(nil, false)
			if e != nil {
				t.Fatalf("ausente: %v", e)
			}
			for _, v := range []any{nil, 1} {
				got, e := run(v, true)
				if e != nil {
					t.Fatalf("%v: %v", v, e)
				}
				eq(t, sprint(v)+" = ausente", got, absent)
			}
			if _, e := run(24, true); e != nil {
				t.Errorf("24: %v", e)
			}
			for _, v := range []int{0, -1, 25} {
				if _, e := run(v, true); e == nil || e.Code != CodeInvalidArgument {
					t.Errorf("%d: quería invalid_argument y dio %v", v, e)
				}
			}
		})
	}
}

// TestVentanasOrden: ventanas que no van ascendentes y disjuntas (start[i] ≤ end[i−1]) son
// invalid_argument por la regla de orden y no por la de cobertura (PROVISIONAL,
// docs/estado/bloqueos/S2-ventanas.md, pregunta 3; hallazgo M1 de REVISION-S2-r2). Con la
// semilla, todos los eventos caen dentro de alguna ventana pedida, así que sin la regla de orden
// estas consultas darían cifras (las incorrectas que reportó la revisión).
func TestVentanasOrden(t *testing.T) {
	seed := seedFromModel(t)
	w := func(l, a, b string) string { return `{"label":"` + l + `","start":"` + a + `","end":"` + b + `"}` }
	p1, p2, p3 := w("P1", "2026-10-10", "2026-11-09"), w("P2", "2026-11-10", "2026-12-09"), w("P3", "2026-12-10", "2027-01-09")
	cases := map[string][]string{
		"[P2, P1]":          {p2, p1},
		"[P3, P1, P2]":      {p3, p1, p2},
		"traslape P1 + X":   {p1, w("X", "2026-10-25", "2026-12-09")},
		"traslape de 1 día": {p1, w("Y", "2026-11-09", "2026-12-09")},
		"A contiene a B":    {w("A", "2026-10-10", "2027-01-09"), w("B", "2026-11-10", "2026-12-09")},
	}
	for name, ws := range cases {
		_, e := Eval(mustJSON(t, seed), "periodos", []byte(`{"windows":[`+strings.Join(ws, ",")+`]}`))
		if e == nil || e.Code != CodeInvalidArgument || !strings.Contains(e.Detail, "ascendentes y disjuntas") {
			t.Errorf("%s: se esperaba invalid_argument por la regla de orden y salió %v", name, e)
		}
	}
	// Control: contiguas y ascendentes (P1, P2, P3) siguen dando R-45…R-47.
	out := q(t, seed, "periodos", `{"windows":[`+p1+","+p2+","+p3+`]}`)
	for i := 0; i < 3; i++ {
		eq(t, fmt.Sprintf("P%d devengado", i+1), Cents(qc(t, out, fmt.Sprintf("[%d].accrued.total_cents", i))), mc(t, fmt.Sprintf("periodos[%d].devengado.total", i)))
		eq(t, fmt.Sprintf("P%d diferencia", i+1), Cents(qc(t, out, fmt.Sprintf("[%d].cumulative_diff_cents", i))), mc(t, fmt.Sprintf("periodos[%d].diferencia_acumulada", i)))
	}
}
