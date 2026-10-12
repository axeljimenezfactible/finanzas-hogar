package motor_cred

import "sort"

// CalendarEntry es una entrada del calendario unificado (§6).
type CalendarEntry struct {
	Date          Date           `json:"date"`
	CreditID      int64          `json:"credit_id"`
	CreditKind    Kind           `json:"credit_kind"`
	Source        string         `json:"source"`
	PurchaseRef   string         `json:"purchase_ref"`
	Seq           int            `json:"seq"`
	StatementDate *Date          `json:"statement_date"`
	Amount        Cents          `json:"amount_cents"`
	Principal     Cents          `json:"principal_cents"`
	Interest      Cents          `json:"interest_cents"`
	Shares        []PaymentShare `json:"shares"`
	Status        string         `json:"status"`
	Provisional   []string       `json:"provisional"`
}

// FortnightTotal es el total de pagos de una quincena.
type FortnightTotal struct {
	PayDay   Date           `json:"pay_day"`
	Total    Cents          `json:"total_cents"`
	ByPerson []PaymentShare `json:"by_person"`
}

// Calendar es el calendario unificado de pagos.
type Calendar struct {
	Entries     []CalendarEntry  `json:"entries"`
	ByFortnight []FortnightTotal `json:"by_fortnight"`
}

func inRange(d, from, to Date) bool { return d.Compare(from) >= 0 && d.Compare(to) <= 0 }

func status(paid, amount Cents) string {
	switch {
	case paid <= 0:
		return "pending"
	case paid < amount:
		return "partial"
	}
	return "paid"
}

func paidFor(st CreditState, ref string, seq int) Cents {
	var s Cents
	for _, r := range st.recordsFor(ref, seq) {
		s += r.Amount
	}
	return s
}

// UnifiedCalendar reúne los pagos de todos los créditos en [from, to] (S3-11). En tarjeta, un
// corte sustituye a los pagos programados de su vencimiento (PROVISIONAL, S3-formula-tarjeta).
func UnifiedCalendar(credits []CreditState, from, to Date) (Calendar, error) {
	if to.Compare(from) < 0 {
		return Calendar{}, invalid("ventana inválida")
	}
	cs, err := sortedCredits(credits)
	if err != nil {
		return Calendar{}, err
	}
	var entries []CalendarEntry
	for _, st := range cs {
		qs, err := st.quoted()
		if err != nil {
			return Calendar{}, err
		}
		c := st.Credit
		dueSt, stSeen := map[Date]bool{}, map[Date]bool{}
		for _, s := range st.Statements {
			if dueSt[s.DueDate] || stSeen[s.StatementDate] {
				return Calendar{}, invalid("dos cortes con la misma fecha")
			}
			dueSt[s.DueDate], stSeen[s.StatementDate] = true, true
		}
		for _, s := range st.Statements {
			if !inRange(s.DueDate, from, to) {
				continue
			}
			sd := s.StatementDate
			e := CalendarEntry{Date: s.DueDate, CreditID: c.ID, CreditKind: c.Kind, Source: "statement",
				StatementDate: &sd, Amount: s.Balance, Provisional: []string{provCard}}
			ext := map[string]Cents{}
			var order []string
			var paid Cents
			for _, q := range qs {
				if !q.live() {
					continue
				}
				for _, p := range q.Purchase.Payments {
					if p.Due != s.DueDate {
						continue
					}
					e.Principal += p.Principal
					paid += paidFor(st, q.Ref, p.Seq)
					for _, sh := range p.Shares {
						if sh.Code == c.Holder {
							continue
						}
						if _, ok := ext[sh.Code]; !ok {
							order = append(order, sh.Code)
						}
						ext[sh.Code] += sh.Amount
					}
				}
			}
			e.Interest = e.Amount - e.Principal
			holder := e.Amount
			sort.Strings(order)
			shares := []PaymentShare{}
			for _, code := range order {
				holder -= ext[code]
				shares = append(shares, PaymentShare{Code: code, Amount: ext[code]})
			}
			if e.Interest < 0 || holder < 0 {
				return Calendar{}, invalid("el corte es menor que los pagos que incluye")
			}
			e.Shares = append([]PaymentShare{{PersonID: c.HolderPersonID, Code: c.Holder, Amount: holder}}, shares...)
			e.Status = status(paid, e.Amount)
			entries = append(entries, e)
		}
		for _, q := range qs {
			if !q.live() {
				continue
			}
			for _, p := range q.Purchase.Payments {
				if !inRange(p.Due, from, to) || dueSt[p.Due] {
					continue
				}
				entries = append(entries, CalendarEntry{Date: p.Due, CreditID: c.ID, CreditKind: c.Kind,
					Source: "scheduled", PurchaseRef: q.Ref, Seq: p.Seq, Amount: p.Amount, Principal: p.Principal,
					Interest: p.Interest, Shares: append([]PaymentShare{}, p.Shares...),
					Status: status(paidFor(st, q.Ref, p.Seq), p.Amount), Provisional: q.Purchase.Provisional})
			}
		}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		x, y := entries[i], entries[j]
		if c := x.Date.Compare(y.Date); c != 0 {
			return c < 0
		}
		if x.CreditID != y.CreditID {
			return x.CreditID < y.CreditID
		}
		if x.Source != y.Source {
			return x.Source == "scheduled"
		}
		if x.PurchaseRef != y.PurchaseRef {
			return x.PurchaseRef < y.PurchaseRef
		}
		return x.Seq < y.Seq
	})
	cal := Calendar{Entries: []CalendarEntry{}, ByFortnight: []FortnightTotal{}}
	cal.Entries = append(cal.Entries, entries...)
	idx := map[Date]int{}
	for _, e := range entries {
		pd := FortnightOf(e.Date).PayDay
		i, ok := idx[pd]
		if !ok {
			i = len(cal.ByFortnight)
			idx[pd] = i
			cal.ByFortnight = append(cal.ByFortnight, FortnightTotal{PayDay: pd})
		}
		cal.ByFortnight[i].Total += e.Amount
		cal.ByFortnight[i].ByPerson = addShares(cal.ByFortnight[i].ByPerson, e.Shares)
	}
	return cal, nil
}

