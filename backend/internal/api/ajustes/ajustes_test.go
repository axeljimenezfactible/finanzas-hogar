package ajustes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/ajustes"
	"finanzas-hogar/internal/testutil"
)

type env struct {
	s                        *testutil.Server
	c                        *testutil.Client
	hogar, karen, conj       int64
	groceries, rent, subGroc int64
}

var keyN int

func idem() []string {
	keyN++
	return []string{"Idempotency-Key", fmt.Sprintf("prueba-ajustes-%08d", keyN)}
}

func exec(t *testing.T, s *testutil.Server, q string, args ...any) int64 {
	t.Helper()
	res, err := s.DB.Exec(q, args...)
	if err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	id, _ := res.LastInsertId()
	return id
}

// setup: Axel, Ximena y Karen (testutil) más las reglas Hogar, Karen y Conjunto 50/50 y dos categorías.
func setup(t *testing.T) *env {
	t.Helper()
	s := testutil.NewServer(t, ajustes.Register)
	e := &env{s: s, c: s.SessionFor(t, testutil.AxelUser)}
	e.hogar = exec(t, s, `INSERT INTO cost_share_rule (name, remainder_person_id, is_default) VALUES ('Hogar', ?, 1)`, s.Fx.AxelPersonID)
	e.karen = exec(t, s, `INSERT INTO cost_share_rule (name, remainder_person_id) VALUES ('Karen', ?)`, s.Fx.KarenPersonID)
	e.conj = exec(t, s, `INSERT INTO cost_share_rule (name, remainder_person_id) VALUES ('Conjunto 50/50', ?)`, s.Fx.AxelPersonID)
	exec(t, s, `INSERT INTO cost_share_rule_part (rule_id, person_id, pct_bp) VALUES (?, ?, 5000)`, e.conj, s.Fx.KarenPersonID)
	e.groceries = exec(t, s, `INSERT INTO category (code, name, kind, sort_order, default_cost_share_rule_id) VALUES ('groceries', 'Despensa', 'variable', 1, ?)`, e.conj)
	e.rent = exec(t, s, `INSERT INTO category (name, kind, sort_order) VALUES ('Renta', 'fixed', 2)`)
	e.subGroc = exec(t, s, `INSERT INTO category (name, parent_id, sort_order) VALUES ('Frutas', ?, 1)`, e.groceries)
	return e
}

func want(t *testing.T, res *testutil.Response, code int) {
	t.Helper()
	if res.StatusCode != code {
		t.Fatalf("esperaba %d y llegó %d: %s", code, res.StatusCode, res.Body)
	}
}

func errCode(t *testing.T, res *testutil.Response) string {
	t.Helper()
	var b api.ErrorBody
	res.JSON(t, &b)
	return b.Code
}

func TestOp_listCategories(t *testing.T) {
	e := setup(t)
	var pg api.CategoryPage
	res := e.c.Do("GET", "/api/v1/categories", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 3 || pg.Items[0].Name != "Despensa" || *pg.Items[0].ExpenseCount != 0 || pg.NextCursor.Valid {
		t.Fatalf("%+v", pg)
	}
	res = e.c.Do("GET", "/api/v1/categories?sort=-recent_usage_count,sort_order&limit=1", nil)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || !pg.NextCursor.Valid {
		t.Fatalf("paginación: %+v", pg)
	}
	want(t, e.c.Do("GET", "/api/v1/categories?cursor="+pg.NextCursor.Value+"&sort=-recent_usage_count,sort_order&limit=1", nil), 200)
	pg = api.CategoryPage{}
	e.c.Do("GET", "/api/v1/categories?q=DESPÉNSA", nil).JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].ID != e.groceries {
		t.Fatalf("q: %+v", pg)
	}
	e.c.Do("GET", "/api/v1/categories?status=archived", nil).JSON(t, &pg)
	if len(pg.Items) != 0 {
		t.Fatalf("archivadas: %+v", pg)
	}
	want(t, e.c.Do("GET", "/api/v1/categories?status=otra", nil), 400)
	want(t, e.c.Do("GET", "/api/v1/categories?limit=201", nil), 400)
	want(t, e.c.Do("GET", "/api/v1/categories?sort=color", nil), 400)
	want(t, e.s.Anonymous(t).Do("GET", "/api/v1/categories", nil), 401)
}

