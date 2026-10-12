# S5-05 · Cobertura de PROFECO (Quién es Quién en los Precios, cortes 2024, 2025 y 2026)

Versión 1.1 · 11-oct-2026 · WF-1, agente `correccion-S5-r1`: se suman los cortes 2024 y 2025 (`docs/estado/wf1/REVISION-S5-r1.md` H-02); la sección 6 se regenera con **un solo comando sobre los tres cortes**; `renglon` y `source_ref` cuentan el ordinal del registro de datos (1 = el primero tras el encabezado; H-07). Las cifras de las secciones 1 (salvo las dos primeras filas y las de tamaño y tiempo) y 4.2 son del **corte 2026** (versión 1.0); las de los tres cortes están en la sección 6. Versión 1.0 · 11-oct-2026 · WF-1, agente `profeco` (Sonnet). Los conteos de las secciones «Archivos», «Totales», «Cobertura por `cadena_comercial`», «Cadenas de interés por estado y municipio» y «Sucursales de Axel» **los genera el cargador** (no están escritos a mano) y se reproducen con un comando (sección «Cómo reproducir»). Definiciones de columnas: `datos/semilla/profeco/QQP_diccionario.csv` (15 columnas, citado por el cargador `backend/internal/precios/profeco/reader.go`). Metadatos de la fuente: `datos/semilla/profeco/QQP_metadatos.csv` (Programa Quién es Quién en los Precios, Enero 2024 – Agosto de 2026; última modificación 31/08/2026; actualización mensual; autor PROFECO).

## 1. Resumen para Axel

| Pregunta | Respuesta |
|---|---|
| ¿Se cargaron los cortes 2024, 2025 y 2026? | **Sí.** 64 archivos CSV (24 de 2024, 24 de 2025 y 16 de 2026; una quincena por archivo), **37,548,200** registros leídos = 37,548,200 guardados + 0 rechazados: 2024 **14,889,219**; 2025 **13,270,524**; 2026 **9,388,457**, cada uno igual a su `wc -l` menos un encabezado por archivo (14,889,243 − 24; 13,270,548 − 24; 9,388,473 − 16). Cero filas perdidas sin explicar. |
| ¿Hasta qué fecha llega? | 2024-01-02 a **2026-08-31** (`MAX(fecha_registro)` = `2026/08/31`), igual que la fuente (`Cobertura temporal` de `QQP_metadatos.csv`). |
| ¿Están Walmart y Bodega Aurrera? | **Sí.** PROFECO las registra con cuatro nombres de `cadena_comercial`: `Wal-mart`, `Wal-mart Express`, `Bodega Aurrera` y `Bodega Aurrera Express` (con guion y en dos formatos de tienda). **Esa es la ortografía real para `store.profeco_name`: Axel debe confirmarla.** |
| ¿Está Sam's Club? | **No como tienda.** No existe ninguna `cadena_comercial` que se llame Sam's Club ni un `nombre_comercial` con ese texto. La regla de empate (ver sección 4) encuentra 24,788 registros, pero son **una tienda Wal-mart de Cancún (Quintana Roo) cuya dirección dice «Costado Sams Club»**: una referencia de ubicación, no precios de Sam's. Para la comparativa, **Sam's Club no tiene precios de PROFECO**. |
| ¿Está Bodega Aurrera Vallejo? | **Hay una tienda candidata, falta que Axel la reconozca.** La regla de empate (texto «vallejo» en `nombre_comercial` o `direccion`) encuentra 4,461 registros de una sola tienda: *Bodega Aurrera Express Sucursal Poniente 140* (Poniente 140 No. 237, Col. Lindavista Vallejo III Sección, Gustavo A. Madero). «Vallejo» aparece solo en la dirección (colonia), no en el nombre, y es formato Express. No se asume que sea la de Axel. |
| ¿Está Walmart Acueducto de Guadalupe? | **No con la regla.** Cero registros con «acueducto de guadalupe». La única coincidencia parcial («acueducto», 21,333 registros) es *Walmart Express Sucursal Virreyes* (Aguilar y Seijas 123, esquina Acueducto Río Hondo, Miguel Hidalgo): otra tienda. En Gustavo A. Madero PROFECO solo mide un Wal-mart: *Walmart Sucursal Tepeyac* (Calz. de Guadalupe 431, Col. Guadalupe Tepeyac; 26,633 registros) y un Walmart Express (*Lindavista*). Axel dice si alguna es la suya (sección 4.3). |
| ¿Qué tan grande es? | RAR 59,971,365 bytes; CSV extraídos 2,957,095,221 bytes (16 archivos de 150 a 234 MB; el mayor está bajo el límite de 1 GB por archivo de `price.profeco_import_max_gb`); base local de consulta `qqp.sqlite` 4,293,079,040 bytes (4.0 GB; con un índice por cadena, municipio y fecha). Corte 2024: RAR 96,058,666 bytes, 24 CSV de 4,696,970,932 bytes en total; corte 2025: RAR 85,141,166 bytes, 24 CSV de 4,172,949,078 bytes. Con los tres cortes, `qqp.sqlite` pesa 17,223,405,568 bytes. |
| ¿Cuánto tarda? | Lectura y conteo de todo el corte: **57 s**. Carga completa a la base local aparte: **1 min 49 s**. Lectura + filtro de cadenas y municipio (lo que hace una importación a `price`): **29 s**. Todo en un solo hilo, leyendo en flujo (1 MB de búfer; el archivo nunca se carga completo). Con los tres cortes (versión 1.1): conteo **4 min 5 s**; carga a la base aparte **8 min 11 s**. |
| ¿Cuántos registros pasarían a `price`? | **Depende de las ligas, que hoy no existen.** `datos/semilla/profeco/vinculos.csv` no existe todavía: sin liga confirmada, todo registro filtrado va a `skipped_count` con «sin presentación ligada» (S5-interfaz §11). Ejemplo medido: las cuatro cadenas de Walmart/Bodega Aurrera en **Gustavo A. Madero** dejan pasar el filtro 74,154 registros (22,485 + 4,461 + 26,633 + 20,575), de 9,388,457 (0.79 %). |

### Lo que falta y quién lo decide
1. **Cortes 2024 y 2025:** cargados en la versión 1.1; antes faltaban. Llegaron en **RAR** (`QQP_2024.rar` y `QQP_2025.rar`), no en CSV como supone S5-interfaz §11 («2024 y 2025 en CSV»). El cargador los leyó sin cambios de lectura (mismas 15 columnas). La corrección del texto de §11 va como solicitud al integrador.
2. **Ortografía de las cadenas** (4 nombres arriba) y **municipios** tal como los escribe PROFECO (por ejemplo `Gustavo a. Madero`, `Benito Juárez` existe en la Ciudad de México y en Quintana Roo): Axel confirma; el cargador compara sin acentos, mayúsculas, guiones ni puntos.
3. **Sucursales de Axel:** ver sección 4. La selección de sucursales cercanas **no se conoce** y queda en Ajustes › Tiendas y precios.
4. **Dónde se guarda todo el archivo** (bloqueo `S5-profeco-almacen`, PROVISIONAL, opción A): se usó una base aparte, pero en `datos/semilla/profeco/db/qqp.sqlite` y no en `datos/profeco/qqp.sqlite` porque el encargo prohíbe escribir datos fuera de `datos/semilla/profeco/`. La carpeta ignora todo su contenido en git.

## 2. Qué se descargó y cómo se trató

- Fuente: portal de datos abiertos de PROFECO (`datos.profeco.gob.mx`), dataset «Quién es Quién en los Precios», corte 2026 en un RAR. Archivo: `datos/semilla/profeco/raw/2026-rar/QQP_2026.rar` (59,971,365 bytes).
- Integridad: `bsdtar -tf` lista los 16 archivos sin error; los tamaños de los 16 CSV extraídos coinciden uno por uno con los del listado del RAR (`bsdtar -tvf`).
- Extracción: `bsdtar -xf` a una carpeta propia, `datos/semilla/profeco/raw/2026-csv/` (sin scripts dentro; `raw/.gitignore` ignora todo). El contenido se trata como dato no confiable: nada se ejecuta, el cargador solo lo lee en flujo y convierte el dinero a centavos desde el texto.
- Formato real: CSV UTF-8 con BOM, 15 columnas en el orden del diccionario, una quincena por archivo (`MM-2026_Q1.csv` = días 1 a 15; `_Q2` = el resto del mes). La fecha viene como `AAAA/MM/DD` (único formato; el diccionario solo dice «Datetime»).
- Versión 1.1 (`correccion-S5-r1`): de la misma página (`https://datos.profeco.gob.mx/datos_abiertos/qqp.php`, enlaces «Quien es Quien en los Precios 2024» y «… 2025»), `raw/2024-rar/QQP_2024.rar` (sha256 `740e1498e9d7fae77a9b4d8c560e5c1feb61308e23721ecad64d981a745e2523`) y `raw/2025-rar/QQP_2025.rar` (sha256 `41ca4cd962e696f9f19cca598566d11ce60b200e37b5a5ef46c3b6a1ae58dc59`), extraídos con `bsdtar -xf` en carpetas propias `raw/2024-csv/QQP_2024/` y `raw/2025-csv/QQP_2025/`. Los 24 + 24 tamaños extraídos coinciden uno por uno con `bsdtar -tvf`. Nombres `MM-AAAA_01.csv` y `MM-AAAA_02.csv` (en 2026, `_Q1` y `_Q2`).
- Red usada: una sola descarga del RAR (la dejó el intento anterior de este agente) desde `datos.profeco.gob.mx` (fuente permitida por el plan §1.7). Esta ejecución no usó red.

## 3. Cargador (`backend/internal/precios/profeco/`)

| Archivo | Qué hace |
|---|---|
| `reader.go`, `text.go`, `parse.go` | Lee el CSV en flujo (un registro a la vez, `encoding/csv`, sin cargar el archivo); BOM, UTF-8 o Windows-1252 (tabla propia), columnas en cualquier orden y sin distinguir mayúsculas; precio a centavos exactos desde el texto (sin flotantes); fecha ISO con los formatos que declara; latitud y longitud como texto decimal. Rechaza con motivo: precio vacío, ilegible, no positivo o con más de dos decimales; fecha ilegible; producto o cadena vacíos; columnas de más o de menos; CSV mal formado. |
| `filter.go` | Filtra por `cadena_comercial` (= `store.profeco_name` de las cadenas) y por municipio (= `store.municipality` de las sucursales, con estado si la sucursal lo trae). Tienda del precio: la sucursal cuyo `profeco_name` = `nombre_comercial`; si no hay, la cadena. |
| `links.go` | Lee `vinculos.csv` con las reglas de S5-interfaz §11: `propuesto` sin `confirmo` ni fecha; `confirmado` con `axel` o `ximena` y fecha ISO. |
| `import.go` | Importación por archivo y por lote. Conteo: `Read = Rejected + FilteredOut + Skipped + Loaded + LoadedProposed + Duplicates` (método `Balanced`). Modo producción: solo ligas confirmadas; modo revisión: también las propuestas, contadas aparte (`LoadedProposed`). Duplicados (misma presentación, tienda y día): gana el primero del archivo. |
| `sqlite.go` | `LoadIntoDB`: una sola transacción sobre la base de la app (todo o nada); inserta en `price` con `source = 'profeco'`, `is_promotional = 0` y `source_ref = "QQP <corte> <archivo>:<renglón>"` (PROFECO nunca compite: K22, K37); lee `price.profeco_import_max_gb` de `app_setting` (K44), nunca de una constante. `QQP`: base local aparte, tabla `profeco_record` (las 15 columnas + `archivo`, `renglon`, `precio_cents`, `fecha`). |
| `survey.go`, `report.go` | Cobertura (este documento) sin guardar los registros. |
| `cobertura/main.go` | Comando (en `-dir` acepta varias carpetas separadas por coma, una por corte): reporte, carga a `qqp.sqlite`, medición del filtro y generación de la muestra. |

