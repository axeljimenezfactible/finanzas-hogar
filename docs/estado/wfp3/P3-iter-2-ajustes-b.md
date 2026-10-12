# P3-iter-2 · ajustes-b · revisión de la vuelta 2

Fecha: 12-oct-2026 · Revisor: `revisor-ajustes-b-r2` (independiente, vuelta 2 de 3; no corrige) · Criterio: P3-07 de `docs/criterios/P3-claude-design.md` 1.3.
Alcance: 14 artboards de `docs/design/maquetas/lienzo/project/ajustes/` (P-26.08 a P-26.14, `-m` y `-e`) y el fragmento `docs/design/maquetas/lienzo/fragmentos/ajustes-b.json`, tal como los dejó `correccion-ajustes-b-r1`.
Contra: P3-08 y P3-12; U5-04, U5-06 y U5-07 (lo verificable en el markup); DESIGN.md (§3.3, §7.2, §7.3, §7.9, §7.10, §7.13, §9.3, §10.5); U1 §1.3, §2.1, §2.2 y §2.3; U4 §0 y P-26.08 a P-26.14; `referencia-tipo/format.md`.
Ruta del informe: la que fijó la tarea (`docs/estado/wfp3/`). La copia a `docs/design/salida/iter-2/` (P3-06) la hace el orquestador.

**Veredicto: NO APROBADO** · 1 alta · 2 medias · 15 bajas.
De los 12 hallazgos de la vuelta 1: 11 cerrados y 1 reabierto en parte (H-11, la mitad K7, que depende de `datos-sinteticos.json`). Hay 1 alta nueva que la vuelta 1 no detectó (texto con `opacity:.85`, también en la copia de iter-1) y B-02 pasa de baja a media.

## Verificación ejecutada (scripts propios, `python3 -I`, en `datos/tmp/rev-ajb-r2/`, borrados al final)

