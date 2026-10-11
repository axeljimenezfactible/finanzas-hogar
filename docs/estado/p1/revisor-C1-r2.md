# Reporte · revisor-C1-r2 (flujo P1, 11-oct-2026)

## Hecho
- Leí el protocolo, §K de `decisiones.md`, `docs/criterios/C1.md` v1.1, el reporte `C1-v1.1.md` y `correccion-C1-r1.md`. Del contrato solo leí por partes: `C1-datos.md` (Cambios, §1, §4, §6, §7.1, §9, §11, §13 y §14) y `openapi.yaml` (grep y scripts).
- Volví a correr las verificaciones (a) a (j) con mis propios comandos y scripts en `datos/tmp/revC1r2/`, y comprobé los hallazgos H1 a H5 de r1.
- Escribí el informe `docs/estado/p1/REVISION-C1-r2.md`. Veredicto: **NO APROBADO** (1 alto, 2 medios y 3 bajos).

## Archivos
- `docs/estado/p1/REVISION-C1-r2.md` (nuevo)
- `docs/estado/p1/revisor-C1-r2.md` (este reporte)

## Verificación
- `build.py --check`: IDÉNTICO en los 21 archivos. El sha256 de `extractos/` sale igual en dos corridas.
- goose 00001+00002 y después `verify_c1.py --db`: 23 OK, 0 FALLA, TODO OK.
- openapi-spec-validator 0.9.0: `openapi.yaml` OK y 17 de 17 extractos OK.
- `verificar_u1.py`: TODO OK, 0 necesidades sin operación ni Local.
- `go build`, `go vet` y `gofmt -l` (go1.27.2): OK, OK, vacío.
- Extractos: 201 = 201 operaciones, 0 repetidas, el mayor de 76,202 B, todos parsean.
- Lista K, marcas, versión y cifras: detalle en el informe.
- Encontré un defecto que `verify_c1` no detecta: `period` tiene dos definiciones (`Period` y `ListWeek`).

## Supuestos
- Las cifras que vienen de K12 (25.8685) y del diccionario o de RN-17 (factores de actividad) cuentan como no inventadas.
- Leer los extractos con un script cuenta como lo que imprimió el script, no como su tamaño.

## Bloqueos
- Ninguno nuevo. Para H-r2-02 sugiero un bloqueo `C1-presupuesto-onward` si K2 no alcanza para decidir.

## Pendientes
- Corregir H-r2-01 (alta) y H-r2-02 (media) en `backend/internal/api/c1gen/` y regenerar.
- H-r2-03: aplicar la solicitud 1 de `correccion-C1-r1.md` en `docs/contratos/diccionario-datos.md`.
- Ronda 3 de revisión.

## Solicitudes fuera de alcance
1. `docs/contratos/diccionario-datos.md` (S1): registrar `this_month`/`onward` y `chat_id_masked` (repite la solicitud 1 de `correccion-C1-r1.md`). Además, alinear la línea 530 (`budget.valid_to_date`) con la regla de `onward` que se fije en H-r2-02.
2. `backend/internal/api/c1gen/verify_c1.py`: que C1-06 compruebe un solo componente por nombre de filtro. Es parte del alcance del corrector de C1.

## Paquetes instalados
- Ninguno. Usé la caché de uv para pyyaml, jsonschema y openapi-spec-validator 0.9.0.

## Consumo aproximado
- 41 acciones de herramienta.
- Leído: unos 135 KB (ningún archivo grande entero).
- Escrito: unos 33 KB (informe 13 KB, este reporte 3 KB y scripts temporales unos 17 KB, ya borrados).
- `datos/tmp/revC1r2/` borrado. `cerrado_por_tope`: false.
