# api-creditos · WF-3a · reporte

Agente api-creditos (Opus) · 12-oct-2026 · alcance: `backend/internal/api/creditos/`, `backend/internal/dominio/creditos_*.go`. Sin `git add` ni commit.

## Hecho
- **HECHO** Las 29 operaciones de C1 v1.4 de Créditos (`creditos-a.yaml`, `creditos-b.yaml` y las 10 de `ajustes-d.yaml`) en `creditos.Register(reg)`, con CONVENCIONES-backend: errores de C1 en español de México, paginación de 1 a 200 con `ParsePage`, filtros (`credit_id`, `family` incluyente como C1-filtro-familia opción A, `period`, `person_id`, `store_id`, `kind`, `status`, `q`, `valid_on`, `merchant_name`), Idempotency-Key (la pone el middleware en createCredit, createCreditTermRate, createPromotion, createCreditPurchase y markScheduledPaymentPaid), JSON Merge Patch con `api.Nullable`, `audit_log` en cada alta, cambio y baja, centavos enteros.
- **HECHO** Cálculos por el motor S3 sin tocarlo: `QuotePurchase` (preview y alta de compras: intereses, partes, serie con el último ajustando centavos), `AvailabilityTimeline` (A/B, vista previa, serie, disponible), `LiveDebt` (deuda viva, tendencia, saldo de tarjeta), `UnifiedCalendar` (totales de quincena y mes, intereses), `Credit.CardMinimum` (pago mínimo cuando el corte no lo trae), `Rate.Display` (rate_pct, K12), `comparador.Compare` con `comparator.timeout_seconds` (`TimeoutFromSetting` + contexto: vencido → 503 `unavailable`) y `optimizer.date_window_days` de app_setting; restricciones globales `optimizer.*` sobreescritas por la petición (K16); topes por persona de `person.fortnight_payment_cap_cents` (K13).
- **HECHO** Q2-67 opción 1: `vat_on_interest_bp` con valor y `vat_included` true, null o ausente (sobre lo guardado en un PATCH) → 422 `validation_failed` en `vat_on_interest_bp`; la API no cambia el otro campo. `missing_terms` de tarjeta y tienda con los 5 datos de K18 (`vat` mientras `vat_included` sea NULL).
- **HECHO** Repositorios en `dominio`: renglones de credit, credit_term_rate, promotion, credit_purchase(+share), scheduled_payment(+share con lo pagado), credit_payment, credit_statement y person (clave, participación, tope).
- **PREGUNTA / PROVISIONAL** 10 puntos en `docs/estado/bloqueos/S6-creditos.md` (gasto vinculado en un paso, ventana de la tendencia, saldo de tarjeta, partes y canastas del comparador, horizonte A/B, promoción sin comercio en la captura, recotización de compras guardadas, paid_off automático, compras canceladas en el calendario, códigos 409 de pagos).

## Archivos
- `backend/internal/dominio/creditos_repo.go`, `creditos_compras.go`, `creditos_repo_test.go`
- `backend/internal/api/creditos/creditos.go` (Register, líneas, totales, tendencia, A/B, tarjeta), `motor.go` (base → entrada del motor, errores del motor → C1), `condiciones.go` (alta y cambio de créditos, tasas, promociones, valores del comparador, comparador), `compras.go` (compras a crédito), `pagos.go` (pagos programados y abonos), `creditos_test.go`
- `docs/estado/bloqueos/S6-creditos.md`, este reporte.

## Verificación
- `go build ./...` y `go vet ./...` sin errores; `gofmt -l .` vacío (dentro de `make verify-wf3-go`).
- `go test ./internal/api/creditos/... ./internal/dominio/... -count=1 -cover` → `creditos` ok 81.8 %, `dominio` ok 70.3 % (el paquete completo con los repositorios de las demás áreas; los de créditos promedian 94 % por función con `TestCreditosRepos`).
- `make verify-wf3-go GOOSE=…` → build, vet, gofmt, `go mod tidy -diff` sin cambios, `go test ./...`: 19 paquetes OK, 8 sin pruebas.
- **S6-03**: 29 pruebas `TestOp_<operationId>`, todas PASS; cada una llama con `testutil.Client.Do`, que valida código y esquema contra `openapi.yaml`, y cubre 2xx y al menos un 4xx (tabla abajo).
- **S6-04**: errores con el esquema único (`validation_failed`, `not_found`, `already_exists`, `in_use`, `has_payments`, `not_revolving`, `not_a_card`, `bad_request`, `unauthorized`, `unavailable`); `limit=201` → 400; dinero no entero (`97.5`) → 400; fechas AAAA-MM-DD (`valid_on=10-oct-2026` → 400); POST sin Idempotency-Key → 400.
- **S6-15**: identificadores en inglés; mensajes en español de México; no se registran contraseñas, cookies ni tokens (los errores 500 van solo al registro con su request_id).
- Canon S3 reproducido en las pruebas: despensa 1 = 309450 + 80050 = 389500 en 4 pagos de 97375 desde 2026-10-30 (último capital 77361, 2026-12-15); tasa 4 quincenas `rate_pct` 25.8685 y 2 quincenas 17.27; disponible hoy A 1000000 y B 690550; tarifa de tabla (term_table) 80050.

