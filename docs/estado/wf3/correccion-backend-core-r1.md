# correccion-backend-core-r1 · corrección de la revisión backend-core ronda 1

12-oct-2026 · agente **correccion-backend-core-r1** (Opus). Fuente: `docs/estado/wf3/REVISION-backend-core-r1.md` (6 altos, 3 medios). Alcance: archivos de backend-base, auth-sistema, semilla, api-gastos-inicio, api-presupuesto y api-comidas-salud que nombran los hallazgos. Sin git add ni commit, sin red, sin Telegram, sin credenciales reales (las pruebas usan `testutil.FakeCredential` y un token falso en un archivo temporal). **Cerrado por tope de acciones** (≈ 60): queda 1 operación pendiente (H2, `getDashboardBestBasket`), que además depende de un cambio fuera del alcance.

## Hecho (por hallazgo)

| Hallazgo | Estado | Qué se hizo |
|---|---|---|
| H1 (alta) S6-02, S6-04 · categoría sin `code` → 500 | **CORREGIDO** | `dominio/presupuesto_motor.go`: `code` se lee en `sql.NullString`; si es NULL, `engineCategoryCode` da `cat#<id>` solo en memoria (no cabe en el CHECK `a-z0-9_`, no choca con un código real; las respuestas regresan a `category_id` con `FinEngine.CategoryID`, así que no cambia nada visible ni ninguna cifra). Nota en `S6-presupuesto-motor` punto 5 (no PROVISIONAL: nada visible). Prueba `TestOp_getBudgetMonth_CategoriaSinCodigo` (mes, categorías, tendencia filtrada por la categoría sin código, quincenas, alertas, deuda viva, ahorro y flujo → 200). |
| H2 (alta) S6-03, C3-03 · 2 operaciones en 501 | **PARCIAL** | `getDashboardWeekMeals` HECHO: `comidas.WeekMeals` (archivo nuevo `internal/api/comidas/inicio.go`, mismo motor S4 y mismo `fillWeek` que `getMealPlanWeek`), manejador y registro en `gastos-inicio/inicio.go` y `gastos.go`; semana del `period` o la semana en curso (K8, no K38). Prueba `TestOp_getDashboardWeekMeals` (costo 1000 = `week_cost_cents`, 20 casillas sin armar, Axel 2 días y 530 kcal, semana vacía, domingo sigue en la semana en curso, 422, 400, 401). Qué es «día con plan» («fuera» cuenta): **PROVISIONAL**, `S6-comidas-calculos` punto 14. `getDashboardBestBasket` **sigue en 501**: necesita exportar la comparación de `internal/api/comprasa` (api-compras-a, fuera del alcance). Ver Solicitudes. |
| H3 (alta) S6-04, S6-02 · borrador sin `lines` | **CORREGIDO** | `gastos-inicio/tickets.go`: `v.Required(d.Lines != nil, "lines")` → 422 `validation_failed` con `fields`. `TestOp_saveTicketScanDraft` comprueba el 422, el campo `lines` y que el borrador guardado no cambia. |
| H4 (alta) S6-03 c, S6-05, C3-03 · `areas.go` vacío | **CORREGIDO** (215/216 por H2) | `cmd/finanzas/areas.go` con las 10 áreas (auth incluye sistema). Prueba nueva `cmd/finanzas/areas_test.go` `TestAreas_SinOperacionesEn501`: recorre `api.Routes` con `Registry.Has` (sin llamar a operaciones que envían avisos), hace `login` con el registro del binario (200) y una operación de área (200). Hoy falla solo por `getDashboardBestBasket`. |
| H5 (alta) S6-01 · gofmt en `internal/reportes` | **Ya resuelto por su dueño** | No lo toqué (es de s7-pdf y lo tiene REVISION-backend-areas-r1 H-3). Al correr mis verificaciones `gofmt -l .` ya sale vacío. |
| H6 (alta) S6-02, S6-04 · `sendReport` y `setStoreProductGroups` | **Ya resuelto por sus dueños** | No los toqué (s7-pdf y api-compras-a; REVISION-backend-areas-r1 H-4 y M-1). En mi corrida `TestContrato_Operaciones` solo falla en `getDashboardBestBasket`. |
| M1 (media) S6-05 · permisos solo con `auth.Register` | **CORREGIDO** (falla solo por H2) | `auth/auth_test.go`: `todasLasAreas` = la lista del binario (Telegram de s7-avisos a un servidor falso con token falso); exige mismo código para Axel y Ximena y ningún 403 ni 501. Resultado: 212 operaciones iguales; única falla `getDashboardBestBasket: Axel 501, Ximena 501`. |
| M2 (media) · paquete `gastos-inicio` vs convención | **CORREGIDO** (opción «actualizar la convención») | CONVENCIONES-backend 1.1 §2 y §14: `internal/api/gastos-inicio`, `package gastos`, import con alias. No se movió el paquete: moverlo obligaba a editar `internal/api/contrato_test.go` (pruebas-contrato) y `internal/reportes/avisos/resumen_test.go` (s7-avisos), fuera del alcance. `S6-gastos-inicio` §8 cerrado. |
| M3 (media) S6-16 · bloqueos `semilla-*` | **CORREGIDO** | Renombrados a `S6-semilla-despensas-futuras.md`, `S6-semilla-sucursales-municipio.md`, `S6-semilla-catalogo-alias.md`; citas actualizadas en `semilla.md`, `internal/semilla/compra.go`, `internal/semilla/otros.go` y en los propios archivos; alta en `bloqueos/INDICE.md` (sección «Bloqueos de la semilla de WF-3a»). `ls docs/estado/bloqueos/S6-*` ya los encuentra. |

