# WF-3a · pruebas-contrato · Reporte

Agente pruebas-contrato (Sonnet) · 12-oct-2026 · repositorio en HEAD «[Puerta 2 bis] K53» (sin commit). Criterios: `docs/criterios/S6.md` (S6-03, S6-04). Convenciones: `docs/estado/wf3/CONVENCIONES-backend.md` §12.

## Resultado en una línea
**191 de 216 operaciones cumplen** (2xx declarada y todas sus 4xx válidas contra `openapi.yaml`). Fallan 25: **21 por un solo defecto** (la categoría creada por la API tiene `category.code` NULL y el motor de presupuesto lo lee como texto: responde 500), **2 sin implementar** (501) y **2 que aceptan un cuerpo vacío** que el contrato declara obligatorio. Con ese defecto salvado solo en la prueba (variable de diagnóstico) cumplen **212 de 216** (quedan las 4 propias). Paridad de rutas: **216 en el contrato, 216 en el servidor, 0 sobrantes, 0 faltantes**. Idempotency-Key a 24 h, página máxima de 200 y capturas atrasadas (K9): **cumplen**. `verify_c1.py --db` = TODO OK y `build.py --check` = 0 (23 archivos).

## HECHO
- **`TestContrato_Operaciones`** (S6-03): recorre las 216 rutas de `api.Routes` con las 10 áreas reales registradas (`auth` incluye `sistema`, `ajustes`, `avisos` con Telegram falso, `comidas`, `comprasa`, `comprasb`, `creditos`, `gastos-inicio`, `presupuesto`, `reportes`), sesión de Axel, servidor TLS de `testutil` y reloj simulado. Datos: todo lo que la API puede crear se crea por la API; solo lo que ninguna operación crea entra por SQL (3 cosas: grupo de producto, regla de reparto, chat detectado y destinatario de Telegram). Cada respuesta (2xx y 4xx) pasa por `testutil.ValidateResponse` (kin-openapi) en la subprueba de su operación. Por operación comprueba, según aplique: 401 sin credencial, 400 sin `Idempotency-Key`, la 2xx declarada, repetición con la misma clave (`Idempotent-Replayed: true`, mismo cuerpo), 409 `idempotency_key_reused` con otro cuerpo, 404 con id inexistente, 400 con campo desconocido, 422 con cuerpo vacío si el esquema tiene obligatorios, y `limit` 200 sí / 201 y 0 = 400. Total: **864 llamadas validadas** en esta prueba.
- **`TestContrato_Rutas`** (S6-03): contrato ↔ `routes.go` ↔ metadatos; operaciones sin implementar; por HTTP, ninguna ruta del contrato da 405 y todo método fuera del contrato da 405 `method_not_allowed` (la excepción `GET /credit-purchases/preview`, que cae en `/{credit_purchase_id}`, se reconoce); ruta desconocida = 404 de C1.
- **`TestContrato_Idempotencia24h`** (S6-04, K6): alta de ingreso; a 1 h y a 23 h 59 min la misma clave repite la respuesta guardada (mismo id, 1 solo renglón); a 24 h 1 min ya no la repite y ejecuta de nuevo; el job de purga (`dominio.PurgeIdempotency`) borra las vencidas (quedan 0).
- **`TestContrato_CapturasAtrasadas`** (S6-04, K9): 3 gastos y 2 pesos capturados con 48 h de atraso se aceptan; reenviados a las 25 h y a las 49 h (con purga entre medias) responden 200 + `Idempotent-Replayed`, mismo id, y la cuenta de renglones no cambia (3 y 2).
- **`TestContrato_PaginaMaxima`** (S6-04): 205 ingresos; `limit=200` = 200 renglones + `next_cursor`; la segunda página trae 5 y cursor nulo; 205 distintos; `limit=201`, `0`, `-1`, `abc` = 400; por omisión 50.
- Salida: variable `CONTRATO_TABLA=<ruta absoluta>` escribe la tabla operación → prueba → resultado; `CONTRATO_CATEGORIA_CON_CODE=1` es **diagnóstico**: pone un `code` a cada categoría creada para ver qué más falla detrás del defecto 1 (no oculta nada en la corrida normal).

## Archivos
Nuevos (todos en `backend/internal/api/`): `contrato_test.go` (runner, generador de cuerpos desde el esquema, rutas), `contrato_datos_test.go` (cuerpos de negocio, filas de apoyo, `logout` y `changeOwnPassword` con sesión propia), `contrato_tiempo_test.go` (24 h, K9, página máxima). Reporte: `docs/estado/wf3/pruebas-contrato.md`. No se tocó código de producción ni de otras áreas; `api_test.go` de backend-base queda igual (sus pruebas pasan).

