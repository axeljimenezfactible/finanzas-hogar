# revisor-U1-r1 · Reporte de agente (flujo P1, 11-oct-2026)

## Hecho
- Revisé U1 1.1 (`docs/design/U1-arquitectura.md`, `docs/design/tools/verificar_u1.py`, el bloqueo `U1-lista-vigente-domingo` y el reporte `docs/estado/p1/U1-1.1.md`) contra `docs/criterios/U1.md` 1.1, contra §K (K2–K5, K8–K10, K12, K13, K16, K21, K23) y contra R2-05 y R3-01.
- **Veredicto: NO APROBADO.** Hay 0 hallazgos altos, 3 medios (M1, M2, M3) y 4 bajos (B1 a B4).
  - M1: N-090 contradice a K21, N-318 y N-091 sobre los meses 7 a 12.
  - M2: F-7 no aplica K16.
  - M3: la regla PROVISIONAL del domingo se repite sin marca en las líneas 353 y 876.
- El detalle está en `docs/estado/p1/REVISION-U1-r1.md`.
- No corregí nada del artefacto revisado.

## Archivos
- `docs/estado/p1/REVISION-U1-r1.md` (nuevo): el informe de revisión.
- `docs/estado/p1/revisor-U1-r1.md` (nuevo): este reporte.
- Temporales en `datos/tmp/rev-u1/`: scripts `rev.py` y `rev2.py`, y las líneas añadidas y quitadas del diff. Ya están borrados.

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py` → código de salida 0, «RESULTADO: TODO OK».
  - Da 27 pantallas, 14 secciones, 203 N-xxx (138 L y 65 E) y Local = N-003, N-301 y N-317.
  - F-1 = 3 toques; hay 8 flujos; marcas PROVISIONAL en 65, 419, 1147 y 1148.
  - `openapi.yaml` 1.0.0: 197 operaciones, 196 con `x-necesidades`.
  - Línea `PENDIENTE C1 v1.1: ['N-316', 'N-318', 'N-319']`.
  - Todo coincide con el reporte.
- `git diff --stat HEAD -- docs/design/U1-arquitectura.md` → 115 líneas añadidas y 81 quitadas.
- Script propio `rev.py` (PyYAML y diff de HEAD):
  - N-xxx: 203 = 138 L + 65 E, la última es N-319, sin duplicados. HEAD tenía 199 y la última era N-315. No se quitó ninguna ni cambió su tipo o pantalla.
  - Pantallas: 27 P-xx + 14 P-26.xx = T 41. La matriz tiene 41 filas, ninguna N-xxx falta en ella y ninguna tiene un tipo distinto al de su columna.
  - Única diferencia: N-242 se define en P-22 y la matriz la pone en P-21. Ya estaba así en v1.0.
  - `openapi.yaml`: las N-xxx sin operación son N-316, N-318 y N-319; la única operación sin N-xxx es `GET /openapi.yaml`.
  - `period` es opcional en las 4 operaciones de la lista; `createExpense` y `createWeightEntry` aceptan `Idempotency-Key`.
  - Las 43 N-xxx cuya fila cambió están todas en «Para C1».
- Script propio `rev2.py`:
  - Encontró por grep las líneas de cada K y de R2-05 y R3-01.
  - Marcas PROVISIONAL: 65 y 1148 citan un bloqueo abierto; 419 y 1147 citan K31.
  - Las apariciones de «vigente» relacionadas con la lista remiten a §0.2.
  - El único tamaño táctil en px es 44 × 44.
  - «Cola» solo aparece con N-317, N-054 y N-240.
  - §2.4 cubre las 27 pantallas sin traslapes.
  - La regla «domingos, la siguiente» aparece fuera de §0.2 en las líneas 353 y 876.
- Grep de F-7 (líneas 947–959): dice «restricciones, tope y pesos por defecto (N-225)» y no menciona «Faltan restricciones».

## Supuestos
1. N-316, N-318 y N-319 sin operación en C1 1.0.0 no cuentan como hallazgo de U1: dependen de C1 v1.1, que está fuera del alcance de U1, y la verificación (h) del prompt los trata como lista esperada.
2. Las marcas PROVISIONAL de las líneas 419 y 1147, que citan un bloqueo resuelto más K31, son válidas, porque K31 ordena dejar el valor PROVISIONAL con la propuesta escrita hasta que Axel lo apruebe.
3. K21 abarca todos los meses 7 a 12 de P-07, porque la propia U1 lo afirma en N-318 y en el supuesto 2 del reporte.

Ninguno de estos supuestos cambia cifras ni esquema.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ronda 2: el agente U1 corrige M1, M2 y M3, y el revisor de la ronda 2 comprueba lo listado en §8 del informe.

## Solicitudes fuera de alcance
- `docs/contratos/` (C1 v1.1): cerrar R3-01 del lado de C1 (`period` obligatorio en las escrituras de la lista, o con el mismo valor por omisión que U1 §0.2) y considerar B3, para no pedir `period` en `updateShoppingListItem`. Motivo: U1-16 depende de C1.
- `docs/estado/bloqueos/INDICE.md`: registrar `U1-lista-vigente-domingo` como abierto. Motivo: ya lo pidió U1-1.1 y sigue sin registrarse.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: 29 de 60.
- Leído: unos 150 KB. Nunca abrí entero ni U1 ni `openapi.yaml`; `openapi.yaml`, U1 y el HEAD de U1 solo los leyeron los scripts.
- Escrito: unos 34 KB (scripts temporales, unos 11.5 KB; informe, 18 KB; este reporte, unos 4.5 KB).