## Archivos

Código: `backend/internal/dominio/presupuesto_motor.go` · `backend/internal/api/presupuesto/presupuesto_test.go` · `backend/internal/api/gastos-inicio/{tickets.go,gastos.go,inicio.go,gastos_test.go}` · `backend/internal/api/comidas/inicio.go` (nuevo) · `backend/internal/api/comidas/inicio_test.go` (nuevo) · `backend/cmd/finanzas/areas.go` · `backend/cmd/finanzas/areas_test.go` (nuevo) · `backend/internal/auth/auth_test.go` · `backend/internal/semilla/{compra.go,otros.go}` (solo el nombre del bloqueo en un texto).
Documentos: `docs/estado/wf3/CONVENCIONES-backend.md` (1.1) · `docs/estado/wf3/semilla.md` · `docs/estado/bloqueos/S6-semilla-*.md` (3, renombrados) · `docs/estado/bloqueos/INDICE.md` · `docs/estado/bloqueos/S6-presupuesto-motor.md` (punto 5) · `docs/estado/bloqueos/S6-comidas-calculos.md` (punto 14) · `docs/estado/bloqueos/S6-gastos-inicio.md` (§8, §11, §15) · este reporte.

## Verificación (comandos y salidas resumidas)

- **S6-01**: `go version` → go1.27.2 · `go build ./...` → OK · `go vet ./...` → OK · `gofmt -l .` → vacío.
- **S6-02**: `go test ./... -count=1` → 29 paquetes `ok`, 3 `FAIL` (`cmd/finanzas`, `internal/api`, `internal/auth`), las 3 **solo** por `getDashboardBestBasket` en 501. Cobertura: dominio **87.0 %** (≥ 85 %), api/comidas 81.7 %, api/gastos-inicio 81.8 %, api/presupuesto 77.6 %, auth 83.8 %.
- **S6-03**: `go test ./internal/api/ -run TestContrato -v` → «215 de 216 operaciones cumplen», «contrato 216, servidor 216, sobrantes 0, faltantes 0, sin implementar 1» (`getDashboardBestBasket`). Antes: 191 de 216 y 2 sin implementar. `TestAreas_SinOperacionesEn501` → «1 de 216 operaciones sin registrar: [getDashboardBestBasket]».
- **S6-04**: `TestOp_saveTicketScanDraft` PASS (422 sin `lines`); `TestOp_getBudgetMonth_CategoriaSinCodigo` PASS (antes 500).
- **S6-05**: `TestPermisos_MismasOperacionesAmbasCuentas` → «212 operaciones con el mismo resultado para las dos cuentas»; 0 respuestas 403; 1 respuesta 501 (H2).
- **S6-16**: `ls docs/estado/bloqueos/S6-*` incluye los 3 `S6-semilla-*.md`.
- `make verify-all` / `make verify-wf3`: no los corrí completos por el tope de acciones; su parte Go (`gofmt`, `vet`, `test`) es la de arriba y falla solo por H2.

