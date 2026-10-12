package ocr

import (
	"regexp"
	"strings"
)

// Profile ajusta el análisis a la forma de un ticket. Ninguno de los perfiles
// está validado contra tickets reales de Bodega Aurrera ni de Walmart: la
// carpeta datos/semilla/tickets/ está vacía (K24) y el contrato prohíbe
// inventar formatos. Cuando lleguen las 5 fotos se afina aquí.
type Profile struct {
	Name string
	// StripBarcodes quita de la descripción los códigos de barras (8 a 14
	// dígitos seguidos), frecuentes antes del nombre en algunos tickets.
	StripBarcodes bool
}

var (
	// ProfileGeneric sigue al pie de la letra S5-interfaz §12.
	ProfileGeneric = Profile{Name: "generico"}
	// ProfileWithCode es ProfileGeneric más la limpieza de códigos de barras.
	// HIPÓTESIS sin validar contra un ticket real.
	ProfileWithCode = Profile{Name: "con-codigo", StripBarcodes: true}
)

var (
	// Importe: dos decimales exactos, con $ y miles opcionales. Los OCR
	// confunden O por 0 e I, l, | por 1; fixDigits los corrige solo dentro de
	// un token que ya tiene forma de importe.
	amountRe  = regexp.MustCompile(`^-?\$?(\d{1,3}(,\d{3})+|\d+)\.\d{2}$`)
	amountFix = regexp.MustCompile(`^-?\$?[\dOoIl|,]+\.[\dOoIl|]{2}$`)
	qtyRe     = regexp.MustCompile(`^\d{1,3}([.,]\d{1,3})?$`)
	barcodeRe = regexp.MustCompile(`^\d{8,14}$`)
	summaryRe = regexp.MustCompile(`^(SUB ?TOTAL|TOTAL|IVA|I\.V\.A|I V A|CAMBIO|EFECTIVO|TARJETA|PAGO|ARTICULOS|NO DE ART|AHORRO|REDONDEO)\b`)
	units     = map[string]bool{
		"KG": true, "KGS": true, "G": true, "GR": true, "GRS": true,
		"L": true, "LT": true, "LTS": true, "ML": true,
		"PZA": true, "PZAS": true, "PZ": true, "PZS": true, "U": true, "UN": true,
	}
)

// fixDigits corrige confusiones letra-dígito en un token con forma de importe.
func fixDigits(tok string) string {
	if !amountFix.MatchString(tok) || !strings.ContainsAny(tok, "0123456789") {
		return tok
	}
	r := strings.NewReplacer("O", "0", "o", "0", "I", "1", "l", "1", "|", "1")
	return r.Replace(tok)
}

// parseFixed convierte "1,234.50" o "0.750" a un entero con `decimals`
// decimales exactos (centavos o milésimas). Sin float. Rechaza más decimales
// de los pedidos y más de 12 dígitos enteros.
func parseFixed(s string, decimals int) (int64, bool) {
	s = strings.TrimSpace(s)
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")
	intPart, frac, _ := strings.Cut(s, ".")
	if intPart == "" || len(intPart) > 12 || len(frac) > decimals {
		return 0, false
	}
	for len(frac) < decimals {
		frac += "0"
	}
	var v int64
	for _, c := range intPart + frac {
		if c < '0' || c > '9' {
			return 0, false
		}
		v = v*10 + int64(c-'0')
	}
	if neg {
		v = -v
	}
	return v, true
}

// parseQty acepta "2", "0.750" y "0,750" (coma decimal solo en cantidades,
// que nunca llevan miles).
func parseQty(tok string) (int64, bool) {
	if !qtyRe.MatchString(tok) {
		return 0, false
	}
	return parseFixed(strings.Replace(tok, ",", ".", 1), 3)
}

// ParseLine interpreta el texto de un renglón (S5-interfaz §12):
//   - importe  = último token con dos decimales (en centavos); lo que venga
//     después (letras de IVA, asteriscos) se descarta;
//   - cantidad = número inicial (entero o hasta 3 decimales) con su unidad
//     opcional, en milésimas; solo si queda texto de descripción después;
//   - descripcion = el resto.
//
// La confianza la fija quien llama (ParseTSV); aquí queda en 0.
func ParseLine(text string, p Profile) Line {
	l := Line{Texto: strings.TrimSpace(text)}
	toks := strings.Fields(text)

	amountAt := -1
	for i := len(toks) - 1; i >= 0; i-- {
		if t := fixDigits(toks[i]); amountRe.MatchString(t) {
			if v, ok := parseFixed(t, 2); ok {
				l.ImporteCents = &v
				amountAt = i
				break
			}
		}
	}
	head := toks
	if amountAt >= 0 {
		head = toks[:amountAt]
	}

	if len(head) >= 2 {
		if q, ok := parseQty(head[0]); ok {
			rest := head[1:]
			unit := ""
			if len(rest) >= 2 && units[strings.ToUpper(rest[0])] {
				unit = strings.ToUpper(rest[0])
				rest = rest[1:]
			}
			if hasLetter(rest) {
				l.CantidadMilli = &q
				l.Unidad = unit
				head = rest
			}
		}
	}

	if p.StripBarcodes {
		kept := head[:0:0]
		for _, t := range head {
			if !barcodeRe.MatchString(t) {
				kept = append(kept, t)
			}
		}
		head = kept
	}
	l.Descripcion = strings.Trim(strings.Join(head, " "), " |_:.,;-")
	l.EsResumen = summaryRe.MatchString(NormalizeDescription(l.Descripcion))
	return l
}

func hasLetter(toks []string) bool {
	for _, t := range toks {
		for _, c := range t {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c >= 0xC0 {
				return true
			}
		}
	}
	return false
}
