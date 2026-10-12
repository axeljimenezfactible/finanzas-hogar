# revisor-gastos-r1: reporte (Pista 3, fase 2 «lienzo», vuelta 1 de 3, 12-oct-2026)

## Hecho
- Revisé los 6 artboards de `docs/design/maquetas/lienzo/project/gastos/` (P-02, P-04 y P-05, en móvil 412 y escritorio 1440) y el fragmento `fragmentos/gastos.json`. Usé como referencia P3 1.3 (P3-07, P3-08, P3-12), U5-04, U5-06 y U5-07, DESIGN.md (§3.3, §7.2, §7.5, §7.14, §9.2 y §9.3), U1 (§2.1, §2.2, §2.4, §2.6 y §3.4 a §3.7) y `referencia-tipo/format.md`.
- El informe está en `docs/estado/wfp3/P3-iter-1-gastos.md`. Trae resumen de hallazgos, verificaciones (a) a (h), recálculo de 12 pares de contraste en los dos temas, tabla de fidelidad contra U1, la lista de 8 puntos de P3-08, la tabla artboard × regla, correcciones para la vuelta 2 y las líneas P3-04 y P3-12.
- **Veredicto: NO APROBADO.** Hay 1 hallazgo alto, 6 medios y 5 bajos:
  - Alto: placeholders sin token, bajo 4.5:1 en claro.
  - Medios: día de la semana mal en P-04, falta la Nota en P-04 móvil, foco invisible en la búsqueda de P-02, enlaces menores de 44 px en P-04 móvil, dinero sin cifras tabulares en 10 lugares y tablero sin «+ Nuevo».

## Verificación
Usé scripts propios con `python3 -I` en `datos/tmp/rev-gastos-r1/`, ya borrados:
- `forma.py` revisó lo siguiente:
  - líneas obligatorias y JSON de `data-props`, con `$preview` igual al fragmento;
  - `html.parser`, que no encontró etiquetas sin cerrar;
  - hex, `rgb()` y `hsl()` literales: 0;
  - `var(--x)` que no están en `tokens.json` o `tokens.css`: 0;
  - grep de P3-12 con las 21 cifras: 0; `$NNNN.NN` sin coma: 0; «-$» ASCII: 0; emoji: 0;
  - href fuera de la lista del lienzo o con ancho distinto: 0;
  - raíz móvil de 412 px y raíz de escritorio fluida con `expand:fill`.
- `contraste.py` hizo una cascada simplificada sobre el helmet y los estilos en línea, con variantes `.on` y `.sel`. Midió 25 pares de texto y 25 de UI; ninguno queda bajo umbral, con mínimo de 4.51:1. `pares.py` recalculó 12 pares y midió el placeholder del navegador: 2.25:1 en Chromium y 3.05:1 en Firefox.
- `tactil.py` midió alturas en móvil con especificidad CSS: 2 enlaces quedan bajo 44 px, los de P-04 móvil.
- `datos.py` comprobó que los montos literales que no están en el JSON son 0 y las fechas ISO que no están en el JSON también. Las 33 filas g01–g33 coinciden con `gastos[]`. La muestra fue mayor que 10: 152 montos, 21 fechas y 33 filas.
- `tabular.py` buscó dinero sin `tabular-nums` (H6) y px literales (H11).
- Un validador de días de la semana encontró «Mié 12 oct 2026», que debe ser «Lun 12 oct 2026».
- Corrí además el comando de fechas de P3-12. Las 4 fechas coincidentes que aparecen en el área van con montos sintéticos; lo revisé a mano.

## Supuestos
- **Placeholder (H1):** el contraste se calculó con el color por defecto del navegador, porque el área no lo define: Chromium `darkGray` #a9a9a9 y Firefox con el texto al 54 %. Lo califiqué alto porque la regla del encargo dice «contraste bajo umbral = alta».
- **Cascada:** la de los scripts es una aproximación. No ejecuté `renderVals()` ni rendericé; los pares que dependen de huecos `{{ }}` en `class` se evaluaron con las variantes `on` y `sel`.
- **Cifras tabulares (H6):** lo califiqué medio y no alto. Las cifras principales (tablas, tarjetas y totales) sí cumplen; solo fallan montos dentro de frases y algunos encabezados.
- **«+ Nuevo» (H7):** lo califiqué medio porque es un elemento secundario de la anatomía de DESIGN §7.5.3.
- **Pestañas de P-02:** Resumen, Tendencias y Calendario con `href="#"` son correctas, porque son subsecciones de P-02 y no tienen artboard en la lista del lienzo.
- **Tabla:** la tabla con `div` y roles ARIA es un equivalente aceptado por la restricción del analizador de HTML. No la cuento como hallazgo.

## Bloqueos
Ninguno.

## Pendientes
- Vuelta 2: que el autor aplique las correcciones 1 a 11 de `P3-iter-1-gastos.md`.
- Para el orquestador:
  - decidir `is_interactive`, que hoy está en `true` y la plantilla pide `false`;
  - medir los `h` reales en el lienzo;
  - unificar «páginas de 50» (U1 N-040) con «200 por página» (DESIGN §7.5.2, K6).
- Para el dueño de los datos: nombres de archivo de tickets con fecha ISO (H9) y los promedios de N-042 que faltan en el JSON.

## Consumo aproximado
- Leído: unos 150 KB. Incluye criterios P3 y U5 por secciones, `format.md` completo, DESIGN §7.5, §7.14, §9.3 y §3.3 por grep, U1 §2.1, §2.2, §2.4 y §3.4 a §3.7, el reporte del autor y extractos de texto de los 6 artboards; no leí completos los 180 KB de los artboards.
- Escrito: unos 28 KB (scripts de unos 15 KB, ya borrados, e informe y reporte de unos 13 KB).
- Acciones de herramienta: unas 40.
