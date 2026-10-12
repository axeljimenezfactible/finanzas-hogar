# P3-iter-3 · ajustes-b · revisión de la vuelta 3

Fecha: 12-oct-2026 · Revisor: `revisor-ajustes-b-r3` (independiente, vuelta 3 de 3; no corrige) · Criterio: P3-07 de `docs/criterios/P3-claude-design.md` 1.3.
Alcance: 14 artboards de `docs/design/maquetas/lienzo/project/ajustes/` (P-26.08 a P-26.14, `-m` y `-e`) y el fragmento `docs/design/maquetas/lienzo/fragmentos/ajustes-b.json`, como los dejó `correccion-ajustes-b-r2`.
Contra: P3-08 y P3-12; U5-04, U5-06 y U5-07 (lo verificable en el markup); DESIGN.md (§3.3, §5.2, §7.2, §7.3, §7.9, §10.5); U1 §2.1, §2.2 y §2.3; U4 §0 y P-26.08 a P-26.14; `referencia-tipo/format.md`.
Ruta del informe: la que fijó la tarea (`docs/estado/wfp3/`). La copia a `docs/design/salida/iter-3/` (P3-06) la hace el orquestador.

**Veredicto: NO APROBADO** · 1 alta · 1 media · 18 bajas.
Hallazgos de la vuelta 2: R2-A1 **cerrado**, R2-M1 **cerrado** y H-11 **sigue abierto** (`datos-sinteticos.json` no cambia desde el 11-oct 23:55). La alta es nueva en severidad, no en contenido: es la antigua B-07 (marcadores de posición sin regla `::placeholder`). Las vueltas 1 y 2 la dejaron en baja sin medirla. En esta vuelta se midió y queda por debajo de 4.5:1.

## Verificación ejecutada (scripts propios, `python3 -I`, en `datos/tmp/rev-ajb-r3/`, borrados al final)

El motor de cascada propio aplica las reglas del `helmet` (selectores de clase, etiqueta, atributo, `>` y `+`, con especificidad y orden) y el `style` en línea. También resuelve `var()` contra `tokens.css` en claro y oscuro y calcula herencia de `color`, fondo efectivo, `data-theme` local y opacidad acumulada.

