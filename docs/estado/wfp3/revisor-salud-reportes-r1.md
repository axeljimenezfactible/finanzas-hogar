# revisor-salud-reportes-r1 · Pista 3 fase 2 (lienzo) · vuelta 1 · 12-oct-2026

## Hecho
- Revisión independiente (sin corregir) de los 8 artboards de `docs/design/maquetas/lienzo/project/salud-reportes/` y de `docs/design/maquetas/lienzo/fragmentos/salud-reportes.json`.
- Informe P3-07 en `docs/estado/wfp3/P3-iter-1-salud-reportes.md`: veredicto **no aprobado**, 0 altos, 5 medios (M01 barra lateral sin botón de colapsar; M02 filtro «Tipo» de avisos con tipos de reporte; M03 P-23 sin selector de vista ni filtro Estado, y filtros de móvil sin «Filtrar» + chips; M04 P-21 sin «Restricciones alimentarias»; M05 «semana 2026-W43» visible en P-24) y 14 bajos.

## Verificación
Scripts propios `python3 -I` en `datos/tmp/rev-sr-r1/` (borrados al terminar): forma y fragmento (líneas obligatorias, `data-props`, `$preview` = fragmento, `html.parser`, stems de los 88 artboards del lienzo), tokens (hex, `rgb()`, `var(--x)` contra `tokens.json` y `tokens.css`), contraste (55 pares en claro y oscuro con la fórmula WCAG), cifras P3-12 y credenciales, procedencia de los 57 importes y las 18 fechas contra `datos-sinteticos.json`, días de la semana, texto (formato de dinero, U+2212, cifras tabulares, fechas ISO, palabras de U5-07, emoji), accesibilidad (nombres, `label`, `svg`, `meter`, IDs, encabezados), objetivos táctiles en móvil (101 controles) y enlaces (120 `href` contra la lista del lienzo). Resultados en el informe.

## Supuestos
1. Informe en `docs/estado/wfp3/` por indicación de la tarea (P3-07 nombra `docs/estado/`).
2. Contraste estático por tokens resueltos, sin DOM renderizado (K45: sin navegador).
3. Severidades según la escala de la tarea; la barra lateral sin «‹‹» cuenta como elemento secundario ausente (media), no como decisión contradicha.
4. `--font-sans` (definido en `tokens.css`, ausente en `tokens.json`) se registra como observación (B14), no como defecto del área.

## Bloqueos
Ninguno.

## Pendientes
- Vuelta 2: corregir M01 a M05 y, si cabe, los bajos B01 a B13.
- Medir en el lienzo el alto real de los 8 artboards (B08) al publicar.
- Copia de este informe en `docs/design/salida/iter-1/` (P3-06): la hace el orquestador.

## Consumo aproximado
- Leído: ≈ 150 KB (criterios P3 y U5 por secciones, U1 y DESIGN por secciones, `format.md` completo, `craft.md` por secciones, partes del JSON sintético, `tokens.css`, artboards por volcado estructurado).
- Escrito: ≈ 23 KB (dos reportes) + ≈ 20 KB de scripts temporales borrados.
- Acciones: ≈ 50 de 60. Herramientas: solo Bash; 0 navegador, 0 Artifact, 0 Telegram, sin git.
