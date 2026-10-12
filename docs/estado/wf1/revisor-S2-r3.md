# revisor-S2-r3 · WF-1 · 11-oct-2026

Revisor independiente de S2, ronda 3 de 3. No corregí nada. El informe completo está en `docs/estado/wf1/REVISION-S2-r3.md`.

## Hecho
- Veredicto: **NO APROBADO**, con 1 hallazgo medio (M1 r3) y 4 bajos que se arrastran de r2 (B1–B4).
- El M1 de r2 (ventanas desordenadas, traslapadas o contenidas) queda **CERRADO**. Lo comprobé con las 5 consultas de la fixture y con 6 casos míos.
- M1 r3: un hueco sin eventos fechados que cruza un devengo de varios meses manda esa parte a la siguiente ventana declarada (C da 50000 contra 0 en modo meses). La regla no está escrita, y contradice la consecuencia «ventanas contiguas» de la opción A del bloqueo `S2-ventanas`.
- Reejecuté (a)–(g) con mis comandos, y (f) con 5 fixtures propias (35 consultas) y cálculos con `Fraction` y `Decimal`.

## Archivos
- `docs/estado/wf1/REVISION-S2-r3.md` (nuevo).
- `docs/estado/wf1/revisor-S2-r3.md` (este reporte).
- Temporales en `datos/tmp/rev-s2-r3/` (fixtures propias, salidas, copia del oráculo para el generador y scripts): borrados al terminar.

## Verificación
| Comando | Salida resumida |
|---|---|
| `gofmt -l internal/motor_fin cmd/motores` · `go vet ./internal/motor_fin/... ./cmd/motores/...` | vacío · limpio. `gofmt -l .` y `go vet ./...` fallan solo en `precios` y `recetas` (B4) |
| `go test ./internal/motor_fin/... ./cmd/motores/... -count=1 -cover` | ok, 89.2 % y 78.4 %. `-v`: 309 PASS, 0 FAIL, 0 SKIP. `TestCanon` 86, `TestOraculoFixtures` 30, `TestOraculoHojas` 159 |
| `python3 -I oraculo/fin_run_fixtures.py --out datos/tmp/…` | 30 fixtures, 30 coinciden. `cmp` contra `salidas/s2.json`: idéntico |
| `go run ./cmd/motores s2 dump --fixtures ../oraculo/fixtures/s2 --out ../datos/tmp/…` | exit 0. Idéntico byte a byte a la salida del oráculo |
| `python3 -I oraculo/compara.py s2 --oraculo … --motor …` | 2320/2320 ×3, 159 de 159, cruces 43/43, provisionales 933/933, 0 diferencias |
| `fin_genera_fixtures.py` en una copia bajo `datos/tmp` | `diff -rq` contra `oraculo/fixtures/s2` vacío. 0 `HALLAZGO` |
| grep de independencia | El oráculo no lee `backend/`, `datos/tmp` ni `motor-*`. `motor_fin` (código no de prueba) solo usa biblioteca estándar. Las pruebas leen `model.json` y `fixtures/s2` |
| `sha256sum` antes y después | `model.json` 1a570666…bfd4 y `verify3.json` d86087f3…5041, sin cambio. El mtime de `~/reporte-compra-latex/model/verify3.json` no cambió |
| Fixtures propias (`rev-ventanas-hueco`, `rev-orden`, `rev-k`, `rev-31`, `rev-bajos`) | Oráculo = motor en 35/35. Cifras iguales a las mías en K2, K8, K19 (corte 2027-01-20: 251738 = −19899 + 271637), K20/K39 (dic 7201/50410 y 8121/56844; feb 62934 devengado contra 62935 efectivo) y K21 (abr–sep apagado). `[P2,P1]`, `[P1,P3,P2]`, `[B,A]`, `[P1,P1]` y la ventana invertida dan `invalid_argument`. Hueco A + C con devengo: C = 50000 (M1) |
| `Decimal` de S2-08 | `219.88 629.35` |
| `grep float32\|float64` (sin pruebas) · `grep t.Skip` · `grep reporte-compra-latex oraculo/*.py` | 0 · 0 · 0 |

Criterios: S2-01 a S2-15 CUMPLEN, con su evidencia propia en la tabla del informe. S2-16 NO CUMPLE (M1).

## Supuestos
- Ninguno que cambie cifras. Para M1 tomo como referencia la propia consecuencia escrita de la opción A del bloqueo («ventanas contiguas») y el devengo «mitad por mes» de decisiones §B. No decido la regla: propongo que se bloquee.

## Bloqueos
- Ninguno nuevo mío. Propongo una pregunta 4 en `docs/estado/bloqueos/S2-ventanas.md` (la escribe el corrector).

## Pendientes
- Corregir M1 r3: `backend/internal/motor_fin/ledger.go` (`windowed`, `uncovered`), `oraculo/fin_devengo.py` (`sin_cubrir`, `resumen`), `oraculo/fin_genera_fixtures.py` (consulta A + C en `s2-ventanas-sin-cubrir`) y `docs/estado/bloqueos/S2-ventanas.md` (pregunta 4).
- B1–B4: según el orquestador.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S2-interfaz.md` §2.12 y §3.2: decir que las ventanas son contiguas, ascendentes y disjuntas (o lo que se decida en la pregunta 4) y citar `S2-ventanas` (PROVISIONAL). Motivo: M1 r3 y B3.
2. `oraculo/README.md`, `docs/estado/bloqueos/INDICE.md` y `motor-S2-go.md` §Supuestos 2: actualizarlos (B3).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Pregunta 4 de `S2-ventanas`: se recomienda exigir ventanas contiguas (sin huecos), ascendentes y disjuntas. Así «n ventanas consecutivas» coincide con «n meses» cuando las ventanas son las del oráculo, y no hay cifras distintas entre el modo meses y el modo ventanas.

## Consumo aproximado
~37 acciones de herramienta; ~200 KB leídos; ~25 KB escritos (dos reportes de ~16 KB más scripts y fixtures temporales).
