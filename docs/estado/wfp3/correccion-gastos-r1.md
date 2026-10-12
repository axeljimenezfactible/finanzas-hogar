# correccion-gastos-r1 (vuelta 2) · Área «gastos»

Fecha: 12-oct-2026. Corrige los hallazgos alta y media de `P3-iter-1-gastos.md` (H1 a H7). Los bajos (H8 a H12) no se tocaron, salvo lo que cae dentro de H7 y H3.

## Hecho
- Copia previa (P3-06): los 6 artboards y `gastos.json` de la vuelta 1 están en `docs/design/salida/iter-1/gastos/`, copiados antes de editar.
- H1 (alta): regla `::placeholder{color:var(--ui-text-secondary);opacity:1}` en el helmet de p02-e y p04-e (también en p04-m, que ahora tiene textarea con placeholder). Contraste 5.40:1 en claro y 10.43:1 en oscuro.
- H2: «Gasto del Mié 12 oct 2026» pasa a «Lun 12 oct 2026» en p04-m y p04-e.
- H3: p04-m gana la sección «Nota» (card con `label for="nota"` y `textarea`, mínimo `var(--size-touch-min)`, mismo placeholder), tras Reparto y antes de Ticket adjunto. Siempre visible: «Mostrar 7 más» no cambia. Alto del artboard 1960 a 2180 (raíz `min-height`, `$preview` y fragmento).
- H4: quitado `outline:0` del `input type=search` de p02-e; ahora rige `:focus-visible`.
- H5: los enlaces «Compra a crédito» y `ticket-walmart-2026-10-12.jpg` de p04-m llevan `display:inline-flex;align-items:center;min-height:var(--size-touch-min)`.
- H6: cada monto citado va en `<span class="num">`: p02-e ($25,126.40 y $2,858.70), p02-m ($25,126.40, $2,858.70 y `{{ a.tot }}`), p04-e ($2,929.00, −$10.60 y los 3 montos de «Totales por cuenta»), p04-m ($2,918.40 del botón de renglones, `{{ r.unit }}` y los 2 de «hogar … Karen …») y p05-e (`{{ L.suma }}` y `{{ L.ajuste }}`).
- H7: botón de texto «+ Nuevo» al pie de cada columna del tablero, con enlace a `../shell/shell-p03-registrar-gasto-e.dc.html` (escritorio) y `-m` (móvil, `min-height:var(--size-touch-min)`).

## Verificación (python3 -I, script propio en /tmp, borrable)
- Forma en 6/6: doctype, `lang="es-MX"`, `support.js`, `tokens.css`, `x-dc`, helmet, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`. JSON de `data-props` válido y `$preview` igual a w/h del fragmento.
- `html.parser` con pila propia: 0 etiquetas sin cerrar ni huérfanas.
- 0 hex literales fuera de `data-props`; 0 `rgb()`/`hsl()`; 0 `innerHTML`, iframes ni URLs de red.
- Todos los `var(--x)` existen en `tokens.json` o `tokens.css`.
- 0 cifras prohibidas de P3-12 (regex tomada de `tools/ds/verificar_sinteticos.py`, con y sin coma); 0 credenciales; 0 «-$» ASCII.
- Heurística de dinero literal fuera de `.num`: marcaba 7 frases en la vuelta 1 (copia en `salida/iter-1/gastos`) y ahora marca 0. Interpolaciones `{{ }}` de dinero revisadas a mano: todas dentro de `.num`.
- Fragmento `gastos.json`: válido, lista exactamente los 6 archivos, `expand:"fill"` solo en escritorio.

## Supuestos
- La sección Nota va siempre visible en móvil (no dentro de «Mostrar N más»); por eso el texto sigue siendo «Mostrar 7 más».
- Alto de p04-m (2180) estimado: Nota suma unos 200 px; no se midió en navegador (sin navegador por K45).
- «+ Nuevo» sin icono, para igualar el texto del hallazgo.

## Bloqueos
Ninguno.

## Pendientes
- Bajos H8 a H12 (aria-current en «hoy», fechas ISO en nombres de archivo, doble «Ver ticket», px literales, `is_interactive`) sin tocar; quedan a decisión del orquestador.

## Consumo aproximado
~55 KB leídos, ~12 KB escritos, 26 acciones.
