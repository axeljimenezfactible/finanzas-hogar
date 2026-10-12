package avisos

import (
	"context"
	"database/sql"
	"errors"

	"finanzas-hogar/internal/dominio"
)

// ErrSinFuente: la cifra todavía no tiene una función del motor o de la API que la dé (PROVISIONAL,
// docs/estado/bloqueos/S7-fuentes-avisos.md). La alerta se omite y queda anotada; nunca se inventa.
var ErrSinFuente = errors.New("fuente de cifras pendiente")

// Pago es un pago programado (scheduled_payment, calculado por motor_cred al registrar la compra).
type Pago struct {
	ID         int64
	Fecha      dominio.Date
	MontoCents int64
	Credito    string
}

// Uso es lo devengado contra lo planeado de un presupuesto en un mes. Clave: "total" o "category:<id>".
type Uso struct {
	Clave, Nombre         string
	PlanCents, UsadoCents int64
}

// Fuente da las cifras de los avisos; solo del motor o de las mismas funciones que usa la API.
type Fuente interface {
	Gastos(ctx context.Context, p dominio.Period) (dominio.ExpenseTotals, error)
	Pagos(ctx context.Context, desde, hasta dominio.Date) ([]Pago, error)
	Casillas(ctx context.Context, lunes dominio.Date) (sinArmar, total int, err error)
	Presupuesto(ctx context.Context, mes dominio.Period) ([]Uso, error)
	Despensa(ctx context.Context, hoy dominio.Date) (dominio.Date, bool, error)
	PreciosViejos(ctx context.Context, despensa dominio.Date) (int, error)
}

// FuenteBD lee de la base con las funciones de la API (dominio.SumExpenses es la de getExpenseTotals).
type FuenteBD struct{ DB *sql.DB }

// Gastos suma los gastos del periodo igual que getExpenseTotals.
func (f FuenteBD) Gastos(ctx context.Context, p dominio.Period) (dominio.ExpenseTotals, error) {
	return dominio.SumExpenses(ctx, f.DB, dominio.ExpenseFilter{Period: &p})
}

// Pagos lista los pagos programados entre dos fechas, inclusive.
func (f FuenteBD) Pagos(ctx context.Context, desde, hasta dominio.Date) ([]Pago, error) {
	rows, err := f.DB.QueryContext(ctx, `SELECT sp.id, sp.due_date, sp.amount_cents, c.name FROM scheduled_payment sp
 JOIN credit_purchase cp ON cp.id = sp.credit_purchase_id JOIN credit c ON c.id = cp.credit_id
 WHERE sp.due_date BETWEEN ? AND ? ORDER BY sp.due_date, sp.id`, desde.String(), hasta.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Pago
	for rows.Next() {
		var p Pago
		var d string
		if err := rows.Scan(&p.ID, &d, &p.MontoCents, &p.Credito); err != nil {
			return nil, err
		}
		if p.Fecha, err = dominio.ParseDate(d); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Casillas cuenta las casillas sin armar de la semana que abre lunes (K43, FIRME): pares día × tiempo
// obligatorio (desayuno, comida, cena: meal_slot.is_optional = 0) en que alguna persona del hogar no tiene
// receta ni «fuera». total = 7 × tiempos obligatorios (21).
func (f FuenteBD) Casillas(ctx context.Context, lunes dominio.Date) (int, int, error) {
	var n, slots int
	err := f.DB.QueryRowContext(ctx, `SELECT count(*) FROM meal_slot WHERE is_optional = 0`).Scan(&slots)
	if err != nil {
		return 0, 0, err
	}
	err = f.DB.QueryRowContext(ctx, `WITH d(n) AS (VALUES (0), (1), (2), (3), (4), (5), (6))
 SELECT count(*) FROM d CROSS JOIN meal_slot s WHERE s.is_optional = 0 AND EXISTS (
   SELECT 1 FROM person p WHERE p.is_active = 1 AND p.participation IN ('holder', 'co_holder') AND NOT EXISTS (
     SELECT 1 FROM meal_plan_entry e WHERE e.plan_date = date(?, '+' || d.n || ' days') AND e.meal_slot_id = s.id
       AND e.person_id = p.id AND (e.recipe_id IS NOT NULL OR e.is_eating_out = 1)))`, lunes.String()).Scan(&n)
	return n, 7 * slots, err
}

// Presupuesto: PROVISIONAL, pendiente de la función de devengo de api-presupuesto (getBudgetAlerts).
func (FuenteBD) Presupuesto(context.Context, dominio.Period) ([]Uso, error) { return nil, ErrSinFuente }

// Despensa: PROVISIONAL, el «día de despensa» no está definido en el plan (bloqueo S7-dia-despensa).
func (FuenteBD) Despensa(context.Context, dominio.Date) (dominio.Date, bool, error) {
	return dominio.Date{}, false, ErrSinFuente
}

// PreciosViejos: PROVISIONAL, pendiente de precios.StalePricesForPlan cargado desde la base (api-compras-b).
func (FuenteBD) PreciosViejos(context.Context, dominio.Date) (int, error) { return 0, ErrSinFuente }
