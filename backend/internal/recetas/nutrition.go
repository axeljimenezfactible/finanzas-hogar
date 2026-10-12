package recetas

import (
	"encoding/csv"
	"io"
	"math/big"
	"os"
	"sort"
	"strings"
)

// Tabla nutricional (S4-interfaz §9.1) y kcal con fuente citada (§9.2, §9.3; S4-08, S4-10, S4-11).

// NutritionHeader es el encabezado exacto de tabla.csv.
const NutritionHeader = "producto_id,kcal_100g,proteina_g,grasa_g,carbohidratos_g,fuente,id_fuente,nombre_fuente,fecha,estado,confirmo,fecha_confirmacion"

var sourceES = map[string]string{"USDA": "usda_fdc", "OFF": "open_food_facts", "manual": "manual"}
var licenseES = map[string]string{"USDA": "CC0", "OFF": "ODbL"}

// RowError es un renglón de tabla.csv que no se cargó (en production, motivo «unconfirmed»).
type RowError struct {
	Line    int
	Product string
	Reason  string
}

// Nutrition son los datos nutricionales preferidos por producto y, si existen, por presentación.
type Nutrition struct {
	facts    map[string]*NutritionFact // clave producto + "|" + presentación ("" = genérico)
	Rejected []RowError
}

// NewNutrition arma la tabla; en production un dato propuesto no entra (queda «sin dato»).
func NewNutrition(facts []NutritionFact, mode Mode) (*Nutrition, error) {
	n := &Nutrition{facts: map[string]*NutritionFact{}}
	for i := range facts {
		f := facts[i]
		k := f.Product + "|" + f.Presentation
		if n.facts[k] != nil || f.KcalMilliPer100g < 0 || strings.TrimSpace(f.SourceRef) == "" {
			return nil, errf(ErrInvalidArgument, "dato nutricional repetido, negativo o sin fuente: %q", f.Product)
		}
		if f.Proposed && mode == Production {
			n.Rejected = append(n.Rejected, RowError{Line: i + 2, Product: f.Product, Reason: ErrUnconfirmed.Error()})
			continue
		}
		n.facts[k] = &f
	}
	return n, nil
}

// fact es el dato que rige: el de la presentación si existe; si no, el genérico del producto (§9.3, S-5).
func (n *Nutrition) fact(product, presentation string) *NutritionFact {
	if n == nil {
		return nil
	}
	if f := n.facts[product+"|"+presentation]; f != nil && presentation != "" {
		return f
	}
	return n.facts[product+"|"]
}

// LoadNutritionCSV lee datos/semilla/nutricion/tabla.csv con las reglas de §9.1.
func LoadNutritionCSV(path string, cat *Catalog, mode Mode) (*Nutrition, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errf(ErrInvalidArgument, "no se pudo leer la tabla nutricional: %v", err)
	}
	defer f.Close()
	return ParseNutritionCSV(f, cat, mode)
}

// ParseNutritionCSV valida encabezado, producto, decimales, fuente, duplicados, estado y confirmación.
func ParseNutritionCSV(r io.Reader, cat *Catalog, mode Mode) (*Nutrition, error) {
	rd := csv.NewReader(r)
	rows, err := rd.ReadAll()
	if err != nil || len(rows) == 0 {
		return nil, errf(ErrInvalidArgument, "tabla nutricional ilegible o vacía: %v", err)
	}
	if strings.Join(rows[0], ",") != NutritionHeader {
		return nil, errf(ErrInvalidArgument, "encabezado distinto en tabla.csv")
	}
	var facts []NutritionFact
	for i, row := range rows[1:] {
		f, err := parseNutritionRow(row, cat)
		if err != nil {
			return nil, errf(err, "renglón %d", i+2)
		}
		facts = append(facts, f)
	}
	return NewNutrition(facts, mode)
}

