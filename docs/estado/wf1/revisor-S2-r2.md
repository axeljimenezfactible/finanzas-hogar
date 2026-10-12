# revisor-S2-r2 · Reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Reejecuté las verificaciones (a)–(g) con mis propios comandos.
- Escribí 34 fixtures propias en `datos/tmp/rev-s2r2/fx/`, en 5 tandas, con esperados calculados aparte con `Fraction` desde `model.json`. No leí ni el motor ni el oráculo para calcularlos.
- Comprobé los 3 hallazgos medios de la ronda 1: M2 y M3 están cerrados. M1 está reabierto en parte (ventanas desordenadas o traslapadas).
- Informe: `docs/estado/wf1/REVISION-S2-r2.md`. Veredicto **NO APROBADO**: 1 hallazgo medio y 4 bajos.

## Archivos
- Creados:
  - `docs/estado/wf1/REVISION-S2-r2.md`
  - `docs/estado/wf1/revisor-S2-r2.md`
- Temporales: `datos/tmp/rev-s2r2/`, borrado al terminar.

## Verificación
- (a) `gofmt -l internal/motor_fin cmd/motores`: vacío. `go vet ./internal/motor_fin/... ./cmd/motores/...`: limpio. Cobertura de `motor_fin`: 89.1 %, con 308 PASS, 0 FAIL y 0 SKIP. En todo el repositorio, `gofmt -l .` y `go vet ./...` fallan por `internal/precios` e `internal/recetas` (B4).
- (b) Oráculo, `fin_run_fixtures.py`: «30 fixtures, 30 coinciden, 0 diferencias», idéntico a `salidas/s2.json`. `motores s2 dump`: exit 0. `compara.py s2`: 2320/2320 en los tres cruces, 159/159 en `model.json`, 43/43 cruces, 0 diferencias. `fin_verify.py`: 159/159 ×2.
- (c) Independencia: los greps salen limpios; no hay lecturas cruzadas.
- (d) sha256 sin cambio (1a570666…bfd4 y d86087f3…5041). 0 escrituras fuera del repositorio.
- (e) Los 16 criterios, con evidencia en la tabla del informe. 15 cumplen y S2-16 no (M1).
- (f) Casos propios, 552 comparaciones con 0 diferencias en cifras:
  - K2: `this_month` en medio de una vigencia, `onward` que borra posteriores, sin `scope` → `onward`, renglón que empieza en el mes, `out_of_range`.
  - K8: 10 fechas, incluido el 31-dic en q2 y 2026-W53.
  - K19: 6 cortes nuevos.
  - K20 y K39: 6 promedios y 2 mensuales.
  - K21: interruptor apagado y encendido, más la proyección de 12 meses desde los eventos de `model.json`.
- (g) Los dos reportes declaran 0 hallazgos de cifras. No hay ajustes.

## Supuestos
- Calculé los esperados de deuda viva y de proyección desde los eventos de `model.json` (despensas, pagos y papel de contado). Primero los crucé contra P1–P6: 12 de 12 iguales.

## Bloqueos
- Ninguno nuevo. M1 pide ampliar `docs/estado/bloqueos/S2-ventanas.md`; eso lo hace el corrector.

## Pendientes
- Ninguno en mi alcance.

## Solicitudes fuera de alcance
1. Para el corrector de S2, M1:
   - `backend/internal/motor_fin/ledger.go` y `oraculo/fin_devengo.py`: rechazar ventanas no ascendentes o traslapadas.
   - `oraculo/fin_genera_fixtures.py`: agregar los 4 casos a `s2-ventanas-sin-cubrir`.
   - `docs/estado/bloqueos/S2-ventanas.md`: agregar la pregunta del orden.
2. Para el dueño de la interfaz:
   - `S2-interfaz.md` §2.12: ventanas ascendentes y disjuntas.
   - §1.1 y §3.15: el código para un mes o una quincena mal formados y para instantes sin «Z» (B1 y B2).
   - Aplicar las solicitudes de `correccion-S2-r1.md` (B3).
3. Para los dueños de `internal/precios` e `internal/recetas`: gofmt y vet (B4).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- No hay valores de K31 en este alcance. Para `S2-ventanas` sugiero ampliar la opción A con «ventanas ascendentes y disjuntas». Es la forma más simple de que la diferencia acumulada nunca salga incompleta ni duplicada en silencio.

## Consumo aproximado
- 40 acciones de herramienta.
- Unos 150 KB leídos: protocolo, §B y §K, criterios S2 sin la tabla, secciones de la interfaz, `correccion-S2-r1`, extractos de reportes y de código, y salidas de comandos.
- Unos 40 KB escritos: scripts y fixtures temporales, informe y este reporte.
