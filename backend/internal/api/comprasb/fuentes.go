package comprasb

import (
	"database/sql"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios/ocr"
	"finanzas-hogar/internal/precios/profeco"
)

// verdictMissing: desde K53 la aplicación no lee sitios y la base no guarda veredictos (S5-terminos.md es registro
// histórico). PROVISIONAL, bloqueo S6-compras-b-fuentes: cada cadena con sitio sale «missing».
const verdictMissing = "missing"

func coverageOut(cs []dominio.CBProfecoCoverage) []api.PriceSourcesProfecoCoverageItem {
	out := make([]api.PriceSourcesProfecoCoverageItem, 0, len(cs))
	for _, c := range cs {
		it := api.PriceSourcesProfecoCoverageItem{StoreID: api.Ptr(c.StoreID), PriceCount: c.PriceCount}
		if c.Municipality.Valid {
			it.Municipality = api.Ptr(c.Municipality.String)
		}
		out = append(out, it)
	}
	return out
}

func (h *handlers) getPriceSources(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	assisted, err := dominio.SettingBool(ctx, h.db(), assistedKey)
	if err != nil {
		return err
	}
	out := api.PriceSources{ManualEnabled: true, OCREngine: "Tesseract " + ocr.RequiredVersion + ", spa", AssistedReadingEnabled: assisted,
		Sites: []api.PriceSourcesSitesItem{}}
	if err := eachRow(ctx, h.db(), `SELECT id, website_url, assisted_reading_enabled FROM store WHERE kind = 'chain' AND website_url IS NOT NULL ORDER BY id`, func(sc scanner) error {
		var it api.PriceSourcesSitesItem
		var url string
		err := sc.Scan(&it.StoreID, &url, &it.AssistedReadingEnabled)
		it.WebsiteURL, it.Verdict = api.Ptr(url), verdictMissing
		out.Sites = append(out.Sites, it)
		return err
	}); err != nil {
		return err
	}
	cov, total, last, err := dominio.CBProfecoStats(ctx, h.db())
	if err != nil {
		return err
	}
	out.ProfecoCoverage, out.ProfecoPriceCount = coverageOut(cov), total
	if last.Valid { // PROVISIONAL (bloqueo S6-compras-b-fuentes): alta más reciente de un precio PROFECO
		out.ProfecoLastImportAt = api.Ptr(api.APITime(last.String))
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// importProfeco recibe el corte (multipart, parte «file»), lo guarda en un temporal de TmpDir con el tope vigente
// price.profeco_import_max_gb (K6, K44) y lo carga con el motor profeco (LoadIntoDB: filtro de tiendas, ligas
// confirmadas de vinculos.csv, sin duplicados). PROVISIONAL (bloqueo S6-compras-b-fuentes): RAR responde 415 porque
// el motor solo lee CSV; las ligas se leen de <DataDir>/profeco/vinculos.csv (sin archivo, ningún registro se liga).
func (h *handlers) importProfeco(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	cfg := h.reg.Deps().Config
	mr, err := r.MultipartReader()
	if err != nil {
		return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Manda el corte como multipart/form-data en el campo «file»."}
	}
	max, err := profeco.MaxBytesFromSettings(ctx, h.db())
	if err != nil {
		return err
	}
	var path string
	for path == "" {
		part, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			return api.BadRequest("Falta el archivo del corte en el campo «file».")
		} else if err != nil {
			return api.BadRequest("El formulario no se pudo leer.")
		}
		if part.FormName() != "file" {
			continue
		}
		mt, _, _ := mime.ParseMediaType(part.Header.Get("Content-Type"))
		switch mt {
		case "text/csv":
		case "application/vnd.rar":
			return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Por ahora solo se importa el CSV del corte: descomprime el RAR y sube el CSV."}
		default:
			return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "El corte debe ser CSV (text/csv) o RAR."}
		}
		f, err := os.CreateTemp(cfg.TmpDir, "profeco-*.csv")
		if err != nil {
			return err
		}
		path = f.Name()
		defer os.Remove(path)
		n, err := io.Copy(f, io.LimitReader(part, max+1))
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return err
		}
		if n > max {
			return &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El corte pasa del tamaño máximo de Ajustes (price.profeco_import_max_gb)."}
		}
	}
	var links *profeco.Links
	if l, err := profeco.ReadLinksFile(filepath.Join(cfg.DataDir, "profeco", "vinculos.csv")); err == nil {
		links = l
	} else if !errors.Is(err, os.ErrNotExist) {
		return dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "invalid_format", Message: "El archivo de ligas PROFECO (vinculos.csv) no es válido."})
	}
	today := h.today().String()
	res, err := profeco.LoadIntoDB(ctx, h.db(), profeco.Options{Corte: today[:7], Links: links, Mode: profeco.ModeProduction, Today: today, MaxBytes: max}, path)
	switch {
	case errors.Is(err, profeco.ErrFileTooLarge):
		return &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large", Message: "El corte pasa del tamaño máximo de Ajustes (price.profeco_import_max_gb)."}
	case errors.Is(err, profeco.ErrProfecoHeader):
		return dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "invalid_format", Message: "El corte no trae las 15 columnas del diccionario de PROFECO: revisa que sea el CSV de datos abiertos."})
	case err != nil:
		return err
	}
	cov, _, _, err := dominio.CBProfecoStats(ctx, h.db())
	if err != nil {
		return err
	}
	out := api.ProfecoImportResult{LoadedCount: res.Loaded, SkippedCount: res.Skipped, Coverage: make([]api.ProfecoImportResultCoverageItem, 0, len(cov))}
	for _, c := range coverageOut(cov) {
		out.Coverage = append(out.Coverage, api.ProfecoImportResultCoverageItem(c))
	}
	return api.WriteJSON(w, http.StatusCreated, out)
}

