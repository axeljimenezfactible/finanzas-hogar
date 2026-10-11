# Reporte WF-0 · S6a-contrato-c1

Fase S6 parte a «Contrato de datos C1 y esqueleto compilable del backend» · flujo WF-0 «Cimientos» · 10-oct-2026 · Estado: **terminado** (sin cierre por tope).

> **Corrección r1 (10-oct-2026, agente `correccion-C1-r1`, hallazgo H-05 de `REVISION-r1.md`).** El generador y el verificador de C1 ya están en el repositorio, en `backend/internal/api/c1gen/` (con `build.py`, `build.py --check` y `README.md`). `openapi.yaml`, `routes.go` y `C1-datos.md` se regeneraron con ellos y solo cambiaron comentarios de cabecera y dos filas de documentación. Las notas «(corrección r1)» de este reporte señalan lo que cambió; el resto queda como lo entregó S6a. Evidencia: `docs/estado/wf0/correccion-C1-r1.md`.

> **Corrección r2 (11-oct-2026, agente `correccion-C1-r2`, hallazgos R2-01 y R2-02 de `REVISION-r2.md`).** C1 se regeneró contra el esquema 1.1 de S1 (54 tablas): 158 rutas, 197 operaciones y 237 esquemas. La lista de compras tiene sus operaciones (N-151, N-152, N-153 y N-165 ya no están bloqueadas; todas PROVISIONALES por la opción A de `C1-lista-compras`), y el presupuesto, las recetas, el gasto, el inventario y el ingreso esperado siguen las columnas nuevas. Las cifras y tablas de este reporte son la foto de C1 v1.0 en WF-0; el estado vigente está en `docs/contratos/C1-datos.md` (§10, §11, §13 y §14) y en `docs/estado/wf0/correccion-C1-r2.md`.

## Hecho

- **`docs/contratos/openapi.yaml`**: OpenAPI **3.1.1**, contrato **C1 v1.0** (`info.version` 1.0.0): **154 rutas, 193 operaciones, 232 esquemas**. Cubre autenticación (login, logout, sesión actual, cambio de contraseña, sesiones por dispositivo y su cierre), personas y reparto, periodos, ingresos fijos y extra, partidas planificadas e importación del CSV, gastos con reparto, renglones y adjuntos, captura de ticket con OCR y confirmación, categorías, presupuesto y flujo, metas y fondo de ahorro, créditos de los 4 tipos y sus condiciones (tasas por plazo, interruptor A/B con vista previa, serie de disponible, estado de tarjeta), compras a crédito con partes por persona, pagos programados y marcado por parte, promociones, comparador de escenarios (entrada: compra planeada y restricciones; salida: los 4 óptimos y la recomendación ponderada), tiendas (cadenas, sucursales, mercados con día de visita, membresía), productos, presentaciones, precios con fuente y frescura, comparador de canasta, recetas e ingredientes, tiempos de comida, plan de comidas, lista de compras, inventario, perfil de salud y peso, nutrición, notificaciones y destinatarios, reportes (generar, listar, descargar PDF, reenviar, vista web, estado de cuenta de Karen), jobs, respaldos, IA opcional, estado del sistema, las 14 secciones de Ajustes y el endpoint genérico de **vistas de base de datos** (`GET /api/v1/views/{entity}` con modo tabla, tablero, calendario o galería, filtros, orden, agrupación y paginación).
- Convenciones únicas: dinero entero en centavos (`MoneyCents`, sufijo `_cents`, 253 propiedades), fechas `Date`/`DateTime` ISO con desplazamiento de America/Mexico_City, **un** esquema de error (`Error`, 874 respuestas de error que lo referencian), paginación por cursor con la envoltura `Page` (39 colecciones), 6 filtros globales con definición única, `Idempotency-Key` obligatoria en los 38 POST de captura (C1 §1 punto 8 lista los 10 POST que no la llevan y por qué), cookie `fh_session` (HttpOnly, Secure, SameSite=Strict), solo login y salud públicos. Cada operación declara `x-tablas` (tablas de S1) y `x-necesidades` (N-xxx de U1); cada métrica del dashboard lleva su fórmula en la descripción de su esquema.
- **`docs/contratos/C1-datos.md`** (C1 v1.0): convenciones, autenticación, modelo de errores con tabla de códigos, fechas y dinero, paginación, filtros globales con matriz recurso × filtro, las 31 entidades → recurso → operaciones → vistas, endpoint de vistas con agrupaciones y fechas por entidad, **17 métricas** (M-01 a M-17) con fórmula, periodo, filtros, tablas, endpoint y ejemplo del oráculo (R-nn de `docs/criterios/S3.md`), la **matriz completa pantalla P-xx → necesidad N-xxx → endpoint** (220 filas, 41 pantallas), cobertura (necesidades, tablas, 14 secciones de Ajustes, 35 capacidades), índice de las 193 operaciones, lista de PROVISIONAL con su bloqueo e inconsistencias S1 ↔ U1.
- **Backend compilable**: `backend/go.mod` (`module finanzas-hogar`, `go 1.26`), `backend/cmd/finanzas/main.go` (servidor HTTP con `GET /health` y `GET /openapi.yaml` reales; 404 y 405 también con el esquema único de error; escucha en 127.0.0.1:8080 por defecto) y `backend/internal/api/routes.go` (las 193 operaciones del contrato registradas en `net/http.ServeMux` como stubs que responden **501** `not_implemented` con el esquema de error y `X-Request-Id`). `go build ./...`, `go vet ./...` y `gofmt -l` limpios; el servidor arranca sin conflicto de patrones.
- 6 bloqueos nuevos con el formato del protocolo (5 de C1 y 1 de S6).
- Todo se generó con un generador propio (esquema de S1 → esquemas de entidad con los mismos `nombre_codigo`; tabla de necesidades de U1 → matriz), para que contrato, rutas y matriz no puedan divergir. Ver «Solicitudes» punto 1. (corrección r1) El generador está en `backend/internal/api/c1gen/`.

