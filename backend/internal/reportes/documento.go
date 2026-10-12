package reportes

// Cell es una celda del documento. Kind: "" (texto), "money", "date" o "date_long". Src dice de qué respuesta de la API
// sale la cifra («respuesta|ruta.json»): la prueba de cruce (S7-02, S7-03, S7-11) la resuelve contra la API real.
type Cell struct {
	Text string
	Kind string
	Src  string
}

// Section es una tabla con título.
type Section struct {
	Title string
	Note  string
	Spec  string // columnas de tabularx: L (texto) y R (cifra)
	Cols  []string
	Rows  [][]Cell
}

// Doc es lo que la plantilla pinta: todo ya formateado y escapado al pintar.
type Doc struct {
	Title    string
	Subtitle []Cell
	Footer   string
	Note     string
	Sections []Section
}

// Cells devuelve las celdas con cifra o fecha en el orden en que aparecen en el PDF.
func (d *Doc) Cells() []Cell {
	var out []Cell
	for _, c := range d.Subtitle {
		if c.Kind != "" {
			out = append(out, c)
		}
	}
	for _, s := range d.Sections {
		for _, r := range s.Rows {
			for _, c := range r {
				if c.Kind != "" {
					out = append(out, c)
				}
			}
		}
	}
	return out
}

func txt(s string) Cell { return Cell{Text: s} }

func money(src string, v int64) Cell { return Cell{Text: Pesos(v), Kind: "money", Src: src} }

func date(src, iso string, larga bool) Cell {
	t, err := FechaAPI(iso, larga)
	if err != nil {
		return Cell{Text: "fecha inválida"}
	}
	k := "date"
	if larga {
		k = "date_long"
	}
	return Cell{Text: t, Kind: k, Src: src}
}

func statusText(s string) string {
	switch s {
	case "paid":
		return "Pagado"
	case "overdue":
		return "Vencido"
	default:
		return "Pendiente"
	}
}
