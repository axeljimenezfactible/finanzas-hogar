# correccion-S5-r1 · WF-1 · corrección del motor de precios y tiendas (S5), ronda 1

Agente `correccion-S5-r1` · 11-oct-2026. Corrige los hallazgos altos y medios de `docs/estado/wf1/REVISION-S5-r1.md` (H-01 a H-08). No toqué los bajos (H-09 a H-12) ni reabrí decisiones de §A–§K. Sin `git commit`.

## Hecho
- **H-01 (alta) · filtro de municipio de `profeco_import`.** Fui por la vía que la revisión deja como alternativa: la fixture. En `oraculo/pre_genera_fixtures.py` la cadena `walmart` ya no lleva municipio y se agrega la sucursal `walmart-tepeyac` (`kind` branch, `profeco_name` «WALMART TEPEYAC», municipio GUSTAVO A. MADERO). Así las dos lecturas de §11 («municipios de las sucursales») dan lo mismo: el motor solo toma municipios de sucursales y el oráculo los de cualquier tienda. b1 sigue cargándose a nombre de la cadena porque ninguna sucursal se llama WALMART ACUEDUCTO. Regeneré las 27 fixtures y solo cambió `s5-profeco-muestra.json` (diff contra la copia previa). No cambié `filter.go` ni el adaptador porque elegir una lectura le toca al integrador. La diferencia de fondo sigue latente: una cadena con municipio y sin sucursales filtra distinto en el motor que en el oráculo. Va como solicitud a §11.
- **H-07 (media) · `rejections[].file` y `row`.** El motor ahora usa la convención del oráculo y de la fixture. `file` es la ruta tal como llega en `files[]` (en el dump, `s5.go` traduce la ruta resuelta a la original). `row` es el ordinal del registro de datos, donde 1 es el primero después del encabezado: el campo nuevo `Record.Row` y `RowError.Row` en `reader.go`. `Line`, la línea física, se queda solo para diagnóstico. La misma convención aplica a `source_ref` («QQP <corte> <archivo>:<renglón>») y a `profeco_record.renglon` (`sqlite.go`), y por eso reconstruí `qqp.sqlite` completa. Ajusté dos pruebas: `import_test.go` (fila 3 → 2) y `sqlite_test.go` (`db.csv:2` → `db.csv:1`).
- **H-05 (media) · orden de `rows[]` en `compare_carts`.** `compare.go` aplica las convenciones 1 a 3 del oráculo (`docs/estado/wf1/oraculo-S5.md`, Supuestos). El orden es: primero los idénticos en el orden de la canasta base, luego los sustitutos por producto en orden de aparición y al final los renglones `only_one`, por canasta y renglón. Un sustituto es un producto con presentaciones sin emparejar (que no están en todas las canastas) en dos o más canastas, y `substitute_count` cuenta solo las tiendas que tienen línea en ese renglón. `stores[]` de los carritos va por total ascendente y, a igual total, en el orden base y luego `others` (`rankCarts`). `winners[]` sigue ese mismo orden.
- **H-06 (media) · orden de `stores[]` en `compare_basket`.** Ahora va por total ascendente con las incompletas incluidas, como dice §7.4 (`rankStores`). `cheapest_store` es la primera tienda completa (`basket.go`) y `difference_vs_cheapest_cents` se mide contra ella, así que en una incompleta puede ser negativa (convención 5 del oráculo). Los desempates PROVISIONALES de §7.4 no cambiaron. Cambié 4 pruebas que tomaban la incompleta por su posición; ahora la buscan por código (`storeRow`) y comprueban el orden por total de §7.4.
- **H-04 (alta) · `compara.py s5`.** Da **diferencias: 0** y código 0. Motor contra esperado: 747 de 747.
- **H-02 (alta) · cortes PROFECO 2024 y 2025.** Los descargué de `datos.profeco.gob.mx`. Vienen en RAR, igual que 2026, y no en CSV como supone §11. Cada uno quedó en su carpeta propia: `raw/2024-rar/`, `raw/2025-rar/`, `raw/2024-csv/` y `raw/2025-csv/`. Revisé su integridad (tamaños extraídos iguales al listado del RAR, archivo por archivo) y los cargué con el mismo comando, que ahora acepta varias carpetas en `-dir`. Después regeneré `docs/estado/S5-cobertura-profeco.md` (ver Verificación). No hizo falta el bloqueo `S5-profeco-cortes`: con los tres cortes cargados se cumplen tanto el plan («2026») como el criterio («2024, 2025 y 2026»).
- **H-03 (alta) · `S5-terminos.md`.** Abrí el bloqueo `docs/estado/bloqueos/S5-terminos.md` con las opciones A a D. **S5-07 no está cumplido.** No escribí `S5-terminos.md` ni entré a ningún sitio de tienda.
- **H-08 (media) · carpeta de tickets.** Moví `datos/semilla/tickets/sinteticos/` a `datos/semilla/tickets-sinteticos/` y `datos/semilla/tickets/README.md` a `datos/semilla/tickets-README.md`, junto a la carpeta, como ya se hace con `productos-README.md`. `datos/semilla/tickets/` solo tiene un `.gitkeep` oculto que `ls` no cuenta, así que hoy `ls datos/semilla/tickets/ | wc -l` = 0 y con las 5 fotos dará 5. Actualicé la ruta en `ocr/integration_test.go` y el texto del README.