func TestOp_createCategory(t *testing.T) {
	e := setup(t)
	var c api.Category
	res := e.c.Do("POST", "/api/v1/categories", map[string]any{"name": " Verduras ", "parent_id": e.groceries, "kind": "variable"}, idem()...)
	want(t, res, 201)
	res.JSON(t, &c)
	if c.Name != "Verduras" || c.SortOrder != 2 || *c.ParentID != e.groceries || !c.IsActive {
		t.Fatalf("%+v", c)
	}
	var n int
	_ = e.s.DB.QueryRow(`SELECT count(*) FROM audit_log WHERE table_name = 'category' AND row_id = ? AND action = 'insert'`, c.ID).Scan(&n)
	if n != 1 {
		t.Fatalf("audit_log: %d", n)
	}
	res = e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "Otra", "code": "groceries"}, idem()...)
	if want(t, res, 409); errCode(t, res) != "already_exists" {
		t.Fatal(string(res.Body))
	}
	want(t, e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "  "}, idem()...), 422)
	want(t, e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "X", "parent_id": 999}, idem()...), 422)
	want(t, e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "X", "code": "Con Espacio"}, idem()...), 422)
	want(t, e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "X", "kind": "otro"}, idem()...), 400)
	want(t, e.c.Do("POST", "/api/v1/categories", map[string]any{"name": "X", "color": "rojo"}, idem()...), 400)
}

