# P3-iter-1-salud-reportes · Revisión de la vuelta 1 · área «salud-reportes»

Revisor: `revisor-salud-reportes-r1` (independiente; no corrige) · 12-oct-2026 · Pista 3, fase 2 (lienzo), vuelta 1 de 3 · Criterio P3-07.

Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/salud-reportes/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/salud-reportes.json`, contra P3 1.3 (P3-06, P3-08, P3-12), U5 (U5-04, U5-06, U5-07, en lo que se ve en el marcado), DESIGN.md (§5.2, §7.5, §7.7, §7.9, §7.10, §7.14, §8.6, §8.7, §10.5), U1 (§1.2, §1.3, §2.1, §2.2, §2.4, §2.5, §3.23 a §3.27, F-5, F-6, §6.4) y `referencia-tipo/format.md` y `craft.md`. U3 y U4 no aplican a esta área (no hay Inicio ni Ajustes).

## Veredicto

**No aprobado.** 0 altos · **5 medios** · 14 bajos. Forma, tokens, contraste, cifras sintéticas, enlaces y anchos están bien en los 8; los medios son de fidelidad (barra lateral sin colapsar, barra de la vista de P-23, filtro de avisos equivocado, restricciones alimentarias) y uno de formato de fecha («2026-W43»).

## Cobertura (P3-06)

| Pantalla U1 | Móvil 412 | Escritorio 1440 | Estados (vacío · carga · error) en escritorio |
|---|---|---|---|
| P-21 Salud | `salud-p21-salud-m` | `salud-p21-salud-e` | sí · sí · sí |
| P-23 Centro de reportes | `reportes-p23-centro-m` | `reportes-p23-centro-e` | sí · sí · sí |
| P-24 Reporte mensual | `reportes-p24-mensual-m` | `reportes-p24-mensual-e` | sí · sí · sí |
| P-25 Estado de cuenta de Karen | `reportes-p25-karen-m` | `reportes-p25-karen-e` | sí · sí · sí |

8 de 8 artboards; faltantes: 0 (P-22 es del área «shell»).

## Verificación propia (scripts `python3 -I` en `datos/tmp/rev-sr-r1/`, borrados al final)

| Punto | Resultado |
|---|---|
| (a) Forma | 8/8 con `<!doctype html>`, `lang="es-MX"`, `<meta charset="utf-8">`, `<title>` propio, `./support.js`, `../ds/finanzas-hogar/tokens.css`, `<x-dc>`, `<helmet><style>`, `data-dc-script` con `class Component extends DCLogic` y `renderVals()`. `data-props` parsea en los 8; `$preview` = (w, h) del fragmento en los 8. `html.parser`: 0 etiquetas sin cerrar, 0 cierres inesperados. 0 IDs duplicados. 88 stems en todo el lienzo, 0 duplicados; nombres válidos. |
| (b) Tokens | 0 hex literales fuera de `data-props` (los `href="#…"` excluidos); 0 `rgb()/hsl()`. `var(--x)`: 45 a 55 por archivo; todos existen en `tokens.css`; todos en `tokens.json` salvo `--font-sans` (lo define `tokens.css` del sistema y lo usan los 88 artboards; ver B14). Espaciado y tipografía con tokens; único literal: `border-bottom:2px` de las pestañas (B09). |
| (c) Contraste | 55 pares recalculados (tabla abajo), claro y oscuro: **0 bajo umbral**. |
| (d) P3-12 | Grep de las 21 cifras canon (con y sin coma) en los 8 y en el fragmento: **0**. Grep de credenciales: 0. Grep P3-08 `\$[0-9]{4,}\.[0-9]{2}`: 0. Los 57 importes distintos están todos en `datos-sinteticos.json` (campos `*_fmt`, `*_cents` o textos); las 18 fechas distintas, todas en el JSON; días de la semana correctos (Mié 21 oct 2026, Dom 25 oct 2026, Dom 1 nov 2026). |
| (e) Fidelidad | Tabla abajo. |
| (f) Accesibilidad | 0 controles sin nombre accesible; 0 campos sin `<label>`; todos los `<svg>` con `aria-hidden` o `role="img"` + `aria-label`; botones, enlaces, `input`, `select` y `label` reales; 0 emoji; sin barra de estado ni teclado falsos; objetivos táctiles en móvil: 101 controles, todos ≥ 44 px (`--size-control-coarse` y `--size-touch-min` = 44 px, barra inferior 56 px). Dinero: 57 importes con formato `$1,234.56`, 0 sin cifras tabulares, signo U+2212 (`−$98.60`), 0 guion ASCII. |
| (g) Anchos | Raíz móvil `width:412px` (box-sizing por `*`) en los 4; escritorio fluido (sin ancho en px; barra `flex:1 1 240px` + contenido `flex:999 1 560px; min-width:0` con `flex-wrap`); `"expand":"fill"` solo en los 4 de escritorio. |
| (h) Enlaces | 120 `href`: 0 fuera de la lista del lienzo, 0 con ancho cruzado (`-m` a `-e`), 0 anclas internas sin destino; 24 a `#` (2 de ellos ya tienen destino, B06). |

