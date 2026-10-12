# revisor-presupuesto-r2 · Pista 3, fase 2 (lienzo), vuelta 2 de 3

Fecha: 12-oct-2026 · Agente: `revisor-presupuesto-r2` (revisor independiente; no corrige) · Área: `presupuesto`

## Hecho

- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/presupuesto/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/presupuesto.json`, después de `correccion-presupuesto-r1`.
  - Los comparé con P3 1.3 (P3-07, P3-08, P3-12), con U5 (U5-04, U5-06, U5-07) y con DESIGN.md: §7.3, §7.7, §8, §9.3 y §10.5.
  - También con U1 (§2.1, §2.2, §2.3, §6.4, N-090 a N-113, N-318 y N-319), con U4 (P-26.02) y con `format.md`.
- Hice un diff contra la copia de la vuelta 1 en `docs/design/salida/iter-1/presupuesto/`. P-06 no cambia. P-07, P-08 y P-09 solo cambian en lo que se pidió.
- Informe P3-07: `docs/estado/wfp3/P3-iter-2-presupuesto.md`.
- **Veredicto: NO APROBADO.** Hay 1 alto, 0 medios y 15 bajos.
  - A1, M1, M2, M3 y M4 de la vuelta 1: **cerrados**.
  - **N1 (alto, nuevo).** En P-07, la línea a guiones del devengado acumulado (`ui-chart-2`) pasa sobre el relleno de la capa 1. Ahí da **2.31:1 en claro**, bajo el 3:1 de UI. La agregó la corrección de M1 (supuesto 2).
  - Además, el tramo inferior de la línea «Hoy» da 2.55:1 sobre ese mismo relleno.
  - Hay dos formas de corregirlo: un halo en `ui-surface-card` debajo de las dos líneas, o quitar la línea auxiliar, que §8.3 no pide.
  - Bajos: siguen abiertos B1 a B9 de la vuelta 1 (B9 en parte) y hay 6 nuevos (N2 a N7).

## Verificación

Los scripts son propios, con `python3 -I` y un ejecutor en node para `renderVals()`. Estaban en `datos/tmp/rev-pres-r2/` (`v1.py`, `v2.py`, `v3.py`, `v4.py`, `c1.py`, `c2.py`, `t1.py` y `run.js`) y los borré al terminar.

**Forma:**
- Las líneas obligatorias están exactas en los 8 archivos.
- `data-props` es JSON válido y `$preview` es igual a `w`/`h` del fragmento.
- `html.parser` no encuentra etiquetas sin cerrar.
- Todos tienen `lang="es-MX"` y título propio.
- Hay 90 stems en el lienzo, sin duplicados.
- No quedan huecos sin resolver después de ejecutar `renderVals()`.
- `sc-for` y `sc-if` llevan `hint-*`.
- Raíz móvil de 412 px. Escritorio fluido, con `expand:"fill"` solo en los 4 de escritorio.

**Tokens:**
- 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`.
- 0 colores con nombre.
- Todos los `var(--x)` existen en `tokens.json` y `tokens.css`.

**Contraste:**
- El motor de cascada sacó 43 pares de texto, icono, borde y hover.
- Ninguno queda bajo el umbral, salvo los divisores y la retícula, que DESIGN §7.1 declara decorativos.
- Los pares de la gráfica los calculé a mano con `color-mix` en sRGB y de ahí sale N1.

**P3-12:**
- El grep de las 21 cifras canon da 0, el de los centavos 0 y el de credenciales 0.
- De las 15 fechas que comparten el modelo real y los datos sintéticos, solo `2026-10-30` aparece. Está en P-08 y va con $14,780.00, que es sintético.
- Los datos incrustados son subconjuntos exactos de `datos-sinteticos.json`.
- Revisé una muestra de 12 valores y todos coinciden.

**Accesibilidad:**
- Todos los inputs tienen label (las opciones radio, por `<label>` envolvente).
- Todos los botones tienen nombre y los SVG llevan `aria-hidden`.
- No hay controles bajo 44 px en móvil.
- 0 emoji, 0 palabras en inglés y 0 fechas ISO visibles.
- El dinero va como `$1,234.56`, tabular, y los negativos llevan U+2212.

**Enlaces:**
- Los `href` estáticos y los de `renderVals()` apuntan a archivos de la lista del lienzo que existen y tienen el mismo ancho, o a `#`.

## Supuestos

1. Los divisores (`ui-border-divider`) y la retícula (`ui-viz-grid`) no cuentan para el 3:1, porque DESIGN §7.1 los declara decorativos.
2. Aplico la regla de §8.6 para áreas: relleno al 20 % y la serie en color pleno. Por eso no cuento como hallazgo ni el relleno de la capa 1 contra la tarjeta (1.36:1, porque el área se lee por su línea) ni las rayas de la capa 2 sobre su propio relleno.
3. N1 es alto por la regla del encargo («contraste bajo umbral»), aunque el dato también está en texto y en la tabla. P3-07 pide 0 pares bajo 3:1 en UI, sin excepción por tener una tabla alterna.
4. `./support.js` y `../ds/finanzas-hogar/tokens.css` los pone el orquestador al publicar. No existen bajo `lienzo/` en disco.
5. Las alturas `$preview` no se pueden medir sin navegador (P3-04), así que no las juzgo (B4).

## Bloqueos

Ninguno.

## Pendientes

- Vuelta 3: corregir N1 en `presupuesto-p07-flujo-{m,e}.dc.html` y repetir el contraste de la gráfica en claro y oscuro.
- El orquestador copia el informe a `docs/design/salida/iter-2/` (P3-06).
- Al publicar: medir las alturas (B4) y confirmar `checked` (B2), `color-mix`, el patrón SVG y los huecos en `d=""` y `style=""`.

## Consumo aproximado

- Leído: ~190 KB. Incluye diffs, el informe de la vuelta 1, `format.md`, DESIGN §7.3, §8 y §9.3, filas de U1, `datos-sinteticos.json` por script y fragmentos de los artboards.
- Escrito: ~14 KB de informe y reporte, más ~22 KB de scripts temporales, ya borrados.
- Acciones de herramienta: ~38.
