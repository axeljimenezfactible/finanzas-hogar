# P2c · C1-puerta2 · C1 v1.3 (puerta 2)

Agente C1-puerta2 (Opus), flujo P2c, 11-oct-2026 (noche). Fuentes: `docs/estado/puerta2/RESPUESTAS.md`, `PUERTA-2.md` §2, §3 y §5, `docs/estado/p2c/S1-puerta2.md` (esquema 1.4, `00004_puerta2.sql`), `docs/estado/p2c/contratos-puerta2.md` («Solicitudes a C1»), S4 1.5 §6.5 y §12.2, S5 1.5 §7.4, S2 1.3 (metas), S3 §8, U3 1.3, los bloqueos `U3-t2-acumulados`, `U4-reparto-retroactivo` y `U4-iva-credito`, y RN-03 y RN-26 de `reglas-negocio.md` 1.4.

## Hecho
Todo sale del generador (`backend/internal/api/c1gen/`). No edité a mano `openapi.yaml`, `C1-datos.md`, `routes.go` ni los extractos.

**Conteos** de `openapi.yaml`:

| | Rutas | Operaciones | Esquemas |
|---|---|---|---|
| Antes (C1 v1.2) | 167 | 210 | 252 |
| Después (C1 v1.3) | 168 | 212 | 255 |

- La ruta nueva es `/api/v1/settings/reports`, con 2 operaciones.
- Los 3 esquemas nuevos son `ReportSettings`, `ReportSettingsUpdate` y `SettingParameterTime`.
- `info.version` y `C1_VERSION` quedan en 1.3.0.
- Extractos: 17 YAML más `INDICE.md`, con las 212 operaciones, y ningún archivo pasa de 80 KB.

1. **Q2-11 B (K47), Q2-72 y Q2-73 (K50).** La hora, el día y la anticipación de reportes y avisos se leen y se cambian por clave, igual que los demás parámetros de K44: `value`, `initial_value`, `is_initial`, `parameters` y `reset_parameters`.
   - Operaciones nuevas: `getReportSettings` (N-296) y `updateReportSettings` (N-297), para las 5 claves `reports.*`.
   - `AlertSettings` suma `alerts.send_time`, `alerts.stale_prices_advance_days` y `alerts.meal_plan_missing_weekday`. Sus operaciones cubren además N-298 y N-300.
   - `params_json` de `Notification` y de `NotificationUpdate` ya no declara `send_time` ni `advance_days` y pierde `x-provisional`. Su descripción dice que con los valores iniciales solo `budget_threshold` conserva un parámetro (`scopes`).
   - `updateNotification` se queda para lo que solo vive en el aviso: si está activo y sus parámetros propios.
   - En §10, N-297 y N-300 pasan a «FIRME (K43, K47)» y ya no llevan marca de lugar. El bloqueo `S1-parametros-avisos` ya no se cita como abierto.
