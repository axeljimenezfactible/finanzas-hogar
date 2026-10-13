# Bloqueo S6-tablero-lineas-archivadas · ¿los pagos de una línea de crédito archivada cuentan en el tablero? (RN-30, M-20, M-21; C1 v1.5)

Abierto el 12-oct-2026 por el agente `correccion-tablero-backend-r1` (WF-3c), por el hallazgo M-2 de `docs/estado/wf3/REVISION-tablero-backend-r1.md`. Lo contesta Axel en la **puerta 4**. La propuesta PROVISIONAL es la opción (a) del punto 1 y la (a1) del punto 2. **Rige en el código desde WF-3c bis** (decisión 2 del orquestador, salida (i) del punto 4; la aplicó `backend-cierre`): `PaymentsByMonth` recorre todas las líneas y `has_credits` sigue la (a1). Hasta la ronda 2 el servidor hacía lo de la opción (b). Actualizado el 12-oct-2026 por `correccion-tablero-backend-r2` (hallazgo M-2 de `REVISION-tablero-backend-r2.md`).

**Estado del código (12-oct-2026, WF-3c bis, `backend-cierre`).**
- Aplicada la solicitud 1 de `correccion-tablero-backend-r2.md` tal cual en `backend/internal/api/creditos/tablero.go`: el calendario del motor S3 se arma con todas las líneas (`s.state(c)` de cada `credit`), `has_credits` = alguna línea activa o algún pago del calendario entre `from` y `to`.
- Quitado el `t.Skip` del subtest `opcion_a_PROVISIONAL` de `TestOp_Dashboard_LineaArchivada`: pasa (RN-30 «de todas las líneas»).
- Aplicada la solicitud 2: el texto de `has_credits` en C1 1.5.1 cita este bloqueo.
- La pregunta que se abrió aquí sobre `getCreditTotals` pasó al **punto 5** (corrección r1 de `backend-cierre`, 12-oct-2026).

**Estado del código (12-oct-2026, ronda 2).**
- La opción PROVISIONAL **no está aplicada**. Vive en `backend/internal/api/creditos/tablero.go` (`PaymentsByMonth`). Ese archivo es del paquete de `api-creditos` (CONVENCIONES-backend §2) y no está en el alcance de `api-tablero` (`runs/wf3c-tablero.js`).
- Hasta que se aplique el cambio, el tablero cuenta solo las líneas activas (`s.states(nil)` filtra `credit.is_active`) y `has_credits` sale `true` fijo cuando hay alguna línea activa. Es lo que hace la opción (b), pero nadie la eligió.
- La prueba `TestOp_Dashboard_LineaArchivada` (`backend/internal/api/tablero/tablero_datos_test.go`) tiene dos partes:
  - **Corre siempre** lo que vale con cualquier opción: el efectivo de S2 no cambia al archivar (700,000), G-1/G-6 y el barómetro dan la misma cifra del mes, y con `has_credits` false los pagos y el porcentaje son 0. Esta parte descarta la (a2).
  - El subtest **`opcion_a_PROVISIONAL`** comprueba la opción (a) y la (a1). Está en SKIP con la referencia a este punto.
- Se comprobó con `go test -overlay`:
  - con el cambio de la solicitud y sin el SKIP, la prueba pasa, y `api`, `api/creditos` y `api/tablero` siguen OK;
  - sin el cambio y sin el SKIP, falla con «pagos 0 (antes 650000…)».

**Qué falta.** RN-30 (PROVISIONAL) dice «pagos = parte del hogar de los pagos programados con fecha en el mes, **de todas las líneas**, pagados o no (la misma cifra que `month_household_cents` de `getCreditTotals`)». Las dos mitades chocan cuando hay una línea archivada: `getCreditTotals` (S3) cuenta solo las activas, y el efectivo de S2 (`dominio.LoadFinEngine`) cuenta todos los `scheduled_payment`. Con el fixture del tablero y la línea archivada, octubre daba `credit_payments_cents` 650,000 → **0**, `debt_ratio_pct` 33 → **0** y `has_credits` true → **false**, mientras `cash_cents` (S2) se quedaba en **700,000**.

