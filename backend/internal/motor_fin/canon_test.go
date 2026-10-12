package motor_fin

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const fixDir = "oraculo/fixtures/s2"

func needFixture(t testing.TB, rel string) any {
	t.Helper()
	p := filepath.Join(repoRoot(t), fixDir, rel)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("falta %s/%s (autor: agente Python de S2)", fixDir, rel)
	}
	return readJSON(t, filepath.Join(fixDir, rel))
}

// rowsCanon arma las 85 filas de S2.md §6 con sus rutas (los valores se leen de model.json).
func rowsCanon() [][]string {
	var r [][]string
	add := func(p ...string) { r = append(r, p) }
	d0, dq := "despensas[0].", "data.json:financiamiento.dos_quincenas."
	add(d0 + "capital")
	add(d0 + "intereses")
	add(d0 + "total")
	add("parametros.tasa_pct")
	add(dq + "tasa_pct")
	add(dq+"total", dq+"intereses")
	for _, p := range []string{"axel", "karen"} {
		for _, k := range []string{"total", "pago_por_quincena", "intereses"} {
			add(d0 + p + "." + k)
		}
	}
	for k := 0; k < 4; k++ {
		D := fmt.Sprintf("despensas[%d].", k)
		if k > 0 {
			add(D+"capital", D+"intereses", D+"total")
		}
		for i := 0; i < 4; i++ {
			x := fmt.Sprintf("%spagos[%d].", D, i)
			ps := []string{x + "fecha", x + "monto", x + "axel", x + "karen"}
			if k == 0 {
				ps = append(ps, x+"capital_repuesto")
			}
			add(ps...)
		}
	}
	for _, k := range []string{"total", "pago", "capital", "intereses", "conjunto"} {
		add("parametros.karen_fija." + k)
	}
	for _, k := range []string{"mes.total", "mes.axel", "mes.karen", "quincena.total", "quincena.axel", "quincena.karen", "componentes.mitad_despensa", "componentes.papel_contado_mes"} {
		add("valor_estable." + k)
	}
	for k := 0; k < 6; k++ {
		P := fmt.Sprintf("periodos[%d].", k)
		add(P+"devengado.total", P+"efectivo.total", P+"diferencia_acumulada")
	}
	for _, k := range []string{"total_devengado", "total_efectivo", "deuda_viva_10_abr", "deuda_viva_detalle.despensa_3_pendiente", "deuda_viva_detalle.despensa_4",
		"promedio_mensual_efectivo", "promedio_mensual_devengado", "diferencia_final", "pagos_despensa_3_en_horizonte", "papel_contado"} {
		add("horizonte." + k)
	}
	add("credito.minimo_A")
	add("credito.minimo_B")
	add("parametros.limite_credito")
	add("credito.eventos[0].disponible_A", "credito.eventos[0].disponible_B")
	add("credito.comparacion_series.esperado_A")
	add("credito.comparacion_series.esperado_B")
	var ev []string
	for i := 0; i < 17; i++ {
		for _, k := range []string{"fecha", "tipo", "monto", "disponible_A", "disponible_B"} {
			ev = append(ev, fmt.Sprintf("credito.eventos[%d].%s", i, k))
		}
	}
	add(ev...)
	for k := 0; k < 6; k++ {
		P := fmt.Sprintf("papel[%d].", k)
		add(P+"compra", P+"quien", P+"forma", P+"monto")
	}
	for k := 0; k < 5; k++ {
		P := fmt.Sprintf("sensibilidad[%d].", k)
		add(P+"intereses", P+"total", P+"pago_quincenal", P+"axel_total")
	}
	for k := 0; k < 7; k++ {
		P := fmt.Sprintf("interes_por_monto[%d].", k)
		add(P+"intereses", P+"total", P+"pago_quincenal")
	}
	return r
}

func isS2Path(p string) (string, bool) {
	for _, root := range []string{"valor_estable", "horizonte", "papel["} {
		if strings.HasPrefix(p, root) {
			return strings.TrimSuffix(root, "["), true
		}
	}
	if strings.HasPrefix(p, "periodos[") {
		return p[:11], true
	}
	return "", false
}

