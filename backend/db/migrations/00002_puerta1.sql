-- =====================================================================
-- Finanzas Hogar · Migración 00002 · Puerta 1: esquema 1.1 → 1.2
-- Agente S1-esquema-1.2 (flujo P1 «Aplicar respuestas de la puerta 1»), 11-oct-2026.
-- Aplica las respuestas de Axel de docs/plan/decisiones.md §K que tocan el esquema:
--   K3 borrador de correcciones del OCR · K4 chats de Telegram por confirmar
--   K9 captura sin conexión (clave de idempotencia y hora de captura)
--   K12 decimales capturados de las tasas · K18 tarjeta bancaria y crédito de tienda
--   K23 objetivo de peso · K6, K11 y K21 ajustes sembrados.
--   K8 la lista de compras empieza en lunes (disparadores; revisión S1-r1, hallazgo H-02).
--   K16: las restricciones duras del comparador (optimizer.*) NO se siembran.
-- Contrato: docs/contratos/diccionario-datos.md 1.2 y docs/contratos/reglas-negocio.md 1.2.
-- Convenciones: las de 00001. Notas técnicas:
--   * Las columnas nuevas entran con ALTER TABLE ... ADD COLUMN. Ahí SQLite no admite
--     UNIQUE, ni NOT NULL sin valor por defecto, ni un DEFAULT con expresión: por eso
--     admiten NULL (cada comentario dice qué significa NULL) y la unicidad va en índices
--     únicos. Cada CHECK nuevo es de la propia columna, para que el Down pueda quitarla.
--   * El Down borra en orden inverso los ajustes sembrados, los índices, la tabla nueva y
--     las columnas (DROP COLUMN, SQLite 3.35 o mayor) y deja el esquema 1.1 idéntico: el
--     texto de sqlite_master vuelve byte por byte (lo prueba backend/db/tools/s1_verifica.py).
--   * Anotaciones que lee backend/db/tools/gen_dict.py además de las de 00001:
--     -- @alter <tabla> (las columnas que siguen son de esa tabla; «descripción | regla»
--     en el comentario de la misma línea), -- @used_new <texto> (quién usa las columnas
--     nuevas) y -- @after <tabla> (lugar de una tabla nueva en el diccionario).
-- =====================================================================

-- +goose Up

-- ---------------------------------------------------------------------
-- K3. Ticket con OCR: «Guardar para después» guarda las correcciones como borrador
-- ---------------------------------------------------------------------
-- @alter attachment
-- @used_new C1 v1.1 (guardar y leer el borrador del ticket). S5 (lectura de tickets). U1 P-05 (N-071, N-075 y N-076). Regla RN-23.
ALTER TABLE attachment ADD COLUMN ocr_draft_json TEXT CHECK (ocr_draft_json IS NULL OR (json_valid(ocr_draft_json) AND json_type(ocr_draft_json) = 'object')); -- Nueva en 1.2 (K3). Borrador de las correcciones del ticket que la persona guardó con «Guardar para después», en JSON con la forma que define C1 v1.1. La lectura cruda (ocr_text) no se toca. Se vacía al confirmar el ticket | objeto JSON válido; NULL = sin borrador
ALTER TABLE attachment ADD COLUMN ocr_draft_saved_at TEXT CHECK ((ocr_draft_saved_at IS NULL) = (ocr_draft_json IS NULL) AND (ocr_draft_saved_at IS NULL OR (ocr_draft_saved_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(ocr_draft_saved_at) IS NOT NULL))); -- Nueva en 1.2 (K3). Cuándo se guardó el borrador por última vez (UTC) | instante ISO 8601 con «Z»; va junto con ocr_draft_json

