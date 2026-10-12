// Package creditos implementa las operaciones de Créditos de C1 (líneas, tasas por plazo, compras a
// crédito, pagos, promociones, valores del comparador y comparador de escenarios). Los cálculos son
// del motor S3 (motor_cred y motor_cred/comparador); aquí solo se leen y escriben renglones.
package creditos

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// Register reclama las 29 operaciones de api-creditos (CONVENCIONES-backend §14).
func Register(reg *api.Registry) {
	for id, h := range map[string]func(*api.Registry) api.HandlerFunc{
		"listCredits": listCredits, "getCreditTotals": getCreditTotals, "getDebtTrend": getDebtTrend,
		"getCredit": getCredit, "previewAvailability": previewAvailability, "getAvailabilitySeries": getAvailabilitySeries,
		"getCardStatus": getCardStatus, "listCreditTermRates": listCreditTermRates, "createCredit": createCredit,
		"updateCredit": updateCredit, "createCreditTermRate": createCreditTermRate, "updateCreditTermRate": updateCreditTermRate,
		"deleteCreditTermRate": deleteCreditTermRate, "listPromotions": listPromotions, "createPromotion": createPromotion,
		"updatePromotion": updatePromotion, "deletePromotion": deletePromotion, "getOptimizerSettings": getOptimizerSettings,
		"updateOptimizerSettings": updateOptimizerSettings, "compareScenarios": compareScenarios,
		"listCreditPurchases": listCreditPurchases, "createCreditPurchase": createCreditPurchase,
		"previewCreditPurchase": previewCreditPurchase, "getCreditPurchase": getCreditPurchase,
		"updateCreditPurchase": updateCreditPurchase, "deleteCreditPurchase": deleteCreditPurchase,
		"listScheduledPayments": listScheduledPayments, "markScheduledPaymentPaid": markScheduledPaymentPaid,
		"deleteCreditPayment": deleteCreditPayment,
	} {
		reg.Handle(id, h(reg))
	}
}

func today(reg *api.Registry) dominio.Date { return dominio.Today(reg.Deps().Clock) }

func snap(reg *api.Registry, r *http.Request) (*snapshot, error) {
	return load(r.Context(), reg.Deps().DB, today(reg))
}

func notFound(err error, msg string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.NotFound(msg)
	}
	return err
}

const msgCredit = "No encontramos ese crédito. Revisa la lista de créditos."

// nullable aplica un campo de JSON Merge Patch (ausente: no cambia; null: borra; valor: cambia).
func nullable[T any](dst **T, n api.Nullable[T]) {
	if n.Set {
		if n.Valid {
			v := n.Value
			*dst = &v
		} else {
			*dst = nil
		}
	}
}

func set[T any](dst *T, p *T) {
	if p != nil {
		*dst = *p
	}
}

