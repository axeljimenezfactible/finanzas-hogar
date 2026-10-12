# REVISION-S2-r1 · Motor financiero S2 (Go `motor_fin` + oráculo `fin_*.py`)

Revisor independiente, ronda 1 de 3 · 11-oct-2026 · Base: `docs/criterios/S2.md` 1.2, `docs/plan/decisiones.md` §B y §K, `docs/contratos/motores/S2-interfaz.md` 1.2. Ninguna cifra la tomé de los reportes: todo se reejecutó con comandos propios y con scripts en `datos/tmp/rev-s2/` (ya borrados).

## Veredicto
**NO APROBADO**: 0 hallazgos altos, **3 medios** y 2 bajos.
Las cifras están bien: canon 85/85, 159/159 hojas contra `model.json` (comparación propia), 27/27 fixtures, oráculo = motor en 2318/2318 hojas (las dos salidas son idénticas byte a byte), y una fixture propia calculada a mano (248 hojas) con 0 diferencias. Lo que falta: dos reglas del motor que cambian cifras en entradas válidas fuera de las fixtures, en las que el motor y el oráculo no coinciden, y una prueba de S2-05 que no cubre lo que pide la Verifica.

## Verificaciones reejecutadas (resumen)
- (a) `cd backend && go vet ./internal/motor_fin/... ./cmd/motores/...`: limpio. `gofmt -l internal/motor_fin cmd/motores`: vacío. `go test ./internal/motor_fin/... -count=1 -cover`: `ok … coverage: 88.8% of statements`; `go tool cover -func … | tail -1`: `total: (statements) 88.8%`; `cmd/motores` 78.4 %. `-v`: 298 `--- PASS`, 0 FAIL, 0 SKIP; `grep t.Skip`: 0. A nivel módulo, `gofmt -l .` y `go vet ./...` fallan **solo** en paquetes de otros agentes (B1).
- (b) `python3 -I oraculo/fin_run_fixtures.py`: `27 fixtures, 27 coinciden, 0 diferencias` (sha256 de `salidas/s2.json` sin cambio: `ba2cbcb1…dce7`). `go run ./cmd/motores s2 dump --fixtures ../oraculo/fixtures/s2 --out ../datos/tmp/motor-s2.json`: exit 0. `cmp` de las dos salidas: idénticas. `python3 -I oraculo/compara.py s2`: oráculo contra esperado 2318/2318, motor contra esperado 2318/2318, **oráculo contra motor 2318/2318**, `model.json (ramas de s2): 159 de 159`, cruces 43/43, provisionales 867/867, diferencias 0, exit 0.
- Script propio `hojas.py`: recorre `model.json`. Cuenta 122 + 12 + 13 + 12 = **159** hojas S2 de 596, convierte los float a centavos con `Decimal(str(v))·100` y compara contra la salida del motor y la del oráculo. Resultado: **motor 159/159, oráculo 159/159**.
- (c) Independencia: `grep backend|motor-|datos/tmp|subprocess` en `oraculo/fin_*.py`: solo un comentario y la ruta por omisión del archivo del motor en `compara.py:240`, que es la herramienta de cruce y no el oráculo. `motor_fin` no hace E/S en el código que no es de prueba; sus importaciones son solo de la biblioteca estándar. Las pruebas leen `oraculo/model.json` y las fixtures, pero ningún `.py` ni `salidas/`. `trim` de `cmd/motores/main.go:218` toma de `esperado` solo la forma, nunca los valores. El método es distinto: el oráculo usa `Fraction` y recorre día por día (`fin_devengo.py`); el motor agrupa con enteros.
- (d) `sha256sum oraculo/model.json oraculo/verify3.json` antes, durante y después: `1a570666…bfd4` y `d86087f3…5041`, sin cambio. `grep -n reporte-compra-latex oraculo/*.py`: 0 líneas. `python3 -I oraculo/model.py` y `verify3.py` escriben solo en `datos/tmp/`, y `cmp` contra los canónicos da iguales. `~/reporte-compra-latex/model/verify3.json` no se toca desde el 10-oct 18:05, y `find ~/reporte-compra-latex -newer …` no devuelve nada.
- (f) Fixture propia `rev-propio`, con datos inventados por mí y esperado calculado a mano: mes de 31 días, gasto el 31-dic, compra a crédito de 120001 en 3 pagos, débito a 2 meses, ingresos, 4 renglones de presupuesto con renglón del total. Motor y oráculo la reproducen en sus **248 hojas, con 0 diferencias**. Fixture propia `rev-cadencia-apagada`: la cadencia de `s2-proyeccion-7-12-cadencia` con el interruptor **apagado**. Da exactamente la proyección por defecto (326 hojas) y `provisional` null en abr–sep 2027.

## Criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S2-01 (alta) | CUMPLE | `TestCanon`: 85 subpruebas R-01…R-85 + `semilla`, todas PASS, 0 SKIP. Los valores se leen de `model.json` en tiempo de prueba (`helpers_test.go:55`). |
| S2-02 (alta) | CUMPLE | `compara.py s2` da 159/159. Mi `hojas.py` da motor 159/159 y oráculo 159/159 contra `model.json`. |
| S2-03 (alta) | CUMPLE (ver M1, M2) | 27/27 fixtures (oráculo y `TestOraculoFixtures`). Oráculo = motor en 2318/2318 hojas. Independencia comprobada. Fuera de las fixtures, el motor y el oráculo divergen en 3 entradas válidas (M1, M2). |
| S2-04 (alta) | CUMPLE | Los sha256 no cambian. `grep` da 0. Las escrituras van solo a `datos/tmp/`, y nada se escribió fuera del repositorio. |
| S2-05 (alta) | PARCIAL (M3) | La propiedad se cumple: mi script sobre la salida del motor en las 27 fixtures revisa **294 invariantes, 0 con diferencia** (Σ `by_person` = total, Σ `shares` = monto, pagos + contado = total, axel + karen = total o monto). Pero `TestSumaPartes` solo recorre la semilla (40 casos), no «todas las fixtures». |
| S2-06 (alta) | CUMPLE | `TestQuincenaFinDeMes`: 4/4. Propio: `calendario {date: 2026-12-31}` da q2 del 16 al 31, paga el 2026-12-30 y cae en la semana **2026-W53** (28-dic a 3-ene). `flujo_quincenal`: q1 `variable_spent` 38333 y q2 17001, que incluye los 10000 del 31-dic. `periodos` de dic-2026: devengado 98668 (axel 59612, karen 39056), efectivo 95334. |
| S2-07 (alta) | CUMPLE | `TestInvarianteDeudaViva`: 7/7 con las cifras del criterio. Propio: corte 31-dic, 40001 = 23335 + 16666 (sin saldo no vale). Corte a media quincena 20-dic, 40001 = −65333 + 105334. |
| S2-08 (alta) | CUMPLE | Con `Decimal`, el comando del criterio imprime `219.88 629.35`. `TestPromedios` usa días reales y R-56/R-57. Propio: dic 95334/31 da 3075 y **21527** (con el diario redondeado ×7 serían 21525). Feb 100001/28 da 3571/25000. Ene, con el mismo total, da 3226/22581: denominadores distintos. 2 centavos en 4 días dan 1/4 (HALF_UP). |
| S2-09 (alta) | CUMPLE | Las pruebas `Proyeccion*` pasan. La semilla no tiene eventos entre los días 1 y 9. Interruptor apagado sin cadencia, y apagado **con** cadencia presente (propio), dan lo mismo. Encendido marca `S2-proyeccion-interruptor` solo en los meses 7–12. |
| S2-10 (alta) | CUMPLE | `TestFlujoQuincenal`: 16 quincenas 59882/37493 y 88375 ×3 + 88377. El contraste da +1 en 9 quincenas y −1 en 3. Propio: 2 quincenas con ingreso fijo y extra, Karen en `external_informative_cents`, saldo acumulado 1444666. |
| S2-11 (alta) | CUMPLE | Salida `s2-alertas-80-100`: 79 % normal; 79.999 % `used_pct` 80 y normal; 80 % warning; 99 % warning; 100 % exceeded; 101 % exceeded. Propio dic-2026: 62 %, 21 % y `no_plan`; total 26 % con `declared_total_row_cents` 500000. |
| S2-12 (media) | CUMPLE (ver B2) | `TestMetasAhorro`. Casos: exacto 100000 ×6; con residuo 44953 ×6 + 44954 = 314672; cumplida 110 %. Va con la marca PROVISIONAL `S2-metas`. |
| S2-13 (alta) | CUMPLE | `grep -rn 'float32\|float64' … --include='*.go' \| grep -v _test.go`: 0 líneas. |
| S2-14 (media) | CUMPLE | `TestZonaHoy`. Propio: `now 2027-01-01T05:59:59Z` da `today 2026-12-31`. |
| S2-15 (alta) | CUMPLE | `vet` y `gofmt` limpios en `motor_fin`; cobertura **88.8 %** (≥ 85 %, K25). |
| S2-16 (media) | PARCIAL (M1, M2) | Los dos reportes tienen SUPUESTOS y los 5 bloqueos `S2-*` existen. Pero los supuestos 1 y 2 del motor cambian cifras con entradas válidas y no se convirtieron en bloqueo (protocolo §2). |

