# correccion-U3-r1 · reporte final (WF-1)

Agente `correccion-U3-r1` · 11-oct-2026 · Artefacto: `docs/design/U3-dashboard.md` (v1.0 → **v1.1**, 462 líneas, 70.9 KB) · Corrige los 7 hallazgos medios de `docs/estado/wf1/REVISION-U3-r1.md` contra `docs/criterios/U3.md` v1.2, sin reabrir §A–§K.

## Hecho

| Hallazgo | Corrección en U3 v1.1 |
|---|---|
| H-1.1 título de tarjeta | §1.2: «Etiqueta (título)» en `font.size.small` y `ui.text.secondary`, sin dos puntos (DESIGN §7.6). El icono de 20 px se queda (`size.icon.sm`, DESIGN §5.2; §5.4 prevé iconos en tarjetas) y se pide a DESIGN que lo sume a la anatomía (solicitud 1). §8, fila «Tipografía», alineada. |
| H-1.2 relleno | §1.2 y §8: relleno `space.4` por debajo de 768 px y `space.6` desde 768 (DESIGN §7.6). El medianil sigue DESIGN §4.2 (16 px hasta 1023, 24 px desde 1024). |
| H-1.3 carga de gráficas | §7, T-1 a T-4: «esqueleto del alto final en `ui.skeleton.on-card`, sin ejes falsos (DESIGN §7.7)». §1.3: el esqueleto de la tarjeta va en `ui.skeleton.on-card`. |
| H-1.4 estados de pago | Nuevo **§1.4** con los 5 estados de DESIGN §7.8 (Pagado, Vencido `calendar-x`, Por vencer `calendar-clock`, Parcial `circle-dot` con «Falta la parte de Karen» según K15, Pendiente `circle-dashed`), con el dato de C1 de cada uno y qué bloque puede mostrarlo. «Por vencer» lee `alerts.payment_due_days_before` con `getAlertSettings` una vez por página (§1.2, K44), nunca como constante. Tarjeta 3, D-3 y §5 citan §1.4. «Sin marcar como pagado» (U1 N-022, C1 M-03) se queda como línea junto a «Vencido»: así U3 no contradice a U1 ni a DESIGN (solicitud 2). |
| H-1.5 destino de la tarjeta | §1.2 «Destino»: U3 se queda con la pantalla destino de cada acción primaria, porque U3-01 exige que la tarjeta 4 lleve a P-10. Se pide a DESIGN §7.6 (solicitud 1). |
| H-2 patrones de §8 | §4: etiquetado selectivo de DESIGN §8.5 y nueva regla de marcas de §8.1 (nunca un contorno). T-1: «Plan» con rayas a 45° y sin contorno. T-2: capa 2 en cero con «Pagado por adelantado», sin capa negativa ni cambio a líneas; línea «Hoy»; etiquetas finales; el subtítulo no rotula la deuda viva (K19, N-091, M-11; solicitud 3). T-4b: Ximena con cuadrado y trazo 8 · 4 (§8.7); la línea de meta pasa a puntos 2 · 4 para no chocar con el trazo de Ximena (supuesto 8, solicitud 4). Supuesto 5 retirado. |
| H-3 ejemplo de D-3 | El 10-oct-2026, el rango de hoy a hoy + 60 termina el 9-dic-2026 y caben 3 pagos: «Vie 30 oct 2026», «Dom 15 nov 2026» y «Lun 30 nov 2026» (R-13 a R-15). El de R-16 (15-dic) queda fuera. El ejemplo del Calendario (oct–dic) no cambia. |
| H-4 petición de D-3 | Dos peticiones explícitas a `listScheduledPayments`: (1) vencidos con `status` = `overdue` y `period` = `[vencido más antiguo]..[ayer]`, donde la fecha es `next_payment.due_date` cuando `is_overdue` = true (M-03); si no hay vencido, esta petición no se hace. (2) Próximos con `status` = `pending` y `period` = `[hoy]..[hoy + 60]`. «Pagado» se quita de D-3 porque no puede llegar; «Parcial» queda definido con `shares[]`. §6 cambia en la fila D-3. |
| H-5 Karen en T-2 | §9 y A-13 suman las columnas informativas de Karen en la tabla de T-2 (`accrued_external_cents`, `cash_external_cents`). |
| H-6 fila «Total» de T-2 | La fila «Total» suma las 12 filas visibles y no lleva cifra de canon. R-51 y R-52 ya no van en la tabla (ver «Hallazgo de causa» abajo). |
| H-7 D-4 y D-5 sin métrica | D-4 y D-5 escriben fórmula, periodo y filtros con la marca «Métrica pendiente en C1-datos §9»; la solicitud está en §15 y abajo (5). |

