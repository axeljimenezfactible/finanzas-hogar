# Revisión C1 v1.2 · ronda 1 de 3 (flujo P1c, puerta 1b)

Revisor independiente revisor-C1-r1, 11-oct-2026. Evalúa solo contra `docs/criterios/C1.md` (versión 1.2), `docs/plan/decisiones.md` §K (K36, K38, K43, K44) y el hallazgo R3-01. No corrigió nada. Todas las verificaciones se reejecutaron con comandos y scripts propios en `datos/tmp/rc1/` (borrados al terminar).

## Veredicto

**NO APROBADO.** Hay 1 hallazgo alto (A-01) y ningún medio. A-01: en C1-23, `send_time` no tiene el patrón `HH:MM` y `advance_days` no tiene `minimum: 0` en el contrato. Solo aparecen en el texto de `NotificationUpdate`, y `params_json` es un objeto libre. Lo demás de la puerta 1b sí cumple: K36, K38, K44 con sus 15 claves, los 10 MB y 413, los 5 s y 503 `unavailable`, las marcas quitadas, los extractos y las verificaciones (a)–(j). Hay además 5 hallazgos bajos que no bloquean. Uno de ellos (B-01) dice que el texto de C1-22, C1-13, C1-21 y C1-12 en los criterios ya no corresponde a lo vigente. Ahí manda K44 y no se cuenta contra C1.

## Verificaciones reejecutadas

