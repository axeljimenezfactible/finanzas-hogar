# S5 · Interfaz del motor de precios y tiendas (`backend/internal/precios`)

Versión 1.4 · 11-oct-2026 · WF-1, agente `contrato-S4S5` (Opus); corregida por `correccion-contrato-S4S5-r1` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r1.md`) , por `correccion-contrato-S4S5-r2` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r2.md`) , por `correccion-contrato-S4S5-r3` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r3.md`) y por `correccion-contrato-S4S5-r4` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r4.md`). Documento de contrato, **sin código**. Lo implementan los agentes de WF-1 del plan §3.4 fila S5 (comparador Opus; PROFECO y OCR Sonnet).

**Fuentes (en este orden):** `docs/plan/decisiones.md` §D, §E, K7, K22, K24, K25, K37, K40, K44; `docs/contratos/reglas-negocio.md` 1.3 (RN-12, RN-13, RN-16, RN-20; C-01, C-02, C-05); esquema 1.3 (`store`, `product_group`, `store_product_group`, `product`, `presentation`, `price`, `attachment`, `app_setting`); C1 v1.2 por extractos (`compras-tiendas-a.yaml`: `BasketComparison`, `PriceAt`, `ShoppingListItem`; `compras-tiendas-b.yaml`: `Price`, `PriceSummary`; `ajustes-c.yaml`: `importProfeco`); criterios `docs/criterios/S5.md` 1.2; `datos/semilla/data.json`; `datos/semilla/profeco/QQP_diccionario.csv`.

**Marcas.** «PROVISIONAL (bloqueo X)» = regla no decidida en §A–§K ni en RN; queda con la opción indicada hasta la puerta 2 y **no cambia ninguna cifra canon**. «ILUSTRATIVO» = número de fixture que no viene del reporte ni del oráculo.

**Cambios de la versión 1.4 (11-oct-2026, ronda 4 de revisión):**
- M-r4-1: en la comparativa por canasta, el total de **cada** tienda se arma con una sola regla: el menor costo de paquetes para cubrir la cantidad a comprar, también en la tienda recomendada del renglón. La recomendación por renglón sigue siendo `QuoteFor` sin elección, la misma de S4. La línea de la tienda recomendada en `prices[]` puede mostrar otra presentación, como sustituto. Cambian el ejemplo de §7.3 y el caso `s5-cotizacion-recomendada`: Walmart cotiza 750 g a 2200 y Bodega, la más barata, 1500. Caso nuevo `s5-canasta-regla-unica` con el contraejemplo de la revisión (§7.3, §14, §14.1, §15, §15.2, §16, §18).
- M-r4-2: el renglón de `shopping_list` trae `quantity_milli` (la cantidad fijada por la persona), y `needed_milli` puede ser nulo (artículo `manual`). La cantidad a comprar sigue la regla de `S4-interfaz.md` §6.6, y con ella se calculan `estimated_price_cents` y el total de cada tienda (§7.3, §15.1).

**Cambios de la versión 1.3 (11-oct-2026, ronda 3 de revisión):**
- M-r3-2: una sola regla para `QuoteFor` (§14.1). La fecha es «hoy» en America/Mexico_City, salvo en §10. Sin elección gana el par (tienda, presentación) de menor precio por unidad base ajustado; con tienda o con presentación, el mismo criterio dentro de esa elección. El desempate no depende de la lista: menor precio del paquete, luego el orden de `stores` y luego el del catálogo (PROVISIONAL, bloqueo `S5-reglas-comparativa`, pregunta 2). `CompareBasket` recomienda cada renglón con `QuoteFor` y cotiza la tienda recomendada con esa misma presentación: hay una sola presentación recomendada (§7.3, §7.4; la 1.4 cambia cómo se cotiza la tienda recomendada, M-r4-1). `StalePricesForPlan` usa la misma regla (§10). Consulta nueva `quote_for` y caso `s5-cotizacion-recomendada` (§15, §15.2, §16, §18).
- M-r3-1: los campos calculados de `ShoppingListItem` los llena la API con la salida de S4 (`S4-interfaz.md` §12.2, traducción a C1), no con `CompareBasket` (§14).

**Cambios de la versión 1.2 (11-oct-2026, ronda 2 de revisión):**
- M-r2-1: `files[]`, `links_file` y `tsv_file` de `profeco_import` y `ocr_parse_tsv` son rutas relativas a la raíz del repositorio, dentro de `oraculo/fixtures/s5/_datos/`, con la misma regla que los `*_file` de §15.1; el evaluador no necesita la carpeta de fixtures y `main.go` no cambia (§15, §15.1, §15.2).
- M-r2-2: `QuoteFor` sin tienda ni presentación entrega la cotización recomendada, que S4 recibe con `recommended = true` (§14).
- M-r2-4: la recomendación por ingrediente no mueve un artículo de vertical `market` a otro grupo de la lista (`S4-interfaz.md` §6.5, PROVISIONAL, bloqueo `S4-asignacion-mercado`) (§7.4).

**Cambios de la versión 1.1 (11-oct-2026, ronda 1 de revisión):**
- M1: el canon es la compra registrada de Bodega contra el carrito de Walmart (`CompareCarts`); nada se siembra como gasto de Walmart; se declara la fuente de cada semilla (K37); cómo arma la API la canasta de las otras tiendas con `basket = expense` es PROVISIONAL (bloqueo `S5-carrito-comparativa`) (§2, §3, §7.2, §14, §15).
- M2: `Input` de las fixtures (§15.1), argumentos y salida de cada función (§15.2) y oráculo `oraculo/s5_run_fixtures.py` (§15.3).
- M4: fotos solas en `datos/semilla/tickets/` (`.jpg`, `.jpeg` o `.png`) y transcripciones en `datos/semilla/tickets-esperado/` (§12).
- M5: `vinculos.csv` con `estado`, `confirmo` y `fecha_confirmacion`; en producción solo cuentan las ligas confirmadas; quién confirma es PROVISIONAL (bloqueo `S4-confirmacion-empates`) (§11).
- M6: `identical_count` = renglones; K-07 sale en `stores[].piece_count` e `identical_pieces`; solicitud a C1 (§7.2).

## 1. Paquete, límites y dependencias

| Tema | Regla |
|---|---|
| Rutas | `backend/internal/precios` (núcleo puro), `backend/internal/precios/profeco` (cargador), `backend/internal/precios/ocr` (Tesseract), `backend/internal/precios/asistida` (interruptor de lectura asistida). Subcomando `s5 dump` del binario común `backend/cmd/motores` (§15, FORMATO §4). |
| Dependencias | Solo biblioteca estándar (+ modernc.org/sqlite solo en el cargador PROFECO y en la capa de datos). OCR y extracción de RAR por `os/exec` **sin shell** con herramientas locales (`tesseract`, `unrar` o `bsdtar`, instaladas y verificadas el 11-oct-2026: tesseract 5.5.3 con `spa`). |
| Red | Ninguna en ejecución. `net/http` solo podría aparecer en `precios/asistida`, que en WF-1 **no implementa cliente HTTP** y está apagada (S5-07, S5-08). Ninguna URL de tiendas en el código. |
| Aritmética | Centavos `int64` y milésimas `int64`; comparaciones de precio por unidad base con **fracciones exactas** (multiplicación cruzada); sin `float32`/`float64` fuera de `_test.go` (S5-13). Redondeo ROUND_HALF_UP solo al publicar. |
| Hoy | Lo pasa el llamador (America/Mexico_City, C-05). La edad de un precio usa días de calendario. |
| Parámetros | Se leen de `app_setting` (K44), nunca constantes en código: `price.stale_after_days` = 30, `price.zero_weight_after_days` = 90, `price.assisted_reading_enabled` = false, `price.profeco_import_max_gb` = 1. |
| Cobertura | ≥ 85 % en `internal/precios` y subpaquetes (K25, S5-13). |

