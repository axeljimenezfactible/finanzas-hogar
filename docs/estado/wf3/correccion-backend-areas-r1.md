# Corrección backend de áreas · ronda 1 · 12-oct-2026

Agente: `correccion-backend-areas-r1` (Opus). Entrada: `docs/estado/wf3/REVISION-backend-areas-r1.md` (H-1 a H-4 altos y M-1 medio). Seguí `CONVENCIONES-backend.md` §3, §4 y §12. No toqué migraciones, motores, `routes.go`/`types_gen.go`, Makefile ni nada de Pista 3. No hice `git add` ni `git commit`. Temporales en `datos/tmp/correccion-ba-r1/`, ya borrados.

## Hecho

| Hallazgo | Qué hice | Archivo |
|---|---|---|
| H-1 (alta) · S7-10, S7-11, S7-02, S6-03 | **La causa ya estaba corregida cuando empecé.** `dominio/presupuesto_motor.go` ya leía `category.code` con `sql.NullString` y `engineCategoryCode(id, code)` (clave `cat#<id>` solo en memoria). El archivo lo cambió otro agente a las 04:26:37 y es de api-presupuesto, fuera de mi alcance, así que no lo toqué. Lo mío fue la prueba de s7-pdf. Agregué `userCategory(t, e)`, que crea «Mascotas» con `createCategory` (201) y comprueba en la base que `code IS NULL`. Se usa en `TestCruceMensual`, `TestVistaMensual` y `TestMensual`, que ahora registran `ajustes.Register`. Con esa categoría: el job mensual termina en `ok`, `getMonthlyReport` da 200 y `generateMonthlyReport` da 202. | `backend/internal/reportes/reportes_test.go` |
| H-2 (alta) · S7-02 | `TestCruceMensual` ahora cruza **P1–P6 completos**. `canon(t, e)` agrega a la despensa 1 de `setup` las despensas 2, 3 y 4: gasto de Despensa con devengo de 2 meses, compra a plazos con partes y 4 pagos con partes. También agrega el papel alterno **de contado** (paquetes 3 y 5, Karen, `payment_method = 'cash'`, `paid_by_person_id` = Karen, como dice el comentario de la columna en `00001_esquema.sql`). **Ninguna cifra está escrita a mano**: todo se lee de `oraculo/model.json` (`despensas`, `disposicion_por_despensa`, `papel`), y la prueba verifica que cada despensa cuadre consigo misma. Hay dos comprobaciones. (1) `getAccrualVsCash` de oct-2026 a mar-2027 = `periodos[0..5]`, es decir R-45 a R-50. (2) Para cada uno de los 6 meses: se genera el PDF, el cruce completo PDF = API (`getMonthlyReport` + `getAccrualVsCash`) da 0 diferencias y el PDF trae las 3 cifras de su fila R. El reloj se pone en el día 1 de cada mes a las 12:00, porque el reporte solo se genera de un mes que ya empezó. Ningún evento de la semilla cae del día 1 al 9 (S2-09). No hizo falta ningún bloqueo. | `backend/internal/reportes/reportes_test.go` |
| H-3 (alta) · S7-14 y verificación global | `gofmt -w` en los dos archivos. Solo cambió la alineación: el `struct env` y el campo `Note:` de `ExternalDoc`. `gofmt -l .` en `backend/` sale vacío. | `backend/internal/reportes/vistas.go`, `backend/internal/reportes/reportes_test.go` |
| H-4 (alta) · S6-04, S6-03 | `setStoreProductGroups` ahora llama a `v.Required(in.IDs != nil, "ids")`. `{}` y `{"ids": null}` responden **422 `validation_failed`** y no tocan los grupos. `{"ids": []}` sigue vaciando la tienda (204). Agregué a `TestOp_setStoreProductGroups` los dos 422. También comprueba con `listStores` que los grupos siguen siendo 2 y que `{"ids": []}` los deja en 0. | `backend/internal/api/comprasa/tiendas.go`, `backend/internal/api/comprasa/comprasa_test.go` |
| M-1 (media) · S6-03 / CONVENCIONES §12 | `TestOp_sendReport` cubre ahora el **202**. Genera el mensual de oct-2026, pone un Enviador falso (`fakeSender`, implementa `reportes.DocumentSender`, responde 200 sin red) y suscribe a Axel, Ximena y Karen a `monthly_report`. `c.Do` valida el 202 contra `openapi.yaml`, con `delivery_status = "sent"`, 2 envíos (Axel y Ximena) y ninguno a Karen. Siguen el 409 (estado de cuenta) y el 404. | `backend/internal/api/reportes/reportes_test.go` |

