package sistema_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"finanzas-hogar/internal/adjuntos"
	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/sistema"
	"finanzas-hogar/internal/auth"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/ia"
	"finanzas-hogar/internal/testutil"
)

func init() { auth.HashParams = testutil.FastArgon2 }

var keyN atomic.Int64

func idem() []string {
	return []string{"Idempotency-Key", fmt.Sprintf("prueba-sistema-%08d", keyN.Add(1))}
}

func errCode(t *testing.T, r *testutil.Response) string {
	t.Helper()
	var e struct {
		Code string `json:"code"`
	}
	r.JSON(t, &e)
	return e.Code
}

func TestOp_getHealth(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var h api.Health
	r := s.Anonymous(t).Do("GET", "/health", nil)
	r.JSON(t, &h)
	if r.StatusCode != 200 || h.Status != "ok" || r.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
}

func TestOp_getOpenApiDocument(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	r := s.SessionFor(t, testutil.AxelUser).Do("GET", "/openapi.yaml", nil)
	if r.StatusCode != 200 || !strings.HasPrefix(r.Header.Get("Content-Type"), "application/yaml") || !bytes.Contains(r.Body, []byte("openapi:")) {
		t.Fatalf("%d %s", r.StatusCode, r.Header)
	}
	if s.Anonymous(t).Do("GET", "/openapi.yaml", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}

func TestOp_getSecuritySettings(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var out api.SecuritySettings
	r := s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/settings/security", nil)
	r.JSON(t, &out)
	p := out.Parameters
	if r.StatusCode != 200 || p.AuthMaxFailedAttempts.Value.Value != 5 || p.AuthLockoutWindowMinutes.Value.Value != 15 || p.AuthPasswordMinLength.Value.Value != 12 || p.SessionMaxAgeDays.Value.Value != 90 {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
}

func TestOp_updateSecuritySettings(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	var out api.SecuritySettings
	r := c.Do("PATCH", "/api/v1/settings/security", map[string]any{"parameters": map[string]any{"auth.password_min_length": 14}})
	r.JSON(t, &out)
	if r.StatusCode != 200 || out.Parameters.AuthPasswordMinLength.Value.Value != 14 || out.Parameters.AuthPasswordMinLength.IsInitial {
		t.Fatalf("cambiar: %d %s", r.StatusCode, r.Body)
	}
	// changeOwnPassword usa el valor vigente desde el siguiente intento
	if r := c.Do("POST", "/api/v1/auth/password", map[string]any{"current_password": testutil.FakeCredential, "new_password": "trece-caract1"}); r.StatusCode != 422 {
		t.Errorf("13 caracteres con mínimo 14: %d", r.StatusCode)
	}
	r = c.Do("PATCH", "/api/v1/settings/security", map[string]any{"reset_parameters": []string{"auth.password_min_length"}})
	r.JSON(t, &out)
	if out.Parameters.AuthPasswordMinLength.Value.Value != 12 || !out.Parameters.AuthPasswordMinLength.IsInitial {
		t.Errorf("restablecer: %s", r.Body)
	}
	for _, bad := range []map[string]any{
		{"parameters": map[string]any{"auth.password_min_length": "doce"}},
		{"parameters": map[string]any{"attachment.max_size_mb": 5}},
		{"parameters": map[string]any{"auth.password_min_length": -1}},
	} {
		if r := c.Do("PATCH", "/api/v1/settings/security", bad); r.StatusCode != 422 && r.StatusCode != 400 {
			t.Errorf("%v: %d %s", bad, r.StatusCode, r.Body)
		}
	}
	if r := c.Do("PATCH", "/api/v1/settings/security", map[string]any{}); r.StatusCode != 400 {
		t.Errorf("sin cambios: %d", r.StatusCode)
	}
}

func TestOp_getDataSettings(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var out api.DataSettings
	r := s.SessionFor(t, testutil.AxelUser).Do("GET", "/api/v1/settings/data", nil)
	r.JSON(t, &out)
	if r.StatusCode != 200 || out.Parameters.AttachmentMaxSizeMB.Value.Value != 10 || out.Parameters.IdempotencyRetentionHours.Value.Value != 24 {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
}

func TestOp_updateDataSettings(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	var out api.DataSettings
	r := c.Do("PATCH", "/api/v1/settings/data", map[string]any{"parameters": map[string]any{"attachment.max_size_mb": 3}})
	r.JSON(t, &out)
	if r.StatusCode != 200 || out.Parameters.AttachmentMaxSizeMB.Value.Value != 3 {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	if r := c.Do("PATCH", "/api/v1/settings/data", map[string]any{"parameters": map[string]any{"session.max_age_days": 5}}); r.StatusCode != 422 {
		t.Errorf("clave de otra sección: %d", r.StatusCode)
	}
}

// --- IA (S6-13): apagada por defecto y sin llamadas de salida ---

type failDoer struct{ calls atomic.Int64 }

func (f *failDoer) Do(*http.Request) (*http.Response, error) {
	f.calls.Add(1)
	return nil, errors.New("prueba: la IA no debe llamar a la red")
}

func TestIA_ApagadaNoLlamaALaRed(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	doer := &failDoer{}
	reg := s.Reg.Deps().IA
	ollama := ia.NewOllama(reg, "http://127.0.0.1:9", "modelo-de-prueba", doer)
	claude := ia.NewClaude(reg, "", doer)
	reg.Set(ollama)
	reg.Set(claude)
	c := s.SessionFor(t, testutil.AxelUser)
	var out api.AiProviders
	r := c.Do("GET", "/api/v1/ai/providers", nil)
	r.JSON(t, &out)
	for _, p := range out.Providers {
		if p.IsEnabled || p.IsAvailable || len(p.AllowedUses) != 2 {
			t.Errorf("por defecto: %+v", p)
		}
	}
	ctx := context.Background()
	if _, err := reg.Active(ctx); !errors.Is(err, ia.ErrDisabled) {
		t.Errorf("Active: %v", err)
	}
	if _, err := ollama.ReadTicket(ctx, ia.TicketImage{}); !errors.Is(err, ia.ErrDisabled) {
		t.Errorf("OCR apagado: %v", err)
	}
	if _, err := claude.SuggestRecipe(ctx, ia.RecipeRequest{}); !errors.Is(err, ia.ErrDisabled) {
		t.Errorf("recetas apagadas: %v", err)
	}
	if err := ollama.Ping(ctx); !errors.Is(err, ia.ErrDisabled) {
		t.Errorf("ping apagado: %v", err)
	}
	if n := doer.calls.Load(); n != 0 {
		t.Fatalf("con la IA apagada hubo %d llamadas de salida", n)
	}
	// encendida: ahora sí se consulta (y el cliente de prueba falla)
	var p api.AiProvider
	r = c.Do("PATCH", "/api/v1/ai/providers/local_model", map[string]any{"is_enabled": true})
	r.JSON(t, &p)
	if r.StatusCode != 200 || !p.IsEnabled || p.IsAvailable || doer.calls.Load() != 1 {
		t.Errorf("encendida: %d %s, llamadas %d", r.StatusCode, r.Body, doer.calls.Load())
	}
}

func TestOp_listAiProviders(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var out api.AiProviders
	r := s.SessionFor(t, testutil.AxelUser).Do("GET", "/api/v1/ai/providers", nil)
	r.JSON(t, &out)
	if r.StatusCode != 200 || len(out.Providers) != 2 || out.Providers[0].Provider != "local_model" || out.Providers[1].IsConfigured {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
}

func TestOp_updateAiProvider(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.XimenaUser)
	var p api.AiProvider
	r := c.Do("PATCH", "/api/v1/ai/providers/claude_api", map[string]any{"is_enabled": true})
	r.JSON(t, &p)
	if r.StatusCode != 200 || !p.IsEnabled || p.IsAvailable { // sin adaptador instalado no responde
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	if _, err := s.Reg.Deps().IA.Active(context.Background()); !errors.Is(err, ia.ErrDisabled) {
		t.Errorf("sin adaptador real no hay proveedor activo: %v", err)
	}
	if r := c.Do("PATCH", "/api/v1/ai/providers/claude_api", map[string]any{"is_enabled": false}); r.StatusCode != 200 {
		t.Errorf("apagar: %d", r.StatusCode)
	}
	if r := c.Raw("PATCH", "/api/v1/ai/providers/gpt", map[string]any{"is_enabled": true}); r.StatusCode != 404 && r.StatusCode != 400 {
		t.Errorf("proveedor desconocido: %d", r.StatusCode)
	}
	if r := c.Do("PATCH", "/api/v1/ai/providers/local_model", map[string]any{}); r.StatusCode != 422 {
		t.Errorf("sin is_enabled: %d", r.StatusCode)
	}
}

func TestOp_getSystemStatus(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	var st api.SystemStatus
	r := s.SessionFor(t, testutil.AxelUser).Do("GET", "/api/v1/system/status", nil)
	r.JSON(t, &st)
	if r.StatusCode != 200 || st.Version != sistema.Version || st.DatabaseSizeBytes <= 0 {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
}

func writePEM(t *testing.T, dir, name, typ string, der []byte) string {
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestOp_getCaCertificate(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	dir := testutil.TempDir(t)
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "CA de prueba"}, NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour), IsCA: true, BasicConstraintsValid: true}
	der, _ := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	kder, _ := x509.MarshalPKCS8PrivateKey(key)
	defer func(p string) { sistema.CACertPath = p }(sistema.CACertPath)
	sistema.CACertPath = ""
	if r := c.Do("GET", "/api/v1/system/ca-certificate", nil); r.StatusCode != 503 {
		t.Errorf("sin configurar: %d", r.StatusCode)
	}
	sistema.CACertPath = writePEM(t, dir, "rootCA-key.pem", "PRIVATE KEY", kder)
	if r := c.Do("GET", "/api/v1/system/ca-certificate", nil); r.StatusCode != 503 || bytes.Contains(r.Body, []byte("PRIVATE")) {
		t.Errorf("nunca sirve la llave: %d", r.StatusCode)
	}
	sistema.CACertPath = writePEM(t, dir, "rootCA.pem", "CERTIFICATE", der)
	r := c.Do("GET", "/api/v1/system/ca-certificate", nil)
	if r.StatusCode != 200 || r.Header.Get("Content-Type") != "application/x-pem-file" || !bytes.HasPrefix(r.Body, []byte("-----BEGIN CERTIFICATE-----")) {
		t.Fatalf("%d %s", r.StatusCode, r.Header)
	}
}

func TestOp_getSeedStatus(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	var out api.SeedStatus
	r := c.Do("GET", "/api/v1/system/seed-status", nil)
	r.JSON(t, &out)
	if r.StatusCode != 200 || len(out.Items) != 3 || out.Items[1].IsLoaded {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	for i := 1; i <= 2; i++ {
		if err := dominio.Audit(context.Background(), s.DB, nil, "insert", "credit", int64(i), nil, []byte(`{}`), "seed:kueski_credit"); err != nil {
			t.Fatal(err)
		}
	}
	r = c.Do("GET", "/api/v1/system/seed-status", nil)
	r.JSON(t, &out)
	if !out.Items[1].IsLoaded || out.Items[1].RowCount != 2 || out.Items[0].IsLoaded {
		t.Errorf("tras cargar: %s", r.Body)
	}
}

func TestOp_listJobs(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	ctx := context.Background()
	for _, j := range [][3]string{{"backup_daily", "Respaldo diario", "15 3 * * *"}, {"idempotency_purge", "Purga", "5 * * * *"}} {
		if err := dominio.EnsureJob(ctx, s.DB, j[0], j[1], j[2]); err != nil {
			t.Fatal(err)
		}
	}
	var pg api.JobPage
	r := s.SessionFor(t, testutil.AxelUser).Do("GET", "/api/v1/jobs?limit=1", nil)
	r.JSON(t, &pg)
	if r.StatusCode != 200 || len(pg.Items) != 1 || pg.Items[0].Code != "backup_daily" || !pg.NextCursor.Valid {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	r = s.SessionFor(t, testutil.AxelUser).Do("GET", "/api/v1/jobs?cursor="+pg.NextCursor.Value+"&limit=1", nil)
	r.JSON(t, &pg)
	if len(pg.Items) != 1 || pg.Items[0].Code != "idempotency_purge" {
		t.Errorf("segunda página: %s", r.Body)
	}
}

func TestOp_updateJob(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	ctx := context.Background()
	if err := dominio.EnsureJob(ctx, s.DB, "backup_daily", "Respaldo diario", "15 3 * * *"); err != nil {
		t.Fatal(err)
	}
	j, _ := dominio.GetJobByCode(ctx, s.DB, "backup_daily")
	c := s.SessionFor(t, testutil.XimenaUser)
	var out api.Job
	r := c.Do("PATCH", fmt.Sprintf("/api/v1/jobs/%d", j.ID), map[string]any{"schedule": "30 4 * * *", "is_enabled": false})
	r.JSON(t, &out)
	if r.StatusCode != 200 || out.Schedule != "30 4 * * *" || out.IsEnabled || out.NextRunAt == nil || !strings.HasPrefix(*out.NextRunAt, "2026-10-13T04:30:00") {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	if r := c.Do("PATCH", fmt.Sprintf("/api/v1/jobs/%d", j.ID), map[string]any{"schedule": "a las 4"}); r.StatusCode != 422 {
		t.Errorf("cron malo: %d", r.StatusCode)
	}
	if r := c.Do("PATCH", "/api/v1/jobs/999", map[string]any{"is_enabled": true}); r.StatusCode != 404 {
		t.Errorf("no existe: %d", r.StatusCode)
	}
	if r := c.Do("PATCH", fmt.Sprintf("/api/v1/jobs/%d", j.ID), map[string]any{}); r.StatusCode != 400 {
		t.Errorf("sin cambios: %d", r.StatusCode)
	}
}

// --- Respaldos (S6-09) ---

func TestOp_createBackup(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	h := idem()
	var b api.Backup
	r := c.Do("POST", "/api/v1/backups", nil, h...)
	r.JSON(t, &b)
	if r.StatusCode != 202 || b.Status != "ok" || *b.IntegrityCheck != "ok" || *b.RestoreTestResult != "ok" || b.Kind != "manual" {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	if _, err := os.Stat(*b.FilePath); err != nil || !strings.HasPrefix(*b.FilePath, s.Reg.Deps().Config.BackupDir) {
		t.Errorf("archivo del respaldo: %v", err)
	}
	if r := c.Do("POST", "/api/v1/backups", nil, h...); r.Header.Get("Idempotent-Replayed") != "true" {
		t.Errorf("repetición: %d %s", r.StatusCode, r.Header)
	}
	if r := c.Do("POST", "/api/v1/backups", nil); r.StatusCode != 400 {
		t.Errorf("sin Idempotency-Key: %d", r.StatusCode)
	}
	if _, err := s.DB.Exec(`INSERT INTO backup (kind, status) VALUES ('scheduled', 'running')`); err != nil {
		t.Fatal(err)
	}
	if r := c.Do("POST", "/api/v1/backups", nil, idem()...); r.StatusCode != 409 {
		t.Errorf("con otro en curso: %d", r.StatusCode)
	}
}

func TestOp_listBackups(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	c.Do("POST", "/api/v1/backups", nil, idem()...)
	s.Clock.Advance(time.Hour)
	c.Do("POST", "/api/v1/backups", nil, idem()...)
	var pg api.BackupPage
	r := c.Do("GET", "/api/v1/backups", nil)
	r.JSON(t, &pg)
	if r.StatusCode != 200 || len(pg.Items) != 2 || pg.Items[0].ID < pg.Items[1].ID {
		t.Fatalf("orden -started_at: %d %s", r.StatusCode, r.Body)
	}
	if r := c.Do("GET", "/api/v1/backups?sort=sha256", nil); r.StatusCode != 400 {
		t.Errorf("orden no permitido: %d", r.StatusCode)
	}
}

func TestOp_restoreBackup(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	var b api.Backup
	c.Do("POST", "/api/v1/backups", nil, idem()...).JSON(t, &b)
	if _, err := s.DB.Exec(`UPDATE person SET display_name = 'Karen después del respaldo' WHERE code = 'karen'`); err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("/api/v1/backups/%d/restore", b.ID)
	if r := c.Do("POST", path, map[string]any{"confirmation": "restaurar"}); r.StatusCode != 422 {
		t.Errorf("confirmación mala: %d", r.StatusCode)
	}
	if r := c.Do("POST", "/api/v1/backups/999/restore", map[string]any{"confirmation": "RESTAURAR"}); r.StatusCode != 404 {
		t.Errorf("no existe: %d", r.StatusCode)
	}
	var out api.RestoreResult
	r := c.Do("POST", path, map[string]any{"confirmation": "RESTAURAR"})
	r.JSON(t, &out)
	if r.StatusCode != 202 || out.IntegrityCheck != "ok" || out.PreRestoreBackupID == 0 {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	var name string
	s.DB.QueryRow(`SELECT display_name FROM person WHERE code = 'karen'`).Scan(&name)
	if name != "Karen" {
		t.Errorf("no se restauró: %q", name)
	}
	// el catálogo de respaldos y la sesión actual se conservan (PROVISIONAL, bloqueo S6-respaldos)
	var pre api.Backup
	var pg api.BackupPage
	c.Do("GET", "/api/v1/backups", nil).JSON(t, &pg)
	for _, it := range pg.Items {
		if it.ID == out.PreRestoreBackupID {
			pre = it
		}
	}
	if pre.Kind != "pre_migration" || pre.Status != "ok" {
		t.Errorf("respaldo previo: %+v", pre)
	}
	os.Remove(*pre.FilePath)
	if r := c.Do("POST", fmt.Sprintf("/api/v1/backups/%d/restore", pre.ID), map[string]any{"confirmation": "RESTAURAR"}); r.StatusCode != 409 {
		t.Errorf("archivo borrado: %d", r.StatusCode)
	}
}

// --- Adjuntos (S6-10) ---

func expense(t *testing.T, s *testutil.Server) int64 {
	t.Helper()
	var cat, id int64
	if err := s.DB.QueryRow(`INSERT INTO category (code, name) VALUES ('groceries', 'Despensa') RETURNING id`).Scan(&cat); err != nil {
		t.Fatal(err)
	}
	if err := s.DB.QueryRow(`INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, payment_method, paid_by_person_id)
		VALUES ('variable', 'Súper', ?, 10000, '2026-10-05', 'cash', ?) RETURNING id`, cat, s.Fx.AxelPersonID).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

var pngHead = []byte("\x89PNG\r\n\x1a\n")

func pngOf(size int) []byte {
	b := make([]byte, size)
	copy(b, pngHead)
	return b
}

func upload(t *testing.T, c *testutil.Client, expenseID int64, name string, content []byte) *testutil.Response {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.WriteField("kind", "ticket_photo")
	fw, _ := mw.CreateFormFile("file", name)
	fw.Write(content)
	mw.Close()
	return c.Do("POST", fmt.Sprintf("/api/v1/expenses/%d/attachments", expenseID), buf.Bytes(), append([]string{"Content-Type", mw.FormDataContentType()}, idem()...)...)
}

func TestOp_uploadExpenseAttachment(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	id := expense(t, s)
	var a api.Attachment
	r := upload(t, c, id, "ticket.png", pngOf(2048))
	r.JSON(t, &a)
	if r.StatusCode != 201 || a.MimeType != "image/png" || a.SizeBytes != 2048 || *a.ExpenseID != id || a.Kind != "ticket_photo" {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	if r := upload(t, c, 999, "ticket.png", pngOf(10)); r.StatusCode != 404 {
		t.Errorf("gasto que no existe: %d", r.StatusCode)
	}
	if r := upload(t, c, id, "nota.txt", []byte("hola, esto es texto")); r.StatusCode != 415 {
		t.Errorf("texto: %d", r.StatusCode)
	}
	if r := c.Do("POST", fmt.Sprintf("/api/v1/expenses/%d/attachments", id), map[string]any{"file": "x"}, idem()...); r.StatusCode != 415 {
		t.Errorf("sin multipart: %d", r.StatusCode)
	}
}

func TestOp_downloadAttachment(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.XimenaUser)
	var a api.Attachment
	content := pngOf(4096)
	upload(t, c, expense(t, s), "foto.png", content).JSON(t, &a)
	r := c.Do("GET", fmt.Sprintf("/api/v1/attachments/%d/file", a.ID), nil)
	if r.StatusCode != 200 || !bytes.Equal(r.Body, content) || r.Header.Get("Content-Type") != "application/octet-stream" {
		t.Fatalf("%d %s", r.StatusCode, r.Header)
	}
	if r := c.Do("GET", "/api/v1/attachments/999/file", nil); r.StatusCode != 404 {
		t.Errorf("no existe: %d", r.StatusCode)
	}
}

func TestOp_deleteAttachment(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	var a api.Attachment
	upload(t, c, expense(t, s), "foto.png", pngOf(100)).JSON(t, &a)
	var rel string
	s.DB.QueryRow(`SELECT file_path FROM attachment WHERE id = ?`, a.ID).Scan(&rel)
	full := filepath.Join(s.Reg.Deps().Config.DataDir, rel)
	if r := c.Do("DELETE", fmt.Sprintf("/api/v1/attachments/%d", a.ID), nil); r.StatusCode != 204 {
		t.Fatalf("%d %s", r.StatusCode, r.Body)
	}
	if _, err := os.Stat(full); !os.IsNotExist(err) {
		t.Errorf("el archivo sigue: %v", err)
	}
	if r := c.Do("DELETE", fmt.Sprintf("/api/v1/attachments/%d", a.ID), nil); r.StatusCode != 404 {
		t.Errorf("ya borrado: %d", r.StatusCode)
	}
}

func TestAdjunto_SubirDescargarConYSinSesion(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	var a api.Attachment
	upload(t, c, expense(t, s), "Mi ticket del súper.png", pngOf(1000)).JSON(t, &a)
	var rel string
	s.DB.QueryRow(`SELECT file_path FROM attachment WHERE id = ?`, a.ID).Scan(&rel)
	if !regexp.MustCompile(`^adjuntos/2026/10/[0-9a-f]{32}\.png$`).MatchString(rel) || *a.OriginalName != "Mi ticket del súper.png" {
		t.Errorf("nombre derivado del usuario: %q", rel)
	}
	path := fmt.Sprintf("/api/v1/attachments/%d/file", a.ID)
	if r := c.Do("GET", path, nil); r.StatusCode != 200 {
		t.Errorf("con sesión: %d", r.StatusCode)
	}
	if r := s.Anonymous(t).Do("GET", path, nil); r.StatusCode != 401 {
		t.Errorf("sin sesión: %d", r.StatusCode)
	}
}

func TestAdjunto_RutaConPuntosRechazada(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	id := expense(t, s)
	if r := upload(t, c, id, "../../fuera.png", pngOf(100)); r.StatusCode != 422 {
		t.Errorf("nombre con ../: %d %s", r.StatusCode, r.Body)
	}
	if r := upload(t, c, id, `..\fuera.png`, pngOf(100)); r.StatusCode != 422 {
		t.Errorf("nombre con ..\\: %d", r.StatusCode)
	}
	// un renglón con ruta que sale de la carpeta (base alterada) no se sirve
	data := s.Reg.Deps().Config.DataDir
	os.MkdirAll(data, 0o700)
	os.WriteFile(filepath.Join(data, "secreto.png"), pngOf(10), 0o600)
	var bad int64
	if err := s.DB.QueryRow(`INSERT INTO attachment (expense_id, kind, file_path, mime_type, size_bytes, sha256) VALUES (?, 'other', 'adjuntos/../secreto.png', 'image/png', 10, ?) RETURNING id`,
		id, strings.Repeat("0", 64)).Scan(&bad); err != nil {
		t.Fatal(err)
	}
	if r := c.Do("GET", fmt.Sprintf("/api/v1/attachments/%d/file", bad), nil); r.StatusCode != 404 {
		t.Errorf("ruta con ../ en la base: %d", r.StatusCode)
	}
	if r := c.Raw("GET", "/api/v1/attachments/..%2F..%2Fsecreto/file", nil); r.StatusCode != 400 && r.StatusCode != 404 {
		t.Errorf("id con ../: %d", r.StatusCode)
	}
}

// Factor PROVISIONAL de C1 (S6-10): 1 MB = adjuntos.BytesPerMB = 1,048,576 bytes.
func TestAdjunto_Limite10MBYLimiteCambiado(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	c := s.SessionFor(t, testutil.AxelUser)
	id := expense(t, s)
	mb := int(adjuntos.BytesPerMB)
	if r := upload(t, c, id, "exacto.png", pngOf(10*mb)); r.StatusCode != 201 {
		t.Errorf("exactamente 10 MB: %d %s", r.StatusCode, r.Body)
	}
	r := upload(t, c, id, "grande.png", pngOf(10*mb+1))
	if r.StatusCode != 413 || errCode(t, r) != "payload_too_large" || !strings.Contains(string(r.Body), "10 MB") {
		t.Errorf("10 MB + 1 byte: %d %s", r.StatusCode, r.Body)
	}
	if r := c.Do("PATCH", "/api/v1/settings/data", map[string]any{"parameters": map[string]any{"attachment.max_size_mb": 1}}); r.StatusCode != 200 {
		t.Fatalf("cambiar el límite: %d %s", r.StatusCode, r.Body)
	}
	if r := upload(t, c, id, "uno.png", pngOf(mb)); r.StatusCode != 201 {
		t.Errorf("1 MB con límite 1: %d", r.StatusCode)
	}
	if r := upload(t, c, id, "dos.png", pngOf(mb+1)); r.StatusCode != 413 || !strings.Contains(string(r.Body), "1 MB") {
		t.Errorf("1 MB + 1 con límite 1: %d %s", r.StatusCode, r.Body)
	}
	var n int
	s.DB.QueryRow(`SELECT COUNT(*) FROM attachment`).Scan(&n)
	if n != 2 {
		t.Errorf("adjuntos guardados: %d", n)
	}
}

// --- Registro (S6-15): sin contraseñas, cookies ni tokens ---

func TestLog_SinContrasenasNiCookies(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	const mala, nueva = "contrasena-que-no-debe-salir", "nueva-contrasena-secreta-1"
	an := s.Anonymous(t)
	an.Do("POST", "/api/v1/auth/login", map[string]any{"username": "axel", "password": mala})
	res := an.Do("POST", "/api/v1/auth/login", map[string]any{"username": "axel", "password": testutil.FakeCredential})
	var token string
	for _, ck := range res.Cookies() {
		if ck.Name == api.SessionCookie {
			token = ck.Value
		}
	}
	c := s.Anonymous(t)
	c.Token = token
	c.Do("GET", "/api/v1/auth/me", nil)
	c.Do("POST", "/api/v1/auth/password", map[string]any{"current_password": testutil.FakeCredential, "new_password": nueva})
	logs := s.Logs.String()
	if token == "" || !strings.Contains(logs, "/api/v1/auth/login") {
		t.Fatalf("sin registro o sin token: %q", logs)
	}
	for _, secret := range []string{mala, nueva, testutil.FakeCredential, token, dominio.HashToken(token), api.SessionCookie + "="} {
		if strings.Contains(logs, secret) {
			t.Errorf("el registro contiene un secreto (%d caracteres)", len(secret))
		}
	}
}
