package recetas

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Cargador y validación del recetario (S4-interfaz §3, S4-01, S4-11, S4-13).

// SeedRecipeCount es el número exacto de recetas de la semilla (S4-01).
const SeedRecipeCount = 60

var idPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// KcalSource es un elemento de fuentes_kcal[] (§3.2): copiado tal cual de tabla.csv.
type KcalSource struct {
	Product string `json:"producto_id"`
	Source  string `json:"fuente"`
	Ref     string `json:"id_fuente"`
}

type bookIngredient struct {
	Producto     string `json:"producto_id"`
	Cantidad     int64  `json:"cantidad_milli"`
	Unidad       string `json:"unidad"`
	Presentacion string `json:"presentacion_id"`
	Texto        string `json:"texto"`
	Opcional     bool   `json:"opcional"`
}

type bookRecipe struct {
	ID           string           `json:"id"`
	Nombre       string           `json:"nombre"`
	Descripcion  string           `json:"descripcion"`
	Tiempos      []string         `json:"tiempos"`
	Porciones    int              `json:"porciones_base"`
	Prep         *int             `json:"minutos_preparacion"`
	Coccion      *int             `json:"minutos_coccion"`
	Origen       string           `json:"origen"`
	Ingredientes []bookIngredient `json:"ingredientes"`
	Pasos        []string         `json:"pasos"`
	Kcal         *int64           `json:"kcal_por_porcion"`
	Fuentes      []KcalSource     `json:"fuentes_kcal"`
}

type bookFile struct {
	Version  string       `json:"version"`
	Fecha    string       `json:"fecha"`
	Catalogo string       `json:"catalogo"`
	Recetas  []bookRecipe `json:"recetas"`
}

// RecipeError es una receta rechazada por el cargador, con su motivo.
type RecipeError struct {
	Recipe string
	Err    error
}

// RecipeBook es el recetario cargado: las recetas válidas y las rechazadas.
type RecipeBook struct {
	Recipes  []*Recipe
	Rejected []RecipeError
}

// SeedComplete exige las 60 recetas de la semilla, todas válidas (S4-01).
func (b *RecipeBook) SeedComplete() error {
	if len(b.Rejected) > 0 {
		return errf(b.Rejected[0].Err, "receta %q rechazada (%d en total)", b.Rejected[0].Recipe, len(b.Rejected))
	}
	if len(b.Recipes) != SeedRecipeCount {
		return errf(ErrInvalidArgument, "la semilla trae %d recetas y deben ser %d", len(b.Recipes), SeedRecipeCount)
	}
	return nil
}

var originES = map[string]string{"claude": OriginSeed, "propia": OriginHousehold}

// LoadRecipeBook lee recetario.json (§3.2, «Cargador del recetario y modos»). La tabla nutricional va
// cargada en el mismo modo; puede ser nil (sin tabla: toda receta sin kcal).
func LoadRecipeBook(path string, cat *Catalog, nut *Nutrition, mode Mode) (*RecipeBook, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errf(ErrInvalidArgument, "no se pudo leer el recetario: %v", err)
	}
	return ParseRecipeBook(b, cat, nut, mode)
}

// ParseRecipeBook valida cada receta (esquema de §3.2 y validación cruzada) y, en modo review, coteja
// kcal_por_porcion y fuentes_kcal[] contra el cálculo del motor (ErrKcalMismatch). En production no coteja.
func ParseRecipeBook(b []byte, cat *Catalog, nut *Nutrition, mode Mode) (book *RecipeBook, err error) {
	defer catch(&err)
	// Solo los campos de §3.1 y §3.2: una llave que el contrato no define (por ejemplo `etiquetas`) es error,
	// igual que en recetario.schema.json (additionalProperties: false).
	var f bookFile
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&f); err != nil {
		return nil, errf(ErrInvalidArgument, "recetario: %v", err)
	}
	book = &RecipeBook{}
	ids, names := map[string]bool{}, map[string]bool{}
	english := map[string]bool{}
	if nut != nil {
		for _, fct := range nut.facts {
			if fct.SourceName != "" {
				english[strings.ToLower(fct.SourceName)] = true
			}
		}
	}
	for _, br := range f.Recetas {
		r, err := bookToRecipe(br, cat)
		if err == nil && (ids[br.ID] || names[strings.ToLower(br.Nombre)]) {
			err = errf(ErrInvalidArgument, "id o nombre repetido")
		}
		if err == nil {
			err = checkEnglish(r, cat, english)
		}
		if err == nil && mode == Review {
			err = checkKcal(br, r, cat, nut)
		}
		if err != nil {
			book.Rejected = append(book.Rejected, RecipeError{br.ID, err})
			continue
		}
		ids[br.ID], names[strings.ToLower(br.Nombre)] = true, true
		book.Recipes = append(book.Recipes, r)
	}
	return book, nil
}

