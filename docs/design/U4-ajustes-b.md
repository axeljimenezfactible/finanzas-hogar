# U4 · Ajustes estilo Notion — parte b (secciones 8 a 14)

Versión 1.5 (11-oct-2026) · Flujo P2c, agente `correccion-U4-r1` (hallazgos M-1, M-2, M-3 y M-5 de `docs/estado/p2c/REVISION-U4-r1.md`): P-26.08 lee y escribe día y hora con `getReportSettings` y `updateReportSettings` de C1 v1.3 y cita el bloqueo abierto `C1-hora-reportes-job`; P-26.12 escribe `scope` (`onward` o `all`) de C1 v1.3, corrige lo que hace C1 v1.2, marca PROVISIONAL la preselección del diálogo (bloqueo nuevo `U4-reparto-preseleccion`) y cita `C1-reparto-alcance`; P-26.13 aplica Q2-16 A, Q2-17 C y Q2-18 A (`S2-metas` RESUELTO) y cita el bloqueo abierto `S2-metas-vencida`. Versión 1.4 (11-oct-2026) · Flujo P2c, agente `U4-puerta2` (decisiones de la puerta 2 en `docs/estado/puerta2/RESPUESTAS.md`: Q2-11 B (K47), P-26.08 y P-26.09 leen y escriben las claves `reports.*` y `alerts.*` de `app_setting` de la migración 00004, con una sola variante; Q2-72 y Q2-73 (K50) dan sus valores iniciales; Q2-66 B (en un mes corto sale el último día); Q2-68 C (P-26.12 pregunta el alcance); Q2-70 (a) (P-26.14: días restantes y consumo por semana); Q2-09 A (R5-1: P-26.13 y P-26.14 siguen la regla general de la parte a §0.2); Q2-76 (K50) `ocr.timeout_seconds` en P-26.11; Q2-82 (K51) interruptor «Atajos de una tecla» en P-26.10; R4-3 (citas en la misma línea)). Versión 1.3 (11-oct-2026) · Flujo WF-1, agente `correccion-U4-r4` (hallazgo R4-1 de `docs/estado/wf1/REVISION-U4-r4.md`: en P-26.13, el interruptor «Activa» de una meta no guarda al tocar: muestra antes la vista previa (3) del total reservado para metas, como la excepción de la parte a §0.2, y `Ctrl+/` la muestra sumada sobre las metas seleccionadas). Versión 1.2 (11-oct-2026) · Flujo WF-1, agente `correccion-U4-r2` (hallazgos R2-1 punto 1 y R2-3 de `docs/estado/wf1/REVISION-U4-r2.md`: en P-26.13, vista previa del IMC y de la meta sugerida con cualquier dato del perfil, y metas que citan `S2-metas` para los meses que faltan, el formato del avance y los estados). Versión 1.1 (11-oct-2026) · Flujo WF-1, agente `correccion-U4-r1` (hallazgos H-1, H-2 y H-3 de `docs/estado/wf1/REVISION-U4-r1.md`: lo común pasa a la parte a §0, vista previa en P-26.13 y P-26.14, una sola regla de reparto en P-26.12). Versión 1.0 del mismo día, agente `U4b-ajustes`. El integrador une esta parte con `U4-ajustes-a.md` (secciones 1 a 7) en `docs/design/U4-ajustes.md`. Fuentes: `docs/plan/decisiones.md` §A–§K, `docs/criterios/U4.md` 1.2, `docs/design/U1-arquitectura.md` 1.2 (§2.3, §2.4, §2.5, §2.6, §3.28, §5), extractos de C1 v1.2 (`docs/contratos/extractos/ajustes-a.yaml` a `ajustes-e.yaml`, `reportes.yaml`, `salud.yaml`, `presupuesto-a.yaml`, `presupuesto-b.yaml`, `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`, `inicio.yaml`), esquema 1.4 (`app_setting` e inventario `@param` de 00003 y 00004) y `docs/design/DESIGN.md` §2 y §4.6 (la §7 de componentes aún no existe; aquí solo se citan nombres de token ya publicados).

## Convenciones de la parte b (formato; lo común rige desde la parte a §0)

- **Encabezado de sección**: `## P-26.xx Nombre`. Conteo: `grep -c '^## P-26\.' docs/design/U4-ajustes-b.md` = 7 (14 al unir con la parte a).
- **Tabla de campos** de cada sección: Campo · Tipo · Validación y mensaje · Valor por defecto · Lectura · Escritura · Necesidad (`N-xxx` de U1 §3.28). «Valor por defecto» lleva solo lo decidido en §A–§K (con su cita) o, si no hay decisión, «Sin valor: lo captura Axel»; lo que no se edita dice «Calculado» o «Fijo». Las operaciones son `operationId` de C1 v1.2; «—» en Escritura = no se edita aquí.
- **Parámetros operativos (K44, U4-12)**: tabla Parámetro · Clave · Valor actual e inicial · Validación y mensaje · Restablecer · Decisión · Lectura y escritura. Reglas de U1 §2.3: el valor vigente y «Valor inicial: …» van en la misma fila; «Restablecer al valor inicial» aparece solo cuando `is_initial` es falso y manda `reset_parameters: [clave]`; un valor vacío o de otro tipo no se guarda (422 `validation_failed`) y el campo vuelve a su valor anterior; con los rangos aprobados en la puerta 2 (parte a §0.5; K50): día de la semana de 1 a 7, día del mes de 1 a 31 y horas «HH:MM» de 00:00 a 23:59; los demás, entero ≥ 0.
- **Lo común a las 14 secciones** está en un solo bloque, `U4-ajustes-a.md` §0, y esta parte no lo repite ni lo cambia: edición en el lugar, guardado y vista previa en línea (§0.2), estados vacío, carga, guardando, guardado, restablecido, error y sin conexión con un solo token y un solo texto cada uno (§0.3), diálogo de acción crítica (§0.4), parámetros operativos (§0.5), catálogo de mensajes V-xx (§0.6) y táctil, teclado y foco visible (§0.7). Lo reversible se guarda sin preguntar; las acciones de U4-08 de esta parte (restaurar un respaldo, activar una IA, descartar un chat detectado, que borra) usan el diálogo de §0.4.
- **Atajos**: solo los de U1 §5.3, sin atajos nuevos (0 duplicados): `↑`/`↓` y `Enter` (filas), `Esc` (cancelar edición o cerrar hoja), `Ctrl+Enter` (guardar hoja o formulario), `Alt+N` (nuevo registro en la base de la página), `Ctrl+/` (editar a la vez filas seleccionadas), `Ctrl+K` (buscar una sección: «Notificaciones», «Respaldos»…), `Ctrl+Shift+L` (tema), `?` (lista). Menú «/» de U1 §5.2: `/reporte` (P-26.08), `/respaldo` (P-26.11), `/tema` (P-26.10) e `/ir`.
- **Puntos de quiebre**: móvil 412 px = lista de secciones y cada sección como página; las tablas se ven como lista de tarjetas con sus propiedades apiladas y cada fila abre una hoja inferior con botón primario de `size.control.primary-sheet`. Escritorio 1440 px = columna interna con las 14 secciones y el contenido a la derecha (U1 §2.3).
- **Formatos** (decisión B): dinero `$1,234.56` con cifras tabulares; fechas `10-oct-2026` (con día de la semana donde importa: `Dom 18 oct 2026`); en el envío, ISO. Las cifras de ejemplo vienen del oráculo, de los criterios o son marcadores (`$T`, `n`, `hh:mm`).

## P-26.08 Reportes y frecuencias

**Propósito.** Programar los tres reportes automáticos de la decisión F (resumen diario en la noche, resumen semanal el domingo, reporte mensual en PDF el día 1), elegir quién los recibe y lanzar a mano lo que va bajo demanda: el reporte mensual ahora, el estado de cuenta de Karen y la vista web mensual (N-296, N-297).

**Vista de base de datos.** Tabla (única vista, U1 §2.3): una fila por reporte programado (`daily_summary`, `weekly_summary`, `monthly_report` de `listNotifications`), con su tarea programada de `listJobs`. Sin tablero, calendario ni galería: son tres filas fijas (no se crean ni se borran reportes). Debajo, el bloque «Bajo demanda».

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Reporte | Texto, solo lectura | — | Resumen diario · Resumen semanal · Reporte mensual en PDF (fijo, sembrado en 00001) | `listNotifications` (`name`) | — | N-296 |
| Activo | Interruptor | — | Encendido (decisión F: son los reportes por defecto; `is_enabled` = 1 al sembrar) | `listNotifications` (`is_enabled`) | `updateNotification` (`is_enabled`) | N-297 |
| Frecuencia | Texto, solo lectura | — | Diario · semanal · mensual (fijo por tipo, decisión F) | `listNotifications` (`code`) | — | N-296 |
| Día de la semana (solo semanal) | Selector de lunes a domingo (1 a 7, ISO 8601) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Domingo (7), decisión F; clave `reports.weekly_weekday` (Q2-11 B, K47) | `getReportSettings` (`reports.weekly_weekday`) | `updateReportSettings` (ver «Escritura de día y hora») | N-297 |
| Día del mes (solo mensual) | Número entero | Entero de 1 a 31 · «Escribe un día del mes, del 1 al 31»; con 29, 30 o 31, en los meses que no tienen ese día sale el último día del mes (Q2-66 B; bloqueo `docs/estado/bloqueos/U4-dia-mes-reporte.md` RESUELTO) | 1, decisión F; clave `reports.monthly_day_of_month` (Q2-11 B, K47) | `getReportSettings` (`reports.monthly_day_of_month`) | `updateReportSettings` (ver «Escritura de día y hora») | N-297 |
| Hora | Hora `HH:MM` de 24 h, hora del centro de México (America/Mexico_City) | De 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Diario 21:30, semanal 19:00 y mensual 08:00 (Q2-72, K50); claves `reports.daily_send_time`, `reports.weekly_send_time` y `reports.monthly_send_time` (Q2-11 B, K47). Debajo del campo, la pista de la decisión F: «Noche», «Domingo» o «Día 1» | `getReportSettings` (la clave de la fila) | `updateReportSettings` (ver «Escritura de día y hora») | N-297 |
| Destinatarios | Selección múltiple: personas con chat confirmado (Axel, Ximena) | Solo chats confirmados en P-26.09; Karen no aparece (sin cuenta, decisión A) · «Primero confirma el chat de esta persona en Notificaciones» | Sin valor: lo captura Axel (decisión A: destinatarios configurables por aviso) | `listNotifications` (`recipient_ids`) y `listNotificationRecipients` | `setNotificationRecipients` | N-297 |
| Próximo envío | Fecha y hora, calculada | — | Calculado por la tarea programada | `listJobs` (`next_run_at`) | — | N-296 |
| Último resultado | Estado y fecha, calculado | — | Calculado («Correcto», «Con error: …», «Omitido») | `listJobs` (`last_status`, `last_run_at`, `last_error`) | — | N-296 |

