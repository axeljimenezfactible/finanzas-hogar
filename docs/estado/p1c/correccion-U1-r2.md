# Corrección U1 · ronda 2 (flujo P1c, 11-oct-2026)

Agente: **correccion-U1-r2**. Encargo: H-r2-01 (media) de `docs/estado/p1c/REVISION-U1-r2.md`. Alcance: `docs/design/U1-arquitectura.md` (por parches), `docs/design/tools/verificar_u1.py` y este reporte.

## Hecho
- **H-r2-01 (media) cerrado** con la «otra salida» que da el revisor, porque `docs/estado/p1c/U1-1.2.md` queda fuera de mi alcance. La sección «Para C1» completa y corregida va al final de este reporte, como «Para C1 (corregida r2)». La línea PENDIENTE de `verificar_u1.py` apunta a ella, y el apuntador es dinámico: si el integrador aplica la sección en `U1-1.2.md` (solicitud 1), el verificador pasa a apuntar ahí. Una comprobación nueva del verificador rechaza cualquier «Para C1» con la viñeta vieja («hoy no tienen valor sembrado»).
- **Desfase encontrado al revisar la coherencia entre artefactos.** `backend/db/migrations/00003_puerta1b.sql` cambió a las 15:23 con la corrección S1-r1 (H-01 de `docs/estado/p1c/REVISION-S1-r1.md`). Eso fue después de la corrección U1 r1 (15:17) y de la viñeta de `correccion-U1-r1.md` (15:18). Los umbrales del presupuesto y los días antes del pago próximo salieron de `notification.params_json` y ahora son claves de `app_setting`: `alerts.budget_warning_pct` (80), `alerts.budget_limit_pct` (100) y `alerts.payment_due_days_before` ([1, 0]), con decisión F y K44. Además, el bloqueo `S1-parametros-avisos` se redujo a `send_time` y `advance_days`. La viñeta de r1 («con la opción A, la que S1 deja aplicada, viven en `params_json`») ya no coincidía con el esquema. La revisión r2 (15:24) cita el encabezado anterior de 00003 (líneas 3 a 14). Por eso la viñeta corregida r2 y U1 nombran ahora las claves `alerts.*`. No se reabre ninguna decisión: el cambio viene de 00003 y de K44.
- **`docs/design/U1-arquitectura.md`, 9 parches** con un script que exige una sola aparición de cada texto:
  - fila 2 de los cambios 1.2 y fila 9 nueva («Corrección r2»);
  - §2.3: la fila de umbrales nombra `alerts.budget_warning_pct` y `alerts.budget_limit_pct`; la de pago próximo nombra `alerts.payment_due_days_before`, con [1, 0] y la necesidad N-323;
  - N-298 y N-300: la anticipación del pago próximo es `alerts.payment_due_days_before` y se edita en N-323. N-300 se queda con `advance_days` y `send_time`;
  - N-323 pasa a «Umbrales del presupuesto y anticipación del pago próximo», con las tres claves, sus valores iniciales y «Restablecer»;
  - §7, fila «Parámetros operativos configurables».

  No cambian la barra, las pantallas, las 14 secciones, §4 (3 toques), la matriz §3.29, los conteos (211 N) ni ninguna cifra.
- **`docs/design/tools/verificar_u1.py`**:
  - docstring de la corrección r2;
  - la comprobación H-01 pasa a H-01/H-r2-01. Lee las líneas `@param` `alerts.*` de 00003 y el origen [1, 0] de 00001. Exige las claves en N-298, N-300, N-323 y las filas de §2.3, y rechaza un `days_before` suelto en necesidades y parámetros;
  - comprobación nueva H-r2-01 sobre la sección «Para C1» que lee C1 v1.2;
  - la línea PENDIENTE apunta a esa sección.