Además: encabezado v1.1 con la lista de cambios; §11.1 y §15 ya no tratan DESIGN §7 y §8 como pendientes; A-02 y A-11 se actualizaron y entraron A-16 (D-3 el 10-oct-2026) y A-17 (estados de §1.4 con reloj simulado); supuestos 6 a 8.

**Hallazgo de causa (H-6), sin ajustar cifras:** en el oráculo, P1 a P6 **no son meses calendario**: van del 10 de un mes al 9 del siguiente (`periodos[0]` = «10-oct-2026 a 9-nov-2026», …, P6 = «10-mar-2027 a 9-abr-2027»). R-51 = Σ devengado P1–P6 = 14,362.48 y R-52 = Σ efectivo P1–P6 = 13,103.78 (recalculados con `model.json`: coinciden). Una fila «oct-2026 a mar-2027» en meses calendario no sería la misma medida, así que no se propone. C1-datos M-01 y M-11 llaman «oct-2026» a P1 (por el hallazgo 4 de S1). En efectivo da lo mismo (no hay pagos del día 1 al 9), pero en devengado depende del método de S2. Va como solicitud 6.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/U3-dashboard.md` (modificado, v1.1)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-U3-r1.md` (este reporte)

## Verificación
Scripts en el scratchpad de la sesión (`edit.py` aplica 37 reemplazos exactos y comprueba antes que cada texto viejo aparezca las veces esperadas; `verify.py` revisa el documento; `peek.py` lee el oráculo), con `python3 -I`:
- `edit.py` → «reemplazos: 37 bytes: 70862»; 0 textos sin coincidencia.
- **U3-01**: `^### Tarjeta` = **6**; tarjeta 4 sin cambios (K38, U1 §0.2).
- **U3-02**: 28 identificadores entre comillas invertidas; 27 operaciones existen en los extractos, incluida la nueva `getAlertSettings` (`ajustes-e.yaml`). El 28.º es la palabra `operationId`. Métricas citadas: M-01 a M-11 y M-13 a M-17, **0 faltantes** en C1-datos §9. D-4 y D-5 llevan fórmula, periodo y filtros propios mientras C1 no agregue su métrica (H-7).
- **U3-03**: §7 con 18 filas, **0 celdas vacías**; las 4 gráficas cargan sin ejes falsos.
- **U3-04**: §1.2 y §8 con el relleno de DESIGN §7.6; columnas 360/412, 768/1024 y 1440 sin cambios; táctil `size.touch.min`.
- **U3-05**: T-2 sigue DESIGN §8.3; 0 donas; quincenas K8 y febrero sin cambios.
- **U3-06**: Karen aparece en el encabezado de reglas, §1.4, la tarjeta 3, §3 (D-3), §3.1 (D-7), T-2 (tabla), §5, §6 (K5), §7 (D-7), §9, §11.2, §12 y §13. Todas las apariciones están dentro de lo que enumeran §9 y A-13.
- **U3-07**: tarjeta 3 con $973.75, «Vie 30 oct 2026», $6,905.50 B y $10,000.00 A, sin cambios.
- **U3-08**: 42 montos, 18 distintos; 17 están en `model.json` o `data.json` y $138.50 es K-03 de `S5.md` (3,233.00 − 3,094.50, recalculado). **0 sin origen.** Hoy + 60 el 10-oct-2026 = 9-dic-2026 (calculado con `datetime`).
- **U3-09**: 0 fechas ISO; 30 fechas con día de la semana, **0 con el día equivocado**.
- **U3-10**: §6 con 18 filas, **0 celdas vacías**.
- **U3-11**: 0 hex y 0 `rgb(`; 64 tokens citados, **0 sin definir** en `C2-tokens.json` o `C2-colores.json` (los nuevos son `ui.text.secondary`, `ui.skeleton.on-card`, `ui.viz.baseline`, `ui.icon.default`, `ui.text.default` y `size.icon.sm`). Iconos nuevos: `calendar-x`, `circle-dashed` y `circle-dot` vienen de DESIGN §7.8; no están en el mapa de §5.4 (solicitud 2).
- **U3-12**: §11, §14 (supuestos 6 a 8; el 5 queda retirado) y §15. No hay bloqueos `U3-*`.
- Tablas de §3 y §1.4: 4 columnas en todas sus filas. Ya no aparecen «U2b-2», «sobre cada punto» ni «marco de gráfica». «Mar 15 dic 2026» solo queda en el ejemplo del Calendario (oct–dic), que es correcto.
- No hubo pruebas Go, oráculo ni tokens que reejecutar: U3 es un documento y su verificación es la de arriba. No se tocó código.

