# P3 · Vuelta 2 · Área «compras» · Informe de revisión (P3-07)

Fecha: 12-oct-2026. Revisor independiente: `revisor-compras-r2` (vuelta 2 de 3). No escribió ni corrigió nada de lo revisado. Fuentes de verdad: DESIGN.md 1.0 y C2 1.0 (Q2-96); U1 §0.4, §1.2, §1.3, §2.1, §2.2, §3.15 a §3.18, §6.4 y F-3; criterios P3 1.3 (P3-07, P3-08 y P3-12); U5-04, U5-06 y U5-07 en lo que se ve en el código; `referencia-tipo/format.md` y `craft.md` (§ «Layout that survives direct manipulation», «Accessible as drawn» y «Mobile prototypes»), que se leyeron como datos de formato.

Objeto: los 8 artboards de `docs/design/maquetas/lienzo/project/compras/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/compras.json`. Reporte del autor: `docs/estado/wfp3/correccion-compras-r1.md`. Vuelta anterior: `docs/estado/wfp3/P3-iter-1-compras.md`.

**Dictamen: APROBADO.** Quedan 0 hallazgos altos y 0 medios. Los 2 altos y 5 medios de la vuelta 1 están cerrados. Quedan 16 bajos que no bloquean: 11 siguen abiertos desde la vuelta 1 (uno de ellos reformulado) y 5 son nuevos.

## 1. Hallazgos de la vuelta 1: estado

| Id r1 | Sev. | Hallazgo | Estado | Evidencia en la vuelta 2 |
|---|---|---|---|---|
| A-01 | alta | `@media (max-width:799px){` sin cerrar en los 4 `-e` | **CERRADO** | En los 8 `<style>` hay tantas `{` como `}`. El único `@media` de cada escritorio es `@media (max-width:799px){.side{max-width:none;flex-basis:100%}}` y solo contiene `.side`. `.side .it`, `.main`, `.bar-t`, `.c13` a `.c16 .row`, `.cell`, `.ax` y `.lb` ya valen en 1440. El eje de P-15 usa `fill: var(--ui-text-secondary)` en 12 px: 5.04:1 en claro y 7.97:1 en oscuro sobre `ui-surface-card`. |
| A-02 | alta | Fechas cortas sin año | **CERRADO** | Hay 0 fechas `d mmm` sin año y 0 rangos `19–25 oct`, tanto en el marcado como en `renderVals()`. Quedan 16 fechas distintas, todas `d-mmm-aaaa` (`2-oct-2026`, `18-sep-2026`, `13-oct-2026`…) o con día de la semana (`Lun 19 oct 2026`, `Dom 25 oct 2026`) en el rango de la semana, como pide U1 §0.4. Hay 0 fechas ISO visibles. |
| M-01 | media | Montos en frases sin cifras tabulares | **CERRADO** | El script `tabular.py` recorre todos los nodos de texto con `$` o con una variable cuyo valor lleva `$`: 51 nodos en total, y los 51 tienen un ancestro con clase `n` (`font-variant-numeric: tabular-nums`). Los `select` de P-14 llevan `class="n"`. |
| M-02 | media | «+» de la barra inferior de 40 × 56 | **CERRADO** | En los 4 `-m`, el `<a class="add">` hereda `min-width: 64px` y `min-height: 56px` de `.bot a` y no lleva `width`, `height` ni `min-width:0`. El círculo es `<span class="mas">` de 40 × 40 (`--size-bottombar-add`, `--radius-full`). Cada destino lleva `<span class="pil">` de 56 × 32. En el activo, la píldora es `ui-surface-selected` con icono `ui-icon-selected` (8.04:1 en oscuro y 5.15:1 en claro). |
| M-03 | media | Barra lateral sin «‹‹», tema ni chevrons | **CERRADO** | En los 4 `-e` está el botón «Colapsar barra lateral», con el icono `panel-left` (`rect` + `M9 3v18`). Hay 9 `<button aria-expanded>` con `chevron-right` (8 en false y 1 en true, en Compras y tiendas) y `aria-label` «Desplegar…» o «Plegar…». El chevron gira 90° (`.boton[aria-expanded=true] .i`). En el pie va el grupo «Tema» con 3 botones `aria-pressed`. La página activa lleva `aria-current="page"`. |
| M-04 | media | h1 de móvil con «· móvil» | **CERRADO** | Los h1 son «Lista de compras», «Comparar canasta», «Productos y precios» e «Inventario». «· móvil» queda solo en `<title>` y en el fragmento. |
| M-05 | media | Historial de P-15: marcadores iguales y sin etiqueta directa | **CERRADO** | Cada tienda lleva el marcador y el trazo de su lugar en el orden fijo (§8.6 y §8.7). Walmart: círculo y línea continua (serie 1). Bodega Aurrera: cuadrado y `8 4` (serie 2). Mercado Portales: rombo y `8 4 2 4` (serie 4). Tianguis: triángulo invertido y `4 4` (serie 5). Sam's no tiene datos en este producto. Las 4 etiquetas directas van en `ui-text-default`, no en el color de la serie (§8.8), con halo del color de la tarjeta. La leyenda repite los marcadores. Los puntos cuadran con `productos[jitomate].historial_precios`: Walmart $37.80, $35.00, $32.70 y $33.70; Tianguis $31.50, $29.70 y $31.30. |

