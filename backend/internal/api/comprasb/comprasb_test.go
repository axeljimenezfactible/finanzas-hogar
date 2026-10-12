package comprasb_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comprasb"
	"finanzas-hogar/internal/testutil"
)

// Credencial de servicio falsa: solo existe en estas pruebas (protocolo §7).
const fakeServiceToken = "credencial-de-servicio-falsa-de-prueba"

type fakeService struct{}

func (fakeService) VerifyServiceToken(_ context.Context, token string) (bool, error) {
	return token == fakeServiceToken, nil
}

var keyN int

// key devuelve una Idempotency-Key nueva (16 a 128 caracteres).
func key() string { keyN++; return fmt.Sprintf("clave-prueba-compras-b-%06d", keyN) }

// setup levanta el servidor con un catálogo mínimo: cadena con sitio y su sucursal, dos productos con una
// presentación cada uno, precios manual (fresco), ticket (viejo) y PROFECO, y dos renglones de inventario.
func setup(t *testing.T) (*testutil.Server, *testutil.Client) {
	t.Helper()
	s := testutil.NewServer(t, comprasb.Register)
	s.Reg.Deps().Service = fakeService{}
	for _, q := range []string{
		`INSERT INTO product_group (id, name, purchase_channel) VALUES (1, 'Abarrotes', 'chain')`,
		`INSERT INTO store (id, kind, name, profeco_name, website_url) VALUES (1, 'chain', 'Bodega Aurrera', 'Bodega Aurrera', 'https://tienda.invalid')`,
		`INSERT INTO store (id, kind, parent_id, name, profeco_name, municipality) VALUES (2, 'branch', 1, 'Bodega Aurrera Centro', 'Bodega Aurrera Centro', 'Cuauhtémoc')`,
		`INSERT INTO product (id, name, product_group_id, base_unit) VALUES (1, 'Aceite', 1, 'ml'), (2, 'Atún', 1, 'g'), (3, 'Arroz', 1, 'g')`,
		`INSERT INTO presentation (id, product_id, name, size_milli) VALUES (1, 1, 'Nutrioli 500 ml', 500000), (2, 2, 'Lata 140 g', 140000)`,
		`INSERT INTO price (id, presentation_id, store_id, price_cents, observed_date, source, source_ref) VALUES
		 (1, 1, 1, 15000, '2026-10-01', 'manual', NULL), (2, 1, 1, 14000, '2026-08-01', 'ticket_manual', NULL), (3, 1, 2, 13000, '2026-09-01', 'profeco', 'QQP 2026')`,
		`INSERT INTO inventory_item (id, product_id, tracking_mode, quantity_milli) VALUES (1, 1, 'stock', 0)`,
		`INSERT INTO inventory_item (id, product_id, tracking_mode, quantity_milli, counted_date, rate_quantity_milli, rate_period_days, reorder_level_milli)
		 VALUES (2, 2, 'consumption_rate', 1000000, '2026-10-05', 140000, 7, 900000)`,
	} {
		if _, err := s.DB.Exec(q); err != nil {
			t.Fatalf("%v\n%s", err, q)
		}
	}
	return s, s.SessionFor(t, testutil.AxelUser)
}

func want(t *testing.T, res *testutil.Response, code int) {
	t.Helper()
	if res.StatusCode != code {
		t.Fatalf("status %d, quería %d: %s", res.StatusCode, code, res.Body)
	}
}

func agent(s *testutil.Server, t *testing.T, method, path string, body any) *testutil.Response {
	return s.Anonymous(t).Do(method, path, body, "Authorization", "Bearer "+fakeServiceToken, "Idempotency-Key", key())
}

func TestOp_listProductPrices(t *testing.T) {
	_, c := setup(t)
	res := c.Do("GET", "/api/v1/products/1/prices", nil)
	want(t, res, 200)
	var pg api.PricePage
	res.JSON(t, &pg)
	if len(pg.Items) != 3 || pg.Items[0].ID != 1 || *pg.Items[0].BaseUnitPriceCents != 30000 || *pg.Items[0].IsStale || *pg.Items[0].AgeDays != 11 {
		t.Fatalf("%+v", pg.Items)
	}
	if !*pg.Items[2].IsStale { // 2026-08-01: 72 días ≥ 30 (K22)
		t.Fatalf("el precio de agosto debe ser viejo: %+v", pg.Items[2])
	}
	res = c.Do("GET", "/api/v1/products/1/prices?source=profeco&store_id=2&limit=1", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].Source != "profeco" {
		t.Fatalf("%+v", pg.Items)
	}
	want(t, c.Do("GET", "/api/v1/products/1/prices?source=sitio", nil), 400)
	want(t, c.Do("GET", "/api/v1/products/99/prices", nil), 404)
}

