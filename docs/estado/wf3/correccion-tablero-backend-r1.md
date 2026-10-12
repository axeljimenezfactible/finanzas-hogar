# correccion-tablero-backend-r1 · Corrección de M-1 y M-2 de la revisión del tablero (backend) · WF-3c · 12-oct-2026

Agente `correccion-tablero-backend-r1` (Opus). Corrige los dos hallazgos medios de `docs/estado/wf3/REVISION-tablero-backend-r1.md` (no había altos). Los bajos B-1 a B-4 no eran parte del encargo y no se tocaron. Alcance: los archivos de `api-tablero` que nombran los hallazgos; `backend/internal/api/creditos/` **no** es de ese alcance (`runs/wf3c-tablero.js`, ALCANCE de api-tablero), así que la parte de M-2 que vive ahí va como solicitud fuera de alcance.

## Hecho
- **HECHO · M-1 (RN-31: solo meses cerrados).** `getBudgetSignals` (`backend/internal/api/tablero/tablero.go`, líneas 284-306):
  - la ventana sigue siendo los `budget.signal_months` meses anteriores a `period`, pero un mes cuenta solo si es **anterior al mes en curso** (`closed := monthIndex(mm) < current`);
  - el mes en curso y los de proyección no salen en `months` y la racha no se completa (`len(cells) == months`), así que la señal es `none`;
  - la fila de cada categoría sigue saliendo aunque no haya meses cerrados en la ventana (`months` vacío);
  - con `period` = mes en curso (el caso por omisión, el que pide la PWA) no cambia nada.
  - PROVISIONAL, bloqueo nuevo `docs/estado/bloqueos/S6-tablero-senal-periodo.md`, opción (a). Las otras: (b) 422 a un `period` posterior al mes siguiente y (c) recorrer la ventana a los últimos meses cerrados.
- **HECHO · M-1, prueba.** `TestOp_getBudgetSignals_PeriodoFuturo` (`tablero_datos_test.go`, línea 326):
  - `period=2026-11` (mes siguiente): mira solo 2026-08 y 2026-09; las 4 filas dan `none`. Sin el corte, Transporte daba `cut` con el 0 % de octubre en curso;
  - `period=2027-01` y `2027-02` (futuro): 4 filas, `months` vacío, `none`. Antes daban `cut` en 3 categorías;
  - `period=2026-10`: sigue mirando jul–sep (Despensa `raise`, Transporte `cut`).
  - Comprobé que la prueba **falla** sin el corte (overlay con `closed` siempre true): 11 errores, entre ellos «2027-01, categoría 1: cut con meses "2026-10 2026-11 2026-12 "».
- **HECHO · M-2, bloqueo PROVISIONAL.** Nuevo `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`, con tres puntos:
  - punto 1: (a) PROVISIONAL, todas las líneas, como dice RN-30 literal y como cuenta el efectivo de S2; (b) solo las activas, corrigiendo RN-30, M-20 y M-21;
  - punto 2: `has_credits` con (a). (a1) PROVISIONAL: true si hay una línea activa o algún pago de una archivada en el mes, para que siga valiendo «con false, pagos = 0». (a2) solo líneas activas;
  - punto 3: `getCreditTotals` no cambia (PROVISIONAL). Solo diferiría del tablero si una línea archivada tiene un pago en el mes en curso.
- **HECHO · M-2, prueba.** `TestOp_Dashboard_LineaArchivada` (`tablero_datos_test.go`, línea 368). Usa el fixture y `UPDATE credit SET is_active = 0`, y compara:
  - octubre en M-21 contra S2: `cash_cents` sigue en 700,000; `credit_payments_cents` sigue en 650,000 = efectivo − 500.00 de contado del mes; total 750,000, 33 % y warning;
  - barómetro: octubre 650,000 en 2 pagos y `has_credits` true; noviembre 325,000, 16 % y true; diciembre 0, 0 % y false.
  - **Queda en SKIP** con la referencia a la solicitud y al bloqueo, porque el cambio que la hace pasar es de `api/creditos/tablero.go` (fuera de alcance). Con `go test -overlay`:
    - sin la solicitud y sin el SKIP **falla**: «pagos 0 (antes 650000; parte de crédito del efectivo S2 650000), total 0, 0 %, normal», y `has_credits` false en octubre y noviembre;
    - con la solicitud y sin el SKIP **pasa**, y `api`, `api/creditos` y `api/tablero` siguen OK.
