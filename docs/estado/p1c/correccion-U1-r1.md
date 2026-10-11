# Corrección U1 1.2 · ronda 1 (flujo P1c, 11-oct-2026)

Agente: **correccion-U1-r1**. Corrige el único hallazgo de severidad alta o media de `docs/estado/p1c/REVISION-U1-r1.md`: **H-01 (media)**. Los 6 hallazgos bajos (H-02 a H-07) no estaban en el encargo y no se tocaron. No se reabre ninguna decisión de §A a §K.

## Hecho
- **H-01 · anticipación del pago próximo** (`docs/design/U1-arquitectura.md`):
  - **N-298** (Alertas): el pago próximo nombra su clave del esquema, `days_before` de `payment_due` (lista de días antes del pago, valor inicial [1, 0], decisión F). `advance_days` (días enteros ≥ 0) queda **solo** para el aviso de precios viejos, como en `00003_puerta1b.sql` línea 46 y en la propuesta aprobada por K43 (`C1-valores-k31`, punto 3). `send_time` sigue igual. Dice que `days_before` y los umbrales están hoy en `notification.params_json` (00001) y que, si pasan a `app_setting`, lo decide el bloqueo abierto `S1-parametros-avisos` sin cambiar la pantalla. Es el mismo trato que el revisor aceptó para N-323, por eso no lleva marca PROVISIONAL.
  - **N-300** (Editar alertas y probar): edita `days_before` del pago próximo (valor inicial [1, 0], decisión F), `advance_days` del aviso de precios viejos y `send_time`. «Restablecer al valor inicial» devuelve el pago próximo a [1, 0].
  - **Tabla «Parámetros operativos» de §2.3**: la fila del pago próximo (antes línea 317) cita la clave `days_before` de `payment_due`, con valor inicial «[1, 0]: 1 día antes y el mismo día (decisión F; N-298)». La fila de `advance_days` (antes 318) dejó de decir «los demás avisos»: ahora es la anticipación del aviso de precios viejos antes del día de despensa, más la hora (`send_time`) de los avisos que la llevan. Sigue «Sin valor inicial decidido (K31)». La columna «Decisión» no cambió: sumar «F» es H-02 (baja), que no estaba en el encargo.
  - **Lista de cambios 1.2**: la fila 2 ya no dice que `advance_days` es la anticipación «de reportes y alertas». Se suma la fila 8, que registra esta corrección.
- **«Para C1»**: la viñeta de N-297, N-298 y N-300 está corregida más abajo. El archivo `docs/estado/p1c/U1-1.2.md` queda fuera de mi alcance, así que dejo la sección completa ya corregida en «Para C1 (corregida)» y la solicitud 1 para que el integrador la sustituya. La viñeta corregida pide que `listNotifications` y `updateNotification` devuelvan el valor inicial de `days_before` ([1, 0]) y de `thresholds_pct` ([80, 100]) para «Restablecer», según la opción con que se resuelva `S1-parametros-avisos`.
- **Verificador** (`docs/design/tools/verificar_u1.py`): nueva comprobación «H-01 (P1c r1)». Lee de 00001 el valor sembrado de `days_before` de `payment_due` y de 00003 la definición de `advance_days`. Exige que N-298, N-300 y la fila del pago próximo de §2.3 lleven `days_before` con ese valor inicial y «decisión F». Exige también que `advance_days` aparezca solo como anticipación del aviso de precios viejos. El docstring registra la corrección.

