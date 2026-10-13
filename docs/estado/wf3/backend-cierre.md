# backend-cierre · WF-3c bis · 12-oct-2026

Agente `backend-cierre` (Opus). HEAD `0197dad`. Aplica las decisiones 1, 2 y 3 del orquestador de WF-3c bis en el backend y el contrato. Cierra H-4 de `REVISION-c3-tablero.md`, M-2 y B-5/B-6 de `REVISION-tablero-backend-r3.md`, M-R3-1 y B-R3-2 (lo del backend) de `REVISION-especificacion-r3.md` y las diferencias 1 y 2 de `u5-tablero.md`. No toqué motores ni cifras de motor, ni el frontend, ni git.

## Hecho
- **HECHO · Deuda viva «hoy» (decisión 1, H-4).** `creditos.LiveDebtAt` (`motor.go`): el motor S3 calcula la deuda (`motor_cred.LiveDebt`) y solo se dejan fuera las compras con fecha posterior al corte, que es la regla de M-12 y de `getLiveDebt` (motor S2). La usan `getCreditTotals.live_debt_cents`, la deuda de cada línea (`listCredits`/`getCredit`), `getCardStatus` y cada cierre de `getDebtTrend`. En el servidor: `getCreditTotals.live_debt_cents` = `getLiveDebt?date=2026-10-12` = **389500**. La línea Kueski también da 389500. La «Deuda viva proyectada» arranca en **292125** al cierre de oct-2026 (2,921.25, como el oráculo) y crece en dic, feb y abr con cada despensa programada. `getLiveDebt?date=2027-04-10` sigue en 629344.
- **HECHO · M-2 (decisión 2).** Apliqué tal cual la solicitud 1 de `correccion-tablero-backend-r2` en `creditos/tablero.go`: `PaymentsByMonth` arma el calendario con todas las líneas y `has_credits` sigue la (a1). Quité el `t.Skip` de `opcion_a_PROVISIONAL` y el subtest pasa. Apliqué también la solicitud 2: el texto de `has_credits` en C1.
- **HECHO · K-3 (M-19, B-R3-2).** `getCreditTotals` llena `month_payment_count` y `month_pending_household_cents`. Falta pagar = por cada pago que el motor no da por pagado, la parte del hogar menos los abonos del hogar, sin bajar de 0 (PROVISIONAL, `S6-deuda-viva-hoy` punto 3). Oct-2026 en desarrollo: 1 pago y 59882 por pagar.
- **HECHO · K-2 (M-18, B-R3-2).** `getIncomeMonth` llena `previous_month` (siempre), `previous_fortnights_cents`, `previous_extras_cents` y `fortnights_difference_cents`. Las tres últimas van ausentes si el mes anterior no tiene quincenas ni extras (`dominio.MonthHasData`, `S6-tablero-sin-datos`). En desarrollo, septiembre no tiene ingreso esperado, así que las tres van ausentes.
- **HECHO · M-R3-1 (decisión 3, opción B).** `scale_max_pct` es ahora entero requerido en `DashboardMonthlyOverview`, con C1 **1.5.1** (`gen2_schemas.py`, `gen4_build.py`, `c1_plantilla.md` con la fila de Cambios, la sección v1.5.1 y M-21, `verify_c1.py` y `README.md` de c1gen). Regeneré con `build.py`. `tablero.go` lo llena con la misma `scaleMaxPct` de B-1 y el servidor da 60.
- **HECHO · bajos.** B-5: el texto de `BudgetSignals` cita `S6-tablero-senal-periodo` (a). B-6: `monthParam` da su propio 400, «…Aquí va un mes con el formato AAAA-MM, por ejemplo 2026-10.», con prueba.
- **HECHO · documentos.** RN-21 suma un párrafo PROVISIONAL sobre el «hoy» (cita `S6-deuda-viva-hoy`). RN-30 aclara activas y archivadas. Hay una nota de corrección en el encabezado de `reglas-negocio.md`. Bloqueo nuevo `S6-deuda-viva-hoy.md`; actualicé `S6-tablero-lineas-archivadas.md` (ya rige la (a); pregunta nueva) y `S6-tablero-sin-datos.md`. `INDICE.md` tiene la sección «WF-3c bis» y la fila de las archivadas al día. `S6-tablero-sin-datos` ya tenía su fila.

## Archivos
- Go: `backend/internal/api/creditos/{tablero.go,motor.go,creditos.go,creditos_test.go}`, `backend/internal/api/presupuesto/{ingresos.go,presupuesto_test.go}` y `backend/internal/api/tablero/{tablero.go,tablero_test.go,tablero_datos_test.go}`.
- Contrato (c1gen): `gen2_schemas.py`, `gen4_build.py`, `c1_plantilla.md`, `verify_c1.py` y `README.md`. Lo regenerado: `docs/contratos/openapi.yaml`, `C1-datos.md`, `extractos/*` (22), `backend/internal/api/routes.go` y `types_gen.go`.
- Documentos: `docs/contratos/reglas-negocio.md`, `docs/estado/bloqueos/{S6-deuda-viva-hoy.md (nuevo),S6-tablero-lineas-archivadas.md,S6-tablero-sin-datos.md,INDICE.md}` y este reporte.

