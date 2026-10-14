package gastos

import (
	"cmp"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios"
)

// Calendario del hogar (N-029, N-086, N-104, N-115, N-192, N-208). eventTypes va en el orden de N-029: es el
// orden dentro de un mismo día.
var eventTypes = []string{"income", "credit_payment", "fixed_expense", "upcoming_expense", "card_statement", "planned_grocery",
	"market_visit", "promotion", "savings_goal", "savings_movement", "meal_plan", "monthly_report"}

// Prefijo del texto de cada tipo (lo demás es el nombre guardado).
var eventPrefix = map[string]string{"credit_payment": "Pago de ", "card_statement": "Corte de ", "market_visit": "Día de mercado: ",
	"promotion": "Promoción: ", "savings_goal": "Fecha meta: "}

func listCalendarEvents(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"date"}, "date")
		if err != nil {
			return err
		}
		t := today(reg)
		per, err := api.QueryPeriod(r)
		if err != nil {
			return err
		}
		if per == nil {
			mp := dominio.MonthPeriod(t) // por defecto: el mes en curso
			per = &mp
		}
		credits, err := api.QueryIDs(r, "credit_id")
		if err != nil {
			return err
		}
		qs := r.URL.Query()
		fam := qs.Get("family")
		if fam != "" && !slices.Contains(families, fam) {
			return api.BadRequest("El filtro family solo admite household, external o joint.")
		}
		want := map[string]bool{}
		if raw := qs.Get("event_type"); raw != "" {
			for _, et := range strings.Split(raw, ",") {
				if !slices.Contains(eventTypes, et) {
					return api.BadRequest("El tipo de evento «" + et + "» no existe. Usa income, credit_payment, fixed_expense, upcoming_expense, card_statement, planned_grocery, market_visit, promotion, savings_goal, savings_movement, meal_plan o monthly_report.")
				}
				want[et] = true
			}
		}
		ctx, db := r.Context(), reg.Deps().DB
		from, to := per.Start.String(), per.End.String()
		rows, err := dominio.GastosCalendarRows(ctx, db, from, to)
		if err != nil {
			return err
		}
		// Pagos de crédito: total, partes y estado (pagado, pendiente o vencido; RN-04).
		cs, err := dominio.QueryCredits(ctx, db, "", "ORDER BY c.id")
		if err != nil {
			return err
		}
		names := map[int64]string{}
		for _, c := range cs {
			names[c.ID] = c.Name
		}
		persons, err := dominio.CreditPersons(ctx, db)
		if err != nil {
			return err
		}
		sched, err := dominio.QueryScheduled(ctx, db, "p.status <> 'cancelled' AND sp.due_date BETWEEN ? AND ?", "ORDER BY sp.due_date, sp.id", from, to)
		if err != nil {
			return err
		}
		spShares, err := dominio.ScheduledShares(ctx, db, 0)
		if err != nil {
			return err
		}
		for _, sp := range sched {
			hh, ext, pending := payShares(sp, spShares[sp.ID], persons)
			st := "paid"
			if pending && sp.DueDate < t.String() {
				st = "overdue"
			} else if pending {
				st = "pending"
			}
			rows = append(rows, dominio.CalendarRow{Type: "credit_payment", Date: sp.DueDate, Title: names[sp.CreditID], AmountCents: api.Ptr(sp.AmountCents),
				HouseholdShareCents: api.Ptr(hh), ExternalShareCents: api.Ptr(ext), Status: api.Ptr(st), CreditID: api.Ptr(sp.CreditID),
				RefTable: "scheduled_payment", RefID: sp.ID})
		}
		for i := range rows {
			switch rows[i].Type {
			case "income", "fixed_expense":
				rows[i].Status = api.Ptr("captured")
			case "upcoming_expense":
				rows[i].Status = api.Ptr(map[bool]string{true: "overdue", false: "pending"}[rows[i].Date < t.String()])
			}
		}
		// Días de mercado (RN-16): las fechas las da precios.NextVisit.
		markets, err := dominio.GastosMarkets(ctx, db)
		if err != nil {
			return err
		}
		end := per.End.StartOfDay()
		for _, m := range markets {
			st := precios.Store{Code: strconv.FormatInt(m.ID, 10), Name: m.Name, Kind: "market", VisitWeekday: int(m.VisitWeekday)}
			for d := per.Start.StartOfDay(); !d.After(end); d = d.AddDate(0, 0, 1) {
				v, err := precios.NextVisit(st, d)
				if err != nil || v.After(end) {
					break
				}
				rows = append(rows, dominio.CalendarRow{Type: "market_visit", Date: dominio.DateOf(v).String(), Title: m.Name, RefTable: "store", RefID: m.ID})
				d = v
			}
		}
		// Envío del reporte mensual: el día reports.monthly_day_of_month (en un mes corto, su último día; Q2-66 B).
		if nid, err := dominio.GastosNotificationID(ctx, db, "monthly_report"); err != nil {
			return err
		} else if nid > 0 {
			day, err := dominio.SettingInt(ctx, db, "reports.monthly_day_of_month")
			if err != nil {
				return err
			}
			for m := time.Date(per.Start.Year, per.Start.Month, 1, 12, 0, 0, 0, time.UTC); !m.After(end); m = m.AddDate(0, 1, 0) {
				d := dominio.NewDate(m.Year(), m.Month(), min(int(day), dominio.DaysIn(m.Year(), m.Month())))
				if !d.Before(per.Start) && !d.After(per.End) {
					rows = append(rows, dominio.CalendarRow{Type: "monthly_report", Date: d.String(), Title: "Envío del reporte mensual", RefTable: "notification", RefID: nid})
				}
			}
		}
		items := []api.CalendarEvent{}
		for _, c := range rows {
			if len(want) > 0 && !want[c.Type] {
				continue
			}
			if len(credits) > 0 && (c.CreditID == nil || !slices.Contains(credits, *c.CreditID)) {
				continue
			}
			if !familyMatch(fam, c) {
				continue
			}
			ev := api.CalendarEvent{EventType: c.Type, Date: c.Date, Title: eventPrefix[c.Type] + c.Title, AmountCents: c.AmountCents,
				HouseholdShareCents: c.HouseholdShareCents, ExternalShareCents: c.ExternalShareCents, Status: c.Status, CreditID: c.CreditID,
				RefTable: c.RefTable, RefID: c.RefID}
			if c.EndDate != nil {
				ev.EndDate = c.EndDate
			}
			items = append(items, ev)
		}
		slices.SortStableFunc(items, func(a, b api.CalendarEvent) int {
			return cmp.Or(cmp.Compare(a.Date, b.Date), cmp.Compare(slices.Index(eventTypes, a.EventType), slices.Index(eventTypes, b.EventType)),
				cmp.Compare(a.RefID, b.RefID))
		})
		if p.Offset >= len(items) {
			items = nil
		} else {
			items = items[p.Offset:min(len(items), p.Offset+p.Limit+1)]
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// familyMatch aplica el filtro family (incluyente, opción A de C1-filtro-familia) a los eventos con partes;
// los eventos sin partes (ingresos, cortes, promociones, ahorro, mercados, reporte) son del hogar
// (PROVISIONAL, bloqueo S6-gastos-inicio §14).
func familyMatch(fam string, c dominio.CalendarRow) bool {
	hh, ext := true, false
	if c.HouseholdShareCents != nil {
		hh, ext = *c.HouseholdShareCents > 0, *c.ExternalShareCents > 0
	}
	switch fam {
	case "household":
		return hh
	case "external":
		return ext
	case "joint":
		return hh && ext
	}
	return true
}
