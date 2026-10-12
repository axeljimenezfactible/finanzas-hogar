# P3 · Vuelta 2 · Área «gastos» · Informe de revisión (P3-07)

Revisor: `revisor-gastos-r2` (independiente, vuelta 2 de 3, no corrige). Fecha: 12-oct-2026.
Objeto: los 6 artboards de `docs/design/maquetas/lienzo/project/gastos/` (P-02 Todos los gastos, P-04 Gasto y P-05 Capturar ticket, en móvil 412 y escritorio 1440) y el fragmento `docs/design/maquetas/lienzo/fragmentos/gastos.json`.
Criterios: P3 1.3 (P3-07, P3-08, P3-12); U5-04, U5-06 y U5-07, en lo que se puede verificar en el markup; DESIGN.md (§3.3, §7.2, §7.5.3, §7.5.4, §9.2, §9.3 y §10.5); U1 §2.1, §2.2, §2.4 y §6.2 (filas de P-02, P-04 y P-05); `referencia-tipo/format.md` (reglas de interactividad, de huecos y de PAGE); `datos-sinteticos.json` 1.2.
Método: escribí scripts propios y los corrí con `python3 -I` en `datos/tmp/rev-gastos-r2/`; ya están borrados. Fueron 4:
- `rev.py`: forma, tokens, P3-12, enlaces, raíz y días de la semana.
- `cas.py`: una cascada simplificada con especificidad, para medir contraste de texto y de UI en claro y oscuro, objetivos táctiles en móvil y dinero sin `tabular-nums`.
- `datos.py`: montos y fechas contra el JSON.
- `ingles.py`: busca texto en inglés.

Además comparé con `diff` contra la copia de la vuelta 1 (`docs/design/salida/iter-1/gastos/`). No usé navegador ni red.

**Veredicto: APROBADO.** Hay 0 hallazgos altos, 0 medios y 8 bajos. Los 7 hallazgos altos y medios de la vuelta 1 (H1 a H7) quedaron cerrados. Siguen abiertos 5 bajos de la vuelta 1 (H8 a H12) y se suman 3 bajos nuevos (N1 a N3).

## 1. Cierre de los hallazgos de la vuelta 1

