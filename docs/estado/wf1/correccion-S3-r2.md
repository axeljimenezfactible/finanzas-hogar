# correccion-S3-r2 · WF-1 · 11-oct-2026

Corrige los dos hallazgos medios de `docs/estado/wf1/REVISION-S3-r2.md` (M1 y M2). No hay hallazgos altos. Los bajos (B1 a B4) no estaban en mi encargo y no los toqué.

## Hecho
- **M1 (S3-02 / S3-12) · salida guardada del oráculo.** Regeneré `oraculo/salidas/s3.json` con `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json`. Ahora tiene 27 fixtures y 79 consultas, incluidas las 6 `s3-comparador-*`. `compara.py s3` con su oráculo por defecto da 0 diferencias y código 0 (antes, 2,036 diferencias y código 1).
- **M2 (S3-12) · el oráculo medía reloj.** Apliqué la opción (a) que recomienda la revisión. `oraculo/cred_escenarios.py` ya no compara contra el reloj: con `timeout_seconds` = 0 responde `timeout`, que es determinista (fixture `s3-comparador-tiempo`), y con un valor positivo calcula completo. Quité las tres revisiones de `time.monotonic()` de `compare()` (antes en las líneas 157, 167 y 181) y dejé la explicación en el docstring y en un comentario. El tiempo real contra los 5 s (K43, configurable por K44) lo mide solo el motor en `TestRestriccionesTiempo`. El motor no cambia.
- La opción (a) no contradice el texto vigente del contrato: §8.7 pide `context.WithTimeout` a la API, a la CLI y a las pruebas del motor, y nunca le pide al oráculo medir reloj. Aun así la revisión pide que el integrador del contrato lo deje escrito, así que va como solicitud (abajo).
- `docs/estado/S3-escenarios.md`: la observación de §7 quedó como resuelta y agregué la §8 con el comportamiento al agotarse el tiempo (motor y oráculo), la prueba con el equipo cargado y la regeneración de la salida.

## Archivos
- Modificado: `/home/deepsec/finanzas-hogar/oraculo/cred_escenarios.py` (solo `compare()` y el docstring).
- Regenerado: `/home/deepsec/finanzas-hogar/oraculo/salidas/s3.json` (lo escribe el script, ningún número se editó a mano).
- Modificado: `/home/deepsec/finanzas-hogar/docs/estado/S3-escenarios.md` (§7, última viñeta, y §8 nueva).
- Nuevo: `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-S3-r2.md` (este reporte).

