-- =====================================================================
-- Finanzas Hogar · Migración 00001 · Esquema base (S1 «Dominio y modelo de datos», WF-0)
-- Contrato: docs/contratos/diccionario-datos.md y docs/contratos/reglas-negocio.md
-- Versión del esquema: 1.1. La migración 00002_puerta1.sql lo lleva a 1.2 (puerta 1, decisiones §K).
-- Convenciones (detalle en reglas-negocio.md, sección «Convenciones»):
--   * Dinero en enteros de centavos MXN con sufijo _cents. Sin REAL para dinero.
--   * Fechas de negocio TEXT ISO AAAA-MM-DD (calendario America/Mexico_City), sufijo _date.
--   * Instantes TEXT ISO 8601 en UTC con «Z», sufijo _at.
--   * Booleanos INTEGER 0/1. Porcentajes en puntos base (_bp, 10000 = 100 %).
--   * Tasas como fracción exacta num/den (nunca redondeadas a 4 decimales).
--   * Cantidades físicas en milésimas de la unidad base (_milli).
--   * Tablas STRICT. id INTEGER PRIMARY KEY. created_at y updated_at en todas.
--   * Sin borrado lógico general: catálogos se archivan con is_active = 0,
--     movimientos se borran de verdad y quedan en audit_log.
-- PRAGMA foreign_keys: goose envuelve esta migración en una transacción y dentro
-- de una transacción el PRAGMA no tiene efecto. El backend debe activarlo en cada
-- conexión (DSN de modernc: _pragma=foreign_keys(1)). Se deja aquí para cuando
-- el archivo se aplica con el cliente sqlite3.
-- =====================================================================

-- +goose Up
PRAGMA foreign_keys = ON;

-- ---------------------------------------------------------------------
-- A. Personas y acceso
-- ---------------------------------------------------------------------

-- @table person
-- @title Persona
-- @entity persona
-- @desc Persona que participa en el dinero del hogar: el titular (Axel), la co-titular (Ximena) o una participante externa sin cuenta (Karen).
-- @used Ajustes «Personas y reparto» y «Cuenta y roles». Motores S2 (reparto), S3 (partes por persona y tope por quincena), S4 (plan de comidas), S7 (estado de cuenta de Karen). Filtro «persona» de C1.
CREATE TABLE person (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL UNIQUE CHECK (length(code) BETWEEN 1 AND 32 AND code NOT GLOB '*[^a-z0-9_]*'), -- Clave estable para semilla e importadores (axel, ximena, karen) | solo a-z, 0-9 y guion bajo
    display_name TEXT NOT NULL CHECK (length(trim(display_name)) > 0), -- Nombre que se muestra en pantallas y reportes | no vacío
    participation TEXT NOT NULL CHECK (participation IN ('holder', 'co_holder', 'external')), -- Papel: holder = titular (Axel), co_holder = co-titular con los mismos permisos (Ximena), external = participante externa sin cuenta (Karen) | enumerado, un solo holder
    pays_share_directly INTEGER NOT NULL DEFAULT 0 CHECK (pays_share_directly IN (0, 1)), -- 1 = paga su parte directo al acreedor (Karen con Kueski). Su parte es informativa para el flujo de efectivo del titular y repone la línea cuando ella paga | booleano
    fortnight_payment_cap_cents INTEGER CHECK (fortnight_payment_cap_cents >= 0), -- Tope de pago por quincena de la persona, opcional. El comparador lo usa como objetivo (pago parejo y bajo el tope), nunca como restricción dura. No hay un tope fijo: Axel lo captura por persona en Ajustes «Personas y reparto» cuando quiera (K13) | centavos, NULL = sin tope
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Garantiza que haya un solo titular.
CREATE UNIQUE INDEX ux_person_single_holder ON person (participation) WHERE participation = 'holder';
-- @idx Garantiza que haya una sola co-titular: con el titular suman las dos cuentas posibles (decisión A).
CREATE UNIQUE INDEX ux_person_single_co_holder ON person (participation) WHERE participation = 'co_holder';

-- @table user_account
-- @title Cuenta de usuario
-- @entity cuenta_usuario
-- @desc Cuenta para iniciar sesión. Solo dos: Axel y Ximena, con los mismos permisos (no hay roles). Karen nunca tiene cuenta.
-- @used Ajustes «Cuenta y roles». S6 (autenticación argon2id, límite de intentos).
CREATE TABLE user_account (
    id INTEGER PRIMARY KEY,
    person_id INTEGER NOT NULL UNIQUE REFERENCES person (id) ON DELETE RESTRICT, -- Persona dueña de la cuenta | una cuenta por persona, nunca una persona externa (disparador)
    username TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(username) BETWEEN 3 AND 64), -- Usuario para iniciar sesión | único sin distinguir mayúsculas, 3 a 64 caracteres
    password_hash TEXT NOT NULL CHECK (password_hash LIKE '$argon2id$%'), -- Hash argon2id en formato PHC. La contraseña nunca se guarda | debe empezar con $argon2id$
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)), -- 0 = cuenta deshabilitada (no puede iniciar sesión) | booleano
    password_changed_at TEXT, -- Último cambio de contraseña (UTC)
    last_login_at TEXT, -- Último inicio de sesión correcto (UTC)
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;

-- @table device_session
-- @title Sesión por dispositivo
-- @entity sesion_dispositivo
-- @desc Sesión larga (meses) de una cuenta en un dispositivo. Se puede cerrar desde Ajustes.
-- @used Ajustes «Cuenta y roles» (cerrar sesiones). S6 (middleware de sesión, purga de vencidas).
CREATE TABLE device_session (
    id INTEGER PRIMARY KEY,
    user_account_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE, -- Cuenta dueña de la sesión
    token_hash TEXT NOT NULL UNIQUE CHECK (length(token_hash) = 64), -- SHA-256 en hexadecimal del token de sesión. El token en claro nunca se guarda | 64 caracteres
    device_label TEXT, -- Nombre legible del dispositivo (Chrome en Android)
    user_agent TEXT, -- Agente de usuario reportado por el navegador
    ip_address TEXT, -- Última IP de la red local desde la que se usó
    last_seen_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')), -- Último uso (UTC)
    expires_at TEXT NOT NULL, -- Vencimiento de la sesión (UTC) | posterior a created_at
    revoked_at TEXT, -- Cierre manual desde Ajustes (UTC). NULL = vigente
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (julianday(expires_at) > julianday(created_at)) -- @rule La sesión vence después de crearse.
) STRICT;
-- @idx Lista las sesiones vigentes de una cuenta (Ajustes, cerrar sesiones).
CREATE INDEX ix_device_session_user ON device_session (user_account_id, revoked_at);
-- @idx Purga de sesiones vencidas.
CREATE INDEX ix_device_session_expires ON device_session (expires_at);

-- @table login_attempt
-- @title Intento de acceso
-- @entity (complementaria) decisión A: límite de 5 intentos en 15 minutos
-- @desc Bitácora de intentos de inicio de sesión para bloquear tras 5 fallos en 15 minutos.
-- @used S6 (autenticación y bloqueo). Job de purga.
CREATE TABLE login_attempt (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL COLLATE NOCASE, -- Usuario tecleado (puede no existir)
    user_account_id INTEGER REFERENCES user_account (id) ON DELETE SET NULL, -- Cuenta que corresponde al usuario, si existe
    ip_address TEXT, -- IP de origen en la red local
    succeeded INTEGER NOT NULL CHECK (succeeded IN (0, 1)), -- 1 = acceso correcto | booleano
    attempted_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')), -- Momento del intento (UTC)
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Cuenta los fallos de un usuario en la ventana de 15 minutos.
CREATE INDEX ix_login_attempt_username ON login_attempt (username, attempted_at);
-- @idx Cuenta los fallos desde una IP en la ventana de 15 minutos.
CREATE INDEX ix_login_attempt_ip ON login_attempt (ip_address, attempted_at);
-- @idx Llave foránea hacia la cuenta.
CREATE INDEX ix_login_attempt_user ON login_attempt (user_account_id);

-- @table app_setting
-- @title Ajuste
-- @entity (complementaria) valores de Ajustes que no son de una entidad (pesos del optimizador, frescura, interruptores de IA)
-- @desc Ajustes del hogar o de una cuenta en formato JSON, con clave por sección.
-- @used Ajustes (varias secciones). S3 (pesos y restricciones del optimizador), S5 (frescura), S6 (límites de acceso, IA opcional).
CREATE TABLE app_setting (
    id INTEGER PRIMARY KEY,
    setting_key TEXT NOT NULL CHECK (setting_key GLOB '[a-z]*' AND setting_key NOT GLOB '*[^a-z0-9_.]*'), -- Clave con prefijo de sección (price.stale_after_days) | minúsculas, dígitos, punto y guion bajo
    user_account_id INTEGER REFERENCES user_account (id) ON DELETE CASCADE, -- NULL = ajuste del hogar. Con valor = preferencia de esa cuenta
    value_json TEXT NOT NULL CHECK (json_valid(value_json)), -- Valor en JSON | JSON válido
    description TEXT, -- Para qué sirve, en español
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Una sola fila por clave y alcance (hogar o cuenta).
CREATE UNIQUE INDEX ux_app_setting_key_scope ON app_setting (setting_key, coalesce(user_account_id, 0));
-- @idx Llave foránea hacia la cuenta.
CREATE INDEX ix_app_setting_user ON app_setting (user_account_id);

-- @table audit_log
-- @title Bitácora de cambios
-- @entity (complementaria) convención de auditoría de cambios
-- @desc Registro inmutable de cada alta, cambio o baja hecha por la API o por un job: quién, cuándo y el renglón antes y después.
-- @used S6 (lo escribe en la misma transacción de cada cambio). Ajustes «Datos y respaldos» (historial).
CREATE TABLE audit_log (
    id INTEGER PRIMARY KEY,
    occurred_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')), -- Momento del cambio (UTC)
    user_account_id INTEGER, -- Cuenta que hizo el cambio. NULL = proceso del sistema. Sin llave foránea para que la bitácora sobreviva a cualquier baja
    action TEXT NOT NULL CHECK (action IN ('insert', 'update', 'delete')), -- Tipo de cambio | enumerado
    table_name TEXT NOT NULL, -- Tabla afectada
    row_id INTEGER NOT NULL, -- id del renglón afectado
    before_json TEXT CHECK (before_json IS NULL OR json_valid(before_json)), -- Renglón antes del cambio (NULL en insert) | JSON válido
    after_json TEXT CHECK (after_json IS NULL OR json_valid(after_json)), -- Renglón después del cambio (NULL en delete) | JSON válido
    request_id TEXT, -- Identificador de la petición HTTP o del job que originó el cambio
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((action = 'insert' AND before_json IS NULL AND after_json IS NOT NULL) OR (action = 'update' AND before_json IS NOT NULL AND after_json IS NOT NULL) OR (action = 'delete' AND before_json IS NOT NULL AND after_json IS NULL)) -- @rule insert lleva solo «después», delete solo «antes», update ambos.
) STRICT;
-- @idx Historial de un renglón.
CREATE INDEX ix_audit_log_row ON audit_log (table_name, row_id, occurred_at);
-- @idx Consulta por fecha.
CREATE INDEX ix_audit_log_time ON audit_log (occurred_at);
-- @idx Cambios hechos por una cuenta.
CREATE INDEX ix_audit_log_user ON audit_log (user_account_id, occurred_at);

-- @table idempotency_key
-- @title Clave de idempotencia
-- @entity (complementaria) encabezado Idempotency-Key de los POST de captura (C1 §1 punto 8): una captura repetida no se duplica
-- @desc Clave que manda el cliente en cada POST de captura. Si la misma cuenta repite la clave con el mismo cuerpo, la API devuelve la respuesta guardada; con otro cuerpo responde 409 idempotency_key_reused. La retención es de 24 h (K6, app_setting idempotency.retention_hours de la migración 00002); la purga la hace un job. Las capturas de la cola del teléfono guardan además su clave en expense y weight_entry, sin purga (K9).
-- @used S6 (API: los POST de captura de C1, entre ellos «Registrar gasto» del flujo F-1 de U1) y jobs de S6 (purga por antigüedad).
CREATE TABLE idempotency_key (
    id INTEGER PRIMARY KEY,
    user_account_id INTEGER NOT NULL REFERENCES user_account (id) ON DELETE CASCADE, -- Cuenta que mandó la clave
    idempotency_key TEXT NOT NULL CHECK (length(idempotency_key) BETWEEN 16 AND 128 AND idempotency_key NOT GLOB '*[^A-Za-z0-9_-]*'), -- Valor del encabezado Idempotency-Key | 16 a 128 caracteres: letras sin acento, dígitos, guion y guion bajo (C1 §1 punto 8)
    http_method TEXT NOT NULL CHECK (http_method IN ('POST', 'PUT', 'PATCH', 'DELETE')), -- Método de la petición | enumerado
    request_path TEXT NOT NULL CHECK (request_path GLOB '/*'), -- Ruta de la petición (/api/v1/expenses) | empieza con /
    request_sha256 TEXT NOT NULL CHECK (length(request_sha256) = 64), -- Huella SHA-256 del cuerpo: distingue repetir la petición de reutilizar la clave con otro cuerpo | 64 caracteres
    response_status INTEGER CHECK (response_status BETWEEN 100 AND 599), -- Código HTTP de la respuesta guardada. NULL = la petición original sigue en proceso | 100 a 599
    response_body TEXT CHECK (response_body IS NULL OR json_valid(response_body)), -- Cuerpo JSON de la respuesta guardada, para repetirla igual | JSON válido
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (user_account_id, idempotency_key) -- @rule Una clave se usa una sola vez por cuenta.
) STRICT;
-- @idx Purga de claves por antigüedad (retención de 24 h, K6).
CREATE INDEX ix_idempotency_key_created ON idempotency_key (created_at);

-- ---------------------------------------------------------------------
-- B. Dinero, periodos y reparto
-- ---------------------------------------------------------------------

-- @table period
-- @title Periodo
-- @entity periodo
-- @desc Mes calendario (periodo base del presupuesto) y sus dos quincenas (subperiodos). Los días de pago (15, y 30 o último día) son FIRMES (oráculo). Los rangos de la quincena (1 a 15 y 16 a fin de mes) son FIRMES (K8). Ambos se validan en la propia tabla.
-- @used S2 (presupuesto mensual, flujo por quincena, proyección). Pantallas Presupuesto e Inicio. Filtro «periodo» de C1.
CREATE TABLE period (
    id INTEGER PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('month', 'fortnight')), -- month = mes calendario, fortnight = quincena | enumerado
    parent_id INTEGER REFERENCES period (id) ON DELETE CASCADE, -- Mes al que pertenece la quincena. NULL en meses
    year INTEGER NOT NULL CHECK (year BETWEEN 2000 AND 2100), -- Año | 2000 a 2100
    month INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12), -- Mes | 1 a 12
    half INTEGER CHECK (half IN (1, 2)), -- Quincena 1 (días 1 a 15) o 2 (día 16 a fin de mes). NULL en meses | 1 o 2. Rangos FIRMES (K8)
    start_date TEXT NOT NULL CHECK (date(start_date) IS start_date), -- Primer día (ISO) | fecha válida
    end_date TEXT NOT NULL CHECK (date(end_date) IS end_date), -- Último día, inclusive (ISO) | fecha válida
    pay_date TEXT CHECK (date(pay_date) IS pay_date), -- Día de pago de la quincena: 15, o 30 o el último día si el mes es más corto (28-feb-2027). NULL en meses | fecha válida
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ( -- @rule Mes: del día 1 al último. Quincena 1: del 1 al 15, paga el 15. Quincena 2: del 16 al último, paga el 30 o el último día si es menor (febrero). Rangos de quincena FIRMES (K8).
        (kind = 'month' AND half IS NULL AND parent_id IS NULL AND pay_date IS NULL
            AND start_date = printf('%04d-%02d-01', year, month)
            AND end_date = date(start_date, '+1 month', '-1 day'))
        OR (kind = 'fortnight' AND half = 1 AND parent_id IS NOT NULL
            AND start_date = printf('%04d-%02d-01', year, month)
            AND end_date = printf('%04d-%02d-15', year, month)
            AND pay_date = end_date)
        OR (kind = 'fortnight' AND half = 2 AND parent_id IS NOT NULL
            AND start_date = printf('%04d-%02d-16', year, month)
            AND end_date = date(printf('%04d-%02d-01', year, month), '+1 month', '-1 day')
            AND pay_date = min(printf('%04d-%02d-30', year, month), end_date))
    )
) STRICT;
-- @idx Un solo mes y una sola quincena 1 y 2 por año y mes.
CREATE UNIQUE INDEX ux_period_slot ON period (kind, year, month, coalesce(half, 0));
-- @idx Ubica el periodo que contiene una fecha.
CREATE INDEX ix_period_range ON period (start_date, end_date);
-- @idx Quincenas de un mes.
CREATE INDEX ix_period_parent ON period (parent_id);

