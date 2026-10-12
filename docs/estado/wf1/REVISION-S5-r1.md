# REVISIÓN S5 · ronda 1 de 3 · motor de precios y tiendas

Revisor independiente `revisor-S5-r1` · WF-1 · 11-oct-2026. No escribí nada de lo revisado y no corregí nada. Referencias únicas: `docs/criterios/S5.md` 1.2, `docs/plan/decisiones.md` §E y §K, `docs/contratos/motores/S5-interfaz.md` 1.4. Toda evidencia de esta tabla sale de comandos y scripts míos (temporales en `datos/tmp/`, ya borrados), no de los reportes.

## Veredicto: NO APROBADO

4 hallazgos de severidad alta, 4 de severidad media y 4 bajos. Lo esencial del motor está bien: el canon (+$138.50 y +$100.50), la unidad base, la frescura K40, las fuentes K37, Sam's, mercados, la cobertura de pruebas y la independencia del oráculo. Lo que impide aprobar es: (1) `python3 -I oraculo/compara.py s5` sale con código 1 y **65 diferencias motor contra oráculo** (nadie lo había corrido: los dos reportes lo dejaron pendiente); (2) faltan los cortes PROFECO 2024 y 2025 (S5-05); (3) no existe `docs/estado/S5-terminos.md` (S5-07) y no hay bloqueo que lo explique.

## Tabla de criterios

| Criterio | Estado | Evidencia mía |
|---|---|---|
| S5-01 (alta) comparativa canon | CUMPLE (totales); ver H-05 por el orden de `rows[]` | `go test ./internal/precios/... -run Comparativa -count=1 -v`: 4 PASS (`TestComparativaCanon` lee `datos/semilla/data.json` y `productos.json`). Cálculo propio con `Decimal` desde `data.json`: idénticos 277350 / 287400 (Σ unitario × cantidad = Σ `*_total`), sustitutos 30000 / 35900, sexta lata 2100, totales 309450 / 323300, Walmart +13850 = +10050 + 5900 − 2100; piezas 61 / 60 / 55; renglones 33 / 3; ganadores 22 / 8 / 3 = `resultados_por_producto`. `compara.py s5`: 0 diferencias en `canon-totales`. |
| S5-02 (alta) unidad base | CUMPLE | `-run Unidad`: 2 PASS. Propio (HALF_UP por 1 kg / 1 l): arroz 2933 vs 2800, Coffee Mate 12692 vs 16406, leche 3375 vs 2833. `s5-unidad-sustitutos`: motor = oráculo = esperado. |
| S5-03 (alta) mejor tienda | CUMPLE con reservas (H-06, H-12) | `-run MejorTienda`: 5 PASS; canasta del 10-oct → Bodega. `compara.py s5`: `s5-mejor-tienda-canasta`, `-empate-*`, `s5-cotizacion-recomendada` y `s5-canasta-regla-unica` sin diferencias; `s5-mejor-tienda-incompleta` con 7 diferencias de orden (mismos totales y misma `cheapest_store` t2). |
| S5-04 (alta) frescura y fuentes | CUMPLE | `-run Frescura` (3 PASS) y `-run Fuente` (4 PASS). Propio con `Fraction`: 29 → 1 (1.0000); 30 → 60/61 (0.9836); 31 → 59/61 (0.9672); 60 → 30/61 (0.4918); 90 → 0. `TestFuenteReferencia`: PROFECO a 1.00 y `site_assisted` a 2.00 contra `manual` a 28.00: gana Walmart (`manual`); Bodega queda con `missing_count` 1. Parámetros desde `app_setting` (`types.go:97-99`, `eval.go:96-98`), sin constantes 30/90 en código. `s5-frescura-fronteras`, `-parametros`, `s5-fuentes-compiten`: sin diferencias. |
| S5-05 (alta) PROFECO | NO CUMPLE (H-02); 2026 sí | RAR en `datos/semilla/profeco/raw/2026-rar/QQP_2026.rar` (58 MB) y 16 CSV (2.8 GB). `wc -l` total 9,388,473 = 9,388,457 + 16 encabezados; `sqlite3 …/db/qqp.sqlite "SELECT COUNT(*),MAX(fecha_registro)"` → 9388457 · 2026/08/31. SQL propio por cadena: Wal-mart 1,500,220 · Wal-mart Express 412,731 · Bodega Aurrera 1,020,792 · Bodega Aurrera Express 33,491; Sam's Club 0. Reproducción con **un** comando (`go run ./internal/precios/profeco/cobertura -dir …/2026`, 55 s): su salida es idéntica (diff de 0 líneas) a la sección 6 de `S5-cobertura-profeco.md`. Cita `QQP_diccionario.csv`; dice SÍ/NO por sucursal. Muestra 615,625 bytes (≤ 2 MB). Cargador en flujo: `bufio` 1 MB + `csv.Reader` con `ReuseRecord`; pico de 29 MB con un CSV de 197 MB. **Faltan 2024 y 2025.** |
| S5-06 (alta) OCR | SIN MEDIR (declarado con honestidad) | `tesseract --version` 5.5.3; `--list-langs` osd, spa. `go run ./internal/precios/ocr/cmd/ocrmedir` → «S5-06 sin medir: hay 0 imágenes y se necesitan 5: faltan las fotos de Axel». `TestOCRSinConfirmar` PASS (`PricesFromOCR(…, false)` → `ErrNotConfirmed`, sin renglones). La carpeta de tickets choca con §12 (H-08). |
| S5-07 (alta) lectura asistida | NO CUMPLE (H-03) | `go test ./internal/precios/asistida/ -run LecturaAsistida` PASS (rechaza sin veredicto); grep de dominios de tiendas en `backend` sin `_test.go`: 0 líneas. **`docs/estado/S5-terminos.md` no existe**: `grep -c Veredicto` da error, se esperaban 3. |
| S5-08 (alta) sin red | CUMPLE | `grep -rln "net/http\|http.Get\|http.Client" internal/precios --include='*.go' \| grep -v _test.go` → vacío. |
| S5-09 (alta) Sam's | CUMPLE | `-run Membresia`: 2 PASS. Propio: 100000 ÷ 12500 = 8; 100000 ÷ 12000 = 8.33 → 9 (hacia arriba, PROVISIONAL p4b). `s5-membresia-*` sin diferencias; sin costo → `missing_fee`. |
| S5-10 (media) trazabilidad | CUMPLE | `TestFuenteTrazabilidad`: 5 fuentes × 1 precio, compiten las 3 primeras; `TestFuentePromocionVencida` PASS; `s5-promocion-vigencia` sin diferencias. |
| S5-11 (media) mercados | CUMPLE | `-run Mercado`: 2 PASS; `s5-mercado-visita` sin diferencias (sábado 10-oct + martes → 13-oct). |
| S5-12 (media) precios viejos | CUMPLE | `-run PreciosViejos` PASS; `s5-precios-viejos` sin diferencias. |
| S5-13 (alta) calidad Go | CUMPLE | `go vet ./internal/precios/... ./cmd/motores/` limpio; `gofmt -l internal/precios cmd/motores` vacío. Cobertura: precios 96.0 %, asistida 100 %, ocr 99.3 % (sin `exec.go`: 295/297 = 99.3 %, también sin la prueba de integración), profeco 95.3 %; total del perfil 88.9 % (incluye los dos `main` sin pruebas). `float64` solo en `ocr/tsv.go` (promedio de confianza 0–100), nada monetario. |

