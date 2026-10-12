package motor_cred

import "math/big"

// Amortize arma la tabla francesa (§5.4). Tasa del periodo = anual / 12, 24 o 52 y, con IVA,
// × (1 + IVA): PROVISIONAL (bloqueo S3-convencion-prestamo, opción a). El último pago ajusta.
func Amortize(principal Cents, annual Rate, vatBP *int, n int, frequency string, first Date) (Purchase, error) {
	per := map[string]int64{"monthly": 12, "fortnightly": 24, "weekly": 52}[frequency]
	if principal <= 0 || n < 1 || !annual.valid() || per == 0 || first.IsZero() || (vatBP != nil && *vatBP < 0) {
		return Purchase{}, invalid("préstamo: capital, plazo, tasa, frecuencia o fecha")
	}
	r := new(big.Rat).Quo(annual.Rat(), big.NewRat(per, 1))
	if vatBP != nil {
		r.Mul(r, big.NewRat(10000+int64(*vatBP), 10000))
	}
	c := new(big.Rat).SetInt64(int64(principal))
	var pay Cents
	if r.Sign() == 0 {
		pay = divHalfUp(principal, int64(n))
	} else {
		pow := big.NewRat(1, 1)
		one := new(big.Rat).Add(big.NewRat(1, 1), r)
		for i := 0; i < n; i++ {
			pow.Mul(pow, one)
		}
		x := new(big.Rat).Mul(c, r)
		x.Mul(x, pow)
		x.Quo(x, new(big.Rat).Sub(pow, big.NewRat(1, 1)))
		pay = roundRat(x)
	}
	dates := seriesFrom(first, n, frequency)
	p := Purchase{Principal: principal, ReferencePayment: pay, Rate: annual, RateDisplay: annual.Display(),
		RateBasis: "annual", Source: "amortization", PurchaseKind: "loan_disbursement",
		Shares: []PersonShare{}, Provisional: provList(provLoan)}
	if vatBP == nil { // «la tasa ya trae IVA»: lectura PROVISIONAL (bloqueo U4-iva-credito)
		p.Provisional = provList(provLoan, provIVA)
	}
	bal := principal
	for i := 0; i < n; i++ {
		in := roundRat(new(big.Rat).Mul(new(big.Rat).SetInt64(int64(bal)), r))
		amount, amort := pay, pay-in
		if i == n-1 || amort > bal {
			amort = bal
			amount = bal + in
		}
		bal -= amort
		p.Interest += in
		p.Payments = append(p.Payments, ScheduledPayment{Seq: i + 1, Due: dates[i], Amount: amount,
			Principal: amort, Interest: in, Shares: []PaymentShare{}})
		if bal == 0 && i < n-1 {
			return Purchase{}, ErrUnsupported
		}
	}
	p.Total = principal + p.Interest
	p.TotalExact = new(big.Rat).SetInt64(int64(p.Total))
	return p, nil
}

func quoteLoan(in PurchaseInput) (Purchase, error) {
	c := in.Credit
	if in.FirstPaymentDate == nil || in.QuotedInterest != nil || in.PromotionID != nil || !validFrequency(c.Frequency) {
		return Purchase{}, invalid("préstamo sin primera fecha o frecuencia, o con cotización o promoción")
	}
	if c.AnnualRate == nil {
		return Purchase{}, ErrMissingCapture
	}
	p, err := Amortize(in.Principal, *c.AnnualRate, c.VATOnInterestBP, in.TermCount, c.Frequency, *in.FirstPaymentDate)
	if err != nil {
		return Purchase{}, err
	}
	p.PromotionID = nil
	shares, err := applyShares(in.Shares, p.Principal, p.Interest, p.Payments)
	if err != nil {
		return Purchase{}, err
	}
	p.Shares = shares
	return p, nil
}
