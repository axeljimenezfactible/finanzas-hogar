# Reporte · redactor-puerta2-fix (Opus) · 11-oct-2026, noche

Encargo: corregir los 2 hallazgos (alto y medio) de `REVISION-puerta2-r1.md`, dentro del alcance de `redactor-puerta2`, sin tocar nada más ni reabrir decisiones. Los 8 bajos (B-1 a B-8) no entraban en el encargo y no se tocaron.

## Hecho

### H-1 (alta) · M1 del motor S2 incompleto
Cambios en `docs/estado/puerta2/PUERTA-2.md`:
- **§2.1:** el resumen de M1 ahora dice que la causa es una sola, que el devengo de n meses se reparte en «n ventanas consecutivas» de la lista y no en n meses (`REVISION-S2-r3.md` l. 39), y separa los **dos casos**:
  - Caso 1, el hueco que cruza un devengo, con el ejemplo del revisor (l. 40–43): $1,000.00 (100000) a 2 meses; con A + C, C da $500.00 (50000); con A, B y C, C da 0; en modo meses, diciembre da 0.
  - Caso 2, el que faltaba, las ventanas contiguas que no duran un mes (l. 45): P1 sola da $973.75 (97375) al 2026-11-09; d = 2026-10-10 más r = 2026-10-11 a 11-09 da $2,921.25 (292125) al mismo corte. Va con la cita de §3.2 de la interfaz.
- **Corrección (1):** ahora copia el texto literal de la pregunta 4 que propone el revisor (l. 48). Solo cambian las comillas internas, de «» a "", para poder anidarlas.
- **Q2-01:** la recomendada sigue siendo A, porque es la autorización de la corrección del revisor, igual que en el resto de §2. Se agregaron dos cosas:
  - Un HECHO: la recomendación de `INTEGRACION.md` §7.2 («exigir ventanas contiguas, ascendentes y disjuntas») **cubre solo el caso 1**. d y r cumplen esa regla y siguen dando $2,921.25 contra $973.75. `revisor-S2-r3.md` l. 51 lo acota («cuando las ventanas son las del oráculo»). Las reglas de ejemplo de la corrección (2) tampoco lo tocan, porque entre d y r no hay hueco.
  - **«Caso 2: sin recomendación»:** ningún reporte propone opciones. Contestar A no cierra ese caso: queda abierto dentro de la pregunta 4 del bloqueo. No se inventaron opciones.
- **§3.7:** la nota ahora dice «La pregunta 4 (devengo con ventanas: hueco que cruza un devengo y ventanas que no duran un mes) es la Q2-01» y remite a §2.1.
- **§8:** la fila de Q2-01 ahora dice «M1 devengo con ventanas (hueco y ventanas que no duran un mes)», con la recomendada «A (el caso 2 va sin recomendación)».

### H-2 (media) · HECHO falso en Q2-11
- **PUERTA-2.md §3.1:** quité el «Ojo». Ahora dice como HECHO que el hallazgo A de `REVISION-criterios-r3` está cerrado:
  - `REVISION-criterios-r4.md` l. 17 comprueba los umbrales como claves `alerts.*` y deja H-1 y H-2 pendientes.
  - La ronda 5 los cierra (`REVISION-criterios-r5.md` l. 8 y 39) y aprueba con 0 altos y 0 medios (l. 6; commit f10d3bb).
  - C1-22 (`docs/criterios/C1.md` l. 124) deja fuera los cuatro parámetros y los marca PROVISIONAL mientras siga abierto el bloqueo.
  - Las opciones A, B y C y la falta de recomendación única no cambiaron.
- **`redactor-puerta2.md`, «Avisos para el orquestador»:** corregí el aviso con las mismas citas y una nota de que lo corrigió este agente. El resto del reporte anterior quedó igual.

## Verificación
Corrí otra vez la verificación completa del agente original:
- `grep -L "Respuesta de Axel" docs/estado/bloqueos/*.md | grep -v INDICE.md | wc -l` da 28. La diferencia es `criterios-nxx.md`, que está resuelto, así que quedan 27.
- `python3 -I <scratchpad>/verifica_puerta2.py /home/deepsec/finanzas-hogar`:
  - Bloqueos abiertos: 27. Encabezados en §3: 27. Faltan []. Sobran [].
  - Numeración 3.N correlativa: True. Orden igual a INDICE: True.
  - Q2 definidos: 98. Únicos: 98. Duplicados: [].
  - Rango Q2-01..Q2-98 sin huecos: True. Orden de aparición correlativo: True.
  - Hoja: 98 filas, 98 únicas, igual a los definidos: True.
  - Hallazgos altos y medios: 11. Ausentes en §2: [].
  - Tamaño: 64.5 KB.
