package avisos

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"finanzas-hogar/internal/api/ajustes"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func a9(d dominio.Date, hh, mm int) time.Time { return testutil.MX(d.Year, d.Month, d.Day, hh, mm) }

// TestAlerta_Matriz (S7-08): fechas × alertas con disparos esperados y no esperados, sin duplicados.
func TestAlerta_Matriz(t *testing.T) {
	ctx := context.Background()
	db, fx := testutil.NewDB(t)
	destinos(t, db, fx.AxelPersonID, fx.XimenaPersonID)
	g := fecha(t, "2026-10-17")
	f := &fuenteFija{FuenteBD: FuenteBD{db},
		pagos:    []Pago{{42, fecha(t, "2026-10-30"), 97375, "Kueski"}, {43, fecha(t, "2027-02-28"), 97375, "Kueski"}},
		usos:     []Uso{{"category:5", "Despensa", 100000, 80000}, {"total", "Total del mes", 500000, 510000}},
		despensa: &g, viejos: 3}
	b := &buzon{}
	m := &Motor{DB: db, Fuente: f, Enviador: b}
	type caso struct {
		dia    string
		pago   int
		precio int
	}
	casos := []caso{{"2026-10-15", 0, 0}, {"2026-10-16", 0, 2}, {"2026-10-17", 0, 0}, {"2026-10-28", 0, 0}, {"2026-10-29", 2, 0},
		{"2026-10-30", 2, 0}, {"2026-10-31", 0, 0}, {"2027-02-26", 0, 0}, {"2027-02-27", 2, 0}, {"2027-02-28", 2, 0}, {"2027-03-01", 0, 0}}
	difs := 0
	for _, c := range casos {
		d := fecha(t, c.dia)
		if err := m.TickAlertas(ctx, a9(d, 8, 59)); err != nil { // antes de alerts.send_time: nada
			t.Fatal(err)
		}
		if len(b.msgs) != 0 {
			t.Fatalf("%s 08:59: %d mensajes antes de la hora", c.dia, len(b.msgs))
		}
		for range 2 { // reintento: no duplica
			if err := m.TickAlertas(ctx, a9(d, 9, 0)); err != nil {
				t.Fatal(err)
			}
		}
		p, s := len(b.con("PAGO PRÓXIMO")), len(b.con("PRECIOS VIEJOS"))
		ok := p == c.pago && s == c.precio
		if !ok {
			difs++
		}
		t.Logf("S7-08 %s pago=%d (esperado %d) precios=%d (esperado %d) %v", c.dia, p, c.pago, s, c.precio, ok)
		b.msgs = nil
	}
	if difs != 0 {
		t.Fatalf("%d diferencias en la matriz", difs)
	}
	// Presupuesto: Despensa al 80 % sale una vez; el total rebasó el 100 %: sale el 100 y el 80 queda skipped.
	var n80, n100, skipped int
	db.QueryRow(`SELECT count(*) FROM notification_delivery WHERE dedupe_key = 'budget_threshold:2026-10:category:5:80' AND status = 'sent'`).Scan(&n80)
	db.QueryRow(`SELECT count(*) FROM notification_delivery WHERE dedupe_key = 'budget_threshold:2026-10:total:100' AND status = 'sent'`).Scan(&n100)
	db.QueryRow(`SELECT count(*) FROM notification_delivery WHERE dedupe_key = 'budget_threshold:2026-10:total:80' AND status = 'skipped'`).Scan(&skipped)
	if n80 != 2 || n100 != 2 || skipped != 2 {
		t.Fatalf("presupuesto: 80=%d 100=%d skipped=%d", n80, n100, skipped)
	}
	f.usos[0].UsadoCents = 100000 // Despensa llega al 100 %
	b.msgs = nil
	m.TickAlertas(ctx, a9(fecha(t, "2026-10-20"), 9, 0))
	if got := b.con("PRESUPUESTO AL 100 %"); len(got) != 2 || !strings.Contains(got[0].text, "Despensa lleva $1,000.00 de $1,000.00") {
		t.Fatalf("100 %%: %+v", got)
	}
	t.Logf("S7-08 presupuesto: categoría 80 %% y 100 %% una vez cada uno; total 100 %% una vez (80 %% anotado sin mandar)")
}

