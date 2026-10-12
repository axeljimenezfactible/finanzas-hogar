Finanzas Hogar es la libreta de dinero de una casa: Axel, Ximena y Karen registran gastos, presupuestos, comidas, compras, créditos y salud en una sola app. La interfaz se parece a una libreta tipo Notion (barra lateral, páginas de bloques, vistas de base de datos, menú `/`, edición en el lugar): superficies casi blancas en claro y casi negras en oscuro, letra neutra y color solo para significar.

Aplica estos ocho principios en cada pantalla:

1. **Completo pero minimalista.** Muestra solo lo que pide la decisión de la pantalla y deja lo demás a un paso («Mostrar 3 más», menú «⋯», vista previa del registro). El Resumen de Inicio lleva 6 tarjetas como máximo.
2. **Una acción primaria por pantalla.** Un solo botón con relleno de marca visible a la vez; una hoja o un diálogo cuentan como pantalla propia.
3. **Móvil primero.** Diseña a 360 px y amplía. Registrar un gasto toma 3 toques.
4. **El color significa, no decora.** Todo color de rol, persona, tienda o serie va con signo, icono, forma o texto; la pantalla se entiende en escala de grises.
5. **Cifras exactas y legibles.** Dinero con cifras tabulares, nunca truncado, partido ni animado.
6. **Calma.** Densidad cómoda, sombras leves solo en claro, movimiento corto, sin rojos ni destellos salvo error y acción crítica.
7. **Claro y oscuro de primera.** Todo color sale de un token con valor claro y oscuro; revisa siempre los dos temas.
8. **Accesible por cálculo.** Texto ≥ 4.5:1; texto grande, controles, foco y gráficas ≥ 3:1.

## Tono y texto

- Escribe en español de México, de tú, con frases cortas, para una familia y no para un banco.
- Pon los botones en infinitivo: «Guardar», «Registrar gasto», «Cambiar semana», «Marcar pago».
- Quita la jerga: «Lo que debes hoy», no «saldo insoluto».
- Di qué pasó y qué hacer: «Falta capturar la tasa de esta tarjeta», no «Error de validación».
- No uses signos de exclamación ni emoji en ningún mensaje. El éxito es discreto: «Guardado · Deshacer».
- Respeta los textos fijados: «Sin conexión con el servidor de casa», «Requiere conexión con el servidor de casa», «Falta capturar · Capturar», «Aún no hay gastos este mes · Registrar gasto».
- Escribe fechas como `5-nov-2026` y `Jue 5 nov 2026`, porcentajes enteros como `82 %` y las tasas tal como se capturaron.
- Nombra las secciones como en la barra lateral: Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes, Ajustes; y las pestañas: Resumen · Detalle · Tendencias · Calendario.

| Así | No así |
|---|---|
| «Guardado · Deshacer» | «¡Gasto guardado exitosamente!» |
| «Te quedan $815.80 de despensa este mes» | «Disponible: 815.8» |
| «Este pago vence mañana» | «Pago próximo a vencer (D-1)» |

## Color

### Proporción 60-30-10

- **60 % neutro.** Pinta la página con `var(--ui-surface-page)`, las tarjetas y la barra lateral con `var(--ui-surface-card)` y el texto con `var(--ui-text-default)` y `var(--ui-text-secondary)`. Todos los neutros llevan un tinte índigo muy bajo.
- **30 % marca en pasos de superficie.** Usa `var(--ui-surface-selected)` para el elemento activo, la fila seleccionada y el chip de filtro activo.
- **10 % énfasis.** Reserva `var(--ui-action-primary-bg)` para la acción primaria, la casilla marcada, el interruptor encendido y el círculo de «hoy»; los roles, las categóricas y las series ocupan el resto del 10 %.
- La marca es índigo 280° y vive solo en controles (botón primario, enlace, anillo de foco, elemento activo). Nunca marca a una persona ni a una tienda.
- No agregues colores: todo color sale de `tokens.json`. Las transparencias y sombras se derivan de los neutros (`var(--shadow-1)`, `var(--shadow-2)`, `var(--ui-overlay-scrim)`).

### Superficies, texto y bordes

