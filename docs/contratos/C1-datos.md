# C1 · Contrato de datos y API · Finanzas Hogar

| | |
|---|---|
| Versión | **C1 v1.3** (11-oct-2026, puerta 2) · `info.version` 1.3.0 de `docs/contratos/openapi.yaml` |
| Autor | Agente S6a-contrato-c1 (flujo WF-0 «Cimientos», fase S6 parte a); corrección r1: agente correccion-C1-r1; corrección r2: agente correccion-C1-r2; v1.1 corrección r1: agente correccion-C1-r1; v1.1 corrección r4: agente correccion-C1-r4; v1.2: agente C1-v1.2; v1.2 corrección r1: agente correccion-C1-r1; v1.2 ronda 3 (P1d): agente correccion-C1-r3; v1.3: agente C1-puerta2 (P2c) |
| Estado | Puerta 2 aplicada (C1 v1.3): respuestas de Axel de `docs/estado/puerta2/RESPUESTAS.md` sobre el esquema 1.4 de S1. Las preguntas sin recomendación quedan con la opción en uso y no cambian C1 |
| Fuentes | `plan.md` §1, §3.4 (S6, S7), §3.5 (U1), §3.7; `decisiones.md`; esquema de S1, versión 1.4 (00001 a 00004: `backend/db/migrations/00001_esquema.sql`, `00002_puerta1.sql`, `00003_puerta1b.sql` y `00004_puerta2.sql`; `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md`); U1 1.2 (`docs/design/U1-arquitectura.md`, necesidades N-xxx); puerta 2: `docs/estado/puerta2/PUERTA-2.md` y `RESPUESTAS.md`, contratos de motores S2 1.3, S4 1.5 y S5 1.5 (`docs/contratos/motores/`) y U3 1.3 |
| Archivo de máquina | `docs/contratos/openapi.yaml` (OpenAPI 3.1.1): 168 rutas, 212 operaciones, 255 esquemas |
| Generador | Este documento, `openapi.yaml` y `backend/internal/api/routes.go` son archivos generados: no se editan a mano. Fuente: `backend/internal/api/c1gen/` (este texto sale de `c1_plantilla.md`; ver su `README.md`). Regenerar: `python3 -I backend/internal/api/c1gen/build.py`. Comprobar que nadie los editó a mano: `python3 -I backend/internal/api/c1gen/build.py --check`. Verificar C1-01 a C1-15, la matriz §10 y S6-03: `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py` |

Este documento fija **cómo** se piden los datos que U1 dice que cada pantalla **necesita**. El OpenAPI es la fuente de verdad de rutas y esquemas; aquí están las convenciones únicas, las métricas con su fórmula, la matriz pantalla → necesidad → endpoint y la cobertura de S1. Todo lo que no sale del plan, de las decisiones o del oráculo va marcado **PROVISIONAL** y apunta a su bloqueo en `docs/estado/bloqueos/`.

## Cambios

| Versión | Fecha | Cambio |
|---|---|---|
| C1 v1.0 | 10-oct-2026 | Primera versión (WF-0). Los cambios posteriores se publican como C1 v1.1, v1.2… (plan §3.7) con esta tabla actualizada y `info.version` igual. |
| C1 v1.0 · corrección r1 | 10-oct-2026 | Sin cambios de rutas, esquemas, convenciones ni matriz. El generador y el verificador pasan al repositorio (`backend/internal/api/c1gen/`) y las cabeceras de `openapi.yaml` y `routes.go` lo citan (hallazgo H-05 de la revisión r1). |
| C1 v1.0 · corrección r2 | 11-oct-2026 | Regenerado contra el esquema 1.1 de S1 (hallazgos R2-01 y R2-02 de la revisión r2; sigue siendo la entrega de WF-0 para la puerta 1). Lista de compras guardada: `generateShoppingList` (N-127), `createShoppingListItem` (N-152), `updateShoppingListItem` (N-151, N-153), `assignShoppingListStore` (N-165) y `getShoppingList` con la lista guardada y el filtro `is_bought` (N-150, N-073), todo PROVISIONAL por la opción A de `C1-lista-compras`. Presupuesto con vigencia (`valid_from_date`, `valid_to_date`; sin `period_id`) y 409 `budget_overlap`. Tiempos aptos de la receta como lista `meal_slot_ids` sobre `recipe_meal_slot` (sin `default_meal_slot_id`). Descripción del gasto opcional sin texto de relleno. Movimientos de inventario guardados en `inventory_movement`, con motivo en los ajustes. Ingreso esperado con `second_half_amount_cents`. `idempotency_key` en `x-tablas` de los POST de captura. Tope por quincena sin cifra. Columna nueva `category.kind`. |
| C1 v1.1 (puerta 1) | 11-oct-2026 | Aplica las respuestas de Axel de la puerta 1 (`decisiones.md` §K) y U1 1.1 contra el esquema 1.2 de S1 (migraciones 00001 y 00002). Detalle en «C1 v1.1 (puerta 1)», abajo. `info.version` = 1.1.0. |
| C1 v1.1 · corrección r1 | 11-oct-2026 | Sin cambios de rutas, parámetros, tipos ni cifras; solo textos (hallazgos H1 a H5 de `docs/estado/p1/REVISION-C1-r1.md`). §7.1 nueva: efecto de cada `scope` de `updateBudget` sobre las vigencias, con un ejemplo de cada uno (K2). §1 punto 8: ciclo de una captura de gasto o de peso antes y después de las 24 h, con ejemplo de reintento (K9). El borrador del ticket no crea gasto, precio ni inventario (K3). Las descripciones de los instantes ya no dicen «(UTC)», que contradecía la salida con el desplazamiento de America/Mexico_City (R2-06, K12). `info.version` sigue en 1.1.0. |
| C1 v1.1 · corrección r2 | 11-oct-2026 | Hallazgos H-r2-01 y H-r2-02 de `docs/estado/p1/REVISION-C1-r2.md`. La semana de la lista de compras ya no se llama `period`: es la llave `week` (`components/parameters/ListWeek`, semana ISO, obligatoria, sin valor por omisión) en `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets` y `getDashboardBestBasket`; el filtro `period` vuelve a tener una sola definición, §6 explica las llaves que no son filtros y la matriz recurso × filtro suma la columna de `week` (C1-06). `scope=onward` de `updateBudget`: el monto nuevo rige desde `month` sin fin y elimina los renglones posteriores de la misma categoría; §7.1 lo dice y trae el ejemplo, y `Budget.valid_to_date`, `BudgetCreate`/`BudgetUpdate` y `createBudget` dicen lo mismo (C1-16, K2). Cambia el nombre de un parámetro de consulta en 7 operaciones; sin cambios de rutas, cuerpos ni cifras. `info.version` sigue en 1.1.0. |
| C1 v1.1 · corrección r4 | 11-oct-2026 | Hallazgos H-r3-01, H-r3-02 y H-r3-04 de `docs/estado/p1/REVISION-C1-r3.md`. `getBasketComparison` recibe la llave `week` (`ListWeek`, obligatoria; 422 `validation_failed` si falta) y ya no recibe el filtro `period`; su descripción dice qué compara con `basket=shopping_list` y con `basket=expense`; la matriz de §6 y la tabla «Llaves que no son filtros» la incluyen. El valor por omisión de `getMealPlanWeek` (los domingos, la semana siguiente) queda `x-provisional` con el bloqueo `U1-lista-vigente-domingo`, en §13 y en N-120 de §10. `Budget.valid_to_date` toma el texto completo del comentario de S1 (que ahora nombra `onward` y `this_month`) y se quitó el parche de `gen2_schemas.py`. `verify_c1.py`: C1-06 exige `week` a toda operación que lee `shopping_list` y C1-15 marca la regla del domingo sin marca |
| C1 v1.2 (puerta 1b) | 11-oct-2026 | Aplica K35 a K44 (`decisiones.md` §K) y la sección «Para C1» de `docs/estado/p1c/U1-1.2.md` contra el esquema 1.3 de S1 (migraciones 00001, 00002 y 00003). Parámetros configurables (K44): cada clave de `app_setting` del inventario de S1 se lee con su valor vigente y su valor inicial y se cambia o restablece por clave; tabla generada en §11.7. Operaciones nuevas: `getSecuritySettings` y `updateSecuritySettings` (N-320), `getAlertSettings` y `updateAlertSettings` (N-323), `dismissTelegramChatCandidate` (N-324), `getDataSettings` y `updateDataSettings` (N-325), `getHealthSettings` y `updateHealthSettings` (N-326); N-321 y N-322 entran en las de comparador y precios; N-327 en `createHealthProfile` y `updateHealthProfile`. K36: `HealthProfile` con el ajuste por persona y el de la casa. K38, K39 y K43: sin marcas (N-023, N-042, N-064, N-120, N-226, N-230, N-297, N-300, N-311 y N-313 FIRMES en §10). `info.version` 1.2.0. Detalle en «C1 v1.2 (puerta 1b)», abajo. |
| C1 v1.2 · corrección r1 | 11-oct-2026 | Hallazgo A-01 de `docs/estado/p1c/REVISION-C1-r1.md` (C1-23, K43). `params_json` de `Notification` y `NotificationUpdate` declara las propiedades `send_time` (`type: string`, patrón `HH:MM` de 24 horas) y `advance_days` (`type: integer`, `minimum: 0`), sin valor por omisión; `weekday` y `day_of_month` siguen por `additionalProperties: true` y la marca de lugar `x-provisional` (bloqueo `S1-parametros-avisos`) se queda. Detalle al final de §11.7. `verify_c1.py` comprueba C1-23. Sin cambios de rutas, operaciones, cifras ni versión (`info.version` 1.2.0). |
| C1 v1.2, ronda 3 (P1d) | 11-oct-2026 | Hallazgos bajos N-01, N-02, N-03, B-02, B-03, B-04, B-05 e I-01 de `docs/estado/p1c/REVISION-C1-r2.md`. N-01: una sola constante de versión (`info.version` 1.2.0) para los encabezados generados de `openapi.yaml`, `routes.go` y los extractos. N-02: `minimum: 0` en `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` de `HealthProfileCreate` y `HealthProfileUpdate` (K36). N-03: el estado de la tarjeta (N-205) cita RN-26 y S3-08 (K18) en vez de un archivo de bloqueo. B-02: `NotificationUpdate` ya no lleva la marca de lugar en el esquema entero (`is_enabled` es FIRME); la llevan `params_json` y sus claves `send_time` y `advance_days`, en `Notification` y en `NotificationUpdate` (bloqueo `S1-parametros-avisos`). B-03: la descripción de `params_json` dice que los umbrales y los días antes del pago viven en `app_setting` desde 1.3 y que ahí quedan `send_time`, `advance_days`, `weekday` y `day_of_month`, con valor FIRME (decisión F) en las dos últimas. B-04: el mensaje de `verify_c1.py` sobre la regla del domingo dice FIRME (K38). B-05: el reporte de C1 v1.2 trae los comandos completos de verificación. I-01: esta cabecera dice puerta 1b aplicada, esquema 1.3 y U1 1.2. La nota de `send_time` en la 00003 dice que lo propone S7 en WF-3 (K31, puerta 4). Sin cambios de rutas, operaciones, cifras ni versión (`info.version` 1.2.0). |
| C1 v1.3 (puerta 2) | 11-oct-2026 | Aplica las respuestas de Axel de la puerta 2 (`docs/estado/puerta2/RESPUESTAS.md`: lo recomendado de `PUERTA-2.md` más Q2-11 B) contra el esquema 1.4 de S1 (migraciones 00001 a 00004). Reportes y avisos: hora, día y anticipación pasan a `app_setting` con lectura y cambio por clave (operaciones nuevas `getReportSettings` y `updateReportSettings`; `AlertSettings` suma tres claves); `params_json` ya no declara `send_time` ni `advance_days`. Rango por clave en el cambio (`minimum`, `maximum`, `pattern`) y tipo nuevo `SettingParameterTime`. `Credit.missing_terms` con `vat` (5 datos) y `credit.vat_included`. Alcance `scope` en `CostShareRuleUpdate` (§7.2). Lista de compras con tienda asignada, fecha de visita y totales de Q2-42 A; comparativa con `cheapest_store_id`; acumulados de T-2; décimas del avance de una meta; segunda tienda de la tarjeta 5; tienda del escenario. Detalle en «C1 v1.3 (puerta 2)», abajo. `info.version` = 1.3.0. |

### C1 v1.3 (puerta 2)

- **Entradas**: el generador aplica `00001` a `00004_puerta2.sql` (esquema 1.4). `parse_schema.py` lee además `@param_rango` (rango aprobado que se suma a su clave) y `@param_sin_valor_cerrado` (un `@param_sin_valor` que ya tiene clave), y toma la última `@col_note` de cada columna. El rango de cada clave no se copia: se prueba contra la base temporal (cada valor con un `UPDATE` que se deshace), así `minimum`, `maximum` y el patrón de la hora salen de los disparadores de S1 y el generador se detiene si el patrón no coincide con la base.
- **Reportes y avisos (Q2-11 B, K47; Q2-72 y Q2-73, K50)**: `getReportSettings` y `updateReportSettings` (`/api/v1/settings/reports`, N-296 y N-297) leen y cambian por clave `reports.daily_send_time` (21:30), `reports.weekly_weekday` (7) y `reports.weekly_send_time` (19:00), `reports.monthly_day_of_month` (1; en un mes más corto sale el último día, Q2-66 B) y `reports.monthly_send_time` (08:00). `AlertSettings` suma `alerts.send_time` (09:00), `alerts.stale_prices_advance_days` (1) y `alerts.meal_plan_missing_weekday` (7), y sus operaciones cubren además N-298 y N-300. `Notification.params_json` y `NotificationUpdate.params_json` ya no declaran `send_time` ni `advance_days` y pierden la marca de lugar: el bloqueo `S1-parametros-avisos` está resuelto. `updateNotification` sigue para activo y parámetros propios (hoy solo `scopes` de `budget_threshold`).
- **Rangos (Q2-74, Q2-77, Q2-66 B)**: cada clave del `PATCH` declara su rango (§11.7, columna «Rango (cambio)»): día de la semana de 1 a 7, día del mes de 1 a 31, `comparator.timeout_seconds` de 1 a 30, `price.stale_after_days` de 1 o más, horas con el patrón `HH:MM` de 00:00 a 23:59 y los demás enteros de 0 o más. `price.zero_weight_after_days` ≥ `price.stale_after_days` la valida la API (422) y el motor de precios. Tipo nuevo `SettingParameterTime`.
- **Otras claves (Q2-75, Q2-76)**: `optimizer.date_window_days` (7) en `parameters` de `getOptimizerSettings` y `updateOptimizerSettings`, y deja de ser campo propio de esos dos esquemas (sigue en `ScenarioRequest` y en `constraints_used`); `ocr.timeout_seconds` (30) en `DataSettings` (sección propuesta por S1).
- **IVA de un crédito (Q2-67 opción 1, K46)**: `Credit`, `CreditCreate` y `CreditUpdate` traen `vat_included` (`true` = la tasa ya incluye IVA; `false` = se suma `vat_on_interest_bp`; ausente = falta capturar). `Credit.missing_terms` lista 5 datos: `annual_rate`, `rate_base`, `vat`, `min_payment_pct_bp` y `min_payment_fixed_cents`.
- **Reparto (Q2-68 C)**: `CostShareRuleUpdate.scope` (`onward` o `all`), obligatorio con `parts` y sin valor por omisión (§7.2).
- **Lista de compras (Q2-39 a Q2-42 A, S4 1.5)**: `ShoppingListItem.assigned_store_id` (tienda asignada, rige la agrupación); `recommended_store_id` es la recomendación de S5 por precio y en vertical `market` no rige la agrupación; `estimated_price_cents` es el costo con la cotización que rige; `groups[].store_id` es la tienda asignada (ausente = «sin mercado dado de alta»), `groups[].visit_date` es nuevo, y los totales suman solo los artículos con precio. Orden de los grupos: mercado, cadena y el grupo sin tienda al final.
- **Comparativa (S5 1.5, §7.4)**: `BasketComparison.cheapest_store_id` (primera tienda completa; ausente si ninguna) y `stores[].difference_vs_cheapest_cents` solo en tiendas completas; `stores[]` lleva primero las completas; `identical_count` cuenta renglones.
- **Inicio, metas y comparador**: `AccrualVsCash.months[]` suma `cumulative_accrued_cents` y `cumulative_cash_cents` (Q2-65 A); `SavingsGoal.progress_tenths_pct` (Q2-17 C); `DashboardBestBasket.second_store_id` (Q2-86 A); `Scenario.store_id` (Q2-33 (b)).

### C1 v1.2 (puerta 1b)

- **Entradas**: el generador aplica `00001_esquema.sql`, `00002_puerta1.sql` y `00003_puerta1b.sql` (esquema 1.3). `parse_schema.py` lee además las anotaciones `@param`, `@param_sin_valor` y `@col_note` de las migraciones y comprueba que cada valor inicial anotado sea el que la migración siembra en `app_setting`.
- **Parámetros (K44)**: lectura y cambio **por clave** de `app_setting` en la operación de settings de su sección: `parameters` (valor vigente, valor inicial, tipo y unidad) y, al escribir, `parameters` y `reset_parameters`. «Restablecer» es un campo, no una operación. Detalle y tabla en §11.7.
- **`PriceSettings`**: `stale_after_days` deja de ser un campo propio y pasa a `parameters` (`price.stale_after_days`), junto a `price.zero_weight_after_days` y `price.profeco_import_max_gb`. `PriceSettingsUpdate` ya no exige `assisted_reading_enabled`.
- **Salud (K36)**: `HealthProfile`, `HealthProfileCreate` y `HealthProfileUpdate` traen `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` (ausente o `null` = el de la casa); `HealthProfile` suma `household_goal_adjust_lose_kcal` y `household_goal_adjust_gain_kcal` y `suggested_kcal_target` sigue RN-17 1.3.
- **Valores de K43**: adjunto de 10 MB (`attachment.max_size_mb`, 413 con el máximo vigente en el mensaje), comparador de 5 s (`comparator.timeout_seconds`, 503 `unavailable`), claves `send_time` y `advance_days` de `params_json` y «casillas sin armar», sin marca.
- **Contraseña**: `PasswordChangeRequest.new_password` ya no fija `minLength` 12 en el contrato: el mínimo es el valor vigente de `auth.password_min_length` (12 al sembrar).
- **Otros**: N-289 sale de `updatePerson` (el tope por persona es de N-310, K13); `Notification.params_json` dice que los umbrales y los días antes del pago próximo pasaron a `app_setting`.

### C1 v1.1 (puerta 1)

- **Entradas**: el generador aplica todas las migraciones de S1 en orden (`00001_esquema.sql` y `00002_puerta1.sql`) en una base temporal y lee de ahí tablas, columnas y llaves; los comentarios de columna salen de los dos archivos.
- **Operaciones nuevas**: `saveTicketScanDraft` (N-075, K3), `listTelegramChatCandidates` (N-299) y `confirmTelegramChatCandidate` (N-316, K4), `getProjectionSettings` (N-318) y `updateProjectionSettings` (N-319, K21). N-317 (cola sin conexión) es Local.
- **Operación retirada**: `createNotificationRecipient` (alta manual del chat). K4 eligió que el bot detecte el chat y la persona lo confirme.
- **Presupuesto (K2)**: `updateBudget` recibe `scope` (`this_month` | `onward`, por omisión `onward`) y `month`; `Budget` suma `is_single_month` y `following_month_amount_cents` (N-275). Efecto de cada alcance sobre las vigencias y un ejemplo de cada uno en §7.1.
- **Ticket (K3)**: `TicketScan` trae `has_draft`, `draft` (`TicketScanDraft`), `display_store_id` y `display_total_cents` (N-071, N-076). `saveTicketScanDraft` solo guarda el borrador: no crea gasto, precio ni inventario; solo `confirmTicketScan` guarda precios (decisiones E: el OCR siempre con confirmación humana) y vacía el borrador.
- **Captura sin conexión (K9)**: `createExpense` y `createWeightEntry` exigen `Idempotency-Key`, aceptan `captured_at` opcional y responden **200 con el mismo recurso** al reintento (`Idempotent-Replayed: true`). La clave queda en `expense.idempotency_key` y `weight_entry.idempotency_key`, sin purga, y no se devuelve en las lecturas. Ciclo completo antes y después de las 24 h y ejemplo de reintento en §1 punto 8.
- **Tasas (K12)**: `rate_decimals` en `CreditTermRate`, `CreditPurchase` y `Promotion`, y `annual_rate_decimals` en `Credit`; `rate_pct` y `annual_rate_pct` muestran la tasa como se capturó (la fracción redondeada ROUND_HALF_UP a esos decimales). Los cálculos usan la fracción.
- **Tope por persona (K13)**: `person.fortnight_payment_cap_cents`, opcional, se lee y edita en `listPersons` y `updatePerson`.
- **Restricciones duras (K16)**: `OptimizerSettings` las devuelve aunque falten, con `missing_constraints`; `ScenarioRequest` las sobreescribe solo para esa comparación; `ScenarioResult` suma `status` (`ok`, `missing_constraints`, `no_valid_scenario`), `missing_constraints` y `constraints_used`, y sin restricciones no recomienda.
- **Tarjeta y crédito de tienda (K18)**: `Credit` suma `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents`, `annual_rate_decimals` y el calculado `missing_terms` («falta capturar»).
- **Proyección (K21)**: `CashFlowFortnights`, `AccrualVsCash` y `CashFlowAssumptions` dicen el modo usado (`projection_mode`).
- **Meta de kcal (K23)**: `HealthProfile` suma `goal` y el calculado `suggested_kcal_target`; confirmar la meta = `updateHealthProfile` con `daily_kcal_target`.
- **Familia (K5)**: `family` con tres valores incluyentes. **Lista de compras (R3-01)**: la semana va en la llave `week` (semana ISO, obligatoria; no es el filtro `period`, corrección r2) en `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getDashboardBestBasket` y `getBasketComparison`; el contrato no tiene semana de lista por omisión.
- **Precios (K22)**: `PriceSettings` pierde `stale_weight`; la regla de peso de un precio viejo quedó fija en K22.
- **Marcas**: se quitan las marcas provisionales que cerraron K1, K2, K3, K4, K5, K6, K7, K8, K11, K12, K13, K15, K16, K17, K20, K21 y K22. Las que quedan citan un bloqueo abierto (§13).
- **Extractos**: `build.py` genera `docs/contratos/extractos/` (un archivo por etiqueta, `componentes-comunes.yaml` e `INDICE.md`).