-- @table cost_share_rule
-- @title Regla de reparto
-- @entity regla_reparto
-- @desc Cómo se divide un monto entre el hogar Axel–Ximena (a nombre del titular) y la participante externa: hogar, Karen o conjunto con su porcentaje. Las partes con porcentaje o monto fijo van en cost_share_rule_part y el resto cae en remainder_person_id.
-- @used Ajustes «Personas y reparto» y «Categorías». S2 (partes por gasto), S3 (partes de la compra a crédito). Filtro «familia» de C1.
CREATE TABLE cost_share_rule (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Nombre visible (Hogar, Karen, Conjunto 50/50) | no vacío
    remainder_person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona que recibe el resto del monto y absorbe los centavos de redondeo (en el hogar, el titular) | no debe repetirse en las partes (lo valida la API)
    is_template INTEGER NOT NULL DEFAULT 1 CHECK (is_template IN (0, 1)), -- 1 = regla reutilizable que se ofrece en capturas. 0 = regla de un solo uso (renglón «mixto» por cantidad) | booleano
    is_default INTEGER NOT NULL DEFAULT 0 CHECK (is_default IN (0, 1)), -- 1 = regla que aplica cuando ni el gasto ni su categoría traen regla (Hogar) | booleano, una sola
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (is_default = 0 OR is_template = 1) -- @rule La regla por defecto debe ser reutilizable.
) STRICT;
-- @idx Una sola regla por defecto.
CREATE UNIQUE INDEX ux_cost_share_rule_default ON cost_share_rule (is_default) WHERE is_default = 1;
-- @idx Nombres únicos entre las reglas reutilizables.
CREATE UNIQUE INDEX ux_cost_share_rule_template_name ON cost_share_rule (name COLLATE NOCASE) WHERE is_template = 1;
-- @idx Llave foránea hacia la persona del resto.
CREATE INDEX ix_cost_share_rule_remainder ON cost_share_rule (remainder_person_id);

-- @table cost_share_rule_part
-- @title Parte de una regla de reparto
-- @entity regla_reparto (detalle)
-- @desc Parte de una persona dentro de una regla: un porcentaje o un monto fijo. Lo que no se asigna aquí va a la persona del resto.
-- @used Ajustes «Personas y reparto». S2 y S3 al calcular partes.
CREATE TABLE cost_share_rule_part (
    id INTEGER PRIMARY KEY,
    rule_id INTEGER NOT NULL REFERENCES cost_share_rule (id) ON DELETE CASCADE, -- Regla a la que pertenece
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona de la parte
    pct_bp INTEGER CHECK (pct_bp BETWEEN 0 AND 10000), -- Porcentaje en puntos base: 5000 = 50.00 % (Conjunto 50/50: Karen 5000) | 0 a 10000, es decir 0 a 100 %
    fixed_cents INTEGER CHECK (fixed_cents >= 0), -- Monto fijo en centavos (parte fija de Karen en despensas futuras: 119150) | centavos
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((pct_bp IS NULL) <> (fixed_cents IS NULL)), -- @rule Cada parte es porcentaje o monto fijo, nunca ambos ni ninguno.
    UNIQUE (rule_id, person_id) -- @rule Una parte por persona en cada regla.
) STRICT;
-- @idx Llave foránea hacia la persona.
CREATE INDEX ix_cost_share_rule_part_person ON cost_share_rule_part (person_id);

-- @table category
-- @title Categoría
-- @entity categoria
-- @desc Categoría de gasto con jerarquía. Puede traer su regla de reparto por defecto (la despensa es Conjunto 50/50).
-- @used Ajustes «Categorías». S2 (presupuesto y real por categoría, alertas 80/100 %). Pantallas Gastos y Presupuesto. Filtro «categoría» de C1.
CREATE TABLE category (
    id INTEGER PRIMARY KEY,
    parent_id INTEGER REFERENCES category (id) ON DELETE RESTRICT, -- Categoría padre. NULL = raíz
    code TEXT UNIQUE CHECK (code IS NULL OR (length(code) BETWEEN 1 AND 40 AND code NOT GLOB '*[^a-z0-9_]*')), -- Clave estable para categorías que usan los motores (groceries, eating_out). NULL en las del usuario | a-z, 0-9 y guion bajo
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Nombre visible (Despensa, Comida fuera) | no vacío
    kind TEXT CHECK (kind IN ('fixed', 'variable')), -- Tipo de la categoría: fijo o variable (U1 N-281). Propone el tipo del gasto al registrarlo. NULL = sin clasificar | enumerado
    default_cost_share_rule_id INTEGER REFERENCES cost_share_rule (id) ON DELETE SET NULL, -- Reparto por categoría. NULL = la regla por defecto
    sort_order INTEGER NOT NULL DEFAULT 0, -- Orden entre hermanas
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (parent_id IS NULL OR parent_id <> id) -- @rule Una categoría no puede ser su propia madre.
) STRICT;
-- @idx Nombres únicos entre hermanas, sin distinguir mayúsculas.
CREATE UNIQUE INDEX ux_category_sibling_name ON category (coalesce(parent_id, 0), name COLLATE NOCASE);
-- @idx Árbol de la barra lateral en orden.
CREATE INDEX ix_category_parent ON category (parent_id, sort_order);
-- @idx Llave foránea hacia la regla de reparto.
CREATE INDEX ix_category_rule ON category (default_cost_share_rule_id);

-- @table budget
-- @title Presupuesto
-- @entity (complementaria) presupuesto mensual por categoría y total (S2, alertas 80/100 %)
-- @desc Monto mensual planeado por categoría o total, con vigencia: rige desde un mes hasta otro (inclusive) o sin fin. El monto de un solo mes es un renglón que empieza y termina en ese mes. Sirve para «real contra plan» y para las alertas al 80 % y 100 %.
-- @used Ajustes «Presupuesto» (N-275, N-278) y pantalla Presupuesto (N-087, una sola fuente). S2 (real contra plan), S7 (alertas).
CREATE TABLE budget (
    id INTEGER PRIMARY KEY,
    category_id INTEGER REFERENCES category (id) ON DELETE RESTRICT, -- Categoría. NULL = presupuesto total del mes
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0), -- Monto planeado por mes en centavos | centavos
    valid_from_date TEXT NOT NULL CHECK (date(valid_from_date) IS valid_from_date AND strftime('%d', valid_from_date) = '01'), -- Primer día del primer mes en que rige | fecha válida, día 1 del mes
    valid_to_date TEXT CHECK (date(valid_to_date) IS valid_to_date AND date(valid_to_date, '+1 day') = date(valid_to_date, '+1 day', 'start of month')), -- Último día del último mes en que rige, inclusive. NULL = rige hasta que se cambie. Cada cambio hecho en Ajustes elige su vigencia (K2): «de este mes en adelante» (preseleccionada) deja NULL; «solo este mes» pone el último día de ese mes y el monto anterior sigue desde el mes siguiente | fecha válida, último día de un mes
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (valid_to_date IS NULL OR valid_to_date >= valid_from_date) -- @rule La vigencia no termina antes de empezar.
) STRICT;
-- @idx Un renglón por categoría y mes de inicio (y uno solo de total por mes de inicio).
CREATE UNIQUE INDEX ux_budget_category_from ON budget (coalesce(category_id, 0), valid_from_date);
-- @idx Presupuesto vigente de una categoría en un mes e historial de cambios.
CREATE INDEX ix_budget_category ON budget (category_id, valid_from_date);

-- ---------------------------------------------------------------------
-- E. Tiendas, productos y presentaciones (se crean antes por dependencias)
-- ---------------------------------------------------------------------

-- @table store
-- @title Tienda
-- @entity tienda
-- @desc Cadena (Walmart, Bodega Aurrera, Sam's Club), sucursal de una cadena, o mercado o tianguis con su día de visita. La membresía y el permiso de lectura asistida viven en la cadena.
-- @used Ajustes «Tiendas y precios» (mercados y membresía Sam's). S5 (mejor tienda, punto de equilibrio de la membresía, PROFECO). Pantalla Compras y tiendas. Filtro «tienda» de C1.
CREATE TABLE store (
    id INTEGER PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('chain', 'branch', 'market')), -- chain = cadena, branch = sucursal de una cadena, market = mercado o tianguis | enumerado
    parent_id INTEGER REFERENCES store (id) ON DELETE RESTRICT, -- Cadena de la sucursal. Solo en branch
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Nombre: cadena (Bodega Aurrera), sucursal (Vallejo) o mercado | no vacío
    profeco_name TEXT, -- Nombre en PROFECO para cruzar precios: cadena_comercial en cadenas, nombre_comercial en sucursales
    website_url TEXT, -- Sitio de la cadena. Solo en chain
    assisted_reading_enabled INTEGER NOT NULL DEFAULT 0 CHECK (assisted_reading_enabled IN (0, 1)), -- Lectura asistida del sitio de la cadena. Queda en 0 hasta que docs/estado/S5-terminos.md tenga veredicto para ese sitio | booleano, solo chain
    address TEXT, -- Dirección (sucursal o mercado)
    municipality TEXT, -- Municipio o alcaldía, para cruzar con PROFECO
    state TEXT, -- Entidad federativa
    visit_weekday INTEGER CHECK (visit_weekday BETWEEN 1 AND 7), -- Día de visita del mercado: 1 = lunes … 7 = domingo (ISO 8601) | obligatorio en market, prohibido en las demás
    requires_membership INTEGER NOT NULL DEFAULT 0 CHECK (requires_membership IN (0, 1)), -- 1 = la cadena exige membresía (Sam's Club) | booleano, solo chain
    has_membership INTEGER NOT NULL DEFAULT 0 CHECK (has_membership IN (0, 1)), -- 1 = la casa ya tiene la membresía (hoy Sam's = 0 y se evalúa) | booleano
    membership_annual_fee_cents INTEGER CHECK (membership_annual_fee_cents >= 0), -- Costo anual de la membresía en centavos (lo captura Axel en Ajustes) | centavos
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((kind = 'branch') = (parent_id IS NOT NULL)), -- @rule Solo las sucursales tienen cadena madre.
    CHECK ((kind = 'market') = (visit_weekday IS NOT NULL)), -- @rule Todo mercado tiene día de visita y solo los mercados lo tienen.
    CHECK (kind = 'chain' OR (requires_membership = 0 AND assisted_reading_enabled = 0 AND website_url IS NULL)), -- @rule Membresía, sitio y lectura asistida son datos de la cadena.
    CHECK (requires_membership = 1 OR (has_membership = 0 AND membership_annual_fee_cents IS NULL)) -- @rule Sin membresía requerida no hay costo ni membresía.
) STRICT;
-- @idx Nombres únicos dentro de la misma cadena (y entre cadenas y mercados).
CREATE UNIQUE INDEX ux_store_name ON store (coalesce(parent_id, 0), name COLLATE NOCASE);
-- @idx Sucursales de una cadena.
CREATE INDEX ix_store_parent ON store (parent_id);
-- @idx Listas de mercados o cadenas activas.
CREATE INDEX ix_store_kind ON store (kind, is_active);

-- @table product_group
-- @title Grupo de producto
-- @entity (complementaria) decisión D: dos verticales de compra (mercados para verduras y carne no especializada, cadenas para el resto)
-- @desc Grupo de compra con su vertical: mercado o cadena. Permite repartir la lista de compras por mercado y por día de visita.
-- @used Ajustes «Tiendas y precios» e «Inventario». S4 (lista de compras por vertical), S5.
CREATE TABLE product_group (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) > 0), -- Grupo (Verduras, Carne no especializada, Abarrotes) | único
    purchase_channel TEXT NOT NULL CHECK (purchase_channel IN ('market', 'chain')), -- Vertical: market = mercado o tianguis, chain = Walmart, Bodega Aurrera o Sam's Club | enumerado
    sort_order INTEGER NOT NULL DEFAULT 0, -- Orden en la lista de compras
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;

-- @table store_product_group
-- @title Qué se compra en cada tienda
-- @entity tienda (detalle «qué se compra ahí»)
-- @desc Grupos de producto que se compran en cada mercado (o tienda). Es el «qué se compra ahí» que Axel captura en Ajustes.
-- @used Ajustes «Tiendas y precios». S4 (asigna renglones de la lista a un mercado y su día).
CREATE TABLE store_product_group (
    id INTEGER PRIMARY KEY,
    store_id INTEGER NOT NULL REFERENCES store (id) ON DELETE CASCADE, -- Mercado o tienda
    product_group_id INTEGER NOT NULL REFERENCES product_group (id) ON DELETE CASCADE, -- Grupo que se compra ahí
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (store_id, product_group_id) -- @rule Cada grupo una vez por tienda.
) STRICT;
-- @idx Tiendas que surten un grupo.
CREATE INDEX ix_store_product_group_group ON store_product_group (product_group_id);

-- @table product
-- @title Producto
-- @entity producto
-- @desc Producto genérico (ingrediente): Aceite de oliva, Arroz grano largo. Todas sus presentaciones son sustitutas entre sí. Define la unidad base de equivalencia.
-- @used Ajustes «Tiendas y precios», «Recetas» e «Inventario». S4 (ingredientes, kcal), S5 (equivalencias, mejor tienda por ingrediente).
CREATE TABLE product (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) > 0), -- Nombre genérico en español | único
    product_group_id INTEGER REFERENCES product_group (id) ON DELETE SET NULL, -- Grupo y vertical de compra
    base_unit TEXT NOT NULL CHECK (base_unit IN ('g', 'ml', 'piece')), -- Unidad base para equivalencias: g, ml o pieza | enumerado
    grams_per_unit_milli INTEGER CHECK (grams_per_unit_milli > 0), -- Gramos por unidad base ×1000 para calcular kcal: pieza = peso promedio, ml = densidad. NULL si base_unit = g
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (base_unit <> 'g' OR grams_per_unit_milli IS NULL) -- @rule En productos por gramo el factor a gramos sobra.
) STRICT;
-- @idx Productos de un grupo.
CREATE INDEX ix_product_group ON product (product_group_id);

