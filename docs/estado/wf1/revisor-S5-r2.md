# revisor-S5-r2 · WF-1 · reporte de agente

Revisor independiente de S5, ronda 2 de 3 · 11-oct-2026. El informe está en `docs/estado/wf1/REVISION-S5-r2.md`. **Veredicto: NO APROBADO.** Hay 1 alto (R2-01, S5-07, bloqueado por Axel), 3 medios (R2-02, R2-03, R2-04) y 3 bajos (R2-05, R2-06, R2-07).

## Hecho
- Leí, por partes, el protocolo, las decisiones §E y §K, los criterios S5 1.2 completos, `correccion-S5-r1.md`, el contrato S5 (§7.4, §11, §12 y extractos de §7.2/§15.2), RN-12/RN-13, el bloqueo `S5-terminos` y las secciones 1–5 y la sección generada de `S5-cobertura-profeco.md`.
- Volví a correr todas las verificaciones (a)–(g) con mis propios comandos y scripts en `datos/tmp/rev-s5r2/` (ya borrados). Los resultados están en el informe.
- Revisé los 8 hallazgos de la ronda 1: 7 cerrados en la conducta (3 de ellos con el texto del contrato pendiente, R2-04) y 1 abierto (S5-07). El filtro de municipio sigue latente como R2-02.
- Comprobé 4 divergencias latentes entre el motor y el oráculo con 4 fixtures propias: cadena con municipio, ortografía «Wal-mart», sustituto con dos presentaciones y `missing_count`.

## Archivos
- Creados: `docs/estado/wf1/REVISION-S5-r2.md` y `docs/estado/wf1/revisor-S5-r2.md`.
- Temporales: `datos/tmp/rev-s5r2/`, borrado al terminar. No modifiqué nada más.

## Verificación
- `go vet ./internal/precios/... ./cmd/motores/` limpio; `gofmt -l internal/precios cmd/motores` vacío.
- `go test -count=1 -coverprofile=… ./internal/precios/...`: precios 96.1 %, asistida 100 %, ocr 99.3 % (`exec.go` 100 %), profeco 95.3 %; total con los `main` 88.8 %.
- `-run` por criterio: Comparativa 4, Unidad 2, MejorTienda 5, Frescura 3, Fuente 4, Membresia 2, Mercado 2, PreciosViejos 1, LecturaAsistida 1. Todo PASS, 0 FAIL.
- `motores s5 dump` (desde la raíz) + `pre_run_fixtures.py --out datos/tmp/…` + `compara.py s5`: 747/747 en los tres cruces, provisionales 225/225, diferencias 0, código 0. Desde `backend/`: 419/747 (R2-05).
- Mi canon desde `data.json`: +13850, +10050, 61/60/55, 33/3; por unidad base 2933/2800, 12692/16406, 3375/2833.
- PROFECO: mi `wc -l` − encabezados da 14,889,219 / 13,270,524 / 9,388,457. `sqlite3 -readonly` da 37548200 y MAX 2026/08/31, con `MIN(renglon)` 1. La regeneración de la cobertura (4 min 2 s) es idéntica a la §6. La muestra (601 KB), regenerada sobre 2026, también es idéntica.
- OCR: tesseract 5.5.3 con `spa`; `ocrmedir` → «sin medir». Frescura: `[0.9836, 0.4918, 0.0]`.
- `data.json`: su sha256 es igual al de HEAD. Grep de red, de dominios de tiendas y de `float` en dinero: 0 líneas.

## Supuestos
- Mantengo S5-07 como alto aunque esté bloqueado correctamente: el criterio es alto y no tiene salida «sin medir».
- S5-03 cumple con los desempates PROVISIONALES del bloqueo `S5-reglas-comparativa`, igual que en la ronda 1.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S5-terminos.md`, que depende de Axel.

## Pendientes
- R2-01: decisión de Axel sobre `S5-terminos`.
- R2-02 y R2-03: el contrato debe fijar §11 y §7.2; luego hay que alinear el motor o el oráculo y añadir fixtures.
- R2-04: aplicar en el contrato S5 las solicitudes de `correccion-S5-r1`.
- S5-06: las 5 fotos de Axel.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S5-interfaz.md` §7.2, §7.4, §11, §12 y §15.2: lo descrito en R2-02, R2-03 y R2-04.
- `docs/criterios/S5.md` S5-05 y plan §3.4: los tres cortes vienen en RAR; «PROFECO 2024–2026».

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna de K31. Decidir el bloqueo `S5-terminos` (opción A, B, C o D).

## Consumo aproximado
~52 acciones de herramienta; ~115 KB leídos; ~40 KB escritos (incluidas las fixtures temporales).
