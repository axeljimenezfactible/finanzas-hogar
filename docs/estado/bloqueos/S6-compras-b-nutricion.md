# Bloqueo S6-compras-b-nutricion

Abierto por `api-compras-b` (WF-3a, 12-oct-2026). Se pregunta en la puerta 4. Nutrición, recetario e inventario. Código: `backend/internal/api/comprasb/fuentes.go`, `inventario.go`, `agente.go`, `backend/internal/dominio/compras-b_repo.go`.

## Qué falta (PREGUNTAS)

1. **`reloadNutritionBase`: qué archivo local y con qué llave.** `product` no tiene código y `datos/semilla/nutricion/tabla.csv` usa el `producto_id` de la semilla (por ejemplo `aceite-oliva`); `raw/` trae los volcados de USDA y Open Food Facts sin un formato acordado. Opciones: (a) un CSV por base en `<DataDir>/nutricion/<source>.csv` con `product.name` como llave; (b) S1 agrega `product.code`; (c) la recarga se hace con `finanzas semilla` y la operación solo informa.
2. **Dato preferido nuevo (`createNutritionFact` con `is_preferred = true`) cuando ya hay uno.** Opciones: (a) 409 `already_exists` (la API no cambia sola el valor que usan los cálculos); (b) el nuevo pasa a preferido y el anterior deja de serlo en la misma transacción.
3. **Pasar un artículo a ritmo de consumo (`updateInventoryItem`) sin fecha de conteo.** El CHECK de S1 pide `counted_date`, que no está en `InventoryItemUpdate`. Opciones: (a) 422 y primero un ajuste de conteo; (b) la API pone hoy como fecha de conteo.
4. **Pasar a `stock` con ritmo guardado.** Opciones: (a) 422 si no se mandan `rate_quantity_milli` y `rate_period_days` en null; (b) la API los borra sola.
5. **Recetas del agente externo (preguntas c y f de S5-agente-externo).** Implementado (f2): `origin = household`, `author_person_id` NULL, `is_active = 0`; `notes` se guarda tal cual (sin marca del agente).
6. **Orden permitido en colecciones** (C1 no lista las propiedades): precios `observed_date`, `price_cents`, `created_at` (por defecto `-observed_date`); nutrición `source`, `created_at` (por defecto `id`); inventario `status` (por defecto, y empate por nombre e id); movimientos `movement_date` (por defecto `-movement_date`).
7. **`createInventoryItem` en ritmo de consumo sin `counted_date`.** Hoy 422 `required` (el CHECK lo exige y C1 no da un valor por defecto).

## Mientras tanto (PROVISIONAL)

1: 409 `conflict` sin tocar `nutrition_fact` · 2 (a) · 3 (a) · 4 (a) · 5 (f2) · 6 como arriba · 7: 422.

## Impacto

`reloadNutritionBase` no tiene camino 2xx hasta resolver 1. Lo demás cambia solo validaciones y sus pruebas.

**Estado: ABIERTO.**
