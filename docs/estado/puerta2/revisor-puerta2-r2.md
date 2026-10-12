# Reporte · revisor-puerta2 (ronda 2) · 11-oct-2026, noche

Encargo: evaluar `docs/estado/puerta2/PUERTA-2.md` después de `redactor-puerta2-fix` y comprobar que H-1 (alta) y H-2 (media) de `REVISION-puerta2-r1.md` están cerrados. Revisor independiente: no corregí nada.

## Hecho
- Escribí el informe `docs/estado/puerta2/REVISION-puerta2-r2.md`. Veredicto: **APROBADO, con 0 altos, 0 medios y 10 bajos**.
- H-1 está CERRADO. §2.1 trae los dos casos de M1, con cifras y la cita de l. 45. Copia la pregunta 4 literal (l. 48). Q2-01 aclara que la recomendación de §7.2 no cubre el caso 2, que va sin recomendación y sin opciones inventadas. §3.7 y §8 están alineados.
- H-2 está CERRADO. El HECHO de Q2-11 ahora es cierto (criterios r4 l. 17; r5 l. 6, 8 y 39; C1 l. 124; f10d3bb), y el aviso del reporte del redactor también se corrigió.
- Volví a comprobar de la (a) a la (g) con mis propios comandos. Muestra de fidelidad: 10 bloqueos completos, más la pregunta de los otros 18. Revisé los 11 altos y medios de los 5 informes, la tabla de §7.1 y las 98 filas de la hoja.

## Verificación (comandos y salidas)
- `grep -L "Respuesta de Axel" docs/estado/bloqueos/*.md | grep -v INDICE.md | wc -l` da 28. Contra los encabezados de §3, `diff` solo marca `criterios-nxx.md` (RESUELTO, l. 2), así que quedan 27 abiertos y los 27 aparecen.
- `python3 -I`: «defs 98 unique 98 correlativo True»; «hoja 98 98 True»; ninguna mención pasa de 98; 66,095 bytes. Por sección: Q2-11 a Q2-71 dan 61, y la cuenta de preguntas en los archivos también da 61.
- La comparación de «Recomendada» entre cuerpo y hoja coincide en las 98 filas. Las filas sin recomendación propia caen bajo «A en las tres», «A en las cuatro» o el encabezado de la tabla de §4.
- En el script de existencia de archivos citados solo faltan `gastos_fijos_ingresos.csv`, `vinculos.csv` y `ops/mkcert.md`. Los tres son entregables futuros, no reportes.
- `grep -n "Pregunta 4 de" wf1/revisor-S2-r3.md` da 51. `sed -n 6,9p REVISION-criterios-r5.md` da «APROBADO. 0 hallazgos altos, 0 medios y 2 bajos». `git log -1` de la r5 da f10d3bb.
- `grep -c "Crea un proyecto"`: 0 en el árbol y 1 en HEAD. `decisiones.md` l. 105 dice «Sustituida por K45» con «WF-3 no espera a la Pista 3»; la l. 115 es K45.
- `grep -i "consumo\|KB le\|acciones de herramienta" PUERTA-2.md` solo encuentra «ritmo de consumo» (l. 72, 204 y 282).
- `git status --short` da 22 archivos modificados fuera de `docs/estado/puerta2/` y `docs/design/pista3-plan.md` sin seguimiento. Sus mtime van de 21:56 a 21:58.

## Supuestos
- Igual que en la ronda 1, conté los altos y medios con el veredicto y los encabezados de cada informe REVISION.
- Juzgué la fidelidad por el sentido y las opciones. Las omisiones de matices que no cambian la decisión las dejé como bajas (B-4), como en la r1.
- El HECHO derivado de Q2-01 («d y r no dejan hueco») lo acepto porque sale de la definición de hueco de `REVISION-S2-r3.md` l. 49, y la cuenta da bien.

## Bloqueos
- Ninguno.

## Pendientes
- Para el orquestador: `INTEGRACION.md` §4 (l. 84) sigue sin el caso 2 de M1. Además, falta decidir cuándo se presenta la pregunta 4 completa de `S2-ventanas`.
- Los bajos B-1 a B-10 están en el informe. Ninguno bloquea.
- Aviso: el encargo decía que el árbol estaba limpio salvo `docs/estado/puerta2/`. No es así. `git status` muestra cambios de otro agente, hechos al mismo tiempo, en `README.md`, `docs/criterios/INDICE.md`, `docs/design/U1-arquitectura.md`, `docs/design/tools/`, `entrada/`, `salida/LEEME.md`, `docs/plan/` y `docs/criterios/P3-claude-design.md`, además de `docs/design/pista3-plan.md`, que es nuevo. Ver B-7. No toqué nada de eso.

## Consumo aproximado
- Acciones de herramienta: unas 40.
- KB leídos: unos 165. PUERTA-2 por partes, 60. Informes REVISION por partes, 30. Bloqueos, 40. INTEGRACION §4 y §7, 9. Reportes del redactor y la r1, 27. Greps y salidas de scripts, unos 10.
- KB escritos: unos 15 (este reporte e informe). Los archivos auxiliares del scratchpad suman menos de 1 KB.
