# Bloqueo WF3-creditos-salud-reportes · P-17..P-21, P-23..P-25 (12-oct-2026)

Autor: agente `front-creditos-salud-reportes`. Todo lo de aquí está implementado como **PROVISIONAL** y marcado en el código; nada se decidió en firme.

| # | Qué falta | Opciones | Impacto | Valor provisional |
|---|---|---|---|---|
| 1 | U1 habla de «parte del hogar» y «parte de Karen»; C1 reparte por `person_id` | (a) persona externa = su nombre, titular y co-titular = «Hogar»; (b) C1 marca qué personas son «hogar» | Etiquetas de partes en P-17/P-18/P-19 y quién es «Hogar» al comparar en P-20 | (a); en P-20 «Hogar» = la persona `holder` |
| 2 | N-210 «Permite desmarcar»: `DELETE /credit-payments/{id}` pide el id y `ScheduledPayment.shares` no lo trae | (a) C1 agrega `credit_payment_ids` a cada parte; (b) una lista de pagos hechos | Hoy solo se desmarca con «Deshacer» justo después de marcar | (a) propuesta, sin implementar |
| 3 | N-246 destinatarios de los próximos envíos y N-247 destinatario de cada aviso: `Job` no trae destinatarios y `recipient_id` no se puede convertir en persona con los extractos | (a) C1 agrega `recipients` (nombres) a `Job` y `recipient_person_id` a `NotificationDelivery`; (b) se quitan de U1 | P-23 no muestra destinatarios | Columna omitida |
| 4 | N-260: `ExternalStatement.payment_shares` no trae fecha | C1 agrega `due_date` | P-25 lista las partes de Karen sin fecha (con número y estado) | Sin fecha |
| 5 | N-256 «peso al inicio y al cierre» no viene en `MonthlyReport` | C1 lo agrega o U1 lo quita | P-24 no lo muestra | Omitido |
| 6 | U1 §2.6 (P-19): capital, intereses, plazo y reparto se editan con vista previa; `CreditPurchaseUpdate` solo acepta descripción, estado y nota | C1 amplía la operación o U1 se ajusta | En P-19 solo se editan comercio y nota | Solo comercio y nota |
| 7 | Rutas propias (amplía `WF3-frontend-rutas`): alta `/creditos/:id/compras/nueva`; P-20 → P-19 con `?monto&fecha&plazo&promocion&comercio`; `?credito=` hacia P-26.05 y P-20; `?mes=` en P-24; `?hoja=pago&pago=<id>`; `?hoja=generar` | Las que fije Axel | Enlaces internos | Esas |
| 8 | N-248 no dice si la casilla «Enviarlo por Telegram» viene marcada | Marcada o sin marcar | Que un reporte manual se envíe sin pedirlo | Sin marcar |
| 9 | `POST /reports/monthly` y `/external-statements` responden 202 | (a) el PDF ya está listo (trae sha256 y tamaño); (b) es asíncrono y hace falta consultar el estado | «Descargar» justo después de generar | (a) |
