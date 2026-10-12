# motor-S5-go · Reporte final (WF-1, reanudación del 11-oct-2026)

Agente `motor-S5-go` (Opus). Contrato: `docs/contratos/motores/S5-interfaz.md` **1.4** (ronda 4: M-r4-1 una sola regla para el total de cada tienda; M-r4-2 cantidad a comprar con `quantity_milli`). Criterios: `docs/criterios/S5.md` 1.2. Los parciales del intento anterior (contrato 1.2) se leyeron, se compilaron y se reescribieron donde no cumplían la 1.4.

## Hecho
- **Núcleo `backend/internal/precios`** (sin red ni base de datos; centavos `int64`, fracciones exactas con `math/big`, HALF_UP solo al publicar):
  - Ingesta por fuente con la validación común de §3 (centavos ≥ 0, fecha no futura, presentación y tienda existentes, llave única → `ErrDuplicatePrice`); `manual`, `ticket_manual` y `ticket_ocr` compiten con el mismo peso; `site_assisted` y `profeco` solo son referencia (K22, K37).
  - Precio por unidad base por 1 kg, 1 l o 1 pieza con HALF_UP y fracción exacta (K7, RN-20); idéntico = misma presentación, sustituto = otra del mismo producto.
  - Frescura con los parámetros de `app_setting` (K44): viejo a los `Dv` días; peso `(D0 − edad)/(D0 − Dv + 1)` = (90 − días)/61 con 30 y 90 (K40); después de `D0` no compite.
  - Precio vigente (§6): último regular que compite (a igual fecha, `created_at` más reciente), referencias aparte, promoción vigente solo en su fecha (PROVISIONAL `S5-reglas-comparativa` p3).
  - `QuoteFor` con la **regla única de §14.1** (unidad base ajustada ÷ peso → menor precio del paquete → orden de `stores` → orden del catálogo; `recommended` solo sin tienda ni presentación).
  - `CompareCarts` (§7.2, canon): consume los renglones en orden de captura (la sexta lata de atún queda como «solo en Bodega»), ganadores y empates solo en idénticos, `piece_count` e `identical_pieces`.
  - `CompareBasket` (§7.3, 1.4): cantidad a comprar = `quantity_milli` o `needed_milli`; total de **cada** tienda con el menor costo de paquetes completos (también la recomendada); recomendación por renglón = `QuoteFor` sin elección; empates de total por renglones recomendados, luego menos viejos, luego orden (PROVISIONAL `S5-reglas-comparativa` p2); una tienda con faltantes no puede ser la más barata.
  - `CompareExpense` (PROVISIONAL `S5-carrito-comparativa` p1, opción A), `Membership` (RN-13: ahorro por renglón donde Sam's es más barato, equilibrio hacia arriba p4b, compras al año p4a opción B, `missing_fee` sin costo), `NextVisit` (RN-16), `StalePricesForPlan` (§10, `QuoteFor` en la fecha de despensa), `Summary` (`PriceSummary`).
  - `eval.go`: `Input` de §15.1 con `unknown_field`, `missing_field`, `invalid_date`, rutas `*_file` validadas (`RepoPath`), `report_file` que arma las canastas `bodega` (expense) y `walmart` (cart) desde `data.json`, y las funciones consultables de §15.2.
- **`backend/internal/precios/asistida`**: interruptor `CanEnable` (global + tienda + veredicto completo por sitio) y lector de `S5-terminos.md`; sin cliente HTTP ni URLs de tiendas.
- **`backend/cmd/motores/s5.go`**: registro `s5` para `motores s5 dump`; `profeco_import` y `ocr_parse_tsv` se atienden con los subpaquetes `profeco` y `ocr` (sin tocarlos), solo con archivos de `oraculo/fixtures/s5/_datos/`.
- Pruebas con los nombres de §16 y cobertura 96.0 % (`asistida` 100 %).

## Archivos
- `backend/internal/precios/types.go`, `pricing.go`, `seed.go`, `extras.go` (parciales revisados y ajustados a la 1.4).
- `backend/internal/precios/compare.go`, `basket.go` (reescritos), `eval.go` (nuevo).
- `backend/internal/precios/precios_test.go`, `eval_test.go` (nuevos).
- `backend/internal/precios/asistida/asistida.go` (parcial revisado, sin cambios), `asistida_test.go` (nuevo).
- `backend/cmd/motores/s5.go` (nuevo).
- `docs/estado/wf1/motor-S5-go.md` (este reporte).

## Verificación
- `gofmt -l internal/precios/*.go internal/precios/asistida cmd/motores/s5.go` → vacío. `go vet ./internal/precios/ ./internal/precios/asistida/ ./cmd/motores/` → limpio.
- `go test ./internal/precios/ -count=1 -cover -v | tail -40` → 31 pruebas PASS, `coverage: 96.0% of statements`. `go test ./internal/precios/asistida/ -cover` → PASS, 100 %. `go test ./cmd/motores/` → ok.
- `sha256sum datos/semilla/data.json` → `2fe7afae5b1be77fafeed787d4c7bd4d68fe469b78a86898b370c1ba1031d31b` antes y después (sin cambio).
- Humo del CLI (fixtures temporales en `datos/tmp/`, ya borradas): `datos/tmp/motores s5 dump --fixtures … --out datos/tmp/motor-s5.json` → código 0; canon 309450 / 323300 / +13850, 277350 / 287400, 61 / 60 / 55, `ties` 8; `ocr_parse_tsv` y `profeco_import` (modo `review` y `production`) responden con conteos que cuadran (`read` 3 = 1 rechazado + 1 filtrado + 1 cargado o saltado por liga sin confirmar).
- Criterios:
  - **S5-01** `go test -run Comparativa` PASS (lee `datos/semilla/data.json` y `productos.json`): K-01 309450, K-02 323300, K-03 +13850 = +10050 + 5900 − 2100, K-04 277350, K-05 287400, K-06 +10050, K-07 61 / 60 / 55, K-08 33 / 3, ganadores 22 / 8 / 3, todo contra los valores leídos de `data.json`; 0 diferencias.
  - **S5-02** `-run Unidad` PASS: arroz 2933 vs 2800 → Walmart; crema 12692 vs 16406 → Bodega; leche 3375 vs 2833 → Walmart; HALF_UP en g, ml y pieza.
  - **S5-03** `-run MejorTienda` PASS: canasta del 10-oct → Bodega; `s5-cotizacion-recomendada` (Walmart 1 kg recomendado, 2800; totales Bodega 1500, Walmart 2200, +700); `s5-canasta-regla-unica` (Walmart 9430, Bodega 9460, +30; aceite 2 paquetes 7980); 2 canastas con empate; mercados y Sam's con membresía supuesta.
  - **S5-04** `-run Frescura` y `-run Fuente` PASS: 29 → 1 sin marca; 30 → 60/61 (0.9836); 31 → 59/61; 60 → 30/61 (0.4918); 90 → 0 no compite; 91 marcado y no compite; recta de 1/61 exacta; tres fuentes que compiten pesan igual; PROFECO y lectura asistida más baratos nunca ganan.
  - **S5-07** `go test ./internal/precios/asistida/ -run LecturaAsistida` PASS; grep de dominios de tiendas en `backend` (sin `_test.go`) → 0 líneas.
  - **S5-08** grep `net/http|http.Get|http.Client` en `internal/precios` (sin `_test.go`) → 0 archivos.
  - **S5-09** `-run Membresia` PASS: 100000 ÷ 12500 = 8; 100000 ÷ 12000 → 9; sin costo → `missing_fee`.
  - **S5-10** `-run Fuente` PASS: 5 fuentes × 1 precio con su marca de compite o referencia; promoción vencida no entra.
  - **S5-11** `-run Mercado` PASS: sábado 10-oct + martes → 2026-10-13; mercado sin el grupo no compite.
  - **S5-12** `-run PreciosViejos` PASS: despensa 2026-11-09 → precio del 10-oct viejo (30 días, 60/61).
  - **S5-13** vet y gofmt limpios, cobertura 96.0 %, grep `float32|float64` en `internal/precios` (sin `_test.go`, sin `profeco/` ni `ocr/`) → 0 líneas.
  - S5-05 y S5-06 son de los agentes PROFECO y OCR; este agente solo los expone en el CLI.

## Supuestos
- S5-go-1: en `shopping_list`, un renglón con `recommended_presentation` nulo no marca sustituto en ninguna tienda (no hay presentación contra cuál marcarlo). No toca cifras canon.
- S5-go-2: «últimos 365 días» de `grocery_purchases` = edad de 0 a 364 días contra «hoy», solo tiendas `chain` o `branch` (incluye Sam's).
- S5-go-3: `CompareExpense` no cotiza una tienda sin ningún precio que compita para los productos de la compra; la fecha y la fuente de un carrito cotizado son las de su precio más antiguo. Sin fixture; PROVISIONAL por `S5-carrito-comparativa`.
- S5-go-4: `compare_basket` con lista vacía devuelve `cheapest_store` nulo.
- S5-go-5: en `membership`, «la mejor tienda sin membresía» = tiendas con `requires_membership = false`; el costo por renglón es el mismo menor costo de paquetes del total de la canasta.
- S5-go-6: `Summary` marca `fresh` si algún precio que compite no es viejo (el contrato no cubre el caso en que el mejor es viejo y otro es fresco).
- S5-go-7: `ErrUnknownProduct` no está en §17; sale como `invalid_argument`.
- S5-go-8 (solo adaptador del CLI): en `profeco_import`, `loaded` = ligas confirmadas + propuestas en modo `review` (§11: «se usan las dos») y se agrega la clave `loaded_proposed`; en `production` solo cuentan las confirmadas. `presentacion_id` del CSV de ligas de fixture puede traer el id numérico o el código de la presentación; los ids del importador son posiciones desde 1 (tienda en `stores`, presentación en el catálogo).

## Bloqueos
Ninguno nuevo. Reglas PROVISIONALES que el código aplica con la opción indicada en el contrato: `docs/estado/bloqueos/S5-reglas-comparativa.md` (p1 peso al ordenar, p2 desempates, p3 promociones, p4a/p4b membresía), `docs/estado/bloqueos/S5-carrito-comparativa.md` (`CompareExpense`, fuente `manual` de la semilla), `docs/estado/bloqueos/S4-asignacion-mercado.md` (citado en el contrato; el motor no mueve artículos entre grupos).

## Pendientes
- Fixtures `oraculo/fixtures/s5/*.json` y `oraculo/fixtures/s5/_datos/` y el oráculo `oraculo/s5_run_fixtures.py` no existen; sin ellos no corre `python3 -I oraculo/compara.py s5`. Los escribe otro agente (§15.3).
- Capa de datos (lectura de `price`, `store`, `store_product_group`, `app_setting price.*`) y traducción a C1 (`BasketComparison`, `PriceAt`, `PriceSummary`): fuera de WF-1 para este motor.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S5-interfaz.md` §15.2 (`profeco_import`): fijar si `loaded` suma las ligas propuestas en modo `review` y si `presentacion_id` del CSV de ligas de fixture es id o código. Motivo: el importador del agente PROFECO usa id entero y las fixtures nombran entidades por código.
- `docs/contratos/motores/S5-interfaz.md` §7.3 y §15.1: decir si un renglón con `recommended_presentation` nulo cuenta como idéntico o sustituto (S5-go-1).
- `docs/contratos/motores/S5-interfaz.md` §6: cerrar `PriceSummary` cuando el mejor precio que compite es viejo y otro es fresco (S5-go-6).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Rango mínimo de los parámetros de frescura (K44, para la puerta 2): `price.stale_after_days` ≥ 1 y `price.zero_weight_after_days` ≥ `price.stale_after_days`. Justificación: la recta (D0 − edad)/(D0 − Dv + 1) solo baja de 1 a 0 si D0 ≥ Dv; con Dv = 0 todo precio nacería viejo. El motor rechaza otros valores con `invalid_argument`.

## Consumo aproximado
~41 acciones de herramienta; ~190 KB leídos; ~110 KB escritos.
