# Revisión U2a · Colorimetría 1.1 · ronda 2 (flujo P1)

Revisor: `revisor-U2a-r2`. Soy independiente: no escribí nada de lo revisado y no corregí nada. Fecha: 11-oct-2026.

Alcance: `docs/criterios/U2.md` (U2-02 a U2-09), decisiones K26 y K27 y hallazgo R3-04.

Qué leí:
- `git diff HEAD --stat` y el diff de `docs/design/tools/colorimetria.py` (solo las partes cambiadas).
- `docs/design/colorimetria.md`: «Cambios 1.1», §7.7, §12 y los renglones citados abajo.
- `docs/design/contraste.csv` y `docs/contratos/C2-colores.json`, ambos leídos por script.
- Los bloqueos `U2a-rosas-ximena-bodega` y `U2a-walmart-oscuro`.
- Los reportes `U2a-1.1.md` y `correccion-U2a-r1.md`.

Todas las cifras de este informe salen de mis propios comandos y scripts, no de los reportes.

## Veredicto
**APROBADO.** No hay hallazgos altos ni medios. Hay 5 bajos: 2 nuevos y 3 de la ronda 1 que siguen abiertos.

El hallazgo **H-01 (media) de la ronda 1 queda cerrado**:
- Existe el bloqueo `U2a-rosas-ximena-bodega`, con el formato del protocolo §2 y 4 opciones.
- Los 4 tokens que quedan bajo la regla frente a Ximena van PROVISIONAL en el JSON y en el md.
- «Cambios 1.1» punto 2 trae el ΔE 1.3, y §12 trae el hallazgo 14 y el bloqueo.

Repliqué por mi cuenta las cifras que sostienen el bloqueo y coinciden. Decidir si basta el refuerzo frente a Ximena le toca a Axel, en ese bloqueo.

## Verificaciones del encargo, (a) a (f)

**(a) `make verify-color`: pasa.**
- Salida: VERIFICACIÓN OK, con código de salida 0.
- md, csv, json y html son idénticos byte a byte a la regeneración en el árbol espejo.
- `contraste.csv` tiene 202 filas, 0 con cumple=false.

**(b) Cálculo propio: pasa.**
- Scripts: `verif.py` y `verif2.py` en `datos/tmp/revisor-u2a-r2`. OKLab de Ottosson (2020), OKLCH→sRGB con recorte por bisección y contraste WCAG 2.x, escritos aparte; no importan `colorimetria.py`.
- **Gasto:** el texto mide 44.73° en claro y 45.03° en oscuro (ver B-01).
- **Karen y Bodega Aurrera frente a los 8 roles, acción crítica y la marca:** las dos cumplen ≥ 20° o ΔL ≥ 0.15. Holgura mínima, con holgura = máx(Δmatiz/20, ΔL/0.15):

  | Color | Claro | Oscuro | Rol que la fija |
  |---|---|---|---|
  | Karen | 1.146 | 1.170 | deuda |
  | Bodega Aurrera | 1.168 | 1.197 | error |

- **Karen–Bodega Aurrera:** 22.95° con ΔL 0.233 en claro; 22.90° con ΔL 0.104 en oscuro. Cumple.
- **Gasto frente a error y advertencia:** gasto–error 25.26° y 24.93°; gasto–advertencia 30.69° y 30.72° (claro y oscuro).
- **Contraste de los 20 tokens cambiados:**
  - Texto de gasto, mínimo 4.5: 5.33, 5.72 y 5.33 en claro; 6.03, 7.71 y 5.90 en oscuro. En cada tema, el orden es: sobre su superficie, sobre el fondo, sobre la tarjeta.
  - Categóricas como UI, mínimo 3.0: el más bajo es Karen claro sobre tarjeta, 3.16.
  - Los 11 pasos de escala de gasto no aparecen en ningún par. Las únicas escalas de rol que el CSV usa son las de error, para el botón destructivo.
