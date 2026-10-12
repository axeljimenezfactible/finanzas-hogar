# ds-tokens-css — reporte

## Hecho
- `docs/design/tools/ds/tokens_css.py`: generador + verificador (`python3 -I tokens_css.py` escribe; `--check` solo verifica contra disco).
- `docs/design/sistema/project/tokens.css` (27.3 KB) generado desde `project/tokens.json` v1 (C2 1.0), formato «tokens.css as compiled»: línea 1 `/* Finanzas Hogar — generated from tokens.json */`; `:root, [data-theme="light"]` (266 colores + 3 sombras); `[data-theme="dark"]` (los mismos 269, alias re-declarados); `:root` con space, radius, font-weight, border, focus, size, layout, icon, motion, breakpoint, viz y `--font-sans` (87 variables); 6 clases (`.title-1 .title-2 .title-3 .body .small .caption`); 2 `@font-face` (Inter normal e italic, `src: url("fonts/…") format("woff2")`).
- Los alias `{x}` de tokens.json salen como `var(--x)`; los nombres no llevan «.» (el escape `\.` está implementado pero no se ejerce).

## Verificación
- Balance de llaves, sin texto fuera de bloques ni anidamiento; selectores iniciales exactos.
- Cada uno de los 266 colores y 3 sombras en los dos temas con el valor esperado (0 diferencias).
- Cada estilo tipográfico con su clase; @font-face = type.fonts (2).
- Conteo: 269 variables por tema + 87 en :root = 356 = 355 tokens de tokens.json + `--font-sans`.
- Todo `var(--x)` usado en el CSS existe. `--check` devuelve rc=0.

## Nombres más usados (por menciones en DESIGN.md; todos existen en tokens.css)
- Texto: `--ui-text-default`, `--ui-text-secondary`, `--ui-text-link`, `--ui-text-link-hover`, `--ui-text-disabled`, `--ui-text-on-primary`, `--ui-text-on-destructive`, `--ui-text-on-critical`.
- Fondo / superficie: `--ui-surface-page` (fondo de página), `--ui-surface-card` (tarjeta), `--ui-surface-raised`, `--ui-surface-field`, `--ui-surface-selected`, `--ui-surface-hover`, `--ui-surface-hover-subtle`, `--ui-surface-pressed`, `--ui-surface-disabled`, `--ui-skeleton-on-card`, `--ui-skeleton-on-page`, `--ui-overlay-scrim`.
- Marca / acción: `--ui-action-primary-bg` (+ `-hover`, `-pressed`), `--ui-action-destructive-bg`, `--ui-action-critical-bg` (+ `-hover`, `-pressed`), `--ui-icon-default`, `--ui-icon-selected`; primitivos `--color-brand-50..950`.
- Bordes: `--ui-border-control` (campos), `--ui-border-divider` (divisores y borde de tarjeta), `--ui-border-selected`, `--ui-border-control-hover`, `--ui-border-error`; grosores `--border-width-default` (1px), `--border-width-focus` (2px).
- Foco: `--ui-border-focus` (anillo), `--border-width-focus`, `--focus-offset` (2px).
- Roles y personas: `--ui-role-income|expense|debt|savings|success|warning|error|info` (+ `-surface`), `--ui-critical` (+ `-surface`), `--ui-person-axel|ximena|karen|shared`, `--ui-store-*`, `--ui-chart-1..8`, `--ui-viz-baseline|track|grid|meter|other`.
- Medidas: `--space-0|1|2|3|4|6|8|12`, `--radius-*` (sm 6px, md 10px…), `--size-touch-min` (44px), `--layout-*`, `--shadow-0|1|2`, `--font-sans`, `--font-weight-regular|medium|semibold|bold`, `--duration-fast|base`. Clases tipográficas: `.title-1/2/3 .body .small .caption`.

## Supuestos
- `--font-<clave>` se emite con la clave de `type.families` (`--font-sans`); las clases usan `var(--font-<family del grupo>)`.
- El orden de familias en :root es el de tokens.json salvo `fontWeight`, que va tras radius (después de shadow no aplica: las sombras viven en los bloques de tema).
- Las fuentes se resuelven con la ruta relativa `fonts/…` desde `project/`; los artboards enlazan `../ds/finanzas-hogar/tokens.css`, así que el orquestador debe copiar/publicar `fonts/` junto al CSS en esa carpeta (o las fuentes caen a la pila de respaldo del sistema).

## Bloqueos
- Ninguno.

## Pendientes
- Orquestador: publicar `tokens.css` + `fonts/` en `ds/finanzas-hogar/`. Regenerar con el script si cambia tokens.json.

## Consumo aproximado
- ~140 KB leídos, ~13 KB escritos, 8 acciones de herramienta.
