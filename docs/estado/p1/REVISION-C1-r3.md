# Revisión C1 v1.1 · ronda 3 de 3 (flujo P1, 11-oct-2026)

Revisor: `revisor-C1-r3` (independiente; no escribió nada de lo revisado y no corrige). Evalúa solo contra `docs/criterios/C1.md` v1.1, las decisiones §K (K1–K7, K9, K11–K13, K16, K18, K21, K23) y los hallazgos R2-06 y R3-01. Todas las verificaciones se volvieron a correr con comandos y scripts propios en `datos/tmp/rev-c1-r3/` (ya borrado).

## Veredicto

**NO APROBADO.** Quedan 2 hallazgos altos y 2 medios.
- **Alto 1 (C1-06, R3-01):** `getBasketComparison` lee la lista de compras sin la llave `week`.
- **Alto 2 (C1-15):** `getMealPlanWeek` aplica sin marca la regla del domingo, que sigue en el bloqueo abierto `U1-lista-vigente-domingo`.
- **Medio 1:** U1 §0.2 todavía llama `period` a la semana de la lista.
- **Medio 2:** el diccionario sigue sin `this_month`, `onward` ni `chat_id_masked`.

Las verificaciones mecánicas (a)–(f), (i) y (j) pasan. Los dos hallazgos de la ronda anterior que tocaban a C1 quedan cerrados en C1. El del diccionario sigue abierto.

## Verificaciones reejecutadas

| # | Comando (propio) | Resultado |
|---|---|---|
| a | `python3 -I backend/internal/api/c1gen/build.py --check` | IDÉNTICO en los 21 archivos: `openapi.yaml` (630,717 B), `routes.go` (25,004 B), `C1-datos.md` (128,721 B), `extractos/INDICE.md` y los 17 YAML. «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| b | `goose v3.28.0 -dir backend/db/migrations sqlite3 datos/tmp/rev-c1-r3/c1.db up` + `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db …` | 00001 y 00002 OK, versión 2. **24 OK, 0 FALLA, «RESULTADO: TODO OK»**, código de salida 0 |
| c | `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` (y cada extracto) | `openapi.yaml: OK`, los 17 extractos OK |
| d | `python3 -I docs/design/tools/verificar_u1.py` | «RESULTADO: TODO OK». N-xxx sin operación ni Local: []. N-316 a N-319 nuevas, N-317 Local. Cruce openapi ↔ C1-datos: [] |
| e | `cd backend && go build ./... && go vet ./... && gofmt -l .` (go1.27.2) | build OK, vet OK, gofmt sin archivos |
| f | script propio `chk1.py` sobre `extractos/` | 17 archivos (con `componentes-comunes.yaml`), el mayor de 76,202 B (≤ 81,920). Suma de operationId = 201 = `grep -c operationId openapi.yaml`. 0 repetidas, 0 faltantes, 0 sobrantes. Todos parsean como YAML. sha256 de la concatenación: `74476ad3…1a26` |
| g | scripts propios `chk2.py`, `chk3.py` y `PRAGMA table_info` | Ver la tabla de §K. Todo está, salvo la semana en `getBasketComparison` (H-r3-01) |
| h | `chk4.py`: recorre `x-provisional` y `PROVISIONAL` | 4 `x-provisional` y 6 líneas PROVISIONAL en openapi. Citan `C1-valores-k31`, `S1-ajuste-objetivo-kcal` y `criterios-promedios-redondeo`, los tres sin «Respuesta de Axel». C1-datos: 18 líneas. `C1-lista-compras` (respondido) solo sale en historia (Cambios v1.0 y §14). **No marcado:** la regla del domingo en `getMealPlanWeek` (H-r3-02) |
| i | `chk1.py`, `grep` | `info.version` 1.1.0. `openapi` 3.1.1. Cambios tiene las filas «C1 v1.1 (puerta 1)», «corrección r1» y «corrección r2» |
| j | `chk5.py`: compara cada número de los ejemplos con `oraculo/model.json` y `datos/semilla/data.json` (en pesos y en centavos) | openapi: 14 ejemplos, 0 cifras sin origen. C1-datos desde §1: 62 cifras y 0 inventadas. «3.29» es una sección. 138.50 = 3,233.00 − 3,094.50 de `data.json` (S5 K-03), igual que en HEAD. Los ejemplos de §7.1 (194750, 251737) salen de `devengado.despensa_total` |

## Criterios de C1.md v1.1

