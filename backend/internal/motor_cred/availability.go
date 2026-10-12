package motor_cred

import "sort"

// TimelineInput: lo arma quien llama (consulta linea_tiempo, §9.2).
type TimelineInput struct {
	State           CreditState
	Cash            []CashEvent
	Today, From, To Date
}

// TimelineEvent es un evento de la línea de tiempo del disponible.
type TimelineEvent struct {
	Date                Date   `json:"date"`
	Kind                string `json:"kind"`
	Description         string `json:"description"`
	PurchaseRef         string `json:"purchase_ref"`
	Seq                 int    `json:"seq"`
	Amount              Cents  `json:"amount_cents"`
	PrincipalRestored   Cents  `json:"principal_restored_cents"`
	InterestNotRestored Cents  `json:"interest_not_restored_cents"`
	AvailableA          Cents  `json:"available_a_cents"`
	AvailableB          Cents  `json:"available_b_cents"`
}

// Timeline es el disponible A y B de una línea (§5.2).
type Timeline struct {
	Events              []TimelineEvent `json:"events"`
	SeriesA             []Cents         `json:"series_a_cents"`
	SeriesB             []Cents         `json:"series_b_cents"`
	MinA                Cents           `json:"min_a_cents"`
	MinB                Cents           `json:"min_b_cents"`
	Reading             Reading         `json:"reading"`
	PendingConfirmation bool            `json:"pending_confirmation"`
	AAnchor             Date            `json:"a_anchor"`
	AAnchorDefaulted    bool            `json:"a_anchor_defaulted"`
}

var kindRank = map[string]int{"start": 0, "purchase": 1, "cash": 2, "payment": 3}

type rankedEvent struct {
	ev  TimelineEvent
	sub int
}

type restoration struct {
	date Date
	k    Cents
}

// paymentEvents aplica K15: un evento por abono; el que completa el pago repone capital e
// interés no repuesto; el resto de un pago no vencido se proyecta en su vencimiento.
func paymentEvents(st CreditState, q quotedPurchase, today Date) ([]rankedEvent, []restoration, error) {
	var evs []rankedEvent
	var rest []restoration
	n := len(q.Purchase.Payments)
	for _, r := range st.Payments {
		if r.PurchaseRef == q.Ref && (r.Seq < 1 || r.Seq > n || r.Amount < 0) {
			return nil, nil, invalid("abono %s/%d", r.PurchaseRef, r.Seq)
		}
	}
	for _, p := range q.Purchase.Payments {
		recs := st.recordsFor(q.Ref, p.Seq)
		var cum Cents
		done := false
		for j, r := range recs {
			cum += r.Amount
			ev := TimelineEvent{Date: r.Date, Kind: "payment", Description: "Abono", PurchaseRef: q.Ref, Seq: p.Seq, Amount: r.Amount}
			if !done && cum >= p.Amount {
				ev.PrincipalRestored, ev.InterestNotRestored = p.Principal, p.Interest
				rest = append(rest, restoration{r.Date, p.Principal})
				done = true
			}
			evs = append(evs, rankedEvent{ev, j})
		}
		if cum < p.Amount && p.Due.Compare(today) >= 0 {
			evs = append(evs, rankedEvent{TimelineEvent{Date: p.Due, Kind: "payment", Description: "Pago programado",
				PurchaseRef: q.Ref, Seq: p.Seq, Amount: p.Amount - cum, PrincipalRestored: p.Principal,
				InterestNotRestored: p.Interest}, len(recs)})
			rest = append(rest, restoration{p.Due, p.Principal})
		}
	}
	return evs, rest, nil
}

func between(d, lo, hi Date) bool { return d.Compare(lo) > 0 && d.Compare(hi) <= 0 }

