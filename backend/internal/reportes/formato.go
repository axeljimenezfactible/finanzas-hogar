// Package reportes genera el reporte mensual y el estado de cuenta de personas externas en PDF (LaTeX con lualatex),
// los guarda versionados y envía el mensual por Telegram (S7-01…S7-03, S7-10…S7-13). Las cifras salen solo del motor
// S2 (motor_fin vía dominio.FinEngine) y de las mismas tablas que lee la API; las plantillas no traen valores fijos.
package reportes

import (
	"strings"
	"time"

	"finanzas-hogar/internal/dominio"
)

var (
	mesesCortos = [...]string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}
	mesesLargos = [...]string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	diasCortos  = [...]string{"Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"}
)

// Pesos formatea centavos como en pantalla: $1,234.56 (S7-12; mismo formato que dominio.FormatPesos).
func Pesos(cents int64) string { return dominio.FormatPesos(cents) }

// Fecha da 10-oct-2026 (decisión B, S7-12).
func Fecha(d dominio.Date) string {
	return itoa(d.Day) + "-" + mesesCortos[d.Month-1] + "-" + itoa(d.Year)
}

// FechaLarga da Vie 30 oct 2026 (con día de la semana, donde importa: fechas de pago).
func FechaLarga(d dominio.Date) string {
	return diasCortos[d.Weekday()] + " " + itoa(d.Day) + " " + mesesCortos[d.Month-1] + " " + itoa(d.Year)
}

// MesLargo da «octubre 2026».
func MesLargo(y int, m time.Month) string { return mesesLargos[m-1] + " " + itoa(y) }

// FechaAPI convierte AAAA-MM-DD de la API al formato visible; una fecha inválida se rechaza (nunca se imprime ISO).
func FechaAPI(s string, larga bool) (string, error) {
	d, err := dominio.ParseDate(s)
	if err != nil {
		return "", err
	}
	if larga {
		return FechaLarga(d), nil
	}
	return Fecha(d), nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

// texReplacer escapa los caracteres especiales de LaTeX y lleva a comandos los signos que la codificación T1 no
// tiene en la misma posición que Unicode (¿, ¡, rayas, comillas), para que el texto salga igual en el PDF.
var texReplacer = strings.NewReplacer(
	`\`, `\textbackslash{}`, `{`, `\{`, `}`, `\}`, `$`, `\$`, `&`, `\&`, `#`, `\#`, `^`, `\textasciicircum{}`,
	`_`, `\_`, `%`, `\%`, `~`, `\textasciitilde{}`, `<`, `\textless{}`, `>`, `\textgreater{}`, `|`, `\textbar{}`,
	`¿`, `\textquestiondown{}`, `¡`, `\textexclamdown{}`, `—`, `\textemdash{}`, `–`, `\textendash{}`,
	`·`, `\textperiodcentered{}`, `“`, `\textquotedblleft{}`, `”`, `\textquotedblright{}`, `«`, `\guillemotleft{}`,
	`»`, `\guillemotright{}`, "\r", "", "\t", " ",
)

// Tex escapa texto para las plantillas; quita caracteres de control.
func Tex(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == '\n' {
			return ' '
		}
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s)
	return texReplacer.Replace(s)
}