## Verificación
- **S6-03 (a) paridad**: `go test ./internal/api -run TestContrato_Rutas -v` → «contrato 216 operaciones, servidor 216, sobrantes 0, faltantes 0, sin implementar 2». La prueba **falla** por las 2 sin implementar (`getDashboardBestBasket`, `getDashboardWeekMeals`).
- **S6-03 (b) validación de respuestas**: `go test ./internal/api -run TestContrato_Operaciones -v` → «S6-03: 191 de 216 operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)». Con `CONTRATO_CATEGORIA_CON_CODE=1`: «212 de 216».
- **S6-03 (c)**: `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/contrato.db` (base migrada con `finanzas migrate`, versión 4) → exit 0, 33 comprobaciones «OK», «RESULTADO: TODO OK». `python3 -I backend/internal/api/c1gen/build.py --check` → exit 0, «los 23 archivos coinciden byte a byte con el generador».
- **S6-04 Idempotency-Key 24 h**: `TestContrato_Idempotencia24h` PASA («clave retenida 23 h 59 min; a las 24 h 1 min ya no; la purga borró 2 vencidas y quedan 0»). Además, cada una de las 44 operaciones idempotentes pasó por 400 sin clave, repetición con `Idempotent-Replayed` y 409 con otro cuerpo (las 3 de agente externo y `login` quedan fuera de la repetición por su autenticación; ver tabla).
- **S6-04 capturas atrasadas (K9)**: `TestContrato_CapturasAtrasadas` PASA («3 gastos y 2 pesos … reenviados a las 25 h y a las 49 h (con purga) no se duplican»).
- **S6-04 página máxima**: `TestContrato_PaginaMaxima` PASA («205 ingresos = 200 + 5, 201/0/-1/abc = 400, por omisión 50») y la prueba por operación comprobó `limit` 200/201/0 en todas las colecciones.
- `go build ./...` y `go vet ./internal/api` limpios; `gofmt -l` vacío; `go test ./... -count=1`: todos los paquetes ok **excepto `internal/api`**, que queda en rojo por los hallazgos de abajo (a propósito: no se corrigen aquí). **`make verify-wf3` queda en rojo hasta que los dueños corrijan.** Segunda corrida igual a la primera (sin intermitencias).

