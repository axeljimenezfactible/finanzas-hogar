// Package profeco carga los datos abiertos de «Quién es Quién en los Precios»
// (PROFECO) como precios de referencia. PROFECO nunca compite en la comparativa
// (K22, K37): aquí solo se leen, se normalizan, se filtran y se ligan al catálogo.
//
// Todo el contenido de los archivos es dato no confiable: nunca se ejecuta, se
// lee en flujo y el dinero se convierte a centavos exactos desde el texto.
package profeco

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Errores de archivo (S5-interfaz §11 y §14).
var (
	ErrProfecoHeader = errors.New("profeco: el encabezado no trae las 15 columnas del diccionario")
	ErrFileTooLarge  = errors.New("profeco: el archivo supera el tamaño máximo permitido")
)

// cp1252High es la tabla Windows-1252 para los bytes 0x80 a 0x9F; el resto de
// los bytes de un byte coincide con Latin-1. Los bytes sin definir se dejan en
// su punto de código C1.
var cp1252High = [32]rune{
	0x20AC, 0x0081, 0x201A, 0x0192, 0x201E, 0x2026, 0x2020, 0x2021,
	0x02C6, 0x2030, 0x0160, 0x2039, 0x0152, 0x008D, 0x017D, 0x008F,
	0x0090, 0x2018, 0x2019, 0x201C, 0x201D, 0x2022, 0x2013, 0x2014,
	0x02DC, 0x2122, 0x0161, 0x203A, 0x0153, 0x009D, 0x017E, 0x0178,
}

// decodeText devuelve s como UTF-8. Si los bytes ya son UTF-8 válido no cambia
// nada; si no, se leen como Windows-1252 (así viene QQP_metadatos.csv).
func decodeText(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c < 0x80:
			b.WriteByte(c)
		case c >= 0x80 && c <= 0x9F:
			b.WriteRune(cp1252High[c-0x80])
		default:
			b.WriteRune(rune(c))
		}
	}
	return b.String()
}

// cleanText recorta y junta los espacios (también saltos de línea dentro de
// comillas) en uno solo.
func cleanText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

var foldAccents = map[rune]rune{
	'á': 'a', 'à': 'a', 'ä': 'a', 'â': 'a', 'é': 'e', 'è': 'e', 'ë': 'e', 'ê': 'e',
	'í': 'i', 'ì': 'i', 'ï': 'i', 'î': 'i', 'ó': 'o', 'ò': 'o', 'ö': 'o', 'ô': 'o',
	'ú': 'u', 'ù': 'u', 'ü': 'u', 'û': 'u', 'ñ': 'n',
}

// Fold deja un nombre listo para comparar: minúsculas, sin acentos, sin guion,
// apóstrofo ni punto («Wal-mart» = «walmart», «Sam's Club» = «sams club») y con
// los demás signos convertidos en un espacio.
func Fold(s string) string {
	var b strings.Builder
	space := false
	for _, r := range strings.ToLower(decodeText(s)) {
		if a, ok := foldAccents[r]; ok {
			r = a
		}
		switch {
		case r == '-' || r == '\'' || r == '’' || r == '.' || r == '`':
			continue
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if space && b.Len() > 0 {
				b.WriteByte(' ')
			}
			space = false
			b.WriteRune(r)
		default:
			space = true
		}
	}
	return b.String()
}
