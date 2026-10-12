# WF-1 · agente ocr (Sonnet) · 11-oct-2026

## Hecho
- Paquete `backend/internal/precios/ocr` (Go, solo biblioteca estándar; `tesseract` y `magick` por `os/exec` sin shell):
  - `Engine.Recognize(ctx, imagen)`: valida la imagen (.jpg/.jpeg/.png, ruta absoluta), verifica **una vez** `tesseract --version` = 5.5.3 y `spa` en `--list-langs` (si no, `ErrOCRUnavailable`), preprocesa y corre Tesseract con tiempo máximo.
  - Preprocesado (ImageMagick, opcional, con respaldo a la imagen original y aviso si falla): `-auto-orient -colorspace Gray -resize 1600x< -deskew 40% -normalize -lat 41x41-6% -type Bilevel` (escala de grises, enderezado, umbral local que resiste sombras de foto).
  - Configuración de Tesseract (documentada en `doc.go`): `tesseract <img> stdout -l spa --psm 6 --oem 1 -c preserve_interword_spaces=1 tsv`. Es el comando del contrato (§12) más `--oem 1` y `preserve_interword_spaces=1` (conserva el hueco entre descripción e importe).
  - `ParseTSV`: agrupa palabras de nivel 5 por página/bloque/párrafo/línea; confianza = promedio HALF_UP.
  - `ParseLine`: importe = último token con dos decimales (`$`, miles, signo `-`; corrige O→0 e I,l,|→1 solo dentro de un token con forma de importe; lo que sigue, como la letra de IVA, se descarta); cantidad = número inicial (entero o hasta 3 decimales, coma o punto) más unidad opcional (KG, G, L, ML, PZA...), solo si queda texto después; descripción = el resto; marca `EsResumen` en TOTAL, SUBTOTAL, IVA, CAMBIO, EFECTIVO, TARJETA, PAGO, ARTICULOS, AHORRO, REDONDEO. Renglones sin importe se devuelven con `importe` nulo para que la persona los descarte. Todo en enteros (centavos y milésimas), sin `float` en dinero.
  - Perfiles: `generico` (al pie de la letra de §12) y `con-codigo` (quita códigos de barras de 8 a 14 dígitos). **Ninguno está validado contra tickets reales** (ver Supuestos).
  - `PricesFromOCR(lineas, confirmado)`: con `confirmado = false` devuelve `ErrNotConfirmed` y ningún renglón (un OCR sin confirmar no crea precios). Confirmado, entrega solo los renglones con importe que no son resumen.
  - Medición S5-06: `Measure`, `Score`, `ReadExpectedCSV`, `NormalizeDescription` (regla congelada en `measure.go`) y la herramienta `cmd/ocrmedir`.
- `datos/semilla/tickets/README.md` pide a Axel las 5 fotos, la forma de fotografiar, el formato de la transcripción y el comando de medición.
- 2 tickets **sintéticos** (`datos/semilla/tickets/sinteticos/sintetico-01.png`, `-02.png`, con transcripción en `sinteticos/esperado/`; generador `backend/internal/precios/ocr/testdata/gen-sinteticos.sh`). No son fotos reales ni imitan el formato de una tienda.

### Regla congelada de «renglón correcto» (S5-06; escrita antes de medir, no se afloja)
1. Esperado = fila de la transcripción a mano (`datos/semilla/tickets-esperado/NN.csv`: `descripcion,cantidad,importe`), con la descripción **impresa**, sin código de barras ni letra de IVA; cantidad solo si el ticket la imprime; importe impreso.
2. Correcto = existe un renglón leído, aún no usado, con descripción igual tras normalizar (mayúsculas, sin acentos, espacios colapsados, sin puntuación final) **y** cantidad exacta (ambas vacías o iguales en milésimas) **y** importe exacto en centavos.
3. Los de resumen y los sin importe no compiten; los leídos de más no restan, se reportan como «sobrantes».
4. Porcentaje = correctos / esperados sumando los 5 tickets; cumple con `correctos × 100 ≥ 80 × esperados`.
5. Sin las 5 imágenes con su transcripción: «sin medir», nunca cumplido.

