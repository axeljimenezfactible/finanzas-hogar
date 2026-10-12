# Integración de WF-1 · Motores y sistema de diseño

Integrador: `integrador-wf1` (Opus) · 11-oct-2026 · repositorio `/home/deepsec/finanzas-hogar` (git local sin remoto). Corrida: continuación del run `wf_98b77133-bf1` (pausado a las 19:00:06 por un reinicio de la sesión del orquestador; ver §3). El hash del commit de WF-1 no puede ir dentro del propio commit: lo reporta el integrador en su salida.

## 0. Resumen

- **La continuación quedó completa.** Los 11 agentes interrumpidos se relanzaron desde cero y los 11 dejaron su reporte; las rondas extra (4 y 5) corrieron antes de los motores que dependían de esos artefactos; el paquete para Claude Design se regeneró y el integrador lo volvió a generar con `U4-ajustes.md` unido.
- **Todo lo verificable pasa:** `make verify-all` (WF-0, C1, S1, U1, criterios y el nuevo `verify-wf1`) termina en «Verificación ampliada OK». Los 4 motores reproducen sus fixtures y `oraculo/model.json` con **0 diferencias** (S2 2320 hojas, S3 906, S4 8451, S5 747; `model.json` 159 + 437 = 596 de 596); `model.json` y `verify3.json` no cambiaron (sha256 `1a570666…` y `d86087f3…`). Cobertura K25 ≥ 85 % en todos los paquetes con pruebas de `motor_fin`, `motor_cred`, `recetas` y `precios` (`dominio` no existe todavía).
- **Quedan 6 artefactos sin aprobar** que se escalan a Axel en la puerta 2 con sus hallazgos abiertos (§4): motor S2 (1 medio), motor S5 (1 alto que solo Axel destraba y 4 medios), contrato S4/S5 (3 medios, tras la ronda 5), U4 (1 medio, tras la ronda 5) y el paquete para Claude Design (1 medio, una frase del LEEME). Nadie los corrigió en la integración.
- **Bloqueos:** 26 nuevos de WF-1, todos abiertos, más `S1-parametros-avisos` (P1c), que sigue abierto (§6). Propuestas para Axel consolidadas en §7.

## 1. Hecho por artefacto

| Artefacto | Rutas principales | Revisión (rondas · veredicto final) | Lo que comprobó el integrador |
|---|---|---|---|
| Contrato S2 y formato común de fixtures | `docs/contratos/motores/S2-interfaz.md` 1.2, `oraculo/fixtures/FORMATO.md` 1.2, `oraculo/compara.py` | 3 · **APROBADO** (13 bajos) · `REVISION-contrato-S2-r3.md` | `compara.py` corre las 4 ramas (s2 a s5) sin error |
| Motor S2 (`motor_fin`) y oráculo | `backend/internal/motor_fin/`, `backend/cmd/motores/{main,s2}.go`, `oraculo/fin_*.py`, `oraculo/fixtures/s2/`, `oraculo/salidas/s2.json` | 3 · **NO APROBADO** (1 medio) · `REVISION-S2-r3.md` | 30 fixtures; 2320 de 2320 hojas en las 3 vías; `model.json` 159 de 159 y 43 cruces; provisionales 933; `fin_verify.py` 159 de 159; cobertura 89.2 % |
| Contrato S3 | `docs/contratos/motores/S3-interfaz.md` 1.4, `docs/estado/S3-escenarios.md` | 5 · **APROBADO en la ronda 5** · `REVISION-contrato-S3-r5.md` | — |
| Motor S3 (`motor_cred`, comparador, evaluar) y oráculo | `backend/internal/motor_cred/`, `backend/cmd/motores/s3.go`, `oraculo/cred_*.py`, `oraculo/fixtures/s3/` | 3 · **APROBADO** (6 bajos) · `REVISION-S3-r3.md` | 27 fixtures; 906 de 906; `model.json` 437 de 437; `data.json` 9 de 9; segunda vía float 174 de 174; cobertura 93.7 / 96.0 / 90.8 / 90.0 % |
| Contrato S4 y S5 | `docs/contratos/motores/S4-interfaz.md` 1.4, `S5-interfaz.md` 1.4 | 5 · **NO APROBADO tras la ronda 5** (3 medios, 4 bajos) · `REVISION-contrato-S4S5-r5.md` | — |
| Catálogo de productos | `datos/semilla/productos.json` (160 productos, 9 grupos), `productos-README.md` 1.1, `docs/contratos/motores/productos-formato.md` 1.1 | 2 · **APROBADO** (7 bajos) · `REVISION-catalogo-r2.md` | lo leen los motores S4 y S5 en sus fixtures sin error |
| Motor S4 (`recetas`), oráculo, recetario y nutrición | `backend/internal/recetas/`, `oraculo/rec_*.py`, `datos/semilla/recetario{,-a,-b}.json`, `recetario.schema.json`, `datos/semilla/nutricion/` | 2 · **APROBADO** (7 bajos) · `REVISION-S4-r2.md` | 26 fixtures; 8451 de 8451; recetario: ver §2.3; cobertura 93.2 % |
| Motor S5 (`precios`, asistida, OCR, PROFECO) y oráculo | `backend/internal/precios/`, `oraculo/pre_*.py`, `datos/semilla/profeco/`, `datos/semilla/tickets-sinteticos/`, `docs/estado/S5-cobertura-profeco.md` | 3 · **NO APROBADO** (1 alto, 4 medios) · `REVISION-S5-r3.md` | 29 fixtures; 747 de 747; `pre_verify.py` TODO COINCIDE; cobertura 96.1 / 100 / 99.3 / 95.3 % |
| U2: sistema de diseño, tokens C2 y contraste | `docs/design/DESIGN.md`, `docs/contratos/C2-tokens.json`, `tokens.css`, `tailwind.theme.css`, `docs/design/tools/tokens/`, `docs/design/contraste-componentes.{csv,md}` | 4 · **APROBADO en la ronda 4** (11 bajos) · `REVISION-U2-r4.md` | generador y compilación reproducibles; contraste regenerado idéntico (§2.4) |
| U3: Inicio (dashboard) | `docs/design/U3-dashboard.md` 1.2 | 3 · **APROBADO** (9 bajos) · `REVISION-U3-r3.md` | entra al paquete como `U3-dashboard.pdf` |
| U4: Ajustes | `docs/design/U4-ajustes-a.md` 1.4, `U4-ajustes-b.md` 1.3 y **`U4-ajustes.md` (unido aquí)** | 5 · **NO APROBADO tras la ronda 5** (1 medio) · `REVISION-U4-r5.md` | unión byte a byte comprobada; 14 secciones P-26 (§2.5) |
| Paquete para Claude Design (Pista 3, K35) | `docs/design/entrada/`, `docs/design/tools/paquete_entrada.sh`, `docs/design/tools/paquete/` | 3 · **NO APROBADO** (1 medio, 5 bajos) · `REVISION-paquete-r3.md` | regenerado con el U4 unido; `sha256sum -c` sin diferencias; sin código (§2.6 y §10) |

