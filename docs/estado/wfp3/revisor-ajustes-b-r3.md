# revisor-ajustes-b-r3 · vuelta 3 de 3 · 12-oct-2026

Agente `revisor-ajustes-b-r3`, revisor independiente que no corrige. Alcance: los 14 artboards de ajustes-b (P-26.08 a P-26.14, `-m` y `-e`) y `fragmentos/ajustes-b.json`. Informe P3-07: `docs/estado/wfp3/P3-iter-3-ajustes-b.md`.

## Hecho

- Comprobé los 3 hallazgos de la vuelta 2. R2-A1 (opacidad en la nota de s13) está **cerrado**. R2-M1 (7 montos sin `tabular-nums`) está **cerrado**. H-11 (K7 en s14) sigue **abierto**, porque `datos-sinteticos.json` no cambió desde el 11-oct 23:55.
- Revisé desde cero la forma, los tokens, el contraste (claro y oscuro), P3-12, la fidelidad a U1/U4, la accesibilidad, la disposición y los enlaces. Para eso escribí un motor propio de cascada CSS y tokens.
- Comparé contra la copia de la vuelta 2 (`docs/design/salida/iter-2/ajustes-b/`). Solo cambiaron los 5 artboards del alcance de la corrección y la altura de s13-m, y no hay regresiones.
- Veredicto: **NO APROBADO**, con 1 alta, 1 media y 18 bajas.
  - **R3-A1 (alta, era B-07):** 49 marcadores de posición visibles sin regla `::placeholder`. Toman el gris del navegador y quedan en 4.42:1 (Chrome), 3.05:1 (Firefox) y 2.25:1 (Safari) en claro, y en 3.41:1 en oscuro con Chrome. DESIGN C-035 pide `ui-text-secondary` (5.40 / 10.43).
  - **H-11 (media):** sigue abierta.
  - **Bajas nuevas:** B-20 (Ximena muestra «0» propio, pero la pantalla dice que deja vacío su ajuste; falta «Usar el de la casa») y B-21 (iconos de 16 px en insignias y literales fuera de la escala de 4 pt).

## Verificación

Scripts propios con `python3 -I` en `datos/tmp/rev-ajb-r3/` (`eng.py`, `chk.py`, `vis.py`), borrados al final.

- Forma 14/14. `data-props` y `$preview` coinciden con el fragmento. `html.parser` da 0 errores. Stems únicos en los 90 artboards.
- 0 hex o funciones de color fuera de `data-props`. Todos los `var()` existen; `--font-sans` sale de `type.families.sans`.
- Contraste: 1,629 nodos por tema y 0 pares bajo umbral por cascada. Se recalcularon 10 pares de texto y 7 de UI. Los marcadores de posición quedan bajo umbral (R3-A1).
- P3-12: 0 cifras canon (con y sin coma, y en centavos) y 0 credenciales. Las 2 fechas que coinciden con el calendario real van sin monto. 65/65 montos son del JSON o derivados exactos, y la muestra de 10 valores coincide con sus rutas.
- 143 objetivos táctiles móviles, todos ≥ 44 px. 0 emoji. 0 palabras en inglés (quedan las de B-19). 65/65 montos con `tabular-nums`.
- 242 enlaces, 0 rotos y 0 de otro ancho. Fragmento válido, con `expand: "fill"` en los 7 escritorios.

## Supuestos

- El color del marcador de posición sale de las hojas de estilo por defecto de cada motor: Chrome `#757575` en claro, Firefox el texto al 54 % y Safari `darkgray`. Las páginas no declaran `color-scheme`, así que en Chrome oscuro también queda `#757575`.
- Mantengo H-11 en media, igual que las vueltas 1 y 2. Es una incoherencia de datos que se ve en Ajustes, y su corrección es del JSON.
- Las alturas del fragmento no las volví a medir, porque no uso navegador. El +164 px de s13-m cuadra con dos campos apilados.

## Bloqueos

- H-11 depende de `docs/design/maquetas/datos-sinteticos.json`, que está fuera del alcance de ajustes-b, y no hay archivo en `docs/estado/bloqueos/`. Lo tiene que encargar el orquestador a correccion-datos, o registrarlo y marcarlo «no adoptar».

## Pendientes

- Para el corrector: R3-A1 (una regla CSS en los 14) y las bajas baratas B-03, B-05, B-12, B-13, B-14, B-20 y B-21.
- Para el orquestador o correccion-datos: H-11, B-11, B-15, B-16 y B-19.
- La copia de este informe a `docs/design/salida/iter-3/` (P3-06) es tarea del orquestador.

## Consumo aproximado

~50 acciones de herramienta, ~170 KB leídos (U4 por secciones, §10.5 de DESIGN, partes de `format.md`, P3, U5, los dos reportes previos y extractos de los artboards) y ~22 KB escritos (los dos informes y los scripts temporales). 0 acciones de navegador y 0 llamadas a Artifact.