| Comprobación | Resultado |
|---|---|
| (a) Forma: empieza con `<!doctype html>`, `lang="es-MX"`, charset, `<title>` propio y único, `<script src="./support.js">` una vez, `<link … ../ds/finanzas-hogar/tokens.css>` una vez en `<head>`, `<body><x-dc><helmet><style>`, un solo `x-dc`, bloque `data-dc-script` después de `</x-dc>` con `class Component extends DCLogic` y `renderVals()`, sin import ni export | 14/14 OK |
| (a) `data-props` parsea (tras decodificar entidades) y `$preview` = `w`×`h` del fragmento | 14/14 OK |
| (a) `html.parser`: etiquetas sin cerrar, cierres huérfanos, ids repetidos, `label for` huérfanos | 0 · 0 · 0 · 0 |
| (a) Nombres `[A-Za-z0-9.-]`, stems únicos en todo el lienzo (90 artboards) | OK · 0 duplicados |
| (b) Hex fuera de `data-props` (sin contar `href="#…"`); `rgb()`, `hsl()`, `oklch()` | 0 · 0 |
| (b) `var(--x)` contra `tokens.json` y `tokens.css`; `var()` con respaldo; propiedades locales | todos existen · 0 · 0 |
| (c) Contraste de texto: cascada CSS propia (reglas del `helmet` + `style` en línea, herencia de `color`, fondo efectivo, `data-theme` local, **opacidad**) en claro y oscuro | 1,620 nodos por tema, 22 pares distintos: **0 bajo 4.5:1 por color; 2 nodos con `opacity:.85` bajan a 3.72:1 en claro** (R2-A1) |
| (c) Contraste de UI 3:1: perilla y pista del interruptor, pista y tarjeta, borde de 141 campos | 0 bajo 3:1 (H-01 cerrado). Los 16 avisos de «botón sin borde» del script son botones de `.seg`, cuyo contenedor lleva `ui-border-control` (3.47 / 6.09): no son fallas |
| (d) P3-12: los 21 patrones de la Verifica, con y sin coma, en centavos y como credenciales | 0 · 0 · 0 |
| (d) Fechas que coinciden con el calendario real (comando de P3-12) | en ajustes-b solo `2026-10-30` (s09-e, «Vie 30 oct 2026», próximo pago) y `2027-02-28` (s08, ejemplo de U4 «28-feb-2027»). Ninguna va con monto: OK |
| (d) Dinero y fechas visibles contra `datos-sinteticos.json` | 65 apariciones de montos: 57 del JSON y 8 derivadas con aritmética comprobada; 46 fechas de texto: todas del JSON. Fechas en `value`: `14-mar-1990` y `22-jun-1993` **no están en el JSON** (B-15); `18-may-2026` y `04-may-2026` sí |
| (e) Fidelidad a U1/U4 | tabla abajo |
| (f) `button` con nombre y `type`; `a` con nombre; campos con `label`/`aria-label`; `svg` decorativo con `aria-hidden`; `role="switch"` con `aria-checked` y nombre | 0 fallas; 30 interruptores OK |
| (f) Objetivos táctiles en móvil (altura por cascada: `height`/`min-height` con tokens resueltos, o relleno + interlineado) | 142 objetivos, **0 bajo 44 px** (H-02 cerrado); barra inferior 56 px |
| (f) Emoji · barra de estado o teclado falsos · `innerHTML`/`iframe`/red/`localStorage` · `type="date"` · fecha ISO visible · `<code>` | 0 · 0 · 0 · 0 · 0 · 0 |
| (f) Inglés visible | 0 claves técnicas (H-03 cerrado). Quedan nombres de tecla (Ctrl, Alt, Enter, Esc) y «1.0.0-demo», «Ollama», «API de Claude» (B-19) |
| (f) Dinero `$1,234.56`, signo «−» U+2212, `tabular-nums` | 0 mal formados · 0 «-$» · **7 montos en prosa sin `tabular-nums`** (R2-M1) |
| (g) Móvil: `.mr{width:412px;min-height:915px}` y `box-sizing:border-box`; barra inferior fija | 7/7 |
| (g) Escritorio: `.pg` flex-wrap sin ancho en px; barra `flex:0 0 var(--layout-sidebar-width)`; `.mn{flex:999 1 560px;min-width:0}`; tablas en `.tw{overflow-x:auto}`; fragmento con `"expand":"fill"` | 7/7 |
| (h) Enlaces: 242 `href`, 30 a `#`; los demás apuntan a archivos que existen y están en la lista del lienzo, del mismo ancho | 0 rotos · 0 de otro ancho |
| Fragmento: JSON válido, `boards` = `order` = los 14 archivos, `page:"ajustes"`, `is_interactive:false`, `w` 412/1440, móvil `h` ≥ 915, títulos `P-26.xx · … · móvil|escritorio` | OK. Las alturas (medidas por el corrector con chromium) no se volvieron a medir: este revisor no usa navegador |

### Recálculo independiente de contraste (P3-07: 10 pares de texto + UI, hex de `tokens.css`)

