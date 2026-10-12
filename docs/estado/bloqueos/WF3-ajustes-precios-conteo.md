# Bloqueo WF3-ajustes-precios-conteo · conteo de precios que pasan de fresco a viejo en P-26.06

Abierto por `correccion-frontend-areas-r1` (continuación 1, WF-3a, 12-oct-2026). Estado: ABIERTO.

## Qué falta
U4 P-26.06, vista previa 1: al cambiar «Días para que un precio sea viejo», el texto de frescura se recalcula (HECHO) y además «dice cuántos precios cambian de fresco a viejo hoy». C1 v1.4 no tiene una operación que dé ese conteo con un valor hipotético: `getPriceSummary` (N-172) solo cuenta con el valor vigente (`fresh_count`, `stale_count`). U4 tampoco dice qué se cuenta: «precios» puede ser cada precio por tienda o el último precio de cada producto (que es lo que cuenta `getPriceSummary`: «Productos con precio fresco»).

## Opciones
- **A (en uso, PROVISIONAL):** la pantalla cuenta en el navegador los productos de `listProducts` (`limit` 200) cuyo último precio (`last_observed_date`) tiene hoy entre el valor nuevo y el vigente menos 1 días. Es la misma unidad que `getPriceSummary`. Si el catálogo trae más de una página, dice «Vista previa pendiente». Con el valor nuevo igual o mayor que el vigente da 0.
- **B:** agregar a C1 un parámetro `stale_after_days` en `getPriceSummary`, o una operación de vista previa, para que el servidor dé el conteo con el valor nuevo. Es un cambio de contrato y de backend.
- **C:** contar cada precio por producto y tienda, no solo el último. Para eso hace falta también la opción B, porque el navegador no tiene todos los precios.

## Impacto
Solo cambia una nota de la vista previa: no guarda nada ni toca cifras de los motores. Mientras tanto, la nota lleva «PROVISIONAL (bloqueo WF3-ajustes-precios-conteo)» (`frontend/src/pantallas/ajustes-a/S06Tiendas.tsx`).
