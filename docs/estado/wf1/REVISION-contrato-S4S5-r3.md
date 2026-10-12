# Revisión del contrato S4/S5 · ronda 3 de 3 (WF-1, 11-oct-2026)

Revisor independiente: `revisor-contrato-S4S5-r3`. No escribí ninguno de los artefactos revisados y no corregí nada. Evalué solo contra `docs/criterios/S4.md` 1.2, `docs/criterios/S5.md` 1.2 y decisiones §D, §E y §K. Artefactos revisados: `docs/contratos/motores/S4-interfaz.md` 1.2, `S5-interfaz.md` 1.2 y `productos-formato.md` 1.1. También leí el reporte `docs/estado/wf1/contrato-S4S5.md` y la corrección `correccion-contrato-S4S5-r2.md`, pero no aprobé nada por lo que dicen: rehíce cada cifra con mi propio script.

## Veredicto
**NO APROBADO.** No hay hallazgos altos. Quedan **2 medios nuevos**: M-r3-1 y M-r3-2, los dos en la frontera S4 ↔ S5 ↔ C1 de la lista de compras. Los 4 medios de r2 están **cerrados**. Todas las cifras canon coinciden con `data.json`, el oráculo y §K.

## Verificación propia
Script `datos/tmp/rev-s4s5-r3/verif.py` (`python3 -I`, solo biblioteca estándar; `Decimal` con `parse_float` y `Fraction`; lee `datos/semilla/data.json` y los `.md` directamente). Resultado: **63 comprobaciones bien** y 1 aviso, que es un falso positivo de mi propio filtro: la línea S4-12 de §15 cita `S4-asignacion-mercado` sin la palabra «bloqueo».

- **(c) Canon de la comparativa.** Lo calculé desde los renglones, no desde `totales`. 0 renglones con unitario × cantidad ≠ total.
  - Idénticos: Bodega 277350 y Walmart 287400, diferencia **+10050 (+100.50)**.
  - Totales: Bodega 309450 (277350 + 30000 + 2100) y Walmart 323300 (287400 + 35900), diferencia **+13850 (+138.50)** = 10050 + 5900 − 2100.
  - Los totales coinciden con `totales.*`. Piezas 61 / 60 / 55; renglones 33 / 3; ganadores 22 / 8 / 3, igual que `resultados_por_producto`.
- **Unidad base**, con el contenido de las presentaciones del catálogo §5.1: 2933 / 2800, 12692 / 16406 y 3375 / 2833.
- **(e) Catálogo.**
  - 36 productos y 39 presentaciones, todas con el prefijo de su producto.
  - Los 39 nombres de `data.json` existen tal cual y son distintos: 33 idénticos, 3 sustitutos de Bodega y 3 de Walmart.
  - La tabla §5.2 es igual, renglón por renglón, a `cuentas.{axel,karen,conjunto}.lineas`: cuenta, texto, piezas e importe en el mismo orden. Suma **39 renglones, 61 piezas y 3,094.50**.
  - El importe de cada renglón = piezas × `bodega_unit` (o `total ÷ cant_b` en los sustitutos), la unidad coincide y las piezas por presentación = `items_identicos[].cantidad` (atún 5 + 1).
- **S4 §8.3:** los renglones 21, 16, 15 y 0 son espagueti, atún, elote y aceite. Resultado 4425, 2213 y 8850.
- **Frescura (K40):**
  - 29 días: peso 1.
  - 30 días: 60/61, que se muestra 0.9836. 31 días: 59/61.
  - 60 días: 30/61 = 0.4918. 90 días: 0. 91 días: no compite.
  - Fechas 2026-11-08, 2026-11-09 y 2027-01-08.
- **Salud (K23, K36):**
  - MSJ en milésimas: 1748750, 1330250 y 1648750.
  - Metas 2099, 1596, 1799, 1846, 2556 y 2056; exactas 2098.5 y 1596.3.
  - Factores exactos 6/5, 11/8, 31/20, 69/40 y 19/10 = 1.2 a 1.9, igual que RN-17.
  - IMC en décimas: 185, 250, 300 y 237; 73990 g da 185 (18.4975 exacto).
