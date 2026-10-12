# Bloqueo S7-estado-cuenta-fechas · fechas de pago del estado de cuenta en C1

Estado: ABIERTO · 12-oct-2026 · agente s7-pdf (WF-3a).

## Qué falta
`ExternalStatement.payment_shares[]` (C1 1.4) trae `person_id`, `amount_cents`, `paid_cents` y `status`, pero **no la fecha de pago ni el pago programado**. S7-03 exige las 4 fechas (30-oct, 15-nov, 30-nov y 15-dic-2026) en el PDF y que crucen con la API.

## Qué se hizo (PROVISIONAL)
El PDF toma cada fecha del mismo renglón `scheduled_payment` y la cita desde `getMonthlyReport` del mes de la fecha (`external_participation[k].payment_shares[j].due_date`), que sí la trae: el cruce de S7-03 da 0 diferencias con 20 de 20 cifras y fechas.

## Opciones
- **A (recomendada)**: C1 1.5 agrega `scheduled_payment_id` y `due_date` a `payment_shares` (igual que `PersonParticipation.payment_shares`). Impacto: regenerar `types_gen.go`, 2 líneas en `datos.go` y la cita en `vistas.go`.
- **B**: quedarse con la cita cruzada a `getMonthlyReport`.