**Punto 1 · qué líneas cuentan en los pagos del tablero (B-1, G-1, G-6).**

| Opción | Qué hace | Ventaja | Riesgo |
|---|---|---|---|
| **(a) PROVISIONAL** | Todas las líneas, activas o archivadas. | Es el texto literal de RN-30 y M-20/M-21. Coincide con el efectivo de S2: al archivar una tarjeta liquidada, su historia sigue en G-1, G-6 y en el barómetro de meses pasados. | Con una línea archivada que tenga un pago con fecha en el mes en curso, el tablero deja de ser igual a `getCreditTotals.month_household_cents` (ver punto 3). |
| (b) | Solo las líneas activas (lo que hace hoy el código). | Igual a `getCreditTotals`. | Contradice «de todas las líneas»: la historia de una línea archivada desaparece del tablero pero no del efectivo de S2. Pide corregir RN-30 (`docs/contratos/reglas-negocio.md`) y M-20/M-21 (`docs/contratos/C1-datos.md` y el texto de `payments_cents` por `c1gen`). |

**Punto 2 · `has_credits` con la opción (a).** C1 dice «true si hay al menos una línea de crédito activa; con false, `payments_cents` = 0 y `ratio_pct` = 0 (cifra real)». Con (a) y todas las líneas archivadas, un mes pasado puede tener pagos con `has_credits` false.

| Opción | Qué hace | Riesgo |
|---|---|---|
| **(a1) PROVISIONAL** | true si hay al menos una línea activa **o** algún pago (de una línea archivada) con fecha en el mes. Se cumple siempre «con false, pagos = 0». | Pide afinar el texto de `has_credits` en C1 (`c1gen`): «…activa, o algún pago de una línea archivada en el mes». |
| (a2) | true solo si hay una línea activa (texto literal de la primera mitad). | Puede salir false con pagos > 0: la PWA mostraría «sin créditos» junto a un porcentaje. |

**Punto 3 · igualdad con `getCreditTotals`.** Con (a), el tablero y `getCreditTotals` solo diferían si una línea archivada tenía un pago con fecha en el mes en curso. Desde la corrección r1 de `backend-cierre` ya no difieren: ver el punto 5.

**Punto 4 · Decisión pendiente del orquestador: el alcance (no la regla).** Las rondas 1 y 2 de corrección son de `api-tablero` y no pueden tocar `api/creditos/` ni `c1gen`. M-2 solo se cierra de una de estas dos maneras:

| Salida | Qué se hace | Quién | Efecto en la propuesta |
|---|---|---|---|
| **(i)** | Aplicar la solicitud 1 de `docs/estado/wf3/correccion-tablero-backend-r2.md` (unas 10 líneas en `PaymentsByMonth`) y quitar el SKIP del subtest. Aplicar también la solicitud 2 (texto de `has_credits` por `c1gen` y regenerar los extractos). | El orquestador amplía el alcance a `backend/internal/api/creditos/tablero.go` y `c1gen`, o lo hace el dueño de `api/creditos` y de `c1gen`. | Ninguno: rige la (a) con la (a1), como se propone. |
| (ii) | Que los artefactos digan lo mismo que el código: este bloqueo declara que hasta la puerta 4 rige la (b), el subtest comprueba la (b) y se pide corregir RN-30 y M-20/M-21. | El orquestador. | Cambia la opción PROVISIONAL propuesta. Por eso no la toma un agente de corrección. |

PROPUESTA: la (i). Es un cambio chico, ya probado con overlay en dos rondas, y no mueve ninguna cifra del motor.

