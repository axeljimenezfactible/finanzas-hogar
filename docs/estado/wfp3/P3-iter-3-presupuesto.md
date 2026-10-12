# P3-iter-3 · Área «presupuesto» · revisión independiente (vuelta 3 de 3)

Fecha: 12-oct-2026 · Revisor: `revisor-presupuesto-r3` (no corrige) · Criterio: P3-07.
Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/presupuesto/` (P-06 a P-09, móvil 412 y escritorio 1440) y `docs/design/maquetas/lienzo/fragmentos/presupuesto.json`, después de `correccion-presupuesto-r2`.
Contra: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06 y U5-07, en lo que se ve en el markup), `docs/design/DESIGN.md` (§7, §8.3, §8.5, §8.6, §10.5), `docs/design/U1-arquitectura.md` (§2.1, §2.2, §6.4), U4 (P-26.02) y `docs/design/maquetas/referencia-tipo/format.md` y `craft.md`.
Diff contra la copia de la vuelta 2 (`docs/design/salida/iter-2/presupuesto/`): P-06, P-08, P-09 y el fragmento no cambian. P-07 (m y e) cambia solo en el SVG de la gráfica y en una línea de `renderVals()` (`hoyPath`).

**Veredicto: APROBADO** · 0 altas · 0 medias · 15 bajas (las mismas 15 de la vuelta 2, todas abiertas y sin cambio) · 1 observación sin severidad.
El único alto de la vuelta 2 (N1) queda **cerrado**. La corrección no trae fallas nuevas en forma, tokens, contraste, P3-12 ni enlaces.

## 1. Hallazgo de la vuelta 2

| Id | Estado | Evidencia |
|---|---|---|
| N1 · contraste de la línea a guiones del devengado acumulado y de «Hoy» sobre el relleno de la capa 1 (U5-04, P3-07, §8.6, §10.5 C-107) | **Cerrado** (opción a) | Se resolvió la plantilla con `renderVals()` en node. El SVG de P-07 pinta, en este orden: relleno de la capa 1, rayas de la capa 2, **halo** con el mismo `d` que el devengado (`stroke:var(--ui-surface-card)`, `calc(var(--size-viz-line) * 3)` = 6 px, `non-scaling-stroke`), **halo** de «Hoy» (mismo ancho), línea «Hoy» (`ui-viz-baseline`, `--border-width-default` = 1 px), línea del efectivo y línea a guiones (`ui-chart-2`, 2 px, `8 4`). El halo de 6 px bajo un trazo de 2 px deja 2 px de tarjeta a cada lado, así que el color contiguo a la línea a guiones es `ui-surface-card`: **3.15:1** en claro y **7.73:1** en oscuro (C-107). En «Hoy», el halo de 6 px bajo 1 px deja 2.5 px de tarjeta a cada lado: **3.47:1** y **6.09:1**. Medido sin halo, como en la vuelta 2: 2.30 y 2.54 en claro. Las series se pintan después de los halos, y el halo no tapa ninguna línea. `hoyPath` vale `M15.05,0 L15.05,200` en móvil y `M16.13,0 L16.13,280` en escritorio. Sale de la misma fracción que `hoyLeft` (5.4 %) y la etiqueta sigue alineada. La línea «Hoy» sigue siendo de 1 px, como pide §8.3. |

## 2. Tabla artboard × regla

Leyenda: OK = cumple; un id remite a §4 (todas las del §4 son bajas).

| Artboard | Forma | Tokens y hex | Contraste U5-04 | P3-12 | P3-08 | U5-06 | U5-07 | U1/U4 | Táctil y a11y | Enlaces |
|---|---|---|---|---|---|---|---|---|---|---|
| `presupuesto-p06-mes-m` | OK | OK | OK | OK | OK | B1 | OK | B8 | OK | OK |
| `presupuesto-p06-mes-e` | OK | OK | OK | OK | OK | B1, B9 | OK | B8 | OK | OK |
| `presupuesto-p07-flujo-m` | OK | OK | OK (N1 cerrado) | OK | OK | N4, N5, N7 | N3 | B5 | N2 | OK |
| `presupuesto-p07-flujo-e` | OK | OK | OK (N1 cerrado) | OK | OK | N4, N5, N7, B9 | N3 | B5, N6 | N2 | OK |
| `presupuesto-p08-ingresos-m` | OK | OK | OK | OK | OK | OK | N3 | B6 | N2 | OK |
| `presupuesto-p08-ingresos-e` | OK | OK | OK | OK | OK | B9 | N3 | OK | B7, N2 | OK |
| `presupuesto-p09-ahorro-m` | OK | OK | OK | OK | OK | OK | N3 | B6, B8 | N2 | OK |
| `presupuesto-p09-ahorro-e` | OK | OK | OK | OK | OK | B9 | N3 | B8 | N2 | OK |
| Fragmento `presupuesto.json` | OK | — | — | OK | — | — | — | — | — | B3, B4 |

## 3. Evidencia por criterio

Todo se comprobó con scripts propios en `datos/tmp/rev-pres-r3/` (`python3 -I` y node 26, borrados al terminar): 751 comprobaciones estáticas sobre la fuente y otra pasada sobre la plantilla resuelta con `renderVals()`.

**(a) Forma (format.md y encargo).**
- Cabecera exacta en los 8: `<!doctype html><html lang="es-MX"><head><meta charset="utf-8"><title>…</title><script src="./support.js"></script><link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css"></head><body><x-dc><helmet><style>`.
- Estilos base: `html,body{margin:0}`, `body{background:var(--ui-surface-page);color:var(--ui-text-default);font-family:var(--font-sans);…}`, más `a{…}` y `a:hover{…}` con tokens.
- Un solo `<x-dc>` y `</style></helmet>`. `<script type="text/x-dc" data-dc-script data-props='…'>` con `class Component extends DCLogic { renderVals() {`. El archivo cierra con `</script></body></html>`.
- `data-props` es JSON válido y su `$preview` es igual a `w`/`h` del fragmento en los 8: 412 × 4400, 3440, 1880 y 2830; 1440 × 2650, 2340, 1480 y 1880. Todas las alturas móviles son ≥ 915.
- `html.parser` no encuentra cierres sueltos ni etiquetas sin cerrar. Hay un `h1` por artboard, ningún id repetido y `hint-*` en todo `sc-for` y `sc-if`.
- `renderVals()` corre en node con un stub de `DCLogic`. La plantilla resuelta deja **0 huecos sin resolver** en los 8.
- Títulos: el nombre de la pantalla, como pide el encargo (m y e comparten título; el fragmento los distingue con «· móvil» y «· escritorio»).
- Nombres: `[a-z0-9.-]`, 90 stems en todo el lienzo, ninguno repetido.

