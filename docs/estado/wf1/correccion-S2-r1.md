# correccion-S2-r1 · Reporte final (WF-1, 11-oct-2026)

Corrige los 3 hallazgos medios de `docs/estado/wf1/REVISION-S2-r1.md` (M1, M2 y M3). No toqué los bajos (B1 y B2), que no estaban en la tarea, ni reabrí decisiones de §A–§K. Ninguna cifra canon cambia y no ajusté ningún número a mano.

## Hecho
- **M1 · Ventanas (S2-03, S2-16).** Abrí el bloqueo `docs/estado/bloqueos/S2-ventanas.md` con 3 opciones para lo anterior a la primera ventana y 3 para los huecos. Mientras se decide, las dos implementaciones aplican la misma regla conservadora, marcada PROVISIONAL con cita al bloqueo: un **costo o un pago** anterior a la primera ventana, o en un hueco entre ventanas, da `invalid_argument`. Lo posterior a la última ventana se sigue admitiendo, porque no entra en ningún periodo ni en ninguna diferencia.
  - Motor: `uncovered` y `(*ledger).checkCovered` en `ledger.go`. Se validan antes de `windowed`. La rama `w0 < 0` de `windowed` queda solo para lo posterior a la última ventana.
  - Oráculo: `sin_cubrir` en `fin_devengo.py`, llamada al inicio de `resumen` en modo ventana. Antes revisaba solo costos anteriores; ahora también revisa pagos y huecos.
  - El caso del revisor, semilla con `{windows:[P2]}`: el motor daba 0 / 194750 / −292125 y ahora da `invalid_argument`, igual que el oráculo.
  - Con esto el supuesto 2 de `motor-S2-go.md` queda sustituido por el bloqueo.
- **M2 · `accrual_months` e `interval_count` (S2-03, S2-16).** El motor ya distingue «ausente» de 0:
  - `Expense.AccrualMonths`, `CreditPurchase.AccrualMonths`, `PlannedItem.IntervalCount` y `PlannedItem.AccrualMonths` pasan a `*int`.
  - `countOK` acepta ausente o null (vale 1) y rechaza 0, negativos y más de 24 con `invalid_argument`. `count1` da el valor en los 5 lugares donde se usa (`ledger.go` ×3, `projection.go`, `plan.go`).
  - Es lo mismo que ya hacía el oráculo (`fin_comun.py`: null → 1, rango 1–24), así que el oráculo no cambió. El supuesto 1 de `motor-S2-go.md` («0 cuenta como 1») queda eliminado; lo que dice sobre `is_active` no cambia.
- **M3 · S2-05.** `TestSumaPartes` tiene 2 subpruebas:
  - `fixtures`: evalúa con `Eval` cada consulta sin error de cada fixture de `oraculo/fixtures/s2/` y recorre la salida completa. Comprueba Σ `by_person` = `total_cents`, `payments_cents + cash_expenses_cents = total_cents`, Σ `shares` = `amount_cents` y, en las vistas, axel + karen = total o monto. Suma con `big.Rat` exacto e informa el conteo.
  - `semilla-12-meses`: la prueba anterior, sin cambios.
- **Fixtures de error** (las escribe `fin_genera_fixtures.py`, único autor según §5.0; las 3 nuevas se generan con el mismo script):
  - `s2-ventanas-sin-cubrir`: PROVISIONAL, cita `S2-ventanas`. Tiene 4 consultas: control P1+P2 contiguas, que reproduce R-45 y R-46 (devengado, efectivo y diferencia fijados desde `model.json`); solo P2; hueco con solo pagos (P1 y P3); hueco con solo un costo (el 10-dic-2026).
  - `s2-rango-devengo-0`: gasto del 5-dic-2026 con `accrual_months: 0`.
  - `s2-rango-intervalo-0`: renta mensual con `interval_count: 0`.
- **Pruebas Go nuevas:**
  - `TestRangoConteos`, con 4 campos. Ausente, null y 1 dan la misma salida; 24 se acepta; 0, −1 y 25 dan `invalid_argument`.
  - 3 casos de ventanas en `TestErrores`.

