# REVISIÓN S3 · ronda 2 de 3 (revisor-S3-r2, 11-oct-2026)

Evaluado solo contra `docs/criterios/S3.md` 1.2, decisiones §C y §K y `docs/contratos/motores/S3-interfaz.md`. Cada verificación la volví a correr con mis comandos y mis scripts en `datos/tmp/rev-s3r2/` (ya borrados). No corregí nada.

## Veredicto
**NO APROBADO: 0 hallazgos altos, 2 medios y 4 bajos.** Los tres hallazgos de la ronda 1 (H1, H2 y H3) están **cerrados**. El cálculo está bien: núcleo y comparador reproducen el canon, el oráculo y mis propias cuentas al centavo. Quedan dos problemas de integración del oráculo, los dos medios:
- **M1.** La salida guardada del oráculo, `oraculo/salidas/s3.json`, sigue siendo la de antes: 21 fixtures y nada del comparador. Si se corre `compara.py s3` con su ruta por defecto, salen 2,036 diferencias y código 1.
- **M2.** El oráculo tarda unos 4.5 s en `s3-comparador-carga` y aplica el mismo límite de reloj de 5 s. Con la CPU compartida responde `timeout` y `compara.py s3` sale con código 1.

## Comandos y resultados (míos)
- `gofmt -l internal/motor_cred cmd/motores`: vacío. `go vet ./internal/motor_cred/... ./cmd/motores/`: limpio.
- `go test ./internal/motor_cred/... -count=1 -cover`: `motor_cred` ok, 93.7 %; `canontest` ok, 90.0 %; `comparador` ok, 96.0 %; `evaluar` ok, 90.8 %. `go tool cover -func … | tail -1`: **total 94.0 %**.
- `-v`: 133 `--- PASS`, 0 `FAIL` y 0 `SKIP`; `t.Skip` en `motor_cred`: 0. `TestCanon85` da 85 de 85 subpruebas. Hay pruebas con nombre para cada patrón `-run` de los criterios: Canon, Series, Calendario, Reposicion, Tasa, Promocion, Tarjeta, Prestamo, Tienda, Unificado, Restricciones, Optimos, Puntaje y Reparto. `TestOraculoFixtures` corre en `evaluar` («27 fixtures corridas») y en `comparador` («6 fixtures del comparador corridas»).
- Tiempo del motor (`S3-TIEMPO`): caso-1 1 ms, caso-2 0 ms, caso-3 1 ms y carga 132 ms con 3,870 candidatos. `TestRestriccionesPropiedad`: 300 comparaciones, 11,730 candidatos y 0 válidos con violación.
- `python3 -I oraculo/cred_escenarios.py`: **código 0**, en 4.67 s. Casos 1 y 3: 31 evaluados, 17 válidos y 14 descartados. Caso 2: 26 evaluados, 12 válidos y 14 descartados. En los tres, 0 aceptados con violación y 4 de 4 óptimos distintos. `sin-restricciones` da `missing_constraints`, `tiempo` da `timeout` y la carga, 3,870 evaluados en 4,516 ms.
- Compilé `motores` en `datos/tmp` y corrí `motores s3 dump --fixtures oraculo/fixtures/s3`: código 0. `cred_run_fixtures.py --out datos/tmp/…`: 27 fixtures y 79 consultas; segunda vía, 174 de 174.
- `compara.py s3 --oraculo <mío> --motor <mío>`: **906 de 906** en las tres direcciones; **model.json (ramas de s3) 437 de 437**; data.json 9 de 9; provisionales 2,718 de 2,718; diferencias 0; código 0.
- Mi regresión Kueski (`kueski.py`: enteros de centavos, `Fraction`, los modos de K14, `sched` propio) contra `model.json` y contra el dump del motor: **218 de 218**. Comprobé:
  - 800.50; 973.75; 598.82/374.93; 492.28/308.22; capital repuesto 773.63 ×3 + 773.61.
  - 1,034.74 y 5,034.74; 1,258.68 ×3 + 1,258.70; 883.75 ×3 + 883.77.
  - Las 16 fechas, con 2027-02-28.
  - Series A y B de 14 valores, con mínimos 8,094.50 y 5,000.00. Con la lectura omitida sale B y `pending_confirmation` es verdadero.
  - Deuda viva 6,293.44 = 1,258.70 + 5,034.74.
  - Sensibilidad (5 renglones) e interés por monto (7 renglones).

  Mi primera corrida tuvo 3 diferencias, y el error era mío: calculé el pago de referencia de las cotizaciones truncando total/4. El canon usa HALF_UP(monto × (1 + r) / 4) (K14, «pago de referencia solo en cotizaciones»). Lo corregí y quedó 218 de 218.
