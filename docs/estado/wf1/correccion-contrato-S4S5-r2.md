# correccion-contrato-S4S5-r2 · reporte final (WF-1, 11-oct-2026)

Corrige los 4 hallazgos medios de `docs/estado/wf1/REVISION-contrato-S4S5-r2.md` (no había altos). No reabre decisiones de §A–§K ni cambia cifras canon. Los bajos B-r2-1 a B-r2-6 y los bajos abiertos de r1 (B1, B2, B3, B4, B6, B7 y B8) quedan fuera, por instrucción.

## Hecho
- **M-r2-1 (S5-05, S5-06; `S5-interfaz.md` §15, §15.1, §15.2, §18):** `files[]`, `links_file` y `tsv_file` pasan a ser rutas **relativas a la raíz del repositorio**, siempre dentro de `oraculo/fixtures/s5/_datos/` (por ejemplo, `oraculo/fixtures/s5/_datos/profeco-muestra.csv`). Siguen la misma regla que los `*_file` del `Input`. El evaluador `func(input json.RawMessage, q Query)` de FORMATO §4.1 abre la ruta tal como viene, porque el motor y el oráculo corren desde la raíz (FORMATO §4: rutas relativas al directorio actual). No necesita la carpeta de fixtures y `main.go` no cambia. Una ruta absoluta, con `..` o fuera de `_datos/` es `invalid_argument`. Se aclara además que no son `$ref`.
- **M-r2-2 (S4-03, S4-06; `S4-interfaz.md` §8.1, §6.2 paso 5, §6.5, §12, §12.1, §13, §15; `S5-interfaz.md` §14):**
  - Cada cotización de `quotes` lleva `recommended` (booleano). Por producto con cotizaciones hay **exactamente una** recomendada y a lo más una por (`product`, `store`, `presentation`); si no, `invalid_argument`.
  - Se fija una sola **«cotización que rige»** para `recipe_detail`, `week_plan` y `shopping_list`:
    1. Si la persona eligió algo, la cotización que casa con su elección (tienda y/o presentación). Si casan varias, la recomendada. Si no casa ninguna, «sin precio», sin caer a la recomendada. Para recetas, plan y artículo calculado cuenta la elección del artículo `calculated`; un artículo `manual` usa la suya.
    2. Si no eligió nada, la recomendada.
    3. Si no hay cotizaciones, «sin precio».
  - Las cotizaciones no recomendadas solo sirven para atender `list_choices`. S4 no compara precios ni repite el peso de S5 §7.1.
  - El paso 5 de §6.2 toma la presentación de esa cotización.
  - S5 §14 dice que `QuoteFor` sin tienda ni presentación entrega la recomendada (§7.4, con su regla de empate).
  - Fixture nueva: `s4-costo-cotizacion-elegida`. Prueba nueva: `TestCostoCotizacionQueRige`.
- **M-r2-3 (S4-01, S4-08; `S4-interfaz.md` §3.2, §13, §14, §15):**
  - El generador escribe `kcal_por_porcion` y `fuentes_kcal[]` con `tabla.csv` en modo `review`. `fuentes_kcal[]` lleva un elemento por ingrediente con renglón en `tabla.csv`, en el orden de `ingredientes[]`, copiado tal cual del CSV.
  - Párrafo nuevo «Cargador del recetario y modos» (`LoadRecipeBook(path, catalog, nutrition, mode)`):
    - En **`review`** recalcula los dos campos y, si no coinciden, devuelve `ErrKcalMismatch` y esa receta no se carga.
    - En **`production`** no coteja e ignora los dos campos: carga las 60 recetas. Las que tengan un ingrediente `propuesto` salen «incompletas» (`kcal` ausente y `missing_nutrition`, §9.3), nunca con `ErrKcalMismatch`.
  - Confirmar un renglón no cambia `kcal_100g`, `fuente` ni `id_fuente`, así que el cotejo en `review` sigue valiendo.
  - Pruebas nuevas: `TestRecetarioProduccion` y `TestKcalCotejoRevision`. En §14, `ErrKcalMismatch` queda «solo en modo `review`».
