# Revisión C1 v1.1 · ronda 4 (flujo P1b, 11-oct-2026)

Revisor: `revisor-C1-r4` (independiente; no escribió nada de lo revisado y no corrige). Evalúa solo contra `docs/criterios/C1.md` v1.1, las decisiones §K (K1–K7, K9, K11–K13, K16, K18, K21, K23) y los hallazgos R2-06 y R3-01, con el foco en los 4 hallazgos altos y medios de la ronda 3 (H-r3-01 a H-r3-04). Base: HEAD `b900420` más el árbol de trabajo (17 archivos modificados, 82 inserciones y 43 borrados según `git diff HEAD --stat`). Todas las verificaciones se volvieron a correr con comandos y scripts propios en `datos/tmp/rev-c1-r4/` (`neg.py`, `chk.py`, `chk2.py`; ya borrados).

## Veredicto

**APROBADO.** No hay hallazgos altos ni medios. Los 4 hallazgos altos y medios de la ronda 3 están cerrados, y las verificaciones (a) a (m) pasan con mis propios comandos. Los 5 hallazgos bajos de r3 siguen abiertos, como estaba previsto, y quedan para WF-1. Hay 3 observaciones bajas nuevas (N-r4-01 a N-r4-03), que no bloquean.

## Verificaciones reejecutadas

