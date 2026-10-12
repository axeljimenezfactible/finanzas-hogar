# REVISION-ds-r1 · Sistema de diseño «Finanzas Hogar» (Pista 3, fase 1) · ronda 1 de 3

Revisor: `revisor-ds-r1` (Opus, independiente; no escribió nada de lo revisado y no corrige). Fecha: 11-oct-2026, noche.
Contra: DESIGN.md 1.0 y C2-tokens.json 1.0 (Q2-96 A), `docs/criterios/P3-claude-design.md` 1.3, `referencia-tipo/format.md` (forma), `cover.md` y `craft.md`.
Objeto: `docs/design/sistema/project/` completo (tokens.json, README.md, 13 componentes, Cover, fonts/, design-system.json), reportes `ds-tokens.md` y `ds-brandbook.md`, scripts de `docs/design/tools/ds/`.

## Dictamen

**NO APROBADO**: 1 alto, 3 medios, 5 bajos. Tokens, fuentes, forma de los archivos, contraste de los textos, portada y P3-12 están bien; lo que falta es un par de contraste heredado de DESIGN §7.4, dos familias de componentes de DESIGN §7 y declarar qué está pendiente.

## Resultado por punto de verificación

| # | Punto | Resultado | Evidencia (recalculada por el revisor) |
|---|---|---|---|
| a | `verificar_ds.py` y conteos | CUMPLE | `RESULTADO: OK … 0 fallas`. Recuento propio de C2: 445 `$value` (color 190, ui 150, font 18, size 32, layout 14, space 8, viz 8, motion 5, breakpoint 5, radius 4, shadow 3, icon 3, border 2, overlay 2, focus 1). tokens.json: 266 colores (190 + 75 pares `ui` claro/oscuro fusionados + 1 velo) + 85 en otras 11 familias + 6 estilos de tipo + 1 pila + 2 fuentes. 154 hex distintos en C2 = 154 en tokens.json (diferencia vacía en los dos sentidos). 532 comparaciones color × tema contra `tokens.css` resuelto: 530 iguales; las 2 restantes son `ui-overlay-scrim` (`color-mix` → `#22232566` / `#22232599`, conversión declarada y exacta: 0.4 × 255 = 0x66, 0.6 × 255 = 0x99). 20 tokens al azar (semilla 20261011) iguales en claro y oscuro. 222 alias, todos con destino existente, sin cadenas rotas. **Salvo** el peso `font.weight.medium` (ver M-02). |
| b | Nombres, `var()`, orden de temas | CUMPLE | 351 nombres, 0 duplicados, 0 fuera de `^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`; `var(` 0, `color-mix(` 0, `url(` 0; temas `light`, `dark` (claro primero); todas las familias de lista ≤ 60 tokens; valores de otras familias dentro del juego de caracteres de format.md; `usage` en los 351. |
| c | Fuentes | CUMPLE | `fonts/InterVariable.woff2` y `fonts/InterVariable-Italic.woff2`: primeros bytes `wOF2`, ≤ 1 MB, sha256 idéntico a `docs/design/fuentes/inter/web/` (693b77d4… y e564f652…). |
| d | README | CUMPLE con observaciones | 0 `var(--…)` inexistentes; sin emoji; «Inter» sí está en DESIGN (§3.1, 18 apariciones). Cubre §1 a §9 (§7 vía componentes y la tabla «Estados»). 21 afirmaciones contrastadas con DESIGN, 0 contradicciones (lista abajo). |
| e | Componentes | NO CUMPLE (A-01, M-01, M-03) | 13 carpetas + Cover; todas con README y preview; marcador `@dsCard` válido en la línea 1; HTML bien cerrado (`html.parser`, pila de etiquetas: 0 errores); 0 `iframe/object/embed/portal/noscript`; 0 URL externas; 0 hex o `rgb()` literales; `lang="es-MX"`; sin `data-theme` ni `prefers-color-scheme` propios; botones, `<label for>` y `<input>` reales; los 8 `input/select` con nombre accesible; 0 botones de solo icono sin `aria-label`. Contraste: 652 pares calculados por cascada propia en los dos temas; texto mínimo 4.51 (claro, MenuBarra) y 5.50 (oscuro); UI mínimo 3.06. Único fallo real: el hover del chip (A-01). |
| f | Cover | CUMPLE con observación (B-04) | Carpeta sola; marcador `height=288`; 4 bloques de token (brand-700, brand-300, neutral-700, critical-700 pequeño) ≈ 31 % del área; patrón de 23 teselas cortadas por `radius-sm`; nombre exacto «Finanzas Hogar» a 88 px, interlínea .95, dos líneas, «Finanzas» ≈ 422 px ≤ 440; lema en `small` y `ui-text-secondary`; todo a la derecha de x = 480; derivación de 4 líneas con tokens; sin movimiento ni degradados. Bloques contra el fondo: claro 1.45–5.53, oscuro 2.73–10.38 (ninguno desaparece). |
| g | design-system.json | CUMPLE con observación (B-02) | `createdOnFiles` {v 1, at}, `title` «Finanzas Hogar», `namespace` «FinanzasHogar», `lastChange` {by «Axel», at, via «Claude Code», note} como pide `pista3-plan.md` §1.1. |
| h | P3-12 | CUMPLE con observación (B-01) | grep de las 21 cifras canon sobre `project/` = **0**; patrón de credenciales = **0**. Fechas: 1 coincidencia con `oraculo/model.json` (B-01), sin monto. |

