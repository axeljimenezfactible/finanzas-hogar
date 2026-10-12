package profeco

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Columns son las 15 columnas de datos/semilla/profeco/QQP_diccionario.csv.
var Columns = []string{
	"producto", "presentacion", "marca", "categoria", "catalogo", "precio",
	"fecha_registro", "cadena_comercial", "giro", "nombre_comercial", "direccion",
	"estado", "municipio", "latitud", "longitud",
}

// Record es un renglón ya normalizado. Los textos van recortados y en UTF-8; el
// dinero en centavos; la fecha en ISO; latitud y longitud como texto decimal.
type Record struct {
	Line            int // línea física donde empieza el registro (1 = encabezado); solo diagnóstico
	Row             int // renglón: ordinal del registro de datos en su archivo (1 = el primero tras el encabezado; S5-interfaz §11)
	Producto        string
	Presentacion    string
	Marca           string
	Categoria       string
	Catalogo        string
	PrecioTexto     string
	PrecioCents     int64
	FechaTexto      string
	Fecha           string // ISO AAAA-MM-DD
	FechaFormato    string // formato de fecha reconocido
	CadenaComercial string
	Giro            string
	NombreComercial string
	Direccion       string
	Estado          string
	Municipio       string
	Latitud         string
	Longitud        string
}

// RowError es un renglón rechazado con su motivo. Nunca se pierde un renglón sin
// explicación: leídos = guardados + rechazados.
type RowError struct {
	Line   int // línea física (1 = encabezado)
	Row    int // ordinal del registro de datos (1 = el primero tras el encabezado)
	Reason string
}

func (e *RowError) Error() string { return fmt.Sprintf("renglón %d: %s", e.Line, e.Reason) }

// Reader lee un CSV de PROFECO en flujo: un registro a la vez, sin cargar el
// archivo en memoria.
type Reader struct {
	cr    *csv.Reader
	idx   [15]int // posición de cada columna del diccionario en el archivo
	ncols int
	Read  int // renglones leídos de datos (sin encabezado), buenos o rechazados
}

// NewReader lee el encabezado (con o sin BOM, columnas en cualquier orden y sin
// distinguir mayúsculas) y deja listo el flujo. Devuelve ErrProfecoHeader si
// falta alguna de las 15 columnas.
func NewReader(r io.Reader) (*Reader, error) {
	br := bufio.NewReaderSize(r, 1<<20)
	if b, err := br.Peek(3); err == nil && bytes.Equal(b, []byte{0xEF, 0xBB, 0xBF}) {
		_, _ = br.Discard(3)
	}
	cr := csv.NewReader(br)
	cr.FieldsPerRecord = -1
	cr.LazyQuotes = true
	cr.ReuseRecord = true
	head, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrProfecoHeader, err)
	}
	pos := map[string]int{}
	for i, h := range head {
		name := strings.ToLower(cleanText(decodeText(h)))
		if _, dup := pos[name]; dup {
			return nil, fmt.Errorf("%w: columna repetida %q", ErrProfecoHeader, name)
		}
		pos[name] = i
	}
	rd := &Reader{cr: cr, ncols: len(head)}
	var missing []string
	for i, c := range Columns {
		p, ok := pos[c]
		if !ok {
			missing = append(missing, c)
		}
		rd.idx[i] = p
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%w: faltan %s", ErrProfecoHeader, strings.Join(missing, ", "))
	}
	return rd, nil
}

// Next devuelve el siguiente registro normalizado, *RowError si el renglón se
// rechaza (el flujo sigue en el siguiente) o io.EOF al terminar.
func (rd *Reader) Next() (Record, error) {
	row, err := rd.cr.Read()
	if errors.Is(err, io.EOF) {
		return Record{}, io.EOF
	}
	rd.Read++
	if err != nil {
		var pe *csv.ParseError
		if errors.As(err, &pe) {
			return Record{}, &RowError{Line: pe.StartLine, Row: rd.Read, Reason: "CSV mal formado"} // texto cerrado, sin el detalle del lector (§11)
		}
		return Record{}, err
	}
	line, _ := rd.cr.FieldPos(0)
	if len(row) != rd.ncols {
		return Record{}, &RowError{Line: line, Row: rd.Read, Reason: fmt.Sprintf("columnas: %d en vez de %d", len(row), rd.ncols)}
	}
	get := func(i int) string { return cleanText(decodeText(row[rd.idx[i]])) }
	rec := Record{
		Line: line, Row: rd.Read, Producto: get(0), Presentacion: get(1), Marca: get(2), Categoria: get(3),
		Catalogo: get(4), PrecioTexto: get(5), FechaTexto: get(6), CadenaComercial: get(7),
		Giro: get(8), NombreComercial: get(9), Direccion: get(10), Estado: get(11),
		Municipio: get(12), Latitud: cleanCoord(get(13)), Longitud: cleanCoord(get(14)),
	}
	if rec.Producto == "" {
		return Record{}, &RowError{Line: line, Row: rd.Read, Reason: "producto vacío"}
	}
	if rec.CadenaComercial == "" {
		return Record{}, &RowError{Line: line, Row: rd.Read, Reason: "cadena comercial vacía"}
	}
	if rec.PrecioCents, err = ParsePriceCents(rec.PrecioTexto); err != nil {
		return Record{}, &RowError{Line: line, Row: rd.Read, Reason: err.Error()}
	}
	if rec.Fecha, rec.FechaFormato, err = ParseDate(rec.FechaTexto); err != nil {
		return Record{}, &RowError{Line: line, Row: rd.Read, Reason: err.Error()}
	}
	return rec, nil
}

// countingReader cuenta los saltos de línea que pasan, para compararlos con
// `wc -l` en el reporte de cobertura.
type countingReader struct {
	r     io.Reader
	lines int64
	bytes int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.bytes += int64(n)
	c.lines += int64(bytes.Count(p[:n], []byte{'\n'}))
	return n, err
}

// File es un CSV abierto con su lector y sus contadores de bytes y líneas.
type File struct {
	*Reader
	f   *os.File
	cnt *countingReader
	// Size es el tamaño en bytes del archivo.
	Size int64
}

// Lines devuelve cuántos saltos de línea se han leído (equivale a `wc -l` al
// terminar el archivo).
func (f *File) Lines() int64 { return f.cnt.lines }

// Close cierra el archivo.
func (f *File) Close() error { return f.f.Close() }

// OpenFile abre un CSV de PROFECO. maxBytes > 0 es el límite por archivo
// (`price.profeco_import_max_gb`, K6): si lo supera devuelve ErrFileTooLarge sin
// leerlo.
func OpenFile(path string, maxBytes int64) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if maxBytes > 0 && st.Size() > maxBytes {
		f.Close()
		return nil, fmt.Errorf("%w: %d bytes > %d", ErrFileTooLarge, st.Size(), maxBytes)
	}
	cnt := &countingReader{r: f}
	rd, err := NewReader(cnt)
	if err != nil {
		f.Close()
		return nil, err
	}
	return &File{Reader: rd, f: f, cnt: cnt, Size: st.Size()}, nil
}
