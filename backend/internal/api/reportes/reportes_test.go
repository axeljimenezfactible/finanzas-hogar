package reportesapi_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"finanzas-hogar/internal/api"
	reportesapi "finanzas-hogar/internal/api/reportes"
	"finanzas-hogar/internal/reportes"
	"finanzas-hogar/internal/testutil"
)

var n atomic.Int64

func key() string { return fmt.Sprintf("s7-pdf-op-%010d", n.Add(1)) }

func srv(t *testing.T) (*testutil.Server, *testutil.Client) {
	s := testutil.NewServer(t, reportesapi.Register)
	return s, s.SessionFor(t, testutil.AxelUser)
}

func want(t *testing.T, res *testutil.Response, code int) {
	t.Helper()
	if res.StatusCode != code {
		t.Fatalf("se esperaba %d y llegó %d: %s", code, res.StatusCode, res.Body)
	}
}

func monthly(t *testing.T, c *testutil.Client) api.GeneratedReport {
	var r api.GeneratedReport
	res := c.Do("POST", "/api/v1/reports/monthly", map[string]any{"month": "2026-10", "send_telegram": false}, "Idempotency-Key", key())
	want(t, res, 202)
	res.JSON(t, &r)
	return r
}

func TestOp_listReports(t *testing.T) {
	s, c := srv(t)
	r := monthly(t, c)
	var pg api.GeneratedReportPage
	res := c.Do("GET", "/api/v1/reports?kind=monthly_pdf&period=2026-10", nil)
	want(t, res, 200)
	res.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].ID != r.ID || *pg.Items[0].DeliveryStatus != "not_sent" {
		t.Fatalf("%+v", pg)
	}
	want(t, c.Do("GET", "/api/v1/reports?kind=otro", nil), 400)
	want(t, s.Anonymous(t).Do("GET", "/api/v1/reports", nil), 401)
}

func TestOp_generateMonthlyReport(t *testing.T) {
	_, c := srv(t)
	if r := monthly(t, c); r.Kind != "monthly_pdf" || r.PeriodStartDate != "2026-10-01" || r.PeriodEndDate != "2026-10-31" || r.SizeBytes == 0 {
		t.Fatalf("%+v", r)
	}
	want(t, c.Do("POST", "/api/v1/reports/monthly", map[string]any{"month": "2026-12", "send_telegram": false}, "Idempotency-Key", key()), 422)
	want(t, c.Raw("POST", "/api/v1/reports/monthly", map[string]any{"month": "2026-10", "send_telegram": false}), 400)
}

func TestOp_generateExternalStatement(t *testing.T) {
	s, c := srv(t)
	body := map[string]any{"person_id": s.Fx.KarenPersonID, "period_start_date": "2026-10-01", "period_end_date": "2026-10-31"}
	want(t, c.Do("POST", "/api/v1/reports/external-statements", body, "Idempotency-Key", key()), 202)
	body["person_id"] = s.Fx.AxelPersonID
	want(t, c.Do("POST", "/api/v1/reports/external-statements", body, "Idempotency-Key", key()), 422)
	body["person_id"], body["period_end_date"] = s.Fx.KarenPersonID, "2026-09-30"
	want(t, c.Do("POST", "/api/v1/reports/external-statements", body, "Idempotency-Key", key()), 422)
}

func TestOp_downloadReportFile(t *testing.T) {
	_, c := srv(t)
	r := monthly(t, c)
	res := c.Do("GET", fmt.Sprintf("/api/v1/reports/%d/file", r.ID), nil)
	want(t, res, 200)
	if int64(len(res.Body)) != r.SizeBytes || string(res.Body[:5]) != "%PDF-" {
		t.Fatalf("PDF de %d bytes", len(res.Body))
	}
	want(t, c.Do("GET", "/api/v1/reports/9999/file", nil), 404)
}

