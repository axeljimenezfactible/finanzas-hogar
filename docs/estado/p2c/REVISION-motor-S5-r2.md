# P2c · REVISION-motor-S5-r2 · 11-oct-2026, noche

Revisor independiente `revisor-motor-S5-r2`, ronda 2 de 3 del flujo P2c «Aplicación de la puerta 2».
- No escribí nada de lo que reviso y no corregí nada.
- Evalué solo contra `docs/criterios/S5.md` 1.2 (completo) y `docs/contratos/motores/S5-interfaz.md` 1.5.
- Consideré las respuestas de Axel: `docs/estado/puerta2/RESPUESTAS.md` y las recomendadas de `PUERTA-2.md`.
- Objeto: el estado actual del trabajo de `motor-S5-puerta2` y de `correccion-motor-S5-r1`. Incluye el diff de `backend/internal/precios/` (con `profeco/`), `backend/cmd/motores/s5.go`, `oraculo/pre_*.py`, las fixtures s5 nuevas y el bloqueo nuevo `docs/estado/bloqueos/S5-motor-oraculo-reglas.md`.
- Todo lo reejecuté con mis propios comandos.
- Mis casos los armé en una raíz aislada, `datos/tmp/rev-s5-r2/root/`, con una copia de `oraculo/*.py` y de `productos.json`. Así no toqué `oraculo/fixtures/s5/_datos/`. Ya la borré.

## Veredicto: NO APROBADO (0 altos, 1 medio, 5 bajos)

**HECHO:**
- Los tres hallazgos medios de la ronda 1 quedaron cerrados.
  - M-1 y M-2 quedan como PROVISIONAL «en uso: el motor», con un bloqueo de 8 preguntas que trae cifras y opciones. Las marcas están en los dos lados.
  - M-3: la nota de S5-09 recuperó su marca PROVISIONAL y la prueba ya comprueba el texto.
- Mis casos reproducen, cifra por cifra, cada divergencia que el bloqueo declara.
- La cobertura, gofmt, vet, las fixtures, `compara.py s5`, el canon, `data.json` y `make verify-all` están bien.

**Lo que impide aprobar (M-1 de esta ronda):** el motor y el oráculo divergen con un `created_at` RFC 3339 que trae zona horaria.
- El oráculo rechaza toda la fixture con `invalid_date`. El motor la acepta, y esa zona decide el desempate S-8.
- La corrección r1 lo encontró, pero lo dejó fuera de su alcance porque supuso que `oraculo/pre_comun.py` lo comparten los oráculos S2 a S5. No es así: solo lo importa el oráculo S5 (`pre_*.py`).

## Verificaciones (evidencia propia)

