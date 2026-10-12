# revisor-S5-r3 · WF-1 · revisión independiente del motor S5 (ronda 3 de 3)

Agente `revisor-S5-r3` · 11-oct-2026. Informe: `docs/estado/wf1/REVISION-S5-r3.md`. Veredicto: **NO APROBADO** (1 alta, 4 medias, 2 bajas).

## Hecho
- Volví a correr todas las verificaciones de S5-01 a S5-13 con mis comandos y escribí cálculos propios: canon desde `data.json` con Decimal y frescura contra las 64 edades del motor.
- Volví a contar PROFECO (`wc -l` de los 64 CSV y `COUNT(*)` en `qqp.sqlite`).
- Armé 7 fixtures propias en una raíz sombra para comparar el motor con el oráculo en filtro, lectura y liga PROFECO, Windows-1252, carritos con 3 tiendas y cantidades a granel.
- Comprobé los 4 hallazgos de la r2. R2-02 y R2-03 quedan cerrados como PROVISIONAL. R2-01 (alta) y R2-04 (media) siguen abiertos.
- Hallazgos nuevos: H-r3-02 (granel), H-r3-03 (lectura y liga de PROFECO) y H-r3-04 (`difference_vs_cheapest_cents`), los tres medios, y H-r3-05 y H-r3-06, bajos.

## Archivos
- Creados: `docs/estado/wf1/REVISION-S5-r3.md` y `docs/estado/wf1/revisor-S5-r3.md`.
- Temporales en `datos/tmp/rev5r3/` (scripts `canon.py`, `fresc.py`, `gen/gen.py` y `gen/diff.py`, la raíz sombra y el binario `motores`): **borrados** al terminar. No toqué nada más.

## Verificación
- `go vet ./internal/precios/... ./cmd/motores/` limpio; `gofmt -l` vacío. `go test -count=1 -coverprofile ./internal/precios/...` da precios 96.1 %, asistida 100 %, ocr 99.3 %, profeco 95.3 % y total 88.8 %, con 78 PASS, 0 FAIL y 0 SKIP.
- `-run` de cada criterio: Comparativa 4, Unidad 2, MejorTienda 5, Frescura 3, Fuente 4, LecturaAsistida 1, Membresia 2, Mercado 2, PreciosViejos 1, OCRSinConfirmar 1, EvalCarritos 2. Todo PASS.
- `canon.py` sobre `data.json`: 309450 / 323300 / +13850; idénticos 277350 / 287400 / +10050; piezas 61/60/55; renglones 33/3; ganadores 22/8/3; sustitutos 29.33/28.00, 126.92/164.06 y 33.75/28.33. Resultado: CANON OK.
- `motores s5 dump` y `compara.py s5` dan 747/747 (oráculo contra esperado, motor contra esperado y oráculo contra motor), 301/301 provisionales y 0 diferencias. `pre_run_fixtures.py` da 29/29 y deja `salidas/s5.json` idéntico (`cmp`). `pre_verify.py` → TODO COINCIDE.
- `fresc.py`: 64 edades con 0 diferencias contra (90 − d)/61; site_assisted y profeco no compiten aunque sean más baratos.
- PROFECO: wc 14,889,243 / 13,270,548 / 9,388,473 y SQLite 14,889,219 / 13,270,524 / 9,388,457 (total 37,548,200; máximo 2026/08/31). Muestra de 601 KB. Los RAR existen y git los ignora.
- OCR: tesseract 5.5.3 con spa; `ocrmedir` → «S5-06 sin medir … faltan las fotos de Axel».
- Sam's: 8, 9 (8.33 hacia arriba), 6 y «sin costo capturado (PROVISIONAL)». Mercado: martes 13 y martes 20; uva sí, avena no.
- S5-07: `docs/estado/S5-terminos.md` no existe. S5-08: grep de red da 0. Grep de dominios de tiendas da 0.
- Independencia: `data.json` y `model.json` sin cambios en git; no hay referencias cruzadas entre el motor y el oráculo.
- Fixtures propias, contadas como diferencias motor contra oráculo: filtro en revisión 0; Windows-1252 0; dos de tres 0; granel 8; granel con sustituto 2; bordes de lectura 13 en producción y 13 en revisión.

## Supuestos
- Clasifiqué H-r3-02, H-r3-03 y H-r3-04 como medios, igual que la r2 trató divergencias del mismo tipo: el contrato no fija la regla y no cambia ninguna cifra canon.
- S5-05 queda como «cumple en sustancia»: los conteos reales se reproducen y los hallazgos medios están en el contrato y en el oráculo.

## Bloqueos
- Ninguno nuevo mío. Siguen abiertos `docs/estado/bloqueos/S5-terminos.md`, `S5-profeco-filtro.md` y `S5-carritos-detalle.md`. Sugiero agregar la pregunta 4 (`difference_vs_cheapest`) a `S5-carritos-detalle` y abrir `S5-carritos-granel` y la regla de lectura de PROFECO si no son técnicas.

## Pendientes
- H-r3-00: decisión de Axel en la puerta 2 sobre `S5-terminos`.
- H-r3-01: el integrador aplica las solicitudes 1 a 6 de `correccion-S5-r2.md` (contrato 1.5) y suma la regla de 3 o más canastas.
- H-r3-02 a H-r3-04: regla en el contrato, alinear el motor o el oráculo y agregar fixtures (detalle en el informe).

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S5-interfaz.md` §7.2, §7.4, §11, §12 y §15.2: lo listado en H-r3-01 a H-r3-05.
2. `oraculo/pre_calculo.py` `_agrupa`: que el precio por unidad base salga de la fracción exacta (H-r3-02).
3. `oraculo/pre_profeco.py` y `backend/internal/precios/profeco/{parse,links}.go`: alinear las reglas de lectura y de liga una vez escritas (H-r3-03).
4. `docs/estado/S5-cobertura-profeco.md` §1 y §4.2: las cifras de los tres cortes (H-r3-06).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna de K31.

## Consumo aproximado
~42 acciones de herramienta; ~150 KB leídos; ~45 KB escritos (los dos informes, los scripts y las fixtures temporales generadas).