Nota sobre el desfase con el contrato: S5-interfaz §11 dice que el RAR se extrae a `datos/tmp/`; el encargo pidió extraer en `raw/` (carpeta propia, ignorada por git). Se siguió el encargo.

## 4. Sucursales de Axel: regla escrita antes de contar

### 4.1 Regla de empate (fijada antes del primer conteo)
Un registro es de la sucursal buscada si (a) la `cadena_comercial`, comparada con `Fold` (minúsculas, sin acentos, sin guion, punto ni apóstrofo), es una de las cadenas de la sucursal y (b) la frase de la sucursal aparece **como palabras completas** en `nombre_comercial` o en `direccion`. Frases: «vallejo» (cadenas `Bodega Aurrera` y `Bodega Aurrera Express`), «acueducto de guadalupe» (cadenas `Wal-mart` y `Wal-mart Express`) y «sams club» (sin exigir cadena; también en `cadena_comercial`). Las coincidencias parciales (una sola palabra: «vallejo», «acueducto», «sams») se listan aparte como pista y **no cuentan**. Además se separa cuántos registros traen la frase en el nombre y cuántos solo en la dirección, porque una dirección puede nombrar un lugar vecino sin ser ese lugar.

### 4.2 Resultado (de la sección generada «Sucursales de Axel» más abajo)
- **Bodega Aurrera Vallejo:** la regla da 4,461 registros, todos en la dirección de *Bodega Aurrera Express Sucursal Poniente 140*; ninguno trae «Vallejo» en el nombre.
- **Walmart Acueducto de Guadalupe:** **no hay registros** con la regla; solo una coincidencia parcial distinta (Virreyes, Miguel Hidalgo).
- **Sam's Club:** la regla da 24,788 registros, todos solo en la dirección de una tienda Wal-mart de Cancún («Costado Sams Club»). **Sam's Club no aparece como tienda en PROFECO.**

### 4.3 Lo que Axel decide
Las tiendas de Axel no se conocen para este cargador. Reconocerlas por nombre y dirección en la lista completa de la sección generada «Todas las sucursales de las cadenas de interés en los municipios vigilados» (Gustavo A. Madero) y en la tabla por municipio. Si ninguna es la suya, la cobertura de PROFECO para esa tienda es cero y la comparativa seguirá con precios capturados (tickets, manual).

## 5. Cómo reproducir los conteos

Desde `~/finanzas-hogar/backend` (los CSV viven en `datos/semilla/profeco/raw/2024-csv/QQP_2024/`, `raw/2025-csv/QQP_2025/` y `raw/2026-csv/2026/`; `D` abrevia las tres carpetas):

```
R=../datos/semilla/profeco/raw; D=$R/2024-csv/QQP_2024,$R/2025-csv/QQP_2025,$R/2026-csv/2026

# Reporte completo de los tres cortes (las secciones generadas de abajo); tarda ~4 min:
go run ./internal/precios/profeco/cobertura -dir $D

# Contra wc -l (14,889,243 + 13,270,548 + 9,388,473 líneas = 37,548,200 registros + 64 encabezados):
for y in $R/2024-csv/QQP_2024 $R/2025-csv/QQP_2025 $R/2026-csv/2026; do wc -l $y/*.csv | tail -1; done

# Carga a la base local aparte (8 min 11 s; borrar antes qqp.sqlite si existe) y comprobación S5-05:
go run ./internal/precios/profeco/cobertura -dir $D -qqp ../datos/semilla/profeco/db/qqp.sqlite
sqlite3 ../datos/semilla/profeco/db/qqp.sqlite "SELECT count(*), min(fecha), max(fecha), max(fecha_registro) FROM profeco_record"
#   -> 37548200|2024-01-02|2026-08-31|2026/08/31
sqlite3 ../datos/semilla/profeco/db/qqp.sqlite "SELECT substr(archivo,4,4), count(*) FROM profeco_record GROUP BY 1"
#   -> 2024|14889219 · 2025|13270524 · 2026|9388457

# Medir el filtro del cargador en el corte 2026 (sin ligas; 29 s):
go run ./internal/precios/profeco/cobertura -dir $R/2026-csv/2026 -cadenas "Wal-mart,Wal-mart Express,Bodega Aurrera,Bodega Aurrera Express" -municipios "Gustavo a. Madero"
#   -> pasan_el_filtro(sin liga)=74154

# Pruebas del cargador (con la muestra de ~2,000 filas):
go test ./internal/precios/profeco/ -cover      # 95.3 %
```

Comprobación cruzada sin Go: `sqlite3 ../datos/semilla/profeco/db/qqp.sqlite "SELECT cadena_comercial, count(*) FROM profeco_record WHERE cadena_comercial IN ('Wal-mart','Wal-mart Express','Bodega Aurrera','Bodega Aurrera Express') GROUP BY 1"` da, en el corte 2026 (versión 1.0), 1,500,220 · 412,731 · 1,020,792 · 33,491; con los tres cortes da las cifras de la tabla por cadena de la sección 6.

Diferencias entre `wc -l` y registros: **ninguna** (columna «diferencia» = 0 en los 64 archivos): en estos cortes no hay saltos de línea dentro de comillas; cada renglón es un registro. Por eso `wc -l` − 64 = 37,548,200.

La muestra de pruebas (`datos/semilla/profeco/muestra/qqp-muestra.csv`, 1,974 filas, 615,625 bytes, las 221 cadenas, con hasta 5 filas cada una salvo las de interés, que reparten el 55 % de las filas) se regenera con `-muestra <ruta> -n 2000` y es determinista.

## 6. Conteos generados

Generado con `go run ./internal/precios/profeco/cobertura -dir $D` sobre los tres cortes (versión 1.1).

### Archivos leídos: conteo contra `wc -l`

Cada archivo trae un encabezado; «diferencia» = `wc -l` − 1 − registros leídos. Cero significa que cada salto de línea es un registro; si no es cero, hay saltos de línea dentro de comillas (se cuentan como un solo registro con el lector CSV).

