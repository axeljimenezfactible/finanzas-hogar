# C1 · Contrato de datos y API · Finanzas Hogar

| | |
|---|---|
| Versión | **C1 v1.5** (12-oct-2026, tablero de control K54; revisión 1.5.1 del mismo día) · `info.version` 1.5.1 de `docs/contratos/openapi.yaml` |
| Autor | Agente S6a-contrato-c1 (flujo WF-0 «Cimientos», fase S6 parte a); corrección r1: agente correccion-C1-r1; corrección r2: agente correccion-C1-r2; v1.1 corrección r1: agente correccion-C1-r1; v1.1 corrección r4: agente correccion-C1-r4; v1.2: agente C1-v1.2; v1.2 corrección r1: agente correccion-C1-r1; v1.2 ronda 3 (P1d): agente correccion-C1-r3; v1.3: agente C1-puerta2 (P2c); v1.4: agente k53-aplicar (P2d); v1.5: agente s1-c1-v15 (WF-3c); v1.5.1: agente backend-cierre (WF-3c bis) |
| Estado | K54 aplicada (C1 v1.5): operaciones del tablero de control de U3 2.0 §2.6 sobre el esquema 1.5 de S1, con lo no decidido PROVISIONAL (bloqueos `C1-tablero` y `U3-tablero-*`). Puerta 2 bis aplicada (C1 v1.4). Puerta 2 aplicada (C1 v1.3): respuestas de Axel de `docs/estado/puerta2/RESPUESTAS.md` sobre el esquema 1.4 de S1. Las preguntas sin recomendación quedan con la opción en uso y no cambian C1 |
| Fuentes | `plan.md` §1, §3.4 (S6, S7), §3.5 (U1), §3.7; `decisiones.md`; esquema de S1, versión 1.5 (00001 a 00005: `backend/db/migrations/00001_esquema.sql`, `00002_puerta1.sql`, `00003_puerta1b.sql`, `00004_puerta2.sql` y `00005_tablero.sql`; `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md`); U1 1.2 (`docs/design/U1-arquitectura.md`, necesidades N-xxx); puerta 2: `docs/estado/puerta2/PUERTA-2.md` y `RESPUESTAS.md`, contratos de motores S2 1.3, S4 1.5 y S5 1.5 (`docs/contratos/motores/`) y U3 1.3 |
| Archivo de máquina | `docs/contratos/openapi.yaml` (OpenAPI 3.1.1): {{N_RUTAS}} rutas, {{N_OPS}} operaciones, {{N_ESQ}} esquemas |
| Generador | Este documento, `openapi.yaml` y `backend/internal/api/routes.go` son archivos generados: no se editan a mano. Fuente: `backend/internal/api/c1gen/` (este texto sale de `c1_plantilla.md`; ver su `README.md`). Regenerar: `python3 -I backend/internal/api/c1gen/build.py`. Comprobar que nadie los editó a mano: `python3 -I backend/internal/api/c1gen/build.py --check`. Verificar C1-01 a C1-15, la matriz §10 y S6-03: `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py` |

Este documento fija **cómo** se piden los datos que U1 dice que cada pantalla **necesita**. El OpenAPI es la fuente de verdad de rutas y esquemas; aquí están las convenciones únicas, las métricas con su fórmula, la matriz pantalla → necesidad → endpoint y la cobertura de S1. Todo lo que no sale del plan, de las decisiones o del oráculo va marcado **PROVISIONAL** y apunta a su bloqueo en `docs/estado/bloqueos/`.

## Cambios

| Versión | Fecha | Cambio |
|---|---|---|
| C1 v1.0 | 10-oct-2026 | Primera versión (WF-0). Los cambios posteriores se publican como C1 v1.1, v1.2… (plan §3.7) con esta tabla actualizada y `info.version` igual. |
| C1 v1.0 · corrección r1 | 10-oct-2026 | Sin cambios de rutas, esquemas, convenciones ni matriz. El generador y el verificador pasan al repositorio (`backend/internal/api/c1gen/`) y las cabeceras de `openapi.yaml` y `routes.go` lo citan (hallazgo H-05 de la revisión r1). |
| C1 v1.0 · corrección r2 | 11-oct-2026 | Regenerado contra el esquema 1.1 de S1 (hallazgos R2-01 y R2-02 de la revisión r2; sigue siendo la entrega de WF-0 para la puerta 1). Lista de compras guardada: `generateShoppingList` (N-127), `createShoppingListItem` (N-152), `updateShoppingListItem` (N-151, N-153), `assignShoppingListStore` (N-165) y `getShoppingList` con la lista guardada y el filtro `is_bought` (N-150, N-073), todo PROVISIONAL por la opción A de `C1-lista-compras`. Presupuesto con vigencia (`valid_from_date`, `valid_to_date`; sin `period_id`) y 409 `budget_overlap`. Tiempos aptos de la receta como lista `meal_slot_ids` sobre `recipe_meal_slot` (sin `default_meal_slot_id`). Descripción del gasto opcional sin texto de relleno. Movimientos de inventario guardados en `inventory_movement`, con motivo en los ajustes. Ingreso esperado con `second_half_amount_cents`. `idempotency_key` en `x-tablas` de los POST de captura. Tope por quincena sin cifra. Columna nueva `category.kind`. |
| C1 v1.1 (puerta 1) | 11-oct-2026 | Aplica las respuestas de Axel de la puerta 1 (`decisiones.md` §K) y U1 1.1 contra el esquema 1.2 de S1 (migraciones 00001 y 00002). Detalle en «C1 v1.1 (puerta 1)», abajo. `info.version` = 1.1.0. |
| C1 v1.1 · corrección r1 | 11-oct-2026 | Sin cambios de rutas, parámetros, tipos ni cifras; solo textos (hallazgos H1 a H5 de `docs/estado/p1/REVISION-C1-r1.md`). §7.1 nueva: efecto de cada `scope` de `updateBudget` sobre las vigencias, con un ejemplo de cada uno (K2). §1 punto 8: ciclo de una captura de gasto o de peso antes y después de las 24 h, con ejemplo de reintento (K9). El borrador del ticket no crea gasto, precio ni inventario (K3). Las descripciones de los instantes ya no dicen «(UTC)», que contradecía la salida con el desplazamiento de America/Mexico_City (R2-06, K12). `info.version` sigue en 1.1.0. |
| C1 v1.1 · corrección r2 | 11-oct-2026 | Hallazgos H-r2-01 y H-r2-02 de `docs/estado/p1/REVISION-C1-r2.md`. La semana de la lista de compras ya no se llama `period`: es la llave `week` (`components/parameters/ListWeek`, semana ISO, obligatoria, sin valor por omisión) en `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets` y `getDashboardBestBasket`; el filtro `period` vuelve a tener una sola definición, §6 explica las llaves que no son filtros y la matriz recurso × filtro suma la columna de `week` (C1-06). `scope=onward` de `updateBudget`: el monto nuevo rige desde `month` sin fin y elimina los renglones posteriores de la misma categoría; §7.1 lo dice y trae el ejemplo, y `Budget.valid_to_date`, `BudgetCreate`/`BudgetUpdate` y `createBudget` dicen lo mismo (C1-16, K2). Cambia el nombre de un parámetro de consulta en 7 operaciones; sin cambios de rutas, cuerpos ni cifras. `info.version` sigue en 1.1.0. |
| C1 v1.1 · corrección r4 | 11-oct-2026 | Hallazgos H-r3-01, H-r3-02 y H-r3-04 de `docs/estado/p1/REVISION-C1-r3.md`. `getBasketComparison` recibe la llave `week` (`ListWeek`, obligatoria; 422 `validation_failed` si falta) y ya no recibe el filtro `period`; su descripción dice qué compara con `basket=shopping_list` y con `basket=expense`; la matriz de §6 y la tabla «Llaves que no son filtros» la incluyen. El valor por omisión de `getMealPlanWeek` (los domingos, la semana siguiente) queda `x-provisional` con el bloqueo `U1-lista-vigente-domingo`, en §13 y en N-120 de §10. `Budget.valid_to_date` toma el texto completo del comentario de S1 (que ahora nombra `onward` y `this_month`) y se quitó el parche de `gen2_schemas.py`. `verify_c1.py`: C1-06 exige `week` a toda operación que lee `shopping_list` y C1-15 marca la regla del domingo sin marca |
| C1 v1.2 (puerta 1b) | 11-oct-2026 | Aplica K35 a K44 (`decisiones.md` §K) y la sección «Para C1» de `docs/estado/p1c/U1-1.2.md` contra el esquema 1.3 de S1 (migraciones 00001, 00002 y 00003). Parámetros configurables (K44): cada clave de `app_setting` del inventario de S1 se lee con su valor vigente y su valor inicial y se cambia o restablece por clave; tabla generada en §11.7. Operaciones nuevas: `getSecuritySettings` y `updateSecuritySettings` (N-320), `getAlertSettings` y `updateAlertSettings` (N-323), `dismissTelegramChatCandidate` (N-324), `getDataSettings` y `updateDataSettings` (N-325), `getHealthSettings` y `updateHealthSettings` (N-326); N-321 y N-322 entran en las de comparador y precios; N-327 en `createHealthProfile` y `updateHealthProfile`. K36: `HealthProfile` con el ajuste por persona y el de la casa. K38, K39 y K43: sin marcas (N-023, N-042, N-064, N-120, N-226, N-230, N-297, N-300, N-311 y N-313 FIRMES en §10). `info.version` 1.2.0. Detalle en «C1 v1.2 (puerta 1b)», abajo. |
| C1 v1.2 · corrección r1 | 11-oct-2026 | Hallazgo A-01 de `docs/estado/p1c/REVISION-C1-r1.md` (C1-23, K43). `params_json` de `Notification` y `NotificationUpdate` declara las propiedades `send_time` (`type: string`, patrón `HH:MM` de 24 horas) y `advance_days` (`type: integer`, `minimum: 0`), sin valor por omisión; `weekday` y `day_of_month` siguen por `additionalProperties: true` y la marca de lugar `x-provisional` (bloqueo `S1-parametros-avisos`) se queda. Detalle al final de §11.7. `verify_c1.py` comprueba C1-23. Sin cambios de rutas, operaciones, cifras ni versión (`info.version` 1.2.0). |
| C1 v1.2, ronda 3 (P1d) | 11-oct-2026 | Hallazgos bajos N-01, N-02, N-03, B-02, B-03, B-04, B-05 e I-01 de `docs/estado/p1c/REVISION-C1-r2.md`. N-01: una sola constante de versión (`info.version` 1.2.0) para los encabezados generados de `openapi.yaml`, `routes.go` y los extractos. N-02: `minimum: 0` en `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` de `HealthProfileCreate` y `HealthProfileUpdate` (K36). N-03: el estado de la tarjeta (N-205) cita RN-26 y S3-08 (K18) en vez de un archivo de bloqueo. B-02: `NotificationUpdate` ya no lleva la marca de lugar en el esquema entero (`is_enabled` es FIRME); la llevan `params_json` y sus claves `send_time` y `advance_days`, en `Notification` y en `NotificationUpdate` (bloqueo `S1-parametros-avisos`). B-03: la descripción de `params_json` dice que los umbrales y los días antes del pago viven en `app_setting` desde 1.3 y que ahí quedan `send_time`, `advance_days`, `weekday` y `day_of_month`, con valor FIRME (decisión F) en las dos últimas. B-04: el mensaje de `verify_c1.py` sobre la regla del domingo dice FIRME (K38). B-05: el reporte de C1 v1.2 trae los comandos completos de verificación. I-01: esta cabecera dice puerta 1b aplicada, esquema 1.3 y U1 1.2. La nota de `send_time` en la 00003 dice que lo propone S7 en WF-3 (K31, puerta 4). Sin cambios de rutas, operaciones, cifras ni versión (`info.version` 1.2.0). |
| C1 v1.3 (puerta 2) | 11-oct-2026 | Aplica las respuestas de Axel de la puerta 2 (`docs/estado/puerta2/RESPUESTAS.md`: lo recomendado de `PUERTA-2.md` más Q2-11 B) contra el esquema 1.4 de S1 (migraciones 00001 a 00004). Reportes y avisos: hora, día y anticipación pasan a `app_setting` con lectura y cambio por clave (operaciones nuevas `getReportSettings` y `updateReportSettings`; `AlertSettings` suma tres claves); `params_json` ya no declara `send_time` ni `advance_days`. Rango por clave en el cambio (`minimum`, `maximum`, `pattern`) y tipo nuevo `SettingParameterTime`. `Credit.missing_terms` con `vat` (5 datos) y `credit.vat_included`. Alcance `scope` en `CostShareRuleUpdate` (§7.2). Lista de compras con tienda asignada, fecha de visita y totales de Q2-42 A; comparativa con `cheapest_store_id`; acumulados de T-2; décimas del avance de una meta; segunda tienda de la tarjeta 5; tienda del escenario. Detalle en «C1 v1.3 (puerta 2)», abajo. `info.version` = 1.3.0. |
| C1 v1.4 (puerta 2 bis) | 12-oct-2026 | Aplica K53 (`docs/plan/decisiones.md` §K; puerta 2 bis): la aplicación no lee sitios de tiendas y un agente externo mete productos, presentaciones, precios y recetas por la API. Cuatro operaciones nuevas sin N-xxx: `getExternalAgentReference`, `importExternalCatalog`, `importExternalPrices` e `importExternalRecipes`, con la etiqueta «Compras y tiendas» y la marca `x-agente-externo`. Esquema de seguridad nuevo `serviceToken` y esquemas `External*`. La lectura asistida queda «no aplica desde K53» en `updateStore`, `updatePriceSettings` y sus esquemas. Remanentes de P2c: `SavingsGoal.status` según Q2-18 A (I-2) y `getCashFlowFortnights` según S2 1.3 §3.11 (I-6). Detalle en «C1 v1.4 (puerta 2 bis)», abajo. `info.version` = 1.4.0. |
| C1 v1.4 · corrección r1 (P2d) | 12-oct-2026 | Hallazgos M-1 y M-2 de `docs/estado/p2c/REVISION-k53-r1.md` (agente `correccion-k53-r1`). `info.version` sigue en 1.4.0; no cambian rutas, parámetros ni tipos, solo textos. M-1: el tamaño máximo de los lotes del agente externo (500 productos, 500 precios y 100 recetas; 422 fuera de rango) era un supuesto sin fuente; ahora cita la pregunta (g) del bloqueo `S5-agente-externo` en `ExternalCatalogImport`, `ExternalPriceImport`, `ExternalRecipeImport` e `importExternalCatalog`. M-2: la prosa de `balance_cents` en `CashFlowFortnights` resta también el gasto variable registrado, como S2 1.3 §3.11 (`spent_cents` incluye `variable_spent_cents`), y aclara que ese dato no viene en la fila. |
| C1 v1.5 (K54, tablero) | 12-oct-2026 | Aplica K54 (`docs/plan/decisiones.md` §K) con las operaciones que pide `docs/design/U3-dashboard.md` 2.0 §2.6, sobre el esquema 1.5 de S1 (migración `00005_tablero.sql`). Cinco operaciones nuevas (`getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSettings`, `updateDashboardSettings`), 14 campos nuevos opcionales en `DashboardMonthVsPlan`, `IncomeMonth` y `CreditTotals`, cinco claves de `app_setting` con `x-app-setting`, las llaves `from` y `to` y las métricas M-18 a M-22. Detalle abajo |
| C1 v1.5.1 (escala de G-6) | 12-oct-2026 | Hallazgo M-R3-1 de `REVISION-especificacion-r3.md`, opción B (decisión 3 del orquestador de WF-3c bis): `DashboardMonthlyOverview` suma `scale_max_pct`, entero requerido, la misma escala de `DashboardDebtCapacity`. Textos: `has_credits` cuenta las líneas archivadas con pagos en el mes (bloqueo `S6-tablero-lineas-archivadas`, opción a1), `BudgetSignals` dice qué pasa con un `period` posterior al mes en curso (bloqueo `S6-tablero-senal-periodo`, opción a), y `CreditTotals` y `DebtTrend` dicen qué compras y qué líneas cuentan (bloqueos `S6-deuda-viva-hoy` y `S6-tablero-lineas-archivadas`, punto 5). No cambia rutas ni operaciones. Detalle abajo |

