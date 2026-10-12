# Corrección backend de áreas · ronda 2 · 12-oct-2026

Agente **correccion-backend-areas-r2** (Opus). Fuente: `docs/estado/wf3/REVISION-backend-areas-r2.md`, con 1 hallazgo alto (H-1) y 0 medios. Los 5 bajos no estaban asignados y no los toqué. No hice git add ni commit, no usé la red ni Telegram, y no hay credenciales reales: las pruebas usan `testutil.FakeCredential`. HEAD = `fae1890`.

## Hecho

| Hallazgo | Estado | Qué cambió | Archivos |
|---|---|---|---|
| H-1 (alta) · S6-03 (c, d) y verificación global: `getDashboardBestBasket` en 501 | **CORREGIDO** | **api-compras-a**: exporta `func BestBasket(ctx context.Context, q dominio.Querier, today, monday dominio.Date) (api.DashboardBestBasket, error)` con la firma del hallazgo, y `func WeekParam(r *http.Request) (dominio.Date, error)`, que lee `week` igual que el área (falta → 422 `required`, otro formato → 400). Para no duplicar el armado, `lista.go` separa tres funciones que ahora usan también `getBasketComparison` y `compareExpense`: `comparableList` (la canasta S5 de la lista), `expenseCart` (la compra base de S5 §7.2) y `firstComplete` (la ganadora de S5 §7.4 en modo compra). `handlers` acepta un «hoy» fijo (`day`), así `buildList` y `loadWorld` corren sin `Registry`. Las cifras salen tal cual de `precios.CompareBasket` y `CompareExpense`: total, cobertura, precios viejos, membresía y diferencia con la siguiente tienda. No se calcula ninguna aparte. **api-gastos-inicio**: el manejador `getDashboardBestBasket` en `inicio.go` (lee `week` con `comprasa.WeekParam` y «hoy» del reloj de la API) y `reg.Handle("getDashboardBestBasket", …)` en `gastos.go`. **Pruebas**: `TestOp_getDashboardBestBasket` y `TestBestBasket_MismaReglaQueLaComparativa` / `TestBestBasket_SemillaSinPrecios`. **Bloqueo**: `S6-gastos-inicio` §15, con el estado y los puntos PROVISIONAL 15.1 a 15.5. | `backend/internal/api/comprasa/comprasa.go`, `backend/internal/api/comprasa/lista.go`, `backend/internal/api/comprasa/mejor_canasta_test.go` (nuevo), `backend/internal/api/gastos-inicio/inicio.go`, `backend/internal/api/gastos-inicio/gastos.go`, `backend/internal/api/gastos-inicio/mejor_canasta_test.go` (nuevo), `docs/estado/bloqueos/S6-gastos-inicio.md` |

## Archivos

- `backend/internal/api/comprasa/comprasa.go`: `handlers.day`, `today()`, `WeekParam`, `BestBasket` y `bestOf`.
- `backend/internal/api/comprasa/lista.go`: se separan `comparableList`, `expenseCart` y `firstComplete`, sin cambiar el comportamiento.
- `backend/internal/api/comprasa/mejor_canasta_test.go` (nuevo)
- `backend/internal/api/gastos-inicio/inicio.go`: el manejador.
- `backend/internal/api/gastos-inicio/gastos.go`: una línea `reg.Handle`.
- `backend/internal/api/gastos-inicio/mejor_canasta_test.go` (nuevo)
- `docs/estado/bloqueos/S6-gastos-inicio.md` (§15)
- este reporte

## Verificación

- **S6-03 (c)**: `go test ./internal/api/ -run TestContrato_Rutas -v` → «S6-03: contrato 216 operaciones, servidor 216, sobrantes 0, faltantes 0, **sin implementar 0**», PASS (antes había 1 sin implementar).
- **S6-03 (d)**: `go test ./internal/api/ -run TestContrato -v` → «S6-03: **216 de 216** operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)» y `--- PASS: TestContrato_Operaciones/getDashboardBestBasket`.
- **S6-03, prueba por operación**: `go test ./internal/api/gastos-inicio/ -run TestOp_getDashboardBestBasket -v` → PASS. Revisa:
  - 401 sin sesión, 422 sin `week` y 400 con `week=2026-42`.
  - 200 validado contra `openapi.yaml` en tres casos: vacío (`shopping_list`, sin `store_id`, `item_count` 0); con la lista vacía (`basket` = `expense`, `expense_id` = la compra); y con la lista (`shopping_list`).
  - En los dos últimos casos, `store_id`, `total_cents`, `covered_count`, `stale_price_count`, `assumes_membership`, `second_store_id` y `difference_vs_second_cents` son iguales a lo que da `getBasketComparison` para la misma canasta. La prueba los compara campo por campo, sin cifras escritas a mano.
  - Ximena recibe 200.