Reglas fijas que se leen en la sección (no son campos): canal Telegram en texto plano, sin formato (decisión F); el reporte mensual se envía por Telegram **y** se guarda en el Centro de reportes (decisión F).

**Escritura de día y hora (Q2-11 B, K47; bloqueo `docs/estado/bloqueos/S1-parametros-avisos.md` RESUELTO).** Día y hora de cada reporte son claves de `app_setting` del hogar (migración 00004; ya no están en `notification.params_json`, que queda `{}` en los tres reportes). Se leen con `getReportSettings` (`parameters`: `value`, `initial_value`, `is_initial`) y se cambian o restablecen por clave con `updateReportSettings` (`parameters` y `reset_parameters`; parte a §0.5), la sección «Reportes y frecuencias» que C1 v1.3 agregó para las 5 claves `reports.*` (`/api/v1/settings/reports`, N-296 y N-297). Las claves `alerts.*` siguen en `getAlertSettings` y `updateAlertSettings` (P-26.09). `updateNotification` solo cambia «Activo». La tarea programada de cada reporte (`listJobs`) tiene además su `schedule`, que C1 deja cambiar con `updateJob`. Qué manda, las claves o el `schedule` de la tarea, y por lo tanto si «Próximo envío» (`next_run_at`) se recalcula al guardar el día o la hora, es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/C1-hora-reportes-job.md`; en uso, `updateJob` acepta `schedule` en todas las tareas). Esta sección no usa `updateJob`.

**Parámetros operativos de la sección (K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Hora del resumen diario | `reports.daily_send_time` | Actual: `value` · Inicial: 21:30 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 21:30 | Q2-72 (K50); lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Día de la semana del resumen semanal | `reports.weekly_weekday` | Actual: `value` · Inicial: 7 (domingo) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Vuelve a 7 | Decisión F; lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Hora del resumen semanal | `reports.weekly_send_time` | Actual: `value` · Inicial: 19:00 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 19:00 | Q2-72 (K50); lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Día del mes del reporte mensual | `reports.monthly_day_of_month` | Actual: `value` · Inicial: 1 | Entero de 1 a 31; en un mes más corto sale el último día (Q2-66 B) · «Escribe un día del mes, del 1 al 31» | Vuelve a 1 | Decisión F, Q2-66 B; lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |
| Hora del reporte mensual | `reports.monthly_send_time` | Actual: `value` · Inicial: 08:00 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 08:00 | Q2-72 (K50); lugar Q2-11 B (K47) | `getReportSettings` · `updateReportSettings` |

**Bajo demanda.**
- «Generar reporte mensual ahora» (acción primaria; también `/reporte`): `generateMonthlyReport` con el mes (selector con los meses de `listMonthlyReportMonths`) y la casilla «Enviarlo también por Telegram a los destinatarios del mensual» (`send_telegram`, sin preselección: la persona decide). Respuesta aceptada: «Reporte en preparación · Ver en el Centro de reportes» (P-23).
- «Enviar ahora a…» en cada fila: `sendTestNotification` con la persona (Axel o Ximena); sirve de disparo manual del resumen.
- «Generar estado de cuenta de Karen» (decisión A: bajo demanda; Axel lo reenvía a mano): `generateExternalStatement` con Karen fija y el periodo (inicio y fin; «La fecha final no puede quedar antes de la inicial»). Al terminar abre P-25.
- «Abrir la vista web mensual»: enlace a P-24 con el último mes de `listMonthlyReportMonths`. Siempre disponible (decisión F); no se configura.

**Vista previa del efecto.** Al cambiar día u hora, debajo de la fila: «Próximo envío: Dom dd mmm aaaa, hh:mm», calculado en el teléfono o navegador con el día y la hora tecleados antes de guardar; al guardar, lo reemplaza `next_run_at` de `listJobs` (**PROVISIONAL**, bloqueo abierto `C1-hora-reportes-job`: ver «Escritura de día y hora»). Con día del mes 29, 30 o 31: «En los meses que no tienen el día {d}, sale el último día del mes (por ejemplo, el 28-feb-2027)» (Q2-66 B). Con el interruptor apagado: «No se enviará mientras esté apagado». Sin destinatarios: «Sin destinatarios: no llegará a nadie» (insignia del rol advertencia de DESIGN §2.9 con icono y texto).

**Estados.** Vacío: no aplica (las tres filas vienen sembradas); si `listNotifications` no trae alguna, «Falta el reporte en la base · Avisar a Axel» (error de instalación). Carga: esqueleto de tres filas. Error al cargar: «No se pudo cargar la configuración de reportes · Reintentar». Error al guardar: el de la parte a §0.3 («No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar»; el campo vuelve al valor anterior). Guardado: automático, «Guardado».

**Edición en el lugar.** Activo, día, hora y destinatarios se editan en la celda (o en la hoja en móvil). No se editan: nombre, frecuencia, próximo envío y último resultado (calculados o fijos).

**Atajos.** `↑`/`↓` entre filas, `Enter` abre la celda o la hoja, `Esc` cancela, `Ctrl+Enter` guarda la hoja; `/reporte` genera el mensual.

**Móvil y escritorio.** 412 px: tres tarjetas (nombre, activo, «Próximo envío») que abren una hoja con día, hora y destinatarios; el bloque «Bajo demanda» como botones de ancho completo. 1440 px: tabla con columnas Reporte · Activo · Día · Hora · Destinatarios · Próximo envío · Último resultado, y el bloque «Bajo demanda» debajo.

## P-26.09 Notificaciones

**Propósito.** Configurar las 4 alertas de la decisión F (pago próximo, presupuesto al 80 % y al 100 %, precios viejos antes del día de despensa, plan de comidas sin armar), quién recibe cada una, sus umbrales y anticipación (K44), y registrar los chats de Telegram de Axel y Ximena en el bot MultiX: el bot detecta el chat cuando la persona le escribe y ella lo confirma aquí o lo descarta (K4). Necesidades N-298, N-299, N-300, N-316, N-323 y N-324.

**Vista de base de datos.** Tres bloques en tabla (U1 §2.3: tabla): «Alertas» (4 filas fijas de `listNotifications`: `payment_due`, `budget_threshold`, `stale_prices`, `meal_plan_missing`), «Chats confirmados» (una fila por persona con cuenta) y «Chats detectados por confirmar» (lista que se vacía al confirmar o descartar). Los avisos ya enviados no viven aquí: están en P-23 (tabla, tablero por tipo y calendario por fecha de envío, U1 §2.5).

**Alertas.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Alerta | Texto, solo lectura | — | Pago próximo · Presupuesto al 80 % y al 100 % · Precios viejos antes del día de despensa · Plan de comidas sin armar (fijo, 00001) | `listNotifications` (`name`) | — | N-298 |
| Activa | Interruptor | — | Encendida (decisión F; `is_enabled` = 1 al sembrar) | `listNotifications` (`is_enabled`) | `updateNotification` (`is_enabled`) | N-300 |
| Destinatarios | Selección múltiple: personas con chat confirmado | Solo chats confirmados · «Primero confirma el chat de esta persona» | Sin valor: lo captura Axel (decisión A) | `listNotifications` (`recipient_ids`) y `listNotificationRecipients` | `setNotificationRecipients` | N-300 |
| Días antes del pago (pago próximo) | Lista de días en fichas: 1 = el día anterior, 0 = el mismo día | Uno o más enteros ≥ 0 · «Escribe al menos un día; 0 es el mismo día del pago» | [1, 0] (decisión F, K44) | `getAlertSettings` (`alerts.payment_due_days_before`) | `updateAlertSettings` | N-323 |
| Primer aviso del presupuesto | Porcentaje entero | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | 80 % (decisión F, K44) | `getAlertSettings` (`alerts.budget_warning_pct`) | `updateAlertSettings` | N-323 |
| Aviso de tope del presupuesto | Porcentaje entero | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | 100 % (decisión F, K44) | `getAlertSettings` (`alerts.budget_limit_pct`) | `updateAlertSettings` | N-323 |
| Días antes del día de despensa (precios viejos) | Número entero de días | Entero ≥ 0 · «Escribe cuántos días antes; 0 es el mismo día» | 1 día (Q2-73, K50); clave `alerts.stale_prices_advance_days` (Q2-11 B, K47) | `getAlertSettings` (`alerts.stale_prices_advance_days`) | `updateAlertSettings` | N-300 |
| Día de revisión (plan sin armar) | Selector de lunes a domingo (1 a 7) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Domingo (7), decisión F; clave `alerts.meal_plan_missing_weekday` (Q2-11 B, K47) | `getAlertSettings` (`alerts.meal_plan_missing_weekday`) | `updateAlertSettings` | N-300 |
| Hora (avisos con hora) | `HH:MM` de 24 h, hora del centro | De 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | 09:00 (Q2-73, K50); clave `alerts.send_time` (Q2-11 B, K47) | `getAlertSettings` (`alerts.send_time`) | `updateAlertSettings` | N-300 |

Reglas fijas que se leen en la fila (no son campos): el aviso de presupuesto mira **cada categoría y el total** (decisión F); «casillas sin armar» = pares día × tiempo obligatorio (desayuno, comida, cena) en que alguna persona no tiene receta ni «fuera», máximo 21 por semana (K43: es una regla, no un parámetro). A qué alertas aplica la hora lo fija S7 (ver SUPUESTOS).

**Chats confirmados** (registro de chat, N-299; nunca se muestran el `chat_id` ni el token del bot, K30).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Persona | Texto, solo lectura | — | Axel y Ximena (las dos personas con cuenta; Karen no recibe avisos, decisión A) | `listNotificationRecipients` (`person_id`) y `listPersons` | — | N-299 |
| Chat registrado | Sí o no, con fecha de comprobación | — | Calculado («No» hasta confirmar) | `listNotificationRecipients` (`is_registered`, `verified_at`) | — | N-299 |
| Etiqueta | Texto opcional | Texto libre; vacío borra la etiqueta | Sin valor: lo captura Axel | `listNotificationRecipients` (`label`) | `updateNotificationRecipient` (`label`) | N-299 |
| Recibe avisos | Interruptor | — | El que deja la confirmación del chat (`is_active` del destinatario creado) | `listNotificationRecipients` (`is_active`) | `updateNotificationRecipient` (`is_active`) | N-299 |

**Chats detectados por confirmar** (K4).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre en Telegram | Texto, solo lectura | — | El que da Telegram | `listTelegramChatCandidates` (`telegram_name`) | — | N-299 |
| Detectado | Fecha y hora, solo lectura | — | Cuando el bot recibió el mensaje | `listTelegramChatCandidates` (`detected_at`) | — | N-299 |
| Terminación del chat | Texto «•••» con los últimos 4 caracteres | — | Calculado, para distinguir dos chats con el mismo nombre | `listTelegramChatCandidates` (`chat_id_masked`) | — | N-299 |
| Confirmar como | Selector de persona (Axel o Ximena) + botón «Confirmar» | Un chat ya confirmado responde 409 · «Este chat ya estaba confirmado; recarga la lista» | La persona del usuario de la sesión (K4, `TelegramChatCandidateConfirm`) | — | `confirmTelegramChatCandidate` (con `Idempotency-Key`) | N-316 |
| Descartar | Botón (acción crítica) | — | — | — | `dismissTelegramChatCandidate` | N-324 |

Texto de ayuda sobre la lista: «Para recibir avisos, escribe un mensaje al bot desde tu Telegram y confírmalo aquí». El texto exacto lo fija S7 (WF-3) según cómo detecte el chat (K4); es solo redacción, no datos.

**Parámetros operativos de la sección (K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Primer aviso del presupuesto | `alerts.budget_warning_pct` | Actual: `value` · Inicial: 80 % | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | Vuelve a 80 % | Decisión F, K44 | `getAlertSettings` · `updateAlertSettings` |
| Aviso de tope del presupuesto | `alerts.budget_limit_pct` | Actual: `value` · Inicial: 100 % | Entero ≥ 0 · «Escribe un porcentaje entero, sin decimales» | Vuelve a 100 % | Decisión F, K44 | `getAlertSettings` · `updateAlertSettings` |
| Días antes del pago próximo | `alerts.payment_due_days_before` | Actual: `value` · Inicial: [1, 0] | Lista de uno o más enteros ≥ 0 · «Escribe al menos un día; 0 es el mismo día del pago» | Vuelve a [1, 0] | Decisión F, K44 | `getAlertSettings` · `updateAlertSettings` |
| Días antes del día de despensa (precios viejos) | `alerts.stale_prices_advance_days` | Actual: `value` · Inicial: 1 día | Entero ≥ 0 · «Escribe cuántos días antes; 0 es el mismo día» | Vuelve a 1 | Q2-73 (K50), K43; lugar Q2-11 B (K47) | `getAlertSettings` · `updateAlertSettings` |
| Hora de los avisos con hora | `alerts.send_time` | Actual: `value` · Inicial: 09:00 | `HH:MM` de 24 h, de 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Vuelve a 09:00 | Q2-73 (K50), K43; lugar Q2-11 B (K47) | `getAlertSettings` · `updateAlertSettings` |
| Día de revisión del plan sin armar | `alerts.meal_plan_missing_weekday` | Actual: `value` · Inicial: 7 (domingo) | Entero de 1 a 7 · «Elige un día de lunes a domingo» | Vuelve a 7 | Decisión F; lugar Q2-11 B (K47) | `getAlertSettings` · `updateAlertSettings` |

**Acciones.** «Enviar aviso de prueba» (acción primaria): elige alerta y persona y llama `sendTestNotification`; respuesta: «Aviso de prueba enviado a Ximena» o el error en palabras simples («El bot no pudo entregar el mensaje · Revisa que el chat esté confirmado»). «Descartar» un chat detectado abre la acción crítica: «Descartar el chat de “nombre en Telegram” (•••nnnn): se borra de esta lista y el bot no le mandará avisos. Si era de alguien de la casa, tendrá que volver a escribirle al bot. · Descartar chat». Los chats ya confirmados no se tocan (N-324).

**Vista previa del efecto.**
- Umbrales: «Con el presupuesto total de este mes ($T), el primer aviso llega al gastar $A y el de tope al gastar $B», con $T de `getBudgetMonth` y $A = $T × porcentaje ÷ 100 en centavos con ROUND_HALF_UP (decisión B); igual por categoría al tocar una categoría de ejemplo.
- Días antes del pago: «El próximo pago (Vie dd mmm aaaa) avisaría el Jue dd mmm y el Vie dd mmm», con la fecha de `getDashboardNextPayment` y un renglón por día de la lista.
- Precios viejos y plan sin armar: «Llegará n días antes del día de despensa» y «Se revisa cada domingo» (o el día elegido); la hora, en los avisos donde S7 la aplique, sale de `alerts.send_time` (09:00 al inicio).

**Estados.** Vacío: «Chats detectados por confirmar» sin filas = «No hay chats nuevos. Escribe al bot desde Telegram y vuelve aquí» (con «Actualizar»); «Chats confirmados» sin chat = «Ximena aún no confirma su chat» en la fila. Carga: esqueleto de las tres tablas. Error: «No se pudieron cargar los avisos · Reintentar»; al guardar, el campo vuelve a su valor. Guardado: automático, «Guardado»; confirmar un chat lo pasa a «Chats confirmados» con «Chat confirmado».

**Edición en el lugar.** Activa, destinatarios, días, porcentajes, hora, etiqueta y «Recibe avisos». No se editan: nombre de la alerta, chat registrado, datos del chat detectado (vienen del bot).

**Atajos.** `↑`/`↓` y `Enter` en las tres tablas; `Ctrl+/` activa o desactiva a la vez las alertas seleccionadas; `Esc` cancela.

**Móvil y escritorio.** 412 px: tarjetas por alerta (nombre, activa y resumen «1 día antes y el mismo día») que abren su hoja; los chats detectados arriba cuando los hay, con «Confirmar» de ancho completo. 1440 px: las tres tablas una debajo de otra; parámetros en la fila de su alerta.

## P-26.10 Apariencia

**Propósito.** Elegir el tema de este dispositivo: según el sistema (por defecto), claro u oscuro (decisión G; N-301), y encender o apagar los atajos de una tecla (Q2-82 A, K51; WCAG 2.1.4, DESIGN §9.6). Dos controles, los dos del dispositivo: el tema (U1 §2.3) y el interruptor que ratificó la puerta 2 (U1 aún no lo lista: solicitud en `docs/estado/p2c/U4-puerta2.md`).

**Vista de base de datos.** Ninguna (un solo control).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Tema | Control segmentado: Sistema · Claro · Oscuro | Uno de los tres | Sistema (decisión G) | Local del dispositivo (U1 §0.3: no usa C1) | Local del dispositivo (también `Ctrl+Shift+L` y `/tema`) | N-301 |
| Atajos de una tecla | Interruptor | — | Encendido (Q2-82 A, K51; DESIGN §9.6) | Local del dispositivo (no usa C1) | Local del dispositivo | N-301 |

Valores fijos del sistema de diseño que se leen en la sección (no son campos; DESIGN §1.2 P4 y P6, §2.1): «Saturación: apagada» y «Densidad: cómoda» (decisión G). No se configuran en esta versión (ver SUPUESTOS).

**Vista previa del efecto.** Atajos de una tecla: al apagarlo, `/`, `?` y las flechas dejan de actuar como atajos fuera de los campos de texto (DESIGN §9, criterio 2.1.4; U1 §5.4); los atajos con `Ctrl` o `Alt` siguen. Tema: el cambio es su propia vista previa: la página cambia de tema al instante (sin transición si el sistema pide movimiento reducido, DESIGN §6.3). **Estados.** Sin vacío, carga ni error (no hay red); guardado inmediato en el dispositivo, «Guardado». Sin conexión funciona igual. **Edición en el lugar:** el control mismo. **Atajos:** `Ctrl+Shift+L`, `/tema`. **Móvil y escritorio:** 412 px, tres botones de 44 px a todo el ancho; 1440 px, el control a la izquierda del contenido, con una muestra de tarjeta en claro y oscuro.

## P-26.11 Datos y respaldos

**Propósito.** Ver y operar los respaldos locales (copia diaria versionada, `integrity_check`, restauración probada; decisión H), importar `gastos_fijos_ingresos.csv`, revisar los datos semilla, encender o apagar la IA opcional (desactivada por defecto, decisión H), ver el estado del servicio y ajustar los parámetros de datos (K44). Necesidades N-302 a N-308 y N-325.

**Vista de base de datos.** Tabla de respaldos (vista por defecto, U1 §2.5; filtro por periodo, orden por fecha descendente). Los demás bloques son fichas: Importar, Datos semilla (3 filas), IA opcional (2 filas), Estado del servicio y Parámetros.

**Respaldos.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Fecha y hora | Fecha y hora, solo lectura | — | Calculado | `listBackups` (`started_at`, `finished_at`) | — | N-302 |
| Tipo | Etiqueta: Diario · Manual · Previo a una migración | — | Calculado | `listBackups` (`kind`) | — | N-302 |
| Resultado | Etiqueta: En curso · Correcto · Falló (con icono) | — | Calculado | `listBackups` (`status`, `error`) | — | N-302 |
| Tamaño | Tamaño en MB o GB | — | Calculado | `listBackups` (`size_bytes`) | — | N-302 |
| Integridad | Texto: «Íntegro» o el error de `integrity_check` | — | Calculado | `listBackups` (`integrity_check`) | — | N-302 |
| Restauración probada | Fecha y resultado | — | Calculado | `listBackups` (`restore_tested_at`, `restore_test_result`) | — | N-302 |
| Ubicación y huella | Ruta en la carpeta o disco de respaldos y huella abreviada | — | Calculado | `listBackups` (`file_path`, `sha256`) | — | N-302 |
| Próximo respaldo diario | Fecha y hora, calculada | — | Calculado | `listJobs` (`next_run_at` de la tarea del respaldo) | — | N-302 |
| Hora del respaldo diario | `HH:MM` de 24 h | De 00:00 a 23:59 · «Escribe la hora con dos dígitos, como hh:mm» | Sin valor inicial: lo proponen S6 y WF-4 (K31 punto 3, puerta 4), PROVISIONAL | `listJobs` (`schedule`) | `updateJob` (`schedule` de la tarea del respaldo), PROVISIONAL (ver SOLICITUDES del reporte) | N-302 |

Encabezado de la tabla: «Último respaldo correcto: dd-mmm-aaaa hh:mm · Próximo: … · Última restauración probada: …».

**Acciones de respaldo.** «Respaldar ahora» (acción primaria; también `/respaldo`): `createBackup`; al aceptar, la fila nueva aparece «En curso» y cambia sola a «Correcto» o «Falló». «Restaurar» en cada fila correcta: `restoreBackup`, con **doble confirmación de acción crítica** (U4-08, N-303): paso 1, vista previa del efecto; paso 2, escribir `RESTAURAR` (`BackupRestoreRequest.confirmation`). Texto: «Restaurar el respaldo del dd-mmm-aaaa hh:mm: la base vuelve a como estaba en ese momento y se pierden los cambios hechos después en los dos teléfonos y en la computadora. Antes se toma un respaldo de hoy por si necesitas regresar. · Restaurar respaldo». Al terminar: «Respaldo restaurado · Íntegro · Respaldo previo guardado» (`RestoreResult`: `integrity_check`, `pre_restore_backup_id`) y todas las sesiones recargan.

**Importar gastos fijos e ingresos.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Archivo | Archivo CSV (plantilla de `plan.md` §8.2) | Solo `.csv` con las columnas de la plantilla · «El archivo no tiene las columnas de la plantilla: …» (renglón y columna en palabras) | Sin valor: lo elige Axel | — | `previewPlannedItemsImport` (vista previa) | N-304 |
| Lo que se creará | Tabla de vista previa, solo lectura | Renglones con error marcados con icono y texto | Calculado | `previewPlannedItemsImport` | — | N-304 |
| Importar | Botón primario | Se activa solo sin errores · «Corrige los renglones marcados antes de importar» | — | — | `importPlannedItems` | N-304 |

Resultado: «n partidas creadas» (`ImportResult.created_count`) con enlace a Presupuesto (P-26.02).

**Datos semilla.** Tres filas fijas de `getSeedStatus` (N-305): compra del 10-oct-2026, crédito Kueski con sus despensas y pagos, CSV de gastos fijos e ingresos; cada una con «Cargada» o «Pendiente» y su número de registros. Solo lectura.

**IA opcional** (decisión H; la IA solo se usa para OCR de tickets y sugerencias de recetas; el núcleo es determinista).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Proveedor | Texto: «Modelo local (Ollama)» · «API de Claude» | — | Fijo (2 filas) | `listAiProviders` (`provider`) | — | N-306 |
| Activa | Interruptor | Encender pide acción crítica | **Apagada** los dos (decisión H; `ai.local_model_enabled` y `ai.claude_api_enabled` = false en 00001) | `listAiProviders` (`is_enabled`) | `updateAiProvider` (`is_enabled`) | N-307 |
| Usos permitidos | Etiquetas: OCR de tickets · Sugerencias de recetas | — | Fijo (decisión H) | `listAiProviders` (`allowed_uses`) | — | N-306 |
| Disponible | Sí o no (solo si está activa) | — | Calculado | `listAiProviders` (`is_available`) | — | N-306 |
| Llave | «Configurada» o «Sin configurar» (solo API de Claude); la llave nunca se muestra | — | Calculado | `listAiProviders` (`is_configured`) | — | N-306 |

Encender abre la acción crítica: modelo local, «Activar el modelo local: el OCR de tickets y las sugerencias de recetas usarán el modelo instalado en esta máquina. Siempre pedirá confirmación antes de guardar un ticket. · Activar modelo local»; API de Claude, «Activar la API de Claude: las fotos de tickets y los pedidos de sugerencias de recetas saldrán de la red de casa hacia el servicio de Claude. Siempre pedirá confirmación antes de guardar un ticket. · Activar API de Claude». Apagar es inmediato y sin diálogo. El OCR local con Tesseract no es IA: vive en Tiendas y precios › Fuentes de precio (P-26.06, parte a).

**Estado del servicio.** `getSystemStatus` (N-308): versión, en marcha desde y tamaño de la base de datos. Solo lectura.

**Parámetros operativos de la sección (K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Tamaño máximo de un adjunto de gasto | `attachment.max_size_mb` | Actual: `value` · Inicial: 10 MB | Entero ≥ 0, en MB · «Escribe un número entero de MB» | Vuelve a 10 MB | K43, K44 | `getDataSettings` · `updateDataSettings` |
| Retención de las claves de idempotencia | `idempotency.retention_hours` | Actual: `value` · Inicial: 24 h | Entero ≥ 0, en horas · «Escribe un número entero de horas» | Vuelve a 24 h | K6, K44 | `getDataSettings` · `updateDataSettings` |
| Tiempo máximo del OCR de un ticket | `ocr.timeout_seconds` | Actual: `value` · Inicial: 30 s | Entero ≥ 0, en segundos · «Escribe un número entero de segundos» | Vuelve a 30 s | Q2-76 (K50); sección: la de su `@param` en la 00004 (PROPUESTA de S1 que la puerta 2 no preguntó, `docs/estado/p2c/S1-puerta2.md`) | `getDataSettings` · `updateDataSettings` (C1 v1.3 la expone; solicitud de S1) |

Parámetros de K31 de esta sección **sin clave ni operación todavía** (no son campos; se piden en SOLICITUDES): cuántas versiones del respaldo diario se guardan (K31 punto 3) y los permisos de los archivos de la base y de los respaldos (K31 punto 8).

**Vista previa del efecto.** Restaurar: paso 1 con fecha, tamaño, integridad del respaldo elegido y la advertencia de lo que se pierde (arriba). Importar: tabla de lo que se creará antes de confirmar. Adjunto: «Con n MB, un ticket o adjunto más pesado se rechaza con “El archivo pesa más de n MB”». Retención: «La API recuerda cada captura repetida n horas para no duplicarla; los gastos y pesos de la cola del teléfono nunca se duplican porque su clave queda en el registro (K9)». IA: el texto del diálogo de activación. OCR: «Tiempo máximo para leer un ticket: {antes} s → {después} s».

**Estados.** Vacío: sin respaldos = «Aún no hay respaldos · Respaldar ahora». Carga: esqueleto de la tabla y de las fichas. Error: «No se pudo cargar la lista de respaldos · Reintentar»; un respaldo que falla queda en la tabla como «Falló» con su error resumido; una restauración fallida dice «No se restauró; la base sigue como estaba» (`integrity_check` distinto de ok). Guardado: automático en parámetros y en el interruptor de IA (tras confirmar).

**Edición en el lugar.** Parámetros, hora del respaldo (PROVISIONAL: sin valor inicial hasta K31 punto 3; `updateJob` según SOLICITUDES del reporte y supuesto 4) e interruptores de IA. No se editan: filas de respaldos, semillas, estado del servicio (calculados).

**Atajos.** `/respaldo`; `↑`/`↓` y `Enter` en la tabla de respaldos (Enter abre el detalle con «Restaurar»); `Esc` cierra el diálogo.

**Móvil y escritorio.** 412 px: «Respaldar ahora» arriba a todo el ancho, respaldos como tarjetas (fecha, resultado, integridad) y las demás fichas apiladas; «Restaurar» dentro del detalle, nunca en la tarjeta. 1440 px: tabla completa y fichas en dos columnas.

## P-26.12 Personas y reparto

**Propósito.** Ver quién participa en el dinero del hogar y cómo: Axel (titular, con cuenta), Ximena (co-titular, mismos permisos) y Karen (participante externa, sin cuenta ni acceso, paga directo a Kueski su parte; decisión A); fijar el reparto por defecto del conjunto (50/50, decisión B) y, si Axel quiere, un tope de pago por quincena por persona (K13). Necesidades N-309 y N-310.

**Vista de base de datos.** Tabla de personas (única vista, U1 §2.5; orden Axel, Ximena, Karen; sin filtros). Debajo, tabla «Repartos reutilizables» (las reglas con `is_template` de `listCostShareRules`: Hogar, Karen, Conjunto 50/50). No se crean ni se borran personas ni reglas aquí (C1 v1.2 no tiene esas operaciones).

**Personas.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre visible | Texto | No vacío · «Escribe un nombre» | Axel · Ximena · Karen (semilla) | `listPersons` (`display_name`) | `updatePerson` (`display_name`) | N-310 |
| Papel | Etiqueta, solo lectura: Titular · Co-titular · Participante externa | — | Fijo (decisión A) | `listPersons` (`participation`) | — | N-309 |
| Cuenta | Sí o no, solo lectura | — | Axel y Ximena sí; Karen nunca (decisión A) | `listPersons` (`has_account`) | — | N-309 |
| Paga su parte directo al acreedor | Sí o no, solo lectura | — | Karen sí (paga a Kueski; su parte es informativa y repone la línea, decisión A y K15) | `listPersons` (`pays_share_directly`) | — | N-309 |
| Qué hace en el sistema | Texto fijo según el papel | — | Titular y co-titular: «Mismos permisos: captura, consulta y ajustes». Karen: «Sin cuenta. Solo aparece en la despensa en común y en sus compras con Kueski; recibe su estado de cuenta en PDF, que Axel reenvía» | `listPersons` (`participation`) | — | N-309 |
| Tope de pago por quincena | Dinero, opcional | Vacío o monto ≥ $0.00 con dos decimales · «Escribe un monto como $1,234.56 o deja el campo vacío» | **Vacío** (sin tope, K13: Axel lo captura cuando quiera); lo usa solo el comparador como objetivo, no como restricción dura (decisión C) | `listPersons` (`fortnight_payment_cap_cents`) | `updatePerson` (`fortnight_payment_cap_cents`; vacío = `null`) | N-310 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listPersons` (`notes`) | `updatePerson` (`notes`) | N-310 |