Además: OCR (`docs/estado/wf1/ocr.md`, `backend/internal/precios/ocr/`, tickets sintéticos) y PROFECO (`backend/internal/precios/profeco/`, base local de 37,548,200 registros) son parte del motor S5.

## 2. Verificaciones del integrador

### 2.1 Go (`backend/`)
- `gofmt -l .` → vacío · `go vet ./...` → sin errores · `go build ./...` → sin errores (Go 1.27.2).
- `go test ./... -count=1 -cover` → todos `ok`. Cobertura por paquete: `motor_fin` 89.2 · `motor_cred` 93.7 · `motor_cred/canontest` 90.0 · `motor_cred/comparador` 96.0 · `motor_cred/evaluar` 90.8 · `recetas` 93.2 · `precios` 96.1 · `precios/asistida` 100.0 · `precios/ocr` 99.3 · `precios/profeco` 95.3 · `cmd/motores` 48.7 (sin umbral) · `cmd/finanzas` y `internal/api` 0.0 (sin pruebas, sin umbral; K25 los mide sin umbral). Programas main sin pruebas dentro de los árboles K25: `precios/ocr/cmd/ocrmedir` y `precios/profeco/cobertura` (se informan). `internal/dominio` no existe todavía: el umbral se aplicará cuando exista.

### 2.2 Oráculos (desde la raíz del repositorio)
- `python3 -I oraculo/{fin,cred,rec,pre}_run_fixtures.py` → 30/30, 27/27 (79 consultas; segunda vía 174/174), 26/26 (11214 hojas de esperado) y 29/29 (1134); `oraculo/salidas/s{2,3,4,5}.json` regeneradas **idénticas** (sha256 antes y después).
- Motor: `datos/tmp/motores <sN> dump --fixtures oraculo/fixtures/<sN> --out …` y `python3 -I oraculo/compara.py <sN>` → **diferencias: 0** en los 4. S2: 2320 hojas, `model.json` 159/159, cruces 43/43, provisionales 933/933. S3: 906, `model.json` 437/437, `data.json` 9/9, provisionales 2718/2718. S4: 8451, provisionales 429/429. S5: 747, provisionales 301/301.
- `fin_verify.py` 159 de 159 (float contra exacto y contra `model.json`) · `pre_verify.py` 15 OK, TODO COINCIDE · `model.py` → «datos/tmp/model.json igual a oraculo/model.json: True» · `verify3.py` reproduce `verify3.json` byte a byte.
- `sha256sum oraculo/model.json oraculo/verify3.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` y `d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041`, iguales a HEAD y a `oraculo/README.md`. Los únicos cambios en archivos ya versionados de `oraculo/` son `model.py` y `verify3.py` (S2-04: importarlos ya no escribe; ejecutados escriben en `datos/tmp/`).

### 2.3 Recetario
- `python3 -I oraculo/rec_genera_recetario.py --verifica` → «recetario.json coincide con lo regenerado» (la unión de a + b la escribe ese generador; el integrador no la reescribió porque ya coincide y el contenido revisado no se cambia).
- JSON Schema (`recetario.schema.json`, Draft 2020-12, `jsonschema` efímero con `uv`): 0 errores · 60 recetas · 60 ids y 60 nombres únicos · `kcal_por_porcion` sin null · orden = 30 de `recetario-a.json` + 30 de `recetario-b.json`, sin ids compartidos. Los parciales se conservan.
- Motor: `go test ./internal/recetas -run 'TestRecetarioValido|TestRecetarioProduccion|TestEscaladoRecetario|TestKcalRecetasAutor'` → ok; la fixture `s4-escalado-recetario` (240 consultas sobre el recetario real) da 0 diferencias en `compara.py s4`.

### 2.4 Tokens C2 y contraste
- `generar_c2_tokens.py` → `C2-tokens.json` idéntico byte a byte · `verificar_c2.py` → RESULTADO: OK · `node sincronia.mjs` (Style Dictionary 5.6.1, dos compilaciones) → `tokens.css` sha256 `2b9429bc…` y `tailwind.theme.css` `f66eb866…`, iguales a los versionados.
- `contraste_componentes.py --check` → «CSV en disco IDENTICO al regenerado» · bajo umbral: 0 · cruce con `contraste.csv`: igual 266, nuevo 50, difiere 0.

### 2.5 U4 unido
- `docs/design/U4-ajustes.md` (143,172 bytes) = encabezado e índice del integrador + `U4-ajustes-a.md` + línea en blanco + `U4-ajustes-b.md`, sin cambios en el texto de las partes; `grep -c '^## P-26\.'` = 14. El encabezado dice que U4 no quedó aprobado y cita el bloqueo `S1-parametros-avisos`.

### 2.6 Paquete para Claude Design
- `bash docs/design/tools/paquete_entrada.sh` → solo se reconstruyó `U4-ajustes.pdf` (ahora sale de `U4-ajustes.md`; 37 páginas, 1,909,203 bytes); `pdftotext | grep -c '/api/v1'` = 0 en los 4 PDF; sin nombres de operación en U3, U4 y C1-resumen; MANIFIESTO.md regenerado (cambian solo la fila y la línea de `U4-ajustes.pdf`); `sha256sum -c` → 0 diferencias en 21 archivos (más MANIFIESTO.md). En `docs/design/entrada/` solo hay pdf, png, md, json, css y docx; 0 scripts.

