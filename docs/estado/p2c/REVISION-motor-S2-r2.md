# REVISION · motor-S2-puerta2 · ronda 2 de 3 (P2c)

Revisor: `revisor-motor-S2-r2` (independiente; no escribí nada de lo revisado y no corregí nada) · 11-oct-2026, noche · Evalúo solo contra `docs/criterios/S2.md` 1.2 completo y `docs/contratos/motores/S2-interfaz.md` 1.3, con `docs/estado/puerta2/RESPUESTAS.md` y `PUERTA-2.md` como mandato. Revisé el diff de `backend/internal/motor_fin/` (más `go doc`), el diff de `oraculo/fin_*.py` y de `fin_genera_fixtures.py`, las fixtures s2 nuevas (`s2-ventanas-a-mas-c` y `s2-metas-vencida`) y las cambiadas, el reporte `motor-S2-puerta2.md` y la corrección `correccion-motor-S2-r1.md`. Corrí todas las verificaciones yo mismo, con mis comandos y con fixtures y scripts propios en `datos/tmp/rev-s2-r2/`. Borré esa carpeta al terminar.

## Veredicto
**APROBADO.** 0 altos y 0 medios. El medio M1 de la ronda 1 quedó **CERRADO**, y lo comprobé con una mutación propia. Siguen abiertos 3 bajos de r1 (B5, B6 y B7) y 2 de r3 (B2 y parte de B3). Ninguno se le asignó a la corrección, y los arreglos de B6, B7 y B3 les tocan a otros dueños. No encontré hallazgos nuevos.

Las cifras cuadran. Motor, oráculo y esperado coinciden en las 32 fixtures (3185 hojas). También coinciden las 159 hojas de `model.json` y los 43 cruces. La cobertura de `motor_fin` es 89.2 % y `make verify-all` termina en verde. La regla 4 (PROVISIONAL) se aplica en los dos lados, y lo comprobé con 12 consultas mías de ventanas. Ninguna respuesta de Axel queda contradicha, y no se decidió nada marcado «sin recomendación».

