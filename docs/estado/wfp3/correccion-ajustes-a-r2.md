# correccion-ajustes-a-r2 · Pista 3 fase 2 · vuelta 3

Fecha: 12-oct-2026. Agente: correccion-ajustes-a-r2. Alcance: artboards de ajustes-a (indice y s01 a s07, m y e), fragmento `ajustes-a.json` y este reporte.

## Hecho
- Copia de la vuelta 2 (P3-06) en `docs/design/salida/iter-2/ajustes-a/` (30 archivos: los .dc.html de ajustes y el fragmento) antes de tocar nada.
- M-01 (diálogos críticos): en s02-m, s02-e, s05-m y s05-e se añadieron al `<helmet>` las reglas `.btn.c`, `.btn.c:hover` y `.dlg.cr h3` (y en s05 también `.dlg` y `.dlg h3`, copiadas de s01-e) y el `<symbol id="ic-shield-alert">` al sprite. El botón de confirmar usa `ui-action-critical-bg`, el encabezado lleva icono y color `ui-critical`, y el diálogo de s05 tiene superficie, borde, radio y sombra.
- M-02 (borde de error): `.inw.er{border-color:var(--ui-border-error)}` en el `<helmet>` de s01-m, s03-m, s04-m, s05-m y s07-m.
- No se tocó el marcado, ni textos, ni alturas, ni el fragmento (sin cambios de contenido).

## Verificación (script propio, python3 -I, 16 archivos de la parte a)
- Forma: doctype, lang es-MX, charset, support.js, tokens.css, x-dc, helmet, data-dc-script, class Component extends DCLogic, renderVals: 16/16.
- data-props JSON válido y `$preview` = w x h del fragmento: 16/16.
- html.parser: 0 etiquetas sin cerrar. 0 hex fuera de data-props. Todos los var(--x) existen en tokens.css/tokens.json.
- 0 cifras prohibidas de P3-12 (con y sin coma). Sin iframe, img, innerHTML, http(s), fetch.
- Todo `<use href>` tiene su `<symbol>`; `.btn.c`, `.dlg`, `.dlg.cr h3` e `.inw.er` presentes donde el marcado los usa.
- Fragmento: 16 boards = 16 order = 16 archivos en disco.

## Supuestos
- Las reglas se añadieron al final del `<helmet>`; con la especificidad de `.btn.c` y `.dlg.cr h3` el orden no cambia el resultado.
- Alturas siguen estimadas (K45, sin navegador).

## Bloqueos
Ninguno.

## Pendientes (bajos de la vuelta 2, fuera del encargo)
- B-01: `.btn.t` sin regla en s04 m/e y s07-m; `.row`, `.tx`, `.cn` sin regla en s04 m/e.
- B-02: «No se pudo cargar los créditos» en s05-e. B-03: leyenda «Enter…» fuera de campo (s02-e, s07-e). B-04: atajos Ctrl+Shift+L y flechas. B-05: alturas estimadas. B-06: primaria no única en s02-e, s05-e, s06-e.

## Consumo aproximado
~35 KB leídos, ~6 KB escritos de contenido nuevo (parches por script a 9 archivos), 9 acciones de herramienta.
