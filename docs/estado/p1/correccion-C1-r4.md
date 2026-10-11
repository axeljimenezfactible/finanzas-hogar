# Corrección C1 v1.1 · ronda 4 (flujo P1b, 11-oct-2026)

Agente: `correccion-C1-r4` (Opus). Alcance: los 4 hallazgos altos y medios de `docs/estado/p1/REVISION-C1-r3.md` (H-r3-01 a H-r3-04). Los bajos (H-r3-05, H-r3-06, H6, H7 y H8) no se tocaron y quedan para WF-1. Sin `git commit`.

## Hecho
1. **H-r3-01 [alta] C1-06, R3-01 · `getBasketComparison`**
   - Recibe `ListWeek` (`week`, obligatoria) y ya no recibe `Period`.
   - La descripción ya no es «Comparar canasta.». Ahora dice:
     - con `basket=shopping_list` (por omisión) compara la lista de la semana `week`;
     - con `basket=expense` compara la compra `expense_id`, y `week` no cambia el resultado;
     - si falta `week`, o falta `expense_id` con `basket=expense`, responde 422 `validation_failed` con `fields`.
   - `basket` y `expense_id` también tienen descripciones precisas.
   - En la matriz §6, «periodo» queda vacío y `week` lleva ✓. La operación entra en la tabla «Llaves que no son filtros» y en §1.
   - `verify_c1.py` C1-06 tiene una regla nueva: toda operación cuyo `x-tablas` incluye `shopping_list` o `shopping_list_item` exige `$ref` a `ListWeek`, y `ListWeek` debe llamarse `week` y ser obligatoria. Además, esas operaciones no pueden recibir `Period`. Hoy son 8 y las 8 cumplen.
2. **H-r3-02 [alta] C1-15 · `getMealPlanWeek`**
   - La semana por omisión lleva «PROVISIONAL: bloqueo U1-lista-vigente-domingo (abierto)» en la descripción y `x-provisional: true` en la operación. `gen4_build.py` pone esa marca solo, porque la descripción contiene PROVISIONAL.
   - §13 de C1-datos.md tiene una fila nueva, en `c1_plantilla.md`.
   - N-120 pasa a «PROVISIONAL `U1-lista-vigente-domingo`» en §10 (`gen6_c1.py`, `PARTIAL`).
   - `verify_c1.py` C1-15 tiene una regla nueva de 14 líneas. Mientras el bloqueo siga abierto:
     - en `openapi.yaml`, falla toda descripción con «los domingos» que no tenga `x-provisional` (en su objeto o en su operación), la marca PROVISIONAL y el nombre del bloqueo;
     - en `C1-datos.md`, falla toda línea con «los domingos» que no cite el bloqueo.
   - La regla busca «domingos» y no «domingo», porque «lunes a domingo» (K8) y «7 = domingo» (ISO) son hechos firmes y no la regla.
3. **H-r3-03 [media] C1-02, R3-01 · nombre de la semana en U1**
   - `docs/design/U1-arquitectura.md` §0.2 (línea 66) dice «lleva la semana (`week` en C1, obligatoria) que la pantalla muestra».
   - La frase sobre el valor por omisión ahora se refiere al filtro `period` del plan de la semana (`getMealPlanWeek`).
   - `docs/criterios/U1.md` U1-16 (líneas 80 y 81) dice «la semana … (`week` en C1, obligatoria)», en el enunciado y en «Verifica».
   - La lista de cambios v1.1 de criterios/U1.md corrige la línea de U1-16 y suma una línea que anota el ajuste. La versión sigue siendo 1.1.
