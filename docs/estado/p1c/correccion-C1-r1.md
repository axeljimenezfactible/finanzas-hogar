# Corrección C1 v1.2 · ronda 1 (flujo P1c, puerta 1b)

Agente correccion-C1-r1, 11-oct-2026. Corrige el único hallazgo alto o medio de `docs/estado/p1c/REVISION-C1-r1.md`: A-01 (C1-23 y C1-22, parámetros de avisos; K43). Los bajos (B-01 a B-05) no estaban en el encargo y no se tocaron. No se reabrió ninguna decisión de §A–§K.

## Hecho
- **A-01 corregido.** `params_json` de `Notification` y de `NotificationUpdate` declara en `properties` las dos claves de K43:
  - `send_time`: `type: string`, `pattern: ^([01][0-9]|2[0-3]):[0-5][0-9]$`. La descripción cita K43: hora local HH:MM de 24 horas, America/Mexico_City.
  - `advance_days`: `type: integer`, `minimum: 0`. La descripción cita K43: días de anticipación del aviso de precios viejos.
  - Ninguna de las dos tiene `default`. No tienen valor sembrado: el valor lo propone S7 (RN-29 1.3) y es configurable (K44).
  - Se conservan `additionalProperties: true`, para `weekday` y `day_of_month` de la decisión F, y la marca de lugar `x-provisional`, que cita el bloqueo abierto `S1-parametros-avisos`. El tipo y el rango son FIRMES (K43). Lo único PROVISIONAL es dónde viven.
- El generador se detiene en dos casos: si `params_json` de cualquiera de los dos esquemas pierde `additionalProperties: true`, la marca o la cita del bloqueo, o si `Notification` y `NotificationUpdate` no traen las mismas dos propiedades.
- C1-datos.md (generado): al final de §11.7 hay un párrafo con la validación de `send_time` y `advance_days`, tomada de las propiedades de openapi. Lleva el patrón literal y `minimum: 0`. También hay una fila «C1 v1.2 · corrección r1» en Cambios y el autor en la cabecera.
- `verify_c1.py` tiene la comprobación nueva **C1-23**, que revisa:
  - el patrón exacto y `minimum: 0` en openapi.yaml, en el extracto que trae `NotificationUpdate` (`ajustes-d.yaml`) y en C1-datos.md;
  - que no haya `default` y que sigan `additionalProperties: true` y la marca de lugar;
  - que el patrón acepte de 00:00 a 23:59 y rechace 24:00, 7:30, 12:60, 12:5, 1230 y «08:00 »;
  - el resto del criterio: 10 MB con 413 en `uploadExpenseAttachment`, 503 `unavailable` en `compareScenarios`, N-023, N-064, N-226, N-297 y N-300 como FIRME (K43) en §10, M-05 con 21 y 0 líneas `PROVISIONAL.*C1-valores-k31`.
- Regenerado con `build.py`. Cambian solo 4 salidas: openapi.yaml, C1-datos.md, extractos/ajustes-d.yaml y extractos/INDICE.md. Para comprobarlo se revirtieron temporalmente las fuentes y se corrió `build.py --check`, que dio «4 archivo(s) difieren». Después se restauraron las fuentes. routes.go y los otros 16 extractos no cambian. `info.version` sigue en 1.2.0 y siguen las 210 operaciones.

## Archivos
- `backend/internal/api/c1gen/gen2_schemas.py`: bloque `SEND_TIME_PATTERN` y `PARAMS_K43` después de `NotificationUpdate`.
- `backend/internal/api/c1gen/gen6_c1.py`: `PARAMS_K43` para §11.7.
- `backend/internal/api/c1gen/c1_plantilla.md`: marcador `{{PARAMS_K43}}`, fila de Cambios y autor.
- `backend/internal/api/c1gen/verify_c1.py`: comprobación C1-23 y descripción.
- `backend/internal/api/c1gen/README.md`: sección «C1 v1.2 · corrección r1».
- Regenerados: `docs/contratos/openapi.yaml`, `docs/contratos/C1-datos.md`, `docs/contratos/extractos/ajustes-d.yaml` (61 871 B, menos de 80 000) y `docs/contratos/extractos/INDICE.md`.
- `docs/estado/p1c/correccion-C1-r1.md` (este reporte).

