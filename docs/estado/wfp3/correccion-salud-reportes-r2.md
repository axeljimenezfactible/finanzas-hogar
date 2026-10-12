# Corrección salud-reportes, vuelta 3 (r2)

Agente: correccion-salud-reportes-r2 · 12-oct-2026

## Hecho

- Copia de la vuelta 2 (8 artboards y el fragmento) en `docs/design/salida/iter-2/salud-reportes/` (P3-06), hecha antes de editar.
- Hallazgo 1 (media, P3-08 punto 5 / U5-06 / DESIGN §3.3 y §7.5.2): en `reportes-p25-karen-e.dc.html` las 3 tablas de P-25 escritorio mostraban 22 importes en `<td class="r">` sin cifras tabulares. Se agregó `font-variant-numeric:tabular-nums` a la regla `.tbl td` del helmet (misma solución que P-25 móvil, P-24 y P-21). No se usó `dl.kv dd` porque el artboard no tiene `<dl>`.
- Es el único archivo que cambió respecto a la copia de iter-2 (comprobado con `cmp`).

## Verificación

Script propio (`python3 -I`, /tmp/vsr/v.py) sobre los 8 archivos: RESULTADO OK.
- Forma: doctype, `lang="es-MX"`, support.js, tokens.css, `x-dc`, helmet, `data-dc-script`, `class Component extends DCLogic`, `renderVals`; JSON de data-props parsea; ancho de $preview 412 (móvil) y 1440 (escritorio).
- html.parser: 0 etiquetas sin cerrar ni cierres inesperados; 0 ids duplicados.
- 0 hex literales fuera de data-props (excluidos `href="#…"`), 0 rgb()/hsl(); todos los `var(--x)` existen en `docs/design/sistema/project/tokens.css`.
- P3-12: 0 de las 21 cifras canon (con y sin coma); 0 credenciales; 0 iframe/innerHTML.
- Fragmento `salud-reportes.json`: válido; lista exactamente los 8 archivos; `expand: fill` solo en escritorio.
- Regla `.tbl td` con tabular-nums presente en P-25 escritorio.

## Supuestos

- Los tokens se validaron contra `docs/design/sistema/project/tokens.css` porque `project/ds/finanzas-hogar/tokens.css` lo coloca el orquestador al armar el lienzo.
- Sin render en este flujo; el cambio es solo tipográfico y no altera anchos de columna de forma relevante.

## Bloqueos

Ninguno.

## Pendientes

- Los hallazgos bajos de rondas anteriores siguen sin tocarse (el encargo era alta y media).
- Vistas Tablero y Calendario de P-23 ofrecidas en el selector, sin dibujar.

## Consumo aproximado

~25 KB leídos, ~4 KB escritos, 7 acciones.
