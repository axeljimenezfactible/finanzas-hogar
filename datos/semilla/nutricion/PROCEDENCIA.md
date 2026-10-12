# Procedencia de la tabla de nutrición (S4 §9.1)

Tabla: `tabla.csv` (150 renglones, uno por producto alimenticio del catálogo `../productos.json`). Se regenera con `python3 -I construir_tabla.py` (salida idéntica byte a byte en corridas repetidas). Todos los renglones salen con `estado = propuesto`: ninguno entra automático (R9, S4-11); Axel los pasa a `confirmado` en la puerta 2 (bloqueo `S4-confirmacion-empates`, opción A).

## Fuentes y licencias
| Fuente | Versión | Licencia | Renglones | Archivo local (ignorado por git) |
|---|---|---|---|---|
| USDA FoodData Central, Foundation Foods | CSV 2026-04-30 | CC0 1.0 (dominio público) | 43 | `raw/usda-foundation-2026-04-30/FoodData_Central_foundation_food_csv_2026-04-30.zip` |
| USDA FoodData Central, SR Legacy | CSV 2018-04 | CC0 1.0 (dominio público) | 99 | `raw/usda-srlegacy-2018-04/FoodData_Central_sr_legacy_food_csv_2018-04.zip` |
| Open Food Facts (volcado completo, filtrado a México) | descargado 2026-10-11 | ODbL 1.0 + DbCL 1.0 (atribución «Open Food Facts contributors» y compartir igual) | 8 | `raw/off-2026-10-11/en.openfoodfacts.org.products.csv.gz` |
| Tablas mexicanas (SMAE) por captura manual | no se usó | licencia sin verificar | 0 | no aplica |

SHA-256 de las descargas: Foundation `70457ee9d9342f43bda2010318c85f04210c689fdeb9cd2da4c513b0e8dbc655`; SR Legacy `b80817294b8850530aaedf2e515c02593b1824f763a0ff356e5c2081643e6fd0`; Open Food Facts `f72687ee8bc6522054fe69dbfda6b91902c16af1ec2e043cde27bc6c29ad8176`.
Sitios: fdc.nal.usda.gov (descargas CSV y consulta por `fdc_id`) y static.openfoodfacts.org. De las tiendas no se tomó nada.

## Reglas de empate producto a alimento
1. Se busca primero el alimento **crudo** en USDA. Cocido o enlatado solo cuando el producto del catálogo ya viene así (atún, sardina, elote en grano, duraznos, frijoles refritos) o la receta lo pida; hoy el catálogo guarda ingredientes crudos, por eso se usa crudo.
2. **Foundation Foods** si trae el alimento equivalente (misma especie y parte comestible) **con energía publicada**; si no, **SR Legacy**. Varios Foundation (p. ej. jitomate, sandía, tortilla de maíz, mayonesa) no publican energía y se pasan a SR Legacy.
3. **Parte comestible**: se usa el registro «crudo, con piel o sin piel» que corresponda a lo que se compra y se pesa; los cortes de carne se toman «con grasa» o «solo magro» según el corte (nota en `tabla_detalle.tsv`). Los kcal van por 100 g de la parte comestible, igual que el catálogo.
4. **Nombre**: el empate se hace por el nombre en inglés de USDA; `nombre_fuente` conserva ese nombre solo como procedencia y nunca se muestra (S4-11).
5. **Open Food Facts** solo para productos empacados sin equivalente USDA: queso manchego, mayonesa con chipotle, yogur bebible, petit suisse, chipotles en adobo, chocolate de mesa, queso panela y café molido. Se busca primero la misma marca y presentación del catálogo; si el catálogo no fija marca, se toma una marca de México y se anota. El `id_fuente` es el código de barras. Se guardan solo las filas usadas en `off_filas_usadas.tsv`.
6. **Energía**: nutriente 1008 (kcal); si falta, 2048 (Atwater específico); si falta, 2047 (Atwater general). Proteína 1003, grasa 1004, carbohidratos 1005. Redondeo ROUND_HALF_UP a 3 decimales. Un macro negativo se deja vacío.
7. **Confianza**: `alta` = mismo alimento y parte comestible; `media` = equivalente cercano o dato de OFF capturado por usuarios; `baja` = aproximación o datos de OFF contradictorios. Toda confianza `media` o `baja` va en `revisar.tsv` para Axel. Nada se inventó: lo que no tiene fuente confiable lo dice la nota.
8. **Sin fuente (10 productos)**: no son alimentos (limpieza, higiene, papel) y no llevan renglón en `tabla.csv`; salen en `revisar.tsv` con motivo «sin fuente» y en `tabla_detalle.tsv`.

## Archivos de esta carpeta
- `tabla.csv`: salida con el formato exacto de S4 §9.1 (cabecera de 12 columnas).
- `empates.tsv`: empate escrito a mano (producto, fuente, id, confianza, nota). Es la única entrada de criterio humano.
- `off_filas_usadas.tsv`: las 8 filas de OFF usadas (generada por `extraer_off.py` desde el volcado, lectura en flujo del .gz).
- `construir_tabla.py`, `extraer_off.py`: scripts reproducibles, solo biblioteca estándar, con `python3 -I`.
- `tabla_detalle.tsv` (los 160 productos con conjunto, energía usada, confianza y nota) y `revisar.tsv` (lista para Axel).
- `fuentes.md`: apunta a este archivo.
