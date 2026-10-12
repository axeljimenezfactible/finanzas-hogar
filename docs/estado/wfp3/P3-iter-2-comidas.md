# P3-iter-2-comidas · Revisión del área «comidas», vuelta 2 de 3

Fecha: 12-oct-2026 · Revisor: `revisor-comidas-r2` (independiente; no corrige) · Criterio: P3-07

Revisé los 6 artboards de `docs/design/maquetas/lienzo/project/comidas/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/comidas.json`, después de la corrección `correccion-comidas-r1`. Los comparé con `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06 y U5-07, en lo que se ve en el markup), `DESIGN.md` (§7.5.1, §7.12, §9.3, §10.5), U1 (§2.1, §2.2, §2.5, §6.2, §6.4, N-120 a N-123, F-2) y `referencia-tipo/format.md` y `craft.md`. Todo se verificó con scripts propios (`python3 -I`, en `datos/tmp/rev-comidas-r2/`, borrados al terminar) y leyendo el markup. No usé navegador ni Artifact.

Alcance de la vuelta: comparé los archivos con la copia de la vuelta 1 (`docs/design/salida/iter-1/comidas/`). `p11-e` y `p12-e` no cambiaron. `p11-m` y `p12-m` solo cambian en la clase `bk` y en los `href`. `p10-m` y `p10-e` suman el selector de recetas; `p10-e` suma además la barra de la vista.

**Veredicto: APROBADO.** 0 altos y 0 medios. Los 4 medios de la vuelta 1 están cerrados. Siguen abiertos 10 bajos de la vuelta 1 y hay 6 bajos nuevos; ninguno impide aprobar.

## Hallazgos de la vuelta 1: cerrados o reabiertos

| Id r1 | Sev. | Estado | Evidencia |
|---|---|---|---|
| M-1 | media | **Cerrado** | 0 enlaces a `-e` en los artboards `-m` y 0 a `-m` en los `-e`. p10-m tiene 5 `href` a `comidas-p12-receta-m.dc.html` y p11-m tiene 12. p12-m enlaza a `comidas-p11-recetario-m.dc.html` |
| M-2 | media | **Cerrado** | El enlace «‹ Comidas» o «‹ Recetario» lleva `class="na bk small row g1"`, con `.bk{min-height:var(--size-touch-min)}` (44 px) y `display:flex` de `.row` (p10-m línea 91, p11-m línea 80, p12-m línea 83). Medí todos los controles de los 3 móviles (96) con la cascada y los tokens: 0 bajo 44 px. Las 12 `.rn` de p11-m se descartan porque `.rn::after{inset:0}` cubre la tarjeta `.gc{position:relative}` |
| M-3 | media | **Cerrado** | p10-m (hoja, líneas 106 a 114) y p10-e (tarjeta, líneas 153 a 159) tienen: título «Elegir receta · Lun · Comida»; fila de 7 días con Lun marcado (`aria-pressed="true"`); persona con «Ambos» marcado; porciones por persona en 1; búsqueda con `<label>`; nota del orden. Las 5 recetas aptas para Comida aparecen ordenadas por `veces_8_semanas` y luego por nombre, con costo y kcal por porción. Las cotejé con `recetas[]` y salen iguales: $30.82/807/6, $27.72/639/4, $14.23/831/2, $34.38/743/1, $12.43/247/1 |
| M-4 | media | **Cerrado** | p10-e, línea 139: grupo «Vista del plan» con Calendario (`aria-pressed="true"`, icono calendar), Tabla (table-2), Tablero (kanban) y Galería (layout-grid). También tiene «Filtrar» y el filtro Persona (Todas, Axel y Ximena). El periodo (semana) está en la sección «Semana elegida», justo encima. Esto coincide con U1 §2.5 (`plan_comida`: Persona, periodo) |
| B-1 | baja | Abierto | Sin base de precio por casilla ni porciones sobrantes (N-120, N-121) |
| B-2 | baja | Abierto | «347 kcal c/u · $37.04» (el costo es de los dos) |
| B-3 | baja | Abierto | p10-e, línea 137: Ximena «Confirmada en Ajustes › Salud», pero en `salud.perfiles` dice `requiere_confirmar: true` y la meta sugerida es 1,816 |
| B-4 | baja | Abierto | «−313 vs meta» (U1 dice «contra su meta») |
| B-5 | baja | Abierto | «Semana 43 · 19 al 25 oct 2026» |
| B-6 | baja | Abierto | Los encabezados saltan de h1 a h3 en p10-m (orden 1,3,3,3,3,3,2), p11-m y p11-e |
| B-7 | baja | Abierto | `gap:2px` en 5 lugares; px literales (56, 28, 104); sin contenedor `max-width` en `main` |
| B-8 | baja | Abierto | `h` estimados. Calculé a mano unos 2,500 px de contenido en p10-m (h 2700) y unos 1,950 en p10-e (h 2300). Hay que medirlos en el lienzo |
| B-9 | baja | Abierto | «Copiar semana anterior» sigue siendo `<a href="#">` en p10-e (líneas 120 y 162) |
| B-10 | baja | Abierto | La galería de P-11 muestra 5 o 6 propiedades por tarjeta |
| B-11 | baja | Abierto | p12-e: «No aplica: la receta ya existe y siempre tiene ingredientes» |
| B-12 | baja | Cerrado (no aplica) | En esta vuelta el autor escribió unos 14 KB |

## Hallazgos nuevos (vuelta 2)

| Id | Sev. | Regla | Artboards | Hallazgo | Corrección sugerida |
|---|---|---|---|---|---|
| N-1 | baja | E (U1 F-2, paso 5; N-123) | p10-m, p10-e | Para asignar hay que elegir y luego confirmar: la receta queda como alternado (`aria-pressed`) y después se pulsa «Poner en el plan · 2 porciones». En F-2, tocar la receta asigna y cierra el selector, así que esto suma 1 toque por receta (el ejemplo de U1 pasa de 30 a 37 toques) | Quitar el alternado y la primaria (tocar la receta asigna), o dejar la primaria solo para cuando cambian persona o porciones, y anotarlo |
| N-2 | baja | DESIGN §7.12 y U1 §6.2 | p10-m, p10-e | Hoja: el título usa `title-3` (DESIGN pide `title-2`) y la primaria mide 44 px (`size-control-coarse`), cuando DESIGN pide 48 (`size-control-primary-sheet`). Diálogo: la primaria va antes de «Cancelar», a la izquierda (DESIGN: acciones a la derecha, primaria al final). Además es una tarjeta a todo el ancho, no un diálogo centrado de 640 px | `title-2` y `min-height:var(--size-control-primary-sheet)` en la hoja; en escritorio, «Cancelar» y luego la primaria, alineadas a la derecha, y un ancho máximo de `layout-dialog-width-long` |
| N-3 | baja | DESIGN §9.3 (separación entre objetivos) | p10-m, p10-e | Entre objetivos de 44 px contiguos hay 4 px (`space-1`): fila de días (`.sdias`), persona (`row g1` en p10-m) y los grupos de vista y persona (`row g1` en p10-e). DESIGN pide `space-2` (8 px) como mínimo. El tamaño sí cumple 44 px | `gap:var(--space-2)` en esos grupos |
| N-4 | baja | D (coherencia con los datos) | p10-m, p10-e | El ejemplo «Lun · Comida» sigue a U1 F-2, pero en `datos-sinteticos.json` esa casilla no está vacía: Ximena tiene Pollo asado y Axel está fuera. Las vacías son Sáb · Cena, Dom · Cena y las colaciones de Mar y Jue | Abrir el selector desde «Sáb · Cena» con las recetas aptas para Cena, o anotar que se cambia una receta ya puesta |
| N-5 | baja | DESIGN §7.5.1 | p10-e (y p10-m) | En p10-e, «Filtrar» usa el icono `funnel`, pero DESIGN pide `list-filter`, que es el que usa P-11. Los botones de vista son secundarios con borde (`btn s`) y DESIGN pide botones de texto (pasa lo mismo en P-11). En móvil, P-10 no tiene «Filtrar» (Persona), aunque DESIGN §7.5.1 lo pide en móvil; U1 §6.4 no lo nombra | Icono `list-filter`; `btn t` con `aria-pressed`; decidir si P-10 móvil lleva «Filtrar» |
| N-6 | baja | T (forma) | p10-m, p10-e, p11-m | La clase `g2` se usa 5 veces y ninguna regla del `<helmet>` la define. No se nota porque `.row` y `.col` ya traen `space-2` | Definir `.g2{gap:var(--space-2)}` o quitar la clase |

## Tabla artboard × regla

Claves: **F** forma del tipo · **T** tokens · **C** contraste (U5-04) · **D** datos sintéticos (P3-12) · **E** fidelidad a U1 · **A** accesibilidad (U5-06, craft) · **I** idioma y formatos (U5-07) · **R** raíz móvil fija o escritorio fluido · **L** enlaces.

| Artboard | F | T | C | D | E | A | I | R | L | Bajos abiertos | Correcciones para la vuelta 3 (opcionales) |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `comidas-p10-semana-m` | ✔ | ✔ (B-7, N-6) | ✔ 40 | ✔ (N-4) | ✔ (B-1, B-2, N-1) | ✔ (B-6, N-2, N-3) | ✔ (B-5) | ✔ 412 | ✔ | B-1, B-2, B-5 a B-8, N-1 a N-6 | N-1, N-2 (hoja), N-3, N-4 |
| `comidas-p10-semana-e` | ✔ | ✔ (B-7, N-6) | ✔ 45 | ✔ (B-3, N-4) | ✔ (B-1, B-2, N-1, N-5) | ✔ (B-9, N-2, N-3) | ✔ (B-4, B-5) | ✔ fluida | ✔ | B-1 a B-5, B-7 a B-9, N-1 a N-6 | N-1, N-2 (diálogo), N-5, B-3 |
| `comidas-p11-recetario-m` | ✔ | ✔ (B-7, N-6) | ✔ 27 | ✔ | ✔ | ✔ (B-6) | ✔ | ✔ 412 | ✔ | B-6, B-7, B-10, N-6 | — |
| `comidas-p11-recetario-e` | ✔ | ✔ | ✔ 28 | ✔ | ✔ | ✔ (B-6) | ✔ | ✔ fluida | ✔ | B-6, B-10 | — |
| `comidas-p12-receta-m` | ✔ | ✔ (B-7) | ✔ 25 | ✔ | ✔ | ✔ | ✔ | ✔ 412 | ✔ | B-7, B-8 | — |
| `comidas-p12-receta-e` | ✔ | ✔ | ✔ 25 | ✔ | ✔ (B-11) | ✔ (B-9) | ✔ | ✔ fluida | ✔ | B-9, B-11 | — |

Lo que cumplen los 6 artboards (script `v1.py`):
- Las 14 piezas obligatorias de la forma, en orden: `<!doctype html>`, `<html lang="es-MX">`, charset, `<title>` propio, `./support.js` antes de `../ds/finanzas-hogar/tokens.css`, un solo `<x-dc>`, `<helmet>`, `data-dc-script` con `class Component extends DCLogic` y `renderVals()`.
- El JSON de `data-props` parsea y `$preview` es igual a `w`×`h` del fragmento.
- `html.parser`: 0 etiquetas sin cerrar y 0 cierres inesperados.
- 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`.
- De 61 a 66 `var(--x)` por archivo, todos en `tokens.json` (`font-sans` es `type.families.sans`) y en `tokens.css`.
- 0 iframes, `innerHTML`, `fetch`, URL externas y `{{`.
- 190 SVG Lucide, todos con `aria-hidden`, `viewBox` de 24, clase `ic` y `currentColor`.
- 0 ids duplicados, 0 `aria-labelledby` rotos y un solo h1 por archivo.
- 90 stems en el lienzo, 0 duplicados.
- El fragmento lista exactamente los 6 archivos, con `order` completo, `expand: "fill"` solo en escritorio, `is_interactive: false` y sin x ni y.
- Móvil: `.mroot{width:412px;min-height:915px}` con `box-sizing: border-box`. Escritorio: `.pg` en `flex-wrap`, `.side` en `flex:0 0 var(--layout-sidebar-width)` y `.main` en `flex:999 1 560px; min-width:0`. La tabla de 1040 px va dentro de `.scx{overflow-x:auto}`.

