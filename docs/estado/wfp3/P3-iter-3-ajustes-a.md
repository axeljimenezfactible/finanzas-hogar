# P3-iter-3 · ajustes-a · revisión independiente (vuelta 3 de 3)

- Fecha: 12-oct-2026. Revisor: `revisor-ajustes-a-r3`. Solo revisa, no corrige.
- Alcance: los 16 artboards de la parte a en `docs/design/maquetas/lienzo/project/ajustes/` (P-26 índice y secciones 01 a 07, móvil `-m` y escritorio `-e`) y `docs/design/maquetas/lienzo/fragmentos/ajustes-a.json`, tal como los dejó `correccion-ajustes-a-r2`.
- Se revisa contra:
  - `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08 y P3-12);
  - `docs/criterios/U5.md`: U5-04, U5-06 y U5-07, en lo que se ve en el marcado;
  - DESIGN.md: §2.5, §7.12, §7.13, §7.14, §9.3 y §10.5 (con `docs/design/contraste-componentes.csv`);
  - U1 §2.1, §2.2 y §2.3;
  - U4 §0 y P-26.01 a P-26.07;
  - `referencia-tipo/format.md`.
- **Veredicto: APROBADO. 0 altos, 0 medios, 8 bajos.** Se cierran los 2 medios de la vuelta 2. De los bajos, 6 vienen de antes y siguen abiertos, y 2 son nuevos. Ninguno impide publicar.

## Hallazgos de la vuelta 2: estado

| Id r2 | Estado | Evidencia (script propio, cascada CSS del `<helmet>` resuelta por elemento) |
|---|---|---|
| M-01 · diálogos críticos de s02 y s05 | **Cerrado** | En s02-m, s02-e, s05-m y s05-e, cada `.dlg.cr` queda así: <br>• **Superficie**: fondo `ui-surface-raised`, radio `radius-lg` y sombra `shadow-2` (en s05 viene de la regla `.dlg` nueva). <br>• **Encabezado** `h3`: color `ui-critical`, con `<use href="#ic-shield-alert">`. El `<symbol>` (3 trazos, `viewBox 0 0 24 24`) está dentro del `<svg>` del sprite. <br>• **Botón** `.btn.c`: fondo `ui-action-critical-bg` y texto `ui-text-on-critical`, más `:hover` con `-bg-hover`. Ninguna regla posterior lo pisa. <br>• **Cancelar**: botón `.btn.t` con `autofocus`. <br>Pares de contraste: C-024 5.60/10.55, C-025 8.23/12.67, C-027 5.37/8.06 y C-139 5.37/6.18 (claro/oscuro). |
| M-02 · borde de error en 5 `-m` | **Cerrado** | `.inw.er{border-color:var(--ui-border-error)}` está en s01-m, s03-m, s04-m, s05-m y s07-m. En las 14 páginas de sección, el campo con `aria-invalid="true"` tiene: <br>• borde `ui-border-error` (C-039: 5.78/7.61); <br>• `aria-describedby` hacia un mensaje con `circle-alert` en `ui-role-error`; <br>• el texto literal del catálogo §0.6 (V-15, V-03, V-04, V-12, V-23, V-25 y V-09). |
| B-01 · `.btn.t`, `.row`, `.tx`, `.cn` sin regla | Sigue | s04 m/e: `btn t` ×13, `row`, `tx` y `cn` ×1. s07-m: `btn t` ×1. |
| B-02 · concordancia | Sigue | s05-e: «No se pudo cargar los créditos». |
| B-03 · leyenda «Enter…» fuera de campo | Sigue | s02-e y s07-e. |
| B-04 · atajos incompletos | Sigue | Los 8 `-e` omiten `Ctrl+Shift+L` y `↑`/`↓`. |
| B-05 · alturas estimadas | Sigue | Sin navegador (K45). |
| B-06 · primaria no única | Sigue | `btn p`: 6 en s02-e, 5 en s05-e y 3 en s06-e. |

El corrector solo tocó el `<helmet>` y el sprite de 9 archivos. Lo comprobé con un diff contra `docs/design/salida/iter-2/ajustes-a/`: no cambió marcado, textos, alturas ni el fragmento.

## Verificación hecha (scripts propios, `python3 -I`, en `datos/tmp/rev-aj-a-r3/`, borrados al final)

| Comprobación | Resultado |
|---|---|
| (a) Forma: doctype primero, `lang="es-MX"`, `<meta charset="utf-8">`, `<title>` propio, `./support.js`, `../ds/finanzas-hogar/tokens.css`, `<x-dc>` → `<helmet><style>` en orden, un solo `<style>`, `body{}`, `a{}` y `a:hover{}` en el helmet, `data-dc-script` con comillas simples, `class Component extends DCLogic` y `renderVals()` | 16/16 |
| `data-props`: JSON válido, solo `$preview`, igual a `w`×`h` del fragmento | 16/16 |
| CSS: llaves y paréntesis balanceados; declaraciones con `:`; sin huecos `{{}}` | 16/16 |
| `html.parser` (árbol propio con elementos vacíos y SVG autocerrado): etiquetas sin cerrar, ids repetidos, `for`/`aria-labelledby`/`aria-describedby`/`aria-controls` sin destino, controles dentro de `<a>` | 0 / 0 / 0 / 0 |
| Nombres `[a-z0-9.-]`; stems únicos en los 9 directorios del lienzo | Bien |
| Fragmento | 16 `boards` = 16 `order` = 16 archivos en disco; `-e` con `w` 1440 y `expand: "fill"`; `-m` con `w` 412 y `h` ≥ 915; `page: "ajustes"`; `is_interactive: false`; títulos «P-26.0X · Nombre · móvil\|escritorio» |
| Sin `iframe`, `img`, `canvas`, `innerHTML`, `http(s)://` ni `fetch` | 0 |
| (b) Hex literales fuera de `data-props` | 0 |
| `var(--x)` inexistentes en `tokens.json` | 0 |
| Clases del marcado que no aparecen en ninguna regla que case con ese elemento | s04 m/e (`t`, `row`, `tx`, `cn`), s06 m/e (`ro`, nuevo) y s07-m (`t`) |
| `<use href>` sin `<symbol>` | 0 |
| (c) Contraste con cascada CSS propia (selectores del `<helmet>`, especificidad y orden, herencia de `color` y fondo efectivo del ancestro) y tokens resueltos en claro y oscuro | 27 pares de texto distintos y 38 de interfaz (bordes de control, iconos, pistas de interruptor). Además, 14 combinaciones color/fondo de reglas de estado (`:hover`, `[aria-current]`, `[disabled]`). **0 bajo umbral** |
| (d) P3-12 | Se buscan las 21 cifras canon (con y sin coma), sus centavos, credenciales, correos, `$` de 4 cifras o más sin separador y fechas ISO: **0** en todo |
| Montos | 66 distintos (texto y `value=`). 60 están en el JSON. Los otros 6 son derivados o ejemplos declarados: $26,686.00 = 27,586 − 900; $27,786.00 = 27,586 − 3,200 + 3,400; $3,400.00 es el valor tecleado del ejemplo; $29,040.00 = 14,260 + 14,780; $29,210.00 = 14,260 + 14,950; $14,950.00 es el valor tecleado |
| Fechas | 25 distintas: 24 están en el JSON y 1 se calcula (14-ene-2027) |
| Fechas del oráculo (comando de P3-12) | En pantalla solo aparece 30-oct-2026, junto a $1,861.58, $6,189.05 y $6,881.60, que son sintéticos. Sin hallazgo |
| (e) Fidelidad | Ver la tabla de fidelidad; sin faltas |
| (f) Accesibilidad e idioma | 0 emoji; 0 palabras de interfaz en inglés («Error», «No», «Total» y «Plan» son español; «Open» es de Open Food Facts); 0 trato de usted; dinero 66/66 `$1,234.56`; negativos con U+2212 (−$17.70, −$212.40); `tabular-nums` en `body`; 0 `svg` sin `aria-hidden`; 0 botones o enlaces sin nombre; 0 campos sin `label` o `aria-label`; sin barra de estado ni teclado falsos |
| Objetivos táctiles en `-m` (alto calculado con reglas y tokens) | Todo `a`, `button`, `input` y `select` mide ≥ 44 px (`size-touch-min` o `size-control-coarse`). Las 14 tarjetas de indice-m no tienen `min-height`, pero llevan 12 + 12 px de relleno y dos líneas: ≥ 64 px |
| (g) Raíz | `-m`: un solo `.m` con `width: var(--breakpoint-phone)` (412 px), `min-height: 915px` y `*{box-sizing:border-box}`. `-e`: `.pg` sin ancho en px, con `flex-wrap: wrap`; `.mw` con `flex: 999 1 560px` y `.ct` con `flex: 999 1 480px`; tablas en `.tw{overflow-x:auto}` |
| (h) Enlaces | 0 rotos y 0 a otro ancho. Cada `href` existe en el lienzo. Hay 8 `href="#"` (sin artboard destino en la lista) y 2 `#instrucciones`, ambos con destino |

