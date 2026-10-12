# revisor-creditos-r3 · vuelta 3 · Créditos (P-17 a P-20)

Agente: `revisor-creditos-r3` (revisor independiente, vuelta 3 de 3; no corrige) · Fecha: 12-oct-2026 · Criterio: P3-07 · Informe: `docs/estado/wfp3/P3-iter-3-creditos.md`.

## Hecho

1. Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/creditos/` y `fragmentos/creditos.json` contra P3-07, P3-08, P3-12, U5-04, U5-06, U5-07, DESIGN, U1 (§2.1, §2.2, §3.19 a §3.22 y §6) y `referencia-tipo/format.md`.
2. Comprobé los hallazgos de la vuelta 2. **A-01 (alto) está cerrado**: el relleno es `--ui-viz-meter` sobre `--ui-viz-track` y da 4.51 en claro y 5.50 en oscuro, que es C-102. **M-01 (medio) está cerrado**: los 13 pagos de la agenda de P-17 m abren la hoja «Marcar pago» o «Ver pago», y los datos coinciden 13 de 13 con `pagos_programados`. Los 11 bajos de la vuelta 2 siguen como estaban.
3. **Veredicto: aprobado, con 0 altos, 0 medios y 15 bajos.** De los bajos, 11 vienen de la vuelta 2 (B-01 a B-11) y 4 son nuevos (B-12 a B-15): nombres accesibles repetidos en los botones de pago, `hint-placeholder-count` faltantes y `onsubmit` en texto, token de esqueleto sobre la página y la tabla de P-20 e que se desplaza a 1440 px.

## Verificación

Usé scripts propios en `datos/tmp/rev-cred3/` con `python3 -I` y node; ya los borré.
- **Forma, tokens, P3-12 y enlaces**: 0 fallas en los 8 artboards y el fragmento. Esto cubre las líneas obligatorias exactas, `html.parser`, `$preview` = fragmento, stems únicos en los 90 artboards del lienzo, 0 hex, 0 `var()` inexistentes, el grep de P3-12 con y sin coma y en centavos, 0 credenciales, `href` existentes y del mismo ancho, y el fragmento exacto.
- **Render**: corrí `renderVals()` en node y exploré estados por BFS (150 por archivo en P-17 y P-18), sin errores. Expandí la plantilla en Python con una cascada CSS propia y medí el contraste en los dos temas, en reposo y en hover: 61,186 evaluaciones de texto y 14,612 de interfaz, con **0 bajo umbral** fuera de los pares decorativos de DESIGN.
- **Accesibilidad**: 0 controles sin nombre, 0 campos sin `label`, 0 interactivos anidados, 0 controles móviles bajo 44 px, 0 inglés, 0 emoji y 0 fechas ISO en el texto visible.
- **Datos**: de 237 montos, 228 salen literalmente del JSON y 9 son derivados cuya aritmética verifiqué. Las 78 fechas visibles salen todas del JSON. La intersección de fechas de P3-12 da 15, ninguna junto a una cifra prohibida.

## Supuestos

- Aplico la escala de severidad del encargo. Los `href="#"` con destino existente fuera de la barra lateral y de los botones principales cuentan como bajo (B-04), no como enlace roto, porque no fallan. Así lo calificaron también las vueltas 1 y 2.
- Tomo como decorativos, igual que DESIGN: la pista sobre la página en móvil (análoga a C-105), el velo de la hoja y los esqueletos.
- En móvil, Tendencias de P-17 muestra el total y un dato por mes, sin el desglose por línea. Lo acepto por U1 §6.2 («filas compactas + 1 propiedad»); el desglose está en P-17 e y en P-18 › Tendencias de cada línea.
- No medí alturas en un navegador (K45: sin navegador). Los `h` del fragmento siguen siendo estimados (B-11).
- `is_interactive: false` en P-20, que tiene enlaces, no lo cuento como hallazgo: la plantilla del encargo lo trae en `false`.

## Bloqueos

- Ninguno.

## Pendientes

- Para el orquestador: copiar `P3-iter-3-creditos.md` a `docs/design/salida/iter-3/` (P3-06), medir las alturas reales y decidir si los 15 bajos se arrastran a la referencia visual o a la fase siguiente.
- Para el dueño de `datos-sinteticos.json`: los huecos de B-03.

## Consumo aproximado

Lectura ≈ 230 KB: criterios, format.md, §10.5 y §3.19 a §3.22 de U1, los informes de la vuelta 2, las consultas al JSON y la salida de los scripts. Escritura ≈ 45 KB: scripts temporales, el informe y este reporte. Unas 55 acciones de herramienta.
