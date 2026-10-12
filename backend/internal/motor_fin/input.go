package motor_fin

import (
	"sort"
	"time"
)

// Input es la raíz de entradas (interfaz §2; en las fixtures, entradas.datos).
type Input struct {
	Params           Params            `json:"params"`
	Persons          []Person          `json:"persons"`
	Categories       []Category        `json:"categories"`
	Expenses         []Expense         `json:"expenses"`
	CreditPurchases  []CreditPurchase  `json:"credit_purchases"`
	Incomes          []Income          `json:"incomes"`
	PlannedItems     []PlannedItem     `json:"planned_items"`
	Budgets          []BudgetRow       `json:"budgets"`
	SavingsGoals     []SavingsGoal     `json:"savings_goals"`
	SavingsMovements []SavingsMovement `json:"savings_movements"`
	Cadence          Cadence           `json:"cadence"`
	StableItems      []StableItem      `json:"stable_items"`

	ready    bool
	today    Date
	holder   string
	keys     []string // titular y externas, en el orden de persons (§1.1)
	byCode   map[string]*Person
	warnPct  int64
	limitPct int64
}

// Params §2.1.
type Params struct {
	Today                *Date      `json:"today"`
	Now                  *time.Time `json:"now"`
	Zone                 string     `json:"zone"`
	Holder               string     `json:"holder"`
	GroceriesCategory    string     `json:"groceries_category"`
	PaperCategory        string     `json:"paper_category"`
	RepeatCadenceEnabled bool       `json:"repeat_cadence_enabled"`
	BudgetWarningPct     *int64     `json:"budget_warning_pct"`
	BudgetLimitPct       *int64     `json:"budget_limit_pct"`
}

// Person §2.2.
type Person struct {
	Code              string `json:"code"`
	DisplayName       string `json:"display_name"`
	Participation     string `json:"participation"`
	PaysShareDirectly bool   `json:"pays_share_directly"`
}

// Category §2.3.
type Category struct {
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Kind     *string `json:"kind"`
	IsActive *bool   `json:"is_active"`
}

// Share es una parte congelada {person, amount_cents}.
type Share struct {
	Person      string `json:"person"`
	AmountCents Cents  `json:"amount_cents"`
}

// Expense §2.4.
type Expense struct {
	ID               int64   `json:"id"`
	Kind             string  `json:"kind"`
	Category         string  `json:"category"`
	AmountCents      Cents   `json:"amount_cents"`
	ExpenseDate      Date    `json:"expense_date"`
	AccrualMonths    *int    `json:"accrual_months"`
	PaymentMethod    string  `json:"payment_method"`
	PaidBy           *string `json:"paid_by"`
	Shares           []Share `json:"shares"`
	CreditPurchaseID *int64  `json:"credit_purchase_id"`
	PlannedItemID    *int64  `json:"planned_item_id"`
	Projected        bool    `json:"projected"`
}

// PurchaseShare es la parte de una persona en una compra.
type PurchaseShare struct {
	Person         string `json:"person"`
	PrincipalCents Cents  `json:"principal_cents"`
	InterestCents  Cents  `json:"interest_cents"`
	TotalCents     Cents  `json:"total_cents"`
}

// Embedded es un artículo marcado dentro de una compra (papel dentro de la despensa).
type Embedded struct {
	Category    string `json:"category"`
	Person      string `json:"person"`
	AmountCents Cents  `json:"amount_cents"`
}

// Paid es un abono registrado a un pago programado.
type Paid struct {
	Person      string `json:"person"`
	AmountCents Cents  `json:"amount_cents"`
	PaidDate    Date   `json:"paid_date"`
}

// ScheduledPayment es un pago del calendario de S3.
type ScheduledPayment struct {
	Seq            int     `json:"seq"`
	DueDate        Date    `json:"due_date"`
	AmountCents    Cents   `json:"amount_cents"`
	PrincipalCents Cents   `json:"principal_cents"`
	InterestCents  Cents   `json:"interest_cents"`
	Shares         []Share `json:"shares"`
	Paid           []Paid  `json:"paid"`
}

