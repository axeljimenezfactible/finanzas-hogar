# U2c-tokens · reporte (WF-1, 11-oct-2026)

Agente Sonnet. Cerrado por tope: no. Tokens DTCG completos (C2 v1.0) y compilación reproducible con Style Dictionary 5.6.1.

## Hecho
- `docs/contratos/C2-tokens.json` (DTCG): 445 tokens. Los 190 primitivos de `C2-colores.json` copiados sin cambiar un valor; 75 semánticos `ui.*` (los 41 de DESIGN §7.1 más 34 pares de rol, crítico, persona, tienda y serie de §10.3) y `overlay.scrim`; tipografía, espaciado, radios, bordes, foco, sombras, medidas, iconos, movimiento, puntos de quiebre y trazos de series.
- `docs/contratos/tokens.css` (25 KB) y `docs/contratos/tailwind.theme.css` (6.7 KB) compilados; claro por defecto, oscuro con `prefers-color-scheme` y `data-theme="dark"`, `prefers-reduced-motion`, `@font-face` de Inter con rutas a `docs/design/fuentes/`.
- `docs/design/tools/tokens/`: `package.json` (versiones fijas), `package-lock.json`, `config.mjs`, `formatos.mjs`, `generar_c2_tokens.py`, `verificar_c2.py`, `sincronia.mjs`, `prueba/`, `README.md` con los comandos, `.gitignore` (node_modules).

## Archivos
docs/contratos/C2-tokens.json · docs/contratos/tokens.css · docs/contratos/tailwind.theme.css · docs/design/tools/tokens/{README.md, package.json, package-lock.json, config.mjs, formatos.mjs, generar_c2_tokens.py, verificar_c2.py, sincronia.mjs, .gitignore, prueba/entrada.css, prueba/prueba.html} · docs/estado/wf1/U2c-tokens.md

