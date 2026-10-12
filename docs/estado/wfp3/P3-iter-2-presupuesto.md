# P3-iter-2 · Área «presupuesto» · revisión independiente (vuelta 2 de 3)

Fecha: 12-oct-2026 · Revisor: `revisor-presupuesto-r2` (no corrige) · Alcance: 8 artboards de `docs/design/maquetas/lienzo/project/presupuesto/` (P-06 a P-09, móvil 412 y escritorio 1440) y `docs/design/maquetas/lienzo/fragmentos/presupuesto.json`, después de `correccion-presupuesto-r1`.
Contra: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06, U5-07, en lo que se ve en el markup), `docs/design/DESIGN.md` (§7.3, §7.7, §8.2 a §8.8, §9.3, §10.5), `docs/design/U1-arquitectura.md` (§2.1, §2.2, §2.3, §6.4, N-090 a N-113, N-318, N-319), U4 (P-26.02) y `docs/design/maquetas/referencia-tipo/format.md`.
Diff contra la copia de la vuelta 1 (`docs/design/salida/iter-1/presupuesto/`): P-06 (m y e) sin cambios; P-07, P-08 y P-09 cambian solo en lo que pedían A1 y M1 a M4.

**Veredicto: NO APROBADO** · 1 alta (nueva) · 0 medias · 15 bajas (9 que siguen de la vuelta 1 y 6 nuevas).
Los 5 hallazgos de la vuelta 1 (A1, M1, M2, M3 y M4) están **cerrados**. El alta nueva (N1) la trae la línea auxiliar que se agregó a la gráfica para cerrar M1.

## 1. Hallazgos de la vuelta 1

| Id | Estado | Evidencia |
|---|---|---|
| A1 · fechas en `<input type="date">` | **Cerrado** | 0 `type="date"` en los 8. Los 6 campos son `type="text"`, con el valor que da `renderVals()`: `31-dic-2026` (P-07, `deuda_viva.a_fecha[1].fecha`), `30-oct-2026` (P-08, `ingresos.quincenas`, la del 30) y `21-oct-2026` (P-09, `meta.hoy`). Cada uno tiene `<label for>`, `aria-describedby` y un botón `calendar` con `aria-label="Elegir fecha"` de 44 × 44 px (`--size-touch-min` y `--size-control-coarse` valen 44px). No queda ninguna fecha ISO visible: las claves `iso` de P-08 no se usan en la plantilla. |
| M1 · forma de la gráfica de P-07 | **Cerrado en lo pedido** | Área apilada acumulada de 12 puntos. Capa 1, abajo: efectivo acumulado, serie 1 lisa. Capa 2: devengado sin pagar, serie 2 con rayas a 45° cada 6 px, con trazo de 1.5 px sobre relleno al 20 %; vale 0 porque `max(0, d − e)` da 0 en los 12 meses. Línea «Hoy · 21-oct-2026» de 1 px en `ui-viz-baseline`, colocada en proporción (5.4 %). Debajo de la gráfica van en texto los valores finales: $257,740.56, $0.00 y $242,520.65. El subtítulo dice «pagado por adelantado» y queda «Ver como tabla». Hay un `role="img"` con un `aria-label` calculado. Pero la línea auxiliar del devengado acumulado abre **N1**. |
| M2 · partes del hogar y de Karen (N-091) | **Cerrado** | Escritorio: tabla de 8 columnas dentro de una caja `overflow-x:auto` (format.md). Móvil: por mes, «Hogar … · Karen …» bajo el devengado y bajo el efectivo. Coteo: el oct-2026 del devengado de Karen da $906.98 y el sep-2027 del efectivo del hogar en móvil da $12,227.50, los dos iguales a `proyeccion.meses`. En los 12 meses, hogar + Karen = total. |
| M3 · interruptor operable en P-07 | **Cerrado** | 0 `role="switch"`. El estado es de solo lectura y sale de `presupuesto.proyeccion_meses_7_a_12.repetir_despensa_y_papel` (false): «Apagado: solo compromisos contratados, gastos fijos e ingresos del CSV.». El enlace «Cambiar en Ajustes › Presupuesto» va a `ajustes-p26-s02-presupuesto-{m,e}` y mide 44 px de alto en línea, más ancho que 44 px. |
| M4 · enlaces a Ajustes | **Cerrado** | P-09 «Nueva meta» va a `ajustes-p26-s13-salud-metas-{m,e}` (U1 §2.3: P-26.13 «Salud y metas de ahorro»). P-07 y P-08 «Ir a Ajustes › Ingresos» van a `ajustes-p26-s03-ingresos-e`. La lectura B de Kueski en P-07 va a `ajustes-p26-s05-creditos-condiciones-e` con el rótulo «Ajustes › Créditos y condiciones». El índice solo queda en el elemento «Ajustes» de la barra lateral. Todos los destinos existen en disco y tienen el mismo ancho. |

