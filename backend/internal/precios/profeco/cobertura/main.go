// Comando cobertura: mide la cobertura de los CSV de PROFECO (Quién es Quién en
// los Precios) y, si se pide, carga todos los registros en una base aparte o
// arma la muestra pequeña de pruebas. Solo lee archivos locales; sin red.
//
// Reproduce los conteos de docs/estado/S5-cobertura-profeco.md:
//
//	cd backend && go run ./internal/precios/profeco/cobertura -dir ../datos/semilla/profeco/raw/2026-csv/2026
package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"finanzas-hogar/internal/precios/profeco"

	_ "modernc.org/sqlite"
)

func main() {
	dir := flag.String("dir", "", "carpeta o carpetas (separadas por coma) con los CSV de PROFECO; una por corte (obligatorio)")
	maxGB := flag.Int64("max-gb", 1, "tamaño máximo por archivo en GB (price.profeco_import_max_gb)")
	qqp := flag.String("qqp", "", "si se da, carga todos los registros en esta base SQLite aparte (tabla profeco_record)")
	muestra := flag.String("muestra", "", "si se da, escribe aquí la muestra CSV de pruebas (en vez del reporte)")
	nMuestra := flag.Int("n", 2000, "tamaño aproximado de la muestra")
	cadenas := flag.String("cadenas", "", "mide el filtro del cargador: cadenas_comerciales separadas por coma (con -municipios)")
	municipios := flag.String("municipios", "", "mide el filtro del cargador: municipios separados por coma (con -cadenas)")
	flag.Parse()
	if *dir == "" {
		fmt.Fprintln(os.Stderr, "falta -dir")
		os.Exit(2)
	}
	var files []string
	for _, d := range strings.Split(*dir, ",") {
		m, err := filepath.Glob(filepath.Join(strings.TrimSpace(d), "*.csv"))
		if err != nil || len(m) == 0 {
			fmt.Fprintln(os.Stderr, "no hay CSV en", d)
			os.Exit(2)
		}
		files = append(files, m...)
	}
	var err error
	sort.Strings(files)
	ctx := context.Background()
	switch {
	case *cadenas != "" || *municipios != "":
		err = filtro(ctx, files, *cadenas, *municipios, *maxGB<<30)
	case *qqp != "":
		err = loadQQP(ctx, files, *qqp, *maxGB<<30)
	case *muestra != "":
		err = buildSample(ctx, files, *muestra, *nMuestra, *maxGB<<30)
	default:
		err = report(ctx, files, *maxGB<<30)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func report(ctx context.Context, files []string, maxBytes int64) error {
	s := profeco.NewSurvey()
	for _, f := range files {
		if err := s.ScanFile(ctx, f, maxBytes, 50); err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(f), err)
		}
	}
	s.WriteMarkdown(os.Stdout)
	fmt.Fprintf(os.Stderr, "tiempo de lectura (solo conteos): %s\n", s.Elapsed.Round(time.Second))
	return nil
}

// filtro corre el cargador real (lectura, filtro y liga) sobre todos los
// archivos, sin ligas, y reporta cuántos registros pasan y cuánto tarda. Sirve
// para medir el costo de una importación y para que Axel vea el tamaño de la
// selección que haga en Ajustes.
func filtro(ctx context.Context, files []string, cadenas, municipios string, maxBytes int64) error {
	var set profeco.StoreSet
	for i, c := range strings.Split(cadenas, ",") {
		if c = strings.TrimSpace(c); c != "" {
			set.Chains = append(set.Chains, profeco.Chain{ID: int64(i + 1), Name: c, ProfecoName: c})
		}
	}
	for _, m := range strings.Split(municipios, ",") {
		if m = strings.TrimSpace(m); m != "" {
			set.Branches = append(set.Branches, profeco.Branch{ID: 100, ChainID: 1, Municipality: m})
		}
	}
	im, err := profeco.NewImporter(profeco.Options{Corte: "2026", Stores: set, Mode: profeco.ModeReview, MaxBytes: maxBytes})
	if err != nil {
		return err
	}
	total, _, err := im.Import(ctx, files...)
	if err != nil {
		return err
	}
	fmt.Printf("leidos=%d rechazados=%d filtrados_fuera=%d pasan_el_filtro(sin liga)=%d tiempo=%s\n",
		total.Read, total.Rejected, total.FilteredOut, total.Skipped, total.Elapsed.Round(time.Second))
	return nil
}

