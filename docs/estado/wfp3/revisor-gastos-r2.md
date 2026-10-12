# revisor-gastos-r2: reporte (Pista 3, fase 2 «lienzo», vuelta 2 de 3, 12-oct-2026)

## Hecho
- Revisé los 6 artboards de `docs/design/maquetas/lienzo/project/gastos/`: P-02, P-04 y P-05, en móvil 412 y escritorio 1440. Revisé también el fragmento `fragmentos/gastos.json`.
- Los comparé contra estas fuentes:
  - P3 1.3 (P3-07, P3-08 y P3-12);
  - U5-04, U5-06 y U5-07;
  - DESIGN.md: §3.3, §7.2, §7.5.3, §7.5.4, §9.2, §9.3 y §10.5;
  - U1: las filas de P-02, P-04 y P-05 en §2.1, §2.2, §2.4 y §6.2;
  - `referencia-tipo/format.md`: interactividad, huecos y PAGE;
  - `datos-sinteticos.json` 1.2.
- Comparé cada archivo con su copia de la vuelta 1 (`docs/design/salida/iter-1/gastos/`) para comprobar el cierre de los 7 hallazgos altos y medios.
- El informe está en `docs/estado/wfp3/P3-iter-2-gastos.md`. Trae:
  - el cierre de H1 a H7 con evidencia;
  - los hallazgos vigentes;
  - las verificaciones (a) a (h);
  - el contraste en los dos temas;
  - una muestra de 12 valores con su ruta en el JSON;
  - la tabla de fidelidad contra U1;
  - los 8 puntos de P3-08;
  - la tabla artboard × regla;
  - las líneas de P3-04 y P3-12.
- **Veredicto: APROBADO.** Hay 0 hallazgos altos, 0 medios y 8 bajos.
  - Los 7 altos y medios de la vuelta 1 están **cerrados**: H1 placeholder, H2 día de la semana, H3 Nota en P-04 móvil, H4 foco de la búsqueda, H5 enlaces de 44 px, H6 cifras tabulares y H7 «+ Nuevo».
  - Siguen abiertos 5 bajos de la vuelta 1 (H8 a H12).
  - Hay 3 bajos nuevos, los tres en p02-e:
    - N1: el calendario usa `role="grid"` sin filas ni celdas.
    - N2: el texto menciona un menú «⋯» que las tarjetas del tablero no dibujan.
    - N3: los números de los días de otro mes no llevan `ui-text-disabled`.

## Verificación
Corrí scripts propios con `python3 -I` en `datos/tmp/rev-gastos-r2/`; ya están borrados.

`rev.py` dio 162 OK y 0 fallas:
- **Forma:** la cabecera y el cierre son exactos. El JSON de `data-props` es válido y `$preview` coincide con el fragmento. `html.parser` no encontró etiquetas sin cerrar.
- **Plantilla y nombres:** los huecos `{{ }}` son búsquedas con punto y salen de `renderVals()`. Los nombres de archivo son válidos y los 90 stems del lienzo son únicos.
- **Tokens y colores:** 0 hex, 0 `rgb()` y 0 `hsl()`. Todas las `var(--x)` están en `tokens.json` y `tokens.css`.
- **P3-12 y formato de dinero:** 0 cifras prohibidas, 0 credenciales, 0 `$NNNN.NN` sin separador y 0 negativos con «-$».
- **Enlaces:** todos van a la lista del lienzo en el mismo ancho o a «#». Todos los destinos existen en disco.
- **Raíz:** la de móvil mide 412 con `border-box`; la de escritorio es fluida, con `expand:fill`.
- **Fechas:** los días de la semana son correctos.

`cas.py` hizo una cascada con especificidad sobre el helmet y los estilos en línea, en reposo y con las variantes `on`/`sel`:
- **Contraste:** 27 pares de texto (mínimo 4.51:1) y 17 de UI (mínimo 3.19:1), sin ninguno bajo umbral en claro ni en oscuro.
- **Placeholders:** 5.40:1 en claro y 10.43:1 en oscuro.
- **Objetivos táctiles en móvil:** solo `a.gal` quedó sin alto declarado, y es una tarjeta con portada 4:3, de mucho más de 44 px.
- **Dinero:** 0 montos fuera de `.num`. Revisé a mano 73 interpolaciones `{{ }}` que no llevan `.num` y ninguna es dinero.

`datos.py` revisó los montos y las fechas:
- Hay 241 montos en el área y 0 faltan en el JSON.
- Ninguna fecha `d-mmm-aaaa` falta en el JSON.
- Con el comando de fechas de P3-12, 4 fechas coinciden con `oraculo/model.json` o con `datos/semilla/data.json`. Son las mismas que en la vuelta 1, son de calendario quincenal y van con montos sintéticos.

`ingles.py` no encontró texto en inglés. Revisó el texto visible, `aria-label`, `placeholder` y las cadenas de datos; la única coincidencia fue `R.filter`, que es código JS.

## Supuestos
- **Contraste:** se calculó con una cascada aproximada, sin ejecutar `renderVals()` ni renderizar. El anillo de foco se midió contra cada fondo usado en el archivo, como conjunto amplio.
- **N1, N2 y N3:** los califiqué bajos:
  - N1: la estructura ARIA del grid la evaluará axe en U5-03, que queda fuera del alcance estático de este encargo.
  - N2: la alternativa al arrastre existe, porque la categoría se edita en la tabla y en P-04, y en Notion el «⋯» aparece al pasar el puntero.
  - N3: es solo de estilo, y el par C-093 cumple 5.40:1.
- **Alto de p04-m:** lo estimé en unos 2,040 px de contenido en estado plegado, contra 2,180 de la raíz. No lo medí en navegador, porque K45 no lo permite.

## Bloqueos
Ninguno.

## Pendientes
- Bajos para la vuelta 3, si el orquestador la abre: H8, H9, H10, H11, N1, N2 y N3. Las correcciones están en §9 del informe.
- Para el orquestador:
  - decidir `is_interactive` (H12); `format.md` respalda `true`;
  - medir los `h` reales en el lienzo;
  - unificar «páginas de 50» (U1 N-040) con «200 por página» (DESIGN §7.5.2).

## Consumo aproximado
- **Leído:** unos 150 KB. Incluye:
  - los reportes y el informe de la vuelta 1;
  - los criterios P3 y U5, por secciones;
  - DESIGN §7.2, §7.5.3, §7.5.4, §9.2, §9.3 y la cabecera de §10.5;
  - las filas de U1;
  - los diffs contra iter-1;
  - extractos de los artboards;
  - la salida de los scripts.
- **Escrito:** unos 40 KB: scripts de unos 22 KB, ya borrados, e informe y reporte de unos 18 KB.
- **Acciones de herramienta:** unas 50.