- **RN-15 y RN-16:** 2838000 − 1351429 = 1486571; martes siguiente al 10-oct = 2026-10-13. La fórmula de kcal (250 g × 350 kcal por 100 g = 875000 milésimas) da lo esperado.
- **Estructura:**
  - S4 §15 cubre S4-01 a S4-15 y S5 §18 cubre S5-01 a S5-13.
  - 2 bloques JSON válidos.
  - Las 17 expresiones `-run` de los criterios encuentran pruebas en §13 o §16, y solo de su tema: 0 mezclas.
  - Todo bloqueo citado existe en `docs/estado/bloqueos/`. La excepción es `S4-meta-ajustes`, que es hipotético según S4-07 1.2.
  - Las claves `price.stale_after_days`, `price.zero_weight_after_days`, `price.assisted_reading_enabled`, `price.profeco_import_max_gb` y `health.goal_adjust_*` existen en 00001 y 00003.

## Estado de los hallazgos de la ronda 2
| Hallazgo r2 | Estado | Evidencia propia |
|---|---|---|
| M-r2-1 (rutas de `profeco_import` y `ocr_parse_tsv`) | **cerrado** | S5 §15.1 y §15.2: las rutas son relativas a la raíz del repositorio, dentro de `oraculo/fixtures/s5/_datos/`. Una ruta absoluta, con `..` o fuera de esa carpeta es `invalid_argument`. Explícitamente no son `$ref`. El CLI de S4 §12 y S5 §15 corre desde la raíz, como permite FORMATO §4 («rutas relativas al directorio actual»), y el `Evaluator` no necesita la carpeta. `backend/cmd/motores/` aún no existe, así que nada lo contradice. Queda el bajo B-r3-2 |
| M-r2-2 (cotización recomendada en `quotes`) | **cerrado** en S4 | S4 §8.1, §12.1 y §6.2 paso 5: `recommended` con exactamente una por producto y a lo más una por (`product`, `store`, `presentation`), y una sola regla de cotización que rige para `recipe_detail`, `week_plan` y `shopping_list`. Fixture `s4-costo-cotizacion-elegida`. El motor S4 queda determinista con las cotizaciones dadas. Del lado de S5 abro M-r3-2: cómo arma `QuoteFor` esa cotización |
| M-r2-3 (cotejo de kcal y modo de carga) | **cerrado** | S4 §3.2: el generador escribe los derivados en modo `review`; el cotejo es solo en `review`; `production` carga las 60 sin cotejo y deja las recetas incompletas con `missing_nutrition`. §13 trae `TestRecetarioProduccion` y `TestKcalCotejoRevision`; §14 dice «`ErrKcalMismatch` solo en modo `review`». Coherente con §9.1. Queda el bajo B-r3-1 |
| M-r2-4 (artículos de vertical mercado) | **cerrado** dentro de S4, como PROVISIONAL | S4 §6.5 pasos 1 a 3, `groups[]`, §11 y §12.1 (`next_visit`). El bloqueo `S4-asignacion-mercado` existe con los 5 apartados del protocolo §2; la fixture `s4-canal-mercados` lleva `provisional`. La opción A contradice C1 v1.2: abro M-r3-1 |

