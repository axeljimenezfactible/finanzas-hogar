# Bloqueo S6-compras-a-motores · Cómo la API de compras (parte a) alimenta y lee los motores S4 y S5

Abierto 12-oct-2026 por **api-compras-a** (WF-3a). Estado: **ABIERTO**. Todo lo de abajo está implementado como PROVISIONAL en `backend/internal/api/comprasa/` y se puede cambiar sin migración.

| # | Qué falta (PREGUNTA) | Opciones | Lo que hace hoy la API (PROVISIONAL) | Impacto |
|---|---|---|---|---|
| 1 | El motor S4 en Go (`recetas.Quote`) no trae `choice_store`/`choice_presentation` de S4-interfaz 1.5 §8.1 (M-r5-1): indexa cotizaciones por (tienda, presentación). | a) actualizar el motor S4 a 1.5; b) dejar la API como está | Manda la recomendada de S5 y una por cada elección distinta, sin repetir (tienda, presentación). | Caso r5 (calculated eligió tienda y manual eligió presentación del mismo producto): el motor puede responder `invalid_argument` → la API da 422/409 con mensaje. |
| 2 | El esquema 1.4 no guarda quién tiene plan de comidas (`Person.HasMealPlan`). | a) columna nueva; b) regla por participación | `participation <> 'external'` (Axel y Ximena sí; Karen no). | Un renglón de plan de Karen haría fallar la lista. |
| 3 | Vertical de un artículo en texto libre (S4 no lo calcula). | a) siempre cadena; b) según la tienda elegida; c) grupo propio | Mercado si la tienda elegida es mercado; si no, cadena (sin tienda si no eligió). | Solo agrupación; no suma precio. |
| 4 | ¿Una tienda dada de baja compite? ¿Y una presentación inactiva? | — | Tiendas inactivas fuera del libro de S5; presentaciones inactivas dentro. | Recomendaciones y comparativa. |
| 5 | `getShoppingList` cuando el plan cambió después de generar. | a) mostrar solo lo guardado; b) agregar lo nuevo | Solo artículos guardados; `needed_milli` = última generación; piezas, tienda y precio con el plan y precios de hoy. | Lo nuevo aparece al volver a generar. |
| 6 | Al regenerar, artículos `calculated` que el plan ya no pide. | a) borrarlos; b) dejarlos en 0 | Se borran si nadie eligió nada en ellos (tienda, presentación, cantidad o comprado); con elección se conservan (inv. 13). | Lista más limpia. |
| 7 | `basket=shopping_list`: `precios.CompareBasket` no entrega `winners`, `ties`, `winner_store_id` ni `difference_cents`. | a) que S5 los dé; b) quitarlos del contrato en esa modalidad | `winners` = [], `ties` = 0 y los dos campos ausentes (no se recalculan a mano). `is_substitute` = algún precio sustituto. | La tarjeta de ganadores solo sale con `basket=expense`. |
| 8 | `basket=expense`: carrito base de S5. | — | `Kind` "expense", fuente `ticket_manual`; renglón sin presentación se omite; sin precio unitario = total ÷ cantidad (ROUND_HALF_UP); compra sin tienda activa → 422; `price_cents` = total del renglón (`RowLine.TotalCents`). | Comparativa de compras capturadas. |
| 9 | Historial de la membresía (N-164). | — | Compras con renglones y tienda activa; solo con `basket=shopping_list`; la primera tienda que pide membresía. | Veredicto de Sam's. |
| 10 | Filtro `source` de `listProducts`. | a) fuente del último precio; b) cualquier precio | Fuente del último precio (el que fija `freshness`, con `precios.Freshness`). | Filtro de catálogo. |
