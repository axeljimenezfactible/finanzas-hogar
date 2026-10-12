# revisor-contrato-S3-r4 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Revisé `docs/contratos/motores/S3-interfaz.md` v1.3 completo, por partes, contra `docs/criterios/S3.md` 1.2 y las decisiones §C y §K. Veredicto: **NO APROBADO**, con 0 hallazgos altos, 1 medio (P1) y 13 bajos.
- Comprobé los 4 medios de la ronda 3: N1, N2 y N3 cerrados; N4 cerrado salvo el residuo P1.
- Reejecuté las verificaciones con un script propio: conteo y cruce de las 437 hojas, reproducción de las 437 con las reglas del contrato, R-05 y R-06, ejemplos, desviación estándar exacta, normalización, ejemplo K15, esquema y extractos de C1.

## Archivos
- `docs/estado/wf1/REVISION-contrato-S3-r4.md` (informe: veredicto, verificaciones, tabla de criterios, estado de r3, hallazgos con corrección sugerida)
- `docs/estado/wf1/revisor-contrato-S3-r4.md` (este reporte)
- Temporal `datos/tmp/rev-s3-r4/verifica.py` (borrado al terminar)

## Verificación
- `sha256sum oraculo/model.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto).
- `python3 -I datos/tmp/rev-s3-r4/verifica.py`:
  - hojas por rama 9/144/143/45/68/28 = 437; S2 159; total 596;
  - 42 renglones de §12: 0 sin dueño, 0 con dos dueños, 0 conteos distintos;
  - reproducción propia: 0 diferencias en 437 hojas, fechas y tipos;
  - R-05 y R-06: 9 de 9;
  - efectivo P1…P6, horizonte 1310378 y deuda viva 629344: iguales;
  - desviación estándar: 81650 en el ejemplo y 0 discrepancias en 20,000 series aleatorias contra `Decimal` de 60 dígitos;
  - normalización: 1, 2/3, 0;
  - ejemplo 1 de §10.3: igual;
  - intereses por pago negativos con intereses totales de 1 a 4 centavos: 14 casos (Q1).
- `sed`/`grep` en 00001 (l. 441, 448, 808), 00002 (`rate_decimals`), 00003 (l. 54–55, 86) y en los extractos `creditos-a.yaml` y `creditos-b.yaml`: confirman las citas del contrato.
- Criterios: S3-01…S3-20 tienen sección y cumplen en el contrato; S3-01 y S3-02 con el medio P1 (tabla en el informe).

## Supuestos
- Mantuve como bajos los bajos de r3 que siguen abiertos (B1–B9), como los clasificó el revisor anterior, porque S3-08, las fixtures o el diseño de los casos los resuelven sin cambiar cifras.
- P1 lo califiqué medio: es una pregunta que el núcleo tendría que adivinar en la prueba de S3-01 y S3-02 (criterios altos), y es parte de lo que pedía r3-N4.

## Bloqueos
Ninguno nuevo. Revisé `docs/estado/bloqueos/S3-serie-sin-intereses.md` (formato correcto) y `docs/estado/bloqueos/U4-iva-credito.md` (citado por el contrato).

## Pendientes
- Para el corrector del contrato: P1 (§10.0 y §11.1); bajos Q1–Q4 y B1–B9 opcionales (rutas y correcciones en el informe).

## Solicitudes fuera de alcance
1. `docs/estado/bloqueos/S3-formula-tarjeta.md` y `S3-convencion-prestamo.md`: actualizar las citas a la v1.3 (B2) y la frase «quedan FIRMES» de los meses sin intereses y abonos (Q4). Motivo: trazabilidad.
2. `docs/estado/wf1/contrato-S3.md`: nota inicial que remita a la v1.3 y a `correccion-contrato-S3-r1..r3.md` (Q4).
3. `docs/estado/bloqueos/INDICE.md` no lista ningún bloqueo de WF-1 de S3 ni `U4-iva-credito`; el orquestador decide si lo agrega para la puerta 2.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (soy revisor; no propongo valores K31).

## Consumo aproximado
27 acciones de herramienta; unos 245 KB leídos (contrato 129.5 KB completo por partes, criterios 35 KB, FORMATO 16 KB, decisiones §C y §K 16 KB, protocolo 5.6 KB, reporte 7 KB, bloqueos 9 KB, partes de r3, esquema y extractos); unos 20 KB escritos.