4. **H-r3-04 [media] C1-16, C1-18 · diccionario**
   - Se aplicaron las solicitudes 1 y 2 de `correccion-C1-r2.md`, cambiando solo comentarios de columna:
     - `00001_esquema.sql`, `budget.valid_to_date`: nombra `scope` de updateBudget, `onward` (preseleccionada y valor por omisión) y `this_month`, y dice que `onward` «elimina los renglones posteriores de la misma categoría, que el monto nuevo reemplaza (updateBudget, C1-datos.md §7.1)»;
     - `00002_puerta1.sql`, `telegram_chat_candidate.chat_id`: «La API nunca lo devuelve completo: lo muestra como el campo calculado chat_id_masked («•••» y sus últimos 4 caracteres; K4, C1-18)».
   - El diccionario se regeneró con `gen_dict.py`.
   - Se quitó de `gen2_schemas.py` el parche con `assert`. Ahora `Budget`, `BudgetCreate` y `BudgetUpdate` toman el mismo texto del comentario de S1.
   - **Dónde quedó `chat_id_masked`:** en el comentario de la columna `chat_id`. Cabe ahí y `gen_dict.py` lo publica en la fila `chat_id` de la tabla `telegram_chat_candidate`. `gen_dict.py` no genera una sección de campos calculados. Como `chat_id` está en `HIDDEN` de `gen1_base.py`, ese comentario no aparece en `openapi.yaml`.
5. Se agregó la fila «C1 v1.1 · corrección r4» en Cambios, el autor en la cabecera de `c1_plantilla.md` y la sección «C1 v1.1 · corrección r4» en `backend/internal/api/c1gen/README.md`.

## Archivos
- Fuente C1: `backend/internal/api/c1gen/gen3_ops.py`, `gen2_schemas.py`, `gen6_c1.py`, `c1_plantilla.md`, `verify_c1.py` y `README.md`.
- Regenerados con `build.py`: `docs/contratos/openapi.yaml`, `docs/contratos/C1-datos.md`, `docs/contratos/extractos/INDICE.md`, `comidas.yaml`, `compras-tiendas-a.yaml` y `presupuesto-a.yaml`. `backend/internal/api/routes.go` se regeneró sin cambios, porque las rutas no cambian.
- Solo comentarios: `backend/db/migrations/00001_esquema.sql` (1 línea) y `00002_puerta1.sql` (1 línea).
- Regenerado con `gen_dict.py`: `docs/contratos/diccionario-datos.md` (2 renglones).
- Parches mínimos: `docs/design/U1-arquitectura.md` (línea 66) y `docs/criterios/U1.md` (U1-16 y su lista de cambios).
- `docs/estado/p1/correccion-C1-r4.md` (este reporte).

## Verificación
- `python3 -I backend/internal/api/c1gen/build.py`: escribió openapi, routes, C1-datos y los 17 extractos más INDICE.
- `python3 -I backend/internal/api/c1gen/build.py --check`: «RESULTADO: los 21 archivos coinciden byte a byte con el generador» (IDÉNTICO en todos).
- `goose v3.28.0 -dir backend/db/migrations sqlite3 datos/tmp/c1.db up`: 00001 OK, 00002 OK, versión 2. Después, `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/c1.db` dio código de salida 0 y **RESULTADO: TODO OK**. Las líneas de las reglas nuevas:
  - «C1-06 operaciones que leen o escriben shopping_list: 8 · sin `week` o con `period`: [] · ListWeek = week obligatoria: True» → OK;
  - «C1-15 regla del domingo (U1-lista-vigente-domingo abierto: True) · en openapi sin marca: [] · en C1-datos.md sin bloqueo: []» → OK;
  - C1-15 general: 7 marcas en openapi (antes 4 `x-provisional` y 6 líneas PROVISIONAL), con 4 bloqueos abiertos citados (incluye `U1-lista-vigente-domingo`).
