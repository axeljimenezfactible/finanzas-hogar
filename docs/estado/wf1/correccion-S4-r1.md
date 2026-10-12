# correccion-S4-r1 · Reporte final (11-oct-2026)
Agente `correccion-S4-r1`. Corrige los hallazgos altos y medios de `docs/estado/wf1/REVISION-S4-r1.md` (H1, H2, H3, M1, M2, M3). No se reabrió ninguna decisión de §A–§K y no se tocó ninguna cifra de `oraculo/model.json`.

## Hecho
- **H1 (alta)**: en `recetario-b.json` cambié las 9 recetas que chocaban con `recetario-a.json` (se quedan las versiones de a): los 7 ids repetidos (`pollo-en-salsa-verde`, `caldo-de-pollo-con-verduras`, `picadillo-de-res-con-papa`, `bistec-a-la-mexicana`, `huevos-rancheros`, `frijoles-de-la-olla`, `garbanzos-con-espinaca`), `tinga-de-pollo-en-tostadas` (mismo nombre que `a:tinga-de-pollo`) y `lentejas-con-verduras` (caso de las lentejas: **se cambia**; tenía 10 de 11 ingredientes en común con `a:lentejas-guisadas`). Las recetas nuevas son platillos distintos: `pozole-rojo-de-cerdo`, `tacos-dorados-de-papa`, `chiles-rellenos-de-queso`, `costillas-de-cerdo-con-verdolagas`, `mole-de-olla-de-res`, `entomatadas-de-pollo`, `rajas-con-crema-y-elote`, `pechuga-a-la-plancha-con-brocoli` y `calabaza-en-tacha` (colación). Su mayor parecido por ingredientes (Jaccard) con cualquier otra receta es 0.67; el de las lentejas era 0.83 (10 de 12). Solo usan productos con renglón en `tabla.csv` y con `gramos_por_unidad` cuando van en ml o pieza.
- Generador nuevo `oraculo/rec_genera_recetario.py`: une a y b en `datos/semilla/recetario.json` con la llave `recetas`. Valida las reglas cruzadas de §3.2 (id y nombre únicos sin distinguir mayúsculas, producto una vez por receta, producto y presentación existentes, unidad base, nada de «al gusto», nombre distinto de todo `nombre_fuente`) y exige 60. Es determinista y tiene `--verifica`. Resultado: 60 ids y 60 nombres únicos.
- **H2 (alta)**: `datos/semilla/recetario.schema.json` (JSON Schema 2020-12) con exactamente los campos de §3.1 y §3.2, `additionalProperties: false` en todos los niveles y, para `origen = claude`, `porciones_base = 2` y `opcional = false`.
- **H3 (alta)**: el generador escribe `kcal_por_porcion` y `fuentes_kcal` desde `tabla.csv` en modo `review`, con las reglas de §9.2 y §9.3: fracciones exactas y un solo HALF_UP. Ninguna de las 60 queda en null. Los 5 valores de referencia del revisor coinciden: chuletas-de-cerdo-en-guajillo 563, tilapia-al-mojo-de-ajo 459, albondigas-en-chipotle 811, licuado-de-platano-y-avena 348 y espagueti-con-atun-y-elote 666.
- **M1 (media)**: §3.2 dice «exactamente los campos de abajo», así que `etiquetas` **no entra** en el recetario oficial ni en el esquema. Los borradores a y b conservan las etiquetas con el mismo formato (`etiquetas` dentro de cada receta). La llave de raíz `etiquetas_por_receta` de b se pasó a cada receta y queda solo como material de autoría. El generador las quita al unir. El cargador Go **ya no las ignora**: `ParseRecipeBook` usa `DisallowUnknownFields` y una llave fuera de §3.1–§3.2 es `invalid_argument`, igual que el esquema (el supuesto S-11 de `motor-S4-go.md` queda sin efecto). Agregué la prueba dentro de `TestRecetarioValido`. Pedí el cambio del contrato aparte, en Solicitudes.
- **M2 (media)**: se generó `oraculo/fixtures/s4/s4-escalado-recetario.json`: 60 recetas × 4 escalas (2→4, 2→1, 2→3, 2→6), 240 consultas, calculadas con `fractions`. Lleva `nutrition_file` porque el motor carga `recipe_book_file` en modo `review` y coteja las kcal.
- **M3 (media)**: `s4-kcal-recetas-autor` usa `nutrition_file: datos/semilla/nutricion/tabla.csv` en modo `review`. La `descripcion` cita el valor, la fuente y el `id_fuente` de cada uno de los 22 productos. Cobertura: 22 con fuente, 0 confirmados y 22 propuestos. Ya no hay `confirmed` inventado. `s4-kcal-sin-dato` lleva los renglones en línea **copiados tal cual de tabla.csv** (valor, fuente, `id_fuente` y estado `proposed`) y omite a propósito `jamon-pavo` para simular el «sin dato». No puede usar `nutrition_file` por dos razones: tabla.csv cubre los 150 productos de comida, y `nutrition_file` junto con `nutrition` es `invalid_argument` (§12.1). Como las dos dependen de renglones `propuesto`, llevan `provisional` → `docs/estado/bloqueos/S4-confirmacion-empates.md`. Las recetas reales de las fixtures se leen ahora de `recetario.json`.
- B1 (baja) quedó resuelto de paso: sin `recetario.json` la prueba falla en lugar de pasar en verde. Además, la unión a+b con etiquetas ya no se podría leer con el cargador estricto.

