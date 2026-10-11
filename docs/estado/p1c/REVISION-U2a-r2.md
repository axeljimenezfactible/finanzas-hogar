# Revisión U2a · ronda 2 de 3 (flujo P1c, puerta 1b, 11-oct-2026)

Revisor independiente `revisor-U2a-r2`. No escribí nada de lo revisado y no corregí nada. Evalúo colorimetría 1.2 (`docs/design/tools/colorimetria.py`, `docs/design/colorimetria.md`, `docs/contratos/C2-colores.json`, `docs/design/contraste.csv`, `docs/design/paleta.html`, `docs/design/paleta.png`) solo contra `docs/criterios/U2.md` (U2-02 a U2-09) y `docs/plan/decisiones.md` §K (K26, K27, K41, K42). Todas las cifras de abajo salen de mis propios comandos y scripts en `datos/tmp/rev-u2a-r2/` (ya borrados), no de los reportes.

## Veredicto

**APROBADO.** No hay hallazgos de severidad alta ni media. El hallazgo M-01 de la ronda 1 está cerrado. Siguen abiertos dos bajos de la ronda 1 (B-01 y B-02) y agrego uno bajo nuevo (B-03). Ninguno cambia valores, tokens ni cumplimiento de criterios.

## Verificaciones del encargo

| # | Verificación | Comando o script propio | Resultado |
|---|---|---|---|
| a | `make verify-color` | `make verify-color` (salida 0) | Regenera en árbol espejo: idénticos byte a byte `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html`; «contraste.csv: 202 filas, 0 con cumple=false»; «VERIFICACIÓN OK». Además, `python3 -I docs/design/tools/colorimetria.py --solo-verificar`: salida 0, «VERIFICACIÓN OK», «versión 1.2: 5 tokens FIRMES por K41 y K42, 0 con otro hex; marcas PROVISIONAL en el JSON: 0», «§12: personas o tiendas FIRMES… entre las propuestas de la puerta 2: 0 de 3». |
| b | Ningún hex cambió contra HEAD | `cmp.py`: aplana `git show HEAD:docs/contratos/C2-colores.json` y el actual | 190 tokens en HEAD y 190 ahora. Diferencias en `$value`: 0; en `$type`: 0; en `$extensions`: 0; en otras llaves: 0; en metadatos de grupo: 0. `$description` cambia solo en los 5 tokens de K41/K42 y solo pierde el sufijo « · PROVISIONAL: bloqueo …». En `colorimetria.md`, el diff contra HEAD no agrega ni quita ningún hex. `contraste.csv`, `paleta.html` y `paleta.png` no aparecen en `git diff HEAD --stat`. |
| c | `grep PROVISIONAL` en los 5 artefactos | `grep -a -c PROVISIONAL` | 0 en `colorimetria.md`, `contraste.csv`, `paleta.html`, `paleta.png` y `C2-colores.json`. En el script quedan 19 líneas, todas del mecanismo (diccionario `PROVISIONAL = {}` y comprobaciones); ninguna cita un bloqueo resuelto como vigente. |
| d | Regla K41 para `DESIGN.md` en «Cambios 1.2» | `sed -n '23,49p' colorimetria.md` | Presente: «Regla nueva para `DESIGN.md` (K41)», punto 1 «Una persona y una tienda nunca se distinguen solo por color en una misma vista» y punto 2 «U2-02: entre categóricas basta el refuerzo de U2-08, también entre grupos distintos». |
| e | Cálculo propio OKLab de los 5 tokens | `oklab.py` (matrices de Ottosson, desde el hex del JSON) | Karen `#B870AC` 333.2°, L 0.6429, C 0.1199; `#D389C6` 333.3°, L 0.7251, C 0.1196. Bodega Aurrera `#762C4C` 356.1°, L 0.4096, C 0.1095; `#FFAAC9` 356.2°, L 0.8292, C 0.1063. Walmart oscuro `#7AA3EE` 261.9°, L 0.7165, C 0.1190 (nominal K42: L 0.716). Frente a Ximena (Δ matiz / ΔL / ΔE×100): Karen claro 17.1° / 0.103 / 10.9; Karen oscuro 16.7° / 0.111 / 11.6; Bodega claro 5.8° / 0.130 / 13.1; Bodega oscuro 6.2° / 0.006 / 1.3. Coinciden con «Cambios 1.2», «Cambios 1.1» punto 2, §7.7 y §12 (hallazgo 14). Tiendas en oscuro: ΔL mínimo 0.1008 en OKLab y en gris (K42 se sostiene). |
| f | Nombres de tokens sin cambios | `cmp.py` | Mismo conjunto de 190 nombres; ninguno solo en HEAD ni solo ahora. |