**Repartos reutilizables.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre | Texto | No vacío · «Escribe un nombre para el reparto» | Hogar · Karen · Conjunto 50/50 (semilla) | `listCostShareRules` (`name`) | `updateCostShareRule` (`name`) | N-310 |
| Parte de cada persona | Por persona: porcentaje (dos decimales) o monto fijo | Cada parte es porcentaje (0.00 % a 100.00 %, dos decimales, V-06) o monto fijo (≥ $0.00, V-03), nunca los dos (00001, `cost_share_rule_part`). Las partes con porcentaje suman a lo más 100 % · V-16: «Las partes con porcentaje suman {suma} %. No pueden pasar de 100 %: lo que falta lo recibe {persona del resto}.» | Conjunto: Karen 50 % y el hogar 50 % (decisión B); Hogar: todo al hogar; Karen: todo a Karen | `listCostShareRules` (`parts`: `pct_bp`, `fixed_cents`) | `updateCostShareRule` (`parts`) | N-310 |
| Quién toma el resto y los centavos | Texto, solo lectura | — | El titular (K14: las partes de personas externas con ROUND_HALF_UP y el titular toma el resto) | `listCostShareRules` (`remainder_person_id`) | — | N-309 |
| Reparto por defecto | Marca, solo lectura | — | «Hogar» (se aplica cuando ni el gasto ni su categoría traen reparto) | `listCostShareRules` (`is_default`) | — | N-309 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listCostShareRules` (`notes`) | `updateCostShareRule` (`notes`) | N-310 |

**Regla única del reparto (RN-03, S1-08, `CostShareRule` de C1).** Se editan solo las partes explícitas (`parts`: `pct_bp` o `fixed_cents` por persona). La persona del resto (`remainder_person_id`, el titular en las tres reglas sembradas) no tiene campo: su parte se muestra al final, calculada y de solo lectura, como «el resto» = 100 % menos las partes con porcentaje (al aplicarse a un gasto, además menos los montos fijos; RN-03: «El resto recibe la suma menos las partes»). Así cada regla suma 100 % por construcción y no hay que cuadrarla a mano; el único error posible es que las partes con porcentaje pasen de 100 % (V-16), y entonces no se guarda. Un monto fijo mayor que el gasto se rechaza al aplicar la regla al gasto (RN-03), no aquí. En pantalla: «Conjunto 50/50: Karen 50.00 % · el resto 50.00 % (Axel, hogar)». Qué reparto usa cada categoría se elige en Categorías (P-26.04, parte a) y cada gasto en P-03 y P-04.

**Vista previa del efecto.** Al cambiar una parte, antes de guardar: «Con un gasto de $M: Karen $K (M × p con ROUND_HALF_UP) y el hogar $H (el resto)», con $M marcador que la persona puede cambiar en la misma vista previa. Al salir del último campo de la regla, antes de guardar, el diálogo de alcance (Q2-68 C; bloqueo `docs/estado/bloqueos/U4-reparto-retroactivo.md` RESUELTO; como el del presupuesto en P-26.02, K2) pregunta «¿A qué gastos aplica el reparto nuevo?» con dos opciones: «Solo a los gastos nuevos» (los gastos ya capturados conservan las partes que se congelaron al guardarlos y los estados de cuenta ya entregados a Karen no cambian) y «También a los gastos ya capturados» («{n} gastos ya capturados con esta regla cambian sus partes, desde el {fecha del más antiguo}; incluye meses anteriores»), con «Confirmar» y «Cancelar»; `Esc` cancela sin cambiar (parte a §0.2). **La preselección es PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/U4-reparto-preseleccion.md`, se pregunta en la puerta 4): ni Q2-68 C ni K46 la fijan y C1 v1.3 no da valor por omisión. En uso, «Solo a los gastos nuevos» preseleccionada; las otras opciones son preseleccionar «También a los gastos ya capturados» o no preseleccionar ninguna, con «Confirmar» inactivo hasta elegir. Escritura: `updateCostShareRule` con `parts` y `scope` (C1 v1.3, `CostShareRuleUpdate.scope`, obligatorio con `parts` y sin valor por omisión): `onward` = «Solo a los gastos nuevos» y `all` = «También a los gastos ya capturados». C1 v1.2 no recibe alcance y ya se comporta como «solo los nuevos»: las partes de los gastos ya guardados están congeladas en `expense_share`; «también los ya capturados» necesita `scope` = `all` de C1 v1.3. Qué gastos abarca `all` (solo los que nombran la regla, o también los que la toman por su categoría o por ser la de por defecto) es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/C1-reparto-alcance.md`), y con él el «{n}» y la fecha de «También a los gastos ya capturados». C1 v1.3 no da ese conteo ni esa fecha (solicitud SC-U4-P2-1 en `docs/estado/p2c/U4-puerta2.md`): mientras no los dé, «{n}» y la fecha se rotulan «Vista previa pendiente» (parte a §0.2). Tope por quincena: «El comparador usará $X como tope de esta persona al evaluar el objetivo “pago por quincena parejo y bajo tope”; no descarta escenarios» (K17: pico contra el tope cuando hay tope; decisión C: no es restricción dura).

**Estados.** Vacío: no aplica (las 3 personas y las 3 reglas vienen sembradas). Carga: esqueleto de las dos tablas. Error al guardar: el de la parte a §0.3, con las partes de vuelta a su valor. Guardado: automático al salir del campo; las partes se guardan juntas al confirmar el alcance (Q2-68 C).

**Edición en el lugar.** Nombre visible, tope, notas, nombre y partes de las reglas. No se editan: papel, cuenta, paga directo, persona del resto, reparto por defecto (decididos o fijos).

**Atajos.** `↑`/`↓` y `Enter` en las dos tablas; `Esc` cancela la edición de una parte.

**Móvil y escritorio.** 412 px: tres tarjetas de persona (nombre, papel, tope) y tres de reparto con la barra de partes etiquetada (nombre y porcentaje en texto, nunca solo color: K41 y DESIGN §2.6, Karen en su color categórico con su etiqueta). 1440 px: las dos tablas con la vista previa a la derecha de la regla que se edita.

## P-26.13 Salud y metas de ahorro

**Propósito.** Los perfiles de salud de Axel y Ximena (plan de comidas y salud son solo de ellos dos, decisión A): estatura, nacimiento, sexo, actividad, objetivo, IMC y la meta de calorías sugerida que la persona confirma o cambia (decisión D, K23); el ajuste de la meta por objetivo de la casa y de cada persona (K36); las metas de ahorro con fecha y la aportación mensual planificada al fondo (decisión B). Necesidades N-311, N-312, N-313, N-326 y N-327.

**Vista de base de datos.** Perfiles: una ficha por persona (2), sin vista de lista (U1 §2.5). Metas de ahorro: **tabla** por defecto en esta sección; también tablero por estado (`status`, motor S2) con tres columnas, «En curso», «Cumplida» y «Vencida» (Q2-18 A; bloqueo `docs/estado/bloqueos/S2-metas.md`, pregunta 3, RESUELTO; `in_progress`, `achieved` y `overdue` en `docs/contratos/motores/S2-interfaz.md` 1.3 §3.13). C1 v1.3 todavía declara `status` con `on_track`, `behind` y `achieved`, y U1 §2.5 nombra las columnas en ruta, atrasada y cumplida: los dos tienen que seguir a Q2-18 A (solicitudes en `docs/estado/p2c/correccion-U4-r1.md`), calendario por fecha meta y galería con el avance (U1 §2.5; la galería es la vista por defecto de P-09). Filtro por defecto: activas; orden: fecha meta ascendente.

**Perfil de salud (una ficha por persona).**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Estatura | Metros con dos decimales (se guarda en mm) | Mayor que 0 · «Escribe la estatura en metros, como 1.75» | Sin valor: lo captura la persona | `listHealthProfiles` (`height_mm`) | `updateHealthProfile` o `createHealthProfile` (`height_mm`) | N-311, N-313 |
| Fecha de nacimiento | Fecha (`10-oct-2026` en pantalla) | No después de hoy · «La fecha de nacimiento no puede ser futura» | Sin valor: lo captura la persona | `listHealthProfiles` (`birth_date`; edad en `age_years`) | `updateHealthProfile` (`birth_date`) | N-311, N-313 |
| Sexo | Selector: Mujer · Hombre | Uno de los dos (lo pide la fórmula) · «Elige una opción para calcular la meta» | Sin valor: lo captura la persona | `listHealthProfiles` (`sex`) | `updateHealthProfile` (`sex`) | N-311, N-313 |
| Actividad | Selector de 5 niveles con su factor: sedentaria 1.2 · ligera 1.375 · moderada 1.55 · activa 1.725 · muy activa 1.9 | Uno de los cinco | Sin valor: lo captura la persona (K23: 5 niveles de 1.2 a 1.9) | `listHealthProfiles` (`activity_level`) | `updateHealthProfile` (`activity_level`) | N-311, N-313 |
| Objetivo | Selector: Bajar · Mantener · Subir de peso | Uno de los tres o vacío | **Vacío** («todavía no lo elige», K23) | `listHealthProfiles` (`goal`) | `updateHealthProfile` (`goal`) | N-311, N-313 |
| Peso actual | Kg con fecha, solo lectura | — | Calculado del historial (se registra en P-22) | `listHealthProfiles` (`current_weight_g`, `current_weight_date`) | — | N-311 |
| IMC | Número con un decimal y clase OMS (bajo peso, normal, sobrepeso, obesidad) | — | Calculado (K7: cortes 18.5, 25 y 30) | `listHealthProfiles` (`bmi_tenths`, `bmi_class`) | — | N-311 |
| Meta sugerida | Kcal por día con desglose, solo lectura | — | Calculado (RN-17: Mifflin-St Jeor × actividad − ajuste de bajar o + ajuste de subir; un redondeo ROUND_HALF_UP al final) | `listHealthProfiles` (`suggested_kcal_target`) | — | N-311 |
| Meta confirmada | Kcal por día (entero) | Entero mayor que 0 · «Escribe las kilocalorías por día, sin decimales» | **Vacía** hasta que la persona confirme la sugerida o escriba otra (K23) | `listHealthProfiles` (`daily_kcal_target`) | `updateHealthProfile` (`daily_kcal_target`); botón «Usar la sugerida» | N-313 |
| Mi ajuste para bajar | Kcal (entero), opcional | Vacío o entero ≥ 0 · «Escribe un número entero de kcal, 0 o más, o deja vacío para usar el de la casa» | **Vacío = usa el de la casa** y lo dice: «El de la casa: n kcal» (K36) | `listHealthProfiles` (`goal_adjust_lose_kcal`, `household_goal_adjust_lose_kcal`) | `updateHealthProfile` (`goal_adjust_lose_kcal`; «Usar el de la casa» = `null`) | N-327 |
| Mi ajuste para subir | Kcal (entero), opcional | Igual que el anterior | **Vacío = usa el de la casa** (K36) | `listHealthProfiles` (`goal_adjust_gain_kcal`, `household_goal_adjust_gain_kcal`) | `updateHealthProfile` (`goal_adjust_gain_kcal`; «Usar el de la casa» = `null`) | N-327 |
| Restricciones alimentarias | Texto opcional | — | **Vacío** (decisión D: el campo existe y va vacío) | `listHealthProfiles` (`dietary_restrictions`) | `updateHealthProfile` (`dietary_restrictions`) | N-311, N-313 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listHealthProfiles` (`notes`) | `updateHealthProfile` (`notes`) | N-313 |

