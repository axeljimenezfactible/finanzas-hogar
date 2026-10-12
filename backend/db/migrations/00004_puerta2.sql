-- =====================================================================
-- Finanzas Hogar · Migración 00004 · Puerta 2: esquema 1.3 → 1.4
-- Agente S1-puerta2 (flujo P2c «Aplicación de la puerta 2»), 11-oct-2026; correcciones correccion-S1-r1 (M-01 de REVISION-S1-r1) y
-- correccion-S1-r2 (N-02 de REVISION-S1-r2: qué es vat_included = 0 sin porcentaje; solo comentarios).
-- Aplica las respuestas de Axel de la puerta 2 que tocan el esquema (docs/estado/puerta2/RESPUESTAS.md,
-- docs/estado/puerta2/PUERTA-2.md y docs/plan/decisiones.md K46, K47 y K50):
--   Q2-11 B (K47): la programación y la anticipación de los avisos son claves de app_setting, una por valor:
--       el día de la semana (7) y el día del mes (1) de la decisión F salen de notification.params_json
--       (el Down los repone con el valor vigente, como days_before en la 00003); la hora de cada reporte
--       (Q2-72: 21:30, domingo 19:00, día 1 a las 08:00) y la hora y la anticipación de los avisos
--       (Q2-73: 09:00 y 1 día) nacen aquí con su valor aprobado. Cierra el bloqueo S1-parametros-avisos.
--   Q2-75 y Q2-76 (K50): optimizer.date_window_days = 7 y ocr.timeout_seconds = 30.
--   Q2-74, Q2-77 y Q2-66 B (K50, K46): rangos aprobados de UNA clave en el esquema (disparadores
--       trg_app_setting_param_range_*): comparator.timeout_seconds de 1 a 30, price.stale_after_days de 1 o más
--       y el día del mes de 1 a 31; de formato: día de la semana de 1 a 7 (ISO 8601) y hora «HH:MM» de 00:00 a
--       23:59 (trg_app_setting_param_time_*). La regla ENTRE dos claves de Q2-77 (price.zero_weight_after_days
--       mayor o igual que price.stale_after_days) no va en el esquema: la valida el motor de precios
--       (backend/internal/precios/types.go, Settings.Validate, invalid_argument) y la API al guardar (RN-12).
--   Q2-67 opción 1 (K46): columna credit.vat_included (NULL = falta capturar); vat_on_interest_bp conserva su
--       significado (el IVA que se suma) y solo lleva valor cuando vat_included = 0: el CHECK rechaza un
--       porcentaje con vat_included 1 y también con vat_included NULL. Acepta vat_included = 0 sin porcentaje
--       («no incluye IVA» con el porcentaje por escribir; U4 captura las dos cosas por separado): falta el
--       porcentaje y es PROVISIONAL (bloqueo abierto S3-iva-lectura, pregunta 2; RN-26 1.4).
--   Sin DDL: Q2-38 (a) (el CHECK interest_cents >= 0 ya admite capital = monto, RN-07), Q2-66 B (qué pasa en
--       un mes corto es regla de S7, RN-29), Q2-68 C (expense_share ya congela las partes de cada gasto; el
--       alcance es un parámetro de C1 y del motor de reparto, RN-03) y el resto de PUERTA-2 §3 y §5
--       (lista en docs/estado/p2c/S1-puerta2.md).
-- Contrato: docs/contratos/diccionario-datos.md 1.4 y docs/contratos/reglas-negocio.md 1.4 (RN-03, RN-07,
-- RN-11, RN-12, RN-26 y RN-29).
-- Convenciones: las de 00001 a 00003. Notas técnicas:
--   * Las horas (tipo «hora HH:MM» del inventario) son texto JSON «HH:MM» de 24 h en America/Mexico_City ('"21:30"'); los días y los segundos,
--     enteros JSON. Cada clave nueva toma, si existe y es válido, el valor que su aviso tenga en params_json
--     al migrar (en una base nueva, 7, 7 y 1, los de la decisión F); si falta o no es válido, el de F.
--   * Los params_json de weekly_summary, meal_plan_missing y monthly_report quedan sin weekday ni day_of_month
--     y conservan updated_at (la casa no cambió nada). send_time y advance_days nunca estuvieron en params_json:
--     el Down solo borra sus claves.
--   * credit.vat_included: los créditos que ya existan pasan con el significado que les daba el esquema 1.3:
--     vat_on_interest_bp con valor → 0 (se suma ese IVA); NULL → 1 («la tasa o la cotización ya incluye IVA»,
--     así lo definía 1.3; Kueski, decisión C). Los créditos nuevos nacen NULL = falta capturar. SQLite (3.37 o
--     posterior) prueba el CHECK de una columna nueva contra las filas que ya existen, y un crédito de 1.3 con
--     porcentaje quedaría en (NULL, porcentaje) durante el ALTER: el porcentaje se aparta en la tabla temporal
--     p2_credit_vat_13 y vuelve junto con vat_included = 0. El Down quita la columna; un crédito con
--     vat_included = 0 y sin porcentaje vuelve a leerse «incluye IVA» en 1.3. En 1.4 ese crédito no incluye IVA
--     y le falta el porcentaje: lo que hacen C1 y S3 con él es PROVISIONAL (bloqueo abierto S3-iva-lectura).
--   * El Down repone weekday y day_of_month al principio de params_json con el valor vigente de app_setting
--     (con los valores iniciales, el texto de 1.3 byte por byte), borra las claves que siembra esta migración
--     (en cualquier alcance), los cuatro disparadores y la columna, y deja el esquema 1.3 idéntico (lo prueba
--     backend/db/tools/s1_verifica.py, DOWN14).
--   * Anotaciones nuevas que leen backend/db/tools/gen_dict.py y s1_verifica.py:
--     -- @param_rango <clave que ya tiene @param> | <rango aprobado, se suma a la unidad> | <decisión que lo aprueba>
--        (aquí, comparator.timeout_seconds, price.stale_after_days y price.zero_weight_after_days; el @param de
--        la 00001 o la 00003 sigue siendo el único de la clave, como pide backend/internal/api/c1gen).
--     -- @param_sin_valor_cerrado <primer campo exacto de un @param_sin_valor anterior> | <dónde quedó>
--     -- @resuelve <bloqueo> | <respuesta>: las marcas «provisional» de migraciones anteriores que citan ese
--        bloqueo son historia (00001 a 00003 no se editan).
--     -- @col_note de una columna que ya tenía nota: la última nota sustituye a la anterior en el diccionario.
--   * Nombres: prefijo reports.* para la programación de los reportes y alerts.* para la de los avisos (como
--     alerts.payment_due_days_before); optimizer.* y ocr.* por su módulo. Ningún DDL de 00001 a 00003 cambia.
-- =====================================================================

