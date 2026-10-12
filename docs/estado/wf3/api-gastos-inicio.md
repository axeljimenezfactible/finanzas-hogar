# api-gastos-inicio · WF-3a

12-oct-2026 · Opus · **cerrado por tope**: hay 13 de 26 operaciones implementadas y probadas; las otras 13 siguen en 501 (ver Pendientes).

## Hecho
- Paquete `backend/internal/api/gastos-inicio/` (`package gastos`, `Register`). Implementa las 13 operaciones: `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary`, `getExpense`, `updateExpense`, `deleteExpense`, `listExpenseItems`, `duplicateExpense`, `markExpensePaid`, `listTicketScans`, `getTicketScan` y `saveTicketScanDraft`.
- Repositorio en `internal/dominio/gastos-inicio_gastos.go`:
  - Gastos, renglones y partes congeladas en `expense_share`.
  - Reparto RN-03: regla del renglón → del gasto → de la categoría → por defecto → «Hogar» implícita. Se agrupa por regla, porcentajes con HALF_UP, fijos tal cual y el resto a la persona del resto. El reparto propio y el renglón «mixto» crean reglas de un solo uso.
  - Los 6 filtros globales más `kind`, `payment_method` y `has_attachment`.
  - `q` busca sin mayúsculas ni acentos en descripción, notas, categoría y tienda.
  - `family` es incluyente (opción A de C1-filtro-familia).
- Repositorio en `internal/dominio/gastos-inicio_tickets.go`: cola de tickets y borrador K3 (RN-23).
- K9 y RN-24: `expense.idempotency_key` y `captured_at`. Si llega otra vez la misma clave, aunque hayan pasado más de 24 h, la respuesta es 200 con el mismo gasto e `Idempotent-Replayed: true`. Con otro cuerpo dentro de las 24 h, 409 (middleware).
- Cada alta, cambio o baja escribe en `audit_log` (`insert`, `update`, `delete`).
- Motores usados, solo con llamadas:
  - `motor_fin.Averages` para el promedio diario y el semanal (RN-27 y K39: 973.75 en 31 días dan 31.41 y 219.88, probado).
  - `motor_fin.MonthlyAverage`.
  - `precios/ocr.ParseLine` para los renglones que propone `getTicketScan`.
  - Todas las cifras van en centavos enteros.

## Archivos
- `backend/internal/api/gastos-inicio/gastos.go`, `resumen.go`, `tickets.go` y `gastos_test.go`
- `backend/internal/dominio/gastos-inicio_gastos.go`, `gastos-inicio_tickets.go` y `gastos-inicio_gastos_test.go`
- `docs/estado/bloqueos/S6-gastos-inicio.md` (8 puntos PROVISIONAL o PREGUNTA)
- `docs/estado/wf3/api-gastos-inicio.md`

## Verificación
- `gofmt -l` sobre mis archivos: vacío. `go vet ./internal/api/gastos-inicio/ ./internal/dominio/`: sin errores.
- `go test ./internal/api/gastos-inicio/... ./internal/dominio/... -count=1 -cover`:
  - `ok … gastos-inicio` con 85.0 % de cobertura.
  - `ok … dominio` con 35.7 % del paquete completo. Mis archivos `gastos-inicio_*` tienen 293 de 339 sentencias cubiertas (86.4 %).
  - El total del paquete baja porque otros agentes están agregando archivos en paralelo (`creditos_*`, `presupuesto_*`, `compras-*`, `comidas-salud_*`, `ajustes_repos.go`) que todavía no tienen pruebas. backend-base reportó 87.3 % antes.
- `go build ./...` y `make verify-wf3-go`: **FALLAN, pero no por este agente.** `internal/api/comidas` y `internal/api/comprasa` están a medio escribir y tienen métodos sin definir. `gofmt -l .` también marca archivos de ajustes, comidas, comprasa, creditos y presupuesto. Mis paquetes compilan y pasan.
- **S6-03**: 13 de 26 operaciones de esta área tienen prueba `TestOp_<operationId>`. Cada prueba valida la respuesta contra `openapi.yaml` con `testutil` (`c.Do`) y cubre 2xx y al menos un 4xx. Las 13 restantes responden 501 (tabla abajo).
- **S6-04**:
  - Errores con el esquema único y `request_id`.
  - Paginación de 1 a 200 (`limit=201` y `limit=0` dan 400) con cursor opaco.
  - Fechas AAAA-MM-DD y los instantes en −06:00.
  - Dinero en centavos enteros (`"20.01"` da 400).
  - Idempotency-Key obligatoria (sin clave, 400) y K9 probado pasadas 25 h.
- **S6-15**: los identificadores van en inglés. Los mensajes al usuario están en español de México y dicen qué hacer. No se registran ni contraseñas ni tokens.

