# Reporte S1-esquema-1.2 · flujo P1 «Aplicar respuestas de la puerta 1»
Fecha: 11-oct-2026 · Esquema 1.1 → 1.2 · Fuente: `docs/plan/decisiones.md` §K y la «Respuesta de Axel» de los bloqueos asignados.

## Hecho
1. **Migración goose `backend/db/migrations/00002_puerta1.sql`** (Up y Down completos, 14.7 KB):
   - K3: `attachment.ocr_draft_json` (objeto JSON) y `attachment.ocr_draft_saved_at` (instante con «Z»; van juntos por CHECK).
   - K4: tabla nueva `telegram_chat_candidate` (`chat_id` único, `telegram_name`, `detected_at`, `confirmed_at`, `person_id`; CHECK confirmado ⇔ con persona; índice de pendientes, índice de la FK y disparador `updated_at`). Solo guarda chat_id y nombre de Telegram, más `id`, `created_at` y `updated_at` por convención.
   - K9: `idempotency_key` (16 a 128, mismo formato que el encabezado) con índices únicos `ux_expense_idempotency_key` y `ux_weight_entry_idempotency_key`, y `captured_at` (NULL = captura en línea) en `expense` y `weight_entry`.
   - K12: `rate_decimals` en `credit_term_rate`, `credit_purchase` (solo con tasa) y `promotion`; `credit.annual_rate_decimals` (solo con tasa anual).
   - K18: `credit.rate_base` (`daily` | `monthly`), `credit.min_payment_pct_bp` (0 a 10000) y `credit.min_payment_fixed_cents` (≥ 0), solo en `bank_card` y `store_credit`. Ya existían `annual_rate_num`/`annual_rate_den` y `vat_on_interest_bp` (si incluye IVA).
   - K23: `health_profile.goal` (`lose` | `maintain` | `gain`).
   - Ajustes sembrados: `projection.repeat_cadence_enabled` = false (K21), `session.max_age_days` = 90 (K11), `idempotency.retention_hours` = 24 (K6). K16: las claves `optimizer.*` de restricciones duras no se siembran.
   - Down en orden inverso (DELETE de los ajustes, DROP COLUMN, DROP INDEX, DROP TABLE): deja `sqlite_master` del esquema 1.1 byte por byte.
2. **`00001_esquema.sql`, solo comentarios** (13 líneas): versión del esquema; redacción neutral de `person.fortnight_payment_cap_cents` (K13, R2-03); retención de 24 h (K6); rangos de `period` FIRMES (K8); vigencia de `budget.valid_to_date` (K2); lista de compras compartida (K1); escala de actividad con sus 5 factores (K23). Quedan 0 PROVISIONAL.
3. **`backend/db/tools/gen_dict.py`** (desde el Anexo A, cierra R2-08): lee todas las migraciones en orden, migra una base en memoria, entiende `-- @alter`, `-- @used_new` y `-- @after`, y toma el texto fijo siempre del diccionario del repositorio. Uso: `python3 -I backend/db/tools/gen_dict.py [salida.md]`.
4. **`backend/db/tools/s1_verifica.py`** (desde el Anexo B, cierra R2-08): aplica 00001 y 00002; nuevas líneas KUESKI (389500 y 239528 sobre 1.2), K12 (17.27 % y 25.8685 % con `rate_decimals`, 1,034.74 con la fracción), NEG12 (24 negativas y 10 positivas de 1.2, ajustes sembrados y K16 sin valor), DOWN (Down de 00002 = esquema 1.1 y ajustes de 1.1; Down de 00001 deja 0 objetos; Up otra vez = 1.2) y GOOSE (up / down-to 0 / up, versión 2, mismas columnas y objetos que con sqlite3). S1-16 exige ahora que cada PROVISIONAL cite un bloqueo **abierto**.
5. **`docs/contratos/diccionario-datos.md` 1.2**: regenerado con gen_dict (55 tablas, 152 índices, 69 disparadores); encabezado y «Cambios en 1.2»; 6 filas de glosario; datos de referencia de la 00002; «Claves previstas sin valor» con K16 (salen `price.stale_weight`, `session.max_age_days` e `idempotency.retention_hours`); filas 1, 2, 5, 9, 12 y 13 de «Ajustes: 14 secciones»; invariantes 14 a 16; notas para S3, S4/C1, S6, S7 y U1/S2.
6. **`docs/contratos/reglas-negocio.md` 1.2** (cierra R2-04): encabezado, «Cambios en 1.2», índice de 28 reglas; C-04 («Z» solo en entradas, K12), C-13 (55 tablas), C-15 (24 h, K6); RN-02 y RN-07 (K14), RN-04 (K15 y K13), RN-06 (K8, semana lunes a domingo), RN-11 (K16 y K17), RN-12 (K22), RN-15 (K1), RN-17 (K23), RN-21 (K19: 6,293.44 = 1,258.70 + 5,034.74), RN-22 (K2); reglas nuevas RN-23 (K3), RN-24 (K9), RN-25 (K12), RN-26 (K18), RN-27 (K20: 62.82 diario en oct-2026; el redondeo del semanal queda PROVISIONAL) y RN-28 (K21). Cada regla tiene ejemplo; las cifras del oráculo se citan con su ruta y las demás dicen «ilustrativo».
7. **Marcas PROVISIONAL restantes**: 00001 = 0, 00002 = 0, diccionario = 1 (la definición de la marca en «Cómo leer»), reglas = 11 líneas (la definición en «Estado de cada regla» y 10 marcas que citan bloqueos abiertos: `S1-fuentes-precio` y `S1-ajuste-objetivo-kcal`, míos, y `criterios-frescura-extremos` y `criterios-promedios-redondeo`, del agente de criterios).

