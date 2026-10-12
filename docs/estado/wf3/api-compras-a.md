# api-compras-a · Reporte WF-3a

Agente **api-compras-a** (Opus) · 12-oct-2026 · Alcance: `compras-tiendas-a.yaml` (15) + `createStore`, `updateStore`, `setStoreProductGroups` = **18 operaciones**.

## Hecho
- `backend/internal/api/comprasa/` (paquete de CONVENCIONES §2) con `Register` y las 18 operaciones de C1 v1.4: esquemas de `types_gen.go`, errores con el esquema único (400 formato, 422 regla, 404, 409 `already_exists` / `assisted_reading_not_allowed` / `conflict`), paginación 1–200 con `api.ParsePage`/`PageOf`, Idempotency-Key en las 6 POST que la piden (middleware), `audit_log` en cada alta, cambio o baja, mensajes en español de México.
- Cálculo solo con motores: S5 `precios.Book` (`QuoteFor`, `CurrentPrice`, `BaseUnitPrice`, `Freshness`, `NextVisit`, `CompareBasket`, `CompareExpense`, `Membership`) y S4 `recetas.Engine.ShoppingListFor` (piezas, `assigned_store_id`, `recommended_store_id`, `recommended_presentation_id`, precio estimado, frescura, motivos y recetas; S4 1.5 §6.5 y Q2-39 a Q2-42 A). Orden de grupos y artículos según C1 (mercado primero, visita, id de tienda, sin tienda al final; categoría y nombre).
- `backend/internal/dominio/compras-a_repo.go`: repositorios bajo `dominio.ShopA` (tiendas y grupos, productos, presentaciones, precios, nutrición, lista y artículos, recetas, plan, personas, inventario, compras con renglones). Prueba propia `compras-a_test.go`.

## Archivos
`backend/internal/api/comprasa/{comprasa.go,tiendas.go,productos.go,motor.go,lista.go,comprasa_test.go}`, `backend/internal/dominio/{compras-a_repo.go,compras-a_test.go}`, `docs/estado/bloqueos/S6-compras-a-motores.md`, este reporte.

## Verificación
- `go build ./...` y `go vet ./...`: sin errores. `gofmt -l .`: vacío.
- `go test ./internal/api/comprasa/... -count=1 -cover`: 18/18 PASS, **cobertura 77.5 %**. `go test ./internal/dominio/... -count=1 -cover`: ok, **87.7 %** (≥ 85 %); funciones de `compras-a_repo.go` 92.7 % en promedio.
- `make verify-wf3-go`: build, vet, gofmt, `go mod tidy -diff` sin cambios, `go test ./...` 21 paquetes OK.
- **S6-03**: 18 de 18 operaciones con prueba `TestOp_<operationId>` que valida contra `openapi.yaml` (testutil) con 2xx y 4xx:

| Operación | Prueba | 2xx | 4xx probados |
|---|---|---|---|
| getShoppingList | TestOp_getShoppingList | 200 (grupos, filtros, semana vacía) | 422, 400 |
| generateShoppingList | TestOp_generateShoppingList | 200 ×2 | 400 sin clave, 400 semana |
| createShoppingListItem | TestOp_createShoppingListItem | 201 producto y texto libre | 422 ×2, 400 |
| updateShoppingListItem | TestOp_updateShoppingListItem | 200 | 422 ×2, 400, 404 |
| assignShoppingListStore | TestOp_assignShoppingListStore | 200 | 404, 422 |
| listBaskets | TestOp_listBaskets | 200 | 422, 400 |
| getBasketComparison | TestOp_getBasketComparison | 200 (lista) | 422 ×3, 400 |
| listStores | TestOp_listStores | 200 (q sin acentos, cursor) | 400, 401 |
| createStore | TestOp_createStore | 201 | 422, 400 ×2 |
| updateStore | TestOp_updateStore | 200 | 409, 422, 404 |
| setStoreProductGroups | TestOp_setStoreProductGroups | 204 | 404, 422 |
| listProductGroups | TestOp_listProductGroups | 200 | 400 |
| listProducts | TestOp_listProducts | 200 | 400 ×2 |
| createProduct | TestOp_createProduct | 201 | 409, 422 |
| getProduct | TestOp_getProduct | 200 (best_price S5) | 404, 400 |
| updateProduct | TestOp_updateProduct | 200 | 422, 409, 404 |
| createPresentation | TestOp_createPresentation | 201 | 409, 422, 404 |
| updatePresentation | TestOp_updatePresentation | 200 | 422, 409, 404 |

- **S6-04**: error único de C1 (`dominio.*`/`api.Problem`), página ≤ 200 (`limit=500` → 400), fechas ISO, dinero en centavos `int64`, Idempotency-Key por middleware (sin clave → 400, probado).
- **S6-15**: identificadores en inglés (`ShopA`, `handlers`…), mensajes al usuario en español de México; no se registran credenciales.

## Supuestos
- Paquete en `internal/api/comprasa` (CONVENCIONES §2 manda sobre «compras-a/» del encargo).
- Los motores usan códigos de texto: el código es el id en decimal.
- Orden por defecto `name` en `listStores` y `sort_order` en `listProductGroups` (C1 no lo fija); `listBaskets` pone la lista de la semana primero (solo si existe) y luego las compras por fecha descendente, paginadas en memoria.
- `q` sin acentos ni mayúsculas se filtra en Go (SQLite no quita acentos).
- Lo de `docs/estado/bloqueos/S6-compras-a-motores.md` (10 puntos PROVISIONALES).

## Bloqueos
- `S6-compras-a-motores` (nuevo): motor S4 sin `choice_store` (1.5), quién tiene plan, texto libre, ganadores de la comparativa de la lista, carrito de compra, membresía, filtro `source`.

## Pendientes
- Línea en `backend/cmd/finanzas/areas.go` (fuera de alcance): import `"finanzas-hogar/internal/api/comprasa"` y `comprasa.Register,` en `areas`.
- Prueba 200 de `getBasketComparison` con `basket=expense` (necesita una compra con renglones; depende del alta de gastos de api-gastos-inicio).
- Cobertura del paquete 77.5 % (sin umbral para `internal/api`): faltan ramas de error de la base.

## Propuestas para Axel
- Decidir el bloqueo S6-compras-a-motores (sobre todo 1 y 7): la opción a) en los dos deja la comparativa y las elecciones exactas como dice S4 1.5.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~54 acciones, ~200 KB leídos, ~125 KB escritos.