## 2. Tabla artboard × regla

Leyenda: OK = cumple; un id remite a §4.

| Artboard | Forma | Tokens y hex | Contraste U5-04 | P3-12 | P3-08 | U5-06 | U5-07 | U1/U4 | Táctil y a11y | Enlaces |
|---|---|---|---|---|---|---|---|---|---|---|
| `presupuesto-p06-mes-m` | OK | OK | OK | OK | OK | B1 | OK | B8 | OK | OK |
| `presupuesto-p06-mes-e` | OK | OK | OK | OK | OK | B1 | OK | B8 | OK | OK |
| `presupuesto-p07-flujo-m` | OK | OK | **N1** | OK | OK | N4, N5, N7 | OK | B5, N6 | N2, N3 | OK |
| `presupuesto-p07-flujo-e` | OK | OK | **N1** | OK | OK | N4, N5, N7 | OK | B5, N6 | N2, N3 | OK |
| `presupuesto-p08-ingresos-m` | OK | OK | OK | OK | OK | OK | OK | B6 | N2, N3 | OK |
| `presupuesto-p08-ingresos-e` | OK | OK | OK | OK | OK | OK | OK | OK | B7, N2, N3 | OK |
| `presupuesto-p09-ahorro-m` | OK | OK | OK | OK | OK | OK | OK | B6, B8 | N2, N3 | OK |
| `presupuesto-p09-ahorro-e` | OK | OK | OK | OK | OK | OK | OK | B8 | N2, N3 | OK |
| Fragmento `presupuesto.json` | OK | — | — | OK | — | — | — | — | — | B3, B4 |

## 3. Evidencia por criterio

**Forma (format.md y encargo).** Lo revisé con scripts propios (`python3 -I`) y en cada archivo comprobé lo siguiente:
- La cabecera exacta: `<!doctype html><html lang="es-MX"><head><meta charset="utf-8"><title>…</title><script src="./support.js"></script><link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css"></head><body><x-dc><helmet><style>`.
- La regla `body{background:var(--…);color:var(--…);font-family:var(--…)}`, junto con `a{…}` y `a:hover{…}`.
- Un solo `<x-dc>` y `</style></helmet>`.
- `<script type="text/x-dc" data-dc-script data-props='…'>` con `class Component extends DCLogic { renderVals() {`, y el documento cerrado con `</script></body></html>`.
- `data-props` es JSON válido y su `$preview` es igual a `w`/`h` del fragmento: 412 × 4400, 3440, 1880 y 2830 en móvil; 1440 × 2650, 2340, 1480 y 1880 en escritorio.
- `html.parser` no encuentra etiquetas sin cerrar ni cierres sueltos. Hay un solo `h1`, los ids son únicos y cada artboard tiene título propio.
- Todos los `sc-for` y `sc-if` llevan `hint-*`.
- Ejecuté `renderVals()` con node, con un stub de `DCLogic`, y resolví la plantilla con su ámbito: 0 huecos sin resolver.
- Nombres válidos. En el lienzo hay 90 stems, sin duplicados.
- La raíz móvil es `.mroot{width:412px;box-sizing:border-box;min-height:915px}`. La de escritorio no tiene ancho en px: `flex-wrap:wrap`, un contenedor con `max-width` y el contenido en `flex:999 1 560px`.
- Fragmento: JSON válido; lista exactamente los 8 archivos y `order` los mismos 8; `expand:"fill"` solo en los 4 de escritorio; `page:"presupuesto"`; `is_interactive:false`; sin x ni y.

**Tokens.** Fuera de `data-props` no hay hex, ni `rgb()`/`hsl()`/`oklch()`, ni colores con nombre. Los SVG de icono usan `currentColor` con `aria-hidden="true"`; los de la gráfica, `var(--…)` dentro de `style` y `color-mix(in srgb, var(--…) 20%, var(--ui-surface-card))`. No hay variables CSS locales. Los `var(--x)` de los 8 archivos existen todos en `tokens.json` (366 nombres más `font-sans`) y en `tokens.css`.