## 2. Tipos del dominio

| Tipo | Campos | Esquema |
|---|---|---|
| `Store` | `ID`, `Kind` (`chain`, `branch`, `market`), `ParentID *`, `Name`, `ProfecoName *`, `Municipality *`, `VisitWeekday *int` (1 = lunes … 7 = domingo), `RequiresMembership`, `HasMembership`, `MembershipFeeCents *int64`, `AssistedReading bool`, `Groups []int64` (de `store_product_group`) | `store`, `store_product_group` |
| `Product`, `Presentation` | como en `S4-interfaz.md` §2 (`BaseUnit`, `SizeMilli`) | `product`, `presentation` |
| `PriceObs` | `ID`, `PresentationID`, `StoreID`, `PriceCents`, `ObservedDate`, `Source` (`ticket_manual`, `ticket_ocr`, `manual`, `site_assisted`, `profeco`), `IsPromotional`, `SourceRef *`, `CreatedAt` | `price` |
| `Settings` | `StaleAfterDays` (30), `ZeroWeightAfterDays` (90), `AssistedReadingEnabled` (false) | `app_setting` |
| `Weight` | fracción exacta en [0, 1] + `Competes bool` + `IsStale bool` + `AgeDays int` | calculado |
| `Quote` | `ProductID`, `StoreID`, `PresentationID`, `PriceCents`, `SizeMilli`, `ObservedDate`, `IsStale`, `Weight` | calculado (lo consume S4 por `PriceSource`) |
| `Cart` | `Code`, `StoreID`, `Kind` (`expense` = compra registrada; `cart` = carrito que no se compró), `ExpenseID *int64` (solo `expense`), `ObservedDate`, `Source` (fuente de sus precios, §3), `Lines []CartLine{PresentationID, QuantityMilli (presentaciones × 1000, como `expense_item.quantity_milli`), UnitPriceCents}` | compra: `expense` + `expense_item`; el carrito no se guarda (§3) |

## 3. Ingesta por fuente (K22, K37; S5-04, S5-10)

| Fuente (`price.source`) | Cómo entra | ¿Compite? |
|---|---|---|
| `ticket_manual` | renglones de un ticket capturados a mano (`expense_item`) al confirmar el gasto | **sí**, peso de frescura |
| `ticket_ocr` | renglones del OCR **después** de la confirmación humana (§12); el borrador vive en `attachment.ocr_draft_json` (K3) y no crea precios | **sí**, el mismo peso que el ticket |
| `manual` | captura a mano sin ticket (`createPrice`) | **sí**, el mismo peso y la misma frescura (K37) |
| `site_assisted` | lectura asistida de sitios; apagada (§13) | **no**: solo referencia hasta que la cadena tenga veredicto de términos (K37) |
| `profeco` | cargador PROFECO (§11) | **nunca**: solo referencia (K22) |

**Validación común:** `price_cents` entero ≥ 0; `observed_date` ISO válida y **no futura** respecto a «hoy» (supuesto S-7: una fecha futura es error de captura); presentación y tienda existentes; `UNIQUE (presentation_id, store_id, observed_date, source, is_promotional)` (repetir devuelve conflicto). Cada precio guarda tienda, fecha, fuente y referencia (S5-10).

**Semilla del 10-oct-2026 (`data.json`) para el motor y sus fixtures.** `data.json` dice «Precios del carrito ya con descuentos e IVA; sin envío» (`supuestos`) y solo la compra de Bodega Aurrera se registró y se repartió (`cuentas`: 39 renglones, 61 piezas, 3,094.50). El canon se modela como **una compra registrada contra un carrito**:
- **Bodega Aurrera, compra registrada** (`Cart` de tipo `expense`: `expense` + `expense_item`, los 39 renglones de `productos-formato.md` §5.2), `observed_date = 2026-10-10`. Fuente de sus precios: PROVISIONAL (bloqueo `S5-carrito-comparativa`, pregunta 2), opción A: `manual` (captura a mano sin ticket, K37), porque no consta un ticket; con `ticket_manual` el canon es el mismo (compiten igual, K37).
- **Walmart, carrito que no se compró** (`Cart` de tipo `cart`): precios `manual` del 2026-10-10 (K37: captura a mano sin ticket) y las cantidades del carrito (las de `items_identicos` y `sustitutos[].cant_w`; 5 latas de atún). Existe **solo como entrada del motor y de la fixture**: **no** se siembra como gasto (un gasto de 3,233.00 que no ocurrió sería un hecho no documentado). Si la semilla de la app guarda además esos precios en `price` es PROVISIONAL (bloqueo `S5-carrito-comparativa`, pregunta 3), opción A: no; el canon no depende de eso.
- Precio de una presentación = `bodega_unit` / `walmart_unit` de `items_identicos` y, en sustitutos, `total ÷ cantidad` (todos dan centavos exactos: verificado, 0 renglones con `unitario × cantidad ≠ total`). La sexta lata de atún de Bodega (`atun_sexta_pieza_bodega` = 21.00) es un renglón más de la compra de Bodega a precio regular.

## 4. Equivalencia por unidad base y sustitutos (RN-20, K7; S5-02)

- Precio por unidad base exacto = `price_cents / size_milli` (fracción). Para **comparar** se usa la fracción (multiplicación cruzada), nunca el valor redondeado.
- `base_unit_price_cents` (C1 `PriceAt`) = precio por **1 kg, 1 l o 1 pieza** con HALF_UP: `HALF_UP(price_cents × 1000000 / size_milli)` en g y ml; `HALF_UP(price_cents × 1000 / size_milli)` en pieza.
- **Idéntico** = misma presentación. **Sustituto** = otra presentación del mismo producto (todas las presentaciones de un producto son sustitutas entre sí: `product`); se marca `is_substitute = true` con nota.
- **Regresión de `data.json` (verificada con `fractions`):**

| Producto | Bodega Aurrera | Walmart | Más barato (exacto) | Nota de `data.json` |
|---|---|---|---|---|
| Arroz | 3 × 750 g a 66.00 → **2933** ($29.33/kg) | 3 × 1 kg a 84.00 → **2800** ($28.00/kg) | Walmart | «$29.33/kg vs $28.00/kg» |
| Sustituto de crema para café | 780 g a 99.00 → **12692** ($126.92/kg) | 640 g a 105.00 → **16406** ($164.06/kg) | Bodega | «$127/kg vs $164/kg» (la nota redondea a pesos; el motor publica 126.92 y 164.06) |
| Leche | 4 × 1 l a 135.00 → **3375** ($33.75/l) | 6 × 1 l a 170.00 → **2833** ($28.33/l) | Walmart | «$33.75/l vs $28.33/l» |

## 5. Frescura y peso (RN-12, K22, K40, K44; S5-04)

- `AgeDays` = días de calendario entre `observed_date` y «hoy». `IsStale` = `AgeDays ≥ Dv`, con `Dv = price.stale_after_days` (30).
- `Weight` exacto, con `D0 = price.zero_weight_after_days` (90): `1` si `AgeDays < Dv`; `(D0 − AgeDays) / (D0 − Dv + 1)` si `Dv ≤ AgeDays ≤ D0`; si `AgeDays > D0`, **no compite** y solo se muestra marcado. Con 30 y 90 es exactamente (90 − días) / 61 (K40).
- `Competes` = fuente que compite (§3) **y** `Weight > 0`. El día 90 el peso es 0: se muestra marcado y no influye (no hay división entre cero en §7).
- Se muestra con 4 decimales HALF_UP; se calcula siempre con la fracción.
- **Fronteras con un precio del 10-oct-2026 (verificadas):**

