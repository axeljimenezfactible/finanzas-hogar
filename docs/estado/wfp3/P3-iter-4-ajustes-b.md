# P3-iter-4 · ajustes-b · revisión de la vuelta 4

Fecha: 12-oct-2026 · Revisor: `revisor-ajustes-b-r4` (independiente; vuelta 4, ronda extra sobre las 3 previstas; no corrige) · Criterio: P3-07 de `docs/criterios/P3-claude-design.md` 1.3.
Alcance: los 14 artboards de `docs/design/maquetas/lienzo/project/ajustes/` de P-26.08 a P-26.14 (`-m` y `-e`) y el fragmento `docs/design/maquetas/lienzo/fragmentos/ajustes-b.json`, tal como los dejó `correccion-ajustes-b-r3`.
Contra: P3-08 y P3-12; U5-04, U5-06 y U5-07 (lo que se puede verificar en el markup); DESIGN.md (§3.3, §5.2, §7, §10.5 C-035); U1 §2.1, §2.2 y §2.3; U4 P-26.08 a P-26.14; `referencia-tipo/format.md`; `datos-sinteticos.json` (regenerado por `correccion-datos-r4`).
Ruta del informe: la que fijó la tarea (`docs/estado/wfp3/`). La copia a `docs/design/salida/iter-4/` (P3-06) la hace el orquestador.

**Veredicto: APROBADO** · 0 altas · 0 medias · 21 bajas (18 que siguen abiertas de antes y 3 nuevas).
Hallazgos de la vuelta 3: **R3-A1 cerrado** y **H-11 cerrado**. No hay regresiones. Contra la copia `docs/design/salida/iter-3/ajustes-b/` solo cambiaron dos cosas: la regla `::placeholder` en los 14 artboards y los tres niveles de s14-m y s14-e. El fragmento es idéntico.

## Verificación ejecutada (scripts propios, `python3 -I`, en `datos/tmp/rev-ajb-r4/`, borrados al final)

Escribí un motor de cascada propio. Aplica las reglas del `helmet` y el `style` en línea, con selectores de etiqueta, clase, id, atributo (`=`, `~=`), `>`, `+`, descendiente, `:first-child`, `:last-child` y `:not()`. Calcula especificidad y orden. No aplica `:hover` ni `:focus`. También resuelve `var()` contra los hex de `tokens.json` en claro y oscuro, hereda `color` y `font-variant-numeric`, compone el fondo efectivo con alfa y respeta el `data-theme` local y la opacidad acumulada.

