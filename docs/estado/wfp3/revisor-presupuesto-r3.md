# revisor-presupuesto-r3 · vuelta 3 de 3 · 12-oct-2026

## Hecho
- Revisé, sin corregir, los 8 artboards de `docs/design/maquetas/lienzo/project/presupuesto/` y `fragmentos/presupuesto.json`. El informe es `docs/estado/wfp3/P3-iter-3-presupuesto.md` (criterio P3-07).
- **Veredicto: APROBADO.** 0 altas, 0 medias, 15 bajas abiertas (las de la vuelta 2, sin cambio) y 1 observación sin severidad (O1).
- **N1 (alta de la vuelta 2) cerrado.** En P-07 (m y e), el halo color tarjeta de 6 px bajo la línea a guiones y bajo «Hoy» deja la tarjeta como color contiguo:
  - `ui-chart-2`: 3.15:1 en claro y 7.73:1 en oscuro (C-107);
  - `ui-viz-baseline`: 3.47:1 y 6.09:1.
- Diff contra `docs/design/salida/iter-2/presupuesto/`: solo cambió P-07, y solo en el SVG y en `hoyPath`.

## Verificación
- Scripts propios en `datos/tmp/rev-pres-r3/` (`python3 -I` y node), borrados al terminar.
- **Forma:** 751 comprobaciones sobre la fuente y otra pasada sobre la plantilla resuelta con `renderVals()`. Se revisó:
  - cabecera exacta, `support.js` y `tokens.css`;
  - `x-dc` y `helmet`;
  - `data-props` JSON válido con `$preview` igual al fragmento;
  - `html.parser` limpio;
  - `hint-*` presentes y 0 huecos sin resolver;
  - 90 stems únicos.

  Las únicas «fallas» del primer barrido fueron falsos positivos del script: `--font-sans` existe en `type.families`, y los `href` con `{{…}}` y las anclas `#id` resueltas son válidos.
- **Tokens:** 0 hex fuera de `data-props`, y todos los `var(--x)` existen en `tokens.json` y en `tokens.css`.
- **Contraste:** calculado con los hex de `tokens.json` en claro y oscuro.
  - Los 13 pares de una misma regla dan ≥ 4.5.
  - En el producto cruzado de 11 colores de texto × 23 fondos, ningún par bajo umbral ocurre con texto.
  - En la gráfica, todos los pares contiguos son ≥ 3.
- **P3-12:**
  - Grep de las 21 cifras canon: 0. En centavos: 0. Credenciales: 0.
  - Fechas coincidentes: 6, ninguna con un monto del reporte.
  - 600 de 616 montos visibles son iguales a un campo de `datos-sinteticos.json`; los otros 16 son sumas o subcampos del mismo archivo.
  - Muestra de 10 valores: cotejada.
- **Accesibilidad:**
  - 274 controles con nombre accesible y 0 controles falsos.
  - 107 controles móviles, 0 bajo 44 px.
  - 0 emoji, 0 ISO visibles, 0 «-$».
- **Enlaces:** todos van a archivos de la lista del mismo ancho, a «#» o a anclas con id existente.

## Supuestos
- El halo que agrega la corrección (no lo pide §8.3) se acepta como tratamiento de contraste: la línea «Hoy» sigue siendo de 1 px.
- El título igual en m y e es el que fija el encargo (nombre de la pantalla). El fragmento los distingue.
- La copia del informe en `docs/design/salida/iter-3/` (P3-06) le toca al orquestador. Este agente solo escribe en `docs/estado/wfp3/`.

## Bloqueos
- Ninguno.

## Pendientes
- En el render publicado, en claro y en oscuro: medir el contraste de la gráfica de P-07 (incluido O1, el arranque de las dos series en el primer tramo móvil) y las alturas reales (B4).
- Las 15 bajas quedan para implementación o para la referencia visual.

## Consumo aproximado
- Leído ≈ 95 KB (informes previos, criterios, format.md, craft.md por secciones, DESIGN §8.3 a §8.6, U1 §2.1/§2.2/§6.4 por grep y P-07 por partes; el resto, por script).
- Escrito ≈ 21 KB (informe y este reporte).
- ≈ 40 acciones. Sin navegador, sin Artifact, sin red, sin git y sin Telegram.
