# revisor-comidas-r2 · Pista 3, fase 2 (lienzo), vuelta 2 de 3

Fecha: 12-oct-2026 · Agente: `revisor-comidas-r2` (revisor independiente; no corrige) · Área: `comidas`

## Hecho

- Revisé los 6 artboards de `docs/design/maquetas/lienzo/project/comidas/` (P-10, P-11 y P-12, en móvil y escritorio) y el fragmento `docs/design/maquetas/lienzo/fragmentos/comidas.json`, después de `correccion-comidas-r1`. Los comparé con P3 1.3 (P3-07, P3-08, P3-12), U5 (U5-04, U5-06, U5-07), DESIGN.md (§7.5.1, §7.12, §9.3, §10.5), U1 (§2.1, §2.2, §2.5, §6.2, §6.4, N-120 a N-123, F-2) y `format.md` y `craft.md`.
- Hice un diff contra la copia de la vuelta 1 (`docs/design/salida/iter-1/comidas/`). `p11-e` y `p12-e` no cambian. `p11-m` y `p12-m` solo cambian en la clase `bk` y en los `href`. `p10-m` y `p10-e` suman el selector de recetas y `p10-e` suma la barra de la vista.
- Informe P3-07: `docs/estado/wfp3/P3-iter-2-comidas.md`.
- **Veredicto: aprobado.** 0 altos y 0 medios.
  - M-1, M-2, M-3 y M-4 de la vuelta 1: **cerrados**, con evidencia en el informe.
  - Siguen abiertos 10 bajos de la vuelta 1 (B-1 a B-11 menos B-12, que ya no aplica).
  - Hay 6 bajos nuevos:
    - N-1: el selector pide confirmar con «Poner en el plan», que es un toque más que F-2.
    - N-2: la hoja y el diálogo se apartan de §7.12 (título `title-3`, primaria de 44 px, orden de las acciones).
    - N-3: hay 4 px entre objetivos contiguos, y DESIGN pide 8.
    - N-4: el ejemplo «Lun · Comida» no es una casilla vacía en los datos.
    - N-5: «Filtrar» usa el icono `funnel`, los botones de vista llevan borde y P-10 móvil no tiene «Filtrar».
    - N-6: la clase `g2` no tiene regla.

## Verificación

Hice scripts propios con `python3 -I` en `datos/tmp/rev-comidas-r2/` (`v1.py`, `cas.py`, `dat.py`, `tap.py` y `a11y.py`) y los borré al terminar.
- Forma: las 14 piezas obligatorias están, en orden, en 6 de 6 archivos. `data-props` es JSON válido y `$preview` es igual al fragmento. `html.parser` no encuentra etiquetas sin cerrar. Todos tienen `lang="es-MX"` y título propio. Hay 90 stems en el lienzo, sin duplicados.
- Tokens: 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`. Todos los `var(--x)` existen en `tokens.json` y `tokens.css`. Los 190 SVG llevan `aria-hidden` y `currentColor`.
- Contraste: un motor de cascada propio sacó 190 pares, en claro y en oscuro, más el hover; 0 bajo el umbral. Coteje 4 pares con `contraste-componentes.csv` (C-001, C-058, C-064, C-085) y coinciden.
- P3-12:
  - El grep de las 21 cifras canon da 0; el de centavos, 0; el de credenciales, 0.
  - Ninguna de las 15 fechas que comparten el modelo real y `datos-sinteticos.json` aparece en los artboards.
  - Todos los montos y kcal salen de los datos. Las 28 casillas × 2 personas, las 5 recetas del selector y las 12 tarjetas de P-11 coinciden sin diferencias.
- Accesibilidad: 96 controles móviles medidos con los tokens; 0 bajo 44 px (las `.rn` de P-11 cuentan como tarjeta completa por `::after`). No hay controles falsos, botones sin nombre ni inputs sin label. Hay 0 emoji, 0 palabras en inglés y 0 fechas ISO visibles. Los 116 importes tienen formato `$1,234.56` dentro de `.tn`, y los negativos van con U+2212.
- Enlaces: todos apuntan a la lista del lienzo o a «#». 0 cruzan de ancho.
- Fragmento: lista exactamente los 6 archivos, con `expand: "fill"` solo en escritorio, `is_interactive: false` y sin x ni y.

## Supuestos

1. Los bajos (B-* y N-*) no bloquean, porque el encargo aprueba con 0 altos y 0 medios. N-1 se deja en baja porque todos los elementos de N-122 y N-123 están dibujados y lo único que cambia es un toque por receta.
2. La hoja móvil y el diálogo de escritorio se dibujan en el flujo de la página, con una nota que explica cuándo aparecen. Para un lienzo estático lo acepto como forma válida de mostrar una superposición.
3. P-10 móvil sin «Filtrar» queda en baja (N-5): U1 §6.4 no lo pide para la lista por día y la vuelta 1 lo dio por no aplicable.
4. `./support.js` y `../ds/finanzas-hogar/tokens.css` los pone el orquestador al publicar; en disco todavía no existen bajo `lienzo/project/`.
5. El informe va en `docs/estado/wfp3/`, como pide el encargo.

## Bloqueos

Ninguno.

## Pendientes

- Vuelta 3 (opcional): N-1 a N-6 y los bajos abiertos que el autor quiera; los que más rinden son N-1, N-2, N-3 y B-3.
- Orquestador: copiar `P3-iter-2-comidas.md` a `docs/design/salida/iter-2/` (P3-06), medir los altos de p10-m y p10-e al publicar (B-8) y decidir `is_interactive`.

## Consumo aproximado

- Leído: unos 150 KB (cuerpos de p10-m y p10-e sin SVG, diffs de la vuelta 1, partes de U1, DESIGN, P3 y U5, `format.md` entero, 3 secciones de `craft.md`, reportes de la vuelta 1, y `tokens.json` y `datos-sinteticos.json` por script).
- Escrito: unos 40 KB (scripts temporales, ya borrados, este reporte y el informe).
- Acciones de herramienta: unas 45. Sin Artifact, navegador, red, git ni Telegram.
