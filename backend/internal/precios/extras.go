package precios

import (
	"fmt"
	"time"
)

// NextVisit (RN-16) es la primera fecha ≥ from cuyo día ISO es el día de visita del mercado
// (si from es día de visita, cuenta: supuesto S-10).
func NextVisit(st Store, from time.Time) (time.Time, error) {
	if st.VisitWeekday < 1 || st.VisitWeekday > 7 {
		return time.Time{}, fmt.Errorf("%w: %s sin día de visita", ErrInvalidArgument, st.Code)
	}
	return from.AddDate(0, 0, (st.VisitWeekday-isoWeekday(from)+7)%7), nil
}

// Purchase es una compra de despensa registrada (para anualizar la membresía).
type Purchase struct {
	Date  time.Time
	Store string
}

// MembershipResult es la evaluación de la membresía (§8, RN-13, S5-09).
type MembershipResult struct {
	Store              string `json:"store"`
	AssumesMembership  bool   `json:"assumes_membership"`
	BasketSavingsCents int64  `json:"basket_savings_cents"`
	AnnualFeeCents     *int64 `json:"annual_fee_cents"`
	BreakEvenPurchases *int64 `json:"break_even_purchases"`
	PurchasesPerYear   int64  `json:"purchases_per_year"`
	AnnualSavingsCents int64  `json:"annual_savings_cents"`
	Verdict            string `json:"verdict"`
	Note               string `json:"note,omitempty"`
}

// Membership evalúa la membresía de una tienda que la exige. El ahorro por compra viene dado
// (savings != nil) o sale de la canasta: Σ max(0, costo en la mejor tienda sin membresía −
// costo en la tienda). Equilibrio = techo(costo anual ÷ ahorro por compra) (FIRME, puerta 2, Q2-63 = A);
// compras al año = despensas en cadenas de los últimos 365 días (FIRME, puerta 2, Q2-62 = B).
func (b *Book) Membership(store string, savings *int64, list []ListItem, hist []Purchase) (*MembershipResult, error) {
	st, ok := b.Store(store)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownStore, store)
	}
	r := &MembershipResult{Store: st.Code, AssumesMembership: st.AssumesMembership(), AnnualFeeCents: st.MembershipFeeCents}
	if savings != nil {
		if *savings < 0 {
			return nil, fmt.Errorf("%w: ahorro negativo", ErrInvalidArgument)
		}
		r.BasketSavingsCents = *savings
	} else if len(list) > 0 {
		cmp, err := b.CompareBasket(list)
		if err != nil {
			return nil, err
		}
		for _, opts := range cmp.opts {
			var mine, other *option
			for k := range opts {
				o := &opts[k]
				switch {
				case o.store.Code == st.Code:
					mine = o
				case !o.store.RequiresMembership && (other == nil || o.cost < other.cost):
					other = o
				}
			}
			if mine != nil && other != nil && other.cost > mine.cost {
				r.BasketSavingsCents += other.cost - mine.cost
			}
		}
	}
	for _, p := range hist {
		ps, ok := b.Store(p.Store)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStore, p.Store)
		}
		age := AgeDays(p.Date, b.Today)
		if ps.Kind != KindMarket && age >= 0 && age < 365 {
			r.PurchasesPerYear++
		}
	}
	var err error
	if r.AnnualSavingsCents, err = mul(r.BasketSavingsCents, r.PurchasesPerYear); err != nil {
		return nil, err
	}
	switch {
	case r.AnnualFeeCents == nil:
		// §8 y criterio S5-09 (alta): con costo ausente la salida dice «sin costo capturado»
		// PROVISIONAL. La marca se queda hasta que el integrador alinee S5-09 con §8 (REVISION-motor-S5-r1, M-3).
		r.Verdict, r.Note = "missing_fee", "sin costo capturado (PROVISIONAL, S5-09)"
	case r.BasketSavingsCents == 0:
		r.Verdict = "not_worth_it"
	default:
		be := (*r.AnnualFeeCents + r.BasketSavingsCents - 1) / r.BasketSavingsCents
		r.BreakEvenPurchases = &be
		r.Verdict = "not_worth_it"
		if r.AnnualSavingsCents >= *r.AnnualFeeCents {
			r.Verdict = "worth_it"
		}
	}
	return r, nil
}

