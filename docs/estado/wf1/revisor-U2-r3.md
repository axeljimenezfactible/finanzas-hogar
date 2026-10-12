# revisor-U2-r3 · Revisor independiente de U2 (DESIGN.md + C2), ronda 3 de 3 · WF-1 · 11-oct-2026

## Hecho
- Leí, por partes y dentro de los límites:
  - el protocolo; decisiones §G, §I y §K completas; los criterios U2 1.2 y C2 completos; plan §4.1–§4.4;
  - REVISION-U2-r2 y correccion-U2-r2;
  - DESIGN.md completo, salvo la tabla de §10.5, que verifiqué por script fila por fila;
  - las verificaciones y supuestos de U2b-1, U2b-2, U2c y U2d; contraste-componentes.md; README y package.json de tokens.
- Reejecuté todo en una copia aislada (`datos/tmp/rev-u2r3/repo`): `npm ci`, generate + build dos veces y `cmp` contra lo versionado, `verify`, `check:sync`, `test:tailwind`, regeneración del CSV de contraste, `--check` y `verificar.py`.
- Escribí 4 scripts propios (`rev.py`, `rev2.py`, `rev3.py` y uno en línea) con conversión OKLCH y fórmula WCAG propias:
  - matiz por rol, L y croma por paso;
  - barrido U2-08 y series U2-09;
  - recálculo de las 316 combinaciones de componentes, resueltas desde el par semántico, y de las 202 de `contraste.csv`;
  - hex y rutas de DESIGN contra C2, valores no de color contra C2, uso de tokens y sha256 de las fuentes.
- Comprobé el hallazgo anterior N-1: **cerrado**.
- Dejé el informe en `docs/estado/wf1/REVISION-U2-r3.md`. Veredicto: **NO APROBADO**, con 0 altos, 1 medio nuevo (M-r3-1) y 10 bajos.

## Archivos
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-U2-r3.md`.
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-U2-r3.md` (este reporte).
- Temporales en `datos/tmp/rev-u2r3/`, ya borrados. No toqué ningún artefacto revisado.

## Verificación
- **N-1**: la línea 893 de DESIGN dice «56 × 32 px (`size.bottombar.indicator-width` × `size.bottombar.indicator-height`)». En C2 valen 56px y 32px. «32 × 56» aparece 0 veces. **Cerrado.**
- **C2-01 / C2-09**: `npm ci` termina con código 0; Style Dictionary 5.6.1. Hice gen y build dos veces, con código 0 y 0 advertencias. Los sha256 de las dos corridas son iguales y `cmp` sale IDÉNTICO en C2-tokens.json, tokens.css y tailwind.theme.css. `check:sync` da OK.
- **C2-02 / C2-03 / C2-06 / C2-08 / C2-10**: 445 tokens, 0 sin tipo y 0 referencias rotas. `--ui-*` da 78/78/78 con 0 diferencias. 0 URL. 0 hojas que no son color sin cita en DESIGN.
- **C2-04 / (b)**: 190 de 190 colores con el mismo hex en C2-tokens. 85 hex distintos en DESIGN, 0 fuera de C2-colores. 158 pares token·hex, 0 diferencias. 254 rutas citadas, 0 inexistentes.
- **C2-07**: Tailwind 4.3.3 da «Done». En la salida aparecen `#FAFAFB`, `#5D609F`, `#457723`, `var(--ui-role-income)`, `border-radius: 10px` y `var(--ui-shadow-1)`.
- **U2-04 / C2-05 / (d)**: el CSV regenerado sale idéntico y `--check` termina con código 0 (316 combinaciones, 0 bajo umbral). Mi recálculo de las 316, resueltas por el par: 0 diferencias de hex y de razón, 0 bajo umbral. Mi recálculo de las 202 de `contraste.csv`: 0 diferencias y 0 bajo umbral. La muestra de 12 pares está en el informe.
- **U2-02 / U2-03**: los matices medidos están dentro de ±2° en todos los roles. Gasto: 44.7 / 45.0. Desviación de L: como máximo 0.0016. Croma de roles: 0.1408 como máximo (B-10).
- **U2-05 / U2-06 / U2-07**: los 16 pares estructurales coinciden con §2.8. El texto de rol tiene un mínimo de 4.73. Hay 0 filas de blanco sobre marca 600.
- **U2-08**: 6 pares cercanos en claro y 5 en oscuro, que son los 11 de §2.6 con su refuerzo. Walmart oscuro: `#7AA3EE`, L 0.7165. 0 PROVISIONAL.
- **U2-09**: ΔL consecutivo mínimo de 0.174 en claro y 0.218 en oscuro. Contraste mínimo sobre tarjeta: 3.145 y 3.108.
- **U2-10 / (e)**: el grep contra terceros da 0. Los 7 sha256 están en FUENTES.md. Licencias OFL 1.1, ISC y MIT. 1,877 símbolos.
- **U2-11 / U2-13**: los valores de C2 coinciden con la lista del criterio. «44 px» aparece 23 veces. `linear|stripe|material` = 0.
- **(g) / (h)**:
  - K26: gasto en 45°.
  - K27: ΔL consecutivo con forma, trazo y patrón (§8.6).
  - K41: la regla está escrita en §2.6, §2.10, §7.4 y §8.7.
  - K42: firme.
  - Notion es la única referencia y la saturación está apagada (croma máximo 0.11 en marca).

## Supuestos
- La tolerancia para el croma de los roles (0.1408 contra 0.14) se toma como redondeo a hex, igual que la de ±2° del matiz. Se reporta como bajo (B-10), no como incumplimiento.

## Bloqueos
- Ninguno nuevo. El bloqueo abierto `docs/estado/bloqueos/S1-parametros-avisos.md` está bien marcado como PROVISIONAL en DESIGN §7.13.

## Pendientes
- M-r3-1: corregir las líneas 3 y 23 de `/home/deepsec/finanzas-hogar/docs/design/contraste-componentes.md`, que deben citar `C2-tokens.json` como única fuente.
- Bajos opcionales:
  - B-1, B-2, B-3, B-4, B-6, B-7 y B-10 en `/home/deepsec/finanzas-hogar/docs/design/DESIGN.md`;
  - B-5 en `/home/deepsec/finanzas-hogar/docs/design/tools/tokens/verificar_c2.py`;
  - B-8 en `/home/deepsec/finanzas-hogar/docs/design/contraste-componentes.md`;
  - B-9 en `/home/deepsec/finanzas-hogar/docs/estado/wf1/U2c-tokens.md`.

## Solicitudes fuera de alcance
- Ninguna. Los arreglos son del corrector de U2.

## Paquetes instalados
- Ninguno del sistema. `npm ci` se corrió solo en la copia aislada, ya borrada.

## Propuestas para Axel
- Ninguna. Este alcance no tiene valores de K31.

## Consumo aproximado
- Acciones de herramienta: 37.
- Leído: unos 330 KB. DESIGN.md sin §10.5 ocupa unos 124 KB; criterios, decisiones y plan, unos 45 KB; reportes y revisiones, unos 45 KB; el resto son salidas de scripts.
- Escrito: unos 32 KB (informe, este reporte y scripts temporales).
