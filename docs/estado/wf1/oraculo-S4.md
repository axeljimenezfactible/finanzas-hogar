# WF-1 · oraculo-S4 · reporte final (11-oct-2026, reanudación 19:00)

Agente `oraculo-S4` (Sonnet). Contrato leído: `docs/contratos/motores/S4-interfaz.md` 1.4 (ronda 4), `oraculo/fixtures/FORMATO.md` 1.2, criterios `docs/criterios/S4.md` 1.2, decisiones K7, K23, K36, K38, K43, K44. No se leyó `backend/internal/recetas`.

## Hecho
- Oráculo independiente de S4 en Python 3.14 (`python3 -I`, solo biblioteca estándar: `json`, `csv`, `fractions`, `decimal`, `datetime`). Parciales del intento anterior revisados y conservados lo que cumplía; reescrito `shopping_list` a la versión 1.4 del contrato (cantidad fijada, artículos `manual`, precio y recomendación por artículo, `quantity_milli`, `estimated_total_cents`, `visit_date`, orden de `groups[]`).
- Funciones del §12.2: `scale_recipe`, `default_week`, `unplanned_slots`, `week_plan`, `shopping_list`, `inventory_status`, `recipe_detail`, `nutrition_coverage`, `health_profile`; validación estricta de `Input` (`unknown_field`, `missing_field`, `invalid_date`, `invalid_argument`) y errores de §14.
- 25 fixtures en `oraculo/fixtures/s4/` (generadas por `oraculo/rec_genera_fixtures.py`, determinista) y `oraculo/salidas/s4.json` (`salida/1`, 64 KB, dos corridas idénticas byte a byte).
- Recetas reales usadas (cita): `huevos-a-la-mexicana`, `huevo-con-jamon-y-frijoles`, `avena-con-manzana-y-canela`, `yogur-griego-con-fresas-y-avena` (de `datos/semilla/recetario-a.json`) y `espagueti-con-atun-y-elote` (de `recetario-b.json`) en `s4-kcal-recetas-autor`; además `huevo-con-jamon-y-frijoles` y `yogur-griego-con-fresas-y-avena` en `s4-kcal-sin-dato` y `espagueti-con-atun-y-elote` en `s4-costo-recetas`. No sustituyen S4-09 (las elige el revisor).

## Archivos
- `oraculo/rec_comun.py` (errores, HALF_UP, fechas, validación), `rec_datos.py` (Input, catálogo/recetario/CSV), `rec_calculo.py` (escalado, cotización que rige, costo, kcal, inventario, salud), `rec_plan.py` (plan semanal, sobrantes, lista), `rec_run_fixtures.py` (ejecutor), `rec_genera_fixtures.py` (autor único de las fixtures).
- `oraculo/fixtures/s4/*.json` (25), `oraculo/salidas/s4.json`, este reporte.

Fixtures: `s4-escalado-2a4`, `s4-lista-sin-existencia`, `-existencia-parcial`, `-existencia-mayor`, `-sobrantes`, `-cantidad-manual` (ejemplo de §6.6: 19380), `-reposicion-ritmo` (PROVISIONAL), `s4-inventario-modos`, `s4-plan-celdas`, `s4-plan-karen-rechazada`, `s4-plan-fuera-con-receta`, `s4-costo-cotizacion-elegida`, `s4-costo-dos-recomendadas`, `s4-canal-mercados` (PROVISIONAL), `s4-canal-sin-grupo`, `s4-costo-recetas` (§8.3: 4425 / 2213 / 8850), `s4-kcal-recetas-autor`, `s4-kcal-sin-dato`, `s4-salud-perfiles`, `s4-salud-ajuste-inicial`, `s4-salud-imc-franja` (PROVISIONAL), `s4-entrada-campo-ausente`, `-campo-desconocido`, `-fecha-invalida`, `-producto-desconocido`.