- **Única regla que no se cumple, frente a Ximena** (cubierta por el bloqueo):

  | Par | Claro | Oscuro |
  |---|---|---|
  | Karen–Ximena | 17.1°, ΔL 0.103, ΔE 10.93 | 16.7°, ΔL 0.111, ΔE 11.59 |
  | Bodega Aurrera–Ximena | 5.8°, ΔL 0.130, ΔE 13.10 | 6.2°, ΔL 0.006, ΔE 1.33 |

**(c) §7.7: pasa.**
- Incluye la marca: 700 en claro y 400 en oscuro.
- Mi barrido encuentra 6 pares con Δ matiz < 20° y ΔL < 0.15 en claro y 5 en oscuro. Son los mismos 11 renglones de §7.7, y todos declaran su refuerzo.
- Hay 3 pares marca–categórica, con refuerzo de forma y etiqueta. No hay ningún par marca–rol ni rol–categórica.
- Los 4 renglones frente a Ximena llevan la marca PROVISIONAL y citan el bloqueo.

**(d) ΔL consecutivo en la paleta de gráficas: pasa.**

| Medida | Claro | Oscuro |
|---|---|---|
| ΔL consecutivo mínimo, L OKLab | 0.174 | 0.218 |
| ΔL consecutivo mínimo, gris | 0.179 | 0.223 |
| Con simulación Machado 2009 (severidad 1) en RGB lineal | 0.162 (deuteranopia) | 0.198 (tritanopia) |
| Con la simulación aplicada en sRGB con gamma | 0.108 | 0.178 |
| 28 pares, gris mínimo (solo información) | 0.002 | 0.001 |
| Contraste mínimo contra fondo y contra tarjeta | 3.37 y 3.15 | 4.06 y 3.11 |

- Las cifras en RGB lineal y las de los 28 pares coinciden con «Cambios 1.1» punto 5.
- Las series no cambian respecto a HEAD.

**(e) Nombres de tokens: pasa.**
- Las 190 rutas son iguales y están en el mismo orden que en `git show HEAD`, y el `$type` no cambia en ninguna.
- `$value` cambia en 20 tokens, y `$extensions` en esos mismos 20.
- `$description` cambia en 5, todos dentro de esos 20: los 4 PROVISIONAL del bloqueo de las rosas y `color.store.walmart.dark`, que es PROVISIONAL por `U2a-walmart-oscuro`.

**(f) «Cambios 1.1» documenta los valores anteriores y los nuevos: pasa.**
- La tabla de parámetros va de 1.0 a 1.1, y la tabla de tokens trae hex, matiz, L y croma antes y después.
- Mi script lee la tabla de tokens: las 20 filas son justo los 20 `$value` que cambian.
- El hex 1.0 coincide con HEAD y el hex 1.1 con el JSON actual.
- Recalculé matiz, L y croma de las 40 celdas: 0 discrepancias.

## Criterios

