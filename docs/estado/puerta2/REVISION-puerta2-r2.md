# REVISION-puerta2-r2 · PUERTA-2.md · 11-oct-2026, noche

Revisor: `revisor-puerta2` (independiente; ronda 2; no corregí nada). Objeto: `docs/estado/puerta2/PUERTA-2.md` (66,095 bytes) después de `redactor-puerta2-fix`, más `redactor-puerta2.md` y `redactor-puerta2-fix.md`. Base: HEAD 0926802.

## Veredicto

**APROBADO: 0 altos, 0 medios y 10 bajos.**

- Los 2 hallazgos de la ronda 1 están **cerrados**: H-1 (alta, M1 de S2 incompleto) y H-2 (media, HECHO falso en Q2-11). Los comprobé contra las fuentes, no contra el reporte del redactor.
- Los 10 bajos no bloquean. Hay 8 que vienen de la ronda 1 y siguen igual, porque la corrección no los incluía. Hay 2 nuevos, menores.

## Hallazgos de la ronda 1

| Hallazgo r1 | Estado | Evidencia propia |
|---|---|---|
| H-1 · alta · M1 de S2 incompleto y pregunta 4 reducida al hueco | **CERRADO** | §2.1 (l. 36–43) trae los **dos casos** de `REVISION-S2-r3.md` l. 38–45. Caso 1: el ejemplo de 100000 con A + C da 50000, con A, B y C da 0, y en modo meses diciembre da 0 (l. 40–43). Caso 2: `[P1]` da 97375 y `[d, r]` da 292125 al 2026-11-09 (l. 45). Las cifras están en pesos y la cifra original va entre paréntesis: 97375 = $973.75 y 292125 = $2,921.25. La corrección (1) copia la pregunta 4 del revisor (l. 48) tal cual; solo cambian las comillas internas. Q2-01 trae un HECHO: la recomendación de `INTEGRACION.md` §7.2 cubre solo el caso 1. La cita de `revisor-S2-r3.md` l. 51 es correcta (`grep -n` da 51). El caso 2 va como «sin recomendación», sin opciones inventadas, y deja claro que contestar A no lo cierra. §3.7 (l. 144) dice ahora «hueco que cruza un devengo y ventanas que no duran un mes». La fila de Q2-01 en §8 dice «A (el caso 2 va sin recomendación)». Las correcciones (2) a (4) coinciden con la l. 49–53 del informe. |
| H-2 · media · Q2-11 con un HECHO falso sobre los criterios | **CERRADO** | §3.1 ya no tiene el «Ojo». Ahora dice como HECHO que el hallazgo A de r3 está cerrado. Comprobé las líneas que cita: `REVISION-criterios-r4.md` l. 17 («Parcial (H-1, H-2)», umbrales como `alerts.*`); `REVISION-criterios-r5.md` l. 6 («APROBADO. 0 hallazgos altos, 0 medios y 2 bajos»), l. 8 (los 6 hallazgos de r4 cerrados) y l. 39 (H-1 cerrado); `docs/criterios/C1.md` l. 124 («salvo los de programación de avisos», con «excepción PROVISIONAL, bloqueo abierto»). `git log -1` de la r5 da f10d3bb. El aviso de `redactor-puerta2.md` (l. 46) también está corregido. `grep -c` de «todavía describen» y de «Sigue abierto el hallazgo A» da 0 en los dos archivos. Las opciones A, B y C de Q2-11 siguen siendo fieles al archivo, y las dos recomendaciones citadas existen (`correccion-S1-r2.md` l. 59 recomienda B; `correccion-criterios-r1.md` l. 36 recomienda A). |

## Verificaciones (HECHO, con comandos propios)