| archivo | bytes | `wc -l` | registros leídos | rechazados | diferencia | primera fecha | última fecha |
|---|---:|---:|---:|---:|---:|---|---|
| 01-2024_01.csv | 176,443,178 | 560,730 | 560,729 | 0 | 0 | 2024-01-02 | 2024-01-15 |
| 01-2024_02.csv | 224,177,524 | 710,527 | 710,526 | 0 | 0 | 2024-01-16 | 2024-01-31 |
| 02-2024_01.csv | 187,912,774 | 596,334 | 596,333 | 0 | 0 | 2024-02-01 | 2024-02-15 |
| 02-2024_02.csv | 196,823,163 | 624,759 | 624,758 | 0 | 0 | 2024-02-16 | 2024-02-29 |
| 03-2024_01.csv | 213,237,521 | 676,519 | 676,518 | 0 | 0 | 2024-03-01 | 2024-03-15 |
| 03-2024_02.csv | 144,940,895 | 458,900 | 458,899 | 0 | 0 | 2024-03-19 | 2024-03-27 |
| 04-2024_01.csv | 201,181,313 | 637,340 | 637,339 | 0 | 0 | 2024-04-01 | 2024-04-15 |
| 04-2024_02.csv | 214,149,971 | 678,353 | 678,352 | 0 | 0 | 2024-04-16 | 2024-04-30 |
| 05-2024_01.csv | 180,735,498 | 573,584 | 573,583 | 0 | 0 | 2024-05-02 | 2024-05-15 |
| 05-2024_02.csv | 219,948,491 | 699,080 | 699,079 | 0 | 0 | 2024-05-16 | 2024-05-31 |
| 06-2024_01.csv | 186,569,870 | 592,208 | 592,207 | 0 | 0 | 2024-06-03 | 2024-06-14 |
| 06-2024_02.csv | 193,695,585 | 614,526 | 614,525 | 0 | 0 | 2024-06-17 | 2024-06-28 |
| 07-2024_01.csv | 218,916,978 | 693,723 | 693,722 | 0 | 0 | 2024-07-01 | 2024-07-15 |
| 07-2024_02.csv | 246,612,173 | 781,175 | 781,174 | 0 | 0 | 2024-07-16 | 2024-07-31 |
| 08-2024_01.csv | 231,838,905 | 735,697 | 735,696 | 0 | 0 | 2024-08-01 | 2024-08-15 |
| 08-2024_02.csv | 220,613,991 | 700,638 | 700,637 | 0 | 0 | 2024-08-16 | 2024-08-30 |
| 09-2024_01.csv | 177,501,265 | 561,264 | 561,263 | 0 | 0 | 2024-09-02 | 2024-09-13 |
| 09-2024_02.csv | 173,041,838 | 549,216 | 549,215 | 0 | 0 | 2024-09-17 | 2024-09-30 |
| 10-2024_01.csv | 177,851,350 | 564,007 | 564,006 | 0 | 0 | 2024-10-02 | 2024-10-15 |
| 10-2024_02.csv | 206,867,532 | 656,834 | 656,833 | 0 | 0 | 2024-10-16 | 2024-10-31 |
| 11-2024_01.csv | 185,655,454 | 586,210 | 586,209 | 0 | 0 | 2024-11-01 | 2024-11-15 |
| 11-2024_02.csv | 188,284,827 | 591,279 | 591,278 | 0 | 0 | 2024-11-16 | 2024-11-29 |
| 12-2024_01.csv | 173,268,465 | 549,253 | 549,252 | 0 | 0 | 2024-12-02 | 2024-12-13 |
| 12-2024_02.csv | 156,702,371 | 497,087 | 497,086 | 0 | 0 | 2024-12-16 | 2024-12-30 |
| 01-2025_01.csv | 172,924,393 | 547,783 | 547,782 | 0 | 0 | 2025-01-02 | 2025-01-15 |
| 01-2025_02.csv | 218,837,634 | 696,183 | 696,182 | 0 | 0 | 2025-01-16 | 2025-01-31 |
| 02-2025_01.csv | 168,933,779 | 536,079 | 536,078 | 0 | 0 | 2025-02-04 | 2025-02-14 |
| 02-2025_02.csv | 190,830,861 | 606,777 | 606,776 | 0 | 0 | 2025-02-17 | 2025-02-28 |
| 03-2025_01.csv | 184,284,322 | 587,059 | 587,058 | 0 | 0 | 2025-03-03 | 2025-03-14 |
| 03-2025_02.csv | 186,930,413 | 594,321 | 594,320 | 0 | 0 | 2025-03-18 | 2025-03-31 |
| 04-2025_01.csv | 186,622,200 | 594,888 | 594,887 | 0 | 0 | 2025-04-01 | 2025-04-15 |
| 04-2025_02.csv | 135,336,065 | 430,974 | 430,973 | 0 | 0 | 2025-04-16 | 2025-04-30 |
| 05-2025_01.csv | 137,312,752 | 437,834 | 437,833 | 0 | 0 | 2025-05-02 | 2025-05-15 |
| 05-2025_02.csv | 165,182,068 | 526,970 | 526,969 | 0 | 0 | 2025-05-16 | 2025-05-30 |
| 06-2025_01.csv | 168,286,465 | 534,835 | 534,834 | 0 | 0 | 2025-06-02 | 2025-06-13 |
| 06-2025_02.csv | 180,845,643 | 573,714 | 573,713 | 0 | 0 | 2025-06-16 | 2025-06-30 |
| 07-2025_01.csv | 209,028,535 | 664,947 | 664,946 | 0 | 0 | 2025-07-01 | 2025-07-15 |
| 07-2025_02.csv | 224,540,038 | 715,247 | 715,246 | 0 | 0 | 2025-07-16 | 2025-07-31 |
| 08-2025_01.csv | 217,515,249 | 692,570 | 692,569 | 0 | 0 | 2025-08-01 | 2025-08-15 |
| 08-2025_02.csv | 153,500,849 | 489,956 | 489,955 | 0 | 0 | 2025-08-18 | 2025-08-29 |
| 09-2025_01.csv | 170,596,551 | 542,923 | 542,922 | 0 | 0 | 2025-09-01 | 2025-09-15 |
| 09-2025_02.csv | 166,954,073 | 534,093 | 534,092 | 0 | 0 | 2025-09-17 | 2025-09-30 |
| 10-2025_01.csv | 166,396,956 | 528,493 | 528,492 | 0 | 0 | 2025-10-01 | 2025-10-15 |
| 10-2025_02.csv | 186,736,618 | 595,063 | 595,062 | 0 | 0 | 2025-10-16 | 2025-10-31 |
| 11-2025_01.csv | 153,833,762 | 486,183 | 486,182 | 0 | 0 | 2025-11-03 | 2025-11-15 |
| 11-2025_02.csv | 138,876,694 | 437,364 | 437,363 | 0 | 0 | 2025-11-16 | 2025-11-29 |
| 12-2025_01.csv | 151,927,226 | 481,459 | 481,458 | 0 | 0 | 2025-12-01 | 2025-12-15 |
| 12-2025_02.csv | 136,715,932 | 434,833 | 434,832 | 0 | 0 | 2025-12-16 | 2025-12-31 |
| 01-2026_Q1.csv | 156,465,328 | 496,727 | 496,726 | 0 | 0 | 2026-01-02 | 2026-01-15 |
| 01-2026_Q2.csv | 198,482,949 | 630,372 | 630,371 | 0 | 0 | 2026-01-16 | 2026-01-30 |
| 02-2026_Q1.csv | 156,476,250 | 496,646 | 496,645 | 0 | 0 | 2026-02-03 | 2026-02-13 |
| 02-2026_Q2.csv | 172,814,585 | 547,812 | 547,811 | 0 | 0 | 2026-02-16 | 2026-02-27 |
| 03-2026_Q1.csv | 176,171,184 | 558,147 | 558,146 | 0 | 0 | 2026-03-02 | 2026-03-13 |
| 03-2026_Q2.csv | 190,545,396 | 606,523 | 606,522 | 0 | 0 | 2026-03-17 | 2026-03-31 |
| 04-2026_Q1.csv | 151,791,389 | 480,843 | 480,842 | 0 | 0 | 2026-04-01 | 2026-04-15 |
| 04-2026_Q2.csv | 187,803,269 | 597,635 | 597,634 | 0 | 0 | 2026-04-16 | 2026-04-30 |
| 05-2026_Q1.csv | 150,920,462 | 480,011 | 480,010 | 0 | 0 | 2026-05-01 | 2026-05-15 |
| 05-2026_Q2.csv | 167,575,168 | 533,215 | 533,214 | 0 | 0 | 2026-05-16 | 2026-05-31 |
| 06-2026_Q1.csv | 196,098,295 | 622,765 | 622,764 | 0 | 0 | 2026-06-01 | 2026-06-15 |
| 06-2026_Q2.csv | 199,900,196 | 633,687 | 633,686 | 0 | 0 | 2026-06-16 | 2026-06-30 |
| 07-2026_Q1.csv | 210,414,850 | 665,910 | 665,909 | 0 | 0 | 2026-07-01 | 2026-07-15 |
| 07-2026_Q2.csv | 234,044,542 | 744,624 | 744,623 | 0 | 0 | 2026-07-16 | 2026-07-31 |
| 08-2026_Q1.csv | 200,612,182 | 635,953 | 635,952 | 0 | 0 | 2026-08-01 | 2026-08-15 |
| 08-2026_Q2.csv | 206,979,176 | 657,603 | 657,602 | 0 | 0 | 2026-08-16 | 2026-08-31 |
| **total** | 11,827,015,231 | 37,548,264 | 37,548,200 | 0 | 0 | 2024-01-02 | 2026-08-31 |

Registros leídos = guardados + rechazados: 37,548,200 = 37,548,200 + 0.

**Formatos de fecha reconocidos** (el diccionario solo dice «Datetime»): `2006/01/02` = 37,548,200 registros;

**Rechazados con motivo:** ninguno.

### Totales

- Registros: 37,548,200 · productos distintos: 820 · presentaciones distintas (producto + presentación + marca): 7,960 · sucursales distintas (nombre comercial + dirección + municipio): 3,504
- Cadenas comerciales distintas (`cadena_comercial`): 269
- Meses cubiertos: 32 (2024-01 a 2026-08)
- Fecha mínima: 2024-01-02 · fecha máxima: 2026-08-31

### Cobertura por `cadena_comercial` (todas, de mayor a menor)