| Hoy | Días | ¿Viejo? | Peso exacto | Peso mostrado |
|---|---|---|---|---|
| 2026-11-08 | 29 | no | 1 | 1.0000 |
| 2026-11-09 | 30 | sí | 60/61 | 0.9836 |
| 2026-11-10 | 31 | sí | 59/61 | 0.9672 |
| 2026-12-09 | 60 | sí | 30/61 | 0.4918 |
| 2027-01-08 | 90 | sí | 0 | 0.0000 (no compite) |
| 2027-01-09 | 91 | sí | — | no compite, marcado |

La diferencia entre días consecutivos de 30 a 90 es exactamente 1/61 y el peso nunca sube con la edad.

## 6. Precio vigente

- Por presentación y tienda, el **último precio regular** (`is_promotional = 0`) de una fuente que compite (RN-12). Si dos precios que compiten tienen la misma fecha, rige el de `created_at` más reciente (supuesto S-8: «la última captura»).
- Los precios de referencia (`profeco`, `site_assisted`) se devuelven aparte (`References []PriceObs`) para mostrarlos; **nunca** entran en ganadores, totales de canasta ni costo de S4.
- **Promociones** (`is_promotional = 1`): RN-12 las excluye de la comparativa. Su vigencia no tiene columna en el esquema (`promotion` es de promociones de crédito); PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 3), opción A: una promoción solo es vigente en su `observed_date`; se muestra junto al precio regular con la leyenda «promoción del <fecha>» y una vencida no se muestra ni cuenta (S5-10).
- `PriceSummary` (C1): por producto, `fresh` si su mejor precio que compite no es viejo, `stale` si todos los que compiten son viejos, `no_price` si no hay ninguno que compita; por tienda, la fecha más reciente.

## 7. Comparativa: por ingrediente y por canasta (S5-01, S5-03, S5-11)

### 7.1 Aplicación del peso al elegir ganador
PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 1), opción A: para **ordenar** se usa el precio por unidad base ajustado = precio por unidad base ÷ peso (un precio de peso 30/61 tiene que ser menos de la mitad de caro para ganar a uno fresco); los **totales, diferencias y pantallas usan siempre el precio real**. Con precios del mismo día en todas las tiendas (canon del 10-oct) el ajuste no cambia nada.

### 7.2 Modo «compra contra carrito» (`CompareCarts`; canon del reporte)
Compara una canasta base (la compra registrada) contra uno o más carritos u otras compras (`Cart`, §2), cada uno de **una** tienda, renglón contra renglón:
1. Se agrupan los renglones de cada canasta por presentación (la sexta lata de atún se suma a las otras 5: 6 en Bodega).
2. **Idéntico:** la misma presentación en las dos canastas; la cantidad emparejada es la menor de las dos (atún: min(6, 5) = 5) y el excedente queda como renglón «solo en una tienda» (1 lata en Bodega a 21.00).
3. **Sustituto:** presentaciones distintas del mismo producto; cada tienda conserva **su propia cantidad y su propio importe** (así lo hace el reporte: 3 × 750 g contra 3 × 1 kg) y el ganador del renglón se decide por precio por unidad base exacto (§4).
4. Total de tienda = Σ de todos sus renglones (idénticos + sustitutos + solo en esa tienda). Total de idénticos = Σ cantidad emparejada × precio de esa tienda.
5. `winner_store_id` por renglón: menor precio unitario (idénticos) o por unidad base (sustitutos); ausente si empata. **`winners[]` y `ties` cuentan solo los renglones idénticos** (así reproducen `data.json: resultados_por_producto`; supuesto S-9); los sustitutos llevan su ganador en el renglón.
6. **Conteos (K-07, K-08).** `stores[].identical_count` (C1) = **renglones** idénticos (33 en cada tienda), como el propio C1 resume el canon en `BasketComparison` («33 idénticos, 3 sustitutos»); `stores[].substitute_count` = renglones sustitutos (3). Las piezas de K-07 salen con dos claves del motor y del dump: `stores[].piece_count` = piezas de la canasta completa de esa tienda (Σ `quantity_milli` ÷ 1000: Bodega **61**, Walmart **60**) e `identical_pieces` = piezas emparejadas de los renglones idénticos (**55**). Un renglón con cantidad no entera (a granel) cuenta como 1 pieza (supuesto S-11; ninguna cifra canon lo usa). C1 no tiene estas dos claves y describe `identical_count` como «Artículos idénticos.»: solicitud a C1 para que diga «Renglones idénticos.» y, si se quiere K-07 en pantalla, agregue `piece_count` e `identical_pieces`.

**Cifras canon y cómo se obtienen (verificado contra `data.json` con aritmética exacta):**

| ID | Cifra | Cómo se obtiene |
|---|---|---|
| K-01 | Bodega **3,094.50** | idénticos 2,773.50 + sustitutos 300.00 (66.00 + 99.00 + 135.00) + sexta lata de atún 21.00 |
| K-02 | Walmart **3,233.00** | idénticos 2,874.00 + sustitutos 359.00 (84.00 + 105.00 + 170.00) |
| K-03 | Walmart **+138.50** | 3,233.00 − 3,094.50 = +100.50 (idénticos) + 59.00 (sustitutos: 359.00 − 300.00) − 21.00 (lata que solo compró Bodega) |
| K-04 / K-05 | idénticos **2,773.50 / 2,874.00** | Σ de los 33 renglones de `items_identicos` (55 piezas) con las mismas cantidades en las dos tiendas |
| K-06 | Walmart **+100.50** en idénticos | 2,874.00 − 2,773.50: la lectura de plan §3.4 S5 es la de **mismos artículos, mismas cantidades**; +138.50 es la de **carrito contra carrito** |
| K-07 | **61 / 60 / 55** piezas | Bodega 55 + 5 (sustitutos: 3 + 1 + 1) + 1 (atún) = 61; Walmart 55 + 5 = 60 |
| K-08 | **33 / 3** renglones | `len(items_identicos)` / `len(sustitutos)` |
| — | ganadores **22 / 8 / 3** | por precio unitario en los 33 idénticos: Bodega gana 22, empatan 8, Walmart gana 3 (= `resultados_por_producto`) |

**API (C1 `getBasketComparison`, `basket = expense`).** C1 recibe **una sola** `expense_id`: la canasta base es esa compra. Cómo se arma la canasta de cada **otra** tienda no está en §A–§K, RN ni C1: PROVISIONAL (bloqueo `S5-carrito-comparativa`, pregunta 1), opción A: `CompareExpense(expense, date)` cotiza cada otra tienda con sus precios vigentes que compiten (§6) para los productos y cantidades de la compra base, elige presentación como en §7.3 y llama a `CompareCarts`. Con la opción A la API no reproduce +138.50 con la semilla (Walmart cotizado con las cantidades de Bodega da otro total); el canon K-01 a K-08 lo reproducen `CompareCarts` y la fixture `s5-comparativa-canon`. Opción B del bloqueo: C1 agrega un segundo parámetro (`compare_expense_id` o un carrito guardado) y la API reproduce el canon.

### 7.3 Modo «canasta por necesidad» (`basket = shopping_list`)
Entrada: los artículos de la lista de compras de S4, un renglón por artículo (producto, `needed_milli`, `quantity_milli` y la presentación que rige; §15.1). **Cantidad a comprar** = `quantity_milli` si no es nula; si es nula, `needed_milli`. Es la misma regla de `S4-interfaz.md` §6.6; un artículo `manual` trae `needed_milli` nulo. Para cada tienda y renglón, entre las presentaciones con precio que compite en esa tienda, se elige la que cubre la cantidad a comprar con el **menor costo de paquetes completos** = techo(cantidad a comprar / contenido) × precio (empate: la de menor precio por unidad base, luego la de menor id). Si la presentación elegida es la recomendada de la lista, el renglón es idéntico; si no, sustituto. Sin precio que compita → `missing_count` de la tienda. `stale_count` = renglones cuyo precio es viejo.