### 2.7 Makefile y verificación completa
- Nuevo `make verify-wf1` (go build, vet y gofmt; `go test ./... -cover` con umbral K25 de 85 %; los 4 oráculos con `compara.py`; `fin_verify`, `pre_verify`, `model.py`, `verify3.py` y sha256; recetario; tokens y contraste; U4 unido; paquete). `verify-all` lo incluye.
- `make verify-all GOOSE=…/gobin/goose` → rc 0: «WF-0: verificación completa OK», C1 «RESULTADO: TODO OK», S1 «RESUMEN TODO OK» (goose v3.28.0 up / down-to 0 / up), U1 «TODO OK», criterios OK, «WF-1: verificación completa OK» y «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## 3. Interrupción, continuación y rondas extra

- **Pausa:** el run `wf_98b77133-bf1` se pausó a las 19:00:06 del 11-oct-2026 por un reinicio de la sesión del orquestador, con 50 agentes terminados y 11 interrumpidos.
- **Relanzados desde cero (11), con aviso de sus archivos parciales:** `motor-S3-nucleo-go`, `motor-S3-comparador-go`, `oraculo-S3`, `motor-S4-go`, `nutricion`, `oraculo-S4`, `motor-S5-go`, `profeco`, `oraculo-S5`, `revisor-S2-r1` y `revisor-paquete-r1`. Los 11 tienen su reporte en `docs/estado/wf1/<agente>.md` (marcados en la tabla de §11) y sus artefactos pasan las verificaciones de §2.
- **Rondas extra** (decisión del orquestador, precedente P1d/P1e, reportada a Axel como HECHO) para los artefactos que no se aprobaron en 3 rondas, antes de arrancar los motores que dependen de ellos:

| Artefacto | Ronda 4 | Ronda 5 | Resultado |
|---|---|---|---|
| Contrato S3 | `REVISION-contrato-S3-r4.md`: NO (0 altos, 1 medio P1, 13 bajos) → `correccion-contrato-S3-r4.md` | `REVISION-contrato-S3-r5.md`: **APROBADO** (0 altos, 0 medios, 3 bajos nuevos) | aprobado |
| Contrato S4/S5 | `REVISION-contrato-S4S5-r4.md`: NO (2 medios, 3 bajos) → `correccion-contrato-S4S5-r4.md` | `REVISION-contrato-S4S5-r5.md`: NO (0 altos, 3 medios, 4 bajos) | **se escala a Axel** |
| U2 | `REVISION-U2-r4.md`: **APROBADO** (0 altos, 0 medios, 11 bajos) | no hizo falta | aprobado |
| U4 | `REVISION-U4-r4.md`: NO (1 medio R4-1, 7 bajos) → `correccion-U4-r4.md` | `REVISION-U4-r5.md`: NO (1 medio R5-1, 2 bajos nuevos) | **se escala a Axel** |

- **Paquete regenerado** (`docs/estado/wf1/paquete-regenerar.md`) tras los cambios de U4 en las rondas extra; después vinieron `correccion-paquete-r1` y `-r2` y la revisión `REVISION-paquete-r3.md` (NO, 1 medio). El integrador lo volvió a generar con `U4-ajustes.md` unido (§2.6).
- Motor S2, motor S5 y paquete se quedaron en 3 rondas: la decisión del orquestador dio rondas extra solo a contrato S3, contrato S4/S5, U2 y U4. También se escalan (§4).

## 4. Artefactos sin aprobar que se escalan a Axel (puerta 2)

Ninguno de estos hallazgos cambia una cifra canon: los 4 motores reproducen `model.json`, `data.json` y sus fixtures con 0 diferencias (§2.2). Son reglas que el contrato no escribe o textos incompletos. El integrador no los corrigió.

| Artefacto | Hallazgos abiertos (altos y medios) | Informe |
|---|---|---|
| Motor S2 | **M1 (media) · S2-16 / S2-03:** con ventanas ascendentes y disjuntas que dejan un hueco sin eventos fechados, un devengo de n meses sigue en «n ventanas consecutivas» y la parte del hueco cae en la siguiente ventana declarada. Motor y oráculo coinciden, pero la regla no está escrita y en modo meses la misma parte cae distinto. Se liga al bloqueo `S2-ventanas` (pregunta 4). | `docs/estado/wf1/REVISION-S2-r3.md` |
| Motor S5 | **H-r3-00 (alta) · S5-07:** `docs/estado/S5-terminos.md` no existe; solo Axel lo destraba (bloqueo `S5-terminos`); la lectura asistida sigue apagada y no hay URLs de tiendas en el código. **H-r3-01 (media):** `S5-interfaz.md` 1.4 desfasado, sin las solicitudes 1 a 6 de `correccion-S5-r2` (R2-04). **H-r3-02 (media):** `compare_carts` con cantidades a granel: motor y oráculo dan cifras distintas y ninguna fixture lo prueba. **H-r3-03 (media):** `profeco_import`: las reglas de lectura y de liga difieren entre motor y oráculo (p. ej. `ATUN DOLORES` contra `ATÚN DOLORES`) y §11 y §15.2 no las fijan. **H-r3-04 (media):** `difference_vs_cheapest_cents` tiene dos significados según el modo. Además S5-06 (OCR) queda «sin medir» con honestidad: faltan las 5 fotos de K24. | `docs/estado/wf1/REVISION-S5-r3.md` |
| Contrato S4/S5 (tras 5 rondas) | **M-r5-1 (media) · S4-03, S4-06:** `S4-interfaz.md` §8.1 paso 1 y §12.1 (`quotes`): la cotización con varias elecciones. **M-r5-2 (media) · S4-03:** §6.2 paso 5 y §12.1: la presentación sin cotización. **M-r5-3 (media) · S4-03, S5-03 y oráculos:** `S4-interfaz.md` §12.2, `S5-interfaz.md` §7.2, §7.4 y §15.2, `FORMATO.md` §4 y §5. | `docs/estado/wf1/REVISION-contrato-S4S5-r5.md` |
| U4 (tras 5 rondas) | **R5-1 (media) · U4-05:** caminos de edición con vista previa que la regla de §0.2 no cubre (P-26.04 «Edición en el lugar» y «Atajos»; P-26.14 y otra sección de la parte b). Bajos abiertos: R5-2, R5-3, R4-2 («cada cuánto» sin el tope 1–24 de `interval_count`), R4-3 (PROVISIONAL sin cita en la misma línea, parte b) y R3-2 a R3-6. | `docs/estado/wf1/REVISION-U4-r5.md` |
| Paquete para Claude Design | **H-1 (media) · P3-03:** el paso 2 del LEEME no dice qué elegir si al crear el proyecto la herramienta pide tipo de proyecto o sistema de diseño (la corrección es una frase en `LEEME.md` y regenerar `MANIFIESTO.md`). 5 bajos. P3-01, P3-02, P3-03 y P3-05 cumplen en su texto. | `docs/estado/wf1/REVISION-paquete-r3.md` |

