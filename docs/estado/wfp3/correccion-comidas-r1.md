# correccion-comidas-r1 · Vuelta 2 del área «comidas»

Fecha: 12-oct-2026 · Agente: `correccion-comidas-r1` · Fuente: `P3-iter-1-comidas.md` (M-1 a M-4)

## Hecho

- Copia de la vuelta 1 (exigida por P3-06) en `docs/design/salida/iter-1/comidas/`: los 6 artboards y `comidas.json`, copiados antes de editar.
- **M-1 (enlaces del mismo ancho).** Los 17 `href` a la receta ahora apuntan a `../comidas/comidas-p12-receta-m.dc.html`: 5 en `comidas-p10-semana-m` y 12 en `comidas-p11-recetario-m`. Ya no queda ningún enlace a `-e` en los artboards móviles de comidas.
- **M-2 (objetivo táctil).** El enlace «‹ Comidas» o «‹ Recetario» de la barra superior de `p10-m`, `p11-m` y `p12-m` lleva la clase `bk`: `min-height: var(--size-touch-min)` y `padding: 0 var(--space-2)`.
- **M-3 (selector de recetas, N-122, N-123, F-2).**
  - Móvil: `p10-m` suma una hoja inferior «Elegir receta · Lun · Comida» al final del contenido, con fila de 7 días (Lun marcado, cada botón de 44 px o más), persona (Ambos por defecto, Axel, Ximena), porciones por persona (1, que con Ambos son 2), filtro de texto y las 5 recetas aptas para Comida.
  - Escritorio: `p10-e` suma una tarjeta con el mismo contenido, en dos columnas (días, persona y porciones a la izquierda; búsqueda y lista a la derecha), antes de la tarjeta «Estados».
  - Las recetas van ordenadas por veces usadas en 8 semanas (descendente) y luego por nombre, con costo y kcal por porción, tomados de `recetas[]` de `datos-sinteticos.json`: Pollo asado con arroz y ensalada ($30.82, 807 kcal, 6 veces), Pasta con atún y jitomate ($27.72, 639, 4), Frijoles de la olla con aguacate ($14.23, 831, 2), Picadillo de res con papa ($34.38, 743, 1) y Sopa de verduras ($12.43, 247, 1).
- **M-4 (barra de la vista en P-10 escritorio).** `p10-e` suma, sobre la cuadrícula, el grupo «Vista del plan» con Calendario (activa, `aria-pressed="true"`), Tabla, Tablero y Galería, con los iconos Lucide calendar, table-2, kanban y layout-grid. Suma también «Filtrar» (embudo, `aria-expanded="false"`) y el filtro «Persona» con Todas, Axel y Ximena. El selector de semana ya estaba justo encima, en la sección «Semana elegida».
- Fragmento `comidas.json`: `p10-m` pasa de h 2000 a 2700 y `p10-e` de h 1700 a 2300. `$preview` de cada artboard coincide con el fragmento.

## Verificación

Script propio `python3 -I` (en `/tmp/vcom/v.py`, fuera del proyecto), resultado: sin fallas. Revisó, en los 6 archivos:
- Líneas obligatorias de la forma, `lang="es-MX"`, `support.js` antes de `tokens.css`, `<x-dc><helmet><style>`, y `class Component extends DCLogic` con `renderVals()`.
- El JSON de `data-props` parsea y `$preview` coincide con `w` y `h` del fragmento (412 en móvil; 1440 y `expand: fill` en escritorio).
- `html.parser`: 0 etiquetas sin cerrar ni cierres inesperados; un solo `h1` por archivo; ids sin duplicar y todos los `aria-labelledby` apuntan a un id existente.
- 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`; 0 iframes, `innerHTML`, URL externas y `{{`; 0 emoji.
- 0 cifras prohibidas de P3-12 (con y sin coma) y 0 coincidencias de credenciales.
- Todos los `var(--x)` existen en `sistema/project/tokens.css`.
- Todo monto tiene la forma `$1,234.56`.
- Ningún artboard móvil enlaza a un `-e`; ningún enlace a comidas está roto.
- Presencia del selector (título, recetas, costos, «Ambos», búsqueda) en `p10-m` y `p10-e`; la barra de la vista y el icono calendar en `p10-e`; la clase `bk` y su regla en los 3 artboards móviles.
- Fragmento válido; lista exactamente los 6 archivos de comidas.

## Supuestos

- La hoja y el diálogo se dibujan en el flujo de la página (no como superposición), con una nota que dice cuándo aparecen; así el lienzo los muestra sin necesidad de interacción.
- «Filtrar» queda como botón más un grupo de chips de persona (Todas, Axel, Ximena) visibles, no como panel desplegado.
- `Poner en el plan · 2 porciones` refleja Ambos con 1 porción por persona (N-123).
- El selector de semana de `p10-e` no se duplicó en la barra de la vista; ya está inmediatamente encima.
- Los hallazgos bajos (B-1 a B-12) no se tocaron, según el encargo (solo alta y media).

## Bloqueos

Ninguno.

## Pendientes

- Medir el alto real de `p10-m` y `p10-e` en el lienzo: los nuevos `h` (2700 y 2300) son estimaciones sin navegador, con margen.
- Los 12 bajos de la vuelta 1 siguen abiertos si la revisión de la vuelta 2 los pide.
- Los artboards suman unos 167 KB en disco por las dos vueltas acumuladas; sobre este agente solo se escribieron unos 14 KB en esta vuelta.

## Consumo aproximado

Unos 75 KB leídos y unos 14 KB escritos (artboards, fragmento y este reporte), con unas 17 acciones de herramienta. Sin Artifact, sin navegador, sin red, sin git ni Telegram.