## Archivos
- Modificados: `backend/internal/precios/compare.go`, `backend/internal/precios/basket.go`, `backend/internal/precios/precios_test.go`, `backend/internal/precios/eval_test.go`, `backend/internal/precios/profeco/reader.go`, `backend/internal/precios/profeco/import.go`, `backend/internal/precios/profeco/sqlite.go`, `backend/internal/precios/profeco/import_test.go`, `backend/internal/precios/profeco/sqlite_test.go`, `backend/internal/precios/profeco/cobertura/main.go`, `backend/internal/precios/ocr/integration_test.go`, `backend/cmd/motores/s5.go`, `oraculo/pre_genera_fixtures.py`, `oraculo/fixtures/s5/s5-profeco-muestra.json` (regenerada), `docs/estado/S5-cobertura-profeco.md` (versión 1.1).
- Creados: `docs/estado/bloqueos/S5-terminos.md`, `datos/semilla/tickets/.gitkeep`, este reporte.
- Movidos: `datos/semilla/tickets/README.md` → `datos/semilla/tickets-README.md`; `datos/semilla/tickets/sinteticos/` → `datos/semilla/tickets-sinteticos/`.
- Datos (ignorados por git): `datos/semilla/profeco/raw/2024-rar/QQP_2024.rar`, `raw/2025-rar/QQP_2025.rar`, `raw/2024-csv/QQP_2024/` (24 CSV), `raw/2025-csv/QQP_2025/` (24 CSV); `datos/semilla/profeco/db/qqp.sqlite` reconstruida con los tres cortes (16 GB).

## Verificación
- `python3 -I oraculo/compara.py s5` (con `motores s5 dump --fixtures oraculo/fixtures/s5`): oráculo contra esperado 747/747; **motor contra esperado 747/747**; oráculo contra motor 747/747; provisionales 225/225; **diferencias: 0**, código 0 (H-04). `python3 -I oraculo/pre_run_fixtures.py`: 27 fixtures, 0 diferencias.
- `go vet ./internal/precios/... ./cmd/motores/`: limpio. `gofmt -l internal/precios cmd/motores`: vacío.
- `go test -coverprofile ./internal/precios/...`: precios 96.1 %, asistida 100 %, ocr 99.3 %, profeco 95.3 % (los cuatro ≥ 85 %, K25). El total del perfil incluye los dos `main` sin pruebas (`ocrmedir` y `cobertura`). Todas las pruebas PASS; `-run` por criterio: Comparativa 4, Unidad 2, MejorTienda 5, Frescura 3, Fuente 4, Membresia 2, Mercado 2, PreciosViejos 1 y LecturaAsistida 1 PASS.
- S5-01 / S5-03: canon sin cambios (+13850 y +10050; 61 / 60 / 55; 22 / 8 / 3) y `rows[]` y `stores[]` en el orden del oráculo (`canon-filas` y `s5-mejor-tienda-incompleta` sin diferencias).
- S5-05: `wc -l` 2024 14,889,243 − 24 = 14,889,219; 2025 13,270,548 − 24 = 13,270,524; 2026 9,388,473 − 16 = 9,388,457. `sqlite3 qqp.sqlite "SELECT COUNT(*), MIN(fecha), MAX(fecha), MAX(fecha_registro) FROM profeco_record"` → `37548200|2024-01-02|2026-08-31|2026/08/31`; por año 14889219 · 13270524 · 9388457. Hay 0 rechazados y 0 diferencias por archivo en los 64. `grep -c cadena_comercial docs/estado/S5-cobertura-profeco.md` = 13 y el documento cita `QQP_diccionario.csv`. Con los tres cortes, Bodega Aurrera Vallejo tiene 23,809 registros (todos solo en la dirección de Bodega Aurrera Express Sucursal Poniente 140). Walmart Acueducto y Sam's: ver la sección 6, «Sucursales de Axel».
- S5-06: `ls datos/semilla/tickets/ | wc -l` = 0. `go run ./internal/precios/ocr/cmd/ocrmedir` → «S5-06 sin medir: hay 0 imágenes y se necesitan 5: faltan las fotos de Axel». `TestIntegracionSinteticos` PASS con la ruta nueva.
- S5-07: **NO CUMPLE**, pendiente del bloqueo `S5-terminos`. `-run LecturaAsistida` PASS: rechaza la activación sin veredicto. El grep de dominios de tiendas en Go sin `_test.go` da 0 líneas.
- S5-08: el grep de `net/http|http.Get|http.Client` en `internal/precios` sin `_test.go` da 0 líneas. URLs consultadas: `https://datos.profeco.gob.mx/datos_abiertos/qqp.php`, `https://datos.profeco.gob.mx/` y los dos `file.php?t=…` de 2024 y 2025 (solo la fuente del plan).

