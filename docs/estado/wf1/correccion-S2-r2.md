# correccion-S2-r2 · WF-1 · 11-oct-2026

Corrige el único hallazgo alto o medio de `docs/estado/wf1/REVISION-S2-r2.md`: **M1** (media; S2-03 y S2-16; reabre en parte el M1 de r1). Los bajos B1–B4 no estaban en la tarea y no se tocaron. No reabre decisiones de §A–§K.

## Hecho
- **Causa real (los dos lados):** el motor (`ledger.go`, `windowed` y `periods`) y el oráculo (`fin_devengo.py`, `resumen`) leen «n ventanas consecutivas» (§2.12) como consecutivas en el orden recibido, y calculan la diferencia acumulada como suma corrida en ese orden. La regla A + H1 (cobertura) no lo impedía: con la semilla, `[P2, P1]`, `[P3, P1, P2]`, P1 + X y A ⊃ B tienen todos sus eventos dentro de alguna ventana pedida. Por eso las dos implementaciones coincidían en cifras que estaban mal. No se ajustó ningún número: las dos rechazan ahora ese caso.
- **Regla nueva PROVISIONAL (bloqueo `S2-ventanas`, pregunta 3, opción O1):** una lista de ventanas en la que alguna no empieza después del fin de la anterior (`start[i] ≤ end[i−1]`) da `invalid_argument`. Esto incluye las desordenadas, las traslapadas (aunque sea un día) y las contenidas. Cada lado la escribe a su manera, como sugirió la observación de la revisión:
  - Motor (`periods`, dentro del mismo recorrido que ya validaba inicio ≤ fin): compara cada inicio con el fin de la ventana anterior. Mensaje: `ventana "X" empieza el … y la anterior "…" termina el …: van ascendentes y disjuntas (PROVISIONAL, S2-ventanas)`.
  - Oráculo (función nueva `desordenadas`, que `resumen` llama antes de `sin_cubrir`): compara cada inicio con el fin más tardío de todas las anteriores. Con ventanas válidas (inicio ≤ fin) las dos formas son equivalentes.
- **Fixtures:** `fin_genera_fixtures.py` (único autor) agrega 5 consultas con error esperado a `s2-ventanas-sin-cubrir`: `desordenadas-p2-p1`, `desordenadas-p3-p1-p2`, `traslape-p1-x` (X 2026-10-25 a 2026-12-09), `traslape-un-dia` (Y empieza el 2026-11-09, último día de P1; frontera `start = end` anterior) y `contenida-a-b` (A 2026-10-10 a 2027-01-09, B 2026-11-10 a 2026-12-09). El campo `provisional.que` y la descripción citan las dos reglas. Siguen siendo 30 fixtures; esa pasa de 4 a 9 consultas.
- **Prueba Go `TestVentanasOrden`** (`cover_test.go`): los 5 casos dan `invalid_argument` y el detalle contiene «ascendentes y disjuntas». Así se demuestra que los rechaza la regla de orden y no la de cobertura. El control P1, P2 y P3 contiguas reproduce R-45…R-47 (devengado y diferencia), leídos de `model.json`.
- **Bloqueo** `docs/estado/bloqueos/S2-ventanas.md`: nueva sección «Pregunta 3 · orden y traslape» con las cifras de la revisión, la causa y dos opciones: O1 rechazar (aplicada) u O2 ordenar por fecha para calcular y devolver en el orden recibido. En las dos se siguen rechazando los traslapes. También tiene el impacto y cómo seguí.

