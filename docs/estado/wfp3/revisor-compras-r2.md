# revisor-compras-r2 · Pista 3, fase 2 · vuelta 2 de 3 · 12-oct-2026

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/compras/` (P-13 a P-16, móvil y escritorio) y el fragmento `docs/design/maquetas/lienzo/fragmentos/compras.json`. Los revisé contra P3-07, P3-08 y P3-12 (criterios 1.3), U5-04, U5-06 y U5-07 (lo que se ve en el código), DESIGN.md (§3.3, §4.5, §4.6, §5.2, §7.9, §7.10, §7.12, §7.14, §8.6 a §8.8 y §10.5), U1 (§0.4, §1.2, §1.3, §2.1, §2.2 y §6.4), `format.md` y `craft.md`.
- Comprobé los 7 hallazgos de la vuelta 1 (2 altos y 5 medios): **los 7 están cerrados**. También cerré B-02 y reformulé B-08.
- Informe: `docs/estado/wfp3/P3-iter-2-compras.md`, con el estado de cada hallazgo, las verificaciones (a) a (h), 47 pares de contraste en los dos temas, la procedencia, la fidelidad, la lista P3-08, la tabla artboard × regla y los bajos para la vuelta 3.
- **Dictamen: APROBADO** (0 altos y 0 medios). Quedan 16 bajos: 11 de la vuelta 1, uno de ellos reformulado, y 5 nuevos (B-13 a B-17).

## Verificación
Scripts propios con `python3 -I` en `datos/tmp/rev-compras-r2/`, borrados al terminar:
- `verif.py` revisó la forma: el prefijo y el cierre obligatorios, `data-props` que parsea, `$preview` igual al fragmento, `html.parser` sin etiquetas abiertas, llaves CSS balanceadas y un solo `@media` con solo `.side`, títulos y stems únicos y nombres válidos. También revisó 0 hex y 0 funciones de color, que todos los `var(--x)` existan en `tokens.json`, el grep de P3-12 con y sin coma y en centavos, el formato del dinero y U+2212, las fechas (0 cortas, 0 ISO y 0 rangos), 0 emoji, los `href` contra la lista del lienzo y el mismo ancho, y el fragmento. Resultado: 0 fallas.
- `contraste.py`: 47 pares × 2 temas, 0 bajo umbral, más la matriz completa de primer plano contra fondo.
- `tabular.py`: 51 de 51 nodos con dinero dentro de `.n`.
- `tactil.py`: 55 controles en móvil, todos de 44 px o más.
- `proc.py`: 741 valores contra `datos-sinteticos.json`, con 0 discrepancias.
- P3-12 con los comandos del criterio: 0 cifras canon, 0 credenciales y 0 montos de 4 cifras sin coma. Hay 15 fechas en el cruce de calendarios y ninguna aparece en compras.

## Supuestos
- `tokens.css` del sistema y `support.js` los instala el orquestador. Que no estén junto a los artboards no cuenta como falla.
- Los enlaces de forma corta dentro del área (`compras-p14-canasta-e.dc.html`) son válidos: format.md resuelve el `href` como ruta relativa al artboard, y llega al mismo archivo que `../compras/…`.
- Los altos de móvil y la posición de las etiquetas directas se estimaron sin navegador (K45: 0 acciones de navegador). B-03 y B-16 se confirman al renderizar.
- «Stock por unidad» es el término de U1 §2.1 y no lo cuento como inglés.
- La marca de 10 px de los mercados sobre el chip en hover (2.11:1 en oscuro) es decorativa, porque es `aria-hidden` y repite el nombre escrito. No la cuento como hallazgo.

## Bloqueos
Ninguno.

## Pendientes
- Corresponde al orquestador decidir si los bajos B-01 y B-03 a B-17 se atienden en la vuelta 3 o pasan a `referencia-visual.md` como notas. B-12 es del responsable de los datos.
- P3-06 pide una copia de este informe en `docs/design/salida/iter-2/`. Le toca al orquestador; yo no escribo fuera de `docs/estado/wfp3/`.

## Consumo aproximado
Leí unos 190 KB: los criterios, format.md, secciones de craft.md, DESIGN y U1, el marcado de los 8 artboards y muestras de los datos. Escribí unos 31 KB: los dos informes y unos 15 KB de scripts temporales. Fueron unas 40 acciones. Hice 0 llamadas a Artifact, 0 de navegador y 0 de red, y no usé git ni Telegram.