## Verificación
- `python3 -I oraculo/rec_genera_fixtures.py` -> `25 fixtures`; `python3 -I oraculo/rec_run_fixtures.py` -> `25 fixtures, 25 coinciden, 1518 hojas de esperado, 0 diferencias`. Los `esperado` de números sintéticos están escritos a mano (aritmética en la descripción y en comentarios del generador) y coincidieron a la primera con el oráculo; los de recetas reales salen de una segunda vía independiente (solo `fractions`, otro código que `rec_calculo.py`).
- Prueba negativa: con 3 valores de `esperado` alterados (9566, 2763, 2555) el ejecutor reporta exactamente 3 `DIF`.
- `python3 -I oraculo/compara.py s4 --sin-motor` -> `25 fixtures`, `oráculo contra esperado: 1218 de 1218`, `provisionales: 210 de 210`, `model.json (ramas de s4): 0 de 0`, `diferencias: 0`, código 0.
- Dos corridas de generador+ejecutor: `cmp` idéntico.
- Salud, dos perfiles ilustrativos a mano (hoy 2026-10-11): hombre 30 años, 80 kg, 175 cm, sedentary: MSJ = 10x80 + 6.25x175 - 5x30 + 5 = 1,748.75; x1.2 = 2,098.50 -> **2,099** (HALF_UP); con «bajar» y 300 de la casa (fixture): 1,798.50 -> **1,799**. Mujer 28 años, 60 kg, 165 cm: MSJ = 600 + 1,031.25 - 140 - 161 = 1,330.25; x1.2 = 1,596.30 -> **1,596**; con «subir» y 250 (fixture): 1,846.30 -> **1,846**. Además RN-17 (35 años, 72.5 kg, 175 cm, moderate): 1,648.75 x1.55 = 2,555.5625 -> 2,556 y con bajar 500 propios 2,056; IMC 72,500 g y 1,750 mm = 236.73 -> 237 normal; cortes 18.5 / 25.0 / 30.0 con 2,000 mm -> normal / overweight / obesity. Con ajustes iniciales 0 y 0 (`s4-salud-ajuste-inicial`): 2,099 y 1,596.
- Receta real a mano: `espagueti-con-atun-y-elote` a 2 porciones, kcal ILUSTRATIVAS: 927.5 + 162.4 + 89.1 + 54 + 24 + 5.96 + 121.06 + 0 + 2.65 = 1,386.67 -> 1,387 (693 por porción); coincide con la fixture.
- Costo §8.3: 1200 + 2100 + 675 + 450 = 4425; por porción HALF_UP(2212.5) = 2213; 4 porciones 8850. Ejemplo §6.6: 19380.
- Criterios cubiertos con fixtures: S4-02 (`s4-escalado-2a4`), S4-03 (`s4-lista-*`), S4-04 (`s4-inventario-modos`), S4-05 (`s4-plan-*`), S4-06 (`s4-costo-*`), S4-07 (`s4-salud-*`), S4-08 (`s4-kcal-*`), S4-12 (`s4-canal-*`). S4-01, S4-09, S4-10, S4-11, S4-13, S4-14 y S4-15 no son del oráculo.

