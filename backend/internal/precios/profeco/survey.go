package profeco

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func h64(parts ...string) uint64 {
	h := uint64(14695981039346656037)
	for _, p := range parts {
		for i := 0; i < len(p); i++ {
			h ^= uint64(p[i])
			h *= 1099511628211
		}
		h ^= 0xff
		h *= 1099511628211
	}
	return h
}

// containsPhrase dice si la frase (ya con Fold) aparece como palabras completas
// dentro del texto (ya con Fold).
func containsPhrase(folded, phrase string) bool {
	return strings.Contains(" "+folded+" ", " "+phrase+" ")
}

// IsFocusChain dice si la cadena (ya con Fold) es de las que interesan a Axel:
// la familia Walmart, la familia Bodega Aurrera o Sam's Club.
func IsFocusChain(folded string) bool {
	return strings.HasPrefix(folded, "walmart") || strings.HasPrefix(folded, "bodega aurrera") || containsPhrase(folded, "sams")
}

// Agg acumula la cobertura de un grupo (una cadena, o una cadena en un municipio).
// Los conjuntos usan un hash de 64 bits del texto (colisión esperada ≈ 1e-9).
type Agg struct {
	Records       int64
	products      map[uint64]struct{}
	presentations map[uint64]struct{}
	branches      map[uint64]struct{}
	Months        map[string]struct{}
	MinDate       string
	MaxDate       string
}

func newAgg() *Agg {
	return &Agg{products: map[uint64]struct{}{}, presentations: map[uint64]struct{}{}, branches: map[uint64]struct{}{}, Months: map[string]struct{}{}}
}

func (a *Agg) add(r *Record) {
	a.Records++
	a.products[h64(r.Producto)] = struct{}{}
	a.presentations[h64(r.Producto, r.Presentacion, r.Marca)] = struct{}{}
	a.branches[h64(r.NombreComercial, r.Direccion, r.Municipio)] = struct{}{}
	a.Months[r.Fecha[:7]] = struct{}{}
	if a.MinDate == "" || r.Fecha < a.MinDate {
		a.MinDate = r.Fecha
	}
	if r.Fecha > a.MaxDate {
		a.MaxDate = r.Fecha
	}
}

// Products es el número de productos distintos (columna producto).
func (a *Agg) Products() int { return len(a.products) }

// Presentations es el número de combinaciones distintas producto + presentación + marca.
func (a *Agg) Presentations() int { return len(a.presentations) }

// Branches es el número de sucursales distintas (nombre_comercial + direccion + municipio).
func (a *Agg) Branches() int { return len(a.branches) }

// NamedBranch es una sucursal concreta que se busca por texto. La regla se fija
// antes de contar: el registro cuenta si la cadena (cuando Chains no está vacía)
// es una de Chains y la frase aparece como palabras completas en nombre_comercial
// o en direccion (o, si AnyField, también en cadena_comercial).
type NamedBranch struct {
	Label    string
	Chains   []string // con Fold; vacío = cualquier cadena
	Phrase   string   // con Fold
	AnyField bool
	Near     string // frase más laxa, solo informativa (coincidencias parciales)
}

// DefaultNamedBranches son las sucursales de Axel (S5-interfaz §11).
var DefaultNamedBranches = []NamedBranch{
	{Label: "Bodega Aurrera Vallejo", Chains: []string{"bodega aurrera", "bodega aurrera express"}, Phrase: "vallejo", Near: "vallejo"},
	{Label: "Walmart Acueducto de Guadalupe", Chains: []string{"walmart", "walmart express"}, Phrase: "acueducto de guadalupe", Near: "acueducto"},
	{Label: "Sam's Club", Phrase: "sams club", AnyField: true, Near: "sams"},
}

// NamedResult es el resultado de buscar una sucursal concreta.
type NamedResult struct {
	NamedBranch
	Agg      *Agg
	DirOnly  *Agg             // parte de Agg donde la frase solo aparece en la dirección (referencia, no el nombre del local)
	Examples map[string]int64 // sucursales que cumplen la regla: nombre | dirección | municipio → registros
	NearAgg  *Agg
	NearEx   map[string]int64
}

// FileStat son los conteos de un archivo leído.
type FileStat struct {
	Name     string
	Size     int64
	Lines    int64
	Read     int64
	Rejected int64
	MinDate  string
	MaxDate  string
	Elapsed  time.Duration
}

// Survey recorre archivos y mide la cobertura de todas las cadenas, sin guardar
// los registros.
type Survey struct {
	Files      []FileStat
	Total      *Agg
	Chains     map[string]*Agg
	ChainMun   map[[3]string]*Agg // cadena, estado, municipio (solo cadenas de interés)
	Formats    map[string]int64
	Rejects    map[string]int64
	Rejections []Rejection
	Named      []*NamedResult
	// WatchMun son los municipios (con Fold) cuyas sucursales de las cadenas de
	// interés se listan completas: ahí viven las tiendas de Axel.
	WatchMun []string
	Watch    map[[2]string]map[string]int64 // municipio, cadena → «nombre | dirección» → registros
	Elapsed  time.Duration
}