## Archivos
- `backend/internal/motor_fin/ledger.go` (modificado: regla de orden en `periods`; comentario de `windowed`).
- `backend/internal/motor_fin/cover_test.go` (modificado: import `fmt`, prueba `TestVentanasOrden`).
- `oraculo/fin_devengo.py` (modificado: `desordenadas` y su llamada en `resumen`).
- `oraculo/fin_genera_fixtures.py` (modificado: `ERR_ORDEN` y 5 consultas en `s2-ventanas-sin-cubrir`).
- `oraculo/fixtures/s2/s2-ventanas-sin-cubrir.json` (regenerado; las otras 29 fixtures y `_datos/` quedan idénticas byte a byte: `diff -rq` solo lista esta).
- `oraculo/salidas/s2.json` (regenerado por `fin_run_fixtures.py`; sha256 `b5adb6b0…1af6`; contra la salida anterior `2f39f0a4…0daa` solo cambia el caso `s2-ventanas-sin-cubrir`).
- `docs/estado/bloqueos/S2-ventanas.md` (sección agregada; la tarea lo pide explícitamente).
- `docs/estado/wf1/correccion-S2-r2.md` (este reporte).
- Sin cambios: `backend/cmd/motores/`, `oraculo/compara.py`, `oraculo/fin_comun.py` y `oraculo/fin_run_fixtures.py`.

## Verificación
| Comando | Salida resumida |
|---|---|
| `python3 -I oraculo/fin_genera_fixtures.py` | `30 fixtures y 3 archivos de _datos escritos`, 0 líneas `HALLAZGO`: el oráculo da `invalid_argument` en las 5 consultas nuevas y R-45/R-46 en el control |
| `diff -rq <copia previa> oraculo/fixtures/s2` | solo difiere `s2-ventanas-sin-cubrir.json` |
| `python3 -I oraculo/fin_run_fixtures.py` | `30 fixtures, 30 coinciden, 0 diferencias` |
| Oráculo consulta por consulta (`fin_run_fixtures.evalua`) | `control-p1-p2` OK. `solo-p2`, `hueco-con-pagos` y `hueco-con-costo` dan el mensaje de cobertura. Las 5 nuevas dan «no empieza después de todas las anteriores; van ascendentes y disjuntas (PROVISIONAL, S2-ventanas)» |
| `go run ./cmd/motores s2 dump --fixtures ../oraculo/fixtures/s2 --out ../datos/tmp/motor-s2.json` | exit 0; `cmp` contra `oraculo/salidas/s2.json`: idénticos |
| `python3 -I oraculo/compara.py s2` | oráculo contra esperado 2320/2320, motor contra esperado 2320/2320, **oráculo contra motor 2320/2320**, `model.json (ramas de s2): 159 de 159`, cruces 43/43, provisionales 933/933 (antes 928: +5 consultas de error en la fixture PROVISIONAL), diferencias 0 |
| `gofmt -l internal/motor_fin cmd/motores` · `go vet ./internal/motor_fin/... ./cmd/motores/...` | vacío · limpio |
| `go test -coverprofile=../datos/tmp/s2.cov ./internal/motor_fin/... -count=1` + `go tool cover -func … \| tail -1` | `ok`, `total: (statements) 89.2%` (antes 89.1 %; K25 ≥ 85 %) |
| `go test ./cmd/motores/... -count=1 -cover` | `ok`, 78.4 % |
| `go test ./internal/motor_fin/... ./cmd/motores/... -count=1 -v` | 310 `--- PASS` (antes 309; +`TestVentanasOrden`), 0 FAIL, 0 SKIP; `grep t.Skip` da 0 |
| `python3 -I oraculo/fin_verify.py` | 159/159 (float contra la vía exacta) y 159/159 (float contra `model.json`) |
| `sha256sum oraculo/model.json oraculo/verify3.json` | `1a570666…bfd4` y `d86087f3…5041`, sin cambio; `~/reporte-compra-latex/model/verify3.json` mtime 2026-10-10 18:05:44, sin cambio |

