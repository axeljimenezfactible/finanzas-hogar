package creditos_test

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/creditos"
	"finanzas-hogar/internal/api/presupuesto"
	"finanzas-hogar/internal/api/tablero"
	"finanzas-hogar/internal/testutil"
)

var keyN atomic.Int64

func key() string { return fmt.Sprintf("creditos-test-key-%08d", keyN.Add(1)) }

func must(t *testing.T, res *testutil.Response, code int, v any) {
	t.Helper()
	if res.StatusCode != code {
		t.Fatalf("esperaba %d, llegó %d: %s", code, res.StatusCode, res.Body)
	}
	if v != nil {
		res.JSON(t, v)
	}
}

func errCode(t *testing.T, res *testutil.Response, status int, code string) {
	t.Helper()
	if res.StatusCode != status || !strings.Contains(string(res.Body), `"code":"`+code+`"`) {
		t.Fatalf("esperaba %d %s, llegó %d: %s", status, code, res.StatusCode, res.Body)
	}
}

type env struct {
	s                       *testutil.Server
	c                       *testutil.Client
	kueski, card, rate, buy int64
	axel, karen, ximena     int64
}

// setup: Kueski (canon semilla 10-oct-2026) con su tasa a 4 quincenas y la despensa 1 repartida Axel/Karen.
func setup(t *testing.T) *env {
	s := testutil.NewServer(t, creditos.Register)
	e := &env{s: s, c: s.SessionFor(t, testutil.AxelUser), axel: s.Fx.AxelPersonID, karen: s.Fx.KarenPersonID, ximena: s.Fx.XimenaPersonID}
	var cr api.Credit
	must(t, e.c.Do("POST", "/api/v1/credits", map[string]any{"name": "Kueski Pay", "kind": "revolving_line", "holder_person_id": e.axel,
		"credit_limit_cents": 1000000, "payment_frequency": "fortnightly", "first_payment_min_days": 15}, "Idempotency-Key", key()), 201, &cr)
	e.kueski = cr.ID
	must(t, e.c.Do("POST", "/api/v1/credits", map[string]any{"name": "Tarjeta de prueba", "kind": "bank_card", "holder_person_id": e.axel,
		"credit_limit_cents": 5000000, "payment_frequency": "monthly", "statement_day": 5, "payment_due_day": 25}, "Idempotency-Key", key()), 201, &cr)
	e.card = cr.ID
	var tr api.CreditTermRate
	must(t, e.c.Do("POST", fmt.Sprintf("/api/v1/credits/%d/term-rates", e.kueski), map[string]any{"term_count": 4, "rate_num": 80050,
		"rate_den": 309450, "source": "quote", "quote_date": "2026-10-10", "rate_decimals": 4}, "Idempotency-Key", key()), 201, &tr)
	e.rate = tr.ID
	if tr.RatePct == nil || *tr.RatePct != "25.8685" {
		t.Fatalf("rate_pct %v", tr.RatePct)
	}
	e.buy = e.purchase(t, 201).ID
	return e
}

func (e *env) purchaseBody() map[string]any {
	return map[string]any{"credit_id": e.kueski, "kind": "installment_plan", "purchase_date": "2026-10-10", "principal_cents": 309450,
		"interest_cents": 80050, "term_count": 4, "shares": []map[string]any{{"person_id": e.axel, "principal_cents": 209450},
			{"person_id": e.karen, "principal_cents": 100000}}}
}

func (e *env) purchase(t *testing.T, code int) api.CreditPurchase {
	var p api.CreditPurchase
	must(t, e.c.Do("POST", "/api/v1/credit-purchases", e.purchaseBody(), "Idempotency-Key", key()), code, &p)
	return p
}

func (e *env) scheduled(t *testing.T, q string) []api.ScheduledPayment {
	var pg api.ScheduledPaymentPage
	must(t, e.c.Do("GET", "/api/v1/scheduled-payments"+q, nil), 200, &pg)
	return pg.Items
}

func (e *env) pay(t *testing.T, spID, person, cents int64, code int) api.CreditPayment {
	var m api.CreditPayment
	res := e.c.Do("POST", fmt.Sprintf("/api/v1/scheduled-payments/%d/payments", spID), map[string]any{"payer_person_id": person,
		"paid_date": "2026-10-12", "amount_cents": cents}, "Idempotency-Key", key())
	if code >= 300 {
		must(t, res, code, nil)
		return m
	}
	must(t, res, code, &m)
	return m
}

func TestOp_createCredit(t *testing.T) {
	e := setup(t)
	body := map[string]any{"name": "Liverpool", "kind": "store_credit", "holder_person_id": e.axel, "payment_frequency": "fortnightly",
		"vat_on_interest_bp": 1600}
	errCode(t, e.c.Do("POST", "/api/v1/credits", body, "Idempotency-Key", key()), 422, "validation_failed") // Q2-67: vat_included ausente
	body["vat_included"] = true
	errCode(t, e.c.Do("POST", "/api/v1/credits", body, "Idempotency-Key", key()), 422, "validation_failed")
	body["vat_included"] = false
	var cr api.Credit
	must(t, e.c.Do("POST", "/api/v1/credits", body, "Idempotency-Key", key()), 201, &cr)
	if len(cr.MissingTerms) != 4 || cr.MissingTerms[0] != "annual_rate" {
		t.Fatalf("missing_terms %v", cr.MissingTerms)
	}
	errCode(t, e.c.Do("POST", "/api/v1/credits", body, "Idempotency-Key", key()), 409, "already_exists")
	errCode(t, e.c.Raw("POST", "/api/v1/credits", body), 400, "bad_request") // sin Idempotency-Key
}

