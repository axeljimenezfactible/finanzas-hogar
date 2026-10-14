package gastos

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
	mf "finanzas-hogar/internal/motor_fin"
	"finanzas-hogar/internal/precios"
	"finanzas-hogar/internal/recetas"
)

// opt es un valor de un enumerado con su etiqueta en español (columnas del tablero en su orden).
type opt struct{ k, l string }

// groupDef es una propiedad de agrupación permitida (C1 §8). Con enum, las columnas salen aunque estén vacías;
// con derived, la clave la pone un motor en record[campo]; con multi, el registro va en cada columna.
type groupDef struct {
	key, label string
	enum       []opt
	derived    bool
	multi      bool
}

// viewDef es la definición de una entidad de queryView (C1 §8). Las expresiones SQL usan el alias t.
type viewDef struct {
	modes                []string
	groups               map[string]groupDef
	dates                map[string]string // date_field → columna
	defDate              string            // fecha del filtro period y del calendario por defecto
	title, amount, image string
	extra                string
	sorts                []string
	defSort              string
	person, credit       string // expresión del filtro global (vacía = el filtro no aplica a la entidad)
	store                string
	status               string // propiedad de record con el estado propio (status)
	statusEnum           []opt
	derivedAmount        bool // el importe principal lo da un motor (ahorrado, deuda viva)
	derive               func(reg *api.Registry, r *http.Request, recs []dominio.ViewRecord) error
}

var (
	allModes      = []string{"table", "board", "calendar", "gallery"}
	lbExpenseKind = []opt{{"fixed", "Fijo"}, {"variable", "Variable"}, {"upcoming", "Venidero"}}
	lbPayMethod   = []opt{{"cash", "Efectivo"}, {"debit", "Débito"}, {"transfer", "Transferencia"}, {"credit", "Crédito"}}
	lbGoalStatus  = []opt{{"in_progress", "En curso"}, {"achieved", "Lograda"}, {"overdue", "Vencida"}}
	lbMoveKind    = []opt{{"deposit", "Aportación"}, {"withdrawal", "Retiro"}}
	lbCreditKind  = []opt{{"revolving_line", "Línea revolvente"}, {"bank_card", "Tarjeta bancaria"}, {"term_loan", "Préstamo a plazo"}, {"store_credit", "Crédito de tienda"}}
	lbPurchase    = []opt{{"active", "Viva"}, {"paid_off", "Liquidada"}, {"cancelled", "Cancelada"}}
	lbSchedStatus = []opt{{"pending", "Pendiente"}, {"overdue", "Vencido"}, {"paid", "Pagado"}}
	lbStoreKind   = []opt{{"chain", "Cadena"}, {"branch", "Sucursal"}, {"market", "Mercado o tianguis"}}
	lbFreshness   = []opt{{"fresh", "Precio vigente"}, {"stale", "Precio viejo"}, {"no_price", "Sin precio"}}
	lbStale       = []opt{{"false", "Vigente"}, {"true", "Viejo"}}
	lbPriceSource = []opt{{"ticket_manual", "Ticket a mano"}, {"ticket_ocr", "Ticket leído"}, {"site_assisted", "Lectura asistida"}, {"profeco", "PROFECO"}, {"manual", "Captura manual"}}
	lbOrigin      = []opt{{"claude_seed", "Recetario inicial"}, {"household", "De la casa"}}
	lbInvStatus   = []opt{{recetas.StatusSufficient, "Suficiente"}, {recetas.StatusReorder, "Reponer"}, {recetas.StatusOut, "Agotado"}}
	lbTracking    = []opt{{"stock", "Existencia"}, {"consumption_rate", "Ritmo de consumo"}}
	lbDelivery    = []opt{{"pending", "Pendiente"}, {"sent", "Enviado"}, {"failed", "Falló"}, {"skipped", "Omitido"}}
	lbJobStatus   = []opt{{"running", "Corriendo"}, {"ok", "Correcto"}, {"error", "Con error"}, {"skipped", "Omitido"}}
	lbBackupState = []opt{{"running", "En curso"}, {"ok", "Correcto"}, {"failed", "Falló"}}
	lbBackupKind  = []opt{{"scheduled", "Programado"}, {"manual", "Manual"}, {"pre_migration", "Antes de migrar"}}
)

