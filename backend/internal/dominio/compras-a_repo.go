package dominio

import (
	"context"
	"database/sql"
	"strings"
)

// Repositorios del área «Compras y tiendas» (agente api-compras-a): tiendas, grupos, productos, presentaciones,
// precios (solo lectura) y lista de compras, más lo que leen los motores S4 y S5 (recetas, plan, inventario y
// compras con renglones). Todo va bajo ShopA para no chocar con los nombres de otras áreas.

type shopRowScanner interface{ Scan(dest ...any) error }

type shopRepo struct{}

// ShopA agrupa los repositorios del área de compras (parte a).
var ShopA shopRepo

// ShopStore es un renglón de store con sus grupos (store_product_group).
type ShopStore struct {
	ID                                               int64
	Kind, Name                                       string
	ParentID, VisitWeekday, MembershipAnnualFeeCents sql.NullInt64
	ProfecoName, WebsiteURL, Address                 sql.NullString
	Municipality, State, Notes                       sql.NullString
	AssistedReadingEnabled, RequiresMembership       bool
	HasMembership, IsActive                          bool
	CreatedAt, UpdatedAt                             string
	GroupIDs                                         []int64
}

// ShopGroup es un renglón de product_group.
type ShopGroup struct {
	ID                    int64
	Name, PurchaseChannel string
	SortOrder             int64
	IsActive              bool
	CreatedAt, UpdatedAt  string
}

// ShopProduct es un renglón de product.
type ShopProduct struct {
	ID                         int64
	Name, BaseUnit             string
	GroupID, GramsPerUnitMilli sql.NullInt64
	IsActive                   bool
	Notes                      sql.NullString
	CreatedAt, UpdatedAt       string
}

// ShopPresentation es un renglón de presentation.
type ShopPresentation struct {
	ID, ProductID, SizeMilli int64
	Name                     string
	Brand, Barcode, Notes    sql.NullString
	IsBulk, IsActive         bool
	CreatedAt, UpdatedAt     string
}

// ShopPrice es un renglón de price (solo lectura en esta área).
type ShopPrice struct {
	ID, PresentationID, StoreID, PriceCents int64
	ObservedDate, Source, CreatedAt         string
	IsPromotional                           bool
	SourceRef                               sql.NullString
}

// ShopNutrition es un renglón de nutrition_fact.
type ShopNutrition struct {
	ID, ProductID, KcalMilliPer100g          int64
	PresentationID, Protein, Fat, Carbs, Fib sql.NullInt64
	Source                                   string
	SourceRef, SourceLicense                 sql.NullString
	IsPreferred                              bool
	CreatedAt, UpdatedAt                     string
}

// ShopList es un renglón de shopping_list.
type ShopList struct {
	ID                   int64
	StartDate, EndDate   string
	GeneratedAt          sql.NullString
	CreatedAt, UpdatedAt string
}

// ShopListItem es un renglón de shopping_list_item.
type ShopListItem struct {
	ID, ListID                                                     int64
	Origin                                                         string
	ProductID, PresentationID, NeededMilli, QuantityMilli, StoreID sql.NullInt64
	FreeText, FreeTextUnit, BoughtAt, Notes                        sql.NullString
	IsBought                                                       bool
	CreatedAt, UpdatedAt                                           string
}

// ShopIngredient es un renglón de recipe_ingredient.
type ShopIngredient struct {
	RecipeID, ProductID, QuantityMilli, SortOrder int64
	PresentationID                                sql.NullInt64
	DisplayText                                   sql.NullString
	IsOptional                                    bool
}

// ShopRecipe es una receta con sus ingredientes (lo que lee S4 para la lista).
type ShopRecipe struct {
	ID, BaseServings int64
	Name, Origin     string
	IsActive         bool
	Ingredients      []ShopIngredient
}

// ShopPlanEntry es un renglón de meal_plan_entry con el código de su tiempo de comida.
type ShopPlanEntry struct {
	Date, Slot              string
	PersonID                int64
	RecipeID, ServingsMilli sql.NullInt64
	IsEatingOut             bool
}

