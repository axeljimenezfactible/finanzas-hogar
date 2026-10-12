# U4b-ajustes · reporte final (WF-1, 11-oct-2026)

## Hecho
- `docs/design/U4-ajustes-b.md` (58.8 KB): especificación estilo Notion de P-26.08 Reportes y frecuencias, P-26.09 Notificaciones, P-26.10 Apariencia, P-26.11 Datos y respaldos, P-26.12 Personas y reparto, P-26.13 Salud y metas de ahorro y P-26.14 Inventario. Por sección: propósito, vista de base de datos, tabla de campos de 7 columnas (campo, tipo, validación y mensaje, valor por defecto con cita, lectura, escritura, `N-xxx`), parámetros K44 con valor actual e inicial y «restablecer», acciones (críticas con el rol `color.critical.*`), vista previa del efecto, estados, edición en el lugar, atajos y móvil/escritorio. Cierra con la tabla de parámetros de la parte b (U4-12), dónde quedan los campos de la puerta 1, la matriz de operaciones (U4-02), la cobertura de criterios y SUPUESTOS.
- K36: ajuste de la casa (`health.goal_adjust_*`, 0 y 0, `getHealthSettings`/`updateHealthSettings`) y de cada persona (`goal_adjust_*` de `health_profile`, vacío = el de la casa, «Usar el de la casa» = restablecer) en P-26.13; «mantener» sin campo.
- Lo que depende de `S1-parametros-avisos` (`send_time`, `advance_days`, `weekday`, `day_of_month`) va con lugar PROVISIONAL y las dos variantes de escritura (`updateNotification` + `updateJob`, o la operación de settings `updateAlertSettings`).
- Dos bloqueos nuevos: `docs/estado/bloqueos/U4-dia-mes-reporte.md` y `docs/estado/bloqueos/U4-reparto-retroactivo.md`.

## Archivos
- `docs/design/U4-ajustes-b.md` (nuevo)
- `docs/estado/bloqueos/U4-dia-mes-reporte.md` (nuevo)
- `docs/estado/bloqueos/U4-reparto-retroactivo.md` (nuevo)
- `docs/estado/wf1/U4b-ajustes.md` (este reporte)

