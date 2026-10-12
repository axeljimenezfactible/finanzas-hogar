package gastos_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"strings"
	"testing"

	"finanzas-hogar/internal/api"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios/ocr"
	"finanzas-hogar/internal/testutil"
)

// fakeOCR devuelve renglones fijos o un error, sin llamar a Tesseract.
type fakeOCR struct {
	lines []ocr.Line
	err   error
	calls int
}

func (f *fakeOCR) Recognize(context.Context, string) ([]ocr.Line, error) {
	f.calls++
	return f.lines, f.err
}

var pngHead = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00")

func form(t *testing.T, field, name string, content []byte) ([]byte, string) {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	if field != "" {
		fw, err := mw.CreateFormFile(field, name)
		if err != nil {
			t.Fatal(err)
		}
		fw.Write(content)
	}
	mw.Close()
	return b.Bytes(), mw.FormDataContentType()
}

func upload(c *testutil.Client, body []byte, ctype string, n int) *testutil.Response {
	return c.Do("POST", "/api/v1/ticket-scans", body, append(key(n), "Content-Type", ctype)...)
}

func TestOp_createTicketScan(t *testing.T) {
	f := &fakeOCR{lines: []ocr.Line{{Texto: "LECHE ENTERA 1L 25.50"}, {Texto: "TOTAL 25.50"}}}
	defer gastos.SetOCR(f)()
	s := testutil.NewServer(t, gastos.Register)
	ctx := context.Background()
	jb := s.Reg.Deps().Jobs
	if err := jb.Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	c := s.SessionFor(t, testutil.AxelUser)

	body, ct := form(t, "file", "ticket.png", pngHead)
	res := upload(c, body, ct, 901)
	want(t, res, 202)
	var ts api.TicketScan
	res.JSON(t, &ts)
	if ts.OCRStatus != "pending" || ts.Kind != "ticket_photo" || ts.ExpenseID != nil || ts.MimeType != "image/png" || *ts.OriginalName != "ticket.png" || len(ts.Lines) != 0 {
		t.Fatalf("pendiente: %+v", ts)
	}
	// El job lee la cola: texto crudo y motor; getTicketScan propone los renglones (sin el TOTAL).
	if n, err := jb.Tick(ctx); err != nil || n != 1 || f.calls != 1 {
		t.Fatalf("tick: %d %v %d", n, err, f.calls)
	}
	res = c.Do("GET", "/api/v1/ticket-scans/"+itoa(ts.ID), nil)
	want(t, res, 200)
	res.JSON(t, &ts)
	if ts.OCRStatus != "done" || ts.OCREngine == nil || *ts.OCREngine != "tesseract 5.5.3 spa" || len(ts.Lines) != 1 || ts.Lines[0].LineTotalCents == nil || *ts.Lines[0].LineTotalCents != 2550 {
		t.Fatalf("leído: %+v", ts)
	}
	var audits int
	s.DB.QueryRow(`SELECT count(*) FROM audit_log WHERE table_name = 'attachment' AND row_id = ?`, ts.ID).Scan(&audits)
	if audits != 2 {
		t.Errorf("audit_log: %d, quiero alta y lectura", audits)
	}

	// Si el OCR falla, el ticket queda failed (se captura a mano).
	f.err = errors.New("tesseract no está instalado")
	body, ct = form(t, "file", "", pngHead)
	res = upload(c, body, ct, 902)
	want(t, res, 202)
	res.JSON(t, &ts)
	if err := jb.RunNow(ctx, "ticket_ocr"); err != nil {
		t.Fatal(err)
	}
	if _, err := jb.Tick(ctx); err != nil {
		t.Fatal(err)
	}
	if tr, err := dominio.GetTicketScan(ctx, s.DB, ts.ID); err != nil || tr.OCRStatus != "failed" || tr.OCRText.Valid {
		t.Fatalf("fallido: %+v %v", tr, err)
	}

	// Errores: PDF o texto = 415; sin multipart = 415; campo desconocido = 400; sin foto = 422; nombre con ruta = 422;
	// muy grande = 413; sin Idempotency-Key = 400; sin sesión = 401.
	pdf, ct := form(t, "file", "t.pdf", []byte("%PDF-1.4\n%prueba\n"))
	want(t, upload(c, pdf, ct, 903), 415)
	txt, ct := form(t, "file", "t.txt", []byte("hola"))
	want(t, upload(c, txt, ct, 904), 415)
	want(t, upload(c, []byte(`{}`), "application/json", 905), 415)
	other, ct := form(t, "photo", "t.png", pngHead)
	want(t, upload(c, other, ct, 906), 400)
	empty, ct := form(t, "", "", nil)
	want(t, upload(c, empty, ct, 907), 422)
	bad, ct := form(t, "file", "../t.png", pngHead)
	want(t, upload(c, bad, ct, 908), 422)
	maxMB, err := dominio.SettingInt(ctx, s.DB, "attachment.max_size_mb")
	if err != nil {
		t.Fatal(err)
	}
	big, ct := form(t, "file", "t.png", append(append([]byte{}, pngHead...), bytes.Repeat([]byte{0}, int(maxMB)<<20)...))
	want(t, upload(c, big, ct, 909), 413)
	body, ct = form(t, "file", "t.png", pngHead)
	want(t, c.Do("POST", "/api/v1/ticket-scans", body, "Content-Type", ct), 400)
	if s.Anonymous(t).Do("POST", "/api/v1/ticket-scans", body, append(key(910), "Content-Type", ct)...).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
	if !strings.Contains(s.Logs.String(), "ticket_ocr") {
		t.Error("la lectura fallida debe quedar en el registro")
	}
}