// TestAlerta_Plan (S7-08 caso 4): 21 casillas armadas, una sin receta ni «fuera», «fuera» cuenta como armada.
func TestAlerta_Plan(t *testing.T) {
	ctx := context.Background()
	db, fx := testutil.NewDB(t)
	destinos(t, db, fx.AxelPersonID, fx.XimenaPersonID)
	b := &buzon{}
	m := &Motor{DB: db, Fuente: FuenteBD{db}, Enviador: b}
	lunes := fecha(t, "2026-10-19")
	for i := range 7 {
		for _, p := range []int64{fx.AxelPersonID, fx.XimenaPersonID} {
			if _, err := db.Exec(`INSERT INTO meal_plan_entry (plan_date, meal_slot_id, person_id, is_eating_out)
 SELECT ?, id, ?, 1 FROM meal_slot WHERE is_optional = 0`, lunes.AddDays(i).String(), p); err != nil {
				t.Fatal(err)
			}
		}
	}
	n, total, err := m.Fuente.Casillas(ctx, lunes)
	if err != nil || n != 0 || total != 21 {
		t.Fatalf("armada: %d de %d %v", n, total, err)
	}
	m.TickAlertas(ctx, a9(fecha(t, "2026-10-18"), 9, 0))
	if got := b.con("PLAN DE COMIDAS"); len(got) != 0 {
		t.Fatal("semana armada no debe avisar")
	}
	// Una casilla de Ximena sin receta ni «fuera» (el miércoles en la cena).
	db.Exec(`DELETE FROM meal_plan_entry WHERE person_id = ? AND plan_date = '2026-10-21' AND meal_slot_id = (SELECT id FROM meal_slot WHERE code = 'dinner')`, fx.XimenaPersonID)
	db.Exec(`DELETE FROM notification_delivery`)
	m.TickAlertas(ctx, a9(fecha(t, "2026-10-18"), 9, 0))
	m.TickAlertas(ctx, a9(fecha(t, "2026-10-18"), 9, 1))
	got := b.con("PLAN DE COMIDAS SIN ARMAR")
	if len(got) != 2 || !strings.Contains(got[0].text, "1 de 21 casillas sin armar") {
		t.Fatalf("%+v", got)
	}
	m.TickAlertas(ctx, a9(fecha(t, "2026-10-17"), 9, 0)) // sábado: no es el día
	if len(b.con("PLAN")) != 0 {
		t.Fatal("solo el día configurado")
	}
	t.Logf("S7-08 (4): 21 armadas → 0 avisos; 1 sin receta ni «fuera» → aviso «1 de 21»; «fuera» cuenta como armada; sábado → 0")
}

