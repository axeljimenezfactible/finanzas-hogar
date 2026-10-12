# Reporte · motor-S3-comparador-go (WF-1, reanudación del 11-oct-2026)

## Hecho
- Paquete nuevo `backend/internal/motor_cred/comparador` (S3-interfaz.md §8, ronda 4): tipos de §8 (`Request`, `Constraints`, `Weights`, `PersonCap`, `BasketQuote`, `Scenario`, `Result`…), `Compare(ctx, req)`, `TimeoutFromSetting`, y `EvalComparar(in, args)` para la consulta «comparar» de las fixtures y de `motores s3 dump`.
- Generación exhaustiva crédito × plazo × fecha ± días × tienda × promoción vigente, cotizada con `motor_cred.QuotePurchase` (sin fórmulas duplicadas); el cargo normal de tarjeta siempre se genera; créditos sin límite → `excluded no_limit`.
- Restricciones duras de K16: colchón (mínimo diario de `AvailabilityTimeline` en la ventana), compras vivas (crédito usado, con la del escenario), deuda en cero (crédito usado); globales con sobreescritura campo por campo; sin valores → `missing_constraints` («faltan restricciones»), sin recomendación.
- 4 objetivos con las medidas de K17 (desviación estándar exacta con `big.Int.Sqrt`, pico, `over_cap`, unidades base por peso), normalización mín-máx solo entre válidos (empate → 1, declarado), puntaje con pesos normalizados, 4 óptimos lado a lado con desempates de §8.6 y recomendación; listas `provisional` de §14.1.
- Tabla a mano independiente (`testdata/tabla_mano.py`, fracciones exactas, día por día) y 5 fixtures con la forma de `FORMATO.md` §2 (`testdata/caso-1.json`, `caso-2.json`, `caso-3.json`, `sin-restricciones.json`, `tiempo.json`), escritas **antes** de correr el comparador; `docs/estado/S3-escenarios.md` con reglas, los 3 casos, la tabla, los tiempos y el resultado.
- Corrida: el motor reproduce la tabla a mano completa en los 3 casos (0 diferencias). Tiempos: ≤ 1 ms por caso y 111–123 ms en la carga de 3,870 candidatos; 5 s alcanza, no hay bloqueo.

## Archivos
- `backend/internal/motor_cred/comparador/request.go`, `generate.go`, `constraints.go`, `objectives.go`, `score.go`, `compare.go`, `timeout.go`, `evaluate.go`
- `backend/internal/motor_cred/comparador/comparador_test.go`, `validate_test.go`
- `backend/internal/motor_cred/comparador/testdata/tabla_mano.py`, `testdata/caso-1.json`, `caso-2.json`, `caso-3.json`, `sin-restricciones.json`, `tiempo.json`
- `docs/estado/S3-escenarios.md` (nuevo)
- `docs/estado/wf1/motor-S3-comparador-go.md` (este reporte)

## Verificación
- `gofmt -l internal/motor_cred/comparador` → vacío; `go vet ./internal/motor_cred/comparador/` → limpio.
- `go test ./internal/motor_cred/comparador/ -count=1 -cover -v | tail -40` → 13 pruebas PASS, `coverage: 96.0% of statements`.
- **S3-12:** `TestRestricciones` (3 casos = tabla a mano: `status`, `evaluated_count` 31/26/31, `valid` 17/12/17, `rejected` 14/14/14 con su restricción, 0 válidos con violación; al menos un descarte por cada restricción; cargo normal generado en las 3 fechas × 2 tiendas), `TestRestriccionesPropiedad` (300 comparaciones al azar con semilla fija, 11,730 candidatos, 0 válidos con violación; cada descartado viola la que dice y ninguna anterior), `TestRestriccionesVivas` (conteo de vivas = `motor_cred.LivePurchases` día por día), `TestRestriccionesFaltan` (4 faltantes, 0 evaluados, sin óptimos ni recomendación; una faltante; `override` que completa la global), `TestTiempoClaveSembrada` (lee `00003_puerta1b.sql`: valor sembrado 5 → 5 s; −1 → `ErrInvalid`; 0 → vence), `TestRestriccionesTiempo` (0 s → `ErrTimeout` y `ResponseFor` = 503 `unavailable`; fixture `tiempo` → código `timeout`; 5 s: `S3-TIEMPO caso-1 0–1 ms`, `caso-2 0–1 ms`, `caso-3 0–1 ms`, `carga 111–123 ms` en 3 corridas; carga = 3,870 candidatos, igual a la cuenta a mano).
- **S3-13:** `TestOptimos` (3 casos × 4 óptimos distintos = tabla a mano; el de intereses es la promoción), `TestObjetivoParejo` (0/100/200 → 81,650; quincena en 0 cuenta en n; `over_cap` por persona; empate en desviación → gana el de menos personas sobre el tope aunque tenga menos puntaje).
- **S3-14:** `TestPuntaje` (100/200/400 → 1, 2/3, 0; empate y nulos → 1; pesos por defecto Σ = 1.00; puntajes 10,000/9,167/7,500 pb; pesos {0,1,0,0} cambian la recomendación al óptimo parejo; la escala sale solo de los válidos aunque un descartado tenga intereses mayores).
- **S3-15:** `grep -rn 'net/http\|database/sql' backend/internal/motor_cred/comparador --include='*.go' | grep -v _test.go` → 0; `grep -i 'highs\|ortools\|or-tools' backend/go.mod backend/go.sum` → 0.
- **S3-17:** `grep -rn 'float32\|float64' backend/internal/motor_cred/comparador --include='*.go'` → 0 (todo en `int64`, `big.Int` y `big.Rat`).
- Validación (`TestValidacion`): 15 entradas inválidas → `ErrInvalid`; 5 argumentos inválidos de `EvalComparar` → código de error. `TestGeneracionTiendas`: tienda fija por id y por comercio, sin tiendas (objetivo 4 nulo), préstamo sin candidatos, MSI con tasa ≠ 0 no generado, `no_valid_scenario`. `TestContextoCancelado`.
- `python3 -I backend/internal/motor_cred/comparador/testdata/tabla_mano.py <salida.md>` → 3 casos con las 4 comprobaciones de §10.1 en OK (un descarte por restricción, un válido sobre el tope, 4 óptimos distintos, intereses = promoción).

