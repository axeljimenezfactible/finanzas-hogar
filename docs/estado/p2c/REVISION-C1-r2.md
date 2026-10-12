# Revisión de C1 v1.3 (puerta 2) · ronda 2 de 3

Revisor: `revisor-C1-r2`, independiente. No escribí nada de lo que reviso y no corregí nada. Flujo P2c, 11-oct-2026, noche.

Evalúo contra:
- `docs/criterios/C1.md` (versión 1.2, la vigente), sobre todo C1-04, C1-09, C1-14, C1-15, C1-21, C1-22 y C1-23;
- `PUERTA-2.md` §2 y §3;
- las claves de `S1-puerta2` (`00004_puerta2.sql`);
- las «Solicitudes a C1» de `contratos-puerta2.md`.

Todas las cifras y salidas de esta revisión salen de mis propios comandos y scripts en `datos/tmp/rev-c1-r2/`, que ya borré. Ningún resultado lo tomé de un reporte.

## Veredicto
**NO APROBADO.** Hay 0 altos, 1 medio y 5 bajos.

- Los 3 medios de la ronda 1 (H-1, H-2 y H-3) están **cerrados**. Lo comprobé con evidencia propia (V-13 a V-15).
- El medio nuevo, **M-1**, es de la misma clase que H-3. El CHECK de la 00004 también rechaza un porcentaje con `vat_included` NULL o ausente, y C1 solo documenta qué responde la API cuando `vat_included` = true.

## Verificaciones (con mis comandos)