### Contraste recalculado (fórmula WCAG 2.x con los hex de `tokens.css`; extracto de 20 de los 55 pares)

| Primer plano / fondo | Tipo | Claro | Oscuro | En §10.5 |
|---|---|---|---|---|
| ui-text-default / ui-surface-page | texto | 11.53 | 14.07 | sí |
| ui-text-secondary / ui-surface-card | texto | 5.04 | 7.97 | sí |
| ui-text-secondary / ui-surface-hover | texto | 4.51 | 5.50 | sí |
| ui-text-link / ui-surface-card | texto | 5.16 | 7.94 | sí |
| ui-text-on-primary / ui-action-primary-bg | texto | 5.77 | 10.47 | sí |
| ui-text-on-primary / ui-action-primary-bg-hover | texto | 8.49 | 12.66 | sí |
| ui-text-default / ui-surface-selected | texto | 10.73 | 10.90 | sí |
| ui-icon-selected / ui-surface-selected | UI | 5.15 | 8.04 | sí |
| ui-role-info / ui-role-info-surface | texto | 4.96 | 6.12 | par de rol |
| ui-role-success / ui-role-success-surface | texto | 4.81 | 6.19 | par de rol |
| ui-role-error / ui-role-error-surface | texto | 5.37 | 6.05 | par de rol |
| ui-role-warning / ui-role-warning-surface | texto | 5.12 | 6.09 | par de rol |
| ui-role-warning / ui-surface-card | texto | 5.13 | 6.04 | calculado |
| ui-text-default / ui-role-error-surface | texto | 10.72 | 11.20 | calculado |
| ui-border-control / ui-surface-card | UI | 3.47 | 6.09 | sí |
| ui-border-control / ui-role-info-surface | UI | 3.49 | 6.00 | calculado |
| ui-border-focus / ui-surface-card | UI | 3.56 | 6.03 | sí |
| ui-person-axel / ui-surface-card | UI | 11.93 | 3.10 | calculado |
| ui-person-ximena / ui-surface-card | UI | 4.85 | 6.93 | calculado |
| ui-role-warning / ui-viz-track | UI | 4.59 | 4.16 | calculado |

El par más justo es `ui-person-axel` sobre tarjeta en oscuro (3.10:1, UI ≥ 3); `ui-text-disabled` sobre `ui-surface-disabled` (5.04 / 6.09) pasa aunque estaría exento. Esqueletos y retícula: decorativos (§10.5).

## P3-08 · Lista de 8 puntos

