# revisor-S4-r1 · reporte de agente (WF-1, 11-oct-2026)

Revisor independiente, ronda 1 de 3. No corregí nada. El veredicto y los hallazgos están en `docs/estado/wf1/REVISION-S4-r1.md`.

## Hecho
- Leí el protocolo, los criterios S4 1.2 completos, las decisiones §D y §K completas, el contrato S4 1.4 en sus §1–§4 y §9–§12, RN-17, y los reportes de motor-S4-go, oraculo-S4, recetario-a, recetario-b y nutricion.
- Repetí cada verificación con mis propios comandos y scripts en `datos/tmp/rev-s4/` (ya borrados):
  - `valida.py`: valida el recetario contra §3.2 y busca repetidos;
  - `fixture.py` y `cmp.py`: una fixture propia con las 60 entradas, que suma 361 consultas de escalado, kcal y cobertura, comparada con mis cálculos hechos con `fractions`;
  - `salud.py`: una fixture propia de salud con 9 personas ficticias;
  - `usda.py`: cotejo contra los CSV oficiales de USDA.

## S4 verificación manual (S4-09)
Elegí las 5 recetas con `random.Random(41151011)` sobre las 60 entradas de a+b. Datos de `tabla.csv` (todos `propuesto`). Gramos de ml = ml × `gramos_por_unidad`. kcal = g × kcal/100 g / 100. Un solo redondeo HALF_UP.
1. **a:chuletas-de-cerdo-en-guajillo** (2 porciones). Chuleta 400 g × 1.70 = 680 (USDA 168238); guajillo 20 × 3.24 = 64.8 (168570); tomate 150 × 0.18 = 27 (170457); ajo 5 × 1.43 = 7.15 (1104647); comino 1 × 3.75 = 3.75 (170923); cebolla 40 × 0.35 = 14 (1104962); aceite de soya 15 ml × 0.921 = 13.815 g × 8.84 = 122.1246 (171411); sal 0 (746775); papa cambray 300 × 0.69 = 207 (170028). Σ = 1125.8246, total **1126**, por porción 562.91 → **563**. El motor da 563. A 4 porciones todo se multiplica por 2 exacto (chuleta 800000, aceite 30000…). Pasos completos (5) y cantidades razonables para 2 adultos.
2. **b:tilapia-al-mojo-de-ajo**. Tilapia 400 × 0.96 = 384 (175176); ajo 20 × 1.43 = 28.6; mantequilla 25 × 7.17 = 179.25 (173410); limón 60 × 0.30 = 18 (168155); perejil 10 × 0.36 = 3.6 (170416); sal 0; pimienta 1 × 2.51 = 2.51 (170931); aceite 10 ml × 0.921 × 8.84 = 81.4164; papa 300 × 0.69 = 207; lechuga 80 × 0.17458 = 13.9664 (2346389). Σ = 918.3428, total **918**, por porción **459**. El motor da 459. 2→4 exacto. Pasos completos (5).
3. **b:albondigas-en-chipotle**. Carne molida 400 × 1.90135 = 760.54 (2514743); huevo 1 pieza × 50 g × 1.43 = 71.5 (171287); pan molido 40 × 3.95 = 158 (174928); jitomate 400 × 0.18 = 72; chipotle 40 × 0.49 = 19.6 (OFF 7501017006021); cebolla 100 × 0.35 = 35; ajo 6 × 1.43 = 8.58; zanahoria 100 × 0.41 = 41 (170393); aceite 122.1246; sal 0; comino 3.75; cilantro 10 × 0.23 = 2.3 (169997); tortilla 150 × 2.18 = 327 (175036). Σ = 1621.3946, total **1621**, por porción 810.70 → **811**. El motor da 811. La sal de los pasos (3 g + 2 g) cuadra con los 5 g de la lista.
4. **a:licuado-de-platano-y-avena**. Leche 400 ml × 1.031 = 412.4 g × 0.61 = 251.564 (171265); plátano 240 × 0.97 = 232.8 (790991); avena 40 × 3.78866 = 151.5464 (2346396); azúcar 15 × 3.85 = 57.75 (334247); canela 1 × 2.47 = 2.47 (171320). Σ = 696.1304, total **696**, por porción **348**. El motor da 348.
5. **b:espagueti-con-atun-y-elote**. Espagueti 250 × 3.71 = 927.5 (169736); atún 140 × 0.90 = 126 (334194); elote 110 × 0.67 = 73.7 (169214); jitomate 300 × 0.18 = 54; cebolla 60 × 0.35 = 21; ajo 4 × 1.43 = 5.72; aceite de oliva 15 ml × 0.913 = 13.695 g × 8.84 = 121.0638 (171413); sal 0; orégano 1 × 2.65 = 2.65 (171328). Σ = 1331.6338, total **1332**, por porción **666**. El motor da 666. La fixture del autor dice 693 porque usa kcal ilustrativas (hallazgo M3).