**Total por tienda: una sola regla.** El total de la canasta de **cada** tienda, también el de la tienda recomendada de un renglón, suma lo que esa tienda cobra con la regla de arriba: el menor costo de paquetes completos para la cantidad a comprar. Así `stores[].total_cents`, `cheapest_store` y `difference_vs_cheapest_cents` comparan lo que cuesta cubrir la lista completa en cada tienda (S5-03: «por canasta (lista de compras completa)»).

**Recomendación por renglón.** La recomendación de cada renglón (`items[].recommended_store`, `recommended_presentation`, `estimated_price_cents` e `is_stale_price`) es la cotización de `QuoteFor(product)` sin elección, con la misma fecha (§14.1). Es la misma cotización `recommended = true` con que S4 calcula `packages_to_buy` y el precio del artículo (`S4-interfaz.md` §6.2 paso 5, §6.6 y §8.1). `estimated_price_cents` = techo(cantidad a comprar / contenido) × precio del paquete de esa cotización, e `items[].quantity_milli` = la cantidad a comprar. Así cada artículo tiene una sola presentación recomendada, que es la que rige en la lista (los campos de `ShoppingListItem` salen de S4, §14). La regla de la canasta no cambia esa recomendación. En `items[].prices[]`, la línea de la tienda recomendada puede mostrar **otra presentación** (la que cubre la cantidad con menos costo), marcada como sustituto (`is_substitute = true`). En ese caso, lo que el renglón suma al total de esa tienda es menor que `estimated_price_cents`.

Ejemplo ILUSTRATIVO (caso `s5-cotizacion-recomendada`): arroz, necesidad 500 g. Walmart tiene 1 kg a 28.00 (2800 por kg) y 750 g a 22.00 (2933 por kg); Bodega, 750 g a 22.00 y 500 g a 15.00.
- Se recomienda Walmart 1 kg: un paquete, `estimated_price_cents` 2800.
- En la canasta, Walmart cotiza la bolsa de 750 g a 2200 y Bodega la de 500 g a 1500, las dos como sustituto. Bodega es la más barata y Walmart queda +7.00.

Contraejemplo de la revisión r4 (caso `s5-canasta-regla-unica`): Walmart tiene 1 kg a 28.00 y 500 g a 14.50; Bodega, 500 g a 14.60.
- Se recomienda Walmart 1 kg (2800 por kg, contra 2900 y 2920).
- Pero Walmart cubre los 500 g por 1450 y Bodega por 1460, así que en ese renglón la canasta favorece a Walmart.
- Si la tienda recomendada se cotizara con la presentación recomendada, Walmart saldría en 2800 y la canasta señalaría a Bodega, que cuesta más.

### 7.4 Mejor tienda por ingrediente y por canasta
- **Por ingrediente:** la tienda con menor precio por unidad base ajustado (§7.1) entre las que compiten para ese producto, con la presentación y el desempate de §14.1 (es la cotización recomendada de `QuoteFor`): todas las cadenas o sucursales con precio y **solo los mercados que declaran el grupo** del producto en `store_product_group` (RN-16, S5-11). Sam's compite como si tuvieran membresía (`assumes_membership = true` si `requires_membership = 1` y `has_membership = 0`). La recomendación por ingrediente no mueve un artículo de vertical `market` a otro grupo de la lista de compras: la tienda y el grupo de cada artículo los fija `S4-interfaz.md` §6.5 (PROVISIONAL, bloqueo `S4-asignacion-mercado`).
- **Por canasta:** tiendas ordenadas por total ascendente; `difference_vs_cheapest_cents` contra la más barata. En la canasta del 10-oct la más barata es Bodega Aurrera (3,094.50 < 3,233.00).
- **Empates y canastas incompletas:** PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 2), opción A: hay empate cuando los precios ajustados son iguales como fracción; por ingrediente, `winner_store_id` queda ausente y suma en `ties`, y la tienda **recomendada** del renglón sale del desempate de §14.1, que no depende de la lista (menor precio del paquete, luego el orden de `stores`, luego el del catálogo). Por canasta, si dos tiendas empatan en el total, va primero la que tiene más renglones recomendados (renglones cuya tienda recomendada es ella), luego la de menos precios viejos, luego la que va antes en `stores` (en la API, la de menor id). Una tienda con `missing_count > 0` se lista pero **no puede ser la más barata** (su total no cubre la canasta).

## 8. Membresía de Sam's (RN-13, K7; S5-09)
- `basket_savings_cents` = HALF_UP(Σ por renglón de max(0, costo en la mejor tienda sin membresía − costo en Sam's)), solo renglones donde Sam's es más barato y con las cantidades de la canasta (RN-13).
- `annual_fee_cents` = `store.membership_annual_fee_cents`; ausente → `verdict = missing_fee` y la salida dice «sin costo capturado» (PROVISIONAL de S5-09); nunca un costo inventado.
- `break_even_purchases` = costo anual ÷ ahorro por compra (S5-09: 1,000.00 ÷ 125.00 = **8**, ILUSTRATIVO). Si la división no es exacta, PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 4b), opción A: **hacia arriba** (8.3 compras → 9: con 8 no se cubre). Ahorro 0 → ausente y `not_worth_it`.
- `annual_savings_cents` y `verdict`: la frecuencia con que se anualiza no está decidida (el bloqueo `C1-campos-calculados` fijó «compras al año» pero no la frecuencia). PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 4a), opción B: compras al año = compras de despensa registradas en cadenas en los últimos 365 días; `annual_savings_cents` = ahorro por compra × ese número; `worth_it` si `annual_savings_cents ≥ annual_fee_cents`, si no `not_worth_it`. Ejemplo ILUSTRATIVO de RN-13: 138.50 × 6 = 831.00.

## 9. Mercados con día de visita (RN-16; S5-11)
- `NextVisit(store, from)` = la primera fecha **≥ `from`** cuyo día ISO es `visit_weekday` (si hoy es día de visita, hoy cuenta: supuesto S-10). RN-16: tianguis de martes (2) desde el sábado 10-oct-2026 → **martes 13-oct-2026**.
- Un mercado compite solo por los grupos que declara; sus precios entran por las mismas fuentes (normalmente `manual` o `ticket_manual`) y la misma frescura.

## 10. Insumo de la alerta de precios viejos (S5-12)
`StalePricesForPlan(products, assignedStores, groceryDate)`: para cada producto de los ingredientes del plan, la cotización de `QuoteFor` (§14.1) **evaluada en la fecha de despensa** (no en «hoy»). Si el producto trae tienda en `assignedStores` (la que eligió la persona en la lista, `shopping_list_item.store_id`), se cotiza en esa tienda; si no, sin tienda (la recomendada). En los dos casos la presentación sale de la regla de §14.1. La función devuelve `{product_id, store_id, presentation_id, observed_date, age_days, weight}` de las que estén viejas en esa fecha y, aparte, los productos sin precio que compita. La fecha de despensa la da S7 (decisión F).