**(b) Tokens.**
- Fuera de `data-props` no hay ningún hex. Tampoco `rgb()`, `hsl()`, `oklch()` ni colores con nombre. Los únicos `#` son `url(#pt-dev)` y anclas.
- No hay variables CSS locales.
- Todos los `var(--x)` existen en `tokens.json` (`--font-sans` sale de `type.families.sans`) y en `tokens.css`.
- Los iconos son SVG con `currentColor` y `aria-hidden="true"`. La gráfica usa `var(--…)` y `color-mix(in srgb, var(--…) 20%, var(--ui-surface-card))`.

**(c) Contraste (U5-04 y P3-07).** Calculé WCAG 2.x con los hex de `tokens.json`, resueltos por alias en claro y en oscuro.
- Los 13 pares de texto y fondo que se declaran en la misma regla dan todos ≥ 4.5:1 en los dos temas.
- Revisé además el producto cruzado de los 11 colores de texto contra los 23 fondos. Los pares bajo 4.5 caen en fondos que no llevan texto (barras, puntos de persona, muestras de leyenda, el tirador de la hoja) o en combinaciones que no ocurren. Comprobé en el markup que esos fondos tienen 0 nodos de texto.
- Pares que sí ocurren:

| Par (primer plano / fondo) | Umbral | Claro | Oscuro | §10.5 / estado |
|---|---|---|---|---|
| `ui-text-default` / `ui-surface-page` | 4.5 | 11.53 | 14.07 | OK |
| `ui-text-secondary` / `ui-surface-card` | 4.5 | 5.04 | 7.97 | C-016, C-125 |
| `ui-text-secondary` / `ui-surface-hover` | 4.5 | 4.51 | 5.50 | OK |
| `ui-text-link` / `ui-surface-page` | 4.5 | 5.53 | 10.38 | OK |
| `ui-text-on-primary` / `ui-action-primary-bg` | 4.5 | 5.77 | 10.47 | C-001 |
| `ui-role-warning` / `ui-role-warning-surface` | 4.5 | 5.12 | 6.09 | OK |
| `ui-role-info` / `ui-role-info-surface` | 4.5 | 4.96 | 6.12 | OK |
| `ui-role-success` / `ui-role-success-surface` | 4.5 | 4.81 | 6.19 | OK |
| `ui-role-error` / `ui-role-error-surface` (icono) | 3 | 5.37 | 6.05 | OK |
| `ui-icon-default` / `ui-surface-field` | 3 | 5.40 | 10.43 | OK |
| `ui-border-control` / `ui-surface-field` | 3 | 3.71 | 7.96 | OK |
| `ui-border-focus` / `ui-surface-card` | 3 | 3.56 | 6.03 | OK |
| `ui-viz-meter` / `ui-viz-track` | 3 | 4.51 | 5.50 | OK |
| `ui-chart-1` / `ui-surface-card` (línea del efectivo) | 3 | 7.23 | 3.22 | C-106 |
| **`ui-chart-2` / `ui-surface-card` (guiones sobre el halo)** | 3 | **3.15** | **7.73** | C-107 · N1 cerrado |
| **`ui-viz-baseline` / `ui-surface-card` («Hoy» sobre el halo)** | 3 | **3.47** | **6.09** | N1 cerrado |
| `ui-chart-2` / relleno de la capa 1 (sin halo, solo como referencia) | 3 | 2.30 | 6.10 | ya no es contiguo |

