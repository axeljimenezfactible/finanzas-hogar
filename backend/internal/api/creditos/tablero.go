package creditos

import (
	"context"
	"fmt"

	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// Pagos de crédito por mes para el tablero de Inicio (K54; RN-30, métricas M-20 y M-21 de C1 v1.5). Es la misma
// cuenta que getCreditTotals (month_household_cents y month_total_cents): calendario unificado del motor S3 con las
// líneas activas; parte del hogar = partes de las personas que no son externas. Aquí no se recalcula nada: solo se
// agrupa por mes lo que entrega el motor. Agregado por correccion-especificacion-r2 (hallazgo A-1 de
// REVISION-especificacion-r2) para que el paquete tablero use exactamente las cifras de S3.

// MonthPayments son los pagos programados con fecha en un mes.
type MonthPayments struct {
	HouseholdCents int64 // parte del hogar (sin Karen)
	TotalCents     int64 // con la parte de Karen
	Count          int64 // pagos programados con fecha en el mes
}

// PaymentsByMonth agrupa por mes («AAAA-MM») los pagos del calendario unificado entre from y to (incluidos).
// hasCredits dice si hay al menos una línea de crédito activa.
func PaymentsByMonth(ctx context.Context, q dominio.Querier, today, from, to dominio.Date) (byMonth map[string]MonthPayments, hasCredits bool, err error) {
	s, err := load(ctx, q, today)
	if err != nil {
		return nil, false, err
	}
	states := s.states(nil)
	out := map[string]MonthPayments{}
	if len(states) == 0 {
		return out, false, nil
	}
	cal, err := mc.UnifiedCalendar(states, mdateOf(from), mdateOf(to))
	if err != nil {
		return nil, true, motorErr(err, "period")
	}
	for _, e := range cal.Entries {
		k := fmt.Sprintf("%04d-%02d", e.Date.Year, int(e.Date.Month))
		mp := out[k]
		mp.Count++
		for _, sh := range e.Shares {
			mp.TotalCents += int64(sh.Amount)
			if !s.isExternal(sh.PersonID) {
				mp.HouseholdCents += int64(sh.Amount)
			}
		}
		out[k] = mp
	}
	return out, true, nil
}
