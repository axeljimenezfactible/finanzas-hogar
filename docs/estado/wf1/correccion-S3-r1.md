# correccion-S3-r1 · Corrección de la revisión S3 ronda 1 (11-oct-2026)

Corrige los hallazgos H1 y H2 (altos) y H3 (medio) de `docs/estado/wf1/REVISION-S3-r1.md`. No toqué los bajos (H4–H7), ningún cálculo ni ninguna decisión de §A–§K. No cambié ningún número de la tabla a mano.

## Hecho
- **H1:** escribí las 9 fixtures del comparador en `oraculo/fixtures/s3/` con un script temporal (`python3 -I`, ya borrado). Las 5 de `comparador/testdata/` se copiaron sin cambiar un número: `entradas.datos` de cada caso pasó a `_datos/comparador-caso-N.json` con `{"$ref": …}`, y `-sin-restricciones` y `-tiempo` apuntan a `_datos/comparador-caso-1.json` (el script comprobó que sus `datos` eran idénticos a los del caso 1). `s3-comparador-carga` lleva en línea la rejilla de `carga()`, volcada con `json.Marshal` desde una prueba temporal que ya borré, sin `today` ni `credits` en `request` (§9.2). Su `esperado` es `status` «ok» y `evaluated_count` 3,870 (cuenta a mano 6 × (15 × 6 × 7 + 15)), y lleva la marca PROVISIONAL `S3-objetivos-medida` igual que los casos 1 a 3 (§10.0, fila de `-carga`).
- **H2:** `evaluar.go` atiende «comparar» con `comparador.Query`, que es nueva y equivale a `EvalComparar` pero devuelve el error con su tipo. `EvalComparar` ahora la envuelve. No hay ciclo de importación. `motores s3 dump` y `compara.py s3` ya incluyen el comparador.
- **H3:** las pruebas del comparador leen `oraculo/fixtures/s3/s3-comparador-*` (resuelven `$ref` y fallan con «falta oraculo/fixtures/s3/<archivo> (autor: optimizador S3)»). Agregué `TestOraculoFixtures` en `comparador`, que corre las 6 fixtures. Borré las copias `testdata/{caso-1,caso-2,caso-3,sin-restricciones,tiempo}.json` y la función `carga()`: `TestRestriccionesTiempo` lee la carga de su fixture. `testdata/tabla_mano.py` ya no escribe copias: recalcula la tabla y comprueba que las fixtures son idénticas (sale con 1 si no lo son). Queda una sola fuente.
- `evaluar_test.go`: el caso «comparar `{}` → unsupported» pasa a `invalid`. Agregué `unknown_field` (campo de más) e `invalid` (`timeout_seconds` −1).
- `docs/estado/S3-escenarios.md`: la secuencia (§0) y el párrafo de la carga (§5) dicen ahora dónde vive la tabla. Agregué la §7, «Integración con el oráculo», con la corrida.

## Archivos
- Nuevos: `oraculo/fixtures/s3/_datos/comparador-caso-1.json`, `_datos/comparador-caso-2.json`, `_datos/comparador-caso-3.json`, `oraculo/fixtures/s3/s3-comparador-caso-1.json`, `s3-comparador-caso-2.json`, `s3-comparador-caso-3.json`, `s3-comparador-sin-restricciones.json`, `s3-comparador-tiempo.json`, `s3-comparador-carga.json`.
- Modificados: `backend/internal/motor_cred/evaluar/evaluar.go`, `backend/internal/motor_cred/evaluar/evaluar_test.go`, `backend/internal/motor_cred/comparador/evaluate.go`, `backend/internal/motor_cred/comparador/comparador_test.go`, `backend/internal/motor_cred/comparador/testdata/tabla_mano.py`, `docs/estado/S3-escenarios.md`, este reporte.
- Borrados: `backend/internal/motor_cred/comparador/testdata/{caso-1,caso-2,caso-3,sin-restricciones,tiempo}.json`.
- Sin cambios: `backend/cmd/motores/s3.go` (ya llamaba `evaluar.Evaluate`), `oraculo/cred_*.py` y `oraculo/model.json`.