2. **Claves nuevas de S1** (`@param`), expuestas por el generador:
   - Tipo nuevo «hora HH:MM» (`SettingParameterTime`, patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`).
   - El rango de cada clave **se lee de la base**. `parse_schema.py` prueba cada valor con un `UPDATE` dentro de un `SAVEPOINT` que se deshace. Si el patrón no coincide con los disparadores de S1, el generador se detiene.
   - Cada clave del `PATCH` declara su rango con `minimum`, `maximum` o `pattern`. Los rangos que quedan:
     - día de la semana: de 1 a 7;
     - día del mes: de 1 a 31;
     - `comparator.timeout_seconds`: de 1 a 30;
     - `price.stale_after_days`: 1 o más;
     - anticipación, ventana y OCR: 0 o más.
   - `price.zero_weight_after_days` ≥ `price.stale_after_days` va en la descripción de `PriceSettingsUpdate`: 422 si no se cumple.
   - Ubicación de las claves que no son de reportes ni de avisos:
     - `optimizer.date_window_days` (7) va en `parameters` de `getOptimizerSettings` y `updateOptimizerSettings`, y deja de ser campo propio de esos dos esquemas. Se queda en `ScenarioRequest` y en `constraints_used`.
     - `ocr.timeout_seconds` (30) va en `DataSettings`.
   - `@param_rango` y `@param_sin_valor_cerrado` se leen. Los 5 parámetros cerrados ya no salen como «sin valor»: §11.7 los lista aparte.
   - La última `@col_note` de una columna sustituye a la anterior, como hace `gen_dict.py` de S1.
3. **Q2-67, opción 1 (K46).**
   - `Credit`, `CreditCreate` y `CreditUpdate` traen `vat_included` (boolean, sale de la columna de la 00004).
   - `Credit.missing_terms` lista 5 datos: `annual_rate`, `rate_base`, `vat`, `min_payment_pct_bp` y `min_payment_fixed_cents`.
   - El caso `vat_included` = 0 sin porcentaje queda PROVISIONAL y cita `S3-iva-lectura`.
4. **Q2-68 C.** `CostShareRuleUpdate.scope` puede ser `onward` («solo los nuevos») o `all` («también los ya capturados»).
   - Es obligatorio con `parts`, no tiene valor por omisión y responde 422 si llega sin `parts`.
   - `updateCostShareRule` suma `expense`, `expense_item` y `expense_share` a `x-tablas`.
   - C1-datos tiene una sección nueva, §7.2, con la tabla y un ejemplo sintético.
   - Qué gastos abarca `all` queda PROVISIONAL (bloqueo nuevo `C1-reparto-alcance`).
5. **Q2-39 a Q2-42 A (S4 1.5).**
   - `ShoppingListItem.assigned_store_id` es la tienda asignada y rige la agrupación.
   - `recommended_store_id` es la recomendación de S5. En la vertical `market` no rige la agrupación.
   - `estimated_price_cents` es el costo con la cotización que rige.
   - En `groups[]`:
     - `store_id` es la tienda asignada; si falta, el grupo es «sin mercado dado de alta»;
     - `visit_weekday` tiene texto nuevo y `visit_date` es campo nuevo;
     - `total_cents` suma solo los artículos con precio.
   - Orden de los grupos: mercado, cadena y, al final, el grupo sin tienda.
6. **Q2-65 A.** `AccrualVsCash.months[]` suma `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents`.
7. **Q2-66 B, Q2-72 y Q2-73.** En C1 solo cambian las claves y sus textos. El mes corto está en `updateReportSettings` y en la unidad de `reports.monthly_day_of_month`.
8. **Otros puntos de PUERTA-2 que nombran a C1:**
   - Q2-17 C: `SavingsGoal.progress_tenths_pct`.
   - Q2-33 (b): `Scenario.store_id`.
   - Q2-86 A: `DashboardBestBasket.second_store_id`.
   - S5 1.5 §7.4 (`contratos-puerta2.md`):
     - campo nuevo `BasketComparison.cheapest_store_id`, ausente si ninguna tienda está completa;
     - `stores[].difference_vs_cheapest_cents` solo aparece en las tiendas completas;
     - `stores[]` lleva primero las completas;
     - `identical_count` cuenta renglones.
9. **Texto de C1-datos.md.**
   - Cabecera, fila de «Cambios» y detalle «C1 v1.3 (puerta 2)».
   - §7.2.
   - §11.7 con la columna «Rango (cambio)», la regla entre claves y los parámetros cerrados.
   - §13 con 4 bloqueos abiertos (`C1-reparto-alcance`, `C1-hora-reportes-job`, `S3-iva-lectura` y `S1-reparto-edicion-gasto`) y el párrafo «Cerrados en la puerta 2».
10. **Bloqueos nuevos:** `docs/estado/bloqueos/C1-reparto-alcance.md` y `docs/estado/bloqueos/C1-hora-reportes-job.md`. Cada uno trae qué falta, las opciones, el impacto y cómo seguí.

**Archivos del generador que toqué** (`backend/internal/api/c1gen/`):
- `parse_schema.py`, `gen2_schemas.py`, `gen3_ops.py`, `gen4_build.py` (`C1_VERSION`), `gen6_c1.py` y `c1_plantilla.md`;
- `verify_c1.py`: 1.3.0, «C1 v1.3» y C1-23 de la puerta 2;
- `README.md`: sección «C1 v1.3 (puerta 2)».

## Verificación
- `python3 -I backend/internal/api/c1gen/build.py` → rc 0. Con `--check`: «los 21 archivos coinciden byte a byte con el generador».
- `make verify-c1` → «comprobaciones OK: 29 · RESULTADO: TODO OK». Las líneas clave:
  - C1-15: 6 marcas en openapi; bloqueos abiertos citados: `C1-hora-reportes-job`, `C1-reparto-alcance` y `S3-iva-lectura`; 0 marcas sin bloqueo.
  - EXT: 212 operaciones en los extractos y 212 en openapi.
  - S6-03: 212 rutas.
  - C1-23: claves de la puerta 2 con un rango distinto del aprobado: [] · filas de §11.7 sin su rango: [] · extracto con `NotificationUpdate` (`ajustes-d.yaml`) sin fallas.
  - K44 OK.
- `make verify-openapi` → `docs/contratos/openapi.yaml: OK`.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → rc 0 y «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». La salida completa quedó en `/tmp/c1p2/verify-all.log`.
- La primera corrida falló en C1-23 porque quedó una línea vieja (`k43c1`) en `verify_c1.py`. La quité y la segunda corrida pasó.

## Supuestos
1. Edité `verify_c1.py` y `README.md`. Están en `c1gen`, pero no los nombra la lista literal del ALCANCE. Sin el cambio de versión y de C1-23, `make verify-c1` no podía pasar.
2. **«Reportes y frecuencias»** es la sección que S1 dio a las claves `reports.*`. No tenía operación de settings, y la regla de K44 de C1 v1.2 pide una lectura y un cambio por clave. La ruta y los nombres (`/api/v1/settings/reports`, `getReportSettings`, `updateReportSettings`) los puse yo, siguiendo las otras secciones.
3. **Ventana de fechas:** `optimizer.date_window_days` sale de los campos propios de `OptimizerSettings` y `OptimizerSettingsUpdate` para no tener dos maneras de cambiar la misma clave. Sigue en:
   - `ScenarioRequest`, para sobreescribirla en una sola comparación;
   - `constraints_used`;
   - el enum de `missing_constraints`, porque S3 usa «el orden de C1».
4. **Alcance del reparto:** `scope` va en el **cuerpo** de `CostShareRuleUpdate`, no en la consulta. C1-06 no deja dos componentes de parámetro con el mismo nombre, y ya existe `BudgetScope` (`scope`). Tampoco tiene valor por omisión, porque Q2-68 C dice «se pregunta en cada cambio».
5. **Ausente en lugar de null:** «admite null» de S5 y S4 se traduce a «campo ausente», que es la convención 6 de C1 §1.
6. Agregué `groups[].visit_date`, que S4 §12.2 pide como solicitud opcional.
7. Apliqué Q2-17 C, Q2-33 (b) y Q2-86 A porque tienen recomendación y nombran a C1.

## Bloqueos
- Nuevo `C1-reparto-alcance`: qué gastos recalcula `scope = all`. Opciones A, B y C; propuesta B.
- Nuevo `C1-hora-reportes-job`: `job.schedule` frente a las claves `reports.*`. Opciones A, B y C; propuesta A. Mientras tanto, `updateJob` sigue aceptando `schedule`.
- Siguen abiertos y citados: `S3-iva-lectura` (pregunta 2) y `S1-reparto-edicion-gasto`.

## Pendientes
- PREGUNTA para la puerta 4, sin decisión en la puerta 2: S3 §8 pide además en `Scenario` y `Result` los campos `std_dev_milli_cents`, `over_cap`, `provisional` y `excluded`, y S5 §7.2 pide `stores[].piece_count` e `identical_pieces`. No los agregué.
- K53 (puerta 2 bis; se registra en P2d) no está aplicado. Siguen en C1:
  - la lectura asistida: `PriceSettings.assisted_reading_enabled`, `PriceSources.sites` y `updateStore` con 409 `assisted_reading_not_allowed`;
  - falta la importación del agente externo.
- Q2-94 («falta por comprar») no tiene recomendación: no se agrega.

## Propuestas para Axel
- `C1-reparto-alcance`: B. «También los ya capturados» abarca todos los gastos que usan la regla, también los que la toman por su categoría o por ser la de por defecto.
- `C1-hora-reportes-job`: A. Mandan las claves `reports.*`, la API reescribe el cron del job y `updateJob` no acepta `schedule` en los jobs de reportes.

## Solicitudes fuera de alcance
1. **Integrador:** agregar `C1-reparto-alcance` y `C1-hora-reportes-job` a `docs/estado/bloqueos/INDICE.md`.
2. **U1** (N-297, N-300 y la tabla K44 de §2.3): todavía dice que `send_time` y `advance_days` viven en `params_json`. Hay que cambiarlo a `reports.*` y `alerts.*` en `app_setting`, «Reportes y frecuencias» → `getReportSettings`.
3. **U4** (P-26.05, P-26.08, P-26.09 y P-26.12): usar `getReportSettings`/`updateReportSettings` y `parameters['optimizer.date_window_days']`, `vat_included` y `scope`, que es obligatorio con `parts`.
4. **Criterios** (`docs/criterios/C1.md`, C1-22 y C1-23): reflejar el lugar FIRME (K47) y los rangos.
5. **S3, S4, S5 y S7 (WF-3):** entregar los campos nuevos y leer las claves nuevas:
   - `store_id` del escenario;
   - `assigned_store_id` y `visit_date`;
   - `cheapest_store`;
   - los acumulados de T-2;
   - `progress_tenths_pct` y `second_store_id`;
   - en S7, las claves de reportes y avisos.

## Consumo aproximado
- Acciones de herramienta: unas 44.
- Leído: unos 210 KB, por partes (`sed`, `grep` y `awk`). No leí `openapi.yaml` ni `DESIGN.md` enteros.
- Escrito a mano: unos 32 KB entre parches del generador, la plantilla, el README, 2 bloqueos y este reporte.
- Regenerado por `build.py` (no a mano): `openapi.yaml` (700 KB), `C1-datos.md` (155 KB), `routes.go` (26 KB) y los extractos.
