# P3-iter-1 · ajustes-b · revisión de la vuelta 1

Fecha: 12-oct-2026 · Revisor: `revisor-ajustes-b-r1` (independiente, vuelta 1 de 3; no corrige) · Criterio: P3-07 de `docs/criterios/P3-claude-design.md` 1.3.
Alcance: 14 artboards de `docs/design/maquetas/lienzo/project/ajustes/` (P-26.08 a P-26.14, `-m` y `-e`) y el fragmento `docs/design/maquetas/lienzo/fragmentos/ajustes-b.json`.
Contra: P3-08 y P3-12; U5-04, U5-06 y U5-07 (lo verificable en el markup); DESIGN.md (§3.2, §7.3, §7.9, §7.13, §7.14, §9.3, §10.5); U1 §1.3, §2.1, §2.2 y §2.3; U4 §0 y P-26.08 a P-26.14; `referencia-tipo/format.md`.
Ruta del informe: la que fijó la tarea (`docs/estado/wfp3/`); P3-07 nombra `docs/estado/P3-iter-N-<area>.md`. La copia a `docs/design/salida/iter-1/` (P3-06) la hace el orquestador.

**Veredicto: NO APROBADO** · 1 alta · 11 medias · 10 bajas.

## Verificación ejecutada (scripts propios, `python3 -I`, en `datos/tmp/rev-ajb-r1/`, borrados al final)

| Comprobación | Resultado |
|---|---|
| (a) Forma: doctype, `lang="es-MX"`, charset, `<title>` propio y único, `<script src="./support.js">` (1 vez), `<link … ../ds/finanzas-hogar/tokens.css>` (1 vez, en `<head>`), `<x-dc>` → `<helmet><style>` primero, bloque `data-dc-script` fuera de `x-dc` con `class Component extends DCLogic` y `renderVals()`, sin import/export | 14/14 OK |
| (a) `data-props` parsea (tras decodificar entidades) y `$preview` = `w`×`h` del fragmento | 14/14 OK |
| (a) `html.parser`: etiquetas sin cerrar, cierres huérfanos, ids repetidos | 0 · 0 · 0 |
| (a) Nombres válidos (`[A-Za-z0-9.-]`), stems únicos en todo el lienzo (89 artboards) | OK · 0 duplicados |
| (b) Hex literales fuera de `data-props` (excluidos `href="#…"`); `rgb()`/`hsl()`/`oklch()` | 0 · 0 |
| (b) `var(--x)` contra `tokens.json` | todos existen (`--font-sans` sale de `type.families.sans`); 0 `var()` con respaldo; 0 propiedades locales |
| (b) Clases usadas sin regla CSS | **`.in` sin regla en s11-m, s11-e, s12-m, s12-e** (H-02) |
| (c) Contraste de texto (DOM con herencia de `color` y fondo efectivo, tokens resueltos en claro y oscuro) | 23 combinaciones, ~1,362 nodos de texto por tema: **0 bajo 4.5:1** |
| (c) Contraste de UI 3:1 | bordes de control y pista OK; **perilla del interruptor bajo 3:1 en oscuro** (H-01) |
| (d) P3-12: 21 patrones de cifras canon (con y sin coma), en centavos, credenciales, montos ≥ 1,000 sin separador | 0 · 0 · 0 · 0 |
| (d) Fechas que coinciden con el calendario real | intersección de 15 fechas quincenales; en ajustes-b solo «Vie 30 oct 2026» (s09-e, próximo pago) y sin monto del reporte ni de Kueski al lado: OK |
| (d) Muestra de 10+ valores contra `datos-sinteticos.json` | todos presentes (tabla abajo) |
| (f) `input`/`select` con `label for` o `aria-label`; `button` con nombre y `type`; `a` con nombre; `svg` decorativos con `aria-hidden` | OK en los 14 |
| (f) Objetivos ≥ 44 px en móvil (alturas de CSS e inline) | **fallan los campos `.in` de s11-m (4) y s12-m (3)** (H-02); el resto ≥ 44 |
| (f) Emoji · barra de estado o teclado falsos · `innerHTML`/`iframe`/red/`localStorage` | 0 · 0 · 0 |
| (f) Inglés visible | **claves técnicas en `<code>`** en 6 artboards (H-03) |
| (f) Dinero `$1,234.56`, signo «−» U+2212 | 0 mal formados; 0 «-$»; «−» solo en «− 400 kcal» (s13) |
| (g) Móvil: raíz única `.mr` `width:412px; min-height:915px`, `box-sizing:border-box` | 7/7 |
| (g) Escritorio: raíz `.pg` fluida (flex-wrap, sin ancho en px), barra 240 px + `.ct{flex:999 1 520px;min-width:0}`, tablas en caja `overflow-x:auto`; fragmento con `"expand":"fill"` | 7/7 |
| (h) Enlaces: 193 `href`; destinos dentro de la lista del lienzo, mismo ancho; resto `#` | 0 rotos · 0 de otro ancho |
| Fragmento: JSON válido, `boards` = `order` = los 14 archivos, `page:"ajustes"`, `is_interactive:false`, `w` 412/1440, móvil `h` ≥ 915 | OK |

