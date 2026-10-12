# correccion-U4-r1 · Reporte final

Flujo WF-1 · 11-oct-2026 · Corrige los hallazgos H-1 (alta), H-2 y H-3 (media) de `docs/estado/wf1/REVISION-U4-r1.md`. No toca los bajos (H-4 a H-8) ni reabre decisiones de §A–§K.

## Hecho
- **H-1 · U4-05 (vista previa antes de confirmar).** Nueva convención «vista previa en línea» en la parte a §0.2: en un campo con efecto numérico, «Antes → Después» bajo el campo mientras el valor tecleado difiere del guardado; `Enter` o salir guarda, `Esc` descarta. Respeta el guardado automático de U1 §2.6, así que U1 no cambia. Casos definidos:
  - P-26.05 vista previa 5, **límite**: Δ = límite nuevo − actual. En revolvente, las dos lecturas y sus 14 valores se mueven exactamente Δ (RN-09; comprobado: A − B es constante en las 14 cifras del oráculo). «Antes» sale de `previewAvailability` y `getAvailabilitySeries`, y «Después» se calcula en el navegador, en centavos. En tarjeta: `available_cents` + Δ (C1: «tarjeta = límite − saldo»). En préstamo y tienda no hay disponible y el bloque lo dice.
  - P-26.05 vista previa 6, **frecuencia de pagos**: `previewCreditPurchase` con `payment_frequency` (el campo ya existe en `CreditPurchaseCreate`), antes y después, con una compra de ejemplo que toma capital y plazo de la última compra; si la línea no tiene compras, los campos nacen vacíos (no se inventa un monto). Semilla: 30-oct, 15-nov, 30-nov y 15-dic-2026, pagos de $973.75, total $3,895.00.
  - P-26.05 vista previa 7, **días mínimos al primer pago**: PROVISIONAL por la solicitud SC-U4-1 a C1, porque `previewCreditPurchase` no recibe `first_payment_min_days`. Ejemplo de RN-06: 30-oct-2026.
  - P-26.05 vista previa 8, **comisiones**: hoy ninguna regla las usa. La vista previa lo dice, y la decisión de si deben entrar a algún cálculo queda PROVISIONAL (bloqueo `U4-vista-previa-formulas`, pregunta 1).
  - P-26.02 vista previa 4, **despensa programada**: próximas 3 fechas, más `previewCreditPurchase` con la línea, el plazo, el monto y la próxima fecha (pago por quincena, fechas y total), más la baja del disponible (RN-09).
  - P-26.13, **metas**: avance antes → después (RN-18, un decimal; ejemplo de RN-18: 50.0 %) y reserva del plan del mes antes → después (RN-18). La aportación necesaria y el estado quedan PROVISIONAL: C1 no da la fórmula (bloqueo `U4-vista-previa-formulas`, pregunta 2, y solicitud SC-U4-2 condicional).
  - P-26.14, **ritmo**: existencia estimada (RN-15) y estado (K7), antes → después. Ejemplo de RN-15: 1,486.571 ml. Días restantes y consumo por semana quedan PROVISIONAL (bloqueo, pregunta 3).
- **H-2 · U4-10 y U4-09 (convenciones).** La parte a §0 pasa a ser el **único** bloque para las 14 secciones:
  - Guardado: `circle-check` y «Guardado» en `color.role.success.light`/`.dark` durante 2 s. Los 2 s ya tienen origen: DESIGN §7.13.
  - Un solo texto de error al guardar: «No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar». Error al cargar: franja «No se pudo cargar {qué} · Reintentar» (DESIGN §7.14).
  - Sin conexión (U1 §2.4) y carga con aviso a los 10 s.
  - Diálogo crítico con `shield-alert` en `color.critical.light`/`.dark` (DESIGN §7.12).
  - §0.7: táctil, teclado y foco visible (`size.touch.min`, `size.control.coarse`/`fine`, `ui.border.focus`).
  - La parte b ya no tiene sus bullets de Guardado, Confirmación, Estados y Táctil: remite a la parte a §0. Los estados de cada sección de b usan los textos únicos.
