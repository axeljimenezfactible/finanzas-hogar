# revisor-C1-r3 · Reporte (flujo P1d «Ronda 4 de criterios y bajos de C1»)

Revisor independiente (Opus), 11-oct-2026. HEAD `82b9fd1`. Revisé el trabajo de `correccion-C1-r3` contra C1.md 1.2 (HEAD), K36, K38, K43 y K44, y contra los bajos N-01, N-02, N-03, B-02, B-03, B-04 y B-05 de REVISION-C1-r2. También revisé I-01 y la línea de `send_time` de la 00003. Informe completo: `docs/estado/p1c/REVISION-C1-r3.md`.

## Hecho
- Leí el protocolo, §K (K36, K38, K43 y K44), el bloqueo `S1-parametros-avisos`, REVISION-C1-r2, el reporte `correccion-C1-r3.md` y los diffs de `c1gen/`, de la 00003 y de `C1-v1.2.md`. De `openapi.yaml`, `C1-datos.md` y el diccionario solo leí lo que sale con grep y `sed -n`.
- Reejecuté las verificaciones (a) a (k) con mis propios comandos. Las 11 pasan.
- Comparé el openapi del árbol con el de HEAD campo por campo (script yaml): no cambia ninguna ruta, ninguna operación ni el conjunto de esquemas; los 17 campos que cambian están todos en la lista.
- Cerré con evidencia, con número de línea, los 8 hallazgos de la lista y la línea de `send_time`.
- Veredicto: **APROBADO**. 0 altos, 0 medios y 3 bajos (R3-B01, R3-B02 y R3-B03).
- No corregí nada y no hice commit.

## Archivos
- Creados: `docs/estado/p1c/REVISION-C1-r3.md` y este reporte.
- Temporales: `datos/tmp/revc1r3/` (base de goose, copia de openapi y de C1.md de HEAD, diccionario regenerado, script y logs). Ya los borré.

## Verificación
- (a) `python3 -I backend/internal/api/c1gen/build.py --check`: IDÉNTICO en los 21 archivos.
- (b) goose v3.28.0 `up` sobre `datos/tmp/revc1r3/c1.db` (versión 3) y después `uv run --offline --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db …`: TODO OK.
- (c) `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator` sobre openapi.yaml y los 17 extractos: 18 de 18 OK.
- (d) `python3 -I docs/design/tools/verificar_u1.py`: TODO OK, 0 PENDIENTE.
- (e) `go build ./...`, `go vet ./...` y `gofmt -l .` con Go 1.27.2: OK, OK y vacío.
- (f) Ningún extracto pasa de 80 000 B (el mayor mide 76 220 B). La suma de operaciones de los extractos (210) es igual a la de openapi (210), con el mismo conjunto.
- (g) Evidencia por hallazgo: tabla del informe.
- (h) Hay 6 `x-provisional`, todos de `S1-parametros-avisos` (abierto), y ninguno a nivel de esquema.
- (i) Contra HEAD: mismas 210 operationId, 167 rutas y 252 esquemas, y 0 operaciones cambiadas.
- (j) En la 00003 cambia 1 línea, la L56, que es un comentario. `gen_dict.py` hacia un temporal da un diccionario idéntico al del árbol. `make verify-s1 GOOSE=…`: rc 0 y RESUMEN TODO OK.
- (k) `make verify-all GOOSE=…`: rc 0 y «Verificación ampliada OK (WF-0, C1, S1, U1 y criterios)». El árbol quedó igual (mismo md5 de `git status --porcelain` antes y después).

## Supuestos
- «C1.md versión 1.2» es la de HEAD (la copié con `git show HEAD:`), porque `correccion-criterios-r4` la está editando en paralelo y el prompt excluye esos cambios.
- Tomé como válidas las marcas de lugar en `send_time` y `advance_days` porque citan el bloqueo abierto, que cubre justo esas dos claves (protocolo §2: lo no decidido va PROVISIONAL). Las anoto como redundantes en R3-B02 (baja).

## Bloqueos
- Ninguno nuevo. El único abierto sigue siendo `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
- Ninguno de mi alcance.

## Solicitudes fuera de alcance
- `backend/db/migrations/00003_puerta1b.sql` L57 y `docs/estado/bloqueos/S1-parametros-avisos.md` L2: cambiar «lo propone S7 (WF-1)» por «S7 (WF-3, K31, puerta 4)», como ya está en L56, y regenerar `diccionario-datos.md` y C1. Motivo: R3-B03; hoy hay dos líneas vecinas que se contradicen. Dueño: S1 y el orquestador.
- `backend/internal/api/c1gen/gen4_build.py` L101–L102, `c1_plantilla.md` L5 y `gen2_schemas.py` L364: R3-B01 y R3-B02, opcionales, para la próxima pasada de C1.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Unas 33 acciones de herramienta. Leído: unos 150 KB (protocolo, §K, bloqueo, REVISION-C1-r2, reporte de corrección, diffs, salidas de verificación y greps). Escrito: unos 19 KB (informe y reporte).
