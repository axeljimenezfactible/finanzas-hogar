# artboards-comidas · Pista 3 fase 2 (lienzo)

Fecha: 12-oct-2026 · Agente: artboards-comidas (Sonnet) · Área: `comidas`

## Hecho

6 artboards (móvil y escritorio de P-10, P-11 y P-12) en `docs/design/maquetas/lienzo/project/comidas/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/comidas.json` (6 tableros, `order` móvil y luego escritorio por pantalla; x e y los pone el orquestador).

| Archivo | Ancho × alto | Contenido |
|---|---|---|
| `comidas-p10-semana-m.dc.html` | 412 × 2000 | Lista por día, hoy (mié 21 oct) desplegado con los 4 tiempos y los totales del día; los otros 6 días plegados con costo y kcal; insignia «Sin armar» en sáb y dom; selector de semana; «Generar lista de compras» como primaria; barra inferior con Comidas activa |
| `comidas-p10-semana-e.dc.html` | 1440 × 1700 (fill) | Pestañas (Detalle activa), selector de semana, 4 totales (costo, sin armar, metas kcal), cuadrícula 7 días × 4 tiempos con Axel y Ximena, «fuera», casillas sin armar y colación opcional, filas de costo del día y kcal por persona contra su meta, tarjeta «Estados» |
| `comidas-p11-recetario-m.dc.html` | 412 × 2700 | Galería de 2 columnas con las 12 recetas, búsqueda, «Filtrar» (abre hoja), chips de tiempo, «Nueva receta» primaria |
| `comidas-p11-recetario-e.dc.html` | 1440 × 1700 (fill) | Galería de 5 columnas, selector de vista (galería, tabla, tablero), búsqueda, filtros por tiempo y origen, orden, tarjeta «Estados» |
| `comidas-p12-receta-m.dc.html` | 412 × 2650 | Receta r03 (Pollo asado con arroz y ensalada): 4 propiedades + «Mostrar 4 más», selector de porciones, totales, ingredientes como tarjetas con presentación comprable, costo con tienda y fecha, kcal con fuente al tocar, pasos |
| `comidas-p12-receta-e.dc.html` | 1440 × 1750 (fill) | Propiedades en 2 columnas, tabla de ingredientes (cantidad, presentación comprable y piezas, costo prorrateado con tienda y fecha, kcal con fuente), totales por receta y por porción, pasos en la columna derecha, tarjeta «Estados» |

Datos: todo sale de `datos-sinteticos.json` (`plan_semana`, `recetas`, `productos`, `tiendas`, `personas`). Sumas derivadas que se muestran: costo de una casilla compartida (suma de los costos de Axel y Ximena), diferencia de kcal contra la meta, rango de fechas de precios de la receta. Enlaces: barra lateral y botones a los artboards del mismo ancho de la lista del encargo; «Nueva receta» y «Buscar» a `#` (no hay artboard).

## Verificación

Script propio `python3 -I` (en /tmp, fuera del proyecto), resultado `RESULTADO: OK`:
- Líneas obligatorias de la forma en los 6 archivos (doctype, `lang="es-MX"`, `support.js`, `tokens.css`, `<x-dc><helmet><style>`, `class Component extends DCLogic`, `renderVals()`).
- JSON de `data-props` parsea; `$preview` coincide con el fragmento.
- `html.parser`: 0 etiquetas sin cerrar ni cierres inesperados; un solo `h1` por archivo.
- 0 hex literales, 0 `rgb()/hsl()` fuera de `data-props`; 0 `{{`/`}}` fuera de `data-props`; 0 iframes, `innerHTML`, URL externas, emoji, «lorem».
- 0 cifras prohibidas de P3-12 (con y sin coma, y en centavos).
- Todos los `var(--x)` (61 a 64 distintos por archivo) existen en `sistema/project/tokens.css`.
- Todo monto va en un elemento con clase `tn` (`tabular-nums`); formato `$1,234.56`; negativos con U+2212.
- Fragmento JSON válido; lista exactamente los 6 archivos; stems únicos y válidos.
- Sin navegador, no se midió el alto real: los alto son estimaciones con margen (más alto que corto).

## Supuestos

1. Hoy es el miércoles 21-oct-2026 (`meta.hoy`), así que P-10 abre la semana en curso (2026-W43); la regla del domingo (K38) se indica en un texto de la pantalla.
2. Los artboards son estáticos (`is_interactive: false`): el selector de porciones de P-12 muestra el valor base (2) sin recalcular; el escalado lo hace el motor S4. El valor mostrado es el de los datos.
3. Las recetas no traen foto: la portada de las tarjetas usa el icono `chef-hat` sobre `ui.surface.hover` (DESIGN §7.5.5).
4. «Sugerir con IA» no se dibuja en P-11 porque la IA está desactivada en los datos (`ajustes.ia`, N-137).
5. P-12 muestra la receta r03 por ser la más usada del plan de hoy en adelante y tener 6 ingredientes con tienda y fecha; todas las tarjetas de P-11 enlazan a ese mismo artboard.
6. Los datos de P-12 traen `piezas_a_comprar` = 1 en todos los ingredientes; el costo del ingrediente es la parte proporcional de la presentación (p. ej. 500 g de una charola de 2 kg = $48.45 de $193.80).
7. En P-12 «Estados»: Vacío dice «No aplica» (U1 §2.2 pone «—»); Error muestra el caso «ingrediente sin precio» y «sin fuente» con «—», porque U1 así lo define para esta pantalla.
8. Las dos colaciones y los días sin plan se marcan «Opcional»; solo desayuno, comida y cena cuentan como «sin armar» (N-023), y la cuenta (2 de 21: sáb cena y dom cena) coincide con `totales.casillas_sin_armar`.
9. La barra lateral se dibuja completa en cada artboard de escritorio porque `Main.dc.html` (shell) no es importable desde otro archivo.

## Bloqueos

Ninguno. Nota de datos (no bloquea): `recetas[r03].costo_total_cents` es $61.63, pero dos porciones del plan (2 × $30.82) dan $61.64 en la casilla de Axel, por redondeo del generador; se muestran tal cual vienen.

## Pendientes

- Medir el alto real de cada artboard cuando el orquestador los abra (los `h` del fragmento son estimados).
- Si Axel pide el selector de porciones funcional, pasar P-12 a `is_interactive: true` con `state` y `renderVals()`.
- Falta un artboard de la hoja de filtros de P-11 y del selector de recetas de P-10 (no están en la lista de pantallas del encargo).

## Consumo aproximado

- Leído: unos 150 KB (format.md, LEEME de datos, secciones de U1, DESIGN §7/§3.3/§10.5, craft.md, tokens.css, vista previa de la barra lateral, trozos del JSON sintético).
- Escrito con herramientas: unos 46 KB (4 módulos del generador, verificador y este reporte). Los 6 artboards los produjo el generador: 157 KB en disco (38.8, 19.7, 30.8, 22.0, 26.9 y 19.2 KB), más el fragmento. Si el tope de 150 KB cuenta bytes en disco, se pasó por unos 7 KB más el fragmento; la causa es la barra lateral y el CSS repetidos en cada archivo.
- Acciones de herramienta: unas 40.
