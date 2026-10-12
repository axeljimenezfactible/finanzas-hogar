// Package comparador es el comparador de escenarios del motor de créditos (S3): genera por
// búsqueda exhaustiva los escenarios crédito × plazo × fecha ± días × promoción de una compra
// planeada, aplica las 3 restricciones duras (K16), mide los 4 objetivos (K17), normaliza
// mín-máx y recomienda con pesos iguales normalizados (decisión C, RN-11, S3-interfaz.md §8).
// Es puro: sin red, sin base de datos y sin solver (S3-15).
package comparador

import (
	"fmt"

	mc "finanzas-hogar/internal/motor_cred"
)

// Request es la petición del comparador (C1 ScenarioRequest más lo que arma la API).
type Request struct {
	Today       mc.Date          `json:"today"`
	Amount      mc.Cents         `json:"amount_cents"`
	DesiredDate mc.Date          `json:"desired_date"`
	StoreID     *int64           `json:"store_id"`
	Merchant    string           `json:"merchant_name"`
	Shares      []mc.ShareInput  `json:"shares"`
	Global      Constraints      `json:"global"`   // Ajustes › Créditos y condiciones (optimizer.*)
	Override    Constraints      `json:"override"` // las de la comparación (K16)
	Weights     Weights          `json:"weights"`  // optimizer.weights; enteros ≥ 0
	Credits     []mc.CreditState `json:"credits"`
	Caps        []PersonCap      `json:"caps"`    // person.fortnight_payment_cap_cents (K13)
	Baskets     []BasketQuote    `json:"baskets"` // canasta por tienda (K17); PROVISIONAL S3-objetivos-medida (1)
}

// Constraints son las restricciones duras; nil = sin capturar (no hay valor por defecto, K16).
type Constraints struct {
	MinAvailableBuffer *mc.Cents `json:"min_available_buffer_cents"`
	MaxLivePurchases   *int      `json:"max_live_purchases"`
	ZeroDebtBy         *mc.Date  `json:"zero_debt_by_date"`
	DateWindowDays     *int      `json:"date_window_days"`
}

// Weights son los pesos del puntaje (se normalizan: w / Σ w).
type Weights struct {
	LowestInterest   int64 `json:"lowest_interest"`
	EvenLowPayment   int64 `json:"even_low_payment"`
	MaxAvailable     int64 `json:"max_available"`
	MostValuePerPeso int64 `json:"most_value_per_peso"`
}

// PersonCap es el tope de pago por quincena de una persona (no es restricción dura).
type PersonCap struct {
	Person   string   `json:"person"`
	PersonID int64    `json:"person_id,omitempty"`
	Cap      mc.Cents `json:"cap_cents"`
}

// BasketQuote: unidades base (milésimas) y costo de la canasta en una tienda.
type BasketQuote struct {
	StoreID        int64    `json:"store_id"`
	Merchant       string   `json:"merchant_name"`
	BaseUnitsMilli int64    `json:"base_units_milli"`
	Cost           mc.Cents `json:"cost_cents"`
}

// Scenario es un escenario evaluado (C1 Scenario más store_id, std_dev, over_cap y provisional).
type Scenario struct {
	CreditID               int64        `json:"credit_id"`
	TermCount              int          `json:"term_count"`
	PurchaseDate           mc.Date      `json:"purchase_date"`
	StoreID                *int64       `json:"store_id"`
	PromotionID            *int64       `json:"promotion_id"`
	Interest               mc.Cents     `json:"interest_cents"`
	MaxFortnightPayment    mc.Cents     `json:"max_fortnight_payment_cents"`
	FortnightPaymentSpread mc.Cents     `json:"fortnight_payment_spread_cents"`
	StdDevMilliCents       int64        `json:"std_dev_milli_cents"`
	MinAvailable           mc.Cents     `json:"min_available_cents"`
	MaxLivePurchases       int          `json:"max_live_purchases"`
	ZeroDebtDate           mc.Date      `json:"zero_debt_date"`
	ValuePerPesoMilli      *int64       `json:"value_per_peso_milli"`
	ScoreBP                *int64       `json:"score_bp"`
	OverCap                []PersonPeak `json:"over_cap"`
	Provisional            []string     `json:"provisional"`
}

// PersonPeak: pico de pago por quincena de una persona contra su tope.
type PersonPeak struct {
	Person string   `json:"person"`
	Peak   mc.Cents `json:"peak_cents"`
	Cap    mc.Cents `json:"cap_cents"`
}

// Rejected es un candidato que viola una restricción dura.
type Rejected struct {
	Scenario           Scenario `json:"scenario"`
	ViolatedConstraint string   `json:"violated_constraint"`
}

// ExcludedCredit es un crédito que no aporta candidatos.
type ExcludedCredit struct {
	CreditID int64  `json:"credit_id"`
	Reason   string `json:"reason"`
}