| Token | Úsalo para | Se lee sobre |
|---|---|---|
| `var(--ui-surface-page)` | fondo de página, barra superior, tablas, tableros y calendario | — |
| `var(--ui-surface-card)` | tarjetas, barra lateral, fondo de gráfica, renglones del calendario | — |
| `var(--ui-surface-raised)` | lo que flota: menús, menú `/`, diálogos, hojas, vista previa, avisos breves y el fondo de los chips | — |
| `var(--ui-surface-field)` | interior de los campos | — |
| `var(--ui-surface-hover)` | hover de elementos de menú, botón secundario, de texto y de icono | — |
| `var(--ui-surface-hover-subtle)` | hover de filas de tabla, celdas del calendario y edición en el lugar sobre la página (así las cifras de rol conservan 4.5:1) | — |
| `var(--ui-surface-pressed)` | presionado del botón secundario, de texto y de icono | — |
| `var(--ui-surface-disabled)` | relleno del control deshabilitado | — |
| `var(--ui-text-default)` | texto normal y cifras neutras | página, tarjeta, elevada, campo, seleccionada y superficies de rol |
| `var(--ui-text-secondary)` | texto secundario, marcador de posición, etiquetas de eje y encabezados de tabla | página, tarjeta, elevada y campo |
| `var(--ui-text-disabled)` | texto del control deshabilitado, siempre con borde punteado | `var(--ui-surface-disabled)` |
| `var(--ui-text-link)` · `var(--ui-text-link-hover)` | enlaces, siempre subrayados | página, tarjeta y elevada |
| `var(--ui-icon-default)` · `var(--ui-icon-selected)` | icono en reposo · icono del elemento activo | página, tarjeta, elevada · seleccionada |
| `var(--ui-border-divider)` | divisores y borde de lo que flota (decorativo) | — |
| `var(--ui-border-control)` · `var(--ui-border-control-hover)` | borde de campo, casilla y botón secundario (≥ 3:1) | página, tarjeta y campo |
| `var(--ui-border-selected)` | pestaña activa, destino de arrastre, chip elegido | elevada, página |
| `var(--ui-border-error)` | borde de campo o celda con error | campo |

En oscuro, `var(--ui-surface-raised)` es el mismo gris que la tarjeta: lo que flota se separa por su borde `var(--ui-border-divider)`, no por sombra.

### Acciones

| Botón | Relleno | Texto | Hover · presionado |
|---|---|---|---|
| Primario | `var(--ui-action-primary-bg)` | `var(--ui-text-on-primary)` | `var(--ui-action-primary-bg-hover)` · `var(--ui-action-primary-bg-pressed)` |
| Secundario | ninguno, borde `var(--ui-border-control)` | `var(--ui-text-default)` | `var(--ui-surface-hover)` · `var(--ui-surface-pressed)` |
| Destructivo | `var(--ui-action-destructive-bg)` | `var(--ui-text-on-destructive)` | `var(--ui-action-destructive-bg-hover)` · `var(--ui-action-destructive-bg-pressed)` |
| Acción crítica | `var(--ui-action-critical-bg)` | `var(--ui-text-on-critical)` | `var(--ui-action-critical-bg-hover)` · `var(--ui-action-critical-bg-pressed)` |

En claro el primario es marca 700 con texto blanco (5.77:1); en oscuro es marca 300 con texto marca 950 (10.47:1). Nunca pongas texto blanco sobre marca 600 (3.98:1).

### Roles

Pinta el texto y el icono de un rol con `var(--ui-role-<rol>)` y su insignia con `var(--ui-role-<rol>-surface)` de fondo. Un rol nunca va solo:

| Rol | Token | Refuerzo obligatorio |
|---|---|---|
| Ingreso | `var(--ui-role-income)` | signo «+» e icono `arrow-down-left` en la cifra |
| Gasto | `var(--ui-role-expense)` | signo «−» e icono `arrow-up-right` en la cifra |
| Ahorro | `var(--ui-role-savings)` | icono `piggy-bank` y etiqueta |
| Deuda / crédito | `var(--ui-role-debt)` | icono `credit-card` y etiqueta |
| Advertencia | `var(--ui-role-warning)` | icono `triangle-alert` y texto |
| Error | `var(--ui-role-error)` | icono `circle-alert` y mensaje; solo en validación y avisos |
| Éxito | `var(--ui-role-success)` | icono `circle-check` y texto («Guardado», «Pagado»); nunca en cifras |
| Información | `var(--ui-role-info)` | icono `info` y texto |