| # | Decisión | Resultado |
|---|---|---|
| 1 | Resumen ≤ 6 tarjetas | No aplica (área sin Inicio) |
| 2 | 14 secciones de Ajustes | No aplica |
| 3 | Barra lateral colapsable | **Falla** (M01: sin botón «‹‹» en los 4 de escritorio) |
| 4 | Español de México | Cumple, con B01 y B03 |
| 5 | Dinero `$1,234.56` tabular | Cumple (57 importes; grep de miles sin coma = 0) |
| 6 | Fechas `10-oct-2026` | **Falla** parcial (M05: «2026-W43» visible); el resto cumple |
| 7 | Paleta de C2 (índigo 280°) | Cumple (solo `var(--…)`) |
| 8 | Notion como única referencia | Cumple |

## P3-12

Línea `P3-12`: 0 cifras canon en los 8 artboards y el fragmento; 0 credenciales; nombres Axel, Ximena y Karen sí (K45.1). Fechas coincidentes con el calendario real (15, de `2026-10-30` a `2027-06-15`): los artboards usan 30-oct, 15-nov, 30-nov y 15-dic de 2026 en P-24 y P-25 con montos $381.84, $453.49, $803.27 y $803.26, que no son del reporte ni de Kueski. Muestra de 10 valores con su origen en `datos-sinteticos.json`: $31,130.00 (`reporte_mensual.resumen.ingreso_fmt`), $20,525.78 (`real_fmt`), $1,750.00 (`ahorro_aportado_fmt`), $42,744.75 (`creditos.deuda_viva_cierre_fmt`), $1,107.07 (`intereses_fmt`), $1,551.70 (`canasta.por_tienda[bodega-aurrera].total_fmt`), $3,802.33 (`estado_cuenta_karen.saldo_pendiente_fmt`), 84.6 kg del 19-oct-2026 (`salud.perfiles[axel]`), 1,816 kcal sugerida de Ximena (`meta_kcal_sugerida`), 21:30 del 21-oct-2026 (`reportes.proximos_envios[0]`).

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presentes | Faltantes o distintos |
|---|---|---|
| P-21 | Dos columnas (escritorio) y pestañas Axel/Ximena (móvil); Resumen·Detalle·Tendencias; Registrar peso (primaria, a P-22); Editar perfil; altura, edad, sexo, actividad, objetivo, peso y fecha, IMC con clasificación; meta confirmada y sugerida con su fórmula; «Confirmar 1,816 kcal» (K23); ajuste para bajar y subir con «El de la casa: 0 kcal» y «Restablecer» (N-327); kcal de hoy con meta y diferencia; gráfica de peso (círculo continuo Axel, cuadrado 8·4 Ximena, §8.6) con «Ver como tabla»; historial con IMC y editar (N-242); estados | **Restricciones alimentarias** (N-230) — M04; Editar perfil a `#` (B06) |
| P-23 | Generar reporte mensual con mes y opción de Telegram (N-248); próximos envíos con hora y destinatarios (N-246); reportes generados con tipo, periodo, generado, origen, estado, envío y Descargar (N-245); fila en error con «Reintentar envío» (N-250); avisos enviados con fecha, tipo, para, resultado y texto (N-247), plegados en móvil; estado de cuenta de Karen «Sin envío · se comparte a mano» (decisión A); estados | **Selector de vista** Tabla·Tablero·Calendario y filtro **Estado** de reportes (U1 §2.1, §2.5) — M03; opciones del filtro **Tipo de avisos** equivocadas — M02 |
| P-24 | Selector de mes con «parcial» (N-255); 7 secciones plegables (resumen, presupuesto por categoría, flujo por quincena, créditos, despensa y precios, comidas y salud, participación de Karen) a ancho de lectura con índice; Generar PDF (no hay PDF de oct) y Compartir deshabilitado; estados por sección con Reintentar | Mes por defecto «último cerrado» (B10, lo impone el JSON); «semana 2026-W43» (M05) |
| P-25 | Periodo (mes en curso, mes anterior, rango); compras particulares con comercio, fecha, capital, intereses, total y pago por quincena; su parte de cada despensa; pagos del periodo y pendientes posteriores; pagado en el periodo; saldo pendiente; nota fija de Kueski; Generar PDF; Compartir el último PDF; estados anteriores (N-262); movimientos plegados en móvil; estados | Vacío sin acción (B02) |

