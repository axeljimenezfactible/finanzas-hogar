package profeco

import (
	"errors"
	"strings"
	"time"
)

// ParsePriceCents convierte el texto del precio (Número 18,2) a centavos exactos, sin
// flotantes, con la gramática de §11 paso 3 (regla técnica, H-r3-03, Q2-04 = A): se quitan
// los espacios de las orillas y un «$» inicial; dígitos con coma de miles opcional (después del
// primer grupo, cada grupo lleva exactamente 3 dígitos: «1,021.00» = 102100) y, opcional, un
// punto con 1 o 2 decimales. Los motivos van en este orden y son la lista cerrada de §11:
// «precio vacío», «precio con comas mal puestas», «precio con más de dos decimales»,
// «precio ilegible» (otro carácter, punto sin decimales o más de 15 dígitos enteros, que
// desbordarían int64) y «precio no positivo» (0.00 nunca se carga).
// PROVISIONAL (bloqueo S5-motor-oraculo-reglas pregunta 7; §11 paso 3 no los fija): un espacio en medio
// da «precio ilegible» (S-1), más de 15 dígitos enteros da «precio ilegible» (S-2) y «más de dos
// decimales» cuenta todo lo que sigue al primer punto, aunque no sean dígitos: «1.2.3» (S-3).
func ParsePriceCents(s string) (int64, error) {
	t := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "$"))
	if t == "" {
		return 0, errors.New("precio vacío")
	}
	intPart, frac, hasDot := strings.Cut(t, ".")
	if strings.Contains(t, ",") {
		if strings.Contains(frac, ",") {
			return 0, errors.New("precio con comas mal puestas")
		}
		for i, p := range strings.Split(intPart, ",") {
			if p == "" || (i > 0 && len(p) != 3) {
				return 0, errors.New("precio con comas mal puestas")
			}
		}
		intPart = strings.ReplaceAll(intPart, ",", "")
	}
	if len(frac) > 2 {
		return 0, errors.New("precio con más de dos decimales")
	}
	if intPart == "" || len(intPart) > 15 || (hasDot && frac == "") {
		return 0, errors.New("precio ilegible")
	}
	var cents int64
	for _, r := range intPart + (frac + "00")[:2] {
		if r < '0' || r > '9' {
			return 0, errors.New("precio ilegible")
		}
		cents = cents*10 + int64(r-'0')
	}
	if cents == 0 {
		return 0, errors.New("precio no positivo")
	}
	return cents, nil
}

// dateLayouts son los formatos de fecha que el cargador acepta. El diccionario
// solo dice «Datetime»; el archivo 2026 trae AAAA/MM/DD. Con año al final se
// lee día/mes/año (convención mexicana). Cualquier otro se rechaza con motivo.
// PROVISIONAL (bloqueo S5-motor-oraculo-reglas pregunta 4): la lista no está en §11 y difiere de la del oráculo.
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
