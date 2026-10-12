package dominio

// Reglas compartidas del tablero de Inicio (K54, C1 v1.5) entre el paquete tablero y las tarjetas de Inicio
// (getDashboardMonthVsPlan). Agente api-tablero (WF-3c).

// MonthHasData dice si el sistema tiene datos de gasto del hogar en un mes, a partir de las dos cifras que el motor S2
// da para ese mes: el devengado (real de M-01) y el efectivo (cash_household_cents de M-11). C1 v1.5 pide omitir las
// cifras de un mes «sin datos» porque «sin datos no es cero» (M-01 ampliada y M-21), pero no define cuándo un mes no
// tiene datos. PROVISIONAL (bloqueo docs/estado/bloqueos/S6-tablero-sin-datos.md, opción a): un mes con las dos cifras
// en cero no tiene datos. Así la usan getDashboardMonthVsPlan (mes anterior) y getDashboardMonthlyOverview.
func MonthHasData(accruedCents, cashCents int64) bool { return accruedCents != 0 || cashCents != 0 }