- **H-3 · U4-07 (reparto).** Una sola regla, la de RN-03, S1-08 y `CostShareRule`:
  - Se editan solo las partes explícitas (`parts`).
  - La persona del resto (`remainder_person_id`) no tiene campo: su parte se muestra calculada como «el resto». La regla suma 100 % por construcción.
  - Único error: que las partes con porcentaje pasen de 100 %. V-16 se reescribe como «Las partes con porcentaje suman {suma} %. No pueden pasar de 100 %: lo que falta lo recibe {persona del resto}.» y P-26.12 lo cita.
  - Se quitan «Las partes suman más de 100 %» y el «recálculo» de la parte del hogar como campo editable.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/U4-ajustes-a.md` (1.1): §0, P-26.02 y P-26.05.
- `/home/deepsec/finanzas-hogar/docs/design/U4-ajustes-b.md` (1.1): convenciones, estados de P-26.08, P-26.11, P-26.12, P-26.13 y P-26.14; P-26.12, P-26.13 y P-26.14.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U4-vista-previa-formulas.md` (nuevo).
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-U4-r1.md` (este reporte).

## Verificación
El script `verifica.py` (en el scratchpad de la sesión, fuera del repositorio) corre con `python3 -I` contra las dos partes, `openapi.yaml` (solo la expresión regular de `operationId`), `C2-colores.json`, `DESIGN.md` y `oraculo/model.json`. Resultados:
- U4-01: 14 secciones P-26.01…P-26.14, en orden.
- U4-02: 95 + 45 `operationId` citados, 0 inexistentes. No se quitó ninguna operación, así que la cobertura de las 73 operaciones de Ajustes no cambia.
- U4-05: «se ve(n) al guardar» = 0. Cada caso de H-1 tiene vista previa definida (arriba); los que dependen de una fórmula o de una operación sin decidir llevan PROVISIONAL con cita.
- U4-07: V-16 aparece 1 vez en el catálogo y 2 veces en P-26.12; «deben sumar 100» = 0, «suman más de 100» = 0, «se recalcula como el resto» = 0.
- U4-09: hex = 0, `rgb(` = 0. 21 tokens citados, 0 inexistentes (los `color.*` en C2; `size.*`, `space.*`, `ui.*`, `border.*` y `focus.*` en DESIGN). La parte a §0 contiene «Sin conexión con el servidor de casa», los 10 s, `shield-alert`, «foco visible» y `size.control.coarse`.
- U4-10: `neutral.700` = 0, «No se guardó» = 0. Queda un solo texto de error al guardar y una sola fila «Guardado».
- U4-11: fechas ISO = 0.
- U4-12: «mínimo de seguridad» en filas de tabla = 0.
- Cifras, reproducidas con Decimal sin diferencias:
  - Pagos de la despensa 1: 4 × 973.75 = 3,895.00 = 3,094.50 + 800.50.
  - 3,094.50 × 80050/309450 = 800.50.
  - `minimo_A` 8094.5 y `minimo_B` 5000.0; `eventos[0]` A 10000.0 y B 6905.5.
  - A − B constante en las 14 cifras de cada serie (sostiene el «+ Δ» del límite).
  - RN-15: 946000 × 30 / 21 → 1351429 (HALF_UP); 2838000 − 1351429 = 1486571.
  - RN-18: 50.0 %.
- PROVISIONAL: 29 líneas en total. 8 no citan bloqueo en la misma línea: 1 es la regla general de §0.2 y 7 son de P-26.08, P-26.09 y P-26.11 de la parte b, que ya estaban antes de esta corrección (hallazgo H-7, baja, fuera de esta ronda). Todas las marcas nuevas citan su bloqueo o su solicitud.
- No aplican pruebas Go, `compara.py` ni la compilación de tokens: U4 es un documento de diseño sin código. Los temporales (`datos/tmp/correccion-U4-r1`) quedaron borrados.

## Supuestos
- La vista previa en línea con guardado al salir del campo se toma como «confirmación» (U4-05) para no cambiar U1 §2.6, que pide guardado automático en P-26.xx. No cambia cifras, esquema, endpoints ni tokens.
- Para P-26.12 se toma como «persona del resto» la que C1 devuelve en `remainder_person_id`; U4 no la fija.

## Bloqueos
- `docs/estado/bloqueos/U4-vista-previa-formulas.md` (nuevo): (1) comisiones del crédito, (2) fórmula de la aportación necesaria, (3) días restantes y consumo por semana. Siguen abiertos los anteriores: `S1-parametros-avisos`, `U4-iva-credito`, `U4-dia-mes-reporte` y `U4-reparto-retroactivo`.

## Pendientes
- Hallazgos bajos H-4 a H-8 de `REVISION-U4-r1.md`: no estaban asignados a esta ronda.
- Al integrar `docs/design/U4-ajustes.md`, el §0 de la parte a queda como bloque común de las 14 secciones.

## Solicitudes fuera de alcance
- **SC-U4-1 · C1** (`docs/contratos/openapi.yaml`, cuerpo de `previewCreditPurchase`, hoy `CreditPurchaseCreate`): agregar el campo opcional `first_payment_min_days` (entero ≥ 0), que solo se usa en la vista previa y no se guarda. Motivo: vista previa 7 de P-26.05 (U4-05).
- **SC-U4-2 · C1, condicional** a la respuesta del bloqueo `U4-vista-previa-formulas` (preguntas 2 y 3): si la fórmula no se calcula en el navegador, una vista previa sin guardar de `required_monthly_cents` y `status` (meta) y de `days_left` y `weekly_consumption_milli` (inventario) con los valores propuestos. Motivo: U4-05 en P-26.13 y P-26.14.
- **SC-U4-3 · C1, opcional**: que `previewAvailability` y `getAvailabilitySeries` acepten `credit_limit_cents` propuesto, para que el motor S3 sea la única fuente. Hoy el cálculo en el navegador es exacto por RN-09, así que no es necesario.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Bloqueo `U4-vista-previa-formulas`, recomendación del agente:
  - Pregunta 1, opción (a): las comisiones no entran a ningún cálculo. Es lo que hace hoy el sistema y no cambia ninguna cifra canon.
  - Pregunta 2, opción (a): meses calendario completos incluidos el actual y el de la fecha meta, con ROUND_HALF_UP al centavo. Es simple, se reproduce en el navegador y no necesita una operación nueva.
  - Pregunta 3, opción (a): días restantes truncados (nunca promete un día de más) y consumo por semana con HALF_UP en milésimas, igual que RN-15.
- No hay valores K31 en esta corrección.

## Consumo aproximado
55 acciones de herramienta · ~190 KB leídos · ~35 KB escritos de contenido propio (scripts de edición, bloqueo y reporte). Las dos partes se reescribieron por script con reemplazos puntuales.