## Archivos
- Nuevos: `backend/db/migrations/00002_puerta1.sql`, `backend/db/tools/gen_dict.py`, `backend/db/tools/s1_verifica.py`, `docs/estado/bloqueos/S1-ajuste-objetivo-kcal.md`, `docs/estado/bloqueos/S1-fuentes-precio.md`, `docs/estado/p1/S1-esquema-1.2.md`.
- Modificados: `backend/db/migrations/00001_esquema.sql` (solo comentarios), `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md`.

## Verificación
- (a) `sqlite3 datos/tmp/s1.db` con el Up de 00001 y 00002: `integrity_check` = ok · `foreign_key_check` = 0 filas · 55 STRICT de 55 tablas · 152 índices · 69 disparadores.
- (b) volcado ordenado de `sqlite_master` (type, name, tbl_name, sql) antes de 00002 contra después del Down de 00002: `diff` vacío (1263 líneas, 111,282 bytes).
- (c) `goose -dir backend/db/migrations sqlite3 datos/tmp/g.db up` → rc 0, «migrated database to version: 2»; `down-to 0` → rc 0, versión 0; `up` → rc 0, versión 2; `status` muestra 00001 y 00002 aplicadas.
- (d) `git show HEAD:…/00001_esquema.sql` contra el archivo actual, sin comentarios (`sed 's/--.*//' | grep -v '^\s*$'`): `diff` vacío (1302 líneas). `git diff --numstat`: 13 líneas cambiadas, todas de comentario.
- (e) `GOOSE=<goose v3.28.0> python3 -I backend/db/tools/s1_verifica.py` → rc 0, 23 líneas OK, 0 FALLA, «RESUMEN TODO OK»; KUESKI: suma de pagos de la despensa 1 = 389500 y parte de Axel = 239528 sobre el esquema 1.2; S1-12 compara 164 valores del oráculo con 0 diferencias; S1-03: 28 reglas con ejemplo y 0 citas que no coinciden.
- (f) `python3 -I backend/db/tools/gen_dict.py datos/tmp/s1-dic.md` y `cmp` contra `docs/contratos/diccionario-datos.md`: idénticos (173,798 bytes). La diferencia del diccionario 1.2 contra el 1.1 solo trae lo esperado (columnas y tabla nuevas, marcas quitadas, conteos).
- `datos/tmp`: borré mis temporales (`s1*`, `g.db`, `g.log`). Quedan `datos/tmp/rev-u1` y `datos/tmp/verify-wf0`, que son de otros agentes.

