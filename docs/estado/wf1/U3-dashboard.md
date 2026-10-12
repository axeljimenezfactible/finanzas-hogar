# U3-dashboard · reporte final (WF-1)

Agente U3-dashboard (Opus) · 11-oct-2026 · Evaluado contra `docs/criterios/U3.md` v1.2.

## Hecho
- `docs/design/U3-dashboard.md` v1.0 (58.6 KB, ≤ 60 KB): estructura de Inicio (pestañas, acción primaria de página «Registrar gasto», anatomía de la «Tarjeta de métrica» con tokens), Resumen con **6 tarjetas** en orden fijo (métrica C1, operationId y campos, fórmula, periodo, textos exactos es-MX, acción primaria única, vacío/error, ejemplo trazable), Detalle (D-1 a D-7, incluida «Participación de Karen»), Tendencias (T-1 barras agrupadas, T-2 área apilada devengado contra efectivo, T-3 y T-4 líneas; orden fijo de series con forma, patrón y etiqueta; «Ver como tabla»; 0 donas), Calendario (agenda/mes/mes con día; franjas de quincena K8, día de pago aparte, febrero), tabla de filtros (18 bloques × 6 filtros), tabla de estados (18 bloques × vacío/carga/error/sin conexión), diseño por tamaño (360/412, 768/1024, 1440), reglas de Karen, color y accesibilidad, componentes para F2, 15 criterios de aceptación para F2 (WF-3), lista de capturas U5, trazabilidad U3-01 a U3-12, supuestos y dependencias.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/U3-dashboard.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/U3-dashboard.md` (este reporte)

## Verificación
Script `python3 -I verify.py` (en el scratchpad de la sesión; extrae operationId, M-xx, montos, fechas ISO, hex y fechas con día de la semana):
- **U3-01**: `grep -c "^### Tarjeta"` = **6** (Mes contra plan, Flujo de la quincena, Próximo pago y disponible, Comidas de la semana, Mejor canasta, Ahorro y metas). Tarjeta 4 cita K38 y U1 §0.2: Dom 11 oct 2026 → tarjeta en `2026-W41` (N-023) y P-10 en `2026-W42`; Mié 14 oct 2026 → ambas `2026-W42` (semanas ISO calculadas con `datetime.isocalendar`).
- **U3-02**: 27 operationId citados, **0 inexistentes** en `docs/contratos/extractos/*.yaml` (210 operaciones). Métricas citadas M-01 a M-11 y M-13 a M-17: **0 faltantes** en `C1-datos.md` §9. Tabla tarjeta/panel → métrica → operación: T1 M-01 `getDashboardMonthVsPlan`; T2 M-02 `getDashboardFortnightFlow`; T3 M-03/M-04 `getDashboardNextPayment`; T4 M-05 `getDashboardWeekMeals`; T5 M-06 `getDashboardBestBasket`; T6 M-07 `getDashboardSavings`; D-1 M-09 `getBudgetByCategory`; D-2 M-10 `getCashFlowFortnights`; D-3 M-13 `listScheduledPayments` + M-04 `previewAvailability`; D-4 N-121 `getMealPlanWeek`; D-5 N-162 `getBasketComparison`; D-6 M-16 `listSavingsGoals`; D-7 M-08 `getPersonParticipation`; T-1 M-14 `getBudgetMonthlyTrend`; T-2 M-11 `getAccrualVsCash`; T-3 M-17 `getDebtTrend`; T-4 M-15 `getMealPlanTrends`; C-1 M-13 `listCalendarEvents`.
- **U3-03**: §7, 18 bloques × 4 estados (vacío, carga, error, sin conexión), 0 celdas vacías; reglas comunes en §1.3.
- **U3-04**: §8 con columnas 360/412, 768/1024, 1440; táctil `size.touch.min` 44 px en §1.2 y §8.
- **U3-05**: §3, §4, §5 con propósito, contenido, métricas, filtros y estados (15 de 15); patrones barras/líneas/área apilada; **0 donas**; quincenas K8, día de pago aparte y «Dom 28 feb 2027» en §5.
- **U3-06**: `grep -c Karen` = 26 líneas; todas en tarjeta 3, D-3, D-7 (§3.1), Calendario, filtro de familia (K5) y §9 (sin cuenta, no ve el dashboard, PDF que Axel reenvía, sin indicadores de su familia).
- **U3-07**: Tarjeta 3 muestra solo la lectura activa (B por defecto); la otra en D-3 con `previewAvailability`; insignia «Pendiente confirmar»; ejemplo $973.75 el Vie 30 oct 2026 (R-13), A $10,000.00 y B $6,905.50 (R-64).
- **U3-08**: 18 montos en el documento, **0 sin origen**: 17 en `model.json`/`data.json` y $138.50 = K-03 de `S5.md` (`totales.walmart − totales.bodega`, recalculado = 138.50). Lo demás son marcadores sin número.
- **U3-09**: `grep "[0-9]{4}-[0-9]{2}-[0-9]{2}"` = **0**; 25 fechas con día de la semana, **0** con el día equivocado (verificado con `datetime`); textos de interfaz en es-MX.
- **U3-10**: §6, tabla de 18 bloques × 6 filtros con S/V/—, sin celdas vacías; periodo por defecto mes calendario con quincena como subperiodo. `listExpenses` y `getExpenseSummary` admiten los 6 filtros (enlace a P-02).
- **U3-11**: `grep "#[0-9A-Fa-f]{3,8}\|rgb("` = **0**; tokens solo por nombre (§1.2, §10); gráficas con orden fijo, forma/patrón, etiquetas en texto y «Ver como tabla» (§4).
- **U3-12**: §11 (F2 en WF-3, componentes, A-01 a A-15), §14 Supuestos, §15 Dependencias; `ls docs/estado/bloqueos/U3-*` = ninguno (sin bloqueos propios).

## Supuestos
1. Rangos cortos en tarjetas («Del 16 al 31 oct», «Lun 12 al Dom 18 oct») y meses en gráficas («oct-2026») derivados de U1 §0.4 (solo texto).
2. Hora de carga sin conexión «[HH:MM]» en 24 h (solo texto).
3. Las tarjetas del Resumen no llevan filtros; los filtros viven en Detalle, Tendencias y Calendario.
4. Tarjeta 3 con varias líneas revolventes: la del próximo pago y «+[n] líneas».
5. T-2 como área apilada acumulada (efectivo acumulado + `cumulative_difference_cents`), con cambio a dos líneas si algún mes es negativo (regla de dibujo).
Ninguno cambia cifras, esquema, operaciones ni tokens.

## Bloqueos
Ninguno propio. `docs/estado/bloqueos/S1-parametros-avisos.md` no afecta a U3 (el dashboard no lee ni escribe parámetros de avisos; el estado de M-01 llega calculado). Bloqueos de S2 abiertos durante WF-1 que tocan cifras mostradas, marcados **PROVISIONAL** en U3 (U3 muestra lo que llegue; no cambia estructura, textos fijos ni operaciones): `docs/estado/bloqueos/S2-presupuesto-total.md` (tarjeta 1), `docs/estado/bloqueos/S2-flujo-quincenal.md` (tarjeta 2, D-2), `docs/estado/bloqueos/S2-metas.md` (tarjeta 6, D-6), `docs/estado/bloqueos/S2-proyeccion-interruptor.md` (T-2, meses 7 a 12 con el interruptor encendido).

## Pendientes
- Integrar con DESIGN §7 y §8 cuando U2b-2 los escriba (nombres de componentes y patrones de U3 §11.1 y §4).

## Solicitudes fuera de alcance
1. `docs/design/DESIGN.md` §7, §8 y §5.4 (U2b-2): incluir la «Tarjeta de métrica» con la anatomía de U3 §1.2; los patrones barras agrupadas, barras horizontales, área apilada acumulada (con la regla de dos líneas) y líneas con meta discontinua, con «Ver como tabla»; orden de marcadores (círculo, cuadrado, triángulo, rombo) y trazos (continuo, discontinuo, punteado); en el mapa de iconos, los de las tarjetas: `wallet`, `banknote`, `calendar-clock`, `utensils`, `shopping-cart`, `piggy-bank`. Motivo: U3-11 y U3-05 dependen de que DESIGN los fije.
2. `docs/contratos/openapi.yaml` (vía `backend/internal/api/c1gen/build.py`): agregar `second_store_id` a `DashboardBestBasket` (getDashboardBestBasket). Motivo: la tarjeta 5 dice «menos que la siguiente tienda» sin poder nombrarla. No bloquea.
3. `docs/design/U1-arquitectura.md` §3.3 N-027 y `x-necesidades` de `previewAvailability`: sumar N-285 a N-027 (y N-027 a la operación). Motivo: D-3 muestra las dos lecturas A/B, que U3-07 exige accesibles fuera de la tarjeta.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Sin valores K31: U3 no define parámetros operativos.
- Para la puerta 2: confirmar el supuesto 5 (dibujo de T-2) y la solicitud 2 (nombre de la segunda tienda en la tarjeta 5).

## Consumo aproximado
~39 acciones de herramienta; ~197 KB leídos; ~80 KB escritos (documento 58.6 KB, reporte 8 KB, scripts de verificación y edición 12 KB en el scratchpad).
