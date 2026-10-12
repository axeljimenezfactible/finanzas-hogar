# revisor-backend-areas-r1 · reporte · 12-oct-2026

## Hecho
- Revisé api-compras-a, api-compras-b, api-creditos, api-ajustes, s7-avisos, s7-pdf y pruebas-contrato contra S6 (S6-03, S6-04, S6-15, S6-16) y contra S7 completo (S7-07 queda para el orquestador). El informe con el veredicto, la tabla y los hallazgos está en `docs/estado/wf3/REVISION-backend-areas-r1.md`.
- Veredicto: **NO APROBADO**, con 4 altos (H-1 a H-4), 1 medio (M-1) y 5 bajos (B-1 a B-5).

## Archivos
- `docs/estado/wf3/REVISION-backend-areas-r1.md` (nuevo)
- `docs/estado/wf3/revisor-backend-areas-r1.md` (este archivo)

Fuera de estos dos, no toqué el repositorio. Las pruebas propias corrieron con `go test -overlay` desde `datos/tmp/rev-ba-r1/`, que ya borré.

## Verificación (comandos propios; salidas resumidas)
- `make verify-wf3-go verify-wf3-contrato` → `go build` y `go vet` limpios. **FALLA gofmt** en `internal/reportes/{reportes_test.go,vistas.go}` (H-3).
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → exit 2. Falla `verify-go` por gofmt, y falla `verify-wf1-cobertura` porque `internal/api` queda en FAIL por `TestContrato_Operaciones`. C1 da TODO OK (`build.py --check` y `verify_c1.py --db`); S1, U1 y C2 dan OK.
- S6-03: el overlay con `Registry.Implemented()` da 109 operaciones (18/23/29/19/10/10), sin duplicados. El overlay de `ValidateResponse` da 109 de 109 validadas en su `TestOp_`. `TestRutasIgualesAlContrato` da «216 operaciones en 172 rutas, 0 diferencias». `TestContrato_Operaciones` filtrada a las 109 da 105 PASS y 4 FAIL.
- S6-04: `limit=201` → 400, `973.75` → 400, error único con mensajes en es-MX. `TestContrato_PaginaMaxima` y `TestContrato_CapturasAtrasadas` dan PASS.
- S6-15: 0 apariciones del bearer falso y 0 del chat completo en los registros. 0 mensajes en inglés.
- S7-01: latexmk lualatex → rc 0, 0 errores y 0 indefinidas en los 2 PDF (2 páginas y 1 página).
- S7-02: `-run Cruce` da 36 de 36 con 0 diferencias. **R-47 a R-50 no se cruzan** (H-2).
- S7-03: 20 de 20. S7-04, S7-05, S7-08 y S7-09 dan PASS, y además las probé con llamadas mías (Enviador falso y HOME falso).
- S7-06: el grep del patrón de token en todo el repositorio da **0**.
- S7-10 y S7-11: dan PASS con la base de prueba, pero dan **500 con una categoría creada por la API** (H-1).
- S7-12: pdftotext da 0 fechas ISO y 0 cifras fuera de formato.
- S7-13: con el script de U2 todos los pares dan ≥ 4.5 y los 7 colores son de C2.
- S7-14: la única red es `api.telegram.org`. LaTeX corre con tiempo límite de 120 s y `go vet` sale limpio. gofmt falla (H-3).
- 5 operaciones al azar por área (semilla `20261012`):
  - comprasa: `updateStore`, `getBasketComparison`, `updatePresentation`, `updateProduct`, `getProduct`.
  - comprasb: `updatePrice`, `getInventorySummary`, `getPriceSettings`, `listInventoryItems`, `importProfeco`.
  - creditos: `getDebtTrend`, `updateCredit`, `deleteCreditTermRate`, `getOptimizerSettings`, `updateCreditTermRate`.
  - ajustes: `listCostShareRules`, `updatePerson`, `updateCostShareRule`, `listCategories`, `updateProjectionSettings`.
  - avisos: `updateNotificationRecipient`, `sendTestNotification`, `updateNotification`, `listNotificationRecipients`, `setNotificationRecipients`.
  - reportes: `updateReportSettings`, `listMonthlyReportMonths`, `sendReport`, `getExternalStatement`, `listReports`.
  - Las 30 responden con un código declarado y un esquema válido (2xx y 4xx propios).
- K53: `importExternalPrices` sin servicio da 401 «Esta operación necesita la credencial del agente externo.». Las 4 operaciones del agente externo llevan `x-provisional: true` en `compras-tiendas-b.yaml`.
- Comparador: `creditos/condiciones.go:663-673` lee `optimizer.date_window_days` y `comparator.timeout_seconds` con `dominio.SettingInt`, y aplica el tiempo con `comparador.TimeoutFromSetting`. `getOptimizerSettings` → 200 con los dos parámetros.

## Supuestos
- Para el veredicto conté como «prueba que falla» que `TestContrato_Operaciones` falle en operaciones de los agentes revisados, y también que se rompa `make verify-all`.
- B-4 no lo cuento como incumplimiento de S7-08, porque el criterio pide fixtures; queda como riesgo de integración.

## Bloqueos
Ninguno nuevo.

## Pendientes
- La ronda 2, después de las correcciones de H-1 a H-4 y M-1.

## Propuestas para Axel
- Incluir `TestContrato_Operaciones` en `verify-wf3-contrato`.
- Asignar H-1 a api-presupuesto (la causa) y a s7-pdf (la prueba).

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 60 acciones, unos 150 KB leídos y unos 30 KB escritos (incluidos los temporales ya borrados).
