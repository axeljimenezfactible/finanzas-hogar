package comparador

import "math/big"

// objective describe uno de los 4 objetivos de K17.
type objective struct {
	value        func(s *Scenario) (int64, bool) // false = nulo
	lessIsBetter bool
}

var objectives = [4]objective{
	{func(s *Scenario) (int64, bool) { return int64(s.Interest), true }, true},
	{func(s *Scenario) (int64, bool) { return s.StdDevMilliCents, true }, true},
	{func(s *Scenario) (int64, bool) { return int64(s.MinAvailable), true }, false},
	{func(s *Scenario) (int64, bool) {
		if s.ValuePerPesoMilli == nil {
			return 0, false
		}
		return *s.ValuePerPesoMilli, true
	}, false},
}

func weightsOf(w Weights) [4]int64 {
	return [4]int64{w.LowestInterest, w.EvenLowPayment, w.MaxAvailable, w.MostValuePerPeso}
}

// normalize: mín-máx entre los escenarios generados (valid); empate o todos nulos → 1 (S3-14).
func normalize(valid []Scenario, o objective) []*big.Rat {
	out := make([]*big.Rat, len(valid))
	var hi, lo int64
	have := false
	for i := range valid {
		if v, ok := o.value(&valid[i]); ok {
			if !have {
				hi, lo, have = v, v, true
			}
			hi, lo = max(hi, v), min(lo, v)
		}
	}
	for i := range valid {
		v, ok := o.value(&valid[i])
		switch {
		case !ok || !have || hi == lo:
			out[i] = big.NewRat(1, 1)
		case o.lessIsBetter:
			out[i] = big.NewRat(hi-v, hi-lo)
		default:
			out[i] = big.NewRat(v-lo, hi-lo)
		}
	}
	return out
}

// scores: Σ wᵢ·nᵢ / Σ wᵢ exacto por escenario; llena ScoreBP = HALF_UP(10000 × puntaje).
func scores(valid []Scenario, w Weights) []*big.Rat {
	ws := weightsOf(w)
	var sum int64
	for _, x := range ws {
		sum += x
	}
	norms := [4][]*big.Rat{}
	for j, o := range objectives {
		norms[j] = normalize(valid, o)
	}
	out := make([]*big.Rat, len(valid))
	for i := range valid {
		acc := new(big.Rat)
		for j := range objectives {
			acc.Add(acc, new(big.Rat).Mul(big.NewRat(ws[j], 1), norms[j][i]))
		}
		acc.Quo(acc, big.NewRat(sum, 1))
		out[i] = acc
		bp := halfUp(new(big.Int).Mul(acc.Num(), big.NewInt(10000)), acc.Denom())
		valid[i].ScoreBP = &bp
	}
	return out
}

// keyCmp compara por la llave canónica (credit_id, term_count, purchase_date, store_id nil
// primero, promotion_id nil primero).
func keyCmp(a, b *Scenario) int {
	ptr := func(x, y *int64) int {
		switch {
		case x == nil && y == nil:
			return 0
		case x == nil:
			return -1
		case y == nil:
			return 1
		case *x < *y:
			return -1
		case *x > *y:
			return 1
		}
		return 0
	}
	switch {
	case a.CreditID != b.CreditID:
		return cmpInt(a.CreditID, b.CreditID)
	case a.TermCount != b.TermCount:
		return cmpInt(int64(a.TermCount), int64(b.TermCount))
	}
	if c := a.PurchaseDate.Compare(b.PurchaseDate); c != 0 {
		return c
	}
	if c := ptr(a.StoreID, b.StoreID); c != 0 {
		return c
	}
	return ptr(a.PromotionID, b.PromotionID)
}

func cmpInt(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// best elige el óptimo del objetivo j: mejor valor (nulo al final); empate → (objetivo 2)
// menos personas sobre el tope → mayor puntaje exacto → llave canónica (§8.6).
func best(valid []Scenario, exact []*big.Rat, j int) int {
	o := objectives[j]
	bi := 0
	for i := 1; i < len(valid); i++ {
		if better(valid, exact, o, j, i, bi) {
			bi = i
		}
	}
	return bi
}

func better(valid []Scenario, exact []*big.Rat, o objective, j, a, b int) bool {
	va, oka := o.value(&valid[a])
	vb, okb := o.value(&valid[b])
	if oka != okb {
		return oka
	}
	if oka && va != vb {
		return (va < vb) == o.lessIsBetter
	}
	if j == 1 && len(valid[a].OverCap) != len(valid[b].OverCap) {
		return len(valid[a].OverCap) < len(valid[b].OverCap)
	}
	if c := exact[a].Cmp(exact[b]); c != 0 {
		return c > 0
	}
	return keyCmp(&valid[a], &valid[b]) < 0
}

// recommend: mayor puntaje exacto; empate → llave canónica (el tope no desempata aquí).
func recommend(valid []Scenario, exact []*big.Rat) int {
	bi := 0
	for i := 1; i < len(valid); i++ {
		if c := exact[i].Cmp(exact[bi]); c > 0 || (c == 0 && keyCmp(&valid[i], &valid[bi]) < 0) {
			bi = i
		}
	}
	return bi
}