## Archivos
- `docs/design/U1-arquitectura.md` (modificado por parches; diff de 28 líneas).
- `docs/design/tools/verificar_u1.py` (modificado).
- `docs/estado/p1c/correccion-U1-r2.md` (este reporte; incluye «Para C1 (corregida r2)»).
- Temporales en `datos/tmp/correccion-u1-r2/` (copias de respaldo, scripts de parche, raíz de prueba): borrados al terminar.

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py` antes de corregir: sale con 0, 71 líneas, 62 OK. La línea PENDIENTE mandaba a «Para C1» de `docs/estado/p1c/U1-1.2.md`, que tiene la viñeta vieja.
- El mismo comando después de corregir: sale con 0, 72 líneas, 63 OK, «RESULTADO: TODO OK». Líneas nuevas:
  - `OK H-01/H-r2-01 (P1c) pago próximo con alerts.payment_due_days_before [1, 0] y umbrales alerts.budget_warning_pct 80 y alerts.budget_limit_pct 100 (00003 tras S1-r1; origen 00001, decisión F) en N-298, N-300, N-323 y la tabla de §2.3, sin days_before suelto; …; mal: []`
  - `OK H-r2-01 «Para C1» coherente con U1 1.2 y 00003 (8 N nuevas; …); fuente: «Para C1 (corregida r2)» de docs/estado/p1c/correccion-U1-r2.md`
  - `PENDIENTE C1 v1.2: 8 N-xxx nuevas … (esperado; van en «Para C1 (corregida r2)» de docs/estado/p1c/correccion-U1-r2.md): ['N-320', …, 'N-327']`, las mismas 8 de antes.
- **Prueba A (negativa)**, en una raíz temporal con la U1 anterior y sin este reporte: sale con 1 y «2 falla(s)».
  - H-01/H-r2-01 marca N-298, N-300 y N-323, la fila de pago próximo con `days_before` suelto y las filas de pago próximo y umbrales de §2.3.
  - H-r2-01 da «fuente: ninguna», porque la viñeta de `U1-1.2.md` es la vieja.
  - La línea PENDIENTE dice «sin fuente coherente: ver H-r2-01».
- **Prueba B (el integrador aplica la solicitud 1)**: la U1 corregida y una copia de `U1-1.2.md` con la sección r2 sustituida. Sale con 0 y «TODO OK», y la fuente y la línea PENDIENTE pasan a «Para C1» de `docs/estado/p1c/U1-1.2.md`.
- `diff` entre «Para C1» de `U1-1.2.md` y «Para C1 (corregida r2)»: solo cambian 2 viñetas, la de N-323 y la de N-297, N-298 y N-300.
- `grep -c PROVISIONAL` en U1: 4 antes y 4 después, todas en filas de historial. La fila 9 no agrega ninguna marca.
- `days_before` suelto (con comilla invertida) en U1: solo en la fila 8 del historial (corrección r1).
- Cifras: todas salen de los archivos. [1, 0] de 00001 (decisión F) y de la línea `@param` de 00003; 80 y 100 de las líneas `@param` de 00003 (decisión F, K44). No hay valores nuevos ni topes.

## Supuestos
- Ninguno que cambie cifras, esquema, endpoints ni tokens. Que N-323 absorba la anticipación del pago próximo (antes en N-300) se sigue de dos cosas: la clave vive en `app_setting` desde 00003, y §2.3 tiene la regla «una necesidad por grupo». Es la rama «operación de settings» que la viñeta r1 ya preveía para N-323. La pantalla P-26.09 no cambia y el nombre de la operación lo decide C1.
- Alineé U1 con 00003 tal como quedó después de la corrección S1-r1. Si una revisión posterior de S1 lo revierte, el verificador lo detecta (lee las líneas `@param` de 00003 y marcaría FALLA).

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`, ya solo sobre `send_time` y `advance_days`. U1 y la viñeta r2 lo citan así.

## Pendientes
- Aplicar la sección «Para C1 (corregida r2)» en `docs/estado/p1c/U1-1.2.md` (solicitud 1). Hasta entonces, el verificador apunta a este reporte.
- C1 v1.2 aplica «Para C1». Mientras no lo haga, el verificador sigue marcando PENDIENTE N-320 a N-327, como se espera.

## Solicitudes fuera de alcance
1. **`docs/estado/p1c/U1-1.2.md`, sección «Para C1»**: sustituir la viñeta de **N-323** y la de **N-297**, **N-298** y **N-300** por las de «Para C1 (corregida r2)» de abajo. El `diff` muestra que son las únicas 2 líneas distintas. Reemplaza a la solicitud 1 de `docs/estado/p1c/correccion-U1-r1.md`, que quedó desfasada por la corrección S1-r1. Motivo: H-r2-01 de `REVISION-U1-r2.md` y coherencia con 00003, donde las claves `alerts.*` viven en `app_setting`.
2. **`docs/contratos/` (C1 v1.2, regenerado con `build.py`)**: reemplaza a la solicitud 2 de `correccion-U1-r1.md`.
   - Las descripciones de `params_json` en `Notification` y `NotificationUpdate` ya no deben nombrar `days_before` ni `thresholds_pct`, porque desde 00003 no van ahí. Sí deben nombrar `send_time` y `advance_days`, este último solo para precios viejos.
   - N-323 necesita una lectura y un cambio de settings para `alerts.budget_warning_pct`, `alerts.budget_limit_pct` y `alerts.payment_due_days_before`, con el valor vigente y el inicial (K44).
   - Motivo: coherencia entre U1 1.2, RN-29 1.3, 00003 y C1.

