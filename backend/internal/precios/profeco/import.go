package profeco

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"time"
)

// Motivos de omisión (skipped_count de C1).
const (
	SkipNoLink      = "sin presentación ligada"
	SkipLinkPending = "liga sin confirmar"
)

// Price es un precio de referencia listo para la tabla price con fuente PROFECO.
type Price struct {
	PresentationID int64
	StoreID        int64
	PriceCents     int64
	ObservedDate   string // ISO
	SourceRef      string // «QQP <corte> <archivo>:<renglón>»
}

// PriceSink recibe los precios filtrados y ligados. Devuelve inserted = false si
// ya existía el mismo precio (presentación, tienda y día): el primero gana y los
// demás se cuentan como duplicados (FIRME, puerta 2, Q2-55 = A).
type PriceSink interface {
	AddPrice(ctx context.Context, p Price) (inserted bool, err error)
}

// RecordSink recibe todos los registros bien formados (FIRME, puerta 2,
// Q2-53 = A: base local aparte qqp.sqlite).
type RecordSink interface {
	AddRecord(ctx context.Context, file string, rec *Record) error
}

// Options configura un Importer.
type Options struct {
	Corte         string // etiqueta del corte (por ejemplo «2026») para source_ref
	Stores        StoreSet
	Links         *Links
	Mode          Mode
	Prices        PriceSink  // opcional
	Records       RecordSink // opcional
	Today         string     // ISO; si no está vacío, rechaza fechas futuras
	MaxBytes      int64      // límite por archivo; 0 = sin límite
	MaxRejections int        // renglones rechazados con detalle; 0 = 100 (PROVISIONAL, bloqueo S5-motor-oraculo-reglas pregunta 6)
}

// Rejection es el detalle de un renglón rechazado.
// Rejection es un renglón rechazado. File es la ruta tal como se pidió importar y Row el
// ordinal del registro de datos (1 = el primero tras el encabezado), la misma convención que
// source_ref y profeco_record.renglon (solicitud a S5-interfaz §11 y §15.2).
type Rejection struct {
	File   string
	Row    int
	Reason string
}

// CoverageRow cuenta los precios cargados por cadena y municipio (C1 coverage).
type CoverageRow struct {
	Chain        string
	Municipality string
	PriceCount   int64
}

// ImportResult son los conteos de una importación. Se cumple siempre:
// Read = Rejected + FilteredOut + Skipped + Loaded + LoadedProposed + Duplicates.
type ImportResult struct {
	File           string
	Size           int64
	Lines          int64 // saltos de línea del archivo (equivale a wc -l)
	Read           int64 // renglones de datos leídos (sin encabezado)
	Rejected       int64
	FilteredOut    int64
	Skipped        int64
	Loaded         int64 // con liga confirmada
	LoadedProposed int64 // con liga propuesta (solo modo revisión)
	Duplicates     int64
	SkippedReasons map[string]int64
	RejectReasons  map[string]int64
	Formats        map[string]int64 // formatos de fecha reconocidos
	Rejections     []Rejection
	Coverage       []CoverageRow
	MinDate        string
	MaxDate        string
	Elapsed        time.Duration
}

// Balanced indica que no se perdió ningún renglón sin explicar.
func (r *ImportResult) Balanced() bool {
	return r.Read == r.Rejected+r.FilteredOut+r.Skipped+r.Loaded+r.LoadedProposed+r.Duplicates
}

// Importer filtra, liga y carga archivos de PROFECO en flujo.
type Importer struct {
	o      Options
	filter *Filter
}

// NewImporter valida el modo y arma el filtro.
func NewImporter(o Options) (*Importer, error) {
	if o.Mode == "" {
		o.Mode = ModeProduction
	}
	if o.Mode != ModeProduction && o.Mode != ModeReview {
		return nil, fmt.Errorf("profeco: modo desconocido %q", o.Mode)
	}
	if o.MaxRejections == 0 {
		o.MaxRejections = 100
	}
	return &Importer{o: o, filter: NewFilter(o.Stores)}, nil
}

// ImportFile procesa un CSV y devuelve sus conteos.
func (im *Importer) ImportFile(ctx context.Context, path string) (*ImportResult, error) {
	start := time.Now()
	f, err := OpenFile(path, im.o.MaxBytes)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	res := &ImportResult{File: filepath.Base(path), Size: f.Size, SkippedReasons: map[string]int64{}, RejectReasons: map[string]int64{}, Formats: map[string]int64{}}
	cov := map[[2]string]int64{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		rec, err := f.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		var re *RowError
		if errors.As(err, &re) {
			res.Read++
			res.Rejected++
			res.RejectReasons[re.Reason]++
			if len(res.Rejections) < im.o.MaxRejections {
				res.Rejections = append(res.Rejections, Rejection{File: path, Row: re.Row, Reason: re.Reason})
			}
			continue
		}
		if err != nil {
			return nil, err
		}
		res.Read++
		if im.o.Today != "" && rec.Fecha > im.o.Today {
			res.Rejected++
			res.RejectReasons["fecha futura"]++
			if len(res.Rejections) < im.o.MaxRejections {
				res.Rejections = append(res.Rejections, Rejection{File: path, Row: rec.Row, Reason: "fecha futura"})
			}
			continue
		}
		res.Formats[rec.FechaFormato]++
		if res.MinDate == "" || rec.Fecha < res.MinDate {
			res.MinDate = rec.Fecha
		}
		if rec.Fecha > res.MaxDate {
			res.MaxDate = rec.Fecha
		}
		if im.o.Records != nil {
			if err := im.o.Records.AddRecord(ctx, res.File, &rec); err != nil {
				return nil, err
			}
		}
		if err := im.route(ctx, res, cov, &rec); err != nil {
			return nil, err
		}
	}
	res.Lines = f.Lines()
	for k, n := range cov {
		res.Coverage = append(res.Coverage, CoverageRow{Chain: k[0], Municipality: k[1], PriceCount: n})
	}
	sortCoverage(res.Coverage)
	res.Elapsed = time.Since(start)
	return res, nil
}

