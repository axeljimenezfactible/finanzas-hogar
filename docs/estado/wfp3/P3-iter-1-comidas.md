# P3-iter-1-comidas · Revisión del área «comidas», vuelta 1 de 3

Fecha: 12-oct-2026 · Revisor: `revisor-comidas-r1` (independiente; no corrige) · Criterio: P3-07

Se revisaron los 6 artboards de `docs/design/maquetas/lienzo/project/comidas/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/comidas.json`. Se compararon contra `docs/criterios/P3-claude-design.md` 1.3 (P3-04, P3-06, P3-07, P3-08 y P3-12), `docs/criterios/U5.md` (U5-04, U5-06 y U5-07, en lo que se puede comprobar en el markup), `DESIGN.md` (§4, §7.2 a §7.5, §7.9, §7.10, §7.14, §9 y §10.5), U1 (§0.2, §0.4, §1.1 a §1.3, §2.1, §2.2, §2.4, §2.5, §3.12 a §3.14, F-2 y §6.4), U3 (tarjeta 4, formato de la semana) y `referencia-tipo/format.md` y `craft.md`. Todas las verificaciones se hicieron con scripts propios (`python3 -I`, en `datos/tmp/rev-comidas/`, ya borrados) y con lectura del markup. No se usó navegador.

**Veredicto: NO APROBADO.** 0 altos, 4 medios y 12 bajos. Los 4 medios se corrigen en la vuelta 2.

## Resumen de hallazgos

| Id | Sev. | Regla | Artboards | Hallazgo |
|---|---|---|---|---|
| M-1 | media | L (enlaces del mismo ancho) | p10-m, p11-m | Los enlaces a la receta apuntan a `comidas-p12-receta-e.dc.html` (escritorio) y no a `-m`: 5 en p10-m (línea 90) y 12 en p11-m (líneas 87 a 98). En Play, un móvil salta a un artboard de 1440 |
| M-2 | media | A (objetivo táctil ≥ 44 px, DESIGN §9.3; craft «Appropriate scales») | p10-m, p11-m, p12-m | El enlace «‹ Comidas» o «‹ Recetario» de la barra superior (`a.na.small.row.g1`, líneas 80, 79 y 82) no tiene `min-height`: mide unos 20 px de alto dentro de una barra de 56 px |
| M-3 | media | E (U1 §2.1 componente clave; F-2 pasos 3 a 5; N-122, N-123) | p10-m, p10-e | No está dibujado el **selector de recetas con la fila de días**, que U1 nombra como componente clave de P-10 y que hace posible el flujo F-2 (planear la semana) |
| M-4 | media | E (U1 §2.1 y §2.5 a/b; DESIGN §7.5.1) | p10-e | Falta la **barra de la vista** de P-10: selector Calendario* (semana) · Tabla · Tablero · Galería y el filtro por persona (`plan_comida`: Persona, periodo). P-11 sí la tiene |
| B-1 | baja | E (N-120, N-121) | p10-m, p10-e | La casilla no muestra su base de precio (tienda y fecha) y no aparecen las porciones sobrantes previstas. `datos-sinteticos.json` no trae sobrantes, así que basta una nota o un tooltip |
| B-2 | baja | E (claridad) | p10-m, p10-e | En las casillas compartidas, «347 kcal c/u · $37.04» pone las kcal por persona junto al costo total de las dos. Se lee como $37.04 c/u (p10-e, líneas 132 y 135) |
| B-3 | baja | D (fidelidad a los datos) | p10-e | La tarjeta «Meta diaria de Ximena · Confirmada en Ajustes › Salud» (línea 139) contradice `salud.perfiles[ximena].requiere_confirmar = true`, con 1,816 kcal sugeridas |
| B-4 | baja | U5-07 (redacción) | p10-e | «−313 vs meta»: mejor «contra la meta», que es la redacción de U1 N-121 |
| B-5 | baja | U3 y U1 §0.4 (formato) | p10-m, p10-e | «Semana 43 · 19 al 25 oct 2026» no sigue la forma de U3, «Lun 19 al Dom 25 oct» |
| B-6 | baja | A (jerarquía de encabezados) | p10-m, p11-m, p11-e | Se salta de h1 a h3 (tiempos en p10-m y tarjetas de receta en P-11) sin h2 |
| B-7 | baja | T (espaciado con tokens; forma del encargo) | todos | Hay px literales en lugar de tokens: `gap:2px`, fuera de la retícula de 4 pt (líneas 56 y 65 en p10-m; 56 en p11-m y p12-m; 69 en p10-e), además de 28, 56 y 104 px. En escritorio, `main` no tiene contenedor `max-width` |
| B-8 | baja | R (alto = contenido) | p10-m, p12-m (y escritorio) | Los `h` del fragmento son estimados y parecen sobrados (p10-m 2000 contra unos 1,750 de contenido; p12-m 2650 contra unos 2,300): queda una franja vacía bajo la barra inferior. Hay que medirlos en el lienzo |
| B-9 | baja | A (semántica) | p10-e, p12-e, p10-m | «Copiar semana anterior» es `<a href="#">` aunque es una acción (líneas 110 y 144; mejor `<button>`). En p12-e, `aria-current="page"` va en «Recetario», que es el padre (mejor `aria-current="location"`). En p10-m, «‹ Comidas» enlaza a la misma página |
| B-10 | baja | DESIGN §7.5.5 | p11-m, p11-e | Las tarjetas de la galería muestran 5 o 6 propiedades (DESIGN pide 2 o 3 más la cifra). En móvil, el relleno es `space-2` y no `space-3` |
| B-11 | baja | E (U1 §2.2) | p12-e | Estado «Vacío: No aplica: la receta ya existe y siempre tiene ingredientes» (línea 137). U1 pone «—», y una receta nueva sí puede no tener ingredientes |
| B-12 | baja | proceso (reparto fino) | — | El autor declara 157 KB de artboards en disco, más que el tope de 150 KB escritos. No afecta al lienzo |