## Criterios
| Criterio | Sev. | Estado | Sección | Evidencia propia |
|---|---|---|---|---|
| S4-01 Recetario | alta | cubierto (B-r3-1) | §3 | llave `recetas`, esquema `datos/semilla/recetario.schema.json`, campos y validación cruzada; las 60 cargan en `production` |
| S4-02 Escalado | alta | cubierto | §4 | fracciones exactas; 250000 → 375000, 7000 → 10500, 1000 → 1500 recalculados |
| S4-03 Lista | alta | cubierto, con observación (M-r3-2) | §6 | Σ − existencia + reposición, techo a paquetes, sobrante a existencia; 4 fixtures |
| S4-04 Inventario | alta | cubierto (B4 r1 abierto) | §7 | 1486571 recalculado; `out` prevalece |
| S4-05 Plan | alta | cubierto | §5 | «fuera», colación, Karen rechazada, K38, K43 (21) |
| S4-06 Costo | alta | cubierto, con observación (M-r3-2) | §8 | 4425 / 2213 / 8850 recalculados; «sin precio», nunca 0 |
| S4-07 Salud | alta | cubierto | §10 | MSJ, metas, IMC y factores recalculados; K36 por persona y casa; §10.4 confirma |
| S4-08 kcal | alta | cubierto | §9 | fórmula en milésimas verificada; `Coverage` con dos conteos |
| S4-09 Cinco a mano | alta | cubierto (lo hace el revisor de la implementación) | §12 | sin recetario todavía; el contrato no permite que el autor lo sustituya |
| S4-10 Licencias | media | cubierto (B1 r1 abierto) | §9.1, §1 | `PROCEDENCIA.md`, 3 fuentes |
| S4-11 Español | alta | cubierto | §3.2, §9.1 | `production` no carga lo `propuesto`; `nombre_fuente` solo como procedencia |
| S4-12 Canal | media | cubierto, con observación (M-r3-1) | §6.5, §11 | `ErrProductWithoutChannel`; mercados PROVISIONAL |
| S4-13 Editables | media | cubierto (B-r2-4 abierto) | §3.3 | funciones puras y consulta SQL |
| S4-14 IA apagada | media | cubierto (B2 r1 abierto) | §1, §13 | `TestSinIA`; no cita `ai.*` |
| S4-15 Calidad Go | alta | cubierto | §1, §13 | sin base en `internal/recetas`; sin flotantes; 85 % |
| S5-01 Comparativa | alta | cubierto | §7.2, §3, §15 | +138.50 y +100.50 recalculados (arriba) |
| S5-02 Unidad base | alta | cubierto | §4 | 2933/2800, 12692/16406 y 3375/2833 recalculados con el catálogo |
| S5-03 Mejor tienda | alta | cubierto, con observación (M-r3-2) | §7.3, §7.4 | Bodega la más barata; empate PROVISIONAL |
| S5-04 Frescura y fuentes | alta | cubierto | §5, §3 | pesos y fechas recalculados; 3 compiten y 2 son referencia |
| S5-05 PROFECO | alta | cubierto (B6 r1 abierto) | §11, §15.1 | leídos = guardados + rechazados; ligas confirmadas; rutas desde la raíz |
| S5-06 OCR | alta | cubierto («sin medir» hasta las fotos) | §12, §15.1 | carpetas de fotos y transcripciones; `ErrNotConfirmed` |
| S5-07 Lectura asistida | alta | cubierto | §13 | `CanEnable` con 3 veredictos |
| S5-08 Sin red | alta | cubierto | §1 | sin cliente HTTP en WF-1 |
| S5-09 Membresía | alta | cubierto (B3 y B-r2-6 abiertos) | §8 | 8 compras; `missing_fee` |
| S5-10 Trazabilidad | media | cubierto | §3, §6 | 5 fuentes; promoción PROVISIONAL |
| S5-11 Mercados | media | cubierto | §9, §7.4 | 2026-10-13 recalculado |
| S5-12 Precios viejos | media | cubierto, con observación (M-r3-2) | §10 | «cotización que rige en la tienda asignada», sin regla de presentación |
| S5-13 Calidad Go | alta | cubierto | §1, §16 | |

**(b) Suficiencia.** Se puede implementar sin preguntas:
- el motor S4, con las cotizaciones dadas;
- el recetario y su esquema;
- el cargador de nutrición;
- el cargador PROFECO (con las opciones PROVISIONAL citadas);
- el OCR y su medición;
- los dos oráculos, con `Input`, funciones y salidas completas.

Hay dos huecos que el implementador de la API y de `QuoteFor` tendría que adivinar: M-r3-1 y M-r3-2.

