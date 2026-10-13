# Reporte del orquestador (Fable 5.1) · WF-3 y WF-4

## S7-07 · Envío real por Telegram (12-oct-2026, 21:15 CST)
- **Cómo:** script `~/finanzas-hogar-orq/s7-07.sh` (fuera del repositorio): segunda instancia del servicio con `HOME` real (token legible por el usuario deepsec, K30) en `127.0.0.1:18443` sobre `datos/finanzas-dev.db`; candidato de chat de Axel anotado con su OK (K4, `ops/telegram-disparo.md` requisito 3, PROVISIONAL por `S7-convivencia-multix`) y confirmado por la API (`POST /api/v1/telegram-chat-candidates/1/confirm` → HTTP 201); sesión de Axel por la API; disparo de `sendTestNotification` para `daily_summary` (id 1) y `payment_due` (id 4).
- **Resultado:** `{'notification_code': 'daily_summary', 'status': 'sent', 'http_status': 200}` y `{'notification_code': 'payment_due', 'status': 'sent', 'http_status': 200}`; log de la instancia con 0 menciones de token, contraseña o cookie; instancia detenida.
- **Confirmación visual de Axel:** solicitada por Telegram a las 21:20 (mensaje `[Finanzas-Hogar] CONSTRUCCION TERMINADA`); pendiente de su respuesta.
- **Ximena:** pendiente de acción (K4 y `S7-convivencia-multix`).

## Decisiones PROVISIONALES del orquestador durante WF-3 (a ratificar en la puerta 4)
Listadas en `docs/estado/puerta4/PUERTA-4.md`; entre ellas: C3-matriz (P-26 Local; excepciones técnicas), C3-flujos punto 2 = (c), WF3a-bis-semilla-dev puntos 4 = (a) y 7 = sí, S6-deuda-viva-hoy = (a), RN-30 pagos de todas las líneas, scale_max_pct en C1, «Parámetros del tablero» en Ajustes › Presupuesto, K54 según AskUserQuestion.

## WF4-11 · Aviso final por Telegram
Pendiente: se registra aquí al cerrar WF-4 (respuesta `ok: true`, sin token).
