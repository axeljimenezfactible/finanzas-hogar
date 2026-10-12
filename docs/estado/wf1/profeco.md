# Reporte WF-1 · agente profeco (Sonnet) · 11-oct-2026 (reanudación)

## Hecho
- Integridad de la descarga previa: `bsdtar -tf` del RAR sin errores; los 16 CSV extraídos en `datos/semilla/profeco/raw/2026-csv/2026/` coinciden en tamaño con el listado del RAR (`bsdtar -tvf`). Contenido tratado como dato no confiable.
- Parciales del intento anterior revisados contra S5-interfaz 1.4 §11 (ronda 4): se conservaron `parse`, `text`, `reader`, `links`, `filter`, `import`, `sqlite`, `survey`, `report` (compilan y pasan pruebas); se ampliaron `survey.go` (separa «frase en el nombre» de «solo en la dirección» y lista las sucursales de los municipios vigilados) y `report.go`.
- Nuevo `cobertura/main.go`: reporte de cobertura, carga a `qqp.sqlite`, medición del filtro y generación de la muestra.
- Pruebas: `parse_test.go`, `import_test.go`, `sqlite_test.go` (esquema real 00001+00002+00003 en SQLite), `survey_test.go`. Incluyen `TestProfecoConteo`, `TestProfecoCodificacion`, `TestProfecoFiltroYLiga`, `TestProfecoLigaSinConfirmar` (nombres de S5-interfaz §15 tabla S5-05).
- `datos/semilla/profeco/muestra/qqp-muestra.csv`: 1,974 filas reales, 221 cadenas, 615,625 bytes (≤ 2 MB), de los 8 meses.
- `qqp.sqlite` (base aparte, 4.0 GB, ignorada por git) con los 9,388,457 registros.
- `docs/estado/S5-cobertura-profeco.md` con resumen, regla de empate escrita antes de contar, y conteos generados por el cargador.

## Archivos
- backend/internal/precios/profeco/{parse,text,reader,links,filter,import,sqlite,survey,report}.go (survey y report modificados)
- backend/internal/precios/profeco/{parse,import,sqlite,survey}_test.go
- backend/internal/precios/profeco/cobertura/main.go
- datos/semilla/profeco/muestra/qqp-muestra.csv
- datos/semilla/profeco/db/.gitignore (y qqp.sqlite, ignorada)
- datos/semilla/profeco/raw/ (RAR, CSV extraídos, .gitignore; ignorados por git)
- docs/estado/S5-cobertura-profeco.md · docs/estado/wf1/profeco.md

## Verificación
- `go test ./internal/precios/profeco/ -cover -count=1` → ok, **95.3 %** (umbral 85 %). `go vet ./internal/precios/...` limpio; `gofmt -l internal/precios` vacío; `go test ./internal/precios/... ` ok (ocr 99.3 %).
- Conteo S5-05: 16 archivos, `wc -l` total 9,388,473 = 9,388,457 registros + 16 encabezados; rechazados 0; diferencia 0 por archivo. `sqlite3 qqp.sqlite "SELECT count(*),min(fecha),max(fecha) FROM profeco_record"` → `9388457|2026-01-02|2026-08-31` (la fuente declara cobertura hasta 2026-08-31).
- Cruce independiente por SQL: Wal-mart 1,500,220 · Wal-mart Express 412,731 · Bodega Aurrera 1,020,792 · Bodega Aurrera Express 33,491, igual que el reporte del cargador. Filtro medido (4 cadenas × Gustavo a. Madero) = 74,154 = 22,485 + 4,461 + 26,633 + 20,575.
- Tiempos: lectura y conteo 57 s; carga a qqp.sqlite 1 min 49 s; lectura + filtro 29 s.
- Criterios: S5-05 parcial (2026 cargado y cobertura reportada; faltan 2024 y 2025, ver Pendientes); S5-08 (`grep net/http` en profeco: sin resultados, el cargador no usa red); S5-10/K22 (todo precio entra con `source='profeco'`, `is_promotional=0`, prueba `TestProfecoCargaEnBase`); S5-13 cobertura ≥ 85 %; sin `float` en dinero (`ParsePriceCents` por texto).

