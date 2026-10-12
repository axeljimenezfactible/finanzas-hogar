package gastos

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/comidas"
	"finanzas-hogar/internal/api/comprasa"
	"finanzas-hogar/internal/dominio"
	mf "finanzas-hogar/internal/motor_fin"
)

// Tarjetas de Inicio (U3, N-020 a N-026), tendencia mensual (N-044) y participación de Karen. Todas las cifras
// salen del motor S2 (motor_fin) con la entrada que arma dominio.LoadFinEngine; aquí no se recalcula nada a mano.

func today(reg *api.Registry) dominio.Date { return dominio.Today(reg.Deps().Clock) }

// refDate es el día de referencia de una tarjeta: el inicio del filtro period o, sin filtro, hoy
// (PROVISIONAL, bloqueo S6-gastos-inicio §10).
func refDate(reg *api.Registry, r *http.Request) (dominio.Date, error) {
	p, err := api.QueryPeriod(r)
	if err != nil || p == nil {
		return today(reg), err
	}
	return p.Start, nil
}

// engineAt carga la entrada del motor S2 con «hoy» = el reloj de la API.
func engineAt(reg *api.Registry, r *http.Request) (*dominio.FinEngine, mf.Input, error) {
	fe, err := dominio.LoadFinEngine(r.Context(), reg.Deps().DB)
	if err != nil {
		return nil, mf.Input{}, err
	}
	in, err := fe.At(today(reg), nil)
	if err != nil {
		return nil, mf.Input{}, fmt.Errorf("motor_fin: entrada inválida: %w", err)
	}
	return fe, in, nil
}

