package creditos

import (
	"database/sql"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

const paidSQL = `(SELECT COALESCE(SUM(m.amount_cents), 0) FROM credit_payment m WHERE m.scheduled_payment_id = sp.id)`

func (s *snapshot) scheduledView(sp dominio.ScheduledRow) api.ScheduledPayment {
	t := s.today.String()
	out := api.ScheduledPayment{ID: sp.ID, CreditPurchaseID: sp.CreditPurchaseID, Seq: sp.Seq, DueDate: sp.DueDate, AmountCents: sp.AmountCents,
		PrincipalCents: sp.PrincipalCents, InterestCents: sp.InterestCents, CreatedAt: api.APITime(sp.CreatedAt), UpdatedAt: api.APITime(sp.UpdatedAt),
		CreditID: api.Ptr(sp.CreditID), Shares: []api.ScheduledPaymentSharesItem{}}
	state := func(paid, amount int64) string {
		switch {
		case paid >= amount:
			return "paid"
		case sp.DueDate < t:
			return "overdue"
		}
		return "pending"
	}
	for _, sh := range s.spShares[sp.ID] {
		out.Shares = append(out.Shares, api.ScheduledPaymentSharesItem{PersonID: sh.PersonID, AmountCents: sh.AmountCents,
			PaidCents: sh.PaidCents, Status: state(sh.PaidCents, sh.AmountCents)})
	}
	out.Status = api.Ptr(state(s.paidBySP[sp.ID], sp.AmountCents))
	if p, ok := s.purchase(sp.CreditPurchaseID); ok { // saldo de la compra después de este pago (serie guardada del motor)
		bal := p.TotalCents
		for _, x := range s.sched {
			if x.CreditPurchaseID == p.ID && x.Seq <= sp.Seq {
				bal -= x.AmountCents
			}
		}
		out.BalanceAfterCents = api.Ptr(bal)
	}
	return out
}

func listScheduledPayments(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		pg, err := api.ParsePage(r, []string{"due_date", "amount_cents"}, "due_date")
		if err != nil {
			return err
		}
		status, err := badEnum(r, "status", "pending", "paid", "overdue")
		if err != nil {
			return err
		}
		t := today(reg).String()
		where, args := []string{"p.status <> 'cancelled'"}, []any{}
		for _, f := range []struct{ param, col string }{{"credit_id", "p.credit_id"}, {"credit_purchase_id", "sp.credit_purchase_id"}} {
			ids, err := api.QueryIDs(r, f.param)
			if err != nil {
				return err
			}
			where, args = inIDs(f.col, ids, where, args)
		}
		persons, err := api.QueryIDs(r, "person_id")
		if err != nil {
			return err
		}
		if len(persons) > 0 {
			w2, a2 := inIDs("x.person_id", persons, nil, nil)
			where, args = append(where, "EXISTS (SELECT 1 FROM scheduled_payment_share x WHERE x.scheduled_payment_id = sp.id AND "+w2[0]+")"), append(args, a2...)
		}
		fam, err := familySQL(r, "EXISTS (SELECT 1 FROM scheduled_payment_share x JOIN person pe ON pe.id = x.person_id WHERE x.scheduled_payment_id = sp.id AND %COND%)")
		if err != nil {
			return err
		}
		if fam != "" {
			where = append(where, fam)
		}
		per, err := api.QueryPeriod(r)
		if err != nil {
			return err
		}
		if per != nil {
			where, args = append(where, "sp.due_date BETWEEN ? AND ?"), append(args, per.Start.String(), per.End.String())
		}
		switch status {
		case "":
			where = append(where, paidSQL+" < sp.amount_cents") // por defecto: con alguna parte pendiente, de hoy en adelante
			if per == nil {
				where, args = append(where, "sp.due_date >= ?"), append(args, t)
			}
		case "pending":
			where = append(where, paidSQL+" < sp.amount_cents")
		case "paid":
			where = append(where, paidSQL+" >= sp.amount_cents")
		case "overdue":
			where, args = append(where, paidSQL+" < sp.amount_cents AND sp.due_date < ?"), append(args, t)
		}
		rows, err := dominio.QueryScheduled(r.Context(), reg.Deps().DB, strings.Join(where, " AND "),
			pg.OrderBy(map[string]string{"due_date": "sp.due_date", "amount_cents": "sp.amount_cents", "id": "sp.id"})+" "+pg.LimitOffset(), args...)
		if err != nil {
			return err
		}
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		items := make([]api.ScheduledPayment, 0, len(rows))
		for _, sp := range rows {
			items = append(items, s.scheduledView(sp))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
	}
}

func paymentView(m dominio.PaymentRow) api.CreditPayment {
	return api.CreditPayment{ID: m.ID, CreditID: m.CreditID, ScheduledPaymentID: m.ScheduledPaymentID, CreditStatementID: m.CreditStatementID,
		PayerPersonID: m.PayerPersonID, PaidDate: m.PaidDate, AmountCents: m.AmountCents, Notes: m.Notes,
		CreatedAt: api.APITime(m.CreatedAt), UpdatedAt: api.APITime(m.UpdatedAt)}
}

// syncPurchaseStatus: active = con pagos pendientes, paid_off = liquidada (C1 CreditPurchase.status).
func syncPurchaseStatus(r *http.Request, tx *sql.Tx, purchaseID int64) error {
	ctx := r.Context()
	p, err := dominio.GetPurchase(ctx, tx, purchaseID)
	if err != nil || p.Status == "cancelled" {
		return err
	}
	pending, err := dominio.CountWhere(ctx, tx, "scheduled_payment sp", "sp.credit_purchase_id = ? AND "+paidSQL+" < sp.amount_cents", purchaseID)
	if err != nil {
		return err
	}
	want := "active"
	if pending == 0 {
		want = "paid_off"
	}
	if want == p.Status {
		return nil
	}
	before, _ := dominio.RowJSON(ctx, tx, "credit_purchase", p.ID)
	p.Status = want
	if err := dominio.UpdatePurchaseText(ctx, tx, p); err != nil {
		return err
	}
	after, _ := dominio.RowJSON(ctx, tx, "credit_purchase", p.ID)
	return dominio.Audit(ctx, tx, api.Actor(r), "update", "credit_purchase", p.ID, before, after, api.RequestID(r))
}

func markScheduledPaymentPaid(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "scheduled_payment_id")
		if err != nil {
			return err
		}
		var in api.CreditPaymentCreate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		sps, err := dominio.QueryScheduled(ctx, db, "sp.id = ?", "", id)
		if err != nil {
			return err
		}
		if len(sps) == 0 {
			return dominio.NotFound("No encontramos ese pago programado.")
		}
		sp := sps[0]
		p, err := dominio.GetPurchase(ctx, db, sp.CreditPurchaseID)
		if err != nil {
			return err
		}
		if p.Status == "cancelled" {
			return dominio.Conflict("conflict", "La compra está cancelada; sus pagos ya no se marcan.")
		}
		var v api.Validator
		v.PositiveCents(in.AmountCents, "amount_cents")
		if _, err := v.Date(in.PaidDate, "paid_date"); err != nil {
			return err
		}
		shares, err := dominio.ScheduledShares(ctx, db, sp.CreditPurchaseID)
		if err != nil {
			return err
		}
		var share *dominio.ScheduledShareRow
		for i, sh := range shares[sp.ID] {
			if sh.PersonID == in.PayerPersonID {
				share = &shares[sp.ID][i]
			}
		}
		v.Check(share != nil, "payer_person_id", "not_allowed", "Esa persona no tiene parte en este pago.")
		if err := v.Err(); err != nil {
			return err
		}
		rest := share.AmountCents - share.PaidCents
		if rest <= 0 {
			return dominio.Conflict("conflict", "Esa parte ya está pagada. Si fue un error, desmarca el pago.")
		}
		if in.AmountCents > rest {
			return dominio.Invalid(dominio.FieldIssue{Field: "amount_cents", Code: "out_of_range",
				Message: "El monto es mayor a lo que falta de esta parte (" + dominio.FormatPesos(rest) + ")."})
		}
		m := dominio.PaymentRow{CreditID: sp.CreditID, ScheduledPaymentID: &sp.ID, PayerPersonID: in.PayerPersonID, PaidDate: in.PaidDate,
			AmountCents: in.AmountCents, Notes: in.Notes}
		var pid int64
		if err := audited(r, reg, "insert", "credit_payment", &pid, func(tx *sql.Tx) (err error) {
			if pid, err = dominio.InsertPayment(ctx, tx, m); err != nil {
				return dbErr(err, "amount_cents", "el pago")
			}
			return nil
		}); err != nil {
			return err
		}
		if err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error { return syncPurchaseStatus(r, tx, p.ID) }); err != nil {
			return err
		}
		ms, err := dominio.QueryPayments(ctx, db, "m.id = ?", pid)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, paymentView(ms[0]))
	}
}

func deleteCreditPayment(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "credit_payment_id")
		if err != nil {
			return err
		}
		ctx := r.Context()
		ms, err := dominio.QueryPayments(ctx, reg.Deps().DB, "m.id = ?", id)
		if err != nil {
			return err
		}
		if len(ms) == 0 {
			return dominio.NotFound("No encontramos ese pago.")
		}
		if err := audited(r, reg, "delete", "credit_payment", &id, func(tx *sql.Tx) error {
			if err := dominio.DeleteRow(ctx, tx, "credit_payment", id); err != nil {
				return err
			}
			if ms[0].PurchaseID != 0 {
				return syncPurchaseStatus(r, tx, ms[0].PurchaseID)
			}
			return nil
		}); err != nil {
			return err
		}
		return api.NoContent(w)
	}
}
