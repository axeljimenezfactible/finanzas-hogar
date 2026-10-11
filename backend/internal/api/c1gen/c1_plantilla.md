# C1 · Contrato de datos y API · Finanzas Hogar

| | |
|---|---|
| Versión | **C1 v1.2** (11-oct-2026, puerta 1b) · `info.version` 1.2.0 de `docs/contratos/openapi.yaml` |
| Autor | Agente S6a-contrato-c1 (flujo WF-0 «Cimientos», fase S6 parte a); corrección r1: agente correccion-C1-r1; corrección r2: agente correccion-C1-r2; v1.1 corrección r1: agente correccion-C1-r1; v1.1 corrección r4: agente correccion-C1-r4; v1.2: agente C1-v1.2; v1.2 corrección r1: agente correccion-C1-r1 |
| Estado | Propuesta para la puerta 1 (C1 + U1 + criterios) |
| Fuentes | `plan.md` §1, §3.4 (S6, S7), §3.5 (U1), §3.7; `decisiones.md`; esquema de S1, versión 1.1 (`backend/db/migrations/00001_esquema.sql`, `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md`); `docs/design/U1-arquitectura.md` (199 necesidades N-xxx) |
| Archivo de máquina | `docs/contratos/openapi.yaml` (OpenAPI 3.1.1): {{N_RUTAS}} rutas, {{N_OPS}} operaciones, {{N_ESQ}} esquemas |
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
8. **Idempotencia**: todo `POST` de captura exige el encabezado `Idempotency-Key` (16 a 128 caracteres `[A-Za-z0-9_-]`, UUID v4 recomendado). Repetir la clave con el mismo cuerpo devuelve la respuesta original con `Idempotent-Replayed: true`; con otro cuerpo, 409 `idempotency_key_reused`. La API guarda cada clave por cuenta en la tabla `idempotency_key` de S1, con la huella SHA-256 del cuerpo y la respuesta (S1 C-15), y la declara en `x-tablas` de cada operación que la exige; un job la purga por antigüedad. Retención de 24 h (K6). Gastos y pesos guardan además su clave en el propio renglón, sin purga, y responden 200 con el mismo recurso al reintento (K9). Lo exigen {{N_IDEM}} operaciones. No la llevan: las consultas por `POST` que no guardan nada (`previewCreditPurchase`, `compareScenarios`, `previewPlannedItemsImport`, `suggestRecipeWithAi`), las acciones de sesión (`login`, `logout`, `changeOwnPassword`, `revokeSessions`), el aviso de prueba (`sendTestNotification`) y la restauración, que lleva su propia confirmación escrita (`restoreBackup`).

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

- **Colección** = lista de registros de una entidad. Toda colección (`{{N_LISTAS}}` operaciones `GET`) usa **los mismos parámetros**: `cursor` (opaco; sin valor = primera página), `limit` (por defecto **50**, como pide U1; máximo **200**, K6) y `sort` (propiedades separadas por coma, `-` = descendente; el `id` cierra siempre el orden para que el cursor sea estable). **Misma envoltura** `Page`: `{"items": [...], "next_cursor": "…" | null}`; **`next_cursor: null` indica que no hay más**. El orden por defecto de cada colección está en su descripción (tomado de U1 §2.5.b).
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

{{FILTROS}}

## 7. Recursos: las 31 entidades de plan §3.4.1

Cada entidad tiene su recurso (lectura y, cuando el usuario la edita, escritura) o declara qué motor u operación la usa. Las vistas siguen U1 §2.5.a y se piden con `GET /api/v1/views/{entity}` (§8).

{{ENTIDADES}}

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

## 8. Vistas de base de datos (tabla, tablero, calendario, galería)

`GET /api/v1/views/{entity}?mode=table|board|calendar|gallery` (decisión G) aplica los filtros globales, `q`, `status`, `sort` y la paginación únicos, y responde `ViewPage` = `Page` + `mode` + `groups`. En **tablero**, `group_by` elige la columna y cada grupo trae `count`, `total_cents` (Σ del importe principal) y `share_bp`; en **calendario**, `date_field` elige la fecha y los grupos son días; en **galería**, cada registro trae `image_attachment_id` cuando hay foto. La vista elegida, filtros y columnas visibles se recuerdan en el dispositivo (U1 §2.5 y su supuesto 2); no hay vistas guardadas compartidas en esta versión.