**(d) Coherencia con §K:**
- **K7:** cortes de IMC 18.5 / 25 / 30; precio por 1 kg, 1 l o 1 pieza con HALF_UP; equilibrio en compras al año; reponer si existencia ≤ nivel y agotado si 0.
- **K22:** viejo a los 30 días; recta hasta el día 90; PROFECO nunca compite; ticket y OCR compiten con el mismo peso. Con K37, `manual` compite y `site_assisted` es referencia.
- **K23:** MSJ × factor de 1.2 a 1.9 ± ajuste; la persona confirma (§10.4).
- **K24:** 5 fotos; renglón correcto = descripción, cantidad e importe exactos; regla escrita antes de medir; «sin medir».

Todo coherente.

## Hallazgos nuevos

**M-r3-1 (media) · S4-12, S4-03 · `docs/contratos/motores/S4-interfaz.md` §6.5 (paso 3, último punto) y §12.2, contra C1 v1.2 (`docs/contratos/extractos/compras-tiendas-a.yaml`, `ShoppingListItem` y `ShoppingList.groups`).**

La opción A PROVISIONAL deja un artículo de vertical `market` en el mercado de visita más cercana (o sin tienda). Dice que la recomendación de S5 «se sigue viendo» en `ShoppingListItem.recommended_store_id`. Pero C1 define otra cosa:
- `recommended_store_id`: «Tienda o mercado que recomienda S5; **rige mientras `store_id` esté ausente**».
- `groups[].store_id`: «Tienda o mercado (**el elegido o el recomendado**)».

Además, el `items[].store` que calcula S4 no tiene campo en C1: `store_id` es la elección de la persona y se conserva (§6.3).

*Escenario:* uva verde (grupo `frutas-verduras`, canal `market`), sin elección. Un tianguis de martes surte el grupo y S5 recomienda Bodega Aurrera. S4 la pone en el grupo (`market`, tianguis, `visit_date` 2026-10-13). La API tiene que publicar `recommended_store_id` = Bodega, que según C1 rige, y un grupo cuyo `store_id` según C1 debería ser «el recomendado». La pantalla P-13 recibe dos tiendas que «rigen». Tampoco cuadran los nombres:
- S4 entrega `visit_date`; C1 tiene `visit_weekday`.
- S4 entrega `estimated_total_cents`; C1 tiene `total_cents`.
- S4 no entrega el total general `ShoppingList.total_cents`.

El bloqueo `S4-asignacion-mercado` no menciona el impacto en C1 y no hay solicitud.

*Corrección sugerida:*
1. En el bloqueo, en la opción A, anotar que exige cambiar C1.
2. Abrir una solicitud a C1:
   - un campo calculado `assigned_store_id` en `ShoppingListItem`, con la tienda de §6.5;
   - `groups[].store_id` = la tienda asignada;
   - `recommended_store_id` con la redacción «recomendación de S5 por precio; en vertical `market` no rige la agrupación».

   La alternativa es que, con la opción A, `recommended_store_id` sea la tienda de §6.5 y la cadena más barata solo se vea en `compare_basket`.
3. En §12.2, declarar la traducción `visit_date` → `visit_weekday` (o pedir `visit_date` a C1), `estimated_total_cents` → `total_cents`, y `ShoppingList.total_cents` = Σ de los grupos.

**M-r3-2 (media) · S4-03, S4-06, S5-03, S5-12 · `docs/contratos/motores/S5-interfaz.md` §14 (`QuoteFor`), §7.3, §7.4, §10 y §15.2; `S4-interfaz.md` §8.1 y §6.2 paso 5.**

`QuoteFor(product, store?, presentation?, date)` no dice **qué presentación** devuelve en dos casos:
- cuando solo hay tienda elegida («la de esa elección»);
- cuando no hay elección («la cotización de la tienda recomendada»), si esa tienda tiene varias presentaciones con precio.

