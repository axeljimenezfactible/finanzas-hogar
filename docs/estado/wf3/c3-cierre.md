# c3-cierre · WF-3a cierre · matriz C3 v1.4 (sin firma)

Fecha: 2026-10-12 · Agente: c3-cierre (Sonnet). HEAD `0197dad` + árbol de trabajo de WF-3c. Alcance: `docs/contratos/C3-integracion.md` y este reporte. No hice `git add` ni commit. No toqué código, criterios, scripts ni capturas.

## Hecho
- **HECHO** · Regeneré la matriz con `frontend/scripts/c3-integracion.py --escribir` y entradas de hoy: Playwright con la API simulada (275 ejecutadas, 271 pasan, 4 omitidas, 0 fallan) y Go (`TestContrato_Operaciones`: 221 de 221 PASS). Resultado: implementado sí 42/42, probado sí 42/42, capturado sí 41/42 (Global, con `C3-capturas-global.md`).
- **HECHO** · C3 pasó a **v1.4, sin firma**: encabezado y fuentes (C1 v1.5.1), §3 reescrita, §4 con la reejecución de los flujos, §5 con C3-09, §6 con los comandos de hoy, §7 supuestos 10 a 12, §8 historial.
- **HECHO** · Fila **P-26.02**: 11 operaciones, incluida `updateDashboardSettings` (PATCH `/api/v1/settings/dashboard` en `S02Presupuesto.tsx:82`), Go 11/11, PW 8 (con las 4 pruebas de `parametros-tablero.escritorio.spec.ts`), 4 capturas. Desapareció de la tabla de operaciones sin llamar; queda solo `listPeriods` (excepción PROVISIONAL con cobertura).
- **HECHO** · **F-3** con la aclaración de U1 §4 (condición (c) de la ruta A, PROVISIONAL del orquestador): `flujos.spec.ts:123` (ruta B con la cifra de K46, y la ruta A con la lista de F-2: vacío «Aún no hay precios para comparar» → «Registrar precios») y `pendientes.spec.ts:119` (F-3A3: «Mejor canasta Vallejo $1,309.00 · 26 de 26 artículos con precio», 3 toques).
- **HECHO** · C3-08 con el servidor real, copia aislada nueva por archivo: `flujos.spec.ts` 7 de 7, `pendientes.spec.ts` 4 de 4, `f7-f8.spec.ts` 2 de 2 = 8 de 8 flujos de U1 (los 4 obligatorios incluidos). F-1: 3 toques con K10.
- **HECHO** · C3-09 contra `https://finanzas.local:8443`: «120 cifras leídas en pantalla · diferencias != 0: 0».
- **HECHO** · El script sale con **0**: «Matriz del documento = matriz generada: sí».