// checkS3Leaf aplica FORMATO §3 a una hoja de S3: dinero exacto, textos, listas y la razón de R-04.
func checkS3Leaf(t *testing.T, path string, v any) any {
	k := path[strings.LastIndex(path, ".")+1:]
	switch x := v.(type) {
	case string:
		return x
	case []any:
		var out []any
		for _, e := range x {
			out = append(out, int64(toCents(t, e)))
		}
		return out
	}
	if k == "tasa_pct" && !strings.HasPrefix(path, "data.json:") {
		r := new(big.Rat).SetFrac64(100*int64(mc(t, "despensas[0].intereses")), int64(mc(t, "despensas[0].capital")))
		f, _ := r.Float64()
		g, _ := v.(json.Number).Float64()
		eq(t, path+" (binary64 de 100·intereses/capital)", g, f)
		return map[string]any{"num": int64(mc(t, "despensas[0].intereses")), "den": int64(mc(t, "despensas[0].capital"))}
	}
	if k == "tasa_pct" {
		tot, cap := toCents(t, leaf(t, readJSON(t, "datos/semilla/data.json"), "financiamiento.dos_quincenas.total")), mc(t, "despensas[0].capital")
		bp := halfUp(10000*int64(tot-cap), int64(cap)) // HALF_UP(100·num/den) a 2 decimales
		eq(t, path+" (HALF_UP a 2 decimales)", fmt.Sprintf("%d.%02d", bp/100, bp%100), v.(json.Number).String())
		return map[string]any{"num": int64(tot - cap), "den": int64(cap), "decimales": int64(2)}
	}
	return int64(toCents(t, v))
}

// TestCanon (S2-01): las 85 cifras canon; filas S2 contra el motor y filas S3 con los pasos de §5.1.
func TestCanon(t *testing.T) {
	seed := seedFromModel(t)
	in, err := DecodeInput(mustJSON(t, seed))
	if err != nil {
		t.Fatalf("semilla inválida: %v", err)
	}
	views := map[string]any{}
	data := readJSON(t, "datos/semilla/data.json")
	for i, paths := range rowsCanon() {
		id := fmt.Sprintf("R-%02d", i+1)
		t.Run(id, func(t *testing.T) {
			for _, p := range paths {
				src, rel := model(t), p
				if strings.HasPrefix(p, "data.json:") {
					src, rel = data, strings.TrimPrefix(p, "data.json:")
				}
				v := leaf(t, src, rel)
				if root, ok := isS2Path(p); ok {
					if views[root] == nil {
						views[root] = viewFor(t, seed, root)
					}
					sub := strings.TrimPrefix(p, root)
					got := leaf(t, views[root], strings.TrimPrefix(sub, "."))
					if s, ok := v.(string); ok {
						eq(t, p, got, s)
					} else {
						eq(t, p, num(t, got), toCents(t, v))
					}
					continue
				}
				conv := checkS3Leaf(t, p, v)
				checkS3Input(t, &in, seed, p, v)
				ins := needFixture(t, "_datos/insumo-s3.json")
				eq(t, p+" contra insumo-s3.json", fmt.Sprint(leaf(t, ins, "filas."+id).(map[string]any)[p]), fmt.Sprint(conv))
			}
		})
	}
	t.Run("semilla", func(t *testing.T) {
		want := needFixture(t, "_datos/semilla-oraculo.json")
		eq(t, "semilla-oraculo.json = seedFromModel", string(mustJSON(t, genericJSON(t, mustJSON(t, seed)))), string(mustJSON(t, want)))
	})
}

// checkS3Input: paso (2) de §5.1, las hojas de S3 que S2 usa como entrada contra sus entradas.
func checkS3Input(t *testing.T, in *Input, seed map[string]any, p string, v any) {
	var k, i int
	switch {
	case strings.HasPrefix(p, "despensas["):
		k, _ = strconv.Atoi(p[10:11])
		cp := in.CreditPurchases[k]
		rest := p[13:]
		if n, err := fmt.Sscanf(rest, "pagos[%d].", &i); n == 1 && err == nil {
			py := cp.Payments[i]
			f := rest[strings.LastIndex(rest, ".")+1:]
			got := map[string]any{"fecha": py.DueDate.String(), "monto": py.AmountCents, "capital_repuesto": py.PrincipalCents}[f]
			for _, s := range py.Shares {
				if s.Person == f {
					got = s.AmountCents
				}
			}
			if f == "fecha" {
				eq(t, p+" contra la entrada", got, v)
			} else {
				eq(t, p+" contra la entrada", got, toCents(t, v))
			}
			return
		}
		got := map[string]Cents{"capital": cp.PrincipalCents, "intereses": cp.InterestCents, "total": cp.TotalCents}
		for _, s := range cp.Shares {
			got[s.Person+".total"], got[s.Person+".intereses"] = s.TotalCents, s.InterestCents
			got[s.Person+".pago_por_quincena"] = cp.Payments[0].Shares[map[string]int{"axel": 0, "karen": 1}[s.Person]].AmountCents
		}
		eq(t, p+" contra la entrada", got[rest], toCents(t, v))
	case strings.HasPrefix(p, "sensibilidad["):
		k, _ = strconv.Atoi(p[13:14])
		S := fmt.Sprintf("sensibilidad[%d].", k)
		items := fmt.Sprintf(`{"items":[{"label":"despensa","cycle_months":2,"total_cents":%d,"shares":[{"person":"karen","total_cents":%d},{"person":"axel","total_cents":%d}]},{"label":"papel_contado","cycle_months":2,"total_cents":%d,"shares":[{"person":"karen","total_cents":%d}]}]}`,
			mc(t, S+"total"), mc(t, "parametros.karen_fija.total"), mc(t, S+"axel_total"), mc(t, "parametros.papel"), mc(t, "parametros.papel"))
		out := q(t, seed, "vista_sensibilidad", items)
		for _, f := range []string{"costo_mensual_total", "axel_mensual", "karen_mensual"} {
			eq(t, S+f+" (cruce)", qc(t, out, f), mc(t, S+f))
		}
	case strings.HasPrefix(p, "interes_por_monto["):
		k, _ = strconv.Atoi(p[18:19])
		out := q(t, interestInput(t, seed), "vista_interes_por_monto", fmt.Sprintf(`{"purchase_id":%d}`, k+1))
		f := p[strings.LastIndex(p, ".")+1:]
		eq(t, p+" (cruce por devengado y efectivo)", qc(t, out, f), toCents(t, v))
	}
}