## Hallazgos

### Medios (bloquean la aprobación)

- **M01 · P3-08 / U5-06 / DESIGN §7.9 / U1 §1.2.7.** La barra lateral de escritorio no tiene el botón «‹‹» (`panel-left`, «Colapsar barra lateral»). Rutas: `salud-p21-salud-e`, `reportes-p23-centro-e`, `reportes-p24-mensual-e`, `reportes-p25-karen-e`. Las áreas comidas, gastos y presupuesto ya lo dibujan. Corrección: en el encabezado de la barra, `<button class="btn bt" type="button" aria-label="Colapsar barra lateral">` con icono `panel-left` de 20 px.
- **M02 · U1 N-247 y §2.5.b.** En «Avisos enviados» de P-23 el filtro «Tipo» ofrece «Mensual PDF» y «Estado de cuenta de Karen» (tipos de reporte, copiados del otro filtro). Rutas: `reportes-p23-centro-e`, `reportes-p23-centro-m`. Corrección: Todos, Resumen diario, Resumen semanal, Reporte mensual, Alerta de pago próximo, Alerta de presupuesto al 80 %, Alerta de presupuesto al 100 %, Alerta de precios viejos, Alerta de plan sin armar; periodo por defecto «Mes en curso» (§2.5.b).
- **M03 · U1 §2.1 («Tabla*, calendario (envíos)»), §2.5.a/b, §1.3; DESIGN §7.5.1.** P-23 no dibuja la barra de la vista: falta el selector de vista (Tabla* · Tablero por tipo · Calendario por fecha, con `aria-pressed`) en «Reportes generados» y «Avisos enviados», y el filtro «Estado» de los reportes (job: periodo, tipo, estado). En móvil los filtros van como dos `select` en línea en vez del botón «Filtrar» que abre una hoja y deja chips quitables. Rutas: `reportes-p23-centro-e`, `reportes-p23-centro-m`. Corrección: fila de vista sobre cada lista (escritorio) y botón «Filtrar» + chips (móvil).
- **M04 · U1 N-230.** El perfil de P-21 no muestra «Restricciones alimentarias» (vacías hoy, decisión D). Rutas: `salud-p21-salud-e`, `salud-p21-salud-m`. Corrección: fila `Restricciones alimentarias` · `Ninguna` en el `dl` de cada persona.
- **M05 · U5-07 («ninguna fecha ISO visible»), P3-08 (fechas).** P-24 muestra «Canasta de la lista de compras, semana 2026-W43 (19 al 25 oct).» Rutas: `reportes-p24-mensual-e`, `reportes-p24-mensual-m`. Corrección: «Canasta de la lista de compras de la semana del 19 al 25 oct 2026.»

### Bajos

