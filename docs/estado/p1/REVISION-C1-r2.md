# REVISIÓN C1 v1.1 · ronda 2 de 3 (flujo P1, 11-oct-2026)

Revisor independiente: `revisor-C1-r2`. Evalúa solo contra `docs/criterios/C1.md` v1.1, `decisiones.md` §K (K1–K7, K9, K11–K13, K16, K18, K21, K23) y los hallazgos R2-06 y R3-01. Todas las verificaciones se volvieron a correr con comandos y scripts propios en `datos/tmp/revC1r2/`, ya borrados. No se corrigió nada.

## Veredicto

**NO APROBADO**: 1 hallazgo alto, 2 medios y 3 bajos.
- Los cuatro hallazgos de r1 que tocaban a C1 (H1 a H4) quedan **cerrados**.
- H5 (el diccionario) **sigue abierto**. No se arregla en C1 porque el diccionario es de S1.
- Hay un hallazgo **alto nuevo** que salió de la corrección de R3-01: el filtro `period` ahora tiene dos definiciones distintas.
- Hay un hallazgo **medio nuevo** sobre el alcance `onward` de K2 cuando ya hay renglones posteriores.

## Verificaciones reejecutadas

| # | Comando (resumen) | Resultado |
|---|---|---|
| a | `python3 -I backend/internal/api/c1gen/build.py --check` | IDÉNTICO en los 21 archivos: `openapi.yaml` 629,639 B, `routes.go` 25,004 B, `C1-datos.md` 124,433 B, `extractos/INDICE.md` y los 17 YAML. Se corrió dos veces y el `sha256sum` de `extractos/` salió igual las dos (9e0fe410…) |
| b | goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/revC1r2/c1.db up` y después `uv run … verify_c1.py --db …` | 00001 OK, 00002 OK, versión 2. **23 OK, 0 FALLA, «TODO OK»** |
| c | `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator` | `openapi.yaml: OK`. Los 17 YAML de `extractos/` también OK |
| d | `python3 -I docs/design/tools/verificar_u1.py` | TODO OK. Faltantes de U1-03: `[]`. Locales: N-003, N-301 y N-317. N-316, N-318 y N-319 tienen operación |
| e | `go build ./... && go vet ./... && gofmt -l .` (go1.27.2) | build OK, vet OK, gofmt vacío |
| f | script propio sobre `extractos/` | 17 archivos, todos parsean y declaran `openapi` 3.1.1. Suma de operationId: 201 = `grep -c operationId openapi.yaml` (201). 0 repetidas, 0 faltantes, 0 sobrantes. El mayor es `gastos.yaml` con 76,202 B (≤ 81,920). `componentes-comunes.yaml` sí está |
| g | script propio con la lista K | Ver la tabla «Lista K» |
| h | recorrido de `x-provisional` y PROVISIONAL | 4 nodos `x-provisional` (adjuntos, `compareScenarios`, `ExpenseSummary.weekly_average_cents` y `HealthProfile.suggested_kcal_target`) y 6 líneas PROVISIONAL. Todas citan `C1-valores-k31`, `criterios-promedios-redondeo` o `S1-ajuste-objetivo-kcal`. Los tres archivos existen y ninguno tiene «Respuesta de Axel» |
| i | `info.version` y Cambios | `1.1.0` (HEAD: `1.0.0`). Filas «C1 v1.1 (puerta 1)» y «C1 v1.1 · corrección r1» |
| j | cifras de las líneas agregadas (git diff contra `oraculo/model.json` y `datos/semilla/data.json`) | 0 cifras inventadas. 199.00 (`parametros.papel`), 1,947.50 y 2,517.37 (`despensa_total`) y 97375 salen del oráculo. 25.8685 viene de K12. 1.375, 1.55 y 1.725 vienen del diccionario y de RN-17. «•••4821» es un marcador. 8601, 7396 y 3.29 son referencias (ISO, RFC y U1 §3.29). Los 14 `example` de `openapi.yaml` validan y ninguno termina en «Z» |

## Lista K (g)

| Decisión | Estado | Evidencia |
|---|---|---|
| K1 lista de compras | OK | Hay 5 operaciones (`getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`), ninguna con `x-provisional` |
| K2 `scope` | OK, con hallazgos H-r2-02 y H-r2-03 | `updateBudget`: `scope` es un enum `[this_month, onward]` con default `onward` y sin `x-provisional`; recibe también `month`. §7.1 tiene la tabla y un ejemplo de cada alcance |
| K3 borrador | OK | `saveTicketScanDraft` (PUT …/draft, N-075, `x-tablas` [attachment, audit_log]) no estaba en HEAD. `TicketScan` trae `draft`, `has_draft` y `ocr_draft_saved_at`. La base migrada tiene `attachment.ocr_draft_json` y `ocr_draft_saved_at`. C1-datos.md línea 27 y §14 fila 4: el borrador no crea precio |
| K4 candidatos | OK | `listTelegramChatCandidates` (N-299) y `confirmTelegramChatCandidate` (N-316, cuerpo `person_id`), las dos con la seguridad global `cookieAuth`. Ningún POST, PUT ni PATCH recibe `chat_id`. Se retira `createNotificationRecipient`. Las respuestas no traen token, hash ni secreto |
| K5 `family` | OK | El enum `[household, external, joint]` describe la inclusión y no tiene `x-provisional` |
| K6 valores | OK (H8 de r1, baja) | `limit` máx. 200, `cursor` 512, `q` 200, clave 16–128, `new_password` minLength 12, PROFECO 1 GB y `SameSite=Strict`, todos sin marca |
| K7 calculados | OK | Ningún campo de K7 tiene `x-provisional` |
| K9 idempotencia | OK | `createExpense` y `createWeightEntry` piden `Idempotency-Key` con required, 16–128 y patrón. `captured_at` hace `$ref` a `DateTime` en la entrada y en la salida. Las respuestas incluyen 200 y 201. La base tiene `expense.idempotency_key` y `weight_entry.idempotency_key`. C1-datos.md §1.8 trae el ciclo y el ejemplo |
| K11 sesión | OK | El login dice «Conteo por usuario y por IP en ventana deslizante de 15 minutos» con 429 `too_many_attempts`. `cookieAuth` dice 90 días renovables |
| K12 decimales y «Z» | OK | `rate_decimals` en la entrada y la salida de `CreditTermRate`, `CreditPurchase` y `Promotion`. `rate_pct` es un texto calculado (17.27 con 2, 25.8685 con 4). La nota «Z» solo en entradas está en §4 y en `DateTime`. 0 propiedades con «(UTC)» |
| K13 tope por persona | OK | `fortnight_payment_cap_cents` opcional en `Person` y `PersonUpdate`, y en la base |
| K16 restricciones | OK | `OptimizerSettings` trae `missing_constraints`. `ScenarioRequest` las sobreescribe. `ScenarioResult.status` es `ok`, `missing_constraints` o `no_valid_scenario` |
| K18 campos de credit | OK | `rate_base` (daily o monthly), `min_payment_pct_bp`, `min_payment_fixed_cents`, `annual_rate_decimals` y el calculado `missing_terms`, todos en la base |
| K21 proyección | OK | `getProjectionSettings` (N-318) y `updateProjectionSettings` (N-319). `projection_mode` en las 3 lecturas |
| K23 objetivo | OK | `HealthProfile.goal` es un enum `[lose, maintain, gain]`. `suggested_kcal_target` queda PROVISIONAL por el bloqueo abierto `S1-ajuste-objetivo-kcal` |
| R3-01 `period` de la lista | Obligatorio: OK. Definición: **falla** (H-r2-01) | 7 operaciones con `period` required y patrón de semana |

## Criterios

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| C1-01 | alta | Cumple | 3.1.1 y 162 rutas. Validador 0.9.0 OK |
| C1-02 | alta | Cumple | verify_c1 C1-02 (6 OK). verificar_u1 TODO OK. `build.py --check` IDÉNTICO |
| C1-03 | alta | Cumple | 55 tablas de 2 migraciones. 0 inexistentes y 0 sin operación |
| C1-04 | alta | Cumple | 100 % de las respuestas de error van a `Error` |
| C1-05 | alta | Cumple | 40 de 40 |
| C1-06 | alta | **No cumple** | `period` tiene dos definiciones: `Period` y `ListWeek`. La matriz §6 perdió la marca de periodo en `getShoppingList` (H-r2-01) |
| C1-07 | alta | Cumple | Solo `date` y `date-time`. 0 ejemplos de respuesta con «Z». La nota está en §4 y en `DateTime`. H4 de r1 cerrado |
| C1-08 | alta | Cumple | 265 propiedades `_cents` por `$ref` y 0 `number` |
| C1-09 | alta | Cumple | §9 con 17 filas M. Las cifras 10,000.00, 6,905.50, 973.75 y 1,947.50 están |
| C1-10 / C1-11 | alta | Cumple | §11.1: 203 N y 41 pantallas. §11.6: 36 filas, ninguna sin operación |
| C1-12 | alta | Cumple | Solo login y salud son públicos. Cookie Strict, 90 días, usuario e IP |
| C1-13 | media | Cumple | v1.1, ids únicas, 9 etiquetas y 1.1.0 |
| C1-14 | media | Cumple con observación | 14 de 14 ejemplos válidos. H6 de r1 (baja) sigue: 202 descripciones son solo el identificador |
| C1-15 | alta | Cumple con observación | Todas las marcas citan un bloqueo abierto y no queda ninguna sobre K6 ni K7. H7 de r1 (baja) sigue |
| C1-16 | alta | **No cumple del todo** | El parámetro y §7.1 están bien. El diccionario no registra `this_month` ni `onward` (H-r2-03), y `onward` contradice `valid_to_date` si hay renglones posteriores (H-r2-02) |
| C1-17 | alta | Cumple | Operación nueva N-075, campo `draft`, columna en la base y texto «no crea precio» en C1-datos.md. H3 de r1 cerrado |
| C1-18 | alta | Cumple con observación | 2 operaciones con sesión y 0 `chat_id` libre. `chat_id_masked` no está en el diccionario (H-r2-03) |
| C1-19 | alta | Cumple | Ver K9. Ejemplo de reintento y qué pasa pasadas las 24 h en §1.8. H2 de r1 cerrado |
| C1-20 | alta | Cumple | Ver K12. La columna está en la base y en el diccionario (7 menciones) |
| C1-21 | alta | Cumple | Ver (f). Los extractos son deterministas: mismo sha256 en las dos corridas |

## Hallazgos

| Id | Sev. | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| H-r2-01 | **alta** | C1-06 (R3-01) | La corrección de R3-01 creó `components/parameters/ListWeek` con `name: period`, `required: true` y solo el patrón de semana. Ya existía `Period`, también `name: period`, opcional y con 5 formatos. El mismo filtro tiene ahora dos definiciones con valores permitidos distintos, y C1-06 exige una sola. Además, la matriz recurso × filtro de §6 marcaba periodo en `getShoppingList` en HEAD y ahora no (línea 172), y `listBaskets` y `getDashboardBestBasket` piden `period` pero no aparecen en la matriz. `verify_c1.py` C1-06 no lo detecta porque `F6` (línea 150) solo cuenta los `$ref` a `Period` | `docs/contratos/openapi.yaml` (`ListWeek`, línea 8954; `Period`, 8888); `docs/contratos/C1-datos.md` §6; fuente `backend/internal/api/c1gen/gen1_base.py`/`gen3_ops.py`, `gen6_c1.py`, `verify_c1.py` | Dejar una sola definición por nombre, con una de dos opciones. (a) Llamar a la semana de la lista con otro nombre que no sea filtro (p. ej. `week`, obligatorio, semana ISO) y documentarla en §6 como la llave de la lista semanal. (b) Volver al `Period` común, opcional, con la semana vigente de U1 §0.2 por omisión (la primera opción de R3-01) y 422 para formatos que no sean semana. En los dos casos: regenerar la matriz §6 con las 3 lecturas y agregar a verify_c1 C1-06 la regla «cada nombre de filtro, un solo componente» |
| H-r2-02 | media | C1-16 (K2) | Hay dos textos sobre `onward` que no coinciden. `Budget.valid_to_date` (viene del comentario de S1) dice que «de este mes en adelante» «deja NULL». `updateBudget` y §7.1 dicen que el renglón nuevo rige «hasta el fin original». Ninguno dice qué pasa con los renglones posteriores de la misma categoría. Ejemplo: después del ejemplo `this_month` de §7.1 quedan el renglón 7 (oct–nov), el 8 (solo dic, 2,517.37) y el 9 (desde ene, 1,947.50). Si se cambia diciembre otra vez con `onward`, §7.1 modifica el renglón 8 y deja su fin en 2026-12-31, así que enero sigue en 1,947.50 y el monto no rige «en adelante». Si se lee como «deja NULL», el 8 se traslapa con el 9 y responde 409 | `docs/contratos/openapi.yaml` (`Budget.valid_to_date`, `updateBudget`); `docs/contratos/C1-datos.md` §7.1; `docs/contratos/diccionario-datos.md` línea 530 | Decir en §7.1 y en `updateBudget` qué hace `onward` con los renglones posteriores (p. ej. los cierra o los reemplaza para que el monto nuevo rija desde `month` sin fin) y agregar ese caso al ejemplo. Igualar la descripción de `Budget.valid_to_date` y pedir a S1 el mismo cambio en el diccionario. Si K2 no alcanza para decidirlo, escribir el bloqueo `C1-presupuesto-onward` |
| H-r2-03 | media | C1-16, C1-18 (H5 de r1) | Sigue abierto. El diccionario no registra `this_month`, `onward` ni `chat_id_masked` (`grep -c` da 0, 0 y 0). La parte de C1 ya está hecha: la solicitud 1 de `correccion-C1-r1.md` lo pide a S1, pero nadie la ha aplicado | `docs/contratos/diccionario-datos.md` | El integrador o S1 aplican la solicitud 1 de `correccion-C1-r1.md`. No hay que cambiar nada en C1 |
| H6 | baja | C1-14 | Sigue de r1: 202 descripciones de propiedad son solo el identificador | `openapi.yaml` (`gen2_schemas.py`) | La que propone r1 |
| H7 | baja | C1-15 | Sigue de r1: las marcas PROVISIONAL de `DashboardWeekMeals` y de `NotificationUpdate` no llevan `x-provisional` | `openapi.yaml` (`gen2_schemas.py`) | La que propone r1 |
| H8 | baja | C1-12/K6 | Sigue de r1: `LoginRequest.username` no declara el rango 3–64 | `openapi.yaml` (`gen2_schemas.py`) | La que propone r1 |

## Estado de los hallazgos de r1

| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H1 (alta, C1-16, K2) | **Cerrado** | §7.1 (líneas 225–262): tabla del efecto de cada alcance sobre las vigencias y un ejemplo de cada uno, con montos del oráculo (194750 y 251737). Queda el caso de los renglones posteriores, que pasa a H-r2-02 |
| H2 (alta, C1-19, K9) | **Cerrado** | §1.8 trae primer envío 201, reintento 200 con el mismo `id`, 409 si cambia el cuerpo dentro de 24 h, qué pasa pasadas las 24 h, y el ejemplo `POST /api/v1/expenses`. El cuerpo es válido: tiene los 4 campos obligatorios y `payment_method` `cash` está en el enum |
| H3 (media, C1-17, K3) | **Cerrado** | La línea 27 de Cambios y la fila 4 de §14 dicen «no crea gasto, precio ni inventario» |
| H4 (media, C1-07, R2-06) | **Cerrado** | 0 propiedades con «(UTC)». El único «(UTC)» es la nota del esquema `DateTime` (parche en `gen1_base.py`, `col_schema`) |
| H5 (media, C1-16/C1-18) | **Abierto** (H-r2-03) | El diccionario sigue sin los tres valores. Ya se pidió en `correccion-C1-r1.md` |
