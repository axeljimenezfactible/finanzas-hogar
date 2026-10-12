# correccion-especificacion-r2 · Corrección de la ronda 2 de la especificación del tablero (WF-3c, K54) · 12-oct-2026

Agente `correccion-especificacion-r2` (Opus). Corrige los dos hallazgos alto y medio de `docs/estado/wf3/REVISION-especificacion-r2.md`: A-1 (H-1 reabierto) y M-N1. Los bajos (B-N1 a B-N5) no entran en esta ronda. No hice `git add` ni `git commit`, no mandé Telegram, no usé sudo y no instalé paquetes. Ninguna cifra, id de criterio, token, migración ni operación de C1 cambia.

## Hecho
- **A-1 (alta) · HECHO.** La corrección r1 solo había puesto al día el reporte. Esta vez se registran las 5 operaciones de C1 v1.5:
  - Paquete nuevo `backend/internal/api/tablero/` con `Register`:
    - `getDashboardDebtCapacity` (B-1, M-20, RN-30).
    - `getDashboardMonthlyOverview` (G-1 y G-6, M-21).
    - `getBudgetSignals` (G-3, M-22, RN-31).
    - `getDashboardSettings` y `updateDashboardSettings` (K44, las 5 claves de la 00005).
  - Las cifras salen de los motores, sin recalcular nada:
    - Del motor S2 (`motor_fin`, entrada de `dominio.LoadFinEngine`): el ingreso (`FortnightFlow` de las dos quincenas del mes: capturadas + esperadas, extras aparte), el devengado (`PlanVsActual`, el real de M-01), el efectivo (`BudgetMonth.CashHouseholdCents`, M-11), el plan y el % usado por categoría (`BudgetMonth`, M-09).
    - Del motor S3: los pagos de crédito, con `creditos.PaymentsByMonth`, archivo nuevo `backend/internal/api/creditos/tablero.go`. Hace la misma cuenta que `getCreditTotals`: calendario unificado de las líneas activas y parte del hogar = personas no externas. Solo agrupa por mes.
    - El paquete solo divide con `dominio.PercentOf` y `MulDivRound` (ROUND_HALF_UP) y compara contra los umbrales vigentes de `app_setting`.
  - Validaciones que pidió el revisor:
    - El PATCH responde 422 en tres casos: `debt.capacity_warning_pct` ≥ `debt.capacity_limit_pct` con los valores que resultarían (en la misma transacción, sin guardar nada), la misma clave en `parameters` y en `reset_parameters`, o un valor fuera de rango (`dominio.ApplyParams` y disparadores de la 00005).
    - `monthly-overview` responde 422 si `from` > `to` o si la ventana pasa de 24 meses. Sin `from`/`to` o con un mes mal escrito responde 400 (formato, C1 §3).
    - `period` con quincena, semana o año responde 422 en `debt-capacity` y `signals`.
  - Registro:
    - `cmd/finanzas/areas.go` suma `tablero.Register`.
    - Las dos listas de áreas de las pruebas (`internal/api/contrato_test.go` e `internal/auth/auth_test.go`) también lo suman.
    - `internal/api/contrato_datos_test.go` suma los datos de prueba de las 2 operaciones que el generador no sabe llenar: `settingsEcho` para el PATCH y `from`/`to` = 2026-10/2026-12.
