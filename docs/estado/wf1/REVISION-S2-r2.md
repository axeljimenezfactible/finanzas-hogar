# Revisión S2 · ronda 2 de 3 (revisor-S2-r2, WF-1, 11-oct-2026)

Revisor independiente: no escribí nada de lo revisado y no corregí nada. Evalué solo contra `docs/criterios/S2.md` 1.2, decisiones §B y §K y `docs/contratos/motores/S2-interfaz.md` 1.2. Todas las verificaciones las corrí yo. Los cálculos propios y las fixtures propias estuvieron en `datos/tmp/rev-s2r2/` y ya se borraron.

## Veredicto
**NO APROBADO.** Hay 1 hallazgo medio y 0 altos. El medio reabre en parte el M1 de la ronda 1: con ventanas desordenadas o traslapadas, el motor y el oráculo dan las mismas cifras incorrectas. Además hay 4 hallazgos bajos.

Los 15 criterios restantes cumplen con evidencia reproducida por mí. El motor coincide con el oráculo en las 30 fixtures (2320 de 2320 hojas) y en las 159 hojas de `model.json`. Lo mismo pasa con mis casos propios de K2, K8, K19, K20, K21 y K39: 552 comparaciones contra esperados que calculé por separado, con 0 diferencias en cifras.

## Estado de los hallazgos de la ronda 1
| Hallazgo r1 | Estado | Evidencia propia |
|---|---|---|
| M1 · ventanas (S2-03, S2-16) | **REABIERTO EN PARTE** (ver M1 abajo) | Lo que se corrigió funciona. Semilla con `{windows:[P2]}`: `invalid_argument` en las dos implementaciones. Ventana que empieza antes del primer evento sin eventos previos (2026-10-01 a 11-09): se admite y da 194750 / 97375 / 97375. Hueco sin eventos (A hasta 11-08, B desde 11-10): se admite y da R-45 y R-46. Ventana de 1 día (2026-10-10): 194750 / 0 / 194750. Bloqueo `S2-ventanas` abierto, con marca PROVISIONAL en `ledger.go:300` y en `fin_devengo.py:74`. **Pero** `[P2, P1]` todavía da P2 = 0 / 194750 / −292125, la misma cifra que reportó la ronda 1 |
| M2 · `accrual_months` e `interval_count` en 0 | **CERRADO** | En las dos implementaciones, un gasto con `accrual_months` 0, −1 o 25 da `invalid_argument`. Con 24 da 20833 en dic-2026 y en ene-2027 (500000/24 HALF_UP; mi cálculo da lo mismo). Con null u omitido equivale a 1 (500000 en dic, 0 en ene). Una compra sin gasto con `accrual_months` 0 da `invalid_argument`, y con 2 reproduce P1 (194750, diferencia 97375). Una partida con `interval_count` 25 da `invalid_argument`, y con 24 da 800000 solo en oct-2026. Una partida con `accrual_months` 0 da `invalid_argument` |
| M3 · prueba de propiedad S2-05 | **CERRADO** | `TestSumaPartes/fixtures`: «30 fixtures, 69 consultas sin error, 303 invariantes con diferencia 0» (Σ by_person 124, pagos + contado 61, Σ shares 93, axel + karen = total 14, axel + karen = monto 11). La semilla da 40 casos más |