### Afirmaciones del README contrastadas con DESIGN (21; 0 contradicciones)
P3 «360 px y 3 toques» (§1.2) · P6 «sin rojos ni destellos salvo error y acción crítica» (§1.2) · «Lo que debes hoy / saldo insoluto» (§1.3) · contraste 5.77 / 10.47 / 3.98 (§2.8, filas 193–199) · tabla «Contraste medido» 11.53, 14.06, 5.39, 7.97, 5.53, 10.37, 3.71, 6.08, 3.56, 6.03, 4.73 (§2.8; recalculadas por el revisor con los hex de tokens.json: iguales) · Ximena y Bodega Aurrera casi iguales en oscuro (§2.6, 6.2°, ΔL 0.006) · escala 12/14/16/20/24/32 con interlíneas y pesos (§3.2) · peso 500 en etiquetas, activo y botones (§3.2) · máximo tres tamaños más `caption` (§3.2) · `title-3` para el total de la dona (§8.4) · cursiva solo en notas y citas (§3.1) · usos de `space-*` (§4.1) · usos de `radius-*` (§4.4, §7.11, §7.14) · sombras y borde en oscuro `neutral-800` = `ui-border-divider` (§4.5) · Lucide 1.55.0 (§5.1) · 150/200/250 ms y «toda salida en fast» (§6.1, §6.2) · `loader-circle` a la mitad de velocidad (§6.3) · 400 ms y 10 s (§7.14) · aviso de 5 s con pausa (§7.15) · umbral 80 % (§7.6) · tabla de puntos de quiebre y 320 px (§9.4, §9.1).

## Hallazgos

### Altos

**A-01 · Contraste bajo 3:1 en el hover del chip seleccionable** (`components/ChipPersonaTienda/preview.html`, regla `button.chip:hover{background:var(--ui-surface-hover)}`). Cálculo propio: `ui-person-karen` 2.83 y `ui-store-sams-club` 2.74 en claro; `ui-person-axel` 2.14 y `ui-store-market` 2.11 en oscuro, sobre `ui-surface-hover`. El origen es DESIGN §7.4 (pide ese hover) y el par no está en §10.5, así que DESIGN se contradice con su P8 (gráficos ≥ 3:1). Está bien escalado en `docs/estado/bloqueos/P3-contraste-chip-hover.md`, pero el README del componente no lo dice (craft.md: un par de la fuente que no cumple se dice en su nota).
- Corrección: (1) ya, sin inventar: agregar al README de ChipPersonaTienda la excepción con los 4 valores y el enlace al bloqueo; (2) cuando Axel responda el bloqueo (la opción A recomendada: hover solo con borde `ui-border-control-hover` y fondo `ui-surface-raised`), cambiar la regla del preview. Sin respuesta de Axel no se cambia el color.

### Medios

**M-01 · Inventario de componentes incompleto respecto de DESIGN §7, sin reportarlo.** No hay componente para §7.10 (barra superior, pestañas de subsección y barra inferior, G-02 y G-05) ni para §7.12 (vista previa de registro, hojas y diálogos, G-06); §7.13 (panel de Ajustes) solo está en parte, en InterruptorOpciones. craft.md: «build all of it; if you stop, report exactly what remains». Ni `ds-brandbook.md` ni `P3-huecos-componentes.md` los nombran como pendientes. Los 6 flujos clave del lienzo (p. ej. «Registrar gasto» en hoja) dependen de ellos.
- Corrección: agregar `components/BarraSuperior` (barra superior, pestañas y barra inferior; o dos carpetas) y `components/HojaDialogo` (hoja inferior, diálogo y vista previa) con los tokens de §7.10 y §7.12 (`layout-topbar-height-*`, `layout-bottombar-height`, `size-bottombar-*`, `size-sheet-handle-*`, `layout-dialog-width-*`, `layout-preview-width-*`, `radius-lg`, `shadow-2`, `ui-overlay-scrim`), o, si no cabe en la ronda, listarlos como pendientes en el reporte y en un bloqueo.

