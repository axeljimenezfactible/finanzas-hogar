# motor-S4-go · Motor de recetas, comidas y salud (WF-1, 11-oct-2026, reanudación de las 19:00)

Contrato: `docs/contratos/motores/S4-interfaz.md` 1.4 (ronda 4). Criterios: `docs/criterios/S4.md` 1.2.

## Hecho
- Revisé los parciales del intento anterior (`errors.go`, `num.go`, `types.go`, `catalog.go`): compilan, cumplen el contrato y los conservé (solo `gofmt`).
- Agregué `recipebook.go` (cargador del recetario en los modos `review` y `production`, validación de §3.2, cotejo de kcal, nombres en inglés y las funciones puras de §3.3), `nutrition.go` (`tabla.csv`, cobertura y kcal), `cost.go` (escalado, cotización que rige y `recipe_detail`), `plan.go` (celdas, «fuera», huecos y sobrantes registrados), `shopping.go` (lista neteada, cantidad fijada, artículos `manual`, mercados, grupos y totales), `inventory.go` (estado y movimientos), `health.go` (edad, IMC, Mifflin-St Jeor × actividad ± ajuste) y `eval.go` (entrada estricta de §12.1 y consultas de §12.2).
- `backend/cmd/motores/s4.go` registra `s4` en el binario común; `s4 dump` corre las 25 fixtures de `oraculo/fixtures/s4/`.
- 41 pruebas con los nombres de §13, más datos de prueba ficticios en `testdata/` (5 recetas propias). `testdata/generar_testdata.py` calcula aparte las kcal y la tabla de escalado con `fractions`.

## Archivos
`backend/internal/recetas/{errors,num,types,catalog,recipebook,nutrition,cost,plan,shopping,inventory,health,eval}.go`, `recetas_test.go`, `testdata/{generar_testdata.py,catalogo.json,recetario.json,tabla.csv,escalado.json}`; `backend/cmd/motores/s4.go`; este reporte.

## Verificación
- `gofmt -l internal/recetas cmd/motores/s4.go`: vacío. `go vet ./internal/recetas/... ./cmd/motores/`: limpio. `go test ./cmd/motores/`: ok.
- `go test ./internal/recetas/ -count=1 -cover -v`: 41 de 41 PASS, **cobertura 93.3 %** (`go tool cover -func`: total 93.3 %).
- `-run` por criterio: Escalado 3, Lista 7, Inventario 3, Plan 5, Costo 5, Salud 4, Kcal 5, Canal 3, todas PASS.
- `datos/tmp/motores s4 dump --fixtures oraculo/fixtures/s4 --out datos/tmp/motor-s4.json`: código 0; dos corridas dan bytes iguales. `python3 -I oraculo/compara.py s4`: oráculo contra esperado 1218 de 1218, **motor contra esperado 1218 de 1218, oráculo contra motor 1218 de 1218**, provisionales 210 de 210, **0 diferencias**. La primera corrida dio 6 diferencias en el orden de `items[].recipes[]`; las corregí ordenando por código (S-9).
- S4-01: testdata 5 de 5 válidas. **Recetario real: no hay `recetario.json`, solo las mitades a y b. Juntas dan 52 válidas y 8 rechazadas por repetidas** (en modo `production`). En modo `review` se rechazan las 60 por `kcal_mismatch`: `kcal_por_porcion` está en null y `fuentes_kcal` vacío en las 60 recetas, y `tabla.csv` ya existe. **S4-01 NO se cumple con los datos actuales** (solicitudes 1 y 2). Mientras no exista el archivo oficial, la prueba informa sin fallar; cuando exista, falla si algo no valida.
- S4-02: 52 recetas reales × sus ingredientes × 4 escalas = 1940 escalados contra multiplicación cruzada exacta, con 0 diferencias. testdata contra `escalado.json` de Python, 0 diferencias. `grep float32|float64` fuera de pruebas: 0.
- S4-03 y S4-06: §6.6 exacto (arroz 3, 2500000, 8400; frijol 1, 0, 3000; aceite `manual` 2, 1892000, 7980; servilletas null; grupo walmart 19380; raíz 19380). §8.3: 4425, 2213 y 8850.
- S4-04: 2838000 − HALF_UP(946000 × 30 / 21) = 1486571; reponer cuando la existencia es igual al nivel; agotado en 0 y prevalece.
- S4-07: 2098.50 → 2099; 1596.30 → 1596; bajar 300 → 1798.50 → 1799; subir 250 → 1846.30 → 1846; «mantener» con 300 y 300 → 2099; RN-17 2556 y 2056; IMC 18.5 `normal`, 25.0 `overweight`, 30.0 `obesity`, 23.7 y la franja de 73990 g (opción A).
- S4-08: cobertura de testdata 13 únicos, 12 con fuente (11 confirmados, 1 propuesto), sin dato [frijol]. **Real (modo `review`): 89 únicos, 89 con fuente, 0 confirmados, 89 propuestos, 0 sin dato.**
- S4-10 y S4-14: `grep http` fuera de pruebas = 0; `TestSinIA` no encuentra lecturas de entorno, red, `database/sql` ni flotantes.
- S4-13: `goose` no está instalado. Apliqué las secciones Up de 00001 a 00003 con `sqlite3` en `datos/tmp/s4-13.db` (sin semilla de S1, porque no existe). El campo existe = 1 y no está vacío en ningún renglón = 0. La base temporal la borré.