## Supuestos
1. Supuestos 6, 7 y 8 de U3 §14: prioridad entre los estados de pago; los textos «Falta la parte del hogar» y «Vence en [d] días»; la línea de meta de T-4b con puntos 2 · 4. Son de presentación y no cambian cifras, esquema, operaciones ni tokens.
2. Calendario, estado «Parcial»: un evento `credit_payment` con `ref_table` = `scheduled_payment` se une con el pago de `listScheduledPayments` cuyo `id` es `ref_id` (C1: «Tabla de origen (nombre_codigo)»). Si no se puede unir, el pago se queda con su estado de `CalendarEvent` (U3 §5). Se pide confirmación a C1 (solicitud 7).

## Bloqueos
Ninguno nuevo. `docs/estado/bloqueos/S1-parametros-avisos.md` no afecta: los días de «Por vencer» son `alerts.payment_due_days_before`, ya decididos (decisión F, K44). Los 4 bloqueos S2 que cita U3 siguen igual.

## Pendientes
Los hallazgos bajos de la revisión quedaron fuera de esta ronda, que solo pidió corregir los altos y medios:
- L-1: `docs/design/U3-dashboard.md`, tarjeta 5: «[covered] de [items]» → «[con precio] de [artículos]».
- L-2: `docs/design/U3-dashboard.md`, T-3: la lista de marcadores y trazos con filtro de crédito debe citar DESIGN §8.6.
- L-3: `docs/design/U3-dashboard.md`, D-3 (b) y fila D-3 de §6: `previewAvailability` solo admite `credit_id`; falta «Este bloque no usa el filtro [nombre]».

## Solicitudes fuera de alcance
1. `docs/design/DESIGN.md` §7.6: la tarjeta entera lleva a la pantalla destino de su acción primaria y no a la pestaña Detalle, porque U3-01 exige que la tarjeta 4 lleve a P-10; sumar a la anatomía el icono de 20 px junto a la etiqueta (§5.2 y §5.4).
2. `docs/design/DESIGN.md` §7.8 y §5.4: «Sin marcar como pagado» como línea junto a «Vencido», como piden U1 N-022 y C1 M-03; la prioridad entre estados; los textos «Falta la parte del hogar» y «Vence en [d] días». Agregar `calendar-x`, `circle-dashed` y `circle-dot` al mapa de iconos de §5.4, porque hoy solo están en §7.8.
3. `docs/design/DESIGN.md` §8.3: quitar del subtítulo «Deuda viva $6,293.44 = …» o justificarlo frente a K19, N-091 y M-11 («no se rotula deuda viva»).
4. `docs/design/DESIGN.md` §8: fijar cómo se dibuja la línea de meta en una gráfica de líneas. U3 T-4b usa el color de la persona, puntos 2 · 4 y sin marcadores.
5. `docs/contratos/C1-datos.md` §9: dos métricas nuevas (ids que asigne C1). «Comidas de hoy» (N-121 vía N-027): casillas de `getMealPlanWeek` con `plan_date` = hoy, por tiempo y persona, con `kcal` y `cost_cents` tal como llegan; periodo: semana ISO que contiene hoy; filtro `person_id`. «Diferencias de precio» (N-162 vía N-027): los 5 renglones de `getBasketComparison` con mayor `difference_cents` y `stores[].total_cents`; periodo: lista vigente (`week` según K38); filtro de tienda en la vista. Motivo: U3-02 (H-7).
6. `docs/contratos/C1-datos.md` M-01 y M-11 (y la descripción de `getAccrualVsCash`): aclarar si las filas mensuales son meses calendario o los periodos del oráculo, que van del 10 al 9. El ejemplo «oct-2026» usa las cifras de P1 (10-oct-2026 a 9-nov-2026). Ver «Hallazgo de causa (H-6)».
7. `docs/contratos/openapi.yaml`, operación `listScheduledPayments` (vía `c1gen/build.py`): aclarar qué pasa con `status` = `overdue` sin `period`, porque la descripción «de hoy en adelante» lo deja vacío; si `pending` incluye los vencidos; y que `ref_table` de un `CalendarEvent` de pago es `scheduled_payment`. U3 evita la ambigüedad con rangos explícitos.
8. Siguen abiertas las solicitudes 2 y 3 del reporte `U3-dashboard.md`: `second_store_id` en `DashboardBestBasket`, y N-285 en N-027.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Sin valores K31: U3 no define parámetros operativos. Los días de «Por vencer» son `alerts.*` (K44).
- Para la puerta 2: ver las solicitudes 1 a 4 a DESIGN, que son los cuatro puntos en que U3 se aparta de DESIGN o lo completa.

## Consumo aproximado
~34 acciones de herramienta · ~180 KB leídos · ~105 KB escritos (documento reescrito 70.9 KB, scripts 23 KB, reporte 10 KB).
