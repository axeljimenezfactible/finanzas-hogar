# correccion-C1-r3 · Reporte (flujo P1d «Ronda 4 de criterios y bajos de C1»)

Agente correccion-C1-r3 (Opus), 11-oct-2026. HEAD `82b9fd1`. Lista de trabajo: los hallazgos bajos de `docs/estado/p1c/REVISION-C1-r2.md` (N-01, N-02, N-03, B-02, B-03, B-04, B-05) más I-01 y el comentario de `send_time` en la 00003. No hay hallazgos altos ni medios abiertos en C1. B-01 (texto de `docs/criterios/C1.md`) no es de este agente.

## Hecho
- **N-01.** `gen4_build.py` define una sola constante `C1_VERSION = '1.2.0'`, que da `info.version`. El encabezado de `openapi.yaml` (L1) y el de `routes.go` (`gen5_go.py`, L1) leen `info.version`, igual que ya lo hacían los extractos (`gen7_extractos.py`). La primera línea de los dos archivos dice ahora «C1 v1.2.0» (antes «v1.1» y «v1.0»).
- **N-02.** `gen2_schemas.py`: `minimum: 0` en `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` de `HealthProfileCreate` y `HealthProfileUpdate` (K36: magnitud ≥ 0; la 00003 tiene `CHECK (… >= 0)`). En Update la propiedad sigue aceptando `null` (`type: [integer, null]`). Los dos esquemas viven en el extracto `ajustes-d.yaml`.
- **N-03.** La descripción de `CardStatus` (N-205) cita RN-26 y el criterio S3-08 (K18) en vez del archivo de bloqueo `criterios-credito-formulas.md`. Ya no queda ninguna mención de ese archivo en `openapi.yaml` ni en `C1-datos.md`.
- **B-02.** `NotificationUpdate` ya no lleva `x-provisional` a nivel de esquema, así que `is_enabled` no queda marcado. La marca de lugar la llevan `params_json` (por `weekday` y `day_of_month`) y sus claves `send_time` y `advance_days`. La descripción de esas dos claves cita `S1-parametros-avisos`. La descripción del esquema ya no dice PROVISIONAL: dice que `is_enabled` es FIRME y dónde va la marca. `gen6_c1.py` exige que `Notification` y `NotificationUpdate` traigan las mismas propiedades en `params_json`, por eso las claves también van marcadas en `Notification` (ver Supuestos).
- **B-03.** Descripción nueva de `params_json` en `Notification` y `NotificationUpdate`, con puntuación correcta. Dice que desde 1.3 los umbrales (80 y 100 %) y los días antes del pago (1 y 0) son claves `alerts.*` de `app_setting` (RN-29, `getAlertSettings`) y que el Down de la 00003 los repone. Dice también que en `params_json` quedan cuatro claves: `send_time` y `advance_days` (K43, sin valor sembrado), y `weekday` (7) y `day_of_month` (1), con valor FIRME de la decisión F. El lugar de las cuatro es PROVISIONAL (`S1-parametros-avisos`). El texto ya no afirma lo que luego niega. No se tocó la descripción de la columna en las migraciones: el contrato la reemplaza.
- **B-04.** `verify_c1.py`: el mensaje de la regla del domingo dice «es FIRME por K38 y no deja marcas pendientes». Ya no dice PROVISIONAL ni nombra `U1-lista-vigente-domingo`. La línea informativa dice «FIRME por K38; bloqueo de origen abierto: False». La lógica de la comprobación no cambió.
- **B-05.** `docs/estado/p1c/C1-v1.2.md` §Verificación trae el comando completo de `verify_c1.py` (`uv run --no-project --with pyyaml --with jsonschema …`), goose con `-dir` y el comando de openapi-spec-validator tal como lo usa `make verify-openapi` (`uvx --offline --from openapi-spec-validator==0.9.0 …`), con los extractos. También tiene una viñeta nueva con las verificaciones de esta ronda.
- **I-01.** Cabecera de `C1-datos.md` (desde `c1_plantilla.md`). «Estado» dice puerta 1b aplicada (K35 a K44, C1 v1.2) con los bajos corregidos en la ronda 3. «Fuentes» dice esquema 1.3 (00001+00002+00003, con sus tres archivos) y U1 1.2. Quité el conteo viejo de «199 necesidades». «Autor» suma este agente.
- **Comentario de la 00003.** En `backend/db/migrations/00003_puerta1b.sql`, L56 (`@param_sin_valor` de `send_time`), «lo propone S7 (WF-1)» pasó a «lo propone S7 (WF-3, K31, puerta 4)». Es la única línea que cambió en el archivo. Regeneré `docs/contratos/diccionario-datos.md` con `gen_dict.py`, y en él cambió solo esa línea. El texto llega también a la tabla de parámetros de `C1-datos.md` por regeneración.
- Fila «C1 v1.2, ronda 3 (P1d)» en la tabla Cambios de `c1_plantilla.md`, que sale en `C1-datos.md`, con estos puntos. También agregué la subsección «C1 v1.2 · ronda 3 (P1d)» en `backend/internal/api/c1gen/README.md`.
- No cambian rutas, operaciones, esquemas (siguen 252), cifras ni `info.version` (1.2.0).

## Archivos
- Fuente: `backend/internal/api/c1gen/gen2_schemas.py`, `gen4_build.py`, `gen5_go.py`, `verify_c1.py`, `c1_plantilla.md` y `README.md`.
- Generados con `build.py`: `docs/contratos/openapi.yaml`, `backend/internal/api/routes.go`, `docs/contratos/C1-datos.md`, `docs/contratos/extractos/INDICE.md`, `ajustes-d.yaml` y `creditos-a.yaml` (los demás extractos regenerados quedaron idénticos).
- `backend/db/migrations/00003_puerta1b.sql`: solo la línea 56, que es un comentario. `docs/contratos/diccionario-datos.md`: una línea, regenerada con `gen_dict.py`.
- `docs/estado/p1c/C1-v1.2.md` (§Verificación) y este reporte.

