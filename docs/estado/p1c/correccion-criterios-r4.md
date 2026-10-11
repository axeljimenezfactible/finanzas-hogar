# correccion-criterios-r4 · P1d, ronda 4 de criterios y bajos de C1 (11-oct-2026)

## Hecho
- A (alta): C1-22, U4-12, S7 (lista de cambios, S7-04, S7-08, S7-10, §5) e INDICE (§6 líneas 68 y 96, §11 línea 174) ya tratan los umbrales 80/100 % y los días antes [1, 0] como claves `alerts.budget_warning_pct`, `alerts.budget_limit_pct` y `alerts.payment_due_days_before` de `app_setting` (00003, «Decisión F, K44», `getAlertSettings`/`updateAlertSettings` con `value`, `initial_value` y `reset_parameters`). Se quitó «no se siembran claves `alerts.*`» y la comprobación contra `params_json` para esos tres. La excepción PROVISIONAL queda solo para `send_time`, `advance_days`, `weekday` y `day_of_month` (lugar según `S1-parametros-avisos`; valor de `send_time` y `advance_days` según K31/S7). U4-12: esos dos no tienen valor inicial y su «restablecer» se activa cuando S7 proponga el valor.
- B: P3-06 cuenta solo lo descargado (sin `siguiente-prompt.md`), el agente Sonnet asigna archivo o página a artboard en `P3-iter-N.md` (PDF de N páginas cuenta N); `P-xx` es sugerencia del LEEME.
- C: P3-04 comprueba 0 llamadas a `mcp__claude-in-chrome__`/`navigate` en el registro del flujo (ruta absoluta) o la declaración «0 acciones de navegador» contrastada con la lista de herramientas.
- D: regex de P3-01, P3-03 y P3-12 exigen valor de 6 o más caracteres tras `:`/`=`; probado: 5 advertencias dan 0 y 6 credenciales de muestra dan 6.
- E: `send_time` y `advance_days` pasan a «los demás… sección que U4 proponga, con justificación en SUPUESTOS».
- F: P3-11 hace `ls -d` de `salida/final/` y entra a la fila de la puerta 3 de INDICE §3.
- G: S6-10 no inventa el factor: C1 v1.2 declara «10 MB» (`attachment.max_size_mb`, `uploadExpenseAttachment`, 413) sin factor en bytes; se escribió «factor de bytes por MB: lo fija C1; propuesta de S1 1 MB = 1,048,576 bytes, a confirmar en la puerta 2».
- I-02: S4-07 dice que `S4-meta-ajustes` es un bloqueo que S4 abriría si hiciera falta (no existe).
- B-01 de REVISION-C1-r2: C1-13 dice C1 v1.2 / `info.version` 1.2.0; C1-21 dice 210 operaciones (total de `extractos/INDICE.md`; `openapi.yaml` tiene 210 `operationId`).
- U3-01: la Verifica de la tarjeta (4) exige que la semana por defecto siga K38 (domingos, semana que empieza el lunes), citando U1 §0.2.
- Lista de cambios «1.2, ronda 4» en C1, U3, U4, S4, S6, S7 y P3; INDICE: encabezado, §3, §6, §10 (U3 pasa a 1.2) y §11 (filas nuevas). Conteo sin cambio: 261 (201 / 60 / 0).

## Archivos
`docs/criterios/` C1.md, S4.md, S6.md, S7.md, U3.md, U4.md, P3-claude-design.md, INDICE.md; este reporte. No se tocó tools/ ni nada generado.

## Verificación
- `bash docs/criterios/tools/verif_final.sh`: 261 criterios (201/60/0), 18 archivos, 0 sin Verifica/Evidencia/§5; refs 261/261; citas 305, 0 sin fuente; N-xx residuos 0; canon 178 / 0 diferencias; S1-09 SQL OK y equivalencia 12 / 0; matriz 210 operaciones; rutas /tmp 0; «INDICE contra conteo propio: 0 diferencias». (Primera corrida dio 1 diferencia por la versión de U3 en §10; corregida a 1.2.)
- `GOOSE=… make verify-criterios`: «las 11 comprobaciones dan lo esperado», 0 diferencias.
- `grep -rn params_json docs/criterios/`: 12 líneas; un script que busca umbral/días antes/thresholds/days_before en la misma oración solo encuentra C1-22 Verifica diciendo que `params_json` YA NO trae `thresholds_pct` ni `days_before`. Ninguna línea dice que vivan ahí.
- `grep -rn PROVISIONAL docs/criterios/`: marcas activas citan `S1-parametros-avisos` (lugar de los 4) o K31 (valor inicial de `send_time`/`advance_days` en S7); el resto son avisos de «deja de ser PROVISIONAL», reglas generales de K31 o texto histórico. Excepción a notar: S4-07 conserva la contingencia ya existente «sin fuente de los factores intermedios, S4 abre `S4-meta-ajustes` y esos factores quedan PROVISIONAL» (condicional, no marca activa, no cita K31).
- Decimal: 219.88 (S2, INDICE), 629.35 (S2, INDICE), 0.9836 y 0.4918 (S5) presentes y recalculados con ROUND_HALF_UP.

## Supuestos
- P3-04: se usa `docs/estado/wf1/INTEGRACION.md` como ruta absoluta del registro del flujo WF-1 por analogía con `docs/estado/p1c/INTEGRACION.md`; el criterio admite otra ruta si el orquestador la declara.
- U3-01: «el resto de los días, la semana en curso» sale de U1 §0.2 (semana por defecto).
- S7-10: se agregó que el día del mes del job lee `day_of_month` (decisión F, lugar PROVISIONAL) para coherencia con el bloqueo ampliado.

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
Ninguno de la lista de r3 y de los extras. Revisor de ronda siguiente debe reejecutar las comprobaciones.

## Solicitudes fuera de alcance
- `docs/contratos/extractos/ajustes-d.yaml` (operación `updateNotification`) y 00003 dicen «send_time y advance_days… FIRMES (K43)» y a la vez lugar PROVISIONAL; coherente con los criterios, sin cambio pedido.
- N-320…N-327 sin operación en C1-datos.md (B-01 de r1, paso 9 de verif_final) ya figura como cerrado en la matriz (210 operaciones, 0 sin operación): sin acción.
- `datos/tmp/c1r3/` existe y no es mío; no lo toqué.

## Propuestas para Axel
- Ninguna nueva. Pendiente de puerta: responder `S1-parametros-avisos` (opciones A/B/C) y confirmar en la puerta 2 el factor de bytes por MB (propuesta de S1: 1 MB = 1,048,576 bytes).

## Consumo aproximado
~25 acciones de herramienta; lectura ≈ 90 KB; escritura ≈ 35 KB.
