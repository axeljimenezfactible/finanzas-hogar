package recetas

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// Aritmética exacta (C-01, C-02, C-11): fracciones big.Rat y un solo redondeo al final.

func r64(n int64) *big.Rat        { return new(big.Rat).SetInt64(n) }
func frac(n, d int64) *big.Rat    { return big.NewRat(n, d) }
func mulR(a, b *big.Rat) *big.Rat { return new(big.Rat).Mul(a, b) }
func addR(a, b *big.Rat) *big.Rat { return new(big.Rat).Add(a, b) }
func subR(a, b *big.Rat) *big.Rat { return new(big.Rat).Sub(a, b) }
func quoR(a, b *big.Rat) *big.Rat { return new(big.Rat).Quo(a, b) }

func toInt64(z *big.Int) int64 {
	if !z.IsInt64() {
		panic(overflowPanic{})
	}
	return z.Int64()
}

// roundHalfUp redondea al entero más cercano; la mitad se aleja de cero (ROUND_HALF_UP, C-02).
func roundHalfUp(r *big.Rat) int64 {
	num := new(big.Int).Abs(r.Num())
	den := r.Denom()
	q, m := new(big.Int).QuoRem(num, den, new(big.Int))
	if m.Lsh(m, 1).Cmp(den) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if r.Sign() < 0 {
		q.Neg(q)
	}
	return toInt64(q)
}

// ceilRat es el techo (hacia +infinito); floorRat, el piso.
func ceilRat(r *big.Rat) int64 {
	q, m := new(big.Int).DivMod(r.Num(), r.Denom(), new(big.Int))
	if m.Sign() != 0 {
		q.Add(q, big.NewInt(1))
	}
	return toInt64(q)
}

func floorRat(r *big.Rat) int64 {
	q, _ := new(big.Int).DivMod(r.Num(), r.Denom(), new(big.Int))
	return toInt64(q)
}

// Frac es una fracción exacta reducida para la salida ({"num": n, "den": d}, FORMATO §3).
type Frac struct {
	Num int64 `json:"num"`
	Den int64 `json:"den"`
}

func fracOf(r *big.Rat) Frac { return Frac{toInt64(r.Num()), toInt64(r.Denom())} }

// Rat devuelve la fracción como big.Rat.
func (f Frac) Rat() *big.Rat { return big.NewRat(f.Num, f.Den) }

func ptr[T any](v T) *T { return &v }

// parseMilli lee un decimal en texto (hasta 3 decimales, sin signo) y lo pasa a milésimas sin flotantes:
// «884» → 884000, «52.5» → 52500.
func parseMilli(s string) (int64, error) {
	ip, fp, dot := strings.Cut(s, ".")
	if !digits(ip) || (dot && (!digits(fp) || len(fp) > 3)) || len(ip) > 12 {
		return 0, fmt.Errorf("decimal inválido %q (hasta 3 decimales, sin signo)", s)
	}
	fp += strings.Repeat("0", 3-len(fp))
	n, err := strconv.ParseInt(ip+fp, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("decimal inválido %q", s)
	}
	return n, nil
}

func digits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Date es una fecha civil (AAAA-MM-DD) sin hora; «hoy» ya viene resuelto en America/Mexico_City (C-05).
type Date struct{ t time.Time }

// ParseDate lee una fecha ISO.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Date{}, errf(ErrInvalidDate, "fecha %q", s)
	}
	return Date{t}, nil
}

func (d Date) String() string               { return d.t.Format(time.DateOnly) }
func (d Date) IsZero() bool                 { return d.t.IsZero() }
func (d Date) AddDays(n int) Date           { return Date{d.t.AddDate(0, 0, n)} }
func (d Date) Before(o Date) bool           { return d.t.Before(o.t) }
func (d Date) DaysUntil(o Date) int64       { return int64(o.t.Sub(d.t) / (24 * time.Hour)) }
func (d Date) inRange(a, b Date) bool       { return !d.Before(a) && !b.Before(d) }
func (d Date) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

// ISOWeekday: 1 = lunes … 7 = domingo.
func (d Date) ISOWeekday() int {
	if w := int(d.t.Weekday()); w != 0 {
		return w
	}
	return 7
}

// UnmarshalJSON lee la fecha ISO; un texto inválido es ErrInvalidDate.
func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errf(ErrInvalidDate, "%s", b)
	}
	x, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = x
	return nil
}

// YearsUntil son los años cumplidos de d (nacimiento) a today (§10.1).
func (d Date) YearsUntil(today Date) int64 {
	y := today.t.Year() - d.t.Year()
	if today.t.Month() < d.t.Month() || (today.t.Month() == d.t.Month() && today.t.Day() < d.t.Day()) {
		y--
	}
	return int64(y)
}

// DefaultWeek es el lunes de la semana de today; si today es domingo, el lunes siguiente (K38, §5.2).
func DefaultWeek(today Date) Date {
	monday := today.AddDays(1 - today.ISOWeekday())
	if today.ISOWeekday() == 7 {
		monday = monday.AddDays(7)
	}
	return monday
}

// WeekCode es la semana ISO «AAAA-Www» de un lunes (K8: lunes a domingo).
func WeekCode(monday Date) string {
	y, w := monday.t.ISOWeek()
	return fmt.Sprintf("%04d-W%02d", y, w)
}

// ParseWeek acepta «AAAA-Www» o la fecha ISO de un lunes (otra fecha es ErrNotMonday).
func ParseWeek(s string) (Date, error) {
	if len(s) == 8 && s[4:6] == "-W" {
		y, e1 := strconv.Atoi(s[:4])
		w, e2 := strconv.Atoi(s[6:])
		if e1 == nil && e2 == nil && w >= 1 && w <= 53 {
			jan4 := Date{time.Date(y, 1, 4, 0, 0, 0, 0, time.UTC)}
			mon := jan4.AddDays(1 - jan4.ISOWeekday() + (w-1)*7)
			if WeekCode(mon) == s {
				return mon, nil
			}
		}
		return Date{}, errf(ErrInvalidDate, "semana %q", s)
	}
	d, err := ParseDate(s)
	if err != nil {
		return Date{}, err
	}
	if d.ISOWeekday() != 1 {
		return Date{}, errf(ErrNotMonday, "%s no es lunes", s)
	}
	return d, nil
}
