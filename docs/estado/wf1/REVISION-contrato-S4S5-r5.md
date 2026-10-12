# Revisión del contrato S4/S5 · ronda 5 (WF-1, 11-oct-2026)

Agente: `revisor-contrato-S4S5-r5` (revisor independiente; no escribí nada de lo revisado y no corregí nada). Evalué solo contra `docs/criterios/S4.md` 1.2, `docs/criterios/S5.md` 1.2 y `docs/plan/decisiones.md` §D, §E y §K. Artefactos: `docs/contratos/motores/S4-interfaz.md` 1.4, `S5-interfaz.md` 1.4 y `productos-formato.md` 1.1. Reportes leídos: `contrato-S4S5.md` y `correccion-contrato-S4S5-r4.md`. No me basé en lo que dicen los reportes: volví a calcular todo con mi propio script (`datos/tmp/rev-s4s5-r5/verif.py`, `python3 -I`, solo `fractions`, `decimal`, `json` y `re`; borrado al terminar).

## Veredicto
**NO APROBADO**: 0 hallazgos altos, **3 medios** y 4 bajos.
- Los dos medios de la ronda 4 quedan **cerrados**.
- Ninguna cifra canon falla: **46 de 46** comprobaciones OK, incluidos +138.50 y +100.50 calculados desde los renglones de `data.json`.
- Los 3 medios son nuevos. Son preguntas que el implementador del motor o del oráculo tendría que adivinar:
  - M-r5-1 viene de la corrección r4: los artículos `manual` usan su propia elección.
  - M-r5-2 y M-r5-3 vienen de antes y ninguna ronda los había señalado.

## Verificación propia (resumen)
`python3 -I datos/tmp/rev-s4s5-r5/verif.py /home/deepsec/finanzas-hogar` → **46 OK, 0 FAIL**:
- **Canon S5-01, desde los renglones** (no desde `totales`):
  - idénticos: Σ cantidad × `bodega_unit` = 277350 y Σ cantidad × `walmart_unit` = 287400; diferencia **+10050 (+$100.50)**;
  - totales: 277350 + 30000 (sustitutos) + 2100 (sexta lata) = **309450** y 287400 + 35900 = **323300**; diferencia **+13850 (+$138.50)** = 10050 + 5900 − 2100;
  - los cuatro totales son iguales a `data.json: totales`;
  - piezas 61 / 60 / 55, renglones 33 / 3 y ganadores 22 / 8 / 3, iguales a `resultados_por_producto`;
  - los 33 renglones cumplen total = cantidad × unitario.
- **S5-02**, con el contenido del catálogo (§5.1 de `productos-formato.md`): 2933 / 2800, 12692 / 16406 y 3375 / 2833.
- **K40:** 1, 60/61, 30/61, 0 y «no compite» en los días 29, 30, 60, 90 y 91; se muestra 0.9836 y 0.4918; la recta baja 1/61 por día entre 30 y 90.
- **S4:**
  - §8.3: 4425 / 2213 / 8850. RN-15: 1486571.
  - S4-07: MSJ 1748.75 y 1330.25; valores exactos 2098.50 y 1596.30; metas 2099 / 1596 / 1799 / 1846. RN-17: 2556 y 2056.
  - IMC: 185 / 250 / 300 / 237 décimas.
- **S4 §6.6, con mi propia implementación:** arroz 3 / 2500000 / 8400; frijol 1 / 0 / 3000; aceite `manual` 2 / 1892000 / 7980; servilletas `null`. Grupo y raíz **19380**. Sin la cantidad fijada: 1 / 500000 / 2800 y 2 / 800000 / 6000.
- **S5 §14.1 y §7.3, con mi propia implementación:**
  - `quote_for`: sin elección → Walmart 1 kg; con Bodega → 750 g a 2200; con 750 g → Walmart, por el orden de `stores`; frijol → Bodega, por el precio del paquete.
  - `s5-cotizacion-recomendada`: totales Walmart 2200 y Bodega 1500, los dos sustitutos; `estimated_price_cents` 2800; `cheapest_store` Bodega.
  - `s5-canasta-regla-unica`: Walmart **9430** y Bodega **9460** (+30); `cheapest_store` Walmart; recomendación Walmart 1 kg (2800) y Walmart 946 ml (7980); arroz sustituto y aceite idéntico en las dos tiendas. Con la regla 1.3 serían 10780.
