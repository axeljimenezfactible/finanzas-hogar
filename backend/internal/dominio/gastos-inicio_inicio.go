package dominio

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
)

// GastosExpenseIDs da los ids de los gastos que cumplen el filtro (sin paginar): sirve para pedir al motor S2
// solo los gastos filtrados (tendencia mensual, N-044).
func GastosExpenseIDs(ctx context.Context, q Querier, f ExpenseFilter) (map[int64]bool, error) {
	where, args := f.Where()
	rows, err := q.QueryContext(ctx, "SELECT e.id FROM expense e WHERE "+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]bool{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out[id] = true
	}
	return out, rows.Err()
}

// GastosIsExternalPerson dice si la persona existe con participation = external (N-026). Otra persona o una
// inexistente dan ErrNotFound.
func GastosIsExternalPerson(ctx context.Context, q Querier, id int64) error {
	var part string
	err := q.QueryRowContext(ctx, `SELECT participation FROM person WHERE id = ?`, id).Scan(&part)
	if err == sql.ErrNoRows || (err == nil && part != "external") {
		return NotFound("Esa persona no es externa o no existe. La participación solo aplica a personas externas, como Karen.")
	}
	return err
}

// gxShareSQL lee las partes de una persona en pagos programados de compras no canceladas, con lo que ella ya abonó
// a cada pago (credit_payment con payer_person_id = ella, RN-04).
const gxShareSQL = `SELECT sp.id, sp.credit_purchase_id, sp.due_date, s.amount_cents,
 (SELECT coalesce(sum(cp.amount_cents), 0) FROM credit_payment cp WHERE cp.scheduled_payment_id = sp.id AND cp.payer_person_id = s.person_id)
 FROM scheduled_payment_share s JOIN scheduled_payment sp ON sp.id = s.scheduled_payment_id
 JOIN credit_purchase pu ON pu.id = sp.credit_purchase_id
 WHERE s.person_id = ? AND s.amount_cents > 0 AND pu.status <> 'cancelled'`

// GastosPaymentShare es la parte de una persona en un pago programado.
type GastosPaymentShare struct {
	ScheduledPaymentID, PurchaseID int64
	DueDate                        string
	AmountCents, PaidCents         int64
}

// Paid dice si la parte ya se cubrió completa.
func (s GastosPaymentShare) Paid() bool { return s.PaidCents >= s.AmountCents }

// PendingCents es lo que falta de la parte (0 si ya se cubrió).
func (s GastosPaymentShare) PendingCents() int64 { return max(s.AmountCents-s.PaidCents, 0) }

