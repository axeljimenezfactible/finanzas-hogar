package reportes

import (
	"strconv"

	"finanzas-hogar/internal/dominio"
)

// Nota fija de los PDF. El contenido de los reportes es PROVISIONAL hasta la puerta 4 (K31).
const notaPropuesta = "Contenido propuesto por el sistema; Axel lo aprueba en la puerta 4. El plan de comidas y la canasta se agregan cuando su cálculo compartido esté listo."

// MonthlyDoc pinta el reporte mensual. Cada cifra y fecha cita su campo de getMonthlyReport («monthly:AAAA-MM») o del
// renglón del mes de getAccrualVsCash («accrual:AAAA-MM»).
func MonthlyDoc(m *Monthly) *Doc {
	r := &m.Report
	k := "monthly:" + r.Month + "|"
	mo, _ := dominio.ParseDate(r.Month + "-01")
	d := &Doc{Title: "Reporte mensual de " + MesLargo(mo.Year, mo.Month), Footer: "Finanzas Hogar · reporte mensual", Note: notaPropuesta}
	if len(r.Fortnights.Fortnights) == 2 {
		d.Subtitle = []Cell{txt("Del "), date(k+"fortnights.fortnights.0.start_date", r.Fortnights.Fortnights[0].StartDate, false),
			txt(" al "), date(k+"fortnights.fortnights.1.end_date", r.Fortnights.Fortnights[1].EndDate, false)}
	}
	b := r.Budget
	res := Section{Title: "Resumen del mes", Spec: "LR", Cols: []string{"Concepto", "Importe"}, Rows: [][]Cell{
		{txt("Ingresos"), money(k+"budget.income_cents", b.IncomeCents)},
		{txt("Gasto planeado"), money(k+"budget.planned_spending_cents", b.PlannedSpendingCents)},
		{txt("Ahorro planeado"), money(k+"budget.planned_savings_cents", b.PlannedSavingsCents)},
		{txt("Sin asignar"), money(k+"budget.unassigned_cents", b.UnassignedCents)},
		{txt("Gastado"), money(k+"budget.actual_cents", b.ActualCents)},
		{txt("Restante"), money(k+"budget.remaining_cents", b.RemainingCents)},
	}}
	if b.UsedPct != nil {
		res.Rows = append(res.Rows, []Cell{txt("Presupuesto usado"), txt(strconv.FormatInt(*b.UsedPct, 10) + " %")})
	}
	d.Sections = append(d.Sections, res)
	a := m.Accrual
	ka := "accrual:" + r.Month + "|"
	d.Sections = append(d.Sections, Section{Title: "Devengado contra efectivo", Spec: "LR", Cols: []string{"Concepto", "Importe"},
		Note: "Devengado: lo que se consumió en el mes. Efectivo: lo que se pagó. La diferencia acumulada es lo consumido que falta pagar.",
		Rows: [][]Cell{
			{txt("Devengado del mes"), money(ka+"accrued_cents", a.AccruedCents)},
			{txt("Efectivo del mes"), money(ka+"cash_cents", a.CashCents)},
			{txt("Diferencia acumulada al cierre"), money(ka+"cumulative_difference_cents", a.CumulativeDifferenceCents)},
		}})
	cat := Section{Title: "Por categoría", Spec: "LRRRR", Cols: []string{"Categoría", "Planeado", "Real", "Diferencia", "Usado"}}
	for i, row := range r.Categories.Rows {
		p := k + "categories.rows." + strconv.Itoa(i) + "."
		c := []Cell{txt(m.Categories[row.CategoryID]), txt("Sin plan"), money(p+"actual_cents", row.ActualCents), txt("—"), txt("—")}
		if row.PlannedCents != nil {
			c[1] = money(p+"planned_cents", *row.PlannedCents)
		}
		if row.DifferenceCents != nil {
			c[3] = money(p+"difference_cents", *row.DifferenceCents)
		}
		if row.UsedPct != nil {
			c[4] = txt(strconv.FormatInt(*row.UsedPct, 10) + " %")
		}
		cat.Rows = append(cat.Rows, c)
	}
	if len(cat.Rows) > 0 {
		d.Sections = append(d.Sections, cat)
	}
	q := Section{Title: "Por quincena", Spec: "LLLLRR", Cols: []string{"Quincena", "Del", "Al", "Día de pago", "Real", "Planeado"}}
	for i, f := range r.Fortnights.Fortnights {
		p := k + "fortnights.fortnights." + strconv.Itoa(i) + "."
		name := "Primera"
		if f.Half == 2 {
			name = "Segunda"
		}
		q.Rows = append(q.Rows, []Cell{txt(name), date(p+"start_date", f.StartDate, false), date(p+"end_date", f.EndDate, false),
			date(p+"pay_date", f.PayDate, true), money(p+"actual_cents", f.ActualCents), money(p+"planned_cents", f.PlannedCents)})
	}
	d.Sections = append(d.Sections, q)
	c := r.Credits
	d.Sections = append(d.Sections, Section{Title: "Créditos", Spec: "LR", Cols: []string{"Concepto", "Importe"}, Rows: [][]Cell{
		{txt("Pagos programados del mes"), money(k+"credits.month_total_cents", c.MonthTotalCents)},
		{txt("Parte del hogar"), money(k+"credits.month_household_cents", c.MonthHouseholdCents)},
		{txt("Parte de personas externas"), money(k+"credits.month_external_cents", c.MonthExternalCents)},
		{txt("Pagos de la segunda quincena"), money(k+"credits.fortnight_total_cents", c.FortnightTotalCents)},
		{txt("Intereses pagados en el mes"), money(k+"credits.interest_paid_month_cents", c.InterestPaidMonthCents)},
		{txt("Intereses pagados en el año"), money(k+"credits.interest_paid_year_cents", c.InterestPaidYearCents)},
		{txt("Deuda viva al cierre del mes"), money(k+"live_debt_end_cents", r.LiveDebtEndCents)},
	}})
	for i, pp := range r.ExternalParticipation {
		p := k + "external_participation." + strconv.Itoa(i) + "."
		who := m.Persons[pp.PersonID]
		d.Sections = append(d.Sections, Section{Title: "Participación de " + who, Spec: "LR", Cols: []string{"Concepto", "Importe"},
			Rows: [][]Cell{{txt("Despensa devengada en el mes"), money(p+"accrued_groceries_cents", pp.AccruedGroceriesCents)},
				{txt("Saldo pendiente de su parte"), money(p+"live_debt_cents", pp.LiveDebtCents)}}})
		if len(pp.PaymentShares) > 0 {
			s := Section{Title: "Pagos de " + who + " en el mes", Spec: "LRL", Cols: []string{"Fecha de pago", "Monto", "Estado"}}
			for j, ps := range pp.PaymentShares {
				pj := p + "payment_shares." + strconv.Itoa(j) + "."
				s.Rows = append(s.Rows, []Cell{date(pj+"due_date", ps.DueDate, true), money(pj+"amount_cents", ps.AmountCents), txt(statusText(ps.Status))})
			}
			d.Sections = append(d.Sections, s)
		}
		if len(pp.Purchases) > 0 {
			s := Section{Title: "Compras a crédito con parte de " + who, Spec: "LLRLR", Cols: []string{"Compra", "Fecha", "Su total", "Pagos", "Saldo"}}
			for j, pu := range pp.Purchases {
				pj := p + "purchases." + strconv.Itoa(j) + "."
				label := "Compra a crédito"
				if pu.Description != nil {
					label = *pu.Description
				}
				s.Rows = append(s.Rows, []Cell{txt(label), date(pj+"purchase_date", pu.PurchaseDate, false), money(pj+"total_cents", pu.TotalCents),
					txt(strconv.FormatInt(pu.PaymentsMade, 10) + " de " + strconv.FormatInt(pu.PaymentsTotal, 10)), money(pj+"balance_cents", pu.BalanceCents)})
			}
			d.Sections = append(d.Sections, s)
		}
	}
	return d
}

