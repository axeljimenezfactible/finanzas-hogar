# s7-pdf · Reportes en PDF, vista web mensual y estado de cuenta de Karen

Agente s7-pdf (Opus) · WF-3a · 12-oct-2026 · sin git add/commit (lo hace el integrador).

## HECHO
- `backend/internal/reportes/`: datos del mes y del estado de cuenta (motor S2 vía `dominio.FinEngine`, mismas llamadas que el área de presupuesto), documento con cada cifra y fecha ligada a su campo de la API (`Cell.Src`), plantillas LaTeX embebidas (`plantillas/comun.tex`, `mensual.tex`, `estado-cuenta.tex`) compiladas con `lualatex` por `exec` con tiempo límite (120 s, PROVISIONAL), `-no-shell-escape`, directorio temporal propio y rechazo de cualquier «!» o «undefined» del registro; guardado versionado `datos/reportes/mensual/AAAA-MM/reporte-mensual-AAAA-MM-vN.pdf` (O_EXCL, nunca pisa) con `sha256`, `size_bytes` y `audit_log`; envío por Telegram (interfaz `DocumentSender`, `sendDocument`, texto plano sin `parse_mode`, nunca a personas externas, `notification_delivery` por destinatario, sin repetir a quien ya lo recibió; límite de 50 MB de la Bot API citado en `MaxTelegramDocumentBytes`); job `monthly_report` que lee `reports.monthly_day_of_month` y `reports.monthly_send_time` (cron `M H D * *`; con 29–31 corre del 28 al 31 y solo actúa en el día efectivo, Q2-66 B); si LaTeX falla, job en `error`, sin PDF ni envío.
- `backend/internal/api/reportes/` (paquete `reportesapi`): las 10 operaciones de s7-pdf (`listReports`, `generateMonthlyReport`, `generateExternalStatement`, `downloadReportFile`, `sendReport`, `listMonthlyReportMonths`, `getMonthlyReport`, `getExternalStatement`, `getReportSettings`, `updateReportSettings`); `updateReportSettings` usa `dominio.ApplyParams` y re-programa el job en la misma transacción (si ya corrió hoy, la siguiente corrida cuenta desde el fin del día).
- Bloqueos nuevos: `docs/estado/bloqueos/S7-lualatex.md`, `S7-reporte-mensual-secciones.md`, `S7-estado-cuenta-fechas.md`.

## Archivos
`backend/internal/reportes/{formato,documento,latex,datos,vistas,servicio}.go`, `reportes_test.go`, `plantillas/{comun,mensual,estado-cuenta}.tex`; `backend/internal/api/reportes/{reportes,reportes_test}.go`; los 3 bloqueos; este reporte.