| cadena_comercial | registros | productos | presentaciones | sucursales | meses | primera | última |
|---|---:|---:|---:|---:|---|---|---|
| Wal-mart ★ | 6,111,448 | 778 | 5,966 | 92 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Hipermercado Soriana | 5,387,519 | 796 | 5,592 | 80 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Bodega Aurrera ★ | 4,065,188 | 753 | 4,814 | 89 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Chedraui | 2,768,451 | 792 | 5,565 | 43 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Mega Soriana | 2,502,644 | 786 | 5,233 | 38 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Wal-mart Express ★ | 1,911,475 | 663 | 3,070 | 41 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Farmacia Guadalajara | 1,336,064 | 541 | 2,090 | 71 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Mercado Soriana | 906,589 | 745 | 3,783 | 20 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| La Comer | 855,593 | 771 | 4,925 | 12 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-31 |
| Soriana Super | 666,401 | 704 | 3,115 | 12 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Farmacias del Ahorro | 621,789 | 515 | 1,416 | 50 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Benavides | 621,297 | 446 | 1,414 | 45 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Coppel | 600,346 | 221 | 1,725 | 94 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Mercado Publico | 599,917 | 218 | 891 | 128 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Ley | 588,162 | 732 | 3,557 | 14 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-31 |
| H.e.b. | 545,287 | 755 | 3,984 | 10 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Chedraui Selecto | 392,818 | 762 | 4,719 | 7 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| Fresko la Comer | 363,622 | 750 | 3,724 | 6 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Superissste | 357,992 | 193 | 835 | 23 | 28 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-31 |
| Farmacia | 338,212 | 383 | 1,068 | 36 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Liverpool | 329,776 | 60 | 2,173 | 54 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Central de Abastos | 315,004 | 209 | 968 | 433 | 28 (2024-01 a 2026-04) | 2024-01-02 | 2026-04-30 |
| Oxxo | 303,615 | 350 | 992 | 47 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Sumesa | 298,372 | 665 | 2,667 | 5 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Elektra | 276,458 | 257 | 1,246 | 80 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Sanborn S Hnos. | 271,663 | 386 | 1,583 | 24 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| S Mart | 221,067 | 582 | 2,225 | 5 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-27 |
| Super Chedraui | 217,080 | 696 | 2,953 | 6 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Sears Roebuck de Mexico | 207,590 | 193 | 1,964 | 34 | 31 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Farmacia Yza | 190,640 | 423 | 1,189 | 15 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| City Market | 176,695 | 664 | 2,544 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-31 |
| Bodega Aurrera Express ★ | 165,997 | 313 | 977 | 9 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-31 |
| Office Depot | 158,008 | 47 | 678 | 34 | 28 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-28 |
| Farmacias de Similares | 137,892 | 310 | 613 | 38 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| Minisuper | 137,217 | 213 | 1,085 | 11 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Soriana Express | 127,662 | 466 | 1,738 | 5 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| Papeleria Tony | 99,777 | 162 | 688 | 19 | 28 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-28 |
| Office Max | 97,134 | 39 | 566 | 25 | 28 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-28 |
| Alsuper | 95,758 | 611 | 2,074 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-12 |
| Alsuper Fresh Market | 95,463 | 594 | 2,046 | 1 | 31 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Farmacia la Union | 87,731 | 380 | 1,041 | 6 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-31 |
| Super Aki | 79,662 | 274 | 1,459 | 4 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-28 |
| Tortillerias Tradicionales | 72,583 | 133 | 465 | 1,052 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Supermercado Gonzalez | 65,691 | 517 | 1,608 | 1 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-17 |
| Farmatodo | 63,756 | 369 | 966 | 6 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Aprecio | 63,393 | 497 | 1,483 | 2 | 30 (2024-03 a 2026-08) | 2024-03-26 | 2026-08-26 |
| Farmacia San Pablo | 51,869 | 457 | 1,292 | 3 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Super Aki Xtra | 50,126 | 437 | 1,535 | 2 | 31 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| 7 Eleven | 49,551 | 104 | 379 | 10 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-28 |
| Farmacia Bazar | 49,500 | 369 | 985 | 4 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| U.n.a.m. | 47,360 | 560 | 1,982 | 2 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Central de Abasto | 45,901 | 185 | 561 | 180 | 4 (2026-05 a 2026-08) | 2026-05-01 | 2026-08-31 |
| Papeleria Tradicional | 39,857 | 34 | 538 | 24 | 26 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-28 |
| Farmacias Similares | 39,814 | 247 | 480 | 9 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-28 |
| Abastecedora Lumen | 39,040 | 36 | 434 | 8 | 28 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-28 |
| Mi Tienda del Ahorro | 37,353 | 295 | 1,194 | 1 | 24 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Salinas y Rocha | 27,919 | 28 | 491 | 7 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| Papeleria Ofix | 26,913 | 36 | 477 | 6 | 26 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-28 |
| Neto (super Precio) | 26,123 | 79 | 237 | 6 | 24 (2024-01 a 2025-12) | 2024-01-03 | 2025-12-29 |
| Panaderias Tradicionales | 25,428 | 5 | 16 | 93 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Ley Express | 25,283 | 179 | 723 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Merkabastos | 25,245 | 208 | 673 | 2 | 26 (2024-01 a 2026-02) | 2024-01-02 | 2026-02-10 |
| Operadora Merco | 24,805 | 202 | 927 | 1 | 24 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Farmacia Especializada | 24,195 | 304 | 701 | 3 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| El Pitico | 22,815 | 157 | 623 | 1 | 31 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Farmacia Santa Cecilia | 22,153 | 349 | 849 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Farmacia Isseg | 21,836 | 371 | 896 | 3 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Farmacia Pey Pharma | 21,155 | 379 | 942 | 2 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Abarrotera Sanchez | 20,933 | 112 | 436 | 1 | 31 (2024-01 a 2026-08) | 2024-01-12 | 2026-08-28 |
| Farmacia Popular Roman | 20,924 | 248 | 610 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Farmacia Comercio | 19,886 | 339 | 822 | 1 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-27 |
| Farmacia Paris | 19,784 | 346 | 821 | 1 | 32 (2024-01 a 2026-08) | 2024-01-10 | 2026-08-25 |
| Palacio de Hierro | 19,473 | 53 | 1,206 | 5 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-31 |
| Farmacia San Miguel | 19,369 | 336 | 747 | 2 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Farmacias Especializadas | 19,223 | 307 | 720 | 3 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Abarrotera Monterrey | 18,480 | 97 | 373 | 1 | 29 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-28 |
| Tiendas Garis | 17,830 | 266 | 907 | 3 | 14 (2024-01 a 2026-01) | 2024-01-04 | 2026-01-16 |
| Mega Elektra | 16,956 | 31 | 515 | 5 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-31 |
| Farmacia Torre Medica | 16,614 | 270 | 659 | 2 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Mercado Emilio Sanchez Piedras | 16,477 | 158 | 478 | 3 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Farmacia Universitaria Alexander Flemin | 16,353 | 344 | 804 | 1 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-25 |
| Farmacia la Luz | 16,023 | 321 | 717 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Tiendas de Descuento Arteli | 15,353 | 483 | 1,032 | 1 | 8 (2024-01 a 2024-08) | 2024-01-08 | 2024-08-12 |
| Farmacia el Sol | 14,995 | 298 | 658 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| Grand Marchand | 14,766 | 35 | 403 | 2 | 28 (2024-01 a 2026-08) | 2024-01-25 | 2026-08-26 |
| Ultrafarma | 14,568 | 371 | 942 | 1 | 25 (2024-07 a 2026-08) | 2024-07-03 | 2026-08-26 |
| Farmacia Abc | 14,505 | 279 | 625 | 2 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Juguetron | 14,165 | 16 | 961 | 19 | 7 (2024-01 a 2026-01) | 2024-01-02 | 2026-01-02 |
| Farmacia San Isidro | 13,834 | 318 | 703 | 2 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Tiendas de Uniformes | 13,508 | 11 | 281 | 24 | 6 (2024-07 a 2026-08) | 2024-07-01 | 2026-08-28 |
| Pescaderias Tradicionales | 13,344 | 42 | 224 | 37 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-28 |
| Mueblerias | 12,691 | 31 | 431 | 6 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-28 |
| Alsuper Store | 12,669 | 203 | 605 | 1 | 8 (2024-01 a 2024-08) | 2024-01-12 | 2024-08-16 |
| Farmacia la Miniatura | 12,124 | 371 | 926 | 2 | 30 (2024-01 a 2026-08) | 2024-01-11 | 2026-08-13 |
| Del Sol | 12,051 | 58 | 970 | 6 | 29 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-25 |
| Gran Bodega | 11,822 | 256 | 763 | 1 | 8 (2024-01 a 2024-08) | 2024-01-08 | 2024-08-16 |
| Dunosusa | 11,793 | 82 | 260 | 1 | 30 (2024-01 a 2026-06) | 2024-01-02 | 2026-06-23 |
| Mercado Malibran | 11,661 | 89 | 299 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Mercado Miguel Hidalgo | 11,348 | 96 | 252 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Super Kompras | 10,707 | 276 | 923 | 1 | 14 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Farmacia San Juan | 10,154 | 191 | 409 | 1 | 28 (2024-01 a 2026-04) | 2024-01-04 | 2026-04-23 |
| Mercado Adolfo Lopez Mateos | 9,877 | 164 | 548 | 3 | 31 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-27 |
| Cimaco | 9,809 | 30 | 349 | 2 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Neto (Super Precio) | 9,620 | 62 | 193 | 5 | 8 (2026-01 a 2026-08) | 2026-01-05 | 2026-08-31 |
| Jugueteria y Papeleria Tony | 9,472 | 33 | 389 | 2 | 25 (2024-01 a 2026-08) | 2024-01-26 | 2026-08-28 |
| Electromuebles | 9,406 | 15 | 164 | 4 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Mercado Publico 5 de Mayo | 9,374 | 102 | 259 | 2 | 30 (2024-03 a 2026-08) | 2024-03-01 | 2026-08-28 |
| Mercado J. Jesus Teran | 9,247 | 153 | 294 | 1 | 21 (2024-01 a 2025-09) | 2024-01-08 | 2025-09-01 |
| Farmacia del Hospital | 9,140 | 234 | 459 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Departamental del Sol | 8,728 | 69 | 983 | 2 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Compucopias | 8,656 | 31 | 355 | 1 | 26 (2024-01 a 2026-08) | 2024-01-23 | 2026-08-27 |
| Papeleria del Ahorro | 8,620 | 32 | 266 | 3 | 26 (2024-01 a 2026-08) | 2024-01-24 | 2026-08-28 |
| Papeleria Dabo | 8,526 | 35 | 334 | 2 | 21 (2024-01 a 2026-02) | 2024-01-22 | 2026-02-16 |
| Farmacias Guelaguetza | 8,448 | 368 | 703 | 2 | 25 (2024-01 a 2026-02) | 2024-01-08 | 2026-02-03 |
| Papeleria Casa Marchand | 8,207 | 34 | 367 | 2 | 20 (2024-01 a 2026-08) | 2024-01-23 | 2026-08-25 |
| Mercado el Meson Estrella | 7,687 | 110 | 288 | 2 | 28 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-24 |
| Julio Cepeda Jugueterias | 7,647 | 16 | 644 | 10 | 6 (2024-01 a 2025-12) | 2024-01-02 | 2025-12-31 |
| Papeleria la Mariposa | 7,623 | 33 | 295 | 1 | 26 (2024-01 a 2026-08) | 2024-01-23 | 2026-08-25 |
| Papeleria Escolar | 7,396 | 34 | 385 | 1 | 24 (2024-01 a 2026-08) | 2024-01-24 | 2026-08-27 |
| Calimax | 7,351 | 499 | 1,011 | 1 | 3 (2024-01 a 2024-03) | 2024-01-02 | 2024-03-19 |
| Farmacia Caballero | 7,211 | 151 | 311 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-24 |
| Farmacia Emergencias | 7,065 | 161 | 338 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Papeleria Vazquez Hermanos y Compañia | 6,818 | 28 | 244 | 1 | 25 (2024-01 a 2026-08) | 2024-01-26 | 2026-08-24 |
| Hiper Lumen | 6,454 | 33 | 406 | 1 | 24 (2024-01 a 2026-08) | 2024-01-24 | 2026-08-26 |
| Librerias Gonvill | 6,334 | 30 | 208 | 1 | 25 (2024-01 a 2026-08) | 2024-01-26 | 2026-08-28 |
| Abarrotes Super Cabrera Class | 5,918 | 76 | 217 | 1 | 32 (2024-01 a 2026-08) | 2024-01-11 | 2026-08-27 |
| Mercado Central de Acapulco | 5,806 | 164 | 484 | 2 | 17 (2024-02 a 2026-08) | 2024-02-22 | 2026-08-26 |
| Sears Roebuck de México | 5,800 | 47 | 682 | 21 | 1 (2025-12 a 2025-12) | 2025-12-01 | 2025-12-30 |
| Farmacia y Drogueria Toluca | 5,761 | 307 | 556 | 2 | 10 (2024-02 a 2026-08) | 2024-02-13 | 2026-08-28 |
| Muebleria el Gran Salto | 5,750 | 22 | 171 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-06 |
| Electrodomesticos Chapur | 5,509 | 48 | 672 | 1 | 24 (2024-01 a 2026-05) | 2024-01-03 | 2026-05-20 |
| Farmacia Canto | 5,421 | 147 | 308 | 1 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Telebodega | 5,393 | 27 | 303 | 1 | 31 (2024-01 a 2026-08) | 2024-01-09 | 2026-08-25 |
| Ryse | 5,191 | 16 | 176 | 2 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Juguetibici | 4,821 | 14 | 653 | 2 | 5 (2024-11 a 2026-01) | 2024-11-04 | 2026-01-02 |
| Papeleria San Felipe | 4,732 | 32 | 296 | 1 | 21 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-24 |
| Ekar de Gas | 4,468 | 18 | 142 | 1 | 29 (2024-01 a 2026-05) | 2024-01-03 | 2026-05-20 |
| Papelerias Super Precios | 4,449 | 32 | 211 | 1 | 22 (2024-01 a 2026-08) | 2024-01-24 | 2026-08-26 |
| Distribuidora Queretana de Medicamentos | 3,835 | 182 | 214 | 1 | 12 (2024-01 a 2024-12) | 2024-01-04 | 2024-12-19 |
| Farmacia San Jose | 3,640 | 164 | 265 | 2 | 24 (2024-01 a 2025-12) | 2024-01-04 | 2025-12-18 |
| Papeleria Paquin | 3,423 | 30 | 192 | 1 | 22 (2024-01 a 2026-08) | 2024-01-25 | 2026-08-27 |
| Mercado del Mar Zapopan | 3,304 | 28 | 109 | 4 | 26 (2024-02 a 2026-08) | 2024-02-01 | 2026-08-27 |
| Papeleria del Estudiante | 3,189 | 29 | 167 | 1 | 24 (2024-01 a 2026-08) | 2024-01-25 | 2026-08-27 |
| Pescaderias | 3,185 | 37 | 184 | 6 | 30 (2024-01 a 2026-08) | 2024-01-12 | 2026-08-28 |
| Abarrotera Guadalupana (frutas) | 3,176 | 108 | 415 | 1 | 4 (2024-01 a 2024-04) | 2024-01-09 | 2024-04-23 |
| Hemsa | 3,153 | 36 | 317 | 1 | 22 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-24 |
| Farmacia Hospital Reforma | 2,920 | 227 | 489 | 2 | 31 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-25 |
| Farmacia Lopez | 2,820 | 265 | 310 | 1 | 5 (2026-04 a 2026-08) | 2026-04-06 | 2026-08-31 |
| Papeleria Saldaña | 2,661 | 32 | 174 | 1 | 22 (2024-01 a 2026-08) | 2024-01-22 | 2026-08-24 |
| Papeleria el Escritorio | 2,599 | 28 | 135 | 1 | 25 (2024-01 a 2026-08) | 2024-01-23 | 2026-08-28 |
| Papeleria Colibri | 2,558 | 27 | 193 | 1 | 12 (2024-01 a 2026-08) | 2024-01-26 | 2026-08-28 |
| Optivosa | 2,520 | 28 | 186 | 2 | 17 (2024-01 a 2026-08) | 2024-01-23 | 2026-08-25 |
| Vinaterias | 2,518 | 9 | 112 | 32 | 7 (2024-01 a 2025-12) | 2024-01-04 | 2025-12-24 |
| Mercado del Mar Felipe Angeles Higueril | 2,479 | 25 | 94 | 3 | 22 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Centro Nacional de Distribucion Jutega | 2,399 | 13 | 447 | 2 | 5 (2024-01 a 2025-12) | 2024-01-03 | 2025-12-17 |
| Game Planet | 2,395 | 2 | 96 | 20 | 7 (2024-01 a 2026-01) | 2024-01-02 | 2026-01-02 |
| Papeleria el Tintero | 2,327 | 25 | 184 | 1 | 22 (2024-01 a 2026-08) | 2024-01-25 | 2026-08-26 |
| Papeleria Dany | 2,311 | 23 | 163 | 1 | 19 (2024-01 a 2026-08) | 2024-01-25 | 2026-08-24 |
| Colchones y Muebles de Campeche | 2,181 | 14 | 70 | 1 | 31 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Farmapronto | 2,170 | 321 | 372 | 3 | 8 (2024-01 a 2024-12) | 2024-01-03 | 2024-12-19 |
| Ultra Hogar | 1,948 | 9 | 58 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-28 |
| Mercado la Cruz | 1,892 | 101 | 166 | 1 | 5 (2024-01 a 2024-05) | 2024-01-10 | 2024-05-08 |
| Mueblerias Portillo | 1,770 | 18 | 99 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Zapaterias | 1,734 | 2 | 78 | 15 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-28 |
| Distribuidora Escolar y de Oficina J.p. | 1,677 | 30 | 176 | 1 | 19 (2024-01 a 2026-08) | 2024-01-26 | 2026-08-26 |
| Uniformes Casa Cruz | 1,657 | 10 | 215 | 1 | 6 (2024-07 a 2026-08) | 2024-07-03 | 2026-08-26 |
| Surtidora Departamental | 1,635 | 32 | 198 | 1 | 28 (2024-01 a 2026-05) | 2024-01-04 | 2026-05-12 |
| Papeleria Snoopy | 1,575 | 29 | 197 | 1 | 23 (2024-02 a 2026-08) | 2024-02-23 | 2026-08-27 |
| Papeleria Tlaxacala | 1,543 | 29 | 211 | 2 | 17 (2024-04 a 2026-08) | 2024-04-19 | 2026-08-28 |
| Farmacias la Salud | 1,517 | 162 | 193 | 1 | 5 (2026-04 a 2026-08) | 2026-04-08 | 2026-08-26 |
| Woolworth | 1,516 | 23 | 488 | 4 | 7 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-26 |
| Papeleria Casa Sakura | 1,411 | 24 | 117 | 1 | 13 (2024-01 a 2026-08) | 2024-01-15 | 2026-08-24 |
| Mercado del Mar San Felipe | 1,356 | 23 | 88 | 2 | 29 (2024-01 a 2026-08) | 2024-01-29 | 2026-08-24 |
| Pescados y Mariscos Miguel Hidalgo | 1,315 | 20 | 77 | 4 | 31 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Novedades Olague S.A. de C.V. | 1,286 | 11 | 178 | 1 | 6 (2024-07 a 2026-08) | 2024-07-18 | 2026-08-27 |
| Papeleria Masan | 1,279 | 23 | 82 | 1 | 18 (2024-01 a 2025-10) | 2024-01-25 | 2025-10-23 |
| Pescaderia Puerto Aventura | 1,250 | 37 | 127 | 1 | 22 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Pescaderia Cepesmar | 1,173 | 30 | 131 | 1 | 31 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-24 |
| Uniformes el Artillero | 1,165 | 10 | 151 | 1 | 6 (2024-07 a 2026-08) | 2024-07-30 | 2026-08-25 |
| Papeleria el Lapiz | 1,117 | 31 | 128 | 1 | 25 (2024-01 a 2026-08) | 2024-01-26 | 2026-08-28 |
| Farmacia Cruz Blanca | 1,113 | 107 | 127 | 1 | 10 (2024-02 a 2025-08) | 2024-02-20 | 2025-08-19 |
| Papeleria Abc | 1,111 | 28 | 153 | 1 | 6 (2026-02 a 2026-08) | 2026-02-27 | 2026-08-14 |
| Uniformes Casa Rosy | 1,081 | 12 | 165 | 1 | 6 (2024-07 a 2026-08) | 2024-07-25 | 2026-08-27 |
| Uniformes Escolares Dany | 1,027 | 9 | 152 | 1 | 6 (2024-07 a 2026-08) | 2024-07-26 | 2026-08-28 |
| Milano | 1,013 | 10 | 103 | 4 | 6 (2024-07 a 2026-08) | 2024-07-01 | 2026-08-28 |
| Almacenes Rodriguez | 1,005 | 10 | 124 | 2 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-26 |
| Uniformes Acosta | 988 | 11 | 174 | 1 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-25 |
| Papeleria Delta | 957 | 18 | 69 | 1 | 24 (2024-01 a 2026-08) | 2024-01-24 | 2026-08-28 |
| Uniformes el Surtidor de Xochimilco | 950 | 11 | 217 | 1 | 6 (2024-07 a 2026-08) | 2024-07-16 | 2026-08-25 |
| Mercado Centralito | 918 | 81 | 129 | 1 | 4 (2026-03 a 2026-06) | 2026-03-30 | 2026-06-29 |
| Issstezac | 898 | 73 | 131 | 1 | 6 (2024-01 a 2024-06) | 2024-01-05 | 2024-06-21 |
| Uniformes la Favorita | 896 | 9 | 130 | 1 | 6 (2024-07 a 2026-08) | 2024-07-18 | 2026-08-26 |
| Uniformes Casa Tino | 892 | 8 | 98 | 1 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-25 |
| Pescaderia Sevilla Carro Verde | 865 | 20 | 109 | 2 | 31 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-24 |
| Pescaderia Guasaveña | 857 | 21 | 85 | 1 | 26 (2024-02 a 2026-08) | 2024-02-01 | 2026-08-27 |
| Uniformes Tejidos Keyri | 833 | 9 | 109 | 1 | 6 (2024-07 a 2026-08) | 2024-07-05 | 2026-08-26 |
| Pescaderia el Foro | 827 | 40 | 99 | 1 | 10 (2024-01 a 2024-11) | 2024-01-31 | 2024-11-27 |
| Pescaderia el Tejano de Queretaro | 825 | 33 | 122 | 2 | 24 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Pescaderia Bahia de Kino | 813 | 18 | 87 | 1 | 31 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Uniformes Erick | 809 | 10 | 127 | 1 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-26 |
| Uniformes Trofeo | 801 | 10 | 144 | 1 | 6 (2024-07 a 2026-08) | 2024-07-05 | 2026-08-21 |
| Uniformes Bekamex | 774 | 8 | 110 | 1 | 6 (2024-07 a 2026-08) | 2024-07-19 | 2026-08-28 |
| Loyep (uniformes) | 768 | 10 | 81 | 1 | 4 (2024-07 a 2025-08) | 2024-07-02 | 2025-08-20 |
| Muebleria Elizondo | 763 | 24 | 124 | 1 | 6 (2024-01 a 2024-10) | 2024-01-08 | 2024-10-14 |
| Pescaderia Yuridia | 750 | 25 | 104 | 2 | 24 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Uniformes Luce | 740 | 10 | 83 | 1 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-26 |
| Uniformes Roa | 727 | 11 | 211 | 2 | 6 (2024-07 a 2026-08) | 2024-07-01 | 2026-08-24 |
| Pescaderia Hermanos Luna | 724 | 22 | 80 | 1 | 29 (2024-01 a 2026-08) | 2024-01-29 | 2026-08-24 |
| Pescaderia el Mangle | 715 | 22 | 69 | 1 | 22 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Mariscos Papillon | 659 | 15 | 67 | 1 | 28 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Operadora Prissa | 639 | 8 | 77 | 6 | 5 (2024-01 a 2025-12) | 2024-01-02 | 2025-12-22 |
| Pescaderia Liz del Mar | 637 | 17 | 65 | 1 | 22 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| La Costa Marisquera | 622 | 18 | 80 | 1 | 28 (2024-01 a 2026-08) | 2024-01-29 | 2026-08-24 |
| La Europea | 621 | 5 | 50 | 9 | 4 (2024-11 a 2025-12) | 2024-11-19 | 2025-12-23 |
| Almacenes el Aguila | 614 | 10 | 133 | 1 | 6 (2024-07 a 2026-08) | 2024-07-16 | 2026-08-25 |
| Pescaderia la Marinera del 23 | 600 | 18 | 63 | 1 | 20 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Pescaderia las Truchas | 597 | 17 | 68 | 1 | 28 (2024-02 a 2026-08) | 2024-02-01 | 2026-08-27 |
| Uniforama | 588 | 8 | 145 | 1 | 3 (2024-08 a 2026-08) | 2024-08-01 | 2026-08-27 |
| Pescaderia Lindo Sinaloa | 583 | 19 | 58 | 1 | 22 (2024-02 a 2026-08) | 2024-02-02 | 2026-08-28 |
| Pescaderia Mari-mich | 578 | 17 | 71 | 1 | 21 (2024-02 a 2026-06) | 2024-02-02 | 2026-06-26 |
| Pescaderia Hermanos Villegas | 555 | 13 | 74 | 2 | 21 (2024-01 a 2026-08) | 2024-01-30 | 2026-08-24 |
| La Colonial | 519 | 13 | 128 | 1 | 5 (2024-01 a 2025-12) | 2024-01-03 | 2025-12-17 |
| Pescados y Marscos Estrella | 492 | 24 | 70 | 1 | 10 (2024-01 a 2024-11) | 2024-01-31 | 2024-11-27 |
| Loyep (Uniformes) | 489 | 11 | 81 | 1 | 2 (2026-07 a 2026-08) | 2026-07-08 | 2026-08-26 |
| Gamers | 488 | 2 | 79 | 3 | 6 (2024-01 a 2025-12) | 2024-01-02 | 2025-12-29 |
| Muebleria Foto Contino | 470 | 5 | 35 | 1 | 8 (2024-01 a 2024-08) | 2024-01-05 | 2024-08-23 |
| Santovalle | 469 | 86 | 137 | 1 | 2 (2024-01 a 2024-02) | 2024-01-03 | 2024-02-07 |
| Uniformes el General | 457 | 7 | 99 | 3 | 6 (2024-07 a 2026-08) | 2024-07-03 | 2026-08-21 |
| Farmacia Cuernavaca | 441 | 267 | 305 | 1 | 2 (2024-01 a 2024-02) | 2024-01-03 | 2024-02-28 |
| Pescaderia Boca del Rio | 435 | 18 | 62 | 1 | 16 (2024-01 a 2026-08) | 2024-01-31 | 2026-08-26 |
| Farmacon | 428 | 331 | 404 | 1 | 1 (2025-09 a 2025-09) | 2025-09-22 | 2025-09-29 |
| Uniformes Casa Castilla | 413 | 7 | 89 | 1 | 6 (2024-07 a 2026-08) | 2024-07-18 | 2026-08-27 |
| Pescaderia Alex | 405 | 16 | 69 | 1 | 17 (2024-01 a 2026-08) | 2024-01-30 | 2026-08-25 |
| Uniformes Novedades Gler | 401 | 10 | 70 | 1 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-28 |
| Uniformes Escolares y Deportivos | 390 | 11 | 131 | 1 | 5 (2024-07 a 2026-08) | 2024-07-05 | 2026-08-28 |
| Muebleria el Importador | 385 | 5 | 25 | 1 | 8 (2024-01 a 2024-08) | 2024-01-04 | 2024-08-15 |
| Zapateria la Perla | 328 | 1 | 38 | 1 | 6 (2024-07 a 2026-08) | 2024-07-05 | 2026-08-28 |
| Pescaderia la Jaiba | 313 | 17 | 54 | 1 | 27 (2024-01 a 2026-08) | 2024-01-29 | 2026-08-24 |
| Vinateria la Playa | 293 | 4 | 53 | 3 | 4 (2024-11 a 2025-12) | 2024-11-21 | 2025-12-22 |
| Pescaderia Gambero | 257 | 12 | 30 | 1 | 21 (2024-01 a 2026-04) | 2024-01-31 | 2026-04-29 |
| Zapateria Pakar | 252 | 1 | 55 | 2 | 6 (2024-07 a 2026-08) | 2024-07-05 | 2026-08-26 |
| Heinss Hardy | 248 | 4 | 53 | 2 | 3 (2024-12 a 2025-12) | 2024-12-19 | 2025-12-23 |
| Pescaderia Juarez | 236 | 9 | 46 | 2 | 27 (2024-01 a 2026-07) | 2024-01-29 | 2026-07-27 |
| Zapateria Coppel Canada | 215 | 1 | 27 | 5 | 7 (2024-07 a 2026-08) | 2024-07-03 | 2026-08-27 |
| Farmacia Cruz Rosa | 204 | 115 | 127 | 1 | 1 (2024-01 a 2024-01) | 2024-01-05 | 2024-01-19 |
| Mundo de los Mariscos | 204 | 15 | 47 | 1 | 19 (2024-02 a 2026-08) | 2024-02-16 | 2026-08-28 |
| La Divina | 180 | 4 | 57 | 3 | 4 (2024-11 a 2025-12) | 2024-11-21 | 2025-12-19 |
| Cristaleria la Violeta | 171 | 10 | 60 | 1 | 5 (2024-01 a 2025-12) | 2024-01-02 | 2025-12-30 |
| Zapateria la Luna | 163 | 1 | 26 | 2 | 6 (2024-07 a 2026-08) | 2024-07-02 | 2026-08-25 |
| Jugueterías | 143 | 2 | 34 | 2 | 5 (2024-01 a 2025-12) | 2024-01-02 | 2025-12-30 |
| Zapateria Estrella | 132 | 1 | 13 | 1 | 6 (2024-07 a 2026-08) | 2024-07-11 | 2026-08-27 |
| Pescaderia Delfin | 123 | 8 | 18 | 1 | 4 (2024-01 a 2024-04) | 2024-01-31 | 2024-04-24 |
| Uniformes Marusa | 112 | 7 | 51 | 1 | 2 (2026-07 a 2026-08) | 2026-07-30 | 2026-08-27 |
| Comercial Universitaria | 95 | 14 | 30 | 1 | 5 (2026-03 a 2026-08) | 2026-03-20 | 2026-08-14 |
| Papeleria Farah | 92 | 17 | 92 | 1 | 1 (2024-01 a 2024-01) | 2024-01-24 | 2024-01-24 |
| Jugueteria Juguetimundo | 85 | 9 | 33 | 1 | 3 (2024-11 a 2025-11) | 2024-11-07 | 2025-11-21 |
| Uniformes el Leon | 70 | 2 | 13 | 1 | 5 (2024-07 a 2025-08) | 2024-07-01 | 2025-08-11 |
| Zapateria el Globito | 66 | 1 | 14 | 1 | 5 (2024-07 a 2026-07) | 2024-07-12 | 2026-07-31 |
| Pescados y Mariscos Miramar | 59 | 8 | 14 | 1 | 5 (2024-01 a 2024-05) | 2024-01-29 | 2024-05-27 |
| Daniel S Uniformes | 50 | 1 | 10 | 1 | 2 (2026-07 a 2026-08) | 2026-07-15 | 2026-08-27 |
| San Francisco de Asis (Ropa y Telas) | 49 | 7 | 21 | 1 | 2 (2026-07 a 2026-08) | 2026-07-10 | 2026-08-28 |
| Vinoteca | 40 | 3 | 16 | 2 | 2 (2024-12 a 2025-12) | 2024-12-13 | 2025-12-16 |
| San Francisco de Asis (ropa y Telas) | 39 | 5 | 21 | 1 | 4 (2024-07 a 2025-08) | 2024-07-31 | 2025-08-08 |
| La Castellana | 38 | 2 | 19 | 1 | 2 (2024-11 a 2024-12) | 2024-11-22 | 2024-12-06 |
| Pescaderia Aldamari | 32 | 11 | 17 | 1 | 8 (2024-01 a 2024-10) | 2024-01-12 | 2024-10-18 |
| Mixup | 23 | 2 | 12 | 2 | 4 (2024-11 a 2025-12) | 2024-11-07 | 2025-12-11 |
| Zapateria Destroyer | 22 | 1 | 11 | 1 | 6 (2024-07 a 2026-08) | 2024-07-31 | 2026-08-19 |
| Super Papelera, S.A. de C.V. | 12 | 5 | 12 | 1 | 1 (2024-07 a 2024-07) | 2024-07-31 | 2024-07-31 |
| Zapateria 3 Hermanos | 3 | 1 | 3 | 1 | 1 (2026-07 a 2026-07) | 2026-07-07 | 2026-07-07 |