## Verificaciones pedidas (a)–(g)

- **(a)** Cobertura ≥ 85 % en los cuatro paquetes con pruebas; ocr sin contar exec 99.3 %. Cumple.
- **(b)** Canon reproducido por mí desde `data.json` (arriba). `compara.py s5`: oráculo contra esperado 747/747; **motor contra esperado 713/747; 65 diferencias; código 1** (H-01, H-05, H-06, H-07). `pre_run_fixtures.py`: 27/27, 0 diferencias; `pre_verify.py`: «TODO COINCIDE».
- **(c)** Frescura K40 y fuentes K37: cumple (tabla).
- **(d)** PROFECO: archivo descargado existe, cobertura reproducible con un comando, muestra ≤ 2 MB, cargador en flujo; faltan 2024 y 2025 (H-02).
- **(e)** OCR: tesseract 5.5.3 spa; «sin medir» declarado; nunca escribe sin confirmar. Cumple.
- **(f)** Sam's y mercados: cumple.
- **(g)** Independencia: `oraculo/pre_*.py` no nombra `backend` ni código Go; `pre_genera_fixtures.py` solo importa biblioteca estándar (no llama al oráculo); las pruebas Go leen `data.json` y comparan contra `totales.*` leídos del archivo. `sha256sum datos/semilla/data.json` = `2fe7afae…1031d31b` y `git status` limpio para ese archivo. Cumple.

## Hallazgos

