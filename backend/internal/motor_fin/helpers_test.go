package motor_fin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// Las pruebas leen oraculo/model.json en tiempo de prueba (S2 §3: nada pegado en el código).

func repoRoot(t testing.TB) string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "oraculo", "model.json")); err == nil {
			return dir
		}
		p := filepath.Dir(dir)
		if p == dir {
			t.Fatal("no encuentro oraculo/model.json subiendo desde el directorio de prueba")
		}
		dir = p
	}
}

func genericJSON(t testing.TB, b []byte) any {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	return v
}

func readJSON(t testing.TB, rel string) any {
	b, err := os.ReadFile(filepath.Join(repoRoot(t), rel))
	if err != nil {
		t.Fatalf("falta %s: %v", rel, err)
	}
	return genericJSON(t, b)
}

var modelCache any

func model(t testing.TB) any {
	if modelCache == nil {
		modelCache = readJSON(t, "oraculo/model.json")
	}
	return modelCache
}

// leaf sigue una ruta «a.b[0].c» (FORMATO §3); ok=false si no existe.
func leafOK(root any, path string) (any, bool) {
	cur := root
	for _, part := range strings.Split(strings.ReplaceAll(path, "[", ".["), ".") {
		if part == "" {
			continue
		}
		if part[0] == '[' {
			i, _ := strconv.Atoi(strings.TrimSuffix(part[1:], "]"))
			l, ok := cur.([]any)
			if !ok || i >= len(l) {
				return nil, false
			}
			cur = l[i]
			continue
		}
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		if cur, ok = m[part]; !ok {
			return nil, false
		}
	}
	return cur, true
}

func leaf(t testing.TB, root any, path string) any {
	v, ok := leafOK(root, path)
	if !ok {
		t.Fatalf("no existe la ruta %s", path)
	}
	return v
}

// toCents: int(valor × 100) exacto (FORMATO §3); fracción de centavo = error de formato.
func toCents(t testing.TB, v any) Cents {
	n, ok := v.(json.Number)
	if !ok {
		t.Fatalf("no es número: %v", v)
	}
	r, ok := new(big.Rat).SetString(n.String())
	if !ok {
		t.Fatalf("número inválido %s", n)
	}
	r.Mul(r, big.NewRat(100, 1))
	if !r.IsInt() {
		t.Fatalf("%s no es exacto en centavos", n)
	}
	return Cents(r.Num().Int64())
}

func num(t testing.TB, v any) int64 {
	n, ok := v.(json.Number)
	if !ok {
		t.Fatalf("no es entero: %#v", v)
	}
	i, err := n.Int64()
	if err != nil {
		t.Fatalf("no es entero: %v", n)
	}
	return i
}

func mc(t testing.TB, path string) Cents  { return toCents(t, leaf(t, model(t), path)) }
func ms(t testing.TB, path string) string { return leaf(t, model(t), path).(string) }

// seedFromModel arma la semilla con la receta exacta de S2-interfaz §5.3, solo con hojas de S3.
func seedFromModel(t testing.TB) map[string]any {
	type M = map[string]any
	pp := []string{"axel", "karen"}
	var exps, purs []any
	for k := 0; k < 4; k++ {
		D := fmt.Sprintf("despensas[%d]", k)
		shares, pshares, emb, pays := []any{}, []any{}, []any{}, []any{}
		for _, p := range pp {
			shares = append(shares, M{"person": p, "amount_cents": mc(t, D+"."+p+".capital")})
			pshares = append(pshares, M{"person": p, "principal_cents": mc(t, D+"."+p+".capital"),
				"interest_cents": mc(t, D+"."+p+".intereses"), "total_cents": mc(t, D+"."+p+".total")})
			if v := mc(t, D+"."+p+".papel_dentro"); v > 0 {
				emb = append(emb, M{"category": "papel", "person": p, "amount_cents": v})
			}
		}
		for i := 0; i < 4; i++ {
			x := fmt.Sprintf("%s.pagos[%d]", D, i)
			ps := []any{}
			for _, p := range pp {
				ps = append(ps, M{"person": p, "amount_cents": mc(t, x+"."+p)})
			}
			pays = append(pays, M{"seq": i + 1, "due_date": ms(t, x+".fecha"), "amount_cents": mc(t, x+".monto"),
				"principal_cents": mc(t, x+".capital_repuesto"), "interest_cents": mc(t, x+".monto") - mc(t, x+".capital_repuesto"),
				"shares": ps, "paid": []any{}})
		}
		exps = append(exps, M{"id": k + 1, "kind": "variable", "category": "despensa", "amount_cents": mc(t, D+".capital"),
			"expense_date": ms(t, D+".fecha"), "accrual_months": 2, "payment_method": "credit", "paid_by": nil,
			"shares": shares, "credit_purchase_id": k + 1, "planned_item_id": nil})
		purs = append(purs, M{"id": k + 1, "credit_id": 1, "expense_id": k + 1, "purchase_date": ms(t, D+".fecha"),
			"principal_cents": mc(t, D+".capital"), "interest_cents": mc(t, D+".intereses"), "total_cents": mc(t, D+".total"),
			"status": "active", "shares": pshares, "embedded": emb, "payments": pays, "projected": false})
	}
	papel := mc(t, "parametros.papel")
	for id, i := range []int{2, 4} {
		who := strings.ToLower(ms(t, fmt.Sprintf("papel[%d].quien", i)))
		exps = append(exps, M{"id": id + 5, "kind": "variable", "category": "papel", "amount_cents": papel,
			"expense_date": ms(t, fmt.Sprintf("papel[%d].compra", i)), "accrual_months": 1, "payment_method": "cash",
			"paid_by": who, "shares": []any{M{"person": who, "amount_cents": papel}}, "credit_purchase_id": nil, "planned_item_id": nil})
	}
	return M{
		"params": M{"today": "2026-10-11", "holder": "axel", "groceries_category": "despensa", "paper_category": "papel", "repeat_cadence_enabled": false},
		"persons": []any{M{"code": "axel", "display_name": "Axel", "participation": "holder", "pays_share_directly": false},
			M{"code": "karen", "display_name": "Karen", "participation": "external", "pays_share_directly": true}},
		"categories": []any{M{"code": "despensa", "name": "Despensa", "kind": "variable", "is_active": true},
			M{"code": "papel", "name": "Papel", "kind": "variable", "is_active": true}},
		"expenses": exps, "credit_purchases": purs, "incomes": []any{}, "planned_items": []any{}, "budgets": []any{},
		"savings_goals": []any{}, "savings_movements": []any{}, "cadence": M{"purchases": []any{}, "expenses": []any{}},
		"stable_items": []any{
			M{"label": "despensa", "cycle_months": 2, "total_cents": mc(t, "despensas[1].total"), "shares": []any{
				M{"person": "karen", "total_cents": mc(t, "despensas[1].karen.total")}, M{"person": "axel", "total_cents": mc(t, "despensas[1].axel.total")}}},
			M{"label": "papel_contado", "cycle_months": 2, "total_cents": papel, "shares": []any{M{"person": "karen", "total_cents": papel}}}},
	}
}