| operationId | prueba | 2xx | 4xx probados |
|---|---|---|---|
| createExpense | TestOp_createExpense | 201, 200 (K9) | 400, 401, 409, 422 |
| listExpenses | TestOp_listExpenses | 200 | 400 |
| getExpenseTotals | TestOp_getExpenseTotals | 200 | 400 |
| getExpenseSummary | TestOp_getExpenseSummary | 200 | 400/422 |
| getExpense | TestOp_getExpense | 200 | 400, 404 |
| updateExpense | TestOp_updateExpense | 200 | 400, 404, 422 |
| deleteExpense | TestOp_deleteExpense | 204 | 404 (el 409 se prueba en dominio) |
| listExpenseItems | TestOp_listExpenseItems | 200 (con «Ajuste del ticket») | 400, 404 |
| duplicateExpense | TestOp_duplicateExpense | 201 | 400, 404 |
| markExpensePaid | TestOp_markExpensePaid | 200 | 404, 409, 422 |
| listTicketScans | TestOp_listTicketScans | 200 | 400 |
| getTicketScan | TestOp_getTicketScan | 200 | 404 |
| saveTicketScanDraft | TestOp_saveTicketScanDraft | 200 | 404, 409, 422 |
| createTicketScan, confirmTicketScan, getExpenseMonthlyTrend, search, listCalendarEvents, queryView, getDashboardMonthVsPlan, getDashboardFortnightFlow, getDashboardNextPayment, getDashboardWeekMeals, getDashboardBestBasket, getDashboardSavings, getPersonParticipation | — (501) | — | — |

## Supuestos
Todos están en el bloqueo, cada uno con su opción PROVISIONAL:
- Reparto RN-03 en dominio.
- Confianza OCR `low`.
- Promedio mensual sobre los meses que toca el periodo.
- Periodo anterior equivalente.
- `accruals` omitido y `cash_flows` solo de contado.
- 409 en gastos que vienen de crédito.
- `category_id` exacto.
- Directorio `gastos-inicio` con `package gastos`.

Además: si no hay regla por defecto, el 100 % va al titular (RN-03, «Hogar»). Un ticket del borrador es `kind = ticket_photo`. La cola de tickets son los que tienen `expense_id` NULL y `ocr_status` pending, done o failed.

## Bloqueos
- `docs/estado/bloqueos/S6-gastos-inicio.md`, puntos 1 a 8.

## Pendientes (13 operaciones, siguen en 501)
- `createTicketScan`: falta el multipart, guardar el archivo en la carpeta de adjuntos de auth-sistema (convención compartida con `uploadExpenseAttachment`), un job y `ocr.Engine.Recognize`.
- `confirmTicketScan`: una transacción que crea gasto, precios `ticket_ocr`, inventario con `inventory_movement` y, si aplica, la compra a crédito con `motor_cred`.
- `getExpenseMonthlyTrend` y las tarjetas `getDashboardMonthVsPlan`, `getDashboardFortnightFlow` y `getDashboardSavings`: necesitan armar el `motor_fin.Input` desde la base (bloqueo, punto 5) y llamar a `PlanVsActual`, `FortnightFlow` o `SavingsGoals`.
- `getDashboardNextPayment` y `getPersonParticipation`: `motor_cred` y las tablas de crédito.
- `getDashboardWeekMeals` y `getDashboardBestBasket`: recetas, precios y la lista semanal.
- `search`, `listCalendarEvents` y `queryView` (C1 §8).
- `created_by_user_account_id` y `updated_by_user_account_id` de Expense, leídos de `audit_log` (son opcionales y no se mandan).
- Borrado: `datos/tmp/agi` se borra al cerrar.

## Propuestas para Axel
- Exportar `motor_fin.SplitRule` para que el reparto de gastos salga del motor S2 (bloqueo, punto 1).
- Guardar la confianza por renglón del OCR con una migración nueva (bloqueo, punto 2).
- Fijar el código de la categoría de papel para poder dar el devengo por gasto (bloqueo, punto 5).

## Paquetes instalados
Ninguno.

## Consumo aproximado
57 acciones, unos 180 KB leídos y unos 100 KB escritos.

## Continuación 1

12-oct-2026 · Opus · **cerrado por tope**: se agregan 6 operaciones, así que van 19 de 26 implementadas y probadas. Siguen 7 en 501 (ver Pendientes de la continuación).

### Hecho
Las 6 operaciones nuevas están en `inicio.go` y se registran en `Register`.
- `getDashboardMonthVsPlan` (U3 tarjeta 1, M-01):
  - real y plan salen de `motor_fin.PlanVsActual(in, mes, 1, "")`;
  - % con `dominio.PercentOf` (HALF_UP); `status` según los umbrales K44 que lee el motor;
  - sin plan, los campos opcionales van ausentes; `days_left` cuenta hoy.
- `getDashboardFortnightFlow` (M-02): todas las cifras salen de `motor_fin.FortnightFlow(in, quincena, 1)`. El saldo al cierre es el del motor, y la prueba comprueba la identidad ingresos − hechos − comprometido.
- `getDashboardSavings` (M-07):
  - totales del mes con `motor_fin.SavingsGoals`;
  - `next_goal` = la meta activa con la fecha meta más próxima, con su avance del motor.