func idGroup(col, table, name string) groupDef {
	return groupDef{key: "t." + col, label: fmt.Sprintf("SELECT %s FROM %s WHERE id = t.%s", name, table, col)}
}

const (
	sqlPaidSched  = `SELECT COALESCE(SUM(m.amount_cents), 0) FROM credit_payment m WHERE m.scheduled_payment_id = t.id`
	sqlSchedCred  = `SELECT credit_id FROM credit_purchase WHERE id = t.credit_purchase_id`
	sqlLastPrice  = `SELECT pr.observed_date || '|' || pr.source FROM price pr JOIN presentation ps ON ps.id = pr.presentation_id WHERE ps.product_id = t.id ORDER BY pr.observed_date DESC, pr.id DESC LIMIT 1`
	sqlHousehold  = `SELECT coalesce(sum(s.amount_cents), 0) FROM expense_share s JOIN person p ON p.id = s.person_id WHERE s.expense_id = t.id AND p.participation = 'holder'`
	sqlExpenseImg = `SELECT a.id FROM attachment a WHERE a.expense_id = t.id AND a.mime_type LIKE 'image/%' ORDER BY a.id LIMIT 1`
)

var viewDefs = map[string]viewDef{
	"expense": {modes: allModes, dates: map[string]string{"expense_date": "t.expense_date"}, defDate: "expense_date",
		groups: map[string]groupDef{"category_id": idGroup("category_id", "category", "name"), "kind": {key: "t.kind", enum: lbExpenseKind},
			"payment_method": {key: "t.payment_method", enum: lbPayMethod}, "person_id": idGroup("person_id", "person", "display_name")},
		title: "COALESCE(t.description, (SELECT name FROM category WHERE id = t.category_id))", amount: sqlHousehold, image: sqlExpenseImg,
		sorts: []string{"expense_date", "amount_cents", "description", "created_at"}, defSort: "-expense_date"},
	"income": {modes: []string{"table", "calendar"}, dates: map[string]string{"income_date": "t.income_date"}, defDate: "income_date",
		title: "t.description", amount: "t.amount_cents", person: "t.person_id", sorts: []string{"income_date", "amount_cents", "description"}, defSort: "-income_date"},
	"category": {modes: []string{"table", "board"}, title: "t.name", sorts: []string{"name", "sort_order"}, defSort: "sort_order",
		groups: map[string]groupDef{"default_cost_share_rule_id": idGroup("default_cost_share_rule_id", "cost_share_rule", "name")}},
	"savings_goal": {modes: allModes, dates: map[string]string{"target_date": "t.target_date"}, defDate: "target_date", title: "t.name",
		groups: map[string]groupDef{"status": {enum: lbGoalStatus, derived: true}}, status: "status", statusEnum: lbGoalStatus,
		sorts: []string{"name", "target_date", "start_date"}, defSort: "name", derive: deriveGoals, derivedAmount: true},
	"savings_movement": {modes: []string{"table", "calendar"}, dates: map[string]string{"movement_date": "t.movement_date"}, defDate: "movement_date",
		title: "(SELECT name FROM savings_goal WHERE id = t.savings_goal_id)", amount: "CASE t.kind WHEN 'withdrawal' THEN -t.amount_cents ELSE t.amount_cents END",
		groups: map[string]groupDef{"savings_goal_id": idGroup("savings_goal_id", "savings_goal", "name"), "kind": {key: "t.kind", enum: lbMoveKind}},
		sorts:  []string{"movement_date", "amount_cents"}, defSort: "-movement_date"},
	"credit": {modes: []string{"table", "board", "gallery"}, title: "t.name", groups: map[string]groupDef{"kind": {key: "t.kind", enum: lbCreditKind}},
		person: "t.holder_person_id", credit: "t.id", sorts: []string{"name"}, defSort: "name", derive: deriveCredits, derivedAmount: true},
	"credit_term_rate": {modes: []string{"table"}, title: "(SELECT name FROM credit WHERE id = t.credit_id) || ' · ' || t.term_count || ' pagos'",
		credit: "t.credit_id", sorts: []string{"term_count", "quote_date"}, defSort: "term_count"},
	"credit_purchase": {modes: []string{"table", "board", "calendar"}, dates: map[string]string{"purchase_date": "t.purchase_date"}, defDate: "purchase_date",
		title: "COALESCE(t.description, (SELECT name FROM credit WHERE id = t.credit_id))", amount: "t.total_cents", credit: "t.credit_id",
		groups: map[string]groupDef{"status": {key: "t.status", enum: lbPurchase}, "credit_id": idGroup("credit_id", "credit", "name")},
		status: "status", statusEnum: lbPurchase, sorts: []string{"purchase_date", "total_cents"}, defSort: "-purchase_date"},
	"scheduled_payment": {modes: []string{"table", "board", "calendar"}, dates: map[string]string{"due_date": "t.due_date"}, defDate: "due_date",
		title:  "COALESCE((SELECT COALESCE(cp.description, c.name) FROM credit_purchase cp JOIN credit c ON c.id = cp.credit_id WHERE cp.id = t.credit_purchase_id), 'Pago') || ' · pago ' || t.seq",
		amount: "t.amount_cents", extra: sqlPaidSched, credit: "(" + sqlSchedCred + ")",
		groups: map[string]groupDef{"status": {enum: lbSchedStatus, derived: true}, "credit_id": {key: sqlSchedCred, label: "SELECT name FROM credit WHERE id = (" + sqlSchedCred + ")"}},
		status: "status", statusEnum: lbSchedStatus, sorts: []string{"due_date", "amount_cents", "seq"}, defSort: "due_date", derive: deriveScheduled},
	"promotion": {modes: []string{"table", "board", "calendar"}, dates: map[string]string{"valid_from": "t.valid_from"}, defDate: "valid_from",
		title: "t.merchant_name", credit: "t.credit_id", store: "t.store_id", groups: map[string]groupDef{"credit_id": idGroup("credit_id", "credit", "name")},
		sorts: []string{"valid_from", "valid_to", "merchant_name"}, defSort: "-valid_from"},
	"store": {modes: allModes, dates: map[string]string{"visit_weekday": ""}, title: "t.name", store: "t.id",
		groups: map[string]groupDef{"kind": {key: "t.kind", enum: lbStoreKind}}, sorts: []string{"name"}, defSort: "name"},
	"product": {modes: []string{"table", "board", "gallery"}, title: "t.name", extra: sqlLastPrice,
		groups: map[string]groupDef{"freshness": {enum: lbFreshness, derived: true}, "product_group_id": idGroup("product_group_id", "product_group", "name")},
		sorts:  []string{"name"}, defSort: "name", derive: deriveProducts},
	"price": {modes: []string{"table", "board"}, dates: map[string]string{"observed_date": "t.observed_date"}, defDate: "observed_date",
		title: "(SELECT p.name || ' · ' || ps.name FROM presentation ps JOIN product p ON p.id = ps.product_id WHERE ps.id = t.presentation_id)", store: "t.store_id",
		groups: map[string]groupDef{"is_stale": {enum: lbStale, derived: true}, "source": {key: "t.source", enum: lbPriceSource}, "store_id": idGroup("store_id", "store", "name")},
		sorts:  []string{"observed_date", "price_cents"}, defSort: "-observed_date", derive: derivePrices},
	"recipe": {modes: []string{"table", "board", "gallery"}, title: "t.name", sorts: []string{"name"}, defSort: "name",
		groups: map[string]groupDef{"meal_slot_id": {key: "SELECT group_concat(meal_slot_id) FROM (SELECT meal_slot_id FROM recipe_meal_slot WHERE recipe_id = t.id ORDER BY meal_slot_id)", multi: true},
			"origin": {key: "t.origin", enum: lbOrigin}}},
	"meal_plan_entry": {modes: allModes, dates: map[string]string{"plan_date": "t.plan_date"}, defDate: "plan_date", person: "t.person_id",
		title:  "COALESCE((SELECT name FROM recipe WHERE id = t.recipe_id), CASE WHEN t.is_eating_out = 1 THEN 'Comer fuera' ELSE 'Sin receta' END)",
		groups: map[string]groupDef{"meal_slot_id": idGroup("meal_slot_id", "meal_slot", "name"), "person_id": idGroup("person_id", "person", "display_name")},
		sorts:  []string{"plan_date"}, defSort: "plan_date"},
	"inventory_item": {modes: []string{"table", "board", "gallery"}, title: "COALESCE((SELECT name FROM product WHERE id = t.product_id), (SELECT name FROM recipe WHERE id = t.recipe_id), 'Artículo')",
		groups: map[string]groupDef{"status": {enum: lbInvStatus, derived: true}, "tracking_mode": {key: "t.tracking_mode", enum: lbTracking}},
		status: "status", statusEnum: lbInvStatus, sorts: []string{"expires_date", "counted_date"}, defSort: "id", derive: deriveInventory},
	"weight_entry": {modes: []string{"table"}, dates: map[string]string{"measured_date": "t.measured_date"}, defDate: "measured_date", person: "t.person_id",
		title: "(SELECT display_name FROM person WHERE id = t.person_id)", groups: map[string]groupDef{"person_id": idGroup("person_id", "person", "display_name")},
		sorts: []string{"measured_date", "weight_g"}, defSort: "-measured_date"},
	"notification_delivery": {modes: []string{"table", "board", "calendar"}, dates: map[string]string{"sent_at": "t.sent_at"}, defDate: "sent_at",
		title:  "(SELECT name FROM notification WHERE id = t.notification_id)",
		groups: map[string]groupDef{"notification_id": idGroup("notification_id", "notification", "name"), "status": {key: "t.status", enum: lbDelivery}},
		status: "status", statusEnum: lbDelivery, sorts: []string{"sent_at", "created_at"}, defSort: "-created_at"},
	"job": {modes: []string{"table", "board", "calendar"}, dates: map[string]string{"next_run_at": "t.next_run_at"}, defDate: "next_run_at", title: "t.name",
		groups: map[string]groupDef{"code": {key: "t.code", label: "t.name"}, "last_status": {key: "t.last_status", enum: lbJobStatus}},
		status: "last_status", statusEnum: lbJobStatus, sorts: []string{"name", "next_run_at", "last_run_at"}, defSort: "name"},
	"person":       {modes: []string{"table"}, title: "t.display_name", person: "t.id", sorts: []string{"display_name"}, defSort: "id"},
	"user_account": {modes: []string{"table"}, title: "t.username", person: "t.person_id", sorts: []string{"username"}, defSort: "username"},
	"device_session": {modes: []string{"table"}, title: "COALESCE(t.device_label, t.user_agent, 'Sesión ' || t.id)",
		groups: map[string]groupDef{"user_account_id": idGroup("user_account_id", "user_account", "username")},
		sorts:  []string{"last_seen_at", "expires_at"}, defSort: "-last_seen_at"},
	"backup": {modes: []string{"table"}, dates: map[string]string{"started_at": "t.started_at"}, defDate: "started_at", title: "'Respaldo ' || t.id",
		groups: map[string]groupDef{"status": {key: "t.status", enum: lbBackupState}, "kind": {key: "t.kind", enum: lbBackupKind}},
		status: "status", statusEnum: lbBackupState, sorts: []string{"started_at"}, defSort: "-started_at"},
}

