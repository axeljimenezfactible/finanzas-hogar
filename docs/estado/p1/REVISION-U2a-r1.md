# Revisión U2a · Colorimetría 1.1 · ronda 1 (flujo P1)

Fecha: 11-oct-2026 · Revisor: `revisor-U2a-r1` (independiente: no escribió nada de lo revisado y no corrige nada) · Se evalúa solo contra `docs/criterios/U2.md` (U2-02 a U2-09), decisiones §K (K26, K27) y el hallazgo R3-04 · Artefactos: `docs/design/tools/colorimetria.py` (diff contra HEAD), `docs/design/colorimetria.md` («Cambios 1.1», §6, §7, §8, §12), `docs/design/contraste.csv`, `docs/contratos/C2-colores.json`, `docs/estado/bloqueos/U2a-walmart-oscuro.md` y el reporte `docs/estado/p1/U2a-1.1.md`.

## Veredicto

**NO APROBADO**: 0 hallazgos altos, 1 medio y 3 bajos.

Todas las cifras se reproducen con mi propia implementación, escrita aparte y sin importar `colorimetria.py`: OKLab de Ottosson 2020, contraste WCAG 2.x, gris = Y^(1/3) como dice §2.7, simulación de Machado 2009 y una réplica del barrido de matices. `make verify-color` pasa.

Lo único que impide aprobar es H-01. La reubicación deja a Bodega Aurrera prácticamente igual a Ximena en oscuro (ΔE 1.3). Se resolvió con un refuerzo declarado, cosa que U2-08 admite. Pero el texto de U2-02 y la respuesta de Axel piden mover estos colores «para conservar ≥ 20° o ΔL ≥ 0.15». Esa diferencia no se llevó a un bloqueo: solo quedó como «Observación para Axel» en los Supuestos del reporte.

## Verificaciones del encargo, (a) a (f)

| | Verificación | Resultado | Evidencia (cálculo propio) |
|---|---|---|---|
| a | `make verify-color` | Pasa | «VERIFICACIÓN OK». Los 4 archivos son idénticos byte a byte: md, csv, json y html. `contraste.csv` tiene 202 filas y 0 con cumple=false. |
| b | Gasto en 40–45°, Karen y Bodega Aurrera frente a cada rol y entre sí, y contraste | Pasa | **Gasto:** matiz medido 44.73° en claro y 45.03° en oscuro, dentro de la tolerancia de ±2° de U2-02. Queda a 25.26° y 24.93° de error y a 30.69° y 30.72° de advertencia. **Karen y Bodega Aurrera:** 0 fallos frente a los 8 roles, acción crítica y la marca. Los pares más cercanos son Karen–deuda (22.93° y 23.40°) y Bodega Aurrera–error (23.37° y 23.93°). **Karen–Bodega Aurrera:** 22.95° con ΔL 0.233 en claro y 22.90° con ΔL 0.104 en oscuro. **Texto de gasto (≥ 4.5):** 5.71, 5.33 y 5.32 en claro; 7.70, 5.89 y 6.02 en oscuro, sobre fondo, tarjeta e insignia. **Elementos de interfaz (≥ 3.0), sobre fondo y tarjeta:** Karen 3.38 y 3.16 en claro, 6.11 y 4.68 en oscuro. Bodega Aurrera 9.01 y 8.40 en claro, 8.87 y 6.79 en oscuro. Walmart en oscuro 6.22 y 4.76. |
| c | §7.7 incluye la marca y declara refuerzos | Pasa | Quedan cerca 3 pares con la marca: Conjunto–marca 700 (0.9°, ΔL 0.094), Walmart–marca 700 (18.0°, ΔL 0.000) y Walmart–marca 400 (18.0°, ΔL 0.069). Los tres llevan refuerzo de forma y etiqueta. Ningún rol queda cerca de la marca. |
| d | ΔL ≥ 0.10 entre series consecutivas | Pasa | **Claro:** OKLab 0.174, gris 0.179, protanopia 0.174, deuteranopia 0.162, tritanopia 0.177. **Oscuro:** 0.218, 0.223, 0.224, 0.205 y 0.198. Los mínimos (0.162 y 0.198) coinciden con §8.3 al milésimo. |
| e | Nombres de token de C2 sin cambios frente a `git show HEAD` | Pasa | 190 tokens en ambas versiones, 0 agregados y 0 quitados. 20 cambian de `$value`. `$type` no cambia en ninguno. `$extensions` y `$description` solo cambian en esos 20 (la descripción, en 3). |
| f | «Cambios 1.1» documenta valores anteriores y nuevos | Pasa | La sección trae una tabla de parámetros (4 filas) y otra de los 20 tokens con hex, matiz, L y croma de 1.0 → 1.1. Los hex de 1.0 coinciden con HEAD y los de 1.1 con el JSON. Matiz, L y croma coinciden con mi medición (al 0.1° y al 0.001). |

