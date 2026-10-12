# correccion-U2-r1 · WF-1 · 11-oct-2026

Corrige los 5 hallazgos altos y medios de `docs/estado/wf1/REVISION-U2-r1.md` (H-1/C2-05, H-2/C2-04, M-1/U2-11, M-2/C2-10, M-3/C2-09·U2-04). No cambia ningún hex, ninguna cifra de contraste ni ninguna decisión de §A–§K. `C2-colores.json` no se tocó.

## Hecho
1. **H-1 · C2-05.** `docs/design/tools/contraste_componentes.py` resuelve cada par **solo** con las referencias de `docs/contratos/C2-tokens.json`: `ui.<x>` → `ui.<x>.light/.dark`; «rol X» → `ui.role.X`; «superficie de rol X» → `ui.role.X-surface`; persona, tienda y serie → `ui.person|store|chart.*`; «acción crítica» → `ui.critical`. Sigue los alias `ui.*` hasta el primer primitivo `color.*` (columna `token_*`) y de ahí hasta el hex (columna `hex_*`). DESIGN §10.5 da solo la lista de pares. Las columnas claro/oscuro de §10.5, §7.1 y `C2-colores.json` quedan como cruces secundarios: avisan y dan código 1. Un lado que no se pueda resolver da código 2. El CSV regenerado sale **idéntico** (`--check`).
2. **H-2 · C2-04.** En DESIGN.md los 4 grupos se citan con sus hojas y su valor: `layout.topbar.height.mobile`/`.desktop`, `layout.preview.width.desktop`/`.wide`, `layout.dialog.width.short`/`.long` y `size.bottombar.item.height`/`.min-width` (§4.3, §4.6, §7.10, §7.12, §9.3, §10.2, §10.6). Se corrigió un 5.º caso del mismo defecto que encontró la verificación nueva: `icon.stroke.user-units` pasa a `.md`/`.sm` (§5.2). Además:
   - §10.1: el tema se llama `tailwind.theme.css` y está la cadena única.
   - §10.2: el conteo real es 75 semánticos (41 + 34) más `overlay.scrim`; se agregan `shadow.0` (se emite `none`) y el alfa en `$extensions`.
   - §10.3: 78 `--ui-*` por tema.
   - §10.4: describe la importación real. El tema importa solo `tokens.css` y la hoja del frontend importa primero `tailwindcss`; se listan las escalas vaciadas.
   - Supuesto 7 actualizado.
3. **M-1 · U2-11.** La barra inferior queda sin ambigüedad: **alto 56 px** = `layout.bottombar.height` y **ancho mínimo 64 px** (DESIGN §4.6, §7.10, §9.3 y Supuesto 10). En C2: `size.bottombar.item.height` = `{layout.bottombar.height}` y `size.bottombar.item.min-width` = 64px. Se quitó `min-height` 64px y `min-width` deja de valer 56px.
4. **M-2 · C2-10.** El generador emite `light`, `dark`, `surface-light` y `surface-dark` de los 8 roles y de la acción crítica como alias de los pasos 700/400/100/900 (36 alias). Antes de reemplazar, comprueba que el hex sea igual. `leer_hex` sigue los alias. `verificar_c2.py` se extendió a los colores: 0 hex literales repetidos sin alias y 0 semánticos literales; antes había 36 grupos. Las medidas no de color iguales se siguen listando como información (Supuesto 12 de U2c). DESIGN §10.2 «Reglas» dice exactamente eso.
5. **M-3 · C2-09/U2-04.** Hay una sola cadena: DESIGN §3–§10.6 + `C2-colores.json` → `generar_c2_tokens.py` → `C2-tokens.json` (no se edita a mano) → Style Dictionary → CSS. Está escrita en el README de tokens, en DESIGN §10.1 y en el docstring del generador. Las dos tablas oficiales de contraste quedan nombradas en DESIGN §1.2 P8, §2.8, §2.10 regla 5 y §10.5 y en el README: `contraste.csv` (paleta, `colorimetria.py`) y `contraste-componentes.csv` (componentes, `contraste_componentes.py`).
6. **Verificación nueva en `verificar_c2.py` (C2-04).** Falla si DESIGN cita como token un grupo sin `$value` cuyas hijas son hojas y que no es un par light/dark. Antes de corregir DESIGN detectaba los 5 casos; ahora detecta 0. C2-04 compara los hex siguiendo alias.
7. El historial de C2 suma la entrada «1.0 · corrección de la revisión U2 r1», con la versión 1.0, en el README y en `$extensions.finanzas.cambios`.

