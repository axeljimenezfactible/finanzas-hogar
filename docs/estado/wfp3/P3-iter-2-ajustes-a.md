# P3-iter-2 · ajustes-a · revisión independiente (vuelta 2 de 3)

- Fecha: 12-oct-2026. Revisor: `revisor-ajustes-a-r2` (no corrige; solo revisa).
- Alcance: los 16 artboards de la parte a en `docs/design/maquetas/lienzo/project/ajustes/` (P-26 índice y secciones 01 a 07, móvil `-m` y escritorio `-e`) y `docs/design/maquetas/lienzo/fragmentos/ajustes-a.json`, tal como los dejó `correccion-ajustes-a-r1`.
- Contra: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06, U5-07, lo que se ve en el marcado), DESIGN.md (§2.5, §7.4, §7.9, §7.12, §9.3, §10.5), U1 §2.1, §2.2 y §2.3, U4 §0 y P-26.01 a P-26.07, y `referencia-tipo/format.md`.
- **Veredicto: NO aprobado.** 0 altos, 2 medios, 6 bajos. De los 8 medios de la vuelta 1 se cierran 6 por completo. Dos quedan a medias, y los dos por la misma causa: en los archivos nuevos o tocados falta la regla CSS o el símbolo que dibuja lo añadido. Las clases del marcado no tienen su regla en el `<helmet>` de ese archivo.

## Verificación hecha (scripts propios, `python3 -I`, en `datos/tmp/rev-aj-a-r2/`, borrados al final)

| Comprobación | Resultado |
|---|---|
| (a) Forma: doctype, `lang="es-MX"`, `<meta charset>`, `./support.js`, `../ds/finanzas-hogar/tokens.css`, `<x-dc>` y `<helmet>` en orden, `data-dc-script` con `class Component extends DCLogic` y `renderVals()` | 16/16 bien |
| `data-props` es JSON válido y `$preview` = `w`×`h` del fragmento | 16/16 bien |
| `html.parser`: etiquetas sin cerrar, ids repetidos, `for`/`aria-*` sin destino, controles dentro de `<a>` | 0 / 0 / 0 / 0 |
| Nombres `[a-z0-9.-]`; stems únicos en todo el lienzo (los 10 directorios); fragmento = 16 `boards` = 16 `order` = archivos en disco de la parte a | bien |
| Sin `iframe`, `img`, `innerHTML`, `http(s)://`, `fetch` | 0 |
| (b) Hex literales fuera de `data-props` | 0 |
| `var(--x)`: todos existen en `tokens.json` | bien (0 inexistentes) |
| **Clases del marcado sin regla CSS en su `<helmet>`** (sin contar las de `tokens.css`: `body`, `caption`, `small`, `title-1..3`) | **s02 m/e: `c`, `cr`; s05 m/e: `c`, `cr`, `dlg`; s01-m, s03-m, s04-m, s05-m y s07-m: `er`; s04 m/e: `t`, `row`, `tx`, `cn`; s07-m: `t`** (M-01, M-02, B-01) |
| **`<use href="#…">` sin `<symbol>` destino** | **`#ic-shield-alert` en s02-m, s02-e, s05-m y s05-e** (M-01) |
| (c) Contraste con los tokens de `tokens.json` resueltos en claro y oscuro: matriz de los 13 tokens de texto y los 7 de UI usados contra las 9 superficies usadas, 450 combinaciones | 10 bajo umbral en la matriz, pero ninguno aparece en el marcado: los colores de rol como texto solo van sobre página, tarjeta o su superficie de rol, nunca sobre `ui-surface-hover`. El anillo de avatar `ui-person-axel` sobre `ui-surface-hover` da 2.14 en oscuro y es la anatomía de DESIGN §7.4, con su par C-059 sobre la página aprobado. **0 pares reales bajo umbral** |
| (d) P3-12: grep de las 21 cifras canon con y sin coma y en centavos; grep de credenciales | 0 / 0 |
| Dinero y fechas contra `datos-sinteticos.json` (los campos `*_cents` /100) | 60 montos distintos: 55 están en el JSON y 5 son sumas o ejemplos declarados ($26,686.00 = 27,586.00 − 900.00; $27,786.00; $29,040.00; $14,950.00; $29,210.00). 25 fechas: 24 están en el JSON y 1 se calcula (14-ene-2027) |
| Muestra de 10 valores | $212.40 · $12,000.00 · $465.00 · $27,586.00 · $766.50 · $17.70 · $1,500.00 · 2-oct-2026 · 1-ene-2027 · 2-nov-2026 |
| Fechas que coinciden con el oráculo (comando de P3-12) | En pantalla solo 30-oct-2026, con $1,861.58 sintético: sin hallazgo |
| (e) Fidelidad: filas de las tablas de campos de U4 P-26.01 a P-26.07 buscadas en el texto de cada artboard | Sin faltas reales (las 4 que marca el script son encabezados o valores sin rótulo: «Comercio», «Vigencia», «Dirección», «Alcance» como diálogo) |
| Columna interna: las 14 secciones en el orden de U1 §2.3, con `aria-current` en la propia | 8/8 `-e` bien; lista del índice móvil con 14 bien; «‹ Ajustes» en los 7 `-m` hacia `ajustes-p26-indice-m` |
| (f) Emoji; guion como signo de negativo; fechas ISO visibles; fechas con cero a la izquierda | 0 / 0 (los negativos usan U+2212: −$17.70, −$212.40) / 0 / 0 |
| `svg` sin `aria-hidden`; botones o enlaces sin nombre; campos sin `label` o `aria-label`; barra de estado o teclado falsos | 0 / 0 / 0 / 0 |
| Iconos: sprite de `<symbol viewBox="0 0 24 24">`, `stroke-width` 1.5, `size-icon-sm` | bien |
| Objetivos táctiles en `-m`: alto calculado de cada `a`, `button`, `input` y `select` con las reglas del `<helmet>` y los tokens | Todos ≥ 44 px. Los `select` miden 40 px dentro de un `.inw` de 44 px, que es su área. La casilla «Ya lo confirmé» va en una fila de 44 px con `label for` de toda la fila |
| (g) Raíz móvil `.m` con `width: var(--breakpoint-phone)` (412 px), `min-height: 915px`, `box-sizing: border-box`; escritorio `.pg` con `flex-wrap`, sin ancho en px, `.mw` `flex: 999 1 560px`, tablas y calendario en `.tw{overflow-x:auto}`; `expand: "fill"` solo en `-e` | bien |
| (h) Enlaces: cada `href` está en la lista del lienzo (ajustes a y b incluidas) y existe en disco, o es `#` o un ancla con destino; el mismo ancho que el origen | 0 rotos, 0 a otro ancho |

