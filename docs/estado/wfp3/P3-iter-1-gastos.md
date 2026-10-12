# P3 · Vuelta 1 · Área «gastos» · Informe de revisión (P3-07)

Revisor: `revisor-gastos-r1` (independiente, no corrige). Fecha: 12-oct-2026.
Objeto: los 6 artboards de `docs/design/maquetas/lienzo/project/gastos/` (P-02 Todos los gastos, P-04 Gasto, P-05 Capturar ticket; móvil 412 y escritorio 1440) y el fragmento `docs/design/maquetas/lienzo/fragmentos/gastos.json`.
Contra: P3 1.3 (P3-07, P3-08, P3-12), U5 (U5-04, U5-06, U5-07, lo verificable en el markup), DESIGN.md (§3.3, §7.2, §7.5, §7.14, §9.2, §9.3), U1 §2.1, §2.2, §2.4, §2.6 y §3.4 a §3.7, `referencia-tipo/format.md` y `datos-sinteticos.json` 1.2.
Método: scripts propios con `python3 -I` en `datos/tmp/rev-gastos-r1/` (forma, tokens, contraste con cascada y especificidad simplificadas, objetivos táctiles, datos, cifras tabulares, fechas); borrados al terminar. Sin navegador ni red.

**Veredicto: NO APROBADO** (1 alta, 6 medias, 5 bajas).

## Resumen de hallazgos

| # | Sev. | Criterio | Artboards | Hallazgo |
|---|---|---|---|---|
| H1 | alta | U5-04 / P3-07 (contraste) | p02-e, p04-e | Los placeholders («Buscar gastos», «Escribe una nota…») no tienen color propio: el navegador pinta su color por defecto. Chromium usa `darkGray` #a9a9a9, que da 2.25:1 sobre `ui-surface-field` (#fafafb) en claro; Firefox usa el texto al 54 %, que da 3.05:1. Los dos quedan bajo 4.5:1. En P-02 el placeholder es la única etiqueta visible del campo. |
| H2 | media | U5-07 (fechas), fidelidad de datos | p04-m, p04-e | Dice «Gasto del **Mié** 12 oct 2026», pero el 12-oct-2026 es **lunes**. El texto está escrito a mano en el markup: no sale del formateador `fd()`, que sí calcula bien el día. |
| H3 | media | U1 §2.1 P-04 «nota (bloques)», N-060, N-062 | p04-m | El móvil no tiene el bloque **Nota**. En escritorio sí está (etiqueta y área de texto). Tampoco aparece dentro de «Mostrar 7 más». |
| H4 | media | U5-06 (foco visible), DESIGN §9.2 | p02-e | El campo de búsqueda lleva `outline:0` en línea. Ese estilo anula la regla `:focus-visible` del helmet, y la etiqueta que lo envuelve no tiene regla `:focus-within`. Resultado: el campo no muestra foco. |
| H5 | media | DESIGN §9.3 (≥ 44 px), U5-06 | p04-m | Los enlaces «Compra a crédito» (`a.small.end`) y «ticket-walmart-2026-10-12.jpg» (`#ticket`) miden lo que mide su línea de texto (unos 16 a 20 px). La fila `.val` sí mide 44 px, pero el objetivo táctil es el `<a>`, no la fila. |
| H6 | media | DESIGN §3.3 «cifras tabulares en **todo** dinero», P3-08, U5-06 | p02-e, p02-m, p04-e, p04-m, p05-e | Hay dinero sin `tabular-nums` en 10 lugares; la lista va en Correcciones. Las tablas, tarjetas y totales principales sí cumplen. |
| H7 | media | DESIGN §7.5.3 (anatomía del tablero) | p02-e, p02-m | Las columnas del tablero no tienen «+ Nuevo» al pie. |
| H8 | baja | DESIGN §7.5.4, U5-07 | p02-e | El día de hoy se marca solo con estilo (`class hoy`), sin `aria-current="date"`. El texto «Hoy es miércoles 21 de octubre de 2026» no sigue el formato «Mié 21 oct 2026». |
| H9 | baja | U5-07 (ninguna fecha ISO visible) | p04-m, p04-e, p05-m, p05-e | Se ven fechas ISO dentro de nombres de archivo («ticket-walmart-2026-10-12.jpg», «foto-ticket-walmart-2026-10-21.jpg»). Los nombres vienen de `datos-sinteticos.json`. |
| H10 | baja | DESIGN §7.2 (primario: «uno solo visible») | p04-e | «Ver ticket» aparece dos veces como primario: en el encabezado y en la sección «Ticket adjunto». |
| H11 | baja | DESIGN §3–§4 (escala de tokens) | todos | Hay medidas y tipografía en px literales: `font-size:12px` en `.badge` y `.bi`, `gap:2px`, `min-height:64px/88px/40px/48px`, `width:140px/72px`. |
| H12 | baja | Encargo (plantilla del fragmento) | fragmento | `is_interactive: true` en los 6 tableros; la plantilla del encargo decía `false`. `format.md` lo respalda porque hay manejadores y navegación, pero lo decide el orquestador. Los `h` son estimados: la raíz móvil usa `min-height` = `h`, así que un contenido más alto quedaría recortado en el lienzo. |