| # | Verificación | Resultado | Evidencia (mis comandos, desde `/home/deepsec/finanzas-hogar`) |
|---|---|---|---|
| a1 | gofmt y vet (S5-13) | CUMPLE | `gofmt -l internal/precios cmd/motores`: vacío. `go vet ./internal/precios/... ./cmd/motores`: limpio |
| a2 | Cobertura K25 ≥ 85 % (S5-13) | CUMPLE | `go test -count=1 -coverprofile`: precios 96.0 %, asistida 100.0 %, ocr 99.3 %, profeco 95.3 %. Perfil total 88.8 % (con los dos `main` sin pruebas) |
| a3 | Sin `float` en dinero ni red (S5-08, S5-13) | CUMPLE | `float64` solo en `ocr/tsv.go:13,64,68`, que es la confianza de Tesseract y no dinero. `net/http`: 0 líneas fuera de pruebas. Ninguna URL de las tiendas en código Go |
| a4 | Pruebas nombradas por los criterios | CUMPLE | `-run 'Comparativa\|Membresia\|Unidad\|MejorTienda'`: 13 PASS (entre ellas TestComparativaCanon, TestComparativaPiezas y TestMembresiaSinCosto). `-run 'Frescura\|Fuente\|LecturaAsistida\|Mercado\|PreciosViejos'`: ok. Python imprime `[0.9836, 0.4918, 0.0]` |
| b1 | Las fixtures se regeneran idénticas | CUMPLE | `pre_genera_fixtures.py` da 32 fixtures y `diff -r` contra mi copia previa sale vacío |
| b2 | `run_fixtures` reproduce `salidas/s5.json` | CUMPLE | `pre_run_fixtures.py --out datos/tmp/rev-s5-r2/oraculo-s5.json`: 32 fixtures, 1336 hojas, 0 diferencias. `cmp` contra `oraculo/salidas/s5.json`: idéntico |
| b3 | `compara.py s5` | CUMPLE | Con el binario propio `datos/tmp/rev-s5-r2/motores s5 dump`: oráculo contra motor 935 de 935, provisionales 315 de 315, **diferencias: 0** |
| c1 | Granel 0.5 a $1.01 y $1.03 | CUMPLE | `rv2-granel` g0: totales 101 y 103, líneas 51 y 52, unidad base 135, 137, 721 y 736, Bodega gana 2. Igual en el motor, el oráculo y mi cálculo a mano |
| c2 | Granel propio | CUMPLE en dinero. `piece_count` diverge solo donde el bloqueo lo declara | g1: tres líneas de 0.5 contra una de 1.5, 151.5 → 152 y 154.5 → 155 (un solo redondeo), igual en los dos lados. g2: 1.01 × 22.01 = 2223.01 → 2223; Walmart 1.01 × 21.99 = 2220.99 → 2221 más un excedente 0.02 → 44; total 2264.97 → 2265; gana Walmart y la diferencia es +42, igual en los dos lados. `piece_count` de Walmart: motor 1 y oráculo 2 (pregunta 1, caso `piezas-2`). p1 (el `piezas-1` del bloqueo): Bodega da motor 2 y oráculo 3. s4: sustituto 0.333 × 21.00 = 699.3 → 699 contra 700, unidad base 2800 = 2800, sin ganador, igual en los dos lados |
| c3 | `profeco_import` con 21 registros propios | CUMPLE | Mis 21 registros cubren los 17 del tipo del revisor: precio 0, «1,021.00», cadena vacía, ATUN/ATÚN y «Wal-Mart»/«Gustavo A. Madero». Además van: producto vacío, precio vacío, «1,02.00», «12.345», «12.», «$ 25.00», fecha vacía, «ayer», 2026/10/13, municipio fuera, liga propuesta, sin liga, duplicado y 14 columnas. Producción: read 21, rejected 12, filtered_out 1, skipped 2, loaded 5, loaded_proposed 0 y duplicates 1. Revisión: skipped 1 y loaded_proposed 1. Se cumple la identidad de §11. `coverage` lleva el texto tal como llega (BODEGA AURRERA 3 o 4, WALMART 1 y Wal-Mart 1, en orden de bytes). Los 12 motivos salen de la lista cerrada. El motor y el oráculo son iguales hoja por hoja, y coinciden con mi esperado escrito antes de correr |
| c4 | `difference_vs_cheapest_cents` con un solo significado | CUMPLE | `rv2-canasta`, con `compare_carts` y con `compare_basket`: bodega 0, walmart +200, tercera (incompleta y con el menor total, 1000) al final y `null`; `cheapest_store` = bodega. `rv2-canasta-ninguna`: todas `null` y `cheapest_store` `null`. `go doc`: `DifferenceVsCheapestCents *int64` |
| c5 | Canon intacto | CUMPLE | Dump de `s5-comparativa-canon`: totales 309450 y 323300 (**+13850**), idénticos 277350 y 287400 (**+10050**), piezas 61/60 e `identical_pieces` 55, renglones 33/3, ganadores 22/3 y `ties` 8 |
| c6 | Divergencias declaradas en `S5-motor-oraculo-reglas` | CONFIRMADAS | Sobre 604 hojas, mi `dif.py` encontró 13 diferencias, todas declaradas: pregunta 2 (8 hojas; además la línea `only_one` da 15714 contra 15000), pregunta 3 (sobra la referencia del 01-oct en el oráculo) y pregunta 1 (2 hojas). Pregunta 4: el motor acepta `2026-8-3`, `2026/8/3`, `3/8/2026`, `20260831` y `…T10:20:30Z`, y el oráculo acepta `2026-08-31 10:20` y `31/08/2026 10:20:30`; los dos aceptan la fracción de segundo y `2026/08/31`. Pregunta 6: 105 rechazos dan 100 en el motor y 105 en el oráculo. Pregunta 5: ver B-5 |
| c7 | `created_at` RFC 3339 (§15.1) | **NO CUMPLE** (M-1) | `rv2-instante-0` (`-06:00`) y `rv2-instante-2` (`+00:00`): el motor da un precio y el oráculo responde `invalid_date`. Con `Z` y con `.5Z` coinciden |
| d1 | `data.json` intacto | CUMPLE | `git diff --quiet HEAD -- datos/semilla/data.json`: sin cambios. sha256 `2fe7afae…1031d31b` |
| e1 | S5-09: nota de membresía | CUMPLE | `extras.go:90` dice «sin costo capturado (PROVISIONAL, S5-09)». `precios_test.go:612` comprueba el texto exacto. El dump emite una sola `note` con ese texto |
| e2 | Marcas PROVISIONAL de las 8 preguntas | CUMPLE | `grep S5-motor-oraculo-reglas`: aparece en compare.go:76 y 105, pricing.go:226, profeco/parse.go:17 y 59, links.go:134, import.go:51 y s5.go:222; y en pre_calculo.py:73, 136 y 146 y pre_profeco.py:22, 67, 154, 174 y 183 |
| e3 | `provisional` de las fixtures nuevas | CUMPLE | `s5-carritos-granel` cita S5-carritos-granel; `s5-profeco-bordes` y `s5-profeco-motivos` citan S5-profeco-filtro. Los dos bloqueos siguen ABIERTOS |
| f1 | Verificación global | CUMPLE | `make verify-all GOOSE=…/goose`: EXIT=0, «WF-1: verificación completa OK» y «Verificación ampliada OK». Después, las fixtures y `salidas/s5.json` seguían idénticas |