func deref[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

func oneOf(v string, opts ...string) bool {
	for _, o := range opts {
		if v == o {
			return true
		}
	}
	return false
}

func badEnum(r *http.Request, name string, opts ...string) (string, error) {
	v := r.URL.Query().Get(name)
	if v != "" && !oneOf(v, opts...) {
		return "", api.BadRequest("El filtro " + name + " no es válido; usa " + strings.Join(opts, ", ") + ".")
	}
	return v, nil
}

// view arma la respuesta Credit con sus campos calculados por el motor (deuda viva y disponible).
// Si el motor no puede calcular (p. ej. falta capturar), el campo se omite: nunca se inventa un número.
func (s *snapshot) view(c dominio.CreditRow) api.Credit {
	out := api.Credit{ID: c.ID, Name: c.Name, Kind: c.Kind, Lender: c.Lender, HolderPersonID: c.HolderPersonID,
		CreditLimitCents: c.CreditLimitCents, PaymentFrequency: c.PaymentFrequency, FirstPaymentMinDays: c.FirstPaymentMinDays,
		AvailabilityReading: c.AvailabilityReading, AvailabilityConfirmed: c.AvailabilityConfirmed,
		ReadingAAnchorDate: c.ReadingAAnchorDate, StatementDay: c.StatementDay, PaymentDueDay: c.PaymentDueDay,
		AnnualRateNum: c.AnnualRateNum, AnnualRateDen: c.AnnualRateDen, VatOnInterestBp: c.VatOnInterestBp,
		AnnualFeeCents: c.AnnualFeeCents, LateFeeCents: c.LateFeeCents, OpeningFeeCents: c.OpeningFeeCents,
		IsActive: c.IsActive, Notes: c.Notes, CreatedAt: api.APITime(c.CreatedAt), UpdatedAt: api.APITime(c.UpdatedAt),
		RateBase: c.RateBase, MinPaymentPctBp: c.MinPaymentPctBp, MinPaymentFixedCents: c.MinPaymentFixedCents,
		AnnualRateDecimals: c.AnnualRateDecimals, VatIncluded: c.VatIncluded}
	st := s.state(c)
	debt, derr := mc.LiveDebt([]mc.CreditState{st}, mdateOf(s.today))
	if derr == nil {
		out.LiveDebtCents = api.Ptr(int64(debt.Total))
	}
	switch c.Kind {
	case "revolving_line":
		if tl, err := s.timeline(c); err == nil && len(tl.Events) > 0 {
			v := tl.Events[0].AvailableB
			if tl.Reading == "A" {
				v = tl.Events[0].AvailableA
			}
			out.AvailableCents = api.Ptr(int64(v))
		}
	case "bank_card":
		if c.CreditLimitCents != nil && derr == nil { // límite − saldo (C1); saldo = deuda viva del motor (PROVISIONAL)
			out.AvailableCents = api.Ptr(*c.CreditLimitCents - int64(debt.Total))
		}
	}
	active := map[int64]bool{}
	var live int64
	for _, p := range s.purchases {
		if p.CreditID == c.ID && p.Status == "active" {
			active[p.ID] = true
			live++
		}
	}
	out.LivePurchaseCount = api.Ptr(live)
	var next string
	var nextCents int64
	for _, sp := range s.sched {
		rest := sp.AmountCents - s.paidBySP[sp.ID]
		if sp.CreditID != c.ID || !active[sp.CreditPurchaseID] || rest <= 0 {
			continue
		}
		if next == "" || sp.DueDate < next {
			next, nextCents = sp.DueDate, 0
		}
		if sp.DueDate == next {
			nextCents += rest
		}
	}
	if next != "" {
		out.NextPaymentDate, out.NextPaymentCents = api.Ptr(next), api.Ptr(nextCents)
	}
	if c.AnnualRateNum != nil && c.AnnualRateDen != nil {
		out.AnnualRatePct = ratePct(*c.AnnualRateNum, *c.AnnualRateDen, c.AnnualRateDecimals)
	}
	if c.Kind == "bank_card" || c.Kind == "store_credit" { // K18, RN-26: los 5 datos del estado de cuenta
		m := []string{}
		if c.AnnualRateNum == nil {
			m = append(m, "annual_rate")
		}
		if c.RateBase == nil {
			m = append(m, "rate_base")
		}
		if c.VatIncluded == nil {
			m = append(m, "vat")
		}
		if c.MinPaymentPctBp == nil {
			m = append(m, "min_payment_pct_bp")
		}
		if c.MinPaymentFixedCents == nil {
			m = append(m, "min_payment_fixed_cents")
		}
		out.MissingTerms = m
	}
	return out
}

const nextPaymentSQL = `(SELECT MIN(sp.due_date) FROM scheduled_payment sp JOIN credit_purchase p2 ON p2.id = sp.credit_purchase_id
 WHERE p2.credit_id = c.id AND p2.status = 'active' AND sp.amount_cents > (SELECT COALESCE(SUM(m.amount_cents), 0) FROM credit_payment m WHERE m.scheduled_payment_id = sp.id))`

func listCredits(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"next_payment_date", "name"}, "next_payment_date")
		if err != nil {
			return err
		}
		kind, err := badEnum(r, "kind", "revolving_line", "bank_card", "term_loan", "store_credit")
		if err != nil {
			return err
		}
		status, err := badEnum(r, "status", "active", "archived")
		if err != nil {
			return err
		}
		where, args := []string{"c.is_active = ?"}, []any{status != "archived"}
		if kind != "" {
			where, args = append(where, "c.kind = ?"), append(args, kind)
		}
		if t := strings.TrimSpace(r.URL.Query().Get("q")); t != "" {
			where, args = append(where, "(c.name LIKE ? OR COALESCE(c.lender, '') LIKE ?)"), append(args, "%"+t+"%", "%"+t+"%")
		}
		rows, err := dominio.QueryCredits(r.Context(), reg.Deps().DB, strings.Join(where, " AND "),
			p.OrderBy(map[string]string{"next_payment_date": nextPaymentSQL, "name": "c.name", "id": "c.id"})+" "+p.LimitOffset(), args...)
		if err != nil {
			return err
		}
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		items := make([]api.Credit, 0, len(rows))
		for _, c := range rows {
			items = append(items, s.view(c))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func getCredit(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		s, c, err := creditFor(reg, r)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, s.view(c))
	}
}

