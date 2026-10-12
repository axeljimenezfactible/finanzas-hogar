# correccion-presupuesto-r2 · vuelta 3 · 12-oct-2026

## Hecho
- Copia previa (P3-06): los 8 artboards de presupuesto y el fragmento `presupuesto.json` de la vuelta 2 quedaron en `docs/design/salida/iter-2/presupuesto/` antes de editar.
- **N1 (alta) cerrado con la opción (a).** En `presupuesto-p07-flujo-m.dc.html` y `presupuesto-p07-flujo-e.dc.html`:
  - Halo bajo la línea a guiones del devengado acumulado: mismo `d` (`chart.devLine`), `stroke:var(--ui-surface-card)`, ancho `calc(var(--size-viz-line) * 3)` = 6 px.
  - La línea «Hoy» pasó de un `div` con `border-left` a un `path` dentro del mismo SVG (`chart.hoyPath`, calculado en `renderVals()` con la misma proporción 5.4 %), con halo de tarjeta de 6 px debajo y trazo `ui-viz-baseline` de `--border-width-default` encima. El rótulo «Hoy · 21-oct-2026» sigue igual, posicionado con `chart.hoyLeft`.
  - Orden de pintura: rellenos, halo del devengado, halo de «Hoy», línea «Hoy», línea del efectivo, línea a guiones. Así las series quedan por encima y el color contiguo a las dos líneas es la tarjeta.
- Fragmento sin cambios (mismos 8 archivos, mismas `w`/`h`).

## Verificación
- Contraste recalculado con `tokens.css` (WCAG, línea contra tarjeta, el color que ahora la rodea): `ui-chart-2` 3.15:1 claro y 7.73:1 oscuro; `ui-viz-baseline` 3.47:1 claro y 6.09:1 oscuro. Todos ≥ 3:1 (P3-07, U5-04; igual que C-107 y base de §10.5).
- Script propio `python3 -I` (`/tmp/p3pres2/verify.py`, fuera del repo): **360 comprobaciones, 0 fallas** sobre los 8 artboards y el fragmento: cabecera exacta, `support.js`, `tokens.css`, `data-props` JSON válido con `$preview` 412/1440, `class Component extends DCLogic`, html.parser sin etiquetas sin cerrar, 0 hex fuera de `data-props`, todos los `var(--x)` existen en tokens, 0 cifras prohibidas de P3-12 (con y sin coma) ni centavos canon, sin `innerHTML`/iframe/red/emoji, huecos `{{}}` resueltos con `renderVals()` en node (`hoyPath` resuelve: `M15.05,0 L15.05,200` móvil, `M16.13,0 L16.13,280` escritorio), enlaces dentro de la lista, fragmento válido.
- No se renderizó en navegador (no permitido); el halo se verificó por estructura y por cálculo de contraste.

## Supuestos
- Se eligió (a) y no (b) para conservar el devengado acumulado visible aun cuando la capa 2 vale 0.
- Halo de 6 px = 3 × `--size-viz-line` (dentro del rango de 4 a 6 px pedido), sin token nuevo.

## Bloqueos
- Ninguno.

## Pendientes
- Medir el contraste real sobre el render publicado en claro y oscuro (la altura estimada B4 sigue por medir).

## Consumo aproximado
- Leído ≈ 60 KB; escrito ≈ 3 KB en artboards (2 parches) más el reporte; ≈ 14 acciones.