## Hallazgos

**M-1 [media] `created_at` con zona horaria: el oráculo rechaza un instante RFC 3339 válido que el motor acepta, y el caso quedó fuera del alcance por una premisa falsa.**
- Rutas:
  - `oraculo/pre_comun.py:39`: `_RFC = …(\.\d+)?Z$`, que solo admite el sufijo `Z`.
  - `oraculo/pre_comun.py:51-54` (`instante`).
  - `backend/internal/precios/eval.go:259`: `time.Parse(time.RFC3339, …)`, que admite `±hh:mm`.
  - `S5-interfaz.md` §15.1 (`created_at` «instante RFC 3339») y §6 (S-8: a igual fecha rige el `created_at` más reciente).
- Evidencia (mis fixtures `rv2-instante-*`): dos precios que compiten, del mismo día y la misma tienda. Uno es `manual` a 2000 con `created_at` 12:30:00Z.
  - Si el otro es `ticket_manual` a 2100 con `07:00:00-06:00` (13:00Z), el motor da 2100 `ticket_manual` y el oráculo da `{"error": "invalid_date"}` para toda la fixture.
  - Con `12:00:00+00:00`, el motor da 2000 `manual` y el oráculo da `invalid_date`.
  - Con `Z` y con `.5Z`, los dos lados coinciden.
- `correccion-motor-S5-r1.md` (HECHO 5 y solicitud 4) dice que `pre_comun.py` «lo comparten los oráculos S2 a S5» y lo deja como solicitud. No es así:
  - `grep -rln pre_comun --include='*.py' .` solo da `pre_calculo.py`, `pre_datos.py`, `pre_profeco.py` y `pre_run_fixtures.py`.
  - Es del oráculo S5 y entra en el alcance `oraculo/pre_*.py`.
  - No hay bloqueo, y ninguna fixture usa zona horaria, así que `compara.py s5` no lo ve (§15.3).
  - No es una regla sin decidir: §15.1 dice RFC 3339.
- **Corrección:**
  - Que `instante` acepte `Z` y `±hh:mm`, con o sin fracción, escrito desde §15.1 sin leer el Go.
  - Agregar una fixture en la que la zona horaria decida S-8, con las cifras de arriba: 2100 `ticket_manual` con `-06:00` y 2000 `manual` con `+00:00`.
  - Si contratos prefiere restringir `created_at` a `Z`, se le hace la solicitud para §15.1, y entonces el que se alinea es el motor.