| Criterio | Estado | Evidencia |
|---|---|---|
| C1-01 (alta) | Cumple | (c). Además `openapi` 3.1.1 y 162 rutas |
| C1-02 (alta) | Cumple | (b) C1-02 da 203 N-xxx, 0 sin operación ni Local (N-003, N-301 y N-317 son Local), 0 pantallas sin operación, matriz §10 de 227 pares = U1 §3.29. La única operación sin N-xxx es `getOpenApiDocument` (estáticos). (d) coincide |
| C1-03 (alta) | Cumple | (b) 55 tablas, 0 inexistentes, 0 sin operación, 1003 propiedades con columna o marcadas como calculado o entrada. `PRAGMA` confirma las columnas de 00002 (ocr_draft_json, idempotency_key y captured_at, rate_decimals, fortnight_payment_cap_cents, goal, rate_base, min_payment_*) |
| C1-04 (alta) | Cumple | 915 de 915 respuestas de error van a `Error` |
| C1-05 (alta) | Cumple | 40 de 40 colecciones con cursor, limit, sort y Page |
| C1-06 (alta) | **No cumple** | Un componente por nombre (`Period` 31 usos, `ListWeek` → `week` 7 usos, `Family` con enum de 3 sin `x-provisional`). Pero `getBasketComparison` lee la lista sin `week` y con `period` (H-r3-01) |
| C1-07 (alta) | Cumple | Solo `date` y `date-time`. `DateTime` y §4 dicen que «Z» solo se acepta en entradas (R2-06). 0 ejemplos de respuesta con «Z» |
| C1-08 (alta) | Cumple | 265 propiedades `_cents` por `$ref` a MoneyCents, 0 `type: number` |
| C1-09 (alta) | Cumple | §9 con 17 métricas M-xx. Ejemplos 10,000.00, 6,905.50, 973.75 (30-oct-2026) y 1,947.50, todos del oráculo. Observación baja en M-06 (H-r3-05) |
| C1-10 (alta) | Cumple | §7 con 31 entidades. El generador no reporta entidades ni vistas sin decisión |
| C1-11 (alta) | Cumple | §11.6. `gen6_c1.py`: «capacidades con id malo: []» |
| C1-12 (alta) | Cumple | Solo `login` y `/health` son públicos. Cookie HttpOnly, Secure y SameSite=Strict, 90 días renovable (K11). Login: 5 fallos, «por usuario y por IP en ventana deslizante de 15 minutos», 429. 0 campos sensibles. Observación baja H8 |
| C1-13 (media) | Cumple | 0 operationId duplicados, 9 etiquetas, `info.version` 1.1.0, filas v1.1 en Cambios |
| C1-14 (media) | Cumple | 14 ejemplos válidos con jsonschema, 2765 descripciones en español. Observación baja H6 |
| C1-15 (alta) | **No cumple** | Las marcas que quedan citan bloqueos abiertos, y no hay marcas sobre K6 ni K7. Pero la regla del domingo de `getMealPlanWeek` no tiene marca (H-r3-02). Observación baja H7 |
| C1-16 (alta) | Cumple en C1; **falta el diccionario** | `BudgetScope`: enum `this_month`/`onward`, `default: onward`, sin `x-provisional`. §7.1 trae tres ejemplos coherentes. El diccionario no registra los valores (H-r3-04) |
| C1-17 (alta) | Cumple | `saveTicketScanDraft` (PUT `/ticket-scans/{attachment_id}/draft`, N-075, x-tablas attachment). `TicketScan.has_draft`/`draft`. «No crea gasto, precio ni inventario». Columna `attachment.ocr_draft_json` |
| C1-18 (alta) | Cumple en C1; **falta el diccionario** | `listTelegramChatCandidates` (N-299) y `confirmTelegramChatCandidate` (N-316), las dos con sesión global. La única escritura con «chat» es la confirmación, cuyo cuerpo trae solo `person_id`. Sin token. `chat_id_masked` falta en el diccionario (H-r3-04) |
| C1-19 (alta) | Cumple | `createExpense` y `createWeightEntry`: `IdempotencyKey` con required, 16–128 y patrón. `captured_at` DateTime. x-tablas con expense y weight_entry. 200 con `Idempotent-Replayed` y 409 `idempotency_key_reused`. Ciclo y reintento en §1 punto 8 |
| C1-20 (alta) | Cumple | `rate_decimals` en CreditTermRate, CreditPurchase y Promotion, y `annual_rate_decimals` en Credit. `rate_pct` en texto HALF_UP. Ejemplos 17.27 → 2 y 25.8685 → 4 (258685/1000000). Columnas en la base |
| C1-21 (alta) | Cumple | (a) y (f): regenerables byte a byte, ≤ 80 KB y cada operación una sola vez |

## Decisiones §K revisadas (g)

