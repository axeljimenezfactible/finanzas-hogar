# S1-esquema-1.3 · reporte final (flujo P1c, 11-oct-2026)

## Hecho
- **Migración nueva `backend/db/migrations/00003_puerta1b.sql`** (esquema 1.2 → 1.3, Up y Down completos):
  - K36: `health_profile.goal_adjust_lose_kcal` y `.goal_adjust_gain_kcal` (INTEGER NULL, `CHECK (≥ 0)` de columna; NULL = usa el valor de la casa) y claves `health.goal_adjust_lose_kcal = 0` y `health.goal_adjust_gain_kcal = 0` en `app_setting`.
  - K44: inventario de 12 parámetros operativos en anotaciones `-- @param` (5 ya existían en 00001 y 00002 y no se tocan; 7 se siembran en 00003): `price.zero_weight_after_days = 90` (K22, K40), `auth.password_min_length = 12` y `price.profeco_import_max_gb = 1` (K6), `attachment.max_size_mb = 10` y `comparator.timeout_seconds = 5` (K43) y los dos de kcal (K36). 7 renglones `-- @param_sin_valor` con lo que existe sin valor o vive en otra tabla (optimizer.* de K16, `optimizer.date_window_days`, `send_time` y `advance_days` de K43 que propone S7, cinco valores de K31 para la puerta 4, la ventana de la cola, que no existe por K9, y los umbrales de avisos, ver bloqueo).
  - Disparadores `trg_app_setting_param_int_ins` y `_upd`: los 12 parámetros enteros solo aceptan un entero JSON ≥ 0 (tipo y signo; sin mínimos razonables, K44).
  - Down: borra disparadores, las 7 claves (en cualquier alcance) y las 2 columnas; deja sqlite_master, `app_setting` y los datos de 1.2 idénticos.
