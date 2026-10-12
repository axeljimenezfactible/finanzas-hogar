# correccion-creditos-r2 · vuelta 3 · Créditos (P-17 a P-20)

Agente: `correccion-creditos-r2` · Fecha: 12-oct-2026 · Fuente: `docs/estado/wfp3/P3-iter-2-creditos.md` (A-01 alto y M-01 medio). Copia de la vuelta 2 (P3-06): `docs/design/salida/iter-2/creditos/` (8 artboards y `creditos.json`, tal como estaban antes de tocar nada).

## Hecho

1. **A-01 (alto) · contraste de las barras de «deuda al cierre de cada mes».** El relleno pasó de `var(--ui-chart-1)` a `var(--ui-viz-meter)` (par C-102 de DESIGN §10.5: 4.51 en claro y 5.50 en oscuro sobre `--ui-viz-track`) en `creditos-p17-lineas-m`, `creditos-p17-lineas-e`, `creditos-p18-credito-m` y `creditos-p18-credito-e`. Solo se tocó el relleno de la barra; la línea de la gráfica A/B de P-18 conserva `--ui-chart-1` (es un trazo, no un relleno sobre pista). 0 apariciones de `ui-chart-1` como relleno de barra.
2. **M-01 (medio) · pagos tocables en la agenda de P-17 móvil.** Cada pago de la pestaña Calendario trae un botón de 44 px «Marcar pago» (o «Ver pago» si ya está pagado), con nombre accesible completo («Marcar pago de Kueski Pay, Vie 30 oct 2026»). Abre la misma hoja de P-18 m (`role="dialog"`, `aria-modal`): compra de origen, monto y fecha programados, capital repuesto, una casilla por parte con monto y estado, fecha real «21-oct-2026», Guardar y Cancelar. «Guardar» marca el pago como pagado en la agenda (ícono de check, «pagada», botón «Ver pago», y los filtros Pendientes/Pagados lo respetan). Las filas de corte no llevan botón. El capital repuesto de cada fila sale de `creditos.pagos_programados` de `datos-sinteticos.json` (cruce por línea, fecha y monto: 13 de 13 pagos con coincidencia única); la compra de origen sale de la descripción de la fila. Se agregó el símbolo `i-circle-check` y la clase `.st-ok` al artboard.
3. `$preview` y `h` de P-17 m de 2200 a 3000 (estimado: la agenda crece ≈ 48 px por fila); `creditos.json` actualizado. Es la única diferencia del fragmento con la copia de la vuelta 2.

Archivos modificados: `creditos-p17-lineas-m`, `-p17-lineas-e`, `-p18-credito-m`, `-p18-credito-e` y `fragmentos/creditos.json`. P-19 y P-20 (4 archivos) quedan idénticos a la vuelta 2 (`cmp`).

## Verificación (script propio `python3 -I` + prueba en node de `renderVals()`)

- Forma, 8 de 8: doctype, `lang="es-MX"`, charset, `./support.js`, `tokens.css`, un `x-dc`, `helmet`, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`; JSON de `data-props` parsea y su `$preview` es igual a `w`×`h` del fragmento.
- `html.parser`: 0 etiquetas sin cerrar o huérfanas; ids únicos.
- 0 hex literales, 0 `rgb()`/`hsl()` fuera de `data-props`; todas las `var(--x)` existen en `docs/design/sistema/project/tokens.css`.
- Cifras prohibidas de P3-12 (lista completa del grep de `docs/criterios/P3-claude-design.md`, con y sin coma): 0 en los 8 artboards y en el fragmento. 0 credenciales, 0 iframes, `innerHTML`, URL externas y «-$».
- Enlaces `../` relativos: todos existen (el de `ds/` lo pone el orquestador). Fragmento: JSON válido, `boards` y `order` listan exactamente los 8 archivos.
- Prueba en node de P-17 m: inicial, Calendario, hoja abierta, casilla alternada, Guardar (fila pasa a pagada), «Ver pago» (casillas deshabilitadas), filtros Pagados (1) / Pendientes (12) / Banorte Oro / Solo Karen (7), Cerrar: 0 errores, 0 variables del markup sin resolver.
- Resultado del script: `RESULTADO: OK`.

## Supuestos

- «Hoy» es 21-oct-2026 (fecha real propuesta de la hoja), igual que en P-18.
- Las casillas por parte alternan en la hoja, pero «Guardar» marca el pago completo (el prototipo no lee el DOM), como en P-18 m.
- Contraste de `--ui-viz-meter` sobre `--ui-viz-track`: se toma de DESIGN §10.5 C-102 y del informe de la revisión; no lo recalculé.
- La altura 3000 de P-17 m es estimada (sin navegador); el orquestador debe medirla.

## Bloqueos

- Ninguno.

## Pendientes

- Los 11 bajos de la revisión de la vuelta 2 (B-01 a B-11) no estaban en el alcance y siguen igual.
- Medir la altura real de P-17 m (pestaña Calendario con la hoja abierta).

## Consumo aproximado

Lectura ≈ 190 KB (informe de la vuelta 2, markup y script de P-17 m y P-18 m, JSON por consultas, criterio P3-12, reporte r1); escritura ≈ 14 KB (scripts de parche y verificación, este reporte); los artboards cambiaron ≈ 6 KB en disco; ≈ 26 acciones de herramienta.