## Criterios

| Criterio | Estado | Evidencia propia |
|---|---|---|
| U2-02 (alta) matiz C y K41 | Cumple | Matiz medido `.light`/`.dark`: ingreso 135.1/134.7 (propuesta 135° documentada), gasto 44.7/45.0 (K26, 40–45°), error 19.5/20.1, éxito 150.0/150.0, advertencia 75.4/75.7, información 239.8/239.6, ahorro 200.0/199.8, deuda 310.2/309.9, acción crítica 100.1/100.1. Karen 333° y Bodega Aurrera 356° con el refuerzo declarado y sin cambio de token; marcas quitadas (verificación c). Contraste de Karen, Bodega Aurrera y Walmart: 12 filas del CSV, todas ≥ 3:1, recalculadas. |
| U2-03 (alta) L fija por paso | Cumple | Marca: desviación máxima de L 0.0010, croma máximo 0.1103; neutro: 0.0009 y 0.0121. |
| U2-04 (alta) 0 pares bajo umbral | Cumple | Recalculé las **202** filas de `contraste.csv` con fórmula WCAG propia: 0 difieren más de ±0.01; 0 hex del CSV distintos del JSON; 0 bajo umbral; 0 con cumple≠true. |
| U2-05 (alta) pares estructurales | Cumple | 900/50 = 11.53 (plan 11.4; §10 documenta «difiere 0.13»); 700/50 = 5.40; 600/100 = 3.47; blanco/marca 600 = 3.98; 100/950 = 14.07; 300/900 = 7.97; marca 300/950 = 10.38; borde neutro 600/50 = 3.71; neutro 400/950 = 7.96. Valores idénticos a HEAD. |
| U2-06 (alta) roles como texto | Cumple | Sin cambio de valores respecto a HEAD (CSV idéntico); las filas de rol como texto están entre las 202 recalculadas, todas sobre umbral. |
| U2-07 (alta) botón primario | Cumple | Marca 700/50 = 5.53; blanco/marca 700 = 5.77; marca 300 con 950 = 10.38. |
| U2-08 (alta) colisiones y refuerzo; K41 y K42 | Cumple | Barrido propio (8 roles, acción crítica, 4 personas, 4 tiendas, marca 700 en claro y 400 en oscuro; 18 entradas por tema): pares con Δ matiz < 20° y ΔL < 0.15 = 6 en claro (ingreso–éxito, Ximena–Karen, Ximena–Bodega Aurrera, Conjunto–Walmart, Conjunto–marca 700, Walmart–marca 700) y 5 en oscuro (ingreso–éxito, Axel–Walmart, Ximena–Karen, Ximena–Bodega Aurrera, Walmart–marca 400). Son exactamente las 11 filas de §7.7, con las mismas cifras, y cada una declara refuerzo. Ninguna categórica ni la marca queda cerca de un rol. Walmart en oscuro `#7AA3EE` sin marca y FIRME por K42 (§7.2, «Cambios 1.1» punto 3, «Cambios 1.2», §12). La prohibición de distinguir persona y tienda solo por color está escrita para `DESIGN.md`. Ver B-02 (redacción). |
| U2-09 (alta) paleta de 8 series; K27 | Cumple | ΔL mínimo entre series consecutivas: claro 0.1745 (OKLab) y 0.1786 (gris); oscuro 0.2184 y 0.2228. Contraste mínimo propio: claro 3.37 (fondo 50 y blanco), oscuro 3.11 (fondos 950 y 900); el script, con todas sus superficies, da 3.14 y 3.10. Series sin cambio respecto a HEAD. |
| K26 | Cumple | Gasto 44.7°/45.0°; error 19.5°/20.1°; gasto–error 25.3° y 24.9° (salida del script, coherente con mis matices). |
| K27 | Cumple | Ver U2-09. |
| K41 | Cumple | 4 tokens sin cambio de valor, sin marca y FIRMES; §12 ya no los propone (M-01 cerrado). |
| K42 | Cumple | `color.store.walmart.dark` = `#7AA3EE`, L nominal 0.716 (medido 0.7165), sin marca. |

