# Bloqueo S6-respaldos · hora, retención y restauración (PROVISIONAL)

12-oct-2026 · agente auth-sistema (WF-3a) · K31 (puerta 4), S6-09.

1. **Hora y retención** (PROPUESTA, la misma de backend-base en `cmd/finanzas`): respaldo diario a las **03:15** de México (casa dormida, después del resumen diario de las 21:30 y antes del reporte mensual de las 08:00) y se conservan **30** diarios (un mes para volver atrás; cada copia pesa ~1 MB vacía). Los manuales y los previos a una restauración (`pre_migration`) **no se borran solos**. Constantes `respaldo.DailySchedule` y `respaldo.KeepScheduled`. Debe copiarse a `ops/runbook.md` (no existe; es de WF-4: solicitud fuera de alcance).
2. **Qué conserva una restauración** (PREGUNTA): `restoreBackup` reemplaza toda la base con el respaldo, salvo dos tablas que se reponen del estado previo: `backup` (para no perder el catálogo ni el respaldo previo que se acaba de tomar) y `device_session` (una sesión cerrada después del respaldo no revive y la sesión de quien restaura sigue abierta). Opciones: a) así; b) restaurar todo tal cual; c) conservar también `user_account` (si no, una contraseña cambiada después del respaldo vuelve a la anterior) y `login_attempt`. Impacto: lista `respaldo.PreservedTables`.
3. **Síncrono**: `createBackup` y `restoreBackup` hacen el trabajo dentro de la petición y responden 202 con el resultado ya listo (la base pesa ~1 MB). Si crece, pasar a job.
4. Un respaldo en curso (`status = running`) da 409 a otro manual.