| # | Comando propio | Resultado |
|---|---|---|
| a | `python3 -I backend/internal/api/c1gen/build.py --check` | «los 21 archivos coinciden byte a byte con el generador», IDÉNTICO en los 21; `git diff --stat` igual antes y después |
| b | goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/rc1/c1.db up`, luego `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/rc1/c1.db` | 00001, 00002 y 00003 OK, versión 3. **RESULTADO: TODO OK** (27 OK). Sin `--db` también da TODO OK. Con `python3 -I` del sistema falla C1-14 porque falta `jsonschema` (B-05) |
| c | `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator openapi.yaml extractos/*.yaml` | 18 de 18 OK |
| d | `python3 -I docs/design/tools/verificar_u1.py` | TODO OK; `grep -c PENDIENTE` = 0; openapi 1.2.0, 210 operaciones, 209 con x-necesidades |
| e | `cd backend && go build ./... && go vet ./... && gofmt -l .` | go1.27.2, BUILD_OK, VET_OK, gofmt vacío |
| f | script propio `check.py` | openapi tiene 210 operationId, 210 únicas. Los extractos suman 210: 0 faltan, 0 sobran, 0 repetidas. Ninguno pasa de 80 000 B; el mayor es gastos.yaml con 76 220 B. Las filas de INDICE.md suman 210 y sus operaciones y bytes cuadran con los archivos. `sha256sum extractos/* \| sha256sum` = `e7fb8de8…69bfa` |
| g | `check.py` y `check3.py` (sqlite) | Las 10 claves de `app_setting` que siembra 00003 están en openapi.yaml y en C1-datos.md, y también todas las sembradas en 00001 a 00003. Las 15 `@param` tienen 1 lectura con `value` e `initial_value`, 1 cambio y 1 `reset_parameters`. Todas son `integer` o `array` con `minimum: 0`, sin `maximum` y sin `x-provisional`. Las 20 descripciones de `app_setting` citan su decisión (`decisión X` o `Kn`) |
| h | `check2.py` | 0 marcas que citen `U1-lista-vigente-domingo`, `C1-valores-k31`, `S1-ajuste-objetivo-kcal` o `criterios-promedios-redondeo`. La única mención está en una fila histórica de Cambios (C1-datos.md L24, v1.1 r4), que narra la v1.1 y no es una marca |
| i | `check2.py` | Hay 3 `x-provisional` en openapi: `Notification.params_json`, `NotificationUpdate.params_json` y `NotificationUpdate`. Las 3 citan `S1-parametros-avisos`, que no tiene «Respuesta de Axel». En C1-datos.md, las líneas PROVISIONAL vigentes (L570, L574, L770–773 y §13) citan solo ese bloqueo. Las demás son leyendas o filas históricas de Cambios |
| j | `check2.py` y `sed -n` | `info.version` = 1.2.0. C1-datos.md trae la cabecera «C1 v1.2», la fila «C1 v1.2 (puerta 1b)» en Cambios (L25) y el detalle en la sección de cambios (L27–35) |
| k | `check3.py` | Hay 14 ejemplos, los mismos de v1.1. `MoneyCents` 97375 sale del oráculo (17 coincidencias en `oraculo/model.json`). Los demás son marcadores o formatos (`2026-W42`, `5000`, `1000`, códigos). Los valores iniciales de las descripciones nuevas salen de la siembra: parse_schema los compara y verify_c1 dice 15 de 15 |

## Criterios

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| C1-01 | alta | Cumple | (c); verify_c1 «C1-01 openapi = 3.1.1» |
| C1-02 | alta | Cumple | 211 N-xxx; 208 con operación; 3 Local (N-003, N-301, N-317). N-320 a N-327 tienen operación. §10 tiene 236 filas, igual a U1 §3.29 |
| C1-03 | alta | Cumple | 55 de 55 tablas; 1011 propiedades sin huérfanas |
| C1-04 | alta | Cumple | 953 de 953 respuestas de error van a `Error`. `payload_too_large` (413) y `unavailable` (503) están en el catálogo de §3 (L121 y L125) |
| C1-05 | alta | Cumple | 40 de 40 con `Page` |
| C1-06 | alta | Cumple | R3-01: las 8 operaciones de `shopping_list`, más `getBasketComparison`, `listBaskets` y `getDashboardBestBasket`, piden `week` obligatoria (`in: query`, `required: true`) |
| C1-07 | alta | Cumple | Solo `date` y `date-time` |
| C1-08 | alta | Cumple | 265 `_cents` van a MoneyCents; 0 `number` |
| C1-09 | alta | Cumple (sin cambios en v1.2) | M-05 trae «máx. 21; K43». Sigue el bajo H-r3-05 de M-06 |
| C1-10, C1-11 | alta | Cumple | verify_c1 TODO OK. Cada sección nueva de K44 tiene su GET y su PATCH de settings |
| C1-12 | alta | Cumple | Solo login y salud son públicos. La única clave con «password» en una respuesta es `auth.password_min_length`, un entero de configuración. El literal del criterio no la contempla (B-01) |
| C1-13 | media | Cumple con observación | C1 v1.2 y 1.2.0; 0 operationId duplicadas; 9 etiquetas. El criterio todavía dice «C1 v1.1» (B-01) |
| C1-14 | media | Cumple | 14 de 14 ejemplos válidos con jsonschema; 2871 descripciones en español. Sigue el bajo H6 |
| C1-15 | alta | Cumple con observación | 3 marcas, todas del bloqueo abierto `S1-parametros-avisos`. 0 marcas sobre K6, K7, K36, K38, K43 o K44. La marca a nivel de esquema en `NotificationUpdate` es más amplia de lo necesario (B-02) |
| C1-16 a C1-20 | alta | Cumple (sin cambios en v1.2) | verify_c1 TODO OK; no se reevaluaron a fondo en esta ronda |
| C1-21 | alta | Cumple | (a) y (f) |
| C1-22 | alta | Cumple con observación | Tabla clave → valor inicial → operación en §11.7, con 15 claves. Ejemplos: `auth.max_failed_attempts` 5, `auth.lockout_window_minutes` 15, `auth.password_min_length` 12 y `session.max_age_days` 90 van por `getSecuritySettings`/`updateSecuritySettings`; `price.stale_after_days` 30, `price.zero_weight_after_days` 90 y `price.profeco_import_max_gb` 1 por `getPriceSettings`/`updatePriceSettings`; `comparator.timeout_seconds` 5 por `getOptimizerSettings`/`updateOptimizerSettings`; `attachment.max_size_mb` 10 e `idempotency.retention_hours` 24 por `getDataSettings`/`updateDataSettings`; `health.goal_adjust_lose_kcal` 0 y `health.goal_adjust_gain_kcal` 0 por `getHealthSettings`/`updateHealthSettings`; `alerts.budget_warning_pct` 80, `alerts.budget_limit_pct` 100 y `alerts.payment_due_days_before` [1, 0] por `getAlertSettings`/`updateAlertSettings`. Avisos: los umbrales y los días antes del pago viven en `app_setting`, como dicen K44, la 00003 y el bloqueo reducido (B-01). `send_time` y `advance_days` solo están declarados en texto (A-01) |
| C1-23 | alta | **No cumple** | 10 MB y 413 `payload_too_large` OK. 503 `ServiceUnavailable` (`unavailable`) en `compareScenarios` OK. N-023, N-064, N-226, N-297 y N-300 son «FIRME (K43)» en §10. `grep PROVISIONAL.*C1-valores-k31` = 0. M-05 con 21 OK. **Falla**: no hay propiedad `send_time` con patrón ni `advance_days` con `minimum: 0` (0 coincidencias de `[01][0-9]` en openapi y en C1-datos) |

Decisiones: **K36** cumple. `HealthProfile`, `HealthProfileCreate` y `HealthProfileUpdate` (este último nullable) traen `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal`; `HealthProfile` trae `household_*`. `suggested_kcal_target` dice RN-17 1.3: se resta o se suma según el objetivo, `maintain` no ajusta y hay un solo ROUND_HALF_UP. `getHealthSettings`/`updateHealthSettings` cubren N-326 y las operaciones de perfil cubren N-327. **K38** cumple. `getMealPlanWeek` no tiene marca; su descripción dice que los domingos abre la semana siguiente y que la lista y comparar canasta llevan `week`. N-120 es FIRME (K38) y §13 no tiene la fila. **K43** cumple, salvo A-01. **K44** cumple: 15 de 15 claves, sin topes ni `x-provisional`, con `reset_parameters`.

## Hallazgos

| Id | Sev. | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| A-01 | alta | C1-23 (y C1-22, avisos) | `send_time` (HH:MM) y `advance_days` (≥ 0) solo están en el texto de `NotificationUpdate`. `Notification.params_json` y `NotificationUpdate.params_json` son `{type: object, additionalProperties: true}` sin propiedades. La verificación de C1-23 pide el patrón `^([01][0-9]\|2[0-3]):[0-5][0-9]$` y `minimum: 0`, y da 0 coincidencias. `verify_c1.py` no lo comprueba, por eso pasa | `backend/internal/api/c1gen/gen2_schemas.py` (`NotificationUpdate`, `Notification`) → `docs/contratos/openapi.yaml`, extractos `ajustes-d.yaml`; `backend/internal/api/c1gen/verify_c1.py` | En `params_json`, declarar `properties`: `send_time` con `type: string` y ese `pattern`, y `advance_days` con `type: integer` y `minimum: 0`, cada una con descripción que cite K43. Conservar `additionalProperties: true` para `weekday` y `day_of_month`, y la marca de lugar `x-provisional` que cita `S1-parametros-avisos`. Sin valor por omisión. Sumar a verify_c1 la comprobación de C1-23 (patrón y mínimo). Regenerar con `build.py` |
| B-01 | baja | C1-22, C1-13, C1-21, C1-12 (texto del criterio) | `docs/criterios/C1.md` 1.2 tiene textos que ya no corresponden. C1-22 pone los umbrales 80/100 y [1, 0] en `params_json` con PROVISIONAL y prohíbe sembrar `alerts.*`, pero K44 los nombra como claves de `app_setting`, la 00003 los siembra, el bloqueo abierto se redujo y ya no los incluye, y U1 1.2 r2 ya los toma de `app_setting`. C1-13 dice «C1 v1.1». C1-21 dice «hoy 201». El grep literal de C1-12 contaría `auth.password_min_length`. Manda K44 y no se cuenta contra C1 | `docs/criterios/C1.md` | Fuera del alcance de C1: el agente de criterios actualiza C1-22 (excepción solo para `send_time`, `advance_days`, `weekday` y `day_of_month`), C1-13 (v1.2 y 1.2.0), C1-21 (210) y C1-12 (excepción `x-app-setting`). Es la solicitud 2 del reporte de C1 |
| B-02 | baja | C1-15, C1-22 | `NotificationUpdate` lleva `x-provisional: true` a nivel de esquema, y eso cubre también `is_enabled`, que está decidido. C1-22 pide que la marca cubra solo el lugar | `gen2_schemas.py` (`S['NotificationUpdate']['x-provisional'] = True`) | Dejar la marca solo en `params_json` (ya la tiene) |
| B-03 | baja | C1-14 | `Notification.params_json` dice «Parámetros del aviso (días antes 1 y 0, umbrales 80 y 100) Desde 1.3 … ya no lleva `days_before`…». El comentario de 00001 y la `@col_note` quedan pegados sin punto, y la primera frase afirma lo que la segunda niega | `backend/internal/api/c1gen/parse_schema.py` (unión de `@col_note`) | Separar con «. » y presentar la nota como cambio, por ejemplo «Hasta 1.2: (…). Desde 1.3: …», o dejar que la nota sustituya el paréntesis |
| B-04 | baja | C1-15 | verify_c1 imprime «OK C1-15 la regla del domingo («los domingos») va marcada PROVISIONAL con el bloqueo U1-lista-vigente-domingo», aunque el bloqueo está resuelto y la regla es FIRME (K38). La comprobación está inactiva y la etiqueta dice lo contrario del estado actual | `verify_c1.py` | Si el bloqueo está resuelto, que la etiqueta diga «regla del domingo FIRME (K38): 0 marcas que lo citen» |
| B-05 | baja | Protocolo §5 | El reporte dice que `python3 -I verify_c1.py --db` da TODO OK, pero con el Python del sistema falla C1-14 por falta de `jsonschema` | `docs/estado/p1c/C1-v1.2.md` §Verificación | Citar el comando que se usó: `uv run --no-project --with pyyaml --with jsonschema python3 -I …` |

## Estado de los anteriores (P1, `docs/estado/p1/REVISION-C1-r4.md`, y R3-01)

| Hallazgo | Estado | Evidencia |
|---|---|---|
| R3-01 (`week` explícita en lista y canasta) | Cerrado y vigente | Ver C1-06; K38 lo confirma |
| H-r3-05 (M-06 «lista vigente») | Sigue (bajo) | C1-datos.md L341: M-06, periodo «lista vigente» |
| H-r3-06 (IVA, K18) | Sigue (bajo) | `Credit.vat_on_interest_bp`: «NULL = … ya incluye IVA» |
| H6 (descripciones que repiten el identificador) | Sigue (bajo) | En la muestra de verify_c1 salen `id`, `created_at` y `updated_at` |
| H7 (marca sin `x-provisional`) | Cerrado | `DashboardWeekMeals` ya no tiene marca (K43). `NotificationUpdate` tiene `x-provisional` que cita el bloqueo abierto |
| H8 (`LoginRequest.username` sin 3–64) | Sigue (bajo) | `{type: string, description}` sin `minLength`/`maxLength` |
| N-r4-01 a N-r4-03 | No reevaluados en esta ronda | — |