## Supuestos
- El orden de `rows[]`, `stores[]` y `winners[]` y la definición de sustituto de las convenciones 1 a 3 del oráculo se toman como la convención vigente mientras §15.2 no la escriba. No cambia ninguna cifra canon.
- `rejections[].file` = ruta tal como llega, y `row`/`renglon`/`source_ref` = ordinal del registro de datos. Es la convención del oráculo mientras §11 y §15.2 no la escriban.

## Bloqueos
- `docs/estado/bloqueos/S5-terminos.md` (nuevo; S5-07). Siguen los anteriores: `S5-reglas-comparativa`, `S5-carrito-comparativa` y `S5-profeco-almacen`.

## Pendientes
- S5-07: depende de que Axel responda el bloqueo `S5-terminos`.
- S5-06: faltan las 5 fotos de Axel en `datos/semilla/tickets/` y sus transcripciones en `datos/semilla/tickets-esperado/`.
- Diferencias latentes entre motor y oráculo que ninguna fixture prueba y que esta corrección no cubre. Se señalan para §7.2 y §11. (1) Cadena con municipio y sin sucursales (H-01, de fondo). (2) Ganador de un sustituto cuando una tienda tiene varias presentaciones: el motor usa costo total ÷ contenido total y el oráculo la mejor presentación. (3) En `compare_carts`, `missing_count` (el motor cuenta productos ausentes y el oráculo da 0) y `stale_count` (el motor cuenta por línea de renglón y el oráculo por presentación agrupada). (4) Filtro PROFECO: el motor compara con `Fold` (sin acentos ni puntuación) y el oráculo por igualdad exacta.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S5-interfaz.md` §15.2: escribir el orden de `rows[]` (idénticos en orden de la canasta base, sustitutos por producto, `only_one` por canasta y renglón), de `stores[]` de `compare_carts` (total y luego orden base/others), de `winners[]` (mismo orden, con `count` 0) y la definición de sustituto (presentaciones sin emparejar en dos o más canastas). Motivo: H-05; hoy solo lo dicen el oráculo y el motor.
- S5-interfaz §7.4: decir que `difference_vs_cheapest_cents` se mide contra la más barata completa (puede ser negativa en una incompleta) y que `cheapest_store` es la primera completa. Motivo: H-06.
- S5-interfaz §11: (a) fijar si la `municipality` de una tienda `chain` entra al filtro (H-01; la fixture ya no depende de eso); (b) `rejections[].file` = ruta tal como llega en `files[]` y `row`/`renglón` = ordinal del registro de datos, donde 1 es el primero tras el encabezado (H-07); (c) los cortes 2024 y 2025 también vienen en RAR; (d) las rutas reales `raw/<corte>-rar/`, `raw/<corte>-csv/` y `datos/semilla/profeco/db/qqp.sqlite` (H-11).
- `docs/criterios/S5.md` S5-05: «2024, 2025 y 2026 (este último en RAR)» → los tres vienen en RAR. Plan §3.4 fila S5: «PROFECO 2026 cargado» → 2024 a 2026 (ya se cumplen los dos).
- S5-interfaz §12: mencionar `datos/semilla/tickets-README.md` y `datos/semilla/tickets-sinteticos/` como lugar de las instrucciones y de las imágenes sintéticas.

## Paquetes instalados
Ninguno (`bsdtar` y `sqlite3` ya estaban).

## Propuestas para Axel
Ninguna de K31. Decidir el bloqueo `S5-terminos` (opción A, B, C o D).

## Consumo aproximado
52 acciones de herramienta; ~180 KB leídos; ~125 KB escritos (de ellos ~59 KB de la sección 6, que genera el cargador).
