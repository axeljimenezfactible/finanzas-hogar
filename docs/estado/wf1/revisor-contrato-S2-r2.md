# WF-1 · revisor-contrato-S2-r2 · Reporte final
Agente: revisor-contrato-S2-r2 (revisor independiente, ronda 2 de 3) · 11-oct-2026 · Alcance de escritura: `docs/estado/wf1/REVISION-contrato-S2-r2.md` y este reporte (más temporales en `datos/tmp/rev-s2-r2/`, ya borrados).

## Hecho
- Revisé `docs/contratos/motores/S2-interfaz.md` 1.1 y `oraculo/fixtures/FORMATO.md` 1.1 contra `docs/criterios/S2.md` 1.2 y decisiones §B y §K. Veredicto: **NO APROBADO**, con 0 altos, 6 medios y 5 bajos.
- Comprobé cerrados los 10 hallazgos medios de la r1.
- Hice cuatro verificaciones propias, sin apoyarme en los reportes: conteo de hojas, expansión de la tabla §4, cálculo independiente de las cifras y validación del esquema de insumo-s3.

## Archivos
- `docs/estado/wf1/REVISION-contrato-S2-r2.md` (nuevo)
- `docs/estado/wf1/revisor-contrato-S2-r2.md` (este reporte)

## Verificación
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`.
- Hojas numéricas (`python3 -I`): 596 en total; S2 = 159 (periodos 122, valor_estable 12, horizonte 13, papel 12); S3 = 437.
- Tabla §4 con mi propio expansor de llaves: 159 rutas distintas, 0 sin dueño, 0 sobrantes, 0 en ramas de S3.
- Cálculo propio con `Fraction`, que parte solo de las hojas S3: 159/159 hojas S2, sensibilidad 15/15, interés por monto 28/28, las 7 filas de K19, promedios de K39, flujo de 16 quincenas con su contraste (+1 ×9, −1 ×3), calendario y zona, presupuesto y metas. Resultado: **242 OK, 0 FALLA**.
- §5.1: 55 filas S3 exactas; R-67 con 85 llaves; R-04 en binary64 igual; R-05 = 17.27.
- `grep reporte-compra-latex oraculo/*.py` → `verify3.py:89` (hallazgo M6).
- Criterios S2-01…S2-16: los 16 tienen sección. Estado y evidencia en la tabla de la revisión.

## Supuestos
- Evalué el contrato como especificación para dos implementadores independientes. Por eso cuento como hallazgo medio cada dato que tendrían que adivinar, como lo pide la tarea.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Corregir M1–M6 (y, si se quiere, B1–B5) en `docs/contratos/motores/S2-interfaz.md` y `oraculo/fixtures/FORMATO.md`. Después, la ronda 3 de revisión.

## Solicitudes fuera de alcance
- `oraculo/verify3.py:89`: asignar la corrección que exige S2-04 (ver M6).
- `docs/contratos/motores/S3-interfaz.md` 1.1: citar FORMATO 1.1 §4.1 (`cmd/motores/s3.go` con `register`). Ya pedido en la solicitud 1 de `correccion-contrato-S2-r1.md`.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: el alcance no tiene valores de K31.

## Consumo aproximado
32 acciones; ~180 KB leídos (protocolo 6, decisiones §B y §K 12, criterios S2 35, interfaz 63, FORMATO 15, reportes 17, extractos, esquema, S3 y salidas ~32); ~35 KB escritos (scripts ~14, revisión 13, reportes ~5).
