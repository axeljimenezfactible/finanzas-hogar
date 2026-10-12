# revisor-S5-r1 · Reporte final (WF-1, 11-oct-2026)

Revisor independiente del motor S5 (precios y tiendas), ronda 1 de 3. Informe: `docs/estado/wf1/REVISION-S5-r1.md`. **Veredicto: NO APROBADO** (4 altas, 4 medias, 4 bajas).

## Hecho
- Lecturas: protocolo completo; decisiones §E y §K; criterios S5 1.2 completos; contrato S5 1.4 (encabezado, §5 a §13, fila `profeco_import` de §15.2); reportes `motor-S5-go`, `oraculo-S5`, `profeco` y `ocr`; partes de `S5-cobertura-profeco.md`, la fixture `s5-profeco-muestra` con sus `_datos/`, y partes de `filter.go`, `import.go`, `cmd/motores/s5.go` y `precios_test.go`.
- Reejecuté cada verificación con mis comandos: vet, gofmt, pruebas `-run` por criterio, cobertura por paquete y sin `exec.go`, greps de red, de dominios y de `float`, el dump del motor con `compara.py s5`, `pre_run_fixtures.py`, `pre_verify.py`, `ocrmedir`, `tesseract`, conteos PROFECO (`wc -l` y SQL propio), reproducción del reporte de cobertura con un comando, memoria del cargador, independencia y sha256 de `data.json`.
- Escribí mi propio cálculo en `datos/tmp/rev_s5.py` (ya borrado): canon desde `data.json` con `Decimal`, unidad base, pesos K40 con `Fraction`, equilibrio de Sam's y comparación por multiconjunto de las filas motor contra oráculo.

## Archivos
- `docs/estado/wf1/REVISION-S5-r1.md` (informe) y `docs/estado/wf1/revisor-S5-r1.md` (este reporte). No modifiqué nada más.

## Verificación
- `go vet ./internal/precios/... ./cmd/motores/` limpio; `gofmt -l internal/precios cmd/motores` vacío.
- `go test ./internal/precios/... -count=1 -cover`: precios 96.0 %, asistida 100 %, ocr 99.3 %, profeco 95.3 %; total del perfil 88.9 %. Ocr sin `exec.go`: 99.3 % (295/297), con y sin la integración.
- `-run` Comparativa, Unidad, MejorTienda, Frescura, Fuente, LecturaAsistida, Membresia, Mercado, PreciosViejos: todo PASS.
- Mi cálculo: 309450 / 323300 / +13850 = +10050 + 5900 − 2100; 277350 / 287400; 61 / 60 / 55; 33 / 3; 22 / 8 / 3; 2933/2800, 12692/16406, 3375/2833; pesos 1, 0.9836, 0.9672, 0.4918, 0; equilibrio 8 y 9.
- `motores s5 dump` + `python3 -I oraculo/compara.py s5`: **código 1**; oráculo contra esperado 747/747; motor contra esperado 713/747; 65 diferencias (canon-filas 34, mejor-tienda-incompleta 7, profeco-muestra 24).
- PROFECO: `wc -l` 9,388,473 − 16 = 9,388,457 = `COUNT(*)`; `MAX(fecha_registro)` 2026/08/31; reproducción con un comando con diff de 0 líneas contra la sección 6; muestra 615,625 bytes; pico de 29 MB con un CSV de 197 MB. Faltan 2024 y 2025.
- OCR: tesseract 5.5.3 con spa; `ocrmedir` dice «sin medir: faltan las fotos de Axel»; `TestOCRSinConfirmar` PASS.
- `docs/estado/S5-terminos.md`: no existe.
- `data.json`: sha256 `2fe7afae5b1be77fafeed787d4c7bd4d68fe469b78a86898b370c1ba1031d31b`; `git status` limpio.

## Supuestos
- Leí «cargador en flujo» como lectura por flujo (sin cargar el archivo entero a memoria) y lo medí con el pico de memoria.
- H-04 (`compara.py s5` falla) junta las diferencias de H-01, H-05, H-06 y H-07. Lo cuento aparte porque la tarea pide esa verificación de forma explícita.

## Bloqueos
- Ninguno escrito por mí: el revisor no escribe bloqueos. En H-02 y H-03 propongo abrir `docs/estado/bloqueos/S5-profeco-cortes.md` (opcional) y `docs/estado/bloqueos/S5-terminos.md`.

## Pendientes
- Para la ronda 2: corregir H-01 a H-08 y volver a correr `compara.py s5` hasta «diferencias: 0».

## Solicitudes fuera de alcance
- `docs/contratos/motores/S5-interfaz.md`: §11 (municipio de una cadena en el filtro; `file` y `row` de los rechazos; rutas reales; sin Sam's en PROFECO), §7.2 y §15.2 (orden de `rows[]`), §7.4 (dónde va una tienda incompleta), §12 (comando real de Tesseract; carpeta de tickets), §15.3 (nombre del ejecutor).
- `docs/criterios/S5.md` S5-05 contra plan §3.4 («2026» contra «2024, 2025 y 2026»): lo decide el integrador o la puerta 2.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Para S5-terminos (H-03): decidir si autoriza leer solo las páginas legales de Walmart, Bodega Aurrera y Sam's Club, o si él mismo dicta el veredicto.

## Consumo aproximado
- 33 acciones de herramienta, unos 165 KB leídos y unos 21 KB escritos.