## Supuestos
- `term_loan` no aporta candidatos: §8.2 no le da plazos (un préstamo no paga una compra planeada). No cambia cifras ni esquema; solicitud al contrato para escribirlo.
- Sin tienda en la petición y sin canastas: un solo candidato sin tienda (valor por peso nulo, objetivo 4 no distingue).
- Un valor por peso nulo entre valores no nulos no puede ocurrir (todos los escenarios comparten el conjunto de tiendas); si ocurriera, vale 1 como el empate.
- `AvailabilityTimeline` se llama sin el contado (`Cash` vacío): el papel de contado no mueve el crédito (S3-05).
- Las compras vivas se cuentan con los intervalos [fecha, último pago] de `QuotePurchase` (misma regla que `LivePurchases`, verificada en `TestRestriccionesVivas`) para no recotizar por día.
- Un error de `AvailabilityTimeline`, `UnifiedCalendar` o de una compra existente (estado inválido o `unsupported`) detiene `Compare` con ese error; solo los errores de cotización del candidato lo descartan sin contarlo.
- Las fixtures del comparador viven en `comparador/testdata/` porque `oraculo/fixtures/` no está en mi alcance (ver solicitudes).

## Bloqueos
- Ninguno nuevo. 5 s alcanza (peor caso 123 ms): no se abre `docs/estado/bloqueos/S3-tiempo-optimizador.md`. Se usan los abiertos `S3-objetivos-medida`, `S3-formula-tarjeta`, `S3-origen-cargo-normal` y `S3-serie-sin-intereses` como marcas PROVISIONAL en las listas `provisional`, según §14.1.

## Pendientes
- Copiar las fixtures a `oraculo/fixtures/s3/` y escribir `s3-comparador-carga` como fixture (fuera de mi alcance; ver solicitudes). La rejilla de carga ya existe como `carga()` en `comparador_test.go`.
- Que `oraculo/cred_escenarios.py` (agente de regresión) reproduzca las fixtures del comparador: lo corre ese agente.

## Solicitudes fuera de alcance
1. `oraculo/fixtures/s3/`: copiar sin cambiar un número `comparador/testdata/caso-1.json`, `caso-2.json`, `caso-3.json`, `sin-restricciones.json` y `tiempo.json` como `s3-comparador-caso-1.json`, `-caso-2`, `-caso-3`, `-sin-restricciones` y `-tiempo` (el campo `caso` ya trae ese nombre); si se quiere el `$ref` de §10.0, separar `entradas.datos` del caso 1 en `_datos/comparador-caso-1.json`. Agregar `s3-comparador-carga` con la rejilla de `carga()` (`evaluated_count` 3,870). Motivo: §10.0 se las asigna al optimizador, pero mi alcance de WF-1 solo cubre `comparador/` y `S3-escenarios.md`.
2. `backend/internal/motor_cred/evaluar/evaluar.go` (núcleo S3): en el despacho de `funcion = "comparar"` llamar `comparador.EvalComparar(input, args)` (devuelve `(Result, "")` o `(nil, código)`); con eso `motores s3 dump` incluye los casos del comparador. No edité `cmd/motores` ni `evaluar/`.
3. `docs/contratos/motores/S3-interfaz.md` §8.2: decir que `term_loan` no aporta candidatos y qué pasa sin tienda ni canastas (un candidato sin tienda), como en los supuestos.
4. Bloqueo `docs/estado/bloqueos/S3-objetivos-medida.md` (pregunta 2): agregar el hallazgo de §6 de `S3-escenarios.md` (una serie de una sola quincena tiene desviación 0 aunque sea un pago único grande).

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- **K31 punto 7, tiempo máximo del comparador:** mantener **5 s** (ya aprobado en K43 y configurable por K44). Justificación: el peor caso medido (rejilla de 3,870 candidatos, mayor que el uso real de una casa) tarda 111–123 ms en un hilo, 2.5 % del tope; los casos reales tardan ≤ 1 ms. Si S1 propone límites de seguridad para `comparator.timeout_seconds`, sugiero mínimo 1 s y máximo 30 s (con menos de 1 s un equipo cargado podría dar 503 sin necesidad; con más de 30 s la pantalla parecería colgada).
- **`optimizer.date_window_days` (sin valor sembrado; mientras falte, «faltan restricciones»):** proponer **7 días** (la misma propuesta del contrato): cubre una semana antes y después de la fecha deseada, alcanza para mover la compra a la quincena siguiente y la carga con ± 7 días tarda ~0.12 s.

## Consumo aproximado
- Acciones de herramienta: 39.
- Leído: ~155 KB (protocolo 6, decisiones C y K 17, criterios S3 36, S3-interfaz por secciones ~55, núcleo Go por partes ~25, RN-11 5, salidas de pruebas y del script ~11).
- Escrito a mano: ~118 KB (Go 54 KB con pruebas, `tabla_mano.py` 24 KB, `S3-escenarios.md` 30 KB, reportes 10 KB). Generado por el script sin pasar por el contexto: fixtures JSON 70 KB y las tablas de §4.1 (15 KB, incluidas en `S3-escenarios.md`).