## Supuestos

- La clave interna `cat#<id>` no es visible (todas las salidas del motor vuelven a `category_id`); por eso no se marca PROVISIONAL.
- `days_with_plan` cuenta días con receta o «fuera»; `avg_kcal_per_day` = Σ kcal ÷ esos días, ROUND_HALF_UP, ausente con 0 días; la meta es la que rige el último día de la semana que la tiene (igual que `getMealPlanTrends`). PROVISIONAL, `S6-comidas-calculos` punto 14.
- La tarjeta de comidas usa la semana en curso sin `period` (K8, U3 tarjeta 4); `period` que no es semana → 422.

## Bloqueos

- `getDashboardBestBasket` depende de una exportación de `internal/api/comprasa` (ver Solicitudes). Registrado en `S6-gastos-inicio` §15.
- PROVISIONAL nuevo: `S6-comidas-calculos` punto 14 (qué es «día con plan»).

## Pendientes

1. `getDashboardBestBasket` (H2): cuando exista la exportación de comprasa, agregar en `backend/internal/api/gastos-inicio/inicio.go` el manejador (parámetro `week` obligatorio, R3-01) que convierte la comparación en `api.DashboardBestBasket`: `store_id` = `cheapest_store_id`; `total_cents`, `stale_price_count` (= `stale_count`) y `assumes_membership` de esa tienda; `second_store_id` y `difference_vs_second_cents` = la tienda siguiente en `stores` (orden del motor S5); `item_count` = artículos de la canasta; `covered_count` = `identical_count + substitute_count` de la ganadora (PROVISIONAL: si un sustituto cuenta como «con precio»); registrarlo en `gastos.go` y agregar `TestOp_getDashboardBestBasket` con la compra semilla del 10-oct-2026 (Bodega Aurrera 309450, Walmart 323300, diferencia 13850; U3 tarjeta 5).
2. Correr `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` y `make verify-wf3` completos después de (1).
3. (Fuera de esta revisión, lo nota REVISION-backend-areas-r1 B-5) `cmd/finanzas/main.go` no llama a `reportes.SetSender(...)` en `serve`.

## Solicitudes fuera de alcance

- **api-compras-a** · archivo nuevo `backend/internal/api/comprasa/inicio.go`: exportar `func WeekBasketComparison(ctx context.Context, reg *api.Registry, monday dominio.Date) (*api.BasketComparison, error)` que devuelve la misma comparación que `getBasketComparison` para la lista de la semana `monday` y, si la lista no tiene artículos con producto y cantidad, la de la última compra con renglones (`dominio.ShopA.ExpenseBaskets(..., "ORDER BY e.expense_date DESC, e.id DESC")`, basket = `expense`); `nil` si no hay nada que comparar. Para no duplicar, separar en `lista.go` el cálculo de `getBasketComparison` (desde `buildList` hasta `membership`) y de `compareExpense` en funciones que devuelven `*api.BasketComparison` sin escribir la respuesta. Lo pide el bloqueo `S6-gastos-inicio` §15 y el hallazgo H2.

## Propuestas para Axel

- PREGUNTA (`S6-comidas-calculos` punto 14): en la tarjeta «Comidas de la semana», ¿un día en que la persona solo tiene «fuera» cuenta como día con plan? A (hoy): sí, y baja el promedio de kcal; B: no, solo días con receta.

## Paquetes instalados

Ninguno.

## Consumo aproximado

≈ 60 acciones · ≈ 150 KB leídos · ≈ 30 KB escritos.

---

## Continuación 1

12-oct-2026 · mismo agente, continuación tras el cierre por tope. Hice **solo** lo pendiente; no rehíce nada de lo anterior.

### Hecho