## Archivos
- `backend/internal/precios/ocr/`: `doc.go`, `types.go`, `parse.go`, `tsv.go`, `engine.go`, `exec.go`, `measure.go`; pruebas `parse_test.go`, `engine_test.go`, `confirm_test.go`, `measure_test.go`, `integration_test.go`; `testdata/gen-sinteticos.sh`; `cmd/ocrmedir/main.go`.
- `datos/semilla/tickets/README.md`; `datos/semilla/tickets/sinteticos/{sintetico-01.png,sintetico-02.png,esperado/sintetico-01.csv,esperado/sintetico-02.csv}`.
- `docs/estado/wf1/ocr.md` (este reporte).

## Verificación
- `tesseract --version` → `tesseract 5.5.3` (leptonica 1.87.0); `tesseract --list-langs` → `osd`, `spa`.
- `cd backend && gofmt -l internal/precios/ocr` → vacío; `go vet ./internal/precios/...` → limpio.
- `go test ./internal/precios/ocr/ -cover -count=1` → ok, **99.3 %** (con la prueba de integración que sí corre el exec real); sin la integración (solo ejecutor simulado) **97.7 %** del paquete y **99.3 %** sin contar `exec.go` (criterio: ≥ 85 % sin contar el exec). Única función sin cubrir sin la integración: `ExecRunner.Run`.
- Criterio **S5-06**: pruebas `TestOCRParseo`, `TestOCRSinConfirmar`, `TestOCRVersion` (nombres del contrato §16) en verde con ejecutor simulado; `TestOCRVersion` cubre versión distinta, sin versión, sin `spa` y tesseract ausente; `TestOCRSinConfirmar`: un OCR sin confirmar no devuelve renglones.
- **Medición S5-06 con las 5 fotos reales: «sin medir: faltan las fotos de Axel»** (`go run ./internal/precios/ocr/cmd/ocrmedir` → `S5-06 sin medir: hay 0 imágenes y se necesitan 5: faltan las fotos de Axel`). No hay ningún porcentaje real que reportar.
- Solo como prueba de la herramienta, sobre los 2 sintéticos (NO cuenta para S5-06): perfil `generico` 8/11 = 72.7 %; perfil `con-codigo` 9/11 = 81.8 %. El preprocesado no cambió el resultado sobre estas imágenes limpias (mismo 8/11 con y sin él), así que su efecto sobre fotos reales (sombras, arrugas, brillo) queda sin evidencia hasta tener las fotos. Fallas: Tesseract leyó `900` como `908` y `200` como `208` (confusión de dígitos que no se corrige a ciegas: sería inventar) y el perfil `generico` conserva el código de barras en la descripción, por diseño.
- Sin `float` en dinero (`grep float` en el paquete: solo el promedio de la confianza 0–100 en `tsv.go`; importes y cantidades son enteros). Sin red.

## Supuestos
1. **Formatos de Bodega Aurrera y Walmart no inventados:** la tarea pedía expresiones regulares por tienda, pero no hay tickets reales ni descripción de su formato en el repositorio y el contrato prohíbe inventar. Se implementó el análisis de §12 literal (perfil `generico`) y un mecanismo de perfiles; el perfil `con-codigo` es una hipótesis (código de barras antes del nombre) sin validar. Se afinan con las 5 fotos (ver Pendientes).
2. Cantidad vacía cuando el ticket no la imprime (no se asume 1); quien confirma decide. El contrato §12 tampoco asume.
3. La unidad impresa (KG, G...) se guarda en `unidad` y no cambia el valor de `cantidad_milli` (milésimas de la unidad impresa).
4. `--oem 1` y `preserve_interword_spaces=1` se añaden al comando del contrato; el resto del comando es idéntico.
5. `ErrNotConfirmed` y `ErrOCRUnavailable` viven en el paquete `ocr` porque `internal/precios` aún no existía al escribirlo; el integrador puede reexportarlos en `precios` (alias `var ErrNotConfirmed = ocr.ErrNotConfirmed`) sin cambiar nada más.
6. El corrector O→0 / I,l,|→1 solo actúa dentro de un token con forma de importe (`dígitos.dd`) que ya trae al menos un dígito real; nunca toca descripciones.

## Bloqueos
- Ninguno nuevo. S1 (`docs/estado/bloqueos/S1-parametros-avisos.md`) no afecta a este agente. `criterios-ocr-tickets` está RESUELTO (K24); lo que falta son las fotos (acción de Axel), no una decisión.

