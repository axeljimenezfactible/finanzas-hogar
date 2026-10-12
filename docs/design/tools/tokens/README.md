# Compilación de los tokens C2 (Style Dictionary 5.6.1)

Contrato **C2 v1.0** (11-oct-2026). `docs/contratos/C2-tokens.json` (DTCG) es la única fuente de valores de diseño para quien lo consume (frontend, PDF de reportes, Pista 3 y la tabla de contraste de componentes). De ella se compilan, sin servicios externos, los dos archivos que consume el frontend:

| Archivo generado | Qué es |
|---|---|
| `docs/contratos/tokens.css` | `@font-face` de Inter autoalojada, variables `--color-*` (primitivos), `--ui-*` (semánticos por tema), `--font-*`, `--space-*`, `--radius-*`, `--shadow-*`, `--duration-*`, `--ease-*`, `--size-*`, `--layout-*`; tema claro por defecto, oscuro con `prefers-color-scheme` y con `data-theme="dark"`; `prefers-reduced-motion` |
| `docs/contratos/tailwind.theme.css` | Tema de Tailwind 4 (`@theme inline`): solo semánticos como utilidades (`bg-surface-page`, `text-role-income`…), escalas cerradas de texto, radios, sombras y quiebres |

## Cadena única: de dónde sale cada valor

```
DESIGN.md §3–§10.6 + C2-colores.json → generar_c2_tokens.py → C2-tokens.json → Style Dictionary → tokens.css y tailwind.theme.css
```

- **Nadie escribe un valor a mano en `C2-tokens.json` ni en los dos CSS.** `npm run check` empieza con `generate`, que regenera `C2-tokens.json`: un cambio hecho a mano ahí se pierde.
- Un **color nuevo** entra por `docs/design/tools/colorimetria.py`, que regenera `C2-colores.json`; su contraste de paleta sale en `docs/design/contraste.csv`.
- Un **semántico `ui.*`** se cambia en la tabla de DESIGN §7.1, que el generador lee.
- Un **valor no de color** se cambia en DESIGN (§3–§6, §9.3–§9.4, §10.6) y en `generar_c2_tokens.py`, que lo transcribe; `verificar_c2.py` comprueba que todo nombre citado en DESIGN exista en C2 como token (no como grupo).
- Después, siempre: `npm run check` y `python3 -I docs/design/tools/contraste_componentes.py --check`. Si cambió un par o un color, se corre sin `--check` para reescribir la tabla de componentes.

**Tablas oficiales de contraste:** `docs/design/contraste.csv` para los pares de paleta (`colorimetria.py`) y `docs/design/contraste-componentes.csv` para los pares de componente de DESIGN §10.5 (`docs/design/tools/contraste_componentes.py`, que lee `C2-tokens.json` como única fuente).

## Requisitos
Node 20 o superior (probado con 26.5.0) y npm. Python 3 (solo biblioteca estándar) para el generador y la verificación. Versiones fijas (sin `^`): `style-dictionary` 5.6.1; `tailwindcss` y `@tailwindcss/cli` 4.3.3 (solo para la prueba C2-07).

## Comandos (desde `docs/design/tools/tokens/`)

```sh
npm ci                  # instala exactamente lo del package-lock.json (node_modules está ignorado por git)
npm run build           # compila C2-tokens.json -> docs/contratos/tokens.css y tailwind.theme.css
npm run verify          # python3 -I verificar_c2.py: C2-02, 03, 04 (con grupos citados como token), 05 (pares básicos), 06, 08 y 10 (con colores repetidos sin alias)
npm run check:sync      # C2-09: compila dos veces a datos/tmp/c2-sync/ y compara byte a byte con los versionados
npm run test:tailwind   # C2-07: compila prueba/entrada.css con Tailwind 4 a datos/tmp/c2-prueba.css (bórralo al terminar)
npm run generate        # regenera C2-tokens.json (generar_c2_tokens.py) desde C2-colores.json, DESIGN.md §7.1 y los valores no de color de DESIGN que transcribe
npm run check           # generate + build + verify + check:sync + test:tailwind
npx style-dictionary --version   # 5.6.1
```

Equivalente directo (C2-01), desde cualquier carpeta: `docs/design/tools/tokens/node_modules/.bin/style-dictionary build --config docs/design/tools/tokens/config.mjs`. La configuración tiene `log.warnings: 'error'` y `brokenReferences: 'throw'`: una referencia rota o una advertencia detiene la compilación.

## Archivos
| Archivo | Para qué |
|---|---|
| `config.mjs` | Configuración de Style Dictionary: fuente, plataforma, formatos propios y registro |
| `formatos.mjs` | Los dos formatos (`finanzas/css`, `finanzas/tailwind`); no contienen valores de diseño |
| `generar_c2_tokens.py` | Arma `C2-tokens.json`: copia los 190 primitivos sin cambiar un hex (`light`, `dark` y `surface-*` de los 8 roles y la acción crítica, como alias de los pasos 700, 400, 100 y 900), lee los 41 semánticos de DESIGN §7.1 (y verifica sus hex) y transcribe los tokens no de color de §3–§6, §9.3–§9.4 y §10.6 |
| `verificar_c2.py` | Verificaciones de C2 (formato, cobertura, igualdad con DESIGN.md y ningún grupo citado como token, contraste básico, CSS, autoalojado, duplicados de color sin alias y lista informativa de medidas iguales) |
| `sincronia.mjs` | Determinismo y sincronía (C2-09) |
| `prueba/` | Hoja de entrada y HTML de prueba de Tailwind (C2-07) |

## Cómo se usa en el frontend (WF-3)
```css
@import "tailwindcss";
@import "<ruta>/docs/contratos/tailwind.theme.css";   /* trae tokens.css */
```
- Las rutas de `@font-face` son relativas a `docs/contratos/` (`../design/fuentes/inter/web/…`); la PWA copia los `woff2` a `/fonts/` y ajusta la ruta en su compilación. Si `tokens.css` se copia a otra carpeta (p. ej. el paquete de Pista 3), hay que copiar también las fuentes o cambiar esas dos rutas.
- Tema: la raíz lleva `data-theme="light"` o `"dark"` solo cuando el interruptor fija el tema; sin atributo manda el sistema.
- Dinero: clase `.money` (o `tabular-nums` de Tailwind) para cifras tabulares (DESIGN §3.3).

## Cambios del contrato
- **C2 v1.0 (11-oct-2026):** primera versión. 190 primitivos de color (colorimetría 1.2: K26, K41, K42) sin cambios de valor; 41 semánticos `ui.*` de DESIGN §7.1 más 34 pares de rol, crítico, persona, tienda y serie y el velo; tipografía, espaciado, radios, bordes, foco, sombras, medidas de control y de estructura, iconos, movimiento, puntos de quiebre y trazos de las series.
- **C2 v1.0 · corrección de la revisión U2 r1 (11-oct-2026, antes de la puerta 2; ningún hex cambia):** `light`, `dark` y `surface-*` de los 8 roles y de la acción crítica pasan a alias de los pasos 700, 400, 100 y 900 (C2-10; en `tokens.css` siguen valiendo el mismo hex). `size.bottombar.item` queda como `height` (alias de `layout.bottombar.height`, 56px) y `min-width` 64px, en lugar de `min-width` 56px y `min-height` 64px (DESIGN §4.6 y §9.3); en CSS, `--size-bottombar-item-min-height` deja de existir y aparece `--size-bottombar-item-height`.