// ShopPerson es lo que S4 necesita de person.
type ShopPerson struct {
	ID                  int64
	Code, Participation string
}

// ShopInventory es un renglón de inventory_item.
type ShopInventory struct {
	ID, QuantityMilli                                                    int64
	Mode                                                                 string
	ProductID, RecipeID, RateQuantityMilli, RatePeriodDays, ReorderLevel sql.NullInt64
	CountedDate, ExpiresDate                                             sql.NullString
}

// ShopExpense es una compra con renglones (canasta, N-160).
type ShopExpense struct {
	ID, AmountCents, ItemCount int64
	Date                       string
	StoreID                    sql.NullInt64
}

// ShopExpenseLine es un renglón de expense_item.
type ShopExpenseLine struct {
	PresentationID, UnitPriceCents sql.NullInt64
	QuantityMilli, LineTotalCents  int64
}

// Insert inserta un renglón (tabla y columnas fijas del código, nunca del usuario) y devuelve su id.
func (shopRepo) Insert(ctx context.Context, q Querier, table string, cols []string, vals ...any) (int64, error) {
	marks := strings.TrimSuffix(strings.Repeat("?, ", len(cols)), ", ")
	res, err := q.ExecContext(ctx, "INSERT INTO "+table+" ("+strings.Join(cols, ", ")+") VALUES ("+marks+")", vals...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update cambia columnas de un renglón por id; sin renglón devuelve ErrNotFound.
func (shopRepo) Update(ctx context.Context, q Querier, table string, id int64, cols []string, vals ...any) error {
	set := make([]string, len(cols))
	for i, c := range cols {
		set[i] = c + " = ?"
	}
	res, err := q.ExecContext(ctx, "UPDATE "+table+" SET "+strings.Join(set, ", ")+" WHERE id = ?", append(vals, id)...)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete borra un renglón por id.
func (shopRepo) Delete(ctx context.Context, q Querier, table string, id int64) error {
	_, err := q.ExecContext(ctx, "DELETE FROM "+table+" WHERE id = ?", id)
	return err
}

func shopCollect[T any](rows *sql.Rows, err error, scan func(shopRowScanner) (T, error)) ([]T, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []T
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

const shopStoreCols = "id, kind, parent_id, name, profeco_name, website_url, assisted_reading_enabled, address, municipality, state, visit_weekday, requires_membership, has_membership, membership_annual_fee_cents, is_active, notes, created_at, updated_at"

func shopScanStore(s shopRowScanner) (ShopStore, error) {
	var r ShopStore
	err := s.Scan(&r.ID, &r.Kind, &r.ParentID, &r.Name, &r.ProfecoName, &r.WebsiteURL, &r.AssistedReadingEnabled, &r.Address,
		&r.Municipality, &r.State, &r.VisitWeekday, &r.RequiresMembership, &r.HasMembership, &r.MembershipAnnualFeeCents,
		&r.IsActive, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}

func (shopRepo) storeGroups(ctx context.Context, q Querier) (map[int64][]int64, error) {
	rows, err := q.QueryContext(ctx, "SELECT store_id, product_group_id FROM store_product_group ORDER BY store_id, product_group_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64][]int64{}
	for rows.Next() {
		var s, g int64
		if err := rows.Scan(&s, &g); err != nil {
			return nil, err
		}
		out[s] = append(out[s], g)
	}
	return out, rows.Err()
}

// Stores devuelve las tiendas (kind vacío = todas) con el ORDER BY que arma la API, cada una con sus grupos.
func (r shopRepo) Stores(ctx context.Context, q Querier, kind, orderBy string) ([]ShopStore, error) {
	query, args := "SELECT "+shopStoreCols+" FROM store", []any{}
	if kind != "" {
		query += " WHERE kind = ?"
		args = append(args, kind)
	}
	rows, err := q.QueryContext(ctx, query+" "+orderBy, args...)
	out, err := shopCollect(rows, err, shopScanStore)
	if err != nil {
		return nil, err
	}
	gs, err := r.storeGroups(ctx, q)
	for i := range out {
		out[i].GroupIDs = gs[out[i].ID]
	}
	return out, err
}

// Store lee una tienda con sus grupos (sql.ErrNoRows si no existe).
func (r shopRepo) Store(ctx context.Context, q Querier, id int64) (ShopStore, error) {
	s, err := shopScanStore(q.QueryRowContext(ctx, "SELECT "+shopStoreCols+" FROM store WHERE id = ?", id))
	if err != nil {
		return s, err
	}
	gs, err := r.storeGroups(ctx, q)
	s.GroupIDs = gs[id]
	return s, err
}

// SaveStore escribe todas las columnas editables de la tienda.
func (r shopRepo) SaveStore(ctx context.Context, q Querier, s ShopStore) error {
	return r.Update(ctx, q, "store", s.ID, []string{"kind", "parent_id", "name", "profeco_name", "website_url", "assisted_reading_enabled",
		"address", "municipality", "state", "visit_weekday", "requires_membership", "has_membership", "membership_annual_fee_cents", "is_active", "notes"},
		s.Kind, s.ParentID, s.Name, s.ProfecoName, s.WebsiteURL, s.AssistedReadingEnabled, s.Address, s.Municipality, s.State,
		s.VisitWeekday, s.RequiresMembership, s.HasMembership, s.MembershipAnnualFeeCents, s.IsActive, s.Notes)
}

// SetStoreGroups reemplaza «qué se compra ahí» de una tienda.
func (shopRepo) SetStoreGroups(ctx context.Context, q Querier, storeID int64, groupIDs []int64) error {
	if _, err := q.ExecContext(ctx, "DELETE FROM store_product_group WHERE store_id = ?", storeID); err != nil {
		return err
	}
	for _, g := range groupIDs {
		if _, err := q.ExecContext(ctx, "INSERT INTO store_product_group (store_id, product_group_id) VALUES (?, ?)", storeID, g); err != nil {
			return err
		}
	}
	return nil
}

// Groups devuelve los grupos de producto con el ORDER BY dado.
func (shopRepo) Groups(ctx context.Context, q Querier, orderBy string) ([]ShopGroup, error) {
	rows, err := q.QueryContext(ctx, "SELECT id, name, purchase_channel, sort_order, is_active, created_at, updated_at FROM product_group "+orderBy)
	return shopCollect(rows, err, func(s shopRowScanner) (ShopGroup, error) {
		var g ShopGroup
		err := s.Scan(&g.ID, &g.Name, &g.PurchaseChannel, &g.SortOrder, &g.IsActive, &g.CreatedAt, &g.UpdatedAt)
		return g, err
	})
}

const shopProductCols = "id, name, product_group_id, base_unit, grams_per_unit_milli, is_active, notes, created_at, updated_at"

func shopScanProduct(s shopRowScanner) (ShopProduct, error) {
	var p ShopProduct
	err := s.Scan(&p.ID, &p.Name, &p.GroupID, &p.BaseUnit, &p.GramsPerUnitMilli, &p.IsActive, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// Products devuelve los productos (groupID 0 = todos) con el ORDER BY dado.
func (shopRepo) Products(ctx context.Context, q Querier, groupID int64, orderBy string) ([]ShopProduct, error) {
	query, args := "SELECT "+shopProductCols+" FROM product", []any{}
	if groupID > 0 {
		query += " WHERE product_group_id = ?"
		args = append(args, groupID)
	}
	rows, err := q.QueryContext(ctx, query+" "+orderBy, args...)
	return shopCollect(rows, err, shopScanProduct)
}

// Product lee un producto (sql.ErrNoRows si no existe).
func (shopRepo) Product(ctx context.Context, q Querier, id int64) (ShopProduct, error) {
	return shopScanProduct(q.QueryRowContext(ctx, "SELECT "+shopProductCols+" FROM product WHERE id = ?", id))
}

// SaveProduct escribe las columnas editables del producto.
func (r shopRepo) SaveProduct(ctx context.Context, q Querier, p ShopProduct) error {
	return r.Update(ctx, q, "product", p.ID, []string{"name", "product_group_id", "base_unit", "grams_per_unit_milli", "is_active", "notes"},
		p.Name, p.GroupID, p.BaseUnit, p.GramsPerUnitMilli, p.IsActive, p.Notes)
}

const shopPresCols = "id, product_id, name, brand, size_milli, is_bulk, barcode, is_active, notes, created_at, updated_at"

func shopScanPres(s shopRowScanner) (ShopPresentation, error) {
	var p ShopPresentation
	err := s.Scan(&p.ID, &p.ProductID, &p.Name, &p.Brand, &p.SizeMilli, &p.IsBulk, &p.Barcode, &p.IsActive, &p.Notes, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// Presentations devuelve las presentaciones de un producto (0 = todas), por id.
func (shopRepo) Presentations(ctx context.Context, q Querier, productID int64) ([]ShopPresentation, error) {
	query, args := "SELECT "+shopPresCols+" FROM presentation", []any{}
	if productID > 0 {
		query += " WHERE product_id = ?"
		args = append(args, productID)
	}
	rows, err := q.QueryContext(ctx, query+" ORDER BY id", args...)
	return shopCollect(rows, err, shopScanPres)
}

// Presentation lee una presentación (sql.ErrNoRows si no existe).
func (shopRepo) Presentation(ctx context.Context, q Querier, id int64) (ShopPresentation, error) {
	return shopScanPres(q.QueryRowContext(ctx, "SELECT "+shopPresCols+" FROM presentation WHERE id = ?", id))
}

// SavePresentation escribe las columnas editables de la presentación.
func (r shopRepo) SavePresentation(ctx context.Context, q Querier, p ShopPresentation) error {
	return r.Update(ctx, q, "presentation", p.ID, []string{"name", "brand", "size_milli", "is_bulk", "barcode", "is_active", "notes"},
		p.Name, p.Brand, p.SizeMilli, p.IsBulk, p.Barcode, p.IsActive, p.Notes)
}

// Prices devuelve todos los precios por id (la entrada de S5).
func (shopRepo) Prices(ctx context.Context, q Querier) ([]ShopPrice, error) {
	rows, err := q.QueryContext(ctx, "SELECT id, presentation_id, store_id, price_cents, observed_date, source, is_promotional, source_ref, created_at FROM price ORDER BY id")
	return shopCollect(rows, err, func(s shopRowScanner) (ShopPrice, error) {
		var p ShopPrice
		err := s.Scan(&p.ID, &p.PresentationID, &p.StoreID, &p.PriceCents, &p.ObservedDate, &p.Source, &p.IsPromotional, &p.SourceRef, &p.CreatedAt)
		return p, err
	})
}

// Nutrition devuelve los datos nutricionales de un producto, por id.
func (shopRepo) Nutrition(ctx context.Context, q Querier, productID int64) ([]ShopNutrition, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, product_id, presentation_id, kcal_milli_per_100g, protein_mg_per_100g, fat_mg_per_100g,
		carbs_mg_per_100g, fiber_mg_per_100g, source, source_ref, source_license, is_preferred, created_at, updated_at
		FROM nutrition_fact WHERE product_id = ? ORDER BY id`, productID)
	return shopCollect(rows, err, func(s shopRowScanner) (ShopNutrition, error) {
		var n ShopNutrition
		err := s.Scan(&n.ID, &n.ProductID, &n.PresentationID, &n.KcalMilliPer100g, &n.Protein, &n.Fat, &n.Carbs, &n.Fib, &n.Source,
			&n.SourceRef, &n.SourceLicense, &n.IsPreferred, &n.CreatedAt, &n.UpdatedAt)
		return n, err
	})
}

// ListByStart lee la lista de la semana que empieza en start (sql.ErrNoRows si no existe).
func (shopRepo) ListByStart(ctx context.Context, q Querier, start string) (ShopList, error) {
	var l ShopList
	err := q.QueryRowContext(ctx, "SELECT id, start_date, end_date, generated_at, created_at, updated_at FROM shopping_list WHERE start_date = ?", start).
		Scan(&l.ID, &l.StartDate, &l.EndDate, &l.GeneratedAt, &l.CreatedAt, &l.UpdatedAt)
	return l, err
}

// EnsureList crea la lista de la semana si falta y devuelve su id y si se creó.
func (r shopRepo) EnsureList(ctx context.Context, q Querier, start, end string) (int64, bool, error) {
	l, err := r.ListByStart(ctx, q, start)
	if err == nil {
		return l.ID, false, nil
	}
	if err != sql.ErrNoRows {
		return 0, false, err
	}
	id, err := r.Insert(ctx, q, "shopping_list", []string{"start_date", "end_date"}, start, end)
	return id, err == nil, err
}

const shopItemCols = "id, shopping_list_id, origin, product_id, free_text, free_text_unit, presentation_id, needed_milli, quantity_milli, store_id, is_bought, bought_at, notes, created_at, updated_at"

func shopScanItem(s shopRowScanner) (ShopListItem, error) {
	var it ShopListItem
	err := s.Scan(&it.ID, &it.ListID, &it.Origin, &it.ProductID, &it.FreeText, &it.FreeTextUnit, &it.PresentationID, &it.NeededMilli,
		&it.QuantityMilli, &it.StoreID, &it.IsBought, &it.BoughtAt, &it.Notes, &it.CreatedAt, &it.UpdatedAt)
	return it, err
}

// Items devuelve los artículos de una lista, por id.
func (shopRepo) Items(ctx context.Context, q Querier, listID int64) ([]ShopListItem, error) {
	rows, err := q.QueryContext(ctx, "SELECT "+shopItemCols+" FROM shopping_list_item WHERE shopping_list_id = ? ORDER BY id", listID)
	return shopCollect(rows, err, shopScanItem)
}

// Item lee un artículo (sql.ErrNoRows si no existe).
func (shopRepo) Item(ctx context.Context, q Querier, id int64) (ShopListItem, error) {
	return shopScanItem(q.QueryRowContext(ctx, "SELECT "+shopItemCols+" FROM shopping_list_item WHERE id = ?", id))
}

// SaveItem escribe las columnas editables del artículo.
func (r shopRepo) SaveItem(ctx context.Context, q Querier, it ShopListItem) error {
	return r.Update(ctx, q, "shopping_list_item", it.ID, []string{"presentation_id", "needed_milli", "quantity_milli", "store_id", "is_bought", "bought_at", "notes"},
		it.PresentationID, it.NeededMilli, it.QuantityMilli, it.StoreID, it.IsBought, it.BoughtAt, it.Notes)
}

// Recipes devuelve las recetas con sus ingredientes (orden por id y sort_order).
func (shopRepo) Recipes(ctx context.Context, q Querier) ([]ShopRecipe, error) {
	rows, err := q.QueryContext(ctx, "SELECT id, name, base_servings, origin, is_active FROM recipe ORDER BY id")
	recs, err := shopCollect(rows, err, func(s shopRowScanner) (ShopRecipe, error) {
		var r ShopRecipe
		err := s.Scan(&r.ID, &r.Name, &r.BaseServings, &r.Origin, &r.IsActive)
		return r, err
	})
	if err != nil {
		return nil, err
	}
	rows, err = q.QueryContext(ctx, `SELECT recipe_id, product_id, presentation_id, quantity_milli, display_text, is_optional, sort_order
		FROM recipe_ingredient ORDER BY recipe_id, sort_order, id`)
	ings, err := shopCollect(rows, err, func(s shopRowScanner) (ShopIngredient, error) {
		var g ShopIngredient
		err := s.Scan(&g.RecipeID, &g.ProductID, &g.PresentationID, &g.QuantityMilli, &g.DisplayText, &g.IsOptional, &g.SortOrder)
		return g, err
	})
	idx := map[int64]int{}
	for i := range recs {
		idx[recs[i].ID] = i
	}
	for _, g := range ings {
		recs[idx[g.RecipeID]].Ingredients = append(recs[idx[g.RecipeID]].Ingredients, g)
	}
	return recs, err
}

// PlanEntries devuelve el plan de comidas entre from y to (inclusive) con el código del tiempo de comida.
func (shopRepo) PlanEntries(ctx context.Context, q Querier, from, to string) ([]ShopPlanEntry, error) {
	rows, err := q.QueryContext(ctx, `SELECT e.plan_date, s.code, e.person_id, e.recipe_id, e.servings_milli, e.is_eating_out
		FROM meal_plan_entry e JOIN meal_slot s ON s.id = e.meal_slot_id WHERE e.plan_date BETWEEN ? AND ? ORDER BY e.plan_date, s.sort_order, e.person_id, e.id`, from, to)
	return shopCollect(rows, err, func(s shopRowScanner) (ShopPlanEntry, error) {
		var e ShopPlanEntry
		err := s.Scan(&e.Date, &e.Slot, &e.PersonID, &e.RecipeID, &e.ServingsMilli, &e.IsEatingOut)
		return e, err
	})
}

// Persons devuelve código y participación de cada persona, por id.
func (shopRepo) Persons(ctx context.Context, q Querier) ([]ShopPerson, error) {
	rows, err := q.QueryContext(ctx, "SELECT id, code, participation FROM person ORDER BY id")
	return shopCollect(rows, err, func(s shopRowScanner) (ShopPerson, error) {
		var p ShopPerson
		err := s.Scan(&p.ID, &p.Code, &p.Participation)
		return p, err
	})
}

// Inventory devuelve el inventario, por id.
func (shopRepo) Inventory(ctx context.Context, q Querier) ([]ShopInventory, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, product_id, recipe_id, tracking_mode, quantity_milli, counted_date, rate_quantity_milli,
		rate_period_days, reorder_level_milli, expires_date FROM inventory_item ORDER BY id`)
	return shopCollect(rows, err, func(s shopRowScanner) (ShopInventory, error) {
		var it ShopInventory
		err := s.Scan(&it.ID, &it.ProductID, &it.RecipeID, &it.Mode, &it.QuantityMilli, &it.CountedDate, &it.RateQuantityMilli,
			&it.RatePeriodDays, &it.ReorderLevel, &it.ExpiresDate)
		return it, err
	})
}

const shopExpenseBasketSQL = `SELECT e.id, e.expense_date, e.store_id, e.amount_cents, COUNT(i.id) FROM expense e
	JOIN expense_item i ON i.expense_id = e.id`

func shopScanExpense(s shopRowScanner) (ShopExpense, error) {
	var e ShopExpense
	err := s.Scan(&e.ID, &e.Date, &e.StoreID, &e.AmountCents, &e.ItemCount)
	return e, err
}

// ExpenseBaskets devuelve las compras con renglones con el ORDER BY dado (columnas e.*).
func (shopRepo) ExpenseBaskets(ctx context.Context, q Querier, orderBy string) ([]ShopExpense, error) {
	rows, err := q.QueryContext(ctx, shopExpenseBasketSQL+" GROUP BY e.id "+orderBy)
	return shopCollect(rows, err, shopScanExpense)
}

// ExpenseBasket lee una compra con renglones y sus renglones (sql.ErrNoRows si no existe o no tiene renglones).
func (shopRepo) ExpenseBasket(ctx context.Context, q Querier, id int64) (ShopExpense, []ShopExpenseLine, error) {
	e, err := shopScanExpense(q.QueryRowContext(ctx, shopExpenseBasketSQL+" WHERE e.id = ? GROUP BY e.id", id))
	if err != nil {
		return e, nil, err
	}
	rows, err := q.QueryContext(ctx, "SELECT presentation_id, unit_price_cents, quantity_milli, line_total_cents FROM expense_item WHERE expense_id = ? ORDER BY line_no", id)
	lines, err := shopCollect(rows, err, func(s shopRowScanner) (ShopExpenseLine, error) {
		var l ShopExpenseLine
		err := s.Scan(&l.PresentationID, &l.UnitPriceCents, &l.QuantityMilli, &l.LineTotalCents)
		return l, err
	})
	return e, lines, err
}

// Exists dice si hay un renglón con ese id en la tabla (nombre fijo del código).
func (shopRepo) Exists(ctx context.Context, q Querier, table string, id int64) (bool, error) {
	var n int
	err := q.QueryRowContext(ctx, "SELECT 1 FROM "+table+" WHERE id = ?", id).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