«Mantener» **no ajusta y no tiene campo** (K36): con ese objetivo, los dos ajustes se ven atenuados con «Con “mantener” no hay ajuste». El campo que corresponde al objetivo elegido va primero y resaltado. «Usar el de la casa» es el «restablecer» del ajuste de la persona (U4-13). Si la persona aún no tiene perfil: «Sin perfil de salud · Crear perfil» (`createHealthProfile` con `person_id`). El mismo ajuste por persona también se edita en P-21 (N-327).

**Ajuste por objetivo de la casa (K36, K44).**

| Parámetro | Clave | Valor actual e inicial | Validación y mensaje | Restablecer | Decisión | Lectura y escritura |
|---|---|---|---|---|---|---|
| Ajuste de la casa para bajar | `health.goal_adjust_lose_kcal` | Actual: `value` · Inicial: 0 kcal | Entero ≥ 0 · «Escribe un número entero de kcal, 0 o más» | Vuelve a 0 kcal | K36, K44 | `getHealthSettings` · `updateHealthSettings` |
| Ajuste de la casa para subir | `health.goal_adjust_gain_kcal` | Actual: `value` · Inicial: 0 kcal | Entero ≥ 0 · «Escribe un número entero de kcal, 0 o más» | Vuelve a 0 kcal | K36, K44 | `getHealthSettings` · `updateHealthSettings` |