## Verificación
- `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json`: «27 fixtures · 79 consultas» y «segunda vía (float, 2 decimales): 174 de 174 cifras coinciden». sha256 `b84c0db7fa80d19dffbd9f9d1834de52b5bef34f9c8f561727e0b81f2b1bbd8f`, igual en una segunda corrida a `datos/tmp`.
- `go build -o datos/tmp/motores ./cmd/motores` y `motores s3 dump --fixtures oraculo/fixtures/s3 --out datos/tmp/motor-s3.json`: código 0.
- **M1:** `python3 -I oraculo/compara.py s3 --motor datos/tmp/motor-s3.json`, con el oráculo por defecto: oráculo contra esperado, motor contra esperado y oráculo contra motor **906 de 906** cada una; model.json (ramas de s3) **437 de 437**; cruces 0 de 0; data.json 9 de 9; provisionales **2,718 de 2,718**; **diferencias 0; código 0**. También sin argumentos (`compara.py s3`): diferencias 0, código 0.
- **M2, con el equipo cargado:** dos procesos ocupando el núcleo 5 y el oráculo en ese mismo núcleo (`taskset -c 5`). `cred_escenarios.py`: la carga tardó **7,395 ms** (más de 5 s) y aun así dio `ok`, 3,870 evaluados, 0 aceptados con violación; código 0. `cred_run_fixtures.py` en esas condiciones: el mismo sha256 `b84c0db7…` que sin carga. `compara.py s3` con esa salida: diferencias 0, código 0. Antes, en esas condiciones, salía `timeout (5003 ms)` y 2 diferencias.
- `python3 -I oraculo/cred_escenarios.py` sin carga: autoprueba bien; caso-1 y caso-3 con 31 evaluados, 17 válidos y 14 descartados; caso-2 con 26, 12 y 14; 0 aceptados con violación y 4 de 4 óptimos distintos en los tres; `sin-restricciones` → `missing_constraints`; `tiempo` → `timeout`; carga `ok` con 3,870 evaluados en 4,449 ms (dato informativo); «escenarios aceptados con violación = 0»; código 0.
- `python3 -I backend/internal/motor_cred/comparador/testdata/tabla_mano.py`: «OK … = tabla a mano» en las 5 fixtures; código 0.
- `gofmt -l internal/motor_cred cmd/motores`: vacío. `go vet ./internal/motor_cred/... ./cmd/motores/`: limpio.
- `go test -count=1 -coverprofile=datos/tmp/s3.cov ./internal/motor_cred/...`: `motor_cred` 93.7 %, `canontest` 90.0 %, `comparador` 96.0 %, `evaluar` 90.8 %; `go tool cover -func … | tail -1`: **total 94.0 %**. `go test ./cmd/motores/`: ok.
- `-v`: 133 `--- PASS`, 0 `FAIL` y 0 `SKIP`; `TestCanon85` 85 de 85 subpruebas; «27 fixtures corridas» (`evaluar`) y «6 fixtures del comparador corridas» (`comparador`). `S3-TIEMPO`: caso-1 0 ms, caso-2 0 ms, caso-3 0 ms y carga 127 ms con 3,870 candidatos.
- Los `grep` de S3-15 y S3-17: `highs|ortools` en go.mod y go.sum 0; `net/http` fuera de pruebas 0; `float32|float64` fuera de pruebas 0. `reporte-compra-latex` en `oraculo/cred_*.py`: 0. Los `cred_*.py` no leen la salida del motor (solo escriben `oraculo/salidas/s3.json`).
- `sha256sum oraculo/model.json`: `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, igual antes y después.
- Criterios cubiertos por esta corrección: **S3-02** (437 de 437 y 906/906 con la salida guardada por defecto); **S3-12** (la vía del oráculo ya no depende del reloj; 0 válidos con violación; motor = oráculo = tabla; 5 s sembrado; 0 s → `timeout` / 503 `unavailable`; peor caso del motor 127 ms, sin bloqueo de tiempo); **S3-16** (sha256 intacto, oráculo independiente); **S3-18** (vet, gofmt y 94.0 %). Los demás criterios siguen como los dejó REVISION-S3-r2 (nada de su código cambió).
- Temporales: `datos/tmp/motores`, `motor-s3.json`, `s3-oraculo-2.json`, `s3-oraculo-cargado.json`, `s3.cov` y `s3-v.txt`, todos borrados. No toqué `datos/tmp/rev-s5r2/` (es de otro agente).

## Supuestos
- Regeneré `oraculo/salidas/s3.json` aunque la carpeta `oraculo/salidas/` no aparece en la lista literal de mi alcance, porque M1 me lo encarga con ese comando. Es una salida del script, sin números editados a mano.
- Opción (a) de M2: el oráculo es un modelo de resultados, no de rendimiento. No cambia cifras, esquema, endpoints ni tokens, y no contradice §8.7 vigente; falta que el integrador del contrato lo deje escrito (solicitud 1).

## Bloqueos
- Ninguno nuevo. Siguen los 5 `docs/estado/bloqueos/S3-*.md` de antes (sin cambios). No hace falta `S3-tiempo-optimizador.md`: el motor tarda 127 ms en la carga, el 2.5 % de los 5 s.

## Pendientes
- B1 (baja, fuera de mi encargo): `cred_escenarios.py` no compara el `status` y el `evaluated_count` de cada fixture con su `esperado`. Con la opción (a) la carga ya no puede salir en `timeout`, pero el código de salida del script sigue sin revisar el `esperado`. Archivo: `/home/deepsec/finanzas-hogar/oraculo/cred_escenarios.py` (`_resumen`).
- B2, B3 y B4 (bajas) siguen abiertos como dice REVISION-S3-r2.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S3-interfaz.md` §8.7 y §10.0 (integrador del contrato). Cambio: escribir que el oráculo no mide reloj: `timeout_seconds` = 0 → `timeout` (determinista) y un valor positivo calcula completo; el tiempo real contra `comparator.timeout_seconds` lo mide solo `TestRestriccionesTiempo` en el motor. Motivo: REVISION-S3-r2 M2. El texto actual no dice qué debe hacer el oráculo, y medir reloj hacía que la verificación oficial dependiera de la carga del equipo. `cred_escenarios.py` ya funciona así.

## Paquetes instalados
- Ninguno (`taskset` ya estaba instalado).

## Propuestas para Axel
- Ninguna nueva. El tiempo máximo del comparador sigue en 5 s (K43, FIRME, configurable por K44); el peor caso del motor tarda 127 ms.

## Consumo aproximado
- Acciones de herramienta: 24.
- Leído: unos 110 KB (protocolo, §K, criterios S3, REVISION-S3-r2, extractos de §8 y §10 del contrato, partes de `cred_escenarios.py` y `compara.py`, salidas de comandos).
- Escrito: unos 12 KB a mano (`cred_escenarios.py`, `S3-escenarios.md` y este reporte). Aparte, el script regeneró `oraculo/salidas/s3.json` (149,709 bytes; salida generada, no escrita a mano).
