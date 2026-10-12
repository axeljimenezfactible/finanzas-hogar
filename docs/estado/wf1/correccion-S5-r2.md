# correccion-S5-r2 · WF-1 · corrección del motor S5 tras `REVISION-S5-r2.md`

Agente `correccion-S5-r2` · 11-oct-2026. Corrige los hallazgos alto y medios de `docs/estado/wf1/REVISION-S5-r2.md` dentro de su alcance; lo que toca contrato, criterios o plan va como solicitud.

## Hecho
- **R2-01 (alta, S5-07):** no se corrige aquí ni puede corregirlo ningún agente. `docs/estado/S5-terminos.md` sigue sin existir a propósito: escribir veredictos exigiría leer términos de sitios de tiendas, que WF-1 no puede (protocolo §7). El bloqueo `docs/estado/bloqueos/S5-terminos.md` (opciones A–D) sigue abierto para la puerta 2. Siguen cumpliendo `-run LecturaAsistida` (rechaza la activación sin veredicto) y el grep de dominios de tiendas (0).
- **R2-02 (media, filtro de `profeco_import`):** el §11 no fija la lectura, así que abrí el bloqueo `docs/estado/bloqueos/S5-profeco-filtro.md` con 3 preguntas y 3 opciones. La opción A es PROVISIONAL: (1) solo cuentan los municipios de las sucursales, como dice el texto de §11 (la cadena y el mercado no amplían el filtro); (2) cadena, municipio y nombre comercial se comparan normalizados (minúsculas, sin acentos, sin guion, apóstrofo ni punto, un solo espacio); (3) la sucursal del precio se busca dentro de la cadena que casó. El motor ya hacía A (`filter.go`, sin cambios). Alineé el oráculo: `oraculo/pre_profeco.py` (`_normaliza`, municipios solo de `branch`, sucursal dentro de su cadena) y `oraculo/pre_datos.py` (la tienda guarda `parent`). Fixture nueva `s5-profeco-filtro`, marcada PROVISIONAL, con los dos casos de la revisión y uno más: f3 es una cadena con municipio sin sucursal, f4 un mercado con municipio, f1 y f2 son «Wal-mart» contra «WALMART» y «WAL-MART», f7 tiene acentos y minúsculas, y f6 y f8 prueban que la sucursal se busca en la otra cadena y el duplicado que resulta.
- **R2-03 (media, `compare_carts`):** abrí el bloqueo `docs/estado/bloqueos/S5-carritos-detalle.md` con 3 preguntas y 3 opciones. La opción A es PROVISIONAL: (1) el ganador de un sustituto con varias presentaciones en una tienda se decide por importe exacto total ÷ contenido total; (2) `missing_count` = productos de la unión de las canastas que faltan por completo en esa tienda; (3) `stale_count` = renglones agrupados por presentación si la canasta es vieja. Alineé el oráculo (`oraculo/pre_calculo.py`: ganador, `missing_count`, y `winners[]` en el orden de `stores[]`, que es su propia convención 1 y no cumplía) y el motor (`backend/internal/precios/compare.go`: `stale_count` contaba líneas de `rows[]`, así que el atún con excedente contaba dos veces). Agregué la fixture nueva `s5-carritos-presentaciones` (PROVISIONAL), con esas claves en el esperado y una tienda con 750 g a 15.00 y 1 kg a 50.00, y la prueba Go `TestEvalCarritosVariasPresentaciones`.
- **R2-04 (media, contrato desfasado):** está fuera de mi alcance (contrato, criterio y plan). Abajo van las solicitudes con el texto exacto.
- Regeneré las fixtures: las 27 anteriores quedan idénticas byte a byte (`diff -rq`) y las 2 nuevas suman 29. La salida del oráculo de los 27 casos anteriores no cambia.

