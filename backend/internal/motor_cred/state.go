package motor_cred

import (
	"encoding/json"
	"io"
	"math/big"
	"sort"
	"strings"
)

// Kind es credit.kind.
type Kind string

// Tipos de crédito (decisión C).
const (
	KindRevolving Kind = "revolving_line"
	KindCard      Kind = "bank_card"
	KindLoan      Kind = "term_loan"
	KindStore     Kind = "store_credit"
)

// Reading es credit.availability_reading: "A", "B" o "" (= B; RN-09, S3-03).
type Reading string

// Credit es el espejo de credit (S3-interfaz.md §5).
type Credit struct {
	ID                  int64       `json:"id"`
	Name                string      `json:"name"`
	Kind                Kind        `json:"kind"`
	Holder              string      `json:"holder"`
	HolderPersonID      int64       `json:"holder_person_id,omitempty"`
	Limit               *Cents      `json:"limit_cents"`
	Frequency           string      `json:"frequency"`
	FirstPaymentMinDays *int        `json:"first_payment_min_days"`
	Reading             Reading     `json:"reading"`
	ReadingConfirmed    bool        `json:"reading_confirmed"`
	ReadingAAnchor      *Date       `json:"reading_a_anchor"`
	StatementDay        *int        `json:"statement_day"`
	PaymentDueDay       *int        `json:"payment_due_day"`
	AnnualRate          *Rate       `json:"annual_rate"`
	RateBase            string      `json:"rate_base"`
	VATOnInterestBP     *int        `json:"vat_on_interest_bp"`
	MinPaymentPctBP     *int        `json:"min_payment_pct_bp"`
	MinPaymentFixed     *Cents      `json:"min_payment_fixed_cents"`
	TermRates           []TermRate  `json:"term_rates"`
	Promotions          []Promotion `json:"promotions"`
}

// TermRate es credit_term_rate.
type TermRate struct {
	TermCount int    `json:"term_count"`
	Rate      Rate   `json:"rate"`
	Source    string `json:"source"`
	QuoteDate *Date  `json:"quote_date"`
}

// Promotion es promotion (captura manual; el motor nunca precarga).
type Promotion struct {
	ID          int64  `json:"id"`
	CreditID    *int64 `json:"credit_id"`
	Merchant    string `json:"merchant"`
	StoreID     *int64 `json:"store_id"`
	TermCount   int    `json:"term_count"`
	Rate        Rate   `json:"rate"`
	MinPurchase Cents  `json:"min_purchase_cents"`
	ValidFrom   Date   `json:"valid_from"`
	ValidTo     Date   `json:"valid_to"`
	Active      bool   `json:"active"`
}

// ShareInput es la parte de una persona en una compra (§4).
type ShareInput struct {
	PersonID  int64  `json:"person_id,omitempty"`
	Code      string `json:"person"`
	External  bool   `json:"external"`
	Principal Cents  `json:"principal_cents"`
	Joint     *Cents `json:"joint_cents"`
	PaperIn   *Cents `json:"paper_in_cents"`
}

// PersonShare es la parte calculada de una persona.
type PersonShare struct {
	PersonID     int64        `json:"person_id,omitempty"`
	Code         string       `json:"person"`
	External     bool         `json:"external"`
	Principal    Cents        `json:"principal_cents"`
	Interest     Cents        `json:"interest_cents"`
	Total        Cents        `json:"total_cents"`
	FirstPayment Cents        `json:"first_payment_cents"`
	Composition  *Composition `json:"composition"`
}

// Composition es la composición informativa de una parte (RN-04).
type Composition struct {
	Joint   Cents `json:"joint_cents"`
	Own     Cents `json:"own_cents"`
	PaperIn Cents `json:"paper_in_cents"`
	Free    Cents `json:"free_cents"`
}