Bajos de la vuelta 1: B-02 (eje de 11 px) está **CERRADO**, porque el eje y las etiquetas van en 12 px. B-08 se **reformula**: §4.5 dice que las tarjetas en reposo son planas (`shadow.0`) y que se separan por la superficie, así que no les falta borde en oscuro. Lo que sí falta es el tratamiento de la vista previa (B-08 abajo). B-01, B-03 a B-07 y B-09 a B-12 siguen abiertos.

## 2. Verificaciones (a)–(h)

Scripts propios con `python3 -I` en `datos/tmp/rev-compras-r2/`: `verif.py` (forma, tokens, fechas, dinero, enlaces y fragmento), `contraste.py`, `tabular.py`, `tactil.py` y `proc.py` (procedencia). Se borraron al terminar.

| | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a | Forma | CUMPLE | Las 8 piezas abren exacto con `<!doctype html><html lang="es-MX"><head><meta charset="utf-8"><title>…</title><script src="./support.js"></script><link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css"></head><body><x-dc><helmet><style>`. Todas tienen 1 `<x-dc>` y 1 `</x-dc>`, y cierran con `<script type="text/x-dc" data-dc-script data-props='…'>class Component extends DCLogic { renderVals() { return …` y `</script></body></html>`. `data-props` parsea en las 8 y `$preview` coincide con `w`×`h` del fragmento en las 8. `html.parser` da 0 etiquetas sin cerrar y 0 cierres inesperados. Las llaves CSS están balanceadas. Los 8 `<title>` son propios y no se repiten en el lienzo. Hay 90 stems en el lienzo, sin duplicados, y todos cumplen `[A-Za-z0-9.-]`. Los holes son lookups con punto (`{{r.c1.t}}`), sin expresiones. |
| b | Tokens | CUMPLE | Fuera de `data-props` hay 0 hex literales y 0 `rgb()`, `hsl()` u `oklch()`. Los SVG usan `currentColor` o `var(--ui-…)`. Las piezas usan entre 52 y 67 `var(--x)` distintos y todos existen en `tokens.json` (más `font-sans`, de `type.families`). |
| c | Contraste | CUMPLE | Se recalcularon 47 pares en claro y en oscuro (94 combinaciones) con los hex de `tokens.json`, resolviendo las referencias `{…}`. 0 quedan bajo umbral (tabla §3). Además, una matriz de todos los primeros planos contra todos los fondos de las 8 piezas no da ningún par real bajo umbral (nota en §3). |
| d | P3-12 | CUMPLE | El grep de las 21 cifras canon (con y sin coma) da 0, y en centavos también 0. Credenciales: 0. Montos de 4 cifras sin separador (`\$[0-9]{4,}\.[0-9]{2}`, P3-08): 0. Fechas del cruce `oraculo/model.json` ∪ `data.json` ∩ `datos-sinteticos.json`: hay 15 (del 30-oct-2026 al 15-jun-2027) y ninguna aparece en compras; la más tardía del área es el 21-oct-2026. Procedencia: `proc.py` hizo **741 comprobaciones con 0 discrepancias** (detalle en §4). |
| e | Fidelidad | CUMPLE | Están todos los elementos de U1 §2.1 y §2.2 y los 3 estados de cada pantalla en el escritorio (tabla §5). Los faltantes son de estilo y van como bajos. |
| f | Accesibilidad | CUMPLE | Los controles son `button`, `a`, `input` y `select` reales. Los campos tienen `label` (visible o `.sr`) y las casillas, `aria-label` interpolado. Los 231 SVG llevan `aria-hidden` o `role="img"` con `aria-label`, y los botones de solo icono, `aria-label`. El foco usa `:focus-visible` con `ui-border-focus`. Móvil: 55 controles medidos con `tactil.py`, todos de 44 px o más (`.btn`, `.ib`, `.tab`, `.mg`, `.chk`, campos y `select`) y los destinos de la barra inferior de 64 × 56. Hay 0 emoji, 0 barra de estado falsa y 0 teclado falso. 0 palabras de interfaz en inglés: «Stock por unidad» es el término de U1 §2.1, y las marcas y USDA están en la lista permitida. Dinero: todos los montos cumplen `$1,234.56` y los negativos llevan U+2212 (−$212.40, −$17.70, −1,500 ml). |
| g | Raíz | CUMPLE | Los móviles tienen `<div class="m" style="width:412px;min-height:H">` con `*{box-sizing:border-box}`, y H es 2,950, 2,250, 1,320 o 1,500 (todos de 915 o más). Los escritorios tienen la raíz `.app` sin ancho en px, con `flex-wrap`, barra lateral de 240 px que pasa a 100 % bajo 800 px, `.main` con `flex: 999 1 560px; min-width:0; max-width:1400px` y tablas en `.tbl` con `overflow-x:auto`. En el fragmento, los 4 `-e` llevan `"expand":"fill"` y los 4 `-m` no. |
| h | Enlaces | CUMPLE | Cada `href` es «#», un ancla interna (`#resumen`, `#detalle` o `#tendencias`) o un artboard de la lista del lienzo con el mismo ancho, y todos existen en disco. Los enlaces dentro del área usan la forma corta (`compras-p14-canasta-e.dc.html`), que format.md acepta como ruta relativa y resuelve al mismo archivo que `../compras/…`. F-3 cuadra: «Comparar tiendas» de P-13 lleva a P-14, y «Usar Bodega Aurrera» de P-14 lleva a P-13. «Simular pago con crédito» lleva a P-20, y «Terminar compra» y «Foto de ticket», a P-05. |

