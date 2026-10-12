package api_test

// Datos de apoyo de la prueba de contrato: cuerpos que el esquema no puede inventar, filas que ninguna operación
// crea (reglas de reparto, grupos de producto, chats detectados) y operaciones con sesión propia.

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"finanzas-hogar/internal/api/sistema"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

// with genera el cuerpo del esquema y lo ajusta.
func with(op string, mut func(k *K, m map[string]any)) func(k *K) any {
	return func(k *K) any {
		meta, _ := apiMeta(op)
		k.root = meta.RequestSchema
		m, _ := k.gen(k.doc.Components.Schemas[meta.RequestSchema], "", 0).(map[string]any)
		if m == nil {
			m = map[string]any{}
		}
		mut(k, m)
		return m
	}
}

func creditPurchase(k *K, m map[string]any) {
	for kk := range m {
		delete(m, kk)
	}
	m["credit_id"] = k.ids["credit_id"]
	m["kind"] = "installment_plan"
	m["purchase_date"] = "2026-10-10"
	m["principal_cents"] = 309450
	m["interest_cents"] = 80050
	m["term_count"] = 4
	m["shares"] = []map[string]any{{"person_id": k.s.Fx.AxelPersonID, "principal_cents": 209450}, {"person_id": k.s.Fx.KarenPersonID, "principal_cents": 100000}}
}

// settingsEcho devuelve el valor vigente del primer parámetro de la sección (GET de la misma sección) para volver a mandarlo.
func settingsEcho(path string) func(k *K) any {
	return func(k *K) any {
		res := k.c.Raw("GET", path, nil)
		var out struct {
			Parameters map[string]struct {
				Value any `json:"value"`
			} `json:"parameters"`
		}
		_ = json.Unmarshal(res.Body, &out)
		keys := make([]string, 0, len(out.Parameters))
		for kk := range out.Parameters {
			keys = append(keys, kk)
		}
		sort.Strings(keys)
		if len(keys) == 0 {
			return map[string]any{"parameters": map[string]any{}}
		}
		return map[string]any{"parameters": map[string]any{keys[0]: out.Parameters[keys[0]].Value}}
	}
}