func (im *Importer) route(ctx context.Context, res *ImportResult, cov map[[2]string]int64, rec *Record) error {
	_, storeID, ok := im.filter.Match(rec)
	if !ok {
		res.FilteredOut++
		return nil
	}
	link, found := im.o.Links.Lookup(rec)
	switch {
	case !found:
		res.Skipped++
		res.SkippedReasons[SkipNoLink]++
		return nil
	case link.Status != LinkConfirmed && im.o.Mode == ModeProduction:
		res.Skipped++
		res.SkippedReasons[SkipLinkPending]++
		return nil
	}
	inserted := true
	if im.o.Prices != nil {
		var err error
		inserted, err = im.o.Prices.AddPrice(ctx, Price{
			PresentationID: link.PresentationID, StoreID: storeID, PriceCents: rec.PrecioCents,
			ObservedDate: rec.Fecha, SourceRef: fmt.Sprintf("QQP %s %s:%d", im.o.Corte, res.File, rec.Row),
		})
		if err != nil {
			return err
		}
	}
	switch {
	case !inserted:
		res.Duplicates++
	case link.Status == LinkConfirmed:
		res.Loaded++
		cov[[2]string{rec.CadenaComercial, rec.Municipio}]++
	default:
		res.LoadedProposed++
		cov[[2]string{rec.CadenaComercial, rec.Municipio}]++
	}
	return nil
}

func sortCoverage(c []CoverageRow) {
	sort.Slice(c, func(i, j int) bool {
		if c[i].Chain != c[j].Chain {
			return c[i].Chain < c[j].Chain
		}
		return c[i].Municipality < c[j].Municipality
	})
}

// Import procesa varios archivos y devuelve el detalle por archivo y el total.
// Se detiene en el primer error de archivo.
func (im *Importer) Import(ctx context.Context, paths ...string) (total *ImportResult, files []*ImportResult, err error) {
	total = &ImportResult{SkippedReasons: map[string]int64{}, RejectReasons: map[string]int64{}, Formats: map[string]int64{}}
	cov := map[[2]string]int64{}
	for _, p := range paths {
		r, err := im.ImportFile(ctx, p)
		if err != nil {
			return total, files, fmt.Errorf("%s: %w", filepath.Base(p), err)
		}
		files = append(files, r)
		total.Size += r.Size
		total.Lines += r.Lines
		total.Read += r.Read
		total.Rejected += r.Rejected
		total.FilteredOut += r.FilteredOut
		total.Skipped += r.Skipped
		total.Loaded += r.Loaded
		total.LoadedProposed += r.LoadedProposed
		total.Duplicates += r.Duplicates
		total.Elapsed += r.Elapsed
		for k, v := range r.SkippedReasons {
			total.SkippedReasons[k] += v
		}
		for k, v := range r.RejectReasons {
			total.RejectReasons[k] += v
		}
		for k, v := range r.Formats {
			total.Formats[k] += v
		}
		for _, x := range r.Rejections {
			if len(total.Rejections) < im.o.MaxRejections {
				total.Rejections = append(total.Rejections, x)
			}
		}
		for _, c := range r.Coverage {
			cov[[2]string{c.Chain, c.Municipality}] += c.PriceCount
		}
		if r.MinDate != "" && (total.MinDate == "" || r.MinDate < total.MinDate) {
			total.MinDate = r.MinDate
		}
		if r.MaxDate > total.MaxDate {
			total.MaxDate = r.MaxDate
		}
	}
	for k, n := range cov {
		total.Coverage = append(total.Coverage, CoverageRow{Chain: k[0], Municipality: k[1], PriceCount: n})
	}
	sortCoverage(total.Coverage)
	return total, files, nil
}

// MemoryPrices es un PriceSink en memoria (pruebas y modo revisión): el primero
// de cada (presentación, tienda, día) gana.
type MemoryPrices struct {
	Rows []Price
	seen map[[3]string]bool
}

// AddPrice implementa PriceSink.
func (m *MemoryPrices) AddPrice(_ context.Context, p Price) (bool, error) {
	if m.seen == nil {
		m.seen = map[[3]string]bool{}
	}
	k := [3]string{fmt.Sprint(p.PresentationID), fmt.Sprint(p.StoreID), p.ObservedDate}
	if m.seen[k] {
		return false, nil
	}
	m.seen[k] = true
	m.Rows = append(m.Rows, p)
	return true, nil
}