**Verificación adicional: réplica del barrido de ventanas.** Convertí de OKLCH a sRGB reduciendo croma, y la conversión reproduce los 7 hex de referencia. Resultados:
- Karen cabe en 331°–359°.
- Bodega Aurrera cabe en 331°–335°, 338°–341°, 344°–348°, 352° y 355°–359°.
- Con Walmart en oscuro en L 0.720 (`#7BA4F0`), Bodega Aurrera no cabe en ningún matiz: en 356° queda a ΔL gris 0.097 de Walmart.

El cambio PROVISIONAL de Walmart se sostiene por cálculo. Además tiene su bloqueo con el formato del protocolo §2 y la marca en el JSON y en el md.

## Criterios

| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| U2-02 Matiz aplicado | alta | Cumple lo que fija K26, con observación (H-01) | Gasto en 45°; error sigue en 20° (su token no cambia); Karen en 333° y Bodega Aurrera en 356°, ambas a ≥ 20° de todos los roles. El ingreso en 135° es una propuesta documentada (§6.4). Marca y neutro no cambian (22/22). Frente a «las demás entradas», Bodega Aurrera–Ximena y Karen–Ximena quedan bajo la regla (H-01). |
| U2-03 Método y saturación | alta | Cumple | En la escala de gasto, \|ΔL\| máximo de 0.0013 (paso 700) frente a la L fija por paso y croma máximo de 0.1400. El texto en oscuro (paso 400) tiene croma 0.118, bajo 0.12. Las constantes del método no aparecen en el diff. |
| U2-04 Contraste: 0 pares bajo umbral | alta | Cumple | Recalculé las 202 filas con mi implementación: diferencia máxima ≤ 0.01 (el CSV trunca a 2 decimales), 0 con cumple=false, y el hex del CSV coincide con el JSON en todas. Frente a HEAD cambian 16 filas, todas de tokens cambiados, y no hay filas nuevas ni quitadas. Mínimos: texto 4.73, interfaz 3.05. Los únicos tokens cambiados sin fila son los pasos 50–950 de la escala, que no se usan como par. |
| U2-05 Pares estructurales | alta | Cumple (sin cambios en 1.1) | 900/50 da 11.53 (§10 declara «difiere 0.13»); 700/50, 5.40; 600/100, 3.47; blanco/marca 600, 3.98; 100/950, 14.07; 300/900, 7.97; marca 300/950, 10.38; 200/50, 1.20; 800/950, 1.90; 600/50, 3.71; 400/950, 7.96. |
| U2-06 Rol como texto | alta | Cumple | El texto de gasto usa el paso 700 en claro (5.71, 5.33 y 5.32) y el 400 en oscuro (7.70, 5.89 y 6.02). El paso 600 (3.65 a 3.92) queda solo para relleno o icono (§6.2). |
| U2-07 Botón primario | alta | Cumple (sin cambios) | Marca 700/50: 5.53. Blanco/marca 700: 5.77. 950/marca 300: 10.38. La marca 600 solo aparece en el CSV como anillo de foco (interfaz), nunca como texto. |
| U2-08 Colisiones | alta | Cumple, con observaciones (H-01 y H-03) | Mi barrido encuentra 6 pares cercanos en claro y 5 en oscuro, los mismos de §7.7, todos con refuerzo declarado. Ningún par es categórica o marca frente a rol. Error–gasto pasa de 5° a 25.3° en claro y 24.9° en oscuro. Las colisiones de decisiones §I se resuelven por escrito en §7.5. Sam's queda a 25° de ahorro (§7.2) e información a 40° de la marca. |
| U2-09 Paleta de gráficas | alta | Cumple | Ver (d). Los 28 pares se reportan solo como información (gris mínimo 0.002 en claro y 0.001 en oscuro, §8.5). Contraste mínimo: 3.15 en claro (el documento dice 3.14 por el truncado) y 3.11 en oscuro. Refuerzo con forma, patrón y etiqueta («Cambios 1.1» §5 y §8.5). Las series no cambian frente a HEAD. Mi simulación de Machado 2009 reproduce los ΔL de §8.3. |
| K26 | — | Cumple, salvo H-01 | Gasto en 40–45°, error en 20°, Karen y Bodega Aurrera reubicadas, la marca en el barrido y todo regenerado por el script con contraste verificado. |
| K27 | — | Cumple | El umbral duro es entre series consecutivas, los 28 pares van como información y el refuerzo está declarado. |
| R3-04 | — | Atendido | La marca 700 (claro) y 400 (oscuro) entra en §7.7 con su refuerzo, y el script aborta si queda un rol cerca de la marca. |

## Hallazgos

