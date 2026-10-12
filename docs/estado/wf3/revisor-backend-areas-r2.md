# revisor-backend-areas-r2 · reporte · 12-oct-2026

Ronda 2 de 3, como revisor independiente. No corregí nada. El informe completo está en `docs/estado/wf3/REVISION-backend-areas-r2.md`.

## Hecho

- **Veredicto: NO APROBADO.** Hay 1 alto, 0 medios y 5 bajos.
- Los 5 hallazgos de la ronda 1 están **cerrados**: H-1 (`category.code` NULL), H-2 (cruce P1–P6), H-3 (gofmt), H-4 (`setStoreProductGroups` con `{}`) y M-1 (`sendReport` 202). Cada uno lo reproduje con mis propios comandos.
- El hallazgo alto nuevo es H-1: `getDashboardBestBasket` sigue respondiendo 501. Por eso falla S6-03 en el servidor completo y `make verify-all` (Makefile:233). Para cerrarlo hacen falta dos cosas: que api-compras-a exporte `BestBasket` y que api-gastos-inicio la registre.

## Archivos

- `docs/estado/wf3/REVISION-backend-areas-r2.md` (informe)
- `docs/estado/wf3/revisor-backend-areas-r2.md` (este reporte)
- Temporales en `datos/tmp/rev-ba-r2/`, ya borrados: mi prueba que entra por overlay, `overlay.json`, la base `s6-c1.db`, la compilación LaTeX, los PNG y los scripts de conteo.

## Verificación

- **S6-03**: `build.py --check` → 23 archivos idénticos. `verify_c1.py --db` → TODO OK. Mi `TestRevR2_AreasOps` → 109 de 109 operaciones con `TestOp_`. `TestContrato_Operaciones` → 215 PASS y 1 FAIL (`getDashboardBestBasket`). `TestContrato_Rutas` → «sin implementar 1». **FALLA global (H-1).**
- **S6-04**: `limit=201` → 400 en 6 colecciones. Idempotency-Key: corta → 400; repetida → `Idempotent-Replayed: true`; con otro cuerpo → 409. `973.75` → 400. Los `IdList` con `{}` → 422. CUMPLE.
- **S6-15**: `-run Log` → 2 PASS. 20 mensajes al azar, todos en es-MX. 0 tokens en el registro. CUMPLE.
- **S6-16**: hay `## Supuestos` en los 7 reportes y 19 bloqueos `S6-*`. CUMPLE (queda B-3).
- **S7-01**: `latexmk` desde cero → rc 0, 0 «!», 0 «undefined», 2 + 1 páginas, 3 PNG. CUMPLE.
- **S7-02**: `-run Cruce` → R-45…R-50 = oráculo, 256 cifras y fechas en 6 PDF, 0 diferencias. CUMPLE.
- **S7-03**: `-run Karen` → 20 de 20. CUMPLE.
- **S7-04**: `-run Resumen` → PASS. Grep de horas literales → 0. CUMPLE.
- **S7-05**: `-run` de `internal/telegram` → PASS, sin `parse_mode` y con trozos ≤ 4096. CUMPLE.
- **S7-06**: grep del patrón de token en todo el repositorio → **0**. `TestToken` PASS. CUMPLE.
- **S7-07**: no se evalúa (lo hace el orquestador).
- **S7-08**: `-run Alerta` → la matriz de 11 fechas coincide y la configuración se restablece. CUMPLE.
- **S7-09**: `-run Destinatarios` → PASS. Karen → 422. CUMPLE.
- **S7-10**: `-run Mensual` → PASS. Mi job con una categoría de usuario → `ok`. CUMPLE (queda B-5).
- **S7-11**: `-run VistaMensual` → 34 de 34. CUMPLE.
- **S7-12**: 0 fechas ISO, 0 cifras fuera del patrón, 0 palabras en inglés. CUMPLE. El signo negativo es una PREGUNTA ya abierta.
- **S7-13**: 7 de 7 colores son de C2. Contraste mínimo 5.04, calculado con el script de U2. CUMPLE.
- **S7-14**: vet limpio. Cobertura: reportes 77.3 %, reportes/avisos 77.3 %, telegram 80.0 %. Única URL: `api.telegram.org`. gofmt vacío. CUMPLE.
- **K53**: las 4 operaciones tienen `x-provisional`. Bearer → 200 y cookie → 401. Repetir el lote no duplica. `source = site_assisted`.
- **Comparador**: lee `optimizer.date_window_days` de `app_setting`. Con 7 evalúa 15 escenarios, con 0 evalúa 1, y al restablecer vuelve a 15. `comparator.timeout_seconds` = 99 → 422.
- **Global**: `make -k verify-all GOOSE=…` → un solo Error 1, en Makefile:233 (H-1).

## Supuestos

- Para S6-03 se evalúa el servidor completo, como dice el criterio («el conjunto de rutas registradas… idéntico»). Por eso cuento como alta la operación en 501. Es de api-gastos-inicio, pero depende de una exportación de api-compras-a, que es un agente de esta revisión. El revisor de la ronda 1 la dejó fuera de alcance. Yo la incluyo porque es lo único que deja S6-03 y `verify-all` en rojo. La marco como compartida con la revisión backend-core.
- Los destinatarios y los grupos de producto de mis pruebas los sembré por SQL en una base temporal, como hacen las pruebas de los agentes. Con datos reales no toqué nada.

## Bloqueos

Ninguno nuevo.

## Pendientes

Son de los dueños, no míos:
- H-1: api-compras-a y api-gastos-inicio.
- B-1: s7-pdf.
- B-2 y B-4: s7-avisos.
- B-3: api-creditos.
- B-5: integrador.

## Propuestas para Axel

- PREGUNTA 4 de `S7-reporte-mensual-secciones`: ¿cómo se escriben los negativos en el PDF y en los mensajes: `-$`, `$-` o `−$` (U+2212, como en DESIGN.md)?

## Paquetes instalados

Ninguno.

## Consumo aproximado

58 acciones, unos 170 KB leídos y unos 30 KB escritos (incluidos los temporales borrados).
