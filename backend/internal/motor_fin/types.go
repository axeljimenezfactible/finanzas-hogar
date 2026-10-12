// Package motor_fin es el motor financiero (S2) de Finanzas Hogar: devengado
// contra efectivo, presupuesto, flujo por quincena, proyección de 12 meses,
// deuda viva, promedios, metas y valor estable (docs/contratos/motores/S2-interfaz.md).
// Es puro y determinista: sin reloj, red ni base de datos. El dinero va en
// centavos enteros (Cents) y solo se redondea donde lo dice la interfaz §1.3.
package motor_fin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Cents es dinero en centavos enteros (S2-13: nunca punto flotante).
type Cents int64

// Códigos de error de la interfaz §3.15.
const (
	CodeUnknownField    = "unknown_field"
	CodeInvalidDate     = "invalid_date"
	CodeInvalidZone     = "invalid_zone"
	CodeInvalidArgument = "invalid_argument"
	CodeSharesMismatch  = "shares_mismatch"
	CodeSumMismatch     = "sum_mismatch"
	CodeMissingPurchase = "missing_purchase"
	CodeIncomeNotHolder = "income_not_holder"
	CodeBudgetOverlap   = "budget_overlap"
	CodeOutOfRange      = "out_of_range"
	CodeCadenceBefore7  = "cadence_before_month_7"
	CodeOracleViewShape = "oracle_view_shape"
)

// DefaultZone es la única zona admitida (decisión B, C-05).
const DefaultZone = "America/Mexico_City"

// Error es el error del motor: código de §3.15 y detalle legible.
type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Detail }

func errf(code, format string, a ...any) *Error {
	return &Error{Code: code, Detail: fmt.Sprintf(format, a...)}
}

// ---------- Fechas, meses, quincenas y semanas ----------

// Date es una fecha civil sin hora.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate normaliza (p. ej. 31 de abril pasa a 1 de mayo).
func NewDate(y int, m time.Month, d int) Date {
	t := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return Date{t.Year(), t.Month(), t.Day()}
}

// ParseDate lee "AAAA-MM-DD".
func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil || len(s) != 10 {
		return Date{}, errf(CodeInvalidDate, "fecha inválida %q", s)
	}
	return Date{t.Year(), t.Month(), t.Day()}, nil
}

func (d Date) ord() int64 {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Unix() / 86400
}
func (d Date) AddDays(n int) Date           { return NewDate(d.Year, d.Month, d.Day+n) }
func (d Date) Before(o Date) bool           { return d.ord() < o.ord() }
func (d Date) After(o Date) bool            { return d.ord() > o.ord() }
func (d Date) MonthOf() Month               { return Month{d.Year, d.Month} }
func (d Date) String() string               { return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day) }
func (d Date) in(a, b Date) bool            { return !d.Before(a) && !d.After(b) }
func daysIncl(a, b Date) int                { return int(b.ord()-a.ord()) + 1 }
func (d Date) display() string              { return fmt.Sprintf("%d-%s-%d", d.Day, monthAbbr[d.Month-1], d.Year) }
func (d Date) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errf(CodeInvalidDate, "fecha no es texto: %s", b)
	}
	p, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = p
	return nil
}

var monthAbbr = [12]string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}

// Month es un mes calendario "AAAA-MM".
type Month struct {
	Year int
	Mon  time.Month
}

func ParseMonth(s string) (Month, error) {
	t, err := time.Parse("2006-01", s)
	if err != nil || len(s) != 7 {
		return Month{}, errf(CodeInvalidArgument, "mes inválido %q", s)
	}
	return Month{t.Year(), t.Month()}, nil
}

func (m Month) idx() int                     { return m.Year*12 + int(m.Mon) - 1 }
func monthFromIdx(i int) Month               { return Month{i / 12, time.Month(i%12 + 1)} }
func (m Month) Add(k int) Month              { return monthFromIdx(m.idx() + k) }
func (m Month) First() Date                  { return Date{m.Year, m.Mon, 1} }
func (m Month) Last() Date                   { return NewDate(m.Year, m.Mon+1, 0) }
func (m Month) Days() int                    { return m.Last().Day }
func (m Month) String() string               { return fmt.Sprintf("%04d-%02d", m.Year, int(m.Mon)) }
func (m Month) Label() string                { return m.First().display() + " a " + m.Last().display() }
func (m Month) IsZero() bool                 { return m.Year == 0 }
func (m Month) MarshalJSON() ([]byte, error) { return json.Marshal(m.String()) }

func (m *Month) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errf(CodeInvalidArgument, "mes no es texto: %s", b)
	}
	p, err := ParseMonth(s)
	if err != nil {
		return err
	}
	*m = p
	return nil
}

// MonthBounds da el primer y el último día del mes.
func MonthBounds(m Month) (start, end Date) { return m.First(), m.Last() }

// Fortnight es una quincena "AAAA-MM-q1" (días 1 a 15) o "AAAA-MM-q2" (16 al último), K8.
type Fortnight struct {
	Year int
	Mon  time.Month
	Half int
}

func ParseFortnight(s string) (Fortnight, error) {
	if len(s) != 10 || s[7:9] != "-q" || (s[9] != '1' && s[9] != '2') {
		return Fortnight{}, errf(CodeInvalidArgument, "quincena inválida %q", s)
	}
	m, err := ParseMonth(s[:7])
	if err != nil {
		return Fortnight{}, errf(CodeInvalidArgument, "quincena inválida %q", s)
	}
	return Fortnight{m.Year, m.Mon, int(s[9] - '0')}, nil
}

