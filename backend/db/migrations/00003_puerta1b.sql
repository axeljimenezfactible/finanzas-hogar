-- =====================================================================
-- Finanzas Hogar · Migración 00003 · Puerta 1b: esquema 1.2 → 1.3
-- Agente S1-esquema-1.3 (flujo P1c «Aplicar respuestas de la puerta 1b»), 11-oct-2026.
-- Aplica las respuestas de Axel de docs/plan/decisiones.md §K que tocan el esquema:
--   K36 ajuste de la meta de calorías por objetivo, por persona y con valor de la casa.
--   K44 todo parámetro operativo existe con su valor decidido y es configurable en Ajustes:
--       cada uno es una clave de app_setting sembrada con el valor de §A–§K (F, K6, K22, K36, K40, K43).
--       Los umbrales del presupuesto (80 y 100 %) y los días antes del pago próximo ([1, 0]) de la decisión F
--       pasan de notification.params_json a app_setting (alerts.*): una sola fuente (corrección S1-r1, H-01).
--       Los días de programación de la decisión F (weekday 7 de weekly_summary y meal_plan_missing, day_of_month 1
--       de monthly_report) siguen en notification.params_json con su valor: dónde viven es PROVISIONAL, igual que
--       send_time, bloqueo S1-parametros-avisos (corrección S1-r2, M-01).
--   K37, K39 y K40 no tocan el esquema (son reglas RN-12 y RN-27); K40 sí siembra el día de peso cero.
-- Contrato: docs/contratos/diccionario-datos.md 1.3 y docs/contratos/reglas-negocio.md 1.3 (RN-29).
-- Convenciones: las de 00001 y 00002. Notas técnicas:
--   * Las columnas nuevas entran con ALTER TABLE ... ADD COLUMN, admiten NULL (el comentario dice qué
--     significa NULL) y su CHECK es de la propia columna, para que el Down pueda quitarla.
--   * Los parámetros son enteros en JSON (sin comillas), salvo alerts.payment_due_days_before, que es una
--     lista JSON de enteros. Los disparadores trg_app_setting_param_int_* y trg_app_setting_param_int_list_*
--     rechazan otro tipo o un negativo en las claves del inventario; NO fijan mínimos ni máximos
--     razonables: esos los propone S1 y se aprueban en la puerta 2 (K44), mientras tanto sin tope.
--   * El Down repone days_before y thresholds_pct en notification.params_json con el valor vigente de
--     app_setting, borra los disparadores, los ajustes que siembra esta migración (en cualquier alcance) y
--     las dos columnas, y deja el esquema 1.2 idéntico: sqlite_master y los datos de 1.2 vuelven byte
--     por byte con los valores iniciales (lo prueba backend/db/tools/s1_verifica.py).
--   * Anotaciones nuevas que lee backend/db/tools/gen_dict.py:
--     -- @param <clave> | <tipo> | <unidad> | <valor inicial JSON> | <decisión> | <sección de Ajustes> | <migración>
--        (inventario de parámetros operativos, K44; incluye los que ya sembraron 00001 y 00002)
--     -- @param_sin_valor <clave o parámetro> | <para qué> | <quién lo propone y por qué no se siembra>
--     -- @col_note <tabla>.<columna> | <nota de 1.3 que el diccionario añade a la descripción de la columna>
--   * Nombres: se conservan las claves que ya existen (auth.*, session.*, idempotency.*, price.*,
--     optimizer.*); las nuevas usan el prefijo de su módulo (health.*, price.*, auth.*, attachment.*,
--     comparator.*, alerts.*). Ningún DDL de 00001 ni de 00002 cambia.
-- =====================================================================