## Archivos
- `docs/design/tools/contraste_componentes.py` (modificado)
- `docs/design/tools/tokens/generar_c2_tokens.py` (modificado)
- `docs/design/tools/tokens/verificar_c2.py` (modificado)
- `docs/design/tools/tokens/README.md` (modificado)
- `docs/design/DESIGN.md` (modificado: 29 líneas cambiadas o agregadas)
- `docs/contratos/C2-tokens.json` (regenerado: 36 colores pasan a alias con el mismo hex; `size.bottombar.item.min-height` sale, entra `size.bottombar.item.height` y `min-width` va de 56 a 64px)
- `docs/contratos/tokens.css` (recompilado: los 36 alias salen como `var(--color-<grupo>-<paso>)`; `--size-bottombar-item-height: var(--layout-bottombar-height)` y `--size-bottombar-item-min-width: 64px`)
- `docs/contratos/tailwind.theme.css`: recompilado, idéntico byte a byte (sha256 `f66eb866…c88a`)
- `docs/design/contraste-componentes.csv`: regenerado, idéntico
- `docs/estado/wf1/correccion-U2-r1.md` (este reporte)

## Verificación
- `cd docs/design/tools/tokens && npm run check` → código 0. Pasos y resultados:
  - generate: escrito, 100073 B.
  - build: 0 advertencias, Style Dictionary 5.6.1.
  - verify: **RESULTADO: OK**.
  - check:sync: tokens.css `2b9429bc…708d` y tailwind.theme.css `f66eb866…c88a`, dos compilaciones iguales y coinciden con los versionados.
  - test:tailwind: tailwindcss 4.3.3 «Done».
- **C2-01**: compila con `log.warnings: 'error'` y `brokenReferences: 'throw'`, código 0.
- **C2-02**: 445 tokens; 0 sin tipo, 0 sin `$value` y 0 referencias rotas. Por tipo: color 342, dimension 79, fontWeight 4, shadow 3, number 3, duration 3, cubicBezier 2, fontFamily 1 y strokeStyle 8.
- **C2-03**: 76 semánticos de color (ui + overlay) con claro y oscuro, 0 incompletos; escalas 22, roles + crítica 135, categóricas 16 y series 16.
- **C2-04**:
  - 190 colores de `C2-colores.json`: 0 diferencias de hex (se siguen los alias) y 0 colores de más.
  - 158 pares «token · hex» de DESIGN, 0 diferencias.
  - 242 tokens citados, 0 inexistentes.
  - **0 grupos sin `$value` citados como token** (antes, 5).
  - `grep 'tailwind-theme\|56 × ≥ 64' DESIGN.md`: solo el Supuesto 10, que cita la expresión vieja a propósito.
- **C2-05 / U2-04**:
  - `python3 -I docs/design/tools/contraste_componentes.py --check` → «CSV en disco IDENTICO».
  - 158 pares, 316 combinaciones, 304 con umbral, 12 decorativas exentas y **0 bajo umbral**.
  - Cruce con contraste.csv: igual 266, nuevo 50 y difiere 0. Cruces secundarios: 0. Código 0.
  - `--independiente`: 10 razones con segundo cálculo, todas coinciden, código 0.
  - **Prueba de mutación**: copia de C2 con `ui.text.default.light` → `{color.neutral.800}`. Cambian 28 filas del CSV y salen 29 avisos de cruce (§7.1 y §10.5), con código 1. El CSV sigue a `C2-tokens.json`.
  - `verificar_c2.py` C2-05: 170 combinaciones, 0 bajo umbral.