**B-1 [baja, sigue de r1] Independencia según la letra de §15.3.**
- La corrección r1 también tocó los dos lados, aunque solo en comentarios: la salida es idéntica.
- Corrección: la misma de r1. El integrador registra la excepción del reparto de P2c.

**B-2 [baja, sigue de r1] Faltan marcas PROVISIONAL.**
- `compare.go` no cita `S5-carritos-granel` (importe a granel, §7.2 paso 8).
- `profeco/links.go` y `profeco/filter.go` no citan `S5-profeco-filtro` (Q2-56 a Q2-58).
- La corrección r1 no las agregó porque no estaban en su encargo.

**B-3 [baja, sigue de r1] Artefactos fuera del alcance del motor.**
- La tabla de §15 no lista `s5-carritos-granel`, `s5-profeco-bordes` ni `s5-profeco-motivos`.
- «Cómo seguí» de `S5-carritos-granel.md` está desactualizado.
- Los dos son del integrador o de contratos.

**B-4 [baja] Índice de bloqueos.**
- `docs/estado/bloqueos/INDICE.md` no lista `S5-motor-oraculo-reglas` (`grep` sin resultados).
- La corrección ya lo pidió al integrador (solicitud 3). No es del alcance del motor.

**B-5 [baja] Pregunta 5 del bloqueo sin medir.**
- El bloqueo dice «No medí qué pasa con un id que no está en el catálogo».
- Lo medí con `rv2-profeco-desconocida`, mis 21 registros y una liga confirmada a `no-existe.x-1l`:
  - El motor responde `{"error": "unknown_presentation"}` para toda la consulta.
  - El oráculo da read 21, rejected 12, skipped 6 y loaded 0.
- Corrección: el integrador o contratos suman esta cifra a la pregunta 5.

## Estado de los hallazgos de la ronda 1

| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| M-1 [media] `piece_count` a granel (S5-01, S5-02; §7.2 paso 6) | CERRADO como PROVISIONAL con bloqueo | El bloqueo `S5-motor-oraculo-reglas`, pregunta 1, trae las opciones A/B/C, las cifras y el impacto, y sigue ABIERTO. Las marcas están en `compare.go:105` y `pre_calculo.py:146`, y se pidió a contratos que fije §7.2 paso 6. Mis g2 y p1 reproducen la tabla del bloqueo: es la única diferencia en dinero cero. Queda una fixture pendiente «con la regla fijada», la misma condición que pidió r1 para M-2 |
| M-2 [media] divergencias y supuestos copiados (§6, §7.2 pasos 1 y 2, §11, §15.3) | CERRADO como PROVISIONAL con bloqueo | (a) pregunta 2: reproduje las 8 hojas. (b) pregunta 3: 2 contra 3. (c) preguntas 4, 5 y 6: reproduje las fechas y el tope; la 5 la medí yo (B-5). (d) preguntas 7 y 8: el oráculo dice «copiado del motor en uso» (`pre_profeco.py:67,154`). Las solicitudes a contratos están en `correccion-motor-S5-r1.md`. El siguiente implementador ya no tiene que adivinar: hay regla en uso declarada y una pregunta abierta |
| M-3 [media] S5-09, nota sin marca | CERRADO | `extras.go:90` y `precios_test.go:612` (texto exacto, PASS). El dump emite «sin costo capturado (PROVISIONAL, S5-09)». La desalineación entre el criterio S5-09 y §8 pasó al integrador (solicitud 2) |
| B-1, B-2 y B-3 [baja] | SIGUEN ABIERTOS | ver arriba |

## Separación
- **HECHO:** todo lo de la tabla de verificaciones, medido por mí.
- **PREGUNTA:** ninguna nueva para Axel. Las preguntas 1 a 3 de `S5-motor-oraculo-reglas` siguen sin recomendación, y no las decido.
- **PROPUESTA:** corregir M-1 en el oráculo (`pre_comun.py`) con una fixture de zona horaria. Las bajas van al integrador.
