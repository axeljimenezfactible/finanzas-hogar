# Revisión U2a · colorimetría 1.2 · ronda 1 de 3 (flujo P1c, 11-oct-2026)

Revisor independiente `revisor-U2a-r1`. No escribí nada de lo revisado y no corregí nada. Evalúo solo contra `docs/criterios/U2.md` 1.2 (U2-02 a U2-09) y `docs/plan/decisiones.md` §K (K26, K27, K41, K42). Todas las cifras de abajo salen de mis propios comandos y scripts en `datos/tmp/rev-u2a-r1/` (ya borrados), no del reporte `docs/estado/p1c/U2a-1.2.md`.

## Veredicto

**NO APROBADO**: 0 hallazgos altos, **1 medio** (M-01) y 2 bajos. Las seis verificaciones del encargo pasan. Ningún valor cambió, los nombres de los tokens son los mismos y no queda ninguna marca PROVISIONAL. Lo que falla es que §12 de `colorimetria.md` todavía pone los matices nuevos de Karen y Bodega Aurrera entre las «Propuestas sujetas a la puerta 2», aunque K41 ya los aceptó y el mismo documento los llama FIRMES. Para corregirlo basta cambiar un texto en el script y regenerar.

## Verificaciones del encargo

| | Verificación | Resultado | Evidencia (comando propio) |
|---|---|---|---|
| a | `make verify-color` | Pasa | Salida 0 (`set -o pipefail`). Da «VERIFICACIÓN OK», 202 filas con 0 en cumple=false, y `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html` idénticos byte a byte. Una segunda regeneración en mi propio árbol espejo (`--sin-png`) también salió idéntica en los 4. `contraste.csv`, `paleta.html` y `paleta.png` no difieren de HEAD (`git diff --quiet`). |
| b | Ningún hex cambió contra `git show HEAD` | Pasa | Script propio que aplana los dos JSON. 190 tokens en HEAD y 190 ahora. `$value`, `$type` y `$extensions`: 0 diferencias. `$description` cambia solo en los 5 tokens de K41 y K42, las descripciones de grupo no cambian y la raíz `$` tampoco. |
| c | PROVISIONAL en los 5 archivos | Pasa | `grep -c PROVISIONAL` da 0 en `colorimetria.md`, `contraste.csv`, `paleta.html`, `paleta.png` y `C2-colores.json`. En `colorimetria.py` hay 19 coincidencias, todas del mecanismo (`PROVISIONAL = {}` y sus comprobaciones), ninguna con la cita de un bloqueo. En todo el repositorio, fuera de `docs/estado/` (histórico), ningún artefacto vigente marca PROVISIONAL citando `U2a-walmart-oscuro` o `U2a-rosas-ximena-bodega`. |
| d | Regla K41 para `DESIGN.md` en «Cambios 1.2» | Pasa (con B-02) | `colorimetria.md`, líneas 43–48: «Una persona y una tienda nunca se distinguen solo por color en una misma vista» y «U2-02: entre categóricas basta el refuerzo de U2-08, también entre grupos distintos». |
| e | Cálculo propio en OKLab de los 5 tokens | Pasa | Uso las matrices de Ottosson desde el hex. Karen: `#B870AC` da 333.2°, L 0.6429 y C 0.1199; `#D389C6` da 333.3°, L 0.7251 y C 0.1196. Bodega Aurrera: `#762C4C` da 356.1°, L 0.4096 y C 0.1095; `#FFAAC9` da 356.2°, L 0.8292 y C 0.1063. Walmart en oscuro: `#7AA3EE` da 261.9°, L 0.7165 y C 0.1190. Con la conversión inversa propia, oklch(0.716, 0.12, 262) da `#7AA3EE` y oklch(0.720, 0.12, 262) da `#7BA4F0`, con ΔE 0.39. Frente a Ximena: Karen a 17.1° / ΔL 0.103 / ΔE 10.9 en claro y 16.7° / 0.111 / 11.6 en oscuro; Bodega Aurrera a 5.8° / 0.130 / 13.1 en claro y 6.2° / 0.006 / 1.3 en oscuro. Las tiendas en oscuro quedan con ΔL mínimo 0.1008 en OKLab y 0.1008 en gris. Todo coincide con «Cambios 1.2», §7.7 y §12. |
| f | Nombres de tokens sin cambios | Pasa | Los conjuntos de nombres de HEAD y del árbol de trabajo son iguales (190). |

Pruebas negativas propias, cada una en un árbol espejo con el script parcheado y una sola sustitución:

