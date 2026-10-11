# REVISIÓN C1 v1.1 · ronda 1 de 3 (flujo P1, 11-oct-2026)

Revisor: `revisor-C1-r1` (independiente; no escribió nada de lo revisado y no corrigió nada).
Contra qué: `docs/criterios/C1.md` versión 1.1 (C1-01 a C1-21), decisiones §K (K1–K7, K9, K11–K13, K16, K18, K21, K23) y los hallazgos R2-06 y R3-01 de `docs/estado/wf0/REVISION-r3.md`.
Objeto: `docs/contratos/openapi.yaml` (info.version 1.1.0), `docs/contratos/C1-datos.md`, `backend/internal/api/routes.go`, `docs/contratos/extractos/` y la fuente `backend/internal/api/c1gen/` (git diff HEAD).

## Veredicto

**NO APROBADO.** 2 hallazgos de severidad alta, 3 media y 3 baja. Todas las verificaciones automáticas pasan (generador idéntico, base migrada, validador, verificadores, Go, extractos). Lo que falla es documentación que los criterios C1-16 y C1-19 exigen **en `C1-datos.md`** y que hoy solo vive (en parte) en las descripciones de `openapi.yaml`. Las correcciones son de texto en `c1_plantilla.md` y del generador. No cambian rutas, cifras ni el esquema de S1.

## Verificaciones reejecutadas (comandos propios; temporales en `datos/tmp/rev-c1`, ya borrados)

