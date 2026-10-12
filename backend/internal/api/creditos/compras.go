package creditos

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// checkShares valida las partes por persona (capital) y las vuelve entrada del motor; vacío = todo al titular.
func checkShares(v *api.Validator, s *snapshot, total int64, in []api.CreditPurchaseShareInput) []mc.ShareInput {
	if len(in) == 0 {
		return nil
	}
	var sum int64
	var out []mc.ShareInput
	for _, sh := range in {
		_, ok := s.persons[sh.PersonID]
		v.Check(ok, "shares", "not_found", "Una de las personas de las partes no existe.")
		v.Check(sh.PrincipalCents >= 0, "shares", "out_of_range", "El capital de cada parte no puede ser negativo.")
		sum += sh.PrincipalCents
		out = append(out, s.shareInput(sh.PersonID, sh.PrincipalCents))
	}
	v.Check(sum == total, "shares", "invalid_sum", "Las partes deben sumar el monto de la compra.")
	return out
}

// purchaseInput valida la captura y arma la entrada del motor; el motor calcula intereses, partes y serie.
func purchaseInput(s *snapshot, in api.CreditPurchaseCreate) (mc.PurchaseInput, dominio.CreditRow, error) {
	var v api.Validator
	c, ok := s.byID[in.CreditID]
	v.Check(ok, "credit_id", "not_found", "No encontramos ese crédito.")
	v.Check(oneOf(in.Kind, "installment_plan", "single_charge", "loan_disbursement"), "kind", "invalid_format", "Elige un tipo de compra válido.")
	v.PositiveCents(in.PrincipalCents, "principal_cents")
	v.Check(in.TermCount >= 1 && in.TermCount <= 360, "term_count", "out_of_range", "El número de pagos va de 1 a 360.")
	v.Check(in.Kind != "single_charge" || in.TermCount == 1, "term_count", "out_of_range", "Un cargo normal es de un solo pago.")
	v.Check(in.InterestCents == nil || *in.InterestCents >= 0, "interest_cents", "out_of_range", "Los intereses no pueden ser negativos.")
	v.Check(in.PaymentFrequency == nil || !ok || *in.PaymentFrequency == c.PaymentFrequency, "payment_frequency", "not_allowed",
		"La frecuencia de los pagos es la del crédito.")
	// PROVISIONAL (bloqueo S6-creditos-gasto-vinculado): el alta del gasto vinculado es de api-gastos-inicio.
	v.Check(in.Expense == nil, "expense", "not_allowed", "Por ahora registra primero el gasto y manda su expense_id.")
	date, err := v.Date(in.PurchaseDate, "purchase_date")
	if err != nil {
		return mc.PurchaseInput{}, c, err
	}
	var first *mc.Date
	if in.FirstPaymentDate != nil {
		d, err := v.Date(*in.FirstPaymentDate, "first_payment_date")
		if err != nil {
			return mc.PurchaseInput{}, c, err
		}
		first = api.Ptr(mdateOf(d))
	}
	shares := checkShares(&v, s, in.PrincipalCents, in.Shares)
	if err := v.Err(); err != nil {
		return mc.PurchaseInput{}, c, err
	}
	if shares == nil {
		shares = []mc.ShareInput{s.shareInput(c.HolderPersonID, in.PrincipalCents)}
	}
	pin := mc.PurchaseInput{Credit: s.motorCredit(c), Date: mdateOf(date), Principal: mc.Cents(in.PrincipalCents),
		TermCount: int(in.TermCount), QuotedInterest: centsPtr(in.InterestCents), RateDecimals: intPtr(in.RateDecimals),
		PromotionID: in.PromotionID, FirstPaymentDate: first, Shares: shares}
	for _, p := range s.promos { // la promoción elegida se compara con su propio comercio (supuesto en el reporte)
		if in.PromotionID != nil && p.ID == *in.PromotionID {
			pin.Merchant, pin.StoreID = p.MerchantName, p.StoreID
		}
	}
	return pin, c, nil
}

