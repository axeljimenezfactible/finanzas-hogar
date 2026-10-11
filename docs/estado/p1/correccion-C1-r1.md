# Reporte · correccion-C1-r1 (flujo P1, 11-oct-2026)

Agente: `correccion-C1-r1`. Corrige los hallazgos H1 a H5 (2 de severidad alta y 3 media) de `docs/estado/p1/REVISION-C1-r1.md` sobre C1 v1.1. No toca H6 a H8 (baja, opcionales) ni reabre decisiones de §A a §K. Todo se cambió en la fuente del generador (`backend/internal/api/c1gen/`) y se regeneró con `build.py`.

## Hecho

- **H1 (alta, C1-16, K2)**: nueva §7.1 en `C1-datos.md` («Presupuesto: alcance de un cambio de monto»). Trae una tabla con el efecto de `onward` (valor por omisión, «de este mes en adelante») y de `this_month` («solo este mes») sobre los renglones de vigencia (`valid_from_date`, `valid_to_date`), el caso en que el renglón empieza en `month`, `following_month_amount_cents`, 409 `budget_overlap` y 422 `out_of_range`. Incluye un ejemplo de cada alcance con los renglones antes y después y la respuesta 200. Los montos son 1,947.50 y 2,517.37 (194750 y 251737 centavos), tomados del oráculo (`periodos[0]` y `periodos[2]`, `devengado.despensa_total`). El texto aclara que son montos de ejemplo y no presupuestos capturados, y que los `id` son identificadores de ejemplo. Las respuestas omiten `valid_to_date` cuando es NULL (§1 punto 6). También remiten a §7.1 la línea de Cambios (K2) y la fila 5 de §14.
- **H2 (alta, C1-19, K9)**: §1 punto 8 trae ahora el ciclo completo de una captura de gasto o de peso: primer envío 201; reintento con la misma clave y el mismo cuerpo, 200, mismo `id` e `Idempotent-Replayed: true`; otro cuerpo dentro de 24 h, 409 `idempotency_key_reused`. Pasadas las 24 h la clave sigue en el renglón: 200 con el registro guardado, solo se compara la clave, otro cuerpo ya no se detecta y el registro no cambia. Lo acompaña el ejemplo `POST /api/v1/expenses` con `Idempotency-Key`, gasto de contado de $199.00 (19900; `parametros.papel` del oráculo) y `captured_at` `2026-10-10T21:30:00-06:00`, y la nota de que `createWeightEntry` sigue el mismo ciclo. La línea de Cambios (K9) remite a §1 punto 8.
- **H3 (media, C1-17, K3)**: la línea de Cambios (Ticket, K3) y la fila 4 de §14 dicen que el borrador no crea gasto, precio ni inventario y que solo `confirmTicketScan` guarda precios (decisiones E) y vacía el borrador.
- **H4 (media, C1-07, R2-06, K12)**: `col_schema` en `gen1_base.py` (ahí se forma la descripción desde el comentario de columna; el informe señalaba `gen2_schemas.py`) quita «(UTC)» cuando la propiedad es `DateTime`. Las 22 descripciones de propiedad con «(UTC)» bajan a 0. La única aparición que queda en `openapi.yaml` (y en `extractos/componentes-comunes.yaml`) es la del esquema base `DateTime` («El sufijo «Z» (UTC) solo se acepta en las entradas»), que es correcta.
- **H5 (media, C1-16, C1-18)**: se pide a S1 registrar en el diccionario los valores de `scope` y el calculado `chat_id_masked` (ver Solicitudes). El diccionario no es de C1, así que no se tocó.
- Encabezado: el renglón Autor suma «v1.1 corrección r1: agente correccion-C1-r1» y Cambios suma la fila «C1 v1.1 · corrección r1». `info.version` sigue en 1.1.0, porque no cambian rutas, parámetros, tipos ni cifras.

## Archivos

Modificados a mano (fuente del generador):
- `backend/internal/api/c1gen/c1_plantilla.md` (Autor, fila de Cambios, líneas K2, K3 y K9, §1 punto 8, §7.1)
- `backend/internal/api/c1gen/c1_inconsistencias.md` (filas 4 y 5 de §14)
- `backend/internal/api/c1gen/gen1_base.py` (`col_schema`: sin «(UTC)» en instantes)

Regenerados con `build.py`:
- `docs/contratos/C1-datos.md` (124,433 bytes)
- `docs/contratos/openapi.yaml` (629,639 bytes)
- `backend/internal/api/routes.go` (sin cambio de contenido)
- `docs/contratos/extractos/` (17 YAML + `INDICE.md`)

