# REVISION-S5-r2 · WF-1 · revisión independiente del motor de precios y tiendas (S5), ronda 2 de 3

Revisor `revisor-S5-r2` · 11-oct-2026. Evalué solo contra `docs/criterios/S5.md` 1.2, `docs/plan/decisiones.md` §E y §K y `docs/contratos/motores/S5-interfaz.md` 1.4. No corregí nada. Todas las verificaciones las volví a correr yo; mis scripts y fixtures de prueba vivieron en `datos/tmp/rev-s5r2/` y ya los borré.

## Veredicto: NO APROBADO

Hay 1 hallazgo alto, 3 medios y 3 bajos. El alto (S5-07, falta `docs/estado/S5-terminos.md`) **no lo puede corregir ningún agente**: depende de que Axel responda el bloqueo `S5-terminos`, que ya está abierto con el formato del protocolo. Los 8 hallazgos de la ronda 1 quedan cerrados en la conducta del motor, salvo ese. Los medios son dos divergencias latentes entre el motor y el oráculo, que comprobé con fixtures propias, y un contrato que no recoge lo que ya hacen el motor y el oráculo.

## Tabla de criterios

| Criterio | Estado | Evidencia propia |
|---|---|---|
| S5-01 (alta) comparativa canon | CUMPLE | `-run Comparativa`: 4 PASS (leen `datos/semilla/data.json`). Mi cálculo con `Decimal` desde `data.json`: Bodega 309450, Walmart 323300, diferencia **+13850**; idénticos 277350 contra 287400, diferencia **+10050**; piezas 61 / 60 / 55; renglones 33 / 3. +138.50 es el total con sustitutos y el sexto atún (21.00); +100.50 es solo lo idéntico. Ver la verificación (b). |
| S5-02 (alta) unidad base | CUMPLE | `-run Unidad`: 2 PASS. Mi cálculo por 1 kg o 1 l con HALF_UP: arroz 2933 contra 2800 (gana Walmart), Coffee Mate 12692 contra 16406 (gana Bodega), leche 3375 contra 2833 (gana Walmart). `s5-unidad-sustitutos`: motor = oráculo = esperado. |
| S5-03 (alta) mejor tienda | CUMPLE (desempates PROVISIONALES, bloqueo `S5-reglas-comparativa`) | `-run MejorTienda`: 5 PASS. `s5-mejor-tienda-*` y `s5-canasta-regla-unica` sin diferencias. `stores[]` va por total ascendente con las incompletas incluidas (§7.4). |
| S5-04 (alta) frescura y fuentes | CUMPLE | `-run Frescura`: 3 PASS. `-run Fuente`: 4 PASS. Dump de `s5-frescura-fronteras`: 29 → 1/1 «1.0000», sin marca de viejo; 30 → 60/61 «0.9836», viejo; 31 → 59/61; 60 → 30/61 «0.4918»; 90 → 0 y no compite; 91 → no compite y queda marcado. `python3 -I -c "print([round((90-d)/61,4) for d in (30,60,90)])"` → `[0.9836, 0.4918, 0.0]`. `s5-fuentes-compiten`: `manual`, `ticket_manual` y `ticket_ocr` a 29.00 compiten. `site_assisted` a 10.00 y `profeco` a 9.00 salen solo en `references` y nunca ganan (`TestFuenteReferencia`). |
| S5-05 (alta) PROFECO | CUMPLE (ver R2-06, baja) | Los 3 RAR están en `raw/2024-rar`, `raw/2025-rar` y `raw/2026-rar`. Mi `wc -l` menos encabezados da 2024 14,889,219, 2025 13,270,524 y 2026 9,388,457. `sqlite3 -readonly qqp.sqlite` da 37548200 y `MAX(fecha_registro)` 2026/08/31. Por año salen las mismas tres cifras, con `MIN(renglon)` = 1. Volví a correr `cobertura -dir <3 cortes>` (4 min 2 s) y la salida es **idéntica byte a byte** a la sección 6 del documento. La muestra pesa 601 KB (≤ 2 MB), y regenerada con `-muestra -n 2000` sobre el corte 2026 es idéntica. El lector trabaja en flujo (`bufio.NewReaderSize(r, 1<<20)` + `csv.NewReader`). El documento cita `QQP_diccionario.csv`. |
| S5-06 (alta) OCR | SIN MEDIR (declarado con honestidad) | `tesseract 5.5.3`, `spa` en `--list-langs`. `ls datos/semilla/tickets/ \| wc -l` = 0 (solo un `.gitkeep` oculto). `ocrmedir` → «S5-06 sin medir: hay 0 imágenes y se necesitan 5: faltan las fotos de Axel». `TestOCRSinConfirmar` PASS: sin confirmar devuelve `ErrNotConfirmed` y ningún renglón. |
| S5-07 (alta) lectura asistida | **NO CUMPLE** (R2-01) | `-run LecturaAsistida` PASS: rechaza la activación sin veredicto. El grep de dominios de tiendas en Go sin `_test.go` da 0. `docs/estado/S5-terminos.md` **no existe**. |
| S5-08 (alta) sin red | CUMPLE | El grep de `net/http\|http.Get\|http.Client` en `internal/precios` sin `_test.go` da 0. Las URLs del reporte son solo `datos.profeco.gob.mx`. |
| S5-09 (alta) Sam's | CUMPLE | `-run Membresia`: 2 PASS. Mi cálculo: 100000 ÷ 12500 = 8; 100000 ÷ 12000 = 8.33, que sube a 9; 100000 ÷ 13850 = 7.22, que sube a 8; 80000 ÷ 13850 = 5.78, que sube a 6. Con 6 compras el ahorro es 83100 ≥ 80000, así que conviene. Sin costo: `null` con «sin costo capturado (PROVISIONAL, S5-09)». |
| S5-10 (media) trazabilidad | CUMPLE | `TestFuenteTrazabilidad` y `TestFuentePromocionVencida` PASS. `s5-promocion-vigencia` sin diferencias. |
| S5-11 (media) mercados | CUMPLE | `-run Mercado`: 2 PASS. `s5-mercado-visita`: visita en martes desde sáb 10, lun 12 y mar 13 → 13-oct; desde mié 14 → 20-oct. La avena no se compra en el tianguis (`quote` nulo). |
| S5-12 (media) precios viejos | CUMPLE | `-run PreciosViejos`: 1 PASS. `s5-precios-viejos` sin diferencias. |
| S5-13 (alta) calidad Go | CUMPLE (ver R2-07, baja) | `go vet ./internal/precios/... ./cmd/motores/` limpio. `gofmt -l` vacío. Cobertura: precios 96.1 %, asistida 100 %, ocr 99.3 % (`exec.go` `Run` 100 %), profeco 95.3 %. El total del perfil, con los dos `main` sin pruebas, da 88.8 %. Hay `float64` solo en `ocr/tsv.go` (confianza, no dinero). |

