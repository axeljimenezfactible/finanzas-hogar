-- =====================================================================
-- Finanzas Hogar · Migración 00005 · Tablero de control (K54): esquema 1.4 → 1.5
-- Agente s1-c1-v15 (flujo WF-3c «Construcción»), 12-oct-2026.
-- Aplica K54 (docs/plan/decisiones.md §K, 12-oct-2026): Inicio › Resumen pasa a ser un tablero de control con el
-- barómetro de capacidad de endeudamiento = pagos mensuales de todos los créditos ÷ ingreso mensual neto × 100
-- (decisión 1), con umbrales configurables: verde hasta 30 %, amarillo de 30 a 40 %, rojo arriba de 40 % (decisión 2,
-- K44). Las claves y sus valores iniciales salen de docs/design/U3-dashboard.md 2.0 §2.6 (agente u3-v2).
--   * Solo siembra cinco claves de app_setting del hogar (user_account_id NULL) y dos disparadores de rango. No crea ni
--     cambia tablas, columnas ni índices; ningún DDL de 00001 a 00004 cambia y ninguna migración aplicada se edita.
--   * Valores: debt.capacity_warning_pct = 30 y debt.capacity_limit_pct = 40 son FIRMES (K54 decisión 2).
--     dashboard.payment_reminder_count = 5, budget.signal_months = 3 y budget.signal_margin_pct = 10 son PROVISIONALES
--     (bloqueos abiertos U3-tablero-recordatorio y U3-tablero-senal-presupuesto, propuesta de U3 2.0).
--   * Rangos (los prueban los disparadores trg_app_setting_param_k54_ins y _upd): de 1 a 100 los dos umbrales, de 1 a 20
--     el número de pagos, de 2 a 12 los meses y de 1 a 50 el margen. Son la propuesta de U3 2.0 §2.6, PROVISIONAL
--     (bloqueo abierto C1-tablero, punto 2). La regla entre dos claves (debt.capacity_warning_pct menor que
--     debt.capacity_limit_pct) la valida la API con 422 validation_failed, no el esquema (como Q2-77 en la 00004).
--   * Sección de Ajustes: «Presupuesto» (P-26.02) para las cinco, PROVISIONAL (bloqueo abierto C1-tablero, punto 1):
--     U3 2.0 §2.6 deja el lugar a U4 y K44 no admite una sección nueva.
--   * Nombre de la clave del recordatorio: dashboard.payment_reminder_count, como lo escribe U3 2.0 §2.6 y el bloqueo
--     U3-tablero-recordatorio (la tarea del orquestador decía payment_reminders_count; C1-tablero, punto 4).
--   * El Down borra las cinco claves (en cualquier alcance) y los dos disparadores, y deja el esquema 1.4 idéntico.
-- =====================================================================

-- Inventario de parámetros operativos (K44) que agrega K54.
-- @param debt.capacity_warning_pct | entero | % del ingreso mensual neto que se va en pagos de crédito del mes, frontera entre verde y amarillo del barómetro, de 1 a 100 | 30 | K54 decisión 2 (valor FIRME); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005
-- @param debt.capacity_limit_pct | entero | % del ingreso mensual neto que se va en pagos de crédito del mes, tope del barómetro (arriba es rojo), de 1 a 100 | 40 | K54 decisión 2 (valor FIRME); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005
-- @param dashboard.payment_reminder_count | entero | pagos que muestra «Pagos próximos» del tablero, de 1 a 20 | 5 | K54; valor PROVISIONAL (bloqueo U3-tablero-recordatorio); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005
-- @param budget.signal_months | entero | meses cerrados seguidos que mira la sugerencia «subir o recortar presupuesto», de 2 a 12 | 3 | K54; valor PROVISIONAL (bloqueo U3-tablero-senal-presupuesto); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005
-- @param budget.signal_margin_pct | entero | puntos de % usado arriba o abajo del 100 % del plan para la sugerencia «subir o recortar presupuesto», de 1 a 50 | 10 | K54; valor PROVISIONAL (bloqueo U3-tablero-senal-presupuesto); rango y sección PROVISIONAL (bloqueo C1-tablero) | Presupuesto | 00005
-- @param_rango debt.capacity_limit_pct | mayor que debt.capacity_warning_pct (lo valida la API con 422 validation_failed) | K54 decisión 2; PROVISIONAL (bloqueo C1-tablero)
-- +goose Up

