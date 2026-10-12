# Bloqueo S4-presentacion-respaldo

Abierto el 11-oct-2026 en la ronda P2c por el agente `contratos-puerta2`, al aplicar Q2-07 = A (hallazgo M-r5-2 de `docs/estado/wf1/REVISION-contrato-S4S5-r5.md`). Axel autorizó la corrección, pero la variante quedó sin recomendación y sin respuesta (PUERTA-2 §2.3 y §8). Se decide en WF-3 o cuando Axel conteste por id.

Contexto: sin cotización de S5, la presentación que rige un artículo de la lista era «la `presentacion_id` del ingrediente» (`S4-interfaz.md` 1.4 §6.2 paso 5). Esa regla no da una sola presentación en tres casos que pasan al arrancar, cuando la mayoría de los productos de las recetas aún no tiene precio: (a) el producto viene de varias recetas con presentación distinta, o de recetas de la casa sin presentación; (b) el artículo solo es de reposición por ritmo (motivo `replenishment`, sin ingrediente); (c) un artículo `calculated` sigue en la lista por una elección con necesidad 0 y sin ingrediente esa semana. Quedan sin definir `presentation`, `packages_to_buy` y `leftover_milli`. Ejemplo del revisor: papel higiénico por ritmo de consumo, sin precio que compita.

HECHO (código leído el 11-oct-2026):
- Motor (`backend/internal/recetas/shopping.go` y `catalog.go`, `fallbackPresentation`): toma la `presentation` del primer ingrediente no nulo del producto, recorriendo las entradas de la semana en orden de fecha, tiempo y persona; si ninguno la trae, la `presentacion_tipica` del producto (solo llega con `catalog_file`); si no, la primera presentación del producto en el catálogo; si el producto no tiene presentaciones, los tres campos en `null`.
- Oráculo (`oraculo/rec_plan.py`, `pres_ing.setdefault`): toma la presentación del primer ingrediente aunque sea nula, y sin ingrediente deja la presentación en `null`, sin respaldo.
- Ninguna fixture prueba los casos (a), (b) y (c).

Qué falta (pregunta para Axel): sin precio, ¿qué presentación rige un artículo de la lista en los casos (a), (b) y (c)?

Opciones (las tres variantes del revisor, sin marcar preferida):
- i) La presentación típica del catálogo (`presentacion_tipica` de `productos.json`, que `productos-formato.md` §2 describe como «la que sugiere la receta y la lista cuando S5 no tiene precio»); se agrega a `catalog.products[]` como `typical_presentation` para el catálogo en línea.
- ii) Para el caso (a), la de la primera receta en orden de fecha y tiempo; para (b) y (c) hace falta además i) o iii).
- iii) `null` en `presentation`, `packages_to_buy` y `leftover_milli`, como el artículo `manual` sin cotización (§6.6); el producto va en `missing_prices[]`.

Impacto: solo la presentación, los paquetes y el sobrante de artículos sin precio; ninguna cifra canon (el canon tiene precio en todos sus productos). Con iii) la lista del arranque muestra «sin precio» sin cantidad de paquetes.

Cómo seguí: `S4-interfaz.md` 1.5 §6.2 paso 5 escribe como PROVISIONAL la regla del motor (ii para el caso (a), luego i, luego la primera presentación del catálogo). El oráculo tiene que alinearse a ella y hace falta una fixture con los tres casos (solicitud a los agentes de S4 en WF-3). Sin recomendación.

- **Estado:** ABIERTO