- **Misma regla (exportación)**: `go test ./internal/api/comprasa/ -run 'TestBestBasket|TestOp_getBasketComparison|TestOp_listBaskets' -v` → 4 PASS.
  - `TestBestBasket_MismaReglaQueLaComparativa` comprueba que coincide con `getBasketComparison` y que `item_count` es igual al de `listBaskets`.
  - `TestBestBasket_SemillaSinPrecios` carga la semilla real `datos/semilla`. Resultado: `basket` = `expense` con la compra del 10-oct-2026 y sin `store_id`, porque la semilla no trae precios (ver PROVISIONAL 15.5).
  - `getBasketComparison` y `listBaskets` siguen en PASS después de separar las funciones.
- **S6-05 (pruebas que fallaban por H-1)**:
  - `go test ./cmd/finanzas/ -run TestAreas_SinOperacionesEn501 -v` → «operaciones registradas: 216 de 216», PASS.
  - `go test ./internal/auth/ -run TestPermisos_MismasOperacionesAmbasCuentas -v` → «212 operaciones con el mismo resultado para las dos cuentas», PASS.
- **Cobertura**: `go test ./internal/api/gastos-inicio/ ./internal/api/comprasa/ -count=1 -cover` → gastos-inicio 81.9 % y comprasa 78.1 % (antes 77.5 %). Ninguno de los dos tiene umbral K25.
- **Calidad**: `gofmt -l` vacío y `go vet` limpio en los dos paquetes.
- **Token**: el grep del patrón del token en los dos paquetes da 0.
- **verify-wf3**: `make verify-wf3-go verify-wf3-cobertura verify-wf3-contrato` → rc 0:
  - «gofmt -l .: sin archivos por formatear», «go mod tidy -diff: sin cambios» y «go test ./...: 32 paquetes OK, 5 sin pruebas».
  - K25 OK: dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 %.
  - «los 23 archivos coinciden byte a byte con el generador» y «216 operaciones en 172 rutas, 0 diferencias».
- **Verificación global**: `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → **rc 0** (antes, Error 1 en Makefile:233). La salida termina con «WF-1: verificación completa OK (Go y cobertura K25, oráculos S2 a S5, recetario, tokens y contraste, U4 y paquete)» y «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». También dice «C1 RESULTADO: TODO OK» y «s2 · 32 fixtures, 32 coinciden, 0 diferencias». Tardó 34 s.

## Supuestos

- **HECHO**:
  - La firma es la del hallazgo (`BestBasket`) y sustituye a la solicitud `WeekBasketComparison` de correccion-backend-core-r1.
  - Agregué `WeekParam` para que la tarjeta lea `week` con la misma regla que el área, en vez de copiar el análisis.
  - La ganadora es `cheapest_store_id` de S5 §7.4 (la primera tienda completa), igual que en la comparativa. Si ninguna tienda está completa, la tarjeta queda vacía.
- **Alcance**:
  - El hallazgo nombra `comprasa.go` e `inicio.go`. Toqué además dos archivos de los mismos dueños, solo porque la corrección los exige:
    - `lista.go` (api-compras-a), para separar sin duplicar las funciones de `getBasketComparison`. Su comportamiento no cambia y sus pruebas siguen en PASS.
    - `gastos.go` (api-gastos-inicio), que es donde se hace el `reg.Handle` que pide el hallazgo.
  - Las dos pruebas van en archivos nuevos para no tocar las de los dueños.
- **PROVISIONAL** (`docs/estado/bloqueos/S6-gastos-inicio.md` §15):
  - 15.1: cuándo se considera vacía la lista.
  - 15.2: `item_count` = el mismo número de `listBaskets`.
  - 15.3: `covered_count` = idénticos + sustitutos de la ganadora.
  - 15.4: si la tienda siguiente está incompleta, no hay segunda tienda.
  - 15.5: si la compra no se puede comparar, la tarjeta responde 200 vacía en lugar de 422.

## Bloqueos

- `docs/estado/bloqueos/S6-gastos-inicio.md` §15, actualizado con el estado y los puntos 15.1 a 15.5. No creé bloqueos nuevos.

## Pendientes

- Ninguno de H-1.
- Los bajos de la revisión (B-1 a B-5) no estaban asignados a esta corrección.
- El comentario de `gastos.go:16` («las demás de esta área siguen en 501») ya no es cierto. Lo dejé así porque no es parte del hallazgo.

## Propuestas para Axel

- **PREGUNTA 15.4**: si la segunda tienda está incompleta, ¿la tarjeta dice «$X menos que [tienda]»? El motor S5 no da esa diferencia.
- **PREGUNTA 15.5**: si la compra no se puede comparar, ¿la tarjeta queda vacía («Aún no hay precios para comparar»)? Hoy pasa con la semilla real, que no trae precios.
- **PROPUESTA**: confirmar 15.1 a 15.3 en la puerta 4, junto con los valores operativos de K31.

## Paquetes instalados

Ninguno.

## Consumo aproximado

Unas 36 acciones, unos 115 KB leídos y unos 30 KB escritos.
