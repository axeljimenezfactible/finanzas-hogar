# Reporte · redactor-puerta2 (Opus) · 11-oct-2026

## Hecho
- Escribí `docs/estado/puerta2/PUERTA-2.md` (61.9 KB), el documento único para que Axel decida la puerta 2. Tiene **98 preguntas, de Q2-01 a Q2-98**, sin huecos ni duplicados:
  - §2, los 6 artefactos sin aprobar: Q2-01 a Q2-10, uno por hallazgo alto o medio. El alto H-r3-00 (S5-07) se contesta con Q2-64, la pregunta del bloqueo `S5-terminos`, como pidió el encargo.
  - §3, los 27 bloqueos abiertos en el orden de INDICE: Q2-11 a Q2-71, uno por pregunta de cada archivo.
  - §4, los valores K31 y K44 de INTEGRACION §7.1: Q2-72 a Q2-77.
  - §5, lo de INTEGRACION §7.3 para ratificar: Q2-78 a Q2-94.
  - §6, propuestas de proceso: Q2-95 a Q2-98.
  - §7, acciones de Axel, y §8, la hoja de respuestas con las 98 filas.
- Cada recomendación cita su fuente: INTEGRACION §7.2 o el reporte del agente, con línea cuando la tengo. Donde no hay recomendación dice «sin recomendación» y nombra la opción que ya está en uso como PROVISIONAL.
- No modifiqué ningún otro archivo.

## Verificación
- `grep -L "Respuesta de Axel" docs/estado/bloqueos/*.md | grep -v INDICE.md | wc -l` → 28. La diferencia es `criterios-nxx.md`, que está resuelto con su «Resolución» (INDICE, sección «Resuelto»). Quedan 27 abiertos, la cifra que da INDICE.
- `python3 -I <scratchpad>/verifica_puerta2.py /home/deepsec/finanzas-hogar`, un script propio que relee los bloqueos, INDICE, los 5 informes y PUERTA-2.md:
  - bloqueos abiertos 27; encabezados en §3: 27; faltan []; sobran []; numeración 3.1–3.27 correlativa; orden igual a INDICE (P1c primero, luego la tabla de WF-1): True.
  - Q2 definidos 98, únicos 98, duplicados []; rango Q2-01..Q2-98 sin huecos; orden de aparición correlativo: True.
  - hoja de respuestas: 98 filas, 98 únicas, igual a los definidos: True.
  - hallazgos altos y medios sacados de los 5 informes: 11 (S2 M1; S5 H-r3-00 a H-r3-04; S4/S5 M-r5-1 a M-r5-3; U4 R5-1; paquete H-1). Ausentes en §2: [].
  - tamaño: 61.9 KB, debajo del tope de 120 KB.
- `git status --short`: mis archivos están solo en `docs/estado/puerta2/`.

## Supuestos
- Orden de §3: `S1-parametros-avisos` va primero porque INDICE pone la sección de P1c antes de la de WF-1. Los otros 26 siguen la tabla «Bloqueos nuevos de WF-1».
- Cuando un archivo hace varias preguntas con un solo juego de opciones, puse un id por pregunta con las mismas opciones y lo dije en el texto. Pasa en `S2-proyeccion-interruptor`, `S3-formula-tarjeta` (preguntas 1 a 3), `S5-carritos-detalle` y `S5-profeco-filtro`.
- `S5-reglas-comparativa` 4a y 4b van como 2 ids, porque tienen opciones y recomendaciones distintas (§7.2: «B en 4a», «A en 4b»).
- La «pregunta 4» de `S2-ventanas` no está en el archivo; es la corrección que pide M1. Va como Q2-01 y no la dupliqué en §3.
- Donde el archivo no pone letras, las letras son de este documento y lo digo: `S3-objetivos-medida` 2 a 4, `S3-formula-tarjeta` 4, y Q2-89 y Q2-90.
- Para Q2-07 (M-r5-2), el revisor da 3 variantes sin preferir ninguna: A queda «variante sin recomendación».
- Recomendaciones de fuera de §7.2, tomadas de los reportes y citadas: `contrato-S2.md` (S2-flujo-quincenal, S2-metas, S2-presupuesto-total, S2-proyeccion-interruptor), `contrato-S3.md` l. 51 (los 3 bloqueos S3 del contrato), `U4b-ajustes.md` l. 63-64, `U4a-ajustes.md` l. 126, `correccion-catalogo-r1.md` l. 47, `correccion-contrato-S4S5-r1.md` l. 73, `correccion-contrato-S2-r2.md` l. 63 y `correccion-U4-r2.md` l. 122.
- Q2-10: recomiendo B porque el árbol de trabajo ya trae el LEEME reescrito por K45, sin commit. Comprobé que «Crea un proyecto» da 0 coincidencias en el árbol y 1 en HEAD. La base es `k45-borrador.md`.
- Q2-85: el supuesto 5 de U3 ya se retiró en U3 1.1 (`docs/design/U3-dashboard.md` l. 459). Lo dejo como pregunta de cierre, sin recomendación.
- Agregué Q2-98 (nombres reales o seudónimos en las maquetas). Viene del supuesto declarado del orquestador en `k45-borrador.md` y no estaba en la lista (a)–(c) del encargo.
- Para §6 y §7 leí `pista3-plan-borrador.md`, y con grep las líneas de `decisiones.md` y `plan.md` que se citan. Las versiones las saqué del encabezado de DESIGN.md («Versión 1.0 en construcción») y de `C2-tokens.json` (1.0).
- La regla «sin respuesta = sigue PROVISIONAL con la opción en uso y pasa a WF-3» de la §0 es una PROPUESTA del redactor y así está marcada.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Axel contesta la hoja de la §8. Después vienen la ronda P2c y el commit «[Puerta 2]» (Q2-97).
- Avisos para el orquestador:
  - Al redactar, el árbol de trabajo ya tenía 5 archivos modificados sin commit que no son míos: `docs/plan/decisiones.md`, `docs/plan/plan.md`, `docs/criterios/P3-claude-design.md`, `docs/design/entrada/LEEME.md` y `MANIFIESTO.md`, con K45 aplicada.
  - §7.2 dice «A, A y A» para las 4 preguntas de `S4-asignacion-mercado`. Puse A en las 4 porque el reporte r2 aplica A en todas, pero conviene confirmarlo.
  - En `S1-parametros-avisos` los agentes no coinciden (B contra A). Corregido por `redactor-puerta2-fix` (H-2 de `REVISION-puerta2-r1`): el hallazgo A, alto, de `REVISION-criterios-r3` ya está cerrado en las rondas 4 y 5 (`REVISION-criterios-r4.md` l. 17; `REVISION-criterios-r5.md` l. 6, 8 y 39), y C1-22 (`docs/criterios/C1.md` l. 124) deja fuera los cuatro parámetros mientras siga abierto el bloqueo.

## Consumo aproximado
- Acciones de herramienta: unas 38.
- KB leídos: unos 195. Bloqueos unos 78; INDICE y INTEGRACION por partes unos 36; los 5 informes por partes unos 32; borradores K45 y Pista 3, 11; greps de reportes, plan y decisiones unos 38.
- KB escritos: unos 136. PUERTA-2.md 62 KB, más una reescritura completa de 62 KB al aplicar 8 correcciones de texto; el script de verificación 4 KB; este reporte unos 6 KB.