-- @table presentation
-- @title Presentación
-- @entity presentacion
-- @desc Presentación comprable de un producto (marca y contenido), con su factor de equivalencia a la unidad base.
-- @used S4 (presentación comprable del ingrediente), S5 (precio por unidad base, sustitutos), Ajustes «Tiendas y precios».
CREATE TABLE presentation (
    id INTEGER PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES product (id) ON DELETE RESTRICT, -- Producto genérico
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Presentación como se ve en tienda (Nutrioli Oli 500 ml) | no vacío
    brand TEXT, -- Marca
    size_milli INTEGER NOT NULL CHECK (size_milli > 0), -- Contenido en milésimas de la unidad base: 500 ml = 500000, 750 g = 750000, 1 pieza = 1000 | factor de equivalencia, mayor que 0
    is_bulk INTEGER NOT NULL DEFAULT 0 CHECK (is_bulk IN (0, 1)), -- 1 = a granel: se vende por kilo o por unidad base | booleano
    barcode TEXT CHECK (barcode IS NULL OR (length(barcode) BETWEEN 8 AND 14 AND barcode NOT GLOB '*[^0-9]*')), -- Código de barras EAN o UPC (cruce con Open Food Facts) | 8 a 14 dígitos
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (product_id, name) -- @rule Nombres únicos dentro del mismo producto.
) STRICT;
-- @idx Un código de barras pertenece a una sola presentación.
CREATE UNIQUE INDEX ux_presentation_barcode ON presentation (barcode) WHERE barcode IS NOT NULL;

-- ---------------------------------------------------------------------
-- C. Créditos (catálogo, tasas y promociones)
-- ---------------------------------------------------------------------

-- @table credit
-- @title Crédito
-- @entity credito
-- @desc Crédito del titular: línea revolvente tipo Kueski, tarjeta bancaria, préstamo personal a plazo fijo o crédito de tienda con abonos. Guarda el límite, la lectura A/B, las fechas de corte y pago y las comisiones.
-- @used Ajustes «Créditos y condiciones». S3 (calendario unificado, disponible A/B, optimizador). Pantalla Créditos. Filtro «crédito» de C1.
CREATE TABLE credit (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) > 0), -- Nombre visible (Kueski Pay) | único
    kind TEXT NOT NULL CHECK (kind IN ('revolving_line', 'bank_card', 'term_loan', 'store_credit')), -- revolving_line = línea revolvente tipo Kueski, bank_card = tarjeta bancaria, term_loan = préstamo personal a plazo fijo, store_credit = crédito de tienda con abonos | enumerado
    lender TEXT, -- Institución que otorga el crédito
    holder_person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Titular del crédito (Axel)
    credit_limit_cents INTEGER CHECK (credit_limit_cents >= 0), -- Límite de la línea (Kueski 1000000) | centavos, obligatorio en revolving_line y bank_card
    payment_frequency TEXT NOT NULL CHECK (payment_frequency IN ('weekly', 'fortnightly', 'monthly')), -- Frecuencia de los pagos. fortnightly cae en los días 15 y 30 (o último día) | enumerado, monthly en tarjetas
    first_payment_min_days INTEGER CHECK (first_payment_min_days BETWEEN 0 AND 60), -- Días mínimos entre la compra y el primer pago (Kueski 15, oráculo) | 0 a 60
    availability_reading TEXT NOT NULL DEFAULT 'B' CHECK (availability_reading IN ('A', 'B')), -- Lectura del disponible (interruptor de Ajustes, solo aplica en revolving_line): B = el límite es el tope total (por defecto), A = el límite está libre en la fecha ancla y los pagos suman encima | A o B
    availability_confirmed INTEGER NOT NULL DEFAULT 0 CHECK (availability_confirmed IN (0, 1)), -- 0 = mostrar el aviso «pendiente confirmar» con el acreedor | booleano
    reading_a_anchor_date TEXT CHECK (date(reading_a_anchor_date) IS reading_a_anchor_date), -- Fecha «hoy» de la lectura A: lo comprado hasta ese día ya está dentro del límite libre (Kueski 2026-10-10) | fecha válida, solo revolving_line
    statement_day INTEGER CHECK (statement_day BETWEEN 1 AND 31), -- Día de corte de la tarjeta | 1 a 31, solo bank_card
    payment_due_day INTEGER CHECK (payment_due_day BETWEEN 1 AND 31), -- Día límite de pago de la tarjeta | 1 a 31, solo bank_card
    annual_rate_num INTEGER CHECK (annual_rate_num >= 0), -- Numerador de la tasa anual ordinaria (tarjeta, préstamo) | junto con annual_rate_den
    annual_rate_den INTEGER CHECK (annual_rate_den > 0), -- Denominador de la tasa anual ordinaria | mayor que 0
    vat_on_interest_bp INTEGER CHECK (vat_on_interest_bp BETWEEN 0 AND 10000), -- IVA que se suma a los intereses cuando la tasa no lo trae (1600 = 16 %). NULL = la tasa o la cotización ya incluye IVA (Kueski) | 0 a 10000
    annual_fee_cents INTEGER CHECK (annual_fee_cents >= 0), -- Comisión anual | centavos
    late_fee_cents INTEGER CHECK (late_fee_cents >= 0), -- Comisión por pago tardío | centavos
    opening_fee_cents INTEGER CHECK (opening_fee_cents >= 0), -- Comisión por apertura | centavos
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (kind NOT IN ('revolving_line', 'bank_card') OR credit_limit_cents IS NOT NULL), -- @rule Línea revolvente y tarjeta tienen límite.
    CHECK ((kind = 'bank_card') = (statement_day IS NOT NULL AND payment_due_day IS NOT NULL)), -- @rule Solo la tarjeta tiene día de corte y día de pago, y siempre los tiene.
    CHECK (kind <> 'bank_card' OR payment_frequency = 'monthly'), -- @rule La tarjeta se paga mensualmente.
    CHECK (kind = 'revolving_line' OR reading_a_anchor_date IS NULL), -- @rule La fecha ancla de la lectura A solo existe en líneas revolventes.
    CHECK (kind <> 'revolving_line' OR availability_reading = 'B' OR reading_a_anchor_date IS NOT NULL), -- @rule Con lectura A hace falta la fecha ancla.
    CHECK ((annual_rate_num IS NULL) = (annual_rate_den IS NULL)) -- @rule La tasa anual va completa (numerador y denominador) o no va.
) STRICT;
-- @idx Créditos del titular.
CREATE INDEX ix_credit_holder ON credit (holder_person_id);
-- @idx Créditos activos por tipo.
CREATE INDEX ix_credit_kind ON credit (kind, is_active);

-- @table credit_term_rate
-- @title Tasa por plazo
-- @entity tasa_por_plazo
-- @desc Tabla editable de tasas por plazo de un crédito, como fracción exacta del capital por todo el plazo con IVA incluido. Kueski: 2 quincenas = 53450/309450 (17.27 %), 4 quincenas = 80050/309450 (25.8685 %).
-- @used Ajustes «Créditos y condiciones». S3 (estimar intereses de compras planeadas, escenarios del optimizador).
CREATE TABLE credit_term_rate (
    id INTEGER PRIMARY KEY,
    credit_id INTEGER NOT NULL REFERENCES credit (id) ON DELETE CASCADE, -- Crédito
    term_count INTEGER NOT NULL CHECK (term_count BETWEEN 1 AND 120), -- Plazo en número de pagos (2, 4, 6, 8, 10 o 12 quincenas en Kueski) | 1 a 120
    rate_num INTEGER NOT NULL CHECK (rate_num >= 0), -- Numerador: si viene de una cotización, los intereses en centavos (80050) | entero mayor o igual a 0
    rate_den INTEGER NOT NULL CHECK (rate_den > 0), -- Denominador: si viene de una cotización, el capital en centavos (309450). Una tasa tecleada como 25.8685 % se guarda 258685/1000000 | mayor que 0
    source TEXT NOT NULL CHECK (source IN ('quote', 'manual')), -- quote = sale de una cotización real, manual = tecleada en la tabla editable | enumerado
    quote_date TEXT CHECK (date(quote_date) IS quote_date), -- Fecha de la cotización (2026-10-10) | obligatoria si source = quote
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (credit_id, term_count), -- @rule Una tasa vigente por plazo y crédito.
    CHECK (source = 'manual' OR quote_date IS NOT NULL) -- @rule Una tasa de cotización dice de qué fecha es.
) STRICT;

-- @table promotion
-- @title Promoción
-- @entity promocion
-- @desc Calendario de promociones capturado a mano: comercio, plazo, tasa (0 % casi siempre), compra mínima y vigencia.
-- @used Ajustes «Créditos y condiciones» (promociones). S3 (escenarios con promoción), S5 (comprar más por peso).
CREATE TABLE promotion (
    id INTEGER PRIMARY KEY,
    credit_id INTEGER REFERENCES credit (id) ON DELETE CASCADE, -- Crédito al que aplica (Kueski). NULL = cualquiera
    merchant_name TEXT NOT NULL CHECK (length(trim(merchant_name)) > 0), -- Comercio tal como lo anuncia la promoción (Chedraui) | no vacío
    store_id INTEGER REFERENCES store (id) ON DELETE SET NULL, -- Tienda registrada que corresponde al comercio, si existe
    term_count INTEGER NOT NULL CHECK (term_count BETWEEN 1 AND 120), -- Plazo en pagos (4 quincenas) | 1 a 120
    rate_num INTEGER NOT NULL DEFAULT 0 CHECK (rate_num >= 0), -- Numerador de la tasa promocional por todo el plazo. 0 = sin intereses | mayor o igual a 0
    rate_den INTEGER NOT NULL DEFAULT 1 CHECK (rate_den > 0), -- Denominador de la tasa promocional | mayor que 0
    min_purchase_cents INTEGER NOT NULL DEFAULT 0 CHECK (min_purchase_cents >= 0), -- Compra mínima (200000 a 250000 en el tablero 2026) | centavos
    valid_from TEXT NOT NULL CHECK (date(valid_from) IS valid_from), -- Inicio de vigencia, inclusive | fecha válida
    valid_to TEXT NOT NULL CHECK (date(valid_to) IS valid_to), -- Fin de vigencia, inclusive (Chedraui: solo 2026-10-30) | fecha válida
    source_note TEXT, -- De dónde se tomó (tablero oficial de promociones Kueski 2026)
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (valid_to >= valid_from) -- @rule La vigencia termina el mismo día o después de empezar.
) STRICT;
-- @idx Promociones vigentes en una fecha.
CREATE INDEX ix_promotion_validity ON promotion (valid_from, valid_to);
-- @idx Promociones de un crédito.
CREATE INDEX ix_promotion_credit ON promotion (credit_id, valid_from);
-- @idx Llave foránea hacia la tienda.
CREATE INDEX ix_promotion_store ON promotion (store_id);

-- ---------------------------------------------------------------------
-- B (cont.). Plan, ingresos y gastos
-- ---------------------------------------------------------------------

-- @table planned_item
-- @title Partida planificada
-- @entity ingreso (fijo quincenal) y gasto (fijo y venidero): parte recurrente. Coincide con datos/semilla/gastos_fijos_ingresos.csv (plan §8.2)
-- @desc Ingreso o gasto esperado con su periodicidad: sueldo quincenal, renta, despensa bimestral con Kueski, gastos venideros. Los motores proyectan sus ocurrencias, no se guardan una por una.
-- @used Ajustes «Ingresos» y «Presupuesto». S2 (plan del mes, proyección 12 meses), S3 (compras planeadas con crédito). Importador del CSV de Axel.
CREATE TABLE planned_item (
    id INTEGER PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('income', 'fixed_expense', 'upcoming_expense')), -- Tipo (columna «tipo» del CSV): ingreso, gasto fijo o gasto venidero | enumerado
    description TEXT NOT NULL CHECK (length(trim(description)) > 0), -- Concepto | no vacío
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0), -- Monto de cada ocurrencia en centavos (sueldo por quincena, renta por mes). En quincenal es el de la quincena del 15, y también el del 30 si second_half_amount_cents es NULL | centavos
    second_half_amount_cents INTEGER CHECK (second_half_amount_cents >= 0), -- Monto de la quincena del 30 (o último día del mes) cuando difiere del de la quincena del 15 (N-279). NULL = las dos quincenas iguales | centavos, solo en frequency = fortnightly
    frequency TEXT NOT NULL CHECK (frequency IN ('weekly', 'fortnightly', 'monthly', 'yearly', 'once')), -- Periodicidad. fortnightly cae en los días de pago 15 y 30 (o último día). Las demás se repiten desde start_date | enumerado
    interval_count INTEGER NOT NULL DEFAULT 1 CHECK (interval_count BETWEEN 1 AND 24), -- Cada cuántas unidades se repite (despensa bimestral = monthly cada 2) | 1 a 24
    start_date TEXT NOT NULL CHECK (date(start_date) IS start_date), -- Primera ocurrencia o fecha única. Su día del mes es el ancla y se recorta al último día en meses cortos | fecha válida
    end_date TEXT CHECK (date(end_date) IS end_date), -- Última fecha posible, inclusive. NULL = sin fin | fecha válida
    category_id INTEGER REFERENCES category (id) ON DELETE RESTRICT, -- Categoría del gasto | obligatoria en gastos, NULL en ingresos
    person_id INTEGER REFERENCES person (id) ON DELETE RESTRICT, -- Persona a quien se atribuye. En ingresos solo el titular (disparador). NULL = hogar en general
    cost_share_rule_id INTEGER REFERENCES cost_share_rule (id) ON DELETE RESTRICT, -- Reparto del gasto. NULL = el de la categoría o el de por defecto
    payment_method TEXT CHECK (payment_method IN ('cash', 'debit', 'transfer', 'credit')), -- Cómo se pagará cada ocurrencia. NULL = sin definir | enumerado, NULL en ingresos
    credit_id INTEGER REFERENCES credit (id) ON DELETE RESTRICT, -- Crédito con que se planea pagar (despensa con Kueski) | obligatorio si payment_method = credit
    term_count INTEGER CHECK (term_count BETWEEN 1 AND 120), -- Plazo planeado en pagos (4 quincenas) | obligatorio con credit_id
    accrual_months INTEGER NOT NULL DEFAULT 1 CHECK (accrual_months BETWEEN 1 AND 24), -- Meses en que se devenga cada ocurrencia (despensa bimestral = 2) | 1 a 24
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (end_date IS NULL OR end_date >= start_date), -- @rule La vigencia no termina antes de empezar.
    CHECK (frequency <> 'once' OR (interval_count = 1 AND (end_date IS NULL OR end_date = start_date))), -- @rule Una partida única ocurre solo en start_date.
    CHECK (kind <> 'income' OR (category_id IS NULL AND cost_share_rule_id IS NULL AND payment_method IS NULL AND credit_id IS NULL AND accrual_months = 1)), -- @rule Un ingreso no lleva categoría, reparto, medio de pago ni devengo.
    CHECK (kind = 'income' OR category_id IS NOT NULL), -- @rule Todo gasto planificado tiene categoría.
    CHECK ((payment_method IS 'credit') = (credit_id IS NOT NULL)), -- @rule Hay crédito si y solo si el medio de pago es crédito.
    CHECK (second_half_amount_cents IS NULL OR frequency = 'fortnightly'), -- @rule El monto de la segunda quincena solo existe en partidas quincenales.
    CHECK ((credit_id IS NULL) = (term_count IS NULL)) -- @rule El plazo acompaña siempre al crédito.
) STRICT;
-- @idx Partidas activas por tipo y fecha (proyección).
CREATE INDEX ix_planned_item_kind ON planned_item (kind, is_active, start_date);
-- @idx Plan por categoría.
CREATE INDEX ix_planned_item_category ON planned_item (category_id);
-- @idx Plan por persona.
CREATE INDEX ix_planned_item_person ON planned_item (person_id);
-- @idx Compras planeadas con un crédito.
CREATE INDEX ix_planned_item_credit ON planned_item (credit_id);
-- @idx Llave foránea hacia la regla de reparto.
CREATE INDEX ix_planned_item_rule ON planned_item (cost_share_rule_id);

