# Revisión C1 v1.2 · ronda 2 de 3 (flujo P1c, puerta 1b)

Revisor independiente `revisor-C1-r2`, 11-oct-2026. Evalúo contra `docs/criterios/C1.md` (versión 1.2), contra las decisiones K36, K38, K43 y K44 de `docs/plan/decisiones.md` §K y contra el hallazgo R3-01. No corregí nada. Corrí todas las verificaciones con mis propios comandos y con tres scripts míos en `datos/tmp/rc2/`, que ya borré. HEAD es `082235c`. Revisé el árbol de trabajo tal como lo dejó `correccion-C1-r1`.

## Veredicto

**APROBADO.** No hay hallazgos altos ni medios. A-01, el único alto de la ronda 1, quedó **cerrado**. `send_time` y `advance_days` ya son propiedades de `params_json` en `Notification` y en `NotificationUpdate`, tanto en `openapi.yaml` como en el extracto `ajustes-d.yaml`. `send_time` lleva el patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$` y `advance_days` lleva `minimum: 0`. Ninguna de las dos trae valor por omisión. `C1-datos.md` §11.7 documenta las dos validaciones. Probé el esquema con jsonschema: acepta `07:30` y `23:59`, y rechaza `24:00`, `7:30`, `-1` y `1.5`. Quedan 8 hallazgos bajos: 5 vienen de la ronda 1 y 3 son nuevos. Ninguno bloquea.

## Verificaciones reejecutadas

| # | Comando (desde la raíz del repo) | Resultado |
|---|---|---|
| a | `python3 -I backend/internal/api/c1gen/build.py --check` | **IDÉNTICO en los 21 archivos**: openapi.yaml de 674 710 B, routes.go, C1-datos.md de 143 281 B, INDICE.md, 16 extractos y componentes-comunes.yaml. Reporta 167 rutas, 210 operaciones, 252 esquemas y 15 parámetros K44 |
| b | goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/rc2/c1.db up`, luego `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/rc2/c1.db` | 00001, 00002 y 00003 OK; versión 3. **RESULTADO: TODO OK**, con las dos líneas nuevas de C1-23 en OK. Leí el código de C1-23 en `verify_c1.py` (L391–L420): compara el tipo, el patrón exacto, `minimum == 0`, que no haya `default`, que la marca de lugar exista en openapi y en el extracto, y la línea de §11.7. No es una comprobación tautológica |
| c | `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator` sobre `openapi.yaml` y sobre los 17 yaml de `extractos/` | 18 de 18 **OK**. `openapi: 3.1.1`, 167 rutas, `info.version` 1.2.0 |
| d | `python3 -I docs/design/tools/verificar_u1.py` | **TODO OK**, 0 PENDIENTE. 210 operaciones; 209 con `x-necesidades`; la única excepción es `GET /openapi.yaml` |
| e | `go build ./...`, `go vet ./...` y `gofmt -l .` en `backend/` | build OK, vet OK, gofmt sin salida. Go 1.27.2 |
| f | script propio | openapi tiene 210 operationId y 0 duplicadas. Los extractos suman 210: 0 faltan, 0 sobran y 0 repetidas. Ninguno pasa de 80 000 B; el mayor es gastos.yaml con 76 220 B. Las 17 filas de INDICE.md coinciden en operaciones y bytes con cada archivo, y suman 210 |
| g | script propio: claves de `app_setting` en la base migrada y `@param` de 00003, contra el contrato | Hay 15 `@param`. Las 15 están sembradas con su valor y su comentario cita la decisión (`decisión A/E/F` o K-n). Cada una tiene un nodo `x-app-setting` en un esquema de lectura de settings, `SettingParameterInteger` o `SettingParameterIntegerList`, con `value` e `initial_value`. También está en su `*SettingsUpdate` como entero o lista de enteros con `minimum: 0` y sin `maximum`. 0 `x-provisional`. Las 9 claves que siembra 00003 aparecen en openapi.yaml y en §11.7. Las 21 claves de `app_setting` aparecen en el contrato; solo `ai.*`, que son interruptores de la decisión H y no parámetros K44, no tienen nodo `x-app-setting`. K36: `HealthProfile`, `HealthProfileCreate` y `HealthProfileUpdate` traen `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` (en Update aceptan `null`); `HealthProfile` trae además `household_goal_adjust_*` y `suggested_kcal_target` según RN-17 1.3, «maintain no ajusta» y un solo HALF_UP. `getHealthSettings` y `updateHealthSettings` existen (N-326) |
| h | `grep -rn` de `U1-lista-vigente-domingo`, `C1-valores-k31` y `S1-ajuste-objetivo-kcal` en openapi.yaml, C1-datos.md, routes.go y extractos | **0 marcas.** La única coincidencia es C1-datos.md L24, una fila histórica de Cambios (v1.1 r4). La fila v1.2 (L25) la corrige: «K38, K39 y K43: sin marcas… N-120 FIRME». `getMealPlanWeek` no lleva `x-provisional` y su descripción trae la regla de K38 («los domingos, la siguiente (la que empieza el lunes)… (K38)»). En §10, N-120 está en «FIRME (K38)» |
| i | script propio sobre todos los `x-provisional` de openapi y las líneas PROVISIONAL de openapi y C1-datos.md | Hay 3 nodos `x-provisional`: `Notification.params_json`, `NotificationUpdate` (a nivel de esquema) y `NotificationUpdate.params_json`. Los 3, y las 3 líneas PROVISIONAL de openapi, citan `S1-parametros-avisos`, que existe y no tiene «Respuesta de Axel». En C1-datos.md, todas las marcas vigentes (§10 N-297 y N-300, §11.7 y §13) citan ese mismo bloqueo |
| j | lectura de openapi `info` y de C1-datos.md §Cambios | `info.version` 1.2.0 y summary/description «C1 v1.2». Están las filas «C1 v1.2 (puerta 1b)» (L25) y «C1 v1.2 · corrección r1» (L26), y la subsección «C1 v1.2 (puerta 1b)» |
| k | script propio sobre los 14 `example` y sobre las cifras con decimales de las líneas que agrega `git diff HEAD` en openapi.yaml y C1-datos.md | Ejemplos: 97375, que es 973.75 y está 17 veces en `oraculo/model.json`; 2026-10-30, que es el próximo pago de C1-09; y marcadores de formato (5000 pb = 50 %, 1000 milésimas, 2026-W42). Las líneas agregadas solo traen valores de §K: 10 MB, 5 s, 24 h, 5 intentos y 90 días. El único número con decimales, «3.29», es la referencia «U1 §3.29». **0 cifras inventadas** |
| R3-01 | script propio | Las 8 operaciones que tocan `shopping_list` o `shopping_list_item` (`getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getBasketComparison` y `getDashboardBestBasket`) piden `week` obligatoria |
| extra | script propio | C1-04: 953 de 953 respuestas de error van a `Error`. C1-08: 0 nodos `type: number`. C1-12: solo `login` y `getHealth` son públicas. C1-13: 0 operationId duplicadas y 9 etiquetas. Las 9 operaciones nuevas exigen sesión y tienen etiqueta «Ajustes», `x-necesidades` (N-320, N-323 a N-326) y `x-tablas`. `uploadExpenseAttachment` dice «10 MB (K43)» y documenta el 413 `PayloadTooLarge`. `compareScenarios` dice «Tiempo máximo 5 s (K43)» y documenta el 503 `ServiceUnavailable` |