// ---------- nutrición y recetario ----------

var nutritionSources = []string{"usda_fdc", "open_food_facts", "mx_table", "store_label", "manual"}

func nutritionOut(n dominio.CBNutritionRow) api.NutritionFact {
	out := api.NutritionFact{ID: n.ID, ProductID: n.ProductID, KcalMilliPer100g: n.KcalMilli, Source: n.Source, IsPreferred: n.IsPreferred,
		CreatedAt: api.APITime(n.CreatedAt), UpdatedAt: api.APITime(n.UpdatedAt)}
	opt := func(v sql.NullInt64) *int64 {
		if v.Valid {
			return api.Ptr(v.Int64)
		}
		return nil
	}
	out.PresentationID, out.ProteinMgPer100g, out.FatMgPer100g, out.CarbsMgPer100g, out.FiberMgPer100g = opt(n.PresentationID), opt(n.Protein), opt(n.Fat), opt(n.Carbs), opt(n.Fiber)
	if n.SourceRef.Valid {
		out.SourceRef = api.Ptr(n.SourceRef.String)
	}
	if n.SourceLicense.Valid {
		out.SourceLicense = api.Ptr(n.SourceLicense.String)
	}
	return out
}

func (h *handlers) productOr404(r *http.Request) (int64, error) {
	pid, err := api.PathID(r, "product_id")
	if err != nil {
		return 0, err
	}
	ok, err := dominio.CBExists(r.Context(), h.db(), "product", pid)
	if err == nil && !ok {
		err = dominio.NotFound("No encontramos ese producto.")
	}
	return pid, err
}