-- @table income
-- @title Ingreso
-- @entity ingreso
-- @desc Ingreso recibido por el titular: fijo quincenal (lo captura Axel cada quincena) o extra capturado a mano. Ximena no aporta ingresos al modelo.
-- @used Ajustes «Ingresos». S2 (real contra plan, flujo por quincena). Pantallas Inicio y Presupuesto.
CREATE TABLE income (
    id INTEGER PRIMARY KEY,
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Quién recibe el ingreso | solo el titular (disparador)
    kind TEXT NOT NULL CHECK (kind IN ('fixed', 'extra')), -- fixed = ingreso fijo quincenal, extra = ingreso extra | enumerado
    planned_item_id INTEGER REFERENCES planned_item (id) ON DELETE SET NULL, -- Partida planificada que cumple (sueldo quincenal). NULL en extras sin plan
    description TEXT NOT NULL CHECK (length(trim(description)) > 0), -- Concepto | no vacío
    amount_cents INTEGER NOT NULL CHECK (amount_cents > 0), -- Monto recibido en centavos | mayor que 0
    income_date TEXT NOT NULL CHECK (date(income_date) IS income_date), -- Fecha en que se recibió (calendario de America/Mexico_City) | fecha válida
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Ingresos de un periodo (rango de fechas).
CREATE INDEX ix_income_date ON income (income_date);
-- @idx Ingresos por persona y fecha.
CREATE INDEX ix_income_person_date ON income (person_id, income_date);
-- @idx Ingresos que cumplen una partida.
CREATE INDEX ix_income_planned ON income (planned_item_id);

-- @table expense
-- @title Gasto
-- @entity gasto
-- @desc Gasto real: fijo, variable o venidero. Guarda la fecha que inicia el devengo, los meses de devengo, el medio de pago (si es crédito el efectivo sale por los pagos programados), quién pagó, a quién se atribuye y su reparto.
-- @used Pantallas Gastos (registro en 3 toques), Inicio y Presupuesto. S2 (devengado contra efectivo, real contra plan), S3 (compras a crédito), S7 (reportes).
CREATE TABLE expense (
    id INTEGER PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('fixed', 'variable', 'upcoming')), -- fixed = gasto fijo, variable = gasto del día a día, upcoming = gasto venidero que ya ocurrió o quedó comprometido | enumerado
    description TEXT CHECK (description IS NULL OR length(trim(description)) > 0), -- Concepto o comercio. NULL = sin descripción: el registro en 3 toques no la pide (N-054) y la pantalla muestra la categoría | opcional; si viene, no vacía
    category_id INTEGER NOT NULL REFERENCES category (id) ON DELETE RESTRICT, -- Categoría
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0), -- Importe de lo comprado en centavos, con IVA y descuentos y sin intereses de financiamiento (despensa 1 = 309450) | centavos
    expense_date TEXT NOT NULL CHECK (date(expense_date) IS expense_date), -- Fecha de la compra o consumo. El devengo empieza en su mes | fecha válida
    accrual_months INTEGER NOT NULL DEFAULT 1 CHECK (accrual_months BETWEEN 1 AND 24), -- Meses calendario en que se devenga desde el mes de expense_date (despensa bimestral = 2, mitad por mes) | 1 a 24
    payment_method TEXT NOT NULL CHECK (payment_method IN ('cash', 'debit', 'transfer', 'credit')), -- Medio de pago. credit = el efectivo sale según los pagos programados de su credit_purchase | enumerado
    paid_by_person_id INTEGER REFERENCES person (id) ON DELETE RESTRICT, -- Quién pagó de su bolsillo (papel de contado de Karen). NULL = el titular | NULL si payment_method = credit
    person_id INTEGER REFERENCES person (id) ON DELETE RESTRICT, -- Persona a quien se atribuye (comida fuera de Ximena). NULL = hogar en general
    cost_share_rule_id INTEGER REFERENCES cost_share_rule (id) ON DELETE RESTRICT, -- Reparto del gasto. NULL = la regla de la categoría o la de por defecto
    store_id INTEGER REFERENCES store (id) ON DELETE RESTRICT, -- Tienda donde se compró
    planned_item_id INTEGER REFERENCES planned_item (id) ON DELETE SET NULL, -- Partida planificada que cumple (gasto fijo o venidero)
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (payment_method <> 'credit' OR paid_by_person_id IS NULL) -- @rule En compras a crédito quién paga lo dicen las partes de cada pago programado.
) STRICT;
-- @idx Gastos de un periodo (rango de fechas).
CREATE INDEX ix_expense_date ON expense (expense_date);
-- @idx Gastos por categoría y periodo (real contra presupuesto).
CREATE INDEX ix_expense_category_date ON expense (category_id, expense_date);
-- @idx Gastos por persona y periodo (comida fuera por persona).
CREATE INDEX ix_expense_person_date ON expense (person_id, expense_date);
-- @idx Gastos por tienda y periodo.
CREATE INDEX ix_expense_store_date ON expense (store_id, expense_date);
-- @idx Gastos pagados por una persona de su bolsillo.
CREATE INDEX ix_expense_paid_by ON expense (paid_by_person_id);
-- @idx Llave foránea hacia la regla de reparto.
CREATE INDEX ix_expense_rule ON expense (cost_share_rule_id);
-- @idx Gastos que cumplen una partida planificada.
CREATE INDEX ix_expense_planned ON expense (planned_item_id);

-- @table expense_item
-- @title Renglón de gasto
-- @entity gasto (detalle: renglones del ticket)
-- @desc Renglón de un gasto con ticket (la despensa): presentación, cantidad, importe y reparto del renglón (cuenta axel, karen, conjunto o mixto de data.json).
-- @used S2 y S3 (el crédito se reparte según lo que cada quien compró), S5 (precio desde ticket), S4 (entrada al inventario). Pantalla Gastos (detalle).
CREATE TABLE expense_item (
    id INTEGER PRIMARY KEY,
    expense_id INTEGER NOT NULL REFERENCES expense (id) ON DELETE CASCADE, -- Gasto al que pertenece
    line_no INTEGER NOT NULL CHECK (line_no >= 1), -- Orden del renglón en el ticket | mayor o igual a 1
    description TEXT NOT NULL CHECK (length(trim(description)) > 0), -- Texto del renglón como viene en el ticket | no vacío
    presentation_id INTEGER REFERENCES presentation (id) ON DELETE RESTRICT, -- Presentación del catálogo. NULL mientras no se liga
    quantity_milli INTEGER NOT NULL CHECK (quantity_milli > 0), -- Cantidad de presentaciones ×1000: 3 piezas = 3000, 1.250 kg a granel = 1250 | mayor que 0
    unit_price_cents INTEGER CHECK (unit_price_cents >= 0), -- Precio por presentación en centavos | centavos
    line_total_cents INTEGER NOT NULL CHECK (line_total_cents >= 0), -- Importe del renglón con descuentos e IVA. La suma de renglones es el importe del gasto | centavos
    cost_share_rule_id INTEGER REFERENCES cost_share_rule (id) ON DELETE RESTRICT, -- Reparto del renglón (Hogar, Karen, Conjunto 50/50 o una regla de un solo uso para «mixto»). NULL = el del gasto
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (expense_id, line_no) -- @rule Un número de renglón por gasto.
) STRICT;
-- @idx Compras de una presentación (historial e inventario).
CREATE INDEX ix_expense_item_presentation ON expense_item (presentation_id);
-- @idx Llave foránea hacia la regla de reparto.
CREATE INDEX ix_expense_item_rule ON expense_item (cost_share_rule_id);

-- @table expense_share
-- @title Parte de un gasto
-- @entity regla_reparto (resultado congelado por gasto)
-- @desc Parte de cada persona en un gasto, ya calculada al centavo. Se congela al guardar para que cambiar una regla no altere gastos pasados. La parte del hogar Axel–Ximena va a nombre del titular.
-- @used S2 (gasto del hogar contra parte de Karen), S7 (estado de cuenta de Karen). Filtros «persona» y «familia» de C1.
CREATE TABLE expense_share (
    id INTEGER PRIMARY KEY,
    expense_id INTEGER NOT NULL REFERENCES expense (id) ON DELETE CASCADE, -- Gasto
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona que asume la parte
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0), -- Parte en centavos (despensa 1: Axel 190300, Karen 119150). La suma de partes es el importe del gasto | centavos
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (expense_id, person_id) -- @rule Una parte por persona y gasto.
) STRICT;
-- @idx Partes de una persona (consultas por persona y familia).
CREATE INDEX ix_expense_share_person ON expense_share (person_id);

-- @table attachment
-- @title Adjunto
-- @entity gasto (adjunto)
-- @desc Archivo adjunto: foto de ticket (con su OCR local), factura o estado de cuenta del banco. Ninguna cifra del OCR entra al sistema sin confirmación humana.
-- @used Pantalla Gastos (adjuntar ticket). S5 (OCR Tesseract con confirmación), S6 (almacenamiento de adjuntos).
CREATE TABLE attachment (
    id INTEGER PRIMARY KEY,
    expense_id INTEGER REFERENCES expense (id) ON DELETE SET NULL, -- Gasto al que pertenece. NULL mientras el ticket espera confirmación
    kind TEXT NOT NULL CHECK (kind IN ('ticket_photo', 'invoice', 'bank_statement', 'other')), -- Foto de ticket, factura, estado de cuenta del banco u otro | enumerado
    file_path TEXT NOT NULL UNIQUE, -- Ruta relativa dentro de la carpeta de datos del sistema | única
    original_name TEXT, -- Nombre original del archivo
    mime_type TEXT NOT NULL, -- Tipo MIME (image/jpeg, application/pdf)
    size_bytes INTEGER NOT NULL CHECK (size_bytes >= 0), -- Tamaño en bytes | mayor o igual a 0
    sha256 TEXT NOT NULL CHECK (length(sha256) = 64), -- Huella del archivo (detecta duplicados y corrupción) | 64 caracteres
    ocr_status TEXT NOT NULL DEFAULT 'none' CHECK (ocr_status IN ('none', 'pending', 'done', 'confirmed', 'failed')), -- Estado del OCR local: none, pending, done (falta confirmar), confirmed, failed | enumerado
    ocr_engine TEXT, -- Motor y versión (tesseract 5.5.3 spa)
    ocr_text TEXT, -- Texto crudo reconocido
    uploaded_by_user_id INTEGER REFERENCES user_account (id) ON DELETE SET NULL, -- Cuenta que lo subió
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Adjuntos de un gasto.
CREATE INDEX ix_attachment_expense ON attachment (expense_id);
-- @idx Detecta archivos repetidos por su huella.
CREATE INDEX ix_attachment_sha ON attachment (sha256);
-- @idx Llave foránea hacia la cuenta que subió.
CREATE INDEX ix_attachment_uploader ON attachment (uploaded_by_user_id);
-- @idx Cola de tickets con OCR pendiente de procesar o de confirmar.
CREATE INDEX ix_attachment_ocr_queue ON attachment (ocr_status, created_at) WHERE ocr_status IN ('pending', 'done');

-- @table price
-- @title Precio
-- @entity precio
-- @desc Precio observado de una presentación en una tienda, con fecha y fuente autorizada. Un precio es viejo a los 30 días y pesa menos en la comparativa.
-- @used S5 (mejor tienda por ingrediente y canasta, frescura, membresía), S4 (costo de receta), S7 (alerta de precios viejos). Pantalla Compras y tiendas.
CREATE TABLE price (
    id INTEGER PRIMARY KEY,
    presentation_id INTEGER NOT NULL REFERENCES presentation (id) ON DELETE RESTRICT, -- Presentación
    store_id INTEGER NOT NULL REFERENCES store (id) ON DELETE RESTRICT, -- Tienda del precio (sucursal, cadena o mercado)
    price_cents INTEGER NOT NULL CHECK (price_cents >= 0), -- Precio de una presentación en centavos con IVA (Nutrioli Oli 500 ml en Bodega Aurrera: 15000) | centavos
    observed_date TEXT NOT NULL CHECK (date(observed_date) IS observed_date), -- Fecha del precio. Es viejo cuando tiene 30 días o más | fecha válida
    source TEXT NOT NULL CHECK (source IN ('ticket_manual', 'ticket_ocr', 'site_assisted', 'profeco', 'manual')), -- Fuente: ticket capturado a mano, ticket por OCR ya confirmado, lectura asistida del sitio, PROFECO u otra captura manual | enumerado
    is_promotional INTEGER NOT NULL DEFAULT 0 CHECK (is_promotional IN (0, 1)), -- 1 = precio de promoción (sexta lata de atún a 2100), no se toma como precio regular | booleano
    source_ref TEXT, -- Referencia de la fuente (mes del dataset PROFECO, URL)
    expense_item_id INTEGER REFERENCES expense_item (id) ON DELETE SET NULL, -- Renglón de ticket que originó el precio
    attachment_id INTEGER REFERENCES attachment (id) ON DELETE SET NULL, -- Foto del ticket
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (presentation_id, store_id, observed_date, source, is_promotional) -- @rule No se repite el mismo precio el mismo día con la misma fuente.
) STRICT;
-- @idx Último precio de una presentación en una tienda.
CREATE INDEX ix_price_latest ON price (presentation_id, store_id, observed_date);
-- @idx Precios de una tienda por fecha (canasta, precios viejos).
CREATE INDEX ix_price_store_date ON price (store_id, observed_date);
-- @idx Precios por fecha (frescura).
CREATE INDEX ix_price_date ON price (observed_date);
-- @idx Llave foránea hacia el renglón de ticket.
CREATE INDEX ix_price_item ON price (expense_item_id);
-- @idx Llave foránea hacia la foto del ticket.
CREATE INDEX ix_price_attachment ON price (attachment_id);

-- ---------------------------------------------------------------------
-- C (cont.). Compras a crédito, pagos programados y pagos hechos
-- ---------------------------------------------------------------------

-- @table credit_purchase
-- @title Compra a crédito
-- @entity compra_credito
-- @desc Compra o disposición con un crédito: capital, intereses cotizados (IVA incluido), plazo, primer pago y promoción. Solo el capital usa y repone la línea.
-- @used S3 (calendario, disponible A/B, compras vivas, regresión Kueski), S2 (devengado incluye intereses), S7 (estado de cuenta de Karen). Pantalla Créditos.
CREATE TABLE credit_purchase (
    id INTEGER PRIMARY KEY,
    credit_id INTEGER NOT NULL REFERENCES credit (id) ON DELETE RESTRICT, -- Crédito usado
    expense_id INTEGER UNIQUE REFERENCES expense (id) ON DELETE RESTRICT, -- Gasto que originó la compra (despensa 1). NULL en disposiciones sin gasto (desembolso de préstamo) | un gasto tiene a lo más una compra a crédito
    kind TEXT NOT NULL CHECK (kind IN ('installment_plan', 'single_charge', 'loan_disbursement')), -- installment_plan = pagos fijos (Kueski, meses sin intereses, abonos de tienda), single_charge = cargo normal de tarjeta, loan_disbursement = préstamo con tabla de amortización | enumerado
    description TEXT, -- Descripción corta
    purchase_date TEXT NOT NULL CHECK (date(purchase_date) IS purchase_date), -- Fecha de la compra (2026-10-10) | fecha válida
    principal_cents INTEGER NOT NULL CHECK (principal_cents >= 0), -- Capital: lo que usa la línea (309450) | centavos
    interest_cents INTEGER NOT NULL DEFAULT 0 CHECK (interest_cents >= 0), -- Intereses de todo el plazo, IVA incluido, tal como los cotizó el acreedor (80050). Nunca reponen la línea | centavos
    total_cents INTEGER NOT NULL, -- Total a pagar (389500) | igual a capital más intereses
    term_count INTEGER NOT NULL CHECK (term_count BETWEEN 1 AND 360), -- Número de pagos (4) | 1 a 360
    payment_frequency TEXT NOT NULL CHECK (payment_frequency IN ('weekly', 'fortnightly', 'monthly')), -- Frecuencia de los pagos (quincenal) | enumerado
    first_payment_date TEXT NOT NULL CHECK (date(first_payment_date) IS first_payment_date), -- Primer pago (2026-10-30) | no antes de la compra
    interest_source TEXT NOT NULL CHECK (interest_source IN ('quoted', 'term_table', 'promotion', 'amortization')), -- Origen de los intereses: quoted = cotización del acreedor al pagar, term_table = tabla por plazo, promotion = promoción, amortization = tabla de amortización | enumerado
    rate_num INTEGER CHECK (rate_num >= 0), -- Numerador de la tasa de esta compra (cotización: intereses/capital = 80050/309450) | junto con rate_den y rate_basis
    rate_den INTEGER CHECK (rate_den > 0), -- Denominador de la tasa | mayor que 0
    rate_basis TEXT CHECK (rate_basis IN ('term_total', 'annual')), -- term_total = sobre el capital por todo el plazo (Kueski), annual = tasa anual (préstamo) | enumerado
    interest_includes_vat INTEGER NOT NULL DEFAULT 1 CHECK (interest_includes_vat IN (0, 1)), -- 1 = los intereses ya traen IVA (Kueski) | booleano
    promotion_id INTEGER REFERENCES promotion (id) ON DELETE SET NULL, -- Promoción aplicada
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paid_off', 'cancelled')), -- active = con pagos pendientes (compra viva), paid_off = liquidada, cancelled = cancelada | enumerado
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (total_cents = principal_cents + interest_cents), -- @rule Total = capital + intereses.
    CHECK (first_payment_date >= purchase_date), -- @rule El primer pago no es anterior a la compra.
    CHECK ((rate_num IS NULL AND rate_den IS NULL AND rate_basis IS NULL) OR (rate_num IS NOT NULL AND rate_den IS NOT NULL AND rate_basis IS NOT NULL)), -- @rule La tasa va completa (numerador, denominador y base) o no va.
    CHECK (interest_source <> 'promotion' OR promotion_id IS NOT NULL), -- @rule Si los intereses salen de una promoción, la promoción se identifica.
    CHECK (kind <> 'single_charge' OR term_count = 1) -- @rule Un cargo normal de tarjeta se paga en un solo pago.
) STRICT;
-- @idx Compras de un crédito por fecha.
CREATE INDEX ix_credit_purchase_credit_date ON credit_purchase (credit_id, purchase_date);
-- @idx Compras vivas de un crédito (restricción de máximo de compras vivas).
CREATE INDEX ix_credit_purchase_credit_status ON credit_purchase (credit_id, status);
-- @idx Llave foránea hacia la promoción.
CREATE INDEX ix_credit_purchase_promotion ON credit_purchase (promotion_id);

-- @table credit_purchase_share
-- @title Parte de una compra a crédito
-- @entity compra_credito (partes por persona)
-- @desc Capital, intereses y total de cada persona en una compra a crédito. El capital sale de lo que cada quien compró y los intereses se reparten en proporción al capital, con el resto para el titular.
-- @used S3 (pagos por persona), S2, S7 (estado de cuenta de Karen).
CREATE TABLE credit_purchase_share (
    id INTEGER PRIMARY KEY,
    credit_purchase_id INTEGER NOT NULL REFERENCES credit_purchase (id) ON DELETE CASCADE, -- Compra
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona
    principal_cents INTEGER NOT NULL CHECK (principal_cents >= 0), -- Capital de la persona (Axel 190300, Karen 119150) | centavos
    interest_cents INTEGER NOT NULL CHECK (interest_cents >= 0), -- Intereses de la persona (Axel 49228, Karen 30822) | centavos
    total_cents INTEGER NOT NULL, -- Total de la persona (Axel 239528, Karen 149972) | igual a capital más intereses
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (total_cents = principal_cents + interest_cents), -- @rule Total de la persona = su capital + sus intereses.
    UNIQUE (credit_purchase_id, person_id) -- @rule Una parte por persona y compra.
) STRICT;
-- @idx Compras a crédito de una persona.
CREATE INDEX ix_credit_purchase_share_person ON credit_purchase_share (person_id);

-- @table scheduled_payment
-- @title Pago programado
-- @entity pago_programado
-- @desc Cada pago de la serie de una compra: fecha, monto, capital que repone e intereses. Se guarda ya calculado al centavo (el último de la serie ajusta centavos).
-- @used S3 (calendario unificado de pagos, disponible A/B), S2 (efectivo por quincena), S7 (alerta de pago próximo). Pantallas Créditos e Inicio (calendario).
CREATE TABLE scheduled_payment (
    id INTEGER PRIMARY KEY,
    credit_purchase_id INTEGER NOT NULL REFERENCES credit_purchase (id) ON DELETE CASCADE, -- Compra a la que pertenece
    seq INTEGER NOT NULL CHECK (seq >= 1), -- Número de pago dentro de la serie (1 a term_count) | mayor o igual a 1
    due_date TEXT NOT NULL CHECK (date(due_date) IS due_date), -- Fecha de pago (2026-10-30, 2026-11-15, 2026-11-30, 2026-12-15) | fecha válida
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0), -- Monto del pago (97375) | centavos
    principal_cents INTEGER NOT NULL CHECK (principal_cents >= 0), -- Capital que repone la línea al pagarse (77363, el último 77361) | centavos
    interest_cents INTEGER NOT NULL CHECK (interest_cents >= 0), -- Intereses del pago, no reponen la línea (20012, el último 20014) | centavos
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (amount_cents = principal_cents + interest_cents), -- @rule Monto = capital + intereses.
    UNIQUE (credit_purchase_id, seq) -- @rule Un número de pago por compra.
) STRICT;
-- @idx Calendario unificado: pagos por fecha (pago próximo, flujo por quincena).
CREATE INDEX ix_scheduled_payment_due ON scheduled_payment (due_date);

