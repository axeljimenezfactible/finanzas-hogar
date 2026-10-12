# REVISIÓN backend-core · ronda 1 de 3

12-oct-2026 · agente **revisor-backend-core-r1** (revisor independiente; no escribió ni corrigió nada de lo revisado). Agentes revisados: backend-base, auth-sistema, semilla, api-gastos-inicio, api-presupuesto, api-comidas-salud. Criterios: `docs/criterios/S6.md` 1.2 (todo) y C3-03, C3-04 y C3-09 en lo que toca al servidor. Repositorio en HEAD `fae1890` («[Puerta 2 bis] K53») con el trabajo de WF-3a sin commit.

## Veredicto

**NO APROBADO** · 6 altos · 3 medios · 4 bajos.

Lo que más pesa: `go test ./...` **falla**. La causa principal es un solo error de api-presupuesto: `LoadFinEngine` no acepta categorías sin `code`, y el esquema dice que **todas las categorías que crea el usuario** tienen `code` NULL. Con la primera categoría que capture Axel, unas 20 operaciones de Presupuesto, Flujo, Inicio y Reportes responden 500. Además, el binario `finanzas serve` no registra ninguna área (`areas.go` vacío): con el binario real, `login` responde 501.

Todas las evidencias de abajo salen de comandos que corrí yo, con temporales en `datos/tmp/rev-bc-r1/`, que ya borré. Mis pruebas de las áreas se metieron con `go test -overlay`, sin escribir en `backend/`.

## Tabla criterio → estado → evidencia del revisor

