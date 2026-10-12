# artboards-inicio: reporte (Pista 3, fase 2 «lienzo»)

12-oct-2026 · agente artboards-inicio · P-01 Inicio, 8 artboards (móvil 412 y escritorio 1440 de cada subsección).

## Hecho
- `docs/design/maquetas/lienzo/project/inicio/`:
  - `inicio-p01-resumen-m` (412 × 1840) y `-e` (1440 × 1250): las 6 tarjetas de U3 §2 en orden fijo, con sus métricas, insignias, enlace único y destino. En escritorio van 3 × 2 y llevan la fila «Estados» (vacío, carga y error de «Mes contra plan»).
  - `inicio-p01-detalle-m` (412 × 4250) y `-e` (1440 × 3000): barra de vista (periodo Mes, 1.ª y 2.ª quincena; filtros en hoja o en la barra), los paneles D-1 a D-7 (U3 §3 y §3.1) y el enlace «Ver todos los gastos con estos filtros». En escritorio van 2 columnas, D-7 a todo el ancho y la fila «Estados» de D-1.
  - `inicio-p01-tendencias-m` (412 × 2450) y `-e` (1440 × 1950): T-1 columnas agrupadas real/plan (plan con rayas a 45°), T-2 área apilada de efectivo acumulado con la línea de devengado acumulado, «Hoy» y proyección, T-3 línea de deuda viva y T-4a/T-4b líneas de costo y kcal por persona con meta. Cada gráfica va en `<figure>` con `<figcaption>`, un resumen de una frase, una leyenda con forma y trazo y «Ver como tabla». En escritorio van 2 por fila y llevan la tarjeta «Estados».
  - `inicio-p01-calendario-m` (412 × 3050): agenda del 21-oct al 20-nov con separadores de quincena, la franja de la promoción en curso y «Ver 30 días más». `-e` (1440 × 1600): mes de octubre de lunes a domingo con franjas de quincena (continua y punteada), hasta 3 eventos por día y «+n más», el panel del día seleccionado (Vie 30 oct 2026) y la tarjeta «Estados».
- `docs/design/maquetas/lienzo/fragmentos/inicio.json`: 8 tableros con w, h, título «P-01 · Inicio › … · móvil|escritorio», page «inicio», `expand: fill` solo en escritorio y `is_interactive: false`, más el orden.
- Todos los datos salen de `datos-sinteticos.json` (`inicio`, `categorias`, `proyeccion`, `deuda_viva`, `presupuesto.tendencia`, `salud.tendencia_semanal`, `canasta.por_tienda`, `calendario_hogar`, `participacion_karen`, `ahorro.metas`, `creditos.kueski.disponible_hoy`). Hoy es Mié 21 oct 2026.

## Verificación
Corre con `python3 -I /tmp/fh-inicio/verificar.py /home/deepsec/finanzas-hogar` y da RESULTADO: OK. Revisa:
- las líneas obligatorias de la forma, que `data-props` parsee y que `$preview` coincida con el fragmento;
- que html.parser no encuentre etiquetas sin cerrar ni cruzadas, y que todo `sc-for` y `sc-if` lleve su `hint-*`;
- que haya 0 hex literales fuera de `data-props` y 0 cifras de P3-12 (con y sin coma);
- que todo `var(--x)` exista en `tokens.css`, generado de `tokens.json`;
- que el dinero tenga el formato `$1,234.56`, con el compacto solo en los ejes, y que el negativo use U+2212;
- que no haya emoji, innerHTML, iframes ni red, y que el bloque `<style>` solo traiga CSS;
- que la raíz móvil mida 412 px fijos y que la de escritorio no tenga ancho fijo;
- que cada enlace apunte a la lista de pantallas y al mismo ancho; además, todos los destinos ya existen en disco.

Además, `node` ejecuta `renderVals()` de los 8 artboards: cada hueco `{{ }}` tiene su valor, hay 35 celdas en el mes, 28 eventos en la agenda, los anchos de D-1 están entre 0 y 100, y la celda de hoy (21) y la del día seleccionado (30) salen marcadas.

## Supuestos
1. **Altos estimados.** No hubo navegador (K45), así que el alto de cada tablero se calculó a mano. La raíz lleva `min-height` igual a ese alto. Conviene que el orquestador lo mida al publicar.
2. **T-2.** La diferencia acumulada es negativa en los 12 meses, porque se pagó por adelantado. Por eso la capa 2 vale cero y el texto lo dice. Agregué la línea discontinua de devengado acumulado (serie 2, cuadrado) para que se vean las dos cantidades. Los acumulados suman `proyeccion.meses` desde oct-2026 y cuadran con `diferencia_acumulada_cents` del último mes (comprobado).
3. **T-1** usa la ventana fija de 12 meses: de nov-2025 a abr-2026 va «Sin datos» (regla del LEEME).
4. **T-4b** toma la meta de kcal de `inicio.tarjetas.comidas_semana` y la aplica a las 12 semanas de `salud.tendencia_semanal`, cuya última semana es la del 12 oct.
5. **D-1** dibuja la barra del real con escala común (el máximo es Vivienda, $7,200.00) y el plan como marca vertical (DESIGN §8.2).
6. **D-2** sale de `proyeccion.flujo_por_quincena` (2026-10-q2), con el formato hecho desde centavos. Leí `pagos_credito_cents` como la parte del hogar.
7. **D-5** no muestra la marca «viejo», porque `articulos_mayor_diferencia` no la trae. El total por tienda incluye Sam's Club (4 de 24) y los dos mercados.
8. El nombre de la línea es «Kueski Pay», como en los datos, y no «Kueski», como en el ejemplo de U3.
9. **Enlaces.** Ajustes en la barra lateral abre 01 Cuenta y roles (U1 §1.2.3). «Capturar presupuesto» va a la s02 y «Capturar ingreso» a la s03. Buscar (G-04) y las promociones enlazan a «#».
10. **Iconos Lucide** escritos de memoria, con la geometría aproximada. Ajustes usa `settings-2`, más corto que `settings`.
11. **Estilos.** Hay clases compartidas en `<helmet>` (además de body, a y a:hover) para no repetir estilos. Los contenedores conservan flex o grid con `gap` en línea. No hay huecos en `style`: los anchos de las barras van como atributos numéricos de `<rect>`.
12. La tarjeta 3 no muestra «Por vencer», porque faltan 9 días y los días de aviso son 1 y 0.

## Bloqueos
- Ninguno.

## Pendientes
- **Orquestador:** medir el alto real de cada tablero y ajustar `$preview`, `min-height` y el fragmento si difieren.
- **Revisión en el runtime de Design** (no se pudo sin navegador): los huecos de atributo numérico en `<rect>`, `sc-for` anidado (días y eventos) y `aria-pressed` con booleano.

## Consumo aproximado
- **Leído:** unos 230 KB, de U1, U3, DESIGN §7–§10.5, format.md, craft.md, LEEME, extractos de datos-sinteticos.json y tokens.css.
- **Escrito:** unos 189 KB de artboards, 2 KB del fragmento y 5 KB de este reporte. Es más que el tope de 150 KB: los 8 tableros completos no cupieron. Aparte quedan unos 45 KB de borradores en `/tmp/fh-inicio` (geometría SVG, ensamblador y verificador).
- **Acciones de herramienta:** unas 50.