## Propuestas para Axel
- Ninguna nueva.

## Para C1 (corregida r2)
**Nuevas (8, todas de escritura; salen PENDIENTE en `verificar_u1.py` hasta C1 v1.2). Nombres de clave según RN-29 de `reglas-negocio.md` 1.3 (S1, en paralelo; confirmarlos al integrar):**
- **N-320** (E, P-26.01) Parámetros de acceso: leer (valor vigente y valor inicial) y cambiar o restablecer `auth.max_failed_attempts` (5), `auth.lockout_window_minutes` (15), `session.max_age_days` (90) y `auth.password_min_length` (12) (K6, K11, K44). Hoy no hay operación de settings para estas claves: hace falta una lectura y un cambio (p. ej. `GET`/`PATCH /api/v1/settings/security`; el nombre lo decide C1). `login`, `changeOwnPassword` y la sesión usan el valor vigente.
- **N-321** (E, P-26.05) Tiempo máximo del comparador: `comparator.timeout_seconds` (5; K43). Cabe en `getOptimizerSettings`/`updateOptimizerSettings` (`/api/v1/settings/optimizer`) si se suma el campo con su valor inicial. `compareScenarios` responde 503 `unavailable` al pasar el valor vigente.
- **N-322** (E, P-26.06) Parámetros de precios: `price.stale_after_days` (30; K22), `price.zero_weight_after_days` (90; K40) y `price.profeco_import_max_gb` (1; K6). Cabe en `getPriceSettings`/`updatePriceSettings` (`/api/v1/settings/prices`); `importProfeco` usa el máximo vigente.
- **N-323** (E, P-26.09) Parámetros de los avisos: leer (valor vigente y valor inicial) y cambiar o restablecer `alerts.budget_warning_pct` (80) y `alerts.budget_limit_pct` (100), umbrales del aviso de presupuesto por categoría y total, y `alerts.payment_due_days_before` ([1, 0]; lista JSON de enteros ≥ 0, días antes del pago), anticipación del aviso de pago próximo (decisión F, K44). Desde el esquema 1.3 (`00003_puerta1b.sql`, corrección S1-r1) son claves de `app_setting` y ya no van en `notification.params_json`: hace falta una lectura y un cambio de settings (p. ej. `GET`/`PATCH /api/v1/settings/alerts`; el nombre lo decide C1). Los avisos `budget_threshold` y `payment_due` usan el valor vigente.
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
- **N-297**, **N-298** y **N-300**: `send_time` (`HH:MM`, America/Mexico_City) y `advance_days` (entero ≥ 0) en `params_json`, claves FIRMES (K43; quitar la marca de `NotificationUpdate`). `advance_days` es solo la anticipación del aviso de precios viejos (`stale_prices`), como define `00003_puerta1b.sql`. `send_time` y `advance_days` no tienen valor sembrado (RN-29 1.3; los propone S7, K31): la respuesta dice que no hay valor inicial y «Restablecer» no aparece. Si al final viven en `params_json` o en `app_setting` lo decide el bloqueo abierto `S1-parametros-avisos`; mientras tanto van en `params_json`, como dice K43. La anticipación del aviso de pago próximo **no** es `advance_days`: es `alerts.payment_due_days_before` de `app_setting` (lista de días antes del pago; valor inicial [1, 0], decisión F). Junto con los umbrales `alerts.budget_warning_pct` (80) y `alerts.budget_limit_pct` (100), se lee, se cambia y se restablece por N-323, con su valor vigente y su valor inicial (K44). Desde el esquema 1.3 (00003), `payment_due` ya no lleva `days_before` en `params_json` ni `budget_threshold` lleva `thresholds_pct`: las descripciones de `Notification` y `NotificationUpdate` dejan de nombrarlos. N-298 muestra en cada aviso la anticipación y los umbrales vigentes: los toma de la lectura de N-323 o C1 los repite en `listNotifications` (a criterio de C1).
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
- Acciones de herramienta: ~34 de 60.
- Leído: ~95 KB (protocolo 6, decisiones §K 15, REVISION-U1-r2 12, U1-1.2 por rangos ~4, correccion-U1-r1 por rangos ~9, verificador por rangos ~12, 00003 por grep ~7, bloqueo S1-parametros-avisos ~3, INTEGRACION §7 por grep ~2, U1 por rangos y grep ~12, salidas del verificador y diffs ~13).
- Escrito: ~38 KB (scripts de parche y de prueba ~16, cambio neto en U1 ~+3 y en el verificador ~+5, este reporte ~14). Las copias de respaldo en `datos/tmp/` (cp) se borraron.