// PurchaseInput es lo capturado de una compra.
type PurchaseInput struct {
	Credit           Credit       `json:"-"`
	Date             Date         `json:"date"`
	Principal        Cents        `json:"principal_cents"`
	TermCount        int          `json:"term_count"`
	QuotedInterest   *Cents       `json:"quoted_interest_cents"`
	RateDecimals     *int         `json:"rate_decimals"`
	PromotionID      *int64       `json:"promotion_id"`
	StoreID          *int64       `json:"store_id"`
	Merchant         string       `json:"merchant"`
	FirstPaymentDate *Date        `json:"first_payment_date"`
	Shares           []ShareInput `json:"shares"`
}

// Purchase es una compra calculada con su serie (§5).
type Purchase struct {
	Principal        Cents              `json:"principal_cents"`
	Interest         Cents              `json:"interest_cents"`
	Total            Cents              `json:"total_cents"`
	TotalExact       *big.Rat           `json:"total_exact"`
	ReferencePayment Cents              `json:"reference_payment_cents"`
	Rate             Rate               `json:"rate"`
	RateDisplay      string             `json:"rate_display"`
	RateBasis        string             `json:"rate_basis"`
	Source           string             `json:"source"`
	PurchaseKind     string             `json:"purchase_kind"`
	PromotionID      *int64             `json:"promotion_id"`
	Payments         []ScheduledPayment `json:"payments"`
	Shares           []PersonShare      `json:"shares"`
	Provisional      []string           `json:"provisional"`
}

// ScheduledPayment es un pago programado.
type ScheduledPayment struct {
	Seq       int            `json:"seq"`
	Due       Date           `json:"due_date"`
	Amount    Cents          `json:"amount_cents"`
	Principal Cents          `json:"principal_cents"`
	Interest  Cents          `json:"interest_cents"`
	Shares    []PaymentShare `json:"shares"`
}

// PaymentShare es la parte de una persona en un pago.
type PaymentShare struct {
	PersonID int64  `json:"person_id,omitempty"`
	Code     string `json:"person"`
	Amount   Cents  `json:"amount_cents"`
}

// CashEvent es una compra de contado del hogar: no mueve ninguna línea.
type CashEvent struct {
	Date        Date   `json:"date"`
	Amount      Cents  `json:"amount_cents"`
	Person      string `json:"person"`
	Description string `json:"description"`
}

// PaymentRecord es un abono registrado (credit_payment).
type PaymentRecord struct {
	PurchaseRef   string `json:"purchase_ref"`
	Seq           int    `json:"seq"`
	Date          Date   `json:"date"`
	Payer         string `json:"payer"`
	PayerPersonID int64  `json:"payer_person_id,omitempty"`
	Amount        Cents  `json:"amount_cents"`
}

// StatementInput es un corte (credit_statement).
type StatementInput struct {
	StatementDate Date    `json:"statement_date"`
	DueDate       Date    `json:"due_date"`
	Balance       Cents   `json:"balance_cents"`
	DailyBalances []Cents `json:"daily_balances_cents"`
	Source        string  `json:"source"`
}

// PurchaseState es una compra registrada en un crédito.
type PurchaseState struct {
	Ref    string        `json:"ref"`
	Status string        `json:"status"`
	Input  PurchaseInput `json:"input"`
}

// CreditState es un crédito con sus compras, abonos y cortes (§5.7).
type CreditState struct {
	Credit     Credit           `json:"credit"`
	Purchases  []PurchaseState  `json:"purchases"`
	Payments   []PaymentRecord  `json:"payments"`
	Statements []StatementInput `json:"statements"`
}

// Input es entradas.datos de las fixtures; la API arma lo mismo desde la base.
type Input struct {
	Today   Date          `json:"today"`
	Credits []CreditState `json:"credits"`
	Cash    []CashEvent   `json:"cash"`
}

// DecodeInput decodifica estricto: un campo desconocido es ErrUnknownField (unknown_field).
func DecodeInput(r io.Reader) (Input, error) {
	var in Input
	if err := DecodeStrict(r, &in); err != nil {
		return Input{}, err
	}
	return in, in.Validate()
}

// DecodeStrict decodifica JSON con DisallowUnknownFields y traduce el error a ErrInvalid.
func DecodeStrict(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			return invalidField(err)
		}
		if strings.Contains(err.Error(), ErrInvalid.Error()) {
			return invalid("%v", err)
		}
		return invalid("JSON: %v", err)
	}
	return nil
}