| Criterio | Estado | Evidencia |
|---|---|---|
| U2-02 (alta) | Cumple. 4 tokens PROVISIONAL por bloqueo abierto | Matiz medido en el texto de cada rol, claro / oscuro: ingreso 135.13° / 134.72° (la propuesta de 135° está en §6.4), gasto 44.73° / 45.03°, ahorro 199.98° / 199.80°, deuda 310.23° / 309.87°, advertencia 75.42° / 75.75°, error 19.47° / 20.10°, éxito 150.00° / 150.03°, información 239.81° / 239.56°, acción crítica 100.05° / 100.07°, marca 700 / 400 280.25° / 279.95°. Todos quedan dentro de ±2°. Los extremos 50 / 950 coinciden con el plan: marca `#F9FAFF` / `#1E203A` y neutro `#FAFAFB` / `#222325`. Karen (333°) y Bodega Aurrera (356°) salen de 50° y cumplen frente a gasto, advertencia y los demás roles. Frente a Ximena, ver el bloqueo `U2a-rosas-ximena-bodega`. Hallazgos relacionados: B-01 y H-02. |
| U2-03 (alta) | Cumple en lo que cambia | En la escala de gasto, la L medida va de 0.9859 a 0.2557, con desviación máxima de 0.0013 frente a la L de cada paso. El croma máximo es 0.1400 (pasos 500 y 600), dentro del tope de 0.14. Gasto oscuro (paso 400) tiene croma 0.118, dentro del tope de 0.12. Karen tiene croma 0.120 y Bodega Aurrera 0.110 en claro y 0.106 en oscuro: el recorte de gamut lo baja, pero sigue sobre el piso de 0.10. Marca y neutro no cambian (punto e). |
| U2-04 (alta) | Cumple | Recalculé las 202 filas con mi implementación, no solo 20. Resultado: 0 diferencias mayores que 0.01, 0 hex desfasados del JSON y 0 pares bajo umbral. Ningún uso de texto tiene umbral 3.0. El CSV trunca a 2 decimales, así que 11 filas quedan 0.01 por debajo de mi redondeo. Las 16 filas que cambian frente a HEAD son justo las de los tokens cambiados. |
| U2-05 (alta) | Cumple. Sin cambios en 1.1 | Los pares estructurales están entre las 202 filas recalculadas. El par 900 / 50 da 11.53 frente a 11.4 del plan, y la diferencia está documentada: «difiere 0.13» (`colorimetria.md`, línea 633). |
| U2-06 (alta) | Cumple | Gasto como texto: 5.33, 5.72 y 5.33 en claro (sobre su superficie, el fondo y la tarjeta); 6.03, 7.71 y 5.90 en oscuro. El refuerzo es que el gasto va solo como cifra con signo «−» (`REFUERZO_ROLES`, K26). |
| U2-07 (alta) | Cumple. Sin cambios en 1.1 | Blanco sobre marca 700: 5.77. Marca 700 sobre fondo 50: 5.53. Texto 950 sobre marca 300: 10.47. Marca 300 sobre fondo 950: 10.37. La marca 600 solo se usa como anillo de foco (3.81 y 3.56, con umbral 3.0), y ningún texto va sobre ella. |
| U2-08 (alta) | Cumple. 2 pares por tema PROVISIONAL | El barrido incluye roles, acción crítica, las 8 categóricas y la marca. Da 11 pares, todos con refuerzo declarado. La colisión error–gasto de la 1.0 (5.5° en claro y 5.0° en oscuro) desaparece: ahora están a 25.26° y 24.93°. Karen y Bodega Aurrera ya reubicadas también pasan el barrido frente a los roles. Frente a Ximena, el refuerzo está declarado y la duda está en el bloqueo. |
| U2-09 (alta) | Cumple | Ver el punto (d). La simulación de daltonismo de las series (§8.4) no cambia, porque las series no cambian. |
| K26 | Cumple | Gasto en 45°: mide 44.73° en claro y 45.03° en oscuro. Error se queda en 20°. Karen y Bodega Aurrera dejan 50°. La marca entra al barrido. El contraste se verificó con script. |
| K27 | Cumple | El umbral duro es el ΔL entre series consecutivas, y los 28 pares van como información. El refuerzo es forma, patrón y etiqueta («Cambios 1.1», punto 5). |
| R3-04 | Atendido | §7.7 incluye la marca 700 / 400 y declara el refuerzo de los 3 pares marca–categórica. No hay pares marca–rol: si los hubiera, `verificar()` abortaría con «SIN RESOLVER». |

## Hallazgos

