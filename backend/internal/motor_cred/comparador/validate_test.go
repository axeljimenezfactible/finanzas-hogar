package comparador

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	mc "finanzas-hogar/internal/motor_cred"
)

func TestValidacion(t *testing.T) {
	in, args, _ := loadFixture(t, "caso-1")
	neg, negInt := mc.Cents(-1), -1
	muts := map[string]func(r *Request){
		"monto 0":        func(r *Request) { r.Amount = 0 },
		"sin today":      func(r *Request) { r.Today = mc.Date{} },
		"pesos en 0":     func(r *Request) { r.Weights = Weights{} },
		"peso negativo":  func(r *Request) { r.Weights.MaxAvailable = -1 },
		"colchón < 0":    func(r *Request) { r.Override.MinAvailableBuffer = &neg },
		"vivas < 0":      func(r *Request) { r.Override.MaxLivePurchases = &negInt },
		"ventana < 0":    func(r *Request) { r.Global.DateWindowDays = &negInt },
		"canasta costo":  func(r *Request) { r.Baskets[0].Cost = 0 },
		"tope negativo":  func(r *Request) { r.Caps[0].Cap = -5 },
		"tope sin dueño": func(r *Request) { r.Caps[0].Person = "" },
		"partes suma":    func(r *Request) { r.Shares[0].Principal-- },
		"dos titulares":  func(r *Request) { r.Shares = append(r.Shares, mc.ShareInput{Code: "karen"}) },
		"parte negativa": func(r *Request) {
			r.Shares = append(r.Shares, mc.ShareInput{Code: "karen", External: true, Principal: -1})
		},
		"crédito doble": func(r *Request) { r.Credits = append(r.Credits, r.Credits[0]) },
		"ref reservada": func(r *Request) {
			cs := r.Credits[0]
			cs.Purchases = append([]mc.PurchaseState{}, cs.Purchases...)
			cs.Purchases[0].Ref = scenarioRef
			r.Credits = append([]mc.CreditState{cs}, r.Credits[1:]...)
		},
	}
	for name, mut := range muts {
		req, _ := requestOf(t, in, args)
		req.Baskets = append([]BasketQuote{}, req.Baskets...)
		req.Caps = append([]PersonCap{}, req.Caps...)
		req.Shares = append([]mc.ShareInput{}, req.Shares...)
		mut(&req)
		if _, err := Compare(context.Background(), req); !errors.Is(err, mc.ErrInvalid) {
			t.Errorf("%s: se espera ErrInvalid, salió %v", name, err)
		}
	}
	for name, a := range map[string]string{
		"campo desconocido": `{"timeout_seconds": 5, "request": {}, "otro": 1}`,
		"sin request":       `{"timeout_seconds": 5}`,
		"timeout negativo":  `{"timeout_seconds": -1, "request": {"amount_cents": 1}}`,
		"request inválido":  `{"timeout_seconds": 5, "request": {"amount_cents": "x"}}`,
		"monto 0":           `{"timeout_seconds": 5, "request": {"amount_cents": 0}}`,
	} {
		if out, code := EvalComparar(in, json.RawMessage(a)); out != nil || code == "" {
			t.Errorf("%s: se espera código de error, salió %q", name, code)
		}
	}
}

// TestGeneracionTiendas: tienda fija por id o por comercio, sin tiendas (objetivo 4 nulo → 1),
// préstamo sin candidatos y promoción de tarjeta con tasa ≠ 0 que no se genera.
func TestGeneracionTiendas(t *testing.T) {
	in, args, _ := loadFixture(t, "caso-1")
	req, _ := requestOf(t, in, args)
	two := int64(2)
	req.StoreID = &two
	res := compareT(t, req)
	for _, s := range res.Valid {
		if s.StoreID == nil || *s.StoreID != 2 || s.ValuePerPesoMilli == nil {
			t.Fatalf("tienda fija por id: %s", keyOf(s))
		}
	}
	req.StoreID, req.Merchant = nil, " tienda cara "
	res = compareT(t, req)
	if len(res.Valid) == 0 || res.Valid[0].StoreID != nil || res.Valid[0].ValuePerPesoMilli == nil || res.Optima.LowestInterest.PromotionID == nil {
		t.Errorf("tienda por comercio sin id: la promoción aplica por comercio (RN-10) y la canasta se encuentra por comercio")
	}
	req.Merchant, req.Baskets = "", nil
	res = compareT(t, req)
	for _, s := range res.Valid {
		if s.StoreID != nil || s.ValuePerPesoMilli != nil {
			t.Fatalf("sin tiendas: %s", keyOf(s))
		}
	}
	// préstamo con límite: sin plazos en §8.2, no aporta candidatos
	rate := mc.Rate{Num: 12, Den: 100}
	loan := mc.CreditState{Credit: mc.Credit{ID: 9, Kind: mc.KindLoan, Holder: "axel", Limit: ptr(mc.Cents(900000)), Frequency: "monthly", AnnualRate: &rate}}
	// tarjeta con promoción de tasa ≠ 0 (meses con intereses: sin regla) y sin promociones vigentes
	sd, pd := 5, 20
	card := mc.CreditState{Credit: mc.Credit{ID: 8, Kind: mc.KindCard, Holder: "axel", Limit: ptr(mc.Cents(900000)), Frequency: "monthly",
		StatementDay: &sd, PaymentDueDay: &pd, Promotions: []mc.Promotion{{ID: 7, CreditID: ptr(int64(8)), Merchant: "x", TermCount: 3,
			Rate: mc.Rate{Num: 1, Den: 10}, ValidFrom: req.Today, ValidTo: req.DesiredDate.AddDays(5), Active: true}}}}
	req.Credits = append(append([]mc.CreditState{}, req.Credits...), loan, card)
	req.Merchant = "x"
	res = compareT(t, req)
	for _, s := range append(res.Valid, rejectedScenarios(res)...) {
		if s.CreditID == 9 || (s.CreditID == 8 && s.TermCount != 1) {
			t.Errorf("no debía generarse %s", keyOf(s))
		}
	}
	// ningún válido → no_valid_scenario sin óptimos
	big := mc.Cents(1 << 40)
	req.Override.MinAvailableBuffer = &big
	if r := compareT(t, req); r.Status != StatusNoValid || r.Optima != nil || r.Recommendation != nil || len(r.Provisional) == 0 {
		t.Errorf("sin válidos: %s", r.Status)
	}
	// sin candidatos (ventana antes de today)
	req.DesiredDate = req.Today.AddDays(-30)
	if r := compareT(t, req); r.Status != StatusNoValid || r.EvaluatedCount != 0 {
		t.Errorf("sin candidatos: %s %d", r.Status, r.EvaluatedCount)
	}
}

func TestContextoCancelado(t *testing.T) {
	in, args, _ := loadFixture(t, "sin-restricciones")
	req, _ := requestOf(t, in, args)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Compare(ctx, req); !errors.Is(err, mc.ErrTimeout) {
		t.Errorf("contexto cancelado con faltantes: %v", err)
	}
	in, args, _ = loadFixture(t, "caso-2")
	req, _ = requestOf(t, in, args)
	if _, err := Compare(ctx, req); !errors.Is(err, mc.ErrTimeout) {
		t.Errorf("contexto cancelado: %v", err)
	}
}
