// Package gastos implementa las operaciones de gastos y de inicio de C1 del agente api-gastos-inicio
// (CONVENCIONES-backend.md §2: paquete gastos, función Register).
package gastos

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

// Register registra las operaciones implementadas; las demás de esta área siguen en 501 (ver docs/estado/wf3/api-gastos-inicio.md).
func Register(reg *api.Registry) {
	reg.Handle("listExpenses", listExpenses(reg))
	reg.Handle("createExpense", createExpense(reg))
	reg.Handle("getExpenseTotals", getExpenseTotals(reg))
	reg.Handle("getExpenseSummary", getExpenseSummary(reg))
	reg.Handle("getExpense", getExpense(reg))
	reg.Handle("updateExpense", updateExpense(reg))
	reg.Handle("deleteExpense", deleteExpense(reg))
	reg.Handle("listExpenseItems", listExpenseItems(reg))
	reg.Handle("duplicateExpense", duplicateExpense(reg))
	reg.Handle("markExpensePaid", markExpensePaid(reg))
	reg.Handle("listTicketScans", listTicketScans(reg))
	reg.Handle("getTicketScan", getTicketScan(reg))
	reg.Handle("saveTicketScanDraft", saveTicketScanDraft(reg))
	reg.Handle("getExpenseMonthlyTrend", getExpenseMonthlyTrend(reg))
	reg.Handle("getDashboardMonthVsPlan", getDashboardMonthVsPlan(reg))
	reg.Handle("getDashboardFortnightFlow", getDashboardFortnightFlow(reg))
	reg.Handle("getDashboardSavings", getDashboardSavings(reg))
	reg.Handle("getPersonParticipation", getPersonParticipation(reg))
	reg.Handle("search", search(reg))
	reg.Handle("getDashboardNextPayment", getDashboardNextPayment(reg))
	reg.Handle("getDashboardWeekMeals", getDashboardWeekMeals(reg))
	reg.Handle("getDashboardBestBasket", getDashboardBestBasket(reg))
	reg.Handle("listCalendarEvents", listCalendarEvents(reg))
	reg.Handle("queryView", queryView(reg))
	reg.Handle("createTicketScan", createTicketScan(reg))
	reg.Handle("confirmTicketScan", confirmTicketScan(reg))
	registerTicketJob(reg)
}

var (
	kinds          = []string{"fixed", "variable", "upcoming"}
	paymentMethods = []string{"cash", "debit", "transfer", "credit"}
	families       = []string{"household", "external", "joint"}
	expenseSort    = map[string]string{"id": "e.id", "kind": "e.kind", "description": "e.description", "category_id": "e.category_id",
		"amount_cents": "e.amount_cents", "expense_date": "e.expense_date", "accrual_months": "e.accrual_months", "payment_method": "e.payment_method",
		"paid_by_person_id": "e.paid_by_person_id", "person_id": "e.person_id", "cost_share_rule_id": "e.cost_share_rule_id", "store_id": "e.store_id",
		"planned_item_id": "e.planned_item_id", "notes": "e.notes", "created_at": "e.created_at", "updated_at": "e.updated_at", "captured_at": "e.captured_at"}
)

// queryEnum lee un parámetro de lista (repetido o separado por coma) y rechaza valores fuera del enumerado con 400.
func queryEnum(r *http.Request, name string, allowed []string) ([]string, error) {
	var out []string
	for _, raw := range r.URL.Query()[name] {
		for _, s := range strings.Split(raw, ",") {
			if !slices.Contains(allowed, s) {
				return nil, api.BadRequest("El filtro " + name + " solo admite: " + strings.Join(allowed, ", ") + ".")
			}
			out = append(out, s)
		}
	}
	return out, nil
}

