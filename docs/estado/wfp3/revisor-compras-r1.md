# revisor-compras-r1 · Pista 3, fase 2 (lienzo) · vuelta 1 de 3 · 12-oct-2026

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/compras/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/compras.json`. No corregí nada.
- Escribí el informe P3-07 en `docs/estado/wfp3/P3-iter-1-compras.md`. Trae:
  - las verificaciones (a)–(h);
  - el recálculo de 45 pares de contraste en claro y oscuro;
  - la tabla de fidelidad contra U1 §2.1/§2.2/§3;
  - la lista de 8 puntos de P3-08;
  - la tabla artboard × regla;
  - las correcciones para la vuelta 2;
  - las secciones P3-04 y P3-12.
- **Dictamen: NO APROBADO.** 2 altos, 5 medios y 12 bajos.
  - A-01: `@media` sin cerrar en el CSS de los 4 artboards de escritorio, que rompe su diseño en 1440.
  - A-02: fechas cortas sin año en lugar de `10-oct-2026`.
  - M-01: dinero en frases sin cifras tabulares.
  - M-02: el «+» de la barra inferior mide 40 px de ancho.
  - M-03: a la barra lateral le faltan «‹‹», el interruptor de tema y los chevrons.
  - M-04: el `h1` de móvil dice «· móvil».
  - M-05: el historial de P-15 no tiene marcadores por serie ni etiquetas directas.

## Verificación
- Usé scripts propios con `python3 -I` en `datos/tmp/rev-compras-r1/` (`verif.py`, `clases.py`, `tabular.py` y `contraste.py`). Los borré al final.
- **Forma**: el prefijo y el sufijo obligatorios coinciden exactos en los 8. `data-props` parsea en los 8 y `$preview` coincide con `w`×`h` del fragmento. `html.parser` no encuentra etiquetas sin cerrar. `lang` es es-MX. Títulos únicos, 0 stems duplicados en el lienzo y nombres válidos. Las clases usadas están definidas en su `<style>`. **Las llaves del CSS no cuadran en los 4 `-e`**: es A-01. Revisé además los 9 directorios del lienzo y solo compras tiene este defecto.
- **Tokens**: 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`. Todos los `var(--x)` existen en `tokens.json`.
- **Contraste**: 90 combinaciones, 0 bajo umbral en el código tal como se quiso escribir. La excepción es el eje de P-15 en oscuro, que queda bajo umbral solo por A-01.
- **P3-12**: 0 cifras canon (con y sin coma y en centavos), 0 montos de 4 cifras sin coma, 0 credenciales y 0 fechas que coincidan con el calendario real. Tracé 15 valores a `datos-sinteticos.json` y crucé las 120 celdas del Detalle de P-14 contra `canasta.detalle`: 0 discrepancias.
- **Accesibilidad**: 0 controles sin nombre o sin `label`, 0 emoji, 0 inglés, 0 fechas ISO, dinero en `$1,234.56` y negativos con U+2212. En móvil medí 44 px en todos los controles salvo el «+» (40 de ancho).
- **Raíz**: móvil de 412 px fijos con border-box; escritorio fluido con `expand: fill` en el fragmento.
- **Enlaces**: todos van a la lista del lienzo con el mismo ancho, o a «#». Los 22 destinos de artboard existen.

## Supuestos
- El encargo fija las rutas `docs/estado/wfp3/P3-iter-1-compras.md` y `docs/estado/wfp3/revisor-compras-r1.md`. Las seguí, aunque P3-07 escribe `docs/estado/P3-iter-N-<area>.md` sin `wfp3/` (la nota de rutas de P3 §1 permite que el orquestador cambie la carpeta).
- Califiqué A-02 como alta porque P3-08 y U5-07 nombran el formato `10-oct-2026` como decisión y la rúbrica del encargo pone «decisión contradicha» en alta. Si el orquestador acepta fechas cortas en celdas densas, que lo declare como excepción y A-02 baja a media.
- Califiqué M-03 como media y no alta: la barra lateral no contradice la decisión, solo no dibuja el control de colapsar. Si la barra canónica es solo `shell-barra-lateral-e`, el orquestador puede declararlo.
- Un par de contraste que no aparece literal en §10.5 se acepta si cumple 4.5:1 (texto) o 3:1 (UI) por cálculo. Uso la fórmula WCAG 2.x sin truncar y resuelvo las referencias `{…}` de `tokens.json`.
- Fuera del conteo de 3:1 quedan, por decorativos (§10.5), los divisores, la retícula y los esqueletos. Las marcas de tienda de 10 px se midieron igual y cumplen.
- La palabra «Error» como título de estado es español y no cuenta como inglés de U5-07. «USDA Foundation Foods» cuenta como nombre propio (USDA está en la lista permitida).

## Bloqueos
Ninguno.

## Pendientes
- Para el orquestador:
  - medir al renderizar los altos reales de los móviles (B-03);
  - confirmar A-01 en el render de 1440 después de la corrección.
- Para el responsable de `datos-sinteticos.json`: la incoherencia de «Leche entera» (B-12).

## Consumo aproximado
- Leído: unos 190 KB entre los criterios P3 y U5, format.md, las secciones de U1 y DESIGN, los artboards (marcado condensado), los datos y el reporte del autor.
- Escrito: unos 34 KB (los 2 informes y los 4 scripts temporales, ya borrados).
- Acciones: unas 42.