## Archivos
- Modificados:
  - `backend/internal/motor_fin/input.go`
  - `backend/internal/motor_fin/ledger.go`
  - `backend/internal/motor_fin/projection.go`
  - `backend/internal/motor_fin/plan.go`
  - `backend/internal/motor_fin/engine_test.go`
  - `backend/internal/motor_fin/cover_test.go`
  - `oraculo/fin_devengo.py`
  - `oraculo/fin_genera_fixtures.py`
  - `oraculo/salidas/s2.json`: lo regeneró `fin_run_fixtures.py`.
- Regeneradas: las 27 fixtures anteriores (mismo contenido de salida) y los 3 archivos de `_datos/` de `oraculo/fixtures/s2/`.
- Nuevos:
  - `oraculo/fixtures/s2/s2-ventanas-sin-cubrir.json`
  - `oraculo/fixtures/s2/s2-rango-devengo-0.json`
  - `oraculo/fixtures/s2/s2-rango-intervalo-0.json`
  - `docs/estado/bloqueos/S2-ventanas.md`
  - este reporte
- Sin cambios: `backend/cmd/motores/`, `oraculo/compara.py` y `oraculo/fin_comun.py`.

## Verificación
| Comando | Resultado |
|---|---|
| `cd backend && go vet ./internal/motor_fin/... ./cmd/motores/...` | limpio |
| `gofmt -l internal/motor_fin cmd/motores` | vacío |
| `go test -coverprofile=… ./internal/motor_fin/... -count=1` + `go tool cover -func … \| tail -1` | `ok`, `total: (statements) 89.1%` (antes 88.8 %; K25 ≥ 85 %) |
| `go test ./cmd/motores/... -count=1 -cover` | `ok`, 78.4 % |
| `go test ./internal/motor_fin/... ./cmd/motores/... -count=1 -v` | 309 `--- PASS`, 0 FAIL, 0 SKIP; `grep t.Skip` da 0 |
| `go test … -run 'Suma\|RangoConteos\|Errores' -v` | `S2-05: 30 fixtures, 69 consultas sin error, 303 invariantes con diferencia 0`: Σ by_person 124, pagos + contado 61, Σ shares 93, axel + karen = total 14, axel + karen = monto 11. La semilla da 40 casos. `TestRangoConteos` 4/4 PASS. `TestErrores` PASS |
| `python3 -I oraculo/fin_genera_fixtures.py` | `30 fixtures y 3 archivos de _datos escritos`, 0 líneas `HALLAZGO`: el oráculo da exactamente los valores fijados (errores y control R-45/R-46) |
| `python3 -I oraculo/fin_run_fixtures.py` | `30 fixtures, 30 coinciden, 0 diferencias`. Contra la salida anterior (sha `ba2cbcb1…`): 27 casos idénticos, 0 cambiados, 3 nuevos |
| `go run ./cmd/motores s2 dump --fixtures ../oraculo/fixtures/s2 --out ../datos/tmp/motor-s2.json` | exit 0; `cmp` contra `oraculo/salidas/s2.json`: idénticos |
| `python3 -I oraculo/compara.py s2` (con el archivo del motor) | oráculo contra esperado 2320/2320, motor contra esperado 2320/2320, **oráculo contra motor 2320/2320**, `model.json (ramas de s2): 159 de 159`, cruces 43/43, provisionales 928/928, diferencias 0, exit 0 |
| `python3 -I oraculo/fin_verify.py` | 159/159 (float contra la vía exacta) y 159/159 (float contra `model.json`) |
| `sha256sum oraculo/model.json oraculo/verify3.json` | `1a570666…bfd4` y `d86087f3…5041`, sin cambio |
| `grep -n reporte-compra-latex oraculo/*.py` | 0 líneas |
| `grep -rn 'float32\|float64' internal/motor_fin --include='*.go' \| grep -v _test.go` | 0 líneas |