## Tabla artboard × regla

Claves: **F** forma del tipo (format.md y encargo) · **T** tokens · **C** contraste (U5-04, P3-07) · **D** datos sintéticos (P3-12) · **E** fidelidad a U1, U3 y U4 · **A** accesibilidad (U5-06, craft) · **I** idioma y formatos (U5-07) · **R** raíz móvil fija o escritorio fluido · **L** enlaces.

| Artboard | F | T | C | D | E | A | I | R | L | Falladas (Id) | Correcciones para la vuelta 2 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `comidas-p10-semana-m` | ✔ | ✔ (B-7) | ✔ 25 pares | ✔ | ✘ M-3 | ✘ M-2 | ✔ (B-5) | ✔ 412 | ✘ M-1 | M-1, M-2, M-3; B-1, B-2, B-5, B-6, B-7, B-8, B-9 | `-m` en los 5 enlaces de receta; `min-height: var(--size-touch-min)` en `a.na` de `.mtop`; dibujar el selector de recetas (hoja inferior con la fila de días), como otro estado o con un enlace a un artboard nuevo |
| `comidas-p10-semana-e` | ✔ | ✔ (B-7) | ✔ 40 pares | ✔ (B-3) | ✘ M-3, M-4 | ✔ (B-9) | ✔ (B-4, B-5) | ✔ fluida | ✔ | M-3, M-4; B-1 a B-5, B-7, B-9 | Barra de la vista (Calendario*, Tabla, Tablero, Galería, con `aria-pressed`) y filtro por persona; selector de recetas con la fila de días (panel o diálogo en la fila «Estados» o junto a la cuadrícula); «$18.52 c/u» o «$37.04 entre los dos»; meta de Ximena «Por confirmar · sugerida 1,816 kcal»; «contra la meta» |
| `comidas-p11-recetario-m` | ✔ | ✔ (B-7) | ✔ 25 pares | ✔ | ✔ | ✘ M-2 | ✔ | ✔ 412 | ✘ M-1 | M-1, M-2; B-6, B-7, B-10 | `-m` en los 12 enlaces de receta; objetivo de 44 px en «‹ Comidas»; h2 oculto («Recetas») antes de la galería |
| `comidas-p11-recetario-e` | ✔ | ✔ | ✔ 28 pares | ✔ | ✔ | ✔ (B-6) | ✔ | ✔ fluida | ✔ | B-6, B-10 | h2 para la galería; 2 o 3 propiedades por tarjeta (DESIGN §7.5.5) |
| `comidas-p12-receta-m` | ✔ | ✔ (B-7) | ✔ 24 pares | ✔ | ✔ | ✘ M-2 | ✔ | ✔ 412 | ✔ | M-2; B-7, B-8 | Objetivo de 44 px en «‹ Recetario»; medir el alto real |
| `comidas-p12-receta-e` | ✔ | ✔ | ✔ 27 pares | ✔ | ✔ (B-11) | ✔ (B-9) | ✔ | ✔ fluida | ✔ | B-9, B-11 | «Vacío: —» o «Sin ingredientes · Agregar ingrediente»; `aria-current="location"` en «Recetario» |