| Cambio | Salida | Lo que dice |
|---|---|---|
| Hex de Karen en claro cambiado en `FIRMES_12` | 1 | «Cambios 1.2: cambió el hex…» |
| `PROVISIONAL` con Walmart en oscuro | 1 | «quedan marcas PROVISIONAL» |
| L de Walmart en oscuro de 0.716 a 0.720 en `TIENDAS` | 1 | Se detiene |
| Bodega Aurrera en oscuro reasignada a K42 en `FIRMES_12` | 1 | Se detiene en `rosas_ximena()` |

## Criterios

| Criterio | Estado | Evidencia |
|---|---|---|
| U2-02 (alta) | Cumple. M-01 es la incoherencia de estado | Matices medidos en claro / oscuro: ingreso 135.1° / 134.7° (la propuesta de 135° está en §6.4), gasto 44.7° / 45.0° (K26), ahorro 200.0° / 199.8°, deuda 310.2° / 309.9°, advertencia 75.4° / 75.7°, error 19.5° / 20.1°, éxito 150.0° / 150.0° e información 239.8° / 239.6°. La marca 700 mide 280.2° y el neutro 500, C 0.0116 en 280.4°. Karen queda en 333.2° / 333.3° y Bodega Aurrera en 356.1° / 356.2°, sin marca y con cita a K41. |
| U2-03 (alta) | No se reevalúa: sin cambios | Hex y `$extensions.oklch` idénticos a HEAD en los 190 tokens. Sigue el estado de la revisión anterior (`docs/estado/p1/REVISION-U2a-r2.md`). |
| U2-04 (alta) | Cumple | 202 filas, 0 en cumple=false y 0 con razón bajo el umbral. Recalculé con una implementación propia 32 pares: los 12 de Karen, Bodega Aurrera y Walmart, más 20 al azar con semilla 20261011. Los 32 quedan dentro de ±0.01: el CSV trunca a 2 decimales, como declara §10 en la línea 645. Los hex del CSV son iguales a los del JSON. |
| U2-05 (alta) | Cumple (sin cambios) | Recalculado: 900/50 da 11.53 (11.4 en el plan; la diferencia está documentada en la línea 660), 700/50 da 5.40, 600/100 da 3.47, blanco/marca 600 da 3.98, 100/950 da 14.07, 300/900 da 7.98, marca 300/950 da 10.38, 600/50 da 3.71 y 400/950 da 7.96. |
| U2-06 (alta) | Cumple (sin cambios) | Los roles como texto dan como mínimo 4.74 en claro (sobre 50 y 100) y 5.82 en oscuro (sobre 950 y 900). |
| U2-07 (alta) | Cumple (sin cambios) | Marca 700/50 da 5.53, blanco/marca 700 da 5.77 y 950/marca 300 da 10.38. Que el par blanco/marca 600 no aparezca en ningún texto se comprueba en `DESIGN.md` (WF-1). |
| U2-08 (alta) | Cumple. B-02 es de redacción | Mi barrido incluye los 8 roles, acción crítica, las 8 categóricas y la marca (700 en claro, 400 en oscuro). Da 6 pares en claro y 5 en oscuro bajo «≥ 20° o ΔL ≥ 0.15», los mismos 11 de la tabla de §7.7 (líneas 455–465). Ningún par queda entre una categórica o la marca y un rol o acción crítica. Karen y Bodega Aurrera aparecen con su refuerzo y citan K41. `color.store.walmart.dark` = `#7AA3EE` en L 0.716, FIRME por K42 (§7.2, «Cambios 1.1» punto 3, «Cambios 1.2» y §12). |
| U2-09 (alta) | Cumple (sin cambios) | Las 8 series no cambian. ΔL mínimo entre series consecutivas en OKLab: 0.174 en claro y 0.218 en oscuro; en gris, 0.179 y 0.223. Contraste mínimo sobre fondo y tarjeta: 3.15 en claro y 3.11 en oscuro. La simulación de daltonismo de §8.4 no cambia (no la recalculé con culori en esta ronda). |
| K26 | Cumple | El gasto mide 44.7° / 45.0° y el error 19.5° / 20.1°. Bloqueo `U2-gasto-error` resuelto. |
| K27 | Cumple | Ver U2-09. Los 28 pares se reportan como información (§8.5). |
| K41 | Cumple, salvo M-01 y B-02 | Se quitaron las 4 marcas y ningún token cambia. La regla para `DESIGN.md` está en «Cambios 1.2». |
| K42 | Cumple | Se quitó la marca. L medida 0.7165 (nominal 0.716). Escalonamiento: 0.1008 ≥ 0.10. |