- `getExpenseMonthlyTrend` (N-044):
  - parte del hogar devengada, gasto por gasto, con `fe.Only` + `motor_fin.Periods` en los 12 meses hasta el mes en curso;
  - filtros de persona, familia y crédito; las 6 categorías principales más «Otras»;
  - empieza en el primer mes con datos.
- `getPersonParticipation` (N-026, M-08):
  - 404 si la persona no es externa;
  - despensa devengada en el mes con el motor S2, igual que `getAccruedGroceries`;
  - sus partes de pagos del mes, pagadas o pendientes (RN-04), su deuda viva y sus compras particulares vivas.
- `search` (N-002):
  - sin mayúsculas ni acentos;
  - 9 tipos, cada uno con su total y máximo 5 resultados;
  - texto vacío = grupos vacíos; `type` fuera del enumerado = 400.
- Repositorio nuevo `internal/dominio/gastos-inicio_inicio.go` con `GastosExpenseIDs`, `GastosIsExternalPerson`, `GastosPersonShares`, `GastosParticularPurchases`, `GastosNextSavingsGoal` y `GastosSearch`. Usa la entrada del motor S2 que ya arma `dominio.LoadFinEngine` (api-presupuesto), así que el bloqueo §5 queda resuelto para las tarjetas.

### Archivos
- Nuevos:
  - `backend/internal/api/gastos-inicio/inicio.go` e `inicio_test.go`;
  - `backend/internal/dominio/gastos-inicio_inicio.go` y `gastos-inicio_inicio_test.go`.
- Cambiados:
  - `backend/internal/api/gastos-inicio/gastos.go`: 6 líneas en `Register`;
  - `docs/estado/bloqueos/S6-gastos-inicio.md`: §9 a §12;
  - este reporte.

### Verificación
- `go build ./...` y `go vet ./...` sin errores (los corre `make verify-wf3-go`). `gofmt -l` sobre mis archivos: vacío.
- `go test ./internal/api/gastos-inicio/... ./internal/dominio/... -count=1 -cover`:
  - `gastos-inicio`: **86.2 %**, con 19 pruebas `TestOp_*` en PASS;
  - `dominio`: **87.6 %**.
- `make verify-wf3-go` **falla, pero solo por gofmt en archivos de otros agentes**: `internal/api/sistema/adjuntos.go` e `internal/reportes/avisos/alertas_test.go`. Build, vet y pruebas pasan.
- **S6-03**: 19 de 26 operaciones tienen prueba. Cada prueba valida contra `openapi.yaml` con `c.Do`.
- **S6-04**:
  - `period` malo → 400; `family` o `type` fuera del enumerado → 400; id no numérico → 400; persona no externa → 404;
  - todo en centavos enteros; % con HALF_UP;
  - todas las operaciones nuevas son de solo lectura, así que no llevan Idempotency-Key.
- **S6-15**: mensajes en español de México que dicen qué hacer; identificadores en inglés.

| operationId | prueba | 2xx | 4xx probados |
|---|---|---|---|
| getExpenseMonthlyTrend | TestOp_getExpenseMonthlyTrend | 200 (con `family=external`) | 400, 401 |
| getDashboardMonthVsPlan | TestOp_getDashboardMonthVsPlan | 200 (sin plan, con plan en warning, mes pasado) | 400, 401 |
| getDashboardFortnightFlow | TestOp_getDashboardFortnightFlow | 200 (q1 y q2) | 400, 401 |
| getDashboardSavings | TestOp_getDashboardSavings | 200 (vacío y con metas) | 401 |
| getPersonParticipation | TestOp_getPersonParticipation | 200 | 400, 401, 404 |
| search | TestOp_search | 200 (vacío, acentos, top 5, type, ajustes) | 400, 401 |

### Supuestos
Todos son PROVISIONAL y están en el bloqueo, §9 a §12:
- compra particular y estado de una parte;
- `period` en las tarjetas;
- relevancia de la búsqueda y claves de las secciones de Ajustes.

Además:
- en la tendencia, la parte del hogar de cada gasto se atribuye a la categoría del gasto, aunque el papel venga dentro de la despensa;
- `next_goal` toma la meta activa con fecha más próxima, aunque ya esté vencida.

### Bloqueos
- **§11 (bloquea en uso real)**: `dominio.LoadFinEngine` responde 500 si hay categorías del usuario con `code` NULL. Afecta a mis 5 operaciones con motor y a las de presupuesto. Cambio pedido a api-presupuesto (ver solicitudes).
- §9, §10 y §12, cada uno con su opción PROVISIONAL.