| Id | Qué | Comando o método | Resultado |
|---|---|---|---|
| V-1 | (a) `make verify-c1` con el TMP compartido | `make verify-c1` | **Falló por el entorno, no por C1**: «grep: …/datos/tmp/verify-wf0/verify_c1.log: No such file or directory», `Error 2` en Makefile:131. Otra corrida en paralelo borró con `$(RM_TMP)` el mismo `datos/tmp/verify-wf0` (ver B-4) |
| V-2 | (a) `make verify-c1` con un TMP propio | `make verify-c1 TMP=…/datos/tmp/rev-c1-r2-mk` | «RESULTADO: los 21 archivos coinciden byte a byte con el generador», «comprobaciones OK: 32» y «RESULTADO: TODO OK» |
| V-3 | (a) OpenAPI válido | `make verify-openapi` | «docs/contratos/openapi.yaml: OK» (openapi-spec-validator 0.9.0) |
| V-4 | (a) `build.py --check` | `python3 -I backend/internal/api/c1gen/build.py --check` | 21 archivos «IDÉNTICO». `openapi.yaml`: 702,557 bytes, sha256 6cc6dd6ef7dc…. `routes.go`: 26,202 bytes, «212 rutas». «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| V-5 | Verificación global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=…/datos/tmp/rev-c1-r2/va` | rc=0, «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |
| V-6 | (d) y (e) versión y conteos | script propio con `yaml.safe_load` sobre `openapi.yaml`, y `grep` | `info.version` 1.3.0 y `C1_VERSION = '1.3.0'` (`gen4_build.py` l. 6). `routes.go` dice «contrato C1 v1.3.0» y tiene 212 entradas `{Method: …}`. **168 rutas, 212 operaciones y 255 esquemas**, igual en `extractos/INDICE.md` («212 operaciones en 168 rutas y 255 esquemas»), en `C1-datos.md` l. 9 y l. 627 y en la tabla de `C1-puerta2.md` |
| V-7 | (b) Claves nuevas de S1 en las operaciones de settings | script propio: busca cada clave como propiedad en `requestBody` y en la respuesta 200, con `$ref` resueltos | Las 10 claves de la 00004 (l. 158–180) son propiedad del 200 del `get…` y del cuerpo y el 200 del `update…`:<br>• `reports.*` (5) → `getReportSettings` y `updateReportSettings`<br>• `alerts.send_time`, `alerts.stale_prices_advance_days` y `alerts.meal_plan_missing_weekday` → `getAlertSettings` y `updateAlertSettings`<br>• `optimizer.date_window_days` → `getOptimizerSettings` y `updateOptimizerSettings`<br>• `ocr.timeout_seconds` → `getDataSettings` y `updateDataSettings`<br>Lo mismo con las 3 claves de `@param_rango`. Todas están en `reset_parameters` |
| V-7b | (b) Rangos del `PATCH` frente a los disparadores de la 00004 (l. 113–150) | el mismo script | horas: `pattern ^([01][0-9]\|2[0-3]):[0-5][0-9]$`; días de la semana: 1–7; `monthly_day_of_month`: 1–31; `comparator.timeout_seconds`: 1–30; `price.stale_after_days` ≥ 1; anticipación, ventana, OCR y peso cero ≥ 0. Coinciden con los disparadores y con §11.7 (l. 778–798) |
| V-8 | (b) `NotificationUpdate` sin `send_time` ni `advance_days` | script propio y `grep` | `NotificationUpdate` tiene solo `is_enabled` y `params_json`. `Notification` no tiene esas propiedades. Las dos tienen 0 `x-provisional`. `grep send_time\|advance_days` en `openapi.yaml`, sin contar `reports.*`, `alerts.send_time` ni `alerts.stale_prices_advance_days`: 0 líneas |
| V-9 | (c) `vat_included` | script propio | Tipos: boolean en `Credit` y `CreditCreate`; boolean o null en `CreditUpdate`. Tiene descripción y `x-provisional` (`S3-iva-lectura`). `CreditCreate` y `CreditUpdate` traen la regla 422. `createCredit` y `updateCredit` traen la versión corta y su respuesta 422 |
| V-10 | (c) Alcance del reparto | script propio | `CostShareRuleUpdate.scope` tiene enum `[onward, all]` y `x-provisional` que cita `C1-reparto-alcance` (abierto). Su texto dice «Obligatorio si viene parts y solo se acepta con parts (si no, 422 validation_failed)», sin valor por omisión. El `x-tablas` de `updateCostShareRule` incluye `expense`, `expense_item` y `expense_share` |
| V-11 | (c) `assigned_store_id` y grupos | script propio (descripciones crudas) | `ShoppingListItem.assigned_store_id` está documentado (Q2-39 a Q2-42 A).<br>`groups[]` trae `store_id`, `visit_weekday`, `visit_date` («Fecha de la próxima visita del mercado … Q2-42 A …», opcional) y `total_cents`. `ShoppingList.total_cents` = «Σ groups[].total_cents, solo artículos con precio».<br>El orden de `groups.items` es el de `S4-interfaz.md` l. 216 |
| V-12 | (c) Acumulados de `AccrualVsCash` | script propio y `grep -A3` (`openapi.yaml` l. 11237–11245) | `months[].required` incluye `cumulative_accrued_cents`, `cumulative_cash_cents` y `cumulative_difference_cents`, cada uno con su descripción y el mismo origen (Q2-65 A) |
| V-13 | H-1 | script propio sobre `openapi.yaml` y los 17 extractos, y `grep -rE` en el repositorio | 0 de 610 propiedades terminan en «_», tanto en `openapi.yaml` como en los extractos. Fuera de `docs/estado/`: 0 apariciones de `visit_date_` o `cheapest_store_id_`. `visit_date` está en `openapi.yaml` l. 13517. La guarda está en `gen2_schemas.py` l. 392–394 |
| V-14 | H-2 | M-11 en `C1-datos.md` l. 372, y recálculo propio desde `oraculo/model.json` con `Decimal` | M-11 nombra los 3 campos, su fórmula, el origen («desde el primer evento… El origen no es el primer mes visible») y un ejemplo.<br>Mi recálculo da P1 1947.50 / 973.75 / 973.75; P2 3895.00 / 2921.25 / 973.75; P6 14362.48 / 13103.78 / 1258.70. Coincide con las cifras de M-11.<br>En los 6 periodos, Σ devengado − Σ efectivo = `diferencia_acumulada` (973.75; 973.75; 1258.69; 1258.70; 1258.69; 1258.70) |
| V-15 | H-3 y M-1: el CHECK real | base propia (Up de 00001 a 00004 con `awk` y `sqlite3`), `INSERT` de un `store_credit` con (`vat_included`, `vat_on_interest_bp`) | (NULL, 1600) → «CHECK constraint failed». (1, 1600) → «CHECK constraint failed». (0, 1600) → aceptado. `UPDATE … SET vat_included = NULL` sobre (0, 1600) → «CHECK constraint failed» |
| V-16 | Extractos | script propio que normaliza `componentes-comunes.yaml#/` → `#/` | 17 YAML con `info.version` 1.3.0. Todos ≤ 80,000 bytes (el mayor, `ajustes-c.yaml`, con 78,546). 212 de 212 operaciones. Las 212 operaciones y las 304 copias de esquemas son **idénticas** a `openapi.yaml`: 0 diferencias |
| V-17 | C1-15: marcas PROVISIONAL | script propio: cada `x-provisional: true` y los bloqueos que cita su texto | Hay 6 marcas, cada una con un bloqueo abierto (sin «Respuesta de Axel»):<br>• `scope` → `C1-reparto-alcance`<br>• `vat_included` (×3) y `missing_terms` → `S3-iva-lectura`<br>• `JobUpdate.schedule` → `C1-hora-reportes-job`<br>Los 2 bloqueos nuevos traen qué falta, opciones, propuesta e impacto. `S1-parametros-avisos` está resuelto y tiene 0 citas en `openapi.yaml` |
| V-18 | «Solicitudes a C1» de `contratos-puerta2.md` (l. 65–72) | lectura de las descripciones crudas | Se cumplen las 7: `assigned_store_id`, `groups[].store_id`, `recommended_store_id`, `estimated_price_cents` («Ausente = sin precio: no suma»), `visit_weekday`/`visit_date`, los totales con «se reconocen por estimated_price_cents ausente», `cheapest_store_id` y `difference_vs_cheapest_cents` ausentes en tiendas incompletas, e `identical_count` = «Renglones idénticos (no piezas)» |
| V-19 | PUERTA-2 §2 y §3, en lo que toca a C1 | lectura de §2, de §3.1 y de §3.13 a §3.25 contra el contrato | Están aplicadas Q2-11 B, Q2-39 a Q2-42 A, Q2-08 A (órdenes y campos ausentes), Q2-65 A, Q2-66 B, Q2-67 opción 1, Q2-68 C, Q2-17 C, Q2-33 (b) y Q2-86 A.<br>Ninguna pregunta «sin recomendación» quedó decidida en C1: Q2-50 a Q2-52 no cambian `missing_count` ni `stale_count`, y Q2-94 no se agregó |

