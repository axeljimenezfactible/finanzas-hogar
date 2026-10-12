# revisor-U3-r1 · reporte final (WF-1)

Revisor independiente de `docs/design/U3-dashboard.md` v1.0, ronda 1 de 3 · 11-oct-2026 · Criterios: `docs/criterios/U3.md` v1.2, decisiones §G y §K.

## Hecho
- Revisé U3 completo, por partes, contra los 12 criterios U3-01 a U3-12 y las comprobaciones mínimas (a) a (g) del orquestador. Cada verificación la volví a correr con mis propios scripts; no me basé en el reporte de U3-dashboard.
- Informe: `docs/estado/wf1/REVISION-U3-r1.md`. Veredicto **NO APROBADO**: 0 hallazgos altos, 7 medios (H-1 a H-7) y 3 bajos (L-1 a L-3).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-U3-r1.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-U3-r1.md` (este reporte)
- Temporales en `/home/deepsec/finanzas-hogar/datos/tmp/revU3/` (`ops.py`, `check.py`, `nums.py`, `tok.py`, `ops.json`), **borrados** al terminar.

## Verificación
- `uv run --with pyyaml python -I ops.py`: 210 operaciones indexadas desde los 17 extractos (parámetros con `$ref` resueltos y campos de la respuesta 2xx).
- `python3 -I check.py`: 27 de 27 operaciones citadas existen; 0 campos o parámetros inexistentes (lo que no coincidía eran valores de enumeración, y comprobé que existen).
- `python3 -I nums.py`: 42 montos, 18 distintos, 17 en `model.json`/`data.json` y $138.50 = K-03 de `S5.md` (recalculado); 0 fechas ISO; 0 hex y 0 `rgb(`; 25 fechas con día de la semana y 0 errores; semanas ISO (11-oct-2026 = W41, día 7; 14-oct-2026 = W42); hoy + 60 desde el 10-oct-2026 = 9-dic-2026 (base de H-3); días que faltan: 22 (mes, contando hoy) y 20 (pago).
- `python3 -I tok.py`: 57 tokens citados, 0 sin definir en C2; 22 iconos, todos en DESIGN §5.
- Revisión de parámetros: `Period` y `Family` (K5: `household|external|joint`), `status` de `listScheduledPayments` (un solo valor: `pending|paid|overdue`, base de H-4) y de `listSavingsGoals` (`active|archived`).
- C1-datos §9 (líneas 336 a 354): M-01 a M-11 y M-13 a M-17 con el endpoint citado; ninguna métrica para N-121 ni N-162 (H-7).
- DESIGN §4.2 (rejilla, coincide), §7.6 a §7.8 y §7.14 (H-1), §8 (H-2).
- Criterios: U3-01, U3-03, U3-07 y U3-12 cumplen; U3-02 (H-7), U3-04 (H-1), U3-05 (H-2, H-6), U3-06 (H-5), U3-08 (H-3, H-6), U3-10 (H-3, H-4) y U3-11 (H-2) cumplen con observaciones medias; U3-09 cumple (L-1). La evidencia de cada uno está en la tabla de `REVISION-U3-r1.md`.

## Supuestos
1. H-7 queda en media y no en alta porque el «Verifica» de U3-02 pasa: N-121 y N-162 están en C1-datos y sus operaciones existen. La falta de métrica M-xx es un hueco de C1 (C1-09) que U3 tenía que señalar.
2. Califico de media las contradicciones con DESIGN §7 y §8 (H-1, H-2) aunque DESIGN se escribió en paralelo: hoy existen, y quien implemente F2 tendría que adivinar.

## Bloqueos
Ninguno nuevo. `S1-parametros-avisos` no toca U3.

## Pendientes
Ninguno de esta revisión. La ronda 2 verifica H-1 a H-7 y L-1 a L-3.

## Solicitudes fuera de alcance
1. `docs/design/DESIGN.md` §8.3 (U2b): el subtítulo «Deuda viva $6,293.44 = …» del área apilada choca con K19, N-091 y M-11 («no se rotula deuda viva»). Quitarlo o justificarlo.
2. `docs/contratos/C1-datos.md` §9: métricas para «Comidas de hoy» (N-121 vía N-027) y «Diferencias de precio» (N-162 vía N-027), con fórmula, periodo y filtros (C1-09). Motivo: H-7.

## Paquetes instalados
Ninguno (pyyaml de forma efímera con `uv run --with`).

## Propuestas para Axel
Ninguna: U3 no define valores operativos (K31).

## Consumo aproximado
~36 acciones de herramienta; ~180 KB leídos (U3 58.5, decisiones 19, criterios 12.6, DESIGN por partes 25, C1-datos por grep 8, reporte U3 7.4, protocolo 5.6, salidas de scripts y grep ~45); ~35 KB escritos (informe 16.2, reporte 4.5, scripts temporales ~12, progreso 0.5).