func TestOp_updateCategory(t *testing.T) {
	e := setup(t)
	var c api.Category
	res := e.c.Do("PATCH", fmt.Sprintf("/api/v1/categories/%d", e.rent), map[string]any{"is_active": false, "kind": nil, "default_cost_share_rule_id": e.hogar})
	want(t, res, 200)
	res.JSON(t, &c)
	if c.IsActive || c.Kind != nil || *c.DefaultCostShareRuleID != e.hogar {
		t.Fatalf("%+v", c)
	}
	want(t, e.c.Do("PATCH", "/api/v1/categories/999", map[string]any{"name": "X"}), 404)
	want(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/categories/%d", e.rent), map[string]any{}), 400)
	want(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/categories/%d", e.groceries), map[string]any{"parent_id": e.subGroc}), 422)
	want(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/categories/%d", e.rent), map[string]any{"code": "groceries"}), 409)
	want(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/categories/%d", e.rent), map[string]any{"kind": "raro"}), 400)
}

func TestOp_reorderCategories(t *testing.T) {
	e := setup(t)
	want(t, e.c.Do("PUT", "/api/v1/categories/order", map[string]any{"ids": []int64{e.rent, e.subGroc, e.groceries}}), 204)
	var pg api.CategoryPage
	e.c.Do("GET", "/api/v1/categories", nil).JSON(t, &pg)
	if pg.Items[0].ID != e.rent || pg.Items[2].ID != e.groceries {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, e.c.Do("PUT", "/api/v1/categories/order", map[string]any{"ids": []int64{e.rent, e.groceries}}), 422)
	want(t, e.c.Do("PUT", "/api/v1/categories/order", map[string]any{"ids": []int64{e.rent, e.groceries, e.subGroc, 999}}), 422)
	want(t, e.c.Do("PUT", "/api/v1/categories/order", map[string]any{"ids": []int64{e.rent, e.rent, e.groceries, e.subGroc}}), 422)
	want(t, e.c.Do("PUT", "/api/v1/categories/order", map[string]any{}), 422)
	want(t, e.c.Do("PUT", "/api/v1/categories/order", "[1,2]"), 400)
}

func TestOp_listPersons(t *testing.T) {
	e := setup(t)
	var pg api.PersonPage
	res := e.s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/persons", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 3 || pg.Items[0].Participation != "holder" || pg.Items[2].Participation != "external" ||
		*pg.Items[2].HasAccount || !*pg.Items[0].HasAccount {
		t.Fatalf("%+v", pg)
	}
	want(t, e.c.Do("GET", "/api/v1/persons?sort=-display_name", nil), 200)
	want(t, e.c.Do("GET", "/api/v1/persons?sort=password_hash", nil), 400)
}

func TestOp_updatePerson(t *testing.T) {
	e := setup(t)
	path := fmt.Sprintf("/api/v1/persons/%d", e.s.Fx.KarenPersonID)
	var p api.Person
	res := e.c.Do("PATCH", path, map[string]any{"fortnight_payment_cap_cents": 150000, "display_name": "Karen M.", "notes": "Paga directo"})
	want(t, res, 200)
	res.JSON(t, &p)
	if *p.FortnightPaymentCapCents != 150000 || p.DisplayName != "Karen M." {
		t.Fatalf("%+v", p)
	}
	var p2 api.Person
	e.c.Do("PATCH", path, map[string]any{"fortnight_payment_cap_cents": nil, "notes": nil}).JSON(t, &p2)
	if p2.FortnightPaymentCapCents != nil || p2.Notes != nil {
		t.Fatalf("null no borró: %+v", p2)
	}
	want(t, e.c.Do("PATCH", path, map[string]any{"display_name": ""}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"fortnight_payment_cap_cents": -1}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"fortnight_payment_cap_cents": 10.5}), 400)
	want(t, e.c.Do("PATCH", path, map[string]any{}), 400)
	want(t, e.c.Do("PATCH", "/api/v1/persons/999", map[string]any{"display_name": "X"}), 404)
}

func TestOp_listCostShareRules(t *testing.T) {
	e := setup(t)
	var pg api.CostShareRulePage
	res := e.c.Do("GET", "/api/v1/cost-share-rules", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 3 || !pg.Items[0].IsDefault || len(pg.Items[2].Parts) != 1 || *pg.Items[2].Parts[0].PctBp != 5000 {
		t.Fatalf("%+v", pg)
	}
	e.c.Do("GET", "/api/v1/cost-share-rules?is_template=false", nil).JSON(t, &pg)
	if len(pg.Items) != 0 {
		t.Fatalf("%+v", pg)
	}
	want(t, e.c.Do("GET", "/api/v1/cost-share-rules?is_template=si", nil), 400)
}

func TestOp_updateCostShareRule(t *testing.T) {
	e := setup(t)
	path := fmt.Sprintf("/api/v1/cost-share-rules/%d", e.conj)
	var r api.CostShareRule
	res := e.c.Do("PATCH", path, map[string]any{"name": "Conjunto 60/40", "scope": "onward",
		"parts": []map[string]any{{"person_id": e.s.Fx.KarenPersonID, "pct_bp": 4000}}})
	want(t, res, 200)
	res.JSON(t, &r)
	if r.Name != "Conjunto 60/40" || len(r.Parts) != 1 || *r.Parts[0].PctBp != 4000 {
		t.Fatalf("%+v", r)
	}
	e.c.Do("PATCH", path, map[string]any{"notes": "nota"}).JSON(t, &r)
	if *r.Notes != "nota" || len(r.Parts) != 1 {
		t.Fatalf("%+v", r)
	}
	kp := []map[string]any{{"person_id": e.s.Fx.KarenPersonID, "pct_bp": 6000}}
	want(t, e.c.Do("PATCH", path, map[string]any{"parts": kp}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"scope": "onward"}), 422)
	if res := e.c.Do("PATCH", path, map[string]any{"parts": kp, "scope": "all"}); res.StatusCode != 409 || errCode(t, res) != "conflict" {
		t.Fatalf("all: %d %s", res.StatusCode, res.Body)
	}
	want(t, e.c.Do("PATCH", path, map[string]any{"scope": "onward", "parts": []map[string]any{
		{"person_id": e.s.Fx.KarenPersonID, "pct_bp": 6000}, {"person_id": e.s.Fx.XimenaPersonID, "pct_bp": 5000}}}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"scope": "onward", "parts": []map[string]any{
		{"person_id": e.s.Fx.KarenPersonID, "pct_bp": 100, "fixed_cents": 100}, {"person_id": 999, "fixed_cents": 5}}}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"scope": "siempre", "parts": kp}), 400)
	want(t, e.c.Do("PATCH", path, map[string]any{}), 400)
	want(t, e.c.Do("PATCH", "/api/v1/cost-share-rules/999", map[string]any{"name": "X"}), 404)
}