{{VISTAS}}

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
| M-11 | Devengado contra efectivo (N-091, N-028) | por mes: devengado y efectivo (totales y partes); diferencia acumulada = devengado acumulado − efectivo acumulado (no se rotula deuda viva; `criterios-invariante-deuda`) | mes en curso a +11 | — | expense, expense_share, planned_item, scheduled_payment, scheduled_payment_share, credit_purchase | GET `/api/v1/cash-flow/accrual-vs-cash` | oct-2026: devengado **1,947.50**, efectivo **973.75**, diferencia acumulada 973.75 (R-45); totales P1–P6 14,362.48 / 13,103.78 (R-51, R-52) |
| M-12 | Deuda viva a una fecha (N-092) | Σ de lo que queda por pagar después de la fecha (capital + intereses) de las compras hechas hasta esa fecha; total, por línea y por compra (motor S3) | fecha elegida | `date` | credit_purchase, scheduled_payment, credit | GET `/api/v1/debt/live` | 10-abr-2027: **6,293.44** (R-53 = 1,258.70 R-54 + 5,034.74 R-55) |
| M-13 | Calendario unificado de pagos (N-192, N-027) | pagos de todas las líneas en el rango con partes y estado; promociones vigentes y cortes de tarjeta | mes visible o hoy..+60 | `period`, `credit_id`, `family` | scheduled_payment, scheduled_payment_share, credit_payment, credit_purchase, promotion, credit_statement | GET `/api/v1/scheduled-payments`, GET `/api/v1/calendar/events` | pagos de la despensa 1: 30-oct, 15-nov, 30-nov y 15-dic-2026 por **973.75** (R-13 a R-16) |
| M-14 | Plan contra real por mes (N-085, N-028) | por mes (máx. 12): plan total y real total (o de una categoría) | 12 meses | `category_id` | budget, category, expense, expense_share | GET `/api/v1/budget/monthly-trend` | marcador (plan lo captura Axel) |
| M-15 | Costo semanal y kcal (N-129, N-233, N-028) | por semana (12): costo de casillas con receta y kcal promedio por día por persona contra su meta | 12 semanas | `person_id` | meal_plan_entry, recipe, recipe_ingredient, nutrition_fact, price, health_profile | GET `/api/v1/meal-plan/trends` | marcador |
| M-16 | Metas (N-111, N-027) | por meta: ahorrado, % de avance entero, aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan (motor S2) y estado | hoy | `status` | savings_goal, savings_movement | GET `/api/v1/savings-goals` | marcador |
| M-17 | Deuda viva proyectada (N-194, N-209) | por mes (12): deuda viva al cierre por línea; intereses pagados y por pagar | 12 meses | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/credits/debt-trend` | punto de control: deuda viva al 10-abr-2027 **6,293.44** (R-53); los meses sin cifra canon llevan la marca de `criterios-proyeccion-12m` |

## 10. Matriz pantalla P-xx → necesidad N-xxx → endpoint(s)

Una fila por necesidad y pantalla (las necesidades reutilizadas aparecen en cada pantalla que las usa, como en U1 §3.29). Estado: **OK** = cubierta; **FIRME (Kxx)** = cubierta y cerrada en la puerta 1b por la decisión citada (antes PROVISIONAL); **PROVISIONAL** = cubierta en parte, con bloqueo; **BLOQUEO** = sin endpoint por un bloqueo abierto (hoy ninguna); **Local** = vive en el dispositivo (U1 la declara sin API).

{{MATRIZ}}

## 11. Cobertura

### 11.1 Resumen
{{RESUMEN}}

### 11.2 Operaciones sin necesidad de U1 (excepciones técnicas, una por una)
{{SIN_N}}

### 11.3 Necesidades sin endpoint
{{N_SIN_EP}}

### 11.4 Tablas de S1 → operaciones ({{N_TABLAS_OP}} de {{N_TABLAS}})
{{TABLAS}}

`idempotency_key` no tiene pantalla propia: es infraestructura de S6. La leen y escriben en la misma transacción todas las operaciones que exigen `Idempotency-Key` (§1 punto 8; S1 C-15), por eso figura en su `x-tablas`, y la purga por antigüedad la hace un job de S6 (retención de 24 h, K6). `inventory_movement`, `recipe_meal_slot`, `shopping_list` y `shopping_list_item` son detalle de entidades del plan (inventario, receta y lista de compras) y tienen sus propias operaciones.

### 11.5 Las 14 secciones de Ajustes → necesidades → operaciones
{{AJUSTES}}

### 11.6 Capacidades exigidas (C1-11) → operaciones
{{CAPACIDADES}}

### 11.7 Parámetros operativos configurables (K44) → operaciones de settings

Todo parámetro operativo es una clave de `app_setting` del hogar que nace con el valor decidido (K44; RN-29 de `reglas-negocio.md` 1.3). Esta tabla sale de las anotaciones `@param` de las migraciones y de los esquemas de settings de `openapi.yaml`; no se escribe a mano. `parse_schema.py` compara cada valor inicial con la siembra real y se detiene si no coinciden.

- **Lectura**: el `GET /api/v1/settings/…` de la sección trae `parameters`, un objeto por clave de `app_setting` (`SettingParameterInteger` o `SettingParameterIntegerList`) con `value` (vigente), `initial_value` (sembrado; `null` si no lo hay), `is_initial`, `value_type` y `unit`.
- **Cambio**: el `PATCH` de la misma ruta recibe `parameters` con el valor nuevo por clave. Un valor vacío, `null` o de otro tipo responde 422 `validation_failed` y no se sustituye por otro. Solo se valida tipo y signo (entero mayor o igual a 0, como los disparadores `trg_app_setting_param_int_*` de 00003); no hay topes mínimos ni máximos hasta la puerta 2.
- **Restablecer** es un campo, no una operación: `reset_parameters` lista las claves que vuelven a su `initial_value`. Una clave no puede venir a la vez en `parameters` y en `reset_parameters` (422 `validation_failed`). Un mismo `PATCH` cambia unos y restablece otros, y `audit_log` registra el cambio.
- **Quién usa el valor vigente**: `login`, `changeOwnPassword` y la renovación de la sesión (`auth.*`, `session.max_age_days`); la frescura y el peso de los precios (`price.stale_after_days`, `price.zero_weight_after_days`, RN-12); `importProfeco`; `compareScenarios` (503 `unavailable`); `uploadExpenseAttachment` (413 con el máximo en el mensaje); la purga de `idempotency_key`; `suggested_kcal_target` (`health.*` y el ajuste por persona de `health_profile`, RN-17); los avisos de presupuesto y de pago próximo (`alerts.*`).

{{PARAMS}}

Parámetros que existen sin valor sembrado o que viven en otra tabla (no se inventa un valor):

{{PARAMS_SV}}

{{PARAMS_K43}}

## 12. Índice de operaciones

{{OPERACIONES}}

## 13. PROVISIONAL y bloqueos

| Tema | Qué queda PROVISIONAL en C1 | Bloqueo |
|---|---|---|
| Lugar de los parámetros de los avisos | dónde viven `send_time` y `advance_days` (claves FIRMES de K43, todavía sin valor sembrado) y `weekday` (7) y `day_of_month` (1) de la decisión F: hoy en `notification.params_json` (`Notification` y `NotificationUpdate`); si la respuesta los pasa a `app_setting`, saldrán por una operación de settings con valor inicial y «restablecer» (N-297, N-300) | `S1-parametros-avisos` (abierto, de S1) |

Cerrados en la puerta 1 y ya sin marca: filtro familia (K5), lista de compras (K1), huecos U1 ↔ S1 (K2, K3, K4), parámetros técnicos (K6), campos calculados (K7), quincena y semana (K8), inicio de sesión y sesión (K11), tasas (K12), tope por persona (K13), reposición de la parte externa (K15), restricciones del comparador (K16), puntaje (K17), promedios diario y mensual (K20), meses 7 a 12 (K21), frescura de precios (K22) y escala de actividad (K23). Cerrados en la puerta 1b (C1 v1.2) y ya sin marca: valores de K31 que tocan a C1 (K43: adjunto de 10 MB, comparador de 5 s, claves `send_time` y `advance_days`, «casillas sin armar»), ajuste por objetivo de la meta de kcal (K36), semana por omisión de `getMealPlanWeek` (K38; la lista de compras sigue con `week` obligatoria, R3-01), promedio semanal (K39) y parámetros configurables (K44, §11.7).

## 14. Inconsistencias S1 ↔ U1 encontradas (sin tocar sus archivos)

En C1 v1.0 (WF-0) C1 tomó en cada caso la lectura más conservadora y la marcó PROVISIONAL, y pidió las correcciones a S1 en `docs/estado/wf0/S6a-contrato-c1.md`. El esquema 1.1 de S1 dio dónde guardar siete de los diez puntos y la puerta 1 cerró los demás (K1 a K5 y K6); el esquema 1.2 (00002) agregó el borrador del ticket y los chats por confirmar. C1 v1.1 se regeneró contra ese esquema. La columna «Estado» dice qué queda abierto.

{{INCONSISTENCIAS}}
