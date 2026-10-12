# P3-iter-1 · Área «presupuesto» · revisión independiente (vuelta 1 de 3)

Fecha: 12-oct-2026 · Revisor: `revisor-presupuesto-r1` (no corrige) · Alcance: 8 artboards de `docs/design/maquetas/lienzo/project/presupuesto/` (P-06 a P-09, móvil 412 y escritorio 1440) y `docs/design/maquetas/lienzo/fragmentos/presupuesto.json`.
Contra: `docs/criterios/P3-claude-design.md` 1.3 (P3-04, P3-06, P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06, U5-07 en lo verificable en el markup), `docs/design/DESIGN.md` (§7.3, §7.5, §7.6, §7.7, §8, §9.3, §10.5), `docs/design/U1-arquitectura.md` (§1.3, §2.1, §2.2, §2.3, §2.6, §3.8 a §3.11, §6.4), `docs/design/U4-ajustes.md` (P-26.02, vista previa 1 y 2) y `docs/design/maquetas/referencia-tipo/format.md`.

**Veredicto: NO APROBADO** · 1 alta · 4 medias · 9 bajas.

## 1. Tabla artboard × regla

Leyenda: OK = cumple; el id remite a §3.

| Artboard | Forma (format.md) | Tokens y hex | Contraste U5-04 | P3-12 | P3-08 | U5-06 | U5-07 | U1/U4 | Táctil y a11y | Enlaces |
|---|---|---|---|---|---|---|---|---|---|---|
| `presupuesto-p06-mes-m` | OK | OK | OK | OK | OK | B1 | OK | B8 | OK | OK |
| `presupuesto-p06-mes-e` | OK | OK | OK | OK | OK | B1 | OK | B8 | OK | OK |
| `presupuesto-p07-flujo-m` | OK | OK | OK | OK | **A1** | M1 | **A1** | M2, M3, B5 | M3 (36 px) | OK |
| `presupuesto-p07-flujo-e` | OK | OK | OK | OK | **A1** | M1 | **A1** | M2, M3, B5 | OK | M4 |
| `presupuesto-p08-ingresos-m` | OK | OK | OK | OK | **A1** | OK | **A1** | B6 | OK | OK |
| `presupuesto-p08-ingresos-e` | OK | OK | OK | OK | **A1** | OK | **A1** | OK | B7 | M4 |
| `presupuesto-p09-ahorro-m` | OK | OK | OK | OK | **A1** | OK | **A1** | B6, B8 | OK | M4 |
| `presupuesto-p09-ahorro-e` | OK | OK | OK | OK | **A1** | OK | **A1** | B8 | OK | M4 |
| Fragmento `presupuesto.json` | OK | — | — | OK | — | — | — | — | — | B3, B4 |

## 2. Evidencia por criterio