## Verificación
- lualatex mínimo con OpenType: **falla** (falta `luaotfload`, paquete `texlive-luatex`; `luaotfload-tool` no existe) → bloqueo S7-lualatex; las plantillas usan Inter Type 1 local y compilan con lualatex.
- **S7-01**: `FH_S7_BUILD=datos/tmp/s7-build go test ./internal/reportes/ -run PlantillasCompilan` y en ese directorio `latexmk -lualatex -interaction=nonstopmode -halt-on-error <r>.tex` → mensual: rc 0, `grep -c "^!"` 0, `grep -c undefined` 0, `Pages: 2`; estado-cuenta: rc 0, 0, 0, `Pages: 1`; `pdftoppm -r 60` generó 3 PNG (2 + 1). Temporales borrados.
- **S7-02**: `go test ./internal/reportes/ -run Cruce -count=1 -v` → PASS, **M = 36 de 36** cifras y fechas iguales, 0 diferencias. Campos: `monthly:2026-10|` fortnights.fortnights.{0.start_date,1.end_date}, budget.{income,planned_spending,planned_savings,unassigned,actual,remaining}_cents, categories.rows.0.actual_cents, fortnights.fortnights.{0,1}.{start_date,end_date,pay_date,actual_cents,planned_cents}, credits.{month_total,month_household,month_external,fortnight_total,interest_paid_month,interest_paid_year}_cents, live_debt_end_cents, external_participation.0.{accrued_groceries_cents,live_debt_cents,payment_shares.0.due_date,payment_shares.0.amount_cents,purchases.0.purchase_date,purchases.0.total_cents,purchases.0.balance_cents}; `accrual:2026-10|`{accrued,cash,cumulative_difference}_cents (respuesta real de `getAccrualVsCash` del área de presupuesto). R-45 (oct-2026): devengado $1,947.50, efectivo $973.75, diferencia acumulada $973.75 = `oraculo/model.json`; R-46: $1,947.50 / $1,947.50 / $973.75 = oráculo. R-47…R-50: **no cubiertos** (piden la semilla completa; no hay importador).
- **S7-03**: `go test ./internal/reportes/ -run Karen -count=1 -v` → PASS, M = 20 de 20; trae $1,191.50, $308.22, $1,499.72, 4 × $374.93 y Vie 30 oct 2026, Dom 15 nov 2026, Lun 30 nov 2026, Mar 15 dic 2026; no trae el ingreso de Axel ($12,345.67 de la prueba) ni la parte de Axel; `sendReport` de un estado de cuenta = 409; Axel (no externo) = 422; Karen suscrita a `monthly_report` no recibe envío (S7-10).
- **S7-10**: `go test ./internal/reportes/ -run Mensual -count=1 -v` (reloj 1-nov-2026 08:00 America/Mexico_City) → PASS: PDF guardado `reportes/mensual/2026-10/reporte-mensual-2026-10-v1.pdf` de 96,214 bytes (límite de sendDocument 50,000,000), 2 `sendDocument` al servidor falso (Axel y Ximena, sin `parse_mode`), 2 entregas `sent` con HTTP 200; plantilla rota → job `error`, 0 PDF nuevos, 0 envíos; hora a 09:30 → horario `30 9 1 * *` y siguiente 1-dic-2026 09:30; día 31 → `30 9 28-31 * *`, nada el 28 ni el 29, sale el 30-nov-2026 (versión 2 sin pisar la 1); `notification.params_json` sin `send_time` ni `day_of_month`.
- **S7-11**: `go test ./internal/reportes/ -run VistaMensual -count=1 -v` → PASS, 33 cifras de `getMonthlyReport` = PDF, 0 diferencias; la vista enlaza el PDF; por defecto el último mes cerrado (2026-09). La pantalla de U1 es del frontend (fuera de alcance).
- **S7-12**: `go test ./internal/reportes/ -run FormatoS712 -v` → 39 cifras con `^\$-?[0-9]{1,3}(,[0-9]{3})*\.[0-9]{2}$` (el signo negativo va antes de `$`: PREGUNTA 4 del bloqueo de secciones), 0 fechas ISO, 0 palabras en inglés, 10 líneas en la salida para el revisor.
- **S7-13**: `grep -n -i definecolor internal/reportes/plantillas/*.tex` → 7 colores, cada uno = su token de C2 (prueba `TestColoresC2`); contraste WCAG: texto/página 15.73:1, atenuado/página 5.63:1, marca/página 8.49:1, texto/banda 15.09:1, marca/banda 8.15:1, texto/encabezado de tabla 14.07:1 (todos ≥ 4.5). Calculado en Go con la fórmula WCAG 2.x; no se corrió el script de U2 (pendiente).
- Operaciones: `go test ./internal/api/reportes/ -v` → 10 de 10 `TestOp_*` PASS (2xx y 4xx de cada una; respuestas validadas contra openapi.yaml por testutil).
- `go vet ./internal/reportes/ ./internal/api/reportes/` OK; cobertura 77.2 % y 75.6 %; `grep -rn "http://\|https://"` en el código (sin pruebas) = 0 coincidencias; `build.py --check` = 0.
- `make verify-all` no se corrió completo: `internal/api/gastos-inicio` no compila en este momento (cambio en curso de otro agente: `mdate redeclared`), y `internal/reportes/avisos` (de otro agente) depende de él.

## Supuestos
- Directorio del paquete HTTP `internal/api/reportes` con nombre de paquete `reportesapi` (evita choque con `internal/reportes`).
- `file_path` se guarda relativo a `DataDir` (sobrevive a restaurar en otra ruta).
- `generateMonthlyReport` y `generateExternalStatement` compilan dentro de la petición (2–3 s) y responden 202 ya con el renglón.

## Bloqueos
S7-lualatex, S7-reporte-mensual-secciones, S7-estado-cuenta-fechas (rutas arriba).

## Pendientes
- Integración: línea en `backend/cmd/finanzas/areas.go` e instalar el cliente de Telegram con `reportes.SetSender(...)` al arrancar `serve`.
- Pantalla P de la vista web mensual (frontend).
- Cruce R-47…R-50 cuando exista el importador de la semilla.
- Contraste con el script de U2.

## Propuestas para Axel
- Instalar `texlive-luatex` para usar Inter OpenType (bloqueo S7-lualatex, opción A).
- C1 1.5: `accrual_vs_cash` en `MonthlyReport` y `due_date` en `ExternalStatement.payment_shares`.
- Contenido del reporte mensual (PROVISIONAL, K31): resumen del mes, devengado contra efectivo, por categoría, por quincena, créditos y participación de Karen; comidas y canasta cuando haya cálculo compartido.

## Paquetes instalados
Ninguno.

## Consumo aproximado
≈ 60 acciones, ≈ 240 KB leídos, ≈ 110 KB escritos.
