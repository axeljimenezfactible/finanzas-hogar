# Bloqueo C1-hora-reportes-job · fase C1 (v1.3, puerta 2)

Contexto: agente C1-puerta2 (P2c, 11-oct-2026). Con Q2-11 B (K47) la hora y el día de los reportes programados son claves de `app_setting` (`reports.daily_send_time`, `reports.weekly_weekday`, `reports.weekly_send_time`, `reports.monthly_day_of_month` y `reports.monthly_send_time`, sembradas por `00004_puerta2.sql`), y C1 v1.3 las lee y cambia con `getReportSettings` y `updateReportSettings`. HECHO: la tabla `job` (00001) tiene `schedule`, «expresión cron de 5 campos en hora de America/Mexico_City», y C1 (desde v1.0) deja cambiarla con `updateJob` (`JobUpdate.schedule`, N-297). Así la hora de un reporte puede quedar en dos lugares.

Qué falta: ¿qué manda la hora y el día de un reporte programado?

Opciones:
- A) Mandan las claves de `app_setting`: la API reescribe `schedule` del job del reporte en la misma transacción de `updateReportSettings`, y `updateJob` no acepta `schedule` en los jobs de reportes (solo `is_enabled`). Consecuencia: un valor en un solo lugar (la razón de Q2-11 B); cambia el texto de `updateJob` en C1 y S7 arma el cron desde las claves.
- B) Manda `job.schedule` y las claves solo informan. Consecuencia: contradice Q2-11 B (las claves dejarían de ser la fuente).
- C) Quitar `schedule` de los jobs de reportes (migración nueva de S1) y que S7 programe solo con las claves. Consecuencia: cambia el esquema.

PROPUESTA del agente: A. Sin decisión de Axel no se aplica.

Impacto: hora y día de los reportes en S7 (WF-3) y la pantalla P-26.08 (U4). Ninguna cifra cambia. Mientras tanto queda la opción en uso: `updateJob` acepta `schedule` en todos los jobs.

Cómo seguí: C1 v1.3 marca PROVISIONAL `JobUpdate.schedule` y la fila de §13, citando este archivo.

- **Estado:** ABIERTO
