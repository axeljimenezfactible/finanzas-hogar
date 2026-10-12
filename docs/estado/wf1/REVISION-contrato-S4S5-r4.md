# Revisión del contrato S4/S5 · ronda 4 (WF-1, 11-oct-2026)

Revisor: `revisor-contrato-S4S5-r4` (independiente; no escribí nada de lo revisado y no corregí nada). Evalué solo contra `docs/criterios/S4.md` 1.2, `docs/criterios/S5.md` 1.2 y decisiones §D, §E y §K. Artefactos: `docs/contratos/motores/S4-interfaz.md` 1.3, `S5-interfaz.md` 1.3 y `productos-formato.md` 1.1; reportes `docs/estado/wf1/contrato-S4S5.md` y `correccion-contrato-S4S5-r3.md`. Contraste con C1 v1.2 por el extracto `docs/contratos/extractos/compras-tiendas-a.yaml` (`ShoppingListItem`, `ShoppingList`, `ShoppingListItemCreate`).

## Veredicto
**NO APROBADO: 0 altos, 2 medios, 3 bajos nuevos.** Los dos medios de la ronda 3 (M-r3-1 y M-r3-2) quedan **cerrados**. Quedan abiertos dos medios nuevos:
- **M-r4-1.** Viene de la corrección de M-r3-2. En la comparativa por canasta, la tienda recomendada se cotiza con una regla y las demás con otra, y eso puede invertir la tienda más barata (S5-03).
- **M-r4-2.** S4 no define la salida de los artículos `manual` ni qué hace la cantidad que fija una persona (`quantity_milli`), y C1 sí la usa (S4-03, S4-06).

Todas las cifras canon coinciden con mi cálculo independiente.

## Verificación propia (scripts en `datos/tmp/rev-s4s5-r4/`, ya borrados)
- `python3 -I verif.py` (solo biblioteca estándar, `Decimal` y `Fraction`; lee `data.json` y `productos-formato.md`): **29 OK, 0 FAIL**.
  - **Canon S5-01, desde los renglones de `data.json`:**
    - idénticos: Bodega **277350**, Walmart **287400**, diferencia **+10050 (+$100.50)**;
    - totales: Bodega **309450** = 277350 + 30000 (sustitutos) + 2100 (sexta lata); Walmart **323300**; diferencia **+13850 (+$138.50)** = 10050 + 5900 − 2100;
    - coinciden con `totales` del archivo;
    - piezas 61 / 60 / 55, renglones 33 / 3, ganadores 22 / 8 / 3 (= `resultados_por_producto`).
  - **S5-02:** 2933 / 2800, 12692 / 16406 y 3375 / 2833.
  - **K40:** `[0.9836, 0.4918, 0.0]`; fracciones 60/61, 59/61 y 30/61. El 2026-11-09 está a 30 días del 10-oct-2026 y el 2027-01-08 a 90.
  - **S4 §8.3:** 4425 / 2213 / 8850, con los precios de `items_identicos[21]`, `[16]`, `[15]` y `[0]`.
  - **S4-07:** MSJ 1748.75 y 1330.25; × 1.2 = 2098.50 → 2099 y 1596.30 → 1596; con −300 → 1799 y con +250 → 1846. RN-17: 2555.5625 → 2556, y con −500 → 2056. RN-15: 1486571.
  - **Catálogo (e):**
    - 39 presentaciones en §5.1, cuyos nombres son exactamente los 33 + 3 + 3 textos de `data.json`;
    - los 39 renglones de §5.2 son iguales a `cuentas` en texto, piezas e importe;
    - **61 piezas y 309450**;
    - el precio unitario de cada renglón es igual al `bodega_unit` (o total ÷ cantidad) de su presentación;
    - las piezas por presentación cuadran con `items_identicos[].cantidad` (atún 5 + 1).
- `python3 -I quote.py` es mi implementación de S5 §14.1 y §7.3 escrita desde el texto. Reproduce el ejemplo de `s5-cotizacion-recomendada`:
  - sin elección → Walmart 1 kg a 2800;
  - `store` = Bodega → 750 g a 2200;
  - `presentation` = 750 g → Walmart, por el orden de `stores`;
  - frijol → Bodega 500 g a 1500, por el precio del paquete;
  - `compare_basket`: Walmart 2800 y Bodega 1500, `cheapest_store` = Bodega.

  También da el contraejemplo de M-r4-1.
- Estructura:
  - S4: 1 bloque JSON válido y cercas pares; S5: sin bloques;
  - S4 §15 cubre S4-01 a S4-15 y S5 §18 cubre S5-01 a S5-13;
  - cada `-run` de los criterios casa solo con las pruebas de su criterio (§13 de S4, §16 de S5);
  - las claves `price.*` y `health.goal_adjust_*` existen en `00001` y `00003`;
  - RN-16: el sábado 10-oct-2026 lleva a `NextVisit` = martes 2026-10-13.

