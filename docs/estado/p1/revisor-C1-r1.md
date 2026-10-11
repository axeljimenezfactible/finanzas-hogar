# revisor-C1-r1 · Reporte final (flujo P1, 11-oct-2026)

## Hecho
- Revisión independiente de C1 v1.1 (ronda 1 de 3) contra `docs/criterios/C1.md` 1.1, decisiones §K (K1–K7, K9, K11–K13, K16, K18, K21, K23) y los hallazgos R2-06 y R3-01. No corregí nada.
- Reejecuté las verificaciones (a) a (j) con comandos y scripts propios en `datos/tmp/rev-c1`.
- Veredicto: **NO APROBADO**, con 2 hallazgos altos (H1 C1-16, H2 C1-19), 3 medios (H3 a H5) y 3 bajos (H6 a H8). El detalle está en `docs/estado/p1/REVISION-C1-r1.md`.

## Archivos
- `docs/estado/p1/REVISION-C1-r1.md` (informe de revisión, nuevo).
- `docs/estado/p1/revisor-C1-r1.md` (este reporte).
- `build.py --extractos` corrido dos veces para la prueba de regeneración de C1-21: reescribió `docs/contratos/extractos/` con bytes idénticos (mismo sha256 antes y después). Ningún otro archivo del repositorio cambió.

## Verificación
- (a) `python3 -I backend/internal/api/c1gen/build.py --check`: IDÉNTICO en los 21 archivos (`openapi.yaml`, `routes.go`, `C1-datos.md`, `extractos/`).
- (b) goose v3.28.0 `up` sobre `datos/tmp/rev-c1/c1.db` (00001 OK, 00002 OK, versión 2), luego `verify_c1.py --db`: TODO OK, 0 FALLA.
- (c) `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator`: `openapi.yaml` OK y 17 de 17 extractos OK.
- (d) `python3 -I docs/design/tools/verificar_u1.py`: TODO OK, 0 N-xxx sin operación ni Local.
- (e) `go build ./...` OK, `go vet ./...` OK, `gofmt -l .` vacío (go1.27.2).
- (f) Extractos: 201 operaciones = `grep -c operationId` 201, 0 repetidas, 0 faltantes. Mayor 76 220 bytes. `componentes-comunes.yaml` presente. 17 de 17 YAML parsean. sha256 de la carpeta igual en dos regeneraciones (`c6b365bc…eb243`).
- (g) Lista K: `scope` (enum de 2, `default: onward`), `saveTicketScanDraft`, los dos endpoints de candidatos, `Idempotency-Key` `required` 16–128 con `captured_at` en `createExpense`/`createWeightEntry`, `rate_decimals`, `fortnight_payment_cap_cents`, `OptimizerSettings`/`missing_constraints`, campos K18 de `Credit`, `ProjectionSettings`, `goal`, `family` de 3 valores y `period` obligatorio en la lista. Todos presentes en `openapi.yaml`, y sus columnas en la base migrada y en el diccionario.
- (h) 6 marcas en `openapi.yaml` y 13 líneas de marca en `C1-datos.md`. Citan solo `C1-valores-k31`, `criterios-promedios-redondeo` y `S1-ajuste-objetivo-kcal`, los tres sin «Respuesta de Axel».
- (i) `info.version` 1.1.0. Fila «C1 v1.1 (puerta 1)» en Cambios.
- (j) Ejemplos: 14 nodos y 20 números. El único sin origen es `Error.request_id` (identificador marcador). En `C1-datos.md` todas las cifras con decimales salen del oráculo o de la semilla.
- Temporales: borré `datos/tmp/rev-c1`. No toqué `datos/tmp/verify-wf0`, que no es mío.

## Supuestos
- H3 (C1-17) va como media y no como alta: el requisito de fondo («el borrador no crea precios») sí está en el contrato (`openapi.yaml`). Solo falta en `C1-datos.md`, que es donde lo busca la verificación.
- H1 y H2 van como altas porque el cuerpo de C1-16 y de C1-19, no solo su verificación, exige el texto y los ejemplos en `C1-datos.md`.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ronda 2: corregir H1 y H2 y luego H3 a H5. Fuentes: `backend/internal/api/c1gen/c1_plantilla.md` y `backend/internal/api/c1gen/gen2_schemas.py`. Además, la solicitud a S1 para `docs/contratos/diccionario-datos.md` (H5). Después regenerar con `build.py` y reejecutar (a) a (j).

## Solicitudes fuera de alcance
1. `docs/contratos/diccionario-datos.md` (S1): registrar los valores de `scope` de `updateBudget` (`this_month` = «solo este mes», `onward` = «de este mes en adelante», por omisión) y el campo calculado `chat_id_masked`. Motivo: C1-16 y C1-18 (H5).
2. S1 (informativa): `credit.vat_on_interest_bp` NULL = «ya incluye IVA», así que no se puede avisar «falta capturar» para el IVA de K18. Si se quiere ese aviso, hace falta otra columna. Motivo: tensión con la nota 16 del diccionario.

## Paquetes instalados
- Ninguno. `uv`/`uvx` usaron su caché (pyyaml, jsonschema, openapi-spec-validator 0.9.0).

## Consumo aproximado
- Acciones de herramienta: 45.
- Leído: unos 155 KB (por partes, grep y sed; ningún archivo grande entero).
- Escrito: unos 35 KB (scripts temporales unos 10 KB, informe 17 KB, reporte y progreso 6 KB).