## Verificación
| Comando | Salida resumida |
|---|---|
| `python3 -I backend/internal/api/c1gen/build.py` | 55 tablas y 15 parámetros (K44) de 3 migraciones; openapi con 167 rutas, 210 operaciones y 252 esquemas; routes.go con 210 rutas |
| `python3 -I backend/internal/api/c1gen/build.py --check` | IDÉNTICO en los 21 archivos: openapi.yaml de 675 925 B, routes.go, C1-datos.md de 144 878 B, INDICE.md, 16 extractos y componentes-comunes.yaml. «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/c1r3/c1.db up` | 00001, 00002 y 00003 OK; versión 3 |
| `uv run --offline --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/c1r3/c1.db` | 29 OK y **RESULTADO: TODO OK**. C1-15: 6 marcas, todas de `S1-parametros-avisos`, ninguna sin bloqueo abierto. Regla del domingo: «FIRME por K38», 0 marcas pendientes. C1-23: nada mal en openapi ni en `ajustes-d.yaml`, y el patrón acepta y rechaza bien |
| `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml docs/contratos/extractos/*.yaml` | 18 de 18 OK |
| `python3 -I docs/design/tools/verificar_u1.py` | RESULTADO: TODO OK, 0 PENDIENTE. openapi 1.2.0, 210 operaciones, 209 con x-necesidades |
| `cd backend && go build ./... && go vet ./... && gofmt -l .` | BUILD_OK, VET_OK, gofmt vacío |
| `head -1` de openapi.yaml y de routes.go | «# Finanzas Hogar · Contrato C1 v1.2.0 (OpenAPI 3.1)…» y «// Code generated … (contrato C1 v1.2.0). DO NOT EDIT.» |
| `grep -c '^    operationId:'` (4 espacios, como lo pide el prompt) | 0, igual que en HEAD: en openapi.yaml `operationId` va con 6 espacios. `grep -c '^      operationId:'` da 210, igual al Total de `extractos/INDICE.md` (210) |
| tamaño de los extractos | ninguno pasa de 80 000 B; el mayor es gastos.yaml con 76 220 B |
| `python3 -I backend/db/tools/gen_dict.py` y luego `git diff --stat docs/contratos/diccionario-datos.md` | versión 1.3, 15 parámetros, 8 sin valor, 0 diferencias con lo sembrado. Diff: 1 línea (la de `send_time`) |
| `make verify-s1 GOOSE=…/gobin/goose` | RESUMEN TODO OK (DDL, Down de la 00002 y de la 00003, negativas y goose up / down-to 0 / up) |
| `make verify-all GOOSE=…/gobin/goose` | «Verificación ampliada OK (WF-0, C1, S1, U1 y criterios)»: C1 con 21 IDÉNTICO y TODO OK, U1 TODO OK, criterios con TOTAL 261 |
| `grep -c x-provisional docs/contratos/openapi.yaml` | 7: las 6 marcas más la mención en `info.description`. Antes eran 4 |
| `git diff --stat` | Mis archivos: los 6 de c1gen, la 00003 (1 línea), los generados (openapi.yaml, routes.go, C1-datos.md, INDICE.md, ajustes-d.yaml, creditos-a.yaml, diccionario-datos.md) y C1-v1.2.md. Además aparecen cambios en `docs/criterios/` (C1.md, INDICE.md, P3-claude-design.md, S4.md, S6.md, S7.md, U3.md, U4.md): no son míos, salieron durante esta sesión y vienen del agente paralelo de criterios del mismo flujo. No los toqué |

## Supuestos
- B-02: además de `NotificationUpdate`, también marqué `send_time` y `advance_days` en `Notification`, que es el esquema de lectura. `gen6_c1.py` exige que los dos esquemas traigan las mismas propiedades en `params_json`, y el lugar de esas claves está en duda igual en la lectura. Mantuve la marca en `params_json` porque `weekday` y `day_of_month` solo existen ahí, por `additionalProperties`, y C1-23 la exige. Ninguna cifra ni endpoint cambia.
- En la cabecera quité «199 necesidades» de la celda Fuentes en vez de poner un número nuevo, para no inventar un conteo.

## Bloqueos
- Sin bloqueos nuevos. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` y es el único que citan las marcas.

## Pendientes
- Ninguno de mi alcance.

## Solicitudes fuera de alcance
- `backend/db/migrations/00003_puerta1b.sql` L57 (`@param_sin_valor` de `advance_days`) y `docs/estado/bloqueos/S1-parametros-avisos.md` (Contexto) siguen diciendo «lo propone S7 (WF-1)». El prompt me pidió cambiar solo la línea de `send_time`. Si `advance_days` también va en WF-3, el orquestador decide si se corrige (S1 y el dueño del bloqueo), regenerando el diccionario y C1.
- B-01 de REVISION-C1-r2 (`docs/criterios/C1.md`: C1-22, C1-13 «C1 v1.1» y C1-21 «hoy 201») le toca al agente de criterios.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Unas 34 acciones de herramienta. Leído: unos 100 KB (protocolo, §K, bloqueo, REVISION-C1-r2, partes del README y del Makefile, fragmentos de gen2/gen4/gen5/gen6/verify_c1 y diffs). Escrito a mano: unos 25 KB (parches, reporte y viñetas). Los generados salen por regeneración.
