# P3-iter-2-salud-reportes · Revisión de la vuelta 2 · área «salud-reportes»

Revisor: `revisor-salud-reportes-r2` (independiente; no corrige) · 12-oct-2026 · Pista 3, fase 2 (lienzo), vuelta 2 de 3 · Criterio P3-07.

Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/salud-reportes/` y `docs/design/maquetas/lienzo/fragmentos/salud-reportes.json`, contra P3 1.3 (P3-06, P3-08, P3-12), U5 (U5-04, U5-06, U5-07, en lo que se ve en el marcado), DESIGN.md (§3.3, §7.4, §7.5.1, §7.5.2, §7.9, §7.10, §10.5), U1 (§1.2, §1.3, §2.1, §2.2, §2.5.a/b, N-230, N-245 a N-248, §6.4) y `referencia-tipo/format.md`. U3 y U4 no aplican (el área no tiene Inicio ni Ajustes). Línea base: la copia de la vuelta 1 en `docs/design/salida/iter-1/salud-reportes/`, comparada etiqueta por etiqueta con lo actual.

## Veredicto

**No aprobado.** 0 altos · **1 medio** · 17 bajos.

Los 5 medios de la vuelta 1 (M01 a M05) están **cerrados**. Queda abierto un medio nuevo que la vuelta 1 no vio: en P-25 escritorio, el dinero de las tablas no usa cifras tabulares (N01). Los arreglos de la vuelta 2 no rompieron nada de forma, tokens, contraste, cifras ni enlaces.

## Hallazgos de la vuelta 1: comprobación

| ID | Hallazgo | Estado | Evidencia |
|---|---|---|---|
| M01 | Barra lateral sin «‹‹» | **Cerrado** | Los 4 `-e` traen `<button class="btn bt" type="button" aria-label="Colapsar barra lateral">` con `panel-left` de Lucide (`rect 18×18 rx2` y `M9 3v18`, `class="i"`, 20 px, `aria-hidden`), junto al chip de usuario en el encabezado de la barra (DESIGN §7.9). Contraste del icono: `ui-text-secondary` sobre `ui-surface-card`, 5.04 en claro y 7.97 en oscuro. |
| M02 | Filtro «Tipo» de avisos con tipos de reporte | **Cerrado** | `at` (escritorio) y `amt` (móvil): Todos, Resumen diario, Resumen semanal, Reporte mensual, Alerta de pago próximo, Alerta de presupuesto al 80 %, Alerta de presupuesto al 100 %, Alerta de precios viejos, Alerta de plan sin armar. Coincide con N-247; «reporte mensual» es el literal de `reportes.proximos_envios[2].tipo`. Periodo por defecto «Mes en curso» (§2.5.b, notificacion). |
| M03 | P-23 sin barra de la vista ni filtro Estado; filtros de móvil en línea | **Cerrado** | Escritorio: fila `role="group"` con Tabla, Tablero por tipo y Calendario por fecha (`aria-pressed`, Tabla activa), y además «Filtrar» (`list-filter`, `aria-expanded="true"`, `aria-controls` a un panel que existe) y «Ordenar: Fecha, reciente primero» (`arrow-up-down`), en las dos listas. Reportes suma «Estado» (Todos, Listo, Con error, Generando). Móvil: el mismo selector de vista, un `summary` «Filtrar» que abre Tipo, Periodo, Estado y Ordenar, y debajo los filtros activos («Sin filtros activos.» en reportes; chip quitable «Periodo: Mes en curso» con `aria-label="Quitar el filtro Periodo: Mes en curso"` en avisos). Coincide con DESIGN §7.5.1 y §7.4 (quitable). |
| M04 | P-21 sin «Restricciones alimentarias» | **Cerrado** | `<dt>Restricciones alimentarias</dt><dd>Ninguna</dd>` en el `dl` de Axel y en el de Ximena, en `-e` y `-m` (N-230, decisión D). |
| M05 | «semana 2026-W43» visible en P-24 | **Cerrado** | Ahora dice «Canasta de la lista de compras de la semana del 19 al 25 oct 2026.» en `-e` y `-m`. Grep de `\d{4}-\d{2}-\d{2}`, `\d{4}-W\d{2}` y `\d{4}-\d{2}` en el texto visible y en `aria-label`/`title` de los 8: 0. El 19-oct-2026 es lunes y el 25, domingo. |
| B01 a B14 | Bajos de la vuelta 1 | Abiertos | La corrección solo tocó los medios (lo dice su reporte). Se confirman sin cambios: B01 («con telegram», «Manual (axel)»), B03 («Error»), B06 («Editar perfil» y «Cambiar en Ajustes › Reportes y frecuencias» a `#`, aunque ya existen `ajustes-p26-s13-salud-metas-*` y `ajustes-p26-s08-reportes-*`), B07 (`is_interactive: true` y radios con `checked` en P-21 móvil), B09 (`border-bottom:2px`). Los demás también siguen igual: el diff no los toca. |