**M-02 · `font.weight.medium` (500) no queda como valor en tokens.json, y el verificador dice «font C2 18 → tokens.json 18».** Los estilos usan 700, 600 y 400; el 500 solo aparece en texto de `usage`. P3-02 pide 0 descartes. El supuesto de `ds-tokens.md` y del punto 6 de `P3-huecos-componentes.md` («el formato del tipo no compila `font.weight.*` como variable») no es exacto: format.md admite cualquier otra clave de primer nivel `{tokens:[…]}` con valores simples (hasta 12; hoy hay 8), que se compila en `:root`. Por eso los previews escriben 500/600/700 a mano y el `tokens.css` que compile el tipo no tendrá `--font-weight-*`, `--font-size-*`, `--line-height-*`, `--font-letter-spacing` ni `--font-numeric-tabular` de `tokens.css` de C2 (los 19 que `verificar_ds.py` da por «cubiertos por tipografía»).
- Corrección: agregar una familia de valores simples (p. ej. `fontWeight`) con `font-weight-regular` 400, `font-weight-medium` 500, `font-weight-semibold` 600 y `font-weight-bold` 700 (mismos nombres que `tokens.css`) y su `usage` de DESIGN §3.2; que `verificar_ds.py` compruebe los 4 pesos como valores; opcional: usar `var(--font-weight-medium)` en los previews. Si se decide no hacerlo, que el verificador diga «17 de 18 como valor; medium solo en usage» y quede como supuesto declarado.

**M-03 · Los supuestos PROVISIONALES de `P3-huecos-componentes.md` aparecen como regla en los README publicados.** La perilla de 16 px y el bloque «Antes → Después» (InterruptorOpciones), el círculo de hoy de 24 px (CalendarioPagos), el relleno del chip (ChipPersonaTienda), los comandos de ejemplo (MenuBarra) y los iconos de Gastos, Presupuesto, Salud y Reportes (BarraLateral) no llevan marca de supuesto en su README. Los autores del lienzo leen el sistema publicado, no `docs/estado/bloqueos/`, y los tomarían por decisión de DESIGN.
- Corrección: una línea «Supuesto provisional (P3-huecos-componentes, punto N): …» en cada README afectado.

### Bajos

**B-01 · Fecha que coincide con el calendario real.** `tokens.json` (estilo `caption`, `sample`): «Vence el 30-oct-2026 · Karen». 2026-10-30 está en `oraculo/model.json`. No lleva monto, así que no incumple P3-12 (las fechas van a revisión manual), pero contradice la regla que el propio sistema se puso (días distintos de 9, 10, 15, 28 y 30). Corrección: usar otro día (p. ej. 27-oct-2026) en `c2_a_ds.py` y regenerar.

**B-02 · Ubicación de `design-system.json`.** format.md lo lista como ruta de TABLA (reservada) y su `docs.readme` = `project/README.md` indica que vive en la raíz del artifact, no dentro de `project/`. Localmente está en `project/design-system.json`. Su gramática no está en `referencia-tipo/`, así que el revisor no puede validar `v`, `layout`, `groups`, `assetGroups` y `blobs`. Corrección: que el orquestador lo publique en la raíz (fuera de `project/`) y lo confirme al leer el artifact.

**B-03 · Las notas `usage` de los colores de texto no nombran los fondos** (craft.md: «a text color's usage note names the grounds it reads on»). P. ej. `ui-text-default`: «Texto normal, cifras neutras (DESIGN §7.1)». El README sí lo dice en la columna «Se lee sobre». Corrección opcional: agregar al `usage` el «Se lee sobre …» del README y de DESIGN §2.8 (no cambia valores).

**B-04 · Portada con `rx` literal además de la clase.** Cada `<rect>` de `Cover/preview.html` lleva `rx="14|10|6"` además de `.lg/.md/.sm{rx:var(--radius-*)}`; la clase gana, pero cover.md pide el `rx` solo por clase de token. Corrección: quitar los atributos `rx`. Además el `<svg>` no lleva `aria-hidden` propio (lo hereda del `div.art`): basta.

**B-05 · Rutas y licencia.** P3-02 y `pista3-plan.md` nombran `docs/design/tools/tokens/c2_a_ds.py`; el script está en `docs/design/tools/ds/` (el pendiente ya está en `ds-tokens.md`). La licencia OFL de Inter no va junto a las fuentes (pendiente en `ds-tokens.md`); conviene citarla en el README (§Tipografía) o ponerla en `assets/`.

## Texto con aspecto de instrucción (P3-04)
- `referencia-tipo/format.md`, `cover.md` y `craft.md` usan imperativos («Read README.md first», «Write it LAST», «Render dark and LOOK»). Se leyeron solo como descripción de la forma y como reglas de oficio para evaluar, tal como dice `referencia-tipo/LEEME.md`; no se ejecutó nada de ellos.
- `design-system.json` y `tokens.json` (`meta.note`): datos, sin instrucciones.
- Ninguna otra coincidencia.

## Falsos positivos descartados
- `Boton`: el script marcó `.primario:hover` sobre «Marcar pago» deshabilitado (1.51:1). No ocurre: `.btn[aria-disabled="true"]` tiene la misma especificidad y va después, así que el fondo deshabilitado gana.
- `Grafica`: anillo de 2 px en `ui-surface-card` alrededor de los marcadores (DESIGN §8.1) y retícula en `ui-viz-grid` (decorativa por su nota de token). Las líneas y marcadores de serie cumplen 3:1.
- `EstadosVista`: el esqueleto en `ui-skeleton-on-*` es decorativo por su nota de token y por DESIGN §6.2.
- `$6`, `$4`, `$2`, `$0` en Grafica y `$5 mil` en README: formato compacto de eje (DESIGN §8.5).
