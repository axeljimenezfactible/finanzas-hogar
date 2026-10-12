# P3-iter-3-salud-reportes · Revisión de la vuelta 3 · área «salud-reportes»

Revisor: `revisor-salud-reportes-r3` (independiente; no corrige) · 12-oct-2026 · Pista 3, fase 2 (lienzo), vuelta 3 de 3 · Criterio P3-07.

Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/salud-reportes/` y `docs/design/maquetas/lienzo/fragmentos/salud-reportes.json`. Los comparé con P3 1.3 (P3-04, P3-06, P3-07, P3-08, P3-12), U5 (U5-04, U5-06 y U5-07, en lo que se puede ver en el marcado), DESIGN.md (§3.3, §4.5, §5.2, §7.1, §7.4, §7.5, §7.6, §7.9, §9.3, §10.5), U1 (§1.3, §2.1, §2.2, §3.23 a §3.27, F-5 y F-6) y `referencia-tipo/format.md`. U3 y U4 no aplican, porque el área no tiene Inicio ni Ajustes. Línea base: la copia de la vuelta 2 en `docs/design/salida/iter-2/salud-reportes/`.

## Veredicto

**Aprobado.** 0 altos · 0 medios · 19 bajos (abiertos, ninguno bloquea).

- El único medio de la vuelta 2 (N01, cifras tabulares en las tablas de P-25 escritorio) está **cerrado**.
- Con `cmp` contra la copia de iter-2, el único archivo que cambió es `reportes-p25-karen-e.dc.html`. El diff es una sola regla del `helmet`: `.tbl td{font-variant-numeric:tabular-nums;…}`. El fragmento no cambió.
- La corrección no rompió nada de forma, tokens, contraste, cifras ni enlaces.

## Hallazgo de la vuelta 2: comprobación

| ID | Hallazgo | Estado | Evidencia |
|---|---|---|---|
| N01 (medio) | P-25 escritorio: 22 importes de las 3 tablas sin `tabular-nums` | **Cerrado** | El `helmet` de `reportes-p25-karen-e.dc.html` define ahora `.tbl td{font-variant-numeric:tabular-nums;padding:var(--space-3);…}`. Usé una cascada CSS propia (especificidad, herencia de `font-variant-numeric` y estilo en línea). Con ella, **0 importes sin cifras tabulares** en P-25 escritorio, y lo mismo en las 11 celdas con fecha o porcentaje. Las 3 tablas son compras particulares, su parte de cada despensa y movimientos detallados. El archivo no trae `<dl>`, así que no le hace falta `dl.kv dd`. |
| N02 a N04, B01 a B14 (bajos) | — | Siguen abiertos | La corrección no los tocó: su encargo eran los altos y los medios. Los vuelvo a confirmar abajo. |

## Verificación propia

Corrí scripts `python3 -I` en `datos/tmp/rev-sr-r3/` y los borré al final.

| Punto | Resultado |
|---|---|
| (a) Forma | **Los 8 cumplen:** <ul><li>Empiezan exactamente con `<!doctype html>\n<html lang="es-MX"><head><meta charset="utf-8"><title>`.</li><li>Cada uno tiene `<title>` propio, y los 8 son distintos.</li><li>Traen `<script src="./support.js"></script>` y `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`.</li><li>Tienen un solo `<x-dc>`, con `<helmet><style>` como primer hijo.</li><li>Tras `</x-dc>` viene `<script type="text/x-dc" data-dc-script data-props='…'>`, con `class Component extends DCLogic` y `renderVals()`.</li><li>`data-props` parsea como JSON en los 8.</li><li>`$preview` coincide con `w`×`h` del fragmento en los 8: P-21 m 412×2300, P-21 e 1440×2400, P-23 m 412×2600, P-23 e 1440×2450, P-24 m 412×3400, P-24 e 1440×3100, P-25 m 412×2000 y P-25 e 1440×2000. Ningún móvil mide menos de 915 de alto.</li></ul>`html.parser` con pila propia: 0 etiquetas sin cerrar y 0 cierres inesperados. IDs duplicados: 0. Tampoco hay `{{`, `<iframe`, `innerHTML`, `http(s)://`, `fetch(` ni `<script>` en línea. En el lienzo hay 90 artboards: 0 stems repetidos y 0 nombres fuera de `[a-z0-9.-]`. |
| (b) Tokens | 0 hex literales fuera de `data-props` y 0 `rgb()`, `hsl()`, `oklch()` o `color-mix()`. Los atributos `fill`/`stroke` del SVG solo valen `none`. En CSS, el SVG usa `currentColor` o `var(--…)`. Hay entre 45 y 55 `var(--x)` por archivo, 67 distintos en total. Todos existen en `tokens.css`, y todos en `tokens.json` salvo `--font-sans` (B14). `--font-sans` sale de `type.families.sans` de `tokens.json`, pero no figura como `name`. Ningún `px` literal en tipografía, márgenes, rellenos ni `gap`. |
| (c) Contraste | La misma cascada, con los hex de `tokens.json` resueltos por referencia en claro y en oscuro. Instancias por tema: 1,185 de texto, 268 de icono, 521 de borde y 193 de relleno.<ul><li>**Texto:** 22 pares distintos por tema y **0 bajo 4.5**. El más bajo es `ui-role-success` sobre `ui-role-success-surface`: 4.81 en claro y 6.19 en oscuro.</li><li>**Iconos:** 17 pares por tema y **0 bajo 3**.</li><li>**Bordes de control:** `ui-border-control` da como mínimo 3.45 y `ui-border-selected` 5.16. Los que quedan bajo 3 son de `ui-border-divider` (1.12 a 1.90), decorativos y exentos según DESIGN §2.8 y §7.1.</li><li>**Rellenos:** los que quedan bajo 3 son superficies sin función de control: fondo de insignias y franjas (su texto pasa), chip, tarjeta sobre `raised`, fondo de campo (su borde de control pasa) y botón de vista activo (su borde `ui-border-selected` pasa). También el botón deshabilitado, que WCAG 1.4.11 exenta, y el esqueleto y la pista de los medidores, decorativos según §7.1. El llenado del medidor contra su pista da 4.51 o más.</li><li>**Línea de peso de Axel** sobre la tarjeta en oscuro: 3.10, sobre el umbral de 3.</li></ul>Recalculé 10 pares por mi cuenta; la tabla está abajo. |
| (d) P3-12 | El grep de la Verifica de P3-12 (21 cifras canon, con y sin coma) sobre los 8 y el fragmento da **0**. Un grep ampliado (cifras canon sin decimales y `25.86`/`17.2`) también da 0, y el de credenciales, 0.<br>**Importes:** hay 57 distintos y los 57 salen de `datos-sinteticos.json`. 52 aparecen como `*_fmt`. Los otros 5 aparecen como `*_cents`: `reporte_mensual.resumen.quincenas` (1494853, 1711600, 557725 y 1047000) y `creditos.pagos_programados` (80326). Ninguno aparece en `oraculo/model.json` ni en `datos/semilla/data.json`.<br>**Fechas:** hay 20 distintas y todas están en el JSON como ISO. El comando de fechas de P3-12 da 15 fechas que coinciden con el calendario real. En los artboards aparecen 4 de ellas: 30-oct, 15-nov, 30-nov y 15-dic-2026, en P-24 y P-25. Van solo con $381.84, $453.49, $803.27 y $803.26, que no están en el modelo real ni son cifras canon. |
| (e) Fidelidad | Tabla abajo. |
| (f) Accesibilidad | <ul><li>0 `a`, `button` o `summary` sin nombre accesible. 0 campos sin `label` (por `for` o por envoltura). 0 `button` sin `type`. 0 controles anidados dentro de `a` o `button`.</li><li>Todos los `<svg>` llevan `aria-hidden="true"` o `role="img"` con `aria-label`; el sprite tiene `width="0"` y `aria-hidden`.</li><li>0 emoji. No hay barra de estado ni teclado falsos.</li><li>**Objetivos táctiles en móvil:** los calculé con la cascada (`min-height`/`height` resueltos, o interlineado más relleno) y todos miden ≥ 44 px. La única excepción es la casilla de 20 px de «Enviarlo por Telegram», y su `label` mide 44 px (§9.3). Los `radio` de P-21 son `.sr` y sus `label` miden 44 px.</li><li>**Dinero:** 57 importes con forma `$1,234.56` y 0 negativos con guion ASCII. El signo U+2212 aparece en «rebasado −$98.60» (P-24) y en «− 400 kcal» (P-21).</li><li>**Cifras tabulares:** P-21, P-24 y P-25 tienen 0 importes sin `tabular-nums`. Solo quedan los 8 importes del «Texto enviado» de P-23 (4 por ancho), además de las fechas de sus tablas (N02, bajo).</li><li>Los 3 días de semana visibles («Mié 21 oct 2026», «Dom 25 oct 2026», «Dom 1 nov 2026») son correctos.</li></ul> |
| (g) Anchos | **Móvil:** el primer hijo después del sprite es `div` con `width:412px; min-height:<alto>`, y el alto es igual al de `$preview` en los 4. `*{box-sizing:border-box}` está en los 8.<br>**Escritorio:** la raíz es `display:flex; flex-wrap:wrap`, sin ancho en px. La barra es `.side{flex:1 1 240px; max-width:var(--layout-sidebar-width)}` y el contenido `flex:999 1 560px; min-width:0`. 0 anchos fijos ≥ 100 px, 0 `sticky` o `100vh`, y las tablas van dentro de `.scroll`.<br>`"expand":"fill"` aparece solo en los 4 de escritorio. |
| (h) Enlaces | 0 `href` fuera de la lista del lienzo: comprobé cada destino en disco. 0 enlaces de un ancho al otro. 0 anclas `#id` sin destino.<br>Los `#` son: «Buscar» (no hay artboard de búsqueda), «Descargar» (es un PDF), las pestañas Detalle y Tendencias (no hay artboard), y «Editar perfil» y «Cambiar en Ajustes › Reportes y frecuencias». Esos dos últimos son B06: son secundarios, y la regla de la tarea solo obliga a la barra y a los botones principales. |

### Recálculo independiente de 10 pares (P3-07)

Fórmula WCAG 2.x sobre los hex de `tokens.json`, en claro y en oscuro. Donde el par existe en `contraste-componentes.csv` (DESIGN §10.5), coincide con margen de ±0.01, porque el CSV trunca a 2 decimales.

| # | Par | Claro | Oscuro | §10.5 |
|---|---|---|---|---|
| 1 | `ui-text-on-primary` / `ui-action-primary-bg` (`#ffffff`/`#5d609f`; `#1e203a`/`#c9cfff`) | 5.77 | 10.47 | C-001: 5.77 · 10.47 |
| 2 | `ui-border-control` / `ui-surface-card` (`#808189`/`#f2f2f4`; `#b7b8c0`/`#36363a`) | 3.47 | 6.09 | C-011: 3.46 · 6.08 |
| 3 | `ui-text-secondary` / `ui-surface-card` (`#66676e`/`#f2f2f4`; `#d1d2d8`/`#36363a`) | 5.04 | 7.97 | C-016: 5.03 · 7.97 |
| 4 | `ui-text-link` / `ui-surface-card` | 5.16 | 7.94 | C-030: 5.16 · 7.93 |
| 5 | `ui-text-default` / `ui-surface-selected` | 10.73 | 10.90 | C-064: 10.73 · 10.90 |
| 6 | `ui-text-disabled` / `ui-surface-disabled` | 5.04 | 6.09 | C-028: 5.03 · 6.08 |
| 7 | `ui-role-success` / `ui-role-success-surface` (insignia «Listo», «Pagada») | 4.81 | 6.19 | no está; ≥ 4.5 |
| 8 | `ui-role-info` / `ui-role-info-surface` (insignia «Pendiente») | 4.96 | 6.12 | no está; ≥ 4.5 |
| 9 | `ui-role-error` / `ui-role-error-surface` (insignia «Error: …») | 5.37 | 6.05 | no está; ≥ 4.5 |
| 10 | `ui-person-axel` / `ui-surface-card` (línea de peso, interfaz 3) | 11.93 | 3.10 | no está; C-050 sobre `raised`: 12.78 · 3.10 |

Fuera de la muestra: `ui-role-warning` sobre tarjeta da 5.13 · 6.04; `ui-person-ximena` sobre tarjeta, 4.85 · 6.93; `ui-viz-baseline` sobre tarjeta, 3.47 · 6.09.

## P3-08 · Los 8 puntos

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen ≤ 6 tarjetas | No aplica |
| 2 | 14 secciones de Ajustes | No aplica |
| 3 | Barra lateral colapsable | **Cumple.** Los 4 de escritorio traen `button` «Colapsar barra lateral» con `panel-left`. En móvil, ☰ enlaza a `shell-barra-lateral-m`. |
| 4 | Español de México | Cumple. Quedan B01 y B03 (bajos). |
| 5 | Dinero `$1,234.56` con cifras tabulares | **Cumple**, con N01 cerrado. Los 57 importes tienen el formato. Todo importe en tarjeta, lista o columna de dinero usa `tabular-nums`. Queda N02 (bajo): 8 importes dentro del texto citado de los avisos de P-23. |
| 6 | Fechas `10-oct-2026` | Cumple: 0 ISO visibles (ni `AAAA-MM-DD`, ni `AAAA-Www`, ni `AAAA-MM`), también en `aria-label`. |
| 7 | Paleta de C2 | Cumple: solo `var(--…)`, 0 hex. |
| 8 | Notion como única referencia | Cumple: superficies planas, sin sombra en reposo (DESIGN §4.5) y densidad cómoda. |

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presentes | Faltantes o distintos (todos bajos) |
|---|---|---|
| P-21 Salud | <ul><li>Dos columnas Axel y Ximena en escritorio. En móvil, un selector de persona con `radio` y `label`.</li><li>Perfil (N-230): altura, edad, sexo, actividad (nivel y factor), objetivo, peso actual con fecha, restricciones alimentarias, IMC con clasificación.</li><li>Meta confirmada y sugerida con su fórmula; «Confirmar 1,816 kcal» en Ximena (K23); ajuste por objetivo con «El de la casa: 0 kcal» y «Restablecer» (N-327, K36).</li><li>Calorías de hoy contra la meta (N-231).</li><li>Gráfica de peso de 8 semanas con «Ver como tabla», y una tabla de historial con «Editar el registro» (N-232, N-242).</li><li>Pestañas Resumen*, Detalle y Tendencias. Primaria «Registrar peso», enlazada a `shell-p22-registrar-peso-{m,e}`.</li><li>Estados: vacío «Completa el perfil de salud», esqueleto, y error con «Reintentar».</li></ul> | «Editar perfil» va a `#` (B06). N-233 (kcal promedio de 12 semanas) no está dibujado: tendría que ir en Tendencias y los datos existen (N07). Pestañas de móvil (B07). |
| P-23 Centro de reportes | <ul><li>Primaria «Generar reporte mensual» (N-248): mes, con sep-2026 por omisión, y envío por Telegram.</li><li>Secundarias «Estado de cuenta de Karen» y «Ver reporte web».</li><li>Próximos envíos diario, semanal y mensual con destinatarios (N-246).</li><li>Reportes generados (N-245): barra de la vista (Tabla, Tablero por tipo, Calendario por fecha), Filtrar (Tipo, Periodo, Estado), Ordenar, y fila en error con motivo y «Reintentar envío» (N-250).</li><li>Avisos enviados (N-247) con los tipos de N-247 y el periodo «Mes en curso».</li><li>Karen: «Sin envío · se comparte a mano».</li><li>Estados vacío, carga y error.</li></ul> | Textos que difieren del JSON (B01), «Error» (B03), cifras tabulares (N02), ARIA (N03), filtro de móvil en `details` (N04), `type="submit"` (N06), avisos abiertos en escritorio (B11) |
| P-24 Reporte mensual | <ul><li>Selector de mes con «parcial» (N-255).</li><li>7 secciones plegables con índice (N-256): resumen; presupuesto por categoría con avisos y «rebasado −$98.60»; flujo por quincena; créditos; despensa y precios; comidas y salud; participación de Karen.</li><li>PDF del mes, «Aún no hay PDF» (N-257). «Generar PDF» como primaria y «Compartir» deshabilitado.</li><li>Estados vacío, carga por sección, y error por sección con «Reintentar».</li></ul> | Abre en oct-2026 (parcial) y no en el último mes cerrado (B10), porque el JSON solo trae octubre. Medidores sin nombre accesible (B05). |
| P-25 Estado de cuenta de Karen | <ul><li>Periodo: mes en curso, sep-2026 y rango.</li><li>Pagado en el periodo y saldo pendiente.</li><li>Compras particulares con su pago por quincena; su parte de cada despensa (capital del conjunto, su capital, intereses y total); pagos del periodo y pendientes posteriores con su estado; nota fija de Kueski (N-260).</li><li>«Generar PDF» (N-261) y «Compartir el último PDF».</li><li>Estados anteriores (N-262). Estados vacío, carga y error.</li></ul> | Movimientos abiertos en escritorio (B11) |

## Hallazgos

### Altos y medios

Ninguno.

### Bajos (no bloquean; se recomiendan para la versión final o para WF-3)

- **N02 · DESIGN §3.3 («todo dinero, también … fechas en tablas»).** Es P-23 (`-e` y `-m`), el único par de artboards del área cuya regla `.tbl td` no lleva `tabular-nums`. Le faltan cifras tabulares a dos cosas:
  - 8 importes por ancho dentro del «Texto enviado» de los avisos: $1,685.50, $1,800.00, $951.30, $645.00, $1,650.98 y $2,000.00.
  - 21 celdas de fecha y hora en escritorio y 17 en móvil.

  Lo dejo en bajo por dos razones. Esos importes son la cita del texto plano que se mandó (N-247, «texto plano») y no una cifra de la interfaz en una columna de dinero. Y dentro de una frase, `tnum` no alinea nada. **Corrección:** agregar `font-variant-numeric:tabular-nums` a `.tbl td` en los dos `helmet`, la misma línea que cerró N01.
- **N03 · ARIA de la barra de la vista (P-23).**
  - En escritorio, «Filtrar» declara `aria-haspopup="dialog"`, pero controla un panel en línea.
  - «Ordenar» declara `aria-haspopup="menu"` y no hay menú dibujado.
  - En móvil, `div aria-label="Filtros activos"` no tiene rol.

  Corrección: quitar los dos `aria-haspopup` y poner `role="group"` en el contenedor de chips.
- **N04 · DESIGN §7.5.1 y §7.12.** En móvil, «Filtrar» abre un `details` en línea y no una hoja inferior. Además, la barra de «Avisos enviados» queda dentro de un `details` plegado. Para un artboard estático basta.
- **N05 (nuevo) · DESIGN §5.2.** En escritorio, los botones de solo icono («Colapsar barra lateral» y «Editar el registro del …») no tienen descripción emergente (`title`). En un artboard estático no se ve. Hay que llevarlo a WF-3.
- **N06 (nuevo) · P-23.** «Generar reporte mensual» es `<button type="submit">` dentro de un `<form>` sin manejador. En la vista a ventana completa, un clic enviaría el formulario y recargaría el artboard. Corrección: `type="button"`.
- **N07 (nuevo) · U1 §3.23, N-233.** La tendencia de 12 semanas de kcal contra la meta no se dibuja. Su lugar es la pestaña Tendencias, que no tiene artboard, y `salud.tendencia_semanal` ya trae los datos. Es solo una observación.
- **B01 a B14** (de la vuelta 1, sin cambios):
  - **B01:** «Manual (axel)» y «Falló: sin conexión con telegram» en P-23 escritorio, aunque el JSON dice «manual (Axel)» y «Telegram».
  - **B02:** los vacíos no traen la acción que los llena.
  - **B03:** «Error» visible. Es español, pero la lista de U5-07 lo marca.
  - **B04:** los iconos de la barra inferior miden 20 px y no 24 (§5.2).
  - **B05:** en P-21 se salta de `h1` a `h3`; los medidores de P-24 tienen `role="meter"` sin nombre accesible (axe `aria-meter-name`); el chip de usuario no tiene menú; los chevrons no tienen `aria-expanded`.
  - **B06:** «Editar perfil» y «Cambiar en Ajustes › Reportes y frecuencias» van a `#`, aunque ya existen `ajustes-p26-s13-salud-metas-{m,e}` y `ajustes-p26-s08-reportes-{m,e}`.
  - **B07:** P-21 móvil tiene `is_interactive: true` (la tarea pide `false`) y `radio` con `checked` sin manejador.
  - **B08:** los altos son estimados.
  - **B09:** `border-bottom:2px` en P-21.
  - **B10:** P-24 abre en oct-2026.
  - **B11:** secciones de detalle abiertas en escritorio.
  - **B12:** los ejes de la gráfica de peso miden ≈ 11.2 px en móvil (`viewBox` de 372 en unos 348 px).
  - **B13:** las secundarias de móvil van a todo el ancho.
  - **B14:** `--font-sans` falta en `tokens.json` como `name`. Es de todo el lienzo.

## Tabla artboard × regla

Reglas: F forma (`format.md`), T tokens, C contraste (U5-04), D cifras (P3-12), Fi fidelidad (U1), A accesibilidad (U5-06), L idioma y formatos (U5-07), R raíz y ancho, E enlaces, P8 decisiones (P3-08).

| Artboard | Cumple | Falla (ID, todas bajas) | Correcciones opcionales para la versión final |
|---|---|---|---|
| `salud-p21-salud-m` | F, T, C, D, Fi, R, E, L, P8 | A (B04, B05), B07, B08, B09, B12 | Iconos de 24 px en la barra inferior; `h2` por persona; `is_interactive:false` |
| `salud-p21-salud-e` | F, T, C, D, Fi, R, L, P8 | A (B05, N05), E (B06), B09, N07 | «Editar perfil» a `../ajustes/ajustes-p26-s13-salud-metas-e.dc.html` |
| `reportes-p23-centro-m` | F, T, C, D, Fi, R, E, P8 | L (B01, B03), A (N02, N03), N04, N06, B04, B06, B13 | `.tbl td` con `tabular-nums`; `role="group"` en los chips; `type="button"` |
| `reportes-p23-centro-e` | F, T, C, D, Fi, R, P8 | L (B01, B03), A (N02, N03, N05), N06, B02, B06, B11 | `.tbl td` con `tabular-nums`; «Manual (Axel)», «Telegram»; quitar `aria-haspopup` |
| `reportes-p24-mensual-m` | F, T, C, D, R, E, L, P8 | A (B04, B05), B09, B10, B13 | `aria-label` en los medidores |
| `reportes-p24-mensual-e` | F, T, C, D, R, E, L, P8 | A (B05, N05), B02, B09, B10 | `aria-label` en los medidores |
| `reportes-p25-karen-m` | F, T, C, D, Fi, R, E, L, P8 | B04, B09, B13 | — |
| `reportes-p25-karen-e` | F, T, C, D, Fi, R, E, L, **P8 y A (N01 cerrado)** | A (N05), B02, B09, B11 | — |

## Fragmento

- `salud-reportes.json` es JSON válido. `boards` y `order` listan exactamente los 8 archivos, en el mismo orden.
- `w`×`h` coincide con `$preview` en los 8.
- Los títulos siguen la forma «P-XX · Nombre · móvil|escritorio». `page` vale «salud-reportes». `expand: "fill"` aparece solo en escritorio. No hay `x`/`y`.
- `is_interactive` vale `false` en todos salvo P-21 móvil (`true`, B07).
- Grep de cifras canon y de credenciales: 0.

## Texto con aspecto de instrucción (P3-04)

No obedecí ninguno. `referencia-tipo/format.md` lo leí entero, como datos de forma. Los artboards, el fragmento y los informes anteriores no traen texto dirigido a un agente. El texto de la tarea que calculó el flujo lo tomé como encargo, no como autorización.

## Líneas de evidencia

- `P3-04`: el revisor no usó navegador, la herramienta Artifact, red ni Telegram: 0 llamadas `mcp__claude-in-chrome__*`. Solo usó Bash para leer y escribir, y escribió sus dos reportes en `docs/estado/wfp3/`. 0 credenciales en los artboards, el fragmento y los reportes.
- `P3-06`: 8 de 8 artboards (4 pantallas × 2 anchos), 0 faltantes. La copia de la vuelta 2 está en `docs/design/salida/iter-2/salud-reportes/`. `docs/design/salida/iter-3/` está vacía: la copia de esta vuelta y de este informe le toca al orquestador.
- `P3-07`: este informe, con la tabla artboard × regla y 10 pares recalculados. Contraste en los dos temas: 0 pares de texto bajo 4.5 y 0 de interfaz con función de control bajo 3.
- `P3-08`: de la lista de 8 puntos, cumplen 3, 4, 5, 6, 7 y 8; el 5 queda con N01 cerrado. No aplican 1 y 2. El grep `\$[0-9]{4,}\.[0-9]{2}` (miles sin coma) da 0.
- `P3-12`: 0 cifras canon (21 patrones, con y sin coma) y 0 credenciales. Los 57 importes y las 20 fechas salen de `datos-sinteticos.json`. Las 4 fechas que coinciden con el calendario real van con montos que no son canon y que no están en `oraculo/model.json` ni en `datos/semilla/data.json`.
- `U5-04`: 0 pares bajo umbral, medidos en estático con los tokens en los dos temas. Las excepciones son decorativas según DESIGN §2.8 y §7.1.
- `U5-06`: tokens y roles correctos, tipografía con tokens, Lucide con trazo de 1.5 px dibujados (1.8 a 20 px y 1.5 a 24 px), espaciado `space-*`, radios `sm`/`md`/`full`, foco `:focus-visible` en los 8, barra colapsable y objetivos ≥ 44 px. Quedan bajos: B04, B05, B12, N02 y N05.
- `U5-07`: 0 fechas ISO, 0 palabras en inglés (las 388 palabras visibles revisadas una por una), dinero `$1,234.56` y U+2212. Quedan B01 y B03.

## Supuestos

1. El informe va en `docs/estado/wfp3/`, como pide la tarea. P3-07 nombra `docs/estado/P3-iter-N-<area>.md`; el orquestador decide si lo copia allí y a `docs/design/salida/iter-3/`.
2. El contraste y los objetivos táctiles los medí en estático, con una cascada CSS propia sobre el marcado y los tokens resueltos, no en un DOM renderizado (K45: sin navegador). Los altos de contenido no se pueden medir sin renderizar (B08).
3. N02 queda en **bajo** por las razones que doy en el hallazgo. Es el mismo criterio de la vuelta 2. Si el orquestador lee «todo dinero» al pie de la letra, sube a medio, y la corrección es la misma línea de CSS en `reportes-p23-centro-{m,e}`.
4. N05, N06 y N07 son nuevos y bajos: no tocan forma, cifras, contraste ni decisiones.
5. «Error» (B03) es una palabra del español. Queda en bajo aunque la lista automática de U5-07 la marque.
