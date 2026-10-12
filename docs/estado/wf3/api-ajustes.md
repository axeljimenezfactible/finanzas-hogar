# api-ajustes · WF-3a «Construcción» · 12-oct-2026

Agente api-ajustes (Opus). Alcance: `backend/internal/api/ajustes/`, `backend/internal/dominio/ajustes_*.go` y este reporte. Contrato: C1 v1.4 (extractos `ajustes-a`, `-b`, `-c`, `-f`); esquema 1.4; CONVENCIONES-backend 1.0.

## Hecho
- Las 19 operaciones del reparto (CONVENCIONES §14) registradas con `ajustes.Register`: categorías (lista con `q` sin acentos, `status`, conteos `expense_count` y `recent_usage_count` de 90 días; alta, cambio/archivo con control de ciclos; reorden), personas (lista holder → co_holder → external con `has_account`; cambio de nombre, tope K13 y notas con Merge Patch), reglas de reparto (lista con partes, filtro `is_template`; cambio de nombre, notas y partes con `scope` obligatorio, Q2-68 C), partidas planificadas (lista con filtros `category_id`, `credit_id`, `kind`; alta y cambio con todas las reglas CHECK y del disparador del titular como 422 por campo), importación de `gastos_fijos_ingresos.csv` (vista previa sin escribir; importación todo-o-nada con 422 por renglón y columna), metas (alta y cambio; ahorrado, avance, aportación necesaria y estado **los da `motor_fin.SavingsGoals`**, no se recalculan), proyección (`projection.repeat_cadence_enabled`, K21) y avisos (6 claves `alerts.*` con `value`, `initial_value`, `is_initial`, `parameters` y `reset_parameters` vía `dominio.ApplyParams`).
- Cada alta, cambio o baja escribe `audit_log` en la misma transacción (`InsertOwned`, `UpdateOwned`, `DeleteOwned`, limitados a las tablas del área). Dinero en centavos `int64`; el CSV convierte pesos a centavos sin punto flotante.
- Idempotency-Key: la aplica el middleware en `createCategory`, `createPlannedItem`, `importPlannedItems` y `createSavingsGoal` (probado: falta de clave → 400; repetición → `Idempotent-Replayed: true`).
- PREGUNTA (bloqueos nuevos): `docs/estado/bloqueos/S6-ajustes-api.md` (next_date, scope = all, reorden, órdenes por defecto, metas inactivas, 400 vs 422) y `docs/estado/bloqueos/S6-importar-partidas-csv.md` (correspondencia columna → campo del CSV).
- Las claves `reports.*` son de `getReportSettings`/`updateReportSettings` (s7-pdf, CONVENCIONES §14): no se tocaron.

## Archivos
- `backend/internal/api/ajustes/ajustes.go` (Register y utilidades), `catalogos.go` (categorías, personas, reglas), `partidas.go`, `importar.go`, `metas.go` (metas, proyección, avisos), `ajustes_test.go`.
- `backend/internal/dominio/ajustes_repos.go` (repositorios de category, person, cost_share_rule(+part), planned_item, savings_goal) y `ajustes_repos_test.go`.
- `docs/estado/bloqueos/S6-ajustes-api.md`, `docs/estado/bloqueos/S6-importar-partidas-csv.md`, este reporte.

## Verificación
- `go build` y `go vet` de todos los paquetes: sin errores (mientras comprasa estaba a medias se midió sin él; al final `make verify-wf3-go` pasó completo).
- `gofmt -l .`: vacío.
- `go test ./internal/api/ajustes/... ./internal/dominio/... -count=1 -cover`: `ok internal/api/ajustes 87.1 %`; `ok internal/dominio 63.8 %` (paquete completo, con código de otros agentes aún sin pruebas); `ajustes_repos.go` 195/226 = **86.3 %**.
- `make verify-wf3-go GOOSE=…`: go build/vet sin errores, gofmt limpio, `go mod tidy -diff` sin cambios, `go test ./...` 19 paquetes OK.
- `python3 -I backend/internal/api/c1gen/build.py --check`: «los 23 archivos coinciden byte a byte».
- **S6-02**: `internal/api/ajustes` 87.1 % (sin umbral, K25); mi parte de `internal/dominio` 86.3 % (el 63.8 % del paquete lo bajan archivos de otras áreas aún sin pruebas).
- **S6-03**: 19/19 operaciones con prueba `TestOp_<operationId>` que valida cada respuesta contra `openapi.yaml` (testutil + kin-openapi), 2xx y 4xx:

| Operación | Prueba | Códigos probados |
|---|---|---|
| listCategories | TestOp_listCategories | 200 (cursor, q, status), 400 ×3, 401 |
| createCategory | TestOp_createCategory | 201 (+audit_log), 409 already_exists, 422 ×3, 400 ×2 |
| updateCategory | TestOp_updateCategory | 200, 404, 400 ×2, 422 ciclo, 409 |
| reorderCategories | TestOp_reorderCategories | 204, 422 ×4, 400 |
| listPersons | TestOp_listPersons | 200 (orden y has_account), 400 |
| updatePerson | TestOp_updatePerson | 200 (valor y null), 422 ×2, 400 ×2, 404 |
| listCostShareRules | TestOp_listCostShareRules | 200 (partes, is_template), 400 |
| updateCostShareRule | TestOp_updateCostShareRule | 200 onward, 422 ×4, 409 scope=all, 400 ×2, 404 |
| listPlannedItems | TestOp_listPlannedItems | 200 (filtros), 400 ×2 |
| createPlannedItem | TestOp_createPlannedItem | 201 ×2, 422 ×4, 400 ×3 (fecha, enum, sin Idempotency-Key) |
| updatePlannedItem | TestOp_updatePlannedItem | 200 (null borra), 422 ×2, 400 ×2, 404 |
| previewPlannedItemsImport | TestOp_previewPlannedItemsImport | 200 (5 válidos, 3 con error, sin escribir), 415 ×2, 422 ×3, 400 |
| importPlannedItems | TestOp_importPlannedItems | 422 por renglón, 201 (separador «;»), repetición idempotente, 415 |
| getProjectionSettings | TestOp_getProjectionSettings | 200, 401 |
| updateProjectionSettings | TestOp_updateProjectionSettings | 200, 422 ×2, 400 ×2 |
| createSavingsGoal | TestOp_createSavingsGoal | 201 (cifras del motor), 409, 422 ×2, 400 ×2 |
| updateSavingsGoal | TestOp_updateSavingsGoal | 200 (ahorrado 1 500 000, 500 décimas), 422 ×2, 400, 404, 409 |
| getAlertSettings | TestOp_getAlertSettings | 200 (80/100/[1,0]/09:00/1/7), 401 |
| updateAlertSettings | TestOp_updateAlertSettings | 200 cambiar+restablecer, 422 ×5 (rango, null, ambos, clave ajena), 400 ×5 |

- **S6-04**: error único de C1 (`code` estable en inglés, `message` es-MX, `fields`, `request_id`); páginas de 1 a 200 con cursor opaco (`limit=201` → 400); fechas AAAA-MM-DD e instantes −06:00 con `api.APITime`; dinero en centavos; Idempotency-Key en las 4 altas.
- **S6-15**: identificadores en inglés; mensajes al usuario en español de México; nada imprime contraseñas, cookies ni tokens (no se registra nada propio).

## Supuestos (PROVISIONAL, con bloqueo)
- `next_date` se omite; `scope = all` → 409; reorden con todas las activas; órdenes por defecto; metas inactivas sin cifras; `{}` y enums → 400 (S6-ajustes-api).
- Correspondencia del CSV, tope de 8 MiB y `text/csv` estricto (S6-importar-partidas-csv).
- `groceries_category`/`paper_category` se pasan como «groceries»/«paper» al motor solo porque `Params` los exige; `SavingsGoals` no los usa.

## Bloqueos
- `docs/estado/bloqueos/S6-ajustes-api.md` (ABIERTO), `docs/estado/bloqueos/S6-importar-partidas-csv.md` (ABIERTO). Depende de `C1-reparto-alcance` (ABIERTO).

## Pendientes
- `backend/cmd/finanzas/areas.go`: agregar `ajustes.Register,` y el import `"finanzas-hogar/internal/api/ajustes"` (fuera de alcance; lo hace el integrador).
- Con aprobación: envoltorios exportados en `motor_fin` (`SplitRule`, `NextOccurrence`) para cerrar los puntos 1 y 2 de S6-ajustes-api.
- Cobertura ≥ 85 % del paquete `internal/dominio` completo depende de las pruebas de las otras áreas.

## Propuestas para Axel
- Decidir `C1-reparto-alcance` (propuesta B) y autorizar los envoltorios de `motor_fin` sin cambio de cifras.
- Aprobar la correspondencia del CSV (opción A) o pedir la columna `code` en `cost_share_rule` (opción B).

## Paquetes instalados
Ninguno.

## Consumo aproximado
~48 acciones; ~200 KB leídos; ~125 KB escritos.
