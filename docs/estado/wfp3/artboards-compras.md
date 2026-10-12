# artboards-compras · Pista 3, fase 2 (lienzo) · 12-oct-2026

## Hecho
8 artboards en `docs/design/maquetas/lienzo/project/compras/` (móvil 412 px y escritorio fluido de P-13, P-14, P-15 y P-16) y el fragmento `docs/design/maquetas/lienzo/fragmentos/compras.json` (8 tableros, `order` en el orden P-13 m, e, P-14 m, e, ...; escritorio con `expand: fill`; `is_interactive: false`).

| Archivo | Contenido |
|---|---|
| compras-p13-lista-m / -e | Lista agrupada por tienda o mercado (5 grupos, 24 artículos, total $1,564.10), casillas de comprado, tienda asignada contra recomendada (5 artículos difieren), motivo, frescura del precio. Escritorio: vistas Tabla/Tablero, resumen, fila «Estados». |
| compras-p14-canasta-m / -e | Selector de canasta, ganadora (Bodega Aurrera $1,551.70), tarjetas por tienda con idénticos/sustitutos/sin precio/viejos, cobertura y diferencia contra la más barata, ganadores, membresía de Sam's con veredicto. Escritorio: pestaña Detalle apilada (24 artículos x 5 tiendas) y «Estados». |
| compras-p15-productos-m / -e | Catálogo de 33 productos (móvil: 12 y «Mostrar 21 más»), precio por kg/L/pza, frescura, fuente, resumen (32 frescos, 1 viejo). Escritorio: vista previa lateral de Jitomate saladet con precios por tienda e historial (SVG, solo los puntos capturados, patrón de trazo por tienda). |
| compras-p16-inventario-m / -e | 14 artículos ordenados por estado (1 agotado, 5 reponer, 8 suficiente), existencia, ritmo, días restantes, última entrada. Escritorio: vista previa con movimientos de Leche entera y formulario de ajuste con motivo. |

Enlaces: barra lateral y botones principales apuntan a los artboards de la lista de pantallas (mismo ancho); «Capturar precio», «Nuevo producto», «Ajustar existencia», «Registrar sobrante» van a «#» (son hojas, sin artboard).

## Verificación
Script propio `python3 -I` (en `/tmp/compras_gen/verify.py`), resultado `RESULTADO: OK`: líneas obligatorias; JSON de data-props parsea y coincide con el fragmento; etiquetas balanceadas con html.parser; 0 hex y 0 rgb/hsl fuera de data-props; 0 cifras de P3-12 (con y sin coma, y en centavos); todos los `var(--x)` existen en `sistema/project/tokens.css`; tabular-nums; sin iframe, innerHTML ni red; sin emoji; sin signo menos ASCII antes de `$`; un solo h1; hrefs a artboards de la lista; fragmento lista exactamente los 8 archivos. Además, todos los huecos `{{ }}` son rutas con punto y existen en los datos de `renderVals()`.

## Supuestos
- Los artboards se generaron con scripts en `/tmp/compras_gen` (lib.py + un script por pantalla) leyendo `datos-sinteticos.json`; el HTML resultante lleva los datos como texto o como datos de `renderVals()`.
- Tablas con `div role=table/row/cell` y CSS grid (no `<table>`): `<sc-for>` dentro de `<tbody>` sería reubicado por el analizador de HTML.
- Escritorio apila las pestañas por defecto y las demás para mostrar la pantalla completa (P-14: Resumen + Detalle; P-15: Resumen + Detalle + vista previa con Tendencias; P-16: Resumen + Detalle). Móvil muestra solo la pestaña por defecto (Resumen en P-14, Detalle en P-15 y P-16).
- «Tienda recomendada» de P-13 sale de `canasta.detalle[].ganadora`; no se afirma ahorro (el dato `diferencia_cents` no dice contra qué).
- Precio por unidad base: `precio_unidad_base_cents` en centavos por g, ml o pza; se muestra por kg, L o pza (x1000 en g/ml). Sam's lleva «con membresía» donde es la mejor tienda.
- El sintético trae 1 producto con frescura «viejo» (Pasta espagueti) y 0 sin precio; los 22 precios viejos se ven en las vistas previas por producto, no en el conteo.
- Alturas (`$preview` y fragmento) estimadas sin navegador (filas de alto fijo con texto recortado); se sobreestiman a propósito. Móvil: 2640, 2250, 1320, 1500; escritorio: 2560, 3350, 2800, 1700. Conviene que el orquestador las reajuste si las ve al renderizar.
- Tablero y galería (P-13, P-15, P-16) solo se muestran como botones de vista; se dibuja la Tabla.

## Bloqueos
Ninguno.

## Pendientes
- Renderizar para confirmar alturas y que `<sc-if>` dentro de `<label>`/`<span>` y los atributos `class="... {{r.cls}}"` y `aria-selected="{{r.selx}}"` se resuelven como dice format.md.
- `support.js` y `../ds/finanzas-hogar/tokens.css` los pone el orquestador (no existen en mi carpeta).
- Casillas e inputs con `checked` / `value` literales: en el runtime pueden quedar de solo lectura (artboards no interactivos).

## Consumo aproximado
Leído ~200 KB; escrito: ~67 KB de scripts y reporte por herramienta, más 193 KB de HTML generado en disco (por encima del tope de 150 KB si se cuenta lo generado: cada archivo repite barra lateral y estilos); 30 acciones aproximadas.
