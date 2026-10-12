# REVISIÓN K45 · ronda 1 (revisor-k45, independiente) · 11-oct-2026, noche

Objeto: el trabajo de `aplicar-k45` (reporte `docs/estado/puerta2/aplicar-k45.md`) contra `docs/estado/puerta2/k45-borrador.md`, `docs/plan/protocolo-agentes.md` §4 y la coherencia interna del repositorio. No corregí nada. Todas las comprobaciones las hice con mis propios comandos (salidas en `docs/estado/puerta2/revisor-k45-r1.md`).

**Veredicto: NO APROBADO** · 7 altos, 4 medios, 5 bajos.

Lectura justa: lo que `aplicar-k45` cambió dentro de su alcance está bien hecho. K45 es idéntica al borrador; K35 está marcada y conserva su texto; P3 1.3 conserva ids y severidades; LEEME ya no tiene pasos de subida; el MANIFIESTO coincide; `pista3-plan.md` conserva los datos técnicos; `git status` está limpio fuera del alcance. No se aprueba por dos razones. Primero, `make verify-criterios` y `make verify-all` no pasan. Segundo, quedan en el repositorio instrucciones vigentes del flujo K35, con subida de archivos o navegador, en archivos que `aplicar-k45` no tocó. Una parte ya la reportó en sus Pendientes. **No reportó**: `plan.md` línea 334; `INDICE.md` línea 9; `README.md`; `U1-arquitectura.md` junto con `verificar_u1.py`; `decisiones.md` §J; dónde va el informe de revisión según K45; y las cifras de Kueski 973.75 y 883.77.

## 1. Resultado por comprobación pedida

| # | Comprobación | Resultado |
|---|---|---|
| a | K45 idéntica al borrador; K35 sustituida sin borrar | **PASA**. Las 5 líneas de K45 se extrajeron del borrador y de `decisiones.md` y el `diff` da 0 diferencias. Sin el prefijo «**Sustituida por K45 (11-oct-2026).** », K35 es byte a byte la de HEAD. |
| b | `plan.md` §2 filas 41/45/63, §3.6, §3.8 y R4 coherentes con K45 y con 42, 43 y 44 | **PASA en esas secciones**. Las tres filas llevan la nota de sustitución; §3.6 cubre K45.1–K45.3; §3.8 refleja K45.3; R4 queda marcado obsoleto. No hay contradicción con 42 (los materiales siguen sin código de la app), 43 (referencia visual) ni 44 (3 vueltas). **Pero** el mismo archivo contradice K45 en las líneas 206, 334 y 409 (H-02). |
| c | P3: ids, severidades, medibles, sin subida ni claude.ai/design, P3-12, verify-criterios | **FALLA por verify-criterios** (H-01). Los ids P3-01 a P3-13 y las severidades son iguales a HEAD (10/3/0); hay 13 Verifica y 13 Evidencia. Las menciones de subir o de claude.ai/design son negativas o históricas. P3-12 trae la lista de 8 cifras canon; su cobertura está incompleta (M-04). |
| d | LEEME sin pasos de subida; MANIFIESTO coincide | **PASA en lo pedido**. `sha256sum -c` da exit 0, con 22 elementos y 21 filas sha. Las descripciones del MANIFIESTO siguen en K35 (M-01). |
| e | `pista3-plan.md` conserva los datos técnicos del borrador | **PASA**. 56 de 57 fichas literales del borrador (rutas, URL, releases, números) aparecen; la que falta es la hora «19» del encabezado. Falta la palabra «ultracode» (B-02). |
| f | `make verify-all GOOSE=…` termina en OK | **FALLA**: exit 2. La causa es la misma que en H-01. Con `-k`: WF-0, C1, S1, U1 y WF-1 en OK; falla solo criterios. |
| g | `git status` solo con archivos del alcance | **PASA**. Cambiaron `P3-claude-design.md`, `entrada/LEEME.md`, `entrada/MANIFIESTO.md`, `decisiones.md` y `plan.md`; son nuevos `docs/design/pista3-plan.md` y `docs/estado/puerta2/`. |

Protocolo §4 (reparto fino) para `aplicar-k45`: declara unas 31 acciones de 60, unos 120 KB leídos de 400 y unos 90 KB escritos de 150. Ningún volcado pasa de 40 KB (el mayor es P3, con 26.5 KB). Los respaldos existen en el scratchpad (`k45/`). No hay incumplimiento.

## 2. Hallazgos

### Altos

**H-01 · verify-criterios y verify-all fallan.** `docs/criterios/INDICE.md` línea 149: `| P3-claude-design.md | 13 (10 / 3 / 0) | 1.2 |` choca con la versión 1.3 del archivo. La salida es `DIFERENCIA INDICE §10 P3-claude-design: ('13', '10', '3', '0', '1.2') contra (13, 10, 3, 0, '1.3')` y luego `make: *** [Makefile:159: verify-criterios] Error 1`. Corrección: cambiar 1.2 por 1.3 en esa fila, añadir la fila K45 en la tabla de decisiones (cerca de la línea 160) y volver a correr `make verify-criterios` y `make verify-all GOOSE=…`.

