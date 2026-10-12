package motor_fin

import (
	"errors"
	"sort"
)

// halfUp redondea num/den a entero con la mitad alejándose del cero (ROUND_HALF_UP, §1.3).
func halfUp(num, den int64) int64 {
	if den < 0 {
		num, den = -num, -den
	}
	neg := num < 0
	if neg {
		num = -num
	}
	q := (2*num + den) / (2 * den)
	if neg {
		return -q
	}
	return q
}

// series da la cuota k (0…n−1) de C en n partes: HALF_UP(C/n) y la última ajusta (RN-02).
func series(c Cents, n, k int) Cents {
	base := Cents(halfUp(int64(c), int64(n)))
	if k < n-1 {
		return base
	}
	return c - Cents(n-1)*base
}

// splitRule reparte un monto con una regla efectiva (RN-03): porcentajes HALF_UP,
// fijos tal cual y la persona del resto absorbe la diferencia.
func splitRule(amount Cents, r *ShareRule) ([]Share, error) {
	var out []Share
	rest := amount
	for _, p := range r.Parts {
		var v Cents
		if p.PctBP != nil {
			v = Cents(halfUp(int64(amount)*(*p.PctBP), 10000))
		} else {
			v = *p.FixedCents
			if v > amount {
				return nil, errors.New("parte fija mayor que el monto")
			}
		}
		out = append(out, Share{p.Person, v})
		rest -= v
	}
	if rest < 0 {
		return nil, errors.New("las partes rebasan el monto")
	}
	return append(out, Share{r.RestPerson, rest}), nil
}

// costItem es un costo que se devenga (gasto, compra o ocurrencia planificada).
type costItem struct {
	date      Date
	months    int
	total     Cents
	shares    map[string]Cents // total por persona
	order     []string
	category  string
	paperCash bool
	seq       int // purchase_seq de despensa, 0 si no aplica
	purchase  int64
}

// cashEvent es una salida de efectivo (pago programado o gasto no crédito).
type cashEvent struct {
	date      Date
	total     Cents
	shares    []Share
	category  string
	payment   bool
	purchase  int64
	seq       int
	pseq      int
	paperCash bool
}

type ledger struct {
	in     *Input
	costs  []costItem
	cash   []cashEvent
	seqOf  map[int64]int
	purchs []*CreditPurchase
	catOf  map[int64]string
}

func (in *Input) newCost(date Date, months int, total Cents, shares []Share, cat string) costItem {
	c := costItem{date: date, months: months, total: total, shares: map[string]Cents{}, category: cat}
	for _, s := range shares {
		if _, ok := c.shares[s.Person]; !ok {
			c.order = append(c.order, s.Person)
		}
		c.shares[s.Person] += s.AmountCents
	}
	return c
}

// slot da el devengo de la cuota k del costo, en total y por persona (§1.3).
func (c *costItem) slot(k int, holder string) (Cents, map[string]Cents) {
	tot := series(c.total, c.months, k)
	pp := map[string]Cents{}
	var others Cents
	for _, p := range c.order {
		if p == holder {
			continue
		}
		v := series(c.shares[p], c.months, k)
		pp[p] += v
		others += v
	}
	pp[holder] += tot - others
	return tot, pp
}

