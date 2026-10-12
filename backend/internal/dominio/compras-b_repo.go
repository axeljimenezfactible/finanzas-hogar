package dominio

// Repositorios de api-compras-b: precios, nutrición, inventario y lotes del agente externo (K53).
// Solo acceso a datos; los cálculos (frescura, precio por unidad base, estado de inventario) son de los motores.

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

// CBExists dice si hay un renglón con ese id. table es siempre una constante del código, nunca entrada.
func CBExists(ctx context.Context, q Querier, table string, id int64) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, `SELECT count(*) FROM `+table+` WHERE id = ?`, id).Scan(&n)
	return n > 0, err
}

// CBAuditWrite escribe audit_log con la foto «después» leída del renglón (insert o update).
func CBAuditWrite(ctx context.Context, q Querier, actor *int64, action, table string, id int64, before json.RawMessage, requestID string) error {
	after, err := RowJSON(ctx, q, table, id)
	if err != nil {
		return err
	}
	return Audit(ctx, q, actor, action, table, id, before, after, requestID)
}

// ---------- precios ----------

// CBPriceRow es un renglón de price con el contenido de su presentación y la unidad base de su producto.
type CBPriceRow struct {
	ID, PresentationID, StoreID, PriceCents int64
	ObservedDate, Source                    string
	IsPromotional                           bool
	SourceRef                               sql.NullString
	ExpenseItemID, AttachmentID             sql.NullInt64
	CreatedAt, UpdatedAt                    string
	SizeMilli                               int64
	BaseUnit                                string
}

const cbPriceSelect = `SELECT p.id, p.presentation_id, p.store_id, p.price_cents, p.observed_date, p.source, p.is_promotional,
 p.source_ref, p.expense_item_id, p.attachment_id, p.created_at, p.updated_at, pr.size_milli, pd.base_unit
 FROM price p JOIN presentation pr ON pr.id = p.presentation_id JOIN product pd ON pd.id = pr.product_id `