| # | Sev. v1 | Hallazgo | Estado | Evidencia (vuelta 2) |
|---|---|---|---|---|
| H1 | alta | Placeholders sin color propio: el navegador los pinta a 2.25:1 en Chromium y 3.05:1 en Firefox | **Cerrado** | El helmet de p02-e, p04-e y p04-m trae `::placeholder{color:var(--ui-text-secondary);opacity:1}`. Calculé `ui-text-secondary` sobre `ui-surface-field` (el fondo del `label` en P-02 y de `.ta` en P-04): **5.40:1 en claro y 10.43:1 en oscuro**. Ningún otro `placeholder=` del área queda sin regla. |
| H2 | media | «Gasto del Mié 12 oct 2026», cuando el 12-oct-2026 es lunes | **Cerrado** | p04-m y p04-e dicen «Gasto del Lun 12 oct 2026». Un validador de día de la semana revisó todas las fechas con día del área, incluidas «Lun 19 oct 2026 20:05» y «Mié 21 oct 2026 18:42» de P-05: 0 errores. El arreglo `DS` de `fd()` empieza en «Dom», así que es correcto con `getDay()`. |
| H3 | media | P-04 móvil no tiene el bloque Nota | **Cerrado** | p04-m tiene la sección «Nota» después de Reparto y antes de Ticket adjunto, con `<label for="nota">` y `<textarea id="nota" class="ta">`, `min-height:var(--size-touch-min)` y el mismo placeholder que en escritorio. Está siempre visible, así que «Mostrar 7 más» sigue cuadrando. |
| H4 | media | `outline:0` en línea en la búsqueda de P-02 | **Cerrado** | `outline:0/none` aparece 0 veces en los 6 archivos. Ahora rige `:focus-visible{outline:var(--border-width-focus) solid var(--ui-border-focus);outline-offset:var(--focus-offset)}`. El anillo sobre `ui-surface-field` da 3.82:1 en claro y 7.89:1 en oscuro, sobre el umbral de 3:1. |
| H5 | media | Enlaces de P-04 móvil de 16 a 20 px de alto | **Cerrado** | «Compra a crédito» y «ticket-walmart-2026-10-12.jpg» llevan `display:inline-flex;align-items:center;min-height:var(--size-touch-min)`, que da 44 px. `cas.py` midió todos los `a`, `button`, `input`, `select` y `textarea` de los 3 artboards móviles. Solo `a.gal` no declara alto, y es la tarjeta de galería con portada 4:3, de mucho más de 44 px. |
| H6 | media | Dinero sin cifras tabulares en 10 lugares | **Cerrado** | Los 10 lugares de la lista de la vuelta 1 van ahora en `<span class="num">`, y `.num{font-variant-numeric:tabular-nums}` existe en los 6 helmets. `cas.py` encontró 0 textos con `$` o `−$` fuera de un ancestro `.num`. Revisé a mano las 73 interpolaciones `{{ }}` que quedan fuera de `.num` (`rep`, `tk`, `cnt`, `det`, `alt`, `altaTxt`, `masTxt` y `abTxt`, entre otras) y ninguna da dinero. |
| H7 | media | Tablero sin «+ Nuevo» | **Cerrado** | En p02-e hay `<a class="btn" href="../shell/shell-p03-registrar-gasto-e.dc.html">+ Nuevo</a>` dentro del `<sc-for>` de columnas, al pie de cada una. En p02-m está al pie de la columna visible, enlazado a `-m` y con `min-height:var(--size-touch-min)`. Es un botón de texto (`.btn` sin `pri`), como pide §7.5.3. |

## 2. Hallazgos vigentes (todos bajos)

| # | Sev. | Criterio | Artboards | Hallazgo |
|---|---|---|---|---|
| H8 | baja | DESIGN §7.5.4, U5-07 | p02-e | Sigue abierto. El día de hoy del calendario se marca solo con la clase `hoy`, sin `aria-current="date"`. El texto «Hoy es miércoles 21 de octubre de 2026» no sigue el formato corto «Mié 21 oct 2026». |
| H9 | baja | U5-07 (ninguna fecha ISO visible) | p04-m, p04-e, p05-m, p05-e | Sigue abierto. Hay fechas ISO dentro de nombres de archivo: «ticket-walmart-2026-10-12.jpg», «foto-ticket-walmart-2026-10-21.jpg» y «foto-ticket-bodega-2026-10-19.jpg». Los nombres salen de `datos-sinteticos.json`. |
| H10 | baja | DESIGN §7.2 (un solo primario visible) | p04-e | Sigue abierto. «Ver ticket» aparece 2 veces como `btn pri`. |
| H11 | baja | DESIGN §3–§4 (escala de tokens) | todos | Sigue abierto. Hay px literales: `font-size:12px` en `.badge` y `.bi`, `gap:2px`, `min-height:40px/36px/48px`, `width:140px/72px` y `.dn` con 24px. |
| H12 | baja | Encargo / `format.md` | fragmento | Sigue abierto y lo decide el orquestador. Los 6 tableros tienen `is_interactive: true`. La plantilla del encargo dice `false`, pero `format.md` pide `true` para artboards con manejadores y navegación, y estos los tienen. Los `h` son estimados. Para p04-m estimé unos 2,040 px de contenido en el estado inicial (plegado), menos que los 2,180 de la raíz. |
| N1 | baja (nuevo) | DESIGN §7.5.4 «Patrón ARIA *grid*» | p02-e | El calendario lleva `role="grid"`, pero sus hijos son `div` y `button` sin `role="row"` ni `gridcell`. Un DOM así daría `aria-required-children` en axe (U5-03, fuera del alcance estático). Hay dos arreglos: quitar `role="grid"` y dejar una lista de botones con nombre completo, o agrupar cada semana en `role="row"` con `role="gridcell"`. |
| N2 | baja (nuevo) | DESIGN §7.5.3, Accesibilidad | p02-e | El texto de ayuda dice «con teclado, usa “Mover a…” en el menú “⋯” de la tarjeta», pero las tarjetas no dibujan ese «⋯»: la tarjeta es un `<button>` y no admite otro botón dentro. La alternativa al arrastre sí existe, porque la Categoría se edita en la tabla y en P-04, y en Notion el «⋯» aparece al pasar el puntero. Aun así conviene dibujarlo en una tarjeta en hover o foco, como hermano del botón. |
| N3 | baja (nuevo) | DESIGN §7.5.4 (días de otro mes en `ui.text.disabled`) | p02-e | `.cb.out{color:var(--ui-text-disabled)}` no llega al número del día, porque `.dn` fija `color:var(--ui-text-secondary)`. Los días de sep y nov se ven igual que los de oct. Se arregla con `.cb.out .dn{color:var(--ui-text-disabled)}`. Es el par C-093 de §10.5, de texto con umbral 4.5, y lo cumple: 5.40:1 en claro y 7.96:1 en oscuro. |