## Fidelidad a U1 (lo que cambió en esta vuelta)

| Pantalla | Elemento de U1 | m | e |
|---|---|---|---|
| P-10 | Selector de receta con fila de días (U1 §2.1): recetas aptas para el tiempo, orden por uso en 8 semanas y nombre, costo y kcal por porción, filtro de texto (N-122) | ✔ | ✔ |
| P-10 | Asignar (N-123): días, tiempo, persona con Ambos por defecto, porciones por persona en 1 («con Ambos son 2») | ✔ | ✔ |
| P-10 | F-2, pasos 3 a 5: la casilla abre el selector con su día marcado y tocar la receta asigna y cierra | ✔ con N-1 (paso extra) | ✔ con N-1 |
| P-10 | Barra de la vista: Calendario* · Tabla · Tablero · Galería (§2.5.a) y filtro Persona y periodo (§2.5.b) | n/a (U1 §6.4: lista por día; ver N-5) | ✔ |
| P-10 | Hoja inferior en 412 y diálogo centrado en 1440 (U1 §6.2) | ✔ dibujada en el flujo con nota | ✔ como tarjeta con nota (N-2) |
| P-10 | Estados de U1 §2.2: «Semana sin armar · Copiar semana anterior o elige recetas», esqueleto de cuadrícula y «La casilla que no se guardó vuelve a vacía con aviso» | n/a | ✔ (líneas 160 a 165) |
| P-10 | Cuadrícula (N-120): 28 casillas × 2 personas cotejadas con `plan_semana.dias` (receta, porciones, costo y kcal; «fuera»; vacío con «Elegir receta») | ✔ (Mié desplegado) | ✔ 0 diferencias |
| P-11, P-12 | Sin cambios de contenido desde la vuelta 1, que los aprobó. Volví a cotejar las 12 tarjetas de p11-m con `recetas[]` (costo y kcal por porción, minutos, porciones base, veces en 8 semanas): 0 diferencias | ✔ | ✔ |