## Hallazgos de los datos
1. Cadenas reales: `Wal-mart`, `Wal-mart Express`, `Bodega Aurrera`, `Bodega Aurrera Express`. **No existe Sam's Club** como cadena ni como nombre comercial; los 24,788 registros que la regla de texto encuentra son una tienda Wal-mart de Cancún cuya dirección dice «Costado Sams Club».
2. Bodega Aurrera Vallejo: una candidata (Express Poniente 140, Gustavo A. Madero; «Vallejo» solo en la colonia). Walmart Acueducto de Guadalupe: cero con la regla; en Gustavo A. Madero PROFECO mide Walmart Tepeyac y Walmart Express Lindavista.
3. Fecha: formato único `AAAA/MM/DD`. Sin renglones rechazados. Precio máximo 107,691.00 (electrodoméstico plausible; sin tope impuesto).

## Supuestos
- Extracción en `raw/2026-csv/` como pidió el encargo (el contrato §11 dice `datos/tmp/`).
- `qqp.sqlite` en `datos/semilla/profeco/db/` (no `datos/profeco/`) por la regla «nada de datos fuera de datos/semilla/profeco/»; sigue PROVISIONAL (bloqueo `S5-profeco-almacen`).
- En modo revisión las ligas propuestas se cuentan en `LoadedProposed`, separado de `Loaded` (contrato: «se cuentan aparte»); el oráculo `profeco_import` solo expone `loaded` y debe decidir si suma ambos.
- «Primero gana» para duplicados es el criterio PROVISIONAL del contrato (bloqueo `S5-profeco-almacen`, pregunta 3).
- Los dos nombres «Express» son cadenas distintas en PROFECO; se tratan como cadenas aparte (no se asume que `Wal-mart` los incluya).

## Bloqueos
- Ninguno nuevo. Siguen abiertos por contrato `docs/estado/bloqueos/S5-profeco-almacen.md` y `S4-confirmacion-empates.md` (no tocados).

## Pendientes
1. Descargar los cortes **2024 y 2025** (CSV) en carpetas propias `datos/semilla/profeco/raw/2024/` y `.../2025/` (sin scripts dentro) y cargarlos con el mismo comando (`-dir` por corte); actualizar `docs/estado/S5-cobertura-profeco.md` (criterio S5-05 completo). Sin cambios de código previstos.
2. `datos/semilla/profeco/vinculos.csv` (ligas producto→presentación) no existe; lo propone el flujo de catálogo y lo confirma Axel (puerta 2). Hasta entonces la importación a `price` carga 0 registros.
3. Axel confirma `profeco_name` de las cadenas (4 nombres) y reconoce sus sucursales en la lista de Gustavo A. Madero.
4. Fixture del oráculo `oraculo/fixtures/s5/_datos/profeco-muestra.csv` (6 renglones, S5-interfaz §15) lo escribe el agente del oráculo S5; el cargador ya cubre esos casos con `TestProfecoCodificacion`/`Conteo`.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S5-interfaz.md` §11: anotar que PROFECO escribe `Wal-mart` (con guion) y que Express son cadenas aparte; que no existe `cadena_comercial` Sam's Club (Sam's sin precios de referencia); y la ruta real de extracción/base (`raw/2026-csv/`, `datos/semilla/profeco/db/qqp.sqlite`). Motivo: el contrato supone una cadena Sam's Club.
- `.gitignore` raíz: agregar `datos/semilla/profeco/raw/` y `datos/semilla/profeco/db/` es opcional (ya tienen `.gitignore` propio que ignora todo).

## Paquetes instalados
- Ninguno (bsdtar, unrar y sqlite3 ya estaban).

## Propuestas para Axel
- Cadenas en el catálogo: `Wal-mart` (profeco_name de Walmart), `Bodega Aurrera`, y las dos variantes `Express` como cadenas aparte si quiere sus precios; Sam's Club sin profeco_name (no hay datos).
- Selección de municipios en Ajustes por defecto: el de cada sucursal registrada (Gustavo A. Madero para las tres de Axel si se confirman).

## Consumo aproximado
- Acciones de herramienta: ~34. Leído: ~120 KB (protocolo, contrato §11, criterios, parciales, salidas). Escrito por herramientas: ~100 KB (código, pruebas, documentos); la muestra de 615 KB y la base la generó el programa.