## Verificaciones (evidencia propia)
| # | Qué | Comando (desde la raíz salvo indicación) | Resultado |
|---|---|---|---|
| a1 | Pruebas y cobertura (S2-15, K25) | `cd backend && go test ./internal/motor_fin/ -count=1 -coverprofile=…/s2.cov`; `go tool cover -func=… \| tail -1` | `ok … coverage: 89.2% of statements`; `total: (statements) 89.2%` (≥ 85 %) |
| a2 | Conteo | `go test ./internal/motor_fin/ -count=1 -v`, contado con awk; `grep -rn t.Skip internal/motor_fin \| wc -l` | 20 pruebas de primer nivel PASS y 294 subpruebas PASS; 0 FAIL y 0 SKIP; 0 `t.Skip` |
| a3 | gofmt | `gofmt -l internal/motor_fin cmd/motores` y `gofmt -l .` | vacío en los dos |
| a4 | vet | `go vet ./internal/motor_fin/... ./cmd/motores/` y `go vet ./...` | limpio en los dos |
| a5 | Pruebas por criterio | `-run` Canon, Suma, Quincena, Invariante, Promedio, Proyeccion, Flujo, Presupuesto, Meta, Zona, Oraculo, VentanasContiguas, InstanteZ | todas PASS. `TestCanon` trae 85 subpruebas `R-xx` PASS (86 con la de la semilla); Oraculo trae 191 subpruebas (159 hojas y 32 fixtures); Invariante 7 de 7 |
| a6 | S2-13 sin float | `grep -rn 'float32\|float64' internal/motor_fin --include='*.go' \| grep -v _test.go \| wc -l` | 0 |
| a7 | Cifras canon pegadas en pruebas (S2 §3) | script propio `canon_lits.py`: junta las 594 hojas de `model.json` que son centavos enteros (133 valores distintos) y busca cada literal de 4 cifras o más en `motor_fin/*_test.go` | Ningún esperado es una hoja canon. Las coincidencias son entradas que se inyectan (`100000`, `19900`, `500000` y otras) o números redondos que coinciden por azar, más `engine_test.go:245-246`, donde 97375 es una entrada de `promedios`. `292125` no es hoja de `model.json`: sale de la interfaz y de REVISION-S2-r3 M1 |
| a8 | Mutación de M1 | copia de `motor_fin` en `datos/tmp/rev-s2-r2/mut`, l. 239 cambiada a `periodos[2]`, `-run 'TestVentanasContiguas$'` | `FAIL … cover_test.go:239: [P1] al 2026-11-09: motor=97375 esperado=125869`: la aserción sí lee `model.json`. Una segunda mutación del esperado de A, B y C (`… 0 1`) también falla en la l. 233 |
| b1 | El generador es determinista | copia de `oraculo/` y `datos/semilla/` en `datos/tmp/rev-s2-r2/copia`; `python3 -I oraculo/fin_genera_fixtures.py` ahí; `diff -r` contra `oraculo/fixtures/s2` | «32 fixtures y 3 archivos de _datos escritos», 0 líneas con HALLAZGO; `diff -r` vacío |
| b2 | run_fixtures reproduce `salidas/s2.json` | `python3 -I oraculo/fin_run_fixtures.py --fixtures oraculo/fixtures/s2 --out datos/tmp/rev-s2-r2/oraculo-s2.json`; `cmp` contra `oraculo/salidas/s2.json` | «32 fixtures, 32 coinciden, 0 diferencias»; `cmp` idéntico |
| b3 | Motor desde la raíz y compara | `go build -o datos/tmp/rev-s2-r2/motores ./cmd/motores`; `datos/tmp/rev-s2-r2/motores s2 dump --fixtures oraculo/fixtures/s2 --out …/motor-s2.json` (exit 0); `python3 -I oraculo/compara.py s2 --oraculo … --motor …` | 3185 de 3185 en oráculo contra esperado, motor contra esperado y oráculo contra motor; `model.json` (ramas de s2) 159 de 159; cruces 43 de 43; provisionales 256 de 256; «diferencias: 0» |
| c1 | Regla 4 en los dos lados, con datos propios | fixture mía `rev-r2-ventanas` (semilla sin compras a crédito y con un solo gasto de contado de 90001 el 2026-11-20 a 3 meses; devengo a mano HALF_UP(90001/3) = 30000, 30000 y 30001); `fin_run_fixtures.py --fixtures`, `motores s2 dump --fixtures` y `compara.py s2 --fixtures` | 12 consultas: motor = oráculo = esperado en todas (abajo). El motivo de cada rechazo lo saqué de los dos lados con `evalua()` del oráculo y con una prueba Go propia en la copia (`Eval` → `e.Detail`) |
| c1.1 | hueco de 1 día (W1 11-01…11-30, W2 12-02…12-31) | — | `invalid_argument` en los dos, **por la regla 4** («van contiguas…» en Go y en Python) |
| c1.2 | hueco entre la 2.ª y la 3.ª ventana (W3 2027-01-05…) | — | `invalid_argument` por la regla 4 en los dos |
| c1.3 | A + C propio (A 11-15…12-14, C 2027-01-15…02-14; hueco sin eventos que cruza el devengo) | — | `invalid_argument` por la regla 4 en los dos |
| c1.4 | hueco de 1 día después de una ventana de 3 meses (W2 2027-02-16) | — | `invalid_argument` por la regla 4 en los dos |
| c1.5 | febrero: W1 termina el 2027-02-27 y W2 empieza el 03-01 (falta el 28) | — | `invalid_argument` por la regla 4 en los dos |
| c1.6 | desordenadas y contiguas; traslape de un día | — | `invalid_argument` por la regla 1 en los dos (la regla 4 no se adelanta a la 1) |
| c1.7 | contiguas del 15 al 14 (3 ventanas) | — | se aceptan: 30000/90001/−60001, 30000/0/−30001 y 30001/0/0 |
| c1.8 | contiguas cortas, caso 2 PROVISIONAL (D de 1 día, E de 5 días, F larga) | — | se aceptan con las mismas cifras: la parte k sigue a la ventana (regla en uso, sin decidir) |
| c1.9 | contiguas en el cambio de año (Y = 2027-01-01, de un día) | — | se aceptan: las mismas cifras |
| c1.10 | contiguas en año bisiesto (L1 …2028-02-28, L2 2028-02-29, L3 2028-03-01…) | — | se aceptan: las mismas cifras en los dos lados |
| c1.11 | meses con hueco `[2026-11, 2027-01]` | — | se aceptan (la regla 4 solo aplica a ventanas): nov 30000/90001/−60001; ene 30001/0/0 |
| c1.12 | A + C de la interfaz (fixture del repositorio `s2-ventanas-a-mas-c`) y `hueco-sin-eventos` y `d-r-caso-2` de `s2-ventanas-sin-cubrir` | b3 | `invalid_argument`, `invalid_argument` y 292125 en los dos |
| c2 | Metas (Q2-16 y Q2-18 FIRMES, meta vencida PROVISIONAL), fixture mía `rev-r2-metas-instante` | c1 | Los dos lados coinciden. Meta vencida en agosto, sin ahorro: `overdue`, 0, 1001 y plan [2026-10, 1001]. Meta que vence hoy: `in_progress`, 1, 3000. Meta con un depósito futuro (2026-10-20 > hoy): `in_progress`, 3 y plan 100/100/100 |
| c3 | Instante | c1 | `2026-11-01T05:59:59Z` da hoy 2026-10-31 en los dos. `2026-10-31T24:00:00Z`: el motor da `invalid_date` y el oráculo lo acepta: es la única diferencia de mis 2 fixtures (B5 de r1, sigue abierto) |
| d | `model.json` y `verify3.json` intactos (S2-04) | `sha256sum` antes de todo y `sha256sum -c` al final; `git diff --quiet`; `grep -n reporte-compra-latex oraculo/*.py` | `1a570666…f1bfd4` = `MODEL_SHA` (Makefile l. 187); `d86087f3…95041` = `VERIFY3_SHA` (l. 188); `data.json` `2fe7afae…31d31b`: OK; sin cambios en git; 0 líneas |
| e | Independencia (S2-03) | `grep -n 'backend\|motor_fin\|motor-s2\|datos/tmp\|cmd/motores' oraculo/fin_*.py`; `grep` de `open(`, `subprocess` y `os.system` | La única mención es el docstring de `fin_comun.py:4` («no lee nada de backend/»). Los únicos archivos que se abren son `model.json`, `data.json`, las fixtures, `_datos/semilla-oraculo.json` y la salida propia `salidas/s2.json` (`fin_verify.py`). Sin `subprocess`. Las formas son distintas: regla 4 con `ord()+1` en Go y `(inicio − fin).days != 1` en Python, y la meta vencida con `max(1, n)` en Go y una rama `faltan == 0` en Python |
| f | Qué cambió después de r1 | `stat -c '%y'` de `motor_fin/*.go`, `fin_*.py`, fixtures s2 y `salidas/s2.json` | Solo `cover_test.go` (2026-10-12 00:06, después de que se escribió r1 a las 00:04). Coincide con lo que declara la corrección |
| g | Verificación global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | exit 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»; `K25 OK motor_fin 89.2 %`; s2 «32 fixtures, 32 coinciden, 0 diferencias · salidas/s2.json idéntica», 3185 de 3185 y 159 de 159; `fin_verify` 159 de 159 en las dos vías |

