# Reporte · api-comidas-salud (WF-3a «Construcción», Opus) · 12-oct-2026

## Hecho
- **26 de 26 operaciones** registradas y probadas: `comidas.yaml` (16), `salud.yaml` (6), `createHealthProfile`, `updateHealthProfile` (ajustes-d) y `getHealthSettings`, `updateHealthSettings` (ajustes-f). Un solo `Register(reg)` en el paquete `internal/api/comidas` (nombre fijado por CONVENCIONES-backend §2).
- Cálculos **solo con los motores**: costo, kcal, escalado, semana, huecos (M-05), sobrantes, IMC, edad, clase de IMC y meta sugerida salen de `internal/recetas` (S4: `WeekPlanFor`, `RecipeDetailFor`, `HealthFor`, `DefaultWeek`, `ValidateRecipe`, `RemoveRecipe`); el precio base sale de `internal/precios` (S5: `NewBook` + `QuoteFor`, conectado como `recetas.PriceSource`). No se tocó ningún motor.
- Plan semanal con semana por defecto de S4 (K38: los domingos, la siguiente), copiar semana anterior solo en casillas vacías, «fuera» sin mezclar con recetas (409 `slot_conflict`), Karen fuera del plan y de salud (422 `external_person`).
- Recetas: alta, cambio (JSON Merge Patch), duplicado «(copia)», escalado con `servings_milli`, costo y kcal por porción, `times_used_8w`, orden por campos calculados, archivo lógico cuando el plan la usa (C-06), 409 `recipe_planned`.
- IA: `suggestRecipeWithAi` responde 409 `ai_disabled` con la IA apagada (o sin registro de IA) **sin ninguna llamada de salida** (S6-13).
- Salud: perfiles con meta sugerida (K23, K36: ajuste propio o el de la casa), ajuste de la casa con `ReadParameter`/`ApplyParams` (K44, cambiar y restablecer, 422 si van juntos), pesos con IMC a la fecha (motor) y **captura atrasada idempotente K9**: la clave queda en `weight_entry.idempotency_key`; pasadas las 24 h un reintento responde 200 con el mismo peso e `Idempotent-Replayed: true`.
- Cada alta, cambio o baja escribe `audit_log` (insert/update/delete) dentro de su transacción. Dinero en centavos enteros; fechas ISO; mensajes en español de México.

## Archivos
- `backend/internal/api/comidas/comidas.go` (Register, adaptador SQLite → S4/S5, ayudantes)
- `backend/internal/api/comidas/plan.go` (tiempos, semana, mes, tendencias, casillas, copiar semana)
- `backend/internal/api/comidas/recetas.go` (recetas, IA, uso, duplicado)
- `backend/internal/api/comidas/salud.go` (perfiles, kcal de hoy, pesos, ajuste de la casa)
- `backend/internal/api/comidas/comidas_test.go` (26 pruebas `TestOp_<operationId>`)
- `backend/internal/dominio/comidas-salud_datos.go` (repositorios de meal_slot, meal_plan_entry, recipe*, health_profile, weight_entry y lectura de catálogo, nutrición, tiendas, precios y sobrantes para los motores)
- `backend/internal/dominio/comidas-salud_datos_test.go`
- `docs/estado/bloqueos/S6-comidas-calculos.md` (13 puntos PROVISIONAL)

## Verificación
| Comando | Salida resumida |
|---|---|
| `go build ./...` · `go vet ./...` | sin errores |
| `gofmt -l .` | vacío |
| `go test ./internal/api/comidas/... ./internal/dominio/... -count=1 -cover` | `ok api/comidas 81.5 %` · `ok dominio 87.7 %` (mi archivo de dominio: 89.3 %) |
| `make verify-wf3-go` | build, vet, gofmt, `go mod tidy -diff` sin cambios, `go test ./...`: 21 paquetes OK, 6 sin pruebas |
| `make verify-all GOOSE=…` | **no concluyó**: `verify-sql` falló dos veces con «attempt to write a readonly database» y «disk I/O error» sobre `datos/tmp/…/sql-*.db` (la base temporal compartida que borran y recrean otros agentes que corren la verificación al mismo tiempo). Este agente no toca migraciones ni SQL de esquema; el integrador la repite en serie |

- **S6-03** · 26 operaciones con su prueba; cada respuesta pasa por `testutil` (kin-openapi contra `openapi.yaml`, código declarado + esquema). Tabla abajo. Las rutas no cambian (las genera c1gen).
- **S6-04** · errores con el esquema único (400 formato, 422 negocio con `fields`, 404, 409 con `slot_conflict`, `recipe_planned`, `already_exists`, `ai_disabled`, `external_person` 422); páginas de 1 a 200 (`limit=201` → 400 probado); Idempotency-Key en las 6 operaciones que la piden (falta → 400 probado); K9 de peso probado con la tabla `idempotency_key` purgada.
- **S6-13** · sin IA activa → 409 `ai_disabled`, cero llamadas de salida; no hay llaves ni tokens en código ni pruebas.
- **S6-15** · identificadores en inglés, mensajes en español de México; el único registro (`Warn` si el motor rechaza una receta) lleva `recipe_id` y el código del motor, nunca credenciales.