| Pendiente | Estado | Qué se hizo |
|---|---|---|
| 1 · `getDashboardBestBasket` (H2) | **SIGUE BLOQUEADO** (fuera de alcance) | HECHO: confirmé que `internal/api/comprasa` sigue sin exportar nada fuera de `Register` (`grep '^func [A-Z]'`), así que la condición del pendiente («cuando exista la exportación») no se cumple. No dupliqué en `gastos-inicio` el armado de la canasta (`loadWorld`, `buildList`, `CompareBasket`, `compareExpense`, `membership`: unas 300 líneas de `comprasa/motor.go` y `comprasa/lista.go`). La tarjeta debe dar las mismas cifras que `getBasketComparison` (U3 D-5 usa la misma canasta), y una copia podría dar otras. Dejé el estado en `S6-gastos-inicio` §15 («Continuación 1»). La solicitud de abajo sustituye a la PROPUESTA `BestBasket` de §15. |
| 2 · `make verify-all` y `make verify-wf3` completos | **HECHO** | Ver Verificación. La única causa de rojo en Go sigue siendo H2. |
| 3 · `reportes.SetSender` en `serve` | No es de esta revisión | Sigue en REVISION-backend-areas-r1 B-5; no lo toqué. |

### Archivos (continuación)

- `docs/estado/bloqueos/S6-gastos-inicio.md` (§15: una viñeta «Continuación 1»).
- Este reporte (esta sección).
- No toqué código.

### Verificación (continuación)

- **S6-01**: `go build ./...` y `go vet ./...` → sin errores; `gofmt -l .` → vacío; `go mod tidy -diff` → sin cambios (`make verify-wf3-go`, líneas 2-5).
- **S6-02**: `go test ./... -count=1` → 29 `ok` y 3 `FAIL` (`cmd/finanzas`, `internal/api` e `internal/auth`). Los 3 fallan solo por `getDashboardBestBasket: 501`. Cobertura K25: dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 % («K25 OK» en `verify-wf3-cobertura`). Paquetes de esta revisión: api/comidas 81.7 %, api/gastos-inicio 81.8 %, api/presupuesto 77.6 %, semilla 79.4 %.
- **S6-03**: `TestContrato_Operaciones` → «215 de 216 operaciones cumplen». `TestContrato_Rutas` → «contrato 216, servidor 216, sobrantes 0, faltantes 0, sin implementar 1». `TestAreas_SinOperacionesEn501` → «1 de 216 … [getDashboardBestBasket]». `verify-wf3-contrato` → «23 archivos coinciden byte a byte», `TestRutasIgualesAlContrato` (216 operaciones en 172 rutas, 0 diferencias) PASS.
- **S6-05**: `TestPermisos_MismasOperacionesAmbasCuentas` → «212 operaciones con el mismo resultado para las dos cuentas»; 0 respuestas 403; única falla: `getDashboardBestBasket: Axel 501, Ximena 501`.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`** → exit 2. Pasan WF-0 (SQL, Go, colorimetría, OpenAPI: «WF-0: verificación completa OK»), C1, S1, U1 y criterios. Se detiene en `verify-wf1-cobertura` (Makefile:233), solo por las 3 pruebas de H2.
- **`make -k verify-wf3`** → exit 2. `verify-wf3-go` (Makefile:326) y `verify-wf3-cobertura` (Makefile:331) fallan solo por H2. `verify-wf3-contrato` → OK. `verify-wf3-frontend` dio 77 fallas, **todas** `net::ERR_CONNECTION_REFUSED` en `localhost:4173` (el servidor de vista previa dejó de responder a media corrida). Las 77 son del entorno, no del código: al repetir solo `make verify-wf3-frontend` → **142 passed, exit 0**.
- Temporales: borré `datos/tmp/cbc-r1/` (mis bitácoras). **No** borré `datos/tmp/verify-wf3/` ni `datos/tmp/playwright-frontend/`: los crean los objetivos del Makefile, y al terminar había 13 procesos de verificación de otros agentes (playwright, vite, make verify) usándolos. Esa concurrencia explica también el `ERR_CONNECTION_REFUSED` en 4173 de la primera corrida.

### Supuestos (continuación)

- Ninguno nuevo. No decidí nada de producto.

### Bloqueos (continuación)

- `getDashboardBestBasket` sigue esperando la exportación de `internal/api/comprasa` (api-compras-a). Ver `S6-gastos-inicio` §15.

### Pendientes (continuación)

1. **api-compras-a** (`backend/internal/api/comprasa/`): exportar la comparación (solicitud 1).
2. **api-gastos-inicio** (`backend/internal/api/gastos-inicio/inicio.go` y `gastos.go`), después del 1:
   - Agregar el manejador `getDashboardBestBasket` y registrarlo.
   - `week` obligatorio (`ListWeek`, R3-01): falta → 422 `validation_failed` con `fields[week]`; mal formado → 400.
   - Convertir `*api.BasketComparison` a `api.DashboardBestBasket` así:
     - `basket` y `expense_id` tal cual.
     - `store_id` = `cheapest_store_id`.
     - De esa tienda en `stores`: `total_cents`, `stale_price_count` = `stale_count` y `assumes_membership`.
     - `second_store_id` y `difference_vs_second_cents` = la tienda siguiente en `stores` (orden del motor S5), ausentes si no hay.
     - `item_count` = `len(items)`.
     - `covered_count` = `identical_count + substitute_count` de la ganadora.
     - Sin comparación: `basket = "shopping_list"`, conteos en 0 y sin tienda (vacío de U3: «Aún no hay precios para comparar»).
   - PROVISIONAL, por escribir en `S6-gastos-inicio` §15: (a) si un sustituto cuenta como «con precio»; (b) en `basket=expense`, la «segunda» es la siguiente en `stores` después de la primera tienda completa (S5 §7.4).
   - Agregar `TestOp_getDashboardBestBasket`: 200 con la compra semilla del 10-oct-2026 (Bodega Aurrera 309450, Walmart 323300, diferencia 13850; U3 tarjeta 5), vacío, 422 sin `week`, 400 y 401.
3. Después, volver a correr `go test ./... -count=1`, `make verify-wf3` y `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: deben quedar en verde, porque no hay otra causa de rojo.