// interestInput arma las 7 compras de s2-interes-por-monto (§5.2) desde model.json.
func interestInput(t testing.TB, seed map[string]any) map[string]any {
	type M = map[string]any
	var purs []any
	for k := 0; k < 7; k++ {
		P := fmt.Sprintf("interes_por_monto[%d].", k)
		capC, tot, pq := mc(t, P+"monto"), mc(t, P+"total"), mc(t, P+"pago_quincenal")
		var pays []any
		for i := 0; i < 4; i++ {
			amt, pr := pq, series(capC, 4, i)
			if i == 3 {
				amt = tot - 3*pq
			}
			pays = append(pays, M{"seq": i + 1, "due_date": ms(t, fmt.Sprintf("despensas[0].pagos[%d].fecha", i)), "amount_cents": amt,
				"principal_cents": pr, "interest_cents": amt - pr, "shares": []any{M{"person": "axel", "amount_cents": amt}}, "paid": []any{}})
		}
		purs = append(purs, M{"id": k + 1, "credit_id": 1, "expense_id": nil, "category": "despensa", "accrual_months": 2,
			"purchase_date": "2026-10-10", "principal_cents": capC, "interest_cents": tot - capC, "total_cents": tot, "status": "active",
			"shares":   []any{M{"person": "axel", "principal_cents": capC, "interest_cents": tot - capC, "total_cents": tot}},
			"embedded": []any{}, "payments": pays, "projected": false})
	}
	return M{"params": seed["params"], "persons": seed["persons"], "categories": seed["categories"], "credit_purchases": purs}
}

// TestOraculoFixtures (S2-03): cada fixture de oraculo/fixtures/s2 contra su esperado.
func TestOraculoFixtures(t *testing.T) {
	dir := filepath.Join(repoRoot(t), fixDir)
	names, _ := filepath.Glob(filepath.Join(dir, "s2-*.json"))
	if len(names) == 0 {
		t.Fatalf("falta %s/s2-*.json (autor: agente Python de S2)", fixDir)
	}
	for _, f := range names {
		fx := readJSON(t, filepath.Join(fixDir, filepath.Base(f))).(map[string]any)
		t.Run(fx["caso"].(string), func(t *testing.T) {
			ent := fx["entradas"].(map[string]any)
			datos := resolveData(t, ent["datos"])
			esp := fx["esperado"].(map[string]any)
			for _, c := range ent["consultas"].([]any) {
				qq := c.(map[string]any)
				id := qq["id"].(string)
				out, e := Eval(mustJSON(t, datos), qq["funcion"].(string), mustJSON(t, qq["argumentos"]))
				var got any = map[string]any{"error": "__ninguno__"}
				if e != nil {
					got = map[string]any{"error": e.Code}
				} else {
					got = genericJSON(t, mustJSON(t, out))
				}
				compareExpected(t, id, esp[id], got)
			}
		})
	}
}

func resolveData(t testing.TB, d any) any {
	m, ok := d.(map[string]any)
	if !ok || m["$ref"] == nil {
		return d
	}
	base := readJSON(t, filepath.Join(fixDir, m["$ref"].(string)))
	if p, ok := m["parche"]; ok {
		base = patch(base, p)
	}
	return base
}

func patch(target, p any) any {
	pm, ok := p.(map[string]any)
	if !ok {
		return p
	}
	tm, ok := target.(map[string]any)
	if !ok {
		tm = map[string]any{}
	}
	for k, v := range pm {
		if v == nil {
			delete(tm, k)
		} else {
			tm[k] = patch(tm[k], v)
		}
	}
	return tm
}