## Verificaciones pedidas (a)–(g)
- **(a)** `go test -count=1 -coverprofile ./internal/precios/...`: los cuatro paquetes ≥ 85 % (ver S5-13) y todo PASS.
- **(b)** El canon de mi script coincide (ver S5-01). Con `motores s5 dump --fixtures oraculo/fixtures/s5` **desde la raíz del repositorio** y el oráculo vuelto a correr (`pre_run_fixtures.py`: 27 fixtures, 0 diferencias; salida idéntica a `oraculo/salidas/s5.json`), `python3 -I oraculo/compara.py s5` da oráculo/esperado 747/747, motor/esperado 747/747, oráculo/motor 747/747, provisionales 225/225, **diferencias: 0**, código 0. Si el dump corre desde `backend/`, da `invalid_argument` en 8 fixtures y 419/747 (R2-05).
- **(c)** La frescura cumple K40 y K37 (ver S5-04). PROFECO nunca compite (`TestFuenteReferencia` y `s5-fuentes-compiten`).
- **(d)** Los archivos existen, la cobertura se reproduce con un comando (idéntica), la muestra pesa 601 KB y el cargador trabaja en flujo.
- **(e)** Tesseract 5.5.3 `spa`, «sin medir» declarado y ningún precio sin confirmar.
- **(f)** El punto de equilibrio y los mercados están correctos (ver S5-09 y S5-11).
- **(g)** `sha256(datos/semilla/data.json)` = `git show HEAD:datos/semilla/data.json` = `2fe7afae…31d31b`, sin cambios en git. Independencia: `oraculo/pre_*.py` no lee nada de `backend/`. El Go solo lee de `oraculo/` los datos de entrada `oraculo/fixtures/s5/_datos` (CLI `s5.go` y la prueba de `RepoPath`), nunca `esperado` ni `salidas`.

## Hallazgos

