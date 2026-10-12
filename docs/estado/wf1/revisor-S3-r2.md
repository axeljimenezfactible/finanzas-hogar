# revisor-S3-r2 · Reporte de agente (11-oct-2026)

## Hecho
- Revisé S3 de forma independiente, ronda 2 de 3. Evalué contra `docs/criterios/S3.md` 1.2, decisiones §C y §K y `S3-interfaz.md` (§8, §10.0 y §11, por partes). No corregí nada.
- Volví a correr todas las verificaciones con mis comandos. Escribí mis propios cálculos para:
  - la regresión Kueski: 218 de 218;
  - el comparador: restricciones, mín-máx, `value_per_peso`, óptimos, recomendación y la salida completa motor contra oráculo;
  - el préstamo francés: 3 de 3 tablas renglón por renglón;
  - el interés de tarjeta: 6,598.
- Comprobé H1, H2 y H3 de r1: los tres están cerrados.
- Encontré 2 hallazgos medios: M1, la salida guardada `oraculo/salidas/s3.json` está desactualizada; M2, el oráculo mide reloj en la fixture de carga y queda al 90 % del límite. También 4 bajos.
- Veredicto: **NO APROBADO**. Informe: `docs/estado/wf1/REVISION-S3-r2.md`.

## Archivos
- Creados: `docs/estado/wf1/REVISION-S3-r2.md` y `docs/estado/wf1/revisor-S3-r2.md`.
- Temporales en `datos/tmp/rev-s3r2/` (binario `motores`, dumps, `kueski.py`, `comp.py`, `loan.py`, cobertura): borrados al terminar.

## Verificación
- (a) gofmt vacío, vet limpio, `go test ./internal/motor_cred/... -count=1 -cover`: PASS en los 4 paquetes (93.7, 90.0, 96.0 y 90.8 %; total 94.0 %), 133 PASS, 0 FAIL y 0 SKIP.
- (b) `compara.py s3` con salidas frescas: 437 de 437 hojas y 906/906 en las tres direcciones, código 0. Mi regresión Kueski da 218 de 218 contra `model.json` y el dump: 800.50, 973.75, 598.82/374.93, 5,034.74, 1,258.68/.70, 883.75/.77, series A y B de 14 valores con mínimos 8,094.50 y 5,000.00, 6,293.44 y 2027-02-28. En mi primera corrida usé la truncación para el pago de referencia de las cotizaciones; el canon usa HALF_UP(monto(1+r)/4), según K14. Corregí mi script.
- (c) Comparador:
  - 3 casos con 4 óptimos distintos, iguales al oráculo y a la tabla a mano.
  - Prueba de propiedad: 300 comparaciones y 0 violaciones.
  - «Faltan restricciones» con las 4 faltantes y sin valores; las claves `optimizer.*` no tienen valor en 00003.
  - Normalización mín-máx recalculada: 46 de 46 puntajes.
  - Tiempo del motor: 0 a 1 ms por caso y 132 ms en la carga. 5 s sembrado, FIRME y configurable. La propuesta K31 de `date_window_days` = 7 días aparece como propuesta, no como valor.
- (d) Tarjeta: 0.00 si se paga el total antes de la fecha de pago o ese mismo día. Préstamo francés: cuadra al centavo, con ajuste en el último pago.
- (e) Independencia comprobada con `grep`; `model.json` con el mismo sha256 antes y después.
- (f) Decisiones aplicadas:
  - K12: «25.8685 %» y «17.27 %».
  - K13: sin tope fijo en el código.
  - K14: modos reproducidos.
  - K15: `TestReposicionParteSeparada`.
  - K16: sin valores por defecto.
  - K17: mín-máx, desviación estándar y unidades por peso.
  - K18: `missing_capture`.
  - K43: 5 s.
  - K44: `TimeoutFromSetting` lee `comparator.timeout_seconds`.
- M1: `compara.py s3` sin `--oraculo` da 2,036 diferencias y código 1. M2: con `taskset` en un núcleo compartido, la carga del oráculo da `timeout (5003 ms)` y `compara.py s3` da 2 diferencias y código 1.

## Supuestos
- Clasifiqué M2 como media y no como alta: con el equipo libre todas las verificaciones pasan, y la carga no es uno de los 3 casos de S3-12 y S3-13. Pero la salida oficial depende de la carga del equipo y el contrato no dice qué debe hacer el oráculo.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- M1: regenerar `oraculo/salidas/s3.json` (integrador).
- M2: decisión en `docs/contratos/motores/S3-interfaz.md` §8.7 y §10.0 sobre el reloj del oráculo, y su aplicación en `oraculo/cred_escenarios.py`.
- B1 a B4: ver el informe.

## Solicitudes fuera de alcance
1. `oraculo/salidas/s3.json`: correr `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json`. Motivo: M1.
2. `docs/contratos/motores/S3-interfaz.md` §8.7: el oráculo no mide reloj cuando `timeout_seconds` es mayor que 0. Motivo: M2.
3. `oraculo/cred_escenarios.py`: aplicar lo que decida el punto 2 y comparar `status` y `evaluated_count` con el `esperado` de cada fixture. Motivo: M2 y B1.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva. El límite de 5 s del comparador (K43) alcanza de sobra: el peor caso medido del motor es 132 ms, el 2.6 % del límite.

## Consumo aproximado
- Acciones de herramienta: ~40.
- Leído: ~150 KB (protocolo; decisiones §C y §K; criterios S3 completos; revisión r1; corrección r1; contrato §8, §10.0 y §11 por partes; salidas de pruebas y de los dumps).
- Escrito: ~28 KB (scripts temporales y los dos reportes).