- **M-N1 (media) · HECHO.** Elegí la opción «ampliar la fila M-01», no crear una M-23. U3 ya decía «M-01 ampliada», `getDashboardMonthVsPlan` es la misma operación y así C1-09 («Cubre M-01…M-22») sigue valiendo sin ids nuevos.
  - En `backend/internal/api/c1gen/c1_plantilla.md` §9, la fila M-01 ahora define los 8 campos:
    - `cash_cents` = `cash_household_cents` de M-11.
    - Mes anterior: `previous_month` es el mes calendario anterior. `previous_actual_cents` y `previous_cash_cents` faltan sin datos.
    - Diferencias con signo: `actual_difference_cents` y `cash_difference_cents`.
    - Conteos: `categories_warning_count` y `categories_exceeded_count`, con el `status` de M-09.
    - También cambian el periodo (mes y mes anterior) y las tablas (más `scheduled_payment` y `scheduled_payment_share`).
    - Ejemplo del oráculo: oct-2026 con efectivo **598.82** (`periodos[0].efectivo.axel`, R-13); nov-2026 con real **1,197.64** y efectivo **1,197.64** (`periodos[1]`), o sea diferencia **0.00** en lo gastado y **598.82** en efectivo.
  - La lista de cambios de C1 v1.5 apunta a esa fila.
  - Regeneré con `build.py`. Solo cambia `docs/contratos/C1-datos.md` (línea 404 y la lista de cambios). `openapi.yaml` sale idéntico con `cmp`, así que los extractos, `routes.go` y `types_gen.go` no cambian.
  - En `docs/design/U3-dashboard.md`:
    - K-1 cita la fila M-01 ampliada de C1-datos §9 (C1-09).
    - Línea de cambios «Corrección r2».
    - §15 marca HECHO el backend del tablero.
  - Ninguna cifra ni fórmula de K-1 cambia.
- Notas «Corrección r2» en `docs/estado/wf3/s1-c1-v15.md` y `docs/estado/wf3/u3-v2.md`.

## Archivos
- Nuevos:
  - `/home/deepsec/finanzas-hogar/backend/internal/api/tablero/tablero.go`
  - `/home/deepsec/finanzas-hogar/backend/internal/api/tablero/tablero_test.go`
  - `/home/deepsec/finanzas-hogar/backend/internal/api/tablero/export_test.go`
  - `/home/deepsec/finanzas-hogar/backend/internal/api/creditos/tablero.go`
- Una línea y su import en cada uno:
  - `/home/deepsec/finanzas-hogar/backend/cmd/finanzas/areas.go`
  - `/home/deepsec/finanzas-hogar/backend/internal/api/contrato_test.go`
  - `/home/deepsec/finanzas-hogar/backend/internal/auth/auth_test.go`
- `/home/deepsec/finanzas-hogar/backend/internal/api/contrato_datos_test.go` (3 líneas de datos de prueba)
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/c1_plantilla.md` y `/home/deepsec/finanzas-hogar/docs/contratos/C1-datos.md` (este, regenerado)
- `/home/deepsec/finanzas-hogar/docs/design/U3-dashboard.md`
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/s1-c1-v15.md`, `/home/deepsec/finanzas-hogar/docs/estado/wf3/u3-v2.md` y este reporte