- Usa cada rol solo para su significado: el terracota del gasto no decora, el verde del éxito no marca ingresos.
- Ingreso y éxito casi no se distinguen lado a lado: el ingreso vive en cifras con «+», el éxito solo en avisos con palomita.
- Gasto y error nunca comparten componente: el gasto es una cifra con «−», el error es un mensaje con icono.
- Un vencimiento o un presupuesto rebasado es advertencia, no error.

### Acción crítica

Usa `var(--ui-critical)` (oliva, el eje opuesto a la marca) solo en el icono `shield-alert` del encabezado del diálogo que confirma algo irreversible, y `var(--ui-action-critical-bg)` en su único botón, con verbo y cantidad: «Borrar 12 gastos», «Restaurar respaldo». No la uses para errores ni advertencias.

### Personas y tiendas

| Grupo | Orden fijo | Tokens | Forma |
|---|---|---|---|
| Personas | Axel, Ximena, Karen, Conjunto | `var(--ui-person-axel)`, `var(--ui-person-ximena)`, `var(--ui-person-karen)`, `var(--ui-person-shared)` | círculo relleno de 10 px; Conjunto, anillo hueco de 2 px; chip en píldora |
| Tiendas | Walmart, Bodega Aurrera, Sam's Club, Mercado | `var(--ui-store-walmart)`, `var(--ui-store-bodega-aurrera)`, `var(--ui-store-sams-club)`, `var(--ui-store-market)` | cuadrado relleno de 10 px; chip con `var(--radius-sm)` |

- No distingas nunca a una persona de una tienda solo por color en una misma vista: cada color va con su nombre en texto y su forma.
- Personas y tiendas nunca comparten gráfica ni leyenda; en tablas van en columnas distintas. (Ximena y Bodega Aurrera son casi el mismo color en oscuro.)
- No escribas texto encima de un color categórico: el nombre va en `var(--ui-text-default)` junto al punto, anillo o marca.
- Los mercados adicionales usan `var(--ui-store-market)`.

### Series de gráficas

Asigna `var(--ui-chart-1)` a `var(--ui-chart-8)` en ese orden y no las recicles; una novena serie se junta en «Otros» con `var(--ui-viz-other)`. Las gráficas por persona o por tienda usan los colores de persona o tienda, no las series.

### Contraste medido

| Par | Claro | Oscuro |
|---|---|---|
| `var(--ui-text-default)` sobre `var(--ui-surface-page)` | 11.53:1 | 14.06:1 (sobre página) |
| `var(--ui-text-secondary)` sobre `var(--ui-surface-page)` · sobre tarjeta en oscuro | 5.39:1 | 7.97:1 |
| `var(--ui-text-link)` sobre `var(--ui-surface-page)` | 5.53:1 | 10.37:1 |
| `var(--ui-border-control)` sobre página · sobre tarjeta en oscuro | 3.71:1 | 6.08:1 |
| `var(--ui-border-focus)` sobre `var(--ui-surface-card)` | 3.56:1 | 6.03:1 |
| Texto de rol más ajustado: `var(--ui-role-success)` sobre tarjeta | 4.73:1 | — |

Antes de usar un par nuevo de primer plano y fondo, calcúlalo en los dos temas.

## Tipografía

Usa una sola familia, Inter variable autoalojada, con la pila `var(--font-sans)`; el tamaño óptico es automático y el espaciado entre letras queda en 0. Usa la cursiva solo en notas y citas.

| Estilo | Tamaño / interlínea · peso | Úsalo para |
|---|---|---|
| `title-1` | 32/40 · 700 | título de página, uno por página |
| `title-2` | 24/32 · 600 | cifra principal de una tarjeta de métrica, título de hoja o diálogo, monto grande de «Registrar gasto» |
| `title-3` | 20/28 · 600 | títulos de sección y de tarjeta, total al centro de una dona |
| `body` | 16/24 · 400 | texto de páginas y notas, campos (también en móvil, para que el teléfono no amplíe), botones en móvil |
| `small` | 14/20 · 400 | barra lateral, tablas, propiedades, menús, chips, botones en escritorio |
| `caption` | 12/16 · 400 | metadatos, ejes, pies de tabla, insignias; 500 en etiquetas |

- Usa el peso 500 (`var(--font-weight-medium)`) en etiquetas de campo, elemento activo y botones; 600 (`var(--font-weight-semibold)`) en títulos de 20 y 24 y en cifras principales; 700 (`var(--font-weight-bold)`) solo en el título de página; 400 (`var(--font-weight-regular)`) en el resto.
- No uses más de tres tamaños por pantalla además de `caption`: la jerarquía la dan el peso y el color.

