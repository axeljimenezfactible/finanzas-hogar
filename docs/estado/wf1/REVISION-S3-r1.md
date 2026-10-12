# REVISIÓN S3 · ronda 1 de 3 (revisor-S3-r1, 11-oct-2026)

Revisor independiente. No escribí nada de lo revisado y no corregí nada. Evalué solo contra `docs/criterios/S3.md` 1.2, decisiones §C y §K y `docs/contratos/motores/S3-interfaz.md` 1.4 (§8, §10.0, §11). Todas las verificaciones las corrí yo, con scripts propios en `datos/tmp/rev-s3/` (borrados al terminar).

## Veredicto
**NO APROBADO.** 2 hallazgos altos, 1 medio y 4 bajos. El cálculo está bien: núcleo y comparador reproducen el canon, el oráculo y mis propias cuentas al centavo. Lo que falla es la integración del comparador con el oráculo: faltan las 9 fixtures del comparador en `oraculo/fixtures/s3/` y la consulta «comparar» no está conectada. Por eso `go test ./internal/motor_cred/...` FALLA (`TestOraculoFixtures`) y `python3 -I oraculo/cred_escenarios.py` sale con código 1, y esas son las verificaciones literales de S3-12, S3-13 y S3-18.

## Comandos y resultados (míos)
- `gofmt -l internal/motor_cred cmd/motores`: vacío. `go vet ./internal/motor_cred/... ./cmd/motores/`: limpio.
- `go test ./internal/motor_cred/... -count=1 -cover`: `motor_cred` ok 93.7 %, `canontest` ok 90.0 %, `comparador` ok 96.0 %, **`evaluar` FAIL** 90.5 % (`TestOraculoFixtures`: «falta oraculo/fixtures/s3/_datos/comparador-caso-1.json (autor: agente de optimizador S3)» ×9; «21 fixtures corridas»). Total 94.0 %. `t.Skip`: 0; `--- SKIP`: 0.
- `-run TestCanon85 -v`: 85 subpruebas PASS. Pruebas `Series`, `Calendario`, `Reposicion*`, `Tasa*`, `PromocionFronteras`, `Tarjeta*`, `PrestamoFrances`, `Tienda*`, `Unificado`, `Reparto`, `CanonHojas437`: todas PASS.
- `motores s3 dump` (compilado en datos/tmp) → 21 casos, código 0; `cred_run_fixtures.py --out datos/tmp/...` → 21 fixtures, 73 consultas, segunda vía 174 de 174; `compara.py s3 --oraculo … --motor …`: oráculo/esperado, motor/esperado y oráculo/motor **893 de 893**; **model.json (ramas de s3) 437 de 437**; data.json 9 de 9; provisionales 695 de 695; diferencias 0. sha256 de `model.json` igual antes y después (`1a5706…bfd4`).
- Regresión Kueski con mi propio cálculo (enteros de centavos, `Fraction`, HALF_UP y base truncada de K14) contra el dump del motor y `model.json`: 17 de 17. 800.50; 973.75; 598.82/374.93; 492.28/308.22; capital repuesto 773.63 ×3 + 773.61; 1,034.74 y 5,034.74; 1,258.68 ×3 + 1,258.70; 883.75 ×3 + 883.77; series A y B de 14 valores (mínimos 8,094.50 y 5,000.00); lectura omitida = B con `pending_confirmation` verdadero; deuda viva 6,293.44 = 1,258.70 + 5,034.74; 16 fechas iguales a `model.json`, con 2027-02-28 (y 2028-02-29 en el motor).
- Préstamo francés recalculado por mí (pago = HALF_UP(P·r/(1−(1+r)^−n)), interés HALF_UP(saldo·r), el último pago liquida el saldo): 3 de 3 tablas iguales al motor fila por fila. Σ amortización = capital, saldo final 0.00, Σ pagos = capital + Σ intereses; hay ajuste en la última fila (156.42 → 156.41 y 964.28 → 964.24).
- Tarjeta (dump): pago total un día antes y el mismo día, con datos y sin datos → `interest_cents` 0; un día después con fixture → 6,598 (base diaria/360, IVA 16 %, PROVISIONAL `S3-formula-tarjeta`) = oráculo; sin datos → `missing_capture`; MSI 333.33 + 333.33 + 333.34 = 1,000.00.
- Comparador. Corrí el oráculo (`cred_lib.evaluate` / `cred_escenarios._resumen`) sobre copias de `comparador/testdata/*.json`: caso-1 633 hojas, caso-2 535, caso-3 635, sin-restricciones 8 y tiempo 1, todas con 0 diferencias contra el `esperado` (tabla a mano). Las pruebas Go comparan el motor contra ese mismo `esperado` con listas completas (`subset`), así que motor = tabla = oráculo en los 3 casos. Con mi script de puntaje: 0 válidos con violación; 14 de 14 descartes por caso violan la restricción que declaran; `score_bp` mín-máx con pesos 1,1,1,1 recalculado 17/17, 12/12 y 17/17; cada óptimo tiene el mejor valor de su objetivo; 4 óptimos distintos en cada caso; la recomendación es el puntaje máximo. Sin restricciones → `missing_constraints` con 4 faltantes, 0 evaluados y `optima`/`recommendation` nulos. `TestRestriccionesPropiedad`: 300 comparaciones, 11,730 candidatos, 0 violaciones. Tiempos (`S3-TIEMPO`): caso-1 1 ms, caso-2 0 ms, caso-3 1 ms, carga 114 ms (3,870 candidatos); 0 s → `ErrTimeout` → 503 `unavailable`.
- `grep highs|ortools` en go.mod y go.sum: 0. `net/http` fuera de pruebas: 0. `float32/float64` fuera de pruebas: 0 (en pruebas solo `ratioFloat` y `modelFloat` de `canon_test.go`, para las hojas flotantes `tasa` y `tasa_pct`). `reporte-compra-latex` en `oraculo/*.py`: 0. `cred_*.py` no leen la salida del motor; `tabla_mano.py` y `cred_genera_fixtures.py` no importan `cred_lib` ni el motor; `canontest` solo lo importan pruebas.
- 00003: `comparator.timeout_seconds` = 5 sembrado y sin marca PROVISIONAL. `optimizer.min_available_buffer_cents`, `max_live_purchases`, `zero_debt_by_date` y `date_window_days` van sin valor (`@param_sin_valor`). Ninguna migración precarga promociones ni tasas por plazo.