| ID | Severidad | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| B-01 | baja | U2-02 (medición del gasto), K26 | **Qué pasa:** el gasto se declara en 45°, el extremo del rango. Medidos a 2 decimales, 10 de los 15 tokens de gasto pasan de 45.00° (son 7 valores distintos): `.dark` y `.400` 45.03°; `.500` 45.09°; `.600` 45.10°; `.surface-dark` y `.900` 45.17°; `.950` 45.21°; `.surface-light` y `.100` 45.26°; y `.50` 48.66°, con croma 0.008. La tabla «Tokens que cambian» muestra esas cifras (de 45.0° a 45.3°, y 48.7°) sin explicarlas. Además, `verificar()` solo comprueba el matiz declarado (`h_gasto` de `ROLES`), no el que se mide en el hex. **Por qué es bajo:** todo queda dentro de la tolerancia de ±2° de U2-02, salvo el paso 50, cuyo matiz no significa nada con croma menor que 0.01. **Por qué cuenta:** la lectura literal «gasto medido entre 40° y 45°» depende de la precisión y de qué token se mida. | `docs/design/colorimetria.md` («Cambios 1.1»: tabla de tokens y punto 1); `docs/design/tools/colorimetria.py` (`GASTO_RANGO`, `verificar`) | En «Cambios 1.1» punto 1, que el script escriba qué tokens cuenta la comprobación de 40–45° (texto y superficie), con qué precisión y que aplica la tolerancia de ±2° de U2-02. También, que el paso 50 no cuenta por tener croma menor que 0.01. Y que `verificar()` compruebe además el matiz medido de `color.role.expense.light` y `.dark`. No cambia ningún valor. |
| B-02 | baja | Protocolo §2 (trazabilidad del bloqueo) | **(1)** El «Impacto» del bloqueo `U2a-rosas-ximena-bodega`, y también `correccion-U2a-r1.md`, dicen que la marca PROVISIONAL va en la «introducción de §7». Ahí no está: las líneas 344–347 no tienen ninguna coincidencia. Está en §7.5 (línea 413); las demás secciones citadas sí la llevan. **(2)** `docs/estado/bloqueos/INDICE.md` no registra ninguno de los dos bloqueos de U2a: no aparece ninguna vez «U2a-». Así, 5 tokens PROVISIONAL quedan fuera del índice que se lleva a la puerta. **(3)** `docs/estado/p1/U2a-1.1.md` sigue con la «Observación para Axel» en Supuestos y no lista el bloqueo nuevo. Los puntos (2) y (3) están fuera del alcance de U2a, que ya los pidió como solicitudes. | `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` (Impacto); `docs/estado/bloqueos/INDICE.md`; `docs/estado/p1/U2a-1.1.md` (Supuestos, Bloqueos) | En el bloqueo, cambiar «introducción de §7» por «§7.5». El integrador registra `U2a-walmart-oscuro` y `U2a-rosas-ximena-bodega` en el INDICE, y anota en `U2a-1.1.md` que la observación pasó al bloqueo. |
| H-02 | baja (sigue abierto) | U2-02 (documentación de desviaciones) | En §6.3, la nota de la fila Gasto (línea 305) sigue en «—». El matiz pasó de 25° (plan, `#B84D49` / `#ED8C84`) a 45° (K26), y la fila de Ingreso sí anota su cambio. | `docs/design/colorimetria.md` línea 305 (la genera `colorimetria.py`) | Que el script anote «25° en el plan; 45° por K26» en esa fila. |
| H-03 | baja (sigue abierto) | U2-08, R3-04 (caso vecino) | Lo recalculé. La marca 300 es el relleno del botón primario en oscuro. Frente a ella, Conjunto queda a 0.0°, ΔL 0.075 y ΔE 8.4, y Walmart a 17.4°, ΔL 0.149 y ΔE 16.0. Ninguno de los dos pares aparece en una tabla. | `docs/design/tools/colorimetria.py` (`PASO_MARCA`); `docs/design/colorimetria.md` §7.7 | Agregar la marca 300 en oscuro al barrido, solo como información, o llevar la nota a `DESIGN.md` en WF-1. |
| H-04 | baja (sigue abierto) | K26 (efecto colateral) | Lo recalculé. Gasto frente a la serie 3 en claro: 0.1°, ΔL 0.074 y ΔE 7.4. Frente a la serie 4 en oscuro: 10.2°, ΔL 0.055 y ΔE 6.0. Las series no entran en el barrido de U2-08. | `docs/design/colorimetria.md` §7.7 y §8 | En `DESIGN.md` (WF-1), declarar que ninguna serie representa un rol y que el gasto siempre va con signo «−». Otra opción: sumar las series al barrido como información. |