### Dinero y cifras

- Pon cifras tabulares (`font-variant-numeric: tabular-nums`, clase `.money`) en todo dinero, porcentaje, cantidad y fecha de tabla.
- Escribe el dinero así: `$4,826.37` (símbolo, coma de miles, punto y dos decimales).
- Ingreso: `+$12,305.00` en `var(--ui-role-income)` con `arrow-down-left`.
- Gasto o negativo: `−$1,037.40`, con el signo menos tipográfico (U+2212) pegado al `$`, nunca el guion corto; color de gasto y `arrow-up-right`.
- Límite rebasado: `−$215.40` en `var(--ui-role-warning)` con `triangle-alert` y la palabra «rebasado».
- Montos neutros (totales, saldos sin rol): sin signo, en `var(--ui-text-default)`.
- Alinea el dinero a la derecha en tablas. No trunques ni partas una cifra: quita primero lo secundario (`$99,999.99` cabe completo a 360 px). No animes una cifra.

## Espaciado y rejilla

Usa solo `var(--space-0)`, `var(--space-1)` (4), `var(--space-2)` (8), `var(--space-3)` (12), `var(--space-4)` (16), `var(--space-6)` (24), `var(--space-8)` (32) y `var(--space-12)` (48) en márgenes, rellenos y separaciones.

- `var(--space-1)`: icono y texto dentro de un control; entre insignias.
- `var(--space-2)`: entre controles relacionados, etiqueta y campo, y como mínimo entre dos objetivos táctiles.
- `var(--space-3)`: relleno horizontal de botones, campos y tarjetas de tablero.
- `var(--space-4)`: margen lateral y medianil en móvil, relleno de tarjeta en móvil, entre campos.
- `var(--space-6)`: relleno de tarjeta desde 768, medianil desde 1024, entre bloques.
- `var(--space-8)` y `var(--space-12)`: margen lateral a 1024 y 1440, aire sobre el título de página.

Rejilla de 4 columnas hasta 767 px, 8 desde 768 y 12 desde 1024. Barra lateral de `var(--layout-sidebar-width)`; páginas de texto a `var(--layout-reading-max-width)` como máximo; las vistas de datos usan todo el ancho.

## Radios, sombras y bordes

- `var(--radius-sm)` (6): botones, campos, chips de tienda, elementos de menú, insignias.
- `var(--radius-md)` (10): tarjetas, menús emergentes, avisos breves, vista previa, franjas.
- `var(--radius-lg)` (14): diálogos, hojas inferiores (solo arriba) y el menú `/`.
- `var(--radius-full)`: avatares, puntos de persona, chips de persona, interruptores y la pista de avance.
- Un radio interior nunca es mayor que el de su contenedor.
- Deja las tarjetas planas (`var(--shadow-0)`): se separan por superficie. Usa `var(--shadow-1)` en menús, la tarjeta que se arrastra y la vista previa, y `var(--shadow-2)` en diálogos, hojas, menú `/` y avisos breves. En oscuro las sombras valen `none`: lo que flota lleva borde `var(--border-width-default)` en `var(--ui-border-divider)`.
- Pon bajo diálogos, hojas y el panel de navegación móvil el velo `var(--ui-overlay-scrim)`.

## Iconografía

- Usa Lucide 1.55.0 autoalojado, de línea, con trazo de 1.5 px dibujados en los dos tamaños: `var(--size-icon-sm)` (20 px, `stroke-width: var(--icon-stroke-user-units-sm)`) junto a texto de 14 y 16, y `var(--size-icon-md)` (24 px, `stroke-width: var(--icon-stroke-user-units-md)`) junto a títulos de 20 px o más y en los botones de solo icono de la barra superior y la barra inferior.
- Dibuja con `stroke: currentColor`, extremos y uniones redondeados, sin relleno: el icono toma el color del texto que acompaña.
- Pon un icono de solo icono en un área de `var(--size-touch-min)` y dale `aria-label`; en escritorio, también descripción emergente.
- No uses emoji como icono ni como adorno.

