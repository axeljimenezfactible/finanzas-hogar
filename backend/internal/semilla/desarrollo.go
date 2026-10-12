package semilla

// Semilla de DESARROLLO (WF-3a bis, agente semilla-dev, 12-oct-2026; C3-matriz punto 5, opción b).
//
// Es un overlay APARTE de la semilla real: llena lo que Inicio y las capturas de U5 necesitan para oct-2026
// (presupuesto, sueldo quincenal, un gasto fijo, una meta de ahorro con aportación, precios de la compra del
// 10-oct y un plan de comidas de la semana). NUNCA se carga en producción: «finanzas semilla --desarrollo» se
// niega salvo en una base de desarrollo (ver devAllowed y ops/runbook.md).
//
// De dónde sale cada valor (regla «no inventes montos»): los montos y las fechas se LEEN en tiempo de ejecución de
// las fixtures aprobadas del oráculo S2 (oraculo/fixtures/s2/*.json); aquí no hay ninguna cifra de dinero escrita
// a mano. Lo que ninguna fixture trae (el plan de comidas) se arma con el recetario y los precios reales ya
// sembrados y se marca como dato de desarrollo.
//
//	Qué                                   Fixture (parche de entradas.datos)                    Id
//	Sueldo quincenal 12,000.00 / 12,500.00  s2-proyeccion-csv · planned_items                      1
//	Renta mensual 8,000.00 (Vivienda)       s2-proyeccion-csv · planned_items (+ categories)       2  (solo con --con-gasto-fijo)
//	Presupuesto Despensa 2,000.00 desde oct s2-presupuesto-scope · budgets                         1
//	Meta de ahorro 6,293.44 al 10-abr-2027  s2-metas · savings_goals (+ savings_movements)         2
//
// Idempotente igual que la semilla real: cada entidad se busca por su llave natural y solo se crea si falta.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"finanzas-hogar/internal/dominio"
)

// RequestIDDev es el request_id con que la semilla de desarrollo firma audit_log (se distingue de la real).
const RequestIDDev = "semilla-desarrollo"

// Sufijos y nombres de desarrollo: se ven en pantalla, para que nadie confunda estos datos con los de la casa.
const (
	devSuffix   = "(desarrollo)"
	devGoalName = "Meta de ahorro " + devSuffix
	devNote     = "Dato de desarrollo (finanzas semilla --desarrollo): no existe en producción."
)

// Fixtures y ids que usa el overlay (ver la tabla de arriba).
const (
	fixPlanned = "s2-proyeccion-csv.json"
	fixBudget  = "s2-presupuesto-scope.json"
	fixGoals   = "s2-metas.json"

	fixSueldoID = 1 // planned_items[id=1] de s2-proyeccion-csv: ingreso quincenal
	fixRentaID  = 2 // planned_items[id=2] de s2-proyeccion-csv: gasto fijo mensual
	fixBudgetID = 1 // budgets[id=1] de s2-presupuesto-scope: Despensa desde 2026-10-01
	fixGoalID   = 2 // savings_goals[id=2] de s2-metas («Con residuo (RN-18)») y su aportación
)

// Plan de comidas (dato de desarrollo, sin fixture): lunes 12 a sábado 17-oct-2026, tiempo «comida», Axel y Ximena,
// una porción cada uno. Son recetas del recetario con tiempo «comida» y con el mayor número de ingredientes que
// tienen precio en la compra del 10-oct.
var (
	devMealPlan = []struct{ Day, Recipe string }{
		{"2026-10-12", "Espagueti con atún y elote"},
		{"2026-10-13", "Arroz con pollo y zanahoria"},
		{"2026-10-14", "Sopa de fideo con jitomate"},
		{"2026-10-15", "Rajas poblanas con crema y elote"},
		{"2026-10-16", "Espagueti con calabacita y queso manchego"},
		{"2026-10-17", "Arroz rojo con huevo estrellado y plátano macho"},
	}
	devMealPersons = []string{"axel", "ximena"}
)