## Archivos

- `backend/internal/reportes/reportes_test.go` (H-1, H-2, H-3)
- `backend/internal/reportes/vistas.go` (H-3, solo gofmt)
- `backend/internal/api/comprasa/tiendas.go` (H-4)
- `backend/internal/api/comprasa/comprasa_test.go` (H-4)
- `backend/internal/api/reportes/reportes_test.go` (M-1)
- `docs/estado/wf3/correccion-backend-areas-r1.md` (este reporte)

## Verificación (comandos de `backend/` o de la raíz)

- **S7-02**: `go test ./internal/reportes/ -run TestCruceMensual -count=1 -v`
  - R-45 a R-50 en la API, iguales al oráculo:

    | Fila | Mes | Devengado | Efectivo | Diferencia |
    |---|---|---|---|---|
    | R-45 | 2026-10 | $1,947.50 | $973.75 | $973.75 |
    | R-46 | 2026-11 | $1,947.50 | $1,947.50 | $973.75 |
    | R-47 | 2026-12 | $2,517.37 | $2,232.43 | $1,258.69 |
    | R-48 | 2027-01 | $2,716.37 | $2,716.36 | $1,258.70 |
    | R-49 | 2027-02 | $2,517.37 | $2,517.38 | $1,258.69 |
    | R-50 | 2027-03 | $2,716.37 | $2,716.36 | $1,258.70 |

  - PDF = API, mes por mes: 38 de 38, 40 de 40, 43 de 43, 43 de 43, 46 de 46 y 46 de 46 cifras y fechas iguales, 0 diferencias. **M = 256 en 6 PDF.** Cada PDF trae las 3 cifras de su fila R. Campos que cruza un PDF (46 en mar-2027): `fortnights.*` (fechas y montos), `budget.*`, `accrual:<mes>|accrued_cents/cash_cents/cumulative_difference_cents`, `categories.rows.*`, `credits.*`, `live_debt_end_cents`, `external_participation.0.*`. PASS.
- **S7-10**: `go test ./internal/reportes/ -run 'TestMensual$' -v` con la categoría de la API → PDF v1 guardado (96,403 bytes); con la plantilla rota, el job queda en error y no hay envíos; con día 31 sale el 30-nov-2026. PASS.
- **S7-11**: `-run TestVistaMensual` → «S7-11 vista web = PDF: 34 cifras de getMonthlyReport, 0 diferencias». PASS.
- **S6-03 / M-1**: `go test ./internal/api/reportes/ -run TestOp_sendReport -v` → PASS (202 validado contra `openapi.yaml`, 409 y 404).
- **S6-04 / H-4**: `go test ./internal/api/comprasa/ -run TestOp_setStoreProductGroups -v` → PASS (204, 404, 422 por id inexistente, 422 con `{}` y con `{"ids": null}`, `{"ids": []}` = 204 y vacía).
- **S6-03 contrato**: `go test ./internal/api -run TestContrato_Operaciones -count=1` → «S6-03: **215 de 216** operaciones cumplen».
  - Antes fallaban 4 de mi grupo: `generateMonthlyReport`, `getMonthlyReport`, `sendReport` y `setStoreProductGroups`. Ahora pasan las 4.
  - La única que falla es `getDashboardBestBasket` (501, sin registrar). Es de api-gastos-inicio y está fuera de esta revisión.
  - Si `setStoreProductGroups` se corre **sola** con `-run …/^setStoreProductGroups$`, da 404. No es un defecto: la tienda la crea `createStore` en la corrida completa.