## Criterios S2 (`docs/criterios/S2.md` 1.2)
| Criterio | Estado | Evidencia |
|---|---|---|
| S2-01 (alta) | CUMPLE | a5: `TestCanon` 85 de 85 `R-xx`; lee `model.json` al correr. a7 y a8: ya no queda ningún esperado canon pegado |
| S2-02 (alta) | CUMPLE | b3: 159 de 159 |
| S2-03 (alta) | CUMPLE en las fixtures | b1 a b3, c1 y e. Fuera de las fixtures quedan B5 y B2 (bajos) |
| S2-04 (alta) | CUMPLE | d |
| S2-05 (alta) | CUMPLE | `TestSumaPartes` PASS (a5) |
| S2-06 (alta) | CUMPLE | `TestQuincenaFinDeMes` (4 subpruebas) PASS; c1.5, c1.9 y c1.10 |
| S2-07 (alta) | CUMPLE | `TestInvarianteDeudaViva` 7 de 7 |
| S2-08 (alta) | CUMPLE | `TestPromedios` PASS |
| S2-09 (alta) | CUMPLE | `TestProyeccionDefecto` y `TestProyeccionCadencia` PASS; `provisional` null en los 12 meses (Q2-14 = A; Q2-21 a Q2-24 = C) |
| S2-10 (alta) | CUMPLE | `TestFlujoQuincenal` PASS; `variable_estimated_cents` y `savings_planned_cents` en null, con su bloqueo citado en `plan.go` |
| S2-11 (alta) | CUMPLE | `TestPresupuestoUmbrales` y `TestPresupuestoScope` PASS; plan 0 FIRME (Q2-20 = A) |
| S2-12 (media) | CUMPLE | `TestMetasAhorro` y `TestMetasEstado` PASS; `s2-metas` ya compara `state`; c2 |
| S2-13 (alta) | CUMPLE | a6 |
| S2-14 (media) | CUMPLE (ver B5) | `TestZonaHoy` y `TestInstanteZ` PASS; c3 |
| S2-15 (alta) | CUMPLE | a1, a3 y a4 |
| S2-16 (media) | CUMPLE (ver B7) | El reporte trae SUPUESTOS. Cada PROVISIONAL cita un bloqueo ABIERTO: `S2-ventanas` pregunta 4 (casos 1 y 2), `S2-metas-vencida`, `S2-ahorro-quincenal` y `S2-flujo-quincenal` pregunta 1 (Q2-12) |
| §3 No aceptable | CUMPLE | a7 y a8: el último esperado canon pegado (M1 de r1) ya se lee de `model.json`. Sin `t.Skip`, sin float y sin red. El oráculo no lee el motor |

