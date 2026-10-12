# P2c · revisor-bloqueos-r1 (11-oct-2026, noche)

Agente: `revisor-bloqueos-r1`, revisor independiente, ronda 1 de 3. Informe: `docs/estado/p2c/REVISION-bloqueos-r1.md`.

## Hecho
- HECHO: revisé los 27 bloqueos tocados, `INDICE.md`, `decisiones.md` §K46–§K52 y `docs/estado/p2c/bloqueos-respuestas.md` contra RESPUESTAS.md, PUERTA-2.md (§0, §2.1, §3, §4, §5, §6, §8), el protocolo §2, §3 y §5 y el formato de la puerta 1.
- HECHO: veredicto **NO APROBADO**: 1 alto (H-r1-01: verify-c1 C1-15 y verify-s1 S1-16/K44 fallan por el cambio y el reporte no lo declara), 3 medios (H-r1-02 regla de conteo que no siguen las herramientas, H-r1-03 cita a un archivo equivocado, H-r1-04 solicitudes y secciones del reporte que faltan) y 2 bajos.
- HECHO: fidelidad sin fallas. 53 de 53 letras iguales a la hoja §8; Q2-11 = B y Q2-64 = A literales; 53 de 53 textos literales en el archivo de HEAD; ningún «sin recomendación» con respuesta; K46 con 0 discrepancias; K50 igual a PUERTA-2 §4; ningún resuelto anterior tocado; INDICE 69/63/6 igual a mi conteo con la regla nueva.

## Archivos
- Creados: `docs/estado/p2c/REVISION-bloqueos-r1.md`, `docs/estado/p2c/revisor-bloqueos-r1.md`.
- Temporales: `datos/tmp/rev-bloq-r1/` (chk.py, cnt.py, k46.py, cites.py, logs y una copia `git archive HEAD`). Borrados al final.

## Verificación
- `git diff --stat HEAD -- docs/estado/bloqueos` → 28 archivos, 665 inserciones y 4 borrados (todos en INDICE). Ningún archivo sin seguimiento en la carpeta.
- `python3 -I datos/tmp/rev-bloq-r1/chk.py` → hoja 98; 62 entradas; ninguna marca «DIF» salvo Q2-11 y Q2-64 (las explícitas) y Q2-32 (letra con nota); TEXTFAIL 0; faltan 11..71: [].
- `python3 -I datos/tmp/rev-bloq-r1/cnt.py` → total 69, resueltos 63, abiertos 6; resueltos en HEAD 42.
- `python3 -I datos/tmp/rev-bloq-r1/k46.py` → 98 entradas, discrepancias [], 15 sin respuesta.
- `make verify-c1` y `make verify-s1` (GOOSE de finanzas-hogar-orq) en el árbol → «FALLA C1-15», «S1-16 FALLA», «K44 FALLA». En la copia `git archive HEAD` → «RESULTADO: TODO OK» y rc=0. En la copia con solo los 28 archivos de bloqueos del árbol → vuelven las 3 fallas.
- `python3 -I datos/tmp/rev-bloq-r1/cites.py` → `INTEGRACION.md` l. 250 = fila `correccion-U4-r2` (cita equivocada); `REVISION-S2-r3.md` l. 45 correcta.
- No corrí `make verify-all` completo: corrí las dos partes que leen `docs/estado/bloqueos/` (verify-c1 y verify-s1).

## Supuestos
- El «protocolo §3 (bloqueos en disco)» del prompt lo leí como el formato de bloqueo del protocolo §2, más el alcance disjunto del §3 y las secciones del reporte del §5.
- La falla de H-r1-01 cuenta como «verificación global rota» aunque se cierre cuando otro agente de P2c aplique Q2-11 = B. La pruebo con la copia de HEAD y la superposición de los bloqueos.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ronda 2: comprobar H-r1-01 a H-r1-06; `make verify-c1 verify-s1` en verde cuando esté el cambio de las marcas de Q2-11.

## Solicitudes fuera de alcance
- Para el integrador: no hacer commit de `docs/estado/bloqueos/` sin el retiro de las marcas PROVISIONAL que citan `S1-parametros-avisos` (H-r1-01), y decidir quién actualiza `verify_c1.py` y `s1_verifica.py` a la regla nueva de conteo (H-r1-02).

## Propuestas para Axel
- PREGUNTA (por el orquestador): Q2-12, ¿C o A? La hoja trae las dos y RESPUESTAS.md no la lista como abierta (H-r1-06).

## Consumo aproximado
Acciones de herramienta: 33 · KB leídos: ~150 · KB escritos: ~14.