**(d) P3-12.**
- Grep de las 21 cifras canon de P3-12 (con y sin coma) sobre los 8 artboards y el fragmento: **0**. Las mismas en centavos: 0. Patrón de credenciales: 0.
- Fechas: el comando de P3-12 lista 15 fechas que coinciden con `model.json` y `data.json`. En presupuesto aparecen 6 de ellas, pero ninguna va con un monto del reporte ni de Kueski:
  - `2026-10-15` y `2026-10-30` en P-08, con $14,260.00 y $14,780.00 de `ingresos`.
  - `2026-11-15` a `2027-01-15` en P-07, como `fin` de las quincenas de `proyeccion.flujo_por_quincena`. No se muestran: la pantalla dice «1.ª quincena · nov-2026».
- Datos incrustados: todos son subconjuntos exactos de `datos-sinteticos.json`:
  - `CAT` ⊂ `categorias`; `MS` ⊂ `proyeccion.meses`;
  - `F` = `proyeccion.saldo_inicial_cents` + `flujo_por_quincena` + `supuestos`;
  - `DV` = `deuda_viva.a_fecha` + nombres de `creditos.lineas`;
  - `SU` = `presupuesto.gastos_fijos_recurrentes`, `despensa_programada`, `metodo_variable` y `proyeccion_meses_7_a_12`;
  - `P` = `presupuesto.{resumen, quincenas, despensa_devengada, alertas}`; `I` = `ingresos`; `A` y `MV` = `ahorro`.
- Montos visibles: 616 en total, 600 iguales a un campo `*_cents` o `*_fmt` del archivo. Los 16 restantes son sumas o campos sin sufijo `_cents`:
  - $257,740.56 y $242,520.65 son los acumulados;
  - $10,774.70 es la suma de los 8 gastos fijos;
  - $7,049.52, $2,797.02, $13,076.08 y $5,300.31 son `por_linea_cents` del 31-dic-2026;
  - $27,684.60 es el plan después del ajuste de U4.
