package ocr

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// REGLA CONGELADA DE MEDICIÓN (S5-06, K24; se escribe antes de medir y no se
// afloja después):
//
//  1. Renglón esperado = una fila de la transcripción a mano
//     (datos/semilla/tickets-esperado/NN.csv: descripcion,cantidad,importe).
//     La descripción se transcribe tal como está IMPRESA en el ticket, sin
//     código de barras ni letra de IVA; la cantidad solo si el ticket la
//     imprime (vacía si no); el importe tal como está impreso.
//  2. Renglón correcto = existe un renglón leído, todavía no usado, con:
//     descripción igual tras normalizar (mayúsculas, sin acentos, espacios
//     colapsados, sin puntuación final) Y cantidad exacta (ambas vacías, o
//     iguales en milésimas) Y importe exacto en centavos.
//  3. Los renglones de resumen (TOTAL, IVA, CAMBIO...) y los sin importe no
//     compiten por una coincidencia. Los renglones leídos de más no restan al
//     porcentaje, pero se reportan como «sobrantes».
//  4. Porcentaje = correctos / esperados, sumando los 5 tickets; cumple con
//     correctos * 100 >= 80 * esperados.
//  5. Sin 5 imágenes 01..05 con su transcripción: «sin medir», nunca cumplido.

// MinPassPercent es el umbral de S5-06.
const MinPassPercent = 80

// RequiredTickets es el número de fotos reales que exige K24.
const RequiredTickets = 5

// Expected es un renglón de la transcripción a mano.
type Expected struct {
	Descripcion   string
	CantidadMilli *int64
	ImporteCents  int64
}

var accents = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n",
	"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U", "Ü", "U", "Ñ", "N",
)

// NormalizeDescription aplica la regla 2 de la medición.
func NormalizeDescription(s string) string {
	s = strings.ToUpper(accents.Replace(s))
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimRight(s, " .,;:-")
}

// ReadExpectedCSV lee la transcripción: encabezado descripcion,cantidad,importe;
// cantidad vacía = no impresa; decimales con punto.
func ReadExpectedCSV(r io.Reader) ([]Expected, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 3
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ocr: transcripción: %w", err)
	}
	if len(rows) == 0 || strings.ToLower(strings.TrimSpace(rows[0][0])) != "descripcion" {
		return nil, errors.New("ocr: transcripción sin encabezado descripcion,cantidad,importe")
	}
	var out []Expected
	for i, row := range rows[1:] {
		imp, ok := parseFixed(strings.TrimSpace(row[2]), 2)
		if !ok || strings.TrimSpace(row[2]) == "" {
			return nil, fmt.Errorf("ocr: transcripción fila %d: importe %q no válido", i+2, row[2])
		}
		e := Expected{Descripcion: strings.TrimSpace(row[0]), ImporteCents: imp}
		if q := strings.TrimSpace(row[1]); q != "" {
			v, ok := parseFixed(q, 3)
			if !ok {
				return nil, fmt.Errorf("ocr: transcripción fila %d: cantidad %q no válida", i+2, q)
			}
			e.CantidadMilli = &v
		}
		out = append(out, e)
	}
	return out, nil
}

// TicketResult es el desglose de un ticket.
type TicketResult struct {
	Name     string
	Expected int
	Correct  int
	Extra    int      // renglones leídos que no casaron con ninguno esperado
	Misses   []string // descripciones esperadas que no salieron correctas
}

