# ds-tokens · Pista 3 fase 1 (sistema de diseño) · 11-oct-2026, noche

Agente `ds-tokens` (Sonnet). Alcance: `docs/design/sistema/project/tokens.json`, `project/fonts/`, `docs/design/tools/ds/c2_a_ds.py`, `docs/design/tools/ds/verificar_ds.py` y este reporte. Nada fuera de ese alcance.

## Hecho
- `c2_a_ds.py` convierte `docs/contratos/C2-tokens.json` 1.0 (DTCG, mapa) al formato lista del tipo Design System (`format.md`, «tokens.json grammar»). No cambia ningún valor; los nombres salen de las variables de `tokens.css` sin «--».
- `tokens.json` (68 KB): `version`, `meta.source = {kind:c2, file, version:1.0, synced:2026-10-11}`, `color` (temas light primero, 266 tokens), `type` (2 fuentes, pila `sans`, 1 grupo «Interfaz» con 6 estilos) y 11 familias de lista: spacing 8, radius 4, shadow 3, border 2, focus 1, size 32, layout 14, icon 3, motion 5, breakpoint 5, viz 8.
- Color: 190 primitivos `color-*` (mismo hex en claro y oscuro; los `*-light`/`*-dark`/`surface-*` de roles y crítica quedan como alias `{color-…-700}` etc.) + 75 semánticos `ui-*` (cada par `ui.X.light`/`ui.X.dark` de C2 es un token con `{light, dark}`, alias al primitivo o a otro `ui-*`) + `ui-overlay-scrim`. Alias DTCG `{a.b.c}` pasan a `{a-b-c}` solo si el destino existe (los 273 alias existen).
- `usage`: la `$description` del token o, si no la trae, la del grupo C2 más cercano; en los primitivos se agrega «Lo usan: …» con los tokens que los referencian por tema (derivado de los propios alias, nada inventado) y, en los pasos de rol, el uso de DESIGN §2.4.
- Tipografía (DESIGN §3.2): `title-1` 32/40 700 · `title-2` 24/32 600 · `title-3` 20/28 600 · `body` 16/24 400 · `small` 14/20 400 · `caption` 12/16 400; espaciado 0px; muestras en español con cifras sintéticas; la nota de cifras tabulares (`font-variant-numeric: tabular-nums`, clase `.money`) va en el `usage` de cada estilo salvo `title-1`.
- Fuentes: copiadas sin modificar a `project/fonts/`: `InterVariable.woff2` (352,240 B) e `InterVariable-Italic.woff2` (387,976 B), pesos `100 900`, estilos normal e italic; mismos nombres, peso y estilo que los `@font-face` de `tokens.css`.
- `verificar_ds.py` (independiente del conversor: lee C2, `tokens.css` y `tokens.json`).

## Verificación
```
$ python3 -I docs/design/tools/ds/c2_a_ds.py
tokens.json: 69926 bytes · color 266 · otras familias 85 · tipos 6 estilos · fuentes 2
$ python3 -I docs/design/tools/ds/verificar_ds.py
  color       C2 190  ->  tokens.json 190   (uno a uno)
  ui          C2 150  ->  tokens.json  75   (2 (light, dark) -> 1 con {light, dark})
  overlay     C2   2  ->  tokens.json   1
  font        C2  18  ->  tokens.json  18   (6 estilos de tamaño+interlínea, 4 pesos, espaciado, pila)
  space 8 · radius 4 · shadow 3 · border 2 · focus 1 · size 32 · layout 14 · icon 3 · motion 5 · breakpoint 5 · viz 8  (iguales)
C2: 445 tokens · tokens.json: 351 tokens en listas (color 266 + otras 85) + 6 estilos de tipo + 1 pila + 2 fuentes
Comparaciones de color contra tokens.css (token × tema): 532
Valores vacíos: 0 · con var(/color-mix(/url(: 0 · fuentes presentes: 2/2
Variables de tokens.css: 372 (cubiertas: 370; copias temáticas declaradas: ui-shadow-1, ui-shadow-2)
RESULTADO: OK · 445 tokens C2 -> 351 en listas de tokens.json + 6 estilos de tipo (18 font.*) + 2 fuentes; 0 fallas
```
- El verificador resuelve alias en las dos listas y compara cada hex (claro: bloque `:root`; oscuro: `:root[data-theme="dark"]`, que además comprueba idéntico al bloque `@media`) con `tokens.css`; compara también cada valor no color (rem a px), las sombras (color-mix a rgba) y el velo (a `#rrggbbaa`).
- Prueba negativa: cambié un hex y un alias a mano en una copia de `tokens.json`; el verificador terminó con `RESULTADO: FALLA · 10 fallas`; regeneré y quedó idéntico (`cmp`) y en OK.
- `grep -c 'var(\|color-mix(' tokens.json` = 0. Grep de cifras canon (P3-12) sobre `docs/design/sistema/project` = 0 coincidencias.
- C2 trae 445 `$value` (la ficha del encargo decía 445; mi primer recuento dio 444 porque `icon.stroke.$root` se salta si se ignoran las claves `$`; el conversor y el verificador lo cuentan).

