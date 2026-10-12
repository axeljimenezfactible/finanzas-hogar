# correccion-shell-r1 (vuelta 2) · Pista 3 fase 2 · área shell

Fecha: 12-oct-2026. Corrige los hallazgos S1-01 (media) y el de enlaces de escritorio (media) de `P3-iter-1-shell.md`.

## Hecho
- Copia de la vuelta 1 (antes de tocar nada): `docs/design/salida/iter-1/shell/` con los 8 artboards (P3-06).
- S1-01, objetivo táctil de los campos en móvil (todos en `docs/design/maquetas/lienzo/project/shell/`):
  - `shell-p03-registrar-gasto-m.dc.html`: `#monto` con `align-self:stretch;min-height:var(--size-control-coarse)`. Llena la caja de 64 px.
  - `shell-p22-registrar-peso-m.dc.html`: `#peso`, mismo cambio.
  - `shell-p00-entrar-m.dc.html`: la regla `.fld input` (afecta a `#usuario` y `#contrasena`) ahora lleva `min-height:var(--size-control-coarse);align-self:stretch`. Cada campo mide 44 px o más.
  - `shell-barra-lateral-m.dc.html`: `#q-m` con `min-height:var(--size-touch-min)` (el `font:inherit` en línea se conserva).
- Enlaces de escritorio de «Ajustes»: 4 href cambiados a `../ajustes/ajustes-p26-s01-cuenta-roles-e.dc.html` (U1 §1.2.3). Están en `shell-barra-lateral-e` (barra lateral y tarjeta de Main) y en las barras de fondo de `shell-p03-registrar-gasto-e` y `shell-p22-registrar-peso-e`. En móvil el enlace al índice (`ajustes-p26-indice-m`) no se tocó.
- Fragmento `fragmentos/shell.json`: sin cambios (mismos 8 archivos, mismas dimensiones).

## Verificación (script propio, `python3 -I`, 8 archivos)
- Líneas obligatorias exactas presentes en los 8 (doctype, `lang="es-MX"`, charset, support.js, tokens.css, x-dc, helmet).
- JSON de `data-props` parsea; `$preview` coincide con `w` × `h` del fragmento; móviles a 412.
- `html.parser`: 0 etiquetas sin cerrar ni desparejadas.
- 0 hex literales en color, fondo, borde, relleno o trazo fuera de `data-props` y SVG.
- 0 cifras de la expresión CANON de `tools/ds/verificar_ds.py` (con y sin coma).
- Todos los `var(--x)` existen en tokens.json, salvo `--font-sans`: está definido en `sistema/project/tokens.css` (7 usos), no figura como entrada con `name` en tokens.json. Es lo mismo que en la vuelta 1.
- 0 enlaces `../` o `shell-` rotos hacia artboards existentes en disco. El `../ds/finanzas-hogar/tokens.css` no existe aún en disco, y lo coloca el orquestador (línea obligatoria).
- Los 4 inputs móviles declaran alto mínimo ≥ 44 px (por estilo en línea o por la regla `.fld input`). 0 href a `ajustes-p26-indice-e` en escritorio.
- Fragmento: JSON válido, lista exactamente los 8 archivos, sin repetidos.

## Supuestos
- Medición táctil por estilo declarado (`min-height` con token de 44 px + `align-self:stretch`), sin navegador (K45). El relleno vertical del input queda dentro de `min-height` por `box-sizing` de border-box del reinicio base. Si el reinicio no lo aplicara, el input mediría más, nunca menos.
- `--font-sans` se considera token válido por existir en tokens.css.

## Bloqueos
Ninguno.

## Pendientes
Ninguno de este alcance. El orquestador debe publicar `ds/finanzas-hogar/tokens.css` y `support.js` junto al lienzo.

## Consumo aproximado
~12 acciones de herramienta, ~75 KB leídos, ~3 KB escritos (más 105 KB de copia de seguridad).