func eqQty(a, b *int64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// Score compara un ticket renglón a renglón (reglas 2 y 3).
func Score(name string, got []Line, exp []Expected) TicketResult {
	res := TicketResult{Name: name, Expected: len(exp)}
	used := make([]bool, len(got))
	for _, e := range exp {
		want := NormalizeDescription(e.Descripcion)
		hit := false
		for i, g := range got {
			if used[i] || g.EsResumen || g.ImporteCents == nil {
				continue
			}
			if NormalizeDescription(g.Descripcion) == want && eqQty(g.CantidadMilli, e.CantidadMilli) && *g.ImporteCents == e.ImporteCents {
				used[i], hit = true, true
				break
			}
		}
		if hit {
			res.Correct++
		} else {
			res.Misses = append(res.Misses, e.Descripcion)
		}
	}
	for i, g := range got {
		if !used[i] && !g.EsResumen && g.ImporteCents != nil {
			res.Extra++
		}
	}
	return res
}

// Recognizer lo cumple *Engine; en las pruebas, un simulado.
type Recognizer interface {
	Recognize(ctx context.Context, imagen string) ([]Line, error)
}

// MeasureOptions configura Measure.
type MeasureOptions struct {
	TicketsDir  string // fotos NN.jpg|.jpeg|.png
	ExpectedDir string // NN.csv (mismo nombre base que la foto)
	// Synthetic mide imágenes de prueba generadas: no exige 5 y el informe
	// dice que NO cuenta para S5-06.
	Synthetic bool
}

// Report es el resultado de la medición.
type Report struct {
	Measured  bool   // false = «sin medir»
	Reason    string // por qué quedó sin medir
	Synthetic bool
	Tickets   []TicketResult
	Expected  int
	Correct   int
	PctTenths int // porcentaje en décimas (833 = 83.3 %), HALF_UP
	Passes    bool
}

func listImages(dir string) ([]string, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		switch strings.ToLower(filepath.Ext(e.Name())) {
		case ".jpg", ".jpeg", ".png":
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// Measure aplica la regla congelada. Nunca declara «cumplido» sin las 5 fotos.
func Measure(ctx context.Context, rec Recognizer, o MeasureOptions) (Report, error) {
	rep := Report{Synthetic: o.Synthetic}
	imgs, err := listImages(o.TicketsDir)
	if err != nil {
		return rep, err
	}
	if o.Synthetic {
		if len(imgs) == 0 {
			rep.Reason = "sin imágenes sintéticas"
			return rep, nil
		}
	} else if len(imgs) != RequiredTickets {
		rep.Reason = fmt.Sprintf("hay %d imágenes y se necesitan %d: faltan las fotos de Axel", len(imgs), RequiredTickets)
		return rep, nil
	}
	type pair struct {
		img string
		exp []Expected
	}
	var pairs []pair
	for _, n := range imgs {
		stem := strings.TrimSuffix(n, filepath.Ext(n))
		f, err := os.Open(filepath.Join(o.ExpectedDir, stem+".csv"))
		if err != nil {
			rep.Reason = fmt.Sprintf("falta la transcripción a mano de %s en %s", n, o.ExpectedDir)
			return rep, nil
		}
		exp, err := ReadExpectedCSV(f)
		f.Close()
		if err != nil {
			return rep, fmt.Errorf("%s: %w", stem, err)
		}
		pairs = append(pairs, pair{filepath.Join(o.TicketsDir, n), exp})
	}
	for _, p := range pairs {
		got, err := rec.Recognize(ctx, p.img)
		if err != nil {
			return rep, fmt.Errorf("%s: %w", filepath.Base(p.img), err)
		}
		tr := Score(strings.TrimSuffix(filepath.Base(p.img), filepath.Ext(p.img)), got, p.exp)
		rep.Tickets = append(rep.Tickets, tr)
		rep.Expected += tr.Expected
		rep.Correct += tr.Correct
	}
	if rep.Expected == 0 {
		rep.Reason = "las transcripciones no tienen renglones"
		return rep, nil
	}
	rep.Measured = true
	rep.PctTenths = (rep.Correct*2000 + rep.Expected) / (2 * rep.Expected)
	rep.Passes = rep.Correct*100 >= MinPassPercent*rep.Expected
	return rep, nil
}

// Format escribe el informe para el reporte (tabla ticket × correctos/totales).
func (r Report) Format() string {
	var b strings.Builder
	if r.Synthetic {
		b.WriteString("SINTETICO: prueba de herramienta, NO cuenta para S5-06\n")
	}
	if !r.Measured {
		fmt.Fprintf(&b, "S5-06 sin medir: %s\n", r.Reason)
		return b.String()
	}
	b.WriteString("ticket | correctos/esperados | sobrantes\n")
	for _, t := range r.Tickets {
		fmt.Fprintf(&b, "%s | %d/%d | %d\n", t.Name, t.Correct, t.Expected, t.Extra)
		for _, m := range t.Misses {
			fmt.Fprintf(&b, "  falla: %s\n", m)
		}
	}
	verdict := "NO cumple (< 80 %)"
	if r.Passes {
		verdict = "cumple (>= 80 %)"
	}
	fmt.Fprintf(&b, "TOTAL %d/%d = %d.%d %% — %s\n", r.Correct, r.Expected, r.PctTenths/10, r.PctTenths%10, verdict)
	return b.String()
}