## 11. Cargador PROFECO (`precios/profeco`; S5-05)
- **Archivos:** cortes de datos abiertos (2024 y 2025 en CSV; 2026 en RAR), cada descarga en su carpeta propia `datos/semilla/profeco/raw/<corte>/` (sin scripts dentro). El RAR se extrae con `unrar x` o `bsdtar -xf` a `datos/tmp/`; nada del archivo se ejecuta. Límite por archivo: `price.profeco_import_max_gb` (1 GB, K6).
- **Codificación:** UTF-8 con o sin BOM (así viene `QQP_diccionario.csv`) o Windows-1252 (así viene `QQP_metadatos.csv`); si los bytes no son UTF-8 válido se decodifica con una tabla Windows-1252 propia (biblioteca estándar).
- **Columnas** (las 15 de `QQP_diccionario.csv`, en cualquier orden, nombres sin distinguir mayúsculas): `producto`, `presentacion`, `marca`, `categoria`, `catalogo`, `precio` (Número 18,2 → centavos exactos desde el texto), `fecha_registro` (Datetime 8: el diccionario no da el formato; el cargador acepta los formatos que encuentre, los declara en el reporte de cobertura y rechaza con motivo los que no reconozca), `cadena_comercial`, `giro`, `nombre_comercial`, `direccion`, `estado`, `municipio`, `latitud`, `longitud` (texto decimal, sin flotantes).
- **Conteo:** renglones leídos (con lector CSV, para no contar saltos de línea dentro de comillas) = guardados + rechazados con motivo; cero pérdidas sin explicar. `MAX(fecha_registro)` debe ser la de la fuente (hasta 31-ago-2026).
- **Almacenamiento:** PROVISIONAL (bloqueo `S5-profeco-almacen`), opción A: todos los registros a una base local aparte `datos/profeco/qqp.sqlite`, tabla `profeco_record` (las 15 columnas + `archivo`, `renglon`, `precio_cents`, `fecha` ISO), fuera del esquema de la app e ignorada por git; ahí cuenta S5-05. Solo los registros **filtrados y ligados** pasan a `price` con `source = profeco` (que nunca compite) y `source_ref = "QQP <corte> <archivo>:<renglón>"`.
- **Filtro:** `cadena_comercial` ∈ nombres PROFECO de las cadenas (`store.profeco_name`) y `municipio` ∈ municipios de las sucursales (`store.municipality`). La ortografía real de las cadenas en PROFECO se descubre en los datos y se lista para que Axel confirme; no se supone.
- **Liga con el catálogo:** `datos/semilla/profeco/vinculos.csv` (`producto,presentacion,marca,presentacion_id,estado,confirmo,fecha_confirmacion`), con confirmación humana como los empates de nutrición (`S4-interfaz.md` §9.1): `estado` = `propuesto` o `confirmado`; `confirmo` (código de la persona, `axel` o `ximena`) y `fecha_confirmacion` (ISO) vacíos si `propuesto` y obligatorios si `confirmado`. En modo **producción** solo cuentan las ligas `confirmado`: un registro cuya liga está `propuesto` va a `skipped_count` con motivo «liga sin confirmar»; en modo **revisión** (pruebas, fixtures y medición de WF-1) se usan las dos y se cuentan aparte. Quién confirma y cuándo: PROVISIONAL (bloqueo `S4-confirmacion-empates`), opción A: Axel en la puerta 2. Registro filtrado sin liga → `skipped_count` (C1: «sin presentación ligada»); con liga que cuenta → `loaded_count`. Tienda del precio: la sucursal cuyo `profeco_name` = `nombre_comercial`; si no hay, la cadena. Varios registros de la misma presentación, tienda y día: se conserva el primero del archivo y los demás se cuentan como duplicados (PROVISIONAL, mismo bloqueo, pregunta 3).
- **Salida** `ImportResult{Read, Rejected, FilteredOut, Skipped, Loaded, Duplicates, Coverage[]}`; C1 `ProfecoImportResult` expone `loaded_count`, `skipped_count` y `coverage` (por cadena y municipio, `price_count`).
- **Reporte** `docs/estado/S5-cobertura-profeco.md` (lo escribe el agente PROFECO): cita `datos/semilla/profeco/QQP_diccionario.csv`; conteo por archivo contra `wc -l`; tabla por `cadena_comercial` y por municipio; fechas mínima y máxima; y para **Bodega Aurrera Vallejo, Walmart Acueducto de Guadalupe y Sam's Club** dice con todas sus letras si hay o no registros, con la regla de empate (texto en `nombre_comercial` o `direccion`) escrita **antes** de contar.

## 12. OCR local (`precios/ocr`; K24, S5-06)
- `Recognize(imagen)`: verifica una vez `tesseract --version` = 5.5.3 y `spa` en `--list-langs` (si no, `ErrOCRUnavailable`); corre `tesseract <imagen> stdout -l spa --psm 6 tsv` sin shell, con tiempo máximo PROVISIONAL (propuesta K31/K44: `ocr.timeout_seconds` = 30, clave que hoy no existe en `app_setting`; solicitud a S1).
- **Salida:** renglones `{descripcion string, cantidad_milli *int64, importe_cents *int64, confianza int}`: palabras agrupadas por bloque, párrafo y línea del TSV; `importe` = último número con dos decimales del renglón (`1,234.50`, `$45.00`) en centavos desde el texto; `cantidad` = número inicial (entero o con decimales y unidad, `0.750 KG`) en milésimas; `descripcion` = el resto; `confianza` = promedio HALF_UP de la confianza de las palabras (0 a 100). Renglones sin importe se devuelven con `importe` nulo para que la persona los descarte.
- **Siempre a confirmación humana:** la salida es un borrador (`attachment.ocr_draft_json`, K3, «Guardar para después»). `PricesFromOCR(lineas, confirmado)` devuelve `ErrNotConfirmed` si `confirmado = false`: un OCR sin confirmar no crea ningún precio (S5-06). Confirmado → renglones de gasto y precios `ticket_ocr`.
- **Medición:** `datos/semilla/tickets/` guarda **solo** las 5 fotos, `01` a `05`, en cualquier formato que Tesseract lea (`.jpg`, `.jpeg` o `.png`; K24 no fija el formato), así `ls datos/semilla/tickets/ | wc -l` = 5 (S5-06). La transcripción a mano de cada ticket va **fuera** de esa carpeta, en `datos/semilla/tickets-esperado/NN.csv` (`descripcion,cantidad,importe`; mismo `NN` que la foto), escrita **antes** de medir; no va en una subcarpeta `tickets/esperado/` porque `ls` la contaría y daría 6. Regla de descripción (propuesta para congelar en el reporte del agente OCR antes de medir): igualdad tras normalizar (mayúsculas, sin acentos, espacios colapsados, sin puntuación final) contra el texto **impreso** en el ticket; cantidad e importe exactos. Hoy la carpeta no existe: **«sin medir: faltan las fotos de Axel»**.