## Tabla de criterios
| Criterio | Estado | Evidencia (corrida por mí) |
|---|---|---|
| S2-01 (alta) | CUMPLE | `go test -run Canon -v`: `TestCanon`, 86 subpruebas PASS (85 filas canon + semilla), 0 FAIL, 0 SKIP. `grep t.Skip` da 0 |
| S2-02 (alta) | CUMPLE | `compara.py s2`: «model.json (ramas de s2): 159 de 159». `fin_verify.py`: 159/159 de float contra la vía exacta y 159/159 de float contra `model.json` |
| S2-03 (alta) | CUMPLE en fixtures; ver M1, B1 y B2 | `fin_run_fixtures.py`: 30 fixtures, 30 coinciden, 0 diferencias. Mi salida es idéntica byte a byte (`cmp`) a `oraculo/salidas/s2.json`. `motores s2 dump`: exit 0. `compara.py`: oráculo contra esperado 2320/2320, motor contra esperado 2320/2320, oráculo contra motor 2320/2320, cruces 43/43, provisionales 928/928, 0 diferencias. `TestOraculoFixtures` + `TestOraculoHojas`: 189 subpruebas PASS. Independencia: `fin_*.py` no menciona `backend/`, `datos/tmp` ni `motor-` (solo un comentario de `fin_comun.py:4`). Los imports son de biblioteca estándar y de los propios `fin_*`. `motor_fin` no lee `oraculo/*.py` ni `salidas/`; las pruebas leen `model.json` y `fixtures/s2` |
| S2-04 (alta) | CUMPLE | sha256, antes y después de todo: `model.json` 1a570666…bfd4 y `verify3.json` d86087f3…5041. `grep reporte-compra-latex oraculo/*.py` da 0. `~/reporte-compra-latex/model/verify3.json` no cambió de mtime (2026-10-10 18:05) |
| S2-05 (alta) | CUMPLE | Ver M3. Además, en mis fixtures de partidas, Σ partes = total en todas las salidas |
| S2-06 (alta) | CUMPLE | `-run Quincena`: PASS. Mi `calendario` independiente (datetime) coincide en 10 fechas: 2026-11-30 q2 16–30 paga 30; 2026-12-31 q2 16–31 paga 30; 12-15 q1; 12-16 q2; dom 2026-10-18 = W42 (12–18); lun 10-19 = W43; 2027-01-01 = 2026-W53 (12-28 a 01-03); 2028-02-29 paga 29; 2027-02-28 y 02-15. Con un gasto del 2026-12-31 y otro del 12-05 (hoy 2027-01-05), `flujo_quincenal` pone 1000 en q1 y 12345 en q2. En `periodos` dic suma 13345 al contado y al devengado, y ene 0 |
| S2-07 (alta) | CUMPLE | `-run Invariante`: las 7 filas son iguales al criterio (292125 = 97375 + 194750 … 629344 = 125870 + 503474). Mi cálculo desde `model.json` (eventos propios, cruzados contra P1–P6) en 6 cortes nuevos (2026-10-09, 11-15, 12-20, 2027-04-30, 06-15, 12-31) da los mismos valores que las dos implementaciones, con la identidad cumplida |
| S2-08 (alta) | CUMPLE | `-run Promedio`: 3141/21988, 6492/45442, 8991/62935, 6282/43976. Míos con Fraction: 100000 en dic (31 d) da 3226/22581; 7 en 14 d da 1/4; ene-2027 efectivo, dic devengado y feb despensa coinciden; mensual 100/8 da 13 y 1310378/6 da 218396 |
| S2-09 (alta) | CUMPLE | 0 de 26 fechas de la semilla caen en los días 1–9. Mi proyección de 12 meses (desde los eventos de `model.json`) coincide en devengado y efectivo con las dos implementaciones. Con el interruptor apagado, `provisional` es null en los 12 meses; encendido, `S2-proyeccion-interruptor` aparece en los meses 7–12. Partidas weekly, monthly día 31, yearly 29-feb, fortnightly cada 2 y fortnightly con inicio el 20: el motor y el oráculo coinciden |
| S2-10 (alta) | CUMPLE | `s2-flujo-quincenal`: 59882 / 37493 por quincena (2026-10-q2…12-q1). Ximena no aparece en los mapas |
| S2-11 (alta) | CUMPLE | `s2-alertas-80-100`: 79000 → sin bandera; 79999 → `used_pct` 80 sin bandera (la frontera es exacta, no redondeada); 80000 → warning; 99000 → warning; 100000 → limit; 101000 → limit |
| S2-12 (media) | CUMPLE | Meta 2: (629344 − 314672) / 7 = 44953 × 6 + 44954 = 314672 al centavo. Meta 1: 6 × 100000. La meta cumplida no tiene plan |
| S2-13 (alta) | CUMPLE | `grep -rn 'float32\|float64' … --include='*.go' \| grep -v _test.go` da 0 |
| S2-14 (media) | CUMPLE; ver B1 | `-run Zona`: PASS. `now` 2026-12-01T05:59:59Z da 11-30, 06:00:00Z da 12-01 y 2027-03-01T05:30:00Z da 02-28, igual en las dos |
| S2-15 (alta) | CUMPLE; ver B4 | `gofmt -l internal/motor_fin cmd/motores` vacío; `go vet ./internal/motor_fin/... ./cmd/motores/...` limpio; cobertura `total: (statements) 89.1%`; `cmd/motores` 78.4 %; 308 PASS en `motor_fin` |
| S2-16 (media) | **NO CUMPLE** (M1) | Bloqueos S2-*: flujo-quincenal, metas, metas-resumen, presupuesto-total, proyeccion-interruptor y ventanas. La regla del orden de las ventanas cambia cifras y no está escrita ni bloqueada |