// GastosPersonShares lee todas las partes de la persona en pagos programados, por fecha.
func GastosPersonShares(ctx context.Context, q Querier, personID int64) ([]GastosPaymentShare, error) {
	rows, err := q.QueryContext(ctx, gxShareSQL+` ORDER BY sp.due_date, sp.id`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []GastosPaymentShare
	for rows.Next() {
		var s GastosPaymentShare
		if err := rows.Scan(&s.ScheduledPaymentID, &s.PurchaseID, &s.DueDate, &s.AmountCents, &s.PaidCents); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// GastosParticularPurchase es una compra a crédito viva en la que la persona es la única con parte (PROVISIONAL,
// bloqueo S6-gastos-inicio §9: «compra particular»).
type GastosParticularPurchase struct {
	ID           int64
	Description  sql.NullString
	PurchaseDate string
	TotalCents   int64
}

// GastosParticularPurchases lee las compras vivas (status active) cuya única parte con monto es de la persona.
func GastosParticularPurchases(ctx context.Context, q Querier, personID int64) ([]GastosParticularPurchase, error) {
	rows, err := q.QueryContext(ctx, `SELECT pu.id,
 coalesce(pu.description, (SELECT st.name FROM expense e JOIN store st ON st.id = e.store_id WHERE e.id = pu.expense_id)),
 pu.purchase_date, pu.total_cents
 FROM credit_purchase pu
 WHERE pu.status = 'active'
 AND EXISTS (SELECT 1 FROM credit_purchase_share s WHERE s.credit_purchase_id = pu.id AND s.person_id = ? AND s.total_cents > 0)
 AND NOT EXISTS (SELECT 1 FROM credit_purchase_share s WHERE s.credit_purchase_id = pu.id AND s.person_id <> ? AND s.total_cents > 0)
 ORDER BY pu.purchase_date, pu.id`, personID, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []GastosParticularPurchase
	for rows.Next() {
		var p GastosParticularPurchase
		if err := rows.Scan(&p.ID, &p.Description, &p.PurchaseDate, &p.TotalCents); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GastosNextSavingsGoal da la meta activa (kind goal) con la fecha meta más próxima; ok = false si no hay.
func GastosNextSavingsGoal(ctx context.Context, q Querier) (SavingsGoalRecord, bool, error) {
	gs, err := ListSavingsGoals(ctx, q, true, "ORDER BY target_date, id")
	if err != nil {
		return SavingsGoalRecord{}, false, err
	}
	for _, g := range gs {
		if g.Kind == "goal" && g.TargetDate.Valid {
			return g, true, nil
		}
	}
	return SavingsGoalRecord{}, false, nil
}

// GastosSearchHit es una coincidencia de la búsqueda global (N-002).
type GastosSearchHit struct {
	ID             string
	Title          string
	Subtitle, Date sql.NullString
}

// gxSearch: por tipo, la consulta que da (id, title, sub, d) y los textos donde se busca (U1 N-002).
var gxSearch = map[string]struct {
	from string
	hay  []string
}{
	"expense": {`SELECT e.id AS id, coalesce(nullif(trim(e.description), ''), st.name, c.name) AS title, c.name AS sub, e.expense_date AS d
 FROM expense e LEFT JOIN store st ON st.id = e.store_id LEFT JOIN category c ON c.id = e.category_id`, []string{"e.description", "st.name", "c.name"}},
	"income": {`SELECT i.id AS id, i.description AS title, NULL AS sub, i.income_date AS d FROM income i WHERE i.kind = 'extra'`, []string{"i.description"}},
	"recipe": {`SELECT r.id AS id, r.name AS title, NULL AS sub, NULL AS d FROM recipe r`, []string{"r.name",
		"(SELECT group_concat(coalesce(ri.display_text, '') || ' ' || p.name, ' ') FROM recipe_ingredient ri JOIN product p ON p.id = ri.product_id WHERE ri.recipe_id = r.id)"}},
	"product": {`SELECT p.id AS id, p.name AS title, NULL AS sub, NULL AS d FROM product p`, []string{"p.name",
		"(SELECT group_concat(coalesce(pr.name, '') || ' ' || coalesce(pr.brand, ''), ' ') FROM presentation pr WHERE pr.product_id = p.id)"}},
	"store":  {`SELECT s.id AS id, s.name AS title, s.municipality AS sub, NULL AS d FROM store s`, []string{"s.name"}},
	"credit": {`SELECT k.id AS id, k.name AS title, NULL AS sub, NULL AS d FROM credit k`, []string{"k.name"}},
	"credit_purchase": {`SELECT pu.id AS id, coalesce(nullif(trim(pu.description), ''), st.name, k.name) AS title, k.name AS sub, pu.purchase_date AS d
 FROM credit_purchase pu JOIN credit k ON k.id = pu.credit_id LEFT JOIN expense e ON e.id = pu.expense_id LEFT JOIN store st ON st.id = e.store_id`,
		[]string{"pu.description", "st.name"}},
	"savings_goal": {`SELECT g.id AS id, g.name AS title, NULL AS sub, g.target_date AS d FROM savings_goal g`, []string{"g.name"}},
}

// GastosSettingsSections son las 14 secciones de Ajustes (U4, P-26.01 a P-26.14); id = clave de la pantalla
// (PROVISIONAL, bloqueo S6-gastos-inicio §12).
var GastosSettingsSections = [][2]string{{"P-26.01", "Cuenta y roles"}, {"P-26.02", "Presupuesto"}, {"P-26.03", "Ingresos"},
	{"P-26.04", "Categorías"}, {"P-26.05", "Créditos y condiciones"}, {"P-26.06", "Tiendas y precios"}, {"P-26.07", "Recetas"},
	{"P-26.08", "Reportes y frecuencias"}, {"P-26.09", "Notificaciones"}, {"P-26.10", "Apariencia"}, {"P-26.11", "Datos y respaldos"},
	{"P-26.12", "Personas y reparto"}, {"P-26.13", "Salud y metas de ahorro"}, {"P-26.14", "Inventario"}}

// GastosSearch busca text (sin mayúsculas ni acentos) en un tipo y devuelve todas las coincidencias: primero las
// que empiezan con el texto en el título, luego por fecha descendente e id (PROVISIONAL, bloqueo §12: «relevancia»).
func GastosSearch(ctx context.Context, q Querier, typ, text string) ([]GastosSearchHit, error) {
	var out []GastosSearchHit
	if typ == "settings_section" {
		needle := GastosFoldText(text)
		for _, s := range GastosSettingsSections {
			if strings.Contains(GastosFoldText(s[1]), needle) {
				out = append(out, GastosSearchHit{ID: s[0], Title: s[1]})
			}
		}
		return out, nil
	}
	d, ok := gxSearch[typ]
	if !ok {
		return nil, BadRequest("Ese tipo de búsqueda no existe.")
	}
	pat := GastosLikePattern(text)
	var w []string
	var args []any
	for _, h := range d.hay {
		w = append(w, GastosFoldSQL("coalesce("+h+", '')")+` LIKE ? ESCAPE '\'`)
		args = append(args, pat)
	}
	glue := " WHERE "
	if strings.Contains(d.from, " WHERE ") {
		glue = " AND "
	}
	query := "SELECT id, title, sub, d FROM (" + d.from + glue + "(" + strings.Join(w, " OR ") + ")) ORDER BY CASE WHEN " +
		GastosFoldSQL("coalesce(title, '')") + ` LIKE ? ESCAPE '\' THEN 0 ELSE 1 END, d IS NULL, d DESC, id DESC`
	rows, err := q.QueryContext(ctx, query, append(args, pat[1:])...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var h GastosSearchHit
		if err := rows.Scan(&id, &h.Title, &h.Subtitle, &h.Date); err != nil {
			return nil, err
		}
		h.ID = strconv.FormatInt(id, 10)
		out = append(out, h)
	}
	return out, rows.Err()
}