-- Inventario de parámetros operativos (K44). Sección de Ajustes: PROPUESTA del orquestador en K44
-- (cada parámetro vive en la sección a la que pertenece, sin sección nueva); se aprueba en la puerta 2.
-- @param auth.max_failed_attempts | entero | intentos | 5 | Decisión A, K11 | Cuenta y roles | 00001
-- @param auth.lockout_window_minutes | entero | minutos | 15 | Decisión A, K11 | Cuenta y roles | 00001
-- @param auth.password_min_length | entero | caracteres | 12 | K6 | Cuenta y roles | 00003
-- @param session.max_age_days | entero | días | 90 | K11 | Cuenta y roles | 00002
-- @param price.stale_after_days | entero | días | 30 | Decisión E, K22 | Tiendas y precios | 00001
-- @param price.zero_weight_after_days | entero | días | 90 | K22, K40 | Tiendas y precios | 00003
-- @param price.profeco_import_max_gb | entero | GB | 1 | K6 | Tiendas y precios | 00003
-- @param comparator.timeout_seconds | entero | segundos | 5 | K43 | Créditos y condiciones | 00003
-- @param attachment.max_size_mb | entero | MB | 10 | K43 | Datos y respaldos | 00003
-- @param idempotency.retention_hours | entero | horas | 24 | K6 | Datos y respaldos | 00002
-- @param health.goal_adjust_lose_kcal | entero | kcal | 0 | K36 | Salud y metas de ahorro | 00003
-- @param health.goal_adjust_gain_kcal | entero | kcal | 0 | K36 | Salud y metas de ahorro | 00003
-- @param alerts.budget_warning_pct | entero | % del presupuesto | 80 | Decisión F, K44 | Notificaciones | 00003
-- @param alerts.budget_limit_pct | entero | % del presupuesto | 100 | Decisión F, K44 | Notificaciones | 00003
-- @param alerts.payment_due_days_before | lista de enteros | días antes del pago | [1, 0] | Decisión F, K44 | Notificaciones | 00003
-- Parámetros que existen sin valor sembrado o que viven en otra tabla (no se inventa un valor):
-- @param_sin_valor optimizer.min_available_buffer_cents, optimizer.max_live_purchases, optimizer.zero_debt_by_date | restricciones duras del comparador | las captura Axel en Ajustes «Créditos y condiciones» (K16); mientras falten, el comparador avisa «faltan restricciones» y no recomienda
-- @param_sin_valor optimizer.date_window_days | «fecha ± días» de los escenarios | valor no dado; lo propone S3 con Axel
-- @param_sin_valor notification.params_json send_time | hora de envío (HH:MM, America/Mexico_City) del resumen diario y de los avisos con hora | clave aprobada en K43 sin valor; lo propone S7 (WF-1). Si vive en params_json (K43) o en app_setting (K44) es PROVISIONAL, bloqueo S1-parametros-avisos
-- @param_sin_valor notification.params_json advance_days | anticipación en días del aviso de precios viejos antes del día de despensa | clave aprobada en K43 sin valor; lo propone S7 (WF-1). Si vive en params_json (K43) o en app_setting (K44) es PROVISIONAL, bloqueo S1-parametros-avisos
-- @param_sin_valor notification.params_json weekday (weekly_summary y meal_plan_missing) | día de la semana del resumen semanal y del aviso de plan de comidas sin armar (1 = lunes … 7 = domingo, ISO 8601) | vive en otra tabla: valor 7 (domingo) de la decisión F, sembrado por 00001 en notification.params_json y configurable en Ajustes; no se siembra en app_setting. Si vive en params_json (como send_time, K43) o en app_setting (K44, como alerts.payment_due_days_before) es PROVISIONAL, bloqueo S1-parametros-avisos
-- @param_sin_valor notification.params_json day_of_month (monthly_report) | día del mes en que se genera el reporte mensual en PDF | vive en otra tabla: valor 1 (día 1) de la decisión F, sembrado por 00001 en notification.params_json y configurable en Ajustes; no se siembra en app_setting. Si vive en params_json (como send_time, K43) o en app_setting (K44, como alerts.payment_due_days_before) es PROVISIONAL, bloqueo S1-parametros-avisos
-- @param_sin_valor reportes, respaldos, cabeceras, PDF y permisos | cinco de los ocho valores operativos de K31 | sin valor; los proponen S6 y WF-4 y se aprueban en la puerta 4 (K31)
-- @param_sin_valor ventana de la cola sin conexión | tiempo máximo que una captura espera en la cola del teléfono | no existe: K9 no pone caducidad (las claves de expense y weight_entry no se purgan); si se quiere un límite, lo propone la PWA en WF-3
-- @col_note notification.params_json | Desde 1.3 (00003, K44) `payment_due` ya no lleva `days_before` ni `budget_threshold` lleva `thresholds_pct`: esos valores son las claves `alerts.payment_due_days_before`, `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `app_setting` (RN-29). El Down los repone aquí. `weekday` (7 en `weekly_summary` y `meal_plan_missing`) y `day_of_month` (1 en `monthly_report`) de la decisión F siguen aquí: dónde viven es PROVISIONAL, bloqueo S1-parametros-avisos

-- +goose Up

-- ---------------------------------------------------------------------
-- K36. Ajuste de la meta de calorías por objetivo, por persona
-- ---------------------------------------------------------------------
-- @alter health_profile
-- @used_new Pantalla Salud y Ajustes «Salud y metas de ahorro» (ajuste por persona; el de la casa en app_setting health.goal_adjust_*). S4 (meta sugerida, regla RN-17). C1 v1.2. U4.
ALTER TABLE health_profile ADD COLUMN goal_adjust_lose_kcal INTEGER CHECK (goal_adjust_lose_kcal >= 0); -- Nueva en 1.3 (K36). Cuántas kcal se restan de la meta sugerida cuando el objetivo (goal) es lose. Es una magnitud: siempre se resta. NULL = usa el valor de la casa (app_setting health.goal_adjust_lose_kcal, 0 al sembrar) | entero mayor o igual a 0, en kcal; NULL = valor de la casa
ALTER TABLE health_profile ADD COLUMN goal_adjust_gain_kcal INTEGER CHECK (goal_adjust_gain_kcal >= 0); -- Nueva en 1.3 (K36). Cuántas kcal se suman a la meta sugerida cuando el objetivo (goal) es gain. Es una magnitud: siempre se suma. maintain no ajusta. NULL = usa el valor de la casa (app_setting health.goal_adjust_gain_kcal, 0 al sembrar) | entero mayor o igual a 0, en kcal; NULL = valor de la casa

-- ---------------------------------------------------------------------
-- K44. Parámetros operativos nuevos, sembrados con su valor decidido
-- (los de 00001 y 00002 ya existen: auth.max_failed_attempts, auth.lockout_window_minutes,
-- price.stale_after_days, session.max_age_days, idempotency.retention_hours)
-- ---------------------------------------------------------------------
INSERT INTO app_setting (setting_key, value_json, description) VALUES
    ('health.goal_adjust_lose_kcal', '0', 'Kilocalorías que se restan de la meta sugerida cuando el objetivo es bajar. Entero en kcal, 0 o más. Valor de la casa: cada persona puede tener el suyo en su perfil de salud (K36, K44).'),
    ('health.goal_adjust_gain_kcal', '0', 'Kilocalorías que se suman a la meta sugerida cuando el objetivo es subir. Entero en kcal, 0 o más. Valor de la casa: cada persona puede tener el suyo en su perfil de salud (K36, K44).'),
    ('price.zero_weight_after_days', '90', 'Días de edad desde los que un precio ya no compite en la comparativa: solo se muestra marcado. Entero en días. Entre price.stale_after_days y este valor el peso baja en línea recta (K22, K40, K44).'),
    ('auth.password_min_length', '12', 'Caracteres mínimos de una contraseña nueva. Entero en caracteres (K6, K44).'),
    ('attachment.max_size_mb', '10', 'Tamaño máximo de un ticket o adjunto. Entero en MB (K43, K44).'),
    ('price.profeco_import_max_gb', '1', 'Tamaño máximo de un archivo de PROFECO al importarlo. Entero en GB (K6, K44).'),
    ('comparator.timeout_seconds', '5', 'Tiempo máximo del comparador de escenarios; si se pasa, la API responde 503 unavailable. Entero en segundos (K43, K44).');

-- Los parámetros enteros del inventario (K44) solo aceptan un entero JSON mayor o igual a 0.
-- No son mínimos razonables (esos se aprueban en la puerta 2): solo el tipo y el signo.
-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_int_ins BEFORE INSERT ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('auth.max_failed_attempts', 'auth.lockout_window_minutes', 'auth.password_min_length', 'session.max_age_days', 'price.stale_after_days', 'price.zero_weight_after_days', 'price.profeco_import_max_gb', 'comparator.timeout_seconds', 'attachment.max_size_mb', 'idempotency.retention_hours', 'health.goal_adjust_lose_kcal', 'health.goal_adjust_gain_kcal', 'alerts.budget_warning_pct', 'alerts.budget_limit_pct')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'integer' AND CAST(NEW.value_json AS INTEGER) >= 0)
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro es un entero mayor o igual a 0 (K44)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_int_upd BEFORE UPDATE OF setting_key, value_json ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('auth.max_failed_attempts', 'auth.lockout_window_minutes', 'auth.password_min_length', 'session.max_age_days', 'price.stale_after_days', 'price.zero_weight_after_days', 'price.profeco_import_max_gb', 'comparator.timeout_seconds', 'attachment.max_size_mb', 'idempotency.retention_hours', 'health.goal_adjust_lose_kcal', 'health.goal_adjust_gain_kcal', 'alerts.budget_warning_pct', 'alerts.budget_limit_pct')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'integer' AND CAST(NEW.value_json AS INTEGER) >= 0)
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro es un entero mayor o igual a 0 (K44)');
END;
-- +goose StatementEnd

-- La anticipación del pago próximo es una lista JSON de enteros mayores o iguales a 0. Tampoco aquí
-- hay mínimos ni cuántos elementos (eso se aprueba en la puerta 2): solo el tipo y el signo.
-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_int_list_ins BEFORE INSERT ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('alerts.payment_due_days_before')
    AND (CASE WHEN NOT json_valid(NEW.value_json) THEN 1
              WHEN json_type(NEW.value_json) <> 'array' THEN 1
              ELSE EXISTS (SELECT 1 FROM json_each(NEW.value_json) WHERE type <> 'integer' OR value < 0) END)
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro es una lista de enteros mayores o iguales a 0 (K44)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_int_list_upd BEFORE UPDATE OF setting_key, value_json ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('alerts.payment_due_days_before')
    AND (CASE WHEN NOT json_valid(NEW.value_json) THEN 1
              WHEN json_type(NEW.value_json) <> 'array' THEN 1
              ELSE EXISTS (SELECT 1 FROM json_each(NEW.value_json) WHERE type <> 'integer' OR value < 0) END)
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro es una lista de enteros mayores o iguales a 0 (K44)');
END;
-- +goose StatementEnd

-- ---------------------------------------------------------------------
-- K44 y decisión F. Umbrales del presupuesto y días antes del pago próximo: pasan de
-- notification.params_json (00001) a app_setting, una sola fuente. Cada clave toma el valor que el
-- aviso tenga al migrar (en una base nueva, 80, 100 y [1, 0], los de la decisión F); si falta o no
-- es válido, el de la decisión F.
-- ---------------------------------------------------------------------
INSERT INTO app_setting (setting_key, value_json, description) VALUES
    ('alerts.budget_warning_pct',
     coalesce((SELECT CAST(json_extract(n.params_json, '$.thresholds_pct[0]') AS TEXT) FROM notification n
               WHERE n.code = 'budget_threshold' AND json_type(n.params_json, '$.thresholds_pct[0]') = 'integer'
                 AND json_extract(n.params_json, '$.thresholds_pct[0]') >= 0), '80'),
     'Porcentaje del presupuesto, por categoría y en el total, en el que llega el primer aviso. Entero en por ciento (decisión F, K44).'),
    ('alerts.budget_limit_pct',
     coalesce((SELECT CAST(json_extract(n.params_json, '$.thresholds_pct[1]') AS TEXT) FROM notification n
               WHERE n.code = 'budget_threshold' AND json_type(n.params_json, '$.thresholds_pct[1]') = 'integer'
                 AND json_extract(n.params_json, '$.thresholds_pct[1]') >= 0), '100'),
     'Porcentaje del presupuesto, por categoría y en el total, en el que llega el aviso de tope. Entero en por ciento (decisión F, K44).'),
    ('alerts.payment_due_days_before',
     coalesce((SELECT replace(json_extract(n.params_json, '$.days_before'), ',', ', ') FROM notification n
               WHERE n.code = 'payment_due' AND json_type(n.params_json, '$.days_before') = 'array'
                 AND NOT EXISTS (SELECT 1 FROM json_each(n.params_json, '$.days_before') WHERE type <> 'integer' OR value < 0)), '[1, 0]'),
     'Días antes de cada pago en que llega el aviso de pago próximo: 1 es el día anterior y 0 el mismo día. Lista de enteros en días (decisión F, K44).');

-- Se quitan de params_json para no dejar dos fuentes. updated_at se conserva porque la casa no cambió
-- nada: va en dos pasos porque trg_notification_updated_at lo renueva si el UPDATE no lo trae.
UPDATE notification SET updated_at = updated_at || '~',
    params_json = CASE WHEN params_json = '{"days_before": [1, 0]}' THEN '{}' ELSE json_remove(params_json, '$.days_before') END
WHERE code = 'payment_due' AND json_type(params_json, '$.days_before') IS NOT NULL;
UPDATE notification SET updated_at = updated_at || '~',
    params_json = CASE WHEN params_json = '{"thresholds_pct": [80, 100], "scopes": ["category", "total"]}' THEN '{"scopes": ["category", "total"]}' ELSE json_remove(params_json, '$.thresholds_pct') END
WHERE code = 'budget_threshold' AND json_type(params_json, '$.thresholds_pct') IS NOT NULL;
UPDATE notification SET updated_at = substr(updated_at, 1, length(updated_at) - 1)
WHERE code IN ('payment_due', 'budget_threshold') AND updated_at LIKE '%~';

-- +goose Down
-- Deshace la 00003 en orden inverso y deja el esquema 1.2 idéntico.
-- Repone days_before y thresholds_pct al principio de notification.params_json con el valor vigente de
-- app_setting (con los valores iniciales, el texto de 1.2 byte por byte) y conserva updated_at.
UPDATE notification SET updated_at = updated_at || '~',
    params_json = (SELECT CASE WHEN (SELECT count(*) FROM json_each(notification.params_json)) = 0 THEN '{"days_before": ' || v.d || '}'
                               WHEN substr(notification.params_json, 1, 1) = '{' THEN '{"days_before": ' || v.d || ', ' || substr(notification.params_json, 2)
                               ELSE json_set(notification.params_json, '$.days_before', json(v.d)) END
                   FROM (SELECT replace(json(coalesce((SELECT value_json FROM app_setting WHERE setting_key = 'alerts.payment_due_days_before' AND user_account_id IS NULL), '[1, 0]')), ',', ', ') AS d) v)
WHERE code = 'payment_due' AND json_type(params_json) = 'object' AND json_type(params_json, '$.days_before') IS NULL;
UPDATE notification SET updated_at = updated_at || '~',
    params_json = (SELECT CASE WHEN (SELECT count(*) FROM json_each(notification.params_json)) = 0 THEN '{"thresholds_pct": ' || v.t || '}'
                               WHEN substr(notification.params_json, 1, 1) = '{' THEN '{"thresholds_pct": ' || v.t || ', ' || substr(notification.params_json, 2)
                               ELSE json_set(notification.params_json, '$.thresholds_pct', json(v.t)) END
                   FROM (SELECT '[' || json(coalesce((SELECT value_json FROM app_setting WHERE setting_key = 'alerts.budget_warning_pct' AND user_account_id IS NULL), '80'))
                                || ', ' || json(coalesce((SELECT value_json FROM app_setting WHERE setting_key = 'alerts.budget_limit_pct' AND user_account_id IS NULL), '100')) || ']' AS t) v)
WHERE code = 'budget_threshold' AND json_type(params_json) = 'object' AND json_type(params_json, '$.thresholds_pct') IS NULL;
UPDATE notification SET updated_at = substr(updated_at, 1, length(updated_at) - 1)
WHERE code IN ('payment_due', 'budget_threshold') AND updated_at LIKE '%~';
DROP TRIGGER IF EXISTS trg_app_setting_param_int_list_upd;
DROP TRIGGER IF EXISTS trg_app_setting_param_int_list_ins;
DROP TRIGGER IF EXISTS trg_app_setting_param_int_upd;
DROP TRIGGER IF EXISTS trg_app_setting_param_int_ins;
DELETE FROM app_setting WHERE setting_key IN ('health.goal_adjust_lose_kcal', 'health.goal_adjust_gain_kcal', 'price.zero_weight_after_days', 'auth.password_min_length', 'attachment.max_size_mb', 'price.profeco_import_max_gb', 'comparator.timeout_seconds', 'alerts.budget_warning_pct', 'alerts.budget_limit_pct', 'alerts.payment_due_days_before');
ALTER TABLE health_profile DROP COLUMN goal_adjust_gain_kcal;
ALTER TABLE health_profile DROP COLUMN goal_adjust_lose_kcal;
