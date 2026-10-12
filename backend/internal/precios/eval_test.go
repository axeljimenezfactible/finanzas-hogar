package precios

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

const evalCatalog = `"catalog": {"groups": [{"code": "abarrotes", "name": "Abarrotes", "purchase_channel": "chain"}],
 "products": [{"code": "arroz", "name": "Arroz", "base_unit": "g", "group": "abarrotes", "grams_per_unit_milli": null},
  {"code": "huevo", "name": "Huevo", "base_unit": "piece", "group": null, "grams_per_unit_milli": null},
  {"code": "frijol", "name": "Frijol", "base_unit": "g", "group": "abarrotes", "grams_per_unit_milli": null}],
 "presentations": [{"code": "arroz-1kg", "product": "arroz", "name": "Arroz 1 kg", "size_milli": 1000000},
  {"code": "arroz-750", "product": "arroz", "name": "Arroz 750 g", "size_milli": 750000},
  {"code": "huevo-12", "product": "huevo", "name": "Huevo 12", "size_milli": 12000},
  {"code": "frijol-1kg", "product": "frijol", "name": "Frijol 1 kg", "size_milli": 1000000}]}`

const evalStores = `"stores": [
 {"code": "walmart", "kind": "chain", "parent": null, "name": "Walmart", "profeco_name": "WALMART", "municipality": null, "visit_weekday": null, "requires_membership": false, "has_membership": false, "membership_fee_cents": null, "assisted_reading": true, "groups": []},
 {"code": "bodega", "kind": "chain", "parent": null, "name": "Bodega Aurrera", "profeco_name": null, "municipality": null, "visit_weekday": null, "requires_membership": false, "has_membership": false, "membership_fee_cents": null, "assisted_reading": false, "groups": []},
 {"code": "sams", "kind": "chain", "name": "Sam's Club", "requires_membership": true, "has_membership": false, "membership_fee_cents": 100000, "assisted_reading": false, "groups": []},
 {"code": "tianguis", "kind": "market", "name": "Tianguis", "visit_weekday": 2, "groups": ["frutas-verduras"]}]`

func evalInput(extra string) json.RawMessage {
	s := `{"today": "2026-10-12", "settings": {"stale_after_days": 30, "zero_weight_after_days": 90, "assisted_reading_enabled": true}, ` +
		evalCatalog + `, ` + evalStores + `, "prices": [
 {"presentation": "arroz-1kg", "store": "walmart", "price_cents": 2800, "observed_date": "2026-10-10", "source": "manual", "is_promotional": false, "source_ref": null, "created_at": "2026-10-10T18:00:00Z"},
 {"presentation": "arroz-1kg", "store": "bodega", "price_cents": 2900, "observed_date": "2026-10-10", "source": "ticket_manual", "is_promotional": false, "source_ref": "gasto 1", "created_at": "2026-10-10T18:00:00Z"},
 {"presentation": "arroz-1kg", "store": "bodega", "price_cents": 1000, "observed_date": "2026-10-12", "source": "manual", "is_promotional": true, "source_ref": null, "created_at": "2026-10-12T18:00:00Z"},
 {"presentation": "arroz-1kg", "store": "bodega", "price_cents": 2000, "observed_date": "2026-10-01", "source": "profeco", "is_promotional": false, "source_ref": "QQP 2026 a.csv:2", "created_at": "2026-10-11T18:00:00Z"},
 {"presentation": "arroz-1kg", "store": "sams", "price_cents": 2500, "observed_date": "2026-10-10", "source": "ticket_ocr", "is_promotional": false, "source_ref": null, "created_at": "2026-10-10T18:00:00Z"}],
 "shopping_list": [{"product": "arroz", "needed_milli": 2000000, "quantity_milli": null, "recommended_presentation": "arroz-1kg"}],
 "grocery_purchases": [{"date": "2026-09-01", "store": "walmart"}, {"date": "2026-10-01", "store": "bodega"}],
 "terms": [{"site": "walmart", "verdict": "permitido", "url": "https://ejemplo.invalid/terminos", "read_date": "2026-10-11"}, {"site": "bodega_aurrera", "verdict": "", "url": "", "read_date": ""}]` + extra + `}`
	return json.RawMessage(s)
}

