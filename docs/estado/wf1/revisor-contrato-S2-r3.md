# WF-1 · revisor-contrato-S2-r3 · Reporte final
Agente: `revisor-contrato-S2-r3` (revisor independiente, ronda 3 de 3) · 11-oct-2026 · Alcance de escritura: `docs/estado/wf1/REVISION-contrato-S2-r3.md`, este reporte y temporales en `datos/tmp/rev-s2-r3/` (ya borrados). No corregí nada de lo revisado.

## Hecho
- Revisé `docs/contratos/motores/S2-interfaz.md` 1.2 completo (por partes) y `oraculo/fixtures/FORMATO.md` 1.2 contra `docs/criterios/S2.md` 1.2 y decisiones §B y §K. Leí `docs/estado/wf1/contrato-S2.md` y `correccion-contrato-S2-r2.md` solo como contexto.
- Reejecuté por mi cuenta (a)–(e) y las cifras: conteo de hojas, la tabla §4 expandida y su cruce con S3 §12, la receta de la semilla con mi propio código (159/159), K19, K39, los cruces, las series de `interes-por-monto`, los ejemplos de presupuesto, metas, flujo y 31-oct, la copia de `verify3.py` redirigida y la zona.
- Veredicto: **APROBADO** (0 altos, 0 medios, 13 bajos). Los 6 medios de la r2 quedaron cerrados.

## Archivos
- `docs/estado/wf1/REVISION-contrato-S2-r3.md` (nuevo, 12.0 KB)
- `docs/estado/wf1/revisor-contrato-S2-r3.md` (este reporte)

## Verificación
- `python3 -I hojas.py oraculo/model.json docs/contratos/motores/S2-interfaz.md` → `por rama: parametros 9, periodos 122, despensas 144, credito 143, valor_estable 12, sensibilidad 45, papel 12, disposicion_por_despensa 68, horizonte 13, interes_por_monto 28 · total 596 · S2 159 S3 437 · renglones 20 rutas 159 distintas 159 duplicadas 0 n mal 0 · no son hoja [] · S2 sin dueño 0 · fuera de S2 []`.
- Ramas de la tabla §12 de `S3-interfaz.md` (awk) → solo credito, despensas, disposicion_por_despensa, interes_por_monto, parametros y sensibilidad: 0 duplicadas con S2.
- `python3 -I recalc.py oraculo/model.json`:
  - `validaciones semilla OK` y `hojas S2 desde la receta: 159 coinciden; 0 difieren`.
  - K19: 7 de 7 con `True` (sin saldo solo en 11-30, 01-31 y 03-31).
  - Promedios 3141/21988, 6492/45442, 8991/62935 y 6282/43976.
  - `sens 0..4 True`; `ipm 0..6` con intereses ≥ 0 e intereses = total − monto.
- `sha256sum oraculo/model.json oraculo/verify3.json` → `1a570666…bfd4` y `d86087f3…5041`, sin cambio. Copia de `verify3.py` con la salida en `datos/tmp` → `rc=0` y `cmp` IGUALES. `grep -n reporte-compra-latex oraculo/*.py` → 1 línea (`verify3.py:89`), pendiente de la implementación como asigna §0.5.
- Script de §9 → 16 de 16 criterios. Secciones citadas: todas existen (§0.3–§0.5 son incisos). Están los 5 bloqueos `S2-*`.
- `zoneinfo`: 2026-10-31T05:59:59Z → 2026-10-30; 06:00:00Z → 2026-10-31; 2027-06-15T05:59:59Z → 2027-06-14 (UTC−6 todo el año).
- Criterios de la fase: S2-01…S2-16 cubiertos por el contrato, con la evidencia de cada fila en la tabla de `REVISION-contrato-S2-r3.md`.

## Supuestos
- S2-04 se evalúa a nivel contrato: la línea de `verify3.py` se corrige en la implementación (asignada en §0.5). No es un hallazgo de esta revisión.
- Que `s2-metas` y `s2-alertas-80-100` tengan varios casos en un archivo cumple S2-11 y S2-12 si se declara la equivalencia (bajo B13).

## Bloqueos
Ninguno nuevo. Siguen abiertos los del contrato: `S2-proyeccion-interruptor`, `S2-flujo-quincenal`, `S2-metas`, `S2-metas-resumen` y `S2-presupuesto-total`. `S1-parametros-avisos` no afecta a S2.

## Pendientes
- Opcional, a juicio del integrador: los 13 bajos B1–B13 de `REVISION-contrato-S2-r3.md`, en `docs/contratos/motores/S2-interfaz.md` (§2, §2.1, §2.7, §3.2, §3.6, §3.8, §3.9, §3.11, §3.12, §5, §5.0, §6, §7, §0.5) y en `oraculo/fixtures/FORMATO.md` (§2.1, §5).
- Implementación (fuera de esta revisión): `oraculo/verify3.py:89` → `datos/tmp/verify3.json` (agente Python de S2).

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S2-interfaz.md` §0.5: prohibir importar `oraculo/model.py` desde `fin_*.py`, porque reescribe `model.json` al importarse (líneas 208–209). Motivo: S2-04 y la independencia de S2-03 (B10).
2. `docs/criterios/S2.md` S2-03: el Verifica corre `fin_run_fixtures.py` sin argumentos. Alinear con §6 o declarar valores por omisión en la interfaz (B9).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: no hay valores K31 en este alcance.

## Consumo aproximado
- 29 acciones de herramienta.
- Lectura ≈ 225 KB: protocolo 6, §B y §K 13, S2.md 35, interfaz 80, FORMATO 16, reportes 20, extractos de S3, RN y bloqueo ≈ 8, salidas de scripts ≈ 15.
- Escritura ≈ 32 KB: scripts temporales ≈ 14, revisión 12, reporte 5.