-- @table scheduled_payment_share
-- @title Parte de un pago programado
-- @entity pago_programado (partes por persona)
-- @desc Parte de cada persona en un pago programado. La de Karen es informativa para el flujo de Axel porque ella paga directo al acreedor.
-- @used S3 y S2 (efectivo del titular contra parte informativa), S7 (estado de cuenta de Karen, alerta de pago).
CREATE TABLE scheduled_payment_share (
    id INTEGER PRIMARY KEY,
    scheduled_payment_id INTEGER NOT NULL REFERENCES scheduled_payment (id) ON DELETE CASCADE, -- Pago programado
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona
    amount_cents INTEGER NOT NULL CHECK (amount_cents >= 0), -- Parte de la persona (Axel 59882, Karen 37493). La suma de partes es el monto del pago | centavos
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (scheduled_payment_id, person_id) -- @rule Una parte por persona y pago.
) STRICT;
-- @idx Pagos por persona (estado de cuenta de Karen).
CREATE INDEX ix_scheduled_payment_share_person ON scheduled_payment_share (person_id);

-- @table credit_statement
-- @title Estado de cuenta de tarjeta
-- @entity credito (fechas de corte y pago, pago mínimo)
-- @desc Corte de una tarjeta bancaria: saldo, pago para no generar intereses, pago mínimo, intereses y comisiones. Se captura del estado de cuenta del banco o lo calcula el motor.
-- @used S3 (tarjeta: corte, pago, mínimo, intereses en cero si se paga el total a tiempo). Pantalla Créditos.
CREATE TABLE credit_statement (
    id INTEGER PRIMARY KEY,
    credit_id INTEGER NOT NULL REFERENCES credit (id) ON DELETE CASCADE, -- Tarjeta
    statement_date TEXT NOT NULL CHECK (date(statement_date) IS statement_date), -- Fecha de corte | fecha válida
    due_date TEXT NOT NULL CHECK (date(due_date) IS due_date), -- Fecha límite de pago | no antes del corte
    balance_cents INTEGER NOT NULL CHECK (balance_cents >= 0), -- Saldo al corte | centavos
    no_interest_payment_cents INTEGER CHECK (no_interest_payment_cents >= 0), -- Pago para no generar intereses | centavos
    minimum_payment_cents INTEGER CHECK (minimum_payment_cents >= 0), -- Pago mínimo | centavos
    interest_cents INTEGER NOT NULL DEFAULT 0 CHECK (interest_cents >= 0), -- Intereses cobrados en el periodo, con IVA | centavos
    fees_cents INTEGER NOT NULL DEFAULT 0 CHECK (fees_cents >= 0), -- Comisiones cobradas en el periodo | centavos
    source TEXT NOT NULL CHECK (source IN ('bank', 'computed')), -- bank = capturado del estado de cuenta, computed = calculado por el motor | enumerado
    attachment_id INTEGER REFERENCES attachment (id) ON DELETE SET NULL, -- PDF del estado de cuenta
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (due_date >= statement_date), -- @rule La fecha límite no es anterior al corte.
    UNIQUE (credit_id, statement_date) -- @rule Un corte por fecha y tarjeta.
) STRICT;
-- @idx Cortes por fecha límite (alerta de pago próximo).
CREATE INDEX ix_credit_statement_due ON credit_statement (due_date);
-- @idx Llave foránea hacia el PDF.
CREATE INDEX ix_credit_statement_attachment ON credit_statement (attachment_id);

-- @table credit_payment
-- @title Pago hecho
-- @entity pago_programado (pago real) y decisión A: Karen paga directo y repone la línea cuando paga
-- @desc Pago realmente hecho a un crédito: quién pagó, cuándo y cuánto, y qué pago programado o corte cubre. Con esto se sabe si Karen ya pagó su parte.
-- @used S3 (reposición de capital, saldo vivo), S7 (estado de cuenta de Karen, alerta de pago). Pantalla Créditos.
CREATE TABLE credit_payment (
    id INTEGER PRIMARY KEY,
    credit_id INTEGER NOT NULL REFERENCES credit (id) ON DELETE RESTRICT, -- Crédito pagado
    scheduled_payment_id INTEGER REFERENCES scheduled_payment (id) ON DELETE RESTRICT, -- Pago programado que cubre. NULL = abono libre (tarjeta)
    credit_statement_id INTEGER REFERENCES credit_statement (id) ON DELETE RESTRICT, -- Corte de tarjeta que paga
    payer_person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Quién pagó (Karen paga directo su parte)
    paid_date TEXT NOT NULL CHECK (date(paid_date) IS paid_date), -- Fecha del pago | fecha válida
    amount_cents INTEGER NOT NULL CHECK (amount_cents > 0), -- Monto pagado en centavos | mayor que 0
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Pagos de un crédito por fecha.
CREATE INDEX ix_credit_payment_credit_date ON credit_payment (credit_id, paid_date);
-- @idx Pagos que cubren un pago programado (¿ya se pagó?).
CREATE INDEX ix_credit_payment_scheduled ON credit_payment (scheduled_payment_id);
-- @idx Llave foránea hacia el corte de tarjeta.
CREATE INDEX ix_credit_payment_statement ON credit_payment (credit_statement_id);
-- @idx Pagos hechos por una persona (Karen) por fecha.
CREATE INDEX ix_credit_payment_payer_date ON credit_payment (payer_person_id, paid_date);

-- ---------------------------------------------------------------------
-- B (cont.). Ahorro
-- ---------------------------------------------------------------------

-- @table savings_goal
-- @title Meta de ahorro
-- @entity meta_ahorro
-- @desc Fondo de ahorro o meta con monto y fecha, con su aportación mensual planificada. El avance se calcula con sus movimientos.
-- @used Ajustes «Salud y metas de ahorro». S2 (aportación en el plan del mes, avance). Pantalla Inicio (tarjeta de ahorro).
CREATE TABLE savings_goal (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) > 0), -- Nombre (Fondo de emergencia) | único
    kind TEXT NOT NULL CHECK (kind IN ('fund', 'goal')), -- fund = fondo de ahorro, goal = meta con monto y fecha | enumerado
    target_cents INTEGER CHECK (target_cents > 0), -- Monto objetivo en centavos | obligatorio en goal
    target_date TEXT CHECK (date(target_date) IS target_date), -- Fecha objetivo | obligatoria en goal, posterior al inicio
    planned_monthly_cents INTEGER CHECK (planned_monthly_cents >= 0), -- Aportación mensual planificada en centavos | centavos
    start_date TEXT NOT NULL CHECK (date(start_date) IS start_date), -- Inicio del fondo o meta | fecha válida
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (kind <> 'goal' OR (target_cents IS NOT NULL AND target_date IS NOT NULL)), -- @rule Una meta tiene monto y fecha.
    CHECK (target_date IS NULL OR target_date > start_date) -- @rule La fecha objetivo es posterior al inicio.
) STRICT;

