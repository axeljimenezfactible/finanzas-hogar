package semilla

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"finanzas-hogar/internal/recetas"
)

// profecoNames son los nombres de cadena_comercial que PROFECO usa en la muestra
// (datos/semilla/profeco/muestra/qqp-muestra.csv: «Wal-mart», «Bodega Aurrera»). PROVISIONAL: se confirman al
// importar PROFECO completo (S5-interfaz §11).
var profecoNames = map[string]string{"Bodega Aurrera": "Bodega Aurrera", "Walmart": "Wal-mart"}

var reStore = regexp.MustCompile(`^(.+?) \(sucursal (.+)\)$`)

// tiendas: las dos cadenas y sus sucursales de data.json (tienda_a y tienda_b). data.json no trae mercados ni
// membresías: Axel los captura en Ajustes (P-26.06).
func (l *loader) tiendas() error {
	d, err := readData(l.opt.Dir)
	if err != nil {
		return err
	}
	for _, full := range []string{d.TiendaA, d.TiendaB} {
		m := reStore.FindStringSubmatch(full)
		if m == nil {
			return fmt.Errorf("semilla: data.json: no entiendo el nombre de tienda %q (se esperaba «Cadena (sucursal X)»)", full)
		}
		chain, branch := m[1], m[2]
		var prof any
		if p, ok := profecoNames[chain]; ok {
			prof = p
		}
		cid, _, err := l.ensure("tiendas (cadenas)", "store", `SELECT id FROM store WHERE kind = 'chain' AND name = ? COLLATE NOCASE`, a(chain),
			`INSERT INTO store (kind, parent_id, name, profeco_name) VALUES ('chain', NULL, ?, ?)`, chain, prof)
		if err != nil {
			return err
		}
		l.ids["store:"+chain] = cid
		bid, _, err := l.ensure("tiendas (sucursales)", "store", `SELECT id FROM store WHERE parent_id = ? AND name = ? COLLATE NOCASE`, a(cid, branch),
			`INSERT INTO store (kind, parent_id, name) VALUES ('branch', ?, ?)`, cid, branch)
		if err != nil {
			return err
		}
		l.ids["store:"+full] = bid
	}
	return nil
}

type catalogFile struct {
	Grupos []struct {
		ID, Nombre, Canal string
		Orden             int
	} `json:"grupos"`
	Productos []struct {
		ID             string  `json:"id"`
		Nombre         string  `json:"nombre"`
		Categoria      string  `json:"categoria"`
		UnidadBase     string  `json:"unidad_base"`
		Gramos         *string `json:"gramos_por_unidad"`
		FuenteGramos   *string `json:"fuente_gramos"`
		Presentaciones []struct {
			ID      string   `json:"id"`
			Nombre  string   `json:"nombre"`
			Marca   *string  `json:"marca"`
			Milli   int64    `json:"contenido_milli"`
			Granel  bool     `json:"granel"`
			Barcode *string  `json:"codigo_barras"`
			Tickets []string `json:"nombres_ticket"`
		} `json:"presentaciones"`
	} `json:"productos"`
}

var channelDB = map[string]string{"mercado": "market", "cadena": "chain"}
var unitDB = map[string]string{"g": "g", "ml": "ml", "pieza": "piece"}

// milli convierte un decimal de texto con hasta 3 decimales a milésimas ("0.913" → 913; "19.94" → 19940).
func milli(s string) (int64, error) {
	ip, fp, dot := strings.Cut(strings.TrimSpace(s), ".")
	if ip == "" || (dot && (fp == "" || len(fp) > 3)) {
		return 0, fmt.Errorf("decimal inválido %q (hasta 3 decimales)", s)
	}
	fp += strings.Repeat("0", 3-len(fp))
	return strconv.ParseInt(ip+fp, 10, 64)
}

