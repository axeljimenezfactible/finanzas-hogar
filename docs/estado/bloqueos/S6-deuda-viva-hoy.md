# Bloqueo S6-deuda-viva-hoy · qué compras cuenta la deuda viva «hoy» y qué falta pagar del mes (K54, C1 v1.5)

Abierto el 12-oct-2026 por el agente `backend-cierre` (WF-3c bis). Lo contesta Axel en la **puerta 4**. Mientras tanto rige la opción PROVISIONAL (a) de cada punto, que eligió el orquestador (decisión 1 de WF-3c bis) para el hallazgo H-4 de `docs/estado/wf3/REVISION-c3-tablero.md` y las diferencias 1 y 2 de `docs/estado/wf3/u5-tablero.md`.

**Qué falta.** RN-21 define la deuda viva a una fecha como «la suma de los pagos programados posteriores a esa fecha de las compras no canceladas», sin decir qué pasa con una compra cuya fecha es **posterior** a la fecha de corte. La semilla trae tres despensas programadas a futuro (10-dic-2026, 10-feb-2027 y 10-abr-2027; bloqueo `S6-semilla-despensas-futuras`, opción A). C1 M-12 sí dice «de las compras hechas hasta esa fecha», y `getLiveDebt` (motor S2) las deja fuera, pero `getCreditTotals`, la deuda de cada línea y `getDebtTrend` (motor S3) las contaban: el 12-oct-2026, K-3 y P-17 decían $18,999.22 y `getLiveDebt` y P-07, $3,895.00.

**1. Deuda viva a una fecha (M-12, RN-21).**

| Opción | Qué hace | Consecuencia |
|---|---|---|
| **(a) PROVISIONAL** | Solo cuentan las compras con fecha igual o anterior al corte. Una compra programada a futuro empieza a contar el día de su fecha. | Una sola cifra hoy: $3,895.00 en K-3, P-17, P-07 y `getLiveDebt`, igual que el oráculo (`model.json: despensas[0].total`). |
| (b) | Cuentan todas las compras activas, aunque su fecha sea futura. | K-3 y P-17 dicen $18,999.22; `getLiveDebt` y M-12 tendrían que cambiar a la misma regla. |

**2. «Deuda viva proyectada» (G-5 y T-3, M-17, `getDebtTrend`).**

| Opción | Qué hace | Consecuencia |
|---|---|---|
| **(a) PROVISIONAL** | Cada cierre de mes aplica el punto 1 a su fecha: parte de la deuda de hoy menos los pagos que vencen antes del cierre, y crece en el mes de cada despensa programada. | Al cierre de oct-2026: $2,921.25 (3,895.00 − 973.75), igual que el oráculo. Al 10-abr-2027 el punto de control sigue en $6,293.44. |
| (b) | Cada cierre cuenta todas las compras activas. | Arranca en $18,025.47, que no es ninguna deuda que exista. |

**3. Pagos del mes de K-3 (M-19, `month_payment_count` y `month_pending_household_cents`).**

| Opción | Qué hace | Consecuencia |
|---|---|---|
| **(a) PROVISIONAL** | Cuántos = pagos del calendario del motor S3 con fecha en el mes en curso (las mismas líneas que `month_household_cents`). Falta pagar = por cada pago que el motor no da por pagado, la parte del hogar menos lo que ya abonaron Axel o Ximena a ese pago, nunca negativa. El abono de Karen no cubre la parte del hogar (RN-04). Un pago de corte de tarjeta, que no tiene compra, cuenta su parte del hogar completa mientras no esté pagado. | Oct-2026 en la base de desarrollo: 1 pago y $598.82 por pagar. |
| (b) | Falta pagar = la parte del hogar completa de los pagos pendientes o parciales. | Con un abono parcial del hogar se ve más de lo que falta. |

**Impacto.** Solo qué compras entran en la suma: no cambia ninguna cifra del motor S3 ni del oráculo. El código vive en una sola función, `LiveDebtAt` (`backend/internal/api/creditos/motor.go`), que usan `getCreditTotals`, `listCredits` y `getCredit` (deuda de cada línea), `getCardStatus` y `getDebtTrend`; el punto 3 vive en `pendingHousehold` (`creditos.go`). **Fuera de este cambio:** `backend/internal/api/gastos-inicio/vistas.go` (`deriveCredits`, la deuda de cada línea en las vistas de BD) sigue llamando a `motor_cred.LiveDebt` directo; está pedido como solicitud fuera de alcance en `docs/estado/wf3/backend-cierre.md`.

**Qué líneas cuentan** (corrección r1, `correccion-backend-cierre-r1`, hallazgo M-1 de `docs/estado/wf3/REVISION-backend-cierre-r1.md`): para que la cifra de hoy sea una sola también con una línea archivada que todavía tiene pagos, `getCreditTotals` y `getDebtTrend` cuentan todas las líneas, activas o archivadas, igual que `getLiveDebt`. Esa parte es PROVISIONAL y su pregunta va en el punto 5 de `S6-tablero-lineas-archivadas`.

- **Estado:** ABIERTO