## Hallazgos
### M1 · media · S2-03 / S2-16 (interfaz §2.12 y §3.2; reabre en parte el M1 de r1)
**Qué pasa.** Con ventanas que no van en orden cronológico, o que se traslapan, el motor y el oráculo dan las mismas cifras, y esas cifras son incorrectas. La regla PROVISIONAL A de `S2-ventanas` no lo impide, porque solo revisa eventos antes de la primera ventana o en huecos. Casos con la semilla:
- `[P2, P1]` da P2 = devengado 0, efectivo 194750 y diferencia −292125. Según §3.2 («desde el primer evento … hasta `end`») y R-46, debería dar 194750 / 194750 / 97375. Es la cifra que la ronda 1 dio por corregida.
- `[P3, P1, P2]` da devengado de P1 = 446487: la segunda mitad de la despensa del 10-dic-2026 se devenga en la ventana de octubre. Las diferencias salen en −263631 y 600849.
- Traslape P1 + X (2026-10-25 a 12-09): el pago del 10-30 se cuenta en las dos ventanas (X tiene efectivo 292125 con 3 pagos).
- Ventana contenida A (10-10 a 01-09) ⊃ B (11-10 a 12-09): B da devengado 446487 y diferencia 600849.

En modo meses no pasa: `[2026-12, 2026-11]` da 125869 y 97375. Las dos implementaciones leen «n ventanas consecutivas» como consecutivas en el orden recibido, y la diferencia acumulada como suma corrida en ese orden. Ese supuesto cambia cifras y no está escrito en la interfaz ni en el bloqueo. El bloqueo además dice que con A «la cifra no puede salir incompleta en silencio».

**Rutas:** `backend/internal/motor_fin/ledger.go:283-316` (`uncovered`, `checkCovered`, `windowed`); `oraculo/fin_devengo.py:70-95` (`sin_cubrir`, `resumen`); `docs/estado/bloqueos/S2-ventanas.md`; `docs/contratos/motores/S2-interfaz.md` §2.12 y §3.2.

**Corrección sugerida:**
1. En las dos implementaciones, rechazar con `invalid_argument` (PROVISIONAL, `S2-ventanas`) unas ventanas que no estén en orden cronológico estricto o que se traslapen, es decir, cuando `start[i] ≤ end[i−1]`. Es coherente con la opción A ya aplicada.
2. Agregar a `s2-ventanas-sin-cubrir` (vía `fin_genera_fixtures.py`) las consultas `[P2,P1]`, `[P3,P1,P2]`, el traslape P1 + X y la ventana contenida, todas con error esperado.
3. Agregar la pregunta del orden y el traslape a `S2-ventanas` con sus opciones: rechazar, u ordenar por fecha para calcular y devolver en el orden recibido.
4. Pedir al dueño de la interfaz que §2.12 diga «ventanas en orden ascendente y disjuntas».