func TestOp_sendReport(t *testing.T) {
	s, c := srv(t)
	var r api.GeneratedReport
	res := c.Do("POST", "/api/v1/reports/external-statements", map[string]any{"person_id": s.Fx.KarenPersonID, "period_start_date": "2026-10-01",
		"period_end_date": "2026-10-31"}, "Idempotency-Key", key())
	want(t, res, 202)
	res.JSON(t, &r)
	want(t, c.Do("POST", fmt.Sprintf("/api/v1/reports/%d/send", r.ID), nil, "Idempotency-Key", key()), 409)
	want(t, c.Do("POST", "/api/v1/reports/9999/send", nil, "Idempotency-Key", key()), 404)

	// 202: un mensual guardado se envía con un Enviador falso (nada sale a la red) a Axel y Ximena, nunca a Karen.
	m := monthly(t, c)
	fake := &fakeSender{}
	reportes.SetSender(fake)
	t.Cleanup(func() { reportes.SetSender(nil) })
	var nid int64
	if err := s.DB.QueryRow(`INSERT INTO notification (code, name) VALUES ('monthly_report', 'Reporte mensual')
		ON CONFLICT (code) DO UPDATE SET is_enabled = 1 RETURNING id`).Scan(&nid); err != nil {
		t.Fatal(err)
	}
	for _, p := range []struct {
		id   int64
		chat string
	}{{s.Fx.AxelPersonID, "1001"}, {s.Fx.XimenaPersonID, "1002"}, {s.Fx.KarenPersonID, "1003"}} {
		var rid int64
		if err := s.DB.QueryRow(`INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, ?) RETURNING id`, p.id, p.chat).Scan(&rid); err != nil {
			t.Fatal(err)
		}
		if _, err := s.DB.Exec(`INSERT INTO notification_subscription (notification_id, recipient_id) VALUES (?, ?)`, nid, rid); err != nil {
			t.Fatal(err)
		}
	}
	var sent api.GeneratedReport
	res = c.Do("POST", fmt.Sprintf("/api/v1/reports/%d/send", m.ID), nil, "Idempotency-Key", key()) // valida el 202 contra openapi.yaml
	want(t, res, 202)
	res.JSON(t, &sent)
	if sent.ID != m.ID || sent.DeliveryStatus == nil || *sent.DeliveryStatus != "sent" {
		t.Fatalf("respuesta del envío: %+v", sent)
	}
	if chats := fake.chats(); len(chats) != 2 || chats[0] == "1003" || chats[1] == "1003" {
		t.Fatalf("envíos a %v; se esperaban Axel y Ximena (Karen nunca)", chats)
	}
}

// fakeSender es el Enviador falso de las pruebas: guarda a quién se mandó el PDF y responde 200 sin salir a la red.
type fakeSender struct {
	mu   sync.Mutex
	sent []string
}

func (f *fakeSender) SendDocument(_ context.Context, chatID, _ string, pdf []byte, _ string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(pdf) == 0 {
		return 400, nil
	}
	f.sent = append(f.sent, chatID)
	return 200, nil
}

func (f *fakeSender) chats() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.sent...)
}

func TestOp_listMonthlyReportMonths(t *testing.T) {
	_, c := srv(t)
	monthly(t, c)
	var m api.ReportMonths
	res := c.Do("GET", "/api/v1/monthly-report/months", nil)
	want(t, res, 200)
	res.JSON(t, &m)
	last := m.Months[len(m.Months)-1]
	if last.Month != "2026-10" || !last.IsPartial || last.GeneratedReportID == nil {
		t.Fatalf("%+v", m)
	}
}

func TestOp_getMonthlyReport(t *testing.T) {
	_, c := srv(t)
	var m api.MonthlyReport
	res := c.Do("GET", "/api/v1/monthly-report", nil)
	want(t, res, 200)
	res.JSON(t, &m)
	if m.Month != "2026-09" || len(m.Fortnights.Fortnights) != 2 {
		t.Fatalf("%+v", m)
	}
	want(t, c.Do("GET", "/api/v1/monthly-report?period=2026-10-01..2026-10-15", nil), 422)
	want(t, c.Do("GET", "/api/v1/monthly-report?period=no", nil), 400)
}

func TestOp_getExternalStatement(t *testing.T) {
	s, c := srv(t)
	want(t, c.Do("GET", fmt.Sprintf("/api/v1/persons/%d/statement", s.Fx.KarenPersonID), nil), 200)
	want(t, c.Do("GET", fmt.Sprintf("/api/v1/persons/%d/statement", s.Fx.AxelPersonID), nil), 422)
	want(t, c.Do("GET", "/api/v1/persons/9999/statement", nil), 404)
}

func TestOp_getReportSettings(t *testing.T) {
	_, c := srv(t)
	var st api.ReportSettings
	res := c.Do("GET", "/api/v1/settings/reports", nil)
	want(t, res, 200)
	res.JSON(t, &st)
	p := st.Parameters
	if p.ReportsMonthlyDayOfMonth.Value.Value != 1 || p.ReportsMonthlySendTime.Value.Value != "08:00" || p.ReportsDailySendTime.Value.Value != "21:30" ||
		p.ReportsWeeklyWeekday.Value.Value != 7 || p.ReportsWeeklySendTime.Value.Value != "19:00" {
		t.Fatalf("%+v", p)
	}
}

func TestOp_updateReportSettings(t *testing.T) {
	_, c := srv(t)
	var st api.ReportSettings
	res := c.Do("PATCH", "/api/v1/settings/reports", map[string]any{"parameters": map[string]any{"reports.monthly_day_of_month": 31}})
	want(t, res, 200)
	res.JSON(t, &st)
	if st.Parameters.ReportsMonthlyDayOfMonth.Value.Value != 31 || st.Parameters.ReportsMonthlyDayOfMonth.IsInitial {
		t.Fatalf("%+v", st.Parameters.ReportsMonthlyDayOfMonth)
	}
	want(t, c.Do("PATCH", "/api/v1/settings/reports", map[string]any{"parameters": map[string]any{"reports.monthly_day_of_month": 32}}), 422)
	want(t, c.Do("PATCH", "/api/v1/settings/reports", map[string]any{"reset_parameters": []string{"reports.monthly_day_of_month"}}), 200)
}