func (f Fortnight) idx() int { return (Month{f.Year, f.Mon}).idx()*2 + f.Half - 1 }
func (f Fortnight) Add(k int) Fortnight {
	i := f.idx() + k
	m := monthFromIdx(i / 2)
	return Fortnight{m.Year, m.Mon, i%2 + 1}
}
func (f Fortnight) String() string               { return fmt.Sprintf("%04d-%02d-q%d", f.Year, int(f.Mon), f.Half) }
func (f Fortnight) MarshalJSON() ([]byte, error) { return json.Marshal(f.String()) }

func (f *Fortnight) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errf(CodeInvalidArgument, "quincena no es texto: %s", b)
	}
	p, err := ParseFortnight(s)
	if err != nil {
		return err
	}
	*f = p
	return nil
}

// FortnightOf: el día 31 cae en la quincena 2 (K8, S2-06).
func FortnightOf(d Date) Fortnight {
	h := 1
	if d.Day > 15 {
		h = 2
	}
	return Fortnight{d.Year, d.Month, h}
}

// FortnightBounds da inicio, fin y día de pago (15; 30 o último día si el mes es más corto).
func FortnightBounds(f Fortnight) (start, end, pay Date) {
	m := Month{f.Year, f.Mon}
	if f.Half == 1 {
		return Date{f.Year, f.Mon, 1}, Date{f.Year, f.Mon, 15}, Date{f.Year, f.Mon, 15}
	}
	last := m.Last()
	pd := 30
	if last.Day < 30 {
		pd = last.Day
	}
	return Date{f.Year, f.Mon, 16}, last, Date{f.Year, f.Mon, pd}
}

// Week es una semana ISO 8601, de lunes a domingo (K8).
type Week struct {
	Year int
	Num  int
}

func (w Week) String() string               { return fmt.Sprintf("%04d-W%02d", w.Year, w.Num) }
func (w Week) MarshalJSON() ([]byte, error) { return json.Marshal(w.String()) }

func WeekOf(d Date) Week {
	y, n := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).ISOWeek()
	return Week{y, n}
}

func WeekBounds(w Week) (monday, sunday Date) {
	jan4 := time.Date(w.Year, 1, 4, 0, 0, 0, 0, time.UTC)
	off := (int(jan4.Weekday()) + 6) % 7
	monday = NewDate(w.Year, 1, 4-off+(w.Num-1)*7)
	return monday, monday.AddDays(6)
}

// Instant es un instante RFC 3339 en UTC con «Z» (§1.1, K12). Otro desfase, aunque sea +00:00, es
// invalid_date, igual que un instante mal formado (REVISION-S2-r3 B1).
type Instant struct{ time.Time }

func (t *Instant) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errf(CodeInvalidDate, "instante no es texto: %s", b)
	}
	if !strings.HasSuffix(s, "Z") {
		return errf(CodeInvalidDate, "instante %q: va en UTC con «Z» (§1.1)", s)
	}
	p, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return errf(CodeInvalidDate, "instante inválido %q", s)
	}
	t.Time = p
	return nil
}

// Today resuelve la fecha civil del instante en la zona (S2-14), con la base de zonas del sistema.
func Today(now time.Time, zone string) (Date, error) {
	if zone == "" {
		zone = DefaultZone
	}
	if zone != DefaultZone {
		return Date{}, errf(CodeInvalidZone, "zona no admitida %q", zone)
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return Date{}, errf(CodeInvalidZone, "sin base de zonas para %q: %v", zone, err)
	}
	t := now.In(loc)
	return Date{t.Year(), t.Month(), t.Day()}, nil
}

// ---------- Mapas ordenados de salida ----------

// CentsMap es un mapa {code: cents} que se serializa en el orden de Keys (§1.1).
type CentsMap struct {
	Keys []string
	Vals map[string]Cents
}

func newCentsMap(keys []string) CentsMap {
	m := CentsMap{Keys: append([]string(nil), keys...), Vals: map[string]Cents{}}
	for _, k := range keys {
		m.Vals[k] = 0
	}
	return m
}

// catMap arma un by_category: solo importes distintos de 0, en orden de bytes del código.
func catMap(v map[string]Cents) CentsMap {
	m := CentsMap{Vals: map[string]Cents{}}
	for k, c := range v {
		if c != 0 {
			m.Keys = append(m.Keys, k)
			m.Vals[k] = c
		}
	}
	sort.Strings(m.Keys)
	return m
}

func (m CentsMap) Get(k string) Cents { return m.Vals[k] }

func (m CentsMap) Sum() Cents {
	var s Cents
	for _, k := range m.Keys {
		s += m.Vals[k]
	}
	return s
}

func (m CentsMap) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range m.Keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		b.Write(kb)
		fmt.Fprintf(&b, ":%d", m.Vals[k])
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// OMap es un objeto JSON con claves en orden (vistas del oráculo con claves de model.json).
type OMap []KV

// KV es una pareja clave-valor de OMap.
type KV struct {
	K string
	V any
}

func (o OMap) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, kv := range o {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(kv.K)
		vb, err := json.Marshal(kv.V)
		if err != nil {
			return nil, err
		}
		b.Write(kb)
		b.WriteByte(':')
		b.Write(vb)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// ---------- Decodificación estricta ----------

func decodeStrict(b []byte, v any) error {
	if len(bytes.TrimSpace(b)) == 0 || string(bytes.TrimSpace(b)) == "null" {
		b = []byte("{}")
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		var me *Error
		if errors.As(err, &me) {
			return me
		}
		if strings.Contains(err.Error(), "unknown field") {
			return errf(CodeUnknownField, "%v", err)
		}
		if strings.Contains(err.Error(), "parsing time") {
			return errf(CodeInvalidDate, "%v", err)
		}
		return errf(CodeInvalidArgument, "%v", err)
	}
	if dec.More() {
		return errf(CodeInvalidArgument, "contenido sobrante después del objeto")
	}
	return nil
}