| Significado | Icono |
|---|---|
| Ingreso · Gasto | `arrow-down-left` · `arrow-up-right` |
| Ahorro · Deuda | `piggy-bank` · `credit-card` |
| Advertencia · Error · Éxito · Información | `triangle-alert` · `circle-alert` · `circle-check` · `info` |
| Acción crítica | `shield-alert` |
| Inicio · Comidas · Registrar · Compras · Buscar | `house` · `utensils` · `plus` · `shopping-cart` · `search` |
| Menú · Más acciones · Colapsar | `menu` · `ellipsis` · `panel-left` |
| Desplegar | `chevron-right` (gira a `chevron-down`) |
| Calendario · Pagos próximos | `calendar` · `calendar-clock` |
| Filtrar · Ordenar · Arrastrar | `list-filter` · `arrow-up-down` · `grip-vertical` |
| Deshacer · Cerrar · Borrar · Editar | `undo-2` · `x` · `trash` · `pencil` |
| Cargando · Sin conexión · Reintentar | `loader-circle` · `wifi-off` · `refresh-cw` |
| Tabla · Tablero · Calendario · Galería | `table-2` · `kanban` · `calendar` · `layout-grid` |
| Ajustes | `settings` |

## Movimiento

- Usa `var(--duration-fast)` (150 ms) para hover, presionado, casillas e interruptores y para toda salida; `var(--duration-base)` (200 ms) para menús, aviso que entra, pestaña y nodo de la barra lateral; `var(--duration-slow)` (250 ms) para diálogos, hojas, panel de navegación y vista previa.
- Entra con `var(--ease-emphasized)` y sal o cambia dentro de la pantalla con `var(--ease-standard)`.
- Anima solo `opacity` y `transform`. El anillo de foco aparece sin transición.
- Con `prefers-reduced-motion: reduce` todo cambio es instantáneo; solo `loader-circle` sigue girando, a la mitad de velocidad.
- Nunca animes cifras de dinero ni gráficas; nada parpadea.

## Estados

| Estado | Cómo se ve |
|---|---|
| Foco | anillo de `var(--border-width-focus)` en `var(--ui-border-focus)`, separado `var(--focus-offset)`; en filas y celdas, por dentro. Nunca lo quites |
| Deshabilitado | `var(--ui-surface-disabled)`, texto `var(--ui-text-disabled)` y borde punteado en `var(--ui-border-control)`; di la razón en texto («Requiere conexión con el servidor de casa») en lugar de esconder el control |
| Error de campo | borde `var(--ui-border-error)`, `circle-alert` y mensaje en `var(--ui-role-error)` que dice qué pasó y qué hacer; lo capturado no se borra |
| Carga | esqueleto estático con la forma final en `var(--ui-skeleton-on-page)` o `var(--ui-skeleton-on-card)`; `loader-circle` solo en acciones de más de 400 ms; a los 10 s, «Esto está tardando más de lo normal · Reintentar». Nunca cifras provisionales |
| Vacío | icono de la entidad en `var(--ui-icon-default)`, una frase en `var(--ui-text-secondary)` y el botón que lo llena: «Aún no hay gastos este mes · Registrar gasto» |
| Error de vista | franja en `var(--ui-role-error-surface)` con `circle-alert`, texto `var(--ui-text-default)` y «Reintentar» |
| Falta capturar · Sin conexión | franja en `var(--ui-role-warning-surface)` con `triangle-alert` o `wifi-off`: «Falta capturar · Capturar», «Sin conexión con el servidor de casa · 2 en cola» |
| Aviso breve | «Guardado · Deshacer» dura 5 s y se pausa con hover o foco; el de error se queda hasta cerrarlo |

## Visualización de datos