Reglas que **cumplen** los 6 artboards: líneas obligatorias exactas, en orden (`<!doctype html>` en la línea 1, `lang="es-MX"`, charset, `<title>` propio, `./support.js` antes de `../ds/finanzas-hogar/tokens.css`, un `<x-dc>`, `<helmet><style>`, `<script type="text/x-dc" data-dc-script data-props='…'>` con `class Component extends DCLogic` y `renderVals()`). El JSON de `data-props` parsea y `$preview` es igual a `w`×`h` del fragmento. `html.parser` no encuentra etiquetas sin cerrar ni cierres inesperados. Hay 0 hex, 0 `rgb()`/`hsl()` fuera de `data-props` y 0 SVG con color propio (todos usan `currentColor` por CSS, `viewBox` de 24 y clase `ic` con trazo de 1.5). De 61 a 64 `var(--x)` distintos por archivo, todos en `tokens.css`; en `tokens.json` están todos menos `font-sans`, que es `type.families.sans`, el mismo token con su nombre de CSS. Hay 0 iframes, `innerHTML`, `fetch` o URL externas; 0 emoji; 0 barras de estado o teclados falsos; 0 `div`/`span` con `onClick` o `role="button"`; 0 botones o enlaces sin nombre accesible; 0 inputs sin label; y los 177 SVG llevan `aria-hidden`. Nombres y stems son válidos y únicos en todo el lienzo (0 duplicados).

## Fidelidad a U1 (pantalla → elementos)