// fixture* leen solo lo que usa el overlay del «parche» de una fixture del oráculo.
type fixtureFile struct {
	Entradas struct {
		Datos struct {
			Parche struct {
				Categories []struct {
					Code     string  `json:"code"`
					Name     string  `json:"name"`
					Kind     *string `json:"kind"`
					IsActive bool    `json:"is_active"`
				} `json:"categories"`
				PlannedItems []fixPlannedItem `json:"planned_items"`
				Budgets      []struct {
					ID          int64   `json:"id"`
					Category    string  `json:"category"`
					AmountCents int64   `json:"amount_cents"`
					ValidFrom   string  `json:"valid_from"`
					ValidTo     *string `json:"valid_to"`
				} `json:"budgets"`
				SavingsGoals []struct {
					ID                  int64   `json:"id"`
					Kind                string  `json:"kind"`
					TargetCents         *int64  `json:"target_cents"`
					TargetDate          *string `json:"target_date"`
					PlannedMonthlyCents *int64  `json:"planned_monthly_cents"`
					StartDate           string  `json:"start_date"`
					IsActive            bool    `json:"is_active"`
				} `json:"savings_goals"`
				SavingsMovements []struct {
					GoalID       int64  `json:"goal_id"`
					MovementDate string `json:"movement_date"`
					Kind         string `json:"kind"`
					AmountCents  int64  `json:"amount_cents"`
				} `json:"savings_movements"`
			} `json:"parche"`
		} `json:"datos"`
	} `json:"entradas"`
}

type fixPlannedItem struct {
	ID              int64   `json:"id"`
	Kind            string  `json:"kind"`
	Description     string  `json:"description"`
	AmountCents     int64   `json:"amount_cents"`
	SecondHalfCents *int64  `json:"second_half_amount_cents"`
	Frequency       string  `json:"frequency"`
	IntervalCount   int64   `json:"interval_count"`
	StartDate       string  `json:"start_date"`
	EndDate         *string `json:"end_date"`
	Category        *string `json:"category"`
	PaymentMethod   *string `json:"payment_method"`
	AccrualMonths   *int64  `json:"accrual_months"`
	Person          *string `json:"person"`
	ShareRule       *struct {
		Parts      []json.RawMessage `json:"parts"`
		RestPerson string            `json:"rest_person"`
	} `json:"share_rule"`
	CreditID *int64 `json:"credit_id"`
	IsActive bool   `json:"is_active"`
}

// fixtureDir resuelve oraculo/fixtures/s2: la ruta dada o Dir/../../oraculo/fixtures/s2.
func fixtureDir(dir, explicit string) (string, error) {
	for _, c := range []string{explicit, filepath.Join(dir, "..", "..", "oraculo", "fixtures", "s2")} {
		if c == "" {
			continue
		}
		if st, err := os.Stat(c); err == nil && st.IsDir() {
			return c, nil
		}
	}
	return "", errors.New("semilla (desarrollo): no encuentro oraculo/fixtures/s2: los montos de la semilla de desarrollo salen de esas fixtures y no se inventan")
}

func readFixture(dir, name string) (*fixtureFile, error) {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return nil, fmt.Errorf("semilla (desarrollo): %w", err)
	}
	var f fixtureFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("semilla (desarrollo): %s: %w", name, err)
	}
	return &f, nil
}

// devAllowed decide si esta base admite la semilla de desarrollo. PROVISIONAL (bloqueo WF3a-bis-semilla-dev, punto 2):
// se niega si FINANZAS_ENTORNO=produccion; si no, solo en una base cuyo nombre de archivo contiene «dev»
// (datos/finanzas-dev.db) o con FINANZAS_SEMILLA_DEV=1 puesto a propósito.
func devAllowed(ctx context.Context, db *sql.DB) error {
	if strings.EqualFold(os.Getenv("FINANZAS_ENTORNO"), "produccion") {
		return errors.New("la semilla de desarrollo nunca se carga en producción (FINANZAS_ENTORNO=produccion)")
	}
	if os.Getenv("FINANZAS_SEMILLA_DEV") == "1" {
		return nil
	}
	var seq int
	var name, file string
	rows, err := db.QueryContext(ctx, `PRAGMA database_list`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&seq, &name, &file); err != nil {
			return err
		}
		if name == "main" && strings.Contains(strings.ToLower(filepath.Base(file)), "dev") {
			return nil
		}
	}
	return errors.New("la semilla de desarrollo solo se carga en una base de desarrollo (nombre con «dev», p. ej. datos/finanzas-dev.db); " +
		"si de verdad es una base de pruebas, pon FINANZAS_SEMILLA_DEV=1. En producción no se carga nunca (ops/runbook.md)")
}

