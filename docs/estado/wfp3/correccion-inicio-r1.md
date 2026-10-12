# correccion-inicio-r1 (vuelta 2, Pista 3, área Inicio)

Fecha: 12-oct-2026. Alcance: artboards de `docs/design/maquetas/lienzo/project/inicio/`, fragmento `inicio.json` y este reporte.

## Hecho
- Copia de la vuelta 1 (P3-06) en `docs/design/salida/iter-1/inicio/`: los 8 `.dc.html` y `inicio.json` (9 archivos), hecha antes de editar.
- A-1 (alta, Q2-18 A, U3 tarjeta 6 y D-6): las 4 insignias de meta pasan al estado «En curso» con rol info (`bd bi`) e icono `info` (círculo con i), sin triángulo de advertencia.
  - `inicio-p01-resumen-m.dc.html` l. 107: «Atrasada» -> «En curso».
  - `inicio-p01-resumen-e.dc.html` l. 144: «Atrasada» -> «En curso».
  - `inicio-p01-detalle-m.dc.html` l. 150 («Atrasada») y l. 156 («En ruta») -> «En curso».
  - `inicio-p01-detalle-e.dc.html` l. 189 («Atrasada») y l. 195 («En ruta») -> «En curso».
- Nada más se tocó. Las cifras, textos y fechas de las metas no cambian (42 %, 78 %, aportación necesaria $2,783.34 y $650.00, fechas 17-ene-2027 y 17-abr-2027). B-1 a B-8 (todas baja) quedan fuera del alcance pedido (solo alta y media).

## Verificación (python3 -I, script propio en /tmp)
- 8/8 artboards: doctype, `lang="es-MX"`, charset, `./support.js`, enlace a `../ds/finanzas-hogar/tokens.css`, `<x-dc>`, `<helmet>`, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`.
- `data-props` parsea y `$preview` = w x h del fragmento en los 8.
- `html.parser`: 0 etiquetas sin cerrar ni cierres inesperados.
- 0 hex fuera de `data-props`; todos los `var(--x)` existen en `tokens.css`; 0 iframes, `innerHTML`, `fetch` ni URL externas.
- P3-12: 0 coincidencias de las 21 cifras canon (con y sin coma) en los 8 artboards y en `inicio.json`.
- 0 apariciones de «Atrasada» o «En ruta» en los artboards.
- Fragmento: JSON válido; `boards` y `order` listan exactamente los 8 archivos.

## Supuestos
- El orquestador corrige `datos-sinteticos.json` (next_goal.status -> in_progress / «En curso»; metas[].estado) y vuelve a correr `verificar_sinteticos.py`; yo no toqué ese archivo (fuera de mi alcance). En la fecha de este reporte aún trae `behind`/`Atrasada` y `atrasada`/`en_ruta`, así que los artboards ya van por delante del JSON hasta que se corrija.
- Las dos metas tienen fecha futura y no están cumplidas, así que ambas son «En curso» (no se muestra «Vencida» ni «Cumplida»).
- Los textos «Aportación mensual necesaria» se mantienen: son un dato, no un estado.

## Bloqueos
- Ninguno.

## Pendientes
- Orquestador: corregir `datos-sinteticos.json` y re-correr `verificar_sinteticos.py`; el cambio afecta también a P-09 (presupuesto, ahorro).
- Hallazgos de severidad baja B-1 a B-8 sin atender.

## Consumo aproximado
~75 KB leídos (sobre todo el JSON sintético y el informe de revisión), ~4 KB de artboards modificados más este reporte, ~10 acciones de herramienta.
