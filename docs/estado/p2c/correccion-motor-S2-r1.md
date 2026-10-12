# correccion-motor-S2-r1 · P2c · corrección de la ronda 1 de motor-S2

Agente: `correccion-motor-S2-r1` · 11-oct-2026, noche · Alcance: `backend/internal/motor_fin/`, `backend/cmd/motores` (s2), `oraculo/fin_*.py`, fixtures s2. Mandato: corregir los hallazgos alto y medio de `docs/estado/p2c/REVISION-motor-S2-r1.md`. Ese informe trae 0 altos y 1 medio (M1). Los bajos B5, B6 y B7 no me los asignaron y no los toqué.

## Hecho
- **M1 (media) · S2 §3 «No aceptable» / S2-interfaz §7, CORREGIDO.** En `backend/internal/motor_fin/cover_test.go:239` (`TestVentanasContiguas`), el esperado de la consulta `[P1] al 2026-11-09` ya no es el literal `97375`. Ahora se lee la hoja canon `model.json:periodos[0].diferencia_acumulada` con `mc()`, igual que en `TestVentanasOrden` (l. 206):
  - antes: `eq(t, "[P1] al 2026-11-09", qc(t, q(…, `{"windows":[`+a+`]}`), "[0].cumulative_diff_cents"), 97375)`
  - ahora: `eq(t, "[P1] al 2026-11-09", Cents(qc(t, q(…, `{"windows":[`+a+`]}`), "[0].cumulative_diff_cents")), mc(t, "periodos[0].diferencia_acumulada"))`
- Comprobé que la ventana `a` de la prueba (2026-10-10 a 2026-11-09) es la misma que `periodos[0]` de `model.json` (`inicio` 2026-10-10, `fin` 2026-11-09, `diferencia_acumulada` 973.75). Así que la hoja leída es la que corresponde.
- Las cifras de la interfaz y de M1 (50000, 100000, −50000, 0 y 292125) se quedan como están, como pide el informe.
- Es el único cambio que hice. No toqué ningún otro archivo de código, del oráculo ni de las fixtures.