## Verificación
- `gofmt -l internal/motor_cred cmd/motores`: vacío. `go vet ./internal/motor_cred/... ./cmd/motores/`: limpio.
- `go test -coverprofile=datos/tmp/s3.cov ./internal/motor_cred/... -count=1`: `motor_cred` ok 93.7 %, `canontest` ok 90.0 %, `comparador` ok 96.0 %, `evaluar` ok 90.8 %; `go tool cover -func … | tail -1` → **total 94.0 %**. `t.Skip` en `motor_cred`: 0.
- `-v`: `evaluar` `TestOraculoFixtures` PASS («27 fixtures corridas»); `comparador` `TestOraculoFixtures` PASS («6 fixtures del comparador corridas»); `TestCanon85` 85/85 PASS; `TestRestricciones` PASS (caso-1 31/17/14, caso-2 26/12/14, caso-3 31/17/14; 0 válidos con violación); `TestRestriccionesPropiedad` PASS (300 comparaciones, 11,730 candidatos, 0 violaciones); `TestRestriccionesFaltan`, `TestRestriccionesVivas` y `TestPuntaje` PASS; `TestOptimos` PASS (4 óptimos distintos por caso, iguales a §6 de S3-escenarios); `TestRestriccionesTiempo` PASS: `S3-TIEMPO` caso-1 1 ms, caso-2 1 ms, caso-3 1 ms, carga 137 ms (3,870 candidatos).
- `python3 -I backend/internal/motor_cred/comparador/testdata/tabla_mano.py`: «OK oraculo/fixtures/s3/<caso> = tabla a mano» ×5; código 0.
- `python3 -I oraculo/cred_escenarios.py`: autoprueba bien; caso-1, caso-2 y caso-3 ok con 0 aceptados con violación y 4 de 4 óptimos distintos; sin-restricciones → `missing_constraints`; tiempo → error `timeout`; carga ok con 3,870 evaluados en 4,549 ms; «escenarios aceptados con violación = 0»; **código 0** (antes 1).
- `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out datos/tmp/…/oraculo-s3.json`: 27 fixtures, 79 consultas; segunda vía 174 de 174.
- `motores s3 dump --fixtures oraculo/fixtures/s3 --out datos/tmp/…/motor-s3.json` (binario compilado en `datos/tmp`): código 0.
- `python3 -I oraculo/compara.py s3 --oraculo … --motor …`: oráculo contra esperado, motor contra esperado y oráculo contra motor **906 de 906** cada una; model.json (ramas de s3) **437 de 437**; data.json 9 de 9; provisionales 2,718 de 2,718; **diferencias 0**; código 0. Las fixtures provisionales (casos 1 a 3 y carga) entran en la línea «provisionales». Por eso comparé además las salidas completas de oráculo y motor en los 6 casos del comparador: son idénticas (1,864 hojas).
- `grep net/http` fuera de pruebas: 0. `grep float32|float64` fuera de pruebas: 0. `grep highs|ortools` en go.mod y go.sum: 0. `reporte-compra-latex` en `oraculo/cred_*.py`: 0. sha256 de `model.json` igual antes y después (`1a57066689ddeed3…`).
- Criterios que cierra esta corrección: **S3-12** (las dos verificaciones literales pasan; conjunto motor/oráculo = tabla, 0 diferencias en los 3 casos; «faltan restricciones»; 503 con 0 s; 5 s sembrado y medido), **S3-13** (12 óptimos = esperado = oráculo por `TestOraculoFixtures` y `compara.py`) y **S3-18** (`go test ./internal/motor_cred/...` ok, vet y gofmt limpios, 94.0 %). Los demás siguen pasando (canon 85/85; 437/437).
- Borré los temporales de `datos/tmp/correccion-s3/` y `datos/tmp/s3.cov`. Lo demás que hay en `datos/tmp` es de otros agentes.

## Supuestos
- La `provisional` de `s3-comparador-carga` copia la del caso 1 (`S3-objetivos-medida` y, como otro bloqueo, `S3-serie-sin-intereses`), porque la rejilla tiene las mismas piezas: líneas tipo Kueski y una promoción 0 % a 4 q. No cambia cifras, porque su `esperado` solo fija `status` y `evaluated_count`.
- Las salidas de verificación las escribí en `datos/tmp/` y no en `oraculo/salidas/s3.json`, que no está en mi alcance.

## Bloqueos
- Ninguno nuevo. Siguen los de S3 citados en las listas `provisional`: `S3-objetivos-medida`, `S3-serie-sin-intereses`, `S3-formula-tarjeta`, `S3-origen-cargo-normal` y `S3-convencion-prestamo`.

## Pendientes
- Oráculo: `oraculo/cred_escenarios.py` tarda unos 4.5 s en `s3-comparador-carga` y aplica el mismo límite de 5 s de reloj (§8.7, líneas 157–182). En un equipo cargado podría responder `timeout` en esa fixture y `compara.py` marcaría diferencia, aunque el motor tarda 137 ms. No lo cambié porque no estaba entre los hallazgos y tocaría la semántica de §8.7 en el oráculo. Opciones para una ronda futura: acelerar el oráculo o, si el contrato lo permite, que no mida reloj en la fixture de carga.
- Hallazgos bajos H4–H7 de la revisión: fuera del encargo de esta corrección.

## Solicitudes fuera de alcance
1. `oraculo/salidas/s3.json` (integrador): regenerarlo con `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json`, para que incluya las 6 fixtures del comparador (27 fixtures, 79 consultas).
2. `docs/contratos/motores/S3-interfaz.md` §11 (integrador del contrato): registrar que `TestOraculoFixtures` existe en `evaluar` (núcleo) y en `comparador`, y que las pruebas del comparador leen `oraculo/fixtures/s3/` (ya no `testdata/`). Lo mismo que la solicitud de H7.
3. `docs/estado/wf1/motor-S3-comparador-go.md` (reporte de otro agente): sus solicitudes 1 y 2 quedan atendidas por esta corrección. No lo edité.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva. El tiempo máximo del comparador sigue en 5 s (K43, FIRME, configurable por K44). Con el motor, la carga tarda 137 ms, el 2.7 % del límite.

## Consumo aproximado
- Acciones de herramienta: ~40.
- Leído: ~125 KB (protocolo, revisión, decisiones §K, criterios S3 §1–§5, contrato §9–§11 por partes, código y salidas).
- Escrito con herramientas: ~25 KB (scripts temporales, ediciones de Go y Python, S3-escenarios §7 y este reporte).
- Generado por script (copias de fixtures existentes, fuera de contexto): ~133 KB.