func optKeys(os []opt) []string {
	out := make([]string, 0, len(os))
	for _, o := range os {
		out = append(out, o.k)
	}
	return out
}

func optLabel(os []opt, k string) string {
	for _, o := range os {
		if o.k == k {
			return o.l
		}
	}
	if k == "" {
		return "Sin valor"
	}
	return k
}

func inIDs(expr string, ids []int64) (string, []any) {
	ph := strings.TrimSuffix(strings.Repeat("?, ", len(ids)), ", ")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	return expr + " IN (" + ph + ")", args
}

func recInt(m map[string]any, k string) (int64, bool) {
	v, ok := m[k].(int64)
	return v, ok
}

func recStr(m map[string]any, k string) string {
	s, _ := m[k].(string)
	return s
}

func dayOf(d dominio.Date) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

// deriveGoals: estado y ahorrado de cada meta salen del motor S2 (SavingsGoals; regla S2 §3.13).
func deriveGoals(reg *api.Registry, r *http.Request, recs []dominio.ViewRecord) error {
	_, in, err := engineAt(reg, r)
	if err != nil {
		return err
	}
	sv, err := mf.SavingsGoals(in)
	if err != nil {
		return err
	}
	by := map[int64]mf.GoalStatus{}
	for _, g := range sv.Goals {
		by[g.GoalID] = g
	}
	for i := range recs {
		if g, ok := by[recs[i].ID]; ok {
			recs[i].Record["status"], recs[i].Record["saved_cents"] = g.State, int64(g.SavedCents)
			recs[i].Amount.Int64, recs[i].Amount.Valid = int64(g.SavedCents), true
		}
	}
	return nil
}

