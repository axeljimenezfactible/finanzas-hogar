# REVISION · motor-S2-puerta2 · ronda 1 de 3 (P2c)

Revisor: `revisor-motor-S2-r1` (independiente; no escribí nada de lo revisado y no corregí nada) · 11-oct-2026, noche · Evalúo solo contra `docs/criterios/S2.md` 1.2 completo y `docs/contratos/motores/S2-interfaz.md` 1.3 (cambios de la puerta 2), con `docs/estado/puerta2/RESPUESTAS.md` y `PUERTA-2.md` como mandato. Revisé el diff de `backend/internal/motor_fin/` (más `go doc`), de `oraculo/fin_*.py`, `fin_genera_fixtures.py`, las fixtures s2 nuevas o cambiadas y el reporte `docs/estado/p2c/motor-S2-puerta2.md`. Todas las verificaciones las corrí yo, con mis propios comandos y fixtures en `datos/tmp/rev-s2-r1/` (borrada al terminar).

## Veredicto
**NO APROBADO.** 0 altos, **1 medio** (M1: una prueba nueva pega una cifra de `model.json`) y 5 bajos. Las cifras cuadran: motor = oráculo = esperado en las 32 fixtures (3185 hojas), 159 de 159 hojas de `model.json`, 43 de 43 cruces, cobertura de `motor_fin` 89.2 %, `make verify-all` en verde. La regla PROVISIONAL de ventanas contiguas se aplica en los dos lados, y lo comprobé también con mis casos. Ninguna respuesta de Axel queda contradicha, y nada «sin recomendación» se decidió.