## Criterios
| Criterio | Estado | Evidencia (mía) |
|---|---|---|
| S3-01 | Cumple | `TestCanon85` 85/85; regresión propia 17/17 |
| S3-02 | Cumple | `compara.py s3`: 437 de 437; 893/893 en las tres direcciones |
| S3-03 | Cumple | series A/B 28/28 iguales a mi cálculo y a R-65/R-66; omitida → B y `pending_confirmation` |
| S3-04 | Cumple | 16/16 fechas; 2027-02-28 |
| S3-05 | Cumple | Σ capital repuesto = 3,094.50; `s3-reposicion-separada`: el pago de Karen (29-oct) repone 0 y el que completa repone 773.63 (K15) |
| S3-06 | Cumple | 6/8/10/12 q → `no_rate`; 4 q → 1,034.74; tasa 80050/309450 con decimales 4 → «25.8685 %» (K12) |
| S3-07 | Cumple | día 28 y día 30 → 0.00; día 27 y 1,999.99 → `promotion_not_applicable` y tasa normal; mínimo exacto aplica |
| S3-08 | Cumple | 0.00 un día antes y el mismo día (con y sin datos); fórmula PROVISIONAL con bloqueo; `missing_capture`; MSI exacto |
| S3-09 | Cumple | 3 tablas francesas propias = motor; ajuste de centavos |
| S3-10 | Cumple | abonos 333.33/333.33/333.34, saldo 0; con interés 46.40 = oráculo (PROVISIONAL); sin datos `missing_capture`; compra con interés `unsupported` según §5.5 |
| S3-11 | Cumple | `TestUnificado` PASS; `s3-unificado*` 0 diferencias en compara |
| S3-12 | **No cumple** | cálculo correcto (mi corrida: 0 violaciones, oráculo = tabla = motor), pero `python3 -I oraculo/cred_escenarios.py` → código 1 por fixtures faltantes y no hay comparación motor/oráculo por la vía oficial (H1, H2) |
| S3-13 | **No cumple (parcial)** | `-run Optimos` PASS con 4 distintos por caso; «coinciden con el oráculo» solo por mi corrida, no por `TestOraculoFixtures` ni `compara.py` (H1, H2) |
| S3-14 | Cumple | `TestPuntaje` PASS; 100/200/400 → 1, 2/3, 0; mis 46 puntajes = motor; empate declarado en §8.5 |
| S3-15 | Cumple | 0 / 0 |
| S3-16 | Cumple | sha256 igual; sin escrituras fuera |
| S3-17 | Cumple | 0 floats monetarios |
| S3-18 | **No cumple** | vet y gofmt limpios, cobertura ≥ 85 % en los 4 paquetes, pero `go test ./internal/motor_cred/...` termina en FAIL (H1, H2) |
| S3-19 | Cumple | 492.28/308.22; Axel + Karen = total en cada pago (d1 y d2) |
| S3-20 | Cumple | §Supuestos en los 3 reportes; bloqueos S3-convencion-prestamo, -formula-tarjeta, -objetivos-medida, -origen-cargo-normal y -serie-sin-intereses citados en las listas `provisional` |
| K12–K18, K43, K44 | Aplicados | K12 decimales; K13 topes solo de `caps`; K14 dos modos; K15; K16 sin valores por defecto; K17 medidas y mín-máx; K18 `missing_capture`; K43 5 s sembrado; K44 clave configurable que lee `TimeoutFromSetting` |