- **B01 · U5-07 / datos.** «Falló: sin conexión con telegram» (el JSON dice «Telegram») en `reportes-p23-centro-e` y `-m`; «Manual (axel)» en `reportes-p23-centro-e` (el JSON y el móvil dicen «manual (Axel)»). Copiar el texto literal del JSON.
- **B02 · U1 §2.4, DESIGN §7.14.** Los vacíos de P-23, P-24 y P-25 no ofrecen la acción que los llena (p. ej. «Generar reporte mensual», «Cambiar mes», «Cambiar periodo») ni el icono de la entidad; P-24 cambia el literal «Sin datos en sep-2026» por «abr-2026» y su carga no es «por sección». Rutas: `reportes-p23-centro-e`, `reportes-p24-mensual-e`, `reportes-p25-karen-e`.
- **B03 · U5-07 (lista de palabras).** La insignia «Error: base de datos ocupada» (`reportes-p23-centro-e`, `-m`) y los rótulos «Error» de la tarjeta Estados harían saltar la búsqueda literal de `Error` de U5-07 (es español, pero el script no distingue). Sugerencia: «No se generó: base de datos ocupada» y «Falla».
- **B04 · DESIGN §7.10.** Los destinos de la barra inferior usan icono de 20 px; DESIGN pide `size.icon.md` (24 px). Rutas: los 4 `-m`.
- **B05 · U5-06 / DESIGN §7.9.** P-21 salta de `h1` a `h3`; las secciones de P-24 llevan el título en `summary` sin encabezado; los 3 `role="meter"` de P-24 no tienen `aria-label`; el chip de usuario es `<span aria-label>` sin rol ni menú (Cuenta, Apariencia, Cerrar sesión); los chevrons no son botones con `aria-expanded`.
- **B06 · enlaces.** «Editar perfil» (P-21) y «Cambiar en Ajustes › Reportes y frecuencias» (P-23) van a `#`, pero ya existen `ajustes/ajustes-p26-s13-salud-metas-{m,e}.dc.html` y `ajustes/ajustes-p26-s08-reportes-{m,e}.dc.html`.
- **B07 · format.md (eventos y estado).** `salud-p21-salud-m` va con `"is_interactive": true` y pestañas por radios CSS con `checked` y `value` estáticos; en el runtime de React pueden quedar como controles controlados y no cambiar. La plantilla pedía `false`; el patrón del tipo es `state` + `sc-if`. Corrección: pestañas con `onClick` y `state`, o `is_interactive: false`.
- **B08 · altos.** Los `h` son estimados (supuesto 1 del autor) y la raíz móvil lleva `min-height` = h con la barra inferior al fondo (`margin-top:auto`): riesgo de hueco grande sobre la barra o de recorte. Medir al publicar.
- **B09 · tokens.** `border-bottom:2px` literal en el indicador de pestañas; usar `var(--border-width-focus)`. Los 8.
- **B10 · U1 N-256, F-5.** P-24 abre en oct-2026 (parcial) y no en el último mes cerrado; en F-5 la fila «Mensual PDF · sep-2026» de P-23 lleva a un P-24 de octubre. Lo impone `reporte_mensual.mes` del JSON; anotarlo en el artboard («Ejemplo con el mes en curso»).
- **B11 · U1 §6.4 («bajo demanda»).** P-23 escritorio muestra los avisos abiertos y P-25 escritorio los movimientos con `open`; §6.4 los pide plegados. P-21 dibuja la serie de peso en Resumen, mientras §2.5.a dice «La serie va en Tendencias» (choca con §2.1 «gráfica de peso»; se acepta y se anota).
- **B12 · legibilidad.** El texto de ejes del gráfico de peso en móvil (`viewBox` 372 en ≈ 348 px) queda en ≈ 11.2 px, bajo `caption` (12 px). Usar `viewBox` de 348 o texto fuera del SVG.
- **B13 · U1 §1.3.** En móvil las acciones secundarias van como botones a todo el ancho («Estado de cuenta de Karen», «Ver reporte web», «Compartir…») y no en el «⋯» del encabezado. Rutas: los 3 `reportes-*-m`.
- **B14 · observación de tokens.** `var(--font-sans)` no está en `tokens.json` (lo define `tokens.css` del sistema y lo usa todo el lienzo). Sin cambio salvo que el orquestador lo decida.

## Tabla artboard × regla

Reglas: F forma (format.md), T tokens, C contraste U5-04, D cifras P3-12, Fi fidelidad U1, A accesibilidad U5-06, L idioma y formatos U5-07, R raíz y ancho, E enlaces, P8 P3-08.