## Verificaciones (evidencia propia)
| # | Qué | Comando (desde la raíz salvo indicación) | Resultado |
|---|---|---|---|
| a1 | Pruebas y cobertura K25 | `cd backend && go test ./internal/motor_fin/ -count=1 -cover` | `ok … coverage: 89.2% of statements` (≥ 85 %) |
| a2 | Conteo de pruebas | `go test ./internal/motor_fin/ -count=1 -v` (conteo con awk) | 20 pruebas de primer nivel PASS y 294 subpruebas PASS; 0 FAIL y 0 SKIP; `grep -rn t.Skip internal/motor_fin` da 0 |
| a3 | gofmt | `gofmt -l internal/motor_fin cmd/motores` y `gofmt -l .` | vacío en los dos |
| a4 | vet | `go vet ./internal/motor_fin/ ./cmd/motores/` y `go vet ./...` | limpio en los dos |
| a5 | S2-13 sin float | `grep -rn 'float32\|float64' internal/motor_fin --include='*.go' \| grep -v _test.go \| wc -l` | 0 |
| a6 | Pruebas nuevas | `-run 'TestVentanasContiguas\|TestInstanteZ\|TestMetasEstado\|TestProyeccionCadencia\|TestCanon$\|TestOraculoFixtures'` | las 6 PASS |
| b1 | El generador es determinista y las fixtures salen de él | copia de `oraculo/` y `datos/semilla/` en `datos/tmp/rev-s2-r1/copia`, `python3 -I oraculo/fin_genera_fixtures.py` ahí, `diff -r` contra `oraculo/fixtures/s2` | «32 fixtures y 3 archivos de _datos escritos», sin ningún HALLAZGO; `diff -r` sin diferencias (32 fixtures) |
| b2 | run_fixtures reproduce `salidas/s2.json` | `python3 -I oraculo/fin_run_fixtures.py --fixtures oraculo/fixtures/s2 --out datos/tmp/rev-s2-r1/oraculo-s2.json`, luego `cmp` | «32 fixtures, 32 coinciden, 0 diferencias»; `cmp` idéntico byte a byte |
| b3 | Motor desde la raíz y compara | `go build -o datos/tmp/rev-s2-r1/motores ./cmd/motores`; `motores s2 dump --fixtures oraculo/fixtures/s2 --out …/motor-s2.json` (exit 0); `python3 -I oraculo/compara.py s2 --oraculo … --motor …` | 3185 de 3185 hojas en oráculo contra esperado, motor contra esperado y oráculo contra motor; `model.json` (ramas de s2) 159 de 159; cruces 43 de 43; provisionales 256 de 256; «diferencias: 0» |
| c1 | Regla 4 (contiguas) en los dos lados, con fixtures mías | 3 fixtures propias (`rev-ventanas-d1`, `rev-ventanas-d2`, `rev-metas-instantes`; 15 consultas, 78 hojas, cifras a mano), corridas con `fin_run_fixtures.py --fixtures`, `motores s2 dump --fixtures` y `compara.py s2 --fixtures` | Motor contra esperado: 78 de 78. Oráculo contra esperado: 77 de 78 (la que falta es B5, un instante, no las ventanas). Las ventanas cuadran en los dos lados: |
| c1.1 | hueco de 1 día (A 10-10…11-09, B 11-11…12-09; gasto 100001 a 2 meses) | — | `invalid_argument` en los dos |
| c1.2 | ventana de un día y hueco (d 10-10, r 10-12…11-09) | — | `invalid_argument` en los dos |
| c1.3 | hueco sin eventos **y sin devengo que lo cruce** (W1 10-01…10-31, W2 11-02…11-30; gasto de 1000 a 1 mes) | — | `invalid_argument` en los dos: se aplica la forma fuerte («toda lista con hueco»), no solo «hueco que cruza un devengo» |
| c1.4 | desordenadas y con hueco | — | `invalid_argument` en los dos |
| c1.5 | contiguas en el cambio de mes (O 10-10…10-31, N 11-01…11-30, D 12-01…2027-01-09) | — | se aceptan: 50001/100001/−50000, 50000/0/0 y 0/0/0 (devengo HALF_UP de 100001 entre 2: 50001 y 50000) |
| c1.6 | contiguas en el cambio de año (X 10-10…12-31, Y 2027-01-01…01-09) | — | se aceptan: X 50001/100001/−50000; Y 50000/0/0 (la parte 1 sigue a la ventana, caso 2 PROVISIONAL) |
| c1.7 | contiguas en febrero (W1 2026-10-01…2027-02-28, W2 2027-03-01…03-31) | — | se aceptan: 20900/20900/0 y 19900/19900/0 |
| c1.8 | meses con hueco `[2026-10, 2026-12]` | — | se aceptan (la regla 4 solo aplica a ventanas): oct 50001, dic 0 |
| c1.9 | A + C de la interfaz (fixture del repositorio `s2-ventanas-a-mas-c`) | b3 | `invalid_argument` en los dos; A, B y C contiguas con las cifras de M1 |
| c2 | Metas (Q2-16 y Q2-18 FIRMES, meta vencida PROVISIONAL), con mi fixture | c1 | Los dos lados dan lo mismo en: meta con fecha de hoy (`in_progress`, 1, 600); vencida ayer, en el mes de hoy (`overdue`, 1, 900); cumplida y vencida (`achieved`, 0, 0, plan vacío); con retiros y ahorro −200 (`progress_tenths_pct` −200, `progress_pct` −20, 3 × 400); vencida 2 meses antes con un depósito futuro (ahorro 0, `overdue`, 0, plan [2026-10, 1001]); y el resumen (1050, 300, 50) |
| c3 | Instantes (B1) con mis casos | c1 | 9 dígitos de fracción con «Z», 06:00Z, «t» y «z» minúsculas y desfase `+06:00`: los dos lados coinciden. `T24:00:00Z`: **difieren** (B5) |
| d | `model.json` y `verify3.json` intactos | `sha256sum` antes y después (`sha256sum -c`); `git diff --quiet` | `1a570666…f1bfd4`, `d86087f3…95041` (iguales a `MODEL_SHA` y `VERIFY3_SHA` del Makefile) y `data.json` `2fe7afae…31d31b`: OK antes y después; sin cambios en git; `grep -n reporte-compra-latex oraculo/*.py` da 0 |
| e | Independencia (S2-03) | `grep -rn 'backend\|motor-s2\|datos/tmp' oraculo/fin_*.py`; lecturas de archivos del oráculo | La única mención de `backend` es el docstring de `fin_comun.py:4` («no lee nada de backend/»). `fin_genera_fixtures.py` solo abre `oraculo/model.json` y `datos/semilla/data.json`. Las pruebas Go no escriben en `oraculo/`: el único `os.WriteFile` es el `--out` de `cmd/motores/main.go:76`. Forma distinta en cada lado: regla 4 con `ord()+1` en Go y `(inicio − fin).days != 1` en Python; meta vencida con `max(1, n)` en Go y una rama explícita en Python |
| g | Verificación global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | exit 0, «Verificación ampliada OK»; `K25 OK motor_fin 89.2 %`; s2, s3, s4 y s5 con `salidas` idénticas y 0 diferencias. La falla de s5 que cita el reporte ya no ocurre |

