# REVISION-ds-r2 · Sistema de diseño «Finanzas Hogar» (Pista 3, fase 1) · ronda 2 de 3

Revisor: `revisor-ds-r2`, independiente: no escribió nada de lo que revisa y no corrige nada. Fecha: 11-oct-2026, noche.

Se evalúa contra:
- DESIGN.md 1.0 y C2-tokens.json 1.0 (Q2-96 A);
- `docs/criterios/P3-claude-design.md` 1.3;
- `referencia-tipo/format.md` (forma), `cover.md` y `craft.md`.

Se revisó:
- `docs/design/sistema/project/` completo: tokens.json, README.md, 18 componentes, Cover, fonts/ y design-system.json;
- los reportes `ds-tokens.md`, `ds-brandbook.md` y `correccion-ds-r1.md`;
- los scripts de `docs/design/tools/ds/`;
- el respaldo previo a la corrección, en `/tmp/correccion-ds-r1-respaldo/`.

Cada cifra de este informe sale de scripts propios, escritos en `datos/tmp/rev-ds-r2/` y borrados al terminar.

## Dictamen

**NO APROBADO: 1 hallazgo alto, 1 medio y 9 bajos** (4 nuevos y 5 que vienen de la r1 sin atender).

Los 4 hallazgos de la r1 (A-01, M-01, M-02 y M-03) quedan **cerrados**. Hay dos hallazgos nuevos que bloquean:
- **A-01 (alto):** en `BarraSuperior`, los iconos de los botones de solo icono miden 20 px. DESIGN pide 24 px.
- **M-01 (medio):** DESIGN §7.5.1 (barra de la vista) y §7.5.5 (galería) no tienen componente, y nada los reporta como pendientes.

## Estado de los hallazgos de la r1

| r1 | Estado | Evidencia recalculada |
|---|---|---|
| A-01 · contraste del chip en hover | **CERRADO** en lo que tocaba al sistema; el resto depende del bloqueo, pendiente de Axel | Recalculé las marcas sobre `ui-surface-hover`: Karen 2.83 y Sam's Club 2.74 en claro; Axel 2.14 y Mercado 2.11 en oscuro. Son los mismos valores del README y de `P3-contraste-chip-hover.md`. La afirmación del README «sobre raised y selected las 8 cumplen» es cierta: el mínimo es 3.05 en claro y 3.06 en oscuro. El README declara la excepción, como pide craft.md («a real source's pair that misses stays exact — say so in its note»). La regla `button.chip:hover` sigue como la pide DESIGN §7.4 hasta que Axel responda. En la cascada del preview solo aparece el par de Karen, porque es el chip seleccionable que se dibuja. |
| M-01 · componentes de §7.10, §7.12 y §7.13 | **CERRADO** | `BarraSuperior`, `BarraInferior`, `HojaDialogo`, `VistaPrevia` y `PanelAjustes` existen, cada uno con README y preview. Se contrastaron con DESIGN: 12 de las afirmaciones de abajo son de ellos. Se abren dos hallazgos aparte: el tamaño de icono (A-01 de la r2) y el inventario de §7.5 (M-01 de la r2). |
| M-02 · `font.weight.medium` como valor | **CERRADO** | La familia `fontWeight` tiene `font-weight-regular` 400, `-medium` 500, `-semibold` 600 y `-bold` 700, iguales a `tokens.css`. Los 14 previews anteriores cambian **solo** sus pesos: hice el diff contra el respaldo, normalizando los pesos, y salen 0 líneas distintas en los 14. `tokens.json` sin `fontWeight` queda idéntico al respaldo. Tamaños e interlíneas siguen en `type.groups[].styles`, que es lo que piden format.md («reads sizes, weights and line heights from styles only») y from-code.md. Ningún preview usa `--font-size-*`, `--line-height-*` ni `--font-numeric-tabular`. |
| M-03 · supuestos sin marcar | **CERRADO** | Hay líneas «Supuesto provisional (P3-huecos-componentes, punto N)» en BarraLateral (1), InterruptorOpciones (2 y 3), CalendarioPagos y ChipPersonaTienda (4), MenuBarra (5), BarraSuperior y BarraInferior (7), HojaDialogo y VistaPrevia (8) y PanelAjustes (9). Los puntos 7 a 9 ya están en el bloqueo. |

## Resultado por punto de verificación