// parseFilter lee los filtros de gastos; full = false solo lee los globales (persona, familia, categoría, tienda, periodo, crédito, q).
func parseFilter(reg *api.Registry, r *http.Request, full bool) (dominio.ExpenseFilter, error) {
	var f dominio.ExpenseFilter
	var err error
	for name, dst := range map[string]*[]int64{"person_id": &f.PersonIDs, "category_id": &f.CategoryIDs, "store_id": &f.StoreIDs, "credit_id": &f.CreditIDs} {
		if *dst, err = api.QueryIDs(r, name); err != nil {
			return f, err
		}
	}
	qs := r.URL.Query()
	if fam := qs.Get("family"); fam != "" {
		if !slices.Contains(families, fam) {
			return f, api.BadRequest("El filtro family solo admite household, external o joint.")
		}
		f.Family = fam
	}
	if f.Q = qs.Get("q"); len([]rune(f.Q)) > 200 {
		return f, api.BadRequest("El texto de búsqueda admite hasta 200 caracteres.")
	}
	if f.Period, err = api.QueryPeriod(r); err != nil {
		return f, err
	}
	if f.Period == nil {
		p := dominio.MonthPeriod(dominio.Today(reg.Deps().Clock)) // periodo por defecto: mes en curso
		f.Period = &p
	}
	if !full {
		return f, nil
	}
	if f.Kinds, err = queryEnum(r, "kind", kinds); err != nil {
		return f, err
	}
	if f.PaymentMethods, err = queryEnum(r, "payment_method", paymentMethods); err != nil {
		return f, err
	}
	switch qs.Get("has_attachment") {
	case "":
	case "true":
		f.HasAttachment = api.Ptr(true)
	case "false":
		f.HasAttachment = api.Ptr(false)
	default:
		return f, api.BadRequest("El filtro has_attachment solo admite true o false.")
	}
	return f, nil
}

func nullInt(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return api.Ptr(n.Int64)
}

func nullStr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	return api.Ptr(n.String)
}

// toExpense arma la respuesta Expense de C1 (columnas NULL = campo ausente).
func toExpense(x dominio.ExpenseRow) api.Expense {
	e := api.Expense{ID: x.ID, Kind: x.Kind, Description: nullStr(x.Description), CategoryID: x.CategoryID, AmountCents: x.AmountCents,
		ExpenseDate: x.ExpenseDate, AccrualMonths: x.AccrualMonths, PaymentMethod: x.PaymentMethod, PaidByPersonID: nullInt(x.PaidByPersonID),
		PersonID: nullInt(x.PersonID), CostShareRuleID: nullInt(x.CostShareRuleID), StoreID: nullInt(x.StoreID), PlannedItemID: nullInt(x.PlannedItemID),
		Notes: nullStr(x.Notes), CreatedAt: api.APITime(x.CreatedAt), UpdatedAt: api.APITime(x.UpdatedAt),
		HouseholdShareCents: api.Ptr(x.HouseholdCents), ExternalShareCents: api.Ptr(x.ExternalCents), HasAttachment: api.Ptr(x.HasAttachment),
		CreditPurchaseID: nullInt(x.CreditPurchaseID)}
	if x.CapturedAt.Valid {
		e.CapturedAt = api.Ptr(api.APITime(x.CapturedAt.String))
	}
	for _, s := range x.Shares {
		e.Shares = append(e.Shares, api.ShareLine{PersonID: s.PersonID, AmountCents: s.AmountCents})
	}
	if x.PaymentMethod != "credit" { // RN-02: de contado, débito o transferencia el efectivo sale en expense_date
		e.CashFlows = []api.ExpenseCashFlowsItem{{Date: x.ExpenseDate, AmountCents: x.AmountCents}}
	}
	return e
}

func writeExpense(reg *api.Registry, w http.ResponseWriter, r *http.Request, status int, id int64) error {
	x, err := dominio.GetExpense(r.Context(), reg.Deps().DB, id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, status, toExpense(x))
}

func listExpenses(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		f, err := parseFilter(reg, r, true)
		if err != nil {
			return err
		}
		allowed := make([]string, 0, len(expenseSort))
		for k := range expenseSort {
			allowed = append(allowed, k)
		}
		p, err := api.ParsePage(r, allowed, "-expense_date,-created_at")
		if err != nil {
			return err
		}
		xs, err := dominio.ListExpenses(r.Context(), reg.Deps().DB, f, p.OrderBy(expenseSort), p.LimitOffset())
		if err != nil {
			return err
		}
		items := make([]api.Expense, 0, len(xs))
		for _, x := range xs {
			items = append(items, toExpense(x))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func getExpenseTotals(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		f, err := parseFilter(reg, r, false)
		if err != nil {
			return err
		}
		t, err := dominio.SumExpenses(r.Context(), reg.Deps().DB, f)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, api.ExpenseTotals{ExpenseCount: t.Count, AmountCents: t.AmountCents,
			HouseholdShareCents: t.HouseholdCents, ExternalShareCents: t.ExternalCents})
	}
}

