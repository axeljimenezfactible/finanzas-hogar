package semilla

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios/profeco"
)

// profeco: solo la MUESTRA de datos/semilla/profeco/muestra/*.csv, filtrada contra las tiendas sembradas y ligada con
// profeco/vinculos.csv si existe (liga confirmada por una persona; en producción nada propuesto cuenta, Q2-43 = A).
// Hoy no hay vinculos.csv y las sucursales no traen municipio (data.json no lo da): ningún precio pasa el filtro y
// el reporte lo dice. Los CSV completos se importan después con importProfeco (Ajustes), no aquí.
func (l *loader) profeco() error {
	files, _ := filepath.Glob(filepath.Join(l.opt.Dir, "profeco", "muestra", "*.csv"))
	sort.Strings(files)
	if len(files) == 0 {
		l.rep.Notices = append(l.rep.Notices, "No hay muestra de PROFECO en profeco/muestra/: no se cargaron precios de referencia.")
		return nil
	}
	var set profeco.StoreSet
	rows, err := l.tx.QueryContext(l.ctx, `SELECT id, kind, coalesce(parent_id, 0), name, coalesce(profeco_name, ''), coalesce(state, ''), coalesce(municipality, '') FROM store ORDER BY id`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id, parent int64
		var kind, name, prof, state, mun string
		if err := rows.Scan(&id, &kind, &parent, &name, &prof, &state, &mun); err != nil {
			rows.Close()
			return err
		}
		switch kind {
		case "chain":
			set.Chains = append(set.Chains, profeco.Chain{ID: id, Name: name, ProfecoName: prof})
		case "branch":
			set.Branches = append(set.Branches, profeco.Branch{ID: id, ChainID: parent, Name: name, ProfecoName: prof, State: state, Municipality: mun})
		}
	}
	rows.Close()
	var links *profeco.Links
	if p := filepath.Join(l.opt.Dir, "profeco", "vinculos.csv"); fileExists(p) {
		if links, err = profeco.ReadLinksFile(p); err != nil {
			return fmt.Errorf("semilla: profeco/vinculos.csv: %w", err)
		}
	}
	clock := l.opt.Clock
	if clock == nil {
		clock = dominio.SystemClock{}
	}
	im, err := profeco.NewImporter(profeco.Options{Corte: "muestra", Stores: set, Links: links, Mode: profeco.ModeProduction,
		Prices: priceSink{l}, Today: dominio.Today(clock).String(), MaxBytes: 64 << 20})
	if err != nil {
		return err
	}
	for _, f := range files {
		res, err := im.ImportFile(l.ctx, f)
		if err != nil {
			return fmt.Errorf("semilla: PROFECO %s: %w", filepath.Base(f), err)
		}
		if !res.Balanced() {
			return fmt.Errorf("semilla: PROFECO %s: los conteos no cuadran (se perdió algún renglón)", filepath.Base(f))
		}
		for i := int64(0); i < res.Loaded; i++ {
			l.rep.count("precios PROFECO (muestra)", true)
		}
		for i := int64(0); i < res.Duplicates; i++ {
			l.rep.count("precios PROFECO (muestra)", false)
		}
		l.rep.Notices = append(l.rep.Notices, fmt.Sprintf("PROFECO muestra %s: %d renglones leídos, %d cargados, %d fuera de las tiendas y municipios sembrados, %d sin liga al catálogo, %d rechazados. "+
			"Sin municipio en las sucursales (data.json no lo trae) y sin profeco/vinculos.csv confirmado no entra ningún precio (bloqueos S6-semilla-sucursales-municipio y S4-empates-sin-confirmar).",
			res.File, res.Read, res.Loaded, res.FilteredOut, res.Skipped, res.Rejected))
	}
	return nil
}

type priceSink struct{ l *loader }

// AddPrice inserta un precio PROFECO; el mismo precio (presentación, tienda y día) cuenta como duplicado (Q2-55 = A).
func (s priceSink) AddPrice(ctx context.Context, p profeco.Price) (bool, error) {
	res, err := s.l.tx.ExecContext(ctx, `INSERT OR IGNORE INTO price (presentation_id, store_id, price_cents, observed_date, source, is_promotional, source_ref)
		VALUES (?, ?, ?, ?, 'profeco', 0, ?)`, p.PresentationID, p.StoreID, p.PriceCents, p.ObservedDate, p.SourceRef)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }

// partidasCSV: gastos_fijos_ingresos.csv lo entrega Axel (plan §8.2). Si falta, la carga sigue sin esa parte y se avisa;
// si existe, se importa con el importador de Ajustes (Options.ImportPlannedCSV) o se avisa que falta engancharlo.
func (l *loader) partidasCSV() error {
	p := filepath.Join(l.opt.Dir, "gastos_fijos_ingresos.csv")
	if !fileExists(p) {
		l.rep.Notices = append(l.rep.Notices, "Falta datos/semilla/gastos_fijos_ingresos.csv (lo entrega Axel; plantilla en plan.md §8.2): sin él no hay ingresos ni gastos fijos. "+
			"Cuando exista se importa en Ajustes → Datos y respaldos (importPlannedItems) o volviendo a correr «finanzas semilla».")
		return nil
	}
	if l.opt.ImportPlannedCSV == nil {
		l.rep.Notices = append(l.rep.Notices, "gastos_fijos_ingresos.csv existe pero esta carga no tiene el importador enganchado (solicitud fuera de alcance a api-ajustes): impórtalo en Ajustes → Datos y respaldos (importPlannedItems).")
		return nil
	}
	n, err := l.opt.ImportPlannedCSV(l.ctx, l.tx, p)
	if err != nil {
		return fmt.Errorf("semilla: gastos_fijos_ingresos.csv: %w", err)
	}
	for i := 0; i < n; i++ {
		l.rep.count("partidas de ingresos y gastos fijos", true)
	}
	return nil
}

// avisos: lo que Axel debe saber de esta carga.
func (l *loader) avisos() error {
	d, err := readData(l.opt.Dir)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(string(d.Promociones))) > 0 && string(d.Promociones) != "null" {
		l.rep.Notices = append(l.rep.Notices, "data.json trae «promociones», pero su formato no está definido en ningún contrato: no se cargaron (bloqueo semilla-promociones).")
	} else {
		l.rep.Notices = append(l.rep.Notices, "data.json no trae promociones: no se sembró ninguna (Axel las captura en Ajustes → Créditos y condiciones).")
	}
	l.rep.Notices = append(l.rep.Notices, "Las contraseñas iniciales son solo para entrar la primera vez: cámbialas en Ajustes → Cuenta y roles (cambiar mi contraseña).")
	return nil
}
