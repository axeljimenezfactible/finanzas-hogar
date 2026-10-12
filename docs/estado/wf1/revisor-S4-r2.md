# revisor-S4-r2 · Reporte de agente (11-oct-2026)
Revisor independiente de S4, ronda 2 de 3. No escribí nada de lo que reviso y no corregí nada. Informe de revisión: `docs/estado/wf1/REVISION-S4-r2.md`. Veredicto: **APROBADO** (0 altos, 0 medios, 7 bajos).

## Hecho
- Leí, por partes, el protocolo, las decisiones §D y §K, los criterios `S4.md` completos, el contrato S4 §1–§4 y §9–§13, el reporte `correccion-S4-r1.md` y los hallazgos bajos de `REVISION-S4-r1.md`.
- Reejecuté S4-01 a S4-15 con mis propios comandos y con 5 scripts en `datos/tmp/rev-s4-r2/`: validación cruzada y recálculo de kcal de las 60 recetas, tabla de escalado, listas neteadas, 5 recetas a mano y cotejo USDA desde los zip. Los scripts ya están borrados.
- Comprobé los 6 hallazgos altos y medios de la ronda 1: los 6 están cerrados. Comprobé también B1 (cerrado) y B2, B3, B4 y B6 (siguen abiertos).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-S4-r2.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-S4-r2.md` (nuevo, este reporte)

## Verificación
- `go vet ./internal/recetas/...` limpio. `gofmt -l internal/recetas cmd/motores` vacío. `go test -coverprofile`: **93.2 %**, 43 PASS, 0 FAIL y 0 SKIP. `go test ./cmd/motores/...` ok. `grep float32|float64`, `grep http` y `grep os.Getenv|LookupEnv|Environ` fuera de pruebas: 0 líneas cada uno.
- S4-01: 60 recetas, 60 ids y 60 nombres distintos. `Draft202012Validator`: 0 errores y 60 de 60 válidas. Las pruebas negativas (llave extra, porciones 4, unidad kg) dan 1 error cada una. Validación cruzada propia: 0 errores. `rec_genera_recetario.py --verifica`: coincide. `go test -run Recetario -v`: review 60 de 60 y production 60 de 60.
- S4-02: tabla propia con `fractions` contra `s4-escalado-recetario.json`: 240 consultas, 2244 ingredientes y 0 diferencias. Con 2→4, el doble exacto en todos.
- S4-03: recalculé las 3 fixtures `s4-lista-*` (Σ − existencia, paquetes y sobrante): 0 diferencias y 0 negativos. El total de la raíz es igual a la Σ de los artículos.
- S4-06: `s4-costo-recetas` contra `data.json`: 4425, 2213, 8850, 2491, 14701 y 18061.
- S4-08: 100 ingredientes únicos, 100 con fuente, 0 sin dato y 0 con fuente vacía (100 en `propuesto`). Recalculé las kcal de las 60 recetas: 0 diferencias. Los renglones de `s4-kcal-sin-dato` son idénticos a `tabla.csv`.
- S4-10: SHA-256 de los zip USDA correcto. 7 valores USDA desde los CSV locales y 1 de OFF por la API oficial: todos iguales.
- S4-11: 0 nombres iguales a `nombre_fuente`. Leí 20 nombres al azar, todos en español. Los 40 empates media o baja están en `revisar.tsv` con los mismos valores que `tabla.csv`.
- S4-12: 0 sin canal. S4-13: el campo existe (1) y está vacío (0) en una base armada desde 00001 a 00003; ahí `health.goal_adjust_*` vale 0 y 0. S4-14: `TestSinIA` pasa.
- Oráculo: dump con `datos/tmp/motores s4 dump`, después `rec_run_fixtures.py` (26 de 26) y `compara.py s4`: 8451 de 8451 en las 3 vías, provisionales 429 de 429, 0 diferencias. Los md5 de las fixtures, de `salidas/s4.json` y de `recetario.json` no cambiaron.

## S4 verificación manual
Las 5 recetas las elegí yo con `random.Random('revisor-S4-r2-2026-10-11')`: 2 de las 9 recetas nuevas de b y 3 del resto. Excluí las de la ronda 1 y las de la fixture del autor. Gramos: `g` = cantidad; `ml` o `pieza` = cantidad × `gramos_por_unidad` del catálogo. kcal = gramos × kcal_100g / 100. Por porción = HALF_UP(Σ / 2).

1. **pechuga-a-la-plancha-con-brocoli** («Pechuga a la plancha con brócoli y arroz blanco»; comida y cena; 15 + 25 min). Ingredientes:
   - pechuga 350 g × 112.202 (USDA 2646170) = 392.707
   - brócoli 300 g × 32 (USDA 321900) = 96
   - arroz 120 g × 365 (168877) = 438
   - zanahoria 80 g × 41 (170393) = 32.8
   - cebolla 40 g × 35 (1104962) = 14
   - ajo 6 g × 143 (1104647) = 8.58
   - aceite de oliva 20 ml × 0.913 = 18.26 g × 884 (171413) = 161.4184
   - limón 40 g × 30 (168155) = 12
   - pimienta 1 g × 251 (170931) = 2.51
   - sal 5 g × 0 (746775) = 0

   Σ = 1158.0154. Por porción 579.0077 → **579** (archivo 579). 2→4: 700000, 600000, 240000, 160000, 80000, 12000, 40000, 80000, 2000 y 10000, el doble exacto. Cantidades razonables para 2 adultos (175 g de pechuga por persona). Son 5 pasos y usan todos los ingredientes.
2. **costillas-de-cerdo-con-verdolagas** (comida y cena; 20 + 50 min). Ingredientes:
   - costilla 500 g × 277 (USDA 167853) = 1385
   - verdolagas 250 g × 20 (169274) = 50
   - tomate verde 400 g × 32 (168566) = 128
   - serrano 15 g × 32 (169395) = 4.8
   - cebolla 80 g × 35 = 28
   - ajo 8 g × 143 = 11.44
   - comino 1 g × 375 (170923) = 3.75
   - aceite de soya 15 ml × 0.921 = 13.815 g × 884 (171411) = 122.1246
   - sal 6 g × 0 = 0
   - tortilla 200 g × 218 (175036) = 436

   Σ = 2169.1146. Por porción 1084.5573 → **1085** (archivo 1085). 2→4 doble exacto. Son 5 pasos completos. Nota B7: la costilla se pesa con hueso.
3. **chilaquiles-verdes-con-pollo** (desayuno y comida; 15 + 25 min). Ingredientes:
   - tortilla 240 g × 218 = 523.2
   - pechuga 200 g × 112.202 = 224.404
   - tomate verde 300 g × 32 = 96
   - serrano 20 g × 32 = 6.4
   - cebolla 100 g × 35 = 35
   - ajo 4 g × 143 = 5.72
   - cilantro 10 g × 23 (169997) = 2.3
   - crema 80 g × 192.744 (2346387) = 154.1952
   - queso fresco 80 g × 297.417 (2647442) = 237.9336
   - aceite 60 ml × 0.921 = 55.26 g × 884 = 488.4984
   - sal 4 g × 0 = 0

   Σ = 1773.6512. Por porción 886.8256 → **887** (archivo 887). 2→4 doble exacto. Son 5 pasos completos. El aceite de freír se cuenta entero, lo que sobreestima un poco y es conservador.
4. **torta-de-jamon-con-aguacate** («Torta de jamón y queso con aguacate»; comida, cena y colación; 10 + 0 min). Ingredientes:
   - bolillo 2 piezas × 57 g = 114 g × 293 (175031) = 334.02
   - jamón de pavo 80 g × 106 (172941) = 84.8
   - manchego 60 g × 366 (OFF 7501077506585) = 219.6
   - aguacate 100 g × 160 (171705) = 160
   - jitomate 80 g × 18 (170457) = 14.4
   - mayonesa 20 g × 680 (171009) = 136
   - jalapeños 20 g × 27 (170080) = 5.4

   Σ = 954.22. Por porción 477.11 → **477** (archivo 477). 2→4 doble exacto (bolillo 4000). Son 5 pasos completos.
5. **fideo-seco-a-la-mexicana** (comida y cena; 10 + 20 min). Ingredientes:
   - fideo 200 g × 371 (169736) = 742
   - jitomate 300 g × 18 = 54
   - cebolla 60 g × 35 = 21
   - ajo 4 g × 143 = 5.72
   - aceite 30 ml × 0.921 = 27.63 g × 884 = 244.2492
   - cubo de caldo 1 pieza × 4 g = 4 g × 198 (171563) = 7.92
   - crema 60 g × 192.744 = 115.6464
   - cotija 40 g × 351.761 (2647443) = 140.7044
   - aguacate 100 g × 160 = 160
   - sal 2 g × 0 = 0

   Σ = 1491.24. Por porción 745.62 → **746** (archivo 746). 2→4 doble exacto. Son 4 pasos completos. 100 g de fideo seco por persona es una porción generosa pero razonable como plato fuerte.

Las 5 coinciden con `recetario.json`, con el motor (review, 0 rechazos) y con la fuente (los 4 valores Foundation coinciden con el CSV de USDA).

## S4-07 cálculo a mano
- Hombre de 30 años, 80 kg y 1750 mm: MSJ = 800 + 1093.75 − 150 + 5 = 1748.75. Por 6/5 da 2098.5 = 4197/2 → 2099.
  - Con «bajar» y la casa en 300: 1798.5 → 1799.
  - Con 100 propios: 1998.5 → 1999.
  - Con «mantener» y 300 y 300: 2099.
  - IMC = 80 / 3.0625 = 26.12 → 261 `overweight`.
- Mujer de 28 años, 60 kg y 1650 mm: MSJ = 600 + 1031.25 − 140 − 161 = 1330.25. Por 6/5 da 1596.3 → 1596.
  - Con «subir» y 250: 1846.3 → 1846.
  - IMC = 22.04 → 220 `normal`.
- RN-17, 35 años, 72.5 kg y `moderate`: MSJ = 1648.75. Por 31/20 da 2555.5625 = 40889/16 → 2556. Con «bajar» y 500 propios: 2056.
- Víspera del cumpleaños (29 años): MSJ = 1753.75. Por 6/5 da 2104.5 → 2105.
- Sin claves de la casa (ajuste inicial, K36): 2099 y 1596.

Todo es igual a `s4-salud-perfiles` y a `s4-salud-ajuste-inicial`.

## Supuestos
1. Cotejé contra los CSV locales de USDA (SHA-256 verificado) porque la API de USDA con DEMO_KEY respondió `OVER_RATE_LIMIT`. El de OFF lo cotejé contra su API oficial.
2. Para S4-13 usé `sqlite3` con las secciones Up porque `goose` no está instalado. Es la misma vía que la ronda 1 (B6).
3. No traté como platillos repetidos los pares de B8, porque tienen nombre y proteína principal distintos.

## Bloqueos
Ninguno nuevo. Siguen vigentes, y los respetan los artefactos, `docs/estado/bloqueos/S4-confirmacion-empates.md`, `S4-reglas-lista.md` y `S4-asignacion-mercado.md`.

## Pendientes
Ninguno de esta revisión. Los bajos B2, B3, B4, B6, B7, B8 y B9 quedan para el integrador o la puerta 2. Dejé `datos/tmp/motores` porque lo comparten otros agentes en curso (también lo usan `motor-s5.json` y `rev-s3r2`). Borré `datos/tmp/rev-s4-r2/`, `datos/tmp/s4r2.cov` y `datos/tmp/motor-s4.json`.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S4-interfaz.md` §3, §12 y §12.3 y `oraculo/README.md`: aplicar las solicitudes 1 y 3 de `correccion-S4-r1.md` (generador `rec_genera_recetario.py`, borradores a y b, nombre real `rec_run_fixtures.py`, orden de regeneración). Motivo: B4.
2. `oraculo/fixtures/s4/s4-kcal-recetas-autor.json` (por medio de `rec_genera_fixtures.py`): cambiar la `descripcion` para que diga que las recetas vienen de `recetario.json`. Motivo: B4.

## Paquetes instalados
Ninguno del sistema. `jsonschema` corrió de forma efímera con `uv run --no-project --with jsonschema`.

## Propuestas para Axel
- Puerta 2, B7 (cortes con hueso): (a) aceptar las kcal como están, sobreestimadas en esas recetas; (b) escribir las cantidades en peso comestible; (c) agregar un factor de parte comestible por producto, lo que exige cambiar el contrato y el esquema.
- Puerta 2, B8: decidir si se cambia alguna receta de los pares parecidos (huevos con nopales y huevos a la mexicana; chilaquiles verdes con huevo y con pollo) para tener más variedad.
- Puerta 2, B2: confirmar que SR Legacy 2018-04 (USDA FDC, CC0) se acepta junto a Foundation Foods.
- No hay valores K31 en este alcance.

## Consumo aproximado
- Acciones de herramienta: ≈ 38.
- Leído: ≈ 200 KB (protocolo 6, decisiones §D y §K 14, criterios 18, contrato S4 por partes 45, reportes 22, salidas de scripts y pruebas ≈ 95).
- Escrito: ≈ 37 KB (5 scripts ≈ 14 KB, ya borrados; `REVISION-S4-r2.md` 13 KB; este reporte ≈ 10 KB).