## Hallazgos

| Id | Severidad | Criterio | Qué pasa | Ruta | Corrección sugerida |
|---|---|---|---|---|---|
| B-01 | baja | K42, U2-08 (claridad) | Sigue abierto de la ronda 1. La tabla «Marcas que se quitan» de «Cambios 1.2» (línea 29) se titula «Valor sin cambio: hex · matiz · L · croma» y da valores medidos en el hex: Walmart en oscuro «L 0.717», Karen en claro «L 0.643», Bodega Aurrera en oscuro «C 0.106». K42 y el resto del documento dan el nominal (L 0.716, 0.642, croma 0.107). Ambos valores son correctos, pero el encabezado no dice que son medidos. | `docs/design/tools/colorimetria.py`, `md_cambios_12()` (encabezado de la tabla); genera `docs/design/colorimetria.md` línea 29 | Encabezado «Valor sin cambio: hex · matiz · L · croma (medidos en el hex)», o agregar el nominal. Regenerar. |
| B-02 | baja | U2-08, K41 (redacción del refuerzo) | Sigue abierto de la ronda 1. En §7.7, las filas Ximena–Karen (líneas 456 y 463) dicen «orden fijo, etiqueta directa y ΔL ≥ 0.10», sin forma, y agregan «también entre grupos» en un par del mismo grupo; las filas persona–tienda no nombran orden fijo ni forma. El punto 1 de «Regla nueva para `DESIGN.md`» pide «nombre en texto y su forma», sin el orden fijo que U2-08 exige cuando una vista mezcla personas y tiendas. El punto 2 sí nombra los tres elementos, así que U2-08 se cumple entre ambos puntos. | `docs/design/tools/colorimetria.py`: cadenas `ref` de `pares_cercanos()` y punto 1 de `md_cambios_12()` | En las 4 filas de K41, «orden fijo, forma y etiqueta directa»; en el punto 1, «orden fijo, forma (punto, anillo de avatar o marca de serie) y nombre en texto»; quitar «también entre grupos» de las filas del mismo grupo. Regenerar. |
| B-03 | baja | K42, U2-02 (claridad de §12) | Nuevo. «Propuestas sujetas a la puerta 2» (línea 720) sigue proponiendo «la resolución de colisiones categóricas (§7.5…)», y §7.5 incluye que «Walmart ajusta su L en oscuro (K42)». «Ya no se proponen» solo excluye a Karen y Bodega Aurrera, no la L 0.716 de Walmart en oscuro, que K42 ya fijó. La línea de bloqueos que sigue la declara FIRME, así que la ambigüedad es menor. El candado del script ya tiene a Walmart en `nombres_firmes_12()`. | `docs/design/tools/colorimetria.py`, párrafo de §12 en `generar_md()` (línea 1722) | Agregar a «Ya no se proponen» «y Walmart en oscuro en L 0.716 (K42)», de preferencia generado desde `FIRMES_12`. El matiz 262° de Walmart sigue como propuesta, porque K42 no lo decide. Regenerar. |

## Estado de los hallazgos anteriores (`docs/estado/p1c/REVISION-U2a-r1.md`)

| Id | Severidad | Estado ahora |
|---|---|---|
| M-01 (§12 proponía de nuevo a Karen y Bodega Aurrera) | media | **Cerrado.** La línea 720 ya no las incluye en la lista: «la resolución de colisiones categóricas (§7.5; en la versión 1.1, gasto 45°, «Cambios 1.1»)… **Ya no se proponen:** Karen en 333° y Bodega Aurrera en 356°, aceptadas por K41 en la puerta 1b; son FIRMES («Cambios 1.2»)». El script lo genera (línea 1722) con los matices del modelo, y `verificar()` (líneas 2086–2096) falla si una persona o tienda de `FIRMES_12` aparece antes de la marca. `--solo-verificar` da «0 de 3 revisadas». |
| B-01 (encabezado sin «medidos en el hex») | baja | Sigue abierto (B-01 de arriba). |
| B-02 (redacción del refuerzo en §7.7 y punto 1) | baja | Sigue abierto (B-02 de arriba). |
| Bajos de la P1 (B-01 gasto > 45.00° en algunos tokens; H-02, H-03, H-04) | baja | Fuera del encargo de esta ronda. Los valores no cambiaron, así que siguen como los dejó la ronda 1. |
