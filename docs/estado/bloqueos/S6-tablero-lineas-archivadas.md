# Bloqueo S6-tablero-lineas-archivadas · ¿los pagos de una línea de crédito archivada cuentan en el tablero? (RN-30, M-20, M-21; C1 v1.5)

Abierto el 12-oct-2026 por el agente `correccion-tablero-backend-r1` (WF-3c), por el hallazgo M-2 de `docs/estado/wf3/REVISION-tablero-backend-r1.md`. Lo contesta Axel en la **puerta 4**. La propuesta PROVISIONAL es la opción (a) del punto 1 y la (a1) del punto 2. **Todavía no rige en el código**: hoy el servidor hace lo de la opción (b). Eso no es una elección: falta un cambio fuera del alcance de `api-tablero` (punto 4). Actualizado el 12-oct-2026 por `correccion-tablero-backend-r2` (hallazgo M-2 de `REVISION-tablero-backend-r2.md`).

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

**Punto 3 · igualdad con `getCreditTotals`.** Con (a), el tablero y `getCreditTotals` solo difieren si una línea archivada tiene un pago con fecha en el mes en curso (raro: se archiva una línea liquidada). PROVISIONAL: `getCreditTotals` no cambia (es de `api-creditos` y de S3). Si Axel quiere que coincidan siempre, el cambio es de `api/creditos/creditos.go` (fuera de este alcance).

**Punto 4 · Decisión pendiente del orquestador: el alcance (no la regla).** Las rondas 1 y 2 de corrección son de `api-tablero` y no pueden tocar `api/creditos/` ni `c1gen`. M-2 solo se cierra de una de estas dos maneras:

| Salida | Qué se hace | Quién | Efecto en la propuesta |
|---|---|---|---|
| **(i)** | Aplicar la solicitud 1 de `docs/estado/wf3/correccion-tablero-backend-r2.md` (unas 10 líneas en `PaymentsByMonth`) y quitar el SKIP del subtest. Aplicar también la solicitud 2 (texto de `has_credits` por `c1gen` y regenerar los extractos). | El orquestador amplía el alcance a `backend/internal/api/creditos/tablero.go` y `c1gen`, o lo hace el dueño de `api/creditos` y de `c1gen`. | Ninguno: rige la (a) con la (a1), como se propone. |
| (ii) | Que los artefactos digan lo mismo que el código: este bloqueo declara que hasta la puerta 4 rige la (b), el subtest comprueba la (b) y se pide corregir RN-30 y M-20/M-21. | El orquestador. | Cambia la opción PROVISIONAL propuesta. Por eso no la toma un agente de corrección. |

PROPUESTA: la (i). Es un cambio chico, ya probado con overlay en dos rondas, y no mueve ninguna cifra del motor.

**Impacto.** Ninguna cifra del motor cambia: S3 (`motor_cred.UnifiedCalendar`) entrega los mismos pagos; solo cambia qué líneas se le pasan. Afecta `payments_cents`, `payments_total_cents`, `payment_count`, `ratio_pct`, `status`, `headroom_cents` y `has_credits` de `getDashboardDebtCapacity`, y `credit_payments_*`, `debt_ratio_pct` y `debt_status` de `getDashboardMonthlyOverview`, solo cuando hay líneas archivadas con pagos en la ventana. Sin migración.

**Pregunta para Axel.** Cuando archives una tarjeta que ya terminaste de pagar, ¿quieres que sus pagos de meses pasados sigan apareciendo en el barómetro y en las gráficas del tablero (a) o que desaparezcan (b)?
