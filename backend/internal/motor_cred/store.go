package motor_cred

// quoteStore: crédito de tienda con abonos (§5.5, §5.8). Tasa 0 en la tabla = abonos sin interés
// (capital = monto); sin regla de intereses: falta capturar o, con datos K18, sin regla (PROVISIONAL).
func quoteStore(in PurchaseInput) (Purchase, error) {
	c := in.Credit
	if !validFrequency(c.Frequency) {
		return Purchase{}, invalid("tienda sin frecuencia")
	}
	if in.FirstPaymentDate == nil && (c.Frequency != "fortnightly" || c.FirstPaymentMinDays == nil || *c.FirstPaymentMinDays < 0) {
		return Purchase{}, invalid("tienda sin primera fecha")
	}
	r, ok, err := resolveInterest(in)
	if err != nil {
		return Purchase{}, err
	}
	if !ok {
		if missingK18(c) {
			return Purchase{}, ErrMissingCapture
		}
		return Purchase{}, ErrUnsupported
	}
	total := in.Principal + r.interest
	amounts := SplitTrunc(total, in.TermCount)
	principals := amounts
	if r.interest != 0 {
		principals = SplitHalfUp(in.Principal, in.TermCount)
	}
	p := Purchase{Interest: r.interest, TotalExact: r.exact, Rate: r.rate, RateBasis: "term_total",
		Source: r.source, PurchaseKind: "installment_plan"}
	return finish(in, p, installmentDates(in), amounts, principals, provList(serieProv(r.interest, in.TermCount)))
}