- **C2-06**: `--ui-*` 78/78/78 (claro, media oscuro y data-theme); 0 diferencias; prefers-color-scheme 1, data-theme 2 y prefers-reduced-motion 1; 0 colores sueltos.
- **C2-07**: en `datos/tmp/c2-prueba.css` aparecen `#fafafb`, `#5d609f`, `#457723`, `var(--ui-role-income)`, `border-radius: 10px`, `var(--ui-shadow-1)` y `48rem` (todos ≥ 1; búsqueda sin distinguir mayúsculas). El archivo se borró.
- **C2-08**: 0 URL; las 2 rutas de `@font-face` existen y la licencia existe.
- **C2-09**: `npm run generate` dos veces → `cmp` idéntico. check:sync OK.
- **C2-10**: 190 primitivos, 36 de ellos alias de su paso; **0 hex repetidos sin alias** (con el C2 anterior el mismo conteo daba 36) y 0 semánticos literales. Quedan 18 grupos de medidas no de color iguales, como información. Quedan 6 `font.line-height.*` sin cita literal; ya estaban antes y no son de esta ronda.
- **U2-01 / U2-10**: `python3 -I docs/design/fuentes/scripts/verificar.py .` → 11 encabezados `##` en orden, 19 contrastes de §2.8 coinciden, 0 referencias a terceros, RESULTADO OK.
- **U2-11**: `size.touch.min` = 44px sin cambio. La barra inferior queda en 56 px de alto, ≤ el alto de la barra, y ≥ 64 px de ancho (5 × 72 px caben en 360).
- `python3 -I -m py_compile` de los 3 scripts → OK. Temporales `datos/tmp/corr-u2r1`, `c2-prueba.css` y `c2-sync` borrados.

## Supuestos
1. **M-1:** «56 × ≥ 64 px» se lee como alto 56 y ancho ≥ 64. La lectura contraria (alto mínimo 64) no cabe en `layout.bottombar.height` = 56px, así que no había dos opciones válidas. No se agrega ningún número que DESIGN no tuviera. Va como Supuesto 10 de DESIGN.
2. Se corrigió `icon.stroke.user-units` (§5.2) aunque el revisor no lo nombró. Es el mismo defecto de H-2 y la verificación nueva lo detecta.
3. La versión de C2 sigue en **1.0**, porque C2-11 la exige y C2 aún no pasa la puerta 2. El cambio queda en el historial como «1.0 · corrección de la revisión U2 r1».
4. Style Dictionary emite los alias de paso como `var()` en `tokens.css`. El valor calculado es el mismo y el tema de Tailwind no cambia.

## Bloqueos
Ninguno nuevo. No se creó ningún `docs/estado/bloqueos/U2-*` ni `C2-*`.

## Pendientes
- Hallazgos bajos de la revisión, fuera de esta ronda (solo se pidieron altos y medios), todos en `docs/design/DESIGN.md`:
  - B-1: §1.2 P5 cita §3.4, que no existe; la referencia correcta es §3.3.
  - B-2: §2.7 debe decir que el ΔL se mide en gris y en las simulaciones.
  - B-3: §2.2 dice ±30°, pero los ángulos reales son −25° y +35°.
- Ronda 2 del revisor: volver a correr `npm run check`, `python3 -I docs/design/tools/contraste_componentes.py --check` y `python3 -I docs/design/fuentes/scripts/verificar.py .`.

## Solicitudes fuera de alcance
1. `docs/design/U3-dashboard.md` línea 316: cita `size.bottombar.item` y `layout.topbar.height`, que son grupos sin `$value`. Cambio: el «+» de la barra es `size.bottombar.add` (40 px), o se citan `size.bottombar.item.height`/`size.bottombar.item.min-width`; en 768 corresponde `layout.topbar.height.desktop` (48 px). Motivo: C2-04, los mismos nombres desdoblados de esta corrección.
2. Agente paquete-claude-design (`docs/design/entrada/`): copiar el `tokens.css` nuevo. `--size-bottombar-item-min-height` ya no existe y los 36 primitivos de rol y de acción crítica usan `var()`, así que hay que copiarlo junto con sus fuentes, como pedía U2c.

## Paquetes instalados
Ninguno. Se usaron los `node_modules` existentes: style-dictionary 5.6.1 y tailwindcss 4.3.3.

## Propuestas para Axel
Ninguna. Esta corrección no toca valores K31.

## Consumo aproximado
- Acciones de herramienta: unas 38.
- Lectura: unos 195 KB, todo por partes. DESIGN.md (160 KB) nunca se leyó entero.
- Escritura propia: unos 45 KB entre scripts de parche, cambios y este reporte. Los archivos regenerados los escribieron los scripts.
