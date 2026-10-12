# Revisión k45 · ronda 2 · 11-oct-2026, noche

Revisor: `revisor-k45` (independiente; no corrige nada). Objeto: el trabajo de `aplicar-k45` y su ronda de corrección `aplicar-k45-fix` (`docs/estado/puerta2/aplicar-k45-fix.md`). Referencias: `docs/estado/puerta2/k45-borrador.md` (fuente de la decisión), `docs/plan/protocolo-agentes.md` §4 (reparto fino) y la coherencia interna del repositorio.

**Veredicto: NO APROBADO.** Quedan 1 alto y 1 medio, además de 8 bajos. Los dos que bloquean se corrigen en pocos minutos:
- **Alto**: `README.md` línea 5 sigue diciendo que el paquete «**sube Axel** (K35)».
- **Medio**: P3-12 no cubre todas las cifras canon de `plan.md` §1.

Lo demás está bien hecho:
- K45 es idéntica al borrador y K35 conserva su texto.
- `make verify-criterios` y `make verify-all` terminan en OK.
- Ya no queda en `plan.md`, INDICE, `prompt.md` ni en los dos `LEEME.md` ninguna instrucción vigente de subir archivos o abrir claude.ai/design.
- El MANIFIESTO coincide.
- `git status` está dentro del alcance.

## 1. Comprobaciones (a) a (g), con mis propios comandos

| # | Qué | Resultado |
|---|---|---|
| a | K45 idéntica al borrador; K35 sustituida sin borrar texto | **PASA**. Comparé con `diff` el bloque K45 (5 líneas, de `- **K45` a `Consecuencias:`) extraído con `awk` de `k45-borrador.md` y de `decisiones.md`: sin diferencias. La línea K35 de HEAD contra la actual, quitando el prefijo `**Sustituida por K45 (11-oct-2026).**`, tampoco tiene diferencias. Además, §J punto 3 (línea 66) lleva «→ sustituido por K45 (11-oct-2026)». |
| b | `plan.md` §2 filas 41/45/63, §3.3, §3.6, §3.8, R4, §6 y §8.2 coherentes con K45; no contradicen 42/43/44 | **PASA**. Filas 41, 45 y 63 con nota de sustitución. §3.3 línea 206, §3.6 (8 pasos), §3.8 puerta 3, línea 334 y §8.2 línea 409 están en K45. §3.6 paso 8 dice «siguen vigentes las decisiones 42, 43 y 44»; el paso 5 tiene máximo 3 vueltas. R4 está marcado «Obsoleto desde K45». Restos menores: B-1, B-2 y B-8. |
| c | P3 conserva P3-01..P3-13 y severidades; medible; nada de subir ni de abrir claude.ai/design; P3-12 con cifras sintéticas y lista prohibida; `make verify-criterios` pasa | **PASA en parte**. Ids y severidades son iguales a HEAD (`diff` vacío): 10 alta y 3 media. Cada criterio tiene Verifica con comando o conteo. §3 «No aceptable» prohíbe pedir a Axel que suba archivos o comparta carpeta. `make verify-criterios` da exit 0, con `INDICE contra conteo propio: 0 diferencias` y `TOTAL 261 (alta 201 / media 60 / baja 0)`. **Falla en**: la lista de cifras de P3-12 está incompleta frente a `plan.md` §1 (M-1). |
| d | `LEEME.md` sin pasos de subida; MANIFIESTO coincide | **PASA**. LEEME: cabecera «NO se sube», «Qué pasa ahora» en 5 puntos y nadie sube nada. `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md \| sha256sum -c --quiet` da exit 0, con 22 elementos. Las descripciones del MANIFIESTO ya están en K45 (cabecera, filas 1, 2 y 11 a 21, y las 3 carpetas). |
| e | `pista3-plan.md` conserva los datos técnicos del borrador | **PASA**. Revisé el `diff` completo contra `pista3-plan-borrador.md`. Siguen las 2 `type_url`, los 2 releases, el formato de `tokens.json` y de `canvas.json`, las reglas del artboard (412×915, 1440, `support.js`, `DCLogic`), los límites de envío, las fases 0 a 2, la estimación de 7 a 9 M y la tabla P3. Faltan «ultracode» (B-02 de r1) y «cómo abrir y comentar» (B-6). |
| f | `make verify-all GOOSE=…/gobin/goose` en OK | **PASA**. `EXIT=0` y `Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`. U1: `comprobaciones OK: 63 · INFO: 5 · PENDIENTE: 0`. Paquete: `sha256sum -c: 21 archivos sin diferencias`. Log: `scratchpad/revk45r2/verify-all.log`. |
| g | `git status` solo con archivos del alcance | **PASA**. Hay 22 archivos modificados: los 14 de los hallazgos de r1 (README, INDICE, P3, U1, `entrada/LEEME.md`, MANIFIESTO, `prompt.md`, `salida/LEEME.md`, `paquete_entrada.sh`, `verificar_u1.py`, `decisiones.md` y `plan.md`) y, del generador, 4 PDF y 6 sellos de `.cache`. Sin seguimiento solo están `docs/design/pista3-plan.md` y `docs/estado/puerta2/`. Las secciones que cambian en `plan.md` y en INDICE son exactamente las que nombra el reporte. |