### Pendientes de la continuación (7, siguen en 501)
- `getDashboardNextPayment`: el disponible A/B de `motor_cred` (la línea de tiempo `AvailabilityTimeline`) solo se arma en `internal/api/creditos/motor.go` (`snapshot.timeline`), que no se exporta. Pedido: exportar desde `internal/dominio` o `internal/api/creditos` una función `AvailabilityToday(ctx, q, today) ([]{CreditID, Reading, AvailableCents, Confirmed}, error)`. El próximo pago sí se puede sacar con `GastosPersonShares` y su consulta.
- `getDashboardWeekMeals` (motor S4, `recetas`) y `getDashboardBestBasket` (motor S5, `precios`, lista `shopping_list`): sin empezar.
- `listCalendarEvents`: unión de 12 tipos de evento (C1 `CalendarEvent`); sin empezar.
- `queryView` (C1 §8): 23 entidades × 4 modos; sin empezar.
- `createTicketScan`: multipart, guardar el adjunto con la convención de auth-sistema (`internal/api/sistema/adjuntos.go`), job y `ocr.Engine.Recognize`.
- `confirmTicketScan`: una transacción con gasto, precios `ticket_ocr`, inventario y compra a crédito.

### Consumo aproximado (continuación)
Unas 55 acciones, unos 150 KB leídos y unos 45 KB escritos.

## Continuación 2

12-oct-2026 · Opus · **cerrado por tope**: se agregan 2 operaciones, así que van **21 de 26** implementadas y probadas. Siguen 5 en 501 (ver Pendientes de la continuación 2).

### Hecho
- `getDashboardNextPayment` (U3, N-022, M-03 y M-04) en `proximo_pago.go`:
  - próximo pago = el pendiente vencido más antiguo o, si no hay, el pendiente más cercano ≥ hoy; suma los pendientes de ese día, con partes de hogar y de Karen (`scheduled_payment_share` y lo pagado en `credit_payment`); `days_left` negativo si está vencido;
  - disponible: una fila por línea revolvente activa, con lectura A/B y monto de hoy de `motor_cred.AvailabilityTimeline` (solo llamada); la entrada del motor se arma igual que api-creditos (bloqueo §13, PROPUESTA de exportarla).
- `listCalendarEvents` (N-029, N-086, N-104, N-115, N-192, N-208) en `calendario.go`:
  - 10 de 12 tipos: ingresos, pagos de crédito con partes y estado, gastos fijos y venideros, cortes, días de mercado (`precios.NextVisit`), promociones (con `end_date`), fechas meta, movimientos de ahorro y envío del reporte mensual;
  - `period` (por defecto el mes en curso), `credit_id`, `family` (incluyente), `event_type` (lista separada por comas; fuera del enumerado = 400), paginación de 1 a 200 con cursor;
  - orden fijo: fecha, luego el orden de tipos de N-029, luego id.
  - `planned_grocery`, los eventos esperados de partidas y `meal_plan` no se emiten todavía: PREGUNTA en el bloqueo §14.
- Repositorio nuevo `internal/dominio/gastos-inicio_calendario.go`: `GastosCalendarRows`, `GastosMarkets`, `GastosNotificationID`, con su prueba.
- `gastos.go`: 2 líneas más en `Register`.

### Archivos
- Nuevos:
  - `backend/internal/api/gastos-inicio/proximo_pago.go`, `calendario.go` y `calendario_test.go`;
  - `backend/internal/dominio/gastos-inicio_calendario.go` y `gastos-inicio_calendario_test.go`.
- Cambiados:
  - `backend/internal/api/gastos-inicio/gastos.go` (`Register`);
  - `docs/estado/bloqueos/S6-gastos-inicio.md` (§13 a §15);
  - este reporte.

### Verificación
- `go build ./...` y `go vet ./...`: sin errores. `gofmt -l .` solo marca `internal/reportes/reportes_test.go` e `internal/reportes/vistas.go` (de s7-pdf, no de este agente); mis archivos: vacío.
- `go test ./internal/api/gastos-inicio/... ./internal/dominio/... -count=1 -cover`:
  - `gastos-inicio`: **84.7 %**, 21 pruebas `TestOp_*` en PASS;
  - `dominio`: **87.1 %**.
- `make verify-wf3-go`: build y vet pasan; **falla solo por gofmt en los 2 archivos de `internal/reportes`** (otro agente).
- **S6-03**: 21 de 26 operaciones con prueba `TestOp_<operationId>`; cada respuesta se valida contra `openapi.yaml` con `c.Do`.
- **S6-04**:
  - `event_type`, `family` o `period` malos = 400; sin sesión = 401;
  - paginación con cursor probada (`limit=2` y segunda página);
  - centavos enteros; las 2 operaciones son de solo lectura (sin Idempotency-Key).
- **S6-15**: mensajes en español de México que dicen qué hacer; identificadores en inglés; no se registran credenciales.