| K | Estado | Evidencia en `openapi.yaml` y `C1-datos.md` |
|---|---|---|
| K1 | OK | 5 operaciones de la lista, sin `x-provisional` |
| K2 | OK en C1 | `scope` con `default: onward`. `updateBudget`, `Budget.valid_to_date`, `BudgetCreate`/`BudgetUpdate` y §7.1 dicen lo mismo, incluida la eliminación de los renglones posteriores |
| K3 | OK | `saveTicketScanDraft` y `TicketScanDraft` |
| K4 | OK | candidatos y confirmación, `chat_id_masked` calculado |
| K5 | OK | `family` household/external/joint, inclusión descrita |
| K6 | OK | Limit máx. 200, Cursor 512, Q 200, Idempotency 16–128 y 24 h, SameSite=Strict, contraseña ≥ 12, PROFECO 1 GB. `LoginRequest.username` sin 3–64 (H8, baja) |
| K7 | OK | calculados sin marca (C1-15 de verify) |
| K9 | OK | ver C1-19 |
| K11 | OK | ver C1-12 |
| K12 | OK | decimales capturados. «Z» solo en entradas |
| K13 | OK | `Person.fortnight_payment_cap_cents` opcional, sin valor fijo y usado como objetivo, no como restricción |
| K16 | OK | `OptimizerSettings` sin valores por omisión y con `missing_constraints`. `ScenarioResult.status` incluye `missing_constraints` (no recomienda) y trae `constraints_used` |
| K18 | OK, con observación | `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents` y `missing_terms`. El IVA va en `vat_on_interest_bp`, donde NULL = «ya incluye IVA» (H-r3-06, baja) |
| K21 | OK | `get/updateProjectionSettings` (N-318 y N-319). `repeat_cadence_enabled` apagado por omisión |
| K23 | OK | `HealthProfile.goal` lose/maintain/gain. `suggested_kcal_target` PROVISIONAL con `S1-ajuste-objetivo-kcal` (abierto) |
| R2-06 | OK | ver C1-07 |
| R3-01 | **Parcial** | 7 operaciones con `week`, pero `getBasketComparison` no (H-r3-01) |

## Hallazgos