### C1 v1.5.1 (escala de G-6)

Agente `backend-cierre` (flujo WF-3c bis, 12-oct-2026). Corrige el medio M-R3-1 de `docs/estado/wf3/REVISION-especificacion-r3.md` con la opción B, que eligió el orquestador: U3 G-6 dibuja el eje Y «de 0 % al mayor entre `scale_max_pct` y el mayor `debt_ratio_pct` a la vista», y ese campo solo venía en `getDashboardDebtCapacity`.

- **Campo nuevo requerido:** `DashboardMonthlyOverview.scale_max_pct`, el múltiplo de 10 igual o mayor que 1.5 × `limit_pct` (60 con 40). Lo calcula la misma función del servidor que el de `getDashboardDebtCapacity`, así que G-1, G-6 y B-1 leen la misma escala sin pedir B-1 ni calcular nada. PROVISIONAL, bloqueo `U3-tablero-barometro`.
- **Textos:** `DashboardDebtCapacity.has_credits` ahora dice que cuenta una línea archivada con pagos en el mes (`S6-tablero-lineas-archivadas`, opción a1, que ya rige en el código). `BudgetSignals` dice que con un `period` posterior al mes en curso solo cuentan los meses cerrados, así que `months` puede venir más corto o vacío y la señal es `none` (`S6-tablero-senal-periodo`, opción a).
- **Textos de la deuda viva y los pagos del mes** (corrección r1, agente `correccion-backend-cierre-r1`, hallazgos M-1 y M-2 de `docs/estado/wf3/REVISION-backend-cierre-r1.md`): `CreditTotals.live_debt_cents` es la deuda de `getLiveDebt` hoy y `DebtTrend.months[].live_debt_cents` la de `getLiveDebt` al cierre, con solo las compras con fecha hasta esa fecha (`S6-deuda-viva-hoy`, puntos 1 y 2). Las dos y los pagos del mes de `CreditTotals` cuentan todas las líneas, activas o archivadas, como el barómetro (`S6-tablero-lineas-archivadas`, punto 5). `month_pending_household_cents` dice su regla: la parte del hogar menos los abonos del hogar, nunca negativa; el abono de Karen no la cubre; un corte de tarjeta pendiente cuenta su parte completa (`S6-deuda-viva-hoy`, punto 3). Todo PROVISIONAL.
- Sin cambios de rutas, operaciones, parámetros ni tablas: 221 operaciones y 176 rutas, como en 1.5.0.

### C1 v1.5 (K54, tablero de control)

Agente `s1-c1-v15` (flujo WF-3c, 12-oct-2026). Fuentes: K54 (palabras de Axel, decisiones 1 a 4 y propósito del sistema) y U3 2.0 §2.6 (agente `u3-v2`). Esquema de S1 1.5: la `00005_tablero.sql` solo siembra cinco claves de `app_setting` y sus dos disparadores de rango; no cambia tablas.

- **Operaciones nuevas de la etiqueta «Inicio»** (la interfaz no calcula cifras; todo sale de S2 y S3):
  - `getDashboardDebtCapacity` (GET `/api/v1/dashboard/debt-capacity`, `period` opcional): barómetro de capacidad de endeudamiento, M-20 y RN-30.
  - `getDashboardMonthlyOverview` (GET `/api/v1/dashboard/monthly-overview`, llaves `from` y `to` obligatorias, máximo 24 meses): G-1 y G-6, M-21.
  - `getBudgetSignals` (GET `/api/v1/budget/signals`, `period` opcional): sugerencia «subir o recortar presupuesto», M-22 y RN-31.
- **Parámetros del tablero** (etiqueta «Ajustes», K44): `getDashboardSettings` y `updateDashboardSettings` (`/api/v1/settings/dashboard`) con `debt.capacity_warning_pct` = 30 y `debt.capacity_limit_pct` = 40 (FIRMES, K54 decisión 2), `dashboard.payment_reminder_count` = 5, `budget.signal_months` = 3 y `budget.signal_margin_pct` = 10. Los umbrales llegan también en las respuestas del barómetro y del resumen mensual (`warning_pct`, `limit_pct`, con `x-app-setting`).
- **Campos nuevos** (opcionales, para no romper P-06, P-08 ni P-17): `DashboardMonthVsPlan` suma `cash_cents`, `previous_month`, `previous_actual_cents`, `previous_cash_cents`, `actual_difference_cents`, `cash_difference_cents`, `categories_warning_count` y `categories_exceeded_count`; `IncomeMonth` suma `previous_month`, `previous_fortnights_cents`, `previous_extras_cents` y `fortnights_difference_cents`; `CreditTotals` suma `month_payment_count` y `month_pending_household_cents`. Su fórmula, periodo y ejemplo del oráculo están en la fila M-01 de §9, ampliada para K-1 (corrección r2, C1-09).
- **Lo PROVISIONAL** (bloqueo `C1-tablero`): la sección «Presupuesto» de las cinco claves, sus rangos (de 1 a 100 los umbrales con el amarillo menor que el tope, de 1 a 20 pagos, de 2 a 12 meses y de 1 a 50 puntos), las necesidades de U1 de las cinco operaciones (cuelgan de N-022, N-028, N-027, N-275 y N-278 hasta que U1 agregue las suyas), el nombre `dashboard.payment_reminder_count` y que `period` acepte solo `AAAA-MM`. Lo que mide el barómetro, N y la regla de la sugerencia siguen los bloqueos `U3-tablero-barometro`, `U3-tablero-recordatorio` y `U3-tablero-senal-presupuesto`.
- **Error de validación.** U3 2.0 §2.6 escribe «400 `validation_failed`» para una ventana inválida; C1 usa 422 `validation_failed` para un parámetro bien formado con valor inválido (§3), y así queda.

