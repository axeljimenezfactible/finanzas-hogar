# correccion-ajustes-b-r2 · vuelta 3 · 12-oct-2026

Agente `correccion-ajustes-b-r2`. Alcance: los 14 artboards de ajustes-b (P-26.08 a P-26.14), `fragmentos/ajustes-b.json` y este reporte. Entrada: `docs/estado/wfp3/P3-iter-2-ajustes-b.md` (alta R2-A1, medias H-11 y R2-M1).

Copia de la vuelta 2 (P3-06): `docs/design/salida/iter-2/ajustes-b/` (14 artboards y `ajustes-b.json`, copiados antes de editar).

## Hecho

- **R2-A1 (alta, U5-04 / P3-07)**: s13-m y s13-e. Se quitó `style="opacity:.85"` de la frase «Con “mantener” no hay ajuste…» (queda en `ui-text-secondary` sobre la tarjeta: 5.04:1 en claro, 7.97 en oscuro). Se dibujaron los dos campos de ajuste de Ximena, «Mi ajuste para bajar (kcal)» y «Mi ajuste para subir (kcal)» (`abximena`, `asximena`), con `readonly`, `aria-disabled="true"` y `aria-describedby` hacia la frase. Regla nueva `.in[aria-disabled="true"]` con `ui-surface-disabled`, `ui-text-disabled` y `ui-border-divider`. Esto cierra también B-17. El valor «0» es el de la casa que la misma pantalla ya muestra.
- **R2-M1 (media, P3-08 punto 5 / DESIGN §3.3)**: 7 montos envueltos en `<span class="m">` (tabular-nums): 3 en s09-e (`$27,586.00` x3) y 2 en cada uno de s12-m y s12-e (`$2,918.40` en las dos filas «Con un gasto de…»).
- **Alturas**: medidas con chromium headless (copia temporal en /tmp, sin `support.js`). s13-m pasó de 3896 a 4060 (`$preview` y fragmento); s13-e 2577, s09-e 2280, s12-m 1502 y s12-e 1902 no cambian.

## Verificación (script propio, `python3 -I`, /tmp/vv/v.py)

| Comprobación | Resultado |
|---|---|
| Forma de los 14 artboards (doctype, `lang="es-MX"`, charset, title, `support.js`, `tokens.css`, `body > x-dc > helmet > style`, bloque `data-dc-script` con `DCLogic` y `renderVals`) | 14/14 |
| `data-props` parsea y `$preview` = `w` x `h` del fragmento | 14/14 |
| html.parser: etiquetas sin cerrar o con cierre cruzado | 0 |
| Ids repetidos · `label for` huérfanos | 0 · 0 |
| Hex, `rgb()`, `hsl()` y `oklch()` fuera de `data-props` (sin contar `href="#…"`) | 0 |
| `var(--x)` contra `sistema/project/tokens.css` | todos existen (incluye `ui-surface-disabled`, `ui-text-disabled`) |
| Cifras prohibidas de P3-12 (21 patrones, con y sin coma) | 0 |
| Dinero en prosa fuera de `.m`/`tabular-nums` | 0 (antes 7) |
| `opacity` en texto | 0 (solo queda la regla `.dis{opacity:.6}` de s08/s10/s11, sin uso en ningún elemento) |
| `innerHTML`, `iframe`, `fetch`, `localStorage` | 0 |
| Fragmento: JSON válido, `boards` = `order` = los 14 archivos | OK |
| Contraste del texto corregido | 5.04:1 en claro (recalculado: #66676e / #f2f2f4) |

## Supuestos

- Los campos deshabilitados de Ximena están exentos de contraste (WCAG 1.4.3) por ser controles inactivos; el texto que los explica sí cumple 4.5:1.
- Las alturas de s13-m cambian con la tipografía final; son las del equipo local.

## Bloqueos

- **H-11 (K7 / U4 P-26.14), sin cerrar**: `docs/design/maquetas/datos-sinteticos.json` sigue como el 11-oct 23:55 (`inventario.articulos`: huevo, leche y tortilla con `umbral_reposicion: null` y `estado: "reponer"`; `resumen` 1/5/8). Está fuera de mi alcance y las pantallas pintan lo que dice el JSON, así que s14-m y s14-e siguen mostrando «Reponer» con «Sin aviso» junto a la regla «reponer si la existencia es igual o menor que el nivel». Requiere que `correccion-datos` dé a esos tres artículos un umbral coherente con la existencia o los pase a «suficiente» y ajuste el resumen; después hay que repintar s14-m/e y `compras-p16-inventario-m/e`.

## Pendientes

- Cerrar H-11 cuando cambie el JSON (repintar s14 y compras-p16).
- Bajas B-01, B-03 a B-16, B-18 y B-19 de `P3-iter-2-ajustes-b.md` no estaban en el alcance de esta vuelta.

## Consumo aproximado

~25 acciones de herramienta, ~60 KB leídos, ~8 KB escritos (más 14 copias en salida/iter-2/ajustes-b).