| # | Comando | Resultado |
|---|---|---|
| a | `python3 -I backend/internal/api/c1gen/build.py --check` | IDÉNTICO en los 21 archivos: `openapi.yaml` (629 771 bytes), `routes.go`, `C1-datos.md` y los 18 de `extractos/`. «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| b | goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/rev-c1/c1.db up` y `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/rev-c1/c1.db` | 00001 OK, 00002 OK, versión 2. **TODO OK**, 0 FALLA. C1-02: 203 N-xxx, 0 sin operación ni Local (Local N-003, N-301, N-317), matriz §10 de 227 pares = U1 §3.29. C1-03: 55 tablas, 0 inexistentes, 0 sin operación. C1-04: 915 de 915 errores → `Error`. C1-05: 40 de 40 colecciones. C1-15: 6 marcas, todas con bloqueo abierto. EXT y S6-03: 201 = 201 |
| c | `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` y lo mismo con los 17 YAML de `extractos/` | `openapi.yaml: OK`, 17 de 17 extractos OK |
| d | `python3 -I docs/design/tools/verificar_u1.py` | TODO OK. U1-03: 0 N-xxx sin operación ni Local. 0 operaciones sin N-xxx fuera de `GET /openapi.yaml`. N-316 a N-319 sin reusar |
| e | `cd backend && go build ./... && go vet ./... && gofmt -l .` (go1.27.2) | build OK, vet OK, gofmt vacío |
| f | script propio `ext.py` | 18 archivos (17 YAML + `INDICE.md`). Suma de operationId en los extractos = 201 = `grep -c operationId openapi.yaml` (201). 0 repetidas, 0 faltantes, 0 sobrantes. Mayor 76 220 bytes (≤ 81 920). `componentes-comunes.yaml` presente (con `cookieAuth` y 14 parámetros). 17 de 17 cargan como YAML. 0 `$ref` internas sin resolver. `build.py --extractos` corrido dos veces: sha256 de la carpeta igual (`c6b365bc…eb243`) antes y después |
| g | scripts propios `k.py`, `k2.py` y `d.py` (lista K) | Ver tabla «Decisiones §K» |
| h | `grep -n "PROVISIONAL\|x-provisional"` en `openapi.yaml` y `C1-datos.md` y `grep -c "Respuesta de Axel"` en los bloqueos citados | `openapi.yaml`: 6 marcas (la línea 9 es la definición de la extensión). Citan solo `C1-valores-k31`, `criterios-promedios-redondeo` y `S1-ajuste-objetivo-kcal`, los tres con 0 «Respuesta de Axel». `C1-datos.md`: 10 filas de la matriz (9 N-xxx únicas, igual que la cifra de §11.1) y 3 filas de §13, con los mismos bloqueos. 0 marcas sobre los valores de K6 o K7 |
| i | `grep` de `info.version` y de Cambios | `info.version: 1.1.0`. Fila «C1 v1.1 (puerta 1)» en Cambios con su subsección. Encabezado «C1 v1.1» |
| j | script propio `ex.py` (cifras de ejemplos contra `oraculo/model.json` y `datos/semilla/data.json`) | 14 nodos `example(s)`, 20 números. Solo queda sin origen `Error.request_id` = `123456789`: es un identificador marcador, no una cifra. En `C1-datos.md`, 36 cifras con decimales: todas en oráculo o semilla, salvo números de sección (26.01 a 26.14, 3.29), 138.50 (diferencia canon de la semilla, 3,233.00 − 3,094.50, citada en U1) y 25.8685 (la tasa del oráculo 0.2586847632897075 redondeada HALF_UP a 4 decimales, K12). Ninguna inventada |

## Criterios de C1.md 1.1

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| C1-01 OpenAPI 3.1 válido | alta | Cumple | 3.1.1, 162 rutas, 201 operaciones. Validador 0.9.0 OK (c) |
| C1-02 N-xxx ↔ operación | alta | Cumple | (b) y (d). Única operación sin N-xxx: `GET /openapi.yaml` (excepción por estáticos) |
| C1-03 Tablas de S1 | alta | Cumple | 55 de 55 tablas, 1003 propiedades con columna o campo calculado (b) |
| C1-04 Errores | alta | Cumple | 915 de 915. 409 `idempotency_key_reused` documentado |
| C1-05 Paginación | alta | Cumple | 40 de 40. `Limit` 50 por omisión (U1 N-040 «Páginas de 50») y máximo 200 (K6) |
| C1-06 Filtros | alta | Cumple | 6 parámetros únicos. `Family` con enum `household`/`external`/`joint`, inclusión descrita y sin `x-provisional`. `Period` con K8 |
| C1-07 Fechas | alta | Cumple con observación | Solo `date` y `date-time`. Ningún ejemplo de respuesta termina en «Z». Nota de «Z» en §4 y en `DateTime`. Hay 22 descripciones de campo que dicen «(UTC)» (H4) |
| C1-08 Dinero | alta | Cumple | 265 propiedades `_cents` por `$ref` a `MoneyCents`. 0 `type: number` |
| C1-09 Métricas | alta | Cumple | §9 con 10,000.00, 6,905.50, 973.75 (×3), 1,947.50 y 30-oct-2026, todas del oráculo (j) |
| C1-10 31 entidades y vistas | alta | Cumple | §7 y §8 sin cambios de fondo. verify_c1 sin fallas |
| C1-11 Cobertura funcional | alta | Cumple | §11.6 regenerada. Las capacidades nuevas (`getProjectionSettings`, candidatos de chat) existen |
| C1-12 Seguridad | alta | Cumple con observación | `cookieAuth`: HttpOnly, Secure, `SameSite=Strict`, 90 días renovable (K11). El login describe 5 intentos en ventana deslizante de 15 min por usuario e IP, con 429. Públicos solo login y salud. 0 campos sensibles. `LoginRequest.username` no declara el 3–64 de K6 que §1.11 enuncia (H8) |
| C1-13 Versionado | media | Cumple | 0 operationId duplicadas, 9 etiquetas, `info.version` 1.1.0, Cambios |
| C1-14 Idioma y ejemplos | media | Cumple con observación | 14 de 14 ejemplos válidos con jsonschema. Muestra de 20 en español. Pero 202 descripciones de propiedad son solo el identificador (`id`, `created_at`, `updated_at`, `notes`, `is_active`), heredado de v1.0 (H6) |
| C1-15 PROVISIONAL | alta | Cumple con observación | (h). Supuestos del reporte revisados: ninguno cambia cifras. El retiro de `createNotificationRecipient` lo exige C1-18 y el `± días` como restricción lo respalda la pregunta (1) del bloqueo `S1-restricciones-optimizador`. Dos marcas sin `x-provisional` (H7) |
| C1-16 `scope` (K2) | alta | **No cumple** | `openapi.yaml`: `updateBudget` (N-087, N-278) con `scope` enum `this_month`/`onward`, `default: onward`, sin marca, y la partición descrita. Pero `C1-datos.md` solo nombra los valores y el valor por omisión (líneas 28 y 857): no describe el efecto de cada alcance ni trae un ejemplo de cada uno (H1). El diccionario no registra los dos valores (H5) |
| C1-17 Borrador del ticket (K3) | alta | Cumple con observación | `saveTicketScanDraft` (PUT, N-075, nueva), `cookieAuth`, errores estándar, `x-tablas` con `attachment`. `TicketScan` con `has_draft`/`draft`. Columnas `ocr_draft_json` y `ocr_draft_saved_at` en la base y en el diccionario. «No crea precio» solo está en `openapi.yaml`, no en `C1-datos.md` (H3) |
| C1-18 Chat de Telegram (K4) | alta | Cumple con observación | `listTelegramChatCandidates` (N-299) y `confirmTelegramChatCandidate` (N-316) con sesión. Ninguna escritura recibe un chat libre. `NotificationRecipient` no expone `chat_id`. Respuestas sin token. `chat_id_masked` no está en el diccionario (H5) |
| C1-19 Captura atrasada (K9) | alta | **No cumple** | `createExpense` y `createWeightEntry`: `Idempotency-Key` `required: true`, 16–128, `captured_at` `date-time`, `x-tablas` con `expense`/`weight_entry`, 200 al reintento y 409 con otro cuerpo. Índices únicos `ux_expense_idempotency_key` y `ux_weight_entry_idempotency_key` en la base. Pero `C1-datos.md` no trae el ejemplo de reintento y no dice qué pasa al vencer las 24 h en gastos y peso (H2) |
| C1-20 Decimales de tasa (K12) | alta | Cumple | `rate_decimals` (entero) y `rate_pct` (texto, `x-calculado`) en `CreditTermRate`, `CreditPurchase` y `Promotion`; `annual_rate_decimals` en `Credit`. Fracción `rate_num`/`rate_den`. 17.27 → 2 y 25.8685 → 4 en las descripciones y en §4. Columnas en la base y en el diccionario |
| C1-21 Extractos (K32) | alta | Cumple | (f). 16 + `componentes-comunes.yaml` + `INDICE.md`, regenerables byte a byte |

## Decisiones §K y hallazgos heredados

| Punto | Estado | Evidencia |
|---|---|---|
| K1 lista de compras | OK | 5 operaciones de `shopping-list`, sin marca |
| K2 `scope` | Parcial | Ver C1-16 / H1 |
| K3 borrador | OK (obs. H3) | Ver C1-17 |
| K4 candidatos | OK (obs. H5) | Ver C1-18. `createNotificationRecipient` retirada también de `routes.go` |
| K5 `family` | OK | Enum de 3, incluyente, sin marca |
| K6 valores técnicos | OK (obs. H8) | `limit` máx. 200, `cursor` ≤ 512, `q` ≤ 200, clave 16–128, contraseña `minLength: 12`, PROFECO 1 GB, `SameSite=Strict`, 24 h en §1.8 y §1.11. Ninguno con `x-provisional` |
| K7 campos calculados | OK | Sin marca |
| K9 idempotencia | Parcial | Ver C1-19 / H2 |
| K11 sesión | OK | Login y `cookieAuth` |
| K12 tasas y «Z» | OK (obs. H4) | Ver C1-20 y C1-07 |
| K13 tope por persona | OK | `fortnight_payment_cap_cents` opcional, descripción neutral («El contrato no fija ninguna cifra»). El 97375 visible es el ejemplo genérico de `MoneyCents` |
| K16 restricciones | OK | `OptimizerSettings` sin valores por omisión («Lo captura Axel; no hay valor por defecto»). `missing_constraints` calculado. `ScenarioRequest` sobreescribe. `ScenarioResult.status` |
| K18 tarjeta y tienda | OK | `rate_base`, `vat_on_interest_bp`, `min_payment_pct_bp`, `min_payment_fixed_cents`, `annual_rate_decimals` y `missing_terms` en `Credit` (y en las altas y cambios, salvo el calculado) |
| K21 proyección | OK | `getProjectionSettings`/`updateProjectionSettings` (N-318, N-319). `repeat_cadence_enabled` false por omisión. `projection_mode` con `contracted_only` por omisión |
| K23 meta de kcal | OK | `goal` (lose/maintain/gain) en `HealthProfile` y en su cambio. `suggested_kcal_target` calculado, PROVISIONAL solo en la magnitud del ajuste (bloqueo abierto `S1-ajuste-objetivo-kcal`) |
| R2-06 «Z» | Cerrado con observación | §4 y `DateTime` lo documentan. Las salidas usan desplazamiento. Ver H4 |
| R3-01 semana de la lista | Cerrado | `period` (`ListWeek`, semana ISO) `required: true` y sin valor por omisión en las 5 operaciones de la lista, `listBaskets` y `getDashboardBestBasket` |
| R3-02 límites sin marca | Cerrado con observación | Sin marca. Falta el 3–64 del usuario en el esquema (H8) |

## Hallazgos

| Id | Sev. | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| H1 | alta | C1-16 (K2) | `C1-datos.md` no describe qué hace cada alcance sobre las vigencias (`valid_from_date`/`valid_to_date`) y no trae un ejemplo de cada uno. Solo dice «`scope` (`this_month` \| `onward`, por omisión `onward`)» en Cambios (línea 28) y en §14 (línea 857). La partición está solo en la descripción de `updateBudget` en `openapi.yaml`, y sin ejemplo | `docs/contratos/C1-datos.md` (fuente: `backend/internal/api/c1gen/c1_plantilla.md`) | Agregar en `c1_plantilla.md` (§7 renglón presupuesto, o una subsección de §6/§11) el efecto de `onward` y de `this_month` sobre los renglones de vigencia, con un ejemplo de cada uno: renglones antes y después y respuesta. Montos de la semilla o del oráculo, o marcador explícito. Regenerar con `build.py` y confirmar con `--check` |
| H2 | alta | C1-19 (K9) | `C1-datos.md` §1.8 documenta la clave, las 24 h y la purga, pero no trae el ejemplo de reintento. Tampoco dice qué pasa al vencer las 24 h en gastos y peso: misma clave → 200 con el mismo registro; otro cuerpo ya no se detecta. Eso solo está en las descripciones de `createExpense`/`createWeightEntry` | `docs/contratos/C1-datos.md` §1 punto 8 (fuente: `c1_plantilla.md`) | Agregar en §1.8 el ciclo completo de gastos y peso (antes y después de 24 h) y un ejemplo: primer `POST /api/v1/expenses` con `Idempotency-Key` → 201 con `id`; reintento con la misma clave y el mismo cuerpo → 200, mismo `id`, `Idempotent-Replayed: true`; misma clave con otro cuerpo dentro de 24 h → 409 `idempotency_key_reused`. Cifras del oráculo o marcador. Regenerar |
| H3 | media | C1-17 (K3) | La verificación de C1-17 pide confirmar en `C1-datos.md` que el borrador no inserta filas de precio. Solo `openapi.yaml` lo dice («No crea gasto, precio ni inventario»). `C1-datos.md` (líneas 27 y 856) no | `docs/contratos/C1-datos.md` (fuente: `c1_plantilla.md`) | Agregar en Cambios (Ticket, K3) y en §14 fila 4: «el borrador no crea gasto, precio ni inventario; solo `confirmTicketScan` guarda precios (decisiones E) y vacía el borrador». Regenerar |
| H4 | media | C1-07 (R2-06, K12) | 22 propiedades `DateTime` llevan una descripción hermana del `$ref` que termina en «(UTC)», tomada del comentario de columna de S1, y así contradicen la regla de salida con el desplazamiento de `America/Mexico_City`. Entre ellas: `Expense.captured_at`, `ExpenseCreate.captured_at`, `WeightEntry.captured_at` y `Attachment.ocr_draft_saved_at` (nuevas en v1.1), `DeviceSession.expires_at`, `UserAccount.last_login_at`, `ShoppingListItem.bought_at`, `GeneratedReport.sent_at` y `NotificationRecipient.verified_at`. Un implementador puede leerlo como «la respuesta va en UTC» | `docs/contratos/openapi.yaml` (fuente: `backend/internal/api/c1gen/gen2_schemas.py`) | En `gen2_schemas.py`, cuando la propiedad es `DateTime`, quitar «(UTC)» del comentario de columna, como ya se quita el prefijo «Nueva en 1.2 (Kn).», o cambiarlo por «(la base guarda UTC; la API responde con el desplazamiento de America/Mexico_City)». Regenerar y confirmar que el conteo queda en 0 |
| H5 | media | C1-16, C1-18 | Inconsistencia C1 ↔ diccionario. C1-16 pide que el diccionario registre los identificadores de los dos alcances (`this_month`, `onward`): `grep` da 0 en `diccionario-datos.md`. `chat_id_masked` (C1-18: «los campos los fija el diccionario») tampoco está (0). El reporte C1-v1.1 no los pide en sus solicitudes (la 2 habla solo de `ocr_draft_json` e `idempotency_key`) | `docs/contratos/diccionario-datos.md` (S1); `docs/estado/p1/C1-v1.1.md` §Solicitudes | Que C1 agregue la solicitud a S1 y que S1 registre en `diccionario-datos.md` los valores de `scope` (`this_month` = «solo este mes», `onward` = «de este mes en adelante», por omisión) y el calculado `chat_id_masked` («•••» y los últimos 4 caracteres) |
| H6 | baja | C1-14 | 202 descripciones de propiedad son solo el identificador: `id` 43, `created_at` 43, `updated_at` 43, `notes` 41 e `is_active` 32. Las dos últimas son palabras en inglés. Viene de v1.0 (41 `notes` en HEAD) | `docs/contratos/openapi.yaml` (fuente: `gen2_schemas.py`) | Texto de respaldo en español en el generador: «Identificador», «Alta (instante)», «Último cambio (instante)», «Notas», «Activo» |
| H7 | baja | C1-15 | Dos marcas PROVISIONAL viven solo en el texto de la descripción y no llevan `x-provisional: true`, a diferencia de las otras cuatro: `DashboardWeekMeals` («casillas sin armar», máx. 21) y `NotificationUpdate` (claves de `params_json`). Quien filtre por la extensión no las ve. La matriz de `C1-datos.md` sí las marca | `docs/contratos/openapi.yaml` líneas 11924 y 17167 (fuente: `gen2_schemas.py`) | Agregar `x-provisional: true` a esos dos esquemas, o a sus propiedades, citando `C1-valores-k31` |
| H8 | baja | C1-12 / K6 | `C1-datos.md` §1.11 enuncia «usuario de 3 a 64 caracteres» (K6) y la base tiene el CHECK 3–64, pero `LoginRequest.username` no declara `minLength: 3` / `maxLength: 64` (`cursor`, `q` y la clave sí llevan su límite) | `docs/contratos/openapi.yaml` `LoginRequest` (fuente: `gen2_schemas.py`) | Agregar `minLength: 3` y `maxLength: 64` a `username` en las entradas, o aclarar en §1.11 que el límite lo aplica S1 al crear cuentas |

## Observaciones fuera del alcance de C1 (sin hallazgo)

- S1: `credit.vat_on_interest_bp` NULL significa «ya incluye IVA». Por eso no se distingue «no capturado» de «incluye IVA», aunque la nota 16 del diccionario lista el IVA entre los datos de K18 que muestran «falta capturar». C1 sigue la semántica de la columna (supuesto 7 del reporte). Si se quiere el aviso también para el IVA, S1 necesita otra columna.
- `openapi.yaml` ya pesa 629 771 bytes (el protocolo §4 dice 597 KB). La regla de leerlo solo por extractos no cambia.

## Para la ronda 2

Corregir H1 y H2 (alta) y H3, H4 y H5 (media) en la fuente del generador (`c1_plantilla.md` y `gen2_schemas.py`), más la solicitud a S1 de H5. Luego regenerar con `build.py` y reejecutar (a) a (j). H6 a H8 son opcionales.