## Verificación
Script `python3 -I verify.py` (en el scratchpad de la sesión, fuera del repositorio) sobre `docs/design/U4-ajustes-b.md`, los extractos y 00003:
- `U4-01`: `grep -c '^## P-26\.'` = 7 (14 al unir con la parte a); 130 filas de tabla con 6 o más columnas, 0 celdas vacías.
- `U4-02`: 45 `operationId` citados, 0 inexistentes en `docs/contratos/extractos/*.yaml` (el script también marcó la palabra `operationId` de las convenciones, que no es una operación). Las 34 operaciones de Ajustes cuyas `x-necesidades` son de la parte b (N-296 a N-316 y N-323 a N-327): 0 sin campo o acción. Esquemas citados (`BackupRestoreRequest`, `ImportResult`, `RestoreResult`, `TelegramChatCandidateConfirm`): 0 inexistentes. Excepción declarada: Tema (N-301) es Local por U1 §0.3 y no tiene operación.
- `U4-03`: P-26.08 diario noche, semanal domingo, mensual día 1, estado de cuenta de Karen bajo demanda, vista web mensual: presentes. P-26.09: 4 alertas, destinatarios por aviso, chats detectados y confirmación (K4): presentes. P-26.10: tema según el sistema con interruptor, saturación apagada, densidad cómoda: presentes. P-26.11: respaldo diario local, integridad, restauración, IA apagada: presentes. P-26.12: % 50/50, qué se comparte con Karen, tope por persona sin valor (K13): presentes. P-26.13: perfil, IMC, objetivo, meta sugerida y confirmada (K23), metas con fecha y aportación: presentes. P-26.14: stock o ritmo, sobrantes como stock: presentes.
- `U4-04`: valores por defecto con cita: IA apagada (H, 00001), tema sistema (G), 50/50 (B), alertas 80/100 y [1, 0] (F, K44), adjunto 10 MB (K43), retención 24 h (K6), kcal 0 y 0 (K36), weekday 7 y day_of_month 1 (F), tope vacío (K13), objetivo vacío (K23), restricciones vacías (D). Lo no decidido dice «Sin valor: lo captura Axel» o «sin valor inicial: lo propone S7 (K31), PROVISIONAL».
- `U4-05`: «Vista previa del efecto» en las 7 secciones (frecuencia → próximo envío; umbrales → montos con `getBudgetMonth`; días antes → fechas con `getDashboardNextPayment`; reparto → ejemplo K14 con marcador; kcal → desglose; restauración e importación → antes de confirmar). El caso A/B es de P-26.05 (parte a).
- `U4-06`: atajos citados `/ir`, `/reporte`, `/respaldo`, `/tema`, `?`, `Alt+N`, `Ctrl+/`, `Ctrl+Enter`, `Ctrl+K`, `Ctrl+Shift+L`, `Enter`, `Esc`, flechas: 0 fuera de U1 §5.2/§5.3, 0 atajos nuevos (0 duplicados). Vistas: tabla en las 6 secciones con lista; tablero, calendario y galería de metas en P-26.13; tablero por modo en P-26.14.
- `U4-07`: columna «Validación y mensaje» en todas las tablas; dinero `$1,234.56`; porcentajes que suman 100 % en repartos (S1-08); fechas `10-oct-2026`.
- `U4-08`: restaurar (doble confirmación + `RESTAURAR`), activar IA (cada proveedor) y descartar chat detectado: diálogo de acción crítica con `color.critical.700`/`color.critical.300` y verbo explícito.
- `U4-09`: `grep` de hex y `rgb(` = 0; «Móvil y escritorio» (412 y 1440 px) en las 7; `size.touch.min` 44 px.
- `U4-10`: «Estados» (vacío, carga, error, guardado) en las 7; guardado automático igual en todas.
- `U4-11`: 0 fechas ISO; SUPUESTOS en el documento; bloqueos `U4-dia-mes-reporte` y `U4-reparto-retroactivo`.
- `U4-12` (filas 1–5 y 12–16 de la lista mínima; 6–11 son de la parte a): claves `alerts.budget_warning_pct` 80, `alerts.budget_limit_pct` 100, `alerts.payment_due_days_before` [1, 0], `attachment.max_size_mb` 10, `idempotency.retention_hours` 24, `health.goal_adjust_lose_kcal` 0, `health.goal_adjust_gain_kcal` 0: las 7 coinciden con el `@param` de 00003. `send_time`, `advance_days`, `weekday`, `day_of_month`: los 4 están en `@param_sin_valor` de 00003; `weekday` = 7 en `weekly_summary` y `meal_plan_missing`, `day_of_month` = 1 en `monthly_report` (sembrados en 00001, líneas 1242–1247); lugar PROVISIONAL citando `S1-parametros-avisos`. 0 mínimos de seguridad en las tablas.
- `U4-13`: desglose con los valores de prueba de S4-07: 1,748.75 × 1.2 = 2,098.50; − 300 = 1,798.50; meta sugerida entera 1,799 (un solo ROUND_HALF_UP al final, `suggested_kcal_target` es entero). Coincide con S4-07 (1,798.50 antes de redondear). Operaciones: `listHealthProfiles`, `createHealthProfile`, `updateHealthProfile` (N-313, N-327), `getHealthSettings`, `updateHealthSettings` (N-326).

## Supuestos
1. Lugar de los parámetros que K44 no ubica: `send_time`, `weekday` y `day_of_month` de los reportes en Reportes y frecuencias; `send_time` y `advance_days` de las alertas y `weekday` del plan sin armar en Notificaciones; retención de idempotencia en Datos y respaldos (como U1). PROVISIONAL hasta la respuesta de `S1-parametros-avisos`.
2. A qué alertas aplica la hora lo confirma S7; el aviso de presupuesto salta al registrar el gasto y no necesita hora.
3. Saturación y densidad son fijas del sistema de diseño y no se configuran (U1: Apariencia tiene un solo control).
4. La hora del respaldo diario se edita con `updateJob` sobre la tarea del respaldo (C1 hoy liga `updateJob` solo a N-297).
5. Detalles de interfaz sin decisión (casilla «Enviarlo también por Telegram» sin preselección; la persona elige el mes) no son valores operativos ni cifras.