func init() {
	over["login"] = func(k *K) any {
		return map[string]any{"username": testutil.AxelUser, "password": testutil.FakeCredential}
	}
	over["createCategory"] = with("createCategory", func(k *K, m map[string]any) {
		m["name"] = fmt.Sprintf("Renta %d", k.n)
	})
	over["createNutritionFact"] = with("createNutritionFact", func(k *K, m map[string]any) {
		m["presentation_id"] = k.ids["presentation_id"]
		m["kcal_milli_per_100g"] = 250000
		m["source"] = "mx_table"
		m["source_ref"] = "Tabla de prueba"
	})
	over["createRecipe"] = with("createRecipe", func(k *K, m map[string]any) {
		m["instructions"] = "Mezclar y cocinar."
		m["meal_slot_ids"] = []int64{1}
		m["ingredients"] = []map[string]any{{"product_id": k.ids["product_id"], "presentation_id": k.ids["presentation_id"], "quantity_milli": 200000}}
	})
	over["createInventoryItem"] = with("createInventoryItem", func(k *K, m map[string]any) {
		m["product_id"] = k.ids["product_id"]
		m["tracking_mode"] = "stock"
		m["quantity_milli"] = 1000
		m["counted_date"] = "2026-10-12"
	})
	over["createCredit"] = with("createCredit", func(k *K, m map[string]any) {
		for kk := range m {
			delete(m, kk)
		}
		m["name"] = fmt.Sprintf("Kueski de prueba %d", k.n)
		m["kind"] = "revolving_line"
		m["holder_person_id"] = k.s.Fx.AxelPersonID
		m["credit_limit_cents"] = 1000000
		m["payment_frequency"] = "fortnightly"
		m["first_payment_min_days"] = 15
	})
	over["createCreditTermRate"] = with("createCreditTermRate", func(k *K, m map[string]any) {
		m["source"] = "quote"
		m["quote_date"] = "2026-10-10"
		m["term_count"] = 4
		m["rate_num"] = 80050
		m["rate_den"] = 309450
		m["rate_decimals"] = 4
	})
	over["createProduct"] = with("createProduct", func(k *K, m map[string]any) {
		m["product_group_id"] = k.ids["product_group_id"]
		m["base_unit"] = "g"
		delete(m, "presentations")
	})
	over["createCreditPurchase"] = with("createCreditPurchase", creditPurchase)
	over["previewCreditPurchase"] = with("previewCreditPurchase", creditPurchase)
	over["createWeightEntry"] = with("createWeightEntry", func(k *K, m map[string]any) { m["weight_g"] = 75000 })
	over["createBudget"] = with("createBudget", func(k *K, m map[string]any) {
		m["valid_from_date"] = "2026-10-01"
		m["category_id"] = k.ids["category_id"]
	})
	over["createMealPlanEntries"] = with("createMealPlanEntries", func(k *K, m map[string]any) {
		m["plan_dates"] = []string{"2026-10-14"}
		m["meal_slot_id"] = 1
		m["person_ids"] = []int64{k.s.Fx.AxelPersonID}
		m["recipe_id"] = k.ids["recipe_id"]
		m["is_eating_out"] = false
	})
	over["createShoppingListItem"] = with("createShoppingListItem", func(k *K, m map[string]any) {
		m["product_id"] = k.ids["product_id"]
		m["quantity_milli"] = 1000
	})
	over["generateExternalStatement"] = with("generateExternalStatement", func(k *K, m map[string]any) {
		m["person_id"] = k.s.Fx.KarenPersonID
		m["period_start_date"] = "2026-10-01"
		m["period_end_date"] = "2026-10-31"
	})
	over["restoreBackup"] = with("restoreBackup", func(k *K, m map[string]any) { m["confirmation"] = "RESTAURAR" })
	over["revokeSessions"] = func(k *K) any { return map[string]any{"scope": "all_except_current"} }
	over["adjustInventoryItem"] = func(k *K) any {
		return map[string]any{"reason": "count", "quantity_milli": 2000, "counted_date": "2026-10-12"}
	}
	over["updateAlertSettings"] = settingsEcho("/api/v1/settings/alerts")
	over["updateDataSettings"] = settingsEcho("/api/v1/settings/data")
	over["updateHealthSettings"] = settingsEcho("/api/v1/settings/health")
	over["updateSecuritySettings"] = settingsEcho("/api/v1/settings/security")
	over["updatePriceSettings"] = settingsEcho("/api/v1/settings/prices")
	over["updateDashboardSettings"] = settingsEcho("/api/v1/settings/dashboard") // C1 v1.5, K54 (correccion-especificacion-r2)
	// Llaves obligatorias from y to (AAAA-MM) del resumen mensual del tablero: el generador no saca un mes de «examples».
	qover["getDashboardMonthlyOverview"] = map[string]string{"from": "2026-10", "to": "2026-12"}
	over["generateMonthlyReport"] = func(k *K) any { return map[string]any{"month": "2026-10", "send_telegram": false} }
}

func apiMeta(op string) (m struct{ RequestSchema string }, ok bool) {
	for _, o := range operations() {
		if o.OperationID == op {
			m.RequestSchema = o.RequestSchema
			return m, true
		}
	}
	return m, false
}

// prepExtra deja listo lo que una operación necesita y que ninguna otra operación puede crear.
func prepExtra(k *K, op string) {
	db := k.s.DB
	exec := func(q string, a ...any) int64 {
		res, err := db.Exec(q, a...)
		if err != nil {
			k.t.Fatalf("preparación %s: %v", op, err)
		}
		id, _ := res.LastInsertId()
		return id
	}
	now := dominio.FormatDBTime(k.s.Clock.Now())
	switch op {
	case "confirmTelegramChatCandidate", "dismissTelegramChatCandidate":
		k.ids["candidate_id"] = exec(`INSERT INTO telegram_chat_candidate (chat_id, telegram_name, detected_at) VALUES (?, 'Prueba', ?)`, fmt.Sprint(5000+k.n), now)
	case "updateNotificationRecipient", "setNotificationRecipients", "sendTestNotification", "listNotificationRecipients":
		if k.ids["recipient_id"] == 0 {
			k.ids["recipient_id"] = exec(`INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, '1001')`, k.s.Fx.AxelPersonID)
		}
	case "updateCostShareRule", "listCostShareRules":
		if k.ids["cost_share_rule_id"] == 0 {
			k.ids["cost_share_rule_id"] = exec(`INSERT INTO cost_share_rule (name, remainder_person_id, is_template, is_default) VALUES ('Hogar de prueba', ?, 1, 0)`, k.s.Fx.AxelPersonID)
		}
	case "revokeSession":
		_ = k.s.SessionFor(k.t, testutil.XimenaUser)
		var id int64
		_ = db.QueryRow(`SELECT id FROM device_session ORDER BY id DESC LIMIT 1`).Scan(&id)
		k.ids["session_id"] = id
	case "getCaCertificate":
		f := filepath.Join(testutil.TempDir(k.t), "rootCA.pem")
		_ = os.WriteFile(f, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("certificado-falso-de-prueba")}), 0o600)
		old := sistema.CACertPath
		sistema.CACertPath = f
		k.t.Cleanup(func() { sistema.CACertPath = old })
	case "getCardStatus":
		if k.ids["card_credit_id"] == 0 {
			body := map[string]any{"name": "Tarjeta de prueba", "kind": "bank_card", "holder_person_id": k.s.Fx.AxelPersonID,
				"credit_limit_cents": 5000000, "payment_frequency": "monthly", "statement_day": 5, "payment_due_day": 25}
			res := k.do(k.c, "POST", "/api/v1/credits", body, "Idempotency-Key", k.key())
			if m := jsonOf(res.Body); m != nil {
				if id, ok := m["id"].(float64); ok {
					k.ids["card_credit_id"] = int64(id)
				}
			}
		}
	case "markExpensePaid":
		if k.ids["upcoming_expense_id"] == 0 {
			body := map[string]any{"kind": "upcoming", "category_id": k.ids["category_id"], "amount_cents": 50000, "expense_date": "2026-10-20", "payment_method": "cash"}
			res := k.do(k.c, "POST", "/api/v1/expenses", body, "Idempotency-Key", k.key())
			if m := jsonOf(res.Body); m != nil {
				if id, ok := m["id"].(float64); ok {
					k.ids["upcoming_expense_id"] = int64(id)
				}
			}
		}
	}
}