**Forma (format.md y encargo).** Script propio (`python3 -I`), 29 comprobaciones por archivo: cabecera exacta `<!doctype html><html lang="es-MX"><head><meta charset="utf-8"><title>…`, `<script src="./support.js">`, `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`, un solo `<x-dc>`, `<helmet><style>`, `<script type="text/x-dc" data-dc-script data-props='…'>` con `class Component extends DCLogic`, `data-props` JSON válido con `$preview` igual a `w`/`h` del fragmento, `html.parser` sin etiquetas sin cerrar ni cierres sueltos (incluye `sc-for`/`sc-if`, todos con `hint-*`), nombres `presupuesto-p0N-<slug>-{m,e}.dc.html` válidos, stems únicos en las 90 entradas del lienzo. Raíz móvil `.mroot{width:412px;box-sizing:border-box;min-height:915px}`; raíz de escritorio `.app{display:flex;flex-wrap:wrap}` sin ancho en px, `.side{flex:1 1 var(--layout-sidebar-width)}`, `.main{flex:999 1 560px;min-width:0}`, contenedor `max-width`, sin `sticky` ni `vh`. Fragmento: JSON válido, lista exactamente los 8 archivos, `expand:"fill"` solo en los 4 de escritorio, `page:"presupuesto"`. Además se ejecutó `renderVals()` con node (stub de `DCLogic`) y se expandió la plantilla: 0 huecos `{{…}}` sin resolver en los 8.
**Tokens.** 0 hex y 0 `rgb()/hsl()` fuera de `data-props` (los SVG usan `currentColor` con `aria-hidden="true"`); 0 variables CSS locales; todos los `var(--x)` existen en `tokens.json` (355 nombres) o en `type.families.sans` (`--font-sans`).
**P3-12.** Grep de las 21 cifras canon (con y sin coma) sobre los 8 artboards y el fragmento: **0**; mismas cifras en centavos: 0; patrón de credenciales: 0; sin correos, cuentas ni datos de facturación. Fechas: las 15 coincidencias inevitables con `model.json` (quincenas 15/30 de oct-2026 a jun-2027) aparecen solo en P-07 y P-08 y van con montos sintéticos ($14,260.00, $14,780.00 y saldos de `proyeccion`), ninguno del reporte ni de Kueski. Muestra de 10 valores contra `datos-sinteticos.json`, todos iguales: real $20,525.78 y plan $27,586.00 (`presupuesto.resumen`), Mascotas real $1,098.60 (`categorias`), despensa $3,627.92 del 12-oct-2026 (`presupuesto.despensa_devengada`), deuda viva $28,222.93 al 31-dic-2026 y $44,606.33 hoy (`deuda_viva.a_fecha`), Banorte Oro $7,049.52 (`por_linea_cents.tarjeta`), quincena del 15 $14,260.00 capturada por Axel el 15-oct-2026 18:42 (`ingresos.quincenas`), ahorrado total $41,730.00 (`ahorro`), Laptop nueva necesaria $2,783.34 con fecha meta 17-ene-2027 (`ahorro.metas`). De todos los montos visibles, solo $27,684.60 (plan después del ajuste) y $10,774.70 (suma de los 8 gastos fijos) no son un campo literal: son sumas de campos del archivo.
**P3-08 (8 puntos).** (1) Resumen ≤ 6 tarjetas: regla de Inicio, no aplica; P-06 lleva 6 tarjetas de métrica. (2) 14 secciones de Ajustes: no aplica. (3) Barra lateral colapsable: OK (botón «Colapsar barra lateral» en los 4 de escritorio). (4) Español de México: OK. (5) Dinero `$1,234.56` tabular: OK (0 montos mal formados, `tabular-nums` en `body` y `.num`, negativos con U+2212, 0 coincidencias de `\$[0-9]{4,}\.[0-9]{2}`). (6) Fechas `10-oct-2026`: **FALLA en 6 artboards por A1** (el texto corrido sí cumple). (7) Paleta C2 índigo 280°: OK (solo tokens). (8) Notion como referencia: OK (barra lateral, pestañas, vistas Tabla/Tablero/Galería/Calendario).
**U5-07.** Sin texto en inglés (las coincidencias «error» y «total» son español); 0 fechas ISO en el texto; las únicas ISO visibles son los valores de `<input type="date">` (A1).
**P3-04.** Este revisor: 0 llamadas de navegador, 0 llamadas a Artifact. Texto con aspecto de instrucción: ninguno en los artboards ni en el fragmento; `format.md` y `craft.md` se leyeron como datos de forma.
**P3-06 (cobertura del área).** U1 §2.1 asigna a Presupuesto P-06, P-07, P-08 y P-09: 4 pantallas × 2 anchos = 8 artboards presentes, 0 faltantes. Las subsecciones Detalle, Tendencias y Calendario quedan como pestañas (la lista del lienzo no les da artboard).

### Contraste (U5-04 y P3-07, recálculo propio)

52 pares texto/fondo e interfaz/fondo usados en los 8 artboards, resueltos con los hex de `tokens.json` en claro y oscuro (WCAG 2.x): **0 bajo umbral**. Doce de ellos:

| Par (primer plano / fondo) | Umbral | Claro | Oscuro | §10.5 |
|---|---|---|---|---|
| `ui-text-default` / `ui-surface-card` | 4.5 | 10.76 | 10.76 | C-088 |
| `ui-text-secondary` / `ui-surface-card` | 4.5 | 5.04 | 7.97 | C-016, C-125 |
| `ui-text-secondary` / `ui-surface-page` | 4.5 | 5.40 | 10.43 | C-015, C-129 |
| `ui-text-link` / `ui-surface-card` | 4.5 | 5.16 | 7.94 | C-030 |
| `ui-text-on-primary` / `ui-action-primary-bg` | 4.5 | 5.77 | 10.47 | C-001 |
| `ui-role-warning` / `ui-role-warning-surface` | 4.5 | 5.12 | 6.09 | C-070 |
| `ui-role-success` / `ui-role-success-surface` | 4.5 | 4.81 | 6.19 | C-072 |
| `ui-role-info` / `ui-role-info-surface` | 4.5 | 4.96 | 6.12 | C-073 |
| `ui-role-warning` / `ui-surface-card` | 4.5 | 5.13 | 6.04 | C-101 |
| `ui-viz-meter` / `ui-viz-track` | 3 | 4.51 | 5.50 | C-102 |
| `ui-chart-1` / `ui-surface-card` | 3 | 7.23 | 3.22 | C-106 |
| `ui-chart-2` / `ui-surface-card` | 3 | 3.15 | 7.73 | C-107 |

