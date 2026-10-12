# api-compras-b · WF-3a «Construcción» · reporte

Agente `api-compras-b` (Opus) · 12-oct-2026 · paquete `backend/internal/api/comprasb` (nombre de CONVENCIONES-backend §2; el encargo decía `compras-b/`, pero un paquete Go no lleva guion) y `backend/internal/dominio/compras-b_repo.go`.

## Hecho

- Las **23 operaciones** de CONVENCIONES §14 registradas con `comprasb.Register` y probadas: 16 de `compras-tiendas-b.yaml` (con las 4 del agente externo, K53) y 7 de `ajustes-c.yaml`.
- Cálculos solo con los motores, sin rehacerlos a mano: precio por unidad base (`precios.BaseUnitPrice`, RN-20/K7), frescura, edad y «viejo» (`precios.Freshness` con los valores vigentes de `price.*`, K22/K40/K44), resumen de precios (`precios.Book.Summary`), regla entre claves de Q2-77 (`precios.Settings.Validate`), estado, existencia estimada, consumo semanal y días restantes (`recetas.Engine.InventoryStatus`), movimientos de inventario (`recetas.AdjustmentMovement`, `LeftoverMovement`) y corte PROFECO (`profeco.MaxBytesFromSettings`, `ReadLinksFile`, `LoadIntoDB`).
- Repositorios en `dominio` con el prefijo `CB` (ya había un `Exists` de api-creditos en el paquete; el prefijo evita choques entre agentes que trabajan al mismo tiempo). Todos reciben `dominio.Querier` o `*sql.Tx`. Cada alta, cambio o baja escribe `audit_log`.
- Lotes del agente externo: cada renglón va en su propia transacción. Un renglón rechazado sale `rejected` con su `FieldError` y no detiene a los demás. Las llaves naturales son `product.name` (sin distinguir mayúsculas) con `presentation.name` dentro del producto, la UNIQUE de `price` y `recipe.name`, así que repetir un lote deja los renglones en `unchanged`. Un lote vacío o más grande que el tope (500, 500 o 100) responde 422.
- Bloqueos nuevos: `docs/estado/bloqueos/S6-compras-b-fuentes.md` y `docs/estado/bloqueos/S6-compras-b-nutricion.md`.

## Tabla operación → prueba (S6-03)

| Operación | Prueba (`backend/internal/api/comprasb/comprasb_test.go`) | 2xx | 4xx probados |
|---|---|---|---|
| listProductPrices | TestOp_listProductPrices | 200 (filtros source, store_id, limit) | 400, 404 |
| getPriceSummary | TestOp_getPriceSummary | 200 | 401 |
| createPrice | TestOp_createPrice | 201 | 409, 422, 400 (tipo, falta Idempotency-Key) |
| updatePrice | TestOp_updatePrice | 200 | 403 (PROFECO), 422, 404 |
| deletePrice | TestOp_deletePrice | 204 (+ audit_log) | 404, 403 |
| getPriceSources | TestOp_getPriceSources | 200 | 401 |
| getPriceSettings | TestOp_getPriceSettings | 200 | 401 |
| updatePriceSettings | TestOp_updatePriceSettings | 200 (cambiar y restablecer) | 409, 422 (regla entre claves y null), 400 |
| importProfeco | TestOp_importProfeco | 201 (sin liga: 0 cargados; con liga: 1) | 415 (RAR), 422 (encabezado), 400 |
| listNutritionFacts | TestOp_listNutritionFacts | 200 | 404 |
| createNutritionFact | TestOp_createNutritionFact | 201 | 409, 422 (×3), 404 |
| getRecipeSettings | TestOp_getRecipeSettings | 200 | 401 |
| reloadNutritionBase | TestOp_reloadNutritionBase | — (PROVISIONAL) | 409, 422 |
| listInventoryItems | TestOp_listInventoryItems | 200 (orden, filtros y cursor) | 400 |
| createInventoryItem | TestOp_createInventoryItem | 201 (+ movimiento count) | 409, 422 (×2) |
| getInventorySummary | TestOp_getInventorySummary | 200 | 401 |
| listInventoryMovements | TestOp_listInventoryMovements | 200 | 404 |
| adjustInventoryItem | TestOp_adjustInventoryItem | 200 | 422 (×2), 404 |
| updateInventoryItem | TestOp_updateInventoryItem | 200 (Merge Patch con null) | 422, 400, 404 |
| getExternalAgentReference | TestOp_getExternalAgentReference | 200 (bearer falso) | 401 (otra credencial y sin credencial) |
| importExternalCatalog | TestOp_importExternalCatalog | 200 (creado, sin cambio, rechazados; repetir no duplica) | 422, 401 |
| importExternalPrices | TestOp_importExternalPrices | 200 (creado, sin cambio, actualizado, rechazados) | 422 |
| importExternalRecipes | TestOp_importExternalRecipes | 200 (borrador `is_active = 0`) | 422 (101 recetas) |

