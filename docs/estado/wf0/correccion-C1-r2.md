# Reporte WF-0 · correccion-C1-r2

Corrección del artefacto C1 según la revisión independiente r2 (`docs/estado/wf0/REVISION-r2.md`) · flujo WF-0 «Cimientos» · 11-oct-2026 · Estado: **terminado**, sin cierre por tope.
Hallazgos asignados: **R2-01** (alta; C1-02, U1-03, C1-11, U1-05) y **R2-02** (alta; C1-03, C1-15, U1-12). Los resolví en una sola regeneración de C1 contra el esquema 1.1 de S1: migración sha256 `747fc8fd88e576f4`, U1 `45fcfb61da0ca09f`.

## Hecho

1. **R2-01 · lista de compras.** Agregué operaciones sobre `shopping_list` y `shopping_list_item` (`gen2_schemas.py` y `gen3_ops.py`). Todas quedan PROVISIONALES por la opción A de `C1-lista-compras`, la misma pregunta que `S1-lista-compras`, hasta la puerta 1.
   - `POST /api/v1/shopping-list/generate` · `generateShoppingList` · N-127. Crea la lista de la semana si falta, pone `generated_at` y guarda `needed_milli` por producto. Conserva lo elegido por una persona (invariante 13 de S1).
   - `POST /api/v1/shopping-list/items` · `createShoppingListItem` · N-152. Agrega un artículo manual (producto o texto libre con unidad). Crea la lista de la semana si no existe.
   - `PATCH /api/v1/shopping-list/items/{shopping_list_item_id}` · `updateShoppingListItem` · N-151 y N-153. Marca comprado (`is_bought`; la API pone o quita `bought_at`) y cambia tienda, presentación o cantidad.
   - `POST /api/v1/shopping-list/assign-store` · `assignShoppingListStore` · N-165 (paso 3 de la ruta A de F-3). Pone `store_id` en los artículos que tienen precio en esa tienda. Responde 404 si la semana no tiene lista.
   - `GET /api/v1/shopping-list` · `getShoppingList` · N-150 y N-073. Devuelve la lista guardada más lo que calculan S4 y S5, y tiene el filtro nuevo `is_bought` para el prellenado de «Terminar compra».
   - Esquemas nuevos: `ShoppingListItem` (entidad con campos calculados marcados), `ShoppingListItemCreate`, `ShoppingListItemUpdate`, `ShoppingListStoreAssignment` y `...Result`. `ShoppingList` usa `start_date` y `end_date`, los nombres de las columnas.
   - `getDashboardBestBasket`, `listBaskets` y `getBasketComparison` declaran las dos tablas.
2. **R2-02 · C1 alineado con el esquema 1.1.**
   - (1) Presupuesto con vigencia: `Budget` y `BudgetCreate` llevan `valid_from_date` y `valid_to_date`, sin `period_id`. `BudgetUpdate` permite cerrar o reabrir la vigencia. Un traslape responde el código nuevo 409 `budget_overlap` (invariante 7 de S1). M-01 suma el renglón vigente en el mes (RN-22).
   - (2) Tiempos de receta: `meal_slot_ids` con `x-columna` `recipe_meal_slot.meal_slot_id` en `Recipe`, `RecipeCreate` y `RecipeUpdate`. `AiRecipeRequest.meal_slot_id`. `recipe_meal_slot` aparece en 7 operaciones. El tablero de recetas se agrupa por tiempo apto.
   - (3) Descripción del gasto opcional y sin relleno (S1 C-16). Quité la marca PROVISIONAL y el texto «la API usa el nombre de la categoría». `Expense` ya no la exige.
   - (4) `InventoryMovement` es la tabla `inventory_movement` (más el campo calculado `expense_id`). `InventoryAdjustmentRequest` escribe `reason` y `quantity_delta_milli`, o la cantidad nueva. Escriben su movimiento `adjustInventoryItem`, `createInventoryItem` y `confirmTicketScan`.
   - (5) `PlannedItem` trae `second_half_amount_cents`. Las descripciones de N-279 y N-280 y `IncomeMonth.expected_cents` lo usan.
   - (6) `idempotency_key`: `gen4_build.py` la agrega a `x-tablas` de las 41 operaciones con `Idempotency-Key`. C1 §1 punto 8, §11.1 y §11.4 la justifican como infraestructura de S6; la retención sigue PROVISIONAL.
   - (7) Tope por quincena sin cifra: la descripción es neutral en `Person` y `PersonUpdate`. Ya no dice «Karen 37493 en el oráculo» y tampoco copia la interpretación de S1 que R2-03 cuestiona.
   - `category.kind` entra solo, desde el esquema.