- `wc -c PUERTA-2.md` da 66,095 bytes (antes 63,366), por debajo de 120 KB.
- `grep -i "consumo\|KB le\|acciones de herramienta"` solo encuentra «ritmo de consumo» (l. 72, 204 y 282), que es tema de las preguntas. No hay cifras de consumo.
- Comprobaciones propias de H-1:
  - `grep -c` da 1 coincidencia para «97375» y 1 para «292125».
  - La nota vieja de §3.7 da 0 coincidencias.
  - Un script `python3 -I` toma la l. 48 de `REVISION-S2-r3.md`, cambia solo las comillas internas y la encuentra tal cual en PUERTA-2: True.
  - `grep` de las preguntas numeradas en `docs/estado/bloqueos/S2-ventanas.md` da solo la 1 y la 2 (l. 12 y 13); la 3 va con sus opciones O1/O2 (l. 44). No hay pregunta 4, así que «las opciones todavía no existen» es cierto.
  - `S2-interfaz.md` l. 16 y 24 dicen que el dinero va en centavos enteros. De ahí sale la conversión 97375 = $973.75 y 292125 = $2,921.25.
- Comprobaciones propias de H-2:
  - `sed -n 1,12p REVISION-criterios-r5.md` muestra APROBADO en l. 6 y la ronda 4 cerrada en l. 8. `grep -n C1-22` muestra H-1 cerrado en l. 39.
  - `grep -n params_json REVISION-criterios-r4.md` muestra la l. 17, «Hallazgo A… Parcial (H-1, H-2)».
  - `sed -n 124p docs/criterios/C1.md` muestra la excepción de los cuatro parámetros.
  - `git log -1 -- REVISION-criterios-r5.md` da f10d3bb.
  - «todavía describen» y «Sigue abierto el hallazgo A» dan 0 coincidencias en los dos archivos.
- Alcance:
  - `diff` contra los respaldos da 16 líneas cambiadas en PUERTA-2.md (6 viejas por 10 nuevas) y 2 en `redactor-puerta2.md` (1 por 1).
  - `git diff --stat HEAD -- docs/estado/wf1 docs/estado/bloqueos` sale vacío.
  - Los respaldos están en `<scratchpad>/p2fix/` (sha256 de antes: PUERTA-2 408a755b…53e6; reporte f25cb559…6cb68476).

## Supuestos
- Las cifras del informe están en centavos, según el contrato S2 (l. 16 y 24). Las puse en pesos, con la cifra original entre paréntesis para que se pueda rastrear. No agregué valores.
- El caso 2 no va como Q2 aparte. No hay opciones escritas para él, y crear una Q2 me habría obligado a inventarlas o a cambiar la numeración de 98. Va dentro de Q2-01 como «sin recomendación», que es una de las dos salidas que propone el revisor.
- «Las reglas de ejemplo de (2) no tocan el caso 2» es una deducción mía. La saqué de la definición de hueco de la l. 49 (`start[i] ≠ end[i−1] + 1`): d termina el 2026-10-10 y r empieza el 2026-10-11, así que entre ellas no hay hueco.
- Para H-2 cito líneas que comprobé yo: r4 l. 17 y r5 l. 6, 8 y 39. El revisor citó «r5 l. 6–9» con la frase «C1-22 (hallazgo A de `REVISION-criterios-r3`)». Esa frase no está literal en la r5 (`grep "hallazgo A"` da 0), pero el sentido sí se sostiene con las líneas citadas.

## Bloqueos
- Ninguno.

## Pendientes
- Para el orquestador, fuera de mi alcance: `docs/estado/wf1/INTEGRACION.md` §4 (l. 84) tiene la misma omisión del caso 2 de M1. No la toqué.
- Los bajos B-1 a B-8 siguen sin aplicar porque el encargo era solo alta y media. Nota: la cita de B-2 (`correccion-contrato-S4S5-r2.md` l. 101–105) resolvería el aviso sobre `S4-asignacion-mercado` del reporte anterior.
- PREGUNTA para el orquestador: ¿cuándo se le presenta a Axel la pregunta 4 completa de `S2-ventanas`, con opciones para los dos casos, una vez escrita en P2c? No está escrito.
- Aviso: `git status` ya muestra más archivos modificados fuera de `docs/estado/puerta2/` que los 5 que listan §1 y Q2-97: `README.md`, `docs/criterios/INDICE.md`, `docs/design/U1-arquitectura.md`, los PDF y `prompt.md` de `entrada/`, `salida/LEEME.md`, `tools/` y los sellos `.cache`, además de `docs/design/pista3-plan.md` sin seguimiento. Los cambió otro agente que trabaja al mismo tiempo, no yo. Se relaciona con B-7.

## Consumo aproximado
- Acciones de herramienta: unas 21.
- KB leídos: unos 80. El reporte y el informe de revisión, 18. REVISION-S2-r3 por partes, 9. PUERTA-2 por partes, 14. Criterios r3, r4 y r5 y C1 l. 124 por partes, 12. El script, 4. Greps y salidas de verificación, unos 23.
- KB escritos: el texto neto nuevo es de unos 4 KB. En bytes de disco, dentro del repositorio: 2 reescrituras completas de PUERTA-2.md (unos 132 KB), `redactor-puerta2.md` (6 KB) y este reporte (unos 7 KB), unos 145 KB. En el scratchpad: respaldos y script, unos 76 KB. Si el tope de 150 KB cuenta también el scratchpad, lo pasé por unos 70 KB, y lo declaro.