-- @resuelve S1-parametros-avisos | Q2-11 B (K47), respuesta explícita de Axel
-- Inventario de parámetros operativos (K44) que agrega o precisa la puerta 2. Sección de Ajustes: la que U4b
-- propone (docs/design/U4-ajustes-b.md, ubicación de U4-12: reportes en «Reportes y frecuencias»; hora y
-- anticipación de las alertas y día del plan sin armar en «Notificaciones»), la de U4a para la ventana de
-- fechas del comparador y, para el OCR, PROPUESTA S1 (junto a la IA opcional; ocr.md la dejó «a decidir»).
-- @param reports.daily_send_time | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "21:30" | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004
-- @param reports.weekly_weekday | entero | día de la semana ISO 8601, de 1 a 7 (7 = domingo) | 7 | Decisión F, lugar K47 | Reportes y frecuencias | 00004
-- @param reports.weekly_send_time | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "19:00" | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004
-- @param reports.monthly_day_of_month | entero | día del mes, de 1 a 31 (en un mes más corto sale su último día, Q2-66 B) | 1 | Decisión F, Q2-66 B, lugar K47 | Reportes y frecuencias | 00004
-- @param reports.monthly_send_time | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "08:00" | Q2-72 (K50), K31 punto 1, lugar K47 | Reportes y frecuencias | 00004
-- @param alerts.send_time | hora HH:MM | hora del centro de México, de 00:00 a 23:59 | "09:00" | Q2-73 (K50), K31 punto 2, K43, lugar K47 | Notificaciones | 00004
-- @param alerts.stale_prices_advance_days | entero | días antes del día de despensa, 0 o más | 1 | Q2-73 (K50), K31 punto 2, K43, lugar K47 | Notificaciones | 00004
-- @param alerts.meal_plan_missing_weekday | entero | día de la semana ISO 8601, de 1 a 7 (7 = domingo) | 7 | Decisión F, lugar K47 | Notificaciones | 00004
-- @param optimizer.date_window_days | entero | días antes y después de la fecha (fecha ± días), 0 o más | 7 | Q2-75 (K50) | Créditos y condiciones | 00004
-- @param ocr.timeout_seconds | entero | segundos por imagen, 0 o más | 30 | Q2-76 (K50); sección: propuesta S1 | Datos y respaldos | 00004
-- @param_rango comparator.timeout_seconds | de 1 a 30 (límites de seguridad, Q2-74) | Q2-74 (K50)
-- @param_rango price.stale_after_days | 1 o más (Q2-77) | Q2-77 (K50)
-- @param_rango price.zero_weight_after_days | igual o mayor que price.stale_after_days (Q2-77; lo valida el motor de precios) | Q2-77 (K50)
-- @param_sin_valor_cerrado optimizer.date_window_days | clave sembrada con 7 (Q2-75, K50) por la 00004
-- @param_sin_valor_cerrado notification.params_json send_time | claves reports.daily_send_time, reports.weekly_send_time, reports.monthly_send_time y alerts.send_time (Q2-72, Q2-73, K47) de la 00004
-- @param_sin_valor_cerrado notification.params_json advance_days | clave alerts.stale_prices_advance_days = 1 (Q2-73, K47) de la 00004
-- @param_sin_valor_cerrado notification.params_json weekday (weekly_summary y meal_plan_missing) | claves reports.weekly_weekday y alerts.meal_plan_missing_weekday = 7 (decisión F, K47) de la 00004
-- @param_sin_valor_cerrado notification.params_json day_of_month (monthly_report) | clave reports.monthly_day_of_month = 1 (decisión F, Q2-66 B, K47) de la 00004
-- @col_note notification.params_json | Desde 1.3 (00003, K44) `payment_due` ya no lleva `days_before` ni `budget_threshold` lleva `thresholds_pct` (son `alerts.payment_due_days_before`, `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `app_setting`), y desde 1.4 (00004, K47) `weekly_summary` y `meal_plan_missing` ya no llevan `weekday` ni `monthly_report` lleva `day_of_month` (son `reports.weekly_weekday`, `alerts.meal_plan_missing_weekday` y `reports.monthly_day_of_month`); la hora y la anticipación nunca estuvieron aquí (`reports.*_send_time`, `alerts.send_time` y `alerts.stale_prices_advance_days`). Con los valores iniciales solo `budget_threshold` conserva un parámetro (`scopes`); los demás quedan con `{}`. Los Down de la 00004 y de la 00003 los reponen aquí (RN-29)
-- @col_note credit.vat_on_interest_bp | Desde 1.4 (00004, Q2-67 opción 1) si la tasa ya incluye IVA lo dice `credit.vat_included`: este porcentaje conserva su significado y solo lleva valor cuando `vat_included` = 0; con `vat_included` = 1 es NULL, y con `vat_included` NULL el dato falta por capturar (RN-26)

-- +goose Up

-- ---------------------------------------------------------------------
-- Q2-67 opción 1 (K18). ¿La tasa incluye IVA? con «falta capturar»
-- ---------------------------------------------------------------------
-- @alter credit
-- @used_new Ajustes «Créditos y condiciones» (P-26.05, «¿la tasa incluye IVA?»). S3 (`Credit.missing_terms` lista `vat` mientras sea NULL en `bank_card` y `store_credit`: 5 datos de K18; con 0 y sin porcentaje también lo lista, PROVISIONAL por el bloqueo abierto `S3-iva-lectura`, pregunta 2, RN-26). C1 v1.3. U4.
-- Los créditos que ya existen conservan lo que el esquema 1.3 decía de ellos; updated_at no cambia (va en dos
-- pasos porque trg_credit_updated_at lo renueva si el UPDATE no lo trae). El porcentaje de 1.3 se aparta antes
-- del ALTER, porque el CHECK de la columna nueva rechazaría (NULL, porcentaje), y vuelve con vat_included = 0.
CREATE TEMP TABLE p2_credit_vat_13 AS SELECT id, vat_on_interest_bp FROM credit WHERE vat_on_interest_bp IS NOT NULL;
UPDATE credit SET updated_at = updated_at || '~', vat_on_interest_bp = NULL;
ALTER TABLE credit ADD COLUMN vat_included INTEGER CHECK (vat_included IN (0, 1) AND (vat_on_interest_bp IS NULL OR vat_included IS 0)); -- Nueva en 1.4 (Q2-67 opción 1, K18). ¿La tasa anual o la cotización ya incluye IVA? 1 = sí, no se suma IVA aparte; 0 = no, se suma el de vat_on_interest_bp; NULL = falta capturar (S3 lo lista como `vat` en missing_terms de `bank_card` y `store_credit`). Los créditos que ya existían al migrar pasan con el significado de 1.3 (porcentaje con valor → 0; NULL → 1). Con 0 y vat_on_interest_bp NULL falta el porcentaje: PROVISIONAL (bloqueo abierto S3-iva-lectura, pregunta 2) `vat` también va en missing_terms y S3 dice «falta capturar»; no calcula con 0 % ni inventa un porcentaje (K18) | booleano o NULL; con 1 o NULL, vat_on_interest_bp es NULL (solo lleva valor con 0); con 0 el CHECK admite el porcentaje en NULL
UPDATE credit SET updated_at = substr(updated_at, 1, length(updated_at) - 1),
    vat_included = CASE WHEN id IN (SELECT id FROM temp.p2_credit_vat_13) THEN 0 ELSE 1 END,
    vat_on_interest_bp = (SELECT s.vat_on_interest_bp FROM temp.p2_credit_vat_13 s WHERE s.id = credit.id)
WHERE updated_at LIKE '%~';
DROP TABLE temp.p2_credit_vat_13;

-- ---------------------------------------------------------------------
-- Q2-74, Q2-77, Q2-66 B y formatos. Rangos de una sola clave
-- ---------------------------------------------------------------------
-- Enteros con su rango: 1 a 30 (comparador, Q2-74), 1 o más (precio viejo, Q2-77), 1 a 7 (día de la semana,
-- ISO 8601), 1 a 31 (día del mes, Q2-66 B) y 0 o más (anticipación, ventana de fechas y OCR: sin rango
-- aprobado, solo el tipo y el signo, como los de la 00003).
-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_range_ins BEFORE INSERT ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('comparator.timeout_seconds', 'price.stale_after_days', 'reports.weekly_weekday', 'alerts.meal_plan_missing_weekday', 'reports.monthly_day_of_month', 'alerts.stale_prices_advance_days', 'optimizer.date_window_days', 'ocr.timeout_seconds')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'integer'
             AND CAST(NEW.value_json AS INTEGER) >= (CASE WHEN NEW.setting_key IN ('alerts.stale_prices_advance_days', 'optimizer.date_window_days', 'ocr.timeout_seconds') THEN 0 ELSE 1 END)
             AND CAST(NEW.value_json AS INTEGER) <= (CASE NEW.setting_key WHEN 'comparator.timeout_seconds' THEN 30 WHEN 'reports.weekly_weekday' THEN 7 WHEN 'alerts.meal_plan_missing_weekday' THEN 7 WHEN 'reports.monthly_day_of_month' THEN 31 ELSE 9223372036854775807 END))
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro es un entero dentro de su rango aprobado (puerta 2, RN-29)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_range_upd BEFORE UPDATE OF setting_key, value_json ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('comparator.timeout_seconds', 'price.stale_after_days', 'reports.weekly_weekday', 'alerts.meal_plan_missing_weekday', 'reports.monthly_day_of_month', 'alerts.stale_prices_advance_days', 'optimizer.date_window_days', 'ocr.timeout_seconds')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'integer'
             AND CAST(NEW.value_json AS INTEGER) >= (CASE WHEN NEW.setting_key IN ('alerts.stale_prices_advance_days', 'optimizer.date_window_days', 'ocr.timeout_seconds') THEN 0 ELSE 1 END)
             AND CAST(NEW.value_json AS INTEGER) <= (CASE NEW.setting_key WHEN 'comparator.timeout_seconds' THEN 30 WHEN 'reports.weekly_weekday' THEN 7 WHEN 'alerts.meal_plan_missing_weekday' THEN 7 WHEN 'reports.monthly_day_of_month' THEN 31 ELSE 9223372036854775807 END))
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro es un entero dentro de su rango aprobado (puerta 2, RN-29)');
END;
-- +goose StatementEnd

-- Horas: texto JSON «HH:MM» de 24 h, de 00:00 a 23:59, hora del centro de México.
-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_time_ins BEFORE INSERT ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('reports.daily_send_time', 'reports.weekly_send_time', 'reports.monthly_send_time', 'alerts.send_time')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'text'
             AND json_extract(NEW.value_json, '$') GLOB '[0-2][0-9]:[0-5][0-9]' AND json_extract(NEW.value_json, '$') <= '23:59')
BEGIN
    SELECT RAISE(ABORT, 'app_setting: esta hora es un texto HH:MM de 00:00 a 23:59, hora del centro de México (puerta 2, RN-29)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_time_upd BEFORE UPDATE OF setting_key, value_json ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('reports.daily_send_time', 'reports.weekly_send_time', 'reports.monthly_send_time', 'alerts.send_time')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'text'
             AND json_extract(NEW.value_json, '$') GLOB '[0-2][0-9]:[0-5][0-9]' AND json_extract(NEW.value_json, '$') <= '23:59')
BEGIN
    SELECT RAISE(ABORT, 'app_setting: esta hora es un texto HH:MM de 00:00 a 23:59, hora del centro de México (puerta 2, RN-29)');
END;
-- +goose StatementEnd

-- ---------------------------------------------------------------------
-- Q2-11 B (K47), Q2-72, Q2-73, Q2-75 y Q2-76 (K50). Claves nuevas con su valor aprobado
-- ---------------------------------------------------------------------
INSERT INTO app_setting (setting_key, value_json, description) VALUES
    ('reports.daily_send_time', '"21:30"', 'Hora a la que sale el resumen diario. Texto HH:MM de 24 h, hora del centro de México (Q2-72, K31, K47).'),
    ('reports.weekly_weekday',
     coalesce((SELECT CAST(json_extract(n.params_json, '$.weekday') AS TEXT) FROM notification n
               WHERE n.code = 'weekly_summary' AND json_type(n.params_json, '$.weekday') = 'integer'
                 AND json_extract(n.params_json, '$.weekday') BETWEEN 1 AND 7), '7'),
     'Día de la semana del resumen semanal: 1 = lunes … 7 = domingo (ISO 8601). Entero de 1 a 7 (decisión F, K47).'),
    ('reports.weekly_send_time', '"19:00"', 'Hora a la que sale el resumen semanal. Texto HH:MM de 24 h, hora del centro de México (Q2-72, K31, K47).'),
    ('reports.monthly_day_of_month',
     coalesce((SELECT CAST(json_extract(n.params_json, '$.day_of_month') AS TEXT) FROM notification n
               WHERE n.code = 'monthly_report' AND json_type(n.params_json, '$.day_of_month') = 'integer'
                 AND json_extract(n.params_json, '$.day_of_month') BETWEEN 1 AND 31), '1'),
     'Día del mes en que se genera el reporte mensual en PDF. Entero de 1 a 31; en un mes que no tiene ese día sale su último día (decisión F, Q2-66, K47).'),
    ('reports.monthly_send_time', '"08:00"', 'Hora a la que se genera el reporte mensual en PDF. Texto HH:MM de 24 h, hora del centro de México (Q2-72, K31, K47).'),
    ('alerts.send_time', '"09:00"', 'Hora a la que salen los avisos con hora. Texto HH:MM de 24 h, hora del centro de México (Q2-73, K31, K43, K47).'),
    ('alerts.stale_prices_advance_days', '1', 'Días antes del día de despensa en que llega el aviso de precios viejos: 0 es el mismo día. Entero en días (Q2-73, K31, K43, K47).'),
    ('alerts.meal_plan_missing_weekday',
     coalesce((SELECT CAST(json_extract(n.params_json, '$.weekday') AS TEXT) FROM notification n
               WHERE n.code = 'meal_plan_missing' AND json_type(n.params_json, '$.weekday') = 'integer'
                 AND json_extract(n.params_json, '$.weekday') BETWEEN 1 AND 7), '7'),
     'Día de la semana del aviso de plan de comidas sin armar: 1 = lunes … 7 = domingo (ISO 8601). Entero de 1 a 7 (decisión F, K47).'),
    ('optimizer.date_window_days', '7', 'Días antes y después de la fecha de una compra planeada que prueba el comparador de escenarios («fecha ± días»). Entero en días (Q2-75, K44).'),
    ('ocr.timeout_seconds', '30', 'Tiempo máximo de la lectura OCR de una imagen de ticket. Entero en segundos (Q2-76, K44).');

-- weekday y day_of_month salen de params_json para no dejar dos fuentes (K47). updated_at se conserva.
UPDATE notification SET updated_at = updated_at || '~',
    params_json = CASE WHEN params_json = '{"weekday": 7}' THEN '{}' ELSE json_remove(params_json, '$.weekday') END
WHERE code IN ('weekly_summary', 'meal_plan_missing') AND json_type(params_json, '$.weekday') IS NOT NULL;
UPDATE notification SET updated_at = updated_at || '~',
    params_json = CASE WHEN params_json = '{"day_of_month": 1}' THEN '{}' ELSE json_remove(params_json, '$.day_of_month') END
WHERE code = 'monthly_report' AND json_type(params_json, '$.day_of_month') IS NOT NULL;
UPDATE notification SET updated_at = substr(updated_at, 1, length(updated_at) - 1)
WHERE code IN ('weekly_summary', 'meal_plan_missing', 'monthly_report') AND updated_at LIKE '%~';

-- +goose Down
-- Deshace la 00004 en orden inverso y deja el esquema 1.3 idéntico.
-- Repone weekday y day_of_month al principio de params_json con el valor vigente de app_setting (con los
-- valores iniciales, el texto de 1.3 byte por byte) y conserva updated_at.
UPDATE notification SET updated_at = updated_at || '~',
    params_json = (SELECT CASE WHEN (SELECT count(*) FROM json_each(notification.params_json)) = 0 THEN '{"weekday": ' || v.w || '}'
                               WHEN substr(notification.params_json, 1, 1) = '{' THEN '{"weekday": ' || v.w || ', ' || substr(notification.params_json, 2)
                               ELSE json_set(notification.params_json, '$.weekday', json(v.w)) END
                   FROM (SELECT json(coalesce((SELECT value_json FROM app_setting WHERE setting_key = 'reports.weekly_weekday' AND user_account_id IS NULL), '7')) AS w) v)
WHERE code = 'weekly_summary' AND json_type(params_json) = 'object' AND json_type(params_json, '$.weekday') IS NULL;
UPDATE notification SET updated_at = updated_at || '~',
    params_json = (SELECT CASE WHEN (SELECT count(*) FROM json_each(notification.params_json)) = 0 THEN '{"weekday": ' || v.w || '}'
                               WHEN substr(notification.params_json, 1, 1) = '{' THEN '{"weekday": ' || v.w || ', ' || substr(notification.params_json, 2)
                               ELSE json_set(notification.params_json, '$.weekday', json(v.w)) END
                   FROM (SELECT json(coalesce((SELECT value_json FROM app_setting WHERE setting_key = 'alerts.meal_plan_missing_weekday' AND user_account_id IS NULL), '7')) AS w) v)
WHERE code = 'meal_plan_missing' AND json_type(params_json) = 'object' AND json_type(params_json, '$.weekday') IS NULL;
UPDATE notification SET updated_at = updated_at || '~',
    params_json = (SELECT CASE WHEN (SELECT count(*) FROM json_each(notification.params_json)) = 0 THEN '{"day_of_month": ' || v.d || '}'
                               WHEN substr(notification.params_json, 1, 1) = '{' THEN '{"day_of_month": ' || v.d || ', ' || substr(notification.params_json, 2)
                               ELSE json_set(notification.params_json, '$.day_of_month', json(v.d)) END
                   FROM (SELECT json(coalesce((SELECT value_json FROM app_setting WHERE setting_key = 'reports.monthly_day_of_month' AND user_account_id IS NULL), '1')) AS d) v)
WHERE code = 'monthly_report' AND json_type(params_json) = 'object' AND json_type(params_json, '$.day_of_month') IS NULL;
UPDATE notification SET updated_at = substr(updated_at, 1, length(updated_at) - 1)
WHERE code IN ('weekly_summary', 'meal_plan_missing', 'monthly_report') AND updated_at LIKE '%~';
DELETE FROM app_setting WHERE setting_key IN ('reports.daily_send_time', 'reports.weekly_weekday', 'reports.weekly_send_time', 'reports.monthly_day_of_month', 'reports.monthly_send_time', 'alerts.send_time', 'alerts.stale_prices_advance_days', 'alerts.meal_plan_missing_weekday', 'optimizer.date_window_days', 'ocr.timeout_seconds');
DROP TRIGGER IF EXISTS trg_app_setting_param_time_upd;
DROP TRIGGER IF EXISTS trg_app_setting_param_time_ins;
DROP TRIGGER IF EXISTS trg_app_setting_param_range_upd;
DROP TRIGGER IF EXISTS trg_app_setting_param_range_ins;
ALTER TABLE credit DROP COLUMN vat_included;
