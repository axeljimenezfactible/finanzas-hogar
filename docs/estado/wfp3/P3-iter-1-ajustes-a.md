# P3-iter-1 · ajustes-a · revisión independiente (vuelta 1 de 3)

- Fecha: 12-oct-2026. Revisor: `revisor-ajustes-a-r1` (no corrige; solo revisa).
- Alcance: los 16 artboards de `docs/design/maquetas/lienzo/project/ajustes/` de la parte a (P-26 índice y secciones 01 a 07, móvil `-m` y escritorio `-e`) y el fragmento `docs/design/maquetas/lienzo/fragmentos/ajustes-a.json`.
- Contra: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06, U5-07, lo que se ve en el marcado), DESIGN.md (§7.9, §9.3, §10.5), U1 §2.1, §2.2 y §2.3, U4 §0 y P-26.01 a P-26.07, y `referencia-tipo/format.md`.
- **Veredicto: NO aprobado.** 0 altos, 8 medios, 7 bajos. Forma, tokens, contraste y cifras sintéticas están bien. Fallan los enlaces, los objetivos táctiles en móvil, los estados de U1 §2.2 y U4 §0.3, el formato de fecha y algunos elementos de U4 y de DESIGN §7.9.

## Verificación hecha (scripts propios, `python3 -I`, en `datos/tmp/rev-aj-a/`, borrados al final)

| Comprobación | Resultado |
|---|---|
| (a) Forma: doctype, `lang="es-MX"`, `<meta charset>`, `./support.js`, `../ds/finanzas-hogar/tokens.css`, `<x-dc>`, `<helmet>` dentro de `x-dc`, `data-dc-script` con `class Component extends DCLogic` y `renderVals()` | 16/16 bien |
| `data-props` parsea como JSON y `$preview` = `w`×`h` del fragmento | 16/16 bien |
| `html.parser`: etiquetas sin cerrar o mal anidadas, ids repetidos, `<button>`/`<input>` dentro de `<a>` | 0 / 0 / 0 |
| Nombres `[a-z0-9.-]`, stems únicos, fragmento = archivos en disco (16 `boards`, 16 `order`) | bien |
| Sin `iframe`, `img`, `innerHTML`, `http(s)://`, `fetch`, `{{` | 0 |
| (b) Hex literales fuera de `data-props` | 0 |
| `var(--x)`: 68 nombres distintos; todos existen en `tokens.json` (`--font-sans` sale de `type.families.sans`, ver la nota de `meta`) | bien |
| (c) Contraste, tokens de `tokens.json` resueltos en claro y oscuro (cascada del `<helmet>` hecha a mano, sin navegador) | 69 pares (30 de texto, 39 de UI) × 2 temas = 138 combinaciones; **0 bajo umbral** |
| (d) P3-12: grep de las 21 cifras canon con y sin coma, más sus centavos; grep de credenciales | 0 / 0 / 0 |
| Dinero y fechas contra `datos-sinteticos.json` | 59 montos distintos: 55 están en el JSON; 4 son ejemplos tecleados o sumas declaradas por el autor ($14,950.00, $27,786.00, $29,040.00, $29,210.00). 25 fechas: 24 están en el JSON y 1 se calcula (14-ene-2027) |
| Muestra de 10 montos encontrados | $15,981.90 · $1,566.00 · $2,600.00 · $6,104.85 · $1,098.60 · $1,100.00 · $349.00 · $1,452.91 · $2,500.00 · $399.00 |
| Fechas que coinciden con el oráculo (comando de P3-12) | En pantalla solo aparece 30-oct-2026, con $1,861.58 sintético (próximo pago de Kueski Pay), que no es cifra del reporte ni de Kueski: sin hallazgo |
| (f) Emoji, guion como signo de negativo, fechas ISO visibles, inglés | 0, 0, 0; inglés: solo las claves `groceries` y `eating_out` (ver la nota en P3-08) |
| SVG con `aria-hidden`, botones sin nombre accesible, campos sin `label` o `aria-label` | 0 sin `aria-hidden`; 0 sin nombre; 0 sin etiqueta |
| (g) Raíz móvil `.m` con `width: var(--breakpoint-phone)` (412 px), `min-height: 915px`, `box-sizing: border-box`; escritorio `.pg` flex-wrap sin ancho en px, `.mw` `flex: 999 1 560px`, `expand: "fill"` solo en `-e` | bien |
| (h) Enlaces | 2 destinos rotos: `ajustes-p26-s08-reportes-frecuencias-{m,e}` en 9 archivos (10 enlaces) y `#instrucciones` en 2 (ver M-01 y M-02) |