- **HECHO · referencia en el manejador.** Una línea de comentario en `getDashboardDebtCapacity` (`tablero.go:148`) cita el bloqueo S6-tablero-lineas-archivadas.
- **PREGUNTA (sin aplicar en código):** la opción PROVISIONAL de M-2 queda escrita y probada, pero **no aplicada** hasta que se acepte la solicitud 1. Mientras tanto, el tablero sigue contando solo las líneas activas.

## Archivos
- Cambiados:
  - `backend/internal/api/tablero/tablero.go`: M-1 en `getBudgetSignals` y el comentario de M-2.
  - `backend/internal/api/tablero/tablero_datos_test.go`: 2 pruebas nuevas al final.
- Nuevos:
  - `docs/estado/bloqueos/S6-tablero-senal-periodo.md`
  - `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`
  - este reporte.
- Sin cambio: `backend/internal/api/creditos/` (solicitud 1), C1, reglas de negocio y migraciones.

## Verificación
- **S6-01:** `make verify-wf3-go` EXIT 0:
  - go1.27.2, go.mod 1.27;
  - build y vet sin errores; `gofmt -l .` vacío; `go mod tidy -diff` sin cambios;
  - `go test ./...`: 33 paquetes OK y 5 sin pruebas.
- **S6-02:**
  - `make verify-wf3-cobertura` EXIT 0. K25 OK, entre ellos precios 88.8 %. Sin umbral: api 82.0, auth 83.6, jobs 93.0, reportes 77.3 y telegram 80.0;
  - `go test -cover ./internal/api/tablero/`: **87.8 %** (antes 87.6 %).
- **S6-03:** `make verify-wf3-contrato` EXIT 0:
  - «los 22 archivos coinciden byte a byte con el generador»;
  - «221 operaciones en 176 rutas, 0 diferencias»;
  - `TestContrato_Operaciones`: «221 de 221 operaciones cumplen».
