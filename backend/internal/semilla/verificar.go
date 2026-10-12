package semilla

import (
	"context"
	"database/sql"
	"fmt"

	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// Cruce son las cifras que motor_cred calcula sobre la base sembrada (C3-09) y las del oráculo con que se cotejan.
type Cruce struct {
	DisponibleA, DisponibleB int64 // 10-oct-2026, tras la despensa 1
	PrimerPago               int64 // pago del 30-oct-2026
	PrimerPagoFecha          string
	DeudaViva                int64 // al 10-abr-2027
	DeudaVivaFecha           string
	Oraculo                  Cruce1
}

// Cruce1 son las cifras del oráculo (oraculo/model.json).
type Cruce1 struct{ A, B, Pago, Deuda int64 }

// CreditInput arma el Input de motor_cred desde la base, como lo hace la API: cada compra se vuelve a cotizar con sus
// intereses ya fijados y su primer pago, así la serie sale igual a la guardada.
func CreditInput(ctx context.Context, q dominio.Querier, today mc.Date) (mc.Input, error) {
	persons := map[int64]struct {
		code string
		ext  bool
	}{}
	rows, err := q.QueryContext(ctx, `SELECT id, code, participation = 'external' FROM person`)
	if err != nil {
		return mc.Input{}, err
	}
	for rows.Next() {
		var id int64
		var code string
		var ext bool
		if err := rows.Scan(&id, &code, &ext); err != nil {
			rows.Close()
			return mc.Input{}, err
		}
		persons[id] = struct {
			code string
			ext  bool
		}{code, ext}
	}
	rows.Close()
	cr, err := q.QueryContext(ctx, `SELECT id, name, kind, holder_person_id, credit_limit_cents, payment_frequency, first_payment_min_days,
		availability_reading, availability_confirmed, reading_a_anchor_date FROM credit WHERE is_active = 1 ORDER BY id`)
	if err != nil {
		return mc.Input{}, err
	}
	var credits []mc.Credit
	for cr.Next() {
		var c mc.Credit
		var holder int64
		var kind, reading string
		var limit, minDays sql.NullInt64
		var anchor sql.NullString
		var conf bool
		if err := cr.Scan(&c.ID, &c.Name, &kind, &holder, &limit, &c.Frequency, &minDays, &reading, &conf, &anchor); err != nil {
			cr.Close()
			return mc.Input{}, err
		}
		c.Kind, c.HolderPersonID, c.Holder, c.Reading, c.ReadingConfirmed = mc.Kind(kind), holder, persons[holder].code, mc.Reading(reading), conf
		if limit.Valid {
			v := mc.Cents(limit.Int64)
			c.Limit = &v
		}
		if minDays.Valid {
			v := int(minDays.Int64)
			c.FirstPaymentMinDays = &v
		}
		if anchor.Valid {
			d := mcDate(anchor.String)
			c.ReadingAAnchor = &d
		}
		c.Promotions = []mc.Promotion{}
		credits = append(credits, c)
	}
	cr.Close()
	in := mc.Input{Today: today, Cash: []mc.CashEvent{}}
	for _, c := range credits {
		tr, err := q.QueryContext(ctx, `SELECT term_count, rate_num, rate_den, coalesce(rate_decimals, -1), source FROM credit_term_rate WHERE credit_id = ? ORDER BY term_count`, c.ID)
		if err != nil {
			return mc.Input{}, err
		}
		for tr.Next() {
			var n, num, den, dec int64
			var src string
			if err := tr.Scan(&n, &num, &den, &dec, &src); err != nil {
				tr.Close()
				return mc.Input{}, err
			}
			c.TermRates = append(c.TermRates, mc.TermRate{TermCount: int(n), Rate: mc.Rate{Num: num, Den: den, Decimals: int(dec)}, Source: src})
		}
		tr.Close()
		st := mc.CreditState{Credit: c, Payments: []mc.PaymentRecord{}, Statements: []mc.StatementInput{}}
		pr, err := q.QueryContext(ctx, `SELECT id, purchase_date, principal_cents, interest_cents, term_count, first_payment_date, rate_decimals, interest_source, status
			FROM credit_purchase WHERE credit_id = ? ORDER BY purchase_date, id`, c.ID)
		if err != nil {
			return mc.Input{}, err
		}
		type prow struct {
			id, principal, interest, terms int64
			date, first, source, status    string
			dec                            sql.NullInt64
		}
		var ps []prow
		for pr.Next() {
			var p prow
			if err := pr.Scan(&p.id, &p.date, &p.principal, &p.interest, &p.terms, &p.first, &p.dec, &p.source, &p.status); err != nil {
				pr.Close()
				return mc.Input{}, err
			}
			ps = append(ps, p)
		}
		pr.Close()
		for _, p := range ps {
			fp := mcDate(p.first)
			pi := mc.PurchaseInput{Credit: c, Date: mcDate(p.date), Principal: mc.Cents(p.principal), TermCount: int(p.terms), FirstPaymentDate: &fp}
			if p.dec.Valid {
				v := int(p.dec.Int64)
				pi.RateDecimals = &v
			}
			if p.source != "amortization" {
				qi := mc.Cents(p.interest)
				pi.QuotedInterest = &qi
			}
			sr, err := q.QueryContext(ctx, `SELECT person_id, principal_cents FROM credit_purchase_share WHERE credit_purchase_id = ? ORDER BY person_id`, p.id)
			if err != nil {
				return mc.Input{}, err
			}
			for sr.Next() {
				var pid, pc int64
				if err := sr.Scan(&pid, &pc); err != nil {
					sr.Close()
					return mc.Input{}, err
				}
				pi.Shares = append(pi.Shares, mc.ShareInput{PersonID: pid, Code: persons[pid].code, External: persons[pid].ext, Principal: mc.Cents(pc)})
			}
			sr.Close()
			st.Purchases = append(st.Purchases, mc.PurchaseState{Ref: fmt.Sprintf("p%09d", p.id), Status: p.status, Input: pi})
		}
		in.Credits = append(in.Credits, st)
	}
	return in, nil
}

// Verify calcula con motor_cred, sobre la base, el disponible inicial A y B, el primer pago y la deuda viva, y los cota con
// oraculo/model.json (C3-09: 10,000.00 / 6,905.50, 973.75 del 30-oct-2026 y 6,293.44 al 10-abr-2027).
func Verify(ctx context.Context, q dominio.Querier, dir, modelPath string) (*Cruce, error) {
	m, err := readModel(dir, modelPath)
	if err != nil {
		return nil, err
	}
	if len(m.Credito.Eventos) == 0 || len(m.Despensas[0].Pagos) == 0 {
		return nil, fmt.Errorf("semilla: model.json no trae credito.eventos ni pagos")
	}
	start := mcDate(m.Credito.Eventos[0].Fecha)
	in, err := CreditInput(ctx, q, start)
	if err != nil {
		return nil, err
	}
	if len(in.Credits) != 1 {
		return nil, fmt.Errorf("semilla: se esperaba 1 crédito activo (Kueski) y hay %d", len(in.Credits))
	}
	st := in.Credits[0]
	last := start
	for _, p := range m.Despensas {
		for _, pay := range p.Pagos {
			if d := mcDate(pay.Fecha); d.Compare(last) > 0 {
				last = d
			}
		}
	}
	tl, err := mc.AvailabilityTimeline(mc.TimelineInput{State: st, Today: start, From: start, To: last})
	if err != nil {
		return nil, err
	}
	cut := mcDate(m.Horizonte.Fin)
	debt, err := mc.LiveDebt(in.Credits, cut)
	if err != nil {
		return nil, err
	}
	c := &Cruce{DeudaViva: int64(debt.Total), DeudaVivaFecha: cut.String(),
		Oraculo: Cruce1{A: int64(m.Credito.Eventos[0].DisponibleA), B: int64(m.Credito.Eventos[0].DisponibleB), Pago: int64(m.Despensas[0].Pagos[0].Monto), Deuda: int64(m.Horizonte.DeudaViva)}}
	if len(tl.Events) == 0 {
		return nil, fmt.Errorf("semilla: la línea de tiempo del disponible salió vacía")
	}
	c.DisponibleA, c.DisponibleB = int64(tl.Events[0].AvailableA), int64(tl.Events[0].AvailableB)
	for _, e := range tl.Events {
		if e.Kind == "payment" {
			c.PrimerPago, c.PrimerPagoFecha = int64(e.Amount), e.Date.String()
			break
		}
	}
	if c.DisponibleA != c.Oraculo.A || c.DisponibleB != c.Oraculo.B || c.PrimerPago != c.Oraculo.Pago || c.PrimerPagoFecha != m.Despensas[0].Pagos[0].Fecha || c.DeudaViva != c.Oraculo.Deuda {
		return c, fmt.Errorf("semilla: HALLAZGO: el motor sobre la base sembrada da A %d, B %d, pago %d (%s) y deuda viva %d; el oráculo dice A %d, B %d, pago %d (%s) y deuda %d",
			c.DisponibleA, c.DisponibleB, c.PrimerPago, c.PrimerPagoFecha, c.DeudaViva, c.Oraculo.A, c.Oraculo.B, c.Oraculo.Pago, m.Despensas[0].Pagos[0].Fecha, c.Oraculo.Deuda)
	}
	return c, nil
}

// cruce es el último paso de la carga: la prueba C3-09 sobre lo que acaba de quedar en la base.
func (l *loader) cruce() error {
	c, err := Verify(l.ctx, l.tx, l.opt.Dir, l.opt.ModelPath)
	if err != nil {
		return err
	}
	l.rep.Checks = append(l.rep.Checks,
		fmt.Sprintf("C3-09 sobre la base: disponible inicial A %s / B %s, pago del %s %s, deuda viva al %s %s = oráculo",
			dominio.FormatPesos(c.DisponibleA), dominio.FormatPesos(c.DisponibleB), c.PrimerPagoFecha, dominio.FormatPesos(c.PrimerPago), c.DeudaVivaFecha, dominio.FormatPesos(c.DeudaViva)))
	return nil
}