| Criterio | Estado | Evidencia (comando del revisor → salida) |
|---|---|---|
| S6-01 (alta) | **NO CUMPLE** | `go version` = go1.27.2 · `go.mod` = `go 1.27` · `go build ./...` OK · `go vet ./...` OK · `gofmt -l .` → `internal/reportes/reportes_test.go`, `internal/reportes/vistas.go` (dueño s7-pdf) |
| S6-02 (alta) | **NO CUMPLE** | `go test ./... -count=1` → **FAIL** `internal/api` (`TestContrato_Operaciones`, `TestContrato_Rutas`); 31 paquetes ok. Con umbral: dominio **87.0 %**, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 %, precios 88.8 % (las 5 ≥ 85 %). Sin umbral: api 81.9 % (el perfil se escribe aunque la prueba falle; el `&&` del comando literal lo omite), auth 83.6 %, jobs 92.3 %, reportes 77.2 %, telegram 80.0 % |
| S6-03 (alta) | **NO CUMPLE** | (a) `build.py --check` → «los 23 archivos coinciden byte a byte», exit 0 · (b) goose up + `verify_c1.py --db` → «rutas en routes.go: 216 · en openapi.yaml: 216 · sobrantes: [] · faltantes: []», «TODO OK», exit 0 · (c) `TestContrato_Rutas`: «contrato 216, servidor 216, sobrantes 0, faltantes 0, **sin implementar 2**» · (d) `TestOp_*` del backend: **214 de 216**; faltan `getDashboardBestBasket` y `getDashboardWeekMeals` (501). Binario: `finanzas serve --dev-http` → `POST /api/v1/auth/login` = **501** `not_implemented` |
| S6-04 (alta) | **NO CUMPLE** | La prueba de contrato que recorre las 216 operaciones da «**191 de 216** operaciones cumplen». Lo que sí cumple: K9 de gasto (`createExpense` con la misma clave: 201, luego 200 con `Idempotent-Replayed: true`; a las +25 h, purgada la tabla `idempotency_key`, 200 y **1** fila en `expense`); K9 de peso (201 → 200 repetido, 1 fila); `limit=200` → 200 y `limit=201` → 400 en gastos, usuarios y recetas. Muestra de 20 mensajes al azar (`shuf` sobre 214 mensajes de las 6 áreas): 0 en inglés |
| S6-05 (alta) | PARCIAL | `go test ./internal/auth/... -v`: 15 PASS. Mi prueba: `login` de axel → 200 con `HttpOnly`, `Secure`, `SameSite=Strict`, `Max-Age=7776000`; contraseña nueva de 11 caracteres → 422 y de 12 → 204; `login` de karen → 401. En la base sembrada: `substr(password_hash,1,10)` = `$argon2id$` ×2 y `COUNT(user_account)` = 2. En contra: la prueba de permisos solo registra `auth.Register` (M1) y el binario no sirve `login` (H4) |
| S6-06 (alta) | CUMPLE | `-run Intentos -v` → 3 PASS. Leí el código: (a) minuto 14 = 429 y minuto 16 = 200; (b) fallos en 0, 1, 2, 3 y 14: minuto 14.5 = 429 y 15.5 = 200; (c) por usuario y por IP, lo que llegue primero a 5, como dice C1-datos.md l.136 |
| S6-07 (alta) | CUMPLE | `-run Sesion -v` → `TestSesion_DosDispositivosYCierre` y `TestSesion_90DiasRenovableConElUso` PASS |
| S6-08 (alta) | CUMPLE | `go test ./internal/jobs/... -v` → `TestJobsDia30Febrero2027EIdempotencia`, `TestJobsErroresApagadoYRunNow` y `TestCron` PASS; la prueba cubre el día 30, febrero de 2027, el doble disparo y el reinicio con un job vencido (l.18 y l.68–82) |
| S6-09 (alta) | CUMPLE (hora y retención PROVISIONAL) | Base migrada y sembrada con el binario → `finanzas respaldo` → `finanzas-20261012-042019-manual-1.db`; copia a `datos/tmp/restaura.db`; `PRAGMA integrity_check` = `ok`; `COUNT(*)` de 56 tablas en el origen y en la restauración: **0 diferencias** (p. ej. audit_log 1211, recipe_ingredient 561, scheduled_payment 16). Ver B2 |
| S6-10 (media) | CUMPLE (factor PROVISIONAL) | `-run Adjunto -v` → 3 PASS (subir y bajar con y sin sesión, `../`, 10 MB y 10 MB + 1, límite cambiado). Factor: 1 MB = 1,048,576 bytes, PROVISIONAL en `S6-adjuntos` (C1 no declara el factor) |
| S6-11 (alta) | NO VERIFICADO | `mkcert` no está instalado; el certificado lo emite el integrador (`ops/mkcert.md`). No se declara cumplido. Android: pendiente de Axel |
| S6-12 (alta) | Fuera de esta ronda | La PWA la revisa el revisor de frontend; el servidor la sirve desde `frontend/dist` (lo dice backend-base; no lo probé) |
| S6-13 (alta) | CUMPLE | `-run IA -v` → `TestIAApagadaSinLlamadas` y `TestIA_ApagadaNoLlamaALaRed` PASS; el `grep` literal de llaves sobre `backend` = **0** |
| S6-14 (media) | Parcial (WF-4) | Mi `serve` en `127.0.0.1:18447` → `ss` muestra solo `127.0.0.1:18447`. La IP de la LAN va en `ops/runbook.md`, que todavía no existe (WF-4) |
| S6-15 (media) | CUMPLE | `-run Log` → ok. Mi prueba: el registro no contiene la credencial falsa, la contraseña nueva ni `fh_session=`; `serve.log` del binario: 0 coincidencias. Identificadores en inglés; mensajes en español (B3) |
| S6-16 (media) | **NO CUMPLE** (semilla) | `ls docs/estado/bloqueos/S6-*` → 16 archivos; los de backend-base, auth-sistema, gastos, presupuesto y comidas tienen marca PROVISIONAL. Los 3 de semilla se llaman `semilla-*.md` y ese `ls` no los encuentra (M3) |
| C3-03 (servidor) | NO CUMPLE | Las rutas están en `routes.go` (216/216), pero el servidor real no registra ninguna área (H4) y 2 operaciones siguen en 501 (H2) |
| C3-04 (servidor) | NO CUMPLE | `go test ./... -count=1` falla en `internal/api`: una fila «probado» que cite `TestContrato_Operaciones` no pasa |
| C3-09 (servidor) | CUMPLE | Semilla real con contraseñas aleatorias que no se imprimieron: «cotejo OK: C3-09 … A $10,000.00 / B $6,905.50, pago del 2026-10-30 $973.75, deuda viva al 2027-04-10 $6,293.44 = oráculo». Cruce mío con `oraculo/model.json`: 16 pagos (fecha y centavos) = `despensas[*].pagos`, **0 diferencias**; primer pago 2026-10-30 97375; Σ de pagos con vencimiento después del 2027-04-10 = 629344; `credito.eventos[0]` = 10000.0 / 6905.5 |

