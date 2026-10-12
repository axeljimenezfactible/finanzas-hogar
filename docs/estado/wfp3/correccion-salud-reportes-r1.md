# correccion-salud-reportes-r1 · vuelta 2 · área «salud-reportes»

12-oct-2026 · Pista 3, fase 2 (lienzo). Corrige los 5 hallazgos medios de `P3-iter-1-salud-reportes.md`. No se rehízo nada más.

## Hecho

- Copia de la vuelta 1 (P3-06): `docs/design/salida/iter-1/salud-reportes/` (8 artboards y `salud-reportes.json`, antes de tocar nada).
1. Barra lateral: botón `Colapsar barra lateral` (icono panel-left, 20 px, `btn bt`) junto al chip de usuario en el encabezado de la barra de los 4 artboards de escritorio (P-21, P-23, P-24, P-25).
2. P-23, «Avisos enviados» (escritorio y móvil): filtro «Tipo» con Todos, Resumen diario, Resumen semanal, Reporte mensual, Alerta de pago próximo, Alerta de presupuesto al 80 %, Alerta de presupuesto al 100 %, Alerta de precios viejos y Alerta de plan sin armar; «Periodo» por defecto «Mes en curso».
3. P-23, barra de la vista en «Reportes generados» y «Avisos enviados». Escritorio: fila con selector de vista (Tabla, Tablero por tipo, Calendario por fecha; `aria-pressed`, Tabla activa), «Filtrar» (`aria-expanded`, controla el panel de filtros que queda dibujado debajo) y «Ordenar: Fecha, reciente primero»; reportes suma el filtro «Estado» (Todos, Listo, Con error, Generando). Móvil: selector de vista, botón «Filtrar» (hoja con Tipo, Periodo, Estado y Ordenar) y chips quitables bajo el título («Periodo: Mes en curso» en avisos; «Sin filtros activos» en reportes).
4. P-21 (escritorio y móvil): fila «Restricciones alimentarias» / «Ninguna» en el perfil de Axel y en el de Ximena.
5. P-24 (escritorio y móvil): «Canasta de la lista de compras de la semana del 19 al 25 oct 2026.»
- Alturas: P-23 móvil 2300 a 2600, P-23 escritorio 2300 a 2450, P-21 móvil 2200 a 2300 (min-height, `$preview` y fragmento coinciden).

## Verificación (`python3 -I`, scripts propios en /tmp, fuera del repositorio)

- Forma: 8/8 con doctype, `lang="es-MX"`, charset, `support.js`, `tokens.css`, `x-dc`, `helmet`, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`; `data-props` parsea y `$preview` = (w, h) del fragmento.
- `html.parser`: 0 etiquetas sin cerrar ni cierres inesperados; 0 IDs duplicados; todo `select` con su `label for`.
- Tokens: 0 hex literales fuera de `data-props` (excluidos `href="#…"`), 0 `rgb()/hsl()`; todos los `var(--x)` existen en `tokens.css`.
- P3-12: grep de las cifras canon (con y sin coma) en los 8 y el fragmento = 0; 0 credenciales; 0 coincidencias de «W43»/«ISO».
- Fragmento: JSON válido; lista exactamente los 8 archivos; móvil sin `expand`, escritorio `"expand":"fill"`.
- Anchos: raíz móvil 412 px; escritorio fluido (sin cambios).

## Supuestos

- «Mes en curso» como Periodo por defecto solo en avisos (U1 §2.5.b, notificacion); en reportes queda «Todos» (U1: job, «según la pantalla»), porque la lista incluye el reporte de ago-2026.
- Opción de Estado «Generando» y de Ordenar «Fecha, antiguo primero»: no vienen en datos-sinteticos.json; son opciones de interfaz, no cifras.
- El panel de filtros de escritorio se dibuja abierto para que las opciones se vean; en móvil la hoja va cerrada (`details`) y sus opciones están en el marcado.
- Alturas móviles estimadas sin navegador (no hay render en este flujo).

## Bloqueos

Ninguno.

## Pendientes

- Los 14 hallazgos bajos de la vuelta 1 (B01 a B14) no se tocaron: el encargo era solo alta y media.
- Las vistas Tablero y Calendario de P-23 se ofrecen en el selector pero no se dibujan (el artboard muestra Tabla).

## Consumo aproximado

~85 KB leídos, ~25 KB escritos (reporte, scripts, parches), ~22 acciones.