**Criterios cubiertos con evidencia:**
- **S2-01:** `TestCanon` 85 + semilla, todos PASS.
- **S2-02:** 159/159.
- **S2-03:** 30/30 fixtures en el oráculo y en `TestOraculoFixtures` (30 subpruebas PASS); oráculo = motor en 2320/2320 hojas. Las 2 divergencias de M1 y las 2 de M2 ya no existen.
- **S2-04:** sha sin cambio y `grep` en 0.
- **S2-05:** prueba de propiedad sobre todas las fixtures, 303 invariantes con 0 diferencias. Es el conteo del revisor, 294, más las 9 de la consulta de control nueva.
- **S2-13:** 0 líneas.
- **S2-15:** 89.1 %, vet y gofmt limpios.
- **S2-16:** los dos supuestos que cambiaban cifras ya no lo hacen. El 1 se eliminó. El 2 quedó como el bloqueo `S2-ventanas`, con marca PROVISIONAL en el código y en la fixture.
- Los demás criterios (S2-06 a S2-12 y S2-14) siguen en PASS en la corrida `-v`.

Temporales: `datos/tmp/corr-s2/` y `datos/tmp/motor-s2.json`, borrados.

## Supuestos
- En las dos implementaciones, «hueco» significa una fecha que no está en ninguna ventana y que no es posterior al fin más tardío de las ventanas declaradas. Es parte de la regla PROVISIONAL de `S2-ventanas` y no cambia cifras: solo decide si se responde `invalid_argument`.
- `null` en `accrual_months` o `interval_count` equivale a ausente (vale 1), igual que en el oráculo; la interfaz dice «no (1)».

## Bloqueos
- `docs/estado/bloqueos/S2-ventanas.md` (nuevo): ventanas que empiezan después del primer evento y eventos en huecos. Se decide en la puerta 2.
- Siguen abiertos sin cambio: `S2-flujo-quincenal`, `S2-metas`, `S2-metas-resumen`, `S2-proyeccion-interruptor` y `S2-presupuesto-total`.

## Pendientes
- Ninguno dentro de mi alcance. Los bajos B1 (paquetes `precios` y `recetas` de otros agentes) y B2 (RN-18, depende de `S2-metas`) siguen como los dejó el revisor.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S2-interfaz.md` §2.12 y §3.2: citar el bloqueo `S2-ventanas` y escribir la regla PROVISIONAL (costo o pago antes de la primera ventana, o en un hueco, da `invalid_argument`; lo posterior a la última ventana se admite). Cuando Axel decida, poner la regla definitiva. Motivo: M1 y protocolo §2.
2. `docs/contratos/motores/S2-interfaz.md` §5 (tabla de fixtures, junto a la línea 339): agregar `s2-ventanas-sin-cubrir` (S2-03, S2-16, PROVISIONAL `S2-ventanas`), `s2-rango-devengo-0` (S2-03, §2.4) y `s2-rango-intervalo-0` (S2-03, §2.7). Motivo: son fixtures nuevas del mismo generador.
3. `oraculo/README.md` líneas 22 y 30: «27 fixtures» pasa a «30 fixtures» (`30 fixtures, 30 coinciden, 0 diferencias`). Motivo: el conteo cambió.
4. `docs/estado/bloqueos/INDICE.md`: agregar `S2-ventanas` (puerta 2). Motivo: es un bloqueo nuevo.
5. Integrador: los supuestos 1 y 2 de `docs/estado/wf1/motor-S2-go.md` y el último supuesto de `oraculo-S2.md` quedan sustituidos por este reporte y por `S2-ventanas`.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Para `S2-ventanas`, propongo la **opción A + H1**, que es la que quedó aplicada: rechazar ventanas que no cubren desde el primer evento y huecos con eventos. Justificación:
  - Las ventanas solo existen para las fixtures con eventos entre los días 1 y 9 (§2.12, S2-09).
  - El periodo base del sistema es el mes calendario (decisiones B), y en meses la diferencia acumulada ya es completa.
  - Con A, la cifra no puede salir incompleta en silencio.
- No hay valores de K31 en este alcance.

## Consumo aproximado
- 36 acciones de herramienta.
- Unos 150 KB leídos: protocolo, revisión, §K, criterios S2 completos, los 2 reportes, secciones de la interfaz y fragmentos del código.
- Unos 40 KB escritos con herramientas: scripts de edición, bloqueo y reporte. El generador escribió además unos 13 KB de fixtures nuevas.