Decisiones §K: **K2** (scope): 5 casos propios. `onward` cierra el renglón el 2026-10-31, abre el id 6 (máx + 1) sin fin y **borra el renglón posterior 2** (`deleted_row_ids [2]`); sin `scope` se comporta como `onward`. `this_month` deja el monto el 2026-11 y crea el id 7 hasta el fin original, sin tocar el 2. Si el renglón empieza en el mes, `onward` lo deja sin fin y borra el posterior, y `this_month` parte la vigencia; el total va primero y luego el orden es por `code`. **K8**, **K19**, **K20**, **K21** y **K39**: comprobadas arriba con los casos propios. Las dos implementaciones coinciden.

## Hallazgos
### M1 (media) · Ventanas que empiezan después del primer evento: el motor da cifras que contradicen §3.2 y el oráculo las rechaza
- Rutas: `backend/internal/motor_fin/ledger.go:257` (`windowed`) y `:439` (`Periods`); `oraculo/fin_devengo.py:78-82`; `docs/contratos/motores/S2-interfaz.md` §2.12 y §3.2.
- Caso reproducido: semilla y `periodos {windows: [{label: P2, start: 2026-11-10, end: 2026-12-09}]}`.
  - **Motor**: `accrued.total_cents 0`, `cash.total_cents 194750` y `cumulative_diff_cents −292125`.
  - **Oráculo**: `invalid_argument` («hay un evento antes de la primera ventana»).
  - §3.2 dice que `cumulative_diff_cents` va «desde el primer evento de las entradas hasta `end` (incluye periodos anteriores no pedidos)». Para esa misma ventana el canon es devengado 194750 y diferencia 97375 (R-46).
  - El motor cuenta el efectivo anterior pero no el devengo anterior, que es su supuesto 2. Tampoco devenga un costo que cae en un hueco entre ventanas (el oráculo hace lo mismo, sin regla escrita).
- Por qué importa: el motor da en silencio una cifra que no coincide con el oráculo ni con el canon. Además, su supuesto cambia cifras y debía ser un bloqueo. El oráculo ya pidió la aclaración (su solicitud 3).
- Corrección sugerida:
  - Abrir el bloqueo `docs/estado/bloqueos/S2-ventanas.md` (ventana posterior al primer evento y hueco entre ventanas) y pedir la aclaración en §2.12 y §3.2 de la interfaz.
  - Mientras se decide, alinear el motor con la regla conservadora del oráculo: `invalid_argument` si hay un costo o pago anterior a la primera ventana, y también si hay uno en un hueco entre ventanas, en las dos implementaciones.
  - Agregar fixtures de error para los dos casos.