- **M-r2-4 (S4-12, RN-16; `S4-interfaz.md` §6.5, §11, §12, §12.1, §12.2, §13, §15; `S5-interfaz.md` §7.4; `productos-formato.md` §3):**
  - Regla de la tienda de cada artículo, en orden:
    1. La tienda elegida por la persona.
    2. Vertical `chain`: la tienda de la cotización que rige; sin cotización, `store` nulo.
    3. Vertical `market`: PROVISIONAL (bloqueo nuevo `S4-asignacion-mercado`), opción A:
       - Entre los mercados que surten el grupo, el de `next_visit` más cercana; en empate, el primero en el orden de `stores` (la API los entrega por `id` ascendente).
       - Si ninguno surte el grupo, `store` nulo en la vertical `market` y la leyenda «sin mercado dado de alta».
       - Una cadena más barata según S5 no mueve el artículo: sigue visible en la comparativa y la persona puede elegirla.
  - El precio y la presentación salen siempre de la cotización que rige.
  - `groups[]`: un grupo por cada par (`vertical`, `store`). `visit_date` es la del mercado. `estimated_total_cents` = Σ `packages_to_buy × price_cents` de los artículos con precio; los que no tienen precio van a `missing_prices[]`. El orden es fijo (PROVISIONAL, mismo bloqueo).
  - En el `Input`, `stores[].next_visit` es obligatoria en mercados. Su `from` en la API (el más tardío entre `today` y el lunes de la semana de la lista) es PROVISIONAL.
  - S5 §7.4 y `productos-formato.md` §3 remiten a S4 §6.5.
  - Fixture nueva: `s4-canal-mercados` (con `provisional`). Prueba nueva: `TestCanalMercados`.

## Archivos
- `docs/contratos/motores/S4-interfaz.md`: 1.1 → 1.2; 56,588 bytes; 66 líneas de diff más el ajuste de los pasos 2 y 3 de §6.5.
- `docs/contratos/motores/S5-interfaz.md`: 1.1 → 1.2; 42,614 bytes; 23 líneas de diff.
- `docs/contratos/motores/productos-formato.md`: 1.0 → 1.1; 19,149 bytes; solo el encabezado y el último punto de §3.
- `docs/estado/bloqueos/S4-asignacion-mercado.md`: nuevo, con los 5 apartados del protocolo §2.
- `docs/estado/wf1/correccion-contrato-S4S5-r2.md`: este reporte.

## Verificación
- **Pruebas Go, cobertura, oráculo y `compara.py`:** no aplican todavía. El artefacto es un contrato sin código. No existen `backend/internal/recetas`, `backend/internal/precios`, `oraculo/compara.py`, `oraculo/fixtures/s4/` ni `oraculo/fixtures/s5/` (lo comprobé con `ls`). Lo que se verifica son las cifras y la estructura.
- **Cifras canon** (`python3 -I verifica.py`, con `Decimal` y `Fraction`, leyendo `datos/semilla/data.json` directamente): **24 de 24 OK**.
  - 0 renglones con unitario × cantidad ≠ total.
  - K-01 = 309450 y K-02 = 323300, iguales a `totales`.
  - K-03 = 13850 = 10050 + 5900 − 2100.
  - K-04 / K-05 = 277350 / 287400; K-06 = 10050.
  - K-07: 61 / 60 / 55. K-08: 33 / 3. Ganadores 22 / 8 / 3, iguales a `resultados_por_producto`.
  - Unidad base: 2933/2800, 12692/16406 y 3375/2833.
  - S4 §8.3: los renglones 21, 16, 15 y 0 son espagueti, atún, elote y aceite. Resultado 4425 / 2213 / 8850.
  - Frescura: 0.9836, 0.4918 y 0.0000. Fechas 2026-11-08, 2026-11-09 y 2027-01-08.
  - MSJ: 1748.75, 1330.25 y 1648.75. Metas: 2099, 1596, 1799, 1846, 2556 y 2056.
  - IMC en décimas: 237, 185 (con 18.4975 exacto), 250 y 300. RN-15: 1486571. RN-16: 2026-10-13.
  - Catálogo de `productos-formato.md` §5: 36 productos y 39 presentaciones. §5.2 tiene 39 renglones, 61 piezas y 309450.
- **Estructura** (`python3 -I estructura.py`):
  - 0 tablas con columnas desparejas en los 3 archivos (12, 12 y 4 tablas).
  - El bloque JSON de S4 §3.1 es válido.
  - Ninguna de las 41 pruebas de S4 mezcla dos temas de `-run`.
  - No queda ninguna frase vieja: «relativas a la carpeta del caso», «CSV en `_datos/`», «si no coincide, error `ErrKcalMismatch` y no se carga», «el cargador verifica que coincide», «S4 pide, por producto», «si S5 no tiene precio, la `presentacion_id`».
  - Todo bloqueo citado existe, salvo dos falsos positivos de S4 que ya estaban en 1.1: `S4-13` es un criterio, y `S4-meta-ajustes` es el bloqueo que S4 abriría solo si faltara la fuente de los factores (criterios S4-07 1.2).
  - Solo hay un PROVISIONAL sin bloqueo en su línea: «PROVISIONAL de S5-09» en S5 §8. Viene de 1.0 y el revisor lo aceptó.