## Archivos
- `/home/deepsec/finanzas-hogar/datos/semilla/recetario.json` (nuevo, generado, 201 KB)
- `/home/deepsec/finanzas-hogar/datos/semilla/recetario.schema.json` (nuevo)
- `/home/deepsec/finanzas-hogar/datos/semilla/recetario-b.json` (9 recetas nuevas y etiquetas por receta)
- `/home/deepsec/finanzas-hogar/oraculo/rec_genera_recetario.py` (nuevo)
- `/home/deepsec/finanzas-hogar/oraculo/rec_genera_fixtures.py` (kcal desde tabla.csv; recetas desde recetario.json; nutrition_file en el escalado)
- `/home/deepsec/finanzas-hogar/oraculo/fixtures/s4/s4-escalado-recetario.json` (nuevo), `s4-kcal-recetas-autor.json` y `s4-kcal-sin-dato.json` (regenerados); las otras 23 se regeneraron sin cambio de cifras
- `/home/deepsec/finanzas-hogar/oraculo/salidas/s4.json` (salida del oráculo)
- `/home/deepsec/finanzas-hogar/backend/internal/recetas/recipebook.go` (cargador estricto), `recetas_test.go` (`realBook` sin la unión provisional; prueba de campo desconocido)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-S4-r1.md` (este reporte)
- Script de apoyo fuera del repositorio: el scratchpad de la sesión, `s4/nuevas_b.py` (las 9 recetas y el cambio de formato de b)

## Verificación
- **S4-01**: `python3 -I -c "...len(d['recetas'])"` = **60**. `uv run --no-project --with jsonschema python -I -m jsonschema -i datos/semilla/recetario.json datos/semilla/recetario.schema.json` termina con 0 y sin errores (solo el aviso de que el CLI está en desuso). Con `Draft202012Validator`: esquema válido, 0 errores y **60 de 60** recetas válidas. Los borradores con etiquetas dan errores, como se espera. `go test -run Recetario -v`: «review 60 válidas y 0 rechazadas; production 60 válidas y 0 rechazadas». `python3 -I oraculo/rec_genera_recetario.py --verifica`: coincide.
- **S4-02**: `TestEscaladoRecetario` cubre 60 recetas y 2244 escalados, 0 diferencias. La tabla independiente `s4-escalado-recetario` da 0 diferencias en `compara`. `grep float32|float64` fuera de pruebas: 0.
- **S4-08**: 100 ingredientes únicos, 100 con renglón y fuente, 0 sin dato, 0 con fuente vacía, 0 confirmados (PROVISIONAL, `S4-confirmacion-empates`). La prueba Go da la misma cobertura.
- **S4-09** (referencia del revisor): 563 / 459 / 811 / 348 / 666, iguales en `recetario.json` y en el motor.
- **S4-11**: 0 nombres visibles (recetas y productos usados) iguales a un `nombre_fuente`. **S4-12**: 0 ingredientes sin canal. **S4-10**: `grep http` fuera de pruebas da 0.
- **S4-15**: `go vet` limpio; `gofmt -l internal/recetas cmd/motores` vacío; `go test -coverprofile`: **93.2 %** (umbral 85 %, K25); 43 PASS, 0 FAIL, 0 SKIP; `go test ./cmd/motores/...` ok.
- **Oráculo**: `python3 -I oraculo/rec_genera_fixtures.py` genera 26 fixtures. `rec_run_fixtures.py`: 26 de 26 coinciden, 0 diferencias. `motores s4 dump` sale con 0. `python3 -I oraculo/compara.py s4`: oráculo contra esperado, motor contra esperado y oráculo contra motor, **8451 de 8451** hojas cada uno; provisionales 429 de 429; **0 diferencias**. Corrí el dump y el oráculo dos veces cada uno y las salidas son idénticas byte a byte.

## Supuestos
1. Alcance: `datos/semilla/recetario.json` y `recetario.schema.json` no estaban en la lista de rutas, pero H1 y H2 los exigen (y el contrato §3 los asigna al agente del recetario).
2. Las cantidades de las 9 recetas nuevas son estimaciones de cocina casera, con los mismos criterios de a y b (gramos exactos, presentación = `presentacion_tipica`). No cambian ninguna cifra canónica.
3. Que las etiquetas queden fuera del recetario oficial es la lectura literal de §3.2. No se inventa un campo nuevo.
4. `jsonschema` corrió de forma efímera con `uv` (Python 3.13 del entorno de uv); los scripts propios corren con `python3 -I` (3.14), solo biblioteca estándar.

## Bloqueos
Ninguno nuevo. Las fixtures de kcal llevan `provisional` citando `docs/estado/bloqueos/S4-confirmacion-empates.md` (existente).

## Pendientes
Ninguno de los hallazgos altos ni medios. No tocados por estar fuera de la lista: B2, B3, B4 y B6 (bajas). Nota para el revisor: entre recetas que ya estaban hay pares muy parecidos por ingredientes que la ronda 1 no señaló: huevos-a-la-mexicana con huevos-con-nopales (0.88), chilaquiles-verdes-con-huevo con chilaquiles-verdes-con-pollo (0.83) y pollo-en-salsa-verde con chilaquiles-verdes-con-pollo (0.82). Sus nombres son distintos y no se tocaron.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S4-interfaz.md` §3 y §3.2: agregar que `recetario.json` lo genera `oraculo/rec_genera_recetario.py` a partir de los borradores `recetario-a.json` y `recetario-b.json`; que las etiquetas de los borradores no forman parte del recetario 1.0; y que el cargador y el esquema rechazan toda llave fuera de §3.1–§3.2. Motivo: M1 pide dejarlo decidido en el contrato.
2. `docs/estado/wf1/motor-S4-go.md`: el supuesto S-11 («el cargador ignora las claves que no están en §3.2») queda sin efecto. Motivo: el cargador ahora las rechaza.
3. `oraculo/README.md`: agregar `python3 -I oraculo/rec_genera_recetario.py [--verifica]` antes de `rec_genera_fixtures.py`. Motivo: orden de regeneración (relacionado con B4).

## Paquetes instalados
Ninguno del sistema (`jsonschema` efímero con `uv run --with`).

## Propuestas para Axel
- Puerta 2: ¿las recetas llevan etiquetas («económica», «rápida», «sobrantes») para filtrar en el plan? (a) No, como hoy: el recetario 1.0 queda sin cambios. (b) Sí, como campo opcional: hace falta la versión 1.5 del contrato S4, el esquema, el cargador y una columna o tabla nueva en la base (el esquema 1.3 no la tiene). Los borradores ya las traen para las 60 recetas. No hay valores K31 en este alcance.

## Consumo aproximado
- Acciones de herramienta: ≈ 38.
- Leído: ≈ 160 KB (protocolo, decisiones §D y §K, criterios S4, contrato S4 §3, §4, §9 y §12, revisión, reportes del recetario, partes del código Go y del generador, salidas de verificación).
- Escrito con herramientas: ≈ 50 KB (scripts, esquema, ediciones Go y del generador, reportes). Generado por scripts, sin pasar por el contexto: `recetario.json` 201 KB, `recetario-b.json` 105 KB, fixtures s4 ≈ 490 KB (`s4-escalado-recetario` 452 KB) y salidas del oráculo.
