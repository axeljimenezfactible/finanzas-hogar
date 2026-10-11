# Corrección C1 v1.1 · ronda 2 (flujo P1, 11-oct-2026)

Agente: `correccion-C1-r2`. Corrige los hallazgos alto y medios de `docs/estado/p1/REVISION-C1-r2.md` (H-r2-01, H-r2-02 y H-r2-03). No toca los bajos (H6, H7, H8) ni reabre decisiones de §A–§K.

## Hecho

1. **H-r2-01 (alta, C1-06, R3-01): opción (a) del revisor.** La semana de la lista de compras ya no se llama `period`. Es la llave `week` (`components/parameters/ListWeek`: obligatoria, semana ISO `^[0-9]{4}-W(0[1-9]|[1-4][0-9]|5[0-3])$`, sin valor por omisión). La usan `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets` y `getDashboardBestBasket`.
   - El filtro `period` vuelve a tener **un solo componente**, `Period`, con sus 5 formatos y 31 usos.
   - Se eligió (a) y no (b) porque (b) devuelve un valor por omisión que depende del bloqueo abierto `U1-lista-vigente-domingo`, lo que obligaría a nuevas marcas PROVISIONAL. También desharía el cierre de R3-01 en r1 (semana obligatoria). Con (a) el contrato sigue sin semana por omisión.
   - §6 de `C1-datos.md` tiene ahora la tabla **«Llaves que no son filtros»**, con `week` (ListWeek, las 7 operaciones) y `month` (BudgetMonth, `updateBudget`). La fila «periodo» dice que la lista no usa ese filtro.
   - La matriz recurso × filtro tiene una columna nueva, «semana de lista (`week`, llave)». Entran `listBaskets`, `getDashboardBestBasket` y las 4 escrituras. `getShoppingList` lleva marca en tienda y en `week`.
   - Las descripciones de `getDashboardBestBasket`, `updateShoppingListItem` (422 `out_of_range`), `listBaskets` y el texto común `SEMANA` dicen `week`.
   - `verify_c1.py` C1-06 tiene una regla nueva que exige tres cosas: cada nombre de parámetro tiene un solo componente; ningún otro componente usa el nombre de un filtro; toda llave con nombre de periodo (`month`, `week`…) aparece en §6.
2. **H-r2-02 (media, C1-16, K2): `onward` con renglones posteriores.** Con `scope=onward`, el monto nuevo rige desde `month` **sin fin** (`valid_to_date` NULL). Los renglones posteriores de la misma categoría (o del total) **se eliminan**, porque el monto nuevo los reemplaza. Quedan en `audit_log`, que `updateBudget` ya declara en `x-tablas`. Con `this_month`, los renglones posteriores no cambian.
   - Por qué así: es lo que dice K2 («de este mes en adelante»: el monto rige este mes y los siguientes) y lo que dice el comentario de S1 en `budget.valid_to_date` («deja NULL», «NULL = rige hasta que se cambie»). Es la única forma de cumplir las dos cosas sin traslape (invariante 7). La otra lectura, respetar los renglones posteriores, deja enero en 1,947.50 en el ejemplo del revisor y contradice K2. Por eso no se escribió el bloqueo `C1-presupuesto-onward`.
   - Ya no hay textos que se contradigan. `updateBudget` dice lo mismo que §7.1 (tabla de alcances y párrafo). `Budget.valid_to_date` completa el texto de S1 con un parche en `gen2_schemas.py` que falla con `assert` si S1 cambia el comentario. `BUDVIG` lleva la misma regla a las descripciones de `BudgetCreate` y `BudgetUpdate`. `createBudget` ya no dice «o hasta el día anterior al siguiente renglón» como lectura de «en adelante»: ahora aclara que el alta solo llena meses sin monto y no elimina renglones.
   - §7.1 tiene un tercer ejemplo, «`scope=onward` con renglones posteriores». Parte del resultado de `this_month` (7: oct–nov, 8: solo dic, 9: desde ene) y aplica `PATCH /api/v1/budgets/7?month=2026-11&scope=onward` `{"amount_cents": 251737}`. Queda 7 (194750, 2026-10-01 a 2026-10-31) más 10 nuevo (251737, desde 2026-11-01, sin fin), y se eliminan el 8 y el 9. La respuesta trae `following_month_amount_cents` 251737. Los montos salen del oráculo: 1,947.50 y 2,517.37 son `periodos[0]` y `periodos[2]` de `devengado.despensa_total`.