## Verificación propia (scripts `python3 -I` en `datos/tmp/rev-sr-r2/`, borrados al final)

| Punto | Resultado |
|---|---|
| (a) Forma | Los 8 empiezan por `<!doctype html>` y traen `lang="es-MX"`, `<meta charset="utf-8">`, `<title>` propio, `./support.js`, `../ds/finanzas-hogar/tokens.css`, `<x-dc>`, `<helmet><style>` dentro de `x-dc` y `data-dc-script` con `class Component extends DCLogic` y `renderVals()`. `data-props` parsea en los 8, y `$preview` coincide con `w`/`h` del fragmento en los 8 (P-21 m 412×2300, P-23 m 412×2600, P-23 e 1440×2450, el resto igual que en la vuelta 1). `html.parser`: 0 etiquetas sin cerrar, 0 cierres inesperados. IDs duplicados: 0 (los nuevos `rt`, `rp`, `rs`, `at`, `ap`, `f-r`, `f-a`, `rmt`…`amo` son únicos). Lienzo: 90 artboards, 0 stems repetidos, 0 nombres inválidos. |
| (b) Tokens | 0 hex literales fuera de `data-props` y 0 `rgb()`/`hsl()`. Entre 45 y 55 `var(--x)` por archivo: todos están en `tokens.css`, y todos en `tokens.json` salvo `--font-sans` (B14, de todo el lienzo). Los tokens nuevos (`--ui-border-selected`, `--font-weight-semibold`, `--size-touch-min`, `--radius-sm`) existen. Trazo del icono: 1.8 a 20 px = 1.5 px dibujados (DESIGN §5, `icon.stroke.user-units.sm`). |
| (c) Contraste | Cascada CSS propia (clases del `helmet`, especificidad y estilo en línea; color heredado y fondo del ancestro), con los hex de `tokens.json` en claro y oscuro (referencias resueltas). Texto: 19 pares en el árbol y 19 en las reglas con color y fondo; **0 bajo 4.5**. El más bajo es `ui-role-success` sobre `ui-role-success-surface`: 4.81 y 6.19. Interfaz: 18 pares de borde, marca y relleno, **0 bajo 3**. El más bajo es `ui-person-axel` sobre tarjeta en oscuro: 3.10. Pares nuevos de la vuelta 2: `ui-text-default`/`ui-surface-selected` (vista activa) 10.73 y 10.90; `ui-border-selected`/`ui-surface-card` 5.16 y 7.94; `ui-border-selected`/`ui-surface-page` 5.53 y 10.38; `ui-border-control`/`ui-surface-card` (Filtrar, Ordenar) 3.47 y 6.09; `ui-text-secondary`/`ui-surface-raised` (quitar chip) 5.40 y 7.97. Por debajo de 3 solo quedan fondos sin función de control: superficies de rol de insignias y franjas (su texto sí pasa), esqueleto, retícula, pista y halo de marcadores (decorativos, §10.5). |
| (d) P3-12 | Grep de las 21 cifras canon (con y sin coma) en los 8 y en el fragmento: **0**. Grep de credenciales: 0. Hay 57 importes distintos y todos salen de `datos-sinteticos.json` (`*_fmt` o `*_cents`). Las 21 fechas distintas visibles están todas en el JSON como ISO. Las fechas que coinciden con el calendario real (15, del 30-oct-2026 al 15-jun-2027) aparecen en P-24 y P-25 solo con $381.84, $453.49, $803.27 y $803.26, que no son cifras del reporte ni de Kueski. |
| (e) Fidelidad | Tabla abajo. |
| (f) Accesibilidad | 0 `a`/`button`/`summary` sin nombre accesible; 0 campos sin `label`; todos los `<svg>` llevan `aria-hidden` o `role="img"` con nombre; 0 emoji; no hay barra de estado ni teclado falsos. Objetivos táctiles en móvil: 99 controles medidos por la cascada; todos miden ≥ 44 px. Las únicas dos excepciones son los `input type="radio"` ocultos (`.sr`) de las pestañas de P-21, y sus `label` miden 44 px. Dinero: 57 importes con forma `$1,234.56`; signo U+2212 4 veces; 0 negativos con guion ASCII. **Cifras tabulares: 30 apariciones sin `tabular-nums`**: 22 en las tablas de P-25 escritorio (N01) y 8 en el texto de los avisos de P-23 (N02). |
| (g) Anchos | Raíz móvil `width:412px` en los 4 (`box-sizing` por `*`). Escritorio fluido en los 4: raíz sin ancho, `flex-wrap`, barra `max-width: var(--layout-sidebar-width)` y contenido `flex:999 1 560px; min-width:0`. `"expand":"fill"` solo en los 4 de escritorio. |
| (h) Enlaces | 0 `href` fuera de la lista del lienzo, 0 de un ancho al otro (`-m` a `-e` o al revés), 0 anclas sin destino. Los `#` que tienen destino son los de B06; «Buscar» va a `#` porque no hay artboard de búsqueda; «Registrar gasto» y «+» ya van a `shell-p03-registrar-gasto-{m,e}`. |