-- @table savings_movement
-- @title Movimiento de ahorro
-- @entity movimiento_ahorro
-- @desc Aportación o retiro de un fondo o meta.
-- @used S2 (avance de metas), Pantalla Inicio.
CREATE TABLE savings_movement (
    id INTEGER PRIMARY KEY,
    savings_goal_id INTEGER NOT NULL REFERENCES savings_goal (id) ON DELETE RESTRICT, -- Fondo o meta
    movement_date TEXT NOT NULL CHECK (date(movement_date) IS movement_date), -- Fecha | fecha válida
    kind TEXT NOT NULL CHECK (kind IN ('deposit', 'withdrawal')), -- deposit = aportación, withdrawal = retiro | enumerado
    amount_cents INTEGER NOT NULL CHECK (amount_cents > 0), -- Monto en centavos | mayor que 0
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Movimientos de una meta por fecha (avance).
CREATE INDEX ix_savings_movement_goal_date ON savings_movement (savings_goal_id, movement_date);
-- @idx Movimientos de un periodo.
CREATE INDEX ix_savings_movement_date ON savings_movement (movement_date);

-- ---------------------------------------------------------------------
-- D. Comida, salud e inventario
-- ---------------------------------------------------------------------

-- @table meal_slot
-- @title Tiempo de comida
-- @entity tiempo_comida
-- @desc Los cuatro tiempos del día: desayuno, comida, cena y colación (cuarto tiempo opcional con costo y calorías). Datos de referencia que crea esta migración.
-- @used S4 (plan semanal por tiempo), Pantalla Comidas, Ajustes «Recetas».
CREATE TABLE meal_slot (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL UNIQUE CHECK (code IN ('breakfast', 'lunch', 'dinner', 'snack')), -- breakfast = desayuno, lunch = comida, dinner = cena, snack = colación | enumerado
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Nombre en pantalla | no vacío
    sort_order INTEGER NOT NULL UNIQUE, -- Orden en el día | único
    is_optional INTEGER NOT NULL DEFAULT 0 CHECK (is_optional IN (0, 1)), -- 1 = tiempo opcional (colación) | booleano
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
INSERT INTO meal_slot (id, code, name, sort_order, is_optional) VALUES
    (1, 'breakfast', 'Desayuno', 1, 0),
    (2, 'lunch', 'Comida', 2, 0),
    (3, 'dinner', 'Cena', 3, 0),
    (4, 'snack', 'Colación', 4, 1);

-- @table recipe
-- @title Receta
-- @entity receta
-- @desc Receta con porciones base (2), tiempos y origen: recetario inicial escrito por Claude o recetas de Axel y Ximena. Editable.
-- @used Ajustes «Recetas». S4 (escalado, costo y kcal por receta). Pantalla Comidas.
CREATE TABLE recipe (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) > 0), -- Nombre | único
    description TEXT, -- Descripción corta
    instructions TEXT, -- Pasos de preparación
    base_servings INTEGER NOT NULL DEFAULT 2 CHECK (base_servings BETWEEN 1 AND 50), -- Porciones para las que están escritas las cantidades (base 2) | 1 a 50
    prep_minutes INTEGER CHECK (prep_minutes >= 0), -- Minutos de preparación | mayor o igual a 0
    cook_minutes INTEGER CHECK (cook_minutes >= 0), -- Minutos de cocción | mayor o igual a 0
    origin TEXT NOT NULL CHECK (origin IN ('claude_seed', 'household')), -- claude_seed = recetario inicial escrito por Claude, household = de Axel o Ximena | enumerado
    author_person_id INTEGER REFERENCES person (id) ON DELETE SET NULL, -- Quién la dio de alta (en household)
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Llave foránea hacia el autor.
CREATE INDEX ix_recipe_author ON recipe (author_person_id);

-- @table recipe_meal_slot
-- @title Tiempo apto de una receta
-- @entity receta (detalle: tiempos aptos)
-- @desc Tiempos de comida para los que sirve una receta (desayuno, comida, cena o colación). Una receta puede servir para varios.
-- @used Pantalla Comidas (selector de recetas de una casilla, N-122). Ajustes «Recetas» (lista con filtro por tiempo N-135, crear receta N-136). S4 (plan semanal).
CREATE TABLE recipe_meal_slot (
    id INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipe (id) ON DELETE CASCADE, -- Receta
    meal_slot_id INTEGER NOT NULL REFERENCES meal_slot (id) ON DELETE RESTRICT, -- Tiempo para el que sirve
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (recipe_id, meal_slot_id) -- @rule Un tiempo aparece una sola vez por receta.
) STRICT;
-- @idx Recetas aptas para un tiempo (selector de la casilla).
CREATE INDEX ix_recipe_meal_slot_slot ON recipe_meal_slot (meal_slot_id, recipe_id);

-- @table recipe_ingredient
-- @title Ingrediente de receta
-- @entity ingrediente_receta
-- @desc Ingrediente exacto de una receta: producto, presentación comprable sugerida y cantidad para las porciones base en milésimas de la unidad base.
-- @used S4 (escalado exacto, lista de compras, kcal por receta), S5 (costo por ingrediente).
CREATE TABLE recipe_ingredient (
    id INTEGER PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipe (id) ON DELETE CASCADE, -- Receta
    product_id INTEGER NOT NULL REFERENCES product (id) ON DELETE RESTRICT, -- Ingrediente (producto genérico)
    presentation_id INTEGER REFERENCES presentation (id) ON DELETE SET NULL, -- Presentación comprable sugerida
    quantity_milli INTEGER NOT NULL CHECK (quantity_milli > 0), -- Cantidad para base_servings en milésimas de la unidad base (250 g = 250000, media pieza = 500) | mayor que 0
    display_text TEXT, -- Cantidad como la escribe la receta (1/2 taza)
    is_optional INTEGER NOT NULL DEFAULT 0 CHECK (is_optional IN (0, 1)), -- 1 = ingrediente opcional | booleano
    sort_order INTEGER NOT NULL DEFAULT 0, -- Orden en la receta
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (recipe_id, product_id) -- @rule Un producto aparece una vez por receta.
) STRICT;
-- @idx Recetas que usan un producto.
CREATE INDEX ix_recipe_ingredient_product ON recipe_ingredient (product_id);
-- @idx Llave foránea hacia la presentación.
CREATE INDEX ix_recipe_ingredient_presentation ON recipe_ingredient (presentation_id);

-- @table meal_plan_entry
-- @title Entrada del plan de comidas
-- @entity plan_comida
-- @desc Un renglón del plan: día × tiempo × persona × receta × porciones, o «fuera» (sin receta ni calorías). Solo Axel y Ximena.
-- @used S4 (plan semanal, lista de compras, costo y kcal por día), S7 (alerta de plan sin armar). Pantalla Comidas.
CREATE TABLE meal_plan_entry (
    id INTEGER PRIMARY KEY,
    plan_date TEXT NOT NULL CHECK (date(plan_date) IS plan_date), -- Día del plan | fecha válida
    meal_slot_id INTEGER NOT NULL REFERENCES meal_slot (id) ON DELETE RESTRICT, -- Tiempo de comida
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Axel o Ximena | nunca una persona externa (disparador)
    recipe_id INTEGER REFERENCES recipe (id) ON DELETE RESTRICT, -- Receta. NULL si come fuera
    servings_milli INTEGER CHECK (servings_milli > 0), -- Porciones que come la persona ×1000 (1 porción = 1000) | NULL si come fuera
    is_eating_out INTEGER NOT NULL DEFAULT 0 CHECK (is_eating_out IN (0, 1)), -- 1 = come fuera: sin receta ni calorías, el gasto va a la categoría «comida fuera» | booleano
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((is_eating_out = 1 AND recipe_id IS NULL AND servings_milli IS NULL) OR (is_eating_out = 0 AND recipe_id IS NOT NULL AND servings_milli IS NOT NULL)) -- @rule Con receta lleva porciones. Fuera no lleva receta ni porciones.
) STRICT;
-- @idx Una receta una sola vez por día, tiempo y persona (y un solo «fuera»).
CREATE UNIQUE INDEX ux_meal_plan_entry_slot ON meal_plan_entry (plan_date, meal_slot_id, person_id, coalesce(recipe_id, 0));
-- @idx Plan de una persona por fecha (kcal por día).
CREATE INDEX ix_meal_plan_entry_person_date ON meal_plan_entry (person_id, plan_date);
-- @idx Plan por fecha (semana, lista de compras).
CREATE INDEX ix_meal_plan_entry_date ON meal_plan_entry (plan_date);
-- @idx Días en que se planeó una receta.
CREATE INDEX ix_meal_plan_entry_recipe ON meal_plan_entry (recipe_id);
-- @idx Llave foránea hacia el tiempo de comida.
CREATE INDEX ix_meal_plan_entry_slot ON meal_plan_entry (meal_slot_id);

-- @table inventory_item
-- @title Renglón de inventario
-- @entity inventario_item
-- @desc Inventario híbrido: existencia por unidad (básicos, no perecederos y sobrantes de recetas) o ritmo de consumo (el resto).
-- @used Ajustes «Inventario». S4 (lista = ingredientes − existencia, sobrantes). Pantalla Compras y tiendas.
CREATE TABLE inventory_item (
    id INTEGER PRIMARY KEY,
    product_id INTEGER REFERENCES product (id) ON DELETE RESTRICT, -- Producto en existencia | exactamente uno de product_id o recipe_id
    recipe_id INTEGER REFERENCES recipe (id) ON DELETE RESTRICT, -- Sobrante de una receta (en porciones)
    tracking_mode TEXT NOT NULL CHECK (tracking_mode IN ('stock', 'consumption_rate')), -- stock = existencia por unidad, consumption_rate = ritmo de consumo | enumerado
    quantity_milli INTEGER NOT NULL DEFAULT 0 CHECK (quantity_milli >= 0), -- Existencia en milésimas de la unidad base (en sobrantes, porciones ×1000): la vigente en stock y la del último conteo en ritmo de consumo. Cada cambio deja su inventory_movement | mayor o igual a 0
    counted_date TEXT CHECK (date(counted_date) IS counted_date), -- Fecha del último conteo, base de la estimación por ritmo | obligatoria en consumption_rate
    rate_quantity_milli INTEGER CHECK (rate_quantity_milli > 0), -- Ritmo: cantidad que se consume (milésimas de la unidad base) | solo en consumption_rate
    rate_period_days INTEGER CHECK (rate_period_days > 0), -- Ritmo: en cuántos días se consume esa cantidad | solo en consumption_rate
    reorder_level_milli INTEGER CHECK (reorder_level_milli >= 0), -- Nivel que pone el producto en la lista de compras. NULL = sin aviso
    expires_date TEXT CHECK (date(expires_date) IS expires_date), -- Caducidad (sobrantes y perecederos) | fecha válida
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((product_id IS NULL) <> (recipe_id IS NULL)), -- @rule Cada renglón es de un producto o de un sobrante de receta, nunca de ambos.
    CHECK ((tracking_mode = 'stock' AND rate_quantity_milli IS NULL AND rate_period_days IS NULL) OR (tracking_mode = 'consumption_rate' AND rate_quantity_milli IS NOT NULL AND rate_period_days IS NOT NULL AND counted_date IS NOT NULL)), -- @rule El ritmo de consumo lleva cantidad, días y fecha de conteo. La existencia por unidad no lleva ritmo.
    CHECK (recipe_id IS NULL OR tracking_mode = 'stock') -- @rule Los sobrantes se llevan por unidad (porciones).
) STRICT;
-- @idx Un solo renglón de inventario por producto.
CREATE UNIQUE INDEX ux_inventory_item_product ON inventory_item (product_id) WHERE product_id IS NOT NULL;
-- @idx Sobrantes de una receta.
CREATE INDEX ix_inventory_item_recipe ON inventory_item (recipe_id);

-- @table inventory_movement
-- @title Movimiento de inventario
-- @entity inventario_item (detalle: movimientos)
-- @desc Cada entrada o salida de un renglón de inventario: compra, sobrante de receta, consumo del plan o ajuste a mano con su motivo (conteo, consumo o merma). La API lo escribe en la misma transacción que cambia la existencia.
-- @used Pantalla Compras y tiendas › Inventario (movimientos N-181, ajustar existencia N-183). S4 (existencia, consumo del plan y sobrantes).
CREATE TABLE inventory_movement (
    id INTEGER PRIMARY KEY,
    inventory_item_id INTEGER NOT NULL REFERENCES inventory_item (id) ON DELETE CASCADE, -- Renglón de inventario
    movement_date TEXT NOT NULL CHECK (date(movement_date) IS movement_date), -- Fecha del movimiento | fecha válida
    kind TEXT NOT NULL CHECK (kind IN ('purchase', 'recipe_leftover', 'plan_consumption', 'adjustment')), -- purchase = entrada por compra, recipe_leftover = sobrante de receta, plan_consumption = consumo del plan, adjustment = ajuste a mano (N-181) | enumerado
    reason TEXT CHECK (reason IN ('count', 'consumption', 'waste')), -- Motivo del ajuste: count = conteo, consumption = consumo, waste = merma (N-183) | obligatorio solo en adjustment
    quantity_delta_milli INTEGER NOT NULL, -- Diferencia en milésimas de la unidad base (en sobrantes, porciones ×1000): positiva entra, negativa sale. Un ajuste por «cantidad nueva» guarda la diferencia | entero con signo
    expense_item_id INTEGER REFERENCES expense_item (id) ON DELETE SET NULL, -- Renglón del ticket que le dio origen (compra)
    meal_plan_entry_id INTEGER REFERENCES meal_plan_entry (id) ON DELETE SET NULL, -- Entrada del plan que le dio origen (consumo o sobrante)
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((kind = 'adjustment') = (reason IS NOT NULL)), -- @rule Solo el ajuste lleva motivo, y siempre lo lleva.
    CHECK (expense_item_id IS NULL OR kind = 'purchase'), -- @rule El renglón de ticket solo acompaña a una compra.
    CHECK (meal_plan_entry_id IS NULL OR kind IN ('recipe_leftover', 'plan_consumption')) -- @rule La entrada del plan solo acompaña a consumos y sobrantes.
) STRICT;
-- @idx Movimientos de un renglón por fecha (N-181).
CREATE INDEX ix_inventory_movement_item_date ON inventory_movement (inventory_item_id, movement_date);
-- @idx Llave foránea hacia el renglón del ticket.
CREATE INDEX ix_inventory_movement_expense_item ON inventory_movement (expense_item_id);
-- @idx Llave foránea hacia la entrada del plan.
CREATE INDEX ix_inventory_movement_meal_plan ON inventory_movement (meal_plan_entry_id);

-- @table shopping_list
-- @title Lista de compras
-- @entity (complementaria) lista de compras de la semana con marcas de comprado y artículos a mano (U1 N-150 a N-153 y N-165), compartida en el servidor entre los dos teléfonos (K1)
-- @desc Lista de compras de una semana del plan de comidas. Guarda lo que no se puede calcular: los artículos agregados a mano, la tienda o cantidad elegida, las marcas de comprado y cuándo se generó. Lo calculado (necesidad del plan − existencia + reposición, precios y frescura) lo dan S4 y S5. Se comparte en el servidor entre los dos teléfonos (K1).
-- @used Pantalla Compras y tiendas › Lista de compras (N-150 a N-153) y comparativa (N-165 «Usar esta tienda»). Comidas (N-127 «Generar lista de compras»). Gastos (N-073 «Terminar compra»). S4 (lista neteada), S5 (tienda recomendada).
CREATE TABLE shopping_list (
    id INTEGER PRIMARY KEY,
    start_date TEXT NOT NULL UNIQUE CHECK (date(start_date) IS start_date), -- Lunes que abre la semana del plan (K8: semana de lunes a domingo) | fecha válida y en lunes (lo exigen trg_shopping_list_monday_ins y trg_shopping_list_monday_upd, de la 00002), una lista por semana
    end_date TEXT NOT NULL CHECK (date(end_date) IS end_date), -- Último día de la semana, inclusive | fecha válida, start_date + 6 días
    generated_at TEXT, -- Última vez que S4 la generó o actualizó desde el plan (N-127), UTC. NULL = todavía no se genera
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (end_date = date(start_date, '+6 days')) -- @rule La lista cubre exactamente una semana de 7 días.
) STRICT;

-- @table shopping_list_item
-- @title Artículo de la lista de compras
-- @entity (complementaria) artículos de la lista de compras, calculados por S4 o agregados a mano, con tienda, cantidad y marca de comprado, compartidos en el servidor (K1)
-- @desc Un artículo de la lista: un producto que calculó S4 al generar la lista o un artículo agregado a mano (producto del catálogo o texto libre). Guarda la tienda o mercado elegido, la cantidad elegida y si ya se compró. Se comparte en el servidor entre los dos teléfonos (K1).
-- @used Pantalla Compras y tiendas › Lista de compras (marcar comprado N-151, agregar artículo N-152, cambiar tienda o cantidad N-153) y comparativa (N-165). Gastos (prellenado N-073). S4 (lista neteada), S5 (mejor tienda y canasta, N-024).
CREATE TABLE shopping_list_item (
    id INTEGER PRIMARY KEY,
    shopping_list_id INTEGER NOT NULL REFERENCES shopping_list (id) ON DELETE CASCADE, -- Lista a la que pertenece
    origin TEXT NOT NULL CHECK (origin IN ('calculated', 'manual')), -- calculated = lo calculó S4 (recetas del plan o reposición por ritmo), manual = lo agregó una persona (N-152) | enumerado
    product_id INTEGER REFERENCES product (id) ON DELETE RESTRICT, -- Producto del catálogo | exactamente uno de product_id o free_text
    free_text TEXT CHECK (free_text IS NULL OR length(trim(free_text)) > 0), -- Artículo escrito a mano que no está en el catálogo | solo en manual; si viene, no vacío
    free_text_unit TEXT CHECK (free_text_unit IN ('g', 'ml', 'piece')), -- Unidad de la cantidad de un artículo en texto libre (un producto usa su product.base_unit) | g, ml o piece; obligatoria con free_text
    presentation_id INTEGER REFERENCES presentation (id) ON DELETE SET NULL, -- Presentación elegida para comprar. NULL = la que recomienda S5
    needed_milli INTEGER CHECK (needed_milli >= 0), -- Necesidad que calculó S4 en la última generación (plan − existencia + reposición), en milésimas de la unidad base | solo en calculated
    quantity_milli INTEGER CHECK (quantity_milli > 0), -- Cantidad a comprar que fijó una persona (N-152, N-153), en milésimas de la unidad base. NULL en calculated = comprar la necesidad | obligatoria en manual
    store_id INTEGER REFERENCES store (id) ON DELETE SET NULL, -- Tienda o mercado elegido (N-153, N-165). NULL = el que recomienda S5
    is_bought INTEGER NOT NULL DEFAULT 0 CHECK (is_bought IN (0, 1)), -- 1 = comprado (N-151). No crea gasto ni precio hasta confirmar el ticket | booleano
    bought_at TEXT, -- Cuándo se marcó como comprado (UTC) | solo con is_bought = 1
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK ((product_id IS NULL) <> (free_text IS NULL)), -- @rule Cada artículo es un producto del catálogo o un texto libre, nunca ambos.
    CHECK ((free_text IS NULL) = (free_text_unit IS NULL)), -- @rule El texto libre lleva su unidad; un producto usa su unidad base.
    CHECK (free_text IS NULL OR presentation_id IS NULL), -- @rule Un texto libre no tiene presentación del catálogo.
    CHECK ((origin = 'calculated' AND product_id IS NOT NULL AND needed_milli IS NOT NULL) OR (origin = 'manual' AND needed_milli IS NULL AND quantity_milli IS NOT NULL)), -- @rule Lo calculado es un producto con su necesidad; lo manual lleva cantidad y no tiene necesidad calculada.
    CHECK (is_bought = 1 OR bought_at IS NULL) -- @rule Solo un artículo comprado lleva la hora de compra.
) STRICT;
-- @idx Un artículo calculado por producto en cada lista (la regeneración lo actualiza en su lugar).
CREATE UNIQUE INDEX ux_shopping_list_item_calculated ON shopping_list_item (shopping_list_id, product_id) WHERE origin = 'calculated';
-- @idx Artículos de una lista, comprados o no.
CREATE INDEX ix_shopping_list_item_list ON shopping_list_item (shopping_list_id, is_bought);
-- @idx Llave foránea hacia el producto.
CREATE INDEX ix_shopping_list_item_product ON shopping_list_item (product_id);
-- @idx Llave foránea hacia la presentación.
CREATE INDEX ix_shopping_list_item_presentation ON shopping_list_item (presentation_id);
-- @idx Artículos asignados a una tienda (lista agrupada por tienda, N-165).
CREATE INDEX ix_shopping_list_item_store ON shopping_list_item (store_id);

-- @table health_profile
-- @title Perfil de salud
-- @entity perfil_salud
-- @desc Perfil de salud por persona (solo Axel y Ximena): estatura, nacimiento, sexo, actividad, meta de calorías y restricciones (vacías por decisión). El peso vive en weight_entry.
-- @used Ajustes «Salud y metas de ahorro». S4 (IMC, kcal por día contra la meta). Pantalla Salud.
CREATE TABLE health_profile (
    id INTEGER PRIMARY KEY,
    person_id INTEGER NOT NULL UNIQUE REFERENCES person (id) ON DELETE RESTRICT, -- Persona | una por persona, nunca externa (disparador)
    height_mm INTEGER CHECK (height_mm BETWEEN 500 AND 2500), -- Estatura en milímetros (1.75 m = 1750) | 500 a 2500
    birth_date TEXT CHECK (date(birth_date) IS birth_date), -- Fecha de nacimiento. La edad se calcula | fecha válida
    sex TEXT CHECK (sex IN ('female', 'male')), -- Sexo para fórmulas de requerimiento | enumerado
    activity_level TEXT CHECK (activity_level IN ('sedentary', 'light', 'moderate', 'active', 'very_active')), -- Nivel de actividad de la escala de 5 niveles que multiplica a Mifflin-St Jeor: sedentary 1.2, light 1.375, moderate 1.55, active 1.725, very_active 1.9 (K23) | enumerado
    daily_kcal_target INTEGER CHECK (daily_kcal_target BETWEEN 500 AND 10000), -- Meta de kilocalorías por día | 500 a 10000
    dietary_restrictions TEXT, -- Restricciones alimentarias (el campo existe y hoy va vacío)
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;

-- @table weight_entry
-- @title Peso histórico
-- @entity peso_historico
-- @desc Historial de peso por persona. El más reciente alimenta el IMC.
-- @used S4 (IMC), Pantalla Salud (tendencia).
CREATE TABLE weight_entry (
    id INTEGER PRIMARY KEY,
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona | nunca externa (disparador)
    measured_date TEXT NOT NULL CHECK (date(measured_date) IS measured_date), -- Fecha de la medición | fecha válida
    weight_g INTEGER NOT NULL CHECK (weight_g BETWEEN 20000 AND 400000), -- Peso en gramos (72.5 kg = 72500) | 20 a 400 kg
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (person_id, measured_date) -- @rule Una medición por persona y día.
) STRICT;

-- @table nutrition_fact
-- @title Dato nutricional
-- @entity nutricion_ingrediente
-- @desc Kilocalorías por 100 g de un producto (o de una presentación específica) con fuente citada. Macronutrientes previstos, se llenan después.
-- @used S4 (kcal por receta y por día), carga local de USDA FoodData Central y Open Food Facts.
CREATE TABLE nutrition_fact (
    id INTEGER PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES product (id) ON DELETE CASCADE, -- Producto
    presentation_id INTEGER REFERENCES presentation (id) ON DELETE CASCADE, -- Presentación específica (etiqueta, Open Food Facts). NULL = valor genérico del producto
    kcal_milli_per_100g INTEGER NOT NULL CHECK (kcal_milli_per_100g >= 0), -- Milésimas de kilocaloría por 100 g (884 kcal = 884000) | mayor o igual a 0
    protein_mg_per_100g INTEGER CHECK (protein_mg_per_100g >= 0), -- Proteína en mg por 100 g (previsto) | mayor o igual a 0
    fat_mg_per_100g INTEGER CHECK (fat_mg_per_100g >= 0), -- Grasa en mg por 100 g (previsto) | mayor o igual a 0
    carbs_mg_per_100g INTEGER CHECK (carbs_mg_per_100g >= 0), -- Carbohidratos en mg por 100 g (previsto) | mayor o igual a 0
    fiber_mg_per_100g INTEGER CHECK (fiber_mg_per_100g >= 0), -- Fibra en mg por 100 g (previsto) | mayor o igual a 0
    source TEXT NOT NULL CHECK (source IN ('usda_fdc', 'open_food_facts', 'mx_table', 'store_label', 'manual')), -- Fuente: USDA FoodData Central, Open Food Facts, tabla mexicana citada, etiqueta de tienda o captura manual | enumerado
    source_ref TEXT, -- Cita: fdc_id, código de barras, tabla y página o URL de la etiqueta | obligatoria salvo en manual
    source_license TEXT, -- Licencia de la fuente (CC0, ODbL)
    is_preferred INTEGER NOT NULL DEFAULT 0 CHECK (is_preferred IN (0, 1)), -- 1 = valor que usan los cálculos | booleano, uno por producto o presentación
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (source = 'manual' OR source_ref IS NOT NULL) -- @rule Toda fuente externa lleva su cita.
) STRICT;
-- @idx Un solo valor preferido por producto o presentación.
CREATE UNIQUE INDEX ux_nutrition_fact_preferred ON nutrition_fact (product_id, coalesce(presentation_id, 0)) WHERE is_preferred = 1;
-- @idx Datos de un producto.
CREATE INDEX ix_nutrition_fact_product ON nutrition_fact (product_id, source);
-- @idx Llave foránea hacia la presentación.
CREATE INDEX ix_nutrition_fact_presentation ON nutrition_fact (presentation_id);

-- ---------------------------------------------------------------------
-- F. Avisos, jobs, respaldos y reportes
-- ---------------------------------------------------------------------

-- @table notification
-- @title Aviso
-- @entity notificacion
-- @desc Cada aviso o reporte configurable: resúmenes diario y semanal, reporte mensual, pago próximo, presupuesto 80/100 %, precios viejos y plan sin armar. Esta migración crea los siete con los parámetros que ya decidió Axel.
-- @used Ajustes «Notificaciones» y «Reportes y frecuencias». S7 (reportes y alertas).
CREATE TABLE notification (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL UNIQUE CHECK (code IN ('daily_summary', 'weekly_summary', 'monthly_report', 'payment_due', 'budget_threshold', 'stale_prices', 'meal_plan_missing')), -- Clave del aviso | enumerado
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Nombre en Ajustes | no vacío
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)), -- 1 = activo | booleano
    params_json TEXT NOT NULL DEFAULT '{}' CHECK (json_valid(params_json)), -- Parámetros del aviso (días antes 1 y 0, umbrales 80 y 100) | JSON válido
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
INSERT INTO notification (code, name, params_json) VALUES
    ('daily_summary', 'Resumen diario', '{}'),
    ('weekly_summary', 'Resumen semanal', '{"weekday": 7}'),
    ('monthly_report', 'Reporte mensual en PDF', '{"day_of_month": 1}'),
    ('payment_due', 'Pago próximo', '{"days_before": [1, 0]}'),
    ('budget_threshold', 'Presupuesto al 80 % y al 100 %', '{"thresholds_pct": [80, 100], "scopes": ["category", "total"]}'),
    ('stale_prices', 'Precios viejos antes del día de despensa', '{}'),
    ('meal_plan_missing', 'Plan de comidas sin armar', '{"weekday": 7}');

-- @table notification_recipient
-- @title Destinatario
-- @entity destinatario
-- @desc Destino de avisos de una persona: su chat de Telegram en el bot MultiX. El chat_id no es credencial. El token del bot nunca se guarda en la base.
-- @used Ajustes «Notificaciones». S7 (envío de avisos a Axel y Ximena).
CREATE TABLE notification_recipient (
    id INTEGER PRIMARY KEY,
    person_id INTEGER NOT NULL REFERENCES person (id) ON DELETE RESTRICT, -- Persona destinataria
    channel TEXT NOT NULL DEFAULT 'telegram' CHECK (channel IN ('telegram')), -- Canal (bot MultiX) | enumerado
    chat_id TEXT NOT NULL CHECK (length(chat_id) BETWEEN 1 AND 32 AND chat_id NOT GLOB '*[^0-9-]*'), -- chat_id de Telegram | dígitos y signo menos
    label TEXT, -- Etiqueta (Celular de Ximena)
    verified_at TEXT, -- Cuándo se comprobó con un mensaje de prueba (UTC)
    is_active INTEGER NOT NULL DEFAULT 1 CHECK (is_active IN (0, 1)),
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (channel, chat_id) -- @rule Un chat se registra una sola vez.
) STRICT;
-- @idx Destinos de una persona.
CREATE INDEX ix_notification_recipient_person ON notification_recipient (person_id);

-- @table notification_subscription
-- @title Suscripción a un aviso
-- @entity destinatario (destinatarios configurables por aviso)
-- @desc Qué destinatario recibe qué aviso. Se configura por aviso en Ajustes.
-- @used Ajustes «Notificaciones». S7.
CREATE TABLE notification_subscription (
    id INTEGER PRIMARY KEY,
    notification_id INTEGER NOT NULL REFERENCES notification (id) ON DELETE CASCADE, -- Aviso
    recipient_id INTEGER NOT NULL REFERENCES notification_recipient (id) ON DELETE CASCADE, -- Destinatario
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (notification_id, recipient_id) -- @rule Una suscripción por aviso y destinatario.
) STRICT;
-- @idx Avisos de un destinatario.
CREATE INDEX ix_notification_subscription_recipient ON notification_subscription (recipient_id);

-- @table notification_delivery
-- @title Envío de aviso
-- @entity notificacion (bitácora de envíos)
-- @desc Cada envío de un aviso a un destinatario, con clave para no mandar dos veces lo mismo y el resultado HTTP.
-- @used S7 (deduplicar alertas, evidencia de envío HTTP 200). Ajustes «Notificaciones» (historial).
CREATE TABLE notification_delivery (
    id INTEGER PRIMARY KEY,
    notification_id INTEGER NOT NULL REFERENCES notification (id) ON DELETE RESTRICT, -- Aviso
    recipient_id INTEGER NOT NULL REFERENCES notification_recipient (id) ON DELETE RESTRICT, -- Destinatario
    dedupe_key TEXT NOT NULL CHECK (length(dedupe_key) > 0), -- Clave que evita repetir el mismo aviso (payment_due:42:2026-10-29) | no vacía
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed', 'skipped')), -- Estado del envío | enumerado
    message_text TEXT, -- Texto plano enviado (sin parse_mode)
    http_status INTEGER, -- Código HTTP de la API de Telegram
    error TEXT, -- Error resumido, sin token ni datos sensibles
    sent_at TEXT, -- Momento del envío (UTC)
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE (dedupe_key, recipient_id) -- @rule El mismo aviso no se manda dos veces al mismo destinatario.
) STRICT;
-- @idx Cola de envíos pendientes o fallidos.
CREATE INDEX ix_notification_delivery_status ON notification_delivery (status, created_at);
-- @idx Historial de un aviso.
CREATE INDEX ix_notification_delivery_notification ON notification_delivery (notification_id, sent_at);
-- @idx Llave foránea hacia el destinatario.
CREATE INDEX ix_notification_delivery_recipient ON notification_delivery (recipient_id);

-- @table job
-- @title Job programado
-- @entity job
-- @desc Tarea programada del servicio (respaldos, resúmenes, alertas, purgas) con su horario y su último resultado.
-- @used S6 (planificador), S7 (reportes y alertas). Ajustes «Reportes y frecuencias».
CREATE TABLE job (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL UNIQUE CHECK (length(code) BETWEEN 1 AND 40 AND code NOT GLOB '*[^a-z0-9_]*'), -- Clave del job (backup_daily, daily_summary) | a-z, 0-9 y guion bajo
    name TEXT NOT NULL CHECK (length(trim(name)) > 0), -- Nombre visible | no vacío
    schedule TEXT NOT NULL CHECK (length(trim(schedule)) > 0), -- Expresión cron de 5 campos en hora de America/Mexico_City | no vacía
    is_enabled INTEGER NOT NULL DEFAULT 1 CHECK (is_enabled IN (0, 1)), -- 1 = activo | booleano
    notification_id INTEGER REFERENCES notification (id) ON DELETE SET NULL, -- Aviso que manda, si manda alguno
    last_run_at TEXT, -- Inicio de la última corrida (UTC)
    last_status TEXT CHECK (last_status IN ('running', 'ok', 'error', 'skipped')), -- Resultado de la última corrida | enumerado
    last_error TEXT, -- Error resumido de la última corrida
    next_run_at TEXT, -- Próxima corrida calculada (UTC)
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
) STRICT;
-- @idx Siguiente job por correr.
CREATE INDEX ix_job_next_run ON job (is_enabled, next_run_at);
-- @idx Llave foránea hacia el aviso.
CREATE INDEX ix_job_notification ON job (notification_id);

-- @table backup
-- @title Respaldo
-- @entity respaldo
-- @desc Respaldo local de la base (VACUUM INTO) con huella, resultado de integrity_check y prueba de restauración.
-- @used S6 (respaldo diario, restauración probada), S8 (operación). Ajustes «Datos y respaldos».
CREATE TABLE backup (
    id INTEGER PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('scheduled', 'manual', 'pre_migration')), -- Diario programado, manual o previo a una migración | enumerado
    status TEXT NOT NULL DEFAULT 'running' CHECK (status IN ('running', 'ok', 'failed')), -- Estado | enumerado
    started_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')), -- Inicio (UTC)
    finished_at TEXT, -- Fin (UTC)
    file_path TEXT, -- Ruta del archivo en la carpeta o disco de respaldos | obligatoria si ok
    size_bytes INTEGER CHECK (size_bytes >= 0), -- Tamaño en bytes | mayor o igual a 0
    sha256 TEXT CHECK (sha256 IS NULL OR length(sha256) = 64), -- Huella del archivo | 64 caracteres, obligatoria si ok
    integrity_check TEXT, -- Resultado de PRAGMA integrity_check sobre la copia (ok o el error) | ok si status = ok
    restore_tested_at TEXT, -- Última restauración de prueba (UTC)
    restore_test_result TEXT, -- Resultado de esa restauración
    pruned_at TEXT, -- Cuándo se borró el archivo por la política de retención (UTC)
    error TEXT, -- Error resumido
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (status <> 'ok' OR (file_path IS NOT NULL AND sha256 IS NOT NULL AND integrity_check = 'ok')) -- @rule Un respaldo correcto tiene archivo, huella e integrity_check = ok.
) STRICT;
-- @idx Respaldos por fecha.
CREATE INDEX ix_backup_started ON backup (started_at);
-- @idx Último respaldo correcto o fallido.
CREATE INDEX ix_backup_status ON backup (status, started_at);