### 5 operaciones al azar por área, probadas por mí (`c.Do` valida contra `openapi.yaml`)

- **auth-sistema**: `listUsers` 200 · `listSessions` 200 · `getSecuritySettings` 200 · `listJobs` 200 · `getDataSettings` 200 · `listUsers?limit=201` 400 · `login` de karen 401 · `login` de axel 200 · `changeOwnPassword` 422/204.
- **api-gastos-inicio**: `createExpense` 201/200/200 (1 fila) · `listExpenses` 200/400 · `getExpenseTotals` 200 · `getExpenseSummary` 200 (promedio diario 3141 y semanal 21988 del motor = 31.41 y 219.88 de RN-27) · `listTicketScans` 200 · `getDashboardSavings` **500** con una categoría del usuario.
- **api-presupuesto**: `listPeriods`, `listBudgets`, `listIncomes`, `getLiveDebt` y `getBudgetMonth` → 200 sin categorías del usuario; después de `INSERT INTO category (name, kind) VALUES ('Mascotas','variable')`, `getBudgetMonth` y `getLiveDebt` → **500** `internal_error`.
- **api-comidas-salud**: `listMealSlots` 200 · `listRecipes` 200 y `limit=201` 400 · `getMealPlanWeek` 200 · `listHealthProfiles` 200 · `createWeightEntry` 201 → 200 repetido; Karen 422 `external_person`.
- **Cifras fuera del motor**: `grep` de aritmética con centavos en las 3 áreas de API. Encontré la resta de B1 y dos repartos PROVISIONAL con bloqueo (`presupuesto.go:573-577`, plan mensual por quincena por días, `S6-presupuesto-quincenas`; `comidas/plan.go:272`, promedio de 7 días, `S6-comidas-calculos`). No encontré cifras que dupliquen una función del motor sin declararlo.

## Hallazgos

### Altos

**H1 · Una categoría del usuario (`code` NULL) tumba Presupuesto, Flujo, Inicio y Reportes (500).**
- Criterios: S6-02 y S6-04.
- Ruta: `backend/internal/dominio/presupuesto_motor.go:226-236`.
- Agente dueño: api-presupuesto.
- Qué pasa: `LoadFinEngine` escanea `category.code` en `mf.Category.Code` (string). El esquema (`00001_esquema.sql:272`) dice «NULL en las del usuario». Cualquier categoría capturada en Ajustes (U4 P-26.04) produce `converting NULL to string`, y las operaciones responden 500.
- Operaciones afectadas en `TestContrato_Operaciones`: getBudgetMonth, getBudgetByCategory, getBudgetFortnights, getBudgetAlerts, getBudgetMonthlyTrend, getAccruedGroceries, getCashFlowFortnights, getCashFlowAssumptions, getAccrualVsCash, getLiveDebt, getSavingsSummary, getSavingsMonthlyTrend, listSavingsGoals; de gastos-inicio, getDashboardFortnightFlow, getDashboardMonthVsPlan, getDashboardSavings, getExpenseMonthlyTrend y getPersonParticipation; de s7-pdf, generateMonthlyReport y getMonthlyReport.
- Lo reproduje con mi propia prueba: antes del INSERT, 200; después, 500.
- Corrección sugerida:
  1. Escanear en `sql.NullString`.
  2. Cuando sea NULL, darle al motor un código interno estable que no choque con el CHECK `a-z0-9_`, solo en memoria (p. ej. `"cat#"+id`). Marcarlo PROVISIONAL en `S6-presupuesto-motor` si cambia lo que ve el usuario.
  3. Agregar una prueba `TestOp_getBudgetMonth` con una categoría sin código.

