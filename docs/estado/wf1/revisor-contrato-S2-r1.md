# WF-1 · revisor-contrato-S2-r1 · Reporte final
Agente: revisor independiente del contrato S2, ronda 1 de 3 · 11-oct-2026 · Alcance de escritura: `docs/estado/wf1/REVISION-contrato-S2-r1.md` y este reporte (temporales en `datos/tmp/rev-s2-r1/`, ya borrados).

## Hecho
- Revisé `docs/contratos/motores/S2-interfaz.md` 1.0, `oraculo/fixtures/FORMATO.md` 1.0, los 4 bloqueos `S2-*` y el reporte `contrato-S2.md` solo contra `docs/criterios/S2.md` 1.2 y decisiones §B y §K.
- Veredicto: **NO APROBADO**, con 0 hallazgos altos, 10 medios y 6 bajos (detalle en `REVISION-contrato-S2-r1.md` §4).

## Archivos
- `docs/estado/wf1/REVISION-contrato-S2-r1.md` (nuevo, 17.2 KB)
- `docs/estado/wf1/revisor-contrato-S2-r1.md` (este reporte; antes fue el archivo de progreso)

## Verificación
- (a) Hojas: `python3 -I hojas.py` → S2 159 (periodos 122, valor_estable 12, horizonte 13, papel 12), S3 437, total 596. `tabla.py` → 20 renglones, 159 rutas distintas, 0 sin dueño, 0 «n» mal declarados. `union.py` contra `S3-interfaz.md` §12 → S2 ∩ S3 = ∅ y unión 596/596.
- (b) S2-01…S2-16 tienen sección (§9 de la interfaz); estado por criterio en `REVISION` §3.
- (c) Las preguntas que le quedarían al implementador están registradas como M1–M8 y M10.
- (d) K2, K8, K14, K19, K20/K39, K21, K44 y K12 se respetan; hay tensión con §B en la regla de reparto nula (M7).
- (e) Las 24 fixtures están enumeradas con fuente; faltan interés por monto (M9) y R-22…R-31 en el flujo (M6).
- Cálculo propio (`calc.py`, día por día con `Fraction`): periodos de oct-2026 a mar-2027 contra `periodos[0..5]` sin diferencias; horizonte 1310378/1436248/39800/218396/239375/125870; identidad K19 en los 7 cierres igual a S2-07; promedios 3141/21988, 6492/45442, 8991/62935, 6282/43976; valor estable 261687/130844 y cruce con `sensibilidad` 15/15; `used_pct` 88; RN-18 500; proyección de los meses 7–12 apagada = §3.12; 0 eventos entre los días 1 y 9.
- `zoneinfo`: medianoche UTC−6 y ISO 2026-W42 como en §3.1. `sha256sum model.json` = `1a5706…bfd4` (sin cambio). `verify3.py:89` todavía escribe en `~/reporte-compra-latex/` (pendiente del agente del oráculo).

## Supuestos
- Evalué la interfaz como contrato: un criterio «cumple» si la interfaz permite implementarlo y probarlo sin adivinar, no si ya hay código (todavía no existe).

## Bloqueos
Ninguno nuevo.

## Pendientes
- Que `contrato-S2` (ronda 2) atienda M1–M10 en `docs/contratos/motores/S2-interfaz.md` y `oraculo/fixtures/FORMATO.md`.
- Que el integrador resuelva la solicitud 1 de `contrato-S2.md` (S2-03 e interés por monto), que decide M9.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S3-interfaz.md` §7 (línea 193): S3 dice que produce `horizonte.deuda_viva_*`; según S2.md lo calcula S2. Alinear (para el revisor de S3).
2. `docs/contratos/motores/S3-interfaz.md` §10: S3 no usa `FORMATO.md` (formato propio y `cred_compara.py`) aunque FORMATO dice que rige para todos los motores; lo decide el integrador (ligado a M2).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: no hay valores operativos de K31 en esta revisión.

## Consumo aproximado
34 acciones de herramienta; ~215 KB leídos (criterios 34.5, interfaz 44.4, FORMATO 7.1, reporte 8.2, decisiones ~14, protocolo 5.6, bloqueos 6.6, extractos del esquema, de C1, de RN y de S3 ~40, salidas de scripts ~25); ~32 KB escritos (revisión 17.2, reporte 4.5, progreso 0.6, scripts temporales ~9).