## Verificación
- **C2-01**: `npm ci` (limpio, sin node_modules previo) y `npm run build`; también `cd docs/contratos && ../design/tools/tokens/node_modules/.bin/style-dictionary build --config ../design/tools/tokens/config.mjs` → código 0; `log.warnings: 'error'` y `brokenReferences: 'throw'`, 0 advertencias, 0 referencias rotas. `npx style-dictionary --version` → `5.6.1`.
- **C2-02**: `npm run verify`: 445 tokens; 0 sin tipo efectivo, 0 sin `$value`, 0 referencias rotas; por tipo: color 342, dimension 79, duration 3, cubicBezier 2, fontFamily 1, fontWeight 4, number 3, shadow 3, strokeStyle 8.
- **C2-03**: por grupo: color 190, ui 150 hojas (75 pares), font 18, space 8, radius 4, border 2, focus 1, shadow 3, overlay 2, size 32, layout 14, icon 3, motion 5, breakpoint 5, viz 8. Escalas brand y neutral de 11 pasos; 8 roles + crítico; 8 categóricos; 8 series; 76 semánticos de color (ui + overlay) con claro y oscuro, 0 incompletos. Espaciado 4/8/12/16/24/32/48, radios 6/10/14, tipografía 12–32, movimiento 150/200/250 ms, quiebres 360/412/768/1024/1440.
- **C2-04**: 190 colores presentes con el mismo hex (0 diferencias, 0 de más); 158 pares «token · hex» de DESIGN.md comparados, 0 diferencias; 233 tokens citados en DESIGN.md, 0 inexistentes en C2-tokens.json (se aceptan las plantillas `color.chart.N.light` y `viz.series.N.dash`). `generar_c2_tokens.py` verifica además cada hex de la tabla §7.1.
- **C2-05 (parcial)**: 170 combinaciones (par × tema: texto, iconos, enlace, botones, bordes, foco, roles, crítico, personas, tiendas y series contra fondo, tarjeta y superficie propia) calculadas desde C2-tokens.json, con la fórmula WCAG truncada a 2 decimales: 0 bajo umbral. La tabla completa de §10.5 (158 pares × 2 temas) es del script de U2d; no se recalculó aquí.
- **C2-06**: `--ui-*` 78 en claro, 78 en oscuro (@media) y 78 en oscuro (`data-theme`), 0 diferencias entre ellos, y el oscuro de @media es idéntico al de `data-theme`; `prefers-color-scheme: dark` 1, `data-theme` 2, `prefers-reduced-motion` 1; variables `--color-` 190, `--space-` 8, `--radius-` 4, `--shadow-` 3, `--duration-` 3, `--font-` 13; 0 colores sueltos fuera de los primitivos.
- **C2-07**: `npm run test:tailwind` (Tailwind 4.3.3) → código 0, 37 KB; en el CSS resultante: `var(--ui-surface-page)` y `#fafafb` presentes, `0.25rem` (spacing), `0.625rem`/`10px` (radio), `var(--ui-shadow-1)`, `cubic-bezier(0.05…`, `var(--ui-role-income)`, `@media (width >= 48rem)` y `64rem` (quiebres), `any-pointer: coarse`, 2 `@font-face`. El archivo de salida se borró.
- **C2-08**: 0 URL `http(s)://` en tokens.css, tailwind.theme.css y C2-tokens.json; los dos `woff2` y `docs/design/fuentes/inter/LICENSE.txt` existen.
- **C2-09**: `npm run check:sync`: dos compilaciones consecutivas idénticas y iguales a los versionados (sha256 tokens.css `495664490ca32c1ef34afa690bb231d623d91277b4de6304209bad8ae4df1c88`, tailwind.theme.css `f66eb866449f391a8799011515ce0bd82365f04dd1177e38bcac86ded4b9c88a`); tras `rm -rf node_modules && npm ci && npm run check`, `cmp` de C2-tokens.json, tokens.css y tailwind.theme.css contra las copias previas: iguales.
- **C2-10**: los 6 `font.line-height.*` no se citan uno por uno en DESIGN.md (los cubre `font.line-height.*` de §10.2 y la columna «Interlínea» de §3.2); los demás no de color/ui están citados. Duplicados de valor sin alias entre tokens que no son de color: 18 grupos, todos coincidencias de medida entre conceptos distintos (listados por `npm run verify`); los de un mismo concepto sí llevan alias (`size.control.coarse` → `size.touch.min`, `layout.toast.max-width` → `layout.dialog.width.short`, y los `ui.*` repetidos). Ver Supuesto 12. `$description` en español de México en todos los tokens que no son de color; nombres en inglés corto sin mayúsculas.
- **C2-11 (parte propia)**: contrato identificado como C2 v1.0 con sección de cambios (README y `$extensions.finanzas` de la raíz). Que `prompt.md` y las plantillas `.tex` coincidan con C2 lo comprueban sus agentes (ver solicitudes).
- **C2-12**: sección Supuestos abajo; ningún archivo `docs/estado/bloqueos/C2-*` (no hubo valores sin decidir).