**P3-12.** Pasé el grep de las 21 cifras canon de P3-12, con y sin coma, sobre los 8 artboards y el fragmento: **0**. Las mismas cifras en centavos: 0. El patrón de credenciales: 0. No hay correos, cuentas ni datos de facturación.
Fechas: el comando de P3-12 lista 15 coincidencias entre `model.json`, `data.json` y `datos-sinteticos.json`. Solo `2026-10-30` aparece aquí, en P-08 (m y e), como la quincena del 30. Va con $14,780.00 esperado (`ingresos.quincenas`), que es sintético y no sale del reporte ni de Kueski.
Los datos incrustados en `renderVals()` son subconjuntos exactos de `datos-sinteticos.json`: `CAT` = `categorias`, `P` = `presupuesto`, `MS` = `proyeccion.meses`, `F.q` = `proyeccion.flujo_por_quincena`, `DV.a_fecha` = `deuda_viva.a_fecha`, `SU` = `presupuesto.gastos_fijos_recurrentes`, `.despensa_programada`, `.metodo_variable` y `.proyeccion_meses_7_a_12`, `I` = `ingresos`, `A` = `ahorro` y `MV` = `ahorro.movimientos`.
Muestra de 12 valores. Todos son iguales al archivo o a una suma de sus campos:
- efectivo acumulado $257,740.56, que es la suma de `efectivo_total_cents`;
- devengado acumulado $242,520.65;
- diferencia acumulada de sep-2027 −$15,219.91, igual a `diferencia_acumulada_cents` del mes 12 y a devengado − efectivo;
- devengado de Karen en oct-2026, $906.98;
- efectivo del hogar en sep-2027, $12,227.50;
- deuda viva $28,222.93 al 31-dic-2026;
- «Hoy · 21-oct-2026» (`meta.hoy`);
- quincena del 30, $14,780.00 al 30-oct-2026;
- ahorrado total $41,730.00 y fondo general $21,480.00;
- plan $27,586.00 e ingreso del mes $31,130.00 (`presupuesto.resumen`).

**P3-08 (8 puntos).**
1. ≤ 6 tarjetas: la regla es de Inicio. Aquí P-06 lleva 6.
2. 14 secciones: no aplica.
3. Barra lateral colapsable: OK, en los 4 de escritorio.
4. Español de México: OK.
5. Dinero `$1,234.56` tabular: OK. Hay `tabular-nums` en `body` y en `.num`, los negativos llevan U+2212 (`−$15,219.91`) y ninguna cadena de `renderVals()` tiene un monto mal formado. El eje usa la forma compacta: «$300 mil» y «$0».
6. Fechas `10-oct-2026`: **OK, ya en los 8** (A1 cerrado).
7. Paleta C2: OK, solo tokens.
8. Notion: OK.

**U5-07.** No hay texto en inglés: «total» y «error» son palabras en español. No hay fechas ISO visibles, ni emoji, ni barra de estado o teclado falsos.

**P3-04.** Este revisor no usó el navegador ni la herramienta Artifact. Leí `format.md` como datos de forma.

### Contraste (U5-04 y P3-07, cálculo propio con los hex de `tokens.json`)

Usé un motor de cascada propio: reglas del `helmet` más estilos en línea, texto, iconos, bordes de control y hover. Encontró 43 pares en los 8 artboards; ninguno de texto ni de control queda bajo el umbral. Los únicos bajo 3:1 son `ui-border-divider` y `ui-viz-grid`, que §7.1 declara «decorativo». Los pares de la gráfica nueva los calculé aparte, con la mezcla `color-mix` en sRGB (20 % de la serie y 80 % de `ui-surface-card`).