### R2-01 (alta) · S5-07: no existe `docs/estado/S5-terminos.md` (bloqueado por Axel)
- Ruta: `docs/estado/S5-terminos.md` (falta); `docs/estado/bloqueos/S5-terminos.md`.
- Qué pasa: `grep -c "Veredicto"` debe dar 3 (Walmart, Bodega Aurrera y Sam's Club, con URL y fecha) y el archivo no existe. La corrección hizo lo que pide el protocolo: abrió el bloqueo `S5-terminos` (opciones A–D), no entró a sitios de tiendas y no escribió un archivo con veredictos falsos. Aun así, el criterio alto sigue sin cumplirse. Lo demás de S5-07 sí cumple: el interruptor rechaza la activación y no hay URLs de tiendas en el código.
- Corrección sugerida: ninguna para los agentes. El orquestador lleva el bloqueo `S5-terminos` a Axel (puerta 2). Si Axel elige (D), el integrador ajusta S5-07 y el plan §3.4. Si elige (A) o (B), se escribe el documento y se vuelve a revisar solo S5-07.

### R2-02 (media) · Divergencia latente entre el motor y el oráculo en el filtro de `profeco_import`
- Ruta: `backend/internal/precios/profeco/filter.go`; `backend/cmd/motores/s5.go`; `oraculo/pre_profeco.py`; `docs/contratos/motores/S5-interfaz.md` §11.
- Qué pasa (comprobado con dos fixtures propias derivadas de `s5-profeco-muestra.json`, que luego borré): (1) Si la cadena `walmart` lleva municipio GUSTAVO A. MADERO y no hay sucursal, en producción el motor da `loaded` 1, `filtered_out` 3 y sin cobertura WALMART·GAM; el oráculo da 2 y 2, con WALMART·GAM. En revisión, el motor da 2/3 y el oráculo 3/2. Es el H-01 de la ronda 1: la fixture se cambió para no tocarlo, pero la diferencia sigue. (2) Si `profeco_name` = «Wal-mart» (la ortografía real de PROFECO) y el CSV dice «WALMART», el motor carga 2 (compara con `Fold`) y el oráculo 1 (igualdad exacta). El §11 no fija ninguna de las dos cosas, y ninguna fixture prueba esos casos.
- Corrección sugerida: que el dueño del contrato fije en §11 (a) si la `municipality` de una tienda `chain` entra al filtro y (b) cómo se comparan `cadena_comercial`/`municipio` con `profeco_name`/`municipality` (`Fold` o exacto). Después, alinear el motor o el oráculo y sumar a `pre_genera_fixtures.py` una fixture con los dos casos.

### R2-03 (media) · Divergencia latente entre el motor y el oráculo en `compare_carts`
- Ruta: `backend/internal/precios/compare.go`; `oraculo/pre_calculo.py`; `docs/contratos/motores/S5-interfaz.md` §7.2 y §15.2.
- Qué pasa (comprobado con fixtures propias derivadas de `s5-unidad-sustitutos.json`): (1) Agregué una tercera tienda con dos presentaciones de arroz (750 g a 15.00 y 1 kg a 50.00). El motor da como ganador del sustituto a `walmart`, porque usa costo total ÷ contenido total (37.14/kg). El oráculo da `tercera`, porque toma la mejor presentación (20.00/kg). (2) `missing_count`: si a Walmart le quito la leche, el motor da 1 y el oráculo 0. Con la tercera tienda, el motor da 2 y el oráculo 0. §15.2 enumera `missing_count` y `stale_count` en la salida, pero §7.2 no los define ni dice cómo se decide un sustituto cuando una tienda tiene varias presentaciones. El `stale_count` por línea o por presentación agrupada lo reporta la corrección y yo no lo medí. Ninguna fixture lo prueba porque el esperado no trae esas claves.
- Corrección sugerida: definir las tres cosas en §7.2 (ganador del sustituto con varias presentaciones, `missing_count` y `stale_count` en carritos). Después, alinear el motor o el oráculo y agregar esas claves y un caso con varias presentaciones a una fixture.

### R2-04 (media) · El contrato S5 1.4 contradice lo implementado (no recoge las solicitudes de `correccion-S5-r1`)
- Ruta: `docs/contratos/motores/S5-interfaz.md` §7.2, §7.4, §11, §12 y §15.2; `docs/criterios/S5.md` S5-05; `docs/plan/plan.md` §3.4 fila S5.
- Qué pasa: el §11 dice «2024 y 2025 en CSV», carpetas `raw/<corte>/`, extracción «a `datos/tmp/`» y base en `datos/profeco/qqp.sqlite`. Lo real es que los tres cortes vienen en RAR, en `raw/<corte>-rar/` y `raw/<corte>-csv/`, y la base está en `datos/semilla/profeco/db/qqp.sqlite`. El §11 y el §15.2 no definen `rejections[].file` (ruta tal como llega) ni `row`, `renglon` y `source_ref` (ordinal del registro de datos, 1 tras el encabezado). El §7.2 y el §15.2 no fijan el orden de `rows[]`, `stores[]` y `winners[]` ni la definición de sustituto. El §7.4 no dice que `cheapest_store` es la primera tienda completa ni que `difference_vs_cheapest_cents` puede ser negativa en una incompleta. El §12 no menciona `tickets-README.md` ni `tickets-sinteticos/`. El motor, el oráculo y las fixtures coinciden entre sí, pero el siguiente implementador (API `importProfeco`, S4, S7) que lea el contrato tendría que adivinar o haría otra cosa.
- Corrección sugerida: que el integrador o el agente del contrato S4S5 apliquen las 5 solicitudes de `docs/estado/wf1/correccion-S5-r1.md` («Solicitudes fuera de alcance») y suban el contrato a 1.5. En el criterio S5-05: «los tres cortes vienen en RAR». En el plan §3.4: «PROFECO 2024–2026».

### R2-05 (baja) · `motores s5 dump` depende del directorio de trabajo
- Ruta: `backend/cmd/motores/s5.go` (`precios.RepoPath`); `backend/cmd/motores/main.go`.
- Qué pasa: corrido desde `backend/` (con `--fixtures ../oraculo/fixtures/s5`), 8 fixtures salen con `invalid_argument` y `compara.py` da 419/747 y código 1. Desde la raíz da 747/747. El código de error esconde «archivo no encontrado», y el mensaje de uso no dice que hay que correrlo desde la raíz.
- Corrección sugerida: resolver las rutas contra la raíz del repositorio (buscar `go.mod` o `.git` hacia arriba) o fallar con un mensaje que nombre la ruta; anotar «desde la raíz» en §15 y en el uso.

### R2-06 (baja) · `S5-cobertura-profeco.md`: §1 y §4.2 con cifras del corte 2026 y §6 con los tres cortes
- Ruta: `docs/estado/S5-cobertura-profeco.md` §1, §4.2 y §5.
- Qué pasa: el «Resumen para Axel» y el §4.2 dan Bodega Vallejo 4,461, Sam's 24,788, Tepeyac 26,633 y 74,154 de 9,388,457. La sección generada da Vallejo 23,809 y «SÍ hay registros de Sam's Club: 94,293», frente al «No como tienda» del §1. La línea de versión lo avisa, pero las filas de la tabla no dicen de qué corte son. El §5 no dice que la muestra se regenera con `-dir` del corte 2026 (lo comprobé: así sale idéntica).
- Corrección sugerida: regenerar el §1 y el §4.2 con los tres cortes, o rotular cada cifra con su corte; aclarar junto al resultado generado de Sam's que es la dirección de un Wal-mart; escribir el `-dir` del comando de la muestra.

### R2-07 (baja) · `profeco/cobertura/main.go` sin pruebas
- Ruta: `backend/internal/precios/profeco/cobertura/main.go` (`filtro`, `loadQQP`, `buildSample`, `scan`: 0 %).
- Qué pasa: genera la evidencia de S5-05 y la muestra, y no tiene pruebas. El total del perfil sigue en 88.8 %, así que no incumple K25.
- Corrección sugerida: mover `buildSample` y `filtro` al paquete `profeco` con pruebas, o probar el `main` sobre la muestra.

## Estado de los hallazgos de la ronda 1
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| [alta] S5-05/(b) filtro de municipio en `profeco_import` | CERRADO en la fixture (0 diferencias); la divergencia de fondo **sigue** → R2-02 (media) | `compara.py s5` 0 diferencias; mi fixture de cadena con municipio reproduce el 1/3 contra 2/2. |
| [alta] S5-05 solo el corte 2026 | CERRADO | 2024, 2025 y 2026 cargados = `wc -l` − encabezados; total 37,548,200; MAX 2026/08/31; §6 reproducida idéntica. |
| [alta] S5-07 falta `S5-terminos.md` | SIGUE ABIERTO → R2-01 (alta) | Ya hay bloqueo `S5-terminos` bien formado; depende de Axel. |
| [alta] (b) `compara.py s5` con 65 diferencias | CERRADO | 747/747, diferencias 0, código 0 (desde la raíz). |
| [media] orden de `rows[]` en `compare_carts` | CERRADO en la conducta; el texto del contrato falta → R2-04 | `s5-comparativa-canon` `canon-filas`: motor = oráculo = esperado. |
| [media] orden de `stores[]` en `compare_basket` | CERRADO; texto de `cheapest_store` → R2-04 | `s5-mejor-tienda-incompleta` sin diferencias; orden por total ascendente. |
| [media] `rejections[].file` y `row` | CERRADO en la conducta; texto → R2-04 | `s5-profeco-muestra` sin diferencias; `MIN(renglon)` = 1 en los tres cortes. |
| [media] carpeta de tickets | CERRADO | `ls datos/semilla/tickets/ \| wc -l` = 0 (`.gitkeep` oculto); `tickets-README.md` y `tickets-sinteticos/` junto a la carpeta. |