## Verificaciones (a)–(h)

| Verificación | Resultado | Detalle |
|---|---|---|
| (a) Forma | Cumple, 6/6 | doctype; `lang="es-MX"`; `<title>` propio («Todos los gastos», «Gasto», «Capturar ticket»); `./support.js`; link a `../ds/finanzas-hogar/tokens.css`; `<x-dc><helmet><style>`; regla de `body` con fondo, color y fuente por token; `a` y `a:hover`; `data-dc-script` con comilla simple, JSON válido y `$preview` igual a w/h del fragmento (412×2870, 1440×2200, 412×1960, 1440×2200, 412×1790, 1440×1600); `class Component extends DCLogic` con `renderVals()`. `html.parser` sin etiquetas sin cerrar ni cruzadas. Nombres válidos y stems únicos. |
| (b) Tokens | Cumple | 0 hex literales fuera de `data-props`. 0 `rgb()`/`hsl()`. Variables por archivo: 62, 48, 59, 53, 57 y 56; todas existen en `tokens.json` y en `tokens.css`. |
| (c) Contraste | Falla por H1 | 25 pares de texto y 25 de UI medidos con la cascada del helmet, todos sobre umbral en los dos temas; el mínimo es 4.51:1 (`ui-text-secondary` o `ui-icon-default` sobre `ui-surface-hover`, en claro). El placeholder sin token queda bajo umbral (H1). |
| (d) P3-12 | Cumple | El grep de las 21 cifras canon da 0. El de credenciales da 0. El de `$NNNN.NN` sin separador (P3-08) da 0. Los montos literales que no están en el JSON son 0, de 8+7+41+40+28+28. Las 33 filas g01–g33 de P-02 coinciden campo por campo con `gastos[]`. |
| (e) Fidelidad | Falla por H3 y H7 | Ver la tabla de abajo. |
| (f) Accesibilidad | Falla por H4 y H5 | Controles reales (`button`, `a`, `input`, `select` con `label for`, `textarea` con `label`). Iconos `aria-hidden` y botones de solo icono con `aria-label`. Sin emoji, sin barra de estado ni teclado falsos, sin texto en inglés. Negativo con U+2212 (0 casos de «-$» ASCII). |
| (g) Raíz | Cumple | Móvil: `width:412px; box-sizing:border-box`. Escritorio: fila `flex-wrap:wrap` sin ancho en px, con `main` en `flex:999 1 560px; min-width:0`, sin sticky ni 100vh; `expand:"fill"` solo en los 3 de escritorio. |
| (h) Enlaces | Cumple | Todos los href apuntan a archivos de la lista del lienzo del mismo ancho, a «#» o a `#ticket`, que es un ancla interna: `id="ticket"` existe en p04-m y p04-e. |

### Recálculo independiente de contraste (P3-07: 12 pares, hex de `tokens.json`)