## Hallazgos
**H1 (alta) · Faltan las 9 fixtures del comparador en `oraculo/fixtures/s3/`.** No existen `_datos/comparador-caso-{1,2,3}.json` ni `s3-comparador-{caso-1,caso-2,caso-3,sin-restricciones,tiempo,carga}.json`, que §10.0 asigna al optimizador S3. Están solo como copias en `backend/internal/motor_cred/comparador/testdata/`, y `-carga` ni siquiera existe como fixture. Consecuencias: `TestOraculoFixtures` FALLA, `go test ./internal/motor_cred/...` da FAIL (S3-18) y `cred_escenarios.py` sale con 1 (S3-12).
Corrección: el integrador copia sin cambiar un número `testdata/caso-N.json` → `s3-comparador-caso-N.json`; separa `entradas.datos` en `_datos/comparador-caso-N.json` con `$ref`, con `-sin-restricciones` y `-tiempo` como variantes del caso 1; escribe `s3-comparador-carga` con la rejilla de `carga()` (`evaluated_count` 3,870) y su `esperado`; vuelve a correr `cred_escenarios.py`, `cred_run_fixtures.py`, `s3 dump` y `compara.py s3`.

**H2 (alta) · «comparar» no está conectada.** `backend/internal/motor_cred/evaluar/evaluar.go:197` devuelve `ErrUnsupported`. Ni `motores s3 dump` ni `compara.py s3` incluyen el comparador, así que no hay comparación motor/oráculo por la vía oficial y `TestOraculoFixtures` no pasará aunque existan las fixtures.
Corrección: en el despacho de `"comparar"` llamar `comparador.EvalComparar(in, raw)` (sin ciclo de importación: `comparador` no importa `evaluar`) y volver a correr la verificación §10.2.

**H3 (media) · El paquete `comparador` no tiene su `TestOraculoFixtures`.** §11 dice que ese test va en «ambos». Las pruebas del comparador leen `comparador/testdata/`, una segunda copia de la tabla a mano que puede divergir de la de `oraculo/fixtures/s3/`.
Corrección: después de H1, que las pruebas del comparador lean `oraculo/fixtures/s3/s3-comparador-*` (o agregar `TestOraculoFixtures` en `comparador`) y dejar una sola fuente de la tabla.

**H4 (baja) · El redondeo de `value_per_peso_milli` crea empates artificiales en el objetivo 4.** Se redondea a milésimas enteras (§8.4 (4)). En caso-1 y caso-3, 26.087 y 25.581 dan 26, y 17.391 y 17.054 dan 17; en caso-2, 24.000 y 23.834 dan 24. En los 3 casos no cambia el óptimo, porque el grupo del máximo tiene un solo escenario. Con otras entradas el desempate (puntaje o llave) podría elegir el peor valor exacto.
Corrección: solicitud al contrato y al bloqueo `S3-objetivos-medida` (1) para escalar más (×10⁶) o normalizar con la fracción exacta.

**H5 (baja) · El caso 3 da los mismos 4 óptimos que el caso 1.** Mismas llaves y mismos conteos 31/17/14; solo cambian los puntajes, por ejemplo la recomendación 7,963 contra 8,014. Cumple «4 distintos en cada caso», pero aporta poca variedad.
Corrección: en una revisión futura de los casos (sin tocar los ya fijados), que el caso del tope u `override` mueva al menos un óptimo.

**H6 (baja) · `term_loan` no aporta candidatos y §8.2 no lo dice.** El supuesto está documentado y no cambia cifras.
Corrección: aplicar la solicitud 3 del comparador a §8.2.

**H7 (baja) · `canontest/` no figura en §2 del contrato, y §11 no dice que el `TestOraculoFixtures` del núcleo vive en `evaluar`.**
Corrección: aplicar la solicitud del núcleo a §2 y §11.

## Estado de los anteriores
Ronda 1: no hay hallazgos previos de esta revisión de motores.