func previewCreditPurchase(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.CreditPurchaseCreate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		pin, _, err := purchaseInput(s, in)
		if err != nil {
			return err
		}
		out := api.CreditPurchasePreview{}
		p, err := mc.QuotePurchase(pin)
		if errors.Is(err, mc.ErrPromotionNotApplicable) {
			out.PromotionIssue = api.Ptr("La promoción no aplica a esta compra (comercio, compra mínima, plazo o vigencia); el cálculo va sin ella.")
			pin.PromotionID, pin.Merchant, pin.StoreID = nil, "", nil
			p, err = mc.QuotePurchase(pin)
		}
		if err != nil {
			return motorErr(err, "credit_id")
		}
		if in.PromotionID != nil {
			out.PromotionApplies = api.Ptr(out.PromotionIssue == nil)
		}
		out.PrincipalCents, out.InterestCents, out.TotalCents = int64(p.Principal), int64(p.Interest), int64(p.Total)
		out.Shares, out.Schedule = []api.CreditPurchasePreviewSharesItem{}, []api.CreditPurchasePreviewScheduleItem{}
		for _, sh := range p.Shares {
			out.Shares = append(out.Shares, api.CreditPurchasePreviewSharesItem{PersonID: sh.PersonID, PrincipalCents: int64(sh.Principal),
				InterestCents: int64(sh.Interest), TotalCents: int64(sh.Total), FortnightPaymentCents: int64(sh.FirstPayment)})
		}
		for _, sp := range p.Payments {
			it := api.CreditPurchasePreviewScheduleItem{Seq: int64(sp.Seq), DueDate: sp.Due.String(), AmountCents: int64(sp.Amount),
				PrincipalCents: int64(sp.Principal), InterestCents: int64(sp.Interest), Shares: []api.ShareLine{}}
			for _, sh := range sp.Shares {
				it.Shares = append(it.Shares, api.ShareLine{PersonID: sh.PersonID, AmountCents: int64(sh.Amount)})
			}
			out.Schedule = append(out.Schedule, it)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func createCreditPurchase(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.CreditPurchaseCreate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		ctx := r.Context()
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		pin, c, err := purchaseInput(s, in)
		if err != nil {
			return err
		}
		p, err := mc.QuotePurchase(pin)
		if err != nil {
			return motorErr(err, "credit_id")
		}
		var v api.Validator
		v.Check(in.InterestSource == nil || *in.InterestSource == p.Source, "interest_source", "not_allowed",
			"El origen de los intereses no coincide con el cálculo ("+p.Source+").")
		if in.ExpenseID != nil {
			ok, err := dominio.Exists(ctx, reg.Deps().DB, "expense", *in.ExpenseID)
			if err != nil {
				return err
			}
			v.Check(ok, "expense_id", "not_found", "No encontramos ese gasto.")
		}
		if err := v.Err(); err != nil {
			return err
		}
		row := dominio.PurchaseRow{CreditID: c.ID, ExpenseID: in.ExpenseID, Kind: in.Kind, Description: in.Description,
			PurchaseDate: in.PurchaseDate, PrincipalCents: int64(p.Principal), InterestCents: int64(p.Interest), TotalCents: int64(p.Total),
			TermCount: in.TermCount, PaymentFrequency: c.PaymentFrequency, FirstPaymentDate: p.Payments[0].Due.String(),
			InterestSource: p.Source, InterestIncludesVat: deref(in.InterestIncludesVat, true), PromotionID: p.PromotionID,
			Status: "active", Notes: in.Notes, RateDecimals: in.RateDecimals}
		if p.RateBasis != "" && p.Rate.Den > 0 {
			row.RateNum, row.RateDen, row.RateBasis = api.Ptr(p.Rate.Num), api.Ptr(p.Rate.Den), api.Ptr(p.RateBasis)
		} else {
			row.RateDecimals = nil
		}
		var id int64
		err = audited(r, reg, "insert", "credit_purchase", &id, func(tx *sql.Tx) error {
			if id, err = dominio.InsertPurchase(ctx, tx, row); err != nil {
				return dbErr(err, "expense_id", "una compra ligada a ese gasto")
			}
			for _, sh := range p.Shares {
				if _, err := dominio.InsertPurchaseShare(ctx, tx, dominio.PurchaseShareRow{CreditPurchaseID: id, PersonID: sh.PersonID,
					PrincipalCents: int64(sh.Principal), InterestCents: int64(sh.Interest), TotalCents: int64(sh.Total)}); err != nil {
					return err
				}
			}
			for _, sp := range p.Payments {
				spID, err := dominio.InsertScheduled(ctx, tx, dominio.ScheduledRow{CreditPurchaseID: id, Seq: int64(sp.Seq),
					DueDate: sp.Due.String(), AmountCents: int64(sp.Amount), PrincipalCents: int64(sp.Principal), InterestCents: int64(sp.Interest)})
				if err != nil {
					return err
				}
				for _, sh := range sp.Shares {
					if err := dominio.InsertScheduledShare(ctx, tx, spID, sh.PersonID, int64(sh.Amount)); err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return writePurchase(reg, w, r, http.StatusCreated, id)
	}
}

func (s *snapshot) purchaseView(p dominio.PurchaseRow) api.CreditPurchase {
	out := api.CreditPurchase{ID: p.ID, CreditID: p.CreditID, ExpenseID: p.ExpenseID, Kind: p.Kind, Description: p.Description,
		PurchaseDate: p.PurchaseDate, PrincipalCents: p.PrincipalCents, InterestCents: p.InterestCents, TotalCents: p.TotalCents,
		TermCount: p.TermCount, PaymentFrequency: p.PaymentFrequency, FirstPaymentDate: p.FirstPaymentDate, InterestSource: p.InterestSource,
		RateNum: p.RateNum, RateDen: p.RateDen, RateBasis: p.RateBasis, InterestIncludesVat: p.InterestIncludesVat, PromotionID: p.PromotionID,
		Status: p.Status, Notes: p.Notes, CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt), RateDecimals: p.RateDecimals,
		Shares: []api.CreditPurchaseShare{}}
	if p.RateNum != nil && p.RateDen != nil {
		out.RatePct = ratePct(*p.RateNum, *p.RateDen, p.RateDecimals)
	}
	for _, sh := range s.shares[p.ID] {
		out.Shares = append(out.Shares, api.CreditPurchaseShare{ID: sh.ID, CreditPurchaseID: p.ID, PersonID: sh.PersonID,
			PrincipalCents: sh.PrincipalCents, InterestCents: sh.InterestCents, TotalCents: sh.TotalCents,
			CreatedAt: api.APITime(sh.CreatedAt), UpdatedAt: api.APITime(sh.UpdatedAt)})
	}
	var made, total, balance int64
	for _, sp := range s.sched {
		if sp.CreditPurchaseID != p.ID {
			continue
		}
		total++
		if rest := sp.AmountCents - s.paidBySP[sp.ID]; rest > 0 {
			balance += rest
		} else {
			made++
		}
		if sp.Seq == 1 {
			for i := range out.Shares {
				out.Shares[i].FortnightPaymentCents = api.Ptr(s.shareAmount(sp.ID, out.Shares[i].PersonID))
			}
		}
	}
	out.PaymentsMade, out.PaymentsTotal, out.BalanceCents = api.Ptr(made), api.Ptr(total), api.Ptr(balance)
	return out
}

func (s *snapshot) shareAmount(spID, personID int64) int64 {
	for _, x := range s.spShares[spID] {
		if x.PersonID == personID {
			return x.AmountCents
		}
	}
	return 0
}

func writePurchase(reg *api.Registry, w http.ResponseWriter, r *http.Request, status int, id int64) error {
	s, err := snap(reg, r)
	if err != nil {
		return err
	}
	p, ok := s.purchase(id)
	if !ok {
		return dominio.NotFound(msgPurchase)
	}
	return api.WriteJSON(w, status, s.purchaseView(p))
}

const msgPurchase = "No encontramos esa compra a crédito."

func (s *snapshot) purchase(id int64) (dominio.PurchaseRow, bool) {
	for _, p := range s.purchases {
		if p.ID == id {
			return p, true
		}
	}
	return dominio.PurchaseRow{}, false
}

func inIDs(col string, ids []int64, where []string, args []any) ([]string, []any) {
	if len(ids) == 0 {
		return where, args
	}
	where = append(where, col+" IN ("+strings.TrimSuffix(strings.Repeat("?, ", len(ids)), ", ")+")")
	for _, id := range ids {
		args = append(args, id)
	}
	return where, args
}

// familySQL: household = alguna parte del hogar; external = alguna parte de Karen; joint = ambas (K5).
func familySQL(r *http.Request, shareExists string) (string, error) {
	f, err := badEnum(r, "family", "household", "external", "joint")
	if err != nil || f == "" {
		return "", err
	}
	hh := strings.ReplaceAll(shareExists, "%COND%", "pe.participation <> 'external'")
	ex := strings.ReplaceAll(shareExists, "%COND%", "pe.participation = 'external'")
	switch f {
	case "household":
		return hh, nil
	case "external":
		return ex, nil
	}
	return hh + " AND " + ex, nil
}

func listCreditPurchases(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		pg, err := api.ParsePage(r, []string{"purchase_date", "total_cents"}, "-purchase_date")
		if err != nil {
			return err
		}
		status, err := badEnum(r, "status", "active", "paid_off", "cancelled")
		if err != nil {
			return err
		}
		where, args := []string{"p.status = ?"}, []any{deref(api.Ptr(status), "active")}
		if status == "" {
			args[0] = "active"
		}
		credits, err := api.QueryIDs(r, "credit_id")
		if err != nil {
			return err
		}
		where, args = inIDs("p.credit_id", credits, where, args)
		persons, err := api.QueryIDs(r, "person_id")
		if err != nil {
			return err
		}
		if len(persons) > 0 {
			w2, a2 := inIDs("s.person_id", persons, nil, nil)
			where, args = append(where, "EXISTS (SELECT 1 FROM credit_purchase_share s WHERE s.credit_purchase_id = p.id AND "+w2[0]+")"), append(args, a2...)
		}
		fam, err := familySQL(r, "EXISTS (SELECT 1 FROM credit_purchase_share s JOIN person pe ON pe.id = s.person_id WHERE s.credit_purchase_id = p.id AND %COND%)")
		if err != nil {
			return err
		}
		if fam != "" {
			where = append(where, fam)
		}
		if per, err := api.QueryPeriod(r); err != nil {
			return err
		} else if per != nil {
			where, args = append(where, "p.purchase_date BETWEEN ? AND ?"), append(args, per.Start.String(), per.End.String())
		}
		rows, err := dominio.QueryPurchases(r.Context(), reg.Deps().DB, strings.Join(where, " AND "),
			pg.OrderBy(map[string]string{"purchase_date": "p.purchase_date", "total_cents": "p.total_cents", "id": "p.id"})+" "+pg.LimitOffset(), args...)
		if err != nil {
			return err
		}
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		items := make([]api.CreditPurchase, 0, len(rows))
		for _, p := range rows {
			items = append(items, s.purchaseView(p))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, pg))
	}
}

func purchaseFor(reg *api.Registry, r *http.Request) (*snapshot, dominio.PurchaseRow, error) {
	id, err := api.PathID(r, "credit_purchase_id")
	if err != nil {
		return nil, dominio.PurchaseRow{}, err
	}
	s, err := snap(reg, r)
	if err != nil {
		return nil, dominio.PurchaseRow{}, err
	}
	p, ok := s.purchase(id)
	if !ok {
		return nil, p, dominio.NotFound(msgPurchase)
	}
	return s, p, nil
}

func getCreditPurchase(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		s, p, err := purchaseFor(reg, r)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, s.purchaseView(p))
	}
}