- Mi comparador (`comp.py`) sobre los casos 1 a 3:
  - Salida completa motor contra oráculo: 709/709, 601/601 y 711/711 hojas, 0 diferencias; también 0 en sin-restricciones, tiempo y carga.
  - 0 válidos con violación de las restricciones de la petición; 0 descartes con la restricción mal declarada.
  - `score_bp` recalculado por mí con mín-máx entre los válidos, empate = 1 y HALF_UP(10000 × Σwn/Σw): 0 diferencias.
  - `value_per_peso_milli` recalculado con la fórmula de §8.4 (4): 0 diferencias.
  - Cada óptimo tiene el mejor valor de su objetivo (4 de 4), los 4 óptimos son distintos y la recomendación tiene el puntaje máximo (8,014, 7,205 y 7,963).
  - `sin-restricciones`: `missing_constraints` con las 4 faltantes en el orden de C1, `evaluated_count` 0, `optima` y `recommendation` nulos y 0 válidos.
- `tabla_mano.py`: «OK … = tabla a mano» en las 5 fixtures de la tabla. No importa `cred_lib` ni el motor.
- Préstamo francés (`loan.py`, aritmética propia con `Fraction`): las 3 tablas son iguales al motor renglón por renglón (12/12, 7/7 y 7/7). En las tres, Σ amortización = capital, saldo final 0 y Σ pagos = capital + Σ intereses. El último pago ajusta: 964.28 → 964.24, 156.42 → 156.41 y 154.51 → 154.53.
- Tarjeta (dump): el pago total un día antes y el mismo día dan `interest_cents` 0, con datos y sin datos. Un día después con la fixture da 6,598, y mi cálculo da lo mismo: Σ saldos diarios × 45.5 %/360 × 1.16, HALF_UP (PROVISIONAL `S3-formula-tarjeta`). Sin datos: `missing_capture`. MSI: 333.33 + 333.33 + 333.34 = 1,000.00. Mínimos: 200.00, 500.00 y 150.00.
- Los `grep`:
  - `highs|ortools` en go.mod y go.sum: 0.
  - `net/http` fuera de pruebas: 0.
  - `float32|float64` fuera de pruebas: 0. En pruebas solo están `ratioFloat` y `modelFloat`, para las hojas `tasa`.
  - `reporte-compra-latex` en `oraculo/*.py`: 0.
  - Go no lee `oraculo/salidas` ni llama Python.
  - `cred_*.py` no leen la salida del motor.
  - No hay literales canon (80050, 97375, 503474, 809450, 629344…) en `*_test.go`.
- `sha256sum oraculo/model.json` antes y después: `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`.
- Migración 00003: `comparator.timeout_seconds` = 5, sin PROVISIONAL (línea 86). `optimizer.min_available_buffer_cents`, `max_live_purchases`, `zero_debt_by_date` y `date_window_days` van como `@param_sin_valor` (líneas 54 y 55).
- **M1:** `compara.py s3 --motor <mío>`, con el oráculo por defecto (`oraculo/salidas/s3.json`): oráculo contra esperado 893 de 906; provisionales 695 de 2,718; **diferencias 2,036; código 1**.
- **M2:** corrí `cred_escenarios.py` y `cred_run_fixtures.py` con `taskset -c 5` mientras otro proceso usaba el mismo núcleo. La carga da `error timeout (5003 ms)` y aun así `cred_escenarios.py` sale con 0. `compara.py s3` con esa salida da provisionales 2,716 de 2,718, **diferencias 2** (`carga.status` y `carga.evaluated_count` del oráculo = `__falta__`) y **código 1**.