| operationId | prueba | 2xx | 4xx probados |
|---|---|---|---|
| getDashboardNextPayment | TestOp_getDashboardNextPayment | 200 (sin créditos, vencido, siguiente tras pagar) | 401 |
| listCalendarEvents | TestOp_listCalendarEvents | 200 (mes por defecto, `family` + `event_type`, `credit_id` + cursor, otro mes) | 400 (3 casos), 401 |

### Supuestos
Todos son PROVISIONAL y están en el bloqueo, §13 y §14:
- reglas de la tarjeta de próximo pago;
- `family` y `credit_id` en el calendario;
- estados de los eventos;
- `ref_table = notification` en el reporte mensual.

### Bloqueos
- §14: los eventos esperados de partidas planificadas necesitan que `motor_fin` exporte sus ocurrencias (OK de Axel).
- §15: `getDashboardWeekMeals` y `getDashboardBestBasket` necesitan que api-comidas-salud y api-compras-a exporten sus armados.

### Pendientes de la continuación 2 (5, siguen en 501)
- `getDashboardWeekMeals` y `getDashboardBestBasket`: esperan las exportaciones del bloqueo §15. Con ellas, son unas 40 líneas cada una en `backend/internal/api/gastos-inicio/inicio.go`.
- `createTicketScan`:
  - multipart igual que `uploadExpenseAttachment` (`internal/api/sistema/adjuntos.go`, `adjuntos.Store.Save`), solo imágenes (PDF = 415);
  - `attachment` con `kind = ticket_photo`, `expense_id` NULL y `ocr_status = pending`; respuesta 202 con `TicketScan`;
  - un job `ticket_ocr` (`reg.Deps().Jobs.Register`) que lea la cola con `ocr.New().Recognize` y guarde `ocr_text`, `ocr_engine` y `ocr_status` (`done` o `failed`);
  - en las pruebas, un `ocr.Runner` falso (por `export_test.go`).
- `confirmTicketScan`: una transacción con gasto, precios `ticket_ocr`, inventario e `inventory_movement` y, si aplica, la compra a crédito con `motor_cred`.
- `queryView` (C1 §8): 23 entidades × 4 modos.
- Calendario: `planned_grocery`, los eventos esperados y `meal_plan` (bloqueo §14).

### Solicitudes fuera de alcance
- `internal/api/creditos/motor.go`: exportar el armado de `CreditState` (bloqueo §13).
- `internal/api/comidas` y `internal/api/comprasa`: exportar `WeekPlan` y `BestBasket` (bloqueo §15).
- `internal/motor_fin/plan.go`: exportar `PlannedOccurrences` (bloqueo §14; necesita el OK de Axel).
- `internal/reportes/reportes_test.go` y `vistas.go`: `gofmt -w` (s7-pdf), para que `make verify-wf3-go` pase.

### Paquetes instalados
Ninguno.

### Consumo aproximado (continuación 2)
Unas 56 acciones, unos 190 KB leídos y unos 40 KB escritos.

## Continuación 3

12-oct-2026 · Opus · **cerrado por tope**. Agregué 2 operaciones: van **23 de 26** implementadas y probadas. Siguen 3 en 501 (ver Pendientes de la continuación 3).

### Hecho
- **`queryView`** (C1 §8, decisión G; N-043, N-045, N-046, N-111, N-128, N-135, N-170, N-171, N-180, N-202, N-247, N-281, N-287, N-290) está en `vistas.go`:
  - cubre las 23 entidades con sus modos, sus `group_by` y sus `date_field` de C1 §8;
  - aplica los filtros globales, `q`, `status`, `sort` y la paginación de 1 a 200 con cursor;
  - en el tablero, cada grupo trae `count`, `total_cents` y `share_bp` con HALF_UP;
  - en el calendario hay un grupo por día; en la galería, `image_attachment_id` es la primera foto del gasto;
  - en `expense` se reutilizan los filtros de `listExpenses`, y el importe principal es la parte del hogar.
- Los valores calculados salen de los motores; aquí solo se hacen las llamadas:
  - estado y ahorrado de las metas: `motor_fin.SavingsGoals`;
  - deuda viva por crédito: `motor_cred.LiveDebt`, con la misma entrada que la tarjeta de próximo pago;
  - estado del inventario: `recetas.Engine.InventoryStatus`;
  - «viejo» de precios y productos: `precios.Freshness`, con `price.stale_after_days` y `price.zero_weight_after_days`;
  - estado de los pagos programados: la regla de api-creditos (bloqueo §16, punto 9).
- **`createTicketScan`** (N-075, K3) está en `ticket_ocr.go`:
  - recibe la foto por multipart y la guarda con `adjuntos.Store`, con el límite `attachment.max_size_mb`;
  - crea un `attachment` `ticket_photo`, sin gasto y `pending`, con `audit_log`; responde 202 con `TicketScan`;
  - el job **`ticket_ocr`** (se registra en `Register`) llama a `ocr.Engine.Recognize` y guarda `ocr_text`, `ocr_engine` y `ocr_status` (`done` o `failed`) con `audit_log`;
  - en las pruebas, el motor falso entra con `export_test.go` (`SetOCR`).
