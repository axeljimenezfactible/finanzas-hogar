package motor_cred

import (
	"encoding/json"
	"math/big"
	"strings"
)

// Cents es dinero en centavos MXN. Nunca hay flotantes en dinero ni en tasas (S3-17).
type Cents int64

// Rate es una tasa como fracción exacta Num/Den con los decimales capturados (K12).
type Rate struct {
	Num      int64 `json:"num"`
	Den      int64 `json:"den"`
	Decimals int   `json:"decimales"` // -1: sin dato de captura
}

// NewRate valida Den > 0, Num >= 0 y Decimals >= -1.
func NewRate(num, den int64, decimals int) (Rate, error) {
	r := Rate{num, den, decimals}
	if !r.valid() {
		return Rate{}, invalid("tasa %d/%d con %d decimales", num, den, decimals)
	}
	return r, nil
}

func (r Rate) valid() bool { return r.Den > 0 && r.Num >= 0 && r.Decimals >= -1 }

// QuotedRate es la tasa de una compra cotizada: intereses / capital (RN-08, rate_basis term_total).
func QuotedRate(interest, principal Cents, decimals int) Rate {
	return Rate{int64(interest), int64(principal), decimals}
}

// Rat devuelve la fracción exacta (0 si la tasa es inválida).
func (r Rate) Rat() *big.Rat {
	if r.Den == 0 {
		return new(big.Rat)
	}
	return big.NewRat(r.Num, r.Den)
}

// Display redondea HALF_UP el porcentaje a los decimales capturados: "25.8685 %" (K12).
// Sin dato de captura (-1) devuelve "": el motor no inventa decimales (§5.1; C1 omite rate_pct).
func (r Rate) Display() string {
	if r.Den <= 0 || r.Decimals < 0 {
		return ""
	}
	d := r.Decimals
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(d)), nil)
	num := new(big.Int).Mul(big.NewInt(r.Num), big.NewInt(100))
	num.Mul(num, scale)
	q := roundHalfUpInt(num, big.NewInt(r.Den))
	s := q.String()
	if d == 0 {
		return s + " %"
	}
	if len(s) <= d {
		s = strings.Repeat("0", d-len(s)+1) + s
	}
	return s[:len(s)-d] + "." + s[len(s)-d:] + " %"
}

// roundHalfUpInt redondea num/den al entero con ROUND_HALF_UP (empates lejos de cero); den > 0.
func roundHalfUpInt(num, den *big.Int) *big.Int {
	n := new(big.Int).Abs(num)
	n.Mul(n, big.NewInt(2))
	n.Add(n, den)
	d := new(big.Int).Mul(den, big.NewInt(2))
	q := n.Quo(n, d)
	if num.Sign() < 0 {
		q.Neg(q)
	}
	return q
}

// RoundHalfUp redondea num/den (en centavos) a centavo con ROUND_HALF_UP; den > 0.
func RoundHalfUp(num, den *big.Int) Cents {
	return Cents(roundHalfUpInt(num, den).Int64())
}

func roundRat(r *big.Rat) Cents { return RoundHalfUp(r.Num(), r.Denom()) }

func divHalfUp(c Cents, n int64) Cents {
	return RoundHalfUp(big.NewInt(int64(c)), big.NewInt(n))
}

// ApplyRate = HALF_UP(c × Num / Den).
func ApplyRate(c Cents, r Rate) Cents {
	if r.Den <= 0 {
		return 0
	}
	n := new(big.Int).Mul(big.NewInt(int64(c)), big.NewInt(r.Num))
	return RoundHalfUp(n, big.NewInt(r.Den))
}

// SplitTrunc reparte total en n montos con base truncada; el último ajusta (RN-07, montos).
func SplitTrunc(total Cents, n int) []Cents {
	if n <= 0 {
		return nil
	}
	return fill(total/Cents(n), total, n)
}

// SplitHalfUp reparte total en n con base HALF_UP; el último ajusta (RN-07, capital y partes).
func SplitHalfUp(total Cents, n int) []Cents {
	if n <= 0 {
		return nil
	}
	return fill(divHalfUp(total, int64(n)), total, n)
}

func fill(base, total Cents, n int) []Cents {
	out := make([]Cents, n)
	for i := range out {
		out[i] = base
	}
	out[n-1] = total - base*Cents(n-1)
	return out
}

// ReferencePayment = HALF_UP(total sin redondear / n): solo cotizaciones y sensibilidad (§1.2).
func ReferencePayment(totalExact *big.Rat, n int) Cents {
	if totalExact == nil || n <= 0 {
		return 0
	}
	den := new(big.Int).Mul(totalExact.Denom(), big.NewInt(int64(n)))
	return RoundHalfUp(totalExact.Num(), den)
}

// grossUp = c × (1 + r) exacto.
func grossUp(c Cents, r Rate) *big.Rat {
	x := new(big.Rat).Add(big.NewRat(1, 1), r.Rat())
	return x.Mul(x, new(big.Rat).SetInt64(int64(c)))
}

// ratJSON es la forma JSON de una fracción de salida: {"num", "den"} reducida.
type ratJSON struct {
	Num *big.Int `json:"num"`
	Den *big.Int `json:"den"`
}

func toRatJSON(r *big.Rat) *ratJSON {
	if r == nil {
		return nil
	}
	return &ratJSON{new(big.Int).Set(r.Num()), new(big.Int).Set(r.Denom())}
}

// MarshalJSON emite TotalExact como {"num", "den"} (S3-interfaz.md §3).
func (p Purchase) MarshalJSON() ([]byte, error) {
	type alias Purchase
	return json.Marshal(struct {
		alias
		TotalExact *ratJSON `json:"total_exact"`
	}{alias(p), toRatJSON(p.TotalExact)})
}