## Archivos

Creados (todos dentro de mi alcance):
- `/home/deepsec/finanzas-hogar/docs/contratos/openapi.yaml` (570 KB, generado)
- `/home/deepsec/finanzas-hogar/docs/contratos/C1-datos.md` (106 KB)
- `/home/deepsec/finanzas-hogar/backend/go.mod`
- `/home/deepsec/finanzas-hogar/backend/cmd/finanzas/main.go`
- `/home/deepsec/finanzas-hogar/backend/internal/api/routes.go` (generado desde el contrato)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-filtro-familia.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-lista-compras.md` (corrección r2) Ya no hay necesidades sin endpoint: el esquema 1.1 trae las tablas y C1 implementa la opción A como PROVISIONAL; la pregunta es la misma de `S1-lista-compras`.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-huecos-u1-s1.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-parametros-tecnicos.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-campos-calculados.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S6-version-go.md`
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/S6a-contrato-c1.md` (este reporte)

Fuera del repositorio (scratchpad de la sesión, efímero): el generador y las herramientas en `/tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/c1/` (`parse_schema.py`, `gen1_base.py` … `gen6_c1.py`, `c1_plantilla.md`, `c1_inconsistencias.md`, `verify_c1.py`, `u1_verify.py`). No modifiqué ningún archivo de otros agentes.

(corrección r1) Esos archivos ya están en el repositorio: `backend/internal/api/c1gen/` (`parse_schema.py`, `gen1_base.py` … `gen6_c1.py`, `c1_plantilla.md`, `c1_inconsistencias.md`, `verify_c1.py`, más los nuevos `build.py`, `c1paths.py` y `README.md`) y `backend/internal/api/generate.go` (`go generate`). `u1_verify.py` es de U1 y no se copió a C1 (ver solicitudes de `correccion-C1-r1.md`).

## Verificación

Script propio `verify_c1.py` (anexo A), ejecutado con `uv run --no-project --with pyyaml --with jsonschema python3 -I verify_c1.py` → **RESULTADO: TODO OK**.

(corrección r1) El script está ahora en `backend/internal/api/c1gen/verify_c1.py` y se llama desde cualquier carpeta. Además revisa la consistencia de la matriz §10 y sigue el texto vigente de C1-02, en el que una necesidad bloqueada no cuenta como cubierta. Por eso hoy reporta **1 falla**: N-151, N-152, N-153 y N-165, del bloqueo `C1-lista-compras` (hallazgo H-01, de S1). Las salidas nuevas están en `correccion-C1-r1.md` §Verificación.

- **C1-01**: `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` → `docs/contratos/openapi.yaml: OK`, código de salida 0. `npx @redocly/cli@2.62.1 lint docs/contratos/openapi.yaml` → «Your API description is valid», 0 errores, **1 advertencia `info-license`** (justificada: proyecto privado del hogar sin licencia; no se inventa una). Carga YAML: `openapi = 3.1.1 · rutas = 154 · operaciones = 193`. El comando literal del criterio (`python3 -I -m openapi_spec_validator` y `python3 -I -c "import yaml…"`) no corre con el Python del sistema porque ni `openapi-spec-validator` ni PyYAML están instalados ahí (Arch bloquea pip global); usé los mismos paquetes en entornos aislados de uv (ver Solicitudes punto 4).
- **C1-02**: `N-xxx en U1: 199 · con endpoint: 193 · N-xxx sin endpoint: 6 ['N-003', 'N-151', 'N-152', 'N-153', 'N-165', 'N-301']` (N-003 y N-301 son locales según U1; las otras 4 están en el bloqueo `C1-lista-compras`) · `operaciones: 193 · operaciones sin N-xxx: 1 ['GET /openapi.yaml']` (excepción estática justificada en C1 §11.2) · rutas de la matriz inexistentes: 0. Cruce con el script de U1 (`u1_verify.py`, anexo de su reporte): `FALLA U1-03 N-xxx sin endpoint en C1: ['N-151', 'N-152', 'N-153', 'N-165']` y `FALLA U1-03 operaciones de openapi.yaml sin N-xxx: ['GET /openapi.yaml']`. Son exactamente las mismas excepciones justificadas (su script no tiene lista de excepciones).
  (corrección r2) Hoy: `N-xxx sin operación ni Local: 0`; solo N-003 y N-301 quedan sin operación, por ser Local según U1.
- **C1-03**: `tablas en sqlite_master: 49 · citadas por x-tablas: 49 · inexistentes: [] · tablas sin operación: []`; 4 operaciones con `x-tablas` vacía porque leen un archivo o el proceso, no la base (`getHealth`, `getOpenApiDocument`, `getSystemStatus`, `getCaCertificate`). `propiedades de esquemas de entidad: 906 · sin columna ni marca de calculado/entrada: 0`.
- **C1-04**: `respuestas de error: 874 · con referencia al esquema único Error: 874 · excepciones: [] · esquemas de error alternativos: []`.
- **C1-05**: `colecciones GET: 39 · con cursor+limit+sort y envoltura Page: 39 · diferencias: []`.
- **C1-06**: `usos por filtro: {'person_id': 12, 'family': 8, 'category_id': 7, 'store_id': 8, 'period': 32, 'credit_id': 11} · definidos inline: [] · sinónimos: []` → 6 de 6 con definición única. La matriz recurso × filtro está en C1 §6.
- **C1-07**: `formatos usados: ['date', 'date-time'] · ejemplos con fecha: 2 · fuera de ISO: []`; C1 §4 declara quincenas 15 y 30 y que en febrero el «30» es el 28 (`2027-02-28`, R-23).
- **C1-08**: `propiedades de dinero (_cents): 253 · sin referencia a MoneyCents: [] · type: number en el archivo: 0`.
- **C1-09**: 17 métricas en C1 §9, las 6 tarjetas de U3 (N-020 a N-025) con su métrica (M-01 a M-07). Cifras de ejemplo buscadas por script: 1,197.64 · 598.82 · 374.93 · 749.86 · 973.75 · 1,258.70 · 1,947.50 · 3,094.50 · 5,000.00 · 5,034.74 · 6,293.44 · 6,905.50 · 8,094.50 · 10,000.00 · 13,103.78 · 14,362.48 en `oraculo/model.json`; 3,233.00 en `datos/semilla/data.json` (`totales.walmart`) y 138.50 = su diferencia (plan §3.4 fila S5). 0 cifras sin origen. Las que el oráculo no trae (plan del presupuesto, ingreso, comidas, ahorro) van como marcador.
- **C1-10**: tabla de las 31 entidades → tablas → operaciones → vistas en C1 §7 (31 de 31; 8 sin lista propia con el motivo de U1) y agrupaciones y fechas por entidad en C1 §8.
- **C1-11**: 35 capacidades → operationId en C1 §11.6 y las 14 secciones de Ajustes → necesidades → operaciones en C1 §11.5 (P-26.10 Apariencia sin API, como dice U1).
- **C1-12**: `operaciones públicas: ['GET /health', 'POST /api/v1/auth/login'] · campos sensibles en respuestas: [] · operaciones de alta de cuentas o personas: []`; cookie `fh_session` con HttpOnly, Secure y SameSite descrita; login describe 5 intentos en 15 minutos y responde 429 `too_many_attempts` con `retry_at`.
- **C1-13**: `operationId duplicadas: 0 · sin etiqueta válida: [] · etiquetas: Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes, Ajustes · info.version: 1.0.0`; «C1 v1.0» y la tabla de cambios en C1.
- **C1-14**: `ejemplos validados con jsonschema: 13 · inválidos: 0`; `descripciones y resúmenes: 2607 · con palabras frecuentes del inglés fuera de código: 0`; el script imprime una muestra fija de 20 descripciones para el revisor.
- **C1-15**: `marcas PROVISIONAL en openapi.yaml: 63 · sin bloqueo existente: [] · líneas PROVISIONAL en C1-datos.md: 44 · sin bloqueo: [] · bloqueos citados en C1: 18 · citados que no existen: []`.
- **S6-01 (parcial)**: `go version go1.26.5-X:nodwarf5 linux/amd64`; `grep -n '^go ' go.mod` → `3:go 1.26`; `gofmt -l .` sin salida; `go build ./...` y `go vet ./...` sin diagnósticos. La versión aprobada (1.27) no está disponible: bloqueo `S6-version-go`.
- **S6-03 (parte de rutas)**: `rutas en routes.go: 193 · en openapi.yaml: 193 · sobrantes: [] · faltantes: []`. Prueba del binario en 127.0.0.1:18089: `GET /health` → 200 `{"status":"ok","time":"2026-10-10T23:00:55-06:00"}`; `GET /openapi.yaml` → 200 `application/yaml`; `GET /api/v1/expenses` → 501 `{"code":"not_implemented","message":"Esta operación todavía no está implementada.","request_id":"…"}`; ruta inexistente → 404 `not_found` y `DELETE /health` → 405 `method_not_allowed`, ambos con el esquema de error. El servidor se detuvo al terminar (puerto libre).

## Cobertura

Toda necesidad N-xxx de U1 → operación(es) (detalle por pantalla en C1 §10) y toda tabla de S1 → recurso (detalle en C1 §11.4). Generado del contrato:

### Necesidades N-xxx → operaciones (199)

(corrección r2) Esta tabla es la de C1 v1.0 en WF-0. La vigente es la matriz §10 de `docs/contratos/C1-datos.md`: N-151, N-152, N-153 y N-165 tienen operación y quedan PROVISIONALES por `C1-lista-compras`.

| N-xxx | Operaciones (operationId) |
|---|---|
| N-001 | `getCurrentSession` |
| N-002 | `search` |
| N-003 | Local (U1), sin API |
| N-004 | `getHealth` |
| N-005 | `listCredits` |
| N-006 | `logout` |
| N-010 | `login` |
| N-011 | `login` |
| N-020 | `getDashboardMonthVsPlan` |
| N-021 | `getDashboardFortnightFlow` |
| N-022 | `getDashboardNextPayment` |
| N-023 | `getDashboardWeekMeals` |
| N-024 | `getDashboardBestBasket` |
| N-025 | `getDashboardSavings` |
| N-026 | `getPersonParticipation` |
| N-027 | `getBudgetByCategory`, `getCashFlowFortnights`, `listSavingsGoals`, `getMealPlanWeek`, `getBasketComparison` |
| N-028 | `getBudgetMonthlyTrend`, `getAccrualVsCash`, `getMealPlanTrends`, `getDebtTrend` |
| N-029 | `listCalendarEvents` |
| N-040 | `listExpenses`, `listPersons` |
| N-041 | `getExpenseTotals` |
| N-042 | `getExpenseSummary` |
| N-043 | `queryView` |
| N-044 | `getExpenseMonthlyTrend` |
| N-045 | `queryView` |
| N-046 | `queryView`, `downloadAttachment` |
| N-047 | `updateExpense` |
| N-050 | `listCategories` |
| N-051 | `listCategories`, `listCostShareRules` |
| N-052 | `getCurrentSession`, `listCostShareRules` |
| N-053 | `listCredits` |
| N-054 | `createExpense` |
| N-055 | `deleteExpense` |
| N-060 | `getExpense`, `downloadAttachment` |
| N-061 | `listExpenseItems` |
| N-062 | `updateExpense` |
| N-063 | `listExpenses`, `deleteExpense` |
| N-064 | `uploadExpenseAttachment`, `deleteAttachment` |
| N-065 | `markExpensePaid` |
| N-066 | `duplicateExpense` |
| N-070 | `createTicketScan` |
| N-071 | `downloadAttachment`, `getTicketScan` |
| N-072 | `search`, `listProducts` |
| N-073 | `getShoppingList` |
| N-074 | `confirmTicketScan` |
| N-075 | `createTicketScan` |
| N-076 | `listTicketScans` |
| N-080 | `listPeriods`, `getBudgetMonth` |
| N-081 | `getBudgetByCategory` |
| N-082 | `listPeriods`, `getBudgetFortnights` |
| N-083 | `getAccruedGroceries` |
| N-084 | `getBudgetAlerts` |
| N-085 | `getBudgetMonthlyTrend` |
| N-086 | `listCalendarEvents` |
| N-087 | `createBudget`, `updateBudget` |
| N-090 | `getCashFlowFortnights` |
| N-091 | `getAccrualVsCash` |
| N-092 | `getLiveDebt` |
| N-093 | `getCashFlowAssumptions` |
| N-100 | `listPeriods`, `listIncomes`, `getIncomeMonth` |
| N-101 | `createIncome` |
| N-102 | `createIncome` |
| N-103 | `getIncomeMonthlyTrend` |
| N-104 | `listCalendarEvents` |
| N-105 | `listIncomes`, `updateIncome`, `deleteIncome` |
| N-110 | `getSavingsSummary` |
| N-111 | `queryView`, `listSavingsGoals` |
| N-112 | `listSavingsMovements` |
| N-113 | `createSavingsMovement` |
| N-114 | `getSavingsMonthlyTrend` |
| N-115 | `listCalendarEvents` |
| N-116 | `listSavingsMovements`, `updateSavingsMovement`, `deleteSavingsMovement` |
| N-120 | `listMealSlots`, `getMealPlanWeek` |
| N-121 | `getMealPlanWeek` |
| N-122 | `listRecipes` |
| N-123 | `createMealPlanEntries`, `updateMealPlanEntry` |
| N-124 | `createMealPlanEntries` |
| N-125 | `deleteMealPlanEntry` |
| N-126 | `copyPreviousWeek` |
| N-127 | `getShoppingList` |
| N-128 | `queryView`, `getMealPlanMonth` |
| N-129 | `getMealPlanTrends` |
| N-135 | `queryView`, `listRecipes` |
| N-136 | `createRecipe` |
| N-137 | `suggestRecipeWithAi` |
| N-140 | `getRecipe` |
| N-141 | `getRecipeUsage` |
| N-142 | `updateRecipe` |
| N-143 | `duplicateRecipe` |
| N-144 | `deleteRecipe` |
| N-150 | `getShoppingList`, `listProductGroups` |
| N-151 | BLOQUEO `C1-lista-compras` |
| N-152 | BLOQUEO `C1-lista-compras` |
| N-153 | BLOQUEO `C1-lista-compras` |
| N-160 | `listBaskets` |
| N-161 | `getBasketComparison` |
| N-162 | `getBasketComparison` |
| N-163 | `getBasketComparison` |
| N-164 | `getBasketComparison` |
| N-165 | BLOQUEO `C1-lista-compras` |
| N-170 | `queryView`, `listProducts` |
| N-171 | `queryView`, `getProduct`, `listProductPrices` |
| N-172 | `getPriceSummary` |
| N-173 | `getProduct`, `listNutritionFacts` |
| N-174 | `createPrice` |
| N-175 | `createProduct`, `updateProduct`, `createPresentation`, `updatePresentation`, `createNutritionFact` |
| N-176 | `updatePrice`, `deletePrice` |
| N-180 | `queryView`, `listInventoryItems` |
| N-181 | `listInventoryMovements` |
| N-182 | `getInventorySummary` |
| N-183 | `adjustInventoryItem` |
| N-184 | `createInventoryItem` |
| N-190 | `listCredits` |
| N-191 | `getCreditTotals` |
| N-192 | `listCalendarEvents`, `listScheduledPayments` |
| N-193 | `getCredit` |
| N-194 | `getDebtTrend` |
| N-200 | `getCredit` |
| N-201 | `getCredit`, `previewAvailability` |
| N-202 | `queryView`, `listCreditPurchases` |
| N-203 | `getAvailabilitySeries` |
| N-204 | `listCreditTermRates` |
| N-205 | `getCardStatus` |
| N-206 | `listScheduledPayments` |
| N-207 | `listScheduledPayments` |
| N-208 | `listCalendarEvents` |
| N-209 | `getDebtTrend`, `getAvailabilitySeries` |
| N-210 | `markScheduledPaymentPaid`, `deleteCreditPayment` |
| N-215 | `listCreditPurchases`, `getCreditPurchase`, `listScheduledPayments` |
| N-216 | `previewCreditPurchase` |
| N-217 | `previewCreditPurchase`, `listPromotions` |
| N-218 | `createCreditPurchase` |
| N-219 | `updateCreditPurchase`, `deleteCreditPurchase` |
| N-225 | `getOptimizerSettings` |
| N-226 | `compareScenarios` |
| N-230 | `listHealthProfiles` |
| N-231 | `getMealPlanWeek`, `getTodayKcal` |
| N-232 | `listWeightEntries` |
| N-233 | `getMealPlanTrends` |
| N-240 | `createWeightEntry` |
| N-241 | `listWeightEntries` |
| N-242 | `updateWeightEntry`, `deleteWeightEntry` |
| N-245 | `listReports`, `listNotificationDeliveries` |
| N-246 | `listJobs` |
| N-247 | `queryView`, `listNotificationDeliveries` |
| N-248 | `generateMonthlyReport` |
| N-249 | `downloadReportFile` |
| N-250 | `sendReport` |
| N-255 | `listMonthlyReportMonths` |
| N-256 | `getMonthlyReport` |
| N-257 | `listReports`, `getMonthlyReport` |
| N-260 | `listScheduledPayments`, `getExternalStatement` |
| N-261 | `generateExternalStatement` |
| N-262 | `listReports` |
| N-270 | `listUsers` |
| N-271 | `listSessions` |
| N-272 | `revokeSession`, `revokeSessions` |
| N-273 | `changeOwnPassword` |
| N-274 | `getCaCertificate` |
| N-275 | `listBudgets` |
| N-276 | `listPlannedItems` |
| N-277 | `listPlannedItems` |
| N-278 | `createBudget`, `updateBudget`, `deleteBudget`, `createPlannedItem`, `updatePlannedItem` |
| N-279 | `listPlannedItems` |
| N-280 | `createPlannedItem`, `updatePlannedItem` |
| N-281 | `queryView`, `listCategories`, `listCostShareRules` |
| N-282 | `createCategory`, `updateCategory`, `reorderCategories` |
| N-283 | `listCredits` |
| N-284 | `getCredit` |
| N-285 | `previewAvailability` |
| N-286 | `listCreditTermRates` |
| N-287 | `queryView`, `listPromotions` |
| N-288 | `getOptimizerSettings` |
| N-289 | `updatePerson`, `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`, `updateOptimizerSettings` |
| N-290 | `queryView`, `listStores`, `listProductGroups` |
| N-291 | `getPriceSources`, `getPriceSettings`, `listJobs` |
| N-292 | `createStore`, `updateStore`, `setStoreProductGroups`, `updatePriceSettings` |
| N-293 | `importProfeco` |
| N-294 | `listMealSlots`, `getRecipeSettings` |
| N-295 | `reloadNutritionBase` |
| N-296 | `listJobs`, `listNotifications` |
| N-297 | `updateJob`, `updateNotification`, `setNotificationRecipients` |
| N-298 | `listNotifications` |
| N-299 | `listNotificationRecipients`, `createNotificationRecipient`, `updateNotificationRecipient` |
| N-300 | `updateNotification`, `setNotificationRecipients`, `sendTestNotification` |
| N-301 | Local (U1), sin API |
| N-302 | `listJobs`, `listBackups` |
| N-303 | `createBackup`, `restoreBackup` |
| N-304 | `previewPlannedItemsImport`, `importPlannedItems` |
| N-305 | `getSeedStatus` |
| N-306 | `listAiProviders` |
| N-307 | `updateAiProvider` |
| N-308 | `getSystemStatus` |
| N-309 | `listPersons`, `listCostShareRules` |
| N-310 | `updatePerson`, `updateCostShareRule` |
| N-311 | `listHealthProfiles` |
| N-312 | `listSavingsGoals` |
| N-313 | `createSavingsGoal`, `updateSavingsGoal`, `createHealthProfile`, `updateHealthProfile` |
| N-314 | `listProducts`, `listInventoryItems` |
| N-315 | `createInventoryItem`, `updateInventoryItem` |

### Tablas de S1 → recurso (49 de 49)

(corrección r2) El esquema 1.1 de S1 tiene 54 tablas; las 54 tienen operación en `C1-datos.md` §11.4 (`idempotency_key` como infraestructura de S6).

| Tabla | Recurso | Operaciones con la tabla en `x-tablas` |
|---|---|---|
| `person` | `/api/v1/persons` | 46 |
| `user_account` | `/api/v1/users` | 7 |
| `device_session` | `/api/v1/sessions` | 9 |
| `login_attempt` | `/api/v1/auth/login` (límite de intentos) | 1 |
| `app_setting` | `/api/v1/settings/optimizer`, `/api/v1/settings/prices`, `/api/v1/ai/providers` | 16 |
| `audit_log` | bitácora que escriben todas las altas, cambios y bajas; se lee en `getExpense` («creado por») y `listInventoryMovements` | 85 |
| `period` | `/api/v1/periods` | 12 |
| `cost_share_rule` | `/api/v1/cost-share-rules` | 21 |
| `cost_share_rule_part` | partes de `/api/v1/cost-share-rules` | 14 |
| `category` | `/api/v1/categories` | 33 |
| `budget` | `/api/v1/budgets` | 12 |
| `store` | `/api/v1/stores` | 32 |
| `product_group` | `/api/v1/product-groups` | 6 |
| `store_product_group` | `/api/v1/stores/{store_id}/product-groups` | 3 |
| `product` | `/api/v1/products` | 22 |
| `presentation` | `/api/v1/products/{product_id}/presentations`, `/api/v1/presentations/{presentation_id}` | 24 |
| `credit` | `/api/v1/credits` | 31 |
| `credit_term_rate` | `/api/v1/credits/{credit_id}/term-rates` | 16 |
| `promotion` | `/api/v1/promotions` | 12 |
| `planned_item` | `/api/v1/planned-items` | 18 |
| `income` | `/api/v1/incomes` | 13 |
| `expense` | `/api/v1/expenses` | 40 |
| `expense_item` | `/api/v1/expenses/{expense_id}/items` | 18 |
| `expense_share` | partes de `/api/v1/expenses` (se congelan al guardar) | 27 |
| `attachment` | `/api/v1/expenses/{expense_id}/attachments`, `/api/v1/ticket-scans` | 18 |
| `price` | `/api/v1/prices`, `/api/v1/products/{product_id}/prices` | 23 |
| `credit_purchase` | `/api/v1/credit-purchases` | 46 |
| `credit_purchase_share` | partes de `/api/v1/credit-purchases` | 16 |
| `scheduled_payment` | `/api/v1/scheduled-payments` | 28 |
| `scheduled_payment_share` | partes de `/api/v1/scheduled-payments` | 23 |
| `credit_statement` | `/api/v1/credits/{credit_id}/card-status` (lectura; motor S3) | 9 |
| `credit_payment` | `/api/v1/scheduled-payments/{scheduled_payment_id}/payments`, `/api/v1/credit-payments/{credit_payment_id}` | 22 |
| `savings_goal` | `/api/v1/savings-goals` | 14 |
| `savings_movement` | `/api/v1/savings-movements` | 11 |
| `meal_slot` | `/api/v1/meal-slots` | 8 |
| `recipe` | `/api/v1/recipes` | 17 |
| `recipe_ingredient` | ingredientes de `/api/v1/recipes/{recipe_id}` | 17 |
| `meal_plan_entry` | `/api/v1/meal-plan-entries`, `/api/v1/meal-plan/week` | 21 |
| `inventory_item` | `/api/v1/inventory-items` | 13 |
| `health_profile` | `/api/v1/health-profiles` | 8 |
| `weight_entry` | `/api/v1/weight-entries` | 7 |
| `nutrition_fact` | `/api/v1/products/{product_id}/nutrition-facts` | 11 |
| `notification` | `/api/v1/notifications` | 7 |
| `notification_recipient` | `/api/v1/notification-recipients` | 8 |
| `notification_subscription` | `/api/v1/notifications/{notification_id}/recipients` | 5 |
| `notification_delivery` | `/api/v1/notification-deliveries` | 6 |
| `job` | `/api/v1/jobs` | 11 |
| `backup` | `/api/v1/backups` | 4 |
| `generated_report` | `/api/v1/reports` | 7 |

## Supuestos

Ninguno cambia cifras canon, el esquema de S1 ni tokens. Lo que fija valores no decididos está en bloqueos y marcado PROVISIONAL. Decisiones de diseño propias de la tarea de C1 (no son preferencias de Axel):
1. Prefijo `/api/v1`, nombre de la cookie `fh_session`, códigos de error en inglés y `operationId`: identificadores del contrato.
2. «Colección» = lista de registros de una entidad (paginada con `Page`); los agregados y métricas son objetos acotados sin paginar (C1 §5).
3. El estado de bloqueo del inicio de sesión (N-011) se entrega en la respuesta 429 del propio login, en lugar de una consulta pública, para no revelar usuarios (C1-12 solo admite login y salud públicos).
4. `GET /openapi.yaml` exige sesión en el contrato; el esqueleto lo sirve sin sesión porque la autenticación aún no existe (pendiente de S6, WF-3).
5. El esqueleto escucha solo en `127.0.0.1:8080` por HTTP hasta que S6 agregue TLS con mkcert y se resuelva `criterios-bind`; no se publica en la LAN.

## Bloqueos

Nuevos:
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-filtro-familia.md`: valores y semántica del filtro «familia».
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-lista-compras.md`: dónde guardar marcas de comprado, artículos a mano y tienda por artículo (N-151, N-152, N-153, N-165 sin endpoint).
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-huecos-u1-s1.md`: 7 huecos o diferencias entre U1 y S1 (descripción del gasto, movimientos de inventario, borrador del OCR, presupuesto por mes, tiempos de receta, ingreso esperado del 15 y del 30, registro del chat).
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-parametros-tecnicos.md`: máximo por página, idempotencia, SameSite, contraseña mínima, tamaño de PROFECO.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-campos-calculados.md`: definición de campos calculados de los motores (IMC, precio por unidad base, membresía, rendimiento por peso, puntaje, estado de inventario).
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S6-version-go.md`: Go 1.27 no disponible sin actualizar todo el sistema.

Existentes que C1 referencia (no modificados): `criterios-login-bloqueo`, `U1-periodos-semana-quincena`, `criterios-umbrales-ausentes`, `criterios-frescura-peso`, `S1-frescura-peso`, `criterios-promedios`, `criterios-optimizador`, `S1-restricciones-optimizador`, `S1-reposicion-parte-externa`, `S1-meta-kcal`, `criterios-proyeccion-12m`, `criterios-credito-formulas`, `criterios-invariante-deuda`.

## Pendientes

1. Axel: responder los 6 bloqueos nuevos (y los existentes referenciados) en la puerta 1; con las respuestas sale C1 v1.1.
2. S1: migración 00002 según las respuestas (solicitud 2); después C1 v1.1 agrega las operaciones de la lista de compras (N-151, N-152, N-153, N-165) y quita las marcas PROVISIONAL resueltas. (corrección r2) Hecho en el esquema 1.1 de S1 (migración 00001) y en C1 regenerado; quedan las 3 preguntas de `C1-huecos-u1-s1` = `S1-huecos-negocio` y la de `C1-lista-compras`.
3. S6 en WF-3: implementar los 193 manejadores sobre `routes.go` (sesión por cookie con argon2id y límite de intentos, TLS con mkcert en finanzas.local:8443, base SQLite con los pragmas que pide S1, `audit_log` en la misma transacción, jobs, respaldos, adjuntos), exigir sesión en `GET /openapi.yaml` y una prueba por operación (S6-03).
4. Go 1.27 (bloqueo `S6-version-go`); hasta entonces `go.mod` declara `go 1.26`.
5. Si el contrato cambia, regenerar juntos `openapi.yaml`, `routes.go` y la matriz de C1 con el generador (solicitud 1). (corrección r1) Comando: `python3 -I backend/internal/api/c1gen/build.py`, seguido de `build.py --check` y `verify_c1.py`.

## Solicitudes fuera de alcance

1. **Integrador** · ruta propuesta `docs/contratos/tools/c1/` · copiar desde `/tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/c1/` los archivos `parse_schema.py`, `gen1_base.py`, `gen2_schemas.py`, `gen3_ops.py`, `gen4_build.py`, `gen5_go.py`, `gen6_c1.py`, `c1_plantilla.md`, `c1_inconsistencias.md` y `verify_c1.py` (orden de uso: `python3 parse_schema.py && python3 gen4_build.py && python3 gen5_go.py && python3 gen6_c1.py`, desde esa carpeta) · motivo: `openapi.yaml`, `routes.go` y la matriz de `C1-datos.md` salen del mismo generador; editarlos a mano por separado rompería S6-03 y C1-02. El scratchpad es efímero. `verify_c1.py` también va en el anexo A. (corrección r1) **Atendida** dentro del alcance de C1: los archivos están en `backend/internal/api/c1gen/`, no en `docs/contratos/tools/c1/`, porque esa ruta queda fuera del alcance de C1.
2. **S1** (`backend/db/migrations/00002_…`, `diccionario-datos.md`) · según las respuestas de Axel: tabla `idempotency_key` (`C1-parametros-tecnicos`); `shopping_list` y `shopping_list_item` (`C1-lista-compras`); `inventory_movement` con motivo; `expense.description` admite NULL o regla de defecto; borrador de correcciones del OCR en `attachment`; monto de presupuesto vigente desde una fecha; tabla receta-tiempo de comida; monto esperado por quincena (`C1-huecos-u1-s1`) · motivo: necesidades de U1 sin dónde guardarse. (corrección r2) Atendida por S1 en el esquema 1.1, salvo los puntos que esperan a Axel (borrador del OCR y registro del chat).
3. **U1** (`docs/design/U1-arquitectura.md`) · (a) alinear «familia» de §0.2 con N-040 según `C1-filtro-familia`; (b) §3.29, columna «Endpoint C1»: remitir a `C1-datos.md` §10; (c) N-011 se cubre con la respuesta 429 del login; (d) el script `u1_verify.py` no admite excepciones: reporta N-151, N-152, N-153, N-165 y `GET /openapi.yaml`, que son las excepciones justificadas de C1 §11.
4. **Criterios** (`docs/criterios/C1.md`) · C1-01: el comando literal usa módulos que no están en el Python del sistema; proponer como equivalente `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator <archivo>` y `uv run --no-project --with pyyaml python3 -I -c …`; C1-02 y C1-10: «N-xx» = necesidades N-xxx de U1 (como ya propuso U1); C1-06: los filtros se llaman `person_id`, `category_id`, `store_id` y `credit_id` (iguales al `nombre_codigo`), más `family` y `period`.
5. **Criterios** (`docs/criterios/S6.md`) · S6-01 espera Go 1.27; depende de `S6-version-go`.
6. **S7** · definir el registro del chat con el bot (N-299) y las claves de `notification.params_json` sin decidir (hora del resumen diario, anticipación de precios viejos), ya expuestas en `updateNotification`.

## Paquetes instalados

- `openapi-spec-validator` **0.9.0** con `uvx` (entorno aislado en la caché de uv; no toca el sistema).
- `@redocly/cli` **2.62.1** con `npx -y` (caché de npm del usuario).
- `pyyaml` y `jsonschema` con `uv run --with` (entorno efímero en la caché de uv) para `verify_c1.py`.
- Nada con pacman. Go 1.26.5 ya estaba instalado; no instalé Go 1.27 (bloqueo `S6-version-go`). El backend solo usa la biblioteca estándar (sin `go.sum`).

## Consumo aproximado

- Unas 80 acciones de herramienta, sin subagentes.
- Leído, unos 260 KB: protocolo y decisiones (16 KB), plan §1, §3.4, §3.5 y §3.7 (~30 KB), criterios C1, S6 y U1 (~38 KB), 15 bloqueos existentes (~22 KB), U1 §0, §1.1, §2.3, §2.5 y §3 completo (~85 KB), reporte de S1 (~20 KB) y extractos del de U1 (~12 KB), resumen del esquema por script (~15 KB), extractos del diccionario (~12 KB), tabla R-nn de S3 (~9 KB).
- Escrito en el repositorio, unos 720 KB (570 KB del YAML generado, 106 KB de C1, 28 KB de Go, 12 KB de bloqueos y este reporte); en el scratchpad, unos 95 KB del generador.
- Contexto de la sesión cerca de 200 mil tokens al cierre; el trabajo terminó completo, sin compresión de contexto.

## Anexo A · `verify_c1.py`

(corrección r1) El script vive en `backend/internal/api/c1gen/verify_c1.py`. Se ejecuta con `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py [--db BASE]`. La versión que estaba en este anexo (sha256 `9d9563c4244e108ba0ee495618e041b7bb363eb1b2c3a46b318eb9c8afff8af2`, 227 líneas, igual a `scratchpad/c1/verify_c1.py`) se trajo tal cual y luego se ajustó: rutas relativas a su carpeta, opciones `--db` y `--root`, consistencia de la matriz §10 (pares de U1 §3.29 y `x-necesidades`), C1-02 según el texto vigente del criterio y código de salida 1 si hay fallas. El código ya no se copia aquí, para que no haya dos versiones.
