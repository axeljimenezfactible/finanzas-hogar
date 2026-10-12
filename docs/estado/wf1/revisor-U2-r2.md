# revisor-U2-r2 · WF-1 · 11-oct-2026

Revisor independiente de U2 (DESIGN.md + contrato C2), ronda 2 de 3. El informe completo está en `docs/estado/wf1/REVISION-U2-r2.md`.

## Hecho
- Lecturas obligatorias: protocolo completo; decisiones §G, §I y §K completa; criterios U2.md completo y C2.md §2–§3; plan §4.1 y §4.4; DESIGN.md por secciones (§1–§10.4, cabecera de §10.5, §10.6 y Supuestos); README, config y sincronía de tokens; FUENTES.md; corrección y revisión r1; encabezados de los reportes U2b-1, U2b-2, U2c y U2d.
- Copia aislada en `datos/tmp/rev-u2r2/repo`: `npm ci`, `generate` dos veces, build dos veces, `cmp` contra los versionados, `verificar_c2.py`, Tailwind, `contraste_componentes.py` con y sin `--check`, `verificar.py` y prueba de mutación.
- Scripts propios con biblioteca estándar (rev.py, rev2.py y otros en línea): aplanado DTCG, hex de DESIGN contra C2-colores, tokens citados, valores no de color, recálculo de las 316 combinaciones, OKLCH (matiz, L y croma), barrido U2-08, series U2-09 (L, gris y contraste), pares estructurales U2-05, duplicados C2-10 y sha256 de fuentes.
- Veredicto: NO APROBADO por 1 hallazgo medio nuevo (N-1). Los 5 hallazgos altos y medios de r1 quedaron cerrados.

## Archivos
- `docs/estado/wf1/REVISION-U2-r2.md` (creado)
- `docs/estado/wf1/revisor-U2-r2.md` (este reporte; antes tenía el progreso)
- Temporales `datos/tmp/rev-u2r2/` (borrados)

## Verificación
- `npm ci --prefer-offline` → código 0. `npx style-dictionary --version` → 5.6.1.
- `python3 -I generar_c2_tokens.py`, dos veces → C2-tokens.json idéntico al versionado (100073 B).
- `style-dictionary build --config config.mjs` → código 0. tokens.css `2b9429bc41836fc5…` y tailwind.theme.css `f66eb866449f391a…`, idénticos a los versionados; dos compilaciones idénticas.
- `verificar_c2.py` → RESULTADO OK. Tailwind 4.3.3 «Done», con los 6 valores presentes.
- `contraste_componentes.py --check` → «CSV en disco IDENTICO», 316 combinaciones, 0 bajo umbral y código 0. La regeneración sin `--check` sale idéntica (`cmp`).
- `verificar.py .` → 11 encabezados en orden, 19 contrastes de §2.8 coinciden, 0 terceros y RESULTADO OK.
- Recálculo propio de las 316 combinaciones: 0 discrepancias de token, hex ni razón (±0.01) y 0 bajo umbral. Holgura mínima: C-138 en claro, 4.51.
- Mutación (ui.text.default.light → neutral 800; brand.300 → #C0C6FF): cambian 42 filas y el código es 1 por los cruces secundarios. El CSV sigue a C2-tokens.
- Hex de DESIGN: 162, todos en C2-colores. 190 colores de C2-colores en C2-tokens con el mismo hex. 242 tokens citados, 0 inexistentes (aparte de los patrones N).
- U2-08: 6 pares en claro y 5 en oscuro, iguales a §2.6. U2-09: ΔL consecutivo ≥ 0.174 (claro) y ≥ 0.218 (oscuro) en OKLCH; contraste ≥ 3.11.
- sha256 de las 8 fuentes y licencias = FUENTES.md. Grep de terceros = 0.
- Criterios cubiertos con evidencia propia: U2-01…U2-14 y C2-01…C2-12. C2-11 queda sin evaluar porque el paquete de la Pista 3 y las plantillas de S7 aún no existen.

## Supuestos
1. La convención «ancho × alto» de DESIGN la deduzco de sus otras medidas (asa 32 × 4, pista 36 × 20). Sobre esa lectura descansa N-1.
2. Las 18 medidas no de color con el mismo valor no incumplen C2-10, porque §10.2 declara la regla y el revisor r1 no lo objetó.

## Bloqueos
Ninguno.

## Pendientes
- Ronda 3: corregir N-1 en `docs/design/DESIGN.md` línea 893. De paso conviene corregir B-1 a B-4 y B-7 (texto de DESIGN.md) y B-5 y B-6 (mejoras opcionales).

## Solicitudes fuera de alcance
Ninguna nueva. El Supuesto 9 de DESIGN (constantes de interfaz fuera de K44) queda anotado como observación fuera de mi alcance.

## Paquetes instalados
Ninguno del sistema. `npm ci` solo dentro de la copia temporal, desde la caché local, con las versiones fijas del package-lock. Se borró con la copia.

## Propuestas para Axel
Ninguna.

## Consumo aproximado
- 49 acciones de herramienta.
- ~280 KB leídos.
- ~32 KB escritos: informe 12.7 KB, este reporte, progreso y scripts temporales.