## Hallazgos

### M-1 (media) · C1-04 y Q2-67 opción 1 · el CHECK también rechaza un porcentaje con `vat_included` ausente o null, y C1 no dice qué responde la API
- **HECHO.** La columna de la 00004 dice: `CHECK (vat_included IN (0, 1) AND (vat_on_interest_bp IS NULL OR vat_included IS 0))`. El comentario de l. 20–21 lo dice en palabras: «el CHECK rechaza un porcentaje con vat_included 1 y también con vat_included NULL». V-15 lo comprueba.
- **HECHO.** `CREDIT_VAT_RULE` y `CREDIT_VAT_RULE_OP` (`gen2_schemas.py` l. 322–332) cubren solo `vat_included = true`. Esos textos salen en:
  - `openapi.yaml`: l. 15487 (`CreditCreate`), l. 15589 (`CreditUpdate`), l. 6164 (`createCredit`) y l. 6326 (`updateCredit`);
  - `extractos/ajustes-c.yaml`: l. 1344 y l. 1446.
- **HECHO.** En `CreditCreate`, `vat_included` no es obligatorio. `C1-datos.md` l. 36 dice «ausente = falta capturar», y `CreditUpdate` dice «null borra el valor».
- **Escenarios que fallan:**
  1. `POST /api/v1/credits` con `{"name": "Tarjeta X", "kind": "store_credit", "holder_person_id": 1, "payment_frequency": "monthly", "vat_on_interest_bp": 1600}` y sin `vat_included`. El cuerpo es válido contra `CreditCreate`, y es justo como se capturaba el IVA en C1 v1.2. La base lo rechaza.
  2. `PATCH /api/v1/credits/{id}` con `{"vat_included": null}` sobre un crédito con (false, 1600). La base lo rechaza.
- **Qué falta.** C1 no dice si la respuesta es 422 ni con qué `fields[].field`. Tampoco dice si la API deduce `vat_included = false` a partir del porcentaje. Quien implemente S6 tendría que adivinarlo. Si S6 mapea solo el caso documentado, la violación del CHECK sale como 5xx. Es el mismo defecto que H-3 de la ronda 1, en la otra combinación que rechaza el CHECK.
- **Corrección sugerida.** En `gen2_schemas.py` (`CREDIT_VAT_RULE` y `CREDIT_VAT_RULE_OP`):
  - Escribir la regla completa del CHECK: `vat_on_interest_bp` con valor exige que la combinación que quedaría guardada tenga `vat_included = false`. Con true, ausente o null, la API responde 422 `validation_failed`, con `fields[].field = vat_on_interest_bp`, y no guarda nada.
  - Agregar los 2 ejemplos de arriba.
  - Conservar «la API no cambia el otro campo por su cuenta». Que la API deduzca `false` sería una regla nueva: si se quiere, es PREGUNTA para Axel, no la decide el corrector.
  - En `verify_c1.py`, agregar una comprobación como la de rangos de `parse_schema.py`: probar contra la base, en un `SAVEPOINT`, las combinaciones (true, n), (NULL, n) y (false, n), y exigir que cada combinación rechazada aparezca en el texto de la regla.
  - Agregar el caso NULL a la solicitud para S6 del reporte de corrección.