## 3. Verificaciones (a)–(h)

| Verificación | Resultado | Detalle |
|---|---|---|
| (a) Forma | Cumple 6/6 | La cabecera es exacta: doctype, `lang="es-MX"`, charset, `<title>` propio («Todos los gastos», «Gasto», «Capturar ticket»), `./support.js`, `../ds/finanzas-hogar/tokens.css`, `<x-dc><helmet><style>`. Hay regla `body` con fondo, color y fuente por token, y reglas `a` y `a:hover`. El cierre `</x-dc>` va seguido de `<script type="text/x-dc" data-dc-script data-props='…'>` y `</body></html>`. El JSON de `data-props` es válido y `$preview` coincide con el fragmento: 412×2870, 1440×2200, 412×2180, 1440×2200, 412×1790 y 1440×1600. Todos tienen `class Component extends DCLogic` con `renderVals()`. `html.parser` con pila propia: 0 etiquetas sin cerrar, 0 cierres huérfanos, 0 `<sc-for>`/`<sc-if>` autocerrados. Los huecos `{{ }}` (99, 63, 14, 21, 68 y 66) son todos búsquedas con punto, sin expresiones. Todo nombre de primer nivel sale de `renderVals()`. Los nombres de archivo son válidos y los 90 stems del lienzo son únicos. |
| (b) Tokens | Cumple | Hex literales fuera de `data-props`: 0. `rgb()`, `hsl()`, `oklch()` y `color-mix()`: 0. Variables por archivo: 62, 49, 59, 55, 57 y 56. Todas existen en `tokens.css`, y en `tokens.json`, donde `--font-sans` sale de `type.families.sans`. |
| (c) Contraste | Cumple | Ver §4. Midió 27 pares de texto, con 0 bajo 4.5:1, y 17 pares de UI, con 0 bajo 3:1, en claro y en oscuro. |
| (d) P3-12 | Cumple | El grep de las 21 cifras canon, con y sin coma, da 0. Credenciales: 0. `$NNNN.NN` sin separador: 0. «-$» o «–$»: 0. De los 241 montos con `$` de los 6 artboards, 0 faltan en el JSON. De las fechas `d-mmm-aaaa`, 0 faltan en el JSON. La muestra está en §5. |
| (e) Fidelidad | Cumple | Ver §6. Los dos elementos que faltaban en la vuelta 1 (Nota en P-04 móvil y «+ Nuevo») ya están. |
| (f) Accesibilidad | Cumple, con bajos | Los controles son reales: `button`, `a`, `input` con `aria-label`, y `textarea` con `label for`. No hay controles sin nombre ni `svg` sin `aria-hidden` o sin nombre. No hay emoji ni barra de estado o teclado falsos, y no hay texto en inglés (`ingles.py`: 0, contando solo texto visible, `aria-label`, `placeholder` y cadenas de datos). Los negativos usan U+2212. Móvil ≥ 44 px (ver H5). El foco es visible (ver H4). Quedan los bajos N1 y N2. |
| (g) Raíz | Cumple | Móvil: `width:412px;…;box-sizing:border-box` en los 3. Escritorio: fila `display:flex;flex-wrap:wrap` sin ancho en px, y `expand:"fill"` solo en los 3 de escritorio. |
| (h) Enlaces | Cumple | Todos los `href` apuntan a «#», a un ancla interna que existe (`#ticket`) o a un archivo de la lista del lienzo del mismo ancho. Todos los destinos ya existen en disco. |

