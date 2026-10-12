package avisos_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/avisos"
	"finanzas-hogar/internal/telegram"
	"finanzas-hogar/internal/testutil"
)

var fakeToken = strings.Repeat("7", 10) + ":" + strings.Repeat("falso_", 6)

type entorno struct {
	s       *testutil.Server
	c       *testutil.Client
	bodies  []string
	mu      sync.Mutex
	axelRec int64
}

// nuevo arma el servidor con el cliente real de Telegram apuntando a un servidor falso y un openclaw.json falso.
func nuevo(t *testing.T, user string) *entorno {
	t.Helper()
	e := &entorno{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		e.mu.Lock()
		e.bodies = append(e.bodies, string(b))
		e.mu.Unlock()
		w.Write([]byte(`{"ok":true,"result":{}}`))
	}))
	t.Cleanup(ts.Close)
	cfg := filepath.Join(testutil.TempDir(t), "openclaw.json")
	os.WriteFile(cfg, []byte(`{"bots":{"multix":{"botToken":"`+fakeToken+`"}}}`), 0o600)
	cl := &telegram.Client{Base: ts.URL, Token: telegram.FileToken(cfg), HTTP: ts.Client()}
	e.s = testutil.NewServer(t, avisos.RegisterWith(avisos.Options{Enviador: cl}))
	e.c = e.s.SessionFor(t, user)
	res, err := e.s.DB.Exec(`INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, '1001')`, e.s.Fx.AxelPersonID)
	if err != nil {
		t.Fatal(err)
	}
	e.axelRec, _ = res.LastInsertId()
	return e
}

func (e *entorno) notif(t *testing.T, code string) string {
	var id int64
	if err := e.s.DB.QueryRow(`SELECT id FROM notification WHERE code = ?`, code).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return strconv.FormatInt(id, 10)
}

func status(t *testing.T, res *testutil.Response, want int) {
	t.Helper()
	if res.StatusCode != want {
		t.Fatalf("esperaba %d, llegó %d: %s", want, res.StatusCode, res.Body)
	}
}

func TestOp_listNotifications(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	var pg api.NotificationPage
	res := e.c.Do("GET", "/api/v1/notifications", nil)
	status(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 7 || pg.Items[0].Code != "daily_summary" || len(pg.Items[0].ParamsJSON) != 0 {
		t.Fatalf("%+v", pg.Items)
	}
	status(t, e.c.Do("GET", "/api/v1/notifications?limit=0", nil), 400)
	status(t, e.s.Anonymous(t).Do("GET", "/api/v1/notifications", nil), 401)
}

func TestOp_updateNotification(t *testing.T) {
	e := nuevo(t, testutil.XimenaUser)
	p := "/api/v1/notifications/" + e.notif(t, "budget_threshold")
	var n api.Notification
	res := e.c.Do("PATCH", p, map[string]any{"is_enabled": false, "params_json": map[string]any{"scopes": []string{"total"}}})
	status(t, res, 200)
	res.JSON(t, &n)
	if n.IsEnabled || n.ParamsJSON["scopes"].([]any)[0] != "total" {
		t.Fatalf("%+v", n)
	}
	status(t, e.c.Do("PATCH", "/api/v1/notifications/"+e.notif(t, "daily_summary"), map[string]any{"params_json": map[string]any{"send_time": "10:00"}}), 422)
	status(t, e.c.Do("PATCH", p, map[string]any{"params_json": map[string]any{"scopes": []string{"x"}}}), 422)
	status(t, e.c.Do("PATCH", "/api/v1/notifications/999", map[string]any{"is_enabled": true}), 404)
	status(t, e.c.Raw("PATCH", p, map[string]any{"is_enabled": "sí"}), 400)
}

func TestOp_setNotificationRecipients(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	id := e.notif(t, "payment_due")
	status(t, e.c.Do("PUT", "/api/v1/notifications/"+id+"/recipients", map[string]any{"ids": []int64{e.axelRec}}), 204)
	var pg api.NotificationPage
	e.c.Do("GET", "/api/v1/notifications", nil).JSON(t, &pg)
	for _, n := range pg.Items {
		if n.Code == "payment_due" && (len(n.RecipientIDs) != 1 || n.RecipientIDs[0] != e.axelRec) {
			t.Fatalf("%+v", n)
		}
	}
	res, _ := e.s.DB.Exec(`INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, '3003')`, e.s.Fx.KarenPersonID)
	karen, _ := res.LastInsertId()
	r := e.c.Do("PUT", "/api/v1/notifications/"+id+"/recipients", map[string]any{"ids": []int64{e.axelRec, karen}})
	status(t, r, 422)
	if !strings.Contains(string(r.Body), "not_allowed") {
		t.Fatal("Karen debe rechazarse con not_allowed")
	}
	status(t, e.c.Do("PUT", "/api/v1/notifications/"+id+"/recipients", map[string]any{"ids": []int64{999}}), 422)
	status(t, e.c.Do("PUT", "/api/v1/notifications/999/recipients", map[string]any{"ids": []int64{}}), 404)
}

func TestOp_listNotificationRecipients(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	res := e.c.Do("GET", "/api/v1/notification-recipients", nil)
	status(t, res, 200)
	var pg api.NotificationRecipientPage
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].PersonID != e.s.Fx.AxelPersonID || *pg.Items[0].IsRegistered {
		t.Fatalf("%+v", pg.Items)
	}
	if strings.Contains(string(res.Body), "chat_id") || strings.Contains(string(res.Body), "1001") {
		t.Fatal("nunca devuelve chat_id")
	}
	status(t, e.c.Do("GET", "/api/v1/notification-recipients?sort=chat_id", nil), 400)
}