**H-02 · `plan.md` sigue dando el flujo de navegador fuera de las secciones editadas.**
- Línea 206 (§3.3): «PUERTA 3: Axel inicia sesión en claude.ai/design y comparte carpeta ──► Pista 3 (orquestador + 2 agentes)».
- Línea 409 (§8.2, acciones de Axel): «Puerta 3: iniciar sesión en claude.ai/design en Chrome y compartir la carpeta de entrada con la sesión.»
- Línea 334 (presupuesto por agente): «Pista 3 | … (orquestador ejecuta el navegador) | … | 2 | … 200K (+~300K sesión principal)». Choca con K45 («sigue sin navegador») y con la estimación de 7 a 9 M de `pista3-plan.md` §3.

Las tres contradicen K45 y §3.6 y §3.8 del mismo archivo. Las líneas 206 y 409 están en los Pendientes de `aplicar-k45`; la 334 no. Corrección (PROPUESTA): en 206 y 409, «PUERTA 3 (al final de la Pista 3): Axel abre la URL del lienzo y aprueba»; en 334, «sin navegador; agentes según `pista3-plan.md` §3».

**H-03 · `docs/criterios/INDICE.md` mantiene K35 como vigente.**
- Línea 9: «… subir el paquete a Claude Design, iterar y avisar «iter-N lista»: la Pista 3 la ejecuta Axel, K35». No está en los Pendientes de `aplicar-k45`.
- Línea 28: «(paquete: orquestador; Pista 3: Axel, K35)».
- Línea 44: puerta 3 = «paquete de entrada listo y revisado», con P3-01, P3-02 y P3-11. En P3 1.3, P3-01 ya es la aprobación final de Axel.

Corrección: reescribir las tres líneas con K45. La puerta 3 queda como «al final de la Pista 3: Axel aprueba la referencia visual en la URL del lienzo». El orquestador elige qué criterios entran a esa fila (PROPUESTA de `aplicar-k45`: P3-01, P3-06 y P3-12).

**H-04 · `docs/design/entrada/prompt.md` contradice K45.1 y K45.2.**
- Línea 101: «No inventes datos financieros. Usa solo las cifras que aparecen en las capturas `reporte-p01.png` a `reporte-p11.png` … no pongas montos inventados». Es lo contrario de K45.1 («Cifras sintéticas en las maquetas»).
- Línea 12: «× 4 combinaciones = 84 artboards», con tema en el nombre. Choca con K45.2 (claro y oscuro por el menú de tema; todas las pantallas de U1).

`LEEME.md` tiene una nota de precedencia, pero `plan.md` §3.6 dice que los agentes leen `prompt.md` como brief. Corrección: reescribir `prompt.md` antes de la Fase 2 de la Pista 3 y regenerar `MANIFIESTO.md` con `paquete_entrada.sh`.

**H-05 · `docs/design/salida/LEEME.md` da instrucciones K35 vigentes a Axel.**
- Línea 3: «Aquí guarda Axel lo que descargue de Claude Design; el paquete que se sube está en …».
- Líneas 10 y 11: «después de pegar `iter-1/siguiente-prompt.md`».
- Línea 24: `siguiente-prompt.md` «se pega como mensaje nuevo».

Corrección: reescribirlo con las rutas de K45 (copia de cada vuelta, `URLS.md`, `final/`) o retirarlo.

**H-06 · `README.md` del repositorio dice que Axel sube el paquete.**
- Línea 57: «paquete para Claude Design …; lo sube Axel (K35)».
- Línea 156: «Pista 3 (Claude Design), la ejecuta Axel (K35). Axel sube el paquete … La puerta 3 queda en «paquete listo y revisado»».

Corrección: actualizar las dos líneas con K45.

**H-07 · `docs/design/U1-arquitectura.md` contradice K45.**
- Línea 27: «La Pista 3 (Claude Design) la ejecuta Axel con el paquete … | K35 | §7».
- Línea 1210: «**Decidido** (K35): lo ejecuta Axel con el paquete de `docs/design/entrada/` …».

Hay un acoplamiento: `docs/design/tools/verificar_u1.py` línea 295 exige `'K35' in s7`. Si se cambia U1 §7 sin cambiar el verificador, se rompe `make verify-u1`. Corrección: cambiar U1 §7 y el verificador en el mismo paso, con la revisión que corresponda a U1, que lleva versión y revisor propio.

### Medios

**M-01 · Las descripciones de `MANIFIESTO.md` siguen en K35.**
- Línea 5: «Carpeta de entrada (lo que se sube)».
- Línea 14 (`LEEME.md`): «Pasos para Axel: qué subir, qué pegar, dónde guardar lo que baje…».
- Línea 15 (`prompt.md`): «Mensaje que se pega en Claude Design».