| # | Par (texto / fondo) | Claro | Oscuro | Umbral |
|---|---|---|---|---|
| 1 | `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 ✓ |
| 2 | `ui-text-secondary` / `ui-surface-card` | 5.04 | 7.97 | 4.5 ✓ |
| 3 | `ui-text-default` / `ui-surface-field` | 11.53 | 14.07 | 4.5 ✓ |
| 4 | `ui-text-secondary` / `ui-surface-selected` | 5.02 | 8.08 | 4.5 ✓ |
| 5 | `ui-role-success` / `ui-role-success-surface` | 4.81 | 6.19 | 4.5 ✓ |
| 6 | `ui-role-warning` / `ui-role-warning-surface` | 5.12 | 6.09 | 4.5 ✓ |
| 7 | `ui-text-on-primary` / `ui-action-primary-bg` | 5.77 | 10.47 | 4.5 ✓ |
| 8 | `ui-text-on-critical` / `ui-action-critical-bg` | 5.60 | 10.55 | 4.5 ✓ |
| 9 | `ui-text-secondary` / `ui-surface-raised` (barra inferior) | 5.40 | 7.97 | 4.5 ✓ |
| 10 | `ui-role-success` / `ui-surface-card` («Guardado») | 4.74 | 6.42 | 4.5 ✓ |
| ✗ | `ui-text-secondary` con `opacity:.85` / `ui-surface-card` (s13: #66676e → #7b7c82 sobre #f2f2f4) | **3.72** | 6.28 | 4.5 ✗ |
| UI | perilla `ui-text-on-primary` / pista encendida `ui-action-primary-bg` (C-043) | 5.77 | 10.47 | 3 ✓ |
| UI | perilla `ui-surface-page` / pista apagada `ui-border-control` (C-045) | 3.71 | 7.96 | 3 ✓ |
| UI | pista encendida / `ui-surface-card` · pista apagada / `ui-surface-card` | 5.16 · 3.47 | 7.94 · 6.09 | 3 ✓ |
| UI | borde de campo `ui-border-control` / `ui-surface-card` · `ui-surface-raised` · `ui-surface-page` | 3.47 · 3.71 · 3.71 | 6.09 · 6.09 · 7.96 | 3 ✓ |

### Muestra de valores contra `datos-sinteticos.json` (P3-12)
`$7,200.00`, `$1,086.00`, `$238.40` y `$684.30` (presupuesto.gastos_fijos_recurrentes, con su categoría y día) · `$14,260.00` y `$14,780.00` (ingresos.esperado q15/q30) · `$21,480.00` y `$800.00` (ahorro.fondo) · `$14,500.00`, `$6,150.00`, `$2,783.34` (meta laptop) · 84.6 kg, IMC 26.7, 2,364 kcal; 1,816 y 1900 kcal · respaldos 21-oct-2026 03:00 de 2.3 MB a 17-oct-2026 de 2.7 MB · «Ximena (tablet)» 20-oct-2026 18:05 · semilla 33 productos y 5 compras de Kueski Pay · 1.0.0-demo, en marcha desde 12-oct-2026. Derivadas: 80 % y 85 % de $27,586.00 = $22,068.80 y $23,448.10 ✓; 40 % y 60 % de $2,918.40 = $1,167.36 y $1,751.04 ✓; días de inventario con piso entero: 2, 3, 3, 20, 16, 18 y 17 ✓; edades a hoy (21-oct-2026) con las fechas pintadas: 36 y 33 ✓.

## Estado de los hallazgos de la vuelta 1

| Id | Sev. r1 | Estado r2 | Evidencia |
|---|---|---|---|
| H-01 perilla del interruptor | alta | **Cerrado** | `.sw i:after{background:var(--ui-surface-page)}`, `.sw.on i:after{background:var(--ui-text-on-primary)}` en los 10 artboards; 30 pares perilla/pista ≥ 3.71 en los dos temas |
| H-02 regla `.in` | media | **Cerrado** | `.in,.sl{min-height:var(--size-control-coarse);…border…ui-border-control;background:ui-surface-field}` en s11, s12 y s14-m; campos de 44 px |
| H-03 claves técnicas en `<code>` | media | **Cerrado** | 0 `<code>`, 0 claves con punto o guion bajo en texto visible |
| H-04 «Valor inicial» junto al campo (s09-e) | media | **Cerrado** | cada campo de la tabla de alertas lleva «Valor inicial: …» en su celda; se quitó la lista duplicada; «Primer aviso» en edición (85, borde `ed`) cuadra con la vista previa 80 → 85 |
| H-05 `type="date"` (s08-e) | media | **Cerrado** | texto `01-sep-2026` y `30-sep-2026` con botón de calendario de 44 px. Queda B-13 (`inputmode="numeric"`) |
| H-06 barra inferior móvil | media | **Cerrado** | 7/7 con Inicio · Comidas · «+» (`shell-p03-registrar-gasto-m`) · Compras · Buscar, 56 px |
| H-07 barra lateral | media | **Cerrado** | avatar y rol, «Ocultar barra lateral» (`panel-left`, Ctrl+\), Buscar y Registrar gasto, 9 destinos en el orden del shell, pie con Tema, `aria-current="page"`. Quedan B-11 y B-12 |
| H-08 muestra clara sin `data-theme` (s10-e) | media | **Cerrado** | `data-theme="light"` y `data-theme="dark"` en las dos muestras |
| H-09 «Lo que se creará» y selector de archivo (s11) | media | **Cerrado** | `input type="file"` oculto accesible con `label` «Elegir archivo» (`.btn`, 44 px) y foco en `.vh:focus-visible+.btn`; tabla de 10 renglones en escritorio y lista apilada en móvil, del JSON |
| H-10 campos faltantes | media | **Cerrado** | s13: Notas del perfil, Inicio y Notas de cada meta; s14: Notas (columna en escritorio, `details` en móvil); s11: «Hora del respaldo diario» 03:00 con «Sin valor inicial · Provisional» y «Huella: —» |
| H-11 días con decimal y K7 | media | **Reabierto en parte** | días en piso entero ✓; huevo, leche y tortilla siguen «Reponer» con nivel «Sin aviso» junto a la regla K7 escrita en la misma pantalla (s14-m y s14-e). `datos-sinteticos.json` no cambió desde 11-oct 23:55 |
| H-12 fecha de nacimiento vacía (s13) | media | **Cerrado en pantalla** | 14-mar-1990 y 22-jun-1993, coherentes con 36 y 33 años. El dato no está en el JSON: B-15 |

## Tabla artboard × regla

Columnas: F forma y fragmento · T tokens · C contraste (U5-04) · D cifras sintéticas (P3-12) · U fidelidad a U1/U4 · A accesibilidad y táctil · I idioma y formatos (U5-07, P3-08) · N navegación (barra lateral o inferior, enlaces) · E estados (U1 §2.2).

| Artboard | F | T | C | D | U | A | I | N | E | Fallas |
|---|---|---|---|---|---|---|---|---|---|---|
| s08-reportes-m | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✓ | ✓ | — | B-05, B-06 |
| s08-reportes-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | B-01, B-11, B-12, B-13, B-18 |
| s09-notificaciones-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05 |
| s09-notificaciones-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | R2-M1, B-04, B-09, B-11, B-12 |
| s10-apariencia-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | — | B-05 |
| s10-apariencia-e | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | n/a | B-05, B-11, B-12 |
| s11-datos-respaldos-m | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | — | B-05, B-19 |
| s11-datos-respaldos-e | ✓ | ✓ | ✓ | ✓ | ✓ | ~ | ~ | ✓ | ✓ | B-11, B-12, B-14, B-19 |
| s12-personas-reparto-m | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✗ | ✓ | — | R2-M1, B-05, B-08, B-10 |
| s12-personas-reparto-e | ✓ | ✓ | ✓ | ✓ | ~ | ✓ | ✗ | ✓ | ✓ | R2-M1, B-10, B-11, B-12 |
| s13-salud-metas-m | ✓ | ✓ | ✗ | ~ | ~ | ✓ | ✓ | ✓ | — | R2-A1, B-03, B-05, B-15, B-16, B-17 |
| s13-salud-metas-e | ✓ | ✓ | ✗ | ~ | ~ | ✓ | ✓ | ✓ | ✓ | R2-A1, B-03, B-11, B-12, B-15, B-16, B-17 |
| s14-inventario-m | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | — | H-11 (K7), B-05 |
| s14-inventario-e | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | H-11 (K7), B-11, B-12 |

«—» en E: el móvil no lleva fila de estados (la tarea la pide en escritorio). «n/a»: P-26.10 no tiene vacío, carga ni error (U4). B-01, B-07 y B-09 son transversales.

## Fidelidad: pantalla → elementos de U1/U4

| Sección | Presentes | Faltantes o distintos |
|---|---|---|
| P-26.08 Reportes | tabla Reporte · Activo · Día · Hora · Destinatarios · Próximo envío · Último resultado; «Enviar ahora a…» por fila; «Valor inicial» en los campos de hora; vista previa del próximo envío; regla de días 29–31; Bajo demanda (mes, casilla de Telegram, estado de cuenta de Karen con fechas en texto, vista web); estados | móvil con «Reintentar el envío» y sin «Enviar ahora a…» (B-06); campos y lista de parámetros duplicados en escritorio (B-18) |
| P-26.09 Notificaciones | 4 alertas con activa, destinatarios, anticipación o umbral y hora, cada campo con «Valor inicial»; chats detectados (nombre, fecha, Confirmar como, Descartar con diálogo crítico); chats confirmados; aviso de prueba; vistas previas de umbral y de pago; estados | 4 campos «Hora» para una sola clave (B-04); terminación del chat sin los 4 caracteres, que el JSON no trae (B-09) |
| P-26.10 Apariencia | tema segmentado con `role="radio"`, atajos de una tecla, valores fijos, muestras clara y oscura con su `data-theme` | «Enter guarda · Esc descarta» en un interruptor que guarda al instante (B-05) |
| P-26.11 Datos y respaldos | encabezado último/próximo/restauración; hora del respaldo diario (provisional); tabla de 5 respaldos con ubicación y huella; restaurar con paso 2 y `RESTAURAR`; importar con «Elegir archivo» y «Lo que se creará» (10 renglones); semilla (3); IA (2, apagadas, llave); servicio; 3 parámetros con vista previa; estados | «Restaurar respaldo» deshabilitado sin aspecto de deshabilitado (B-14) |
| P-26.12 Personas y reparto | Axel, Ximena y Karen con papel, cuenta, paga directo, qué hace, tope y notas; 3 repartos con resto y «Reparto por defecto»; V-16; vista previa con $2,918.40; diálogo de alcance (Q2-68 C) | sin control visible para «Parte de cada persona» (B-10); montos en prosa sin `tabular-nums` (R2-M1) |
| P-26.13 Salud y metas | fichas completas con fecha de nacimiento, notas, desglose de la meta, «Usar la sugerida»; ajuste de la casa (2) con «Valor inicial»; metas con Inicio, Notas, Activa y vista previa; selector Tabla · Tablero · Calendario · Galería; «Nueva meta»; estados | nota atenuada bajo contraste (R2-A1); ajustes de Ximena no dibujados, solo descritos (B-17); avance entero y estados «Atrasada» y «En ruta» frente a Q2-17 C y Q2-18 A (B-16); «Mi ajuste para bajar» muestra 400 y la vista previa dice 400 → 500 (B-03) |
| P-26.14 Inventario | 14 artículos con modo, unidad, nivel, ritmo, caducidad, existencia, estado, días enteros y notas; resumen 1/5/8; Tabla · Tablero por modo; «Agregar producto»; vista previa; estados | «Reponer» con «Sin aviso» frente a la regla K7 de la misma pantalla (H-11) |

## Líneas de evidencia

- **P3-04**: 0 acciones de navegador y 0 llamadas a Artifact de este revisor; 0 credenciales en los 14 artboards. `referencia-tipo/format.md` trae guía en imperativo del publicador del tipo y se leyó como dato que fija la forma. En los artboards y el fragmento no hay texto con aspecto de instrucción.
- **P3-06 (cobertura de ajustes-b)**: 7 secciones × 2 anchos = 14/14; 0 faltantes. La copia de la vuelta 1 está en `docs/design/salida/iter-1/ajustes-b/` (14 + fragmento).
- **P3-07**: tabla artboard × regla arriba; 10 pares de texto y 6 de UI recalculados; 1 par de texto bajo 4.5:1 por opacidad (R2-A1); 0 pares de UI bajo 3:1.
- **P3-08** (8 puntos): (1) Resumen ≤ 6 tarjetas: no aplica · (2) 14 secciones de Ajustes: ✓ (la columna interna enlaza las 14) · (3) barra lateral colapsable: ✓ (botón «Ocultar barra lateral», `panel-left`) · (4) español de México: ✓, con B-19 · (5) dinero `$1,234.56` tabular: ✗, 7 montos en prosa sin `tabular-nums` (R2-M1) · (6) fechas `10-oct-2026`: ✓ (0 ISO, 0 `type="date"`) · (7) paleta C2 índigo 280°: ✓, solo tokens · (8) Notion como única referencia: ✓.
- **P3-12**: 0 cifras canon (21 patrones, con y sin coma, y en centavos); 0 credenciales; fechas coincidentes revisadas a mano, ninguna con monto; muestra de valores del JSON. Las fechas de nacimiento pintadas no están en el JSON (B-15). Se buscaron `14-mar-1990`, `22-jun-1993` y sus formas compactas en la memoria del proyecto, el oráculo y la semilla real: 0 archivos, así que no son datos reales. **Cumple**, con la observación B-15.
- **U5-04**: 0 pares de color bajo umbral; 1 texto con opacidad baja a 3.72:1 en claro (R2-A1). **Falla.**
- **U5-06**: tokens ✓; iconos Lucide de 20 px con trazo 1.5 ✓; barra lateral colapsable ✓; objetivos de 44 px ✓; foco visible ✓; cifras tabulares ✗ (R2-M1); tipografía de 14 px en campos móviles (B-01).
- **U5-07**: 0 claves en inglés ni fechas ISO; dinero bien formado; «1.0.0-demo», «Ollama» y «Claude» no están en la lista permitida (B-19).

## Hallazgos

### Alta
- **R2-A1 · U5-04, P3-07 (texto ≥ 4.5:1).** En s13-m y s13-e, `<p class="mu" style="opacity:.85">Con “mantener” no hay ajuste: los ajustes para bajar y subir se ven atenuados.</p>` mezcla `ui-text-secondary` (#66676e) con la tarjeta (#f2f2f4) y queda en #7b7c82: **3.72:1 en claro** (oscuro, 6.28:1). No es un control deshabilitado, así que no tiene la excepción de WCAG. Ya estaba en la copia de iter-1 y la vuelta 1 no lo midió. **Corrección:** quitar la opacidad del texto (queda en 5.04 / 7.97). La atenuación va en los dos campos de ajuste de Ximena, dibujados y deshabilitados (`aria-disabled="true"`, `ui-text-disabled`), con este texto como explicación (U4 P-26.13, ver B-17).

### Media
- **H-11 (reabierto en parte) · U4 P-26.14 y K7.** En s14-m y s14-e, huevo, leche y tortilla dicen «Reponer» con nivel «Sin aviso», y la misma pantalla escribe la regla «reponer si la existencia es igual o menor que el nivel». Según U4, nivel vacío significa sin aviso, así que la regla no puede producir «Reponer». El origen es `inventario.articulos` del JSON, que no cambió. Los días ya salen enteros. **Corrección:** correccion-datos da a esos artículos un `umbral_reposicion` coherente con su existencia (o los pasa a «suficiente» y ajusta `inventario.resumen`), y el autor pinta lo que el JSON diga. Lo mismo aplica a `compras-p16-inventario`. Mientras el JSON no cambie, este hallazgo no cierra.
- **R2-M1 (era B-02) · P3-08 punto 5, DESIGN §3.3 («cifras tabulares en todo dinero»), U5-06.** Hay 7 montos en prosa sin `tabular-nums`. En s09-e: «Llega al gastar (plan $27,586.00)», «Con el presupuesto total de este mes ($27,586.00)» y «…llega al gastar $27,586.00». En s12-m y s12-e: «Con un gasto de $2,918.40 · Karen» y «… · hogar (el resto)». Sube de baja a media porque P3-08 lo pone entre las decisiones y la vuelta 2 no lo tocó. **Corrección:** envolver cada monto en `<span class="m">`.

### Baja
- **B-01 · DESIGN §7.2, §7.3 y §7.10.** El texto base es de 14 px también en los campos y botones móviles (§7.2 y §7.3 piden 16 px). El título de página de escritorio mide 24 px y `title-1` es de 32 px. `.mx{max-width:1180px}` es un literal y no un token de layout.
- **B-03 · U4 §0.2.** En s13, «Mi ajuste para bajar» muestra 400 con borde de edición y la vista previa dice 400 → 500. El campo debe mostrar 500, como ya se hizo en s09 (85) y s08 (20:00).
- **B-04.** s09-e tiene cuatro campos «Hora», uno por alerta, para una sola clave (`alerts.send_time`). Conviene un campo único o el rótulo «Misma hora para todos los avisos».
- **B-05 · U4 §0.7 («En Android no hay atajos») y P-26.10.** Las vistas previas de los 7 móviles dicen «Enter guarda · Esc descarta». En s10 (m y e) el interruptor de atajos guarda al instante. Quitar la línea en móvil y en los controles de guardado inmediato.
- **B-06 · U4 P-26.08.** s08-m trae «Reintentar el envío», que U4 no tiene, y no trae «Enviar ahora a…».
- **B-07.** No hay `::placeholder{color:var(--ui-text-secondary)}`, así que «Sin notas», «Sin tope» y «dd-mmm-aaaa» toman el gris del navegador.
- **B-08.** En la barra de reparto de s12-m, los segmentos de Karen y Axel van pegados y entre sí miden 1.51:1 en oscuro. El dato también va en texto, pero un hueco de 2 px ayudaría.
- **B-09.** La columna interna de las 14 secciones va sin los iconos de §5.4. La terminación del chat es «•••» sin 4 caracteres, que el JSON no trae. En s14-e hay 14 etiquetas «Modo» iguales. Las alturas del fragmento no se volvieron a medir en esta vuelta.
- **B-10 · U4 P-26.12.** «Parte de cada persona» no tiene control de edición visible, aunque la vista previa muestra una edición (50.00 % → 40.00 %).
- **B-11 (nueva) · U4 P-26.12 y `personas[].rol`.** El encabezado de la barra lateral dice «Axel · Co-titular» en los 7 escritorios, pero s12-e y el JSON dicen «Titular». U1 §2.3 (P-26.01) dice «co-titulares», así que los documentos no coinciden; el texto viene de `shell-barra-lateral`. El orquestador debe unificarlo a lo que diga el JSON.
- **B-12 (nueva) · DESIGN §7.9.** La barra lateral no tiene `chevron-right` en las 9 secciones y es `<nav aria-label="Principal">`, no `"Secciones"`.
- **B-13 (nueva).** En s08-e, los campos de fecha de Karen tienen `inputmode="numeric"`, pero el formato lleva letras (`01-sep-2026`), y en un teléfono el teclado numérico no las tiene. Quitar `inputmode` o dejarlo en `text`.
- **B-14 (nueva) · DESIGN §7.2.** En s11-e, «Restaurar respaldo» tiene `aria-disabled="true"` hasta que se escriba RESTAURAR, pero se pinta con el color crítico completo y no se ve deshabilitado.
- **B-15 (nueva) · P3-12 (un solo archivo de datos).** Las fechas de nacimiento 14-mar-1990 y 22-jun-1993 se derivaron de `edad` y no están en `datos-sinteticos.json`. No coinciden con datos reales. Pedir a correccion-datos `salud.perfiles[].fecha_nacimiento` con esos valores, para que P-21 y Ajustes cuadren.
- **B-16 (nueva) · U4 P-26.13 (Q2-17 C, Q2-18 A) frente a DESIGN §3.3 y U1 §2.5.** El avance sale entero (42 %, 78 %) y el estado dice «Atrasada» y «En ruta». U4 pide un decimal en Metas y los estados En curso, Cumplida y Vencida. DESIGN §3.3 (avance entero) y U1 §2.5 (en ruta y atrasada) dicen otra cosa, y el JSON trae los valores viejos. El orquestador decide y correccion-datos ajusta.
- **B-17 (nueva) · U4 P-26.13.** Para Ximena (objetivo «Mantener») no se dibujan los dos ajustes atenuados: solo hay una frase que lo describe («se ven atenuados»). Ese texto es una descripción de la maqueta, no copia del producto.
- **B-18 (nueva).** En s08-e, día y hora de cada reporte se editan en dos sitios: la tabla y la lista «Parámetros operativos». La hora semanal se ve 20:00 en la tabla (en edición) y 19:00 en la lista, y los dos campos de día no llevan «Valor inicial» en la tabla. Basta con un solo sitio, como en s09-e.
- **B-19 (nueva) · U5-07.** «1.0.0-demo» (de `ajustes.servicio`), «Modelo local (Ollama)» y «API de Claude» (nombres de U4) no están en la lista permitida de U5-07. Hay que agregarlos a la lista o cambiar «demo» en el JSON.

## Correcciones para la vuelta 3 (orden sugerido)
1. R2-A1: quitar `opacity:.85` del texto de s13-m y s13-e y dibujar los ajustes atenuados de Ximena (cierra también B-17).
2. R2-M1: envolver en `.m` los 7 montos en prosa de s09-e, s12-m y s12-e.
3. H-11 (K7): escalar a correccion-datos (`inventario.articulos` y `resumen`) y repintar s14-m y s14-e con el JSON corregido. Si el JSON no cambia en esta vuelta, el orquestador debe registrar el bloqueo en `docs/estado/bloqueos/`.
4. Bajas baratas: B-03, B-05, B-07, B-12, B-13 y B-14 (CSS y texto).
5. Para el orquestador o correccion-datos: B-11 (rol en la barra lateral del shell), B-15 (fecha de nacimiento en el JSON), B-16 (Q2-17 C y Q2-18 A frente a DESIGN y U1) y B-19 (lista permitida).
