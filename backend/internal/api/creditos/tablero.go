package creditos

import (
	"context"
	"fmt"

	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// Pagos de crédito por mes para el tablero de Inicio (K54; RN-30, métricas M-20 y M-21 de C1 v1.5). Es la misma
// cuenta que getCreditTotals (month_household_cents y month_total_cents), pero con todas las líneas, activas o
// archivadas (RN-30 «de todas las líneas», igual que el efectivo de S2; PROVISIONAL, bloqueo
// S6-tablero-lineas-archivadas, opción a, aplicada por backend-cierre con la decisión 2 del orquestador de WF-3c bis):
// calendario unificado del motor S3; parte del hogar = partes de las personas que no son externas. Aquí no se
// recalcula nada: solo se agrupa por mes lo que entrega el motor. Agregado por correccion-especificacion-r2 (hallazgo A-1 de
// REVISION-especificacion-r2) para que el paquete tablero use exactamente las cifras de S3.

// MonthPayments son los pagos programados con fecha en un mes.
type MonthPayments struct {
	HouseholdCents int64 // parte del hogar (sin Karen)
	TotalCents     int64 // con la parte de Karen
	Count          int64 // pagos programados con fecha en el mes
}

// PaymentsByMonth agrupa por mes («AAAA-MM») los pagos del calendario unificado entre from y to (incluidos).
// hasCredits dice si hay al menos una línea de crédito activa o algún pago (de una línea archivada) con fecha entre
// from y to; con false no hay pagos (PROVISIONAL, bloqueo S6-tablero-lineas-archivadas, opción a1).
func PaymentsByMonth(ctx context.Context, q dominio.Querier, today, from, to dominio.Date) (byMonth map[string]MonthPayments, hasCredits bool, err error) {
	s, err := load(ctx, q, today)
	if err != nil {
		return nil, false, err
	}
	var states []mc.CreditState
	for _, c := range s.credits {
		states = append(states, s.state(c))
		hasCredits = hasCredits || c.IsActive
	}
	out := map[string]MonthPayments{}
	if len(states) == 0 {
		return out, false, nil
	}
	cal, err := mc.UnifiedCalendar(states, mdateOf(from), mdateOf(to))
	if err != nil {
		return nil, hasCredits, motorErr(err, "period")
	}
	hasCredits = hasCredits || len(cal.Entries) > 0
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
	return out, hasCredits, nil
}
