package avisos

import (
	"fmt"
	"strconv"

	"finanzas-hogar/internal/dominio"
)

var (
	diasCortos  = [...]string{"Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"}
	mesesCortos = [...]string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}
)

// fechaLarga: «Vie 30 oct 2026» (decisión B, S7-12).
func fechaLarga(d dominio.Date) string {
	return fmt.Sprintf("%s %d %s %d", diasCortos[d.Weekday()], d.Day, mesesCortos[d.Month-1], d.Year)
}

// fechaCorta: «10-oct-2026».
func fechaCorta(d dominio.Date) string {
	return fmt.Sprintf("%d-%s-%d", d.Day, mesesCortos[d.Month-1], d.Year)
}

// mes: «oct 2026».
func mes(d dominio.Date) string { return mesesCortos[d.Month-1] + " " + strconv.Itoa(d.Year) }

// pesos: «$1,234.56».
func pesos(c int64) string { return dominio.FormatPesos(c) }

func gastos(n int64) string {
	if n == 1 {
		return "1 gasto"
	}
	return strconv.FormatInt(n, 10) + " gastos"
}