// RunDev es la firma de «finanzas semilla --desarrollo»: aplica el candado, carga el overlay y imprime el reporte.
//
// conGastoFijo agrega la renta mensual (gasto fijo). Por omisión NO se carga: el motor S2 suma las ocurrencias de las
// partidas de gasto al devengado y al efectivo de cada mes, así que con la renta P1–P6 dejan de ser las cifras de
// model.json (+8,000.00 por mes) y C3-09 no cuadraría. Ver el bloqueo WF3a-bis-semilla-dev, punto 1.
func RunDev(ctx context.Context, db *sql.DB, dir string, conGastoFijo bool) error {
	if err := devAllowed(ctx, db); err != nil {
		return fmt.Errorf("semilla --desarrollo: %w", err)
	}
	rep, err := RunDevWith(ctx, db, Options{Dir: dir, DevFixedExpense: conGastoFijo})
	if rep != nil {
		fmt.Print(rep.String())
	}
	return err
}

// RunDevWith carga el overlay con opciones explícitas (sin el candado de devAllowed: lo aplica RunDev). Exige que la
// semilla real ya esté cargada.
func RunDevWith(ctx context.Context, db *sql.DB, opt Options) (*Report, error) {
	if opt.Dir == "" {
		return nil, errors.New("semilla (desarrollo): falta la carpeta de la semilla (datos/semilla)")
	}
	fdir, err := fixtureDir(opt.Dir, opt.FixturesDir)
	if err != nil {
		return nil, err
	}
	rep := &Report{}
	err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		l := &loader{ctx: ctx, tx: tx, opt: opt, rep: rep, ids: map[string]int64{}, pres: map[string]int64{}, argo: dominio.DefaultArgon2, req: RequestIDDev}
		for _, s := range []func() error{
			l.devBase,
			func() error {
				if !opt.DevFixedExpense {
					return nil
				}
				return l.devCategoria(fdir)
			},
			func() error { return l.devPresupuesto(fdir) },
			func() error { return l.devPartidas(fdir) },
			func() error { return l.devAhorro(fdir) },
			l.devPrecios,
			l.devComidas,
		} {
			if err := s(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if opt.DevFixedExpense {
		rep.Notices = append(rep.Notices, "Con el gasto fijo (renta 8,000.00 al mes) el devengado y el efectivo de cada mes suben 8,000.00: P1 y P2 de getAccrualVsCash dejan de ser las cifras de model.json (C3-09 no cuadraría con esta carga).")
	} else {
		rep.Notices = append(rep.Notices, "Sin gasto fijo (por omisión): getAccrualVsCash queda igual que con la semilla real sola (P1 y P2 = model.json; P3 a P6 ya difieren sin este overlay, bloqueo WF3a-bis-semilla-dev punto 4). Para agregar la renta: finanzas semilla --desarrollo --con-gasto-fijo.")
	}
	rep.Notices = append(rep.Notices,
		"Semilla de DESARROLLO: no se carga en producción. Montos tomados de las fixtures del oráculo S2 (oraculo/fixtures/s2); plan de comidas armado con el recetario y los precios sembrados.",
		"No hay ingresos capturados: la quincena del 15-oct aún no llega, así que Inicio los muestra como esperados (income_is_expected).")
	return rep, nil
}

// devBase comprueba que la semilla real ya está y junta los ids que se usan después.
func (l *loader) devBase() error {
	for _, code := range []string{"axel", "ximena"} {
		id, ok, err := l.find(`SELECT id FROM person WHERE code = ?`, code)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("semilla (desarrollo): falta la semilla real (no hay personas): corre primero «finanzas semilla»")
		}
		l.ids["person:"+code] = id
	}
	id, ok, err := l.find(`SELECT id FROM category WHERE code = 'groceries'`)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("semilla (desarrollo): falta la categoría Despensa: corre primero «finanzas semilla»")
	}
	l.ids["category:groceries"] = id
	if _, ok, err := l.find(`SELECT id FROM expense WHERE expense_date = '2026-10-10' AND category_id = ?`, id); err != nil {
		return err
	} else if !ok {
		return errors.New("semilla (desarrollo): falta la despensa del 10-oct-2026: corre primero «finanzas semilla»")
	}
	return nil
}