// StaleItem es un precio viejo en la fecha de despensa (insumo de la alerta de S7).
type StaleItem struct {
	Product      string `json:"product"`
	Store        string `json:"store"`
	Presentation string `json:"presentation"`
	ObservedDate string `json:"observed_date"`
	AgeDays      int    `json:"age_days"`
	Weight       Frac   `json:"weight"`
}

// StaleReport lista los precios viejos y los productos sin precio que compita (§10).
type StaleReport struct {
	Stale        []StaleItem `json:"stale"`
	WithoutPrice []string    `json:"without_price"`
}

// StalePricesForPlan (§10, S5-12) cotiza cada producto del plan con la regla de QuoteFor
// evaluada en la fecha de despensa: en la tienda que eligió la persona (assigned) o, si no hay,
// la recomendada. Devuelve las cotizaciones viejas en esa fecha y, aparte, los productos sin
// precio que compita (un precio de más de D0 días ya no compite: queda «sin precio»).
func (b *Book) StalePricesForPlan(products []string, assigned map[string]string, grocery time.Time) (*StaleReport, error) {
	r := &StaleReport{Stale: []StaleItem{}, WithoutPrice: []string{}}
	seen := map[string]bool{}
	for _, pc := range products {
		if seen[pc] { // §15.2: en el orden de products[], sin repetidos
			continue
		}
		seen[pc] = true
		q, ok, err := b.QuoteFor(pc, assigned[pc], "", grocery)
		switch {
		case err != nil:
			return nil, err
		case !ok:
			r.WithoutPrice = append(r.WithoutPrice, pc)
		case q.IsStale:
			r.Stale = append(r.Stale, StaleItem{pc, q.Store, q.Presentation,
				q.ObservedDate.Format(time.DateOnly), q.Weight.AgeDays, q.Weight.Frac})
		}
	}
	return r, nil
}

// ProductStatus es el estado de precio de un producto (C1 PriceSummary).
type ProductStatus struct {
	Product string `json:"product"`
	Status  string `json:"status"` // fresh · stale · no_price
}

// StoreLatest es la fecha más reciente de precio que compite por tienda.
type StoreLatest struct {
	Store      string  `json:"store"`
	LatestDate *string `json:"latest_date"`
}

// Summary arma el resumen de frescura (C1 PriceSummary) en «hoy»: fresh si algún precio que
// compite no es viejo, stale si todos los que compiten son viejos, no_price si ninguno compite;
// por tienda, la fecha más reciente de un precio que compite.
func (b *Book) Summary() ([]ProductStatus, []StoreLatest) {
	var ps []ProductStatus
	for _, prod := range b.Catalog.Products {
		status := "no_price"
		for _, o := range b.candidates(prod, "", "", b.Today) {
			if !o.weight.IsStale {
				status = "fresh"
				break
			}
			status = "stale"
		}
		ps = append(ps, ProductStatus{prod.Code, status})
	}
	var sl []StoreLatest
	for _, st := range b.Stores {
		l := StoreLatest{Store: st.Code}
		var last *PriceObs
		for i := range b.Prices {
			p := &b.Prices[i]
			if c, _ := Competes(p.Source); c && p.Store == st.Code && (last == nil || p.ObservedDate.After(last.ObservedDate)) {
				last = p
			}
		}
		if last != nil {
			d := last.ObservedDate.Format(time.DateOnly)
			l.LatestDate = &d
		}
		sl = append(sl, l)
	}
	return ps, sl
}