Reporte: `docs/estado/p1/correccion-C1-r1.md`.

## Verificación

| Comando | Salida resumida |
|---|---|
| `python3 -I backend/internal/api/c1gen/build.py` | Escribe `openapi.yaml`, `routes.go`, `C1-datos.md` y los 18 archivos de `extractos/` |
| `python3 -I backend/internal/api/c1gen/build.py --check` | «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| `grep -c '(UTC)' docs/contratos/openapi.yaml` | 1: el esquema base `DateTime`; propiedades con «(UTC)»: 0 (antes 22). Ejemplo: `Expense.captured_at` = «Instante en que la persona capturó el gasto en el teléfono. En una captura atrasada…» |
| goose v3.28.0 `-dir backend/db/migrations sqlite3 datos/tmp/c1r1/c1.db up` | 00001 OK, 00002 OK, versión 2 |
| `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/c1r1/c1.db` | 23 OK, 0 FALLA, «RESULTADO: TODO OK». C1-15: 6 marcas, 3 bloqueos abiertos citados, 0 sin bloqueo, 0 citados que no existen. EXT: 201 = 201. S6-03: 201 = 201 |
| `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator` sobre `openapi.yaml` y los 17 YAML de `extractos/` | `openapi.yaml: OK`; 17 de 17 extractos OK |
| `python3 -I docs/design/tools/verificar_u1.py` | «RESULTADO: TODO OK» (U1-03 sin faltantes) |
| `cd backend && go build ./... && go vet ./... && gofmt -l .` (go1.27.2) | build OK, vet OK, gofmt vacío |
| Cifras nuevas de `C1-datos.md` | 199.00 (`parametros.papel`), 1,947.50 y 2,517.37 (`devengado.despensa_total` de P1 y P3) están en `oraculo/model.json`. Los centavos 19900, 194750 y 251737 son exactos. No se agregó ninguna cifra inventada |

Temporales en `datos/tmp/c1r1`, ya borrados.

## Supuestos

1. En el ejemplo de `this_month`, el renglón original conserva su `id` para la parte anterior al mes y los renglones nuevos llevan `id` nuevos. La descripción de `updateBudget` ya lo dice así para `onward` («el renglón termina… y otro… rige»), y el ejemplo lo marca como identificadores de ejemplo. No cambia cifras ni esquema.
2. «Pasadas las 24 h… el registro guardado no cambia» se lee de lo que ya dice `openapi.yaml` («responde 200 con el mismo gasto… sin crear otro»). No es regla nueva.

## Bloqueos

Ninguno nuevo.

## Pendientes

- H6, H7 y H8 (severidad baja, opcionales según el informe) no se corrigieron. Siguen en `docs/estado/p1/REVISION-C1-r1.md`, para que el orquestador decida si los reparte.

## Solicitudes fuera de alcance

1. **Ruta**: `docs/contratos/diccionario-datos.md` (dueño S1). **Cambio**: registrar los valores del parámetro `scope` de `updateBudget`: `this_month` = «solo este mes» y `onward` = «de este mes en adelante» (por omisión, preseleccionado; K2). Registrar también el campo calculado `chat_id_masked` de `telegram_chat_candidate`: «•••» y los últimos 4 caracteres del `chat_id`, por ejemplo «•••4821»; el `chat_id` completo nunca se devuelve (K4). **Motivo**: hallazgo H5 de `REVISION-C1-r1.md`. C1-16 pide que el diccionario registre los dos alcances y C1-18 dice que los campos los fija el diccionario. Hoy `grep` da 0 para los tres.
2. **Ruta**: `docs/estado/p1/C1-v1.1.md` §Solicitudes (reporte del agente C1-v1.1). **Cambio**: anotar que la solicitud 1 de arriba completa la solicitud 2 de ese reporte, que hoy solo pide `ocr_draft_json` e `idempotency_key`. **Motivo**: H5 señala que el reporte C1-v1.1 no la pedía. Este agente no escribe ese archivo.

## Paquetes instalados

Ninguno.

## Consumo aproximado

- Acciones de herramienta: 26.
- Leído: unos 105 KB (protocolo 6, §K 8, informe de revisión 17, plantilla e inconsistencias por partes unos 25, generador por partes unos 15, oráculo y semilla unos 10, salidas de verificación unos 20). No se leyó entero ningún archivo grande.
- Escrito a mano: unos 18 KB (parche de unos 10 KB con texto nuevo de unos 8.5 KB, más este reporte de unos 7 KB). Los archivos regenerados los escribe `build.py`.
