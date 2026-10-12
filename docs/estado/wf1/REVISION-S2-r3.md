# REVISION-S2-r3 · Motor financiero S2 · WF-1 · 11-oct-2026

Revisor independiente `revisor-S2-r3` (ronda 3 de 3). Evalúo solo contra `docs/criterios/S2.md` 1.2, decisiones §B y §K y `docs/contratos/motores/S2-interfaz.md`. No corregí nada. Reejecuté cada verificación con mis comandos. Usé mis propias fixtures en `datos/tmp/rev-s2-r3/` (5 fixtures, 35 consultas, ya borradas) y mis cálculos con `Fraction` y `Decimal`.

## Veredicto
**NO APROBADO.** Hay 1 hallazgo medio y 0 altos. El M1 de r2 (ventanas desordenadas, traslapadas o contenidas) queda **CERRADO**. El hallazgo medio nuevo es de la misma familia: con un hueco entre ventanas que no tiene eventos fechados pero sí lo cruza un devengo de varios meses, la parte del hueco se devenga en la siguiente ventana declarada. Esa regla cambia cifras y no está escrita ni en la interfaz ni en el bloqueo. Además, el bloqueo dice que las ventanas deben ser contiguas, y ninguna de las dos implementaciones lo exige. Siguen abiertos los 4 bajos de r2 (B1–B4).

## Estado de los hallazgos de la ronda 2
| Hallazgo r2 | Estado | Evidencia propia |
|---|---|---|
| M1 · orden y traslape de ventanas (S2-03, S2-16) | **CERRADO** | Las 5 consultas nuevas de `s2-ventanas-sin-cubrir` esperan `invalid_argument`, y el oráculo y el motor las dan. Con mi fixture `rev-orden` (semilla), el motor y el oráculo dan `invalid_argument` en los 6 casos: `[P1,P3,P2]`, `[B,A]` (A ⊃ B al revés), `[P1,P1]`, `[P2,P1]` y `[P1, Z invertida, P3]`. El control `[P1,P2,P3]` da 194750/97375/97375, 194750/194750/97375 y 251737/223243/125869 (= R-45…R-47). La regla lleva la marca PROVISIONAL que cita `S2-ventanas` en `ledger.go:478-487` y en `fin_devengo.py:82-94`, y la pregunta 3 está en el bloqueo. Cada lado está escrito con su propia forma (anterior inmediata contra máximo acumulado). Pero ver el M1 de r3 |
| B1 · desfase distinto de «Z» | ABIERTO (baja) | `now` = `2026-10-31T05:59:59+00:00`: el oráculo da `invalid_date` y el motor da today 2026-10-30. Con `Z` las dos dan 2026-10-30 |
| B2 · `invalid_date` contra `invalid_argument` | ABIERTO (baja) | Con `presupuesto_mes` y mes `2026-13`, y con `flujo_quincenal` y quincena `2026-10-q3`: el oráculo da `invalid_date` y el motor `invalid_argument` |
| B3 · solicitudes sin aplicar | ABIERTO (baja) | La interfaz §2.12 y §3.2 no cita `S2-ventanas` ni la regla «ascendentes y disjuntas». `oraculo/README.md:22,30` sigue diciendo 27 fixtures. `bloqueos/INDICE.md` no lista `S2-ventanas`. `motor-S2-go.md` §Supuestos 2 está desactualizado |
| B4 · gofmt y vet fuera del paquete | ABIERTO (baja) | `gofmt -l .` lista 4 archivos de `precios` y `recetas`. `go vet ./...` da `internal/recetas/types.go:229:13: undefined: Nutrition` |

