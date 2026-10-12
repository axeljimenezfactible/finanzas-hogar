# motor-S2-go · Reporte final (WF-1, 11-oct-2026)

## Hecho
- `backend/internal/motor_fin`: todas las firmas de `S2-interfaz.md` 1.2 (§3.1–§3.15), con enteros de centavos y `halfUp` exacto. Incluye Calendario, Periods (meses y ventanas), Horizon, LiveDebt (K19), Averages/MonthlyAverage (K20, K39), StableValue, PaperPlan, BudgetMonth (80/100 %), ApplyBudgetChange (K2), PlanVsActual, FortnightFlow (K8), Projection (K21 con interruptor), SavingsGoals, las 6 vistas del oráculo y `Eval` (evaluador de consultas con argumentos estrictos).
- `backend/cmd/motores`: `main.go` (registro `register`, `dump`, `$ref` + parche RFC 7396, recorte FORMATO §4, salida determinista), `s2.go` (`init` registra «s2») y `main_test.go`.
- Pruebas: TestCanon (85), TestOraculoHojas (159), TestOraculoFixtures (27), TestSumaPartes, TestQuincenaFinDeMes (4), TestInvarianteDeudaViva (7), TestPromedios, TestProyeccionDefecto, TestProyeccionCadencia, TestFlujoQuincenal, TestPresupuestoUmbrales, TestPresupuestoScope, TestMetasAhorro, TestZonaHoy y TestErrores. Todas leen `oraculo/model.json` y las fixtures en tiempo de prueba. `seedFromModel` sigue la receta §5.3.

## Archivos
backend/internal/motor_fin/{types,input,ledger,summary,plan,projection,eval}.go y {helpers,canon,engine,cover}_test.go; backend/cmd/motores/{main,s2,main_test}.go; este reporte. No escribí en `oraculo/` porque la interfaz §5.0 lo prohíbe. Las fixtures las escribió el agente Python durante mi corrida y las pruebas las usan.

## Verificación
- `gofmt -l .` sale vacío. `go vet ./internal/motor_fin/... ./cmd/motores/...` sale limpio. `go vet ./...` falla solo en `internal/precios/ocr` (otro agente; `undefined: NormalizeDescription`).
- `go test ./internal/motor_fin/... -count=1 -cover -v`: PASS, cobertura **88.8 %**, 0 SKIP. `cmd/motores`: PASS, 78.4 %.
- S2-01: TestCanon 85/85 PASS. Las 30 filas S2 salen del motor. Las 55 filas S3 pasan los pasos (1)–(3) de §5.1 contra `insumo-s3.json`, y `semilla-oraculo.json` es igual a `seedFromModel`.
- S2-02: TestOraculoHojas 159/159 hojas iguales. `compara.py` todavía no existe (dueño: agente Python).
- S2-03: TestOraculoFixtures 27/27 fixtures iguales a su `esperado`. Los cruces `vista_sensibilidad` (15 hojas) e `interes_por_monto` (28) también coinciden.
- S2-04: sha256 de `model.json` = 1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4 antes y después.
- S2-05: 40 casos con Σ partes = total, diferencia 0.
- S2-06: feb-2027 (q2 del 16 al 28, paga el 28), mes de 30 días, mes de 31 días (el gasto del 31-oct cae en q2: 50000, y el efectivo de octubre da 147375) y semana 2026-W42 (12 al 18).
- S2-07: los 7 cierres dan 292125=97375+194750, 97375=97375+0, 377606=125869+251737, 125870=125870+0, 377606=125869+251737, 125870=125870+0 y 629344=125870+503474. La forma sin saldo vale solo en P2, P4 y P6.
- S2-08: diario/semanal de oct 3141/21988, nov 6492/45442, feb 8991/62935 y despensa de oct 6282/43976. Los mensuales son 218396 y 239375.
- S2-09: los meses 1–6 dan 12 de 12 cifras iguales a P1–P6. Con el interruptor apagado, los meses 7–12 traen solo los pagos 2027-04-15…06-15. Con el interruptor encendido aparece la marca PROVISIONAL. Una cadencia antes del mes 7 da error.
- S2-10: 16 quincenas: 59882/37493 y 88375×3+88377. Cada ciclo da diferencia 0. Valor estable − pago: +1 en 9 quincenas y −1 en 3 (se informa, no se corrige).
- S2-11: fronteras exactas (79, 79.999, 80, 99, 100, 101 %) y RN-22 en 88 %.
- S2-12: Σ del plan = objetivo, el último pago ajusta y RN-18 da 50.0 %.
- S2-13: 0 líneas con float.
- S2-14: medianoche local.
- `go run ./cmd/motores s2 dump --fixtures ../oraculo/fixtures/s2 --out ../datos/tmp/motor-s2.json`: código de salida 0, 27 casos, 70 consultas y 0 marcas `__`. Los 2 errores (`out_of_range`, `invalid_argument`) los pide la propia fixture. Borré los temporales.

## Hallazgos
Ninguna cifra canon difiere.

## Supuestos (ninguno cambia cifras canon)
1. `is_active` ausente cuenta como verdadero. `accrual_months` e `interval_count` ausentes o 0 cuentan como 1.
2. Con ventanas, un costo fuera de toda ventana no se devenga, y la diferencia acumulada cuenta desde la primera ventana pedida.
3. La primera ocurrencia quincenal es la primera quincena cuyo día de pago es ≥ `start_date`.
4. Una compra sin gasto a crédito que la cite se devenga sola.
5. En `metas` (PROVISIONAL, S2-metas): `months_left` cuenta el mes en curso, y `state` ∈ open/in_progress/achieved/overdue.
6. El sufijo de `deuda_viva_<día>_<mes>` lleva el día sin cero a la izquierda.

## Bloqueos
No abrí bloqueos nuevos. Uso las marcas de los bloqueos existentes: S2-proyeccion-interruptor, S2-flujo-quincenal, S2-metas, S2-metas-resumen y S2-presupuesto-total.

## Pendientes
- Correr `python3 -I oraculo/compara.py s2` cuando exista.

## Solicitudes fuera de alcance
- `internal/precios/ocr/parse.go:146`: `NormalizeDescription` no está definida y rompe `go vet ./...`. Corresponde al dueño de precios/OCR.
- El prompt de esta tarea permitía escribir fixtures si la carpeta estaba vacía, pero la interfaz §5.0 lo prohíbe. Seguí la interfaz. El integrador decide si hay que alinear el prompt.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: los umbrales de S2 ya están decididos (K44).

## Consumo aproximado
36 acciones; ~215 KB leídos; ~148 KB escritos.