- **Criterios tocados → evidencia:**
  - S4-01 → §3.2 (cotejo solo en `review`; `TestRecetarioProduccion`).
  - S4-03 y S4-06 → §8.1, §6.2 paso 5, §12.1 (`recommended`) y `s4-costo-cotizacion-elegida`.
  - S4-08 → §3.2 y §9.1 (el generador y el cotejo usan `review`; `TestKcalCotejoRevision`).
  - S4-12 → §6.5, §11 y `s4-canal-mercados` (PROVISIONAL `S4-asignacion-mercado`).
  - S5-05 y S5-06 → §15.1 y §15.2 (rutas desde la raíz).
  - S5-03 y S5-11 → §7.4 (remite a S4 §6.5).
  - El resto de los criterios queda con el estado que les dio el revisor r2.
- Temporales en `datos/tmp/corr-s4s5-r2/` (respaldos, parches y scripts de verificación): borrados al terminar.

## Supuestos
- S-r2-1 (S4 §8.1): las reglas «si casan varias, la recomendada; si siguen varias, `invalid_argument`» y «a lo más una cotización por (`product`, `store`, `presentation`)» son convenciones de las fixtures para que el motor y el oráculo den lo mismo. La API pide a S5 exactamente una cotización por elección. No cambian cifras, esquema ni C1.
- S-r2-2 (S4 §6.5): un grupo lleva la vertical del producto aunque la persona elija una tienda de otro tipo, como dice C1: `purchase_channel` = `product_group.purchase_channel`.
- S-r2-3: `LoadRecipeBook` es un identificador nuevo para el cargador del recetario, que ya existía sin nombre.

## Bloqueos
- `docs/estado/bloqueos/S4-asignacion-mercado.md` (nuevo; M-r2-4; 4 preguntas, opción A PROVISIONAL en cada una).
- Siguen vigentes: `S4-reglas-lista`, `S4-confirmacion-empates`, `S5-carrito-comparativa`, `S5-reglas-comparativa` y `S5-profeco-almacen`.

## Pendientes
- Bajos de r2 sin tocar, por instrucción: B-r2-1 (claves de `compare_carts` y `compare_basket`), B-r2-2 (detalles del TSV del OCR), B-r2-3 (solicitud condicionada a C1 `BasketComparison`), B-r2-4 (consulta de restricciones sobre una tabla vacía), B-r2-5 (sexta lata y 37 contra 39 renglones) y B-r2-6 (extremos de los 365 días).
- Bajos de r1 abiertos: B1, B2, B3, B4, B6, B7 y B8.

## Solicitudes fuera de alcance
1. `docs/estado/bloqueos/INDICE.md`: indexar `S4-asignacion-mercado`, junto con `S4-confirmacion-empates` y `S5-carrito-comparativa` si siguen sin índice. Motivo: bloqueo nuevo.
2. C1, vía `backend/internal/api/c1gen`, solo si en la puerta 2 queda la opción A de `S4-asignacion-mercado`:
   - En `ShoppingList.groups[].store_id`, agregar «ausente = artículos de mercado sin mercado dado de alta».
   - En `groups[].total_cents`, agregar «suma solo los artículos con precio».
   - Motivo: S4 §6.5.
3. Dueño de P-13 (U1/U4, lista de compras): leyenda «sin mercado dado de alta» para el grupo de mercado sin tienda. Motivo: S4 §6.5, PROVISIONAL.
4. `oraculo/fixtures/FORMATO.md` §4.1 (integrador; opcional): escribir que un evaluador que lee archivos auxiliares lo hace con rutas relativas al directorio actual, que es la raíz del repositorio. Así quedaría por escrito la convención que ya usan S4 §12.1 y S5 §15.1. No es necesario para cerrar M-r2-1.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- **`S4-asignacion-mercado`:**
  - Pregunta 1: opción A. Es el mercado al que van antes, sin columnas nuevas.
  - Pregunta 2: opción A. Mientras no haya mercados, la lista recuerda darlos de alta y el total usa el mejor precio de S5. Si prefiere que la lista se parezca a la compra de hoy (la uva en Bodega), opción B.
  - Pregunta 3: opción A. RN-16 manda en la lista y la persona puede elegir la cadena más barata con un toque.
  - Pregunta 4: opción A. No cambia C1.
- Sin valores K31 nuevos.

## Consumo aproximado
- ~41 acciones de herramienta.
- ~170 KB leídos (por partes): protocolo, §D, §E y §K de decisiones, criterios S4 y S5, revisión r2, reporte r1, secciones de S4, S5 y `productos-formato.md`, FORMATO, RN-16 y RN-17, y extracto C1.
- ~45 KB escritos con herramientas: parches ~20 KB, bloqueo 3.6 KB, scripts de verificación ~8 KB y este reporte ~10 KB. Los parches reescribieron completos S4 (56.6 KB), S5 (42.6 KB) y `productos-formato.md` (19.1 KB) con `python3 -I`.