## Estado de los hallazgos de la ronda 3
| Hallazgo | Estado | Evidencia |
|---|---|---|
| M-r3-1 (S4-12, S4-03; tienda asignada frente a C1) | **cerrado** | S4 §6.5 paso 3 ya dice que la opción A exige cambiar C1 y trae el ejemplo de la uva. §12.2 traduce `visit_date` → `visit_weekday`, `estimated_total_cents` → `total_cents` y `ShoppingList.total_cents` = Σ de los grupos, y propone `items[].store` → `assigned_store_id` (solicitud), con una variante sin campo nuevo. La raíz entrega `estimated_total_cents`. El bloqueo `S4-asignacion-mercado` (opciones 2A y 3A, Impacto y «Cómo seguí») y la solicitud 1 del reporte r3 citan el impacto en C1. Comprobé en el extracto que C1 v1.2 tiene `recommended_store_id` («rige mientras store_id esté ausente»), `groups[].store_id` («el elegido o el recomendado»), `visit_weekday` y `total_cents` obligatorio, y que la traducción los respeta |
| M-r3-2 (S4-03, S4-06, S5-03, S5-12; regla de `QuoteFor`) | **cerrado** | S5 §14.1 fija la fecha (`today`; la única excepción es §10), los candidatos, el orden (1)–(4) sin depender de la lista (PROVISIONAL, citado) y `recommended`. Además, `quote_for` está en §15.2 y el caso `s5-cotizacion-recomendada` en §15. S4 §6.2 paso 5 y §8.1 usan esa regla; con tienda y sin presentación, la elige §14.1. El ejemplo del arroz ya da una sola presentación (1 kg, 2800) en S4 y en S5, y mi implementación lo reproduce. Efecto colateral de la corrección: M-r4-1 |

Bajos anteriores sin cambios, que no bloquean: B-r3-1 (cotejo de kcal sin tabla nutricional en el `Input`) y B-r3-2 (cómo hallan la raíz del repositorio las pruebas Go; `grep -c go.mod` = 0 en los dos contratos), más los de r1 y r2.

## Criterios
| Criterio | Sev. | Estado | Sección | Mi evidencia |
|---|---|---|---|---|
| S4-01 Recetario | alta | cubierto | §3 | esquema en `datos/semilla/recetario.schema.json`, llave `recetas`, validación cruzada, modos `review` y `production`; el bloque JSON de §3.1 es válido |
| S4-02 Escalado exacto | alta | cubierto | §4 | fracción reducida y tabla `fractions`; prohíbe `float` (§1) |
| S4-03 Lista = Σ − stock | alta | **parcial** | §6, §12.2 | artículos `calculated` completos (pasos 1–7, sobrantes, cotización que rige); artículos `manual` y `quantity_milli` sin definir (**M-r4-2**) |
| S4-04 Inventario híbrido | alta | cubierto | §7 | RN-15 = 1486571 (verificado); agotado prevalece |
| S4-05 Plan semanal | alta | cubierto | §5 | K38 (domingo → lunes siguiente), K43 (máximo 21), Karen rechazada, «fuera» |
| S4-06 Costo | alta | cubierto | §8 | §8.3 = 4425 / 2213 / 8850 (verificado); «sin precio» ≠ 0; fecha `today` y regla de S5 §14.1. El costo estimado de la lista con cantidad fijada a mano entra en M-r4-2 |
| S4-07 Salud, IMC, meta | alta | cubierto | §10 | 9 cifras verificadas (arriba); factores 6/5, 11/8, 31/20, 69/40, 19/10 = RN-17; K36 por persona con valor de la casa |
| S4-08 kcal con fuente | alta | cubierto | §9 | `tabla.csv` con `estado` y `confirmo`; `Coverage` con dos conteos |
| S4-09 Cinco recetas a mano | alta | cubierto (procedimiento) | §12 | las elige el revisor de la implementación; no aplica a esta revisión de contrato (no hay recetario) |
| S4-10 Cargas locales | media | cubierto | §9.1, §1 | `PROCEDENCIA.md` con CC0, ODbL y SMAE manual; 0 red |
| S4-11 Español | alta | cubierto | §3.2, §9.1 | `nombre_fuente` solo como procedencia; `production` no carga lo `propuesto` |
| S4-12 Canal | media | cubierto (PROVISIONAL citado) | §6.5, §11, §12.2 | M-r3-1 cerrado |
| S4-13 Recetas editables | media | cubierto | §3.3 | funciones puras y consulta SQL del campo vacío |
| S4-14 IA apagada | media | cubierto | §1, §13 | `TestSinIA` |
| S4-15 Calidad Go | alta | cubierto | §1, §13 | sin SQLite en `internal/recetas`; cobertura ≥ 85 % |
| S5-01 Comparativa canon | alta | cubierto | §7.2, §3, §15 | **+138.50 y +100.50 verificados** (arriba) |
| S5-02 Unidad base | alta | cubierto | §4 | 2933 / 2800, 12692 / 16406, 3375 / 2833 |
| S5-03 Mejor tienda | alta | **parcial** | §7.3, §7.4, §14.1 | por ingrediente, correcto (4 resultados reproducidos); por canasta, regla asimétrica (**M-r4-1**) |
| S5-04 Frescura y fuentes | alta | cubierto | §5, §3 | K40 exacto; 3 fuentes compiten y 2 son referencia (K22, K37) |
| S5-05 PROFECO | alta | cubierto | §11 | conteo leído = guardados + rechazados; reporte por cadena y municipio; 3 sucursales con regla escrita antes de contar |
| S5-06 OCR | alta | cubierto | §12 | 5 fotos en `tickets/`, transcripciones fuera; «sin medir» sin fotos (K24) |
| S5-07 Lectura asistida | alta | cubierto | §13 | `ErrNoTermsVerdict`; sigue como referencia |
| S5-08 Sin red | alta | cubierto | §1 | `asistida` sin cliente HTTP |
| S5-09 Membresía | alta | cubierto | §8 | 100000 ÷ 12500 = 8; `missing_fee` |
| S5-10 Trazabilidad | media | cubierto | §3, §6 | 5 fuentes con su marca; promoción vencida (PROVISIONAL citado) |
| S5-11 Mercados | media | cubierto | §9, §7.4 | RN-16 → 2026-10-13 (verificado) |
| S5-12 Precios viejos | media | cubierto | §10, §14.1 | evalúa en la fecha de despensa (ver B-r4-1) |
| S5-13 Calidad Go | alta | cubierto | §1, §16 | sin `float`; cobertura ≥ 85 % |

