# artboards-gastos — reporte (Pista 3, fase 2 «lienzo», 12-oct-2026)

## Hecho
6 artboards en `docs/design/maquetas/lienzo/project/gastos/` (móvil 412 y escritorio 1440 de cada pantalla) y el fragmento `docs/design/maquetas/lienzo/fragmentos/gastos.json` (6 tableros, orden m/e por pantalla, `expand: fill` solo en escritorio).

| Archivo | Tamaño | Contenido |
|---|---|---|
| gastos-p02-todos-m.dc.html | 26 KB | P-02 móvil: lista de filas compactas (33 gastos), tablero de una columna con selector de categoría, agenda por día, galería de tickets; «Filtrar» (hoja), chips, totales, barra inferior |
| gastos-p02-todos-e.dc.html | 42 KB | P-02 escritorio: barra lateral, migas, Registrar gasto (primaria) y Foto de ticket, pestañas Resumen · Detalle · Tendencias · Calendario, resumen del periodo (N-042), barra de vista (Tabla · Tablero · Calendario · Galería, Filtrar, Ordenar, Agrupar, búsqueda), tabla con 11 columnas y fila de totales, tablero por categoría, calendario de octubre con detalle del día, galería, vista previa lateral, tarjeta «Estados» |
| gastos-p04-gasto-m.dc.html | 18 KB | P-04 móvil: 4 propiedades + «Mostrar 7 más», reparto, ticket adjunto, renglones plegados con total, compra a crédito, acciones |
| gastos-p04-gasto-e.dc.html | 30 KB | P-04 escritorio: propiedades en 2 columnas, reparto con devengo y efectivo, nota, ticket adjunto, compra a crédito vinculada, acciones (Duplicar con aviso «Deshacer»; Eliminar deshabilitado con su razón), 12 renglones del ticket con ajuste y totales por cuenta, «Estados» (carga, error) |
| gastos-p05-ticket-m.dc.html | 26 KB | P-05 móvil: pasos 1 a 4 con estado, renglones como tarjetas (alta plegados, baja desplegados con motivo y alternativa), Confirmar ticket, Guardar para después, tickets por confirmar |
| gastos-p05-ticket-e.dc.html | 38 KB | P-05 escritorio: pasos 1 a 4, foto a la izquierda y tabla de renglones a la derecha, borrador del ticket ocr2 con sus 2 correcciones (K3), «Estados» (vacío, carga, error) |

Todas las cifras salen de `datos-sinteticos.json` (gastos g01–g33, `gastos_resumen`, `tickets.confirmados[0]` = g13, `tickets.lectura_pendiente` = ocr1, `tickets.borradores` = ocr2, `creditos.kueski.compras` kc1). Las tablas llevan los datos en `renderVals()` y el dinero se formatea ahí con U+2212; el resto es texto literal. Enlaces entre artboards con `../<área>/<archivo>.dc.html` al destino del mismo ancho (shell, gastos, creditos, inicio, presupuesto, comidas, compras, salud-reportes, ajustes); Buscar enlaza a «#».

Interactivos (`is_interactive: true` en el fragmento): P-02 cambia de vista y selecciona gasto o día; P-04 despliega propiedades y renglones y duplica; P-05 recorre los pasos, abre el borrador y guarda para después.

## Verificación
Script propio `/tmp/wfp3-gastos/verify.py` (`python3 -I`): RESULTADO FINAL OK, 0 fallas. Comprueba por archivo: líneas obligatorias de la forma; JSON de `data-props` y ancho 412/1440; `html.parser` sin etiquetas sin cerrar; 0 hex literales fuera de `data-props` y de SVG; 0 `rgb()/hsl()`; 0 cifras de P3-12 (con y sin coma) y 0 valores en centavos prohibidos; todos los `var(--x)` existen en `sistema/project/tokens.css` (356 variables); sin iframe, innerHTML, http(s) ni emoji; sin signo «-» ASCII antes de `$`; `tabular-nums` presente; sin expresiones dentro de `{{ }}`; todos los agujeros `{{ x }}` tienen dato en `renderVals()` (ejecutado en Node con 13 estados distintos: vistas, pasos, borrador ocr2, etc.); el fragmento es JSON válido, lista exactamente los 6 archivos y su w/h coincide con `$preview`.

## Supuestos
- Alto de cada tablero ESTIMADO por cálculo de bloques (sin navegador): P-02 m 2870, e 2200; P-04 m 1960, e 2200; P-05 m 1790, e 1600. El orquestador debe ajustar `h` si el render difiere (la raíz usa `min-height` igual al alto declarado).
- Tablas con filas repetidas = `div` con roles ARIA (`table/row/cell/columnheader`), no `<table>`: un `<sc-for>` dentro de `<table>` lo saca el analizador de HTML. La tabla P-02 usa `aria-sort` en Fecha.
- `is_interactive: true` (la plantilla del encargo decía `false`) porque tablas, pasos y pestañas de vista funcionan; si prefieren estático, cambiarlo en el fragmento.
- Totales de P-02: la fila de totales usa `gastos_resumen` (29 gastos realizados: $24,413.20 · hogar $21,078.00 · Karen $3,335.20; los 4 venideros $2,858.70 no se suman). El «Gasto del hogar en oct-2026» ($20,525.78) es el devengado (suma de `por_categoria`) y la diferencia contra septiembre es la de `anterior`; son dos cifras distintas del JSON y las etiquetas lo dicen.
- No hay promedios diario/semanal/mensual (N-042) en el JSON: no se dibujaron ni se calcularon.
- Todas las filas de P-02 enlazan a P-04 (solo g13 trae renglones; es un prototipo). El detalle de P-04 siempre muestra g13.
- Forma de pago por defecto de P-05 «De contado» y paso 4 sin id de gasto: son estados de interfaz, no cifras del JSON. Ocr2 no trae categoría ni confianza: se muestra «Elegir categoría» y la insignia «Corregido/Revisado» según las correcciones del borrador. Las partes hogar/Karen del «Ajuste del ticket» propuesto de ocr1 se dejan en «—» (el JSON no las define).
- Iconos Lucide tomados de `fuentes/lucide/sprite.svg` (v1.55.0) como SVG inline con `stroke-width: var(--icon-stroke-user-units-sm)` (equivale a 1.5 px a 20 px). Barra lateral dibujada en cada escritorio con las mismas 9 secciones (Gastos desplegado); Main.dc.html es del shell.
- Hover solo bajo `@media (hover:hover)`; hojas (filtros, menú ⋯, confirmación) no dibujadas: solo sus disparadores.

## Bloqueos
Ninguno.

## Pendientes
- Medir alturas reales en el lienzo y corregir `h` del fragmento si hace falta.
- Confirmar que `support.js` y `ds/finanzas-hogar/tokens.css` (con `fonts/`) están publicados junto al lienzo; sin ellos los artboards se ven sin tokens.
- Verificar en el editor que `sc-if` dentro de `sc-for` y los `onClick` por elemento se comportan como describe `format.md` (no se pudo abrir navegador).

## Consumo aproximado
- Leído: ~230 KB (U1 por secciones, DESIGN §3.3 y §7, format/craft, LEEME y volcados parciales de datos-sinteticos.json, README de componentes del sistema).
- Escrito: ~117 KB redactados por mí (generadores Python ~104 KB, verificador ~7 KB, reporte ~5 KB); los 6 artboards generados suman 184 KB en disco (el CSS compartido y los iconos se repiten por archivo).
- Acciones de herramienta: ~52.
