# api-tablero · Manejadores del tablero de Inicio (C1 v1.5, K54) y HTTP/2 · WF-3c · 12-oct-2026

Agente `api-tablero` (Opus). Al empezar, el paquete `backend/internal/api/tablero/` ya existía (lo escribió `correccion-especificacion-r2` para cerrar el hallazgo A-1). Por eso este trabajo **completa y prueba con datos** esas 5 operaciones, no las reescribe. Además llena los 8 campos nuevos de `DashboardMonthVsPlan` (es el hallazgo B-R3-2 en la parte que me toca) y activa HTTP/2.

## Hecho
- **HECHO · HTTP/2.**
  - `cmd/finanzas/tls.go`: `loadTLS` ahora pone `NextProtos: ["h2", "http/1.1"]` (variable `alpnProtos`).
  - `http.Server` ya traía el manejador de h2, porque su `TLSConfig` es nil. Lo que faltaba era que el listener TLS anunciara h2 por ALPN.
  - `parseAddrs` y `checkLAN` no cambian (S6-14).
  - Prueba nueva `TestEscucha_NegociaHTTP2`, sobre el mismo `serveAll` de producción:
    - con un cliente TLS crudo, ALPN negocia `h2`;
    - un cliente HTTP con h2 llega como `HTTP/2.0`;
    - un cliente que solo habla HTTP/1.1 llega como `HTTP/1.1`.
  - Comprobé que la prueba **falla** si se quita `NextProtos`.
- **HECHO · K-1 «Gasto del mes»** (`getDashboardMonthVsPlan`, en `api/gastos-inicio/inicio.go`).
  - La función nueva `addMonthComparison` llena los 8 campos opcionales de C1 v1.5: `cash_cents`, `previous_month`, `previous_actual_cents`, `previous_cash_cents`, `actual_difference_cents`, `cash_difference_cents`, `categories_warning_count` y `categories_exceeded_count`.
  - Todo sale de `mf.BudgetMonth` y `mf.PlanVsActual` (motor S2). Aquí solo se resta y se cuenta.
- **HECHO · Regla «sin datos» compartida.** Está en `dominio.MonthHasData` (`dominio/tablero_datos.go`).
  - La usan el mes anterior de M-01 y el devengado y el efectivo de M-21.
  - Es PROVISIONAL, con el bloqueo nuevo `S6-tablero-sin-datos`.
  - Antes la regla estaba en línea dentro de `tablero.go` y citaba un punto que no existe en el bloqueo C1-tablero.
- **HECHO · Escala del barómetro.** `scaleMaxPct(limit)` queda aparte, con su prueba (40→60, 45→70, 32→50, 100→150). Si C1 agrega `scale_max_pct` a M-21 (M-R3-1), basta una línea.
- **HECHO · Pruebas con datos** (`api/tablero/tablero_datos_test.go`). El fixture tiene:
  - sueldo capturado de julio a septiembre y esperado en octubre;
  - un extra en octubre;
  - una compra a 3 pagos repartida entre el hogar y Karen;
  - planes y gastos de contado en 4 categorías.
  
  Con él las pruebas comprueban:
  - **Contra el motor:** que `payments_cents` y `payments_total_cents` sean iguales a `getCreditTotals` (S3), y que el devengado, el efectivo y el plan de octubre en M-21 sean iguales a `getDashboardMonthVsPlan` (S2).
  - **Contra la cuenta a mano de RN-30 con ROUND_HALF_UP:**
    - 32.5 % → 33, warning; 37.5 % → 38; 16.25 % → 16;
    - margen 1,500.00; con umbrales 20/32, −100.00.
  - **Contra RN-31:**
    - 110.5 % sube a 111 y cuenta para raise;
    - 110.4 % baja a 110 y no cuenta;
    - un mes con plan y sin gasto cuenta como 0 % y da cut;
    - un mes sin plan corta la racha.
  - **El color sale de `app_setting`:** cambiar a 35/45 da normal y escala 70; 20/32 da exceeded.
  - **Serie de 12 meses** (de nov-2025 a oct-2026, y de nov-2026 a oct-2027 como proyección).
  - **Ajustes:** `audit_log` del cambio; un 422 no guarda nada; un cuerpo malformado da 400.