| # | Punto | Resultado | Evidencia propia |
|---|---|---|---|
| a | `verificar_ds.py` y conteos | CUMPLE | <ul><li>`python3 -I docs/design/tools/ds/verificar_ds.py` termina en `RESULTADO: OK … 0 fallas`.</li><li>Recuento propio de C2: 445 `$value`. Mi recorrido da 444 más el `$root` de `icon.stroke`. Por familia: color 190, ui 150, font 18, size 32, layout 14, space 8, viz 8, motion 5, breakpoint 5, radius 4, shadow 3, icon 3, border 2, overlay 2 y focus 1.</li><li>tokens.json tiene 266 colores (190 + 75 pares `ui` + 1 velo) y 89 tokens en otras 12 familias, más 6 estilos de tipo, 1 pila y 2 fuentes.</li><li>Comparé los 266 colores × 2 temas (532) contra `tokens.css`, resolviendo `var()` y `color-mix()`: **0 diferencias**.</li><li>Muestra al azar de 20 tokens (semilla 20261012), en claro y oscuro: 0 diferencias.</li><li>222 referencias de alias, 0 rotas, 0 valores finales que no sean hex.</li><li>Conjunto de hex: 154 en C2 y 154 en tokens.json; la diferencia está vacía en los dos sentidos.</li><li>Otras familias: `space-*` en px es igual al `$value` de C2. `tokens.css` usa rem, pero es la misma cifra. `shadow-1` y `shadow-2` en oscuro valen `none`, como DESIGN §4.5.</li></ul> |
| b | Nombres, `var()`, orden de temas | CUMPLE | <ul><li>355 nombres en listas y 6 estilos: 0 repetidos y 0 fuera de `^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`.</li><li>`var(` 0, `color-mix(` 0 y `url(` 0.</li><li>Temas `light` y `dark`, en ese orden.</li><li>`usage` en los 355.</li><li>Valores de las otras familias dentro del juego de caracteres de format.md, con paréntesis balanceados.</li><li>Todas las familias tienen 60 tokens o menos.</li><li>Tope de familias: ver B-02.</li></ul> |
| c | Fuentes | CUMPLE | `fonts/InterVariable.woff2` (352,240 B) y `fonts/InterVariable-Italic.woff2` (387,976 B): empiezan con `wOF2` y su sha256 es idéntico al de `docs/design/fuentes/inter/web/` (693b77d4f32e… y e564f652916d…). `type.fonts[].file` apunta a los dos. |
| d | README | CUMPLE con observaciones (A-01, B-01) | <ul><li>Sin título `#`, sin emoji.</li><li>122 `var(--…)` distintos: 0 inexistentes. Nombra los 6 estilos.</li><li>Cubre DESIGN §1 a §9. §7 se cubre con la tabla «Estados» y con los componentes.</li><li>«Inter» sí está en DESIGN §3.1.</li><li>17 afirmaciones contrastadas con DESIGN (lista abajo): 15 coinciden, 1 queda incompleta en su redacción (A-01) y 1 es una inferencia sin marcar (B-01).</li></ul> |
| e | Componentes | NO CUMPLE (A-01, M-01) | <ul><li>18 carpetas más Cover, todas con README y preview.</li><li>Marcador válido en la línea 1 de los 19 previews.</li><li>HTML bien cerrado según `html.parser` con pila de etiquetas: 0 errores.</li><li>0 `iframe`, `frame`, `object`, `embed`, `portal` o `noscript`; 0 `http(s)://`; 0 `<script>` y 0 `<link>`.</li><li>0 hex o `rgb()` literales; 0 `data-theme` o `prefers-color-scheme` propios; `lang="es-MX"` en los 19.</li><li>Los 143 `var(--…)` distintos de los previews existen.</li><li>0 botones o elementos con rol de botón sin nombre accesible.</li><li>Todo `input`/`select` tiene `<label for>` o `aria-label`.</li><li>Los 101 iconos de 24 × 24 coinciden **byte a byte** (sin contar espacios) con un `<symbol>` de `fuentes/lucide/sprite.svg`.</li><li>Contraste: mi propia cascada (selectores, especificidad, herencia, composición alfa) recorre los 19 previews × 2 temas en reposo, hover, presionado y foco. Son **2,466 pares**:<ul><li>texto: mínimo 4.51 en claro (MenuBarra) y 5.04 en oscuro;</li><li>interfaz y gráfico: mínimo 3.28 en claro y 3.06 en oscuro, sin contar la excepción declarada del chip (2.83).</li></ul></li><li>0 fallas nuevas.</li></ul> |
| f | Cover | CUMPLE con observación (B-04 de la r1) | <ul><li>La carpeta solo tiene `preview.html`.</li><li>`height=288`.</li><li>Un SVG en línea con clases ligadas a tokens (`color-brand-700`, `-brand-300`, `color-neutral-700`, `color-critical-700`; `rx` por `radius-*`).</li><li>Derivación de 4 líneas con tokens; ≈ 31 % del área.</li><li>Teselas a paso `space-6`.</li><li>Nombre exacto «Finanzas Hogar» a 88 px, interlínea .95, dentro de 440 px. Bloques a la derecha de x = 480.</li><li>Sin `animation`, `transition` ni degradados.</li><li>Siguen los `rx` literales (B-04 de la r1).</li></ul> |
| g | design-system.json | CUMPLE con observación (B-02 de la r1) | `createdOnFiles` {v 1, at 2026-10-12T04:57:11Z}, `title` «Finanzas Hogar», `namespace` «FinanzasHogar», `lastChange` {by, at 2026-10-12T05:20:19Z, via, note de la corrección}. |
| h | P3-12 | CUMPLE con observación (B-01 de la r1) | <ul><li>`grep -rnoE` de las 21 cifras canon sobre `project/`: **0**.</li><li>Patrón de credenciales: **0**.</li><li>10 fechas en `project/`; solo 2026-10-30 coincide con `oraculo/model.json` (muestra de `caption`, sin monto).</li><li>41 montos distintos: ninguno con centavos aparece en `oraculo/model.json` ni en `datos/semilla/data.json`. `$0`, `$2`, `$4`, `$5` y `$6` son marcas de eje (DESIGN §8.5).</li></ul> |

