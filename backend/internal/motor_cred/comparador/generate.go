package comparador

import (
	"context"
	"sort"
	"strings"

	mc "finanzas-hogar/internal/motor_cred"
)

// storeChoice es la tienda de un candidato y su canasta (nil si no hay).
type storeChoice struct {
	id       *int64
	merchant string
	basket   *BasketQuote
}

// candidate es un escenario enumerado y ya cotizado con el núcleo.
type candidate struct {
	ci    int // índice del crédito en la lista ordenada por id
	in    mc.PurchaseInput
	quote mc.Purchase
	store storeChoice
}

type option struct {
	term  int
	promo *int64
}

// sameMerchant: sin espacios a los lados y sin mayúsculas ASCII (como COLLATE NOCASE).
func sameMerchant(a, b string) bool {
	low := func(s string) string {
		bs := []byte(strings.TrimSpace(s))
		for i, ch := range bs {
			if ch >= 'A' && ch <= 'Z' {
				bs[i] = ch + 32
			}
		}
		return string(bs)
	}
	return low(a) == low(b)
}

// storeChoices: la tienda de la petición si trae store_id o merchant_name; si no, cada canasta
// (PROVISIONAL S3-objetivos-medida (1)); sin ninguna, un candidato sin tienda.
func storeChoices(r Request) []storeChoice {
	if r.StoreID != nil || r.Merchant != "" {
		sc := storeChoice{id: r.StoreID, merchant: r.Merchant}
		for i := range r.Baskets {
			b := &r.Baskets[i]
			if (r.StoreID != nil && b.StoreID == *r.StoreID) || (r.StoreID == nil && sameMerchant(b.Merchant, r.Merchant)) {
				sc.basket = b
				break
			}
		}
		return []storeChoice{sc}
	}
	if len(r.Baskets) == 0 {
		return []storeChoice{{}}
	}
	out := make([]storeChoice, 0, len(r.Baskets))
	for i := range r.Baskets {
		id := r.Baskets[i].StoreID
		out = append(out, storeChoice{id: &id, merchant: r.Baskets[i].Merchant, basket: &r.Baskets[i]})
	}
	return out
}

func applicable(c mc.Credit, p mc.Promotion, d mc.Date, amount mc.Cents, sc storeChoice) bool {
	for _, q := range mc.ApplicablePromotions(c, d, amount, p.TermCount, sc.id, sc.merchant) {
		if q.ID == p.ID {
			return true
		}
	}
	return false
}

func hasRate(c mc.Credit, term int) bool {
	for _, t := range c.TermRates {
		if t.TermCount == term {
			return true
		}
	}
	return false
}

// options: plazos y promociones por tipo de crédito (§8.2, §5.8).
func options(c mc.Credit, d mc.Date, amount mc.Cents, sc storeChoice) []option {
	promos := append([]mc.Promotion(nil), c.Promotions...)
	sort.SliceStable(promos, func(i, j int) bool { return promos[i].ID < promos[j].ID })
	var out []option
	switch c.Kind {
	case mc.KindCard: // el cargo normal siempre se genera; más cada MSI aplicable
		out = append(out, option{1, nil})
		for _, p := range promos {
			if applicable(c, p, d, amount, sc) {
				id := p.ID
				out = append(out, option{p.TermCount, &id})
			}
		}
	case mc.KindRevolving, mc.KindStore:
		terms := map[int]bool{}
		for _, t := range c.TermRates {
			terms[t.TermCount] = true
		}
		var appl []mc.Promotion
		for _, p := range promos {
			if applicable(c, p, d, amount, sc) {
				terms[p.TermCount] = true
				appl = append(appl, p)
			}
		}
		keys := make([]int, 0, len(terms))
		for t := range terms {
			keys = append(keys, t)
		}
		sort.Ints(keys)
		for _, t := range keys {
			if hasRate(c, t) {
				out = append(out, option{t, nil})
			}
			for _, p := range appl {
				if p.TermCount == t {
					id := p.ID
					out = append(out, option{t, &id})
				}
			}
		}
	}
	return out
}

// generate enumera y cotiza los candidatos; los que QuotePurchase rechaza no se generan.
func generate(ctx context.Context, r Request, w int, credits []mc.CreditState) ([]candidate, []ExcludedCredit, error) {
	var cands []candidate
	excluded := []ExcludedCredit{}
	stores := storeChoices(r)
	for ci, cs := range credits {
		c := cs.Credit
		if c.Limit == nil {
			excluded = append(excluded, ExcludedCredit{c.ID, ReasonNoLimit})
			continue
		}
		for k := -w; k <= w; k++ {
			d := r.DesiredDate.AddDays(k)
			if d.Compare(r.Today) < 0 {
				continue
			}
			for _, sc := range stores {
				for _, op := range options(c, d, r.Amount, sc) {
					if err := checkCtx(ctx); err != nil {
						return nil, nil, err
					}
					in := mc.PurchaseInput{Credit: c, Date: d, Principal: r.Amount, TermCount: op.term,
						PromotionID: op.promo, StoreID: sc.id, Merchant: sc.merchant, Shares: r.Shares}
					q, err := mc.QuotePurchase(in)
					if err != nil || len(q.Payments) == 0 {
						continue
					}
					cands = append(cands, candidate{ci, in, q, sc})
				}
			}
		}
	}
	return cands, excluded, nil
}
