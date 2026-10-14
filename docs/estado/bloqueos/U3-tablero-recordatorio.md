# Bloqueo U3-tablero-recordatorio · N y rango de «Pagos próximos» (U3 2.0, K-4)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54, palabras de Axel: «una vista rapida de recordatorio de pagos del mes, puede ser de pagos N o de creditos proximos». Axel no fijó N ni el rango. Mientras tanto rige la opción PROVISIONAL, que usa `docs/design/U3-dashboard.md` 2.0, K-4.

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **N, cuántos pagos se ven.** | (a) PROVISIONAL: **5**, como D-3 (U1 N-027 «próximos 5 pagos»), en la clave nueva de `app_setting` `dashboard.payment_reminder_count` (valor inicial 5, de 1 a 20, K44). Mientras C1 v1.5 no publique la clave, la interfaz usa 5. (b) 3, para que quepa en un vistazo a 412 px. (c) Todos los del rango. | Alto de K-4 y una petición con `limit` = N. |
| 2 | **Rango.** | (a) PROVISIONAL: primero los vencidos y luego los pendientes de hoy a hoy + 30 días. Cubre «los pagos del mes» y no se queda vacío al final del mes. (b) Solo el mes calendario en curso. (c) De hoy a hoy + 60, como D-3. | Qué pagos se ven; con la semilla el 10-oct-2026, (a) da 1 pago y (c) da 3. |
| 3 | **Pagos o créditos.** «pagos N o créditos próximos». | (a) PROVISIONAL: una fila por pago (fecha, línea, monto y estado) y varios pagos del mismo día en filas separadas. (b) Una fila por línea con su próximo pago. | Forma de la lista. |

**Qué cambia al contestar:** U3 2.0 K-4 y §6, el esquema 1.5 (valor inicial de la clave) y C1 v1.5 (`DashboardSettings`).