Fragmento: `boards` y `order` listan exactamente los 8 archivos, sin duplicados ni `x`/`y`. Todos llevan `page` «compras» e `is_interactive` false, con títulos `P-XX · Nombre · móvil|escritorio`.

## 3. Recálculo de contraste (U5-04 y P3-07)

Fórmula WCAG 2.x sobre los hex de `tokens.json`. Muestra de los 47 pares; los 23 restantes de la vuelta 1 dan los mismos valores.

| Par (primer plano / fondo) | Umbral | Claro | Oscuro | En §10.5 | Dónde |
|---|---|---|---|---|---|
| ui-text-default / ui-surface-page | 4.5 | 11.53 | 14.07 | sí | cuerpo y tabla |
| ui-text-default / ui-surface-card | 4.5 | 10.76 | 10.76 | sí | tarjetas, barra lateral y etiquetas directas `.lb` |
| ui-text-secondary / ui-surface-card | 4.5 | 5.04 | 7.97 | sí | `.cap.sec` y eje `.ax` del historial |
| ui-text-secondary / ui-surface-raised | 4.5 | 5.40 | 7.97 | sí | nombres de la barra inferior |
| ui-text-default / ui-surface-raised | 4.5 | 11.53 | 10.76 | sí | destino activo |
| ui-text-default / ui-surface-selected | 4.5 | 10.73 | 10.90 | sí | subpágina activa de la barra lateral |
| ui-text-on-primary / ui-action-primary-bg | 4.5 | 5.77 | 10.47 | sí | botón primario y `plus` del «+» |
| ui-text-on-primary / ui-action-primary-bg-hover | 4.5 | 8.49 | 12.66 | sí | primario en hover |
| ui-text-link / ui-role-info-surface | 4.5 | 5.19 | 7.81 | no (cumple) | «Planear» y «Capturar precios» en vacíos |
| ui-role-warning / ui-role-warning-surface | 4.5 | 5.12 | 6.09 | no (cumple) | insignias «Viejo» y «Reponer» |
| ui-role-success / ui-role-success-surface | 4.5 | 4.81 | 6.19 | no (cumple) | «La más barata», «Fresco» y «Suficiente» |
| ui-role-error / ui-role-error-surface | 4.5 | 5.37 | 6.05 | no (cumple) | «Agotado» |
| ui-text-default / ui-role-error-surface | 4.5 | 10.72 | 11.20 | no (cumple) | franja de error |
| ui-icon-default / ui-surface-card | 3 | 5.04 | 7.97 | sí | chevrons y botones de tema (nuevo en r2) |
| ui-icon-selected / ui-surface-selected | 3 | 5.15 | 8.04 | sí | tema pulsado y píldora activa (nuevo en r2) |
| ui-border-selected / ui-surface-card | 3 | 5.16 | 7.94 | no (cumple) | borde del tema pulsado (nuevo en r2) |
| ui-action-primary-bg / ui-surface-raised | 3 | 5.53 | 7.94 | sí | círculo «+» sobre la barra inferior |
| ui-border-control / ui-surface-card | 3 | 3.47 | 6.09 | sí | campos de «Ajustar existencia» |
| ui-border-focus / ui-surface-card | 3 | 3.56 | 6.03 | sí | anillo de foco en tarjeta |
| ui-viz-meter / ui-viz-track | 3 | 4.51 | 5.50 | sí | medidor de cobertura |
| ui-store-walmart / ui-surface-card | 3 | 5.11 | 4.76 | no (cumple) | línea, marcador y marca de Walmart |
| ui-store-bodega-aurrera / ui-surface-card | 3 | 8.40 | 6.79 | no (cumple) | ídem, Bodega Aurrera |
| ui-store-sams-club / ui-surface-card | 3 | 3.06 | 10.18 | no (cumple) | marca de Sam's |
| ui-store-market / ui-surface-card | 3 | 12.82 | 3.06 | no (cumple) | líneas y marcadores de los mercados |
| ui-person-axel / ui-surface-card | 3 | 11.93 | 3.10 | no (cumple) | anillo del avatar |

