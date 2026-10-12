# Bloqueo S7-lualatex · lualatex sin luaotfload (fuentes OpenType)

Estado: ABIERTO · 12-oct-2026 · agente s7-pdf (WF-3a) · PROVISIONAL hasta respuesta de Axel.

## Qué falta
- `lualatex` (TeX Live 2026, LuaHBTeX 1.24.0) **no carga fuentes OpenType**: falta el módulo `luaotfload` (paquete de Arch `texlive-luatex`, no instalado). Salida del documento mínimo con `fontspec` + Latin Modern:
  - `[\directlua]:1: module 'luaotfload-main' not found` · `Error in luaotfload: reverting to OT1`
  - `! Font \TU/lmr/m/n/10=[lmroman10-regular]:+tlig; at 10pt not loadable: metric data not found or bad.` · `Fatal error occurred, no output PDF file produced!`
- `luaotfload-tool --update` no se pudo correr: el comando no existe (viene en el mismo paquete). El agente solo tiene permiso de `sudo pacman -S --noconfirm mkcert` (decisión 64 y G), así que no instaló nada.

## Qué se hizo mientras (PROVISIONAL)
- Las plantillas (`backend/internal/reportes/plantillas/*.tex`) compilan **con lualatex** usando **Inter** (la familia de C2, SIL OFL 1.1) en **Type 1** con `\usepackage[T1]{fontenc}` y `\renewcommand*\familydefault{Inter-TLF}` (cifras tabulares), instalada localmente por `texlive-fontsextra`. Sin red. S7-01 pasa: `latexmk -lualatex` = 0, 0 líneas «!», 0 «undefined».
- `¿ ¡ — – · « »` se convierten a comandos LaTeX en `reportes.Tex` porque T1 no los tiene en la posición Unicode.

## Opciones
- **A (recomendada)**: Axel instala `sudo pacman -S texlive-luatex` (≈ 15 MB, sin red en ejecución); se cambia el preámbulo a `fontspec` + Inter OpenType (mejor interletrado y Unicode completo). Costo: 1 cambio de 4 líneas en `comun.tex`.
- **B**: quedarse con Inter Type 1 (funciona hoy; T1 cubre español de México).

## Impacto
Ninguno sobre cifras ni criterios: S7-01, S7-12 y S7-13 pasan con B. Solo cambia la calidad tipográfica.