Los sha256 coinciden, pero el texto contradice `LEEME.md` y K45. Lo genera `docs/design/tools/paquete_entrada.sh` (líneas 181, 182 y 203). Corrección: actualizar esas cadenas en el script y regenerar. Como advierte `aplicar-k45`, cambiar el script cambia los sellos del caché, se reconstruyen los PDF y cambian sus sha256. Conviene hacerlo junto con H-04.

**M-02 · `decisiones.md` §J punto 3 no está marcado.** La línea 66 dice: «3. Antes de Pista 3: Axel inicia sesión en claude.ai/design y comparte `docs/design/entrada/`», sin nota. P3 §5 (línea 132) afirma «ambos sustituidos por K45». K45 sustituye esa puerta («la puerta 3 pasa a ser…»), pero el registro no lo dice donde está el texto. Corrección: añadir «→ sustituido por K45 (11-oct-2026)» sin borrar el texto. Necesita el OK del orquestador, porque `decisiones.md` es el registro de Axel.

**M-03 · Dónde va el informe de revisión de cada vuelta.** K45 (consecuencias) dice: «`docs/design/salida/iter-N/` guarda las copias en disco de los artboards de cada vuelta y el informe de revisión». Sin embargo, el informe solo va en `docs/estado/P3-iter-N-<area>.md` en estos lugares:
- P3-06 (línea 67), que copia solo `canvas.json` y los artboards;
- P3-07 (línea 71);
- `plan.md` §3.6, paso 5;
- `pista3-plan.md` §5 («Informes de revisión»).

P3 §5 traza esa frase de K45 a P3-06 y P3-07, pero ninguno de los dos la exige. El origen es una diferencia entre K45 y el borrador técnico del propio orquestador. PREGUNTA: ¿el informe va en los dos lugares? PROPUESTA mínima, sin tocar K45: P3-06 exige además una copia del informe (o un índice que apunte a él) en `iter-N/`.

**M-04 · P3-12 no cubre todas las cifras de Kueski.** P3-12 (línea 91) prohíbe 8 cifras. `plan.md` línea 218 (S3, regresión Kueski) lista además **973.75**, **1,258.70** y **883.77**, y K19 (`decisiones.md` línea 89) cita 1,258.70. El criterio dice «ni cifras del reporte real o de Kueski», pero su Verifica no detectaría esas tres. Tampoco comprueba fechas reales: K45.1 habla de «montos, fechas y saldos». `aplicar-k45` preguntó solo por 1,258.70. Corrección: añadir `973\.75|1,?258\.70|883\.77` al grep, o que el orquestador deje escrito por qué se excluyen. Para las fechas, comparar `datos-sinteticos.json` contra las fechas de `data.json`, o declararlas fuera de la comprobación automática.

### Bajos

**B-01** · La Verifica de P3-02 lee `iter-1/project/ds/finanzas-hogar/tokens.json`, pero P3-06 solo exige copiar `canvas.json` y los artboards. Corrección: añadir `project/ds/` a la lista de P3-06.

**B-02** · `pista3-plan.md` §3 omite «ultracode», que el borrador pone en el encabezado del flujo WF-P3.

**B-03** · `pista3-plan.md` línea 64 dice «512 archivos y 256 MB por lienzo», fiel al borrador. La descripción de la herramienta Artifact que ve este revisor dice «a version may hold up to 511 files and 256MB». PREGUNTA al orquestador: confirmar el límite.

**B-04** · En `plan.md` R4 (línea 383), la mitigación vieja («carpeta compartida en puerta 3») va antes de la nota «Obsoleto desde K45», así que se lee como vigente hasta el final de la celda. Sugerencia: poner «OBSOLETO (K45)» al inicio de la celda.

**B-05** · No está claro quién guarda la exportación final. `pista3-plan.md` §3 dice «la descarga Axel o la pide el orquestador si la herramienta la da como archivo». `LEEME.md` (paso 4) y `plan.md` §3.6 (paso 7) no dicen quién. Si es Axel, es un paso manual que su guía no menciona. PREGUNTA.

## 3. Qué falta para aprobar (ronda 2)

1. H-01: INDICE §10 a 1.3 y la fila K45; que `make verify-criterios` y `make verify-all GOOSE=…` den OK.
2. H-02, H-03, H-05, H-06 y H-07: quitar las instrucciones K35 vigentes. En H-07, junto con `verificar_u1.py`.
3. H-04 y M-01: reescribir `prompt.md`, actualizar las cadenas de `paquete_entrada.sh` y regenerar el MANIFIESTO (cambian los sha256 de los PDF).
4. M-02, M-03 y M-04: marcar §J.3; resolver dónde va el informe por vuelta; completar o justificar las cifras de P3-12.