Resultado: **0 de 47 pares bajo umbral en los dos temas.** Nota de la matriz completa: el único cruce bajo 3:1 que puede darse es la marca de 10 px de los mercados sobre `ui-surface-hover`, en oscuro (2.11:1), solo mientras el puntero está sobre el chip «Tienda asignada» de P-13 e. La marca es `aria-hidden` y repite el nombre escrito al lado, así que cuenta como decorativa y no es hallazgo. Divisores, retícula y esqueletos están exentos por ser decorativos (§10.5).

## 4. P3-12: procedencia (muestra de 10 de 741)

| Valor en el artboard | Fuente en `datos-sinteticos.json` |
|---|---|
| Total de la lista $1,564.10 (P-13 m/e) | `lista_compras.total_cents` 156410 |
| Walmart $318.60, «cadena» (grupo de P-13) | `lista_compras.grupos[walmart].total_cents` 31860 |
| Detergente líquido $114.20 · «2-oct-2026 · hace 19 d» | `items[detergente].subtotal_cents` 11420, `precio_fecha` 2026-10-02, `meta.hoy` 2026-10-21 |
| Bodega Aurrera $1,551.70 y Walmart $1,638.30 | `canasta.por_tienda[].total_cents` 155170 y 163830 |
| +$86.60 y +$80.20 en idénticos | `diferencia_contra_la_mas_barata[walmart]` 8660 y 8020 |
| $650.00 · −$17.70 · −$212.40 · punto de equilibrio «—» | `membresia_sams`: 65000, −1770, −21240 y `null` |
| Yogur natural: Walmart $79.60 del 19-oct-2026; Sam's $109.30 sustituto; diferencia $29.70 | `canasta.detalle[yogur]` |
| Aguacate Hass en Tianguis: $146.80 del 18-sep-2026, viejo | `canasta.detalle[aguacate].precios[tianguis-narvarte]` |
| Jitomate saladet: $31.30 en Tianguis, 13-oct-2026, ticket; historial de 4 tiendas | `productos[jitomate].precios` e `historial_precios` |
| Inventario 1 · 5 · 8; Leche entera 1,500 ml, 3,500 ml / sem, 3.0 días | `inventario.resumen` y `articulos[leche]` |