| Comprobación | Resultado |
|---|---|
| (a) Forma: `<!doctype html>` y `<html lang="es-MX">`, `<meta charset="utf-8">`, `<title>` propio, `<script src="./support.js"></script>`, `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`, `</head><body><x-dc><helmet><style>`, `</style></helmet>`, `</x-dc>` seguido del `data-dc-script` con `class Component extends DCLogic { renderVals() {`, y cierre `</script></body></html>`. Cada marcador aparece una sola vez. Sin `import` ni `export` | 14/14 OK |
| (a) `data-props` parsea y `$preview` = `w`×`h` del fragmento | 14/14 OK |
| (a) `html.parser`: etiquetas sin cerrar · cierres cruzados · ids repetidos · `label for` huérfanos | 0 · 0 · 0 · 0 |
| (a) Nombres `[A-Za-z0-9.-]` · stems únicos en el lienzo (90 artboards) · títulos repetidos | OK · 0 repetidos · 0 |
| (a) Plantillas: `{{…}}` y `<sc-…>` · `renderVals()` | 0 · `return {}` en los 14. Todo el texto es literal, así que no hay enlaces de datos que puedan quedar vacíos |
| (b) Hex fuera de `data-props` (sin contar `href="#…"`) · `rgb()`, `hsl()`, `oklch()` · `fill`/`stroke` de SVG distintos de `none`/`currentColor` | 0 · 0 · 0 |
| (b) `var(--x)`: todos los nombres existen en `tokens.json` (356 nombres; `font-sans` sale de `type.families.sans`) y en `tokens.css` (356 nombres, el mismo conjunto) · `var()` con respaldo | OK · 0 |
| (c) Contraste de texto por cascada | 1,583 nodos de texto por tema y 42 pares (color, fondo, tema): **0 bajo 4.5:1**. El mínimo es 4.74 (`ui-role-success` sobre `ui-surface-card`, claro) |
| (c) Marcadores de posición (R3-A1) | 49 visibles («Sin notas» 37, «Sin tope» 6, «Sin restricciones» 4, «Sin etiqueta» 2). Los 14 artboards tienen `::placeholder{color:var(--ui-text-secondary);opacity:1}`. Resultado sobre el fondo efectivo de cada campo: **5.40 en claro y 10.43 en oscuro**, el par C-035 de DESIGN §10.5. 0 bajo umbral |
| (c) Contraste de UI 3:1 | tabla abajo: 0 bajo 3:1 |
| (d) P3-12: los 21 patrones de la Verifica (con y sin coma) · los 19 equivalentes en centavos · patrones de credencial | 0 · 0 · 0 |
| (d) Fechas que coinciden con el calendario real (comando de P3-12: 15 coincidencias) | En ajustes-b solo aparecen dos. `2026-10-30` está en s09-e, «El próximo pago (Vie 30 oct 2026) avisaría…», sin monto. `2027-02-28` está en s08-e, «por ejemplo, el 28-feb-2027», el ejemplo de U4, sin monto. OK |
| (d) Dinero visible contra `datos-sinteticos.json` | 65 apariciones con 30 valores distintos. 26 están en el JSON (como `_fmt` o como `_cents`). Los otros 4 son derivados y su aritmética está comprobada: 80 % y 85 % de $27,586.00 dan $22,068.80 y $23,448.10, y 40 % y 60 % de $2,918.40 dan $1,167.36 y $1,751.04 |
| (e) Fidelidad a U1/U4 | lista de comprobación por sección (61 elementos obligatorios): 0 faltantes en los escritorios. En móvil, lo que falta es justo lo que U4 deja fuera del móvil («Restaurar» va en el detalle y no en la tarjeta; la hoja de cada alerta va cerrada). Tabla abajo |
| (e) Estados (U1 §2.2 «según sección» y U4) | Los 6 escritorios con estados (s08, s09, s11, s12, s13, s14) traen la tarjeta «Estados» con Vacío, Carga (3 barras `.sk` de esqueleto) y Error. s10-e dice «Vacío, carga y error: no aplican», como U4 P-26.10 |
| (f) `button` con nombre y `type` · `a` con nombre · campos con `label` o `aria-label` · `svg` con `aria-hidden` o nombre · `role="switch"` con `aria-checked` · interactivo dentro de `a`/`button` | 0 fallas en todos |
| (f) Objetivos táctiles en móvil (alto por cascada: `height`/`min-height` resueltos o relleno + interlineado + borde; casillas medidas por su `label`) | 144 objetivos y 0 bajo 44 px. La única medición por debajo es el `input type="file"` oculto de s11-m (`.vh`, 1 px), y su objetivo real es `label.btn` «Elegir archivo» (44 px). El foco se dibuja con `.vh:focus-visible+.btn` |
| (f) Emoji · barra de estado o teclado falsos · `innerHTML`/`iframe`/`fetch`/`localStorage` · `type="date"` · fecha ISO visible | 0 · 0 · 0 · 0 · 0 |
| (f) Dinero con formato `$1,234.56` · «-$» con guion · montos sin `tabular-nums` | 65 de 65 bien formados · 0 · 0 |
| (f) Inglés visible | 0 palabras de interfaz en inglés. Siguen dentro de B-19 las teclas (Ctrl, Alt, Shift, Enter, Esc), «1.0.0-demo», «Ollama», «API de Claude» y el concepto «Streaming» (viene del dato compartido) |
| (g) Móvil: raíz `.mr{width:412px;min-height:915px}` con `*{box-sizing:border-box}`. Escritorio: `.pg{display:flex;flex-wrap:wrap}` sin ancho en px, `.mn{flex:999 1 560px;min-width:0}`, barra lateral con `flex:0 0 var(--layout-sidebar-width)`, 10 tablas dentro de `overflow-x:auto`, 0 `sticky`/`fixed`/`100vh`, y `"expand":"fill"` en el fragmento | 7/7 · 7/7 |
| (h) Enlaces | 242 `href`: 212 van a archivos de la lista del lienzo del mismo ancho y 30 a `#` (Buscar, «Abrir» de las hojas de alerta, «Ver detalle», «Editar», «Ver los n restantes»; ninguno tiene artboard propio). 0 rotos y 0 de otro ancho |
| Fragmento | JSON válido. `boards` = `order` = los 14 archivos, sin solaparse con otros fragmentos. `page:"ajustes"` e `is_interactive:false`. `w` es 412 o 1440. En móvil `h` ≥ 915 y sin `expand`. Títulos `P-26.xx · … · móvil|escritorio`, únicos. Sin `x`/`y`. Sin cambios contra iter-3 |
| Diferencias contra `docs/design/salida/iter-3/ajustes-b/` (diff por tokens) | Los 14 artboards cambian solo en la regla `::placeholder` al final del `helmet`. s14-e cambia además «Sin aviso» → `12 pza`, `2,000 ml` y `1,000 g` (ahora en `span.m`). s14-m cambia además «Sin aviso» → «nivel 12 pza», «nivel 2,000 ml» y «nivel 1,000 g». 0 regresiones |

