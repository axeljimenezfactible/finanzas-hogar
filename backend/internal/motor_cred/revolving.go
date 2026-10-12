package motor_cred

import (
	"math/big"
	"sort"
	"strings"
)

// Listas de bloqueos (§14.1), en orden alfabético por bytes.
const (
	provLoan       = "S3-convencion-prestamo"
	provCard       = "S3-formula-tarjeta"
	provCardSource = "S3-origen-cargo-normal"
	provSerie      = "S3-serie-sin-intereses"
	provIVA        = "U4-iva-credito"
)

// provList arma una lista Provisional sin repetidos y en orden alfabético por bytes (§14.1).
func provList(names ...string) []string {
	out := []string{}
	for _, n := range names {
		if n == "" {
			continue
		}
		i := sort.SearchStrings(out, n)
		if i < len(out) && out[i] == n {
			continue
		}
		out = append(out, "")
		copy(out[i+1:], out[i:])
		out[i] = n
	}
	return out
}

// serieProv: con intereses 0 y n >= 2 la serie de capital es PROVISIONAL (bloqueo S3-serie-sin-intereses).
func serieProv(interest Cents, n int) string {
	if interest == 0 && n >= 2 {
		return provSerie
	}
	return ""
}

func validFrequency(f string) bool { return f == "weekly" || f == "fortnightly" || f == "monthly" }

// QuotePurchase despacha por in.Credit.Kind (§5.8).
func QuotePurchase(in PurchaseInput) (Purchase, error) {
	if err := validateCommon(in); err != nil {
		return Purchase{}, err
	}
	switch in.Credit.Kind {
	case KindRevolving:
		return quoteRevolving(in)
	case KindCard:
		return quoteCard(in)
	case KindLoan:
		return quoteLoan(in)
	case KindStore:
		return quoteStore(in)
	}
	return Purchase{}, invalid("tipo de crédito %q", in.Credit.Kind)
}

func validateCommon(in PurchaseInput) error {
	if in.Principal <= 0 || in.TermCount < 1 || in.Date.IsZero() {
		return invalid("capital, plazo o fecha")
	}
	if in.FirstPaymentDate != nil && in.FirstPaymentDate.Compare(in.Date) < 0 {
		return invalid("primer pago antes de la compra")
	}
	if in.QuotedInterest != nil && *in.QuotedInterest < 0 {
		return invalid("intereses cotizados negativos")
	}
	if in.RateDecimals != nil && *in.RateDecimals < 0 {
		return invalid("rate_decimals negativo")
	}
	holders := 0
	var sum Cents
	codes := map[string]bool{}
	for _, s := range in.Shares {
		if s.Code == "" || codes[s.Code] || s.Principal < 0 {
			return invalid("parte %q", s.Code)
		}
		codes[s.Code] = true
		if !s.External {
			holders++
		}
		sum += s.Principal
	}
	if holders != 1 || sum != in.Principal {
		return invalid("partes: %d titulares, suma %d de %d", holders, sum, in.Principal)
	}
	return nil
}

func promoMatches(p Promotion, c Credit, date Date, principal Cents, term int, storeID *int64, merchant string) bool {
	if !p.Active || (p.CreditID != nil && *p.CreditID != c.ID) {
		return false
	}
	if p.StoreID != nil && storeID != nil {
		if *p.StoreID != *storeID {
			return false
		}
	} else if !sameMerchant(p.Merchant, merchant) {
		return false
	}
	return p.ValidFrom.Compare(date) <= 0 && date.Compare(p.ValidTo) <= 0 &&
		principal >= p.MinPurchase && term == p.TermCount
}

// sameMerchant compara sin espacios a los lados y sin mayúsculas ASCII (COLLATE NOCASE).
func sameMerchant(a, b string) bool {
	low := func(s string) string {
		s = strings.TrimSpace(s)
		bs := []byte(s)
		for i, ch := range bs {
			if ch >= 'A' && ch <= 'Z' {
				bs[i] = ch + 32
			}
		}
		return string(bs)
	}
	return low(a) == low(b)
}