3. **Marcas PROVISIONAL resueltas, retiradas.** Salen de la matriz §10 N-054, N-122, N-135, N-136, N-181, N-183, N-275, N-279 y N-280.
   - De los huecos U1 ↔ S1 quedan PROVISIONALES solo los 3 de `S1-huecos-negocio`: N-087 y N-278 (qué vigencia escribe un cambio), N-075 y N-299.
   - En `C1-datos.md` actualicé §1 punto 8, §3 (409 `budget_overlap`), §10 (leyenda), §11.1, §11.3, §11.4 (con 54 de 54 y el párrafo de justificación), §13 y §14 (tabla con columna «Estado»).
   - En «Cambios» agregué la fila «C1 v1.0 · corrección r2».
4. **Bloqueos de C1 coherentes con los de S1.**
   - `C1-lista-compras`: la misma pregunta y las mismas opciones que `S1-lista-compras`, y «Cómo seguí» con la opción A ya implementada.
   - `C1-huecos-u1-s1`: solo los 3 puntos de `S1-huecos-negocio`.
   - `C1-parametros-tecnicos` punto 2: ya no pide la tabla, que existe; solo pregunta la retención. Quité la opción C, que guardaba las claves en memoria.
5. **El generador y el verificador.**
   - `gen1_base.py`: la marca PROVISIONAL de una columna cita el bloqueo que trae su comentario de S1. Antes citaba siempre `S1-meta-kcal`; por ejemplo, `Period.half` ahora cita `U1-periodos-semana-quincena`. Además, la marca se conserva en los campos que admiten null.
   - `gen6_c1.py`: estados y conteos dinámicos, ENT y VIS, capacidades, y una regex de ids sin falsos positivos («lista», «market»).
   - `verify_c1.py`: en C1-03 la cuenta de propiedades sin columna ahora **sí falla** (antes solo se imprimía), y cada `x-columna` se compara con la base. Lo comprobé con una prueba negativa.
   - Actualicé el `README.md`.
6. Agregué notas «(corrección r2)» en `S6a-contrato-c1.md`: una nota general, C1-02, las tablas de cobertura, Bloqueos, Pendientes 2 y la solicitud 2.

## Archivos

