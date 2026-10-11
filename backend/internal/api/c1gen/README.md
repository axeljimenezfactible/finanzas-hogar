# c1gen · Generador y verificador del contrato C1

Esta carpeta es la **fuente** del contrato C1. De aquí salen tres archivos que **no se editan a mano**:

| Salida | Etapa que la escribe | Qué es |
|---|---|---|
| `docs/contratos/openapi.yaml` | `gen4_build.py` (con `gen1_base.py`, `gen2_schemas.py` y `gen3_ops.py`) | OpenAPI 3.1.1, contrato C1 |
| `backend/internal/api/routes.go` | `gen5_go.py` | Rutas del enrutador `net/http`, idénticas a las del contrato (S6-03) |
| `docs/contratos/C1-datos.md` | `gen6_c1.py` (texto de `c1_plantilla.md` y `c1_inconsistencias.md`) | Documento de C1, con la matriz §10 pantalla → necesidad → endpoint |

Entradas que se leen del repositorio: el esquema de S1 (`backend/db/migrations/00001_esquema.sql`) y las necesidades de U1 (`docs/design/U1-arquitectura.md`).

## Uso

Desde la raíz del repositorio (da igual la carpeta; las rutas se resuelven desde esta carpeta):

```sh
python3 -I backend/internal/api/c1gen/build.py            # regenera los tres archivos en el repositorio
python3 -I backend/internal/api/c1gen/build.py --check    # regenera en una carpeta temporal y compara byte a byte; no escribe
python3 -I backend/internal/api/c1gen/build.py --out DIR  # escribe los tres archivos bajo DIR (árbol espejo)
cd backend && go generate ./internal/api                  # lo mismo que la primera línea (generate.go)
```

`--check` sale con código 0 si los tres archivos coinciden con el generador y 1 si alguno difiere (muestra las primeras líneas de la diferencia). Los intermedios `schema.json` y `openapi.json` van a una carpeta temporal bajo `datos/tmp/` que se borra al terminar. Solo usa la biblioteca estándar de Python 3.11 o posterior.

Orden de las etapas (lo aplica `build.py`): `parse_schema.py` (SQL → `schema.json`) → `gen4_build.py` (→ `openapi.yaml` y `openapi.json`) → `gen5_go.py` (→ `routes.go`) → `gen6_c1.py` (→ `C1-datos.md`). `c1paths.py` resuelve las rutas; ninguna etapa tiene rutas absolutas fijas.

## Cómo se cambia el contrato

1. Edita la fuente, nunca las salidas: esquemas en `gen2_schemas.py`, operaciones (ruta, `operationId`, etiqueta, `x-necesidades`, `x-tablas`) en `gen3_ops.py`, parámetros, encabezados, respuestas y armado del documento en `gen4_build.py`, tablas generadas de C1 en `gen6_c1.py`, texto de C1 en `c1_plantilla.md` y la sección 14 en `c1_inconsistencias.md`.
2. Si cambia la versión del contrato, actualiza `info.version` (en `gen4_build.py`) y la tabla «Cambios» de `c1_plantilla.md` (C1 v1.1, v1.2…, plan §3.7).
3. Regenera con `build.py`, corre `build.py --check` y la verificación de abajo.

## Verificación

```sh
uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py
# C1-03 contra una base migrada con goose (se abre solo lectura):
goose -dir backend/db/migrations sqlite3 datos/tmp/c1.db up
uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/c1.db
uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml   # C1-01
cd backend && go build ./... && go vet ./... && gofmt -l .
```

`verify_c1.py` revisa C1-01 a C1-15, la consistencia de la matriz §10 (pares pantalla-necesidad de U1 §3.29 y pares necesidad-operación de `x-necesidades`) y la paridad de rutas `routes.go` = `openapi.yaml` (S6-03). Sale con código 1 si hay alguna falla.

C1-02 sigue el texto vigente del criterio: solo acepta como excepción las necesidades que U1 declara «Local» (las lee de U1; hoy N-003 y N-301) y la operación `GET /openapi.yaml` (C1 §11.2). La pantalla P-26 se cubre con sus secciones P-26.xx. Una necesidad bloqueada **no cuenta como cubierta**. Desde la corrección r2 (esquema 1.1 de S1) ninguna está bloqueada: N-151, N-152, N-153 y N-165 tienen operación sobre `shopping_list_item`, PROVISIONAL por la opción A de `C1-lista-compras`. C1-03 exige que toda tabla figure en algún `x-tablas`: `idempotency_key` la agrega `gen4_build.py` a cada operación con `Idempotency-Key`.

## Procedencia

El agente S6a-contrato-c1 (WF-0) escribió estos archivos en su carpeta temporal de sesión. La corrección r1 (hallazgo H-05 de `docs/estado/wf0/REVISION-r1.md`) los trajo al repositorio. Antes de cualquier cambio se comprobó que regeneran byte a byte el `openapi.yaml`, el `routes.go` y el `C1-datos.md` revisados en la ronda 1. Cambios hechos al traerlos:

- Rutas: `parse_schema.py`, `gen1_base.py`, `gen4_build.py`, `gen5_go.py` y `gen6_c1.py` toman sus rutas de `c1paths.py`. La lógica no cambió y `gen2_schemas.py`, `gen3_ops.py` y `c1_inconsistencias.md` siguen idénticos.
- Texto: las cabeceras de `openapi.yaml` (comentario) y de `routes.go` (comentario «Code generated … DO NOT EDIT.») ahora citan esta carpeta. `c1_plantilla.md` suma la fila «Generador» y una fila en «Cambios». No cambió ninguna ruta, esquema ni fila de la matriz.
- `verify_c1.py`: rutas relativas, opciones `--db` y `--root`, consistencia de la matriz §10 y código de salida 1 si hay fallas.
- Nuevos: `build.py`, `c1paths.py`, este `README.md` y `../generate.go`.

Corrección r2 (hallazgos R2-01 y R2-02 de `docs/estado/wf0/REVISION-r2.md`): regeneración contra el esquema 1.1 de S1 (sha256 `747fc8fd…`).

- `gen2_schemas.py`: lista de compras guardada (`ShoppingListItem`, `ShoppingListItemCreate`, `ShoppingListItemUpdate`, `ShoppingListStoreAssignment`), presupuesto con vigencia, `meal_slot_ids` de receta sobre `recipe_meal_slot`, descripción del gasto sin relleno, `InventoryMovement` y `InventoryAdjustmentRequest` sobre `inventory_movement`, tope por quincena sin cifra.
- `gen3_ops.py`: `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore` y `getShoppingList` con la lista guardada; tablas nuevas en `x-tablas`; textos de presupuesto, ingreso esperado e inventario.
- `gen4_build.py`: `idempotency_key` en `x-tablas` de toda operación con `Idempotency-Key`; código 409 `budget_overlap`.
- `gen1_base.py`: la marca PROVISIONAL de una columna cita el bloqueo que trae el comentario de S1 (antes citaba siempre `S1-meta-kcal`) y se conserva en los campos que admiten null.
- `gen6_c1.py`, `c1_plantilla.md` y `c1_inconsistencias.md`: estados de la matriz §10, conteos de tablas, §11.4, §13 y §14 al día con el esquema 1.1.
