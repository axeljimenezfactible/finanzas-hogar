# Diccionario de datos · Finanzas Hogar

**Versión:** 1.5 · migraciones `backend/db/migrations/00001_esquema.sql`, `backend/db/migrations/00002_puerta1.sql`, `backend/db/migrations/00003_puerta1b.sql`, `backend/db/migrations/00004_puerta2.sql` y `backend/db/migrations/00005_tablero.sql` · 11-oct-2026 (respuestas de la puerta 1 y de la puerta 1b, `docs/plan/decisiones.md` §K)
**Autor:** agente S1-modelo-datos (WF-0, fase S1 «Dominio y modelo de datos»); versión 1.2: agente S1-esquema-1.2 (flujo P1); versión 1.3: agente S1-esquema-1.3 (flujo P1c)
**Fuente de verdad:** las migraciones, en orden. Este documento se genera con `python3 -I backend/db/tools/gen_dict.py` a partir de los comentarios estructurados de las migraciones (`-- @table`, `-- @idx`, comentario por columna; desde 1.2 también `-- @alter`, `-- @used_new` y `-- @after`; desde 1.3 `-- @param` y `-- @param_sin_valor`; desde 1.4 `-- @param_sin_valor_cerrado` y `-- @resuelve`) y con la introspección de la base ya migrada (`PRAGMA table_xinfo`, `foreign_key_list`, `index_list`). Si una migración cambia o se agrega otra, el diccionario se regenera con ese comando. No se edita a mano por separado.
**Reglas de negocio con ejemplos numéricos:** `docs/contratos/reglas-negocio.md`.
**Cambios en 1.1** (hallazgos H-01, H-02, H-04 y H-11 de `docs/estado/wf0/REVISION-r1.md`): tablas nuevas `shopping_list` y `shopping_list_item` (provisionales en 1.1; firmes desde 1.2 por K1), `recipe_meal_slot` (sustituye a `recipe.default_meal_slot_id`), `inventory_movement` e `idempotency_key`; `expense.description` admite NULL; `budget` cambia `period_id` por la vigencia `valid_from_date`/`valid_to_date`; columnas nuevas `planned_item.second_half_amount_cents` y `category.kind`; una sola co-titular (`ux_person_single_co_holder`) y disparador `trg_person_external_without_access`; sección «Ajustes: 14 secciones → tabla.columna»; Apariencia queda local por dispositivo, como en U1 y C1.
**Cambios en 1.2** (migración `00002_puerta1.sql`, respuestas de Axel de `docs/plan/decisiones.md` §K): columnas nuevas `attachment.ocr_draft_json` y `.ocr_draft_saved_at` (K3, borrador de «Guardar para después»); tabla nueva `telegram_chat_candidate` (K4, chats que el bot detecta y la persona confirma); `expense.idempotency_key` y `.captured_at`, y `weight_entry.idempotency_key` y `.captured_at`, con sus índices únicos (K9, cola del teléfono); `credit_term_rate.rate_decimals`, `credit_purchase.rate_decimals`, `promotion.rate_decimals` y `credit.annual_rate_decimals` (K12, tasas como se capturaron); `credit.rate_base`, `.min_payment_pct_bp` y `.min_payment_fixed_cents` (K18, tarjeta bancaria y crédito de tienda); `health_profile.goal` (K23); ajustes sembrados `projection.repeat_cadence_enabled` = false (K21), `session.max_age_days` = 90 (K11) e `idempotency.retention_hours` = 24 (K6); las restricciones `optimizer.*` siguen sin valor (K16). Quedan FIRMES la lista de compras (K1), la retención de idempotencia (K6), los rangos de `period` (K8), la vigencia de `budget` (K2) y la escala de actividad (K23). Redacción neutral de `person.fortnight_payment_cap_cents` (K13, hallazgo R2-03). El generador vive en `backend/db/tools/gen_dict.py` y la verificación en `backend/db/tools/s1_verifica.py` (hallazgo R2-08).
**Cambios en 1.3** (migración `00003_puerta1b.sql`, respuestas de Axel de la puerta 1b, `docs/plan/decisiones.md` K36 a K44): columnas nuevas `health_profile.goal_adjust_lose_kcal` y `.goal_adjust_gain_kcal` (K36, ajuste de la meta de calorías por persona; NULL = valor de la casa); claves nuevas de `app_setting` sembradas con su valor decidido (K44): `health.goal_adjust_lose_kcal` = 0 y `health.goal_adjust_gain_kcal` = 0 (K36), `price.zero_weight_after_days` = 90 (K22, K40), `auth.password_min_length` = 12 y `price.profeco_import_max_gb` = 1 (K6), `attachment.max_size_mb` = 10 y `comparator.timeout_seconds` = 5 (K43), y `alerts.budget_warning_pct` = 80, `alerts.budget_limit_pct` = 100 y `alerts.payment_due_days_before` = [1, 0] (decisión F), que salen de `notification.params_json` (`payment_due` deja `days_before` y `budget_threshold` deja `thresholds_pct`; el Down los repone); disparadores `trg_app_setting_param_int_ins` y `trg_app_setting_param_int_upd` (los parámetros enteros solo aceptan un entero JSON mayor o igual a 0) y `trg_app_setting_param_int_list_ins` y `trg_app_setting_param_int_list_upd` (los días antes del pago próximo, una lista JSON de enteros mayores o iguales a 0); sección nueva «Parámetros operativos de app_setting (K44)». Ninguna tabla nueva; el DDL de 00001 y 00002 no cambia. Corrección de la revisión S1-r1 de P1c (`docs/estado/p1c/REVISION-S1-r1.md`, agente correccion-S1-r1): H-01, los umbrales del presupuesto y los días antes del pago próximo son claves de `app_setting` (K44); H-02, cada fila de «Ajustes: 14 secciones» nombra los parámetros de su sección con su valor inicial. Corrección de la revisión S1-r2 de P1c (`docs/estado/p1c/REVISION-S1-r2.md`, agente correccion-S1-r2): M-01, el inventario lista `weekday` (7) y `day_of_month` (1) de la decisión F, que en 1.3 seguían en `notification.params_json` en espera del bloqueo `S1-parametros-avisos` (resuelto en la puerta 2, K47: desde 1.4 son claves de `app_setting`). Ningún valor cambia.
**Cambios en 1.4** (migración `00004_puerta2.sql`, respuestas de Axel de la puerta 2: `docs/estado/puerta2/RESPUESTAS.md` y `docs/plan/decisiones.md` K46, K47 y K50; agente S1-puerta2, flujo P2c): columna nueva `credit.vat_included` (Q2-67 opción 1: 1 = la tasa ya incluye IVA, 0 = se suma el de `vat_on_interest_bp`, NULL = falta capturar; `Credit.missing_terms` de C1 pasa a 5 datos con `vat`; los créditos que ya existían pasan con el significado de 1.3: con porcentaje 0 y sin porcentaje 1); claves nuevas de `app_setting` con su valor aprobado (Q2-11 B y K47 para el lugar; decisión F y K50 para el valor): `reports.daily_send_time` = "21:30", `reports.weekly_weekday` = 7, `reports.weekly_send_time` = "19:00", `reports.monthly_day_of_month` = 1, `reports.monthly_send_time` = "08:00", `alerts.send_time` = "09:00", `alerts.stale_prices_advance_days` = 1, `alerts.meal_plan_missing_weekday` = 7, `optimizer.date_window_days` = 7, `ocr.timeout_seconds` = 30; `weekday` y `day_of_month` salen de `notification.params_json` de `weekly_summary`, `meal_plan_missing` y `monthly_report` (el Down los repone con el valor vigente); disparadores `trg_app_setting_param_range_ins` y `_upd` (rangos aprobados de una sola clave: `comparator.timeout_seconds` de 1 a 30, Q2-74; `price.stale_after_days` de 1 o más, Q2-77; días de la semana de 1 a 7; día del mes de 1 a 31, Q2-66 B; anticipación, ventana de fechas y OCR, 0 o más) y `trg_app_setting_param_time_ins` y `_upd` (horas «HH:MM» de 00:00 a 23:59); la regla entre dos claves de Q2-77 (`price.zero_weight_after_days` mayor o igual que `price.stale_after_days`) la validan el motor de precios y la API, no el esquema (RN-12, invariante 17). Se cierra el bloqueo `S1-parametros-avisos` (K47). Ninguna tabla nueva; el DDL de 00001 a 00003 no cambia; Q2-38, Q2-66 (mes corto), Q2-68 y el resto de la puerta 2 no cambian el esquema (reporte `docs/estado/p2c/S1-puerta2.md`).
**Correcciones de la revisión S1-r1** (11-oct-2026, agente correccion-S1-r1; hallazgos H-02 y H-03 de `docs/estado/p1/REVISION-S1-r1.md`): `shopping_list.start_date` es el lunes que abre la semana del plan (K8) y la base lo exige con los disparadores `trg_shopping_list_monday_ins` y `trg_shopping_list_monday_upd` de la 00002; las columnas del borrador de `attachment` citan a U1 P-05 (N-071, N-075 y N-076) y `telegram_chat_candidate` suma N-316 «Confirmar chat detectado».

Resumen: **55 tablas STRICT**, 152 índices (incluye los automáticos de UNIQUE) y 81 disparadores (27 de reglas de negocio y 54 de `updated_at`).

## Cómo leer este documento
- **nombre_codigo**: nombre en el código (SQL, Go, TypeScript). Los términos de dominio del plan se traducen aquí una sola vez (ver «Glosario de nombre_codigo»).
- **tipo**: tipo STRICT de SQLite (INTEGER o TEXT). No hay REAL ni BLOB, así que ningún monto puede guardarse con decimales.
- **nulo**: «sí» si la columna admite NULL.
- **regla**: PK, llave foránea (FK) con su acción al borrar la madre, unicidad, valor por defecto (def) y la restricción CHECK resumida.
- **Sufijos:** `_cents` dinero en centavos MXN · `_date` fecha de negocio ISO `AAAA-MM-DD` del calendario de America/Mexico_City · `_at` instante UTC ISO 8601 con «Z» · `_bp` puntos base (10000 = 100 %) · `_milli` milésimas de la unidad base · `rate_num`/`rate_den` tasa como fracción exacta · `is_`/`has_` booleano 0/1.
- **PROVISIONAL** marca lo que depende de un bloqueo abierto en `docs/estado/bloqueos/`.

## Mapa entidad del plan → tabla
Las 31 entidades de `plan.md` §3.4.1 están presentes (las que tienen detalle aparecen también como «(detalle)»):

| Entidad del plan (§3.4.1) | nombre_codigo | Título |
|---|---|---|
| persona | `person` | Persona |
| cuenta_usuario | `user_account` | Cuenta de usuario |
| sesion_dispositivo | `device_session` | Sesión por dispositivo |
| periodo | `period` | Periodo |
| regla_reparto | `cost_share_rule` | Regla de reparto |
| regla_reparto (detalle) | `cost_share_rule_part` | Parte de una regla de reparto |
| categoria | `category` | Categoría |
| tienda | `store` | Tienda |
| tienda (detalle «qué se compra ahí») | `store_product_group` | Qué se compra en cada tienda |
| producto | `product` | Producto |
| presentacion | `presentation` | Presentación |
| credito | `credit` | Crédito |
| tasa_por_plazo | `credit_term_rate` | Tasa por plazo |
| promocion | `promotion` | Promoción |
| ingreso (fijo quincenal) y gasto (fijo y venidero): parte recurrente. Coincide con datos/semilla/gastos_fijos_ingresos.csv (plan §8.2) | `planned_item` | Partida planificada |
| ingreso | `income` | Ingreso |
| gasto | `expense` | Gasto |
| gasto (detalle: renglones del ticket) | `expense_item` | Renglón de gasto |
| regla_reparto (resultado congelado por gasto) | `expense_share` | Parte de un gasto |
| gasto (adjunto) | `attachment` | Adjunto |
| precio | `price` | Precio |
| compra_credito | `credit_purchase` | Compra a crédito |
| compra_credito (partes por persona) | `credit_purchase_share` | Parte de una compra a crédito |
| pago_programado | `scheduled_payment` | Pago programado |
| pago_programado (partes por persona) | `scheduled_payment_share` | Parte de un pago programado |
| credito (fechas de corte y pago, pago mínimo) | `credit_statement` | Estado de cuenta de tarjeta |
| pago_programado (pago real) y decisión A: Karen paga directo y repone la línea cuando paga | `credit_payment` | Pago hecho |
| meta_ahorro | `savings_goal` | Meta de ahorro |
| movimiento_ahorro | `savings_movement` | Movimiento de ahorro |
| tiempo_comida | `meal_slot` | Tiempo de comida |
| receta | `recipe` | Receta |
| receta (detalle: tiempos aptos) | `recipe_meal_slot` | Tiempo apto de una receta |
| ingrediente_receta | `recipe_ingredient` | Ingrediente de receta |
| plan_comida | `meal_plan_entry` | Entrada del plan de comidas |
| inventario_item | `inventory_item` | Renglón de inventario |
| inventario_item (detalle: movimientos) | `inventory_movement` | Movimiento de inventario |
| perfil_salud | `health_profile` | Perfil de salud |
| peso_historico | `weight_entry` | Peso histórico |
| nutricion_ingrediente | `nutrition_fact` | Dato nutricional |
| notificacion | `notification` | Aviso |
| destinatario | `notification_recipient` | Destinatario |
| destinatario (destinatarios configurables por aviso) | `notification_subscription` | Suscripción a un aviso |
| notificacion (bitácora de envíos) | `notification_delivery` | Envío de aviso |
| job | `job` | Job programado |
| respaldo | `backup` | Respaldo |

Tablas complementarias: no son entidades de §3.4.1, pero cada una sale de una decisión cerrada o del plan.

| nombre_codigo | Título | Por qué existe |
|---|---|---|
| `login_attempt` | Intento de acceso | decisión A: límite de 5 intentos en 15 minutos |
| `app_setting` | Ajuste | valores de Ajustes que no son de una entidad (pesos del optimizador, frescura, interruptores de IA) |
| `audit_log` | Bitácora de cambios | convención de auditoría de cambios |
| `idempotency_key` | Clave de idempotencia | encabezado Idempotency-Key de los POST de captura (C1 §1 punto 8): una captura repetida no se duplica |
| `budget` | Presupuesto | presupuesto mensual por categoría y total (S2, alertas 80/100 %) |
| `product_group` | Grupo de producto | decisión D: dos verticales de compra (mercados para verduras y carne no especializada, cadenas para el resto) |
| `shopping_list` | Lista de compras | lista de compras de la semana con marcas de comprado y artículos a mano (U1 N-150 a N-153 y N-165), compartida en el servidor entre los dos teléfonos (K1) |
| `shopping_list_item` | Artículo de la lista de compras | artículos de la lista de compras, calculados por S4 o agregados a mano, con tienda, cantidad y marca de comprado, compartidos en el servidor (K1) |
| `telegram_chat_candidate` | Chat de Telegram por confirmar | chats que el bot detecta cuando alguien le escribe y que la persona confirma en Ajustes (K4) |
| `generated_report` | Reporte generado | decisión F: el reporte mensual PDF se envía y se guarda; estado de cuenta de Karen en PDF bajo demanda |

## Glosario de nombre_codigo
| Término de dominio (es-MX) | nombre_codigo |
|---|---|
| titular / co-titular / participante externa | `holder` / `co_holder` / `external` en `person.participation` |
| hogar Axel–Ximena (parte de dinero) | parte a nombre del titular en `expense_share`, `credit_purchase_share` y `scheduled_payment_share` |
| reparto Hogar / Karen / Conjunto 50/50 | `cost_share_rule` con sus `cost_share_rule_part`. El resto y los centavos van a `remainder_person_id` |
| paga directo al acreedor | `person.pays_share_directly` |
| tope de pago por quincena (objetivo, no restricción dura) | `person.fortnight_payment_cap_cents` |
| mes calendario / quincena / quincena 1 y 2 | `period.kind` = `month` / `fortnight`, `period.half` = 1 / 2 |
| día de pago 15 y 30 (o último día del mes) | `period.pay_date` |
| devengado / efectivo | devengo: `expense.expense_date` + `expense.accrual_months` (+ intereses de su `credit_purchase`). Efectivo: fecha de pago de contado o `scheduled_payment.due_date` |
| gasto fijo / variable / venidero | `expense.kind` = `fixed` / `variable` / `upcoming` |
| partida planificada (CSV de Axel, plan §8.2) | `planned_item.kind` = `income` / `fixed_expense` / `upcoming_expense` |
| ingreso fijo quincenal / ingreso extra | `income.kind` = `fixed` / `extra` |
| línea revolvente / tarjeta / préstamo a plazo fijo / crédito de tienda | `credit.kind` = `revolving_line` / `bank_card` / `term_loan` / `store_credit` |
| lectura A/B del disponible y su fecha «hoy» | `credit.availability_reading`, `credit.reading_a_anchor_date`, aviso `credit.availability_confirmed` |
| tasa por plazo | `credit_term_rate` con `rate_num`/`rate_den` |
| compra a crédito · capital · intereses · compra viva | `credit_purchase` · `principal_cents` · `interest_cents` · `status = 'active'` |
| pago programado / pago hecho | `scheduled_payment` / `credit_payment` |
| corte · fecha límite · pago mínimo | `credit_statement.statement_date` · `due_date` · `minimum_payment_cents` |
| promoción (comercio, plazo, compra mínima, vigencia) | `promotion` |
| cadena / sucursal / mercado o tianguis | `store.kind` = `chain` / `branch` / `market` |
| día de visita | `store.visit_weekday` (1 = lunes … 7 = domingo) |
| vertical de compra | `product_group.purchase_channel` = `market` / `chain` |
| producto genérico / presentación / unidad base | `product` / `presentation` / `product.base_unit` = `g` / `ml` / `piece` |
| sustituto | otra `presentation` del mismo `product` |
| fuente del precio | `price.source` = `ticket_manual` / `ticket_ocr` / `site_assisted` / `profeco` / `manual` |
| desayuno / comida / cena / colación | `meal_slot.code` = `breakfast` / `lunch` / `dinner` / `snack` |
| comer fuera | `meal_plan_entry.is_eating_out` |
| existencia por unidad / ritmo de consumo / sobrante | `inventory_item.tracking_mode` = `stock` / `consumption_rate`, sobrante = renglón con `recipe_id` |
| meta kcal · IMC | `health_profile.daily_kcal_target` · se calcula con `health_profile.height_mm` y el último `weight_entry.weight_g` |
| fondo de ahorro / meta con fecha | `savings_goal.kind` = `fund` / `goal` |
| aviso / destinatario / suscripción | `notification` / `notification_recipient` / `notification_subscription` |
| chat de Telegram por confirmar | `telegram_chat_candidate`; al confirmarlo se crea su `notification_recipient` |
| borrador del ticket («Guardar para después») | `attachment.ocr_draft_json` y `attachment.ocr_draft_saved_at` |
| captura sin conexión (cola del teléfono) | `idempotency_key` y `captured_at` en `expense` y `weight_entry` |
| decimales capturados de una tasa | `rate_decimals` en `credit_term_rate`, `credit_purchase` y `promotion`; `credit.annual_rate_decimals` |
| base de la tasa · pago mínimo | `credit.rate_base` (`daily` / `monthly`) · `credit.min_payment_pct_bp` y `credit.min_payment_fixed_cents` |
| objetivo de peso (bajar / mantener / subir) | `health_profile.goal` = `lose` / `maintain` / `gain` |
| respaldo · reporte guardado | `backup` · `generated_report` |

## A. Personas y acceso

### `person` · Persona

- **Entidad del plan (§3.4.1):** persona
- **Descripción:** Persona que participa en el dinero del hogar: el titular (Axel), la co-titular (Ximena) o una participante externa sin cuenta (Karen).
- **Usado por:** Ajustes «Personas y reparto» y «Cuenta y roles». Motores S2 (reparto), S3 (partes por persona y tope por quincena), S4 (plan de comidas), S7 (estado de cuenta de Karen). Filtro «persona» de C1.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `code` | TEXT | no | Clave estable para semilla e importadores (axel, ximena, karen) | único; solo a-z, 0-9 y guion bajo |
| `display_name` | TEXT | no | Nombre que se muestra en pantallas y reportes | no vacío |
| `participation` | TEXT | no | Papel: holder = titular (Axel), co_holder = co-titular con los mismos permisos (Ximena), external = participante externa sin cuenta (Karen) | enumerado, un solo holder |
| `pays_share_directly` | INTEGER | no | 1 = paga su parte directo al acreedor (Karen con Kueski). Su parte es informativa para el flujo de efectivo del titular y repone la línea cuando ella paga | def: 0; booleano |
| `fortnight_payment_cap_cents` | INTEGER | sí | Tope de pago por quincena de la persona, opcional. El comparador lo usa como objetivo (pago parejo y bajo el tope), nunca como restricción dura. No hay un tope fijo: Axel lo captura por persona en Ajustes «Personas y reparto» cuando quiera (K13) | centavos, NULL = sin tope |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- Sin llaves foráneas salientes.
- La referencian: `cost_share_rule.remainder_person_id`, `cost_share_rule_part.person_id`, `credit.holder_person_id`, `credit_payment.payer_person_id`, `credit_purchase_share.person_id`, `expense.paid_by_person_id`, `expense.person_id`, `expense_share.person_id`, `generated_report.person_id`, `health_profile.person_id`, `income.person_id`, `meal_plan_entry.person_id`, `notification_recipient.person_id`, `planned_item.person_id`, `recipe.author_person_id`, `scheduled_payment_share.person_id`, `telegram_chat_candidate.person_id`, `user_account.person_id`, `weight_entry.person_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_person_1` | UNIQUE (code) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |
| `ux_person_single_co_holder` | UNIQUE `person (participation) WHERE participation = 'co_holder'` | Garantiza que haya una sola co-titular: con el titular suman las dos cuentas posibles (decisión A). |
| `ux_person_single_holder` | UNIQUE `person (participation) WHERE participation = 'holder'` | Garantiza que haya un solo titular. |

