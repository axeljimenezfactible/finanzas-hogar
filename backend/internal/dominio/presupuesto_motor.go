package dominio

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	mf "finanzas-hogar/internal/motor_fin"
)

// Códigos de categoría que el motor S2 pide en params (S2-interfaz §2.1: «category.code, semilla de S1/S6»).
// PROVISIONAL (bloqueo docs/estado/bloqueos/S6-presupuesto-motor.md, pregunta 1): la semilla aún no fija estos códigos.
const (
	GroceriesCategoryCode = "groceries"
	PaperCategoryCode     = "paper"
)

// FinEngine es la entrada del motor S2 (motor_fin.Input) armada desde la base, con los mapas id ↔ código
// que necesitan las respuestas de la API. Se arma una vez por petición y se prepara con At u Only.
type FinEngine struct {
	raw            mf.Input
	Holder         string
	Externals      []string
	CategoryCode   map[int64]string
	CategoryID     map[string]int64
	PersonCode     map[int64]string
	PurchaseCredit map[int64]int64
	RepeatCadence  bool
	WarningPct     int64
	LimitPct       int64
	coHolders      map[string]bool
}

// EngineDate convierte una fecha de dominio a la del motor.
func EngineDate(d Date) mf.Date { return mf.NewDate(d.Year, d.Month, d.Day) }

// FromEngineDate convierte una fecha del motor a la de dominio.
func FromEngineDate(d mf.Date) Date { return NewDate(d.Year, d.Month, d.Day) }

// EngineMonth es el mes del motor que contiene la fecha.
func EngineMonth(d Date) mf.Month { return EngineDate(d).MonthOf() }

// MonthStart es el día 1 del mes (normaliza meses fuera de 1–12).
func MonthStart(y int, m time.Month) Date { return DateOf(time.Date(y, m, 1, 12, 0, 0, 0, MX)) }

// MonthEnd es el último día del mes de la fecha.
func MonthEnd(d Date) Date { return MonthStart(d.Year, d.Month+1).AddDays(-1) }

// Expenses devuelve los gastos registrados tal como los recibe el motor.
func (fe *FinEngine) Expenses() []mf.Expense { return fe.raw.Expenses }

// Purchase busca una compra a crédito por id.
func (fe *FinEngine) Purchase(id int64) (mf.CreditPurchase, bool) {
	for _, c := range fe.raw.CreditPurchases {
		if c.ID == id {
			return c, true
		}
	}
	return mf.CreditPurchase{}, false
}

// At prepara la entrada con «hoy» = today; con upTo solo entran gastos y compras con fecha ≤ upTo.
func (fe *FinEngine) At(today Date, upTo *Date) (mf.Input, error) {
	keepE := func(e mf.Expense) bool { return upTo == nil || !e.ExpenseDate.After(EngineDate(*upTo)) }
	keepP := func(c mf.CreditPurchase) bool { return upTo == nil || !c.PurchaseDate.After(EngineDate(*upTo)) }
	return fe.build(today, keepE, keepP, false)
}

// Only prepara una entrada con un solo gasto (y su compra a crédito), sin ingresos ni partidas planificadas:
// sirve para pedir al motor el devengado de ese gasto en un mes.
func (fe *FinEngine) Only(today Date, expenseID int64) (mf.Input, error) {
	return fe.build(today, func(e mf.Expense) bool { return e.ID == expenseID }, func(mf.CreditPurchase) bool { return false }, true)
}

func (fe *FinEngine) build(today Date, keepE func(mf.Expense) bool, keepP func(mf.CreditPurchase) bool, bare bool) (mf.Input, error) {
	r := &fe.raw
	td := EngineDate(today)
	in := mf.Input{Params: r.Params, Persons: r.Persons, Categories: r.Categories, Budgets: r.Budgets,
		SavingsGoals: r.SavingsGoals, SavingsMovements: r.SavingsMovements, Cadence: r.Cadence, StableItems: r.StableItems}
	if !bare {
		in.Incomes, in.PlannedItems = r.Incomes, r.PlannedItems
	}
	in.Params.Today = &td
	ref := map[int64]bool{}
	for _, e := range r.Expenses {
		if keepE(e) {
			in.Expenses = append(in.Expenses, e)
			if e.CreditPurchaseID != nil {
				ref[*e.CreditPurchaseID] = true
			}
		}
	}
	for _, c := range r.CreditPurchases {
		if ref[c.ID] || (c.ExpenseID == nil && keepP(c)) {
			in.CreditPurchases = append(in.CreditPurchases, c)
		}
	}
	err := in.Prepare()
	return in, err
}