- **Catálogo (e):**
  - §5.1 tiene 36 productos y 39 presentaciones;
  - los 39 nombres de `data.json` (33 + 3 + 3) existen tal cual y no sobra ninguno;
  - §5.2 es igual a `cuentas.{axel,karen,conjunto}.lineas` en orden, cuenta, texto, piezas e importe: **39 renglones, 61 piezas, 3,094.50**;
  - cada presentación de §5.2 existe y es de su producto, y su precio unitario es `bodega_unit` (o total ÷ cantidad en sustitutos);
  - piezas por presentación = `cantidad` (atún 6 = 5 + 1).
- **Estructura:**
  - los 7 bloqueos citados existen (`S4-reglas-lista`, `S4-asignacion-mercado`, `S4-confirmacion-empates`, `S5-reglas-comparativa`, `S5-carrito-comparativa`, `S5-profeco-almacen`, `C1-campos-calculados`);
  - cada `-run` de los criterios (S4: Canal, Costo, Escalado, Inventario, Kcal, Lista, Plan, Salud; S5: Comparativa, Frescura, Fuente, LecturaAsistida, MejorTienda, Membresia, Mercado, PreciosViejos, Unidad) tiene pruebas en S4 §13 y S5 §16;
  - C1 `ShoppingListItem.reasons` admite `manual` y `needed_milli` es opcional (extracto `compras-tiendas-a.yaml`).

## Tabla de criterios

