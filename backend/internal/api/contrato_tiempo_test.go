package api_test

// S6-04 con reloj simulado: Idempotency-Key retenida 24 h, página máxima de 200 y capturas atrasadas sin duplicar (K9).

import (
	"context"
	"fmt"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func (k *K) body(t *testing.T, op string) map[string]any {
	t.Helper()
	meta, _ := api.Meta(op)
	b, _ := k.bodyFor(meta)
	m, _ := b.(map[string]any)
	if m == nil {
		t.Fatalf("%s: no pude armar el cuerpo", op)
	}
	return m
}

func (k *K) count(t *testing.T, table string) int {
	t.Helper()
	var n int
	if err := k.s.DB.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func (k *K) mustCategory(t *testing.T) {
	t.Helper()
	meta, _ := api.Meta("createCategory")
	b, _ := k.bodyFor(meta)
	res := k.do(k.c, "POST", "/api/v1/categories", b, "Idempotency-Key", k.key())
	if res.StatusCode != 201 {
		t.Fatalf("categoría de apoyo: %d %s", res.StatusCode, short(res.Body))
	}
	k.ids["category_id"] = int64(jsonOf(res.Body)["id"].(float64))
}

// TestContrato_Idempotencia24h: la clave se retiene 24 h; pasado el plazo ya no repite la respuesta (S6-04, K6).
func TestContrato_Idempotencia24h(t *testing.T) {
	k := newK(t)
	body := k.body(t, "createIncome")
	key := "contrato-24h-ingreso-0001"
	post := func() *testutil.Response { return k.do(k.c, "POST", "/api/v1/incomes", body, "Idempotency-Key", key) }

	t0 := k.s.Clock.Now()
	r1 := post()
	if r1.StatusCode != 201 {
		t.Fatalf("alta: %d %s", r1.StatusCode, short(r1.Body))
	}
	id1 := jsonOf(r1.Body)["id"]
	// 1 h después: repetición idéntica
	k.s.Clock.Advance(1 * time.Hour)
	if r := post(); r.StatusCode != 201 || r.Header.Get("Idempotent-Replayed") != "true" || jsonOf(r.Body)["id"] != id1 {
		t.Errorf("a la 1 h debe repetir la respuesta guardada: %d replayed=%q", r.StatusCode, r.Header.Get("Idempotent-Replayed"))
	}
	// 23 h 59 min desde el alta: sigue retenida
	k.s.Clock.Set(t0.Add(23*time.Hour + 59*time.Minute))
	if r := post(); r.Header.Get("Idempotent-Replayed") != "true" || jsonOf(r.Body)["id"] != id1 {
		t.Errorf("a las 23 h 59 min debe seguir retenida: %d replayed=%q", r.StatusCode, r.Header.Get("Idempotent-Replayed"))
	}
	if n := k.count(t, "income"); n != 1 {
		t.Errorf("hasta aquí debe haber 1 ingreso y hay %d", n)
	}
	// 24 h 1 min: ya no se retiene; el servidor la trata como una petición nueva (y por eso el cliente no debe depender de ella)
	k.s.Clock.Set(t0.Add(24*time.Hour + time.Minute))
	r := post()
	if r.Header.Get("Idempotent-Replayed") == "true" {
		t.Errorf("a las 24 h 1 min la clave ya no debe repetir la respuesta guardada")
	}
	if r.StatusCode != 201 || jsonOf(r.Body)["id"] == id1 {
		t.Errorf("pasadas las 24 h la petición se ejecuta de nuevo: %d %s", r.StatusCode, short(r.Body))
	}
	// el job idempotency_purge borra lo vencido
	n, err := dominio.PurgeIdempotency(context.Background(), k.s.DB, k.s.Clock.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	// el job de purga borra las claves de más de 24 h aunque nadie las repita
	body["description"] = "Otro ingreso"
	r2 := k.do(k.c, "POST", "/api/v1/incomes", body, "Idempotency-Key", "contrato-24h-ingreso-0002")
	if r2.StatusCode != 201 {
		t.Fatalf("alta 2: %d %s", r2.StatusCode, short(r2.Body))
	}
	k.s.Clock.Advance(25 * time.Hour)
	n2, err := dominio.PurgeIdempotency(context.Background(), k.s.DB, k.s.Clock.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if left := k.count(t, "idempotency_key"); left != 0 || n2 < 1 {
		t.Errorf("la purga debía borrar las claves vencidas: borró %d y quedan %d", n2, left)
	}
	t.Logf("S6-04: clave retenida 23 h 59 min; a las 24 h 1 min ya no; la purga borró %d vencidas y quedan 0 (antes borró %d)", n2, n)
}

// TestContrato_CapturasAtrasadas (K9): una captura hecha sin conexión y reenviada pasadas las 24 h no duplica gastos ni pesos.
func TestContrato_CapturasAtrasadas(t *testing.T) {
	k := newK(t)
	k.mustCategory(t)
	t0 := k.s.Clock.Now()
	type cap struct {
		path, key string
		body      map[string]any
		id        any
	}
	var caps []*cap
	for i := 0; i < 3; i++ {
		b := k.body(t, "createExpense")
		b["captured_at"] = t0.Add(-48 * time.Hour).Format(time.RFC3339)
		b["expense_date"] = t0.Add(-48 * time.Hour).Format("2006-01-02")
		b["amount_cents"] = 10000 + i*100
		caps = append(caps, &cap{"/api/v1/expenses", fmt.Sprintf("contrato-k9-gasto-%04d", i), b, nil})
	}
	for i := 0; i < 2; i++ {
		b := k.body(t, "createWeightEntry")
		b["captured_at"] = t0.Add(-48 * time.Hour).Format(time.RFC3339)
		b["measured_date"] = t0.Add(-time.Duration(48+24*i) * time.Hour).Format("2006-01-02")
		b["weight_g"] = 75000 + i*100
		caps = append(caps, &cap{"/api/v1/weight-entries", fmt.Sprintf("contrato-k9-peso-%04d", i), b, nil})
	}
	// primer envío: capturas con 2 días de atraso se aceptan
	for _, c := range caps {
		r := k.do(k.c, "POST", c.path, c.body, "Idempotency-Key", c.key)
		if r.StatusCode != 201 && r.StatusCode != 200 {
			t.Fatalf("%s primer envío: %d %s", c.path, r.StatusCode, short(r.Body))
		}
		c.id = jsonOf(r.Body)["id"]
	}
	ge, pe := k.count(t, "expense"), k.count(t, "weight_entry")
	if ge != 3 || pe != 2 {
		t.Fatalf("tras el primer envío: %d gastos, %d pesos", ge, pe)
	}
	// la cola sin conexión reenvía todo 25 h y 49 h después, antes y después de la purga
	for _, horas := range []time.Duration{25, 49} {
		k.s.Clock.Set(t0.Add(horas * time.Hour))
		if horas == 49 {
			if _, err := dominio.PurgeIdempotency(context.Background(), k.s.DB, k.s.Clock.Now().Add(-24*time.Hour)); err != nil {
				t.Fatal(err)
			}
		}
		for _, c := range caps {
			r := k.do(k.c, "POST", c.path, c.body, "Idempotency-Key", c.key)
			if r.StatusCode != 200 || r.Header.Get("Idempotent-Replayed") != "true" || jsonOf(r.Body)["id"] != c.id {
				t.Errorf("%s a las %d h: %d replayed=%q id=%v (era %v): %s", c.path, horas, r.StatusCode, r.Header.Get("Idempotent-Replayed"), jsonOf(r.Body)["id"], c.id, short(r.Body))
			}
		}
		if g, p := k.count(t, "expense"), k.count(t, "weight_entry"); g != ge || p != pe {
			t.Errorf("a las %d h hay %d gastos y %d pesos, debían seguir en %d y %d (duplicó)", horas, g, p, ge, pe)
		}
	}
	t.Logf("S6-04/K9: 3 gastos y 2 pesos capturados con 48 h de atraso; reenviados a las 25 h y a las 49 h (con purga) no se duplican")
}

// TestContrato_PaginaMaxima: limit 200 devuelve a lo más 200 con cursor; 201 y 0 son 400 (C1 §5).
func TestContrato_PaginaMaxima(t *testing.T) {
	k := newK(t)
	body := k.body(t, "createIncome")
	for i := 0; i < 205; i++ {
		body["description"] = fmt.Sprintf("Ingreso %d", i)
		if r := k.c.Raw("POST", "/api/v1/incomes", body, "Idempotency-Key", fmt.Sprintf("contrato-pagina-%05d", i)); r.StatusCode != 201 {
			t.Fatalf("alta %d: %d %s", i, r.StatusCode, short(r.Body))
		}
	}
	type page struct {
		Items      []map[string]any `json:"items"`
		NextCursor *string          `json:"next_cursor"`
	}
	get := func(q string) (int, page) {
		r := k.do(k.c, "GET", "/api/v1/incomes"+q, nil)
		var p page
		if r.StatusCode == 200 {
			r.JSON(t, &p)
		}
		return r.StatusCode, p
	}
	code, p1 := get("?limit=200")
	if code != 200 || len(p1.Items) != 200 || p1.NextCursor == nil {
		t.Fatalf("limit=200: %d, %d renglones, cursor %v", code, len(p1.Items), p1.NextCursor)
	}
	code, p2 := get("?limit=200&cursor=" + *p1.NextCursor)
	if code != 200 || len(p2.Items) != 5 || p2.NextCursor != nil {
		t.Errorf("segunda página: %d, %d renglones (se esperaban 5), cursor %v", code, len(p2.Items), p2.NextCursor)
	}
	seen := map[any]bool{}
	for _, it := range append(p1.Items, p2.Items...) {
		seen[it["id"]] = true
	}
	if len(seen) != 205 {
		t.Errorf("entre las dos páginas deben salir 205 ingresos distintos y salen %d", len(seen))
	}
	for _, q := range []string{"?limit=201", "?limit=0", "?limit=-1", "?limit=abc"} {
		if c, _ := get(q); c != 400 {
			t.Errorf("%s dio %d, se esperaba 400", q, c)
		}
	}
	if c, p := get(""); c != 200 || len(p.Items) != 50 {
		t.Errorf("sin limit deben salir 50 y salen %d", len(p.Items))
	}
	t.Logf("S6-04: página máxima 200 (205 ingresos = 200 + 5), 201/0/-1/abc = 400, por omisión 50")
}