## Archivos
- `/home/deepsec/finanzas-hogar/docs/contratos/C3-integracion.md` (v1.4)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/c3-cierre.md`
- Efectos de las corridas de prueba en archivos ya versionados (los reescribe la propia prueba, como en rondas anteriores): `frontend/pruebas-u5/flujos/capturas/*.jpg` y `frontend/dist/` (recompilado por `npm run build`, ignorado por git).

## Verificación
- **C3-01** · `grep -c "^| P-" docs/contratos/C3-integracion.md` = 41 (27 P-xx + 14 P-26.xx) más la fila Global = 42; T de U1 = 41. El encabezado trae las 3 definiciones (§1).
- **C3-02** · «C3-02 pantallas U1: 41 (P-xx 27 + P-26.xx 14); filas: 42»; «pantallas sin fila: 0 []; filas sin operación ni Local: 0 []»; «operaciones: 221; sin pantalla: 5 (excepciones: 5, de ellas K53: 4; sin justificar: 0 [])»; «U1 §3.29 211 N-xxx (3 Local); C1 x-necesidades 208; de U1 sin operación ni Local: 0 []; de C1 sin fila en U1: 0 []». Script, salida 0.
- **C3-03** · «C3-03 frontend: 215 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (routes.go: 221; openapi: 221 operaciones)»; «Operaciones con pantalla que el frontend no llama: 1 ['listPeriods']; sin excepción justificada o sin cobertura: 0 []»; `npm run build` rc 0 (PWA 80 entradas). Script, salida 0.
- **C3-04** · `go test ./... -count=1 -v`: 33 paquetes ok (5 sin pruebas), 1151 PASS, 0 FAIL, 0 SKIP, 221/221 `TestContrato_Operaciones`. Playwright (API simulada): 275 ejecutadas, 271 pasan, 4 omitidas, 0 fallan. Servidor real: 8 flujos y C3-09 (ver C3-08 y C3-09).
- **C3-05** · «C3-05 capturas: 164 esperadas (4 × 41 P-); filas con 4/4: 41; filas en «no»: ['Global']». Las 16 de las subsecciones de P-01 se citan aparte (16/16).
- **C3-06** · `grep -c "| no |"` = 1 (Global, con `docs/estado/bloqueos/C3-capturas-global.md`). Sin otras filas abiertas.
- **C3-07** · no aplica a esta versión: sin firma; la agrega el revisor independiente.
- **C3-08** · 8 de 8 flujos pasan en el servidor real (F-1 `flujos.spec.ts:13`, F-3A0 `:71`, F-2 `:88`, F-3 `:123`, F-4 `:200`, F-5 `:235`, F-6 `:257`; `pendientes.spec.ts` 4 de 4; F-7 `f7-f8.spec.ts:16` y F-8 `:141`). Tabla flujo → pasos → pantallas → resultado en C3 §4.
- **C3-09** · 120 cifras de pantalla, 0 diferencias con `oraculo/model.json` (13 del tablero y P-17; ≥ 10 exigidas). Con la API (según `u5-cierre`): deuda viva $3,895.00 y próximo pago $973.75 (vence 2026-10-30; hogar $598.82, Karen $374.93).
- **C3-10** · nombres y notas en español de México, fechas ISO. **C3-11** · versión v1.4 y fecha; supuestos 10 a 12; bloqueos `C3-ajustes-tablero`, `C3-flujos`, `C3-matriz`, `C3-capturas-global`.

## Supuestos
- La tabla se generó con las capturas que ya estaban en disco (164 de 164, regeneradas por `u5-cierre`).
- Las pruebas citadas son las del repositorio tal cual; «probado» conserva la regla PROVISIONAL del punto 4 de `C3-matriz` (título de la prueba y API simulada).
- F-3 ruta A depende de la condición (c) de U1 §4, PROVISIONAL hasta la puerta 4. Lo mismo el lugar y los rangos de «Parámetros del tablero» (P-26.02).
- `make verify-all` y `verify-wf3` no los corrí: `u5-cierre` informa verde (verify-wf3, verify-all). Mi alcance no cambió nada que los afecte.

## Bloqueos
Ninguno nuevo. Siguen abiertos (puerta 4): `C3-capturas-global`, `C3-matriz`, `C3-flujos` punto 2, `C3-ajustes-tablero` (opción A aplicada, PROVISIONAL).

## Incidente propio, ya corregido
- **HECHO** · Una primera pasada de los flujos de C3-08 la lancé sin `FH_BASE=https://127.0.0.1:8444`, así que escribió en la base de desarrollo (`datos/finanzas-dev.db`, puerto 8443) en vez de en la copia aislada: 2 gastos, 1 compra a crédito, 1 peso, 1 reporte, 1 ticket con 22 precios, una lista de compras, 2 comidas y 1 pago de crédito marcado. Lo detecté porque `c3-09.spec.ts` dio 19 diferencias.
- **HECHO** · Reparación: detuve el servidor, aparté la base contaminada y los 2 archivos que creé (el reporte `reporte-mensual-2026-10-v2.pdf` y un adjunto PNG), y la reconstruí con `finanzas migrate`, `finanzas semilla` y `finanzas semilla --desarrollo` (el procedimiento de `docs/estado/wf3/semilla-dev.md`). Comparé tabla por tabla: 6 gastos, 4 compras a crédito, 72 precios, 12 comidas, 1 presupuesto, 1 partida, 1 meta, 0 reportes, es decir, el estado documentado. Se perdieron las sesiones de la base de desarrollo (hay que volver a entrar) y se reinició `ops/dev-run.sh` (PID nuevo, `/health` ok). Después repetí los flujos con la variable bien puesta (los resultados citados son de esa pasada) y C3-09 dio 0 diferencias. Después de todo eso la base sigue en 6/4/72.
- El orden de la pasada equivocada no invalida la evidencia: no se cita ninguna cifra de ella.

## Pendientes
- Firma del revisor independiente (C3-07) y su nuevo cálculo del SHA-256 tras la última edición.
- Ratificar en la puerta 4 las decisiones PROVISIONALES de arriba.

## Solicitudes fuera de alcance
- `frontend/pruebas-integracion/humo-inicio.spec.ts` (dueño `u5-flujos` o el agente de Inicio): falla con «tarjeta 1 visible» porque espera las 6 tarjetas. Cambio: leer K-1 a K-4 y B-1 del tablero de U3 2.0 y comparar con la API (`getDashboardNextPayment`, `getCreditTotals`).
- `frontend/scripts/c3-integracion.py` (dueño de C3): quitar el diccionario `PENDIENTES_PANTALLA` (ya sin uso, `updateDashboardSettings` tiene pantalla) y la entrada del texto de ayuda.
- `frontend/pruebas-u5/flujos/playwright.config.ts` y `Makefile` (comentario de `verify-u5-servidor`): el comando documentado `bash servidor-prueba.sh arrancar && npx playwright test -c playwright.config.ts` no apunta a la copia aislada si no se define `FH_BASE=https://127.0.0.1:8444`. Cambio: que `servidor-prueba.sh arrancar` imprima la variable, o que la configuración use 8444 cuando existe el pid de la copia, o que las pruebas que escriben se nieguen a correr contra 8443.

## Propuestas para Axel
- Ninguna nueva. Para la puerta 4, ratificar F-3 ruta A (condición (c)) y el bloque «Parámetros del tablero» de P-26.02.

## Paquetes instalados
Ninguno.

## Consumo aproximado
≈ 50 acciones, ≈ 210 KB leídos, ≈ 28 KB escritos. Temporales de `datos/tmp/` (c3, dev-contaminada, u5-flujos, playwright-*) borrados.