// deriveCredits: deuda viva por crédito con motor_cred.LiveDebt (misma entrada que la tarjeta de próximo pago).
func deriveCredits(reg *api.Registry, r *http.Request, recs []dominio.ViewRecord) error {
	t := today(reg)
	ci, err := loadCredInput(r.Context(), reg.Deps().DB, t)
	if err != nil {
		return err
	}
	for i := range recs {
		for _, c := range ci.credits {
			if c.ID != recs[i].ID {
				continue
			}
			if debt, derr := mc.LiveDebt([]mc.CreditState{ci.state(c)}, cdateOf(t)); derr == nil { // sin cálculo posible: se omite
				recs[i].Record["live_debt_cents"] = int64(debt.Total)
				recs[i].Amount.Int64, recs[i].Amount.Valid = int64(debt.Total), true
			}
		}
	}
	return nil
}

// deriveScheduled: estado del pago programado con la regla de api-creditos (pagado si lo abonado cubre el monto,
// vencido si la fecha ya pasó, si no pendiente).
func deriveScheduled(reg *api.Registry, _ *http.Request, recs []dominio.ViewRecord) error {
	t := today(reg).String()
	for i := range recs {
		paid, _ := strconv.ParseInt(recs[i].Extra.String, 10, 64)
		amount, _ := recInt(recs[i].Record, "amount_cents")
		st := "pending"
		switch {
		case paid >= amount:
			st = "paid"
		case recStr(recs[i].Record, "due_date") < t:
			st = "overdue"
		}
		recs[i].Record["status"] = st
	}
	return nil
}