### C1 v1.4 (puerta 2 bis)

Agente `k53-aplicar` (flujo P2d, 12-oct-2026). Fuente: `docs/estado/puerta2/RESPUESTAS.md`, sección «Puerta 2 bis», registrada como K53 en `docs/plan/decisiones.md`. No cambia el esquema de S1: sigue el 1.4, migraciones 00001 a 00004.

- **Agente externo (K53).** La aplicación no lee sitios de tiendas. Una sesión de Claude Code de Axel con Claude in Chrome obtiene precios y productos, y puede agregar recetas. Los mete con estas operaciones de la etiqueta «Compras y tiendas», que no tienen N-xxx porque no tienen pantalla (§11.2):
  - `getExternalAgentReference`: ids de tiendas, grupos de producto y tiempos de comida.
  - `importExternalCatalog`: productos y presentaciones por nombre.
  - `importExternalPrices`: precios con tienda, fecha y `external_ref`, que se guarda en `price.source_ref`.
  - `importExternalRecipes`: recetas en borrador.
- **Idempotencia y resultado.** Las tres importaciones son idempotentes por llave natural y por `Idempotency-Key`. Responden `ExternalImportResult`, renglón por renglón.
- **Lo PROVISIONAL del agente externo** (bloqueo `S5-agente-externo`):
  - `source` = `site_assisted`, que la documentación llama «agente externo», y esos precios compiten como `manual`.
  - Las recetas entran como borrador, con origen «external_agent». El esquema 1.4 todavía no guarda ni el borrador ni ese origen.
  - El agente no escribe inventario.
  - Tamaño máximo de un lote: 500 productos, 500 precios o 100 recetas; un lote vacío o más grande responde 422. Es supuesto de `k53-aplicar`, no decisión de Axel (bloqueo `S5-agente-externo`, pregunta g).
- **Credencial de servicio.** Esquema de seguridad nuevo `serviceToken` (`Authorization: Bearer`), solo en las operaciones con `x-agente-externo`. El esquema 1.4 no tiene ninguna credencial de servicio que reutilizar. Dónde vive y a nombre de quién quedan `audit_log` e `idempotency_key` es PROVISIONAL (bloqueo `S5-agente-externo`, pregunta e).
- **Lectura asistida retirada (K53).** Ningún veredicto la habilita, así que `true` responde siempre 409 `assisted_reading_not_allowed` en `updateStore` y en `updatePriceSettings`. `Store`, `PriceSources`, `PriceSettings` y `PriceSettingsUpdate` la describen como «no aplica desde K53». Las columnas siguen en el esquema 1.4 hasta que S1 las retire.
- **Remanentes de P2c** (`docs/estado/p2c/INTEGRACION.md` §5):
  - I-2: `SavingsGoal.status` pasa a `in_progress`, `achieved` y `overdue`, con la regla de S2 1.3 §3.13 (Q2-18 A).
  - I-6: en `getCashFlowFortnights`, `planned_savings_cents` deja de ser obligatoria. S2 la entrega en null y es PROVISIONAL (bloqueo `S2-ahorro-quincenal`). `balance_cents` ya no resta el ahorro planificado ni el gasto variable estimado (S2 1.3 §3.11).

### C1 v1.3 (puerta 2)

- **Entradas**: el generador aplica `00001` a `00004_puerta2.sql` (esquema 1.4). `parse_schema.py` lee además `@param_rango` (rango aprobado que se suma a su clave) y `@param_sin_valor_cerrado` (un `@param_sin_valor` que ya tiene clave), y toma la última `@col_note` de cada columna. El rango de cada clave no se copia: se prueba contra la base temporal (cada valor con un `UPDATE` que se deshace), así `minimum`, `maximum` y el patrón de la hora salen de los disparadores de S1 y el generador se detiene si el patrón no coincide con la base.
- **Reportes y avisos (Q2-11 B, K47; Q2-72 y Q2-73, K50)**: `getReportSettings` y `updateReportSettings` (`/api/v1/settings/reports`, N-296 y N-297) leen y cambian por clave `reports.daily_send_time` (21:30), `reports.weekly_weekday` (7) y `reports.weekly_send_time` (19:00), `reports.monthly_day_of_month` (1; en un mes más corto sale el último día, Q2-66 B) y `reports.monthly_send_time` (08:00). `AlertSettings` suma `alerts.send_time` (09:00), `alerts.stale_prices_advance_days` (1) y `alerts.meal_plan_missing_weekday` (7), y sus operaciones cubren además N-298 y N-300. `Notification.params_json` y `NotificationUpdate.params_json` ya no declaran `send_time` ni `advance_days` y pierden la marca de lugar: el bloqueo `S1-parametros-avisos` está resuelto. `updateNotification` sigue para activo y parámetros propios (hoy solo `scopes` de `budget_threshold`).
- **Rangos (Q2-74, Q2-77, Q2-66 B)**: cada clave del `PATCH` declara su rango (§11.7, columna «Rango (cambio)»): día de la semana de 1 a 7, día del mes de 1 a 31, `comparator.timeout_seconds` de 1 a 30, `price.stale_after_days` de 1 o más, horas con el patrón `HH:MM` de 00:00 a 23:59 y los demás enteros de 0 o más. `price.zero_weight_after_days` ≥ `price.stale_after_days` la valida la API (422) y el motor de precios. Tipo nuevo `SettingParameterTime`.
- **Otras claves (Q2-75, Q2-76)**: `optimizer.date_window_days` (7) en `parameters` de `getOptimizerSettings` y `updateOptimizerSettings`, y deja de ser campo propio de esos dos esquemas (sigue en `ScenarioRequest` y en `constraints_used`); `ocr.timeout_seconds` (30) en `DataSettings` (sección propuesta por S1).
- **IVA de un crédito (Q2-67 opción 1, K46)**: `Credit`, `CreditCreate` y `CreditUpdate` traen `vat_included` (`true` = la tasa ya incluye IVA; `false` = se suma `vat_on_interest_bp`; ausente = falta capturar). `Credit.missing_terms` lista 5 datos: `annual_rate`, `rate_base`, `vat`, `min_payment_pct_bp` y `min_payment_fixed_cents`.
- **Reparto (Q2-68 C)**: `CostShareRuleUpdate.scope` (`onward` o `all`), obligatorio con `parts` y sin valor por omisión (§7.2).
- **Lista de compras (Q2-39 a Q2-42 A, S4 1.5)**: `ShoppingListItem.assigned_store_id` (tienda asignada, rige la agrupación); `recommended_store_id` es la recomendación de S5 por precio y en vertical `market` no rige la agrupación; `estimated_price_cents` es el costo con la cotización que rige; `groups[].store_id` es la tienda asignada (ausente = «sin mercado dado de alta»), `groups[].visit_date` es nuevo, y los totales suman solo los artículos con precio. Orden de los grupos: mercado, cadena y el grupo sin tienda al final.
- **Comparativa (S5 1.5, §7.4)**: `BasketComparison.cheapest_store_id` (primera tienda completa; ausente si ninguna) y `stores[].difference_vs_cheapest_cents` solo en tiendas completas; `stores[]` lleva primero las completas; `identical_count` cuenta renglones.
- **Inicio, metas y comparador**: `AccrualVsCash.months[]` suma `cumulative_accrued_cents` y `cumulative_cash_cents` (Q2-65 A); `SavingsGoal.progress_tenths_pct` (Q2-17 C); `DashboardBestBasket.second_store_id` (Q2-86 A); `Scenario.store_id` (Q2-33 (b)).

### C1 v1.2 (puerta 1b)

- **Entradas**: el generador aplica `00001_esquema.sql`, `00002_puerta1.sql` y `00003_puerta1b.sql` (esquema 1.3). `parse_schema.py` lee además las anotaciones `@param`, `@param_sin_valor` y `@col_note` de las migraciones y comprueba que cada valor inicial anotado sea el que la migración siembra en `app_setting`.
- **Parámetros (K44)**: lectura y cambio **por clave** de `app_setting` en la operación de settings de su sección: `parameters` (valor vigente, valor inicial, tipo y unidad) y, al escribir, `parameters` y `reset_parameters`. «Restablecer» es un campo, no una operación. Detalle y tabla en §11.7.
- **`PriceSettings`**: `stale_after_days` deja de ser un campo propio y pasa a `parameters` (`price.stale_after_days`), junto a `price.zero_weight_after_days` y `price.profeco_import_max_gb`. `PriceSettingsUpdate` ya no exige `assisted_reading_enabled`.
- **Salud (K36)**: `HealthProfile`, `HealthProfileCreate` y `HealthProfileUpdate` traen `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` (ausente o `null` = el de la casa); `HealthProfile` suma `household_goal_adjust_lose_kcal` y `household_goal_adjust_gain_kcal` y `suggested_kcal_target` sigue RN-17 1.3.
- **Valores de K43**: adjunto de 10 MB (`attachment.max_size_mb`, 413 con el máximo vigente en el mensaje), comparador de 5 s (`comparator.timeout_seconds`, 503 `unavailable`), claves `send_time` y `advance_days` de `params_json` y «casillas sin armar», sin marca.
- **Contraseña**: `PasswordChangeRequest.new_password` ya no fija `minLength` 12 en el contrato: el mínimo es el valor vigente de `auth.password_min_length` (12 al sembrar).
- **Otros**: N-289 sale de `updatePerson` (el tope por persona es de N-310, K13); `Notification.params_json` dice que los umbrales y los días antes del pago próximo pasaron a `app_setting`.

### C1 v1.1 (puerta 1)