### Operación → prueba (todas PASS)
| Operación | Prueba | 2xx | 4xx probados |
|---|---|---|---|
| listMealSlots | TestOp_listMealSlots | 200 | 400, 401 |
| getMealPlanWeek | TestOp_getMealPlanWeek | 200 (costo 500, kcal 530) | 422, 400 |
| getMealPlanMonth | TestOp_getMealPlanMonth | 200 | 422 |
| getMealPlanTrends | TestOp_getMealPlanTrends | 200 | 400 |
| createMealPlanEntries | TestOp_createMealPlanEntries | 201 | 409 slot_conflict, 422 external_person, 422, 400 sin clave |
| updateMealPlanEntry | TestOp_updateMealPlanEntry | 200 | 422, 404 |
| deleteMealPlanEntry | TestOp_deleteMealPlanEntry | 204 | 404 |
| copyPreviousWeek | TestOp_copyPreviousWeek | 201 | 400 |
| listRecipes | TestOp_listRecipes | 200 | 400 |
| createRecipe | TestOp_createRecipe | 201 | 409, 422, 400 |
| suggestRecipeWithAi | TestOp_suggestRecipeWithAi | — (IA apagada) | 409 ai_disabled, 422 |
| getRecipe | TestOp_getRecipe | 200 (×2 porciones: 2000) | 422, 400, 404 |
| updateRecipe | TestOp_updateRecipe | 200 | 422, 404 |
| deleteRecipe | TestOp_deleteRecipe | 204 (archivo y borrado) | 409 recipe_planned, 404 |
| getRecipeUsage | TestOp_getRecipeUsage | 200 | 404 |
| duplicateRecipe | TestOp_duplicateRecipe | 201 | 409, 404 |
| listHealthProfiles | TestOp_listHealthProfiles | 200 (meta 2664, IMC 261) | 400 |
| createHealthProfile | TestOp_createHealthProfile | 201 | 409, 422 external_person, 422 |
| updateHealthProfile | TestOp_updateHealthProfile | 200 (lose 500 → 2164) | 422, 400, 404 |
| getTodayKcal | TestOp_getTodayKcal | 200 | 401 |
| listWeightEntries | TestOp_listWeightEntries | 200 | 400 |
| createWeightEntry | TestOp_createWeightEntry | 201, 200 repetido (K9) | 409, 422 ×2, 400 |
| updateWeightEntry | TestOp_updateWeightEntry | 200 | 422, 404 |
| deleteWeightEntry | TestOp_deleteWeightEntry | 204 | 404 |
| getHealthSettings | TestOp_getHealthSettings | 200 | 401 |
| updateHealthSettings | TestOp_updateHealthSettings | 200 (cambiar y restablecer) | 422 ×2 |

Las cifras esperadas de las pruebas son de una persona y un catálogo ficticios (nunca datos de Axel o Ximena) y se comprobaron a mano contra las fórmulas de S4 §8–§10.

## Supuestos (HECHO = lo que se hizo; ninguno cambia cifras de los motores)
- Paquete `internal/api/comidas` (CONVENCIONES §2) en lugar de `api/comidas-salud/` (un nombre de paquete Go no lleva guion).
- Códigos internos del adaptador: producto `p<id>`, presentación `pr<id>`, grupo `g<id>`, tienda `s<id>`, receta `r<id>`; persona = `person.code`. Solo viven en memoria.
- Tiendas inactivas y precios con fecha futura no entran al libro de precios de S5; solo datos nutricionales `is_preferred = 1`.
- Los puntos de decisión abierta están en el bloqueo `S6-comidas-calculos` (13 puntos PROVISIONAL).

## Bloqueos
- `docs/estado/bloqueos/S6-comidas-calculos.md` (nuevo).

## Pendientes
- Registrar el área en `backend/cmd/finanzas/areas.go` (fuera de mi alcance; ver solicitud).
- `getRecipe`: `packages_to_buy`, `price_store_id` y `price_observed_date` por ingrediente quedan ausentes (opcionales en C1): `RecipeDetailFor` de S4 no los entrega por ingrediente.
- Elecciones de la lista (S4 §8.1) en costos del plan: ver bloqueo, punto 7.

## Propuestas para Axel
- Elegir en el bloqueo los puntos 1, 2, 6, 7 y 8 (afectan lo que se ve en Comidas); los demás pueden quedar como están.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~45 acciones, ~330 KB leídos, ~125 KB escritos.