## Criterios
| Criterio | Estado | Evidencia (mía) |
|---|---|---|
| S3-01 | Cumple | `TestCanon85` 85/85; mi regresión 218/218 contra `model.json` y el motor |
| S3-02 | Cumple | `compara.py s3` con salidas frescas: 437 de 437; 906/906 en las tres direcciones (con la salida guardada, M1) |
| S3-03 | Cumple | 28/28 iguales a mi cálculo, a R-65 y a R-66; mínimos 8,094.50 y 5,000.00; con la lectura omitida, B y `pending_confirmation` |
| S3-04 | Cumple | 16/16 fechas iguales a mi `sched` y a `model.json`; 2027-02-28 (y 2028-02-29) |
| S3-05 | Cumple | Σ capital repuesto = capital en las 4 compras; `TestReposicionParteSeparada` PASS (K15) |
| S3-06 | Cumple | 800.50/3,094.50 se guarda exacta; 4,000 → 1,034.74; plazos 6, 8, 10 y 12 → `no_rate`; «25.8685 %» y «17.27 %» como se capturaron (K12) |
| S3-07 | Cumple | promo-28 y promo-30 → 0.00; promo-27 y mínimo − 0.01 → no aplica y tasa normal; mínimo exacto aplica |
| S3-08 | Cumple | 0.00 un día antes y el mismo día, con datos y sin datos; 6,598 = mi cálculo = oráculo; sin datos, `missing_capture`; MSI exacto (fórmula PROVISIONAL con bloqueo) |
| S3-09 | Cumple | 3/3 tablas iguales a mi cálculo renglón por renglón, con ajuste de centavos; cuadran |
| S3-10 | Cumple | abonos 333.33/333.33/333.34; con datos, 46.40 y mínimo 100.00 = oráculo; sin datos, `missing_capture` |
| S3-11 | Cumple | `TestUnificado` PASS; `s3-unificado*` 906/906 |
| S3-12 | Cumple, con M2 en la vía del oráculo | 0 válidos con violación (mío y prueba de propiedad); conjunto motor/oráculo/tabla = 0 diferencias en los 3 casos; «faltan restricciones» sin valores; 5 s sembrado; 0 s → 503 `unavailable`; peor caso 132 ms; no hace falta el bloqueo de tiempo |
| S3-13 | Cumple | 12 óptimos = esperado = oráculo; 4 distintos por caso; cada uno con el mejor valor de su objetivo (mío) |
| S3-14 | Cumple | `score_bp` recalculado 46/46; empate = 1; `TestPuntaje` PASS (100/200/400 → 1, 2/3, 0; los pesos cambian la recomendación) |
| S3-15 | Cumple | `grep`: 0 y 0 |
| S3-16 | Cumple | sha256 igual; oráculo y motor independientes (`grep`); `reporte-compra-latex`: 0 |
| S3-17 | Cumple | `float` fuera de pruebas: 0; tasa como razón exacta |
| S3-18 | Cumple | vet y gofmt limpios; 93.7/90.0/96.0/90.8 %, total 94.0 % |
| S3-19 | Cumple | 1,903.00/1,191.50 → 492.28/308.22; Axel + Karen = total en cada pago (mío) |
| S3-20 | Cumple | 5 bloqueos `S3-*`; §Supuestos en los 4 reportes; marcas PROVISIONAL en las salidas |