## 13. Lectura asistida (`precios/asistida`; S5-07)
`CanEnable(cadena, terminos)`: exige `price.assisted_reading_enabled = true`, `store.assisted_reading_enabled = 1` y, en `docs/estado/S5-terminos.md`, una sección por sitio (Walmart, Bodega Aurrera, Sam's Club) con las líneas `Veredicto:`, `URL:` y `Fecha de lectura:`; si falta cualquiera, `ErrNoTermsVerdict` y el interruptor no se activa. Aun activada, `site_assisted` sigue siendo referencia mientras rija K37 para esa cadena.

## 14. Funciones que consumen S4, S7 y la API

| Función | Entradas | Salida (C1) |
|---|---|---|
| `QuoteFor(product, store?, presentation?, date)` | catálogo, precios, tiendas, `Settings` | `Quote` o «sin precio» (`PriceSource` de S4), con la regla de §14.1. Sin `store` ni `presentation`: el par (tienda, presentación) recomendado, que S4 recibe con `recommended = true`. Con tienda, con presentación o con las dos: el mismo criterio dentro de esa elección (`S4-interfaz.md` §8.1 y §12.1) |
| `BaseUnitPrice(price, presentation)` | | `base_unit_price_cents` (`PriceAt`) |
| `Freshness(observed, today, settings)` | | `Weight` (`is_stale`) |
| `CompareCarts(base, others[], date)` | la compra registrada y los carritos u otras compras (`Cart`) | comparativa de §7.2 (canon), con `piece_count` e `identical_pieces` |
| `CompareExpense(expense, date)` | la compra `expense_id` y los precios vigentes | `BasketComparison` (`basket = expense`); PROVISIONAL (bloqueo `S5-carrito-comparativa`) |
| `CompareBasket(lista, date)` | lista de S4, con la cantidad a comprar (§7.3) | `BasketComparison` (`basket = shopping_list`): total de cada tienda con el menor costo de paquetes para la cantidad a comprar, la misma regla en todas las tiendas, y recomendación por renglón = `QuoteFor` sin elección (§7.3). Los campos calculados de `ShoppingListItem` (`recommended_store_id`, `recommended_presentation_id`, `estimated_price_cents`, `is_stale_price`) los llena la API con la salida de S4 (`S4-interfaz.md` §12.2, traducción a C1), que usa las mismas cotizaciones; si la persona no eligió nada, coinciden con esta recomendación |
| `Membership(comparison, store, historial)` | | `membership` |
| `NextVisit(store, from)` | | fecha |
| `StalePricesForPlan(...)` | | lista para S7 |
| `Summary(date)` | | `PriceSummary` |

### 14.1 Regla de `QuoteFor`: fecha, presentación y desempate (S5-03, S5-12; S4-03, S4-06)
- **Fecha.** La API llama `QuoteFor` con `date` = «hoy» en America/Mexico_City (`today`, C-05) para el costo de recetas, del plan y de la lista, también cuando la semana es futura: la edad y el peso de un precio se miden contra «hoy» (§5, S5-04) y no hay precios con fecha futura (`ErrFutureDate`). La única excepción es `StalePricesForPlan`, que evalúa en la fecha de despensa (§10).
- **Candidatos.** Cada par (tienda, presentación del producto) con precio vigente que compite en `date` (§6; `Competes`, §5), en las tiendas que compiten por el producto según §7.4: cadenas y sucursales, mercados solo por los grupos que declaran y Sam's como si tuvieran membresía. Con tienda elegida, solo los pares de esa tienda; con presentación elegida, solo los de esa presentación; con las dos, solo ese par. Sin candidatos: «sin precio».
- **Orden (gana el primero).** (1) Menor precio por unidad base ajustado = (`price_cents` / `size_milli`) ÷ peso, como fracción exacta (§7.1). (2) Menor precio del paquete (`price_cents`). (3) La tienda que va antes en `stores`. (4) La presentación que va antes en el catálogo (orden de `catalog.presentations` o de las presentaciones del producto en `productos.json`). La API entrega tiendas y presentaciones por `id` ascendente, así que (3) y (4) equivalen a «menor id». Los pasos 2 a 4 son el desempate PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 2, opción A). No dependen de la lista de compras ni de otros productos: `QuoteFor` no recibe la lista.
- **`recommended`** = `true` solo en la llamada sin tienda ni presentación. La necesidad no cuenta: la presentación se elige por unidad base (S5-03: «la tienda más barata por ingrediente (en unidad base)»). El menor costo de paquetes para la cantidad a comprar (§7.3) solo arma el total de cada tienda en la comparativa por canasta, con la misma regla en **todas** las tiendas, también la recomendada; no cambia la recomendación del renglón.
- **Ejemplo** ILUSTRATIVO (caso `s5-cotizacion-recomendada`). Hoy es 2026-10-12, los precios son del 2026-10-10 (peso 1) y `stores` = [`walmart`, `bodega`]. Arroz: Walmart tiene 1 kg a 2800 y 750 g a 2200; Bodega, 750 g a 2200 y 500 g a 1500.
  - Sin elección: Walmart 1 kg (2800 por kg, contra 2933 y 3000), con `recommended = true`.
  - Con `store` = `bodega`: 750 g a 2200 (2933 por kg contra 3000), aunque la de 500 g cueste menos para una necesidad de 500 g.
  - Con `presentation` = 750 g: Walmart y Bodega empatan en (1) y en (2); gana Walmart por (3).
  - Frijol: Walmart 1 kg a 3000 y Bodega 500 g a 1500 empatan en (1), los dos a 3000 por kg; gana Bodega por (2).

## 15. Fixtures `oraculo/fixtures/s5/` y CLI `s5 dump`
Formato común `oraculo/fixtures/FORMATO.md` 1.0 (`fixture/1`), igual que `S4-interfaz.md` §12, con `motor` = `s5` y casos `s5-<tema>`; los ILUSTRATIVOS llevan `fuente` `manual`; los que tocan `S5-reglas-comparativa` o `S5-profeco-almacen` llevan `provisional`. **Funciones consultables:** `base_unit_price` (§4), `freshness` (§5), `current_price` (§6), `quote_for` (§14.1), `compare_carts` (§7.2), `compare_basket` (§7.3, §7.4), `membership` (§8), `next_visit` (§9), `stale_prices_for_plan` (§10), `profeco_import` (§11, sobre CSV de `oraculo/fixtures/s5/_datos/`), `ocr_parse_tsv` (§12, sobre un TSV de `oraculo/fixtures/s5/_datos/`), `assisted_can_enable` (§13). `Input` en §15.1; argumentos y salida de cada una en §15.2.

| Caso | Fuente | Qué fija |
|---|---|---|
| `s5-comparativa-canon` | `data.json` (lo lee directo por `report_file`, sin copias) + `productos.json` | compra de Bodega contra carrito de Walmart (§3): K-01 a K-08 (con `piece_count` e `identical_pieces`), ganadores 22/8/3, la descomposición de +138.50 y +100.50 de §7.2 |
| `s5-unidad-sustitutos` | `data.json: sustitutos` | 2933/2800, 12692/16406, 3375/2833 y el más barato |
| `s5-frescura-fronteras` | K40, RN-12 | la tabla de §5 y la recta de 1/61 |
| `s5-fuentes-compiten` | K22, K37 | 5 fuentes × 1 precio: `manual`, `ticket_manual`, `ticket_ocr` misma edad pesan igual y compiten; `site_assisted` y `profeco` más baratos nunca ganan; promoción vencida no entra |
| `s5-mejor-tienda-canasta` | `data.json` + ILUSTRATIVO | canasta del 10-oct (Bodega) y 2 canastas con empate |
| `s5-cotizacion-recomendada` | ILUSTRATIVO; `provisional` (`S5-reglas-comparativa`, pregunta 2) | el ejemplo de §14.1: `quote_for` sin elección, con `store` y con `presentation` (empate que resuelve el orden de `stores`), y el empate del frijol (gana el menor precio de paquete). `compare_basket` con arroz (`needed_milli` 500000, `quantity_milli` `null`, `recommended_presentation` 1 kg): `recommended_store` Walmart, `recommended_presentation` 1 kg y `estimated_price_cents` 2800; Walmart cotiza 750 g a 2200 y Bodega 500 g a 1500, los dos sustitutos (§7.3: una sola regla en todas las tiendas); totales Bodega 1500 y Walmart 2200 (+700); `cheapest_store` = Bodega (por ingrediente manda la unidad base; por canasta, el total) |
| `s5-canasta-regla-unica` | ILUSTRATIVO (contraejemplo de M-r4-1, `docs/estado/wf1/REVISION-contrato-S4S5-r4.md`) | hoy 2026-10-12, precios del 2026-10-10 (peso 1), `stores` = [`walmart`, `bodega`]. Arroz (`needed_milli` 500000, `quantity_milli` `null`, `recommended_presentation` 1 kg): Walmart 1 kg a 2800 y 500 g a 1450, Bodega 500 g a 1460 → recomendado Walmart 1 kg, `estimated_price_cents` 2800; Walmart cotiza 500 g a 1450 y Bodega 500 g a 1460, los dos sustitutos. Aceite de soya `manual` (`needed_milli` `null`, `quantity_milli` 1892000, `recommended_presentation` 946 ml): Walmart 3990 y Bodega 4000 por botella → recomendado Walmart, 2 paquetes, `estimated_price_cents` 7980, idéntico en las dos tiendas (Bodega 8000). Totales: Walmart **9430** y Bodega **9460** (+30); `cheapest_store` = Walmart. Si la tienda recomendada se cotizara con la presentación recomendada, Walmart saldría en 10780 y la canasta señalaría a Bodega |
| `s5-membresia-equilibrio` | S5-09, ILUSTRATIVO | costo 100000 y ahorro 12500 → 8; costo ausente → `missing_fee` |
| `s5-mercado-visita` | RN-16 | sábado 10-oct-2026 + martes → 2026-10-13; mercado sin el grupo no compite |
| `s5-precios-viejos` | S5-12, ILUSTRATIVO | despensa el 2026-11-09: precios del 10-oct ya viejos |
| `s5-lectura-asistida` | S5-07 | interruptor sin veredicto rechazado |
| `s5-profeco-muestra` | ILUSTRATIVO | CSV de 6 renglones (uno Windows-1252, una fecha ilegible, uno sin liga, un duplicado) → conteos |
| `s5-ocr-parseo` | ILUSTRATIVO | TSV de muestra → renglones |

