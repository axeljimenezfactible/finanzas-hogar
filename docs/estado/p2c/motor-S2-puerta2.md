# P2c · motor-S2-puerta2 (Opus) · 11-oct-2026, noche

Fuente: `docs/contratos/motores/S2-interfaz.md` 1.3 (cambios de la puerta 2), `docs/estado/p2c/contratos-puerta2.md` («Cambios que deben implementar los motores», bloque S2), `docs/estado/wf1/REVISION-S2-r3.md` (M1, B1–B4). Cada lado está escrito con su propia forma; no se copió código entre Go y Python. Ninguna cifra se ajustó a mano: las cifras fijadas en las fixtures nuevas salen de la interfaz y de REVISION-S2-r3 M1, y el oráculo las reproduce sin intervención (el generador no reporta ningún HALLAZGO).

## Hecho
**Motor Go (`backend/internal/motor_fin`)**
- §2.12 regla 4 (PROVISIONAL, `S2-ventanas` pregunta 4, caso 1; Q2-01 = A): `periods` rechaza con `invalid_argument` toda lista con hueco (`start[i] ≠ end[i−1] + 1 día`), aunque el hueco no tenga eventos (`ledger.go`, después de la regla 1). Las reglas 1 a 3 quedan como FIRMES en comentarios y mensajes (Q2-25 a Q2-27). `checkCovered` se conserva: la regla 3 es FIRME y la 4 es PROVISIONAL, así que no dependen una de otra.
- §3.2 caso 2 (PROVISIONAL, sin decidir): sin cambio de cálculo; `windowed` documenta que la parte k sigue a la ventana y no al mes.
- `flujo_quincenal` y `proyeccion`: `provisional` siempre null (Q2-14 = A; Q2-21 a Q2-24 = C). Se quitaron las constantes `blockFlow` y `blockProjection` y la marca por mes ligado. `variable_estimated_cents` y `savings_planned_cents` siguen null con su bloqueo citado (Q2-12 sin respuesta; `S2-ahorro-quincenal`).
- `metas` (§3.13, FIRMES Q2-16 = A y Q2-18 = A): `state` solo `achieved` (solo exige `target_cents` y `saved ≥ target`), `overdue` (`target_date` < hoy) o `in_progress`; se quitó `open`. `months_left` = `max(0, mes(target) − mes(hoy) + 1)`: una fecha anterior a hoy dentro del mes de hoy da 1 (antes daba 0). Meta vencida sin cumplir: todo el restante en el mes de hoy, marcada PROVISIONAL (`S2-metas-vencida`).
- B1 cerrado: tipo nuevo `Instant` (`types.go`) para `params.now` y `calendario {now}`; un instante sin «Z» final (por ejemplo `+00:00` o `-06:00`) es `invalid_date`, igual que en el oráculo (§1.1: «RFC 3339 en UTC con "Z"»).
- Comentarios PROVISIONAL de A, H1, O1, la liga, la cadencia, las metas, el resumen y el presupuesto pasan a FIRME con su Q2.
- Pruebas: `TestVentanasContiguas` (cifras M1: A + C `invalid_argument`; A, B, C = 50000/100000/−50000, 50000, 0 y diferencia 0; meses [oct, dic] dic = 0; `[d, r]` = 292125 y `[P1]` = 97375), `TestInstanteZ`, `TestMetasEstado`; `TestProyeccionCadencia` ahora espera `provisional` null.

**Oráculo Python (`oraculo/fin_*.py`)**
- `fin_devengo.py`: `con_hueco` (regla 4, PROVISIONAL), medida en días entre fin e inicio; docstrings de `sin_cubrir` y `desordenadas` como FIRMES; nota PROVISIONAL del caso 2 en el reparto `k0 + j`.
- `fin_flujo.py` y `fin_proyeccion.py`: `provisional` None; se quitaron `BLOQUEO_LIGA`, `BLOQUEO_CADENCIA` y el conjunto `marcados`.
- `fin_metas.py`: meta vencida (`months_left` 0, sin cumplir) con `plan` = [{mes de hoy, restante}] y `needed_monthly_cents` = restante (PROVISIONAL, `S2-metas-vencida`; antes daba plan vacío y null).
- `fin_presupuesto.py`: comentario de plan 0 FIRME (Q2-20 = A).
- B1 ya estaba bien en el oráculo (exige `Z`).

**Fixtures (`fin_genera_fixtures.py`, 30 → 32; `_datos/` sin cambios)**
- `s2-liga-provisional`, `s2-proyeccion-7-12-cadencia` y `s2-metas`: `provisional` null y fijos `provisional` null en todas las filas; `s2-metas` ya compara `state` (se quitó el recorte `sin_estado`): in_progress, in_progress, achieved, in_progress.
- `s2-ventanas-sin-cubrir`: `provisional` = `S2-ventanas` pregunta 4; consultas nuevas `hueco-sin-eventos` (A1 2026-10-10 a 11-15 y C1 2026-11-30 a 2027-01-09: hueco 11-16 a 11-29 sin eventos que cruza la despensa 1 → `invalid_argument`) y `d-r-caso-2` (fijo `[1].cumulative_diff_cents` = 292125).
- `s2-ventanas-a-mas-c` (nueva): la consulta A + C del contrato con los datos de M1 (semilla sin despensas a crédito más 100000 de contado el 2026-10-10 a 2 meses): `a-mas-c` = `invalid_argument`, `a-b-c` y `meses-oct-dic` con las cifras de M1 fijadas.
- `s2-metas-vencida` (nueva, `provisional` = `S2-metas-vencida`): vencida este mes (months_left 1, overdue, 400), vencida el mes anterior (0, overdue, plan [2026-10, 400]), fondo cumplido sin fecha (achieved), fondo sin monto con fecha pasada (overdue) y fondo sin nada (in_progress).
- `s2-zona-medianoche`: consultas `desfase-cero` (`+00:00`) y `desfase-cdmx` (`-06:00`) = `invalid_date` (B1).
- `oraculo/salidas/s2.json` regenerado con `fin_run_fixtures.py`; `oraculo/README.md` dice 1.3 y 32 fixtures.