func plannedBody(e *env, kind string) map[string]any {
	if kind == "income" {
		return map[string]any{"kind": "income", "description": "Sueldo", "amount_cents": 1500000, "second_half_amount_cents": 1450000,
			"frequency": "fortnightly", "start_date": "2026-10-15", "person_id": e.s.Fx.AxelPersonID}
	}
	return map[string]any{"kind": "fixed_expense", "description": "Renta", "amount_cents": 950000, "frequency": "monthly",
		"start_date": "2026-11-01", "category_id": e.rent, "cost_share_rule_id": e.hogar, "payment_method": "transfer"}
}

func TestOp_createPlannedItem(t *testing.T) {
	e := setup(t)
	var p api.PlannedItem
	res := e.c.Do("POST", "/api/v1/planned-items", plannedBody(e, "income"), idem()...)
	want(t, res, 201)
	res.JSON(t, &p)
	if p.AmountCents != 1500000 || *p.SecondHalfAmountCents != 1450000 || p.IntervalCount != 1 || p.AccrualMonths != 1 || !p.IsActive {
		t.Fatalf("%+v", p)
	}
	want(t, e.c.Do("POST", "/api/v1/planned-items", plannedBody(e, "fixed"), idem()...), 201)
	bad := plannedBody(e, "income")
	bad["person_id"] = e.s.Fx.XimenaPersonID
	want(t, e.c.Do("POST", "/api/v1/planned-items", bad, idem()...), 422)
	bad = plannedBody(e, "fixed")
	delete(bad, "category_id")
	bad["second_half_amount_cents"] = 5
	bad["payment_method"] = "credit"
	want(t, e.c.Do("POST", "/api/v1/planned-items", bad, idem()...), 422)
	bad = plannedBody(e, "fixed")
	bad["frequency"] = "once"
	bad["interval_count"] = 2
	bad["end_date"] = "2026-10-01"
	bad["cost_share_rule_id"] = 999
	want(t, e.c.Do("POST", "/api/v1/planned-items", bad, idem()...), 422)
	want(t, e.c.Do("POST", "/api/v1/planned-items", map[string]any{"kind": "income"}, idem()...), 422)
	bad = plannedBody(e, "fixed")
	bad["start_date"] = "2026-02-30"
	want(t, e.c.Do("POST", "/api/v1/planned-items", bad, idem()...), 400)
	bad = plannedBody(e, "fixed")
	bad["frequency"] = "daily"
	want(t, e.c.Do("POST", "/api/v1/planned-items", bad, idem()...), 400)
	want(t, e.c.Do("POST", "/api/v1/planned-items", plannedBody(e, "fixed")), 400) // sin Idempotency-Key
}