func getExpense(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		return writeExpense(reg, w, r, http.StatusOK, id)
	}
}

// partsIn convierte el reparto propio de C1 al del dominio.
func partsIn(ps []api.CostShareRulePartInput) []dominio.ExpenseSharePart {
	var out []dominio.ExpenseSharePart
	for _, p := range ps {
		out = append(out, dominio.ExpenseSharePart{PersonID: p.PersonID, PctBP: p.PctBp, FixedCents: p.FixedCents})
	}
	return out
}

func itemsIn(its []api.ExpenseItemInput) []dominio.ExpenseItemIn {
	out := make([]dominio.ExpenseItemIn, 0, len(its))
	for _, it := range its {
		out = append(out, dominio.ExpenseItemIn{LineNo: it.LineNo, Description: strings.TrimSpace(it.Description), PresentationID: it.PresentationID,
			QuantityMilli: it.QuantityMilli, UnitPriceCents: it.UnitPriceCents, LineTotalCents: it.LineTotalCents, CostShareRuleID: it.CostShareRuleID,
			Parts: partsIn(it.CostShareRuleParts)})
	}
	return out
}

// validate revisa las reglas de negocio de un gasto ya mezclado (422 con todos los campos).
func validate(r *http.Request, q dominio.Querier, in dominio.ExpenseIn) error {
	return validateExpense(r, q, in, false)
}

// validateExpense es validate con la forma de pago crédito permitida o no: solo confirmTicketScan la permite,
// porque crea la compra a crédito en la misma transacción (N-074, N-218).
func validateExpense(r *http.Request, q dominio.Querier, in dominio.ExpenseIn, allowCredit bool) error {
	var v api.Validator
	ctx := r.Context()
	v.Check(slices.Contains(kinds, in.Kind), "kind", "invalid_format", "El tipo debe ser fijo, variable o venidero.")
	v.Check(slices.Contains(paymentMethods, in.PaymentMethod), "payment_method", "invalid_format", "La forma de pago debe ser efectivo, débito, transferencia o crédito.")
	v.Check(allowCredit || in.PaymentMethod != "credit", "payment_method", "not_allowed", "Las compras a crédito se registran en Créditos, con su plan de pagos.")
	v.PositiveCents(in.AmountCents, "amount_cents")
	v.Check(in.AccrualMonths >= 1 && in.AccrualMonths <= 24, "accrual_months", "out_of_range", "Los meses de devengo van de 1 a 24.")
	v.Check(in.Description == nil || strings.TrimSpace(*in.Description) != "", "description", "invalid_format", "La descripción no puede quedar en blanco; bórrala si no la necesitas.")
	v.Check(len(in.Parts) == 0 || in.CostShareRuleID == nil, "cost_share_rule_parts", "conflict", "Usa una regla de reparto o un reparto propio, no los dos.")
	refs := []struct {
		id         *int64
		table, fld string
	}{{&in.CategoryID, "category", "category_id"}, {in.PaidByPersonID, "person", "paid_by_person_id"}, {in.PersonID, "person", "person_id"},
		{in.CostShareRuleID, "cost_share_rule", "cost_share_rule_id"}, {in.StoreID, "store", "store_id"}, {in.PlannedItemID, "planned_item", "planned_item_id"}}
	for _, ref := range refs {
		if ref.id == nil {
			continue
		}
		ok, err := dominio.ExpenseRefExists(ctx, q, ref.table, *ref.id)
		if err != nil {
			return err
		}
		v.Check(ok, ref.fld, "not_found", "No existe el registro elegido en "+ref.fld+".")
	}
	if in.ItemsSet && len(in.Items) > 0 {
		sum, seen := int64(0), map[int64]bool{}
		for _, it := range in.Items {
			v.Check(it.LineNo >= 1 && !seen[it.LineNo], "items", "invalid_format", "Cada renglón lleva un número de renglón distinto, desde 1.")
			v.Check(it.Description != "", "items", "required", "Cada renglón lleva su texto.")
			v.Check(it.QuantityMilli > 0, "items", "out_of_range", "La cantidad de cada renglón debe ser mayor que 0.")
			v.Check(it.LineTotalCents >= 0 && (it.UnitPriceCents == nil || *it.UnitPriceCents >= 0), "items", "out_of_range", "Los importes no pueden ser negativos.")
			v.Check(len(it.Parts) == 0 || it.CostShareRuleID == nil, "items", "conflict", "Cada renglón usa una regla de reparto o un reparto mixto, no los dos.")
			seen[it.LineNo] = true
			sum += it.LineTotalCents
		}
		v.Check(sum == in.AmountCents, "items", "out_of_range", "La suma de los renglones debe ser igual al importe del gasto.")
	}
	return v.Err()
}