## Tabla de criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S2-01 (alta) | CUMPLE | `go test -v`: `TestCanon` con 86 subpruebas PASS (85 filas y la semilla), 0 FAIL y 0 SKIP. `grep t.Skip` da 0 |
| S2-02 (alta) | CUMPLE | Con mis salidas, `compara.py s2` da «model.json (ramas de s2): 159 de 159» y cruces 43/43. `TestOraculoHojas`: 159 subpruebas PASS |
| S2-03 (alta) | CUMPLE en fixtures; ver M1 | `fin_run_fixtures.py --out <mío>`: 30 fixtures, 30 coinciden. Mi salida es idéntica (`cmp`) a `salidas/s2.json`. `motores s2 dump` sale con exit 0 y su salida es idéntica byte a byte a la del oráculo. `compara.py`: 2320/2320 en las tres comparaciones, provisionales 933/933 y 0 diferencias. En mis 5 fixtures, el oráculo y el motor dan lo mismo en 35 de 35 consultas. El generador, corrido en una copia, reproduce las 30 fixtures y `_datos/` byte a byte (`diff -rq` vacío). Independencia: `fin_*.py` solo menciona `backend/` en el comentario de `fin_comun.py:4` y solo importa biblioteca estándar y `fin_*`. El código no de prueba de `motor_fin` solo importa biblioteca estándar. Las pruebas leen `model.json` y `fixtures/s2`, no `salidas/` ni `*.py` |
| S2-04 (alta) | CUMPLE | sha256 antes y después: `model.json` 1a570666…bfd4 y `verify3.json` d86087f3…5041. `grep reporte-compra-latex oraculo/*.py` da 0. El mtime de `~/reporte-compra-latex/model/verify3.json` no cambió (2026-10-10 18:05:44) |
| S2-05 (alta) | CUMPLE | Log de `TestSumaPartes`: «30 fixtures, 69 consultas sin error, 303 invariantes con diferencia 0» y «40 casos». En mi proyección de abr-2027, devengado 176751 + 74986 = 251737 y efectivo 176752 + 74986 = 251738 |
| S2-06 (alta) | CUMPLE | `TestQuincenaFinDeMes`: 4 subpruebas PASS. Calendario: 2026-12-31 → q2 del 16 al 31 con pago el 30 y semana W53 (12-28 a 01-03); 2026-11-30 → q2 del 16 al 30 con pago el 30; 2027-02-14 → q1 con pago el 15 y semana 02-08 a 02-14; 2026-10-14 → semana 10-12 a 10-18. En ene-2027 (31 días), el gasto del 01-31 cae en q2 (5321 = 1000 + 4321) y el del 01-15 en q1 (777). `periodos` de enero: gastos de contado 25998 = 19900 + 6098 |
| S2-07 (alta) | CUMPLE | `TestInvarianteDeudaViva` da las 7 filas del criterio. Corte propio 2027-01-20, a mano con §3.4 (meses cerrados): devengado 389500 + 251737 = 641237; efectivo 661136; saldo 912874 − 641237 = 271637; deuda 125868 + 125870 = 251738 = −19899 + 271637. El motor y el oráculo dan exactamente eso. El 2027-04-10 da 629344 = 125870 + 503474 |
| S2-08 (alta) | CUMPLE | El `Decimal` del criterio imprime `219.88 629.35`. Log: 3141/21988, 6492/45442 y 8991/62935. Míos con `Fraction`: dic-2026 efectivo 7201/50410, dic devengado 8121/56844, feb-2027 devengado 8991/62934 contra efectivo 62935. El motor y el oráculo coinciden |
| S2-09 (alta) | CUMPLE | oct-2026…mar-2027 da las 12 cifras de R-45…R-50 y las diferencias 97375, 97375, 125869, 125870, 125869 y 125870. Con el interruptor apagado: abr 251737/251738 (04-15 + 04-30), may 251737/251736, jun 0/125870, jul–sep 0/0, sin despensa nueva y con `provisional` null |
| S2-10 (alta) | CUMPLE | `TestFlujoQuincenal` PASS. Mi flujo de ene-2027 q1 da pagos de crédito axel 88375 y karen 37493 (R-19). Ximena no aparece |
| S2-11 (alta) | CUMPLE | `TestPresupuestoUmbrales` PASS. `s2-alertas-80-100` entra en las 2320/2320 |
| S2-12 (media) | CUMPLE | `TestMetasAhorro` PASS. `s2-metas` entra en las 2320/2320 |
| S2-13 (alta) | CUMPLE | `grep 'float32\|float64' --include='*.go' \| grep -v _test.go` da 0 |
| S2-14 (media) | CUMPLE; ver B1 | `TestZonaHoy` PASS. `2026-10-31T05:59:59Z` da 2026-10-30 en las dos |
| S2-15 (alta) | CUMPLE; ver B4 | `gofmt -l internal/motor_fin cmd/motores` vacío. `go vet` de los dos paquetes limpio. Cobertura de `motor_fin`: 89.2 % (`cmd/motores` 78.4 %). 309 PASS, 0 FAIL |
| S2-16 (media) | **NO CUMPLE** (M1) | Bloqueos S2-*: flujo-quincenal, metas, metas-resumen, presupuesto-total, proyeccion-interruptor y ventanas. El devengo que cruza un hueco cambia cifras y no está escrito ni bloqueado |

## Hallazgos
### M1 · media · S2-16 / S2-03 (interfaz §2.12 y §3.2; bloqueo `S2-ventanas`)
Con ventanas ascendentes y disjuntas que dejan un hueco sin eventos fechados, un costo de n meses sigue devengándose en «n ventanas consecutivas» de la lista. La parte que corresponde al hueco cae en la siguiente ventana declarada. El motor y el oráculo coinciden en esas cifras, pero la regla no está escrita, y en modo meses la misma parte cae en el mes del hueco.
- Datos propios: la semilla sin compras a crédito, más un gasto de contado de 100000 el 2026-10-10 con `accrual_months` 2.
  - Ventanas A (2026-10-10 a 11-09) y C (2026-12-10 a 2027-01-09): A da 50000/100000/−50000 y C da **devengado 50000**, efectivo 0 y diferencia 0.
  - Con A, B y C contiguas, C da devengado **0** (B lleva los 50000).
  - En modo meses `[2026-10, 2026-12]`, diciembre da devengado **0**: la mitad va a noviembre, que no se pidió.