## 5. Hallazgos bajos registrados

- **De las revisiones finales** (abiertos, no bloquean; el detalle y la corrección sugerida están en cada informe): contrato S2 13 (`REVISION-contrato-S2-r3.md`); contrato S3 3 nuevos R1 a R3 más los que seguían (`REVISION-contrato-S3-r5.md`); motor S3 6, B1 a B6 (`REVISION-S3-r3.md`); motor S4 7, B2, B3, B4, B6, B7 (kcal de cortes con hueso), B8 (recetas parecidas) y B9 (`REVISION-S4-r2.md`); catálogo 7 (`REVISION-catalogo-r2.md`); contrato S4/S5 4 (`REVISION-contrato-S4S5-r5.md`); motor S5 H-r3-05 (opciones extra de Tesseract) y H-r3-06 (cifras de un solo corte en `S5-cobertura-profeco.md`) (`REVISION-S5-r3.md`); U2 11 (`REVISION-U2-r4.md`); U3 9, L-1 a L-6 y N-1 a N-3 (`REVISION-U3-r3.md`); U4 los de §4; paquete 5 (`REVISION-paquete-r3.md`).
- **Del integrador:**
  - **I-1 · motor fuera de la raíz.** `go run ./cmd/motores sN dump` corrido desde `backend/` responde `invalid_argument` en 14 fixtures (4 de S4 y 10 de S5) y aun así sale con código 0; `compara.py` da entonces 7758 y 463 diferencias. Causa: las fixtures citan `datos/semilla/…` relativo a la raíz del repositorio. Desde la raíz: 0 diferencias. `verify-wf1` construye el binario y lo corre desde la raíz. Sugerencia (dueño de `main.go`, WF-3): resolver las rutas contra la raíz o documentarlo en `FORMATO.md` §4.1.
  - **I-2 · `go.mod`.** `modernc.org/sqlite` y sus dependencias aparecen como `// indirect`, aunque `internal/precios/profeco/sqlite.go` y `profeco/cobertura/main.go` lo importan; `go mod tidy` lo corrige. No afecta build, vet ni pruebas.
  - **I-3 · PDF del paquete.** lualatex falla por la métrica de `lmroman10`; los 4 PDF salen por html + chromium y no son reproducibles byte a byte (fecha de Chromium); el script los guarda en caché por sello de fuentes. Ya estaba registrado por los agentes del paquete.
  - **I-4 · `FORMATO.md` §6** cita `oraculo/cred_compara.py` como verificación adicional de S3, pero ese archivo no existe; las 437 hojas de `model.json` las cubre `compara.py s3`.
  - **I-5 · `.gitignore`.** La regla `datos/semilla/nutricion/*.csv` de WF-0 dejaba fuera `nutricion/tabla.csv`, que leen `rec_genera_recetario.py` y las fixtures de S4; se agregó la excepción. También se ignora `.obsidian/` (configuración local de Obsidian).

## 6. Bloqueos nuevos de WF-1

26 archivos nuevos en `docs/estado/bloqueos/`, **todos abiertos** (ninguno trae «Respuesta de Axel»); lo que depende de ellos lleva la marca PROVISIONAL con su cita. Sigue abierto también `S1-parametros-avisos.md` (P1c), que nadie decidió en WF-1. Se deciden en la puerta 2; `docs/estado/bloqueos/INDICE.md` los resume en su sección «Bloqueos nuevos de WF-1».