| Par | Claro | Oscuro |
|---|---|---|
| ui-text-default / ui-surface-page | 11.53 (#36363a / #fafafb) | 14.07 (#f2f2f4 / #222325) |
| ui-text-secondary / ui-surface-card | 5.04 (#66676e / #f2f2f4) | 7.97 (#d1d2d8 / #36363a) |
| ui-text-secondary / ui-surface-hover | 4.51 (#66676e / #e5e6e9) | 5.50 (#d1d2d8 / #4d4e53) |
| ui-icon-default / ui-surface-hover | 4.51 | 5.50 |
| ui-role-expense / ui-surface-card | 5.33 (#a04b20 / #f2f2f4) | 5.90 (#f8a078 / #36363a) |
| ui-role-expense / ui-surface-selected | 5.32 (#a04b20 / #f0f1ff) | 5.97 (#f8a078 / #303258) |
| ui-role-warning / ui-role-warning-surface | 5.12 (#8b5d00 / #fff0dc) | 6.09 (#e5ae5b / #4c3000) |
| ui-role-info / ui-role-info-surface | 4.96 (#006fa4 / #e7f5ff) | 6.12 (#6bc3fd / #003b5a) |
| ui-text-on-primary / ui-action-primary-bg | 5.77 (#ffffff / #5d609f) | 10.47 (#1e203a / #c9cfff) |
| ui-text-link / ui-surface-card | 5.16 (#5d609f / #f2f2f4) | 7.94 (#c9cfff / #36363a) |
| ui-border-control / ui-surface-card (UI, 3:1) | 3.47 (#808189 / #f2f2f4) | 6.09 (#b7b8c0 / #36363a) |
| ui-text-default / ui-role-error-surface | 10.72 (#36363a / #ffeeed) | 11.20 (#f2f2f4 / #5a2124) |
| **placeholder del navegador / ui-surface-field** | **2.25** (Chromium #a9a9a9 / #fafafb); **3.05** (Firefox, 54 %) | 6.69 / 5.11 |
| placeholder propuesto `ui-text-secondary` / ui-surface-field | 5.40 | 10.43 |

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presente | Falta o difiere |
|---|---|---|
| P-02 escritorio | Subsecciones Resumen · Detalle* · Tendencias · Calendario; las 4 vistas (Tabla*, Tablero por categoría, Calendario lun–dom con detalle del día, Galería de tickets); Filtrar (persona, familia, categoría, tienda, forma de pago, tipo, ticket, periodo), Ordenar y Agrupar; búsqueda; tabla con las 10 columnas de N-040 y casilla; `aria-sort` en Fecha; fila de totales N-041 (29 gastos, $24,413.20, $21,078.00, $3,335.20); resumen N-042 (hogar, anterior, diferencia, categoría con más gasto); vista previa lateral; primaria «Registrar gasto» y «Foto de ticket»; tarjeta Estados con vacío (2 variantes), carga (esqueleto) y error (fila con Reintentar y edición fallida). | «+ Nuevo» por columna del tablero (H7). Los promedios diario, semanal y mensual de N-042 no están: el JSON no los trae y el autor lo declaró; queda para el dueño de los datos y no es hallazgo del área. |
| P-02 móvil | Filas compactas (descripción y monto; categoría y fecha debajo; persona); botón Filtrar; chips de periodo y orden; tablero de 1 columna con selector; agenda por día; galería de 2 columnas; totales; «+» en la barra inferior (Inicio · Comidas · + · Compras · Buscar). | «+ Nuevo» (H7). Los chips de filtro no tienen botón para quitarlos; es menor y no se reporta como hallazgo. |
| P-04 escritorio | Propiedades editables de N-060 en 2 columnas (con creado por, último cambio y «se calcula»); reparto hogar/Karen, devengo y efectivo; nota (bloques); ticket adjunto con máximo de 10 MB (K43); compra a crédito vinculada; Duplicar con «Deshacer»; Eliminar deshabilitado con su razón (N-063); «Marcar como pagado» explicado; 12 renglones con «Ajuste del ticket» y totales por cuenta (N-061); Estados con carga y error («Este gasto ya no existe» + Volver); vacío «No aplica», que concuerda con U1 («—»). | Día de la semana mal (H2); primario duplicado (H10). |
| P-04 móvil | 4 propiedades + «Mostrar 7 más»; reparto, devengo y efectivo; ticket adjunto; renglones plegados con su total; compra a crédito; acciones. | **Nota** (H3); día mal (H2); objetivos táctiles (H5). |
| P-05 escritorio | Pasos 1 a 4; Tomar foto o Elegir de la galería; foto visible a la izquierda; tabla de renglones a la derecha con texto leído, producto del catálogo, cantidad, precio, importe, cuenta y confianza; baja confianza desplegada con alternativa; forma de pago; tickets por confirmar con «Con borrador» (K3); Confirmar, Agregar renglón, Guardar para después, Cancelar y Repetir foto; Estados vacío, carga y error con los textos de U1 §2.2. | — |
| P-05 móvil | Pasos en página completa; renglones como tarjetas (alta plegados); Confirmar (48 px); Guardar para después; tickets por confirmar. | — |

## P3-08: lista de 8 puntos (área gastos, vuelta 1)

| Punto | Estado |
|---|---|
| Resumen ≤ 6 tarjetas | No aplica (área sin Inicio). |
| 14 secciones de Ajustes | No aplica. La barra lateral enlaza a `ajustes-p26-indice`. |
| Barra lateral colapsable | Cumple: botón «Colapsar barra lateral» y secciones plegables en los 3 de escritorio. |
| Español de México | Cumple: sin texto en inglés, con trato de tú. |
| Dinero `$1,234.56` tabular | Falla parcial (H6). El formato `$1,234.56` y el signo U+2212 se cumplen. |
| Fechas `10-oct-2026` | Falla parcial (H2, H8, H9). |
| Paleta C2 (índigo 280°) | Cumple: solo tokens C2 y 0 hex. |
| Notion como única referencia | Cumple: vistas de base de datos de Notion, propiedades y vista previa lateral. |

## Tabla artboard × regla

| Artboard | a Forma | b Tokens | c Contraste | d P3-12 | e Fidelidad | f Accesib. | g Raíz | h Enlaces | Fechas | Tabulares | Veredicto |
|---|---|---|---|---|---|---|---|---|---|---|---|
| gastos-p02-todos-m | ✓ | ✓ | ✓ | ✓ | H7 | ✓ | ✓ | ✓ | ✓ | H6 | falla (media) |
| gastos-p02-todos-e | ✓ | ✓ | **H1** | ✓ | H7 | H4 | ✓ | ✓ | H8 | H6 | falla (alta) |
| gastos-p04-gasto-m | ✓ | ✓ | ✓ | ✓ | H3 | H5 | ✓ | ✓ | H2, H9 | H6 | falla (media) |
| gastos-p04-gasto-e | ✓ | ✓ | **H1** | ✓ | H10 | ✓ | ✓ | ✓ | H2, H9 | H6 | falla (alta) |
| gastos-p05-ticket-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | H9 | ✓ | cumple (solo baja) |
| gastos-p05-ticket-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | H9 | H6 | falla (media) |

## Correcciones para la vuelta 2

1. **H1**, en p02-e y p04-e: añadir al helmet `::placeholder{color:var(--ui-text-secondary);opacity:1}`, como ya hacen los artboards de comidas. Da 5.40:1 en claro y 10.43:1 en oscuro.
2. **H2**, en p04-m y p04-e: cambiar «Gasto del Mié 12 oct 2026» por «Gasto del Lun 12 oct 2026». Mejor aún, sacarlo de `fd('2026-10-12', true)` en `renderVals()` para que no se repita el error.
3. **H3**, en p04-m: añadir el bloque «Nota» con `<label for>` y `<textarea>`, el mismo placeholder (ya con el token de H1) y área de 44 px como mínimo. Puede ir después de Reparto o dentro de «Mostrar N más»; en ese caso el contador pasa a «Mostrar 8 más».
4. **H4**, en p02-e: quitar `outline:0` del `<input type="search">` o añadir `label:focus-within{outline:var(--border-width-focus) solid var(--ui-border-focus);outline-offset:var(--focus-offset)}`.
5. **H5**, en p04-m: dar a los dos enlaces dentro de `.val` `display:inline-flex;align-items:center;min-height:var(--size-touch-min)`, o convertir la fila completa en el enlace.
6. **H6**: envolver en `<span class="num">` estos montos:
   - p02-e: «($25,126.40)» y «($2,858.70)».
   - p02-m: «($2,858.70)» y `{{ a.tot }}` en el encabezado de la agenda.
   - p04-e: «$2,929.00», «−$10.60» y los montos de «Totales por cuenta».
   - p04-m: «12 renglones · $2,918.40», `{{ r.unit }}` y «hogar $1,459.20 y Karen $1,459.20».
   - p05-e: `{{ L.suma }}` y `{{ L.ajuste }}` en la fila de suma.
7. **H7**, en p02-e y p02-m: botón de texto «+ Nuevo» al pie de cada columna del tablero, enlazado a `../shell/shell-p03-registrar-gasto-{m,e}.dc.html` (alto de 44 px en móvil).
8. **H8**, en p02-e: añadir `aria-current="{{ d.cur }}"` en la celda, con `cur:'date'` solo en hoy. Cambiar el texto a «Hoy es Mié 21 oct 2026».
9. **H9**: mostrar «Foto del ticket · 12-oct-2026» en lugar del nombre de archivo, o pedir al dueño de `datos-sinteticos.json` nombres sin ISO.
10. **H10**, en p04-e: dejar «Ver ticket» de la sección «Ticket adjunto» como secundario (`btn sc`).
11. **H11**: usar `.caption` en lugar de `font-size:12px`, y tokens `space-*`, `size-*` y `layout-*` en lugar de px literales donde exista el token.
12. **H12**, para el orquestador: decidir `is_interactive` y medir los `h` reales en el lienzo.

## Observaciones (no son hallazgos del área)

- DESIGN §7.5.2 dice «Cargar 200 más (200 registros por página, K6)», pero U1 N-040 dice «Páginas de 50». El autor siguió a U1 («páginas de 50»). Hay que unificar los documentos.
- La tabla se arma con `div` y roles ARIA (`table`, `row`, `columnheader`, `cell`) en lugar de `<table>`/`<th scope>`, porque el analizador de HTML saca un `<sc-for>` que esté dentro de `<table>`. Es un equivalente accesible válido y queda justificado en los supuestos del autor.
- «Gasto del hogar en oct-2026» ($20,525.78, devengado de `por_categoria`) y la «Parte del hogar» de la fila de totales ($21,078.00, `suma_hogar_cents`) son dos cifras distintas en la misma pantalla. Las dos salen del JSON y llevan etiqueta propia. Conviene que el texto de apoyo diga por qué difieren (devengo de la compra a crédito).

## P3-04 · Texto con aspecto de instrucción

Ninguno: no hay texto en los artboards dirigido a un agente.

## P3-12

El grep de las 21 cifras canon de `plan.md` §1.2 y K19 (con y sin separador) sobre `project/gastos/` y `fragmentos/gastos.json` da 0. El de credenciales da 0.

Para las fechas, el comando de P3-12 da 15 fechas que coinciden entre `oraculo/model.json`, `datos/semilla/data.json` y `datos-sinteticos.json`. En el área solo aparecen 30-oct-2026, 15-nov-2026, 30-nov-2026 y 15-dic-2026, en el calendario de efectivo de P-04. Van con montos sintéticos ($453.49 por parte, compra de $2,918.40, intereses de $709.52, tasa de 24.3120 %), ninguno del reporte ni de Kueski. Lo confirmé a mano.

Revisé a mano 3 artboards (p02-e, p04-e, p05-e): no hay correos, nombres de cuenta ni datos de facturación. Los nombres son Axel, Ximena y Karen (Q2-98).