// Optima son los 4 óptimos lado a lado.
type Optima struct {
	LowestInterest   Scenario `json:"lowest_interest"`
	EvenLowPayment   Scenario `json:"even_low_payment"`
	MaxAvailable     Scenario `json:"max_available"`
	MostValuePerPeso Scenario `json:"most_value_per_peso"`
}

// Result es el espejo de C1 ScenarioResult más excluded y provisional.
type Result struct {
	Status             string           `json:"status"`
	MissingConstraints []string         `json:"missing_constraints"`
	ConstraintsUsed    Constraints      `json:"constraints_used"`
	EvaluatedCount     int              `json:"evaluated_count"`
	Valid              []Scenario       `json:"valid"`
	Optima             *Optima          `json:"optima"`
	Recommendation     *Scenario        `json:"recommendation"`
	Rejected           []Rejected       `json:"rejected"`
	Excluded           []ExcludedCredit `json:"excluded"`
	Provisional        []string         `json:"provisional"`
}

// Valores de status, violated_constraint y reason (C1).
const (
	StatusOK       = "ok"
	StatusMissing  = "missing_constraints"
	StatusNoValid  = "no_valid_scenario"
	ViolBuffer     = "min_available_buffer"
	ViolLive       = "max_live_purchases"
	ViolZeroDebt   = "zero_debt_by_date"
	ReasonNoLimit  = "no_limit"
	provObjectives = "S3-objetivos-medida"
	scenarioRef    = "_escenario"
)

func invalid(format string, a ...any) error {
	return fmt.Errorf("%w: %s", mc.ErrInvalid, fmt.Sprintf(format, a...))
}

func (r Request) validate() error {
	if r.Today.IsZero() || r.DesiredDate.IsZero() {
		return invalid("faltan today o desired_date")
	}
	if r.Amount <= 0 {
		return invalid("amount_cents debe ser > 0")
	}
	w := r.Weights
	if w.LowestInterest < 0 || w.EvenLowPayment < 0 || w.MaxAvailable < 0 || w.MostValuePerPeso < 0 {
		return invalid("pesos negativos")
	}
	if w.LowestInterest+w.EvenLowPayment+w.MaxAvailable+w.MostValuePerPeso == 0 {
		return invalid("la suma de los pesos es 0")
	}
	for _, c := range []Constraints{r.Global, r.Override} {
		if (c.MinAvailableBuffer != nil && *c.MinAvailableBuffer < 0) || (c.MaxLivePurchases != nil && *c.MaxLivePurchases < 0) ||
			(c.DateWindowDays != nil && *c.DateWindowDays < 0) || (c.ZeroDebtBy != nil && c.ZeroDebtBy.IsZero()) {
			return invalid("restricción fuera de rango")
		}
	}
	holders, sum := 0, mc.Cents(0)
	for _, s := range r.Shares {
		if s.Principal < 0 {
			return invalid("parte negativa")
		}
		if !s.External {
			holders++
		}
		sum += s.Principal
	}
	if holders != 1 || sum != r.Amount {
		return invalid("partes: %d titulares, suma %d de %d", holders, sum, r.Amount)
	}
	for _, c := range r.Caps {
		if c.Person == "" || c.Cap < 0 {
			return invalid("tope por persona inválido")
		}
	}
	for _, b := range r.Baskets {
		if b.Cost <= 0 || b.BaseUnitsMilli < 0 {
			return invalid("canasta de la tienda %d inválida", b.StoreID)
		}
	}
	for _, cs := range r.Credits {
		for _, p := range cs.Purchases {
			if p.Ref == scenarioRef {
				return invalid("ref reservada %q", scenarioRef)
			}
		}
	}
	return nil
}

func pick[T any](o, g *T) *T {
	if o != nil {
		return o
	}
	return g
}

// merge aplica K16: la de la comparación si viene, si no la global, campo por campo; devuelve
// las faltantes en el orden de C1.
func merge(g, o Constraints) (Constraints, []string) {
	u := Constraints{pick(o.MinAvailableBuffer, g.MinAvailableBuffer), pick(o.MaxLivePurchases, g.MaxLivePurchases),
		pick(o.ZeroDebtBy, g.ZeroDebtBy), pick(o.DateWindowDays, g.DateWindowDays)}
	missing := []string{}
	if u.MinAvailableBuffer == nil {
		missing = append(missing, "min_available_buffer_cents")
	}
	if u.MaxLivePurchases == nil {
		missing = append(missing, "max_live_purchases")
	}
	if u.ZeroDebtBy == nil {
		missing = append(missing, "zero_debt_by_date")
	}
	if u.DateWindowDays == nil {
		missing = append(missing, "date_window_days")
	}
	return u, missing
}