func TestOp_getPriceSummary(t *testing.T) {
	s, c := setup(t)
	res := c.Do("GET", "/api/v1/prices/summary", nil)
	want(t, res, 200)
	var sm api.PriceSummary
	res.JSON(t, &sm)
	if sm.FreshCount != 1 || sm.StaleCount != 0 || sm.NoPriceCount != 2 || len(sm.Stores) != 2 || *sm.Stores[0].LatestObservedDate != "2026-10-01" || sm.Stores[1].LatestObservedDate != nil {
		t.Fatalf("%+v", sm)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/prices/summary", nil), 401)
}

func TestOp_createPrice(t *testing.T) {
	_, c := setup(t)
	body := map[string]any{"presentation_id": 2, "store_id": 1, "price_cents": 2100, "observed_date": "2026-10-12", "is_promotional": true}
	res := c.Do("POST", "/api/v1/prices", body, "Idempotency-Key", key())
	want(t, res, 201)
	var p api.Price
	res.JSON(t, &p)
	if p.Source != "manual" || !p.IsPromotional || *p.BaseUnitPriceCents != 15000 {
		t.Fatalf("%+v", p)
	}
	want(t, c.Do("POST", "/api/v1/prices", body, "Idempotency-Key", key()), 409)
	want(t, c.Do("POST", "/api/v1/prices", map[string]any{"presentation_id": 99, "store_id": 1, "price_cents": 1, "observed_date": "2026-10-13"}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/prices", map[string]any{"presentation_id": 1, "store_id": 1, "price_cents": "973.75", "observed_date": "2026-10-01"}, "Idempotency-Key", key()), 400)
	want(t, c.Do("POST", "/api/v1/prices", body), 400)
}

func TestOp_updatePrice(t *testing.T) {
	_, c := setup(t)
	res := c.Do("PATCH", "/api/v1/prices/1", map[string]any{"price_cents": 15500})
	want(t, res, 200)
	var p api.Price
	res.JSON(t, &p)
	if p.PriceCents != 15500 {
		t.Fatalf("%+v", p)
	}
	want(t, c.Do("PATCH", "/api/v1/prices/3", map[string]any{"price_cents": 1}), 403)
	want(t, c.Do("PATCH", "/api/v1/prices/1", map[string]any{"observed_date": "2026-12-01"}), 422)
	want(t, c.Do("PATCH", "/api/v1/prices/2", map[string]any{"observed_date": "2026-10-01", "price_cents": 15000}), 200)
	want(t, c.Do("PATCH", "/api/v1/prices/99", map[string]any{"price_cents": 1}), 404)
}

func TestOp_deletePrice(t *testing.T) {
	s, c := setup(t)
	want(t, c.Do("DELETE", "/api/v1/prices/2", nil), 204)
	want(t, c.Do("DELETE", "/api/v1/prices/2", nil), 404)
	want(t, c.Do("DELETE", "/api/v1/prices/3", nil), 403)
	var n int
	if err := s.DB.QueryRow(`SELECT count(*) FROM audit_log WHERE table_name = 'price' AND action = 'delete' AND row_id = 2`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("audit_log %d %v", n, err)
	}
}

func TestOp_getPriceSources(t *testing.T) {
	s, c := setup(t)
	res := c.Do("GET", "/api/v1/price-sources", nil)
	want(t, res, 200)
	var ps api.PriceSources
	res.JSON(t, &ps)
	if !ps.ManualEnabled || ps.AssistedReadingEnabled || len(ps.Sites) != 1 || ps.Sites[0].Verdict != "missing" || ps.ProfecoPriceCount != 1 ||
		len(ps.ProfecoCoverage) != 1 || *ps.ProfecoCoverage[0].StoreID != 1 || *ps.ProfecoCoverage[0].Municipality != "Cuauhtémoc" {
		t.Fatalf("%+v", ps)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/price-sources", nil), 401)
}

func TestOp_getPriceSettings(t *testing.T) {
	s, c := setup(t)
	res := c.Do("GET", "/api/v1/settings/prices", nil)
	want(t, res, 200)
	var ps api.PriceSettings
	res.JSON(t, &ps)
	if ps.AssistedReadingEnabled || ps.Parameters.PriceStaleAfterDays.Value.Value != 30 || ps.Parameters.PriceZeroWeightAfterDays.Value.Value != 90 ||
		!ps.Parameters.PriceProfecoImportMaxGB.IsInitial {
		t.Fatalf("%+v", ps)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/settings/prices", nil), 401)
}

func TestOp_updatePriceSettings(t *testing.T) {
	_, c := setup(t)
	res := c.Do("PATCH", "/api/v1/settings/prices", map[string]any{"parameters": map[string]any{"price.stale_after_days": 20}, "assisted_reading_enabled": false})
	want(t, res, 200)
	var ps api.PriceSettings
	res.JSON(t, &ps)
	if ps.Parameters.PriceStaleAfterDays.Value.Value != 20 || ps.Parameters.PriceStaleAfterDays.IsInitial {
		t.Fatalf("%+v", ps)
	}
	want(t, c.Do("PATCH", "/api/v1/settings/prices", map[string]any{"assisted_reading_enabled": true}), 409)
	want(t, c.Do("PATCH", "/api/v1/settings/prices", map[string]any{"parameters": map[string]any{"price.zero_weight_after_days": 10}}), 422)
	want(t, c.Do("PATCH", "/api/v1/settings/prices", map[string]any{"parameters": map[string]any{"price.stale_after_days": nil}}), 422)
	want(t, c.Do("PATCH", "/api/v1/settings/prices", map[string]any{"reset_parameters": []string{"price.stale_after_days"}}), 200)
	want(t, c.Do("PATCH", "/api/v1/settings/prices", map[string]any{}), 400)
}

func multipartFile(t *testing.T, ctype, name string, data []byte) ([]byte, string) {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, name))
	h.Set("Content-Type", ctype)
	pw, err := mw.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}
	pw.Write(data)
	mw.Close()
	return b.Bytes(), mw.FormDataContentType()
}

func TestOp_importProfeco(t *testing.T) {
	s, c := setup(t)
	s.Reg.Deps().Config.DataDir = testutil.TempDir(t)
	s.Reg.Deps().Config.TmpDir = testutil.TempDir(t)
	csv := "producto,presentacion,marca,categoria,catalogo,precio,fecha_registro,cadena_comercial,giro,nombre_comercial,direccion,estado,municipio,latitud,longitud\n" +
		"Aceite,Botella 500 ml,Nutrioli,Aceites,Basicos,31.5,2026/10/01,Bodega Aurrera,Supermercado,Bodega Aurrera Centro,Calle 1,Ciudad de México,Cuauhtémoc,19.43,-99.13\n"
	body, ct := multipartFile(t, "text/csv", "qqp.csv", []byte(csv))
	res := c.Do("POST", "/api/v1/profeco-imports", body, "Content-Type", ct, "Idempotency-Key", key())
	want(t, res, 201)
	var out api.ProfecoImportResult
	res.JSON(t, &out)
	if out.LoadedCount != 0 || len(out.Coverage) != 1 { // sin vinculos.csv ningún registro se liga
		t.Fatalf("%+v", out)
	}
	rar, ct2 := multipartFile(t, "application/vnd.rar", "qqp.rar", []byte("Rar!"))
	want(t, c.Do("POST", "/api/v1/profeco-imports", rar, "Content-Type", ct2, "Idempotency-Key", key()), 415)
	bad, ct3 := multipartFile(t, "text/csv", "malo.csv", []byte("a,b\n1,2\n"))
	want(t, c.Do("POST", "/api/v1/profeco-imports", bad, "Content-Type", ct3, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/profeco-imports", body, "Content-Type", ct), 400)
	// Con una liga confirmada en <DataDir>/profeco/vinculos.csv el registro se carga como precio PROFECO.
	dir := filepath.Join(s.Reg.Deps().Config.DataDir, "profeco")
	links := "producto,presentacion,marca,presentacion_id,estado,confirmo,fecha_confirmacion\nAceite,Botella 500 ml,Nutrioli,1,confirmado,axel,2026-10-11\n"
	if err := os.MkdirAll(dir, 0o700); err != nil || os.WriteFile(filepath.Join(dir, "vinculos.csv"), []byte(links), 0o600) != nil {
		t.Fatal(err)
	}
	res = c.Do("POST", "/api/v1/profeco-imports", body, "Content-Type", ct, "Idempotency-Key", key())
	want(t, res, 201)
	res.JSON(t, &out)
	if out.LoadedCount != 1 || out.SkippedCount != 0 || out.Coverage[0].PriceCount != 2 {
		t.Fatalf("con liga: %s", res.Body)
	}
}

func TestOp_listNutritionFacts(t *testing.T) {
	s, c := setup(t)
	if _, err := s.DB.Exec(`INSERT INTO nutrition_fact (product_id, kcal_milli_per_100g, source, source_ref, is_preferred) VALUES (1, 884000, 'usda_fdc', '171413', 1)`); err != nil {
		t.Fatal(err)
	}
	res := c.Do("GET", "/api/v1/products/1/nutrition-facts", nil)
	want(t, res, 200)
	var pg api.NutritionFactPage
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].KcalMilliPer100g != 884000 || !pg.Items[0].IsPreferred {
		t.Fatalf("%+v", pg)
	}
	want(t, c.Do("GET", "/api/v1/products/99/nutrition-facts", nil), 404)
}

func TestOp_createNutritionFact(t *testing.T) {
	_, c := setup(t)
	body := map[string]any{"kcal_milli_per_100g": 884000, "source": "store_label", "source_ref": "etiqueta Nutrioli", "presentation_id": 1, "is_preferred": true}
	res := c.Do("POST", "/api/v1/products/1/nutrition-facts", body, "Idempotency-Key", key())
	want(t, res, 201)
	var n api.NutritionFact
	res.JSON(t, &n)
	if n.ProductID != 1 || *n.PresentationID != 1 || !n.IsPreferred {
		t.Fatalf("%+v", n)
	}
	body["source_ref"] = "otra etiqueta"
	want(t, c.Do("POST", "/api/v1/products/1/nutrition-facts", body, "Idempotency-Key", key()), 409)
	want(t, c.Do("POST", "/api/v1/products/1/nutrition-facts", map[string]any{"kcal_milli_per_100g": 1, "source": "usda_fdc", "source_ref": "1"}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/products/1/nutrition-facts", map[string]any{"kcal_milli_per_100g": 1, "source": "mx_table"}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/products/2/nutrition-facts", map[string]any{"kcal_milli_per_100g": 1, "source": "manual", "presentation_id": 1}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/products/99/nutrition-facts", map[string]any{"kcal_milli_per_100g": 1, "source": "manual"}, "Idempotency-Key", key()), 404)
}

func TestOp_getRecipeSettings(t *testing.T) {
	s, c := setup(t)
	if _, err := s.DB.Exec(`INSERT INTO nutrition_fact (product_id, kcal_milli_per_100g, source, source_ref) VALUES (1, 884000, 'usda_fdc', '171413'), (1, 880000, 'usda_fdc', '171414')`); err != nil {
		t.Fatal(err)
	}
	res := c.Do("GET", "/api/v1/settings/recipes", nil)
	want(t, res, 200)
	var rs api.RecipeSettings
	res.JSON(t, &rs)
	if rs.DefaultBaseServings != 2 || len(rs.MealSlots) != 4 || len(rs.RecipeCounts) != 2 || len(rs.NutritionSources) != 1 || rs.NutritionSources[0].FoodCount != 1 {
		t.Fatalf("%+v", rs)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/settings/recipes", nil), 401)
}

func TestOp_reloadNutritionBase(t *testing.T) {
	_, c := setup(t)
	res := c.Do("POST", "/api/v1/nutrition-facts/reload", map[string]any{"source": "usda_fdc"}, "Idempotency-Key", key())
	want(t, res, 409) // PROVISIONAL: bloqueo S6-compras-b-nutricion
	want(t, c.Do("POST", "/api/v1/nutrition-facts/reload", map[string]any{"source": "mx_table"}, "Idempotency-Key", key()), 422)
}

func TestOp_listInventoryItems(t *testing.T) {
	_, c := setup(t)
	res := c.Do("GET", "/api/v1/inventory-items", nil)
	want(t, res, 200)
	var pg api.InventoryItemPage
	res.JSON(t, &pg)
	// Aceite: stock 0 → agotado. Atún: 1 000 000 − 140 000 × 7 / 7 = 860 000 ≤ 900 000 → reponer (motor S4).
	if len(pg.Items) != 2 || *pg.Items[0].Status != "out" || *pg.Items[1].Status != "reorder" || *pg.Items[1].EstimatedQuantityMilli != 860000 ||
		*pg.Items[1].WeeklyConsumptionMilli != 140000 || *pg.Items[1].DaysLeft != 43 {
		t.Fatalf("%s", res.Body)
	}
	res = c.Do("GET", "/api/v1/inventory-items?status=reorder&tracking_mode=consumption_rate&q=at&limit=1", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].ID != 2 || pg.NextCursor.Valid {
		t.Fatalf("%s", res.Body)
	}
	res = c.Do("GET", "/api/v1/inventory-items?limit=1&sort=-status", nil)
	res.JSON(t, &pg)
	if pg.Items[0].ID != 2 || !pg.NextCursor.Valid {
		t.Fatalf("%s", res.Body)
	}
	want(t, c.Do("GET", "/api/v1/inventory-items?status=lleno", nil), 400)
}

func TestOp_createInventoryItem(t *testing.T) {
	s, c := setup(t)
	res := c.Do("POST", "/api/v1/inventory-items", map[string]any{"product_id": 3, "tracking_mode": "stock", "quantity_milli": 2000000, "reorder_level_milli": 500000}, "Idempotency-Key", key())
	want(t, res, 201)
	var it api.InventoryItem
	res.JSON(t, &it)
	var kind, reason string
	var delta int64
	if err := s.DB.QueryRow(`SELECT kind, reason, quantity_delta_milli FROM inventory_movement WHERE inventory_item_id = ?`, it.ID).Scan(&kind, &reason, &delta); err != nil ||
		kind != "adjustment" || reason != "count" || delta != 2000000 || *it.Status != "sufficient" {
		t.Fatalf("%s %s %d %v %+v", kind, reason, delta, err, it)
	}
	want(t, c.Do("POST", "/api/v1/inventory-items", map[string]any{"product_id": 3, "tracking_mode": "stock"}, "Idempotency-Key", key()), 409)
	want(t, c.Do("POST", "/api/v1/inventory-items", map[string]any{"product_id": 3, "recipe_id": 1, "tracking_mode": "stock"}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/inventory-items", map[string]any{"product_id": 99, "tracking_mode": "consumption_rate"}, "Idempotency-Key", key()), 422)
}

func TestOp_getInventorySummary(t *testing.T) {
	s, c := setup(t)
	res := c.Do("GET", "/api/v1/inventory/summary", nil)
	want(t, res, 200)
	var sm api.InventorySummary
	res.JSON(t, &sm)
	if sm.OutCount != 1 || sm.ReorderCount != 1 || sm.SufficientCount != 0 {
		t.Fatalf("%+v", sm)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/inventory/summary", nil), 401)
}

func TestOp_listInventoryMovements(t *testing.T) {
	_, c := setup(t)
	want(t, c.Do("POST", "/api/v1/inventory-items/1/adjustments", map[string]any{"reason": "count", "quantity_milli": 1500000, "movement_date": "2026-10-10"}, "Idempotency-Key", key()), 200)
	want(t, c.Do("POST", "/api/v1/inventory-items/1/adjustments", map[string]any{"reason": "waste", "quantity_delta_milli": -500000}, "Idempotency-Key", key()), 200)
	res := c.Do("GET", "/api/v1/inventory-items/1/movements", nil)
	want(t, res, 200)
	var pg api.InventoryMovementPage
	res.JSON(t, &pg)
	if len(pg.Items) != 2 || pg.Items[0].MovementDate != "2026-10-12" || *pg.Items[0].Reason != "waste" || pg.Items[0].QuantityDeltaMilli != -500000 {
		t.Fatalf("%s", res.Body)
	}
	want(t, c.Do("GET", "/api/v1/inventory-items/99/movements", nil), 404)
}

func TestOp_adjustInventoryItem(t *testing.T) {
	_, c := setup(t)
	res := c.Do("POST", "/api/v1/inventory-items/2/adjustments", map[string]any{"reason": "count", "quantity_milli": 700000, "notes": "conteo"}, "Idempotency-Key", key())
	want(t, res, 200)
	var it api.InventoryItem
	res.JSON(t, &it)
	if it.QuantityMilli != 700000 || *it.CountedDate != "2026-10-12" || *it.EstimatedQuantityMilli != 700000 || *it.Status != "reorder" {
		t.Fatalf("%+v", it)
	}
	want(t, c.Do("POST", "/api/v1/inventory-items/2/adjustments", map[string]any{"reason": "count", "quantity_milli": 1, "quantity_delta_milli": 1}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/inventory-items/1/adjustments", map[string]any{"reason": "waste", "quantity_delta_milli": -1}, "Idempotency-Key", key()), 422)
	want(t, c.Do("POST", "/api/v1/inventory-items/99/adjustments", map[string]any{"reason": "count", "quantity_milli": 1}, "Idempotency-Key", key()), 404)
}

func TestOp_updateInventoryItem(t *testing.T) {
	_, c := setup(t)
	res := c.Do("PATCH", "/api/v1/inventory-items/2", map[string]any{"reorder_level_milli": nil, "notes": "lata chica", "expires_date": "2027-01-31"})
	want(t, res, 200)
	var it api.InventoryItem
	res.JSON(t, &it)
	if it.ReorderLevelMilli != nil || *it.Notes != "lata chica" || *it.Status != "sufficient" {
		t.Fatalf("%+v", it)
	}
	want(t, c.Do("PATCH", "/api/v1/inventory-items/1", map[string]any{"tracking_mode": "consumption_rate"}), 422)
	want(t, c.Do("PATCH", "/api/v1/inventory-items/1", map[string]any{}), 400)
	want(t, c.Do("PATCH", "/api/v1/inventory-items/99", map[string]any{"notes": "x"}), 404)
}

func TestOp_getExternalAgentReference(t *testing.T) {
	s, _ := setup(t)
	res := s.Anonymous(t).Do("GET", "/api/v1/external-agent/reference", nil, "Authorization", "Bearer "+fakeServiceToken)
	want(t, res, 200)
	var ref api.ExternalAgentReference
	res.JSON(t, &ref)
	if len(ref.Stores) != 2 || *ref.Stores[1].ParentID != 1 || len(ref.ProductGroups) != 1 || len(ref.MealSlots) != 4 {
		t.Fatalf("%+v", ref)
	}
	want(t, s.Anonymous(t).Do("GET", "/api/v1/external-agent/reference", nil, "Authorization", "Bearer otra-credencial-falsa"), 401)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/external-agent/reference", nil), 401)
}

func TestOp_importExternalCatalog(t *testing.T) {
	s, _ := setup(t)
	lote := map[string]any{"items": []any{
		map[string]any{"name": "Frijol", "base_unit": "g", "product_group_id": 1, "presentations": []any{map[string]any{"name": "Bolsa 1 kg", "size_milli": 1000000, "barcode": "7501000000001"}}},
		map[string]any{"name": "aceite", "base_unit": "ml", "product_group_id": 1, "presentations": []any{map[string]any{"name": "Nutrioli 500 ml", "size_milli": 500000}}},
		map[string]any{"name": "Atún", "base_unit": "piece", "presentations": []any{map[string]any{"name": "Lata", "size_milli": 1000}}},
		map[string]any{"name": "Sal", "base_unit": "g", "presentations": []any{map[string]any{"name": "Bolsa", "size_milli": 1000000, "barcode": "7501000000001"}}},
	}}
	res := agent(s, t, "POST", "/api/v1/external-agent/catalog", lote)
	want(t, res, 200)
	var out api.ExternalImportResult
	res.JSON(t, &out)
	if out.CreatedCount != 1 || out.UnchangedCount != 1 || out.RejectedCount != 2 || out.Items[0].Status != "created" || len(out.Items[0].PresentationIDs) != 1 ||
		*out.Items[1].ID != 1 || out.Items[2].Error.Field != "base_unit" || out.Items[3].Error.Field != "presentations[0].barcode" {
		t.Fatalf("%s", res.Body)
	}
	res = agent(s, t, "POST", "/api/v1/external-agent/catalog", lote)
	res.JSON(t, &out)
	if out.CreatedCount != 0 || out.UnchangedCount != 2 {
		t.Fatalf("repetir el lote no debe duplicar: %s", res.Body)
	}
	want(t, agent(s, t, "POST", "/api/v1/external-agent/catalog", map[string]any{"items": []any{}}), 422)
	want(t, s.Anonymous(t).Do("POST", "/api/v1/external-agent/catalog", lote, "Idempotency-Key", key()), 401)
}

func TestOp_importExternalPrices(t *testing.T) {
	s, _ := setup(t)
	item := map[string]any{"presentation_id": 1, "store_id": 1, "price_cents": 14900, "observed_date": "2026-10-10", "external_ref": "https://tienda.invalid/p/1"}
	lote := map[string]any{"items": []any{item,
		map[string]any{"presentation_id": 1, "store_id": 99, "price_cents": 1, "observed_date": "2026-10-10"},
		map[string]any{"presentation_id": 1, "store_id": 1, "price_cents": 1, "observed_date": "2026-10-13"}}}
	res := agent(s, t, "POST", "/api/v1/external-agent/prices", lote)
	want(t, res, 200)
	var out api.ExternalImportResult
	res.JSON(t, &out)
	if out.CreatedCount != 1 || out.RejectedCount != 2 || out.Items[1].Error.Field != "store_id" || out.Items[2].Error.Field != "observed_date" {
		t.Fatalf("%s", res.Body)
	}
	var src string
	if err := s.DB.QueryRow(`SELECT source FROM price WHERE id = ?`, *out.Items[0].ID).Scan(&src); err != nil || src != "site_assisted" {
		t.Fatalf("%s %v", src, err)
	}
	res = agent(s, t, "POST", "/api/v1/external-agent/prices", map[string]any{"items": []any{item}})
	res.JSON(t, &out)
	if out.UnchangedCount != 1 {
		t.Fatalf("%s", res.Body)
	}
	item["price_cents"] = 15100
	res = agent(s, t, "POST", "/api/v1/external-agent/prices", map[string]any{"items": []any{item}})
	res.JSON(t, &out)
	if out.UpdatedCount != 1 {
		t.Fatalf("%s", res.Body)
	}
	want(t, agent(s, t, "POST", "/api/v1/external-agent/prices", map[string]any{"items": []any{}}), 422)
}

func TestOp_importExternalRecipes(t *testing.T) {
	s, _ := setup(t)
	rec := map[string]any{"name": "Atún con arroz", "meal_slot_ids": []int{2}, "ingredients": []any{
		map[string]any{"product_id": 2, "presentation_id": 2, "quantity_milli": 140000}, map[string]any{"product_id": 3, "quantity_milli": 200000}}}
	lote := map[string]any{"items": []any{rec,
		map[string]any{"name": "Receta rara", "ingredients": []any{map[string]any{"product_id": 99, "quantity_milli": 1}}}}}
	res := agent(s, t, "POST", "/api/v1/external-agent/recipes", lote)
	want(t, res, 200)
	var out api.ExternalImportResult
	res.JSON(t, &out)
	if out.CreatedCount != 1 || out.RejectedCount != 1 || out.Items[1].Error.Field != "ingredients[0].product_id" {
		t.Fatalf("%s", res.Body)
	}
	var active int
	var origin string
	if err := s.DB.QueryRow(`SELECT is_active, origin FROM recipe WHERE id = ?`, *out.Items[0].ID).Scan(&active, &origin); err != nil || active != 0 || origin != "household" {
		t.Fatalf("borrador: %d %s %v", active, origin, err)
	}
	res = agent(s, t, "POST", "/api/v1/external-agent/recipes", map[string]any{"items": []any{rec}})
	res.JSON(t, &out)
	if out.RejectedCount != 1 || out.Items[0].Error.Code != "conflict" {
		t.Fatalf("%s", res.Body)
	}
	many := make([]any, 101)
	for i := range many {
		many[i] = map[string]any{"name": fmt.Sprint("R", i)}
	}
	want(t, agent(s, t, "POST", "/api/v1/external-agent/recipes", map[string]any{"items": many}), 422)
}
