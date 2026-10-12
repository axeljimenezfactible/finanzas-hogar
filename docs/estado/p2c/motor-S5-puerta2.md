# P2c · motor-S5-puerta2 (Opus) · 11-oct-2026, noche

Alcance: `backend/internal/precios/` (con `profeco/`), `backend/cmd/motores/s5.go`, `oraculo/pre_*.py`, `oraculo/fixtures/s5/` y este reporte. Contrato: `docs/contratos/motores/S5-interfaz.md` 1.5. No toqué `ocr/` (ninguna regla lo exige), ni `docs/design/`, ni `docs/estado/wfp3/`. No hice `git add` ni `git commit`.

## Hecho

**HECHO 1 · Un solo significado de la mejor tienda (§7.4; H-r3-04, Q2-05 = A; M-r5-3, Q2-08 = A).**
- Motor: `rankCarts` y `rankStores` usan ahora la misma función `rankComplete` (`compare.go`): primero las tiendas completas (`missing_count` 0) por total ascendente y después las incompletas por total ascendente. A igual total, `compare_carts` pone primero la base y luego `others`; `compare_basket` usa más renglones recomendados, luego menos precios viejos y luego el orden de `stores` (FIRME, Q2-60 = A). `StoreTotals.DifferenceVsCheapestCents` pasa a `*int64`: vale total − total de la primera completa y es `null` en las incompletas, o en todas si ninguna está completa. `CompareBasket.cheapest_store` = primera completa o `null`.
- Oráculo: `_diferencias` en `pre_calculo.py` aplica la misma regla en `compare_carts` y en `compare_basket`.
- Canasta vacía: por la letra de §7.4, todas las tiendas tienen `missing_count` 0, así que la más barata es la primera. Antes el motor daba `null` (protección `len(list) > 0`) y el oráculo daba la primera. Alineé el motor con el texto del contrato (PREGUNTA P-1).

**HECHO 2 · Granel (§7.2 paso 8; PROVISIONAL, bloqueo `S5-carritos-granel`, regla en uso = la del motor).** El motor ya hacía esto. El oráculo deja de redondear por línea: `_agrupa` guarda solo el importe exacto, `_unit` = importe exacto ÷ cantidad, y el total y `identical_total_cents` de cada tienda son HALF_UP de la suma exacta, con un solo redondeo. `lines[].total_cents` es el HALF_UP de su línea y solo se muestra.

**HECHO 3 · Lectura PROFECO (§11; regla técnica, H-r3-03, Q2-04 = A).**
- Motor (`profeco/reader.go`, `profeco/parse.go`):
  - «CSV mal formado» sale sin el detalle del lector.
  - Gramática del precio de §11: una coma en los decimales («1.0,0») ahora es «precio con comas mal puestas»; antes se aceptaba como 1.00.
  - «más de dos decimales» se revisa antes que «ilegible», así que «.123» da «más de dos decimales», en el orden que lista el contrato.
- Oráculo (`pre_profeco.py`):
  - Lee cada registro con la lista cerrada y en el orden de §11: «CSV mal formado», «columnas: n en vez de m», «producto vacío», «cadena comercial vacía» (se rechaza, ya no cuenta como filtrada), gramática del precio (acepta «1,021.00» = 102100 y rechaza 0.00 y la cadena vacía), «fecha vacía», «fecha ilegible» y «fecha futura». La fecha futura antes no se revisaba.
  - Liga normalizada (PROVISIONAL, `S5-profeco-filtro` pregunta 2 = Q2-57): ATUN DOLORES = ATÚN DOLORES. Una liga repetida después de normalizar da `invalid_argument`, igual que el motor.
  - Emite `loaded_proposed`.
- `s5 dump` (`cmd/motores/s5.go`): `loaded` cuenta ahora solo las ligas confirmadas. Antes era `Loaded + LoadedProposed`. Así se cumple la identidad de §11: Read = Rejected + FilteredOut + Skipped + Loaded + LoadedProposed + Duplicates. `coverage` sigue contando las dos.

**HECHO 4 · Orden de las listas de salida (§15.2).**
- `current_price.references[]` va por `observed_date` descendente y luego por `source` en orden de bytes, en el motor (`pricing.go`; antes solo por `source`) y en el oráculo.
- `stale_prices_for_plan` sigue el orden de `products[]` y quita los repetidos en los dos.
- Revisé el resto de listas contra el código: `stores[]` (HECHO 1), `winners[]` (en el orden de `stores[]`, con 0), `rows[]` y `lines[]`, `items[]` y `prices[]`, `coverage[]` (orden de bytes de cadena y municipio) y `rejections[]` (por archivo y renglón). Coinciden en el motor y en el oráculo.