3. **H-r2-03 (media, C1-16 y C1-18): no hay nada que cambiar en C1.** El diccionario es de S1. La solicitud se repite abajo (solicitud 1). `grep -c` sigue dando 0 para `this_month`, `onward` y `chat_id_masked`.
4. Se agregó la fila «C1 v1.1 · corrección r2» en Cambios y se actualizó el resumen de v1.1 (línea de la lista de compras). `info.version` sigue en 1.1.0, igual que en la corrección r1, porque v1.1 todavía no está aprobado.
5. Se agregó la sección «C1 v1.1 · corrección r2» en `backend/internal/api/c1gen/README.md`.

## Archivos

Fuente (`backend/internal/api/c1gen/`):
- `gen4_build.py`: `ListWeek`, nombre `week`.
- `gen3_ops.py`: `SEMANA`, descripciones de 3 operaciones de la lista, `updateBudget` y `createBudget`.
- `gen2_schemas.py`: `BUDVIG` y `Budget.valid_to_date`.
- `gen6_c1.py`: matriz con la columna `week`.
- `c1_plantilla.md`: Cambios, resumen, §6 y §7.1.
- `verify_c1.py`: regla C1-06.
- `README.md`.

Regenerados con `build.py`:
- `docs/contratos/openapi.yaml` (630,717 B).
- `docs/contratos/C1-datos.md` (128,721 B).
- `docs/contratos/extractos/` (los 17 YAML y `INDICE.md`).
- `backend/internal/api/routes.go` (25,004 B, sin cambio de contenido).

Reporte: `docs/estado/p1/correccion-C1-r2.md`.

## Verificación

| Comando | Salida resumida |
|---|---|
| `python3 -I backend/internal/api/c1gen/build.py` | 162 rutas, 201 operaciones, 242 esquemas. routes.go: 201 rutas. N sin endpoint: N-003, N-301, N-317. Tablas sin operación: []. 16 extractos más `componentes-comunes.yaml`, el mayor de 76,202 B |
| `python3 -I backend/internal/api/c1gen/build.py --check` | «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| `goose -dir backend/db/migrations sqlite3 datos/tmp/c1r2/c1.db up` (v3.28.0) | 00001 OK, 00002 OK, versión 2 |
| `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/c1r2/c1.db` | **24 OK, 0 FALLA, «RESULTADO: TODO OK»**: los 23 de r2 más la regla C1-06 nueva. C1-06: «nombres con más de un componente: [] · otros componentes con nombre de filtro: [] · llaves que no son filtro: BudgetMonth → month, ListWeek → week · sin documentar en §6: []». C1-14: 14 ejemplos válidos. C1-15: 6 marcas, 3 bloqueos abiertos. EXT y S6-03 OK |
| Prueba negativa de la regla nueva (script propio: `ListWeek` vuelve a llamarse `period`) | Detecta «period: ['Period', 'ListWeek']» y «otros: ['ListWeek']», o sea, sí habría fallado con el estado de r2 |
| `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator` | `openapi.yaml: OK` y los 17 YAML de `extractos/` OK |
| `python3 -I docs/design/tools/verificar_u1.py` | «RESULTADO: TODO OK». Cruce openapi ↔ C1-datos: [] |
| `cd backend && go build ./... && go vet ./... && gofmt -l .` (go1.27.2) | build OK, vet OK, gofmt sin archivos |
| Script propio sobre `openapi.yaml` y `C1-datos.md` | Componentes con `name: period`: solo `Period`. 7 operaciones con `week`. `x-provisional`: 4, sin cambio. operationId: 201. `updateBudget` habla de los posteriores y deja «fin original» solo para `this_month`. Los 2 cuerpos JSON de §7.1 parsean. Cifras de §7.1: 1,947.50 / 194750 / 2,517.37 / 251737, todas del oráculo |

