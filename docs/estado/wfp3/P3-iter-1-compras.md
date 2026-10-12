# P3 · Vuelta 1 · Área «compras» · Informe de revisión (P3-07)

Fecha: 12-oct-2026. Revisor independiente: `revisor-compras-r1` (vuelta 1 de 3). No escribió ni corrigió nada de lo revisado. Fuentes de verdad: DESIGN.md 1.0 y C2 1.0 (Q2-96); U1 §0.4, §1.2, §1.3, §2.1, §2.2, §2.4, §3.15 a §3.18, §6.4 y F-3; criterios P3 1.3 (P3-08, P3-12); U5-04, U5-06 y U5-07 en lo que se ve en el código; `referencia-tipo/format.md`.

Objeto: los 8 artboards de `docs/design/maquetas/lienzo/project/compras/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/compras.json`. Reporte del autor: `docs/estado/wfp3/artboards-compras.md`.

**Dictamen: NO APROBADO.** Hay 2 hallazgos altos y 5 medios, más 12 bajos que no bloquean.

## Resumen de hallazgos

| Id | Sev. | Hallazgo | Artboards |
|---|---|---|---|
| A-01 | alta | El bloque `@media (max-width:799px){` del `<style>` nunca se cierra. Quedan dentro entre 17 y 21 reglas que en 1440 no se aplican, y el escritorio se rompe. | los 4 `-e` |
| A-02 | alta | Fechas cortas sin año («2 oct», «19 oct», «18 sep», «19–25 oct 2026») en lugar de `10-oct-2026`. Contradice la decisión B (U1 §0.4, DESIGN §3.3, P3-08 y U5-07). | p13-m, p13-e, p14-e, p15-m, p15-e (y la opción del selector en p14-m) |
| M-01 | media | Montos dentro de frases sin cifras tabulares (DESIGN §3.3, «**todo** dinero»). | p14-m, p14-e |
| M-02 | media | El «+» de la barra inferior mide 40 px de ancho (`min-width:0`) y queda en una cápsula de 40 × 56. Es un objetivo táctil de menos de 44 px y un destino de menos de 64 px (DESIGN §7.10 y §9.3). | los 4 `-m` |
| M-03 | media | A la barra lateral le faltan el botón «‹‹» (`panel-left`) para colapsarla, el interruptor de tema del pie y los chevrons ▸ como botones propios en cada sección (DESIGN §7.9, U1 §1.2 y P3-08 «barra lateral colapsable»). | los 4 `-e` |
| M-04 | media | El título visible de la página lleva la etiqueta del lienzo: «Lista de compras · móvil», «Comparar canasta · móvil», «Productos y precios · móvil» e «Inventario · móvil». | los 4 `-m` |
| M-05 | media | Historial de P-15: todos los marcadores son círculos. Mercado Portales y Tianguis comparten color y marcador y solo los distingue el guion. Con 4 series falta la etiqueta directa (DESIGN §8.6 y §8.7, U5-06). | p15-e |

Los bajos (B-01 a B-12) se detallan en la sección «Correcciones para la vuelta 2».

## Verificaciones (a)–(h)

Scripts propios con `python3 -I` en `datos/tmp/rev-compras-r1/` (`verif.py`, `clases.py`, `tabular.py` y `contraste.py`). Se borraron al terminar.

| | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a | Forma | **FALLA (A-01)** | Las 8 piezas tienen exacto el prefijo obligatorio: doctype, `lang="es-MX"`, `<title>` propio, `./support.js`, el enlace a `../ds/finanzas-hogar/tokens.css`, `x-dc` y `helmet`. Todas cierran con `</x-dc>`, `data-dc-script`, `class Component extends DCLogic` y `renderVals()`. `data-props` parsea en las 8 y `$preview` = `w`×`h` del fragmento en las 8. `html.parser` no encuentra etiquetas sin cerrar ni cierres inesperados. Los 8 títulos son únicos, no hay stems duplicados en todo el lienzo y los nombres cumplen `[A-Za-z0-9.-]`. **Pero el CSS de los 4 `-e` tiene 1 llave sin cerrar** (ver A-01). |
| b | Tokens | CUMPLE | 0 hex literales fuera de `data-props` (los SVG usan `currentColor` o `var(--…)`) y 0 `rgb()`, `hsl()` y `oklch()`. Las 8 piezas usan entre 50 y 65 `var(--x)` distintos y todos existen en `tokens.json` (nombres de las listas, más `font-sans` de `type.families`). Los tamaños de letra van en px literales (12, 14, 16, 20, 24 y uno de 11, B-02): tokens.json no trae tokens de tamaño. |
| c | Contraste | CUMPLE en el código tal como se quiso escribir | 45 pares × 2 temas = 90 combinaciones recalculadas con los hex de `tokens.json`, resolviendo las referencias `{…}`, con la fórmula WCAG 2.x. 0 quedan bajo umbral (tabla abajo). 26 pares están literalmente en §10.5 y los 19 restantes cumplen por cálculo. **Ojo:** por A-01, en 1440 la clase `.ax` del eje de P-15 no se aplica y el texto del eje queda en el `fill` negro por defecto, así que en tema oscuro sale negro sobre `ui-surface-card` (#36363a), por debajo de 4.5:1. Se corrige con A-01. |
| d | P3-12 | CUMPLE | El grep de las 21 cifras canon de P3-12 (con y sin coma) sobre `project/compras/` y `fragmentos/compras.json` da 0. También hay 0 en centavos, 0 montos de 4 cifras sin separador (P3-08) y 0 coincidencias de credenciales. Las fechas del cruce `oraculo/model.json` ∪ `data.json` ∩ `datos-sinteticos.json` van del 30-oct-2026 en adelante. Ningún artboard de compras muestra fechas posteriores al 21-oct-2026, así que no hay coincidencias. Muestra de procedencia (todas cuadran con `datos-sinteticos.json`): total de la lista $1,564.10 (`lista_compras.total_cents` 156410); Walmart $318.60 (grupo); Bodega Aurrera $1,551.70 y Walmart $1,638.30 (`canasta.por_tienda`); +$86.60 y +$80.20 (`diferencia_contra_la_mas_barata`); −$212.40, −$17.70 y $650.00 (`membresia_sams`); «Compra del 12-oct-2026 · Walmart · $2,918.40» (`gastos` g13 y `tickets` tk1); Jitomate $31.30 / kg en Tianguis, el 13-oct (`productos[].precios`); «40 de ticket, 62 manuales, 3 de PROFECO» (conteo de `precios[].fuente`); 32 frescos y 1 viejo; inventario 1 · 5 · 8; Leche 1,500 ml, 3,500 ml / sem y 3.0 días; ganadores 11 · 5 · 3 · 3 · 2 con 0 empates. Además se cruzaron las 120 celdas del Detalle de P-14 (24 artículos × 5 tiendas) y la diferencia de cada fila: 0 discrepancias, orden descendente correcto. |
| e | Fidelidad | Ver tabla de fidelidad | Están todos los elementos de U1 §2.1 y todos los estados de §2.2. Faltan elementos secundarios (M-03, M-05 y los bajos). |
| f | Accesibilidad | **FALLA (M-02, M-04)** | Hay `button`, `a`, `input` y `select` reales. 0 botones o enlaces sin nombre accesible, 0 controles sin `label` (los de búsqueda usan `label.sr`; las casillas, `aria-label` interpolado). Iconos decorativos con `aria-hidden`; los de solo icono llevan `aria-label`. Foco visible con `outline` de `ui-border-focus`. 0 emoji, 0 barra de estado falsa, 0 teclado falso. 0 palabras de interfaz en inglés (la lista de U5-07 y más). «Error» como título de estado es español; «USDA Foundation Foods» es nombre propio de la lista permitida. 0 fechas ISO visibles. Todos los montos cumplen `$1,234.56` y los negativos llevan U+2212 (−$212.40, −$17.70, −1,500 ml). En móvil se midieron 44 px en `.btn`, `.ib`, `.chk`, `.mg`, `.tab`, campos y `select` (`size-control-coarse` o `size-touch-min`, 44) y 56 × 64 en los destinos de la barra inferior, **salvo el «+»**, que mide 40 px de ancho (M-02). Cifras tabulares en todas las celdas y tarjetas de dinero, **salvo en frases** (M-01). |
| g | Raíz | CUMPLE | Los móviles tienen raíz `width:412px` con `*{box-sizing:border-box}`. Los de escritorio tienen raíz `.app` sin ancho en px (`flex-wrap`, barra lateral de 240 y `.main` con `flex: 999 1 560px`, aunque esta regla está atrapada por A-01). En el fragmento, el escritorio lleva `"expand":"fill"` y el móvil no. El fragmento tiene `boards` y `order` con exactamente los 8 archivos, `page` «compras» y `is_interactive` false. |
| h | Enlaces | CUMPLE | Todos los `href` apuntan a artboards de la lista del lienzo con el mismo ancho, o a «#» (Buscar, Detalle en móvil, Nuevo producto). Los 22 artboards de destino existen hoy en disco; `../ds/finanzas-hogar/tokens.css` todavía no, porque lo instala el orquestador. F-3 Ruta A cuadra: P-14 «Usar Bodega Aurrera» lleva a P-13 y P-13 «Comparar tiendas» lleva a P-14. «Simular pago con crédito» lleva a P-20 y «Terminar compra» y «Foto de ticket», a P-05. |

### Recálculo de contraste (U5-04, P3-07): 45 pares, claro y oscuro

| Par (primer plano / fondo) | Umbral | Claro | Oscuro | En §10.5 | Dónde |
|---|---|---|---|---|---|
| ui-text-default / ui-surface-page | 4.5 | 11.53 | 14.07 | sí | cuerpo, tabla |
| ui-text-default / ui-surface-card | 4.5 | 10.76 | 10.76 | sí | tarjetas, barra lateral, fila de grupo |
| ui-text-secondary / ui-surface-page | 4.5 | 5.40 | 10.43 | sí | subtítulos, encabezado de tabla |
| ui-text-secondary / ui-surface-card | 4.5 | 5.04 | 7.97 | sí | `.cap.sec` en tarjetas |
| ui-text-secondary / ui-surface-raised | 4.5 | 5.40 | 7.97 | sí | barra inferior |
| ui-text-default / ui-surface-raised | 4.5 | 11.53 | 10.76 | sí | chips `.ch`, destino activo |
| ui-text-default / ui-surface-selected | 4.5 | 10.73 | 10.90 | sí | celda ganadora, fila seleccionada |
| ui-text-secondary / ui-surface-selected | 4.5 | 5.02 | 8.08 | sí | fecha bajo el precio ganador |
| ui-text-default / ui-surface-hover | 4.5 | 9.64 | 7.42 | sí | insignia neutra `.bn` |
| ui-text-default / ui-surface-hover-subtle | 4.5 | 10.76 | 10.76 | sí | fila en hover |
| ui-text-secondary / ui-surface-hover-subtle | 4.5 | 5.04 | 7.97 | no (cumple) | `.cap.sec` en fila en hover |
| ui-text-on-primary / ui-action-primary-bg | 4.5 | 5.77 | 10.47 | sí | botón primario |
| ui-text-on-primary / ui-action-primary-bg-hover | 4.5 | 8.49 | 12.66 | sí | botón primario en hover |
| ui-text-link / ui-surface-page | 4.5 | 5.53 | 10.38 | sí | enlace |
| ui-text-link / ui-surface-card | 4.5 | 5.16 | 7.94 | sí | enlace en tarjeta |
| ui-text-link / ui-role-info-surface | 4.5 | 5.19 | 7.81 | no (cumple) | enlace en franja de vacío |
| ui-text-link-hover / ui-surface-card | 4.5 | 7.59 | 9.59 | sí | enlace en hover |
| ui-role-warning / ui-role-warning-surface | 4.5 | 5.12 | 6.09 | no (cumple) | insignia `.bw` |
| ui-role-success / ui-role-success-surface | 4.5 | 4.81 | 6.19 | no (cumple) | insignia `.bg` |
| ui-role-error / ui-role-error-surface | 4.5 | 5.37 | 6.05 | no (cumple) | insignia `.be` |
| ui-text-default / ui-role-info-surface | 4.5 | 10.83 | 10.59 | no (cumple) | franja `.co` |
| ui-text-default / ui-role-warning-surface | 4.5 | 10.74 | 10.86 | no (cumple) | franja `.co.w` |
| ui-text-default / ui-role-success-surface | 4.5 | 10.93 | 10.37 | no (cumple) | franja `.co.ok` |
| ui-text-default / ui-role-error-surface | 4.5 | 10.72 | 11.20 | no (cumple) | franja `.co.x` |
| ui-role-info / ui-role-info-surface | 3 | 4.96 | 6.12 | no (cumple) | icono de franja |
| ui-icon-selected / ui-surface-selected | 3 | 5.15 | 8.04 | sí | icono activo de la barra lateral |
| ui-icon-selected / ui-surface-raised | 3 | 5.53 | 7.94 | no (cumple) | icono activo de la barra inferior |
| ui-border-control / ui-surface-page | 3 | 3.71 | 7.96 | sí | borde de botón secundario y campo |
| ui-border-control / ui-surface-card | 3 | 3.47 | 6.09 | sí | campos del ajuste en tarjeta |
| ui-border-focus / ui-surface-page | 3 | 3.82 | 7.89 | sí | anillo de foco |
| ui-border-focus / ui-surface-card | 3 | 3.56 | 6.03 | sí | anillo de foco en tarjeta |
| ui-border-selected / ui-surface-page | 3 | 5.53 | 10.38 | sí | indicador de pestaña |
| ui-viz-meter / ui-viz-track | 3 | 4.51 | 5.50 | sí | medidor |
| ui-viz-meter / ui-surface-card | 3 | 5.04 | 7.97 | sí | medidor en tarjeta |
| ui-store-walmart / ui-surface-card | 3 | 5.11 | 4.76 | no (cumple) | marca y línea Walmart |
| ui-store-bodega-aurrera / ui-surface-card | 3 | 8.40 | 6.79 | no (cumple) | marca y línea Bodega |
| ui-store-sams-club / ui-surface-card | 3 | 3.06 | 10.18 | no (cumple) | marca Sam's |
| ui-store-market / ui-surface-card | 3 | 12.82 | 3.06 | no (cumple) | marca y línea de mercados |
| ui-store-walmart / ui-surface-page | 3 | 5.48 | 6.22 | no (cumple) | marca en tabla |
| ui-store-market / ui-surface-page | 3 | 13.74 | 4.00 | no (cumple) | marca en tabla |
| ui-store-bodega-aurrera / ui-surface-selected | 3 | 8.38 | 6.88 | no (cumple) | marca en celda ganadora |
| ui-action-primary-bg / ui-surface-raised | 3 | 5.53 | 7.94 | sí | «+» de la barra inferior |
| ui-action-primary-bg / ui-surface-page | 3 | 5.53 | 10.38 | sí | casilla marcada (`accent-color`) |
| ui-person-axel / ui-surface-card | 3 | 11.93 | 3.10 | no (cumple) | anillo del avatar |
| ui-text-default / ui-surface-field | 4.5 | 11.53 | 14.07 | sí | texto de campo |

«En §10.5 = no» quiere decir que mi búsqueda literal de «`x` sobre `y`» no encontró el par. §10.5 nombra varios pares como «rol X» o «tienda». Los 19 cumplen por cálculo. Exentos por decorativos (§10.5): divisores (`ui-border-divider`), retícula (`ui-viz-grid`) y esqueletos (`ui-skeleton-on-card`).

## Fidelidad: pantalla → elementos de U1

| Pantalla | Elemento (fuente) | m | e |
|---|---|---|---|
| P-13 | Grupos por vertical: cadenas por tienda y mercados con día de visita (§2.1, N-150) | ✓ grupos plegables con `aria-expanded` | ✓ filas de grupo con total |
| P-13 | Casillas de comprado (§2.1, N-151) | ✓ 44 px | ✓ |
| P-13 | Totales por tienda y total general (§2.1) | ✓ $1,564.10 | ✓ |
| P-13 | Tienda asignada contra recomendada; motivo; precio con frescura (N-150) | ✓ «Mejor en …» y motivo por aviso (B-06) | ✓ columnas Motivo, Tienda asignada y «Recomendada» |
| P-13 | Acción primaria «Terminar compra» y su paso a P-05; secundarias «Comparar tiendas», «Agregar artículo» y «Cambiar semana» (§2.2) | ✓ (secundarias a la vista, no en «⋯», B-05) | ✓ |
| P-13 | Estados: «Nada que comprar…», «Sin plan esta semana · Planear», esqueleto y Reintentar (§2.2) | n/a | ✓ |
| P-13 | Fechas `10-oct-2026` (§0.4) | ✗ A-02 | ✗ A-02 |
| P-14 | Selector de canasta (§2.1, N-160) | ✓ `select` con `label` | ✓ |
| P-14 | Totales por tienda con idénticos, sustitutos, sin precio y viejos, más la diferencia contra la más barata en total y en idénticos (N-161) | ✓ | ✓ |
| P-14 | Ganadores y empates (N-163) | ✓ | ✓ |
| P-14 | Membresía de Sam's con punto de equilibrio y veredicto (N-164) | ✓ plegada (U1 §6.4) | ✓ costo, ahorro por canasta, compras al año, ahorro anual, punto de equilibrio y veredicto |
| P-14 | Detalle por artículo: precio, fecha, ganadora, sustituto y viejo (N-162) | pestaña «#» (U1 §6.4) | ✓ (sin fuente ni precio por unidad base, B-09) |
| P-14 | Primaria «Usar Bodega Aurrera» (lleva a P-13) y «Simular pago con crédito» (lleva a P-20) (§2.2, F-3) | ✓ | ✓ |
| P-14 | Estados: «La canasta no tiene artículos», «Sin precios para comparar · Capturar precios», esqueleto y Reintentar | n/a | ✓ |
| P-15 | Precio por unidad base, marca de viejo, fuente e historial (§2.1, N-170 y N-171) | ✓ filas compactas, como pide §6.4 | ✓ tabla, vista previa e historial (M-05) |
| P-15 | Resumen de frescura y fecha más reciente por tienda (N-172) | ✓ insignias | ✓ |
| P-15 | Primaria «Capturar precio»; «Foto de ticket», «Nuevo producto» y «Filtrar» (§2.2) | ✓ (`⋯` y Filtrar) | ✓ con filtros de tienda, categoría, frescura y fuente |
| P-15 | Estados: «Sin productos», esqueleto y Reintentar | n/a | ✓ |
| P-16 | Suficiente, reponer o agotado; días restantes; movimientos (§2.1, N-180 a N-182) | ✓ ordenado por estado | ✓ con vista previa y movimientos |
| P-16 | Primaria «Ajustar existencia» con motivo (conteo, consumo o merma); «Registrar sobrante»; «Filtrar» (§2.2, N-183 y N-184) | ✓ | ✓ formulario con `label` |
| P-16 | Estados: «Inventario vacío: se llena al confirmar tickets», esqueleto y Reintentar | n/a | ✓ |
| Global | Barra lateral: orden fijo de 9 secciones, «Buscar» y «Registrar gasto», encabezado con avatar (U1 §1.2) | n/a (☰) | ✓ orden; ✗ «‹‹», tema y chevrons (M-03) |
| Global | Barra inferior: Inicio · Comidas · «+» · Compras · Buscar (G-05) | ✓ orden; ✗ «+» (M-02) | n/a |

Cobertura de P3-06 en el área: 4 pantallas (P-13 a P-16) × 2 anchos = 8 artboards. Están las 8 y el fragmento tiene 8. Faltantes: 0.

## P3-08: lista de 8 puntos (área compras, vuelta 1)

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen con ≤ 6 tarjetas | n/a (no es Inicio). Las filas de métricas del área tienen 4 tarjetas. |
| 2 | 14 secciones de Ajustes | n/a |
| 3 | Barra lateral colapsable | **FALLA**: no se dibuja el control «‹‹» (M-03) |
| 4 | Español de México | CUMPLE (M-04 no es inglés, pero es texto que no corresponde) |
| 5 | Dinero `$1,234.56` con cifras tabulares | Formato CUMPLE (0 montos de 4 cifras sin coma). Tabulares: **FALLA parcial** (M-01) |
| 6 | Fechas `10-oct-2026` | **FALLA** (A-02). P-16 e y la vista previa de P-15 sí cumplen |
| 7 | Paleta de C2 (índigo 280°) | CUMPLE: solo `var(--ui-*)`, 0 hex |
| 8 | Notion como única referencia | CUMPLE: barra lateral, vistas de tabla, franjas tipo callout, propiedades y vista previa lateral |

## Tabla artboard × regla

Leyenda: ✓ cumple · ✗ falla (Id) · n/a no aplica.

| Artboard | a Forma | b Tokens | c Contraste | d P3-12 | e Fidelidad | f Accesibilidad | g Raíz | h Enlaces | Fechas | Tabulares | Veredicto |
|---|---|---|---|---|---|---|---|---|---|---|---|
| compras-p13-lista-m | ✓ | ✓ | ✓ | ✓ | ✓ (B-05, B-06) | ✗ M-02, M-04 | ✓ (B-03) | ✓ | ✗ A-02 | ✓ | no aprobado |
| compras-p13-lista-e | ✗ A-01 | ✓ | ✓* | ✓ | ✗ M-03 | ✓ | ✓ | ✓ | ✗ A-02 | ✓ | no aprobado |
| compras-p14-canasta-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ M-02, M-04 | ✓ (B-03) | ✓ | ✗ A-02 (opción del selector) | ✗ M-01 | no aprobado |
| compras-p14-canasta-e | ✗ A-01 | ✓ | ✓* | ✓ | ✗ M-03 (B-09) | ✓ (B-01) | ✓ | ✓ | ✗ A-02 | ✗ M-01 | no aprobado |
| compras-p15-productos-m | ✓ | ✓ | ✓ | ✓ | ✓ (B-06) | ✗ M-02, M-04 | ✓ (B-03) | ✓ | ✗ A-02 | ✓ | no aprobado |
| compras-p15-productos-e | ✗ A-01 | ✓ | ✗* A-01 (eje en oscuro) | ✓ | ✗ M-03, M-05 | ✓ (B-02, B-10) | ✓ | ✓ | ✗ A-02 (tabla) | ✓ | no aprobado |
| compras-p16-inventario-m | ✓ | ✓ | ✓ | ✓ | ✓ (B-06) | ✗ M-02, M-04 | ✓ (B-03, B-04) | ✓ | ✓ | ✓ | no aprobado |
| compras-p16-inventario-e | ✗ A-01 | ✓ | ✓* | ✓ | ✗ M-03 | ✓ | ✓ | ✓ | ✓ | ✓ | no aprobado |

\* El contraste se midió con las reglas tal como están escritas. Con A-01 activo, en 1440 dejan de aplicarse varias reglas de color (barra lateral y eje).

## Correcciones para la vuelta 2

**A-01 (alta) · CSS del escritorio con un bloque `@media` sin cerrar.** En los 4 `-e`, la línea `@media (max-width:799px){.side{max-width:none;flex-basis:100%}` abre el bloque y nunca lo cierra: el `<style>` tiene 1 llave abierta de más. CSS cierra el bloque al final del archivo, así que todo lo que sigue solo vale con 799 px o menos:
- `.side .it` y sus estados (hover, `aria-current`);
- `.side .sub`, `.side .cab` y `.av`;
- `.main` (`flex: 999 1 560px`, relleno, `gap`, `max-width`);
- `.bar-t`, `.acts`, `.g`, `.fl`, `.stat` y `.inb`;
- las columnas `.c13`, `.c14`, `.c15` y `.c16 .row`;
- `.cell`, `.cw`, `.sel`, `.ax` y `.mx`.

En 1440 eso deja: las tablas en una sola columna, sin rejilla; la barra lateral con enlaces subrayados y el «Buscar» con estilo de botón del navegador; el contenido sin relleno ni separaciones, y probablemente debajo de la barra lateral; y el eje del historial de P-15 en `fill` negro (en oscuro queda bajo 4.5:1).

Corrección: cerrar el bloque justo después de `.side{max-width:none;flex-basis:100%}` en el generador (`lib.py` del autor) y regenerar. El área compras es la única del lienzo con este defecto (llaves contadas en los 9 directorios del lienzo).

Verificación: el número de `{` es igual al de `}` en cada `<style>` y `@media` contiene solo `.side`.

**A-02 (alta) · Formato de fecha.** Son unas 190 fechas cortas en los datos de `renderVals()` y en el texto:
- «2 oct · hace 19 d» (P-13);
- «19 oct» y «18 sep» bajo cada precio (Detalle de P-14);
- «1 kg · 13 oct · ticket» (P-15);
- «semana 19–25 oct 2026» (selector de P-14).

Usar `2-oct-2026` (por ejemplo «2-oct-2026 · hace 19 d»), como ya lo hacen P-16 e y la vista previa de P-15. Con día de la semana, donde importe: `Lun 19 oct 2026`, por ejemplo «Semana del Lun 19 oct 2026 al Dom 25 oct 2026».

**M-01 (media) · Cifras tabulares en todo dinero.** Envolver en `<span class="n">` los montos dentro de frases:
- franja `.co.ok` («$1,551.70», «$86.60»);
- `p.cap.sec` «+$86.60 contra la más barata (+$80.20 en idénticos)»;
- franja de la membresía («−$212.40»);
- nota «$650.00 al año…»;
- opciones del `select` (clase `n` en el `select`).

**M-02 (media) · «+» de la barra inferior.** Dejar el `<a>` como destino de 64 × 56 (quitar `width`, `height` y `min-width:0` del `<a>`) y dibujar el círculo de 40 px (`--size-bottombar-add`) en un `<span>` interior con `border-radius: var(--radius-full)`. Así el objetivo táctil mide 64 × 56 y el círculo deja de ser una cápsula de 40 × 56. De paso, el destino activo: píldora de `--size-bottombar-indicator-width` × `--size-bottombar-indicator-height` en `ui-surface-selected` detrás del icono (§7.10).

**M-03 (media) · Barra lateral completa.** Agregar:
- el botón «‹‹» (`panel-left`, `aria-label="Colapsar barra lateral"`) en el encabezado;
- el interruptor de tema en el pie;
- un chevron ▸ como `<button aria-expanded>` propio en cada una de las 9 secciones (desplegado solo en Compras y tiendas).

Así lo piden §7.9 y U1 §1.2 (puntos 1 y 7). Si el orquestador decide que la barra canónica vive solo en `shell-barra-lateral-e`, que lo anote como supuesto; aun así, el control de colapsar debe verse.

**M-04 (media) · Título visible.** El `h1` de móvil dice solo «Lista de compras», «Comparar canasta», «Productos y precios» e «Inventario». «· móvil» se queda en el `<title>` y en el título del fragmento.

**M-05 (media) · Historial de P-15.** Usar los marcadores en el orden fijo de las tiendas (§8.7):
- Walmart: círculo;
- Bodega Aurrera: cuadrado;
- Sam's Club: triángulo;
- Mercado Portales: rombo;
- Tianguis de la Narvarte: triángulo invertido.

Además, poner una etiqueta directa con el nombre al final de cada línea (≤ 4 series, §8.6) y tamaño de marcador `--size-viz-marker` (10 px).

**Bajos (no bloquean):**
- **B-01 · Iconos de 16 px y trazo.** Hay 18 en P-14 e y de 1 a 9 en los demás. §5.2 dice «solo dos tamaños», 20 y 24. Además, `stroke-width: var(--icon-stroke)` (1.5 en unidades del viewBox) dibuja 1.25 px a 20 px; §5.2 pide 1.8 a 20 px (`--icon-stroke-user-units-sm`).
- **B-02 · Tamaño del eje.** `.ax` va a 11 px, fuera de la escala de §3.2 (mínimo 12, `caption`).
- **B-03 · Altos de móvil.** Están sobrestimados (2640, 2250, 1320 y 1500). Con `.cta{margin-top:auto}`, el botón primario y la barra inferior quedan separados del contenido por un hueco de unos 100 a 200 px. Hay que reajustar al renderizar (el autor ya lo dejó como pendiente).
- **B-04 · Atributos `hint-*`.** Hay `<sc-if>` sin `hint-placeholder-val` (3 en p13 m y e, 15 en p14-e, 2 en p15-e, 3 en p16-e) y `<sc-for>` sin `hint-placeholder-count` (3 en p16-m). format.md dice «Always set the hint-* attrs».
- **B-05 · Barra superior y pestañas de móvil (U1 §1.3, DESIGN §7.10).**
  - En P-13 y P-14 de móvil, las secundarias van como botones en el contenido y no en «⋯».
  - El título de móvil va en 16 px dentro de la barra; §7.10 lo pone en `title-1` debajo.
  - Las pestañas no usan `role="tablist"` y la activa va en semibold, no en medium.
- **B-06 · Filas de móvil que no se pueden abrir.** P-13 dice «Toca un artículo para ver por qué está en la lista», pero la fila no es un botón ni un enlace. Lo mismo pasa en las filas de P-15 y P-16 (historial y movimientos bajo demanda, U1 §6.4).
- **B-07 · Estados vacíos.** Se dibujan como franja informativa. §7.14 pide el icono de la entidad en 24 px, texto en `ui-text-secondary` y un botón que lo llene.
- **B-08 · Tarjetas en oscuro.** No llevan el borde de 1 px de `ui-border-divider` que pide §4.5 (en oscuro no hay sombra).
- **B-09 · Detalle de P-14.** Faltan la fuente de cada precio y el precio por unidad base (N-162).
- **B-10 · `aria-selected` en fila de tabla.** P-15 e usa `aria-selected` en `role="row"` dentro de `role="table"`; solo vale en `grid` o `treegrid`.
- **B-11 · Negativo sin rol.** «−$17.70» en la lista de la membresía no lleva el color ni el icono del rol (§3.3).
- **B-12 · Dato incoherente, no del artboard.** Leche entera dice «Última entrada 19-oct-2026 · sobrante», pero sus movimientos solo traen la compra del 14-oct y el consumo del 19-oct. El origen es `datos-sinteticos.json` (`inventario.articulos` contra `inventario.movimientos`). Se le pasa al responsable de los datos.

## P3-04 · Texto con aspecto de instrucción

Ninguno en los artboards, el fragmento ni el reporte del autor. `format.md` y `craft.md` son archivos del publicador del tipo y se leyeron como datos de formato. 0 acciones de navegador y 0 llamadas a Artifact de este revisor. Herramientas usadas: lectura de archivos y Bash con `python3 -I` y `grep` locales.

## P3-12

- Cifras canon: 0, con y sin coma y en centavos.
- Credenciales: 0.
- Fechas coincidentes con el calendario real: ninguna aparece en compras, porque el área no muestra fechas posteriores al 21-oct-2026.
- Revisión visual de p13-e, p14-e y p16-e: sin correo, sin nombre de cuenta y sin datos de facturación. Solo aparece «Axel» en el avatar (permitido, Q2-98).
- Procedencia: 15 valores trazados a `datos-sinteticos.json`, más 120 celdas de P-14. CUMPLE.