func audit(r *http.Request, tx *sql.Tx, action string, id int64, before json.RawMessage) error {
	var after json.RawMessage
	if action != "delete" {
		var err error
		if after, err = dominio.RowJSON(r.Context(), tx, "expense", id); err != nil {
			return err
		}
	}
	return dominio.Audit(r.Context(), tx, api.Actor(r), action, "expense", id, before, after, api.RequestID(r))
}

// newExpenseIn arma el alta de un gasto (createExpense y confirmTicketScan): variable o, con fecha futura, venidero;
// un mes de devengo si no se dice otra cosa; captured_at en UTC.
func newExpenseIn(reg *api.Registry, b api.ExpenseCreate, d dominio.Date) (dominio.ExpenseIn, error) {
	in := dominio.ExpenseIn{Kind: "variable", Description: b.Description, CategoryID: b.CategoryID, AmountCents: b.AmountCents, ExpenseDate: d.String(),
		AccrualMonths: 1, PaymentMethod: b.PaymentMethod, PaidByPersonID: b.PaidByPersonID, PersonID: b.PersonID, CostShareRuleID: b.CostShareRuleID,
		StoreID: b.StoreID, PlannedItemID: b.PlannedItemID, Notes: b.Notes, Parts: partsIn(b.CostShareRuleParts), Items: itemsIn(b.Items), ItemsSet: true}
	if d.After(dominio.Today(reg.Deps().Clock)) {
		in.Kind = "upcoming"
	}
	if b.Kind != nil {
		in.Kind = *b.Kind
	}
	if b.AccrualMonths != nil {
		in.AccrualMonths = *b.AccrualMonths
	}
	if b.CapturedAt != nil {
		t, err := dominio.ParseDateTime(*b.CapturedAt)
		if err != nil {
			return in, api.BadRequest("captured_at debe ser un instante ISO 8601 con zona.")
		}
		in.CapturedAt = api.Ptr(dominio.FormatDBTime(t))
	}
	return in, nil
}

func createExpense(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var b api.ExpenseCreate
		if err := reg.DecodeJSON(r, &b); err != nil {
			return err
		}
		var v api.Validator
		d, err := v.Date(b.ExpenseDate, "expense_date")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		key := api.IdempotencyKey(r)
		if id, ok, err := dominio.ExpenseIDByIdempotencyKey(ctx, db, key); err != nil {
			return err
		} else if ok { // K9: reintento de una captura ya registrada, aunque hayan pasado las 24 h
			api.MarkReplayed(w)
			return writeExpense(reg, w, r, http.StatusOK, id)
		}
		in, err := newExpenseIn(reg, b, d)
		if err != nil {
			return err
		}
		if key != "" {
			in.IdempotencyKey = &key
		}
		if err := validate(r, db, in); err != nil {
			return err
		}
		var id int64
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			var err error
			if id, err = dominio.SaveExpense(ctx, tx, 0, in, true); err != nil {
				return err
			}
			return audit(r, tx, "insert", id, nil)
		})
		if err != nil && dominio.IsUnique(err) && key != "" { // la misma clave llegó dos veces a la vez
			if id, ok, _ := dominio.ExpenseIDByIdempotencyKey(ctx, db, key); ok {
				api.MarkReplayed(w)
				return writeExpense(reg, w, r, http.StatusOK, id)
			}
		}
		if err != nil {
			return err
		}
		return writeExpense(reg, w, r, http.StatusCreated, id)
	}
}