## Supuestos

- Cuando `onward` elimina renglones posteriores, eso queda en `audit_log`. Es lo mismo que cualquier escritura de la API, y `updateBudget` ya declara `audit_log` en `x-tablas`. No cambia cifras, esquema ni rutas.

## Bloqueos

Ninguno nuevo. No se escribió `C1-presupuesto-onward`: K2 y el comentario de S1 alcanzan para decidir (ver Hecho 2).

## Pendientes

- Observación para el revisor y el orquestador, fuera de los hallazgos y sin cambios: `getBasketComparison` (`GET /api/v1/basket-comparison`) usa el filtro común `period` (opcional), y con `basket=shopping_list` no dice de qué semana compara la lista. Si debe usar la llave `week` igual que `listBaskets`, eso cambia una operación que no está en los hallazgos de r2, y por eso no se tocó.
- Los hallazgos bajos H6, H7 y H8 siguen abiertos, porque esta corrección no los incluía.

## Solicitudes fuera de alcance

1. **Ruta**: `docs/contratos/diccionario-datos.md` (dueño S1). **Cambio**: aplicar la solicitud 1 de `docs/estado/p1/correccion-C1-r1.md`. Hay que registrar `scope` de `updateBudget` (`this_month` = «solo este mes»; `onward` = «de este mes en adelante», por omisión y preseleccionado; K2) y el campo calculado `chat_id_masked` de `telegram_chat_candidate` («•••» más los últimos 4 caracteres; el `chat_id` completo nunca se devuelve; K4). **Motivo**: H-r2-03 (antes H5 de r1). C1-16 y C1-18 lo piden y `grep -c` da 0, 0 y 0.
2. **Ruta**: `docs/contratos/diccionario-datos.md` línea 530 y el comentario de `budget.valid_to_date` en `backend/db/migrations/00001_esquema.sql` línea 299 (dueño S1). **Cambio**: después de ««de este mes en adelante» (preseleccionada) deja NULL», agregar «y elimina los renglones posteriores de la misma categoría, que el monto nuevo reemplaza (updateBudget, C1-datos.md §7.1)». **Motivo**: H-r2-02 pide igualar el diccionario. C1 ya completa ese texto en `Budget.valid_to_date`, y el `assert` de `gen2_schemas.py` falla si S1 reescribe el comentario sin el texto «deja NULL;». Si S1 aplica este cambio, hay que quitar ese parche.
3. **Ruta**: `docs/criterios/U1.md`, U1-16, líneas 79–80 (dueño: agente de criterios o integrador). **Cambio**: donde dice «manda el `period` que la pantalla está mostrando… C1 vuelve `period` obligatorio» y «`period` en las operaciones de escritura de la lista es obligatorio», decir «la semana que la pantalla está mostrando (`week` en C1, obligatoria)». **Motivo**: H-r2-01 dejó la semana de la lista en la llave `week` para que el filtro `period` tenga una sola definición (C1-06). La intención de U1-16 se cumple, porque la semana es explícita y obligatoria en las 4 escrituras y en las 3 lecturas, pero el texto todavía nombra `period`.

## Paquetes instalados

Ninguno.

## Consumo aproximado

- Acciones de herramienta: unas 34 de 60.
- Leído: unos 115 KB. Incluye protocolo, §K, el informe r2, extractos de los generadores con grep y sed, salidas de build y de verificación, y un grep de U1 y de criterios. No se leyó ningún archivo grande completo.
- Escrito a mano: unos 24 KB (dos scripts de parche en el scratchpad, el README y este reporte). Los artefactos regenerados los escribió `build.py`.
- Temporales: `datos/tmp/c1r2/`, borrado al terminar.