// catalogo: 9 grupos, 160 productos y 186 presentaciones de productos.json (validado antes con recetas.LoadCatalogFile).
func (l *loader) catalogo() error {
	path := filepath.Join(l.opt.Dir, "productos.json")
	cat, err := recetas.LoadCatalogFile(path)
	if err != nil {
		return fmt.Errorf("semilla: productos.json: %w", err)
	}
	l.cat = cat
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var f catalogFile
	if err := json.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("semilla: productos.json: %w", err)
	}
	for _, g := range f.Grupos {
		id, _, err := l.ensure("grupos de producto", "product_group", `SELECT id FROM product_group WHERE name = ? COLLATE NOCASE`, a(g.Nombre),
			`INSERT INTO product_group (name, purchase_channel, sort_order) VALUES (?, ?, ?)`, g.Nombre, channelDB[g.Canal], g.Orden)
		if err != nil {
			return err
		}
		l.ids["group:"+g.ID] = id
	}
	for _, p := range f.Productos {
		var grams, notes any
		if p.Gramos != nil {
			m, err := milli(*p.Gramos)
			if err != nil {
				return fmt.Errorf("semilla: producto %s: %w", p.ID, err)
			}
			grams = m
		}
		if p.FuenteGramos != nil {
			notes = *p.FuenteGramos
		}
		pid, _, err := l.ensure("productos", "product", `SELECT id FROM product WHERE name = ? COLLATE NOCASE`, a(p.Nombre),
			`INSERT INTO product (name, product_group_id, base_unit, grams_per_unit_milli, notes) VALUES (?, ?, ?, ?, ?)`,
			p.Nombre, l.id("group:"+p.Categoria), unitDB[p.UnidadBase], grams, notes)
		if err != nil {
			return err
		}
		l.ids["product:"+p.ID] = pid
		for _, x := range p.Presentaciones {
			var brand, bar any
			if x.Marca != nil {
				brand = *x.Marca
			}
			if x.Barcode != nil {
				bar = *x.Barcode
			}
			xid, _, err := l.ensure("presentaciones", "presentation", `SELECT id FROM presentation WHERE product_id = ? AND name = ?`, a(pid, x.Nombre),
				`INSERT INTO presentation (product_id, name, brand, size_milli, is_bulk, barcode) VALUES (?, ?, ?, ?, ?, ?)`,
				pid, x.Nombre, brand, x.Milli, x.Granel, bar)
			if err != nil {
				return err
			}
			l.ids["pres:"+x.ID] = xid
			l.pres[x.Nombre] = xid
			for _, t := range x.Tickets {
				l.pres[t] = xid
			}
		}
	}
	return nil
}

// nutricion: solo los renglones «confirmado» de nutricion/tabla.csv entran a la base (Q2-43 = A: «nada propuesto entra a la
// base de la app»). Hoy los 150 están «propuesto»: se cuentan y se avisa. Lo confirma Axel (bloqueo S4-empates-sin-confirmar).
func (l *loader) nutricion() error {
	path := filepath.Join(l.opt.Dir, "nutricion", "tabla.csv")
	if _, err := os.Stat(path); err != nil {
		l.rep.Notices = append(l.rep.Notices, "No encontré nutricion/tabla.csv: no se cargó la tabla nutricional.")
		return nil
	}
	nut, err := recetas.LoadNutritionCSV(path, l.cat, recetas.Review) // valida encabezado, productos, decimales y confirmación
	if err != nil {
		return fmt.Errorf("semilla: nutricion/tabla.csv: %w", err)
	}
	l.nut = nut
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()
	rows, err := csv.NewReader(fh).ReadAll()
	if err != nil {
		return fmt.Errorf("semilla: nutricion/tabla.csv: %w", err)
	}
	proposed := 0
	srcDB := map[string]string{"USDA": "usda_fdc", "OFF": "open_food_facts", "manual": "manual"}
	licDB := map[string]any{"USDA": "CC0", "OFF": "ODbL"}
	for _, r := range rows[1:] {
		if r[9] != "confirmado" {
			proposed++
			continue
		}
		kcal, err := milli(r[1])
		if err != nil {
			return fmt.Errorf("semilla: tabla.csv %s: %w", r[0], err)
		}
		macro := func(s string) (any, error) {
			if s == "" {
				return nil, nil
			}
			return milli(s)
		}
		var m [3]any
		for i := range m {
			if m[i], err = macro(r[2+i]); err != nil {
				return fmt.Errorf("semilla: tabla.csv %s: %w", r[0], err)
			}
		}
		pid := l.id("product:" + r[0])
		if _, _, err := l.ensure("datos nutricionales", "nutrition_fact",
			`SELECT id FROM nutrition_fact WHERE product_id = ? AND presentation_id IS NULL AND source = ? AND source_ref = ?`, a(pid, srcDB[r[5]], r[6]),
			`INSERT INTO nutrition_fact (product_id, presentation_id, kcal_milli_per_100g, protein_mg_per_100g, fat_mg_per_100g, carbs_mg_per_100g, source, source_ref, source_license, is_preferred)
			 VALUES (?, NULL, ?, ?, ?, ?, ?, ?, ?, 1)`, pid, kcal, m[0], m[1], m[2], srcDB[r[5]], r[6], licDB[r[5]]); err != nil {
			return err
		}
	}
	if proposed > 0 {
		l.rep.Notices = append(l.rep.Notices, fmt.Sprintf("%d renglones de nutricion/tabla.csv están «propuesto» y NO se cargaron: nada propuesto entra a la base hasta que Axel o Ximena lo confirmen (Q2-43 = A; bloqueo S4-empates-sin-confirmar, Q2-90 sin respuesta). Mientras tanto las recetas salen sin kcal.", proposed))
	}
	return nil
}