Además, la regla de empate de §7.4 que cita («la que ya tiene asignados más renglones de la lista») depende del estado de una lista que `QuoteFor` no recibe. Tampoco dice qué `date` pasa la API para costos de recetas, plan y lista. §15.2 no tiene una consulta `quote_for`, así que ningún oráculo lo fija.

Por otro lado, `CompareBasket` (§7.3) elige la presentación por tienda con el **menor costo de paquetes completos para la necesidad**. C1 dice que su `recommended_presentation_id` «rige mientras `presentation_id` esté ausente», mientras que S4 calcula `packages_to_buy` y el costo con la presentación de la cotización `recommended = true` (§7.4: menor precio por unidad base).

*Escenario:* arroz, necesidad 500 g. Walmart tiene 1 kg a 28.00 (2800 por kg) y 750 g a 22.00 (2933 por kg).
- Por unidad base (§7.4) gana la bolsa de 1 kg: S4 calcula 1 paquete de 1 kg y precio 28.00.
- Por costo de paquetes (§7.3) gana la de 750 g: S5 publica `recommended_presentation_id` = 750 g y `estimated_price_cents` = 2200.
- El mismo `ShoppingListItem` lleva dos presentaciones que «rigen», y el motor de S5 y su oráculo pueden armar `quotes` distintos.

*Corrección sugerida:*
1. En S5 §14, fijar la regla de presentación de `QuoteFor`:
   - sin elección: el par (tienda, presentación) de menor precio por unidad base ajustado;
   - con tienda: el mismo criterio dentro de esa tienda;
   - con empate: un desempate que no dependa de la lista, por ejemplo menor precio de paquete y luego menor código.
2. Decir que `CompareBasket` usa, para la tienda recomendada, la presentación de la lista de S4 (o que `QuoteFor` recibe `needed_milli` y aplica §7.3), de modo que haya una sola presentación recomendada.
3. Fijar la `date` de referencia: hoy en America/Mexico_City.
4. Agregar `quote_for` a §15.2 con un caso `s5-cotizacion-recomendada`: dos presentaciones en una tienda y un empate.

**B-r3-1 (baja) · S4-01, S4-02 · `S4-interfaz.md` §3.2 y §12.1.**

En las fixtures, `LoadRecipeBook` corre en modo `review` y coteja los derivados contra la tabla nutricional del `Input`. Pero `nutrition_file` y `nutrition` son opcionales.
- Una fixture con `recipe_book_file` y sin tabla, como `s4-escalado-recetario`, que solo escala, daría `ErrKcalMismatch` en toda receta con `kcal_por_porcion` no nulo.
- Con `nutrition` en línea, `source` (`usda_fdc`) no es el `fuente` (`USDA`) de `fuentes_kcal[]`.

*Corrección:* «`recipe_book_file` exige `nutrition_file`» o «sin tabla nutricional en el `Input`, el cargador no coteja».

**B-r3-2 (baja) · S5-01, S5-05 · `S5-interfaz.md` §15.1 y §16; `S4-interfaz.md` §12.1.**

`go test ./internal/...` corre en la carpeta del paquete. El contrato fija que las rutas `*_file`, `files[]`, `links_file` y `tsv_file` son relativas a la raíz, pero no dice cómo la encuentran las pruebas Go (`TestComparativaCanon`, `TestProfeco*`).

*Corrección:* una línea, por ejemplo «las pruebas suben hasta la carpeta que contiene `backend/go.mod` y usan su padre como raíz».

## Bajos que siguen abiertos de rondas anteriores (no bloquean)
- r1: B1 (fuentes `mx_table` y `store_label`), B2 (claves `ai.*` en S4-14), B3 (RN-13 contra K7), B4 (prioridad de «agotado» en RN-15), B6 (`MAX(fecha_registro)` sobre texto), B7 (nombre «(5 pzas)» del atún) y B8 (cifra de bloques JSON en el reporte del autor). Los comprobé con `grep`: siguen igual.
- r2: B-r2-1 a B-r2-6, sin cambios.
