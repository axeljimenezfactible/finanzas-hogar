package motor_cred

import "math/big"

// DebtLine es la deuda viva de una compra o de un crédito.
type DebtLine struct {
	CreditID     int64  `json:"credit_id"`
	PurchaseRef  string `json:"purchase_ref"`
	Pending      Cents  `json:"pending_cents"`
	PaymentCount int    `json:"payment_count"`
}

// DebtSummary es la deuda viva a una fecha (RN-21, K19).
type DebtSummary struct {
	AsOf       Date       `json:"as_of"`
	Total      Cents      `json:"total_cents"`
	ByPurchase []DebtLine `json:"by_purchase"`
	ByCredit   []DebtLine `json:"by_credit"`
}

// LiveDebt = Σ montos de pagos programados con vencimiento posterior a asOf de compras no canceladas.
func LiveDebt(credits []CreditState, asOf Date) (DebtSummary, error) {
	cs, err := sortedCredits(credits)
	if err != nil {
		return DebtSummary{}, err
	}
	out := DebtSummary{AsOf: asOf, ByPurchase: []DebtLine{}, ByCredit: []DebtLine{}}
	for _, st := range cs {
		qs, err := st.quoted()
		if err != nil {
			return DebtSummary{}, err
		}
		cl := DebtLine{CreditID: st.Credit.ID}
		for _, q := range qs {
			if !q.live() {
				continue
			}
			pl := DebtLine{CreditID: st.Credit.ID, PurchaseRef: q.Ref}
			for _, p := range q.Purchase.Payments {
				if p.Due.Compare(asOf) > 0 {
					pl.Pending += p.Amount
					pl.PaymentCount++
				}
			}
			if pl.Pending > 0 {
				out.ByPurchase = append(out.ByPurchase, pl)
				cl.Pending += pl.Pending
				cl.PaymentCount += pl.PaymentCount
			}
		}
		if cl.Pending > 0 {
			out.ByCredit = append(out.ByCredit, cl)
			out.Total += cl.Pending
		}
	}
	return out, nil
}

// LivePurchases: compras vivas en at, desde su fecha hasta su último pago inclusive (RN-11, RN-21).
func LivePurchases(state CreditState, at Date) ([]string, error) {
	qs, err := state.quoted()
	if err != nil {
		return nil, err
	}
	out := []string{}
	for _, q := range qs {
		ps := q.Purchase.Payments
		if q.live() && len(ps) > 0 && inRange(at, q.Input.Date, ps[len(ps)-1].Due) {
			out = append(out, q.Ref)
		}
	}
	return out, nil
}

// AmountQuote es la cotización por monto (interes_por_monto).
type AmountQuote struct {
	Amount           Cents `json:"amount_cents"`
	Interest         Cents `json:"interest_cents"`
	Total            Cents `json:"total_cents"`
	ReferencePayment Cents `json:"reference_payment_cents"`
}

// QuoteByAmount: intereses = ApplyRate; pago de referencia desde el total exacto (§1.2, §7).
func QuoteByAmount(amount Cents, r Rate, n int) AmountQuote {
	i := ApplyRate(amount, r)
	return AmountQuote{amount, i, amount + i, ReferencePayment(grossUp(amount, r), n)}
}

// SensitivityInput es el insumo de la tabla de sensibilidad.
type SensitivityInput struct {
	Amounts                []Cents    `json:"amounts_cents"`
	Rate                   Rate       `json:"rate"`
	TermCount              int        `json:"term_count"`
	External               ShareInput `json:"external"`
	HolderCurrentPrincipal Cents      `json:"holder_current_principal_cents"`
	AccrualMonths          int        `json:"accrual_months"`
	CashMonthly            Cents      `json:"cash_monthly_cents"`
}

// SensitivityRow es un renglón de la tabla de sensibilidad.
type SensitivityRow struct {
	Amount             Cents `json:"amount_cents"`
	Interest           Cents `json:"interest_cents"`
	Total              Cents `json:"total_cents"`
	ReferencePayment   Cents `json:"reference_payment_cents"`
	HolderTotal        Cents `json:"holder_total_cents"`
	HolderPrincipal    Cents `json:"holder_principal_cents"`
	HolderBelowCurrent bool  `json:"holder_below_current"`
	MonthlyPurchase    Cents `json:"monthly_purchase_cents"`
	MonthlyTotal       Cents `json:"monthly_total_cents"`
	HolderMonthly      Cents `json:"holder_monthly_cents"`
	ExternalMonthly    Cents `json:"external_monthly_cents"`
}

// Sensitivity arma la tabla de sensibilidad desde el total exacto t = N × (1 + tasa) (§7).
func Sensitivity(in SensitivityInput) ([]SensitivityRow, error) {
	if !in.Rate.valid() || in.TermCount < 1 || in.AccrualMonths < 1 || in.External.Principal < 0 {
		return nil, invalid("sensibilidad: tasa, plazo o meses")
	}
	kt := in.External.Principal + ApplyRate(in.External.Principal, in.Rate)
	months := big.NewRat(int64(in.AccrualMonths), 1)
	rows := make([]SensitivityRow, 0, len(in.Amounts))
	for _, n := range in.Amounts {
		t := grossUp(n, in.Rate)
		i := ApplyRate(n, in.Rate)
		row := SensitivityRow{Amount: n, Interest: i, Total: n + i, ReferencePayment: ReferencePayment(t, in.TermCount)}
		row.HolderTotal = row.Total - kt
		row.HolderPrincipal = n - in.External.Principal
		row.HolderBelowCurrent = row.HolderPrincipal < in.HolderCurrentPrincipal
		row.MonthlyPurchase = roundRat(new(big.Rat).Quo(t, months))
		row.MonthlyTotal = row.MonthlyPurchase + in.CashMonthly
		h := new(big.Rat).Sub(t, new(big.Rat).SetInt64(int64(kt)))
		row.HolderMonthly = roundRat(h.Quo(h, months))
		row.ExternalMonthly = divHalfUp(kt, int64(in.AccrualMonths)) + in.CashMonthly
		rows = append(rows, row)
	}
	return rows, nil
}