| Criterio | Estado en el contrato | Sección | Evidencia propia |
|---|---|---|---|
| S4-01 Recetario | cubierto | S4 §3 | esquema 2020-12 en `recetario.schema.json`, llave `recetas` = 60, validación cruzada, modos `review` y `production` |
| S4-02 Escalado exacto | cubierto | §4, §13 | fracción reducida, sin redondeo intermedio; ejemplos de RN-14 (375000, 10500, 1500) correctos |
| S4-03 Lista = Σ − stock | cubierto, con M-r5-1 y M-r5-2 | §6, §6.6, §12.2 | §6.6 reproducido (19380). Abiertos: la cotización con varias elecciones (M-r5-1) y la presentación sin cotización (M-r5-2) |
| S4-04 Inventario híbrido | cubierto | §7 | 1486571 recalculado; agotado prevalece, razonado |
| S4-05 Plan por persona | cubierto | §5 | Karen → `ErrPersonWithoutPlan`; «fuera» sin kcal; K38 y K43 |
| S4-06 Costo | cubierto, con M-r5-1 | §8, §6.6 | 4425 / 2213 / 8850; «sin precio» ≠ 0 |
| S4-07 IMC y meta | cubierto | §10 | MSJ, factores 6/5 … 19/10 (K23), ajuste K36, un redondeo; todas las cifras recalculadas |
| S4-08 kcal con fuente | cubierto | §9 | `tabla.csv` con `estado`; cobertura `WithSource = Unique − Missing` |
| S4-09 5 recetas a mano | cubierto (lo hace el revisor de la implementación) | §12 | el contrato no lo sustituye con fixture del autor |
| S4-10 Cargas locales | cubierto | §9.1, §1 | `PROCEDENCIA.md`, 3 fuentes con licencia, 0 red |
| S4-11 Español | cubierto | §3.2, §9.1 | `nombre_fuente` solo como procedencia; `production` no carga `propuesto` |
| S4-12 Canal | cubierto (mercado PROVISIONAL citado) | §11, §6.5 | bloqueo `S4-asignacion-mercado` existe |
| S4-13 Editables | cubierto | §3.3 | funciones puras y consulta SQL del campo vacío |
| S4-14 IA apagada | cubierto | §1, §13 | `TestSinIA` |
| S4-15 Calidad Go | cubierto | §1, §13 | 85 % (K25), sin `float` |
| S5-01 Comparativa canon | cubierto | §7.2, §3, §15 | +138.50 y +100.50 recalculados desde los renglones (arriba) |
| S5-02 Unidad base | cubierto | §4 | 2933 / 2800, 12692 / 16406, 3375 / 2833 |
| S5-03 Mejor tienda | cubierto (el hallazgo de la r4 está cerrado) | §7.3, §7.4, §14.1 | 9430 / 9460 → Walmart; 1500 / 2200 → Bodega; empates PROVISIONAL citados |
| S5-04 Frescura y fuentes | cubierto | §5, §3, §7.1 | K40 recalculado; K37 y K22 en §3 |
| S5-05 PROFECO | cubierto | §11 | conteo leídos = guardados + rechazados; ligas confirmadas |
| S5-06 OCR | cubierto | §12 | K24: 5 fotos, renglón correcto, «sin medir» |
| S5-07 Lectura asistida | cubierto | §13 | `CanEnable` exige veredicto, URL y fecha |
| S5-08 Sin red | cubierto | §1 | `asistida` sin cliente HTTP en WF-1 |
| S5-09 Membresía | cubierto | §8 | 100000 ÷ 12500 = 8; `missing_fee` |
| S5-10 Trazabilidad | cubierto | §3, §6 | 5 fuentes con «compite / referencia»; promoción PROVISIONAL |
| S5-11 Mercados | cubierto | §9, §7.4 | 2026-10-13 (martes siguiente al 10-oct) |
| S5-12 Precios viejos | cubierto (B-r4-1 abierto) | §10 | evalúa en la fecha de despensa |
| S5-13 Calidad Go | cubierto | §1, §16 | 85 %, sin `float` |

**(d) Coherencia con §K:**
- **K7:** IMC 18.5 / 25 / 30 (S4 §10.1); precio por 1 kg, 1 l o 1 pieza con HALF_UP (S5 §4); equilibrio en compras al año (S5 §8); reponer ≤ nivel y agotado en 0 (S4 §7).
- **K22:** 30 días, recta de 30 a 90, PROFECO nunca compite (S5 §3 y §5).
- **K23:** MSJ × 5 factores de 1.2 a 1.9 con ajuste por objetivo; la persona confirma (S4 §10.2 y §10.4).
- **K24:** 5 fotos en `datos/semilla/tickets/`; renglón correcto = descripción + cantidad + importe; «sin medir» (S5 §12).
- Sin contradicciones.

## Hallazgos

**M-r5-1 (media) · S4-03, S4-06 · `S4-interfaz.md` §8.1 paso 1 y §12.1 (`quotes`).**
- **El problema.** Desde la 1.4, un artículo `manual` usa su propia elección (§6.6 y §8.1), así que un producto puede tener varias elecciones en la misma lista: la del artículo `calculated` y la de uno o más `manual`. Cada elección trae su `QuoteFor`, y `quotes` las guarda sin decir a qué elección responde cada una. El paso 1 busca «la cotización cuya tienda es la elegida». Con una elección de solo tienda, ahí caben las cotizaciones que pidieron otras elecciones en esa misma tienda. Si la recomendada no está entre ellas, el resultado es `invalid_argument` con una entrada válida.
- **Ejemplo** (precios del mismo día, `stores` = [walmart, bodega]): arroz con Walmart 1 kg a 2800 y 750 g a 2250, y Bodega 500 g a 1350 y 750 g a 2300.
  - La recomendada es Bodega 500 g (2700 por kg).
  - El artículo `calculated` eligió `store` = walmart → `QuoteFor` = Walmart 1 kg.
  - Un artículo `manual` eligió `presentation` = 750 g → `QuoteFor` = Walmart 750 g.
  - En `quotes` quedan Bodega 500 (recomendada), Walmart 1 kg y Walmart 750 g.
  - Para el artículo `calculated` casan Walmart 1 kg y Walmart 750 g, y la recomendada no está entre ellas → `invalid_argument`. Lo correcto es Walmart 1 kg (2800).
  - S4 tampoco puede desempatar con el orden de `QuoteFor`, porque `quotes` no trae el peso.