**H2 · Dos operaciones de C1 sin implementar ni probar.**
- Criterios: S6-03 y C3-03.
- Ruta: `backend/internal/api/gastos-inicio/inicio.go`.
- Agente dueño: api-gastos-inicio.
- Qué pasa: `getDashboardWeekMeals` y `getDashboardBestBasket` responden 501. `TestOp_*` cubre 214 de 216 operaciones.
- Corrección sugerida: que `internal/api/comidas` exporte la semana del plan y `internal/api/comprasa` la mejor canasta (bloqueo `S6-gastos-inicio` §15), o que se armen desde `dominio` y los motores S4/S5. Agregar `TestOp_getDashboardWeekMeals` y `TestOp_getDashboardBestBasket`.

**H3 · `saveTicketScanDraft` acepta un cuerpo sin `lines`, que es obligatorio.**
- Criterios: S6-04 y S6-02.
- Ruta: `backend/internal/api/gastos-inicio/tickets.go:112-123`.
- Agente dueño: api-gastos-inicio.
- Qué pasa: `TicketScanDraft.required = [lines]` en `openapi.yaml`. Con `{}` el servidor responde 200 y guarda un borrador vacío. `TestContrato_Operaciones/saveTicketScanDraft` falla: «cuerpo vacío con campos obligatorios: 200, se esperaba 422».
- Corrección sugerida: `v.Required(d.Lines != nil, "lines")` antes de guardar, para responder 422 `validation_failed` con `fields`, y una prueba del caso.

**H4 · El binario no registra ninguna área.**
- Criterios: S6-03 (c), S6-05 y C3-03.
- Ruta: `backend/cmd/finanzas/areas.go:7`.
- Agente dueño: backend-base, dueño de `areas.go`; lo cierra el integrador.
- Qué pasa: `var areas = []func(*api.Registry){}`. Con el binario (`serve --dev-http` en 127.0.0.1:18447), `POST /api/v1/auth/login` responde 501 `not_implemented`. Nadie puede entrar, y las 214 operaciones de las áreas solo funcionan dentro de las pruebas.
- Corrección sugerida:
  1. Agregar a `areas` las líneas `auth.Register`, `gastos.Register` (import `finanzas-hogar/internal/api/gastos-inicio`), `presupuesto.Register`, `comidas.Register` y las de las áreas restantes (comprasa, comprasb, creditos, ajustes, avisos/telegram, reportes), con sus imports.
  2. Agregar en `cmd/finanzas` una prueba que arme el registro con `areas` y compruebe que ninguna operación responde 501.

**H5 · `gofmt -l` no está limpio.**
- Criterio: S6-01.
- Rutas: `backend/internal/reportes/reportes_test.go` y `backend/internal/reportes/vistas.go`.
- Agente dueño: s7-pdf. No es un agente revisado aquí, pero S6-01 es global.
- Qué pasa: rompe `make verify-wf3` y `make verify-all`.
- Corrección sugerida: `gofmt -w` sobre los dos archivos.

**H6 · Fallas de la prueba de contrato en áreas no revisadas aquí.**
- Criterios: S6-02 y S6-04.
- Agentes dueños: s7-pdf y api-compras-a.
- Qué pasa:
  - `sendReport` responde 409 `conflict` («El estado de cuenta no se envía por Telegram…») donde C1 declara 202 para el caso de la prueba. Dueño: s7-pdf.
  - `setStoreProductGroups` acepta un cuerpo vacío con campos obligatorios y responde 204 en lugar de 422. Dueño: api-compras-a.