### Recálculo independiente de contraste (P3-07: 10 pares de texto + UI)

| # | Par (texto / fondo) | Claro | Oscuro | Umbral |
|---|---|---|---|---|
| 1 | `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 ✓ |
| 2 | `ui-text-secondary` / `ui-surface-card` | 5.04 | 7.97 | 4.5 ✓ |
| 3 | `ui-text-default` / `ui-surface-field` | 11.53 | 14.07 | 4.5 ✓ |
| 4 | `ui-text-secondary` / `ui-surface-selected` (vista previa) | 5.02 | 8.08 | 4.5 ✓ |
| 5 | `ui-role-success` / `ui-role-success-surface` | 4.81 | 6.19 | 4.5 ✓ |
| 6 | `ui-role-warning` / `ui-role-warning-surface` | 5.12 | 6.09 | 4.5 ✓ |
| 7 | `ui-text-on-primary` / `ui-action-primary-bg` | 5.77 | 10.47 | 4.5 ✓ |
| 8 | `ui-role-error` / `ui-role-error-surface` | 5.37 | 6.05 | 4.5 ✓ |
| 9 | `ui-text-on-critical` / `ui-action-critical-bg` | 5.60 | 10.55 | 4.5 ✓ |
| 10 | `ui-text-link` / `ui-surface-card` | 5.16 | 7.94 | 4.5 ✓ |
| UI | `ui-border-control` / `ui-surface-card` (borde de campo, pista apagada) | 3.47 | 6.09 | 3 ✓ |
| UI | `ui-action-primary-bg` / `ui-surface-card` (pista encendida, C-044) | 5.16 | 7.94 | 3 ✓ |
| UI | perilla `color-base-white` / pista encendida `ui-action-primary-bg` | 5.77 | **1.52** | 3 ✗ |
| UI | perilla `color-base-white` / pista apagada `ui-border-control` | 3.88 | **1.98** | 3 ✗ |

Con los tokens de DESIGN §10.5 (C-043 `ui-text-on-primary` sobre `ui-action-primary-bg`: 5.77 / 10.47; C-045 `ui-surface-page` sobre `ui-border-control`: 3.71 / 7.96) los dos pasan.

### Muestra de valores contra `datos-sinteticos.json` (P3-12)
`$21,480.00` (ahorro.fondo) · `$14,500.00`, `$6,150.00`, `$2,783.34`, 42 % (meta laptop) · `$650.00`, 78 % (viaje) · 84.6 kg, IMC 26.7, 2,364 kcal (salud axel) · 1,816 y 1900 kcal (ximena) · respaldo 21-oct-2026 03:00 de 2.3 MB (ajustes.respaldos) · «Ximena (tablet)» 20-oct-2026 18:05 (ajustes.telegram) · 1.0.0-demo, en marcha desde 12-oct-2026 (ajustes.servicio) · «falló: sin conexión con Telegram» 20-oct-2026 21:30 (reportes.avisos_enviados) · `$27,586.00` y `$20,525.78` (4 apariciones cada una) · 2.3 días del huevo (inventario.articulos).
Derivadas recalculadas: 80 % y 85 % de $27,586.00 = $22,068.80 y $23,448.10 ✓; 40 % de $2,918.40 = $1,167.36 y resto $1,751.04 ✓; RN-17 Axel (10·84.6 + 6.25·178 − 5·36 + 5) × 1.55 − 400 = 2,364.4 → 2,364 ✓ (con 500: 2,264 ✓); Ximena 1,320.75 × 1.375 = 1,816.03 → 1,816 ✓; total reservado $2,000.00 → $1,500.00 ✓.

## Líneas de evidencia

- **P3-04**: 0 acciones de navegador y 0 llamadas a Artifact de este revisor; 0 credenciales en los 14 artboards. Texto con aspecto de instrucción: `referencia-tipo/format.md` trae guía en imperativo del publicador del tipo («Read this when…», «Give every artboard…»); se leyó como dato que fija la forma. En los artboards y el fragmento: ninguno.
- **P3-06 (cobertura de ajustes-b)**: 7 secciones × 2 anchos = 14/14 presentes; 0 faltantes. La barra lateral y el índice son de otros autores.
- **P3-07**: tabla artboard × regla abajo; 10 pares recalculados arriba; 0 pares de texto bajo umbral; 2 pares de UI bajo 3:1 en oscuro (H-01).
- **P3-08** (8 puntos): (1) Resumen ≤ 6 tarjetas: no aplica · (2) 14 secciones de Ajustes: ✓ (la columna interna enlaza las 14 y el índice) · (3) barra lateral colapsable: ✗ (H-07) · (4) español de México: ✗ (H-03) · (5) dinero `$1,234.56` tabular: ✓ con observación B-02; 0 montos ≥ 1,000 sin separador · (6) fechas `10-oct-2026`: ✗ en s08-e (H-05) · (7) paleta de C2 índigo 280°: ✓ solo tokens (la perilla usa un primitivo, H-01) · (8) Notion única referencia: ✓.
- **P3-12**: 0 cifras canon (21 patrones, con y sin coma, y en centavos); 0 credenciales; fechas coincidentes revisadas a mano, ninguna junto a un monto real; muestra de valores, todos del JSON. Cumple.
- **U5-04**: 0 pares de texto bajo umbral; UI: perilla del interruptor bajo 3:1 en oscuro (H-01). Falla.
- **U5-06**: tokens ✓, iconos Lucide 20 px con trazo de 1.5 ✓, radios ✓, foco visible (`:focus-visible` con `ui-border-focus`) ✓; fallan barra lateral colapsable (H-07), objetivos táctiles (H-02), tipografía (B-01).
- **U5-07**: dinero ✓; fallan el texto en inglés (H-03) y la fecha en formato del navegador (H-05).

## Tabla artboard × regla

Columnas: F forma y fragmento · T tokens · C contraste (U5-04) · D cifras sintéticas (P3-12) · U fidelidad a U1/U4 · A accesibilidad y táctil · I idioma y formatos (U5-07, P3-08) · N navegación (barra lateral/inferior, enlaces) · E estados (U1 §2.2).

| Artboard | F | T | C | D | U | A | I | N | E | Fallas |
|---|---|---|---|---|---|---|---|---|---|---|
| s08-reportes-m | ✓ | ✓ | ✗ | ✓ | ~ | ✓ | ✓ | ✗ | — | H-01, H-06, B-06 |
| s08-reportes-e | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✗ | ✗ | ✓ | H-01, H-03, H-05, H-07 |
| s09-notificaciones-m | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | ✗ | — | H-01, H-06 |
| s09-notificaciones-e | ✓ | ✓ | ✗ | ✓ | ✗ | ✓ | ✗ | ✗ | ✓ | H-01, H-03, H-04, H-07, B-02, B-03, B-04 |
| s10-apariencia-m | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ | ✓ | ✗ | — | H-01, H-06, B-05 |
| s10-apariencia-e | ✓ | ✓ | ✗ | ✓ | ✗ | ✓ | ✓ | ✗ | ✓ | H-01, H-07, H-08, B-05 |
| s11-datos-respaldos-m | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | — | H-01, H-02, H-03, H-06, H-09, H-10 |
| s11-datos-respaldos-e | ✓ | ✗ | ✗ | ✓ | ✗ | ~ | ✗ | ✗ | ✓ | H-01, H-02, H-03, H-07, H-09, H-10 |
| s12-personas-reparto-m | ✓ | ✗ | ✓ | ✓ | ✓ | ✗ | ✓ | ✗ | — | H-02, H-06, B-02, B-08 |
| s12-personas-reparto-e | ✓ | ✗ | ✓ | ✓ | ~ | ~ | ✓ | ✗ | ✓ | H-02, H-07, B-02, B-10 |
| s13-salud-metas-m | ✓ | ✓ | ✗ | ✓ | ✗ | ✓ | ✗ | ✗ | — | H-01, H-03, H-06, H-10, H-12, B-03 |
| s13-salud-metas-e | ✓ | ✓ | ✗ | ✓ | ✗ | ✓ | ✗ | ✗ | ✓ | H-01, H-03, H-07, H-10, H-12, B-03 |
| s14-inventario-m | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✗ | — | H-06, H-11 |
| s14-inventario-e | ✓ | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✗ | ✓ | H-07, H-10, H-11, B-09 |

«—» en E: el móvil no lleva fila de estados (la tarea la pide en escritorio). B-01, B-07 y B-09 son transversales.

## Fidelidad: pantalla → elementos de U1/U4

| Sección | Presentes | Faltantes o distintos |
|---|---|---|
| P-26.08 Reportes | tabla con Reporte · Activo · Día · Hora · Destinatarios · Próximo envío · Último resultado; «Enviar ahora a…»; 5 parámetros `reports.*` con «Valor inicial»; vista previa del próximo envío; regla de días 29–31; Bajo demanda (mensual ahora con mes y casilla Telegram, estado de cuenta de Karen, vista web); estados | fechas de Karen con `type="date"` (H-05); móvil sin «Enviar ahora a…» y con acción «Reintentar el envío» que U4 no tiene (B-06) |
| P-26.09 Notificaciones | 4 alertas con activa, destinatarios, anticipación o umbral, día de revisión, hora; chats detectados (nombre, fecha, Confirmar como, Descartar con diálogo crítico, foco en «Cancelar»); chats confirmados (registrado, etiqueta, recibe avisos); aviso de prueba; vista previa de umbrales y del próximo pago; estados | «Valor inicial» lejos del campo vigente (H-04); terminación del chat sin los 4 caracteres (no está en el JSON, B-09) |
| P-26.10 Apariencia | tema segmentado (Sistema · Claro · Oscuro, `role="radio"`), atajos de una tecla, valores fijos, muestra en claro y oscuro (escritorio), nota de estados | muestra clara sin `data-theme="light"` (H-08) |
| P-26.11 Datos y respaldos | encabezado último/próximo/restauración probada; tabla de 5 respaldos con tipo, resultado, tamaño, integridad; «Restaurar» con doble confirmación (paso 2, `RESTAURAR`, botón crítico con `aria-disabled`); datos semilla (3); IA opcional (2, apagadas, llave); estado del servicio; 3 parámetros con vista previa; estados | tabla «Lo que se creará» (H-09); hora del respaldo diario y huella (H-10); campos sin estilo (H-02) |
| P-26.12 Personas y reparto | Axel/Ximena/Karen con papel, cuenta, paga directo, qué hace, tope (vacío), notas; 3 repartos con resto calculado y «Reparto por defecto»; V-16; vista previa con $2,918.40; diálogo de alcance (Q2-68 C) con «Vista previa pendiente»; barra etiquetada en móvil; estados | campos sin estilo (H-02); parte sin control de edición visible (B-10) |
| P-26.13 Salud y metas | fichas de Axel y Ximena (estatura, nacimiento, sexo, actividad, objetivo, peso, IMC y clase, meta sugerida con desglose, confirmada, «Usar la sugerida», mis ajustes, restricciones); ajuste de la casa (2 parámetros); metas con selector de vista, avance, necesaria, estado, «Activa» con vista previa (3); «Nueva meta»; estados | notas del perfil, inicio y notas de la meta (H-10); fecha de nacimiento vacía con sugerencia calculada (H-12) |
| P-26.14 Inventario | 14 artículos con modo, unidad, nivel, ritmo, caducidad, existencia, estado con icono y texto, días; resumen 1/5/8; Tabla · Tablero por modo; «Agregar producto al inventario»; vista previa de estado; estados | días con decimal y «Reponer» sin nivel (H-11); notas (H-10); filtros de la vista no dibujados (baja, sin id) |

## Hallazgos

### Alta
- **H-01 · U5-04, P3-07, DESIGN §7.3 y §10.5 (C-043, C-045).** La perilla del interruptor (`.sw i:after`) usa `var(--color-base-white)` en los dos temas. En oscuro queda a 1.52:1 sobre la pista encendida (`ui-action-primary-bg` = brand-300) y a 1.98:1 sobre la apagada (`ui-border-control` = neutral-400): la posición, que según DESIGN dice el estado, no se distingue. 10 artboards: s08, s09, s10, s11 y s13, `-m` y `-e`. **Corrección:** `.sw.on i:after{background:var(--ui-text-on-primary)}` y `.sw i:after{background:var(--ui-surface-page)}` (C-043: 5.77 / 10.47; C-045: 3.71 / 7.96); de paso, la pista de 36 × 20 px dentro de 44 px de §7.3.

### Media
- **H-02 · U5-06, DESIGN §7.3 y §9.3.** La regla `.in` desapareció con la poda de CSS en s11-m, s11-e, s12-m y s12-e, pero los campos siguen con `class="in"`: archivo CSV, 3 parámetros y `RESTAURAR` (s11) y tope por quincena, nombre visible, notas y nombre del reparto (s12) se pintan con el estilo del navegador. Quedan de unos 21 px de alto (menos de 44 px en móvil), sin borde `ui-border-control`, sin fondo `ui-surface-field` y sin tema oscuro. **Corrección:** restituir `.in,.sl{min-height:var(--size-control-coarse);padding:0 var(--space-3);border:var(--border-width-default) solid var(--ui-border-control);border-radius:var(--radius-sm);background:var(--ui-surface-field);color:var(--ui-text-default);font:inherit;width:100%;min-width:0}` como en s08.
- **H-03 · U5-07, P3-08 (español de México).** Hay claves técnicas visibles en `<code>`: `reports.daily_send_time`, `alerts.budget_warning_pct`, `attachment.max_size_mb`, `idempotency.retention_hours`, `ocr.timeout_seconds`, `health.goal_adjust_lose_kcal` y otras. Aparecen en s08-e, s09-e, s11-m, s11-e, s13-m y s13-e. Es texto en inglés y jerga que U4 usa solo como documentación. **Corrección:** quitar los `<code>` y dejar el nombre en español y «Valor inicial: …».
- **H-04 · U4 §0.5, U1 §2.3 («junto a cada campo se lee “Valor inicial: …”»), DESIGN §7.13.** En s09-e los campos vigentes (80, 100, 1, 09:00, domingo, en la tabla de alertas) no llevan «Valor inicial». La lista «Parámetros operativos» trae «Valor inicial» pero no el campo. **Corrección:** en la fila de cada alerta, «Valor inicial: 80 %» (100 %, [1, 0], 1 día, 09:00, domingo) junto a su campo, como pide U4 para 1440 («parámetros en la fila de su alerta»). La lista se quita o lleva campo y valor inicial en la misma fila.
- **H-05 · U5-07, P3-08 (fechas `10-oct-2026`).** En s08-e, «Generar estado de cuenta de Karen» usa `<input type="date" value="2026-09-01">` y `…-30`. El navegador lo pinta en su formato local (01/09/2026), no en el de la decisión B. **Corrección:** campo de texto con `1-sep-2026` y `30-sep-2026` (y botón de calendario), como el `dd-mmm-aaaa` de s13.
- **H-06 · U1 §1.3 y §1.4 (G-05).** Los 7 móviles no tienen la barra inferior fija (Inicio · Comidas · «+» · Compras · Buscar); los de comidas sí la traen. **Corrección:** agregarla (56 px de alto, destinos de 64 px o más de ancho, «+» → `../shell/shell-p03-registrar-gasto-m.dc.html`), con enlaces a los `-m` de la lista. El mismo faltante se ve en el índice y en las secciones de ajustes-a: conviene que el orquestador lo unifique.
- **H-07 · P3-08 («barra lateral colapsable»), DESIGN §7.9.** La barra lateral de los 7 escritorios no tiene botón para colapsar (`panel-left` «‹‹»). Tampoco tiene las acciones fijas «Buscar» y «Registrar gasto», ni el avatar, ni el tema en el pie, y usa `aria-current="true"` en lugar de `"page"`. **Corrección:** reutilizar la barra de `shell/shell-barra-lateral-e` o agregar esos elementos; «Registrar gasto» → `../shell/shell-p03-registrar-gasto-e.dc.html`.
- **H-08 · U4 P-26.10 («muestra de tarjeta en claro y oscuro»).** En s10-e solo la muestra oscura lleva `data-theme`. Con el lienzo en oscuro, la muestra «Claro» hereda el tema oscuro y las dos se ven iguales. **Corrección:** `data-theme="light"` en el contenedor de la muestra clara.
- **H-09 · U4 P-26.11 (Importar).** Falta la tabla de vista previa «Lo que se creará»: solo hay una frase con «10 registros». Además, «Archivo» es un campo de texto y no un selector de archivo. Afecta a s11-m y s11-e. **Corrección:** tabla con los renglones que creará el CSV, tomados de `presupuesto.gastos_fijos_recurrentes` e `ingresos.esperado` del JSON, más el botón «Elegir archivo». En móvil, como lista de filas apiladas.
- **H-10 · U4 (campos de sección).** Faltan campos de U4: P-26.13 «Notas» del perfil, e «Inicio» y «Notas» de la meta (Inicio es obligatorio); P-26.14 «Notas»; P-26.11 «Hora del respaldo diario» (PROVISIONAL: va con «Sin valor inicial») y la huella del respaldo. **Corrección:** agregarlos como columnas o campos. En móvil, dentro de la hoja o con «Mostrar N más». Sin dato en el JSON van vacíos o con «—».
- **H-11 · U4 P-26.14, Q2-70 (a) y K7.** «Días restantes» sale con un decimal (2.3, 3.0, 3.5, 20.2, 16.8, 18.7, 17.5), pero Q2-70 (a) los define como piso entero (2, 3, 3, 20, 16, 18, 17). Además, huevo, leche y tortilla dicen «Reponer» con nivel «Sin aviso», lo que contradice la regla K7 escrita en la misma pantalla. El origen es `inventario.articulos` del JSON, que el autor copió tal cual. **Corrección:** pedir a correccion-datos días enteros y un nivel de reposición (o un estado conforme a K7) en los artículos de ritmo. Mientras tanto, pintar el piso entero, que es aritmética sobre el JSON. Lo mismo aplica a `compras-p16-inventario`.
- **H-12 · U4 P-26.13 (RN-17).** «Fecha de nacimiento» está vacía (`dd-mmm-aaaa`), pero la meta sugerida se calcula con la edad. RN-17 dice que, sin ese dato, no hay sugerencia («Falta la fecha de nacimiento»). El JSON trae `edad` y no `birth_date`. **Corrección:** pedir una fecha sintética coherente con `edad` en `datos-sinteticos.json` y pintarla, o mostrar el estado «Falta…» de forma consistente.

### Baja
- **B-01 · DESIGN §3.2, §7.2, §7.3 y §7.10.** El texto base es de 14 px en todo, también en los campos y botones de móvil, donde §7.2 y §7.3 piden 16 px. El título de página de escritorio mide 24 px y `title-1` es de 32 px. El ancho máximo es `1180px` literal y no un token de layout.
- **B-02 · DESIGN §3.3.** Hay dinero en prosa sin `tabular-nums`: en s09-e, «Llega al gastar (plan $27,586.00)» y «…este mes ($27,586.00)…»; en s12-m y s12-e, «Con un gasto de $2,918.40 · …».
- **B-03 · U4 §0.2.** La vista previa no coincide con el campo. En s09-e el campo «Primer aviso» muestra 80 y la vista previa, 80 % → 85 %. En s13-m y s13-e, «Mi ajuste para bajar» muestra 400 y la vista previa, 400 → 500. El campo en edición debe mostrar lo tecleado con el borde de edición, como s08 (20:00) y s11 (15).
- **B-04.** s09-e tiene cuatro campos «Hora», uno por alerta, para una sola clave (`alerts.send_time`). Conviene un campo único o el rótulo «Misma hora para todos los avisos».
- **B-05 · U4 P-26.10 y §0.7.** En s10-m y s10-e el interruptor «Atajos de una tecla» muestra «Enter guarda · Esc descarta», pero U4 pide guardado inmediato. En móvil (Android) no hay atajos.
- **B-06 · U4 P-26.08.** En s08-m aparece «Reintentar el envío», que no está en U4, y falta «Enviar ahora a…». La tarjeta semanal se despliega en la página en vez de abrir una hoja inferior.
- **B-07.** Los textos de ejemplo de los campos («Sin tope», «Sin notas», «Sin etiqueta», «Sin restricciones», «dd-mmm-aaaa») no tienen color de token: el navegador pone el suyo (#757575 en Chrome, unos 4.4:1 en claro). Agregar `::placeholder{color:var(--ui-text-secondary)}`.
- **B-08.** En la barra de reparto de s12-m, los segmentos de Karen y de Axel van pegados y entre sí miden 1.51:1 en oscuro. La información está en texto (K41 se cumple), pero conviene un hueco de 2 px de `ui-surface-card`.
- **B-09.** La columna interna de las 14 secciones va sin los iconos de §5.4 (DESIGN §7.13). La terminación del chat sale «•••» sin los 4 caracteres, que el JSON no trae. En s14-e hay 14 etiquetas «Modo» iguales (mejor «Modo: Frijol negro»). Las alturas del fragmento son estimadas: el orquestador debe medirlas de nuevo.
- **B-10 · U4 P-26.12.** «Parte de cada persona» no tiene control de edición visible, aunque la vista previa muestra una edición (50.00 % → 40.00 %). En escritorio, la vista previa va debajo de la tabla y no a la derecha de la regla.

## Correcciones para la vuelta 2 (orden sugerido)
1. H-01 (perilla) y H-02 (regla `.in`): son cambios de CSS en las plantillas del generador y afectan a 10 y a 4 artboards.
2. H-03, H-04, H-05: texto y formato.
3. H-06, H-07: barra inferior móvil y barra lateral colapsable, en coordinación con shell y ajustes-a.
4. H-08 a H-10: elementos de U4.
5. H-11 y H-12: escalar a correccion-datos (datos-sinteticos.json) y, mientras tanto, pintar valores coherentes con U4.
6. Bajas, si el tope de la vuelta lo permite.