## 4. Contraste (P3-07, U5-04), con hex de `tokens.json` y alias resueltos

`cas.py` tomó el color de texto y el fondo efectivo de cada nodo de texto: reglas del helmet con especificidad, estilos en línea y ancestros. Lo hizo en reposo y con las variantes `on` y `sel`. En la columna «§10.5», «no» significa que el par no está listado, pero cumple el umbral por cálculo propio.

| Par (texto / fondo) | Claro | Oscuro | §10.5 |
|---|---|---|---|
| ui-text-secondary / ui-surface-hover (y ui-icon-default) | **4.51** | 5.50 | sí |
| ui-text-secondary / ui-surface-selected | 5.02 | 8.08 | sí |
| ui-text-secondary / ui-role-warning-surface | 5.03 | 8.05 | no |
| ui-text-secondary / ui-surface-card | 5.04 | 7.97 | sí |
| ui-text-disabled / ui-surface-disabled (Eliminar deshabilitado) | 5.04 | 6.09 | sí |
| ui-text-link / ui-surface-card | 5.16 | 7.94 | sí |
| ui-role-warning / ui-role-warning-surface | 5.12 | 6.09 | no |
| ui-role-info / ui-role-info-surface | 4.96 | 6.12 | no |
| ui-role-expense / card · page · raised · selected | 5.32–5.72 | 5.90–7.71 | no |
| ui-text-secondary / ui-surface-page · raised | 5.40 | 10.43 · 7.97 | sí |
| ui-text-on-primary / ui-action-primary-bg | 5.77 | 10.47 | sí |
| ui-text-default / (12 fondos, incluidas las superficies de rol) | 9.64–11.53 | 7.42–14.07 | sí / no |
| **Placeholder** ui-text-secondary / ui-surface-field (p02-e, p04-e, p04-m) | **5.40** | 10.43 | — |
| UI: ui-border-control / page · card · selected · info · error | 3.45–3.71 | 6.00–7.96 | sí / no |
| UI: ui-border-focus / 12 fondos usados | **3.19** (sobre hover) a 3.82 | 4.16–7.89 | sí / no |

Mínimos: texto 4.51:1 y UI 3.19:1. Pares bajo umbral: 0.

## 5. Muestra de datos (P3-12; los valores salen de `datos-sinteticos.json`)

| Valor en el artboard | Ruta en el JSON |
|---|---|
| $25,126.40 (sep-2026) | `gastos_resumen.anterior.hogar_fmt` · `presupuesto.tendencia[4].real_cents` |
| −$4,600.62 | `gastos_resumen.anterior.diferencia_fmt` |
| $20,525.78 | `presupuesto.resumen.real_cents` |
| $2,858.70 (venideros) | `gastos_resumen.venideros_cents` |
| $2,918.40 / $1,459.20 / $1,459.20 | `gastos[12]` (g13): `monto_cents`, `hogar_cents`, `karen_cents` |
| $2,929.00 y −$10.60 | `tickets.confirmados[0].suma_productos_fmt` y `.ajuste_fmt` |
| $453.49 | `gastos[12].efectivo[0].hogar_cents` |
| $906.98 | `presupuesto.despensa_devengada[0].hogar_mes_cents` |
| $709.52 / $3,627.92 | `creditos.kueski.compras[1].intereses_cents` · `presupuesto.despensa_devengada[0].total_con_intereses_cents` |
| 12-oct-2026 (Lun) | `gastos[12].fecha` = `tickets.confirmados[0].fecha` = 2026-10-12 |
| Hoy 21-oct-2026 (Mié) | `meta.hoy` = 2026-10-21 · `meta.dia_semana` = miércoles |