### Afirmaciones contrastadas con DESIGN (17)

1. **BarraSuperior.** Alto `layout-topbar-height-mobile` 56 px y `-desktop` 48 px (§7.10). Coincide.
2. **Pestañas.** Van en `small`. Reposo `ui-text-secondary`, hover `ui-text-default`. La activa va en `ui-text-default` con `font-weight-medium` y un indicador de 2 px en `ui-border-selected`, con `role="tablist"` (§7.10). Coincide.
3. **BarraInferior.** 5 destinos en orden fijo, con `house`, `utensils`, `plus`, `shopping-cart` y `search` (§7.10; los iconos se comprobaron en el sprite). Coincide.
4. **Destinos de la barra inferior.** Cada uno mide 56 px de alto y al menos 64 de ancho. La píldora mide 56 × 32 sobre `ui-surface-selected`. El «+» es un círculo de 40 px en `ui-action-primary-bg` con `plus` en `ui-text-on-primary` (§7.10). Coincide.
5. **Barra inferior en general.** Fondo `ui-surface-raised` con borde superior `ui-border-divider`, icono `size-icon-md` y nombre en `caption` (§7.10). Coincide.
6. **Hoja.** Asa de 32 × 4 en `ui-border-control`, título en `title-2` y primaria de 48 px. `radius-lg` solo arriba, `shadow-2` y velo (§7.12). Coincide.
7. **Diálogo.** Ancho de 480 o 640. Cuerpo en `ui-text-secondary`. Acciones a la derecha, con la primaria al final (§7.12). Coincide.
8. **Acción crítica.** `shield-alert` en `ui-critical`, que resuelve a `color-critical-light`/`-dark` (C-139). El foco inicial va a «Cancelar», de texto. El botón dice «Borrar 12 gastos» (§7.12 y §7.2). Coincide.
9. **VistaPrevia.** `ui-surface-raised`, `shadow-1`, `radius-md` (§4.4) y anchos de 400/480. Lleva 4 propiedades más «Mostrar N más», y `x` con `aria-label` (§7.12). Coincide.
10. **PanelAjustes, secciones.** Las 14 secciones van en orden, con `user-cog` … `package` (§5.4, filas «Ajustes ›»; iconos comprobados en el sprite). Coincide.
11. **PanelAjustes, filas.** Muestra «Valor inicial: 80 %». Muestra «Restablecer al valor inicial» con `rotate-ccw` solo si el valor difiere, y su nombre accesible dice qué restablece. «Guardado» lleva `circle-check` en `ui-role-success`, que resuelve a `color-role-success-light`/`-dark` (§7.13). Coincide.
12. **PanelAjustes, medidas.** Lista de 240 px sobre `ui-surface-page` y contenido a `layout-reading-max-width`, 720 (§7.13 y §4.2). Coincide.
13. **README, rejilla.** 4, 8 y 12 columnas, con cortes en 768 y 1024 (§4.2). Coincide.
14. **README, sombras.** `shadow-0` en tarjetas. `shadow-1` en menús, en la tarjeta que se arrastra y en la vista previa. `shadow-2` en diálogos, hojas, menú `/` y avisos. En oscuro no hay sombra: va un borde `neutral-800`, que es `ui-border-divider` (§4.5). Coincide.
15. **README, movimiento.** Usos de `duration-fast`, `-base` y `-slow`, incluido el «nodo de la barra lateral» en `base`. Entradas con `ease-emphasized` (§6.1 y §6.2). Coincide.
16. **README, iconografía:** «`size-icon-md` junto a títulos y solo en barras». **Queda incompleta.** DESIGN §5.2 dice «junto a títulos de 20 px o más y solos en barras: barra superior, barra inferior, botones de solo icono». Así escrito, «solo en barras» se lee como «únicamente». Ver A-01.
17. **VistaPrevia README:** hover de solo borde sobre lo que flota. **Es una inferencia**: §7.3 solo fija el hover sobre la página y sobre una tarjeta. Ver B-01.

