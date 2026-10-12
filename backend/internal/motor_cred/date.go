package motor_cred

import (
	"encoding/json"
	"fmt"
	"time"
)

// Date es una fecha civil sin hora; la zona (America/Mexico_City) la decide quien llama.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// ParseDate lee AAAA-MM-DD y rechaza fechas que no existen (2027-02-30).
func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, invalid("fecha %q", s)
	}
	d := fromTime(t)
	if d.String() != s {
		return Date{}, invalid("fecha %q", s)
	}
	return d, nil
}

func fromTime(t time.Time) Date { return Date{t.Year(), t.Month(), t.Day()} }

func (d Date) t() time.Time { return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC) }

// IsZero dice si la fecha no se capturó.
func (d Date) IsZero() bool { return d.Year == 0 && d.Month == 0 && d.Day == 0 }

// String devuelve la fecha ISO.
func (d Date) String() string { return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day) }

// AddDays suma n días.
func (d Date) AddDays(n int) Date { return fromTime(d.t().AddDate(0, 0, n)) }

// Compare devuelve -1, 0 o 1.
func (d Date) Compare(o Date) int { return d.t().Compare(o.t()) }

// DaysBetween = b − a en días.
func DaysBetween(a, b Date) int { return int(b.t().Sub(a.t()) / (24 * time.Hour)) }

// MarshalJSON emite "AAAA-MM-DD" (null si no hay fecha).
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.String())
}

// UnmarshalJSON lee "AAAA-MM-DD" o null.
func (d *Date) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*d = Date{}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return invalid("fecha %s", b)
	}
	v, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = v
	return nil
}

func lastDay(y int, m time.Month) int { return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day() }

// clipDay arma la fecha con el día recortado al fin de mes; normaliza meses fuera de rango.
func clipDay(y int, m time.Month, day int) Date {
	first := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	y, m = first.Year(), first.Month()
	if l := lastDay(y, m); day > l {
		day = l
	}
	return Date{y, m, day}
}

// Fortnight es una quincena: Half 1 (días 1–15, paga el 15) o 2 (16 al último, paga el 30 o el último).
type Fortnight struct {
	Year   int        `json:"year"`
	Month  time.Month `json:"month"`
	Half   int        `json:"half"`
	PayDay Date       `json:"pay_day"`
}

// FortnightOf aplica RN-06 y K8.
func FortnightOf(d Date) Fortnight {
	if d.Day <= 15 {
		return Fortnight{d.Year, d.Month, 1, Date{d.Year, d.Month, 15}}
	}
	return Fortnight{d.Year, d.Month, 2, clipDay(d.Year, d.Month, 30)}
}

// nextPayDay es el primer día de pago de quincena estrictamente posterior a d.
func nextPayDay(d Date) Date {
	f := FortnightOf(d)
	if f.PayDay.Compare(d) > 0 {
		return f.PayDay
	}
	if f.Half == 1 {
		return clipDay(d.Year, d.Month, 30)
	}
	return clipDay(d.Year, d.Month+1, 15)
}

// fortnightSchedule: primer pago = primera candidata a >= minDays de la compra; luego consecutivas (RN-06).
func fortnightSchedule(purchase Date, minDays, n int) []Date {
	out := make([]Date, 0, n)
	c := nextPayDay(purchase.AddDays(-1))
	for len(out) < n {
		if len(out) > 0 || DaysBetween(purchase, c) >= minDays {
			out = append(out, c)
		}
		c = nextPayDay(c)
	}
	return out
}

// seriesFrom: primera fecha capturada; luego mismo día del mes (recortado), quincenas consecutivas o +7 días.
func seriesFrom(first Date, n int, frequency string) []Date {
	out := make([]Date, 0, n)
	c := first
	for i := 0; i < n; i++ {
		switch frequency {
		case "monthly":
			c = clipDay(first.Year, first.Month+time.Month(i), first.Day)
		case "weekly":
			c = first.AddDays(7 * i)
		default:
			if i > 0 {
				c = nextPayDay(c)
			}
		}
		out = append(out, c)
	}
	return out
}