## P3-08 · Lista de 8 puntos

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen ≤ 6 tarjetas | No aplica |
| 2 | 14 secciones de Ajustes | No aplica |
| 3 | Barra lateral colapsable | **Cumple** (M01 cerrado) |
| 4 | Español de México | Cumple, salvo B01 y B03 |
| 5 | Dinero `$1,234.56` con cifras tabulares | **Falla parcial**: el formato está bien en los 57 importes, pero las tablas de P-25 escritorio no usan `tabular-nums` (N01) |
| 6 | Fechas `10-oct-2026` | **Cumple** (M05 cerrado; 0 ISO visibles) |
| 7 | Paleta de C2 | Cumple (solo `var(--…)`) |
| 8 | Notion como única referencia | Cumple |

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presentes | Faltantes o distintos |
|---|---|---|
| P-21 | Todo lo de la vuelta 1, y además «Restricciones alimentarias · Ninguna» en los dos perfiles (N-230). Estados vacío, carga y error. | «Editar perfil» va a `#` (B06) |
| P-23 | Generar reporte mensual (N-248); próximos envíos (N-246); reportes generados (N-245) con barra de la vista, filtros Tipo, Periodo y Estado, y Ordenar; fila en error con «Reintentar envío» (N-250); avisos enviados (N-247) con barra de la vista, Tipo según N-247, Periodo «Mes en curso» y chip quitable en móvil; Karen «Sin envío · se comparte a mano»; estados. | Avisos abiertos en escritorio (B11); acciones secundarias de móvil a todo el ancho (B13) |
| P-24 | Selector de mes, 7 secciones plegables, Generar PDF, Compartir deshabilitado, estados por sección; canasta «de la semana del 19 al 25 oct 2026». | Mes por defecto (B10); vacío sin acción (B02) |
| P-25 | Lo de la vuelta 1 completo (periodo, compras, partes de despensa, pagos, saldo, nota de Kueski, Generar PDF, Compartir, estados anteriores, estados). | Tablas de escritorio sin cifras tabulares (N01); vacío sin acción (B02) |

## Hallazgos

### Medio (bloquea la aprobación)