-- ---------------------------------------------------------------------
-- K54. Rangos de las cinco claves (entero JSON dentro de su rango; PROVISIONAL, bloqueo C1-tablero)
-- ---------------------------------------------------------------------
-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_k54_ins BEFORE INSERT ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('debt.capacity_warning_pct', 'debt.capacity_limit_pct', 'dashboard.payment_reminder_count', 'budget.signal_months', 'budget.signal_margin_pct')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'integer'
             AND CAST(NEW.value_json AS INTEGER) >= (CASE NEW.setting_key WHEN 'budget.signal_months' THEN 2 ELSE 1 END)
             AND CAST(NEW.value_json AS INTEGER) <= (CASE NEW.setting_key WHEN 'dashboard.payment_reminder_count' THEN 20 WHEN 'budget.signal_months' THEN 12 WHEN 'budget.signal_margin_pct' THEN 50 ELSE 100 END))
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro del tablero es un entero dentro de su rango (K54, RN-29)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_param_k54_upd BEFORE UPDATE OF setting_key, value_json ON app_setting
FOR EACH ROW WHEN NEW.setting_key IN ('debt.capacity_warning_pct', 'debt.capacity_limit_pct', 'dashboard.payment_reminder_count', 'budget.signal_months', 'budget.signal_margin_pct')
    AND NOT (json_valid(NEW.value_json) AND json_type(NEW.value_json) = 'integer'
             AND CAST(NEW.value_json AS INTEGER) >= (CASE NEW.setting_key WHEN 'budget.signal_months' THEN 2 ELSE 1 END)
             AND CAST(NEW.value_json AS INTEGER) <= (CASE NEW.setting_key WHEN 'dashboard.payment_reminder_count' THEN 20 WHEN 'budget.signal_months' THEN 12 WHEN 'budget.signal_margin_pct' THEN 50 ELSE 100 END))
BEGIN
    SELECT RAISE(ABORT, 'app_setting: este parámetro del tablero es un entero dentro de su rango (K54, RN-29)');
END;
-- +goose StatementEnd

-- ---------------------------------------------------------------------
-- K54. Claves nuevas con su valor inicial
-- ---------------------------------------------------------------------
INSERT INTO app_setting (setting_key, value_json, description) VALUES
    ('debt.capacity_warning_pct', '30', 'Frontera entre verde y amarillo del barómetro de capacidad de endeudamiento: % del ingreso mensual neto que se va en pagos de crédito del mes. Con el entero redondeado, hasta este valor es verde. Entero de 1 a 100, menor que debt.capacity_limit_pct (K54, RN-30).'),
    ('debt.capacity_limit_pct', '40', 'Tope del barómetro de capacidad de endeudamiento: arriba de este % del ingreso mensual neto es rojo; de debt.capacity_warning_pct a este valor es amarillo. Entero de 1 a 100 (K54, RN-30).'),
    ('dashboard.payment_reminder_count', '5', 'Cuántos pagos muestra «Pagos próximos» en el tablero de Inicio. Entero de 1 a 20 (K54; PROVISIONAL, bloqueo U3-tablero-recordatorio).'),
    ('budget.signal_months', '3', 'Meses cerrados seguidos que mira la sugerencia «subir o recortar presupuesto» de cada categoría. Entero de 2 a 12 (K54, RN-31; PROVISIONAL, bloqueo U3-tablero-senal-presupuesto).'),
    ('budget.signal_margin_pct', '10', 'Margen de la sugerencia «subir o recortar presupuesto»: subir si el % usado pasó de 100 + margen en todos los meses que mira; recortar si quedó abajo de 100 − margen. Entero de 1 a 50 (K54, RN-31; PROVISIONAL, bloqueo U3-tablero-senal-presupuesto).');

-- +goose Down
DELETE FROM app_setting WHERE setting_key IN ('debt.capacity_warning_pct', 'debt.capacity_limit_pct', 'dashboard.payment_reminder_count', 'budget.signal_months', 'budget.signal_margin_pct');
DROP TRIGGER IF EXISTS trg_app_setting_param_k54_upd;
DROP TRIGGER IF EXISTS trg_app_setting_param_k54_ins;