Modificados (generador, fuente):
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/gen1_base.py` (`c5875d1e08ce6445`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/gen2_schemas.py` (`ae1709f661fdd28c`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/gen3_ops.py` (`621fb8ba781667c4`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/gen4_build.py` (`47ca0b028c84e873`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/gen6_c1.py` (`efb13a7565de0c29`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/c1_plantilla.md`
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/c1_inconsistencias.md` (`df2d83330beaf9bc`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/verify_c1.py` (`658da67129f10c0e`)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/README.md`

Regenerados con `build.py` (no se editan a mano):
- `/home/deepsec/finanzas-hogar/docs/contratos/openapi.yaml`: 597527 B, `95c8866589d1`. 158 rutas, 197 operaciones, 237 esquemas.
- `/home/deepsec/finanzas-hogar/backend/internal/api/routes.go`: 24507 B, `b4b4d33296eb`, 197 rutas.
- `/home/deepsec/finanzas-hogar/docs/contratos/C1-datos.md`: 113439 B, `e3ec3c882ee4`.

Bloqueos de C1 actualizados:
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-lista-compras.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-huecos-u1-s1.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-parametros-tecnicos.md`

Reportes:
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/S6a-contrato-c1.md`: notas «(corrección r2)».
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/correccion-C1-r2.md`: este reporte.

Sin cambios: `backend/go.mod` (go 1.26), `backend/cmd/finanzas/main.go`, `generate.go` y `build.py`. No hay `go.sum`, porque solo se usa la biblioteca estándar.

## Verificación

- **Regeneración.** `python3 -I backend/internal/api/c1gen/build.py` y luego `build.py --check`, corrido 3 veces: `IDÉNTICO` ×3 (`openapi.yaml`, `routes.go`, `C1-datos.md`), código 0.
  - Salida del generador: `N sin endpoint: ['N-003', 'N-301']`, que son Local. También `ops sin N: ['getOpenApiDocument']`, `tablas sin op: []`, `tablas inexistentes: []` y `capacidades con id malo: []`.
  - Antes de corregir daba `DIFIERE openapi.yaml (570515 → 573729)`, `DIFIERE C1-datos.md` y `tablas sin op: [idempotency_key, recipe_meal_slot, inventory_movement, shopping_list, shopping_list_item]`.
- **SQL en base temporal.** `goose v3.28.0 -dir backend/db/migrations sqlite3 datos/tmp/correccion-C1-r2/c1.db up/down/up`: `OK 00001_esquema.sql`; tras el down quedan 0 tablas; tras el segundo up hay 54 tablas, `integrity_check` = ok y `foreign_key_check` = 0 filas.
  - Ida y vuelta de las escrituras del contrato sobre una copia de la base: lista de la semana con un artículo manual en texto libre, marca de comprado con `bought_at`, `generated_at`, receta con 2 tiempos en `recipe_meal_slot`, sobrante con su movimiento `recipe_leftover`, ajuste `adjustment`/`consumption` de −500 y presupuesto con vigencia (monto de prueba). Resultado: `listas=1 articulos=1 tiempos=2 movimientos=2 budget=1`, `fk_rows=0`.
  - Negativas rechazadas por la base: ajuste sin motivo (CHECK), artículo manual sin cantidad (CHECK), vigencia que no empieza el día 1 (CHECK) y tiempo repetido en una receta (UNIQUE).
  - Al final borré `datos/tmp` completa; no existía antes.
- **C1-01 a C1-15, matriz y S6-03.** `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/correccion-C1-r2/c1.db` → código 0, **22 OK, 0 FALLA, `RESULTADO: TODO OK`**.
  - C1-02: `N-xxx en U1: 199 · con operación en la matriz §10: 197 · Local: ['N-003', 'N-301'] · N-xxx sin operación ni Local: 0 [] (con estado BLOQUEO en §10: [])`. También 0 pantallas sin operación, matriz de 220 filas = 220 pares de U1 y 0 diferencias con `x-necesidades`.
  - C1-03: `tablas en sqlite_master: 54 · citadas por x-tablas: 54 · tablas sin operación: []`, `propiedades de esquemas de entidad: 967 · sin columna …: 0 []`, `x-columna comprobadas contra la base: 8 · que no existen: []`.
  - C1-04: 895 de 895. C1-05: 39 de 39. C1-06: 6 de 6. C1-07: solo `date` y `date-time`. C1-08: 256 `_cents`, 0 `number`. C1-12: públicas solo salud y login. C1-13: 0 ids duplicados, 9 etiquetas, «C1 v1.0», 1.0.0. C1-14: 13 ejemplos válidos y 2672 descripciones en español. C1-15: 73 marcas en el YAML y 42 líneas en C1, todas con un bloqueo existente; 20 bloqueos citados, todos existen.
  - S6-03: 197 rutas = 197.
- **Prueba negativa del verificador** (árbol espejo en el scratchpad con `Budget.period_id` reinsertado y una `x-columna` que no existe): `FALLA C1-03 … ['Budget.period_id']`, `['Recipe.meal_slot_ids -> recipe.default_meal_slot_id']`, código 1.
- **Cruces independientes.**
  - `matriz_r2.py` del revisor (`scratchpad/rev-r2/`): `N-xxx sin operación ni Local: 0 []`, 41 pantallas sin faltantes, `C1 §10: filas = 220 · pares U1 §3.29 = 220`, `filas BLOQUEO: []`. Estados: OK 198, Local 2, PROVISIONAL 20 (7 `C1-lista-compras`, 4 `C1-huecos-u1-s1`, más los de login, familia, umbrales, reposición y optimizador).
  - `u1_verify.py` de U1: `OK U1-03 N-xxx sin endpoint en C1: []`. Su única falla es `GET /openapi.yaml` sin N-xxx, la excepción «estáticos» que C1-02 admite (§11.2) y que ese script no reconoce.
- **Validadores.** `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` → `OK`, código 0. `npx -y @redocly/cli@2.62.1 lint` → «Your API description is valid», 1 advertencia `info-license` (la misma de r1), código 0. Al cargar el YAML: `3.1.1 158 197 237 1.0.0`.
- **Esquemas contra cargas de ejemplo** (jsonschema, Draft 2020-12): 18 casos. Son válidos los 12 previstos, como artículo con producto o con texto libre, presupuesto con vigencia, ajuste con cantidad nueva o con diferencia, receta con `meal_slot_ids`, gasto sin descripción y partida quincenal con dos montos. Se rechazan los 6 previstos: sin cantidad, con `origin`, con `period_id`, ajuste sin motivo, receta con `default_meal_slot_id` y PATCH vacío. Además, `ExpenseCreate.description` ya no lleva x-provisional y no es requerida en `Expense`.
- **Textos contrarios a S1 1.1.** En `openapi.yaml` hay 0 apariciones de `37493`, «nombre de la categoría», «S1 lo exige», `period_id`, `default_meal_slot_id` o «derivado». En `C1-datos.md`, `period_id` y `default_meal_slot_id` solo aparecen en la fila «Cambios», que dice que se quitaron.
- **Cifras del oráculo** en C1 §9 sin cambios: 10,000.00, 6,905.50, 973.75, 1,947.50, 598.82, 374.93 y 749.86 están en `C1-datos.md` y en `oraculo/model.json`. No agregué ninguna cifra nueva al contrato.
- **Go** (go1.26.5), en `backend/`: `go build ./...` código 0, `go vet ./...` código 0, `gofmt -l .` sin salida, `go list ./...` da 2 paquetes, `go generate -n ./internal/api` da `python3 -I c1gen/build.py`, y no hay `__pycache__` ni `.pyc`.
  - Binario en 127.0.0.1:18094: `GET /health` 200; las 5 rutas de la lista, los movimientos y `GET /api/v1/budgets` responden 501 `not_implemented`; `DELETE /api/v1/shopping-list/items/7` 405; una ruta inexistente 404; `/openapi.yaml` 200, 597527 B. Detuve el servidor y el puerto quedó libre.
- **No aplica a C1:** `colorimetria.py` es de U2a.

## Supuestos

1. El contrato sigue llamándose **C1 v1.0**, con una fila «corrección r2» en «Cambios». Sigue siendo la entrega de WF-0 para la puerta 1: C1-13 exige «C1 v1.0» e `info.version` 1.0.0, y el plan §3.7 versiona como v1.1 los cambios posteriores.
2. La forma de la API de la lista de compras está cubierta por el bloqueo `C1-lista-compras` y marcada PROVISIONAL:
   - La semana se nombra con el filtro único `period` (semana ISO; PROVISIONAL `U1-periodos-semana-quincena`).
   - Si la semana no tiene lista, la consulta responde la lista vacía.
   - Agregar un artículo crea la lista de la semana.
   - «Usar esta tienda» responde 404 si no hay lista.
   - No hay operación para quitar un artículo, porque U1 no tiene esa necesidad.
3. El código 409 `budget_overlap` es convención propia de C1 (tabla única de códigos, §3) para la invariante 7 de S1. No cambia cifras ni esquema.
4. `createInventoryItem` registra la existencia inicial de un producto como movimiento `adjustment` con motivo `count` (conteo inicial). La invariante 12 pide un movimiento por cada cambio, pero no nombra el tipo para la existencia inicial. No cambia esquema ni rutas.
5. Declarar `idempotency_key` en `x-tablas` de toda operación con `Idempotency-Key` describe lo que ya dice S1 C-15: la API lee y escribe la clave en la misma transacción.

## Bloqueos

No hay bloqueos nuevos. Actualicé tres para que no contradigan al esquema 1.1:
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-lista-compras.md`: la misma pregunta de `S1-lista-compras`.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-huecos-u1-s1.md`: los 3 puntos de `S1-huecos-negocio`.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-parametros-tecnicos.md`: solo la retención de las claves.

Siguen abiertos sin cambios `C1-filtro-familia`, `C1-campos-calculados` y `S6-version-go`.

## Pendientes

1. **R2-06** (baja, C1; no venía en esta asignación).
   - El patrón de `DateTime` acepta «Z», pero C1 §4 dice que los instantes llevan el desplazamiento de America/Mexico_City. Se resuelve documentando que «Z» solo se acepta en entradas.
   - `rate_pct` con 4 decimales mostraría 17.2725 % a 2 quincenas, frente al 17.27 % de la decisión C. Para esta parte hay que decidir cuántos decimales se muestran, porque la misma decisión C usa 25.8685 % a 4 quincenas: si se muestran 2 decimales, esa tasa se vería como 25.87 %. Conviene preguntarlo a Axel antes de tocarlo.
2. Go 1.27: `backend/go.mod` sigue en `go 1.26` (bloqueo `S6-version-go`).
3. Cuando Axel responda en la puerta 1 `C1-lista-compras` (A, B o C) y los 3 puntos de `C1-huecos-u1-s1`, hay que regenerar C1 con `gen3_ops.py` y `gen2_schemas.py`, y quitar las marcas PROVISIONAL.

## Solicitudes fuera de alcance

1. **U1** (`docs/estado/wf0/U1-arquitectura.md` y su `u1_verify.py`):
   - El reporte de U1 aún dice «U1-03: PENDIENTE». Hoy `u1_verify.py` da `N-xxx sin endpoint en C1: []`.
   - Su segunda comprobación de U1-03 debería admitir `GET /openapi.yaml` como excepción «estáticos», como lo hace C1-02.
2. **S1** (R2-03, `backend/db/migrations/00001_esquema.sql` y diccionario, `person.fortnight_payment_cap_cents`): la descripción de la columna sigue afirmando que «el oráculo no trae un tope». C1 ya no copia ese texto (usa uno neutral, sin cifra) y no cambia con la respuesta.
3. **U1** (opcional): no existe la necesidad «Quitar artículo» de la lista. Si U1 la agrega, C1 agrega `DELETE /api/v1/shopping-list/items/{shopping_list_item_id}`.
4. **Criterios** (`docs/criterios/C1.md` C1-02 y C1-03, `docs/criterios/S6.md` S6-03), parte de R2-07: citar `backend/internal/api/c1gen/verify_c1.py --db` y `build.py --check`. Desde esta ronda, el C1-03 de ese script falla también por propiedades sin columna y por `x-columna` inexistentes.

## Paquetes instalados

Ninguno en el sistema. Usé lo que ya estaba en cachés de usuario:
- `uvx --offline --from openapi-spec-validator==0.9.0`
- `npx -y @redocly/cli@2.62.1`
- `uv run --no-project --with pyyaml --with jsonschema`, en un entorno efímero
- `goose` v3.28.0 de `scratchpad/gobin/`

## Consumo aproximado

- Unas 100 acciones de herramienta, sin subagentes.
- Leí unos 250 KB, en lecturas dirigidas: protocolo, decisiones, REVISION-r2, criterios C1, reporte de r1, bloqueos, secciones del diccionario, RN-15 y RN-22, y extractos de U1 y del generador.
- Edité unos 30 KB de fuente y regeneré unos 735 KB (YAML, Go y C1).
- El contador de la sesión marca unos 370 K tokens acumulados, por encima de los 300 K de referencia. El trabajo quedó completo y consistente, sin compresión de contexto.