## Supuestos (ninguno cambia cifras canon; el contrato no los fija)
- S-O1: «ausente» y `null` son lo mismo en la salida (se escribe `null`).
- S-O2: orden de listas que el contrato no fija. `cells[]` por fecha, tiempo (desayuno, comida, cena, colación), orden de `persons` y orden de entrada; `days[].persons[]` solo personas con plan, en el orden de `persons`; `items[]` primero los `calculated` por código de producto y luego los `manual` por código; `recipes[]`, `missing_prices[]` y `missing[]` ordenados por código; `reasons[]` en el orden `recipe`, `replenishment`. Las fixtures escriben listas completas, así que el motor Go debe coincidir o el integrador ajusta el contrato.
- S-O3: una celda con algún ingrediente sin precio no declara `price_store` ni `price_observed_date` (ninguna fixture lo afirma).
- S-O4: `missing_inputs` usa `sex`, `birth_date`, `height_mm`, `weight`, `activity`, `goal`; el contrato no da los nombres y ninguna fixture los afirma.
- S-O5: `recipe_detail.ingredients[].kcal` = HALF_UP de cada ingrediente por separado; las fixtures de kcal no lo afirman (solo `quantity`, totales y por porción).
- S-O6: una fila de `list_choices` de un `calculated` cuenta como «elección de persona» (§6.2 paso 7): el artículo sigue aunque su necesidad sea 0.
- S-O7: los cinco últimos campos de `inventory[]` y los de `health_profiles[]` pueden faltar o ser `null` en el oráculo; las fixtures los escriben siempre, explícitos.
- S-O8: un sobrante caducado no cubre entradas pero sí aparece en `leftovers[]` (opción A de `S4-reglas-lista`, pregunta 1).
- S-O9: una tienda de cotización que no está en `stores` se ordena después de las declaradas; las fixtures declaran todas.

## Bloqueos
Ninguno nuevo. Las fixtures provisionales citan `docs/estado/bloqueos/S4-reglas-lista.md` (`s4-lista-reposicion-ritmo`, `s4-salud-imc-franja`) y `S4-asignacion-mercado.md` (`s4-canal-mercados`).

## Hallazgos para el integrador
- `datos/semilla/recetario-a.json` y `recetario-b.json` repiten 7 ids (`bistec-a-la-mexicana`, `caldo-de-pollo-con-verduras`, `frijoles-de-la-olla`, `garbanzos-con-espinaca`, `huevos-rancheros`, `picadillo-de-res-con-papa`, `pollo-en-salsa-verde`): juntos suman 53 ids distintos, no 60. El `recetario.json` de 60 recetas únicas (S4-01) exige renombrar o reemplazar 7.
- `datos/semilla/nutricion/tabla.csv` aún no existe (el agente de nutrición no terminó). `s4-kcal-recetas-autor` y `s4-kcal-sin-dato` llevan kcal ILUSTRATIVAS en línea (`nutrition`, fuente `manual`, redondeadas, NO son datos de fuente); hay que regenerarlas con `nutrition_file` cuando exista la tabla.

## Pendientes (rutas exactas)
- `oraculo/fixtures/s4/s4-escalado-recetario.json` (60 recetas x 4 escalas): `python3 -I oraculo/rec_genera_fixtures.py` lo genera solo cuando exista `datos/semilla/recetario.json`; hoy imprime `PENDIENTE`.
- Comparar contra el motor Go: `python3 -I oraculo/compara.py s4` cuando exista `datos/tmp/motor-s4.json`.
- Sustituir las kcal ILUSTRATIVAS por `datos/semilla/nutricion/tabla.csv` (arriba).

## Solicitudes fuera de alcance
- `docs/contratos/motores/S4-interfaz.md` §12.3 y §12 nombran `oraculo/s4_run_fixtures.py`; el ejecutor entregado es `oraculo/rec_run_fixtures.py` (nombre de la tarea). Cambiar el contrato o crear un alias.
- `oraculo/README.md`: agregar la tabla de `rec_*.py` y los comandos `python3 -I oraculo/rec_genera_fixtures.py`, `python3 -I oraculo/rec_run_fixtures.py`, `python3 -I oraculo/compara.py s4`.
- `docs/contratos/motores/S4-interfaz.md` §5.3/§6.5/§10.2: fijar el orden de `cells[]`, `items[]` y los nombres de `missing_inputs` (S-O2, S-O4).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: no se tocaron valores K31 (todo lo operativo de S4 nace en 0 y 0 por K36/K44).

## Consumo aproximado
26 acciones de herramienta; ~215 KB leídos; ~70 KB escritos con herramientas (generador, ajustes a rec_plan y rec_datos, este reporte); las 25 fixtures (275 KB) y salidas/s4.json (64 KB) las genera el script.
