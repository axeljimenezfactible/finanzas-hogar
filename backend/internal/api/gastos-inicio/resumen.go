package gastos

import (
	"net/http"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/motor_fin"
)

// previousPeriod es el periodo anterior equivalente: mes, quincena o año calendario anterior; cualquier otro rango
// (semana ISO o rango libre) se recorre su propio número de días hacia atrás.
func previousPeriod(p dominio.Period) dominio.Period {
	last := dominio.MonthPeriod(p.Start).End
	switch {
	case p.Start.Day == 1 && p.Start.Month == time.January && p.End.Day == 31 && p.End.Month == time.December && p.End.Year == p.Start.Year:
		return dominio.Period{Kind: p.Kind, Start: dominio.Date{Year: p.Start.Year - 1, Month: time.January, Day: 1}, End: p.Start.AddDays(-1)}
	case p.Start.Day == 1 && p.End == last:
		return dominio.MonthPeriod(p.Start.AddDays(-1))
	case (p.Start.Day == 1 && p.End.Day == 15 && p.End.Month == p.Start.Month) || (p.Start.Day == 16 && p.End == last):
		f := dominio.FortnightOf(p.Start.AddDays(-1))
		return dominio.Period{Kind: p.Kind, Start: f.Start(), End: f.End()}
	}
	n := p.Days()
	return dominio.Period{Kind: p.Kind, Start: p.Start.AddDays(-n), End: p.Start.AddDays(-1)}
}

func mdate(d dominio.Date) motor_fin.Date { return motor_fin.NewDate(d.Year, d.Month, d.Day) }

// monthsTouched cuenta los meses calendario que toca el periodo (PROVISIONAL, bloqueo S6-gastos-inicio §3).
func monthsTouched(p dominio.Period) int {
	return (p.End.Year-p.Start.Year)*12 + int(p.End.Month) - int(p.Start.Month) + 1
}

func getExpenseSummary(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		f, err := parseFilter(reg, r, false)
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		cur, err := dominio.SumExpenses(ctx, db, f)
		if err != nil {
			return err
		}
		per := *f.Period
		prev := previousPeriod(per)
		fp := f
		fp.Period = &prev
		old, err := dominio.SumExpenses(ctx, db, fp)
		if err != nil {
			return err
		}
		out := api.ExpenseSummary{PeriodStart: per.Start.String(), PeriodEnd: per.End.String(), HouseholdShareCents: cur.HouseholdCents,
			PreviousHouseholdShareCents: old.HouseholdCents, DifferenceCents: cur.HouseholdCents - old.HouseholdCents}
		// Promedios del motor S2 (K20, K39, RN-27): diario HALF_UP, semanal = total × 7 ÷ días con un solo redondeo.
		avg, err := motor_fin.Averages(motor_fin.Cents(cur.HouseholdCents), mdate(per.Start), mdate(per.End))
		if err != nil {
			return err
		}
		out.DailyAverageCents, out.WeeklyAverageCents = api.Ptr(int64(avg.DailyCents)), api.Ptr(int64(avg.WeeklyCents))
		mon, err := motor_fin.MonthlyAverage(motor_fin.Cents(cur.HouseholdCents), monthsTouched(per))
		if err != nil {
			return err
		}
		out.MonthlyAverageCents = api.Ptr(int64(mon))
		if id, cents, ok, err := dominio.TopHouseholdCategory(ctx, db, f); err != nil {
			return err
		} else if ok {
			out.TopCategoryID, out.TopCategoryHouseholdShareCents = api.Ptr(id), api.Ptr(cents)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}