- Repositorios nuevos:
  - `internal/dominio/gastos-inicio_vistas.go`: `ViewSQL`, `ViewRecord` y `GastosViewRows`, la consulta genérica. NULL queda ausente, los 0/1 salen como booleano, los instantes en −06:00 y nunca salen credenciales.
  - En `gastos-inicio_tickets.go`: `GastosInsertTicket`, `GastosPendingTickets` y `GastosSetTicketOCR`.

### Archivos
- Nuevos:
  - `backend/internal/api/gastos-inicio/vistas.go`, `vistas_test.go`, `ticket_ocr.go`, `ticket_ocr_test.go` y `export_test.go`;
  - `backend/internal/dominio/gastos-inicio_vistas.go` y `gastos-inicio_vistas_test.go`.
- Cambiados:
  - `backend/internal/api/gastos-inicio/gastos.go` (3 líneas en `Register`);
  - `backend/internal/dominio/gastos-inicio_tickets.go` (3 funciones);
  - `docs/estado/bloqueos/S6-gastos-inicio.md` (§16 y §17);
  - este reporte.

### Verificación
- `go build ./...` y `go vet ./...`: sin errores. `gofmt -l .` solo marca `internal/reportes/reportes_test.go` y `internal/reportes/vistas.go`, que son de s7-pdf, no de este agente.
- `go test ./internal/api/gastos-inicio/... ./internal/dominio/... -count=1 -cover`:
  - `gastos-inicio`: **83.3 %**, con 23 pruebas `TestOp_*` en PASS;
  - `dominio`: **87.2 %**.
- `make verify-wf3-go`: build y vet pasan. **Falla solo por gofmt en esos 2 archivos de `internal/reportes`**, que son de otro agente.
- `python3 -I backend/internal/api/c1gen/build.py --check` = 0.
- **S6-03**: 23 de 26 operaciones tienen su prueba `TestOp_<operationId>`, y cada respuesta se valida contra `openapi.yaml` con `c.Do`.
- **S6-04**:
  - los errores salen con el esquema único: 400 de formato, 422 de las reglas de C1 §8, 404, 413 y 415;
  - la paginación con cursor está probada;
  - todo va en centavos enteros y `share_bp` usa HALF_UP;
  - `createTicketScan` exige Idempotency-Key: sin clave responde 400.
- **S6-15**: los identificadores van en inglés y los mensajes en español de México dicen qué hacer. `password_hash` y `token_hash` nunca salen en `record`, y el registro del job no lleva datos sensibles.

| operationId | prueba | 2xx | 4xx probados |
|---|---|---|---|
| queryView | TestOp_queryView | 200: table, board con categoría, enumerado y estado, calendar con días y mercado recurrente, gallery con foto, cursor; 15 entidades más | 400 (6 casos), 401, 404, 422 (5 casos) |
| createTicketScan | TestOp_createTicketScan | 202 (pendiente), job `done` y `failed` | 400 (campo desconocido, sin Idempotency-Key), 401, 413, 415 (PDF, texto, sin multipart), 422 (sin foto, nombre con ruta) |

### Supuestos
Todos son PROVISIONAL o PREGUNTA y están en el bloqueo, §16 (12 puntos) y §17 (3 puntos). Los principales:
- sin periodo por defecto en `queryView`;
- los filtros globales que no aplican no filtran;
- etiquetas en español de los enumerados, que propuse yo;
- `product.freshness` sale del último precio;
- el job `ticket_ocr` corre cada 5 minutos y además se pide «ahora»;
- si el OCR falla, el ticket queda `failed`.

### Bloqueos
- §16 (PREGUNTA en el punto 6) y §17 (PREGUNTA en el punto 2).
- Siguen abiertos §14 y §15 de la continuación 2.

### Pendientes de la continuación 3 (3, siguen en 501)
- **`confirmTicketScan`**: no lo empecé. Necesita lo siguiente:
  - leer `TicketConfirmRequest` y `TicketConfirmResult` en `docs/contratos/extractos/gastos.yaml`, con `grep -n` de los esquemas;
  - una sola `dominio.WithTx` que haga cuatro cosas:
    1. crear el gasto con renglones y la cuenta de cada renglón, reutilizando `dominio.SaveExpense` y `ExpenseComputeShares`;
    2. crear los precios `ticket_ocr` con la fecha del ticket y la tienda;
    3. crear `inventory_item` e `inventory_movement` `purchase` ligado al renglón;
    4. si aplica, crear la compra a crédito. Su serie de pagos sale de `motor_cred`; hay que ver cómo la arma `internal/api/creditos/compras.go`, que no está exportado, y pedir que se exporte.
  - después, poner `attachment.expense_id`, `ocr_status = confirmed` y vaciar `ocr_draft_json`;
  - devolver 409 si el ticket ya está confirmado.