## Archivos
- `docs/design/U1-arquitectura.md` (por parche con `str.replace`, 6 sustituciones; cada texto viejo aparece una sola vez).
- `docs/design/tools/verificar_u1.py` (por parche, 2 sustituciones).
- `docs/estado/p1c/correccion-U1-r1.md` (este reporte).
- Temporales en `datos/tmp/correccion-u1-r1/` (scripts de parche, respaldos y salidas). Se borraron al terminar.

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py` → código de salida 0, **«RESULTADO: TODO OK»**, 71 líneas. Líneas clave:
  - `OK    K43 N-023, N-064, N-226, N-297, N-298 y N-300 FIRMES con sus valores (10 MB, 5 s, send_time, advance_days); mal: []`
  - `OK    K44 tabla «Parámetros operativos» (§2.3): 15 filas, cada parámetro en una de las 14 secciones con su valor inicial, decisión y necesidad; reglas comunes escritas; faltan: []`
  - `OK    H-01 (P1c r1) pago próximo con days_before [1, 0] (00001, decisión F) en N-298, N-300 y la tabla de §2.3; advance_days solo del aviso de precios viejos (00003); mal: []`
  - `OK    K44 una necesidad de escritura por grupo (N-320 … N-326) …; mal: []`
  - `PENDIENTE C1 v1.2: 8 N-xxx nuevas … ['N-320', …, 'N-327']`: es lo esperado, las mismas 8 de «Para C1».
- **Prueba en negativo**: corrí el verificador nuevo sobre la U1 anterior a la corrección (respaldo puesto en su lugar y luego restaurado; `cmp` confirmó que el archivo corregido volvió intacto). Resultado: `FALLA H-01 … mal: ['N-298', 'N-300', 'fila «pago próximo» de §2.3', 'fila advance_days de §2.3']` y «RESULTADO: 1 falla(s)». Al restaurar, «TODO OK».
- `diff` contra el respaldo: solo cambian las 5 líneas de la lista anterior y se suma la fila 8. Marcas PROVISIONAL en U1: 4, las mismas de antes, todas en las listas de cambios. En las líneas nuevas no hay topes («≥ 1», «365», «no menos»). Ninguna cifra en $ cambió (la comprobación de cifras contra `oraculo/model.json` y `data.json` da OK).

## Supuestos
- Ninguno que cambie cifras, esquema, endpoints ni tokens. Nombrar `days_before` en U1 no es un supuesto: la clave y su valor [1, 0] existen en 00001 (esquema FIRME desde la puerta 1) y vienen de la decisión F. Que se queden ahí o pasen a `app_setting` lo decide el bloqueo abierto `S1-parametros-avisos`, que U1 cita, y la pantalla no cambia.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (de S1). La viñeta corregida de «Para C1» depende de él solo en un punto: qué operación expone el valor inicial.

## Pendientes
- Sustituir en `docs/estado/p1c/U1-1.2.md` la sección «Para C1» por la de abajo (solicitud 1).
- C1 v1.2 aplica «Para C1». Mientras no lo haga, el verificador sigue marcando PENDIENTE N-320 a N-327.

## Solicitudes fuera de alcance
1. `docs/estado/p1c/U1-1.2.md`, sección «Para C1»: sustituir la viñeta «**N-297**, **N-298** y **N-300**» por la de la sección corregida de abajo; el resto queda igual. Motivo: H-01 (media) de `REVISION-U1-r1.md`. La viñeta actual dice que las tres «no tienen valor sembrado» y no pide exponer el valor inicial [1, 0] de `days_before` para «Restablecer», que K44 exige.
2. `docs/contratos/` (C1 v1.2, regenerado con `build.py`): las descripciones de `params_json` en `Notification` y `NotificationUpdate` dicen hoy «días antes 1 y 0» y «anticipación de precios viejos» sin nombrar claves. Conviene nombrar `days_before` (pago próximo), `thresholds_pct` (presupuesto), `advance_days` (precios viejos) y `send_time`, para que U1, RN-29 y C1 digan lo mismo. Motivo: coherencia que pide H-01. Solo cambian descripciones.

## Propuestas para Axel
- Ninguna nueva.

## Para C1 (corregida)
**Nuevas (8, todas de escritura; salen PENDIENTE en `verificar_u1.py` hasta C1 v1.2). Nombres de clave según RN-29 de `reglas-negocio.md` 1.3 (S1, en paralelo; confirmarlos al integrar):**
- **N-320** (E, P-26.01) Parámetros de acceso: leer (valor vigente y valor inicial) y cambiar o restablecer `auth.max_failed_attempts` (5), `auth.lockout_window_minutes` (15), `session.max_age_days` (90) y `auth.password_min_length` (12) (K6, K11, K44). Hoy no hay operación de settings para estas claves: hace falta una lectura y un cambio (p. ej. `GET`/`PATCH /api/v1/settings/security`; el nombre lo decide C1). `login`, `changeOwnPassword` y la sesión usan el valor vigente.
- **N-321** (E, P-26.05) Tiempo máximo del comparador: `comparator.timeout_seconds` (5; K43). Cabe en `getOptimizerSettings`/`updateOptimizerSettings` (`/api/v1/settings/optimizer`) si se suma el campo con su valor inicial. `compareScenarios` responde 503 `unavailable` al pasar el valor vigente.
- **N-322** (E, P-26.06) Parámetros de precios: `price.stale_after_days` (30; K22), `price.zero_weight_after_days` (90; K40) y `price.profeco_import_max_gb` (1; K6). Cabe en `getPriceSettings`/`updatePriceSettings` (`/api/v1/settings/prices`); `importProfeco` usa el máximo vigente.
- **N-323** (E, P-26.09) Umbrales del aviso de presupuesto (80 % y 100 %; K44). La operación depende del bloqueo abierto `S1-parametros-avisos`. Con la opción A, la que S1 deja aplicada: `listNotifications`/`updateNotification` con `thresholds_pct` de `budget_threshold` en `params_json`, con su valor inicial. Con la opción B: una operación de settings.
- **N-324** (E, P-26.09) Descartar chat detectado: `DELETE /api/v1/telegram-chat-candidates/{candidate_id}`, la que C1 v1.1 dejó condicionada en su solicitud 1. Borra el renglón de `telegram_chat_candidate` (S1 1.2) y no toca chats confirmados ni destinatarios.
- **N-325** (E, P-26.11) Parámetros de datos: `attachment.max_size_mb` (10; K43) y `idempotency.retention_hours` (24; K6). Necesita una operación de settings nueva (p. ej. `/api/v1/settings/data`); `uploadExpenseAttachment` y la retención de `Idempotency-Key` usan el valor vigente.
- **N-326** (E, P-26.13) Ajuste por objetivo de la casa: `health.goal_adjust_lose_kcal` (0) y `health.goal_adjust_gain_kcal` (0) (K36). Necesita una operación de settings nueva (p. ej. `/api/v1/settings/health`).
- **N-327** (E, P-21 y P-26.13) Ajuste por objetivo de una persona: `updateHealthProfile` (y `createHealthProfile`) aceptan `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal`: enteros ≥ 0 o `null` (`null` = usa el de la casa; «Restablecer» = `null`) (K36, 00003).

**Regla común de N-320 a N-323, N-325 y N-326 (K44):**
- La lectura devuelve, por parámetro, el valor vigente y el valor inicial sembrado, o dice que no lo hay.
- La escritura rechaza un valor vacío o de tipo no válido con error de validación y no lo sustituye por otro.
- No hay topes mínimos ni máximos hasta la puerta 2.
- Restablecer = volver al valor inicial (como campo o como acción, a criterio de C1).

**Cambiadas, con cambio de contrato:**
- **N-230** (`listHealthProfiles`) y **N-311**: por persona, `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` (`null` = el de la casa), más los valores de la casa para mostrar «El de la casa: n kcal». La meta sugerida usa el ajuste del objetivo elegido (RN-17 FIRME, K36).
- **N-297**, **N-298** y **N-300**: `send_time` (`HH:MM`, America/Mexico_City) y `advance_days` (entero ≥ 0) en `params_json`, claves FIRMES (K43; quitar la marca de `NotificationUpdate`). `advance_days` es solo la anticipación del aviso de precios viejos (`stale_prices`), como define `00003_puerta1b.sql`; `send_time` y `advance_days` hoy no tienen valor sembrado (RN-29 1.3; los propone S7, K31) y la respuesta dice que no hay valor inicial. La anticipación del aviso de pago próximo **no** es `advance_days`: es `days_before` de `payment_due` (lista de días antes del pago; valor inicial [1, 0] sembrado en 00001, decisión F). Por K44, `listNotifications` y la respuesta de `updateNotification` devuelven, por aviso, el valor vigente y el valor inicial de `days_before` de `payment_due` ([1, 0]) y de `thresholds_pct` de `budget_threshold` ([80, 100]) para «Restablecer al valor inicial» (el nombre del campo lo decide C1), y `NotificationUpdate` permite volver a ese valor. Sujeto a la opción que resuelva el bloqueo abierto `S1-parametros-avisos`: con la opción A (la que S1 deja aplicada) viven en `params_json` y los exponen estas dos operaciones; con la opción B pasan a claves de `app_setting` y los expone una operación de settings con la misma regla.
- **N-064** (`uploadExpenseAttachment`): 10 MB FIRME (K43; quitar la marca). El máximo es el vigente de N-325 y el error dice cuál es.
- **N-226** (`compareScenarios`): 5 s FIRME con 503 `unavailable` (K43; quitar la marca). El tiempo es el vigente de N-321.
- **N-289: quitar `updatePerson`** de sus operaciones (hoy `x-necesidades` de `updatePerson` lista N-310 y N-289). El tope por quincena por persona es de N-310 (K13). N-289 se queda con `createCredit`, `updateCredit`, tasas, promociones y `updateOptimizerSettings` (B6 y B7).

**Cambiadas sin cambio de contrato (solo descripción):**
- **N-023** (`DashboardWeekMeals`): «casillas sin armar» FIRME (K43); quitar la marca.
- **N-120** (`getMealPlanWeek`): sin `period`, la semana en curso y, los domingos, la siguiente: FIRME (K38); quitar la marca. N-150 y N-160 siguen con `week` obligatorio (R3-01).
- **N-299** (`listTelegramChatCandidates`): cada candidato trae el id que usa N-324.

**Adenda K18 (de `docs/estado/p1/correccion-U1-r2.md`, fusionada aquí; ya está en C1 v1.1 y `verificar_u1.py` da OK):**
- **N-200, N-201** (`getCredit`) y **N-283** (`listCredits`): en `bank_card` y `store_credit`, `Credit` expone `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents` y `annual_rate_decimals`, junto a `annual_rate_num`/`annual_rate_den` y `vat_on_interest_bp`. Vacío = «falta capturar», y la respuesta dice cuáles faltan (K18, K12, RN-25, RN-26).
- **N-289** (`createCredit`, `updateCredit`): `CreditCreate` y `CreditUpdate` aceptan esos cuatro campos, opcionales y sin valor por defecto. La base, el porcentaje y el mínimo fijo solo se aceptan en esos dos tipos (CHECK de 00002). La tasa se guarda con los decimales tecleados.
- **N-201** y **N-205** (`getCardStatus`): si falta alguno de los cinco datos, intereses y pago mínimo llegan sin valor (nunca en cero) y con «falta capturar» (RN-26). En esos tipos, N-201 lleva también el pago mínimo. Regla general (§2.4): cualquier respuesta con intereses o pago mínimo de esa línea sigue la misma regla.
- **Si incluye IVA**: hoy `vat_on_interest_bp` NULL significa «ya trae IVA», así que la API no distingue «falta capturar». Lo resuelve S1 (INTEGRACION §7, punto 2), no U1.

## Consumo aproximado
- Acciones de herramienta: 23.
- Leído: ~100 KB (protocolo 6, decisiones §K 15, REVISION-U1-r1 11, bloqueo S1-parametros-avisos 3, C1-valores-k31 y RN-29 por grep 3, U1 por rangos ~14, U1-1.2 «Para C1» ~9, verificador por rangos ~12, extracto ajustes-d por rangos ~6, 00001 y 00003 por grep ~4, salidas del verificador ~6, diff ~4).
- Escrito: ~32 KB (scripts de parche ~12, «Para C1» temporal 7, cambio neto en U1 +2.4 y en el verificador +2.1, este reporte ~12).