## Líneas de evidencia

- **P3-07**: revisión hecha por artboard (tabla de abajo). El contraste se calculó con los tokens resueltos en los dos temas: 0 pares bajo 4.5:1 en texto y bajo 3:1 en UI. Los mínimos son 3.46 en claro (borde de control sobre tarjeta) y 4.05 en oscuro (punto de persona Axel sobre la página, UI). Los hallazgos van en la sección «Hallazgos».
- **P3-08**: 14 secciones de Ajustes ✓ (índice y columna interna); español de México ✓; dinero `$1,234.56` con `tabular-nums` en `body` ✓ y signo U+2212 (−$17.70, −$212.40) ✓; paleta solo por tokens C2 ✓; Notion como referencia ✓. **Barra lateral colapsable ✗**: no hay botón `panel-left` ni las demás piezas de DESIGN §7.9 (M-07). **Fechas ✗**: 17 llevan el día con cero a la izquierda (M-06). Nota: las claves `groceries` y `eating_out` se ven porque U4 P-26.04 las pide («Clave que usan los cálculos») y U5-07 no las tiene en su lista permitida. Es un choque entre especificaciones, no un error del autor: queda para el orquestador.
- **P3-12**: 0 cifras canon (con y sin coma, y en centavos), 0 credenciales. Montos y fechas salen de `datos-sinteticos.json` o se calculan a partir de él. La revisión manual de fechas encuentra solo 30-oct-2026, con un monto sintético. Cumple.
- **U5-04** (lo que se puede ver en el marcado): 138 combinaciones, 0 bajo umbral. Por tema: claro, 69 pares y 0 bajo umbral; oscuro, 69 pares y 0 bajo umbral. La medición con estilos computados en el DOM la hace U5 en su ronda.
- **U5-06**: tokens de rol ✓; tipografía con las clases `title-1`, `title-3`, `body`, `small` y `caption` de `tokens.css` ✓; cifras tabulares ✓; Lucide con trazo 1.5 y 20 px ✓; espaciado y radios por token ✓; foco visible (`:focus-visible` y `.inw:focus-within`) ✓. **Barra lateral colapsable ✗** (M-07). **Objetivos táctiles ≥ 44 px ✗**: 5 controles (M-03).
- **U5-07**: sin inglés, salvo las claves que pide U4 ✓; dinero ✓; sin fechas ISO visibles ✓. **Fechas ✗**: «01-oct-2026», «02-nov-2026», etc. U4 y el resto del lienzo escriben «9-nov-2026» y «1-nov-2026» (M-06). Hay errores de concordancia (B-02).

## Fidelidad por pantalla (U1 §2.1/§2.2/§2.3 y U4)