**CLI** (FORMATO §4): `datos/tmp/motores s5 dump --fixtures oraculo/fixtures/s5 --out datos/tmp/motor-s5.json [--caso <caso>]` (mismo binario `backend/cmd/motores` que S2, S3 y S4); comparación `python3 -I oraculo/compara.py s5`. Sin red; `s5-ocr-parseo` usa el TSV de la fixture y no corre Tesseract; `s5-comparativa-canon` lee `datos/semilla/data.json` con la clave `report_file` del `Input` (ruta relativa a la raíz del repositorio, §15.1), sin copias.

### 15.1 `Input` de S5 (raíz de `entradas.datos`)
Mismas convenciones que `S4-interfaz.md` §12.1: claves en inglés y snake_case, `_cents` y `_milli`, fechas ISO, entidades por código (la API traduce a ids), `unknown_field` y `missing_field` en el motor y en el oráculo, rutas `*_file` relativas a la raíz del repositorio. Toda lista puede venir vacía.

| Clave | Tipo | Oblig. | Contenido | Origen en la API |
|---|---|---|---|---|
| `today` | fecha | sí | «Hoy» en America/Mexico_City (C-05) | API |
| `settings` | objeto | no | {`stale_after_days`, `zero_weight_after_days`, `assisted_reading_enabled`}; lo ausente toma el valor inicial (30, 90, `false`) | `app_setting price.*` (K44) |
| `catalog_file` / `catalog` | ruta / objeto | uno de los dos | como en `S4-interfaz.md` §12.1 | semilla, `product`, `presentation` |
| `stores` | lista | no | {`code`, `kind` (`chain` · `branch` · `market`), `parent` (o `null`), `name`, `profeco_name` (o `null`), `municipality` (o `null`), `visit_weekday` (1 a 7 o `null`), `requires_membership`, `has_membership`, `membership_fee_cents` (o `null`), `assisted_reading`, `groups[]`} | `store`, `store_product_group` |
| `prices` | lista | no | {`presentation`, `store`, `price_cents`, `observed_date`, `source` (los 5 de §3), `is_promotional`, `source_ref` (o `null`), `created_at` (instante RFC 3339)} | `price` |
| `carts` | lista | no | canastas de §7.2: {`code`, `store`, `kind` (`expense` · `cart`), `observed_date`, `source`, `lines[]` {`presentation`, `quantity_milli`, `unit_price_cents`}} | compra: `expense` + `expense_item`; carrito: solo motor (§3) |
| `report_file` | ruta | no | `datos/semilla/data.json`. El motor y el oráculo arman con él dos `carts` del 2026-10-10 con las fuentes de §3: `bodega` (`expense`: los 33 idénticos con `bodega_unit`, la sexta lata de atún y los 3 sustitutos de Bodega) y `walmart` (`cart`: los 33 idénticos con `walmart_unit` y los 3 sustitutos de Walmart), emparejando cada texto de `items_identicos[].producto` y `sustitutos[].bodega` / `.walmart` con el `name` exacto de una presentación del catálogo (`productos-formato.md` §2). No se combina con `carts` (`invalid_argument`) | semilla |
| `shopping_list` | lista | no | canasta por necesidad de §7.3: un renglón por artículo de la lista de S4: {`product`, `needed_milli` (o `null`: artículo `manual`), `quantity_milli` (o `null`: sin cantidad fijada por la persona), `recommended_presentation` (o `null`)}; la cantidad a comprar es `quantity_milli` si no es nula, y si no `needed_milli` (§7.3); si las dos son nulas, `invalid_argument`; `recommended_presentation` es la presentación que rige en la lista de S4 (`items[].presentation`, `S4-interfaz.md` §12.2) y solo decide si el renglón de cada tienda es idéntico o sustituto | lista de S4 (`items[].needed_milli`, `quantity_milli` y `presentation`, `S4-interfaz.md` §12.2) |
| `grocery_purchases` | lista | no | compras de despensa registradas en cadenas, para anualizar la membresía (§8): {`date`, `store`} | `expense` |
| `terms` | lista | no | veredictos ya leídos de `docs/estado/S5-terminos.md`: {`site` (`walmart` · `bodega_aurrera` · `sams_club`), `verdict`, `url`, `read_date`}; un campo vacío cuenta como faltante | archivo de términos |

Los CSV de PROFECO, el `vinculos.csv` de prueba y el TSV del OCR no van en el `Input`: los nombran los argumentos de `profeco_import` y `ocr_parse_tsv` con rutas **relativas a la raíz del repositorio**, la misma regla que los `*_file` de esta tabla, y siempre dentro de `oraculo/fixtures/s5/_datos/` (por ejemplo, `oraculo/fixtures/s5/_datos/profeco-muestra.csv`). El evaluador de FORMATO §4.1 (`func(input json.RawMessage, q Query)`) no recibe la carpeta de fixtures y `main.go` no cambia: el motor y el oráculo corren desde la raíz del repositorio (CLI de §15; FORMATO §4: rutas relativas al directorio actual) y abren la ruta tal como viene. Una ruta absoluta, con `..` o fuera de `oraculo/fixtures/s5/_datos/` es `invalid_argument` en el motor y en el oráculo. No son `$ref`: `main.go` solo resuelve el `$ref` de `entradas.datos`.

### 15.2 Funciones consultables: argumentos y salida
Salidas con claves en inglés y entidades por su código con la clave sin `_id` (`store`, `product`, `presentation`), como en `S4-interfaz.md` §12.2. Pesos y precios exactos como fracciones `{"num": n, "den": d}`. Los errores de §17 salen como `{"error": "<código>"}` en snake_case sin `Err` (`ErrNotConfirmed` → `not_confirmed`), además de `unknown_field`, `missing_field`, `invalid_date` e `invalid_argument`.