| # | Comprobación | Resultado | Evidencia |
|---|---|---|---|
| a | Cobertura | CUMPLE | `grep -L "Respuesta de Axel" docs/estado/bloqueos/*.md`, sin contar INDICE, da 28 archivos. El que sobra es `criterios-nxx`, que está resuelto («Estado: RESUELTO el 10-oct-2026», l. 2). Los encabezados `### 3.N` dan 27 nombres, y `diff` contra la lista de abiertos solo marca `criterios-nxx`. No se coló ningún bloqueo resuelto. Conté las preguntas de cada archivo («Qué falta» y las preguntas añadidas): son 61, y Q2-11 a Q2-71 también son 61, con el mismo reparto por archivo. Por ejemplo: S2-ventanas tiene 3 (la 4 todavía no existe en el archivo), S3-formula-tarjeta 4 (con la «adicional»), S5-reglas-comparativa 5 (1, 2, 3, 4a y 4b) y U4-vista-previa-formulas 2 (la 2 está retirada). |
| b | Fidelidad | CUMPLE (solo matices bajos) | Leí y comparé pregunta y opciones de 10 bloqueos: S5-terminos, S2-ventanas, S1-parametros-avisos, S3-serie-sin-intereses, S4-asignacion-mercado, S3-formula-tarjeta, S3-objetivos-medida, S5-reglas-comparativa, U4-vista-previa-formulas y S2-proyeccion-interruptor. También revisé cómo se reparten las opciones en S5-carritos-detalle y S5-profeco-filtro, y la pregunta de los otros 16. No cambia ningún sentido y no hay opciones inventadas. Las cifras coinciden: 773.62 ×3 + 773.64 y 773.63 ×3 + 773.61; 1,191.50 y 1,499.72; 101/102 y 103/104 de H-r3-02; 37,548,200. Las letras que pone el documento se declaran como suyas (Q2-32, Q2-34 a Q2-36, Q2-89 y Q2-90). |
| c | Recomendaciones | CUMPLE | Las de §7.2 coinciden una por una: S2-ventanas A + H1 + O1; S2-proyeccion C; S2-flujo C o A, B y A; S2-metas-resumen A; S3-origen (a); S3-serie (a); S4-reglas A; S5-reglas A, A, A, B y A; S5-profeco-almacen A; S5-carrito 1 A y 3 A; S4-asignacion A; U3-t2 A; U4-vista-previa (a); S5-terminos A o B, sin recomendación única. Para la pregunta 4 de S2-ventanas, el documento reporta §7.2 de forma transparente y dice qué caso no cubre. Revisé en su línea las citas de fuera de §7.2: `correccion-S1-r2` l. 59, `correccion-criterios-r1` l. 36, `contrato-S2` l. 59–62, `contrato-S3` l. 51, `correccion-contrato-S2-r2` l. 63, `U4b-ajustes` l. 63 y 64, `U4a-ajustes` l. 126, `correccion-catalogo-r1` l. 47, `correccion-contrato-S4S5-r1` l. 73, `correccion-U4-r2` l. 122 y `U3-dashboard` l. 51. Un script comprobó que existen todos los reportes citados entre comillas invertidas. Solo faltan archivos que todavía no tienen por qué existir: `gastos_fijos_ingresos.csv` (lo entrega Axel), `vinculos.csv` (está como opción) y `ops/mkcert.md` (ver B-10). La columna «Recomendada» de la §8 coincide con el cuerpo en las 98 filas. |
| d | 6 artefactos de §4 | CUMPLE | Altos y medios según el veredicto y los encabezados de los 5 informes: S2-r3 M1 (ahora completo); S5-r3 H-r3-00 [alta] y H-r3-01 a H-r3-04 [media] (l. 39–67); contrato-S4S5-r5 M-r5-1 a M-r5-3 (l. 6 dice «0 altos, 3 medios»); U4-r5 R5-1 (l. 6); paquete-r3 H-1 (l. 6 y 34). Son 11 de 11, y cada uno tiene su bloque y su Q2. Las correcciones de M-r5-1 (l. 94), M-r5-2 (l. 103–107: 3 variantes sin preferida), M-r5-3 (l. 119–125) y R5-1 (l. 50–57, con la «otra opción») coinciden con el texto del revisor. |
| e | Tabla K31/K44 | CUMPLE | Los 6 renglones de §4 (Q2-72 a Q2-77) coinciden con `INTEGRACION.md` §7.1 (l. 135–143) en valor, razón y agente. Solo cambia la redacción. |
| f | Numeración, hoja y tamaño | CUMPLE | Con un script `python3 -I`: hay 98 definiciones, únicas y correlativas del 1 al 98. La hoja de §8 tiene 98 filas únicas, en el mismo orden, con la columna de respuesta vacía. Ninguna mención pasa de Q2-98. `wc -c` da 66,095 bytes, por debajo de 120 KB. `grep -i "consumo\|KB le\|acciones de herramienta"` solo encuentra «ritmo de consumo» (l. 72, 204 y 282), que es tema de las preguntas. No hay cifras de consumo. |
| g | HECHO / PREGUNTA / PROPUESTA y K45 | CUMPLE | Las etiquetas se mantienen: §1 es HECHO, §2 y §3 son PREGUNTAS y §4 a §6 son PROPUESTA. La regla de §0 sobre lo que no se conteste está marcada como PROPUESTA del redactor. K45 va como «Lo que ya decidiste hoy… sustituye a K35», y `decisiones.md` l. 105 y 115 lo confirman en el árbol. Q2-95 cita bien la l. 105. Q2-98 pregunta por el supuesto declarado de K45.1 y no reabre K45. Los HECHO de Q2-10 son ciertos: «Crea un proyecto» da 0 en el árbol y 1 en HEAD. Los de Q2-96 también: DESIGN l. 3 dice «Versión 1.0 en construcción» y `C2-tokens.json` dice `"version": "1.0"`. El HECHO nuevo de Q2-11 es cierto. El nuevo de Q2-01 («d y r no dejan hueco») sale de la definición de hueco de la l. 49 (`start[i] ≠ end[i−1] + 1`), y la cuenta da bien. |

