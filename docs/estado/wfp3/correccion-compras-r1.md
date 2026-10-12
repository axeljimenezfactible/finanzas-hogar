# correccion-compras-r1 · Pista 3, fase 2 · vuelta 2 · 12-oct-2026

## Hecho
Copia de la vuelta 1 (P3-06) en `docs/design/salida/iter-1/compras/`: los 8 artboards y `compras.json`, tal como estaban antes de tocar nada.

Se corrigió el generador (`/tmp/compras_gen/lib.py`, `p13.py` a `p15.py`, `frag.py`) y se regeneraron los 8 artboards en `docs/design/maquetas/lienzo/project/compras/` y el fragmento `fragmentos/compras.json`. Hallazgos:

| Id | Corrección |
|---|---|
| 1 (alta) | La causa estaba en la función `shake` del generador: su expresión regular partía el `@media` y soltaba la llave de cierre. Ahora el bloque queda `@media (max-width:799px){.side{max-width:none;flex-basis:100%}}` y cada `<style>` de los 8 archivos tiene tantas `{` como `}` (escritorio: 95, 97, 100 y 92). Con eso vuelven a valer en 1440 la barra lateral, `.main`, las columnas `.c13` a `.c16`, `.cell`, `.cw`, `.sel` y `.ax` (el eje de P-15 ya no cae en fill negro). El generador además falla si un bloque queda desbalanceado. |
| 2 (alta) | Todas las fechas llevan año: «2-oct-2026 · hace 19 d» (P-13), fechas del Detalle de P-14 (`19-oct-2026`, `18-sep-2026`), «1 kg · 13-oct-2026 · ticket» (P-15). Semana: «Semana del Lun 19 oct 2026 al Dom 25 oct 2026» en el encabezado de P-13 (e y m) y en el selector de P-14. En el móvil de P-13 la fecha va en su propia línea bajo la presentación y el precio. En P-15 e se ensanchó la columna «Último precio» (200 px mínimo; tabla de 700 px) para que la fecha completa no se recorte. |
| 3 (media) | En P-14 (m y e) los montos dentro de frases van en `<span class="n">`: franja de la ganadora ($1,551.70 y $86.60), «+$86.60 contra la más barata (+$80.20 en idénticos)», franja de la membresía (−$212.40), «$650.00 al año» y el `select` con `class="n"` para las opciones. |
| 4 (media) | Barra inferior de los 4 móviles: el `<a>` del «+» es un destino de 64 × 56 (sin `width`, `height` ni `min-width:0`); el círculo de 40 px (`--size-bottombar-add`, `--radius-full`) vive en un `<span class="mas">` interior. Cada destino lleva su icono dentro de `<span class="pil">` de 56 × 32 (`--size-bottombar-indicator-width` y `-height`); el activo (Compras) pinta esa píldora en `--ui-surface-selected` con icono `--ui-icon-selected`. |
| 5 (media) | Barra lateral de los 4 escritorios, siguiendo el componente BarraLateral: las 9 secciones con `<button aria-expanded>` y chevron-right (gira 90° al desplegar; desplegado solo Compras y tiendas, con `aria-label` «Plegar»/«Desplegar» + nombre); pie con grupo «Tema» (sistema pulsado, claro, oscuro) y botón «‹‹» (`panel-left`, `aria-label="Colapsar barra lateral"`). |
| 6 (media) | El `h1` de los 4 móviles queda sin «· móvil» (Lista de compras, Comparar canasta, Productos y precios, Inventario). «· móvil» se queda en `<title>` y en el título del fragmento. |
| 7 (media) | Historial de Jitomate saladet (P-15 e): marcador por posición de la tienda (Walmart círculo, Bodega Aurrera cuadrado, Sam's triángulo, Mercado Portales rombo, Tianguis triángulo invertido; Sam's no tiene serie en este producto, así que se dibujan 4), tamaño 10 (`--size-viz-marker`; el círculo lee el token en `r`), relleno `--ui-surface-card` y trazo de la tienda. Etiqueta directa al final de cada línea (Walmart, Bodega, Portales, Tianguis), con separación vertical mínima y contorno del color de la tarjeta para que se lea sobre las líneas. La leyenda repite el mismo marcador. Tamaño de letra de eje y etiquetas: 12 px (antes 11, B-02 de la revisión). |