**(b) Suficiencia.** El motor S5 (salvo M-r4-1), PROFECO, OCR, recetario, nutrición y los dos oráculos (§12.1–§12.3, §15.1–§15.3) se pueden implementar sin preguntas. En el motor S4 queda la pregunta de M-r4-2.

**(d) Coherencia con §K.**
- **K7:** IMC 18.5 / 25 / 30 (S4 §10.1); precio por 1 kg, 1 l o 1 pieza con HALF_UP (S5 §4); equilibrio en compras al año (S5 §8); reponer y agotado (S4 §7).
- **K22:** viejo a los 30 días o más; recta hasta el día 90; desde el 91 no compite; PROFECO nunca compite; ticket y OCR compiten con el mismo peso.
- **K23:** Mifflin-St Jeor × 5 factores de 1.2 a 1.9, ajuste por objetivo; la persona confirma (S4 §10.4).
- **K24:** 5 tickets reales; renglón correcto = descripción, cantidad e importe exactos; arranca sin fotos.

Todo coherente, y también con K36, K37, K38, K40, K43 y K44.

## Hallazgos nuevos

**M-r4-1 (media) · S5-03 · `docs/contratos/motores/S5-interfaz.md` §7.3 (párrafo «Recomendación por renglón y tienda recomendada») y §14.1 («`recommended`… El menor costo de paquetes para la necesidad (§7.3) solo elige la presentación de las tiendas no recomendadas»).**
En una misma `compare_basket` hay dos reglas:
- la tienda recomendada de cada renglón se cotiza con la presentación de menor precio por unidad base, «aunque otra de la misma tienda cubra la necesidad con menos costo de paquetes»;
- las demás tiendas, con el menor costo de paquetes para la necesidad.

El total por canasta de una tienda se infla justo en los renglones donde es la más barata por unidad base. Así, `cheapest_store` y `difference_vs_cheapest_cents` pueden señalar una tienda que no es la más barata para la lista completa, que es lo que pide S5-03 («la tienda más barata… por canasta (lista de compras completa)»).

*Escenario (lo calculé con mi implementación de §14.1 y §7.3).* Arroz, necesidad 500 g. Walmart tiene 1 kg a 28.00 (2800 por kg) y 500 g a 14.50 (2900 por kg); Bodega, 500 g a 14.60.
- Se recomienda Walmart 1 kg.
- La canasta cotiza Walmart en **2800** y Bodega en **1460**: `cheapest_store` = Bodega y Walmart sale +13.40.
- Pero Walmart cubre la necesidad por **1450**: es más barata por 0.10.

Con muchos renglones, la tienda que más renglones gana por unidad base es la más castigada en el total de la canasta. Ninguna cifra canon cambia, porque el canon usa `CompareCarts`.