`proc.py` revisó 741 valores: los 24 artículos de P-13 (subtotal, fecha y comprado), en móvil y en escritorio; las 120 celdas de P-14 (precio, fecha, viejo, sustituto y ganadora) y las 24 diferencias; los 10 totales y montos de membresía; la última entrada de los 14 artículos de P-16; y los 33 productos de P-15 e más los 12 de P-15 m. Dio 0 discrepancias. Revisión visual de p13-e, p14-e y p16-e: no hay correo, nombre de cuenta ni datos de facturación. Solo aparece «Axel» en el avatar, y Q2-98 lo permite.

## 5. Fidelidad: pantalla → elementos de U1

| Pantalla | Elemento (fuente) | m | e |
|---|---|---|---|
| P-13 | Grupos por vertical (cadenas por tienda; mercados con «visita el sábado» o «visita el martes») (§2.1, N-150) | ✓ grupos plegables con `aria-expanded` | ✓ filas de grupo con total |
| P-13 | Casillas de comprado, con 3 de 24 marcadas (N-151) | ✓ 44 × 44 | ✓ |
| P-13 | Totales por tienda y total general | ✓ $1,564.10 | ✓ |
| P-13 | Tienda asignada contra recomendada, motivo y precio con fecha | ✓ «Mejor en …» (motivo bajo demanda, B-06) | ✓ columnas Motivo, Tienda asignada y «Recomendada» |
| P-13 | Primaria «Terminar compra» (P-05); «Comparar tiendas», «Agregar artículo» y «Cambiar semana» (§2.2) | ✓ (secundarias fuera de «⋯», B-05) | ✓ |
| P-13 | Vistas tabla y tablero (§2.1) | n/a | ✓ |
| P-13 | Estados vacío (2 textos), carga y error con Reintentar | n/a | ✓ |
| P-14 | Selector de canasta con `label` (N-160) | ✓ | ✓ |
| P-14 | Totales por tienda con idénticos, sustitutos, sin precio y viejos; cobertura; diferencia en total y en idénticos (N-161) | ✓ | ✓ |
| P-14 | Ganadores por tienda y empates (N-163) | ✓ | ✓ |
| P-14 | Membresía de Sam's: costo, ahorro por canasta, compras al año, ahorro anual, punto de equilibrio y veredicto (N-164) | ✓ plegada (§6.4) | ✓ |
| P-14 | Detalle artículos × tiendas con precio, fecha, ganadora, sustituto y viejo (N-162) | pestaña Detalle (§6.4) | ✓ (sin fuente ni unidad base, B-09) |
| P-14 | Primaria «Usar Bodega Aurrera»; «Cambiar canasta»; «Simular pago con crédito» (P-20) | ✓ | ✓ |
| P-14 | Estados vacío (2 textos), carga y error | n/a | ✓ |
| P-15 | Precio por unidad base, marca de viejo, fuente e historial (N-170 y N-171) | ✓ filas compactas (§6.4) | ✓ tabla, vista previa e historial con marcadores |
| P-15 | Frescura (32 · 1 · 0) y fecha más reciente por tienda (N-172) | ✓ insignias | ✓ |
| P-15 | Primaria «Capturar precio»; «Foto de ticket», «Nuevo producto» y «Filtrar» | ✓ (`⋯` y Filtrar) | ✓ con filtros de tienda, categoría, frescura y fuente |
| P-15 | Vistas tabla, tablero y galería; pestañas Resumen, Detalle y Tendencias | ✓ pestañas | ✓ |
| P-15 | Estados vacío, carga y error | n/a | ✓ |
| P-16 | Suficiente, reponer o agotado; días restantes; movimientos (N-180 a N-182) | ✓ ordenado por estado | ✓ con vista previa y movimientos |
| P-16 | Primaria «Ajustar existencia» con motivo (conteo, consumo o merma); «Registrar sobrante»; «Filtrar» (N-183 y N-184) | ✓ | ✓ formulario con `label` |
| P-16 | Estados vacío, carga y error | n/a | ✓ |
| Global | Barra lateral: encabezado con avatar, Buscar y Registrar gasto, 9 secciones en orden fijo con ▸ propio, subpáginas, tema y «‹‹» (U1 §1.2 y DESIGN §7.9) | n/a (☰ «Abrir menú») | ✓ |
| Global | Barra inferior: Inicio · Comidas · «+» · Compras · Buscar (G-05 y §7.10) | ✓ | n/a |