// recetario: las 60 recetas (origen claude_seed) con sus tiempos e ingredientes. Se validan con recetas.LoadRecipeBook
// (modo revisión: cotejo de kcal contra el motor S4) y deben ser las 60 completas.
func (l *loader) recetario() error {
	path := filepath.Join(l.opt.Dir, "recetario.json")
	if _, err := os.Stat(path); err != nil {
		l.rep.Notices = append(l.rep.Notices, "No encontré recetario.json: no se cargó el recetario.")
		return nil
	}
	book, err := recetas.LoadRecipeBook(path, l.cat, l.nut, recetas.Review)
	if err != nil {
		return fmt.Errorf("semilla: recetario.json: %w", err)
	}
	if err := book.SeedComplete(); err != nil {
		return fmt.Errorf("semilla: recetario.json: %w", err)
	}
	slots := map[string]int64{}
	rows, err := l.tx.QueryContext(l.ctx, `SELECT code, id FROM meal_slot`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var c string
		var id int64
		if err := rows.Scan(&c, &id); err != nil {
			rows.Close()
			return err
		}
		slots[c] = id
	}
	rows.Close()
	for _, r := range book.Recipes {
		var cook, desc any
		if r.CookMinutes != nil {
			cook = *r.CookMinutes
		}
		if r.Description != "" {
			desc = r.Description
		}
		rid, created, err := l.ensure("recetas", "recipe", `SELECT id FROM recipe WHERE name = ? COLLATE NOCASE`, a(r.Name),
			`INSERT INTO recipe (name, description, instructions, base_servings, prep_minutes, cook_minutes, origin) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			r.Name, desc, strings.Join(r.Steps, "\n"), r.BaseServings, r.PrepMinutes, cook, r.Origin)
		if err != nil {
			return err
		}
		if !created {
			continue
		}
		for _, s := range r.Slots {
			if _, err := l.insert("recipe_meal_slot", `INSERT INTO recipe_meal_slot (recipe_id, meal_slot_id) VALUES (?, ?)`, rid, slots[string(s)]); err != nil {
				return err
			}
		}
		for _, in := range r.Ingredients {
			var pres any
			if in.Presentation != "" {
				pres = l.id("pres:" + in.Presentation)
			}
			var txt any
			if in.DisplayText != "" {
				txt = in.DisplayText
			}
			if _, err := l.insert("recipe_ingredient", `INSERT INTO recipe_ingredient (recipe_id, product_id, presentation_id, quantity_milli, display_text, is_optional, sort_order)
				VALUES (?, ?, ?, ?, ?, ?, ?)`, rid, l.id("product:"+in.Product), pres, in.QuantityMilli, txt, in.Optional, in.SortOrder); err != nil {
				return err
			}
			l.rep.count("ingredientes de receta", true)
		}
	}
	return nil
}
