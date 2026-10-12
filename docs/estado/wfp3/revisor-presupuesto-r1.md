# revisor-presupuesto-r1 · Pista 3, fase 2 (lienzo), vuelta 1 de 3

Fecha: 12-oct-2026 · Agente: revisor independiente del área «presupuesto» (no corrige).

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/presupuesto/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/presupuesto.json`.
- Escribí el informe P3-07 en `docs/estado/wfp3/P3-iter-1-presupuesto.md`. Trae la tabla artboard × regla, la evidencia de P3-04, P3-06, P3-08, P3-12, U5-04, U5-06 y U5-07, la tabla de fidelidad contra U1 y U4, y los hallazgos con su corrección.
- Veredicto: **no aprobado**, con 1 alta (A1, campos de fecha nativos), 4 medias (M1 gráfica de P-07 fuera de DESIGN §8.3; M2 sin partes hogar/Karen en N-091; M3 interruptor N-318 operable en P-07 y de 36 px; M4 enlaces al índice de Ajustes en vez de la sección) y 9 bajas.

## Verificación
Hice scripts propios con `python3 -I` y uno con node para ejecutar `renderVals()` con un stub de `DCLogic` y expandir `sc-for`, `sc-if` y los huecos. Estuvieron en `datos/tmp/rev-pres-r1/` y los borré al terminar. Resultados:
- (a) Forma: 29 comprobaciones por archivo, 0 fallas. Cubren las líneas obligatorias, `data-props` con `$preview` igual al fragmento, `html.parser` balanceado, `lang` es-MX, título propio, nombres válidos y stems únicos en las 90 entradas del lienzo. No quedó ningún hueco sin resolver al expandir.
- (b) Tokens: 0 hex y 0 `rgb()/hsl()` fuera de `data-props`. Todos los `var(--x)` están en `tokens.json`; `--font-sans` sale de `type.families.sans`.
- (c) Contraste: 52 pares en claro y oscuro con los hex de `tokens.json`, 0 bajo umbral. Doce de ellos están en la tabla del informe.
- (d) P3-12: el grep de las 21 cifras canon, con y sin coma y en centavos, da 0. El de credenciales da 0. Comparé 10 valores con `datos-sinteticos.json` y todos coinciden. Las 15 fechas que coinciden con `model.json` van con montos sintéticos.
- (e) Fidelidad: la tabla por pantalla está en el informe. Los estados vacío, carga y error aparecen en los 4 de escritorio, con los textos de U1 §2.2.
- (f) Accesibilidad: 0 controles sin nombre, 0 campos sin etiqueta, 0 SVG sin `aria-hidden`. Sin emoji, sin barra de estado falsa y sin texto en inglés. Todo el dinero va en `$1,234.56` con U+2212. Los objetivos móviles miden 44 px o más, salvo el interruptor de P-07 (36 px de ancho, M3).
- (g) Móvil con raíz fija de 412 px y `box-sizing`. Escritorio fluido, con `expand: fill` solo en esos 4.
- (h) Los 153 `href` expandidos apuntan a archivos del lienzo, a «#» o a anclas que existen, siempre del mismo ancho.

## Supuestos
1. El informe va en `docs/estado/wfp3/P3-iter-1-presupuesto.md`, como pide el encargo. P3-07 dice `docs/estado/P3-iter-N-<area>.md`. El orquestador decide si lo copia ahí y a `docs/design/salida/iter-1/` (P3-06).
2. Puse A1 como alta porque P3-08 lista el formato de fecha `10-oct-2026` como decisión, y un `<input type="date">` nunca lo pinta así. Al texto corrido no le encontré ninguna falla de fecha.
3. Puse M1 como media y no como alta: la gráfica está, pero su forma no sigue DESIGN §8.3. No es una de las 8 decisiones de P3-08.
4. No medí alturas ni vi el render, porque no se usa navegador. B2 (preselección de la radio) y B4 (alturas) se confirman al publicar.

## Bloqueos
Ninguno.

## Pendientes
- Vuelta 2: que el autor corrija A1 y M1 a M4 y que se vuelva a revisar el área.
- Orquestador: decidir B3 (`is_interactive`), medir las alturas (B4) y, si lo cree útil, pedir al dueño de `datos-sinteticos.json` la deuda viva al fin del horizonte y por compra (B5).

## Consumo aproximado
- Leído: unos 170 KB (secciones de U1, DESIGN §7, §8, §9 y §10.5, U4 P-26.02, P3 y U5 por secciones, format.md, partes de craft.md, el reporte del autor y los artboards por partes y renderizados).
- Escrito: unos 22 KB (este reporte y el informe P3-07). Los scripts temporales ya se borraron.
- Acciones de herramienta: unas 50.