## Hallazgos

| Id | Severidad | Criterio | Qué pasa | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| M-01 | media | K41, U2-02 (incoherencia entre artefactos) | §12, línea 720, «Propuestas sujetas a la puerta 2», todavía incluye «en la versión 1.1, gasto 45°, **Karen y Bodega Aurrera en sus matices nuevos**». K41 ya aceptó Karen en 333° y Bodega Aurrera en 356° («se aceptan… ningún token cambia»). El mismo documento los declara «FIRME por K41» en §7.1, §7.2, «Cambios 1.2» y la línea de bloqueos de §12. En la puerta 2 se le volvería a pedir a Axel aprobar algo que ya decidió. El texto no cambió desde la 1.1. | `docs/design/tools/colorimetria.py`, línea 1709 (genera `docs/design/colorimetria.md` §12, línea 720) | En el script, sacar a Karen y Bodega Aurrera de esa lista, o escribir «Karen 333° y Bodega Aurrera 356°, aceptadas por K41 en la puerta 1b», y regenerar. Si se quiere un candado, que la verificación final falle cuando un token de `FIRMES_12` aparezca en ese párrafo como propuesta. |
| B-01 | baja | K42, U2-08 (claridad) | La tabla «Marcas que se quitan» de «Cambios 1.2» (líneas 29–35) da matiz, L y croma medidos en el hex: Walmart en oscuro L 0.717, Karen en claro L 0.643 y Bodega Aurrera en oscuro C 0.106. K42, el punto K42 de la misma sección y §7.1–§7.2 dan el nominal (L 0.716, 0.642, croma 0.107). El encabezado no dice «medidos en el hex», así que «L 0.717» parece contradecir K42 «L 0.716». Los dos valores son correctos. | `docs/design/tools/colorimetria.py`, `md_cambios_12()` (encabezado de la tabla) | Encabezado «Valor sin cambio: hex · matiz · L · croma (medidos en el hex)», como en «Tokens que cambian (medidos en el hex)» de «Cambios 1.1», o agregar el nominal. |
| B-02 | baja | U2-08, K41 (redacción del refuerzo) | El punto 2 de «Regla nueva para `DESIGN.md`» dice que §7.7 lista cada par con su refuerzo «(orden fijo, forma y etiqueta directa)». Las filas no lo dicen así. Ximena–Karen (líneas 456 y 463) dice «orden fijo, etiqueta directa y ΔL ≥ 0.10», sin forma, y le agrega «también entre grupos» aunque es un par del mismo grupo. Las filas persona–tienda (457 y 464) no nombran el orden fijo ni la forma. El punto 1 pide «nombre en texto y su forma», pero no el orden fijo que U2-08 exige cuando una vista mezcla personas y tiendas. Entre las dos reglas sí están los tres elementos y cada par declara un refuerzo, así que U2-08 se cumple. | `docs/design/tools/colorimetria.py` (cadenas `ref` del barrido de §7.7 y `md_cambios_12()`, punto 1) | Unificar el texto: en las 4 filas de K41, «orden fijo, forma y etiqueta directa», y en el punto 1, «orden fijo, forma (punto, anillo de avatar o marca de serie) y nombre en texto». Quitar «también entre grupos» de las filas del mismo grupo. Regenerar. |

## Estado de los hallazgos anteriores (`docs/estado/p1/REVISION-U2a-r2.md`)

| Id | Severidad | Estado ahora |
|---|---|---|
| B-01 (el gasto mide más de 45.00° en 10 tokens) | baja | Sigue. Los valores no cambiaron (`.light` 44.7° y `.dark` 45.0°, dentro de ±2°). |
| B-02 (1) (la marca PROVISIONAL en la «introducción de §7») | baja | Superado: las marcas se quitaron y «Cambios 1.2» cita bien sus lugares de la 1.1 («Cambios 1.1», §7.1, §7.2, §7.5, §7.7 y §12). |
| B-02 (2) (`INDICE.md` no registraba los bloqueos) | baja | Cerrado: `docs/estado/bloqueos/INDICE.md`, líneas 320–321, los registra con K41 y K42. |
| H-02 (§6.3, nota de la fila Gasto en «—») | baja | Sigue abierto (línea 332). |
| H-03 (marca 300 en oscuro fuera del barrido) | baja | Sigue abierto. |
| H-04 (gasto frente a las series 3 y 4) | baja | Sigue abierto. Se atiende en `DESIGN.md` (WF-1). |