type presScanner interface{ Scan(dest ...any) error }

func presEach(ctx context.Context, q Querier, query string, fn func(presScanner) error, args ...any) error {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := fn(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

func presOptDate(s sql.NullString) (*mf.Date, error) {
	if !s.Valid {
		return nil, nil
	}
	d, err := mf.ParseDate(s.String)
	return &d, err
}

func presOptInt(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return &v
}

func presOptID(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	v := n.Int64
	return &v
}

func presOptCents(n sql.NullInt64) *mf.Cents {
	if !n.Valid {
		return nil
	}
	v := mf.Cents(n.Int64)
	return &v
}

func presOptStr(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	v := s.String
	return &v
}

// money da el código de persona de un campo de dinero: la co-titular va a nombre del titular (C-09, S2-interfaz §1).
func (fe *FinEngine) money(id int64) string {
	c := fe.PersonCode[id]
	if fe.coHolders[c] {
		return fe.Holder
	}
	return c
}

func (fe *FinEngine) catCode(id sql.NullInt64) *string {
	if !id.Valid {
		return nil
	}
	c, ok := fe.CategoryCode[id.Int64]
	if !ok {
		return nil
	}
	return &c
}

// engineCategoryCode da la clave con la que el motor S2 identifica la categoría. Las categorías que crea el
// usuario tienen code NULL (00001_esquema.sql, tabla category): para ellas se usa "cat#<id>", que solo vive en
// memoria y no puede chocar con un code real porque el CHECK del esquema solo admite a-z, 0-9 y guion bajo.
// Las respuestas traducen de vuelta a category_id con CategoryID, así que la clave no se ve en la API
// (bloqueo S6-presupuesto-motor, punto 5).
func engineCategoryCode(id int64, code sql.NullString) string {
	if code.Valid && code.String != "" {
		return code.String
	}
	return "cat#" + strconv.FormatInt(id, 10)
}

// LoadFinEngine lee de la base todo lo que el motor S2 necesita (S2-interfaz §2, columna «Origen»).
// La cadencia de los meses 7 a 12 (Cadence) va vacía: PROVISIONAL, bloqueo S6-presupuesto-motor, pregunta 2.
func LoadFinEngine(ctx context.Context, q Querier) (*FinEngine, error) {
	fe := &FinEngine{CategoryCode: map[int64]string{}, CategoryID: map[string]int64{}, PersonCode: map[int64]string{},
		PurchaseCredit: map[int64]int64{}, WarningPct: 80, LimitPct: 100, coHolders: map[string]bool{}}
	if b, err := SettingBool(ctx, q, "projection.repeat_cadence_enabled"); err == nil {
		fe.RepeatCadence = b
	}
	if n, err := SettingInt(ctx, q, "alerts.budget_warning_pct"); err == nil {
		fe.WarningPct = n
	}
	if n, err := SettingInt(ctx, q, "alerts.budget_limit_pct"); err == nil {
		fe.LimitPct = n
	}
	warn, limit := fe.WarningPct, fe.LimitPct
	in := &fe.raw
	in.Params = mf.Params{GroceriesCategory: GroceriesCategoryCode, PaperCategory: PaperCategoryCode,
		RepeatCadenceEnabled: fe.RepeatCadence, BudgetWarningPct: &warn, BudgetLimitPct: &limit}
	steps := []struct {
		q  string
		fn func(presScanner) error
	}{}
	add := func(query string, fn func(presScanner) error) {
		steps = append(steps, struct {
			q  string
			fn func(presScanner) error
		}{query, fn})
	}

	add(`SELECT id, code, display_name, participation, pays_share_directly FROM person ORDER BY id`, func(s presScanner) error {
		var id int64
		var p mf.Person
		if err := s.Scan(&id, &p.Code, &p.DisplayName, &p.Participation, &p.PaysShareDirectly); err != nil {
			return err
		}
		fe.PersonCode[id] = p.Code
		switch p.Participation {
		case "holder":
			fe.Holder = p.Code
		case "external":
			fe.Externals = append(fe.Externals, p.Code)
		case "co_holder":
			fe.coHolders[p.Code] = true
		}
		in.Persons = append(in.Persons, p)
		return nil
	})
	catRule := map[int64]sql.NullInt64{}
	add(`SELECT id, code, name, kind, is_active, default_cost_share_rule_id FROM category ORDER BY id`, func(s presScanner) error {
		var id int64
		var code, kind sql.NullString
		var active bool
		var rule sql.NullInt64
		var c mf.Category
		if err := s.Scan(&id, &code, &c.Name, &kind, &active, &rule); err != nil {
			return err
		}
		c.Code = engineCategoryCode(id, code)
		c.Kind, c.IsActive = presOptStr(kind), &active
		fe.CategoryCode[id], fe.CategoryID[c.Code] = c.Code, id
		catRule[id] = rule
		in.Categories = append(in.Categories, c)
		return nil
	})
	rules := map[int64]*mf.ShareRule{}
	var ruleOrder []int64
	add(`SELECT id, remainder_person_id FROM cost_share_rule ORDER BY id`, func(s presScanner) error {
		var id int64
		var rest sql.NullInt64
		if err := s.Scan(&id, &rest); err != nil {
			return err
		}
		rules[id] = &mf.ShareRule{RestPerson: fe.money(rest.Int64)}
		ruleOrder = append(ruleOrder, id)
		return nil
	})
	add(`SELECT rule_id, person_id, pct_bp, fixed_cents FROM cost_share_rule_part ORDER BY rule_id, id`, func(s presScanner) error {
		var rule, person int64
		var bp, fixed sql.NullInt64
		if err := s.Scan(&rule, &person, &bp, &fixed); err != nil {
			return err
		}
		if r := rules[rule]; r != nil {
			r.Parts = append(r.Parts, mf.RulePart{Person: fe.money(person), PctBP: presOptID(bp), FixedCents: presOptCents(fixed)})
		}
		return nil
	})
	var defRule sql.NullInt64
	add(`SELECT id FROM cost_share_rule WHERE is_default = 1 ORDER BY id LIMIT 1`, func(s presScanner) error { return s.Scan(&defRule) })
	shares := map[int64][]mf.Share{}
	add(`SELECT expense_id, person_id, amount_cents FROM expense_share ORDER BY expense_id, id`, func(s presScanner) error {
		var e, p int64
		var a mf.Cents
		if err := s.Scan(&e, &p, &a); err != nil {
			return err
		}
		shares[e] = append(shares[e], mf.Share{Person: fe.money(p), AmountCents: a})
		return nil
	})
	add(`SELECT e.id, e.kind, e.category_id, e.amount_cents, e.expense_date, e.accrual_months, e.payment_method, e.paid_by_person_id, e.planned_item_id,
		(SELECT cp.id FROM credit_purchase cp WHERE cp.expense_id = e.id ORDER BY cp.id LIMIT 1) FROM expense e ORDER BY e.id`, func(s presScanner) error {
		var x mf.Expense
		var cat, acc, paidBy, planned, purchase sql.NullInt64
		var date string
		if err := s.Scan(&x.ID, &x.Kind, &cat, &x.AmountCents, &date, &acc, &x.PaymentMethod, &paidBy, &planned, &purchase); err != nil {
			return err
		}
		d, err := mf.ParseDate(date)
		if err != nil {
			return err
		}
		x.ExpenseDate, x.AccrualMonths, x.PlannedItemID, x.CreditPurchaseID, x.Shares = d, presOptInt(acc), presOptID(planned), presOptID(purchase), shares[x.ID]
		if c := fe.catCode(cat); c != nil {
			x.Category = *c
		}
		if paidBy.Valid {
			p := fe.money(paidBy.Int64)
			x.PaidBy = &p
		}
		in.Expenses = append(in.Expenses, x)
		return nil
	})
	pshares := map[int64][]mf.PurchaseShare{}
	add(`SELECT credit_purchase_id, person_id, principal_cents, interest_cents, total_cents FROM credit_purchase_share ORDER BY credit_purchase_id, id`, func(s presScanner) error {
		var c, p int64
		var sh mf.PurchaseShare
		if err := s.Scan(&c, &p, &sh.PrincipalCents, &sh.InterestCents, &sh.TotalCents); err != nil {
			return err
		}
		sh.Person = fe.money(p)
		pshares[c] = append(pshares[c], sh)
		return nil
	})
	paid := map[int64][]mf.Paid{}
	add(`SELECT scheduled_payment_id, payer_person_id, amount_cents, paid_date FROM credit_payment WHERE scheduled_payment_id IS NOT NULL ORDER BY paid_date, id`, func(s presScanner) error {
		var sp int64
		var payer sql.NullInt64
		var pd mf.Paid
		var date string
		if err := s.Scan(&sp, &payer, &pd.AmountCents, &date); err != nil {
			return err
		}
		d, err := mf.ParseDate(date)
		if err != nil {
			return err
		}
		pd.PaidDate, pd.Person = d, fe.Holder
		if payer.Valid {
			pd.Person = fe.money(payer.Int64)
		}
		paid[sp] = append(paid[sp], pd)
		return nil
	})
	spShares := map[int64][]mf.Share{}
	add(`SELECT scheduled_payment_id, person_id, amount_cents FROM scheduled_payment_share ORDER BY scheduled_payment_id, id`, func(s presScanner) error {
		var sp, p int64
		var a mf.Cents
		if err := s.Scan(&sp, &p, &a); err != nil {
			return err
		}
		spShares[sp] = append(spShares[sp], mf.Share{Person: fe.money(p), AmountCents: a})
		return nil
	})
	pays := map[int64][]mf.ScheduledPayment{}
	add(`SELECT id, credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents FROM scheduled_payment ORDER BY credit_purchase_id, seq`, func(s presScanner) error {
		var id, cp int64
		var p mf.ScheduledPayment
		var date string
		if err := s.Scan(&id, &cp, &p.Seq, &date, &p.AmountCents, &p.PrincipalCents, &p.InterestCents); err != nil {
			return err
		}
		d, err := mf.ParseDate(date)
		if err != nil {
			return err
		}
		p.DueDate, p.Shares, p.Paid = d, spShares[id], paid[id]
		pays[cp] = append(pays[cp], p)
		return nil
	})
	add(`SELECT cp.id, cp.credit_id, cp.expense_id, e.category_id, e.accrual_months, cp.purchase_date, cp.principal_cents, cp.interest_cents, cp.total_cents, cp.status
		FROM credit_purchase cp LEFT JOIN expense e ON e.id = cp.expense_id ORDER BY cp.id`, func(s presScanner) error {
		var c mf.CreditPurchase
		var exp, cat, acc sql.NullInt64
		var date string
		if err := s.Scan(&c.ID, &c.CreditID, &exp, &cat, &acc, &date, &c.PrincipalCents, &c.InterestCents, &c.TotalCents, &c.Status); err != nil {
			return err
		}
		d, err := mf.ParseDate(date)
		if err != nil {
			return err
		}
		c.PurchaseDate, c.ExpenseID, c.Category, c.AccrualMonths = d, presOptID(exp), fe.catCode(cat), presOptInt(acc)
		c.Shares, c.Payments = pshares[c.ID], pays[c.ID]
		fe.PurchaseCredit[c.ID] = c.CreditID
		in.CreditPurchases = append(in.CreditPurchases, c)
		return nil
	})
	add(`SELECT id, person_id, kind, amount_cents, income_date, planned_item_id FROM income ORDER BY income_date, id`, func(s presScanner) error {
		var x mf.Income
		var p int64
		var planned sql.NullInt64
		var date string
		if err := s.Scan(&x.ID, &p, &x.Kind, &x.AmountCents, &date, &planned); err != nil {
			return err
		}
		d, err := mf.ParseDate(date)
		if err != nil {
			return err
		}
		x.IncomeDate, x.Person, x.PlannedItemID = d, fe.PersonCode[p], presOptID(planned)
		in.Incomes = append(in.Incomes, x)
		return nil
	})
	add(`SELECT id, kind, description, amount_cents, second_half_amount_cents, frequency, interval_count, start_date, end_date, category_id,
		person_id, cost_share_rule_id, payment_method, credit_id, accrual_months, is_active FROM planned_item ORDER BY id`, func(s presScanner) error {
		var p mf.PlannedItem
		var second, interval, cat, person, rule, credit, acc sql.NullInt64
		var start string
		var end, method sql.NullString
		var active bool
		if err := s.Scan(&p.ID, &p.Kind, &p.Description, &p.AmountCents, &second, &p.Frequency, &interval, &start, &end, &cat,
			&person, &rule, &method, &credit, &acc, &active); err != nil {
			return err
		}
		d, err := mf.ParseDate(start)
		if err != nil {
			return err
		}
		if p.EndDate, err = presOptDate(end); err != nil {
			return err
		}
		p.StartDate, p.SecondHalfAmountCents, p.IntervalCount, p.AccrualMonths = d, presOptCents(second), presOptInt(interval), presOptInt(acc)
		p.Category, p.PaymentMethod, p.CreditID, p.IsActive = fe.catCode(cat), presOptStr(method), presOptID(credit), &active
		if person.Valid {
			c := fe.PersonCode[person.Int64]
			p.Person = &c
		}
		if p.Kind != "income" {
			id := rule
			if !id.Valid && cat.Valid {
				id = catRule[cat.Int64]
			}
			if !id.Valid {
				id = defRule
			}
			if id.Valid && rules[id.Int64] != nil {
				p.ShareRule = rules[id.Int64]
			}
		}
		in.PlannedItems = append(in.PlannedItems, p)
		return nil
	})
	add(`SELECT id, category_id, amount_cents, valid_from_date, valid_to_date FROM budget ORDER BY id`, func(s presScanner) error {
		var b mf.BudgetRow
		var cat sql.NullInt64
		var from string
		var to sql.NullString
		if err := s.Scan(&b.ID, &cat, &b.AmountCents, &from, &to); err != nil {
			return err
		}
		d, err := mf.ParseDate(from)
		if err != nil {
			return err
		}
		if b.ValidTo, err = presOptDate(to); err != nil {
			return err
		}
		b.ValidFrom, b.Category = d, fe.catCode(cat)
		in.Budgets = append(in.Budgets, b)
		return nil
	})
	add(`SELECT id, name, kind, target_cents, target_date, planned_monthly_cents, start_date, is_active FROM savings_goal ORDER BY id`, func(s presScanner) error {
		var g mf.SavingsGoal
		var target, planned sql.NullInt64
		var tdate sql.NullString
		var start string
		var active bool
		if err := s.Scan(&g.ID, &g.Name, &g.Kind, &target, &tdate, &planned, &start, &active); err != nil {
			return err
		}
		d, err := mf.ParseDate(start)
		if err != nil {
			return err
		}
		if g.TargetDate, err = presOptDate(tdate); err != nil {
			return err
		}
		g.StartDate, g.TargetCents, g.PlannedMonthlyCents, g.IsActive = d, presOptCents(target), presOptCents(planned), &active
		in.SavingsGoals = append(in.SavingsGoals, g)
		return nil
	})
	add(`SELECT savings_goal_id, movement_date, kind, amount_cents FROM savings_movement ORDER BY movement_date, id`, func(s presScanner) error {
		var m mf.SavingsMovement
		var date string
		if err := s.Scan(&m.GoalID, &date, &m.Kind, &m.AmountCents); err != nil {
			return err
		}
		d, err := mf.ParseDate(date)
		m.MovementDate = d
		in.SavingsMovements = append(in.SavingsMovements, m)
		return err
	})
	for _, st := range steps {
		if err := presEach(ctx, q, st.q, st.fn); err != nil {
			return nil, err
		}
	}
	in.Params.Holder = fe.Holder
	return fe, nil
}