| Id | Severidad | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| H-r3-01 | alta | C1-06, R3-01 | `getBasketComparison` (GET `/api/v1/basket-comparison`, la única operación de P-14 para N-161–N-164 y de P-01 para N-162) tiene `basket=shopping_list` por omisión y lee `shopping_list` y `shopping_list_item`. No recibe `week`: usa el filtro `period`, opcional y con 5 formatos, sin decir qué significa ahí (la descripción es «Comparar canasta.»). §6 dice «La lista de compras no usa este filtro: su semana va en la llave `week`», y U1 §0.2 dice «Toda lectura y escritura de la lista lleva la semana que la pantalla muestra». `listBaskets` devuelve la lista de la semana X, pero la comparación no puede recibir X, así que con la lista de compras la operación no sabe qué lista comparar. La fuente no cambió desde HEAD (`gen3_ops.py` línea 154), y la matriz §6 la marca en «periodo» | `docs/contratos/openapi.yaml` (getBasketComparison), `docs/contratos/C1-datos.md` §6, `backend/internal/api/c1gen/gen3_ops.py` línea 154, `verify_c1.py` (C1-06) | Agregar `ListWeek` a `getBasketComparison`, obligatoria con `basket=shopping_list` (422 si falta). Quitar `Period` o definir qué filtra con `basket=expense`. Describir la operación. Matriz §6: «periodo» vacío y `week` ✓. En `verify_c1.py`, exigir `ListWeek` a toda operación que lea `shopping_list` |
| H-r3-02 | alta | C1-15 | `getMealPlanWeek` (P-10, N-120) toma por omisión «la semana en curso y, los domingos, la siguiente» (`openapi.yaml` línea 3762, `gen3_ops.py` línea 123). La regla del domingo está en el bloqueo **abierto** `U1-lista-vigente-domingo`, sin respuesta de Axel. Ese bloqueo dice «Si C1 conserva un valor por omisión de `period`, debe ser esta misma regla», y U1 §0.2 la marca PROVISIONAL. En C1 no lleva `x-provisional` ni la marca PROVISIONAL, no cita el bloqueo y no está en §13. Viene de v1.0 y las rondas anteriores no la detectaron. `verify_c1.py` C1-15 no la ve porque solo revisa las marcas que ya existen | `docs/contratos/openapi.yaml` línea 3762, `docs/contratos/C1-datos.md` §13, `backend/internal/api/c1gen/gen3_ops.py` línea 123 | Opción 1: marcar el valor por omisión con `x-provisional` y «PROVISIONAL: bloqueo U1-lista-vigente-domingo», y agregar la fila en §13. Opción 2: quitar el valor por omisión y hacer `period` obligatorio en `getMealPlanWeek`, igual que `week` en la lista, para que la regla quede en la PWA (U1 §0.2) |
| H-r3-03 | media | coherencia C1 ↔ U1 (C1-02, R3-01) | La corrección r2 cambió el nombre de la semana de la lista a `week`. Pero `docs/design/U1-arquitectura.md` línea 66 (§0.2, definición única de R3-01 y U1-16) sigue diciendo «lleva la semana (`period`) que la pantalla muestra» y «Si C1 conserva un valor por omisión de `period`…». La solicitud 3 de `correccion-C1-r2.md` solo cubre `docs/criterios/U1.md` U1-16 y no menciona U1-arquitectura.md, así que los dos artefactos nombran distinto el mismo parámetro | `docs/design/U1-arquitectura.md` línea 66, `docs/criterios/U1.md` líneas 79–80 | Solicitud al dueño de U1 y de los criterios: decir «la semana (`week` en C1, obligatoria)» en U1 §0.2 y en U1-16. Mientras tanto, registrarlo en C1-datos §14 como inconsistencia pendiente |
| H-r3-04 | media | C1-16, C1-18 (H5 de r1, H-r2-03) | `diccionario-datos.md` sigue sin `this_month`, `onward` y `chat_id_masked` (`grep -c` = 0, 0 y 0). La línea 530 (`budget.valid_to_date`) dice «deja NULL», pero no dice que se eliminan los renglones posteriores. C1-16 pide que el diccionario registre los valores, y C1-18 que fije los campos del candidato | `docs/contratos/diccionario-datos.md` (línea 530 y la tabla de `telegram_chat_candidate`), `backend/db/migrations/00001_esquema.sql` línea 299 | El integrador o S1 aplican las solicitudes 1 y 2 de `docs/estado/p1/correccion-C1-r2.md`. Después se quita el parche con `assert` de `gen2_schemas.py` |
| H-r3-05 | baja | C1-09 | En §9, M-06 (línea 329) tiene periodo «lista vigente» y filtros «—». `getDashboardBestBasket` exige `week` y el contrato no define una lista «vigente» por omisión | `docs/contratos/C1-datos.md` línea 329 (`c1_plantilla.md` o `gen6_c1.py`) | Periodo: «semana `week` que manda la pantalla (llave, §6)» |
| H-r3-06 | baja | K18 | En `Credit.vat_on_interest_bp`, NULL significa «la tasa ya incluye IVA». Por eso no se distingue «falta capturar si incluye IVA», y `missing_terms` no lo incluye, aunque K18 pide «falta capturar» cuando no hay dato | `docs/contratos/openapi.yaml` (Credit), columna de S1 | Solicitud a S1: tres estados (incluye, no incluye y sin dato) o una columna `interest_includes_vat` que admita NULL. Luego agregarla a `missing_terms` |
| H6 | baja | C1-14 | Siguen descripciones que solo repiten el identificador (ejemplos de la muestra: `notes`, `is_active`, `id`, `created_at`) | `openapi.yaml` (`gen2_schemas.py`) | La que propone r1 |
| H7 | baja | C1-15 | Siguen las marcas PROVISIONAL de `DashboardWeekMeals` (línea 11924) y de `NotificationUpdate` (línea 17167) sin `x-provisional` | `openapi.yaml` (`gen2_schemas.py`) | La que propone r1 |
| H8 | baja | C1-12, K6 | `LoginRequest.username` sigue sin el rango 3–64 | `openapi.yaml` (`gen2_schemas.py`) | La que propone r1 |

## Estado de los hallazgos anteriores

| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| [alta] C1-06 (R3-01): `ListWeek` y `Period` con el mismo nombre | **Cerrado** en su forma original | `ListWeek` se llama `week` y hay un componente por nombre (`param names multi: {}`). `Period` tiene 31 usos y `ListWeek` 7. §6 trae la tabla «Llaves que no son filtros» y la columna `week`, con `listBaskets` y `getDashboardBestBasket`. `getShoppingList` lleva marca en tienda y en `week`. verify_c1 C1-06 tiene la regla nueva. **R3-01 se reabre en parte** por `getBasketComparison` (H-r3-01) |
| [media] C1-16 (K2): textos de `onward` que no coincidían | **Cerrado en C1** | `updateBudget`, `Budget.valid_to_date` (parche con assert), `BudgetCreate`/`BudgetUpdate` y §7.1 dicen todos «sin fin y elimina los renglones posteriores». El tercer ejemplo (7: oct, 10: desde nov, se eliminan 8 y 9) es coherente y sin traslape, con cifras del oráculo. Lo que falta en el diccionario pasa a H-r3-04 |
| [media] C1-16, C1-18 (H5): diccionario sin `scope` ni `chat_id_masked` | **Sigue abierto** | `grep -c` = 0, 0 y 0. Pasa a H-r3-04 |
| [baja] H6, H7, H8 | Siguen abiertos | Ver la tabla de hallazgos |