## Recálculo de 10 pares de contraste (P3-07; comparados con `contraste-componentes.csv`, ±0.01 por truncado)

| Par | Claro | Oscuro | Id §10.5 |
|---|---|---|---|
| `ui-text-on-critical` / `ui-action-critical-bg` | 5.60 | 10.55 | C-024 (5.59 / 10.54) |
| `ui-text-on-critical` / `ui-action-critical-bg-hover` | 8.23 | 12.67 | C-025 (8.22 / 12.67) |
| `ui-action-critical-bg` / `ui-surface-raised` (UI) | 5.37 | 8.06 | C-027 (5.36 / 8.06) |
| `ui-critical` / `ui-surface-raised` (icono y título) | 5.37 | 6.18 | C-139 (5.36 / 6.18) |
| `ui-border-error` / `ui-surface-page` (UI) | 5.78 | 7.61 | C-039 (5.78 / 7.60) |
| `ui-text-secondary` / `ui-surface-page` | 5.40 | 10.43 | C-140 (5.39 / 10.42) |
| `ui-text-on-primary` / `ui-action-primary-bg` | 5.77 | 10.47 | C-001 (5.77 / —) |
| `ui-role-error` / `ui-surface-page` (mensaje) | 5.78 | 7.61 | ≥ 4.5 |
| `ui-text-secondary` / `ui-role-error-surface` | 5.01 | 8.30 | ≥ 4.5 |
| `ui-border-control` / `ui-surface-card` (UI) | 3.47 | 6.09 | ≥ 3 |