- `getDashboardWeekMeals` y `getDashboardBestBasket`: siguen esperando las exportaciones del bloqueo §15, porque ni `internal/api/comidas` ni `internal/api/comprasa` las tienen todavía.
- Calendario: `planned_grocery`, eventos esperados y `meal_plan` (bloqueo §14).

### Solicitudes fuera de alcance
- `backend/cmd/finanzas/areas.go`: falta la línea `gastos.Register,` con el import `gastos "finanzas-hogar/internal/api/gastos-inicio"`. Sin ella, las 23 operaciones y el job `ticket_ocr` responden 501 o no corren en `serve`.
- `internal/reportes/reportes_test.go` y `vistas.go`: correr `gofmt -w` (s7-pdf) para que pase `make verify-wf3-go`.
- Exportaciones que propongo:
  - `ProductFreshness` (api-compras-a, §16.8);
  - la regla de estado de `scheduled_payment` (api-creditos, §16.9);
  - el armado de la compra a crédito para `confirmTicketScan` (api-creditos).

### Paquetes instalados
Ninguno.

### Consumo aproximado (continuación 3)
Unas 55 acciones, unos 150 KB leídos y unos 60 KB escritos.

## Continuación 4

12-oct-2026 · Opus · terminada sin llegar al tope. Agregué `confirmTicketScan`: van **24 de 26** operaciones implementadas y probadas. Siguen 2 en 501 porque esperan a otras áreas (ver Pendientes de la continuación 4).

### Hecho
- **`confirmTicketScan`** (N-074, RN-23, K3, decisión E) está en `ticket_confirm.go`. Hace todo en una sola `dominio.WithTx`:
  1. **Gasto**: el mismo alta que `createExpense`. Comparten `newExpenseIn` y `validateExpense`, que se sacaron de `gastos.go` sin cambiar su comportamiento. Lleva renglones, la cuenta de cada renglón y las partes RN-03, con `audit_log`. Aquí se acepta la forma de pago crédito, porque la compra se crea en la misma transacción.
  2. **Precios** `ticket_ocr`:
     - uno por presentación, con la fecha del ticket, la tienda del gasto, `expense_item_id` y `attachment_id`;
     - si falta `unit_price_cents`, el precio es el importe entre las presentaciones, con HALF_UP;
     - si ese día ya había precio del mismo origen, se actualiza;
     - con `audit_log`.
  3. **Inventario** (RN-15):
     - solo productos con `inventory_item` en `stock`;
     - cantidad × contenido con `recetas.PurchaseMovement` (solo la llamada), ÷1000 con HALF_UP;
     - `inventory_movement` `purchase` ligado al renglón;
     - `audit_log` del renglón de inventario y del movimiento.
  4. **Compra a crédito**, si `payment_method = credit`:
     - valida como `createCreditPurchase`, y además pide que el capital sea el total del ticket y la fecha, la del ticket;
     - `motor_cred.QuotePurchase` da las partes y la serie («el último ajusta centavos» es del motor);
     - guarda `credit_purchase`, `credit_purchase_share`, `scheduled_payment` y `scheduled_payment_share` con los repositorios de api-creditos en `dominio`, ligada al gasto y con `audit_log`.
  5. **Ticket**: queda `confirmed`, con su `expense_id` y el borrador vacío (las dos columnas en NULL), con `audit_log`.
  - Ya confirmado = 409, también si dos confirmaciones llegan a la vez (`UPDATE … WHERE expense_id IS NULL`). Cualquier 4xx deja la base igual que antes; la prueba lo comprueba.
- Repositorios nuevos en `internal/dominio/gastos-inicio_tickets.go`: `GastosTicketPresentations`, `GastosExpenseItemIDs`, `GastosStockItem`, `GastosTicketStockIn`, `GastosSaveTicketPrice` y `GastosConfirmTicket`. Tienen su prueba propia en `gastos-inicio_confirmar_test.go`.
- `proximo_pago.go`: el armado del crédito para el motor S3 pasó a `credInput.credit`, que usan `state` y la compra del ticket. No cambia ninguna cifra.

### Archivos
- Nuevos:
  - `backend/internal/api/gastos-inicio/ticket_confirm.go` y `ticket_confirm_test.go`;
  - `backend/internal/dominio/gastos-inicio_confirmar_test.go`.
- Cambiados:
  - `backend/internal/api/gastos-inicio/gastos.go`: `Register` (+1 línea), `newExpenseIn` y `validateExpense`;
  - `backend/internal/api/gastos-inicio/proximo_pago.go`: `credInput.credit`;
  - `backend/internal/dominio/gastos-inicio_tickets.go`: 6 funciones;
  - `docs/estado/bloqueos/S6-gastos-inicio.md`: §18;
  - este reporte.