## Criterios S2 (`docs/criterios/S2.md` 1.2)
| Criterio | Estado | Evidencia |
|---|---|---|
| S2-01 (alta) | CUMPLE | `TestCanon` PASS: lee `model.json` en tiempo de prueba. Ver M1 para una prueba **nueva** que pega una cifra de `model.json` |
| S2-02 (alta) | CUMPLE | b3: 159 de 159 |
| S2-03 (alta) | CUMPLE en las fixtures | b1 a b3 y e; con mis fixtures, 1 divergencia fuera de contrato (B5); B2 sigue abierto |
| S2-04 (alta) | CUMPLE | d |
| S2-05 (alta) | CUMPLE | `TestSumaPartes` PASS (a2) |
| S2-06 (alta) | CUMPLE | `TestQuincenaFinDeMes` PASS; c1.5 a c1.7 (fin de mes, de año y febrero) |
| S2-07 (alta) | CUMPLE | `TestInvarianteDeudaViva` PASS |
| S2-08 (alta) | CUMPLE | `TestPromedios` PASS |
| S2-09 (alta) | CUMPLE | `s2-proyeccion-7-12-*`, `s2-proyeccion-csv` y `s2-liga-provisional` en b3; `provisional` null en los 12 meses (Q2-14 = A; Q2-21 a Q2-24 = C) |
| S2-10 (alta) | CUMPLE | `s2-flujo-quincenal` en b3; `provisional` null (Q2-14 = A). `variable_estimated_cents` y `savings_planned_cents` siguen null, con bloqueo abierto |
| S2-11 (alta) | CUMPLE | `s2-alertas-80-100` y `s2-presupuesto-scope` en b3; plan 0 FIRME (Q2-20 = A) |
| S2-12 (media) | CUMPLE | `s2-metas`, que ya compara `state`; `s2-metas-vencida`; c2 |
| S2-13 (alta) | CUMPLE | a5 |
| S2-14 (media) | CUMPLE (ver B5) | `TestZonaHoy` y `TestInstanteZ` PASS; c3 |
| S2-15 (alta) | CUMPLE | a1, a3 y a4 |
| S2-16 (media) | CUMPLE (ver B7) | El reporte tiene SUPUESTOS. Cada PROVISIONAL cita un bloqueo ABIERTO: `S2-ventanas` pregunta 4 (casos 1 y 2), `S2-metas-vencida`, `S2-ahorro-quincenal` y `S2-flujo-quincenal` pregunta 1 (Q2-12) |
| §3 No aceptable | **NO CUMPLE (M1)** | «Valores esperados pegados en el código de prueba en vez de leídos de model.json» |

Mandato de la puerta 2: Q2-01 = A se aplica (regla 4 PROVISIONAL en los dos lados, fixture A + C, comentarios). Caso 2: sin recomendación, así que sigue el comportamiento en uso con su marca, sin decidir. Q2-12: sin respuesta, `variable_estimated_cents` sigue null. Q2-13 = B: la base sigue bloqueada y el campo null. Q2-14 = A, Q2-15 = A, Q2-16 = A, Q2-17 = C (de pantalla), Q2-18 = A, Q2-19 = B, Q2-20 = A, Q2-21 a Q2-24 = C y Q2-25 a Q2-27 = A/H1/O1: todas reflejadas. **Ninguna contradicción.**

## Hallazgos
### M1 · media · S2 §3 «No aceptable» / interfaz §7 («ningún valor esperado se pega en el código»)
- **Ruta:** `backend/internal/motor_fin/cover_test.go:239` (`TestVentanasContiguas`).
- **Qué:** `eq(t, "[P1] al 2026-11-09", …, 97375)` pega como esperado la cifra canon `model.json:periodos[0].diferencia_acumulada` (973.75). El resto de la suite no pega ninguna cifra de `model.json`: `TestVentanasOrden` (l. 206) compara esa misma hoja con `mc(t, "periodos[%d].diferencia_acumulada")`, y un `grep` de 12 cifras canon en `*_test.go` solo encuentra 97375 en `engine_test.go:245`, donde es una entrada y no un esperado. Las otras cifras de la prueba (50000, 100000, −50000, 0 y 292125) son de la interfaz y de M1, no hojas de `model.json`: están bien.
- **Corrección sugerida:** cambiar el literal por `Cents(mc(t, "periodos[0].diferencia_acumulada"))`, con la misma forma que `TestVentanasOrden`.

### B5 · baja · S2-03 / S2-14 (nuevo; misma familia que B1 de r3)
- **Ruta:** `oraculo/fin_comun.py:48-54` (`instante`).
- **Qué:** con `calendario {now: "2026-10-31T24:00:00Z"}`, el motor da `invalid_date` y el oráculo acepta el instante: Python 3.14.6 `fromisoformat` lo lee como 2026-11-01T00:00Z, y la consulta no falla (`compara`: `DIF … hora-24.error … oraculo="__falta__" motor="invalid_date"`). RFC 3339 §5.6 limita la hora a 00-23, así que el motor tiene razón. El caso queda fuera de las fixtures y no toca ninguna cifra, pero el reporte da B1 por cerrado «igual que en el oráculo».
- **Corrección sugerida:** en la expresión de `instante`, que la hora sea `([01]\d|2[0-3])`, y agregar la consulta `hora-24` → `invalid_date` a `s2-zona-medianoche` en `fin_genera_fixtures.py`.