- **S6-04:** las pruebas nuevas usan `c.Do`, que valida cada respuesta contra `openapi.yaml`. `months: []` cumple el esquema (no tiene `minItems`).
- **S6-14:** sin cambios en `tls.go`. Tras reiniciar, `ss -ltnp` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`, y curl negocia HTTP/2.
- **S6-15:** identificadores en inglés (`closed`, `current`); mensajes y comentarios en español de México. `datos/dev/run/finanzas.log` tiene 0 coincidencias de password, contraseña, cookie o token.
- **S2 y S3 (cifras):** no cambian. Servidor de desarrollo reiniciado con `ops/dev-stop.sh` y `ops/dev-run.sh` (el binario de las 15:52 es posterior a los fuentes). Sesión de axel por HTTP/2; la contraseña se leyó del archivo sin imprimirse; logout 204 y luego 401.
  - `debt-capacity`: 59,882 · 97,375 · 2 % · normal · `has_credits` true;
  - `monthly-overview` de oct a dic: pagos 59,882 · 119,764 · 148,257, con Karen 97,375 · 194,750 · 223,243; devengado de octubre 119,764 y efectivo 59,882. Son las del oráculo de la revisión.
- **M-1 en vivo** (`GET /api/v1/budget/signals`, 200 HTTP/2, 3 categorías):
  - sin `period` y con `2026-10`: jul–sep, `none`, como antes;
  - `2026-11`: solo 2026-08 y 2026-09, ya sin octubre;
  - `2027-01` y `2027-02`: `months` vacío y `none`. Antes, 2027-01 daba `cut` en la categoría 1.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:** EXIT 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## Supuestos
- `api/creditos/` está fuera del alcance de api-tablero (`runs/wf3c-tablero.js`), así que el cambio de `PaymentsByMonth` va como solicitud y no se aplicó. El encargo lo prevé.
- Para M-1, «mes cerrado» = anterior al mes calendario en curso del reloj de la API (`dominio.Today`). Es la definición que ya usa `is_projection` en `getDashboardMonthlyOverview`.
- La fila de cada categoría se conserva aunque no haya meses cerrados («una fila por categoría activa», C1).

## Bloqueos
- **Nuevo:** `S6-tablero-senal-periodo`. PROVISIONAL (a): los meses no cerrados no cuentan. Otras: (b) 422 y (c) ventana recorrida.
- **Nuevo:** `S6-tablero-lineas-archivadas`. PROVISIONAL (a): todas las líneas; (a1): `has_credits` = activa o con pago en el mes; `getCreditTotals` sin cambio. Otra: (b) solo las activas. Pendiente de la solicitud 1 para quedar en código.
- **Siguen vigentes:** C1-tablero, S6-tablero-sin-datos y U3-tablero-*.

## Solicitudes fuera de alcance
1. **`backend/internal/api/creditos/tablero.go`**, `PaymentsByMonth`, para aplicar la opción PROVISIONAL (a) y (a1) de S6-tablero-lineas-archivadas:
   - cambiar `states := s.states(nil)` por un ciclo sobre `s.credits` que agregue `s.state(c)` de **todas** las líneas y lleve `hasCredits = hasCredits || c.IsActive`;
   - si `UnifiedCalendar` falla, devolver `hasCredits` en lugar de `true`;
   - después del calendario, `hasCredits = hasCredits || len(cal.Entries) > 0`;
   - devolver `out, hasCredits, nil` en lugar de `out, true, nil`;
   - actualizar los comentarios: «con todas las líneas, activas o archivadas (RN-30)…, PROVISIONAL, bloqueo S6-tablero-lineas-archivadas».
   
   Después, quitar la línea `t.Skip("PENDIENTE: …")` de `TestOp_Dashboard_LineaArchivada` (`backend/internal/api/tablero/tablero_datos_test.go:369`). Ya probado con `go test -overlay`: pasa, junto con `api`, `api/creditos` y `api/tablero`, y `go vet` sin errores. `motor.go` (`states`) y `getCreditTotals` no cambian.
2. **`backend/internal/api/c1gen/`** (texto de `DashboardDebtCapacity.has_credits`, si se mantiene la opción a1). Hoy dice «true si hay al menos una línea de crédito activa; …». Debe decir «true si hay al menos una línea de crédito activa o algún pago de una línea archivada con fecha en el mes (PROVISIONAL, bloqueo S6-tablero-lineas-archivadas); …». Regenerar los extractos con `build.py`.
3. **`docs/estado/bloqueos/INDICE.md`**, tabla de S6. Agregar dos filas:
   - `| \`S6-tablero-senal-periodo.md\` | \`correccion-tablero-backend-r1\` | Qué meses mira la sugerencia de presupuesto con un period posterior al mes en curso | sí | puerta 4 |`;
   - `| \`S6-tablero-lineas-archivadas.md\` | \`correccion-tablero-backend-r1\` | Si los pagos de una línea archivada cuentan en el barómetro y la serie del tablero, y qué es has_credits | sí | puerta 4 |`.
   
   Esto completa B-3, que sigue sin listar S6-tablero-sin-datos, C1-tablero ni U3-tablero-*.

## Pendientes
- Aplicar la solicitud 1 y quitar el SKIP. Mientras no se haga, M-2 queda con bloqueo y prueba, pero el código sigue con la opción (b).
- B-1 a B-4 de la revisión: no eran parte del encargo.

## Propuestas para Axel
- **PROPUESTA (S6-tablero-senal-periodo):** dejar la (a). La sugerencia solo debe salir de meses que ya terminaron; si pides un mes futuro, no hay sugerencia.
- **PROPUESTA (S6-tablero-lineas-archivadas):** dejar la (a). Archivar una tarjeta liquidada no debería borrar sus pagos del barómetro de meses pasados, porque el efectivo del tablero sí los cuenta.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 33 acciones, unos 150 KB leídos y unos 25 KB escritos. Los temporales en `datos/tmp/correccion-tablero-backend-r1/` (overlays, diff y log de verify-all) se borraron.
