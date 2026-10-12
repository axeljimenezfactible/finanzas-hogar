# revisor-inicio-r2 · reporte

Fecha: 12-oct-2026. Pista 3, fase 2 (lienzo). Revisión independiente del área `inicio`, vuelta 2 de 3. Revisé 8 artboards de P-01 y el fragmento `inicio.json`. No corrijo nada.

## Hecho
- Informe de revisión (P3-07) en `docs/estado/wfp3/P3-iter-2-inicio.md`. Trae:
  - el cierre de los hallazgos de la vuelta 1;
  - los hallazgos nuevos;
  - la tabla artboard × regla;
  - la fidelidad contra U1 y U3, el contraste y las decisiones de P3-08;
  - las líneas de evidencia de P3-04, P3-06, P3-07, P3-08, P3-12, U5-04, U5-06 y U5-07;
  - la sección «Texto con aspecto de instrucción», sin citas;
  - las correcciones para la vuelta 3.
- **Veredicto: NO APROBADO.** Hay 1 alto, 0 medios y 9 bajos.
  - **A-1 (alto de la vuelta 1): cerrado en los artboards.** Las 6 insignias de meta dicen «En curso», con rol info e icono `info`. «Atrasada» y «En ruta» aparecen 0 veces en Inicio, y el `diff` contra `iter-1` toca solo esas 6 líneas.
  - **A-2 (alto nuevo).** D-6 «Metas» muestra el avance en entero («42 %» y «78 %»). U3 D-6 y la decisión Q2-17 C piden un decimal en Metas: «42.4 %» y «78.3 %». Pasa en detalle-m y en detalle-e. El origen es `LEEME.md` regla 4 («los porcentajes de avance son enteros») y el JSON, que no trae las décimas.
  - **Bajos.** Siguen abiertos B-1, B-3, B-4, B-5, B-6 y B-7. B-2 se cierra porque fue un falso positivo de la vuelta 1: la fila de D-3 sí lleva el círculo punteado. Hay 3 bajos nuevos:
    - B-9: iconos de insignia de 16 px.
    - B-10: la etiqueta «Error» choca con la lista de U5-07.
    - B-11: el JSON y `LEEME.md` siguen con «Atrasada» y «En ruta», y P-09 y P-26.13 los muestran, así que el lienzo no cuadra entre áreas.

## Verificación
Usé scripts propios (`python3 -I` y `node` para `renderVals()`) en `datos/tmp/rev-inicio-r2/`. Ya están borrados.
- **Forma**:
  - 8/8 tienen las líneas obligatorias exactas y en orden.
  - `data-props` parsea y `$preview` es igual a `w`×`h` del fragmento.
  - `html.parser` da 0 pendientes.
  - Huecos solo con rutas, `hint-*` presentes y `renderVals()` ejecutable.
  - Stems únicos en los 90 del lienzo.
  - Raíz móvil de 412 px con `border-box` y `min-height` igual a `h`. Escritorio fluido con `expand: "fill"`.
  - El fragmento lista exactamente los 8 archivos.
- **Tokens**: 0 hex fuera de `data-props`. Todos los `var()` están en `tokens.css`. En `tokens.json` están todos salvo `font-sans`, que es `type.families.sans`.
- **Contraste**:
  - 108 combinaciones (54 por tema), 0 bajo umbral.
  - 104 coinciden con `contraste.csv` y `contraste-componentes.csv` al ±0.01. Las 4 restantes son decorativas o están por encima de 3:1.
  - Mínimo en texto: 4.80. Mínimo en UI: 3.06.
- **P3-12**:
  - Grep oficial: 0. En centavos, 0. Credenciales, 0.
  - Todos los montos salen de `datos-sinteticos.json`. Las excepciones son los ejes compactos y 2 sumas de T-2, que recalculé.
  - Muestra de 10 montos y 3 fechas trazada a su ruta en el JSON.
  - Las coincidencias de fecha (30 oct, 15 y 30 nov) van con montos sintéticos.
- **Accesibilidad estática**:
  - 70 controles en móvil, 0 por debajo de 44 px.
  - 0 controles anidados o sin nombre.
  - 207 SVG con `aria-hidden` o `role="img"` y nombre.
  - 0 emoji, 0 fechas ISO, 0 negativos con guion y 0 montos sin `tabular-nums`.

## Supuestos
1. El contraste es estático: una cascada propia sobre el CSS del `<helmet>` y los tokens de los dos temas, sin navegador (K45). La medición con estilos computados la hace U5-04.
2. A-2 va como **alto** por el mismo criterio que A-1 en la vuelta 1: contradice una decisión de la puerta 2 (Q2-17 C) transcrita en U3. Si el orquestador considera que el formato del porcentaje es un elemento secundario, bajaría a medio. Con cualquiera de las dos severidades, el área queda sin aprobar.
3. Q2-17 C se redondea con ROUND_HALF_UP a décimas: 615000 ÷ 1450000 = 42.41 %, que da 42.4 %; 1410000 ÷ 1800000 = 78.33 %, que da 78.3 %. El formato «42.4 %» sigue RN-18 («50.0 %»).
4. «Kueski Pay» es un nombre propio y «Error» es español. Solo se señala (B-10) porque el script de U5-07 lo cuenta como inglés.
5. Los halos de 2 px de los marcadores, la rejilla y el carril son decorativos (DESIGN §8.1) y quedan fuera del umbral.

## Bloqueos
Ninguno.

## Pendientes
- **Orquestador**:
  - Agregar a `datos-sinteticos.json` las décimas de avance de las metas.
  - Cambiar los estados de meta a «En curso» (`in_progress`) en `ahorro.metas[]` y en `inicio.tarjetas.ahorro.next_goal`.
  - Corregir `LEEME.md`: regla 4 y líneas 29 y 56.
  - Volver a correr `verificar_sinteticos.py`.
  - Avisar a las áreas presupuesto (P-09) y ajustes (P-26.13).
- **Autor de Inicio (vuelta 3)**: A-2 en `inicio-p01-detalle-m` y `inicio-p01-detalle-e`, en el texto, en `aria-valuenow` y en `aria-valuetext`. Los bajos que pueda.
- Copiar este informe a `docs/design/salida/iter-2/` (P3-06): le toca al orquestador.
- Medir los altos reales al publicar (B-7).

## Consumo aproximado
- **Leído**: unos 150 KB.
  - Criterios P3 y U5, por secciones.
  - U3 §1 a §4 y §7, por secciones.
  - U1: solo las filas de P-01.
  - DESIGN: §5.2, §7.4 y el encabezado de §10.5.
  - `format.md` entero.
  - Informe y corrección de la vuelta 1.
  - Fragmentos de los artboards y salidas de los scripts.
- **Escrito**: unos 22 KB entre el informe y este reporte, más unos 25 KB de scripts temporales, ya borrados.
- Unas 50 acciones de herramienta.