- Paquetes tocados: `go test ./internal/reportes/... ./internal/api/reportes/ ./internal/api/comprasa/ -count=1 -cover` → todo ok. Cobertura: reportes 77.3 %, reportes/avisos 77.3 %, api/reportes 76.1 %, api/comprasa 77.5 % (son paquetes sin umbral K25).
- **S7-14 / H-3**: `make verify-wf3-go` → «go build ./... y go vet ./...: sin errores», «**gofmt -l .: sin archivos por formatear**», «go mod tidy -diff: sin cambios». `go test ./...` → 31 paquetes ok. Sale con Error 1 en Makefile:326 **solo** por `internal/api`: `TestContrato_Operaciones` y `TestContrato_Rutas` fallan por `getDashboardBestBasket` (501), fuera de mi alcance.
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:
  - Pasan: verify-go (antes Error 1 en Makefile:114; ahora «gofmt -l .: sin archivos por formatear»), WF-0, C1 (`build.py --check`: 23 archivos iguales; `verify_c1.py --db`: TODO OK), S1, U1, criterios, oráculos S2 a S5, recetario, C2 y U4.
  - Queda **un solo** Error 1, en `verify-wf1-cobertura` (Makefile:233). La causa es la misma: `getDashboardBestBasket` en 501.
- `make verify-wf3-contrato` → «23 archivos coinciden», «S6-03: 216 operaciones en 172 rutas, 0 diferencias», PASS ×2.
- `make verify-wf3-cobertura` → K25 OK en dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 %. Corta en `api` por la misma operación en 501.
- **S7-06**: el grep del patrón de token en los 3 paquetes tocados → 0. El Enviador falso no usa token.

## Supuestos

- H-1: no edité `presupuesto_motor.go`. Ya estaba corregido con `sql.NullString` cuando llegué y es de api-presupuesto. Si ese cambio se revierte, mis pruebas (`userCategory` en las 3 pruebas) lo detectan.
- H-2: así representé las despensas 2 a 4 y el papel de contado en la fixture:
  - Cada despensa es un gasto de Despensa a crédito, con devengo de 2 meses y su `credit_purchase`, igual que la despensa 1 de la semilla.
  - El papel de contado es un gasto de la categoría `paper` (`dominio.PaperCategoryCode`), en efectivo y pagado por Karen.

  Es una fixture de prueba y no toca la semilla real ni el bloqueo `semilla-despensas-futuras`, que en la base real siembra las despensas 2 a 4 sin gasto. Que las 12 cifras y los 6 cruces den 0 diferencias confirma que esa representación es la que el motor S2 espera.

## Bloqueos

Ninguno nuevo. No hizo falta ningún valor sin decidir.

## Pendientes

- `getDashboardBestBasket` en 501. Es lo único que mantiene en rojo `make verify-wf3-go`, `make verify-wf3-cobertura` y `make verify-all` (Makefile:233). Viene de la revisión core-r1 (api-gastos-inicio, con apoyo de api-compras-a). No está en mis hallazgos, así que no lo toqué. Va en «solicitudes fuera de alcance».
- Los hallazgos bajos B-1 a B-5 no se me asignaron.

## Propuestas para Axel

- Ninguna que cambie decisiones. Para el orquestador repito la propuesta del revisor: subir `TestContrato_Operaciones` a `verify-wf3-contrato`. El Makefile no es mío.

## Paquetes instalados

Ninguno.

## Consumo aproximado

~38 acciones; ~170 KB leídos; ~30 KB escritos (código, pruebas y este reporte).