func priceSettings(ctx context.Context, q dominio.Querier) (precios.Settings, error) {
	st, err := dominio.SettingInt(ctx, q, "price.stale_after_days")
	if err != nil {
		return precios.Settings{}, err
	}
	zw, err := dominio.SettingInt(ctx, q, "price.zero_weight_after_days")
	return precios.Settings{StaleAfterDays: int(st), ZeroWeightAfterDays: int(zw)}, err
}

func freshness(observed, source string, t time.Time, s precios.Settings) (bool, bool) {
	d, err := time.Parse(time.DateOnly, observed)
	if err != nil {
		return false, false
	}
	w, err := precios.Freshness(d, t, s, source)
	if err != nil {
		return false, false
	}
	return w.IsStale, true
}

// derivePrices: «viejo» de cada precio con precios.Freshness (K22) y los parámetros vigentes.
func derivePrices(reg *api.Registry, r *http.Request, recs []dominio.ViewRecord) error {
	s, err := priceSettings(r.Context(), reg.Deps().DB)
	if err != nil {
		return err
	}
	t := dayOf(today(reg))
	for i := range recs {
		if stale, ok := freshness(recStr(recs[i].Record, "observed_date"), recStr(recs[i].Record, "source"), t, s); ok {
			recs[i].Record["is_stale"] = stale
		}
	}
	return nil
}

// deriveProducts: frescura del producto según su último precio (PROVISIONAL, bloqueo §16: api-compras-a la da con QuoteFor).
func deriveProducts(reg *api.Registry, r *http.Request, recs []dominio.ViewRecord) error {
	s, err := priceSettings(r.Context(), reg.Deps().DB)
	if err != nil {
		return err
	}
	t := dayOf(today(reg))
	for i := range recs {
		f := "no_price"
		if date, src, ok := strings.Cut(recs[i].Extra.String, "|"); ok {
			if stale, ok := freshness(date, src, t, s); ok {
				f = map[bool]string{true: "stale", false: "fresh"}[stale]
			}
		}
		recs[i].Record["freshness"] = f
	}
	return nil
}