func bookToRecipe(br bookRecipe, cat *Catalog) (*Recipe, error) {
	origin, ok := originES[br.Origen]
	if !ok {
		return nil, errf(ErrInvalidArgument, "origen %q", br.Origen)
	}
	if br.Prep == nil || br.Coccion == nil {
		return nil, errf(ErrMissingField, "minutos_preparacion y minutos_coccion son obligatorios")
	}
	r := &Recipe{Code: br.ID, Name: br.Nombre, Description: br.Descripcion, BaseServings: br.Porciones,
		PrepMinutes: *br.Prep, CookMinutes: br.Coccion, Origin: origin, Steps: br.Pasos, IsActive: true}
	for _, s := range br.Tiempos {
		r.Slots = append(r.Slots, SlotCode(s))
	}
	for i, in := range br.Ingredientes {
		if p := cat.Products[in.Producto]; p != nil {
			if u, ok := unitES[in.Unidad]; !ok || u != p.BaseUnit {
				return nil, errf(ErrUnitMismatch, "%q: unidad %q y unidad base %q", in.Producto, in.Unidad, p.BaseUnit)
			}
		}
		r.Ingredients = append(r.Ingredients, Ingredient{Product: in.Producto, Presentation: in.Presentacion,
			QuantityMilli: in.Cantidad, DisplayText: in.Texto, Optional: in.Opcional, SortOrder: i})
	}
	return r, validateRecipe(r, cat, true)
}

// checkEnglish: ningún nombre visible (receta o producto usado) igual a un nombre_fuente (S4-11).
func checkEnglish(r *Recipe, cat *Catalog, english map[string]bool) error {
	if english[strings.ToLower(r.Name)] {
		return errf(ErrInvalidArgument, "nombre en inglés %q", r.Name)
	}
	for _, in := range r.Ingredients {
		if english[strings.ToLower(cat.Products[in.Product].Name)] {
			return errf(ErrInvalidArgument, "producto con nombre en inglés %q", in.Product)
		}
	}
	return nil
}

// checkKcal recalcula kcal_por_porcion y fuentes_kcal[] (modo review, §3.2).
func checkKcal(br bookRecipe, r *Recipe, cat *Catalog, nut *Nutrition) error {
	var want *int64
	var srcs []KcalSource
	sum, missing := r64(0), false
	for _, in := range r.Ingredients {
		f := nut.fact(in.Product, in.Presentation)
		if f == nil {
			missing = true
			continue
		}
		srcs = append(srcs, KcalSource{in.Product, f.SourceLabel, f.SourceRef})
		k, ok := ingredientKcal(cat, f, in.Product, r64(in.QuantityMilli))
		if !ok {
			missing = true
			continue
		}
		sum = addR(sum, k)
	}
	if !missing {
		want = ptr(roundHalfUp(quoR(sum, r64(1000*int64(r.BaseServings)))))
	}
	if (want == nil) != (br.Kcal == nil) || (want != nil && *want != *br.Kcal) {
		return errf(ErrKcalMismatch, "%q: kcal_por_porcion escrito %v, calculado %v", r.Code, show(br.Kcal), show(want))
	}
	if len(srcs) != len(br.Fuentes) {
		return errf(ErrKcalMismatch, "%q: fuentes_kcal con %d elementos y se esperaban %d", r.Code, len(br.Fuentes), len(srcs))
	}
	for i := range srcs {
		if srcs[i] != br.Fuentes[i] {
			return errf(ErrKcalMismatch, "%q: fuentes_kcal[%d] distinto", r.Code, i)
		}
	}
	return nil
}