| Artboard | Presente | Falta o difiere |
|---|---|---|
| índice m/e | Las 14 secciones en el orden de U1 §2.3; columna interna (e) y lista (m); tarjetas con un resumen del JSON; «Estados» (e) | Título «Ajustes · Ajustes» (B-01); vacío inventado «Sin secciones que mostrar»; enlace a s08 roto (M-01) |
| s01 Cuenta y roles | Usuarios (2, solo lectura), fila de Karen, contraseña (3 campos y «Al menos 12 caracteres»), tabla de sesiones con «Esta sesión», cierre con alcance y diálogo crítico (Cancelar primero), Este dispositivo, 4 parámetros con «Valor inicial» y «Restablecer», vista previa de la sesión y de los intentos | Vacío «Sin sesiones que mostrar» contra U4 («no ocurre») y error al guardar ausente (M-04); `#instrucciones` roto y de menos de 44 px (M-02, M-03); la primaria de U1 «Cerrar otras sesiones» va en estilo secundario y «Cambiar contraseña» en primario (B-06) |
| s02 Presupuesto | Interruptor de proyección arriba (apagado) con vista previa, tabla por categoría con el total, diálogo de alcance con «De este mes en adelante» preseleccionado y el total antes → después, gastos fijos (8), vista previa de desactivar, despensa programada | Acción por fila «Quitar el monto de {categoría}» (M-08); vacío «Sin gastos fijos · Nuevo gasto fijo»; error al guardar (M-04); 3 botones primarios (B-06) |
| s03 Ingresos | Persona, frecuencia, montos del 15 y del 30, «Rige desde», fila de Ximena, enlace a P-08, vista previa | Error al guardar (M-04); fecha «01-oct-2026» (M-06) |
| s04 Categorías | Tabla con Nombre, Clave, Padre, Tipo, Reparto (lista), Activa y Gastos; selector Tabla/Tablero; vista previa del reparto; aviso al archivar Despensa | «Activa» como texto y no como interruptor (M-08); el asa de reordenar es un icono decorativo y no un control (B-04); error al guardar (M-04) |
| s05 Créditos y condiciones | 4 líneas en tabla, condiciones de la revolvente, A/B lado a lado con B vigente y «Usar la lectura A», «Pendiente confirmar» con «Ya lo confirmé», condiciones de tarjeta y préstamo, tasas por plazo con vista previa, promociones, restricciones con «Faltan restricciones», 4 pesos con «Pesos iguales», tope en P-26.12, tiempo máximo (1 a 30) y ventana ± días | Vista **calendario** de promociones, que es la vista por defecto (U1 §2.3 «calendario*», U4): solo hay tabla (M-08); «Eliminar promoción» (M-08); bloques plegables en móvil (M-08); vacíos «Sin tasas» y «Sin promociones» y error al guardar (M-04); casilla y enlace de menos de 44 px en móvil (M-03) |
| s06 Tiendas y precios | Tabla de tiendas (Tabla por defecto, con Galería, Tablero y Calendario), visitas a mercados, membresía de Sam's con punto de equilibrio en −$, lectura asistida por cadena con V-19, las 6 fuentes de precio, 3 parámetros con «Restablecer», vista previa de frescura | Vacío «Sin cortes importados · Importar corte PROFECO» y error al guardar (M-04); enlace de menos de 44 px en móvil (M-03); propiedades por tienda (PROFECO, sitio, dirección, notas), aceptable en el listado |
| s07 Recetas | Recetario con nombre editable, tiempos, origen, uso, Duplicar y Eliminar; porciones y tiempos; recetas por origen; 3 bases con «Recargar»; restricciones vacías; IA apagada | Error al guardar (M-04); enlace de menos de 44 px en móvil (M-03); la leyenda «Enter o salir del campo guarda» bajo vistas previas que no son campos (B-03) |

## Por artboard: reglas cumplidas / falladas / correcciones

Común a los 16: cumplen la forma de format.md, P3-12, U5-04 (estático), el uso de tokens, `lang` y el título propio, los iconos con `aria-hidden` y los campos con etiqueta.

| Artboard | Falladas (ID) | Corrección para la vuelta 2 |
|---|---|---|
| indice-e | M-01, M-04 (vacío inventado), M-07, B-01 | Slug `s08-reportes`; quitar el vacío «Sin secciones» (U1 §2.2: «Según sección») o ponerlo en «no aplica»; barra lateral de §7.9; título «Ajustes» |
| indice-m | M-01, B-01 | Slug `s08-reportes`; título |
| s01-e | M-01, M-02, M-04, M-05, M-06, M-07, B-02, B-06, B-07 | Ver los hallazgos |
| s01-m | M-02, M-03, M-05, M-06, B-06 | `#instrucciones` con destino real y de 44 px; fechas sin cero |
| s02-e | M-01, M-04, M-05, M-07, M-08, B-06 | «Quitar monto» por fila; vacío de gastos fijos; error al guardar; un solo primario |
| s02-m | M-05, M-08, B-06 | Igual, con la primaria en el encabezado |
| s03-e | M-01, M-04, M-05, M-06, M-07 | Error al guardar; «1-oct-2026» |
| s03-m | M-05, M-06 | Igual |
| s04-e | M-01, M-04, M-05, M-07, M-08, B-02, B-04 | Interruptor «Activa»; asa como botón con `aria-label` |
| s04-m | M-05, M-08, B-04 | Igual |
| s05-e | M-01, M-04, M-05, M-06, M-07, M-08, B-02, B-06 | Vista calendario por defecto, «Eliminar promoción», vacíos de tasas y promociones |
| s05-m | M-03, M-05, M-06, M-08, B-05 | Bloques plegables; casilla de 44 px; enlace de 44 px |
| s06-e | M-01, M-04, M-05, M-06, M-07, B-02 | Vacío de PROFECO; error al guardar |
| s06-m | M-03, M-05, M-06 | Enlace de 44 px |
| s07-e | M-01, M-04, M-05, M-07, B-02, B-03 | Error al guardar; quitar la leyenda «Enter…» de las vistas previas de acción |
| s07-m | M-03, M-05, B-03 | Enlace de 44 px |
| ajustes-a.json | B-05 | Medir `h` al publicar |