## Hallazgos

### Altos

**A-01 · El tamaño de icono de los botones de solo icono de la barra superior no es el de DESIGN** (`components/BarraSuperior/preview.html`).
- Los botones «⋯» (`aria-label="Más acciones"`, en móvil y en escritorio) y Buscar llevan `svg.icon.i20`, es decir `var(--size-icon-sm)`, 20 px, con trazo `sm`.
- DESIGN §5.2 fija `size.icon.md`, 24 px, «solos en barras: barra superior, barra inferior, botones de solo icono».
- DESIGN §7.2 (Tamaños) dice: «De icono: … `size.icon.sm`, o `size.icon.md` en barras».
- `BarraInferior` sí usa `i24`. El README de BarraSuperior no dice el tamaño, y el README del sistema (Iconografía) lo redacta incompleto (afirmación 16).
- Corrección:
  - en el preview, cambiar `i20` por `i24` en los tres `.btn-icono` de la barra, lo que lleva `stroke-width: var(--icon-stroke-user-units-md)`;
  - agregar al README de BarraSuperior «iconos de solo icono en `var(--size-icon-md)`» y sumar `var(--size-icon-md)` a sus Tokens;
  - en el README del sistema, escribir «`var(--size-icon-md)` junto a títulos de 20 px o más y en los botones de solo icono de la barra superior y la inferior».
- Las migas y el icono de «+ Registrar» van junto a texto `small`, así que se quedan en `sm`.

### Medios

**M-01 · Inventario de §7.5 incompleto, sin reportarlo** (`components/`).
- DESIGN §7.5.1, la barra de la vista, no tiene componente: selector Tabla · Tablero · Calendario · Galería con `table-2`, `kanban`, `calendar` y `layout-grid` y `aria-pressed`; «Filtrar» (`list-filter`), «Ordenar» (`arrow-up-down`) y «Agrupar»; en móvil, un botón «Filtrar» que abre una hoja, y los filtros activos como chips quitables.
- DESIGN §7.5.5, la galería, tampoco tiene componente: portada 4:3, o el icono de la entidad sobre `ui-surface-hover`; tarjeta `ui-surface-card`, `radius-md` y `space-3`; hover con borde `ui-border-control-hover`; seleccionada en `ui-surface-selected` con `ui-border-selected`; «Foto no disponible»; columnas por ancho.
- El calendario genérico de §7.5.4 solo aparece a través de `CalendarioPagos` (§7.8).
- La fila «Tabla, tablero, calendario, galería» de la matriz §7.16 queda cubierta solo en parte. Aun así, `correccion-ds-r1.md` afirma «Ahora quedan cubiertas todas las filas de la matriz de DESIGN §7.16», y ningún bloqueo lo lista. craft.md: «build all of it; if you stop, report exactly what remains».
- Lo pasó por alto también la r1.
- Corrección:
  - agregar `components/Galeria` y la barra de la vista, sea como `components/BarraVista` o dentro de `TablaBD`, con su README, solo con los tokens citados;
  - decir en el README de `CalendarioPagos` si representa también §7.5.4;
  - si no cabe en la ronda, listarlos como pendientes en `P3-huecos-componentes.md` y en el reporte.