-- ---------------------------------------------------------------------
-- K4. Chats de Telegram que el bot detecta y la persona confirma en Ajustes
-- ---------------------------------------------------------------------
-- @table telegram_chat_candidate
-- @after notification_recipient
-- @title Chat de Telegram por confirmar
-- @entity (complementaria) chats que el bot detecta cuando alguien le escribe y que la persona confirma en Ajustes (K4)
-- @desc Chat que el bot detectó cuando alguien le escribió y que espera confirmación en Ajustes «Notificaciones». Guarda solo el chat_id y el nombre que da Telegram. Al confirmarlo se liga a la persona y la API crea su notification_recipient en la misma transacción; descartarlo borra el renglón.
-- @used Ajustes «Notificaciones» (candidatos pendientes y confirmación, C1 v1.1). S7 (detección cuando la persona escribe al bot, K4). U1 P-26.09: N-299 (chats detectados) y N-316 (Confirmar chat detectado).
CREATE TABLE telegram_chat_candidate (
    id INTEGER PRIMARY KEY,
    chat_id TEXT NOT NULL UNIQUE CHECK (length(chat_id) BETWEEN 1 AND 32 AND chat_id NOT GLOB '*[^0-9-]*'), -- chat_id de Telegram del que llegó el mensaje | dígitos y signo menos
    telegram_name TEXT NOT NULL CHECK (length(trim(telegram_name)) > 0), -- Nombre que da Telegram para ese chat | no vacío
    detected_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')) CHECK (detected_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(detected_at) IS NOT NULL), -- Cuándo lo detectó el bot (UTC) | instante ISO 8601 con «Z»
    confirmed_at TEXT CHECK (confirmed_at IS NULL OR (confirmed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(confirmed_at) IS NOT NULL)), -- Cuándo la persona lo confirmó en Ajustes (UTC). NULL = pendiente | instante ISO 8601 con «Z»
    person_id INTEGER REFERENCES person (id) ON DELETE RESTRICT, -- Persona a quien se ligó al confirmar | obligatoria si está confirmado
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((confirmed_at IS NULL) = (person_id IS NULL)) -- @rule Un chat confirmado dice de quién es; uno pendiente todavía no tiene persona.
) STRICT;
-- @idx Candidatos pendientes de confirmar.
CREATE INDEX ix_telegram_chat_candidate_pending ON telegram_chat_candidate (detected_at) WHERE confirmed_at IS NULL;
-- @idx Llave foránea hacia la persona.
CREATE INDEX ix_telegram_chat_candidate_person ON telegram_chat_candidate (person_id);

