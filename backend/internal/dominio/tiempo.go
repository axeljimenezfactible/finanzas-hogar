package dominio

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // America/Mexico_City aunque el sistema no tenga tzdata
)

// MX es la zona de todas las fechas de negocio (decisión B; C1 §4).
var MX = mustLoad("America/Mexico_City")

func mustLoad(name string) *time.Location {
	l, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return l
}

// Formatos de C1 §4 y de la base (strftime('%Y-%m-%dT%H:%M:%fZ')).
const (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02T15:04:05-07:00"
	DBTimeLayout   = "2006-01-02T15:04:05.000Z"
)

var (
	reDate     = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`)
	reDateTime = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]{1,6})?(Z|[+-][0-9]{2}:[0-9]{2})$`)
	reMonth    = regexp.MustCompile(`^([0-9]{4})-(0[1-9]|1[0-2])$`)
	reHalf     = regexp.MustCompile(`^([0-9]{4})-(0[1-9]|1[0-2])-q([12])$`)
	reWeek     = regexp.MustCompile(`^([0-9]{4})-W(0[1-9]|[1-4][0-9]|5[0-3])$`)
	reYear     = regexp.MustCompile(`^[0-9]{4}$`)
	reRange    = regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2})\.\.([0-9]{4}-[0-9]{2}-[0-9]{2})$`)
)

// Date es una fecha de negocio del calendario de America/Mexico_City (sin hora).
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate normaliza (2026-02-31 pasa a 2026-03-03, como time.Date).
func NewDate(y int, m time.Month, d int) Date {
	return dateOf(time.Date(y, m, d, 12, 0, 0, 0, time.UTC))
}

func dateOf(t time.Time) Date { return Date{t.Year(), t.Month(), t.Day()} }

// DateOf da la fecha de negocio de un instante (en America/Mexico_City).
func DateOf(t time.Time) Date { return dateOf(t.In(MX)) }

// Today es la fecha de negocio de hoy según el reloj.
func Today(c Clock) Date { return DateOf(c.Now()) }

func (d Date) t() time.Time { return time.Date(d.Year, d.Month, d.Day, 12, 0, 0, 0, time.UTC) }

// String da AAAA-MM-DD (formato de la API).
func (d Date) String() string { return d.t().Format(DateLayout) }

// IsZero dice si la fecha no tiene valor.
func (d Date) IsZero() bool { return d.Year == 0 && d.Month == 0 && d.Day == 0 }

// AddDays suma días (negativos restan).
func (d Date) AddDays(n int) Date { return dateOf(d.t().AddDate(0, 0, n)) }

// Weekday es el día de la semana.
func (d Date) Weekday() time.Weekday { return d.t().Weekday() }

// Before, After y Compare comparan fechas.
func (d Date) Before(o Date) bool { return d.Compare(o) < 0 }

// After dice si d es posterior a o.
func (d Date) After(o Date) bool { return d.Compare(o) > 0 }

// Compare devuelve -1, 0 o 1.
func (d Date) Compare(o Date) int { return d.t().Compare(o.t()) }

// DaysUntil cuenta los días de d a o (negativo si o es anterior).
func (d Date) DaysUntil(o Date) int { return int(o.t().Sub(d.t()).Hours() / 24) }

// StartOfDay es el instante 00:00 de la fecha en America/Mexico_City.
func (d Date) StartOfDay() time.Time { return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, MX) }

// ISOWeek da la semana ISO 8601 en la sintaxis del filtro period: 2026-W42.
func (d Date) ISOWeek() string {
	y, w := d.t().ISOWeek()
	return fmt.Sprintf("%04d-W%02d", y, w)
}

// MonthKey da AAAA-MM.
func (d Date) MonthKey() string { return fmt.Sprintf("%04d-%02d", d.Year, int(d.Month)) }

// DaysIn da los días del mes (febrero bisiesto incluido).
func DaysIn(y int, m time.Month) int { return time.Date(y, m+1, 0, 12, 0, 0, 0, time.UTC).Day() }

// ParseDate lee AAAA-MM-DD estricto; 2026-02-30 o 10-oct-2026 son error (400 en la API, C1 §3).
func ParseDate(s string) (Date, error) {
	if !reDate.MatchString(s) {
		return Date{}, fmt.Errorf("fecha %q: se espera AAAA-MM-DD", s)
	}
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return Date{}, fmt.Errorf("fecha %q no existe", s)
	}
	return dateOf(t), nil
}

// ParseDateTime lee un instante ISO 8601 con «Z» o con desplazamiento (C1 §4: «Z» solo en entradas).
func ParseDateTime(s string) (time.Time, error) {
	if !reDateTime.MatchString(s) {
		return time.Time{}, fmt.Errorf("instante %q: se espera AAAA-MM-DDTHH:MM:SS con Z o desplazamiento", s)
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("instante %q no existe", s)
	}
	return t, nil
}

// FormatDateTime da el instante con el desplazamiento de America/Mexico_City (2026-10-10T21:30:00-06:00).
func FormatDateTime(t time.Time) string { return t.In(MX).Format(DateTimeLayout) }

// FormatDBTime da el instante en UTC con milisegundos y «Z», igual que los DEFAULT de la base.
func FormatDBTime(t time.Time) string { return t.UTC().Format(DBTimeLayout) }

// ParseDBTime lee un instante guardado en la base (UTC con Z; acepta también RFC 3339).
func ParseDBTime(s string) (time.Time, error) {
	if t, err := time.Parse(DBTimeLayout, s); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339Nano, s)
}

// DBTimeToAPI convierte un instante de la base al formato de respuesta de C1; si no se puede leer, lo deja igual.
func DBTimeToAPI(s string) string {
	t, err := ParseDBTime(s)
	if err != nil {
		return s
	}
	return FormatDateTime(t)
}

// Fortnight es una quincena: Half 1 = del 1 al 15; Half 2 = del 16 al último día (el 31 cae en la 2; K8).
type Fortnight struct {
	Year  int
	Month time.Month
	Half  int
}

// FortnightOf da la quincena de una fecha.
func FortnightOf(d Date) Fortnight {
	if d.Day <= 15 {
		return Fortnight{d.Year, d.Month, 1}
	}
	return Fortnight{d.Year, d.Month, 2}
}

// Start es el primer día de la quincena.
func (f Fortnight) Start() Date {
	if f.Half == 1 {
		return Date{f.Year, f.Month, 1}
	}
	return Date{f.Year, f.Month, 16}
}

// End es el último día de la quincena.
func (f Fortnight) End() Date {
	if f.Half == 1 {
		return Date{f.Year, f.Month, 15}
	}
	return Date{f.Year, f.Month, DaysIn(f.Year, f.Month)}
}

// PayDay es el día de pago: el 15, y el 30 o el último día si el mes es más corto (febrero: 28 o 29; R-23).
func (f Fortnight) PayDay() Date {
	if f.Half == 1 {
		return Date{f.Year, f.Month, 15}
	}
	return Date{f.Year, f.Month, min(30, DaysIn(f.Year, f.Month))}
}

// Next es la quincena siguiente.
func (f Fortnight) Next() Fortnight { return FortnightOf(f.End().AddDays(1)) }

// Prev es la quincena anterior.
func (f Fortnight) Prev() Fortnight { return FortnightOf(f.Start().AddDays(-1)) }

// String da la sintaxis del filtro period: 2026-10-q1.
func (f Fortnight) String() string { return fmt.Sprintf("%04d-%02d-q%d", f.Year, int(f.Month), f.Half) }

// Clases de periodo del filtro global period (C1 §6).
const (
	PeriodMonth     = "month"
	PeriodFortnight = "fortnight"
	PeriodWeek      = "week"
	PeriodYear      = "year"
	PeriodRange     = "range"
)

// Period es un rango inclusivo de fechas con la clase y el texto con que se pidió.
type Period struct {
	Kind  string
	Key   string
	Start Date
	End   Date
}

// Contains dice si la fecha cae en el periodo.
func (p Period) Contains(d Date) bool { return !d.Before(p.Start) && !d.After(p.End) }

// Days cuenta los días del periodo.
func (p Period) Days() int { return p.Start.DaysUntil(p.End) + 1 }

// MonthPeriod es el mes calendario de una fecha.
func MonthPeriod(d Date) Period {
	return Period{PeriodMonth, d.MonthKey(), Date{d.Year, d.Month, 1}, Date{d.Year, d.Month, DaysIn(d.Year, d.Month)}}
}

// ParsePeriod lee el filtro global period (2026-10, 2026-10-q1, 2026-W42, 2026, 2026-10-10..2026-11-09).
// Formato incorrecto → *Error 400 bad_request; rango con fin antes del inicio → *Error 422 (C1 §3).
func ParsePeriod(s string) (Period, error) {
	bad := func() (Period, error) {
		return Period{}, BadRequest(fmt.Sprintf("El periodo «%s» no tiene un formato válido. Usa, por ejemplo, 2026-10 o 2026-10-q1.", s))
	}
	atoi := func(x string) int { n, _ := strconv.Atoi(x); return n }
	switch {
	case reMonth.MatchString(s):
		m := reMonth.FindStringSubmatch(s)
		p := MonthPeriod(Date{atoi(m[1]), time.Month(atoi(m[2])), 1})
		return p, nil
	case reHalf.MatchString(s):
		m := reHalf.FindStringSubmatch(s)
		f := Fortnight{atoi(m[1]), time.Month(atoi(m[2])), atoi(m[3])}
		return Period{PeriodFortnight, s, f.Start(), f.End()}, nil
	case reWeek.MatchString(s):
		m := reWeek.FindStringSubmatch(s)
		y, w := atoi(m[1]), atoi(m[2])
		jan4 := Date{y, time.January, 4}
		monday := jan4.AddDays(-((int(jan4.Weekday()) + 6) % 7))
		start := monday.AddDays(7 * (w - 1))
		if start.ISOWeek() != s {
			return bad() // semana 53 en un año de 52
		}
		return Period{PeriodWeek, s, start, start.AddDays(6)}, nil
	case reYear.MatchString(s):
		y := atoi(s)
		return Period{PeriodYear, s, Date{y, time.January, 1}, Date{y, time.December, 31}}, nil
	case reRange.MatchString(s):
		m := reRange.FindStringSubmatch(s)
		a, err1 := ParseDate(m[1])
		b, err2 := ParseDate(m[2])
		if err1 != nil || err2 != nil {
			return bad()
		}
		if b.Before(a) {
			return Period{}, Invalid(FieldIssue{"period", "out_of_range", "La fecha final del periodo no puede ser anterior a la inicial."})
		}
		return Period{PeriodRange, s, a, b}, nil
	}
	return bad()
}

// ParseHHMM lee una hora HH:MM (00:00 a 23:59), como las claves *.send_time de app_setting.
func ParseHHMM(s string) (hour, minute int, err error) {
	h, m, ok := strings.Cut(s, ":")
	if !ok || len(h) != 2 || len(m) != 2 {
		return 0, 0, fmt.Errorf("hora %q: se espera HH:MM", s)
	}
	hour, err1 := strconv.Atoi(h)
	minute, err2 := strconv.Atoi(m)
	if err1 != nil || err2 != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("hora %q fuera de 00:00 a 23:59", s)
	}
	return hour, minute, nil
}