| Comprobación | Resultado |
|---|---|
| (a) Forma: `<!doctype html>`, `<html lang="es-MX">`, `<meta charset="utf-8">`, `<title>` propio, `<script src="./support.js"></script>` y `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">` una vez cada uno en `<head>`, `<body><x-dc><helmet><style>`, un solo `x-dc`, bloque `data-dc-script` justo después de `</x-dc>` con `class Component extends DCLogic` y `renderVals()`, sin import ni export | 14/14 OK |
| (a) `data-props` parsea y `$preview` = `w`×`h` del fragmento (s13-m 412×4060, ya actualizado) | 14/14 OK |
| (a) `html.parser`: sin cerrar · cierres cruzados · ids repetidos · `label for` huérfanos | 0 · 0 · 0 · 0 |
| (a) Nombres `[A-Za-z0-9.-]` · stems únicos en el lienzo (90 artboards) · títulos repetidos | OK · 0 · 0 |
| (a) `format.md`: control dentro de `<a>` · interactivo dentro de `<button>` · `sticky`, `fixed` o `100vh` en escritorio | 0 · 0 · 0 |
| (b) Hex fuera de `data-props` (sin contar `href="#…"`) · `rgb()`, `hsl()`, `oklch()` · `fill`/`stroke` de SVG distintos de `none`/`currentColor` | 0 · 0 · 0 |
| (b) `var(--x)`: todos existen en `tokens.css`. Todos están en `tokens.json` por nombre, salvo `font-sans`, que en `tokens.json` es `type.families.sans` y `tokens.css` emite como `--font-sans` (lo usa todo el lienzo) · `var()` con respaldo | OK · 0 |
| (c) Contraste de texto por cascada, en claro y oscuro | 1,629 nodos por tema, 24 pares distintos: **0 bajo 4.5:1 y 0 con opacidad** (R2-A1 cerrado). Fuera de la cascada quedan **49 marcadores de posición visibles sin `::placeholder`** (R3-A1) |
| (c) Contraste de UI 3:1 | bordes de campo, interruptor (perilla y pista), anillo de foco, borde de edición y barras de avance: 0 bajo 3:1 (tabla abajo) |
| (d) P3-12: los 21 patrones de la Verifica, con y sin coma · en centavos · credenciales | 0 · 0 · 0 |
| (d) Fechas que coinciden con el calendario real (comando de P3-12) | En ajustes-b solo aparecen `2026-10-30` (s09-e, «Vie 30 oct 2026», próximo pago, sin monto) y `2027-02-28` (s08-e, ejemplo de U4 «28-feb-2027»). OK |
| (d) Dinero visible contra `datos-sinteticos.json` | 65 apariciones y 30 valores distintos: 26 del JSON y 4 derivados con la aritmética comprobada (80 % y 85 % de $27,586.00; 40 % y 60 % de $2,918.40) |
| (e) Fidelidad a U1/U4 | tabla abajo |
| (f) `button` con nombre y `type` · `a` con nombre · campos con `label` o `aria-label` · `svg` con `aria-hidden` o nombre · `role="switch"` con `aria-checked` | 0 fallas |
| (f) Objetivos táctiles en móvil (alto por cascada: `height`/`min-height` con tokens o relleno + interlineado + borde; casillas por su `label`) | 143 objetivos, **0 bajo 44 px**; 0 enlaces en línea; barra inferior 56 px |
| (f) Emoji · barra de estado o teclado falsos · `innerHTML`/`iframe`/red/`localStorage` · `type="date"` · fecha ISO visible | 0 · 0 · 0 · 0 · 0 |
| (f) Dinero `$1,234.56` · «-$» · montos sin `tabular-nums` | 0 mal formados · 0 · **0** (R2-M1 cerrado; antes 7) |
| (f) Inglés visible | 0 palabras en inglés. Quedan los nombres de tecla (Ctrl, Alt, Shift, Enter, Esc), «1.0.0-demo», «Ollama» y «API de Claude» (B-19) |
| (g) Móvil: raíz `.mr{width:412px;min-height:915px}` con `*{box-sizing:border-box}` · Escritorio: `.pg` flex-wrap sin ancho en px, contenido `flex:999 1 560px;min-width:0`, 10 tablas dentro de `overflow-x:auto`, `"expand":"fill"` en el fragmento | 7/7 · 7/7 |
| (h) Enlaces | 242 `href`: 212 a archivos de la lista del lienzo del mismo ancho y 30 a `#`. 0 rotos y 0 de otro ancho |
| Fragmento | JSON válido; `boards` = `order` = los 14 archivos; `page:"ajustes"`; `is_interactive:false`; `w` 412/1440; móvil `h` ≥ 915; títulos `P-26.xx · … · móvil|escritorio`; sin `x`/`y`. Las alturas no se volvieron a medir porque este revisor no usa navegador. El +164 px de s13-m cuadra con dos campos apilados (2 × 68 + 16 + 12) |
| Diferencias contra `docs/design/salida/iter-2/ajustes-b/` | Solo cambiaron s09-e, s12-m, s12-e, s13-m y s13-e (los `span.m` de R2-M1, la opacidad quitada y los dos campos de Ximena), más la altura de s13-m en el fragmento. No hay regresiones |

### Recálculo independiente de contraste (P3-07: 10 pares de texto y 7 de UI, hex de `tokens.css`)