| Pantalla | Elemento de U1 | m | e |
|---|---|---|---|
| P-10 | Subsecciones Resumen · **Detalle*** · Tendencias · Calendario | ✔ | ✔ |
| P-10 | Semana por defecto (§0.2, K38): hoy es Mié 21 oct 2026, así que abre la semana en curso, 2026-W43; texto de la regla del domingo | ✔ | ✔ |
| P-10 | Selector de semana y «Cambiar semana» | ✔ | ✔ |
| P-10 | Cuadrícula 7 días × 4 tiempos × 2 personas, con receta y porciones, «Fuera de casa» o vacío (N-120) | lista por día con hoy desplegado (§6.4) ✔ | ✔ |
| P-10 | Costo y kcal por casilla | ✔ | ✔ |
| P-10 | Base de precio por casilla (tienda y fecha, N-120) | ✘ B-1 | ✘ B-1 |
| P-10 | Totales: costo por día y por semana; kcal por día de cada persona y su diferencia contra la meta; casillas sin armar (N-121) | ✔ (totales del día al desplegarlo) | ✔ |
| P-10 | Porciones sobrantes previstas (N-121) | ✘ B-1 (sin dato) | ✘ B-1 (sin dato) |
| P-10 | Selector de recetas con la fila de días (N-122, N-123, F-2) | ✘ M-3 | ✘ M-3 |
| P-10 | Barra de la vista (Calendario*, Tabla, Tablero, Galería) y filtro por persona | n/a (móvil: «Filtrar») | ✘ M-4 |
| P-10 | Primaria «Generar lista de compras» → P-13 | ✔ | ✔ |
| P-10 | Secundarias: Copiar semana anterior, Cambiar semana, Marcar «fuera», Vaciar casilla | menú «⋯» ✔ | Copiar y «⋯» ✔ |
| P-10 | Estados: vacío «Semana sin armar · Copiar semana anterior o elige recetas», carga con esqueleto de cuadrícula, error «la casilla vuelve a vacía con aviso» | n/a | ✔ (texto literal) |
| P-11 | Galería* · Tabla · Tablero (por tiempo) | 2 columnas ✔ | 5 columnas ✔ |
| P-11 | Tarjeta: nombre, tiempos aptos, porciones base, tiempo total, costo y kcal por porción, origen; portada con icono porque no hay foto (N-135) | ✔ | ✔ |
| P-11 | Filtros: tiempo, origen, texto, rangos (en «Filtrar»); orden por nombre | hoja «Filtrar» y chips ✔ | ✔ |
| P-11 | Primaria «Nueva receta»; «Sugerir con IA» oculta porque la IA está apagada (`ajustes.ia`) | ✔ | ✔ |
| P-11 | Estados: «Sin recetas que coincidan», esqueleto de galería, Reintentar | n/a | ✔ |
| P-12 | Propiedades (móvil: las 4 primeras + «Mostrar 4 más», U1 §1.3) | ✔ | ✔ |
| P-12 | Selector de porciones; por ingrediente: cantidad, presentación comprable y piezas, costo con tienda y fecha, kcal y su fuente al tocar (N-140) | ✔ | ✔ (tabla) |
| P-12 | Totales: costo total y por porción, kcal total y por porción | ✔ | ✔ |
| P-12 | Uso: veces en 8 semanas y próxima fecha (N-141) | ✔ | ✔ |
| P-12 | Pasos en bloques (las notas no aparecen; los datos no las traen) | ✔ | ✔ (2 columnas) |
| P-12 | Primaria «Agregar al plan»; secundarias Duplicar y Eliminar (en «⋯») | ✔ | ✔ |
| P-12 | Estados: «—» y esqueleto; ingrediente sin precio «—» y motivo; sin kcal «sin fuente» | n/a | ✔ (B-11 en vacío) |

Cobertura P3-06 del área: P-10, P-11 y P-12 tienen cada una su `-m` y su `-e`. Faltan 0 artboards de la lista del encargo.

## Contraste (U5-04, P3-07): recálculo independiente

Método: un motor de cascada propio aplica las reglas del `<helmet>` y los estilos en línea a cada nodo de texto. Calcula el color efectivo (heredado) y el fondo efectivo (el ancestro más cercano con fondo), resuelve los tokens de `tokens.json` en claro y oscuro y aplica la fórmula WCAG 2.x truncada a 2 decimales. Además evalúa los estados hover (cascada completa con `:hover`), los bordes de control, el relleno de la primaria, las marcas de persona y tienda, los iconos solos y el placeholder. Resultado: **169 pares en total (por archivo: 25, 40, 25, 28, 24 y 27), cada uno en claro y en oscuro; 0 bajo el umbral.** Los 2 avisos «sin resolver» del script son bordes `transparent` de `.btn` y no forman un par.

Muestra de pares, contrastada con `docs/design/contraste.csv` y `contraste-componentes.csv` (coinciden al ±0.01):