Debajo: «Lo usan las personas que dejan vacío su propio ajuste: hoy Axel y Ximena» (según `goal_adjust_*` nulos de `listHealthProfiles`).

**Metas de ahorro y fondo.**

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Nombre | Texto | No vacío · «Escribe un nombre para la meta» | Sin valor: lo captura Axel | `listSavingsGoals` (`name`) | `createSavingsGoal` / `updateSavingsGoal` (`name`) | N-312, N-313 |
| Tipo | Selector: Fondo de ahorro · Meta con fecha | Obligatorio · «Elige fondo o meta» | Sin valor: lo elige Axel | `listSavingsGoals` (`kind`) | `createSavingsGoal` / `updateSavingsGoal` (`kind`) | N-312, N-313 |
| Monto objetivo | Dinero | Obligatorio en una meta; monto > $0.00 · «Escribe el monto de la meta, como $1,234.56» | Sin valor: lo captura Axel | `listSavingsGoals` (`target_cents`) | ídem (`target_cents`) | N-312, N-313 |
| Fecha meta | Fecha | Obligatoria en una meta; no antes del inicio · «La fecha meta no puede quedar antes del inicio» | Sin valor: lo captura Axel | `listSavingsGoals` (`target_date`) | ídem (`target_date`) | N-312, N-313 |
| Aportación mensual planificada | Dinero | Vacío o ≥ $0.00 · «Escribe un monto como $1,234.56» | Sin valor: lo captura Axel (decisión B: aportación mensual planificada) | `listSavingsGoals` (`planned_monthly_cents`) | ídem (`planned_monthly_cents`) | N-312, N-313 |
| Inicio | Fecha | Obligatoria · «Escribe la fecha de inicio» | Sin valor: lo captura Axel | `listSavingsGoals` (`start_date`) | ídem (`start_date`) | N-313 |
| Activa | Interruptor con vista previa en línea (parte a §0.2: no guarda al tocar; vista previa 3 de metas) | — | Encendida (valor de la base al crear; `is_active` de `savings_goal`) | `listSavingsGoals` (`is_active`) | `updateSavingsGoal` (`is_active`; «Desactivar», no se borra) | N-313 |
| Ahorrado, avance, aportación necesaria y estado | Solo lectura | — | Calculados (motor S2; `S2-metas` RESUELTO): avance con un decimal (`progress_tenths_pct` de C1 v1.3, décimas con ROUND_HALF_UP; Q2-17 C: un decimal en Metas y entero solo en la tarjeta de Inicio); aportación necesaria = (monto meta − ahorrado) ÷ meses que faltan (C1-datos M-16), con los meses del mes de hoy al de la fecha meta, los dos incluidos (Q2-16 A), base ROUND_HALF_UP y la última aportación del plan ajusta el centavo (`S2-interfaz.md` 1.3 §3.13); estado con tres valores, en curso, cumplida y vencida (Q2-18 A). La aportación de una meta vencida sin cumplir es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/S2-metas-vencida.md`; en uso, todo el restante en el mes de hoy) | `listSavingsGoals` (`saved_cents`, `progress_tenths_pct`, `required_monthly_cents`, `status`) | — | N-312 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listSavingsGoals` (`notes`) | ídem (`notes`) | N-313 |

