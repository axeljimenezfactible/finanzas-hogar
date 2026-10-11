# revisor-U1-r2 · Reporte final (flujo P1, 11-oct-2026)

Revisor independiente de U1 1.1, ronda 2 de 3. No corrigió nada. El veredicto y los hallazgos están en `docs/estado/p1/REVISION-U1-r2.md`.

## Hecho
- Leí las lecturas obligatorias, todas por partes: el protocolo completo, decisiones §K, los criterios U1 1.1, R2-05 y R3-01 de `REVISION-r3.md`, los reportes `U1-1.1.md` y `correccion-U1-r1.md`, los hallazgos de `REVISION-U1-r1.md` y `verificar_u1.py`.
- Volví a correr `verificar_u1.py` e hice mis propias verificaciones con scripts propios: recuento de pantallas, N-xxx y matriz §3.29; toques de F-1; «vigente»; «Para C1» contra `openapi.yaml`; `period` en las operaciones de la lista.
- Comprobé con grep, línea por línea, las decisiones K2–K5, K8–K10, K12, K13, K16, K21 y K23, y los hallazgos R2-05 y R3-01.
- Revisé el cierre de M1, M2 y M3: los tres están **cerrados**.
- Veredicto: **NO APROBADO**, con 1 alta, 0 medias y 5 bajas.
  - **A1 (alta):** U1-09 pide la capacidad de K18 (parámetros de la tarjeta bancaria y del crédito de tienda, con «falta capturar») y U1 no la tiene. El esquema 1.2 ya agrega esas columnas en `credit`.
  - **Bajas:** B1 a B4 vienen de la ronda 1 y siguen abiertas; B5 es nueva (el supuesto 2 y las líneas de `U1-1.1.md` quedaron viejos).

## Archivos
- `docs/estado/p1/REVISION-U1-r2.md` (nuevo: el informe de revisión)
- `docs/estado/p1/revisor-U1-r2.md` (este reporte; reemplaza el de progreso)
- Temporales `datos/tmp/revisor-u1-r2/` (`ctx.py`, `recuento.py`, `extra.py`): creados y borrados

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py` → «RESULTADO: TODO OK», salida 0. Coincide con lo reportado:
  - 27 pantallas, 14 secciones, 203 N-xxx (138 L y 65 E) y Local N-003, N-301 y N-317;
  - F-1 = 3 y 8 flujos;
  - PROVISIONAL en 65, 419, 1149 y 1150;
  - 197 operaciones (196 con `x-necesidades`);
  - `PENDIENTE C1 v1.1: ['N-316', 'N-318', 'N-319']`.
- `python3 -I datos/tmp/revisor-u1-r2/recuento.py`:
  - P-xx = 27 y P-26.xx = 14, así que T = 41;
  - 203 N-xxx (138 L y 65 E), igual a lo declarado en §3.29;
  - en la matriz, 0 faltantes, 0 sobrantes y 0 en la columna equivocada;
  - sin operación y no Local: N-316, N-318 y N-319; ninguna Local tiene operación;
  - única operación sin N-xxx: `getOpenApiDocument`;
  - `period` es opcional, sin valor por omisión, en las 4 operaciones de la lista de compras de C1 1.0.0.
- `python3 -I datos/tmp/revisor-u1-r2/extra.py`:
  - F-1 suma 1 → 1 → 2 → 3, sin incrementos incoherentes;
  - «lista de compras vigente» en negritas solo en la línea 65, y toda mención de la lista remite a §0.2;
  - «Para C1 · Nuevas» = N-316, N-317, N-318 y N-319, y las 45 N-xxx citadas existen.
- Grep de K18 en U1: 0 coincidencias. El único «Falta capturar» está en la línea 732 y es de K16. En `00002_puerta1.sql:86-92` están las columnas de K18 en `credit`.

## Supuestos
1. Revisé U1-09 completo, incluida K18, porque está en `docs/criterios/U1.md`, aunque K18 no esté en la lista de decisiones del prompt. No cambia cifras ni esquema; solo el alcance de la revisión.
2. No conté como falla de U1-03 que N-316, N-318 y N-319 aún no tengan operación: el prompt (punto h) las da por pendientes de C1 v1.1 y están en «Para C1».

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/U1-lista-vigente-domingo.md`, que no es de este agente.

## Pendientes
- Ronda 3: corregir A1 en `docs/design/U1-arquitectura.md`: §2.3 (P-26.05), N-283, N-289, P-18 (N-205), la lista de cambios y «Para C1». Agregar una comprobación de K18 en `docs/design/tools/verificar_u1.py`. Opcional: B1 a B5.

## Solicitudes fuera de alcance
1. **Orquestador:** incluir K18 en la lista de decisiones de la corrección de la ronda 3 de U1. Motivo: U1-09 la exige y el prompt de U1 1.1 no la incluía.
2. **`docs/estado/p1/U1-1.1.md:72`:** marcar el supuesto 2 como sustituido por `correccion-U1-r1.md` (es la solicitud 1 de ese reporte; B5).
3. **C1 v1.1:** cuando se aplique «Para C1», incluir los campos de K18 de `credit` (`rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents`, `annual_rate_decimals`) en cuanto U1 les dé necesidad.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: 53 de 60.
- Leído: unos 187 KB, todo por partes. U1 por rangos y grep (unos 60 KB); `verificar_u1.py` completo (18 KB); criterios U1 (21 KB); el resto en reportes y salidas de scripts. `openapi.yaml` solo lo leyeron los scripts.
- Escrito: unos 33 KB (informe 16 KB, este reporte 5 KB, progreso 1 KB y scripts temporales 8 KB, ya borrados).