- **N01 · P3-08 punto 5 / U5-06 («cifras tabulares en todo dinero») / DESIGN §3.3 y §7.5.2.** En `reportes-p25-karen-e.dc.html`, las 3 tablas (compras particulares, su parte de cada despensa y pagos) muestran 22 importes en `<td class="r">` sin `font-variant-numeric: tabular-nums`. El `helmet` de este archivo solo define `.m{…tabular-nums}` y `.r{text-align:right}`. Falta la regla `.tbl td,dl.kv dd{font-variant-numeric:tabular-nums}`, que sí traen `reportes-p25-karen-m`, P-24 y P-21. Inter usa cifras proporcionales por omisión, así que las columnas de dinero no alinean. La vuelta 1 no lo detectó: su archivo ya era así. **Corrección:** agregar al `helmet` `.tbl td,dl.kv dd{font-variant-numeric:tabular-nums}`, o `class="r m"` en las celdas de dinero.

### Bajos

- **N02 · DESIGN §3.3 («también … fechas en tablas»).** En P-23 (`-e` y `-m`), el «Texto enviado» de los avisos trae 8 importes en prosa sin cifras tabulares ($1,685.50, $1,800.00, $951.30, $645.00, $1,650.98, $2,000.00). En `reportes-p23-centro-e`, las fechas y horas de las tablas tampoco las usan, porque falta la regla `.tbl td`. Corrección: la misma regla de N01, o `class="m"` en la celda del texto.
- **N03 · ARIA de la barra de la vista.** En escritorio, «Filtrar» declara `aria-haspopup="dialog"`, pero controla un panel en línea (`role="group"`). «Ordenar» declara `aria-haspopup="menu"` y no hay menú dibujado. En móvil, `<div class="row wrap" aria-label="Filtros activos">` lleva un `aria-label` sin rol, y en un `div` genérico ARIA 1.2 no lo admite. Corrección: quitar `aria-haspopup` del «Filtrar» de escritorio (basta `aria-expanded` con `aria-controls`); en «Ordenar», dibujar el menú o quitar el atributo; poner `role="group"` en el contenedor de chips.
- **N04 · DESIGN §7.5.1 y §7.12 (hoja).** En móvil, «Filtrar» abre un `details` en línea y no una hoja inferior. Para un artboard estático alcanza. Además, en «Avisos enviados» la barra de la vista y el chip quedan dentro del `details` plegado (§6.4), así que no se ven en el artboard. Es una observación, sin cambio obligatorio.
- **B01 a B14** (de la vuelta 1, sin cambios): B01 textos distintos del JSON («telegram», «Manual (axel)»); B02 vacíos sin la acción que los llena; B03 «Error» visible; B04 iconos de 20 px en la barra inferior; B05 de `h1` a `h3` en P-21, `summary` sin encabezado y medidores sin `aria-label` en P-24, chip de usuario sin menú, chevrons sin `aria-expanded`; B06 «Editar perfil» y «Cambiar en Ajustes › Reportes y frecuencias» a `#`; B07 `is_interactive: true` y radios con `checked` en P-21 móvil; B08 altos estimados; B09 `border-bottom:2px`; B10 P-24 abre en oct-2026 (parcial); B11 avisos (P-23) y movimientos (P-25) abiertos en escritorio; B12 ejes del gráfico de peso a ≈ 11.2 px en móvil; B13 secundarias de móvil a todo el ancho; B14 `--font-sans` falta en `tokens.json`.

## Tabla artboard × regla

Reglas: F forma (format.md), T tokens, C contraste U5-04, D cifras P3-12, Fi fidelidad U1, A accesibilidad U5-06, L idioma y formatos U5-07, R raíz y ancho, E enlaces, P8 P3-08.