Cobertura de P3-06 en el área: 4 pantallas × 2 anchos = 8 artboards. Están los 8 y el fragmento lista 8. Faltantes: 0.

## 6. P3-08: lista de 8 puntos (área compras, vuelta 2)

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen con ≤ 6 tarjetas | n/a (no es Inicio). Las filas de métricas del área tienen 4 tarjetas. |
| 2 | 14 secciones de Ajustes | n/a |
| 3 | Barra lateral colapsable | CUMPLE: «‹‹» `panel-left` con `aria-label`, ▸ con `aria-expanded` en las 9 secciones y tema en el pie |
| 4 | Español de México | CUMPLE: 0 palabras de interfaz en inglés y h1 sin etiquetas del lienzo |
| 5 | Dinero `$1,234.56` con cifras tabulares | CUMPLE: 0 montos mal formados, 0 de 4 cifras sin coma y 51 de 51 nodos con dinero en `.n` |
| 6 | Fechas `10-oct-2026` | CUMPLE: 0 fechas cortas y 0 ISO; día de la semana solo en el rango de la semana |
| 7 | Paleta de C2 (índigo 280°) | CUMPLE: solo `var(--ui-*)` y 0 hex |
| 8 | Notion como única referencia | CUMPLE: barra lateral con ▸, vistas de tabla, tablero y galería, franjas tipo callout, propiedades y vista previa lateral |

## 7. Tabla artboard × regla

Leyenda: ✓ cumple · ✗ falla (Id) · n/a no aplica. Entre paréntesis, los bajos que no bloquean.