Cobertura P3-06 del área: 6 artboards de 6 esperados. Faltan 0.

## Contraste (U5-04, P3-07): recálculo independiente

Método: escribí un motor de cascada propio (`cas.py`). Aplica las reglas del `<helmet>`, por especificidad y orden, y los estilos en línea a cada nodo de texto; hereda `color` y toma como fondo el del ancestro más cercano que lo tiene. Resuelve los tokens de `tokens.json`, siguiendo las referencias `{…}`, en claro y en oscuro, y usa la fórmula WCAG 2.x truncada a 2 decimales. Evalúa además el hover de cada `a`, `button` y `summary`, y los bordes de control, los rellenos de acción, persona y tienda, y los iconos. Resultado: **190 pares (p10-m 40, p10-e 45, p11-m 27, p11-e 28, p12-m 25, p12-e 25), cada uno en claro y en oscuro: 0 bajo el umbral.** El esqueleto (`ui-skeleton-on-card`, 1.11 y 1.44) queda fuera porque es decorativo y lleva `aria-hidden`.

Muestra de pares nuevos de esta vuelta (hoja, selector y barra de la vista), cotejados con `contraste-componentes.csv`:

| # | Par (primer plano / fondo) | Uso | Claro | Oscuro | Umbral | §10.5 |
|---|---|---|---|---|---|---|
| 1 | `ui-text-default` / `ui-surface-raised` | título y recetas de la hoja | 11.53 | 10.75 | 4.5 | C-058, C-119 (igual) |
| 2 | `ui-text-secondary` / `ui-surface-raised` | nota, orden y uso en la hoja | 5.39 | 7.97 | 4.5 | (calculado) |
| 3 | `ui-text-default` / `ui-surface-selected` | día, persona y receta elegidos | 10.73 | 10.90 | 4.5 | C-064, C-082 (igual) |
| 4 | `ui-text-secondary` / `ui-surface-selected` | costo y kcal de la receta elegida | 5.02 | 8.08 | 4.5 | C-085 (igual) |
| 5 | `ui-icon-selected` / `ui-surface-selected` | marca de la receta elegida | 5.14 | 8.04 | 3 | (calculado) |
| 6 | `ui-border-selected` / `ui-surface-raised` | borde del día y de la persona elegidos | 5.53 | 7.93 | 3 | (calculado) |
| 7 | `ui-border-control` / `ui-surface-raised` | borde de los botones ± y de los días | 3.71 | 6.08 | 3 | (calculado) |
| 8 | `ui-border-control` / `ui-surface-card` | botones de vista y «Filtrar» en p10-e | 3.46 | 6.08 | 3 | C-011 |
| 9 | `ui-action-primary-bg` / `ui-surface-raised` | relleno de «Poner en el plan» | 5.53 | 7.93 | 3 | (calculado) |
| 10 | `ui-text-on-primary` / `ui-action-primary-bg` | texto de la primaria | 5.77 | 10.47 | 4.5 | C-001 (igual) |
| 11 | `ui-text-secondary` / `ui-surface-field` | placeholder «Buscar receta» | 5.39 | 10.42 | 4.5 | (calculado) |
| 12 | `ui-person-axel` / `ui-surface-selected` | marca de Axel en «Ambos» elegido | 11.89 | 3.14 | 3 | (calculado) |