-- +goose StatementBegin
CREATE TRIGGER trg_telegram_chat_candidate_updated_at AFTER UPDATE ON telegram_chat_candidate
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE telegram_chat_candidate SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- ---------------------------------------------------------------------
-- K8. La lista de compras cubre una semana de lunes a domingo (revisión S1-r1, hallazgo H-02)
-- ---------------------------------------------------------------------
-- shopping_list.start_date debe caer en lunes. SQLite no deja agregar un CHECK a una tabla que ya existe
-- sin reconstruirla, así que lo exigen dos disparadores: uno al dar de alta y otro al cambiar la fecha.
-- +goose StatementBegin
CREATE TRIGGER trg_shopping_list_monday_ins BEFORE INSERT ON shopping_list
FOR EACH ROW WHEN strftime('%w', NEW.start_date) IS NOT '1'
BEGIN
    SELECT RAISE(ABORT, 'shopping_list: la semana de la lista empieza en lunes (K8)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_shopping_list_monday_upd BEFORE UPDATE OF start_date ON shopping_list
FOR EACH ROW WHEN strftime('%w', NEW.start_date) IS NOT '1'
BEGIN
    SELECT RAISE(ABORT, 'shopping_list: la semana de la lista empieza en lunes (K8)');
END;
-- +goose StatementEnd

-- ---------------------------------------------------------------------
-- K9. Captura sin conexión: cola del teléfono para gastos y peso, sin duplicados
-- ---------------------------------------------------------------------
-- @alter expense
-- @used_new C1 v1.1 (capturas atrasadas con Idempotency-Key y fecha de captura). S6 (busca la clave antes de insertar). PWA en WF-3 (cola del teléfono, estado «Sin conexión»). Regla RN-24.
ALTER TABLE expense ADD COLUMN idempotency_key TEXT CHECK (idempotency_key IS NULL OR (length(idempotency_key) BETWEEN 16 AND 128 AND idempotency_key NOT GLOB '*[^A-Za-z0-9_-]*')); -- Nueva en 1.2 (K9). Clave que el teléfono genera para cada captura de su cola, igual al encabezado Idempotency-Key. Si llega otra vez, la API devuelve este gasto y no crea otro. A diferencia de la tabla idempotency_key (24 h), no se purga | 16 a 128 caracteres: letras sin acento, dígitos, guion y guion bajo; única (ux_expense_idempotency_key); NULL = sin clave
ALTER TABLE expense ADD COLUMN captured_at TEXT CHECK (captured_at IS NULL OR (captured_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(captured_at) IS NOT NULL)); -- Nueva en 1.2 (K9). Instante en que la persona capturó el gasto en el teléfono (UTC). En una captura atrasada es anterior a created_at. NULL = captura en línea: la hora de captura es created_at | instante ISO 8601 con «Z»
-- @idx Una captura de la cola no se duplica: un solo gasto por clave.
CREATE UNIQUE INDEX ux_expense_idempotency_key ON expense (idempotency_key);

-- @alter weight_entry
-- @used_new C1 v1.1 (pesos capturados sin conexión). S6 (busca la clave antes de insertar). PWA en WF-3 (cola del teléfono). Regla RN-24.
ALTER TABLE weight_entry ADD COLUMN idempotency_key TEXT CHECK (idempotency_key IS NULL OR (length(idempotency_key) BETWEEN 16 AND 128 AND idempotency_key NOT GLOB '*[^A-Za-z0-9_-]*')); -- Nueva en 1.2 (K9). Clave de la captura en la cola del teléfono. Si llega otra vez, la API devuelve este peso y no crea otro. No se purga | 16 a 128 caracteres: letras sin acento, dígitos, guion y guion bajo; única (ux_weight_entry_idempotency_key); NULL = sin clave
ALTER TABLE weight_entry ADD COLUMN captured_at TEXT CHECK (captured_at IS NULL OR (captured_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(captured_at) IS NOT NULL)); -- Nueva en 1.2 (K9). Instante en que la persona capturó el peso en el teléfono (UTC). NULL = captura en línea: la hora de captura es created_at | instante ISO 8601 con «Z»
-- @idx Una captura de la cola no se duplica: un solo peso por clave.
CREATE UNIQUE INDEX ux_weight_entry_idempotency_key ON weight_entry (idempotency_key);

-- ---------------------------------------------------------------------
-- K18 y K12. Tarjeta bancaria y crédito de tienda; tasas como se capturaron
-- ---------------------------------------------------------------------
-- @alter credit
-- @used_new Ajustes «Créditos y condiciones» (Axel los captura de su estado de cuenta, K18). S3 (fórmula general; sin datos muestra «falta capturar»). C1 v1.1. U4. Reglas RN-25 y RN-26.
ALTER TABLE credit ADD COLUMN rate_base TEXT CHECK (rate_base IS NULL OR (rate_base IN ('daily', 'monthly') AND kind IN ('bank_card', 'store_credit'))); -- Nueva en 1.2 (K18). Base con que el estado de cuenta aplica la tasa anual: daily = diaria, monthly = mensual. NULL = falta capturar | enumerado, solo bank_card y store_credit
ALTER TABLE credit ADD COLUMN min_payment_pct_bp INTEGER CHECK (min_payment_pct_bp IS NULL OR (min_payment_pct_bp BETWEEN 0 AND 10000 AND kind IN ('bank_card', 'store_credit'))); -- Nueva en 1.2 (K18). Porcentaje del pago mínimo que dice el estado de cuenta, en puntos base. NULL = falta capturar | 0 a 10000, solo bank_card y store_credit
ALTER TABLE credit ADD COLUMN min_payment_fixed_cents INTEGER CHECK (min_payment_fixed_cents IS NULL OR (min_payment_fixed_cents >= 0 AND kind IN ('bank_card', 'store_credit'))); -- Nueva en 1.2 (K18). Mínimo fijo del pago mínimo que dice el estado de cuenta. NULL = falta capturar | centavos, solo bank_card y store_credit
ALTER TABLE credit ADD COLUMN annual_rate_decimals INTEGER CHECK (annual_rate_decimals IS NULL OR (annual_rate_decimals >= 0 AND annual_rate_num IS NOT NULL)); -- Nueva en 1.2 (K12). Decimales con que se capturó el porcentaje de la tasa anual; la pantalla lo muestra así y los cálculos usan annual_rate_num/annual_rate_den. NULL = sin tasa o sin dato | entero mayor o igual a 0, solo con tasa anual

-- @alter credit_term_rate
-- @used_new C1 v1.1 (CreditTermRate muestra la tasa con los decimales capturados). U4 (tabla de tasas por plazo). Regla RN-25.
ALTER TABLE credit_term_rate ADD COLUMN rate_decimals INTEGER CHECK (rate_decimals >= 0); -- Nueva en 1.2 (K12). Decimales con que se capturó o se leyó de la cotización el porcentaje (17.27 % → 2; 25.8685 % → 4). La pantalla redondea la fracción a esos decimales con HALF_UP; los cálculos usan rate_num/rate_den. NULL = sin dato | entero mayor o igual a 0

-- @alter credit_purchase
-- @used_new C1 v1.1 (tasa de la compra como se capturó o cotizó). U4. Regla RN-25.
ALTER TABLE credit_purchase ADD COLUMN rate_decimals INTEGER CHECK (rate_decimals IS NULL OR (rate_decimals >= 0 AND rate_num IS NOT NULL)); -- Nueva en 1.2 (K12). Decimales con que se muestra la tasa de esta compra, los mismos con que se capturó o cotizó. NULL = sin tasa o sin dato | entero mayor o igual a 0, solo si la compra lleva tasa

-- @alter promotion
-- @used_new C1 v1.1 (tasa promocional como se capturó). U4. Regla RN-25.
ALTER TABLE promotion ADD COLUMN rate_decimals INTEGER CHECK (rate_decimals >= 0); -- Nueva en 1.2 (K12). Decimales con que se capturó la tasa promocional (0 % → 0). NULL = sin dato | entero mayor o igual a 0

-- ---------------------------------------------------------------------
-- K23. Meta de calorías: objetivo de peso de cada persona
-- ---------------------------------------------------------------------
-- @alter health_profile
-- @used_new Pantalla Salud y Ajustes «Salud y metas de ahorro». S4 (meta sugerida con Mifflin-St Jeor, actividad y objetivo; regla RN-17). U4.
ALTER TABLE health_profile ADD COLUMN goal TEXT CHECK (goal IN ('lose', 'maintain', 'gain')); -- Nueva en 1.2 (K23). Objetivo que ajusta la meta sugerida de kilocalorías: lose = bajar de peso, maintain = mantener, gain = subir. NULL = todavía no lo elige | enumerado

-- ---------------------------------------------------------------------
-- Ajustes del hogar decididos en la puerta 1 (K6, K11 y K21)
-- K16: optimizer.min_available_buffer_cents, optimizer.max_live_purchases y
-- optimizer.zero_debt_by_date NO se siembran: los captura Axel y, mientras falten,
-- el comparador avisa «faltan restricciones» y no recomienda.
-- ---------------------------------------------------------------------
INSERT INTO app_setting (setting_key, value_json, description) VALUES
    ('projection.repeat_cadence_enabled', 'false', 'Proyección de los meses 7 a 12: agrega la cadencia bimestral de despensa a crédito y el papel alterno de $199. Apagado por defecto: solo compromisos contratados, gastos fijos e ingresos del CSV (K21).'),
    ('session.max_age_days', '90', 'Días que dura la sesión de un dispositivo; se renueva con cada uso (K11).'),
    ('idempotency.retention_hours', '24', 'Horas que se guarda una clave de idempotency_key antes de purgarla (K6).');

-- +goose Down
-- Deshace la 00002 en orden inverso y deja el esquema 1.1 idéntico.
DELETE FROM app_setting WHERE user_account_id IS NULL AND setting_key IN ('projection.repeat_cadence_enabled', 'session.max_age_days', 'idempotency.retention_hours');
ALTER TABLE health_profile DROP COLUMN goal;
ALTER TABLE promotion DROP COLUMN rate_decimals;
ALTER TABLE credit_purchase DROP COLUMN rate_decimals;
ALTER TABLE credit_term_rate DROP COLUMN rate_decimals;
ALTER TABLE credit DROP COLUMN annual_rate_decimals;
ALTER TABLE credit DROP COLUMN min_payment_fixed_cents;
ALTER TABLE credit DROP COLUMN min_payment_pct_bp;
ALTER TABLE credit DROP COLUMN rate_base;
DROP INDEX IF EXISTS ux_weight_entry_idempotency_key;
ALTER TABLE weight_entry DROP COLUMN captured_at;
ALTER TABLE weight_entry DROP COLUMN idempotency_key;
DROP INDEX IF EXISTS ux_expense_idempotency_key;
ALTER TABLE expense DROP COLUMN captured_at;
ALTER TABLE expense DROP COLUMN idempotency_key;
DROP TRIGGER IF EXISTS trg_shopping_list_monday_upd;
DROP TRIGGER IF EXISTS trg_shopping_list_monday_ins;
DROP TABLE IF EXISTS telegram_chat_candidate;
ALTER TABLE attachment DROP COLUMN ocr_draft_saved_at;
ALTER TABLE attachment DROP COLUMN ocr_draft_json;