func getDashboardMonthVsPlan(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ref, err := refDate(reg, r)
		if err != nil {
			return err
		}
		fe, in, err := engineAt(reg, r)
		if err != nil {
			return err
		}
		em := dominio.EngineMonth(ref)
		rows, err := mf.PlanVsActual(in, em, 1, "")
		if err != nil {
			return err
		}
		out := api.DashboardMonthVsPlan{Month: em.String(), ActualCents: int64(rows[0].ActualCents)}
		if pl := rows[0].PlanCents; pl != nil && *pl > 0 {
			plan := int64(*pl)
			pct := dominio.PercentOf(out.ActualCents, plan)
			st := "normal"
			if pct >= fe.LimitPct {
				st = "exceeded"
			} else if pct >= fe.WarningPct {
				st = "warning"
			}
			out.PlannedCents, out.UsedPct, out.RemainingCents, out.Status = api.Ptr(plan), api.Ptr(pct), api.Ptr(plan-out.ActualCents), api.Ptr(st)
		}
		// Días que faltan contando hoy: mes en curso = hasta el último día; mes pasado = 0; mes futuro = todos.
		t, first, last := today(reg), dominio.MonthStart(ref.Year, ref.Month), dominio.MonthEnd(ref)
		switch {
		case t.After(last):
			out.DaysLeft = 0
		case t.Before(first):
			out.DaysLeft = int64(first.DaysUntil(last) + 1)
		default:
			out.DaysLeft = int64(t.DaysUntil(last) + 1)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func getDashboardFortnightFlow(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ref, err := refDate(reg, r)
		if err != nil {
			return err
		}
		_, in, err := engineAt(reg, r)
		if err != nil {
			return err
		}
		rows, err := mf.FortnightFlow(in, mf.FortnightOf(dominio.EngineDate(ref)), 1)
		if err != nil {
			return err
		}
		x := rows[0]
		return api.WriteJSON(w, http.StatusOK, api.DashboardFortnightFlow{StartDate: x.Start.String(), EndDate: x.End.String(),
			IncomeCents: int64(x.Income.TotalCents), IncomeIsExpected: x.Income.IncomeIsExpected, SpentCents: int64(x.SpentCents),
			CommittedCents: int64(x.CommittedCents), ClosingBalanceCents: int64(x.BalanceCents)})
	}
}

func getDashboardSavings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		_, in, err := engineAt(reg, r)
		if err != nil {
			return err
		}
		sv, err := mf.SavingsGoals(in)
		if err != nil {
			return err
		}
		ms := sv.MonthSummary
		out := api.DashboardSavings{SavedTotalCents: int64(ms.SavedCents), ContributedMonthCents: int64(ms.ContributedThisMonthCents),
			PlannedMonthCents: int64(ms.PlannedThisMonthCents)}
		g, ok, err := dominio.GastosNextSavingsGoal(r.Context(), reg.Deps().DB)
		if err != nil {
			return err
		}
		if ok {
			it := api.SavingsGoal{ID: g.ID, Name: g.Name, Kind: g.Kind, TargetCents: nullInt(g.TargetCents), TargetDate: nullStr(g.TargetDate),
				PlannedMonthlyCents: nullInt(g.PlannedMonthlyCents), StartDate: g.StartDate, IsActive: g.IsActive, Notes: nullStr(g.Notes),
				CreatedAt: api.APITime(g.CreatedAt), UpdatedAt: api.APITime(g.UpdatedAt)}
			for _, s := range sv.Goals {
				if s.GoalID == g.ID {
					it.SavedCents, it.ProgressTenthsPct, it.ProgressPct, it.Status = api.Ptr(int64(s.SavedCents)), s.ProgressTenthsPct, s.ProgressPct, api.Ptr(s.State)
					if s.NeededMonthlyCents != nil {
						it.RequiredMonthlyCents = api.Ptr(int64(*s.NeededMonthlyCents))
					}
				}
			}
			out.NextGoal = &it
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// accruedFor pide al motor el devengado de un solo gasto (con su compra a crédito) en los meses dados.
func accruedFor(fe *dominio.FinEngine, t dominio.Date, id int64, months []mf.Month) ([]mf.PeriodSummary, error) {
	in, err := fe.Only(t, id)
	if err != nil {
		return nil, fmt.Errorf("motor_fin: entrada inválida: %w", err)
	}
	return mf.Periods(in, mf.PeriodSpec{Months: months})
}

func getExpenseMonthlyTrend(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		f, err := parseFilter(reg, r, false)
		if err != nil {
			return err
		}
		// Solo aplican los filtros de C1 (persona, familia, crédito); el rango es fijo: los 12 meses hasta el mes en curso.
		f.Period, f.CategoryIDs, f.StoreIDs, f.Q = nil, nil, nil, ""
		ctx, db := r.Context(), reg.Deps().DB
		ids, err := dominio.GastosExpenseIDs(ctx, db, f)
		if err != nil {
			return err
		}
		fe, err := dominio.LoadFinEngine(ctx, db)
		if err != nil {
			return err
		}
		t := today(reg)
		cur := dominio.EngineMonth(t)
		months := make([]mf.Month, 12)
		for i := range months {
			months[i] = cur.Add(i - 11)
		}
		acc := make([]map[int64]int64, 12) // mes → categoría (0 = sin id) → parte del hogar devengada
		for i := range acc {
			acc[i] = map[int64]int64{}
		}
		total := map[int64]int64{}
		for _, e := range fe.Expenses() {
			if !ids[e.ID] || e.ExpenseDate.After(cur.Last()) {
				continue
			}
			ps, err := accruedFor(fe, t, e.ID, months)
			if err != nil {
				return err
			}
			cat := fe.CategoryID[e.Category]
			for i, p := range ps {
				if v := int64(p.Accrued.ByPerson.Get(fe.Holder)); v != 0 {
					acc[i][cat] += v
					total[cat] += v
				}
			}
		}
		// Las 6 categorías con más gasto del rango; el resto (y lo que no tiene categoría) va a «Otras».
		cats := make([]int64, 0, len(total))
		for c := range total {
			if c != 0 {
				cats = append(cats, c)
			}
		}
		slices.SortFunc(cats, func(a, b int64) int { return cmp.Or(cmp.Compare(total[b], total[a]), cmp.Compare(a, b)) })
		top := cats[:min(6, len(cats))]
		others := len(cats) > len(top) || total[0] != 0
		first := 11
		for i := range acc {
			if len(acc[i]) > 0 {
				first = i
				break
			}
		}
		out := api.ExpenseMonthlyTrend{Months: []api.ExpenseMonthlyTrendMonthsItem{}}
		for i := first; i < 12; i++ {
			it := api.ExpenseMonthlyTrendMonthsItem{Month: months[i].String(), Categories: []api.ExpenseMonthlyTrendMonthsItemCategoriesItem{}}
			rest := acc[i][0]
			for c, v := range acc[i] {
				if c != 0 && !slices.Contains(top, c) {
					rest += v
				}
			}
			for _, c := range top {
				it.Categories = append(it.Categories, api.ExpenseMonthlyTrendMonthsItemCategoriesItem{CategoryID: api.Ptr(c), HouseholdShareCents: acc[i][c]})
			}
			if others {
				it.Categories = append(it.Categories, api.ExpenseMonthlyTrendMonthsItemCategoriesItem{HouseholdShareCents: rest})
			}
			out.Months = append(out.Months, it)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func getPersonParticipation(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		pid, err := api.PathID(r, "person_id")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		if err := dominio.GastosIsExternalPerson(ctx, db, pid); err != nil {
			return err
		}
		fe, err := dominio.LoadFinEngine(ctx, db)
		if err != nil {
			return err
		}
		t := today(reg)
		em := dominio.EngineMonth(t)
		out := api.PersonParticipation{PersonID: pid, Month: em.String(), Purchases: []api.PersonParticipationPurchasesItem{},
			PaymentShares: []api.PersonParticipationPaymentSharesItem{}}
		// (1) Su parte devengada de la despensa en el mes (motor S2), gasto por gasto como getAccruedGroceries.
		code := fe.PersonCode[pid]
		for _, e := range fe.Expenses() {
			if e.Category != dominio.GroceriesCategoryCode || e.ExpenseDate.After(em.Last()) {
				continue
			}
			ps, err := accruedFor(fe, t, e.ID, []mf.Month{em})
			if err != nil {
				return err
			}
			out.AccruedGroceriesCents += int64(ps[0].Accrued.ByPerson.Get(code))
		}
		// (3) y (4) Sus partes de pagos: las del mes con estado; deuda viva = Σ de lo pendiente (RN-04).
		shares, err := dominio.GastosPersonShares(ctx, db, pid)
		if err != nil {
			return err
		}
		first, last := em.First().String(), em.Last().String()
		for _, s := range shares {
			out.LiveDebtCents += s.PendingCents()
			if s.DueDate >= first && s.DueDate <= last {
				st := "pending"
				if s.Paid() {
					st = "paid"
				}
				out.PaymentShares = append(out.PaymentShares, api.PersonParticipationPaymentSharesItem{ScheduledPaymentID: s.ScheduledPaymentID,
					DueDate: s.DueDate, AmountCents: s.AmountCents, Status: st})
			}
		}
		// (2) Sus compras particulares vivas: pagos hechos de cuántos y saldo (lo que le falta).
		ps, err := dominio.GastosParticularPurchases(ctx, db, pid)
		if err != nil {
			return err
		}
		for _, p := range ps {
			it := api.PersonParticipationPurchasesItem{CreditPurchaseID: p.ID, Description: nullStr(p.Description), PurchaseDate: p.PurchaseDate, TotalCents: p.TotalCents}
			for _, s := range shares {
				if s.PurchaseID == p.ID {
					it.PaymentsTotal++
					it.BalanceCents += s.PendingCents()
					if s.Paid() {
						it.PaymentsMade++
					}
				}
			}
			out.Purchases = append(out.Purchases, it)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

var searchTypes = []string{"expense", "income", "recipe", "product", "store", "credit", "credit_purchase", "savings_goal", "settings_section"}

func search(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		types, err := queryEnum(r, "type", searchTypes)
		if err != nil {
			return err
		}
		text := strings.TrimSpace(r.URL.Query().Get("q"))
		if len([]rune(text)) > 200 {
			return api.BadRequest("El texto de búsqueda admite hasta 200 caracteres.")
		}
		out := api.SearchResult{Groups: []api.SearchResultGroupsItem{}}
		if text == "" { // texto vacío: grupos vacíos (los recientes viven en el dispositivo, N-003)
			return api.WriteJSON(w, http.StatusOK, out)
		}
		if len(types) == 0 {
			types = searchTypes
		}
		for _, typ := range searchTypes {
			if !slices.Contains(types, typ) {
				continue
			}
			hits, err := dominio.GastosSearch(r.Context(), reg.Deps().DB, typ, text)
			if err != nil {
				return err
			}
			if len(hits) == 0 {
				continue
			}
			g := api.SearchResultGroupsItem{Type: typ, Total: int64(len(hits)), Items: []api.SearchResultGroupsItemItemsItem{}}
			for _, h := range hits[:min(5, len(hits))] {
				g.Items = append(g.Items, api.SearchResultGroupsItemItemsItem{ID: h.ID, Title: h.Title, Subtitle: nullStr(h.Subtitle), Date: nullStr(h.Date)})
			}
			out.Groups = append(out.Groups, g)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// getDashboardWeekMeals es la tarjeta «Comidas de la semana» (N-023, M-05). El armado del motor S4 es de
// api-comidas-salud (comidas.WeekMeals, bloqueo S6-gastos-inicio §15); aquí solo se elige la semana: la del
// filtro period (debe ser una semana) o, sin filtro, la semana en curso de lunes a domingo (K8; K38 no aplica a
// esta tarjeta, U3 tarjeta 4).
func getDashboardWeekMeals(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.QueryPeriod(r)
		if err != nil {
			return err
		}
		day := today(reg)
		if p != nil {
			var v api.Validator
			v.Check(p.Kind == dominio.PeriodWeek, "period", "invalid_format", "Elige una semana, por ejemplo 2026-W42.")
			if err := v.Err(); err != nil {
				return err
			}
			day = p.Start
		}
		monday := day.AddDays(-((int(day.Weekday()) + 6) % 7))
		out, err := comidas.WeekMeals(r.Context(), reg, monday)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// getDashboardBestBasket es la tarjeta «Mejor canasta» (N-024, M-06). La canasta y su comparación son de api-compras-a
// (comprasa.BestBasket: misma regla y mismo motor S5 que getBasketComparison; bloqueo S6-gastos-inicio §15); aquí solo
// se lee la llave obligatoria week (R3-01, regla de comprasa.WeekParam): la PWA manda la lista vigente (K38, U3 tarjeta 5).
func getDashboardBestBasket(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		monday, err := comprasa.WeekParam(r)
		if err != nil {
			return err
		}
		out, err := comprasa.BestBasket(r.Context(), reg.Deps().DB, today(reg), monday)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}