- Elige la forma por la pregunta: barras horizontales para «¿cuánto en cada categoría?», líneas para el tiempo, columnas para pocos meses (hasta 3 series), área apilada para devengado contra pagado, dona para un total de hasta 5 partes y tarjeta de métrica para una sola cifra.
- Dale a cada serie color, marcador, trazo y patrón: 1 círculo, continuo, liso; 2 cuadrado, `var(--viz-series-2-dash)`, rayas a 45°; 3 triángulo, `var(--viz-series-3-dash)`, rayas a 135°; 4 rombo, `var(--viz-series-4-dash)`, puntos; 5 triángulo invertido, `var(--viz-series-5-dash)`, cuadrícula diagonal; 6 cruz, `var(--viz-series-6-dash)`, rayas anchas a 45°; 7 aspa, `var(--viz-series-7-dash)`, rayas anchas a 135°; 8 anillo, `var(--viz-series-8-dash)`, puntos grandes.
- Con 2 o más series pon leyenda siempre; con 4 o menos, también etiqueta directa.
- Barras de `var(--size-viz-bar-max)` como máximo con extremo redondeado de `var(--size-viz-bar-radius)`; líneas de `var(--size-viz-line)`; marcadores de `var(--size-viz-marker)` con anillo de 2 px en `var(--ui-surface-card)`; `var(--size-viz-gap)` entre segmentos.
- Retícula horizontal en `var(--ui-viz-grid)`, línea base en `var(--ui-viz-baseline)`, un solo eje Y, sin retícula vertical salvo «Hoy».
- Escribe etiquetas, valores y leyendas en `var(--ui-text-default)` o `var(--ui-text-secondary)`, nunca en el color de la serie.
- Escribe cada monto en texto: al final de la barra, último valor de la línea, leyenda de la dona. El formato compacto (`$5 mil`) va solo en ejes. «Ver como tabla» siempre está.
- Barra de avance: pista `var(--ui-viz-track)` de `var(--size-meter-height)`; llenado `var(--ui-viz-meter)`; desde el primer umbral de Ajustes (80 % al inicio), `var(--ui-role-warning)` con `triangle-alert` y «82 % usado»; el excedente se dice en texto.
- No hagas doble eje, 3D, arcoíris, donas de más de 5 partes, montos solo en la descripción emergente, animación de entrada ni series atenuadas al filtrar.

## Accesibilidad

- Cumple WCAG 2.2 AA calculado en los dos temas: texto ≥ 4.5:1; texto grande, bordes de control, foco, iconos que significan y marcas de gráfica ≥ 3:1.
- Haz todo control táctil de al menos `var(--size-touch-min)` (44 × 44 px) aunque el dibujo sea menor; con puntero grueso los controles miden `var(--size-control-coarse)`, el primario de hoja `var(--size-control-primary-sheet)` y, solo con puntero fino, `var(--size-control-fine)`. Deja `var(--space-2)` entre objetivos.
- Muestra el foco con `:focus-visible` y nunca lo ocultes; con `forced-colors: active` usa `Highlight` y conserva bordes, forma, trazo y patrón.
- Usa HTML nativo y los patrones ARIA de cada componente: *grid*, *tabs*, *combobox*, *dialog*, *meter*, *switch*. Toda etiqueta es `<label for>`; el marcador de posición no la sustituye.
- Ofrece siempre una alternativa a arrastrar («Mover a…»).
- Declara `lang="es-MX"`.

## Puntos de quiebre

Diseña a `var(--breakpoint-base)` y amplía con consultas `min-width`:

| Desde | Cambia |
|---|---|
| `var(--breakpoint-base)` (360) | 4 columnas, barra inferior de 5 destinos, hojas a pantalla completa, tabla de 2 líneas, agenda de pagos |
| `var(--breakpoint-phone)` (412) | un dato secundario más por tarjeta y fila; galería de 2 columnas |
| `var(--breakpoint-tablet)` (768) | 8 columnas, sin barra inferior, tarjetas en 2 columnas, diálogos centrados, calendario mensual, barra lateral superpuesta de `var(--layout-sidebar-overlay-width)` |
| `var(--breakpoint-desktop)` (1024) | 12 columnas, barra lateral fija y colapsable, tabla completa, vista previa lateral de `var(--layout-preview-width-desktop)` |
| `var(--breakpoint-wide)` (1440) | más aire, texto a `var(--layout-reading-max-width)`, vista previa de `var(--layout-preview-width-wide)` |

A 320 px todo pasa a una columna sin desplazamiento horizontal de la página; solo tabla, tablero y fila de pestañas se desplazan en su propio contenedor. Las variables `--breakpoint-*` documentan el valor: en una consulta `@media` escribe el número.

## Datos sintéticos

- En maquetas, vistas previas y ejemplos usa solo cifras inventadas con la forma y la magnitud de las reales: `$4,826.37`, `−$1,037.40`, `+$12,305.00`, `82 %`.
- No copies montos, saldos, tasas ni fechas de reportes o estados de cuenta reales.
- Toma todas las cifras de una pantalla y de sus pantallas hermanas del mismo juego de datos, para que Inicio, Detalle y Reportes cuadren entre sí.
- Los nombres Axel, Ximena y Karen sí aparecen; las tiendas son Walmart, Bodega Aurrera, Sam's Club y Mercado.
- Nunca escribas credenciales, correos ni números de cuenta.
