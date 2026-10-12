package comparador

import mc "finanzas-hogar/internal/motor_cred"

// span es la vida de una compra: de su fecha a su último pago, ambos incluidos (RN-21).
type span struct{ from, to mc.Date }

// creditCtx guarda lo que no cambia entre escenarios de un crédito.
type creditCtx struct {
	spans   []span
	lastDue mc.Date
	entries []mc.CalendarEntry // calendario del crédito sin escenario en el rango global
}

func prepare(cs mc.CreditState, from, to mc.Date) (creditCtx, error) {
	var out creditCtx
	for _, ps := range cs.Purchases {
		if ps.Status == "cancelled" {
			continue
		}
		in := ps.Input
		in.Credit = cs.Credit
		p, err := mc.QuotePurchase(in)
		if err != nil {
			return out, err
		}
		if n := len(p.Payments); n > 0 {
			last := p.Payments[n-1].Due
			out.spans = append(out.spans, span{in.Date, last})
			if out.lastDue.IsZero() || last.Compare(out.lastDue) > 0 {
				out.lastDue = last
			}
		}
	}
	cal, err := mc.UnifiedCalendar([]mc.CreditState{cs}, from, to)
	if err != nil {
		return out, err
	}
	out.entries = cal.Entries
	return out, nil
}

// withScenario agrega la compra del escenario al estado del crédito (sin tocar el original).
func withScenario(cs mc.CreditState, in mc.PurchaseInput) mc.CreditState {
	out := cs
	out.Purchases = append(append(make([]mc.PurchaseState, 0, len(cs.Purchases)+1), cs.Purchases...),
		mc.PurchaseState{Ref: scenarioRef, Status: "active", Input: in})
	return out
}

// minAvailable: mínimo del disponible por día en [p, l]; el de un día es el que queda después
// del último evento con fecha ≤ ese día (§8.3), en la lectura del crédito (vacío = B).
func minAvailable(tl mc.Timeline, p, l mc.Date) mc.Cents {
	val := func(e mc.TimelineEvent) mc.Cents {
		if tl.Reading == "A" {
			return e.AvailableA
		}
		return e.AvailableB
	}
	var res mc.Cents
	have := false
	evs := tl.Events
	for i, e := range evs {
		if i+1 < len(evs) && evs[i+1].Date.Compare(e.Date) == 0 {
			continue // solo cuenta el valor al cierre del día
		}
		v := val(e)
		switch {
		case e.Date.Compare(p) <= 0:
			res, have = v, true // valor con que abre la ventana
		case e.Date.Compare(l) <= 0:
			if !have || v < res {
				res, have = v, true
			}
		}
	}
	return res
}

// liveMax: máximo de compras vivas simultáneas en [p, l]; el máximo cae en p o en el inicio de
// alguna compra dentro de la ventana.
func liveMax(spans []span, p, l mc.Date) int {
	count := func(d mc.Date) int {
		n := 0
		for _, s := range spans {
			if s.from.Compare(d) <= 0 && d.Compare(s.to) <= 0 {
				n++
			}
		}
		return n
	}
	best := count(p)
	for _, s := range spans {
		if s.from.Compare(p) > 0 && s.from.Compare(l) <= 0 {
			best = max(best, count(s.from))
		}
	}
	return best
}

// violation devuelve la primera restricción violada en el orden de C1 o "".
func violation(s Scenario, u Constraints) string {
	switch {
	case s.MinAvailable < *u.MinAvailableBuffer:
		return ViolBuffer
	case s.MaxLivePurchases > *u.MaxLivePurchases:
		return ViolLive
	case s.ZeroDebtDate.Compare(*u.ZeroDebtBy) > 0:
		return ViolZeroDebt
	}
	return ""
}
