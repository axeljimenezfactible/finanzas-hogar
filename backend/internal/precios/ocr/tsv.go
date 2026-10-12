package ocr

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type lineAcc struct {
	words []string
	conf  []float64
}

// ParseTSV agrupa las palabras del TSV de tesseract por página, bloque,
// párrafo y línea (columnas 2 a 5), en el orden en que aparecen, y analiza
// cada renglón con ParseLine. La confianza de un renglón es el promedio de la
// de sus palabras, redondeado HALF_UP (0 a 100).
func ParseTSV(r io.Reader, p Profile) ([]Line, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	var order []string
	acc := map[string]*lineAcc{}
	header := true
	for sc.Scan() {
		row := strings.TrimRight(sc.Text(), "\r")
		if row == "" {
			continue
		}
		cols := strings.Split(row, "\t")
		if header {
			header = false
			if len(cols) < 12 || cols[0] != "level" {
				return nil, fmt.Errorf("ocr: encabezado TSV inesperado: %q", truncate(row, 60))
			}
			continue
		}
		if len(cols) < 12 || cols[0] != "5" {
			continue // solo nivel 5 = palabra
		}
		text := strings.TrimSpace(cols[11])
		conf, err := strconv.ParseFloat(cols[10], 64)
		if text == "" || err != nil || conf < 0 {
			continue
		}
		key := strings.Join(cols[1:5], "/")
		a, ok := acc[key]
		if !ok {
			a = &lineAcc{}
			acc[key] = a
			order = append(order, key)
		}
		a.words = append(a.words, text)
		a.conf = append(a.conf, conf)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("ocr: leyendo TSV: %w", err)
	}
	lines := make([]Line, 0, len(order))
	for _, k := range order {
		a := acc[k]
		l := ParseLine(strings.Join(a.words, " "), p)
		var sum float64
		for _, c := range a.conf {
			sum += c
		}
		l.Confianza = int(sum/float64(len(a.conf)) + 0.5)
		lines = append(lines, l)
	}
	return lines, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