// CreditPurchase §2.5.
type CreditPurchase struct {
	ID             int64              `json:"id"`
	CreditID       int64              `json:"credit_id"`
	ExpenseID      *int64             `json:"expense_id"`
	Category       *string            `json:"category"`
	PurchaseDate   Date               `json:"purchase_date"`
	PrincipalCents Cents              `json:"principal_cents"`
	InterestCents  Cents              `json:"interest_cents"`
	TotalCents     Cents              `json:"total_cents"`
	AccrualMonths  *int               `json:"accrual_months"`
	Status         string             `json:"status"`
	Shares         []PurchaseShare    `json:"shares"`
	Embedded       []Embedded         `json:"embedded"`
	Payments       []ScheduledPayment `json:"payments"`
	Projected      bool               `json:"projected"`
}

// Income §2.6.
type Income struct {
	ID            int64  `json:"id"`
	Person        string `json:"person"`
	Kind          string `json:"kind"`
	AmountCents   Cents  `json:"amount_cents"`
	IncomeDate    Date   `json:"income_date"`
	PlannedItemID *int64 `json:"planned_item_id"`
}

// RulePart es una parte de una regla de reparto: porcentaje en puntos base o monto fijo.
type RulePart struct {
	Person     string `json:"person"`
	PctBP      *int64 `json:"pct_bp"`
	FixedCents *Cents `json:"fixed_cents"`
}

// ShareRule es la regla efectiva ya resuelta (RN-03).
type ShareRule struct {
	Parts      []RulePart `json:"parts"`
	RestPerson string     `json:"rest_person"`
}

// PlannedItem §2.7.
type PlannedItem struct {
	ID                    int64      `json:"id"`
	Kind                  string     `json:"kind"`
	Description           string     `json:"description"`
	AmountCents           Cents      `json:"amount_cents"`
	SecondHalfAmountCents *Cents     `json:"second_half_amount_cents"`
	Frequency             string     `json:"frequency"`
	IntervalCount         *int       `json:"interval_count"`
	StartDate             Date       `json:"start_date"`
	EndDate               *Date      `json:"end_date"`
	Category              *string    `json:"category"`
	PaymentMethod         *string    `json:"payment_method"`
	AccrualMonths         *int       `json:"accrual_months"`
	Person                *string    `json:"person"`
	ShareRule             *ShareRule `json:"share_rule"`
	CreditID              *int64     `json:"credit_id"`
	IsActive              *bool      `json:"is_active"`
}

// BudgetRow §2.8.
type BudgetRow struct {
	ID          int64   `json:"id"`
	Category    *string `json:"category"`
	AmountCents Cents   `json:"amount_cents"`
	ValidFrom   Date    `json:"valid_from"`
	ValidTo     *Date   `json:"valid_to"`
}

// BudgetChange §2.8 (argumento de ApplyBudgetChange).
type BudgetChange struct {
	RowID       int64  `json:"row_id"`
	Month       Month  `json:"month"`
	AmountCents Cents  `json:"amount_cents"`
	Scope       string `json:"scope"`
}

// SavingsGoal §2.9.
type SavingsGoal struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	Kind                string `json:"kind"`
	TargetCents         *Cents `json:"target_cents"`
	TargetDate          *Date  `json:"target_date"`
	PlannedMonthlyCents *Cents `json:"planned_monthly_cents"`
	StartDate           Date   `json:"start_date"`
	IsActive            *bool  `json:"is_active"`
}

// SavingsMovement §2.9.
type SavingsMovement struct {
	GoalID       int64  `json:"goal_id"`
	MovementDate Date   `json:"movement_date"`
	Kind         string `json:"kind"`
	AmountCents  Cents  `json:"amount_cents"`
}

// Cadence §2.10 (PROVISIONAL, docs/estado/bloqueos/S2-proyeccion-interruptor.md).
type Cadence struct {
	Purchases []CreditPurchase `json:"purchases"`
	Expenses  []Expense        `json:"expenses"`
}

// StableShare es la parte de una persona en una partida del valor estable.
type StableShare struct {
	Person     string `json:"person"`
	TotalCents Cents  `json:"total_cents"`
}

// StableItem §2.11.
type StableItem struct {
	Label       string        `json:"label"`
	CycleMonths int           `json:"cycle_months"`
	TotalCents  Cents         `json:"total_cents"`
	Shares      []StableShare `json:"shares"`
}

// Window §2.12: ventana declarada, inclusive.
type Window struct {
	Label string `json:"label"`
	Start Date   `json:"start"`
	End   Date   `json:"end"`
}

func active(b *bool) bool { return b == nil || *b }

// DecodeInput lee el JSON de entradas con campos estrictos y lo valida.
func DecodeInput(b []byte) (Input, error) {
	var in Input
	if err := decodeStrict(b, &in); err != nil {
		return in, err
	}
	if err := in.Prepare(); err != nil {
		return in, err
	}
	return in, nil
}

