# P2c · bloqueos-respuestas (11-oct-2026, reanudación 22:35)

Agente: bloqueos-respuestas (Sonnet). Alcance: los 27 bloqueos abiertos de `docs/estado/bloqueos/`, `INDICE.md`, `docs/plan/decisiones.md` §K46 en adelante y este reporte. El intento anterior quedó a medias sin reporte; se trató como borrador no verificado y se comparó contra `RESPUESTAS.md` y `PUERTA-2.md` §3, §4, §5, §6 y §8.

## Hecho
- HECHO: los 27 bloqueos tienen la sección de la puerta 2: 25 con «Respuesta de Axel (puerta 2, 11-oct-2026)» y 5 con «Puerta 2 (11-oct-2026): sin respuesta» (3 de ellos conviven con respuestas en el mismo archivo: `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`; `S5-carritos-detalle` y `S5-profeco-filtro` solo traen la nota de sin respuesta). Cada entrada trae pregunta, letra, texto literal de la opción, Consecuencia, Qué cambia y Estado, y cada archivo cierra con «Estado del bloqueo (puerta 2): RESUELTO|ABIERTO».
- HECHO: revisión del borrador contra las fuentes: la letra de las 53 entradas con respuesta (Q2-11 a Q2-71 son 61 preguntas menos 8 sin respuesta) coincide con la recomendada de la hoja §8 (o con la explícita de RESPUESTAS.md en Q2-11 = B y Q2-64 = A); el texto literal de cada opción aparece en el archivo de bloqueo justo después de su letra; las «Consecuencia» que citan el «Impacto» del archivo aparecen en el original.
- HECHO (corrección): en `S5-terminos.md` la sección «Resultado» estaba sin acentos y mezclaba opciones; se reescribió con el resultado real (3 sitios «no verificado» por desafío anti-robot, S5-07 NO cumple en el fondo) y la pregunta a Axel: B, C o D. Las URL legales directas se anotan como autorización nueva que no está dada.
- HECHO: `INDICE.md` con cuenta nueva (63 resueltos, 6 abiertos), filas del resumen actualizadas y sección «Respuestas de la puerta 2 (11-oct-2026)» (tabla Q2 → bloqueo → letra → estado, y tabla de Q2-72 a Q2-77).
- HECHO: `decisiones.md` §K46 a §K52 (modo, regla de los «sin recomendación» con tabla de las 98, Q2-11 B con la aclaración de RESPUESTAS.md, Q2-64 A, Q2-96 A y Q2-98 A, K31/K44 FIRMES con valores, ratificaciones Q2-78 a Q2-94, Q2-95 y Q2-97). Solo hay citas de Axel donde RESPUESTAS.md las trae.

## Verificación
Script propio `/tmp/p2c/verify.py` (más `check.py`, `check2.py`, `check3.py`, `check4.py` en el mismo directorio):
- `total 69 · tocados 27 · resueltos tocados: []` y el conjunto tocado = los 27 abiertos de HEAD. (El script cuenta 28 «abiertos» en HEAD porque `criterios-nxx.md` no lleva «Respuesta de Axel» pero ya estaba resuelto; no se tocó. Resueltos previos: 42 = 41 + `criterios-nxx`.)
- Línea de estado en los 27: 21 RESUELTO y 6 ABIERTO (`S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-profeco-filtro`, `S5-terminos`). 42 + 21 = 63 resueltos; INDICE dice «63 resueltos y 6 abiertos»: coincide.
- `MISSING []`: cada Q2-12 a Q2-71 aparece en algún archivo de bloqueo (Q2-12 a Q2-71 = 60 preguntas más el caso 2 de Q2-01 en `S2-ventanas`).
- Tabla del INDICE contra las secciones de los archivos: 0 diferencias de letra (la fila de Q2-01 caso 2 usa otro encabezado, a ojo correcta).
- Tabla de §K46 (98 entradas) contra la hoja §8: 0 discrepancias; 15 sin respuesta: Q2-12, 48, 50 a 52, 56 a 58, 85, 87 a 91 y 94. 98 − 15 = 83 con respuesta (Q2-01 y Q2-07 con parte abierta).
- No se corrió `make verify-all`: solo documentos de estado y `decisiones.md`, sin código.

## Supuestos
- «Texto literal de esa opción tal como está en el archivo» = el texto del archivo de bloqueo, no la paráfrasis de PUERTA-2 (hay diferencias menores; las entradas Q2-13 y similares lo señalan).
- Cuando el archivo no separa la consecuencia por opción, se cita su «Impacto» y se dice así.
- Un bloqueo con alguna pregunta sin respuesta, o con decisión tomada pero sin resultado (S5-terminos), cuenta como ABIERTO.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Q2-12: la recomendada de la hoja es «C o A» (dos opciones). No la resolví (NO ASUMAS); queda ABIERTA, PROVISIONAL con `variable_estimated_cents` = null, y pasa a WF-3. Axel puede contestarla por id.
- Caso 2 de Q2-01 (ventanas contiguas que no duran un mes): sin recomendación; la pregunta 4 de `S2-ventanas` aún no está redactada, la redacta quien aplique Q2-01.
- Los 8 «sin recomendación» de bloqueos (Q2-48, 50 a 52, 56 a 58) y los de ratificaciones (Q2-85, 87 a 91, 94): PROVISIONALES, pasan a WF-3.
- S5-terminos: decisión de Axel sobre B, C o D.

## Propuestas para Axel
- PROPUESTA: contestar Q2-12 (C o A). A usa solo datos que ya capturas; C deja el gasto variable en 0 hasta que decidas.
- PREGUNTA: S5-terminos, qué sigue tras el desafío anti-robot: B (tú lees los términos y dictas el veredicto), C (pendiente, S5-07 «NO CUMPLE») o D (descartar la lectura asistida).

## Consumo aproximado
Acciones de herramienta: 19 · KB leídos: ~190 · KB escritos: ~7.
