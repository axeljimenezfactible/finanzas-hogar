# revisor-contrato-S4S5-r4 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Revisé de forma independiente `S4-interfaz.md` 1.3, `S5-interfaz.md` 1.3 y `productos-formato.md` 1.1, solo contra los criterios S4 1.2, S5 1.2 y decisiones §D, §E y §K. Informe: `docs/estado/wf1/REVISION-contrato-S4S5-r4.md`.
- Veredicto: **NO APROBADO**, con 0 altos, 2 medios (M-r4-1 y M-r4-2) y 3 bajos nuevos (B-r4-1 a B-r4-3).
- M-r3-1 y M-r3-2 quedan **cerrados**. M-r4-1 es un efecto colateral de la corrección de M-r3-2.
- Recalculé con mis propios scripts el canon de S5-01 (+138.50 y +100.50), S5-02, K40, S4 §8.3, S4-07, RN-15 y RN-16, y la cobertura de las 61 piezas en el catálogo. También implementé `QuoteFor` (§14.1) y `compare_basket` (§7.3) desde el texto.

## Archivos
- `docs/estado/wf1/REVISION-contrato-S4S5-r4.md` (creado)
- `docs/estado/wf1/revisor-contrato-S4S5-r4.md` (este reporte)
- Temporales `datos/tmp/rev-s4s5-r4/verif.py` y `quote.py`: borrados al terminar.

## Verificación
- `python3 -I datos/tmp/rev-s4s5-r4/verif.py`: **29 OK, 0 FAIL**.
  - Canon S5-01: idénticos 277350 / 287400 (+10050); totales 309450 / 323300 (+13850 = 10050 + 5900 − 2100); piezas 61 / 60 / 55; renglones 33 / 3; ganadores 22 / 8 / 3.
  - S5-02: 2933 / 2800, 12692 / 16406, 3375 / 2833.
  - K40: [0.9836, 0.4918, 0.0].
  - S4 §8.3: 4425 / 2213 / 8850.
  - S4-07: MSJ 1748.75 y 1330.25; metas 2099, 1596, 1799 y 1846. RN-17: 2556 y 2056. RN-15: 1486571.
  - Catálogo: 39 presentaciones con los textos exactos de `data.json`; 39 renglones iguales a `cuentas`; 61 piezas y 309450; los precios unitarios cuadran.
- `python3 -I datos/tmp/rev-s4s5-r4/quote.py`:
  - reproduce los 4 resultados de `quote_for` del ejemplo de §14.1 y su `compare_basket` (Walmart 2800, Bodega 1500, `cheapest_store` Bodega);
  - contraejemplo de M-r4-1: totales Walmart 2800 y Bodega 1460, más barata Bodega; el mínimo real de Walmart es 1450.
- `python3 -I -c` con `datetime`: el sábado 10-oct-2026 lleva a `NextVisit` = 2026-10-13; 30 y 90 días son el 2026-11-09 y el 2027-01-08.
- Estructura:
  - bloques JSON válidos (1 en S4, 0 en S5);
  - las tablas de §15 y §18 cubren los 28 criterios;
  - las claves `price.*` y `health.goal_adjust_*` existen en `00001` y `00003`;
  - `grep -c go.mod` = 0 (B-r3-2 sigue abierto).
- Criterios cubiertos: los 15 de S4 y los 13 de S5 tienen sección. Hay dos parciales:
  - S4-03, por M-r4-2;
  - S5-03, por M-r4-1.
  
  La evidencia por criterio está en la tabla del informe.

## Supuestos
- S-1: leo S5-03, «por canasta (lista de compras completa)», como el costo de cubrir la lista en cada tienda con su mejor opción. En esa lectura se apoya M-r4-1, y la corrección sugerida deja abierta la vía del bloqueo si Axel prefiere otra.

## Bloqueos
Ninguno nuevo. Sigue abierto, sin tocar: `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
- Del autor del contrato: corregir M-r4-1 y M-r4-2 (rutas y corrección en el informe).
- Bajos opcionales: B-r4-1 a B-r4-3, B-r3-1, B-r3-2 y los de r1 y r2.

## Solicitudes fuera de alcance
Ninguna.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna de K31.

## Consumo aproximado
- Unas 33 acciones de herramienta.
- Unos 270 KB leídos: protocolo; decisiones §D, §E y §K; criterios S4 y S5; S4, S5 y el catálogo completos por partes; reportes del autor y de la corrección r3; parte de la revisión r3; extracto C1 de `ShoppingList`; 2 bloqueos.
- Unos 26 KB escritos: informe 15.6 KB, este reporte y scripts.