// deriveInventory: estado de cada artículo con el motor S4 (recetas.InventoryStatus, K7).
func deriveInventory(reg *api.Registry, _ *http.Request, recs []dominio.ViewRecord) error {
	td, err := recetas.ParseDate(today(reg).String())
	if err != nil {
		return err
	}
	e := &recetas.Engine{Today: td}
	opt := func(m map[string]any, k string) *int64 {
		if v, ok := recInt(m, k); ok {
			return &v
		}
		return nil
	}
	for _, it := range recs {
		q, _ := recInt(it.Record, "quantity_milli")
		mi := recetas.InventoryItem{Code: strconv.FormatInt(it.ID, 10), Mode: recStr(it.Record, "tracking_mode"), QuantityMilli: q,
			RateQuantityMilli: opt(it.Record, "rate_quantity_milli"), RatePeriodDays: opt(it.Record, "rate_period_days"), ReorderLevelMilli: opt(it.Record, "reorder_level_milli")}
		if cd := recStr(it.Record, "counted_date"); cd != "" {
			d, err := recetas.ParseDate(cd)
			if err != nil {
				return err
			}
			mi.CountedDate = &d
		}
		e.Inventory = append(e.Inventory, mi)
	}
	if len(e.Inventory) == 0 {
		return nil
	}
	rep, err := e.InventoryStatus()
	if err != nil {
		return err
	}
	for i, s := range rep.Items {
		recs[i].Record["status"] = s.Status
	}
	return nil
}

var (
	weekdaysES = []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
	monthsES   = []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
)

func dayLabel(s string) string {
	d, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%s %d de %s de %d", weekdaysES[d.Weekday()], d.Day(), monthsES[d.Month()-1], d.Year())
}

func unprocessable(field, msg string) error {
	return dominio.Invalid(dominio.FieldIssue{Field: field, Code: "invalid_value", Message: msg})
}