## Criterios de C1.md 1.2

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| C1-01 | alta | Cumple | (c) |
| C1-02 | alta | Cumple | (a), (b) y (d): 0 N sin operación ni Local, 0 operaciones sin N salvo `getOpenApiDocument` (justificada en §11.2), 0 N desconocidas |
| C1-03 | alta | Cumple | build: «tablas sin op: [] · tablas inexistentes: []»; verify_c1 TODO OK |
| C1-04 | alta | Cumple | (extra) 953 de 953 |
| C1-05 | alta | Cumple | verify_c1 TODO OK; v1.2 no cambia la paginación |
| C1-06 | alta | Cumple | verify_c1 TODO OK; R3-01 (8 de 8 con `week` obligatoria) |
| C1-07 | alta | Cumple | verify_c1 TODO OK; los ejemplos de fecha son ISO y DateTime trae `-06:00` |
| C1-08 | alta | Cumple | (extra) 0 `number`; 97375 sale del oráculo |
| C1-09 | alta | Cumple | M-05 trae «máx. 21; K43»; los ejemplos vienen del oráculo |
| C1-10, C1-11 | alta | Cumple | verify_c1 TODO OK; cada sección de K44 tiene su GET y su PATCH |
| C1-12 | alta | Cumple | Solo login y salud son públicas; 0 campos sensibles. `auth.password_min_length` es un entero de configuración (B-01) |
| C1-13 | media | Cumple con observación | `info.version` 1.2.0, ids únicas y 9 etiquetas. Dos encabezados generados todavía dicen v1.1 y v1.0 (N-01, baja) |
| C1-14 | media | Cumple | 14 de 14 ejemplos válidos; 2875 descripciones en español |
| C1-15 | alta | Cumple | (h) e (i): 3 marcas, todas del bloqueo abierto; 0 marcas sobre K6, K7, K36, K37, K38, K43 o K44. El reporte tiene SUPUESTOS y ninguno cambia endpoints sin bloqueo |
| C1-16 a C1-20 | alta | Cumple (sin cambios en v1.2) | verify_c1 TODO OK |
| C1-21 | alta | Cumple | (a) y (f) |
| C1-22 | alta | Cumple con observación | (g): 15 claves con lectura (`value` e `initial_value`), cambio y restablecer, sin mínimos ni máximos de seguridad, y la tabla de §11.7. Los umbrales 80/100 y [1, 0] están como claves `alerts.*` de `app_setting` (`getAlertSettings`), sin marca. Eso sigue K44, que los nombra como claves de `app_setting`, y el bloqueo reducido `S1-parametros-avisos`. El texto literal de C1-22 los pone en `params_json` con PROVISIONAL; ese texto quedó desfasado (B-01) |
| C1-23 | alta | **Cumple** (A-01 cerrado) | 10 MB con 413, 503 `unavailable`, N-023, N-064, N-226, N-297 y N-300 en «FIRME (K43)», M-05 con 21, `grep PROVISIONAL.*C1-valores-k31` = 0, y el patrón y el `minimum: 0` en openapi, en el extracto y en C1-datos.md |