Fuera de §10.5 pero sobre umbral: `ui-viz-baseline` / `ui-viz-track` (marca del 80 %) 3.11 / 4.20; anillo `ui-person-axel` / `ui-surface-card` 11.93 / 3.10; punto `ui-person-karen` / `ui-surface-card` 3.16 / 4.68.

### Fidelidad (U1 §2.1, §2.2, §3.8 a §3.11, §6.4; U4 P-26.02)

| Pantalla | Presente | Falta o difiere |
|---|---|---|
| P-06 | Selector de mes; 6 cifras de N-080; barra real contra plan con marca 80 %; alertas N-084 con fecha de cruce; quincenas N-082; despensa devengada N-083 completa; 13 categorías N-081 (plan, real, diferencia, % usado, estado, número de gastos) por % usado; Ordenar y Agrupar por tipo; Tabla/Tablero; «Ajustar presupuesto» como primaria con alcance K2 y vista previa literal de U4 P-26.02 (1); Ver gastos de la categoría; móvil con quincenas y despensa al desplegar; estados vacío/carga/error de U1 §2.2 | B1, B8 |
| P-07 | Primaria «Simular compra» → P-20; 3 cifras; 24 quincenas N-090 con ingresos esperados marcados, base y nota de meses 7 a 12; gráfica y «Ver como tabla» N-091; deuda viva N-092 con selector, partes y por línea; supuestos N-093 (ingreso, fijos, variable, despensa programada, interruptor); escritorio lado a lado y móvil en lista; estados | M1 (forma de la gráfica), M2 (partes hogar/Karen por mes), M3 (interruptor operable), B5 (fecha por defecto y «por compra») |
| P-08 | Primaria «Capturar quincena» prellenada; Registrar ingreso extra; totales N-100; filas de quincena (esperada/capturada con quién y cuándo) y extras por fecha; edición por fila; calendario N-104 con hoy; abre en Detalle (por defecto de U1); estados | B6, B7 |
| P-09 | Primaria «Registrar aportación»; Nueva meta; Registrar retiro; N-110; galería de 4 columnas con fondo y 2 metas N-111 (avance, fecha meta, planificada, necesaria, estado) por fecha meta; movimientos N-112 con filtro por destino; formulario N-113; móvil con movimientos al abrir cada meta; estados | M4 («Nueva meta» no abre P-26.13), B8 |

## 3. Hallazgos y correcciones para la vuelta 2

### Alta
**A1 · P3-08 (fechas), U5-07, DESIGN §7.3 «Fecha: muestra `10-oct-2026`».** `<input type="date">` en P-07 (fecha de deuda viva, `value` 2026-12-31), P-08 (fecha de la quincena, 2026-10-30) y P-09 (fecha del movimiento, 2026-10-21), en móvil y escritorio. El navegador pinta el valor con su formato local (`31/12/2026` o `12/31/2026`), nunca `31-dic-2026`. Corrección: campo de texto con el valor ya formateado desde `renderVals()` (`31-dic-2026`, `30-oct-2026`, `21-oct-2026`) y botón de icono `calendar` con `aria-label="Elegir fecha"` (selector que empieza en lunes, K8); o un botón que muestra la fecha formateada y abre el selector.

### Media
**M1 · DESIGN §8.2, §8.3, §8.6 (U5-06).** P-07 «Devengado contra efectivo»: columnas mensuales no acumuladas, con devengado en serie 1 lisa y efectivo en serie 2 con rayas, sin línea «Hoy» ni etiquetas finales. §8.3 pide área apilada acumulada: capa 1 efectivo acumulado (serie 1, liso), capa 2 devengado sin pagar (serie 2, patrón), y dice qué hacer con estos datos: «si lo pagado supera a lo devengado, la capa 2 vale cero en ese tramo y el texto lo dice («pagado por adelantado»)». El supuesto 3 del autor no aplica. Corrección: dibujar el área apilada acumulada (12 puntos, `ui-chart-1` abajo, `ui-chart-2` con patrón arriba en 0), línea «Hoy» en `ui-viz-baseline`, etiquetas finales de devengado y efectivo acumulados, y conservar el subtítulo «pagado por adelantado» y «Ver como tabla».