func show(p *int64) any {
	if p == nil {
		return "null"
	}
	return *p
}

// ValidateRecipe aplica las reglas de §3.2 a una receta nueva o editada (S4-13).
func ValidateRecipe(r *Recipe, cat *Catalog) error { return validateRecipe(r, cat, true) }

func validateRecipe(r *Recipe, cat *Catalog, checkSteps bool) error {
	n := utf8.RuneCountInString(r.Name)
	switch {
	case !idPattern.MatchString(r.Code):
		return errf(ErrInvalidArgument, "id de receta %q", r.Code)
	case n < 3 || n > 80:
		return errf(ErrInvalidArgument, "nombre de 3 a 80 caracteres: %q", r.Name)
	case utf8.RuneCountInString(r.Description) > 200:
		return errf(ErrInvalidArgument, "descripción de más de 200 caracteres")
	case len(r.Slots) == 0:
		return errf(ErrInvalidArgument, "%q sin tiempos", r.Code)
	case r.PrepMinutes < 0 || (r.CookMinutes != nil && *r.CookMinutes < 0):
		return errf(ErrInvalidArgument, "minutos negativos en %q", r.Code)
	case len(r.Ingredients) == 0:
		return errf(ErrInvalidArgument, "%q sin ingredientes", r.Code)
	}
	switch r.Origin {
	case OriginSeed:
		if r.BaseServings != 2 {
			return errf(ErrInvalidServings, "el recetario de Claude es para 2 porciones (%q)", r.Code)
		}
	case OriginHousehold:
		if r.BaseServings < 1 || r.BaseServings > 50 {
			return errf(ErrInvalidServings, "porciones base de 1 a 50 (%q)", r.Code)
		}
	default:
		return errf(ErrInvalidArgument, "origen %q", r.Origin)
	}
	seen := map[SlotCode]bool{}
	for _, s := range r.Slots {
		if _, ok := slotOrder[s]; !ok || seen[s] {
			return errf(ErrInvalidArgument, "tiempo %q inválido o repetido en %q", s, r.Code)
		}
		seen[s] = true
	}
	prods := map[string]bool{}
	for _, in := range r.Ingredients {
		if cat.Products[in.Product] == nil {
			return errf(ErrUnknownProduct, "%q en %q", in.Product, r.Code)
		}
		if _, err := cat.presentationOf(in.Product, in.Presentation); err != nil {
			return err
		}
		if prods[in.Product] || in.QuantityMilli <= 0 || (r.Origin == OriginSeed && in.Optional) {
			return errf(ErrInvalidArgument, "ingrediente %q repetido, sin cantidad u opcional en el recetario de Claude (%q)", in.Product, r.Code)
		}
		prods[in.Product] = true
	}
	if checkSteps {
		if len(r.Steps) == 0 {
			return errf(ErrInvalidArgument, "%q sin pasos", r.Code)
		}
		for _, s := range r.Steps {
			if strings.TrimSpace(s) == "" {
				return errf(ErrInvalidArgument, "paso vacío en %q", r.Code)
			}
		}
	}
	return nil
}

// ApplyRecipeEdit devuelve la receta editada y validada sin cambiar su código, su id ni su origen (§3.3).
func ApplyRecipeEdit(old, edit *Recipe, cat *Catalog) (*Recipe, error) {
	r := *edit
	r.ID, r.Code, r.Origin = old.ID, old.Code, old.Origin
	r.Slots = append([]SlotCode(nil), edit.Slots...)
	r.Ingredients = append([]Ingredient(nil), edit.Ingredients...)
	r.Steps = append([]string(nil), edit.Steps...)
	if err := ValidateRecipe(&r, cat); err != nil {
		return nil, err
	}
	return &r, nil
}

// Resultados de RemoveRecipe (C-06).
const (
	RemoveArchive = "archive" // is_active = 0: el plan la usa
	RemoveDelete  = "delete"  // sin historia
)

// RemoveRecipe decide archivo lógico o borrado (§3.3, C-06).
func RemoveRecipe(r *Recipe, usedByPlan bool) string {
	if usedByPlan {
		return RemoveArchive
	}
	return RemoveDelete
}