## 2. Estado de los hallazgos de la ronda 1

| Hallazgo r1 | Estado | Cómo lo comprobé |
|---|---|---|
| H-01 INDICE §10 en 1.2 rompe verify | **CERRADO** | Línea 149 dice `1.3`; `make verify-criterios` y `verify-all` dan exit 0. |
| H-02 `plan.md` 206, 334 y 409 con navegador | **CERRADO** | Las tres líneas están en K45 (sin navegador, URL del lienzo). La 334 conserva cifras viejas marcadas (B-2). |
| H-03 INDICE con K35 vigente (9, 28 y 44) | **CERRADO** | Las tres líneas están en K45. La puerta 3 queda como P3-01, P3-06 y P3-12, marcada PROPUESTA. |
| H-04 `prompt.md` contra K45.1 y K45.2 | **CERRADO** | «Usa solo las cifras» y «84 artboards»: 0 coincidencias. Hay sección «Datos: cifras sintéticas obligatorias (K45.1)», cobertura de todas las pantallas de U1 con 412 y 1440, y «No publiques» y «No abras navegador». 0 cifras canon ni `17.27`. |
| H-05 `salida/LEEME.md` en K35 | **CERRADO** | Reescrito con K45; sin pasos de pegar ni subir. Sigue la duda de quién guarda la exportación (B-3). |
| H-06 README con K35 | **ABIERTO EN PARTE → A-1** | Las líneas 57 y 156 están en K45, pero la línea 5 sigue con «**sube Axel** (K35)». |
| H-07 U1 y su verificador | **CERRADO** | U1 tiene la línea 27 con nota y §7 «Decidido (K45…)». `verificar_u1.py:295` exige `'K45'` y verify-u1 da 63 OK. Restos en comentarios: B-5. |
| M-01 descripciones del MANIFIESTO | **CERRADO** | La cabecera, las descripciones y las líneas de carpetas están en K45, y el `sha256sum -c` da OK. |
| M-02 decisiones §J.3 | **CERRADO** | La línea 66 lleva «→ sustituido por K45 (11-oct-2026)». |
| M-03 dónde va el informe de vuelta | **CERRADO** (PROPUESTA) | P3-06 exige la copia en `iter-N/` y la comprueba con `cmp`. `plan.md` §3.6 paso 5 y `pista3-plan.md` §5 dicen lo mismo. |
| M-04 cifras de Kueski en P3-12 | **ABIERTO EN PARTE → M-1** | Se añadieron 973.75, 1,258.70 y 883.77. Faltan otras cifras canon de `plan.md` §1 (ver M-1). |

## 3. Hallazgos de esta ronda

### Alto
**A-1 · instrucción vigente de subir archivos, contradice K45** (`/home/deepsec/finanzas-hogar/README.md`, línea 5).
- La cabecera «**Estado al 11-oct-2026, cierre de WF-1:**» dice: «… y el paquete para Claude Design en `docs/design/entrada/`, que **sube Axel** (K35). … **Siguiente: puerta 2.**».
- Es lo primero que lee quien abre el repositorio y está en presente. Contradice K45 («no hay subida manual de archivos») y las líneas 57 y 156 del mismo README, que ya dicen «no se sube (K45)».
- Es la misma clase de problema que H-06 de r1 (alto). Ni r1 ni `aplicar-k45-fix` lo encontraron: la búsqueda de huellas del fix buscaba «Axel sube», no «sube Axel».
- Corrección sugerida: «… y el brief de referencia de la Pista 3 en `docs/design/entrada/` (no se sube; K45 sustituye a K35: la Pista 3 la ejecuta el orquestador con Artifacts)».

### Medio
**M-1 · criterio no medible por completo: faltan cifras canon en P3-12** (`/home/deepsec/finanzas-hogar/docs/criterios/P3-claude-design.md`, línea 93, P3-12; también §4 y §5).
- La Verifica de P3-12 busca 11 cifras y las llama «las 8 cifras canon del reporte y las 3 de la regresión Kueski».
- Pero `plan.md` §1, línea 24 («Cifras canon para regresión (S3)»), lista más montos reales de la compra Kueski y del reporte que el grep no detecta:
  - **$2,395.28** (Axel) y **$1,499.72** (Karen);
  - **$2,616.87**, **$1,767.51** y **$849.36** (valor estable);
  - **$8,094.50** (lectura A del disponible);
  - **$14,362.48** y **$13,103.78** (totales P1–P6);
  - las tasas de Kueski **25.8685 %** y **17.27 %**.