// buildLedger arma costos y efectivo con los registros dados más los extras (ocurrencias).
func buildLedger(in *Input, expenses []Expense, purchases []CreditPurchase, extraCost []costItem, extraCash []cashEvent) *ledger {
	l := &ledger{in: in, seqOf: map[int64]int{}, catOf: map[int64]string{}}
	byID := map[int64]*CreditPurchase{}
	for i := range purchases {
		p := &purchases[i]
		if p.Status == "cancelled" {
			continue
		}
		byID[p.ID] = p
		l.purchs = append(l.purchs, p)
	}
	covered := map[int64]*Expense{}
	for i := range expenses {
		e := &expenses[i]
		if e.PaymentMethod == "credit" && e.CreditPurchaseID != nil && byID[*e.CreditPurchaseID] != nil {
			covered[*e.CreditPurchaseID] = e
		}
	}
	for _, p := range l.purchs {
		if e := covered[p.ID]; e != nil {
			l.catOf[p.ID] = e.Category
		} else if p.Category != nil {
			l.catOf[p.ID] = *p.Category
		}
	}
	var groc []*CreditPurchase
	for _, p := range l.purchs {
		if l.catOf[p.ID] == in.Params.GroceriesCategory {
			groc = append(groc, p)
		}
	}
	sort.SliceStable(groc, func(a, b int) bool {
		if groc[a].PurchaseDate != groc[b].PurchaseDate {
			return groc[a].PurchaseDate.Before(groc[b].PurchaseDate)
		}
		return groc[a].ID < groc[b].ID
	})
	for i, p := range groc {
		l.seqOf[p.ID] = i + 1
	}
	pshares := func(p *CreditPurchase) []Share {
		var s []Share
		for _, x := range p.Shares {
			s = append(s, Share{x.Person, x.TotalCents})
		}
		return s
	}
	for i := range expenses {
		e := &expenses[i]
		if e.PaymentMethod == "credit" {
			p := byID[*e.CreditPurchaseID]
			if p == nil {
				continue // compra cancelada: se ignora
			}
			c := in.newCost(e.ExpenseDate, count1(e.AccrualMonths), p.TotalCents, pshares(p), e.Category)
			c.seq, c.purchase = l.seqOf[p.ID], p.ID
			l.costs = append(l.costs, c)
			continue
		}
		paper := e.Category == in.Params.PaperCategory
		c := in.newCost(e.ExpenseDate, count1(e.AccrualMonths), e.AmountCents, e.Shares, e.Category)
		c.paperCash = paper
		l.costs = append(l.costs, c)
		payer := in.holder
		if e.PaidBy != nil {
			payer = *e.PaidBy
		}
		l.cash = append(l.cash, cashEvent{date: e.ExpenseDate, total: e.AmountCents, shares: []Share{{payer, e.AmountCents}}, category: e.Category, paperCash: paper})
	}
	for _, p := range l.purchs {
		if covered[p.ID] == nil {
			c := in.newCost(p.PurchaseDate, count1(p.AccrualMonths), p.TotalCents, pshares(p), l.catOf[p.ID])
			c.seq, c.purchase = l.seqOf[p.ID], p.ID
			l.costs = append(l.costs, c)
		}
		for _, pay := range p.Payments {
			l.cash = append(l.cash, cashEvent{date: pay.DueDate, total: pay.AmountCents, shares: pay.Shares,
				category: l.catOf[p.ID], payment: true, purchase: p.ID, seq: pay.Seq, pseq: l.seqOf[p.ID]})
		}
	}
	l.costs = append(l.costs, extraCost...)
	l.cash = append(l.cash, extraCash...)
	return l
}

// agg acumula el devengo de un periodo.
type agg struct {
	total    Cents
	byPerson map[string]Cents
	byCat    map[string]Cents
	catPers  map[string]map[string]Cents
	paper    Cents
	maxSeq   int
}

func newAgg() *agg {
	return &agg{byPerson: map[string]Cents{}, byCat: map[string]Cents{}, catPers: map[string]map[string]Cents{}}
}

func (a *agg) add(c *costItem, tot Cents, pp map[string]Cents) {
	a.total += tot
	a.byCat[c.category] += tot
	if a.catPers[c.category] == nil {
		a.catPers[c.category] = map[string]Cents{}
	}
	for p, v := range pp {
		a.byPerson[p] += v
		a.catPers[c.category][p] += v
	}
	if c.paperCash {
		a.paper += tot
	}
	if c.seq > a.maxSeq {
		a.maxSeq = c.seq
	}
}

// monthly reparte el devengo de todos los costos por índice de mes.
func (l *ledger) monthly() map[int]*agg {
	out := map[int]*agg{}
	for i := range l.costs {
		c := &l.costs[i]
		m0 := c.date.MonthOf().idx()
		for k := 0; k < c.months; k++ {
			tot, pp := c.slot(k, l.in.holder)
			a := out[m0+k]
			if a == nil {
				a = newAgg()
				out[m0+k] = a
			}
			a.add(c, tot, pp)
		}
	}
	return out
}

// windowed reparte el devengo en ventanas: n ventanas consecutivas desde la que contiene la fecha.
// periods ya comprobó que las ventanas son ascendentes y disjuntas (PROVISIONAL, S2-ventanas).
func (l *ledger) windowed(ws []Window) []*agg {
	out := make([]*agg, len(ws))
	for i := range out {
		out[i] = newAgg()
	}
	for i := range l.costs {
		c := &l.costs[i]
		w0 := -1
		for j, w := range ws {
			if c.date.in(w.Start, w.End) {
				w0 = j
				break
			}
		}
		if w0 < 0 {
			continue // posterior a la última ventana: checkCovered ya rechazó lo anterior y los huecos
		}
		for k := 0; k < c.months && w0+k < len(ws); k++ {
			tot, pp := c.slot(k, l.in.holder)
			out[w0+k].add(c, tot, pp)
		}
	}
	return out
}

