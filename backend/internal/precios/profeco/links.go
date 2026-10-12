package profeco

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Mode decide qué ligas cuentan (S5-interfaz §11). En producción solo las
// confirmadas por una persona; en revisión (pruebas, fixtures, medición de
// WF-1) también las propuestas, contadas aparte.
type Mode string

const (
	ModeProduction Mode = "production"
	ModeReview     Mode = "review"
)

// Estados de una liga.
const (
	LinkProposed  = "propuesto"
	LinkConfirmed = "confirmado"
)

// LinkColumns son las columnas de datos/semilla/profeco/vinculos.csv.
var LinkColumns = []string{"producto", "presentacion", "marca", "presentacion_id", "estado", "confirmo", "fecha_confirmacion"}

// ErrLinks indica un vinculos.csv inválido.
var ErrLinks = errors.New("profeco: vinculos.csv inválido")

// Link liga un producto de PROFECO con una presentación del catálogo.
type Link struct {
	PresentationID int64
	Status         string // propuesto | confirmado
	ConfirmedBy    string // axel | ximena, solo si confirmado
	ConfirmedOn    string // ISO, solo si confirmado
}

// Links es el índice de ligas por (producto, presentación, marca) comparados con Fold.
type Links struct {
	byKey map[string]Link
}

// LinkKey es la llave de comparación de una liga o de un registro.
func LinkKey(producto, presentacion, marca string) string {
	return Fold(producto) + "|" + Fold(presentacion) + "|" + Fold(marca)
}

// Lookup busca la liga de un registro.
func (l *Links) Lookup(rec *Record) (Link, bool) {
	if l == nil {
		return Link{}, false
	}
	k, ok := l.byKey[LinkKey(rec.Producto, rec.Presentacion, rec.Marca)]
	return k, ok
}

// Len es el número de ligas.
func (l *Links) Len() int {
	if l == nil {
		return 0
	}
	return len(l.byKey)
}

// ReadLinks lee un vinculos.csv: estado «propuesto» o «confirmado»; si está
// propuesto, confirmo y fecha_confirmacion van vacíos; si está confirmado son
// obligatorios (confirmo = axel o ximena; fecha ISO).
func ReadLinks(r io.Reader) (*Links, error) {
	br := bufio.NewReader(r)
	if b, err := br.Peek(3); err == nil && string(b) == "\xEF\xBB\xBF" {
		_, _ = br.Discard(3)
	}
	cr := csv.NewReader(br)
	cr.FieldsPerRecord = -1
	head, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: sin encabezado", ErrLinks)
	}
	pos := map[string]int{}
	for i, h := range head {
		pos[strings.ToLower(cleanText(decodeText(h)))] = i
	}
	for _, c := range LinkColumns {
		if _, ok := pos[c]; !ok {
			return nil, fmt.Errorf("%w: falta la columna %s", ErrLinks, c)
		}
	}
	out := &Links{byKey: map[string]Link{}}
	for {
		row, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		line, _ := cr.FieldPos(0)
		if err != nil {
			return nil, fmt.Errorf("%w: línea %d: %v", ErrLinks, line, err)
		}
		if len(row) != len(head) {
			return nil, fmt.Errorf("%w: línea %d: columnas %d en vez de %d", ErrLinks, line, len(row), len(head))
		}
		get := func(c string) string { return cleanText(decodeText(row[pos[c]])) }
		id, err := strconv.ParseInt(get("presentacion_id"), 10, 64)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("%w: línea %d: presentacion_id inválido", ErrLinks, line)
		}
		if get("producto") == "" {
			return nil, fmt.Errorf("%w: línea %d: producto vacío", ErrLinks, line)
		}
		lk := Link{PresentationID: id, Status: strings.ToLower(get("estado")), ConfirmedBy: strings.ToLower(get("confirmo")), ConfirmedOn: get("fecha_confirmacion")}
		switch lk.Status {
		case LinkProposed:
			if lk.ConfirmedBy != "" || lk.ConfirmedOn != "" {
				return nil, fmt.Errorf("%w: línea %d: una liga propuesta no lleva confirmo ni fecha", ErrLinks, line)
			}
		case LinkConfirmed:
			if lk.ConfirmedBy != "axel" && lk.ConfirmedBy != "ximena" {
				return nil, fmt.Errorf("%w: línea %d: confirmo debe ser axel o ximena", ErrLinks, line)
			}
			if _, e := time.Parse("2006-01-02", lk.ConfirmedOn); e != nil {
				return nil, fmt.Errorf("%w: línea %d: fecha_confirmacion debe ser ISO", ErrLinks, line)
			}
		default:
			return nil, fmt.Errorf("%w: línea %d: estado debe ser propuesto o confirmado", ErrLinks, line)
		}
		k := LinkKey(get("producto"), get("presentacion"), get("marca"))
		// Liga repetida (después de normalizar) = error de todo el archivo: PROVISIONAL, bloqueo
		// S5-motor-oraculo-reglas pregunta 8 (§11 no lo dice).
		if _, dup := out.byKey[k]; dup {
			return nil, fmt.Errorf("%w: línea %d: liga repetida", ErrLinks, line)
		}
		out.byKey[k] = lk
	}
}

// ReadLinksFile abre y lee un vinculos.csv.
func ReadLinksFile(path string) (*Links, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadLinks(f)
}