## 1. Convenciones generales

1. **Base**: `https://finanzas.local:8443` (solo red local, decisión H). Recursos bajo `/api/v1`; fuera de él solo `GET /health` y `GET /openapi.yaml`.
2. **Formato**: JSON UTF-8 (`application/json`), salvo adjuntos y PDF (binario) y cargas de archivo (`multipart/form-data`).
3. **Nombres**: rutas, parámetros y propiedades en inglés corto, **iguales al `nombre_codigo` del diccionario de S1**. Descripciones y mensajes en español de México.
4. **Extensiones del contrato**: `x-tablas` (tablas de S1 que lee o escribe cada operación), `x-necesidades` (necesidades N-xxx de U1 que cubre), `x-tabla` (tabla de un esquema de entidad), `x-calculado` (campo calculado por la API o un motor; no es columna), `x-entrada` (dato de una acción que no se guarda, p. ej. la contraseña o la confirmación de una restauración), `x-columna` (columna a la que corresponde una entrada en lista, p. ej. `plan_dates` → `plan_date`), `x-app-setting` (clave de `app_setting`) y `x-provisional` (depende de un bloqueo abierto).
5. **Booleanos**: las columnas 0/1 de SQLite se exponen como `true`/`false`.
6. **Valores nulos**: en las respuestas, una columna NULL **se omite** (campo ausente = sin valor; la pantalla muestra «—», U1 §0.2 «sin datos no es cero»). Los esquemas de lectura marcan como `required` solo las columnas NOT NULL.
7. **Altas** (`POST`) aceptan las columnas escribibles; `id`, `created_at` y `updated_at` los pone el servidor. **Cambios** (`PATCH`) siguen JSON Merge Patch (RFC 7396): solo se cambian los campos enviados y `null` borra un valor que admite NULL. Reemplazos completos de listas (orden de categorías, destinatarios, grupos de producto) usan `PUT` con `{"ids": [...]}`.
8. **Idempotencia**: todo `POST` de captura exige el encabezado `Idempotency-Key` (16 a 128 caracteres `[A-Za-z0-9_-]`, UUID v4 recomendado). Repetir la clave con el mismo cuerpo devuelve la respuesta original con `Idempotent-Replayed: true`; con otro cuerpo, 409 `idempotency_key_reused`. La API guarda cada clave por cuenta en la tabla `idempotency_key` de S1, con la huella SHA-256 del cuerpo y la respuesta (S1 C-15), y la declara en `x-tablas` de cada operación que la exige; un job la purga por antigüedad. Retención de 24 h (K6). Gastos y pesos guardan además su clave en el propio renglón, sin purga, y responden 200 con el mismo recurso al reintento (K9). Lo exigen 41 operaciones. No la llevan: las consultas por `POST` que no guardan nada (`previewCreditPurchase`, `compareScenarios`, `previewPlannedItemsImport`, `suggestRecipeWithAi`), las acciones de sesión (`login`, `logout`, `changeOwnPassword`, `revokeSessions`), el aviso de prueba (`sendTestNotification`) y la restauración, que lleva su propia confirmación escrita (`restoreBackup`).

   **Ciclo de una captura de gasto o de peso** (`createExpense` y `createWeightEntry`; K9 y K6). La cola del teléfono manda cada captura con su propia clave y su `captured_at`:
   - **Primer envío**: 201 con el recurso nuevo y su `id`. La clave queda en la tabla `idempotency_key` (con la huella del cuerpo y la respuesta) y en la columna `idempotency_key` del renglón (`expense` o `weight_entry`).
   - **Reintento con la misma clave y el mismo cuerpo**: 200 con el mismo recurso (mismo `id`) e `Idempotent-Replayed: true`; no se crea otro registro.
   - **Misma clave con otro cuerpo dentro de las 24 h**: 409 `idempotency_key_reused`; no se crea ni se cambia nada.
   - **Pasadas las 24 h** (el job ya purgó la clave de la tabla `idempotency_key`): la clave sigue en el renglón, así que la misma clave responde 200 con el registro guardado, mismo `id` e `Idempotent-Replayed: true`, sin crear otro. Ya no hay huella del cuerpo y solo se compara la clave: un cuerpo distinto ya no se detecta (no hay 409) y el registro guardado no cambia.

   Ejemplo con un gasto de contado de $199.00 (el monto de `parametros.papel` en `oraculo/model.json`), capturado sin red el 10-oct-2026 a las 21:30 y enviado al volver a la red de casa. `id` y `category_id` son identificadores de ejemplo; las respuestas van abreviadas:

   ```http
   POST /api/v1/expenses
   Idempotency-Key: 3f6c2a8e-7b1d-4c9a-9e2f-5a7d1b0c4e81

   {"category_id": 5, "amount_cents": 19900, "expense_date": "2026-10-10", "payment_method": "cash", "captured_at": "2026-10-10T21:30:00-06:00"}
   ```

   - Primer envío → **201** `{"id": 812, "category_id": 5, "amount_cents": 19900, "expense_date": "2026-10-10", "payment_method": "cash", "captured_at": "2026-10-10T21:30:00-06:00", …}`.
   - La cola no recibió la respuesta y reintenta con la misma clave y el mismo cuerpo → **200**, encabezado `Idempotent-Replayed: true` y el mismo cuerpo de respuesta, con `"id": 812`.
   - Misma clave con otro cuerpo (por ejemplo, otro `amount_cents`) dentro de las 24 h → **409** `{"code": "idempotency_key_reused", "message": "…", "request_id": "…"}`.
   - Misma clave después de las 24 h, con el mismo cuerpo o con otro → **200**, `Idempotent-Replayed: true` y el gasto 812 tal como se guardó.

   `createWeightEntry` (`POST /api/v1/weight-entries`) sigue el mismo ciclo con el peso.

9. **Identificador de petición**: toda respuesta lleva `X-Request-Id`; el mismo valor va en `Error.request_id` y en `audit_log.request_id` de los cambios que haga la petición.
10. **Cálculos**: la API no recalcula reglas en el cliente. Las cifras marcadas «motor Sx» las calcula ese motor con las reglas de `reglas-negocio.md`; la interfaz solo muestra, filtra, ordena y suma lo visible (U1 §0.2).
11. **Valores técnicos aprobados (K6)**: 200 registros por página como máximo (`limit`); `Idempotency-Key` retenida 24 h; cookie `SameSite=Strict`; contraseña nueva de 12 caracteres o más; carga de PROFECO de 1 GB como máximo; `cursor` de 512 caracteres como máximo; búsqueda (`q`) de 200 caracteres como máximo; usuario de 3 a 64 caracteres; clave de idempotencia de 16 a 128 caracteres. Desde C1 v1.2, la retención de `Idempotency-Key`, la contraseña mínima y el máximo de PROFECO son parámetros configurables que nacen con esos valores (K44, §11.7).

## 2. Autenticación, sesiones y permisos