-- @table generated_report
-- @title Reporte generado
-- @entity (complementaria) decisión F: el reporte mensual PDF se envía y se guarda; estado de cuenta de Karen en PDF bajo demanda
-- @desc PDF generado por el sistema: reporte mensual LaTeX o estado de cuenta de la participante externa (Karen), con su archivo y huella.
-- @used S7 (reporte mensual día 1, estado de cuenta de Karen). Pantalla Reportes.
CREATE TABLE generated_report (
    id INTEGER PRIMARY KEY,
    kind TEXT NOT NULL CHECK (kind IN ('monthly_pdf', 'external_statement_pdf')), -- monthly_pdf = reporte mensual, external_statement_pdf = estado de cuenta de la participante externa | enumerado
    period_start_date TEXT NOT NULL CHECK (date(period_start_date) IS period_start_date), -- Inicio del periodo que cubre | fecha válida
    period_end_date TEXT NOT NULL CHECK (date(period_end_date) IS period_end_date), -- Fin del periodo que cubre, inclusive | no antes del inicio
    person_id INTEGER REFERENCES person (id) ON DELETE RESTRICT, -- Persona del estado de cuenta (Karen) | solo en external_statement_pdf
    file_path TEXT NOT NULL UNIQUE, -- Ruta del PDF guardado | única
    sha256 TEXT NOT NULL CHECK (length(sha256) = 64), -- Huella del PDF | 64 caracteres
    size_bytes INTEGER NOT NULL CHECK (size_bytes >= 0), -- Tamaño en bytes | mayor o igual a 0
    generated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')), -- Momento de generación (UTC)
    sent_at TEXT, -- Envío por Telegram (reporte mensual) (UTC)
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    CHECK (period_end_date >= period_start_date), -- @rule El periodo no termina antes de empezar.
    CHECK ((kind = 'external_statement_pdf') = (person_id IS NOT NULL)) -- @rule Solo el estado de cuenta lleva persona.
) STRICT;
-- @idx Reportes por tipo y periodo.
CREATE INDEX ix_generated_report_kind_period ON generated_report (kind, period_start_date);
-- @idx Estados de cuenta de una persona.
CREATE INDEX ix_generated_report_person ON generated_report (person_id);