**Punto 5 · deuda viva y pagos del mes de K-3, P-17 y G-5 (`getCreditTotals` y `getDebtTrend`).** Abierto el 12-oct-2026 por `correccion-backend-cierre-r1`, por el hallazgo M-1 de `docs/estado/wf3/REVISION-backend-cierre-r1.md`. Hasta esa corrección, `getCreditTotals` y `getDebtTrend` contaban solo las líneas activas (`states` en `backend/internal/api/creditos/motor.go` filtraba `credit.is_active`), mientras `getLiveDebt` (motor S2) y el barómetro cuentan todas. Archivar una línea con deuda está permitido (`updateCredit` no lo impide). Con una línea archivada con pagos pendientes, la prueba del revisor daba: deuda viva de K-3 11,250.00 contra `getLiveDebt` hoy 13,250.00; G-5 al cierre de octubre 3,750.00 contra `getLiveDebt` del 31-oct 4,750.00; y pagos del hogar del mes de K-3 6,500.00 contra el barómetro 7,500.00.

| Opción | Qué hace | Consecuencia |
|---|---|---|
| **(a) PROVISIONAL, aplicada en el código** | `getCreditTotals` (la deuda viva, los pagos del mes y de la quincena, cuántos, lo que falta pagar y los intereses) y `getDebtTrend` sin `credit_id` cuentan todas las líneas, activas o archivadas. Con `credit_id`, cuentan las líneas nombradas aunque estén archivadas. | Una sola cifra: la deuda viva de K-3 y P-17 es la de `getLiveDebt` hoy, cada cierre de G-5 es la de `getLiveDebt` ese día y los pagos del hogar de K-3 son los del barómetro, como dicen RN-21 (FIRME, «de las compras no canceladas»), RN-30, U3 K-3 («de todas las líneas»), la decisión 1 del orquestador de WF-3c bis y el texto de `DashboardDebtCapacity.payments_cents` en C1. Archivar una línea no cambia ninguna cifra mientras le queden pagos. |
| (b) | Solo las líneas activas, y la tarjeta K-3 y P-17 dicen que no cuentan las archivadas. | K-3, P-17 y G-5 dejan de coincidir con `getLiveDebt`, P-07 y el barómetro cuando hay una línea archivada con pagos pendientes. Pide corregir RN-21, RN-30, el texto de `CreditTotals` y `DebtTrend` en C1 y U3 K-3. |

Por qué la (a): es la única que cumple a la vez la decisión 1 («la misma semántica de `getLiveDebt?date=hoy`») y la decisión 2 («de TODAS las líneas, como dice C1») del orquestador sin cambiar ningún texto FIRME. Ninguna cifra del motor cambia (solo qué líneas se le pasan) y la base de desarrollo no tiene líneas archivadas, así que sus cifras no se mueven (389500 hoy). Código: `states` en `backend/internal/api/creditos/motor.go`; prueba `TestOp_getCreditTotals_LineaArchivada` en `creditos_test.go` (archivar no cambia `getCreditTotals` ni `getDebtTrend`; deuda viva = `getLiveDebt` hoy y en los 12 cierres; pagos del mes = barómetro). C1 1.5.1 lo marca en `CreditTotals.live_debt_cents`, `CreditTotals.month_total_cents`, `DebtTrend.months[].live_debt_cents` y M-17/M-19.

**Pregunta para Axel (punto 5).** Si archivas una línea que todavía tiene pagos por hacer, ¿quieres que su deuda y sus pagos sigan sumando en «Deuda viva», en «Créditos del mes» y en la gráfica de deuda proyectada (a), o que salgan aunque falte pagarlos (b)?

**Impacto.** Ninguna cifra del motor cambia: S3 (`motor_cred.UnifiedCalendar`) entrega los mismos pagos; solo cambia qué líneas se le pasan. Afecta `payments_cents`, `payments_total_cents`, `payment_count`, `ratio_pct`, `status`, `headroom_cents` y `has_credits` de `getDashboardDebtCapacity`, y `credit_payments_*`, `debt_ratio_pct` y `debt_status` de `getDashboardMonthlyOverview`, solo cuando hay líneas archivadas con pagos en la ventana. Sin migración.

**Pregunta para Axel.** Cuando archives una tarjeta que ya terminaste de pagar, ¿quieres que sus pagos de meses pasados sigan apareciendo en el barómetro y en las gráficas del tablero (a) o que desaparezcan (b)?
