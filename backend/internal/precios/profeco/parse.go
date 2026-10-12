package profeco

import (
	"errors"
	"strings"
	"time"
)

// ParsePriceCents convierte el texto del precio (Número 18,2) a centavos
// exactos, sin flotantes. Acepta «28», «34.5», «$1,234.50». Rechaza vacío,
// negativos, más de dos decimales, texto y valores que desbordan int64.
func ParsePriceCents(s string) (int64, error) {
	t := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "$"))
	if t == "" {
		return 0, errors.New("precio vacío")
	}
	if strings.Contains(t, ",") {
		// Coma de miles solo si cada grupo posterior trae 3 dígitos.
		parts := strings.Split(strings.SplitN(t, ".", 2)[0], ",")
		for i, p := range parts {
			if (i > 0 && len(p) != 3) || p == "" {
				return 0, errors.New("precio con comas mal puestas")
			}
		}
		t = strings.ReplaceAll(t, ",", "")
	}
	intPart, frac, _ := strings.Cut(t, ".")
	if intPart == "" || len(intPart) > 15 {
		return 0, errors.New("precio ilegible")
	}
	if len(frac) > 2 {
		return 0, errors.New("precio con más de dos decimales")
	}
	if strings.Contains(t, ".") && frac == "" {
		return 0, errors.New("precio ilegible")
	}
	var cents int64
	for _, r := range intPart + (frac + "00")[:2] {
		if r < '0' || r > '9' {
			return 0, errors.New("precio ilegible")
		}
		cents = cents*10 + int64(r-'0')
	}
	if cents <= 0 {
		return 0, errors.New("precio no positivo")
	}
	return cents, nil
}

// dateLayouts son los formatos de fecha que el cargador acepta. El diccionario
// solo dice «Datetime»; el archivo 2026 trae AAAA/MM/DD. Con año al final se
// lee día/mes/año (convención mexicana). Cualquier otro se rechaza con motivo.
var dateLayouts = []string{
	"2006/01/02",
	"2006-01-02",
	"2006/1/2",
	"2006-1-2",
	"02/01/2006",
	"2/1/2006",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	"20060102",
}

// ParseDate devuelve la fecha ISO (AAAA-MM-DD) y el formato que reconoció.
func ParseDate(s string) (iso, layout string, err error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return "", "", errors.New("fecha vacía")
	}
	for _, l := range dateLayouts {
		if d, e := time.Parse(l, t); e == nil {
			return d.Format("2006-01-02"), l, nil
		}
	}
	return "", "", errors.New("fecha ilegible")
}

// cleanCoord valida una coordenada como texto decimal (sin flotantes) y la
// devuelve sin cambios; si no es decimal devuelve vacío.
func cleanCoord(s string) string {
	t := strings.TrimSpace(s)
	if t == "" {
		return ""
	}
	body := strings.TrimPrefix(t, "-")
	intPart, frac, hasDot := strings.Cut(body, ".")
	if intPart == "" || len(intPart) > 3 || (hasDot && frac == "") {
		return ""
	}
	for _, r := range intPart + frac {
		if r < '0' || r > '9' {
			return ""
		}
	}
	return t
}