func TestOp_updateCredit(t *testing.T) {
	e := setup(t)
	var cr api.Credit
	must(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credits/%d", e.kueski), map[string]any{"availability_reading": "A",
		"reading_a_anchor_date": "2026-10-10", "lender": nil}), 200, &cr)
	if cr.AvailabilityReading != "A" || cr.AvailableCents == nil || *cr.AvailableCents != 1000000 {
		t.Fatalf("lectura A: %+v", cr)
	}
	errCode(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credits/%d", e.card), map[string]any{"vat_on_interest_bp": 1600}), 422, "validation_failed")
	errCode(t, e.c.Do("PATCH", "/api/v1/credits/999", map[string]any{"name": "x"}), 404, "not_found")
}

func TestOp_listCredits(t *testing.T) {
	e := setup(t)
	var pg api.CreditPage
	must(t, e.c.Do("GET", "/api/v1/credits?sort=name&limit=200", nil), 200, &pg)
	if len(pg.Items) != 2 || pg.Items[0].Name != "Kueski Pay" {
		t.Fatalf("%+v", pg.Items)
	}
	must(t, e.c.Do("GET", "/api/v1/credits?kind=bank_card", nil), 200, &pg)
	if len(pg.Items) != 1 || pg.Items[0].ID != e.card {
		t.Fatalf("filtro kind: %+v", pg.Items)
	}
	errCode(t, e.c.Do("GET", "/api/v1/credits?kind=otro", nil), 400, "bad_request")
	errCode(t, e.c.Do("GET", "/api/v1/credits?limit=201", nil), 400, "bad_request")
	errCode(t, e.s.Anonymous(t).Do("GET", "/api/v1/credits", nil), 401, "unauthorized")
}

func TestOp_getCredit(t *testing.T) {
	e := setup(t)
	var cr api.Credit
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d", e.kueski), nil), 200, &cr)
	// Canon despensa 1: 4 pagos de 97375 desde 2026-10-30; deuda viva 389500; B = 1000000 − 309450.
	if *cr.LiveDebtCents != 389500 || *cr.NextPaymentDate != "2026-10-30" || *cr.NextPaymentCents != 97375 ||
		*cr.AvailableCents != 690550 || *cr.LivePurchaseCount != 1 {
		t.Fatalf("%+v", cr)
	}
	errCode(t, e.c.Do("GET", "/api/v1/credits/999", nil), 404, "not_found")
}

func TestOp_getCreditTotals(t *testing.T) {
	e := setup(t)
	var tt api.CreditTotals
	must(t, e.c.Do("GET", "/api/v1/credits/totals", nil), 200, &tt)
	if tt.LiveDebtCents != 389500 || tt.MonthTotalCents != 97375 || tt.MonthHouseholdCents+tt.MonthExternalCents != 97375 ||
		tt.MonthExternalCents == 0 || tt.FortnightTotalCents != 0 { // hoy 12-oct: quincena 1–15; el pago es el 30
		t.Fatalf("%+v", tt)
	}
	errCode(t, e.s.Anonymous(t).Do("GET", "/api/v1/credits/totals", nil), 401, "unauthorized")
}

func TestOp_getDebtTrend(t *testing.T) {
	e := setup(t)
	var dt api.DebtTrend
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/debt-trend?credit_id=%d", e.kueski), nil), 200, &dt)
	if len(dt.Months) != 12 || dt.Months[0].Month != "2026-10" || dt.Months[0].LiveDebtCents != 292125 || dt.Months[2].LiveDebtCents != 0 {
		t.Fatalf("%+v", dt.Months[:3])
	}
	errCode(t, e.c.Do("GET", "/api/v1/credits/debt-trend?credit_id=abc", nil), 400, "bad_request")
}