| ID | Sev. | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| H-01 | media | U2-02, K26 (respuesta de Axel en `U2-gasto-error`) | **Qué pasa:** la reubicación deja dos rosas fuera de la regla «≥ 20° o ΔL ≥ 0.15 frente a… las demás entradas». Bodega Aurrera–Ximena: en oscuro 6.2° y ΔL 0.006, ΔE 1.33 (`#FFAAC9` frente a `#FFACD3`, a simple vista iguales); en claro 5.8° y ΔL 0.130. Karen–Ximena: 17.1° y 16.7°, ΔL 0.103 y 0.111. **Por qué no es alto:** U2-08 admite el refuerzo declarado y §7.7 lo declara. **Por qué cuenta:** es una colisión nueva de la 1.1 (en la 1.0 el par cercano entre grupos más parecido tenía ΔE 10.5). Con las L actuales es inevitable: mi réplica confirma que ningún matiz de las ventanas separa a Bodega Aurrera de Ximena. U2-02 admite dos lecturas en este punto («las demás entradas» frente al refuerzo de U2-08), y el protocolo §2 pide en ese caso un bloqueo. No lo hay: solo está la «Observación para Axel» en los Supuestos del reporte, sin cifra de ΔE y sin hallazgo en §12. | `docs/design/colorimetria.md` («Cambios 1.1» §2, §7.7, §12); `docs/design/tools/colorimetria.py` (PERSONAS, TIENDAS); `docs/estado/p1/U2a-1.1.md` (Supuestos) | Escribir el bloqueo `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` con opciones. Por ejemplo: (1) aceptar con el refuerzo declarado; (2) mover a Ximena de matiz o de L; (3) reordenar la L de las tiendas en oscuro; (4) otra. Mientras se decide, marcar PROVISIONAL los tokens que dependan de la opción (al menos `color.store.bodega-aurrera.dark`), citando ese bloqueo. Agregar la cifra de ΔE 1.3 en «Cambios 1.1» §2 y un hallazgo en §12. Otra vía: pedir que `docs/criterios/U2.md` (U2-02) aclare si «las demás entradas» admite el refuerzo de U2-08 entre categóricas de grupos distintos. |
| H-02 | baja | U2-02 (documentación de desviaciones) | En la tabla de §6.3, comparación con el plan §4.3, la fila de Gasto lleva «—» en la nota. El matiz cambió de 25° (plan) a 45° (K26), mientras que la fila de Ingreso sí anota «150° en el plan». | `docs/design/colorimetria.md` línea 303 (la genera `colorimetria.py`) | Hacer que el script anote «25° en el plan; 45° por K26» en esa fila. |
| H-03 | baja | U2-08, R3-04 (adyacente) | El barrido usa la marca 400 en oscuro, como pide el criterio. Pero el relleno del botón primario en oscuro es la marca 300 (U2-07), y frente a ella Conjunto queda a 0.0° con ΔL 0.075 (ΔE 8.4) y Walmart a 17.4° con ΔL 0.149. El texto del refuerzo de marca ya menciona el botón primario, pero esos pares no aparecen en ninguna tabla. | `docs/design/tools/colorimetria.py` (`PASO_MARCA`); `docs/design/colorimetria.md` §7.7 | Agregar la marca 300 en oscuro al barrido, solo como información, o llevar la nota a `DESIGN.md` en WF-1. |
| H-04 | baja | K26 (efecto colateral fuera del barrido de U2-08) | Con gasto en 45°, el rol de gasto coincide en matiz con la serie 3 de gráficas, «óxido» (45°). En claro quedan a 0.1°, ΔL 0.074 y ΔE 7.4; en la 1.0 había 19.9° entre ellos. En oscuro, gasto frente a la serie 4, «naranja», queda a 10.2°, ΔL 0.055 y ΔE 6.0. Las series no entran en el barrido de U2-08. | `docs/design/colorimetria.md` §7.7 y §8 | En `DESIGN.md` (WF-1), declarar que ninguna serie representa un rol y que el gasto siempre va con signo «−». Otra opción: sumar las series al barrido como información. |

## Estado de los hallazgos anteriores

- **R3-04** (WF-0, revisión 3, baja): **atendido** en la colorimetría 1.1. §7.7 incluye la marca y declara el refuerzo, y el script lo verifica. Cerrarlo en `REVISION-r3.md` le toca al integrador, como ya pidió U2a.
- **Bloqueo `U2-gasto-error`:** resuelto por K26 y aplicado (gasto en 45°, a 25.3° y 24.9° de error).
- **Bloqueo `criterios-deltal`:** resuelto por K27 y aplicado («Cambios 1.1» §5 y §8.5).
- **Bloqueo nuevo `U2a-walmart-oscuro`** (PROVISIONAL): la forma es correcta y el cálculo que lo justifica está verificado. Queda pendiente de la decisión de Axel y de registrarse en `docs/estado/bloqueos/INDICE.md` (lo hace el orquestador).
- Es la primera ronda de U2a en P1: no hay hallazgos de rondas anteriores de este flujo.

## Comandos ejecutados

- `make verify-color`: pasa.
- `python3 -I datos/tmp/rev-u2a-r1/rev.py`: compara tokens contra HEAD, mide los tokens cambiados, revisa separaciones y pares cercanos, escalonamiento y gráficas.
- `python3 -I datos/tmp/rev-u2a-r1/rev2.py`: recalcula las 202 filas del CSV, compara contra HEAD, revisa la escala de gasto y los pares estructurales, y replica el barrido de ventanas, incluida la prueba de Walmart.
- Los scripts se borraron al terminar.