La aportación planificada al fondo general es la fila de tipo «Fondo de ahorro» con su aportación mensual. Acción primaria: «Nueva meta» (`Alt+N`).

**Vista previa del efecto.** **Perfil de salud** (vista previa en línea, parte a §0.2, calculada en el navegador con RN-17; al guardar manda el servidor): al cambiar cualquier dato del perfil que entra a una fórmula, el bloque «Antes → Después» bajo el campo muestra las cifras que cambian. **Estatura** → IMC con su clase OMS (K7 y RN-17: `weight_g` × 1000 ÷ `height_mm`², con el último peso registrado, un decimal con ROUND_HALF_UP; cortes 18.5, 25 y 30) y meta sugerida. **Fecha de nacimiento** (edad en años cumplidos a hoy, America/Mexico_City), **sexo** (+ 5 en hombres, − 161 en mujeres), **actividad**, **objetivo** y **mi ajuste para bajar o subir** → meta sugerida (el IMC no depende de ellos y no se muestra). La meta sugerida va con su desglose: basal de Mifflin-St Jeor × factor de actividad, menos el ajuste de bajar o más el de subir (el de la persona o, si lo deja vacío, el de la casa), con un solo redondeo ROUND_HALF_UP al final (RN-17, K36). Si falta un insumo (sexo, fecha de nacimiento, estatura, peso, actividad u objetivo), el lado sin cálculo dice «Falta {dato}» (RN-17: sin sugerencia mientras falte). Ejemplos ilustrativos de RN-17 (no canon): hombre de 35 años, 72.5 kg, actividad moderada y objetivo mantener; al capturar la estatura de 1.75 m: «IMC: Falta la estatura → 23.7 (normal) · Meta sugerida: Falta la estatura → 1,648.75 × 1.55 = 2,555.5625 → 2,556 kcal»; al cambiar el objetivo a bajar con su ajuste de 500 kcal: «Meta sugerida: 2,556 → 2,056 kcal» (2,555.5625 − 500 = 2,055.5625). Con los valores de prueba del hombre de S4-07, al cambiar el ajuste para bajar: «Basal 1,748.75 × 1.2 = 2,098.50 − ajuste para bajar 300 = 1,798.50 → meta sugerida 1,799 kcal» (con los valores iniciales 0 y 0 la meta queda en 2,098.50 antes de redondear). La meta confirmada no cambia sola: la persona confirma con «Usar la sugerida» o escribe otra (K23). Ajuste de la casa: «Con n kcal, la meta sugerida de Axel pasaría de a a b» por cada persona que usa el de la casa (fórmula RN-17 en el navegador; al guardar manda el servidor). **Metas** (vista previa en línea, parte a §0.2, calculada en el navegador): (1) al cambiar el monto objetivo, «Avance: {ahorrado} de {objetivo} = {p} %» antes → después, con `saved_cents` de `listSavingsGoals` y la fracción exacta ahorrado ÷ monto objetivo (RN-18), redondeada en décimas con ROUND_HALF_UP, en el mismo formato que la columna «avance» de la tabla de metas (`progress_tenths_pct`, un decimal; Q2-17 C, bloqueo `docs/estado/bloqueos/S2-metas.md` pregunta 2, RESUELTO), para que «Antes» sea igual a lo que ya se ve. Ejemplo ilustrativo de RN-18: meta de $6,293.44 con $3,146.72 ahorrados = 50.0 % (314,672 × 1,000 ÷ 629,344 = 500 décimas); (2) al cambiar la aportación mensual planificada, «El plan de {mes en curso} reserva {aportación nueva} para esta meta (antes {aportación actual}); el total reservado para metas en el mes pasa de {suma actual} a {suma nueva}» (RN-18: la aportación planificada se reserva en el plan del mes; la suma es la de `planned_monthly_cents` de las metas activas); (3) al tocar «Activa» (interruptor con vista previa en línea, excepción de la parte a §0.2): «Total reservado para metas en {mes en curso}: {antes} → {después}», la misma suma de `planned_monthly_cents` de las metas activas con esta meta sumada (al encender) o quitada (al apagar); una meta sin aportación planificada no cambia el total y el bloque lo dice: «Sin aportación planificada: el total reservado no cambia»; el interruptor no guarda al tocar: `Enter` o salir de él guarda y `Esc` lo devuelve a como estaba. Con `Ctrl+/` sobre varias metas, el mismo total sumado sobre las seleccionadas antes de aplicar. **Aportación necesaria y estado** (`required_monthly_cents`, `status`, motor S2; FIRMES desde la puerta 2, bloqueo `docs/estado/bloqueos/S2-metas.md`, preguntas 1 y 3, RESUELTO): al cambiar el monto objetivo o la fecha meta, «Aportación necesaria: {antes} → {después} al mes · Estado: {antes} → {después}», calculado en el navegador en centavos enteros con la regla de `docs/contratos/motores/S2-interfaz.md` 1.3 §3.13 (sin la solicitud SC-U4-2 a C1, que solo hacía falta si no se podía calcular aquí): meses que faltan = del mes de hoy al mes de la fecha meta, los dos incluidos (Q2-16 A; una fecha meta de un mes anterior da 0); aportación necesaria = ROUND_HALF_UP((monto meta − ahorrado) ÷ meses que faltan) (C1-datos M-16; la última aportación del plan ajusta el centavo); estado (Q2-18 A) = «Cumplida» si lo ahorrado alcanza el monto meta (aportación $0.00); si no, «Vencida» si la fecha meta ya pasó; si no, «En curso». Ejemplo ilustrativo (no canon), hoy 11-oct-2026, con la meta de $6,293.44 y $3,146.72 ahorrados: con fecha meta en diciembre de 2026 faltan 3 meses y la aportación necesaria es $1,048.91 al mes (314,672 ÷ 3 = 104,890.67 → 104,891 centavos; la última, $1,048.90); al cambiar la fecha meta a noviembre de 2026, $1,048.91 → $1,573.36 al mes, estado «En curso» → «En curso». **PROVISIONAL** solo la aportación de una meta vencida sin cumplir (meses que faltan 0 y restante mayor que $0.00; bloqueo abierto `docs/estado/bloqueos/S2-metas-vencida.md`; en uso, la del motor: todo el restante en el mes de hoy): ahí el bloque rotula «Aportación necesaria: vista previa pendiente; se recalcula al guardar» (parte a §0.2) y muestra el estado «Vencida».