| Artboard | a Forma | b Tokens | c Contraste | d P3-12 | e Fidelidad | f Accesibilidad | g Raíz | h Enlaces | Fechas | Tabulares | Veredicto |
|---|---|---|---|---|---|---|---|---|---|---|---|
| compras-p13-lista-m | ✓ (B-04) | ✓ | ✓ | ✓ | ✓ (B-05, B-06) | ✓ (B-01) | ✓ (B-03) | ✓ | ✓ | ✓ | aprobado |
| compras-p13-lista-e | ✓ (B-04) | ✓ | ✓ | ✓ | ✓ (B-07, B-13, B-15) | ✓ (B-01, B-14) | ✓ | ✓ | ✓ | ✓ | aprobado |
| compras-p14-canasta-m | ✓ | ✓ | ✓ | ✓ | ✓ (B-05) | ✓ (B-01) | ✓ (B-03) | ✓ | ✓ | ✓ | aprobado |
| compras-p14-canasta-e | ✓ (B-04) | ✓ | ✓ | ✓ | ✓ (B-09, B-11, B-13) | ✓ (B-01, B-14) | ✓ | ✓ | ✓ | ✓ | aprobado |
| compras-p15-productos-m | ✓ | ✓ | ✓ | ✓ | ✓ (B-06) | ✓ (B-01) | ✓ (B-03) | ✓ | ✓ | ✓ | aprobado |
| compras-p15-productos-e | ✓ (B-04) | ✓ | ✓ | ✓ | ✓ (B-08, B-16) | ✓ (B-01, B-10, B-14) | ✓ | ✓ | ✓ | ✓ | aprobado |
| compras-p16-inventario-m | ✓ (B-04) | ✓ | ✓ | ✓ | ✓ (B-06) | ✓ (B-01) | ✓ (B-03) | ✓ | ✓ | n/a (sin dinero) | aprobado |
| compras-p16-inventario-e | ✓ (B-04) | ✓ | ✓ | ✓ | ✓ (B-08, B-12) | ✓ (B-01, B-14) | ✓ | ✓ | ✓ | n/a (sin dinero) | aprobado |

## 8. Bajos para la vuelta 3 (no bloquean)

Abiertos desde la vuelta 1:
- **B-01 · Iconos de 16 px y trazo (§5.2).** Hay iconos de 16 px: 2 en p13-e, 1 en p13-m, 18 en p14-e, 8 en p14-m, 12 en p15-e, 2 en p15-m, 9 en p16-e y 3 en p16-m. Los de la leyenda del historial son marcadores, no iconos, y valen. `stroke-width: var(--icon-stroke)` (1.5 en unidades del viewBox) dibuja 1.25 px a 20 px. Para 20 px, usar `--icon-stroke-user-units-sm` (1.8).
- **B-03 · Altos de móvil sobrestimados (empeoró en P-13).** Con `.cta{margin-top:auto}`, el sobrante queda como franja vacía antes de «Terminar compra». La estimación propia de P-13 m da unos 2,430 a 2,500 px contra 2,950 declarados, un hueco de 450 a 500 px. La línea de fecha nueva casi no suma altura, porque las filas ya median 64 px. P-15 m da unos 1,210 contra 1,320 y P-16 m unos 1,380 contra 1,500. Hay que ajustar al renderizar.
- **B-04 · Atributos `hint-*` (format.md).** Hay `<sc-if>` sin `hint-placeholder-val`: 3 en p13-m, 3 en p13-e, 15 en p14-e, 2 en p15-e y 3 en p16-e. Hay 3 `<sc-for>` sin `hint-placeholder-count` en p16-m.
- **B-05 · Barra superior y pestañas de móvil (U1 §1.3 y DESIGN §7.10).** En P-13 y P-14, las secundarias van a la vista y no en «⋯». El título va en 16 px dentro de la barra y no en `title-1` debajo. Las pestañas no usan `role="tablist"` y la activa va en semibold, no en medium.
- **B-06 · Filas de móvil que no se pueden abrir (U1 §6.4).** P-13 m dice «Toca un artículo para ver por qué está en la lista», pero la fila no es botón ni enlace. Lo mismo pasa en P-15 m (historial) y P-16 m (movimientos).
- **B-07 · Estados vacíos (§7.14).** Van como franja informativa. §7.14 pide el icono de la entidad en 24 px, el texto en `ui-text-secondary` y el botón que lo llena.
- **B-08 (reformulado) · Vista previa lateral (§7.12 y §4.5).** El `<aside>` de P-15 e y P-16 e usa `.card` (`ui-surface-card`, sin sombra). La vista previa va en `ui-surface-raised` con `shadow-1` en claro y borde de 1 px en oscuro. Las tarjetas en reposo sí son planas (`shadow.0`), así que retiro esa parte de B-08 de la vuelta 1.
- **B-09 · Detalle de P-14 (N-162).** Faltan la fuente de cada precio y el precio por unidad base.
- **B-10 · `aria-selected` en `role="row"` dentro de `role="table"` (P-15 e).** Solo vale en `grid` o `treegrid`. Se puede usar `aria-current="true"` o el rol `grid`.
- **B-11 · Negativo sin rol (§3.3).** «−$17.70» y «−$212.40» de la lista de la membresía (P-14 e) no llevan color ni icono de rol.
- **B-12 · Dato incoherente, de los datos y no del artboard.** En Leche entera, la «Última entrada» es el 19-oct-2026 · sobrante, pero los movimientos son +6,000 ml el 14-oct y −1,500 ml el 19-oct, y la existencia es 1,500 ml. El artboard copia bien `datos-sinteticos.json`. Va para el responsable de los datos.