- **Corrección sugerida.** Que cada cotización diga a qué elección responde, con `choice_store` y `choice_presentation` (`null` y `null` en la recomendada), y que el paso 1 tome la cotización con exactamente esa elección. Otra opción: `list_choices[]` trae su propia cotización. Ajustar §12.1, el caso `s4-costo-cotizacion-elegida` (con un artículo `manual` que elige otra cosa en el mismo producto) y `TestCostoCotizacionQueRige`.

**M-r5-2 (media) · S4-03 · `S4-interfaz.md` §6.2 paso 5, §12.1 (`catalog`, `inventory`, `recipes`).**
- **El problema.** Sin cotización, la presentación que rige es «la `presentacion_id` del ingrediente». Esa regla no da una sola presentación en tres casos que pasan al arranque, cuando la mayoría de los productos de las recetas (verduras, carne, huevo, especias) todavía no tiene precio:
  - (a) el producto viene de varias recetas con `presentation` distinta (la unicidad es por receta, y en recetas de la casa `Ingredient.PresentationID` puede ser nulo, §2);
  - (b) el artículo solo es de reposición por ritmo (motivo `replenishment`, sin ingrediente);
  - (c) un artículo `calculated` que sigue en la lista por una elección (cantidad o tienda) con necesidad 0 y sin ingrediente esa semana.
- **Qué queda sin definir.** `presentation`, `packages_to_buy` y `leftover_milli`. El `Input` de §12.1 no trae la `presentacion_tipica` del catálogo (`productos-formato.md` §2: «la que sugiere la receta y la lista cuando S5 no tiene precio») ni una presentación en `inventory`.
- **Ejemplo:** papel higiénico en `consumption_rate` con reposición de una semana y sin precio que compita. ¿Cuántos paquetes, de qué presentación y qué sobrante? El motor y el oráculo pueden dar respuestas distintas.
- **Corrección sugerida.** Una regla única de respaldo, por ejemplo:
  - la `presentacion_tipica` del producto, que se agrega a `catalog.products[]` como `typical_presentation`;
  - o, para (a), la primera receta en orden de fecha y tiempo;
  - o, como en el `manual`, `null` en los tres campos.
  Agregar el caso a `s4-lista-sin-existencia` o a una fixture nueva.

**M-r5-3 (media) · S4-03, S5-03 y oráculos · `S4-interfaz.md` §12.2, `S5-interfaz.md` §7.2, §7.4 y §15.2; `oraculo/fixtures/FORMATO.md` §4 y §5.**
- **El problema.** FORMATO compara las listas **por índice**. Las escribe «en su orden» calculado, y el largo `.#` también es hoja. Los contratos no fijan el orden de varias listas de salida que el motor Go y el oráculo Python, escritos por agentes distintos, tienen que emitir igual:
  - S4 `shopping_list.items[]`: junta recetas, reposición y renglones `manual`. Solo `groups[]` tiene orden.
  - S4 `missing_prices[]` de `shopping_list` y `recipe_detail`: tampoco dice si un producto con dos artículos sin precio sale una vez o dos.
  - S4 `week_plan.cells[]` y `leftovers[]`, e `inventory_status.items[]`.
  - S5 `compare_basket.items[]` y `items[].prices[]`.
  - S5 `compare_carts.rows[]`, `rows[].lines[]` y `winners[]`.
  - S5 `stale_prices_for_plan.stale[]` y `without_price[]`, y `profeco_import.coverage[]` y `rejections[]`.
  - S5 `stores[]` de `compare_basket` con tiendas incompletas. «Ordenadas por total ascendente» pondría primero a una tienda con `missing_count > 0` (su total es parcial), aunque «no puede ser la más barata». Tampoco se dice su `difference_vs_cheapest_cents`, que puede salir negativa, ni `cheapest_store` cuando ninguna tienda está completa.