func updateExpense(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		x, err := dominio.GetExpense(ctx, db, id)
		if err != nil {
			return err
		}
		var p api.ExpenseUpdate
		if err := reg.DecodeJSON(r, &p); err != nil {
			return err
		}
		in, err := dominio.ExpenseInOf(ctx, db, x)
		if err != nil {
			return err
		}
		recompute := false
		if p.Kind != nil {
			in.Kind = *p.Kind
		}
		if p.CategoryID != nil {
			in.CategoryID, recompute = *p.CategoryID, true
		}
		if p.AmountCents != nil {
			in.AmountCents, recompute = *p.AmountCents, true
		}
		if p.ExpenseDate != nil {
			var v api.Validator
			d, err := v.Date(*p.ExpenseDate, "expense_date")
			if err != nil {
				return err
			}
			in.ExpenseDate = d.String()
		}
		if p.AccrualMonths != nil {
			in.AccrualMonths = *p.AccrualMonths
		}
		if p.PaymentMethod != nil {
			in.PaymentMethod = *p.PaymentMethod
		}
		setStr := func(n api.Nullable[string], dst **string) {
			if n.Set {
				*dst = nil
				if n.Valid {
					*dst = api.Ptr(n.Value)
				}
			}
		}
		setID := func(n api.Nullable[int64], dst **int64) {
			if n.Set {
				*dst = nil
				if n.Valid {
					*dst = api.Ptr(n.Value)
				}
			}
		}
		setStr(p.Description, &in.Description)
		setStr(p.Notes, &in.Notes)
		setID(p.PaidByPersonID, &in.PaidByPersonID)
		setID(p.PersonID, &in.PersonID)
		setID(p.StoreID, &in.StoreID)
		setID(p.PlannedItemID, &in.PlannedItemID)
		if p.CostShareRuleID.Set {
			setID(p.CostShareRuleID, &in.CostShareRuleID)
			recompute = true
		}
		if p.CostShareRuleParts != nil {
			in.Parts, recompute = partsIn(p.CostShareRuleParts), true
			if !p.CostShareRuleID.Set {
				in.CostShareRuleID = nil // el reparto propio sustituye a la regla guardada
			}
		}
		if p.Items != nil {
			in.Items, in.ItemsSet, recompute = itemsIn(p.Items), true, true
		}
		if x.CreditPurchaseID.Valid && (in.AmountCents != x.AmountCents || in.PaymentMethod != x.PaymentMethod || in.ExpenseDate != x.ExpenseDate) {
			return dominio.Conflict("has_credit_purchase", "Este gasto viene de una compra a crédito: su importe, fecha y forma de pago se cambian en la compra.")
		}
		if x.PaymentMethod == "credit" && in.PaymentMethod == "credit" {
			in.PaymentMethod = "cash" // solo para pasar la validación de alta: el método guardado no cambia
		}
		if err := validate(r, db, in); err != nil {
			return err
		}
		if x.PaymentMethod == "credit" {
			in.PaymentMethod = "credit"
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			before, err := dominio.RowJSON(ctx, tx, "expense", id)
			if err != nil {
				return err
			}
			if _, err := dominio.SaveExpense(ctx, tx, id, in, recompute); err != nil {
				return err
			}
			return audit(r, tx, "update", id, before)
		})
		if err != nil {
			return err
		}
		return writeExpense(reg, w, r, http.StatusOK, id)
	}
}

func deleteExpense(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			x, err := dominio.GetExpense(ctx, tx, id)
			if err != nil {
				return err
			}
			before, err := dominio.RowJSON(ctx, tx, "expense", id)
			if err != nil {
				return err
			}
			if err := dominio.DeleteExpense(ctx, tx, x); err != nil {
				return err
			}
			return audit(r, tx, "delete", id, before)
		})
		if err != nil {
			return err
		}
		return api.NoContent(w)
	}
}

