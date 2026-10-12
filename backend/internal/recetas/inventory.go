package recetas

// Inventario híbrido (§7; RN-15, K7; S4-04).

// estimateAt es la existencia a una fecha: stock = quantity_milli; consumption_rate = contada −
// HALF_UP(rate_quantity_milli × días desde el conteo / rate_period_days), nunca menor que 0 (RN-15).
func estimateAt(it InventoryItem, d Date) int64 {
	if it.Mode != ModeRate {
		return it.QuantityMilli
	}
	days := max(it.CountedDate.DaysUntil(d), 0)
	return max(it.QuantityMilli-roundHalfUp(frac(*it.RateQuantityMilli*days, *it.RatePeriodDays)), 0)
}

// weeklyConsumption = HALF_UP(rate_quantity_milli × 7 / rate_period_days).
func weeklyConsumption(it InventoryItem) int64 {
	return roundHalfUp(frac(*it.RateQuantityMilli*7, *it.RatePeriodDays))
}

// Estados de inventario (K7): agotado prevalece sobre reponer.
const (
	StatusSufficient = "sufficient"
	StatusReorder    = "reorder"
	StatusOut        = "out"
)

// ItemStatus es un artículo de inventory_status.
type ItemStatus struct {
	Code                   string `json:"code"`
	EstimatedQuantityMilli int64  `json:"estimated_quantity_milli"`
	Status                 string `json:"status"`
	WeeklyConsumptionMilli *int64 `json:"weekly_consumption_milli"`
	DaysLeft               *int64 `json:"days_left"`
}

// InventorySummary es el conteo por estado.
type InventorySummary struct {
	Sufficient int `json:"sufficient"`
	Reorder    int `json:"reorder"`
	Out        int `json:"out"`
}

// InventoryReport es la salida de inventory_status.
type InventoryReport struct {
	Items   []ItemStatus     `json:"items"`
	Summary InventorySummary `json:"summary"`
}

// StatusOf: out si la existencia es 0; reorder si hay nivel y existencia ≤ nivel; si no, sufficient.
func StatusOf(estimated int64, reorderLevel *int64) string {
	switch {
	case estimated == 0:
		return StatusOut
	case reorderLevel != nil && estimated <= *reorderLevel:
		return StatusReorder
	}
	return StatusSufficient
}

// InventoryStatus calcula cada artículo a «hoy» (§7).
func (e *Engine) InventoryStatus() (out *InventoryReport, err error) {
	defer catch(&err)
	out = &InventoryReport{Items: []ItemStatus{}}
	for _, it := range e.Inventory {
		est := estimateAt(it, e.Today)
		s := ItemStatus{Code: it.Code, EstimatedQuantityMilli: est, Status: StatusOf(est, it.ReorderLevelMilli)}
		if it.Mode == ModeRate {
			s.WeeklyConsumptionMilli = ptr(weeklyConsumption(it))
			s.DaysLeft = ptr(floorRat(frac(est**it.RatePeriodDays, *it.RateQuantityMilli)))
		}
		switch s.Status {
		case StatusOut:
			out.Summary.Out++
		case StatusReorder:
			out.Summary.Reorder++
		default:
			out.Summary.Sufficient++
		}
		out.Items = append(out.Items, s)
	}
	return out, nil
}

// Tipos y motivos de movimiento (inventory_movement).
const (
	MovePurchase       = "purchase"
	MoveLeftover       = "recipe_leftover"
	MovePlanConsumed   = "plan_consumption"
	MoveAdjustment     = "adjustment"
	ReasonCount        = "count"
	ReasonConsumption  = "consumption"
	ReasonWaste        = "waste"
	errAdjustmentLabel = "motivo de ajuste"
)

// PurchaseMovement: + paquetes × contenido.
func PurchaseMovement(item string, d Date, packages, sizeMilli int64) Movement {
	return Movement{Item: item, Date: d, Kind: MovePurchase, DeltaMilli: packages * sizeMilli}
}

// PlanConsumptionMovement: − cantidad escalada sin bajar de 0; devuelve además el faltante.
func PlanConsumptionMovement(item string, d Date, current, qtyMilli int64) (Movement, int64) {
	used := min(current, qtyMilli)
	return Movement{Item: item, Date: d, Kind: MovePlanConsumed, DeltaMilli: -used}, qtyMilli - used
}

// LeftoverMovement: + porciones × 1000 de un sobrante de receta.
func LeftoverMovement(item string, d Date, servingsMilli int64) Movement {
	return Movement{Item: item, Date: d, Kind: MoveLeftover, DeltaMilli: servingsMilli}
}

// AdjustmentMovement: la diferencia contra la cantidad nueva, con motivo count, consumption o waste.
func AdjustmentMovement(item string, d Date, reason string, current, newQty int64) (Movement, error) {
	if reason != ReasonCount && reason != ReasonConsumption && reason != ReasonWaste {
		return Movement{}, errf(ErrInvalidArgument, "%s %q", errAdjustmentLabel, reason)
	}
	return Movement{Item: item, Date: d, Kind: MoveAdjustment, Reason: reason, DeltaMilli: newQty - current}, nil
}