## Hallazgos
**M1 (media) · `oraculo/salidas/s3.json` no tiene el comparador.** El archivo tiene 21 fixtures (73 consultas) y ninguna `s3-comparador-*`; las fixtures actuales son 27, con 79 consultas. `compara.py s3` lee esa ruta cuando no se le pasa `--oraculo` (línea 239): sin `--oraculo` salen 2,036 diferencias y código 1. La salida guardada no corresponde a las fixtures. correccion-S3-r1 lo dejó como solicitud al integrador y nadie lo ha regenerado.
Corrección: correr `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json` y confirmar que `python3 -I oraculo/compara.py s3 --motor <dump>` da 0 diferencias con su oráculo por defecto.

**M2 (media) · El oráculo mide reloj en la fixture de carga y queda al 90 % del límite.** `oraculo/cred_escenarios.py` (líneas 157 a 182) aplica el límite real de `timeout_seconds` = 5 y en la carga tarda unos 4.5 s (4,516 y 4,549 ms en dos corridas); el motor tarda 132 ms. Con la CPU compartida, la carga responde `timeout` y `compara.py s3` falla (2 diferencias, código 1). La salida oficial depende de la carga del equipo, y aquí corren varios agentes a la vez. El contrato (§8.7) fija el límite para el motor, pero no dice si el oráculo debe medir reloj. El siguiente implementador tendría que adivinarlo.
Corrección: que el integrador del contrato decida en §8.7 y §10.0 una de dos opciones:
- (a) El oráculo trata `timeout_seconds` 0 como `timeout`, que es determinista, y no mide reloj con valores positivos. El tiempo real lo mide solo `TestRestriccionesTiempo` en el motor. Recomiendo esta.
- (b) Acelerar el oráculo hasta dejar un margen amplio, por ejemplo menos de 1 s.

**B1 (baja) · `cred_escenarios.py` no compara la carga con su `esperado`.** Con la carga en `timeout` el script sale con 0, así que su código de salida no detecta que el resultado de la carga no es `ok` con 3,870 evaluados.
Corrección: que compare el `status` y el `evaluated_count` de cada fixture con su `esperado` y salga con 1 si difieren.

**B2 (baja; H4 de r1, sigue abierto) · `value_per_peso_milli` redondeado a milésimas crea empates en el objetivo 4.** En los 3 casos no cambia el óptimo.
Corrección: la misma de r1 (escala ×10⁶ o la fracción exacta, vía el contrato y el bloqueo `S3-objetivos-medida`).

**B3 (baja; H5 de r1, sigue abierto) · El caso 3 da los mismos 4 óptimos que el caso 1.** Las llaves son idénticas; solo cambian los puntajes.
Corrección: igual que en r1, en una revisión futura de los casos.

**B4 (baja; H6 y H7 de r1, siguen abiertos) · Al contrato le faltan dos cosas.**
- La §2 no menciona `canontest/` (0 coincidencias).
- La §8.2 no dice que `term_loan` no aporta candidatos.

Corrección: aplicar las solicitudes del núcleo y del comparador a §2, §8.2 y §11.

## Estado de los anteriores
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H1 (alta): faltan las 9 fixtures del comparador | **Cerrado** | Existen `_datos/comparador-caso-{1,2,3}.json` y las 6 `s3-comparador-*`; `TestOraculoFixtures` PASS; `cred_escenarios.py` sale con 0; `tabla_mano.py` da OK |
| H2 (alta): «comparar» devuelve `ErrUnsupported` | **Cerrado** | `evaluar.go:197` llama `comparador.Query`; el dump trae los 6 casos; `compara.py s3` da 906/906 y la salida completa motor = oráculo |
| H3 (media): el comparador no tiene `TestOraculoFixtures` y hay copia en `testdata/` | **Cerrado** | `comparador/TestOraculoFixtures` PASS (6 fixtures); en `testdata/` solo queda `tabla_mano.py`; las pruebas leen `oraculo/fixtures/s3/` |
| H4, H5, H6 y H7 (bajas) | Abiertos | Ahora son B2, B3 y B4 |