## Hallazgos

| Id | Sev. | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| N-01 | baja | C1-13 | El encabezado generado de `openapi.yaml` (L1) dice «Contrato C1 v1.1» y el de `routes.go` (L1) dice «contrato C1 v1.0». `info.version` es 1.2.0 y los extractos ya dicen «C1 v1.2.0» | `backend/internal/api/c1gen/gen4_build.py` L111; `backend/internal/api/c1gen/gen5_go.py` L15 | Sacar la versión de una sola constante (la de `info.version`) en los dos encabezados y regenerar |
| N-02 | baja | C1-03 y K36 | `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` de `HealthProfileCreate` y `HealthProfileUpdate` no declaran `minimum: 0`. K36 dice «magnitud en kcal ≥ 0» y la base tiene `CHECK (… >= 0)`. El valor de la casa (`SettingParameterInteger.value`) sí lo declara. Hoy el contrato no lleva los CHECK a `minimum` (tampoco en `height_mm` ni en `daily_kcal_target`); por eso queda como mejora | `backend/internal/api/c1gen/gen2_schemas.py` (HealthProfile*) → `openapi.yaml` y `extractos/salud.yaml` | Agregar `minimum: 0` a esas dos propiedades en Create y Update, o decir «entero ≥ 0» en la descripción |
| N-03 | baja | C1-15 | La descripción del estado de la tarjeta (N-205, openapi L15851) remite a `docs/estado/bloqueos/criterios-credito-formulas.md` como fuente de fórmulas. Ese bloqueo está RESUELTO por K18. No es una marca PROVISIONAL, pero apunta a un bloqueo como si fuera la fuente | `backend/internal/api/c1gen/gen2_schemas.py` → `openapi.yaml` | Citar K18 (y RN o S3-08/S3-10) en vez del archivo de bloqueo |

## Estado de los hallazgos de la ronda 1

| Hallazgo | Estado | Evidencia |
|---|---|---|
| A-01 (alta) | **Cerrado** | Verificación (b), script propio sobre openapi y `ajustes-d.yaml`, prueba con jsonschema y §11.7 de C1-datos.md |
| B-01 (baja) | Sigue abierto, fuera del alcance de C1 | `docs/criterios/C1.md` sigue en 1.2 con el texto viejo de C1-22 (umbrales en `params_json` y prohibición de sembrar `alerts.*`), C1-13 («C1 v1.1») y C1-21 («hoy 201»). Le toca al agente de criterios |
| B-02 (baja) | Sigue abierto | `gen2_schemas.py` L361 `S['NotificationUpdate']['x-provisional'] = True`: la marca a nivel de esquema también cubre `is_enabled` |
| B-03 (baja) | Sigue abierto | `Notification.params_json` dice «Parámetros del aviso (días antes 1 y 0, umbrales 80 y 100) Desde 1.3…». No hay punto entre las frases y la primera afirma lo que la segunda niega |
| B-04 (baja) | Sigue abierto | verify_c1 imprime «OK C1-15 la regla del domingo… va marcada PROVISIONAL con el bloqueo U1-lista-vigente-domingo» aunque K38 la hizo FIRME |
| B-05 (baja) | Parcial | `correccion-C1-r1.md` ya cita el comando con `uv run … --with jsonschema`. `C1-v1.2.md` §Verificación sigue sin ese comando |
