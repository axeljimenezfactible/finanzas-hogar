package dominio

import (
	"math/big"
	"strconv"
	"strings"
)

// Dinero: siempre int64 en centavos de MXN (MoneyCents de C1; protocolo §6). Nunca float.
// El redondeo es ROUND_HALF_UP (la mitad se aleja del cero, como decimal.ROUND_HALF_UP de Python).
// Las series de pagos («el último ajusta centavos») las calculan los motores (motor_fin, motor_cred); aquí no se repiten.

// RoundHalfUp da num/den redondeado ROUND_HALF_UP. den no puede ser 0.
func RoundHalfUp(num, den int64) int64 {
	if den == 0 {
		panic("dominio.RoundHalfUp: denominador 0")
	}
	n, d := big.NewInt(num), big.NewInt(den)
	return roundBig(n, d)
}

// MulDivRound da a·num/den con ROUND_HALF_UP y sin desbordar (aritmética exacta).
func MulDivRound(a, num, den int64) int64 {
	if den == 0 {
		panic("dominio.MulDivRound: denominador 0")
	}
	n := new(big.Int).Mul(big.NewInt(a), big.NewInt(num))
	return roundBig(n, big.NewInt(den))
}

func roundBig(n, d *big.Int) int64 {
	if d.Sign() < 0 {
		n, d = new(big.Int).Neg(n), new(big.Int).Neg(d)
	}
	neg := n.Sign() < 0
	abs := new(big.Int).Abs(n)
	q, r := new(big.Int).QuoRem(abs, d, new(big.Int))
	if new(big.Int).Mul(r, big.NewInt(2)).Cmp(d) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if neg {
		q.Neg(q)
	}
	return q.Int64()
}

// ApplyBasisPoints da cents·bp/10000 con ROUND_HALF_UP (10000 bp = 100 %).
func ApplyBasisPoints(cents, bp int64) int64 { return MulDivRound(cents, bp, 10000) }

// PercentOf da el porcentaje entero part/total·100 con ROUND_HALF_UP (used_pct, progress_pct); 0 si total = 0.
func PercentOf(part, total int64) int64 {
	if total == 0 {
		return 0
	}
	return MulDivRound(part, 100, total)
}

// FormatPesos da el importe como en pantalla y en los mensajes: $1,258.68 (negativo: -$1,258.68).
func FormatPesos(cents int64) string {
	sign := ""
	u := uint64(cents)
	if cents < 0 {
		sign, u = "-", uint64(-cents)
	}
	ent, frac := u/100, u%100
	s := strconv.FormatUint(ent, 10)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return sign + "$" + b.String() + "." + strconv.FormatUint(frac/10, 10) + strconv.FormatUint(frac%10, 10)
}