func updateCreditPurchase(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.CreditPurchaseUpdate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		_, p, err := purchaseFor(reg, r)
		if err != nil {
			return err
		}
		nullable(&p.Description, in.Description)
		nullable(&p.Notes, in.Notes)
		set(&p.Status, in.Status)
		var v api.Validator
		v.Check(oneOf(p.Status, "active", "paid_off", "cancelled"), "status", "invalid_format", "El estado es activa, liquidada o cancelada.")
		if err := v.Err(); err != nil {
			return err
		}
		if err := audited(r, reg, "update", "credit_purchase", &p.ID, func(tx *sql.Tx) error {
			return dbErr(dominio.UpdatePurchaseText(r.Context(), tx, p), "status", "la compra")
		}); err != nil {
			return err
		}
		return writePurchase(reg, w, r, http.StatusOK, p.ID)
	}
}

func deleteCreditPurchase(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		_, p, err := purchaseFor(reg, r)
		if err != nil {
			return err
		}
		n, err := dominio.CountWhere(r.Context(), reg.Deps().DB, "credit_payment",
			"scheduled_payment_id IN (SELECT id FROM scheduled_payment WHERE credit_purchase_id = ?)", p.ID)
		if err != nil {
			return err
		}
		if n > 0 {
			return dominio.Conflict("has_payments", "Esta compra ya tiene pagos registrados. Desmárcalos antes o cancélala.")
		}
		if err := audited(r, reg, "delete", "credit_purchase", &p.ID, func(tx *sql.Tx) error {
			return dominio.DeleteRow(r.Context(), tx, "credit_purchase", p.ID)
		}); err != nil {
			return err
		}
		return api.NoContent(w)
	}
}
