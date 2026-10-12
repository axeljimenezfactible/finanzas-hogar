# correccion-presupuesto-r1 · Pista 3, fase 2 · área «presupuesto» (vuelta 2)

Fecha: 12-oct-2026 · Agente: correccion-presupuesto-r1 · Alcance: `docs/design/maquetas/lienzo/project/presupuesto/` (P-07, P-08, P-09; P-06 sin cambios), `docs/design/maquetas/lienzo/fragmentos/presupuesto.json` y este reporte.
Copia de la vuelta 1 (P3-06): `docs/design/salida/iter-1/presupuesto/` (8 artboards + `presupuesto.json`), hecha con `cp -n` antes de tocar nada.

## Hecho
1. **A1 (alta) · Fechas.** Los 6 `<input type="date">` (P-07 2026-12-31, P-08 2026-10-30, P-09 2026-10-21, móvil y escritorio) pasan a campo de texto con el valor formateado desde `renderVals()` (`31-dic-2026`, `30-oct-2026`, `21-oct-2026`) más un botón de icono `calendar` con `aria-label="Elegir fecha"` (44 px de ancho) y una nota «Formato 10-oct-2026. El selector empieza en lunes.».
2. **M1 · Gráfica de P-07.** Las columnas se sustituyen por un área apilada acumulada de 12 puntos (SVG inline con rutas calculadas en `renderVals()`): capa 1 efectivo acumulado en `ui-chart-1` liso (relleno al 20 % y borde de 2 px), capa 2 devengado sin pagar en `ui-chart-2` con patrón de rayas a 45 cada 6 px (vale $0.00 en todo el periodo), línea punteada de devengado acumulado (serie 2, guion 8·4), línea vertical «Hoy · 21-oct-2026» en `ui-viz-baseline` y etiquetas finales en texto (efectivo acumulado $257,740.56, devengado sin pagar $0.00, devengado acumulado $242,520.65). Se conservan el subtítulo con «pagado por adelantado» (calculado: $15,219.91) y «Ver como tabla». Eje Y de $0 a $300 mil.
3. **M2 · Partes del hogar y de Karen (N-091).** Escritorio: tabla de 8 columnas (mes, devengado total / hogar / Karen, efectivo total / hogar / Karen, diferencia acumulada) en caja con desplazamiento horizontal. Móvil: por mes, una línea con el total y otra con «Hogar … · Karen …» para devengado y para efectivo, más la diferencia acumulada; ya no hay abreviaturas «Dev.» y «Efe.» (B9 parcial).
4. **M3 · Interruptor.** Se quita el `role="switch"`; P-07 muestra el estado como texto de solo lectura («Apagado: solo compromisos contratados, gastos fijos e ingresos del CSV.», leído de `repetir.repetir_despensa_y_papel`) y el enlace «Cambiar en Ajustes › Presupuesto» (44 px de alto) a `ajustes-p26-s02-presupuesto-{m,e}`. Se retira el CSS del interruptor.
5. **M4 · Enlaces a Ajustes.** P-09 «Nueva meta» (cabecera y tarjeta de la galería, y estado vacío en escritorio; tarjeta en móvil) → `ajustes-p26-s13-salud-metas-{m,e}`; «Ir a Ajustes › Ingresos» de P-07 y P-08 (estado vacío) → `ajustes-p26-s03-ingresos-e`; en P-07 la lectura B de Kueski → `ajustes-p26-s05-creditos-condiciones-e` con el rótulo «Ajustes › Créditos y condiciones». Los únicos enlaces al índice que quedan son el elemento «Ajustes» de la barra lateral.
6. Alturas (`$preview` y `h` del fragmento, estimadas): P-07 3440 y 2340, P-08 1880 y 1480, P-09 2830 y 1880 (móvil y escritorio).

## Verificación
Script propio `python3 -I` (`/tmp/p3pres2/verify.py`, fuera del repositorio; el de la vuelta 1 más 6 comprobaciones nuevas): **360 comprobaciones, 0 fallas** sobre los 8 artboards y el fragmento. Cubre cabecera exacta, `support.js`, `tokens.css`, `<x-dc><helmet><style>`, `data-props` JSON con `$preview` 412 / 1440, `class Component extends DCLogic`, `html.parser` sin etiquetas sin cerrar (incluye `sc-for`, `sc-if`, `svg`, `pattern`), 0 hex y 0 rgb/hsl fuera de `data-props`, todos los `var(--x)` existen, 0 cifras prohibidas de P3-12 (con y sin coma) y 0 valores canon en centavos, sin `innerHTML`, `iframe`, red ni emoji, un solo `h1`, ids únicos, todos los enlaces dentro de la lista de pantallas (se agregaron las 14 secciones de Ajustes), cada `{{hueco}}` resuelve contra `renderVals()` ejecutado con node, y fragmento válido con los 8 archivos, `w`/`h` iguales a `$preview` y `expand: "fill"` solo en escritorio. Nuevas: 0 `type="date"`, 0 `role="switch"`, 0 enlaces al índice de Ajustes dentro del contenido, área apilada y partes hogar/Karen presentes en P-07, campo de fecha con botón en P-07/P-08/P-09 y ninguna fecha ISO en el valor del campo.
Cotejo de datos: devengado y efectivo acumulados salen de `proyeccion.meses` (suma de `*_total_cents`); la diferencia final ($15,219.91) coincide con `diferencia_acumulada_cents` del mes 12; hogar + Karen = total en cada mes.

## Supuestos
1. La capa «devengado sin pagar» vale 0 en los 12 puntos porque el efectivo acumulado supera siempre al devengado (§8.3: «pagado por adelantado»); se dibuja igualmente (ruta de área de alto cero) para que el patrón aparezca si cambian los datos, y la leyenda lo declara.
2. Se añadió una línea punteada de devengado acumulado (serie 2, trazo 8·4) para que el devengado se vea aunque la capa 2 sea cero; §8.3 no la pide, pero sin ella la cima de la pila no sería el devengado.
3. «Hoy» es el 21-oct-2026 (`meta.hoy`); su posición es proporcional dentro del mes (5.4 %) sobre un eje de 12 meses de ancho uniforme con los valores al fin de cada mes.
4. El SVG usa `preserveAspectRatio="none"` y `vector-effect: non-scaling-stroke`; el patrón se deforma como mucho levemente al escalar y hoy es invisible por valer cero.
5. El campo de fecha es de texto libre con formato `10-oct-2026`; el botón de calendario es un dibujo (el selector real con inicio en lunes, K8, lo da la app).
6. Alturas estimadas a mano, sin navegador (P3-04); hay que medirlas al publicar (B4 de la revisión).
7. Bajas B1, B2, B5 a B8 y B9 (salvo las abreviaturas) no se tocaron: la vuelta 2 pedía solo alta y media.

## Bloqueos
Ninguno.

## Pendientes
- Medir alturas reales en el lienzo y corregir `h` y `$preview`.
- Revisar en el lienzo que `color-mix`, el patrón SVG y los huecos en `d=""` / `style=""` se interpretan como se espera.
- Bajas de la revisión (B1 `role="meter"` en P-06, B6 secundarias en «⋯» en móvil, B7 `aria-current` del calendario de P-08, B8 modal de P-06 y vista Tablero de P-09).

## Consumo aproximado
- Leído: ~170 KB (informe de revisión, artboards P-07 y fragmentos de P-08/P-09, DESIGN §8, format.md).
- Escrito: ~12 KB de script de parche (`/tmp/p3pres2/patch.py`) y ~3 KB de reporte; artboards modificados 6 de 8 (+10 KB netos).
- Acciones de herramienta: ~22.