## Verificación
- **Punto (f) · S6-03 · C1-02 · `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0** («Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»).
  - C1: 22 archivos idénticos al generador, TODO OK.
  - S1: RESUMEN TODO OK.
  - U1: TODO OK.
  - Criterios: 261, 0 citados inexistentes.
  - K25 OK.
  - Oráculos S2 a S5: 0 diferencias.
- **Las 4 pruebas de A-1 pasan:**
  - `TestAreas_SinOperacionesEn501`: «operaciones registradas: 221 de 221».
  - `TestContrato_Operaciones`: «S6-03: 221 de 221 operaciones cumplen».
  - `TestContrato_Rutas`: «contrato 221 operaciones, servidor 221, sobrantes 0, faltantes 0, sin implementar 0».
  - `TestPermisos_MismasOperacionesAmbasCuentas`: PASS.
- **S6-02 / S6-03 · `go test ./internal/api/tablero/`: ok, cobertura 71.2 %** (`api` no tiene umbral K25). Hay una prueba por operación:
  - `TestOp_getDashboardDebtCapacity`: 200 con 30/40/60, 422 con quincena, 400 con un periodo mal escrito y 401.
  - `TestOp_getDashboardMonthlyOverview`: 4 meses, nov y dic marcados como proyección y sin devengado ni efectivo; 422 con `from` > `to` y con 25 meses; 200 con 24 meses; 400 sin `from`; 401.
  - `TestOp_getBudgetSignals`: 3 meses jul a sep, sin el mes pedido; 422 con semana; 401.
  - `TestOp_getDashboardSettings`: 30/40/5/3/10.
  - `TestOp_updateDashboardSettings`: 200; 422 con amarillo = tope y con amarillo > tope; 422 al cambiar y restablecer a la vez; 422 fuera de rango; restablecer vuelve a 40 y 3; 400 vacío; 401.
  - `TestDebtStatus_Fronteras` (RN-30): 30 normal, 31 y 40 warning, 41 exceeded.
  - Todas las respuestas se validan contra `openapi.yaml`.
- **`make -k verify-wf3-go verify-wf3-cobertura verify-wf3-contrato verify-wf3-secretos`: EXIT 0.** 33 paquetes OK; «S6-03: 221 operaciones en 176 rutas, 0 diferencias».
- **C1-09 · U3-02:**
  - `grep -n "^| M-01 " docs/contratos/C1-datos.md` da la línea 404, con los 8 campos y el ejemplo.
  - Comprobé las cifras con `python3 -I` sobre `oraculo/model.json`: `periodos[0]` 1197.64 / 598.82 y `periodos[1]` 1197.64 / 1197.64.
  - `build.py --check`: «los 22 archivos coinciden byte a byte».
  - `cmp openapi.yaml` antes y después: idéntico.
- `gofmt -l` vacío y `go vet` limpio en los paquetes tocados.

## Supuestos
1. Escribí los manejadores aquí, no en un agente `api-tablero` aparte. La tarea pide corregir A-1 y que `make verify-all` dé EXIT 0, y en la ronda r1 el revisor reabrió el hallazgo porque solo se había puesto al día el reporte. Si `api-tablero` corre después, encuentra las 5 operaciones con dueño: `reg.Handle` entra en pánico si alguien las registra dos veces.
2. `creditos/tablero.go` es un archivo nuevo dentro del paquete de `api-creditos`. Solo agrupa la salida de S3 que ya usa `getCreditTotals`, para que `payments_cents` sea por construcción la misma cifra que `month_household_cents` (lo pide C1). No toca nada de lo que ya existía.
3. «Mes sin datos» en `monthly-overview`: un mes que no es proyección y tiene devengado y efectivo en 0 lleva `accrued_cents` y `cash_cents` ausentes. Es PROVISIONAL dentro del bloqueo abierto `C1-tablero`: C1 dice «ausente en meses sin datos» y no define «sin datos».
4. Para RN-31 se usa el entero `used_pct` del motor y la comparación estricta (110 no pasa de 110, como el ejemplo ilustrativo). Un plan en 0 no tiene % usado y corta la racha (Q2-20 A).
5. `scale_max_pct` = el múltiplo de 10 igual o mayor que 1.5 × tope, en enteros: (3 × tope + 19) ÷ 20 × 10. Con 40 da 60.

## Bloqueos
No abrí ninguno nuevo. Siguen abiertos `C1-tablero` y los cinco `U3-tablero-*`.

## Pendientes
- Los 14 campos opcionales nuevos de C1 v1.5 (los 8 de `DashboardMonthVsPlan`, 4 de `IncomeMonth` y 2 de `CreditTotals`) no los llenan todavía los manejadores de `gastos-inicio`, `presupuesto` y `creditos`. Son opcionales, así que el contrato pasa, y la interfaz cae a su regla de §2. Le tocan a sus dueños o a `api-tablero`. Ver solicitudes.
- Bajos B-N1 a B-N5: fuera de esta ronda. B-N1 (U3 «M-01 a M-17» y la viñeta de U1 en §15), B-N2 (G-3 con valores fijos), B-N3 (M-19 en K-3), B-N4 («esquema 1.4» en el Makefile y en `s1_verifica.py`) y B-N5 (punta de la aguja).
- Ninguna prueba con un crédito real compara las cifras de `debt-capacity` con `getCreditTotals`. La igualdad hoy se garantiza porque las dos usan el mismo código de S3.

## Propuestas para Axel
1. **PROPUESTA** (`C1-tablero`): fijar qué es «mes sin datos» en el resumen mensual. Lo implementado es «devengado y efectivo en 0». La alternativa es «antes del primer movimiento registrado».

## Paquetes instalados
Ninguno.

## Consumo aproximado
58 acciones. Unos 175 KB leídos, por secciones (`openapi.yaml` solo por extractos; C1-datos y U3 con `grep`/`sed`). Unos 38 KB escritos a mano; `C1-datos.md` sale del generador.