// Today devuelve el «hoy» resuelto (válido después de Prepare).
func (in *Input) Today() Date { return in.today }

// Prepare resuelve hoy, valores por omisión y valida las entradas (§1.4, §2).
func (in *Input) Prepare() error {
	if in.ready {
		return nil
	}
	p := &in.Params
	if p.Zone != "" && p.Zone != DefaultZone {
		return errf(CodeInvalidZone, "zona no admitida %q", p.Zone)
	}
	switch {
	case p.Today != nil && p.Now == nil:
		in.today = *p.Today
	case p.Now != nil && p.Today == nil:
		t, err := Today(*p.Now, p.Zone)
		if err != nil {
			return err
		}
		in.today = t
	default:
		return errf(CodeInvalidArgument, "params: exactamente uno de today o now")
	}
	if p.GroceriesCategory == "" || p.PaperCategory == "" || p.Holder == "" {
		return errf(CodeInvalidArgument, "params: holder, groceries_category y paper_category son obligatorios")
	}
	in.warnPct, in.limitPct = 80, 100
	if p.BudgetWarningPct != nil {
		in.warnPct = *p.BudgetWarningPct
	}
	if p.BudgetLimitPct != nil {
		in.limitPct = *p.BudgetLimitPct
	}
	if in.warnPct < 0 || in.limitPct < 0 {
		return errf(CodeInvalidArgument, "params: umbrales negativos")
	}
	in.byCode = map[string]*Person{}
	holders := 0
	for i := range in.Persons {
		ps := &in.Persons[i]
		if _, dup := in.byCode[ps.Code]; dup || ps.Code == "" {
			return errf(CodeInvalidArgument, "persons[%d]: código vacío o repetido %q", i, ps.Code)
		}
		in.byCode[ps.Code] = ps
		switch ps.Participation {
		case "holder":
			holders++
			in.holder = ps.Code
			in.keys = append(in.keys, ps.Code)
		case "external":
			in.keys = append(in.keys, ps.Code)
		case "co_holder":
		default:
			return errf(CodeInvalidArgument, "persons[%d]: participation %q", i, ps.Participation)
		}
	}
	if holders != 1 || in.holder != p.Holder {
		return errf(CodeInvalidArgument, "persons: debe haber exactamente un titular y ser params.holder")
	}
	for i := range in.Expenses {
		if err := in.checkExpense(&in.Expenses[i], "expenses", i, in.CreditPurchases); err != nil {
			return err
		}
	}
	for i := range in.CreditPurchases {
		if err := in.checkPurchase(&in.CreditPurchases[i], "credit_purchases", i); err != nil {
			return err
		}
	}
	allP := append(append([]CreditPurchase(nil), in.CreditPurchases...), in.Cadence.Purchases...)
	for i := range in.Cadence.Purchases {
		if err := in.checkPurchase(&in.Cadence.Purchases[i], "cadence.purchases", i); err != nil {
			return err
		}
	}
	for i := range in.Cadence.Expenses {
		if err := in.checkExpense(&in.Cadence.Expenses[i], "cadence.expenses", i, allP); err != nil {
			return err
		}
	}
	for i, x := range in.Incomes {
		if x.Person != in.holder {
			return errf(CodeIncomeNotHolder, "incomes[%d]: %q no es el titular", i, x.Person)
		}
		if (x.Kind != "fixed" && x.Kind != "extra") || x.AmountCents <= 0 {
			return errf(CodeInvalidArgument, "incomes[%d]: kind o monto inválido", i)
		}
	}
	for i := range in.PlannedItems {
		if err := in.checkPlanned(&in.PlannedItems[i], i); err != nil {
			return err
		}
	}
	if err := checkBudgets(in.Budgets); err != nil {
		return err
	}
	for i, g := range in.SavingsGoals {
		if g.Kind != "fund" && g.Kind != "goal" {
			return errf(CodeInvalidArgument, "savings_goals[%d]: kind %q", i, g.Kind)
		}
		if g.Kind == "goal" && (g.TargetCents == nil || g.TargetDate == nil) {
			return errf(CodeInvalidArgument, "savings_goals[%d]: una meta lleva monto y fecha", i)
		}
	}
	for i, m := range in.SavingsMovements {
		if (m.Kind != "deposit" && m.Kind != "withdrawal") || m.AmountCents <= 0 {
			return errf(CodeInvalidArgument, "savings_movements[%d]: kind o monto inválido", i)
		}
	}
	if err := in.checkStable(in.StableItems); err != nil {
		return err
	}
	in.ready = true
	return nil
}

