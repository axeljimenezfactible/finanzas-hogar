# Reporte · revisor-puerta2 (ronda 1) · 11-oct-2026, noche

## Hecho
- Revisé `docs/estado/puerta2/PUERTA-2.md` contra los bloqueos, INTEGRACION §4, §6 y §7, los 5 informes REVISION y los reportes citados. No corregí nada.
- Escribí el informe `docs/estado/puerta2/REVISION-puerta2-r1.md`. Veredicto: **NO APROBADO**, con 1 alto, 1 medio y 8 bajos.
  - H-1 (alta): el M1 del motor S2 aparece incompleto. Falta el caso de ventanas contiguas que no duran un mes ([P1] da 97375 y [d, r] da 292125 al 2026-11-09), y la pregunta 4 de `S2-ventanas` queda reducida al hueco. La recomendación «contiguas» no cierra ese caso.
  - H-2 (media): Q2-11 dice que sigue abierto el hallazgo A, alto, de `REVISION-criterios-r3`, pero las rondas 4 y 5 lo cerraron (`REVISION-criterios-r5.md`: APROBADO; `C1.md` l. 124).
  - Bajos: el conteo «6 artefactos» con 5 bloques; la cita de S4-asignacion-mercado (debe ser l. 101–105, no l. 82); una alternativa sin letra en Q2-34; matices omitidos en 6 opciones; Q2-97 sin archivo fuente; la precisión de K45 sobre las vueltas; `docs/design/pista3-plan.md` sin mencionar; y B2 de `revisor-S4-r2` sin escalar, heredado de §7.3.
- No modifiqué ningún otro archivo.

## Verificación
- `grep -L "Respuesta de Axel" docs/estado/bloqueos/*.md | grep -v INDICE.md | wc -l` → 28. `criterios-nxx.md` l. 2 dice «Estado: RESUELTO el 10-oct-2026»; quedan 27.
- §3 de PUERTA-2 (nombres de los `### 3.N`), ordenada y comparada con `diff` contra los 27 → sin diferencias. Bloqueos resueltos mencionados: solo `criterios-nxx`, como resuelto.
- Script `python3 -I` sobre PUERTA-2 → 98 definiciones, correlativo True, hoja de §8 con 98 filas únicas e igual a las definiciones, y 98 respuestas vacías. `wc -c` → 63366.
- Fidelidad: leí «Qué falta» y «Opciones» de los 27 bloqueos y los comparé con su Q2. No encontré ninguna opción inventada ni ningún sentido cambiado en los bloqueos. El único cambio de alcance está en §2.1 y §3.7 (H-1).
- Altos y medios de los 5 informes: con `grep` de severidad, y `sed -n` de los bloques de hallazgos → 11. Los 11 están en §2.
- Citas: comprobé con `sed -n` las líneas de los 11 reportes citados, y todas existen y dicen lo que PUERTA-2 atribuye. `correccion-contrato-S4S5-r2.md` l. 101–105 trae la recomendación A de las 4 preguntas de S4-asignacion-mercado.
- H-2: `sed -n '6,12p' docs/estado/p1c/REVISION-criterios-r5.md` → «APROBADO. 0 hallazgos altos, 0 medios y 2 bajos» y «C1-22 (hallazgo A de REVISION-criterios-r3)». `git log -1 -- …r5.md` → f10d3bb.
- HECHO de Q2-10: `grep -c 'Crea un proyecto'` → 0 en el árbol y 1 en `git show HEAD:docs/design/entrada/LEEME.md`.
- HECHO de Q2-95 y Q2-96: `decisiones.md` l. 105 trae «WF-3 no espera a la Pista 3 (plan R4…)»; `DESIGN.md` l. 3 dice «Versión 1.0 en construcción»; `C2-tokens.json` trae `"version": "1.0"`.

## Supuestos
- Califiqué H-1 como alto porque la regla dice «hallazgo faltante o sentido cambiado». Lo que falta es una parte del hallazgo M1, con cifras que la regla recomendada no corrige. Si el orquestador prefiere tratarlo como medio, el veredicto no cambia.
- Califiqué H-2 como medio: es un HECHO falso dentro de una pregunta, pero no cambia la pregunta ni sus opciones.
- El orden de §3 (según INDICE) no lo pedía el encargo y no lo revisé a fondo.

## Bloqueos
- Ninguno nuevo. Dejo una PREGUNTA para el orquestador en el informe: ¿los «6 artefactos» de INTEGRACION cuentan el contrato S4/S5 como dos?

## Pendientes
- El redactor corrige H-1 y H-2 (y los bajos que el orquestador quiera), y después viene una ronda 2 de revisión.
- Corregir también el aviso sobre `REVISION-criterios-r3` en `redactor-puerta2.md`.

## Consumo aproximado
- Acciones de herramienta: unas 42.
- KB leídos: unos 235. PUERTA-2 62; INTEGRACION por partes 22; los 27 bloqueos, casi todos en su parte de preguntas y opciones, 70; los 5 REVISION por partes 30; reportes citados y propuestas 25; criterios r5, C1 y K45 10; el resto, salidas de comandos.
- KB escritos: unos 18. El informe 12.8 y este reporte unos 4.5.
