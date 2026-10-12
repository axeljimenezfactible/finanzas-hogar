# revisor-U3-r2 · reporte de agente (WF-1)

Revisor independiente, ronda 2 de 3, de `docs/design/U3-dashboard.md` v1.1, solo contra `docs/criterios/U3.md` v1.2 y decisiones §G y §K. No corregí nada ni escribí nada de lo revisado.

## Hecho
- Leí el protocolo completo, decisiones §G y §K, el criterio U3 completo, U3 v1.1 completo por partes, DESIGN §4.2, §5.4 (grep), §7.5.4, §7.6 a §7.8 y §8, los extractos (por `operationId`, con script), C1-datos M-01, M-08 y M-11 (grep), `model.json` y `data.json`, el reporte de U3-dashboard (§Solicitudes) y el reporte `correccion-U3-r1.md`.
- Reejecuté por mi cuenta todas las verificaciones de U3-01 a U3-12 y comprobé los 7 hallazgos medios de la ronda 1: los 7 quedan cerrados.
- Encontré 1 hallazgo alto (H-1, capas de T-2), 2 medios (H-2, P1 frente a «octubre»; H-3, resumen de gráfica de DESIGN §7.7) y 6 bajos (L-1 a L-3 siguen de la ronda 1; L-4 a L-6 son nuevos).
- Veredicto: **NO APROBADO**. Detalle en `docs/estado/wf1/REVISION-U3-r2.md`.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-U3-r2.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-U3-r2.md` (este reporte)
- Temporales en `datos/tmp/revU3r2/` (`ops.py`, `fields.py`, `verify.py`, `t2.py`, `tokens.py`): borrados al terminar.

## Verificación
Todo con `python3 -I`:
- `ops.py`: 27 operaciones citadas existen en los extractos, más la palabra `operationId`; las 13 rutas citadas coinciden en método y ruta (el prefijo `/api/v1` lo separa el patrón del script). Parámetros revisados en `listScheduledPayments`, `getMealPlanWeek`, `getDebtTrend`, `listCalendarEvents`, `getAlertSettings`, `getBudgetByCategory`, `previewAvailability`, `getBasketComparison`, `getDashboardBestBasket`, `getDashboardWeekMeals`, `getDashboardFortnightFlow`, `getAccrualVsCash`, `getBudgetMonthlyTrend`, `getMealPlanTrends`, `listSavingsGoals`, `getCashFlowFortnights` y `getPersonParticipation`.
- `fields.py`: de 146 nombres citados en minúsculas, 0 son campos inventados; los únicos que no aparecen son iconos, `tnum`, `weekday` y `day_of_month`. Esquemas revisados: `DashboardNextPayment` (`days_left`: «negativo si está vencido»), `DashboardMonthVsPlan` (`days_left`: «contando hoy»), `ScheduledPayment`, `AccrualVsCash` (`cumulative_difference_cents` = `periodos[i].diferencia_acumulada`; `is_canonical` oct-2026 a mar-2027), `CalendarEvent` y `AlertSettings`.
- `verify.py`:
  - **U3-01**: 6 tarjetas.
  - Tablas: §7 18×5, §6 18×7, §3 7×4, §8 8×4 y §1.4 5×4, todas con 0 celdas vacías.
  - **U3-09**: 0 fechas ISO; 30 fechas con día de la semana, 0 erróneas.
  - **U3-11**: 0 hex o `rgb(`.
  - Fechas: hoy + 60 = 2026-12-09, hoy + 30 = 2026-11-09; faltan 22 días del mes contando hoy y 20 hasta el 30-oct. Semanas ISO del 10, 11, 12 y 14 de octubre: 41, 41, 42 y 42. El 28-feb-2027 es domingo.
  - **U3-08**: 42 montos, 18 distintos; 17 están en `model.json` o `data.json` y $138.50 es K-03 (recalculado).
  - Pagos de D-3 en el rango: 30-oct, 15-nov y 30-nov. Σ P1–P6 = 14,362.48 y 13,103.78.
- `t2.py`: con la ventana desde P1, la cima es el devengado acumulado. Con la ventana desde P2, la cima es 2,921.25 contra 3,895.00 del oráculo (1,947.50 dentro de la ventana) y la capa 1 es 1,947.50 contra 2,921.25; en P6, 13,388.73 contra 14,362.48 → H-1. Textos «…»: 452, 0 con palabras en inglés de mi lista.
- `tokens.py`: 64 tokens citados, 0 sin definir en C2-tokens, C2-colores o la tabla `ui.*` de DESIGN; las 32 combinaciones de `color.role.*` existen. DESIGN dice «202 pares, 202 cumplen» (línea 186).
- Karen: 40 menciones, todas dentro de §9 y A-13.
- Criterios: U3-01, 03, 04, 06, 07, 09 y 10 cumplen. U3-02, 08, 11 y 12 cumplen con observaciones (H-2, H-3). U3-05 no cumple (H-1).

## Supuestos
1. Tomé `cumulative_difference_cents` como acumulado desde el origen del oráculo (P1), porque C1 lo define como `periodos[i].diferencia_acumulada` de `model.json`. Si C1 lo reiniciara al empezar la ventana, H-1 cambiaría de forma, pero U3 seguiría sin decir el origen de las capas.

## Bloqueos
Ninguno nuevo de mi parte. Los revisores no escriben bloqueos de diseño; la corrección de H-1 y H-2 propone abrir `docs/estado/bloqueos/U3-t2-acumulados.md` o `docs/estado/bloqueos/U3-mes-calendario-vs-periodo.md` si no se alinea con C1.

## Pendientes
Corregir en `docs/design/U3-dashboard.md` (ronda 2 de corrección):
- H-1: T-2, líneas 233 y 236.
- H-2: línea 238, y coherencia con las líneas 90, 204, 392 y 397.
- H-3: §4, línea 211, y T-1 a T-4.
- Bajos L-1 a L-6: líneas 146, 243, 293, 185, 180, 63 y 261.

## Solicitudes fuera de alcance
1. `docs/contratos/openapi.yaml` (vía `backend/internal/api/c1gen/build.py`), esquema `AccrualVsCash`: decir el origen de `cumulative_difference_cents` y considerar `cumulative_accrued_cents` y `cumulative_cash_cents` con ese mismo origen. Motivo: H-1, T-2 no puede dibujar el área apilada de DESIGN §8.3 con un solo origen.
2. `docs/design/DESIGN.md` §7.8: aclarar si «dentro de los días de aviso» es rango o lista exacta (L-6). Motivo: con valores no consecutivos U3 y DESIGN darían estados distintos.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna sobre valores K31: U3 no define parámetros operativos.

## Consumo aproximado
~39 acciones de herramienta · ~170 KB leídos (U3 71 KB; DESIGN por secciones ~25 KB; criterio, protocolo y decisiones ~33 KB; reporte de corrección 11 KB; salidas de script y greps ~30 KB) · ~32 KB escritos (revisión 13.5 KB, reporte 5 KB, scripts ~13 KB).