// money valida que una persona de un campo de dinero sea el titular o externa (C-09).
func (in *Input) money(code, where string) error {
	ps, ok := in.byCode[code]
	if !ok || ps.Participation == "co_holder" {
		return errf(CodeInvalidArgument, "%s: persona %q no es titular ni externa", where, code)
	}
	return nil
}

// countOK valida un entero opcional 1–24 (§2.4, §2.5, §2.7: accrual_months e interval_count).
// Ausente o null vale 1; un 0 explícito o más de 24 es invalid_argument, como en el oráculo.
func countOK(p *int) bool { return p == nil || (*p >= 1 && *p <= 24) }

// count1 devuelve un entero opcional 1–24 ya validado por countOK; ausente o null vale 1.
func count1(p *int) int {
	if p == nil {
		return 1
	}
	return *p
}

func (in *Input) checkExpense(e *Expense, list string, i int, purchases []CreditPurchase) error {
	w := func(f string) string { return list + "[" + itoa(i) + "]." + f }
	if !countOK(e.AccrualMonths) {
		return errf(CodeInvalidArgument, "%s fuera de 1–24", w("accrual_months"))
	}
	if e.Kind != "fixed" && e.Kind != "variable" && e.Kind != "upcoming" {
		return errf(CodeInvalidArgument, "%s %q", w("kind"), e.Kind)
	}
	if e.PaidBy != nil {
		if err := in.money(*e.PaidBy, w("paid_by")); err != nil {
			return err
		}
	}
	var s Cents
	for _, sh := range e.Shares {
		if err := in.money(sh.Person, w("shares")); err != nil {
			return err
		}
		s += sh.AmountCents
	}
	if s != e.AmountCents {
		return errf(CodeSharesMismatch, "%s suman %d y el importe es %d", w("shares"), s, e.AmountCents)
	}
	switch e.PaymentMethod {
	case "cash", "debit", "transfer":
	case "credit":
		if e.CreditPurchaseID == nil {
			return errf(CodeMissingPurchase, "%s falta", w("credit_purchase_id"))
		}
		var cp *CreditPurchase
		for k := range purchases {
			if purchases[k].ID == *e.CreditPurchaseID {
				cp = &purchases[k]
			}
		}
		if cp == nil {
			return errf(CodeMissingPurchase, "%s %d no existe", w("credit_purchase_id"), *e.CreditPurchaseID)
		}
		if cp.PrincipalCents != e.AmountCents {
			return errf(CodeSumMismatch, "%s distinto del capital de la compra %d", w("amount_cents"), cp.ID)
		}
	default:
		return errf(CodeInvalidArgument, "%s %q", w("payment_method"), e.PaymentMethod)
	}
	return nil
}

func (in *Input) checkPurchase(c *CreditPurchase, list string, i int) error {
	w := func(f string) string { return list + "[" + itoa(i) + "]." + f }
	if !countOK(c.AccrualMonths) {
		return errf(CodeInvalidArgument, "%s fuera de 1–24", w("accrual_months"))
	}
	switch c.Status {
	case "active", "paid_off":
	case "cancelled":
		return nil
	default:
		return errf(CodeInvalidArgument, "%s %q", w("status"), c.Status)
	}
	if c.PrincipalCents+c.InterestCents != c.TotalCents {
		return errf(CodeSumMismatch, "%s capital + intereses ≠ total", w("total_cents"))
	}
	var sp, si, st Cents
	for _, sh := range c.Shares {
		if err := in.money(sh.Person, w("shares")); err != nil {
			return err
		}
		if sh.PrincipalCents+sh.InterestCents != sh.TotalCents {
			return errf(CodeSumMismatch, "%s parte de %s no suma", w("shares"), sh.Person)
		}
		sp, si, st = sp+sh.PrincipalCents, si+sh.InterestCents, st+sh.TotalCents
	}
	if len(c.Shares) > 0 && (sp != c.PrincipalCents || si != c.InterestCents || st != c.TotalCents) {
		return errf(CodeSumMismatch, "%s no suman la compra", w("shares"))
	}
	for _, e := range c.Embedded {
		if err := in.money(e.Person, w("embedded")); err != nil {
			return err
		}
	}
	var sa, sc Cents
	for k, p := range c.Payments {
		pw := w("payments[" + itoa(k) + "]")
		if p.PrincipalCents+p.InterestCents != p.AmountCents {
			return errf(CodeSumMismatch, "%s monto ≠ capital + intereses", pw)
		}
		var ss Cents
		for _, sh := range p.Shares {
			if err := in.money(sh.Person, pw+".shares"); err != nil {
				return err
			}
			ss += sh.AmountCents
		}
		if ss != p.AmountCents {
			return errf(CodeSumMismatch, "%s partes ≠ monto", pw)
		}
		for _, pd := range p.Paid {
			if err := in.money(pd.Person, pw+".paid"); err != nil {
				return err
			}
		}
		sa, sc = sa+p.AmountCents, sc+p.PrincipalCents
	}
	if len(c.Payments) > 0 && (sa != c.TotalCents || sc != c.PrincipalCents) {
		return errf(CodeSumMismatch, "%s Σ montos ≠ total o Σ capital ≠ capital", w("payments"))
	}
	return nil
}