**HECHO 5 · Fixtures (generadas por `oraculo/pre_genera_fixtures.py`; el `esperado` es aritmética escrita ahí con `Fraction`, sin llamar al oráculo).**
- Nueva `s5-carritos-granel` (PROVISIONAL `S5-carritos-granel`), con los casos del revisor rehechos:
  - `granel`: totales 101 y 103 (con redondeo por línea daban 102 y 104); líneas mostradas de 51 y 52; precio por unidad base 135, 137, 721 y 736.
  - `granel-sustituto`: 0.333 × 750 g a 22.00 contra 0.667 × 1 kg a 28.00. Da 2933 y 2800, gana Walmart y los totales son 733 y 1868.
- Nueva `s5-profeco-bordes` (PROVISIONAL `S5-profeco-filtro`) con `_datos/profeco-bordes.csv`, 17 registros:
  - Producción: rechazados 4, filtrados 1, omitidos 1, cargados 10, propuestos 0 y duplicados 1.
  - Revisión: 4, 1, 0, 10, 1 y 1.
  - Los motivos se comparan.
- Nueva `s5-profeco-motivos` con `_datos/profeco-motivos.csv`: el resto de la lista cerrada (precio vacío, dos casos de comas, dos de ilegible, producto vacío, fecha vacía y fecha futura) más un renglón de control que carga.
- `s5-carritos-presentaciones` regenerada: `stores` = bodega (completa, diferencia 0), tercera (null) y walmart (null); `winners` en ese orden.
- `s5-mejor-tienda-incompleta`: t2, t1 y al final t3 con `null`.
- `s5-profeco-muestra`:
  - Ahora compara los motivos.
  - En revisión a6 pasa a `loaded_proposed`: cargados 2 y propuestos 1.
  - `provisional` = `S5-profeco-filtro`, como pide la fila de §15.
- `s5-profeco-filtro` lleva `loaded_proposed` 0.
- `provisional` en `null` en 7 fixtures de `S5-reglas-comparativa` (FIRME, Q2-59 a Q2-63): `s5-canasta-regla-unica`, `s5-cotizacion-recomendada`, `s5-membresia-anual`, `s5-promocion-vigencia` y las 3 `s5-mejor-tienda-empate-*`. `s5-mejor-tienda-incompleta` también queda en `null`.
- `oraculo/salidas/s5.json` regenerada con `pre_run_fixtures.py`: 32 fixtures (antes 29).

**HECHO 6 · Marcas y pruebas Go.**
- Comentarios que decían PROVISIONAL y ahora son FIRME, cada uno con su Q2: `basket.go` (Q2-59, Q2-60, Q2-47), `extras.go` (Q2-62 = B, Q2-63), `pricing.go` (Q2-61), `profeco/sqlite.go` e `import.go` (Q2-53, Q2-55).
- La nota de membresía sin costo dice «sin costo capturado» (§8). La clave `note` no está en §15.2 y no se compara.
- Pruebas ajustadas: `TestEvalCarritos`, `TestEvalCarritosVariasPresentaciones` y `TestMejorTiendaCanasta` (orden y `null`).
- Pruebas nuevas: `TestEvalCarritosGranel` y `TestParsePriceMotivos` (los textos exactos de la lista cerrada).

**HECHO 7 · Cifras de `docs/estado/S5-cobertura-profeco.md`: no cambian.** El único cambio de lectura que podía mover un conteo real es la coma en los decimales. `sqlite3 -readonly qqp.sqlite` da 37,548,200 registros y 0 precios con coma, así que no toqué el archivo.

**HECHO 8 · Discrepancia con el resumen de `contratos-puerta2.md`.** El resumen dice «`tercera` va al final con `null`». Por §7.4, Walmart también está incompleta (no compró leche, `missing_count` 1), así que el orden es bodega, tercera (165.00) y walmart (296.50), y las dos incompletas llevan `null`. Seguí el texto de §7.4, que es el que manda.