`testutil.Client.Do` valida cada respuesta (código declarado y esquema) contra `openapi.yaml` con kin-openapi. La credencial de servicio de las pruebas es falsa (`fakeService`, protocolo §7). También están `TestComprasBRepositorios` y `TestComprasBRepositoriosErrores` en `backend/internal/dominio/compras-b_repo_test.go`.

## Verificación

- `go build` y `go vet` de los 26 paquetes sin `internal/api/comprasa`: sin errores. `go build ./...` completo **falla en `internal/api/comprasa`** (otro agente a medio trabajo: «h.getShoppingList undefined»). Por eso `make verify-wf3-go` sale con Error 1 en este momento. No es de este alcance.
- `gofmt -l` de mis archivos: vacío. `gofmt -l .` muestra archivos de comidas, comprasa, creditos y presupuesto (otros agentes). `go mod tidy -diff`: sin cambios.
- `go test ./internal/api/comprasb/... ./internal/dominio/... -count=1 -cover`: `ok comprasb 83.1 %`, `ok dominio`. Las 23 `TestOp_*` pasan (`go test -v`).
- `compras-b_repo.go`: 87.2 % de sus sentencias cubiertas (177 de 203).
- **S6-03**: `TestRutasIgualesAlContrato` da «216 operaciones en 172 rutas, 0 diferencias». `build.py --check` da «los 23 archivos coinciden byte a byte». Hay una prueba por operación (tabla de arriba).
- **S6-04**: los errores salen con el esquema único (`code` en inglés y `message` es-MX, por `api.Problem` y `dominio.Error`). La paginación usa `api.ParsePage` (de 1 a 200; cursor con huella). El dinero va en centavos `int64` y las fechas en ISO. `Idempotency-Key` la exige el middleware: la prueba sin clave da 400.
- **S6-15**: identificadores en inglés, mensajes en español de México, sin tokens ni contraseñas en el código ni en los mensajes. El bearer de prueba es una cadena falsa.
- `make verify-all GOOSE=…`: **falla en `verify-wf1-cobertura`**: `internal/dominio` tiene 72.4 % (< 85 %). Medido archivo por archivo, la caída viene de `comidas-salud_datos.go` (0 %), `compras-a_repo.go` (0 %) y `presupuesto_datos.go` (69.8 %), todos de otros agentes. `compras-b_repo.go` queda en 87.2 %.

## Supuestos (todos PROVISIONAL, con bloqueo)

- El paquete se llama `comprasb` (CONVENCIONES §2 manda sobre el texto del encargo).
- `site_assisted` es el origen de los precios del agente. Recetas del agente con la opción f2: `household`, sin autor, `is_active = 0`. `audit_log` del agente sin usuario (`S5-agente-externo`).
- `getPriceSummary` respeta el motor aprobado: `site_assisted` no compite, aunque C1 1.4 diga PROVISIONAL que compite como manual. Es la pregunta 7 de `S6-compras-b-fuentes`.
- Sitios con `verdict = missing`. `profeco_last_import_at` = alta más reciente de un precio PROFECO. RAR responde 415. Ligas en `<DataDir>/profeco/vinculos.csv`. Corte = AAAA-MM. Sin `audit_log` por renglón en PROFECO (`S6-compras-b-fuentes`).
- `reloadNutritionBase` responde 409. Un preferido repetido responde 409. Pasar a ritmo sin conteo responde 422. Órdenes permitidos según el bloqueo (`S6-compras-b-nutricion`).
- `default_base_servings` sale del DEFAULT del esquema (`pragma_table_info`), no de una constante. `recipe_counts` trae siempre los dos orígenes de C1, con 0 si no hay recetas.

## Bloqueos

- `docs/estado/bloqueos/S6-compras-b-fuentes.md` (nuevo, 7 preguntas).
- `docs/estado/bloqueos/S6-compras-b-nutricion.md` (nuevo, 7 preguntas).
- Siguen abiertos `S5-agente-externo` (a, c, e, f, g).

## Pendientes

- `backend/cmd/finanzas/areas.go`: agregar `comprasb.Register,` y el import `finanzas-hogar/internal/api/comprasb` (fuera de mi alcance; va en solicitudes).
- `reloadNutritionBase` 2xx: depende de la pregunta 1 de `S6-compras-b-nutricion`.
- Cuando WF-3 complete comprasa, comidas y presupuesto, correr otra vez `make verify-wf3` y `make verify-all`.

## Propuestas para Axel

- Fuentes 1 (a), 2 (b), 3 (a), 4 (a), 5 (a), 6 (a).
- Nutrición 1 (b): un `product.code` estable, porque liga la semilla, las bases nutricionales y el agente externo sin depender del nombre. 2 (a), 3 (a), 4 (a).

## Paquetes instalados

Ninguno.

## Consumo aproximado

Unas 52 acciones; unos 260 KB leídos; unos 125 KB escritos (código 107 KB, pruebas incluidas; bloqueos y reporte 18 KB).