Mandato de la puerta 2: Q2-01 = A se aplica. La regla 4 queda PROVISIONAL en los dos lados, hay fixture de A + C y la pregunta 4 está redactada y ABIERTA. El caso 2 no tiene recomendación, así que se queda el comportamiento en uso con su marca y sin decidir. Q2-12 no tiene respuesta. Q2-13 = B, Q2-14 a Q2-16 = A, Q2-17 = C, Q2-18 = A, Q2-19 = B, Q2-20 = A, Q2-21 a Q2-24 = C y Q2-25 a Q2-27 = A/H1/O1: todas están reflejadas en el código y en los comentarios. **No hay ninguna contradicción.**

## Hallazgos de esta ronda
No hay hallazgos nuevos.

## Estado de los hallazgos anteriores
| Hallazgo | Estado | Evidencia mía |
|---|---|---|
| **M1 r1 · media** · `cover_test.go:239` pegaba 97375 (= `periodos[0].diferencia_acumulada`) | **CERRADO** | La l. 239 ahora dice `eq(t, "[P1] al 2026-11-09", Cents(qc(…, "[0].cumulative_diff_cents")), mc(t, "periodos[0].diferencia_acumulada"))`, con la misma forma que `TestVentanasOrden` (l. 206). La mutación a8 falla con `esperado=125869`, así que sí lee `model.json`. a7: no queda ninguna otra cifra canon pegada como esperado. La ventana `a` de la prueba (2026-10-10 a 2026-11-09) es la de `periodos[0]` |
| B5 r1 · baja · `oraculo/fin_comun.py:48-54` acepta la hora 24 | **ABIERTO** (no asignado) | c3: `hora-24` da `invalid_date` en el motor y el oráculo lo acepta. Corrección: la de r1 (hora `([01]\d\|2[0-3])` en `instante` y la consulta `hora-24` en `s2-zona-medianoche`) |
| B6 r1 · baja · `types.go:229` `Instant` contra §1.1 y §2.1 («`time.Time`») | **ABIERTO** (no asignado; el dueño es el contrato) | `go doc ./internal/motor_fin Params` → `Now *Instant`; `S2-interfaz.md:28` sigue diciendo `time.Time`. Corrección: solicitud al dueño del contrato para que §1.1 y §2.1 digan «`Instant` (envuelve `time.Time`; exige «Z»)» |
| B7 r1 · baja · `docs/estado/bloqueos/S2-metas-vencida.md` HECHO desfasado | **ABIERTO** (no asignado; el dueño son los bloqueos) | La l. 7 sigue diciendo que el oráculo da `plan` vacío y null y que ninguna fixture trae una meta vencida, y `bloqueos/INDICE.md:405` dice «El motor y el oráculo difieren». Hoy los dos dan A (c2 y `s2-metas-vencida`). Corrección: la de r1 |
| B2 r3 · baja · `invalid_date` contra `invalid_argument` en mes y quincena mal formados | **ABIERTO** (PREGUNTA para §3.15, sin decidir) | Sin cambios de código desde r1 (f). El reporte lo deja como PREGUNTA, y es correcto |
| B3 r3 · baja · la tabla de §5 de la interfaz | **ABIERTO en parte** (dueño: el contrato) | `S2-interfaz.md:359` sigue con «Total: **27 fixtures**» y no nombra `s2-ventanas-a-mas-c` ni `s2-metas-vencida` (`grep` da 0 líneas). Corrección: solicitud al dueño del contrato |
| B4 r3 · baja · gofmt y vet fuera del paquete | **CERRADO** | a3 y a4 |

## Observación (HECHO, no es hallazgo)
La misma de r1. En P2c, un mismo agente cambió el lado Go y el lado Python, y eso se aparta del §0 de la interfaz («sin ver el trabajo del otro»). Aun así, S2-03 se cumple tal como está escrito: el oráculo no lee el motor, la forma de cada lado es distinta (e) y mis 2 fixtures propias, con cifras calculadas a mano, coinciden en los dos lados salvo B5. La decisión sobre el reparto le toca al orquestador.