## Verificación
Todo desde la raíz del repositorio (`/home/deepsec/finanzas-hogar`):
```
$ python3 -I oraculo/pre_genera_fixtures.py
32 fixtures en /home/deepsec/finanzas-hogar/oraculo/fixtures/s5
$ python3 -I oraculo/pre_run_fixtures.py --out oraculo/salidas/s5.json
32 fixtures, 32 coinciden, 1336 hojas de esperado, 0 diferencias
$ python3 -I oraculo/pre_run_fixtures.py --out datos/tmp/oraculo-s5.json && cmp datos/tmp/oraculo-s5.json oraculo/salidas/s5.json
salidas/s5.json idéntica
$ go -C backend build -o ../datos/tmp/motores ./cmd/motores && datos/tmp/motores s5 dump --fixtures oraculo/fixtures/s5 --out datos/tmp/motor-s5.json
$ python3 -I oraculo/compara.py s5 --motor datos/tmp/motor-s5.json
oráculo contra esperado: 935 de 935 hojas coinciden
motor contra esperado: 935 de 935 hojas coinciden
oráculo contra motor: 935 de 935 hojas coinciden
provisionales: 315 de 315 hojas coinciden
diferencias: 0
$ python3 -I oraculo/pre_verify.py | tail -1
TODO COINCIDE
$ (cd backend && gofmt -l ./internal/precios ./cmd/motores && go vet ./internal/precios/... ./cmd/motores)
(sin salida)
$ (cd backend && go test ./internal/precios/... -count=1 -cover)
ok  finanzas-hogar/internal/precios          coverage: 95.9% of statements
ok  finanzas-hogar/internal/precios/asistida coverage: 100.0% of statements
ok  finanzas-hogar/internal/precios/ocr      coverage: 99.3% of statements
ok  finanzas-hogar/internal/precios/profeco  coverage: 95.3% of statements
$ (cd backend && go test ./... -count=1 -cover)   # todos ok; sin FAIL
$ sha256sum oraculo/model.json oraculo/verify3.json datos/semilla/data.json
1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4  oraculo/model.json     (= MODEL_SHA del Makefile)
d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041  oraculo/verify3.json   (= VERIFY3_SHA)
2fe7afae5b1be77fafeed787d4c7bd4d68fe469b78a86898b370c1ba1031d31b  datos/semilla/data.json (git diff HEAD: sin cambios)
```
Además comparé la salida completa del oráculo contra la del motor, hoja por hoja, también las hojas fuera de `esperado`. Hay una sola diferencia, que ya existía: la clave `note` de `membership` sin costo, que solo emite el motor y que §15.2 no lista.

## Supuestos
- S-1. «Se quitan los espacios y un `$` inicial» (§11 paso 3): quito los espacios de las orillas, después un `$` y después otra vez los espacios. Un espacio en medio («1 021.00») da «precio ilegible». Es la lectura del motor en uso, y el oráculo hace lo mismo (PREGUNTA P-3).
- S-2. Más de 15 dígitos enteros da «precio ilegible»: es la protección del motor contra el desborde de int64, y el oráculo la copia. Usa un motivo de la lista cerrada; el contrato no la menciona.
- S-3. «Más de dos decimales» cuenta los caracteres después del primer punto, aunque no sean dígitos («1.2.3» da «más de dos decimales»). Así lo hace el motor y el oráculo lo sigue.
- S-4. En `s5-profeco-bordes` el renglón 17 es un duplicado. El contrato fija rechazados 4, filtrados 1, omitidos 1 y cargados 10, y el revisor no dejó su CSV (lo borró), así que 17 = 16 + 1 duplicado.

## Bloqueos
No abrí bloqueos nuevos: `docs/estado/bloqueos/` no es de mi alcance. Siguen abiertos y citados en las fixtures: `S5-carritos-granel`, `S5-profeco-filtro` (Q2-56 a Q2-58), `S5-carritos-detalle` (Q2-50 a Q2-52) y `S5-carrito-comparativa` pregunta 2 (Q2-48).

## Pendientes
Diferencias entre motor y oráculo que ya existían, que ninguna fixture prueba y que el contrato no fija. No las decidí:
- PREGUNTA P-1: con canasta vacía, ¿`cheapest_store` es la primera tienda (la letra de §7.4, que es lo que quedó) o `null`?
- PREGUNTA P-2: `current_price.references[]`: el motor da una referencia por fuente (la más reciente) y el oráculo da todas. §6 no lo dice.
- PREGUNTA P-3: los espacios dentro del precio (S-1).
- «CSV mal formado» no se alcanza con el lector en uso (`LazyQuotes = true`), así que no tiene fixture.
- El motor guarda como máximo 100 `rejections[]` con detalle y el oráculo no tiene tope.
- Los formatos de fecha aceptados no son iguales:
  - Solo el motor acepta `2006/1/2`, `2/1/2006` y `20060102`.
  - Solo el oráculo acepta `HH:MM` sin segundos y fracciones de segundo.
- Una liga con `presentacion_id` desconocido: el motor da `unknown_presentation` y el oráculo omite el registro.
- Misma presentación en varios renglones de una canasta con precios distintos: el motor consume en orden de captura y el oráculo usa el precio promedio. §7.2 no lo define.
- H-r3-06 (baja): las etiquetas «(corte 2026)» de `S5-cobertura-profeco.md` §1 y §4.2 no son cifras que cambien, así que quedan fuera de mi alcance.

## Propuestas para Axel
- Contestar `S5-carritos-granel`: ¿el importe de un renglón a granel se redondea por línea, como en un ticket, o se queda la regla en uso (un solo redondeo)?
- Contestar `S5-profeco-filtro` (Q2-56 a Q2-58), `S5-carritos-detalle` (Q2-50 a Q2-52) y Q2-48 (`S5-carrito-comparativa` pregunta 2).

## Consumo aproximado
Unas 50 acciones de herramienta · unos 230 KB leídos · unos 150 KB escritos (incluye unos 95 KB de JSON generado: 3 fixtures nuevas, 10 regeneradas y `salidas/s5.json`).