### B-1 (baja) · faltan la fila y la nota de la corrección r1 en `C1-datos.md`
- **HECHO.** La tabla «Cambios» (l. 16–28) tiene una fila por corrección en las versiones anteriores (v1.0 r1 y r2; v1.1 r1, r2 y r4; v1.2 r1 y r3), pero no tiene «C1 v1.3 · corrección r1».
- **HECHO.** El punto de IVA del detalle v1.3 (l. 36) no menciona la regla 422.
- **Ruta:** `backend/internal/api/c1gen/c1_plantilla.md`.
- **Corrección sugerida:** agregar la fila (H-1, H-2, H-3, las 3 comprobaciones nuevas y la guarda de `_ins`) y una frase sobre el 422 en l. 36.

### B-2 (baja) · resumen impreciso del orden de los grupos
- **HECHO.** `C1-datos.md` l. 38 dice «Orden de los grupos: mercado, cadena y el grupo sin tienda al final». `C1-puerta2.md` dice lo mismo.
- **HECHO.** `openapi.yaml` (`ShoppingList.groups.items`) y `S4-interfaz.md` l. 216 ponen el grupo sin tienda al final **de cada vertical**. Así, el grupo «sin mercado dado de alta» va antes que los grupos de cadena.
- **Ruta:** `backend/internal/api/c1gen/c1_plantilla.md`.
- **Corrección sugerida:** «primero mercado y luego cadena; dentro de cada vertical, el grupo sin tienda al final».

### B-3 (baja) · `params_json` acepta claves que ya viven en `app_setting`
- **HECHO.** `Notification.params_json` y `NotificationUpdate.params_json` declaran `additionalProperties: true`. Su texto dice que la hora, el día y la anticipación ya no van ahí, pero no dice qué hace la API si un cliente de v1.2 todavía manda `send_time`, `advance_days`, `weekday` o `day_of_month`. ¿Responde 422, o lo guarda sin efecto?
- **Ruta:** `gen2_schemas.py` (texto de `params_json`).
- **Corrección sugerida:** decir que esas 4 claves responden 422 `validation_failed` con `fields[].field = params_json.<clave>`, o declarar la lista de claves propias que se admiten (hoy solo `scopes`).

### B-4 (baja, fuera de C1) · `Makefile`: el TMP fijo se comparte entre corridas en paralelo
- **HECHO.** `TMP := $(TMP_ROOT)/verify-wf0` y `$(RM_TMP)` borra ese directorio al terminar. Con varios agentes a la vez, una corrida borra el log de otra (V-1). Con un TMP propio, la misma verificación pasa (V-2 y V-5).
- **Corrección sugerida** (integrador): un TMP por corrida (`mktemp -d` bajo `datos/tmp`), o escribir en el protocolo que en P2c se pasa `TMP_ROOT=` con una ruta propia.

### B-5 (baja) · los bajos de la ronda 1 siguen abiertos y sin cambio
No venían en el encargo de `correccion-C1-r1`:
- B-1 de r1: la marca de más en `vat_included`.
- B-2 de r1: el texto de `vat_on_interest_bp` en `Credit` y `CreditCreate` todavía empieza con «NULL = la tasa o la cotización ya incluye IVA (Kueski)». En 1.4, ese NULL significa «incluye IVA» solo con `vat_included` = 1.
- B-3 de r1: anotar en `S3-iva-lectura` lo que tomó C1.
- B-4 de r1: `docs/criterios/C1.md` sigue en la versión 1.2.
- B-5 de r1: `bloqueos/INDICE.md` no tiene `C1-reparto-alcance` ni `C1-hora-reportes-job`.
- Además, `C1-puerta2.md` todavía cita «comprobaciones OK: 29». Hoy son 32.

## Estado de los hallazgos de la ronda 1

| Hallazgo r1 | Estado | Evidencia propia |
|---|---|---|
| H-1 (media), C1-14: `visit_date_` y `cheapest_store_id_` | **CERRADO** | V-13. Los nombres están bien en `openapi.yaml`, `compras-tiendas-a.yaml` y `reportes.yaml`; hay 0 propiedades con «_» final y la guarda en `_ins` |
| H-2 (media), Q2-65 A y C1-09: M-11 sin los acumulados | **CERRADO** | V-14. La fila nombra los campos, su fórmula, el origen y el ejemplo, con cifras que recalculé desde `model.json` |
| H-3 (media), Q2-67 opción 1 y C1-04: respuesta a `{vat_included: true, vat_on_interest_bp: 1600}` | **CERRADO en lo que pedía** | V-9 y V-15. Están documentados el 422, `fields[].field` y el Merge Patch sobre lo guardado. La otra combinación que rechaza el CHECK (NULL o ausente con porcentaje) queda abierta como **M-1** |