func parseNutritionRow(row []string, cat *Catalog) (NutritionFact, error) {
	var f NutritionFact
	if len(row) != 12 {
		return f, errf(ErrInvalidArgument, "se esperaban 12 columnas")
	}
	f.Product, f.SourceLabel, f.SourceRef, f.SourceName = row[0], row[5], strings.TrimSpace(row[6]), row[7]
	if cat.Products[f.Product] == nil {
		return f, errf(ErrUnknownProduct, "%q", f.Product)
	}
	var err error
	if f.KcalMilliPer100g, err = parseMilli(row[1]); err != nil {
		return f, errf(ErrInvalidArgument, "kcal_100g: %v", err)
	}
	for j, dst := range []**int64{&f.ProteinMg, &f.FatMg, &f.CarbsMg} {
		if row[2+j] != "" {
			v, err := parseMilli(row[2+j])
			if err != nil {
				return f, errf(ErrInvalidArgument, "macronutriente: %v", err)
			}
			*dst = &v
		}
	}
	src, ok := sourceES[f.SourceLabel]
	if !ok {
		return f, errf(ErrInvalidArgument, "fuente %q", f.SourceLabel)
	}
	f.Source, f.License = src, licenseES[f.SourceLabel]
	switch {
	case f.SourceRef == "":
		return f, errf(ErrInvalidArgument, "id_fuente vacío")
	case f.SourceLabel == "USDA" && !digits(f.SourceRef):
		return f, errf(ErrInvalidArgument, "fdc_id %q", f.SourceRef)
	case f.SourceLabel == "OFF" && (!digits(f.SourceRef) || len(f.SourceRef) < 8 || len(f.SourceRef) > 14):
		return f, errf(ErrInvalidArgument, "código de barras %q", f.SourceRef)
	}
	if f.ObservedDate, err = ParseDate(row[8]); err != nil {
		return f, err
	}
	switch row[9] {
	case "propuesto":
		f.Proposed = true
		if row[10] != "" || row[11] != "" {
			return f, errf(ErrInvalidArgument, "un renglón propuesto no lleva confirmación")
		}
	case "confirmado":
		if row[10] != "axel" && row[10] != "ximena" {
			return f, errf(ErrInvalidArgument, "confirmo %q", row[10])
		}
		d, err := ParseDate(row[11])
		if err != nil || d.Before(f.ObservedDate) {
			return f, errf(ErrInvalidArgument, "fecha_confirmacion %q", row[11])
		}
		f.ConfirmedBy, f.ConfirmedDate = row[10], &d
	default:
		return f, errf(ErrInvalidArgument, "estado %q", row[9])
	}
	return f, nil
}

// Coverage cuenta los ingredientes del recetario con valor y fuente (§9.1, S4-08).
type Coverage struct {
	Unique     int      `json:"unique"`
	WithSource int      `json:"with_source"`
	Confirmed  int      `json:"confirmed"`
	Proposed   int      `json:"proposed"`
	Missing    []string `json:"missing"`
}

// NutritionCoverage: productos únicos de los ingredientes; Missing en orden de código (supuesto S-7).
func NutritionCoverage(recipes []*Recipe, nut *Nutrition) Coverage {
	c := Coverage{Missing: []string{}}
	seen := map[string]bool{}
	for _, r := range recipes {
		for _, in := range r.Ingredients {
			if seen[in.Product] {
				continue
			}
			seen[in.Product] = true
			c.Unique++
			switch f := nut.fact(in.Product, ""); {
			case f == nil:
				c.Missing = append(c.Missing, in.Product)
			case f.Proposed:
				c.Proposed++
			default:
				c.Confirmed++
			}
		}
	}
	c.WithSource = c.Confirmed + c.Proposed
	sort.Strings(c.Missing)
	return c
}

// gramsMilli: g = la cantidad; ml o pieza = cantidad × grams_per_unit_milli / 1000; sin dato → false (§9.2).
func gramsMilli(p *Product, qty *big.Rat) (*big.Rat, bool) {
	if p.BaseUnit == UnitG {
		return qty, true
	}
	if p.GramsPerUnitMilli == nil {
		return nil, false
	}
	return mulR(qty, frac(*p.GramsPerUnitMilli, 1000)), true
}

// ingredientKcal son las milésimas de kcal exactas: gramos_milli × kcal_milli_per_100g / 100000 (§9.3).
func ingredientKcal(cat *Catalog, f *NutritionFact, product string, qty *big.Rat) (*big.Rat, bool) {
	g, ok := gramsMilli(cat.Products[product], qty)
	if !ok {
		return nil, false
	}
	return mulR(g, frac(f.KcalMilliPer100g, 100000)), true
}

// TargetPercentTenths es el % de la meta con un decimal, en décimas: HALF_UP(kcal × 1000 / meta) (§9.3).
// RN-17: 1,850 contra 2,000 = 925 (92.5 %).
func TargetPercentTenths(kcal, target int64) *int64 {
	if target <= 0 {
		return nil
	}
	return ptr(roundHalfUp(frac(kcal*1000, target)))
}