func mustJSON(t testing.TB, v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// q evalúa una consulta y devuelve la salida genérica (números exactos).
func q(t testing.TB, in map[string]any, fn, args string) any {
	t.Helper()
	out, e := Eval(mustJSON(t, in), fn, []byte(args))
	if e != nil {
		t.Fatalf("%s %s: %v", fn, args, e)
	}
	return genericJSON(t, mustJSON(t, out))
}

func qc(t testing.TB, out any, path string) int64 { t.Helper(); return num(t, leaf(t, out, path)) }

func eq(t testing.TB, what string, got, want any) {
	t.Helper()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("%s: motor=%v esperado=%v", what, got, want)
	}
}

// walkLeaves enumera las hojas numéricas de un nodo de model.json (rutas ordenadas).
func walkLeaves(v any, path string, f func(path string, n json.Number)) {
	switch x := v.(type) {
	case map[string]any:
		ks := make([]string, 0, len(x))
		for k := range x {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		for _, k := range ks {
			walkLeaves(x[k], path+"."+k, f)
		}
	case []any:
		for i, e := range x {
			walkLeaves(e, fmt.Sprintf("%s[%d]", path, i), f)
		}
	case json.Number:
		f(path, x)
	}
}

// viewFor calcula la vista de §3.14 que reproduce una raíz S2 de model.json.
func viewFor(t testing.TB, seed map[string]any, root string) any {
	switch {
	case strings.HasPrefix(root, "periodos["):
		k, _ := strconv.Atoi(root[9 : len(root)-1])
		m := (Month{2026, 10}).Add(k)
		return q(t, seed, "vista_periodo", `{"month":"`+m.String()+`"}`)
	case root == "valor_estable":
		return q(t, seed, "vista_valor_estable", `{}`)
	case root == "papel":
		return q(t, seed, "vista_papel", `{}`)
	}
	return q(t, seed, "vista_horizonte", fmt.Sprintf(`{"start":%q,"end":%q,"months":["2026-10","2026-11","2026-12","2027-01","2027-02","2027-03"]}`,
		ms(t, "horizonte.inicio"), ms(t, "horizonte.fin")))
}

// TestOraculoHojas (S2-02): las 159 hojas numéricas de las ramas de S2 contra el motor.
func TestOraculoHojas(t *testing.T) {
	seed := seedFromModel(t)
	n := 0
	for _, root := range []string{"periodos[0]", "periodos[1]", "periodos[2]", "periodos[3]", "periodos[4]", "periodos[5]", "valor_estable", "horizonte", "papel"} {
		view := viewFor(t, seed, root)
		walkLeaves(leaf(t, model(t), root), "", func(rel string, v json.Number) {
			n++
			t.Run(root+rel, func(t *testing.T) {
				want := toCents(t, v)
				if k := rel[strings.LastIndex(rel, ".")+1:]; k == "n_pagos" || k == "despensa" || k == "despensa_que_cubre" || k == "paquete" {
					want = Cents(num(t, v))
				}
				got, ok := leafOK(view, strings.TrimPrefix(rel, "."))
				if !ok {
					t.Fatalf("%s%s: falta en la salida del motor (model.json=%v)", root, rel, v)
				}
				eq(t, root+rel, num(t, got), want)
			})
		})
	}
	if n != 159 {
		t.Errorf("hojas numéricas de S2 en model.json: %d, se esperaban 159", n)
	}
}

// compareExpected compara esperado con la salida (subconjunto de claves; listas completas).
func compareExpected(t testing.TB, path string, want, got any) int {
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			t.Errorf("%s: tipo distinto (motor=%v)", path, got)
			return 1
		}
		d := 0
		for k, v := range w {
			gv, ok := g[k]
			if !ok {
				t.Errorf("%s.%s: __falta__", path, k)
				d++
				continue
			}
			d += compareExpected(t, path+"."+k, v, gv)
		}
		return d
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			t.Errorf("%s: largo o tipo distinto (motor=%v)", path, got)
			return 1
		}
		d := 0
		for i := range w {
			d += compareExpected(t, fmt.Sprintf("%s[%d]", path, i), w[i], g[i])
		}
		return d
	}
	if !reflect.DeepEqual(fmt.Sprint(want), fmt.Sprint(got)) {
		t.Errorf("%s: esperado=%v motor=%v", path, want, got)
		return 1
	}
	return 0
}