## Pendientes
1. **Axel:** fotografiar los 5 tickets en `datos/semilla/tickets/01..05` y transcribirlos (o pedir que alguien lo haga con él) en `datos/semilla/tickets-esperado/NN.csv`, siguiendo el README. Después correr `go run ./internal/precios/ocr/cmd/ocrmedir` y copiar la tabla al reporte `S5-06`.
2. Con las fotos reales: revisar los renglones fallidos y, si hace falta, afinar `Profile` (regex por tienda, p. ej. línea con descuento o con peso en dos renglones). Los cambios al perfil deben hacerse **antes** de congelar la medición definitiva, no después de ver el porcentaje; si se ajusta algo, documentar ambos porcentajes.
3. Conectar `ocr.Engine` al borrador `attachment.ocr_draft_json` (K3) y a `uploadExpenseAttachment`: lo hace la API; este paquete solo entrega `[]Line`.
4. Si el integrador lo quiere, exponer `ocr_parse_tsv` de la fixture `s5-ocr-parseo` sobre `ParseTSV` (claves en inglés: `description`, `quantity_milli`, `amount_cents`, `confidence`); no se tocó `backend/cmd/motores` por estar fuera de alcance.

## Solicitudes fuera de alcance
1. `docs/criterios/S5.md` (S5-06, línea «Verifica») y `docs/contratos/motores/S5-interfaz.md` §12: el conteo `ls datos/semilla/tickets/ | wc -l` = 5 deja de cumplirse en cuanto la carpeta tenga el `README.md` y `sinteticos/` que esta tarea manda crear (daría 7). Cambio propuesto: contar imágenes (`ls datos/semilla/tickets/*.jp*g datos/semilla/tickets/*.png | wc -l`), o mover `README.md` y `sinteticos/` a `datos/semilla/tickets-sinteticos/`. Motivo: contradicción entre el alcance asignado y el criterio. `ocrmedir` ya cuenta solo imágenes del primer nivel.
2. `datos/semilla/tickets-esperado/` (alcance de quien transcriba): crear ahí `NN.csv` antes de medir. No lo creé porque no estaba en mi alcance y las transcripciones reales son de Axel.
3. S1 / esquema (`app_setting`): clave `ocr.timeout_seconds` (valor inicial 30, K44) y sección de Ajustes (Datos y respaldos o Tiendas y precios, a decidir en la puerta 2). Hoy el motor usa la constante `DefaultTimeoutSeconds = 30`, PROVISIONAL.
4. `backend/internal/precios/types.go` (de otro agente) no está formateado: `gofmt -l internal/precios` lo lista. No es del paquete `ocr`.
5. Instrucción al integrador: los temporales del preprocesado se crean en el directorio temporal del sistema (`os.MkdirTemp`) y se borran al terminar; para que un servicio con `PrivateTmp` no falle, se puede fijar `Engine.TempDir`.

## Paquetes instalados
Ninguno (tesseract 5.5.3, ImageMagick y Go 1.27.2 ya estaban).

## Propuestas para Axel
- **K31/K44 · tiempo máximo del OCR por imagen (`ocr.timeout_seconds`): 30 s.** Justificación: en esta máquina la prueba de integración completa (2 imágenes sintéticas de 760×560, con preprocesado) tarda unos 2 s entre preprocesado y lectura; una foto real de teléfono (12 MP) en un servidor doméstico más lento puede tardar 5–10 s; 30 s deja margen de 3× sin dejar la petición colgada. Configurable en Ajustes; se aprueba en la puerta que toque.
- **Fotos de los tickets:** que sean de la compra real (K24) y se tomen como dice el README; si un ticket sale ilegible a simple vista, repetir la foto antes de transcribir (la medición evalúa el OCR, no la foto mal tomada), pero sin tirar las que simplemente salgan feas.

## Consumo aproximado
- Acciones de herramienta: ~23 de 60.
- Leído: ~65 KB (protocolo 5.6, criterios S5 19, decisiones §K ~15, extractos de S5-interfaz y consultas ~25).
- Escrito: ~75 KB (código y pruebas ~60, README y reporte ~15; las 2 imágenes sintéticas pesan ~300 KB cada una y no se cuentan).