## Verificación
| Comando | Salida resumida |
|---|---|
| `python3 -I backend/internal/api/c1gen/build.py --check` | «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| Fuentes revertidas temporalmente, luego `build.py --check` | DIFIERE solo openapi.yaml, C1-datos.md, extractos/INDICE.md y extractos/ajustes-d.yaml (4). Con las fuentes restauradas: los 21 coinciden |
| goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/c1r1/c1.db up` | 00001, 00002 y 00003 OK, versión 3 |
| `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/c1r1/c1.db` | 29 OK (antes 27, más las 2 de C1-23), «RESULTADO: TODO OK». Línea C1-23: mal en openapi `{}`, extracto `['ajustes-d.yaml']` mal `{}`, patrón acepta y rechaza bien `True`, C1-datos.md sin la clave `[]`, 10 MB/413/503 `True`, filas §10 `[]`, M-05 `True`, grep `PROVISIONAL.*C1-valores-k31` = 0 |
| El mismo `verify_c1.py` sin `--db` (migraciones en memoria) | «RESULTADO: TODO OK» |
| Prueba en negativo: copia temporal de openapi.yaml con el patrón cambiado a `^[0-9]{2}:[0-9]{2}$` y `verify_c1.py --root <copia>` | «FALLA C1-23 send_time con patrón HH:MM…», «RESULTADO: 1 falla(s)» |
| `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator openapi.yaml extractos/*.yaml` | 18 de 18 OK |
| `python3 -I docs/design/tools/verificar_u1.py` | «RESULTADO: TODO OK» |
| `cd backend && go build ./... && go vet ./... && gofmt -l .` | BUILD_OK, VET_OK, gofmt vacío |
| `grep -c '\[01\]\[0-9\]'` | openapi.yaml 2, ajustes-d.yaml 2, C1-datos.md 1 (antes 0) |
| `grep -c x-provisional openapi.yaml` | 4 líneas: las mismas 3 marcas de antes más la mención en `info.description`. No se agregó ninguna marca |
| `sha256sum docs/contratos/extractos/* \| sha256sum` | `12a95b591cb1198d…` (cambia por ajustes-d.yaml e INDICE.md) |

## Supuestos
- Ninguno que cambie cifras, esquema ni endpoints. Las descripciones de `send_time` y `advance_days` repiten el texto de K43 y del bloqueo `S1-parametros-avisos` («resumen diario y avisos con hora», «aviso de precios viejos antes del día de despensa»). No agregan valores.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`, que decide dónde viven `send_time`, `advance_days`, `weekday` y `day_of_month`. La marca de lugar se conservó.

## Pendientes
- Ninguno del encargo. No se corrigieron los bajos B-01 a B-05 de la revisión, porque el encargo era solo de alta y media. En particular, B-02 (la marca `x-provisional` a nivel de esquema en `NotificationUpdate`) y B-04 (la etiqueta de la regla del domingo en `verify_c1.py`) son de este alcance y quedan para quien lo decida.

## Solicitudes fuera de alcance
- Ninguna nueva. B-01, el texto de C1-22, C1-13, C1-21 y C1-12 en `docs/criterios/C1.md`, sigue siendo del agente de criterios, como dice la revisión.

## Propuestas para Axel
- Ninguna. Los valores de `send_time` y `advance_days` los propone S7 (K31 y RN-29). Este agente no propone valores.

## Consumo aproximado
- Acciones de herramienta: unas 29.
- Leído: unos 115 KB: protocolo 6, revisión 13, §K 15, bloqueos y criterios 12, fragmentos de generadores, plantilla, verificador y README 35, fragmentos de extractos y C1-datos 10, y salidas de comandos y diffs 25.
- Escrito a mano: unos 14 KB (parches a 5 fuentes y este reporte). El generador reescribió, sin que pasaran por el contexto, openapi.yaml, C1-datos.md y los extractos.
- Temporales en `datos/tmp/c1r1/` (base SQLite y copia de las fuentes): borrados al terminar.