| Par (primer plano / fondo) | Umbral | Claro | Oscuro | §10.5 / estado |
|---|---|---|---|---|
| `ui-text-default` / `ui-surface-card` | 4.5 | 10.76 | 10.76 | C-088 |
| `ui-text-secondary` / `ui-surface-card` | 4.5 | 5.04 | 7.97 | C-016, C-125 |
| `ui-text-link` / `ui-surface-card` («Cambiar en Ajustes › Presupuesto») | 4.5 | 5.16 | 7.94 | C-030 |
| `ui-text-on-primary` / `ui-action-primary-bg` | 4.5 | 5.77 | 10.47 | C-001 |
| `ui-icon-default` / `ui-surface-field` (botón `calendar`) | 3 | 5.40 | 10.43 | OK |
| `ui-icon-default` / `ui-surface-hover` (hover del botón) | 3 | 4.51 | 5.50 | OK |
| `ui-border-control` / `ui-surface-card` (borde del botón y del campo) | 3 | 3.47 | 6.09 | OK |
| `ui-chart-1` / `ui-surface-card` (línea del efectivo) | 3 | 7.23 | 3.22 | C-106 |
| `ui-chart-2` / `ui-surface-card` | 3 | 3.15 | 7.73 | C-107 |
| `ui-viz-baseline` / `ui-surface-card` («Hoy», la mayor parte de su largo) | 3 | 3.47 | 6.09 | OK |
| **`ui-chart-2` / relleno de la capa 1** (línea a guiones del devengado acumulado) | 3 | **2.31** | 6.12 | **N1** |
| `ui-viz-baseline` / relleno de la capa 1 (tramo inferior de «Hoy», cerca del 7 % de su alto) | 3 | 2.55 | 4.83 | N1, secundario |
| `ui-chart-2` / relleno de la capa 2 (rayas) | — | 2.51 | 4.86 | así lo dice §8.6; hoy la capa mide 0 de alto |
| relleno de la capa 1 / `ui-surface-card` | — | 1.36 | 1.26 | el área se lee por su línea (C-106); ver N4 |

### Fidelidad (U1 §2.1, §2.2, §6.4, N-090 a N-113; U4 P-26.02)

| Pantalla | Presente | Falta o difiere |
|---|---|---|
| P-06 (sin cambios) | Todo lo que dio la vuelta 1: selector de mes, 6 cifras de N-080, barras con marca del 80 %, alertas 80/100, quincenas, despensa devengada, 13 categorías, «Ajustar presupuesto» con alcance K2 y vista previa de U4, estados | B1, B2, B8, B9 |
| P-07 | «Simular compra» a P-20; 3 cifras; 24 quincenas (N-090); área apilada con «Hoy» y valores finales, y «Ver como tabla» de 12 meses con totales, partes hogar/Karen y diferencia acumulada (N-091); deuda viva con campo de fecha `31-dic-2026`, partes y por línea (N-092); supuestos con el estado del interruptor de solo lectura (N-093); estados vacío/carga/error de §2.2 | B5, N1, N4 a N7 |
| P-08 | «Capturar quincena» prellenada con la fecha `30-oct-2026`; extras; totales N-100; calendario N-104; estados; vacío que lleva a P-26.03 | B6, B7 |
| P-09 | «Registrar aportación», «Nueva meta» (a P-26.13), «Registrar retiro»; N-110, N-111 y N-112; formulario N-113 con fecha `21-oct-2026`; estados | B6, B8 |

**Táctil (móvil, medido con los tokens).** Botón de fecha: 44 × 44. `.btn`: `--size-control-coarse` (44). `.btn.icb`: 44 × 44. `.tab`: 44. `details>summary`: 44. Enlaces sueltos en línea: `min-height:var(--size-touch-min)`. Barra inferior: 56. Opciones radio: 44 con su `<label>`. Campos: 44. Ningún control queda bajo 44 px.

## 4. Hallazgos y correcciones para la vuelta 3

### Alta
**N1 · U5-04 y P3-07 (UI ≥ 3:1 en los dos temas); DESIGN §8.6 y §10.5 (C-107 mide la serie 2 sobre la tarjeta, no sobre otro relleno).**
- **Dónde:** P-07, móvil y escritorio. Para que el devengado se vea aunque la capa 2 valga 0, la corrección agregó una línea a guiones del devengado acumulado (`ui-chart-2`, 2 px, 8·4; supuesto 2 de la corrección).
- **El problema:** el devengado acumulado queda debajo del efectivo en los 12 meses, de 4 a 14 px por debajo de la línea del efectivo. Así, la línea a guiones corre entera sobre el relleno de la capa 1, `color-mix(in srgb, ui-chart-1 20%, ui-surface-card)`, y no sobre la tarjeta.
- **Medido:** en claro, **2.31:1**, bajo 3:1; en oscuro, 6.12:1. En menor medida, el tramo inferior de la línea «Hoy» (`ui-viz-baseline`) cruza ese relleno a 2.55:1 en claro.
- **Corrección (cualquiera de las dos):**
  - (a) Dibujar debajo de la línea a guiones un halo con el mismo `d`, `stroke:var(--ui-surface-card)` y 4 a 6 px de ancho, y hacer lo mismo con «Hoy». Así el color junto a la línea es la tarjeta: C-107 da 3.15 y 7.73; la base, 3.47 y 6.09.
  - (b) Quitar la línea auxiliar, que §8.3 no pide, y dejar el devengado acumulado solo como valor final en texto. Para «Hoy» se aplica el halo de (a).