Alturas: el móvil de P-13 pasó de 2640 a 2950 (la fecha ocupa una línea más en 24 filas); el resto de `$preview` y del fragmento no cambió.

## Verificación
Scripts propios con `python3 -I` en `/tmp/compras_gen/` (`verify.py`, de la vuelta 1, y `verify2.py`, nuevo). Resultado: `RESULTADO: OK` y `RESULTADO2: OK`.
- Forma: doctype, `lang="es-MX"`, `<title>`, `./support.js`, enlace a `tokens.css`, `x-dc`, `helmet`, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`; `data-props` parsea y `$preview` coincide con el fragmento en los 8; `html.parser` sin etiquetas sin cerrar.
- CSS: llaves balanceadas en cada `<style>` y como máximo 1 `@media`, que contiene solo `.side`.
- Tokens: todos los `var(--x)` existen en `tokens.css` o en `tokens.json`; 0 hex fuera de `data-props` y de `href="#"`; 0 `rgb()` y `hsl()`.
- Cifras prohibidas de P3-12 (con y sin coma, y en centavos): 0.
- Fechas: 0 fechas cortas sin año y 0 rangos tipo «19–25 oct» en cuerpo y datos de los 8.
- Tabulares: todo texto con `$` (y cada campo de `renderVals()` con montos, usado en una plantilla) está dentro de un elemento con clase `n`; los `select` con montos llevan `class="n"`.
- h1 único y sin «móvil» ni «escritorio»; los `<title>` móviles conservan «· móvil».
- Móviles: barra inferior con `.mas` (40 px, círculo), 4 `.pil`, destino 64 × 56 y píldora activa; sin regla vieja `.bot .add`.
- Escritorios: botón «Colapsar barra lateral», 9 botones `aria-expanded` (1 en true), grupo de tema, fila «Estados»; P-15 e con 4 etiquetas directas.
- Fragmento: 8 rutas exactas en `boards` y `order`, sin duplicados; `expand: fill` solo en los `-e`; `is_interactive` false; `page` «compras»; títulos P-XX · Nombre · móvil|escritorio.

## Supuestos
- El tamaño del marcador es 10 unidades del viewBox (= `--size-viz-marker`, 10 px); el viewBox mide 380 y se pinta a ancho de la vista previa, así que en pantalla queda en unos 9 px. Los SVG no admiten `var()` en atributos geométricos; solo el círculo lo lee del token.
- Con día de la semana se usa «Lun 19 oct 2026» tal como lo propuso la revisión; el resto de fechas va como `2-oct-2026`.
- Los chevrons de las secciones distintas de Compras solo dibujan el botón; no se dibujan sus subpáginas.
- El móvil de P-13 pide 2950 de alto estimado sin navegador; conviene que el orquestador lo reajuste al renderizar.
- `support.js` y `../ds/finanzas-hogar/tokens.css` los instala el orquestador.

## Bloqueos
Ninguno.

## Pendientes
- Renderizar para confirmar alturas, la posición de las etiquetas directas de P-15 y que `aria-selected`, `class` y `<sc-if>` dentro de `<label>` se resuelven.
- Los hallazgos bajos B-01 a B-12 de la vuelta 1 no se tocaron (el encargo era alta y media), salvo B-02 (letra de 11 a 12 px en el eje de P-15).

## Consumo aproximado
Leído ~230 KB; escrito ~10 KB de ediciones al generador y scripts de verificación más el reporte; los 8 HTML regenerados suman ~215 KB en disco (el generador repite barra lateral y estilos en cada uno); ~22 acciones.