func creditFor(reg *api.Registry, r *http.Request) (*snapshot, dominio.CreditRow, error) {
	id, err := api.PathID(r, "credit_id")
	if err != nil {
		return nil, dominio.CreditRow{}, err
	}
	s, err := snap(reg, r)
	if err != nil {
		return nil, dominio.CreditRow{}, err
	}
	c, ok := s.byID[id]
	if !ok {
		return nil, c, dominio.NotFound(msgCredit)
	}
	return s, c, nil
}

func monthRange(d dominio.Date) (mc.Date, mc.Date) {
	return mc.Date{Year: d.Year, Month: d.Month, Day: 1}, mc.Date{Year: d.Year, Month: d.Month, Day: dominio.DaysIn(d.Year, d.Month)}
}

func getCreditTotals(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		states := s.states(nil)
		debt, err := s.liveDebt(states)
		if err != nil {
			return motorErr(err, "credit_id")
		}
		out := api.CreditTotals{LiveDebtCents: int64(debt.Total)}
		from, to := monthRange(s.today)
		fFrom, fTo := from, mc.Date{Year: from.Year, Month: from.Month, Day: 15} // quincena del 1 al 15 y del 16 al último (K8)
		if s.today.Day > 15 {
			fFrom, fTo = mc.Date{Year: from.Year, Month: from.Month, Day: 16}, to
		}
		cal, err := mc.UnifiedCalendar(states, from, to)
		if err != nil {
			return motorErr(err, "credit_id")
		}
		for _, e := range cal.Entries {
			inF := e.Date.Compare(fFrom) >= 0 && e.Date.Compare(fTo) <= 0
			for _, sh := range e.Shares {
				ext := s.isExternal(sh.PersonID)
				out.MonthTotalCents += int64(sh.Amount)
				if inF {
					out.FortnightTotalCents += int64(sh.Amount)
				}
				switch {
				case ext:
					out.MonthExternalCents += int64(sh.Amount)
					if inF {
						out.FortnightExternalCents += int64(sh.Amount)
					}
				default:
					out.MonthHouseholdCents += int64(sh.Amount)
					if inF {
						out.FortnightHouseholdCents += int64(sh.Amount)
					}
				}
			}
			if e.Status == "paid" {
				out.InterestPaidMonthCents += int64(e.Interest)
			}
		}
		ycal, err := mc.UnifiedCalendar(states, mc.Date{Year: s.today.Year, Month: time.January, Day: 1}, mc.Date{Year: s.today.Year, Month: time.December, Day: 31})
		if err != nil {
			return motorErr(err, "credit_id")
		}
		for _, e := range ycal.Entries {
			if e.Status == "paid" {
				out.InterestPaidYearCents += int64(e.Interest)
			}
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// getDebtTrend: 12 meses desde el mes en curso (proyección; PROVISIONAL, bloqueo S6-creditos).
func getDebtTrend(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ids, err := api.QueryIDs(r, "credit_id")
		if err != nil {
			return err
		}
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		states := s.states(ids)
		out := api.DebtTrend{Months: []api.DebtTrendMonthsItem{}}
		d := dominio.NewDate(s.today.Year, s.today.Month, 1)
		for i := 0; i < 12; i++ {
			from, to := monthRange(d)
			debt, err := mc.LiveDebt(states, to)
			if err != nil {
				return motorErr(err, "credit_id")
			}
			cal, err := mc.UnifiedCalendar(states, from, to)
			if err != nil {
				return motorErr(err, "credit_id")
			}
			it := api.DebtTrendMonthsItem{Month: d.MonthKey(), LiveDebtCents: int64(debt.Total), Credits: []api.DebtTrendMonthsItemCreditsItem{}}
			for _, l := range debt.ByCredit {
				it.Credits = append(it.Credits, api.DebtTrendMonthsItemCreditsItem{CreditID: l.CreditID, LiveDebtCents: int64(l.Pending)})
			}
			for _, e := range cal.Entries {
				if e.Status == "paid" {
					it.InterestPaidCents += int64(e.Interest)
				} else {
					it.InterestPendingCents += int64(e.Interest)
				}
			}
			out.Months = append(out.Months, it)
			d = d.AddDays(dominio.DaysIn(d.Year, d.Month))
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func previewAvailability(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		s, c, err := creditFor(reg, r)
		if err != nil {
			return err
		}
		if c.Kind != "revolving_line" {
			return &api.Problem{Status: http.StatusUnprocessableEntity, Code: "not_revolving",
				Message: "La vista previa A/B solo aplica a líneas revolventes."}
		}
		tl, err := s.timeline(c)
		if err != nil {
			return motorErr(err, "credit_id")
		}
		ev := tl.Events[0]
		return api.WriteJSON(w, http.StatusOK, api.AvailabilityPreview{CreditID: c.ID, AvailabilityReading: string(tl.Reading),
			AvailabilityConfirmed: c.AvailabilityConfirmed,
			ReadingA:              api.AvailabilityPreviewReadingA{AvailableTodayCents: int64(ev.AvailableA), HorizonMinCents: int64(tl.MinA)},
			ReadingB:              api.AvailabilityPreviewReadingB{AvailableTodayCents: int64(ev.AvailableB), HorizonMinCents: int64(tl.MinB)}})
	}
}

func getAvailabilitySeries(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		s, c, err := creditFor(reg, r)
		if err != nil {
			return err
		}
		if c.CreditLimitCents == nil {
			return &api.Problem{Status: http.StatusUnprocessableEntity, Code: "not_revolving",
				Message: "La serie de disponible solo aplica a créditos con límite."}
		}
		tl, err := s.timeline(c)
		if err != nil {
			return motorErr(err, "credit_id")
		}
		out := api.AvailabilitySeries{CreditID: c.ID, MinACents: int64(tl.MinA), MinBCents: int64(tl.MinB), Events: []api.AvailabilitySeriesEventsItem{}}
		for _, e := range tl.Events {
			if e.Kind == "cash" {
				continue
			}
			out.Events = append(out.Events, api.AvailabilitySeriesEventsItem{EventDate: e.Date.String(), EventType: e.Kind,
				Description: e.Description, AmountCents: int64(e.Amount), PrincipalRestoredCents: int64(e.PrincipalRestored),
				InterestNotRestoredCents: int64(e.InterestNotRestored), AvailableACents: int64(e.AvailableA), AvailableBCents: int64(e.AvailableB)})
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func getCardStatus(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		s, c, err := creditFor(reg, r)
		if err != nil {
			return err
		}
		if c.Kind != "bank_card" {
			return &api.Problem{Status: http.StatusUnprocessableEntity, Code: "not_a_card",
				Message: "El estado de la tarjeta solo aplica a tarjetas bancarias."}
		}
		st := s.state(c)
		out := api.CardStatus{CreditID: c.ID, InstallmentPlans: []api.CardStatusInstallmentPlansItem{}}
		if debt, err := mc.LiveDebt([]mc.CreditState{st}, mdateOf(s.today)); err == nil {
			out.CurrentBalanceCents = int64(debt.Total) // saldo actual = deuda viva del motor (PROVISIONAL)
		}
		for _, t := range s.statements {
			if t.CreditID != c.ID || t.StatementDate > s.today.String() {
				continue
			}
			out.StatementDate, out.DueDate, out.StatementBalanceCents = api.Ptr(t.StatementDate), api.Ptr(t.DueDate), api.Ptr(t.BalanceCents)
			out.NoInterestPaymentCents, out.MinimumPaymentCents = t.NoInterestCents, t.MinimumCents
			if t.MinimumCents == nil {
				in := mc.StatementInput{StatementDate: mdate(t.StatementDate), DueDate: mdate(t.DueDate), Balance: mc.Cents(t.BalanceCents), Source: t.Source}
				if m, _, err := st.Credit.CardMinimum(in); err == nil {
					out.MinimumPaymentCents = api.Ptr(int64(m))
				}
			}
		}
		for _, p := range s.purchases {
			if p.CreditID != c.ID || p.Status != "active" || p.TermCount < 2 {
				continue
			}
			it := api.CardStatusInstallmentPlansItem{CreditPurchaseID: p.ID, Description: p.Description}
			for _, sp := range s.sched {
				if sp.CreditPurchaseID == p.ID && sp.AmountCents > s.paidBySP[sp.ID] {
					if it.RemainingCount == 0 {
						it.InstallmentCents = sp.AmountCents
					}
					it.RemainingCount++
				}
			}
			out.InstallmentPlans = append(out.InstallmentPlans, it)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}