// CBQueryPrices lee precios con un WHERE y una cola (orden y límite) armados por la API con parámetros.
func CBQueryPrices(ctx context.Context, q Querier, where string, args []any, tail string) ([]CBPriceRow, error) {
	rows, err := q.QueryContext(ctx, cbPriceSelect+where+" "+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CBPriceRow
	for rows.Next() {
		var p CBPriceRow
		if err := rows.Scan(&p.ID, &p.PresentationID, &p.StoreID, &p.PriceCents, &p.ObservedDate, &p.Source, &p.IsPromotional,
			&p.SourceRef, &p.ExpenseItemID, &p.AttachmentID, &p.CreatedAt, &p.UpdatedAt, &p.SizeMilli, &p.BaseUnit); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// CBGetPrice lee un precio; 404 si no existe.
func CBGetPrice(ctx context.Context, q Querier, id int64) (CBPriceRow, error) {
	ps, err := CBQueryPrices(ctx, q, "WHERE p.id = ?", []any{id}, "")
	if err != nil {
		return CBPriceRow{}, err
	}
	if len(ps) == 0 {
		return CBPriceRow{}, NotFound("No encontramos ese precio. Revisa la lista de precios del producto.")
	}
	return ps[0], nil
}

// CBInsertPrice da de alta un precio y devuelve su id.
func CBInsertPrice(ctx context.Context, q Querier, p CBPriceRow) (int64, error) {
	res, err := q.ExecContext(ctx, `INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source, is_promotional, source_ref)
 VALUES (?, ?, ?, ?, ?, ?, ?)`, p.PresentationID, p.StoreID, p.PriceCents, p.ObservedDate, p.Source, p.IsPromotional, p.SourceRef)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// CBProfecoCoverage cuenta los precios PROFECO por cadena (la sucursal sube a su cadena) y municipio.
type CBProfecoCoverage struct {
	StoreID      int64
	Municipality sql.NullString
	PriceCount   int64
}

// CBProfecoStats devuelve la cobertura, el total y el alta más reciente de precios PROFECO.
func CBProfecoStats(ctx context.Context, q Querier) ([]CBProfecoCoverage, int64, sql.NullString, error) {
	rows, err := q.QueryContext(ctx, `SELECT coalesce(s.parent_id, s.id), s.municipality, count(*) FROM price p JOIN store s ON s.id = p.store_id
 WHERE p.source = 'profeco' GROUP BY 1, 2 ORDER BY 1, 2`)
	if err != nil {
		return nil, 0, sql.NullString{}, err
	}
	defer rows.Close()
	var out []CBProfecoCoverage
	var total int64
	for rows.Next() {
		var c CBProfecoCoverage
		if err := rows.Scan(&c.StoreID, &c.Municipality, &c.PriceCount); err != nil {
			return nil, 0, sql.NullString{}, err
		}
		total += c.PriceCount
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, sql.NullString{}, err
	}
	var last sql.NullString
	err = q.QueryRowContext(ctx, `SELECT max(created_at) FROM price WHERE source = 'profeco'`).Scan(&last)
	return out, total, last, err
}

// ---------- nutrición ----------

// CBNutritionRow es un renglón de nutrition_fact.
type CBNutritionRow struct {
	ID, ProductID              int64
	PresentationID             sql.NullInt64
	KcalMilli                  int64
	Protein, Fat, Carbs, Fiber sql.NullInt64
	Source                     string
	SourceRef, SourceLicense   sql.NullString
	IsPreferred                bool
	CreatedAt, UpdatedAt       string
}

const cbNutritionSelect = `SELECT id, product_id, presentation_id, kcal_milli_per_100g, protein_mg_per_100g, fat_mg_per_100g, carbs_mg_per_100g,
 fiber_mg_per_100g, source, source_ref, source_license, is_preferred, created_at, updated_at FROM nutrition_fact `

// CBQueryNutrition lee datos nutricionales con un WHERE y una cola armados por la API.
func CBQueryNutrition(ctx context.Context, q Querier, where string, args []any, tail string) ([]CBNutritionRow, error) {
	rows, err := q.QueryContext(ctx, cbNutritionSelect+where+" "+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CBNutritionRow
	for rows.Next() {
		var n CBNutritionRow
		if err := rows.Scan(&n.ID, &n.ProductID, &n.PresentationID, &n.KcalMilli, &n.Protein, &n.Fat, &n.Carbs, &n.Fiber,
			&n.Source, &n.SourceRef, &n.SourceLicense, &n.IsPreferred, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// ---------- inventario ----------

// CBInventoryRow es un renglón de inventory_item con el nombre de su producto o receta.
type CBInventoryRow struct {
	ID                                                   int64
	ProductID, RecipeID                                  sql.NullInt64
	TrackingMode                                         string
	QuantityMilli                                        int64
	CountedDate                                          sql.NullString
	RateQuantityMilli, RatePeriodDays, ReorderLevelMilli sql.NullInt64
	ExpiresDate, Notes                                   sql.NullString
	CreatedAt, UpdatedAt, Name                           string
}

const cbInventorySelect = `SELECT i.id, i.product_id, i.recipe_id, i.tracking_mode, i.quantity_milli, i.counted_date, i.rate_quantity_milli,
 i.rate_period_days, i.reorder_level_milli, i.expires_date, i.notes, i.created_at, i.updated_at, coalesce(pd.name, r.name, '')
 FROM inventory_item i LEFT JOIN product pd ON pd.id = i.product_id LEFT JOIN recipe r ON r.id = i.recipe_id `

// CBQueryInventory lee renglones de inventario con un WHERE armado por la API.
func CBQueryInventory(ctx context.Context, q Querier, where string, args ...any) ([]CBInventoryRow, error) {
	rows, err := q.QueryContext(ctx, cbInventorySelect+where+" ORDER BY i.id", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CBInventoryRow
	for rows.Next() {
		var i CBInventoryRow
		if err := rows.Scan(&i.ID, &i.ProductID, &i.RecipeID, &i.TrackingMode, &i.QuantityMilli, &i.CountedDate, &i.RateQuantityMilli,
			&i.RatePeriodDays, &i.ReorderLevelMilli, &i.ExpiresDate, &i.Notes, &i.CreatedAt, &i.UpdatedAt, &i.Name); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

// CBGetInventory lee un renglón de inventario; 404 si no existe.
func CBGetInventory(ctx context.Context, q Querier, id int64) (CBInventoryRow, error) {
	is, err := CBQueryInventory(ctx, q, "WHERE i.id = ?", id)
	if err != nil {
		return CBInventoryRow{}, err
	}
	if len(is) == 0 {
		return CBInventoryRow{}, NotFound("No encontramos ese artículo de inventario.")
	}
	return is[0], nil
}

// CBMovementRow es un renglón de inventory_movement con el gasto del renglón de ticket de origen.
type CBMovementRow struct {
	ID, ItemID                                int64
	Date, Kind                                string
	Reason                                    sql.NullString
	Delta                                     int64
	ExpenseItemID, MealPlanEntryID, ExpenseID sql.NullInt64
	Notes                                     sql.NullString
	CreatedAt, UpdatedAt                      string
}

// CBQueryMovements lee los movimientos de un renglón con la cola (orden y límite) armada por la API.
func CBQueryMovements(ctx context.Context, q Querier, itemID int64, tail string) ([]CBMovementRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT m.id, m.inventory_item_id, m.movement_date, m.kind, m.reason, m.quantity_delta_milli,
 m.expense_item_id, m.meal_plan_entry_id, ei.expense_id, m.notes, m.created_at, m.updated_at
 FROM inventory_movement m LEFT JOIN expense_item ei ON ei.id = m.expense_item_id WHERE m.inventory_item_id = ? `+tail, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CBMovementRow
	for rows.Next() {
		var m CBMovementRow
		if err := rows.Scan(&m.ID, &m.ItemID, &m.Date, &m.Kind, &m.Reason, &m.Delta, &m.ExpenseItemID, &m.MealPlanEntryID,
			&m.ExpenseID, &m.Notes, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// CBInsertMovement escribe un inventory_movement (misma transacción que el cambio de existencia, invariante 12 de S1).
func CBInsertMovement(ctx context.Context, q Querier, itemID int64, date, kind string, reason *string, delta int64, notes *string) (int64, error) {
	res, err := q.ExecContext(ctx, `INSERT INTO inventory_movement (inventory_item_id, movement_date, kind, reason, quantity_delta_milli, notes)
 VALUES (?, ?, ?, ?, ?, ?)`, itemID, date, kind, reason, delta, notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ---------- lotes del agente externo (K53) ----------

// Estados de un renglón de lote (ExternalImportResult).
const (
	CBRowCreated   = "created"
	CBRowUpdated   = "updated"
	CBRowUnchanged = "unchanged"
)

// CBCatalogPresentation es una presentación de un renglón de catálogo.
type CBCatalogPresentation struct {
	Name      string
	SizeMilli int64
	Brand     *string
	IsBulk    bool
	Barcode   *string
}

// CBCatalogItem es un producto de un lote de catálogo.
type CBCatalogItem struct {
	Name              string
	BaseUnit          string
	GramsPerUnitMilli *int64
	ProductGroupID    *int64
	Presentations     []CBCatalogPresentation
}

// CBUpsertCatalogItem da de alta o actualiza un producto y sus presentaciones por nombre (llave natural, sin
// distinguir mayúsculas). Devuelve el estado del renglón, el id del producto y los de sus presentaciones en orden.
// Un *Error de validación (Invalid/Conflict) indica rechazo del renglón.
func CBUpsertCatalogItem(ctx context.Context, tx *sql.Tx, it CBCatalogItem, actor *int64, requestID string) (string, int64, []int64, error) {
	status := CBRowUnchanged
	var id int64
	var unit string
	var grams, group sql.NullInt64
	err := tx.QueryRowContext(ctx, `SELECT id, base_unit, grams_per_unit_milli, product_group_id FROM product WHERE name = ?`, it.Name).Scan(&id, &unit, &grams, &group)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		res, err := tx.ExecContext(ctx, `INSERT INTO product (name, base_unit, grams_per_unit_milli, product_group_id) VALUES (?, ?, ?, ?)`,
			it.Name, it.BaseUnit, it.GramsPerUnitMilli, it.ProductGroupID)
		if err != nil {
			return "", 0, nil, err
		}
		if id, err = res.LastInsertId(); err != nil {
			return "", 0, nil, err
		}
		status = CBRowCreated
		if err := CBAuditWrite(ctx, tx, actor, "insert", "product", id, nil, requestID); err != nil {
			return "", 0, nil, err
		}
	case err != nil:
		return "", 0, nil, err
	default:
		if unit != it.BaseUnit {
			return "", 0, nil, Invalid(FieldIssue{"base_unit", "conflict", "El producto ya existe con otra unidad base; la unidad base no cambia por lote."})
		}
		ng, ngr := cbNullInt(it.GramsPerUnitMilli), cbNullInt(it.ProductGroupID)
		if ng != grams || ngr != group {
			before, err := RowJSON(ctx, tx, "product", id)
			if err != nil {
				return "", 0, nil, err
			}
			if _, err := tx.ExecContext(ctx, `UPDATE product SET grams_per_unit_milli = ?, product_group_id = ? WHERE id = ?`, ng, ngr, id); err != nil {
				return "", 0, nil, err
			}
			status = CBRowUpdated
			if err := CBAuditWrite(ctx, tx, actor, "update", "product", id, before, requestID); err != nil {
				return "", 0, nil, err
			}
		}
	}
	ids := make([]int64, 0, len(it.Presentations))
	for i, p := range it.Presentations {
		pid, st, err := cbUpsertPresentation(ctx, tx, id, p, actor, requestID)
		if err != nil {
			if e, ok := AsError(err); ok {
				for j := range e.Fields {
					e.Fields[j].Field = "presentations[" + cbItoa(i) + "]." + e.Fields[j].Field
				}
			}
			return "", 0, nil, err
		}
		if st != CBRowUnchanged && status == CBRowUnchanged {
			status = CBRowUpdated
		}
		ids = append(ids, pid)
	}
	return status, id, ids, nil
}

func cbUpsertPresentation(ctx context.Context, tx *sql.Tx, productID int64, p CBCatalogPresentation, actor *int64, requestID string) (int64, string, error) {
	if p.Barcode != nil {
		var other int64
		err := tx.QueryRowContext(ctx, `SELECT id FROM presentation WHERE barcode = ? AND NOT (product_id = ? AND name = ?)`, *p.Barcode, productID, p.Name).Scan(&other)
		if err == nil {
			return 0, "", Invalid(FieldIssue{"barcode", "conflict", "Ese código de barras ya pertenece a otra presentación."})
		} else if !errors.Is(err, sql.ErrNoRows) {
			return 0, "", err
		}
	}
	var id, size int64
	var brand, barcode sql.NullString
	var bulk bool
	err := tx.QueryRowContext(ctx, `SELECT id, size_milli, brand, is_bulk, barcode FROM presentation WHERE product_id = ? AND name = ?`, productID, p.Name).
		Scan(&id, &size, &brand, &bulk, &barcode)
	if errors.Is(err, sql.ErrNoRows) {
		res, err := tx.ExecContext(ctx, `INSERT INTO presentation (product_id, name, brand, size_milli, is_bulk, barcode) VALUES (?, ?, ?, ?, ?, ?)`,
			productID, p.Name, p.Brand, p.SizeMilli, p.IsBulk, p.Barcode)
		if err != nil {
			return 0, "", err
		}
		if id, err = res.LastInsertId(); err != nil {
			return 0, "", err
		}
		return id, CBRowCreated, CBAuditWrite(ctx, tx, actor, "insert", "presentation", id, nil, requestID)
	} else if err != nil {
		return 0, "", err
	}
	nb, nc := cbNullStr(p.Brand), cbNullStr(p.Barcode)
	if size == p.SizeMilli && brand == nb && bulk == p.IsBulk && barcode == nc {
		return id, CBRowUnchanged, nil
	}
	before, err := RowJSON(ctx, tx, "presentation", id)
	if err != nil {
		return 0, "", err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE presentation SET size_milli = ?, brand = ?, is_bulk = ?, barcode = ? WHERE id = ?`, p.SizeMilli, nb, p.IsBulk, nc, id); err != nil {
		return 0, "", err
	}
	return id, CBRowUpdated, CBAuditWrite(ctx, tx, actor, "update", "presentation", id, before, requestID)
}

// CBUpsertSourcedPrice guarda un precio por su llave UNIQUE (presentación, tienda, fecha, fuente y promoción):
// lo crea, lo actualiza si cambió el importe o la referencia, o lo deja igual.
func CBUpsertSourcedPrice(ctx context.Context, tx *sql.Tx, p CBPriceRow, actor *int64, requestID string) (string, int64, error) {
	var id, cents int64
	var ref sql.NullString
	err := tx.QueryRowContext(ctx, `SELECT id, price_cents, source_ref FROM price WHERE presentation_id = ? AND store_id = ? AND observed_date = ?
 AND source = ? AND is_promotional = ?`, p.PresentationID, p.StoreID, p.ObservedDate, p.Source, p.IsPromotional).Scan(&id, &cents, &ref)
	if errors.Is(err, sql.ErrNoRows) {
		if id, err = CBInsertPrice(ctx, tx, p); err != nil {
			return "", 0, err
		}
		return CBRowCreated, id, CBAuditWrite(ctx, tx, actor, "insert", "price", id, nil, requestID)
	} else if err != nil {
		return "", 0, err
	}
	if cents == p.PriceCents && ref == p.SourceRef {
		return CBRowUnchanged, id, nil
	}
	before, err := RowJSON(ctx, tx, "price", id)
	if err != nil {
		return "", 0, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE price SET price_cents = ?, source_ref = ? WHERE id = ?`, p.PriceCents, p.SourceRef, id); err != nil {
		return "", 0, err
	}
	return CBRowUpdated, id, CBAuditWrite(ctx, tx, actor, "update", "price", id, before, requestID)
}

// CBDraftRecipe es una receta de un lote del agente externo.
type CBDraftRecipe struct {
	Name                             string
	Description, Instructions, Notes *string
	BaseServings                     int64
	PrepMinutes, CookMinutes         *int64
	MealSlotIDs                      []int64
	Ingredients                      []CBDraftIngredient
}

// CBDraftIngredient es un ingrediente de CBDraftRecipe.
type CBDraftIngredient struct {
	ProductID, QuantityMilli int64
	PresentationID           *int64
	DisplayText              *string
	IsOptional               bool
	SortOrder                int64
}

// CBInsertDraftRecipe da de alta una receta del agente externo con sus tiempos e ingredientes.
// PROVISIONAL (bloqueo S5-agente-externo, preguntas c y f; opción f2): el esquema 1.4 no guarda el borrador ni el
// origen «external_agent»: entra con origin = household, sin autor e is_active = 0 hasta que una persona la confirme.
// Una receta con el mismo nombre (sin distinguir mayúsculas) no se sobrescribe: 409 already_exists en su renglón.
func CBInsertDraftRecipe(ctx context.Context, tx *sql.Tx, r CBDraftRecipe, actor *int64, requestID string) (int64, error) {
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM recipe WHERE name = ?`, r.Name).Scan(&n); err != nil {
		return 0, err
	}
	if n > 0 {
		return 0, Conflict("already_exists", "Ya hay una receta con ese nombre; no se sobrescribe.")
	}
	res, err := tx.ExecContext(ctx, `INSERT INTO recipe (name, description, instructions, base_servings, prep_minutes, cook_minutes, origin, is_active, notes)
 VALUES (?, ?, ?, ?, ?, ?, 'household', 0, ?)`, r.Name, r.Description, r.Instructions, r.BaseServings, r.PrepMinutes, r.CookMinutes, r.Notes)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, s := range r.MealSlotIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO recipe_meal_slot (recipe_id, meal_slot_id) VALUES (?, ?)`, id, s); err != nil {
			return 0, err
		}
	}
	for _, g := range r.Ingredients {
		if _, err := tx.ExecContext(ctx, `INSERT INTO recipe_ingredient (recipe_id, product_id, presentation_id, quantity_milli, display_text, is_optional, sort_order)
 VALUES (?, ?, ?, ?, ?, ?, ?)`, id, g.ProductID, g.PresentationID, g.QuantityMilli, g.DisplayText, g.IsOptional, g.SortOrder); err != nil {
			return 0, err
		}
	}
	return id, CBAuditWrite(ctx, tx, actor, "insert", "recipe", id, nil, requestID)
}

// CBPresentationOf devuelve el producto de una presentación (sql.ErrNoRows si no existe).
func CBPresentationOf(ctx context.Context, q Querier, id int64) (int64, error) {
	var pid int64
	err := q.QueryRowContext(ctx, `SELECT product_id FROM presentation WHERE id = ?`, id).Scan(&pid)
	return pid, err
}

func cbNullInt(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

func cbNullStr(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func cbItoa(i int) string {
	b, _ := json.Marshal(i)
	return strings.TrimSpace(string(b))
}

// CBStrPtr devuelve un puntero a s (para campos opcionales de los repositorios del área).
func CBStrPtr(s string) *string { return &s }