Criterios de S2 cubiertos (con la corrida de arriba):
- **S2-01:** `TestCanon`, 86 subpruebas PASS, 0 SKIP.
- **S2-02:** 159 de 159 en `compara.py` y en `fin_verify.py`.
- **S2-03:** 30 de 30 fixtures; oráculo = motor en 2320/2320 hojas y en los códigos de error de las 9 consultas de `s2-ventanas-sin-cubrir`; `TestOraculoFixtures` + `TestOraculoHojas`, 189 subpruebas PASS. Independencia: `fin_*.py` solo menciona `backend/` en el comentario de `fin_comun.py:4`. La regla de orden está escrita con dos formas distintas.
- **S2-04:** los sha256 de arriba; `grep reporte-compra-latex oraculo/*.py` da 0.
- **S2-05:** `S2-05: 30 fixtures, 69 consultas sin error, 303 invariantes con diferencia 0`; semilla 40 casos.
- **S2-06:** `-run Quincena` PASS.
- **S2-07:** las 7 filas: 292125 = 97375 + 194750 … 629344 = 125870 + 503474.
- **S2-08:** 21988 y 62935 (219.88 y 629.35); el `Decimal` del criterio imprime `219.88 629.35`.
- **S2-09, S2-10, S2-11, S2-12 y S2-14:** `Proyeccion`, `Flujo`, `Presupuesto`, `Meta` y `Zona` PASS.
- **S2-13:** `grep float32|float64` fuera de `_test.go` da 0.
- **S2-15:** gofmt vacío, vet limpio, 89.2 %.
- **S2-16:** la regla nueva lleva la marca PROVISIONAL que cita `S2-ventanas` en `ledger.go` (`periods`), en `fin_devengo.py` (`desordenadas` y el mensaje) y en el campo `provisional` de la fixture; la pregunta está en el bloqueo.

## Supuestos
- Ninguno que cambie cifras. La regla de orden es PROVISIONAL y está en el bloqueo (O1). Ninguna cifra canon cambia, y ninguna consulta que antes daba cifras correctas deja de darlas: solo se rechazan listas de ventanas desordenadas, traslapadas o contenidas.
- Observación de proceso: las dos formas de la regla las escribió el mismo agente, porque la tarea asigna los dos lados a un solo agente. Son formulaciones distintas: anterior inmediata contra máximo acumulado.

## Bloqueos
- `docs/estado/bloqueos/S2-ventanas.md`: sigue abierto. Ahora tiene las preguntas 1 y 2 (cobertura; opciones A/B/C y H1/H2/H3) y la 3 (orden y traslape; O1 aplicada como PROVISIONAL, u O2). Se decide en la puerta 2.

## Pendientes
- Ninguno de M1 dentro de mi alcance.
- Los hallazgos bajos B1 (desfase distinto de «Z»), B2 (código `invalid_date` contra `invalid_argument` en mes o quincena mal formados), B3 (solicitudes de r1 sin aplicar) y B4 (gofmt y vet de `precios` y `recetas`) no estaban en la tarea: quedan para quien los asigne el orquestador.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S2-interfaz.md` §2.12. Cambio: decir que las ventanas van «en orden ascendente y disjuntas» (`start[i] > end[i−1]`; si no, `invalid_argument`) y citar `docs/estado/bloqueos/S2-ventanas.md` (pregunta 3, PROVISIONAL). Motivo: M1 de REVISION-S2-r2; hoy la interfaz no lo dice y las dos implementaciones ya lo exigen.
2. `docs/contratos/motores/S2-interfaz.md` §3.2. Cambio: citar `S2-ventanas` en la frase «incluye periodos anteriores no pedidos» cuando se usan ventanas (con ventanas, lo anterior a la primera es `invalid_argument` por la regla A). Motivo: se arrastra de la solicitud 1 de `correccion-S2-r1.md`; B3 de r2.
3. `docs/contratos/motores/S2-interfaz.md` §5 (tabla de fixtures). Cambio: anotar que `s2-ventanas-sin-cubrir` tiene 9 consultas (4 de cobertura y control, y 5 de orden y traslape). Motivo: trazabilidad de la fixture modificada.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Bloqueo `S2-ventanas`, pregunta 3: se recomienda **O1 (rechazar)**. Es coherente con la opción A ya aplicada a la cobertura: quien llama declara las ventanas contiguas desde el primer evento y en orden. Además mantiene una sola forma de leer «n ventanas consecutivas» y la suma corrida. O2 solo agrega comodidad (aceptar el orden inverso) y aun así tiene que rechazar los traslapes.

## Consumo aproximado
~33 acciones de herramienta; ~150 KB leídos; ~22 KB escritos (código, fixture regenerada, bloqueo y reporte).