## Archivos
- Creados: `docs/estado/bloqueos/S5-profeco-filtro.md`, `docs/estado/bloqueos/S5-carritos-detalle.md`, `oraculo/fixtures/s5/s5-profeco-filtro.json`, `oraculo/fixtures/s5/s5-carritos-presentaciones.json`, `oraculo/fixtures/s5/_datos/profeco-filtro.csv`, `docs/estado/wf1/correccion-S5-r2.md`.
- Modificados: `oraculo/pre_profeco.py`, `oraculo/pre_datos.py`, `oraculo/pre_calculo.py`, `oraculo/pre_genera_fixtures.py` (`profeco_filtro()`, `carritos()`, `NOMBRE_PROD`), `oraculo/salidas/s5.json` (regenerado), `backend/internal/precios/compare.go`, `backend/internal/precios/eval_test.go`.

## Verificación
- `python3 -I oraculo/pre_genera_fixtures.py` → «29 fixtures»; `diff -rq` contra la copia previa: solo las 2 fixtures nuevas y `profeco-filtro.csv`.
- `python3 -I oraculo/pre_run_fixtures.py` → «29 fixtures, 29 coinciden, 1134 hojas de esperado, 0 diferencias».
- Desde la raíz: `go build -o datos/tmp/motores ./cmd/motores` (en `backend/`), luego `datos/tmp/motores s5 dump --fixtures oraculo/fixtures/s5 --out datos/tmp/motor-s5.json` y `python3 -I oraculo/compara.py s5` → oráculo/esperado 747/747, motor/esperado 747/747, oráculo/motor 747/747, **provisionales 301/301** (antes 225; +76 hojas de las 2 fixtures nuevas, todas PROVISIONALES), **diferencias: 0**, código 0.
- Salida del motor en las fixtures nuevas: `s5-profeco-filtro` read 8, filtered_out 3, loaded 4, duplicates 1, 4 filas de coverage; `s5-carritos-presentaciones` tercera/walmart/bodega con 16500/29650/42600, missing 2/1/0, stale 0/0/4 y el sustituto del arroz para walmart.
- `python3 -I oraculo/pre_verify.py` → «TODO COINCIDE» (ganadores 22/3/8; Walmart +138.50 = 13850).
- S5-13: `go vet ./internal/precios/... ./cmd/motores/` limpio; `gofmt -l internal/precios cmd/motores` vacío; `go test -count=1 -coverprofile ./internal/precios/...` PASS: precios 96.1 %, asistida 100 %, ocr 99.3 %, profeco 95.3 %, total del perfil 88.8 %. Ninguna prueba FAIL. Con `-run 'Comparativa|Unidad|MejorTienda|Frescura|Fuente|LecturaAsistida|Membresia|Mercado|PreciosViejos|OCRSinConfirmar|EvalCarritos'` hay 27 PASS (S5-01, 02, 03, 04, 07, 09, 10, 11, 12 y 06 sin confirmar), incluida `TestEvalCarritosVariasPresentaciones`.
- S5-08: el grep de `net/http|http.Get|http.Client` en `internal/precios` sin `_test.go` da 0. S5-07: el grep de dominios de tiendas da 0. `float64` solo aparece en `ocr/tsv.go` (confianza, no dinero).
- S5-07 sigue **NO CUMPLE**: `docs/estado/S5-terminos.md` no existe (bloqueo `S5-terminos`, de Axel).
- Criterios S5-05 y S5-06: sin cambios respecto de r2 (no toqué el cargador ni el OCR). S5-06 sigue «sin medir» porque faltan las 5 fotos de Axel.

## Supuestos
- Ninguno que cambie cifras canon. Las 2 opciones A se eligieron así: en el filtro, el punto (1) es la lectura literal de §11 («municipios de las sucursales»), y los puntos (2) y (3) son lo que ya hacía el motor y deja capturar el municipio sin la ortografía exacta de PROFECO. En carritos, el punto (1) es la lectura literal de §7.2 paso 3 («su propia cantidad y su propio importe»), el (2) da a `missing_count` el sentido de §7.3 y el (3) es la convención 4 del oráculo. Todo va marcado PROVISIONAL citando su bloqueo.