## Bajos (no bloquean)

Siguen abiertos desde la ronda 1, porque la corrección solo cubría el alto y el medio:
- **B-1 · (d).** §1 (l. 21) y el encabezado de §2 dicen «6 artefactos», pero hay 5 bloques (2.1 a 2.5). INTEGRACION §0 y §4 tienen el mismo desfase.
- **B-2 · (c).** §3.13 cita `correccion-contrato-S4S5-r2.md` l. 82. La recomendación explícita de las 4 preguntas está en la l. 101–105 («Propuestas para Axel»).
- **B-3 · (b).** Q2-34: la alternativa «serie de la compra sola» no tiene letra.
- **B-4 · (b).** Faltan matices que no cambian el sentido: Q2-25 B no dice «o con meses calendario»; Q2-40 B no dice «sin precio en ninguna cadena, como A»; Q2-60 C no dice «(mercado en su día)». Siguen también los demás que listó la r1.
- **B-5 · (c).** Q2-97 dice «propuesta del orquestador» sin archivo fuente.
- **B-6 · (g).** §1, K45: dice «los agentes revisan las vueltas, con un máximo de 3». `k45-borrador.md` dice «agentes Sonnet revisan las vueltas 1 y 2; máximo 3 iteraciones».
- **B-7 · (g), ampliado.** §1 y Q2-97 listan 5 archivos modificados por K45. Ahora `git status --short` muestra 22 archivos modificados fuera de `docs/estado/puerta2/`, más `docs/design/pista3-plan.md` sin seguimiento. Entre los modificados están `README.md`, `docs/criterios/INDICE.md`, `docs/design/U1-arquitectura.md`, `docs/design/tools/*`, los 4 PDF y `prompt.md` de `entrada/`, `salida/LEEME.md` y 6 sellos `.cache`. Sus mtime van de 21:56 a 21:58, así que otro agente los está cambiando al mismo tiempo. Corrección: que P2c parta de `git status` y no de la lista de §1.
- **B-8 · informativo.** `revisor-S4-r2.md` l. 130 propone ratificar SR Legacy 2018-04 (USDA, CC0). No está ni en §7.3 ni en PUERTA-2.

Nuevos:
- **B-9 · (d).** El resumen de H-r3-03 (§2.2) no menciona que el motor emite `loaded_proposed`, una clave que §15.2 no lista (`REVISION-S5-r3.md`, H-r3-03). La corrección del revisor tampoco la nombra, así que no falta ningún hallazgo.
- **B-10 · (g).** §7, punto 4, dice que los pasos «van en `ops/mkcert.md`». Ese archivo todavía no existe: `plan.md` l. 199 y l. 382 lo prevén para WF-4. Corrección opcional: «irán».

## PREGUNTA (para el orquestador, no para Axel)
- `INTEGRACION.md` §4 (l. 84) sigue sin el caso 2 de M1. PUERTA-2 ya lo corrige con la fuente. ¿Se corrige también INTEGRACION, o queda como historia de WF-1? No está escrito.
- ¿En qué momento se le presenta a Axel la pregunta 4 completa de `S2-ventanas`, con opciones para los dos casos, una vez que se escriba en P2c? Tampoco está escrito. Lo levanta también `redactor-puerta2-fix.md`.

## Comandos usados (resumen)
- `grep -L "Respuesta de Axel" docs/estado/bloqueos/*.md | grep -v INDICE.md`, más `diff` contra `grep -oE '^### 3\.[0-9]+ \`[^\`]+\`' PUERTA-2.md`.
- `python3 -I` (heredoc): definiciones Q2, conteo por sección, hoja de §8, comparación de «Recomendada» entre cuerpo y hoja, y existencia de los archivos citados.
- `sed -n '/^Qué falta/,/^Opciones/p'` sobre los 18 bloqueos que no leí completos, y `cat` o `sed` de los 10 de la muestra.
- `sed -n` de `REVISION-S2-r3.md` l. 1–12 y 36–60, `REVISION-S5-r3.md` l. 37–68, `REVISION-contrato-S4S5-r5.md` l. 83–130, `REVISION-U4-r5.md` l. 37–57, e INTEGRACION l. 78–89 y 133–158.
- `sed -n` de las líneas citadas (11 reportes), `grep -n` en `revisor-S2-r3.md`, `REVISION-criterios-r4/r5`, `C1.md` l. 124 y `git log -1`.
- `git status --short` y `stat` de los archivos modificados.