- `17.27` aparece en `datos/semilla/data.json`, y `aplicar-k45-fix` la quitó de `prompt.md` justamente por ser dato real, pero no la puso en P3-12.
- K45.1 dice «nada del reporte real ni de Kueski sale de casa». Con la lista actual, un artboard con «$2,395.28» o «25.8685 %» pasa P3-12.
- Corrección sugerida (PROPUESTA):
  - sumar al patrón `2,?395\.28|1,?499\.72|2,?616\.87|1,?767\.51|849\.36|8,?094\.50|14,?362\.48|13,?103\.78|25\.8685|17\.27`;
  - o, mejor, cruzar con un comando, como ya se hace con las fechas, todos los montos con decimales de `oraculo/model.json` y `datos/semilla/data.json` contra los artboards, dejando fuera los redondos (`5,000.00`, `4,000`) como PREGUNTA;
  - poner al día §4 (umbral) y §5 (traza a `plan.md` §1 «Cifras canon para regresión (S3)»).

### Bajos (no bloquean)
- **B-1** · `plan.md` línea 43 (§1.5 HECHOS): «→ DESIGN.md se sube convertido a PDF…» no lleva nota de K45. La línea 45 (`file_upload`, R4) tampoco. Es una sección de hechos, pero se lee en presente. Sugerencia: «(flujo anterior a K45)».
- **B-2** · `plan.md` §6: la fila Pista 3 (línea 334) se queda con 2 agentes y 200K, «estimación anterior a K45», y el TOTAL (línea 343) sigue en ≈ 8.2 M. Pero `pista3-plan.md` §3 estima 7 a 9 M **solo** para la Pista 3, así que el total del proyecto no está conciliado. PREGUNTA al orquestador: ¿se recalcula el TOTAL o basta una nota?
- **B-3** (B-05 de r1, sigue abierto) · `salida/LEEME.md`: la tabla dice que `final/` lo escribe «el orquestador si la herramienta la da como archivo; si no, la descarga Axel». Su propia sección «Qué hace Axel» («Solo al final… abre la URL… y aprueba») y `entrada/LEEME.md`, que es la guía de Axel, no mencionan esa descarga. PREGUNTA: ¿quién guarda la exportación final?
- **B-4** · P3-12, fechas (PROVISIONAL, declarado): la confirmación manual de que «ninguna fecha coincidente va con un monto real» repite el grep de montos. Si los montos reales dan 0, la condición se cumple sola, así que en la práctica se aceptan fechas reales, cuando K45.1 dice «fechas inventadas». Ya es la PREGUNTA 2 de `aplicar-k45-fix`; aquí solo se deja constancia.
- **B-5** · Restos menores de K35:
  - `INDICE.md` línea 8 nombra `P3-iter-N.md`; en P3 1.3 es `P3-iter-N-<area>.md`.
  - La cabecera de INDICE sigue en «Versión 1.2» aunque cambiaron §1, §3, §10 y §11.
  - `verificar_u1.py` líneas 20 y 218 (docstring y comentario) dicen «K35 (Pista 3 en §7)» y la comprobación ya es K45.
  - En U1 §7, la frase «El orquestador no abre la herramienta» queda pegada a la «Historia (K35)» y puede leerse como vigente.
  - U1 sigue en 1.2 sin revisión formal (pendiente 4 de `aplicar-k45-fix`).
- **B-6** · El borrador (§4) dice que P3-03 pide que `LEEME.md` explique a Axel «cómo abrir y comentar». P3 1.3, `pista3-plan.md` §4 y `LEEME.md` dicen «cómo avisa», y la respuesta es en la sesión de Claude Code. La vía de los comentarios en el lienzo desaparece de la guía, aunque P3-04 los menciona como dato. PREGUNTA al orquestador: ¿Axel comenta en el lienzo, contesta en la sesión, o las dos?
- **B-7** · `README.md` línea 85 («Versiones vigentes»): «Paquete para Claude Design … **no aprobado** (ronda 3: 1 medio, una frase del LEEME)» describe el LEEME anterior a K45. Sugerencia: «Brief de referencia de la Pista 3 (K45), reescrito el 11-oct-2026».
- **B-8** · Siguen abiertos los bajos B-01 a B-04 de r1:
  - P3-06 no copia `project/ds/`, que la Verifica de P3-02 lee;
  - «ultracode» omitido;
  - 512 contra 511 archivos por versión (la descripción de la herramienta Artifact que ve este revisor dice «up to 511 files»);
  - en R4, la mitigación vieja va antes de la nota «Obsoleto».

## 4. Para aprobar en la ronda 3
1. A-1: una frase en `README.md` línea 5.
2. M-1: completar la lista o el método de P3-12, con §4 y §5 al día. Después, volver a correr `make verify-criterios` (la traza nueva de §5 debe tener fuente) y `make verify-all`.
Los bajos y las PREGUNTAS (B-2, B-3, B-4 y B-6) pueden pasar a la puerta 2 sin bloquear.