func init() {
	special["logout"] = func(k *K, st *testing.T, f *fila) {
		cl := k.s.SessionFor(st, testutil.XimenaUser)
		res := k.do(cl, "POST", "/api/v1/auth/logout", nil)
		if res.StatusCode != 204 {
			f.Fallos = append(f.Fallos, fmt.Sprintf("logout: %d", res.StatusCode))
			st.Errorf("logout: %d %s", res.StatusCode, k.cause(res))
			return
		}
		f.Verifs = append(f.Verifs, "204")
		if r2 := k.do(cl, "GET", "/api/v1/auth/me", nil); r2.StatusCode != 401 {
			f.Fallos = append(f.Fallos, fmt.Sprintf("tras cerrar la sesión /auth/me dio %d", r2.StatusCode))
			st.Errorf("tras cerrar la sesión /auth/me dio %d", r2.StatusCode)
		} else {
			f.Verifs = append(f.Verifs, "la sesión cerrada da 401")
		}
		if r3 := k.do(k.anon, "POST", "/api/v1/auth/logout", nil); r3.StatusCode != 401 {
			f.Fallos = append(f.Fallos, fmt.Sprintf("logout sin sesión dio %d", r3.StatusCode))
			st.Errorf("logout sin sesión dio %d", r3.StatusCode)
		} else {
			f.Verifs = append(f.Verifs, "401 sin sesión")
		}
	}
	special["changeOwnPassword"] = func(k *K, st *testing.T, f *fila) {
		cl := k.s.SessionFor(st, testutil.XimenaUser)
		nueva := "otra-contrasena-falsa-9"
		res := k.do(cl, "POST", "/api/v1/auth/password", map[string]any{"current_password": testutil.FakeCredential, "new_password": nueva})
		if res.StatusCode != 204 {
			f.Fallos = append(f.Fallos, fmt.Sprintf("cambio de contraseña: %d", res.StatusCode))
			st.Errorf("cambio de contraseña: %d %s", res.StatusCode, k.cause(res))
			return
		}
		f.Verifs = append(f.Verifs, "204")
		if r2 := k.do(k.anon, "POST", "/api/v1/auth/login", map[string]any{"username": testutil.XimenaUser, "password": nueva}); r2.StatusCode != 200 {
			f.Fallos = append(f.Fallos, fmt.Sprintf("login con la contraseña nueva: %d", r2.StatusCode))
			st.Errorf("login con la contraseña nueva: %d %s", r2.StatusCode, k.cause(r2))
		} else {
			f.Verifs = append(f.Verifs, "la contraseña nueva entra")
		}
		if r3 := k.do(cl, "POST", "/api/v1/auth/password", map[string]any{"current_password": "no-es-la-actual", "new_password": nueva + "x"}); r3.StatusCode != 422 && r3.StatusCode != 401 {
			f.Fallos = append(f.Fallos, fmt.Sprintf("contraseña actual incorrecta: %d", r3.StatusCode))
			st.Errorf("contraseña actual incorrecta: %d %s", r3.StatusCode, k.cause(r3))
		} else {
			f.Verifs = append(f.Verifs, fmt.Sprintf("%d con la contraseña actual incorrecta", r3.StatusCode))
		}
	}
}