## Hallazgos

### Medios (8)

- **M-01 · Enlace roto a la sección 08** (criterio h, P3-07). Hay 10 `href` a `ajustes-p26-s08-reportes-frecuencias-{e,m}.dc.html`, archivo que no existe. Están en la columna interna de los 7 `-e`, en las 2 tarjetas de `indice-e` y en 1 de `indice-m`. El archivo que publicó la parte b, y que lista `ajustes-b.json`, es `ajustes-p26-s08-reportes-{e,m}.dc.html`. Corrección: cambiar el slug en los 9 archivos.
- **M-02 · Ancla sin destino**: `href="#instrucciones"` en s01-m y s01-e, sin ningún `id="instrucciones"`. Corrección: añadir el bloque de instrucciones con ese `id` o usar `href="#"`.
- **M-03 · Objetivos táctiles de menos de 44 px en móvil** (DESIGN §9.3, U4 §0.7, U5-06). Son 5 controles: «Ver instrucciones» (s01-m); la casilla «Ya lo confirmé», de 20 × 20 dentro de `.swr` sin alto mínimo (s05-m); «Ir a Personas y reparto» (s05-m); «Ir a Datos y respaldos» (s06-m y s07-m). Los enlaces sueltos en `.cn` miden lo que mide una línea, unos 24 px. Corrección: darles `min-height: var(--size-touch-min)` y `display: inline-flex; align-items: center`, o la clase `btn t`; a la casilla, una `label` envolvente o una fila `.swr` con `min-height: var(--size-touch-min)`.
- **M-04 · Estados de U1 §2.2 y U4 incompletos o contradichos.** U1 §2.2 define para P-26 el error «El campo que no se guardó vuelve a su valor anterior con aviso» (U4 §0.3, «Error al guardar»), y ningún artboard de la parte a lo dibuja (`No se pudo guardar` aparece 0 veces; la parte b sí lo trae). Además, s01 dibuja el vacío «Sin sesiones que mostrar», que U4 P-26.01 declara imposible («sesiones vacías no ocurre»), e `indice-e` inventa «Sin secciones que mostrar». Faltan los vacíos de U4 «Sin gastos fijos · Nuevo gasto fijo» (s02), «Sin tasas · Nueva tasa» y «Sin promociones · Nueva promoción» (s05) y «Sin cortes importados · Importar corte PROFECO» (s06). Corrección: en la tarjeta «Estados» de cada `-e`, añadir «Error al guardar» con `circle-alert` y el texto de §0.3 («No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar»); en s01, «Vacío: no aplica (siempre está la sesión actual)»; añadir los vacíos que faltan.
- **M-05 · No se dibuja el error de validación de un campo** (U4 §0.3 y §0.6, criterio e «campos y validaciones de U4»). La regla `.inw.er` existe y no se usa en ningún artboard, y `ui-border-error` no aparece en ningún par. De los mensajes del catálogo solo se ve V-19 (s06). Corrección: al menos un campo por sección en estado de error, con borde `ui-border-error`, `circle-alert` y su mensaje literal. Por ejemplo: s01 V-15 («Las dos contraseñas nuevas no coinciden.»), s02 V-03, s03 V-04, s04 V-12, s05 V-23 («Escribe un número entero del 1 al 30.»), s06 V-25, s07 V-09.
- **M-06 · Fechas con cero a la izquierda** (U5-07, P3-08). Hay 17 apariciones en 8 artboards: 01-oct-2026 (×6), 02-oct-2026, 03-oct-2026, 02-nov-2026, 06-nov-2026 y 01-ene-2027 (×2 cada una), en s01, s03, s05 y s06, en m y e. U4 escribe «9-nov-2026» y «8-ene-2027», las especificaciones no tienen ni un ejemplo con cero, y en el resto del lienzo hay 50 fechas sin cero. Corrección: escribir el día sin relleno (1-oct-2026, 6-nov-2026, 1-ene-2027).
- **M-07 · Barra lateral sin colapso ni anatomía de DESIGN §7.9** (P3-08 «barra lateral colapsable», U5-06). En los 8 `-e`, `.sb` trae marca, Buscar, Registrar gasto y las 9 secciones, pero no el botón `panel-left` («‹‹», `Ctrl+\`), ni el avatar con su menú (Cuenta, Apariencia, Cerrar sesión), ni el interruptor de tema al pie. La línea de atajos tampoco incluye `Ctrl+\`. Corrección: añadir un `<button aria-label="Contraer barra lateral">` con `panel-left` en el encabezado de `.sb`, el avatar y el interruptor de tema al pie, igual que `shell/shell-barra-lateral-e`.
- **M-08 · Elementos secundarios de U4 ausentes.** (1) s05-e: la vista **Calendario** de promociones es la vista por defecto (U1 §2.3 «calendario*» y U4 «calendario por defecto»), y el artboard muestra Tabla presionada sin ningún calendario dibujado. (2) s05: falta la acción crítica «Eliminar promoción». (3) s02: falta la acción por fila crítica «Quitar el monto de {categoría}». (4) s04: «Activa» se muestra como texto y no como interruptor (U4: «Interruptor; apagar = archivar»). (5) s05-m: U4 pide «Bloques plegables en este orden» y no hay ningún control plegable (0 `aria-expanded`); el artboard mide 9,610 px. Corrección: dibujar el calendario como vista activa en s05-e, con franjas de inicio a fin y «Vigente hoy», y dejar la tabla como alternativa; añadir las dos acciones con su diálogo crítico (§0.4); poner `.sw role="switch"` en «Activa»; usar encabezados de bloque `<button aria-expanded>` de 44 px en s05-m.

### Bajos (7)

- **B-01**: el título del índice es «Ajustes · Ajustes». Debe ser «Ajustes».
- **B-02**: concordancia en «No se pudo cargar las sesiones / las categorías / los créditos / las tiendas / las recetas». Debe ser «No se pudieron cargar …».
- **B-03**: la leyenda «Enter o salir del campo guarda · Esc descarta» aparece bajo vistas previas que no son de un campo: eliminar receta y recargar base (s07), interruptor de proyección (s02). Esas confirman con un diálogo o al elegir.
- **B-04**: el asa de reordenar de s04 es un `svg` con `aria-hidden` dentro de `.mt` y no un control real. Debe ser un botón con `aria-label="Mover {categoría}"`, de 44 px en móvil.
- **B-05**: las alturas del fragmento son estimadas y, según el autor, «a propósito altas» (s05-m 9,610, s02-m 7,060). Puede quedar una cola vacía en el lienzo. Hay que medirlas al publicar.
- **B-06**: la acción primaria no es única ni está siempre visible (U1 §2.2). s02-e tiene 3 `btn p`, s05 tiene 2 («Nuevo crédito» y «Usar la lectura A»). En s01, «Cambiar contraseña» va en primario y la primaria de U1, «Cerrar otras sesiones», en secundario. En los `-m` la primaria no va en el encabezado `.tbar`.
- **B-07**: el disparador «Cerrar sesiones» (s01) usa el relleno crítico fuera del diálogo. U4 §0.4 reserva ese énfasis para el botón que confirma dentro del diálogo.

## Supuestos del revisor

1. El contraste se calculó con una cascada propia sobre el CSS del `<helmet>` y los tokens de `tokens.json` (claro y oscuro), sin navegador (K45). La medición con estilos computados en el DOM queda para U5-04.
2. `--font-sans` cuenta como token de `tokens.json`: `type.families.sans`, según `meta.note` («Nombres = variables de tokens.css sin “--”»).
3. Doy por buenos los ejemplos tecleados y las sumas que el autor declara (supuestos 2 a 7 de su reporte), porque se derivan del JSON y no son cifras del reporte.
4. `is_interactive: false` sigue la orden del orquestador, aunque format.md lo pide en `true` para artboards con navegación.
5. Según P3-06, la copia del informe en `docs/design/salida/iter-1/` la hace el orquestador. Este revisor no escribe fuera de `docs/estado/wfp3/`.
