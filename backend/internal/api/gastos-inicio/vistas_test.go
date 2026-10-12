package gastos_test

import (
	"strings"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/testutil"
)

func view(t *testing.T, c *testutil.Client, path string, code int) api.ViewPage {
	t.Helper()
	res := c.Do("GET", "/api/v1/views/"+path, nil) // valida contra openapi.yaml
	want(t, res, code)
	var pg api.ViewPage
	if code == 200 {
		res.JSON(t, &pg)
	}
	return pg
}

func groupOf(pg api.ViewPage, key string) *api.ViewGroup {
	for i := range pg.Groups {
		if pg.Groups[i].Key == key {
			return &pg.Groups[i]
		}
	}
	return nil
}

func TestOp_queryView(t *testing.T) {
	s, c, f := inicioFx(t) // gastos de 1,000.00 (Despensa, 50/50 con Karen) y 500.00 (Varios, hogar)
	newExpense(t, c, map[string]any{"category_id": f.varios, "amount_cents": 20000, "expense_date": "2026-09-30", "payment_method": "debit", "accrual_months": 1, "description": "Ferretería"}, 503)

	// Tabla de gastos: importe principal = parte del hogar; sin periodo no se filtra por fecha (bloqueo §16).
	pg := view(t, c, "expense?mode=table", 200)
	if pg.Mode != "table" || len(pg.Items) != 3 || len(pg.Groups) != 0 || pg.NextCursor.Valid || pg.Items[0].Record["expense_date"] != "2026-10-06" {
		t.Fatalf("tabla: %+v", pg)
	}
	if a := pg.Items[0].AmountCents; a == nil || *a != 50000 || pg.Items[0].Title != "Despensa" {
		t.Fatalf("parte del hogar: %+v", pg.Items[0])
	}
	if _, ok := pg.Items[0].Record["idempotency_key"]; ok {
		t.Error("idempotency_key no sale en la vista")
	}
	// Tablero por categoría: conteo, total del hogar y participación (HALF_UP en basis points).
	pg = view(t, c, "expense?mode=board&group_by=category_id&period=2026-10", 200)
	g := groupOf(pg, itoa(f.despensa))
	if len(pg.Items) != 2 || len(pg.Groups) != 2 || g == nil || g.Label != "Despensa" || g.Count != 1 || *g.TotalCents != 50000 || *g.ShareBp != 5000 {
		t.Fatalf("tablero categoría: %+v", pg)
	}
	// Tablero por tipo: las columnas del enumerado salen aunque estén vacías.
	pg = view(t, c, "expense?mode=board&group_by=payment_method", 200)
	if len(pg.Groups) != 4 || groupOf(pg, "transfer").Count != 0 || groupOf(pg, "debit").Label != "Débito" || groupOf(pg, "cash").Count != 2 {
		t.Fatalf("tablero medio de pago: %+v", pg.Groups)
	}
	// Calendario: grupos por día con etiqueta en español.
	pg = view(t, c, "expense?mode=calendar&date_field=expense_date&q=ferreteria", 200)
	if len(pg.Items) != 1 || len(pg.Groups) != 1 || pg.Groups[0].Key != "2026-09-30" || pg.Groups[0].Label != "miércoles 30 de septiembre de 2026" || *pg.Items[0].GroupKey != "2026-09-30" {
		t.Fatalf("calendario: %+v", pg)
	}
	// Galería con foto de ticket.
	ins(t, s.DB, `INSERT INTO attachment (expense_id, kind, file_path, mime_type, size_bytes, sha256, ocr_status) VALUES ((SELECT min(id) FROM expense), 'ticket_photo', 'x/a.jpg', 'image/jpeg', 10, lower(hex(zeroblob(32))), 'none')`)
	pg = view(t, c, "expense?mode=gallery&sort=expense_date", 200)
	if pg.Items[0].ImageAttachmentID != nil || pg.Items[2].ImageAttachmentID == nil {
		t.Fatalf("galería: %+v", pg.Items)
	}
	// Paginación con cursor.
	pg = view(t, c, "expense?mode=table&limit=2", 200)
	if len(pg.Items) != 2 || !pg.NextCursor.Valid {
		t.Fatalf("página 1: %+v", pg)
	}
	pg = view(t, c, "expense?mode=table&limit=2&cursor="+pg.NextCursor.Value, 200)
	if len(pg.Items) != 1 || pg.NextCursor.Valid {
		t.Fatalf("página 2: %+v", pg)
	}

	// Créditos: deuda viva (motor_cred), pagos por estado y compras por estado.
	cr, _ := credFx(t, s)
	pg = view(t, c, "credit?mode=board&group_by=kind", 200)
	if len(pg.Items) != 1 || pg.Items[0].Record["live_debt_cents"] == nil || groupOf(pg, "revolving_line").Count != 1 || len(pg.Groups) != 4 {
		t.Fatalf("créditos: %+v", pg)
	}
	pg = view(t, c, "scheduled_payment?mode=board&group_by=status&credit_id="+itoa(cr), 200)
	if groupOf(pg, "overdue").Count != 1 || groupOf(pg, "pending").Count != 1 || *groupOf(pg, "overdue").TotalCents != 50000 || *groupOf(pg, "pending").ShareBp != 5000 {
		t.Fatalf("pagos por estado: %+v", pg.Groups)
	}
	pg = view(t, c, "scheduled_payment?mode=table&status=overdue", 200)
	if len(pg.Items) != 1 || pg.Items[0].Record["due_date"] != "2026-10-10" || !strings.HasPrefix(pg.Items[0].Title, "Despensa · pago 1") {
		t.Fatalf("filtro status: %+v", pg.Items)
	}
	pg = view(t, c, "credit_purchase?mode=calendar&period=2026-09", 200)
	if len(pg.Groups) != 1 || pg.Groups[0].Key != "2026-09-28" || *pg.Groups[0].TotalCents != 100000 {
		t.Fatalf("compras: %+v", pg.Groups)
	}

	// Metas: estado y ahorrado del motor S2.
	goal := ins(t, s.DB, `INSERT INTO savings_goal (name, kind, target_cents, target_date, start_date) VALUES ('Viaje', 'goal', 100000, '2027-06-30', '2026-10-01')`)
	ins(t, s.DB, `INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (?, '2026-10-02', 'deposit', 30000)`, goal)
	pg = view(t, c, "savings_goal?mode=board&group_by=status", 200)
	if len(pg.Items) != 1 || pg.Items[0].Record["saved_cents"] == nil || *pg.Items[0].AmountCents != 30000 || len(pg.Groups) != 3 {
		t.Fatalf("metas: %+v", pg)
	}
	pg = view(t, c, "savings_movement?mode=calendar", 200)
	if len(pg.Groups) != 1 || *pg.Groups[0].TotalCents != 30000 {
		t.Fatalf("movimientos: %+v", pg.Groups)
	}

	// Mercado: aparece en cada día de visita del mes en curso (ISO 6 = sábado: 3, 10, 17, 24 y 31 de octubre).
	ins(t, s.DB, `INSERT INTO store (kind, name, visit_weekday) VALUES ('market', 'Tianguis del sábado', 6)`)
	pg = view(t, c, "store?mode=calendar", 200)
	if len(pg.Items) != 5 || pg.Groups[0].Key != "2026-10-03" || pg.Groups[4].Key != "2026-10-31" {
		t.Fatalf("mercado: %+v", pg.Groups)
	}
	pg = view(t, c, "store?mode=board&group_by=kind", 200)
	if groupOf(pg, "market").Count != 1 || groupOf(pg, "chain").Count != 1 || groupOf(pg, "chain").TotalCents != nil {
		t.Fatalf("tiendas: %+v", pg.Groups)
	}

	// Vistas sin datos y sin importe; credenciales nunca salen.
	for _, p := range []string{"product?mode=board&group_by=freshness", "price?mode=board&group_by=is_stale", "inventory_item?mode=board",
		"recipe?mode=board&group_by=meal_slot_id", "meal_plan_entry?mode=calendar", "notification_delivery?mode=board&group_by=status",
		"job?mode=calendar", "backup?mode=table&status=ok", "weight_entry?mode=table", "credit_term_rate?mode=table", "promotion?mode=table",
		"category?mode=board", "income?mode=calendar", "device_session?mode=table", "person?mode=table&person_id=1"} {
		view(t, c, p, 200)
	}
	pg = view(t, c, "user_account?mode=table", 200)
	if len(pg.Items) != 2 || pg.Items[0].Record["password_hash"] != nil || pg.Items[0].Record["is_active"] != true {
		t.Fatalf("usuarios: %+v", pg.Items)
	}

	// Errores: formato = 400; regla de C1 §8 = 422; entidad inexistente = 404; sin sesión = 401.
	for p, code := range map[string]int{
		"expense":                            400,
		"expense?mode=list":                  400,
		"expense?mode=table&group_by=Bad":    400,
		"expense?mode=table&family=friends":  400,
		"expense?mode=table&limit=201":       400,
		"expense?mode=table&period=2026-13":  400,
		"income?mode=board":                  422,
		"expense?mode=board&group_by=notes":  422,
		"expense?mode=calendar&date_field=x": 422,
		"expense?mode=table&status=paid":     422,
		"backup?mode=table&status=paid":      422,
		"budget?mode=table":                  404,
	} {
		view(t, c, p, code)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/views/expense?mode=table", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}
