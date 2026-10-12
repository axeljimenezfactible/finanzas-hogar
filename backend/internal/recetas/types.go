package recetas

// Tipos del dominio (S4-interfaz §2). Las entidades se citan por su código (texto); la API traduce
// códigos e ids numéricos (§12.1).

// Unit es la unidad base de un producto.
type Unit string

// Unidades base (product.base_unit).
const (
	UnitG     Unit = "g"
	UnitML    Unit = "ml"
	UnitPiece Unit = "piece"
)

// Channel es la vertical de compra de un grupo (product_group.purchase_channel).
type Channel string

// Verticales de compra (decisiones D, RN-16).
const (
	ChannelMarket Channel = "market"
	ChannelChain  Channel = "chain"
)

// SlotCode es un tiempo de comida (meal_slot.code).
type SlotCode string

// Los 4 tiempos (RN-14); la colación es opcional.
const (
	Breakfast SlotCode = "breakfast"
	Lunch     SlotCode = "lunch"
	Dinner    SlotCode = "dinner"
	Snack     SlotCode = "snack"
)

var slotOrder = map[SlotCode]int{Breakfast: 0, Lunch: 1, Dinner: 2, Snack: 3}

// mandatorySlots son los tiempos obligatorios que cuentan como «casillas sin armar» (K43).
var mandatorySlots = []SlotCode{Breakfast, Lunch, Dinner}

// Mode es el modo de carga de tabla.csv y del recetario (§3.2, §9.1).
type Mode int

// Modos de carga.
const (
	Review     Mode = iota // pruebas, fixtures y revisión: carga propuestos y confirmados, coteja kcal
	Production             // base de la app: solo lo confirmado, sin cotejo
)

// Orígenes de receta (recipe.origin) y modos de inventario (inventory_item.mode).
const (
	OriginSeed      = "claude_seed"
	OriginHousehold = "household"
	ModeStock       = "stock"
	ModeRate        = "consumption_rate"
)

// Group es un grupo de compra.
type Group struct {
	Code, Name string
	Channel    Channel
	SortOrder  int
}

// Product es un producto del catálogo.
type Product struct {
	ID                  int64
	Code, Name          string
	BaseUnit            Unit
	GramsPerUnitMilli   *int64 // ml o pieza → gramos × 1000
	Group               string // código del grupo; "" = sin grupo (error de catálogo en la lista, §6.5)
	TypicalPresentation string
}

// Presentation es una presentación comprable.
type Presentation struct {
	ID                  int64
	Code, Product, Name string
	SizeMilli           int64 // contenido en milésimas de la unidad base
}

// Ingredient es un renglón de receta, en milésimas de la unidad base para BaseServings.
type Ingredient struct {
	Product, Presentation string
	QuantityMilli         int64
	DisplayText           string
	Optional              bool
	SortOrder             int
}

// Recipe es una receta (recetario de Claude o de la casa).
type Recipe struct {
	ID                      int64
	Code, Name, Description string
	Slots                   []SlotCode
	BaseServings            int
	PrepMinutes             int
	CookMinutes             *int
	Origin                  string
	Steps                   []string
	Ingredients             []Ingredient
	IsActive                bool
}

// Person es una persona; solo las que tienen plan (titular y co-titular) entran al plan y a salud.
type Person struct {
	Code        string
	HasMealPlan bool
}

// PlanEntry es un renglón del plan: receta + porciones, o «fuera» (RN-19).
type PlanEntry struct {
	Date          Date
	Slot          SlotCode
	Person        string
	Recipe        string // "" en «fuera»
	ServingsMilli *int64
	EatingOut     bool
}

// InventoryItem es un artículo del inventario híbrido (RN-15); con Recipe es un sobrante en porciones × 1000.
type InventoryItem struct {
	Code, Product, Recipe string
	Mode                  string
	QuantityMilli         int64
	CountedDate           *Date
	RateQuantityMilli     *int64
	RatePeriodDays        *int64
	ReorderLevelMilli     *int64
	ExpiresDate           *Date
}

