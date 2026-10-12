# correccion-ajustes-b-r1 (vuelta 2)

Fecha: 12-oct-2026. Copia de la vuelta 1 (P3-06): `docs/design/salida/iter-1/ajustes-b/` (14 artboards + `ajustes-b.json`).

## Hecho (12 hallazgos)
1. [alta] Perilla del interruptor: `.sw i:after` usa `--ui-surface-page` y `.sw.on i:after` usa `--ui-text-on-primary` (C-043/C-045). 10 artboards (s08, s09, s10, s11, s13, m y e). Revisado en captura oscura de s10-e.
2. [media] `.in,.sl` restituida en s11 y s12 (m y e) y en s14-m (que también usaba `.in` al añadir Notas).
3. [media] Quitados todos los `<code>` con claves técnicas (s08-e, s09-e, s11-m/e, s13-m/e).
4. [media] s09-e: cada campo vigente lleva «Valor inicial» en su celda (80 %, 100 %, 1 día, domingo, 09:00 x4, y 1 día antes y el mismo día). La lista duplicada «Parámetros operativos» se quitó. El campo «Primer aviso» pasó a 85 para cuadrar con la vista previa 80 → 85 que ya existía.
5. [media] s08-e: fechas del estado de cuenta de Karen como texto «01-sep-2026» y «30-sep-2026» con botón de calendario (44 px).
6. [media] Barra inferior fija (Inicio, Comidas, «+», Compras, Buscar; 56 px, destinos de 64 px o más; «+» a `../shell/shell-p03-registrar-gasto-m.dc.html`) en los 7 móviles, con el mismo marcado que comidas.
7. [media] Barra lateral nueva en los 7 escritorios: avatar y rol, botón de colapsar (panel-left), «Buscar» y «Registrar gasto» (enlazado a `shell-p03-registrar-gasto-e`), 9 destinos, pie con Tema (Sistema, Claro, Oscuro) y `aria-current="page"`.
8. [media] s10-e: `data-theme="light"` en la muestra clara (comprobado en captura con el lienzo oscuro).
9. [media] s11: «Elegir archivo» (input file oculto accesible + label) y «Lo que se creará» como tabla de 10 renglones (8 gastos fijos de `presupuesto.gastos_fijos_recurrentes` y 2 quincenas de `ingresos.esperado`); en móvil, lista de filas apiladas. La ficha de importar ocupa todo el ancho en escritorio.
10. [media] Campos añadidos: s13 Notas del perfil (Axel y Ximena), Inicio y Notas de cada meta (columnas en escritorio, campos en la tarjeta en móvil; Fondo solo Notas); s14 Notas (columna en escritorio; en móvil un `details` «Notas» por tarjeta, el primero abierto); s11 Hora del respaldo diario (03:00, «Sin valor inicial · Provisional») y huella («Huella: —», sin dato en el JSON).
11. [media] s14-e: días restantes en piso entero (2, 3, 3, 20, 16, 18, 17).
12. [media] s13: fecha de nacimiento pintada (opción 1 de la sugerencia), coherente con la edad.

Extras: `<ul>` con `<p>` hijo corregido en s13; texto de «Último respaldo» en s11 ya no es `nowrap` (desbordaba 412 px); `Valor inicial` de la alerta de pago separado en bloque propio.

Alturas: medidas con chromium headless local sobre una copia temporal (tokens.css copiado a /tmp, sin `support.js`). `$preview` y fragmento actualizados: s08 1508/2265, s09 1823/2280, s10 915/972, s11 3251/2946, s12 1502/1902, s13 3896/2577, s14 2102/1762 (móvil/escritorio).

## Verificación (`python3 -I`, 14 archivos, 0 fallas)
- Forma: doctype, lang es-MX, charset, support.js, tokens.css, `<x-dc>` → `<helmet><style>`, bloque `data-dc-script` con `class Component extends DCLogic` y `renderVals()`; `data-props` parsea y `$preview` = w x h del fragmento.
- html.parser: 0 etiquetas sin cerrar o huérfanas; 0 ids repetidos; 0 `for` huérfanos.
- 0 hex/rgb/hsl/oklch literales fuera de `data-props`; todos los `var(--x)` existen en `tokens.css`; 0 `var()` con respaldo.
- P3-12: 0 coincidencias de las 21 cifras canon (con y sin coma) y 0 credenciales; 0 innerHTML, iframe, red o localStorage; 0 `type="date"`; 0 signos menos mal formados.
- Clases usadas con regla (`.in`, `.sl`, `.sw`, `.bb`, `.sb`, etc.); móviles con raíz 412 x min 915 y barra inferior; escritorios con barra lateral completa y sin `aria-current="true"`.
- Fragmento: JSON válido, `boards` = `order` = los 14 archivos, `expand: "fill"` solo en escritorio.
- Desbordamiento horizontal móvil: 0 elementos pasan de 412 px en los 7 (medido).
- Capturas revisadas: s10-m, s10-e oscuro, s11-m, s11-e, s09-e, s13-e, s14-e.

## Supuestos
- Fechas de nacimiento sintéticas derivadas de `edad` (el JSON no las trae): Axel 14-mar-1990 (36 años), Ximena 22-jun-1993 (33 años).
- «Inicio» de metas derivado del primer movimiento del JSON: Laptop 18-may-2026, Viaje 04-may-2026.
- «Hora del respaldo diario» 03:00 tomada de la hora de los respaldos del JSON.
- Alturas medidas con las fuentes del equipo local; pueden variar unos píxeles con la tipografía final.

## Bloqueos
- Ninguno.

## Pendientes (para correccion-datos, fuera de mi alcance)
- `datos-sinteticos.json`: agregar `birth_date` por persona, `start_date` y `notes` de metas, hash abreviado de respaldos, y `dias_restantes` entero (`inventario.articulos`); hoy trae 2.3, 3.5, 20.2...
- K7: huevo, leche y tortilla salen «Reponer» con nivel «Sin aviso»; pedir nivel de reposición o estado según K7. No se tocó la pantalla para no contradecir el JSON.
- Revisar `compras/compras-p16-inventario-*` (mismos días decimales y mismo estado) y `ajustes-a` si la barra inferior o la barra lateral deben unificarse con la nueva.
- s11-e/s13-e: las tablas anchas se desplazan horizontalmente dentro de su tarjeta (`overflow-x:auto`), por diseño.

## Consumo aproximado
~170 KB leídos, ~110 KB escritos (artboards modificados en sitio), ~45 acciones.