**Estados.** Vacío: sin perfil = «Sin perfil de salud · Crear perfil»; sin metas = «Sin metas de ahorro · Nueva meta»; sin peso = «Sin peso registrado · Registrar peso» (IMC y meta sugerida dicen «Falta el peso»). Carga: esqueleto de las dos fichas y de la tabla. Error al guardar: el de la parte a §0.3 (el campo vuelve a su valor anterior). Guardado: automático.

**Edición en el lugar.** Todos los campos del perfil y de las metas que no son calculados; los que tienen vista previa (los datos del perfil que entran a RN-17, mis ajustes, el ajuste de la casa, el monto objetivo, la aportación planificada y el interruptor «Activa», que no guarda al tocar) guardan después de ella (regla general de la parte a §0.2, Q2-09 A: sexo, actividad y objetivo son listas con vista previa); en la tabla de metas, arrastrar una tarjeta entre columnas del tablero no cambia el estado (es calculado). No se editan: edad, peso actual, IMC, meta sugerida, ahorrado, avance y estado.

**Atajos.** `Alt+N` nueva meta; `↑`/`↓` y `Enter` en la tabla; `Ctrl+/` cambiar a la vez «Activa» de varias metas, con la vista previa (3) sumada sobre las seleccionadas antes de aplicar (`Enter` aplica a todas y `Esc` cancela sin cambiar ninguna); `←`/`→` mes anterior o siguiente en el calendario de metas.

**Móvil y escritorio.** 412 px: dos fichas apiladas (Axel, Ximena) con el desglose de la meta bajo el campo, luego el ajuste de la casa y las metas como tarjetas con barra de avance etiquetada. 1440 px: las dos fichas lado a lado, el ajuste de la casa debajo y la tabla de metas con su selector de vista.

## P-26.14 Inventario

**Propósito.** Decidir cómo se lleva cada producto en el inventario híbrido (decisión D): existencia por unidad (básicos y no perecederos) o ritmo de consumo (el resto), con su nivel de reposición. Regla fija: los sobrantes de recetas entran como existencia (decisión D). Necesidades N-314 y N-315.

**Vista de base de datos.** Tabla (por defecto) y tablero agrupado por modo (U1 §2.3); filtros: categoría, modo, estado y texto; orden: agotado, reponer, suficiente y luego nombre (U1 §2.5). En el tablero, arrastrar una tarjeta de «Existencia por unidad» a «Ritmo de consumo» cambia su modo después de la vista previa (parte a §0.2, punto 2).

| Campo | Tipo | Validación y mensaje | Valor por defecto | Lectura | Escritura | Necesidad |
|---|---|---|---|---|---|---|
| Producto | Texto, solo lectura (o sobrante de receta, en porciones) | — | — | `listInventoryItems` (`product_id`, `recipe_id`) y `listProducts` | — | N-314 |
| Modo | Selector: Existencia por unidad · Ritmo de consumo | Uno de los dos | Sin valor: lo elige Axel por producto (los sobrantes, existencia por unidad, decisión D) | `listInventoryItems` (`tracking_mode`) | `updateInventoryItem` (`tracking_mode`) | N-315 |
| Unidad | Texto, solo lectura: g, ml o pieza (porción en sobrantes) | — | La unidad base del producto | `listProducts` | — | N-314 |
| Nivel de reposición | Cantidad en la unidad base (hasta 3 decimales) | Vacío o ≥ 0 · «Escribe una cantidad en g, ml o piezas, o deja vacío para no avisar» | **Vacío = sin aviso** (C1 `reorder_level_milli`) | `listInventoryItems` (`reorder_level_milli`) | `updateInventoryItem` (`reorder_level_milli`) | N-315 |
| Ritmo: cantidad | Cantidad en la unidad base | Con modo «ritmo», > 0 · «Escribe cuánto se consume» | Sin valor: lo captura Axel | `listInventoryItems` (`rate_quantity_milli`) | `updateInventoryItem` (`rate_quantity_milli`) | N-315 |
| Ritmo: cada cuántos días | Entero de días | Con modo «ritmo», entero > 0 · «Escribe en cuántos días se consume esa cantidad» | Sin valor: lo captura Axel | `listInventoryItems` (`rate_period_days`) | `updateInventoryItem` (`rate_period_days`) | N-315 |
| Caducidad | Fecha opcional | — | Sin valor | `listInventoryItems` (`expires_date`) | `updateInventoryItem` (`expires_date`) | N-315 |
| Existencia, estado y días restantes | Solo lectura | — | Calculados (K7: reponer si existencia ≤ nivel; agotado si 0) | `listInventoryItems` (`quantity_milli`, `estimated_quantity_milli`, `status`, `days_left`) | — (la existencia se ajusta en P-16) | N-314 |
| Notas | Texto opcional | Vacío borra | Sin valor | `listInventoryItems` (`notes`) | `updateInventoryItem` (`notes`) | N-315 |

Acción secundaria: «Agregar producto al inventario» (`Alt+N`): `createInventoryItem` con producto, modo y existencia inicial. Con modo «ritmo» y sin cantidad o días, la fila dice «Falta el ritmo» y no muestra existencia estimada.

**Vista previa del efecto.** Al cambiar nivel o modo, antes de guardar (el modo es una lista con vista previa: muestra el bloque al elegir y guarda con `Enter` o al salir, parte a §0.2): por fila, el estado antes → después, «Con nivel n g y existencia m g: Suficiente → Reponer» (regla K7, en el navegador), y al pasar a «Ritmo de consumo» también la existencia estimada de hoy (RN-15; sin ritmo capturado, «Falta el ritmo»). Al soltar una tarjeta en la otra columna del tablero, el mismo bloque de su fila; con `Ctrl+/` sobre {n} filas (modo o nivel), «{k} pasan a Reponer, {j} a Agotado y {s} a Suficiente» y la lista de las filas que cambian; `Enter` o «Aplicar» guarda todas y `Esc` cancela sin cambiar ninguna (parte a §0.2); al cambiar el ritmo (cantidad o cada cuántos días), vista previa en línea (parte a §0.2) calculada en el navegador con RN-15: «Existencia estimada hoy: {antes} → {después}» = existencia contada − HALF_UP(cantidad del ritmo × días desde el conteo ÷ días del ritmo), nunca menor que 0, y su estado con la regla K7 (reponer si la existencia ≤ nivel de reposición; agotado si es 0), antes → después. Ejemplo de RN-15 (ritmo ilustrativo): 3 botellas de 946 ml = 2,838 ml contadas; con una botella cada 21 días, a los 30 días la existencia estimada es 2,838 − 1,351.429 = 1,486.571 ml. **Días restantes y consumo por semana** (`days_left`, `weekly_consumption_milli`, motor S4), antes → después, en el navegador con la regla de Q2-70 (a) (bloqueo `docs/estado/bloqueos/U4-vista-previa-formulas.md`, pregunta 3, RESUELTO; `docs/contratos/motores/S4-interfaz.md`): días restantes = piso(existencia estimada × días del ritmo ÷ cantidad del ritmo) y consumo por semana = HALF_UP(cantidad del ritmo × 7 ÷ días del ritmo), en milésimas de la unidad base; solo en modo «ritmo». Con el ejemplo de RN-15: piso(1,486.571 × 21 ÷ 946) = 32 días (el cociente es 32.99999…) y HALF_UP(946 × 7 ÷ 21) = 315.333 ml por semana.