func run(t *testing.T, in json.RawMessage, fn, args string) map[string]any {
	t.Helper()
	out, err := Eval(in, fn, json.RawMessage(args))
	if err != nil {
		t.Fatalf("%s(%s): %v", fn, args, err)
	}
	b, _ := json.Marshal(out)
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func code(t *testing.T, in json.RawMessage, fn, args, want string) {
	t.Helper()
	_, err := Eval(in, fn, json.RawMessage(args))
	if err == nil || Code(err) != want {
		t.Errorf("%s(%s): error %v (%s), se esperaba %s", fn, args, err, Code(err), want)
	}
}

func js(v any) string { b, _ := json.Marshal(v); return string(b) }

func TestEvalComparativaCanon(t *testing.T) {
	t.Chdir(repo)
	in := json.RawMessage(`{"today": "2026-10-10", "catalog_file": "datos/semilla/productos.json", "report_file": "datos/semilla/data.json"}`)
	m := run(t, in, "compare_carts", `{"base": "bodega", "others": ["walmart"]}`)
	st := m["stores"].([]any)
	eq(t, "stores[0]", js(st[0]), `{"assumes_membership":false,"difference_vs_cheapest_cents":0,"identical_count":33,"identical_total_cents":277350,"missing_count":0,"piece_count":61,"stale_count":0,"store":"bodega","substitute_count":3,"total_cents":309450}`)
	eq(t, "stores[1].total", js(st[1].(map[string]any)["difference_vs_cheapest_cents"]), "13850")
	eq(t, "identical_pieces", js(m["identical_pieces"]), "55")
	eq(t, "winners", js(m["winners"]), `[{"count":22,"store":"bodega"},{"count":3,"store":"walmart"}]`)
	code(t, in, "compare_carts", `{"base": "nada", "others": ["walmart"]}`, "invalid_argument")
	code(t, in, "compare_carts", `{"base": "bodega", "others": ["nada"]}`, "invalid_argument")
	code(t, in, "compare_carts", `{"base": "bodega", "others": []}`, "invalid_argument")
	code(t, in, "compare_carts", `{"others": []}`, "missing_field")
	code(t, json.RawMessage(`{"today": "2026-10-10", "catalog_file": "datos/semilla/productos.json", "report_file": "datos/semilla/data.json", "carts": [{"code": "x", "store": "walmart", "kind": "cart", "observed_date": "2026-10-10", "source": "manual", "lines": []}]}`), "compare_carts", `{}`, "invalid_argument")
	code(t, json.RawMessage(`{"today": "2026-10-10", "catalog_file": "../productos.json"}`), "compare_basket", `{}`, "invalid_argument")
	code(t, json.RawMessage(`{"today": "2026-10-10", "catalog_file": "/etc/productos.json"}`), "compare_basket", `{}`, "invalid_argument")
	code(t, json.RawMessage(`{"today": "2026-10-10", "catalog_file": "datos/semilla/productos.json", "report_file": "datos/nada.json"}`), "compare_basket", `{}`, "invalid_argument")
	code(t, json.RawMessage(`{"today": "2026-10-10", "catalog_file": "datos/semilla/nada.json"}`), "compare_basket", `{}`, "invalid_argument")
}

func TestEvalFunciones(t *testing.T) {
	in := evalInput("")
	eq(t, "base_unit_price", js(run(t, in, "base_unit_price", `{"presentation": "arroz-750", "price_cents": 2200}`)),
		`{"base_unit_price_cents":2933,"exact":{"den":3,"num":8800},"per":"kg"}`)
	eq(t, "freshness", js(run(t, in, "freshness", `{"observed_date": "2026-09-12", "source": "manual"}`)),
		`{"age_days":30,"competes":true,"is_stale":true,"weight":{"den":61,"num":60},"weight_display":"0.9836"}`)
	eq(t, "current_price", js(run(t, in, "current_price", `{"presentation": "arroz-1kg", "store": "bodega"}`)),
		`{"price":{"observed_date":"2026-10-10","price_cents":2900,"source":"ticket_manual"},"promotion":{"observed_date":"2026-10-12","price_cents":1000},"references":[{"observed_date":"2026-10-01","price_cents":2000,"source":"profeco"}]}`)
	eq(t, "quote_for", js(run(t, in, "quote_for", `{"product": "arroz"}`)),
		`{"quote":{"is_stale":false,"observed_date":"2026-10-10","presentation":"arroz-1kg","price_cents":2500,"product":"arroz","recommended":true,"size_milli":1000000,"store":"sams"}}`)
	eq(t, "quote_for con tienda", js(run(t, in, "quote_for", `{"product": "arroz", "store": "walmart", "presentation": "arroz-1kg"}`)["quote"].(map[string]any)["recommended"]), "false")
	eq(t, "quote_for sin precio", js(run(t, in, "quote_for", `{"product": "huevo"}`)), `{"quote":null}`)
	cb := run(t, in, "compare_basket", `{}`)
	eq(t, "compare_basket cheapest", js(cb["cheapest_store"]), `"sams"`)
	eq(t, "compare_basket estimated", js(cb["items"].([]any)[0].(map[string]any)["estimated_price_cents"]), "5000")
	m := run(t, in, "membership", `{"store": "sams"}`)
	eq(t, "membership ahorro", js(m["basket_savings_cents"]), "600")
	eq(t, "membership compras", js(m["purchases_per_year"]), "2")
	eq(t, "membership equilibrio dado", js(run(t, in, "membership", `{"store": "sams", "savings_per_purchase_cents": 12500}`)["break_even_purchases"]), "8")
	eq(t, "next_visit", js(run(t, in, "next_visit", `{"store": "tianguis", "from": "2026-10-10"}`)), `{"date":"2026-10-13"}`)
	eq(t, "stale_prices_for_plan", js(run(t, in, "stale_prices_for_plan", `{"products": ["arroz", "huevo"], "assigned_stores": {"arroz": "walmart"}, "grocery_date": "2026-11-09"}`)),
		`{"stale":[{"age_days":30,"observed_date":"2026-10-10","presentation":"arroz-1kg","product":"arroz","store":"walmart","weight":{"den":61,"num":60}}],"without_price":["huevo"]}`)
	eq(t, "assisted_can_enable", js(run(t, in, "assisted_can_enable", `{"store": "walmart"}`)), `{"can_enable":true}`)
	code(t, in, "assisted_can_enable", `{"store": "bodega"}`, "no_terms_verdict")
	code(t, in, "assisted_can_enable", `{"store": "sams"}`, "no_terms_verdict")
}

func TestEvalErrores(t *testing.T) {
	in := evalInput("")
	for _, c := range [][3]string{
		{"nada", `{}`, "invalid_argument"},
		{"compare_basket", `{"x": 1}`, "unknown_field"},
		{"base_unit_price", `{"presentation": "arroz-1kg"}`, "missing_field"},
		{"base_unit_price", `{"presentation": "nada", "price_cents": 1}`, "unknown_presentation"},
		{"base_unit_price", `{"presentation": "arroz-1kg", "price_cents": -1}`, "invalid_argument"},
		{"base_unit_price", `{"presentation": 5, "price_cents": 1}`, "invalid_argument"},
		{"freshness", `{"observed_date": "2026-13-01", "source": "manual"}`, "invalid_date"},
		{"freshness", `{"observed_date": "2026-10-13", "source": "manual"}`, "future_date"},
		{"freshness", `{"observed_date": "2026-10-10"}`, "missing_field"},
		{"freshness", `{"source": "manual"}`, "missing_field"},
		{"freshness", `{"observed_date": "2026-10-10", "source": "otra"}`, "invalid_argument"},
		{"freshness", `{"x": 1}`, "unknown_field"},
		{"current_price", `{"presentation": "arroz-1kg", "store": "nada"}`, "unknown_store"},
		{"current_price", `{"presentation": "nada", "store": "walmart"}`, "unknown_presentation"},
		{"current_price", `{"presentation": "arroz-1kg"}`, "missing_field"},
		{"current_price", `{"x": 1}`, "unknown_field"},
		{"quote_for", `{}`, "missing_field"},
		{"quote_for", `{"product": "nada"}`, "invalid_argument"},
		{"quote_for", `{"product": "arroz", "store": "nada"}`, "unknown_store"},
		{"quote_for", `{"x": 1}`, "unknown_field"},
		{"compare_carts", `{"x": 1}`, "unknown_field"},
		{"membership", `{}`, "missing_field"},
		{"membership", `{"store": "nada"}`, "unknown_store"},
		{"membership", `{"x": 1}`, "unknown_field"},
		{"next_visit", `{"store": "walmart", "from": "2026-10-10"}`, "invalid_argument"},
		{"next_visit", `{"store": "nada", "from": "2026-10-10"}`, "unknown_store"},
		{"next_visit", `{"from": "2026-10-10"}`, "missing_field"},
		{"next_visit", `{"store": "tianguis"}`, "missing_field"},
		{"next_visit", `{"x": 1}`, "unknown_field"},
		{"stale_prices_for_plan", `{"products": ["arroz"]}`, "missing_field"},
		{"stale_prices_for_plan", `{"products": ["arroz"], "assigned_stores": {"arroz": "nada"}, "grocery_date": "2026-11-09"}`, "unknown_store"},
		{"stale_prices_for_plan", `{"x": 1}`, "unknown_field"},
		{"assisted_can_enable", `{}`, "missing_field"},
		{"assisted_can_enable", `{"store": "nada"}`, "unknown_store"},
		{"assisted_can_enable", `{"x": 1}`, "unknown_field"},
	} {
		code(t, in, c[0], c[1], c[2])
	}
	for _, c := range [][2]string{
		{`{"x": 1}`, "unknown_field"},
		{`{"today": "2026-10-12"}`, "missing_field"},
		{`{"today": "12-oct-2026", ` + evalCatalog + `}`, "invalid_date"},
		{`{` + evalCatalog + `}`, "missing_field"},
		{`{"today": "2026-10-12", "catalog_file": "x", ` + evalCatalog + `}`, "invalid_argument"},
		{`{"today": "2026-10-12", "catalog": {"products": [{"code": "a", "base_unit": "oz"}]}}`, "invalid_argument"},
		{`{"today": "2026-10-12", "settings": {"stale_after_days": 0}, ` + evalCatalog + `}`, "invalid_argument"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "stores": [{"kind": "chain"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "stores": [{"code": "a", "kind": "otro"}]}`, "invalid_argument"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, ` + evalStores + `, "prices": [{"presentation": "arroz-1kg", "store": "walmart", "price_cents": 1, "source": "manual"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, ` + evalStores + `, "prices": [{"presentation": "arroz-1kg", "store": "walmart"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, ` + evalStores + `, "prices": [{"presentation": "arroz-1kg", "store": "walmart", "price_cents": 1, "observed_date": "2026-10-10", "source": "manual", "created_at": "ayer"}]}`, "invalid_date"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, ` + evalStores + `, "prices": [{"presentation": "arroz-1kg", "store": "walmart", "price_cents": 1, "observed_date": "2026-10-13", "source": "manual"}]}`, "future_date"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, ` + evalStores + `, "prices": [{"presentation": "arroz-1kg", "store": "walmart", "price_cents": 1, "observed_date": "2026-10-10", "source": "manual"}, {"presentation": "arroz-1kg", "store": "walmart", "price_cents": 2, "observed_date": "2026-10-10", "source": "manual"}]}`, "duplicate_price"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "carts": [{"code": "a"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "carts": [{"code": "a", "store": "w", "kind": "otro", "source": "manual"}]}`, "invalid_argument"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "carts": [{"code": "a", "store": "w", "kind": "cart", "source": "manual"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "carts": [{"code": "a", "store": "w", "kind": "cart", "source": "manual", "observed_date": "2026-10-10", "lines": [{"presentation": "arroz-1kg"}]}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "shopping_list": [{"needed_milli": 1}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "grocery_purchases": [{"date": "2026-10-01"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "grocery_purchases": [{"store": "w"}]}`, "missing_field"},
		{`{"today": "2026-10-12", ` + evalCatalog + `, "terms": [{"verdict": "x"}]}`, "missing_field"},
	} {
		code(t, json.RawMessage(c[0]), "compare_basket", `{}`, c[1])
	}
	if Code(errors.New("otro")) != "invalid_argument" || Code(ErrNotConfirmed) != "not_confirmed" || Code(ErrOverflow) != "overflow" {
		t.Error("Code")
	}
	if _, err := RepoPath("oraculo/fixtures/s5/_datos/a.csv", "oraculo/fixtures/s5/_datos"); err != nil {
		t.Error(err)
	}
	for _, p := range []string{"", "datos/a.csv", "oraculo/fixtures/s5/_datos/../x.csv"} {
		if _, err := RepoPath(p, "oraculo/fixtures/s5/_datos"); err == nil {
			t.Errorf("RepoPath(%q) aceptada", p)
		}
	}
}

func TestEvalCarritos(t *testing.T) {
	in := evalInput(`, "carts": [
 {"code": "w", "store": "walmart", "kind": "expense", "observed_date": "2026-10-10", "source": "manual", "lines": [
  {"presentation": "arroz-1kg", "quantity_milli": 2000, "unit_price_cents": 2800}, {"presentation": "huevo-12", "quantity_milli": 1000, "unit_price_cents": 4000}]},
 {"code": "b", "store": "bodega", "kind": "cart", "observed_date": "2026-09-01", "source": "ticket_manual", "lines": [
  {"presentation": "arroz-750", "quantity_milli": 3000, "unit_price_cents": 2200}, {"presentation": "frijol-1kg", "quantity_milli": 1000, "unit_price_cents": 3000}]}]`)
	m := run(t, in, "compare_carts", `{"base": "w", "others": ["b"]}`)
	st := m["stores"].([]any)
	eq(t, "empate de totales: gana la de menos viejos", js(st[0].(map[string]any)["store"]), `"walmart"`)
	eq(t, "Bodega", js(st[1]), `{"assumes_membership":false,"difference_vs_cheapest_cents":0,"identical_count":0,"identical_total_cents":0,"missing_count":1,"piece_count":4,"stale_count":2,"store":"bodega","substitute_count":1,"total_cents":9600}`)
	rows := m["rows"].([]any)
	eq(t, "renglones", len(rows), 3)
	eq(t, "sustituto", js(rows[0].(map[string]any)["winner_store"]), `"walmart"`)
	eq(t, "solo Walmart", js(rows[1].(map[string]any)["kind"]), `"only_one"`)
	bad := strings.Replace(string(in), `"source": "ticket_manual", "lines"`, `"source": "otra", "lines"`, 1)
	code(t, json.RawMessage(bad), "compare_carts", `{"base": "w", "others": ["b"]}`, "invalid_argument")
	bad = strings.Replace(string(in), `"store": "bodega", "kind": "cart"`, `"store": "nada", "kind": "cart"`, 1)
	code(t, json.RawMessage(bad), "compare_carts", `{"base": "w", "others": ["b"]}`, "unknown_store")
	bad = strings.Replace(string(in), `"observed_date": "2026-09-01"`, `"observed_date": "2026-10-20"`, 1)
	code(t, json.RawMessage(bad), "compare_carts", `{"base": "w", "others": ["b"]}`, "future_date")
	bad = strings.Replace(string(in), `{"presentation": "frijol-1kg", "quantity_milli": 1000`, `{"presentation": "nada", "quantity_milli": 1000`, 1)
	code(t, json.RawMessage(bad), "compare_carts", `{"base": "w", "others": ["b"]}`, "unknown_presentation")
	bad = strings.Replace(string(in), `{"presentation": "frijol-1kg", "quantity_milli": 1000`, `{"presentation": "frijol-1kg", "quantity_milli": 0`, 1)
	code(t, json.RawMessage(bad), "compare_carts", `{"base": "w", "others": ["b"]}`, "invalid_argument")
	bad = strings.Replace(string(in), `"code": "b"`, `"code": "w"`, 1)
	code(t, json.RawMessage(bad), "compare_carts", `{}`, "invalid_argument")
}

func TestCompraContraPrecios(t *testing.T) {
	b := newBook(t, []Store{stBodega, stWalmart, stTianguis}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("arroz-750", "walmart", 2200, "2026-10-05", SourceTicketOCR),
		pr("huevo-12", "walmart", 4000, "2026-10-10", SourceManual),
	}, "2026-10-12")
	base := Cart{Code: "b", Store: "bodega", Kind: "expense", ObservedDate: d("2026-10-10"), Source: SourceManual,
		Lines: []CartLine{{"arroz-750", 3000, 2200}, {"frijol-1kg", 1000, 3000}}}
	c, err := b.CompareExpense(base)
	if err != nil {
		t.Fatal(err)
	}
	eq(t, "renglón idéntico cotizado", c.Rows[0].Kind, "identical")
	eq(t, "empate", c.Ties, 1)
	eq(t, "Walmart sin frijol", storeRow(t, c.Stores, "walmart").MissingCount, 1)
	if _, err := b.CompareExpense(Cart{Store: "bodega", Source: SourceManual, Lines: []CartLine{{"nada", 1000, 1}}}); err == nil {
		t.Error("presentación desconocida aceptada")
	}
	// Empate de costo de paquetes: gana la de menor precio por unidad base (§7.3).
	b = newBook(t, []Store{stWalmart}, []PriceObs{
		pr("arroz-1kg", "walmart", 3000, "2026-10-10", SourceManual),
		pr("arroz-750", "walmart", 1500, "2026-10-10", SourceManual),
		pr("arroz-500", "walmart", 1500, "2026-10-10", SourceManual),
	}, "2026-10-12")
	k, _ := b.CompareBasket([]ListItem{{Product: "arroz", NeededMilli: i64(900000)}})
	eq(t, "empate de paquetes", k.Items[0].Prices[0].Presentation, "arroz-750")
	k, _ = b.CompareBasket([]ListItem{{Product: "arroz", NeededMilli: i64(400000)}})
	eq(t, "empate total por catálogo", k.Items[0].Prices[0].Presentation, "arroz-750")
	b = newBook(t, []Store{stWalmart}, []PriceObs{
		pr("arroz-1kg", "walmart", 2000, "2026-10-10", SourceManual),
		pr("arroz-500", "walmart", 1000, "2026-10-10", SourceManual),
	}, "2026-10-12")
	eq(t, "orden de catálogo (paso 4)", quote(t, b, "arroz", "walmart", "").Presentation, "arroz-500")
}

func TestResumenFrescura(t *testing.T) {
	b := newBook(t, []Store{stWalmart, stBodega}, []PriceObs{
		pr("arroz-1kg", "walmart", 2800, "2026-10-10", SourceManual),
		pr("frijol-1kg", "bodega", 3000, "2026-09-01", SourceManual),
		pr("huevo-12", "bodega", 3000, "2026-10-01", SourceProfeco),
	}, "2026-10-12")
	ps, sl := b.Summary()
	got := map[string]string{}
	for _, p := range ps {
		got[p.Product] = p.Status
	}
	eq(t, "arroz", got["arroz"], "fresh")
	eq(t, "frijol", got["frijol"], "stale")
	eq(t, "huevo (solo PROFECO)", got["huevo"], "no_price")
	eq(t, "Walmart", *sl[0].LatestDate, "2026-10-10")
	eq(t, "Bodega", *sl[1].LatestDate, "2026-09-01")
}

func TestCatalogoYSemilla(t *testing.T) {
	g := []Group{{"a", "A", "chain"}}
	for _, c := range []struct {
		p  []Product
		pr []Presentation
	}{
		{[]Product{{"x", "X", "g", "a"}, {"x", "X", "g", "a"}}, nil},
		{[]Product{{"x", "X", "oz", "a"}}, nil},
		{[]Product{{"x", "X", "g", "a"}}, []Presentation{{"p", "x", "P", 1}, {"p", "x", "P", 1}}},
		{[]Product{{"x", "X", "g", "a"}}, []Presentation{{"p", "y", "P", 1}}},
		{[]Product{{"x", "X", "g", "a"}}, []Presentation{{"p", "x", "P", 0}}},
	} {
		if _, err := NewCatalog(g, c.p, c.pr); err == nil {
			t.Errorf("catálogo inválido aceptado: %+v", c)
		}
	}
	for _, s := range []string{"", "-1", "1.234", "abc", "+3"} {
		if _, err := ParseCents(s); err == nil {
			t.Errorf("ParseCents(%q) aceptado", s)
		}
	}
	for s, want := range map[string]int64{"3094.5": 309450, "21.0": 2100, "150": 15000, "0.05": 5} {
		if v, _ := ParseCents(s); v != want {
			t.Errorf("ParseCents(%q) = %d", s, v)
		}
	}
	if _, err := NewBook(testCatalog(t), []Store{stWalmart, stWalmart}, nil, DefaultSettings(), d("2026-10-10")); err == nil {
		t.Error("tienda repetida aceptada")
	}
	if _, err := NewBook(testCatalog(t), nil, nil, Settings{}, d("2026-10-10")); err == nil {
		t.Error("parámetros vacíos aceptados")
	}
	if _, err := NewBook(testCatalog(t), []Store{stWalmart}, []PriceObs{pr("nada", "walmart", 1, "2026-10-10", SourceManual)}, DefaultSettings(), d("2026-10-10")); err == nil {
		t.Error("precio inválido aceptado")
	}
	eq(t, "NewFrac(0, 0)", NewFrac(0, 0), Frac{0, 1})
	eq(t, "NewFrac(-2, 4)", NewFrac(-2, 4), Frac{-1, 2})
	cat := testCatalog(t)
	if _, err := LoadReport(repo+"/datos/semilla/data.json", cat); err == nil {
		t.Error("data.json con un catálogo sin sus presentaciones aceptado")
	}
	if _, err := LoadCatalogFile(repo + "/docs/plan/protocolo-agentes.md"); err == nil {
		t.Error("texto leído como catálogo")
	}
	if _, err := LoadReport(repo+"/datos/semilla/nada.json", cat); err == nil {
		t.Error("archivo inexistente aceptado")
	}
	if _, err := LoadReport(repo+"/docs/plan/protocolo-agentes.md", cat); err == nil {
		t.Error("texto leído como reporte")
	}
}

// Bloqueo S5-carritos-detalle, opción A (PROVISIONAL): el ganador de un sustituto con varias
// presentaciones en una tienda usa importe total ÷ contenido total; missing_count cuenta los
// productos que faltan; stale_count cuenta cada presentación agrupada una sola vez.
func TestEvalCarritosVariasPresentaciones(t *testing.T) {
	in := evalInput(`, "carts": [
 {"code": "b", "store": "bodega", "kind": "expense", "observed_date": "2026-09-01", "source": "manual", "lines": [
  {"presentation": "huevo-12", "quantity_milli": 2000, "unit_price_cents": 4000}, {"presentation": "arroz-750", "quantity_milli": 3000, "unit_price_cents": 2200},
  {"presentation": "frijol-1kg", "quantity_milli": 1000, "unit_price_cents": 3000}]},
 {"code": "w", "store": "walmart", "kind": "cart", "observed_date": "2026-10-10", "source": "manual", "lines": [
  {"presentation": "huevo-12", "quantity_milli": 1000, "unit_price_cents": 4100}, {"presentation": "arroz-1kg", "quantity_milli": 3000, "unit_price_cents": 2800}]},
 {"code": "s", "store": "sams", "kind": "cart", "observed_date": "2026-10-10", "source": "manual", "lines": [
  {"presentation": "huevo-12", "quantity_milli": 1000, "unit_price_cents": 3900}, {"presentation": "arroz-750", "quantity_milli": 1000, "unit_price_cents": 1500},
  {"presentation": "arroz-1kg", "quantity_milli": 1000, "unit_price_cents": 5000}]}]`)
	m := run(t, in, "compare_carts", `{"base": "b", "others": ["w", "s"]}`)
	got := map[string][3]any{}
	var order []string
	for _, x := range m["stores"].([]any) {
		r := x.(map[string]any)
		order = append(order, r["store"].(string))
		got[r["store"].(string)] = [3]any{r["total_cents"], r["missing_count"], r["stale_count"]}
	}
	eq(t, "orden por total", strings.Join(order, ","), "sams,walmart,bodega")
	eq(t, "bodega total/faltantes/viejos", js(got["bodega"]), `[17600,0,3]`)
	eq(t, "walmart total/faltantes/viejos", js(got["walmart"]), `[12500,1,0]`)
	eq(t, "sams total/faltantes/viejos", js(got["sams"]), `[10400,1,0]`)
	eq(t, "winners en el orden de stores", js(m["winners"]), `[{"count":1,"store":"sams"},{"count":0,"store":"walmart"},{"count":0,"store":"bodega"}]`)
	rows := m["rows"].([]any)
	eq(t, "renglones", len(rows), 4)
	// arroz: Bodega 6600/2.25 kg = 29.33; Walmart 8400/3 kg = 28.00; Sam's 6500/1.75 kg = 37.14 (su 750 g sola daría 20.00)
	eq(t, "sustituto con dos presentaciones en Sam's", js(rows[1].(map[string]any)["winner_store"]), `"walmart"`)
	eq(t, "idéntico", js(rows[0].(map[string]any)["winner_store"]), `"sams"`)
}