### Verificación
- `go build ./...` y `go vet ./...`: sin errores.
- `gofmt -l .`: solo marca `internal/reportes/reportes_test.go` e `internal/reportes/vistas.go`, que son de s7-pdf. Son 8 líneas de alineación de espacios. Mis archivos salen limpios.
- `go test ./internal/api/gastos-inicio/... ./internal/dominio/... -count=1 -cover`:
  - `gastos-inicio`: **82.6 %**, con **24** pruebas `TestOp_*` en PASS;
  - `dominio`: **87.0 %**.
- `go test ./... -count=1` (todo el backend): pasa. `go mod tidy -diff` = 0.
- `make verify-wf3-go` y `make verify-all GOOSE=…`: build y vet pasan. **Las dos fallan solo por el gofmt de esos 2 archivos de `internal/reportes`**, que ya fallaban antes de esta continuación. Borré el temporal `datos/tmp/verify-wf0` que dejó la corrida fallida.
- `python3 -I backend/internal/api/c1gen/build.py --check` = 0.
- **S6-03**: 24 de 26 operaciones tienen su prueba `TestOp_<operationId>`. Cada respuesta se valida contra `openapi.yaml` con `c.Do`.
- **S6-04**:
  - el formato malo es 400: sin Idempotency-Key o con un campo desconocido;
  - la regla de negocio es 422, con la ruta JSON en `fields` (`expense.items`, `expense.store_id`, `credit_purchase.*`);
  - 401, 404 y 409 (ya confirmado);
  - todo va en centavos enteros, y el redondeo HALF_UP es de `dominio`;
  - la serie del crédito es la del motor: canon despensa 1, 309450 + 80050 = 389500 en 4 pagos, el primero de 97375.
- **S6-15**: identificadores en inglés y mensajes en español de México que dicen qué hacer. No se registra ninguna credencial.

| operationId | prueba | 2xx | 4xx probados |
|---|---|---|---|
| confirmTicketScan | TestOp_confirmTicketScan | 201 de contado: 5 renglones, 3 precios (uno actualizado), 2 inventarios (pieza y granel), ticket `confirmed` sin borrador, gasto con parte de Karen. 201 a crédito: Kueski, ticket `failed`, 4 pagos que suman 389500 y compra ligada al gasto | 400 (sin clave, campo desconocido), 401, 404, 409 (reconfirmar), 422 (sin tienda, suma, presentación inexistente, compra con contado, crédito sin compra, capital ≠ total, crédito inexistente) |

### Supuestos
Todos están en el bloqueo §18. Hay 8 puntos PROVISIONAL y 2 PREGUNTA. Los principales:
- sin tienda, la respuesta es 422;
- el precio sale de importe entre presentaciones;
- un precio por presentación y ticket;
- el precio del mismo día se actualiza;
- el capital del crédito es el total del ticket;
- se pueden confirmar tickets `pending`, `done` y `failed`.

### Bloqueos
- §18: PREGUNTA a (inventario sin renglón o en ritmo de consumo) y b (de dónde salen las partes de la compra a crédito).
- Siguen abiertos §14 y §15 (continuación 2) y §16 y §17 (continuación 3).

### Pendientes de la continuación 4 (2, siguen en 501)
- `getDashboardWeekMeals` y `getDashboardBestBasket`: siguen esperando las exportaciones del bloqueo §15. Lo revisé el 12-oct: `internal/api/comidas` y `internal/api/comprasa` todavía solo exportan `Register`. Con ellas, son unas 40 líneas cada una en `backend/internal/api/gastos-inicio/inicio.go`.
- Calendario: `planned_grocery`, eventos esperados y `meal_plan` (bloqueo §14; necesita `motor_fin.PlannedOccurrences` con el OK de Axel).

### Solicitudes fuera de alcance
- `backend/cmd/finanzas/areas.go` todavía está vacío. Hay que agregar la línea `gastos.Register,` con el import `gastos "finanzas-hogar/internal/api/gastos-inicio"`. Sin ella, las 24 operaciones y el job `ticket_ocr` responden 501 o no corren en `serve`.
- `internal/reportes/reportes_test.go` y `vistas.go`: correr `gofmt -w` (s7-pdf). Rompe `make verify-wf3-go` y `make verify-all`.
- Exportaciones:
  - `internal/api/comidas`: `WeekPlan` (§15);
  - `internal/api/comprasa`: `BestBasket` (§15);
  - `internal/api/creditos`: el armado del crédito y de la compra para el motor S3 (§13 y §18.6), para que haya un solo armado.

### Propuestas para Axel
- Contestar el §18 a. Propongo la opción (A), la de hoy, más un aviso en P-05: «Este producto no lleva existencia; actívala en Inventario».
- Contestar el §18 b. Propongo que P-05 mande las partes del crédito ya calculadas con la cuenta de cada renglón, y que la API solo valide que sumen el capital (así funciona hoy).

### Paquetes instalados
Ninguno.

### Consumo aproximado (continuación 4)
Unas 40 acciones, unos 115 KB leídos y unos 45 KB escritos.
