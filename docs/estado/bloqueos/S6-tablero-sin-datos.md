# Bloqueo S6-tablero-sin-datos · cuándo un mes «no tiene datos» en el tablero (K54, C1 v1.5)

Abierto el 12-oct-2026 por el agente `api-tablero` (WF-3c). Lo contesta Axel en la **puerta 4**. Mientras tanto rige la opción PROVISIONAL (a), que aplica una sola función: `dominio.MonthHasData` (`backend/internal/dominio/tablero_datos.go`).

**Qué falta.** C1 v1.5 pide omitir cifras cuando «el sistema no tiene datos de ese mes» porque «sin datos no es cero». Lo pide en tres lugares: el mes anterior de `getDashboardMonthVsPlan` (M-01 ampliada: `previous_actual_cents`, `previous_cash_cents` y las dos diferencias), el devengado y el efectivo de cada mes de `getDashboardMonthlyOverview` (M-21: `accrued_cents` y `cash_cents`) y el ingreso de los dos (`income_*`, que ya se resuelve: sin quincena capturada ni esperada no hay ingreso). Pero ni C1, ni U3 2.0, ni las reglas de negocio dicen **cuándo** un mes de gasto «no tiene datos».

| Opción | Qué hace | Ventaja | Riesgo |
|---|---|---|---|
| **(a) PROVISIONAL** | Un mes sin datos es aquel en que el motor S2 da **devengado del hogar = 0 y efectivo del hogar = 0**. | No necesita consulta nueva: sale de las mismas dos cifras del motor. Es lo que ya hacía `getDashboardMonthlyOverview` (corrección r2). | Un mes real sin ningún gasto del hogar (muy raro en una casa) se mostraría como «—» en lugar de 0.00. |
| (b) | Sin datos = mes **anterior al primer mes con algún gasto o pago registrado** (igual que la tendencia mensual, N-044). Después de ese mes, un cero es cero. | Distingue «no se usaba el sistema» de «no se gastó». | Una consulta más por petición; un gasto viejo capturado tarde cambia qué meses tienen datos. |
| (c) | Nunca se omite: siempre se manda la cifra, aunque sea 0. | Lo más simple. | Contradice «sin datos no es cero» de C1 v1.5 y de U3. |

**Actualizado en WF-3c bis (`backend-cierre`).** La misma función decide un cuarto lugar: el mes anterior de `getIncomeMonth` (M-18, K-2). Si el mes anterior no tiene quincenas (capturadas ni esperadas) ni extras, `previous_fortnights_cents`, `previous_extras_cents` y `fortnights_difference_cents` van ausentes; `previous_month` va siempre.

**Impacto.** Solo cambia qué campos opcionales vienen ausentes y, por eso, si la PWA dibuja «—» o «$0.00» en K-1 (comparación con el mes anterior), G-1 y G-6. No cambia ninguna cifra del motor. Cambiar de opción toca solo `dominio.MonthHasData` (y, en la (b), una consulta en `dominio/tablero_*.go`); no hace falta migración ni cambio de C1.