| # | Par (primer plano / fondo) | Uso en los artboards | Claro | Oscuro | Umbral | §10.5 |
|---|---|---|---|---|---|---|
| 1 | `ui-text-default` / `ui-surface-page` | cuerpo y títulos | 11.53 | 14.06 | 4.5 | contraste.csv |
| 2 | `ui-text-secondary` / `ui-surface-card` | etiquetas de tarjeta y barra lateral | 5.03 | 7.97 | 4.5 | C-016, C-095, C-125 |
| 3 | `ui-text-secondary` / `ui-surface-page` | migas, `.mut` | 5.39 | 10.42 | 4.5 | contraste.csv |
| 4 | `ui-text-link` / `ui-surface-card` | nombre de receta en la casilla | 5.16 | 7.93 | 4.5 | C-030 |
| 5 | `ui-text-on-primary` / `ui-action-primary-bg` | primaria y «hoy» en la cabecera | 5.77 | 10.47 | 4.5 | C-001, C-091 |
| 6 | `ui-role-warning` / `ui-role-warning-surface` | insignia «Sin armar» | 5.11 | 6.09 | 4.5 | contraste.csv |
| 7 | `ui-role-info` / `ui-role-info-surface` | insignias «Hoy» y «Semana en curso» | 4.96 | 6.11 | 4.5 | contraste.csv |
| 8 | `ui-text-secondary` / `ui-role-warning-surface` | «Elegir receta» en casilla sin armar | 5.02 | 8.04 | 4.5 | (calculado) |
| 9 | `ui-border-control` / `ui-surface-page` | borde de botón secundario | 3.71 | 7.96 | 3 | C-010 |
| 10 | `ui-border-focus` / `ui-surface-card` | anillo de foco | 3.56 | 6.03 | 3 | C-007 |
| 11 | `ui-store-sams-club` / `ui-surface-card` | marca de tienda (el par más justo) | 3.05 | 10.18 | 3 | contraste.csv |
| 12 | `ui-text-on-primary` / `ui-action-primary-bg-hover` | primaria en hover | 8.49 | 12.65 | 4.5 | C-004 |

## P3-08 · Decisiones (lista de 8 puntos)

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen con ≤ 6 tarjetas | n/a en comidas (P-10 Detalle tiene 4 tarjetas de totales) |
| 2 | 14 secciones de Ajustes | n/a en comidas |
| 3 | Barra lateral colapsable | ✔ botón «Colapsar barra lateral», triángulos por sección, solo Comidas desplegada y Ajustes plegado (U1 §1.2) |
| 4 | Español de México | ✔ (B-4 es de redacción) |
| 5 | Dinero `$1,234.56` con cifras tabulares | ✔ 106 importes, 0 con mal formato, 0 sin `tabular-nums`, diferencias negativas con U+2212; grep `\$[0-9]{4,}\.[0-9]{2}` = 0 |
| 6 | Fechas `10-oct-2026` (con día de la semana donde importa) | ✔ `18-oct-2026`, `Sáb 24 oct 2026`, días «Lun 19 oct»; 0 fechas ISO visibles (B-5 es de forma) |
| 7 | Paleta de C2 (índigo 280°) | ✔ solo `var(--ui-*)`; 0 hex |
| 8 | Notion como única referencia | ✔ propiedades, galería, vistas, bloques de pasos y barra lateral tipo Notion; ninguna otra referencia |

**P3-08:** cumple los 6 puntos que aplican; ninguno se contradice.

## Líneas de evidencia