**Disparadores:**
- `trg_person_external_without_access`: rechaza con «person: una persona con cuenta, salud o plan de comidas no puede pasar a externa (decisión A)».
- `trg_person_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `user_account` · Cuenta de usuario

- **Entidad del plan (§3.4.1):** cuenta_usuario
- **Descripción:** Cuenta para iniciar sesión. Solo dos: Axel y Ximena, con los mismos permisos (no hay roles). Karen nunca tiene cuenta.
- **Usado por:** Ajustes «Cuenta y roles». S6 (autenticación argon2id, límite de intentos).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `person_id` | INTEGER | no | Persona dueña de la cuenta | FK → `person.id` al borrar: RESTRICT; único; una cuenta por persona, nunca una persona externa (disparador) |
| `username` | TEXT | no | Usuario para iniciar sesión | único; único sin distinguir mayúsculas, 3 a 64 caracteres |
| `password_hash` | TEXT | no | Hash argon2id en formato PHC. La contraseña nunca se guarda | debe empezar con $argon2id$ |
| `is_active` | INTEGER | no | 0 = cuenta deshabilitada (no puede iniciar sesión) | def: 1; booleano |
| `password_changed_at` | TEXT | sí | Último cambio de contraseña (UTC) |  |
| `last_login_at` | TEXT | sí | Último inicio de sesión correcto (UTC) |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `app_setting.user_account_id`, `attachment.uploaded_by_user_id`, `device_session.user_account_id`, `idempotency_key.user_account_id`, `login_attempt.user_account_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_user_account_1` | UNIQUE (person_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |
| `sqlite_autoindex_user_account_2` | UNIQUE (username) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_user_account_not_external_ins`: rechaza con «user_account: una persona externa no puede tener cuenta (decisión A)».
- `trg_user_account_not_external_upd`: rechaza con «user_account: una persona externa no puede tener cuenta (decisión A)».
- `trg_user_account_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `device_session` · Sesión por dispositivo

- **Entidad del plan (§3.4.1):** sesion_dispositivo
- **Descripción:** Sesión larga (meses) de una cuenta en un dispositivo. Se puede cerrar desde Ajustes.
- **Usado por:** Ajustes «Cuenta y roles» (cerrar sesiones). S6 (middleware de sesión, purga de vencidas).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `user_account_id` | INTEGER | no | Cuenta dueña de la sesión | FK → `user_account.id` al borrar: CASCADE |
| `token_hash` | TEXT | no | SHA-256 en hexadecimal del token de sesión. El token en claro nunca se guarda | único; 64 caracteres |
| `device_label` | TEXT | sí | Nombre legible del dispositivo (Chrome en Android) |  |
| `user_agent` | TEXT | sí | Agente de usuario reportado por el navegador |  |
| `ip_address` | TEXT | sí | Última IP de la red local desde la que se usó |  |
| `last_seen_at` | TEXT | no | Último uso (UTC) | def: ahora (UTC) |
| `expires_at` | TEXT | no | Vencimiento de la sesión (UTC) | posterior a created_at |
| `revoked_at` | TEXT | sí | Cierre manual desde Ajustes (UTC). NULL = vigente |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- La sesión vence después de crearse.

**Relaciones:**
- `user_account_id` → `user_account.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_device_session_expires` | `device_session (expires_at)` | Purga de sesiones vencidas. |
| `ix_device_session_user` | `device_session (user_account_id, revoked_at)` | Lista las sesiones vigentes de una cuenta (Ajustes, cerrar sesiones). |
| `sqlite_autoindex_device_session_1` | UNIQUE (token_hash) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_device_session_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `login_attempt` · Intento de acceso

- **Entidad del plan (§3.4.1):** (complementaria) decisión A: límite de 5 intentos en 15 minutos
- **Descripción:** Bitácora de intentos de inicio de sesión para bloquear tras 5 fallos en 15 minutos.
- **Usado por:** S6 (autenticación y bloqueo). Job de purga.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `username` | TEXT | no | Usuario tecleado (puede no existir) |  |
| `user_account_id` | INTEGER | sí | Cuenta que corresponde al usuario, si existe | FK → `user_account.id` al borrar: SET NULL |
| `ip_address` | TEXT | sí | IP de origen en la red local |  |
| `succeeded` | INTEGER | no | 1 = acceso correcto | booleano |
| `attempted_at` | TEXT | no | Momento del intento (UTC) | def: ahora (UTC) |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `user_account_id` → `user_account.id` · al borrar la madre: SET NULL (queda sin liga)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_login_attempt_ip` | `login_attempt (ip_address, attempted_at)` | Cuenta los fallos desde una IP en la ventana de 15 minutos. |
| `ix_login_attempt_user` | `login_attempt (user_account_id)` | Llave foránea hacia la cuenta. |
| `ix_login_attempt_username` | `login_attempt (username, attempted_at)` | Cuenta los fallos de un usuario en la ventana de 15 minutos. |

**Disparadores:**
- `trg_login_attempt_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `app_setting` · Ajuste

- **Entidad del plan (§3.4.1):** (complementaria) valores de Ajustes que no son de una entidad (pesos del optimizador, frescura, interruptores de IA)
- **Descripción:** Ajustes del hogar o de una cuenta en formato JSON, con clave por sección.
- **Usado por:** Ajustes (varias secciones). S3 (pesos y restricciones del optimizador), S5 (frescura), S6 (límites de acceso, IA opcional).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `setting_key` | TEXT | no | Clave con prefijo de sección (price.stale_after_days) | minúsculas, dígitos, punto y guion bajo |
| `user_account_id` | INTEGER | sí | NULL = ajuste del hogar. Con valor = preferencia de esa cuenta | FK → `user_account.id` al borrar: CASCADE |
| `value_json` | TEXT | no | Valor en JSON | JSON válido |
| `description` | TEXT | sí | Para qué sirve, en español |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `user_account_id` → `user_account.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_app_setting_user` | `app_setting (user_account_id)` | Llave foránea hacia la cuenta. |
| `ux_app_setting_key_scope` | UNIQUE `app_setting (setting_key, coalesce(user_account_id, 0))` | Una sola fila por clave y alcance (hogar o cuenta). |

**Disparadores:**
- `trg_app_setting_param_int_ins`: rechaza con «app_setting: este parámetro es un entero mayor o igual a 0 (K44)».
- `trg_app_setting_param_int_list_ins`: rechaza con «app_setting: este parámetro es una lista de enteros mayores o iguales a 0 (K44)».
- `trg_app_setting_param_int_list_upd`: rechaza con «app_setting: este parámetro es una lista de enteros mayores o iguales a 0 (K44)».
- `trg_app_setting_param_int_upd`: rechaza con «app_setting: este parámetro es un entero mayor o igual a 0 (K44)».
- `trg_app_setting_param_k54_ins`: rechaza con «app_setting: este parámetro del tablero es un entero dentro de su rango (K54, RN-29)».
- `trg_app_setting_param_k54_upd`: rechaza con «app_setting: este parámetro del tablero es un entero dentro de su rango (K54, RN-29)».
- `trg_app_setting_param_range_ins`: rechaza con «app_setting: este parámetro es un entero dentro de su rango aprobado (puerta 2, RN-29)».
- `trg_app_setting_param_range_upd`: rechaza con «app_setting: este parámetro es un entero dentro de su rango aprobado (puerta 2, RN-29)».
- `trg_app_setting_param_time_ins`: rechaza con «app_setting: esta hora es un texto HH:MM de 00:00 a 23:59, hora del centro de México (puerta 2, RN-29)».
- `trg_app_setting_param_time_upd`: rechaza con «app_setting: esta hora es un texto HH:MM de 00:00 a 23:59, hora del centro de México (puerta 2, RN-29)».
- `trg_app_setting_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `audit_log` · Bitácora de cambios

- **Entidad del plan (§3.4.1):** (complementaria) convención de auditoría de cambios
- **Descripción:** Registro inmutable de cada alta, cambio o baja hecha por la API o por un job: quién, cuándo y el renglón antes y después.
- **Usado por:** S6 (lo escribe en la misma transacción de cada cambio). Ajustes «Datos y respaldos» (historial).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `occurred_at` | TEXT | no | Momento del cambio (UTC) | def: ahora (UTC) |
| `user_account_id` | INTEGER | sí | Cuenta que hizo el cambio. NULL = proceso del sistema. Sin llave foránea para que la bitácora sobreviva a cualquier baja |  |
| `action` | TEXT | no | Tipo de cambio | enumerado |
| `table_name` | TEXT | no | Tabla afectada |  |
| `row_id` | INTEGER | no | id del renglón afectado |  |
| `before_json` | TEXT | sí | Renglón antes del cambio (NULL en insert) | JSON válido |
| `after_json` | TEXT | sí | Renglón después del cambio (NULL en delete) | JSON válido |
| `request_id` | TEXT | sí | Identificador de la petición HTTP o del job que originó el cambio |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- insert lleva solo «después», delete solo «antes», update ambos.

**Relaciones:**
- Sin llaves foráneas salientes.

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_audit_log_row` | `audit_log (table_name, row_id, occurred_at)` | Historial de un renglón. |
| `ix_audit_log_time` | `audit_log (occurred_at)` | Consulta por fecha. |
| `ix_audit_log_user` | `audit_log (user_account_id, occurred_at)` | Cambios hechos por una cuenta. |

**Disparadores:**
- `trg_audit_log_no_delete`: rechaza con «audit_log es inmutable».
- `trg_audit_log_no_update`: rechaza con «audit_log es inmutable».

### `idempotency_key` · Clave de idempotencia

- **Entidad del plan (§3.4.1):** (complementaria) encabezado Idempotency-Key de los POST de captura (C1 §1 punto 8): una captura repetida no se duplica
- **Descripción:** Clave que manda el cliente en cada POST de captura. Si la misma cuenta repite la clave con el mismo cuerpo, la API devuelve la respuesta guardada; con otro cuerpo responde 409 idempotency_key_reused. La retención es de 24 h (K6, app_setting idempotency.retention_hours de la migración 00002); la purga la hace un job. Las capturas de la cola del teléfono guardan además su clave en expense y weight_entry, sin purga (K9).
- **Usado por:** S6 (API: los POST de captura de C1, entre ellos «Registrar gasto» del flujo F-1 de U1) y jobs de S6 (purga por antigüedad).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `user_account_id` | INTEGER | no | Cuenta que mandó la clave | FK → `user_account.id` al borrar: CASCADE |
| `idempotency_key` | TEXT | no | Valor del encabezado Idempotency-Key | 16 a 128 caracteres: letras sin acento, dígitos, guion y guion bajo (C1 §1 punto 8) |
| `http_method` | TEXT | no | Método de la petición | enumerado |
| `request_path` | TEXT | no | Ruta de la petición (/api/v1/expenses) | empieza con / |
| `request_sha256` | TEXT | no | Huella SHA-256 del cuerpo: distingue repetir la petición de reutilizar la clave con otro cuerpo | 64 caracteres |
| `response_status` | INTEGER | sí | Código HTTP de la respuesta guardada. NULL = la petición original sigue en proceso | 100 a 599 |
| `response_body` | TEXT | sí | Cuerpo JSON de la respuesta guardada, para repetirla igual | JSON válido |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Una clave se usa una sola vez por cuenta.

**Relaciones:**
- `user_account_id` → `user_account.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_idempotency_key_created` | `idempotency_key (created_at)` | Purga de claves por antigüedad (retención de 24 h, K6). |
| `sqlite_autoindex_idempotency_key_1` | UNIQUE (user_account_id, idempotency_key) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_idempotency_key_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## B. Dinero, periodos y reparto

### `period` · Periodo

- **Entidad del plan (§3.4.1):** periodo
- **Descripción:** Mes calendario (periodo base del presupuesto) y sus dos quincenas (subperiodos). Los días de pago (15, y 30 o último día) son FIRMES (oráculo). Los rangos de la quincena (1 a 15 y 16 a fin de mes) son FIRMES (K8). Ambos se validan en la propia tabla.
- **Usado por:** S2 (presupuesto mensual, flujo por quincena, proyección). Pantallas Presupuesto e Inicio. Filtro «periodo» de C1.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `kind` | TEXT | no | month = mes calendario, fortnight = quincena | enumerado |
| `parent_id` | INTEGER | sí | Mes al que pertenece la quincena. NULL en meses | FK → `period.id` al borrar: CASCADE |
| `year` | INTEGER | no | Año | 2000 a 2100 |
| `month` | INTEGER | no | Mes | 1 a 12 |
| `half` | INTEGER | sí | Quincena 1 (días 1 a 15) o 2 (día 16 a fin de mes). NULL en meses | 1 o 2. Rangos FIRMES (K8) |
| `start_date` | TEXT | no | Primer día (ISO) | fecha válida |
| `end_date` | TEXT | no | Último día, inclusive (ISO) | fecha válida |
| `pay_date` | TEXT | sí | Día de pago de la quincena: 15, o 30 o el último día si el mes es más corto (28-feb-2027). NULL en meses | fecha válida |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Mes: del día 1 al último. Quincena 1: del 1 al 15, paga el 15. Quincena 2: del 16 al último, paga el 30 o el último día si es menor (febrero). Rangos de quincena FIRMES (K8).

**Relaciones:**
- `parent_id` → `period.id` · al borrar la madre: CASCADE (se borra con la madre)
- La referencian: `period.parent_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_period_parent` | `period (parent_id)` | Quincenas de un mes. |
| `ix_period_range` | `period (start_date, end_date)` | Ubica el periodo que contiene una fecha. |
| `ux_period_slot` | UNIQUE `period (kind, year, month, coalesce(half, 0))` | Un solo mes y una sola quincena 1 y 2 por año y mes. |

**Disparadores:**
- `trg_period_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `cost_share_rule` · Regla de reparto

- **Entidad del plan (§3.4.1):** regla_reparto
- **Descripción:** Cómo se divide un monto entre el hogar Axel–Ximena (a nombre del titular) y la participante externa: hogar, Karen o conjunto con su porcentaje. Las partes con porcentaje o monto fijo van en cost_share_rule_part y el resto cae en remainder_person_id.
- **Usado por:** Ajustes «Personas y reparto» y «Categorías». S2 (partes por gasto), S3 (partes de la compra a crédito). Filtro «familia» de C1.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `name` | TEXT | no | Nombre visible (Hogar, Karen, Conjunto 50/50) | no vacío |
| `remainder_person_id` | INTEGER | no | Persona que recibe el resto del monto y absorbe los centavos de redondeo (en el hogar, el titular) | FK → `person.id` al borrar: RESTRICT; no debe repetirse en las partes (lo valida la API) |
| `is_template` | INTEGER | no | 1 = regla reutilizable que se ofrece en capturas. 0 = regla de un solo uso (renglón «mixto» por cantidad) | def: 1; booleano |
| `is_default` | INTEGER | no | 1 = regla que aplica cuando ni el gasto ni su categoría traen regla (Hogar) | def: 0; booleano, una sola |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- La regla por defecto debe ser reutilizable.

**Relaciones:**
- `remainder_person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `category.default_cost_share_rule_id`, `cost_share_rule_part.rule_id`, `expense.cost_share_rule_id`, `expense_item.cost_share_rule_id`, `planned_item.cost_share_rule_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_cost_share_rule_remainder` | `cost_share_rule (remainder_person_id)` | Llave foránea hacia la persona del resto. |
| `ux_cost_share_rule_default` | UNIQUE `cost_share_rule (is_default) WHERE is_default = 1` | Una sola regla por defecto. |
| `ux_cost_share_rule_template_name` | UNIQUE `cost_share_rule (name COLLATE NOCASE) WHERE is_template = 1` | Nombres únicos entre las reglas reutilizables. |

**Disparadores:**
- `trg_cost_share_rule_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `cost_share_rule_part` · Parte de una regla de reparto