## Supuestos
1. Las dimensiones van como cadena (`"12px"`, `"150ms"`), formato clásico que Style Dictionary 5.6.1 compila y que va con los colores en hex de `C2-colores.json`; no se usó el objeto `{value, unit}` de DTCG 2025.10.
2. Alfa de sombras y velo: el color es la referencia `{color.neutral.950}` y el alfa va en `$extensions.finanzas.alpha` (DESIGN §10.2: «con alfa derivado del token»); el CSS usa `color-mix(in srgb, …)` como en §4.5.
3. `font.size.*`, `font.line-height.*` y `space.*` salen en `rem` (px ÷ 16): igual valor con base de 16 px y respetan el tamaño de texto del usuario (zoom al 200 %, §9.1). Los quiebres salen en px en tokens.css (referencia) y en rem en Tailwind (§10.4).
4. `icon.stroke` (1.5 px) usa `$root` de DTCG y `icon.stroke.user-units.{md,sm}` (1.5 y 1.8) cuelga del mismo grupo, para conservar ambos nombres de DESIGN §5.2.
5. DESIGN da dos valores bajo un solo nombre en cuatro tokens; se desdoblaron: `layout.dialog.width.{short,long}` (480, 640), `layout.topbar.height.{mobile,desktop}` (56, 48), `layout.preview.width.{desktop,wide}` (400, 480) y `size.bottombar.item.{min-width,min-height}` (56 × ≥ 64, leído como ancho mínimo 56 y alto mínimo 64).
6. `layout.sidebar.mobile-width` (85 %) es `number` 85 con `unit: "%"` en la extensión (DTCG no tiene dimensión en porcentaje).
7. `shadow.0` es una sombra de una capa nula con alfa 0 que se emite como `none`; las sombras 1 y 2 valen `none` en oscuro (`--ui-shadow-N`) y la superficie usa el borde de `ui.border.divider` (§4.5).
8. `--ui-shadow-N` se emite con el valor literal y no con `var(--shadow-N)`: Tailwind define `--shadow-N: var(--ui-shadow-N)` y la referencia mutua sería circular.
9. El tema de Tailwind se llama `tailwind.theme.css` (nombre de la tarea; DESIGN §10.1 proponía `tailwind-theme.css`). No importa `tailwindcss` (resolvería mal desde `docs/contratos/`): la hoja de entrada lo importa primero y luego el tema, que trae `tokens.css`. Se vaciaron las escalas por defecto (`--color-*`, `--text-*`, `--radius-*`, `--shadow-*`, `--ease-*`, `--breakpoint-*`, `--font-*`, `--font-weight-*`) y se dejaron `transparent` y `currentColor` (palabras clave, no colores).
10. Se agregaron a tokens.css `--font-numeric-tabular` y la clase `.money` (cifras tabulares de §3.3) y los `--viz-series-N-dash` (trazos de §8.6); marcador y patrón de cada serie viven en `$extensions.finanzas` de `viz.series.N.dash`. No cambian ningún valor.
11. `tailwindcss` y `@tailwindcss/cli` 4.3.3 (la versión de plan.md) quedan como devDependencies exactas solo para la prueba C2-07; el script de instalación de `@parcel/watcher` queda bloqueado por npm y no se necesita.
12. Duplicados de valor (C2-10): se tomó «sin alias que los relacione» como misma idea de diseño, no mismo número. Si el revisor exige alias para toda coincidencia de medida (p. ej. 24 px = título-2 = icono md = chip), es una decisión de Axel o del integrador, no un cambio de valores.
13. El generador lee la tabla de DESIGN §7.1 para los semánticos; si DESIGN cambia, se corre `npm run generate` y luego `npm run check`.

## Bloqueos
Ninguno propio. El bloqueo abierto único `docs/estado/bloqueos/S1-parametros-avisos.md` no toca estos tokens.

## Pendientes
Ninguno de mi alcance. La tabla de contraste completa (158 pares × 2 temas) es de U2d (`docs/design/tools/contraste_componentes.py`) y debe leer `C2-tokens.json`: las hojas `ui.*.light` y `ui.*.dark` son referencias en cadena (`{ui.text.secondary.light}` → `{color…}`).

## Solicitudes fuera de alcance
1. `docs/design/DESIGN.md` §10.1–§10.4 y §4.3, §4.6: reflejar los nombres desdoblados (Supuesto 5), `tailwind.theme.css` en lugar de `tailwind-theme.css`, que el tema ya no importa `tailwindcss` (Supuesto 9) y que `ui.*` tiene 75 pares y no 41 (§10.2). Motivo: C2-04, que DESIGN y C2 coincidan.
2. `docs/design/entrada/` (agente paquete-claude-design): copiar `C2-tokens.json` y `tokens.css` junto con `docs/design/fuentes/inter/web/*.woff2` o ajustar las dos rutas de `@font-face` (son relativas a `docs/contratos/`). Motivo: tokens.css copiado suelto apuntaría a archivos inexistentes.
3. Agentes de `prompt.md` (P3-02) y plantillas `.tex` (S7-13): tomar los hex de `C2-tokens.json` (C2-11).
4. `.gitignore` raíz: sin cambio; basta el `.gitignore` local de `docs/design/tools/tokens/` (`node_modules/`).

## Paquetes instalados
Del sistema: ninguno. De npm, solo en `docs/design/tools/tokens/node_modules` (ignorado): style-dictionary 5.6.1, tailwindcss 4.3.3, @tailwindcss/cli 4.3.3 (y sus dependencias, fijadas en `package-lock.json`).

## Propuestas para Axel
Ninguna: no se necesitó ningún valor operativo de K31.

## Consumo aproximado
36 acciones de herramienta; lectura ≈ 100 KB (protocolo, §K, C2.md, DESIGN §3–§6, §7.1, §8.6, §9.4–§10, extracto de C2-colores.json); escritura ≈ 62 KB escritos a mano (generador 19.5, verificación 13, formatos 10, README 4, reporte, config) más 131 KB generados por herramientas (C2-tokens.json 99, tokens.css 25, tailwind.theme.css 7).