## Bloqueos
- `docs/estado/bloqueos/U4-dia-mes-reporte.md` (nuevo): días 29 a 31 del reporte mensual.
- `docs/estado/bloqueos/U4-reparto-retroactivo.md` (nuevo): si un cambio de reparto alcanza gastos ya capturados.
- `docs/estado/bloqueos/S1-parametros-avisos.md` (abierto, citado; no se decide en WF-1).

## Pendientes
- Integrador: unir `docs/design/U4-ajustes-b.md` con `docs/design/U4-ajustes-a.md` en `docs/design/U4-ajustes.md` y unificar las convenciones (encabezado `## P-26.xx`, tabla de 7 columnas, SUPUESTOS único). El total de U4-01/U4-02 sobre las 14 secciones se mide al unir.
- Cuando se responda `S1-parametros-avisos`: quitar la variante que no aplique en P-26.08 y P-26.09 y la marca PROVISIONAL de lugar.
- Cuando se aprueben los valores de K31 (puerta 4): llenar «Valor inicial» de `send_time`, `advance_days` y hora del respaldo, y activar su «Restablecer».

## Solicitudes fuera de alcance
1. `docs/contratos/openapi.yaml` (vía `backend/internal/api/c1gen/build.py`), `updateJob` y `listJobs`: agregar N-302 a `x-necesidades` de `updateJob` (hora del respaldo diario) y aclarar qué operación escribe la hora de un reporte (`updateNotification.params_json.send_time` o `updateJob.schedule`) para que no haya dos escritores del mismo dato. Motivo: hoy las dos lo permiten (N-297).
2. S1 (`backend/db/migrations/`, siguiente migración) y C1 (`DataSettings`): claves `@param_sin_valor` explícitas para la hora y las versiones que se guardan del respaldo diario (K31 punto 3) y para el tamaño máximo del PDF (K31 punto 6), expuestas por `getDataSettings`/`updateDataSettings` (y la de PDF donde decida C1). Motivo: K44 pide que todo parámetro operativo tenga lugar en Ajustes; hoy no tienen clave ni operación y U4 no puede mostrarlos como campos.
3. C1, si la respuesta a `S1-parametros-avisos` es A o C: exponer el valor inicial de `weekday` (7) y `day_of_month` (1) en `Notification` (como `initial_value` de las operaciones de settings) para que «Restablecer» no dependa de un valor escrito en la pantalla.
4. `docs/design/DESIGN.md` §7 (U2b-2): componentes que esta parte usa: fila de parámetro con «Valor inicial» y «Restablecer al valor inicial», fichas de lista de días, control segmentado de tema, diálogo de acción crítica con texto escrito (`RESTAURAR`).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- K31 punto 1 (hora de reportes): resumen diario 21:30, semanal domingo 19:00, mensual día 1 a las 08:00. Por qué: de noche ya se capturó el gasto del día; el domingo a esa hora se planea la semana; el día 1 en la mañana el mes anterior ya cerró.
- K31 punto 2 (alertas): `advance_days` de precios viejos = 1 (un día antes de la despensa da tiempo de actualizar precios); hora de avisos con hora 09:00 (pago próximo en horario para pagar) y plan sin armar domingo 18:00 (antes del resumen semanal).
- K31 punto 3 (respaldo diario): 03:00 (la máquina está encendida y nadie la usa) y conservar 30 copias diarias más la primera de cada mes por 12 meses. Por qué: la base es un solo archivo pequeño en disco local; da un mes para deshacer día por día y un año por mes.
- K31 punto 6 (PDF): 10 MB, igual que el adjunto, muy por debajo del límite de envío de documentos de los bots de Telegram (50 MB).
- K31 punto 8 (permisos): base y respaldos 0600 y sus carpetas 0700, del usuario del servicio; no se muestra en Ajustes (es del sistema operativo).
- `U4-dia-mes-reporte`: opción B (en meses cortos sale el último día), para que nunca falte el reporte.
- `U4-reparto-retroactivo`: opción C (preguntar el alcance, como K2), para no cambiar sin querer estados de cuenta ya entregados a Karen.

## Consumo aproximado
38 acciones de herramienta; ~215 KB leídos; ~75 KB escritos (documento 58.8 KB, bloqueos 2.9 KB, reportes ~13 KB).