## Verificación (comandos y salidas)
| # | Comando | Salida |
|---|---|---|
| 1 | `diff` del archivo contra el respaldo previo | una sola línea cambiada, la 239 |
| 2 | Mutación: cambié la hoja leída a `periodos[2]` (1258.69), corrí `go test ./internal/motor_fin/ -count=1 -run 'TestVentanasContiguas$'` y luego restauré | `FAIL … cover_test.go:239: [P1] al 2026-11-09: motor=97375 esperado=125869`. La aserción sí lee `model.json`. Restauré y `cmp` dio idéntico. Nota: `periodos[1]` también vale 973.75, así que no sirve para la mutación |
| 3 | `cd backend && go test ./internal/motor_fin/ -count=1 -cover` | `ok finanzas-hogar/internal/motor_fin 0.168s coverage: 89.2% of statements` (K25 ≥ 85 %) |
| 4 | `go test ./internal/motor_fin/ -count=1 -v`, contado con awk; `grep -rn t.Skip internal/motor_fin \| wc -l` | 20 pruebas de primer nivel PASS y 294 subpruebas PASS; 0 FAIL y 0 SKIP; 0 `t.Skip` |
| 5 | `gofmt -l internal/motor_fin cmd/motores` y `gofmt -l .` | vacío en los dos |
| 6 | `go vet ./internal/motor_fin/ ./cmd/motores/ && go vet ./...` | limpio (`VET_OK`) |
| 7 | `grep -rn 'float32\|float64' internal/motor_fin --include='*.go' \| grep -v _test.go \| wc -l` | 0 (S2-13) |
| 8 | `-run 'TestVentanasContiguas\|TestInstanteZ\|TestMetasEstado\|TestProyeccionCadencia\|TestCanon$\|TestOraculoFixtures\|TestVentanasOrden' -v` | las 7 PASS |
| 9 | Copia de `oraculo/` y `datos/semilla/` en `datos/tmp/corr-s2-r1/copia`, `python3 -I oraculo/fin_genera_fixtures.py` ahí, `diff -r` contra `oraculo/fixtures/s2` | «32 fixtures y 3 archivos de _datos escritos en oraculo/fixtures/s2»; `diff -r` sin diferencias (32 json) |
| 10 | `python3 -I oraculo/fin_run_fixtures.py --fixtures oraculo/fixtures/s2 --out datos/tmp/corr-s2-r1/oraculo-s2.json`; `cmp` contra `oraculo/salidas/s2.json` | «32 fixtures, 32 coinciden, 0 diferencias»; `cmp` idéntico |
| 11 | `go build -o datos/tmp/corr-s2-r1/motores ./cmd/motores`; `motores s2 dump --fixtures oraculo/fixtures/s2 --out …/motor-s2.json`; `python3 -I oraculo/compara.py s2 --oraculo … --motor …` | exit 0. Oráculo contra esperado, motor contra esperado y oráculo contra motor: 3185 de 3185 hojas cada uno. `model.json` (ramas de s2): 159 de 159. Cruces: 43 de 43. Provisionales: 256 de 256. «diferencias: 0» |
| 12 | `sha256sum` de `oraculo/model.json`, `oraculo/verify3.json` y `datos/semilla/data.json` antes y después (`sha256sum -c`); `git diff --quiet` de esos tres | `1a570666…f1bfd4` = `MODEL_SHA`, `d86087f3…95041` = `VERIFY3_SHA` y `2fe7afae…`: OK antes y después; sin cambios en git |
| 13 | `grep -n '97375\|125869\|125870\|1310378' backend/internal/motor_fin/*_test.go` | Solo queda `engine_test.go:245-246`, y ahí 97375 es una **entrada** de `promedios` (`total_cents`), no un esperado. El revisor ya lo había visto y no es hallazgo |
| 14 | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | exit 0. «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». `K25 OK finanzas-hogar/internal/motor_fin 89.2 %`. s2: «32 fixtures, 32 coinciden, 0 diferencias · salidas/s2.json idéntica», oráculo contra motor 3185 de 3185 y `model.json` 159 de 159. s3, s4 y s5: salidas idénticas y 0 diferencias. `fin_verify`: 159 de 159 |

Borré la carpeta temporal `datos/tmp/corr-s2-r1/` al terminar. No toqué las carpetas temporales de otros agentes.

## Supuestos
1. El informe sugiere `Cents(mc(t, "periodos[0].diferencia_acumulada"))`, «con la misma forma que `TestVentanasOrden`». `mc()` ya devuelve `Cents`, así que envolverlo otra vez no cambia nada. Usé exactamente la forma de `TestVentanasOrden`: `Cents(...)` del lado del motor y `mc(...)` del esperado. Con `eq` (compara `fmt.Sprint`, y `Cents` no tiene `String()`) las dos formas dan lo mismo. Es un detalle de redacción, no una decisión.

## Bloqueos
- Ninguno nuevo. Los PROVISIONAL de la prueba (`S2-ventanas` pregunta 4, casos 1 y 2) siguen como estaban y citan sus bloqueos abiertos. No decidí nada marcado «sin recomendación».

## Pendientes
- Del mismo informe quedan los bajos B5, B6 y B7 y los abiertos de r3 (B2 y parte de B3). No me los asignaron y no los toqué. B6, B7 y lo que falta de B3 son solicitudes a los dueños del contrato y de los bloqueos.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones de herramienta: unas 16.
- KB leídos: unos 30 (el informe de revisión de 14 KB, fragmentos de `cover_test.go` y `helpers_test.go`, y la bitácora de `verify-all`).
- KB escritos: unos 5 (este reporte y una línea de `cover_test.go`).