Nuevos en la vuelta 2:
- **B-13 · Escala tipográfica (§3.2, §7.10 y U5-06).** El h1 de escritorio usa `.t1` en 24/32 y 600; `title-1` es 32/40 y 700, y es lo que usan las demás áreas (31 artboards con `title-1`). Los títulos de sección usan `.t3` en 16/24; `title-3` es 20/28. El cuerpo y los botones de móvil van en 14 px; `body` es 16 px («botones en móvil»). La cifra principal de las tarjetas de métrica de escritorio va en 20 px; `title-2` es 24 px.
- **B-14 · Alturas de puntero fino sin `@media` (§4.6).** En los `-e`, `.e .btn`, `.e .fld input` (36 px), `.side .it` y `.boton` (32 px) se fijan sin `@media (pointer: fine) and (not (any-pointer: coarse))`. Cuando la página fluida se apila en un teléfono, quedan bajo 44 px. Los artboards `-m` no tienen este problema.
- **B-15 · Semántica de la barra lateral y migas (§7.9 y §7.10).** El `<nav>` usa `div` en lugar de lista anidada. La miga «Compras y tiendas» es un `p.cap`: no es enlace, no lleva `aria-current` y va en 12 px en lugar de `small`.
- **B-16 · Etiquetas directas del historial (P-15 e).** «Bodega» (y = 104.7) queda pegada a la línea de Tianguis, y «Portales» empieza en x = 241.6, junto al marcador de Walmart (borde en 240.9). Hay que confirmarlo al renderizar. La nota «círculo, cuadrado, triángulo, rombo y triángulo invertido» nombra el triángulo de Sam's, que no tiene serie en este producto.
- **B-17 · Esqueleto sobre página.** En la tarjeta «Estados», el esqueleto se pinta sobre `ui-surface-page` con `.sk` en `ui-skeleton-on-card`. Lo correcto es `ui-skeleton-on-page` (§7.14). Es decorativo.

## 9. P3-04 · Texto con aspecto de instrucción

No hay ninguno en los artboards, el fragmento ni el reporte del autor. `format.md` y `craft.md` se leyeron como datos de formato del publicador del tipo. Este revisor hizo 0 acciones de navegador y 0 llamadas a Artifact. Herramientas usadas: lectura de archivos y Bash con `python3 -I` y `grep` locales.

## 10. P3-12

- Cifras canon: 0, con y sin coma y en centavos.
- Credenciales: 0.
- Fechas coincidentes con el calendario real: hay 15 en el cruce y 0 aparecen en compras.
- Revisión visual de 3 artboards (p13-e, p14-e y p16-e): sin correo, cuenta ni facturación.
- Procedencia: 741 valores trazados a `datos-sinteticos.json`, con 0 discrepancias. **CUMPLE.**