| # | Par (texto / fondo) | Claro | Oscuro | Umbral |
|---|---|---|---|---|
| 1 | `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 ✓ |
| 2 | `ui-text-secondary` / `ui-surface-card` (incluye la nota «Con “mantener”…» de s13, ya sin opacidad) | 5.04 | 7.97 | 4.5 ✓ |
| 3 | `ui-text-default` / `ui-surface-field` | 11.53 | 14.07 | 4.5 ✓ |
| 4 | `ui-text-secondary` / `ui-surface-selected` (vistas previas) | 5.02 | 8.08 | 4.5 ✓ |
| 5 | `ui-role-success` / `ui-role-success-surface` | 4.81 | 6.19 | 4.5 ✓ |
| 6 | `ui-role-warning` / `ui-role-warning-surface` | 5.12 | 6.09 | 4.5 ✓ |
| 7 | `ui-text-on-primary` / `ui-action-primary-bg` | 5.77 | 10.47 | 4.5 ✓ |
| 8 | `ui-text-disabled` / `ui-surface-disabled` (campos nuevos de Ximena, C-042) | 5.04 | 6.09 | 4.5 ✓ (además exento por estar deshabilitado) |
| 9 | `ui-role-error` / `ui-role-error-surface` | 5.37 | 6.05 | 4.5 ✓ |
| 10 | `ui-text-link` / `ui-surface-card` | 5.16 | 7.94 | 4.5 ✓ |
| ✗ | Marcador de posición sin regla, color del navegador / `ui-surface-field` (#fafafb · #222325): Chrome `#757575` | **4.42** | **3.41** (sin `color-scheme`) | 4.5 ✗ |
| ✗ | ídem, Firefox (texto al 54 %) · Safari (`darkgray` #a9a9a9) | **3.05** · **2.25** | 5.11 · 6.69 | 4.5 ✗ |
| — | Lo que pide DESIGN C-035: `ui-text-secondary` / `ui-surface-field` | 5.40 | 10.43 | 4.5 ✓ |
| UI | `ui-border-control` / `ui-surface-card` · `-field` · `-page` · `-raised` · `-selected` | 3.47 · 3.71 · 3.71 · 3.71 · 3.46 | 6.09 · 7.96 · 7.96 · 6.09 · 6.17 | 3 ✓ |
| UI | perilla apagada `ui-surface-page` / `ui-border-control` (C-045) · perilla encendida (C-043) | 3.71 · 5.77 | 7.96 · 10.47 | 3 ✓ |
| UI | `ui-action-primary-bg` / `ui-surface-card` (pista encendida) | 5.16 | 7.94 | 3 ✓ |
| UI | `ui-border-focus` / `ui-surface-card` · `-page` · `-selected` | 3.56 · 3.82 · 3.55 | 6.03 · 7.89 · 6.11 | 3 ✓ |
| UI | `ui-border-selected` / `ui-surface-field` (campo en edición) | 5.53 | 10.38 | 3 ✓ |
| UI | `ui-viz-meter`/`ui-role-savings` sobre `ui-viz-track` (avance de metas) | 4.51 · 4.34 | 5.50 · 4.43 | 3 ✓ |

### Muestra de 10 valores contra `datos-sinteticos.json` (P3-12)
`$7,200.00` (`categorias[0].plan_fmt`) · `$14,260.00` (`ingresos.esperado.q15_cents`) · `$27,586.00` (`presupuesto.resumen.plan_gasto_fmt`) · `$20,525.78` (`presupuesto.resumen.real_fmt`) · `$399.00` (`presupuesto.gastos_fijos_recurrentes[1]`) · `$21,480.00` (`ahorro.fondo.ahorrado_fmt`) · `$2,783.34` (`ahorro.metas[0].aportacion_necesaria_fmt`) · `$18,000.00` (`ahorro.metas[1].monto_meta_fmt`) · `$2,918.40` y `$1,459.20` (`creditos.kueski.compras[1]`, capital y mitades; no son cifras canon) · respaldos 21-oct-2026 03:00 (2.3 MB) a 17-oct-2026 (2.7 MB) (`ajustes.respaldos`) · 84.6 kg, IMC 26.7 y 2,364 kcal; 62.8 kg, 23.6, 1,816 y 1900 kcal; ajuste de Ximena 0 y de la casa 0/0 (`salud.perfiles`, `salud.casa`). Todos coinciden.

## Estado de los hallazgos de la vuelta 2

| Id | Sev. r2 | Estado r3 | Evidencia |
|---|---|---|---|
| R2-A1 nota atenuada con `opacity:.85` (s13) | alta | **Cerrado** | `<p class="mu" id="mxximena">` sin `style`; 5.04 / 7.97. Ningún texto con opacidad en los 14 (la regla `.dis{opacity:.6}` de s08, s10 y s11 no la usa ningún elemento). Los dos ajustes de Ximena están dibujados (`abximena`, `asximena`) con `readonly`, `aria-disabled="true"`, `aria-describedby` y `.in[aria-disabled="true"]` (`ui-surface-disabled`, `ui-text-disabled`). Cierra también B-17 |
| R2-M1 montos en prosa sin `tabular-nums` | media | **Cerrado** | s09-e: los 3 `$27,586.00` en `span.m`. s12-m y s12-e: los 2 `$2,918.40` en `span.m`. El script cuenta 65 montos de 65 con `tabular-nums` |
| H-11 «Reponer» con nivel «Sin aviso» frente a K7 | media | **Abierto (3.ª vuelta)** | s14-e: «Huevo blanco … Sin aviso … 6 pza Reponer 2 días», «Leche entera … Sin aviso … 1,500 ml Reponer 3 días», «Tortilla de maíz … Sin aviso … 700 g Reponer 3 días», junto a «Regla: reponer si la existencia es igual o menor que el nivel». s14-m igual. `inventario.articulos` sigue con `umbral_reposicion: null` y `estado: "reponer"` en los tres; `resumen` 1/5/8. No hay archivo en `docs/estado/bloqueos/` |
| B-07 sin `::placeholder` | baja | **Sube a alta (R3-A1)** | se midió por primera vez: abajo |
| B-17 ajustes de Ximena no dibujados | baja | **Cerrado** | ver R2-A1 |
| B-01, B-03 a B-06, B-08 a B-16, B-18, B-19 | baja | **Abiertos** | fuera del alcance de la vuelta 3; se comprobó que siguen igual |

## Tabla artboard × regla

Columnas: F forma y fragmento · T tokens · C contraste (U5-04) · D cifras sintéticas (P3-12) · U fidelidad a U1/U4 · A accesibilidad y táctil · I idioma y formatos (U5-07, P3-08) · N navegación (barra lateral o inferior, enlaces) · E estados (U1 §2.2).

| Artboard | F | T | C | D | U | A | I | N | E | Fallas |
|---|---|---|---|---|---|---|---|---|---|---|
| s08-reportes-m | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✓ | ✓ | — | B-05, B-06 |
| s08-reportes-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | B-11, B-12, B-13, B-18 |
| s09-notificaciones-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05 |
| s09-notificaciones-e | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | R3-A1 (2), B-04, B-09, B-11, B-12 |
| s10-apariencia-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05 |
| s10-apariencia-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | n/a | B-05, B-11, B-12 |
| s11-datos-respaldos-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | — | B-05, B-19 |
| s11-datos-respaldos-e | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ~ | ✓ | ✓ | B-11, B-12, B-14, B-19 |
| s12-personas-reparto-m | ✓ | ✓ | ✗ | ✓ | ~ | ✓ | ✓ | ✓ | — | R3-A1 (3), B-05, B-08, B-10 |
| s12-personas-reparto-e | ✓ | ✓ | ✗ | ✓ | ~ | ✓ | ✓ | ✓ | ✓ | R3-A1 (9), B-10, B-11, B-12 |
| s13-salud-metas-m | ✓ | ✓ | ✗ | ~ | ~ | ✓ | ✓ | ✓ | — | R3-A1 (7), B-03, B-05, B-15, B-16, B-20 |
| s13-salud-metas-e | ✓ | ✓ | ✗ | ~ | ~ | ✓ | ✓ | ✓ | ✓ | R3-A1 (7), B-03, B-11, B-12, B-15, B-16, B-20 |
| s14-inventario-m | ✓ | ✓ | ✗ | ✓ | ✗ | ✓ | ✓ | ✓ | — | R3-A1 (7), H-11, B-05 |
| s14-inventario-e | ✓ | ✓ | ✗ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | R3-A1 (14), H-11, B-11, B-12 |

«—» en E: el móvil no lleva fila de estados (la tarea la pide en escritorio). «n/a»: P-26.10 no tiene vacío, carga ni error (U4). B-01, B-09 y B-21 son transversales. Entre paréntesis, los marcadores de posición visibles de R3-A1.

## Fidelidad: pantalla → elementos de U1/U4

| Sección | Presentes | Faltantes o distintos |
|---|---|---|
| P-26.08 Reportes | Tabla Reporte · Activo · Día · Hora · Destinatarios · Próximo envío · Último resultado; «Enviar ahora a…»; «Valor inicial»; regla de los días 29 a 31; Bajo demanda; estados. En móvil, tres tarjetas y Bajo demanda a todo el ancho (U4) | En móvil, «Reintentar el envío» en lugar de «Enviar ahora a…» (B-06). Día y hora duplicados en escritorio (B-18) |
| P-26.09 Notificaciones | 4 alertas con activa, destinatarios, anticipación o umbral y hora, cada una con «Valor inicial»; chats detectados con Confirmar como y Descartar; chats confirmados; aviso de prueba; vistas previas (80 → 85 %, $22,068.80 → $23,448.10); estados. En móvil, tarjetas por alerta y chat detectado arriba (U4) | 4 campos «Hora» para una sola clave (B-04); «•••» sin los 4 caracteres (B-09) |
| P-26.10 Apariencia | Tema segmentado, atajos de una tecla con vista previa, valores fijos, muestras clara y oscura | «Enter guarda · Esc descarta» en un control de guardado inmediato (B-05) |
| P-26.11 Datos y respaldos | Último, próximo y restauración probada; hora del respaldo diario (provisional); tabla de 5 respaldos; restaurar con RESTAURAR; importar con «Elegir archivo» y «Lo que se creará»; semilla; IA (Ollama y API de Claude, usos y llave); servicio; 3 parámetros; estados | «Restaurar respaldo» deshabilitado que no se ve deshabilitado (B-14) |
| P-26.12 Personas y reparto | Personas con papel, cuenta, paga directo, qué hace, tope y notas; 3 repartos con resto y «Reparto por defecto»; vista previa con $2,918.40 (tabular); diálogo de alcance. En móvil, tarjetas con barra de partes etiquetada (U4) | Sin control visible para «Parte de cada persona» (B-10) |
| P-26.13 Salud y metas | Fichas completas; ajustes de Ximena dibujados y atenuados con su explicación (U4: «los dos ajustes se ven atenuados»); ajuste de la casa con «Valor inicial» y «Lo usan…»; metas con inicio, notas, activa y vista previa; Tabla · Tablero · Calendario · Galería; «Nueva meta»; estados | Los ajustes de Ximena muestran «0» como valor propio, pero la misma pantalla dice que ella deja vacío su ajuste y usa el de la casa (B-20). «Mi ajuste para bajar» muestra 400 y la vista previa dice 400 → 500 (B-03). Avance entero y estados «Atrasada»/«En ruta» frente a Q2-17 C y Q2-18 A (B-16). Sin el botón «Usar el de la casa» (U4-13) en el ajuste de Axel (incluido en B-20) |
| P-26.14 Inventario | 14 artículos con modo, unidad, nivel, ritmo, caducidad, existencia, estado, días enteros y notas; resumen 1/5/8; Tabla · Tablero por modo; «Agregar producto»; vista previa; regla K7; estados | «Reponer» con «Sin aviso» frente a la regla K7 de la misma pantalla (H-11) |

## Líneas de evidencia

- **P3-04**: 0 acciones de navegador y 0 llamadas a Artifact de este revisor; 0 credenciales. `referencia-tipo/format.md` se leyó como dato que fija la forma (lo escribió el publicador del tipo). En los artboards y el fragmento no hay texto con aspecto de instrucción.
- **P3-06 (cobertura de ajustes-b)**: 7 secciones × 2 anchos = 14/14; 0 faltantes. La copia de la vuelta 2 está en `docs/design/salida/iter-2/ajustes-b/` (14 + fragmento) y se usó para el diff.
- **P3-07**: Tabla artboard × regla arriba. Se recalcularon 10 pares de texto y 7 de UI. Los pares de color por cascada tienen 0 bajo umbral. Los marcadores de posición sin regla quedan bajo 4.5:1 en claro en los tres motores (R3-A1).
- **P3-08** (8 puntos): (1) Resumen ≤ 6 tarjetas: no aplica · (2) 14 secciones de Ajustes: ✓ (la columna interna «Secciones de Ajustes» enlaza el índice y las 14, con `aria-current`) · (3) barra lateral colapsable: ✓ («Ocultar barra lateral») · (4) español de México: ✓, con B-19 · (5) dinero `$1,234.56` tabular: ✓ (65/65; 0 montos de 1,000 o más sin separador) · (6) fechas `10-oct-2026`: ✓ (0 ISO, 0 `type="date"`) · (7) paleta C2 índigo 280°: ✓, solo tokens · (8) Notion como única referencia: ✓.
- **P3-12**: 0 cifras canon (21 patrones, con y sin coma, y en centavos); 0 credenciales. Se revisaron a mano las fechas coincidentes y ninguna va con monto. La muestra de 10 valores coincide con el JSON. **Cumple**, con la observación B-15 (las fechas de nacimiento no están en el JSON).
- **U5-04**: 0 pares de color bajo umbral (1,629 nodos por tema). Los 49 marcadores de posición visibles toman el gris del navegador y quedan bajo 4.5:1 (R3-A1). **Falla.**
- **U5-06**: Cumplen los tokens, la barra lateral colapsable, los objetivos de 44 px, el foco visible y las cifras tabulares (✓, R2-M1 cerrado). Los iconos de 16 px en insignias no siguen los 20/24 px (B-21), y la escala tipográfica y el espaciado llevan literales (B-01, B-21).
- **U5-07**: 0 palabras en inglés, 0 fechas ISO, dinero bien formado; «1.0.0-demo», «Ollama» y «Claude» no están en la lista permitida (B-19).

## Hallazgos

### Alta
- **R3-A1 (era B-07) · U5-04, P3-07 (texto ≥ 4.5:1), DESIGN §10.5 C-035.** Ninguno de los 14 artboards tiene regla `::placeholder`. Hay 49 marcadores visibles en campos vacíos: «Sin notas» (37), «Sin tope» (6), «Sin restricciones» (4) y «Sin etiqueta» (2), en s09-e, s12-m, s12-e, s13-m, s13-e, s14-m y s14-e. Todos toman el color del navegador. Sobre `ui-surface-field` (#fafafb en claro) quedan así: Chrome `#757575` = **4.42:1**, Firefox (texto al 54 %) = **3.05:1** y Safari (`darkgray`) = **2.25:1**. En oscuro, Chrome sin `color-scheme` da **3.41:1** sobre #222325. WCAG 1.4.3 incluye el texto de marcador de posición. DESIGN fija el par C-035 (`ui-text-secondary` sobre `ui-surface-field`: 5.40 / 10.43), y comidas y gastos ya lo aplican. Las vueltas 1 y 2 la dejaron en baja sin medirla. **Corrección:** agregar `::placeholder{color:var(--ui-text-secondary);opacity:1}` al `helmet` de los 14 artboards (también a los que hoy no tienen marcadores, para que no se repita).

### Media
- **H-11 (abierto, 3.ª vuelta) · U4 P-26.14 y K7.** En s14-m y s14-e, huevo, leche y tortilla dicen «Reponer» con nivel «Sin aviso», y la misma pantalla escribe «Regla: reponer si la existencia es igual o menor que el nivel; agotado si es 0». Según U4, nivel vacío significa sin aviso, así que la regla no puede dar «Reponer». El origen es `inventario.articulos` de `datos-sinteticos.json` (`umbral_reposicion: null`, `estado: "reponer"`), que no cambió. Los días ya salen enteros. **Corrección:** el orquestador encarga a correccion-datos una de dos salidas. (a) Dar a esos tres artículos un `umbral_reposicion` mayor o igual que su existencia (por ejemplo 6 pza, 1,500 ml y 700 g o más). (b) Pasarlos a «suficiente» y ajustar `inventario.resumen`. Después el autor repinta s14-m/e y `compras-p16-inventario-m/e`. Si no se hace en esta iteración, hay que registrar el bloqueo en `docs/estado/bloqueos/` y marcar el estado como «no adoptar» en la referencia visual (P3-08).

### Baja
- **B-01 · DESIGN §7.2, §7.3 y §7.10.** El texto base y los campos móviles miden 14 px (§7.2 y §7.3 piden 16 px). El título de escritorio mide 24 px. `.mx{max-width:1180px}` y los tamaños de letra son literales, no tokens.
- **B-03 · U4 §0.2.** En s13, «Mi ajuste para bajar» de Axel muestra 400 con borde de edición y la vista previa dice 400 → 500. El campo debe mostrar 500.
- **B-04.** s09-e tiene 4 campos «Hora», uno por alerta, para la sola clave `alerts.send_time`.
- **B-05 · U4 §0.7.** «Enter guarda · Esc descarta» aparece en los 7 móviles y en controles de guardado inmediato (s10).
- **B-06 · U4 P-26.08.** s08-m trae «Reintentar el envío», que U4 no tiene, y no trae «Enviar ahora a…».
- **B-08.** En la barra de reparto de s12-m, los segmentos de Karen y Axel van pegados (1.51:1 entre sí en oscuro). El dato también va en el texto.
- **B-09.** La columna interna va sin los iconos de §5.4. La terminación del chat es «•••» sin los 4 caracteres (el JSON no los trae). s14-e tiene 14 rótulos «Modo» iguales.
- **B-10 · U4 P-26.12.** «Parte de cada persona» no tiene control de edición visible.
- **B-11.** La barra lateral dice «Axel · Co-titular» en los 7 escritorios y s12-e dice «Titular». Viene del shell y lo unifica el orquestador.
- **B-12 · DESIGN §7.9.** La barra lateral va sin `chevron-right` y con `aria-label="Principal"` en lugar de «Secciones».
- **B-13.** En s08-e, los campos de fecha de Karen llevan `inputmode="numeric"`, pero el formato tiene letras.
- **B-14 · DESIGN §7.2.** En s11-e, «Restaurar respaldo» tiene `aria-disabled="true"` y se pinta con el color crítico completo.
- **B-15 · P3-12.** Las fechas de nacimiento 14-mar-1990 y 22-jun-1993 no están en `datos-sinteticos.json`.
- **B-16 · U4 P-26.13 (Q2-17 C, Q2-18 A) frente a DESIGN §3.3 y U1 §2.5.** El avance sale entero (42 %, 78 %) y los estados dicen «Atrasada» y «En ruta». Decide el orquestador.
- **B-18.** En s08-e, día y hora se editan en la tabla y en «Parámetros operativos»: 20:00 en un sitio y 19:00 en el otro.
- **B-19 · U5-07.** «1.0.0-demo», «Modelo local (Ollama)» y «API de Claude» no están en la lista permitida.
- **B-20 (nueva) · U4 P-26.13 (K36, U4-13).** Los dos campos nuevos de Ximena muestran «0» como valor propio. La misma pantalla dice «Lo usan las personas que dejan vacío su propio ajuste: hoy Ximena (bajar y subir) y Axel (subir)», y en U4 vacío significa «usa el de la casa». Para cuadrar, esos campos deben ir vacíos con «El de la casa: 0 kcal», como el «Mi ajuste para subir» de Axel. Falta además el botón «Usar el de la casa» (el «restablecer» de U4-13) en el ajuste de Axel, que tiene valor propio (400).
- **B-21 (nueva) · DESIGN §5.2, U5-06 (iconos de 20/24 px).** `.tg .ic{width:16px;height:16px}` pinta 63 iconos de 16 px en las insignias de 13 artboards. DESIGN §5.2 dice «Solo dos tamaños: un icono de 16 px a 1.5 px pierde detalle», y con `stroke-width` de 1.8 el trazo dibujado queda en 1.2 px. También hay literales fuera de la escala de 4 pt: `gap:2px`, `gap:6px` y `padding:2px` en fichas y la barra inferior, y `border-radius:12px` en la pista del interruptor. **Corrección:** iconos de insignia con `var(--size-icon-sm)` y espaciados con tokens.

## Correcciones sugeridas (orden)
1. R3-A1: una regla `::placeholder{color:var(--ui-text-secondary);opacity:1}` en el `helmet` de los 14 artboards (es CSS y no cambia las alturas).
2. H-11: correccion-datos ajusta `inventario.articulos` y `resumen`, y después se repintan s14-m/e y compras-p16. Si no se puede en esta iteración, registrar el bloqueo y marcar «no adoptar».
3. Bajas baratas de CSS y texto: B-03, B-05, B-12, B-13, B-14, B-20 y B-21.
4. Para el orquestador o correccion-datos: B-11, B-15, B-16 y B-19.