- Fechas visibles: 59 en total; 53 son fechas del archivo. Las otras 6 son la misma: «10-oct-2026», el ejemplo de formato de la ayuda (ver N3).
- Muestra de 10 valores cotejados:
  - plan $27,586.00 e ingreso del mes $31,130.00 (`presupuesto.resumen`);
  - saldo inicial $12,000.00;
  - deuda viva $28,222.93 al 31-dic-2026 y $44,606.33 hoy;
  - Karen oct-2026 $906.98;
  - quincena del 30 $14,780.00 al 30-oct-2026;
  - ahorrado $41,730.00 y fondo $21,480.00;
  - «Hoy · 21-oct-2026».

**P3-08.**
1. ≤ 6 tarjetas: se pide en Inicio; P-06 lleva 6.
2. 14 secciones: no aplica.
3. Barra lateral colapsable: OK, con «Colapsar barra lateral» y `aria-expanded` en los 4 de escritorio.
4. Español de México: OK.
5. Dinero `$1,234.56` tabular: OK. `tabular-nums` está en `body` y en `.num`. Los negativos llevan U+2212 (0 casos de «-$»). Fuera del eje, 0 montos mal formados; el eje usa la forma compacta «$300 mil», que §8.5 permite.
6. Fechas `10-oct-2026`: OK; 0 fechas ISO visibles y 0 `type="date"`.
7. Paleta C2: OK.
8. Notion: OK.

**U5-06 (markup).** Serie 1 y 2 en su orden con su trazo y patrón. Valores finales en texto y «Ver como tabla». Iconos Lucide de 20/24 px con trazo de los tokens. Radios y espaciado con tokens. Las desviaciones son N4, N5, N7, B1 y B9 (bajas).

**U5-07.** 0 palabras en inglés visibles («CSV» es sigla), 0 fechas ISO visibles, 0 emoji, sin barra de estado ni teclado falsos.

**(e) Fidelidad (U1 §2.1, §2.2, §6.4; U4 P-26.02).**

| Pantalla | Presente | Falta o difiere (bajas) |
|---|---|---|
| P-06 | Selector de mes, pestañas Resumen · Detalle · Tendencias · Calendario, barras con umbral 80/100, alertas, quincenas, despensa devengada, «Ajustar presupuesto» con alcance K2 («De este mes en adelante» preseleccionado) y vista previa de U4 («Plan total del mes: antes $27,586.00, después $27,684.60»), «Ver gastos de la categoría» (en móvil, cada fila de categoría enlaza a P-02), estados vacío/carga/error en escritorio | B1, B8 |
| P-07 | «Simular compra» a P-20, tabla de quincenas, gráfica de área apilada con «Hoy» y valores finales, «Ver como tabla» con partes del hogar y de Karen, deuda viva con campo de fecha y partes, supuestos con el interruptor en solo lectura, estados en escritorio | B5, N4 a N7 |
| P-08 | «Capturar quincena», «Registrar ingreso extra», quincenas esperadas y capturadas, extras, calendario, estados (vacío «Sin ingreso esperado · Ir a Ajustes › Ingresos») | B6, B7 |
| P-09 | «Registrar aportación», «Nueva meta» a P-26.13, «Registrar retiro», fondo, metas con barra y estado (en ruta y atrasada; los datos no traen ninguna cumplida), movimientos, estados | B6, B8 |

En los 4 de escritorio está la tarjeta «Estados» (vacío, esqueleto y «Reintentar»), como pide U1 §2.2.

**(f) Accesibilidad.**
- 274 controles en la plantilla resuelta, todos con nombre accesible (texto o `aria-label`).
- No hay `role="button"`, `role="link"`, `role="switch"` ni `onclick`.
- Cada `input` tiene su `label for` o `aria-label`.
- Táctil en móvil, medido con los tokens resueltos: 107 controles; **0 bajo 44 px**. `.btn`, `.btn.icb` y los campos usan `--size-control-coarse` o `--size-touch-min` (44px). Los enlaces en línea llevan `min-height:var(--size-touch-min)`. Cada radio está dentro de un `label` de 44 px.

**(g) Raíces.**
- Móvil: `.mroot{width:412px;box-sizing:border-box;min-height:915px}`.
- Escritorio: sin ancho en px en la raíz. `.app{display:flex;flex-wrap:wrap}`, `.side{flex:1 1 var(--layout-sidebar-width)}`, `.main{flex:999 1 560px;min-width:0}` y un contenedor con `max-width` (literal, B9).
- En el fragmento, `"expand":"fill"` va solo en los 4 de escritorio.