func TestOp_listPlannedItems(t *testing.T) {
	e := setup(t)
	want(t, e.c.Do("POST", "/api/v1/planned-items", plannedBody(e, "income"), idem()...), 201)
	want(t, e.c.Do("POST", "/api/v1/planned-items", plannedBody(e, "fixed"), idem()...), 201)
	var pg api.PlannedItemPage
	res := e.c.Do("GET", "/api/v1/planned-items", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 2 || pg.Items[0].Kind != "income" {
		t.Fatalf("%+v", pg)
	}
	e.c.Do("GET", fmt.Sprintf("/api/v1/planned-items?kind=fixed_expense&category_id=%d&sort=-amount_cents", e.rent), nil).JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].Description != "Renta" {
		t.Fatalf("%+v", pg)
	}
	e.c.Do("GET", "/api/v1/planned-items?credit_id=7", nil).JSON(t, &pg)
	if len(pg.Items) != 0 {
		t.Fatalf("%+v", pg)
	}
	want(t, e.c.Do("GET", "/api/v1/planned-items?kind=otro", nil), 400)
	want(t, e.c.Do("GET", "/api/v1/planned-items?category_id=a", nil), 400)
}

func TestOp_updatePlannedItem(t *testing.T) {
	e := setup(t)
	var p api.PlannedItem
	e.c.Do("POST", "/api/v1/planned-items", plannedBody(e, "income"), idem()...).JSON(t, &p)
	path := fmt.Sprintf("/api/v1/planned-items/%d", p.ID)
	res := e.c.Do("PATCH", path, map[string]any{"second_half_amount_cents": nil, "is_active": false, "end_date": "2027-12-31", "notes": "sube en enero"})
	want(t, res, 200)
	var u api.PlannedItem
	res.JSON(t, &u)
	if u.SecondHalfAmountCents != nil || u.IsActive || *u.EndDate != "2027-12-31" {
		t.Fatalf("%+v", u)
	}
	want(t, e.c.Do("PATCH", path, map[string]any{"frequency": "monthly", "second_half_amount_cents": 10}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"category_id": e.rent}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"end_date": "31-12-2027"}), 400)
	want(t, e.c.Do("PATCH", path, map[string]any{}), 400)
	want(t, e.c.Do("PATCH", "/api/v1/planned-items/999", map[string]any{"is_active": true}), 404)
}

