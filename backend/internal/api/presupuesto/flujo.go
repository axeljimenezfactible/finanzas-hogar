package presupuesto

import (
	"database/sql"
	"net/http"
	"sort"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
)

func (h *handlers) getCashFlowFortnights(w http.ResponseWriter, r *http.Request) error {
	fe, in, err := h.engine(r.Context())
	if err != nil {
		return err
	}
	t := h.today()
	start := dominio.EngineDate(dominio.MonthStart(t.Year, t.Month))
	rows, err := mf.FortnightFlow(in, mf.FortnightOf(start), 24)
	if err != nil {
		return err
	}
	out := api.CashFlowFortnights{ProjectionMode: projectionMode(fe), Rows: []api.CashFlowFortnightsRowsItem{}}
	td := dominio.EngineDate(t)
	for _, x := range rows {
		it := api.CashFlowFortnightsRowsItem{StartDate: x.Start.String(), EndDate: x.End.String(), PayDate: x.PayDate.String(),
			IncomeCents: cents(x.Income.TotalCents), IncomeIsExpected: x.Income.IncomeIsExpected,
			CreditPaymentsCents: cents(x.CreditPayments.HouseholdCents), FixedAndUpcomingCents: cents(x.FixedUpcoming.PlannedCents),
			PlannedSavingsCents: centsPtr(x.SavingsPlannedCents), BalanceCents: cents(x.BalanceCents),
			CumulativeBalanceCents: cents(x.CumulativeBalanceCents), IsActual: x.End.Before(td)}
		// PROVISIONAL (bloqueo S2-flujo-quincenal, pregunta 1): el motor entrega el gasto variable estimado en null; va 0.
		if x.VariableEstimatedCents != nil {
			it.VariableEstimateCents = cents(*x.VariableEstimatedCents)
		}
		out.Rows = append(out.Rows, it)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// isCanonical: meses con cifra canon en model.json (oct-2026 a mar-2027, C1 AccrualVsCash).
func isCanonical(m string) bool { return m >= "2026-10" && m <= "2027-03" }

func (h *handlers) getAccrualVsCash(w http.ResponseWriter, r *http.Request) error {
	fe, in, err := h.engine(r.Context())
	if err != nil {
		return err
	}
	pr, err := mf.Projection(in, dominio.EngineMonth(h.today()))
	if err != nil {
		return err
	}
	out := api.AccrualVsCash{ProjectionMode: projectionMode(fe), Months: []api.AccrualVsCashMonthsItem{}}
	for _, pm := range pr.Months {
		// Acumulados con el mismo origen (Q2-65 A): el motor da el devengado acumulado al cierre del mes (LiveDebt) y la
		// diferencia acumulada de la proyección (periodos[i].diferencia_acumulada); el efectivo acumulado es su resta.
		ld, err := mf.LiveDebt(in, pm.Month.Last())
		if err != nil {
			return err
		}
		if ld.CumulativeDiffCents != pm.CumulativeDiffCents {
			h.reg.Deps().Log.Warn("presupuesto: diferencia acumulada distinta entre LiveDebt y Projection", "month", pm.Month.String())
		}
		it := api.AccrualVsCashMonthsItem{Month: pm.Month.String(), AccruedCents: cents(pm.Accrued.TotalCents), CashCents: cents(pm.Cash.TotalCents),
			AccruedHouseholdCents: cents(pm.Accrued.ByPerson.Get(fe.Holder)), CashHouseholdCents: cents(pm.Cash.ByPerson.Get(fe.Holder)),
			CumulativeAccruedCents: cents(ld.AccruedToDateCents), CumulativeDifferenceCents: cents(pm.CumulativeDiffCents),
			IsCanonical: isCanonical(pm.Month.String())}
		it.CumulativeCashCents = it.CumulativeAccruedCents - it.CumulativeDifferenceCents
		for _, x := range fe.Externals {
			it.AccruedExternalCents += cents(pm.Accrued.ByPerson.Get(x))
			it.CashExternalCents += cents(pm.Cash.ByPerson.Get(x))
		}
		out.Months = append(out.Months, it)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func optID(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return api.Ptr(n.Int64)
}

func optStr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return api.Ptr(s.String)
}

func (h *handlers) getCashFlowAssumptions(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	fe, err := dominio.LoadFinEngine(ctx, h.db())
	if err != nil {
		return err
	}
	items, err := dominio.ActivePlannedItems(ctx, h.db())
	if err != nil {
		return err
	}
	out := api.CashFlowAssumptions{ProjectionMode: projectionMode(fe), PlannedItems: []api.PlannedItem{},
		// PROVISIONAL (bloqueo S2-flujo-quincenal, pregunta 1, Q2-12): el motor S2 aún no fija el método; sin monto.
		VariableEstimateMethod: "Sin método definido todavía: el gasto variable estimado está pendiente de decisión (PROVISIONAL)."}
	for _, p := range items {
		out.PlannedItems = append(out.PlannedItems, api.PlannedItem{ID: p.ID, Kind: p.Kind, Description: p.Description, AmountCents: p.AmountCents,
			SecondHalfAmountCents: optID(p.SecondHalfAmountCents), Frequency: p.Frequency, IntervalCount: p.IntervalCount, StartDate: p.StartDate,
			EndDate: optStr(p.EndDate), CategoryID: optID(p.CategoryID), PersonID: optID(p.PersonID), CostShareRuleID: optID(p.CostShareRuleID),
			PaymentMethod: optStr(p.PaymentMethod), CreditID: optID(p.CreditID), TermCount: optID(p.TermCount), AccrualMonths: p.AccrualMonths,
			IsActive: p.IsActive, Notes: optStr(p.Notes), CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt)})
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) getLiveDebt(w http.ResponseWriter, r *http.Request) error {
	d, err := api.QueryDate(r, "date")
	if err != nil {
		return err
	}
	t := h.today()
	// PROVISIONAL (bloqueo S6-presupuesto-motor, pregunta 3): «fin del horizonte proyectado» = último día del mes en curso + 11.
	cut := dominio.MonthEnd(dominio.MonthStart(t.Year, t.Month+11))
	if d != nil {
		cut = *d
	}
	fe, in, err := h.engine(r.Context())
	if err != nil {
		return err
	}
	ld, err := mf.LiveDebt(in, dominio.EngineDate(cut))
	if err != nil {
		return err
	}
	out := api.LiveDebt{Date: cut.String(), TotalCents: cents(ld.LiveDebtCents), Credits: []api.LiveDebtCreditsItem{}, Purchases: []api.LiveDebtPurchasesItem{}}
	byCredit := map[int64]int64{}
	for _, p := range ld.ByPurchase {
		if p.PendingCents == 0 {
			continue
		}
		c := fe.PurchaseCredit[p.PurchaseID]
		out.Purchases = append(out.Purchases, api.LiveDebtPurchasesItem{CreditPurchaseID: p.PurchaseID, CreditID: c, TotalCents: cents(p.PendingCents)})
		byCredit[c] += cents(p.PendingCents)
	}
	for c, v := range byCredit {
		out.Credits = append(out.Credits, api.LiveDebtCreditsItem{CreditID: c, TotalCents: v})
	}
	sort.Slice(out.Credits, func(a, b int) bool { return out.Credits[a].CreditID < out.Credits[b].CreditID })
	return api.WriteJSON(w, http.StatusOK, out)
}