// ApplicablePromotions devuelve, en orden de ID, las promociones capturadas que aplican (RN-10).
func ApplicablePromotions(c Credit, date Date, principal Cents, term int, storeID *int64, merchant string) []Promotion {
	out := []Promotion{}
	for _, p := range c.Promotions {
		if promoMatches(p, c, date, principal, term, storeID, merchant) {
			out = append(out, p)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func chosenPromotion(in PurchaseInput) (Promotion, error) {
	for _, p := range in.Credit.Promotions {
		if p.ID == *in.PromotionID {
			if !p.Rate.valid() {
				return Promotion{}, invalid("tasa de la promoción %d", p.ID)
			}
			if promoMatches(p, in.Credit, in.Date, in.Principal, in.TermCount, in.StoreID, in.Merchant) {
				return p, nil
			}
		}
	}
	return Promotion{}, ErrPromotionNotApplicable
}

func termRate(c Credit, n int) (TermRate, bool) {
	for _, t := range c.TermRates {
		if t.TermCount == n {
			return t, true
		}
	}
	return TermRate{}, false
}

type interestRule struct {
	interest Cents
	rate     Rate
	source   string
	exact    *big.Rat
}

// resolveInterest aplica (1) cotización, (2) promoción y (3) tabla por plazo (RN-08, RN-10).
func resolveInterest(in PurchaseInput) (interestRule, bool, error) {
	n := in.TermCount
	if in.QuotedInterest != nil {
		dec := -1 // los decimales de un TermRate del mismo plazo no se usan (K12, §5.1)
		if in.RateDecimals != nil {
			dec = *in.RateDecimals
		}
		q := *in.QuotedInterest
		ex := new(big.Rat).SetInt64(int64(in.Principal + q))
		return interestRule{q, QuotedRate(q, in.Principal, dec), "quoted", ex}, true, nil
	}
	if in.PromotionID != nil {
		p, err := chosenPromotion(in)
		if err != nil {
			return interestRule{}, false, err
		}
		return interestRule{ApplyRate(in.Principal, p.Rate), p.Rate, "promotion", grossUp(in.Principal, p.Rate)}, true, nil
	}
	if t, ok := termRate(in.Credit, n); ok {
		if !t.Rate.valid() {
			return interestRule{}, false, invalid("tasa del plazo %d", n)
		}
		return interestRule{ApplyRate(in.Principal, t.Rate), t.Rate, "term_table", grossUp(in.Principal, t.Rate)}, true, nil
	}
	return interestRule{}, false, nil
}

// installmentDates: primera fecha capturada o, en quincenal, la regla de RN-06.
func installmentDates(in PurchaseInput) []Date {
	if in.FirstPaymentDate != nil {
		return seriesFrom(*in.FirstPaymentDate, in.TermCount, in.Credit.Frequency)
	}
	return fortnightSchedule(in.Date, *in.Credit.FirstPaymentMinDays, in.TermCount)
}

func quoteRevolving(in PurchaseInput) (Purchase, error) {
	c := in.Credit
	if c.Limit == nil || !validFrequency(c.Frequency) {
		return Purchase{}, invalid("línea sin límite o frecuencia")
	}
	if in.FirstPaymentDate == nil {
		if c.Frequency != "fortnightly" {
			return Purchase{}, ErrUnsupported
		}
		if c.FirstPaymentMinDays == nil || *c.FirstPaymentMinDays < 0 {
			return Purchase{}, invalid("falta first_payment_min_days")
		}
	}
	r, ok, err := resolveInterest(in)
	if err != nil {
		return Purchase{}, err
	}
	if !ok {
		return Purchase{}, ErrNoRate
	}
	total := in.Principal + r.interest
	p := Purchase{Interest: r.interest, TotalExact: r.exact, Rate: r.rate, RateBasis: "term_total",
		Source: r.source, PurchaseKind: "installment_plan"}
	return finish(in, p, installmentDates(in), SplitTrunc(total, in.TermCount), SplitHalfUp(in.Principal, in.TermCount),
		provList(serieProv(r.interest, in.TermCount)))
}

// finish arma la serie, las partes y los campos comunes de una compra.
func finish(in PurchaseInput, p Purchase, dates []Date, amounts, principals []Cents, prov []string) (Purchase, error) {
	p.Principal = in.Principal
	p.Total = in.Principal + p.Interest
	if p.TotalExact == nil {
		p.TotalExact = new(big.Rat).SetInt64(int64(p.Total))
	}
	if p.ReferencePayment == 0 {
		p.ReferencePayment = ReferencePayment(p.TotalExact, in.TermCount)
	}
	p.RateDisplay = p.Rate.Display()
	p.PromotionID = in.PromotionID
	if prov == nil {
		prov = []string{}
	}
	p.Provisional = prov
	p.Payments = make([]ScheduledPayment, len(amounts))
	for i := range amounts {
		p.Payments[i] = ScheduledPayment{Seq: i + 1, Due: dates[i], Amount: amounts[i],
			Principal: principals[i], Interest: amounts[i] - principals[i], Shares: []PaymentShare{}}
	}
	shares, err := applyShares(in.Shares, p.Principal, p.Interest, p.Payments)
	if err != nil {
		return Purchase{}, err
	}
	p.Shares = shares
	return p, nil
}

// applyShares reparte intereses y pagos por persona (§4, RN-07, RN-08).
func applyShares(in []ShareInput, principal, interest Cents, pays []ScheduledPayment) ([]PersonShare, error) {
	ss := append([]ShareInput(nil), in...)
	sort.SliceStable(ss, func(i, j int) bool {
		if ss[i].External != ss[j].External {
			return !ss[i].External
		}
		if ss[i].PersonID != ss[j].PersonID {
			return ss[i].PersonID < ss[j].PersonID
		}
		return ss[i].Code < ss[j].Code
	})
	out := make([]PersonShare, len(ss))
	n := len(pays)
	parts := make([][]Cents, len(ss))
	var extInterest Cents
	for i, s := range ss {
		out[i] = PersonShare{PersonID: s.PersonID, Code: s.Code, External: s.External, Principal: s.Principal}
		if s.Joint != nil || s.PaperIn != nil {
			var j, pi Cents
			if s.Joint != nil {
				j = *s.Joint
			}
			if s.PaperIn != nil {
				pi = *s.PaperIn
			}
			out[i].Composition = &Composition{Joint: j, Own: s.Principal - j, PaperIn: pi, Free: s.Principal - pi}
		}
		if i > 0 {
			num := new(big.Int).Mul(big.NewInt(int64(interest)), big.NewInt(int64(s.Principal)))
			out[i].Interest = RoundHalfUp(num, big.NewInt(int64(principal)))
			extInterest += out[i].Interest
		}
	}
	out[0].Interest = interest - extInterest
	for i := range out {
		out[i].Total = out[i].Principal + out[i].Interest
		if i > 0 {
			parts[i] = SplitHalfUp(out[i].Total, n)
		}
	}
	parts[0] = make([]Cents, n)
	for k := range pays {
		rest := pays[k].Amount
		for i := 1; i < len(out); i++ {
			rest -= parts[i][k]
		}
		if rest < 0 {
			return nil, invalid("parte del titular negativa en el pago %d", k+1)
		}
		parts[0][k] = rest
		pays[k].Shares = make([]PaymentShare, len(out))
		for i := range out {
			pays[k].Shares[i] = PaymentShare{PersonID: out[i].PersonID, Code: out[i].Code, Amount: parts[i][k]}
		}
	}
	for i := range out {
		if n > 0 {
			out[i].FirstPayment = parts[i][0]
		}
	}
	return out, nil
}