// queryView (C1 §8, decisión G): vista genérica de tabla, tablero, calendario o galería de una entidad.
func queryView(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		entity := r.PathValue("entity")
		def, ok := viewDefs[entity]
		if !ok {
			return dominio.NotFound("No existe la vista «" + entity + "». Revisa el nombre de la entidad.")
		}
		qs := r.URL.Query()
		mode := qs.Get("mode")
		if !slices.Contains(allModes, mode) {
			return api.BadRequest("El modo de la vista debe ser table, board, calendar o gallery.")
		}
		if !slices.Contains(def.modes, mode) {
			return unprocessable("mode", "Esta vista no tiene el modo «"+mode+"». Usa "+strings.Join(def.modes, ", ")+".")
		}
		groupBy, dateField, status := qs.Get("group_by"), qs.Get("date_field"), qs.Get("status")
		for name, v := range map[string]string{"group_by": groupBy, "date_field": dateField} {
			if v != "" && !viewIdent(v) {
				return api.BadRequest("El parámetro " + name + " no tiene un formato válido.")
			}
		}
		if status != "" && !viewIdent(status) {
			return api.BadRequest("El filtro status no tiene un formato válido.")
		}
		var gd groupDef
		if groupBy != "" {
			if gd, ok = def.groups[groupBy]; !ok {
				return unprocessable("group_by", "No se puede agrupar esta vista por «"+groupBy+"».")
			}
		} else if mode == "board" {
			keys := make([]string, 0, len(def.groups))
			for k := range def.groups {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			groupBy, gd = keys[0], def.groups[keys[0]] // PROVISIONAL (bloqueo §16): sin group_by, la primera propiedad permitida
		}
		if dateField != "" {
			if _, ok := def.dates[dateField]; !ok {
				return unprocessable("date_field", "Esta vista no tiene la fecha «"+dateField+"».")
			}
		} else {
			dateField = def.defDate
			if dateField == "" {
				for k := range def.dates {
					dateField = k
				}
			}
		}
		if status != "" && (def.status == "" || !slices.Contains(optKeys(def.statusEnum), status)) {
			return unprocessable("status", "El estado «"+status+"» no aplica a esta vista.")
		}
		fam := qs.Get("family")
		if fam != "" && !slices.Contains(families, fam) {
			return api.BadRequest("El filtro family solo admite household, external o joint.")
		}
		text := qs.Get("q")
		if len([]rune(text)) > 200 {
			return api.BadRequest("El texto de búsqueda admite hasta 200 caracteres.")
		}
		period, err := api.QueryPeriod(r)
		if err != nil {
			return err
		}
		p, err := api.ParsePage(r, def.sorts, def.defSort)
		if err != nil {
			return err
		}
		cols := map[string]string{"id": "t.id"}
		for _, s := range def.sorts {
			cols[s] = "t." + s
		}
		s := dominio.ViewSQL{Table: entity, Title: def.title, Amount: def.amount, Image: def.image, Extra: def.extra,
			Group: gd.key, GroupLabel: gd.label, OrderBy: p.OrderBy(cols)}
		if mode != "gallery" { // la foto solo se pide en la galería
			s.Image = ""
		}
		if dateField != "" {
			s.Date = def.dates[dateField]
		}
		sqlFilter := entity == "expense" // gastos: los mismos filtros que listExpenses (C1-filtro-familia opción A)
		if sqlFilter {
			f, err := parseFilter(reg, r, false)
			if err != nil {
				return err
			}
			f.Period = period
			where, args := f.Where()
			s.Alias, s.Where, s.Args = "e", []string{where}, args
			s.Title, s.Amount = strings.ReplaceAll(s.Title, "t.", "e."), strings.ReplaceAll(sqlHousehold, "t.id", "e.id")
			s.Image = strings.ReplaceAll(s.Image, "t.id", "e.id")
			s.Group, s.GroupLabel = strings.ReplaceAll(s.Group, "t.", "e."), strings.ReplaceAll(s.GroupLabel, "t.", "e.")
			s.Date = strings.ReplaceAll(s.Date, "t.", "e.")
			s.OrderBy = strings.ReplaceAll(s.OrderBy, "t.", "e.")
		} else {
			for _, fl := range []struct{ param, expr string }{{"person_id", def.person}, {"credit_id", def.credit}, {"store_id", def.store}} {
				ids, err := api.QueryIDs(r, fl.param)
				if err != nil {
					return err
				}
				if len(ids) > 0 && fl.expr != "" { // un filtro global que no aplica a la entidad no la filtra (PROVISIONAL, bloqueo §16)
					wh, args := inIDs(fl.expr, ids)
					s.Where, s.Args = append(s.Where, wh), append(s.Args, args...)
				}
			}
			if _, err := api.QueryIDs(r, "category_id"); err != nil {
				return err
			}
		}
		recs, err := dominio.GastosViewRows(r.Context(), reg.Deps().DB, s)
		if err != nil {
			return err
		}
		if def.derive != nil {
			if err := def.derive(reg, r, recs); err != nil {
				return err
			}
		}
		recs, err = viewExpandStore(reg, entity, mode, period, recs)
		if err != nil {
			return err
		}
		items, groups := viewProject(recs, def, mode, groupBy, gd, dateField, status, text, period, sqlFilter)
		end := min(len(items), p.Offset+p.Limit+1)
		pg := api.PageOf(items[min(p.Offset, len(items)):end], p)
		return api.WriteJSON(w, http.StatusOK, api.ViewPage{Items: pg.Items, NextCursor: pg.NextCursor, Mode: mode, Groups: groups})
	}
}

func viewIdent(s string) bool {
	for i, c := range s {
		if !(c >= 'a' && c <= 'z' || c == '_' || i > 0 && c >= '0' && c <= '9') {
			return false
		}
	}
	return s != ""
}

// viewExpandStore: en el calendario, un mercado aparece en cada fecha de su día de visita (ISO 1 = lunes) dentro del
// periodo pedido o, sin periodo, del mes en curso (PROVISIONAL, bloqueo §16).
func viewExpandStore(reg *api.Registry, entity, mode string, period *dominio.Period, recs []dominio.ViewRecord) ([]dominio.ViewRecord, error) {
	if entity != "store" || mode != "calendar" {
		return recs, nil
	}
	pr := dominio.MonthPeriod(today(reg))
	if period != nil {
		pr = *period
	}
	out := []dominio.ViewRecord{}
	for d := dayOf(pr.Start); !d.After(dayOf(pr.End)); d = d.AddDate(0, 0, 1) {
		iso := int64(d.Weekday())
		if iso == 0 {
			iso = 7
		}
		for _, rc := range recs {
			if wd, ok := recInt(rc.Record, "visit_weekday"); ok && wd == iso {
				rc.Date.String, rc.Date.Valid = d.Format(time.DateOnly), true
				out = append(out, rc)
			}
		}
	}
	return out, nil
}

