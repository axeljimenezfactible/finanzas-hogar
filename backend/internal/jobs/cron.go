package jobs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/dominio"
)

// Schedule es una expresión cron de 5 campos (minuto hora día-del-mes mes día-de-la-semana) en hora de
// America/Mexico_City (columna job.schedule). Admite *, listas, rangos y pasos (*/15, 1-5, 0,30). Domingo = 0 o 7.
// Si día del mes y día de la semana están restringidos, basta con que coincida uno (regla clásica de cron).
type Schedule struct {
	min, hour, dom, month, dow uint64
	domAny, dowAny             bool
}

// ParseSchedule lee la expresión.
func ParseSchedule(expr string) (Schedule, error) {
	f := strings.Fields(expr)
	if len(f) != 5 {
		return Schedule{}, fmt.Errorf("cron %q: se esperan 5 campos", expr)
	}
	var s Schedule
	var err error
	lims := [5][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 7}}
	dst := [5]*uint64{&s.min, &s.hour, &s.dom, &s.month, &s.dow}
	for i := range f {
		if *dst[i], err = parseField(f[i], lims[i][0], lims[i][1]); err != nil {
			return Schedule{}, fmt.Errorf("cron %q: %w", expr, err)
		}
	}
	if s.dow&(1<<7) != 0 {
		s.dow |= 1
	}
	s.domAny, s.dowAny = f[2] == "*", f[4] == "*"
	return s, nil
}

func parseField(f string, lo, hi int) (uint64, error) {
	var bits uint64
	for _, part := range strings.Split(f, ",") {
		rng, stepS, hasStep := strings.Cut(part, "/")
		step := 1
		if hasStep {
			n, err := strconv.Atoi(stepS)
			if err != nil || n < 1 {
				return 0, fmt.Errorf("paso inválido %q", part)
			}
			step = n
		}
		a, b := lo, hi
		if rng != "*" {
			x, y, isRange := strings.Cut(rng, "-")
			var err1, err2 error
			a, err1 = strconv.Atoi(x)
			b, err2 = a, nil
			if isRange {
				b, err2 = strconv.Atoi(y)
			} else if hasStep {
				b = hi
			}
			if err1 != nil || err2 != nil || a < lo || b > hi || a > b {
				return 0, fmt.Errorf("valor fuera de %d-%d: %q", lo, hi, part)
			}
		}
		for v := a; v <= b; v += step {
			bits |= 1 << uint(v)
		}
	}
	return bits, nil
}

func (s Schedule) matchDay(t time.Time) bool {
	d := s.dom&(1<<uint(t.Day())) != 0
	w := s.dow&(1<<uint(t.Weekday())) != 0
	switch {
	case s.domAny && s.dowAny:
		return true
	case s.domAny:
		return w
	case s.dowAny:
		return d
	}
	return d || w
}

// Next da el primer instante estrictamente posterior a after que cumple la expresión, calculado en America/Mexico_City.
func (s Schedule) Next(after time.Time) time.Time {
	t := after.In(dominio.MX).Truncate(time.Minute).Add(time.Minute)
	limit := t.AddDate(5, 0, 0)
	for t.Before(limit) {
		if s.month&(1<<uint(t.Month())) == 0 {
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, dominio.MX)
			continue
		}
		if !s.matchDay(t) {
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, dominio.MX)
			continue
		}
		if s.hour&(1<<uint(t.Hour())) == 0 {
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, dominio.MX)
			continue
		}
		if s.min&(1<<uint(t.Minute())) == 0 {
			t = t.Add(time.Minute)
			continue
		}
		return t
	}
	return time.Time{}
}