| Operación | Prueba | 4xx cubiertos |
|---|---|---|
| listCredits | TestOp_listCredits | 400 kind, 400 limit, 401 |
| getCreditTotals | TestOp_getCreditTotals | 401 |
| getDebtTrend | TestOp_getDebtTrend | 400 credit_id |
| getCredit | TestOp_getCredit | 404 |
| previewAvailability | TestOp_previewAvailability | 422 not_revolving, 404 |
| getAvailabilitySeries | TestOp_getAvailabilitySeries | 404 |
| getCardStatus | TestOp_getCardStatus | 422 not_a_card |
| listCreditTermRates | TestOp_listCreditTermRates | 404 |
| createCredit | TestOp_createCredit | 422 IVA (Q2-67) ×2, 409, 400 sin clave |
| updateCredit | TestOp_updateCredit | 422 IVA, 404 |
| createCreditTermRate | TestOp_createCreditTermRate | 409, 422 |
| updateCreditTermRate | TestOp_updateCreditTermRate | 404 |
| deleteCreditTermRate | TestOp_deleteCreditTermRate | 404 |
| listPromotions | TestOp_listPromotions | 400 valid_on |
| createPromotion | TestOp_createPromotion | 422 vigencia |
| updatePromotion | TestOp_updatePromotion | 404 |
| deletePromotion | TestOp_deletePromotion | 409 in_use, 404 |
| getOptimizerSettings | TestOp_getOptimizerSettings | 401 |
| updateOptimizerSettings | TestOp_updateOptimizerSettings | 422 null, 422 fuera de rango |
| compareScenarios | TestOp_compareScenarios | 422 monto (y 200 missing_constraints / ok) |
| listCreditPurchases | TestOp_listCreditPurchases | 400 family |
| createCreditPurchase | TestOp_createCreditPurchase | 422 partes, 400 monto no entero |
| previewCreditPurchase | TestOp_previewCreditPurchase | 422 sin tasa del plazo |
| getCreditPurchase | TestOp_getCreditPurchase | 404 |
| updateCreditPurchase | TestOp_updateCreditPurchase | 422 status, 404 |
| deleteCreditPurchase | TestOp_deleteCreditPurchase | 409 has_payments, 404 |
| listScheduledPayments | TestOp_listScheduledPayments | 400 status |
| markScheduledPaymentPaid | TestOp_markScheduledPaymentPaid | 409, 422 sin parte, 404 |
| deleteCreditPayment | TestOp_deleteCreditPayment | 404 |

## Supuestos
Los 10 de `S6-creditos.md` (PROVISIONAL). Además: sin `status`, `listCredits` da los activos; `listScheduledPayments` sin `status` ni `period` da los pendientes de hoy en adelante y con `period` los pendientes del periodo; los totales de quincena y mes reparten cada parte en hogar (holder y co_holder) o Karen (external); el gasto vinculado por `expense_id` debe existir (422 si no).

## Bloqueos
- `docs/estado/bloqueos/S6-creditos.md` (nuevo, 10 puntos).
- Fuera de alcance: falta registrar el área en `backend/cmd/finanzas/areas.go` (lo hace el integrador o backend-base); mientras tanto el servidor real responde 501 en estas 29 rutas.

## Pendientes
- Línea en `areas.go` (abajo). Punto 1 del bloqueo (gasto en un paso) cuando api-gastos-inicio exponga su alta en `dominio`. Punto 4 (canastas) cuando compras exponga el costo de canasta por tienda.

## Propuestas para Axel
1. Decidir la ventana de la tendencia de deuda (punto 2): propongo mes en curso + 11 (proyección), porque la deuda viva del motor mira hacia adelante.
2. Decidir si una compra pasa sola a «liquidada» al completar sus pagos (punto 8): propongo que sí, como hoy.
3. Capturar en Ajustes › Créditos y condiciones las 3 restricciones del comparador; sin ellas responde «faltan restricciones» (K16).

## Paquetes instalados
Ninguno.

## Consumo aproximado
~52 acciones; ~200 KB leídos; ~140 KB escritos (código 130 KB, reporte y bloqueo 12 KB).