**Estados.** Vacío: «Aún no hay productos en el inventario · Agregar producto». Carga: esqueleto de la tabla. Error al guardar: el de la parte a §0.3 (la celda vuelve a su valor anterior). Guardado: automático.

**Edición en el lugar.** Modo, nivel, ritmo, caducidad y notas (celda o arrastre en el tablero; el modo y el arrastre, con la vista previa antes de guardar, parte a §0.2). No se editan: producto, unidad, existencia (P-16), estado y días restantes.

**Atajos.** `Alt+N`; `Ctrl+/` cambia a la vez el modo o el nivel de las filas seleccionadas, con la vista previa sumada antes de aplicar (`Enter` aplica a todas, `Esc` cancela); `↑`/`↓` y `Enter`.

**Móvil y escritorio.** 412 px: tarjetas con producto, modo, estado (icono y texto) y nivel; el tablero se ve como dos listas plegables. 1440 px: tabla completa o tablero de dos columnas.

## Parámetros operativos de la parte b (K44, U4-12)

Filas de la lista mínima de U4-12 que viven en estas secciones (las demás, 6 a 11, van en la parte a: precio viejo y peso cero en P-26.06, sesión, intentos y ventana en P-26.01, comparador en P-26.05).

| # U4-12 | Parámetro | Sección | Clave | Inicial | Decisión |
|---|---|---|---|---|---|
| 1 | Primer aviso del presupuesto | P-26.09 | `alerts.budget_warning_pct` | 80 % | Decisión F, K44 |
| 2 | Aviso de tope | P-26.09 | `alerts.budget_limit_pct` | 100 % | Decisión F, K44 |
| 3 | Días antes del pago próximo | P-26.09 | `alerts.payment_due_days_before` | [1, 0] | Decisión F, K44 |
| 4 | `send_time` | P-26.08 (reportes) y P-26.09 (avisos con hora) | `reports.daily_send_time`, `reports.weekly_send_time`, `reports.monthly_send_time` y `alerts.send_time` (Q2-11 B, K47) | 21:30, 19:00, 08:00 y 09:00 | Q2-72 y Q2-73 (K50), K43 |
| 5 | `advance_days` | P-26.09 | `alerts.stale_prices_advance_days` (Q2-11 B, K47) | 1 día | Q2-73 (K50), K43 |
| 12 | Adjunto máximo | P-26.11 | `attachment.max_size_mb` | 10 MB | K43 |
| — | Retención de claves de idempotencia | P-26.11 | `idempotency.retention_hours` | 24 h | K6, K44 |
| — | Tiempo máximo del OCR | P-26.11 | `ocr.timeout_seconds` | 30 s | Q2-76 (K50) |
| 13 | Ajuste de la casa para bajar | P-26.13 | `health.goal_adjust_lose_kcal` | 0 kcal | K36 |
| 14 | Ajuste de la casa para subir | P-26.13 | `health.goal_adjust_gain_kcal` | 0 kcal | K36 |
| 15 | `weekday` | P-26.08 (`weekly_summary`) y P-26.09 (`meal_plan_missing`) | `reports.weekly_weekday` y `alerts.meal_plan_missing_weekday` (Q2-11 B, K47) | 7 (domingo) | Decisión F |
| 16 | `day_of_month` | P-26.08 (`monthly_report`) | `reports.monthly_day_of_month` (Q2-11 B, K47) | 1 | Decisión F, Q2-66 B |

## Campos de la puerta 1: dónde quedan

Tope por persona (K13) → P-26.12 · objetivo de salud y meta confirmada (K23) y ajuste por persona y de la casa (K36) → P-26.13 · candidatos de chat (K4) → P-26.09 · IA desactivable (decisión H) → P-26.11 · restricciones duras (K16), interruptor A/B (decisión C), promociones, tasas como se capturaron (K12) y datos de tarjeta (K18) → P-26.05 · proyección de los meses 7 a 12 (K21) → P-26.02 · mercados, membresía de Sam's y OCR local como fuente de precio → P-26.06 (parte a).

## Matriz de operaciones de la parte b (U4-02)

- P-26.08: `listNotifications`, `updateNotification`, `setNotificationRecipients`, `listNotificationRecipients`, `listJobs`, `sendTestNotification`, `generateMonthlyReport`, `listMonthlyReportMonths`, `generateExternalStatement`, `getReportSettings`, `updateReportSettings` (día y hora, Q2-11 B; C1 v1.3).
- P-26.09: `listNotifications`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `listNotificationRecipients`, `updateNotificationRecipient`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `getAlertSettings`, `updateAlertSettings`, `listPersons`; vista previa: `getBudgetMonth`, `getDashboardNextPayment`.
- P-26.10: sin operación (Local, U1 §0.3, N-301).
- P-26.11: `listBackups`, `createBackup`, `restoreBackup`, `listJobs`, `updateJob`, `previewPlannedItemsImport`, `importPlannedItems`, `getSeedStatus`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getDataSettings`, `updateDataSettings`.
- P-26.12: `listPersons`, `updatePerson`, `listCostShareRules`, `updateCostShareRule`.
- P-26.13: `listHealthProfiles`, `createHealthProfile`, `updateHealthProfile`, `getHealthSettings`, `updateHealthSettings`, `listSavingsGoals`, `createSavingsGoal`, `updateSavingsGoal`.
- P-26.14: `listInventoryItems`, `listProducts`, `updateInventoryItem`, `createInventoryItem`.

## Cobertura de criterios U4 en la parte b

U4-01: 7 secciones `## P-26.xx` con tabla de campos de 7 columnas. U4-02: matriz de arriba. U4-03: decisión F en P-26.08 y P-26.09; K4 en P-26.09; decisión G en P-26.10; decisión H en P-26.11; decisión B y K13 en P-26.12; decisión D, K23 y metas en P-26.13; decisión D en P-26.14. U4-04: valores por defecto con su cita en cada tabla. U4-05: «Vista previa del efecto» en las 7; en P-26.13, IMC y meta sugerida con cualquier dato del perfil que entra a RN-17, y en metas la aportación necesaria y el estado siguen Q2-16 A, Q2-17 C y Q2-18 A (`S2-metas`, RESUELTO), con PROVISIONAL solo la aportación de una meta vencida sin cumplir (`S2-metas-vencida`, abierto); el interruptor «Activa» de una meta y `Ctrl+/` muestran el total reservado para metas antes de guardar (R4-1). U4-06: vistas, edición en el lugar y atajos por sección. U4-07: columna «Validación y mensaje». U4-08: restaurar, activar IA y descartar chat. U4-09: «Móvil y escritorio» por sección y tokens de DESIGN. U4-10: «Estados» por sección. U4-11: SUPUESTOS y bloqueos `U4-dia-mes-reporte`, `U4-reparto-retroactivo` y `U4-vista-previa-formulas` (RESUELTOS en la puerta 2), `U4-reparto-preseleccion` (abierto, P-26.12) y `S3-iva-lectura` (abierto, parte a); las vistas previas y lecturas que dependen de un bloqueo abierto de otro artefacto lo citan: `C1-hora-reportes-job` en P-26.08, `C1-reparto-alcance` en P-26.12, `S2-metas-vencida` en P-26.13, y `S3-formula-tarjeta-pdf` y `S3-iva-lectura` en P-26.05 (parte a). U4-12: tabla de parámetros. U4-13: P-26.13.

## SUPUESTOS

1. Sección de cada clave: la que da su `@param` (00003 y 00004), que coincide con la que proponía esta parte: hora y día de los reportes en **Reportes y frecuencias**; hora, anticipación y día del plan sin armar en **Notificaciones**; retención de claves de idempotencia en **Datos y respaldos** (como U1 §2.3). `ocr.timeout_seconds` en **Datos y respaldos** es PROPUESTA de S1 que la puerta 2 no preguntó (`docs/estado/p2c/S1-puerta2.md`, «Propuestas para Axel»; la otra opción es Tiendas y precios).
2. A qué alertas aplica la hora lo fija S7; el campo aparece donde S7 lo confirme (el aviso de presupuesto salta al registrar el gasto que cruza el umbral y no necesita hora).
3. Saturación y densidad son propiedades fijas del sistema de diseño (DESIGN §1.2 P4 y P6) y no se configuran (U1: Apariencia tiene un solo control).
4. La hora del respaldo diario se edita con `updateJob` sobre la tarea del respaldo; C1 hoy liga `updateJob` solo a N-297.
5. Rangos: los aprobados en la puerta 2 (K50; parte a §0.5). No hay otros mínimos de seguridad (p. ej. no menos de 3 intentos).
6. La preselección del diálogo de alcance de P-26.12 no la fija ninguna fuente: queda **PROVISIONAL** con la opción en uso, «Solo a los gastos nuevos», y la decide Axel en la puerta 4 (bloqueo abierto `docs/estado/bloqueos/U4-reparto-preseleccion.md`). Con «solo los nuevos», con qué partes se recalcula un gasto anterior cuyo importe se edita después es PREGUNTA abierta para C1 y el motor de reparto (bloqueo abierto `docs/estado/bloqueos/S1-reparto-edicion-gasto.md`, WF-3), no de esta pantalla.
7. Las claves `reports.*` se leen y escriben con `getReportSettings`/`updateReportSettings`, la sección «Reportes y frecuencias» de C1 v1.3 (`/api/v1/settings/reports`); las `alerts.*` siguen en `getAlertSettings`/`updateAlertSettings` (P-26.09). Qué manda entre esas claves y el `schedule` de la tarea del reporte es PROVISIONAL (bloqueo abierto `docs/estado/bloqueos/C1-hora-reportes-job.md`).
8. El interruptor «Atajos de una tecla» (P-26.10) usa la necesidad de Apariencia (N-301) porque U1 §3.28 no tiene una propia; U1 lo agrega (solicitud del reporte).