| # | Comando (propio) | Resultado |
|---|---|---|
| a | `python3 -I backend/internal/api/c1gen/build.py --check` | IDÉNTICO en los 21 archivos: `openapi.yaml` (632,336 B), `routes.go` (25,004 B), `C1-datos.md` (130,064 B), `extractos/INDICE.md` y los 17 YAML. «RESULTADO: los 21 archivos coinciden byte a byte con el generador», código 0 |
| b | goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/rev-c1-r4/c1.db up` + `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db …` | 00001 OK y 00002 OK, versión 2. **26 OK y 0 FALLA, «RESULTADO: TODO OK»**, código 0. Reglas nuevas: «C1-06 operaciones que leen o escriben shopping_list: 8 · sin `week` o con `period`: [] · ListWeek = week obligatoria: True» y «C1-15 regla del domingo (U1-lista-vigente-domingo abierto: True) · en openapi sin marca: [] · en C1-datos.md sin bloqueo: []» |
| b′ | **Prueba negativa** con `neg.py`: arma tres raíces temporales con enlaces al repositorio y un `openapi.yaml` alterado, y corre `verify_c1.py --root` | (1) Sin la línea `$ref …/ListWeek` de `getShoppingList` (línea 4438): código 1, «FALLA C1-06 las 8 operaciones que leen shopping_list exigen la llave `week`», detecta `['getShoppingList']`. (2) Sin la de `getBasketComparison` (línea 4742): código 1, FALLA C1-06 con `['getBasketComparison']`. (3) Sin `x-provisional: true` de `getMealPlanWeek` (línea 3763): código 1, «FALLA C1-15 la regla del domingo…» con `['/paths//api/v1/meal-plan/week/get']`. En cada caso es la única falla |
| c | `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml docs/contratos/extractos/*.yaml` | `openapi.yaml: OK` y los 17 extractos OK, código 0 |
| d | `python3 -I docs/design/tools/verificar_u1.py` | «RESULTADO: TODO OK», código 0. Las marcas PROVISIONAL de las líneas 66 y 1152 citan `U1-lista-vigente-domingo` (abierto). Cruce openapi ↔ C1-datos: [] |
| e | `cd backend && go build ./... && go vet ./... && gofmt -l .` (go1.27.2) | build OK, vet OK y gofmt sin archivos |
| f | `chk.py` sobre `extractos/` | 17 archivos. El mayor es `gastos.yaml` con 76,202 B; todos ≤ 80,000 B (el tope de INDICE.md) y ≤ 81,920 B. Suma de operationId en los extractos = 201 = `grep -c operationId openapi.yaml` = 201 operaciones. 0 repetidas, 0 faltantes y 0 sobrantes. En los extractos, `getBasketComparison` (`compras-tiendas-a.yaml`) cita `ListWeek` y trae 422, y `getMealPlanWeek` (`comidas.yaml`) trae `x-provisional: true` |
| g | `chk.py` (parámetros, descripción y respuestas de `getBasketComparison`), `grep` en C1-datos §1 y §6, y U1 línea 66 | Ver H-r3-01 abajo: cerrado |
| h | `chk.py` y `chk2.py` (`getMealPlanWeek`, §13 y §10) | Ver H-r3-02 abajo: cerrado |
| i | `git diff HEAD -- docs/design/U1-arquitectura.md docs/criterios/U1.md`, escaneo propio de `` `period` ``/`` `week` `` en U1, `make verify-criterios` (dentro de `verify-all`) y (d) | Ver H-r3-03 abajo: cerrado |
| j | `grep -c` en el diccionario, `gen_dict.py datos/tmp/rev-c1-r4/dic.md` + `cmp`, `diff` sin comentarios de las dos migraciones contra HEAD, `s1_verifica.py` con `GOOSE=…` y `grep assert gen2_schemas.py` | Ver H-r3-04 abajo: cerrado |
| k | `make verify-all GOOSE=…/gobin/goose` | Código 0, 0 líneas con FALLA. «Verificación ampliada OK (WF-0, C1, S1, U1 y criterios)». Criterios: 18 archivos y 257 criterios, «cifras comparadas: 12 \| diferencias: 0», «N-xxx sin operacion ni Local: 0 []», «matriz C1 (seccion 10): filas = 227 \| N-xxx distintas = 203» y «las 10 comprobaciones dan lo esperado» |
| l | `chk2.py`: cada línea PROVISIONAL y cada `x-provisional` de `openapi.yaml`, contra el estado de los archivos de bloqueo | 7 líneas PROVISIONAL y 5 objetos con `x-provisional` (más la mención de la extensión en `info`). Todos citan un bloqueo **abierto**: `C1-valores-k31`, `U1-lista-vigente-domingo`, `criterios-promedios-redondeo` y `S1-ajuste-objetivo-kcal`. Ninguno tiene «Respuesta de Axel» ni «Estado: RESUELTO». En C1-datos hay 19 líneas: las 6 que no citan un bloqueo abierto son títulos, encabezados o texto explicativo (líneas 12, 345, 923, 925 y 936) y la fila histórica de Cambios v1.0 con `C1-lista-compras` (resuelto), igual que en r3. Ninguna marca cae sobre valores de K6 o K7 |
| m | `chk2.py`: compara los números de las líneas agregadas por `git diff -U0 HEAD` (contrato, diccionario, U1, criterios, migraciones, extractos y generador) con los de las líneas quitadas, sin contar identificadores | Números nuevos: `4` («•••» y sus últimos 4 caracteres, que ya estaba en `TelegramChatCandidate.chat_id_masked` en HEAD; `git show HEAD:… \| grep -c` = 1), `422` (código HTTP de §3), `10` (cobertura PROVISIONAL de §11, antes 9, ahora con N-120), `7` y `80` (código de `verify_c1.py`), tamaños de extractos en INDICE.md (generados) y `b900420` (hash). **0 cifras de dinero, tasas o plazos nuevas; 0 inventadas** |

## Criterios de C1.md v1.1

| Criterio | Estado | Evidencia propia |
|---|---|---|
| C1-01 (alta) | Cumple | (c). `openapi` 3.1.1, 162 rutas y 201 operaciones |
| C1-02 (alta) | Cumple | (b), (d) y (k). 203 N-xxx, 0 sin operación ni Local, y la matriz §10 tiene 227 filas. N-120 queda «PROVISIONAL `U1-lista-vigente-domingo`» en §10 (línea 437) |
| C1-03 (alta) | Cumple | (b): 55 tablas, 0 tablas sin operación y 0 inexistentes (también lo dice `gen6_c1.py` en (a)). Las migraciones solo cambiaron comentarios (j) |
| C1-04 (alta) | Cumple | (b) sin fallas. El 422 de `getBasketComparison` va a `#/components/responses/UnprocessableEntity` (esquema `Error`) |
| C1-05 (alta) | Cumple | (b) sin fallas. Esta ronda no cambió la paginación |
| C1-06 (alta) | **Cumple** | (b): 6 de 6 filtros con `$ref` único y 0 sinónimos. 14 componentes de parámetro, uno por nombre. 2 llaves documentadas en §6 (`BudgetMonth → month`, `ListWeek → week`). Las 8 operaciones cuyo `x-tablas` tiene `shopping_list` llevan `ListWeek` (`week`, `required: true`) y ninguna lleva `Period` (`chk.py`). `Period` tiene 30 usos y `ListWeek` 8. La prueba negativa (b′) confirma que la regla nueva falla si una de las 8 pierde `week`. Matriz §6 (línea 188): `getBasketComparison` con «periodo» vacío y `week` ✓ |
| C1-07 (alta) | Cumple | (b) sin fallas. Esta ronda no tocó fechas |
| C1-08 (alta) | Cumple | (b) sin fallas. Esta ronda no tocó dinero |
| C1-09 (alta) | Cumple | (b) y (k), cifras con 0 diferencias. Sigue la observación baja H-r3-05 sobre M-06 |
| C1-10 (alta) | Cumple | (a): el generador no reporta entidades ni tablas sin operación |
| C1-11 (alta) | Cumple | (a): «capacidades con id malo: []» |
| C1-12 (alta) | Cumple | (b) sin fallas. `chat_id` sigue oculto: el comentario nuevo de la columna no llega a `openapi.yaml` y `chat_id_masked` es `x-calculado`. Sigue la observación baja H8 |
| C1-13 (media) | Cumple | 0 operationId duplicados (f). `info.version` 1.1.0. Cambios tiene la fila «C1 v1.1 · corrección r4» (línea 24) |
| C1-14 (media) | Cumple | (b): ejemplos válidos. Las descripciones nuevas de `getBasketComparison`, `basket` y `expense_id` están en español y son concretas. Sigue la observación baja H6 |
| C1-15 (alta) | **Cumple** | (l): toda marca cita un bloqueo abierto. `getMealPlanWeek` lleva `x-provisional: true` y la marca PROVISIONAL con `U1-lista-vigente-domingo`, y está en §13 (línea 930) y en §10 (N-120). La prueba negativa (b′) confirma que la regla nueva falla sin la marca. 0 marcas sobre K6 y K7. Sigue la observación baja H7 |
| C1-16 (alta) | **Cumple** | `scope` sin cambios desde r3 (enum `this_month`/`onward`, `default: onward`). El diccionario registra los dos valores (línea 530, `grep -c` = 1 y 1). `Budget`, `BudgetCreate` y `BudgetUpdate.valid_to_date` traen el mismo texto, que nombra los dos valores y la eliminación de los renglones posteriores (`chk`: (True, True, True) en los tres) |
| C1-17 (alta) | Cumple | Esta ronda no lo tocó. (b) sin fallas |
| C1-18 (alta) | **Cumple** | El diccionario fija `chat_id_masked` en la fila `chat_id` de `telegram_chat_candidate` (línea 2038; «•••» y sus últimos 4 caracteres; K4, C1-18). Coincide con `TelegramChatCandidate.chat_id_masked` de `openapi.yaml` (línea 17317) |
| C1-19 (alta) | Cumple | Esta ronda no lo tocó. (b) sin fallas |
| C1-20 (alta) | Cumple | Esta ronda no lo tocó. (b) sin fallas |
| C1-21 (alta) | Cumple | (a) y (f): regenerables byte a byte, ≤ 80 KB y cada operación una sola vez |

## Decisiones §K, R2-06 y R3-01

| Ref. | Estado | Evidencia |
|---|---|---|
| K1 | OK | Las operaciones de la lista no tienen `x-provisional`. `C1-lista-compras` está resuelto y solo aparece en la historia |
| K2 | OK | C1-16 |
| K3 | OK | Sin cambios. (b) |
| K4 | OK | C1-18 |
| K5 | OK | (b): `family` con enum de 3 y sin `x-provisional` |
| K6 / K7 | OK | (l): 0 marcas sobre sus valores |
| K9, K11, K12, K13, K16, K18, K21, K23 | OK | Sin cambios en esta ronda. (b) y (k) sin fallas. Siguen las observaciones bajas H-r3-06 (K18) y H8 (K6) |
| R2-06 | OK | Sin cambios. (b) C1-07 |
| R3-01 | **OK** | Las 8 operaciones que leen o escriben la lista exigen `week`, y U1 §0.2 y U1-16 la nombran igual (H-r3-01 y H-r3-03) |

## Estado de los 4 hallazgos altos y medios de r3

| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| **H-r3-01** [alta] C1-06, R3-01: `getBasketComparison` sin `week` | **Cerrado** | Parámetros: `basket` (inline), `expense_id` (inline) y `$ref ListWeek` → `week`, `in: query`, `required: true`, patrón ISO de semana. No recibe `Period`. La descripción ya no es «Comparar canasta.»: explica las dos modalidades y dice «Si falta `week`, o falta `expense_id` con `basket=expense`, responde 422 `validation_failed` con `fields`». Las respuestas incluyen `422` → `UnprocessableEntity`. C1-datos: la fila de la matriz §6 tiene «periodo» vacío y `week` ✓. La tabla «Llaves que no son filtros» (línea 149) y §1 (línea 40) la incluyen. Es coherente con U1 §0.2: «Toda lectura y escritura de la lista lleva la semana (`week` en C1, obligatoria)». `verify_c1.py` tiene la regla nueva y la prueba negativa la hace fallar. Que `week` sea obligatoria también con `basket=expense` está documentado como consecuencia de C1-06 (un componente por nombre), y no cambia cifras ni rutas |
| **H-r3-02** [alta] C1-15: regla del domingo sin marca en `getMealPlanWeek` | **Cerrado** | La operación tiene `x-provisional: true` (línea 3763). La descripción dice «la parte de los domingos es PROVISIONAL: bloqueo U1-lista-vigente-domingo (abierto)». §13 tiene la fila «Semana por omisión del plan de comidas» con ese bloqueo (línea 930). En §10, N-120 está marcada PROVISIONAL (línea 437). El bloqueo sigue abierto: no tiene «Respuesta de Axel». `verify_c1.py` C1-15 tiene la regla nueva y la prueba negativa la hace fallar. Las otras menciones de «domingo» en `openapi.yaml` son K8 («lunes a domingo»), ISO («7 = domingo») o el día de la semana («Domingo.»), y ninguna aplica la regla |
| **H-r3-03** [media] U1 §0.2 llamaba `period` a la semana de la lista | **Cerrado** | `U1-arquitectura.md` línea 66: «Toda lectura y escritura de la lista lleva la semana (`week` en C1, obligatoria)». Ahí solo queda un `` `period` ``, el del plan de la semana (`getMealPlanWeek`), que es correcto. `criterios/U1.md` U1-16 (líneas 80–81) dice «`week` en C1, obligatoria» en el enunciado y en Verifica, con el ajuste anotado en la lista de cambios (línea 12). `make verify-criterios` (dentro de `verify-all`) y `verificar_u1.py` pasan |
| **H-r3-04** [media] diccionario sin `this_month`, `onward` ni `chat_id_masked` | **Cerrado** | `grep -c` en el diccionario: `this_month` 1, `onward` 1, `chat_id_masked` 1 y `valid_to_date` 3. La línea 530 dice «deja NULL y elimina los renglones posteriores de la misma categoría…». `gen_dict.py` regenera idéntico (`cmp` sin diferencias, «columnas sin descripción: []»). Quitando los comentarios, las dos migraciones son iguales a HEAD (`diff` vacío). En el diff `-U0`, las líneas agregadas y quitadas son iguales sin comentarios (0 únicas). `s1_verifica.py`: «RESUMEN TODO OK» (NEG 46/46, NEG12 26/26, DOWN igual a 1.1 con 232 objetos, GOOSE [0, 0, 0]). En `gen2_schemas.py` el parche con `assert` ya no existe: la única aparición de «assert» es el comentario de la línea 112 |

## Estado de los 5 hallazgos bajos de r3 (no bloquean; quedan para WF-1)

| Hallazgo | Estado | Evidencia |
|---|---|---|
| H-r3-05 (C1-09, M-06 «lista vigente») | Sigue | `C1-datos.md` línea 330: M-06 dice «lista de compras vigente» y su periodo «lista vigente», aunque `getDashboardBestBasket` exige `week` |
| H-r3-06 (K18, IVA) | Sigue | `Credit.vat_on_interest_bp`: «NULL = la tasa o la cotización ya incluye IVA». `missing_terms` no distingue si falta el dato del IVA |
| H6 (C1-14) | Sigue | Muchas propiedades (`id`, `created_at`, `updated_at`, `notes` e `is_active` son las más repetidas) tienen descripciones que solo repiten el identificador o van vacías |
| H7 (C1-15) | Sigue | `DashboardWeekMeals` (línea 11925) y `NotificationUpdate` (línea 17168) tienen la marca PROVISIONAL con `C1-valores-k31`, pero no `x-provisional` |
| H8 (C1-12, K6) | Sigue | `LoginRequest.username` = `{type: string, description}`, sin `minLength` 3 ni `maxLength` 64 |

## Hallazgos nuevos

| Id | Severidad | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| N-r4-01 | baja | C1-04, R3-01 (coherencia con U1-16) | El 422 por falta de `week` solo está escrito en `getBasketComparison`. El componente `ListWeek` y las otras 7 operaciones no dicen qué código responde si falta, y `generateShoppingList` (POST, la escritura «generar la lista») no tiene `422` entre sus respuestas (200, 400, 401, 409 y default). En cambio, la Verifica de U1-16, que cambió en esta ronda, espera «una escritura de la lista sin `week` se rechaza (422)». No es una contradicción: §3 asigna a 422 la «validación de campos», y `default` cubre el código con `Error`. Pero un implementador de WF-3 podría responder 400 | `backend/internal/api/c1gen/gen4_build.py` línea 19 (componente `ListWeek`), `gen3_ops.py` línea 145 (`generateShoppingList`), `docs/criterios/U1.md` línea 81 | Agregar a la descripción única de `ListWeek` «si falta, 422 `validation_failed` con `fields`» y agregar `422` a las respuestas de `generateShoppingList`. Otra opción es que U1-16 diga «se rechaza (4xx)» |
| N-r4-02 | baja | C1-14, R3-01 (misma raíz que H-r3-05) | Quedan textos que hablan de una lista «vigente» que el contrato ya no define: `Basket.basket` (`openapi.yaml` línea 13238, «Lista vigente o compra registrada.») y `DashboardBestBasket` (línea 11971, «lista de compras vigente»). Además, el esquema de `basket` en `getBasketComparison` dice «(por defecto)» pero no declara `default: shopping_list` (viene de HEAD) | `backend/internal/api/c1gen/gen2_schemas.py` líneas 211 (Basket) y 153 (DashboardBestBasket), `gen3_ops.py` línea 154 | Decir «la lista de la semana `week`». Agregar `default: shopping_list` al esquema de `basket` |
| N-r4-03 | baja | C1-09, C1-15 (precisión) | `getDashboardWeekMeals` (P-01, tarjeta «Comidas de la semana») recibe `Period` sin decir qué semana toma si no se manda, y su descripción es solo «Tarjeta «Comidas de la semana».». No aplica ninguna regla sin marca, así que no hay valor inventado. Pero, a diferencia de `getMealPlanWeek`, no remite a la semana por defecto de U1 §0.2, y U1 §0.2 solo nombra P-10, P-13 y P-14 | `backend/internal/api/c1gen/gen3_ops.py` línea 42 (`getDashboardWeekMeals`), `docs/design/U1-arquitectura.md` §0.2 | En WF-1, decir qué semana usa la tarjeta sin `period`. Si es la de U1 §0.2, marcarla con `U1-lista-vigente-domingo` como `getMealPlanWeek`. Si no, hacer `period` obligatorio en esa operación |