// Movement es un movimiento de inventario que la API guarda junto con el cambio de existencia.
type Movement struct {
	Item       string
	Date       Date
	Kind       string // purchase, recipe_leftover, plan_consumption, adjustment
	Reason     string // count, consumption, waste (solo en ajustes)
	DeltaMilli int64
}

// NutritionFact es el dato nutricional preferido de un producto (o de una presentación).
type NutritionFact struct {
	Product, Presentation     string
	KcalMilliPer100g          int64
	ProteinMg, FatMg, CarbsMg *int64
	Source, SourceRef         string // usda_fdc, open_food_facts, manual, mx_table, store_label
	SourceLabel               string // como viene en tabla.csv: USDA, OFF o manual
	SourceName                string // nombre en la fuente (procedencia; nunca se muestra, S4-11)
	License                   string
	ObservedDate              Date
	Proposed                  bool
	ConfirmedBy               string
	ConfirmedDate             *Date
}

// HealthProfile es el perfil de salud de una persona (00001 + 00002 + 00003).
type HealthProfile struct {
	Person                         string
	HeightMM                       *int64
	BirthDate                      *Date
	Sex, Activity, Goal            *string
	AdjustLoseKcal, AdjustGainKcal *int64
	DailyKcalTarget                *int64
	DietaryRestrictions            *string // existe y está vacío; el núcleo no lo lee (decisiones D)
}

// WeightEntry es un peso del historial.
type WeightEntry struct {
	Person  string
	Date    Date
	WeightG int64
}

// Settings son los valores de la casa (app_setting health.goal_adjust_*; K36, K44): 0 y 0 al nacer.
type Settings struct {
	HouseAdjustLoseKcal, HouseAdjustGainKcal int64
}

// Store es una tienda: cadena, sucursal o mercado con su próxima visita (S5 §9).
type Store struct {
	Code, Kind string
	Groups     []string
	NextVisit  *Date
}

// Quote es una cotización vigente que entrega S5 por PriceSource (§8.1).
type Quote struct {
	Product, Store, Presentation string
	PriceCents, SizeMilli        int64
	ObservedDate                 Date
	IsStale, Recommended         bool
}

// PriceSource es la entrada de precios de S5 (§8.1): por producto y fecha de referencia, la cotización
// recomendada (exactamente una) y las de las elecciones de la persona; vacío = «sin precio».
// S4 no importa internal/precios: la API conecta la implementación de S5.
type PriceSource interface {
	Quotes(product string, ref Date) []Quote
}

// QuoteList es un PriceSource en memoria (fixtures y pruebas).
type QuoteList []Quote

// Quotes devuelve las cotizaciones del producto.
func (l QuoteList) Quotes(product string, _ Date) []Quote {
	var out []Quote
	for _, q := range l {
		if q.Product == product {
			out = append(out, q)
		}
	}
	return out
}

// ListChoice es una elección de la persona en la lista (shopping_list_item): se conserva al regenerar.
type ListChoice struct {
	Product, Kind, Store, Presentation string // Kind: calculated o manual
	QuantityMilli                      *int64
	IsBought                           bool
}

// Engine reúne los datos ya leídos que usan las funciones del motor. Lo arma la API (S6) desde SQLite
// o Eval desde una fixture; Today lo calcula el llamador (C-05).
type Engine struct {
	Today     Date
	Catalog   *Catalog
	Recipes   []*Recipe
	Nutrition *Nutrition
	Persons   []Person
	Stores    []Store
	Plan      []PlanEntry
	Inventory []InventoryItem
	Choices   []ListChoice
	Prices    PriceSource
	Profiles  []HealthProfile
	Weights   []WeightEntry
	Settings  Settings
}

func (e *Engine) recipe(code string) (*Recipe, error) {
	for _, r := range e.Recipes {
		if r.Code == code {
			return r, nil
		}
	}
	return nil, errf(ErrInvalidArgument, "receta desconocida %q", code)
}

func (e *Engine) person(code string) (*Person, error) {
	for i := range e.Persons {
		if e.Persons[i].Code == code {
			return &e.Persons[i], nil
		}
	}
	return nil, errf(ErrInvalidArgument, "persona desconocida %q", code)
}