func (in *Input) checkPlanned(p *PlannedItem, i int) error {
	w := "planned_items[" + itoa(i) + "]"
	if !countOK(p.IntervalCount) || !countOK(p.AccrualMonths) {
		return errf(CodeInvalidArgument, "%s: interval_count o accrual_months fuera de 1–24", w)
	}
	switch p.Frequency {
	case "weekly", "fortnightly", "monthly", "yearly", "once":
	default:
		return errf(CodeInvalidArgument, "%s: frequency %q", w, p.Frequency)
	}
	switch p.Kind {
	case "income":
		if p.Person != nil && *p.Person != in.holder {
			return errf(CodeIncomeNotHolder, "%s: ingreso de %q", w, *p.Person)
		}
	case "fixed_expense", "upcoming_expense":
		if p.Person != nil {
			if _, ok := in.byCode[*p.Person]; !ok {
				return errf(CodeInvalidArgument, "%s: persona %q", w, *p.Person)
			}
		}
		if p.CreditID != nil {
			return nil
		}
		r := p.ShareRule
		if r == nil {
			return errf(CodeInvalidArgument, "%s: share_rule obligatoria en gastos", w)
		}
		if err := in.money(r.RestPerson, w+".share_rule"); err != nil {
			return err
		}
		for _, pt := range r.Parts {
			if err := in.money(pt.Person, w+".share_rule"); err != nil {
				return err
			}
			if (pt.PctBP == nil) == (pt.FixedCents == nil) {
				return errf(CodeInvalidArgument, "%s: cada parte lleva pct_bp o fixed_cents", w)
			}
			if pt.PctBP != nil && (*pt.PctBP < 0 || *pt.PctBP > 10000) {
				return errf(CodeInvalidArgument, "%s: pct_bp fuera de 0–10000", w)
			}
		}
		if _, err := splitRule(p.AmountCents, r); err != nil {
			return errf(CodeInvalidArgument, "%s: %v", w, err)
		}
	default:
		return errf(CodeInvalidArgument, "%s: kind %q", w, p.Kind)
	}
	return nil
}

func (in *Input) checkStable(items []StableItem) error {
	for i, it := range items {
		if it.CycleMonths < 1 {
			return errf(CodeInvalidArgument, "stable_items[%d]: cycle_months < 1", i)
		}
		for _, sh := range it.Shares {
			if err := in.money(sh.Person, "stable_items["+itoa(i)+"]"); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkBudgets(rows []BudgetRow) error {
	by := map[string][]BudgetRow{}
	for _, r := range rows {
		k := "\x00"
		if r.Category != nil {
			k = *r.Category
		}
		if r.ValidTo != nil && r.ValidTo.Before(r.ValidFrom) {
			return errf(CodeInvalidArgument, "budgets: renglón %d termina antes de empezar", r.ID)
		}
		by[k] = append(by[k], r)
	}
	for _, rs := range by {
		sort.Slice(rs, func(a, b int) bool { return rs[a].ValidFrom.Before(rs[b].ValidFrom) })
		for k := 1; k < len(rs); k++ {
			if rs[k-1].ValidTo == nil || !rs[k-1].ValidTo.Before(rs[k].ValidFrom) {
				return errf(CodeBudgetOverlap, "budgets: renglones %d y %d se traslapan", rs[k-1].ID, rs[k].ID)
			}
		}
	}
	return nil
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}