// catID traduce el código de categoría de las fixtures («despensa») al de la base («groceries»).
func (l *loader) catID(code string) (int64, error) {
	if code == "despensa" {
		code = "groceries"
	}
	id, ok, err := l.find(`SELECT id FROM category WHERE code = ?`, code)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("semilla (desarrollo): la categoría «%s» no existe", code)
	}
	return id, nil
}

// devCategoria: «Vivienda» (fija) de s2-proyeccion-csv · categories, que necesita la renta.
func (l *loader) devCategoria(fdir string) error {
	f, err := readFixture(fdir, fixPlanned)
	if err != nil {
		return err
	}
	for _, c := range f.Entradas.Datos.Parche.Categories {
		if c.Code != "vivienda" {
			continue
		}
		kind := any(nil)
		if c.Kind != nil {
			kind = *c.Kind
		}
		id, _, err := l.ensure("[desarrollo] categorías", "category", `SELECT id FROM category WHERE code = ?`, a(c.Code),
			`INSERT INTO category (parent_id, code, name, kind, default_cost_share_rule_id, sort_order, is_active) VALUES (NULL, ?, ?, ?, NULL, 3, ?)`,
			c.Code, c.Name, kind, c.IsActive)
		if err != nil {
			return err
		}
		l.ids["category:vivienda"] = id
		return nil
	}
	return fmt.Errorf("semilla (desarrollo): %s ya no trae la categoría «vivienda»", fixPlanned)
}

// devPresupuesto: presupuesto mensual de Despensa de s2-presupuesto-scope · budgets[1] (2,000.00 desde 2026-10-01).
func (l *loader) devPresupuesto(fdir string) error {
	f, err := readFixture(fdir, fixBudget)
	if err != nil {
		return err
	}
	for _, b := range f.Entradas.Datos.Parche.Budgets {
		if b.ID != fixBudgetID {
			continue
		}
		cat, err := l.catID(b.Category)
		if err != nil {
			return err
		}
		var to any
		if b.ValidTo != nil {
			to = *b.ValidTo
		}
		_, _, err = l.ensure("[desarrollo] presupuesto por categoría", "budget", `SELECT id FROM budget WHERE category_id = ? AND valid_from_date = ?`, a(cat, b.ValidFrom),
			`INSERT INTO budget (category_id, amount_cents, valid_from_date, valid_to_date, notes) VALUES (?, ?, ?, ?, ?)`,
			cat, b.AmountCents, b.ValidFrom, to, devNote+" Fixture "+fixBudget+", budgets[1].")
		return err
	}
	return fmt.Errorf("semilla (desarrollo): %s ya no trae el presupuesto %d", fixBudget, fixBudgetID)
}