### H-01 (alta) · `profeco_import`: el motor carga y filtra distinto que el oráculo
- Ruta: `backend/cmd/motores/s5.go:131-139` (pasa municipio solo a las sucursales) y `backend/internal/precios/profeco/filter.go:43-56` («sin sucursales con municipio no pasa ningún registro»), contra `oraculo/pre_profeco.py:71` (municipio de cualquier tienda) y `oraculo/fixtures/s5/s5-profeco-muestra.json`.
- Qué pasa: en la fixture, `walmart` es una cadena con `municipality` GUSTAVO A. MADERO y no hay sucursal de Walmart. Producción: motor `loaded` 1 y `filtered_out` 3; esperado y oráculo 2 y 2, y falta en el motor la cobertura WALMART · GUSTAVO A. MADERO. Revisión: motor 2 y 3; esperado 3 y 2. Son cifras que no coinciden con el oráculo. §11 dice «municipios de las sucursales (`store.municipality`)», que admite las dos lecturas.
- Corrección sugerida: el integrador fija en S5-interfaz §11 si la `municipality` de una tienda `chain` entra al filtro; luego se alinea el adaptador y `filter.go`, o la fixture (poner el municipio en una sucursal), y se vuelve a correr `compara.py s5` hasta 0 diferencias.

### H-02 (alta) · S5-05: faltan los cortes PROFECO 2024 y 2025
- Ruta: `datos/semilla/profeco/raw/` (solo `2026-rar/` y `2026-csv/`), `docs/estado/S5-cobertura-profeco.md` (título «corte 2026»; §1 «Lo que falta», punto 1).
- Qué pasa: el criterio exige 2024, 2025 y 2026 cargados con conteo contra fuente. El agente PROFECO lo marca «parcial» como pendiente, sin bloqueo. Plan §3.4 fila S5 dice «PROFECO 2026 cargado», así que el plan y el criterio no dicen lo mismo.
- Corrección sugerida: descargar 2024 y 2025 de datos.profeco.gob.mx en `raw/2024/` y `raw/2025/`, cargarlos con el mismo comando y regenerar el reporte con su conteo. Otra salida: si el orquestador toma el plan («2026») como lo que vale, abrir el bloqueo `docs/estado/bloqueos/S5-profeco-cortes.md` y que la puerta 2 ajuste S5-05.