## Operaciones que fallan (causa y dueño)
| Operaciones | Causa | Dueño |
|---|---|---|
| 21: `getBudgetMonth`, `getBudgetByCategory`, `getBudgetFortnights`, `getBudgetAlerts`, `getBudgetMonthlyTrend`, `getAccruedGroceries`, `getAccrualVsCash`, `getCashFlowAssumptions`, `getCashFlowFortnights`, `getLiveDebt`, `getSavingsSummary`, `getSavingsMonthlyTrend`, `listSavingsGoals`, `getExpenseMonthlyTrend`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow`, `getDashboardSavings`, `getPersonParticipation`, `getMonthlyReport`, `generateMonthlyReport` y, por arrastre, `sendReport` | 500 `internal_error`: «sql: Scan error on column index 1, name "code": converting NULL to string is unsupported». `backend/internal/dominio/presupuesto_motor.go:226-231` lee `category.code` en un `string`, pero el esquema permite NULL («NULL en las del usuario», migración 00001) y `createCategory` no pone código. Basta **una** categoría creada por la API (o por el CSV de partidas) para que todo el presupuesto, el ahorro, los dashboards y los reportes respondan 500. | api-presupuesto |
| 2: `getDashboardBestBasket`, `getDashboardWeekMeals` | Sin `reg.Handle` en ninguna área: responden 501 `not_implemented`. | api-gastos-inicio (inicio.yaml, CONVENCIONES §14) |
| `saveTicketScanDraft` | Con `{}` responde 200 y guarda `{"lines":[]}`; `TicketScanDraft` declara `lines` obligatorio (debería ser 422). | api-gastos-inicio |
| `setStoreProductGroups` | Con `{}` responde 204 y vacía los grupos; `IdList` declara `ids` obligatorio (debería ser 422). | api-compras-a |

## Excepciones documentadas (cuentan como cumplidas; se valida el código alterno contra el contrato)
`suggestRecipeWithAi` → 409 `ai_disabled` (IA apagada y sin proveedor real, S6-13); `reloadNutritionBase` → 409 (falta el archivo de la base nutricional; bloqueo `S6-compras-b-nutricion`). Su 2xx no es alcanzable en el entorno de pruebas.

## Tabla operación → prueba → resultado
Prueba de cada fila: `go test ./internal/api -run 'TestContrato_Operaciones/<operationId>' -v` (la subprueba lleva el operationId). «Llamadas» = peticiones validadas contra el contrato. Claves de la última columna: `401` sin sesión; `400-clave` sin Idempotency-Key; `200/201/202/204` la 2xx declarada; `repite` misma clave = respuesta guardada; `409-otro` misma clave con otro cuerpo; `404` id inexistente; `400-campo` campo desconocido; `422-vacío` cuerpo vacío con obligatorios; `límite` limit 200 sí, 201 y 0 = 400; `409*` excepción documentada. Columna «Con código en la categoría» = misma prueba con `CONTRATO_CATEGORIA_CON_CODE=1`.

| # | Operación | Método y ruta | Llamadas | Resultado | Con código en la categoría | Qué se comprobó / causa |
|---|---|---|---|---|---|---|
| 1 | `createCategory` | POST `/categories` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 2 | `createStore` | POST `/stores` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 3 | `createProduct` | POST `/products` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 4 | `createPresentation` | POST `/products/{product_id}/presentations` | 8 | OK | OK | 401 400-clave 201 repite 409-otro 404 400-campo 422-vacío |
| 5 | `createNutritionFact` | POST `/products/{product_id}/nutrition-facts` | 8 | OK | OK | 401 400-clave 201 repite 409-otro 404 400-campo 422-vacío |
| 6 | `createPrice` | POST `/prices` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 400-vacío |
| 7 | `createRecipe` | POST `/recipes` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 8 | `createInventoryItem` | POST `/inventory-items` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 9 | `createCredit` | POST `/credits` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 10 | `createCreditTermRate` | POST `/credits/{credit_id}/term-rates` | 8 | OK | OK | 401 400-clave 201 repite 409-otro 404 400-campo 422-vacío |
| 11 | `createPromotion` | POST `/promotions` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 12 | `createCreditPurchase` | POST `/credit-purchases` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 400-vacío |
| 13 | `createHealthProfile` | POST `/health-profiles` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 14 | `createWeightEntry` | POST `/weight-entries` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 400-vacío |
| 15 | `createExpense` | POST `/expenses` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 400-vacío |
| 16 | `uploadExpenseAttachment` | POST `/expenses/{expense_id}/attachments` | 5 | OK | OK | 401 400-clave 201 repite 404 |
| 17 | `createTicketScan` | POST `/ticket-scans` | 4 | OK | OK | 401 400-clave 202 repite |
| 18 | `saveTicketScanDraft` | PUT `/ticket-scans/{attachment_id}/draft` | 5 | **FALLA** | **FALLA** | cuerpo vacío con campos obligatorios: 200, se esperaba 422: {"id":2,"kind":"ticket_photo","original_name":"ticket.png","mime_type":"image/png","size_bytes":2048,"sha256":"c99ee08818588a5a4d3aac3a2c8553c971bd25d227db2cc61d33fc097714151b","oc |
| 19 | `createIncome` | POST `/incomes` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 400-vacío |
| 20 | `createBudget` | POST `/budgets` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 21 | `createSavingsGoal` | POST `/savings-goals` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 22 | `createSavingsMovement` | POST `/savings-movements` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 400-vacío |
| 23 | `createPlannedItem` | POST `/planned-items` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 24 | `createShoppingListItem` | POST `/shopping-list/items` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 25 | `createMealPlanEntries` | POST `/meal-plan-entries` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 400-campo 422-vacío |
| 26 | `createBackup` | POST `/backups` | 4 | OK | OK | 401 400-clave 202 repite |
| 27 | `generateMonthlyReport` | POST `/reports/monthly` | 5 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [202]: {…} / LOG: time=2026-10-12T04:14:43.264-06:00 level=ERROR msg="error interno" operation=generateMonthlyReport request_id=7987c4beb46e14141099700c error="sql: Scan error on column index 1, name "co |
| 28 | `generateExternalStatement` | POST `/reports/external-statements` | 7 | OK | OK | 401 400-clave 202 repite 409-otro 400-campo 422-vacío |
| 29 | `listAiProviders` | GET `/ai/providers` | 2 | OK | OK | 401 200 |
| 30 | `downloadAttachment` | GET `/attachments/{attachment_id}/file` | 3 | OK | OK | 401 200 404 |
| 31 | `getCurrentSession` | GET `/auth/me` | 2 | OK | OK | 401 200 |
| 32 | `listBackups` | GET `/backups` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 33 | `getBasketComparison` | GET `/basket-comparison` | 2 | OK | OK | 401 200 |
| 34 | `listBaskets` | GET `/baskets` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 35 | `getAccruedGroceries` | GET `/budget/accrued-groceries` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.544-06:00 level=ERROR msg="error interno" operation=getAccruedGroceries request_id=a7edd8e58fc1b7243a9b63ae error="sql: Scan error on column index 1, name "code |
| 36 | `getBudgetAlerts` | GET `/budget/alerts` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.545-06:00 level=ERROR msg="error interno" operation=getBudgetAlerts request_id=e5979ab79ab9a1e33923161e error="sql: Scan error on column index 1, name "code": c |
| 37 | `getBudgetByCategory` | GET `/budget/categories` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.546-06:00 level=ERROR msg="error interno" operation=getBudgetByCategory request_id=cabf813d34885e2efeb3696e error="sql: Scan error on column index 1, name "code |
| 38 | `getBudgetFortnights` | GET `/budget/fortnights` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.548-06:00 level=ERROR msg="error interno" operation=getBudgetFortnights request_id=e89c3e1609cc29423738e2a3 error="sql: Scan error on column index 1, name "code |
| 39 | `getBudgetMonth` | GET `/budget/month` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.549-06:00 level=ERROR msg="error interno" operation=getBudgetMonth request_id=b3a1f357e886b8722ef75471 error="sql: Scan error on column index 1, name "code": co |
| 40 | `getBudgetMonthlyTrend` | GET `/budget/monthly-trend` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.550-06:00 level=ERROR msg="error interno" operation=getBudgetMonthlyTrend request_id=acc2350944a7be1b3cc8d3cd error="sql: Scan error on column index 1, name "co |
| 41 | `listBudgets` | GET `/budgets` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 42 | `listCalendarEvents` | GET `/calendar/events` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 43 | `getAccrualVsCash` | GET `/cash-flow/accrual-vs-cash` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.558-06:00 level=ERROR msg="error interno" operation=getAccrualVsCash request_id=459d8651d2d7d1aba25b04c4 error="sql: Scan error on column index 1, name "code":  |
| 44 | `getCashFlowAssumptions` | GET `/cash-flow/assumptions` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.559-06:00 level=ERROR msg="error interno" operation=getCashFlowAssumptions request_id=93934c22ac23874b112836ca error="sql: Scan error on column index 1, name "c |
| 45 | `getCashFlowFortnights` | GET `/cash-flow/fortnights` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.560-06:00 level=ERROR msg="error interno" operation=getCashFlowFortnights request_id=978c6e50ecfe458454ed222b error="sql: Scan error on column index 1, name "co |
| 46 | `listCategories` | GET `/categories` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 47 | `listCostShareRules` | GET `/cost-share-rules` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 48 | `getCreditPurchase` | GET `/credit-purchases/{credit_purchase_id}` | 3 | OK | OK | 401 200 404 |
| 49 | `listCreditPurchases` | GET `/credit-purchases` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 50 | `getDebtTrend` | GET `/credits/debt-trend` | 2 | OK | OK | 401 200 |
| 51 | `getCreditTotals` | GET `/credits/totals` | 2 | OK | OK | 401 200 |
| 52 | `previewAvailability` | GET `/credits/{credit_id}/availability-preview` | 3 | OK | OK | 401 200 404 |
| 53 | `getAvailabilitySeries` | GET `/credits/{credit_id}/availability-series` | 3 | OK | OK | 401 200 404 |
| 54 | `getCardStatus` | GET `/credits/{credit_id}/card-status` | 4 | OK | OK | 401 200 404 |
| 55 | `listCreditTermRates` | GET `/credits/{credit_id}/term-rates` | 6 | OK | OK | 401 200 404 límite 201 y 0 = 400 |
| 56 | `getCredit` | GET `/credits/{credit_id}` | 3 | OK | OK | 401 200 404 |
| 57 | `listCredits` | GET `/credits` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 58 | `getDashboardBestBasket` | GET `/dashboard/best-basket` | 0 | **FALLA** | **FALLA** | no está registrada en ninguna área (responde 501) |
| 59 | `getDashboardFortnightFlow` | GET `/dashboard/fortnight-flow` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.602-06:00 level=ERROR msg="error interno" operation=getDashboardFortnightFlow request_id=37326699013f1d3bd96054ca error="sql: Scan error on column index 1, name |
| 60 | `getDashboardMonthVsPlan` | GET `/dashboard/month-vs-plan` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.603-06:00 level=ERROR msg="error interno" operation=getDashboardMonthVsPlan request_id=150526f009725ecb366c34a4 error="sql: Scan error on column index 1, name " |
| 61 | `getDashboardNextPayment` | GET `/dashboard/next-payment` | 2 | OK | OK | 401 200 |
| 62 | `getDashboardSavings` | GET `/dashboard/savings` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.607-06:00 level=ERROR msg="error interno" operation=getDashboardSavings request_id=4688d52cea8fd5bda9ba23cb error="sql: Scan error on column index 1, name "code |
| 63 | `getDashboardWeekMeals` | GET `/dashboard/week-meals` | 0 | **FALLA** | **FALLA** | no está registrada en ninguna área (responde 501) |
| 64 | `getLiveDebt` | GET `/debt/live` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.608-06:00 level=ERROR msg="error interno" operation=getLiveDebt request_id=599f9a94116b7d3ba3245a40 error="sql: Scan error on column index 1, name "code": conve |
| 65 | `getExpenseMonthlyTrend` | GET `/expenses/monthly-trend` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.609-06:00 level=ERROR msg="error interno" operation=getExpenseMonthlyTrend request_id=405d104c7cf04dab8f786a65 error="sql: Scan error on column index 1, name "c |
| 66 | `getExpenseSummary` | GET `/expenses/summary` | 2 | OK | OK | 401 200 |
| 67 | `getExpenseTotals` | GET `/expenses/totals` | 2 | OK | OK | 401 200 |
| 68 | `listExpenseItems` | GET `/expenses/{expense_id}/items` | 6 | OK | OK | 401 200 404 límite 201 y 0 = 400 |
| 69 | `getExpense` | GET `/expenses/{expense_id}` | 3 | OK | OK | 401 200 404 |
| 70 | `listExpenses` | GET `/expenses` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 71 | `getExternalAgentReference` | GET `/external-agent/reference` | 2 | OK | OK | 401 200 |
| 72 | `listHealthProfiles` | GET `/health-profiles` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 73 | `getTodayKcal` | GET `/health/today-kcal` | 2 | OK | OK | 401 200 |
| 74 | `getIncomeMonth` | GET `/incomes/month` | 2 | OK | OK | 401 200 |
| 75 | `getIncomeMonthlyTrend` | GET `/incomes/monthly-trend` | 2 | OK | OK | 401 200 |
| 76 | `listIncomes` | GET `/incomes` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 77 | `listInventoryMovements` | GET `/inventory-items/{inventory_item_id}/movements` | 6 | OK | OK | 401 200 404 límite 201 y 0 = 400 |
| 78 | `listInventoryItems` | GET `/inventory-items` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 79 | `getInventorySummary` | GET `/inventory/summary` | 2 | OK | OK | 401 200 |
| 80 | `listJobs` | GET `/jobs` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 81 | `getMealPlanMonth` | GET `/meal-plan/month` | 2 | OK | OK | 401 200 |
| 82 | `getMealPlanTrends` | GET `/meal-plan/trends` | 2 | OK | OK | 401 200 |
| 83 | `getMealPlanWeek` | GET `/meal-plan/week` | 2 | OK | OK | 401 200 |
| 84 | `listMealSlots` | GET `/meal-slots` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 85 | `listMonthlyReportMonths` | GET `/monthly-report/months` | 2 | OK | OK | 401 200 |
| 86 | `getMonthlyReport` | GET `/monthly-report` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.656-06:00 level=ERROR msg="error interno" operation=getMonthlyReport request_id=453963efeeeb7a0de3bf880c error="sql: Scan error on column index 1, name "code":  |
| 87 | `listNotificationDeliveries` | GET `/notification-deliveries` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 88 | `listNotificationRecipients` | GET `/notification-recipients` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 89 | `listNotifications` | GET `/notifications` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 90 | `listPeriods` | GET `/periods` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 91 | `getPersonParticipation` | GET `/persons/{person_id}/participation` | 3 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.686-06:00 level=ERROR msg="error interno" operation=getPersonParticipation request_id=48b1687d1840a477542a8ca1 error="sql: Scan error on column index 1, name "c |
| 92 | `getExternalStatement` | GET `/persons/{person_id}/statement` | 3 | OK | OK | 401 200 404 |
| 93 | `listPersons` | GET `/persons` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 94 | `listPlannedItems` | GET `/planned-items` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 95 | `getPriceSources` | GET `/price-sources` | 2 | OK | OK | 401 200 |
| 96 | `getPriceSummary` | GET `/prices/summary` | 2 | OK | OK | 401 200 |
| 97 | `listProductGroups` | GET `/product-groups` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 98 | `listNutritionFacts` | GET `/products/{product_id}/nutrition-facts` | 6 | OK | OK | 401 200 404 límite 201 y 0 = 400 |
| 99 | `listProductPrices` | GET `/products/{product_id}/prices` | 6 | OK | OK | 401 200 404 límite 201 y 0 = 400 |
| 100 | `getProduct` | GET `/products/{product_id}` | 3 | OK | OK | 401 200 404 |
| 101 | `listProducts` | GET `/products` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 102 | `listPromotions` | GET `/promotions` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 103 | `getRecipeUsage` | GET `/recipes/{recipe_id}/usage` | 3 | OK | OK | 401 200 404 |
| 104 | `getRecipe` | GET `/recipes/{recipe_id}` | 3 | OK | OK | 401 200 404 |
| 105 | `listRecipes` | GET `/recipes` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 106 | `downloadReportFile` | GET `/reports/{report_id}/file` | 3 | OK | OK | 401 200 404 |
| 107 | `listReports` | GET `/reports` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 108 | `listSavingsGoals` | GET `/savings-goals` | 5 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.729-06:00 level=ERROR msg="error interno" operation=listSavingsGoals request_id=2d19a454593146568b5ba2ad error="sql: Scan error on column index 1, name "code":  |
| 109 | `listSavingsMovements` | GET `/savings-movements` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 110 | `getSavingsMonthlyTrend` | GET `/savings/monthly-trend` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.736-06:00 level=ERROR msg="error interno" operation=getSavingsMonthlyTrend request_id=bd3825abb2c0cb8b66db326e error="sql: Scan error on column index 1, name "c |
| 111 | `getSavingsSummary` | GET `/savings/summary` | 2 | **FALLA** | OK | 2xx: llegó 500 y el contrato declara [200]: {…} / LOG: time=2026-10-12T04:14:43.737-06:00 level=ERROR msg="error interno" operation=getSavingsSummary request_id=3d652a9e9813fffc7159361e error="sql: Scan error on column index 1, name "code": |
| 112 | `listScheduledPayments` | GET `/scheduled-payments` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 113 | `search` | GET `/search` | 2 | OK | OK | 401 200 |
| 114 | `listSessions` | GET `/sessions` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 115 | `getAlertSettings` | GET `/settings/alerts` | 2 | OK | OK | 401 200 |
| 116 | `getDataSettings` | GET `/settings/data` | 2 | OK | OK | 401 200 |
| 117 | `getHealthSettings` | GET `/settings/health` | 2 | OK | OK | 401 200 |
| 118 | `getOptimizerSettings` | GET `/settings/optimizer` | 2 | OK | OK | 401 200 |
| 119 | `getPriceSettings` | GET `/settings/prices` | 2 | OK | OK | 401 200 |
| 120 | `getProjectionSettings` | GET `/settings/projection` | 2 | OK | OK | 401 200 |
| 121 | `getRecipeSettings` | GET `/settings/recipes` | 2 | OK | OK | 401 200 |
| 122 | `getReportSettings` | GET `/settings/reports` | 2 | OK | OK | 401 200 |
| 123 | `getSecuritySettings` | GET `/settings/security` | 2 | OK | OK | 401 200 |
| 124 | `getShoppingList` | GET `/shopping-list` | 2 | OK | OK | 401 200 |
| 125 | `listStores` | GET `/stores` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 126 | `getCaCertificate` | GET `/system/ca-certificate` | 2 | OK | OK | 401 200 |
| 127 | `getSeedStatus` | GET `/system/seed-status` | 2 | OK | OK | 401 200 |
| 128 | `getSystemStatus` | GET `/system/status` | 2 | OK | OK | 401 200 |
| 129 | `listTelegramChatCandidates` | GET `/telegram-chat-candidates` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 130 | `getTicketScan` | GET `/ticket-scans/{attachment_id}` | 3 | OK | OK | 401 200 404 |
| 131 | `listTicketScans` | GET `/ticket-scans` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 132 | `listUsers` | GET `/users` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 133 | `queryView` | GET `/views/{entity}` | 6 | OK | OK | 401 200 404 límite 201 y 0 = 400 |
| 134 | `listWeightEntries` | GET `/weight-entries` | 5 | OK | OK | 401 200 límite 201 y 0 = 400 |
| 135 | `getHealth` | GET `/health` | 1 | OK | OK | 200 |
| 136 | `getOpenApiDocument` | GET `/openapi.yaml` | 2 | OK | OK | 401 200 |
| 137 | `updateAiProvider` | PATCH `/ai/providers/{provider}` | 5 | OK | OK | 401 200 404 400-campo 422-vacío |
| 138 | `login` | POST `/auth/login` | 1 | OK | OK | 200 |
| 139 | `updateBudget` | PATCH `/budgets/{budget_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 140 | `reorderCategories` | PUT `/categories/order` | 4 | OK | OK | 401 204 400-campo 422-vacío |
| 141 | `updateCategory` | PATCH `/categories/{category_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 142 | `updateCostShareRule` | PATCH `/cost-share-rules/{cost_share_rule_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 143 | `previewCreditPurchase` | POST `/credit-purchases/preview` | 4 | OK | OK | 401 200 400-campo 400-vacío |
| 144 | `updateCreditPurchase` | PATCH `/credit-purchases/{credit_purchase_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 145 | `updateCreditTermRate` | PATCH `/credit-term-rates/{credit_term_rate_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 146 | `updateCredit` | PATCH `/credits/{credit_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 147 | `duplicateExpense` | POST `/expenses/{expense_id}/duplicate` | 5 | OK | OK | 401 400-clave 201 repite 404 |
| 148 | `markExpensePaid` | POST `/expenses/{expense_id}/mark-paid` | 9 | OK | OK | 401 400-clave 200 repite 409-otro 404 400-campo 400-vacío |
| 149 | `updateExpense` | PATCH `/expenses/{expense_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 150 | `importExternalCatalog` | POST `/external-agent/catalog` | 5 | OK | OK | 401 400-clave 200 400-campo 422-vacío |
| 151 | `importExternalPrices` | POST `/external-agent/prices` | 5 | OK | OK | 401 400-clave 200 400-campo 422-vacío |
| 152 | `importExternalRecipes` | POST `/external-agent/recipes` | 5 | OK | OK | 401 400-clave 200 400-campo 422-vacío |
| 153 | `updateHealthProfile` | PATCH `/health-profiles/{health_profile_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 154 | `previewPlannedItemsImport` | POST `/imports/planned-items/preview` | 2 | OK | OK | 401 200 |
| 155 | `importPlannedItems` | POST `/imports/planned-items` | 4 | OK | OK | 401 400-clave 201 repite |
| 156 | `updateIncome` | PATCH `/incomes/{income_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 157 | `adjustInventoryItem` | POST `/inventory-items/{inventory_item_id}/adjustments` | 8 | OK | OK | 401 400-clave 200 repite 409-otro 404 400-campo 422-vacío |
| 158 | `updateInventoryItem` | PATCH `/inventory-items/{inventory_item_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 159 | `updateJob` | PATCH `/jobs/{job_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 160 | `updateMealPlanEntry` | PATCH `/meal-plan-entries/{meal_plan_entry_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 161 | `copyPreviousWeek` | POST `/meal-plan/copy-previous-week` | 4 | OK | OK | 401 400-clave 201 repite |
| 162 | `updateNotificationRecipient` | PATCH `/notification-recipients/{recipient_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 163 | `setNotificationRecipients` | PUT `/notifications/{notification_id}/recipients` | 5 | OK | OK | 401 204 404 400-campo 422-vacío |
| 164 | `sendTestNotification` | POST `/notifications/{notification_id}/test` | 5 | OK | OK | 401 202 404 400-campo 422-vacío |
| 165 | `updateNotification` | PATCH `/notifications/{notification_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 166 | `reloadNutritionBase` | POST `/nutrition-facts/reload` | 5 | OK | OK | 401 400-clave 409* 400-campo 422-vacío |
| 167 | `updatePerson` | PATCH `/persons/{person_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 168 | `updatePlannedItem` | PATCH `/planned-items/{planned_item_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 169 | `updatePresentation` | PATCH `/presentations/{presentation_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 170 | `updatePrice` | PATCH `/prices/{price_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 171 | `updateProduct` | PATCH `/products/{product_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 172 | `importProfeco` | POST `/profeco-imports` | 4 | OK | OK | 401 400-clave 201 repite |
| 173 | `updatePromotion` | PATCH `/promotions/{promotion_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 174 | `suggestRecipeWithAi` | POST `/recipes/ai-suggestions` | 4 | OK | OK | 401 409* 400-campo 422-vacío |
| 175 | `duplicateRecipe` | POST `/recipes/{recipe_id}/duplicate` | 5 | OK | OK | 401 400-clave 201 repite 404 |
| 176 | `updateRecipe` | PATCH `/recipes/{recipe_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 177 | `sendReport` | POST `/reports/{report_id}/send` | 4 | **FALLA** | OK | 2xx: llegó 409 y el contrato declara [202]: {…} |
| 178 | `updateSavingsGoal` | PATCH `/savings-goals/{savings_goal_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 179 | `updateSavingsMovement` | PATCH `/savings-movements/{savings_movement_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 180 | `compareScenarios` | POST `/scenario-comparisons` | 4 | OK | OK | 401 200 400-campo 400-vacío |
| 181 | `markScheduledPaymentPaid` | POST `/scheduled-payments/{scheduled_payment_id}/payments` | 8 | OK | OK | 401 400-clave 201 repite 409-otro 404 400-campo 400-vacío |
| 182 | `updateAlertSettings` | PATCH `/settings/alerts` | 3 | OK | OK | 401 200 400-campo |
| 183 | `updateDataSettings` | PATCH `/settings/data` | 3 | OK | OK | 401 200 400-campo |
| 184 | `updateHealthSettings` | PATCH `/settings/health` | 3 | OK | OK | 401 200 400-campo |
| 185 | `updateOptimizerSettings` | PATCH `/settings/optimizer` | 3 | OK | OK | 401 200 400-campo |
| 186 | `updatePriceSettings` | PATCH `/settings/prices` | 3 | OK | OK | 401 200 400-campo |
| 187 | `updateProjectionSettings` | PATCH `/settings/projection` | 4 | OK | OK | 401 200 400-campo 422-vacío |
| 188 | `updateReportSettings` | PATCH `/settings/reports` | 3 | OK | OK | 401 200 400-campo |
| 189 | `updateSecuritySettings` | PATCH `/settings/security` | 3 | OK | OK | 401 200 400-campo |
| 190 | `assignShoppingListStore` | POST `/shopping-list/assign-store` | 7 | OK | OK | 401 400-clave 200 repite 409-otro 400-campo 400-vacío |
| 191 | `generateShoppingList` | POST `/shopping-list/generate` | 4 | OK | OK | 401 400-clave 200 repite |
| 192 | `updateShoppingListItem` | PATCH `/shopping-list/items/{shopping_list_item_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 193 | `setStoreProductGroups` | PUT `/stores/{store_id}/product-groups` | 5 | **FALLA** | **FALLA** | cuerpo vacío con campos obligatorios: 204, se esperaba 422: |
| 194 | `updateStore` | PATCH `/stores/{store_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 195 | `confirmTelegramChatCandidate` | POST `/telegram-chat-candidates/{candidate_id}/confirm` | 7 | OK | OK | 401 400-clave 201 repite 409-otro 404 400-campo |
| 196 | `confirmTicketScan` | POST `/ticket-scans/{attachment_id}/confirm` | 8 | OK | OK | 401 400-clave 201 repite 409-otro 404 400-campo 400-vacío |
| 197 | `updateWeightEntry` | PATCH `/weight-entries/{weight_entry_id}` | 4 | OK | OK | 401 200 404 400-campo |
| 198 | `deleteAttachment` | DELETE `/attachments/{attachment_id}` | 3 | OK | OK | 401 204 404 |
| 199 | `deleteBudget` | DELETE `/budgets/{budget_id}` | 3 | OK | OK | 401 204 404 |
| 200 | `deleteCreditPayment` | DELETE `/credit-payments/{credit_payment_id}` | 3 | OK | OK | 401 204 404 |
| 201 | `deleteCreditPurchase` | DELETE `/credit-purchases/{credit_purchase_id}` | 3 | OK | OK | 401 204 404 |
| 202 | `deleteCreditTermRate` | DELETE `/credit-term-rates/{credit_term_rate_id}` | 3 | OK | OK | 401 204 404 |
| 203 | `deleteExpense` | DELETE `/expenses/{expense_id}` | 3 | OK | OK | 401 204 404 |
| 204 | `deleteIncome` | DELETE `/incomes/{income_id}` | 3 | OK | OK | 401 204 404 |
| 205 | `deleteMealPlanEntry` | DELETE `/meal-plan-entries/{meal_plan_entry_id}` | 3 | OK | OK | 401 204 404 |
| 206 | `deletePrice` | DELETE `/prices/{price_id}` | 3 | OK | OK | 401 204 404 |
| 207 | `deletePromotion` | DELETE `/promotions/{promotion_id}` | 3 | OK | OK | 401 204 404 |
| 208 | `deleteRecipe` | DELETE `/recipes/{recipe_id}` | 3 | OK | OK | 401 204 404 |
| 209 | `deleteSavingsMovement` | DELETE `/savings-movements/{savings_movement_id}` | 3 | OK | OK | 401 204 404 |
| 210 | `revokeSession` | DELETE `/sessions/{session_id}` | 3 | OK | OK | 401 204 404 |
| 211 | `dismissTelegramChatCandidate` | DELETE `/telegram-chat-candidates/{candidate_id}` | 3 | OK | OK | 401 204 404 |
| 212 | `deleteWeightEntry` | DELETE `/weight-entries/{weight_entry_id}` | 3 | OK | OK | 401 204 404 |
| 213 | `restoreBackup` | POST `/backups/{backup_id}/restore` | 5 | OK | OK | 401 202 404 400-campo 422-vacío |
| 214 | `revokeSessions` | POST `/sessions/revoke` | 4 | OK | OK | 401 200 400-campo 422-vacío |
| 215 | `changeOwnPassword` | POST `/auth/password` | 3 | OK | OK | 204 la contraseña nueva entra 422 con la contraseña actual incorrecta |
| 216 | `logout` | POST `/auth/logout` | 3 | OK | OK | 204 la sesión cerrada da 401 401 |

## Supuestos
1. Los fixtures `_test` de cada área no se pueden importar desde otro paquete (están en paquetes `*_test`); se reconstruyeron por la API y con SQL mínimo. Si se prefiere reutilizarlos, hace falta que alguien los exporte a `testutil` (PROPUESTA).
2. Se registran las áreas a mano en la prueba (no se usa `cmd/finanzas/areas.go`, que está vacío): `auth.Register` ya incluye `sistema.Register`. `avisos` se registra con `RegisterWith` y un Telegram falso (`httptest`); `reportes.SetSender` recibe un enviador falso; `Deps.Service` un verificador falso para el agente externo (K53). Nunca sale nada a la red ni se usa el token real.
3. Un 404 con id inexistente exige el cuerpo válido de la operación (si el cuerpo es inválido el 400/422 va primero; no se reporta).
4. Una operación «cumple» si la 2xx declarada se alcanza (o es excepción documentada) y ninguna de sus pruebas 4xx falla ni deja de ser válida contra `openapi.yaml`. «Cuerpo vacío = 422 (o 400)» solo se exige cuando el esquema declara campos obligatorios.
5. `getCaCertificate` se prueba con un archivo PEM falso en `sistema.CACertPath` (restaurado al terminar).

## Bloqueos
Ninguno nuevo. Relacionados: `S6-compras-b-nutricion` (reloadNutritionBase), `S5-agente-externo` (credencial de servicio PROVISIONAL).

## Pendientes
- Cuando se corrijan los 4 hallazgos: volver a correr `go test ./internal/api -run TestContrato -v`; debe dar «216 de 216».
- Segunda pasada sugerida: correr la misma prueba con la semilla real (`internal/semilla`) para validar respuestas con datos ricos (no se hizo: la semilla crea sus propias personas y cuentas y choca con las fixtures de `testutil.NewDB`).

## Solicitudes fuera de alcance (resumen; el JSON las trae completas)
1. **api-presupuesto** · `backend/internal/dominio/presupuesto_motor.go` líneas 226-231: leer `code` con `sql.NullString`; si no es válido, no registrar `fe.CategoryID[""]` (varias categorías sin código chocarían en la llave vacía) y mandar `Code` vacío al motor solo si el motor lo acepta (PREGUNTA al dueño del motor S2: cómo trata una categoría sin `code`).
2. **api-gastos-inicio** · registrar `getDashboardBestBasket` y `getDashboardWeekMeals`; en `saveTicketScanDraft` exigir `lines` (422).
3. **api-compras-a** · en `setStoreProductGroups` exigir `ids` (422).
4. **Integrador / backend-base** · `backend/cmd/finanzas/areas.go` sigue con `var areas = []func(*api.Registry){}`: **el servidor real respondería 501 en todo**. Líneas: `auth.Register`, `ajustes.Register`, `avisos.Register`, `comidas.Register`, `comprasa.Register`, `comprasb.Register`, `creditos.Register`, `gastos.Register` (import `gastos "finanzas-hogar/internal/api/gastos-inicio"`), `presupuesto.Register`, `reportesapi.Register` (import `reportesapi "finanzas-hogar/internal/api/reportes"`). Además `main.go` no llama a `reportes.SetSender(...)` ni instala `Deps.Service`: `sendReport` daría 503 y el agente externo 401 en producción (PROVISIONAL conocido).

## Propuestas para Axel
- Ninguna decisión nueva. Para el integrador: no hacer commit de WF-3a con `make verify-wf3` en rojo; los 4 hallazgos son pequeños.

## Paquetes instalados
Ninguno (usa `kin-openapi` y `testutil`, ya presentes). Sin npm, sin sudo.

## Consumo aproximado
~48 acciones; ~230 KB leídos; ~85 KB escritos (3 archivos de prueba ~52 KB, reporte ~33 KB). Temporales de `datos/tmp/` (base migrada, tablas, bitácoras) borrados.