★ = cadenas de interés (familia Walmart, familia Bodega Aurrera, Sam's Club).

### Cadenas de interés por estado y municipio

Las sucursales cercanas a Axel no se conocen: se reporta por municipio y la selección queda en Ajustes › Tiendas y precios.

| estado | municipio | cadena_comercial | registros | productos | sucursales | meses | primera | última |
|---|---|---|---:|---:|---:|---|---|---|
| Aguascalientes | Aguascalientes | Bodega Aurrera | 97,873 | 648 | 1 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Aguascalientes | Aguascalientes | Wal-mart | 104,044 | 706 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Baja California | Tijuana | Bodega Aurrera | 2,538 | 24 | 1 | 16 (2025-05 a 2026-08) | 2025-05-20 | 2026-08-25 |
| Baja California | Tijuana | Wal-mart | 108,866 | 673 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Baja California Sur | La Paz | Bodega Aurrera | 89,053 | 570 | 2 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Baja California Sur | La Paz | Wal-mart | 118,903 | 662 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Campeche | Campeche | Bodega Aurrera | 94,755 | 589 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Campeche | Campeche | Wal-mart | 120,301 | 689 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Chiapas | Tuxtla Gutiérrez | Bodega Aurrera | 82,515 | 614 | 1 | 32 (2024-01 a 2026-08) | 2024-01-30 | 2026-08-25 |
| Chiapas | Tuxtla Gutiérrez | Wal-mart | 101,319 | 674 | 2 | 32 (2024-01 a 2026-08) | 2024-01-31 | 2026-08-26 |
| Chihuahua | Chihuahua | Bodega Aurrera | 33,598 | 202 | 1 | 25 (2024-08 a 2026-08) | 2024-08-21 | 2026-08-20 |
| Chihuahua | Chihuahua | Wal-mart | 57,802 | 235 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-31 |
| Chihuahua | Juárez | Bodega Aurrera | 106,042 | 644 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Chihuahua | Juárez | Wal-mart | 119,584 | 675 | 1 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Ciudad de Mexico | Azcapotzalco | Bodega Aurrera Express | 10,000 | 127 | 1 | 21 (2024-01 a 2025-10) | 2024-01-11 | 2025-10-30 |
| Ciudad de Mexico | Azcapotzalco | Wal-mart | 144,454 | 711 | 2 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-20 |
| Ciudad de Mexico | Benito Juárez | Bodega Aurrera Express | 13,337 | 133 | 1 | 22 (2024-01 a 2025-11) | 2024-01-08 | 2025-11-10 |
| Ciudad de Mexico | Benito Juárez | Wal-mart | 135,503 | 711 | 2 | 23 (2024-01 a 2025-11) | 2024-01-08 | 2025-11-24 |
| Ciudad de Mexico | Benito Juárez | Wal-mart Express | 198,631 | 626 | 4 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-24 |
| Ciudad de Mexico | Coyoacán | Bodega Aurrera | 58,448 | 646 | 1 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-27 |
| Ciudad de Mexico | Coyoacán | Wal-mart | 183,118 | 708 | 3 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-27 |
| Ciudad de Mexico | Coyoacán | Wal-mart Express | 121,014 | 617 | 3 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-26 |
| Ciudad de Mexico | Cuauhtémoc | Bodega Aurrera | 65,745 | 642 | 2 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-28 |
| Ciudad de Mexico | Cuauhtémoc | Bodega Aurrera Express | 17,353 | 134 | 1 | 21 (2024-01 a 2025-10) | 2024-01-04 | 2025-10-30 |
| Ciudad de Mexico | Cuauhtémoc | Wal-mart | 73,981 | 700 | 1 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-24 |
| Ciudad de Mexico | Cuauhtémoc | Wal-mart Express | 163,780 | 625 | 3 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-28 |
| Ciudad de Mexico | Gustavo a. Madero | Bodega Aurrera | 67,461 | 635 | 1 | 23 (2024-01 a 2025-11) | 2024-01-12 | 2025-11-28 |
| Ciudad de Mexico | Gustavo a. Madero | Bodega Aurrera Express | 19,348 | 137 | 1 | 22 (2024-01 a 2025-10) | 2024-01-05 | 2025-10-31 |
| Ciudad de Mexico | Gustavo a. Madero | Wal-mart | 80,852 | 684 | 1 | 23 (2024-01 a 2025-11) | 2024-01-12 | 2025-11-28 |
| Ciudad de Mexico | Gustavo a. Madero | Wal-mart Express | 60,758 | 584 | 1 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-07 |
| Ciudad de Mexico | Iztacalco | Bodega Aurrera | 92,539 | 659 | 2 | 23 (2024-01 a 2025-11) | 2024-01-09 | 2025-11-26 |
| Ciudad de Mexico | Iztapalapa | Bodega Aurrera | 121,727 | 673 | 2 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-26 |
| Ciudad de Mexico | Iztapalapa | Wal-mart | 78,680 | 700 | 1 | 23 (2024-01 a 2025-11) | 2024-01-09 | 2025-11-26 |
| Ciudad de Mexico | Iztapalapa | Wal-mart Express | 46,245 | 599 | 1 | 23 (2024-01 a 2025-11) | 2024-01-09 | 2025-11-26 |
| Ciudad de Mexico | Magdalena Contreras | Wal-mart Express | 46,430 | 582 | 1 | 23 (2024-01 a 2025-11) | 2024-01-12 | 2025-11-10 |
| Ciudad de Mexico | Miguel Hidalgo | Bodega Aurrera | 124,511 | 651 | 2 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-28 |
| Ciudad de Mexico | Miguel Hidalgo | Wal-mart | 98,961 | 709 | 1 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-24 |
| Ciudad de Mexico | Miguel Hidalgo | Wal-mart Express | 154,278 | 615 | 3 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-25 |
| Ciudad de Mexico | Tlalpan | Bodega Aurrera | 33,831 | 598 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Ciudad de Mexico | Tlalpan | Wal-mart | 70,947 | 678 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Ciudad de Mexico | Tlalpan | Wal-mart Express | 100,681 | 604 | 2 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Ciudad de Mexico | Venustiano Carranza | Bodega Aurrera | 57,798 | 650 | 1 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-24 |
| Ciudad de Mexico | Venustiano Carranza | Wal-mart | 147,795 | 711 | 2 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-24 |
| Ciudad de Mexico | Xochimilco | Bodega Aurrera | 46,344 | 599 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-11 |
| Ciudad de Mexico | Álvaro Obregón | Bodega Aurrera | 77,380 | 626 | 1 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-24 |
| Ciudad de Mexico | Álvaro Obregón | Wal-mart Express | 59,011 | 592 | 1 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-10 |
| Ciudad de México | Azcapotzalco | Bodega Aurrera Express | 2,872 | 118 | 1 | 7 (2026-01 a 2026-08) | 2026-01-15 | 2026-08-20 |
| Ciudad de México | Azcapotzalco | Wal-mart | 54,992 | 659 | 2 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-27 |
| Ciudad de México | Benito Juárez | Bodega Aurrera Express | 3,505 | 104 | 1 | 7 (2025-12 a 2026-08) | 2025-12-08 | 2026-08-31 |
| Ciudad de México | Benito Juárez | Wal-mart | 60,647 | 665 | 2 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Ciudad de México | Benito Juárez | Wal-mart Express | 72,712 | 585 | 6 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Ciudad de México | Coyoacán | Bodega Aurrera | 25,558 | 592 | 1 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-27 |
| Ciudad de México | Coyoacán | Wal-mart | 72,642 | 677 | 4 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-27 |
| Ciudad de México | Coyoacán | Wal-mart Express | 27,996 | 579 | 3 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-25 |
| Ciudad de México | Cuauhtémoc | Bodega Aurrera | 25,441 | 589 | 1 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-31 |
| Ciudad de México | Cuauhtémoc | Bodega Aurrera Express | 4,209 | 121 | 1 | 8 (2026-01 a 2026-08) | 2026-01-08 | 2026-08-27 |
| Ciudad de México | Cuauhtémoc | Wal-mart | 20,519 | 625 | 1 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-13 |
| Ciudad de México | Cuauhtémoc | Wal-mart Express | 51,868 | 578 | 5 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Ciudad de México | Gustavo a. Madero | Bodega Aurrera | 23,158 | 587 | 1 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Ciudad de México | Gustavo a. Madero | Bodega Aurrera Express | 4,461 | 118 | 1 | 8 (2026-01 a 2026-08) | 2026-01-16 | 2026-08-21 |
| Ciudad de México | Gustavo a. Madero | Wal-mart | 27,722 | 635 | 1 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Ciudad de México | Gustavo a. Madero | Wal-mart Express | 20,575 | 547 | 1 | 8 (2026-01 a 2026-08) | 2026-01-02 | 2026-08-28 |
| Ciudad de México | Iztacalco | Bodega Aurrera | 25,677 | 601 | 3 | 9 (2025-12 a 2026-08) | 2025-12-10 | 2026-08-26 |
| Ciudad de México | Iztapalapa | Bodega Aurrera | 51,853 | 623 | 3 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-25 |
| Ciudad de México | Iztapalapa | Wal-mart | 35,874 | 647 | 2 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-26 |
| Ciudad de México | Iztapalapa | Wal-mart Express | 19,373 | 522 | 2 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-26 |
| Ciudad de México | Magdalena Contreras | Wal-mart Express | 7,512 | 548 | 2 | 8 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Ciudad de México | Miguel Hidalgo | Bodega Aurrera | 45,149 | 606 | 2 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Ciudad de México | Miguel Hidalgo | Wal-mart | 37,653 | 654 | 2 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-25 |
| Ciudad de México | Miguel Hidalgo | Wal-mart Express | 52,779 | 585 | 4 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Ciudad de México | Tlalpan | Bodega Aurrera | 8,801 | 548 | 1 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-18 |
| Ciudad de México | Tlalpan | Wal-mart | 20,715 | 613 | 1 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-25 |
| Ciudad de México | Tlalpan | Wal-mart Express | 23,875 | 564 | 3 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-18 |
| Ciudad de México | Venustiano Carranza | Bodega Aurrera | 23,121 | 576 | 1 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-27 |
| Ciudad de México | Venustiano Carranza | Wal-mart | 55,732 | 655 | 3 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-27 |
| Ciudad de México | Xochimilco | Bodega Aurrera | 15,498 | 574 | 1 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-25 |
| Ciudad de México | Álvaro Obregón | Bodega Aurrera | 30,588 | 583 | 3 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Ciudad de México | Álvaro Obregón | Wal-mart Express | 16,957 | 538 | 2 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-21 |
| Coahuila | Saltillo | Bodega Aurrera | 101,963 | 620 | 3 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Coahuila | Saltillo | Wal-mart | 113,743 | 665 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Durango | Durango | Wal-mart | 102,835 | 678 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-31 |
| Estado de Mexico | Atizapán de Zaragoza | Bodega Aurrera | 37,403 | 634 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-27 |
| Estado de Mexico | Coacalco | Wal-mart | 147,917 | 706 | 2 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-26 |
| Estado de Mexico | Cuautitlán | Bodega Aurrera | 62,353 | 648 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Estado de Mexico | Cuautitlán Izcalli | Wal-mart | 81,361 | 696 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Estado de Mexico | Ecatepec de Morelos | Bodega Aurrera | 104,068 | 669 | 2 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Estado de Mexico | Ecatepec de Morelos | Bodega Aurrera Express | 10,365 | 110 | 1 | 22 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-28 |
| Estado de Mexico | Ecatepec de Morelos | Wal-mart | 141,179 | 700 | 2 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Estado de Mexico | Huixquilucan | Wal-mart | 72,442 | 715 | 1 | 20 (2024-04 a 2025-11) | 2024-04-08 | 2025-11-25 |
| Estado de Mexico | Huixquilucan | Wal-mart Express | 52,389 | 602 | 1 | 23 (2024-01 a 2025-11) | 2024-01-08 | 2025-11-25 |
| Estado de Mexico | Metepec | Wal-mart | 9,614 | 561 | 1 | 6 (2024-01 a 2024-06) | 2024-01-08 | 2024-06-17 |
| Estado de Mexico | Naucalpan de Juárez | Bodega Aurrera | 40,957 | 634 | 1 | 23 (2024-01 a 2025-11) | 2024-01-08 | 2025-11-25 |
| Estado de Mexico | Naucalpan de Juárez | Wal-mart | 142,479 | 703 | 2 | 23 (2024-01 a 2025-11) | 2024-01-08 | 2025-11-28 |
| Estado de Mexico | Naucalpan de Juárez | Wal-mart Express | 109,807 | 613 | 2 | 23 (2024-01 a 2025-11) | 2024-01-08 | 2025-11-28 |
| Estado de Mexico | Nezahualcóyotl | Bodega Aurrera | 68,323 | 643 | 2 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-28 |
| Estado de Mexico | Tlalnepantla | Bodega Aurrera | 51,955 | 573 | 1 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-27 |
| Estado de Mexico | Tlalnepantla | Wal-mart | 169,968 | 709 | 2 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-27 |
| Estado de Mexico | Tlalnepantla | Wal-mart Express | 66,382 | 577 | 1 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-27 |
| Estado de Mexico | Tultitlán | Bodega Aurrera | 70,217 | 645 | 1 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-26 |
| Estado de Mexico | Tultitlán | Bodega Aurrera Express | 18,132 | 124 | 1 | 21 (2024-01 a 2025-10) | 2024-01-10 | 2025-10-29 |
| Estado de México | Atizapán de Zaragoza | Bodega Aurrera | 10,220 | 520 | 2 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-17 |
| Estado de México | Coacalco | Wal-mart | 42,305 | 663 | 4 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-19 |
| Estado de México | Cuautitlán | Bodega Aurrera | 21,611 | 585 | 2 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-18 |
| Estado de México | Cuautitlán Izcalli | Wal-mart | 33,883 | 647 | 2 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-25 |
| Estado de México | Ecatepec de Morelos | Bodega Aurrera | 44,743 | 621 | 4 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-26 |
| Estado de México | Ecatepec de Morelos | Bodega Aurrera Express | 2,073 | 103 | 1 | 8 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Estado de México | Ecatepec de Morelos | Wal-mart | 58,934 | 662 | 4 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-26 |
| Estado de México | Huixquilucan | Wal-mart | 34,877 | 673 | 2 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Estado de México | Huixquilucan | Wal-mart Express | 16,037 | 558 | 1 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Estado de México | Metepec | Wal-mart | 26,360 | 652 | 2 | 9 (2025-12 a 2026-08) | 2025-12-08 | 2026-08-31 |
| Estado de México | Naucalpan de Juárez | Bodega Aurrera | 11,878 | 530 | 2 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Estado de México | Naucalpan de Juárez | Wal-mart | 41,648 | 635 | 3 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Estado de México | Naucalpan de Juárez | Wal-mart Express | 31,608 | 566 | 2 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-31 |
| Estado de México | Nezahualcóyotl | Bodega Aurrera | 30,163 | 577 | 4 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Estado de México | Tlalnepantla | Bodega Aurrera | 15,909 | 518 | 3 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-26 |
| Estado de México | Tlalnepantla | Wal-mart | 57,551 | 661 | 3 | 9 (2025-12 a 2026-08) | 2025-12-01 | 2026-08-27 |
| Estado de México | Tlalnepantla | Wal-mart Express | 21,609 | 528 | 1 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-27 |
| Estado de México | Tultitlán | Bodega Aurrera | 27,545 | 540 | 1 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-26 |
| Estado de México | Tultitlán | Bodega Aurrera Express | 4,515 | 117 | 1 | 7 (2026-01 a 2026-08) | 2026-01-07 | 2026-08-19 |
| Guanajuato | León | Bodega Aurrera | 230,010 | 682 | 4 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Guanajuato | León | Wal-mart | 253,261 | 724 | 2 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Guanajuato | León | Wal-mart Express | 24,617 | 566 | 1 | 8 (2024-01 a 2024-08) | 2024-01-02 | 2024-08-15 |
| Guerrero | Acapulco de Juárez | Bodega Aurrera | 76,126 | 691 | 2 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Guerrero | Acapulco de Juárez | Wal-mart | 74,538 | 702 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Hidalgo | Pachuca de Soto | Bodega Aurrera | 20,133 | 599 | 2 | 22 (2024-01 a 2026-04) | 2024-01-08 | 2026-04-06 |
| Hidalgo | Pachuca de Soto | Wal-mart | 24,061 | 632 | 1 | 25 (2024-02 a 2026-08) | 2024-02-20 | 2026-08-18 |
| Jalisco | Guadalajara | Bodega Aurrera | 115,227 | 593 | 1 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Jalisco | Guadalajara | Bodega Aurrera Express | 55,827 | 295 | 3 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Jalisco | Guadalajara | Wal-mart | 134,934 | 653 | 2 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-31 |
| Jalisco | Guadalajara | Wal-mart Express | 148,490 | 571 | 3 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Jalisco | Zapopan | Wal-mart | 142,504 | 659 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Michoacan | Morelia | Bodega Aurrera | 3,362 | 175 | 1 | 16 (2024-08 a 2025-11) | 2024-08-22 | 2025-11-27 |
| Michoacan | Morelia | Wal-mart | 71,598 | 659 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Michoacan | Morelia | Wal-mart Express | 17,466 | 528 | 1 | 8 (2024-01 a 2024-08) | 2024-01-04 | 2024-08-15 |
| Michoacán | Morelia | Bodega Aurrera | 1,406 | 24 | 1 | 9 (2025-12 a 2026-08) | 2025-12-04 | 2026-08-28 |
| Michoacán | Morelia | Wal-mart | 29,631 | 637 | 1 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-25 |
| Morelos | Cuernavaca | Bodega Aurrera | 89,272 | 641 | 2 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-31 |
| Morelos | Cuernavaca | Wal-mart | 92,612 | 685 | 1 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Nuevo Leon | Monterrey | Bodega Aurrera | 34,966 | 568 | 1 | 14 (2024-01 a 2025-03) | 2024-01-05 | 2025-03-28 |
| Nuevo Leon | Monterrey | Wal-mart | 157,514 | 686 | 3 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-28 |
| Nuevo Leon | San Nicolás de los Garza | Bodega Aurrera | 67,471 | 620 | 1 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-26 |
| Nuevo León | Monterrey | Wal-mart | 32,765 | 545 | 1 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Nuevo León | San Nicolás de los Garza | Bodega Aurrera | 26,426 | 543 | 2 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-26 |
| Oaxaca | Oaxaca de Juárez | Wal-mart | 110,847 | 696 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Oaxaca | Santa Cruz Xoxocotlán | Bodega Aurrera | 92,669 | 653 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Puebla | Puebla | Bodega Aurrera | 99,660 | 619 | 1 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Puebla | Puebla | Wal-mart | 125,354 | 685 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Queretaro | Santiago de Querétaro | Bodega Aurrera | 97,604 | 661 | 2 | 23 (2024-01 a 2025-11) | 2024-01-05 | 2025-11-28 |
| Queretaro | Santiago de Querétaro | Wal-mart | 69,500 | 697 | 1 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Queretaro | Santiago de Querétaro | Wal-mart Express | 26,757 | 573 | 1 | 12 (2024-01 a 2024-12) | 2024-01-03 | 2024-12-18 |
| Querétaro | Santiago de Querétaro | Bodega Aurrera | 23,785 | 599 | 2 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Querétaro | Santiago de Querétaro | Wal-mart | 19,837 | 601 | 2 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-26 |
| Quintana Roo | Benito Juárez | Wal-mart | 94,293 | 603 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Quintana Roo | Solidaridad | Wal-mart | 70,371 | 711 | 2 | 29 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| San Luis Potosi | San Luis Potosí | Wal-mart | 74,103 | 672 | 1 | 23 (2024-01 a 2025-11) | 2024-01-04 | 2025-11-28 |
| San Luis Potosí | San Luis Potosí | Wal-mart | 25,200 | 534 | 1 | 9 (2025-12 a 2026-08) | 2025-12-05 | 2026-08-28 |
| Sinaloa | Culiacán | Bodega Aurrera | 28,564 | 463 | 2 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Sinaloa | Culiacán | Wal-mart | 33,484 | 252 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Sonora | Hermosillo | Bodega Aurrera | 72,496 | 648 | 3 | 30 (2024-02 a 2026-08) | 2024-02-14 | 2026-08-26 |
| Sonora | Hermosillo | Wal-mart | 91,237 | 678 | 3 | 31 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-27 |
| Tabasco | Centro | Bodega Aurrera | 222,777 | 661 | 4 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-26 |
| Tabasco | Centro | Wal-mart | 145,995 | 686 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Tabasco | Centro | Wal-mart Express | 81,282 | 595 | 1 | 30 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-26 |
| Tamaulipas | Reynosa | Bodega Aurrera | 32,394 | 548 | 2 | 7 (2026-02 a 2026-08) | 2026-02-20 | 2026-08-28 |
| Tamaulipas | Tampico | Bodega Aurrera | 56,928 | 613 | 1 | 25 (2024-08 a 2026-08) | 2024-08-19 | 2026-08-31 |
| Tamaulipas | Tampico | Wal-mart | 93,164 | 670 | 1 | 32 (2024-01 a 2026-08) | 2024-01-02 | 2026-08-25 |
| Tlaxcala | Tlaxcala | Bodega Aurrera | 96,681 | 645 | 1 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-27 |
| Tlaxcala | Tlaxcala | Wal-mart | 189,451 | 696 | 2 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-31 |
| Veracruz | Boca del Río | Wal-mart | 119,079 | 680 | 1 | 32 (2024-01 a 2026-08) | 2024-01-05 | 2026-08-28 |
| Veracruz | Orizaba | Bodega Aurrera | 4,707 | 25 | 1 | 25 (2024-01 a 2026-01) | 2024-01-04 | 2026-01-15 |
| Veracruz | Orizaba | Wal-mart | 4,610 | 180 | 2 | 25 (2024-01 a 2026-01) | 2024-01-02 | 2026-01-12 |
| Veracruz | Veracruz | Bodega Aurrera | 89,594 | 625 | 2 | 32 (2024-01 a 2026-08) | 2024-01-08 | 2026-08-31 |
| Yucatan | Mérida | Bodega Aurrera | 69,252 | 627 | 1 | 23 (2024-01 a 2025-11) | 2024-01-03 | 2025-11-26 |
| Yucatan | Mérida | Wal-mart | 161,841 | 709 | 2 | 23 (2024-01 a 2025-11) | 2024-01-02 | 2025-11-25 |
| Yucatan | Mérida | Wal-mart Express | 59,418 | 610 | 1 | 22 (2024-01 a 2025-10) | 2024-01-05 | 2025-10-24 |
| Yucatán | Mérida | Bodega Aurrera | 24,366 | 518 | 1 | 9 (2025-12 a 2026-08) | 2025-12-03 | 2026-08-26 |
| Yucatán | Mérida | Wal-mart | 51,770 | 655 | 3 | 9 (2025-12 a 2026-08) | 2025-12-02 | 2026-08-25 |
| Yucatán | Mérida | Wal-mart Express | 11,138 | 552 | 1 | 5 (2026-01 a 2026-05) | 2026-01-09 | 2026-05-22 |
| Zacatecas | Zacatecas | Bodega Aurrera | 163,002 | 665 | 5 | 32 (2024-01 a 2026-08) | 2024-01-04 | 2026-08-28 |
| Zacatecas | Zacatecas | Wal-mart | 209,192 | 700 | 2 | 32 (2024-01 a 2026-08) | 2024-01-03 | 2026-08-31 |

### Sucursales de Axel

**Regla de empate (escrita antes de contar):** un registro es de la sucursal si la `cadena_comercial` (con `Fold`: minúsculas, sin acentos, sin guion, punto ni apóstrofo) es una de las listadas y la frase aparece como palabras completas en `nombre_comercial` o en `direccion`. Para Sam's Club no se exige cadena: basta que «sams club» aparezca en `cadena_comercial`, `nombre_comercial` o `direccion`. Las «coincidencias parciales» (una sola palabra) se listan aparte solo como pista y no cuentan como registros de la sucursal.

**Bodega Aurrera Vallejo** — cadenas: bodega aurrera · bodega aurrera express · frase: «vallejo» · parcial: «vallejo»

- Resultado: **SÍ hay registros** de Bodega Aurrera Vallejo: 23,809 registros, 139 productos, 1 sucursales, meses 30 (2024-01 a 2026-08), del 2024-01-05 al 2026-08-21.
- De esos, la frase aparece **en `cadena_comercial` o `nombre_comercial`**: 0 registros; **solo en `direccion`** (referencia de ubicación, no el nombre del local): 23,809 registros.
  - Sucursales que cumplen la regla (hasta 15 de 1, nombre comercial | dirección | municipio):
    - Bodega Aurrera Express Sucursal Poniente 140 / Poniente 140 No. 237. Entre Norte 27 y Norte 25 a. Col. Lindavista Vallejo Iii Sección. Cp. 07720 / Gustavo a. Madero — 23,809 registros
- Coincidencias parciales («vallejo»): ninguna.

**Walmart Acueducto de Guadalupe** — cadenas: walmart · walmart express · frase: «acueducto de guadalupe» · parcial: «acueducto»

- Resultado: **NO hay registros** de Walmart Acueducto de Guadalupe con esta regla.
- Coincidencias parciales («acueducto»), no cuentan: 77,782 registros en 2 sucursales.
  - Parciales (hasta 8 de 2, nombre comercial | dirección | municipio):
    - Walmart Express Sucursal Virreyes / Aguilar y Seijas 123. Esq. Acueducto Rio Hondo. Col. Molino del Rey. Cp. 11040 / Miguel Hidalgo — 76,395 registros
    - Walmart Express Sucursal Virreyes / Aguilar y Seijas 123. Esq. Acueducto Río Hondo. Col. Molino del Rey. Cp. 11040 / Miguel Hidalgo — 1,387 registros

**Sam's Club** — cadenas: cualquiera · frase: «sams club» · parcial: «sams»

- Resultado: **SÍ hay registros** de Sam's Club: 94,293 registros, 603 productos, 1 sucursales, meses 32 (2024-01 a 2026-08), del 2024-01-02 al 2026-08-25.
- De esos, la frase aparece **en `cadena_comercial` o `nombre_comercial`**: 0 registros; **solo en `direccion`** (referencia de ubicación, no el nombre del local): 94,293 registros.
  - Sucursales que cumplen la regla (hasta 15 de 1, nombre comercial | dirección | municipio):
    - Wal Mart / Coba Sm. 21 Mz. 2 Lt. 2. Costado Sams Club. Col. Centro. Cp. 77500 / Benito Juárez — 94,293 registros
- Coincidencias parciales («sams»): ninguna.

### Todas las sucursales de las cadenas de interés en los municipios vigilados

Municipios: gustavo a madero. Sirve para que Axel reconozca sus tiendas por nombre y dirección; el cargador no adivina cuál es.

  - Gustavo a. Madero · Bodega Aurrera (hasta 20 de 1, nombre comercial | dirección | municipio):
    - Bodega Aurrera Sucursal Ferrocarril Hidalgo / Venturina 354. Entre Tesoro y Victoria. Col. Estrella. Cp. 07810 — 90,619 registros
  - Gustavo a. Madero · Bodega Aurrera Express (hasta 20 de 1, nombre comercial | dirección | municipio):
    - Bodega Aurrera Express Sucursal Poniente 140 / Poniente 140 No. 237. Entre Norte 27 y Norte 25 a. Col. Lindavista Vallejo Iii Sección. Cp. 07720 — 23,809 registros
  - Gustavo a. Madero · Wal-mart (hasta 20 de 1, nombre comercial | dirección | municipio):
    - Walmart Sucursal Tepeyac / Calz. de Guadalupe 431. Entre Henry Ford y Victoria. Col. Guadalupe Tepeyac. Cp. 07840 — 108,574 registros
  - Gustavo a. Madero · Wal-mart Express (hasta 20 de 1, nombre comercial | dirección | municipio):
    - Walmart Express Sucursal Lindavista / Montevideo 360. Esq. Instituto Politecnico Nacional. Col. Lindavista Sur. Cp. 07300 — 81,333 registros