func invalidField(err error) error { return &wrapErr{ErrUnknownField, err.Error()} }

type wrapErr struct {
	base error
	msg  string
}

func (e *wrapErr) Error() string { return e.base.Error() + ": " + e.msg }
func (e *wrapErr) Unwrap() error { return e.base }

// Validate revisa las reglas de §5.7 en todos los créditos.
func (in Input) Validate() error {
	seen := map[int64]bool{}
	for _, cs := range in.Credits {
		if seen[cs.Credit.ID] {
			return invalid("crédito %d repetido", cs.Credit.ID)
		}
		seen[cs.Credit.ID] = true
		if err := cs.validate(); err != nil {
			return err
		}
	}
	return nil
}

// CreditByID busca un crédito de la entrada.
func (in Input) CreditByID(id int64) (CreditState, error) {
	for _, cs := range in.Credits {
		if cs.Credit.ID == id {
			return cs, nil
		}
	}
	return CreditState{}, invalid("crédito %d inexistente", id)
}

func (s CreditState) validate() error {
	refs := map[string]bool{}
	for _, p := range s.Purchases {
		if p.Ref == "" || refs[p.Ref] {
			return invalid("compra %q vacía o repetida", p.Ref)
		}
		switch p.Status {
		case "active", "paid_off", "cancelled":
		default:
			return invalid("estado %q", p.Status)
		}
		refs[p.Ref] = true
	}
	for _, r := range s.Payments {
		if !refs[r.PurchaseRef] {
			return invalid("abono de una compra inexistente %q", r.PurchaseRef)
		}
	}
	if len(s.Statements) > 0 && s.Credit.Kind != KindCard {
		return invalid("cortes en un crédito que no es tarjeta")
	}
	return nil
}

type quotedPurchase struct {
	Ref      string
	Status   string
	Input    PurchaseInput
	Purchase Purchase
}

func (q quotedPurchase) live() bool { return q.Status != "cancelled" }

// quoted calcula cada compra en orden (Input.Date, Ref).
func (s CreditState) quoted() ([]quotedPurchase, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	ps := append([]PurchaseState(nil), s.Purchases...)
	sort.SliceStable(ps, func(i, j int) bool {
		if c := ps[i].Input.Date.Compare(ps[j].Input.Date); c != 0 {
			return c < 0
		}
		return ps[i].Ref < ps[j].Ref
	})
	out := make([]quotedPurchase, 0, len(ps))
	for _, p := range ps {
		in := p.Input
		in.Credit = s.Credit
		pu, err := QuotePurchase(in)
		if err != nil {
			return nil, err
		}
		out = append(out, quotedPurchase{p.Ref, p.Status, in, pu})
	}
	return out, nil
}

// Quote aplica QuotePurchase a cada compra, en orden (Input.Date, Ref).
func (s CreditState) Quote() ([]Purchase, error) {
	qs, err := s.quoted()
	if err != nil {
		return nil, err
	}
	out := make([]Purchase, len(qs))
	for i, q := range qs {
		out[i] = q.Purchase
	}
	return out, nil
}

// recordsFor devuelve los abonos de un pago en orden (Date, Payer, Amount).
func (s CreditState) recordsFor(ref string, seq int) []PaymentRecord {
	var out []PaymentRecord
	for _, r := range s.Payments {
		if r.PurchaseRef == ref && r.Seq == seq {
			out = append(out, r)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if c := out[i].Date.Compare(out[j].Date); c != 0 {
			return c < 0
		}
		if out[i].Payer != out[j].Payer {
			return out[i].Payer < out[j].Payer
		}
		return out[i].Amount < out[j].Amount
	})
	return out
}

func sortedCredits(credits []CreditState) ([]CreditState, error) {
	cs := append([]CreditState(nil), credits...)
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].Credit.ID < cs[j].Credit.ID })
	for i := 1; i < len(cs); i++ {
		if cs[i].Credit.ID == cs[i-1].Credit.ID {
			return nil, invalid("crédito %d repetido", cs[i].Credit.ID)
		}
	}
	return cs, nil
}