// devPartidas: sueldo quincenal (ingreso) y renta mensual (gasto fijo) de s2-proyeccion-csv · planned_items 1 y 2.
// Las dos quincenas de octubre son las ocurrencias del 15 (12,000.00) y del 30 (12,500.00) de la partida del sueldo.
func (l *loader) devPartidas(fdir string) error {
	f, err := readFixture(fdir, fixPlanned)
	if err != nil {
		return err
	}
	seen, want := 0, 1
	if l.opt.DevFixedExpense {
		want = 2
	}
	for _, p := range f.Entradas.Datos.Parche.PlannedItems {
		if p.ID != fixSueldoID && !(p.ID == fixRentaID && l.opt.DevFixedExpense) {
			continue
		}
		seen++
		desc := strings.TrimSuffix(p.Description, " (dato de fixture)") + " " + devSuffix
		var cat, person any
		if p.Category != nil {
			id, err := l.catID(*p.Category)
			if err != nil {
				return err
			}
			cat = id
		}
		if p.Person != nil {
			id, ok := l.ids["person:"+*p.Person]
			if !ok {
				return fmt.Errorf("semilla (desarrollo): la persona «%s» de la partida %d no es Axel ni Ximena", *p.Person, p.ID)
			}
			person = id
		}
		// La fixture trae la regla «hogar» ({parts: [], rest_person: axel}): en la base es la regla por defecto (NULL).
		if p.ShareRule != nil && (len(p.ShareRule.Parts) != 0 || p.ShareRule.RestPerson != "axel") {
			return fmt.Errorf("semilla (desarrollo): la partida %d trae un reparto distinto de «Hogar»; esta semilla solo sabe el de la regla por defecto", p.ID)
		}
		acc := int64(1)
		if p.AccrualMonths != nil {
			acc = *p.AccrualMonths
		}
		var second, end, pay any
		if p.SecondHalfCents != nil {
			second = *p.SecondHalfCents
		}
		if p.EndDate != nil {
			end = *p.EndDate
		}
		if p.PaymentMethod != nil {
			pay = *p.PaymentMethod
		}
		if p.CreditID != nil {
			return fmt.Errorf("semilla (desarrollo): la partida %d trae crédito y no se carga", p.ID)
		}
		_, _, err := l.ensure("[desarrollo] partidas planificadas", "planned_item", `SELECT id FROM planned_item WHERE kind = ? AND description = ?`, a(p.Kind, desc),
			`INSERT INTO planned_item (kind, description, amount_cents, second_half_amount_cents, frequency, interval_count, start_date, end_date,
				category_id, person_id, cost_share_rule_id, payment_method, credit_id, term_count, accrual_months, is_active, notes)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, NULL, NULL, ?, ?, ?)`,
			p.Kind, desc, p.AmountCents, second, p.Frequency, p.IntervalCount, p.StartDate, end, cat, person, pay, acc, p.IsActive,
			devNote+fmt.Sprintf(" Fixture %s, planned_items[%d].", fixPlanned, p.ID))
		if err != nil {
			return err
		}
	}
	if seen != want {
		return fmt.Errorf("semilla (desarrollo): %s ya no trae las partidas %d y %d que se pidieron", fixPlanned, fixSueldoID, fixRentaID)
	}
	return nil
}