- Corrección sugerida: que cada dueño ajuste el manejador o su caso en C1 (si el 409 es correcto, la prueba de contrato debe mandar un caso válido; lo decide el dueño de C1). Volver a correr `go test ./internal/api/ -run TestContrato -count=1`.

### Medios

**M1 · La prueba de «mismos permisos» no recorre las operaciones reales.**
- Criterio: S6-05.
- Ruta: `backend/internal/auth/auth_test.go:270-274`.
- Agente dueño: auth-sistema.
- Qué pasa: `TestPermisos_MismasOperacionesAmbasCuentas` arma el servidor solo con `auth.Register`. Unas 185 de las 212 operaciones responden 501 a las dos cuentas, así que «el mismo resultado» no prueba permisos iguales sobre las operaciones de lectura y escritura de C1.
- Corrección sugerida: usar la misma lista `areas` del binario, una vez completa (H4), y exigir que ninguna responda 501 ni 403.

**M2 · El paquete de gastos no está donde dicen las convenciones.**
- Criterio: inconsistencia entre convenciones y código.
- Ruta: `backend/internal/api/gastos-inicio/` (package `gastos`).
- Agente dueño: api-gastos-inicio (o backend-base, si se cambia la convención).
- Qué pasa: CONVENCIONES-backend §2 fija `internal/api/gastos`. El paquete vive en `gastos-inicio/`, así que el import y la línea de `areas.go` no son los que dicen las convenciones.
- Corrección sugerida: mover el paquete a `backend/internal/api/gastos/`, o que backend-base actualice §2 y §14 con la ruta real.

**M3 · Los bloqueos de la semilla no siguen el nombre `S6-<tema>`.**
- Criterio: S6-16.
- Rutas: `docs/estado/bloqueos/semilla-despensas-futuras.md`, `semilla-sucursales-municipio.md` y `semilla-catalogo-alias.md`.
- Agente dueño: semilla.
- Qué pasa: el criterio pide `S6-<tema>.md` para todo supuesto que cambie cifras. Las despensas 2–4 sembradas cambian cifras, y `ls docs/estado/bloqueos/S6-*` no encuentra esos archivos.
- Corrección sugerida: renombrarlos a `S6-semilla-*.md`, actualizar las citas en `semilla.md` y en el código, y agregarlos a `docs/estado/bloqueos/INDICE.md`.

### Bajos

- **B1** · `backend/internal/api/presupuesto/flujo.go:68` (api-presupuesto): `cumulative_cash_cents` se calcula con una resta en la API (`devengado − diferencia`), aunque el motor ya entrega `LiveDebtResult.CashToDateCents`. Sugerencia: tomar el valor del motor y comprobar la identidad (C1: «mismo origen», Q2-65 A).
- **B2** · S6-09 (auth-sistema y backend-base): la hora (03:15) y la retención (30) PROVISIONAL están solo en `S6-respaldos.md` y `S6-valores-operativos-base.md`; `ops/runbook.md` no existe. Sugerencia: copiarlas allí en WF-4, como pide el criterio.
- **B3** · Mensajes 400 de filtros con valores del enum en inglés, p. ej. «El filtro family solo admite household, external o joint.» y «El tipo de ingreso es fixed (quincena) o extra.» (api-gastos-inicio y api-presupuesto). Son valores de la API; si la PWA llega a mostrarlos, conviene nombrarlos en español.
- **B4** · PREGUNTA para el dueño de C1 y para Axel: `auth.password_min_length` tiene `minimum: 0` en C1-datos.md l.804, y con eso se puede bajar de los 12 caracteres de K6. auth-sistema ya lo propuso; no es defecto de los agentes revisados.

## Estado de los anteriores

Es la primera ronda de revisión del backend core: no hay hallazgos anteriores que seguir. Los reportes de los agentes revisados ya mencionaban el gofmt de `internal/reportes` y `areas.go`. **No** mencionaban H1, H3 ni la falla de `go test ./...`: api-presupuesto y api-gastos-inicio reportan «go test ./... pasa», pero la prueba de contrato, que llegó después, ya no pasa.