func loadQQP(ctx context.Context, files []string, path string, maxBytes int64) error {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=journal_mode(WAL)&_pragma=synchronous(OFF)")
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	q, err := profeco.NewQQP(ctx, db)
	if err != nil {
		return err
	}
	im, err := profeco.NewImporter(profeco.Options{Corte: "2026", Records: q, Mode: profeco.ModeReview, MaxBytes: maxBytes})
	if err != nil {
		return err
	}
	start := time.Now()
	total, per, err := im.Import(ctx, files...)
	if err != nil {
		return err
	}
	if err := q.Flush(); err != nil {
		return err
	}
	for _, r := range per {
		fmt.Printf("%s leidos=%d rechazados=%d lineas=%d\n", r.File, r.Read, r.Rejected, r.Lines)
	}
	fmt.Printf("total leidos=%d rechazados=%d guardados=%d min=%s max=%s tiempo=%s\n",
		total.Read, total.Rejected, q.Rows, total.MinDate, total.MaxDate, time.Since(start).Round(time.Second))
	return nil
}

// buildSample escribe una muestra determinista: por cada cadena un cupo de
// registros (más grande para las cadenas de interés), repartidos entre todos
// los archivos. Dos pasadas: la primera cuenta, la segunda elige.
func buildSample(ctx context.Context, files []string, out string, n int, maxBytes int64) error {
	counts := map[string]int64{}
	for _, f := range files {
		if err := scan(ctx, f, maxBytes, func(r *profeco.Record, _ string) { counts[r.CadenaComercial]++ }); err != nil {
			return err
		}
	}
	quota := map[string]int64{}
	focus, other := 0, 0
	for c := range counts {
		if profeco.IsFocusChain(profeco.Fold(c)) {
			focus++
		} else {
			other++
		}
	}
	// Presupuesto: 55 % a las cadenas de interés, 45 % al resto.
	for c, total := range counts {
		var q int64
		if profeco.IsFocusChain(profeco.Fold(c)) {
			q = int64(n) * 55 / 100 / int64(focus)
		} else {
			q = int64(n) * 45 / 100 / int64(other)
		}
		if q < 5 {
			q = 5
		}
		if q > total {
			q = total
		}
		quota[c] = q
	}
	taken := map[string]int64{}
	of, err := os.Create(out)
	if err != nil {
		return err
	}
	defer of.Close()
	w := csv.NewWriter(of)
	if err := w.Write(profeco.Columns); err != nil {
		return err
	}
	for _, f := range files {
		err := scan(ctx, f, maxBytes, func(r *profeco.Record, file string) {
			c := r.CadenaComercial
			// Selección determinista: cada cadena toma un renglón de cada
			// counts/quota (paso fijo) para repartir entre meses y municipios.
			step := counts[c] / quota[c]
			if step < 1 {
				step = 1
			}
			if taken[c] < quota[c] && int64(profecoHash(file, r.Line))%step == 0 {
				taken[c]++
				_ = w.Write([]string{r.Producto, r.Presentacion, r.Marca, r.Categoria, r.Catalogo, r.PrecioTexto, r.FechaTexto,
					r.CadenaComercial, r.Giro, r.NombreComercial, r.Direccion, r.Estado, r.Municipio, r.Latitud, r.Longitud})
			}
		})
		if err != nil {
			return err
		}
	}
	w.Flush()
	var tot int64
	for _, t := range taken {
		tot += t
	}
	fmt.Fprintf(os.Stderr, "muestra: %d renglones, %d cadenas\n", tot, len(taken))
	return w.Error()
}

func profecoHash(file string, line int) uint64 {
	h := uint64(14695981039346656037)
	for i := 0; i < len(file); i++ {
		h ^= uint64(file[i])
		h *= 1099511628211
	}
	h ^= uint64(line)
	h *= 1099511628211
	h ^= h >> 29
	return h
}

func scan(ctx context.Context, path string, maxBytes int64, fn func(*profeco.Record, string)) error {
	f, err := profeco.OpenFile(path, maxBytes)
	if err != nil {
		return err
	}
	defer f.Close()
	name := filepath.Base(path)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		rec, err := f.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			var re *profeco.RowError
			if errors.As(err, &re) {
				continue
			}
			return err
		}
		fn(&rec, name)
	}
}