Resultado: 5 de 5 coinciden, en kcal por porción y en escalado 2→4.

## S4-07 a mano (personas ficticias, hoy 2026-10-11)
- **h45**: hombre, nació el 1981-03-15 (45 años), último peso 92400 g (el del 2026-10-05, no el de 95000 g del 2026-09-01), 1780 mm, actividad `light`, objetivo `lose`, ajuste propio 400. MSJ = 924 + 1112.5 − 225 + 5 = 1816.5. × 11/8 = 2497.6875. − 400 = 2097.6875, meta **2098**. IMC = 92400 × 10000 / 1780² = 291.63, se muestra **29.2**, sobrepeso.
- **m32**: mujer, nació el 1993-10-12 (cumple mañana, así que tiene **32**), 55000 g, 1600 mm, `very_active`, `gain`, ajuste NULL, casa = 150. MSJ = 550 + 1000 − 160 − 161 = 1229. × 1.9 = 2335.1. + 150 = 2485.1, meta **2485**. IMC 21.5, normal.
- **h30**: hombre que cumple 30 hoy, 70000 g, 1700 mm, `moderate`, `maintain`, ajustes 300/300. 1617.5 × 1.55 = 2507.125, meta **2507** (no ajusta).
- **conf**: igual que h45 pero sin ajuste propio y con meta escrita 1900. Rige 1900 (`confirmed`). La sugerida es 2497.6875 → 2498 con el ajuste de la casa en 0. **Me equivoqué yo:** primero esperaba 2098, como si tuviera el 400 de h45; el motor tenía razón.
- **sinact**: le falta la actividad, así que no hay sugerencia y `missing_inputs` = [activity].
- **Cortes con 1700 mm**: 53176 g da 18.4 (bajo peso), 53465 g da 18.5 (normal), 72250 g da 25.0 (sobrepeso), 86700 g da 30.0 (obesidad).

Resultado: 9 de 9 coinciden.

## 5 valores USDA (S4-10, (e))
La API `api.nal.usda.gov` con DEMO_KEY respondió `OVER_RATE_LIMIT`. Coteje contra los CSV oficiales locales (los zip de FDC; su SHA-256 coincide con `PROCEDENCIA.md`), leídos con mi propio script y no con `construir_tabla.py`. Elegí 5 renglones al azar, distintos de los 6 que ya revisó el agente:

| fdc_id | Tipo | Alimento | kcal | P / G / C |
|---|---|---|---|---|
| 172438 | SR | Refried beans | 90 | 4.98 / 2.01 / 13.55 |
| 170457 | SR | Tomatoes, red, ripe, raw | 18 | 0.88 / 0.2 / 3.89 |
| 169231 | SR | Ginger root, raw | 80 | 1.82 / 0.75 / 17.77 |
| 327198 | Foundation | Melons, cantaloupe, raw | 34.0 | 0.82 / 0.18 / 8.16 |
| 171893 | SR | Coffee, instant, powder | 353 | 12.2 / 0.5 / 75.4 |

Resultado: los 5 coinciden con `tabla.csv`.