### Solicitudes fuera de alcance (continuación)

1. **api-compras-a** · archivo nuevo `backend/internal/api/comprasa/inicio.go`:
   - Firma: `func WeekBasketComparison(ctx context.Context, reg *api.Registry, monday dominio.Date) (*api.BasketComparison, error)`.
   - Devuelve la misma comparación que `getBasketComparison` con `basket=shopping_list` para la semana `monday`.
   - Si la lista no tiene artículos con producto y cantidad a comprar > 0, devuelve la de la última compra con renglones de una tienda activa: `dominio.ShopA.ExpenseBaskets(ctx, q, "ORDER BY e.expense_date DESC, e.id DESC")`, igual que `basket=expense`.
   - Devuelve `nil, nil` si no hay nada que comparar.
   - Sin duplicar código: en `lista.go`, separar el cuerpo de `getBasketComparison` (de `buildList` a `membership`) y el de `compareExpense` en dos funciones que devuelvan `*api.BasketComparison` sin escribir la respuesta. Los manejadores actuales solo las llaman y escriben el JSON; sus cifras no cambian (`TestOp_getBasketComparison` debe seguir igual).
   - Lo piden H2 y `S6-gastos-inicio` §15.
2. **Orquestador** · para cerrar H2 en una sola ronda, dar a **un** agente el alcance de `backend/internal/api/comprasa/` (solicitud 1) y de `backend/internal/api/gastos-inicio/{inicio.go,gastos.go,inicio_test.go}` (pendiente 2). No hace falta tocar `cmd/finanzas/areas.go`: ya registra `comprasa.Register` y `gastos.Register`.

### Propuestas para Axel (continuación)

- Siguen en pie las de arriba (PREGUNTA de `S6-comidas-calculos` punto 14).
- PREGUNTA nueva, ligada a la PROVISIONAL del pendiente 2a, para cuando se construya la tarjeta. En «Mejor canasta», ¿un artículo que la tienda solo tiene como **sustituto** cuenta como «con precio» en «N de M artículos con precio»?
  - A: sí, y cuenta igual que uno idéntico.
  - B: no, solo cuentan los idénticos.

### Paquetes instalados (continuación)

Ninguno.

### Consumo aproximado (continuación)

≈ 16 acciones · ≈ 45 KB leídos · ≈ 6 KB escritos.