- **Entradas**: el generador aplica todas las migraciones de S1 en orden (`00001_esquema.sql` y `00002_puerta1.sql`) en una base temporal y lee de ahí tablas, columnas y llaves; los comentarios de columna salen de los dos archivos.
- **Operaciones nuevas**: `saveTicketScanDraft` (N-075, K3), `listTelegramChatCandidates` (N-299) y `confirmTelegramChatCandidate` (N-316, K4), `getProjectionSettings` (N-318) y `updateProjectionSettings` (N-319, K21). N-317 (cola sin conexión) es Local.
- **Operación retirada**: `createNotificationRecipient` (alta manual del chat). K4 eligió que el bot detecte el chat y la persona lo confirme.
- **Presupuesto (K2)**: `updateBudget` recibe `scope` (`this_month` | `onward`, por omisión `onward`) y `month`; `Budget` suma `is_single_month` y `following_month_amount_cents` (N-275). Efecto de cada alcance sobre las vigencias y un ejemplo de cada uno en §7.1.
- **Ticket (K3)**: `TicketScan` trae `has_draft`, `draft` (`TicketScanDraft`), `display_store_id` y `display_total_cents` (N-071, N-076). `saveTicketScanDraft` solo guarda el borrador: no crea gasto, precio ni inventario; solo `confirmTicketScan` guarda precios (decisiones E: el OCR siempre con confirmación humana) y vacía el borrador.
- **Captura sin conexión (K9)**: `createExpense` y `createWeightEntry` exigen `Idempotency-Key`, aceptan `captured_at` opcional y responden **200 con el mismo recurso** al reintento (`Idempotent-Replayed: true`). La clave queda en `expense.idempotency_key` y `weight_entry.idempotency_key`, sin purga, y no se devuelve en las lecturas. Ciclo completo antes y después de las 24 h y ejemplo de reintento en §1 punto 8.
- **Tasas (K12)**: `rate_decimals` en `CreditTermRate`, `CreditPurchase` y `Promotion`, y `annual_rate_decimals` en `Credit`; `rate_pct` y `annual_rate_pct` muestran la tasa como se capturó (la fracción redondeada ROUND_HALF_UP a esos decimales). Los cálculos usan la fracción.
- **Tope por persona (K13)**: `person.fortnight_payment_cap_cents`, opcional, se lee y edita en `listPersons` y `updatePerson`.
- **Restricciones duras (K16)**: `OptimizerSettings` las devuelve aunque falten, con `missing_constraints`; `ScenarioRequest` las sobreescribe solo para esa comparación; `ScenarioResult` suma `status` (`ok`, `missing_constraints`, `no_valid_scenario`), `missing_constraints` y `constraints_used`, y sin restricciones no recomienda.
- **Tarjeta y crédito de tienda (K18)**: `Credit` suma `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents`, `annual_rate_decimals` y el calculado `missing_terms` («falta capturar»).
- **Proyección (K21)**: `CashFlowFortnights`, `AccrualVsCash` y `CashFlowAssumptions` dicen el modo usado (`projection_mode`).
- **Meta de kcal (K23)**: `HealthProfile` suma `goal` y el calculado `suggested_kcal_target`; confirmar la meta = `updateHealthProfile` con `daily_kcal_target`.
- **Familia (K5)**: `family` con tres valores incluyentes. **Lista de compras (R3-01)**: la semana va en la llave `week` (semana ISO, obligatoria; no es el filtro `period`, corrección r2) en `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getDashboardBestBasket` y `getBasketComparison`; el contrato no tiene semana de lista por omisión.
- **Precios (K22)**: `PriceSettings` pierde `stale_weight`; la regla de peso de un precio viejo quedó fija en K22.
- **Marcas**: se quitan las marcas provisionales que cerraron K1, K2, K3, K4, K5, K6, K7, K8, K11, K12, K13, K15, K16, K17, K20, K21 y K22. Las que quedan citan un bloqueo abierto (§13).
- **Extractos**: `build.py` genera `docs/contratos/extractos/` (un archivo por etiqueta, `componentes-comunes.yaml` e `INDICE.md`).

## 1. Convenciones generales

1. **Base**: `https://finanzas.local:8443` (solo red local, decisión H). Recursos bajo `/api/v1`; fuera de él solo `GET /health` y `GET /openapi.yaml`.
2. **Formato**: JSON UTF-8 (`application/json`), salvo adjuntos y PDF (binario) y cargas de archivo (`multipart/form-data`).
3. **Nombres**: rutas, parámetros y propiedades en inglés corto, **iguales al `nombre_codigo` del diccionario de S1**. Descripciones y mensajes en español de México.
4. **Extensiones del contrato**: `x-tablas` (tablas de S1 que lee o escribe cada operación), `x-necesidades` (necesidades N-xxx de U1 que cubre), `x-tabla` (tabla de un esquema de entidad), `x-calculado` (campo calculado por la API o un motor; no es columna), `x-entrada` (dato de una acción que no se guarda, p. ej. la contraseña o la confirmación de una restauración), `x-columna` (columna a la que corresponde una entrada en lista, p. ej. `plan_dates` → `plan_date`), `x-app-setting` (clave de `app_setting`) y `x-provisional` (depende de un bloqueo abierto).
5. **Booleanos**: las columnas 0/1 de SQLite se exponen como `true`/`false`.
6. **Valores nulos**: en las respuestas, una columna NULL **se omite** (campo ausente = sin valor; la pantalla muestra «—», U1 §0.2 «sin datos no es cero»). Los esquemas de lectura marcan como `required` solo las columnas NOT NULL.
7. **Altas** (`POST`) aceptan las columnas escribibles; `id`, `created_at` y `updated_at` los pone el servidor. **Cambios** (`PATCH`) siguen JSON Merge Patch (RFC 7396): solo se cambian los campos enviados y `null` borra un valor que admite NULL. Reemplazos completos de listas (orden de categorías, destinatarios, grupos de producto) usan `PUT` con `{"ids": [...]}`.
8. **Idempotencia**: todo `POST` de captura exige el encabezado `Idempotency-Key` (16 a 128 caracteres `[A-Za-z0-9_-]`, UUID v4 recomendado). Repetir la clave con el mismo cuerpo devuelve la respuesta original con `Idempotent-Replayed: true`; con otro cuerpo, 409 `idempotency_key_reused`. La API guarda cada clave por cuenta en la tabla `idempotency_key` de S1, con la huella SHA-256 del cuerpo y la respuesta (S1 C-15), y la declara en `x-tablas` de cada operación que la exige; un job la purga por antigüedad. Retención de 24 h (K6). Gastos y pesos guardan además su clave en el propio renglón, sin purga, y responden 200 con el mismo recurso al reintento (K9). Lo exigen {{N_IDEM}} operaciones. No la llevan: las consultas por `POST` que no guardan nada (`previewCreditPurchase`, `compareScenarios`, `previewPlannedItemsImport`, `suggestRecipeWithAi`), las acciones de sesión (`login`, `logout`, `changeOwnPassword`, `revokeSessions`), el aviso de prueba (`sendTestNotification`) y la restauración, que lleva su propia confirmación escrita (`restoreBackup`).

   **Ciclo de una captura de gasto o de peso** (`createExpense` y `createWeightEntry`; K9 y K6). La cola del teléfono manda cada captura con su propia clave y su `captured_at`:
   - **Primer envío**: 201 con el recurso nuevo y su `id`. La clave queda en la tabla `idempotency_key` (con la huella del cuerpo y la respuesta) y en la columna `idempotency_key` del renglón (`expense` o `weight_entry`).
   - **Reintento con la misma clave y el mismo cuerpo**: 200 con el mismo recurso (mismo `id`) e `Idempotent-Replayed: true`; no se crea otro registro.
   - **Misma clave con otro cuerpo dentro de las 24 h**: 409 `idempotency_key_reused`; no se crea ni se cambia nada.
   - **Pasadas las 24 h** (el job ya purgó la clave de la tabla `idempotency_key`): la clave sigue en el renglón, así que la misma clave responde 200 con el registro guardado, mismo `id` e `Idempotent-Replayed: true`, sin crear otro. Ya no hay huella del cuerpo y solo se compara la clave: un cuerpo distinto ya no se detecta (no hay 409) y el registro guardado no cambia.

   Ejemplo con un gasto de contado de $199.00 (el monto de `parametros.papel` en `oraculo/model.json`), capturado sin red el 10-oct-2026 a las 21:30 y enviado al volver a la red de casa. `id` y `category_id` son identificadores de ejemplo; las respuestas van abreviadas:

   ```http
   POST /api/v1/expenses
   Idempotency-Key: 3f6c2a8e-7b1d-4c9a-9e2f-5a7d1b0c4e81

   {"category_id": 5, "amount_cents": 19900, "expense_date": "2026-10-10", "payment_method": "cash", "captured_at": "2026-10-10T21:30:00-06:00"}
   ```

   - Primer envío → **201** `{"id": 812, "category_id": 5, "amount_cents": 19900, "expense_date": "2026-10-10", "payment_method": "cash", "captured_at": "2026-10-10T21:30:00-06:00", …}`.
   - La cola no recibió la respuesta y reintenta con la misma clave y el mismo cuerpo → **200**, encabezado `Idempotent-Replayed: true` y el mismo cuerpo de respuesta, con `"id": 812`.
   - Misma clave con otro cuerpo (por ejemplo, otro `amount_cents`) dentro de las 24 h → **409** `{"code": "idempotency_key_reused", "message": "…", "request_id": "…"}`.
   - Misma clave después de las 24 h, con el mismo cuerpo o con otro → **200**, `Idempotent-Replayed: true` y el gasto 812 tal como se guardó.

   `createWeightEntry` (`POST /api/v1/weight-entries`) sigue el mismo ciclo con el peso.

9. **Identificador de petición**: toda respuesta lleva `X-Request-Id`; el mismo valor va en `Error.request_id` y en `audit_log.request_id` de los cambios que haga la petición.
10. **Cálculos**: la API no recalcula reglas en el cliente. Las cifras marcadas «motor Sx» las calcula ese motor con las reglas de `reglas-negocio.md`; la interfaz solo muestra, filtra, ordena y suma lo visible (U1 §0.2).
11. **Valores técnicos aprobados (K6)**: 200 registros por página como máximo (`limit`); `Idempotency-Key` retenida 24 h; cookie `SameSite=Strict`; contraseña nueva de 12 caracteres o más; carga de PROFECO de 1 GB como máximo; `cursor` de 512 caracteres como máximo; búsqueda (`q`) de 200 caracteres como máximo; usuario de 3 a 64 caracteres; clave de idempotencia de 16 a 128 caracteres. Desde C1 v1.2, la retención de `Idempotency-Key`, la contraseña mínima y el máximo de PROFECO son parámetros configurables que nacen con esos valores (K44, §11.7).

## 2. Autenticación, sesiones y permisos

