package dominio

import "testing"

func TestMonthHasData(t *testing.T) {
	for _, c := range []struct {
		accrued, cash int64
		want          bool
	}{{0, 0, false}, {119764, 0, true}, {0, 59882, true}, {119764, 119764, true}, {-500, 500, true}} {
		if got := MonthHasData(c.accrued, c.cash); got != c.want {
			t.Errorf("MonthHasData(%d, %d) = %v, se esperaba %v", c.accrued, c.cash, got, c.want)
		}
	}
}