- **Entidad del plan (§3.4.1):** regla_reparto (detalle)
- **Descripción:** Parte de una persona dentro de una regla: un porcentaje o un monto fijo. Lo que no se asigna aquí va a la persona del resto.
- **Usado por:** Ajustes «Personas y reparto». S2 y S3 al calcular partes.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `rule_id` | INTEGER | no | Regla a la que pertenece | FK → `cost_share_rule.id` al borrar: CASCADE |
| `person_id` | INTEGER | no | Persona de la parte | FK → `person.id` al borrar: RESTRICT |
| `pct_bp` | INTEGER | sí | Porcentaje en puntos base: 5000 = 50.00 % (Conjunto 50/50: Karen 5000) | 0 a 10000, es decir 0 a 100 % |
| `fixed_cents` | INTEGER | sí | Monto fijo en centavos (parte fija de Karen en despensas futuras: 119150) | centavos |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Cada parte es porcentaje o monto fijo, nunca ambos ni ninguno.
- Una parte por persona en cada regla.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `rule_id` → `cost_share_rule.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_cost_share_rule_part_person` | `cost_share_rule_part (person_id)` | Llave foránea hacia la persona. |
| `sqlite_autoindex_cost_share_rule_part_1` | UNIQUE (rule_id, person_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_cost_share_rule_part_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `category` · Categoría

- **Entidad del plan (§3.4.1):** categoria
- **Descripción:** Categoría de gasto con jerarquía. Puede traer su regla de reparto por defecto (la despensa es Conjunto 50/50).
- **Usado por:** Ajustes «Categorías». S2 (presupuesto y real por categoría, alertas 80/100 %). Pantallas Gastos y Presupuesto. Filtro «categoría» de C1.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `parent_id` | INTEGER | sí | Categoría padre. NULL = raíz | FK → `category.id` al borrar: RESTRICT |
| `code` | TEXT | sí | Clave estable para categorías que usan los motores (groceries, eating_out). NULL en las del usuario | único; a-z, 0-9 y guion bajo |
| `name` | TEXT | no | Nombre visible (Despensa, Comida fuera) | no vacío |
| `kind` | TEXT | sí | Tipo de la categoría: fijo o variable (U1 N-281). Propone el tipo del gasto al registrarlo. NULL = sin clasificar | enumerado |
| `default_cost_share_rule_id` | INTEGER | sí | Reparto por categoría. NULL = la regla por defecto | FK → `cost_share_rule.id` al borrar: SET NULL |
| `sort_order` | INTEGER | no | Orden entre hermanas | def: 0 |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Una categoría no puede ser su propia madre.

**Relaciones:**
- `default_cost_share_rule_id` → `cost_share_rule.id` · al borrar la madre: SET NULL (queda sin liga)
- `parent_id` → `category.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `budget.category_id`, `category.parent_id`, `expense.category_id`, `planned_item.category_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_category_parent` | `category (parent_id, sort_order)` | Árbol de la barra lateral en orden. |
| `ix_category_rule` | `category (default_cost_share_rule_id)` | Llave foránea hacia la regla de reparto. |
| `sqlite_autoindex_category_1` | UNIQUE (code) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |
| `ux_category_sibling_name` | UNIQUE `category (coalesce(parent_id, 0), name COLLATE NOCASE)` | Nombres únicos entre hermanas, sin distinguir mayúsculas. |

**Disparadores:**
- `trg_category_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `budget` · Presupuesto

- **Entidad del plan (§3.4.1):** (complementaria) presupuesto mensual por categoría y total (S2, alertas 80/100 %)
- **Descripción:** Monto mensual planeado por categoría o total, con vigencia: rige desde un mes hasta otro (inclusive) o sin fin. El monto de un solo mes es un renglón que empieza y termina en ese mes. Sirve para «real contra plan» y para las alertas al 80 % y 100 %.
- **Usado por:** Ajustes «Presupuesto» (N-275, N-278) y pantalla Presupuesto (N-087, una sola fuente). S2 (real contra plan), S7 (alertas).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `category_id` | INTEGER | sí | Categoría. NULL = presupuesto total del mes | FK → `category.id` al borrar: RESTRICT |
| `amount_cents` | INTEGER | no | Monto planeado por mes en centavos | centavos |
| `valid_from_date` | TEXT | no | Primer día del primer mes en que rige | fecha válida, día 1 del mes |
| `valid_to_date` | TEXT | sí | Último día del último mes en que rige, inclusive. NULL = rige hasta que se cambie. Cada cambio hecho en Ajustes elige su vigencia (K2) con el alcance `scope` de updateBudget: «de este mes en adelante» (`onward`, preseleccionada y valor por omisión) deja NULL y elimina los renglones posteriores de la misma categoría, que el monto nuevo reemplaza (updateBudget, C1-datos.md §7.1); «solo este mes» (`this_month`) pone el último día de ese mes y el monto anterior sigue desde el mes siguiente | fecha válida, último día de un mes |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- La vigencia no termina antes de empezar.

**Relaciones:**
- `category_id` → `category.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_budget_category` | `budget (category_id, valid_from_date)` | Presupuesto vigente de una categoría en un mes e historial de cambios. |
| `ux_budget_category_from` | UNIQUE `budget (coalesce(category_id, 0), valid_from_date)` | Un renglón por categoría y mes de inicio (y uno solo de total por mes de inicio). |

**Disparadores:**
- `trg_budget_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## E. Tiendas, productos y presentaciones (se crean antes por dependencias)

### `store` · Tienda

- **Entidad del plan (§3.4.1):** tienda
- **Descripción:** Cadena (Walmart, Bodega Aurrera, Sam's Club), sucursal de una cadena, o mercado o tianguis con su día de visita. La membresía y el permiso de lectura asistida viven en la cadena.
- **Usado por:** Ajustes «Tiendas y precios» (mercados y membresía Sam's). S5 (mejor tienda, punto de equilibrio de la membresía, PROFECO). Pantalla Compras y tiendas. Filtro «tienda» de C1.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `kind` | TEXT | no | chain = cadena, branch = sucursal de una cadena, market = mercado o tianguis | enumerado |
| `parent_id` | INTEGER | sí | Cadena de la sucursal. Solo en branch | FK → `store.id` al borrar: RESTRICT |
| `name` | TEXT | no | Nombre: cadena (Bodega Aurrera), sucursal (Vallejo) o mercado | no vacío |
| `profeco_name` | TEXT | sí | Nombre en PROFECO para cruzar precios: cadena_comercial en cadenas, nombre_comercial en sucursales |  |
| `website_url` | TEXT | sí | Sitio de la cadena. Solo en chain |  |
| `assisted_reading_enabled` | INTEGER | no | Lectura asistida del sitio de la cadena. Queda en 0 hasta que docs/estado/S5-terminos.md tenga veredicto para ese sitio | def: 0; booleano, solo chain |
| `address` | TEXT | sí | Dirección (sucursal o mercado) |  |
| `municipality` | TEXT | sí | Municipio o alcaldía, para cruzar con PROFECO |  |
| `state` | TEXT | sí | Entidad federativa |  |
| `visit_weekday` | INTEGER | sí | Día de visita del mercado: 1 = lunes … 7 = domingo (ISO 8601) | obligatorio en market, prohibido en las demás |
| `requires_membership` | INTEGER | no | 1 = la cadena exige membresía (Sam's Club) | def: 0; booleano, solo chain |
| `has_membership` | INTEGER | no | 1 = la casa ya tiene la membresía (hoy Sam's = 0 y se evalúa) | def: 0; booleano |
| `membership_annual_fee_cents` | INTEGER | sí | Costo anual de la membresía en centavos (lo captura Axel en Ajustes) | centavos |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Solo las sucursales tienen cadena madre.
- Todo mercado tiene día de visita y solo los mercados lo tienen.
- Membresía, sitio y lectura asistida son datos de la cadena.
- Sin membresía requerida no hay costo ni membresía.

**Relaciones:**
- `parent_id` → `store.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `expense.store_id`, `price.store_id`, `promotion.store_id`, `shopping_list_item.store_id`, `store.parent_id`, `store_product_group.store_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_store_kind` | `store (kind, is_active)` | Listas de mercados o cadenas activas. |
| `ix_store_parent` | `store (parent_id)` | Sucursales de una cadena. |
| `ux_store_name` | UNIQUE `store (coalesce(parent_id, 0), name COLLATE NOCASE)` | Nombres únicos dentro de la misma cadena (y entre cadenas y mercados). |

**Disparadores:**
- `trg_store_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `product_group` · Grupo de producto

- **Entidad del plan (§3.4.1):** (complementaria) decisión D: dos verticales de compra (mercados para verduras y carne no especializada, cadenas para el resto)
- **Descripción:** Grupo de compra con su vertical: mercado o cadena. Permite repartir la lista de compras por mercado y por día de visita.
- **Usado por:** Ajustes «Tiendas y precios» e «Inventario». S4 (lista de compras por vertical), S5.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `name` | TEXT | no | Grupo (Verduras, Carne no especializada, Abarrotes) | único; único |
| `purchase_channel` | TEXT | no | Vertical: market = mercado o tianguis, chain = Walmart, Bodega Aurrera o Sam's Club | enumerado |
| `sort_order` | INTEGER | no | Orden en la lista de compras | def: 0 |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- Sin llaves foráneas salientes.
- La referencian: `product.product_group_id`, `store_product_group.product_group_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_product_group_1` | UNIQUE (name) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_product_group_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `store_product_group` · Qué se compra en cada tienda

- **Entidad del plan (§3.4.1):** tienda (detalle «qué se compra ahí»)
- **Descripción:** Grupos de producto que se compran en cada mercado (o tienda). Es el «qué se compra ahí» que Axel captura en Ajustes.
- **Usado por:** Ajustes «Tiendas y precios». S4 (asigna renglones de la lista a un mercado y su día).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `store_id` | INTEGER | no | Mercado o tienda | FK → `store.id` al borrar: CASCADE |
| `product_group_id` | INTEGER | no | Grupo que se compra ahí | FK → `product_group.id` al borrar: CASCADE |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Cada grupo una vez por tienda.

**Relaciones:**
- `product_group_id` → `product_group.id` · al borrar la madre: CASCADE (se borra con la madre)
- `store_id` → `store.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_store_product_group_group` | `store_product_group (product_group_id)` | Tiendas que surten un grupo. |
| `sqlite_autoindex_store_product_group_1` | UNIQUE (store_id, product_group_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_store_product_group_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `product` · Producto

- **Entidad del plan (§3.4.1):** producto
- **Descripción:** Producto genérico (ingrediente): Aceite de oliva, Arroz grano largo. Todas sus presentaciones son sustitutas entre sí. Define la unidad base de equivalencia.
- **Usado por:** Ajustes «Tiendas y precios», «Recetas» e «Inventario». S4 (ingredientes, kcal), S5 (equivalencias, mejor tienda por ingrediente).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `name` | TEXT | no | Nombre genérico en español | único; único |
| `product_group_id` | INTEGER | sí | Grupo y vertical de compra | FK → `product_group.id` al borrar: SET NULL |
| `base_unit` | TEXT | no | Unidad base para equivalencias: g, ml o pieza | enumerado |
| `grams_per_unit_milli` | INTEGER | sí | Gramos por unidad base ×1000 para calcular kcal: pieza = peso promedio, ml = densidad. NULL si base_unit = g |  |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- En productos por gramo el factor a gramos sobra.

**Relaciones:**
- `product_group_id` → `product_group.id` · al borrar la madre: SET NULL (queda sin liga)
- La referencian: `inventory_item.product_id`, `nutrition_fact.product_id`, `presentation.product_id`, `recipe_ingredient.product_id`, `shopping_list_item.product_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_product_group` | `product (product_group_id)` | Productos de un grupo. |
| `sqlite_autoindex_product_1` | UNIQUE (name) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_product_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `presentation` · Presentación

- **Entidad del plan (§3.4.1):** presentacion
- **Descripción:** Presentación comprable de un producto (marca y contenido), con su factor de equivalencia a la unidad base.
- **Usado por:** S4 (presentación comprable del ingrediente), S5 (precio por unidad base, sustitutos), Ajustes «Tiendas y precios».

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `product_id` | INTEGER | no | Producto genérico | FK → `product.id` al borrar: RESTRICT |
| `name` | TEXT | no | Presentación como se ve en tienda (Nutrioli Oli 500 ml) | no vacío |
| `brand` | TEXT | sí | Marca |  |
| `size_milli` | INTEGER | no | Contenido en milésimas de la unidad base: 500 ml = 500000, 750 g = 750000, 1 pieza = 1000 | factor de equivalencia, mayor que 0 |
| `is_bulk` | INTEGER | no | 1 = a granel: se vende por kilo o por unidad base | def: 0; booleano |
| `barcode` | TEXT | sí | Código de barras EAN o UPC (cruce con Open Food Facts) | 8 a 14 dígitos |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Nombres únicos dentro del mismo producto.

**Relaciones:**
- `product_id` → `product.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `expense_item.presentation_id`, `nutrition_fact.presentation_id`, `price.presentation_id`, `recipe_ingredient.presentation_id`, `shopping_list_item.presentation_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_presentation_1` | UNIQUE (product_id, name) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |
| `ux_presentation_barcode` | UNIQUE `presentation (barcode) WHERE barcode IS NOT NULL` | Un código de barras pertenece a una sola presentación. |

**Disparadores:**
- `trg_presentation_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## C. Créditos (catálogo, tasas y promociones)

### `credit` · Crédito

- **Entidad del plan (§3.4.1):** credito
- **Descripción:** Crédito del titular: línea revolvente tipo Kueski, tarjeta bancaria, préstamo personal a plazo fijo o crédito de tienda con abonos. Guarda el límite, la lectura A/B, las fechas de corte y pago y las comisiones.
- **Usado por:** Ajustes «Créditos y condiciones». S3 (calendario unificado, disponible A/B, optimizador). Pantalla Créditos. Filtro «crédito» de C1.
- **Usado por las columnas nuevas de 1.2:** Ajustes «Créditos y condiciones» (Axel los captura de su estado de cuenta, K18). S3 (fórmula general; sin datos muestra «falta capturar»). C1 v1.1. U4. Reglas RN-25 y RN-26.
- **Usado por las columnas nuevas de 1.4:** Ajustes «Créditos y condiciones» (P-26.05, «¿la tasa incluye IVA?»). S3 (`Credit.missing_terms` lista `vat` mientras sea NULL en `bank_card` y `store_credit`: 5 datos de K18; con 0 y sin porcentaje también lo lista, PROVISIONAL por el bloqueo abierto `S3-iva-lectura`, pregunta 2, RN-26). C1 v1.3. U4.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `name` | TEXT | no | Nombre visible (Kueski Pay) | único; único |
| `kind` | TEXT | no | revolving_line = línea revolvente tipo Kueski, bank_card = tarjeta bancaria, term_loan = préstamo personal a plazo fijo, store_credit = crédito de tienda con abonos | enumerado |
| `lender` | TEXT | sí | Institución que otorga el crédito |  |
| `holder_person_id` | INTEGER | no | Titular del crédito (Axel) | FK → `person.id` al borrar: RESTRICT |
| `credit_limit_cents` | INTEGER | sí | Límite de la línea (Kueski 1000000) | centavos, obligatorio en revolving_line y bank_card |
| `payment_frequency` | TEXT | no | Frecuencia de los pagos. fortnightly cae en los días 15 y 30 (o último día) | enumerado, monthly en tarjetas |
| `first_payment_min_days` | INTEGER | sí | Días mínimos entre la compra y el primer pago (Kueski 15, oráculo) | 0 a 60 |
| `availability_reading` | TEXT | no | Lectura del disponible (interruptor de Ajustes, solo aplica en revolving_line): B = el límite es el tope total (por defecto), A = el límite está libre en la fecha ancla y los pagos suman encima | def: 'B'; A o B |
| `availability_confirmed` | INTEGER | no | 0 = mostrar el aviso «pendiente confirmar» con el acreedor | def: 0; booleano |
| `reading_a_anchor_date` | TEXT | sí | Fecha «hoy» de la lectura A: lo comprado hasta ese día ya está dentro del límite libre (Kueski 2026-10-10) | fecha válida, solo revolving_line |
| `statement_day` | INTEGER | sí | Día de corte de la tarjeta | 1 a 31, solo bank_card |
| `payment_due_day` | INTEGER | sí | Día límite de pago de la tarjeta | 1 a 31, solo bank_card |
| `annual_rate_num` | INTEGER | sí | Numerador de la tasa anual ordinaria (tarjeta, préstamo) | junto con annual_rate_den |
| `annual_rate_den` | INTEGER | sí | Denominador de la tasa anual ordinaria | mayor que 0 |
| `vat_on_interest_bp` | INTEGER | sí | IVA que se suma a los intereses cuando la tasa no lo trae (1600 = 16 %). NULL = la tasa o la cotización ya incluye IVA (Kueski). Desde 1.4 (00004, Q2-67 opción 1) si la tasa ya incluye IVA lo dice `credit.vat_included`: este porcentaje conserva su significado y solo lleva valor cuando `vat_included` = 0; con `vat_included` = 1 es NULL, y con `vat_included` NULL el dato falta por capturar (RN-26). | 0 a 10000 |
| `annual_fee_cents` | INTEGER | sí | Comisión anual | centavos |
| `late_fee_cents` | INTEGER | sí | Comisión por pago tardío | centavos |
| `opening_fee_cents` | INTEGER | sí | Comisión por apertura | centavos |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `rate_base` | TEXT | sí | Nueva en 1.2 (K18). Base con que el estado de cuenta aplica la tasa anual: daily = diaria, monthly = mensual. NULL = falta capturar | enumerado, solo bank_card y store_credit |
| `min_payment_pct_bp` | INTEGER | sí | Nueva en 1.2 (K18). Porcentaje del pago mínimo que dice el estado de cuenta, en puntos base. NULL = falta capturar | 0 a 10000, solo bank_card y store_credit |
| `min_payment_fixed_cents` | INTEGER | sí | Nueva en 1.2 (K18). Mínimo fijo del pago mínimo que dice el estado de cuenta. NULL = falta capturar | centavos, solo bank_card y store_credit |
| `annual_rate_decimals` | INTEGER | sí | Nueva en 1.2 (K12). Decimales con que se capturó el porcentaje de la tasa anual; la pantalla lo muestra así y los cálculos usan annual_rate_num/annual_rate_den. NULL = sin tasa o sin dato | entero mayor o igual a 0, solo con tasa anual |
| `vat_included` | INTEGER | sí | Nueva en 1.4 (Q2-67 opción 1, K18). ¿La tasa anual o la cotización ya incluye IVA? 1 = sí, no se suma IVA aparte; 0 = no, se suma el de vat_on_interest_bp; NULL = falta capturar (S3 lo lista como `vat` en missing_terms de `bank_card` y `store_credit`). Los créditos que ya existían al migrar pasan con el significado de 1.3 (porcentaje con valor → 0; NULL → 1). Con 0 y vat_on_interest_bp NULL falta el porcentaje: PROVISIONAL (bloqueo abierto S3-iva-lectura, pregunta 2) `vat` también va en missing_terms y S3 dice «falta capturar»; no calcula con 0 % ni inventa un porcentaje (K18) | booleano o NULL; con 1 o NULL, vat_on_interest_bp es NULL (solo lleva valor con 0); con 0 el CHECK admite el porcentaje en NULL |

**Reglas de tabla (CHECK / UNIQUE):**
- Línea revolvente y tarjeta tienen límite.
- Solo la tarjeta tiene día de corte y día de pago, y siempre los tiene.
- La tarjeta se paga mensualmente.
- La fecha ancla de la lectura A solo existe en líneas revolventes.
- Con lectura A hace falta la fecha ancla.
- La tasa anual va completa (numerador y denominador) o no va.

**Relaciones:**
- `holder_person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `credit_payment.credit_id`, `credit_purchase.credit_id`, `credit_statement.credit_id`, `credit_term_rate.credit_id`, `planned_item.credit_id`, `promotion.credit_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_credit_holder` | `credit (holder_person_id)` | Créditos del titular. |
| `ix_credit_kind` | `credit (kind, is_active)` | Créditos activos por tipo. |
| `sqlite_autoindex_credit_1` | UNIQUE (name) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_credit_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `credit_term_rate` · Tasa por plazo

- **Entidad del plan (§3.4.1):** tasa_por_plazo
- **Descripción:** Tabla editable de tasas por plazo de un crédito, como fracción exacta del capital por todo el plazo con IVA incluido. Kueski: 2 quincenas = 53450/309450 (17.27 %), 4 quincenas = 80050/309450 (25.8685 %).
- **Usado por:** Ajustes «Créditos y condiciones». S3 (estimar intereses de compras planeadas, escenarios del optimizador).
- **Usado por las columnas nuevas de 1.2:** C1 v1.1 (CreditTermRate muestra la tasa con los decimales capturados). U4 (tabla de tasas por plazo). Regla RN-25.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_id` | INTEGER | no | Crédito | FK → `credit.id` al borrar: CASCADE |
| `term_count` | INTEGER | no | Plazo en número de pagos (2, 4, 6, 8, 10 o 12 quincenas en Kueski) | 1 a 120 |
| `rate_num` | INTEGER | no | Numerador: si viene de una cotización, los intereses en centavos (80050) | entero mayor o igual a 0 |
| `rate_den` | INTEGER | no | Denominador: si viene de una cotización, el capital en centavos (309450). Una tasa tecleada como 25.8685 % se guarda 258685/1000000 | mayor que 0 |
| `source` | TEXT | no | quote = sale de una cotización real, manual = tecleada en la tabla editable | enumerado |
| `quote_date` | TEXT | sí | Fecha de la cotización (2026-10-10) | obligatoria si source = quote |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `rate_decimals` | INTEGER | sí | Nueva en 1.2 (K12). Decimales con que se capturó o se leyó de la cotización el porcentaje (17.27 % → 2; 25.8685 % → 4). La pantalla redondea la fracción a esos decimales con HALF_UP; los cálculos usan rate_num/rate_den. NULL = sin dato | entero mayor o igual a 0 |

**Reglas de tabla (CHECK / UNIQUE):**
- Una tasa vigente por plazo y crédito.
- Una tasa de cotización dice de qué fecha es.

**Relaciones:**
- `credit_id` → `credit.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_credit_term_rate_1` | UNIQUE (credit_id, term_count) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_credit_term_rate_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `promotion` · Promoción

- **Entidad del plan (§3.4.1):** promocion
- **Descripción:** Calendario de promociones capturado a mano: comercio, plazo, tasa (0 % casi siempre), compra mínima y vigencia.
- **Usado por:** Ajustes «Créditos y condiciones» (promociones). S3 (escenarios con promoción), S5 (comprar más por peso).
- **Usado por las columnas nuevas de 1.2:** C1 v1.1 (tasa promocional como se capturó). U4. Regla RN-25.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_id` | INTEGER | sí | Crédito al que aplica (Kueski). NULL = cualquiera | FK → `credit.id` al borrar: CASCADE |
| `merchant_name` | TEXT | no | Comercio tal como lo anuncia la promoción (Chedraui) | no vacío |
| `store_id` | INTEGER | sí | Tienda registrada que corresponde al comercio, si existe | FK → `store.id` al borrar: SET NULL |
| `term_count` | INTEGER | no | Plazo en pagos (4 quincenas) | 1 a 120 |
| `rate_num` | INTEGER | no | Numerador de la tasa promocional por todo el plazo. 0 = sin intereses | def: 0; mayor o igual a 0 |
| `rate_den` | INTEGER | no | Denominador de la tasa promocional | def: 1; mayor que 0 |
| `min_purchase_cents` | INTEGER | no | Compra mínima (200000 a 250000 en el tablero 2026) | def: 0; centavos |
| `valid_from` | TEXT | no | Inicio de vigencia, inclusive | fecha válida |
| `valid_to` | TEXT | no | Fin de vigencia, inclusive (Chedraui: solo 2026-10-30) | fecha válida |
| `source_note` | TEXT | sí | De dónde se tomó (tablero oficial de promociones Kueski 2026) |  |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `rate_decimals` | INTEGER | sí | Nueva en 1.2 (K12). Decimales con que se capturó la tasa promocional (0 % → 0). NULL = sin dato | entero mayor o igual a 0 |

**Reglas de tabla (CHECK / UNIQUE):**
- La vigencia termina el mismo día o después de empezar.

**Relaciones:**
- `store_id` → `store.id` · al borrar la madre: SET NULL (queda sin liga)
- `credit_id` → `credit.id` · al borrar la madre: CASCADE (se borra con la madre)
- La referencian: `credit_purchase.promotion_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_promotion_credit` | `promotion (credit_id, valid_from)` | Promociones de un crédito. |
| `ix_promotion_store` | `promotion (store_id)` | Llave foránea hacia la tienda. |
| `ix_promotion_validity` | `promotion (valid_from, valid_to)` | Promociones vigentes en una fecha. |

**Disparadores:**
- `trg_promotion_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## B (cont.). Plan, ingresos y gastos

### `planned_item` · Partida planificada

- **Entidad del plan (§3.4.1):** ingreso (fijo quincenal) y gasto (fijo y venidero): parte recurrente. Coincide con datos/semilla/gastos_fijos_ingresos.csv (plan §8.2)
- **Descripción:** Ingreso o gasto esperado con su periodicidad: sueldo quincenal, renta, despensa bimestral con Kueski, gastos venideros. Los motores proyectan sus ocurrencias, no se guardan una por una.
- **Usado por:** Ajustes «Ingresos» y «Presupuesto». S2 (plan del mes, proyección 12 meses), S3 (compras planeadas con crédito). Importador del CSV de Axel.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `kind` | TEXT | no | Tipo (columna «tipo» del CSV): ingreso, gasto fijo o gasto venidero | enumerado |
| `description` | TEXT | no | Concepto | no vacío |
| `amount_cents` | INTEGER | no | Monto de cada ocurrencia en centavos (sueldo por quincena, renta por mes). En quincenal es el de la quincena del 15, y también el del 30 si second_half_amount_cents es NULL | centavos |
| `second_half_amount_cents` | INTEGER | sí | Monto de la quincena del 30 (o último día del mes) cuando difiere del de la quincena del 15 (N-279). NULL = las dos quincenas iguales | centavos, solo en frequency = fortnightly |
| `frequency` | TEXT | no | Periodicidad. fortnightly cae en los días de pago 15 y 30 (o último día). Las demás se repiten desde start_date | enumerado |
| `interval_count` | INTEGER | no | Cada cuántas unidades se repite (despensa bimestral = monthly cada 2) | def: 1; 1 a 24 |
| `start_date` | TEXT | no | Primera ocurrencia o fecha única. Su día del mes es el ancla y se recorta al último día en meses cortos | fecha válida |
| `end_date` | TEXT | sí | Última fecha posible, inclusive. NULL = sin fin | fecha válida |
| `category_id` | INTEGER | sí | Categoría del gasto | FK → `category.id` al borrar: RESTRICT; obligatoria en gastos, NULL en ingresos |
| `person_id` | INTEGER | sí | Persona a quien se atribuye. En ingresos solo el titular (disparador). NULL = hogar en general | FK → `person.id` al borrar: RESTRICT |
| `cost_share_rule_id` | INTEGER | sí | Reparto del gasto. NULL = el de la categoría o el de por defecto | FK → `cost_share_rule.id` al borrar: RESTRICT |
| `payment_method` | TEXT | sí | Cómo se pagará cada ocurrencia. NULL = sin definir | enumerado, NULL en ingresos |
| `credit_id` | INTEGER | sí | Crédito con que se planea pagar (despensa con Kueski) | FK → `credit.id` al borrar: RESTRICT; obligatorio si payment_method = credit |
| `term_count` | INTEGER | sí | Plazo planeado en pagos (4 quincenas) | obligatorio con credit_id |
| `accrual_months` | INTEGER | no | Meses en que se devenga cada ocurrencia (despensa bimestral = 2) | def: 1; 1 a 24 |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- La vigencia no termina antes de empezar.
- Una partida única ocurre solo en start_date.
- Un ingreso no lleva categoría, reparto, medio de pago ni devengo.
- Todo gasto planificado tiene categoría.
- Hay crédito si y solo si el medio de pago es crédito.
- El monto de la segunda quincena solo existe en partidas quincenales.
- El plazo acompaña siempre al crédito.

**Relaciones:**
- `credit_id` → `credit.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `cost_share_rule_id` → `cost_share_rule.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `category_id` → `category.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `expense.planned_item_id`, `income.planned_item_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_planned_item_category` | `planned_item (category_id)` | Plan por categoría. |
| `ix_planned_item_credit` | `planned_item (credit_id)` | Compras planeadas con un crédito. |
| `ix_planned_item_kind` | `planned_item (kind, is_active, start_date)` | Partidas activas por tipo y fecha (proyección). |
| `ix_planned_item_person` | `planned_item (person_id)` | Plan por persona. |
| `ix_planned_item_rule` | `planned_item (cost_share_rule_id)` | Llave foránea hacia la regla de reparto. |

**Disparadores:**
- `trg_planned_item_income_holder_ins`: rechaza con «planned_item: un ingreso planificado solo puede ser del titular (decisión B)».
- `trg_planned_item_income_holder_upd`: rechaza con «planned_item: un ingreso planificado solo puede ser del titular (decisión B)».
- `trg_planned_item_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `income` · Ingreso

- **Entidad del plan (§3.4.1):** ingreso
- **Descripción:** Ingreso recibido por el titular: fijo quincenal (lo captura Axel cada quincena) o extra capturado a mano. Ximena no aporta ingresos al modelo.
- **Usado por:** Ajustes «Ingresos». S2 (real contra plan, flujo por quincena). Pantallas Inicio y Presupuesto.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `person_id` | INTEGER | no | Quién recibe el ingreso | FK → `person.id` al borrar: RESTRICT; solo el titular (disparador) |
| `kind` | TEXT | no | fixed = ingreso fijo quincenal, extra = ingreso extra | enumerado |
| `planned_item_id` | INTEGER | sí | Partida planificada que cumple (sueldo quincenal). NULL en extras sin plan | FK → `planned_item.id` al borrar: SET NULL |
| `description` | TEXT | no | Concepto | no vacío |
| `amount_cents` | INTEGER | no | Monto recibido en centavos | mayor que 0 |
| `income_date` | TEXT | no | Fecha en que se recibió (calendario de America/Mexico_City) | fecha válida |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `planned_item_id` → `planned_item.id` · al borrar la madre: SET NULL (queda sin liga)
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_income_date` | `income (income_date)` | Ingresos de un periodo (rango de fechas). |
| `ix_income_person_date` | `income (person_id, income_date)` | Ingresos por persona y fecha. |
| `ix_income_planned` | `income (planned_item_id)` | Ingresos que cumplen una partida. |

**Disparadores:**
- `trg_income_holder_only_ins`: rechaza con «income: solo se registran ingresos del titular (decisión B)».
- `trg_income_holder_only_upd`: rechaza con «income: solo se registran ingresos del titular (decisión B)».
- `trg_income_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `expense` · Gasto

- **Entidad del plan (§3.4.1):** gasto
- **Descripción:** Gasto real: fijo, variable o venidero. Guarda la fecha que inicia el devengo, los meses de devengo, el medio de pago (si es crédito el efectivo sale por los pagos programados), quién pagó, a quién se atribuye y su reparto.
- **Usado por:** Pantallas Gastos (registro en 3 toques), Inicio y Presupuesto. S2 (devengado contra efectivo, real contra plan), S3 (compras a crédito), S7 (reportes).
- **Usado por las columnas nuevas de 1.2:** C1 v1.1 (capturas atrasadas con Idempotency-Key y fecha de captura). S6 (busca la clave antes de insertar). PWA en WF-3 (cola del teléfono, estado «Sin conexión»). Regla RN-24.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `kind` | TEXT | no | fixed = gasto fijo, variable = gasto del día a día, upcoming = gasto venidero que ya ocurrió o quedó comprometido | enumerado |
| `description` | TEXT | sí | Concepto o comercio. NULL = sin descripción: el registro en 3 toques no la pide (N-054) y la pantalla muestra la categoría | opcional; si viene, no vacía |
| `category_id` | INTEGER | no | Categoría | FK → `category.id` al borrar: RESTRICT |
| `amount_cents` | INTEGER | no | Importe de lo comprado en centavos, con IVA y descuentos y sin intereses de financiamiento (despensa 1 = 309450) | centavos |
| `expense_date` | TEXT | no | Fecha de la compra o consumo. El devengo empieza en su mes | fecha válida |
| `accrual_months` | INTEGER | no | Meses calendario en que se devenga desde el mes de expense_date (despensa bimestral = 2, mitad por mes) | def: 1; 1 a 24 |
| `payment_method` | TEXT | no | Medio de pago. credit = el efectivo sale según los pagos programados de su credit_purchase | enumerado |
| `paid_by_person_id` | INTEGER | sí | Quién pagó de su bolsillo (papel de contado de Karen). NULL = el titular | FK → `person.id` al borrar: RESTRICT; NULL si payment_method = credit |
| `person_id` | INTEGER | sí | Persona a quien se atribuye (comida fuera de Ximena). NULL = hogar en general | FK → `person.id` al borrar: RESTRICT |
| `cost_share_rule_id` | INTEGER | sí | Reparto del gasto. NULL = la regla de la categoría o la de por defecto | FK → `cost_share_rule.id` al borrar: RESTRICT |
| `store_id` | INTEGER | sí | Tienda donde se compró | FK → `store.id` al borrar: RESTRICT |
| `planned_item_id` | INTEGER | sí | Partida planificada que cumple (gasto fijo o venidero) | FK → `planned_item.id` al borrar: SET NULL |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `idempotency_key` | TEXT | sí | Nueva en 1.2 (K9). Clave que el teléfono genera para cada captura de su cola, igual al encabezado Idempotency-Key. Si llega otra vez, la API devuelve este gasto y no crea otro. A diferencia de la tabla idempotency_key (24 h), no se purga | 16 a 128 caracteres: letras sin acento, dígitos, guion y guion bajo; única (ux_expense_idempotency_key); NULL = sin clave |
| `captured_at` | TEXT | sí | Nueva en 1.2 (K9). Instante en que la persona capturó el gasto en el teléfono (UTC). En una captura atrasada es anterior a created_at. NULL = captura en línea: la hora de captura es created_at | instante ISO 8601 con «Z» |

**Reglas de tabla (CHECK / UNIQUE):**
- En compras a crédito quién paga lo dicen las partes de cada pago programado.

**Relaciones:**
- `planned_item_id` → `planned_item.id` · al borrar la madre: SET NULL (queda sin liga)
- `store_id` → `store.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `cost_share_rule_id` → `cost_share_rule.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `paid_by_person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `category_id` → `category.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `attachment.expense_id`, `credit_purchase.expense_id`, `expense_item.expense_id`, `expense_share.expense_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_expense_category_date` | `expense (category_id, expense_date)` | Gastos por categoría y periodo (real contra presupuesto). |
| `ix_expense_date` | `expense (expense_date)` | Gastos de un periodo (rango de fechas). |
| `ix_expense_paid_by` | `expense (paid_by_person_id)` | Gastos pagados por una persona de su bolsillo. |
| `ix_expense_person_date` | `expense (person_id, expense_date)` | Gastos por persona y periodo (comida fuera por persona). |
| `ix_expense_planned` | `expense (planned_item_id)` | Gastos que cumplen una partida planificada. |
| `ix_expense_rule` | `expense (cost_share_rule_id)` | Llave foránea hacia la regla de reparto. |
| `ix_expense_store_date` | `expense (store_id, expense_date)` | Gastos por tienda y periodo. |
| `ux_expense_idempotency_key` | UNIQUE `expense (idempotency_key)` | Una captura de la cola no se duplica: un solo gasto por clave. |

**Disparadores:**
- `trg_expense_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `expense_item` · Renglón de gasto

- **Entidad del plan (§3.4.1):** gasto (detalle: renglones del ticket)
- **Descripción:** Renglón de un gasto con ticket (la despensa): presentación, cantidad, importe y reparto del renglón (cuenta axel, karen, conjunto o mixto de data.json).
- **Usado por:** S2 y S3 (el crédito se reparte según lo que cada quien compró), S5 (precio desde ticket), S4 (entrada al inventario). Pantalla Gastos (detalle).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `expense_id` | INTEGER | no | Gasto al que pertenece | FK → `expense.id` al borrar: CASCADE |
| `line_no` | INTEGER | no | Orden del renglón en el ticket | mayor o igual a 1 |
| `description` | TEXT | no | Texto del renglón como viene en el ticket | no vacío |
| `presentation_id` | INTEGER | sí | Presentación del catálogo. NULL mientras no se liga | FK → `presentation.id` al borrar: RESTRICT |
| `quantity_milli` | INTEGER | no | Cantidad de presentaciones ×1000: 3 piezas = 3000, 1.250 kg a granel = 1250 | mayor que 0 |
| `unit_price_cents` | INTEGER | sí | Precio por presentación en centavos | centavos |
| `line_total_cents` | INTEGER | no | Importe del renglón con descuentos e IVA. La suma de renglones es el importe del gasto | centavos |
| `cost_share_rule_id` | INTEGER | sí | Reparto del renglón (Hogar, Karen, Conjunto 50/50 o una regla de un solo uso para «mixto»). NULL = el del gasto | FK → `cost_share_rule.id` al borrar: RESTRICT |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Un número de renglón por gasto.

**Relaciones:**
- `cost_share_rule_id` → `cost_share_rule.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `presentation_id` → `presentation.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `expense_id` → `expense.id` · al borrar la madre: CASCADE (se borra con la madre)
- La referencian: `inventory_movement.expense_item_id`, `price.expense_item_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_expense_item_presentation` | `expense_item (presentation_id)` | Compras de una presentación (historial e inventario). |
| `ix_expense_item_rule` | `expense_item (cost_share_rule_id)` | Llave foránea hacia la regla de reparto. |
| `sqlite_autoindex_expense_item_1` | UNIQUE (expense_id, line_no) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_expense_item_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `expense_share` · Parte de un gasto

- **Entidad del plan (§3.4.1):** regla_reparto (resultado congelado por gasto)
- **Descripción:** Parte de cada persona en un gasto, ya calculada al centavo. Se congela al guardar para que cambiar una regla no altere gastos pasados. La parte del hogar Axel–Ximena va a nombre del titular.
- **Usado por:** S2 (gasto del hogar contra parte de Karen), S7 (estado de cuenta de Karen). Filtros «persona» y «familia» de C1.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `expense_id` | INTEGER | no | Gasto | FK → `expense.id` al borrar: CASCADE |
| `person_id` | INTEGER | no | Persona que asume la parte | FK → `person.id` al borrar: RESTRICT |
| `amount_cents` | INTEGER | no | Parte en centavos (despensa 1: Axel 190300, Karen 119150). La suma de partes es el importe del gasto | centavos |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Una parte por persona y gasto.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `expense_id` → `expense.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_expense_share_person` | `expense_share (person_id)` | Partes de una persona (consultas por persona y familia). |
| `sqlite_autoindex_expense_share_1` | UNIQUE (expense_id, person_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_expense_share_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `attachment` · Adjunto

- **Entidad del plan (§3.4.1):** gasto (adjunto)
- **Descripción:** Archivo adjunto: foto de ticket (con su OCR local), factura o estado de cuenta del banco. Ninguna cifra del OCR entra al sistema sin confirmación humana.
- **Usado por:** Pantalla Gastos (adjuntar ticket). S5 (OCR Tesseract con confirmación), S6 (almacenamiento de adjuntos).
- **Usado por las columnas nuevas de 1.2:** C1 v1.1 (guardar y leer el borrador del ticket). S5 (lectura de tickets). U1 P-05 (N-071, N-075 y N-076). Regla RN-23.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `expense_id` | INTEGER | sí | Gasto al que pertenece. NULL mientras el ticket espera confirmación | FK → `expense.id` al borrar: SET NULL |
| `kind` | TEXT | no | Foto de ticket, factura, estado de cuenta del banco u otro | enumerado |
| `file_path` | TEXT | no | Ruta relativa dentro de la carpeta de datos del sistema | único; única |
| `original_name` | TEXT | sí | Nombre original del archivo |  |
| `mime_type` | TEXT | no | Tipo MIME (image/jpeg, application/pdf) |  |
| `size_bytes` | INTEGER | no | Tamaño en bytes | mayor o igual a 0 |
| `sha256` | TEXT | no | Huella del archivo (detecta duplicados y corrupción) | 64 caracteres |
| `ocr_status` | TEXT | no | Estado del OCR local: none, pending, done (falta confirmar), confirmed, failed | def: 'none'; enumerado |
| `ocr_engine` | TEXT | sí | Motor y versión (tesseract 5.5.3 spa) |  |
| `ocr_text` | TEXT | sí | Texto crudo reconocido |  |
| `uploaded_by_user_id` | INTEGER | sí | Cuenta que lo subió | FK → `user_account.id` al borrar: SET NULL |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `ocr_draft_json` | TEXT | sí | Nueva en 1.2 (K3). Borrador de las correcciones del ticket que la persona guardó con «Guardar para después», en JSON con la forma que define C1 v1.1. La lectura cruda (ocr_text) no se toca. Se vacía al confirmar el ticket | objeto JSON válido; NULL = sin borrador |
| `ocr_draft_saved_at` | TEXT | sí | Nueva en 1.2 (K3). Cuándo se guardó el borrador por última vez (UTC) | instante ISO 8601 con «Z»; va junto con ocr_draft_json |

**Relaciones:**
- `uploaded_by_user_id` → `user_account.id` · al borrar la madre: SET NULL (queda sin liga)
- `expense_id` → `expense.id` · al borrar la madre: SET NULL (queda sin liga)
- La referencian: `credit_statement.attachment_id`, `price.attachment_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_attachment_expense` | `attachment (expense_id)` | Adjuntos de un gasto. |
| `ix_attachment_ocr_queue` | `attachment (ocr_status, created_at) WHERE ocr_status IN ('pending', 'done')` | Cola de tickets con OCR pendiente de procesar o de confirmar. |
| `ix_attachment_sha` | `attachment (sha256)` | Detecta archivos repetidos por su huella. |
| `ix_attachment_uploader` | `attachment (uploaded_by_user_id)` | Llave foránea hacia la cuenta que subió. |
| `sqlite_autoindex_attachment_1` | UNIQUE (file_path) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_attachment_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `price` · Precio

- **Entidad del plan (§3.4.1):** precio
- **Descripción:** Precio observado de una presentación en una tienda, con fecha y fuente autorizada. Un precio es viejo a los 30 días y pesa menos en la comparativa.
- **Usado por:** S5 (mejor tienda por ingrediente y canasta, frescura, membresía), S4 (costo de receta), S7 (alerta de precios viejos). Pantalla Compras y tiendas.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `presentation_id` | INTEGER | no | Presentación | FK → `presentation.id` al borrar: RESTRICT |
| `store_id` | INTEGER | no | Tienda del precio (sucursal, cadena o mercado) | FK → `store.id` al borrar: RESTRICT |
| `price_cents` | INTEGER | no | Precio de una presentación en centavos con IVA (Nutrioli Oli 500 ml en Bodega Aurrera: 15000) | centavos |
| `observed_date` | TEXT | no | Fecha del precio. Es viejo cuando tiene 30 días o más | fecha válida |
| `source` | TEXT | no | Fuente: ticket capturado a mano, ticket por OCR ya confirmado, lectura asistida del sitio, PROFECO u otra captura manual | enumerado |
| `is_promotional` | INTEGER | no | 1 = precio de promoción (sexta lata de atún a 2100), no se toma como precio regular | def: 0; booleano |
| `source_ref` | TEXT | sí | Referencia de la fuente (mes del dataset PROFECO, URL) |  |
| `expense_item_id` | INTEGER | sí | Renglón de ticket que originó el precio | FK → `expense_item.id` al borrar: SET NULL |
| `attachment_id` | INTEGER | sí | Foto del ticket | FK → `attachment.id` al borrar: SET NULL |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- No se repite el mismo precio el mismo día con la misma fuente.

**Relaciones:**
- `attachment_id` → `attachment.id` · al borrar la madre: SET NULL (queda sin liga)
- `expense_item_id` → `expense_item.id` · al borrar la madre: SET NULL (queda sin liga)
- `store_id` → `store.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `presentation_id` → `presentation.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_price_attachment` | `price (attachment_id)` | Llave foránea hacia la foto del ticket. |
| `ix_price_date` | `price (observed_date)` | Precios por fecha (frescura). |
| `ix_price_item` | `price (expense_item_id)` | Llave foránea hacia el renglón de ticket. |
| `ix_price_latest` | `price (presentation_id, store_id, observed_date)` | Último precio de una presentación en una tienda. |
| `ix_price_store_date` | `price (store_id, observed_date)` | Precios de una tienda por fecha (canasta, precios viejos). |
| `sqlite_autoindex_price_1` | UNIQUE (presentation_id, store_id, observed_date, source, is_promotional) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_price_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## C (cont.). Compras a crédito, pagos programados y pagos hechos

### `credit_purchase` · Compra a crédito

- **Entidad del plan (§3.4.1):** compra_credito
- **Descripción:** Compra o disposición con un crédito: capital, intereses cotizados (IVA incluido), plazo, primer pago y promoción. Solo el capital usa y repone la línea.
- **Usado por:** S3 (calendario, disponible A/B, compras vivas, regresión Kueski), S2 (devengado incluye intereses), S7 (estado de cuenta de Karen). Pantalla Créditos.
- **Usado por las columnas nuevas de 1.2:** C1 v1.1 (tasa de la compra como se capturó o cotizó). U4. Regla RN-25.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_id` | INTEGER | no | Crédito usado | FK → `credit.id` al borrar: RESTRICT |
| `expense_id` | INTEGER | sí | Gasto que originó la compra (despensa 1). NULL en disposiciones sin gasto (desembolso de préstamo) | FK → `expense.id` al borrar: RESTRICT; único; un gasto tiene a lo más una compra a crédito |
| `kind` | TEXT | no | installment_plan = pagos fijos (Kueski, meses sin intereses, abonos de tienda), single_charge = cargo normal de tarjeta, loan_disbursement = préstamo con tabla de amortización | enumerado |
| `description` | TEXT | sí | Descripción corta |  |
| `purchase_date` | TEXT | no | Fecha de la compra (2026-10-10) | fecha válida |
| `principal_cents` | INTEGER | no | Capital: lo que usa la línea (309450) | centavos |
| `interest_cents` | INTEGER | no | Intereses de todo el plazo, IVA incluido, tal como los cotizó el acreedor (80050). Nunca reponen la línea | def: 0; centavos |
| `total_cents` | INTEGER | no | Total a pagar (389500) | igual a capital más intereses |
| `term_count` | INTEGER | no | Número de pagos (4) | 1 a 360 |
| `payment_frequency` | TEXT | no | Frecuencia de los pagos (quincenal) | enumerado |
| `first_payment_date` | TEXT | no | Primer pago (2026-10-30) | no antes de la compra |
| `interest_source` | TEXT | no | Origen de los intereses: quoted = cotización del acreedor al pagar, term_table = tabla por plazo, promotion = promoción, amortization = tabla de amortización | enumerado |
| `rate_num` | INTEGER | sí | Numerador de la tasa de esta compra (cotización: intereses/capital = 80050/309450) | junto con rate_den y rate_basis |
| `rate_den` | INTEGER | sí | Denominador de la tasa | mayor que 0 |
| `rate_basis` | TEXT | sí | term_total = sobre el capital por todo el plazo (Kueski), annual = tasa anual (préstamo) | enumerado |
| `interest_includes_vat` | INTEGER | no | 1 = los intereses ya traen IVA (Kueski) | def: 1; booleano |
| `promotion_id` | INTEGER | sí | Promoción aplicada | FK → `promotion.id` al borrar: SET NULL |
| `status` | TEXT | no | active = con pagos pendientes (compra viva), paid_off = liquidada, cancelled = cancelada | def: 'active'; enumerado |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `rate_decimals` | INTEGER | sí | Nueva en 1.2 (K12). Decimales con que se muestra la tasa de esta compra, los mismos con que se capturó o cotizó. NULL = sin tasa o sin dato | entero mayor o igual a 0, solo si la compra lleva tasa |

**Reglas de tabla (CHECK / UNIQUE):**
- Total = capital + intereses.
- El primer pago no es anterior a la compra.
- La tasa va completa (numerador, denominador y base) o no va.
- Si los intereses salen de una promoción, la promoción se identifica.
- Un cargo normal de tarjeta se paga en un solo pago.

**Relaciones:**
- `promotion_id` → `promotion.id` · al borrar la madre: SET NULL (queda sin liga)
- `expense_id` → `expense.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `credit_id` → `credit.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `credit_purchase_share.credit_purchase_id`, `scheduled_payment.credit_purchase_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_credit_purchase_credit_date` | `credit_purchase (credit_id, purchase_date)` | Compras de un crédito por fecha. |
| `ix_credit_purchase_credit_status` | `credit_purchase (credit_id, status)` | Compras vivas de un crédito (restricción de máximo de compras vivas). |
| `ix_credit_purchase_promotion` | `credit_purchase (promotion_id)` | Llave foránea hacia la promoción. |
| `sqlite_autoindex_credit_purchase_1` | UNIQUE (expense_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_credit_purchase_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `credit_purchase_share` · Parte de una compra a crédito

- **Entidad del plan (§3.4.1):** compra_credito (partes por persona)
- **Descripción:** Capital, intereses y total de cada persona en una compra a crédito. El capital sale de lo que cada quien compró y los intereses se reparten en proporción al capital, con el resto para el titular.
- **Usado por:** S3 (pagos por persona), S2, S7 (estado de cuenta de Karen).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_purchase_id` | INTEGER | no | Compra | FK → `credit_purchase.id` al borrar: CASCADE |
| `person_id` | INTEGER | no | Persona | FK → `person.id` al borrar: RESTRICT |
| `principal_cents` | INTEGER | no | Capital de la persona (Axel 190300, Karen 119150) | centavos |
| `interest_cents` | INTEGER | no | Intereses de la persona (Axel 49228, Karen 30822) | centavos |
| `total_cents` | INTEGER | no | Total de la persona (Axel 239528, Karen 149972) | igual a capital más intereses |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Total de la persona = su capital + sus intereses.
- Una parte por persona y compra.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `credit_purchase_id` → `credit_purchase.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_credit_purchase_share_person` | `credit_purchase_share (person_id)` | Compras a crédito de una persona. |
| `sqlite_autoindex_credit_purchase_share_1` | UNIQUE (credit_purchase_id, person_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_credit_purchase_share_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `scheduled_payment` · Pago programado

- **Entidad del plan (§3.4.1):** pago_programado
- **Descripción:** Cada pago de la serie de una compra: fecha, monto, capital que repone e intereses. Se guarda ya calculado al centavo (el último de la serie ajusta centavos).
- **Usado por:** S3 (calendario unificado de pagos, disponible A/B), S2 (efectivo por quincena), S7 (alerta de pago próximo). Pantallas Créditos e Inicio (calendario).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_purchase_id` | INTEGER | no | Compra a la que pertenece | FK → `credit_purchase.id` al borrar: CASCADE |
| `seq` | INTEGER | no | Número de pago dentro de la serie (1 a term_count) | mayor o igual a 1 |
| `due_date` | TEXT | no | Fecha de pago (2026-10-30, 2026-11-15, 2026-11-30, 2026-12-15) | fecha válida |
| `amount_cents` | INTEGER | no | Monto del pago (97375) | centavos |
| `principal_cents` | INTEGER | no | Capital que repone la línea al pagarse (77363, el último 77361) | centavos |
| `interest_cents` | INTEGER | no | Intereses del pago, no reponen la línea (20012, el último 20014) | centavos |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Monto = capital + intereses.
- Un número de pago por compra.

**Relaciones:**
- `credit_purchase_id` → `credit_purchase.id` · al borrar la madre: CASCADE (se borra con la madre)
- La referencian: `credit_payment.scheduled_payment_id`, `scheduled_payment_share.scheduled_payment_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_scheduled_payment_due` | `scheduled_payment (due_date)` | Calendario unificado: pagos por fecha (pago próximo, flujo por quincena). |
| `sqlite_autoindex_scheduled_payment_1` | UNIQUE (credit_purchase_id, seq) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_scheduled_payment_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `scheduled_payment_share` · Parte de un pago programado

- **Entidad del plan (§3.4.1):** pago_programado (partes por persona)
- **Descripción:** Parte de cada persona en un pago programado. La de Karen es informativa para el flujo de Axel porque ella paga directo al acreedor.
- **Usado por:** S3 y S2 (efectivo del titular contra parte informativa), S7 (estado de cuenta de Karen, alerta de pago).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `scheduled_payment_id` | INTEGER | no | Pago programado | FK → `scheduled_payment.id` al borrar: CASCADE |
| `person_id` | INTEGER | no | Persona | FK → `person.id` al borrar: RESTRICT |
| `amount_cents` | INTEGER | no | Parte de la persona (Axel 59882, Karen 37493). La suma de partes es el monto del pago | centavos |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Una parte por persona y pago.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `scheduled_payment_id` → `scheduled_payment.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_scheduled_payment_share_person` | `scheduled_payment_share (person_id)` | Pagos por persona (estado de cuenta de Karen). |
| `sqlite_autoindex_scheduled_payment_share_1` | UNIQUE (scheduled_payment_id, person_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_scheduled_payment_share_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `credit_statement` · Estado de cuenta de tarjeta

- **Entidad del plan (§3.4.1):** credito (fechas de corte y pago, pago mínimo)
- **Descripción:** Corte de una tarjeta bancaria: saldo, pago para no generar intereses, pago mínimo, intereses y comisiones. Se captura del estado de cuenta del banco o lo calcula el motor.
- **Usado por:** S3 (tarjeta: corte, pago, mínimo, intereses en cero si se paga el total a tiempo). Pantalla Créditos.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_id` | INTEGER | no | Tarjeta | FK → `credit.id` al borrar: CASCADE |
| `statement_date` | TEXT | no | Fecha de corte | fecha válida |
| `due_date` | TEXT | no | Fecha límite de pago | no antes del corte |
| `balance_cents` | INTEGER | no | Saldo al corte | centavos |
| `no_interest_payment_cents` | INTEGER | sí | Pago para no generar intereses | centavos |
| `minimum_payment_cents` | INTEGER | sí | Pago mínimo | centavos |
| `interest_cents` | INTEGER | no | Intereses cobrados en el periodo, con IVA | def: 0; centavos |
| `fees_cents` | INTEGER | no | Comisiones cobradas en el periodo | def: 0; centavos |
| `source` | TEXT | no | bank = capturado del estado de cuenta, computed = calculado por el motor | enumerado |
| `attachment_id` | INTEGER | sí | PDF del estado de cuenta | FK → `attachment.id` al borrar: SET NULL |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- La fecha límite no es anterior al corte.
- Un corte por fecha y tarjeta.

**Relaciones:**
- `attachment_id` → `attachment.id` · al borrar la madre: SET NULL (queda sin liga)
- `credit_id` → `credit.id` · al borrar la madre: CASCADE (se borra con la madre)
- La referencian: `credit_payment.credit_statement_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_credit_statement_attachment` | `credit_statement (attachment_id)` | Llave foránea hacia el PDF. |
| `ix_credit_statement_due` | `credit_statement (due_date)` | Cortes por fecha límite (alerta de pago próximo). |
| `sqlite_autoindex_credit_statement_1` | UNIQUE (credit_id, statement_date) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_credit_statement_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `credit_payment` · Pago hecho

- **Entidad del plan (§3.4.1):** pago_programado (pago real) y decisión A: Karen paga directo y repone la línea cuando paga
- **Descripción:** Pago realmente hecho a un crédito: quién pagó, cuándo y cuánto, y qué pago programado o corte cubre. Con esto se sabe si Karen ya pagó su parte.
- **Usado por:** S3 (reposición de capital, saldo vivo), S7 (estado de cuenta de Karen, alerta de pago). Pantalla Créditos.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `credit_id` | INTEGER | no | Crédito pagado | FK → `credit.id` al borrar: RESTRICT |
| `scheduled_payment_id` | INTEGER | sí | Pago programado que cubre. NULL = abono libre (tarjeta) | FK → `scheduled_payment.id` al borrar: RESTRICT |
| `credit_statement_id` | INTEGER | sí | Corte de tarjeta que paga | FK → `credit_statement.id` al borrar: RESTRICT |
| `payer_person_id` | INTEGER | no | Quién pagó (Karen paga directo su parte) | FK → `person.id` al borrar: RESTRICT |
| `paid_date` | TEXT | no | Fecha del pago | fecha válida |
| `amount_cents` | INTEGER | no | Monto pagado en centavos | mayor que 0 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `payer_person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `credit_statement_id` → `credit_statement.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `scheduled_payment_id` → `scheduled_payment.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `credit_id` → `credit.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_credit_payment_credit_date` | `credit_payment (credit_id, paid_date)` | Pagos de un crédito por fecha. |
| `ix_credit_payment_payer_date` | `credit_payment (payer_person_id, paid_date)` | Pagos hechos por una persona (Karen) por fecha. |
| `ix_credit_payment_scheduled` | `credit_payment (scheduled_payment_id)` | Pagos que cubren un pago programado (¿ya se pagó?). |
| `ix_credit_payment_statement` | `credit_payment (credit_statement_id)` | Llave foránea hacia el corte de tarjeta. |

**Disparadores:**
- `trg_credit_payment_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## B (cont.). Ahorro

### `savings_goal` · Meta de ahorro

- **Entidad del plan (§3.4.1):** meta_ahorro
- **Descripción:** Fondo de ahorro o meta con monto y fecha, con su aportación mensual planificada. El avance se calcula con sus movimientos.
- **Usado por:** Ajustes «Salud y metas de ahorro». S2 (aportación en el plan del mes, avance). Pantalla Inicio (tarjeta de ahorro).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `name` | TEXT | no | Nombre (Fondo de emergencia) | único; único |
| `kind` | TEXT | no | fund = fondo de ahorro, goal = meta con monto y fecha | enumerado |
| `target_cents` | INTEGER | sí | Monto objetivo en centavos | obligatorio en goal |
| `target_date` | TEXT | sí | Fecha objetivo | obligatoria en goal, posterior al inicio |
| `planned_monthly_cents` | INTEGER | sí | Aportación mensual planificada en centavos | centavos |
| `start_date` | TEXT | no | Inicio del fondo o meta | fecha válida |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Una meta tiene monto y fecha.
- La fecha objetivo es posterior al inicio.

**Relaciones:**
- Sin llaves foráneas salientes.
- La referencian: `savings_movement.savings_goal_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_savings_goal_1` | UNIQUE (name) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_savings_goal_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `savings_movement` · Movimiento de ahorro

- **Entidad del plan (§3.4.1):** movimiento_ahorro
- **Descripción:** Aportación o retiro de un fondo o meta.
- **Usado por:** S2 (avance de metas), Pantalla Inicio.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `savings_goal_id` | INTEGER | no | Fondo o meta | FK → `savings_goal.id` al borrar: RESTRICT |
| `movement_date` | TEXT | no | Fecha | fecha válida |
| `kind` | TEXT | no | deposit = aportación, withdrawal = retiro | enumerado |
| `amount_cents` | INTEGER | no | Monto en centavos | mayor que 0 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `savings_goal_id` → `savings_goal.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_savings_movement_date` | `savings_movement (movement_date)` | Movimientos de un periodo. |
| `ix_savings_movement_goal_date` | `savings_movement (savings_goal_id, movement_date)` | Movimientos de una meta por fecha (avance). |

**Disparadores:**
- `trg_savings_movement_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## D. Comida, salud e inventario

### `meal_slot` · Tiempo de comida

- **Entidad del plan (§3.4.1):** tiempo_comida
- **Descripción:** Los cuatro tiempos del día: desayuno, comida, cena y colación (cuarto tiempo opcional con costo y calorías). Datos de referencia que crea esta migración.
- **Usado por:** S4 (plan semanal por tiempo), Pantalla Comidas, Ajustes «Recetas».

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `code` | TEXT | no | breakfast = desayuno, lunch = comida, dinner = cena, snack = colación | único; enumerado |
| `name` | TEXT | no | Nombre en pantalla | no vacío |
| `sort_order` | INTEGER | no | Orden en el día | único; único |
| `is_optional` | INTEGER | no | 1 = tiempo opcional (colación) | def: 0; booleano |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- Sin llaves foráneas salientes.
- La referencian: `meal_plan_entry.meal_slot_id`, `recipe_meal_slot.meal_slot_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_meal_slot_1` | UNIQUE (code) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |
| `sqlite_autoindex_meal_slot_2` | UNIQUE (sort_order) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_meal_slot_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `recipe` · Receta

- **Entidad del plan (§3.4.1):** receta
- **Descripción:** Receta con porciones base (2), tiempos y origen: recetario inicial escrito por Claude o recetas de Axel y Ximena. Editable.
- **Usado por:** Ajustes «Recetas». S4 (escalado, costo y kcal por receta). Pantalla Comidas.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `name` | TEXT | no | Nombre | único; único |
| `description` | TEXT | sí | Descripción corta |  |
| `instructions` | TEXT | sí | Pasos de preparación |  |
| `base_servings` | INTEGER | no | Porciones para las que están escritas las cantidades (base 2) | def: 2; 1 a 50 |
| `prep_minutes` | INTEGER | sí | Minutos de preparación | mayor o igual a 0 |
| `cook_minutes` | INTEGER | sí | Minutos de cocción | mayor o igual a 0 |
| `origin` | TEXT | no | claude_seed = recetario inicial escrito por Claude, household = de Axel o Ximena | enumerado |
| `author_person_id` | INTEGER | sí | Quién la dio de alta (en household) | FK → `person.id` al borrar: SET NULL |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `author_person_id` → `person.id` · al borrar la madre: SET NULL (queda sin liga)
- La referencian: `inventory_item.recipe_id`, `meal_plan_entry.recipe_id`, `recipe_ingredient.recipe_id`, `recipe_meal_slot.recipe_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_recipe_author` | `recipe (author_person_id)` | Llave foránea hacia el autor. |
| `sqlite_autoindex_recipe_1` | UNIQUE (name) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_recipe_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `recipe_meal_slot` · Tiempo apto de una receta

- **Entidad del plan (§3.4.1):** receta (detalle: tiempos aptos)
- **Descripción:** Tiempos de comida para los que sirve una receta (desayuno, comida, cena o colación). Una receta puede servir para varios.
- **Usado por:** Pantalla Comidas (selector de recetas de una casilla, N-122). Ajustes «Recetas» (lista con filtro por tiempo N-135, crear receta N-136). S4 (plan semanal).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `recipe_id` | INTEGER | no | Receta | FK → `recipe.id` al borrar: CASCADE |
| `meal_slot_id` | INTEGER | no | Tiempo para el que sirve | FK → `meal_slot.id` al borrar: RESTRICT |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Un tiempo aparece una sola vez por receta.

**Relaciones:**
- `meal_slot_id` → `meal_slot.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `recipe_id` → `recipe.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_recipe_meal_slot_slot` | `recipe_meal_slot (meal_slot_id, recipe_id)` | Recetas aptas para un tiempo (selector de la casilla). |
| `sqlite_autoindex_recipe_meal_slot_1` | UNIQUE (recipe_id, meal_slot_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_recipe_meal_slot_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `recipe_ingredient` · Ingrediente de receta

- **Entidad del plan (§3.4.1):** ingrediente_receta
- **Descripción:** Ingrediente exacto de una receta: producto, presentación comprable sugerida y cantidad para las porciones base en milésimas de la unidad base.
- **Usado por:** S4 (escalado exacto, lista de compras, kcal por receta), S5 (costo por ingrediente).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `recipe_id` | INTEGER | no | Receta | FK → `recipe.id` al borrar: CASCADE |
| `product_id` | INTEGER | no | Ingrediente (producto genérico) | FK → `product.id` al borrar: RESTRICT |
| `presentation_id` | INTEGER | sí | Presentación comprable sugerida | FK → `presentation.id` al borrar: SET NULL |
| `quantity_milli` | INTEGER | no | Cantidad para base_servings en milésimas de la unidad base (250 g = 250000, media pieza = 500) | mayor que 0 |
| `display_text` | TEXT | sí | Cantidad como la escribe la receta (1/2 taza) |  |
| `is_optional` | INTEGER | no | 1 = ingrediente opcional | def: 0; booleano |
| `sort_order` | INTEGER | no | Orden en la receta | def: 0 |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Un producto aparece una vez por receta.

**Relaciones:**
- `presentation_id` → `presentation.id` · al borrar la madre: SET NULL (queda sin liga)
- `product_id` → `product.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `recipe_id` → `recipe.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_recipe_ingredient_presentation` | `recipe_ingredient (presentation_id)` | Llave foránea hacia la presentación. |
| `ix_recipe_ingredient_product` | `recipe_ingredient (product_id)` | Recetas que usan un producto. |
| `sqlite_autoindex_recipe_ingredient_1` | UNIQUE (recipe_id, product_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_recipe_ingredient_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `meal_plan_entry` · Entrada del plan de comidas

- **Entidad del plan (§3.4.1):** plan_comida
- **Descripción:** Un renglón del plan: día × tiempo × persona × receta × porciones, o «fuera» (sin receta ni calorías). Solo Axel y Ximena.
- **Usado por:** S4 (plan semanal, lista de compras, costo y kcal por día), S7 (alerta de plan sin armar). Pantalla Comidas.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `plan_date` | TEXT | no | Día del plan | fecha válida |
| `meal_slot_id` | INTEGER | no | Tiempo de comida | FK → `meal_slot.id` al borrar: RESTRICT |
| `person_id` | INTEGER | no | Axel o Ximena | FK → `person.id` al borrar: RESTRICT; nunca una persona externa (disparador) |
| `recipe_id` | INTEGER | sí | Receta. NULL si come fuera | FK → `recipe.id` al borrar: RESTRICT |
| `servings_milli` | INTEGER | sí | Porciones que come la persona ×1000 (1 porción = 1000) | NULL si come fuera |
| `is_eating_out` | INTEGER | no | 1 = come fuera: sin receta ni calorías, el gasto va a la categoría «comida fuera» | def: 0; booleano |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Con receta lleva porciones. Fuera no lleva receta ni porciones.

**Relaciones:**
- `recipe_id` → `recipe.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `meal_slot_id` → `meal_slot.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `inventory_movement.meal_plan_entry_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_meal_plan_entry_date` | `meal_plan_entry (plan_date)` | Plan por fecha (semana, lista de compras). |
| `ix_meal_plan_entry_person_date` | `meal_plan_entry (person_id, plan_date)` | Plan de una persona por fecha (kcal por día). |
| `ix_meal_plan_entry_recipe` | `meal_plan_entry (recipe_id)` | Días en que se planeó una receta. |
| `ix_meal_plan_entry_slot` | `meal_plan_entry (meal_slot_id)` | Llave foránea hacia el tiempo de comida. |
| `ux_meal_plan_entry_slot` | UNIQUE `meal_plan_entry (plan_date, meal_slot_id, person_id, coalesce(recipe_id, 0))` | Una receta una sola vez por día, tiempo y persona (y un solo «fuera»). |

**Disparadores:**
- `trg_meal_plan_entry_not_external_ins`: rechaza con «meal_plan_entry: el plan de comidas es solo para Axel y Ximena (decisión A)».
- `trg_meal_plan_entry_not_external_upd`: rechaza con «meal_plan_entry: el plan de comidas es solo para Axel y Ximena (decisión A)».
- `trg_meal_plan_entry_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `inventory_item` · Renglón de inventario

- **Entidad del plan (§3.4.1):** inventario_item
- **Descripción:** Inventario híbrido: existencia por unidad (básicos, no perecederos y sobrantes de recetas) o ritmo de consumo (el resto).
- **Usado por:** Ajustes «Inventario». S4 (lista = ingredientes − existencia, sobrantes). Pantalla Compras y tiendas.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `product_id` | INTEGER | sí | Producto en existencia | FK → `product.id` al borrar: RESTRICT; exactamente uno de product_id o recipe_id |
| `recipe_id` | INTEGER | sí | Sobrante de una receta (en porciones) | FK → `recipe.id` al borrar: RESTRICT |
| `tracking_mode` | TEXT | no | stock = existencia por unidad, consumption_rate = ritmo de consumo | enumerado |
| `quantity_milli` | INTEGER | no | Existencia en milésimas de la unidad base (en sobrantes, porciones ×1000): la vigente en stock y la del último conteo en ritmo de consumo. Cada cambio deja su inventory_movement | def: 0; mayor o igual a 0 |
| `counted_date` | TEXT | sí | Fecha del último conteo, base de la estimación por ritmo | obligatoria en consumption_rate |
| `rate_quantity_milli` | INTEGER | sí | Ritmo: cantidad que se consume (milésimas de la unidad base) | solo en consumption_rate |
| `rate_period_days` | INTEGER | sí | Ritmo: en cuántos días se consume esa cantidad | solo en consumption_rate |
| `reorder_level_milli` | INTEGER | sí | Nivel que pone el producto en la lista de compras. NULL = sin aviso |  |
| `expires_date` | TEXT | sí | Caducidad (sobrantes y perecederos) | fecha válida |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Cada renglón es de un producto o de un sobrante de receta, nunca de ambos.
- El ritmo de consumo lleva cantidad, días y fecha de conteo. La existencia por unidad no lleva ritmo.
- Los sobrantes se llevan por unidad (porciones).

**Relaciones:**
- `recipe_id` → `recipe.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `product_id` → `product.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `inventory_movement.inventory_item_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_inventory_item_recipe` | `inventory_item (recipe_id)` | Sobrantes de una receta. |
| `ux_inventory_item_product` | UNIQUE `inventory_item (product_id) WHERE product_id IS NOT NULL` | Un solo renglón de inventario por producto. |

**Disparadores:**
- `trg_inventory_item_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `inventory_movement` · Movimiento de inventario

- **Entidad del plan (§3.4.1):** inventario_item (detalle: movimientos)
- **Descripción:** Cada entrada o salida de un renglón de inventario: compra, sobrante de receta, consumo del plan o ajuste a mano con su motivo (conteo, consumo o merma). La API lo escribe en la misma transacción que cambia la existencia.
- **Usado por:** Pantalla Compras y tiendas › Inventario (movimientos N-181, ajustar existencia N-183). S4 (existencia, consumo del plan y sobrantes).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `inventory_item_id` | INTEGER | no | Renglón de inventario | FK → `inventory_item.id` al borrar: CASCADE |
| `movement_date` | TEXT | no | Fecha del movimiento | fecha válida |
| `kind` | TEXT | no | purchase = entrada por compra, recipe_leftover = sobrante de receta, plan_consumption = consumo del plan, adjustment = ajuste a mano (N-181) | enumerado |
| `reason` | TEXT | sí | Motivo del ajuste: count = conteo, consumption = consumo, waste = merma (N-183) | obligatorio solo en adjustment |
| `quantity_delta_milli` | INTEGER | no | Diferencia en milésimas de la unidad base (en sobrantes, porciones ×1000): positiva entra, negativa sale. Un ajuste por «cantidad nueva» guarda la diferencia | entero con signo |
| `expense_item_id` | INTEGER | sí | Renglón del ticket que le dio origen (compra) | FK → `expense_item.id` al borrar: SET NULL |
| `meal_plan_entry_id` | INTEGER | sí | Entrada del plan que le dio origen (consumo o sobrante) | FK → `meal_plan_entry.id` al borrar: SET NULL |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Solo el ajuste lleva motivo, y siempre lo lleva.
- El renglón de ticket solo acompaña a una compra.
- La entrada del plan solo acompaña a consumos y sobrantes.

**Relaciones:**
- `meal_plan_entry_id` → `meal_plan_entry.id` · al borrar la madre: SET NULL (queda sin liga)
- `expense_item_id` → `expense_item.id` · al borrar la madre: SET NULL (queda sin liga)
- `inventory_item_id` → `inventory_item.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_inventory_movement_expense_item` | `inventory_movement (expense_item_id)` | Llave foránea hacia el renglón del ticket. |
| `ix_inventory_movement_item_date` | `inventory_movement (inventory_item_id, movement_date)` | Movimientos de un renglón por fecha (N-181). |
| `ix_inventory_movement_meal_plan` | `inventory_movement (meal_plan_entry_id)` | Llave foránea hacia la entrada del plan. |

**Disparadores:**
- `trg_inventory_movement_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `shopping_list` · Lista de compras

- **Entidad del plan (§3.4.1):** (complementaria) lista de compras de la semana con marcas de comprado y artículos a mano (U1 N-150 a N-153 y N-165), compartida en el servidor entre los dos teléfonos (K1)
- **Descripción:** Lista de compras de una semana del plan de comidas. Guarda lo que no se puede calcular: los artículos agregados a mano, la tienda o cantidad elegida, las marcas de comprado y cuándo se generó. Lo calculado (necesidad del plan − existencia + reposición, precios y frescura) lo dan S4 y S5. Se comparte en el servidor entre los dos teléfonos (K1).
- **Usado por:** Pantalla Compras y tiendas › Lista de compras (N-150 a N-153) y comparativa (N-165 «Usar esta tienda»). Comidas (N-127 «Generar lista de compras»). Gastos (N-073 «Terminar compra»). S4 (lista neteada), S5 (tienda recomendada).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `start_date` | TEXT | no | Lunes que abre la semana del plan (K8: semana de lunes a domingo) | único; fecha válida y en lunes (lo exigen trg_shopping_list_monday_ins y trg_shopping_list_monday_upd, de la 00002), una lista por semana |
| `end_date` | TEXT | no | Último día de la semana, inclusive | fecha válida, start_date + 6 días |
| `generated_at` | TEXT | sí | Última vez que S4 la generó o actualizó desde el plan (N-127), UTC. NULL = todavía no se genera |  |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- La lista cubre exactamente una semana de 7 días.

**Relaciones:**
- Sin llaves foráneas salientes.
- La referencian: `shopping_list_item.shopping_list_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_shopping_list_1` | UNIQUE (start_date) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_shopping_list_monday_ins`: rechaza con «shopping_list: la semana de la lista empieza en lunes (K8)».
- `trg_shopping_list_monday_upd`: rechaza con «shopping_list: la semana de la lista empieza en lunes (K8)».
- `trg_shopping_list_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `shopping_list_item` · Artículo de la lista de compras

- **Entidad del plan (§3.4.1):** (complementaria) artículos de la lista de compras, calculados por S4 o agregados a mano, con tienda, cantidad y marca de comprado, compartidos en el servidor (K1)
- **Descripción:** Un artículo de la lista: un producto que calculó S4 al generar la lista o un artículo agregado a mano (producto del catálogo o texto libre). Guarda la tienda o mercado elegido, la cantidad elegida y si ya se compró. Se comparte en el servidor entre los dos teléfonos (K1).
- **Usado por:** Pantalla Compras y tiendas › Lista de compras (marcar comprado N-151, agregar artículo N-152, cambiar tienda o cantidad N-153) y comparativa (N-165). Gastos (prellenado N-073). S4 (lista neteada), S5 (mejor tienda y canasta, N-024).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `shopping_list_id` | INTEGER | no | Lista a la que pertenece | FK → `shopping_list.id` al borrar: CASCADE |
| `origin` | TEXT | no | calculated = lo calculó S4 (recetas del plan o reposición por ritmo), manual = lo agregó una persona (N-152) | enumerado |
| `product_id` | INTEGER | sí | Producto del catálogo | FK → `product.id` al borrar: RESTRICT; exactamente uno de product_id o free_text |
| `free_text` | TEXT | sí | Artículo escrito a mano que no está en el catálogo | solo en manual; si viene, no vacío |
| `free_text_unit` | TEXT | sí | Unidad de la cantidad de un artículo en texto libre (un producto usa su product.base_unit) | g, ml o piece; obligatoria con free_text |
| `presentation_id` | INTEGER | sí | Presentación elegida para comprar. NULL = la que recomienda S5 | FK → `presentation.id` al borrar: SET NULL |
| `needed_milli` | INTEGER | sí | Necesidad que calculó S4 en la última generación (plan − existencia + reposición), en milésimas de la unidad base | solo en calculated |
| `quantity_milli` | INTEGER | sí | Cantidad a comprar que fijó una persona (N-152, N-153), en milésimas de la unidad base. NULL en calculated = comprar la necesidad | obligatoria en manual |
| `store_id` | INTEGER | sí | Tienda o mercado elegido (N-153, N-165). NULL = el que recomienda S5 | FK → `store.id` al borrar: SET NULL |
| `is_bought` | INTEGER | no | 1 = comprado (N-151). No crea gasto ni precio hasta confirmar el ticket | def: 0; booleano |
| `bought_at` | TEXT | sí | Cuándo se marcó como comprado (UTC) | solo con is_bought = 1 |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Cada artículo es un producto del catálogo o un texto libre, nunca ambos.
- El texto libre lleva su unidad; un producto usa su unidad base.
- Un texto libre no tiene presentación del catálogo.
- Lo calculado es un producto con su necesidad; lo manual lleva cantidad y no tiene necesidad calculada.
- Solo un artículo comprado lleva la hora de compra.

**Relaciones:**
- `store_id` → `store.id` · al borrar la madre: SET NULL (queda sin liga)
- `presentation_id` → `presentation.id` · al borrar la madre: SET NULL (queda sin liga)
- `product_id` → `product.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `shopping_list_id` → `shopping_list.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_shopping_list_item_list` | `shopping_list_item (shopping_list_id, is_bought)` | Artículos de una lista, comprados o no. |
| `ix_shopping_list_item_presentation` | `shopping_list_item (presentation_id)` | Llave foránea hacia la presentación. |
| `ix_shopping_list_item_product` | `shopping_list_item (product_id)` | Llave foránea hacia el producto. |
| `ix_shopping_list_item_store` | `shopping_list_item (store_id)` | Artículos asignados a una tienda (lista agrupada por tienda, N-165). |
| `ux_shopping_list_item_calculated` | UNIQUE `shopping_list_item (shopping_list_id, product_id) WHERE origin = 'calculated'` | Un artículo calculado por producto en cada lista (la regeneración lo actualiza en su lugar). |

**Disparadores:**
- `trg_shopping_list_item_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `health_profile` · Perfil de salud

- **Entidad del plan (§3.4.1):** perfil_salud
- **Descripción:** Perfil de salud por persona (solo Axel y Ximena): estatura, nacimiento, sexo, actividad, meta de calorías y restricciones (vacías por decisión). El peso vive en weight_entry.
- **Usado por:** Ajustes «Salud y metas de ahorro». S4 (IMC, kcal por día contra la meta). Pantalla Salud.
- **Usado por las columnas nuevas de 1.2:** Pantalla Salud y Ajustes «Salud y metas de ahorro». S4 (meta sugerida con Mifflin-St Jeor, actividad y objetivo; regla RN-17). U4.
- **Usado por las columnas nuevas de 1.3:** Pantalla Salud y Ajustes «Salud y metas de ahorro» (ajuste por persona; el de la casa en app_setting health.goal_adjust_*). S4 (meta sugerida, regla RN-17). C1 v1.2. U4.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `person_id` | INTEGER | no | Persona | FK → `person.id` al borrar: RESTRICT; único; una por persona, nunca externa (disparador) |
| `height_mm` | INTEGER | sí | Estatura en milímetros (1.75 m = 1750) | 500 a 2500 |
| `birth_date` | TEXT | sí | Fecha de nacimiento. La edad se calcula | fecha válida |
| `sex` | TEXT | sí | Sexo para fórmulas de requerimiento | enumerado |
| `activity_level` | TEXT | sí | Nivel de actividad de la escala de 5 niveles que multiplica a Mifflin-St Jeor: sedentary 1.2, light 1.375, moderate 1.55, active 1.725, very_active 1.9 (K23) | enumerado |
| `daily_kcal_target` | INTEGER | sí | Meta de kilocalorías por día | 500 a 10000 |
| `dietary_restrictions` | TEXT | sí | Restricciones alimentarias (el campo existe y hoy va vacío) |  |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `goal` | TEXT | sí | Nueva en 1.2 (K23). Objetivo que ajusta la meta sugerida de kilocalorías: lose = bajar de peso, maintain = mantener, gain = subir. NULL = todavía no lo elige | enumerado |
| `goal_adjust_lose_kcal` | INTEGER | sí | Nueva en 1.3 (K36). Cuántas kcal se restan de la meta sugerida cuando el objetivo (goal) es lose. Es una magnitud: siempre se resta. NULL = usa el valor de la casa (app_setting health.goal_adjust_lose_kcal, 0 al sembrar) | entero mayor o igual a 0, en kcal; NULL = valor de la casa |
| `goal_adjust_gain_kcal` | INTEGER | sí | Nueva en 1.3 (K36). Cuántas kcal se suman a la meta sugerida cuando el objetivo (goal) es gain. Es una magnitud: siempre se suma. maintain no ajusta. NULL = usa el valor de la casa (app_setting health.goal_adjust_gain_kcal, 0 al sembrar) | entero mayor o igual a 0, en kcal; NULL = valor de la casa |

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_health_profile_1` | UNIQUE (person_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_health_profile_not_external_ins`: rechaza con «health_profile: salud solo para Axel y Ximena (decisión A)».
- `trg_health_profile_not_external_upd`: rechaza con «health_profile: salud solo para Axel y Ximena (decisión A)».
- `trg_health_profile_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `weight_entry` · Peso histórico

- **Entidad del plan (§3.4.1):** peso_historico
- **Descripción:** Historial de peso por persona. El más reciente alimenta el IMC.
- **Usado por:** S4 (IMC), Pantalla Salud (tendencia).
- **Usado por las columnas nuevas de 1.2:** C1 v1.1 (pesos capturados sin conexión). S6 (busca la clave antes de insertar). PWA en WF-3 (cola del teléfono). Regla RN-24.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `person_id` | INTEGER | no | Persona | FK → `person.id` al borrar: RESTRICT; nunca externa (disparador) |
| `measured_date` | TEXT | no | Fecha de la medición | fecha válida |
| `weight_g` | INTEGER | no | Peso en gramos (72.5 kg = 72500) | 20 a 400 kg |
| `notes` | TEXT | sí | Notas libres. |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |
| `idempotency_key` | TEXT | sí | Nueva en 1.2 (K9). Clave de la captura en la cola del teléfono. Si llega otra vez, la API devuelve este peso y no crea otro. No se purga | 16 a 128 caracteres: letras sin acento, dígitos, guion y guion bajo; única (ux_weight_entry_idempotency_key); NULL = sin clave |
| `captured_at` | TEXT | sí | Nueva en 1.2 (K9). Instante en que la persona capturó el peso en el teléfono (UTC). NULL = captura en línea: la hora de captura es created_at | instante ISO 8601 con «Z» |

**Reglas de tabla (CHECK / UNIQUE):**
- Una medición por persona y día.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_weight_entry_1` | UNIQUE (person_id, measured_date) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |
| `ux_weight_entry_idempotency_key` | UNIQUE `weight_entry (idempotency_key)` | Una captura de la cola no se duplica: un solo peso por clave. |

**Disparadores:**
- `trg_weight_entry_not_external_ins`: rechaza con «weight_entry: salud solo para Axel y Ximena (decisión A)».
- `trg_weight_entry_not_external_upd`: rechaza con «weight_entry: salud solo para Axel y Ximena (decisión A)».
- `trg_weight_entry_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `nutrition_fact` · Dato nutricional

- **Entidad del plan (§3.4.1):** nutricion_ingrediente
- **Descripción:** Kilocalorías por 100 g de un producto (o de una presentación específica) con fuente citada. Macronutrientes previstos, se llenan después.
- **Usado por:** S4 (kcal por receta y por día), carga local de USDA FoodData Central y Open Food Facts.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `product_id` | INTEGER | no | Producto | FK → `product.id` al borrar: CASCADE |
| `presentation_id` | INTEGER | sí | Presentación específica (etiqueta, Open Food Facts). NULL = valor genérico del producto | FK → `presentation.id` al borrar: CASCADE |
| `kcal_milli_per_100g` | INTEGER | no | Milésimas de kilocaloría por 100 g (884 kcal = 884000) | mayor o igual a 0 |
| `protein_mg_per_100g` | INTEGER | sí | Proteína en mg por 100 g (previsto) | mayor o igual a 0 |
| `fat_mg_per_100g` | INTEGER | sí | Grasa en mg por 100 g (previsto) | mayor o igual a 0 |
| `carbs_mg_per_100g` | INTEGER | sí | Carbohidratos en mg por 100 g (previsto) | mayor o igual a 0 |
| `fiber_mg_per_100g` | INTEGER | sí | Fibra en mg por 100 g (previsto) | mayor o igual a 0 |
| `source` | TEXT | no | Fuente: USDA FoodData Central, Open Food Facts, tabla mexicana citada, etiqueta de tienda o captura manual | enumerado |
| `source_ref` | TEXT | sí | Cita: fdc_id, código de barras, tabla y página o URL de la etiqueta | obligatoria salvo en manual |
| `source_license` | TEXT | sí | Licencia de la fuente (CC0, ODbL) |  |
| `is_preferred` | INTEGER | no | 1 = valor que usan los cálculos | def: 0; booleano, uno por producto o presentación |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Toda fuente externa lleva su cita.

**Relaciones:**
- `presentation_id` → `presentation.id` · al borrar la madre: CASCADE (se borra con la madre)
- `product_id` → `product.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_nutrition_fact_presentation` | `nutrition_fact (presentation_id)` | Llave foránea hacia la presentación. |
| `ix_nutrition_fact_product` | `nutrition_fact (product_id, source)` | Datos de un producto. |
| `ux_nutrition_fact_preferred` | UNIQUE `nutrition_fact (product_id, coalesce(presentation_id, 0)) WHERE is_preferred = 1` | Un solo valor preferido por producto o presentación. |

**Disparadores:**
- `trg_nutrition_fact_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## F. Avisos, jobs, respaldos y reportes

### `notification` · Aviso

- **Entidad del plan (§3.4.1):** notificacion
- **Descripción:** Cada aviso o reporte configurable: resúmenes diario y semanal, reporte mensual, pago próximo, presupuesto 80/100 %, precios viejos y plan sin armar. Esta migración crea los siete con los parámetros que ya decidió Axel.
- **Usado por:** Ajustes «Notificaciones» y «Reportes y frecuencias». S7 (reportes y alertas).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `code` | TEXT | no | Clave del aviso | único; enumerado |
| `name` | TEXT | no | Nombre en Ajustes | no vacío |
| `is_enabled` | INTEGER | no | 1 = activo | def: 1; booleano |
| `params_json` | TEXT | no | Parámetros del aviso (días antes 1 y 0, umbrales 80 y 100). Desde 1.3 (00003, K44) `payment_due` ya no lleva `days_before` ni `budget_threshold` lleva `thresholds_pct` (son `alerts.payment_due_days_before`, `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `app_setting`), y desde 1.4 (00004, K47) `weekly_summary` y `meal_plan_missing` ya no llevan `weekday` ni `monthly_report` lleva `day_of_month` (son `reports.weekly_weekday`, `alerts.meal_plan_missing_weekday` y `reports.monthly_day_of_month`); la hora y la anticipación nunca estuvieron aquí (`reports.*_send_time`, `alerts.send_time` y `alerts.stale_prices_advance_days`). Con los valores iniciales solo `budget_threshold` conserva un parámetro (`scopes`); los demás quedan con `{}`. Los Down de la 00004 y de la 00003 los reponen aquí (RN-29). | def: '{}'; JSON válido |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- Sin llaves foráneas salientes.
- La referencian: `job.notification_id`, `notification_delivery.notification_id`, `notification_subscription.notification_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `sqlite_autoindex_notification_1` | UNIQUE (code) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_notification_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `notification_recipient` · Destinatario

- **Entidad del plan (§3.4.1):** destinatario
- **Descripción:** Destino de avisos de una persona: su chat de Telegram en el bot MultiX. El chat_id no es credencial. El token del bot nunca se guarda en la base.
- **Usado por:** Ajustes «Notificaciones». S7 (envío de avisos a Axel y Ximena).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `person_id` | INTEGER | no | Persona destinataria | FK → `person.id` al borrar: RESTRICT |
| `channel` | TEXT | no | Canal (bot MultiX) | def: 'telegram'; enumerado |
| `chat_id` | TEXT | no | chat_id de Telegram | dígitos y signo menos |
| `label` | TEXT | sí | Etiqueta (Celular de Ximena) |  |
| `verified_at` | TEXT | sí | Cuándo se comprobó con un mensaje de prueba (UTC) |  |
| `is_active` | INTEGER | no | 0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia. | def: 1 |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Un chat se registra una sola vez.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- La referencian: `notification_delivery.recipient_id`, `notification_subscription.recipient_id`

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_notification_recipient_person` | `notification_recipient (person_id)` | Destinos de una persona. |
| `sqlite_autoindex_notification_recipient_1` | UNIQUE (channel, chat_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_notification_recipient_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `telegram_chat_candidate` · Chat de Telegram por confirmar

- **Entidad del plan (§3.4.1):** (complementaria) chats que el bot detecta cuando alguien le escribe y que la persona confirma en Ajustes (K4)
- **Descripción:** Chat que el bot detectó cuando alguien le escribió y que espera confirmación en Ajustes «Notificaciones». Guarda solo el chat_id y el nombre que da Telegram. Al confirmarlo se liga a la persona y la API crea su notification_recipient en la misma transacción; descartarlo borra el renglón.
- **Usado por:** Ajustes «Notificaciones» (candidatos pendientes y confirmación, C1 v1.1). S7 (detección cuando la persona escribe al bot, K4). U1 P-26.09: N-299 (chats detectados) y N-316 (Confirmar chat detectado).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `chat_id` | TEXT | no | chat_id de Telegram del que llegó el mensaje. La API nunca lo devuelve completo: lo muestra como el campo calculado chat_id_masked («•••» y sus últimos 4 caracteres; K4, C1-18) | único; dígitos y signo menos |
| `telegram_name` | TEXT | no | Nombre que da Telegram para ese chat | no vacío |
| `detected_at` | TEXT | no | Cuándo lo detectó el bot (UTC) | def: ahora (UTC); instante ISO 8601 con «Z» |
| `confirmed_at` | TEXT | sí | Cuándo la persona lo confirmó en Ajustes (UTC). NULL = pendiente | instante ISO 8601 con «Z» |
| `person_id` | INTEGER | sí | Persona a quien se ligó al confirmar | FK → `person.id` al borrar: RESTRICT; obligatoria si está confirmado |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Un chat confirmado dice de quién es; uno pendiente todavía no tiene persona.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_telegram_chat_candidate_pending` | `telegram_chat_candidate (detected_at) WHERE confirmed_at IS NULL` | Candidatos pendientes de confirmar. |
| `ix_telegram_chat_candidate_person` | `telegram_chat_candidate (person_id)` | Llave foránea hacia la persona. |
| `sqlite_autoindex_telegram_chat_candidate_1` | UNIQUE (chat_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_telegram_chat_candidate_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `notification_subscription` · Suscripción a un aviso

- **Entidad del plan (§3.4.1):** destinatario (destinatarios configurables por aviso)
- **Descripción:** Qué destinatario recibe qué aviso. Se configura por aviso en Ajustes.
- **Usado por:** Ajustes «Notificaciones». S7.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `notification_id` | INTEGER | no | Aviso | FK → `notification.id` al borrar: CASCADE |
| `recipient_id` | INTEGER | no | Destinatario | FK → `notification_recipient.id` al borrar: CASCADE |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Una suscripción por aviso y destinatario.

**Relaciones:**
- `recipient_id` → `notification_recipient.id` · al borrar la madre: CASCADE (se borra con la madre)
- `notification_id` → `notification.id` · al borrar la madre: CASCADE (se borra con la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_notification_subscription_recipient` | `notification_subscription (recipient_id)` | Avisos de un destinatario. |
| `sqlite_autoindex_notification_subscription_1` | UNIQUE (notification_id, recipient_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_notification_subscription_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `notification_delivery` · Envío de aviso

- **Entidad del plan (§3.4.1):** notificacion (bitácora de envíos)
- **Descripción:** Cada envío de un aviso a un destinatario, con clave para no mandar dos veces lo mismo y el resultado HTTP.
- **Usado por:** S7 (deduplicar alertas, evidencia de envío HTTP 200). Ajustes «Notificaciones» (historial).

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `notification_id` | INTEGER | no | Aviso | FK → `notification.id` al borrar: RESTRICT |
| `recipient_id` | INTEGER | no | Destinatario | FK → `notification_recipient.id` al borrar: RESTRICT |
| `dedupe_key` | TEXT | no | Clave que evita repetir el mismo aviso (payment_due:42:2026-10-29) | no vacía |
| `status` | TEXT | no | Estado del envío | def: 'pending'; enumerado |
| `message_text` | TEXT | sí | Texto plano enviado (sin parse_mode) |  |
| `http_status` | INTEGER | sí | Código HTTP de la API de Telegram |  |
| `error` | TEXT | sí | Error resumido, sin token ni datos sensibles |  |
| `sent_at` | TEXT | sí | Momento del envío (UTC) |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- El mismo aviso no se manda dos veces al mismo destinatario.

**Relaciones:**
- `recipient_id` → `notification_recipient.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)
- `notification_id` → `notification.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_notification_delivery_notification` | `notification_delivery (notification_id, sent_at)` | Historial de un aviso. |
| `ix_notification_delivery_recipient` | `notification_delivery (recipient_id)` | Llave foránea hacia el destinatario. |
| `ix_notification_delivery_status` | `notification_delivery (status, created_at)` | Cola de envíos pendientes o fallidos. |
| `sqlite_autoindex_notification_delivery_1` | UNIQUE (dedupe_key, recipient_id) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_notification_delivery_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `job` · Job programado

- **Entidad del plan (§3.4.1):** job
- **Descripción:** Tarea programada del servicio (respaldos, resúmenes, alertas, purgas) con su horario y su último resultado.
- **Usado por:** S6 (planificador), S7 (reportes y alertas). Ajustes «Reportes y frecuencias».

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `code` | TEXT | no | Clave del job (backup_daily, daily_summary) | único; a-z, 0-9 y guion bajo |
| `name` | TEXT | no | Nombre visible | no vacío |
| `schedule` | TEXT | no | Expresión cron de 5 campos en hora de America/Mexico_City | no vacía |
| `is_enabled` | INTEGER | no | 1 = activo | def: 1; booleano |
| `notification_id` | INTEGER | sí | Aviso que manda, si manda alguno | FK → `notification.id` al borrar: SET NULL |
| `last_run_at` | TEXT | sí | Inicio de la última corrida (UTC) |  |
| `last_status` | TEXT | sí | Resultado de la última corrida | enumerado |
| `last_error` | TEXT | sí | Error resumido de la última corrida |  |
| `next_run_at` | TEXT | sí | Próxima corrida calculada (UTC) |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Relaciones:**
- `notification_id` → `notification.id` · al borrar la madre: SET NULL (queda sin liga)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_job_next_run` | `job (is_enabled, next_run_at)` | Siguiente job por correr. |
| `ix_job_notification` | `job (notification_id)` | Llave foránea hacia el aviso. |
| `sqlite_autoindex_job_1` | UNIQUE (code) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_job_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `backup` · Respaldo

- **Entidad del plan (§3.4.1):** respaldo
- **Descripción:** Respaldo local de la base (VACUUM INTO) con huella, resultado de integrity_check y prueba de restauración.
- **Usado por:** S6 (respaldo diario, restauración probada), S8 (operación). Ajustes «Datos y respaldos».

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `kind` | TEXT | no | Diario programado, manual o previo a una migración | enumerado |
| `status` | TEXT | no | Estado | def: 'running'; enumerado |
| `started_at` | TEXT | no | Inicio (UTC) | def: ahora (UTC) |
| `finished_at` | TEXT | sí | Fin (UTC) |  |
| `file_path` | TEXT | sí | Ruta del archivo en la carpeta o disco de respaldos | obligatoria si ok |
| `size_bytes` | INTEGER | sí | Tamaño en bytes | mayor o igual a 0 |
| `sha256` | TEXT | sí | Huella del archivo | 64 caracteres, obligatoria si ok |
| `integrity_check` | TEXT | sí | Resultado de PRAGMA integrity_check sobre la copia (ok o el error) | ok si status = ok |
| `restore_tested_at` | TEXT | sí | Última restauración de prueba (UTC) |  |
| `restore_test_result` | TEXT | sí | Resultado de esa restauración |  |
| `pruned_at` | TEXT | sí | Cuándo se borró el archivo por la política de retención (UTC) |  |
| `error` | TEXT | sí | Error resumido |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- Un respaldo correcto tiene archivo, huella e integrity_check = ok.

**Relaciones:**
- Sin llaves foráneas salientes.

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_backup_started` | `backup (started_at)` | Respaldos por fecha. |
| `ix_backup_status` | `backup (status, started_at)` | Último respaldo correcto o fallido. |

**Disparadores:**
- `trg_backup_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

### `generated_report` · Reporte generado

- **Entidad del plan (§3.4.1):** (complementaria) decisión F: el reporte mensual PDF se envía y se guarda; estado de cuenta de Karen en PDF bajo demanda
- **Descripción:** PDF generado por el sistema: reporte mensual LaTeX o estado de cuenta de la participante externa (Karen), con su archivo y huella.
- **Usado por:** S7 (reporte mensual día 1, estado de cuenta de Karen). Pantalla Reportes.

| nombre_codigo | tipo | nulo | descripción | regla |
|---|---|---|---|---|
| `id` | INTEGER | no | Identificador interno (alias de rowid). | PK |
| `kind` | TEXT | no | monthly_pdf = reporte mensual, external_statement_pdf = estado de cuenta de la participante externa | enumerado |
| `period_start_date` | TEXT | no | Inicio del periodo que cubre | fecha válida |
| `period_end_date` | TEXT | no | Fin del periodo que cubre, inclusive | no antes del inicio |
| `person_id` | INTEGER | sí | Persona del estado de cuenta (Karen) | FK → `person.id` al borrar: RESTRICT; solo en external_statement_pdf |
| `file_path` | TEXT | no | Ruta del PDF guardado | único; única |
| `sha256` | TEXT | no | Huella del PDF | 64 caracteres |
| `size_bytes` | INTEGER | no | Tamaño en bytes | mayor o igual a 0 |
| `generated_at` | TEXT | no | Momento de generación (UTC) | def: ahora (UTC) |
| `sent_at` | TEXT | sí | Envío por Telegram (reporte mensual) (UTC) |  |
| `created_at` | TEXT | no | Alta del registro, instante UTC ISO 8601 con «Z». | def: ahora (UTC); ISO 8601 UTC |
| `updated_at` | TEXT | no | Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`. | def: ahora (UTC); ISO 8601 UTC |

**Reglas de tabla (CHECK / UNIQUE):**
- El periodo no termina antes de empezar.
- Solo el estado de cuenta lleva persona.

**Relaciones:**
- `person_id` → `person.id` · al borrar la madre: RESTRICT (no se puede borrar la madre)

**Índices:**

| índice | definición | propósito |
|---|---|---|
| `ix_generated_report_kind_period` | `generated_report (kind, period_start_date)` | Reportes por tipo y periodo. |
| `ix_generated_report_person` | `generated_report (person_id)` | Estados de cuenta de una persona. |
| `sqlite_autoindex_generated_report_1` | UNIQUE (file_path) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |

**Disparadores:**
- `trg_generated_report_updated_at`: pone updated_at en el instante actual si el UPDATE no lo trae.

## Datos de referencia que crea la migración
- `meal_slot`: 4 renglones fijos: 1 desayuno, 2 comida, 3 cena, 4 colación (opcional).
- `notification`: los 7 avisos de la decisión F, activos y solo con parámetros ya decididos: `weekly_summary` y `meal_plan_missing` el domingo (`weekday` 7), `monthly_report` el día 1, `payment_due` el día anterior y el mismo día (`days_before` [1, 0]), `budget_threshold` al 80 y 100 % por categoría y total. `daily_summary` («noche») y `stale_prices` («antes del día de despensa») quedan con `{}` porque la hora y los días exactos no están decididos.
- `app_setting`: `price.stale_after_days` = 30, `auth.max_failed_attempts` = 5, `auth.lockout_window_minutes` = 15, `optimizer.weights` (4 pesos iguales), `ai.local_model_enabled` = false, `ai.claude_api_enabled` = false, `price.assisted_reading_enabled` = false.
- `app_setting` sembrado por la migración 00002 (puerta 1): `projection.repeat_cadence_enabled` = false (K21), `session.max_age_days` = 90 (K11), `idempotency.retention_hours` = 24 (K6).
- `app_setting` sembrado por la migración 00003 (puerta 1b, K44): `auth.password_min_length` = 12 (K6), `price.zero_weight_after_days` = 90 (K22, K40, Q2-77 (K50)), `price.profeco_import_max_gb` = 1 (K6), `comparator.timeout_seconds` = 5 (K43, Q2-74 (K50)), `attachment.max_size_mb` = 10 (K43), `health.goal_adjust_lose_kcal` = 0 (K36), `health.goal_adjust_gain_kcal` = 0 (K36), `alerts.budget_warning_pct` = 80 (Decisión F, K44), `alerts.budget_limit_pct` = 100 (Decisión F, K44), `alerts.payment_due_days_before` = [1, 0] (Decisión F, K44).
- `notification` en la migración 00003 (K44): `budget_threshold` queda con `{"scopes": ["category", "total"]}` y `payment_due` queda con `{}`; los días antes y los umbrales viven en `alerts.payment_due_days_before`, `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `app_setting` (el Down de la 00003 los repone).
- `app_setting` sembrado por la migración 00004 (puerta 2, K47 y K50): `reports.daily_send_time` = "21:30" (Q2-72 (K50), K31 punto 1, lugar K47), `reports.weekly_weekday` = 7 (Decisión F, lugar K47), `reports.weekly_send_time` = "19:00" (Q2-72 (K50), K31 punto 1, lugar K47), `reports.monthly_day_of_month` = 1 (Decisión F, Q2-66 B, lugar K47), `reports.monthly_send_time` = "08:00" (Q2-72 (K50), K31 punto 1, lugar K47), `alerts.send_time` = "09:00" (Q2-73 (K50), K31 punto 2, K43, lugar K47), `alerts.stale_prices_advance_days` = 1 (Q2-73 (K50), K31 punto 2, K43, lugar K47), `alerts.meal_plan_missing_weekday` = 7 (Decisión F, lugar K47), `optimizer.date_window_days` = 7 (Q2-75 (K50)), `ocr.timeout_seconds` = 30 (Q2-76 (K50); sección: propuesta S1).
- `notification` en la migración 00004 (K47): `weekly_summary` queda con `{}` y `monthly_report` queda con `{}` y `meal_plan_missing` queda con `{}`; el día de la semana y el día del mes viven en `reports.weekly_weekday`, `alerts.meal_plan_missing_weekday` y `reports.monthly_day_of_month` de `app_setting` (el Down de la 00004 los repone).

No se siembran personas, cuentas, reglas de reparto, categorías, tiendas ni créditos. Esa semilla la carga S8 (WF-4) desde `datos/semilla/data.json` y `oraculo/model.json`.

## Parámetros operativos de app_setting (K44)
Inventario de K44 («todo parámetro operativo existe y es configurable»), generado desde las anotaciones `-- @param`, `-- @param_sin_valor` y `-- @param_sin_valor_cerrado` de las migraciones (y `-- @param_rango` para los rangos aprobados de claves que ya tenían `@param`; 00005_tablero.sql es la última) y comprobado contra lo sembrado en la base migrada (30 de 30 valores iniciales coinciden). Todos son del hogar (`user_account_id` NULL). Tipos: enteros JSON (disparadores `trg_app_setting_param_int_ins` y `_upd` de la 00003; con su rango, `trg_app_setting_param_range_ins` y `_upd` de la 00004), la lista JSON de enteros de `alerts.payment_due_days_before` (`trg_app_setting_param_int_list_ins` y `_upd`) y las horas, texto JSON «HH:MM» de 00:00 a 23:59 en hora del centro de México (`trg_app_setting_param_time_ins` y `_upd`). Rangos aprobados en la puerta 2 (K50, Q2-66): los de una sola clave están en el esquema y se dicen en la columna «unidad»; la regla entre dos claves (`price.zero_weight_after_days` mayor o igual que `price.stale_after_days`, Q2-77) la validan el motor de precios y la API. Los demás enteros solo exigen 0 o más: no hay otros mínimos ni máximos aprobados. La sección de Ajustes es la PROPUESTA del orquestador en K44, la de U4 para lo que K44 no ubica (U4-12) y, para `ocr.timeout_seconds`, la propuesta de S1 (columna «decisión»). Regla de negocio: RN-29.

| clave | tipo | unidad | valor inicial | decisión | sección de Ajustes | migración | descripción sembrada |
|---|---|---|---|---|---|---|---|
| `auth.max_failed_attempts` | entero | intentos | 5 | Decisión A, K11 | Cuenta y roles | 00001 | Intentos fallidos permitidos dentro de la ventana antes de bloquear (decisión A). |
| `auth.lockout_window_minutes` | entero | minutos | 15 | Decisión A, K11 | Cuenta y roles | 00001 | Ventana en minutos para contar intentos fallidos (decisión A). |
| `auth.password_min_length` | entero | caracteres | 12 | K6 | Cuenta y roles | 00003 | Caracteres mínimos de una contraseña nueva. Entero en caracteres (K6, K44). |
| `session.max_age_days` | entero | días | 90 | K11 | Cuenta y roles | 00002 | Días que dura la sesión de un dispositivo; se renueva con cada uso (K11). |
| `price.stale_after_days` | entero | días, 1 o más (Q2-77) | 30 | Decisión E, K22, Q2-77 (K50) | Tiendas y precios | 00001 | Un precio es viejo cuando tiene esta edad en días o más (decisión E). |
| `price.zero_weight_after_days` | entero | días, igual o mayor que price.stale_after_days (Q2-77; lo valida el motor de precios) | 90 | K22, K40, Q2-77 (K50) | Tiendas y precios | 00003 | Días de edad desde los que un precio ya no compite en la comparativa: solo se muestra marcado. Entero en días. Entre price.stale_after_days y este valor el peso baja en línea recta (K22, K40, K44). |
| `price.profeco_import_max_gb` | entero | GB | 1 | K6 | Tiendas y precios | 00003 | Tamaño máximo de un archivo de PROFECO al importarlo. Entero en GB (K6, K44). |
| `comparator.timeout_seconds` | entero | segundos, de 1 a 30 (límites de seguridad, Q2-74) | 5 | K43, Q2-74 (K50) | Créditos y condiciones | 00003 | Tiempo máximo del comparador de escenarios; si se pasa, la API responde 503 unavailable. Entero en segundos (K43, K44). |
| `attachment.max_size_mb` | entero | MB | 10 | K43 | Datos y respaldos | 00003 | Tamaño máximo de un ticket o adjunto. Entero en MB (K43, K44). |
| `idempotency.retention_hours` | entero | horas | 24 | K6 | Datos y respaldos | 00002 | Horas que se guarda una clave de idempotency_key antes de purgarla (K6). |
| `health.goal_adjust_lose_kcal` | entero | kcal | 0 | K36 | Salud y metas de ahorro | 00003 | Kilocalorías que se restan de la meta sugerida cuando el objetivo es bajar. Entero en kcal, 0 o más. Valor de la casa: cada persona puede tener el suyo en su perfil de salud (K36, K44). |
| `health.goal_adjust_gain_kcal` | entero | kcal | 0 | K36 | Salud y metas de ahorro | 00003 | Kilocalorías que se suman a la meta sugerida cuando el objetivo es subir. Entero en kcal, 0 o más. Valor de la casa: cada persona puede tener el suyo en su perfil de salud (K36, K44). |
| `alerts.budget_warning_pct` | entero | % del presupuesto | 80 | Decisión F, K44 | Notificaciones | 00003 | Porcentaje del presupuesto, por categoría y en el total, en el que llega el primer aviso. Entero en por ciento (decisión F, K44). |
| `alerts.budget_limit_pct` | entero | % del presupuesto | 100 | Decisión F, K44 | Notificaciones | 00003 | Porcentaje del presupuesto, por categoría y en el total, en el que llega el aviso de tope. Entero en por ciento (decisión F, K44). |
| `alerts.payment_due_days_before` | lista de enteros | días antes del pago | [1, 0] | Decisión F, K44 | Notificaciones | 00003 | Días antes de cada pago en que llega el aviso de pago próximo: 1 es el día anterior y 0 el mismo día. Lista de enteros en días (decisión F, K44). |
| `reports.daily_send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "21:30" | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004 | Hora a la que sale el resumen diario. Texto HH:MM de 24 h, hora del centro de México (Q2-72, K31, K47). |
| `reports.weekly_weekday` | entero | día de la semana ISO 8601, de 1 a 7 (7 = domingo) | 7 | Decisión F, lugar K47 | Reportes y frecuencias | 00004 | Día de la semana del resumen semanal: 1 = lunes … 7 = domingo (ISO 8601). Entero de 1 a 7 (decisión F, K47). |
| `reports.weekly_send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "19:00" | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004 | Hora a la que sale el resumen semanal. Texto HH:MM de 24 h, hora del centro de México (Q2-72, K31, K47). |
| `reports.monthly_day_of_month` | entero | día del mes, de 1 a 31 (en un mes más corto sale su último día, Q2-66 B) | 1 | Decisión F, Q2-66 B, lugar K47 | Reportes y frecuencias | 00004 | Día del mes en que se genera el reporte mensual en PDF. Entero de 1 a 31; en un mes que no tiene ese día sale su último día (decisión F, Q2-66, K47). |
| `reports.monthly_send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "08:00" | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004 | Hora a la que se genera el reporte mensual en PDF. Texto HH:MM de 24 h, hora del centro de México (Q2-72, K31, K47). |
| `alerts.send_time` | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "09:00" | Q2-73 (K50), K31 punto 2, K43, lugar K47 | Notificaciones | 00004 | Hora a la que salen los avisos con hora. Texto HH:MM de 24 h, hora del centro de México (Q2-73, K31, K43, K47). |
| `alerts.stale_prices_advance_days` | entero | días antes del día de despensa, 0 o más | 1 | Q2-73 (K50), K31 punto 2, K43, lugar K47 | Notificaciones | 00004 | Días antes del día de despensa en que llega el aviso de precios viejos: 0 es el mismo día. Entero en días (Q2-73, K31, K43, K47). |
| `alerts.meal_plan_missing_weekday` | entero | día de la semana ISO 8601, de 1 a 7 (7 = domingo) | 7 | Decisión F, lugar K47 | Notificaciones | 00004 | Día de la semana del aviso de plan de comidas sin armar: 1 = lunes … 7 = domingo (ISO 8601). Entero de 1 a 7 (decisión F, K47). |
| `optimizer.date_window_days` | entero | días antes y después de la fecha (fecha ± días), 0 o más | 7 | Q2-75 (K50) | Créditos y condiciones | 00004 | Días antes y después de la fecha de una compra planeada que prueba el comparador de escenarios («fecha ± días»). Entero en días (Q2-75, K44). |
| `ocr.timeout_seconds` | entero | segundos por imagen, 0 o más | 30 | Q2-76 (K50); sección: propuesta S1 | Datos y respaldos | 00004 | Tiempo máximo de la lectura OCR de una imagen de ticket. Entero en segundos (Q2-76, K44). |
| `debt.capacity_warning_pct` | entero | % del ingreso mensual neto que se va en pagos de crédito del mes, frontera entre verde y amarillo del barómetro, de 1 a 100 | 30 | K54 decisión 2 (valor FIRME); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005 | Frontera entre verde y amarillo del barómetro de capacidad de endeudamiento: % del ingreso mensual neto que se va en pagos de crédito del mes. Con el entero redondeado, hasta este valor es verde. Entero de 1 a 100, menor que debt.capacity_limit_pct (K54, RN-30). |
| `debt.capacity_limit_pct` | entero | % del ingreso mensual neto que se va en pagos de crédito del mes, tope del barómetro (arriba es rojo), de 1 a 100, mayor que debt.capacity_warning_pct (lo valida la API con 422 validation_failed) | 40 | K54 decisión 2 (valor FIRME); rango y sección PROVISIONAL (bloqueo C1-tablero), K54 decisión 2; PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005 | Tope del barómetro de capacidad de endeudamiento: arriba de este % del ingreso mensual neto es rojo; de debt.capacity_warning_pct a este valor es amarillo. Entero de 1 a 100 (K54, RN-30). |
| `dashboard.payment_reminder_count` | entero | pagos que muestra «Pagos próximos» del tablero, de 1 a 20 | 5 | K54; valor PROVISIONAL (bloqueo U3-tablero-recordatorio); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005 | Cuántos pagos muestra «Pagos próximos» en el tablero de Inicio. Entero de 1 a 20 (K54; PROVISIONAL, bloqueo U3-tablero-recordatorio). |
| `budget.signal_months` | entero | meses cerrados seguidos que mira la sugerencia «subir o recortar presupuesto», de 2 a 12 | 3 | K54; valor PROVISIONAL (bloqueo U3-tablero-senal-presupuesto); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005 | Meses cerrados seguidos que mira la sugerencia «subir o recortar presupuesto» de cada categoría. Entero de 2 a 12 (K54, RN-31; PROVISIONAL, bloqueo U3-tablero-senal-presupuesto). |
| `budget.signal_margin_pct` | entero | puntos de % usado arriba o abajo del 100 % del plan para la sugerencia «subir o recortar presupuesto», de 1 a 50 | 10 | K54; valor PROVISIONAL (bloqueo U3-tablero-senal-presupuesto); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005 | Margen de la sugerencia «subir o recortar presupuesto»: subir si el % usado pasó de 100 + margen en todos los meses que mira; recortar si quedó abajo de 100 − margen. Entero de 1 a 50 (K54, RN-31; PROVISIONAL, bloqueo U3-tablero-senal-presupuesto). |

**Existen sin valor sembrado o viven en otra tabla** (no se inventa un valor):

| parámetro | para qué | por qué no se siembra |
|---|---|---|
| optimizer.min_available_buffer_cents, optimizer.max_live_purchases, optimizer.zero_debt_by_date | restricciones duras del comparador | las captura Axel en Ajustes «Créditos y condiciones» (K16); mientras falten, el comparador avisa «faltan restricciones» y no recomienda |
| reportes, respaldos, cabeceras, PDF y permisos | cinco de los ocho valores operativos de K31 | sin valor; los proponen S6 y WF-4 y se aprueban en la puerta 4 (K31) |
| ventana de la cola sin conexión | tiempo máximo que una captura espera en la cola del teléfono | no existe: K9 no pone caducidad (las claves de expense y weight_entry no se purgan); si se quiere un límite, lo propone la PWA en WF-3 |

## Claves de app_setting previstas sin valor
| clave | para qué | por qué no se siembra |
|---|---|---|
| `optimizer.min_available_buffer_cents` | colchón mínimo de disponible en la línea que se usa (restricción dura) | lo captura Axel en Ajustes «Créditos y condiciones» (K16); global y ajustable en cada comparación. Sin valor por defecto: mientras falte, el comparador avisa «faltan restricciones» y no recomienda |
| `optimizer.max_live_purchases` | máximo de compras vivas a la vez en el crédito que se usa (restricción dura) | igual que la anterior (K16) |
| `optimizer.zero_debt_by_date` | deuda en cero antes de una fecha en el crédito que se usa (restricción dura) | igual que la anterior (K16) |

Desde 1.2 ya no están en esta lista: `price.stale_weight` (el peso de un precio viejo es la regla RN-12, K22), `session.max_age_days` (90, K11) e `idempotency.retention_hours` (24, K6), que siembra la migración 00002. Desde 1.4 tampoco está `optimizer.date_window_days`: la migración 00004 la siembra con 7 (Q2-75, K50).

## Ajustes: 14 secciones → campos → tabla.columna (criterio S1-13)
Campos tomados de U1 §2 (P-26.01 a P-26.14) y §3 (N-270 a N-315). «Sin dato en la base» quiere decir que el campo no se guarda en SQLite y se dice dónde vive.

| # | Sección | Campos | Destino (tabla.columna) |
|---|---|---|---|
| 1 | Cuenta y roles | usuarios Axel y Ximena con los mismos permisos (nombre, usuario, papel, último acceso) · cambio de la contraseña propia · sesiones por dispositivo y su cierre · límite de intentos · «Este dispositivo» | `person.display_name`, `person.participation`, `user_account.username`, `user_account.is_active`, `user_account.last_login_at` (no hay columna de rol) · `user_account.password_hash`, `user_account.password_changed_at` · `device_session.device_label`, `.user_agent`, `.ip_address`, `.created_at`, `.last_seen_at`, `.expires_at`, `.revoked_at` · `login_attempt` + `app_setting` `auth.max_failed_attempts` (5), `auth.lockout_window_minutes` (15, ventana deslizante por usuario e IP, K11) y `session.max_age_days` (90, se renueva con el uso, K11) · sin dato en la base: si la app está instalada lo sabe el navegador (Local) y el certificado de la CA lo sirve S6 como archivo · parámetros de K44 en `app_setting` (valor inicial, RN-29): `auth.max_failed_attempts` (5), `auth.lockout_window_minutes` (15), `auth.password_min_length` (12), `session.max_age_days` (90) |
| 2 | Presupuesto | monto mensual por categoría y total, con vigencia · gastos fijos recurrentes (concepto, categoría, monto, día, reparto, activo) · despensa bimestral programada (monto base, día, frecuencia, línea y plazo) · proyección de los meses 7 a 12 con o sin la cadencia bimestral y el papel (K21) | `budget.category_id`, `.amount_cents`, `.valid_from_date`, `.valid_to_date` · `planned_item.kind` = `fixed_expense`, `.description`, `.category_id`, `.amount_cents`, `.frequency`, `.start_date` (día ancla), `.cost_share_rule_id`, `.is_active` · `planned_item.amount_cents`, `.frequency` = `monthly` con `.interval_count` = 2, `.start_date`, `.credit_id`, `.term_count`, `.accrual_months` = 2 (la próxima fecha la calcula S2) · `app_setting` `projection.repeat_cadence_enabled` (false; en qué pantalla de Ajustes va el interruptor lo decide U4) · parámetros de K44 en `app_setting` (valor inicial, RN-29): `debt.capacity_warning_pct` (30), `debt.capacity_limit_pct` (40), `dashboard.payment_reminder_count` (5), `budget.signal_months` (3), `budget.signal_margin_pct` (10) |
| 3 | Ingresos | ingreso quincenal esperado de Axel: monto de la quincena del 15, de la del 30 y desde cuándo rige | `planned_item.kind` = `income`, `.frequency` = `fortnightly`, `.amount_cents` (15), `.second_half_amount_cents` (30, si difiere), `.start_date`, `.end_date`, `.person_id` (solo el titular, disparador) |
| 4 | Categorías | nombre, tipo (fijo o variable), reparto por defecto, orden, activa | `category.name`, `category.kind`, `category.default_cost_share_rule_id`, `category.sort_order`, `category.is_active` (número de gastos: se cuenta en `expense`) |
| 5 | Créditos y condiciones (incluye promociones) | líneas por tipo y sus condiciones (tarjeta y crédito de tienda con los datos de su estado de cuenta, K18) · interruptor A/B con aviso «pendiente confirmar» · tasas por plazo · promociones · valores del comparador (colchón, compras vivas, deuda en cero, tope por persona, pesos) | `credit.kind`, `.credit_limit_cents`, `.annual_rate_num`/`.annual_rate_den`, `.statement_day`, `.payment_due_day`, `.payment_frequency`, `.first_payment_min_days`, comisiones `*_fee_cents`, `.vat_included` (00004, Q2-67 opción 1: el sí o no de «¿la tasa incluye IVA?»; NULL = falta capturar), `.vat_on_interest_bp` (el porcentaje de IVA que se suma, solo con `.vat_included` = 0), `.rate_base`, `.min_payment_pct_bp`, `.min_payment_fixed_cents`, `.annual_rate_decimals` (K18 y K12), `.is_active` · `credit.availability_reading` (def. B), `.availability_confirmed`, `.reading_a_anchor_date` · `credit_term_rate.term_count`, `.rate_num`/`.rate_den`, `.rate_decimals` · `promotion` (con `.rate_decimals`) · `app_setting` `optimizer.min_available_buffer_cents`, `optimizer.max_live_purchases`, `optimizer.zero_debt_by_date` (sin valor por defecto: las captura Axel, K16), `person.fortnight_payment_cap_cents` (se captura en «Personas y reparto», K13), `app_setting` `optimizer.weights` · parámetros de K44 en `app_setting` (valor inicial, RN-29): `comparator.timeout_seconds` (5), `optimizer.date_window_days` (7) |
| 6 | Tiendas y precios (incluye mercados y membresía Sam's) | cadenas y sucursales · mercados y tianguis (días de visita, qué se compra ahí) · membresía de Sam's (costo anual) · fuentes de precio y lectura asistida | `store.kind`, `.parent_id`, `.name`, `.address`, `.municipality`, `.state`, `.profeco_name` · `store.visit_weekday`, `store_product_group.product_group_id` · `store.requires_membership`, `.has_membership`, `.membership_annual_fee_cents` · `price.source`, `store.assisted_reading_enabled` + `app_setting` `price.assisted_reading_enabled` (ambos apagados), `app_setting` `price.stale_after_days` (30) · parámetros de K44 en `app_setting` (valor inicial, RN-29): `price.stale_after_days` (30), `price.zero_weight_after_days` (90), `price.profeco_import_max_gb` (1) |
| 7 | Recetas | porciones base, tiempos de comida, recetario inicial, bases nutricionales locales | `recipe.base_servings` (def. 2) · `meal_slot` (4 tiempos) y `recipe_meal_slot` · `recipe.origin` = `claude_seed` · `nutrition_fact.source`, `.source_ref`, `.source_license` (fecha de carga: `nutrition_fact.created_at`) |
| 8 | Reportes y frecuencias | resumen diario, semanal y mensual PDF (activo, día, hora), vista web mensual, destinatarios | `notification.is_enabled`, `job.schedule`, `job.next_run_at` (desde 1.4 el día y la hora de cada reporte son claves de `app_setting`, K47; `notification.params_json` ya no los lleva) · la vista web mensual no tiene ajuste propio (se calcula al consultar; los PDF quedan en `generated_report`) · `notification_subscription.notification_id`, `.recipient_id` · parámetros de K44 en `app_setting` (valor inicial, RN-29): `reports.daily_send_time` ("21:30"), `reports.weekly_weekday` (7), `reports.weekly_send_time` ("19:00"), `reports.monthly_day_of_month` (1), `reports.monthly_send_time` ("08:00") |
| 9 | Notificaciones | 4 alertas (activa, destinatarios, anticipación) · registro del chat de Axel y de Ximena | `notification.is_enabled`, `notification.params_json` (`scopes` de `budget_threshold`; desde 1.3 los umbrales y los días antes, y desde 1.4 la hora, la anticipación de precios viejos y el día del plan sin armar, son claves de `app_setting`, K44 y K47), `notification_subscription` · `notification_recipient.person_id`, `.chat_id`, `.verified_at`, `.is_active` · chats que detecta el bot y la persona confirma (K4): `telegram_chat_candidate.chat_id`, `.telegram_name`, `.detected_at`, `.confirmed_at`, `.person_id` · parámetros de K44 en `app_setting` (valor inicial, RN-29): `alerts.budget_warning_pct` (80), `alerts.budget_limit_pct` (100), `alerts.payment_due_days_before` ([1, 0]), `alerts.send_time` ("09:00"), `alerts.stale_prices_advance_days` (1), `alerts.meal_plan_missing_weekday` (7) |
| 10 | Apariencia | tema: sistema (por defecto), claro u oscuro | sin dato en la base: el tema vive en cada dispositivo (U1 N-301 «Local», C1 §11.3). Saturación apagada y densidad cómoda son decisiones fijas del diseño (C2), no ajustes |
| 11 | Datos y respaldos | respaldos (lista, integridad, restauración) · importar `gastos_fijos_ingresos.csv` · estado de la semilla · IA opcional · estado del servicio | `backup.kind`, `.status`, `.started_at`, `.finished_at`, `.file_path`, `.size_bytes`, `.integrity_check`, `.restore_tested_at`, `.restore_test_result` y `job.next_run_at` · el CSV se importa a `planned_item` · la semilla se deduce de lo cargado (`expense`, `credit_purchase`, `planned_item`) · `app_setting` `ai.local_model_enabled` y `ai.claude_api_enabled` (false) · sin dato en la base: versión y tiempo en marcha los da S6; la llave de Claude nunca se guarda en la base · parámetros de K44 en `app_setting` (valor inicial, RN-29): `attachment.max_size_mb` (10), `idempotency.retention_hours` (24), `ocr.timeout_seconds` (30) |
| 12 | Personas y reparto | Axel, Ximena y Karen (nombre, papel, paga directo, tope por quincena opcional) · reparto por defecto del conjunto (50/50) | `person.display_name`, `.participation`, `.pays_share_directly`, `.fortnight_payment_cap_cents` (opcional, K13) · `cost_share_rule.is_default`, `.is_template`, `.remainder_person_id`, `cost_share_rule_part.person_id`, `.pct_bp` |
| 13 | Salud y metas de ahorro | perfiles de Axel y Ximena (con objetivo de peso, K23) · metas con fecha · aportación planificada al fondo | `health_profile.height_mm`, `.birth_date`, `.sex`, `.activity_level`, `.goal`, `.goal_adjust_lose_kcal`, `.goal_adjust_gain_kcal` (K36; NULL = valor de la casa en `app_setting` `health.goal_adjust_lose_kcal` y `health.goal_adjust_gain_kcal`), `.daily_kcal_target`, `.dietary_restrictions` y `weight_entry` · `savings_goal.kind` = `goal`, `.name`, `.target_cents`, `.target_date`, `.planned_monthly_cents`, `.is_active` · `savings_goal.kind` = `fund`, `.planned_monthly_cents` · parámetros de K44 en `app_setting` (valor inicial, RN-29): `health.goal_adjust_lose_kcal` (0), `health.goal_adjust_gain_kcal` (0) |
| 14 | Inventario | modo por producto, unidad, umbral de reposición o ritmo inicial · regla de sobrantes | `inventory_item.tracking_mode`, `product.base_unit`, `inventory_item.reorder_level_milli`, `.rate_quantity_milli`, `.rate_period_days` · regla fija de la decisión D: un sobrante es un `inventory_item` con `recipe_id` y `tracking_mode = 'stock'` |

## Invariantes que SQLite no puede expresar (las validan la API y los motores)
1. Σ `expense_item.line_total_cents` = `expense.amount_cents` cuando el gasto tiene renglones.
2. Σ `expense_share.amount_cents` = `expense.amount_cents`.
3. Σ `scheduled_payment.amount_cents` = `credit_purchase.total_cents`, Σ `principal_cents` = capital de la compra y Σ `interest_cents` = intereses de la compra.
4. Σ `scheduled_payment_share.amount_cents` de un pago = `amount_cents` del pago, y Σ por persona en la serie = `credit_purchase_share.total_cents` de esa persona.
5. Σ de `credit_purchase_share` (capital, intereses, total) = valores de la compra.
6. `expense.payment_method = 'credit'` si y solo si existe una `credit_purchase` con ese `expense_id`, y su `principal_cents` = `expense.amount_cents`.
7. `period.parent_id` de una quincena apunta al mes del mismo año y mes. Los renglones de `budget` de una misma categoría (o del total) no se traslapan en su vigencia: el monto de un mes es el del renglón cuya vigencia contiene ese mes.
8. En `cost_share_rule`, la persona del resto no aparece en sus partes y Σ `pct_bp` ≤ 10000.
9. Las partes de dinero (`*_share`, `cost_share_rule_part`) solo usan al titular y a personas externas. La co-titular no lleva parte propia porque el hogar va a nombre del titular.
10. El mismo día, tiempo y persona no mezcla «fuera» con recetas.
11. Las fechas de `scheduled_payment` siguen el calendario del crédito (Kueski: días de pago de quincena y primer pago al menos 15 días después de la compra).
12. En `inventory_item` con `tracking_mode = 'stock'`, `quantity_milli` es la existencia vigente y cada cambio escribe su `inventory_movement` en la misma transacción. Desde 1.1 el paso de una persona a `external` con cuenta, salud o plan de comidas lo rechaza el disparador `trg_person_external_without_access`.
13. En `shopping_list_item`, la regeneración de la lista (N-127) actualiza `needed_milli` de los artículos calculados y conserva lo que eligió una persona (`quantity_milli`, `store_id`, `presentation_id`, `is_bought`). «Usar esta tienda» (N-165) pone `store_id` en los artículos de la lista vigente que tienen precio en esa tienda.
14. (1.2) Un `telegram_chat_candidate` confirmado tiene su `notification_recipient` con el mismo `chat_id` y la misma persona; la API hace los dos cambios en una transacción (K4).
15. (1.2) `expense.idempotency_key` y `weight_entry.idempotency_key` llevan el valor del encabezado `Idempotency-Key` de la captura; una captura repetida devuelve el registro que ya existe aunque su renglón de `idempotency_key` ya se haya purgado (K9).
16. (1.2) En `credit` de tipo `bank_card` o `store_credit`, la API guarda juntos los datos de K18 (tasa anual con sus decimales, `rate_base`, si incluye IVA en `vat_included`, desde 1.4, con su porcentaje en `vat_on_interest_bp` cuando no lo incluye, `min_payment_pct_bp` y `min_payment_fixed_cents`); si falta alguno, S3 muestra «falta capturar» y no inventa tasa ni mínimo (desde 1.4 son 5 datos: `vat_included` NULL también falta, Q2-67). El esquema acepta `vat_included` = 0 con `vat_on_interest_bp` NULL («no incluye IVA» con el porcentaje por escribir; U4 captura las dos cosas por separado). **PROVISIONAL** (bloqueo abierto `S3-iva-lectura`, pregunta 2, sin recomendación): con ese par falta el porcentaje, `Credit.missing_terms` lista `vat` y S3 muestra «falta capturar» (`ErrMissingCapture`, la lectura de `S3-interfaz.md` 1.5 §5.3); no lo toma como 0 % ni inventa un porcentaje (K18). Si Axel elige la opción B del bloqueo, no se podrá guardar 0 sin porcentaje (con un CHECK en una migración posterior o con una validación de C1) y S3 responderá `ErrInvalid` si lo recibe.
17. (1.4) `price.zero_weight_after_days` es mayor o igual que `price.stale_after_days`, que es 1 o más (Q2-77): el esquema solo prueba el 1 o más; la regla entre las dos claves la valida el motor de precios (`backend/internal/precios/types.go`, `Settings.Validate`, `invalid_argument`) y la API al guardar cualquiera de las dos (RN-12).
18. (1.4) Alcance de un cambio de regla de reparto reutilizable (Q2-68 C): la pantalla pregunta «solo los nuevos» o «también los ya capturados». El esquema no cambia: las partes de cada gasto ya están congeladas en `expense_share`; con «solo los nuevos» no se tocan, y con «también los ya capturados» el motor de reparto las recalcula. El parámetro de alcance es de C1 y el recálculo del motor; qué gastos abarca lo fijan ellos en WF-3 (RN-03). **PROVISIONAL** (bloqueo `S1-reparto-edicion-gasto`): con qué partes se recalcula un gasto anterior que se edita después de un cambio con «solo los nuevos» no está decidido (`expense_share` guarda solo `amount_cents`, sin porcentaje ni versión de la regla). Mientras tanto queda la opción en uso: `updateExpense` (C1 v1.2) recalcula con la regla vigente. Según la respuesta puede hacer falta una columna nueva en una migración posterior.

## Notas para otras fases
- **S6:** activar `foreign_keys`, `journal_mode=WAL` y `busy_timeout` en cada conexión (DSN de modernc `_pragma=...`). Dentro de la transacción de goose el `PRAGMA foreign_keys` de la migración no tiene efecto. Escribir `audit_log` en la misma transacción de cada cambio.
- **S3/S2:** `credit_purchase`, `scheduled_payment` y sus partes se guardan ya calculadas al centavo con las reglas de `reglas-negocio.md`. Los motores las generan, la base no las recalcula.
- **U1/S6:** la tabla `idempotency_key` guarda 24 h las claves del encabezado `Idempotency-Key` (C1 §1 punto 8, K6). Para la captura sin conexión (K9) la 00002 agrega `idempotency_key` (única, sin purga) y `captured_at` a `expense` y `weight_entry`: la API busca la clave antes de insertar y responde con el registro que ya existe.
- **S4/C1:** la lista de compras se comparte en el servidor (K1) y guarda solo lo que no se calcula (artículos a mano, tienda y cantidad elegidas, marcas de comprado).
- **U1/S2:** semana de lunes a domingo y quincenas del 1 al 15 y del 16 al último día (K8): el CHECK de `period` queda como está.
- **S5:** si PROFECO necesita una tabla cruda de carga, va en una migración nueva. Aquí solo entran precios ya ligados a una `presentation`.
- **S7:** `telegram_chat_candidate` guarda solo el `chat_id` y el nombre que da Telegram (K4). Cómo convive la detección con MultiX, que ya consume las actualizaciones del bot, lo verifica S7; si no hay forma limpia deja bloqueo (K4).
- **S3:** una tarjeta o un crédito de tienda sin los datos de K18 muestra «falta capturar»; sin las tres restricciones `optimizer.*`, el comparador avisa «faltan restricciones» y no recomienda (K16).
