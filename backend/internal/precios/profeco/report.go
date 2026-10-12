package profeco

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// money-free: todos los conteos son enteros.

func thousands(n int64) string {
	s := fmt.Sprint(n)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	var out []byte
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, s[i])
	}
	if neg {
		return "-" + string(out)
	}
	return string(out)
}

func cell(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "|", "/"), "\n", " ")
}

func monthsSpan(a *Agg) string {
	m := make([]string, 0, len(a.Months))
	for k := range a.Months {
		m = append(m, k)
	}
	sort.Strings(m)
	if len(m) == 0 {
		return "0"
	}
	return fmt.Sprintf("%d (%s a %s)", len(m), m[0], m[len(m)-1])
}

// WriteMarkdown escribe las secciones de conteos del reporte de cobertura. Es
// determinista: mismos archivos, mismo texto (los tiempos no se incluyen).
func (s *Survey) WriteMarkdown(w io.Writer) {
	p := func(f string, a ...any) { fmt.Fprintf(w, f, a...) }

	p("### Archivos leídos: conteo contra `wc -l`\n\n")
	p("Cada archivo trae un encabezado; «diferencia» = `wc -l` − 1 − registros leídos. Cero significa que cada salto de línea es un registro; si no es cero, hay saltos de línea dentro de comillas (se cuentan como un solo registro con el lector CSV).\n\n")
	p("| archivo | bytes | `wc -l` | registros leídos | rechazados | diferencia | primera fecha | última fecha |\n|---|---:|---:|---:|---:|---:|---|---|\n")
	var tb, tl, tr, trej int64
	for _, f := range s.Files {
		p("| %s | %s | %s | %s | %s | %d | %s | %s |\n", f.Name, thousands(f.Size), thousands(f.Lines), thousands(f.Read), thousands(f.Rejected), f.Lines-1-f.Read, f.MinDate, f.MaxDate)
		tb += f.Size
		tl += f.Lines
		tr += f.Read
		trej += f.Rejected
	}
	p("| **total** | %s | %s | %s | %s | %d | %s | %s |\n\n", thousands(tb), thousands(tl), thousands(tr), thousands(trej), tl-int64(len(s.Files))-tr, s.Total.MinDate, s.Total.MaxDate)
	p("Registros leídos = guardados + rechazados: %s = %s + %s.\n\n", thousands(tr), thousands(tr-trej), thousands(trej))

	p("**Formatos de fecha reconocidos** (el diccionario solo dice «Datetime»):")
	fk := make([]string, 0, len(s.Formats))
	for k := range s.Formats {
		fk = append(fk, k)
	}
	sort.Strings(fk)
	for _, k := range fk {
		p(" `%s` = %s registros;", k, thousands(s.Formats[k]))
	}
	p("\n\n")
	if len(s.Rejects) == 0 {
		p("**Rechazados con motivo:** ninguno.\n\n")
	} else {
		p("**Rechazados con motivo:**\n\n| motivo | renglones |\n|---|---:|\n")
		rk := make([]string, 0, len(s.Rejects))
		for k := range s.Rejects {
			rk = append(rk, k)
		}
		sort.Strings(rk)
		for _, k := range rk {
			p("| %s | %s |\n", cell(k), thousands(s.Rejects[k]))
		}
		p("\n")
	}

	p("### Totales\n\n- Registros: %s · productos distintos: %s · presentaciones distintas (producto + presentación + marca): %s · sucursales distintas (nombre comercial + dirección + municipio): %s\n- Cadenas comerciales distintas (`cadena_comercial`): %d\n- Meses cubiertos: %s\n- Fecha mínima: %s · fecha máxima: %s\n\n",
		thousands(s.Total.Records), thousands(int64(s.Total.Products())), thousands(int64(s.Total.Presentations())), thousands(int64(s.Total.Branches())), len(s.Chains), monthsSpan(s.Total), s.Total.MinDate, s.Total.MaxDate)

	p("### Cobertura por `cadena_comercial` (todas, de mayor a menor)\n\n")
	p("| cadena_comercial | registros | productos | presentaciones | sucursales | meses | primera | última |\n|---|---:|---:|---:|---:|---|---|---|\n")
	for _, n := range s.sortedChains() {
		a := s.Chains[n]
		mark := ""
		if IsFocusChain(Fold(n)) {
			mark = " ★"
		}
		p("| %s%s | %s | %s | %s | %s | %s | %s | %s |\n", cell(n), mark, thousands(a.Records), thousands(int64(a.Products())), thousands(int64(a.Presentations())), thousands(int64(a.Branches())), monthsSpan(a), a.MinDate, a.MaxDate)
	}
	p("\n★ = cadenas de interés (familia Walmart, familia Bodega Aurrera, Sam's Club).\n\n")

	s.writeByMunicipality(w)
	s.writeNamed(w)
	s.writeWatch(w)
}