### M2 (media) · `accrual_months` e `interval_count` en 0: el motor los acepta como 1; el oráculo y el contrato los rechazan
- Rutas: `backend/internal/motor_fin/input.go:383-386`, `:433-436` y `:498-499` (el `int` usa 0 como «ausente»); `oraculo/fin_comun.py:320`, `:344`, `:413` y `:415` (rango 1–24); `S2-interfaz.md` §2.4, §2.5 y §2.7 (entero 1–24).
- Casos reproducidos:
  - Un gasto con `accrual_months: 0`: el motor devuelve `periodos` dic-2026 devengado 98668; el oráculo da `invalid_argument`.
  - Una partida mensual con `interval_count: 0`: el motor devuelve 800000 cada mes; el oráculo da `invalid_argument`.
- Por qué importa: los dos verificadores divergen en una entrada fuera de rango. El supuesto 1 del motor cambia cifras sin bloqueo.
- Corrección sugerida: distinguir «ausente» de 0 (`*int` o detección de la clave) y devolver `invalid_argument` con 0 o con más de 24. Agregar una fixture de error.

### M3 (media) · S2-05: `TestSumaPartes` no recorre las fixtures
- Ruta: `backend/internal/motor_fin/engine_test.go:26-44`.
- La Verifica de S2-05 pide una «prueba de propiedad sobre todas las fixtures». La prueba solo recorre 12 meses de la semilla (40 casos).
- La propiedad sí se cumple: mi script sobre la salida del motor de las 27 fixtures da 294 invariantes y 0 diferencias. Lo que falta es la prueba.
- Corrección sugerida: que `TestSumaPartes` evalúe con `Eval` cada fixture de `oraculo/fixtures/s2/` y compruebe en toda salida:
  - Σ `by_person` = `total_cents`;
  - Σ `shares` = `amount_cents`;
  - `payments_cents + cash_expenses_cents = total_cents`;
  - en las vistas, axel + karen = total o monto.
  
  Que informe el conteo, como hace el oráculo con sus 267.

### B1 (baja) · `gofmt -l .` y `go vet ./...` del módulo no salen limpios por paquetes ajenos a S2
- `gofmt` lista `internal/precios/seed.go`, `internal/precios/types.go`, `internal/recetas/num.go` y `internal/recetas/types.go`. `vet` falla en `internal/recetas/types.go:229: undefined: Nutrition`.
- `motor_fin` y `cmd/motores` están limpios.
- Corrección: la deben hacer los dueños de esos paquetes. El integrador repite `gofmt -l . && go vet ./...` al juntar WF-1.

### B2 (baja) · S2-12: RN-18 no trae la fórmula del plan de aportaciones
- Ruta: `docs/contratos/reglas-negocio.md` RN-18.
- RN-18 trae el porcentaje con su ejemplo (50.0 %), pero no el reparto del restante en los meses (HALF_UP, el último ajusta) que S2-12 pide en `reglas-negocio.md`. Depende del bloqueo abierto `S2-metas`, y el motor y el oráculo lo llevan PROVISIONAL.
- Además, los «3 fixtures» del criterio son 3 casos dentro de `s2-metas`.
- Corrección: al cerrar `S2-metas`, agregar la fórmula y el ejemplo 44953 × 6 + 44954 = 314672 a RN-18.

## Notas sin hallazgo
- (g) No hay cifras ajustadas. Los dos reportes declaran 0 hallazgos de cifra y mi comparación directa contra `model.json` lo confirma. Para los ids libres, el esperado lo genera el oráculo (§5.0); la prueba independiente es motor contra oráculo, que coincide en 2318/2318 hojas.
- El reporte `motor-S2-go.md` dice que «`compara.py` todavía no existe» y deja pendiente correrlo. Ya existe y lo corrí con el motor: 2318/2318 y 159/159. Ese pendiente queda cerrado.

## Estado de los anteriores
Es la ronda 1 de la revisión de la implementación S2, así que no hay hallazgos previos que seguir.