// TestAlerta_Configuracion (S7-08): con updateAlertSettings la matriz cambia y al restablecer vuelve.
func TestAlerta_Configuracion(t *testing.T) {
	ctx := context.Background()
	s := testutil.NewServer(t, ajustes.Register)
	destinos(t, s.DB, s.Fx.AxelPersonID, s.Fx.XimenaPersonID)
	g := fecha(t, "2026-10-17")
	f := &fuenteFija{FuenteBD: FuenteBD{s.DB}, pagos: []Pago{{42, fecha(t, "2026-10-30"), 97375, "Kueski"}, {43, fecha(t, "2027-02-28"), 97375, "Kueski"}},
		usos: []Uso{{"category:5", "Despensa", 100000, 75000}}, despensa: &g, viejos: 2}
	b := &buzon{}
	m := &Motor{DB: s.DB, Fuente: f, Enviador: b}
	c := s.SessionFor(t, testutil.AxelUser)
	res := c.Do("PATCH", "/api/v1/settings/alerts", map[string]any{"parameters": map[string]any{
		"alerts.send_time": "10:15", "alerts.stale_prices_advance_days": 2, "alerts.meal_plan_missing_weekday": 5,
		"alerts.budget_warning_pct": 70, "alerts.payment_due_days_before": []int{2}}})
	if res.StatusCode != 200 {
		t.Fatalf("updateAlertSettings: %d %s", res.StatusCode, res.Body)
	}
	tick := func(dia string, hh, mm int) { m.TickAlertas(ctx, a9(fecha(t, dia), hh, mm)) }
	tick("2026-10-28", 9, 0) // hora vieja: nada
	if len(b.msgs) != 0 {
		t.Fatalf("a las 09:00 ya no debe salir nada: %d", len(b.msgs))
	}
	tick("2026-10-28", 10, 15)
	if got := b.con("PAGO PRÓXIMO"); len(got) != 2 || !strings.Contains(got[0].text, "En 2 días") {
		t.Fatalf("pago con [2]: %+v", got)
	}
	if got := b.con("PRESUPUESTO AL 70 %"); len(got) != 2 {
		t.Fatalf("umbral 70: %d", len(got))
	}
	tick("2026-10-29", 10, 15)
	if len(b.con("PAGO PRÓXIMO")) != 0 {
		t.Fatal("con [2] el día anterior ya no avisa")
	}
	tick("2026-10-15", 10, 15)
	tick("2026-10-16", 10, 15)
	if got := b.con("PRECIOS VIEJOS"); len(got) != 2 {
		t.Fatalf("precios con anticipación 2: %d", len(got))
	}
	if got := b.con("PLAN DE COMIDAS"); len(got) != 2 { // el 16-oct-2026 es viernes
		t.Fatalf("plan el viernes 16: %d", len(got))
	}
	tick("2026-10-23", 10, 15) // viernes
	tick("2026-10-25", 10, 15) // domingo
	if got := b.con("PLAN DE COMIDAS"); len(got) != 2 || !strings.Contains(got[0].text, "26-oct-2026") {
		t.Fatalf("plan el viernes y no el domingo: %+v", got)
	}
	keys := []string{"alerts.send_time", "alerts.stale_prices_advance_days", "alerts.meal_plan_missing_weekday", "alerts.budget_warning_pct", "alerts.budget_limit_pct", "alerts.payment_due_days_before"}
	if res := c.Do("PATCH", "/api/v1/settings/alerts", map[string]any{"reset_parameters": keys}); res.StatusCode != 200 {
		t.Fatalf("restablecer: %d", res.StatusCode)
	}
	a, err := LeerAjustes(ctx, s.DB)
	if err != nil || a.AlertTime != "09:00" || a.StaleAdvance != 1 || a.MealDay != 7 || a.Warn != 80 || a.Limit != 100 || fmt.Sprint(a.DaysBefore) != "[1 0]" {
		t.Fatalf("restablecidos: %+v %v", a, err)
	}
	b.msgs = nil
	for _, d := range []string{"2027-02-26", "2027-02-27", "2027-02-28"} {
		tick(d, 9, 0)
	}
	if got := b.con("PAGO PRÓXIMO"); len(got) != 4 {
		t.Fatalf("matriz original tras restablecer: %d", len(got))
	}
	// params_json ya no trae hora, día ni anticipación (K47).
	rows, _ := s.DB.Query(`SELECT code, params_json FROM notification`)
	defer rows.Close()
	for rows.Next() {
		var code, raw string
		rows.Scan(&code, &raw)
		var p map[string]any
		json.Unmarshal([]byte(raw), &p)
		for _, k := range []string{"send_time", "advance_days", "weekday", "day_of_month", "days_before", "thresholds_pct"} {
			if _, ok := p[k]; ok {
				t.Fatalf("%s.params_json trae %s", code, k)
			}
		}
	}
	t.Logf("S7-08 configuración: send_time 10:15, anticipación 2, plan el viernes, umbral 70, días [2] → la matriz cambia; restablecer → 1, 09:00, 7, 80, 100, [1 0]; params_json sin send_time/advance_days/weekday")
}