// uncovered dice si d cae antes de la primera ventana o en un hueco entre ventanas: no está en
// ninguna y no es posterior a la última. Lo posterior a la última ventana no entra en ningún
// periodo ni en ninguna diferencia acumulada, así que se admite.
func uncovered(ws []Window, d Date) bool {
	hi := ws[0].End
	for _, w := range ws {
		if d.in(w.Start, w.End) {
			return false
		}
		if w.End.After(hi) {
			hi = w.End
		}
	}
	return !d.After(hi)
}

// checkCovered rechaza con invalid_argument un costo o un pago anterior a la primera ventana o
// en un hueco entre ventanas. PROVISIONAL (docs/estado/bloqueos/S2-ventanas.md): §3.2 pide que
// cumulative_diff_cents incluya lo anterior no pedido, pero §2.12 y §3.2 no dicen cómo se devenga
// un costo que no cae en ninguna ventana; mientras Axel decide, la regla es la del oráculo.
func (l *ledger) checkCovered(ws []Window) error {
	for i := range l.costs {
		if d := l.costs[i].date; uncovered(ws, d) {
			return errf(CodeInvalidArgument, "windows: costo del %s antes de la primera ventana o en un hueco entre ventanas (PROVISIONAL, S2-ventanas)", d)
		}
	}
	for i := range l.cash {
		if d := l.cash[i].date; uncovered(ws, d) {
			return errf(CodeInvalidArgument, "windows: pago del %s antes de la primera ventana o en un hueco entre ventanas (PROVISIONAL, S2-ventanas)", d)
		}
	}
	return nil
}

func accruedUpTo(ms map[int]*agg, last int) Cents {
	var s Cents
	for i, a := range ms {
		if i <= last {
			s += a.total
		}
	}
	return s
}

func (l *ledger) cashUpTo(d Date) Cents {
	var s Cents
	for _, e := range l.cash {
		if !e.date.After(d) {
			s += e.total
		}
	}
	return s
}

func (l *ledger) cashBetween(a, b Date) Cents {
	var s Cents
	for _, e := range l.cash {
		if e.date.in(a, b) {
			s += e.total
		}
	}
	return s
}

// Accrued es el bloque de devengado de un periodo (§3.2).
type Accrued struct {
	TotalCents Cents    `json:"total_cents"`
	ByPerson   CentsMap `json:"by_person"`
	ByCategory CentsMap `json:"by_category"`
	paper      Cents
}

// PaymentLine es un pago programado dentro del periodo.
type PaymentLine struct {
	PurchaseID  int64    `json:"purchase_id"`
	PurchaseSeq *int     `json:"purchase_seq"`
	Seq         int      `json:"seq"`
	DueDate     Date     `json:"due_date"`
	AmountCents Cents    `json:"amount_cents"`
	Shares      CentsMap `json:"shares"`
}

// Cash es el bloque de efectivo de un periodo (§3.2).
type Cash struct {
	TotalCents        Cents         `json:"total_cents"`
	PaymentsCents     Cents         `json:"payments_cents"`
	CashExpensesCents Cents         `json:"cash_expenses_cents"`
	ByPerson          CentsMap      `json:"by_person"`
	ByCategory        CentsMap      `json:"by_category"`
	PaymentCount      int           `json:"payment_count"`
	Payments          []PaymentLine `json:"payments"`
	paper             Cents
}

// PeriodSummary §3.2.
type PeriodSummary struct {
	Label                string  `json:"label"`
	Start                Date    `json:"start"`
	End                  Date    `json:"end"`
	Days                 int     `json:"days"`
	Accrued              Accrued `json:"accrued"`
	Cash                 Cash    `json:"cash"`
	CumulativeDiffCents  Cents   `json:"cumulative_diff_cents"`
	CoveringGroceriesSeq *int    `json:"covering_groceries_seq"`
}

// PeriodSpec: meses calendario o ventanas declaradas (exactamente uno).
type PeriodSpec struct {
	Months  []Month
	Windows []Window
}

func seqPtr(s int) *int {
	if s == 0 {
		return nil
	}
	return &s
}

func (l *ledger) accruedBlock(a *agg) Accrued {
	if a == nil {
		a = newAgg()
	}
	bp := newCentsMap(l.in.keys)
	for p, v := range a.byPerson {
		bp.Vals[p] += v
	}
	return Accrued{TotalCents: a.total, ByPerson: bp, ByCategory: catMap(a.byCat), paper: a.paper}
}