## Supuestos
- K12 se aplicó también a `credit.annual_rate_decimals` y `promotion.rate_decimals`, porque K12 es general («tasas en pantalla como se capturaron») y K18 hace que Axel capture la tasa anual. Si el orquestador solo quiere las dos pedidas, basta quitar esas dos columnas de la 00002 (y del Down) antes del commit.
- Las columnas nuevas admiten NULL porque `ALTER TABLE ADD COLUMN` no acepta NOT NULL sin valor fijo ni UNIQUE; cada comentario dice qué significa NULL. `rate_decimals` solo exige ≥ 0 (no inventé un tope).
- `session.max_age_days` = 90 e `idempotency.retention_hours` = 24 se sembraron porque K11 y K6 fijan esos valores y las claves ya estaban previstas en el diccionario 1.1.
- RN-12 y RN-27 siguen, como PROVISIONAL, las lecturas que usa el agente de criterios en sus bloqueos abiertos `criterios-frescura-extremos` (lectura (a), peso = (90 − edad)/60) y `criterios-promedios-redondeo` (opción (a), un solo redondeo), para que reglas y criterios digan lo mismo. Para el peso S1 recomienda la lectura (b), (90 − edad)/61: es el texto de la opción C que eligió Axel («peso 1 hasta los 29 días, que baja a 0 a los 90 días») y hace que todo precio viejo pese menos.
- RN-17: los factores 1.375, 1.55 y 1.725 salen de la opción B del bloqueo `S1-meta-kcal`, que la D elegida incluye; «mantener» no ajusta.
- RN-22: «solo este mes» deja el monto anterior desde el mes siguiente (es lo que significa «solo este mes» con renglones de vigencia que no se traslapan).
- El interruptor K21 quedó en la fila «Presupuesto» del mapa de Ajustes, con la nota de que U4 decide la pantalla.

## Bloqueos
- `docs/estado/bloqueos/S1-ajuste-objetivo-kcal.md`: cuánto ajustan «bajar» y «subir» la meta sugerida (RN-17 marcado PROVISIONAL en ese punto).
- `docs/estado/bloqueos/S1-fuentes-precio.md`: si compiten `site_assisted` y `manual` en la comparativa (RN-12 marcado PROVISIONAL en ese punto).
- Además RN-12 y RN-27 citan dos bloqueos abiertos que escribió en paralelo el agente de criterios: `docs/estado/bloqueos/criterios-frescura-extremos.md` y `docs/estado/bloqueos/criterios-promedios-redondeo.md`.

## Pendientes
- Ninguno dentro de mi alcance. Cuando Axel responda los dos bloqueos nuevos, se quitan sus marcas en RN-12 y RN-17 (sin cambio de esquema salvo la opción D de `S1-ajuste-objetivo-kcal`).

## Solicitudes fuera de alcance
1. `docs/contratos/C1-datos.md` y `docs/contratos/openapi.yaml` (generador `backend/internal/api/c1gen/`): exponer `rate_decimals` en CreditTermRate, CreditPurchase y Promotion; `annual_rate_decimals`, `rate_base`, `min_payment_pct_bp` y `min_payment_fixed_cents` en Credit; `goal` en HealthProfile; `captured_at` y la clave de la cola en Expense y WeightEntry; el borrador `ocr_draft_json` del ticket; la lectura y confirmación de `telegram_chat_candidate`. Motivo: columnas nuevas del esquema 1.2.
2. `docs/estado/bloqueos/INDICE.md`: agregar `S1-ajuste-objetivo-kcal` y `S1-fuentes-precio`. Motivo: bloqueos nuevos abiertos.
3. `docs/estado/wf0/S1-modelo-datos.md`, Anexos A y B: cambiar el código por un enlace a `backend/db/tools/`. Motivo: que haya una sola copia (R2-08).
4. `docs/criterios/S4.md` (S4-07): citar RN-17 1.2 y el bloqueo `S1-ajuste-objetivo-kcal`. `docs/criterios/S5.md` (S5-04): RN-12 ya usa su misma lectura provisional (a); falta sumar el bloqueo `S1-fuentes-precio` (qué fuentes compiten). Motivo: criterios y reglas 1.2 alineados.
5. U1/U4: ubicar en Ajustes el interruptor `projection.repeat_cadence_enabled` (K21). Motivo: el mapa de Ajustes lo deja en «Presupuesto» a reserva de U4.

## Paquetes instalados
- Ninguno (goose v3.28.0 del scratchpad del orquestador; sqlite3 3.53.4 y Python del sistema).

## Consumo aproximado
- Acciones de herramienta: 33.
- Leído: unos 230 KB (protocolo 6, §K 8, bloqueos 15, Anexos A y B 45, reglas 37, partes de la migración ~45, diccionario ~30, salidas de verificación ~30).
- Escrito a mano: unos 100 KB (00002 15, comentarios de 00001 4, parches de los scripts 19, parches del diccionario 12, reglas 36, bloqueos 4, este reporte 10). El diccionario (174 KB) lo escribe `gen_dict.py`.
