# REVISIÓN S3 · ronda 3 de 3 (revisor-S3-r3, 11-oct-2026)

Revisor independiente. No escribí nada de lo revisado y no corregí nada. Evalué solo contra `docs/criterios/S3.md` 1.2, decisiones §C y §K y `docs/contratos/motores/S3-interfaz.md`. Volví a correr cada verificación con mis comandos y con tres scripts propios en `datos/tmp/rev-s3r3/`, ya borrados: `kueski.py` (regresión desde las entradas del plan, con `Fraction`), `comp.py` (restricciones, normalización, puntaje, óptimos y recomendación) y `loan.py` (tabla francesa).

## Veredicto
**APROBADO.** No hay hallazgos altos ni medios. Los dos medios de r2 (M1 y M2) están cerrados y comprobados por mí. Quedan 6 hallazgos bajos: 4 que siguen abiertos desde antes (B1 a B4) y 2 nuevos (B5 y B6).

## Comandos y resultados (míos)
- `sha256sum oraculo/model.json`: `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` al empezar y al terminar.
- `gofmt -l internal/motor_cred cmd/motores`: vacío. `go vet ./internal/motor_cred/... ./cmd/motores/`: limpio.
- `go test ./internal/motor_cred/... -count=1 -cover`: PASS. Cobertura: `motor_cred` 93.7 %, `canontest` 90.0 %, `comparador` 96.0 %, `evaluar` 90.8 %; `go tool cover -func | tail -1` da **total 94.0 %**.
- `go test ./internal/motor_cred/... -count=1 -v`: 133 `--- PASS`, 0 `FAIL`, 0 `SKIP`; `t.Skip` en el código: 0. `TestCanon85`: 85 de 85 subpruebas. Los filtros del criterio (Canon, Series, Calendario, Reposicion, Tasa, Promocion, Tarjeta, Prestamo, Tienda, Unificado, Restricciones, Optimos, Puntaje, Reparto) tienen al menos una prueba cada uno y todas pasan. Líneas de registro: `S3-03: 28 de 28`, `S3-04: 16 de 16`, `S3-12 propiedad: 300 comparaciones, 11730 candidatos, 0 válidos con violación`, `S3-TIEMPO carga 125 ms (3870 candidatos)`, casos de 0 a 1 ms.
- `motores s3 dump --fixtures oraculo/fixtures/s3` (binario recién compilado): código 0.
- `cred_run_fixtures.py` a `datos/tmp`: 27 fixtures, 79 consultas, segunda vía 174 de 174. El sha256 `b84c0db7…bd8f` es **idéntico** al de `oraculo/salidas/s3.json`.
- `compara.py s3 --motor <dump>`, con el oráculo por defecto: 906 de 906 en las tres direcciones; **model.json (ramas de s3) 437 de 437**; data.json 9 de 9; provisionales 2,718 de 2,718; **diferencias 0; código 0**. Sin argumentos también da diferencias 0 y código 0.
- **Prueba con el equipo cargado (M2):** dos procesos ocupando el núcleo 3 y el oráculo en ese mismo núcleo (`taskset -c 3`). Con `cred_escenarios.py`, la carga tardó **7,828 ms** y aun así dio `ok`, 3,870 evaluados y 0 aceptados con violación, con código 0. La salida de `cred_run_fixtures.py` en esas condiciones tiene el **mismo sha256** `b84c0db7…` y `compara.py s3` con ella da diferencias 0 y código 0.
- **Mi regresión Kueski (`kueski.py`) contra el dump y `model.json`: 117 de 117 comprobaciones.** Incluye 800.50 · 973.75 · 598.82/374.93 · 492.28/308.22 · 2,395.28/1,499.72 · 773.63×3 + 773.61 · 5,034.74 · 1,258.68/1,258.70 · 883.75/883.77, las 16 fechas (2027-02-28 en la despensa 3), las series A y B de 14 valores (mínimos 8,094.50 y 5,000.00), la lectura omitida (B con `pending_confirmation`), la deuda viva 6,293.44 al 10-abr-2027, la sensibilidad (5 × 4) y el interés por monto (7 × 3). Comprobé los dos modos de K14 por separado: los montos con base truncada (1,258.685 da 1,258.68 y el último pago 1,258.70) y el capital con base HALF_UP (773.625 da 773.63).
- **Mi comparador (`comp.py`): 0 errores en los 3 casos.** En los válidos, 0 violan las restricciones de la fixture. El caso 3 usa `override`: colchón 300,000 y deuda en cero 2027-02-15, encima del valor global. En cada descartado, `violated_constraint` es la primera violada. Válidos + descartados = `evaluated_count` (31, 26 y 31). El conjunto `valid` y `rejected` del motor es igual al del oráculo. Recalculé `score_bp` de cada válido con mín-máx entre los válidos y con empate = 1, y coincide. Los 4 óptimos coinciden con mi búsqueda exhaustiva, incluidos los desempates de §8.6. La recomendación también coincide. En cada caso hay **4 óptimos distintos**.
- Fórmula de `std_dev_milli_cents` (§8.4): la probé con 20,000 series al azar y es HALF_UP(1000 × desviación estándar) exacta en todas. La autoprueba del oráculo da 0/100/200 → 81,650.
- **Tabla francesa propia (`loan.py`):** las 3 tablas de la fixture coinciden renglón por renglón. En las 3, Σ amortización = capital, saldo final 0 y Σ pagos = capital + Σ intereses. El último pago ajusta centavos (156.41 contra 156.42; 154.53 contra 154.51; 964.24 contra 964.28).
- **Tarjeta:** el pago total un día antes y el mismo día dan 0.00, con datos y sin datos. Un día después, con la fixture capturada, da 65.98: mi cuenta es 1,500.00 × 30 días × 45.5 % / 360 × 1.16 = 65.975, que redondea a 65.98, igual que el oráculo. El mínimo da 200.00, 500.00 y 150.00 (topado al saldo). Sin datos, `missing_capture`. MSI: 333.33 + 333.33 + 333.34 = 1,000.00.
- `grep`: `highs|ortools` en go.mod y go.sum, 0; `net/http` fuera de las pruebas, 0; `float32|float64` fuera de las pruebas, 0; `reporte-compra-latex` en `oraculo/*.py`, 0. Los `cred_*.py` no leen la salida del motor: solo escriben `oraculo/salidas/s3.json`. El código Go que no es de pruebas no llama a Python ni lee salidas del oráculo; `canontest` solo lo importan las pruebas.