- La opción A del bloqueo dice: «quien llama tiene que declarar ventanas contiguas desde el primer evento». Pero `uncovered`/`sin_cubrir` solo rechazan eventos *fechados* en el hueco, y admiten un hueco que cruza un devengo.
- La misma causa aparece con ventanas contiguas que no duran un mes: la diferencia acumulada al mismo `end` depende de cómo se corten las ventanas. `[P1]` da 97375 al 2026-11-09. `[d = 2026-10-10, r = 2026-10-11 a 11-09]` da 292125 al 2026-11-09, porque la segunda mitad de la despensa 1 cae en r. §3.2 define esa diferencia «desde el primer evento… hasta `end`».
- Rutas: `backend/internal/motor_fin/ledger.go:259-281` (`windowed`) y `:287-298` (`uncovered`); `oraculo/fin_devengo.py:71-80` (`sin_cubrir`) y `:110-118` (`k0 + j`); `docs/estado/bloqueos/S2-ventanas.md:16,21`; `docs/contratos/motores/S2-interfaz.md:147,169`.
- **Corrección sugerida:**
  1. Agregar a `S2-ventanas` una pregunta 4: con ventanas, ¿el devengo «mitad por mes» sigue a la ventana o al mes, y qué pasa con un hueco sin eventos que cruza un devengo y con ventanas que no duran un mes? Poner opciones con sus cifras.
  2. Mientras se decide, aplicar en los dos lados (cada uno con su propia forma) una regla PROVISIONAL coherente con la consecuencia escrita de A. Por ejemplo, rechazar con `invalid_argument` toda lista con hueco (`start[i] ≠ end[i−1] + 1`), o al menos todo hueco que cruce un devengo.
  3. Agregar la consulta A + C de arriba a `s2-ventanas-sin-cubrir`.
  4. Pedir al dueño de la interfaz que §2.12 diga «contiguas, ascendentes y disjuntas (PROVISIONAL, S2-ventanas)».

  Si se prefiere conservar el comportamiento actual, hay que escribirlo como regla PROVISIONAL, con las cifras, en el bloqueo y en la solicitud a la interfaz.

### B1 · baja · S2-03 / S2-14 (se arrastra de r2)
El motor acepta un desfase distinto de «Z» y el oráculo lo rechaza (`input.go:35`, `eval.go:266` contra `fin_comun.py:47`). **Corrección:** la de r2.

### B2 · baja · S2-03 (se arrastra de r2)
Un mes o una quincena mal formados dan `invalid_argument` en el motor (`types.go:110-160`) y `invalid_date` en el oráculo (`fin_comun.py:57-70`). **Corrección:** la de r2.

### B3 · baja · S2-16 (se arrastra de r2; ahora incluye la regla de orden)
Faltan aplicar las solicitudes de `correccion-S2-r1.md` y de `correccion-S2-r2.md`:
- la interfaz §2.12 y §3.2 con `S2-ventanas` y la regla «ascendentes y disjuntas»;
- la tabla de §5 con las fixtures nuevas y las 9 consultas;
- `oraculo/README.md:22,30` (dice 27);
- `bloqueos/INDICE.md` sin `S2-ventanas`;
- `motor-S2-go.md` §Supuestos 2.

**Corrección:** que las aplique el integrador o el dueño de la interfaz.

### B4 · baja · S2-15 (fuera del paquete)
`gofmt -l .` y `go vet ./...` no salen limpios por `internal/precios` y `internal/recetas`. **Corrección:** que las resuelvan sus dueños.

## Hallazgos de los reportes (g)
- `motor-S2-go.md` §Hallazgos dice «Ninguna cifra canon difiere».
- `oraculo-S2.md` declara sus supuestos PROVISIONALES con su bloqueo.
- El generador corrido en copia no da ninguna línea `HALLAZGO` y reproduce las fixtures byte a byte.

No encontré cifras ajustadas: mis esperados propios coinciden sin intervención (canon, deuda viva en el corte 2027-01-20, promedios de dic y feb, proyección 7–12, K2 `onward`/`this_month`/omitido, calendario de 31, 30 y 28 días y semana). Las cifras del M1 no son un ajuste: las dos implementaciones aplican la misma lectura no escrita.

## Casos propios de K2, K8, K19, K20, K21 y K39 (f)
- **K2.** Renglones `[1: 200000 de 2026-10-01 a 2026-12-31, 2: 300000 desde 2027-01-01]`, con cambio en nov a 250000:
  - `onward` cierra 1 el 10-31, abre 3 desde el 11-01 sin fin y elimina el 2.
  - `this_month` deja 3 solo en nov y 4 (200000) del 12-01 al 12-31, y conserva el 2.
  - Sin `scope` el resultado es igual a `onward` (preseleccionado).
- **K8.** Ver S2-06.
- **K19.** Ver S2-07.
- **K20 y K39.** Ver S2-08.
- **K21.** Ver S2-09 (apagado por defecto: `repeat_cadence_enabled` false).
