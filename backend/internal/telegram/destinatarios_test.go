package telegram

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func alta(t *testing.T, db *sql.DB, personID int64, chat string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, ?)`, personID, chat)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return id
}

func suscribir(t *testing.T, db *sql.DB, notif string, rid int64) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO notification_subscription (notification_id, recipient_id) SELECT id, ? FROM notification WHERE code = ?`, rid, notif); err != nil {
		t.Fatal(err)
	}
}

func notifID(t *testing.T, db *sql.DB, code string) int64 {
	var id int64
	if err := db.QueryRow(`SELECT id FROM notification WHERE code = ?`, code).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func TestDestinatarios(t *testing.T) { // S7-09
	ctx := context.Background()
	db, fx := testutil.NewDB(t)
	now := testutil.MX(2026, time.October, 12, 12, 0)
	axel := alta(t, db, fx.AxelPersonID, "1001")
	suscribir(t, db, "daily_summary", axel)
	n := notifID(t, db, "daily_summary")

	// Ximena sin chat: se omite y se reporta; el de Axel sale.
	ds, om, err := Destinatarios(ctx, db, n)
	if err != nil || len(ds) != 1 || ds[0].PersonID != fx.AxelPersonID || len(om) != 1 || om[0].PersonID != fx.XimenaPersonID {
		t.Fatalf("%+v %+v %v", ds, om, err)
	}
	t.Logf("S7-09 Ximena sin chat: omitida (%s); Axel sale", om[0].Razon)

	// Karen: no puede ser destinataria (lista de ids ni confirmación de un chat).
	var karenRecipient int64
	if res, err := db.Exec(`INSERT INTO notification_recipient (person_id, chat_id) VALUES (?, '3003')`, fx.KarenPersonID); err == nil {
		karenRecipient, _ = res.LastInsertId()
	}
	if karenRecipient != 0 {
		var de *dominio.Error
		if err := ValidarDestinatarios(ctx, db, []int64{axel, karenRecipient}); !errors.As(err, &de) || de.Fields[0].Code != "not_allowed" {
			t.Fatalf("Karen en la lista debe dar 422 not_allowed: %v", err)
		}
		suscribir(t, db, "daily_summary", karenRecipient) // aun suscrita a mano, no sale
		if ds, _, _ := Destinatarios(ctx, db, n); len(ds) != 1 {
			t.Fatalf("Karen no debe salir: %+v", ds)
		}
	}
	if err := PersonaDelHogar(ctx, db, fx.KarenPersonID, "person_id"); err == nil {
		t.Fatal("Karen debe ser rechazada")
	}
	t.Logf("S7-09 Karen rechazada como destinataria")

	// Mensaje entrante de un chat nuevo: queda candidato pendiente y no recibe avisos.
	upd := []byte(`{"update_id":7,"message":{"message_id":1,"chat":{"id":2002,"type":"private","first_name":"Ximena"},"text":"hola"}}`)
	cid, nuevo, err := RegistrarActualizacion(ctx, db, upd, now)
	if err != nil || !nuevo {
		t.Fatal(err)
	}
	if _, nuevo, _ := RegistrarActualizacion(ctx, db, upd, now); nuevo {
		t.Fatal("el mismo chat no se anota dos veces")
	}
	if ds, _, _ := Destinatarios(ctx, db, n); len(ds) != 1 {
		t.Fatal("un candidato pendiente no recibe avisos")
	}
	// Confirmar con Karen: rechazado; con Ximena: pasa a destinatario.
	err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		_, err := ConfirmarCandidato(ctx, tx, cid, fx.KarenPersonID, now, nil, "t")
		return err
	})
	if err == nil {
		t.Fatal("confirmar a Karen debe fallar")
	}
	var rid int64
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) (err error) {
		rid, err = ConfirmarCandidato(ctx, tx, cid, fx.XimenaPersonID, now, nil, "t")
		return err
	}); err != nil {
		t.Fatal(err)
	}
	suscribir(t, db, "daily_summary", rid)
	ds, om, _ = Destinatarios(ctx, db, n)
	if len(ds) != 2 || ds[1].PersonID != fx.XimenaPersonID || ds[1].ChatID != "2002" || len(om) != 0 {
		t.Fatalf("%+v %+v", ds, om)
	}
	if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		_, err := ConfirmarCandidato(ctx, tx, cid, fx.XimenaPersonID, now, nil, "t")
		return err
	}); err == nil {
		t.Fatal("confirmar dos veces debe dar 409")
	}
	if MaskChatID("2002") != "•••2002" || MaskChatID("-1001234567") != "•••4567" {
		t.Fatal("máscara")
	}
	t.Logf("S7-09 chat nuevo → candidato pendiente sin avisos; confirmado por Ximena → destinatario; Axel y Ximena reciben")
}