// addShares suma partes por persona, en orden de person.
func addShares(acc []PaymentShare, add []PaymentShare) []PaymentShare {
	for _, s := range add {
		found := false
		for i := range acc {
			if acc[i].Code == s.Code {
				acc[i].Amount += s.Amount
				found = true
			}
		}
		if !found {
			acc = append(acc, PaymentShare{PersonID: s.PersonID, Code: s.Code, Amount: s.Amount})
		}
	}
	sort.SliceStable(acc, func(i, j int) bool { return acc[i].Code < acc[j].Code })
	return acc
}

// CashWindow es el efectivo de una ventana (§6.1).
type CashWindow struct {
	From         Date           `json:"from"`
	To           Date           `json:"to"`
	Payments     Cents          `json:"payments_cents"`
	Cash         Cents          `json:"cash_cents"`
	Total        Cents          `json:"total_cents"`
	PaymentCount int            `json:"payment_count"`
	ByPerson     []PaymentShare `json:"by_person"`
}

// CashFlow suma los pagos del calendario y el contado en [from, to].
func CashFlow(cal Calendar, cash []CashEvent, from, to Date) CashWindow {
	w := CashWindow{From: from, To: to, ByPerson: []PaymentShare{}}
	for _, e := range cal.Entries {
		if inRange(e.Date, from, to) {
			w.Payments += e.Amount
			w.PaymentCount++
			w.ByPerson = addShares(w.ByPerson, e.Shares)
		}
	}
	for _, ce := range cash {
		if inRange(ce.Date, from, to) {
			w.Cash += ce.Amount
			w.ByPerson = addShares(w.ByPerson, []PaymentShare{{Code: ce.Person, Amount: ce.Amount}})
		}
	}
	w.Total = w.Payments + w.Cash
	return w
}