## Hallazgos de esta vuelta

### Altos y medios: ninguno

### Bajos (8)

- **B-01** (sigue): sin `.btn.t` en s04 `-m`/`-e` (los 13 asas «Mover {categoría}») ni en s07-m («Ir a Datos y respaldos»). Salen con borde de control en lugar de terciarios; el alto de 44 px sí se cumple. En s04, el bloque «Así se ve un valor no válido» usa `.row`, `.tx` y `.cn` sin regla. **Corrección**: copiar las reglas de s03.
- **B-02** (sigue): s05-e dice «No se pudo cargar los créditos». **Corrección**: «No se pudieron cargar los créditos».
- **B-03** (sigue): la leyenda «Enter o salir del campo guarda · Esc descarta» aparece bajo vistas previas que no son de un campo (s02-e y s07-e).
- **B-04** (sigue): la línea de atajos de los 8 `-e` omite `Ctrl+Shift+L` (claro y oscuro) y `↑`/`↓` (U4 §0.7).
- **B-05** (sigue): las alturas del fragmento y de `$preview` son estimadas. Hay que medirlas al publicar.
- **B-06** (sigue): la acción primaria no es única: s02-e tiene 6 `btn p` (2 en tarjetas de «Estados»), s05-e 5 y s06-e 3.
- **B-07** (nuevo): en s06 `-m` y `-e`, la clase `ro` (valor de solo lectura: «Tianguis de la Narvarte», «Mercado Portales» y el veredicto del punto de equilibrio de Sam's) no tiene regla en el `<helmet>`. Se dibuja como texto suelto, sin el alto de control ni la alineación que tiene en s05 (`.ro{min-height:var(--size-control-coarse);display:flex;align-items:center}`). Las revisiones anteriores no lo vieron porque solo leían `class="…"` con comillas dobles, y aquí va con comillas simples. **Corrección**: copiar `.ro` de s05.
- **B-08** (nuevo): la anatomía de los diálogos de ejemplo no sigue del todo DESIGN §7.12 ni U4 §0.4.
  - En los `-m` (s01, s02 y s05), «En móvil el diálogo se vuelve hoja inferior». Pero se dibujan como caja con `radius-lg` en las 4 esquinas, sin asa de 32 × 4 y sin primaria de 48 px abajo.
  - En los `-e`, la caja no tiene el ancho `layout-dialog-width-short` (480 px) y ocupa todo el contenido.
  - El cuerpo va en `ui-text-default`, no en `ui-text-secondary`.
  - Falta `aria-modal="true"`.
  - s02 `-m` y `-e` llevan 2 `autofocus` en el mismo documento.

  Lo crítico (icono, color, botón, Cancelar con el foco) sí está bien. **Corrección**: en `-m`, hoja inferior; en `-e`, `max-width: var(--layout-dialog-width-short)`; `.dlg p{color:var(--ui-text-secondary)}`; un solo `autofocus` por artboard.

## Líneas de evidencia

- **P3-07**: hay revisión por artboard (tabla de abajo). El contraste se midió en el HTML con cascada propia y tokens resueltos en los dos temas: 0 pares de texto bajo 4.5:1 y 0 de interfaz bajo 3:1. Recalculé 10 pares y coinciden con `contraste-componentes.csv` (±0.01). No hay altos ni medios.
- **P3-08** (8 puntos):
  1. Resumen ≤ 6 tarjetas: no aplica a Ajustes.
  2. 14 secciones de Ajustes ✓: columna interna de los 8 `-e` en el orden de U1 §2.3, con `aria-current="page"` en la propia; lista de 14 en indice-m; «‹ Ajustes» en los 7 `-m`.
  3. Barra lateral colapsable ✓: «Contraer barra lateral (Ctrl+\)», menú del avatar y grupo «Tema».
  4. Español de México ✓.
  5. `$1,234.56` tabular ✓: el grep `\$[0-9]{4,}\.[0-9]{2}` da 0.
  6. Fechas `d-mmm-aaaa` sin cero ✓.
  7. Paleta C2 solo por tokens ✓: 0 hex.
  8. Notion como única referencia ✓.
- **P3-12**: 0 cifras canon (con y sin coma, y en centavos), 0 credenciales y 0 correos. Montos y fechas salen de `datos-sinteticos.json` o se derivan de él. Muestra de 10:
  - $212.40: `canasta.membresia_sams.ahorro_anual_cents` (−21240);
  - $17.70: `…ahorro_por_canasta_cents` (−1770);
  - $12,000.00: `creditos.lineas[0].limite_cents`;
  - $399.00: `presupuesto.gastos_fijos_recurrentes[1].monto_cents`;
  - $250.00: `creditos.lineas[0].comision_tardio_cents`;
  - $766.50: `categorias[7].real_cents`;
  - $14,260.00: `ingresos.esperado.q15_cents`;
  - 48,210 registros: `ajustes.profeco.registros`;
  - 1-ene-2027: `ajustes.sesiones[2].vence`;
  - 2-nov-2026: `creditos.promociones[0].inicio`.

  «Hoy es 21-oct-2026» coincide con `meta.hoy`. **Cumple.**
- **U5-04** (estático): 0 pares bajo umbral. La medición en el DOM la hace U5.
- **U5-06**:
  - Tokens de rol ✓ y error en rojo ✓.
  - **Acción crítica con el matiz opuesto ✓** en s01, s02 y s05: era el M-01 y queda cerrado.
  - Tipografía de `tokens.css` ✓ y cifras tabulares ✓.
  - Lucide con trazo 1.5 y 20 px ✓.
  - Espaciado y radios por token ✓: los únicos px literales son bases de `flex`, `min-height` 900/915 y la rejilla `minmax(min(260px,100%),1fr)`.
  - Foco visible ✓, barra lateral colapsable ✓ y 44 px ✓.
- **U5-07**: sin inglés ✓, dinero ✓, fechas ✓ y sin ISO ✓. Queda la concordancia de B-02.

## Por artboard: reglas cumplidas / falladas / correcciones

Los 16 cumplen lo común: forma de format.md, P3-12, U5-04 estático, tokens, `lang` y título, nombres y etiquetas accesibles, ancho (412 fijo o página fluida) y enlaces. Además, M-01 y M-02 están cerrados donde aplicaban.

| Artboard | Falladas (ID, todas bajas) | Corrección opcional |
|---|---|---|
| indice-m | — | — |
| indice-e | B-04 | Atajos completos |
| s01-m | B-08 | Diálogo como hoja inferior |
| s01-e | B-04, B-08 | Ancho 480 px y cuerpo secundario |
| s02-m | B-08 | Hoja inferior; un solo `autofocus` |
| s02-e | B-03, B-04, B-06, B-08 | Un solo primario; leyenda; ancho del diálogo |
| s03-m | — | — |
| s03-e | B-04 | — |
| s04-m | B-01 | `.btn.t`, `.row`, `.tx` y `.cn` |
| s04-e | B-01, B-04 | Ídem |
| s05-m | B-08 | Hoja inferior |
| s05-e | B-02, B-04, B-06, B-08 | «No se pudieron cargar» |
| s06-m | B-07 | `.ro` |
| s06-e | B-04, B-06, B-07 | `.ro` |
| s07-m | B-01 | `.btn.t` |
| s07-e | B-03, B-04 | — |
| ajustes-a.json | B-05 | Medir `h` al publicar |

## Fidelidad por pantalla (U1 §2.1, §2.2 y §2.3, y U4)

| Artboard | Presente | Falta o difiere |
|---|---|---|
| índice m/e | 14 secciones en orden, cada una con su resumen tomado del JSON. Columna interna (e) y lista (m). «Estados» (e): vacío «No aplica», carga con esqueleto, error al cargar y error al guardar | — |
| s01 | Usuarios (co-titulares, último acceso); Karen sin cuenta; contraseña con V-15; sesiones; «Cerrar sesiones» con diálogo crítico; «Este dispositivo» con `#instrucciones`; 4 parámetros K44; «Estados» (e) | B-08 (anatomía) |
| s02 | Interruptor de proyección con vista previa; montos por categoría con «Quitar monto» y diálogo crítico; diálogo de alcance; gastos fijos con «Nuevo gasto fijo»; despensa; V-03; vacíos; «Estados» (e) | B-03, B-06, B-08 |
| s03 | Quincenas del 15 y del 30, «Rige desde», Ximena sin ingresos al modelo, enlace a P-08, vista previa Antes → Después, V-04, «Estados» (e) | — |
| s04 | Tabla y Tablero, «Activa» como `switch`, «Mover {categoría}», «Nueva categoría», vista previa del reparto, V-12, «Estados» (e) | B-01 |
| s05 | Líneas con «Falta capturar»; A/B con vista previa; tasas; calendario de promociones (e) o agenda (m) con «Vigente hoy»; «Eliminar promoción» con diálogo crítico; comparador con pesos y V-23; vacíos de créditos, tasas y promociones; «Estados» (e) | B-02, B-06, B-08 |
| s06 | Tiendas en Tabla, Galería y Tablero; visitas a mercados; membresía de Sam's con punto de equilibrio; fuentes con V-19; PROFECO; parámetros K44 con V-25; «Nueva tienda o mercado»; «Estados» (e) | B-06, B-07 |
| s07 | Recetario, porciones, bases nutricionales con «Recargar base nutricional», restricciones, IA, V-09, «Nueva receta», «Estados» (e) | B-01 (m), B-03 (e) |

En las tablas de campos de U4, los campos que no tienen rótulo propio están como valores. Por ejemplo: comercio, tasa, vigencia y fuente de cada promoción en la franja del calendario; «Cada mes · Día 1» como periodicidad del gasto fijo. Son vistas de base de datos con propiedades visibles, no formularios completos. No cuentan como faltas.