## Conversiones que el formato obliga (declaradas, ningún valor inventado)
1. **Velo** `overlay.scrim` (neutral.950 al 40 % claro / 60 % oscuro): el tipo no admite `color-mix()` y los demás temas no admiten valores por tema fuera de color y sombra; queda como color `ui-overlay-scrim` = `#22232566` / `#22232599` (alfa exacta 0.4 y 0.6 × 255, redondeada a 66 y 99).
2. **Sombras** `shadow-1/2`: capas de `color-mix` de `tokens.css` pasan a `rgba(34, 35, 37, α)` con las alfas de C2; oscuro = `none` (lo dice `finanzas.dark` de C2 y DESIGN §4.5). Diferencia de forma con `tokens.css`: allí `--shadow-1/2` valen siempre el valor claro y el oscuro lo da `--ui-shadow-1/2`; aquí `shadow-1/2` ya trae las dos temas y no se crean `ui-shadow-1/2` (el verificador los declara como copias).
3. **`viz-series-N-dash`**: C2 guarda `solid` o `dashArray` con `px`; se escribe como en `tokens.css` (`none`, `8 4`, …). La marca y el patrón de cada serie (circle/square…, stripes/dots…) van en el `usage`, porque el tipo no los admite como valor.
4. **Alias en medidas** (`size-control-coarse`, `size-bottombar-item-height`, `layout-toast-max-width`): el tipo solo admite alias de color; se resuelve al valor literal (44px, 56px, 480px) y el `usage` dice «Alias de … en C2».
5. **Unidades en `%`**: `layout-sidebar-mobile-width` es `85` con `finanzas.unit: "%"` en C2; se escribe `85%` (lo detectó el verificador).
6. **Valores en `px`** donde `tokens.css` usa `rem` (0.25rem = 4px, etc.): mismo valor; el verificador convierte para comparar. `space-0` queda `0px`.
7. `font.*` no son tokens de lista: tamaño, interlínea, peso y espaciado van dentro de los estilos de `type`; la pila va en `type.families.sans` (igual a `--font-sans`, con comillas como en `tokens.css`).

## Supuestos
- Los primitivos de color llevan el mismo hex en `light` y `dark` (en C2 no cambian con el tema); se escriben los dos temas por claridad.
- Familias: `focus` se deja separada de `border` como pidió el encargo (11 familias de lista, bajo el tope de 12). La clave `meta` y la familia `motion` son claves extra que el tipo conserva con nota; `motion` no tiene vista propia («No motion family»), por eso `duration-*` y `ease-*` solo se verán como sección genérica.
- `version: 1` como en `demo.json`; no hay `name` en la raíz (no está en la gramática): el nombre «Finanzas Hogar» va en `meta.name`.
- Los nombres de estilo de tipo son los de C2 (`title-1`…`caption`); un estilo `caption` a 500 (etiquetas) no se crea porque DESIGN §3.2 lo describe como excepción de peso y no como estilo propio; queda en el `usage`.
- Muestras de los estilos: cifras sintéticas (`$4,826.37`, `−$1,037.40`) y fechas inventadas; aparecen Axel/Ximena/Karen solo como en K45.1 (Ximena y Karen en dos muestras).
- Solo se copiaron los `woff2` (los que cargan las vistas previas, 340 y 380 KB). Los `ttf` (880 y 910 KB, también ≤ 1 MB) no se copiaron: DESIGN §3.1 los reserva para PDF locales; si el publicador los quiere, basta añadirlos al conversor.

## Bloqueos
- Ninguno que detenga la tarea. No se creó archivo en `docs/estado/bloqueos/`.

## Pendientes (para otros agentes / el orquestador)
- Falta `project/design-system.json` (índice) y el README, componentes y `Cover`: no eran de este alcance.
- Falta la licencia OFL de Inter junto a las fuentes (`docs/design/fuentes/inter/LICENSE.txt`, 4,380 B). El tipo no tiene sitio claro para un `.txt` dentro de `fonts/`; decidir si va en `assets/` o se cita en el README (la OFL pide acompañar la redistribución).
- Comprobar, ya publicado, que la Typography view lee los 6 estilos y que `ui-overlay-scrim` se pinta con alfa (el formato lo permite: `#rrggbbaa`).
- `INDICE.md`/`pista3-plan.md` §1.1 dicen `docs/design/tools/tokens/c2_a_ds.py`; los scripts quedaron en `docs/design/tools/ds/` por el encargo (actualizar el plan).

## Consumo aproximado
- Leído ≈ 135 KB (format.md 16 KB, FUENTES.md 3 KB, demo.json extraído por script ≈ 6 KB, DESIGN §2–§4 ≈ 28 KB, C2/tokens.css por script y por trozos ≈ 70 KB, criterios/plan ≈ 12 KB).
- Escrito ≈ 105 KB de texto (tokens.json 68 KB, scripts 35 KB, este reporte 8 KB) + 740 KB de fuentes copiadas.
- ≈ 22 acciones de herramienta.