## Hallazgos de la vuelta 1: estado

| Id r1 | Estado | Evidencia |
|---|---|---|
| M-01 enlaces a `s08-reportes-frecuencias` | **Cerrado** | 0 `href` fuera de la lista; los 10 enlaces apuntan a `ajustes-p26-s08-reportes-{e,m}.dc.html`, que existe |
| M-02 `#instrucciones` | **Cerrado** | El ancla tiene destino (`id="instrucciones"`) en s01 `-e` y `-m` |
| M-03 objetivos táctiles | **Cerrado** | «Ver instrucciones» y los «Ir a …» llevan `btn` (44 px); la casilla, una fila de 44 px con `label` |
| M-04 estados | **Cerrado** | Tarjeta «Error al guardar» con el texto literal de U4 §0.3 en los 7 `-e` de sección. En el índice y en el vacío de s01, «No aplica». Vacíos de U4 presentes: «Sin gastos fijos · Nuevo gasto fijo» (s02), «Sin tasas · Nueva tasa» y «Sin promociones · Nueva promoción» (s05), «Sin cortes importados · Importar corte PROFECO» (s06) |
| M-05 validaciones | **A medias → M-02 de esta vuelta** | El campo inválido con `aria-invalid`, `aria-describedby`, icono y mensaje literal del catálogo está en las 14 páginas de sección (V-15, V-03, V-04, V-12, V-23, V-25 y V-09). En 5 `-m` falta la regla del borde de error |
| M-06 fechas con cero | **Cerrado** | 0 fechas `0d-mmm-aaaa` |
| M-07 barra lateral | **Cerrado** | En los 8 `-e` están el botón del avatar (`aria-haspopup="menu"`: Cuenta, Apariencia, Cerrar sesión), el botón `panel-left` «Contraer barra lateral» (Ctrl+\), el grupo «Tema» (Sistema, Claro, Oscuro) y `Ctrl+\` en la línea de atajos |
| M-08 elementos de U4 | **A medias → M-01 de esta vuelta** | (1) Calendario de s05-e por defecto, con franjas y «Vigente hoy»: cerrado; móvil en agenda por mes. (3) «Quitar el monto de {categoría}» en las 13 filas de s02 m/e: cerrado. (4) «Activa» como `role="switch"` en las 13 filas de s04 m/e: cerrado. (5) 7 bloques plegables `aria-expanded` en s05-m: cerrado. (2) y (3): los diálogos críticos nuevos de s02 y s05 no se dibujan como críticos |
| B-01 título | Cerrado | `<title>Ajustes</title>` |
| B-02 concordancia | Sigue en s05-e | «No se pudo cargar los créditos» |
| B-03 leyenda «Enter…» | Sigue (el corrector la dejó pendiente) | — |
| B-04 asa de reordenar | Cerrado en función | Botón «Mover {categoría}» de 44 px; el estilo terciario no se aplica (B-01 de esta vuelta) |
| B-05 alturas | Sigue | Alturas estimadas, sin navegador (K45) |
| B-06 primaria única | Sigue en parte | s01 bien; s02-e tiene 6 `btn p` (2 en ejemplos de estado), s05-e 5, s06-e 3 |
| B-07 relleno crítico fuera del diálogo | Cerrado | — |

## Hallazgos de esta vuelta

### Medios (2)

- **M-01 · Los diálogos de acción crítica de s02 y s05 no se dibujan como críticos** (U4 §0.4, DESIGN §2.5 y §7.12, U5-06 «acción crítica con el matiz opuesto», P3-07). Archivos: `ajustes-p26-s02-presupuesto-{m,e}` y `ajustes-p26-s05-creditos-condiciones-{m,e}`. El marcado usa `class="dlg cr"`, `class="btn c"` y `<use href="#ic-shield-alert"/>`. En esos 4 archivos faltan tres cosas: las reglas `.btn.c` y `.dlg.cr h3`, el `<symbol id="ic-shield-alert">` y, en s05, también la regla `.dlg`. Resultado: «Quitar monto» y «Eliminar promoción» salen como botones secundarios con borde, sin el relleno `ui-action-critical-bg`. El encabezado no lleva el icono `shield-alert` (el `<use>` queda vacío) ni el color `ui-critical`. En s05, el diálogo no tiene superficie, borde, radio ni sombra: se lee como texto suelto. Corrección: copiar de `ajustes-p26-s01-cuenta-roles-e.dc.html` las reglas `.btn.c`, `.btn.c:hover`, `.dlg`, `.dlg h3` y `.dlg.cr h3` y el `<symbol id="ic-shield-alert">` a los 4 archivos.
- **M-02 · Falta el borde de error en 5 artboards móviles** (U4 §0.3 «borde de error, icono `circle-alert` y el mensaje»; criterio (e), validaciones). Archivos: `s01-cuenta-roles-m`, `s03-ingresos-m`, `s04-categorias-m`, `s05-creditos-condiciones-m` y `s07-recetas-m`. El campo inválido lleva `<div class="inw er">`, pero en su `<helmet>` no hay ninguna regla `.inw.er`: `ui-border-error` aparece 0 veces en esos 5 archivos. El campo se ve con el borde normal `ui-border-control`; solo el mensaje marca el error. Corrección: añadir `.inw.er{border-color:var(--ui-border-error)}` al `<helmet>` de los 5, como en s06-m y en los `-e`.

### Bajos (6)

- **B-01**: `.btn.t` no está definido en s04 `-m` y `-e` ni en s07-m. Los asas «Mover {categoría}» y «Ir a Datos y respaldos» salen con borde de control y no como terciarios. El alto de 44 px sí se cumple, porque viene de `.btn`. En s04 `-m` y `-e`, el bloque «Así se ve un valor no válido» usa `.row`, `.tx` y `.cn` sin regla, así que se apila sin la maqueta de fila. Hay que copiar las reglas de s03.
- **B-02**: queda en s05-e «No se pudo cargar los créditos». Debe decir «No se pudieron cargar los créditos».
- **B-03**: sigue la leyenda «Enter o salir del campo guarda · Esc descarta» bajo vistas previas que no son de un campo (s02-e y s07-e). El corrector la declaró pendiente.
- **B-04**: la línea de atajos de los `-e` omite `Ctrl+Shift+L` (claro y oscuro) y `↑`/`↓` de U4 §0.7.
- **B-05**: las alturas del fragmento y de `$preview` siguen siendo estimadas (s05-e 8,510 y s02-m 7,940). Hay que medirlas al publicar.
- **B-06**: la acción primaria no es única. s02-e tiene 6 `btn p` (2 en tarjetas de ejemplo de «Estados»), s05-e 5 y s06-e 3.

## Líneas de evidencia

- **P3-07**: hay una revisión por artboard (tabla de abajo). El contraste se calculó con los tokens resueltos en los dos temas: 0 pares reales bajo 4.5:1 en texto o bajo 3:1 en UI. Hallazgos: M-01, M-02 y los bajos B-01 a B-06.
- **P3-08**: 14 secciones de Ajustes ✓ (índice, columna interna y lista móvil). Barra lateral colapsable con `panel-left` y Ctrl+\ ✓. Español de México ✓. Dinero `$1,234.56` con `tabular-nums` en `body` y U+2212 ✓. Fechas `d-mmm-aaaa` sin cero ✓. Paleta solo por tokens de C2 ✓. Notion como referencia ✓. La acción crítica con el matiz opuesto no se dibuja en s02 ni en s05 (M-01). Las claves `groceries` y `eating_out` siguen siendo un choque entre U4 y U5-07: va para el orquestador y no es un error del autor.
- **P3-12**: 0 cifras canon (con y sin coma, y en centavos) y 0 credenciales. Montos y fechas salen de `datos-sinteticos.json` o se calculan a partir de él (la muestra de 10 está arriba). La revisión manual de fechas encuentra solo 30-oct-2026, con un monto sintético. Cumple.
- **U5-04** (estático): 450 combinaciones de la matriz y 0 pares reales bajo umbral en los dos temas. La medición en el DOM la hace U5.
- **U5-06**: tokens de rol ✓; tipografía de `tokens.css` ✓; cifras tabulares ✓; Lucide con trazo 1.5 y 20 px ✓; espaciado y radios por token ✓; foco visible ✓; barra lateral colapsable ✓; objetivos de 44 px ✓. **Acción crítica con el matiz opuesto ✗** en s02 y s05 (M-01).
- **U5-07**: sin inglés, salvo las claves que pide U4 ✓; dinero ✓; fechas ✓; sin ISO visible ✓. Concordancia en s05-e (B-02).

## Por artboard: reglas cumplidas / falladas / correcciones

Común a los 16: cumplen la forma de format.md, P3-12, U5-04 estático, el uso de tokens, `lang` y el título, la accesibilidad de nombres y etiquetas, el ancho (412 fijo o página fluida) y los enlaces.

| Artboard | Falladas (ID) | Corrección para la vuelta 3 |
|---|---|---|
| indice-m | — | — |
| indice-e | B-04 | Atajos completos |
| s01-m | M-02 | `.inw.er` en el `<helmet>` |
| s01-e | B-04 | — |
| s02-m | M-01 | Reglas `.btn.c`, `.dlg.cr h3` y `<symbol id="ic-shield-alert">` |
| s02-e | M-01, B-03, B-04, B-06 | Igual que s02-m; un solo primario |
| s03-m | M-02 | `.inw.er` |
| s03-e | B-04 | — |
| s04-m | M-02, B-01 | `.inw.er`, `.btn.t`, `.row`, `.tx` y `.cn` |
| s04-e | B-01, B-04 | `.btn.t`, `.row`, `.tx` y `.cn` |
| s05-m | M-01, M-02 | `.dlg`, `.dlg h3`, `.dlg.cr h3`, `.btn.c`, `.inw.er` y el símbolo |
| s05-e | M-01, B-02, B-04, B-06 | `.dlg`, `.dlg h3`, `.dlg.cr h3`, `.btn.c` y el símbolo; «No se pudieron cargar» |
| s06-m | — | — |
| s06-e | B-04, B-06 | — |
| s07-m | M-02, B-01 | `.inw.er` y `.btn.t` |
| s07-e | B-03, B-04 | — |
| ajustes-a.json | B-05 | Medir `h` al publicar |

## Fidelidad por pantalla (U1 §2.1, §2.2 y §2.3, y U4)

| Artboard | Presente | Falta o difiere |
|---|---|---|
| índice m/e | 14 secciones en orden, con una tarjeta de resumen por sección tomada del JSON; columna interna (e) y lista (m); «Estados» (e) con vacío «No aplica», carga con esqueleto, error al cargar y error al guardar «No aplica» | — |
| s01 | Usuarios, contraseña con V-15 en el campo real, sesiones, cierre con diálogo crítico completo, «Este dispositivo» con `#instrucciones`, 4 parámetros K44 | Borde de error en `-m` (M-02) |
| s02 | Interruptor de proyección, montos por categoría con «Quitar el monto de {categoría}», diálogo de alcance, gastos fijos, despensa, V-03, vacíos de montos y de gastos fijos, los 4 estados | El diálogo crítico «Quitar el monto de Ocio» no se ve como crítico (M-01) |
| s03 | Persona, frecuencia, montos lado a lado, «Rige desde», Ximena, enlace a P-08, vista previa, V-04, estados | Borde de error en `-m` (M-02) |
| s04 | Tabla con interruptor «Activa», asa «Mover», Tabla y Tablero, vista previa del reparto, V-12, estados | Borde de error en `-m` (M-02); estilos de `.btn.t` y de la fila (B-01) |
| s05 | Líneas, condiciones, A/B, tasas, calendario por defecto (e) o agenda (m), «Eliminar promoción», bloques plegables (m), comparador, V-23, vacíos de créditos, tasas y promociones | Diálogo crítico sin estilo ni icono (M-01); borde de error en `-m` (M-02); B-02 |
| s06 | Tiendas, mercados, Sam's, fuentes, V-19, parámetros con V-25, vacíos de tiendas y de PROFECO, estados | — |
| s07 | Recetario, porciones, bases, restricciones, IA, V-09, estados | Borde de error en `-m` (M-02); `.btn.t` (B-01) |