- **Esquema**: `cookieAuth`, cookie `fh_session` **HttpOnly; Secure; SameSite=Strict; Path=/**. El token en claro solo viaja en la cookie; la base guarda su SHA-256 (`device_session.token_hash`). Ninguna respuesta incluye hashes, tokens, `chat_id` ni llaves (C1-12). Los chats que detecta el bot se muestran con el nombre de Telegram y el `chat_id` enmascarado (`chat_id_masked`: «•••» y los últimos 4 caracteres; K4).
- **Operaciones públicas**: solo `POST /api/v1/auth/login` y `GET /health` (sin datos). Todo lo demás, incluidos adjuntos, PDF y `GET /openapi.yaml`, exige sesión; sin sesión, vencida o cerrada responde **401 `unauthorized`**.
- **Inicio de sesión** (`login`, N-010): usuario y contraseña (argon2id en la base). Error de credenciales: 401 `invalid_credentials` con el mensaje genérico «Usuario o contraseña incorrectos». **Límite (decisión A): 5 intentos fallidos en 15 minutos** (`app_setting auth.max_failed_attempts = 5`, `auth.lockout_window_minutes = 15`). El intento siguiente responde **429 `too_many_attempts`**, mensaje «Demasiados intentos. Vuelve a intentar a las 21:45.» y `retry_at` (instante con -06:00), **aunque la contraseña sea correcta**; así se cubre también el «estado de bloqueo» (N-011) sin una consulta pública que revele usuarios. Aprobado en la puerta 1 (K11): ventana deslizante de 15 minutos; se cuentan los fallos por usuario y por IP (lo que llegue primero a 5); `retry_at` = hora del fallo más antiguo de la ventana + 15 minutos.
- **Sesión larga por dispositivo** (decisión A, «meses»): 90 días desde el último uso, renovable con cada petición (K11; `device_session.expires_at`; `app_setting session.max_age_days = 90`, sembrado en 00002). `logout` cierra la del dispositivo; Ajustes lista las sesiones de **ambos** usuarios y cualquiera puede cerrarlas (`listSessions`, `revokeSession`, `revokeSessions`, N-271 y N-272). Una sesión cerrada responde 401 en la siguiente petición.
- **Permisos**: Axel y Ximena son co-titulares con **los mismos permisos**; ninguna operación filtra ni oculta datos según quién inició sesión (U1 §0.5). El usuario de la sesión solo da valores por defecto (persona del gasto, del peso) y el «creado por» de `audit_log`. **No existe operación para crear cuentas**: Karen no tiene cuenta ni acceso (decisión A); `person` no se crea desde la API.
- **CSRF**: con `SameSite=Strict` (K6) el navegador no envía la cookie en peticiones que vienen de otro sitio.
- **Contraseña mínima** del cambio de contraseña (N-273): 12 caracteres al sembrar (K6); se usa el valor vigente de `auth.password_min_length` (K44, §11.7).

## 3. Modelo de errores (esquema único)

Todas las respuestas 4xx y 5xx de todas las operaciones usan `components/schemas/Error` mediante `components/responses/*` (más `default` para cualquier otro código):

```json
{"code": "validation_failed", "message": "Revisa los campos marcados.", "fields": [{"field": "amount_cents", "code": "out_of_range", "message": "El monto debe ser mayor que cero."}], "request_id": "01JABCDEF0123456789"}
```

- `code`: estable, en inglés, `^[a-z][a-z0-9_]*$`. `message`: español de México, lenguaje simple (qué pasó y qué hacer). `fields`: detalle por campo (nombre_codigo). `request_id`: siempre. `retry_at`: solo en 429.

| HTTP | Uso | Códigos estables |
|---|---|---|
| 400 | Petición mal formada: JSON inválido, parámetro con formato incorrecto (p. ej. `period=oct-2026`, fecha en formato de pantalla `10-oct-2026`) | `bad_request` |
| 401 | Sin sesión, vencida o cerrada; credenciales incorrectas en el inicio de sesión | `unauthorized`, `invalid_credentials` |
| 403 | Acción prohibida sobre el registro (precio de PROFECO) | `forbidden`, `read_only_source` |
| 404 | El registro o la ruta no existe (también una persona no externa en el estado de cuenta) | `not_found` |
| 405 | Método no admitido en una ruta existente | `method_not_allowed` |
| 409 | Conflicto con el estado: ya existe, en uso, gasto con compra a crédito, receta planeada, casilla con «fuera», vigencia de presupuesto traslapada, IA apagada, lectura asistida sin veredicto, clave de idempotencia reutilizada | `conflict`, `already_exists`, `in_use`, `has_credit_purchase`, `has_payments`, `recipe_planned`, `slot_conflict`, `budget_overlap`, `ai_disabled`, `assisted_reading_not_allowed`, `idempotency_key_reused` |
| 413 | Archivo demasiado grande | `payload_too_large` |
| 415 | Tipo de archivo no admitido | `unsupported_media_type` |
| 422 | Validación de negocio o de campos (con `fields`) | `validation_failed`, `invalid_current_password`, `not_revolving`, `not_a_card`, `external_person` |
| 429 | Límite de intentos de inicio de sesión (con `retry_at`) | `too_many_attempts` |
| 500 / 501 / 503 | Error interno / operación aún no implementada (esqueleto WF-0) / servicio o archivo no disponible | `internal_error`, `not_implemented`, `unavailable` |

**Validación de dinero**: un importe que no es entero (`973.75`, `"973.75"`) es 400 `bad_request`; un entero fuera de regla (cero o negativo donde se pide mayor que 0, partes que no suman el total) es 422 `validation_failed` con `fields[].code = out_of_range` o `invalid_sum`. **Validación de fechas**: una fecha que no es `AAAA-MM-DD` válida (`2026-02-30`, `10-oct-2026`) es 400; una fecha válida fuera de regla (pago antes de la compra, fin antes del inicio) es 422.

## 4. Fechas y dinero

- **Fechas de negocio** (`Date`, sufijo `_date`, y `valid_from`/`valid_to`): ISO 8601 `2026-10-10`, calendario de `America/Mexico_City`. **Instantes** (`DateTime`, sufijo `_at`): ISO 8601 con desplazamiento explícito de `America/Mexico_City`, hoy `-06:00` (`2026-10-10T21:30:00-06:00`); la base los guarda en UTC con «Z» y la API los convierte. **«Z» solo se acepta en las entradas** (R2-06, K12): la API recibe `2026-10-11T03:30:00Z` o `2026-10-10T21:30:00-06:00` y responde siempre con el desplazamiento de `America/Mexico_City`. **Mes** (`Month`): `2026-10`. Ningún formato de pantalla (`10-oct-2026`, `Vie 30 oct 2026`) aparece en la API; la interfaz los forma.
- **Quincenas**: los días de pago son el **15 y el 30**; en meses más cortos el «30» es el **último día** (en febrero, el 28: `2027-02-28`, R-23). Rangos para agrupar: Q1 del 1 al 15 y Q2 del 16 al último día; el 31 cae en Q2 (K8, firme). **Semana**: lunes a domingo (K8).
- **Dinero** (`MoneyCents`): entero exacto de **centavos MXN**, sufijo `_cents`, nunca `number` de punto flotante. `97375` = $973.75. Redondeo **ROUND_HALF_UP**; en una serie de pagos **el último pago ajusta centavos** (5,034.74 → 1,258.68 ×3 + 1,258.70, R-18 a R-21). Los dos modos de base de la serie que usa el oráculo y su generalización están en `S1-redondeo-series` (las series se guardan ya calculadas, el contrato solo las expone). En pantalla: `$`, separador de miles y dos decimales.
- **Porcentajes** (`BasisPoints`, sufijo `_bp`): 10000 = 100 %. Los porcentajes de avance calculados (`used_pct`, `progress_pct`) son enteros ROUND_HALF_UP. **Tasas**: fracción exacta `rate_num`/`rate_den` (S1; la tasa redondeada no reproduce el oráculo) y, para mostrar, `rate_pct` en texto **como se capturó** (K12): la fracción redondeada ROUND_HALF_UP a `rate_decimals` decimales (`17.27`, `25.8685`). **Cantidades** (`Milli`, sufijo `_milli`): milésimas de la unidad base (g, ml, pieza) o de porción.

## 5. Paginación y orden (una sola convención)

- **Colección** = lista de registros de una entidad. Toda colección (`40` operaciones `GET`) usa **los mismos parámetros**: `cursor` (opaco; sin valor = primera página), `limit` (por defecto **50**, como pide U1; máximo **200**, K6) y `sort` (propiedades separadas por coma, `-` = descendente; el `id` cierra siempre el orden para que el cursor sea estable). **Misma envoltura** `Page`: `{"items": [...], "next_cursor": "…" | null}`; **`next_cursor: null` indica que no hay más**. El orden por defecto de cada colección está en su descripción (tomado de U1 §2.5.b).
- **Agregados y métricas** (tarjetas, resúmenes, series de 12 meses o 24 quincenas, comparativas) son objetos calculados con tamaño acotado; no se paginan.

## 6. Filtros globales (definición única)

Los seis filtros del plan (§3.7) tienen **una sola definición** en `components/parameters` y todas las operaciones los usan por `$ref`:

| Filtro | Parámetro | Tipo y formato | Significado |
|---|---|---|---|
| persona | `person_id` | ids separados por coma | `person.id`; en gastos, `expense.person_id`; en partes de dinero, la persona de la parte |
| familia | `family` | `household` · `external` · `joint` | Tres valores incluyentes (K5): `household` = Hogar, reparto Hogar o Conjunto; `external` = Karen, reparto Karen o Conjunto; `joint` = solo Conjunto |
| categoría | `category_id` | ids separados por coma | `category.id` |
| tienda | `store_id` | ids separados por coma | `store.id` (en promociones, el comercio ligado) |
| periodo | `period` | `2026-10` mes · `2026-10-q1`/`q2` quincena · `2026-W42` semana · `2026` año · `2026-10-10..2026-11-09` rango | mes calendario como base; quincena del 1 al 15 y del 16 al último día; semana de lunes a domingo (K8). La lista de compras no usa este filtro: su semana va en la llave `week` (abajo) |
| crédito | `credit_id` | ids separados por coma | `credit.id`; en gastos, la forma de pago con esa línea |

Además se reutilizan `q` (texto) y los parámetros de paginación. Los filtros propios de una entidad (`status`, `kind`, `source`, etc.) se documentan en cada operación.

**Llaves que no son filtros.** Dos parámetros hablan de un periodo pero no filtran: dicen sobre qué recurso actúa la operación. Cada uno tiene una sola definición en `components/parameters`, su nombre no es el de ningún filtro y no acepta los demás formatos de `period`:

| Llave | Componente | Tipo y formato | Operaciones | Significado |
|---|---|---|---|---|
| `week` | `ListWeek` | semana ISO `2026-W42` (lunes a domingo, K8); obligatoria, sin valor por omisión | `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getDashboardBestBasket`, `getBasketComparison` | La lista de compras semanal que se lee o se escribe: la que muestra la pantalla (R3-01). La semana que se abre por defecto la decide la PWA (U1 §0.2) |
| `month` | `BudgetMonth` | `AAAA-MM`; sin valor, el mes en curso | `updateBudget` | El mes al que aplica un cambio de monto (§7.1, K2) |

**Matriz recurso × filtro** (operaciones que admiten cada filtro; la última columna marca la llave `week` de la lista de compras, que no es filtro):

| Operación | Ruta | persona | familia | categoría | tienda | periodo | crédito | semana de lista (`week`, llave) |
|---|---|---|---|---|---|---|---|---|
| `listSessions` | GET `/api/v1/sessions` | ✓ |  |  |  |  |  |  |
| `listCalendarEvents` | GET `/api/v1/calendar/events` |  | ✓ |  |  | ✓ | ✓ |  |
| `queryView` | GET `/api/v1/views/{entity}` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |  |
| `getDashboardMonthVsPlan` | GET `/api/v1/dashboard/month-vs-plan` |  |  |  |  | ✓ |  |  |
| `getDashboardFortnightFlow` | GET `/api/v1/dashboard/fortnight-flow` |  |  |  |  | ✓ |  |  |
| `getDashboardWeekMeals` | GET `/api/v1/dashboard/week-meals` |  |  |  |  | ✓ |  |  |
| `getDashboardBestBasket` | GET `/api/v1/dashboard/best-basket` |  |  |  |  |  |  | ✓ |
| `listExpenses` | GET `/api/v1/expenses` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |  |
| `getExpenseTotals` | GET `/api/v1/expenses/totals` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |  |
| `getExpenseSummary` | GET `/api/v1/expenses/summary` | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |  |
| `getExpenseMonthlyTrend` | GET `/api/v1/expenses/monthly-trend` | ✓ | ✓ |  |  |  | ✓ |  |
| `listPeriods` | GET `/api/v1/periods` |  |  |  |  | ✓ |  |  |
| `listBudgets` | GET `/api/v1/budgets` |  |  | ✓ |  | ✓ |  |  |
| `getBudgetMonth` | GET `/api/v1/budget/month` |  |  |  |  | ✓ |  |  |
| `getBudgetByCategory` | GET `/api/v1/budget/categories` |  |  |  |  | ✓ |  |  |
| `getBudgetFortnights` | GET `/api/v1/budget/fortnights` |  |  |  |  | ✓ |  |  |
| `getAccruedGroceries` | GET `/api/v1/budget/accrued-groceries` |  |  |  |  | ✓ |  |  |
| `getBudgetMonthlyTrend` | GET `/api/v1/budget/monthly-trend` |  |  | ✓ |  |  |  |  |
| `listIncomes` | GET `/api/v1/incomes` |  |  |  |  | ✓ |  |  |
| `getIncomeMonth` | GET `/api/v1/incomes/month` |  |  |  |  | ✓ |  |  |
| `listPlannedItems` | GET `/api/v1/planned-items` |  |  | ✓ |  |  | ✓ |  |
| `listSavingsMovements` | GET `/api/v1/savings-movements` |  |  |  |  | ✓ |  |  |
| `getMealPlanWeek` | GET `/api/v1/meal-plan/week` | ✓ |  |  |  | ✓ |  |  |
| `getMealPlanMonth` | GET `/api/v1/meal-plan/month` |  |  |  |  | ✓ |  |  |
| `getMealPlanTrends` | GET `/api/v1/meal-plan/trends` | ✓ |  |  |  |  |  |  |
| `copyPreviousWeek` | POST `/api/v1/meal-plan/copy-previous-week` |  |  |  |  | ✓ |  |  |
| `getShoppingList` | GET `/api/v1/shopping-list` |  |  |  | ✓ |  |  | ✓ |
| `generateShoppingList` | POST `/api/v1/shopping-list/generate` |  |  |  |  |  |  | ✓ |
| `createShoppingListItem` | POST `/api/v1/shopping-list/items` |  |  |  |  |  |  | ✓ |
| `updateShoppingListItem` | PATCH `/api/v1/shopping-list/items/{shopping_list_item_id}` |  |  |  |  |  |  | ✓ |
| `assignShoppingListStore` | POST `/api/v1/shopping-list/assign-store` |  |  |  |  |  |  | ✓ |
| `listBaskets` | GET `/api/v1/baskets` |  |  |  |  |  |  | ✓ |
| `getBasketComparison` | GET `/api/v1/basket-comparison` |  |  |  |  |  |  | ✓ |
| `listProducts` | GET `/api/v1/products` |  |  |  | ✓ |  |  |  |
| `listProductPrices` | GET `/api/v1/products/{product_id}/prices` |  |  |  | ✓ | ✓ |  |  |
| `getDebtTrend` | GET `/api/v1/credits/debt-trend` |  |  |  |  |  | ✓ |  |
| `listCreditPurchases` | GET `/api/v1/credit-purchases` | ✓ | ✓ |  |  | ✓ | ✓ |  |
| `listScheduledPayments` | GET `/api/v1/scheduled-payments` | ✓ | ✓ |  |  | ✓ | ✓ |  |
| `listPromotions` | GET `/api/v1/promotions` |  |  |  | ✓ | ✓ | ✓ |  |
| `listWeightEntries` | GET `/api/v1/weight-entries` | ✓ |  |  |  | ✓ |  |  |
| `listReports` | GET `/api/v1/reports` | ✓ |  |  |  | ✓ |  |  |
| `getMonthlyReport` | GET `/api/v1/monthly-report` |  |  |  |  | ✓ |  |  |
| `getExternalStatement` | GET `/api/v1/persons/{person_id}/statement` |  |  |  |  | ✓ |  |  |
| `listNotificationDeliveries` | GET `/api/v1/notification-deliveries` |  |  |  |  | ✓ |  |  |
| `listBackups` | GET `/api/v1/backups` |  |  |  |  | ✓ |  |  |

## 7. Recursos: las 31 entidades de plan §3.4.1

Cada entidad tiene su recurso (lectura y, cuando el usuario la edita, escritura) o declara qué motor u operación la usa. Las vistas siguen U1 §2.5.a y se piden con `GET /api/v1/views/{entity}` (§8).

| # | Entidad (§3.4.1) | Tabla(s) S1 | Operaciones | Vistas (U1 §2.5.a; * = por defecto) |
|---|---|---|---|---|
| 1 | persona | `person` | listPersons, updatePerson, getPersonParticipation | Tabla* (P-26.12) |
| 2 | cuenta_usuario | `user_account` | listUsers, login, getCurrentSession, changeOwnPassword (sin alta: Karen no tiene cuenta) | Tabla* (P-26.01) |
| 3 | sesion_dispositivo | `device_session` | listSessions, revokeSession, revokeSessions, logout | Tabla* (P-26.01) |
| 4 | periodo | `period` | listPeriods (selector); motor S2 (meses y quincenas de los cálculos; el presupuesto usa vigencia por fechas) | Sin lista: dato derivado (U1) |
| 5 | ingreso | `income + planned_item` | listIncomes, getIncomeMonth, createIncome, updateIncome, deleteIncome, getIncomeMonthlyTrend; listPlannedItems, createPlannedItem, updatePlannedItem | Tabla*, calendario (income_date) |
| 6 | categoria | `category` | listCategories, createCategory, updateCategory, reorderCategories | Tabla*, tablero (default_cost_share_rule_id) |
| 7 | gasto | `expense + expense_item + expense_share + attachment` | listExpenses, getExpense, createExpense, updateExpense, deleteExpense, duplicateExpense, markExpensePaid, listExpenseItems, uploadExpenseAttachment, downloadAttachment, deleteAttachment, ticket-scans | Tabla*, tablero (category_id), calendario (expense_date), galería (tickets) |
| 8 | regla_reparto | `cost_share_rule + cost_share_rule_part` | listCostShareRules, updateCostShareRule; reglas de un solo uso al crear gasto o renglón | Sin lista: propiedad de la categoría y del gasto (U1) |
| 9 | meta_ahorro | `savings_goal` | listSavingsGoals, createSavingsGoal, updateSavingsGoal, getSavingsSummary | Tabla (P-26.13), tablero (status), calendario (target_date), galería* (P-09) |
| 10 | movimiento_ahorro | `savings_movement` | listSavingsMovements, createSavingsMovement, updateSavingsMovement, deleteSavingsMovement | Tabla*, calendario (movement_date) |
| 11 | credito | `credit + credit_statement` | listCredits, getCredit, createCredit, updateCredit, getCreditTotals, previewAvailability, getAvailabilitySeries, getCardStatus, getDebtTrend | Tabla (P-26.05), tablero (kind), galería* (P-17) |
| 12 | tasa_por_plazo | `credit_term_rate` | listCreditTermRates, createCreditTermRate, updateCreditTermRate, deleteCreditTermRate | Tabla* |
| 13 | compra_credito | `credit_purchase + credit_purchase_share` | listCreditPurchases, getCreditPurchase, createCreditPurchase, previewCreditPurchase, updateCreditPurchase, deleteCreditPurchase | Tabla*, tablero (status), calendario (purchase_date) |
| 14 | pago_programado | `scheduled_payment + scheduled_payment_share + credit_payment` | listScheduledPayments, markScheduledPaymentPaid, deleteCreditPayment, listCalendarEvents | Tabla, tablero (status), calendario* (due_date) |
| 15 | promocion | `promotion` | listPromotions, createPromotion, updatePromotion, deletePromotion | Tabla, tablero (credit_id), calendario* (valid_from a valid_to) |
| 16 | tienda | `store + store_product_group + product_group` | listStores, createStore, updateStore, setStoreProductGroups, listProductGroups | Tabla*, tablero (kind), calendario (visit_weekday, recurrente), galería |
| 17 | producto | `product` | listProducts, getProduct, createProduct, updateProduct | Tabla*, tablero (freshness), galería |
| 18 | presentacion | `presentation` | createPresentation, updatePresentation; dentro de getProduct | Sin lista propia: tabla dentro del producto (U1) |
| 19 | precio | `price` | listProductPrices, createPrice, updatePrice, deletePrice, getPriceSummary, importProfeco | Tabla*, tablero (freshness) |
| 20 | receta | `recipe + recipe_meal_slot` | listRecipes, getRecipe, createRecipe, updateRecipe, duplicateRecipe, deleteRecipe, getRecipeUsage, suggestRecipeWithAi | Tabla, tablero (tiempos aptos, recipe_meal_slot), galería* |
| 21 | ingrediente_receta | `recipe_ingredient` | dentro de getRecipe, createRecipe y updateRecipe | Sin lista propia: tabla dentro de la receta (U1) |
| 22 | tiempo_comida | `meal_slot` | listMealSlots (solo lectura) | Sin lista: catálogo fijo de 4 (U1) |
| 23 | plan_comida | `meal_plan_entry` | getMealPlanWeek, getMealPlanMonth, getMealPlanTrends, createMealPlanEntries, updateMealPlanEntry, deleteMealPlanEntry, copyPreviousWeek | Tabla, tablero (meal_slot_id), calendario* (plan_date), galería (recipe_id) |
| 24 | inventario_item | `inventory_item + inventory_movement` | listInventoryItems, createInventoryItem, updateInventoryItem, adjustInventoryItem, listInventoryMovements, getInventorySummary | Tabla*, tablero (status), galería |
| 25 | perfil_salud | `health_profile` | listHealthProfiles, createHealthProfile, updateHealthProfile, getTodayKcal | Sin lista: una ficha por persona (U1) |
| 26 | peso_historico | `weight_entry` | listWeightEntries, createWeightEntry, updateWeightEntry, deleteWeightEntry | Tabla* |
| 27 | nutricion_ingrediente | `nutrition_fact` | listNutritionFacts, createNutritionFact, reloadNutritionBase, getRecipeSettings | Sin lista: propiedad del producto (U1) |
| 28 | notificacion | `notification + notification_delivery` | listNotifications, updateNotification, sendTestNotification, listNotificationDeliveries | Tabla*, tablero (notification_id), calendario (sent_at) en avisos enviados |
| 29 | destinatario | `notification_recipient + notification_subscription + telegram_chat_candidate` | listNotificationRecipients, listTelegramChatCandidates, confirmTelegramChatCandidate, updateNotificationRecipient, setNotificationRecipients | Sin lista: se elige dentro de cada aviso (U1) |
| 30 | job | `job` | listJobs, updateJob | Tabla*, tablero (code), calendario (next_run_at) |
| 31 | respaldo | `backup` | listBackups, createBackup, restoreBackup | Tabla* |

### 7.1 Presupuesto: alcance de un cambio de monto (`scope`, K2)

`updateBudget` (`PATCH /api/v1/budgets/{budget_id}`, N-087 y N-278) recibe dos parámetros de consulta cuando cambia `amount_cents`: `month` (AAAA-MM, el mes que muestra la pantalla; sin valor, el mes en curso), que debe caer en la vigencia del renglón, y `scope`, con dos valores. La pantalla preselecciona «de este mes en adelante».

| `scope` | En pantalla | Efecto sobre los renglones de vigencia (`valid_from_date`, `valid_to_date`) | Responde |
|---|---|---|---|
| `onward` (por omisión) | «de este mes en adelante» | El renglón termina el último día del mes anterior a `month`; un renglón nuevo con el monto nuevo rige desde el día 1 de `month` sin fin (`valid_to_date` NULL); los renglones posteriores de la misma categoría (o del total) se eliminan porque el monto nuevo los reemplaza (quedan en `audit_log`) | El renglón nuevo |
| `this_month` | «solo este mes» | El renglón termina el último día del mes anterior a `month`; un renglón nuevo con el monto nuevo rige solo en `month` (`is_single_month` = true) y otro renglón nuevo con el monto anterior rige desde el mes siguiente hasta el fin original; los renglones posteriores no cambian | El renglón de `month` |

Si el renglón empieza en `month`, el monto cambia en ese mismo renglón (con `onward` queda sin fin y también se eliminan los posteriores; con `this_month`, el resto de su vigencia pasa a un renglón nuevo con el monto anterior). Así, con `onward` todos los meses desde `month` tienen el monto nuevo, como dice K2 («de este mes en adelante») y como lo guarda S1 (`valid_to_date` NULL = rige hasta que se cambie), y nunca queda un traslape. `createBudget` solo llena meses sin monto y no elimina renglones. En los dos casos la respuesta es el renglón que rige en `month`, con `following_month_amount_cents` = el monto que rige el mes siguiente. Los cambios de `valid_to_date` o `notes` no usan `scope` ni `month`. Un traslape con otro renglón de la misma categoría responde 409 `budget_overlap`; `month` fuera de la vigencia, 422 `out_of_range`.

**Ejemplo.** Los montos se toman del oráculo y no son presupuestos capturados: 1,947.50 y 2,517.37 son el devengado de despensa de P1 y de P3 en `oraculo/model.json` (`periodos[0]` y `periodos[2]`, `devengado.despensa_total`). `id` y `category_id` son identificadores de ejemplo. Antes del cambio, la categoría 5 tiene un solo renglón:

| `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
|---|---|---|---|
| 7 | 194750 | 2026-10-01 | — (sin fin) |

En P-06, con diciembre a la vista, se cambia el monto a 2,517.37: `PATCH /api/v1/budgets/7?month=2026-12&scope=…` con el cuerpo `{"amount_cents": 251737}`.

- **`scope=onward`** (o sin `scope`). Después:

  | `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
  |---|---|---|---|
  | 7 | 194750 | 2026-10-01 | 2026-11-30 |
  | 8 (nuevo) | 251737 | 2026-12-01 | — (sin fin) |

  Respuesta **200** (abreviada; `valid_to_date` se omite porque es NULL, §1 punto 6): `{"id": 8, "category_id": 5, "amount_cents": 251737, "valid_from_date": "2026-12-01", "is_single_month": false, "following_month_amount_cents": 251737, …}`.

- **`scope=this_month`**. Después:

  | `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
  |---|---|---|---|
  | 7 | 194750 | 2026-10-01 | 2026-11-30 |
  | 8 (nuevo) | 251737 | 2026-12-01 | 2026-12-31 |
  | 9 (nuevo) | 194750 | 2027-01-01 | — (sin fin) |

  Respuesta **200** (abreviada): `{"id": 8, "category_id": 5, "amount_cents": 251737, "valid_from_date": "2026-12-01", "valid_to_date": "2026-12-31", "is_single_month": true, "following_month_amount_cents": 194750, …}`.

- **`scope=onward` con renglones posteriores** (sigue del resultado de `this_month`: el 7 de octubre a noviembre, el 8 solo en diciembre y el 9 desde enero). Con noviembre a la vista se cambia el monto a 2,517.37 «de este mes en adelante»: `PATCH /api/v1/budgets/7?month=2026-11&scope=onward` con el cuerpo `{"amount_cents": 251737}`. Después:

  | `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
  |---|---|---|---|
  | 7 | 194750 | 2026-10-01 | 2026-10-31 |
  | 10 (nuevo) | 251737 | 2026-11-01 | — (sin fin) |

  El 8 (solo diciembre) y el 9 (desde enero) se eliminan: si se quedaran, enero seguiría en 1,947.50 y el monto no regiría «en adelante», y el 10 sin fin se traslaparía con los dos (409). Respuesta **200** (abreviada): `{"id": 10, "category_id": 5, "amount_cents": 251737, "valid_from_date": "2026-11-01", "is_single_month": false, "following_month_amount_cents": 251737, …}`.

### 7.2 Reparto: alcance de un cambio de partes (`scope`, Q2-68 C)

Al cambiar las partes de una regla de reparto (`updateCostShareRule` con `parts`) la pantalla pregunta el alcance, como el presupuesto en K2 (RN-03 1.4; P-26.12). `CostShareRuleUpdate.scope` es obligatorio cuando viene `parts`, no tiene valor por omisión y no se acepta sin `parts` (422 `validation_failed` en los dos casos). Cambiar solo `name` o `notes` no lleva `scope`.

| `scope` | En pantalla | Qué hace |
|---|---|---|
| `onward` | «Solo los nuevos» | Cambia las partes de la regla. Los gastos ya guardados conservan sus partes, que están congeladas en `expense_share` (S1 1.4); los que se capturen después usan las nuevas. |
| `all` | «También los ya capturados» | Cambia las partes y el motor de reparto recalcula con ellas el `expense_share` de los gastos ya capturados, en la misma transacción, y `audit_log` lo registra. Qué gastos abarca (los que usan la regla en el gasto o en un renglón, o también los que la toman por su categoría o por ser la de por defecto) es PROVISIONAL: bloqueo `C1-reparto-alcance`. |

Ejemplo ilustrativo (cifras sintéticas): «Conjunto 50/50» pasa a 60/40. Un gasto ya guardado de $100.00 con esa regla tiene 5000 y 5000 en `expense_share`. Con `onward` se queda así y un gasto nuevo de $100.00 queda en 6000 y 4000; con `all`, si el gasto queda dentro del alcance, pasa a 6000 y 4000.

## 8. Vistas de base de datos (tabla, tablero, calendario, galería)

`GET /api/v1/views/{entity}?mode=table|board|calendar|gallery` (decisión G) aplica los filtros globales, `q`, `status`, `sort` y la paginación únicos, y responde `ViewPage` = `Page` + `mode` + `groups`. En **tablero**, `group_by` elige la columna y cada grupo trae `count`, `total_cents` (Σ del importe principal) y `share_bp`; en **calendario**, `date_field` elige la fecha y los grupos son días; en **galería**, cada registro trae `image_attachment_id` cuando hay foto. La vista elegida, filtros y columnas visibles se recuerdan en el dispositivo (U1 §2.5 y su supuesto 2); no hay vistas guardadas compartidas en esta versión.

| `entity` | Modos | `group_by` permitidos | `date_field` | Importe principal (`total_cents`) | Necesidades |
|---|---|---|---|---|---|
| `expense` | table, board, calendar, gallery | category_id, kind, payment_method, person_id | expense_date | parte del hogar | N-043, N-045, N-046 |
| `income` | table, calendar | — | income_date | amount_cents | N-104 (también listCalendarEvents) |
| `category` | table, board | default_cost_share_rule_id | — | — | N-281 |
| `savings_goal` | table, board, calendar, gallery | status | target_date | saved_cents | N-111 |
| `savings_movement` | table, calendar | savings_goal_id, kind | movement_date | amount_cents (aportación +, retiro −) | N-115 |
| `credit` | table, board, gallery | kind | — | live_debt_cents | N-190 |
| `credit_term_rate` | table | — | — | — | N-286 |
| `credit_purchase` | table, board, calendar | status, credit_id | purchase_date | total_cents | N-202 |
| `scheduled_payment` | table, board, calendar | status, credit_id | due_date | amount_cents | N-192 |
| `promotion` | table, board, calendar | credit_id | valid_from (franja hasta valid_to) | — | N-287 |
| `store` | table, board, calendar, gallery | kind | visit_weekday (recurrente) | — | N-290 |
| `product` | table, board, gallery | freshness, product_group_id | — | — | N-170 |
| `price` | table, board | is_stale, source, store_id | observed_date | — | N-171 |
| `recipe` | table, board, gallery | meal_slot_id (tiempos aptos de recipe_meal_slot: una receta con varios tiempos aparece en cada columna), origin | — | — | N-135 |
| `meal_plan_entry` | table, board, calendar, gallery | meal_slot_id, person_id | plan_date | — | N-128 |
| `inventory_item` | table, board, gallery | status, tracking_mode | — | — | N-180 |
| `notification_delivery` | table, board, calendar | notification_id, status | sent_at (día local) | — | N-247 |
| `job` | table, board, calendar | code, last_status | next_run_at (día local) | — | N-246, N-302 |
| `person` | table | — | — | — | N-309 |
| `user_account` | table | — | — | — | N-270 |
| `device_session` | table | user_account_id | — | — | N-271 |
| `weight_entry` | table | person_id | measured_date | — | N-232 |
| `backup` | table | status, kind | started_at (día local) | — | N-302 |

Otro valor de `group_by` o `date_field` responde 422 `validation_failed`. `queryView` declara en `x-necesidades` las necesidades cuya vista principal es de tablero, calendario o galería; las vistas de tabla de cada entidad también se pueden pedir aquí, aunque cada pantalla usa la colección tipada de su recurso.

## 9. Métricas del dashboard (fórmula, periodo, filtros, tablas, endpoint, ejemplo del oráculo)

Las 6 tarjetas de Inicio › Resumen (U1 N-020 a N-025, U3) y las métricas de sus subsecciones. Periodo por defecto: **mes calendario**; la quincena es subperiodo. «Parte del hogar» = parte a nombre del titular; toda cifra de presupuesto, flujo y ahorro usa la parte del hogar (U1 §0.2). Los ejemplos usan cifras del oráculo (`oraculo/model.json`, identificadores R-nn de `docs/criterios/S3.md`) o un marcador cuando el oráculo no trae el dato. En la API van en centavos.

| Id | Tarjeta o subsección | Fórmula | Periodo | Filtros | Tablas de origen | Endpoint | Ejemplo (oráculo) |
|---|---|---|---|---|---|---|---|
| M-01 | Mes vs plan (N-020) | real = Σ parte del hogar devengada en el mes (motor S2); plan = Σ `budget.amount_cents` del renglón vigente en el mes (RN-22) de las categorías activas; % usado = real ÷ plan × 100 entero ROUND_HALF_UP; restante = plan − real; estado normal < 80, aviso ≥ 80, excedido ≥ 100; días que faltan contando hoy | mes en curso (`period`) | `period` | budget, category, period, expense, expense_share, credit_purchase | GET `/api/v1/dashboard/month-vs-plan` | oct-2026: real (parte del hogar devengada) **1,197.64** = `periodos[0].devengado.axel` (P1 = oct-2026, hallazgo 4 de S1); plan: marcador (lo captura Axel) |
| M-02 | Flujo quincenal (N-021) | ingresos = quincena capturada (o esperada, marcada) + extras de la quincena; egresos hechos = efectivo pagado hasta hoy (contado + partes del hogar pagadas); comprometido = partes del hogar pendientes de la quincena + fijos y venideros sin pagar; saldo al cierre = ingresos − egresos − comprometido | quincena en curso | `period` (`2026-10-q2`) | income, planned_item, expense, expense_share, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/dashboard/fortnight-flow` | Q2 oct-2026: comprometido incluye la parte del hogar del pago del 30-oct **598.82** (R-13); ingreso: marcador |
| M-03 | Próximo pago (N-022 a) | pago con alguna parte pendiente y fecha mínima ≥ hoy en todas las líneas; si hay pendientes con fecha < hoy, primero el más antiguo («sin marcar como pagado»); pagos del mismo día se suman | hoy en adelante | — | scheduled_payment, scheduled_payment_share, credit_payment, credit | GET `/api/v1/dashboard/next-payment` | **30-oct-2026, 973.75** (hogar 598.82, Karen 374.93) (R-13) |
| M-04 | Disponible A/B (N-022 b, N-201, N-285) | disponible de hoy de cada línea revolvente con la lectura activa (B por defecto) y aviso «pendiente confirmar»; vista previa con ambas lecturas y el mínimo del horizonte (motor S3) | hoy y horizonte | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/dashboard/next-payment`, GET `/api/v1/credits/{credit_id}/availability-preview` | 10-oct-2026: **A 10,000.00 · B 6,905.50** (R-64); mínimos **8,094.50 · 5,000.00** (R-61, R-62) |
| M-05 | Comidas de la semana (N-023) | costo = Σ costo de casillas con receta («fuera» no suma); kcal promedio por día = Σ kcal de la persona ÷ días con al menos un tiempo asignado; casillas sin armar = pares día × tiempo obligatorio sin receta ni «fuera» para alguna persona (máx. 21; K43) | semana en curso (lunes a domingo, K8) | `period` (`2026-W42`) | meal_plan_entry, meal_slot, recipe, recipe_ingredient, nutrition_fact, price, presentation, health_profile | GET `/api/v1/dashboard/week-meals` | marcador (el oráculo no trae comidas) |
| M-06 | Mejor canasta (N-024) | lista de compras vigente (o la última compra con renglones) comparada por tienda (motor S5); ganadora = menor total; diferencia contra la segunda; cobertura = artículos con precio ÷ artículos; precios viejos usados; Sam's como si tuvieran membresía | lista vigente | — | shopping_list, shopping_list_item, meal_plan_entry, recipe_ingredient, inventory_item, price, presentation, product, store, expense, expense_item | GET `/api/v1/dashboard/best-basket` | compra semilla 10-oct-2026: Bodega Aurrera **3,094.50** (R-01) contra Walmart 3,233.00 (`data.json totales.walmart`), diferencia **138.50** (plan §3.4 fila S5) |
| M-07 | Ahorro (N-025) | ahorrado = Σ aportaciones − Σ retiros (fondo y metas); aportado del mes contra planificado del mes (Σ `planned_monthly_cents` activos); meta más próxima con % = ahorrado ÷ monto meta × 100 entero | mes en curso | — | savings_goal, savings_movement | GET `/api/v1/dashboard/savings` | marcador (el oráculo no trae ahorro) |
| M-08 | Participación de Karen (N-026) | su parte devengada de la despensa del mes; sus compras vivas; sus partes de pagos del mes con estado; deuda viva = Σ partes pendientes | mes en curso | — | person, expense, expense_share, credit_purchase, credit_purchase_share, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/persons/{person_id}/participation` | oct-2026: parte devengada **749.86** (`periodos[0].devengado.karen`); su parte del pago del 30-oct **374.93** (R-11) |
| M-09 | Presupuesto por categoría (N-081, N-027) | por categoría: plan, real (parte del hogar devengada), plan − real, % usado y estado (umbrales de M-01) | mes (`period`) | `period` | budget, category, period, expense, expense_share | GET `/api/v1/budget/categories` | marcador (plan por categoría lo captura Axel) |
| M-10 | Flujo por quincena (N-090) | por quincena (24 hacia adelante): ingresos, partes del hogar de pagos, fijos y venideros, variable estimado (método de S2), ahorro planificado, saldo y saldo acumulado; efectivo | desde Q1 del mes en curso | — | income, planned_item, expense, expense_share, scheduled_payment, scheduled_payment_share, credit_payment, savings_goal | GET `/api/v1/cash-flow/fortnights` | Q2 oct-2026: pagos de crédito, parte del hogar **598.82** (R-13) |
| M-11 | Devengado contra efectivo (N-091, N-028) | por mes: devengado y efectivo (totales y partes); al cierre de cada mes, con el mismo origen para los tres (desde el primer evento de las entradas, S2-interfaz §3.2; Q2-25 A y Q2-65 A): devengado acumulado `cumulative_accrued_cents` = Σ devengado desde el primer evento hasta el cierre del mes, efectivo acumulado `cumulative_cash_cents` = Σ efectivo en el mismo tramo, y diferencia acumulada `cumulative_difference_cents` = devengado acumulado − efectivo acumulado (no se rotula deuda viva; `criterios-invariante-deuda`). El origen no es el primer mes visible: los meses anteriores a la ventana cuentan, así que la cima de T-2 es el devengado acumulado sea cual sea el mes en curso | mes en curso a +11 | — | expense, expense_share, planned_item, scheduled_payment, scheduled_payment_share, credit_purchase | GET `/api/v1/cash-flow/accrual-vs-cash` | oct-2026 (P1): devengado **1,947.50**, efectivo **973.75**, acumulados 1,947.50 y 973.75, diferencia acumulada 973.75 (R-45). Ventana que empieza en nov-2026 (reloj en nov-2026): su primer mes (P2) trae devengado 1,947.50 y efectivo 1,947.50 (R-46), y acumulados **3,895.00** y **2,921.25** (oct + nov; diferencia 973.75); sumar solo lo visible daría 1,947.50 y 1,947.50. mar-2027 (P6): acumulados **14,362.48** y **13,103.78** = totales P1–P6 (R-51, R-52), diferencia acumulada 1,258.70 (R-50). Cifras sumadas de `periodos[i]` de `oraculo/model.json` (Decimal) |
| M-12 | Deuda viva a una fecha (N-092) | Σ de lo que queda por pagar después de la fecha (capital + intereses) de las compras hechas hasta esa fecha; total, por línea y por compra (motor S3) | fecha elegida | `date` | credit_purchase, scheduled_payment, credit | GET `/api/v1/debt/live` | 10-abr-2027: **6,293.44** (R-53 = 1,258.70 R-54 + 5,034.74 R-55) |
| M-13 | Calendario unificado de pagos (N-192, N-027) | pagos de todas las líneas en el rango con partes y estado; promociones vigentes y cortes de tarjeta | mes visible o hoy..+60 | `period`, `credit_id`, `family` | scheduled_payment, scheduled_payment_share, credit_payment, credit_purchase, promotion, credit_statement | GET `/api/v1/scheduled-payments`, GET `/api/v1/calendar/events` | pagos de la despensa 1: 30-oct, 15-nov, 30-nov y 15-dic-2026 por **973.75** (R-13 a R-16) |
| M-14 | Plan contra real por mes (N-085, N-028) | por mes (máx. 12): plan total y real total (o de una categoría) | 12 meses | `category_id` | budget, category, expense, expense_share | GET `/api/v1/budget/monthly-trend` | marcador (plan lo captura Axel) |
| M-15 | Costo semanal y kcal (N-129, N-233, N-028) | por semana (12): costo de casillas con receta y kcal promedio por día por persona contra su meta | 12 semanas | `person_id` | meal_plan_entry, recipe, recipe_ingredient, nutrition_fact, price, health_profile | GET `/api/v1/meal-plan/trends` | marcador |
| M-16 | Metas (N-111, N-027) | por meta: ahorrado, % de avance entero, aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan (motor S2) y estado | hoy | `status` | savings_goal, savings_movement | GET `/api/v1/savings-goals` | marcador |
| M-17 | Deuda viva proyectada (N-194, N-209) | por mes (12): deuda viva al cierre por línea; intereses pagados y por pagar | 12 meses | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/credits/debt-trend` | punto de control: deuda viva al 10-abr-2027 **6,293.44** (R-53); los meses sin cifra canon llevan la marca de `criterios-proyeccion-12m` |

## 10. Matriz pantalla P-xx → necesidad N-xxx → endpoint(s)

Una fila por necesidad y pantalla (las necesidades reutilizadas aparecen en cada pantalla que las usa, como en U1 §3.29). Estado: **OK** = cubierta; **FIRME (Kxx)** = cubierta y cerrada en la puerta 1b por la decisión citada (antes PROVISIONAL); **PROVISIONAL** = cubierta en parte, con bloqueo; **BLOQUEO** = sin endpoint por un bloqueo abierto (hoy ninguna); **Local** = vive en el dispositivo (U1 la declara sin API).

| Pantalla | Necesidad | Tipo | Qué necesita (U1) | Operación(es) C1 | Estado |
|---|---|---|---|---|---|
| Global | N-001 | L | Sesión actual | GET `/api/v1/auth/me` | OK |
| Global | N-002 | L | Búsqueda global | GET `/api/v1/search` | OK |
| Global | N-003 | L | Recientes | Recientes en el dispositivo (U1: Local) | Local |
| Global | N-004 | L | Estado del servidor | GET `/health` | OK |
| Global | N-005 | L | Líneas para la barra lateral | GET `/api/v1/credits` | OK |
| Global | N-317 | L | Cola sin conexión | Cola sin conexión en el teléfono (U1: Local; envía con N-054 y N-240, K9) | Local |
| Global | N-006 | E | Cerrar sesión | POST `/api/v1/auth/logout` | OK |
| P-00 | N-011 | L | Estado de bloqueo | POST `/api/v1/auth/login` | OK |
| P-00 | N-010 | E | Iniciar sesión | POST `/api/v1/auth/login` | OK |
| P-01 | N-020 | L | Tarjeta «Mes vs plan» | GET `/api/v1/dashboard/month-vs-plan` | OK |
| P-01 | N-021 | L | Tarjeta «Flujo quincenal» | GET `/api/v1/dashboard/fortnight-flow` | OK |
| P-01 | N-022 | L | Tarjeta «Próximo pago y disponible» | GET `/api/v1/dashboard/next-payment` | OK |
| P-01 | N-023 | L | Tarjeta «Comidas de la semana» | GET `/api/v1/dashboard/week-meals` | FIRME (K43) |
| P-01 | N-024 | L | Tarjeta «Mejor canasta» | GET `/api/v1/dashboard/best-basket` | OK |
| P-01 | N-025 | L | Tarjeta «Ahorro» | GET `/api/v1/dashboard/savings` | OK |
| P-01 | N-026 | L | Participación de Karen | GET `/api/v1/persons/{person_id}/participation` | OK |
| P-01 | N-027 | L | Segundo nivel de las tarjetas | GET `/api/v1/budget/categories` · GET `/api/v1/cash-flow/fortnights` · GET `/api/v1/savings-goals` · GET `/api/v1/meal-plan/week` · GET `/api/v1/basket-comparison` | OK |
| P-01 | N-028 | L | Gráficas de Inicio | GET `/api/v1/budget/monthly-trend` · GET `/api/v1/cash-flow/accrual-vs-cash` · GET `/api/v1/meal-plan/trends` · GET `/api/v1/credits/debt-trend` | OK |
| P-01 | N-029 | L | Calendario del hogar | GET `/api/v1/calendar/events` | OK |
| P-01 | N-081 | L | Presupuesto por categoría | GET `/api/v1/budget/categories` | OK |
| P-01 | N-085 | L | Plan contra real por mes | GET `/api/v1/budget/monthly-trend` | OK |
| P-01 | N-090 | L | Flujo por quincena | GET `/api/v1/cash-flow/fortnights` | OK |
| P-01 | N-091 | L | Devengado contra efectivo, 12 meses | GET `/api/v1/cash-flow/accrual-vs-cash` | OK |
| P-01 | N-111 | L | Metas | GET `/api/v1/views/{entity}` · GET `/api/v1/savings-goals` | OK |
| P-01 | N-121 | L | Totales del plan | GET `/api/v1/meal-plan/week` | OK |
| P-01 | N-129 | L | Tendencias de comidas | GET `/api/v1/meal-plan/trends` | OK |
| P-01 | N-162 | L | Detalle por artículo | GET `/api/v1/basket-comparison` | OK |
| P-01 | N-192 | L | Calendario unificado | GET `/api/v1/calendar/events` · GET `/api/v1/scheduled-payments` | OK |
| P-01 | N-194 | L | Tendencias de deuda | GET `/api/v1/credits/debt-trend` | OK |
| P-01 | N-054 | E | Crear gasto | POST `/api/v1/expenses` | OK |
| P-02 | N-040 | L | Lista de gastos | GET `/api/v1/expenses` · GET `/api/v1/persons` | OK |
| P-02 | N-041 | L | Totales del filtro | GET `/api/v1/expenses/totals` | OK |
| P-02 | N-042 | L | Resumen del periodo | GET `/api/v1/expenses/summary` | FIRME (K39) |
| P-02 | N-043 | L | Gasto por categoría | GET `/api/v1/views/{entity}` | OK |
| P-02 | N-044 | L | Tendencia mensual | GET `/api/v1/expenses/monthly-trend` | OK |
| P-02 | N-045 | L | Calendario de gastos | GET `/api/v1/views/{entity}` | OK |
| P-02 | N-046 | L | Galería de tickets | GET `/api/v1/views/{entity}` · GET `/api/v1/attachments/{attachment_id}/file` | OK |
| P-02 | N-047 | E | Editar en el lugar | PATCH `/api/v1/expenses/{expense_id}` | OK |
| P-02 | N-062 | E | Actualizar gasto | PATCH `/api/v1/expenses/{expense_id}` | OK |
| P-02 | N-063 | E | Eliminar gasto | GET `/api/v1/expenses` · DELETE `/api/v1/expenses/{expense_id}` | OK |
| P-03 | N-050 | L | Categorías sugeridas | GET `/api/v1/categories` | OK |
| P-03 | N-051 | L | Todas las categorías | GET `/api/v1/categories` · GET `/api/v1/cost-share-rules` | OK |
| P-03 | N-052 | L | Valores por defecto | GET `/api/v1/auth/me` · GET `/api/v1/cost-share-rules` | OK |
| P-03 | N-053 | L | Líneas para «forma de pago» | GET `/api/v1/credits` | OK |
| P-03 | N-317 | L | Cola sin conexión | Cola sin conexión en el teléfono (U1: Local; envía con N-054 y N-240, K9) | Local |
| P-03 | N-054 | E | Crear gasto | POST `/api/v1/expenses` | OK |
| P-03 | N-055 | E | Deshacer | DELETE `/api/v1/expenses/{expense_id}` | OK |
| P-03 | N-102 | E | Registrar ingreso extra | POST `/api/v1/incomes` | OK |
| P-03 | N-240 | E | Registrar peso | POST `/api/v1/weight-entries` | OK |
| P-04 | N-060 | L | Gasto | GET `/api/v1/expenses/{expense_id}` · GET `/api/v1/attachments/{attachment_id}/file` | OK |
| P-04 | N-061 | L | Renglones del ticket | GET `/api/v1/expenses/{expense_id}/items` | OK |
| P-04 | N-062 | E | Actualizar gasto | PATCH `/api/v1/expenses/{expense_id}` | OK |
| P-04 | N-063 | E | Eliminar gasto | GET `/api/v1/expenses` · DELETE `/api/v1/expenses/{expense_id}` | OK |
| P-04 | N-064 | E | Adjuntar o quitar ticket | POST `/api/v1/expenses/{expense_id}/attachments` · DELETE `/api/v1/attachments/{attachment_id}` | FIRME (K43) |
| P-04 | N-065 | E | Marcar venidero como pagado | POST `/api/v1/expenses/{expense_id}/mark-paid` | OK |
| P-04 | N-066 | E | Duplicar gasto | POST `/api/v1/expenses/{expense_id}/duplicate` | OK |
| P-05 | N-071 | L | Resultado de la lectura | GET `/api/v1/attachments/{attachment_id}/file` · GET `/api/v1/ticket-scans/{attachment_id}` | OK |
| P-05 | N-072 | L | Buscar en catálogo | GET `/api/v1/search` · GET `/api/v1/products` | OK |
| P-05 | N-073 | L | Prellenado desde la lista | GET `/api/v1/shopping-list` | OK |
| P-05 | N-076 | L | Tickets por confirmar | GET `/api/v1/ticket-scans` | OK |
| P-05 | N-070 | E | Leer ticket | POST `/api/v1/ticket-scans` | OK |
| P-05 | N-074 | E | Confirmar ticket | POST `/api/v1/ticket-scans/{attachment_id}/confirm` | OK |
| P-05 | N-075 | E | Guardar para después | PUT `/api/v1/ticket-scans/{attachment_id}/draft` | OK |
| P-06 | N-080 | L | Resumen del mes | GET `/api/v1/periods` · GET `/api/v1/budget/month` | OK |
| P-06 | N-081 | L | Presupuesto por categoría | GET `/api/v1/budget/categories` | OK |
| P-06 | N-082 | L | Quincenas del mes | GET `/api/v1/periods` · GET `/api/v1/budget/fortnights` | OK |
| P-06 | N-083 | L | Despensa devengada | GET `/api/v1/budget/accrued-groceries` | OK |
| P-06 | N-084 | L | Alertas activas | GET `/api/v1/budget/alerts` | OK |
| P-06 | N-085 | L | Plan contra real por mes | GET `/api/v1/budget/monthly-trend` | OK |
| P-06 | N-086 | L | Calendario del mes | GET `/api/v1/calendar/events` | OK |
| P-06 | N-087 | E | Ajustar presupuesto | POST `/api/v1/budgets` · PATCH `/api/v1/budgets/{budget_id}` | OK |
| P-07 | N-090 | L | Flujo por quincena | GET `/api/v1/cash-flow/fortnights` | OK |
| P-07 | N-091 | L | Devengado contra efectivo, 12 meses | GET `/api/v1/cash-flow/accrual-vs-cash` | OK |
| P-07 | N-092 | L | Deuda viva a una fecha | GET `/api/v1/debt/live` | OK |
| P-07 | N-093 | L | Supuestos de la proyección | GET `/api/v1/cash-flow/assumptions` | OK |
| P-08 | N-100 | L | Ingresos del mes | GET `/api/v1/periods` · GET `/api/v1/incomes` · GET `/api/v1/incomes/month` | OK |
| P-08 | N-103 | L | Ingreso por mes | GET `/api/v1/incomes/monthly-trend` | OK |
| P-08 | N-104 | L | Calendario de ingresos | GET `/api/v1/calendar/events` | OK |
| P-08 | N-101 | E | Capturar quincena | POST `/api/v1/incomes` | OK |
| P-08 | N-102 | E | Registrar ingreso extra | POST `/api/v1/incomes` | OK |
| P-08 | N-105 | E | Corregir o eliminar ingreso | GET `/api/v1/incomes` · PATCH `/api/v1/incomes/{income_id}` · DELETE `/api/v1/incomes/{income_id}` | OK |
| P-09 | N-110 | L | Fondo y totales | GET `/api/v1/savings/summary` | OK |
| P-09 | N-111 | L | Metas | GET `/api/v1/views/{entity}` · GET `/api/v1/savings-goals` | OK |
| P-09 | N-112 | L | Movimientos de ahorro | GET `/api/v1/savings-movements` | OK |
| P-09 | N-114 | L | Tendencia de ahorro | GET `/api/v1/savings/monthly-trend` | OK |
| P-09 | N-115 | L | Calendario de ahorro | GET `/api/v1/calendar/events` | OK |
| P-09 | N-113 | E | Registrar movimiento | POST `/api/v1/savings-movements` | OK |
| P-09 | N-116 | E | Corregir o eliminar movimiento | GET `/api/v1/savings-movements` · PATCH `/api/v1/savings-movements/{savings_movement_id}` · DELETE `/api/v1/savings-movements/{savings_movement_id}` | OK |
| P-10 | N-120 | L | Plan de la semana | GET `/api/v1/meal-slots` · GET `/api/v1/meal-plan/week` | FIRME (K38) |
| P-10 | N-121 | L | Totales del plan | GET `/api/v1/meal-plan/week` | OK |
| P-10 | N-122 | L | Recetas para una casilla | GET `/api/v1/recipes` | OK |
| P-10 | N-128 | L | Calendario mensual | GET `/api/v1/views/{entity}` · GET `/api/v1/meal-plan/month` | OK |
| P-10 | N-129 | L | Tendencias de comidas | GET `/api/v1/meal-plan/trends` | OK |
| P-10 | N-123 | E | Asignar receta | POST `/api/v1/meal-plan-entries` · PATCH `/api/v1/meal-plan-entries/{meal_plan_entry_id}` | OK |
| P-10 | N-124 | E | Marcar «fuera» | POST `/api/v1/meal-plan-entries` | OK |
| P-10 | N-125 | E | Vaciar casilla | DELETE `/api/v1/meal-plan-entries/{meal_plan_entry_id}` | OK |
| P-10 | N-126 | E | Copiar semana anterior | POST `/api/v1/meal-plan/copy-previous-week` | OK |
| P-10 | N-127 | E | Generar lista de compras | POST `/api/v1/shopping-list/generate` | OK |
| P-11 | N-135 | L | Recetas | GET `/api/v1/views/{entity}` · GET `/api/v1/recipes` | OK |
| P-11 | N-136 | E | Crear receta | POST `/api/v1/recipes` | OK |
| P-11 | N-137 | E | Sugerir receta con IA | POST `/api/v1/recipes/ai-suggestions` | OK |
| P-12 | N-140 | L | Receta escalada | GET `/api/v1/recipes/{recipe_id}` | OK |
| P-12 | N-141 | L | Uso de la receta | GET `/api/v1/recipes/{recipe_id}/usage` | OK |
| P-12 | N-123 | E | Asignar receta | POST `/api/v1/meal-plan-entries` · PATCH `/api/v1/meal-plan-entries/{meal_plan_entry_id}` | OK |
| P-12 | N-142 | E | Actualizar receta | PATCH `/api/v1/recipes/{recipe_id}` | OK |
| P-12 | N-143 | E | Duplicar receta | POST `/api/v1/recipes/{recipe_id}/duplicate` | OK |
| P-12 | N-144 | E | Eliminar receta | DELETE `/api/v1/recipes/{recipe_id}` | OK |
| P-13 | N-150 | L | Lista de compras de la semana | GET `/api/v1/shopping-list` · GET `/api/v1/product-groups` | OK |
| P-13 | N-151 | E | Marcar comprado | PATCH `/api/v1/shopping-list/items/{shopping_list_item_id}` | OK |
| P-13 | N-152 | E | Agregar artículo | POST `/api/v1/shopping-list/items` | OK |
| P-13 | N-153 | E | Cambiar tienda o cantidad de un artículo | PATCH `/api/v1/shopping-list/items/{shopping_list_item_id}` | OK |
| P-14 | N-160 | L | Canastas disponibles | GET `/api/v1/baskets` | OK |
| P-14 | N-161 | L | Comparativa por tienda | GET `/api/v1/basket-comparison` | OK |
| P-14 | N-162 | L | Detalle por artículo | GET `/api/v1/basket-comparison` | OK |
| P-14 | N-163 | L | Ganadores | GET `/api/v1/basket-comparison` | OK |
| P-14 | N-164 | L | Membresía Sam's | GET `/api/v1/basket-comparison` | OK |
| P-14 | N-165 | E | Usar esta tienda | POST `/api/v1/shopping-list/assign-store` | OK |
| P-15 | N-170 | L | Catálogo | GET `/api/v1/views/{entity}` · GET `/api/v1/products` | OK |
| P-15 | N-171 | L | Precios de un producto | GET `/api/v1/views/{entity}` · GET `/api/v1/products/{product_id}` · GET `/api/v1/products/{product_id}/prices` | OK |
| P-15 | N-172 | L | Resumen de precios | GET `/api/v1/prices/summary` | OK |
| P-15 | N-173 | L | Nutrición del producto | GET `/api/v1/products/{product_id}` · GET `/api/v1/products/{product_id}/nutrition-facts` | OK |
| P-15 | N-174 | E | Capturar precio | POST `/api/v1/prices` | OK |
| P-15 | N-175 | E | Crear o editar producto | POST `/api/v1/products` · PATCH `/api/v1/products/{product_id}` · POST `/api/v1/products/{product_id}/presentations` · PATCH `/api/v1/presentations/{presentation_id}` · POST `/api/v1/products/{product_id}/nutrition-facts` | OK |
| P-15 | N-176 | E | Corregir o eliminar precio | PATCH `/api/v1/prices/{price_id}` · DELETE `/api/v1/prices/{price_id}` | OK |
| P-16 | N-180 | L | Inventario | GET `/api/v1/views/{entity}` · GET `/api/v1/inventory-items` | OK |
| P-16 | N-181 | L | Movimientos de un artículo | GET `/api/v1/inventory-items/{inventory_item_id}/movements` | OK |
| P-16 | N-182 | L | Resumen | GET `/api/v1/inventory/summary` | OK |
| P-16 | N-183 | E | Ajustar existencia | POST `/api/v1/inventory-items/{inventory_item_id}/adjustments` | OK |
| P-16 | N-184 | E | Registrar sobrante | POST `/api/v1/inventory-items` | OK |
| P-17 | N-190 | L | Líneas | GET `/api/v1/credits` | OK |
| P-17 | N-191 | L | Totales | GET `/api/v1/credits/totals` | OK |
| P-17 | N-192 | L | Calendario unificado | GET `/api/v1/calendar/events` · GET `/api/v1/scheduled-payments` | OK |
| P-17 | N-193 | L | Aviso A/B | GET `/api/v1/credits/{credit_id}` | OK |
| P-17 | N-194 | L | Tendencias de deuda | GET `/api/v1/credits/debt-trend` | OK |
| P-18 | N-200 | L | Línea | GET `/api/v1/credits/{credit_id}` | OK |
| P-18 | N-201 | L | Resumen de la línea | GET `/api/v1/credits/{credit_id}` · GET `/api/v1/credits/{credit_id}/availability-preview` | OK |
| P-18 | N-202 | L | Compras de la línea | GET `/api/v1/views/{entity}` · GET `/api/v1/credit-purchases` | OK |
| P-18 | N-203 | L | Serie de disponible | GET `/api/v1/credits/{credit_id}/availability-series` | OK |
| P-18 | N-204 | L | Tasas por plazo | GET `/api/v1/credits/{credit_id}/term-rates` | OK |
| P-18 | N-205 | L | Estado de la tarjeta | GET `/api/v1/credits/{credit_id}/card-status` | OK |
| P-18 | N-206 | L | Amortización | GET `/api/v1/scheduled-payments` | OK |
| P-18 | N-207 | L | Abonos | GET `/api/v1/scheduled-payments` | OK |
| P-18 | N-208 | L | Calendario de la línea | GET `/api/v1/calendar/events` | OK |
| P-18 | N-209 | L | Tendencias de la línea | GET `/api/v1/credits/debt-trend` · GET `/api/v1/credits/{credit_id}/availability-series` | OK |
| P-18 | N-210 | E | Marcar pago | POST `/api/v1/scheduled-payments/{scheduled_payment_id}/payments` · DELETE `/api/v1/credit-payments/{credit_payment_id}` | OK |
| P-19 | N-215 | L | Compra | GET `/api/v1/credit-purchases` · GET `/api/v1/credit-purchases/{credit_purchase_id}` · GET `/api/v1/scheduled-payments` | OK |
| P-19 | N-216 | L | Vista previa del cálculo | POST `/api/v1/credit-purchases/preview` | OK |
| P-19 | N-217 | L | Promociones aplicables | POST `/api/v1/credit-purchases/preview` · GET `/api/v1/promotions` | OK |
| P-19 | N-210 | E | Marcar pago | POST `/api/v1/scheduled-payments/{scheduled_payment_id}/payments` · DELETE `/api/v1/credit-payments/{credit_payment_id}` | OK |
| P-19 | N-218 | E | Crear compra | POST `/api/v1/credit-purchases` | OK |
| P-19 | N-219 | E | Editar o eliminar compra | PATCH `/api/v1/credit-purchases/{credit_purchase_id}` · DELETE `/api/v1/credit-purchases/{credit_purchase_id}` | OK |
| P-20 | N-225 | L | Valores por defecto | GET `/api/v1/settings/optimizer` | OK |
| P-20 | N-226 | L | Resultado del comparador | POST `/api/v1/scenario-comparisons` | FIRME (K43) |
| P-21 | N-230 | L | Perfiles | GET `/api/v1/health-profiles` | FIRME (K36) |
| P-21 | N-231 | L | Calorías de hoy | GET `/api/v1/meal-plan/week` · GET `/api/v1/health/today-kcal` | OK |
| P-21 | N-232 | L | Historial de peso | GET `/api/v1/weight-entries` | OK |
| P-21 | N-233 | L | kcal promedio por semana | GET `/api/v1/meal-plan/trends` | OK |
| P-21 | N-242 | E | Corregir o eliminar peso | PATCH `/api/v1/weight-entries/{weight_entry_id}` · DELETE `/api/v1/weight-entries/{weight_entry_id}` | OK |
| P-21 | N-313 | E | Editar perfiles, metas y fondo | POST `/api/v1/savings-goals` · PATCH `/api/v1/savings-goals/{savings_goal_id}` · POST `/api/v1/health-profiles` · PATCH `/api/v1/health-profiles/{health_profile_id}` | FIRME (K36) |
| P-21 | N-327 | E | Ajuste por objetivo de una persona | POST `/api/v1/health-profiles` · PATCH `/api/v1/health-profiles/{health_profile_id}` | OK |
| P-22 | N-241 | L | Último peso | GET `/api/v1/weight-entries` | OK |
| P-22 | N-317 | L | Cola sin conexión | Cola sin conexión en el teléfono (U1: Local; envía con N-054 y N-240, K9) | Local |
| P-22 | N-240 | E | Registrar peso | POST `/api/v1/weight-entries` | OK |
| P-23 | N-245 | L | Reportes generados | GET `/api/v1/reports` · GET `/api/v1/notification-deliveries` | OK |
| P-23 | N-246 | L | Próximos envíos | GET `/api/v1/jobs` | OK |
| P-23 | N-247 | L | Avisos enviados | GET `/api/v1/views/{entity}` · GET `/api/v1/notification-deliveries` | OK |
| P-23 | N-249 | L | Descargar archivo | GET `/api/v1/reports/{report_id}/file` | OK |
| P-23 | N-248 | E | Generar reporte mensual | POST `/api/v1/reports/monthly` | OK |
| P-23 | N-250 | E | Reintentar envío | POST `/api/v1/reports/{report_id}/send` | OK |
| P-24 | N-255 | L | Meses disponibles | GET `/api/v1/monthly-report/months` | OK |
| P-24 | N-256 | L | Reporte del mes | GET `/api/v1/monthly-report` | OK |
| P-24 | N-257 | L | PDF del mes | GET `/api/v1/reports` · GET `/api/v1/monthly-report` | OK |
| P-24 | N-249 | L | Descargar archivo | GET `/api/v1/reports/{report_id}/file` | OK |
| P-24 | N-248 | E | Generar reporte mensual | POST `/api/v1/reports/monthly` | OK |
| P-25 | N-260 | L | Estado de cuenta | GET `/api/v1/scheduled-payments` · GET `/api/v1/persons/{person_id}/statement` | OK |
| P-25 | N-262 | L | Estados anteriores | GET `/api/v1/reports` | OK |
| P-25 | N-249 | L | Descargar archivo | GET `/api/v1/reports/{report_id}/file` | OK |
| P-25 | N-261 | E | Generar PDF | POST `/api/v1/reports/external-statements` | OK |
| P-26.01 | N-270 | L | Usuarios | GET `/api/v1/users` | OK |
| P-26.01 | N-271 | L | Sesiones | GET `/api/v1/sessions` | OK |
| P-26.01 | N-274 | L | Este dispositivo | GET `/api/v1/system/ca-certificate` | OK |
| P-26.01 | N-272 | E | Cerrar sesiones | DELETE `/api/v1/sessions/{session_id}` · POST `/api/v1/sessions/revoke` | OK |
| P-26.01 | N-273 | E | Cambiar contraseña | POST `/api/v1/auth/password` | OK |
| P-26.01 | N-320 | E | Parámetros de acceso | GET `/api/v1/settings/security` · PATCH `/api/v1/settings/security` | OK |
| P-26.02 | N-275 | L | Presupuesto mensual | GET `/api/v1/budgets` | OK |
| P-26.02 | N-276 | L | Gastos fijos recurrentes | GET `/api/v1/planned-items` | OK |
| P-26.02 | N-277 | L | Despensa programada | GET `/api/v1/planned-items` | OK |
| P-26.02 | N-318 | L | Proyección de los meses 7 a 12 | GET `/api/v1/settings/projection` | OK |
| P-26.02 | N-278 | E | Editar presupuesto, fijos y despensa | POST `/api/v1/budgets` · PATCH `/api/v1/budgets/{budget_id}` · DELETE `/api/v1/budgets/{budget_id}` · POST `/api/v1/planned-items` · PATCH `/api/v1/planned-items/{planned_item_id}` | OK |
| P-26.02 | N-319 | E | Cambiar el interruptor de proyección | PATCH `/api/v1/settings/projection` | OK |
| P-26.03 | N-279 | L | Ingreso esperado | GET `/api/v1/planned-items` | OK |
| P-26.03 | N-280 | E | Editar ingreso esperado | POST `/api/v1/planned-items` · PATCH `/api/v1/planned-items/{planned_item_id}` | OK |
| P-26.04 | N-281 | L | Categorías | GET `/api/v1/views/{entity}` · GET `/api/v1/categories` · GET `/api/v1/cost-share-rules` | OK |
| P-26.04 | N-282 | E | Editar categorías | POST `/api/v1/categories` · PATCH `/api/v1/categories/{category_id}` · PUT `/api/v1/categories/order` | OK |
| P-26.05 | N-283 | L | Créditos | GET `/api/v1/credits` | OK |
| P-26.05 | N-284 | L | Lectura A/B | GET `/api/v1/credits/{credit_id}` | OK |
| P-26.05 | N-285 | L | Vista previa A/B | GET `/api/v1/credits/{credit_id}/availability-preview` | OK |
| P-26.05 | N-286 | L | Tasas por plazo | GET `/api/v1/credits/{credit_id}/term-rates` | OK |
| P-26.05 | N-287 | L | Promociones | GET `/api/v1/views/{entity}` · GET `/api/v1/promotions` | OK |
| P-26.05 | N-288 | L | Valores del comparador | GET `/api/v1/settings/optimizer` | OK |
| P-26.05 | N-289 | E | Editar créditos y condiciones | POST `/api/v1/credits` · PATCH `/api/v1/credits/{credit_id}` · POST `/api/v1/credits/{credit_id}/term-rates` · PATCH `/api/v1/credit-term-rates/{credit_term_rate_id}` · DELETE `/api/v1/credit-term-rates/{credit_term_rate_id}` · POST `/api/v1/promotions` · PATCH `/api/v1/promotions/{promotion_id}` · DELETE `/api/v1/promotions/{promotion_id}` · PATCH `/api/v1/settings/optimizer` | OK |
| P-26.05 | N-321 | E | Tiempo máximo del comparador | GET `/api/v1/settings/optimizer` · PATCH `/api/v1/settings/optimizer` | OK |
| P-26.06 | N-290 | L | Tiendas | GET `/api/v1/views/{entity}` · GET `/api/v1/stores` · GET `/api/v1/product-groups` | OK |
| P-26.06 | N-291 | L | Fuentes de precio | GET `/api/v1/price-sources` · GET `/api/v1/settings/prices` · GET `/api/v1/jobs` | OK |
| P-26.06 | N-292 | E | Editar tiendas y fuentes | POST `/api/v1/stores` · PATCH `/api/v1/stores/{store_id}` · PUT `/api/v1/stores/{store_id}/product-groups` · PATCH `/api/v1/settings/prices` | OK |
| P-26.06 | N-293 | E | Importar corte PROFECO | POST `/api/v1/profeco-imports` | OK |
| P-26.06 | N-322 | E | Parámetros de precios | GET `/api/v1/settings/prices` · PATCH `/api/v1/settings/prices` | OK |
| P-26.07 | N-294 | L | Recetario y nutrición | GET `/api/v1/meal-slots` · GET `/api/v1/settings/recipes` | OK |
| P-26.07 | N-295 | E | Recargar base nutricional | POST `/api/v1/nutrition-facts/reload` | OK |
| P-26.08 | N-296 | L | Reportes programados | GET `/api/v1/jobs` · GET `/api/v1/notifications` · GET `/api/v1/settings/reports` | OK |
| P-26.08 | N-297 | E | Editar reportes programados | PATCH `/api/v1/jobs/{job_id}` · PATCH `/api/v1/notifications/{notification_id}` · PUT `/api/v1/notifications/{notification_id}/recipients` · PATCH `/api/v1/settings/reports` | FIRME (K43, K47) |
| P-26.08 | N-248 | E | Generar reporte mensual | POST `/api/v1/reports/monthly` | OK |
| P-26.09 | N-298 | L | Alertas | GET `/api/v1/notifications` · GET `/api/v1/settings/alerts` | OK |
| P-26.09 | N-299 | L | Registro de chat | GET `/api/v1/notification-recipients` · GET `/api/v1/telegram-chat-candidates` · PATCH `/api/v1/notification-recipients/{recipient_id}` | OK |
| P-26.09 | N-300 | E | Editar alertas y probar | PATCH `/api/v1/notifications/{notification_id}` · PUT `/api/v1/notifications/{notification_id}/recipients` · POST `/api/v1/notifications/{notification_id}/test` · PATCH `/api/v1/settings/alerts` | FIRME (K43, K47) |
| P-26.09 | N-316 | E | Confirmar chat detectado | POST `/api/v1/telegram-chat-candidates/{candidate_id}/confirm` | OK |
| P-26.09 | N-323 | E | Umbrales del presupuesto y anticipación del pago próximo | GET `/api/v1/settings/alerts` · PATCH `/api/v1/settings/alerts` | OK |
| P-26.09 | N-324 | E | Descartar chat detectado | DELETE `/api/v1/telegram-chat-candidates/{candidate_id}` | OK |
| P-26.10 | N-301 | L | Tema | Tema por dispositivo (U1: Local) | Local |
| P-26.11 | N-302 | L | Respaldos | GET `/api/v1/jobs` · GET `/api/v1/backups` | OK |
| P-26.11 | N-305 | L | Datos semilla | GET `/api/v1/system/seed-status` | OK |
| P-26.11 | N-306 | L | IA opcional | GET `/api/v1/ai/providers` | OK |
| P-26.11 | N-308 | L | Estado del servicio | GET `/api/v1/system/status` | OK |
| P-26.11 | N-303 | E | Respaldar y restaurar | POST `/api/v1/backups` · POST `/api/v1/backups/{backup_id}/restore` | OK |
| P-26.11 | N-304 | E | Importar gastos fijos e ingresos | POST `/api/v1/imports/planned-items/preview` · POST `/api/v1/imports/planned-items` | OK |
| P-26.11 | N-307 | E | Activar o desactivar IA | PATCH `/api/v1/ai/providers/{provider}` | OK |
| P-26.11 | N-325 | E | Parámetros de datos | GET `/api/v1/settings/data` · PATCH `/api/v1/settings/data` | OK |
| P-26.12 | N-309 | L | Personas y reparto | GET `/api/v1/persons` · GET `/api/v1/cost-share-rules` | OK |
| P-26.12 | N-310 | E | Editar personas y reparto | PATCH `/api/v1/persons/{person_id}` · PATCH `/api/v1/cost-share-rules/{cost_share_rule_id}` | OK |
| P-26.13 | N-311 | L | Perfiles de salud | GET `/api/v1/health-profiles` | FIRME (K36) |
| P-26.13 | N-312 | L | Metas y fondo | GET `/api/v1/savings-goals` | OK |
| P-26.13 | N-313 | E | Editar perfiles, metas y fondo | POST `/api/v1/savings-goals` · PATCH `/api/v1/savings-goals/{savings_goal_id}` · POST `/api/v1/health-profiles` · PATCH `/api/v1/health-profiles/{health_profile_id}` | FIRME (K36) |
| P-26.13 | N-326 | E | Ajuste por objetivo de la casa | GET `/api/v1/settings/health` · PATCH `/api/v1/settings/health` | OK |
| P-26.13 | N-327 | E | Ajuste por objetivo de una persona | POST `/api/v1/health-profiles` · PATCH `/api/v1/health-profiles/{health_profile_id}` | OK |
| P-26.14 | N-314 | L | Configuración de inventario | GET `/api/v1/products` · GET `/api/v1/inventory-items` | OK |
| P-26.14 | N-315 | E | Editar inventario | POST `/api/v1/inventory-items` · PATCH `/api/v1/inventory-items/{inventory_item_id}` | OK |

## 11. Cobertura

### 11.1 Resumen
- Necesidades de U1: **211** (N-001 a N-327). Con operación: **208**. Locales sin API (U1): **3** (N-003, N-301, N-317). Sin endpoint por bloqueo: **0**. Con cobertura parcial por un bloqueo abierto de §13: **0**. FIRMES desde la puerta 1b (antes con marca): **10**; en 0 de ellas solo queda pendiente dónde viven `send_time` y `advance_days` (bloqueo abierto `S1-parametros-avisos`, §13).
- Operaciones del contrato: **212** en **168** rutas. Con al menos una N-xxx: **211**. Excepciones técnicas sin N-xxx: **1** (§11.2). N-xxx citadas por el contrato que no existen en U1: **0**.
- Pantallas de la matriz de U1 §3.29: **41** (Global, P-00 a P-25 y P-26.01 a P-26.14); filas de la matriz: **236**.
- Tablas de S1: **55**; con al menos una operación: **55** (§11.4). `idempotency_key` es infraestructura de S6: la declaran en `x-tablas` las **41** operaciones con `Idempotency-Key` (§1 punto 8).

### 11.2 Operaciones sin necesidad de U1 (excepciones técnicas, una por una)
| Operación | Motivo |
|---|---|
| GET `/openapi.yaml` (`getOpenApiDocument`) | Estático y técnico: sirve este mismo contrato para desarrollo y para la prueba de contrato de S6 (S6-03). Requiere sesión. Excepción «estáticos» de C1-02. |

Nota: `POST /api/v1/auth/login` (N-010, N-011) y `GET /health` (N-004, franja «Sin conexión») sí tienen necesidad en U1, aunque C1-02 los admite como excepciones.

### 11.3 Necesidades sin endpoint
| Necesidad | Qué es | Motivo |
|---|---|---|
| N-003 | Recientes | Recientes en el dispositivo (U1: Local) |
| N-301 | Tema | Tema por dispositivo (U1: Local) |
| N-317 | Cola sin conexión | Cola sin conexión en el teléfono (U1: Local; envía con N-054 y N-240, K9) |

### 11.4 Tablas de S1 → operaciones (55 de 55)
| Tabla S1 | Operaciones | Ejemplos |
|---|---|---|
| `person` | 46 | `login`, `getCurrentSession`, `listUsers`, `queryView`, `getDashboardFortnightFlow` … |
| `user_account` | 7 | `login`, `getCurrentSession`, `changeOwnPassword`, `listUsers`, `listSessions` … |
| `device_session` | 9 | `login`, `logout`, `getCurrentSession`, `changeOwnPassword`, `listUsers` … |
| `login_attempt` | 1 | `login` |
| `app_setting` | 28 | `login`, `getDashboardBestBasket`, `getProjectionSettings`, `updateProjectionSettings`, `suggestRecipeWithAi` … |
| `audit_log` | 96 | `changeOwnPassword`, `revokeSession`, `revokeSessions`, `createExpense`, `getExpense` … |
| `idempotency_key` | 41 | `createExpense`, `duplicateExpense`, `markExpensePaid`, `uploadExpenseAttachment`, `createTicketScan` … |
| `period` | 10 | `getDashboardMonthVsPlan`, `listPeriods`, `getBudgetMonth`, `getBudgetByCategory`, `getBudgetFortnights` … |
| `cost_share_rule` | 21 | `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary`, `getExpenseMonthlyTrend` … |
| `cost_share_rule_part` | 14 | `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary`, `getExpenseMonthlyTrend` … |
| `category` | 33 | `search`, `queryView`, `getDashboardMonthVsPlan`, `listExpenses`, `createExpense` … |
| `budget` | 12 | `getDashboardMonthVsPlan`, `listBudgets`, `createBudget`, `updateBudget`, `deleteBudget` … |
| `store` | 36 | `search`, `listCalendarEvents`, `queryView`, `getDashboardBestBasket`, `listExpenses` … |
| `product_group` | 7 | `getShoppingList`, `generateShoppingList`, `listStores`, `setStoreProductGroups`, `listProductGroups` … |
| `store_product_group` | 4 | `getShoppingList`, `generateShoppingList`, `listStores`, `setStoreProductGroups` |
| `product` | 24 | `search`, `queryView`, `getDashboardBestBasket`, `listExpenseItems`, `getTicketScan` … |
| `presentation` | 27 | `search`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `listExpenseItems`, `getTicketScan` … |
| `credit` | 31 | `search`, `listCalendarEvents`, `queryView`, `getDashboardFortnightFlow`, `getDashboardNextPayment` … |
| `credit_term_rate` | 16 | `queryView`, `confirmTicketScan`, `listCredits`, `getCreditTotals`, `getDebtTrend` … |
| `promotion` | 12 | `listCalendarEvents`, `queryView`, `confirmTicketScan`, `listCreditPurchases`, `createCreditPurchase` … |
| `planned_item` | 18 | `listCalendarEvents`, `getDashboardFortnightFlow`, `getBudgetMonth`, `getBudgetFortnights`, `getCashFlowFortnights` … |
| `income` | 13 | `search`, `listCalendarEvents`, `queryView`, `getDashboardFortnightFlow`, `getBudgetMonth` … |
| `expense` | 41 | `search`, `listCalendarEvents`, `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow` … |
| `expense_item` | 19 | `getDashboardBestBasket`, `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary` … |
| `expense_share` | 28 | `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow`, `getPersonParticipation`, `listExpenses` … |
| `attachment` | 19 | `queryView`, `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary` … |
| `price` | 25 | `queryView`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `deleteAttachment`, `getTicketScan` … |
| `credit_purchase` | 46 | `search`, `listCalendarEvents`, `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow` … |
| `credit_purchase_share` | 16 | `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getPersonParticipation`, `confirmTicketScan`, `getAccruedGroceries` … |
| `scheduled_payment` | 28 | `listCalendarEvents`, `queryView`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getPersonParticipation` … |
| `scheduled_payment_share` | 23 | `listCalendarEvents`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getPersonParticipation`, `confirmTicketScan` … |
| `credit_statement` | 9 | `listCalendarEvents`, `deleteAttachment`, `listCredits`, `getCreditTotals`, `getDebtTrend` … |
| `credit_payment` | 22 | `listCalendarEvents`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getPersonParticipation`, `getCashFlowFortnights` … |
| `savings_goal` | 14 | `search`, `listCalendarEvents`, `queryView`, `getDashboardSavings`, `getBudgetMonth` … |
| `savings_movement` | 11 | `listCalendarEvents`, `queryView`, `getDashboardSavings`, `getSavingsSummary`, `getSavingsMonthlyTrend` … |
| `meal_slot` | 11 | `getDashboardWeekMeals`, `listMealSlots`, `getMealPlanWeek`, `getMealPlanMonth`, `getMealPlanTrends` … |
| `recipe` | 17 | `search`, `queryView`, `getDashboardWeekMeals`, `getMealPlanWeek`, `getMealPlanTrends` … |
| `recipe_meal_slot` | 7 | `queryView`, `listRecipes`, `createRecipe`, `getRecipe`, `updateRecipe` … |
| `recipe_ingredient` | 17 | `search`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `getMealPlanWeek`, `getMealPlanTrends` … |
| `meal_plan_entry` | 22 | `queryView`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `getMealPlanWeek`, `getMealPlanMonth` … |
| `inventory_item` | 14 | `queryView`, `getDashboardBestBasket`, `confirmTicketScan`, `getMealPlanWeek`, `getMealPlanTrends` … |
| `inventory_movement` | 4 | `confirmTicketScan`, `createInventoryItem`, `listInventoryMovements`, `adjustInventoryItem` |
| `shopping_list` | 8 | `getDashboardBestBasket`, `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem` … |
| `shopping_list_item` | 8 | `getDashboardBestBasket`, `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem` … |
| `health_profile` | 8 | `getDashboardWeekMeals`, `getMealPlanWeek`, `getMealPlanTrends`, `listHealthProfiles`, `createHealthProfile` … |
| `weight_entry` | 7 | `queryView`, `listHealthProfiles`, `listWeightEntries`, `createWeightEntry`, `updateWeightEntry` … |
| `nutrition_fact` | 11 | `getDashboardWeekMeals`, `getMealPlanWeek`, `getMealPlanTrends`, `listRecipes`, `getRecipe` … |
| `notification` | 7 | `listReports`, `generateMonthlyReport`, `listNotificationDeliveries`, `listJobs`, `listNotifications` … |
| `notification_recipient` | 8 | `generateMonthlyReport`, `sendReport`, `listNotificationDeliveries`, `setNotificationRecipients`, `sendTestNotification` … |
| `notification_subscription` | 5 | `generateMonthlyReport`, `sendReport`, `listJobs`, `listNotifications`, `setNotificationRecipients` |
| `notification_delivery` | 6 | `queryView`, `listReports`, `generateMonthlyReport`, `sendReport`, `listNotificationDeliveries` … |
| `job` | 11 | `listCalendarEvents`, `queryView`, `createTicketScan`, `getPriceSources`, `importProfeco` … |
| `backup` | 4 | `queryView`, `listBackups`, `createBackup`, `restoreBackup` |
| `generated_report` | 7 | `listReports`, `generateMonthlyReport`, `generateExternalStatement`, `downloadReportFile`, `sendReport` … |
| `telegram_chat_candidate` | 3 | `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate` |

`idempotency_key` no tiene pantalla propia: es infraestructura de S6. La leen y escriben en la misma transacción todas las operaciones que exigen `Idempotency-Key` (§1 punto 8; S1 C-15), por eso figura en su `x-tablas`, y la purga por antigüedad la hace un job de S6 (retención de 24 h, K6). `inventory_movement`, `recipe_meal_slot`, `shopping_list` y `shopping_list_item` son detalle de entidades del plan (inventario, receta y lista de compras) y tienen sus propias operaciones.

### 11.5 Las 14 secciones de Ajustes → necesidades → operaciones
| Sección | Nombre | Necesidades | Operaciones |
|---|---|---|---|
| P-26.01 | Cuenta y roles | N-270, N-271, N-274, N-272, N-273, N-320 | `listUsers`, `listSessions`, `getCaCertificate`, `revokeSession`, `revokeSessions`, `changeOwnPassword`, `getSecuritySettings`, `updateSecuritySettings` |
| P-26.02 | Presupuesto | N-275, N-276, N-277, N-318, N-278, N-319 | `listBudgets`, `listPlannedItems`, `getProjectionSettings`, `createBudget`, `updateBudget`, `deleteBudget`, `createPlannedItem`, `updatePlannedItem`, `updateProjectionSettings` |
| P-26.03 | Ingresos | N-279, N-280 | `listPlannedItems`, `createPlannedItem`, `updatePlannedItem` |
| P-26.04 | Categorías | N-281, N-282 | `queryView`, `listCategories`, `listCostShareRules`, `createCategory`, `updateCategory`, `reorderCategories` |
| P-26.05 | Créditos y condiciones (incluye promociones) | N-283, N-284, N-285, N-286, N-287, N-288, N-289, N-321 | `listCredits`, `getCredit`, `previewAvailability`, `listCreditTermRates`, `queryView`, `listPromotions`, `getOptimizerSettings`, `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`, `updateOptimizerSettings` |
| P-26.06 | Tiendas y precios (incluye mercados y membresía Sam's) | N-290, N-291, N-292, N-293, N-322 | `queryView`, `listStores`, `listProductGroups`, `getPriceSources`, `getPriceSettings`, `listJobs`, `createStore`, `updateStore`, `setStoreProductGroups`, `updatePriceSettings`, `importProfeco` |
| P-26.07 | Recetas | N-294, N-295 | `listMealSlots`, `getRecipeSettings`, `reloadNutritionBase` |
| P-26.08 | Reportes y frecuencias | N-296, N-297, N-248 | `listJobs`, `listNotifications`, `getReportSettings`, `updateJob`, `updateNotification`, `setNotificationRecipients`, `updateReportSettings`, `generateMonthlyReport` |
| P-26.09 | Notificaciones | N-298, N-299, N-300, N-316, N-323, N-324 | `listNotifications`, `getAlertSettings`, `listNotificationRecipients`, `listTelegramChatCandidates`, `updateNotificationRecipient`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `updateAlertSettings`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate` |
| P-26.10 | Apariencia | N-301 | Sin API: el tema vive en el dispositivo (N-301, Local) |
| P-26.11 | Datos y respaldos | N-302, N-305, N-306, N-308, N-303, N-304, N-307, N-325 | `listJobs`, `listBackups`, `getSeedStatus`, `listAiProviders`, `getSystemStatus`, `createBackup`, `restoreBackup`, `previewPlannedItemsImport`, `importPlannedItems`, `updateAiProvider`, `getDataSettings`, `updateDataSettings` |
| P-26.12 | Personas y reparto | N-309, N-310 | `listPersons`, `listCostShareRules`, `updatePerson`, `updateCostShareRule` |
| P-26.13 | Salud y metas de ahorro | N-311, N-312, N-313, N-326, N-327 | `listHealthProfiles`, `listSavingsGoals`, `createSavingsGoal`, `updateSavingsGoal`, `createHealthProfile`, `updateHealthProfile`, `getHealthSettings`, `updateHealthSettings` |
| P-26.14 | Inventario | N-314, N-315 | `listProducts`, `listInventoryItems`, `createInventoryItem`, `updateInventoryItem` |

### 11.6 Capacidades exigidas (C1-11) → operaciones
| # | Capacidad | Operaciones |
|---|---|---|
| 1 | Las 14 secciones de Ajustes | §11.5 |
| 2 | Sesiones por dispositivo y su cierre | listSessions, revokeSession, revokeSessions, logout |
| 3 | Interruptor A/B y su vista previa | getCredit, updateCredit, previewAvailability |
| 4 | Comparador de escenarios: 4 objetivos, restricciones, pesos y recomendación | compareScenarios, getOptimizerSettings, updateOptimizerSettings |
| 5 | Promociones | listPromotions, createPromotion, updatePromotion, deletePromotion |
| 6 | Tasa por plazo | listCreditTermRates, createCreditTermRate, updateCreditTermRate, deleteCreditTermRate |
| 7 | Calendario unificado de pagos | listScheduledPayments, listCalendarEvents |
| 8 | Créditos de los 4 tipos y sus condiciones | listCredits, createCredit, updateCredit, getCardStatus, listScheduledPayments (amortización y abonos) |
| 9 | Compras a crédito con partes por persona | createCreditPurchase, previewCreditPurchase, getCreditPurchase, markScheduledPaymentPaid, deleteCreditPayment |
| 10 | Plan de comidas | getMealPlanWeek, createMealPlanEntries, updateMealPlanEntry, deleteMealPlanEntry, copyPreviousWeek, getMealPlanMonth, getMealPlanTrends |
| 11 | Lista de compras | getShoppingList, generateShoppingList, createShoppingListItem, updateShoppingListItem, assignShoppingListStore (compartida en el servidor, K1) |
| 12 | Recetas e ingredientes | listRecipes, getRecipe, createRecipe, updateRecipe, duplicateRecipe, deleteRecipe, getRecipeUsage, suggestRecipeWithAi |
| 13 | Perfil de salud y peso | listHealthProfiles, createHealthProfile, updateHealthProfile, getTodayKcal, listWeightEntries, createWeightEntry, updateWeightEntry, deleteWeightEntry |
| 14 | Precios con fuente y frescura | listProductPrices, createPrice, updatePrice, deletePrice, getPriceSummary, getPriceSettings, getPriceSources |
| 15 | OCR con confirmación humana | createTicketScan, listTicketScans, getTicketScan, confirmTicketScan |
| 16 | Mercados con día de visita | listStores (kind=market), createStore, updateStore, setStoreProductGroups |
| 17 | Membresía de Sam's (punto de equilibrio) | getBasketComparison (membership), updateStore (membership_annual_fee_cents) |
| 18 | Gastos con reparto y adjuntos | createExpense, updateExpense, listCostShareRules, uploadExpenseAttachment, downloadAttachment, deleteAttachment |
| 19 | Ingresos fijos y extra | createIncome, listIncomes, getIncomeMonth, updateIncome, deleteIncome, listPlannedItems |
| 20 | Categorías | listCategories, createCategory, updateCategory, reorderCategories |
| 21 | Metas y fondo de ahorro | listSavingsGoals, createSavingsGoal, updateSavingsGoal, getSavingsSummary, createSavingsMovement |
| 22 | Tiendas, productos, presentaciones y equivalencias | listStores, listProducts, createProduct, getProduct, createPresentation, updatePresentation |
| 23 | Comparador de canasta | listBaskets, getBasketComparison, assignShoppingListStore |
| 24 | Inventario | listInventoryItems, createInventoryItem, updateInventoryItem, adjustInventoryItem, listInventoryMovements, getInventorySummary |
| 25 | Nutrición por ingrediente | listNutritionFacts, createNutritionFact, reloadNutritionBase |
| 26 | Notificaciones y destinatarios | listNotifications, updateNotification, setNotificationRecipients, listNotificationRecipients, listTelegramChatCandidates, confirmTelegramChatCandidate, listNotificationDeliveries |
| 27 | Disparo manual del resumen diario y de una alerta | sendTestNotification |
| 28 | Reporte mensual PDF (generar, listar, descargar, reenviar) | generateMonthlyReport, listReports, downloadReportFile, sendReport |
| 29 | Estado de cuenta de Karen bajo demanda | getExternalStatement, generateExternalStatement |
| 30 | Vista web mensual | getMonthlyReport, listMonthlyReportMonths |
| 31 | Jobs | listJobs, updateJob |
| 32 | Respaldos (listar, ejecutar, restaurar) | listBackups, createBackup, restoreBackup |
| 33 | IA opcional apagada por defecto | listAiProviders, updateAiProvider |
| 34 | Vistas de base de datos genéricas | queryView |
| 35 | Métricas del dashboard | getDashboardMonthVsPlan, getDashboardFortnightFlow, getDashboardNextPayment, getDashboardWeekMeals, getDashboardBestBasket, getDashboardSavings |

### 11.7 Parámetros operativos configurables (K44) → operaciones de settings

Todo parámetro operativo es una clave de `app_setting` del hogar que nace con el valor decidido (K44; RN-29 de `reglas-negocio.md` 1.4). Esta tabla sale de las anotaciones `@param` de las migraciones y de los esquemas de settings de `openapi.yaml`; no se escribe a mano. `parse_schema.py` compara cada valor inicial con la siembra real y se detiene si no coinciden.

- **Lectura**: el `GET /api/v1/settings/…` de la sección trae `parameters`, un objeto por clave de `app_setting` (`SettingParameterInteger`, `SettingParameterIntegerList` o `SettingParameterTime`) con `value` (vigente), `initial_value` (sembrado; `null` si no lo hay), `is_initial`, `value_type` y `unit`.
- **Cambio**: el `PATCH` de la misma ruta recibe `parameters` con el valor nuevo por clave. Un valor vacío, `null` o de otro tipo responde 422 `validation_failed` y no se sustituye por otro. Se valida el tipo y el rango de la clave, el de la columna «Rango (cambio)»: es el de los disparadores de S1 (`trg_app_setting_param_int_*` de 00003; `trg_app_setting_param_range_*` y `trg_app_setting_param_time_*` de 00004) y lo declara la propiedad de la clave en el `PATCH` (`minimum`, `maximum` o `pattern`; puerta 2, K50). Fuera de rango: 422 `validation_failed`.
- **Restablecer** es un campo, no una operación: `reset_parameters` lista las claves que vuelven a su `initial_value`. Una clave no puede venir a la vez en `parameters` y en `reset_parameters` (422 `validation_failed`). Un mismo `PATCH` cambia unos y restablece otros, y `audit_log` registra el cambio.
- **Quién usa el valor vigente**: `login`, `changeOwnPassword` y la renovación de la sesión (`auth.*`, `session.max_age_days`); la frescura y el peso de los precios (`price.stale_after_days`, `price.zero_weight_after_days`, RN-12); `importProfeco`; `compareScenarios` (503 `unavailable`); `uploadExpenseAttachment` (413 con el máximo en el mensaje); la purga de `idempotency_key`; `suggested_kcal_target` (`health.*` y el ajuste por persona de `health_profile`, RN-17); los avisos de presupuesto y de pago próximo y la hora, la anticipación y el día de los demás avisos (`alerts.*`); los reportes programados (`reports.*`, S7); la ventana de fechas del comparador (`optimizer.date_window_days`) y el OCR (`ocr.timeout_seconds`).

| Clave de `app_setting` | Tipo | Unidad | Valor inicial | Rango (cambio) | Decisión | Sección de Ajustes (propuesta de K44) | Sembrada en | Lectura | Cambio y «restablecer» |
|---|---|---|---|---|---|---|---|---|---|
| `auth.max_failed_attempts` | entero | intentos | `5` | `minimum: 0` | Decisión A, K11 | Cuenta y roles | 00001 | `getSecuritySettings` | `updateSecuritySettings` |
| `auth.lockout_window_minutes` | entero | minutos | `15` | `minimum: 0` | Decisión A, K11 | Cuenta y roles | 00001 | `getSecuritySettings` | `updateSecuritySettings` |
| `auth.password_min_length` | entero | caracteres | `12` | `minimum: 0` | K6 | Cuenta y roles | 00003 | `getSecuritySettings` | `updateSecuritySettings` |
| `session.max_age_days` | entero | días | `90` | `minimum: 0` | K11 | Cuenta y roles | 00002 | `getSecuritySettings` | `updateSecuritySettings` |
| `price.stale_after_days` | entero | días | `30` | `minimum: 1` | Decisión E, K22 | Tiendas y precios | 00001 | `getPriceSettings` | `updatePriceSettings` |
| `price.zero_weight_after_days` | entero | días | `90` | `minimum: 0` | K22, K40 | Tiendas y precios | 00003 | `getPriceSettings` | `updatePriceSettings` |
| `price.profeco_import_max_gb` | entero | GB | `1` | `minimum: 0` | K6 | Tiendas y precios | 00003 | `getPriceSettings` | `updatePriceSettings` |
| `comparator.timeout_seconds` | entero | segundos | `5` | `minimum: 1`, `maximum: 30` | K43 | Créditos y condiciones | 00003 | `getOptimizerSettings` | `updateOptimizerSettings` |
| `attachment.max_size_mb` | entero | MB | `10` | `minimum: 0` | K43 | Datos y respaldos | 00003 | `getDataSettings` | `updateDataSettings` |
| `idempotency.retention_hours` | entero | horas | `24` | `minimum: 0` | K6 | Datos y respaldos | 00002 | `getDataSettings` | `updateDataSettings` |
| `health.goal_adjust_lose_kcal` | entero | kcal | `0` | `minimum: 0` | K36 | Salud y metas de ahorro | 00003 | `getHealthSettings` | `updateHealthSettings` |
| `health.goal_adjust_gain_kcal` | entero | kcal | `0` | `minimum: 0` | K36 | Salud y metas de ahorro | 00003 | `getHealthSettings` | `updateHealthSettings` |
| `alerts.budget_warning_pct` | entero | % del presupuesto | `80` | `minimum: 0` | Decisión F, K44 | Notificaciones | 00003 | `getAlertSettings` | `updateAlertSettings` |
| `alerts.budget_limit_pct` | entero | % del presupuesto | `100` | `minimum: 0` | Decisión F, K44 | Notificaciones | 00003 | `getAlertSettings` | `updateAlertSettings` |
| `alerts.payment_due_days_before` | lista de enteros | días antes del pago | `[1, 0]` | elementos `minimum: 0` | Decisión F, K44 | Notificaciones | 00003 | `getAlertSettings` | `updateAlertSettings` |
| `reports.daily_send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | `"21:30"` | `pattern: ^([01][0-9]\|2[0-3]):[0-5][0-9]$` | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004 | `getReportSettings` | `updateReportSettings` |
| `reports.weekly_weekday` | entero | día de la semana ISO 8601, de 1 a 7 (7 = domingo) | `7` | `minimum: 1`, `maximum: 7` | Decisión F, lugar K47 | Reportes y frecuencias | 00004 | `getReportSettings` | `updateReportSettings` |
| `reports.weekly_send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | `"19:00"` | `pattern: ^([01][0-9]\|2[0-3]):[0-5][0-9]$` | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004 | `getReportSettings` | `updateReportSettings` |
| `reports.monthly_day_of_month` | entero | día del mes, de 1 a 31 (en un mes más corto sale su último día, Q2-66 B) | `1` | `minimum: 1`, `maximum: 31` | Decisión F, Q2-66 B, lugar K47 | Reportes y frecuencias | 00004 | `getReportSettings` | `updateReportSettings` |
| `reports.monthly_send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | `"08:00"` | `pattern: ^([01][0-9]\|2[0-3]):[0-5][0-9]$` | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004 | `getReportSettings` | `updateReportSettings` |
| `alerts.send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | `"09:00"` | `pattern: ^([01][0-9]\|2[0-3]):[0-5][0-9]$` | Q2-73 (K50), K31 punto 2, K43, lugar K47 | Notificaciones | 00004 | `getAlertSettings` | `updateAlertSettings` |
| `alerts.stale_prices_advance_days` | entero | días antes del día de despensa, 0 o más | `1` | `minimum: 0` | Q2-73 (K50), K31 punto 2, K43, lugar K47 | Notificaciones | 00004 | `getAlertSettings` | `updateAlertSettings` |
| `alerts.meal_plan_missing_weekday` | entero | día de la semana ISO 8601, de 1 a 7 (7 = domingo) | `7` | `minimum: 1`, `maximum: 7` | Decisión F, lugar K47 | Notificaciones | 00004 | `getAlertSettings` | `updateAlertSettings` |
| `optimizer.date_window_days` | entero | días antes y después de la fecha (fecha ± días), 0 o más | `7` | `minimum: 0` | Q2-75 (K50) | Créditos y condiciones | 00004 | `getOptimizerSettings` | `updateOptimizerSettings` |
| `ocr.timeout_seconds` | entero | segundos por imagen, 0 o más | `30` | `minimum: 0` | Q2-76 (K50); sección: propuesta S1 | Datos y respaldos | 00004 | `getDataSettings` | `updateDataSettings` |

Parámetros que existen sin valor sembrado o que viven en otra tabla (no se inventa un valor):

| Parámetro | Para qué | Por qué no tiene valor sembrado aquí |
|---|---|---|
| optimizer.min_available_buffer_cents, optimizer.max_live_purchases, optimizer.zero_debt_by_date | restricciones duras del comparador | las captura Axel en Ajustes «Créditos y condiciones» (K16); mientras falten, el comparador avisa «faltan restricciones» y no recomienda |
| reportes, respaldos, cabeceras, PDF y permisos | cinco de los ocho valores operativos de K31 | sin valor; los proponen S6 y WF-4 y se aprueban en la puerta 4 (K31) |
| ventana de la cola sin conexión | tiempo máximo que una captura espera en la cola del teléfono | no existe: K9 no pone caducidad (las claves de expense y weight_entry no se purgan); si se quiere un límite, lo propone la PWA en WF-3 |

Reglas de rango que no caben en una sola clave (`@param_rango` de las migraciones; las valida la API con 422 `validation_failed`):

- `price.zero_weight_after_days`: igual o mayor que price.stale_after_days (Q2-77; lo valida el motor de precios); decisión Q2-77 (K50).

Parámetros que antes no tenían valor o vivían en `notification.params_json` y ya son claves de `app_setting` (`@param_sin_valor_cerrado`):

- optimizer.date_window_days: clave sembrada con 7 (Q2-75, K50) por la 00004.
- notification.params_json send_time: claves reports.daily_send_time, reports.weekly_send_time, reports.monthly_send_time y alerts.send_time (Q2-72, Q2-73, K47) de la 00004.
- notification.params_json advance_days: clave alerts.stale_prices_advance_days = 1 (Q2-73, K47) de la 00004.
- notification.params_json weekday (weekly_summary y meal_plan_missing): claves reports.weekly_weekday y alerts.meal_plan_missing_weekday = 7 (decisión F, K47) de la 00004.
- notification.params_json day_of_month (monthly_report): clave reports.monthly_day_of_month = 1 (decisión F, Q2-66 B, K47) de la 00004.

## 12. Índice de operaciones

| Método y ruta | operationId | Etiqueta | N-xxx | Tablas S1 (`x-tablas`) |
|---|---|---|---|---|
| GET `/api/v1/ai/providers` | `listAiProviders` | Ajustes | N-306 | app_setting |
| PATCH `/api/v1/ai/providers/{provider}` | `updateAiProvider` | Ajustes | N-307 | app_setting, audit_log |
| DELETE `/api/v1/attachments/{attachment_id}` | `deleteAttachment` | Gastos | N-064 | attachment, audit_log, credit_statement, price |
| GET `/api/v1/attachments/{attachment_id}/file` | `downloadAttachment` | Gastos | N-046, N-060, N-071 | attachment |
| POST `/api/v1/auth/login` | `login` | Ajustes | N-010, N-011 | app_setting, device_session, login_attempt, person, user_account |
| POST `/api/v1/auth/logout` | `logout` | Ajustes | N-006 | device_session |
| GET `/api/v1/auth/me` | `getCurrentSession` | Inicio | N-001, N-052 | device_session, person, user_account |
| POST `/api/v1/auth/password` | `changeOwnPassword` | Ajustes | N-273 | audit_log, device_session, user_account |
| GET `/api/v1/backups` | `listBackups` | Ajustes | N-302 | backup, job |
| POST `/api/v1/backups` | `createBackup` | Ajustes | N-303 | backup, idempotency_key |
| POST `/api/v1/backups/{backup_id}/restore` | `restoreBackup` | Ajustes | N-303 | backup |
| GET `/api/v1/basket-comparison` | `getBasketComparison` | Compras y tiendas | N-161, N-162, N-163, N-164, N-027 | app_setting, expense, expense_item, inventory_item, meal_plan_entry, presentation, price, product, recipe_ingredient, shopping_list, shopping_list_item, store |
| GET `/api/v1/baskets` | `listBaskets` | Compras y tiendas | N-160 | expense, expense_item, meal_plan_entry, shopping_list, shopping_list_item, store |
| GET `/api/v1/budget/accrued-groceries` | `getAccruedGroceries` | Presupuesto | N-083 | category, credit_purchase, credit_purchase_share, expense, expense_share |
| GET `/api/v1/budget/alerts` | `getBudgetAlerts` | Presupuesto | N-084 | budget, category, credit_purchase, expense, expense_share, period |
| GET `/api/v1/budget/categories` | `getBudgetByCategory` | Presupuesto | N-081, N-027 | budget, category, credit_purchase, expense, expense_share, period |
| GET `/api/v1/budget/fortnights` | `getBudgetFortnights` | Presupuesto | N-082 | budget, category, credit_purchase, expense, expense_share, period, planned_item |
| GET `/api/v1/budget/month` | `getBudgetMonth` | Presupuesto | N-080 | budget, category, credit_purchase, expense, expense_share, income, period, planned_item, savings_goal |
| GET `/api/v1/budget/monthly-trend` | `getBudgetMonthlyTrend` | Presupuesto | N-085, N-028 | budget, category, credit_purchase, expense, expense_share, period |
| GET `/api/v1/budgets` | `listBudgets` | Presupuesto | N-275 | budget, category |
| POST `/api/v1/budgets` | `createBudget` | Presupuesto | N-087, N-278 | audit_log, budget, category, idempotency_key |
| DELETE `/api/v1/budgets/{budget_id}` | `deleteBudget` | Presupuesto | N-278 | audit_log, budget |
| PATCH `/api/v1/budgets/{budget_id}` | `updateBudget` | Presupuesto | N-087, N-278 | audit_log, budget |
| GET `/api/v1/calendar/events` | `listCalendarEvents` | Inicio | N-029, N-086, N-104, N-115, N-192, N-208 | credit, credit_payment, credit_purchase, credit_statement, expense, income, job, planned_item, promotion, savings_goal, savings_movement, scheduled_payment, scheduled_payment_share, store |
| GET `/api/v1/cash-flow/accrual-vs-cash` | `getAccrualVsCash` | Presupuesto | N-091, N-028 | credit, credit_payment, credit_purchase, credit_purchase_share, expense, expense_share, person, planned_item, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/cash-flow/assumptions` | `getCashFlowAssumptions` | Presupuesto | N-093 | expense, planned_item |
| GET `/api/v1/cash-flow/fortnights` | `getCashFlowFortnights` | Presupuesto | N-090, N-027 | credit, credit_payment, credit_purchase, credit_purchase_share, expense, expense_share, income, person, planned_item, savings_goal, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/categories` | `listCategories` | Ajustes | N-050, N-051, N-281 | category, cost_share_rule, expense |
| POST `/api/v1/categories` | `createCategory` | Ajustes | N-282 | audit_log, category, cost_share_rule, idempotency_key |
| PUT `/api/v1/categories/order` | `reorderCategories` | Ajustes | N-282 | audit_log, category |
| PATCH `/api/v1/categories/{category_id}` | `updateCategory` | Ajustes | N-282 | audit_log, category, cost_share_rule |
| GET `/api/v1/cost-share-rules` | `listCostShareRules` | Ajustes | N-309, N-051, N-052, N-281 | cost_share_rule, cost_share_rule_part, person |
| PATCH `/api/v1/cost-share-rules/{cost_share_rule_id}` | `updateCostShareRule` | Ajustes | N-310 | audit_log, cost_share_rule, cost_share_rule_part, expense, expense_item, expense_share |
| DELETE `/api/v1/credit-payments/{credit_payment_id}` | `deleteCreditPayment` | Créditos | N-210 | audit_log, credit_payment |
| GET `/api/v1/credit-purchases` | `listCreditPurchases` | Créditos | N-202, N-215 | credit_payment, credit_purchase, credit_purchase_share, expense, promotion, scheduled_payment |
| POST `/api/v1/credit-purchases` | `createCreditPurchase` | Créditos | N-218 | audit_log, credit, credit_purchase, credit_purchase_share, credit_term_rate, expense, expense_share, idempotency_key, person, promotion, scheduled_payment, scheduled_payment_share |
| POST `/api/v1/credit-purchases/preview` | `previewCreditPurchase` | Créditos | N-216, N-217 | credit, credit_term_rate, person, promotion, scheduled_payment |
| DELETE `/api/v1/credit-purchases/{credit_purchase_id}` | `deleteCreditPurchase` | Créditos | N-219 | audit_log, credit_payment, credit_purchase, credit_purchase_share, expense, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/credit-purchases/{credit_purchase_id}` | `getCreditPurchase` | Créditos | N-215 | credit_payment, credit_purchase, credit_purchase_share, expense, promotion, scheduled_payment, scheduled_payment_share |
| PATCH `/api/v1/credit-purchases/{credit_purchase_id}` | `updateCreditPurchase` | Créditos | N-219 | audit_log, credit_purchase |
| DELETE `/api/v1/credit-term-rates/{credit_term_rate_id}` | `deleteCreditTermRate` | Ajustes | N-289 | audit_log, credit_term_rate |
| PATCH `/api/v1/credit-term-rates/{credit_term_rate_id}` | `updateCreditTermRate` | Ajustes | N-289 | audit_log, credit_term_rate |
| GET `/api/v1/credits` | `listCredits` | Créditos | N-005, N-053, N-190, N-283 | credit, credit_payment, credit_purchase, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| POST `/api/v1/credits` | `createCredit` | Ajustes | N-289 | audit_log, credit, idempotency_key, person |
| GET `/api/v1/credits/debt-trend` | `getDebtTrend` | Créditos | N-194, N-209, N-028 | credit, credit_payment, credit_purchase, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/credits/totals` | `getCreditTotals` | Créditos | N-191 | credit, credit_payment, credit_purchase, credit_purchase_share, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/credits/{credit_id}` | `getCredit` | Créditos | N-200, N-201, N-193, N-284 | credit, credit_payment, credit_purchase, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| PATCH `/api/v1/credits/{credit_id}` | `updateCredit` | Ajustes | N-289 | audit_log, credit |
| GET `/api/v1/credits/{credit_id}/availability-preview` | `previewAvailability` | Créditos | N-285, N-201 | credit, credit_payment, credit_purchase, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/credits/{credit_id}/availability-series` | `getAvailabilitySeries` | Créditos | N-203, N-209 | credit, credit_payment, credit_purchase, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/credits/{credit_id}/card-status` | `getCardStatus` | Créditos | N-205 | credit, credit_payment, credit_purchase, credit_statement, credit_term_rate, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/credits/{credit_id}/term-rates` | `listCreditTermRates` | Créditos | N-204, N-286 | credit_term_rate |
| POST `/api/v1/credits/{credit_id}/term-rates` | `createCreditTermRate` | Ajustes | N-289 | audit_log, credit_term_rate, idempotency_key |
| GET `/api/v1/dashboard/best-basket` | `getDashboardBestBasket` | Inicio | N-024 | app_setting, expense, expense_item, inventory_item, meal_plan_entry, presentation, price, product, recipe_ingredient, shopping_list, shopping_list_item, store |
| GET `/api/v1/dashboard/fortnight-flow` | `getDashboardFortnightFlow` | Inicio | N-021 | credit, credit_payment, credit_purchase, credit_purchase_share, expense, expense_share, income, person, planned_item, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/dashboard/month-vs-plan` | `getDashboardMonthVsPlan` | Inicio | N-020 | budget, category, credit_purchase, expense, expense_share, period |
| GET `/api/v1/dashboard/next-payment` | `getDashboardNextPayment` | Inicio | N-022 | credit, credit_payment, credit_purchase, credit_purchase_share, person, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/dashboard/savings` | `getDashboardSavings` | Inicio | N-025 | savings_goal, savings_movement |
| GET `/api/v1/dashboard/week-meals` | `getDashboardWeekMeals` | Inicio | N-023 | health_profile, meal_plan_entry, meal_slot, nutrition_fact, person, presentation, price, recipe, recipe_ingredient |
| GET `/api/v1/debt/live` | `getLiveDebt` | Presupuesto | N-092 | credit, credit_purchase, scheduled_payment |
| GET `/api/v1/expenses` | `listExpenses` | Gastos | N-040, N-063 | attachment, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, store |
| POST `/api/v1/expenses` | `createExpense` | Gastos | N-054 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, idempotency_key, person, store |
| GET `/api/v1/expenses/monthly-trend` | `getExpenseMonthlyTrend` | Gastos | N-044 | attachment, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, store |
| GET `/api/v1/expenses/summary` | `getExpenseSummary` | Gastos | N-042 | attachment, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, store |
| GET `/api/v1/expenses/totals` | `getExpenseTotals` | Gastos | N-041 | attachment, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, store |
| DELETE `/api/v1/expenses/{expense_id}` | `deleteExpense` | Gastos | N-055, N-063 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, store |
| GET `/api/v1/expenses/{expense_id}` | `getExpense` | Gastos | N-060 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, scheduled_payment, store |
| PATCH `/api/v1/expenses/{expense_id}` | `updateExpense` | Gastos | N-047, N-062 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, person, store |
| POST `/api/v1/expenses/{expense_id}/attachments` | `uploadExpenseAttachment` | Gastos | N-064 | attachment, audit_log, expense, idempotency_key |
| POST `/api/v1/expenses/{expense_id}/duplicate` | `duplicateExpense` | Gastos | N-066 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, idempotency_key, person, store |
| GET `/api/v1/expenses/{expense_id}/items` | `listExpenseItems` | Gastos | N-061 | cost_share_rule, cost_share_rule_part, expense_item, presentation, product |
| POST `/api/v1/expenses/{expense_id}/mark-paid` | `markExpensePaid` | Gastos | N-065 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, expense, expense_item, expense_share, idempotency_key, person, store |
| GET `/api/v1/health-profiles` | `listHealthProfiles` | Salud | N-230, N-311 | health_profile, person, weight_entry |
| POST `/api/v1/health-profiles` | `createHealthProfile` | Ajustes | N-313, N-327 | audit_log, health_profile, idempotency_key, person |
| PATCH `/api/v1/health-profiles/{health_profile_id}` | `updateHealthProfile` | Ajustes | N-313, N-327 | audit_log, health_profile |
| GET `/api/v1/health/today-kcal` | `getTodayKcal` | Salud | N-231 | health_profile, meal_plan_entry, nutrition_fact, recipe, recipe_ingredient |
| POST `/api/v1/imports/planned-items` | `importPlannedItems` | Ajustes | N-304 | audit_log, category, cost_share_rule, credit, idempotency_key, person, planned_item |
| POST `/api/v1/imports/planned-items/preview` | `previewPlannedItemsImport` | Ajustes | N-304 | category, cost_share_rule, credit, person, planned_item |
| GET `/api/v1/incomes` | `listIncomes` | Presupuesto | N-100, N-105 | income, planned_item |
| POST `/api/v1/incomes` | `createIncome` | Presupuesto | N-101, N-102 | audit_log, idempotency_key, income, person, planned_item |
| GET `/api/v1/incomes/month` | `getIncomeMonth` | Presupuesto | N-100 | audit_log, income, period, planned_item |
| GET `/api/v1/incomes/monthly-trend` | `getIncomeMonthlyTrend` | Presupuesto | N-103 | income, planned_item |
| DELETE `/api/v1/incomes/{income_id}` | `deleteIncome` | Presupuesto | N-105 | audit_log, income |
| PATCH `/api/v1/incomes/{income_id}` | `updateIncome` | Presupuesto | N-105 | audit_log, income |
| GET `/api/v1/inventory-items` | `listInventoryItems` | Compras y tiendas | N-180, N-314 | inventory_item, meal_plan_entry, product, recipe, recipe_ingredient |
| POST `/api/v1/inventory-items` | `createInventoryItem` | Compras y tiendas | N-184, N-315 | audit_log, idempotency_key, inventory_item, inventory_movement, product, recipe |
| PATCH `/api/v1/inventory-items/{inventory_item_id}` | `updateInventoryItem` | Ajustes | N-315 | audit_log, inventory_item |
| POST `/api/v1/inventory-items/{inventory_item_id}/adjustments` | `adjustInventoryItem` | Compras y tiendas | N-183 | audit_log, idempotency_key, inventory_item, inventory_movement |
| GET `/api/v1/inventory-items/{inventory_item_id}/movements` | `listInventoryMovements` | Compras y tiendas | N-181 | expense_item, inventory_item, inventory_movement, meal_plan_entry |
| GET `/api/v1/inventory/summary` | `getInventorySummary` | Compras y tiendas | N-182 | inventory_item, meal_plan_entry, product, recipe, recipe_ingredient |
| GET `/api/v1/jobs` | `listJobs` | Reportes | N-246, N-302, N-296, N-291 | job, notification, notification_subscription |
| PATCH `/api/v1/jobs/{job_id}` | `updateJob` | Ajustes | N-297 | audit_log, job |
| POST `/api/v1/meal-plan-entries` | `createMealPlanEntries` | Comidas | N-123, N-124 | audit_log, idempotency_key, meal_plan_entry, meal_slot, person, recipe |
| DELETE `/api/v1/meal-plan-entries/{meal_plan_entry_id}` | `deleteMealPlanEntry` | Comidas | N-125 | audit_log, meal_plan_entry |
| PATCH `/api/v1/meal-plan-entries/{meal_plan_entry_id}` | `updateMealPlanEntry` | Comidas | N-123 | audit_log, meal_plan_entry |
| POST `/api/v1/meal-plan/copy-previous-week` | `copyPreviousWeek` | Comidas | N-126 | audit_log, idempotency_key, meal_plan_entry |
| GET `/api/v1/meal-plan/month` | `getMealPlanMonth` | Comidas | N-128 | meal_plan_entry, meal_slot |
| GET `/api/v1/meal-plan/trends` | `getMealPlanTrends` | Comidas | N-129, N-028, N-233 | health_profile, inventory_item, meal_plan_entry, meal_slot, nutrition_fact, person, presentation, price, product, recipe, recipe_ingredient |
| GET `/api/v1/meal-plan/week` | `getMealPlanWeek` | Comidas | N-120, N-121, N-027, N-231 | health_profile, inventory_item, meal_plan_entry, meal_slot, nutrition_fact, person, presentation, price, product, recipe, recipe_ingredient |
| GET `/api/v1/meal-slots` | `listMealSlots` | Comidas | N-120, N-294 | meal_slot |
| GET `/api/v1/monthly-report` | `getMonthlyReport` | Reportes | N-256, N-257 | budget, category, credit, credit_payment, credit_purchase, credit_purchase_share, expense, expense_share, generated_report, income, meal_plan_entry, period, person, planned_item, price, savings_goal, savings_movement, scheduled_payment, scheduled_payment_share, weight_entry |
| GET `/api/v1/monthly-report/months` | `listMonthlyReportMonths` | Reportes | N-255 | expense, generated_report |
| GET `/api/v1/notification-deliveries` | `listNotificationDeliveries` | Reportes | N-247, N-245 | notification, notification_delivery, notification_recipient |
| GET `/api/v1/notification-recipients` | `listNotificationRecipients` | Ajustes | N-299 | notification_recipient, person |
| PATCH `/api/v1/notification-recipients/{recipient_id}` | `updateNotificationRecipient` | Ajustes | N-299 | audit_log, notification_recipient |
| GET `/api/v1/notifications` | `listNotifications` | Ajustes | N-296, N-298 | job, notification, notification_subscription |
| PATCH `/api/v1/notifications/{notification_id}` | `updateNotification` | Ajustes | N-297, N-300 | audit_log, notification |
| PUT `/api/v1/notifications/{notification_id}/recipients` | `setNotificationRecipients` | Ajustes | N-297, N-300 | audit_log, notification_recipient, notification_subscription |
| POST `/api/v1/notifications/{notification_id}/test` | `sendTestNotification` | Ajustes | N-300 | notification, notification_delivery, notification_recipient |
| POST `/api/v1/nutrition-facts/reload` | `reloadNutritionBase` | Ajustes | N-295 | idempotency_key, job, nutrition_fact, product |
| GET `/api/v1/periods` | `listPeriods` | Presupuesto | N-080, N-082, N-100 | period |
| GET `/api/v1/persons` | `listPersons` | Ajustes | N-309, N-040 | person, user_account |
| PATCH `/api/v1/persons/{person_id}` | `updatePerson` | Ajustes | N-310 | audit_log, person |
| GET `/api/v1/persons/{person_id}/participation` | `getPersonParticipation` | Inicio | N-026 | credit, credit_payment, credit_purchase, credit_purchase_share, expense, expense_share, person, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/persons/{person_id}/statement` | `getExternalStatement` | Reportes | N-260 | credit, credit_payment, credit_purchase, credit_purchase_share, expense, expense_share, person, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/planned-items` | `listPlannedItems` | Ajustes | N-276, N-277, N-279 | category, credit, planned_item |
| POST `/api/v1/planned-items` | `createPlannedItem` | Ajustes | N-278, N-280 | audit_log, category, cost_share_rule, credit, idempotency_key, planned_item |
| PATCH `/api/v1/planned-items/{planned_item_id}` | `updatePlannedItem` | Ajustes | N-278, N-280 | audit_log, planned_item |
| PATCH `/api/v1/presentations/{presentation_id}` | `updatePresentation` | Compras y tiendas | N-175 | audit_log, presentation |
| GET `/api/v1/price-sources` | `getPriceSources` | Ajustes | N-291 | app_setting, job, price, store |
| POST `/api/v1/prices` | `createPrice` | Compras y tiendas | N-174 | audit_log, idempotency_key, presentation, price, store |
| GET `/api/v1/prices/summary` | `getPriceSummary` | Compras y tiendas | N-172 | app_setting, presentation, price, product, store |
| DELETE `/api/v1/prices/{price_id}` | `deletePrice` | Compras y tiendas | N-176 | audit_log, price |
| PATCH `/api/v1/prices/{price_id}` | `updatePrice` | Compras y tiendas | N-176 | audit_log, price |
| GET `/api/v1/product-groups` | `listProductGroups` | Compras y tiendas | N-290, N-150 | product_group |
| GET `/api/v1/products` | `listProducts` | Compras y tiendas | N-170, N-072, N-314 | app_setting, presentation, price, product, product_group, store |
| POST `/api/v1/products` | `createProduct` | Compras y tiendas | N-175 | audit_log, idempotency_key, presentation, product |
| GET `/api/v1/products/{product_id}` | `getProduct` | Compras y tiendas | N-171, N-173 | app_setting, nutrition_fact, presentation, price, product, product_group, store |
| PATCH `/api/v1/products/{product_id}` | `updateProduct` | Compras y tiendas | N-175 | audit_log, product |
| GET `/api/v1/products/{product_id}/nutrition-facts` | `listNutritionFacts` | Compras y tiendas | N-173 | nutrition_fact |
| POST `/api/v1/products/{product_id}/nutrition-facts` | `createNutritionFact` | Compras y tiendas | N-175 | audit_log, idempotency_key, nutrition_fact |
| POST `/api/v1/products/{product_id}/presentations` | `createPresentation` | Compras y tiendas | N-175 | audit_log, idempotency_key, presentation |
| GET `/api/v1/products/{product_id}/prices` | `listProductPrices` | Compras y tiendas | N-171 | app_setting, presentation, price, store |
| POST `/api/v1/profeco-imports` | `importProfeco` | Ajustes | N-293 | audit_log, idempotency_key, job, presentation, price, product, store |
| GET `/api/v1/promotions` | `listPromotions` | Créditos | N-217, N-287 | credit, promotion, store |
| POST `/api/v1/promotions` | `createPromotion` | Ajustes | N-289 | audit_log, idempotency_key, promotion |
| DELETE `/api/v1/promotions/{promotion_id}` | `deletePromotion` | Ajustes | N-289 | audit_log, promotion |
| PATCH `/api/v1/promotions/{promotion_id}` | `updatePromotion` | Ajustes | N-289 | audit_log, promotion |
| GET `/api/v1/recipes` | `listRecipes` | Comidas | N-122, N-135 | meal_plan_entry, meal_slot, nutrition_fact, presentation, price, recipe, recipe_ingredient, recipe_meal_slot |
| POST `/api/v1/recipes` | `createRecipe` | Comidas | N-136 | audit_log, idempotency_key, meal_slot, presentation, product, recipe, recipe_ingredient, recipe_meal_slot |
| POST `/api/v1/recipes/ai-suggestions` | `suggestRecipeWithAi` | Comidas | N-137 | app_setting, meal_slot, presentation, product |
| DELETE `/api/v1/recipes/{recipe_id}` | `deleteRecipe` | Comidas | N-144 | audit_log, meal_plan_entry, recipe, recipe_ingredient, recipe_meal_slot |
| GET `/api/v1/recipes/{recipe_id}` | `getRecipe` | Comidas | N-140 | nutrition_fact, presentation, price, product, recipe, recipe_ingredient, recipe_meal_slot |
| PATCH `/api/v1/recipes/{recipe_id}` | `updateRecipe` | Comidas | N-142 | audit_log, meal_slot, recipe, recipe_ingredient, recipe_meal_slot |
| POST `/api/v1/recipes/{recipe_id}/duplicate` | `duplicateRecipe` | Comidas | N-143 | audit_log, idempotency_key, recipe, recipe_ingredient, recipe_meal_slot |
| GET `/api/v1/recipes/{recipe_id}/usage` | `getRecipeUsage` | Comidas | N-141 | meal_plan_entry |
| GET `/api/v1/reports` | `listReports` | Reportes | N-245, N-257, N-262 | generated_report, notification, notification_delivery |
| POST `/api/v1/reports/external-statements` | `generateExternalStatement` | Reportes | N-261 | audit_log, credit, credit_payment, credit_purchase, credit_purchase_share, generated_report, idempotency_key, person, scheduled_payment, scheduled_payment_share |
| POST `/api/v1/reports/monthly` | `generateMonthlyReport` | Reportes | N-248 | audit_log, budget, category, credit_purchase, expense, expense_share, generated_report, idempotency_key, job, notification, notification_delivery, notification_recipient, notification_subscription, period |
| GET `/api/v1/reports/{report_id}/file` | `downloadReportFile` | Reportes | N-249 | generated_report |
| POST `/api/v1/reports/{report_id}/send` | `sendReport` | Reportes | N-250 | generated_report, idempotency_key, notification_delivery, notification_recipient, notification_subscription |
| GET `/api/v1/savings-goals` | `listSavingsGoals` | Presupuesto | N-111, N-312, N-027 | savings_goal, savings_movement |
| POST `/api/v1/savings-goals` | `createSavingsGoal` | Ajustes | N-313 | audit_log, idempotency_key, savings_goal |
| PATCH `/api/v1/savings-goals/{savings_goal_id}` | `updateSavingsGoal` | Ajustes | N-313 | audit_log, savings_goal |
| GET `/api/v1/savings-movements` | `listSavingsMovements` | Presupuesto | N-112, N-116 | savings_goal, savings_movement |
| POST `/api/v1/savings-movements` | `createSavingsMovement` | Presupuesto | N-113 | audit_log, idempotency_key, savings_goal, savings_movement |
| DELETE `/api/v1/savings-movements/{savings_movement_id}` | `deleteSavingsMovement` | Presupuesto | N-116 | audit_log, savings_movement |
| PATCH `/api/v1/savings-movements/{savings_movement_id}` | `updateSavingsMovement` | Presupuesto | N-116 | audit_log, savings_movement |
| GET `/api/v1/savings/monthly-trend` | `getSavingsMonthlyTrend` | Presupuesto | N-114 | savings_goal, savings_movement |
| GET `/api/v1/savings/summary` | `getSavingsSummary` | Presupuesto | N-110 | savings_goal, savings_movement |
| POST `/api/v1/scenario-comparisons` | `compareScenarios` | Créditos | N-226 | app_setting, credit, credit_purchase, credit_term_rate, person, price, promotion, scheduled_payment, scheduled_payment_share, store |
| GET `/api/v1/scheduled-payments` | `listScheduledPayments` | Créditos | N-192, N-206, N-207, N-215, N-260 | credit, credit_payment, credit_purchase, credit_purchase_share, person, scheduled_payment, scheduled_payment_share |
| POST `/api/v1/scheduled-payments/{scheduled_payment_id}/payments` | `markScheduledPaymentPaid` | Créditos | N-210 | audit_log, credit, credit_payment, idempotency_key, person, scheduled_payment, scheduled_payment_share |
| GET `/api/v1/search` | `search` | Inicio | N-002, N-072 | category, credit, credit_purchase, expense, income, presentation, product, recipe, recipe_ingredient, savings_goal, store |
| GET `/api/v1/sessions` | `listSessions` | Ajustes | N-271 | device_session, user_account |
| POST `/api/v1/sessions/revoke` | `revokeSessions` | Ajustes | N-272 | audit_log, device_session |
| DELETE `/api/v1/sessions/{session_id}` | `revokeSession` | Ajustes | N-272 | audit_log, device_session |
| GET `/api/v1/settings/alerts` | `getAlertSettings` | Ajustes | N-298, N-323 | app_setting |
| PATCH `/api/v1/settings/alerts` | `updateAlertSettings` | Ajustes | N-300, N-323 | app_setting, audit_log |
| GET `/api/v1/settings/data` | `getDataSettings` | Ajustes | N-325 | app_setting |
| PATCH `/api/v1/settings/data` | `updateDataSettings` | Ajustes | N-325 | app_setting, audit_log |
| GET `/api/v1/settings/health` | `getHealthSettings` | Ajustes | N-326 | app_setting |
| PATCH `/api/v1/settings/health` | `updateHealthSettings` | Ajustes | N-326 | app_setting, audit_log |
| GET `/api/v1/settings/optimizer` | `getOptimizerSettings` | Ajustes | N-225, N-288, N-321 | app_setting, person |
| PATCH `/api/v1/settings/optimizer` | `updateOptimizerSettings` | Ajustes | N-289, N-321 | app_setting, audit_log |
| GET `/api/v1/settings/prices` | `getPriceSettings` | Ajustes | N-291, N-322 | app_setting |
| PATCH `/api/v1/settings/prices` | `updatePriceSettings` | Ajustes | N-292, N-322 | app_setting, audit_log |
| GET `/api/v1/settings/projection` | `getProjectionSettings` | Ajustes | N-318 | app_setting |
| PATCH `/api/v1/settings/projection` | `updateProjectionSettings` | Ajustes | N-319 | app_setting, audit_log |
| GET `/api/v1/settings/recipes` | `getRecipeSettings` | Ajustes | N-294 | meal_slot, nutrition_fact, recipe |
| GET `/api/v1/settings/reports` | `getReportSettings` | Ajustes | N-296 | app_setting |
| PATCH `/api/v1/settings/reports` | `updateReportSettings` | Ajustes | N-297 | app_setting, audit_log |
| GET `/api/v1/settings/security` | `getSecuritySettings` | Ajustes | N-320 | app_setting |
| PATCH `/api/v1/settings/security` | `updateSecuritySettings` | Ajustes | N-320 | app_setting, audit_log |
| GET `/api/v1/shopping-list` | `getShoppingList` | Compras y tiendas | N-150, N-073 | inventory_item, meal_plan_entry, presentation, price, product, product_group, recipe_ingredient, shopping_list, shopping_list_item, store, store_product_group |
| POST `/api/v1/shopping-list/assign-store` | `assignShoppingListStore` | Compras y tiendas | N-165 | audit_log, idempotency_key, presentation, price, shopping_list, shopping_list_item, store |
| POST `/api/v1/shopping-list/generate` | `generateShoppingList` | Compras y tiendas | N-127 | audit_log, idempotency_key, inventory_item, meal_plan_entry, presentation, price, product, product_group, recipe_ingredient, shopping_list, shopping_list_item, store, store_product_group |
| POST `/api/v1/shopping-list/items` | `createShoppingListItem` | Compras y tiendas | N-152 | audit_log, idempotency_key, presentation, product, shopping_list, shopping_list_item, store |
| PATCH `/api/v1/shopping-list/items/{shopping_list_item_id}` | `updateShoppingListItem` | Compras y tiendas | N-151, N-153 | audit_log, presentation, shopping_list, shopping_list_item, store |
| GET `/api/v1/stores` | `listStores` | Compras y tiendas | N-290 | product_group, store, store_product_group |
| POST `/api/v1/stores` | `createStore` | Ajustes | N-292 | audit_log, idempotency_key, store |
| PATCH `/api/v1/stores/{store_id}` | `updateStore` | Ajustes | N-292 | audit_log, store |
| PUT `/api/v1/stores/{store_id}/product-groups` | `setStoreProductGroups` | Ajustes | N-292 | audit_log, product_group, store_product_group |
| GET `/api/v1/system/ca-certificate` | `getCaCertificate` | Ajustes | N-274 | — (archivo, no base) |
| GET `/api/v1/system/seed-status` | `getSeedStatus` | Ajustes | N-305 | audit_log, credit_purchase, expense, expense_item, planned_item |
| GET `/api/v1/system/status` | `getSystemStatus` | Ajustes | N-308 | — (archivo, no base) |
| GET `/api/v1/telegram-chat-candidates` | `listTelegramChatCandidates` | Ajustes | N-299 | telegram_chat_candidate |
| DELETE `/api/v1/telegram-chat-candidates/{candidate_id}` | `dismissTelegramChatCandidate` | Ajustes | N-324 | audit_log, telegram_chat_candidate |
| POST `/api/v1/telegram-chat-candidates/{candidate_id}/confirm` | `confirmTelegramChatCandidate` | Ajustes | N-316 | audit_log, idempotency_key, notification_recipient, person, telegram_chat_candidate |
| GET `/api/v1/ticket-scans` | `listTicketScans` | Gastos | N-076 | attachment, store |
| POST `/api/v1/ticket-scans` | `createTicketScan` | Gastos | N-070 | attachment, idempotency_key, job |
| GET `/api/v1/ticket-scans/{attachment_id}` | `getTicketScan` | Gastos | N-071 | attachment, category, cost_share_rule, expense_item, presentation, price, product, store |
| POST `/api/v1/ticket-scans/{attachment_id}/confirm` | `confirmTicketScan` | Gastos | N-074 | attachment, audit_log, category, cost_share_rule, cost_share_rule_part, credit_purchase, credit_purchase_share, credit_term_rate, expense, expense_item, expense_share, idempotency_key, inventory_item, inventory_movement, person, presentation, price, promotion, scheduled_payment, scheduled_payment_share, store |
| PUT `/api/v1/ticket-scans/{attachment_id}/draft` | `saveTicketScanDraft` | Gastos | N-075 | attachment, audit_log |
| GET `/api/v1/users` | `listUsers` | Ajustes | N-270 | device_session, person, user_account |
| GET `/api/v1/views/{entity}` | `queryView` | Inicio | N-043, N-045, N-046, N-111, N-128, N-135, N-170, N-171, N-180, N-202, N-247, N-281, N-287, N-290 | attachment, backup, category, credit, credit_purchase, credit_term_rate, device_session, expense, expense_share, income, inventory_item, job, meal_plan_entry, notification_delivery, person, price, product, promotion, recipe, recipe_meal_slot, savings_goal, savings_movement, scheduled_payment, store, user_account, weight_entry |
| GET `/api/v1/weight-entries` | `listWeightEntries` | Salud | N-232, N-241 | health_profile, person, weight_entry |
| POST `/api/v1/weight-entries` | `createWeightEntry` | Salud | N-240 | audit_log, idempotency_key, person, weight_entry |
| DELETE `/api/v1/weight-entries/{weight_entry_id}` | `deleteWeightEntry` | Salud | N-242 | audit_log, weight_entry |
| PATCH `/api/v1/weight-entries/{weight_entry_id}` | `updateWeightEntry` | Salud | N-242 | audit_log, weight_entry |
| GET `/health` | `getHealth` | Inicio | N-004 | — (archivo, no base) |
| GET `/openapi.yaml` | `getOpenApiDocument` | Ajustes | — | — (archivo, no base) |

## 13. PROVISIONAL y bloqueos

| Tema | Qué queda PROVISIONAL en C1 | Bloqueo |
|---|---|---|
| Alcance «también los ya capturados» de un cambio de reparto | qué gastos recalcula `scope=all` de `updateCostShareRule` (`CostShareRuleUpdate.scope`, §7.2) es PROVISIONAL | `C1-reparto-alcance` (abierto, de C1) |
| Hora de los reportes y `job.schedule` | cómo se relaciona `schedule` de los jobs de los reportes (`JobUpdate`, `updateJob`) con las claves `reports.*` de `updateReportSettings` es PROVISIONAL | `C1-hora-reportes-job` (abierto, de C1) |
| «No incluye IVA» sin porcentaje | `Credit.missing_terms` lista `vat` con `vat_included` = 0 y sin `vat_on_interest_bp` (lectura en uso de S3 1.5 §5.3); PROVISIONAL | `S3-iva-lectura` (abierto, pregunta 2) |
| Edición de un gasto anterior después de «solo los nuevos» | `updateExpense` recalcula sus partes con la regla vigente (opción en uso, RN-03 1.4); PROVISIONAL | `S1-reparto-edicion-gasto` (abierto) |

Cerrados en la puerta 1 y ya sin marca: filtro familia (K5), lista de compras (K1), huecos U1 ↔ S1 (K2, K3, K4), parámetros técnicos (K6), campos calculados (K7), quincena y semana (K8), inicio de sesión y sesión (K11), tasas (K12), tope por persona (K13), reposición de la parte externa (K15), restricciones del comparador (K16), puntaje (K17), promedios diario y mensual (K20), meses 7 a 12 (K21), frescura de precios (K22) y escala de actividad (K23). Cerrados en la puerta 1b (C1 v1.2) y ya sin marca: valores de K31 que tocan a C1 (K43: adjunto de 10 MB, comparador de 5 s, claves `send_time` y `advance_days`, «casillas sin armar»), ajuste por objetivo de la meta de kcal (K36), semana por omisión de `getMealPlanWeek` (K38; la lista de compras sigue con `week` obligatoria, R3-01), promedio semanal (K39) y parámetros configurables (K44, §11.7).

Cerrados en la puerta 2 (C1 v1.3) y ya sin marca: lugar de la hora, el día y la anticipación de reportes y avisos (`S1-parametros-avisos`, Q2-11 B, K47), sus valores (Q2-72, Q2-73), la ventana de fechas y el OCR (Q2-75, Q2-76), los rangos (Q2-74, Q2-77, Q2-66 B), el IVA de un crédito (`U4-iva-credito`, Q2-67), el alcance de un cambio de reparto como pregunta (`U4-reparto-retroactivo`, Q2-68 C), la tienda asignada de la lista (`S4-asignacion-mercado`, Q2-39 a Q2-42) y los acumulados de T-2 (`U3-t2-acumulados`, Q2-65).

## 14. Inconsistencias S1 ↔ U1 encontradas (sin tocar sus archivos)

En C1 v1.0 (WF-0) C1 tomó en cada caso la lectura más conservadora y la marcó PROVISIONAL, y pidió las correcciones a S1 en `docs/estado/wf0/S6a-contrato-c1.md`. El esquema 1.1 de S1 dio dónde guardar siete de los diez puntos y la puerta 1 cerró los demás (K1 a K5 y K6); el esquema 1.2 (00002) agregó el borrador del ticket y los chats por confirmar. C1 v1.1 se regeneró contra ese esquema. La columna «Estado» dice qué queda abierto.

| # | Dónde | S1 1.1 dice | U1 dice | Cómo lo cubre C1 | Estado |
|---|---|---|---|---|---|
| 1 | Descripción del gasto | `expense.description` admite NULL (C-16) | N-054: descripción o comercio opcional; F-1 (3 toques) no la pide | Opcional en `ExpenseCreate`; si falta, el gasto queda sin descripción y la pantalla muestra la categoría. La API no inventa un texto de relleno | Resuelto en S1 1.1 |
| 2 | Lista de compras | `shopping_list` y `shopping_list_item` (opción A) | N-150 a N-153 y N-165: marcas de comprado, artículos a mano, tienda y cantidad por artículo; N-127 generar; N-073 terminar compra | `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem` y `assignShoppingListStore` | Resuelto en la puerta 1 (K1): `C1-lista-compras` |
| 3 | Inventario | `inventory_movement` con tipo (`kind`) y motivo del ajuste (`reason`) | N-181 movimientos con tipo y origen; N-183 ajuste con motivo | `listInventoryMovements` lee la tabla; `adjustInventoryItem`, `createInventoryItem` y `confirmTicketScan` escriben su movimiento en la misma transacción | Resuelto en S1 1.1 |
| 4 | Ticket por confirmar | `attachment.ocr_text`, `ocr_status` y, desde 1.2, `ocr_draft_json` y `ocr_draft_saved_at` | N-075 «Guardar para después» | La lectura queda guardada al crearse; `saveTicketScanDraft` guarda las correcciones como borrador y `TicketScan.draft` las devuelve; el borrador no crea gasto, precio ni inventario: solo `confirmTicketScan` guarda precios (decisiones E) y vacía el borrador | Resuelto en la puerta 1 (K3) y el esquema 1.2 |
| 5 | Presupuesto | `budget` con vigencia `valid_from_date`/`valid_to_date`; los renglones de una categoría no se traslapan (invariante 7) | P-26.02 y N-087: monto mensual por categoría | Renglones con vigencia en `Budget`; un traslape responde 409 `budget_overlap` | Resuelto en la puerta 1 (K2): `updateBudget` recibe `scope` (`this_month` u `onward`, por omisión `onward`) y `month`; efecto de cada alcance y ejemplo en §7.1 |
| 6 | Tiempos de una receta | `recipe_meal_slot` (varios por receta) | N-122, N-135, N-136: «tiempos aptos» | `meal_slot_ids` en `Recipe`, `RecipeCreate` y `RecipeUpdate`; filtro `meal_slot_id` de `listRecipes` | Resuelto en S1 1.1 |
| 7 | Ingreso esperado | `planned_item.second_half_amount_cents` (monto del 30 si difiere del 15) | N-279: el 15 y el 30 pueden diferir | Una partida quincenal con sus dos montos en `PlannedItem` | Resuelto en S1 1.1 |
| 8 | Filtro familia | partes del titular y de personas externas | §0.2 «Hogar o Karen»; N-040 «Hogar, Karen, Conjunto» | `family` = household, external o joint, tres valores incluyentes | Resuelto en la puerta 1 (K5) |
| 9 | Registro del chat | `notification_recipient.chat_id` NOT NULL | N-299: instrucciones para registrarlo en el bot (S7) | El bot detecta el chat (`telegram_chat_candidate`, esquema 1.2); `listTelegramChatCandidates` lo muestra enmascarado y `confirmTelegramChatCandidate` lo liga a la persona | Resuelto en la puerta 1 (K4) |
| 10 | Idempotencia | tabla `idempotency_key` (C-15) | — (la pide C1 §1 punto 8) | Encabezado obligatorio en los POST de captura, tabla declarada en su `x-tablas` | Resuelto en S1 1.1; retención de 24 h aprobada (K6) |