- **`docs/contratos/reglas-negocio.md` 1.3** (parches): RN-12 FIRME con K37 (compiten `ticket_manual`, `ticket_ocr` y `manual`; `site_assisted` y PROFECO solo referencia) y K40 (peso = (90 − edad) / 61: día 30 = 0.9836, día 60 = 0.4918, día 90 = 0; forma general con los dos parámetros); RN-17 FIRME con K36 (meta = round_half_up(MSJ × actividad − ajuste_bajar) o + ajuste_subir; `maintain` no ajusta; 0 y 0 al inicio; por persona con el valor de la casa; ejemplo 2,556 y 2,056); RN-27 FIRME con K39 (219.88 y 629.35 como cuentas citadas de `model.json`, más 439.76); regla nueva RN-29 (tabla clave · valor inicial · unidad · decisión · sección de Ajustes · migración); índice y «Cambios en 1.3». PROVISIONAL: 14 apariciones (11 líneas) antes → 6 (5 líneas) después, ninguna cita un bloqueo resuelto.
- **`backend/db/tools/gen_dict.py`**: lee la versión de la cabecera de cada migración, `@used_new` por versión, `@param` y `@param_sin_valor`; genera encabezado 1.3, «Cambios en 1.3», la semilla de la 00003, la sección «Parámetros operativos de app_setting (K44)» (comprueba cada valor inicial contra lo sembrado y sale con 1 si no coincide) y la fila 13 de Ajustes con las columnas nuevas. Todo idempotente.
- **`docs/contratos/diccionario-datos.md` 1.3** regenerado (55 tablas, 152 índices, 73 disparadores); segunda corrida `cmp` idéntica.
- **`backend/db/tools/s1_verifica.py`**: separa E11/E12/esquema final; pruebas nuevas K44, DDL y DOWN13; S1-03 acepta cuentas encadenadas («ruta × 7 / 31 = valor», un solo HALF_UP al final, K39); GOOSE y DOWN comparan contra el esquema final.
- **Bloqueo nuevo** `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Archivos
- Creados: `backend/db/migrations/00003_puerta1b.sql`, `docs/estado/bloqueos/S1-parametros-avisos.md`, `docs/estado/p1c/S1-esquema-1.3.md`.
- Modificados: `docs/contratos/reglas-negocio.md`, `docs/contratos/diccionario-datos.md` (solo por regeneración), `backend/db/tools/gen_dict.py`, `backend/db/tools/s1_verifica.py`.
- Sin cambios: `00001_esquema.sql` y `00002_puerta1.sql` (ni un byte; no aparecen en `git diff`).

## Verificación
- `sqlite3` CLI con las 3 migraciones en `datos/tmp`: `integrity_check` = ok · `foreign_key_check` = 0 filas · tablas/STRICT = 55/55.
- goose v3.28.0: `up-to 2` → volcado de sqlite_master y `app_setting`; `up` → «successfully migrated database to version: 3»; 7 claves con su valor (10, 12, 5, 0, 0, 1, 90); `UPDATE app_setting … '-1'` → «app_setting: este parámetro es un entero mayor o igual a 0 (K44)»; `down-to 2` → `cmp` de sqlite_master (1256 líneas) y de `app_setting` idénticos a 1.2.
- `GOOSE=…/gobin/goose python3 -I backend/db/tools/s1_verifica.py` → exit 0, **RESUMEN TODO OK** (26 pruebas OK). Destacadas: S1-02 (55 STRICT, 152 índices, 73 disparadores) · S1-03 (29 reglas con ejemplo; RN-27 citas 6/6, cuentas citadas 6) · S1-16 (8 líneas PROVISIONAL, solo citan `S1-parametros-avisos`) · KUESKI 389500 / 239528, 0 diferencias · K44 (2 columnas, 7 claves nuevas + 5 previas con su valor, 12 anotaciones = lo sembrado, 9 negativas rechazadas incluido ajuste −1, 5 positivas aceptadas, cifras 2556/2056, 0.9836/0.4918/0, 219.88/629.35 presentes en las reglas, 0 marcas que citen bloqueos resueltos, los 5 bloqueos de la puerta con respuesta, 0 claves alerts.*) · DDL (00001 y 00002 sin comentarios = git HEAD) · DOWN13 (sqlite_master = 1.2 con 240 objetos, `app_setting` = 1.2 con 10 claves, 130 filas de 54 tablas iguales) · DOWN (1.1 idéntico) · GOOSE (códigos [0, 0, 0], tablas [55, 0, 55], versión 3).
- `gen_dict.py` dos veces → `cmp` idéntico; «valores iniciales que no coinciden con lo sembrado: []»; «columnas sin descripción: []».
- Cifras con Decimal ROUND_HALF_UP: 973.75 × 7 / 31 = 219.88 (diario 31.41); 2,517.38 × 7 / 28 = 629.35 (diario 89.91); 1,947.50 × 7 / 31 = 439.76; (90 − 30)/61 = 0.9836, (90 − 60)/61 = 0.4918; MSJ 1,648.75 × 1.55 = 2,555.5625 → 2,556 y − 500 → 2,056.
- `grep -c PROVISIONAL`: 00001 = 0, 00002 = 0, 00003 = 1, diccionario = 2 (la leyenda general y la fila de avisos), reglas = 5 líneas (leyenda general y RN-29). Bloqueos citados: solo `S1-parametros-avisos` (abierto).
- `datos/tmp` borrado.

## Supuestos
1. Se conservan los nombres de las claves que ya existían (`auth.*` en lugar de `login.*`, `optimizer.*`); las nuevas usan `health.*`, `price.*`, `auth.*`, `attachment.*` y `comparator.*`. Renombrar habría cambiado datos de 1.2 y C1.
2. El valor inicial para «restablecer» no se guarda en una columna nueva: la forma de `app_setting` no cambia y el valor lo conocen C1 (`default` por propiedad) y el diccionario (ver propuesta).
3. Con los dos días configurables, RN-12 usa peso = (D0 − edad) / (D0 − Dv + 1); con 30 y 90 es exactamente K40.
4. Las constantes de fórmula (IMC OMS, factores de actividad, Mifflin-St Jeor) y los límites de formato del contrato (página 50/200, cursor 512, búsqueda 200, usuario 3–64, clave 16–128) no son parámetros de K44; se listan en RN-29 y se proponen así a Axel.
5. El disparador de tipo (entero JSON ≥ 0) es integridad de tipo, no un mínimo razonable.
6. El «comparador máximo 5 s» de K43 es el «tiempo máximo del optimizador» de K31.
7. Sección de Ajustes para contraseña, archivo de PROFECO y retención de idempotencia: propuesta de S1 con el criterio de K44.

## Bloqueos
- `docs/estado/bloqueos/S1-parametros-avisos.md` (nuevo): umbrales 80/100 % y días antes del pago próximo ya viven en `notification.params_json` (00001); ¿se quedan ahí o pasan a `alerts.*` de `app_setting`? Seguí con la opción A (sin cambios); RN-29 y el diccionario los marcan PROVISIONAL.
- Resueltos y aplicados: `S1-ajuste-objetivo-kcal` (K36), `S1-fuentes-precio` (K37), `criterios-promedios-redondeo` (K39), `criterios-frescura-extremos` (K40), `C1-valores-k31` (K43).

## Pendientes
- `docs/estado/p1/INTEGRACION.md` §7.2 (`interest_includes_vat` o tercer estado del IVA, K18) y §7.7 (renombrar `credit.rate_base` o `credit_purchase.rate_basis`, R2-07): piden migración de S1 y no estaban en el alcance de P1c.
- Si la puerta 2 elige la opción B del bloqueo `S1-parametros-avisos`: migración 00004 (claves `alerts.*` y quitar los valores de `params_json`, con Down).
- Mínimos y máximos de los parámetros: se aprueban en la puerta 2 (propuestas abajo).

## Solicitudes fuera de alcance
1. C1 v1.2 (`backend/internal/api/c1gen/`, regenerar `openapi.yaml` y extractos): exponer con `x-app-setting`, tipo, unidad y `default` (para «restablecer») las 7 claves nuevas: `PriceSettings.zero_weight_after_days` y `profeco_import_max_gb`; ajustes de comparador `timeout_seconds`; `attachment.max_size_mb` en `uploadExpenseAttachment`; `auth.password_min_length` en `changeOwnPassword` (el `minLength: 12` estático pasa a ser el valor inicial); `HealthProfile.goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` (enteros ≥ 0, admiten null = valor de la casa) y ajustes de salud de la casa. Motivo: K36, K43 y K44.
2. U1 1.2 y U4 (`docs/design/U1-arquitectura.md`, U4 en WF-1): cada parámetro en su sección de Ajustes con valor inicial y «restablecer»; en Salud, el ajuste por persona con el de la casa como inicial. Motivo: K36 y K44.
3. Criterios (`docs/criterios/S4.md` S4-07, `S5.md` S5-04, `S2.md` S2-08, `S1.md`): citar RN-17, RN-12 y RN-27 1.3 y agregar un criterio para RN-29 (12 claves con su valor, tipo entero ≥ 0, Down de la 00003). Motivo: K36, K39, K40 y K44.
4. Orquestador: llevar el bloqueo `S1-parametros-avisos` a la puerta 2.

## Propuestas para Axel (K44, puerta 2: mínimos y máximos razonables; hoy sin tope en el esquema)
- `auth.max_failed_attempts`: 3 a 10. Menos de 3 bloquea a la casa por un error de dedo en el teléfono; más de 10 debilita el freno a quien adivina contraseñas.
- `auth.lockout_window_minutes`: 5 a 60. Con menos de 5 el límite casi no frena (5 intentos cada 5 min); con más de 60 un error deja fuera demasiado tiempo.
- `auth.password_min_length`: 12 a 64. No bajar del valor aprobado en K6; NIST SP 800-63B rev. 4 recomienda 15 cuando la contraseña es el único factor y admitir al menos 64.
- `session.max_age_days`: 1 a 365 (el tope de 365 es el ejemplo de K44); la sesión se renueva con el uso.
- `price.stale_after_days` y `price.zero_weight_after_days`: 1 a 365, con la regla stale ≤ zero_weight (si no, la recta de RN-12 no existe). La API la valida al guardar.
- `attachment.max_size_mb`: 1 a 25. Una foto de ticket de Android pesa unos cuantos MB; más de 25 solo llena el disco. Propuesta: 1 MB = 1,048,576 bytes.
- `price.profeco_import_max_gb`: 1 a 4, por el disco del servidor de casa.
- `comparator.timeout_seconds`: 1 a 30. Más de 30 s deja al usuario esperando una respuesta HTTP sin aviso.
- `idempotency.retention_hours`: 1 a 168 (7 días). La cola del teléfono no depende de esto (sus claves en `expense` y `weight_entry` no se purgan).
- `health.goal_adjust_lose_kcal` y `goal_adjust_gain_kcal`: 0 a 1,000 kcal. Un ajuste mayor de 1,000 kcal al día no se recomienda sin supervisión médica; la meta resultante ya tiene su rango de captura (500 a 10,000 kcal).
- «Restablecer»: que C1 lleve el `default` de cada parámetro (opción sin cambio de esquema). Otra opción: columna `app_setting.default_value_json` en una 00004.
- Confirmar que las constantes de fórmula y los límites de formato del contrato no son parámetros configurables (RN-29).

## Consumo aproximado
- Acciones de herramienta: 38. Leído: ~205 KB. Escrito: ~60 KB a mano (migración 10.6 KB, bloqueo 2.6 KB, reglas ~21 KB de parches, gen_dict y s1_verifica ~17 KB de parches, reportes ~9 KB) más el diccionario regenerado por script (+~10 KB netos).
