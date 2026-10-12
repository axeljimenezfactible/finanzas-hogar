# revisor-salud-reportes-r2 · Pista 3 fase 2 (lienzo) · vuelta 2 · 12-oct-2026

## Hecho
- Revisé, sin corregir, los 8 artboards de `docs/design/maquetas/lienzo/project/salud-reportes/` y `docs/design/maquetas/lienzo/fragmentos/salud-reportes.json`. La línea base fue la copia de la vuelta 1 en `docs/design/salida/iter-1/salud-reportes/`, comparada etiqueta por etiqueta.
- El informe P3-07 está en `docs/estado/wfp3/P3-iter-2-salud-reportes.md`. Veredicto: **no aprobado**, con 0 altos, **1 medio** y 17 bajos.
  - M01 a M05 de la vuelta 1: los 5 **cerrados**.
  - Medio nuevo **N01**: las tablas de `reportes-p25-karen-e` muestran 22 importes sin `tabular-nums`, porque al `helmet` le falta `.tbl td,dl.kv dd{font-variant-numeric:tabular-nums}`, que sí tienen los otros artboards del área. Va contra P3-08 punto 5 y U5-06. Se corrige con una línea.
  - Bajos nuevos: N02 (dinero en prosa de los avisos y fechas de las tablas de P-23 escritorio sin cifras tabulares), N03 (`aria-haspopup` sin popup y `aria-label` en un `div` sin rol) y N04 (en móvil, «Filtrar» abre un `details` en línea y no una hoja; es observación). B01 a B14 siguen abiertos, porque la corrección solo tocó los medios.

## Verificación
Hice scripts propios con `python3 -I` en `datos/tmp/rev-sr-r2/` y los borré al terminar:
- Diff etiqueta por etiqueta entre la vuelta 1 y la vuelta 2.
- Forma: líneas obligatorias, `data-props`, `$preview` contra el fragmento, `html.parser`, IDs, stems de los 90 artboards del lienzo.
- Tokens: hex, `rgb()`/`hsl()`, `var(--x)` contra `tokens.json` y `tokens.css`.
- Contraste con una cascada CSS propia y los hex resueltos de `tokens.json` en claro y oscuro: 38 pares de texto y 26 de interfaz, 0 bajo umbral en lo que no es decorativo.
- P3-12: grep canon y de credenciales = 0; 57 importes y 21 fechas rastreados en `datos-sinteticos.json`; fechas coincidentes con el calendario real revisadas a mano.
- Texto: formato `$1,234.56`, U+2212, cifras tabulares, ISO y semana ISO, palabras de U5-07, emoji.
- Accesibilidad: nombres accesibles, `label`, `svg`; objetivos táctiles de 99 controles móviles.
- Enlaces y anchos.

## Supuestos
1. N01 cuenta como medio (una regla omitida en un artboard, no un diseño que contradiga la decisión). Es el mismo criterio que la vuelta 1 aplicó a M01.
2. N01 ya estaba en la vuelta 1 y aquella revisión no lo vio. Lo registro porque la aprobación exige 0 medios verificados.
3. El contraste se midió en estático, por tokens (K45: sin navegador).
4. «Reporte mensual», como tipo de aviso, se acepta por ser el literal del JSON.

## Bloqueos
Ninguno.

## Pendientes
- Vuelta 3: corregir N01, una línea en el `helmet` de `reportes-p25-karen-e.dc.html`. Si cabe, también N02, N03 y los bajos B01, B06 y B07.
- Medir en el lienzo el alto real de los 8 artboards al publicar (B08).
- Copiar este informe en `docs/design/salida/iter-2/` (P3-06): lo hace el orquestador.

## Consumo aproximado
- Leído: ≈ 110 KB (reportes de la vuelta 1 y de la corrección, DESIGN §7.4, §7.5, §7.9, §7.10 y §10.5 por secciones, U1 §2.1, §2.2 y §2.5 y N-230 a N-248 por grep, P3-12, U5-06 y U5-07, format.md por secciones, `reportes` del JSON sintético y volcados estructurados de los artboards).
- Escrito: ≈ 21 KB (dos reportes), más ≈ 18 KB de scripts temporales ya borrados.
- Acciones: ≈ 40 de 60. Herramientas: solo Bash. 0 navegador, 0 Artifact, 0 red, 0 Telegram, sin git.