## Criterios
| Criterio | Estado | Evidencia mía |
|---|---|---|
| S3-01 | Cumple | `TestCanon85` 85 de 85; mi regresión 117 de 117 contra el dump y `model.json` |
| S3-02 | Cumple | `compara.py s3` con la salida guardada: 437 de 437 y diferencias 0 |
| S3-03 | Cumple | 28 de 28; mínimos 8,094.50 y 5,000.00; con la lectura omitida, B y `pending_confirmation` = true |
| S3-04 | Cumple | 16 de 16 fechas iguales a mi `sched`; 2027-02-28 y 2028-02-29 |
| S3-05 | Cumple | Σ capital repuesto = capital; K15: la parte de Karen del 29-oct repone 0 y la línea recupera 773.63 cuando se completa el pago, el 31-oct |
| S3-06 | Cumple | La tasa se guarda como 80050/309450; 4,000 → 1,034.74; plazos 6, 8, 10 y 12 → `no_rate`; se muestra «25.8685 %» (K12) |
| S3-07 | Cumple | promo-28 y promo-30 → 0.00; promo-27 y mínimo − 0.01 → no aplica y se usa la tasa normal; el mínimo exacto aplica |
| S3-08 | Cumple | 0.00 un día antes y el mismo día; 65.98 = mi cálculo = oráculo; sin datos, `missing_capture`; MSI exacto. La fórmula es PROVISIONAL con su bloqueo |
| S3-09 | Cumple | 3 de 3 tablas iguales a mi cálculo y cuadran; `cred_run_fixtures.py` código 0 |
| S3-10 | Cumple | Abonos 333.33/333.33/333.34; sin datos, `missing_capture`; con datos, `compara` 906/906 |
| S3-11 | Cumple | `TestUnificado` PASS. En el dump: orden creciente, 0 duplicados, Σ por quincena = Σ de pagos (7,966.75 y 7,941.75), créditos 1 a 4 |
| S3-12 | Cumple | 0 válidos con violación (mío y prueba de propiedad); motor = oráculo en los 3 casos; «faltan restricciones» → 0 evaluados y sin recomendación; 5 s sembrado en 00003; 0 s → 503 `unavailable`; peor caso del motor 125 ms; no hace falta el bloqueo de tiempo; el oráculo ya no mide reloj |
| S3-13 | Cumple | 12 óptimos = mi búsqueda exhaustiva = oráculo; 4 distintos en cada caso |
| S3-14 | Cumple | `score_bp` recalculado en los 46 válidos; empate = 1; `TestPuntaje` PASS |
| S3-15 | Cumple | `grep` 0 y 0 |
| S3-16 | Cumple | sha256 igual; independencia comprobada con `grep`; `reporte-compra-latex` 0 |
| S3-17 | Cumple | `float` fuera de las pruebas: 0; la tasa se guarda como razón exacta |
| S3-18 | Cumple | vet y gofmt limpios; total 94.0 %, y cada paquete ≥ 90.0 % |
| S3-19 | Cumple | 1,903.00/1,191.50 → 492.28/308.22; Axel + Karen = monto en los 16 pagos |
| S3-20 | Cumple | 5 bloqueos `S3-*`; §Supuestos en los reportes; marcas PROVISIONAL en las salidas |

§K: K12, K14, K15, K16, K17, K18, K43 (5 s iniciales, 503) y K44 (`comparator.timeout_seconds` configurable con trigger de validación) están aplicados. La propuesta de K31 está escrita como propuesta en `motor-S3-comparador-go.md`: mantener 5 s, con mínimo de 1 s y máximo de 30 s si S1 pide límites, y 7 días para `date_window_days`.

## Hallazgos
No hay hallazgos altos ni medios.