## Supuestos (ninguno cambia cifras canon: S4 no tiene ramas en `model.json`)
- S-6: un artículo calculado sin cotización toma la presentación del primer ingrediente que lo aporta (orden fecha, tiempo, persona). Si solo viene de reposición o de una elección, toma `presentacion_tipica` o la primera del catálogo.
- S-7: `nutrition_coverage.missing[]` va en orden de código.
- S-8: una celda «fuera» lleva `cost_cents` y `kcal` null y no entra en las sumas del día, así que el día no queda null por ella.
- S-9: `items[].recipes[]` va por código; `items[]` primero los `calculated` por código de producto y luego los `manual` en el orden de `list_choices`; `missing_prices[]` en el orden de `items[]`, sin repetir.
- S-10: las cotizaciones (exactamente una recomendada, ninguna repetida) se validan al usarlas: solo fallan las consultas que tocan ese producto.
- S-11: el cargador ignora las claves que no están en §3.2 (`etiquetas`, `etiquetas_por_receta`).
- S-12: la cuenta de 60 la exige `SeedComplete()`, no el cargador, para que testdata y las fixtures usen recetarios chicos.
- S-13: `recipe_detail.ingredients[]` excluye los opcionales (S-3); `scale_recipe` los incluye.
- S-14: `missing_inputs` usa `sex, birth_date, height_mm, weight, activity, goal`, en ese orden.
- S-15: `health_profile.daily_kcal_target` es la meta que rige, igual que en `week_plan`.
- S-16: si `recipe_book_file` trae alguna receta rechazada, falla toda la entrada.
- S-17: con un conteo posterior a la fecha de cálculo, los días transcurridos cuentan como 0.

## Bloqueos
Ninguno nuevo. PROVISIONAL en el código, con su cita: `S4-reglas-lista` (preguntas 1, 2 y 3), `S4-asignacion-mercado` y `S4-confirmacion-empates`.

## Pendientes (de otros alcances)
- `datos/semilla/recetario.json` (unión de a y b), `datos/semilla/recetario.schema.json` y `oraculo/fixtures/s4/s4-escalado-recetario.json` no existen.
- Confirmar los 89 renglones `propuesto` de `tabla.csv` (Axel en la puerta 2). Mientras tanto, en `production` las recetas salen «incompletas».

## Solicitudes fuera de alcance
1. `datos/semilla/recetario-b.json` (agente del recetario o integrador): repite por id 7 recetas de a (`bistec-a-la-mexicana`, `caldo-de-pollo-con-verduras`, `frijoles-de-la-olla`, `garbanzos-con-espinaca`, `huevos-rancheros`, `picadillo-de-res-con-papa`, `pollo-en-salsa-verde`) y por nombre 1 (`tinga-de-pollo-en-tostadas` = «tinga-de-pollo» de a). Hay que sustituir 8 recetas para llegar a 60 únicas (S4-01).
2. Mismas rutas: regenerar `kcal_por_porcion` y `fuentes_kcal[]` con `tabla.csv` en modo `review` (§3.2). Sin eso, `TestRecetarioValido` y toda fixture con `recipe_book_file` dan `kcal_mismatch`.
3. Contrato S4 §3.1 o esquema: declarar o quitar `etiquetas` y `etiquetas_por_receta`.
4. Integrador: `roundHalfUp` de `num.go` se puede unificar con el de `internal/dominio`, si existe. `goose` no está instalado (S6-03 lo usa).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: ninguno de los ocho valores de K31 es de S4.

## Consumo aproximado
33 acciones de herramienta. Leído: ~190 KB. Escrito: ~145 KB con herramientas (código 90 KB, pruebas 47 KB, generador 7 KB y este reporte), más ~22 KB que generó el script en `testdata/`. **Desviación: `recetas_test.go` (47 KB) se escribió en una sola escritura, por encima de los 40 KB del protocolo §4.**