### Recálculo independiente de contraste (hex de `tokens.json`)

| # | Par (texto / fondo) | Claro | Oscuro | Umbral |
|---|---|---|---|---|
| 1 | `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 ✓ |
| 2 | `ui-text-default` / `ui-surface-field` · `-page` | 11.53 | 14.07 | 4.5 ✓ |
| 3 | `ui-text-default` / `ui-surface-selected` · `-hover` | 10.73 · 9.64 | 10.90 · 7.42 | 4.5 ✓ |
| 4 | `ui-text-secondary` / `ui-surface-card` | 5.04 | 7.97 | 4.5 ✓ |
| 5 | `ui-text-secondary` / `ui-surface-selected` (vistas previas) | 5.02 | 8.08 | 4.5 ✓ |
| 6 | **`ui-text-secondary` / `ui-surface-field` (marcadores de posición, C-035; R3-A1)** | **5.40** | **10.43** | 4.5 ✓ |
| 7 | `ui-text-link` / `ui-surface-card` · `-page` | 5.16 · 5.53 | 7.94 · 10.38 | 4.5 ✓ |
| 8 | `ui-role-success` / `-success-surface` · `ui-surface-card` · `-field` | 4.81 · 4.74 · 5.08 | 6.19 · 6.42 · 8.40 | 4.5 ✓ |
| 9 | `ui-role-warning` / `-warning-surface` · `ui-role-error` / `-error-surface` · `ui-role-error` / `ui-surface-card` | 5.12 · 5.37 · 5.78 | 6.09 · 6.05 · 5.82 | 4.5 ✓ |
| 10 | `ui-text-on-primary` / `ui-action-primary-bg` · `ui-text-on-critical` / `ui-action-critical-bg` | 5.77 · 5.60 | 10.47 · 10.55 | 4.5 ✓ |
| UI | `ui-border-control` / `ui-surface-card` · `-field` · `-page` · `-raised` · `-selected` | 3.47 · 3.71 · 3.71 · 3.71 · 3.46 | 6.09 · 7.96 · 7.96 · 6.09 · 6.17 | 3 ✓ |
| UI | perilla apagada `ui-surface-page` / pista `ui-border-control` · perilla encendida `ui-text-on-primary` / `ui-action-primary-bg` | 3.71 · 5.77 | 7.96 · 10.47 | 3 ✓ |
| UI | pista encendida `ui-action-primary-bg` / `ui-surface-card` | 5.16 | 7.94 | 3 ✓ |
| UI | `ui-border-focus` / `ui-surface-card` · `-page` · `-selected` | 3.56 · 3.82 · 3.55 | 6.03 · 7.89 · 6.11 | 3 ✓ |
| UI | `ui-border-selected` / `ui-surface-field` · `-selected` | 5.53 · 5.15 | 10.38 · 8.04 | 3 ✓ |
| UI | `ui-viz-meter` · `ui-role-savings` / `ui-viz-track` (avance de metas) | 4.51 · 4.34 | 5.50 · 4.43 | 3 ✓ |

### Muestra de 10 valores contra `datos-sinteticos.json` (P3-12)
1. `$7,200.00` (`categorias[0].plan_fmt`).
2. `$14,260.00` (`ingresos.esperado.q15_cents`).
3. `$27,586.00` (`presupuesto.resumen.plan_gasto_fmt`).
4. `$20,525.78` (`presupuesto.resumen.real_fmt`).
5. `$21,480.00` (`ahorro.fondo.ahorrado_fmt`).
6. `$2,783.34` (`ahorro.metas[0].aportacion_necesaria_fmt`).
7. `$18,000.00` y `$14,100.00` (`ahorro.metas[1].monto_meta_fmt` y `.ahorrado_fmt`).
8. `$2,918.40` y `$1,459.20` (`creditos.kueski.compras[1].capital_fmt` y `.hogar.capital_fmt`; no son cifras canon).
9. **Inventario (H-11):** huevo 6 pza con nivel 12, leche 1,500 ml con nivel 2,000 y tortilla 700 g con nivel 1,000 (`inventario.articulos[2]`, `[3]` y `[5]`, `umbral_reposicion` y `existencia`). Resumen 1 / 5 / 8 (`inventario.resumen`).
10. Fechas de las metas: inicio 18-may-2026 y 04-may-2026, meta 17-ene-2027 y 17-abr-2027 (`ahorro.metas`). Pesos de salud: 84.6 kg con IMC 26.7 y 62.8 kg con IMC 23.6 (`salud.perfiles`).

Todos coinciden. La única fecha visible que no está en el JSON es «30-sep-2026», el fin del periodo de ejemplo del estado de cuenta de Karen (fin de mes, sin monto).

## Estado de los hallazgos de la vuelta 3

| Id | Sev. r3 | Estado r4 | Evidencia |
|---|---|---|---|
| R3-A1 marcadores de posición sin `::placeholder` | alta | **Cerrado** | `::placeholder{color:var(--ui-text-secondary);opacity:1}` está en el `helmet` de los 14 (14/14; también en los que no tienen marcadores). Los 49 marcadores visibles quedan en 5.40 en claro y 10.43 en oscuro sobre `ui-surface-field`, el par aprobado C-035 (DESIGN §10.5, línea 1237). Con `opacity:1`, Firefox ya no los atenúa. No cambian las alturas |
| H-11 «Reponer» con nivel «Sin aviso» frente a K7 | media | **Cerrado** | El dato lo corrigió `correccion-datos-r4`: `umbral_reposicion` vale 12 pza en huevo, 2,000 ml en leche y 1,000 g en tortilla, con estado «reponer». En s14-e: «Huevo blanco … 12 pza · 18 pza cada 7 días · 6 pza · Reponer · 2 días», «Leche entera … 2,000 ml … 1,500 ml · Reponer», «Tortilla de maíz … 1,000 g … 700 g · Reponer». En s14-m: «Ritmo de consumo · 6 pza · nivel 12 pza» y las otras dos igual. Aplicada la regla K7 de la pantalla, las tres filas dan «Reponer» (6 ≤ 12, 1,500 ≤ 2,000 y 700 ≤ 1,000). Los cuatro «Sin aviso» que quedan (croquetas, detergente, jabón para trastes y papel higiénico) tienen `umbral_reposicion: null` y estado «Suficiente», lo que coincide con U4 («vacío = sin aviso»). La vista previa de la pasta (400 g con nivel 300 g: Reponer → Suficiente) también cumple la regla. `compras-p16-inventario-m/e` muestran los mismos niveles («Reponer desde 12 pza», etc.) |
| B-01, B-03 a B-06, B-08 a B-16, B-18 a B-21 | baja | **Siguen abiertos** | Fuera del alcance de la ronda. B-08 y B-20 están sin cambios porque el diff de esas zonas es vacío. Los demás se comprobaron uno por uno con grep: B-01 (`body` a 14 px en 14/14), B-03 (`abaxel` `value="400"`), B-04 (4 campos «Hora»), B-05 (8 «Enter guarda»), B-06, B-09, B-10, B-11 (14 «Axel · Co-titular»), B-12, B-13 (2 `inputmode="numeric"`), B-14, B-15, B-16, B-19 y B-21 (13 `.tg .ic{16px}`) |

## Tabla artboard × regla

Columnas: F forma y fragmento · T tokens · C contraste (U5-04) · D cifras sintéticas (P3-12) · U fidelidad a U1/U4 · A accesibilidad y táctil · I idioma y formatos (U5-07, P3-08) · N navegación y enlaces · E estados (U1 §2.2).

| Artboard | F | T | C | D | U | A | I | N | E | Fallas (todas bajas) |
|---|---|---|---|---|---|---|---|---|---|---|
| s08-reportes-m | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✓ | ✓ | — | B-05, B-06 |
| s08-reportes-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ~ | ✓ | B-11, B-12, B-13, B-18, B-22, B-24 |
| s09-notificaciones-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05, B-09 |
| s09-notificaciones-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | B-04, B-09, B-11, B-12 |
| s10-apariencia-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05 |
| s10-apariencia-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | n/a | B-05, B-11, B-12 |
| s11-datos-respaldos-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | — | B-05, B-19 |
| s11-datos-respaldos-e | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ~ | ✓ | ✓ | B-11, B-12, B-14, B-19 |
| s12-personas-reparto-m | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✓ | ✓ | — | B-05, B-08, B-10 |
| s12-personas-reparto-e | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✓ | ✓ | ✓ | B-10, B-11, B-12 |
| s13-salud-metas-m | ✓ | ✓ | ✓ | ~ | ~ | ✓ | ~ | ✓ | — | B-03, B-05, B-15, B-16, B-20, B-22 |
| s13-salud-metas-e | ✓ | ✓ | ✓ | ~ | ~ | ✓ | ~ | ✓ | ✓ | B-03, B-11, B-12, B-15, B-16, B-20, B-22 |
| s14-inventario-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05, B-23 |
| s14-inventario-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | B-09, B-11, B-12, B-23 |

Notas de la tabla:
- «—» en E: el móvil no lleva fila de estados (la tarea la pide en escritorio).
- «n/a»: P-26.10 no tiene vacío, carga ni error (U4), y el artboard lo dice.
- B-01 y B-21 son transversales a los 14 artboards.

## Fidelidad: pantalla → elementos de U1/U4

| Sección | Presentes | Faltantes o distintos (bajas) |
|---|---|---|
| P-26.08 Reportes | Tabla Reporte · Activo · Día · Hora · Destinatarios · Próximo envío · Último resultado; «Enviar ahora a…» por fila; «Valor inicial»; 5 parámetros operativos; regla de los días 29 a 31 con el 28-feb-2027; Bajo demanda (mensual ahora con mes y casilla de Telegram, estado de cuenta de Karen con periodo, vista web mensual); estados. En móvil: tres tarjetas, hoja abierta con vista previa y Bajo demanda a todo el ancho | En móvil, «Reintentar el envío» en lugar de «Enviar ahora a…» (B-06). Hora semanal 20:00 en la tabla frente a 19:00 en Parámetros (B-18). «Generar estado de cuenta de Karen» no enlaza a P-25 (B-24) |
| P-26.09 Notificaciones | 4 alertas con activa, destinatarios, anticipación o umbral y hora, con «Valor inicial»; chats confirmados (etiqueta, recibe avisos) y detectados (Confirmar como, Descartar); «Enviar aviso de prueba»; vistas previas de 80 % y 100 % ($22,068.80 y $23,448.10) y del próximo pago; estados. En móvil: tarjetas por alerta y el chat detectado arriba | 4 campos «Hora» para la sola clave `alerts.send_time` (B-04); «•••» sin los 4 caracteres (B-09) |
| P-26.10 Apariencia | Tema segmentado Sistema · Claro · Oscuro; atajos de una tecla con su vista previa; «Saturación: apagada» y «Densidad: cómoda»; muestras clara y oscura; «no aplican» | «Enter guarda · Esc descarta» en controles de guardado inmediato (B-05) |
| P-26.11 Datos y respaldos | Encabezado con último, próximo y restauración probada; hora del respaldo diario (provisional); tabla de 5 respaldos con integridad, tamaño y ubicación; «Restaurar» con RESTAURAR; importar con «Elegir archivo» y «Lo que se creará» (10 registros del JSON); semilla; IA (Ollama y API de Claude, apagadas); servicio; 3 parámetros (10 MB, 24 h, 30 s); estados. En móvil: «Respaldar ahora» arriba y tarjetas sin «Restaurar» (U4) | «Restaurar respaldo» deshabilitado pintado como crítico (B-14) |
| P-26.12 Personas y reparto | Personas con nombre, papel, cuenta, paga directo, qué hace, tope (vacío = «Sin tope») y notas; 3 repartos con «Reparto por defecto», el resto calculado y la regla única; vista previa con $2,918.40; diálogo de alcance (Q2-68 C); estados. En móvil: tarjetas con barra de partes etiquetada | No hay control visible para editar «Parte de cada persona» (B-10) |
| P-26.13 Salud y metas | Fichas completas (estatura, nacimiento, sexo, actividad, objetivo, peso, IMC con clase, meta sugerida con desglose, meta confirmada con «Usar la sugerida», ajustes propios, restricciones y notas); ajuste de la casa con «Valor inicial»; metas con Tabla · Tablero · Calendario · Galería, «Nueva meta» y vista previa de «Activa»; estados | B-03, B-16 y B-20 siguen igual |
| P-26.14 Inventario | 14 productos con modo, unidad, nivel, ritmo, caducidad, existencia, estado, días y notas; resumen 1/5/8; Tabla · Tablero por modo; «Agregar producto al inventario»; vista previa K7; regla escrita; estados. En móvil: tarjetas con modo, estado y nivel. **H-11 resuelto** | Días restantes truncados a entero (3.5 → «3 días»), mientras compras-p16 los muestra con un decimal (B-23) |

## Líneas de evidencia

- **P3-04**: este revisor hizo 0 acciones de navegador y 0 llamadas a Artifact, y no encontró credenciales. `referencia-tipo/format.md` se leyó como dato que fija la forma. En los artboards y en el fragmento no hay texto con aspecto de instrucción.
- **P3-06 (cobertura de ajustes-b)**: 7 secciones × 2 anchos = 14/14, sin faltantes. El diff se hizo contra la copia de la vuelta 3 (`docs/design/salida/iter-3/ajustes-b/`). La copia de esta vuelta la hace el orquestador.
- **P3-07**: la tabla artboard × regla está arriba. Recalculé 10 pares de texto y 7 de UI. Por cascada hay 0 pares bajo umbral (1,583 nodos por tema). Los 49 marcadores de posición están en 5.40 / 10.43.
- **P3-08** (8 puntos):
  1. Resumen con ≤ 6 tarjetas: no aplica.
  2. 14 secciones de Ajustes: ✓. La columna «Secciones de Ajustes» enlaza las 14 y marca `aria-current`.
  3. Barra lateral colapsable: ✓ («Ocultar barra lateral»).
  4. Español de México: ✓, con la observación B-19.
  5. Dinero `$1,234.56` tabular: ✓ (65 de 65).
  6. Fechas `10-oct-2026`: ✓ (0 ISO, 0 `type="date"`), con la observación B-22.
  7. Paleta C2: ✓, solo tokens y 0 hex.
  8. Notion como única referencia: ✓.
- **P3-12**: 0 cifras canon (21 patrones con y sin coma, más 19 en centavos) y 0 credenciales. Las 2 fechas coincidentes van sin monto. La muestra de 10 valores coincide con el JSON. **Cumple**, con la observación B-15.
- **U5-04**: 0 pares de texto bajo 4.5:1 y 0 de UI bajo 3:1, en claro y en oscuro, incluidos los marcadores de posición. **Cumple** (en lo verificable en el markup).
- **U5-06**: cumplen tokens, barra lateral colapsable, objetivos de 44 px, foco visible y cifras tabulares. Quedan B-01 y B-21 (tamaños tipográficos y algunos espaciados literales, iconos de insignia de 16 px).
- **U5-07**: 0 palabras de interfaz en inglés, 0 fechas ISO y dinero bien formado. Quedan nombres fuera de la lista permitida (B-19).

## Hallazgos

### Alta
Ninguno.

### Media
Ninguno.

### Baja (no bloquean)
Abiertos de vueltas anteriores, sin cambios. El detalle está en `P3-iter-3-ajustes-b.md`.
- **B-01 · DESIGN §7.2, §7.3 y §7.10.** Texto base y campos a 14 px. Título de escritorio a 24 px. `.mx{max-width:1180px}` y los tamaños de letra son literales.
- **B-03 · U4 §0.2.** En s13, «Mi ajuste para bajar» de Axel muestra 400 y la vista previa dice 400 → 500. El campo debe mostrar 500.
- **B-04.** s09-e tiene 4 campos «Hora» para la sola clave `alerts.send_time`.
- **B-05 · U4 §0.7.** «Enter guarda · Esc descarta» aparece en móviles y en controles de guardado inmediato (s10).
- **B-06 · U4 P-26.08.** s08-m dice «Reintentar el envío» en lugar de «Enviar ahora a…».
- **B-08.** En la barra de reparto de s12-m, los segmentos de Karen y Axel van pegados (1.51:1 entre sí en oscuro). El dato también va en el texto.
- **B-09.** La columna interna va sin iconos. La terminación del chat es «•••» sin los 4 caracteres. s14-e tiene 14 rótulos «Modo» iguales.
- **B-10 · U4 P-26.12.** «Parte de cada persona» no tiene control de edición visible.
- **B-11.** La barra lateral dice «Axel · Co-titular» (sale de `ajustes.usuarios[0].rol`), pero s12-e y `personas[0].rol` dicen «Titular». Viene del shell y lo unifica el orquestador.
- **B-12 · DESIGN §7.9.** La barra lateral no tiene `chevron-right` y usa `aria-label="Principal"`.
- **B-13.** Los campos de fecha de s08-e tienen `inputmode="numeric"` con un formato que lleva letras.
- **B-14 · DESIGN §7.2.** «Restaurar respaldo» tiene `aria-disabled="true"` pero se pinta con el color crítico completo.
- **B-15 · P3-12.** Las fechas de nacimiento 14-mar-1990 y 22-jun-1993 no están en `datos-sinteticos.json`.
- **B-16 · U4 P-26.13 (Q2-17 C, Q2-18 A) frente a DESIGN §3.3.** El avance es entero (42 %, 78 %) y los estados son «Atrasada» y «En ruta», como en el JSON (`ahorro.metas[].avance_pct`, `etiqueta_estado`). Decide el orquestador.
- **B-18.** En s08-e, la hora semanal dice 20:00 en la tabla y 19:00 en «Parámetros operativos».
- **B-19 · U5-07.** Fuera de la lista permitida: «1.0.0-demo», «Modelo local (Ollama)» y «API de Claude». También el concepto «Streaming» de s11 (dato compartido `presupuesto.gastos_fijos_recurrentes[4].concepto`, que usan 10 artboards del lienzo; es un extranjerismo, y si se quiere cambiar lo hace correccion-datos).
- **B-20 · U4 P-26.13 (K36, U4-13).** Los ajustes de Ximena muestran «0» como valor propio, pero la pantalla dice que ella usa el de la casa. Falta «Usar el de la casa» en el ajuste propio de Axel.
- **B-21 · DESIGN §5.2.** `.tg .ic{width:16px;height:16px}` en 13 artboards. Hay literales fuera de la escala de 4 pt: `gap:2px`, `gap:6px`, `padding:2px` y `border-radius:12px`.

Nuevos en esta vuelta (también bajos):
- **B-22 · U1 §0.4 (fechas `10-oct-2026`), coherencia del lienzo.** Tres fechas llevan el día con cero a la izquierda: «01-oct-2026» (último resultado del mensual en s08-e), `value="01-sep-2026"` (inicio del periodo de Karen en s08-e) y `value="04-may-2026"` (inicio de «Viaje a Oaxaca» en s13-m y s13-e). El resto del lienzo escribe el día sin cero: «1-oct-2026», «2-nov-2026», «5-oct-2026», etc. (más de 70 apariciones). **Corrección:** «1-oct-2026», «1-sep-2026» y «4-may-2026».
- **B-23 · coherencia entre áreas (decide el orquestador).** s14-m/e muestran los días restantes truncados a entero: 2.3 → «2 días», 3.0 → «3 días», 3.5 → «3 días», 20.2 → «20 días». `compras-p16-inventario-m/e` muestran el mismo dato con un decimal («2.3 días restantes», «3.0 días»). Ni U1 ni U4 fijan el redondeo. Conviene un solo criterio en el lienzo y, si es entero, decir cuál (truncar o ROUND_HALF_UP; con ROUND_HALF_UP la tortilla daría 4).
- **B-24 · (h) prototipo.** U4 dice que «Generar estado de cuenta de Karen» «al terminar abre P-25», y el artboard `salud-reportes/reportes-p25-karen-{m,e}` existe. Aun así, el botón de s08-m/e no enlaza ahí, y tampoco hay enlace de «Reporte en preparación · Ver en el Centro de reportes» a P-23 desde el bloque Bajo demanda. Es opcional porque son acciones y no navegación. **Corrección:** si el orquestador lo quiere, agregar un `<a>` secundario «Ver en el Centro de reportes» o «Abrir el estado de cuenta» con el `href` del mismo ancho.

## Correcciones sugeridas (ninguna bloquea)
1. Baratas, de CSS y texto: B-03, B-05, B-12, B-13, B-14, B-20, B-21 y B-22.
2. Para el orquestador o correccion-datos: B-11, B-15, B-16, B-19 y B-23.
3. Del prototipo: B-24.