- **Ejemplo:** una lista con arroz (receta), papel higiénico (reposición) y aceite (`manual`). Si el motor ordena por código de producto y el oráculo por orden de aparición, todas las hojas de `items[0..2]` salen como `DIF`, aunque los dos cálculos sean correctos.
- **Corrección sugerida.** Una línea de orden por lista, por ejemplo:
  - `items[]` por `(vertical, store)` en el orden de `groups[]`, luego `kind` (`calculated` antes que `manual`), luego código de producto y luego el orden de `list_choices`;
  - `missing_prices[]` por código, sin repetir;
  - `cells[]` por fecha, tiempo (`breakfast` … `snack`) y persona en el orden de `persons`;
  - `compare_basket.items[]` en el orden de `shopping_list` y `prices[]` en el orden de `stores`;
  - `rows[]` en el orden de `items_identicos`, luego `sustitutos` y luego `only_one`;
  - en `stores[]`, primero las completas por total y después las incompletas, con la diferencia `null` y `cheapest_store` `null` si no hay ninguna completa (PROVISIONAL, `S5-reglas-comparativa` pregunta 2).

**B-r5-1 (baja) · S4 §12.2 / S5 §15.2.** La clave `items[].quantity_milli` significa cosas distintas en las dos salidas. En S4 `shopping_list` es la elección guardada de la persona (o `null`); en S5 `compare_basket` es la cantidad a comprar ya resuelta. Está escrito, pero es fácil confundirlas en la API. *Corrección:* en S5, renombrarla a `buy_milli` o `purchase_quantity_milli`.

**Bajos de la ronda 4 que siguen abiertos (no se corrigieron en la r4 y no cambian):** B-r4-1 (S5 §10 sin la presentación elegida), B-r4-2 (S5 §7.3 «luego la de menor id» en vez de «la que va antes en el catálogo») y B-r4-3 (S4 §12.2, variante sin campo nuevo sin mercado dado de alta).

## Estado de los hallazgos anteriores

| Hallazgo (r4) | Estado | Evidencia |
|---|---|---|
| M-r4-1 [media] S5-03: la tienda recomendada se cotizaba con la presentación de menor precio por unidad base e inflaba su total | **cerrado** | S5 §7.3 «Total por tienda: una sola regla» (menor costo de paquetes en **todas** las tiendas); §14.1 viñeta `recommended`; §14 fila `CompareBasket`. El contraejemplo da Walmart 1450 contra Bodega 1460 → Walmart; `s5-canasta-regla-unica` da 9430 / 9460 → Walmart, recalculado por mí. La recomendación por renglón sigue siendo `QuoteFor` por unidad base, como pide S5-03 («por ingrediente (en unidad base)») |
| M-r4-2 [media] S4-03 y S4-06: `manual` y `quantity_milli` sin reglas | **cerrado** | (1) S4 §6.6: el `manual` lleva `needed_milli` `null`, `reasons` [`manual`], paquetes y costo de su `quantity_milli` obligatoria, y suma en `groups[]` y en la raíz. (2) En `calculated`, la cantidad a comprar es `quantity_milli` si no es nula, como dice C1 («NULL en calculated = comprar la necesidad»), y cambia los paquetes y el costo. (3) S5 §15.1 trae `quantity_milli` y `needed_milli` puede ser nulo, y §7.3 usa la misma regla. Ejemplo de 19380 recalculado. Efecto colateral nuevo: M-r5-1 |