Fechas (P3-12, revisión manual): 4 fechas del área coinciden con `oraculo/model.json` o con `datos/semilla/data.json`: 30-oct-2026, 15-nov-2026, 30-nov-2026 y 15-dic-2026. Son las mismas de la vuelta 1, del calendario de efectivo de P-04, y van con montos sintéticos ($453.49 por parte). No hay ninguna cifra del reporte ni de Kueski.

## 6. Fidelidad: pantalla → elementos de U1 §2.1, §2.2, §2.4 y §6.2

| Pantalla | Presente | Falta o difiere |
|---|---|---|
| P-02 escritorio | Pestañas Resumen · Detalle* · Tendencias · Calendario. Vistas Tabla*, Tablero, Calendario y Galería de tickets. Filtrar, que nombra persona, familia, categoría, tienda, forma de pago, tipo, ticket y periodo. Ordenar, Agrupar y Buscar. Tabla con fila de totales, vista previa lateral y primario «Registrar gasto». **«+ Nuevo» al pie de cada columna**. Tarjeta Estados con vacío («Sin gastos en oct-2026» + Registrar gasto; «Ningún gasto coincide» + Quitar filtros), carga (esqueleto, «8 filas» en el texto) y error (Reintentar; edición fallida que vuelve al valor anterior). | — (bajos H8, N1, N2 y N3) |
| P-02 móvil | Filas compactas, Filtrar (hoja), tablero de 1 columna con selector y **«+ Nuevo»**, agenda por día, galería de 2 columnas y «+» en la barra inferior. | Los estados van solo en escritorio, como pide el encargo. |
| P-04 escritorio | Propiedades en 2 columnas; Reparto con devengo y efectivo; Nota; Ticket adjunto; Duplicar con Deshacer; Eliminar deshabilitado con su razón; «Abrir compra a crédito»; 12 renglones con ajuste y totales por cuenta; Estados (carga, y error «Este gasto ya no existe» + Volver). | — (bajos H9 y H10) |
| P-04 móvil | 4 propiedades + «Mostrar 7 más»; Reparto; **Nota**; Ticket adjunto; renglones plegados con su total; compra a crédito; acciones. | «Marcar como pagado» no aplica, porque g13 está realizado. |
| P-05 escritorio | Pasos 1 a 4; Tomar foto o Elegir de la galería; foto a la izquierda y tabla de renglones con confianza a la derecha; forma de pago; tickets por confirmar con «Con borrador» (K3); Confirmar ticket, Agregar renglón, Guardar para después, Cancelar y Repetir foto; Estados vacío, carga («Leyendo el ticket…») y error («No se pudo leer… captura el ticket a mano»). | — |
| P-05 móvil | Pasos en página completa, renglones como tarjetas con los de confianza alta plegados, Confirmar ticket de 48 px, Guardar para después y tickets por confirmar. | — |

## 7. P3-08: lista de 8 puntos (área gastos, vuelta 2)

| Punto | Estado |
|---|---|
| Resumen con ≤ 6 tarjetas | No aplica: el área no tiene Inicio. |
| 14 secciones de Ajustes | No aplica. La barra lateral enlaza a `ajustes-p26-indice-e`. |
| Barra lateral colapsable | Cumple: botón «Colapsar barra lateral» en los 3 artboards de escritorio. |
| Español de México | Cumple: 0 textos en inglés, con trato de tú. |
| Dinero `$1,234.56` tabular | **Cumple**: H6 cerrado, 0 montos fuera de `.num` y signo U+2212. |
| Fechas `10-oct-2026` | Cumple, con bajos H8 y H9. Los días de la semana son correctos (H2 cerrado). |
| Paleta C2 (índigo 280°) | Cumple: solo tokens de C2 y 0 hex. |
| Notion como única referencia | Cumple: vistas de base de datos, propiedades y vista previa lateral. |

