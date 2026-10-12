# revisor-contrato-S3-r5 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Revisión independiente, ronda 5, de `docs/contratos/motores/S3-interfaz.md` v1.4 contra `docs/criterios/S3.md` 1.2 y decisiones §C y §K. No corregí nada del artefacto.
- Veredicto: **APROBADO**, con 0 altos, 0 medios y 3 bajos nuevos (R1–R3). El medio P1 de la ronda 4 queda cerrado. Informe: `docs/estado/wf1/REVISION-contrato-S3-r5.md`.

## Archivos
- Creados: `docs/estado/wf1/REVISION-contrato-S3-r5.md`, `docs/estado/wf1/revisor-contrato-S3-r5.md` (este reporte).
- Temporal: `datos/tmp/rev-s3-r5/verifica.py`, ya borrado.

## Verificación
- `python3 -I datos/tmp/rev-s3-r5/verifica.py` (propio; `fractions.Fraction`; lee `model.json`, `data.json` y los patrones de §12 del texto del contrato):
  - hojas: 596 en total; S3 437 (9, 144, 143, 45, 68, 28); S2 159. Patrones de §12: 42 renglones, 0 sin dueño, 0 con dos dueños, 0 conteos distintos (criterio (a), S3-02).
  - motor propio con las reglas de §4, §5.1, §5.2 y §7: 364 de 364 hojas desde la semilla de §10.0, 73 de 73 desde el insumo `sensitivityFromModel` (P1), R-37…R-43 7 de 7, R-05/R-06 9 de 9 contra `data.json`, efectivo P1…P6 6 de 6, horizonte 1310378 y deuda viva 629344 iguales (S3-01, S3-02, S3-03, S3-04, S3-05, S3-19). 0 diferencias.
- Lectura de §13: S3-01…S3-20 con sección (criterio (b)).
- `grep` en `00003_puerta1b.sql` y `00001`: `comparator.timeout_seconds` = 5; `optimizer.*` de restricciones sin valor; `optimizer.weights` 1,1,1,1 (K16, K43, K44).
- `ls docs/estado/bloqueos/`: los 6 bloqueos que cita el contrato existen (S3-20).
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` al inicio y al final (S3-16).
- Código: no se evalúa en esta revisión (el contrato no tiene código; `motor_cred` está a medio escribir y `oraculo/fixtures/s3/` está vacío).

## Supuestos
- Evalué S3-01…S3-20 como criterios sobre el contrato: si alcanza para implementar y si sus cifras son exactas. Las corridas Go y de `compara.py s3` les tocan al revisor de motores cuando exista el código.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Bajos opcionales para el corrector o el integrador: R1, R2, R3 (este informe), Q1–Q4 (ronda 4) y B1–B9 (ronda 3).

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones de herramienta: 31.
- Leído: unos 255 KB (protocolo 5.6; decisiones §C y §K 16.5; criterios S3 35; contrato completo por partes 134; reporte del contrato 7.3; corrección r4 9; partes de la revisión r4 6; salidas de scripts y `grep` unos 40).
- Escrito: unos 26 KB (script temporal 10; informe 10.7; este reporte 3).