- **Esquema**: `cookieAuth`, cookie `fh_session` **HttpOnly; Secure; SameSite=Strict; Path=/**. El token en claro solo viaja en la cookie; la base guarda su SHA-256 (`device_session.token_hash`). Ninguna respuesta incluye hashes, tokens, `chat_id` ni llaves (C1-12). Los chats que detecta el bot se muestran con el nombre de Telegram y el `chat_id` enmascarado (`chat_id_masked`: «•••» y los últimos 4 caracteres; K4).
- **Operaciones públicas**: solo `POST /api/v1/auth/login` y `GET /health` (sin datos). Todo lo demás (salvo las cuatro operaciones del agente externo, abajo), incluidos adjuntos, PDF y `GET /openapi.yaml`, exige sesión; sin sesión, vencida o cerrada responde **401 `unauthorized`**.
- **Credencial de servicio del agente externo** (K53, C1 v1.4): las cuatro operaciones con `x-agente-externo` (`/api/v1/external-agent/*`) piden el esquema `serviceToken`, un encabezado `Authorization: Bearer` con un token propio, y no aceptan la cookie de sesión. Ninguna otra operación acepta ese token. Sin token, o con uno inválido, responden **401 `unauthorized`**. El token en claro nunca se guarda ni se registra. Dónde vive la credencial y a nombre de quién quedan `audit_log` e `idempotency_key` es PROVISIONAL (bloqueo `S5-agente-externo`, pregunta e).
- **Inicio de sesión** (`login`, N-010): usuario y contraseña (argon2id en la base). Error de credenciales: 401 `invalid_credentials` con el mensaje genérico «Usuario o contraseña incorrectos». **Límite (decisión A): 5 intentos fallidos en 15 minutos** (`app_setting auth.max_failed_attempts = 5`, `auth.lockout_window_minutes = 15`). El intento siguiente responde **429 `too_many_attempts`**, mensaje «Demasiados intentos. Vuelve a intentar a las 21:45.» y `retry_at` (instante con -06:00), **aunque la contraseña sea correcta**; así se cubre también el «estado de bloqueo» (N-011) sin una consulta pública que revele usuarios. Aprobado en la puerta 1 (K11): ventana deslizante de 15 minutos; se cuentan los fallos por usuario y por IP (lo que llegue primero a 5); `retry_at` = hora del fallo más antiguo de la ventana + 15 minutos.
- **Sesión larga por dispositivo** (decisión A, «meses»): 90 días desde el último uso, renovable con cada petición (K11; `device_session.expires_at`; `app_setting session.max_age_days = 90`, sembrado en 00002). `logout` cierra la del dispositivo; Ajustes lista las sesiones de **ambos** usuarios y cualquiera puede cerrarlas (`listSessions`, `revokeSession`, `revokeSessions`, N-271 y N-272). Una sesión cerrada responde 401 en la siguiente petición.
- **Permisos**: Axel y Ximena son co-titulares con **los mismos permisos**; ninguna operación filtra ni oculta datos según quién inició sesión (U1 §0.5). El usuario de la sesión solo da valores por defecto (persona del gasto, del peso) y el «creado por» de `audit_log`. **No existe operación para crear cuentas**: Karen no tiene cuenta ni acceso (decisión A); `person` no se crea desde la API.
- **CSRF**: con `SameSite=Strict` (K6) el navegador no envía la cookie en peticiones que vienen de otro sitio.
- **Contraseña mínima** del cambio de contraseña (N-273): 12 caracteres al sembrar (K6); se usa el valor vigente de `auth.password_min_length` (K44, §11.7).

## 3. Modelo de errores (esquema único)

Todas las respuestas 4xx y 5xx de todas las operaciones usan `components/schemas/Error` mediante `components/responses/*` (más `default` para cualquier otro código):

```json
{"code": "validation_failed", "message": "Revisa los campos marcados.", "fields": [{"field": "amount_cents", "code": "out_of_range", "message": "El monto debe ser mayor que cero."}], "request_id": "01JABCDEF0123456789"}
```

- `code`: estable, en inglés, `^[a-z][a-z0-9_]*$`. `message`: español de México, lenguaje simple (qué pasó y qué hacer). `fields`: detalle por campo (nombre_codigo). `request_id`: siempre. `retry_at`: solo en 429.

| HTTP | Uso | Códigos estables |
|---|---|---|
| 400 | Petición mal formada: JSON inválido, parámetro con formato incorrecto (p. ej. `period=oct-2026`, fecha en formato de pantalla `10-oct-2026`) | `bad_request` |
| 401 | Sin sesión, vencida o cerrada; credenciales incorrectas en el inicio de sesión | `unauthorized`, `invalid_credentials` |
| 403 | Acción prohibida sobre el registro (precio de PROFECO) | `forbidden`, `read_only_source` |
| 404 | El registro o la ruta no existe (también una persona no externa en el estado de cuenta) | `not_found` |
| 405 | Método no admitido en una ruta existente | `method_not_allowed` |
| 409 | Conflicto con el estado: ya existe, en uso, gasto con compra a crédito, receta planeada, casilla con «fuera», vigencia de presupuesto traslapada, IA apagada, lectura asistida sin veredicto, clave de idempotencia reutilizada | `conflict`, `already_exists`, `in_use`, `has_credit_purchase`, `has_payments`, `recipe_planned`, `slot_conflict`, `budget_overlap`, `ai_disabled`, `assisted_reading_not_allowed`, `idempotency_key_reused` |
| 413 | Archivo demasiado grande | `payload_too_large` |
| 415 | Tipo de archivo no admitido | `unsupported_media_type` |
| 422 | Validación de negocio o de campos (con `fields`) | `validation_failed`, `invalid_current_password`, `not_revolving`, `not_a_card`, `external_person` |
| 429 | Límite de intentos de inicio de sesión (con `retry_at`) | `too_many_attempts` |
| 500 / 501 / 503 | Error interno / operación aún no implementada (esqueleto WF-0) / servicio o archivo no disponible | `internal_error`, `not_implemented`, `unavailable` |

**Validación de dinero**: un importe que no es entero (`973.75`, `"973.75"`) es 400 `bad_request`; un entero fuera de regla (cero o negativo donde se pide mayor que 0, partes que no suman el total) es 422 `validation_failed` con `fields[].code = out_of_range` o `invalid_sum`. **Validación de fechas**: una fecha que no es `AAAA-MM-DD` válida (`2026-02-30`, `10-oct-2026`) es 400; una fecha válida fuera de regla (pago antes de la compra, fin antes del inicio) es 422.

## 4. Fechas y dinero

- **Fechas de negocio** (`Date`, sufijo `_date`, y `valid_from`/`valid_to`): ISO 8601 `2026-10-10`, calendario de `America/Mexico_City`. **Instantes** (`DateTime`, sufijo `_at`): ISO 8601 con desplazamiento explícito de `America/Mexico_City`, hoy `-06:00` (`2026-10-10T21:30:00-06:00`); la base los guarda en UTC con «Z» y la API los convierte. **«Z» solo se acepta en las entradas** (R2-06, K12): la API recibe `2026-10-11T03:30:00Z` o `2026-10-10T21:30:00-06:00` y responde siempre con el desplazamiento de `America/Mexico_City`. **Mes** (`Month`): `2026-10`. Ningún formato de pantalla (`10-oct-2026`, `Vie 30 oct 2026`) aparece en la API; la interfaz los forma.
- **Quincenas**: los días de pago son el **15 y el 30**; en meses más cortos el «30» es el **último día** (en febrero, el 28: `2027-02-28`, R-23). Rangos para agrupar: Q1 del 1 al 15 y Q2 del 16 al último día; el 31 cae en Q2 (K8, firme). **Semana**: lunes a domingo (K8).
- **Dinero** (`MoneyCents`): entero exacto de **centavos MXN**, sufijo `_cents`, nunca `number` de punto flotante. `97375` = $973.75. Redondeo **ROUND_HALF_UP**; en una serie de pagos **el último pago ajusta centavos** (5,034.74 → 1,258.68 ×3 + 1,258.70, R-18 a R-21). Los dos modos de base de la serie que usa el oráculo y su generalización están en `S1-redondeo-series` (las series se guardan ya calculadas, el contrato solo las expone). En pantalla: `$`, separador de miles y dos decimales.
- **Porcentajes** (`BasisPoints`, sufijo `_bp`): 10000 = 100 %. Los porcentajes de avance calculados (`used_pct`, `progress_pct`) son enteros ROUND_HALF_UP. **Tasas**: fracción exacta `rate_num`/`rate_den` (S1; la tasa redondeada no reproduce el oráculo) y, para mostrar, `rate_pct` en texto **como se capturó** (K12): la fracción redondeada ROUND_HALF_UP a `rate_decimals` decimales (`17.27`, `25.8685`). **Cantidades** (`Milli`, sufijo `_milli`): milésimas de la unidad base (g, ml, pieza) o de porción.

## 5. Paginación y orden (una sola convención)

- **Colección** = lista de registros de una entidad. Toda colección (`{{N_LISTAS}}` operaciones `GET`) usa **los mismos parámetros**: `cursor` (opaco; sin valor = primera página), `limit` (por defecto **50**, como pide U1; máximo **200**, K6) y `sort` (propiedades separadas por coma, `-` = descendente; el `id` cierra siempre el orden para que el cursor sea estable). **Misma envoltura** `Page`: `{"items": [...], "next_cursor": "…" | null}`; **`next_cursor: null` indica que no hay más**. El orden por defecto de cada colección está en su descripción (tomado de U1 §2.5.b).
- **Agregados y métricas** (tarjetas, resúmenes, series de 12 meses o 24 quincenas, comparativas) son objetos calculados con tamaño acotado; no se paginan.

## 6. Filtros globales (definición única)

Los seis filtros del plan (§3.7) tienen **una sola definición** en `components/parameters` y todas las operaciones los usan por `$ref`:

| Filtro | Parámetro | Tipo y formato | Significado |
|---|---|---|---|
| persona | `person_id` | ids separados por coma | `person.id`; en gastos, `expense.person_id`; en partes de dinero, la persona de la parte |
| familia | `family` | `household` · `external` · `joint` | Tres valores incluyentes (K5): `household` = Hogar, reparto Hogar o Conjunto; `external` = Karen, reparto Karen o Conjunto; `joint` = solo Conjunto |
| categoría | `category_id` | ids separados por coma | `category.id` |
| tienda | `store_id` | ids separados por coma | `store.id` (en promociones, el comercio ligado) |
| periodo | `period` | `2026-10` mes · `2026-10-q1`/`q2` quincena · `2026-W42` semana · `2026` año · `2026-10-10..2026-11-09` rango | mes calendario como base; quincena del 1 al 15 y del 16 al último día; semana de lunes a domingo (K8). La lista de compras no usa este filtro: su semana va en la llave `week` (abajo) |
| crédito | `credit_id` | ids separados por coma | `credit.id`; en gastos, la forma de pago con esa línea |

Además se reutilizan `q` (texto) y los parámetros de paginación. Los filtros propios de una entidad (`status`, `kind`, `source`, etc.) se documentan en cada operación.

**Llaves que no son filtros.** Cuatro parámetros hablan de un periodo pero no filtran: dicen sobre qué recurso actúa la operación. Cada uno tiene una sola definición en `components/parameters`, su nombre no es el de ningún filtro y no acepta los demás formatos de `period`:

| Llave | Componente | Tipo y formato | Operaciones | Significado |
|---|---|---|---|---|
| `week` | `ListWeek` | semana ISO `2026-W42` (lunes a domingo, K8); obligatoria, sin valor por omisión | `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getDashboardBestBasket`, `getBasketComparison` | La lista de compras semanal que se lee o se escribe: la que muestra la pantalla (R3-01). La semana que se abre por defecto la decide la PWA (U1 §0.2) |
| `month` | `BudgetMonth` | `AAAA-MM`; sin valor, el mes en curso | `updateBudget` | El mes al que aplica un cambio de monto (§7.1, K2) |
| `from` | `OverviewFrom` | `AAAA-MM`; obligatoria | `getDashboardMonthlyOverview` | Primer mes de la ventana del resumen mensual del tablero (K54); `from` ≤ `to` y máximo 24 meses, si no 422 |
| `to` | `OverviewTo` | `AAAA-MM`; obligatoria; puede ser posterior al mes en curso | `getDashboardMonthlyOverview` | Último mes de la ventana, incluido; los meses posteriores al en curso son proyección |

**Matriz recurso × filtro** (operaciones que admiten cada filtro; la última columna marca la llave `week` de la lista de compras, que no es filtro):

{{FILTROS}}

## 7. Recursos: las 31 entidades de plan §3.4.1

Cada entidad tiene su recurso (lectura y, cuando el usuario la edita, escritura) o declara qué motor u operación la usa. Las vistas siguen U1 §2.5.a y se piden con `GET /api/v1/views/{entity}` (§8).

{{ENTIDADES}}

### 7.1 Presupuesto: alcance de un cambio de monto (`scope`, K2)

`updateBudget` (`PATCH /api/v1/budgets/{budget_id}`, N-087 y N-278) recibe dos parámetros de consulta cuando cambia `amount_cents`: `month` (AAAA-MM, el mes que muestra la pantalla; sin valor, el mes en curso), que debe caer en la vigencia del renglón, y `scope`, con dos valores. La pantalla preselecciona «de este mes en adelante».

| `scope` | En pantalla | Efecto sobre los renglones de vigencia (`valid_from_date`, `valid_to_date`) | Responde |
|---|---|---|---|
| `onward` (por omisión) | «de este mes en adelante» | El renglón termina el último día del mes anterior a `month`; un renglón nuevo con el monto nuevo rige desde el día 1 de `month` sin fin (`valid_to_date` NULL); los renglones posteriores de la misma categoría (o del total) se eliminan porque el monto nuevo los reemplaza (quedan en `audit_log`) | El renglón nuevo |
| `this_month` | «solo este mes» | El renglón termina el último día del mes anterior a `month`; un renglón nuevo con el monto nuevo rige solo en `month` (`is_single_month` = true) y otro renglón nuevo con el monto anterior rige desde el mes siguiente hasta el fin original; los renglones posteriores no cambian | El renglón de `month` |

Si el renglón empieza en `month`, el monto cambia en ese mismo renglón (con `onward` queda sin fin y también se eliminan los posteriores; con `this_month`, el resto de su vigencia pasa a un renglón nuevo con el monto anterior). Así, con `onward` todos los meses desde `month` tienen el monto nuevo, como dice K2 («de este mes en adelante») y como lo guarda S1 (`valid_to_date` NULL = rige hasta que se cambie), y nunca queda un traslape. `createBudget` solo llena meses sin monto y no elimina renglones. En los dos casos la respuesta es el renglón que rige en `month`, con `following_month_amount_cents` = el monto que rige el mes siguiente. Los cambios de `valid_to_date` o `notes` no usan `scope` ni `month`. Un traslape con otro renglón de la misma categoría responde 409 `budget_overlap`; `month` fuera de la vigencia, 422 `out_of_range`.

**Ejemplo.** Los montos se toman del oráculo y no son presupuestos capturados: 1,947.50 y 2,517.37 son el devengado de despensa de P1 y de P3 en `oraculo/model.json` (`periodos[0]` y `periodos[2]`, `devengado.despensa_total`). `id` y `category_id` son identificadores de ejemplo. Antes del cambio, la categoría 5 tiene un solo renglón:

| `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
|---|---|---|---|
| 7 | 194750 | 2026-10-01 | — (sin fin) |

En P-06, con diciembre a la vista, se cambia el monto a 2,517.37: `PATCH /api/v1/budgets/7?month=2026-12&scope=…` con el cuerpo `{"amount_cents": 251737}`.

- **`scope=onward`** (o sin `scope`). Después:

  | `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
  |---|---|---|---|
  | 7 | 194750 | 2026-10-01 | 2026-11-30 |
  | 8 (nuevo) | 251737 | 2026-12-01 | — (sin fin) |

  Respuesta **200** (abreviada; `valid_to_date` se omite porque es NULL, §1 punto 6): `{"id": 8, "category_id": 5, "amount_cents": 251737, "valid_from_date": "2026-12-01", "is_single_month": false, "following_month_amount_cents": 251737, …}`.

- **`scope=this_month`**. Después:

  | `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
  |---|---|---|---|
  | 7 | 194750 | 2026-10-01 | 2026-11-30 |
  | 8 (nuevo) | 251737 | 2026-12-01 | 2026-12-31 |
  | 9 (nuevo) | 194750 | 2027-01-01 | — (sin fin) |

  Respuesta **200** (abreviada): `{"id": 8, "category_id": 5, "amount_cents": 251737, "valid_from_date": "2026-12-01", "valid_to_date": "2026-12-31", "is_single_month": true, "following_month_amount_cents": 194750, …}`.

- **`scope=onward` con renglones posteriores** (sigue del resultado de `this_month`: el 7 de octubre a noviembre, el 8 solo en diciembre y el 9 desde enero). Con noviembre a la vista se cambia el monto a 2,517.37 «de este mes en adelante»: `PATCH /api/v1/budgets/7?month=2026-11&scope=onward` con el cuerpo `{"amount_cents": 251737}`. Después:

  | `id` | `amount_cents` | `valid_from_date` | `valid_to_date` |
  |---|---|---|---|
  | 7 | 194750 | 2026-10-01 | 2026-10-31 |
  | 10 (nuevo) | 251737 | 2026-11-01 | — (sin fin) |

  El 8 (solo diciembre) y el 9 (desde enero) se eliminan: si se quedaran, enero seguiría en 1,947.50 y el monto no regiría «en adelante», y el 10 sin fin se traslaparía con los dos (409). Respuesta **200** (abreviada): `{"id": 10, "category_id": 5, "amount_cents": 251737, "valid_from_date": "2026-11-01", "is_single_month": false, "following_month_amount_cents": 251737, …}`.

### 7.2 Reparto: alcance de un cambio de partes (`scope`, Q2-68 C)

Al cambiar las partes de una regla de reparto (`updateCostShareRule` con `parts`) la pantalla pregunta el alcance, como el presupuesto en K2 (RN-03 1.4; P-26.12). `CostShareRuleUpdate.scope` es obligatorio cuando viene `parts`, no tiene valor por omisión y no se acepta sin `parts` (422 `validation_failed` en los dos casos). Cambiar solo `name` o `notes` no lleva `scope`.

| `scope` | En pantalla | Qué hace |
|---|---|---|
| `onward` | «Solo los nuevos» | Cambia las partes de la regla. Los gastos ya guardados conservan sus partes, que están congeladas en `expense_share` (S1 1.4); los que se capturen después usan las nuevas. |
| `all` | «También los ya capturados» | Cambia las partes y el motor de reparto recalcula con ellas el `expense_share` de los gastos ya capturados, en la misma transacción, y `audit_log` lo registra. Qué gastos abarca (los que usan la regla en el gasto o en un renglón, o también los que la toman por su categoría o por ser la de por defecto) es PROVISIONAL: bloqueo `C1-reparto-alcance`. |

Ejemplo ilustrativo (cifras sintéticas): «Conjunto 50/50» pasa a 60/40. Un gasto ya guardado de $100.00 con esa regla tiene 5000 y 5000 en `expense_share`. Con `onward` se queda así y un gasto nuevo de $100.00 queda en 6000 y 4000; con `all`, si el gasto queda dentro del alcance, pasa a 6000 y 4000.

## 8. Vistas de base de datos (tabla, tablero, calendario, galería)

`GET /api/v1/views/{entity}?mode=table|board|calendar|gallery` (decisión G) aplica los filtros globales, `q`, `status`, `sort` y la paginación únicos, y responde `ViewPage` = `Page` + `mode` + `groups`. En **tablero**, `group_by` elige la columna y cada grupo trae `count`, `total_cents` (Σ del importe principal) y `share_bp`; en **calendario**, `date_field` elige la fecha y los grupos son días; en **galería**, cada registro trae `image_attachment_id` cuando hay foto. La vista elegida, filtros y columnas visibles se recuerdan en el dispositivo (U1 §2.5 y su supuesto 2); no hay vistas guardadas compartidas en esta versión.

{{VISTAS}}

## 9. Métricas del dashboard (fórmula, periodo, filtros, tablas, endpoint, ejemplo del oráculo)

Los bloques del tablero de Inicio › Resumen (U3 2.0, K54; antes 6 tarjetas, U1 N-020 a N-025) y las métricas de sus subsecciones. Periodo por defecto: **mes calendario**; la quincena es subperiodo. «Parte del hogar» = parte a nombre del titular; toda cifra de presupuesto, flujo y ahorro usa la parte del hogar (U1 §0.2). Los ejemplos usan cifras del oráculo (`oraculo/model.json`, identificadores R-nn de `docs/criterios/S3.md`) o un marcador cuando el oráculo no trae el dato. En la API van en centavos.

| Id | Tarjeta o subsección | Fórmula | Periodo | Filtros | Tablas de origen | Endpoint | Ejemplo (oráculo) |
|---|---|---|---|---|---|---|---|
| M-01 | Mes vs plan (N-020; K-1 «Gasto del mes» del tablero, K54) | real = Σ parte del hogar devengada en el mes (motor S2); plan = Σ `budget.amount_cents` del renglón vigente en el mes (RN-22) de las categorías activas; % usado = real ÷ plan × 100 entero ROUND_HALF_UP; restante = plan − real; estado normal < 80, aviso ≥ 80, excedido ≥ 100; días que faltan contando hoy. **Ampliada en C1 v1.5 (K-1 del tablero, K54; campos opcionales de `DashboardMonthVsPlan`)**: efectivo del hogar `cash_cents` = `cash_household_cents` de M-11 del mismo mes (contado + partes del hogar de pagos de crédito con fecha en el mes); mes anterior `previous_month` = el mes calendario anterior, y `previous_actual_cents` y `previous_cash_cents` = las mismas dos cifras (real y efectivo) de ese mes, ausentes si el sistema no tiene datos de ese mes (sin datos no es cero); diferencias `actual_difference_cents` = real − `previous_actual_cents` y `cash_difference_cents` = `cash_cents` − `previous_cash_cents`, con signo, ausentes sin mes anterior; conteos `categories_warning_count` y `categories_exceeded_count` = categorías activas del mes con `status` `warning` y `exceeded` en M-09 (motor S2) | mes en curso (`period`) y el mes calendario anterior | `period` | budget, category, period, expense, expense_share, credit_purchase; en C1 v1.5 también scheduled_payment, scheduled_payment_share (efectivo de M-11) | GET `/api/v1/dashboard/month-vs-plan` | oct-2026: real (parte del hogar devengada) **1,197.64** = `periodos[0].devengado.axel` (P1 = oct-2026, hallazgo 4 de S1); plan: marcador (lo captura Axel). C1 v1.5: oct-2026 efectivo del hogar **598.82** = `periodos[0].efectivo.axel` (parte del hogar del pago del 30-oct, R-13); nov-2026 (P2) real **1,197.64** = `periodos[1].devengado.axel` y efectivo **1,197.64** = `periodos[1].efectivo.axel`, con oct-2026 como mes anterior: `actual_difference_cents` **0.00** (1,197.64 − 1,197.64) y `cash_difference_cents` **598.82** (1,197.64 − 598.82) más que octubre; conteos: marcador (el plan por categoría lo captura Axel) |
| M-02 | Flujo quincenal (N-021) | ingresos = quincena capturada (o esperada, marcada) + extras de la quincena; egresos hechos = efectivo pagado hasta hoy (contado + partes del hogar pagadas); comprometido = partes del hogar pendientes de la quincena + fijos y venideros sin pagar; saldo al cierre = ingresos − egresos − comprometido | quincena en curso | `period` (`2026-10-q2`) | income, planned_item, expense, expense_share, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/dashboard/fortnight-flow` | Q2 oct-2026: comprometido incluye la parte del hogar del pago del 30-oct **598.82** (R-13); ingreso: marcador |
| M-03 | Próximo pago (N-022 a) | pago con alguna parte pendiente y fecha mínima ≥ hoy en todas las líneas; si hay pendientes con fecha < hoy, primero el más antiguo («sin marcar como pagado»); pagos del mismo día se suman | hoy en adelante | — | scheduled_payment, scheduled_payment_share, credit_payment, credit | GET `/api/v1/dashboard/next-payment` | **30-oct-2026, 973.75** (hogar 598.82, Karen 374.93) (R-13) |
| M-04 | Disponible A/B (N-022 b, N-201, N-285) | disponible de hoy de cada línea revolvente con la lectura activa (B por defecto) y aviso «pendiente confirmar»; vista previa con ambas lecturas y el mínimo del horizonte (motor S3) | hoy y horizonte | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/dashboard/next-payment`, GET `/api/v1/credits/{credit_id}/availability-preview` | 10-oct-2026: **A 10,000.00 · B 6,905.50** (R-64); mínimos **8,094.50 · 5,000.00** (R-61, R-62) |
| M-05 | Comidas de la semana (N-023) | costo = Σ costo de casillas con receta («fuera» no suma); kcal promedio por día = Σ kcal de la persona ÷ días con al menos un tiempo asignado; casillas sin armar = pares día × tiempo obligatorio sin receta ni «fuera» para alguna persona (máx. 21; K43) | semana en curso (lunes a domingo, K8) | `period` (`2026-W42`) | meal_plan_entry, meal_slot, recipe, recipe_ingredient, nutrition_fact, price, presentation, health_profile | GET `/api/v1/dashboard/week-meals` | marcador (el oráculo no trae comidas) |
| M-06 | Mejor canasta (N-024) | lista de compras vigente (o la última compra con renglones) comparada por tienda (motor S5); ganadora = menor total; diferencia contra la segunda; cobertura = artículos con precio ÷ artículos; precios viejos usados; Sam's como si tuvieran membresía | lista vigente | — | shopping_list, shopping_list_item, meal_plan_entry, recipe_ingredient, inventory_item, price, presentation, product, store, expense, expense_item | GET `/api/v1/dashboard/best-basket` | compra semilla 10-oct-2026: Bodega Aurrera **3,094.50** (R-01) contra Walmart 3,233.00 (`data.json totales.walmart`), diferencia **138.50** (plan §3.4 fila S5) |
| M-07 | Ahorro (N-025) | ahorrado = Σ aportaciones − Σ retiros (fondo y metas); aportado del mes contra planificado del mes (Σ `planned_monthly_cents` activos); meta más próxima con % = ahorrado ÷ monto meta × 100 entero | mes en curso | — | savings_goal, savings_movement | GET `/api/v1/dashboard/savings` | marcador (el oráculo no trae ahorro) |
| M-08 | Participación de Karen (N-026) | su parte devengada de la despensa del mes; sus compras vivas; sus partes de pagos del mes con estado; deuda viva = Σ partes pendientes | mes en curso | — | person, expense, expense_share, credit_purchase, credit_purchase_share, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/persons/{person_id}/participation` | oct-2026: parte devengada **749.86** (`periodos[0].devengado.karen`); su parte del pago del 30-oct **374.93** (R-11) |
| M-09 | Presupuesto por categoría (N-081, N-027) | por categoría: plan, real (parte del hogar devengada), plan − real, % usado y estado (umbrales de M-01) | mes (`period`) | `period` | budget, category, period, expense, expense_share | GET `/api/v1/budget/categories` | marcador (plan por categoría lo captura Axel) |
| M-10 | Flujo por quincena (N-090) | por quincena (24 hacia adelante): ingresos, partes del hogar de pagos, fijos y venideros, variable estimado (método de S2), ahorro planificado, saldo y saldo acumulado; efectivo | desde Q1 del mes en curso | — | income, planned_item, expense, expense_share, scheduled_payment, scheduled_payment_share, credit_payment, savings_goal | GET `/api/v1/cash-flow/fortnights` | Q2 oct-2026: pagos de crédito, parte del hogar **598.82** (R-13) |
| M-11 | Devengado contra efectivo (N-091, N-028) | por mes: devengado y efectivo (totales y partes); al cierre de cada mes, con el mismo origen para los tres (desde el primer evento de las entradas, S2-interfaz §3.2; Q2-25 A y Q2-65 A): devengado acumulado `cumulative_accrued_cents` = Σ devengado desde el primer evento hasta el cierre del mes, efectivo acumulado `cumulative_cash_cents` = Σ efectivo en el mismo tramo, y diferencia acumulada `cumulative_difference_cents` = devengado acumulado − efectivo acumulado (no se rotula deuda viva; `criterios-invariante-deuda`). El origen no es el primer mes visible: los meses anteriores a la ventana cuentan, así que la cima de T-2 es el devengado acumulado sea cual sea el mes en curso | mes en curso a +11 | — | expense, expense_share, planned_item, scheduled_payment, scheduled_payment_share, credit_purchase | GET `/api/v1/cash-flow/accrual-vs-cash` | oct-2026 (P1): devengado **1,947.50**, efectivo **973.75**, acumulados 1,947.50 y 973.75, diferencia acumulada 973.75 (R-45). Ventana que empieza en nov-2026 (reloj en nov-2026): su primer mes (P2) trae devengado 1,947.50 y efectivo 1,947.50 (R-46), y acumulados **3,895.00** y **2,921.25** (oct + nov; diferencia 973.75); sumar solo lo visible daría 1,947.50 y 1,947.50. mar-2027 (P6): acumulados **14,362.48** y **13,103.78** = totales P1–P6 (R-51, R-52), diferencia acumulada 1,258.70 (R-50). Cifras sumadas de `periodos[i]` de `oraculo/model.json` (Decimal) |
| M-12 | Deuda viva a una fecha (N-092) | Σ de lo que queda por pagar después de la fecha (capital + intereses) de las compras hechas hasta esa fecha; total, por línea y por compra (motor S3) | fecha elegida | `date` | credit_purchase, scheduled_payment, credit | GET `/api/v1/debt/live` | 10-abr-2027: **6,293.44** (R-53 = 1,258.70 R-54 + 5,034.74 R-55) |
| M-13 | Calendario unificado de pagos (N-192, N-027) | pagos de todas las líneas en el rango con partes y estado; promociones vigentes y cortes de tarjeta | mes visible o hoy..+60 | `period`, `credit_id`, `family` | scheduled_payment, scheduled_payment_share, credit_payment, credit_purchase, promotion, credit_statement | GET `/api/v1/scheduled-payments`, GET `/api/v1/calendar/events` | pagos de la despensa 1: 30-oct, 15-nov, 30-nov y 15-dic-2026 por **973.75** (R-13 a R-16) |
| M-14 | Plan contra real por mes (N-085, N-028) | por mes (máx. 12): plan total y real total (o de una categoría) | 12 meses | `category_id` | budget, category, expense, expense_share | GET `/api/v1/budget/monthly-trend` | marcador (plan lo captura Axel) |
| M-15 | Costo semanal y kcal (N-129, N-233, N-028) | por semana (12): costo de casillas con receta y kcal promedio por día por persona contra su meta | 12 semanas | `person_id` | meal_plan_entry, recipe, recipe_ingredient, nutrition_fact, price, health_profile | GET `/api/v1/meal-plan/trends` | marcador |
| M-16 | Metas (N-111, N-027) | por meta: ahorrado, % de avance entero, aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan (motor S2) y estado | hoy | `status` | savings_goal, savings_movement | GET `/api/v1/savings-goals` | marcador |
| M-17 | Deuda viva proyectada (N-194, N-209) | por mes (12): deuda viva al cierre (M-12 al último día del mes, PROVISIONAL, bloqueo `S6-deuda-viva-hoy`, punto 2), total de todas las líneas, activas o archivadas (PROVISIONAL, bloqueo `S6-tablero-lineas-archivadas`, punto 5), y por línea; intereses pagados y por pagar | 12 meses | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/credits/debt-trend` | punto de control: deuda viva al 10-abr-2027 **6,293.44** (R-53); los meses sin cifra canon llevan la marca de `criterios-proyeccion-12m` |
| M-18 | Ingreso del mes (K-2 del tablero, N-100) | quincenas del mes (capturadas o esperadas, marcadas) + extras; quincenas y extras del mes anterior; diferencia de quincenas = mes − mes anterior, con signo (motor S2) | mes (`period`) | `period` | income, planned_item, period | GET `/api/v1/incomes/month` (campos `previous_*` y `fortnights_difference_cents`, C1 v1.5) | marcador (el oráculo no trae ingresos) |
| M-19 | Créditos del mes (K-3 del tablero, N-191) | pagos con fecha en el mes en curso: total, parte del hogar, parte de Karen, cuántos (`month_payment_count`) y parte del hogar que falta pagar (`month_pending_household_cents`: la parte del hogar menos los abonos del hogar, nunca negativa; PROVISIONAL, bloqueo `S6-deuda-viva-hoy`, punto 3), de todas las líneas, activas o archivadas (PROVISIONAL, bloqueo `S6-tablero-lineas-archivadas`, punto 5) (motor S3) | mes en curso | — | credit, credit_purchase, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/credits/totals` | oct-2026: pago del 30-oct **973.75**, parte del hogar **598.82** (R-13) |
| M-20 | Capacidad de endeudamiento (B-1 del tablero, RN-30; K54) | pagos del mes (parte del hogar) ÷ ingreso mensual neto × 100, entero ROUND_HALF_UP; estado normal ≤ `debt.capacity_warning_pct` (30), aviso ≤ `debt.capacity_limit_pct` (40), excedido arriba; margen = tope × ingreso ÷ 100 − pagos; ingreso sin extras y pagos del hogar PROVISIONALES (bloqueos `U3-tablero-barometro` y `C1-tablero`) | mes (`period`, por omisión el en curso) | `period` | income, planned_item, period, scheduled_payment, scheduled_payment_share, credit, credit_purchase, app_setting | GET `/api/v1/dashboard/debt-capacity` | pagos del hogar oct, nov y dic-2026 **598.82**, **1,197.64** y **1,482.57** (`periodos[0..2].efectivo.axel`); ingreso: marcador |
| M-21 | Resumen mensual del tablero (G-1 y G-6) | por mes de `from` a `to` (máximo 24): ingreso (quincenas + extras), devengado del hogar (M-01), efectivo del hogar (M-11), plan, pagos de crédito (hogar y total) y la capacidad de M-20 con su estado y la misma escala del dibujo (`scale_max_pct`, C1 v1.5.1); los meses posteriores al en curso son proyección, sin devengado ni efectivo (motores S2 y S3) | `from`..`to` | — | income, planned_item, period, budget, category, expense, expense_share, credit_purchase, scheduled_payment, scheduled_payment_share, credit, app_setting | GET `/api/v1/dashboard/monthly-overview` | devengado del hogar oct, nov y dic-2026 **1,197.64**, **1,197.64** y **1,767.51** (`periodos[0..2].devengado.axel`) |
| M-22 | Sugerencia de presupuesto (G-3 del tablero, RN-31) | por categoría activa: % usado de los `budget.signal_months` (3) meses cerrados anteriores al mes; `raise` si todos pasan de 100 + `budget.signal_margin_pct` (10), `cut` si todos quedan abajo de 100 − margen, si no `none`; PROVISIONAL (bloqueos `U3-tablero-senal-presupuesto` y `C1-tablero`) | 3 meses cerrados antes de `period` | `period` | budget, category, period, expense, expense_share, credit_purchase, app_setting | GET `/api/v1/budget/signals` | marcador (el plan por categoría lo captura Axel) |

## 10. Matriz pantalla P-xx → necesidad N-xxx → endpoint(s)

Una fila por necesidad y pantalla (las necesidades reutilizadas aparecen en cada pantalla que las usa, como en U1 §3.29). Estado: **OK** = cubierta; **FIRME (Kxx)** = cubierta y cerrada en la puerta 1b por la decisión citada (antes PROVISIONAL); **PROVISIONAL** = cubierta en parte, con bloqueo; **BLOQUEO** = sin endpoint por un bloqueo abierto (hoy ninguna); **Local** = vive en el dispositivo (U1 la declara sin API).

{{MATRIZ}}

## 11. Cobertura

### 11.1 Resumen
{{RESUMEN}}

### 11.2 Operaciones sin necesidad de U1 (excepciones técnicas, una por una)
{{SIN_N}}

### 11.3 Necesidades sin endpoint
{{N_SIN_EP}}

### 11.4 Tablas de S1 → operaciones ({{N_TABLAS_OP}} de {{N_TABLAS}})
{{TABLAS}}

`idempotency_key` no tiene pantalla propia: es infraestructura de S6. La leen y escriben en la misma transacción todas las operaciones que exigen `Idempotency-Key` (§1 punto 8; S1 C-15), por eso figura en su `x-tablas`, y la purga por antigüedad la hace un job de S6 (retención de 24 h, K6). `inventory_movement`, `recipe_meal_slot`, `shopping_list` y `shopping_list_item` son detalle de entidades del plan (inventario, receta y lista de compras) y tienen sus propias operaciones.

### 11.5 Las 14 secciones de Ajustes → necesidades → operaciones
{{AJUSTES}}

### 11.6 Capacidades exigidas (C1-11) → operaciones
{{CAPACIDADES}}

### 11.7 Parámetros operativos configurables (K44) → operaciones de settings

Todo parámetro operativo es una clave de `app_setting` del hogar que nace con el valor decidido (K44; RN-29 de `reglas-negocio.md` 1.4). Esta tabla sale de las anotaciones `@param` de las migraciones y de los esquemas de settings de `openapi.yaml`; no se escribe a mano. `parse_schema.py` compara cada valor inicial con la siembra real y se detiene si no coinciden.

- **Lectura**: el `GET /api/v1/settings/…` de la sección trae `parameters`, un objeto por clave de `app_setting` (`SettingParameterInteger`, `SettingParameterIntegerList` o `SettingParameterTime`) con `value` (vigente), `initial_value` (sembrado; `null` si no lo hay), `is_initial`, `value_type` y `unit`.
- **Cambio**: el `PATCH` de la misma ruta recibe `parameters` con el valor nuevo por clave. Un valor vacío, `null` o de otro tipo responde 422 `validation_failed` y no se sustituye por otro. Se valida el tipo y el rango de la clave, el de la columna «Rango (cambio)»: es el de los disparadores de S1 (`trg_app_setting_param_int_*` de 00003; `trg_app_setting_param_range_*` y `trg_app_setting_param_time_*` de 00004) y lo declara la propiedad de la clave en el `PATCH` (`minimum`, `maximum` o `pattern`; puerta 2, K50). Fuera de rango: 422 `validation_failed`.
- **Restablecer** es un campo, no una operación: `reset_parameters` lista las claves que vuelven a su `initial_value`. Una clave no puede venir a la vez en `parameters` y en `reset_parameters` (422 `validation_failed`). Un mismo `PATCH` cambia unos y restablece otros, y `audit_log` registra el cambio.
- **Quién usa el valor vigente**: `login`, `changeOwnPassword` y la renovación de la sesión (`auth.*`, `session.max_age_days`); la frescura y el peso de los precios (`price.stale_after_days`, `price.zero_weight_after_days`, RN-12); `importProfeco`; `compareScenarios` (503 `unavailable`); `uploadExpenseAttachment` (413 con el máximo en el mensaje); la purga de `idempotency_key`; `suggested_kcal_target` (`health.*` y el ajuste por persona de `health_profile`, RN-17); los avisos de presupuesto y de pago próximo y la hora, la anticipación y el día de los demás avisos (`alerts.*`); los reportes programados (`reports.*`, S7); la ventana de fechas del comparador (`optimizer.date_window_days`) y el OCR (`ocr.timeout_seconds`).

{{PARAMS}}

Parámetros que existen sin valor sembrado o que viven en otra tabla (no se inventa un valor):

{{PARAMS_SV}}

{{PARAMS_RANGO}}

## 12. Índice de operaciones

{{OPERACIONES}}

## 13. PROVISIONAL y bloqueos

| Tema | Qué queda PROVISIONAL en C1 | Bloqueo |
|---|---|---|
| Alcance «también los ya capturados» de un cambio de reparto | qué gastos recalcula `scope=all` de `updateCostShareRule` (`CostShareRuleUpdate.scope`, §7.2) es PROVISIONAL | `C1-reparto-alcance` (abierto, de C1) |
| Hora de los reportes y `job.schedule` | cómo se relaciona `schedule` de los jobs de los reportes (`JobUpdate`, `updateJob`) con las claves `reports.*` de `updateReportSettings` es PROVISIONAL | `C1-hora-reportes-job` (abierto, de C1) |
| «No incluye IVA» sin porcentaje | `Credit.missing_terms` lista `vat` con `vat_included` = 0 y sin `vat_on_interest_bp` (lectura en uso de S3 1.5 §5.3); PROVISIONAL | `S3-iva-lectura` (abierto, pregunta 2) |
| Edición de un gasto anterior después de «solo los nuevos» | `updateExpense` recalcula sus partes con la regla vigente (opción en uso, RN-03 1.4); PROVISIONAL | `S1-reparto-edicion-gasto` (abierto) |
| Agente externo (K53) | origen `site_assisted` y peso de sus precios (`importExternalPrices`), borrador y origen de sus recetas (`importExternalRecipes`), qué escribe el agente (`ExternalCatalogImport`), el tamaño máximo de sus lotes (los tres esquemas de importación) y dónde vive la credencial `serviceToken`; PROVISIONAL | `S5-agente-externo` (abierto, de P2d) |
| Tablero de control (K54) | sección «Presupuesto» y rangos de `debt.capacity_warning_pct`, `debt.capacity_limit_pct`, `dashboard.payment_reminder_count`, `budget.signal_months` y `budget.signal_margin_pct` (`DashboardSettings`); necesidades de U1 de `getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSettings` y `updateDashboardSettings`; nombre de `dashboard.payment_reminder_count`; `period` solo `AAAA-MM`; PROVISIONAL | `C1-tablero` (abierto, de C1) |
| Qué mide el barómetro y la sugerencia | ingreso neto sin extras, pagos del hogar, redondeo y escala (`DashboardDebtCapacity`, `DashboardMonthlyOverview`); valores 5, 3 y 10 de las claves del tablero; regla de `getBudgetSignals`; PROVISIONAL | `U3-tablero-barometro`, `U3-tablero-recordatorio` y `U3-tablero-senal-presupuesto` (abiertos, de U3) |
| Ahorro planificado por quincena | `planned_savings_cents` de `getCashFlowFortnights` es opcional, S2 1.3 §3.11 la entrega en null y no entra a `balance_cents`; PROVISIONAL | `S2-ahorro-quincenal` (abierto) |

Cerrados en la puerta 1 y ya sin marca: filtro familia (K5), lista de compras (K1), huecos U1 ↔ S1 (K2, K3, K4), parámetros técnicos (K6), campos calculados (K7), quincena y semana (K8), inicio de sesión y sesión (K11), tasas (K12), tope por persona (K13), reposición de la parte externa (K15), restricciones del comparador (K16), puntaje (K17), promedios diario y mensual (K20), meses 7 a 12 (K21), frescura de precios (K22) y escala de actividad (K23). Cerrados en la puerta 1b (C1 v1.2) y ya sin marca: valores de K31 que tocan a C1 (K43: adjunto de 10 MB, comparador de 5 s, claves `send_time` y `advance_days`, «casillas sin armar»), ajuste por objetivo de la meta de kcal (K36), semana por omisión de `getMealPlanWeek` (K38; la lista de compras sigue con `week` obligatoria, R3-01), promedio semanal (K39) y parámetros configurables (K44, §11.7).

Cerrados en la puerta 2 (C1 v1.3) y ya sin marca: lugar de la hora, el día y la anticipación de reportes y avisos (`S1-parametros-avisos`, Q2-11 B, K47), sus valores (Q2-72, Q2-73), la ventana de fechas y el OCR (Q2-75, Q2-76), los rangos (Q2-74, Q2-77, Q2-66 B), el IVA de un crédito (`U4-iva-credito`, Q2-67), el alcance de un cambio de reparto como pregunta (`U4-reparto-retroactivo`, Q2-68 C), la tienda asignada de la lista (`S4-asignacion-mercado`, Q2-39 a Q2-42) y los acumulados de T-2 (`U3-t2-acumulados`, Q2-65).

Cerrado en la puerta 2 bis (C1 v1.4): el veredicto de términos de las tiendas (`S5-terminos`, K53). La lectura asistida no aplica y C1 no tenía marca que quitar.

## 14. Inconsistencias S1 ↔ U1 encontradas (sin tocar sus archivos)

En C1 v1.0 (WF-0) C1 tomó en cada caso la lectura más conservadora y la marcó PROVISIONAL, y pidió las correcciones a S1 en `docs/estado/wf0/S6a-contrato-c1.md`. El esquema 1.1 de S1 dio dónde guardar siete de los diez puntos y la puerta 1 cerró los demás (K1 a K5 y K6); el esquema 1.2 (00002) agregó el borrador del ticket y los chats por confirmar. C1 v1.1 se regeneró contra ese esquema. La columna «Estado» dice qué queda abierto.

{{INCONSISTENCIAS}}