// viewProject aplica los filtros calculados (status, q, period), arma los ViewItem y los grupos con conteo,
// total y participación (basis points, ROUND_HALF_UP).
func viewProject(recs []dominio.ViewRecord, def viewDef, mode, groupBy string, gd groupDef, dateField, status, text string,
	period *dominio.Period, sqlFiltered bool) ([]api.ViewItem, []api.ViewGroup) {
	needle := dominio.GastosFoldText(text)
	items := []api.ViewItem{}
	type agg struct {
		label        string
		count, total int64
	}
	byKey := map[string]*agg{}
	order := []string{}
	if mode != "calendar" && groupBy != "" {
		for _, o := range gd.enum {
			byKey[o.k], order = &agg{label: o.l}, append(order, o.k)
		}
	}
	var viewTotal int64
	for _, rc := range recs {
		if status != "" && fmt.Sprint(rc.Record[def.status]) != status {
			continue
		}
		date := ""
		if rc.Date.Valid {
			date = rc.Date.String
			if strings.HasSuffix(dateField, "_at") { // instante → día local (America/Mexico_City)
				date = dominio.DBTimeToAPI(date)
			}
			if len(date) > 10 {
				date = date[:10]
			}
		}
		if !sqlFiltered {
			if period != nil && dateField != "" && def.defDate != "" && (date == "" || date < period.Start.String() || date > period.End.String()) {
				continue
			}
			if needle != "" && !viewMatches(rc, needle) {
				continue
			}
		}
		keys := []string{""}
		labels := map[string]string{}
		switch {
		case mode == "calendar":
			if date == "" {
				continue
			}
			keys = []string{date}
			labels[date] = dayLabel(date)
		case groupBy == "":
			keys = nil
		case gd.derived:
			if v, ok := rc.Record[groupBy]; ok {
				keys = []string{fmt.Sprint(v)}
			}
		case gd.multi:
			if rc.Group.String != "" {
				keys = strings.Split(rc.Group.String, ",")
			}
		default:
			keys = []string{rc.Group.String}
			if rc.GroupLabel.Valid {
				labels[rc.Group.String] = rc.GroupLabel.String
			}
		}
		base := api.ViewItem{ID: rc.ID, Title: rc.Title, Record: rc.Record}
		if date != "" {
			base.Date = api.Ptr(date)
		}
		if rc.Amount.Valid {
			base.AmountCents = api.Ptr(rc.Amount.Int64)
		}
		if rc.Image.Valid {
			base.ImageAttachmentID = api.Ptr(rc.Image.Int64)
		}
		if keys == nil {
			items = append(items, base)
			viewTotal += rc.Amount.Int64
			continue
		}
		for _, k := range keys {
			it := base
			it.GroupKey = api.Ptr(k)
			items = append(items, it)
			viewTotal += rc.Amount.Int64
			a, ok := byKey[k]
			if !ok {
				l := labels[k]
				if l == "" {
					l = optLabel(gd.enum, k)
				}
				a = &agg{label: l}
				byKey[k], order = a, append(order, k)
			}
			a.count++
			a.total += rc.Amount.Int64
		}
	}
	if mode == "calendar" {
		sort.Strings(order)
	} else if gd.enum == nil {
		sort.SliceStable(order, func(i, j int) bool {
			return dominio.GastosFoldText(byKey[order[i]].label) < dominio.GastosFoldText(byKey[order[j]].label)
		})
	}
	groups := []api.ViewGroup{}
	for _, k := range order {
		a := byKey[k]
		g := api.ViewGroup{Key: k, Label: a.label, Count: a.count}
		if def.amount != "" || def.derivedAmount {
			g.TotalCents = api.Ptr(a.total)
			if viewTotal != 0 {
				g.ShareBp = api.Ptr(dominio.MulDivRound(a.total, 10000, viewTotal))
			}
		}
		groups = append(groups, g)
	}
	return items, groups
}

func viewMatches(rc dominio.ViewRecord, needle string) bool {
	if strings.Contains(dominio.GastosFoldText(rc.Title), needle) {
		return true
	}
	for k, v := range rc.Record {
		if s, ok := v.(string); ok && !strings.HasSuffix(k, "_at") && strings.Contains(dominio.GastosFoldText(s), needle) {
			return true
		}
	}
	return false
}