### Bajos (nuevos)

**B-01 · Inferencia sin marcar en VistaPrevia** (`components/VistaPrevia/README.md`).
- El README dice «edición en el lugar con hover de solo borde `var(--ui-border-control-hover)` (sobre lo que flota, como sobre una tarjeta, DESIGN §7.3)».
- §7.3 fija el hover solo sobre la página (`hover-subtle`) y sobre una tarjeta (borde). Sobre `ui-surface-raised` es una analogía, que no está en el punto 8 del bloqueo.
- El contraste cumple: lo calculé.
- Corrección: pasarlo a su línea «Supuesto provisional (punto 8)» y al bloqueo.

**B-02 · Tope de familias, contado distinto en cada archivo** (`tokens.json` y `tools/ds/verificar_ds.py:126`).
- `tokens.json` tiene 13 familias de lista más `type`; sin contar `color` ni `type` son 12.
- format.md dice «ANY OTHER top-level `{tokens:[…]}` key (…, ≤12)», y en *caps*, «≤12 families».
- `verificar_ds.py` excluye `color` y `type`, así que da 12: justo en el tope.
- `correccion-ds-r1.md` decía que ya había 12 familias de lista contando `color`, y aun así agregó `fontWeight`.
- En la lectura más coherente (como mucho 12 claves «otras») cumple, pero sin margen.
- Corrección: al publicar, el orquestador lee el `tokens.css` compilado y confirma que existen `--viz-*`, `--breakpoint-*`, `--motion`/`--duration-*` y `--font-weight-*`. No agregar más familias. Unificar el conteo en los reportes.

**B-03 · Foco de las pestañas** (`components/BarraSuperior/preview.html`).
- `.tab:focus-visible{outline-offset:-2px}` dibuja un anillo interior, con número literal. DESIGN §7.10 pide «foco con anillo común»; el anillo interior es solo para la barra inferior.
- Lo más probable es que el `overflow-x:auto` de la fila recorte el anillo exterior.
- Corrección: declararlo como supuesto con esa razón, o dar a la fila el aire que necesita el anillo común.

**B-04 · Tokens con otro significado en marcos de ejemplo** (`components/BarraSuperior/preview.html`).
- `.escritorio .marco{max-width:var(--layout-dialog-width-long)}` usa el ancho de diálogo como ancho del marco de muestra.
- Quien copie el preview al lienzo arrastrará un uso que no es el del token.
- Corrección: usar un número de maqueta, o el ancho real de la barra desde 768.

### Bajos que vienen de la r1, sin atender (no estaban en el encargo de la corrección)

- **B-01 de la r1:** la muestra de `caption` dice «Vence el 30-oct-2026», una fecha que coincide con `oraculo/model.json`.
- **B-02 de la r1:** `design-system.json` está dentro de `project/`. format.md lo trata como ruta de tabla reservada; debe publicarse en la raíz.
- **B-03 de la r1:** el `usage` de los colores de texto no nombra los fondos sobre los que se lee.
- **B-04 de la r1:** la portada lleva `rx` literal además de la clase.
- **B-05 de la r1:** P3-02 cita la ruta `tools/tokens`, pero el script está en `tools/ds`. Falta la licencia OFL junto a las fuentes.

## Texto con aspecto de instrucción (P3-04)

- format.md, cover.md, craft.md y from-code.md tienen frases en imperativo. Se leyeron solo como forma y reglas de oficio; no se ejecutó nada de ellos.
- `meta.note` de tokens.json y `lastChange.note` de design-system.json son datos.

## Falsos positivos descartados

- **Contraste del chip:** mi cascada solo detecta el par de Karen en claro, porque el preview dibuja ese chip seleccionable. Los otros 3 pares de la excepción se comprobaron aparte, con los hex de tokens.json.
- **`shadow-1` y `shadow-2` en oscuro:** valen `none` en tokens.json, mientras que `--shadow-*` de `tokens.css` no cambia en oscuro. Coincide con DESIGN §4.5 y con la copia temática `--ui-shadow-*` que declara el verificador. No es un valor distinto.
- **«Descartar cambios» como primario en el diálogo de confirmación:** DESIGN §7.2 reserva el destructivo para «borrar un registro o un renglón».
- **Ancho fijo en la hoja:** `min-width:280px` y `height:400px` de `HojaDialogo` son medidas del marco de muestra, no reglas del componente.