### B1 · baja · S2-03 / S2-14 (interfaz §1.1)
Un instante con desfase distinto de «Z» lo acepta el motor y el oráculo lo rechaza con `invalid_date`. Ejemplos: `2026-10-31T05:59:59+00:00` da en el motor today 2026-10-30, y `2026-10-31T00:30:00-06:00` da 2026-10-31. §1.1 dice «RFC 3339 en UTC con "Z"». La causa es que `*time.Time` usa el RFC 3339 completo (`input.go:35`, `eval.go:266`), mientras `fin_comun.py:47` exige `Z`. El día que da el motor es correcto y ninguna fixture lo cubre. **Corrección:** que el motor rechace todo lo que no termine en `Z` (o que la interfaz admita desfases y el oráculo los acepte), y agregar una consulta de error a `s2-zona-medianoche`.

### B2 · baja · S2-03 (interfaz §3.15)
Un mes o una quincena mal formados (`2026-1`, `2026-13`, `2026-10-01`, `202610`, `2026-10-q3`, `2026-10-Q1`, `2026-10-q0`) dan `invalid_argument` en el motor (`types.go:110-160`) y `invalid_date` en el oráculo (`fin_comun.py:57-70`). §3.15 enumera los códigos pero no define cuándo va `invalid_date`. **Corrección:** fijar el código en la interfaz (§3 dice «tipo o rango inválido → invalid_argument»), alinear la implementación que no siga ese código y agregar una fixture de error.

### B3 · baja · S2-16 (consistencia de artefactos)
Siguen sin aplicarse las solicitudes 1 a 5 de `correccion-S2-r1.md`:
- La interfaz 1.2 §2.12 y §3.2 no cita `S2-ventanas`, y §3.2 sigue diciendo «incluye periodos anteriores no pedidos» también para ventanas.
- La tabla de §5 no lista las 3 fixtures nuevas.
- `oraculo/README.md:22,30` dice 27 fixtures.
- `docs/estado/bloqueos/INDICE.md` no lista `S2-ventanas`.
- `motor-S2-go.md` §Supuestos 1 y 2 están desactualizados.

**Corrección:** que el integrador o el dueño de la interfaz apliquen esas solicitudes.

### B4 · baja · S2-15 (fuera del paquete; se arrastra de B1 de r1)
`gofmt -l .` en `backend/` lista `internal/precios/seed.go`, `internal/precios/types.go`, `internal/recetas/num.go` e `internal/recetas/types.go`. `go vet ./...` falla con `internal/recetas/types.go:229:13: undefined: Nutrition`. No son de S2, pero el comando (a) de la tarea no sale limpio en el repositorio. **Corrección:** que los dueños de `precios` y `recetas` lo resuelvan.

## Observación (sin hallazgo)
La regla de cobertura de ventanas (`uncovered` y `sin_cubrir`) la escribió el mismo agente (`correccion-S2-r1`) en los dos lados, y es la misma lógica de pertenencia. Como es una regla de validación que fija el bloqueo, no un algoritmo de cálculo, no la cuento como independencia rota. Conviene que la regla de orden de M1 la escriba cada lado por separado.

## Hallazgos de los reportes (g)
`motor-S2-go.md` §Hallazgos dice «Ninguna cifra canon difiere». `oraculo-S2.md` dice «0 hallazgos: ninguna cifra fijada difiere del oráculo». El generador de fixtures da 0 líneas `HALLAZGO`. No encontré ninguna cifra ajustada para que cuadre: mis esperados propios (canon, deuda viva en cortes nuevos, promedios, proyección, alcance de presupuesto, calendario) coinciden sin intervención.
