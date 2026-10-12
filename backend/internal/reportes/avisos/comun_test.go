package avisos

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"sync"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/telegram"
)

type enviado struct{ chat, text string }

// buzon es el servidor falso: registra lo que se mandaría a Telegram.
type buzon struct {
	mu   sync.Mutex
	msgs []enviado
}

func (b *buzon) SendMessage(_ context.Context, chat, text string) (telegram.Result, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.msgs = append(b.msgs, enviado{chat, text})
	return telegram.Result{HTTPStatus: 200, Parts: 1}, nil
}

// con cuenta y saca los mensajes cuyo encabezado empieza con h.
func (b *buzon) con(h string) []enviado {
	b.mu.Lock()
	defer b.mu.Unlock()
	var out, rest []enviado
	for _, m := range b.msgs {
		if strings.HasPrefix(m.text, h) {
			out = append(out, m)
		} else {
			rest = append(rest, m)
		}
	}
	b.msgs = rest
	return out
}

// fuenteFija usa la base para gastos y casillas y fixtures para lo demás.
type fuenteFija struct {
	FuenteBD
	pagos    []Pago
	usos     []Uso
	despensa *dominio.Date
	viejos   int
}

func (f *fuenteFija) Pagos(_ context.Context, desde, hasta dominio.Date) ([]Pago, error) {
	var out []Pago
	for _, p := range f.pagos {
		if !p.Fecha.Before(desde) && !p.Fecha.After(hasta) {
			out = append(out, p)
		}
	}
	return out, nil
}
func (f *fuenteFija) Presupuesto(context.Context, dominio.Period) ([]Uso, error) { return f.usos, nil }
func (f *fuenteFija) Despensa(context.Context, dominio.Date) (dominio.Date, bool, error) {
	if f.despensa == nil {
		return dominio.Date{}, false, nil
	}
	return *f.despensa, true, nil
}
func (f *fuenteFija) PreciosViejos(context.Context, dominio.Date) (int, error) { return f.viejos, nil }

// destinos da chat a Axel (1001) y Ximena (2002) y los suscribe a los siete avisos.
func destinos(t *testing.T, db *sql.DB, axel, ximena int64) {
	t.Helper()
	for _, q := range []string{
		`INSERT INTO notification_recipient (person_id, chat_id) VALUES (` + itoa(axel) + `, '1001'), (` + itoa(ximena) + `, '2002')`,
		`INSERT INTO notification_subscription (notification_id, recipient_id) SELECT n.id, r.id FROM notification n, notification_recipient r`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }

func fecha(t *testing.T, s string) dominio.Date {
	t.Helper()
	d, err := dominio.ParseDate(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