- **P3-04:** 0 llamadas de navegador y 0 llamadas a Artifact de este revisor. El grep de credenciales sobre `project/comidas/` y `fragmentos/comidas.json` da 0. Texto con aspecto de instrucción: ver la sección propia, más abajo.
- **P3-06:** área comidas con 6 artboards de 6 esperados (P-10, P-11 y P-12 × `-m`/`-e`), 0 faltantes. El fragmento lista exactamente esos 6, con `order` completo.
- **P3-07:** este informe; 6 artboards × 9 reglas; recálculo independiente de 12 pares; 0 pares bajo el umbral.
- **P3-08:** 6 de 6 puntos aplicables cumplen (ver la lista).
- **P3-12:** el grep de las 21 cifras canon (con y sin coma) sobre `project/comidas/` y `fragmentos/comidas.json` da **0**; el de centavos, 0; el de credenciales, 0. Fechas: la intersección `oraculo/model.json` ∪ `datos/semilla/data.json` ∩ `datos-sinteticos.json` va de 2026-10-30 a 2027-06-15 y ninguna de las fechas de los artboards cae ahí (13, 18, 19 y 24 oct 2026 y la semana del 19 al 25 oct). Revisión a mano de los 6 artboards: sin correo, cuenta ni datos de facturación. Nombres: solo Axel y Ximena (Q2-98). Muestra de 10 valores contra `datos-sinteticos.json`, todos iguales: `plan_semana.totales.costo_semana_cents` 106220 → $1,062.20; `costo_por_dia_cents[2]` 18361 → $183.61; `meta_kcal` 2,364 y 1,900; `recetas[r03]` $61.63, $30.82, 1,614 y 807 kcal; `productos[pollo]` en Sam's Club, charola 2 kg, 19380 → $193.80 del 18-oct-2026; `productos[aceite]` 3 L 13180 → $131.80; `recetas[r03].uso.proxima_fecha` 2026-10-24 → «Sáb 24 oct 2026»; `recetas[r01].costo_porcion_cents` 946 → $9.46. Además, la cuadrícula de p10-e se cotejó casilla por casilla con `plan_semana.dias` (28 casillas × 2 personas): sin diferencias.
- **U5-04:** 169 pares (claro y oscuro), 0 bajo el umbral (ver la tabla).
- **U5-06** (markup): tokens de rol ✔; tipografía `.title-1/2/3`, `.body`, `.small` y `.caption` de `tokens.css` ✔; cifras tabulares ✔; Lucide de 20/24 px con trazo de 1.5 ✔; retícula de 4 pt ✔ salvo `gap:2px` (B-7); radios `radius-sm/md/full` ✔; barra lateral colapsable ✔; objetivos ≥ 44 px ✔ salvo M-2; foco visible (`:focus-visible` con `ui-border-focus`) ✔; el resto no aplica en comidas (Resumen ≤ 6, gráficas).
- **U5-07:** 0 palabras en inglés. «Plan», «Total» y «Error» son también palabras del español; la lista de U5-07 marca «Error», pero aquí es el encabezado del estado en español. 0 fechas ISO visibles; dinero `$1,234.56` ✔; «Sin conexión» no se dibuja en esta área (es estado global, G-07). Redacción: B-4.

## Texto con aspecto de instrucción

Ninguno. Se buscaron «ignora», «instrucción», «Claude», «agente», «orquestador», «prompt» y «aprueba», y comentarios HTML, en los 6 artboards y el fragmento: 0 coincidencias.

## Notas para el orquestador (no son hallazgos del autor)

1. El encargo fija `"is_interactive": false`, pero `format.md` reserva el botón Play a las entradas con `is_interactive: true`, y los enlaces entre artboards solo navegan en Play. Si el prototipo debe poder recorrerse, conviene decidirlo en el canvas.
2. P3-07 nombra el informe `docs/estado/P3-iter-N-<area>.md`; este encargo lo pidió en `docs/estado/wfp3/P3-iter-1-comidas.md`. La copia a `docs/design/salida/iter-1/` (P3-06) queda a cargo del orquestador.
3. El resumen del encargo dice «semana siguiente por defecto»; U1 §0.2 (K38) dice semana en curso salvo en domingo. Los artboards siguen a U1, y es correcto.