## Tabla operación → prueba (S6-03)
| Operación | Pruebas (`go test -v`) | Paquete |
|---|---|---|
| getDashboardDebtCapacity | `TestOp_getDashboardDebtCapacity`, `TestOp_getDashboardDebtCapacity_ConDatos`, `TestDebtStatus_Fronteras`, `TestScaleMaxPct` | api/tablero |
| getDashboardMonthlyOverview | `TestOp_getDashboardMonthlyOverview`, `TestOp_getDashboardMonthlyOverview_Serie12` | api/tablero |
| getBudgetSignals | `TestOp_getBudgetSignals`, `TestOp_getBudgetSignals_ConDatos` | api/tablero |
| getDashboardSettings | `TestOp_getDashboardSettings` | api/tablero |
| updateDashboardSettings | `TestOp_updateDashboardSettings`, `TestOp_updateDashboardSettings_Bitacora` | api/tablero |
| getDashboardMonthVsPlan (campos v1.5) | `TestOp_getDashboardMonthVsPlan_C1v15` (más la previa `TestOp_getDashboardMonthVsPlan`) | api/gastos-inicio |
| HTTP/2 (tls.go) | `TestEscucha_NegociaHTTP2` (S6-14: `TestEscucha_SoloLANYLocal` y `TestEscucha_HTTPSConCertificado` siguen igual) | cmd/finanzas |

## Archivos
- Nuevos:
  - `backend/internal/api/tablero/tablero_datos_test.go`
  - `backend/internal/dominio/tablero_datos.go`
  - `backend/internal/dominio/tablero_datos_test.go`
  - `backend/cmd/finanzas/tls_h2_test.go`
  - `docs/estado/bloqueos/S6-tablero-sin-datos.md`
  - este reporte
- Cambiados:
  - `backend/cmd/finanzas/tls.go`
  - `backend/internal/api/tablero/tablero.go` (escala compartida y regla «sin datos»)
  - `backend/internal/api/tablero/export_test.go`
  - `backend/internal/api/gastos-inicio/inicio.go` (`addMonthComparison`)
  - `backend/internal/api/gastos-inicio/inicio_test.go`
- Sin cambio: `areas.go` (ya traía `tablero.Register`) y `contrato_test.go`.

## Verificación
- **S6-01:** `make verify-wf3-go`:
  - go1.27.2, go.mod 1.27;
  - `go build ./...` y `go vet ./...` sin errores;
  - `gofmt -l .` vacío;
  - `go mod tidy -diff` sin cambios;
  - `go test ./...`: 33 paquetes OK y 5 sin pruebas.
- **S6-02:** `go test ./... -count=1 -cover` en lo tocado:
  - api/tablero 87.6 % (antes 71.2 %);
  - api/gastos-inicio 82.0 %;
  - dominio 87.0 %;
  - cmd/finanzas 26.2 %;
  - internal/api 88.8 %.
  
  `make verify-wf3-cobertura`: K25 OK (dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8). Sin umbral: api 82.0.
- **S6-03:** `make verify-wf3-contrato`:
  - «los 22 archivos coinciden byte a byte con el generador»;
  - «221 operaciones en 176 rutas, 0 diferencias» (216 + 5);
  - `TestAreas_SinOperacionesEn501` «221 de 221»;
  - `TestContrato_Operaciones` «221 de 221 operaciones cumplen»;
  - `TestContrato_Rutas` «sin implementar 0».
