# c1gen · Generador y verificador del contrato C1

Esta carpeta es la **fuente** del contrato C1. De aquí salen estos archivos, que **no se editan a mano**:

| Salida | Etapa que la escribe | Qué es |
|---|---|---|
| `docs/contratos/openapi.yaml` | `gen4_build.py` (con `gen1_base.py`, `gen2_schemas.py` y `gen3_ops.py`) | OpenAPI 3.1.1, contrato C1 |
| `backend/internal/api/routes.go` | `gen5_go.py` | Rutas del enrutador `net/http`, idénticas a las del contrato (S6-03) |
| `docs/contratos/C1-datos.md` | `gen6_c1.py` (texto de `c1_plantilla.md` y `c1_inconsistencias.md`) | Documento de C1, con la matriz §10 pantalla → necesidad → endpoint |
| `docs/contratos/extractos/` | `gen7_extractos.py` | Un extracto por etiqueta, `componentes-comunes.yaml` e `INDICE.md` (lo único del contrato que leen los agentes de WF-1 y WF-3) |

Entradas que se leen del repositorio: el esquema de S1 (todas las migraciones de `backend/db/migrations/`, aplicadas en orden) y las necesidades de U1 (`docs/design/U1-arquitectura.md`).

## Uso

Desde la raíz del repositorio (da igual la carpeta; las rutas se resuelven desde esta carpeta):