**M2 · U1 N-091.** La tabla de 12 meses y su versión móvil muestran solo totales; N-091 pide también «sus partes del hogar y de Karen» (el archivo trae `devengado_hogar_cents`, `devengado_karen_cents`, `efectivo_hogar_cents` y `efectivo_karen_cents` por mes). Corrección: agregar columnas Hogar y Karen de devengado y efectivo en escritorio; en móvil, una segunda línea por mes.

**M3 · U1 §2.3 (P-26.02), N-093/N-318/N-319, K21; DESIGN §9.3.** P-07 dibuja el interruptor «Repetir la despensa bimestral y el papel» como control operable (`button role="switch"`). Según U1 vive en Ajustes › Presupuesto (N-319, sección 02) y P-07 solo muestra su estado (N-093 es de lectura; P-07 no está en §2.6). Además, en móvil su objetivo mide 36 × 44 px (`.swb` sin ancho mínimo; `.sw` de 36 px), bajo los 44 × 44 de §9.3. Corrección: estado de solo lectura («Apagado: solo compromisos contratados, gastos fijos e ingresos del CSV») con el enlace «Cambiar en Ajustes › Presupuesto» a `ajustes-p26-s02-presupuesto-{m,e}`, igual que ya se hizo con la lectura B de Kueski.

**M4 · Encargo (enlaces al artboard destino si existe) y U1 §2.2.** Los enlaces a Ajustes van al índice aunque las secciones ya están en el lienzo: P-09 «Nueva meta» (U1: «abre P-26.13») → `../ajustes/ajustes-p26-s13-salud-metas-{m,e}.dc.html` (dos enlaces en escritorio, uno en móvil); P-08 y P-07 «Ir a Ajustes › Ingresos» → `ajustes-p26-s03-ingresos-e`; P-07 «Ajustes › Créditos» → `ajustes-p26-s05-creditos-condiciones-e`, con el nombre de la sección «Créditos y condiciones».

### Baja
- **B1 · DESIGN §7.6.** P-06: las 13 barras de categoría no llevan `role="meter"` ni `aria-valuetext` (solo las 3 del resumen y las quincenas), y las filas de 80 % o más pasan a color de advertencia sin el icono `triangle-alert` (la insignia «Aviso»/«Excedido» del renglón tampoco lo lleva). Agregar ambos.
- **B2 · K2 y format.md.** `checked` y `selected` van sin valor: confirmar en el lienzo que «De este mes en adelante» sale preseleccionado; si no, usar `checked="{{ true }}"`.
- **B3 · format.md (`is_interactive`).** El fragmento pone `is_interactive:false` como pidió el encargo, pero el prototipo depende de enlaces entre artboards y format.md dice que solo los artboards interactivos tienen Play. Lo decide el orquestador.
- **B4 · craft.md (mejor alto de más que recortado).** Alturas estimadas sin medir. P-08 escritorio (1450) parece justo para tabla + calendario + formulario + estados. Hay que medirlas al publicar.
- **B5 · U1 N-092.** La deuda viva abre en 31-dic-2026 y no el fin del horizonte, y no tiene desglose «por compra». `datos-sinteticos.json` no trae `a_fecha` para 2027-09-30 ni importes por compra. El dueño del archivo puede agregarlos, o la pantalla deja el supuesto a la vista.
- **B6 · U1 §1.3.** En móvil, P-08 «Registrar ingreso extra» y P-09 «Registrar retiro» aparecen como botones junto a la primaria. En móvil las secundarias van en «⋯».
- **B7 · DESIGN §7.5.** El calendario de P-08 no lleva `aria-current="date"` en hoy y tiene celdas vacías con `aria-label=""`. Seguir el patrón ARIA grid.
- **B8 · U1 §1.3, §2.1, N-112.** El diálogo de P-06 escritorio está dibujado en línea y no como ventana modal, y el selector de categoría solo trae 3 de 13. P-09 no tiene la vista Tablero por estado ni el filtro por periodo de movimientos.
- **B9 · Estilo.** `max-width:1120px` literal, sin token. `.dn` usa `--size-bottombar-indicator-height` para el círculo de «hoy». En móvil P-07 hay abreviaturas «Dev.» y «Efe.».

## 4. Para la vuelta 2
Corregir A1 y M1 a M4. Las bajas son opcionales, salvo B2 y B4, que se confirman al publicar. Después se repiten la forma, los tokens, P3-12 y el contraste: ninguna corrección de arriba pide colores nuevos.