## Verificación
- **S6-01:** `make verify-wf3-go` EXIT 0 (go1.27.2; build, vet, gofmt y tidy limpios; 33 paquetes OK y 5 sin pruebas).
- **S6-02:** `make verify-wf3-cobertura` EXIT 0. Con K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Con `go test -cover`: api/tablero 87.8, api/creditos 80.9 (antes 80.4), api/presupuesto 77.7 e internal/api 88.8.
- **S6-03:** `make verify-wf3-contrato` EXIT 0: «22 archivos coinciden byte a byte» y «221 operaciones en 176 rutas, 0 diferencias». `TestContrato_Operaciones`: «221 de 221 operaciones cumplen».
- **S6-04:** pruebas nuevas `TestOp_getCreditTotals_DeudaVivaHoy` (compra a futuro, igualdad con `getLiveDebt` hoy, 30-nov y 31-dic, abono parcial del hogar y de Karen), `TestOp_getIncomeMonth` ampliada (mes anterior con datos y sin datos) y `TestOp_getDashboardDebtCapacity_PeriodoMalFormado` (B-6). `scale_max_pct` = 60 en `TestOp_getDashboardMonthlyOverview`. El subtest `opcion_a_PROVISIONAL` corre sin SKIP.
- **S6-14 / S6-15:** el servidor se reinició con `ops/dev-stop.sh`, `go build` y `ops/dev-run.sh`, y escucha por HTTP/2. La sesión de axel usó la contraseña leída del archivo sin imprimirla; al final, logout 204 y `auth/me` 401. `finanzas.log` tiene 0 coincidencias de password, contraseña, cookie o token.
- **Cifras contra el servidor** (sonda propia, de solo lectura):
  - `getCreditTotals.live_debt_cents` 389500 = `getLiveDebt?date=2026-10-12` 389500.
  - debt-capacity de octubre: 59882 · 97375 · 1, ratios 2 y 4, normal, margen 920118, escala 60.
  - Serie de jul-2026 a jun-2027: hogar 0·0·0·59882·119764·148257·176750·176752·176750·176752·176750·88377, ratios 2·5·6·7·7·7·7·7·4. Son las mismas que en la ronda 3.
  - **P1–P6** (oct-2026 a mar-2027): los pagos del hogar coinciden con `periodos[0..5].efectivo.axel`. El devengado de oct (119764), el efectivo (59882), el plan (200000) y el 60 % no cambian.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0, «Verificación ampliada OK». S1-16 cita `S6-deuda-viva-hoy`; S1-03 OK.

## Supuestos
- PROPUESTA (interpretación de la decisión 1, sin respuesta del orquestador): «arranca en la deuda de hoy» se lee como que cada cierre de `getDebtTrend` cuenta solo las compras hechas hasta ese cierre. El primer punto es el cierre de octubre (292125), no el 12-oct (389500): C1 M-17 dice «al cierre», y u5-tablero da 2,921.25 como cifra del oráculo.
- La misma regla se aplica a la deuda de cada línea y a `getCardStatus`, para que K-3, P-17 y la línea den una sola cifra.

## Bloqueos
- Nuevo: `S6-deuda-viva-hoy` (3 puntos, (a) PROVISIONAL).
- Pregunta nueva en `S6-tablero-lineas-archivadas`: si `getCreditTotals` también cuenta las archivadas en las cifras del mes.

## Solicitudes fuera de alcance
1. `backend/internal/api/gastos-inicio/vistas.go:235` (`deriveCredits`): cambiar `mc.LiveDebt([]mc.CreditState{ci.state(c)}, cdateOf(t))` por `creditos.LiveDebtAt(...)` con los mismos argumentos, importando `finanzas-hogar/internal/api/creditos`. Si no, la vista de BD de créditos sigue contando las compras futuras.
2. Frontend: en `frontend/src/pantallas/inicio/Tablero.tsx`, que G-6 tome el eje Y de `scale_max_pct` de `getDashboardMonthlyOverview`. En `frontend/src/pantallas/inicio/mocks.ts:31`, agregar `scale_max_pct: 60`.
3. `docs/design/U3-dashboard.md` G-6 y §2.6: citar `DashboardMonthlyOverview.scale_max_pct` (C1 v1.5.1).
4. `README.md:88` (raíz): «C1 v1.5.0» → «C1 v1.5.1».
5. B-1 y B-2 siguen abiertos (el texto 422 de `SETW` es común a muchas operaciones). Lo decide quien sea dueño de `gen3_ops.py`/`SETW`.

## Pendientes
- C3-09 y las capturas U5 tienen que volver a correr con el servidor nuevo (lo hace el agente que corresponda).

## Propuestas para Axel
- Puerta 4: confirmar `S6-deuda-viva-hoy` (a). Así una despensa programada se suma a la deuda el día de su compra.
- Puerta 4: decidir si K-3 cuenta las líneas archivadas igual que el barómetro (`S6-tablero-lineas-archivadas`).

## Paquetes instalados
Ninguno.

## Consumo aproximado
57 acciones, unos 150 KB leídos y unos 30 KB escritos a mano (sin contar lo que regeneró `build.py`).
