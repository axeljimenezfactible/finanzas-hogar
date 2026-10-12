# revisor-U2-r4 · Revisor independiente del sistema de diseño (ronda 4) · WF-1 · 11-oct-2026

Veredicto: **APROBADO (0 altos, 0 medios, 11 bajos)**. El informe completo está en `docs/estado/wf1/REVISION-U2-r4.md`. No corregí nada.

## Hecho
- Leí lo obligatorio: el protocolo completo; decisiones §G, §I y §K completa; los criterios U2.md y C2.md completos; plan §4.1 y §4.4; REVISION-U2-r3; correccion-U2-r3; las secciones de Supuestos, Bloqueos y Pendientes de U2b-1, U2b-2, U2c y U2d; `contraste-componentes.md` completo; DESIGN.md por partes (unos 105 KB de 160 KB, y lo demás por script).
- Hice una copia aislada de `docs/contratos` y `docs/design` en `datos/tmp/rev-u2r4/repo`:
  - `npm ci --offline`.
  - Regeneré C2-tokens.json: IDÉNTICO.
  - Compilé dos veces con Style Dictionary 5.6.1: IDÉNTICO entre sí y frente a lo versionado.
  - Corrí Tailwind 4.3.3, `verificar_c2.py` y `contraste_componentes.py`: el CSV sale IDÉNTICO y `--check` termina con código 0.
- Mi script `datos/tmp/rev-u2r4/scripts/rev.py` (biblioteca estándar, `python3 -I`):
  - Conteos DTCG y referencias.
  - C2-colores ⊂ C2-tokens.
  - Hex y rutas de DESIGN ⊂ C2.
  - Recálculo de las 316 combinaciones y de las 202 de `contraste.csv`.
  - OKLCH propio: matiz por rol, L por paso, croma, barrido de pares cercanos y series.
  - Pares estructurales.
- Comprobé cerrado M-r3-1 y el estado de B-1 a B-10. Encontré B-11 y amplié B-8 y B-9.

## Archivos
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-U2-r4.md`.
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-U2-r4.md` (este reporte).
- Temporales en `datos/tmp/rev-u2r4/`, borrados al terminar.

## Verificación
- C2-01 / C2-09:
  - `npm ci --offline` con código 0; `style-dictionary --version` da 5.6.1.
  - `python3 -I generar_c2_tokens.py` escribe 100073 B; `cmp` IDÉNTICO.
  - Build 1 y build 2 con código 0 y ninguna línea de advertencia o error.
  - sha256 `2b9429bc41836fc5…` y `f66eb866449f391a…`; `cmp` IDÉNTICO contra `docs/contratos`.
- C2-07: `tailwindcss -i prueba/entrada.css` da «Done» (v4.3.3). Los 6 valores aparecen en `c2-prueba.css`.
- C2-02: 445 `$value`, 0 sin tipo y 0 referencias rotas.
- C2-03: 101 pares light/dark, ninguno incompleto.
- C2-04: 190/190 colores con el mismo hex. DESIGN: 85 hex distintos, 0 fuera de C2-colores; 158 pares token·hex con 0 diferencias; 0 rutas inexistentes, salvo 2 plantillas.
- C2-06: tokens.css con 1/2/1 y `--ui-*` 78/78/78, 0 diferencias.
- C2-08: 0 URL; las fuentes de `@font-face` existen.
- C2-10: `verificar_c2.py` da RESULTADO: OK.
- C2-11: `prompt.md` tiene 90 hex distintos, 0 fuera de C2. Las copias en `entrada/` son IDÉNTICAS a las de `docs/contratos`.
- U2-04 / C2-05:
  - `contraste_componentes.py` con código 0: 158 pares, 316 combinaciones, 304 con umbral, 12 exentas y 0 bajo umbral; cruces 0.
  - CSV regenerado IDÉNTICO.
  - Recálculo propio de las 316: 0 diferencias de hex, de razón y de «cumple». Resolución propia `ui.* sobre ui.*`: 194/194.
  - `contraste.csv`: 202 filas, 0 diferencias, 0 bajo umbral, mínimos 4.738 y 3.057.
- U2-02 / U2-03 / U2-08 / U2-09:
  - Matices 700/400 dentro de ±2° de lo declarado; gasto 44.7 / 45.0.
  - Desviación de L ≤ 0.0016; croma: marca 0.1103, neutro 0.0121, roles en claro 0.1408 (B-10), roles en oscuro ≤ 0.1200.
  - Barrido: 6 pares en claro y 5 en oscuro, igual que §2.6. `walmart.dark` = `#7AA3EE` con L 0.7165. 0 PROVISIONAL.
  - Series: ΔL consecutivo mínimo 0.174 / 0.218. Contraste mínimo 3.145 / 3.108 sobre tarjeta.
- U2-05 / U2-07: 11.533, 5.397, 3.466, 3.985, 14.068, 7.975, 10.375, 1.196, 1.895, 3.715 y 6.090. Botón: 5.770 / 5.532 y 10.471 / 10.375. 0 filas con blanco sobre `#7579C2`.
- U2-10:
  - El grep de terceros del criterio da 0 archivos.
  - Los 7 sha256 coinciden con FUENTES.md.
  - Licencias OFL 1.1 e ISC; sprite con 1,877 símbolos.
- U2-11:
  - Los valores de C2 coinciden con la lista del criterio.
  - «44 px» aparece 23 veces y no hay ningún mínimo táctil menor de 44.
  - `prefers-reduced-motion` está en DESIGN y en tokens.css.
- U2-12, U2-13 y U2-14:
  - U2-12: la matriz §7.16 cubre los 12 grupos; los estados de §7.2 a §7.15 los conté por script.
  - U2-13: el grep `linear|stripe|material|figma|apple|airbnb|fluent|carbon` en DESIGN da 0.
  - U2-14: los bloqueos U2 están RESUELTOS.
- Criterios cubiertos: U2-01 a U2-14 y C2-01 a C2-12 cumplen. C2-11 se evaluó solo en lo que ya existe, porque las plantillas `.tex` de S7 todavía no existen.

## Supuestos
- C2-10: acepto la lectura de U2c (Supuesto 12). Que dos medidas valgan lo mismo sin alias no es un duplicado cuando expresan ideas de diseño distintas, como 12px de letra y 12px de espacio.
- B-11 (`colorimetria.md`) es bajo y no medio, porque no es un artefacto de U2. DESIGN §1.1 lo excluye explícitamente y el paquete de la Pista 3 no lo contiene.

## Bloqueos
- Ninguno nuevo. El abierto `docs/estado/bloqueos/S1-parametros-avisos.md` aparece como PROVISIONAL en DESIGN §7.13 y no afecta los criterios U2 ni C2.

## Pendientes
- Ninguno de revisión. Los bajos B-1 a B-11 son ediciones de texto opcionales; las rutas están en el informe.

## Solicitudes fuera de alcance
- `docs/design/colorimetria.md` línea 209: quitar «/Linear» de «tipo Notion/Linear». Motivo: decisiones §G, Notion es la única referencia visual (B-11).
- `docs/estado/wf1/U2c-tokens.md` y `docs/estado/wf1/U2d-contraste.md`: agregar la nota «superado por correccion-U2-r1/r3». Motivo: B-9.

## Paquetes instalados
- Ninguno del sistema. `npm ci --offline` solo dentro de la copia temporal (borrada), desde la caché local de npm.

## Propuestas para Axel
- Ninguna: esta revisión no toca valores de K31.

## Consumo aproximado
- Acciones de herramienta: 46.
- Leído: unos 285 KB.
- Escrito: unos 24 KB (los dos reportes y los scripts temporales).