// NewSurvey prepara un recorrido con las sucursales de Axel por defecto.
func NewSurvey() *Survey {
	s := &Survey{Total: newAgg(), Chains: map[string]*Agg{}, ChainMun: map[[3]string]*Agg{}, Formats: map[string]int64{}, Rejects: map[string]int64{},
		WatchMun: []string{"gustavo a madero"}, Watch: map[[2]string]map[string]int64{}}
	for _, nb := range DefaultNamedBranches {
		s.Named = append(s.Named, &NamedResult{NamedBranch: nb, Agg: newAgg(), DirOnly: newAgg(), Examples: map[string]int64{}, NearAgg: newAgg(), NearEx: map[string]int64{}})
	}
	return s
}

// Observe suma un registro.
func (s *Survey) Observe(r *Record) {
	s.Formats[r.FechaFormato]++
	s.Total.add(r)
	a := s.Chains[r.CadenaComercial]
	if a == nil {
		a = newAgg()
		s.Chains[r.CadenaComercial] = a
	}
	a.add(r)
	fc := Fold(r.CadenaComercial)
	if IsFocusChain(fc) {
		k := [3]string{r.CadenaComercial, r.Estado, r.Municipio}
		m := s.ChainMun[k]
		if m == nil {
			m = newAgg()
			s.ChainMun[k] = m
		}
		m.add(r)
		if mun := Fold(r.Municipio); inList(s.WatchMun, mun) {
			wk := [2]string{r.Municipio, r.CadenaComercial}
			if s.Watch[wk] == nil {
				s.Watch[wk] = map[string]int64{}
			}
			s.Watch[wk][r.NombreComercial+" | "+r.Direccion]++
		}
	}
	if len(s.Named) == 0 {
		return
	}
	nom, dir := Fold(r.NombreComercial), Fold(r.Direccion)
	for _, n := range s.Named {
		if len(n.Chains) > 0 && !inList(n.Chains, fc) {
			if !(n.AnyField && (containsPhrase(nom, n.Phrase) || containsPhrase(dir, n.Phrase) || containsPhrase(fc, n.Phrase))) {
				continue
			}
		}
		id := r.NombreComercial + " | " + r.Direccion + " | " + r.Municipio
		inName := containsPhrase(nom, n.Phrase) || (n.AnyField && containsPhrase(fc, n.Phrase))
		if inName || containsPhrase(dir, n.Phrase) {
			n.Agg.add(r)
			n.Examples[id]++
			if !inName {
				n.DirOnly.add(r)
			}
		} else if n.Near != "" && (containsPhrase(nom, n.Near) || containsPhrase(dir, n.Near) || containsPhrase(fc, n.Near)) {
			n.NearAgg.add(r)
			n.NearEx[id]++
		}
	}
}

func inList(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

// ScanFile lee un archivo y suma sus registros. Los renglones rechazados se
// cuentan con su motivo.
func (s *Survey) ScanFile(ctx context.Context, path string, maxBytes int64, maxRejections int) error {
	start := time.Now()
	f, err := OpenFile(path, maxBytes)
	if err != nil {
		return err
	}
	defer f.Close()
	st := FileStat{Name: filepath.Base(path), Size: f.Size}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		rec, err := f.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		var re *RowError
		if errors.As(err, &re) {
			st.Read++
			st.Rejected++
			s.Rejects[re.Reason]++
			if len(s.Rejections) < maxRejections {
				s.Rejections = append(s.Rejections, Rejection{File: st.Name, Row: re.Line, Reason: re.Reason})
			}
			continue
		}
		if err != nil {
			return err
		}
		st.Read++
		if st.MinDate == "" || rec.Fecha < st.MinDate {
			st.MinDate = rec.Fecha
		}
		if rec.Fecha > st.MaxDate {
			st.MaxDate = rec.Fecha
		}
		s.Observe(&rec)
	}
	st.Lines = f.Lines()
	st.Elapsed = time.Since(start)
	s.Elapsed += st.Elapsed
	s.Files = append(s.Files, st)
	return nil
}

// sortedChains devuelve las cadenas por número de registros (y nombre).
func (s *Survey) sortedChains() []string {
	names := make([]string, 0, len(s.Chains))
	for n := range s.Chains {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		a, b := s.Chains[names[i]].Records, s.Chains[names[j]].Records
		if a != b {
			return a > b
		}
		return names[i] < names[j]
	})
	return names
}