- **S6-04:** cada `c.Do` valida la respuesta contra `openapi.yaml`, también los 4xx propios: 400 de formato, 401, 422 de rango, 422 de amarillo ≥ tope y 422 de quincena o semana.
- **S6-14:** `TestEscucha_SoloLANYLocal` PASS y `TestEscucha_HTTPSConCertificado` PASS. Se suma `TestEscucha_NegociaHTTP2` PASS.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0** («Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»).
- **Servidor de desarrollo.** Se reinició con `ops/dev-stop.sh` y `ops/dev-run.sh`, que compila `backend/bin/finanzas`. Las peticiones fueron con curl, la CA de mkcert y la sesión de axel; la contraseña de desarrollo se leyó del archivo sin imprimirse. Todas negocian **HTTP/2** (`%{http_version}` = 2):
  - login → 200 HTTP/2.
  - `GET debt-capacity` → 200 HTTP/2: ingreso 24,500.00 esperado, pagos del hogar 598.82, 2 %, normal, escala 60.
  - `?period=2026-11` → 200: 1,197.64 en 2 pagos, 5 %.
  - `GET monthly-overview` (12 meses) → 200 HTTP/2: octubre con devengado **1,197.64**, efectivo **598.82**, plan 2,000.00 y pagos 598.82. Los meses previos vienen sin cifras.
  - 24 meses → 200 en 0.009 s.
  - `GET budget/signals` → 200 HTTP/2: 3 categorías activas, todas `none` (no hay historia antes de octubre).
  - `GET settings/dashboard` → 200 HTTP/2 (30/40/5/3/10, iniciales).
  - `PATCH` 1/2 → 200; luego warning 45 → 422 `out_of_range`; luego `reset_parameters` → 200. La base de desarrollo quedó en sus valores iniciales.
  - Sin sesión → 401. Con `--http1.1` → 200 HTTP/1.1.
  - Las cifras de octubre coinciden con los ejemplos del oráculo que C1-datos cita en M-01: `periodos[0].devengado.axel` 1,197.64 y `efectivo.axel` 598.82. `payments_cents` (59,882) es igual a `getCreditTotals.month_household_cents` (59,882).

## Supuestos
- El «backend/internal/api/inicio/» del ALCANCE es `api/gastos-inicio/inicio.go`, donde vive `getDashboardMonthVsPlan` (CONVENCIONES §2). Solo toqué esa operación.
- El efectivo `cash_cents` es el `cash_household_cents` del motor tal cual. Con gastos de contado repartidos con Karen, el motor cuenta lo pagado completo: 1,500.00 de efectivo contra 1,000.00 devengados en la prueba. No lo cambié; es la definición de M-11.
- En `updateDashboardSettings`, un tipo incorrecto (`"4"`) responde **400**, como dicen C1 §3 y CONVENCIONES §4. La descripción de `DashboardSettingsUpdate` dice «de otro tipo… 422»: ver la PREGUNTA 2.

## Bloqueos
- **Nuevo:** `docs/estado/bloqueos/S6-tablero-sin-datos.md`, sobre cuándo un mes «no tiene datos».
  - PROVISIONAL (a): devengado y efectivo del hogar los dos en 0.
  - (b): antes del primer mes con algún registro.
  - (c): nunca se omite.
- **Siguen vigentes:** C1-tablero (puntos 1 a 5), U3-tablero-barometro, U3-tablero-senal-presupuesto, U3-tablero-horizonte y U3-tablero-recordatorio.

## Pendientes (fuera de mi alcance; van en solicitudes)
- **`getCreditTotals`** (`api/creditos/creditos.go`): faltan `month_payment_count` y `month_pending_household_cents`.
- **`getIncomeMonth`** (`api/presupuesto/ingresos.go`): faltan `previous_month`, `previous_fortnights_cents`, `previous_extras_cents` y `fortnights_difference_cents`.
- **M-R3-1:** `scale_max_pct` en `DashboardMonthlyOverview`. Va por C1 (c1gen) o por U3 (opción A).
- **`docs/estado/bloqueos/INDICE.md`:** falta listar `S6-tablero-sin-datos`.

## PREGUNTAS y PROPUESTAS para Axel
- **PREGUNTA 1** (bloqueo S6-tablero-sin-datos): cuando un mes no tiene gastos, ¿lo vemos como «—» o como $0.00? Propongo «—» solo antes de que empezaras a capturar (opción b) y $0.00 después.
- **PREGUNTA 2:** si en Ajustes del tablero mandan un texto donde va un número, la API responde «formato incorrecto» (400), igual que en el resto del sistema. La descripción del contrato dice 422. Propongo corregir la descripción en C1 y no la API.
- **PROPUESTA:** en el servidor de desarrollo, el barómetro de octubre marca **2 %** (598.82 de 24,500.00), verde. Con los datos sembrados aún no hay historia para la sugerencia «subir o recortar»: aparece a partir de enero de 2027, con 3 meses cerrados.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 45 acciones de herramienta; unos 165 KB leídos y unos 40 KB escritos (código, pruebas, bloqueo y este reporte). Los temporales de `datos/tmp/api-tablero/` se borraron al terminar.