## P3-08 · Decisiones (lista de 8 puntos)

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen con ≤ 6 tarjetas | n/a en comidas |
| 2 | 14 secciones de Ajustes | n/a en comidas |
| 3 | Barra lateral colapsable | ✔ «Colapsar barra lateral»; Comidas desplegada |
| 4 | Español de México | ✔ 0 palabras en inglés (B-4 es de redacción) |
| 5 | Dinero `$1,234.56` con cifras tabulares | ✔ 116 importes, 0 con mal formato y 0 fuera de `.tn`. Negativos con U+2212 (8 en p10-e); 0 guiones ASCII ante cifras visibles |
| 6 | Fechas `10-oct-2026` | ✔ `13-oct-2026`, `18-oct-2026`, `19-oct-2026` y «Sáb 24 oct 2026»; 0 ISO visibles (B-5 es de forma) |
| 7 | Paleta de C2 (índigo 280°) | ✔ solo `var(--ui-*)`; 0 hex |
| 8 | Notion como única referencia | ✔ vistas, propiedades, hoja y galería al estilo Notion |

**P3-08:** cumple los 6 puntos que aplican; ninguno se contradice.

## Líneas de evidencia

- **P3-04:** 0 llamadas de navegador y 0 a Artifact. El grep de credenciales sobre `project/comidas/` y `fragmentos/comidas.json` da 0.
- **P3-06:** 6 de 6 artboards. La copia de la vuelta 1 está en `docs/design/salida/iter-1/comidas/` (6 artboards y `comidas.json`). La copia de este informe a `docs/design/salida/iter-2/` queda a cargo del orquestador.
- **P3-07:** este informe; 6 artboards × 9 reglas; 190 pares en los dos temas; 0 bajo el umbral.
- **P3-08:** 6 de 6 puntos aplicables cumplen.
- **P3-12:**
  - El grep de las 21 cifras canon (con y sin coma) sobre `project/comidas/` y `fragmentos/comidas.json` da **0**; el de esas mismas cifras en centavos, 0; el de credenciales, 0.
  - Fechas: el comando de P3-12 da 15 fechas comunes, de 2026-10-30 a 2027-06-15. Ninguna aparece en los artboards (13, 18, 19 y 24 oct 2026 y la semana del 19 al 25 oct).
  - Nombres: solo Axel y Ximena (Q2-98).
  - Trazabilidad: todos los montos y kcal de los 6 artboards salen de `datos-sinteticos.json`. Los únicos importes que no aparecen tal cual son $31.92 y $48.54 de p10-e: son sumas de casillas compartidas (2 × 1596 y 2 × 2427), y el cotejo casilla por casilla da 0 diferencias.
  - Muestra de 10 valores, todos iguales a los datos:
    - `plan_semana.totales.costo_semana_cents` 106220 → $1,062.20
    - `costo_por_dia_cents[2]` 18361 → $183.61
    - `kcal_por_dia.axel[4]` 3801 → 3,801
    - `meta_kcal` 2,364 y 1,900
    - `recetas[r03]` 3082 → $30.82, 807 kcal, 6 veces
    - `recetas[r04]` 2772 → $27.72, 639 kcal, 4 veces
    - `recetas[r01]` 946 → $9.46, 404 kcal
    - Mié · Cena de Axel: 2 porciones, 6218 → $62.18, 1,662 kcal
    - Mié · Comida de Ximena: picadillo, 3438 → $34.38, 743 kcal
    - `recetas[r05].tiempo_total_min` → 100 min