// AvailabilityTimeline calcula las dos lecturas del disponible (RN-09, K15, S3-03, S3-05).
func AvailabilityTimeline(in TimelineInput) (Timeline, error) {
	c := in.State.Credit
	if c.Limit == nil || in.To.Compare(in.From) < 0 || in.From.IsZero() {
		return Timeline{}, invalid("línea sin límite o ventana inválida")
	}
	reading := c.Reading
	if reading == "" {
		reading = "B"
	}
	if reading != "A" && reading != "B" {
		return Timeline{}, invalid("lectura %q", c.Reading)
	}
	anchor, defaulted := in.From, true
	if c.ReadingAAnchor != nil {
		anchor, defaulted = *c.ReadingAAnchor, false
	} else if reading == "A" {
		return Timeline{}, invalid("lectura A sin fecha ancla")
	}
	if in.From.Compare(anchor) < 0 {
		return Timeline{}, invalid("from antes del ancla")
	}
	qs, err := in.State.quoted()
	if err != nil {
		return Timeline{}, err
	}
	limit := *c.Limit
	a, b := limit, limit
	var evs []rankedEvent
	for _, q := range qs {
		if !q.live() {
			continue
		}
		d, k := q.Input.Date, q.Purchase.Principal
		if d.Compare(in.From) <= 0 {
			b -= k
		}
		if between(d, anchor, in.From) {
			a -= k
		}
		if between(d, in.From, in.To) {
			evs = append(evs, rankedEvent{TimelineEvent{Date: d, Kind: "purchase", Description: "Compra",
				PurchaseRef: q.Ref, Amount: k}, 0})
		}
		pe, rest, err := paymentEvents(in.State, q, in.Today)
		if err != nil {
			return Timeline{}, err
		}
		for _, r := range rest {
			if r.date.Compare(in.From) <= 0 {
				b += r.k
			}
			if between(r.date, anchor, in.From) {
				a += r.k
			}
		}
		for _, e := range pe {
			if between(e.ev.Date, in.From, in.To) {
				evs = append(evs, e)
			}
		}
	}
	for _, ce := range in.Cash {
		if between(ce.Date, in.From, in.To) {
			evs = append(evs, rankedEvent{TimelineEvent{Date: ce.Date, Kind: "cash", Description: ce.Description, Amount: ce.Amount}, 0})
		}
	}
	sort.SliceStable(evs, func(i, j int) bool {
		x, y := evs[i], evs[j]
		if c := x.ev.Date.Compare(y.ev.Date); c != 0 {
			return c < 0
		}
		if kindRank[x.ev.Kind] != kindRank[y.ev.Kind] {
			return kindRank[x.ev.Kind] < kindRank[y.ev.Kind]
		}
		if x.ev.PurchaseRef != y.ev.PurchaseRef {
			return x.ev.PurchaseRef < y.ev.PurchaseRef
		}
		if x.ev.Seq != y.ev.Seq {
			return x.ev.Seq < y.ev.Seq
		}
		return x.sub < y.sub
	})
	tl := Timeline{Reading: reading, PendingConfirmation: !c.ReadingConfirmed, AAnchor: anchor,
		AAnchorDefaulted: defaulted, SeriesA: []Cents{}, SeriesB: []Cents{}, MinA: a, MinB: b}
	tl.Events = append(make([]TimelineEvent, 0, len(evs)+1),
		TimelineEvent{Date: in.From, Kind: "start", Description: "Inicio", AvailableA: a, AvailableB: b})
	for _, e := range evs {
		ev := e.ev
		var delta Cents
		switch ev.Kind {
		case "purchase":
			delta = -ev.Amount
		case "payment":
			delta = ev.PrincipalRestored
		}
		a += delta
		b = min(limit, b+delta)
		ev.AvailableA, ev.AvailableB = a, b
		tl.Events = append(tl.Events, ev)
		if ev.Kind != "cash" {
			tl.SeriesA = append(tl.SeriesA, a)
			tl.SeriesB = append(tl.SeriesB, b)
		}
		tl.MinA, tl.MinB = min(tl.MinA, a), min(tl.MinB, b)
	}
	return tl, nil
}