| Archivo | Fase | Pregunta (primera línea de «Qué falta») |
|---|---|---|
| `S2-flujo-quincenal.md` | S2 | (1) método del gasto variable estimado de una quincena futura; (2) cómo se reparte el ahorro mensual planificado (`planned_monthly_cents`) entre las dos quincenas; (3) cuándo una ocurrencia planificada (sueldo, renta) se da… |
| `S2-metas-resumen.md` | S2 | ¿el resumen del mes cuenta solo las metas y fondos activos o también los inactivos que aún tienen saldo? M-07 dice «ahorrado = Σ aportaciones − Σ retiros (fondo y metas)» sin decir «activos», y en la misma fila dice… |
| `S2-metas.md` | S2 | (1) cómo se cuentan los «meses que faltan» de M-16 («aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan»); (2) qué porcentaje de avance se muestra: RN-18 dice «con un decimal» (50.0 %) y M-07/M-16 dicen… |
| `S2-presupuesto-total.md` | S2 | (1) si existe un renglón del total, ¿el plan total es ese monto o la suma de las categorías? (2) Con un plan de 0.00 en una categoría, ¿qué porcentaje y qué banderas se muestran? |
| `S2-proyeccion-interruptor.md` | S2 | (1) monto de cada despensa proyectada: ¿capital de la despensa base 4,000.00 (despensas 2–4) con la tasa vigente de la tabla por plazo, o el de la última despensa comprada, o una partida planificada que Axel captura? (2)… |
| `S2-ventanas.md` | S2 | (pregunta para Axel y para el dueño de la interfaz): 1. Con ventanas, ¿qué pasa con un costo o un pago anterior a la primera ventana pedida? §3.2 dice que `cumulative_diff_cents` va «desde el primer evento de las entradas… |
| `S3-convencion-prestamo.md` | S3 | ¿Cómo se calcula la tasa por periodo de un préstamo a plazo fijo y cómo entra el IVA? |
| `S3-formula-tarjeta.md` | S3 | (1) Intereses de tarjeta cuando no se paga el total: ¿sobre qué saldo (saldo promedio diario del periodo, saldo al corte, saldo insoluto después del pago) y con qué días (días reales del periodo / 360, / 365, o tasa anual / 12… |
| `S3-objetivos-medida.md` | S3 | (1) Canasta y tienda (objetivo 4): ¿la tienda es una dimensión más del escenario (cada tienda con precios vigentes) y qué es «la canasta» de una compra planeada? (a) la lista de compras de la semana (P-13) con los precios… |
| `S3-origen-cargo-normal.md` | S3 | ¿qué `interest_source` lleva un cargo normal de tarjeta? |
| `S3-serie-sin-intereses.md` | S3 | con intereses 0, ¿qué capital lleva cada pago? K14 y RN-07 fijan montos con base truncada y capital con base HALF_UP, y el último pago ajusta. Para 3,094.50 a 4 quincenas eso da montos 773.62 ×3 + 773.64 y capital 773.63 ×3 +… |
| `S4-asignacion-mercado.md` | S4 | (cuatro preguntas para Axel): 1. Si varios mercados surten el grupo del producto, ¿a cuál va el artículo de la lista? |
| `S4-confirmacion-empates.md` | S4 | ¿quién confirma los empates de `nutricion/tabla.csv` y las ligas de `profeco/vinculos.csv`, y en qué momento? |
| `S4-reglas-lista.md` | S4 | (tres preguntas para Axel): 1. **Sobrante previsto** (`MealPlanWeek.leftovers`, RN-15): ¿el plan cocina exactamente las porciones planeadas o la receta completa (múltiplos de sus 2 porciones base)? |
| `S5-carrito-comparativa.md` | S5 | (preguntas para Axel): 1. Con `basket = expense`, ¿contra qué se compara la compra registrada en las otras tiendas? |
| `S5-carritos-detalle.md` | S5 | tres preguntas para Axel (o para el dueño del contrato): 1. Ganador de un sustituto cuando una tienda tiene varias presentaciones del producto: ¿precio por unidad base de **lo que esa tienda cobró** por el producto (importe… |
| `S5-profeco-almacen.md` | S5 | (tres preguntas para Axel): 1. ¿Dónde viven los registros crudos de PROFECO (millones de renglones, la mayoría sin producto del catálogo)? |
| `S5-profeco-filtro.md` | S5 | tres preguntas para Axel (o para el dueño del contrato): 1. ¿El `municipality` de una tienda que **no** es sucursal (una cadena o un mercado) entra al filtro de municipios? |
| `S5-reglas-comparativa.md` | S5 | (cuatro preguntas para Axel): 1. ¿Cómo «pesa menos» un precio viejo al elegir la tienda más barata? |
| `S5-terminos.md` | S5 | Axel, ¿cómo se obtiene el veredicto de términos de los 3 sitios (walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx)? |
| `U3-t2-acumulados.md` | U3 | ¿C1 agrega a `AccrualVsCash.months[]` el devengado acumulado y el efectivo acumulado con el mismo origen que `cumulative_difference_cents` (opción A), o T-2 se dibuja relativo a la ventana con las sumas de lo visible (opción B)? |
| `U4-dia-mes-reporte.md` | U4 | si Axel cambia el día del mes del reporte mensual a 29, 30 o 31, ¿qué pasa en los meses que no tienen ese día (febrero, abril, junio, septiembre, noviembre)? |
| `U4-iva-credito.md` | U4 | C1 v1.2 y el esquema 1.3 guardan «si incluye IVA» en `credit.vat_on_interest_bp`, donde NULL significa «la tasa ya incluye IVA»; por eso `Credit.missing_terms` solo puede listar 4 datos (`annual_rate`, `rate_base`,… |
| `U4-reparto-retroactivo.md` | U4 | cuando Axel cambia el porcentaje de una regla de reparto reutilizable (p. ej. «Conjunto 50/50»), ¿el cambio vale para los gastos ya capturados o solo para los nuevos? |
| `U4-vista-previa-formulas.md` | U4 | (dos preguntas vigentes para Axel; la 2 se retiró): 1. **Comisiones de un crédito** (anual, por pago tardío y por apertura; `credit.annual_fee_cents`, `late_fee_cents`, `opening_fee_cents`): ¿deben entrar al costo de una… |
| `WF1-catalogo-huevo-por-kilo.md` | WF1 | ¿Cómo debe registrarse un precio por kilo de un producto cuya unidad base es `pieza` (huevo)? Es decisión de modelo (S5/Ajustes), no del catálogo. |

## 7. Propuestas para Axel (puerta 2), consolidadas de los reportes

### 7.1 Valores operativos (K31 y K44): PROVISIONALES hasta que Axel los apruebe
| Valor | Propuesta | Por qué (resumen) | Agente |
|---|---|---|---|
| Hora de los reportes (K31 punto 1) | diario 21:30; semanal domingo 19:00; mensual día 1 a las 08:00 | de noche ya se capturó el día; el domingo se planea la semana; el día 1 el mes ya cerró. Dónde se guarda depende de `S1-parametros-avisos` | `U4b-ajustes` |
| Avisos (K31 punto 2) | `advance_days` de precios viejos = 1; hora de avisos 09:00 (detalle en el reporte) | un día antes de la despensa da tiempo de actualizar precios | `U4b-ajustes` |
| Tiempo máximo del comparador (K31 punto 7) | mantener 5 s (K43); límites de seguridad 1 s a 30 s | peor caso medido 111–123 ms (3,870 candidatos), 2.5 % del tope; casos reales ≤ 1 ms | `motor-S3-comparador-go`, `oraculo-S3` |
| `optimizer.date_window_days` | 7 días | cubre una semana y mueve el primer pago a la quincena vecina; peor caso ≈ 7,600 escenarios, muy por debajo de 5 s | `contrato-S3` |
| `ocr.timeout_seconds` | 30 s | 2 imágenes sintéticas tardan ≈ 2 s; una foto de 12 MP en un equipo lento, 5–10 s; margen de 3× | `ocr`, `contrato-S4S5` |
| Rango de frescura (K44) | `price.stale_after_days` ≥ 1 y `price.zero_weight_after_days` ≥ `price.stale_after_days`; fórmula general del peso (Z − días) ÷ (Z − S + 1), que con 30 y 90 da K40 | con otros valores la recta no baja de 1 a 0; el motor los rechaza con `invalid_argument` | `motor-S5-go`, `U4a-ajustes` |

### 7.2 Opciones recomendadas para los bloqueos
- `S2-ventanas`: opción A + H1 (aplicada como PROVISIONAL), pregunta 3 O1 (rechazar) y pregunta 4: exigir ventanas contiguas, ascendentes y disjuntas (`correccion-S2-r1`, `correccion-S2-r2`, `revisor-S2-r3`).
- `S2-proyeccion-interruptor`: opción C; `S2-flujo-quincenal`: (1) C con 0 o A, (2) B mitad y mitad, (3) A (`contrato-S2`); `S2-metas-resumen`: A, solo metas activas (`correccion-contrato-S2-r2`).
- `S3-origen-cargo-normal`: (a) `quoted`; `S3-serie-sin-intereses`: (a) capital = monto cuando los intereses son 0 (`correccion-contrato-S3-r2`, `-r3`).
- `S4-reglas-lista`: A; `S5-reglas-comparativa`: A en 1 a 3 y 4b, B en 4a; `S5-profeco-almacen`: A (`contrato-S4S5`); `S5-carrito-comparativa`: 1 A y 3 A (`correccion-contrato-S4S5-r1`); `S4-asignacion-mercado`: A, A y A (`correccion-contrato-S4S5-r2`).
- `U3-t2-acumulados`: A (`correccion-U3-r2`); `U4-vista-previa-formulas`: (a) en las preguntas vigentes (`correccion-U4-r1`).
- `S5-terminos`: decidir si se autoriza leer solo las páginas legales de Walmart, Bodega Aurrera y Sam's Club o si Axel dicta el veredicto (`revisor-S5-r1`).

### 7.3 Otras decisiones para ratificar
- Diseño (`U2b-1-design`, `U2b-2-design`): negativos con signo menos tipográfico U+2212 antes de `$`; «rebasado» con icono en advertencia; trato de tú; ingreso en 135°; interruptor «Atajos de una tecla» (WCAG 2.1.4); indicador «Guardado» de 2 s; las constantes de interfaz (5 s, 2 s, 400 ms, 10 s) como diseño y no como parámetros de K44.
- U3: supuesto 5 (dibujo de T-2) y nombre de la segunda tienda en la tarjeta 5 (`U3-dashboard`, `correccion-U3-r1`).
- Recetas: etiquetas («económica», «rápida», «sobrantes») sí o no (`correccion-S4-r1`); kcal de cortes con hueso (a/b/c) y recetas parecidas (`revisor-S4-r2`); confirmar los 40 empates de `datos/semilla/nutricion/revisar.tsv` (`nutricion`).
- PROFECO: cadenas del catálogo (`Wal-mart`, `Bodega Aurrera`, variantes `Express`) y municipios por defecto (`profeco`); grupos iniciales del catálogo (`contrato-S4S5`); subtotal opcional «falta por comprar» en la lista (`correccion-contrato-S4S5-r4`).

## 8. Pendientes para WF-3 (y para la puerta 2)
1. Puerta 2: los 6 artefactos de §4, los 26 bloqueos de §6 más `S1-parametros-avisos`, y los valores de §7.1.
2. S5-06: medir el OCR con las 5 fotos reales de Axel en `datos/semilla/tickets/` (K24) con `go run ./internal/precios/ocr/cmd/ocrmedir`.
3. S5-07: `docs/estado/S5-terminos.md` con el veredicto de los 3 sitios; hasta entonces la lectura asistida sigue apagada.
4. `backend/internal/dominio`: nace con cobertura ≥ 85 % (K25); `make verify-wf1` ya lo exige en cuanto exista.
5. I-1 (rutas del motor contra la raíz) e I-2 (`go mod tidy`).
6. Base de PROFECO (`datos/semilla/profeco/db/qqp.sqlite`, 17 GB): dónde vive en producción (`S5-profeco-almacen`); confirmar empates de nutrición y ligas de PROFECO en estado `propuesto` (`S4-confirmacion-empates`).
7. Las solicitudes fuera de alcance de cada reporte que tocan artefactos no aprobados se resuelven con las respuestas de la puerta 2; las demás quedan anotadas en cada `docs/estado/wf1/<agente>.md`.
8. WF-3 (PWA): cola sin conexión (K9), pantallas sobre los tokens C2 y DESIGN.md, y lo que devuelva la Pista 3 en `docs/design/salida/`.

## 9. Paquetes instalados y descargas
- **Del sistema:** ninguno en WF-1 (tesseract 5.5.3, ImageMagick, bsdtar, unrar, sqlite3 y Go 1.27.2 ya estaban; reportes `ocr` y `profeco`).
- **npm** (solo en `docs/design/tools/tokens/node_modules`, ignorado por git): style-dictionary 5.6.1, tailwindcss 4.3.3 y @tailwindcss/cli 4.3.3, fijados en `package-lock.json` (`U2c-tokens`).
- **Python efímero con `uv`:** `jsonschema` (validación del recetario en `verify-wf1`); nada se instala en el sistema.
- **Descargas** (PROFECO, USDA FoodData Central y Open Food Facts; todas ignoradas por git; procedencia en `datos/semilla/nutricion/PROCEDENCIA.md` y `docs/estado/S5-cobertura-profeco.md`):

| Ruta | Tamaño |
|---|---|
| `datos/semilla/profeco/raw/2024-rar/` | 91.6 MB |
| `datos/semilla/profeco/raw/2025-rar/` | 81.2 MB |
| `datos/semilla/profeco/raw/2026-rar/` | 57.2 MB |
| `datos/semilla/profeco/raw/2024-csv/` | 4.4 GB |
| `datos/semilla/profeco/raw/2025-csv/` | 3.9 GB |
| `datos/semilla/profeco/raw/2026-csv/` | 2.8 GB |
| `datos/semilla/profeco/db/` | 16.0 GB |
| `datos/semilla/nutricion/raw/usda-foundation-2026-04-30/` | 3.6 MB |
| `datos/semilla/nutricion/raw/usda-srlegacy-2018-04/` | 5.8 MB |
| `datos/semilla/nutricion/raw/off-2026-10-11/` | 1.2 GB |

Lo que sí entra al repositorio de esas fuentes: `datos/semilla/profeco/muestra/` (≈ 604 KB), `datos/semilla/nutricion/tabla.csv` y sus tablas de trabajo, los scripts de procesamiento (fuera de `raw/`) y los `.gitignore` de `raw/` y `db/`.

## 10. Paquete para Claude Design (Pista 3, K35: Axel lo sube él mismo)
- **Carpeta:** `/home/deepsec/finanzas-hogar/docs/design/entrada/` · 22 archivos · 9,634,361 bytes en total. La guía es `LEEME.md`; las huellas, `MANIFIESTO.md` (`sha256sum -c` sin diferencias). Estado: **no aprobado** en la ronda 3 por una frase del paso 2 (§4); el resto de los criterios P3 cumple.

| Archivo | Bytes |
|---|---:|
| `C1-resumen.pdf` | 280,159 |
| `C2-tokens.json` | 100,073 |
| `DESIGN.docx` | 82,222 |
| `DESIGN.pdf` | 2,859,638 |
| `LEEME.md` | 8,212 |
| `MANIFIESTO.md` | 6,956 |
| `paleta.png` | 865,657 |
| `prompt.md` | 12,284 |
| `reporte-p01.png` | 193,383 |
| `reporte-p02.png` | 126,390 |
| `reporte-p03.png` | 217,623 |
| `reporte-p04.png` | 224,544 |
| `reporte-p05.png` | 213,182 |
| `reporte-p06.png` | 174,292 |
| `reporte-p07.png` | 247,150 |
| `reporte-p08.png` | 189,586 |
| `reporte-p09.png` | 194,116 |
| `reporte-p10.png` | 229,803 |
| `reporte-p11.png` | 310,313 |
| `tokens.css` | 26,121 |
| `U3-dashboard.pdf` | 1,163,454 |
| `U4-ajustes.pdf` | 1,909,203 |

**Pasos de `LEEME.md`, resumidos:** (1) abrir `claude.ai/design` con la sesión de Axel; (2) crear el proyecto «Finanzas Hogar»; (3) copiar el texto de `prompt.md` al cuadro de mensaje sin enviarlo; (4) subir `DESIGN.pdf`, `U3-dashboard.pdf`, `U4-ajustes.pdf`, `C1-resumen.pdf`, `C2-tokens.json`, `tokens.css`, `paleta.png` y las 11 páginas `reporte-p01..p11.png`; (5) comprobar 18 adjuntos y enviar; (6) esperar sin reenviar y, si Claude Design pregunta o propone un plan, no contestar: copiarlo y avisar en la sesión de Claude Code; (7) exportar en PDF (y HTML independiente o imágenes si se puede); (8) guardar en `docs/design/salida/iter-1/` sin renombrar; (9) avisar «iter-1 lista». Hasta 3 vueltas: un agente revisa cada una y deja `siguiente-prompt.md`; lo aceptado va a `docs/design/salida/final/`. Si algo falla 2 o 3 veces, detenerse y avisar; `DESIGN.docx` es el respaldo si no acepta el PDF; el proyecto sigue con DESIGN.md aunque la Pista 3 no salga.

## 11. Acciones y KB por agente (de la sección «Consumo aproximado» de cada reporte)

97 reportes de agente en `docs/estado/wf1/` más este integrador. Cifras aproximadas declaradas por cada agente (las de 11 reportes con formato libre se leyeron a mano). Suma: 3333 acciones, ≈ 17990 KB leídos y ≈ 5046 KB escritos. Ninguno declara `cerrado_por_tope`. Por encima de algún límite del protocolo §4 (60 acciones, 400 KB leídos, 150 KB escritos): ninguno. El consumo de tokens lo calcula el orquestador aparte.

| Agente | Acciones | KB leídos | KB escritos | Cerrado por tope | Relanzado tras la pausa |
|---|---:|---:|---:|---|---|
| `U2b-1-design` | 35 | 165 | 68 |  |  |
| `U2b-2-design` | 30 | 165 | 135 |  |  |
| `U2c-tokens` | 36 | 100 | 62 |  |  |
| `U2d-contraste` | 16 | 115 | 36 |  |  |
| `U3-dashboard` | 39 | 197 | 80 |  |  |
| `U4a-ajustes` | 44 | 205 | 78 |  |  |
| `U4b-ajustes` | 38 | 215 | 75 |  |  |
| `catalogo-productos` | 22 | 95 | 135 |  |  |
| `contrato-S2` | 38 | 180 | 67 |  |  |
| `contrato-S3` | 45 | 185 | 66 |  |  |
| `contrato-S4S5` | 51 | 210 | 150 |  |  |
| `correccion-S2-r1` | 36 | 150 | 40 |  |  |
| `correccion-S2-r2` | 33 | 150 | 22 |  |  |
| `correccion-S3-r1` | 40 | 125 | 25 |  |  |
| `correccion-S3-r2` | 24 | 110 | 12 |  |  |
| `correccion-S4-r1` | 38 | 160 | 50 |  |  |
| `correccion-S5-r1` | 52 | 180 | 125 |  |  |
| `correccion-S5-r2` | 37 | 170 | 38 |  |  |
| `correccion-U2-r1` | 38 | 195 | 45 |  |  |
| `correccion-U2-r2` | 15 | 85 | 5 |  |  |
| `correccion-U2-r3` | 16 | 95 | 7 |  |  |
| `correccion-U3-r1` | 34 | 180 | 105 |  |  |
| `correccion-U3-r2` | 31 | 118 | 34 |  |  |
| `correccion-U4-r1` | 55 | 190 | 35 |  |  |
| `correccion-U4-r2` | 34 | 170 | 42 |  |  |
| `correccion-U4-r3` | 29 | 125 | 32 |  |  |
| `correccion-U4-r4` | 21 | 125 | 30 |  |  |
| `correccion-catalogo-r1` | 15 | 55 | 20 |  |  |
| `correccion-contrato-S2-r1` | 38 | 215 | 75 |  |  |
| `correccion-contrato-S2-r2` | 31 | 195 | 62 |  |  |
| `correccion-contrato-S3-r1` | 35 | 178 | 97 |  |  |
| `correccion-contrato-S3-r2` | 32 | 205 | 57 |  |  |
| `correccion-contrato-S3-r3` | 41 | 250 | 50 |  |  |
| `correccion-contrato-S3-r4` | 29 | 140 | 19 |  |  |
| `correccion-contrato-S4S5-r1` | 37 | 185 | 63 |  |  |
| `correccion-contrato-S4S5-r2` | 41 | 170 | 45 |  |  |
| `correccion-contrato-S4S5-r3` | 33 | 200 | 45 |  |  |
| `correccion-contrato-S4S5-r4` | 30 | 190 | 45 |  |  |
| `correccion-paquete-r1` | 32 | 75 | 30 |  |  |
| `correccion-paquete-r2` | 8 | 45 | 12 |  |  |
| `motor-S2-go` | 36 | 215 | 148 |  |  |
| `motor-S3-comparador-go` | 39 | 155 | 118 |  | sí |
| `motor-S3-nucleo-go` | 42 | 285 | 100 |  | sí |
| `motor-S4-go` | 33 | 190 | 145 |  | sí |
| `motor-S5-go` | 41 | 190 | 110 |  | sí |
| `nutricion` | 27 | 190 | 55 |  | sí |
| `ocr` | 23 | 65 | 75 |  |  |
| `oraculo-S2` | 34 | 232 | 141 |  |  |
| `oraculo-S3` | 41 | 300 | 70 |  | sí |
| `oraculo-S4` | 26 | 215 | 70 |  | sí |
| `oraculo-S5` | 40 | 175 | 110 |  | sí |
| `paquete-claude-design` | 38 | 150 | 80 |  |  |
| `paquete-regenerar` | 4 | 35 | 3 |  |  |
| `profeco` | 34 | 120 | 100 |  | sí |
| `recetario-a` | 12 | 85 | 125 |  |  |
| `recetario-b` | 18 | 150 | 150 |  |  |
| `revisor-S2-r1` | 44 | 190 | 35 |  | sí |
| `revisor-S2-r2` | 40 | 150 | 40 |  |  |
| `revisor-S2-r3` | 37 | 200 | 25 |  |  |
| `revisor-S3-r1` | 40 | 185 | 17 |  |  |
| `revisor-S3-r2` | 40 | 150 | 28 |  |  |
| `revisor-S3-r3` | 42 | 165 | 22 |  |  |
| `revisor-S4-r1` | 37 | 215 | 40 |  |  |
| `revisor-S4-r2` | 38 | 200 | 37 |  |  |
| `revisor-S5-r1` | 33 | 165 | 21 |  |  |
| `revisor-S5-r2` | 52 | 115 | 40 |  |  |
| `revisor-S5-r3` | 42 | 150 | 45 |  |  |
| `revisor-U2-r1` | 41 | 280 | 42 |  |  |
| `revisor-U2-r2` | 49 | 280 | 32 |  |  |
| `revisor-U2-r3` | 37 | 330 | 32 |  |  |
| `revisor-U2-r4` | 46 | 285 | 24 |  |  |
| `revisor-U3-r1` | 36 | 180 | 35 |  |  |
| `revisor-U3-r2` | 39 | 170 | 30 |  |  |
| `revisor-U3-r3` | 46 | 240 | 27 |  |  |
| `revisor-U4-r1` | 33 | 265 | 30 |  |  |
| `revisor-U4-r2` | 45 | 330 | 27 |  |  |
| `revisor-U4-r3` | 41 | 290 | 30 |  |  |
| `revisor-U4-r4` | 31 | 265 | 30 |  |  |
| `revisor-U4-r5` | 27 | 290 | 33 |  |  |
| `revisor-catalogo-r1` | 19 | 150 | 25 |  |  |
| `revisor-catalogo-r2` | 29 | 180 | 35 |  |  |
| `revisor-contrato-S2-r1` | 34 | 215 | 32 |  |  |
| `revisor-contrato-S2-r2` | 32 | 180 | 35 |  |  |
| `revisor-contrato-S2-r3` | 29 | 225 | 32 |  |  |
| `revisor-contrato-S3-r1` | 26 | 160 | 31 |  |  |
| `revisor-contrato-S3-r2` | 29 | 255 | 31 |  |  |
| `revisor-contrato-S3-r3` | 47 | 225 | 33 |  |  |
| `revisor-contrato-S3-r4` | 27 | 245 | 20 |  |  |
| `revisor-contrato-S3-r5` | 31 | 255 | 26 |  |  |
| `revisor-contrato-S4S5-r1` | 39 | 215 | 27 |  |  |
| `revisor-contrato-S4S5-r2` | 39 | 265 | 30 |  |  |
| `revisor-contrato-S4S5-r3` | 33 | 230 | 36 |  |  |
| `revisor-contrato-S4S5-r4` | 33 | 270 | 26 |  |  |
| `revisor-contrato-S4S5-r5` | 32 | 215 | 24 |  |  |
| `revisor-paquete-r1` | 37 | 225 | 23 |  | sí |
| `revisor-paquete-r2` | 43 | 160 | 22 |  |  |
| `revisor-paquete-r3` | 41 | 215 | 17 |  |  |
| `integrador-wf1` | ≈ 52 | ≈ 175 | ≈ 60 (más 143 KB de `U4-ajustes.md` generado por concatenación) |  |  |
