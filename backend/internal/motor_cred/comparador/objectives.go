package comparador

import (
	"math/big"
	"sort"
	"time"

	mc "finanzas-hogar/internal/motor_cred"
)

// fq es una quincena (RN-06, K8).
type fq struct {
	y int
	m time.Month
	h int
}

func fqOf(d mc.Date) fq {
	f := mc.FortnightOf(d)
	return fq{f.Year, f.Month, f.Half}
}

func (q fq) next() fq {
	if q.h == 1 {
		return fq{q.y, q.m, 2}
	}
	if q.m == time.December {
		return fq{q.y + 1, time.January, 1}
	}
	return fq{q.y, q.m + 1, 1}
}

func (q fq) first() mc.Date {
	if q.h == 1 {
		return mc.Date{Year: q.y, Month: q.m, Day: 1}
	}
	return mc.Date{Year: q.y, Month: q.m, Day: 16}
}

func (q fq) last() mc.Date {
	if q.h == 1 {
		return mc.Date{Year: q.y, Month: q.m, Day: 15}
	}
	return mc.Date{Year: q.y, Month: q.m, Day: time.Date(q.y, q.m+1, 0, 0, 0, 0, 0, time.UTC).Day()}
}

// stdDevMilli = ⌊(⌊isqrt(4N) / n⌋ + 1) / 2⌋ con N = 10⁶ × (n·Σx² − (Σx)²): HALF_UP de 1000 ×
// la desviación estándar poblacional, exacto (§8.4).
func stdDevMilli(xs []mc.Cents) int64 {
	n := big.NewInt(int64(len(xs)))
	sum, sq := new(big.Int), new(big.Int)
	for _, x := range xs {
		bx := big.NewInt(int64(x))
		sum.Add(sum, bx)
		sq.Add(sq, new(big.Int).Mul(bx, bx))
	}
	N := new(big.Int).Mul(n, sq)
	N.Sub(N, new(big.Int).Mul(sum, sum))
	N.Mul(N, big.NewInt(1_000_000))
	r := new(big.Int).Sqrt(N.Lsh(N, 2))
	r.Quo(r, n)
	r.Add(r, big.NewInt(1))
	return r.Rsh(r, 1).Int64()
}

// series arma la serie del objetivo 2: totales del hogar y por persona en cada quincena de
// fqOf(p) a fqOf(l), incluidas las de 0, y la unión de las Provisional de sus entradas.
func series(entries []mc.CalendarEntry, p, l mc.Date) (tot []mc.Cents, byPerson map[string][]mc.Cents, prov []string) {
	idx := map[fq]int{}
	for q := fqOf(p); ; q = q.next() {
		idx[q] = len(tot)
		tot = append(tot, 0)
		if q == fqOf(l) {
			break
		}
	}
	byPerson = map[string][]mc.Cents{}
	for _, e := range entries {
		i, ok := idx[fqOf(e.Date)]
		if !ok {
			continue
		}
		tot[i] += e.Amount
		prov = append(prov, e.Provisional...)
		for _, sh := range e.Shares {
			if byPerson[sh.Code] == nil {
				byPerson[sh.Code] = make([]mc.Cents, len(tot))
			}
			byPerson[sh.Code][i] += sh.Amount
		}
	}
	return tot, byPerson, prov
}

func maxMin(xs []mc.Cents) (mc.Cents, mc.Cents) {
	hi, lo := xs[0], xs[0]
	for _, x := range xs[1:] {
		hi, lo = max(hi, x), min(lo, x)
	}
	return hi, lo
}

// overCap marca a cada persona con tope cuyo pico por quincena pasa su tope (orden por persona).
func overCap(caps []PersonCap, byPerson map[string][]mc.Cents) []PersonPeak {
	out := []PersonPeak{}
	for _, c := range caps {
		var peak mc.Cents
		for _, v := range byPerson[c.Person] {
			peak = max(peak, v)
		}
		if peak > c.Cap {
			out = append(out, PersonPeak{c.Person, peak, c.Cap})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Person < out[j].Person })
	return out
}

// halfUp: HALF_UP(num/den) con num ≥ 0 y den > 0.
func halfUp(num, den *big.Int) int64 {
	n := new(big.Int).Lsh(num, 1)
	n.Add(n, den)
	return n.Quo(n, new(big.Int).Lsh(den, 1)).Int64()
}

// valuePerPeso = HALF_UP(monto × U × 100 / (C × (monto + intereses))); nil sin canasta (§8.4).
func valuePerPeso(amount, interest mc.Cents, b *BasketQuote) *int64 {
	if b == nil {
		return nil
	}
	num := new(big.Int).Mul(big.NewInt(int64(amount)), big.NewInt(b.BaseUnitsMilli))
	num.Mul(num, big.NewInt(100))
	den := new(big.Int).Mul(big.NewInt(int64(b.Cost)), big.NewInt(int64(amount+interest)))
	v := halfUp(num, den)
	return &v
}

// provList: sin repetidos y en orden alfabético por bytes (§14.1).
func provList(names ...[]string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, ns := range names {
		for _, n := range ns {
			if !seen[n] {
				seen[n] = true
				out = append(out, n)
			}
		}
	}
	sort.Strings(out)
	return out
}