### Media
Ninguna.

### Baja
**Siguen de la vuelta 1:**
- **B1 · DESIGN §7.6.** P-06 no lleva `role="meter"` ni `aria-valuetext` en las 13 barras de categoría, ni `triangle-alert` en las de 80 % o más.
- **B2 · K2 y format.md.** `checked` va sin valor en P-06 (alcance) y en P-09 (tipo): confirmar al publicar que salen preseleccionados.
- **B3 · format.md.** `is_interactive:false` en un prototipo con enlaces: lo decide el orquestador.
- **B4 · craft.md.** Las alturas son estimadas y hay que medirlas al publicar. P-07 subió a 3440 y 2340 px.
- **B5 · U1 N-092.** La deuda viva abre en 31-dic-2026 y no en el fin del horizonte, y no tiene el desglose «por compra». Ahora el archivo trae `deuda_viva.por_defecto_fin_de_horizonte` = `2027-09-30` y `serie_mensual[11]`, que vale $0.00 y no tiene líneas. Abrir en esa fecha mostraría $0.00: el dueño de los datos decide si se agrega un valor útil o si la pantalla deja el supuesto a la vista.
- **B6 · U1 §1.3.** En móvil, «Registrar ingreso extra» (P-08) y «Registrar retiro» (P-09) siguen siendo botones visibles, no van en «⋯».
- **B7 · DESIGN §7.5.** El calendario de P-08 no tiene `aria-current="date"`; hay 0 coincidencias.
- **B8 · U1 §1.3, §2.1, N-112.** El diálogo de P-06 se dibuja en línea y su selector trae 3 categorías. En P-09 no están la vista Tablero ni el filtro por periodo.
- **B9 · Estilo (parcial).** Ya no están «Dev.» ni «Efe.». Siguen `max-width:1120px` literal, en los 4 de escritorio, y `.dn`.

**Nuevas:**
- **N2 · DESIGN §7.3 («`inputmode` abre el teclado correcto»).** Los 6 campos de fecha tienen `inputmode="numeric"`, pero su valor es `31-dic-2026`: el teclado numérico no deja escribir el mes. Usar `inputmode="text"` o quitar el atributo.
- **N3 · U5-07 (texto para el usuario).** La ayuda del campo dice «Formato 10-oct-2026. El selector empieza en lunes.». La segunda frase es una nota de diseño (K8), no texto para Axel. Dejar «Formato 10-oct-2026».
- **N4 · DESIGN §8.6.** La leyenda no reproduce la marca. La muestra de «Efectivo acumulado (liso)» va en color pleno, pero el área va al 20 %. La de «rayas» son franjas de color pleno sobre la tarjeta, pero la capa 2 lleva rayas plenas sobre un relleno al 20 %. Hacer que la muestra sea igual a la marca.
- **N5 · DESIGN §8.3 y §8.6.** Los valores finales van en una lista bajo la gráfica, no como etiqueta directa al final de cada serie (con ≤ 4 series se pide las dos cosas). El subtítulo no da la deuda viva con su suma, aunque N-092 está al lado.
- **N6 · Revisión visual del lienzo.** La tabla de 12 meses con las partes hogar/Karen (M2) está en un `<details>` cerrado en los dos anchos, así que en el artboard estático no se ve. Poner `open` en escritorio.
- **N7 · DESIGN §8.2 (eje X en el tiempo).** Los puntos están al fin de cada mes (`X = (i+1)/12`) y las etiquetas de mes están centradas en su celda: quedan desfasadas media celda. Alinear las etiquetas al final de la celda o poner los puntos al centro.

## 5. Para la vuelta 3
- Corregir N1. Es un cambio local en el SVG de P-07 (m y e) y no pide colores nuevos.
- Después repetir: forma, tokens, P3-12 y el contraste de la gráfica en claro y oscuro.
- Las bajas son opcionales. B2 y B4 se confirman al publicar.