*Corrección sugerida.* Una sola regla para el total de la canasta en todas las tiendas: el menor costo de paquetes para cubrir la necesidad (la regla base de §7.3, también en la tienda recomendada). La recomendación por renglón (`items[].recommended_*`, `estimated_price_cents`) sigue siendo `QuoteFor` sin elección, igual a la de S4. Como §14 ya dice que los campos de `ShoppingListItem` los llena la API con la salida de S4, la presentación que rige sigue siendo una sola. Hay que decir que la línea de esa tienda en `prices[]` puede mostrar otra presentación, como sustituto. Hay que ajustar el ejemplo de §7.3, la frase de §14.1 y el caso `s5-cotizacion-recomendada`: con la regla única, Walmart cotiza la de 750 g a 2200 y `cheapest_store` sigue siendo Bodega con 1500. Si se prefiere otra lectura, hay que abrir la pregunta en `S5-reglas-comparativa` en vez de dejar dos reglas.

**M-r4-2 (media) · S4-03, S4-06 · `docs/contratos/motores/S4-interfaz.md` §6.2, §6.5 y §12.1–§12.2; `S5-interfaz.md` §7.3 y §15.1; contra C1 v1.2 (`ShoppingListItem.quantity_milli`, `ShoppingListItemCreate`).**
`list_choices` trae `kind` = `manual` y `quantity_milli`, y `shopping_list` devuelve `items[].kind`, pero ninguna regla dice:
1. Qué `needed_milli`, `packages_to_buy`, `leftover_milli` y `estimated_cents` lleva un artículo `manual`. C1 exige `quantity_milli` al crearlo y no tiene necesidad calculada. Tampoco dice si un artículo `manual` suma en `groups[].estimated_total_cents` y en el total de la raíz.
2. Si en un artículo `calculated` la cantidad fijada a mano cambia `packages_to_buy` y `estimated_cents`. C1 dice «`quantity_milli`… NULL en calculated = comprar la necesidad», es decir, si no es NULL se compra esa cantidad. S4 §6.2 paso 5 calcula los paquetes siempre con `needed_milli`, y `quantity_milli` no se usa en ninguna parte.
3. Del lado de S5, §7.3 calcula `estimated_price_cents` con `quantity_milli`, pero el `Input` `shopping_list` (§15.1) solo trae `needed_milli`.

*Escenario.* Ximena agrega a mano «aceite de soya», 2 botellas (`quantity_milli` = 1892000), y fija 3 kg de arroz en un artículo calculado que necesita 500 g. El implementador tiene que adivinar si el arroz lleva 1 paquete o 3 (o 4 de 750 g), qué total de grupo sale y si el aceite cuenta. El motor y el oráculo pueden divergir.

*Corrección sugerida.* En S4 §6.2 (o en un §6.6 nuevo):
- cantidad a comprar = `quantity_milli` si no es nula; si es nula, `needed_milli`;
- `packages_to_buy`, `leftover_milli` y `estimated_cents` se calculan con esa cantidad;
- los artículos `manual` tienen `needed_milli` ausente y la misma regla con su `quantity_milli`;
- los dos suman en `groups[]` y en la raíz; aclarar si `is_bought` resta o no.

Agregar una fixture en `s4-lista-*` o `s4-costo-cotizacion-elegida`. En S5 §15.1, agregar `quantity_milli` al renglón de `shopping_list`, o decir que §7.3 usa `needed_milli`.

**B-r4-1 (baja) · S5-12 · `S5-interfaz.md` §10.** `StalePricesForPlan` cotiza con la tienda elegida o sin elección, pero no con la presentación elegida. Si la persona eligió solo la presentación, S4 calcula con `QuoteFor(presentation)` (S4 §8.1), y la alerta revisa otra cotización: puede no avisar de un precio viejo que la lista sí usa. *Corrección:* pasar también la presentación elegida (`assignedPresentations`) y llamar `QuoteFor` con las dos.

**B-r4-2 (baja) · S5-03 · `S5-interfaz.md` §7.3.** El empate entre presentaciones de una tienda dice «luego la de menor id». Las fixtures citan por código. *Corrección:* «luego la que va antes en el catálogo (en la API, menor id)», como en §14.1 (4).

**B-r4-3 (baja) · S4-12 · `S4-interfaz.md` §12.2, «Variante sin campo nuevo».** Si no hay mercado dado de alta (`items[].store` nulo en vertical `market`), la variante publica `recommended_store_id` ausente. En cambio, `recommended_presentation_id` y `estimated_price_cents` salen de la cadena de la cotización, y C1 diría que rige «el que recomienda S5» sin decir cuál. *Corrección:* una línea para ese caso, por ejemplo publicar la cadena de la cotización o dejar el precio sin tienda con la leyenda «sin mercado dado de alta».