-- ---------------------------------------------------------------------
-- Ajustes del hogar ya decididos por Axel (sin valores provisionales)
-- ---------------------------------------------------------------------
INSERT INTO app_setting (setting_key, value_json, description) VALUES
    ('price.stale_after_days', '30', 'Un precio es viejo cuando tiene esta edad en días o más (decisión E).'),
    ('auth.max_failed_attempts', '5', 'Intentos fallidos permitidos dentro de la ventana antes de bloquear (decisión A).'),
    ('auth.lockout_window_minutes', '15', 'Ventana en minutos para contar intentos fallidos (decisión A).'),
    ('optimizer.weights', '{"lowest_interest": 1, "even_low_payment": 1, "max_available": 1, "most_value_per_peso": 1}', 'Pesos del puntaje del comparador de escenarios. Iguales y se normalizan (decisión C).'),
    ('ai.local_model_enabled', 'false', 'Modelo local (Ollama) para OCR y sugerencias. Desactivado por defecto (decisión H).'),
    ('ai.claude_api_enabled', 'false', 'API de Claude para OCR y sugerencias. Desactivada por defecto. La llave nunca se guarda en la base (decisión H).'),
    ('price.assisted_reading_enabled', 'false', 'Interruptor general de la lectura asistida de sitios. Además cada cadena necesita su veredicto de términos (decisión E).');

-- ---------------------------------------------------------------------
-- Disparadores de reglas de negocio (decisiones A y B)
-- ---------------------------------------------------------------------

-- Decisión A: Karen (externa) no tiene cuenta.
-- +goose StatementBegin
CREATE TRIGGER trg_user_account_not_external_ins BEFORE INSERT ON user_account
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'user_account: una persona externa no puede tener cuenta (decisión A)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_user_account_not_external_upd BEFORE UPDATE OF person_id ON user_account
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'user_account: una persona externa no puede tener cuenta (decisión A)');
END;
-- +goose StatementEnd

-- Decisión B: solo entran ingresos del titular (Axel).
-- +goose StatementBegin
CREATE TRIGGER trg_income_holder_only_ins BEFORE INSERT ON income
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS NOT 'holder'
BEGIN
    SELECT RAISE(ABORT, 'income: solo se registran ingresos del titular (decisión B)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_income_holder_only_upd BEFORE UPDATE OF person_id ON income
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS NOT 'holder'
BEGIN
    SELECT RAISE(ABORT, 'income: solo se registran ingresos del titular (decisión B)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_planned_item_income_holder_ins BEFORE INSERT ON planned_item
FOR EACH ROW WHEN NEW.kind = 'income' AND NEW.person_id IS NOT NULL
    AND (SELECT participation FROM person WHERE id = NEW.person_id) IS NOT 'holder'
BEGIN
    SELECT RAISE(ABORT, 'planned_item: un ingreso planificado solo puede ser del titular (decisión B)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_planned_item_income_holder_upd BEFORE UPDATE OF kind, person_id ON planned_item
FOR EACH ROW WHEN NEW.kind = 'income' AND NEW.person_id IS NOT NULL
    AND (SELECT participation FROM person WHERE id = NEW.person_id) IS NOT 'holder'
BEGIN
    SELECT RAISE(ABORT, 'planned_item: un ingreso planificado solo puede ser del titular (decisión B)');
END;
-- +goose StatementEnd

-- Decisión A: plan de comidas y salud solo para Axel y Ximena.
-- +goose StatementBegin
CREATE TRIGGER trg_health_profile_not_external_ins BEFORE INSERT ON health_profile
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'health_profile: salud solo para Axel y Ximena (decisión A)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_health_profile_not_external_upd BEFORE UPDATE OF person_id ON health_profile
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'health_profile: salud solo para Axel y Ximena (decisión A)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_weight_entry_not_external_ins BEFORE INSERT ON weight_entry
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'weight_entry: salud solo para Axel y Ximena (decisión A)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_weight_entry_not_external_upd BEFORE UPDATE OF person_id ON weight_entry
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'weight_entry: salud solo para Axel y Ximena (decisión A)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_meal_plan_entry_not_external_ins BEFORE INSERT ON meal_plan_entry
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'meal_plan_entry: el plan de comidas es solo para Axel y Ximena (decisión A)');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_meal_plan_entry_not_external_upd BEFORE UPDATE OF person_id ON meal_plan_entry
FOR EACH ROW WHEN (SELECT participation FROM person WHERE id = NEW.person_id) IS 'external'
BEGIN
    SELECT RAISE(ABORT, 'meal_plan_entry: el plan de comidas es solo para Axel y Ximena (decisión A)');
END;
-- +goose StatementEnd

-- Decisión A: quien tiene cuenta, salud o plan de comidas no puede pasar a externa.
-- +goose StatementBegin
CREATE TRIGGER trg_person_external_without_access BEFORE UPDATE OF participation ON person
FOR EACH ROW WHEN NEW.participation = 'external' AND OLD.participation IS NOT 'external' AND (
    EXISTS (SELECT 1 FROM user_account WHERE person_id = NEW.id)
    OR EXISTS (SELECT 1 FROM health_profile WHERE person_id = NEW.id)
    OR EXISTS (SELECT 1 FROM weight_entry WHERE person_id = NEW.id)
    OR EXISTS (SELECT 1 FROM meal_plan_entry WHERE person_id = NEW.id))
BEGIN
    SELECT RAISE(ABORT, 'person: una persona con cuenta, salud o plan de comidas no puede pasar a externa (decisión A)');
END;
-- +goose StatementEnd

-- Convención de auditoría: la bitácora no se modifica ni se borra.
-- +goose StatementBegin
CREATE TRIGGER trg_audit_log_no_update BEFORE UPDATE ON audit_log
BEGIN
    SELECT RAISE(ABORT, 'audit_log es inmutable');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_audit_log_no_delete BEFORE DELETE ON audit_log
BEGIN
    SELECT RAISE(ABORT, 'audit_log es inmutable');
END;
-- +goose StatementEnd

-- ---------------------------------------------------------------------
-- updated_at: se actualiza solo si el UPDATE no lo trae (todas las tablas menos audit_log)
-- ---------------------------------------------------------------------

-- +goose StatementBegin
CREATE TRIGGER trg_person_updated_at AFTER UPDATE ON person
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE person SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_user_account_updated_at AFTER UPDATE ON user_account
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE user_account SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_device_session_updated_at AFTER UPDATE ON device_session
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE device_session SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_login_attempt_updated_at AFTER UPDATE ON login_attempt
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE login_attempt SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_app_setting_updated_at AFTER UPDATE ON app_setting
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE app_setting SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_period_updated_at AFTER UPDATE ON period
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE period SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_cost_share_rule_updated_at AFTER UPDATE ON cost_share_rule
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE cost_share_rule SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_cost_share_rule_part_updated_at AFTER UPDATE ON cost_share_rule_part
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE cost_share_rule_part SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_category_updated_at AFTER UPDATE ON category
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE category SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_budget_updated_at AFTER UPDATE ON budget
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE budget SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_store_updated_at AFTER UPDATE ON store
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE store SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_product_group_updated_at AFTER UPDATE ON product_group
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE product_group SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_store_product_group_updated_at AFTER UPDATE ON store_product_group
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE store_product_group SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_product_updated_at AFTER UPDATE ON product
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE product SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_presentation_updated_at AFTER UPDATE ON presentation
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE presentation SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_credit_updated_at AFTER UPDATE ON credit
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE credit SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_credit_term_rate_updated_at AFTER UPDATE ON credit_term_rate
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE credit_term_rate SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_promotion_updated_at AFTER UPDATE ON promotion
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE promotion SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_planned_item_updated_at AFTER UPDATE ON planned_item
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE planned_item SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_income_updated_at AFTER UPDATE ON income
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE income SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_expense_updated_at AFTER UPDATE ON expense
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE expense SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_expense_item_updated_at AFTER UPDATE ON expense_item
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE expense_item SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_expense_share_updated_at AFTER UPDATE ON expense_share
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE expense_share SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_attachment_updated_at AFTER UPDATE ON attachment
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE attachment SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_price_updated_at AFTER UPDATE ON price
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE price SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_credit_purchase_updated_at AFTER UPDATE ON credit_purchase
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE credit_purchase SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_credit_purchase_share_updated_at AFTER UPDATE ON credit_purchase_share
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE credit_purchase_share SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_scheduled_payment_updated_at AFTER UPDATE ON scheduled_payment
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE scheduled_payment SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_scheduled_payment_share_updated_at AFTER UPDATE ON scheduled_payment_share
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE scheduled_payment_share SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_credit_statement_updated_at AFTER UPDATE ON credit_statement
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE credit_statement SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_credit_payment_updated_at AFTER UPDATE ON credit_payment
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE credit_payment SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_savings_goal_updated_at AFTER UPDATE ON savings_goal
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE savings_goal SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_savings_movement_updated_at AFTER UPDATE ON savings_movement
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE savings_movement SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_meal_slot_updated_at AFTER UPDATE ON meal_slot
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE meal_slot SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_recipe_updated_at AFTER UPDATE ON recipe
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE recipe SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_recipe_ingredient_updated_at AFTER UPDATE ON recipe_ingredient
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE recipe_ingredient SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_meal_plan_entry_updated_at AFTER UPDATE ON meal_plan_entry
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE meal_plan_entry SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_inventory_item_updated_at AFTER UPDATE ON inventory_item
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE inventory_item SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_health_profile_updated_at AFTER UPDATE ON health_profile
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE health_profile SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_weight_entry_updated_at AFTER UPDATE ON weight_entry
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE weight_entry SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_nutrition_fact_updated_at AFTER UPDATE ON nutrition_fact
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE nutrition_fact SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_notification_updated_at AFTER UPDATE ON notification
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE notification SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_notification_recipient_updated_at AFTER UPDATE ON notification_recipient
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE notification_recipient SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_notification_subscription_updated_at AFTER UPDATE ON notification_subscription
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE notification_subscription SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_notification_delivery_updated_at AFTER UPDATE ON notification_delivery
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE notification_delivery SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_job_updated_at AFTER UPDATE ON job
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE job SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_backup_updated_at AFTER UPDATE ON backup
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE backup SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_generated_report_updated_at AFTER UPDATE ON generated_report
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE generated_report SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_idempotency_key_updated_at AFTER UPDATE ON idempotency_key
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE idempotency_key SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_recipe_meal_slot_updated_at AFTER UPDATE ON recipe_meal_slot
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE recipe_meal_slot SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_inventory_movement_updated_at AFTER UPDATE ON inventory_movement
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE inventory_movement SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_shopping_list_updated_at AFTER UPDATE ON shopping_list
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE shopping_list SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_shopping_list_item_updated_at AFTER UPDATE ON shopping_list_item
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE shopping_list_item SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
-- Borra en orden inverso a la creación (hijas antes que madres). Índices y disparadores se van con su tabla.
DROP TABLE IF EXISTS generated_report;
DROP TABLE IF EXISTS backup;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS notification_delivery;
DROP TABLE IF EXISTS notification_subscription;
DROP TABLE IF EXISTS notification_recipient;
DROP TABLE IF EXISTS notification;
DROP TABLE IF EXISTS nutrition_fact;
DROP TABLE IF EXISTS weight_entry;
DROP TABLE IF EXISTS health_profile;
DROP TABLE IF EXISTS shopping_list_item;
DROP TABLE IF EXISTS shopping_list;
DROP TABLE IF EXISTS inventory_movement;
DROP TABLE IF EXISTS inventory_item;
DROP TABLE IF EXISTS meal_plan_entry;
DROP TABLE IF EXISTS recipe_ingredient;
DROP TABLE IF EXISTS recipe_meal_slot;
DROP TABLE IF EXISTS recipe;
DROP TABLE IF EXISTS meal_slot;
DROP TABLE IF EXISTS savings_movement;
DROP TABLE IF EXISTS savings_goal;
DROP TABLE IF EXISTS credit_payment;
DROP TABLE IF EXISTS credit_statement;
DROP TABLE IF EXISTS scheduled_payment_share;
DROP TABLE IF EXISTS scheduled_payment;
DROP TABLE IF EXISTS credit_purchase_share;
DROP TABLE IF EXISTS credit_purchase;
DROP TABLE IF EXISTS price;
DROP TABLE IF EXISTS attachment;
DROP TABLE IF EXISTS expense_share;
DROP TABLE IF EXISTS expense_item;
DROP TABLE IF EXISTS expense;
DROP TABLE IF EXISTS income;
DROP TABLE IF EXISTS planned_item;
DROP TABLE IF EXISTS promotion;
DROP TABLE IF EXISTS credit_term_rate;
DROP TABLE IF EXISTS credit;
DROP TABLE IF EXISTS presentation;
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS store_product_group;
DROP TABLE IF EXISTS product_group;
DROP TABLE IF EXISTS store;
DROP TABLE IF EXISTS budget;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS cost_share_rule_part;
DROP TABLE IF EXISTS cost_share_rule;
DROP TABLE IF EXISTS period;
DROP TABLE IF EXISTS idempotency_key;
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS app_setting;
DROP TABLE IF EXISTS login_attempt;
DROP TABLE IF EXISTS device_session;
DROP TABLE IF EXISTS user_account;
DROP TABLE IF EXISTS person;