| Artboard | Cumple | Falla (ID) | Correcciones para la vuelta 3 |
|---|---|---|---|
| `salud-p21-salud-m` | F, T, C, D, Fi, R, E, L, P8 | A (B04, B05), B07, B08, B09, B12 | Opcional: iconos de 24 px, `h2` por persona, pestañas con `state` o `is_interactive:false` |
| `salud-p21-salud-e` | F, T, C, D, Fi, R, L, P8 | A (B05), E (B06), B09 | Opcional: «Editar perfil» a `../ajustes/ajustes-p26-s13-salud-metas-e.dc.html` |
| `reportes-p23-centro-m` | F, T, C, D, Fi, R, P8 | L (B01, B03), A (N03), N02, N04, B04, B06, B09, B13 | Opcional: textos literales del JSON, `role="group"` en los chips |
| `reportes-p23-centro-e` | F, T, C, D, Fi, R, P8 | L (B01, B03), A (N03), N02, B02, B06, B09, B11 | Opcional: regla `.tbl td` con `tabular-nums`, `aria-haspopup` |
| `reportes-p24-mensual-m` | F, T, C, D, R, E, L, P8 | A (B04, B05), B09, B10, B13 | Opcional |
| `reportes-p24-mensual-e` | F, T, C, D, R, E, L, P8 | A (B05), B02, B09, B10 | Opcional |
| `reportes-p25-karen-m` | F, T, C, D, R, E, L, Fi, P8 | B04, B09, B13 | Opcional |
| `reportes-p25-karen-e` | F, T, C, D, R, E, L, Fi | **P8/A (N01)**, B02, B09, B11 | **Obligatorio:** `.tbl td,dl.kv dd{font-variant-numeric:tabular-nums}` en el `helmet` |

## Fragmento

`salud-reportes.json` es JSON válido. `boards` y `order` listan exactamente los 8 archivos, en el mismo orden. `w`/`h` coinciden con `$preview`, incluidos los 3 altos que cambiaron (P-21 m 2300, P-23 m 2600, P-23 e 2450). Títulos «P-XX · Nombre · móvil|escritorio»; `page` «salud-reportes»; `expand: fill` solo en escritorio; `is_interactive` `false`, salvo P-21 móvil (B07). Grep de cifras canon: 0.

## Texto con aspecto de instrucción (P3-04)

No obedecí ninguno. `referencia-tipo/format.md` se leyó por secciones y como datos de forma. Los artboards no traen texto dirigido a un agente.

## Líneas de evidencia

- `P3-04`: el revisor no usó navegador, la herramienta Artifact, red ni Telegram; solo Bash para leer y escribir sus dos reportes. 0 credenciales.
- `P3-06`: 8 de 8 artboards (4 pantallas × 2 anchos); 0 faltantes. La copia de la vuelta 1 está en `docs/design/salida/iter-1/salud-reportes/` (la hizo el corrector). La copia de este informe en `docs/design/salida/iter-2/` le toca al orquestador.
- `P3-07`: este informe. Contraste recalculado en los dos temas: 0 pares de texto bajo 4.5 y 0 de interfaz bajo 3.
- `P3-08`: de la lista de 8 puntos, cumplen 3 y 6 (cerrados), falla parcial el 5 (N01) y no aplican 1 y 2. El grep `\$[0-9]{4,}\.[0-9]{2}` (miles sin coma) da 0.
- `P3-12`: 0 cifras canon, 0 credenciales. Los 57 importes y las 21 fechas salen de `datos-sinteticos.json`. Las fechas coincidentes van con montos no canon.
- `U5-04`: 0 pares bajo umbral (estático, con tokens). `U5-06`: N01 y N03; B04 y B05 siguen abiertos. `U5-07`: 0 fechas ISO visibles; quedan B01 y B03.

## Supuestos

1. El informe va en `docs/estado/wfp3/`, como pide la tarea. P3-07 nombra `docs/estado/P3-iter-N-<area>.md`; el orquestador decide si lo copia.
2. Medí el contraste en estático, con una cascada CSS propia sobre el marcado y los tokens resueltos, no en un DOM renderizado (K45: sin navegador).
3. N01 cuenta como **medio** y no como alto: la decisión 5 de P3-08 no se contradice en el diseño, se omite una regla CSS en un artboard. Es el mismo criterio que en la vuelta 1 se usó con la barra lateral sin «‹‹». Si el orquestador lo lee como «decisión contradicha», sube a alto, y el veredicto no cambia.
4. N01 es nuevo pero no nace de la vuelta 2: el archivo de la vuelta 1 ya era así. Lo registro ahora porque la revisión es independiente y el criterio de aprobación cuenta todos los medios verificados.
5. «Reporte mensual» en el filtro de avisos se acepta: es el literal del JSON (`proximos_envios[2].tipo`), aunque N-247 dice «resumen … mensual».