| `funcion` | `argumentos` | Salida |
|---|---|---|
| `base_unit_price` | `presentation`, `price_cents` | {`base_unit_price_cents`, `per` (`kg` · `l` · `piece`), `exact` (fracción: centavos por 1 kg, 1 l o 1 pieza, sin redondear)} (§4) |
| `freshness` | `observed_date`, `source` | {`age_days`, `is_stale`, `weight` (fracción), `weight_display` (texto con 4 decimales HALF_UP), `competes`} (§5) |
| `current_price` | `presentation`, `store` | {`price` ({`price_cents`, `observed_date`, `source`} o `null`), `references[]` {`source`, `price_cents`, `observed_date`}, `promotion` ({`price_cents`, `observed_date`} o `null`)} (§6) |
| `quote_for` | `product`; `store` y `presentation` opcionales; la fecha es `today` (§14.1) | {`quote` ({`product`, `store`, `presentation`, `price_cents`, `size_milli`, `observed_date`, `is_stale`, `recommended`}, las mismas claves que `quotes` de `S4-interfaz.md` §12.1, o `null` = «sin precio»)} (§14.1) |
| `compare_carts` | `base` (código de `carts`), `others[]` (códigos) | {`stores[]` {`store`, `total_cents`, `identical_count`, `identical_total_cents`, `substitute_count`, `piece_count`, `missing_count`, `stale_count`, `difference_vs_cheapest_cents`, `assumes_membership`}, `identical_pieces`, `winners[]` {`store`, `count`}, `ties`, `rows[]` {`kind` (`identical` · `substitute` · `only_one`), `product`, `lines[]` {`store`, `presentation`, `quantity_milli`, `total_cents`, `base_unit_price_cents`}, `winner_store`}} (§7.2) |
| `compare_basket` | ninguno (usa `shopping_list`) | {`stores[]` (mismas claves que en `compare_carts`), `items[]` {`product`, `quantity_milli` (la cantidad a comprar, §7.3), `recommended_store`, `recommended_presentation`, `estimated_price_cents`, `is_stale_price`, `prices[]` {`store`, `presentation`, `price_cents`, `base_unit_price_cents`, `observed_date`, `is_substitute`}}, `cheapest_store`} (§7.3, §7.4; el total de cada tienda usa una sola regla y la recomendación por renglón es `QuoteFor` sin elección, §14.1) |
| `membership` | `store`; `savings_per_purchase_cents` opcional (si viene, no se calcula con `compare_basket`) | {`store`, `assumes_membership`, `basket_savings_cents`, `annual_fee_cents`, `break_even_purchases`, `purchases_per_year`, `annual_savings_cents`, `verdict` (`worth_it` · `not_worth_it` · `missing_fee`)} (§8) |
| `next_visit` | `store`, `from` | {`date`} (§9) |
| `stale_prices_for_plan` | `products[]`, `assigned_stores` (objeto producto → tienda), `grocery_date` | {`stale[]` {`product`, `store`, `presentation`, `observed_date`, `age_days`, `weight` (fracción)}, `without_price[]`} (§10) |
| `profeco_import` | `files[]` (rutas de CSV desde la raíz del repositorio, dentro de `oraculo/fixtures/s5/_datos/`, §15.1), `links_file` (ruta con la misma regla; CSV con las columnas de `vinculos.csv`), `mode` (`production` · `review`) | {`read`, `rejected`, `filtered_out`, `skipped`, `loaded`, `duplicates`, `coverage[]` {`chain`, `municipality`, `price_count`}, `rejections[]` {`file`, `row`, `reason`}} (§11; filtra con `stores`) |
| `ocr_parse_tsv` | `tsv_file` (ruta desde la raíz del repositorio, dentro de `oraculo/fixtures/s5/_datos/`, §15.1) | {`lines[]` {`description`, `quantity_milli`, `amount_cents`, `confidence`}} (§12: son `descripcion`, `cantidad_milli`, `importe_cents` y `confianza` con clave en inglés) |
| `assisted_can_enable` | `store` (cadena) | {`can_enable`} o `{"error": "no_terms_verdict"}` (§13) |

### 15.3 Oráculo y verificación independiente
- `oraculo/s5_run_fixtures.py` (Python 3.14 con `-I`, solo biblioteca estándar: `json`, `csv`, `fractions`, `decimal`, `datetime`, `codecs`) implementa las funciones de §15.2 desde este contrato, **sin leer el código Go**, y escribe `oraculo/salidas/s5.json` (FORMATO §4) con **todas** las consultas de **todos** los casos de `oraculo/fixtures/s5/`; `python3 -I oraculo/compara.py s5` lo enfrenta con el motor y con `esperado` (FORMATO §5). Lo escribe un agente distinto del que implementa `internal/precios`.
- Independencia por caso: `s5-comparativa-canon` y `s5-unidad-sustitutos` los calcula el oráculo desde `data.json` y deben coincidir con los totales y notas del reporte (K-01 a K-08, 22 / 8 / 3, 2933 / 2800, 12692 / 16406, 3375 / 2833); `s5-frescura-fronteras` contra K40 y la línea de Python de S5-04; los ILUSTRATIVOS (`fuente` `manual`) escriben en `descripcion` su cálculo a mano.
- `s5-ocr-parseo` no mide S5-06: la medición con las 5 fotos de Axel la hace el script del agente OCR contra `datos/semilla/tickets-esperado/` (§12).

## 16. Pruebas (nombres que casan con los `-run` de los criterios)

| Criterio | Pruebas Go |
|---|---|
| S5-01 | `TestComparativaCanon` (lee `data.json`: compra de Bodega contra carrito de Walmart), `TestComparativaGanadores`, `TestComparativaPiezas` (61 / 60 / 55 y 33 / 3) |
| S5-02 | `TestUnidadSustitutos`, `TestUnidadRedondeo` |
| S5-03 | `TestMejorTiendaCanasta`, `TestMejorTiendaEmpate`, `TestMejorTiendaIngrediente`, `TestMejorTiendaCotizacion` (`quote_for` y, en `compare_basket`, la recomendación por renglón frente a la línea de la tienda recomendada, §14.1 y §7.3), `TestMejorTiendaCanastaReglaUnica` (§7.3: una sola regla para el total de cada tienda y cantidad a comprar con `quantity_milli`; caso `s5-canasta-regla-unica`) |
| S5-04 | `TestFrescuraFronteras`, `TestFrescuraRecta`, `TestFuenteMismoPeso`, `TestFuenteReferencia` |
| S5-05 | `TestProfecoConteo`, `TestProfecoCodificacion`, `TestProfecoFiltroYLiga`, `TestProfecoLigaSinConfirmar` |
| S5-06 | `TestOCRParseo`, `TestOCRSinConfirmar`, `TestOCRVersion` (con ejecutor simulado); la medición con las 5 fotos la hace el script del agente OCR contra `datos/semilla/tickets-esperado/` (§12) |
| S5-07 | `TestLecturaAsistidaSinVeredicto` |
| S5-09 | `TestMembresiaEquilibrio`, `TestMembresiaSinCosto` |
| S5-10 | `TestFuenteTrazabilidad`, `TestFuentePromocionVencida` |
| S5-11 | `TestMercadoVisita`, `TestMercadoGrupo` |
| S5-12 | `TestPreciosViejosDespensa` |
| S5-13 | `go vet`, `gofmt -l`, cobertura ≥ 85 %, `grep float` sin líneas de dinero |

## 17. Errores
`ErrUnknownPresentation`, `ErrUnknownStore`, `ErrFutureDate`, `ErrDuplicatePrice`, `ErrNotConfirmed`, `ErrOCRUnavailable`, `ErrNoTermsVerdict`, `ErrProfecoHeader`, `ErrFileTooLarge`, `ErrOverflow`. Estados que no son error: «sin precio», «sin costo capturado», «sin medir».

## 18. Verificación · criterio → sección

| Criterio | Sección |
|---|---|
| S5-01 Comparativa reproducida | §7.2, §3, §15 (`report_file`, `compare_carts`) |
| S5-02 Equivalencia por unidad base | §4 |
| S5-03 Mejor tienda por ingrediente y canasta | §7.3 (una sola regla para el total de cada tienda; cantidad a comprar), §7.4, §14.1 (`QuoteFor`) |
| S5-04 Frescura y fuentes que compiten | §5, §3, §7.1 |
| S5-05 PROFECO y cobertura | §11 (ligas confirmadas), §15.1 (rutas de `_datos/` desde la raíz) |
| S5-06 OCR ≥ 80 % | §12 (carpetas de fotos y de transcripciones), §15.1 (`tsv_file` desde la raíz) |
| S5-07 Lectura asistida desactivada | §13 |
| S5-08 Sin scraping ni red | §1, §13 |
| S5-09 Membresía de Sam's | §8 |
| S5-10 Fuentes y trazabilidad | §3, §6 |
| S5-11 Mercados | §9, §7.4 |
| S5-12 Precios viejos para S7 | §10, §14.1 |
| S5-13 Calidad Go | §1, §16 |