## 8. Tabla artboard × regla

| Artboard | a | b | c | d | e | f | g | h | Fechas | Tabulares | Veredicto |
|---|---|---|---|---|---|---|---|---|---|---|---|
| gastos-p02-todos-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | cumple |
| gastos-p02-todos-e | ✓ | ✓ | ✓ | ✓ | ✓ | N1, N2 | ✓ | ✓ | H8 | ✓ | cumple (solo bajos) |
| gastos-p04-gasto-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | H9 | ✓ | cumple (solo bajos) |
| gastos-p04-gasto-e | ✓ | ✓ | ✓ | ✓ | H10 | ✓ | ✓ | ✓ | H9 | ✓ | cumple (solo bajos) |
| gastos-p05-ticket-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | H9 | ✓ | cumple (solo bajos) |
| gastos-p05-ticket-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | H9 | ✓ | cumple (solo bajos) |

H11 (px literales) aplica a los 6 y H12 al fragmento.

## 9. Correcciones sugeridas (opcionales: solo bajos)

1. **H8 y N1**, en p02-e: añadir `aria-current="{{ d.cur }}"` con `cur:'date'` solo en hoy, y cambiar el texto a «Hoy es Mié 21 oct 2026». Para el grid, quitar `role="grid"` o envolver cada semana en `role="row"` con celdas `role="gridcell"`.
2. **N2**, en p02-e: dibujar en una tarjeta (estado hover o foco) un botón «⋯» hermano del botón de la tarjeta, con `aria-label="Más acciones de la tarjeta"`. Otra opción es cambiar el texto de ayuda para que remita a la edición de Categoría.
3. **N3**, en p02-e: `.cb.out .dn{color:var(--ui-text-disabled)}`.
4. **H9**: mostrar «Foto del ticket · 12-oct-2026» o pedir nombres sin ISO al dueño de `datos-sinteticos.json`.
5. **H10**, en p04-e: dejar el «Ver ticket» de la sección Ticket adjunto como `btn sc`.
6. **H11**: usar `.caption` y los tokens `space-*`, `size-*` y `layout-*` donde existan.
7. **H12**, para el orquestador: decidir `is_interactive` (`format.md` respalda `true`) y medir los `h` reales en el lienzo.

## 10. Observaciones (no son hallazgos del área)

- Sigue la diferencia entre DESIGN §7.5.2 («200 por página», K6) y U1 N-040 («páginas de 50»). El área sigue a U1.
- La tabla usa `div` con roles ARIA (`table`, `row`, `columnheader`, `cell`) por la restricción de `<sc-for>` dentro de `<table>`. Es un equivalente válido.
- Siguen en la misma pantalla «Gasto del hogar en oct-2026» ($20,525.78, devengado) y «Parte del hogar» de la fila de totales ($21,078.00). Las dos salen del JSON y tienen su etiqueta.

## P3-04 · Texto con aspecto de instrucción

Ninguno. Los artboards no tienen texto dirigido a un agente. Leí `format.md` solo como dato de formato.

## P3-12

- El grep de las 21 cifras canon de `plan.md` §1.2 y K19, con y sin separador, sobre `project/gastos/` y `fragmentos/gastos.json` da 0. El de credenciales también da 0.
- Fechas: hay 4 coincidencias, todas de calendario quincenal y con montos sintéticos (§5). No vienen del reporte ni de Kueski.
- Los nombres que aparecen son Axel, Ximena y Karen (Q2-98). No hay correos, cuentas reales ni datos de facturación.
