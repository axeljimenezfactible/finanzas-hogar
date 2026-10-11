# revisor-U2a-r2 · Revisión independiente de colorimetría 1.2, ronda 2 de 3 (flujo P1c, 11-oct-2026)

## Hecho
- Revisé colorimetría 1.2 contra `docs/criterios/U2.md` (U2-02 a U2-09) y §K (K26, K27, K41, K42). Reejecuté las verificaciones con mis propios comandos y scripts. No corregí nada.
- Veredicto: **APROBADO**. Hay 0 hallazgos altos y 0 medios. M-01 de la ronda 1 queda cerrado. Siguen abiertos dos bajos (B-01 y B-02) y agrego uno bajo nuevo (B-03). El informe completo está en `docs/estado/p1c/REVISION-U2a-r2.md`.

## Archivos
- `docs/estado/p1c/REVISION-U2a-r2.md` (informe de revisión, nuevo)
- `docs/estado/p1c/revisor-U2a-r2.md` (este reporte, nuevo)
- Temporales en `datos/tmp/rev-u2a-r2/` (`cmp.py`, `oklab.py`, `crit.py`, `u235.py`, copia de C2 en HEAD y diff de md): borrados al terminar.

## Verificación
- `make verify-color`: salida 0. Regenera en espejo y da idénticos byte a byte `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html`; 202 filas y 0 con cumple=false.
- `python3 -I docs/design/tools/colorimetria.py --solo-verificar`: salida 0, «VERIFICACIÓN OK». Da 5 tokens FIRMES con 0 hex distintos, 0 marcas PROVISIONAL y «§12: … 0 de 3 revisadas».
- `cmp.py` contra `git show HEAD:docs/contratos/C2-colores.json`: 190/190 tokens con los mismos nombres. Hay 0 diferencias en `$value`, `$type` y `$extensions`. `$description` cambia solo en los 5 tokens de K41/K42, que pierden la marca. El diff de `colorimetria.md` no toca ningún hex.
- `grep -a -c PROVISIONAL`: 0 en los 5 artefactos. En el script hay 19 líneas, todas del mecanismo.
- `oklab.py`: OKLab propio de los 5 tokens y de Ximena. Los pares frente a Ximena y el escalonamiento de tiendas en oscuro (ΔL mínimo 0.1008) coinciden con el documento.
- `crit.py`:
  - Recalculé las 202 filas de contraste: 0 difieren más de ±0.01 y 0 quedan bajo umbral.
  - Los matices de los roles cumplen U2-02.
  - Mi barrido U2-08 da 6 pares en claro y 5 en oscuro, idénticos a las 11 filas de §7.7.
  - U2-09: ΔL consecutivo mínimo de 0.1745 en claro y 0.2184 en oscuro; contraste ≥ 3.11.
- `u235.py`: la L por paso desvía como máximo 0.001. Los pares estructurales de U2-05 y U2-07 coinciden con §10.

## Supuestos
- Por «los 5 archivos» de la verificación (c) entendí los 5 artefactos generados (`colorimetria.md`, `contraste.csv`, `paleta.html`, `paleta.png`, `C2-colores.json`). Además revisé el script.
- Dejé en su lugar `datos/tmp/verify-wf0/`, la carpeta de trabajo que el propio `Makefile` usa para `verify-color`, por si otro agente corre `make` al mismo tiempo. Mis temporales sí los borré.

## Bloqueos
- Ninguno.

## Pendientes
- Ninguno en mi alcance.

## Solicitudes fuera de alcance
- `docs/design/tools/colorimetria.py` (corrector U2a, opcional porque son bajos): B-01, encabezado «(medidos en el hex)» en `md_cambios_12()`; B-02, unificar el refuerzo de las filas de K41 en §7.7 y el punto 1 de la regla para `DESIGN.md`; B-03, agregar Walmart en oscuro L 0.716 (K42) a «Ya no se proponen» de §12. Después se regenera.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Unas 34 acciones de herramienta.
- Leído: unos 140 KB. Lo forman protocolo (6), §K (15), bloqueos (6), U2.md por partes (8), reportes U2a-1.2 y correccion-U2a-r1 (12), hallazgos de la ronda 1 (3), diff del script (~30), secciones de `colorimetria.md` (~35) y salidas de comandos (~25).
- Escrito: unos 26 KB (scripts temporales ~10, informe ~10, este reporte ~4).