## Bloqueos
- `docs/estado/bloqueos/S5-terminos.md` (abierto, de Axel; R2-01).
- `docs/estado/bloqueos/S5-profeco-filtro.md` (nuevo; opción A PROVISIONAL).
- `docs/estado/bloqueos/S5-carritos-detalle.md` (nuevo; opción A PROVISIONAL).

## Pendientes
- S5-07: lo resuelve el orquestador con Axel en la puerta 2. Con la opción D, el integrador ajusta S5-07 y el plan §3.4. Con A o B, se escribe `docs/estado/S5-terminos.md` y se revisa solo S5-07.
- Observación que no verifiqué y que no está en los hallazgos: en `compare_carts`, el motor redondea el total de la tienda una sola vez (suma exacta) y el oráculo suma importes de línea ya redondeados (`_agrupa`). Con cantidades a granel no enteras podrían diferir en un centavo. Ninguna fixture usa cantidades no enteras en carritos. Se lo dejo al revisor.
- Hallazgos bajos R2-05, R2-06 y R2-07: no venían en el encargo y no los toqué.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S5-interfaz.md` §11, pasar a 1.5 (R2-02 y R2-04). Cambios: (a) en «Archivos» dejar «los tres cortes (2024, 2025 y 2026) en RAR, en `datos/semilla/profeco/raw/<corte>-rar/` y extraídos a `raw/<corte>-csv/`»; (b) en «Almacenamiento» la base es `datos/semilla/profeco/db/qqp.sqlite`; (c) en «Filtro» escribir la opción A de `S5-profeco-filtro` con la marca «PROVISIONAL (bloqueo `S5-profeco-filtro`)»: municipios solo de sucursales, comparación normalizada con la regla del bloqueo, sucursal dentro de su cadena y `coverage[]` con el texto tal como llega; (d) definir `rejections[].file` como la ruta tal como llega en `files[]`, y `row`, `renglon` y el número de `source_ref` como el ordinal del registro de datos (1 = el primero tras el encabezado). Motivo: hoy el §11 contradice lo implementado y no fija el filtro.
2. Mismo contrato, §7.2 y §15.2 (R2-03 y R2-04). Cambios: escribir la opción A de `S5-carritos-detalle` con la marca PROVISIONAL (ganador del sustituto = importe total ÷ contenido total; `missing_count`; `stale_count` por presentación agrupada); el orden de `rows[]` (idénticos en el orden de la canasta base, sustitutos por producto en orden de aparición, `only_one` por canasta y renglón); `stores[]` por total ascendente y, a igual total, en el orden base y luego `others`; `winners[]` en el orden de `stores[]`, con `count` 0; y la definición de sustituto (producto con presentaciones sin emparejar en dos o más canastas). Motivo: hoy esto solo lo dicen el motor y el oráculo.
3. Mismo contrato, §7.4: «`cheapest_store` es la primera tienda completa (sin faltantes) y `difference_vs_cheapest_cents` se mide contra ella, así que puede ser negativa en una tienda incompleta». §12: mencionar `datos/semilla/tickets-README.md` (instrucciones) y `datos/semilla/tickets-sinteticos/` (imágenes sintéticas). Motivo: R2-04.
4. `docs/criterios/S5.md` S5-05: «Los datasets … 2024, 2025 y 2026 (los tres en RAR)». Motivo: R2-04; hoy dice «este último en RAR».
5. `docs/plan/plan.md` §3.4 fila S5: «PROFECO 2024–2026 cargado». Motivo: R2-04.
6. `docs/estado/wf1/oraculo-S5.md` Supuestos, convención 11: que diga «filtro según la opción A del bloqueo `S5-profeco-filtro`» en vez de «igualdad exacta; municipio de cualquier tienda». Motivo: ya no describe el oráculo.
7. Orquestador: llevar `S5-terminos`, `S5-profeco-filtro` y `S5-carritos-detalle` a la puerta 2.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna de K31. Las opciones A de los dos bloqueos nuevos son propuestas para la puerta 2.

## Consumo aproximado
~37 acciones de herramienta; ~170 KB leídos; ~38 KB escritos.