**B1 (baja; viene de r2 y sigue abierto) · `cred_escenarios.py` no compara su resultado con el `esperado` de cada fixture.** `_resumen` (líneas 231 a 259) solo cuenta los válidos con violación. El código de salida (línea 266) no revisa el `status` ni el `evaluated_count` de la carga contra lo esperado (`ok`, 3,870). Desde M2 la carga ya no puede salir en `timeout` por el reloj, así que el riesgo es menor.
Corrección: comparar `status` y `evaluated_count` (o la salida completa) de cada consulta con su `esperado` y salir con 1 si difieren. Ruta: `/home/deepsec/finanzas-hogar/oraculo/cred_escenarios.py`.

**B2 (baja; viene de r1 como H4 y sigue abierto) · `value_per_peso_milli` en milésimas produce empates en el objetivo 4.** En el caso 2 hay 3 válidos empatados en el mejor valor (30) y el desempate es el puntaje. El óptimo no cambia en ninguno de los 3 casos.
Corrección: usar escala ×10⁶ o la fracción exacta, por medio del contrato y del bloqueo `S3-objetivos-medida`. Rutas: `/home/deepsec/finanzas-hogar/docs/contratos/motores/S3-interfaz.md` §8.4 y `/home/deepsec/finanzas-hogar/backend/internal/motor_cred/comparador/objectives.go`.

**B3 (baja; viene de r1 como H5 y sigue abierto) · El caso 3 da las mismas 4 llaves de óptimo que el caso 1.** Con el `override`, el caso 3 queda con las mismas restricciones que el caso 1 (colchón 300,000, 2 vivas, deuda en cero 2027-02-15); solo cambian el monto y el tope, y por eso cambian los puntajes. Cumple S3-13 (4 distintos dentro de cada caso), pero aporta poca cobertura nueva.
Corrección: en una revisión futura de los casos, escribir antes de correr un caso 3 con restricciones efectivas distintas. Ruta: `/home/deepsec/finanzas-hogar/docs/estado/S3-escenarios.md` §4.

**B4 (baja; viene de r1 como H6 y H7 y sigue abierto) · Al contrato le faltan dos menciones.** La §2 no nombra `canontest/` (0 coincidencias) y la §8.2 no dice que `term_loan` no aporta candidatos. El contrato no cambió desde las 19:41.
Corrección: que el integrador del contrato aplique las solicitudes del núcleo y del comparador a §2, §8.2 y §11. Ruta: `/home/deepsec/finanzas-hogar/docs/contratos/motores/S3-interfaz.md`.

**B5 (baja; nuevo) · El contrato no dice que el oráculo no mide reloj.** M2 se resolvió con la opción (a) en `cred_escenarios.py` y quedó escrito en el docstring y en `S3-escenarios.md` §8. Pero §8.7 y §10.0 de `S3-interfaz.md` siguen sin decir qué hace el oráculo con `timeout_seconds` > 0. La solicitud 1 de `correccion-S3-r2.md` sigue sin aplicar. No hay contradicción: el contrato nunca le pidió al oráculo medir reloj, y la respuesta está en los entregables de S3. Por eso la clasifico como baja.
Corrección: agregar en §8.7 y §10.0 que, en el oráculo, `timeout_seconds` = 0 → `timeout` (determinista) y un valor positivo calcula completo; el tiempo real lo mide solo `TestRestriccionesTiempo`. Ruta: `/home/deepsec/finanzas-hogar/docs/contratos/motores/S3-interfaz.md`.

**B6 (baja; nuevo) · Referencia de sección equivocada en `S3-escenarios.md`.** La §8, primera viñeta, dice «líneas `S3-TIEMPO` de §6», pero los tiempos medidos están en la §5 («Tiempos medidos»); la §6 es «Resultado de la corrida contra la tabla a mano».
Corrección: cambiar «§6» por «§5». Ruta: `/home/deepsec/finanzas-hogar/docs/estado/S3-escenarios.md` §8.

## Estado de los anteriores
| Hallazgo r2 | Estado | Evidencia mía |
|---|---|---|
| M1 (media, S3-02 / S3-12): `oraculo/salidas/s3.json` sin el comparador | **Cerrado** | La salida guardada tiene 27 fixtures y 79 consultas. Su sha256 `b84c0db7…bd8f` es igual al de mi regeneración en `datos/tmp`. `compara.py s3` con el oráculo por defecto da 906/906 en las tres direcciones, 437 de 437, diferencias 0 y código 0; sin argumentos también da 0 |
| M2 (media, S3-12): el oráculo medía reloj en la carga | **Cerrado** | `compare()` ya no llama a `time.monotonic()`; solo `timeout_seconds` = 0 da `timeout` (líneas 141 y 142). Con el núcleo compartido por dos procesos, la carga tardó 7,828 ms y dio `ok` con 3,870 evaluados; el sha256 de la salida no cambió y `compara.py s3` dio diferencias 0. Lo que falta escribir en el contrato queda como B5 (baja) |
| B1 (baja) | Abierto | Sigue como B1 |
| B2, B3 y B4 (bajas) | Abiertos | Siguen como B2, B3 y B4 |