| Artboard | Cumple | Falla (ID) | Correcciones para la vuelta 2 |
|---|---|---|---|
| `salud-p21-salud-m` | F, T, C, D, R, E, L | Fi (M04), A (B04, B05), B07, B08, B09, B12 | Restricciones; iconos 24 px; `h2` por persona; pestañas con `state` o `is_interactive:false`; `viewBox` |
| `salud-p21-salud-e` | F, T, C, D, R, E, L | P8 (M01), Fi (M04), A (B05), B06, B09 | Botón colapsar; restricciones; enlazar Editar perfil a P-26.13 |
| `reportes-p23-centro-m` | F, T, C, D, R, A | Fi (M02, M03), L (B01, B03), B04, B06, B09, B13 | Filtro de avisos; «Filtrar» + chips + selector de vista; textos del JSON |
| `reportes-p23-centro-e` | F, T, C, D, R, A | P8 (M01), Fi (M02, M03), L (B01, B03), B02, B06, B09, B11 | Botón colapsar; barra de la vista con Estado; filtro de avisos; vacío con acción |
| `reportes-p24-mensual-m` | F, T, C, D, R, E | L/P8 (M05), A (B04, B05), B09, B10, B13 | «semana del 19 al 25 oct 2026»; `h2` en secciones; `aria-label` en medidores |
| `reportes-p24-mensual-e` | F, T, C, D, R, E | P8 (M01), L/P8 (M05), A (B05), B02, B09, B10 | Botón colapsar; fecha de la canasta; vacío literal con acción |
| `reportes-p25-karen-m` | F, T, C, D, R, E, L, Fi | B04, B09, B13 | Iconos 24 px; secundarias a «⋯» |
| `reportes-p25-karen-e` | F, T, C, D, R, E, L, Fi | P8 (M01), B02, B09, B11 | Botón colapsar; vacío con «Cambiar periodo»; movimientos plegados |

## Fragmento

`salud-reportes.json` es JSON válido; `boards` y `order` listan exactamente los 8 archivos y en el mismo orden; `w`/`h` = `$preview`; títulos con la forma «P-XX · Nombre · móvil|escritorio»; `page` = «salud-reportes»; `expand: fill` solo en escritorio; `is_interactive` `false` salvo P-21 móvil (B07).

## Texto con aspecto de instrucción (P3-04)

Ninguno obedecido. `referencia-tipo/format.md` y `craft.md` traen indicaciones del publicador del tipo («Read this when…», «Give every artboard…», «No fake chrome») que se leyeron como datos de forma. Los artboards no contienen texto dirigido a un agente.

## Líneas de evidencia

- `P3-04`: el revisor no usó navegador ni la herramienta Artifact (solo Bash con lectura de archivos y escritura de sus dos reportes); 0 credenciales encontradas.
- `P3-06`: 8 de 8 artboards del área (4 pantallas × 2 anchos); faltantes 0.
- `P3-07`: este informe; 55 pares de contraste recalculados, 0 bajo umbral.
- `P3-08`: lista de 8 puntos arriba; fallan 3 (M01) y 6 (M05); grep de miles sin coma = 0.
- `P3-12`: 0 cifras canon, 0 credenciales; fechas coincidentes revisadas a mano sin montos reales.
- `U5-04`: 0 pares bajo umbral en claro y oscuro (estático, por tokens). `U5-06`: M01, B04, B05. `U5-07`: M05, B01, B03.

## Supuestos

1. El informe va en `docs/estado/wfp3/` (lo pide la tarea); P3-07 nombra `docs/estado/P3-iter-N-<area>.md`. El orquestador decide si lo copia.
2. Contraste medido estáticamente: pares por clase y estilo en línea resueltos con `tokens.css` (bloques claro y oscuro), no en un DOM renderizado (K45: sin navegador).
3. «Elemento secundario ausente» (media) incluye la barra lateral sin botón de colapsar: el lienzo no la contradice, la omite; el artboard de la barra lateral (`shell`) es el que la define.
4. Las pestañas Detalle y Tendencias de P-21 no son artboards de la lista: su ausencia no cuenta como faltante.