- **Prueba negativa:** el `verify_c1.py` nuevo, corrido con `--root` sobre `git archive HEAD` (b900420), da «FALLA C1-06 … ['getBasketComparison']» y «FALLA C1-15 regla del domingo … ['/paths//api/v1/meal-plan/week/get']» (2 fallas). Las reglas detectan lo que deben.
- `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml docs/contratos/extractos/*.yaml`: openapi.yaml OK y los 17 extractos OK.
- `python3 -I docs/design/tools/verificar_u1.py`: «RESULTADO: TODO OK». El cruce openapi ↔ C1-datos da [].
- `cd backend && go build ./... && go vet ./... && gofmt -l .`: build OK, vet OK, gofmt vacío.
- Extractos: 17 archivos. La suma de operationId es 201 y `grep -c operationId openapi.yaml` también da 201. El mayor es gastos.yaml con 76,202 B; todos ≤ 81,920 B.
- `GOOSE=… python3 -I backend/db/tools/s1_verifica.py`: «RESUMEN TODO OK». DOWN: el sqlite_master de la 00002 es igual al esquema 1.1 (232 objetos, 0 diferencias). GOOSE: up / down-to 0 / up dan [0, 0, 0].
- Migraciones:
  - en `git diff -U0 HEAD -- backend/db/migrations/`, las líneas quitadas y las agregadas son iguales una vez quitados los comentarios (`sed 's/--.*//'`): diff vacío;
  - `diff <(git show HEAD:… | sed 's/--.*//') <(sed 's/--.*//' …)` también da vacío en los dos archivos;
  - el DDL no cambió ni un byte.
- `python3 -I backend/db/tools/gen_dict.py` y otra corrida a `datos/tmp/dic2.md`: `cmp` idéntico, «columnas sin descripción: []».
- `grep -c` en el diccionario: this_month 1, onward 1, chat_id_masked 1.
- `x-provisional`: `openapi.yaml` líneas 3758 a 3763, `operationId: "getMealPlanWeek"` … `x-provisional: true`.
- Matriz §6: `| getBasketComparison | GET /api/v1/basket-comparison |  |  |  |  |  |  | ✓ |` (periodo vacío, `week` ✓).
- `make verify-all GOOSE=…`: código 0 y «Verificación ampliada OK (WF-0, C1, S1, U1 y criterios)». Los criterios dan 0 N-xxx sin operación, matriz de 227 filas, cifras con 0 diferencias.
- `datos/tmp` se borró al terminar.

## Supuestos
- **`week` obligatoria en las dos modalidades de `getBasketComparison`.** La corrección pedía `week` obligatoria con `basket=shopping_list`. En OpenAPI, un parámetro por `$ref` no puede cambiar `required`. Además, C1-06 de `verify_c1.py` prohíbe `week` inline (está en la lista de sinónimos) y un segundo componente con el mismo nombre. Por eso `week` es obligatoria siempre, y con `basket=expense` no cambia la canasta. P-14 siempre muestra una semana (U1 §0.2), así que la PWA siempre la tiene. No cambia cifras, esquema ni rutas.
- El 422 `validation_failed` por falta de `week` se escribió solo en la descripción de `getBasketComparison`, como pedía la corrección. El componente `ListWeek` no se tocó.
- Se quitó `Period` de `getBasketComparison` en vez de definir qué filtra con `basket=expense`: con `expense_id` la canasta ya queda determinada, y cualquier significado de `period` ahí sería un valor inventado.

## Bloqueos
- No hay bloqueos nuevos. Se cita el bloqueo existente y abierto `docs/estado/bloqueos/U1-lista-vigente-domingo.md`.

## Pendientes
- No queda nada pendiente de los 4 hallazgos. Los bajos H-r3-05 (M-06 «lista vigente»), H-r3-06 (IVA de `vat_on_interest_bp`, solicitud a S1), H6, H7 y H8 siguen registrados en `REVISION-C1-r3.md` para WF-1.

## Solicitudes fuera de alcance
- Ninguna. Los cambios en U1, criterios, migraciones y diccionario estaban en el alcance asignado.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Unas 40 acciones de herramienta.
- Leído: unos 120 KB (lecturas parciales con grep y sed; ningún archivo grande se leyó entero).
- Escrito por mí: unos 22 KB (scripts de parche, filas de plantilla, README y este reporte). Las salidas regeneradas las escribieron `build.py` y `gen_dict.py`.