- **U5-04:** 190 pares en claro y oscuro (móvil y escritorio, es decir, las 4 combinaciones), 0 bajo el umbral.
- **U5-06** (markup): tokens de rol ✔; tipografía de `tokens.css` ✔; cifras tabulares ✔; Lucide de 20/24 px con trazo de 1.5 ✔; retícula de 4 pt ✔ salvo `gap:2px` (B-7); radios ✔; `shadow-2` en la hoja ✔; asa de 32 × 4 en `ui-border-control` ✔; barra lateral colapsable ✔; objetivos de 44 px o más en móvil ✔ (96 controles medidos, 0 bajo 44); separación de 8 px entre objetivos: N-3; foco visible con `ui-border-focus` ✔.
- **U5-07:** 0 palabras en inglés («total» es español); 0 emoji; 0 fechas ISO visibles; dinero ✔.

## Texto con aspecto de instrucción

Ninguno. Busqué «ignora», «instrucción», «Claude», «agente», «orquestador», «prompt», «aprueba» y comentarios HTML en los 6 artboards y el fragmento: 0 coincidencias. Las notas «Así se ve la hoja…» y «Así se ve el diálogo…» explican el dibujo; no dan órdenes.

## Notas para el orquestador (no son hallazgos del autor)

1. `./support.js` y `../ds/finanzas-hogar/tokens.css` no existen todavía en `lienzo/project/`; supongo que el orquestador los pone al publicar.
2. Sigue abierta la nota de la vuelta 1 sobre `is_interactive: false`: los enlaces solo navegan en Play.
3. Los altos de p10-m (2700) y p10-e (2300) son estimaciones (B-8); conviene medirlos al publicar.