func (s *Survey) writeByMunicipality(w io.Writer) {
	p := func(f string, a ...any) { fmt.Fprintf(w, f, a...) }
	keys := make([][3]string, 0, len(s.ChainMun))
	for k := range s.ChainMun {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		if a[1] != b[1] {
			return a[1] < b[1]
		}
		if a[2] != b[2] {
			return a[2] < b[2]
		}
		return a[0] < b[0]
	})
	p("### Cadenas de interés por estado y municipio\n\n")
	if len(keys) == 0 {
		p("Ninguna de las cadenas de interés tiene registros.\n\n")
		return
	}
	p("Las sucursales cercanas a Axel no se conocen: se reporta por municipio y la selección queda en Ajustes › Tiendas y precios.\n\n")
	p("| estado | municipio | cadena_comercial | registros | productos | sucursales | meses | primera | última |\n|---|---|---|---:|---:|---:|---|---|---|\n")
	for _, k := range keys {
		a := s.ChainMun[k]
		p("| %s | %s | %s | %s | %s | %s | %s | %s | %s |\n", cell(k[1]), cell(k[2]), cell(k[0]), thousands(a.Records), thousands(int64(a.Products())), thousands(int64(a.Branches())), monthsSpan(a), a.MinDate, a.MaxDate)
	}
	p("\n")
}

func (s *Survey) writeNamed(w io.Writer) {
	p := func(f string, a ...any) { fmt.Fprintf(w, f, a...) }
	p("### Sucursales de Axel\n\n")
	p("**Regla de empate (escrita antes de contar):** un registro es de la sucursal si la `cadena_comercial` (con `Fold`: minúsculas, sin acentos, sin guion, punto ni apóstrofo) es una de las listadas y la frase aparece como palabras completas en `nombre_comercial` o en `direccion`. Para Sam's Club no se exige cadena: basta que «sams club» aparezca en `cadena_comercial`, `nombre_comercial` o `direccion`. Las «coincidencias parciales» (una sola palabra) se listan aparte solo como pista y no cuentan como registros de la sucursal.\n\n")
	for _, n := range s.Named {
		chains := "cualquiera"
		if len(n.Chains) > 0 {
			chains = strings.Join(n.Chains, " · ")
		}
		p("**%s** — cadenas: %s · frase: «%s» · parcial: «%s»\n\n", n.Label, chains, n.Phrase, n.Near)
		if n.Agg.Records == 0 {
			p("- Resultado: **NO hay registros** de %s con esta regla.\n", n.Label)
		} else {
			p("- Resultado: **SÍ hay registros** de %s: %s registros, %s productos, %s sucursales, meses %s, del %s al %s.\n", n.Label, thousands(n.Agg.Records), thousands(int64(n.Agg.Products())), thousands(int64(n.Agg.Branches())), monthsSpan(n.Agg), n.Agg.MinDate, n.Agg.MaxDate)
			p("- De esos, la frase aparece **en `cadena_comercial` o `nombre_comercial`**: %s registros; **solo en `direccion`** (referencia de ubicación, no el nombre del local): %s registros.\n", thousands(n.Agg.Records-n.DirOnly.Records), thousands(n.DirOnly.Records))
			writeExamples(w, "Sucursales que cumplen la regla", n.Examples, 15)
		}
		if n.NearAgg.Records == 0 {
			p("- Coincidencias parciales («%s»): ninguna.\n\n", n.Near)
		} else {
			p("- Coincidencias parciales («%s»), no cuentan: %s registros en %d sucursales.\n", n.Near, thousands(n.NearAgg.Records), len(n.NearEx))
			writeExamples(w, "Parciales", n.NearEx, 8)
			p("\n")
		}
	}
}

// writeWatch lista, por municipio vigilado, todas las sucursales de las
// cadenas de interés (para que Axel reconozca las suyas).
func (s *Survey) writeWatch(w io.Writer) {
	if len(s.Watch) == 0 {
		return
	}
	keys := make([][2]string, 0, len(s.Watch))
	for k := range s.Watch {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i][0] != keys[j][0] {
			return keys[i][0] < keys[j][0]
		}
		return keys[i][1] < keys[j][1]
	})
	fmt.Fprintf(w, "### Todas las sucursales de las cadenas de interés en los municipios vigilados\n\nMunicipios: %s. Sirve para que Axel reconozca sus tiendas por nombre y dirección; el cargador no adivina cuál es.\n\n", strings.Join(s.WatchMun, ", "))
	for _, k := range keys {
		writeExamples(w, k[0]+" · "+k[1], s.Watch[k], 20)
	}
	fmt.Fprintln(w)
}

func writeExamples(w io.Writer, title string, ex map[string]int64, limit int) {
	keys := make([]string, 0, len(ex))
	for k := range ex {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if ex[keys[i]] != ex[keys[j]] {
			return ex[keys[i]] > ex[keys[j]]
		}
		return keys[i] < keys[j]
	})
	fmt.Fprintf(w, "  - %s (hasta %d de %d, nombre comercial | dirección | municipio):\n", title, limit, len(keys))
	for i, k := range keys {
		if i >= limit {
			break
		}
		fmt.Fprintf(w, "    - %s — %s registros\n", cell(k), thousands(ex[k]))
	}
}