## Estado de los hallazgos anteriores

**H-01 (media; U2-02 y K26): CERRADO.** Lo comprobé así:
- **Bloqueo:** `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md` trae las 5 partes del protocolo §2 (Contexto, Qué falta, Opciones, Impacto, Cómo seguí) y 4 opciones, cada una con su consecuencia.
- **Cifras del bloqueo, replicadas con `verif2.py`:**
  - Con la L y el croma finales, las ventanas salen idénticas a las del md: Karen 331°–359°; Bodega Aurrera 331°–335°, 338°–341°, 344°–348°, 352° y 355°–359°. Los hex en el matiz entregado coinciden con el JSON.
  - Ningún matiz de esas ventanas las separa de Ximena en los dos temas.
  - Busqué también con matiz, L y croma libres: L de 0.25 a 0.75 en claro y de 0.55 a 1.00 en oscuro, croma 0.10, 0.11 o 0.12, y las otras tiendas fijas. Bodega Aurrera solo cabe en 331°–359°, y en ningún caso se separa de Ximena. Esto confirma la opción 4.
  - Opción 3: con L 0.615 en oscuro, Bodega Aurrera da `#B86787`, ΔL 0.221 y contraste 3.05 sobre tarjeta, como dice el bloqueo.
- **Marca PROVISIONAL:** está en la `$description` de los 4 tokens. En `colorimetria.md` aparece en las líneas 32–33, 99, 354, 362, 413, 429–430, 436–437, 683 y 695. Walmart en oscuro tiene la suya por su propio bloqueo (líneas 34, 101, 103, 361 y 682).
- **ΔE 1.3:** está en «Cambios 1.1» punto 2. **Hallazgo 14 y línea de Bloqueos:** están en §12.
- **Guarda del script:** en un espejo quité la marca de `color.store.bodega-aurrera.dark`. `colorimetria.py` abortó con código 1 y el mensaje «Bloqueo U2a-rosas-ximena-bodega: bajo la regla frente a Ximena …».
- **Pendiente fuera de U2a:** índice y reporte U2a-1.1 (ver B-02).

**H-02, H-03 y H-04 (bajos): siguen abiertos.** No estaban en el encargo de corrección. Los volví a medir y siguen igual (tabla de hallazgos).

**R3-04 (WF-0, revisión 3): atendido** en la colorimetría 1.1 (punto c). Cerrarlo en `docs/estado/wf0/REVISION-r3.md` le toca al integrador.

## Comandos ejecutados
- `make verify-color`: VERIFICACIÓN OK, 4 archivos idénticos byte a byte, 202 filas, 0 con cumple=false, código de salida 0.
- `python3 -I datos/tmp/revisor-u2a-r2/verif.py`: puntos (b) a (f), recálculo de las 202 filas del CSV y comparación del JSON con `git show HEAD`.
- `python3 -I datos/tmp/revisor-u2a-r2/verif2.py` (5.6 s): ventanas de matiz, búsqueda libre de Bodega Aurrera, cifra de la opción 3 y matiz medido por rol.
- Prueba de la guarda: copia de `colorimetria.py` en `datos/tmp/revisor-u2a-r2/espejo`, con `PROVISIONAL.pop("color.store.bodega-aurrera.dark")` añadido, ejecutada con `--sin-png`. Salió con código 1 y el mensaje esperado.
- `grep` y `sed -n` sobre `colorimetria.md` (encabezados, «Cambios 1.1», §7.7, §12, §6.3, §2.4, §2.7, línea 633), `contraste.csv`, `Makefile` e `INDICE.md`.
- Los temporales de `datos/tmp/revisor-u2a-r2` se borraron al terminar. `make verify-color` también usa `datos/tmp/verify-wf0`: es el directorio común del Makefile y lo dejé.
