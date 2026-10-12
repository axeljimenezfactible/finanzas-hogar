# correccion-inicio-r2 (vuelta 3, Pista 3, área Inicio)

Fecha: 12-oct-2026. Alcance: artboards de `docs/design/maquetas/lienzo/project/inicio/`, fragmento `inicio.json` y este reporte.

## Hecho
- Copia de la vuelta 2 (P3-06) en `docs/design/salida/iter-2/inicio/`: los 8 `.dc.html` y `inicio.json` (9 archivos), hecha antes de editar.
- Hallazgo 1 (alta, P3-08 / U3 §3 D-6 / Q2-17 C): D-6 «Metas» muestra el avance con un decimal (ROUND_HALF_UP a décimas: 615000/1450000 = 42.4; 1410000/1800000 = 78.3).
  - `inicio-p01-detalle-m.dc.html` l. 151-152 y 157-158; `inicio-p01-detalle-e.dc.html` l. 190-191 y 196-197.
  - Texto «42.4 % · $6,150.00 de $14,500.00» y «78.3 % · $14,100.00 de $18,000.00»; `aria-valuenow` 42.4 y 78.3; `aria-valuetext` «42.4 % de la meta» y «78.3 % de la meta»; el ancho de la barra pasa a 42.4 % y 78.3 %.
  - La tarjeta 6 del Resumen (m y e) se queda con «42 %» entero, como pide el hallazgo.

## Verificación (python3 -I, /tmp/fh-inicio/verificar.py)
- 8/8 artboards: líneas obligatorias, `data-props` parsea y `$preview` = w x h del fragmento, `html.parser` sin etiquetas sin cerrar, 0 hex fuera de `data-props`, todos los `var(--x)` existen en tokens.css, 0 cifras prohibidas P3-12 (con y sin coma), dinero bien formado, sin emoji/iframe/innerHTML/red, raíz móvil 412 px y escritorio sin ancho fijo, enlaces válidos.
- Fragmento `inicio.json`: válido; `boards` y `order` listan exactamente los 8 archivos (sin cambios). RESULTADO: OK.
- `diff` contra la copia de la vuelta 2: solo difieren los dos `detalle`.

## Supuestos
- El ancho de la barra sigue el valor con decimal (coherente con aria-valuenow).

## Bloqueos
- Ninguno.

## Pendientes
- Orquestador: agregar las décimas (424 y 783, `progress_tenths_pct`) a `datos-sinteticos.json` y corregir `LEEME.md` regla 4 (entero en la tarjeta de Inicio, un decimal en Metas). No toqué esos archivos (fuera de alcance).

## Consumo aproximado
~15 KB leídos, ~2 KB de artboards modificados más este reporte, 7 acciones de herramienta.