func (h *handlers) listNutritionFacts(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"source", "created_at"}, "id")
	if err != nil {
		return err
	}
	pid, err := h.productOr404(r)
	if err != nil {
		return err
	}
	rows, err := dominio.CBQueryNutrition(r.Context(), h.db(), "WHERE product_id = ?", []any{pid},
		p.OrderBy(map[string]string{"source": "source", "created_at": "created_at", "id": "id"})+" "+p.LimitOffset())
	if err != nil {
		return err
	}
	items := make([]api.NutritionFact, 0, len(rows))
	for _, n := range rows {
		items = append(items, nutritionOut(n))
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

// createNutritionFact captura un dato con cita (N-175: manual, mx_table o store_label; las bases USDA y Open Food
// Facts entran por reloadNutritionBase). Si ya hay un valor preferido para el producto (o la presentación), uno nuevo
// marcado como preferido responde 409: la API no cambia por su cuenta cuál usan los cálculos.
func (h *handlers) createNutritionFact(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	pid, err := h.productOr404(r)
	if err != nil {
		return err
	}
	var in api.NutritionFactCreate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	var v api.Validator
	v.Check(oneOf(in.Source, "manual", "mx_table", "store_label"), "source", "not_allowed", "A mano solo se captura con fuente manual, tabla mexicana o etiqueta de tienda.")
	v.Check(in.Source == "manual" || (in.SourceRef != nil && *in.SourceRef != ""), "source_ref", "required", "Escribe la cita de la fuente (tabla y página o URL de la etiqueta).")
	v.Check(in.KcalMilliPer100g >= 0, "kcal_milli_per_100g", "out_of_range", "Las kilocalorías no pueden ser negativas.")
	for f, x := range map[string]*int64{"protein_mg_per_100g": in.ProteinMgPer100g, "fat_mg_per_100g": in.FatMgPer100g, "carbs_mg_per_100g": in.CarbsMgPer100g, "fiber_mg_per_100g": in.FiberMgPer100g} {
		v.Check(x == nil || *x >= 0, f, "out_of_range", "El valor no puede ser negativo.")
	}
	if in.PresentationID != nil {
		owner, err := dominio.CBPresentationOf(ctx, h.db(), *in.PresentationID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		v.Check(err == nil && owner == pid, "presentation_id", "not_found", "Esa presentación no es de este producto.")
	}
	if err := v.Err(); err != nil {
		return err
	}
	var id int64
	err = dominio.WithTx(ctx, h.db(), func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `INSERT INTO nutrition_fact (product_id, presentation_id, kcal_milli_per_100g, protein_mg_per_100g, fat_mg_per_100g,
 carbs_mg_per_100g, fiber_mg_per_100g, source, source_ref, source_license, is_preferred) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			pid, in.PresentationID, in.KcalMilliPer100g, in.ProteinMgPer100g, in.FatMgPer100g, in.CarbsMgPer100g, in.FiberMgPer100g,
			in.Source, in.SourceRef, in.SourceLicense, in.IsPreferred != nil && *in.IsPreferred)
		if dominio.IsUnique(err) {
			return dominio.Conflict("already_exists", "Ya hay un dato preferido para este producto; captúralo sin marcarlo como preferido.")
		} else if err != nil {
			return err
		}
		if id, err = res.LastInsertId(); err != nil {
			return err
		}
		return dominio.CBAuditWrite(ctx, tx, api.Actor(r), "insert", "nutrition_fact", id, nil, api.RequestID(r))
	})
	if err != nil {
		return err
	}
	rows, err := dominio.CBQueryNutrition(ctx, h.db(), "WHERE id = ?", []any{id}, "")
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusCreated, nutritionOut(rows[0]))
}

func (h *handlers) getRecipeSettings(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	out := api.RecipeSettings{MealSlots: []api.MealSlot{}, NutritionSources: []api.RecipeSettingsNutritionSourcesItem{}}
	// Porciones base por defecto: el DEFAULT de recipe.base_servings en el esquema (2), no una constante del código.
	if err := h.db().QueryRowContext(ctx, `SELECT CAST(dflt_value AS INTEGER) FROM pragma_table_info('recipe') WHERE name = 'base_servings'`).Scan(&out.DefaultBaseServings); err != nil {
		return err
	}
	if err := eachRow(ctx, h.db(), `SELECT id, code, name, sort_order, is_optional, created_at, updated_at FROM meal_slot ORDER BY sort_order`, func(sc scanner) error {
		var m api.MealSlot
		err := sc.Scan(&m.ID, &m.Code, &m.Name, &m.SortOrder, &m.IsOptional, &m.CreatedAt, &m.UpdatedAt)
		m.CreatedAt, m.UpdatedAt = api.APITime(m.CreatedAt), api.APITime(m.UpdatedAt)
		out.MealSlots = append(out.MealSlots, m)
		return err
	}); err != nil {
		return err
	}
	for _, o := range []string{"claude_seed", "household"} {
		it := api.RecipeSettingsRecipeCountsItem{Origin: o}
		if err := h.db().QueryRowContext(ctx, `SELECT count(*) FROM recipe WHERE origin = ?`, o).Scan(&it.Count); err != nil {
			return err
		}
		out.RecipeCounts = append(out.RecipeCounts, it)
	}
	if err := eachRow(ctx, h.db(), `SELECT source, count(DISTINCT product_id), max(created_at) FROM nutrition_fact GROUP BY source ORDER BY source`, func(sc scanner) error {
		var it api.RecipeSettingsNutritionSourcesItem
		var at sql.NullString
		err := sc.Scan(&it.Source, &it.FoodCount, &at)
		if at.Valid {
			it.LoadedAt = api.Ptr(api.APITime(at.String))
		}
		out.NutritionSources = append(out.NutritionSources, it)
		return err
	}); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// reloadNutritionBase: PROVISIONAL (bloqueo S6-compras-b-nutricion). No hay decisión sobre el archivo local de cada
// base (USDA FDC, Open Food Facts) ni sobre la llave que liga sus alimentos con product (product no tiene código);
// mientras tanto valida la fuente y responde 409 sin tocar nutrition_fact.
func (h *handlers) reloadNutritionBase(w http.ResponseWriter, r *http.Request) error {
	var in api.NutritionReloadRequest
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	var v api.Validator
	v.Check(oneOf(in.Source, "usda_fdc", "open_food_facts"), "source", "invalid_format", "La base debe ser usda_fdc u open_food_facts.")
	if err := v.Err(); err != nil {
		return err
	}
	return dominio.Conflict("conflict", "Todavía no está lista la recarga de la base nutricional local: falta definir su archivo. Puedes capturar datos con cita mientras tanto.")
}
