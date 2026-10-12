package motor_cred

import (
	"math/big"
	"time"
)

// cardDueDates: corte = primera fecha >= compra con el día de corte (recortado); pago = primera
// fecha > corte con el día de pago (recortado); n cortes sucesivos. PROVISIONAL (S3-formula-tarjeta, pregunta 4).
func cardDueDates(c Credit, date Date, n int) []Date {
	sd, pd := *c.StatementDay, *c.PaymentDueDay
	st := clipDay(date.Year, date.Month, sd)
	if st.Compare(date) < 0 {
		st = clipDay(date.Year, date.Month+1, sd)
	}
	out := make([]Date, 0, n)
	for i := 0; i < n; i++ {
		s := clipDay(st.Year, st.Month+time.Month(i), sd)
		due := clipDay(s.Year, s.Month, pd)
		if due.Compare(s) <= 0 {
			due = clipDay(s.Year, s.Month+1, pd)
		}
		out = append(out, due)
	}
	return out
}

func validDay(p *int) bool { return p != nil && *p >= 1 && *p <= 31 }

func quoteCard(in PurchaseInput) (Purchase, error) {
	c := in.Credit
	if c.Limit == nil || !validDay(c.StatementDay) || !validDay(c.PaymentDueDay) || in.QuotedInterest != nil {
		return Purchase{}, invalid("tarjeta sin límite, corte o pago, o con intereses cotizados")
	}
	n := in.TermCount
	if in.PromotionID != nil {
		promo, err := chosenPromotion(in)
		if err != nil {
			return Purchase{}, err
		}
		if promo.Rate.Num != 0 {
			return Purchase{}, ErrUnsupported
		}
		amounts := SplitTrunc(in.Principal, n)
		p := Purchase{Rate: promo.Rate, RateBasis: "term_total", Source: "promotion", PurchaseKind: "installment_plan"}
		return finish(in, p, cardDueDates(c, in.Date, n), amounts, amounts, provList(provCard, serieProv(0, n)))
	}
	if n != 1 {
		return Purchase{}, ErrNoRate
	}
	amounts := []Cents{in.Principal}
	p := Purchase{Rate: Rate{0, 1, 0}, RateBasis: "", Source: "quoted", PurchaseKind: "single_charge"}
	return finish(in, p, cardDueDates(c, in.Date, 1), amounts, amounts, provList(provCard, provCardSource))
}

func missingK18(c Credit) bool {
	return c.AnnualRate == nil || c.RateBase == "" || c.MinPaymentPctBP == nil || c.MinPaymentFixed == nil
}

func checkK18(c Credit) error {
	if c.Kind != KindCard && c.Kind != KindStore {
		return invalid("el crédito %d no es tarjeta ni tienda", c.ID)
	}
	if missingK18(c) {
		return ErrMissingCapture
	}
	if !c.AnnualRate.valid() || (c.RateBase != "daily" && c.RateBase != "monthly") ||
		*c.MinPaymentPctBP < 0 || *c.MinPaymentFixed < 0 || (c.VATOnInterestBP != nil && *c.VATOnInterestBP < 0) {
		return invalid("datos K18 del crédito %d", c.ID)
	}
	return nil
}

// CardInterest devuelve el interés de un periodo y su lista provisional (§5.3, §14.1).
// Pago total a más tardar en la fecha de pago: 0.00 FIRME, con o sin datos capturados.
func (c Credit) CardInterest(st StatementInput, paid Cents, paidDate Date) (Cents, []string, error) {
	if c.Kind != KindCard && c.Kind != KindStore {
		return 0, nil, invalid("el crédito %d no es tarjeta ni tienda", c.ID)
	}
	if st.Balance < 0 {
		return 0, nil, invalid("saldo negativo")
	}
	if paid >= st.Balance && paidDate.Compare(st.DueDate) <= 0 {
		return 0, []string{}, nil
	}
	if err := checkK18(c); err != nil {
		return 0, nil, err
	}
	// Fórmula general PROVISIONAL (bloqueo S3-formula-tarjeta, opción b).
	var base Cents
	r := c.AnnualRate.Rat()
	if c.RateBase == "daily" {
		if len(st.DailyBalances) == 0 {
			return 0, nil, invalid("base diaria sin saldos diarios")
		}
		var sum Cents
		for _, b := range st.DailyBalances {
			sum += b
		}
		x := new(big.Rat).Mul(new(big.Rat).SetInt64(int64(sum)), r)
		base = roundRat(x.Quo(x, big.NewRat(360, 1)))
	} else {
		x := new(big.Rat).Mul(new(big.Rat).SetInt64(int64(st.Balance)), r)
		base = roundRat(x.Quo(x, big.NewRat(12, 1)))
	}
	if c.VATOnInterestBP != nil {
		base += RoundHalfUp(big.NewInt(int64(base)*int64(*c.VATOnInterestBP)), big.NewInt(10000))
	}
	return base, provList(provCard, provIVA), nil
}

// CardMinimum es CardMinimumPayment con su lista Provisional (todo mínimo usa la fórmula, §14.1).
func (c Credit) CardMinimum(st StatementInput) (Cents, []string, error) {
	m, err := c.CardMinimumPayment(st)
	if err != nil {
		return 0, nil, err
	}
	return m, provList(provCard, provIVA), nil
}

// CardPeriodInterest es CardInterest sin la lista provisional.
func (c Credit) CardPeriodInterest(st StatementInput, paidCents Cents, paidDate Date) (Cents, error) {
	v, _, err := c.CardInterest(st, paidCents, paidDate)
	return v, err
}

// CardMinimumPayment = máx(HALF_UP(saldo × pct / 10000), fijo) sin pasar del saldo. PROVISIONAL.
func (c Credit) CardMinimumPayment(st StatementInput) (Cents, error) {
	if err := checkK18(c); err != nil {
		return 0, err
	}
	if st.Balance < 0 {
		return 0, invalid("saldo negativo")
	}
	m := RoundHalfUp(big.NewInt(int64(st.Balance)*int64(*c.MinPaymentPctBP)), big.NewInt(10000))
	if f := *c.MinPaymentFixed; f > m {
		m = f
	}
	if m > st.Balance {
		m = st.Balance
	}
	return m, nil
}
