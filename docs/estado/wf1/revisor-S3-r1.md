# revisor-S3-r1 · reporte final (WF-1, 11-oct-2026)

## Hecho
- Revisión independiente de S3 (núcleo `motor_cred`, `comparador`, `evaluar`, `canontest`, oráculo `cred_*.py`, fixtures `oraculo/fixtures/s3/`, `S3-escenarios.md`) solo contra `docs/criterios/S3.md` 1.2, decisiones §C/§K y `S3-interfaz.md` 1.4.
- Reejecuté gofmt, vet, pruebas con cobertura, `s3 dump`, oráculo y `compara.py s3`. Escribí mis propios scripts: regresión Kueski, préstamo francés, puntaje mín-máx con restricciones y óptimos, oráculo sobre las fixtures del comparador, y empates del objetivo 4.
- Veredicto: **NO APROBADO** (2 altos, 1 medio, 4 bajos). Informe: `docs/estado/wf1/REVISION-S3-r1.md`.

## Archivos
- `docs/estado/wf1/REVISION-S3-r1.md` (nuevo)
- `docs/estado/wf1/revisor-S3-r1.md` (este)

## Verificación
- gofmt vacío; vet limpio. `go test ./internal/motor_cred/... -cover`: motor_cred 93.7 %, canontest 90.0 %, comparador 96.0 %, evaluar 90.5 % con **FAIL** en `TestOraculoFixtures` por las 9 fixtures del comparador faltantes; total 94.0 %.
- `compara.py s3` con dump propio: 893/893 en las tres direcciones, 437/437 de `model.json`, 9/9 de data.json, 0 diferencias; sha256 de `model.json` intacto.
- Regresión Kueski propia: 17/17. Préstamo francés propio: 3/3. Tarjeta: 0.00 antes y el mismo día; `missing_capture` sin datos.
- Oráculo sobre las fixtures del comparador (copias de testdata): 0 diferencias en 5 fixtures (1,812 hojas). Puntajes, restricciones y óptimos recalculados por mí: 0 violaciones, 46/46 puntajes, 4 óptimos distintos por caso. Tiempos ≤ 1 ms por caso y 114 ms en la carga.
- `cred_escenarios.py` sin argumentos: código 1 por las fixtures faltantes.
- Criterios: cumplen S3-01…S3-11, S3-14…S3-17, S3-19 y S3-20. No cumplen S3-12, S3-13 (parcial) y S3-18, todos por H1/H2.

## Supuestos
- La equivalencia motor = oráculo del comparador la establecí de forma transitiva: el motor iguala al `esperado` por la prueba Go con listas completas, y el oráculo iguala al mismo `esperado` por mi corrida. No sustituye la vía oficial (H1, H2).

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Para el integrador o el corrector: H1 (fixtures en `oraculo/fixtures/s3/`), H2 (`evaluar/evaluar.go:197`), H3 (`TestOraculoFixtures` en `comparador`). Detalle en el informe.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S3-interfaz.md` §8.4 (4): precisión de `value_per_peso_milli` (H4). §8.2: `term_loan` (H6). §2 y §11: `canontest/` y ubicación de `TestOraculoFixtures` (H7).

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva. La propuesta K31 del comparador (mantener 5 s; límites de seguridad de 1 s a 30 s; `date_window_days` = 7) es consistente con lo que medí: 114 ms en el peor caso.

## Consumo aproximado
- 40 acciones de herramienta; ~185 KB leídos; ~17 KB escritos. Temporales `datos/tmp/rev-s3/` borrados.