func TestOp_previewAvailability(t *testing.T) {
	e := setup(t)
	var pv api.AvailabilityPreview
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d/availability-preview", e.kueski), nil), 200, &pv)
	if pv.AvailabilityReading != "B" || pv.ReadingB.AvailableTodayCents != 690550 || pv.ReadingA.AvailableTodayCents != 1000000 || pv.AvailabilityConfirmed {
		t.Fatalf("%+v", pv)
	}
	errCode(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d/availability-preview", e.card), nil), 422, "not_revolving")
	errCode(t, e.c.Do("GET", "/api/v1/credits/999/availability-preview", nil), 404, "not_found")
}

func TestOp_getAvailabilitySeries(t *testing.T) {
	e := setup(t)
	var sr api.AvailabilitySeries
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d/availability-series", e.kueski), nil), 200, &sr)
	if len(sr.Events) != 5 || sr.Events[0].EventType != "start" || sr.Events[4].AvailableBCents != 1000000 {
		t.Fatalf("%+v", sr)
	}
	errCode(t, e.c.Do("GET", "/api/v1/credits/999/availability-series", nil), 404, "not_found")
}

func TestOp_getCardStatus(t *testing.T) {
	e := setup(t)
	var cs api.CardStatus
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d/card-status", e.card), nil), 200, &cs)
	if cs.CreditID != e.card || cs.CurrentBalanceCents != 0 || cs.StatementDate != nil {
		t.Fatalf("%+v", cs)
	}
	errCode(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d/card-status", e.kueski), nil), 422, "not_a_card")
}

func TestOp_listCreditTermRates(t *testing.T) {
	e := setup(t)
	var pg api.CreditTermRatePage
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d/term-rates", e.kueski), nil), 200, &pg)
	if len(pg.Items) != 1 || pg.Items[0].TermCount != 4 {
		t.Fatalf("%+v", pg)
	}
	errCode(t, e.c.Do("GET", "/api/v1/credits/999/term-rates", nil), 404, "not_found")
}

func TestOp_createCreditTermRate(t *testing.T) {
	e := setup(t)
	path := fmt.Sprintf("/api/v1/credits/%d/term-rates", e.kueski)
	var tr api.CreditTermRate
	must(t, e.c.Do("POST", path, map[string]any{"term_count": 2, "rate_num": 1727, "rate_den": 10000, "source": "manual", "rate_decimals": 2},
		"Idempotency-Key", key()), 201, &tr)
	if *tr.RatePct != "17.27" {
		t.Fatalf("%+v", tr)
	}
	errCode(t, e.c.Do("POST", path, map[string]any{"term_count": 2, "rate_num": 1, "rate_den": 10, "source": "manual"}, "Idempotency-Key", key()), 409, "already_exists")
	errCode(t, e.c.Do("POST", path, map[string]any{"term_count": 6, "rate_num": 1, "rate_den": 10, "source": "quote"}, "Idempotency-Key", key()), 422, "validation_failed")
}

func TestOp_updateCreditTermRate(t *testing.T) {
	e := setup(t)
	var tr api.CreditTermRate
	must(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credit-term-rates/%d", e.rate), map[string]any{"notes": "cotización", "rate_decimals": nil}), 200, &tr)
	if tr.RatePct != nil || *tr.Notes != "cotización" {
		t.Fatalf("%+v", tr)
	}
	errCode(t, e.c.Do("PATCH", "/api/v1/credit-term-rates/999", map[string]any{"notes": "x"}), 404, "not_found")
}

func TestOp_deleteCreditTermRate(t *testing.T) {
	e := setup(t)
	must(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-term-rates/%d", e.rate), nil), 204, nil)
	errCode(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-term-rates/%d", e.rate), nil), 404, "not_found")
}

func (e *env) promo(t *testing.T, from, to string, code int) api.Promotion {
	var p api.Promotion
	res := e.c.Do("POST", "/api/v1/promotions", map[string]any{"credit_id": e.kueski, "merchant_name": "Chedraui", "term_count": 4,
		"min_purchase_cents": 200000, "valid_from": from, "valid_to": to, "rate_decimals": 0}, "Idempotency-Key", key())
	must(t, res, code, nil)
	if code == 201 {
		res.JSON(t, &p)
	}
	return p
}

func TestOp_createPromotion(t *testing.T) {
	e := setup(t)
	p := e.promo(t, "2026-10-30", "2026-10-30", 201)
	if p.RateNum != 0 || p.RateDen != 1 || *p.RatePct != "0" || *p.IsCurrent {
		t.Fatalf("%+v", p)
	}
	e.promo(t, "2026-10-30", "2026-10-01", 422)
}

func TestOp_listPromotions(t *testing.T) {
	e := setup(t)
	e.promo(t, "2026-10-30", "2026-10-30", 201)
	e.promo(t, "2026-09-01", "2026-09-30", 201) // vencida: no sale por defecto
	var pg api.PromotionPage
	must(t, e.c.Do("GET", "/api/v1/promotions", nil), 200, &pg)
	if len(pg.Items) != 1 {
		t.Fatalf("%+v", pg.Items)
	}
	must(t, e.c.Do("GET", "/api/v1/promotions?valid_on=2026-09-15&merchant_name=chedraui", nil), 200, &pg)
	if len(pg.Items) != 1 || pg.Items[0].ValidFrom != "2026-09-01" {
		t.Fatalf("%+v", pg.Items)
	}
	errCode(t, e.c.Do("GET", "/api/v1/promotions?valid_on=10-oct-2026", nil), 400, "bad_request")
}

func TestOp_updatePromotion(t *testing.T) {
	e := setup(t)
	p := e.promo(t, "2026-10-01", "2026-10-30", 201)
	var up api.Promotion
	must(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/promotions/%d", p.ID), map[string]any{"is_active": false, "credit_id": nil}), 200, &up)
	if up.IsActive || up.CreditID != nil || *up.IsCurrent {
		t.Fatalf("%+v", up)
	}
	errCode(t, e.c.Do("PATCH", "/api/v1/promotions/999", map[string]any{"is_active": false}), 404, "not_found")
}

func TestOp_deletePromotion(t *testing.T) {
	e := setup(t)
	p := e.promo(t, "2026-10-01", "2026-10-30", 201)
	body := e.purchaseBody()
	body["principal_cents"], body["interest_cents"], body["shares"], body["promotion_id"] = 250000, nil, nil, p.ID
	delete(body, "interest_cents")
	delete(body, "shares")
	var cp api.CreditPurchase
	must(t, e.c.Do("POST", "/api/v1/credit-purchases", body, "Idempotency-Key", key()), 201, &cp)
	if cp.InterestSource != "promotion" || cp.InterestCents != 0 {
		t.Fatalf("%+v", cp)
	}
	errCode(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/promotions/%d", p.ID), nil), 409, "in_use")
	q := e.promo(t, "2026-11-01", "2026-11-30", 201)
	must(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/promotions/%d", q.ID), nil), 204, nil)
	errCode(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/promotions/%d", q.ID), nil), 404, "not_found")
}

func TestOp_getOptimizerSettings(t *testing.T) {
	e := setup(t)
	var os api.OptimizerSettings
	must(t, e.c.Do("GET", "/api/v1/settings/optimizer", nil), 200, &os)
	if len(os.MissingConstraints) != 3 || os.Weights.LowestInterest != 1 || os.Parameters.ComparatorTimeoutSeconds.Value.Value != 5 ||
		os.Parameters.OptimizerDateWindowDays.Value.Value != 7 {
		t.Fatalf("%+v", os)
	}
	errCode(t, e.s.Anonymous(t).Do("GET", "/api/v1/settings/optimizer", nil), 401, "unauthorized")
}

func TestOp_updateOptimizerSettings(t *testing.T) {
	e := setup(t)
	var os api.OptimizerSettings
	must(t, e.c.Do("PATCH", "/api/v1/settings/optimizer", map[string]any{"min_available_buffer_cents": 100000, "max_live_purchases": 3,
		"zero_debt_by_date": "2027-03-31", "parameters": map[string]any{"optimizer.date_window_days": 3}}), 200, &os)
	if len(os.MissingConstraints) != 0 || *os.MaxLivePurchases != 3 || os.Parameters.OptimizerDateWindowDays.Value.Value != 3 {
		t.Fatalf("%+v", os)
	}
	errCode(t, e.c.Do("PATCH", "/api/v1/settings/optimizer", map[string]any{"max_live_purchases": nil}), 422, "validation_failed")
	errCode(t, e.c.Do("PATCH", "/api/v1/settings/optimizer", map[string]any{"parameters": map[string]any{"comparator.timeout_seconds": 99}}), 422, "validation_failed")
}

func TestOp_compareScenarios(t *testing.T) {
	e := setup(t)
	req := map[string]any{"amount_cents": 250000, "desired_date": "2026-10-30", "merchant_name": "Chedraui"}
	var res api.ScenarioResult
	must(t, e.c.Do("POST", "/api/v1/scenario-comparisons", req), 200, &res)
	if res.Status != "missing_constraints" || len(res.MissingConstraints) != 3 || res.Recommendation != nil {
		t.Fatalf("%+v", res)
	}
	req["min_available_buffer_cents"], req["max_live_purchases"], req["zero_debt_by_date"], req["date_window_days"] = 0, 5, "2027-12-31", 1
	must(t, e.c.Do("POST", "/api/v1/scenario-comparisons", req), 200, &res)
	if res.Status != "ok" || res.EvaluatedCount == 0 || res.Recommendation == nil || res.ConstraintsUsed == nil {
		t.Fatalf("%+v", res)
	}
	errCode(t, e.c.Do("POST", "/api/v1/scenario-comparisons", map[string]any{"amount_cents": 0, "desired_date": "2026-10-30"}), 422, "validation_failed")
}

func TestOp_listCreditPurchases(t *testing.T) {
	e := setup(t)
	var pg api.CreditPurchasePage
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credit-purchases?family=external&period=2026-10&person_id=%d", e.karen), nil), 200, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].PaymentsTotal != 4 || *pg.Items[0].BalanceCents != 389500 || len(pg.Items[0].Shares) != 2 {
		t.Fatalf("%+v", pg.Items)
	}
	must(t, e.c.Do("GET", "/api/v1/credit-purchases?status=cancelled", nil), 200, &pg)
	if len(pg.Items) != 0 {
		t.Fatalf("%+v", pg.Items)
	}
	errCode(t, e.c.Do("GET", "/api/v1/credit-purchases?family=otra", nil), 400, "bad_request")
}

func TestOp_createCreditPurchase(t *testing.T) {
	e := setup(t)
	var pg api.CreditPurchasePage
	must(t, e.c.Do("GET", "/api/v1/credit-purchases", nil), 200, &pg)
	p := pg.Items[0]
	if p.TotalCents != 389500 || p.InterestSource != "quoted" || p.FirstPaymentDate != "2026-10-30" || *p.RateNum != 80050 || *p.RateDen != 309450 || p.RatePct != nil {
		t.Fatalf("%+v", p)
	}
	sp := e.scheduled(t, fmt.Sprintf("?credit_purchase_id=%d", e.buy))
	if len(sp) != 4 || sp[0].AmountCents != 97375 || sp[3].DueDate != "2026-12-15" || sp[3].PrincipalCents != 77361 {
		t.Fatalf("%+v", sp)
	}
	body := e.purchaseBody()
	body["shares"] = []map[string]any{{"person_id": e.axel, "principal_cents": 1}}
	errCode(t, e.c.Do("POST", "/api/v1/credit-purchases", body, "Idempotency-Key", key()), 422, "validation_failed")
	body = e.purchaseBody()
	body["principal_cents"] = 97.5
	errCode(t, e.c.Raw("POST", "/api/v1/credit-purchases", body, "Idempotency-Key", key()), 400, "bad_request")
}

func TestOp_previewCreditPurchase(t *testing.T) {
	e := setup(t)
	body := e.purchaseBody()
	delete(body, "interest_cents") // sin cotización: tasa de la tabla por plazo (term_table)
	var pv api.CreditPurchasePreview
	must(t, e.c.Do("POST", "/api/v1/credit-purchases/preview", body), 200, &pv)
	if pv.InterestCents != 80050 || pv.TotalCents != 389500 || len(pv.Schedule) != 4 || len(pv.Shares) != 2 {
		t.Fatalf("%+v", pv)
	}
	body["term_count"] = 6
	errCode(t, e.c.Do("POST", "/api/v1/credit-purchases/preview", body), 422, "validation_failed")
}

func TestOp_getCreditPurchase(t *testing.T) {
	e := setup(t)
	var p api.CreditPurchase
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credit-purchases/%d", e.buy), nil), 200, &p)
	if p.ID != e.buy || *p.PaymentsMade != 0 {
		t.Fatalf("%+v", p)
	}
	errCode(t, e.c.Do("GET", "/api/v1/credit-purchases/999", nil), 404, "not_found")
}

func TestOp_updateCreditPurchase(t *testing.T) {
	e := setup(t)
	var p api.CreditPurchase
	must(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credit-purchases/%d", e.buy), map[string]any{"description": "Despensa 1", "notes": nil}), 200, &p)
	if *p.Description != "Despensa 1" {
		t.Fatalf("%+v", p)
	}
	errCode(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credit-purchases/%d", e.buy), map[string]any{"status": "otro"}), 422, "validation_failed")
	errCode(t, e.c.Do("PATCH", "/api/v1/credit-purchases/999", map[string]any{"notes": "x"}), 404, "not_found")
}

func TestOp_deleteCreditPurchase(t *testing.T) {
	e := setup(t)
	sp := e.scheduled(t, fmt.Sprintf("?credit_purchase_id=%d", e.buy))
	e.pay(t, sp[0].ID, e.karen, 1000, 201)
	errCode(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-purchases/%d", e.buy), nil), 409, "has_payments")
	other := e.purchase(t, 201)
	must(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-purchases/%d", other.ID), nil), 204, nil)
	errCode(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-purchases/%d", other.ID), nil), 404, "not_found")
}

func TestOp_listScheduledPayments(t *testing.T) {
	e := setup(t)
	sp := e.scheduled(t, "")
	if len(sp) != 4 || *sp[0].Status != "pending" || len(sp[0].Shares) != 2 || *sp[3].BalanceAfterCents != 0 {
		t.Fatalf("%+v", sp)
	}
	if got := e.scheduled(t, "?status=paid"); len(got) != 0 {
		t.Fatalf("%+v", got)
	}
	errCode(t, e.c.Do("GET", "/api/v1/scheduled-payments?status=otro", nil), 400, "bad_request")
}

func TestOp_markScheduledPaymentPaid(t *testing.T) {
	e := setup(t)
	sp := e.scheduled(t, "")[0]
	for _, sh := range sp.Shares {
		m := e.pay(t, sp.ID, sh.PersonID, sh.AmountCents, 201)
		if m.AmountCents != sh.AmountCents || *m.ScheduledPaymentID != sp.ID {
			t.Fatalf("%+v", m)
		}
	}
	if got := e.scheduled(t, "?status=paid"); len(got) != 1 || *got[0].Status != "paid" {
		t.Fatalf("%+v", got)
	}
	errCode(t, e.c.Do("POST", fmt.Sprintf("/api/v1/scheduled-payments/%d/payments", sp.ID), map[string]any{"payer_person_id": e.axel,
		"paid_date": "2026-10-12", "amount_cents": 1}, "Idempotency-Key", key()), 409, "conflict")
	errCode(t, e.c.Do("POST", fmt.Sprintf("/api/v1/scheduled-payments/%d/payments", sp.ID), map[string]any{"payer_person_id": e.ximena,
		"paid_date": "2026-10-12", "amount_cents": 1}, "Idempotency-Key", key()), 422, "validation_failed")
	e.pay(t, 999, e.axel, 1, 404)
}

func TestOp_deleteCreditPayment(t *testing.T) {
	e := setup(t)
	sp := e.scheduled(t, "")[0]
	m := e.pay(t, sp.ID, e.karen, 500, 201)
	must(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-payments/%d", m.ID), nil), 204, nil)
	errCode(t, e.c.Do("DELETE", fmt.Sprintf("/api/v1/credit-payments/%d", m.ID), nil), 404, "not_found")
}

// Deuda viva «hoy» (M-12 a hoy, K-3 y P-17; PROVISIONAL, bloqueo S6-deuda-viva-hoy) y M-19 (K-3, C1 v1.5): una compra
// programada a futuro no cuenta hoy, pero sí en el cierre del mes de su compra; la cifra es la misma que getLiveDebt
// (motor S2) a la misma fecha. Pagos del mes: cuántos y la parte del hogar que falta pagar.
func TestOp_getCreditTotals_DeudaVivaHoy(t *testing.T) {
	s := testutil.NewServer(t, creditos.Register, presupuesto.Register)
	e := &env{s: s, c: s.SessionFor(t, testutil.AxelUser), axel: s.Fx.AxelPersonID, karen: s.Fx.KarenPersonID}
	var cr api.Credit
	must(t, e.c.Do("POST", "/api/v1/credits", map[string]any{"name": "Kueski Pay", "kind": "revolving_line", "holder_person_id": e.axel,
		"credit_limit_cents": 1000000, "payment_frequency": "fortnightly", "first_payment_min_days": 15}, "Idempotency-Key", key()), 201, &cr)
	e.kueski = cr.ID
	must(t, e.c.Do("POST", fmt.Sprintf("/api/v1/credits/%d/term-rates", e.kueski), map[string]any{"term_count": 4, "rate_num": 80050,
		"rate_den": 309450, "source": "quote", "quote_date": "2026-10-10", "rate_decimals": 4}, "Idempotency-Key", key()), 201, nil)
	e.purchase(t, 201)
	future := e.purchaseBody()
	future["purchase_date"] = "2026-12-10" // despensa programada (como las 2 a 4 de la semilla, S6-semilla-despensas-futuras)
	must(t, e.c.Do("POST", "/api/v1/credit-purchases", future, "Idempotency-Key", key()), 201, nil)
	live := func(date string) int64 {
		var ld api.LiveDebt
		must(t, e.c.Do("GET", "/api/v1/debt/live?date="+date, nil), 200, &ld)
		return ld.TotalCents
	}
	var tt api.CreditTotals
	must(t, e.c.Do("GET", "/api/v1/credits/totals", nil), 200, &tt)
	if tt.LiveDebtCents != 389500 || tt.LiveDebtCents != live("2026-10-12") {
		t.Fatalf("deuda viva hoy %d; getLiveDebt hoy %d; esperaba 389500 (sin la compra del 10-dic)", tt.LiveDebtCents, live("2026-10-12"))
	}
	if tt.MonthPaymentCount == nil || *tt.MonthPaymentCount != 1 || tt.MonthPendingHouseholdCents == nil ||
		*tt.MonthPendingHouseholdCents != tt.MonthHouseholdCents || tt.MonthHouseholdCents == 0 {
		t.Fatalf("pagos del mes: %+v", tt)
	}
	var c api.Credit
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/%d", e.kueski), nil), 200, &c)
	if c.LiveDebtCents == nil || *c.LiveDebtCents != 389500 {
		t.Fatalf("deuda viva de la línea: %v", c.LiveDebtCents)
	}
	var dt api.DebtTrend
	must(t, e.c.Do("GET", "/api/v1/credits/debt-trend", nil), 200, &dt)
	if dt.Months[0].LiveDebtCents != 292125 || dt.Months[1].LiveDebtCents != 97375 || dt.Months[1].LiveDebtCents != live("2026-11-30") ||
		dt.Months[2].LiveDebtCents <= dt.Months[1].LiveDebtCents || dt.Months[2].LiveDebtCents != live("2026-12-31") {
		t.Fatalf("deuda viva proyectada: oct %d, nov %d, dic %d (getLiveDebt 30-nov %d, 31-dic %d)", dt.Months[0].LiveDebtCents,
			dt.Months[1].LiveDebtCents, dt.Months[2].LiveDebtCents, live("2026-11-30"), live("2026-12-31"))
	}
	// Abono parcial de la parte del hogar al pago del 30-oct: falta solo el resto; la parte de Karen no lo cubre.
	var oct int64
	for _, p := range e.scheduled(t, "") {
		if p.DueDate == "2026-10-30" {
			oct = p.ID
		}
	}
	if oct == 0 {
		t.Fatal("no está el pago programado del 30-oct")
	}
	e.pay(t, oct, e.karen, 10000, 201)
	e.pay(t, oct, e.axel, 20000, 201)
	var t2 api.CreditTotals
	must(t, e.c.Do("GET", "/api/v1/credits/totals", nil), 200, &t2)
	if *t2.MonthPendingHouseholdCents != t2.MonthHouseholdCents-20000 || *t2.MonthPaymentCount != 1 {
		t.Fatalf("falta del hogar tras abonar 200.00: %d (parte del hogar %d)", *t2.MonthPendingHouseholdCents, t2.MonthHouseholdCents)
	}
}

// Línea archivada con pagos pendientes (hallazgo M-1 de REVISION-backend-cierre-r1; PROVISIONAL, bloqueo
// S6-tablero-lineas-archivadas, punto 5, y S6-deuda-viva-hoy): getCreditTotals (K-3, P-17) y getDebtTrend (G-5) cuentan
// todas las líneas, activas o archivadas, como getLiveDebt (motor S2) y el barómetro (RN-30). Una sola cifra hoy.
func TestOp_getCreditTotals_LineaArchivada(t *testing.T) {
	s := testutil.NewServer(t, creditos.Register, presupuesto.Register, tablero.Register)
	e := &env{s: s, c: s.SessionFor(t, testutil.AxelUser), axel: s.Fx.AxelPersonID, karen: s.Fx.KarenPersonID}
	line := func(name string) int64 {
		var cr api.Credit
		must(t, e.c.Do("POST", "/api/v1/credits", map[string]any{"name": name, "kind": "revolving_line", "holder_person_id": e.axel,
			"credit_limit_cents": 1000000, "payment_frequency": "fortnightly", "first_payment_min_days": 15}, "Idempotency-Key", key()), 201, &cr)
		must(t, e.c.Do("POST", fmt.Sprintf("/api/v1/credits/%d/term-rates", cr.ID), map[string]any{"term_count": 4, "rate_num": 80050,
			"rate_den": 309450, "source": "quote", "quote_date": "2026-10-10", "rate_decimals": 4}, "Idempotency-Key", key()), 201, nil)
		e.kueski = cr.ID
		e.purchase(t, 201) // 3,895.00 de deuda hoy; pago del 30-oct 973.75, parte del hogar 598.82 (R-13)
		return cr.ID
	}
	line("Kueski Pay")
	archived := line("Línea archivada")
	get := func(path string) string {
		res := e.c.Do("GET", path, nil)
		must(t, res, 200, nil)
		return string(res.Body)
	}
	beforeTotals, beforeTrend := get("/api/v1/credits/totals"), get("/api/v1/credits/debt-trend")
	must(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credits/%d", archived), map[string]any{"is_active": false}), 200, nil)
	// Archivar no cambia ninguna cifra de K-3, P-17 ni G-5.
	if a, b := get("/api/v1/credits/totals"), get("/api/v1/credits/debt-trend"); a != beforeTotals || b != beforeTrend {
		t.Fatalf("archivar cambió las cifras:\ntotales antes %s\ndespués %s\nserie antes %s\ndespués %s", beforeTotals, a, beforeTrend, b)
	}
	live := func(date string) int64 {
		var ld api.LiveDebt
		must(t, e.c.Do("GET", "/api/v1/debt/live?date="+date, nil), 200, &ld)
		return ld.TotalCents
	}
	var tt api.CreditTotals
	must(t, e.c.Do("GET", "/api/v1/credits/totals", nil), 200, &tt)
	if tt.LiveDebtCents != 779000 || tt.LiveDebtCents != live("2026-10-12") {
		t.Fatalf("deuda viva hoy %d; getLiveDebt hoy %d; esperaba 779000 (2 × 389500, con la línea archivada)", tt.LiveDebtCents, live("2026-10-12"))
	}
	var d api.DashboardDebtCapacity
	must(t, e.c.Do("GET", "/api/v1/dashboard/debt-capacity", nil), 200, &d)
	// Pagos del mes: los dos pagos del 30-oct (2 × 973.75), nada pagado; la misma cifra que el barómetro (RN-30).
	if tt.MonthTotalCents != 2*97375 || tt.MonthHouseholdCents == 0 || tt.MonthHouseholdCents != d.PaymentsCents ||
		tt.MonthTotalCents != d.PaymentsTotalCents || tt.MonthPaymentCount == nil || *tt.MonthPaymentCount != 2 || d.PaymentCount != 2 ||
		tt.MonthPendingHouseholdCents == nil || *tt.MonthPendingHouseholdCents != tt.MonthHouseholdCents {
		t.Fatalf("pagos del mes: K-3 hogar %d total %d; barómetro %d %d (%d pagos)", tt.MonthHouseholdCents, tt.MonthTotalCents,
			d.PaymentsCents, d.PaymentsTotalCents, d.PaymentCount)
	}
	var dt api.DebtTrend
	must(t, e.c.Do("GET", "/api/v1/credits/debt-trend", nil), 200, &dt)
	if len(dt.Months) != 12 || dt.Months[0].LiveDebtCents != 2*292125 {
		t.Fatalf("deuda viva proyectada: %+v", dt.Months)
	}
	for _, m := range dt.Months {
		mt, err := time.Parse("2006-01", string(m.Month))
		if err != nil {
			t.Fatal(err)
		}
		end := mt.AddDate(0, 1, -1).Format("2006-01-02")
		if l := live(end); m.LiveDebtCents != l {
			t.Errorf("cierre de %s: getDebtTrend %d, getLiveDebt %d", m.Month, m.LiveDebtCents, l)
		}
	}
	// Con credit_id, la línea nombrada cuenta aunque esté archivada.
	must(t, e.c.Do("GET", fmt.Sprintf("/api/v1/credits/debt-trend?credit_id=%d", archived), nil), 200, &dt)
	if dt.Months[0].LiveDebtCents != 292125 || len(dt.Months[0].Credits) != 1 || dt.Months[0].Credits[0].CreditID != archived {
		t.Fatalf("deuda viva proyectada de la línea archivada: %+v", dt.Months[0])
	}
}

// A-1 de REVISION-backend-cierre-r2: el comparador recibe solo las líneas activas. Una línea archivada con límite, con
// tasa más baja y con su compra viva, no se ofrece ni se recomienda, y el resultado es el mismo que con solo la línea
// activa (misma cuenta de evaluados, mismos válidos y descartados, misma recomendación y mismo score_bp, porque la
// normalización de S3-14 es entre los escenarios generados). getCreditTotals sí la sigue contando (RN-30, M-1).
func TestOp_compareScenarios_LineaArchivada(t *testing.T) {
	s := testutil.NewServer(t, creditos.Register, presupuesto.Register) // presupuesto: getLiveDebt
	e := &env{s: s, c: s.SessionFor(t, testutil.AxelUser), axel: s.Fx.AxelPersonID, karen: s.Fx.KarenPersonID}
	// Una línea revolvente con límite, tasa a 4 quincenas rateNum/309450 y la despensa del 10-oct (R-13).
	line := func(name string, rateNum int64) int64 {
		var cr api.Credit
		must(t, e.c.Do("POST", "/api/v1/credits", map[string]any{"name": name, "kind": "revolving_line", "holder_person_id": e.axel,
			"credit_limit_cents": 1000000, "payment_frequency": "fortnightly", "first_payment_min_days": 15}, "Idempotency-Key", key()), 201, &cr)
		must(t, e.c.Do("POST", fmt.Sprintf("/api/v1/credits/%d/term-rates", cr.ID), map[string]any{"term_count": 4, "rate_num": rateNum,
			"rate_den": 309450, "source": "quote", "quote_date": "2026-10-10", "rate_decimals": 4}, "Idempotency-Key", key()), 201, nil)
		e.kueski = cr.ID
		e.purchase(t, 201)
		return cr.ID
	}
	req := map[string]any{"amount_cents": 250000, "desired_date": "2026-10-30", "merchant_name": "Chedraui",
		"min_available_buffer_cents": 0, "max_live_purchases": 5, "zero_debt_by_date": "2027-12-31", "date_window_days": 1}
	compare := func() (string, api.ScenarioResult) {
		var res api.ScenarioResult
		r := e.c.Do("POST", "/api/v1/scenario-comparisons", req)
		must(t, r, 200, &res)
		return string(r.Body), res
	}
	onLine := func(res api.ScenarioResult, id int64) (n int) {
		for _, v := range res.Valid {
			if v.CreditID == id {
				n++
			}
		}
		for _, r := range res.Rejected {
			if r.Scenario.CreditID == id {
				n++
			}
		}
		return n
	}
	kueski := line("Kueski Pay", 80050) // 25.8685 %
	baseBody, base := compare()
	if base.Status != "ok" || base.Recommendation == nil || base.Recommendation.CreditID != kueski || onLine(base, kueski) != int(base.EvaluatedCount) {
		t.Fatalf("solo la línea activa: %s", baseBody)
	}
	// Segunda línea con tasa más baja y la misma compra: mientras está activa, el comparador la ofrece y la recomienda.
	archived := line("Línea archivada", 40000)
	activeBody, active := compare()
	if onLine(active, archived) == 0 || active.Recommendation == nil || active.Recommendation.CreditID != archived ||
		active.EvaluatedCount <= base.EvaluatedCount {
		t.Fatalf("con las dos líneas activas esperaba escenarios y recomendación en la línea %d: %s", archived, activeBody)
	}
	must(t, e.c.Do("PATCH", fmt.Sprintf("/api/v1/credits/%d", archived), map[string]any{"is_active": false}), 200, nil)
	afterBody, after := compare()
	if n := onLine(after, archived); n != 0 {
		t.Fatalf("la línea archivada %d sigue en %d escenarios (válidos o descartados): %s", archived, n, afterBody)
	}
	if after.Recommendation == nil || after.Recommendation.CreditID == archived {
		t.Fatalf("recomendación en la línea archivada: %s", afterBody)
	}
	ra, rb := after.Recommendation, base.Recommendation
	sameRec := ra.CreditID == rb.CreditID && ra.TermCount == rb.TermCount && ra.PurchaseDate == rb.PurchaseDate &&
		ra.ScoreBp != nil && rb.ScoreBp != nil && *ra.ScoreBp == *rb.ScoreBp
	if after.EvaluatedCount != base.EvaluatedCount || !sameRec || afterBody != baseBody {
		t.Fatalf("con la línea archivada el resultado no es el de solo la línea activa:\nsolo activa %s\ncon archivada %s", baseBody, afterBody)
	}
	// La misma línea archivada sí cuenta en la deuda viva de hoy y en los pagos del mes (RN-30, getLiveDebt).
	var tt api.CreditTotals
	must(t, e.c.Do("GET", "/api/v1/credits/totals", nil), 200, &tt)
	var ld api.LiveDebt
	must(t, e.c.Do("GET", "/api/v1/debt/live?date=2026-10-12", nil), 200, &ld)
	if tt.LiveDebtCents != 2*389500 || tt.LiveDebtCents != ld.TotalCents || tt.MonthTotalCents != 2*97375 {
		t.Fatalf("getCreditTotals: deuda viva %d (getLiveDebt %d), pagos del mes %d; esperaba 779000 y 194750 con la línea archivada",
			tt.LiveDebtCents, ld.TotalCents, tt.MonthTotalCents)
	}
}