### H-03 (alta) · S5-07: no existe `docs/estado/S5-terminos.md` y no hay bloqueo
- Ruta: `docs/estado/S5-terminos.md` (falta); entregable de S5 según criterios §1 y plan §3.4.
- Qué pasa: `grep -c "Veredicto"` debe dar 3 (Walmart, Bodega Aurrera, Sam's Club, con URL y fecha) y el archivo no existe. Ningún agente lo tiene en su reporte ni hay bloqueo. Las reglas de WF-1 prohíben entrar a sitios de tiendas («nada de sitios de tiendas»), así que ningún agente puede escribir el veredicto sin una autorización de Axel. El interruptor sí queda apagado (bien), pero el criterio pide el documento.
- Corrección sugerida: abrir `docs/estado/bloqueos/S5-terminos.md` (formato del protocolo §2) con opciones como: (A) Axel autoriza a leer solo las páginas legales de los 3 sitios; (B) Axel lee los términos y dicta el veredicto; (C) S5-07 se evalúa sin el documento hasta la puerta 2, con la lectura asistida apagada. Dejar S5-07 como «pendiente de bloqueo» y no como cumplido.

### H-04 (alta) · `compara.py s5` falla: el motor nunca se comparó con el oráculo
- Ruta: `oraculo/compara.py s5` con `datos/tmp/motor-s5.json` (generado con `motores s5 dump --fixtures oraculo/fixtures/s5`).
- Qué pasa: código 1 y 65 diferencias en 3 fixtures: `s5-comparativa-canon` · `canon-filas` 34 (H-05), `s5-mejor-tienda-incompleta` 7 (H-06), `s5-profeco-muestra` 24 (H-01 y H-07). `motor-S5-go.md` y `oraculo-S5.md` dejaron esta comparación como pendiente; ningún reporte la cierra.
- Corrección sugerida: resolver H-01, H-05, H-06 y H-07 y adjuntar la salida de `compara.py s5` con «diferencias: 0» al reporte del motor.

### H-05 (media) · `compare_carts`: el orden de `rows[]` difiere (mismo contenido)
- Ruta: `backend/internal/precios/compare.go` contra `oraculo/pre_calculo.py` y `oraculo/fixtures/s5/s5-comparativa-canon.json` (`canon-filas`); S5-interfaz §7.2 y §15.2.
- Qué pasa: lo comparé por mi cuenta. Los 37 renglones son el mismo multiconjunto y solo cambia el orden. El motor pone el atún `only_one` justo después de su idéntico (fila 17) y los sustitutos al final (34–36). El oráculo pone los idénticos, luego los sustitutos (33–35) y al final el `only_one` (36). El contrato no fija el orden (supuesto 2 del oráculo), así que el implementador siguiente tendría que adivinarlo.
- Corrección sugerida: escribir el orden en §15.2 (la solicitud del oráculo, convenciones 1–3) y alinear el lado que no cumpla.

### H-06 (media) · `compare_basket`: una tienda incompleta va al final en el motor; §7.4 ordena por total
- Ruta: `backend/internal/precios/basket.go` contra §7.4 («tiendas ordenadas por total ascendente») y `s5-mejor-tienda-incompleta`.
- Qué pasa: el motor da t2 3100, t1 3100, t3 500 (faltante 1). El esperado y el oráculo dan t3 500, t2 3100, t1 3100. Los totales y `cheapest_store` = t2 coinciden.
- Corrección sugerida: ordenar por total ascendente como dice §7.4. Si se quiere que las incompletas vayan al final, el integrador lo escribe en §7.4 y se ajusta la fixture.

### H-07 (media) · `profeco_import.rejections[]`: `file` y `row` con otra convención
- Ruta: `backend/internal/precios/profeco/import.go:125,141` (`filepath.Base` y número de línea física con encabezado) contra el oráculo (ruta tal como llega en `files[]` y ordinal del registro de datos).
- Qué pasa: motor `profeco-muestra-a.csv` · fila 4; esperado `oraculo/fixtures/s5/_datos/profeco-muestra-a.csv` · fila 3. §11 y §15.2 no definen ninguna de las dos cosas, y lo mismo afecta a `source_ref "<archivo>:<renglón>"`.
- Corrección sugerida: fijar en §11 y §15.2 qué es `file` y desde dónde se cuenta `row` (con o sin encabezado). Alinear motor o fixture.

### H-08 (media) · `datos/semilla/tickets/` no guarda solo las fotos
- Ruta: `datos/semilla/tickets/README.md` y `datos/semilla/tickets/sinteticos/`, contra S5-interfaz §12 («guarda **solo** las 5 fotos… así `ls … | wc -l` = 5») y la verificación de S5-06.
- Qué pasa: hoy `ls datos/semilla/tickets/ | wc -l` = 2 y con las 5 fotos daría 7, así que la verificación literal del criterio fallaría aunque Axel deje las 5 fotos. El agente OCR lo pidió como solicitud y sigue sin resolverse.
- Corrección sugerida: mover `README.md` y `sinteticos/` a `datos/semilla/tickets-sinteticos/` (o a `docs/`) y dejar la carpeta vacía para las 5 fotos. La otra opción es cambiar la verificación para que cuente imágenes.

### H-09 (baja) · Comando de Tesseract distinto del contrato
- Ruta: `backend/internal/precios/ocr/engine.go:92-93`. Agrega `--oem 1 -c preserve_interword_spaces=1` al comando de §12. Está documentado en `doc.go`.
- Corrección sugerida: que el integrador actualice §12 con el comando real.

### H-10 (baja) · Nombre del ejecutor del oráculo
- Ruta: S5-interfaz §15.3 nombra `oraculo/s5_run_fixtures.py`; el archivo es `oraculo/pre_run_fixtures.py`.
- Corrección sugerida: ajustar la mención o renombrar el archivo.

### H-11 (baja) · Rutas y supuestos de §11 desactualizados
- Ruta: S5-interfaz §11 dice extraer el RAR en `datos/tmp/` y guardar la base en `datos/profeco/qqp.sqlite`. Lo real es `datos/semilla/profeco/raw/2026-csv/` y `datos/semilla/profeco/db/qqp.sqlite`. Además §11 supone una cadena Sam's Club que PROFECO no publica (0 registros, comprobado por SQL).
- Corrección sugerida: actualizar §11 (la solicitud del agente PROFECO).

### H-12 (baja) · El desempate de canasta no está en `reglas-negocio.md`
- Ruta: `docs/contratos/reglas-negocio.md` (sin «empate»). S5-03 pide la regla de desempate «escrita en `reglas-negocio.md`». Hoy vive en S5-interfaz §7.4 como PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 2).
- Corrección sugerida: cuando la puerta 2 cierre la pregunta 2, pasar la regla a RN, o citar desde RN la regla PROVISIONAL de §7.4.

## Estado de hallazgos anteriores
No aplica: es la primera revisión del motor S5. Las revisiones `REVISION-contrato-S4S5-r1..r5` son del contrato, no de este código.