// ExternalDoc pinta el estado de cuenta: solo despensa en común, sus compras particulares y sus pagos (S7-03).
func ExternalDoc(x *External) *Doc {
	st := &x.Statement
	k := "statement|"
	d := &Doc{Title: "Estado de cuenta de " + x.Name, Footer: "Finanzas Hogar · estado de cuenta",
		Note:     "Paga su parte directamente al acreedor. Este estado de cuenta no se le envía automáticamente: Axel lo reenvía a mano.",
		Subtitle: []Cell{txt("Del "), date(k+"period_start_date", st.PeriodStartDate, false), txt(" al "), date(k+"period_end_date", st.PeriodEndDate, false)}}
	g := Section{Title: "Su parte de la despensa en común", Spec: "LRRR", Cols: []string{"Fecha de compra", "Capital", "Intereses", "Total"}}
	for i, s := range st.GroceryShares {
		p := k + "grocery_shares." + strconv.Itoa(i) + "."
		g.Rows = append(g.Rows, []Cell{date(p+"purchase_date", s.PurchaseDate, false), money(p+"principal_cents", s.PrincipalCents),
			money(p+"interest_cents", s.InterestCents), money(p+"total_cents", s.TotalCents)})
	}
	if len(g.Rows) > 0 {
		d.Sections = append(d.Sections, g)
	}
	pc := Section{Title: "Sus compras a crédito", Spec: "LLR", Cols: []string{"Compra", "Fecha", "Total"}}
	for i, p := range st.Purchases {
		pp := k + "purchases." + strconv.Itoa(i) + "."
		pc.Rows = append(pc.Rows, []Cell{txt(x.Labels[i]), date(pp+"purchase_date", p.PurchaseDate, false), money(pp+"total_cents", p.TotalCents)})
	}
	if len(pc.Rows) > 0 {
		d.Sections = append(d.Sections, pc)
	}
	ps := Section{Title: "Sus pagos", Note: "Pagos del periodo y pendientes posteriores.", Spec: "LRRL", Cols: []string{"Fecha de pago", "Monto", "Pagado", "Estado"}}
	for i, s := range st.PaymentShares {
		p := k + "payment_shares." + strconv.Itoa(i) + "."
		ps.Rows = append(ps.Rows, []Cell{date(x.DueRefs[i], x.DueDates[i], true), money(p+"amount_cents", s.AmountCents),
			money(p+"paid_cents", s.PaidCents), txt(statusText(s.Status))})
	}
	if len(ps.Rows) > 0 {
		d.Sections = append(d.Sections, ps)
	}
	d.Sections = append(d.Sections, Section{Title: "Resumen", Spec: "LR", Cols: []string{"Concepto", "Importe"}, Rows: [][]Cell{
		{txt("Pagado en el periodo"), money(k+"paid_in_period_cents", st.PaidInPeriodCents)},
		{txt("Saldo pendiente"), money(k+"pending_total_cents", st.PendingTotalCents)}}})
	return d
}