func TestOp_updateNotificationRecipient(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	p := "/api/v1/notification-recipients/" + strconv.FormatInt(e.axelRec, 10)
	var x api.NotificationRecipient
	res := e.c.Do("PATCH", p, map[string]any{"label": "Celular de Axel", "is_active": false})
	status(t, res, 200)
	res.JSON(t, &x)
	if *x.Label != "Celular de Axel" || x.IsActive {
		t.Fatalf("%+v", x)
	}
	res = e.c.Do("PATCH", p, map[string]any{"label": nil})
	var y api.NotificationRecipient
	res.JSON(t, &y)
	if y.Label != nil {
		t.Fatal("null borra la etiqueta")
	}
	status(t, e.c.Do("PATCH", p, map[string]any{"label": "  "}), 422)
	status(t, e.c.Do("PATCH", "/api/v1/notification-recipients/999", map[string]any{"is_active": true}), 404)
}

func TestOp_sendTestNotification(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	for _, code := range []string{"daily_summary", "payment_due"} { // disparo manual de S7-07 contra el servidor falso
		res := e.c.Do("POST", "/api/v1/notifications/"+e.notif(t, code)+"/test", map[string]any{"person_id": e.s.Fx.AxelPersonID})
		status(t, res, 202)
		var d api.NotificationDelivery
		res.JSON(t, &d)
		if d.Status != "sent" || *d.HTTPStatus != 200 || *d.NotificationCode != code {
			t.Fatalf("%+v", d)
		}
	}
	if len(e.bodies) != 2 || strings.Contains(e.bodies[0], "parse_mode") || !strings.Contains(e.bodies[0], "RESUMEN DEL D") || !strings.Contains(e.bodies[1], "PAGO PR") {
		t.Fatalf("%v", e.bodies)
	}
	var ver string
	e.s.DB.QueryRow(`SELECT coalesce(verified_at, '') FROM notification_recipient WHERE id = ?`, e.axelRec).Scan(&ver)
	if ver == "" {
		t.Fatal("un envío correcto marca verified_at")
	}
	id := e.notif(t, "daily_summary")
	status(t, e.c.Do("POST", "/api/v1/notifications/"+id+"/test", map[string]any{"person_id": e.s.Fx.KarenPersonID}), 422)
	status(t, e.c.Do("POST", "/api/v1/notifications/"+id+"/test", map[string]any{"person_id": e.s.Fx.XimenaPersonID}), 409)
	status(t, e.c.Do("POST", "/api/v1/notifications/999/test", map[string]any{"person_id": e.s.Fx.AxelPersonID}), 404)
	status(t, e.c.Raw("POST", "/api/v1/notifications/"+id+"/test", map[string]any{"person_id": "axel"}), 400)
	if strings.Contains(e.s.Logs.String(), fakeToken) {
		t.Fatal("el token salió en el registro")
	}
	t.Logf("S7-07 (servidor falso) resumen diario y alerta de pago: 202, status sent, HTTP 200, sin parse_mode, token fuera del registro")
}

