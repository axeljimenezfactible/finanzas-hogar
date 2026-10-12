# api-presupuesto · WF-3a · reporte

12-oct-2026 · Opus · cerrado_por_tope = false.

## Hecho
- Las **28 operaciones** de `presupuesto-a.yaml` y `presupuesto-b.yaml` registradas en `presupuesto.Register` (una sola función, CONVENCIONES §2). Errores con el esquema único, paginación 1–200 con `api.ParsePage`, `period` con `api.QueryPeriod`, Idempotency-Key automática (createBudget, createIncome, createSavingsMovement), `audit_log` en cada alta, cambio y baja, dinero en centavos `int64`.
- Cálculos con el **motor S2** sin reescribirlo: `BudgetMonth` (resumen, categorías, alertas), `ApplyBudgetChange` (K2 `scope` onward/this_month; se guarda la diferencia de renglones; prueba con el ejemplo de C1-datos §7.1), `Periods` (quincenas y despensa devengada), `PlanVsActual`, `FortnightFlow`, `Projection` + `LiveDebt` (acumulados Q2-65 A), `LiveDebt` (deuda viva), `SavingsGoals` (metas, resumen, tendencia).
- `dominio/presupuesto_motor.go`: `LoadFinEngine` arma `motor_fin.Input` desde la base (personas, categorías, reglas de reparto resueltas, gastos y partes, compras, partes, pagos programados y abonos, ingresos, partidas, presupuesto, metas, movimientos, umbrales de `app_setting`); `At(hoy, corte)` y `Only(hoy, gasto)` preparan copias. La co-titular en campos de dinero va a nombre del titular (C-09).
- `dominio/presupuesto_datos.go`: repositorios de budget, income, savings_movement, period (y lectura de savings_goal y planned_item) con `audit_log`.

## Archivos
`backend/internal/api/presupuesto/{presupuesto.go, flujo.go, ingresos.go, presupuesto_test.go}`, `backend/internal/dominio/{presupuesto_motor.go, presupuesto_datos.go, presupuesto_test.go}`, `docs/estado/bloqueos/S6-presupuesto-{motor,periodos,quincenas}.md`, este reporte.

## Verificación
- `go build` y `go vet` de `./internal/api/presupuesto/ ./internal/dominio/`: sin errores. `gofmt -l`: vacío. `build.py --check` = 0.
- `make verify-wf3-go`: build/vet/gofmt/tidy OK; `go test ./...`: 19 paquetes OK, 8 sin pruebas.
- `go test ./internal/api/presupuesto/... ./internal/dominio/... -count=1 -cover`: presupuesto **77.2 %**; dominio 75.2 % en total (mis archivos: `presupuesto_datos.go` 87.5 %, `presupuesto_motor.go` 89.5 %; lo bajan `comidas-salud_datos.go` y `compras-a_repo.go`, 0 %, de otros agentes en curso).
- **S6-03**: 28 de 28 operaciones con prueba `TestOp_<operationId>` que valida cada respuesta contra `openapi.yaml` (`testutil.Client.Do`), 2xx y al menos un 4xx:

| Operación | Prueba | 4xx |
|---|---|---|
| listPeriods | TestOp_listPeriods | 400, 422 |
| listBudgets | TestOp_listBudgets | 400 |
| createBudget | TestOp_createBudget | 409, 422, 400 |
| updateBudget | TestOp_updateBudget | 422, 400, 404 |
| deleteBudget | TestOp_deleteBudget | 404 |
| getBudgetMonth | TestOp_getBudgetMonth | 422 |
| getBudgetByCategory | TestOp_getBudgetByCategory | 400 |
| getBudgetFortnights | TestOp_getBudgetFortnights | 422 |
| getAccruedGroceries | TestOp_getAccruedGroceries | 400 |
| getBudgetAlerts | TestOp_getBudgetAlerts | 401 |
| getBudgetMonthlyTrend | TestOp_getBudgetMonthlyTrend | 422 |
| getCashFlowFortnights | TestOp_getCashFlowFortnights | 401 |
| getAccrualVsCash | TestOp_getAccrualVsCash | 401 |
| getCashFlowAssumptions | TestOp_getCashFlowAssumptions | 401 |
| getLiveDebt | TestOp_getLiveDebt | 400 |
| listIncomes | TestOp_listIncomes | 400 |
| createIncome | TestOp_createIncome | 422, 400 |
| getIncomeMonth | TestOp_getIncomeMonth | 422 |
| getIncomeMonthlyTrend | TestOp_getIncomeMonthlyTrend | 401 |
| updateIncome | TestOp_updateIncome | 404, 422 |
| deleteIncome | TestOp_deleteIncome | 404 |
| getSavingsSummary | TestOp_getSavingsSummary | 401 |
| getSavingsMonthlyTrend | TestOp_getSavingsMonthlyTrend | 401 |
| listSavingsGoals | TestOp_listSavingsGoals | 400 |
| listSavingsMovements | TestOp_listSavingsMovements | 400 |
| createSavingsMovement | TestOp_createSavingsMovement | 422 |
| updateSavingsMovement | TestOp_updateSavingsMovement | 404 |
| deleteSavingsMovement | TestOp_deleteSavingsMovement | 404 |

- **S6-04**: errores con el esquema único de C1 §3 (formato 400, regla 422, traslape 409 `budget_overlap`); páginas ≤ 200 (`limit=500` → 400 probado); cursor con `next_cursor` probado; fechas ISO, dinero en centavos; Idempotency-Key exigida (sin clave → 400 probado en createIncome).
- **S6-15**: identificadores en inglés; mensajes en español de México; no se registran credenciales (sin datos sensibles en el área).
- Cifras canon probadas: RN-18 (3,146.72 de 6,293.44 = 500 décimas, 50 %); RN-21 (deuda viva = pagos posteriores al corte); C1-datos §7.1 (tres pasos de K2, quedan 2 renglones).

## Supuestos
- Orden por defecto de `listIncomes`: `-income_date` (C1 no lo dice). `getBudgetMonthlyTrend`, `getIncomeMonthlyTrend` y `getSavingsMonthlyTrend`: los 12 meses que terminan en el mes en curso.
- `updateBudget` con monto y vigencia/notas en el mismo cuerpo: primero vigencia y notas, luego el monto con su alcance.
- `getBudgetAlerts.crossed_date`: el primer día del mes en que el motor, con los gastos hasta ese día, da % ≥ umbral (80/100 de `app_setting`).
- `getAccrualVsCash`: acumulado devengado de `LiveDebt` al cierre del mes, diferencia acumulada de `Projection`, efectivo acumulado = devengado − diferencia (si las dos diferencias del motor no coinciden se registra un aviso).
- `getIncomeMonth.expected_cents` sale de `planned_item` (regla literal de C1); `getBudgetMonth.income_cents` = total de `getIncomeMonth`.
- Abono sin pagador o de la co-titular → a nombre del titular (C-09).

## Bloqueos
`S6-presupuesto-motor` (códigos de despensa/papel, cadencia, fecha por omisión de deuda viva, gasto variable estimado), `S6-presupuesto-periodos`, `S6-presupuesto-quincenas`.

## Pendientes
- Cadencia de los meses 7–12 (bloqueo S6-presupuesto-motor, punto 2). `PlannedItem.next_date` (opcional) no se llena en `getCashFlowAssumptions`.
- Cobertura total de `internal/dominio` < 85 % por archivos de otros agentes (mis dos archivos pasan de 85 %).
- `backend/cmd/finanzas/areas.go`: falta la línea del área (solicitud fuera de alcance).

## Propuestas para Axel
1. Fijar en la semilla los códigos `groceries` y `paper` (o decir cuáles son).
2. Decidir cómo reparte el plan mensual entre quincenas (por días, mitad o función del motor).
3. Decidir la fecha por omisión de la deuda viva (fin de la proyección o 10-abr-2027).

## Paquetes instalados
Ninguno.

## Consumo aproximado
~56 acciones · ~185 KB leídos · ~125 KB escritos.