### B6 · baja · inconsistencia contrato ↔ código (S2-interfaz §1.1, §2.1 y §3.1)
- **Ruta:** `backend/internal/motor_fin/types.go:229` (`type Instant struct{ time.Time }`) y `input.go` (`Params.Now *Instant`).
- **Qué:** la tabla de §1.1 dice que el instante es `time.Time` en Go, y el código usa el tipo nuevo `Instant`. Es lo que cierra B1, pero el contrato no lo dice. S6 arma `Params` y lo descubriría al compilar.
- **Corrección sugerida:** solicitud al dueño del contrato para que §1.1 y §2.1 digan «`Instant` (envuelve `time.Time`; exige «Z»)». El implementador lo anota en sus solicitudes fuera de alcance.

### B7 · baja · S2-16 (texto de bloqueo desfasado)
- **Ruta:** `docs/estado/bloqueos/S2-metas-vencida.md`, párrafo «HECHO».
- **Qué:** sigue diciendo que el oráculo da `plan` vacío y `needed_monthly_cents` null, y que ninguna fixture trae una meta vencida. Hoy los dos lados dan la opción A y existe `s2-metas-vencida`. El bloqueo sigue ABIERTO, así que la marca PROVISIONAL es válida.
- **Corrección sugerida:** solicitud al dueño de los bloqueos para agregar una línea con el estado al cierre de P2c (los dos lados con A PROVISIONAL; fixture `s2-metas-vencida`).

## Estado de los hallazgos anteriores (`docs/estado/wf1/REVISION-S2-r3.md`)
| Hallazgo | Estado | Evidencia mía |
|---|---|---|
| M1 r3 · hueco que cruza un devengo y ventanas que no duran un mes (media) | **CERRADO** | Regla 4 en `ledger.go` y `fin_devengo.py` (c1.1 a c1.9). Pregunta 4 redactada y ABIERTA en `S2-ventanas.md`. §2.12 dice «contiguas, ascendentes y disjuntas» con la 4 PROVISIONAL. El caso 2 queda escrito como PROVISIONAL con sus cifras (consulta `d-r-caso-2` = 292125) |
| B1 r3 · desfase distinto de «Z» (baja) | **CERRADO** para el caso reportado | `+00:00`, `-06:00` y `+06:00` dan `invalid_date` en los dos lados. La hora 24 se abre aparte como B5 |
| B2 r3 · `invalid_date` contra `invalid_argument` (baja) | **ABIERTO** | Mi fixture `rev-b2`: `presupuesto_mes {month: "2026-13"}` y `flujo_quincenal {from: "2026-10-q3"}` dan `invalid_date` en el oráculo y `invalid_argument` en el motor. El reporte lo deja como PREGUNTA para §3.15, sin decidir: correcto |
| B3 r3 · solicitudes sin aplicar (baja) | **ABIERTO en parte** | Ya cerrado: `oraculo/README.md` dice 1.3 y 32 fixtures, y `bloqueos/INDICE.md` lista `S2-ventanas` pregunta 4, `S2-ahorro-quincenal` y `S2-metas-vencida`. Sigue abierto: la tabla de §5 de la interfaz dice «Total: 27 fixtures», no lista `s2-ventanas-a-mas-c` ni `s2-metas-vencida`, y pone la consulta A + C dentro de `s2-ventanas-sin-cubrir` (§5 y §8). El implementador la puso en una fixture propia y lo justifica en su Supuesto 2: con la semilla el hueco de A + C tiene pagos. Corrección: solicitud al dueño del contrato |
| B4 r3 · gofmt y vet fuera del paquete (baja) | **CERRADO** | `gofmt -l .` vacío y `go vet ./...` limpio en `backend/` |

## Observación (HECHO, no es hallazgo)
En P2c, un mismo agente cambió los dos lados (Go y Python). Eso se aparta de la regla de §0 de la interfaz («dos agentes… sin ver el trabajo del otro»), pero S2-03, tal como está escrito, se cumple: el oráculo no lee nada del motor, no copia su algoritmo, y las cifras nuevas están fijadas a mano en las fixtures. Mis 78 hojas calculadas a mano coinciden en el motor y, salvo B5, en el oráculo. La decisión de reparto es del orquestador.
