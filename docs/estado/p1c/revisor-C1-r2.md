# Reporte · revisor-C1-r2 (flujo P1c, ronda 2 de 3)

## Hecho
- Revisé C1 v1.2 después de `correccion-C1-r1`. Lo evalué contra `docs/criterios/C1.md` 1.2, contra K36, K38, K43 y K44 y contra el hallazgo R3-01. No corregí nada.
- Corrí por mi cuenta las verificaciones (a) a (k), R3-01 y algunas extra de C1-04, C1-08, C1-12 y C1-13, con tres scripts míos en `datos/tmp/rc2/`, que ya borré.
- Veredicto: **APROBADO**, sin altos ni medios. A-01 está cerrado. Quedan 8 bajos: B-01 a B-05 de la ronda 1, de los que B-05 está parcial, y los nuevos N-01 a N-03.

## Archivos
- `docs/estado/p1c/REVISION-C1-r2.md` (informe: veredicto, verificaciones, tabla de criterios, hallazgos y estado de los anteriores)
- `docs/estado/p1c/revisor-C1-r2.md` (este reporte)

## Verificación
- `python3 -I backend/internal/api/c1gen/build.py --check`: IDÉNTICO en los 21 archivos.
- goose v3.28.0 `up` sobre `datos/tmp/rc2/c1.db` (versión 3), luego `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/rc2/c1.db`: TODO OK, incluidas las dos líneas de C1-23.
- `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator`: OK en openapi.yaml y en los 17 extractos. 3.1.1, 167 rutas, 1.2.0.
- `python3 -I docs/design/tools/verificar_u1.py`: TODO OK, 0 PENDIENTE.
- `go build ./...`, `go vet ./...` y `gofmt -l .`: OK, OK y sin salida.
- Scripts propios:
  - Extractos: 210 de 210, 0 duplicadas, ninguno pasa de 80 000 B, INDICE cuadra.
  - K44: 15 claves con `value`, `initial_value`, cambio y restablecer, sin `maximum` y sin `x-provisional`; las 9 de 00003 están en el contrato.
  - K36: los campos y `getHealthSettings`/`updateHealthSettings` existen.
  - Marcas: 0 citan `U1-lista-vigente-domingo` o `C1-valores-k31`; los 3 `x-provisional` citan `S1-parametros-avisos`, que sigue abierto.
  - C1-23: probé el patrón y el mínimo con jsonschema.
  - Ejemplos: 0 cifras inventadas.
  - R3-01: 8 de 8 operaciones con `week` obligatoria.
  - Extra: 953 de 953 errores van a `Error`; 0 `number`.

## Supuestos
- La fila de Cambios de v1.1 r4 (C1-datos.md L24) es historial, no una marca. La fila v1.2 (L25) declara K38 FIRME.
- C1-22 dice al pie de la letra que los umbrales van en `params_json` con PROVISIONAL. Lo evalué con K44, que manda: K44 nombra los umbrales como claves de `app_setting`, y el bloqueo `S1-parametros-avisos` ya se redujo. La ronda 1 hizo lo mismo; queda como B-01, que corresponde al agente de criterios.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (de S1).

## Pendientes
- Ninguno en mi alcance.

## Solicitudes fuera de alcance
1. `docs/criterios/C1.md`: actualizar C1-22 para que la excepción de avisos cubra solo `send_time`, `advance_days`, `weekday` y `day_of_month`, y los umbrales y [1, 0] queden como `alerts.*` de `app_setting` por K44. Actualizar también C1-13 (v1.2 y 1.2.0), C1-21 (210 operaciones) y C1-12 (excepción `x-app-setting`). Motivo: B-01; el texto quedó desfasado respecto de K44 y del esquema 1.3.
2. C1, en la siguiente corrección: N-01 (`gen4_build.py` L111 y `gen5_go.py` L15, versión desde una constante), N-02 (`minimum: 0` en `goal_adjust_*_kcal` de Create y Update), N-03 (citar K18 en vez del bloqueo resuelto), B-02, B-03 y B-04. Son bajos y no bloquean.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones de herramienta: 39.
- Lectura: unos 185 KB. Incluye el protocolo, §K, C1.md 1.2 y, por partes con grep o sed, la ronda 1, `gen2_schemas.py`, el diff de los generadores, C1-datos.md (Cambios, §10, §11.7 y §13), los reportes de C1 y la salida de las verificaciones.
- Escritura: unos 26 KB (scripts temporales, ya borrados, y los dos informes).
- `cerrado_por_tope`: false.
