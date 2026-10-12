# revisor-U2-r1 · reporte final (WF-1, 11-oct-2026)

Revisor independiente de U2 y C2, ronda 1 de 3. No escribí nada de lo revisado y no corregí nada. Veredicto: **NO APROBADO** (2 altos, 4 medios, 3 bajos). El informe completo está en `docs/estado/wf1/REVISION-U2-r1.md`. cerrado_por_tope: no.

## Hecho
- Lecturas: protocolo, decisiones §G, §I y §K, plan §4, criterios U2 1.2 y C2 1.2, reportes de U2b-1, U2b-2, U2c y U2d, y DESIGN.md por partes (§1–§9, §10.1–§10.4, §10.6, Supuestos y cabecera de §10.5). También fragmentos de colorimetría (§2.7, §5.1, §8.2–§8.4), tokens.css, tailwind.theme.css, README y config de tokens.
- Scripts propios en `datos/tmp/rev-u2r1/` (borrados al terminar): OKLab/OKLCH y contraste WCAG escritos sin importar nada del repo. Con ellos revisé el matiz, la L y el croma de las 11 escalas, los pares estructurales, los roles como texto, el barrido de pares cercanos (U2-08), la paleta de series, el formato DTCG y las referencias, los hex y nombres de DESIGN contra C2, el mapeo de §7.1, el recálculo de las 316 filas desde `C2-tokens.json`, el recálculo de los 202 pares de `contraste.csv` y los conjuntos `--ui-*` por tema.
- Compilación reproducible: `rm -rf node_modules && npm ci`, luego `build` y `cmp` con los entregados; `generate` y `cmp` de C2-tokens.json; `check:sync`; `test:tailwind`.

## Archivos
- `docs/estado/wf1/REVISION-U2-r1.md` (nuevo)
- `docs/estado/wf1/revisor-U2-r1.md` (este)
- Sin cambio de contenido, aunque se reescribieron al correr sus herramientas: `docs/contratos/tokens.css`, `docs/contratos/tailwind.theme.css`, `docs/contratos/C2-tokens.json` y `docs/design/contraste-componentes.csv`. Los cuatro dan `cmp` idéntico contra copias previas. `docs/design/tools/tokens/node_modules` se reinstaló con `npm ci` y está ignorado por git.

## Verificación
- `npm ci --prefer-offline` y `npm run build`: código 0; `style-dictionary --version` = 5.6.1. `cmp` da tokens.css y tailwind.theme.css idénticos. `npm run generate` y `cmp` dan C2-tokens.json idéntico (sha256 `37ad40f5…570d`).
- `npm run check:sync`: OK con sha256 `495664…1c88` y `f66eb8…c88a`. `npm run test:tailwind`: código 0, con los valores de token presentes; salida borrada.
- `npm run verify`: RESULTADO OK. `python3 -I docs/design/fuentes/scripts/verificar.py .`: RESULTADO OK. `python3 -I docs/design/tools/contraste_componentes.py`: 158/316/304/12, 0 bajo umbral, salida 0, y el CSV queda idéntico (`cmp`). `--check` también sale con 0.
- Recálculo propio: 316 de 316 filas resueltas desde `C2-tokens.json`, 0 diferencias y 0 bajo umbral. `contraste.csv`: 202 filas, 0 bajo umbral, diferencia 0.
- U2-02: matices dentro de ±2° y gasto en 44.7–45.0°. U2-03: \|ΔL\| ≤ 0.002 y croma máx. 0.1103/0.0121/0.1408 por el redondeo a hex. U2-05: 8 pares que coinciden o tienen la diferencia documentada. U2-06: texto de rol con mínimo 4.73. U2-07: 5.77, 5.53, 10.47 y 10.37. U2-08: 6 y 5 pares cercanos, iguales a la tabla de DESIGN y con refuerzo; Walmart oscuro `#7AA3EE` L 0.717; 0 PROVISIONAL. U2-09: ΔL consecutivo 0.174/0.218 y contraste 3.145/3.108.
- C2-02: 445 tokens, 0 sin valor, 0 sin tipo, 0 referencias rotas. C2-04: 190/190 hex iguales, 162 hex de DESIGN con 0 fuera, 158 pares sin diferencias, §7.1 con 41/41; pero 4 nombres son grupos (H-2). C2-06: `--ui-*` 78/78/78, iguales. C2-08: 0 URL. C2-10: 36 grupos de color duplicados sin alias (M-2).
- U2-10: `woff2` con firma `wOF2`, TTF con `tnum`, licencias OFL e ISC, 1,877 símbolos y 0 trazos fijos; grep de terceros = 0. U2-13: grep de otras referencias visuales = 0.

## Supuestos
1. Para U2-03 tomé como evidencia de la L medida mi propio cálculo y colorimetría §5. DESIGN §2.3 lista la L fija y no la medida, y remite a colorimetría.
2. «Toda variable de color en ambos temas» (C2-06) la apliqué a las temáticas `--ui-*`. Los primitivos `--color-*` son estáticos por diseño (§10.3).
3. C2-11: el consumo (`prompt.md` y `.tex`) lo verifican sus agentes porque esos archivos no existen todavía.

## Bloqueos
Ninguno nuevo. El abierto `docs/estado/bloqueos/S1-parametros-avisos.md` está bien citado como PROVISIONAL en DESIGN §7.13.

## Pendientes
- Ronda 2: corregir H-1, H-2 y M-1 a M-3 (rutas en el informe §2) y volver a correr los tres verificadores.

## Solicitudes fuera de alcance
1. `docs/design/tools/contraste_componentes.py`: que lea `C2-tokens.json` (H-1). Motivo: C2-05.
2. `docs/design/DESIGN.md` §4.3, §4.6, §9.3, §10.1, §10.2, §10.4 y §10.6: sincronizar con C2 v1.0 (H-2) y aclarar la barra inferior (M-1). Motivo: C2-04 y U2-11.
3. `docs/design/tools/tokens/generar_c2_tokens.py` y `verificar_c2.py`: alias en los colores duplicados y revisión de duplicados en los colores (M-2). Motivo: C2-10.
4. `docs/design/tools/tokens/README.md` y DESIGN §10.1, §10.5 y §2.10: una sola cadena de fuente y la ruta de la tabla de contraste de componentes (M-3).

## Paquetes instalados
Ninguno del sistema. `npm ci` reinstaló las dependencias fijadas en `docs/design/tools/tokens/package-lock.json` (desde la caché de npm, `--prefer-offline`).

## Propuestas para Axel
Ninguna de K31. Una sola cosa para la puerta 2: confirmar la medida de los destinos de la barra inferior (56 de alto × ≥ 64 de ancho, como sugiere DESIGN) cuando U2 la deje sin ambigüedad.

## Consumo aproximado
41 acciones de herramienta · ~280 KB leídos · ~42 KB escritos (scripts temporales ~16 KB, informe 16 KB, reporte ~6 KB, progreso ~1 KB).