## Verificación
- `cd backend && go test ./internal/motor_fin/ ./cmd/motores/ -count=1 -cover` → `motor_fin` ok, cobertura 89.2 % (K25 ≥ 85 %); `cmd/motores` ok 48.7 % (sin umbral en el Makefile; `s5.go` lo cambió otro agente).
- `go test ./... -count=1 -cover` (todo `backend`) → todos ok.
- `gofmt -l internal/motor_fin cmd/motores` → vacío. `go vet ./internal/motor_fin/ ./cmd/motores/` → limpio.
- `python3 -I oraculo/fin_genera_fixtures.py` → «32 fixtures y 3 archivos de _datos escritos», sin HALLAZGO; corrido dos veces, sha256 de las 32 idénticos (determinista).
- `python3 -I oraculo/fin_run_fixtures.py` → «32 fixtures, 32 coinciden, 0 diferencias»; con `--out datos/tmp/oraculo-s2.json`, `cmp` contra `salidas/s2.json` idéntico.
- Desde la raíz: `datos/tmp/motores s2 dump --fixtures oraculo/fixtures/s2 --out datos/tmp/motor-s2.json` → exit 0; `python3 -I oraculo/compara.py s2 --motor datos/tmp/motor-s2.json` → 3185/3185 en las tres comparaciones, model.json 159/159 y cruces 43/43, provisionales 256/256, «diferencias: 0».
- `make verify-wf1-oraculos GOOSE=…`: s2, s3 y s4 en verde; se detiene en s5 (34 diferencias en `s5-carritos-presentaciones`, `s5-mejor-tienda-incompleta` y `s5-profeco-muestra`): es el trabajo en curso del agente de S5 (`pre_calculo.py`, `pre_profeco.py`), fuera de mi alcance.
- `fin_verify.py` → 159 de 159 en las dos vías.
- sha256: `oraculo/model.json` 1a570666…f1bfd4 (= `MODEL_SHA` del Makefile), `oraculo/verify3.json` d86087f3…95041 (= `VERIFY3_SHA`), `datos/semilla/data.json` 2fe7afae…31d31b; `git diff --quiet` de los tres: sin cambios.

## Supuestos
1. HECHO con lectura propia: el código de un instante sin «Z» es `invalid_date`, el mismo que los dos lados ya daban a un instante mal formado; §3.15 enumera el código pero no lo define.
2. HECHO con lectura propia: la consulta A + C del contrato va en una fixture nueva (`s2-ventanas-a-mas-c`) y no en `s2-ventanas-sin-cubrir`, porque una fixture tiene un solo `datos` (FORMATO §2): con la semilla, el hueco de A + C tiene pagos (es la consulta `hueco-con-pagos`, que ya fallaba por la regla 3) y no prueba la regla 4. En `s2-ventanas-sin-cubrir` quedó el equivalente con la semilla (`hueco-sin-eventos`) y la consulta `[d, r]`.
3. Fixture extra `s2-metas-vencida`: el contrato no la pide por nombre, pero la regla PROVISIONAL de la meta vencida y los estados nuevos cambian cifras en los dos lados y sin ella no se cruzan.

## Bloqueos
Ninguno nuevo. Siguen PROVISIONALES, sin decidir: `S2-ventanas` pregunta 4 (casos 1 y 2), `S2-metas-vencida`, `S2-ahorro-quincenal` y `S2-flujo-quincenal` pregunta 1 (Q2-12).

## Pendientes
- B2 sigue abierto (PREGUNTA): un mes o una quincena mal formados dan `invalid_argument` en el motor y `invalid_date` en el oráculo. La regla general de §3 («tipo o rango inválido → `invalid_argument`») apunta al motor, pero con esa misma lectura una fecha mal formada también sería `invalid_argument`, y los dos lados dan `invalid_date`. §3.15 no define `invalid_date`, así que no lo decidí y no agregué la fixture de error.
- B3: solo cerré la parte de `oraculo/README.md`. La tabla de §5 de la interfaz (dos fixtures nuevas y las consultas nuevas) es del dueño del contrato.
- B4 (gofmt y vet de `precios` y `recetas`): fuera de mi alcance. Hoy `go vet ./...` y `gofmt -l` no los revisé fuera de mis paquetes.

## Propuestas para Axel
- Contestar `S2-ventanas` pregunta 4 (casos 1 y 2) y `S2-metas-vencida`; con eso las fixtures `s2-ventanas-sin-cubrir`, `s2-ventanas-a-mas-c` y `s2-metas-vencida` pasan a `provisional` null.
- PROPUESTA para el dueño del contrato: definir en §3.15 cuándo va `invalid_date` (instante, fecha, mes y quincena mal formados) para cerrar B2.

## Consumo aproximado
46 acciones de herramienta · unos 135 KB leídos · unos 45 KB escritos a mano (Go, Python, generador y este reporte), más unos 70 KB de fixtures y salida que escriben los scripts.