func duplicateExpense(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		var nid int64
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			x, err := dominio.GetExpense(ctx, tx, id)
			if err != nil {
				return err
			}
			if x.PaymentMethod == "credit" {
				return dominio.Conflict("has_credit_purchase", "Una compra a crédito no se duplica aquí: regístrala en Créditos con su plan de pagos.")
			}
			in, err := dominio.ExpenseInOf(ctx, tx, x)
			if err != nil {
				return err
			}
			in.ExpenseDate, in.ItemsSet = dominio.Today(reg.Deps().Clock).String(), true
			if nid, err = dominio.SaveExpense(ctx, tx, 0, in, false); err != nil {
				return err
			}
			if err := dominio.CopyExpenseShares(ctx, tx, id, nid); err != nil { // copia: mismas partes congeladas
				return err
			}
			return audit(r, tx, "insert", nid, nil)
		})
		if err != nil {
			return err
		}
		return writeExpense(reg, w, r, http.StatusCreated, nid)
	}
}

func markExpensePaid(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		var b api.ExpenseMarkPaidRequest
		if err := reg.DecodeJSON(r, &b); err != nil {
			return err
		}
		var v api.Validator
		d, err := v.Date(b.ExpenseDate, "expense_date")
		if err != nil {
			return err
		}
		v.PositiveCents(b.AmountCents, "amount_cents")
		if err := v.Err(); err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			x, err := dominio.GetExpense(ctx, tx, id)
			if err != nil {
				return err
			}
			if x.Kind != "upcoming" {
				return dominio.Conflict("conflict", "Este gasto no es venidero: ya está registrado como pagado.")
			}
			in, err := dominio.ExpenseInOf(ctx, tx, x)
			if err != nil {
				return err
			}
			pk, err := dominio.ExpensePlannedItemKind(ctx, tx, in.PlannedItemID)
			if err != nil {
				return err
			}
			in.Kind, in.AmountCents, in.ExpenseDate = "variable", b.AmountCents, d.String()
			if pk == "fixed_expense" { // N-065: pasa a fijo o variable según su partida
				in.Kind = "fixed"
			}
			before, err := dominio.RowJSON(ctx, tx, "expense", id)
			if err != nil {
				return err
			}
			if _, err := dominio.SaveExpense(ctx, tx, id, in, true); err != nil {
				return err
			}
			return audit(r, tx, "update", id, before)
		})
		if err != nil {
			return err
		}
		return writeExpense(reg, w, r, http.StatusOK, id)
	}
}

func listExpenseItems(reg *api.Registry) api.HandlerFunc {
	cols := map[string]string{"id": "i.id", "line_no": "i.line_no", "description": "i.description", "quantity_milli": "i.quantity_milli",
		"unit_price_cents": "i.unit_price_cents", "line_total_cents": "i.line_total_cents", "created_at": "i.created_at"}
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		x, err := dominio.GetExpense(ctx, db, id)
		if err != nil {
			return err
		}
		allowed := make([]string, 0, len(cols))
		for k := range cols {
			allowed = append(allowed, k)
		}
		p, err := api.ParsePage(r, allowed, "line_no")
		if err != nil {
			return err
		}
		rows, err := dominio.ListExpenseItems(ctx, db, id, p.OrderBy(cols), p.LimitOffset())
		if err != nil {
			return err
		}
		items := make([]api.ExpenseItem, 0, len(rows)+1)
		for _, it := range rows {
			items = append(items, api.ExpenseItem{ID: it.ID, ExpenseID: it.ExpenseID, LineNo: it.LineNo, Description: it.Description,
				PresentationID: nullInt(it.PresentationID), QuantityMilli: it.QuantityMilli, UnitPriceCents: nullInt(it.UnitPriceCents),
				LineTotalCents: it.LineTotalCents, CostShareRuleID: nullInt(it.CostShare), CreatedAt: api.APITime(it.CreatedAt), UpdatedAt: api.APITime(it.UpdatedAt)})
		}
		sum, maxLine, n, err := dominio.ExpenseItemsSum(ctx, db, id)
		if err != nil {
			return err
		}
		if n > 0 && sum < x.AmountCents && len(rows) <= p.Limit { // última página: la diferencia viene como «Ajuste del ticket»
			items = append(items, api.ExpenseItem{ID: id, ExpenseID: id, LineNo: maxLine + 1, Description: "Ajuste del ticket", QuantityMilli: 1000,
				LineTotalCents: x.AmountCents - sum, CreatedAt: api.APITime(x.CreatedAt), UpdatedAt: api.APITime(x.UpdatedAt)})
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}