**(h) Enlaces.**
- Todos los `href` resueltos van a un archivo existente de la lista del lienzo y del mismo ancho, a «#», o a un ancla interna con id existente (`#ajustar`, `#capturar`, `#aportar`; format.md lo permite).
- Fragmento: JSON válido; `boards` y `order` listan exactamente los 8 archivos; `page:"presupuesto"`; `is_interactive:false`; sin x ni y.

## 4. Hallazgos

### Alta
Ninguna.

### Media
Ninguna.

### Baja (siguen abiertas, sin cambio desde la vuelta 2)
- **B1 · DESIGN §7.6.** Las 13 barras de categoría de P-06 no llevan `role="meter"` ni `aria-valuetext`, que solo están en las 3 barras de resumen y de quincena. Tampoco llevan `triangle-alert` al 80 % o más.
- **B2 · K2 y format.md.** El `checked` sin valor en P-06 (alcance) y en P-09 (tipo) es HTML válido. Confirmar al publicar que salen preseleccionados.
- **B3 · format.md.** `is_interactive:false` en un prototipo con enlaces: lo decide el orquestador.
- **B4 · craft.md.** Las alturas del fragmento son estimadas: medirlas al publicar.
- **B5 · U1 N-092.** La deuda viva abre en 31-dic-2026 y no en el fin del horizonte, y no tiene desglose «por compra». Lo decide el dueño de los datos.
- **B6 · U1 §1.3.** En móvil, «Registrar ingreso extra» (P-08) y «Registrar retiro» (P-09) van como botones visibles y no en «⋯».
- **B7 · DESIGN §7.5.** El calendario de P-08 no tiene `aria-current="date"`.
- **B8 · U1 §1.3, §2.1, N-112.** El diálogo de P-06 se dibuja en línea y su selector trae 3 categorías. En P-09 faltan la vista Tablero y el filtro por periodo.
- **B9 · Estilo.** `max-width:1120px` literal en los 4 de escritorio, y la clase `.dn`.
- **N2 · DESIGN §7.3.** Los 6 campos de fecha llevan `inputmode="numeric"` con un valor `31-dic-2026`.
- **N3 · U5-07.** La ayuda «El selector empieza en lunes.» es una nota de diseño. Dejar «Formato 10-oct-2026».
- **N4 · DESIGN §8.6.** Las muestras de la leyenda no reproducen la marca: va color pleno donde el área está al 20 %.
- **N5 · DESIGN §8.3 y §8.6.** Valores finales en lista y no como etiqueta directa. El subtítulo no da la deuda viva con su suma.
- **N6 · Lienzo.** «Ver como tabla» va cerrado también en escritorio (`<details>` sin `open`), así que la tabla de 12 meses no se ve en el artboard estático.
- **N7 · DESIGN §8.2.** Los puntos están al final de cada mes y las etiquetas centradas en la celda: quedan desfasados media celda.

### Observación (sin severidad; para la medición en el render publicado)
- **O1 · Arranque de las dos series en P-07 móvil.** Las dos acumuladas parten de (0, 0). En el primer tramo (oct-2026, unos 24 px), la línea a guiones queda a ≈ 2.2 px perpendiculares de la del efectivo. Ahí su borde superior toca el trazo de `ui-chart-1` (2.30 en claro y 2.40 en oscuro entre series) y el inferior conserva ≥ 2 px de halo de tarjeta. Desde nov-2026 la separación es de 8 a 10 px y hay tarjeta a los dos lados. Es inherente a dos series que empiezan en cero y no baja el contraste contra la superficie. Se anota para mirarlo en la captura.
- Efecto visual del halo: abre una franja de 6 px color tarjeta dentro del relleno de la capa 1, debajo de la línea a guiones. Es lo que pedía la opción (a); revisarlo en la captura junto con N4.

## 5. Para la publicación
- Medir en el render publicado, en claro y en oscuro: el contraste de la gráfica de P-07 (incluido O1) y las alturas reales de los 8 artboards (B4).
- Las 15 bajas pueden ir a la referencia visual como «ajustar en implementación». Ninguna bloquea la aprobación.