func candidato(t *testing.T, e *entorno, chat int64, name string) string {
	id, _, err := telegram.RegistrarActualizacion(t.Context(), e.s.DB, []byte(`{"update_id":1,"message":{"chat":{"id":`+strconv.FormatInt(chat, 10)+`,"first_name":"`+name+`"}}}`), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return strconv.FormatInt(id, 10)
}

func TestOp_listTelegramChatCandidates(t *testing.T) {
	e := nuevo(t, testutil.XimenaUser)
	candidato(t, e, 987654321, "Ximena")
	res := e.c.Do("GET", "/api/v1/telegram-chat-candidates", nil)
	status(t, res, 200)
	var pg api.TelegramChatCandidatePage
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].ChatIDMasked != "•••4321" || strings.Contains(string(res.Body), "987654321") {
		t.Fatalf("%s", res.Body)
	}
	status(t, e.c.Do("GET", "/api/v1/telegram-chat-candidates?limit=500", nil), 400)
}

func TestOp_confirmTelegramChatCandidate(t *testing.T) {
	e := nuevo(t, testutil.XimenaUser)
	id := candidato(t, e, 987654321, "Ximena")
	p := "/api/v1/telegram-chat-candidates/" + id + "/confirm"
	res := e.c.Do("POST", p, map[string]any{}, "Idempotency-Key", "0b6f6a52-1c1e-4a57-9a9e-8f1f6c1d2e01")
	status(t, res, 201)
	var x api.NotificationRecipient
	res.JSON(t, &x)
	if x.PersonID != e.s.Fx.XimenaPersonID {
		t.Fatalf("por omisión, la persona de la sesión: %+v", x)
	}
	status(t, e.c.Do("POST", p, map[string]any{}, "Idempotency-Key", "0b6f6a52-1c1e-4a57-9a9e-8f1f6c1d2e02"), 409)
	otro := candidato(t, e, 555123, "Karen")
	status(t, e.c.Do("POST", "/api/v1/telegram-chat-candidates/"+otro+"/confirm", map[string]any{"person_id": e.s.Fx.KarenPersonID}, "Idempotency-Key", "0b6f6a52-1c1e-4a57-9a9e-8f1f6c1d2e03"), 422)
	status(t, e.c.Do("POST", "/api/v1/telegram-chat-candidates/999/confirm", map[string]any{}, "Idempotency-Key", "0b6f6a52-1c1e-4a57-9a9e-8f1f6c1d2e04"), 404)
}

func TestOp_dismissTelegramChatCandidate(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	p := "/api/v1/telegram-chat-candidates/" + candidato(t, e, 4242, "Desconocido")
	status(t, e.c.Do("DELETE", p, nil), 204)
	status(t, e.c.Do("DELETE", p, nil), 404)
}

func TestOp_listNotificationDeliveries(t *testing.T) {
	e := nuevo(t, testutil.AxelUser)
	id := e.notif(t, "daily_summary")
	status(t, e.c.Do("POST", "/api/v1/notifications/"+id+"/test", map[string]any{"person_id": e.s.Fx.AxelPersonID}), 202)
	var pg api.NotificationDeliveryPage
	res := e.c.Do("GET", "/api/v1/notification-deliveries?period=2026-10&notification_id="+id, nil)
	status(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || *pg.Items[0].NotificationCode != "daily_summary" || !strings.HasPrefix(*pg.Items[0].MessageText, "RESUMEN DEL DÍA") {
		t.Fatalf("%s", res.Body)
	}
	e.c.Do("GET", "/api/v1/notification-deliveries?period=2026-09", nil).JSON(t, &pg)
	if len(pg.Items) != 0 {
		t.Fatal("otro periodo no trae el envío")
	}
	status(t, e.c.Do("GET", "/api/v1/notification-deliveries?period=octubre", nil), 400)
}