Empates por revisar: `revisar.tsv` tiene 50 renglones, de los que 40 son de confianza media o baja y 10 son productos sin fuente (no alimentos). Ninguno tiene valor inventado y los 150 renglones de la tabla siguen en `propuesto`.

## Verificación (comandos y salida resumida)
- **Calidad Go**:
  - `go vet ./internal/recetas/... ./cmd/motores/`: limpio.
  - `gofmt -l internal/recetas cmd/motores`: vacío.
  - `go test -count=1 -coverprofile ./internal/recetas/...`: ok, **93.3 %**.
  - Por cada `-run`, todos PASS: Escalado 3, Lista 7, Inventario 3, Plan 5, Costo 5, Salud 4, Kcal 5, Canal 3, Recetario 3, RecetaEditable 1, Restricciones 1, SinIA 1.
- **Motor contra oráculo**:
  - `go build -o datos/tmp/rev-s4/motores ./cmd/motores`, y `s4 dump` dos veces: código 0, salidas idénticas con `cmp`.
  - `python3 -I oraculo/rec_run_fixtures.py`: 25 de 25, 0 diferencias.
  - `python3 -I oraculo/compara.py s4 --motor datos/tmp/rev-s4/motor-s4.json`: 1218 de 1218 en las tres comparaciones, 210 de 210 provisionales, 0 diferencias, código 0.
- **Mis fixtures**:
  - `s4-revisor-recetario` (361 consultas): 0 diferencias entre el motor y mis cálculos.
  - `s4-revisor-salud`: 9 casos, 0 diferencias después de corregir mi esperado de `conf`.
- **Recetario** (`valida.py`): 60 entradas y 53 ids distintos. Errores de §3.2: 30, todos por `etiquetas`, más `etiquetas_por_receta` en la raíz de b. 0 productos o presentaciones inexistentes, 0 unidades distintas de la base, 0 porciones distintas de 2, 0 «al gusto», 0 ingredientes sin canal. Nombres iguales a un nombre de USDA: 0 en recetas y 1 en productos (`catsup`).
- **Construcciones prohibidas**: `grep float32|float64`, `http`, `os.Getenv|database/sql|modernc|net`, todo fuera de `_test.go`: 0.
- **Independencia**: `grep backend|internal/recetas|.go` en `oraculo/rec_*.py`: 0.
- **S4-13**: base temporal con `sqlite3` y las secciones Up de 00001 a 00003. `pragma_table_info` da 1 y los renglones con restricción dan 0. En `app_setting`, `health.goal_adjust_*` = 0 y `ai.*` = false.

## Supuestos
- Conté las 60 entradas de a+b como «el recetario» para S4-02, S4-08 y S4-09, porque `recetario.json` no existe. A las repetidas les puse el sufijo `-b` solo dentro de mi fixture temporal.
- Las lentejas las marqué como «probable repetido» por mi criterio: los nombres son casi iguales y comparten 10 de 11 ingredientes. La decisión es del autor o del integrador.

## Bloqueos
Ninguno nuevo.

## Pendientes
Corregir H1, H2, H3 y M1 a M3 (rutas en `REVISION-S4-r1.md`). La ronda 2 los vuelve a revisar.

## Solicitudes fuera de alcance
Las correcciones de la tabla de hallazgos. El revisor no corrige.

## Paquetes instalados
Ninguno. `goose` no está instalado y usé `sqlite3`, que ya estaba en el sistema.

## Propuestas para Axel
Ninguna sobre K31. Para la puerta 2:
- confirmar los 150 renglones `propuesto` de `tabla.csv` (`S4-confirmacion-empates`);
- aceptar SR Legacy como parte de «USDA FoodData Central» (B2).

## Consumo aproximado
- **Acciones de herramienta**: unas 37.
- **Leído**: unos 215 KB. Lo principal: protocolo 6, criterios 18, decisiones 16, contrato unos 50 KB por partes, reportes 34, pruebas y código unos 16, y salidas de verificación.
- **Escrito**: unos 40 KB, entre scripts y fixtures temporales (las fixtures las generó un script, unos 250 KB que no pasan por herramienta) y estos dos reportes (unos 22 KB).
