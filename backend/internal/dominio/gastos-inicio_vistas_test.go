package dominio_test

import (
	"context"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func TestGastosViewRowsAndTicketQueue(t *testing.T) {
	db, fx := testutil.NewDB(t)
	ctx := context.Background()
	rows, err := dominio.GastosViewRows(ctx, db, dominio.ViewSQL{Table: "user_account", Title: "t.username", Group: "t.person_id",
		GroupLabel: "SELECT display_name FROM person WHERE id = t.person_id", Where: []string{"t.person_id = ?"}, Args: []any{fx.AxelPersonID}, OrderBy: "ORDER BY t.id"})
	if err != nil || len(rows) != 1 {
		t.Fatalf("%v %+v", err, rows)
	}
	r := rows[0]
	if r.ID == 0 || r.Title == "" || !r.Group.Valid || !r.GroupLabel.Valid || r.Amount.Valid || r.Date.Valid || r.Image.Valid {
		t.Fatalf("proyección: %+v", r)
	}
	if _, ok := r.Record["password_hash"]; ok {
		t.Error("password_hash nunca sale en una vista")
	}
	if r.Record["is_active"] != true {
		t.Errorf("is_active debe ser booleano: %#v", r.Record["is_active"])
	}
	if _, ok := r.Record["last_login_at"]; ok {
		t.Error("NULL = propiedad ausente")
	}
	if ca, _ := r.Record["created_at"].(string); len(ca) < 25 || ca[19:] != "-06:00" {
		t.Errorf("instante en −06:00: %q", ca)
	}
	if _, err := dominio.GastosViewRows(ctx, db, dominio.ViewSQL{Table: "no_existe"}); err == nil {
		t.Error("tabla inexistente debe fallar")
	}

	// Cola del OCR: alta pendiente, lectura y fuera de la cola.
	id, err := dominio.GastosInsertTicket(ctx, db, "adjuntos/2026/10/a.png", nil, "image/png", 10, "0000000000000000000000000000000000000000000000000000000000000000", nil)
	if err != nil {
		t.Fatal(err)
	}
	pend, err := dominio.GastosPendingTickets(ctx, db)
	if err != nil || len(pend) != 1 || pend[0].ID != id || pend[0].FilePath != "adjuntos/2026/10/a.png" {
		t.Fatalf("cola: %+v %v", pend, err)
	}
	txt, eng := "LECHE 25.50", "tesseract 5.5.3 spa"
	if err := dominio.GastosSetTicketOCR(ctx, db, id, "done", &eng, &txt, "2026-10-12 18:00:00"); err != nil {
		t.Fatal(err)
	}
	if pend, _ = dominio.GastosPendingTickets(ctx, db); len(pend) != 0 {
		t.Fatalf("la cola debe quedar vacía: %+v", pend)
	}
	if tr, err := dominio.GetTicketScan(ctx, db, id); err != nil || tr.OCRStatus != "done" || tr.OCRText.String != txt {
		t.Fatalf("%+v %v", tr, err)
	}
}