// devAhorro: meta con aportación de s2-metas · savings_goals[2] («Con residuo (RN-18)») y su depósito del 5-oct.
func (l *loader) devAhorro(fdir string) error {
	f, err := readFixture(fdir, fixGoals)
	if err != nil {
		return err
	}
	var gid int64
	found := false
	for _, g := range f.Entradas.Datos.Parche.SavingsGoals {
		if g.ID != fixGoalID {
			continue
		}
		found = true
		var target, tdate, monthly any
		if g.TargetCents != nil {
			target = *g.TargetCents
		}
		if g.TargetDate != nil {
			tdate = *g.TargetDate
		}
		if g.PlannedMonthlyCents != nil {
			monthly = *g.PlannedMonthlyCents
		}
		gid, _, err = l.ensure("[desarrollo] metas de ahorro", "savings_goal", `SELECT id FROM savings_goal WHERE name = ? COLLATE NOCASE`, a(devGoalName),
			`INSERT INTO savings_goal (name, kind, target_cents, target_date, planned_monthly_cents, start_date, is_active, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			devGoalName, g.Kind, target, tdate, monthly, g.StartDate, g.IsActive, devNote+fmt.Sprintf(" Fixture %s, savings_goals[%d].", fixGoals, g.ID))
		if err != nil {
			return err
		}
	}
	if !found {
		return fmt.Errorf("semilla (desarrollo): %s ya no trae la meta %d", fixGoals, fixGoalID)
	}
	n := 0
	for _, m := range f.Entradas.Datos.Parche.SavingsMovements {
		if m.GoalID != fixGoalID {
			continue
		}
		n++
		if _, _, err := l.ensure("[desarrollo] aportaciones de ahorro", "savings_movement",
			`SELECT id FROM savings_movement WHERE savings_goal_id = ? AND movement_date = ? AND kind = ? AND amount_cents = ?`, a(gid, m.MovementDate, m.Kind, m.AmountCents),
			`INSERT INTO savings_movement (savings_goal_id, movement_date, kind, amount_cents, notes) VALUES (?, ?, ?, ?, ?)`,
			gid, m.MovementDate, m.Kind, m.AmountCents, devNote+" Fixture "+fixGoals+", savings_movements."); err != nil {
			return err
		}
	}
	if n == 0 {
		return fmt.Errorf("semilla (desarrollo): %s ya no trae la aportación de la meta %d", fixGoals, fixGoalID)
	}
	return nil
}

// presAliases junta nombre de presentación y nombres de ticket → id, igual que lo hace el catálogo de la semilla real.
func (l *loader) presAliases() (map[string]int64, error) {
	b, err := os.ReadFile(filepath.Join(l.opt.Dir, "productos.json"))
	if err != nil {
		return nil, err
	}
	var f catalogFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("semilla (desarrollo): productos.json: %w", err)
	}
	out := map[string]int64{}
	for _, p := range f.Productos {
		for _, x := range p.Presentaciones {
			id, ok, err := l.find(`SELECT pr.id FROM presentation pr JOIN product pd ON pd.id = pr.product_id WHERE pd.name = ? COLLATE NOCASE AND pr.name = ?`, p.Nombre, x.Nombre)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			out[x.Nombre] = id
			for _, t := range x.Tickets {
				out[t] = id
			}
		}
	}
	return out, nil
}

// walmartFile son los precios de Walmart del carrito comparado del 10-oct (datos/semilla/data.json).
type walmartFile struct {
	Fecha          string `json:"fecha"`
	ItemsIdenticos []struct {
		Producto     string `json:"producto"`
		Cantidad     int64  `json:"cantidad"`
		WalmartTotal money  `json:"walmart_total"`
	} `json:"items_identicos"`
	Sustitutos []struct {
		Walmart      string `json:"walmart"`
		CantW        int64  `json:"cant_w"`
		WalmartTotal money  `json:"walmart_total"`
	} `json:"sustitutos"`
}

// devPrecios: precios de la compra del 10-oct ya sembrada. Bodega Aurrera sale de los renglones del ticket
// (unidad exacta, ticket_manual, ligados a su renglón); Walmart, del carrito comparado de data.json (manual).
// Un renglón cuyo importe no se divide exacto entre las piezas no tiene precio unitario y se omite (se avisa).
func (l *loader) devPrecios() error {
	rows, err := l.tx.QueryContext(l.ctx, `SELECT ei.id, ei.presentation_id, ei.unit_price_cents, e.store_id, e.expense_date
		FROM expense_item ei JOIN expense e ON e.id = ei.expense_id
		WHERE e.expense_date = '2026-10-10' AND ei.presentation_id IS NOT NULL AND e.store_id IS NOT NULL ORDER BY ei.id`)
	if err != nil {
		return err
	}
	type line struct {
		item, pres, store int64
		unit              sql.NullInt64
		date              string
	}
	var lines []line
	for rows.Next() {
		var x line
		if err := rows.Scan(&x.item, &x.pres, &x.unit, &x.store, &x.date); err != nil {
			rows.Close()
			return err
		}
		lines = append(lines, x)
	}
	rows.Close()
	if len(lines) == 0 {
		return errors.New("semilla (desarrollo): la compra del 10-oct no tiene renglones con presentación")
	}
	skipped := 0
	for _, x := range lines {
		if !x.unit.Valid {
			skipped++
			continue
		}
		if err := l.price("[desarrollo] precios Bodega (ticket 10-oct)", x.pres, x.store, x.unit.Int64, x.date, "ticket_manual", "ticket 10-oct-2026", x.item); err != nil {
			return err
		}
	}
	if skipped > 0 {
		l.rep.Notices = append(l.rep.Notices, fmt.Sprintf("Bodega: %d renglón(es) del ticket del 10-oct no tienen precio unitario exacto (importe no divisible entre las piezas) y no generaron precio.", skipped))
	}
	b, err := os.ReadFile(filepath.Join(l.opt.Dir, "data.json"))
	if err != nil {
		return err
	}
	var w walmartFile
	if err := json.Unmarshal(b, &w); err != nil {
		return fmt.Errorf("semilla (desarrollo): data.json: %w", err)
	}
	wstore, ok, err := l.find(`SELECT s.id FROM store s JOIN store c ON c.id = s.parent_id WHERE s.kind = 'branch' AND c.name = 'Walmart' ORDER BY s.id LIMIT 1`)
	if err != nil {
		return err
	}
	if !ok {
		l.rep.Notices = append(l.rep.Notices, "No hay sucursal de Walmart sembrada: no se cargaron precios de Walmart.")
		return nil
	}
	al, err := l.presAliases()
	if err != nil {
		return err
	}
	omitted := 0
	put := func(name string, qty int64, total money) error {
		id, ok := al[name]
		if !ok || qty <= 0 || int64(total)%qty != 0 {
			omitted++
			return nil
		}
		return l.price("[desarrollo] precios Walmart (carrito 10-oct)", id, wstore, int64(total)/qty, w.Fecha, "manual", "data.json carrito comparado 10-oct-2026", 0)
	}
	for _, it := range w.ItemsIdenticos {
		if err := put(it.Producto, it.Cantidad, it.WalmartTotal); err != nil {
			return err
		}
	}
	for _, s := range w.Sustitutos {
		if err := put(s.Walmart, s.CantW, s.WalmartTotal); err != nil {
			return err
		}
	}
	if omitted > 0 {
		l.rep.Notices = append(l.rep.Notices, fmt.Sprintf("Walmart: %d renglón(es) del carrito no se cargaron (sin presentación en el catálogo o sin precio unitario exacto).", omitted))
	}
	return nil
}

// price crea un precio si no existe uno igual (presentación, tienda, día, fuente, sin promoción).
func (l *loader) price(entry string, pres, store, cents int64, date, source, ref string, item int64) error {
	var itemArg any
	if item > 0 {
		itemArg = item
	}
	_, _, err := l.ensure(entry, "price", `SELECT id FROM price WHERE presentation_id = ? AND store_id = ? AND observed_date = ? AND source = ? AND is_promotional = 0`, a(pres, store, date, source),
		`INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source, is_promotional, source_ref, expense_item_id) VALUES (?, ?, ?, ?, ?, 0, ?, ?)`,
		pres, store, cents, date, source, ref, itemArg)
	return err
}

// devComidas: plan de la semana del 12 al 17-oct para Axel y Ximena con 6 recetas del recetario (sin fixture:
// dato de desarrollo). El sábado 17 sí; el domingo 18 queda sin planear.
func (l *loader) devComidas() error {
	slot, ok, err := l.find(`SELECT id FROM meal_slot WHERE code = 'lunch'`)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("semilla (desarrollo): no existe el tiempo de comida «comida»")
	}
	for _, m := range devMealPlan {
		rid, ok, err := l.find(`SELECT id FROM recipe WHERE name = ? COLLATE NOCASE`, m.Recipe)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("semilla (desarrollo): la receta «%s» no está en el recetario sembrado", m.Recipe)
		}
		for _, p := range devMealPersons {
			pid := l.ids["person:"+p]
			if _, _, err := l.ensure("[desarrollo] plan de comidas", "meal_plan_entry",
				`SELECT id FROM meal_plan_entry WHERE plan_date = ? AND meal_slot_id = ? AND person_id = ? AND recipe_id = ?`, a(m.Day, slot, pid, rid),
				`INSERT INTO meal_plan_entry (plan_date, meal_slot_id, person_id, recipe_id, servings_milli, is_eating_out, notes) VALUES (?, ?, ?, ?, 1000, 0, ?)`,
				m.Day, slot, pid, rid, devNote); err != nil {
				return err
			}
		}
	}
	return nil
}
