# correccion-ajustes-b-r3 · vuelta 4 (ronda extra, mecánica) · 12-oct-2026

Agente `correccion-ajustes-b-r3`. Entrada: `docs/estado/wfp3/P3-iter-3-ajustes-b.md` (alta R3-A1, media H-11) y `correccion-datos-r4` (umbrales nuevos).
Copia previa (P3-06): `docs/design/salida/iter-3/ajustes-b/` (14 artboards de ajustes-b, `ajustes-b.json` y `compras-p16-inventario-{m,e}`), copiada antes de editar.

## Hecho

- **R3-A1 (alta, U5-04 / P3-07 / DESIGN C-035)**: regla `::placeholder{color:var(--ui-text-secondary);opacity:1}` agregada al `helmet` de los 14 artboards (s08 a s14, `-m` y `-e`). Contraste calculado por el revisor: 5.40:1 en claro y 10.43:1 en oscuro sobre `ui-surface-field`. No cambia alturas.
- **H-11 / K7 (media, U4 P-26.14)**, con el dato corregido (`umbral_reposicion` huevo 12 pza, leche 2,000 ml, tortilla 1,000 g; estado «reponer»; resumen 1/5/8 sin cambio):
  - `ajustes-p26-s14-inventario-m`: «nivel 12 pza», «nivel 2,000 ml» y «nivel 1,000 g» en lugar de «Sin aviso». Quedan 4 «Sin aviso» legítimos en `-e` (croquetas, detergente, jabón de trastes y papel higiénico, `umbral_reposicion` null) y ninguno en `-m` (que solo muestra los primeros 7 productos).
  - `ajustes-p26-s14-inventario-e`: columna «Nivel de reposición» con `12 pza`, `2,000 ml` y `1,000 g` (en `span.m`, tabulares).
  - `compras-p16-inventario-m`: «Ritmo 18 pza / sem · Reponer desde 12 pza» (y los equivalentes de leche y tortilla). Se acortó «Ritmo de consumo» a «Ritmo» porque la línea es `nowrap` con elipsis y el texto largo se habría cortado.
  - `compras-p16-inventario-e`: columna de nivel «Reponer desde 12 pza / 2,000 ml / 1,000 g» en lugar de «Calculado por ritmo» en esas tres filas. El panel de detalle no cambió (no muestra nivel).
- Fragmento `ajustes-b.json`: sin cambios (alturas iguales). s14-m medida con chromium headless sin `support.js`: 2102 px, igual que `$preview`.

## Verificación (script propio, `python3 -I`, /tmp/chk.py)

| Comprobación | Resultado |
|---|---|
| Forma de los 14 artboards (doctype, `lang="es-MX"`, charset, `support.js`, `tokens.css`, `x-dc > helmet > style`, `data-dc-script` con `DCLogic` y `renderVals`) y la regla `::placeholder` | 14/14 |
| `data-props` parsea y `$preview` = `w` x `h` del fragmento | 14/14 |
| html.parser: etiquetas sin cerrar o con cierre cruzado (14 ajustes + 2 compras) | 0 |
| Hex literales fuera de `data-props` (sin `href="#…"`) | 0 |
| `var(--x)` contra `tokens.css` | todos existen |
| Cifras prohibidas de P3-12 (21 patrones, con y sin coma) y credenciales | 0 · 0 |
| Montos de 4 o más cifras sin separador de miles | 0 |
| `innerHTML`, `iframe`, `fetch` | 0 |
| Fragmento: JSON válido, `boards` = `order` = los 14 archivos | OK |
| Los 2 artboards de compras-p16: forma y `data-props` | OK |

## Supuestos

- La regla de estado del dato es la del generador: con nivel y existencia <= nivel, «reponer»; así las tres filas siguen en «Reponer» y el texto «reponer si la existencia es igual o menor que el nivel» ya es coherente.
- En compras-p16 móvil se abrevia «Ritmo de consumo» a «Ritmo» solo en las tres filas con nivel; el resto de filas de ritmo no cambia.

## Bloqueos

- Ninguno. H-11 queda cerrado por el dato corregido (no se registra bloqueo en `docs/estado/bloqueos/`).

## Pendientes

- Alturas de compras-p16 no remedidas (el texto se limitó para no cambiar la altura; el panel de detalle no cambió).
- Las 18 bajas del informe no estaban en el alcance de esta ronda.

## Consumo aproximado

~110 KB leídos (varios volcados grandes de HTML), ~12 KB escritos, ~22 acciones.