func (l *ledger) cashBlock(start, end Date) Cash {
	c := Cash{ByPerson: newCentsMap(l.in.keys), Payments: []PaymentLine{}}
	cats := map[string]Cents{}
	for _, e := range l.cash {
		if !e.date.in(start, end) {
			continue
		}
		c.TotalCents += e.total
		cats[e.category] += e.total
		for _, s := range e.shares {
			c.ByPerson.Vals[s.Person] += s.AmountCents
		}
		if e.payment {
			c.PaymentsCents += e.total
			sh := newCentsMap(l.in.keys)
			for _, s := range e.shares {
				sh.Vals[s.Person] += s.AmountCents
			}
			c.Payments = append(c.Payments, PaymentLine{e.purchase, seqPtr(e.pseq), e.seq, e.date, e.total, sh})
		} else {
			c.CashExpensesCents += e.total
			if e.paperCash {
				c.paper += e.total
			}
		}
	}
	sort.SliceStable(c.Payments, func(a, b int) bool {
		x, y := c.Payments[a], c.Payments[b]
		if x.DueDate != y.DueDate {
			return x.DueDate.Before(y.DueDate)
		}
		if x.PurchaseID != y.PurchaseID {
			return x.PurchaseID < y.PurchaseID
		}
		return x.Seq < y.Seq
	})
	c.PaymentCount = len(c.Payments)
	c.ByCategory = catMap(cats)
	return c
}

func (l *ledger) periods(spec PeriodSpec) ([]PeriodSummary, error) {
	if (len(spec.Months) == 0) == (len(spec.Windows) == 0) {
		return nil, errf(CodeInvalidArgument, "periodos: exactamente uno de months o windows, no vacío")
	}
	var out []PeriodSummary
	if len(spec.Months) > 0 {
		ms := l.monthly()
		for _, m := range spec.Months {
			s, e := MonthBounds(m)
			a := ms[m.idx()]
			ps := PeriodSummary{Label: m.Label(), Start: s, End: e, Days: m.Days(), Accrued: l.accruedBlock(a), Cash: l.cashBlock(s, e)}
			ps.CumulativeDiffCents = accruedUpTo(ms, m.idx()) - l.cashUpTo(e)
			if a != nil {
				ps.CoveringGroceriesSeq = seqPtr(a.maxSeq)
			}
			out = append(out, ps)
		}
		return out, nil
	}
	for i, w := range spec.Windows {
		if w.End.Before(w.Start) {
			return nil, errf(CodeInvalidArgument, "ventana %q termina antes de empezar", w.Label)
		}
		// PROVISIONAL (docs/estado/bloqueos/S2-ventanas.md, pregunta 3): las ventanas van en orden
		// ascendente y disjuntas. El devengo en «n ventanas consecutivas» y la diferencia acumulada
		// como suma corrida solo valen en ese orden; desordenadas, traslapadas o contenidas cambian
		// cifras sin regla escrita, así que se rechazan (start[i] ≤ end[i−1]).
		if i == 0 {
			continue
		}
		if prev := spec.Windows[i-1]; !w.Start.After(prev.End) {
			return nil, errf(CodeInvalidArgument, "ventana %q empieza el %s y la anterior %q termina el %s: van ascendentes y disjuntas (PROVISIONAL, S2-ventanas)", w.Label, w.Start, prev.Label, prev.End)
		}
	}
	if err := l.checkCovered(spec.Windows); err != nil {
		return nil, err
	}
	ws := l.windowed(spec.Windows)
	var cum Cents
	for i, w := range spec.Windows {
		cum += ws[i].total
		ps := PeriodSummary{Label: w.Label, Start: w.Start, End: w.End, Days: daysIncl(w.Start, w.End),
			Accrued: l.accruedBlock(ws[i]), Cash: l.cashBlock(w.Start, w.End), CoveringGroceriesSeq: seqPtr(ws[i].maxSeq)}
		ps.CumulativeDiffCents = cum - l.cashUpTo(w.End)
		out = append(out, ps)
	}
	return out, nil
}

func prepared(in Input) (*Input, error) {
	p := in
	if err := p.Prepare(); err != nil {
		return nil, err
	}
	return &p, nil
}

func records(p *Input) *ledger { return buildLedger(p, p.Expenses, p.CreditPurchases, nil, nil) }

// Periods calcula devengado y efectivo por mes calendario o por ventana (§3.2).
func Periods(in Input, spec PeriodSpec) ([]PeriodSummary, error) {
	p, err := prepared(in)
	if err != nil {
		return nil, err
	}
	return records(p).periods(spec)
}