func multipartCSV(t *testing.T, ctype, csv string) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="file"; filename="gastos_fijos_ingresos.csv"`)
	h.Set("Content-Type", ctype)
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte(csv))
	_ = w.Close()
	return buf.Bytes(), w.FormDataContentType()
}

const csvHeader = "tipo,concepto,monto,periodicidad,dia_o_fecha,categoria,persona,porcentaje_conjunto,notas\n"

const csvGood = "\xef\xbb\xbf" + csvHeader +
	"ingreso,Sueldo Axel,\"15,000.00\",quincenal,,,axel,,\n" +
	"gasto_fijo,Renta,$9500,mensual,1,Renta,hogar,,depósito\n" +
	"gasto_fijo,Internet,599.5,mensual,31,renta,ximena,,\n" +
	"gasto_venidero,Despensa,2383,unica,2026-11-15,groceries,conjunto,50,\n" +
	"gasto_fijo,Seguro,4800,anual,2027-01-10,Renta,karen,,\n"

func TestOp_previewPlannedItemsImport(t *testing.T) {
	e := setup(t)
	csv := csvGood + "ingreso,Bono,abc,semanal,,Renta,ximena,10,\n" +
		"gasto_fijo,Luz,300,mensual,40,Inexistente,conjunto,60,\n" +
		"gasto_fijo,,300,unica,15,Renta,nadie,,\n"
	body, ct := multipartCSV(t, "text/csv", csv)
	res := e.c.Do("POST", "/api/v1/imports/planned-items/preview", body, "Content-Type", ct)
	want(t, res, 200)
	var pv api.ImportPreview
	res.JSON(t, &pv)
	if pv.ValidCount != 5 || pv.ErrorCount != 3 || pv.Rows[0].LineNo != 2 || pv.Rows[0].PlannedItem.StartDate != "2026-10-15" ||
		pv.Rows[1].PlannedItem.StartDate != "2026-11-01" || pv.Rows[2].PlannedItem.StartDate != "2026-10-31" ||
		pv.Rows[0].PlannedItem.AmountCents != 1500000 || pv.Rows[2].PlannedItem.AmountCents != 59950 ||
		*pv.Rows[3].PlannedItem.CostShareRuleID != e.conj || *pv.Rows[4].PlannedItem.CostShareRuleID != e.karen {
		t.Fatalf("%s", res.Body)
	}
	var n int
	_ = e.s.DB.QueryRow(`SELECT count(*) FROM planned_item`).Scan(&n)
	if n != 0 {
		t.Fatal("la vista previa escribió")
	}
	body, ct = multipartCSV(t, "application/json", csvGood)
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items/preview", body, "Content-Type", ct), 415)
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items/preview", []byte(csvGood), "Content-Type", "text/csv"), 415)
	body, ct = multipartCSV(t, "text/csv", "tipo;concepto;monto\ningreso;Sueldo;100\n")
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items/preview", body, "Content-Type", ct), 422)
	body, ct = multipartCSV(t, "text/csv", csvHeader)
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items/preview", body, "Content-Type", ct), 422)
	body, ct = multipartCSV(t, "text/csv", "\"sin cerrar\n")
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items/preview", body, "Content-Type", ct), 422)
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("otro", "x")
	_ = w.Close()
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items/preview", buf.Bytes(), "Content-Type", w.FormDataContentType()), 400)
}

func TestOp_importPlannedItems(t *testing.T) {
	e := setup(t)
	body, ct := multipartCSV(t, "text/csv", csvGood+"gasto_fijo,Luz,300,mensual,40,Renta,hogar,,\n")
	res := e.c.Do("POST", "/api/v1/imports/planned-items", body, append(idem(), "Content-Type", ct)...)
	if want(t, res, 422); !strings.Contains(string(res.Body), "rows[7].dia_o_fecha") {
		t.Fatalf("%s", res.Body)
	}
	body, ct = multipartCSV(t, "text/csv", strings.ReplaceAll(strings.ReplaceAll(csvGood, ",", ";"), `"15;000.00"`, "15000"))
	h := append(idem(), "Content-Type", ct)
	res = e.c.Do("POST", "/api/v1/imports/planned-items", body, h...)
	want(t, res, 201)
	var ir api.ImportResult
	res.JSON(t, &ir)
	if ir.CreatedCount != 5 || len(ir.PlannedItemIDs) != 5 {
		t.Fatalf("%+v", ir)
	}
	res = e.c.Do("POST", "/api/v1/imports/planned-items", body, h...)
	if res.Header.Get("Idempotent-Replayed") != "true" {
		t.Fatalf("repetición: %d %v", res.StatusCode, res.Header)
	}
	want(t, e.c.Do("POST", "/api/v1/imports/planned-items", []byte("{}"), append(idem(), "Content-Type", "application/json")...), 415)
}

func TestOp_getProjectionSettings(t *testing.T) {
	e := setup(t)
	var ps api.ProjectionSettings
	res := e.c.Do("GET", "/api/v1/settings/projection", nil)
	want(t, res, 200)
	if res.JSON(t, &ps); ps.RepeatCadenceEnabled {
		t.Fatal("K21: apagado por defecto")
	}
	want(t, e.s.Anonymous(t).Do("GET", "/api/v1/settings/projection", nil), 401)
}

func TestOp_updateProjectionSettings(t *testing.T) {
	e := setup(t)
	var ps api.ProjectionSettings
	res := e.c.Do("PATCH", "/api/v1/settings/projection", map[string]any{"repeat_cadence_enabled": true})
	want(t, res, 200)
	if res.JSON(t, &ps); !ps.RepeatCadenceEnabled {
		t.Fatal("no cambió")
	}
	want(t, e.c.Do("PATCH", "/api/v1/settings/projection", map[string]any{}), 422)
	want(t, e.c.Do("PATCH", "/api/v1/settings/projection", map[string]any{"repeat_cadence_enabled": nil}), 422)
	want(t, e.c.Do("PATCH", "/api/v1/settings/projection", map[string]any{"repeat_cadence_enabled": "si"}), 400)
	want(t, e.c.Do("PATCH", "/api/v1/settings/projection", map[string]any{"otro": true}), 400)
}

func goalBody() map[string]any {
	return map[string]any{"name": "Fondo de emergencia", "kind": "goal", "target_cents": 3000000, "target_date": "2027-06-30",
		"planned_monthly_cents": 300000, "start_date": "2026-10-01"}
}

func TestOp_createSavingsGoal(t *testing.T) {
	e := setup(t)
	var g api.SavingsGoal
	res := e.c.Do("POST", "/api/v1/savings-goals", goalBody(), idem()...)
	want(t, res, 201)
	res.JSON(t, &g)
	if *g.Status != "in_progress" || *g.SavedCents != 0 || *g.ProgressPct != 0 || *g.RequiredMonthlyCents != 333333 {
		t.Fatalf("%s", res.Body)
	}
	if res := e.c.Do("POST", "/api/v1/savings-goals", goalBody(), idem()...); res.StatusCode != 409 || errCode(t, res) != "already_exists" {
		t.Fatalf("%d %s", res.StatusCode, res.Body)
	}
	want(t, e.c.Do("POST", "/api/v1/savings-goals", map[string]any{"name": "Viaje", "kind": "goal", "start_date": "2026-10-01"}, idem()...), 422)
	want(t, e.c.Do("POST", "/api/v1/savings-goals", map[string]any{"name": "Viaje", "kind": "fund"}, idem()...), 422)
	want(t, e.c.Do("POST", "/api/v1/savings-goals", map[string]any{"name": "Viaje", "kind": "meta", "start_date": "2026-10-01"}, idem()...), 400)
	want(t, e.c.Do("POST", "/api/v1/savings-goals", map[string]any{"name": "Viaje", "kind": "fund", "start_date": "2026-13-01"}, idem()...), 400)
	var f api.SavingsGoal
	e.c.Do("POST", "/api/v1/savings-goals", map[string]any{"name": "Colchón", "kind": "fund", "start_date": "2026-10-01"}, idem()...).JSON(t, &f)
	if f.ProgressPct != nil || *f.Status != "in_progress" {
		t.Fatalf("fondo: %+v", f)
	}
}

func TestOp_updateSavingsGoal(t *testing.T) {
	e := setup(t)
	var g api.SavingsGoal
	e.c.Do("POST", "/api/v1/savings-goals", goalBody(), idem()...).JSON(t, &g)
	exec(t, e.s, `INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents) VALUES (?, '2026-10-05', 'deposit', 1500000)`, g.ID)
	path := fmt.Sprintf("/api/v1/savings-goals/%d", g.ID)
	res := e.c.Do("PATCH", path, map[string]any{"notes": "medio camino", "planned_monthly_cents": nil})
	want(t, res, 200)
	g = api.SavingsGoal{}
	res.JSON(t, &g)
	if *g.SavedCents != 1500000 || *g.ProgressTenthsPct != 500 || *g.ProgressPct != 50 || g.PlannedMonthlyCents != nil {
		t.Fatalf("%s", res.Body)
	}
	g = api.SavingsGoal{}
	e.c.Do("PATCH", path, map[string]any{"is_active": false}).JSON(t, &g)
	if g.IsActive || g.Status != nil {
		t.Fatalf("inactiva: %+v", g)
	}
	want(t, e.c.Do("PATCH", path, map[string]any{"target_date": "2026-09-01"}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{"target_cents": nil}), 422)
	want(t, e.c.Do("PATCH", path, map[string]any{}), 400)
	want(t, e.c.Do("PATCH", "/api/v1/savings-goals/999", map[string]any{"name": "X"}), 404)
	e.c.Do("POST", "/api/v1/savings-goals", map[string]any{"name": "Otra", "kind": "fund", "start_date": "2026-10-01"}, idem()...)
	want(t, e.c.Do("PATCH", path, map[string]any{"name": "otra"}), 409)
}

func TestOp_getAlertSettings(t *testing.T) {
	e := setup(t)
	var a api.AlertSettings
	res := e.c.Do("GET", "/api/v1/settings/alerts", nil)
	want(t, res, 200)
	res.JSON(t, &a)
	p := a.Parameters
	if p.AlertsBudgetWarningPct.Value.Value != 80 || p.AlertsBudgetLimitPct.Value.Value != 100 || !p.AlertsBudgetLimitPct.IsInitial ||
		len(p.AlertsPaymentDueDaysBefore.Value.Value) != 2 || p.AlertsSendTime.Value.Value != "09:00" ||
		p.AlertsStalePricesAdvanceDays.Value.Value != 1 || p.AlertsMealPlanMissingWeekday.InitialValue.Value != 7 {
		t.Fatalf("%s", res.Body)
	}
	want(t, e.s.Anonymous(t).Do("GET", "/api/v1/settings/alerts", nil), 401)
}

func TestOp_updateAlertSettings(t *testing.T) {
	e := setup(t)
	var a api.AlertSettings
	res := e.c.Do("PATCH", "/api/v1/settings/alerts", map[string]any{
		"parameters":       map[string]any{"alerts.budget_warning_pct": 70, "alerts.send_time": "08:30", "alerts.payment_due_days_before": []int{2, 0}},
		"reset_parameters": []string{"alerts.budget_limit_pct"}})
	want(t, res, 200)
	res.JSON(t, &a)
	if a.Parameters.AlertsBudgetWarningPct.Value.Value != 70 || a.Parameters.AlertsBudgetWarningPct.IsInitial ||
		a.Parameters.AlertsSendTime.Value.Value != "08:30" || a.Parameters.AlertsPaymentDueDaysBefore.Value.Value[0] != 2 {
		t.Fatalf("%s", res.Body)
	}
	a = api.AlertSettings{}
	e.c.Do("PATCH", "/api/v1/settings/alerts", map[string]any{"reset_parameters": []string{"alerts.budget_warning_pct"}}).JSON(t, &a)
	if !a.Parameters.AlertsBudgetWarningPct.IsInitial {
		t.Fatal("no restableció")
	}
	for _, body := range []map[string]any{
		{"parameters": map[string]any{"alerts.meal_plan_missing_weekday": 8}},
		{"parameters": map[string]any{"alerts.budget_warning_pct": nil}},
		{"parameters": map[string]any{"alerts.send_time": "25:00"}},
		{"parameters": map[string]any{"alerts.budget_warning_pct": 50}, "reset_parameters": []string{"alerts.budget_warning_pct"}},
		{"reset_parameters": []string{"auth.max_failed_attempts"}},
	} {
		if res := e.c.Do("PATCH", "/api/v1/settings/alerts", body); res.StatusCode != 422 {
			b, _ := json.Marshal(body)
			t.Fatalf("%s → %d %s", b, res.StatusCode, res.Body)
		}
	}
	for _, body := range []map[string]any{{}, {"parameters": map[string]any{}}, {"reset_parameters": []string{}},
		{"parameters": map[string]any{"alerts.otra": 1}}, {"reset_parameters": []string{"alerts.send_time", "alerts.send_time"}}} {
		if res := e.c.Do("PATCH", "/api/v1/settings/alerts", body); res.StatusCode != 400 {
			b, _ := json.Marshal(body)
			t.Fatalf("%s → %d %s", b, res.StatusCode, res.Body)
		}
	}
}