```sh
python3 -I backend/internal/api/c1gen/build.py            # regenera los tres archivos en el repositorio
python3 -I backend/internal/api/c1gen/build.py --check    # regenera en una carpeta temporal y compara byte a byte; no escribe
python3 -I backend/internal/api/c1gen/build.py --extractos  # regenera en una carpeta temporal y solo escribe docs/contratos/extractos/
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

`verify_c1.py` revisa C1-01 a C1-15 y C1-23, la consistencia de la matriz §10 (pares pantalla-necesidad de U1 §3.29 y pares necesidad-operación de `x-necesidades`) y la paridad de rutas `routes.go` = `openapi.yaml` (S6-03). Sale con código 1 si hay alguna falla.

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

## C1 v1.1 (puerta 1)

Agente C1-v1.1 (flujo P1, 11-oct-2026). Aplica las respuestas de Axel (`docs/plan/decisiones.md` §K) y la sección «Para C1» de `docs/estado/p1/U1-1.1.md` contra el esquema 1.2 de S1. Reporte: `docs/estado/p1/C1-v1.1.md`.

- `parse_schema.py`: aplica **todas** las migraciones de `backend/db/migrations/` en orden (hoy `00001_esquema.sql` y `00002_puerta1.sql`) en una base temporal bajo `datos/tmp/` y lee `sqlite_master`, `PRAGMA table_info` y `PRAGMA foreign_key_list`; los comentarios de columna se buscan en todos los archivos (`CREATE TABLE` y `ALTER TABLE … ADD COLUMN`). Ninguna ruta fija a `00001`; `build.py` y `verify_c1.py` (sin `--db`) hacen lo mismo.
- `gen7_extractos.py` (etapa nueva): `docs/contratos/extractos/`, un archivo por etiqueta (≤ 80 KB; si no cabe, partes `-a`, `-b`…), `componentes-comunes.yaml` e `INDICE.md`. `build.py --extractos` solo escribe esa carpeta; `build.py` la regenera con lo demás y `--check` la compara.
- `gen2_schemas.py`, `gen3_ops.py` y `gen4_build.py`: K2 (`scope` y `month` en `updateBudget`), K3 (`saveTicketScanDraft`, `TicketScanDraft`, borrador en `TicketScan`), K4 (`listTelegramChatCandidates`, `confirmTelegramChatCandidate`; sale `createNotificationRecipient`), K5 (`family`), K9 (`captured_at` y 200 al reintento en `createExpense` y `createWeightEntry`), K12 (`rate_decimals`, `rate_pct`, `annual_rate_pct`), K13 (tope por persona), K16 (restricciones vacías con `missing_constraints` y `status` del comparador), K18 (campos de tarjeta y crédito de tienda, `missing_terms`), K21 (`getProjectionSettings`, `updateProjectionSettings`, `projection_mode`), K23 (`goal`, `suggested_kcal_target`) y R3-01 (`period` de semana obligatorio en la lista de compras). `info.version` = 1.1.0.
- `gen6_c1.py`, `c1_plantilla.md` y `c1_inconsistencias.md`: fila «C1 v1.1 (puerta 1)» en Cambios y su detalle, convenciones de K6 en §1, «Z» solo en entradas en §4, §13 y §14 al día; §10 se regenera con N-316, N-318 y N-319 (con operación) y N-317 (Local).
- `verify_c1.py`: sin `--db` aplica todas las migraciones, exige `info.version` 1.1.0, C1-15 exige que toda marca PROVISIONAL cite un bloqueo **abierto** (sin «Respuesta de Axel») y la comprobación nueva EXT revisa los extractos (cada operación en uno solo, ≤ 80 KB, referencias resueltas).
- Marcas PROVISIONAL que quedan: solo las de los bloqueos abiertos `C1-valores-k31` (valores operativos de K31), `criterios-promedios-redondeo` y `S1-ajuste-objetivo-kcal`.

### C1 v1.1 · corrección r2

Agente correccion-C1-r2 (flujo P1, 11-oct-2026), hallazgos H-r2-01 y H-r2-02 de `docs/estado/p1/REVISION-C1-r2.md`. Reporte: `docs/estado/p1/correccion-C1-r2.md`.

- `gen4_build.py`: `components/parameters/ListWeek` se llama `week` (semana ISO, obligatoria, sin valor por omisión). Es la llave de la lista de compras semanal, no un filtro: el filtro `period` vuelve a tener un solo componente (`Period`).
- `gen3_ops.py`: las descripciones de la lista, `listBaskets` y `getDashboardBestBasket` dicen `week`; `updateBudget` con `scope=onward` deja el monto nuevo sin fin y elimina los renglones posteriores de la categoría; `createBudget` aclara que el alta no elimina renglones.
- `gen2_schemas.py`: `Budget.valid_to_date` completa el comentario de S1 («deja NULL») con la eliminación de los renglones posteriores; `BUDVIG` dice lo mismo en `BudgetCreate` y `BudgetUpdate`.
- `gen6_c1.py` y `c1_plantilla.md`: la matriz recurso × filtro de §6 suma la columna de la llave `week` (entran `listBaskets`, `getDashboardBestBasket` y las escrituras de la lista); §6 trae la tabla «Llaves que no son filtros» (`week` y `month`); §7.1 dice qué hace `onward` con los renglones posteriores y trae el ejemplo; fila «C1 v1.1 · corrección r2» en Cambios.
- `verify_c1.py` C1-06: regla nueva, cada nombre de parámetro tiene un solo componente, ningún otro componente usa el nombre de un filtro y toda llave con nombre de periodo (`month`, `week`…) está documentada en §6.

### C1 v1.1 · corrección r4

Agente correccion-C1-r4 (flujo P1b, 11-oct-2026), hallazgos H-r3-01, H-r3-02, H-r3-03 y H-r3-04 de `docs/estado/p1/REVISION-C1-r3.md`. Reporte: `docs/estado/p1/correccion-C1-r4.md`.

- `gen3_ops.py`: `getBasketComparison` recibe `ListWeek` (`week`, obligatoria) en lugar de `Period` y su descripción dice qué compara con `basket=shopping_list` y con `basket=expense` (422 `validation_failed` si falta `week`, o `expense_id` con `basket=expense`). `getMealPlanWeek` marca su semana por omisión «PROVISIONAL: bloqueo U1-lista-vigente-domingo», y `gen4_build.py` le pone `x-provisional: true`.
- `gen2_schemas.py`: sin el parche con `assert` de `Budget.valid_to_date`. El comentario de S1 en `00001_esquema.sql` ya trae el texto completo y nombra `onward` y `this_month`, así que `Budget`, `BudgetCreate` y `BudgetUpdate` lo toman igual.
- `gen6_c1.py` y `c1_plantilla.md`: N-120 con cobertura PROVISIONAL `U1-lista-vigente-domingo` en §10, fila nueva en §13, `getBasketComparison` en la tabla «Llaves que no son filtros» y en §1, y fila «C1 v1.1 · corrección r4» en Cambios.
- `verify_c1.py`: C1-06 exige `ListWeek` (obligatoria) y prohíbe `Period` en toda operación cuyo `x-tablas` tenga `shopping_list` o `shopping_list_item`. C1-15 marca toda descripción de `openapi.yaml` con «los domingos» que no tenga `x-provisional` (en su objeto o en su operación), la marca PROVISIONAL y el bloqueo, y toda línea de `C1-datos.md` con «los domingos» que no cite el bloqueo, mientras `U1-lista-vigente-domingo` siga abierto. Las dos reglas fallan contra el contrato de HEAD (b900420).

## C1 v1.2 (puerta 1b)

Agente C1-v1.2 (flujo P1c, 11-oct-2026). Aplica K35 a K44 de `docs/plan/decisiones.md` y la sección «Para C1» de `docs/estado/p1c/U1-1.2.md` contra el esquema 1.3 de S1 (`00001`, `00002` y `00003_puerta1b.sql`). Reporte: `docs/estado/p1c/C1-v1.2.md`.

- `parse_schema.py`: sigue aplicando todas las migraciones en orden (hoy tres). Lee además las anotaciones `@param`, `@param_sin_valor` y `@col_note` de cualquier migración (cabecera incluida, sin la sección Down) y escribe `params.json` en la carpeta temporal. Cada valor inicial anotado se compara con la siembra real de `app_setting` en la base temporal; si no coinciden, el generador se detiene. Las notas `@col_note` se suman a la descripción de su columna (hoy `notification.params_json`, con la marca del bloqueo abierto `S1-parametros-avisos`).
- `gen1_base.py`: expone `PARAMS` (el inventario de `params.json`).
- `gen2_schemas.py`: K44 con `SettingParameterInteger` y `SettingParameterIntegerList` (valor vigente, valor inicial, `is_initial`, tipo y unidad). Cada sección de Ajustes con parámetros tiene `parameters` en su lectura y `parameters` y `reset_parameters` en su cambio: «restablecer» es un campo, no una operación. Esquemas nuevos `SecuritySettings`, `DataSettings`, `HealthSettings` y `AlertSettings` (con sus `…Update`); `PriceSettings` y `OptimizerSettings` suman `parameters` (`PriceSettings.stale_after_days` pasa a `parameters`). K36: `HealthProfile` suma `household_goal_adjust_lose_kcal` y `household_goal_adjust_gain_kcal`, y `suggested_kcal_target` cita RN-17 1.3 (las columnas por persona salen solas de la 00003). Sin las marcas de K43 y K39; `NotificationUpdate` con `send_time` y `advance_days` FIRMES y el lugar PROVISIONAL (`S1-parametros-avisos`). `PasswordChangeRequest` sin `minLength` 12 fijo (el mínimo es `auth.password_min_length`).
- `gen3_ops.py`: `getSecuritySettings`/`updateSecuritySettings` (N-320), `getAlertSettings`/`updateAlertSettings` (N-323), `dismissTelegramChatCandidate` (N-324), `getDataSettings`/`updateDataSettings` (N-325), `getHealthSettings`/`updateHealthSettings` (N-326); N-321 en las del comparador, N-322 en las de precios y N-327 en `createHealthProfile` y `updateHealthProfile`; N-289 sale de `updatePerson`. `getMealPlanWeek` FIRME (K38); adjunto 10 MB y comparador 5 s FIRMES con su clave de parámetro (K43, K44); `compareScenarios` declara 503.
- `gen4_build.py`: `info.version` 1.2.0, respuesta 503 `ServiceUnavailable` y textos de K44 (idempotencia y convenciones).
- `gen6_c1.py` y `c1_plantilla.md`: §10 con estado **FIRME (Kxx)** para lo que cerró la puerta 1b; §11.7 nueva, generada de `params.json` y de los esquemas de settings (una lectura y un cambio por clave o el generador se detiene); §13 con solo el bloqueo abierto `S1-parametros-avisos`; fila «C1 v1.2 (puerta 1b)» en Cambios.
- `verify_c1.py`: exige 1.2.0 y «C1 v1.2»; comprobación nueva **K44** (cada `@param` con una sola lectura y un solo cambio con «restablecer», en §11.7 con su valor inicial y sembrado con ese valor en la base); en C1-12, una clave de `app_setting` expuesta como parámetro (`x-app-setting` igual a su nombre, p. ej. `auth.password_min_length`) no cuenta como campo sensible.

### C1 v1.2 · corrección r1

Agente correccion-C1-r1 (flujo P1c, 11-oct-2026). Hallazgo A-01 de `docs/estado/p1c/REVISION-C1-r1.md` (C1-23, K43). Reporte: `docs/estado/p1c/correccion-C1-r1.md`.

- `gen2_schemas.py`: `params_json` de `Notification` y de `NotificationUpdate` declara en `properties` las claves de K43: `send_time` (`type: string`, `pattern: ^([01][0-9]|2[0-3]):[0-5][0-9]$`) y `advance_days` (`type: integer`, `minimum: 0`), cada una con descripción que cita K43 y sin valor por omisión. Sigue `additionalProperties: true` (para `weekday` y `day_of_month`) y la marca de lugar `x-provisional` del bloqueo `S1-parametros-avisos`; el generador se detiene si falta alguna de las dos.
- `gen6_c1.py` y `c1_plantilla.md`: al final de §11.7, la validación de `send_time` y `advance_days` sale de esas propiedades de openapi (el generador exige que `Notification` y `NotificationUpdate` traigan las mismas); fila «C1 v1.2 · corrección r1» en Cambios.
- `verify_c1.py`: comprobación nueva **C1-23**: el patrón y el mínimo en openapi.yaml, en el extracto que trae `NotificationUpdate` y en C1-datos.md, sin `default`, con `additionalProperties: true` y la marca de lugar; el patrón acepta 00:00 a 23:59 y rechaza 24:00, 7:30 y 12:60; además, 10 MB con 413 en `uploadExpenseAttachment`, 503 `unavailable` en `compareScenarios`, N-023, N-064, N-226, N-297 y N-300 FIRMES (K43) en §10, M-05 con el máximo de 21 y 0 líneas `PROVISIONAL.*C1-valores-k31`.
