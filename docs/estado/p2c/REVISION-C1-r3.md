# Revisión de C1 v1.3 (puerta 2) · ronda 3 de 3

Revisor: `revisor-C1-r3`, independiente. No escribí nada de lo que reviso y no corregí nada. Flujo P2c, 11-oct-2026, noche.

Evalúo contra:
- `docs/criterios/C1.md` (versión 1.2, la vigente), sobre todo C1-04, C1-13, C1-15, C1-21, C1-22 y C1-23;
- `PUERTA-2.md` §2 y §3 (Q2-11 B, Q2-39 a Q2-42 A, Q2-65 A, Q2-66 B, Q2-67 opción 1, Q2-68 C, Q2-72 a Q2-77);
- las claves de `S1-puerta2` (`00004_puerta2.sql`);
- las «Solicitudes a C1» de `contratos-puerta2.md`.

Todas las cifras y salidas salen de mis propios comandos y scripts en `datos/tmp/rev-c1-r3/`, que ya borré. Ningún resultado lo tomé de un reporte. Cada `make` lo corrí con un `TMP_ROOT` propio (B-4 de la ronda 2).

## Veredicto
**APROBADO.** Hay 0 altos, 0 medios y 6 bajos.

- El único medio de la ronda 2, **M-1**, está **cerrado**. Lo comprobé contra la base (V-8) y contra los 8 textos del contrato (V-9).
- Los bajos B-1 a B-5 de la ronda 2 siguen abiertos sin cambio (no venían en el encargo de `correccion-C1-r2`). Agrego 2 bajos nuevos (N-1 y N-2). Ninguno impide aprobar.

## Verificaciones (con mis comandos)

| Id | Qué | Comando o método | Resultado |
|---|---|---|---|
| V-1 | (a) `make verify-c1` | `make verify-c1 GOOSE=…/gobin/goose TMP_ROOT=…/datos/tmp/rev-c1-r3/mk` | rc 0. «RESULTADO: los 21 archivos coinciden byte a byte con el generador», «comprobaciones OK: 33» y «RESULTADO: TODO OK» |
| V-2 | (a) `verify_c1.py` completo (sin `--db`, migraciones en memoria) | `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py` | rc 0. 33 líneas `OK`, 0 `FALLA`, «RESULTADO: TODO OK». C1-04 IVA: alta con 1600 {true: rechaza, null: rechaza, false: acepta, ausente: rechaza}; sin porcentaje, las 3 aceptan; Merge Patch sobre (false, 1600) {true: rechaza, null: rechaza, true + null: acepta}; 8 textos revisados; 0 faltantes |
| V-3 | (a) OpenAPI válido | `make verify-openapi` | rc 0. «docs/contratos/openapi.yaml: OK» (openapi-spec-validator 0.9.0) |
| V-4 | (a) `build.py --check` | `python3 -I backend/internal/api/c1gen/build.py --check` | rc 0. 21 «IDÉNTICO» y «RESULTADO: los 21 archivos coinciden byte a byte con el generador». Entradas: 00001 a 00004 (00004 sha256 c2aa03512be5…). `openapi.yaml`: 703,105 bytes, sha256 09f8a6c63877…. `routes.go`: 26,202 bytes |
| V-5 | Verificación global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=…/datos/tmp/rev-c1-r3/va` | rc 0. «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». C1 con 33 OK; U1 63 OK |
| V-6 | (d) versión | `grep` y `yaml.safe_load` | `info.version` 1.3.0. `C1_VERSION = '1.3.0'` en `gen4_build.py` l. 6. `routes.go` l. 1: «contrato C1 v1.3.0». Los 17 extractos: `info.version` 1.3.0 |
| V-7 | (e) conteos | script propio sobre `openapi.yaml`; `grep -c '{Method:'` en `routes.go`; `grep` en `INDICE.md`, `C1-datos.md` y `C1-puerta2.md` | **168 rutas, 212 operaciones y 255 esquemas**. `routes.go`: 212 entradas. Lo mismo en `extractos/INDICE.md` l. 7, `C1-datos.md` l. 9 y l. 627, y la tabla de `C1-puerta2.md` l. 13. Consistente |
| V-8 | M-1: el CHECK real | base propia (Up de 00001 a 00004 con `awk` y `sqlite3`), `INSERT`/`UPDATE` de un `store_credit` dentro de `BEGIN … ROLLBACK` | Alta: (1, 1600), (NULL, 1600) y sin la columna con 1600 → «CHECK constraint failed». (0, 1600) y las 4 sin porcentaje → aceptadas. Sobre (0, 1600): `vat_included = 1` y `= NULL` → rechazadas; `= 1` con porcentaje NULL y `= NULL` con porcentaje NULL → aceptadas. Sobre (NULL, NULL): `vat_on_interest_bp = 1600` → rechazada; junto con `vat_included = 0` → aceptada |
| V-9 | M-1: los textos del contrato | script propio (`yaml.safe_load`) sobre `CreditCreate`, `CreditUpdate`, `createCredit` y `updateCredit` | Los 4 textos dicen la regla completa: «vat_on_interest_bp con valor solo se admite con vat_included = false» y «… y vat_included = true, null o ausente, la API responde 422 validation_failed … (field = vat_on_interest_bp) y no guarda nada». Traen los ejemplos del alta sin `vat_included` y de `{"vat_included": null}` sobre (false, 1600), y «tampoco deduce false a partir del porcentaje». Las dos operaciones tienen respuesta 422. Cada combinación que rechaza la base en V-8 queda cubierta, también `{"vat_on_interest_bp": 1600}` sobre (NULL, NULL), por «el cuerpo sobre lo ya guardado». Solo `createCredit` y `updateCredit` escriben esos campos: ningún otro esquema tiene `vat_included` ni `vat_on_interest_bp` |
| V-10 | (b) claves nuevas de S1 en settings | script propio: claves de `app_setting` de mi base, propiedades de `parameters` del 200 del `get…`, del cuerpo y del 200 del `update…`, con `$ref` resueltos | Las 10 claves nuevas de la 00004 (l. 66–75) y las 3 de `@param_rango` (l. 76–78) están en las 3 posiciones: `reports.*` (5) en `getReportSettings` y `updateReportSettings`; `alerts.send_time`, `alerts.stale_prices_advance_days` y `alerts.meal_plan_missing_weekday` en `get`/`updateAlertSettings`; `optimizer.date_window_days` y `comparator.timeout_seconds` en `get`/`updateOptimizerSettings`; `ocr.timeout_seconds` en `get`/`updateDataSettings`; `price.*` en `get`/`updatePriceSettings`. Todas `required` en el `get`, con su decisión citada (Q2-72, Q2-73, Q2-75, Q2-76, Decisión F, Q2-66 B, K47) y su valor inicial, igual al sembrado (21:30, 7, 19:00, 1, 08:00, 09:00, 1, 7, 7 y 30). Están en el `enum` de `reset_parameters` |
| V-11 | (b) rangos del `PATCH` | el mismo script | Horas: `pattern ^([01][0-9]\|2[0-3]):[0-5][0-9]$`; días de la semana 1–7; día del mes 1–31; `comparator.timeout_seconds` 1–30; `price.stale_after_days` ≥ 1; anticipación, ventana, OCR y peso cero ≥ 0. Coinciden con los disparadores de la 00004 (l. 113–151) |
| V-12 | (b) `NotificationUpdate` | script propio | `NotificationUpdate`: solo `is_enabled` y `params_json`. `Notification`: sin `send_time` ni `advance_days`. `params_json` sin propiedades declaradas, y su texto remite a `getReportSettings` y `getAlertSettings` |
| V-13 | (c) `vat_included` | script propio | `boolean` en `Credit` y `CreditCreate`; `boolean` o `null` en `CreditUpdate`; descripción y `x-provisional` (`S3-iva-lectura`). No es obligatorio en `CreditCreate` (`required` = name, kind, holder_person_id y payment_frequency), y la regla 422 cubre ese caso (V-9) |
| V-14 | (c) alcance del reparto | script propio | `CostShareRuleUpdate.scope`: `enum [onward, all]`, sin valor por omisión, «Obligatorio si viene parts y solo se acepta con parts (si no, 422 validation_failed)», `x-provisional` → `C1-reparto-alcance` (abierto). `updateCostShareRule`: 422 y `x-tablas` con `expense`, `expense_item` y `expense_share` |
| V-15 | (c) `assigned_store_id` | script propio | `ShoppingListItem.assigned_store_id` documentado (Q2-39 a Q2-42 A, «Ausente = sin tienda»). `groups[]` con `purchase_channel`, `store_id`, `visit_weekday`, `visit_date` y `total_cents`; orden: «primero los de mercado y luego los de cadena; dentro de cada vertical … con el grupo sin tienda al final» |
| V-16 | (c) acumulados de `AccrualVsCash` | script propio | `months[].required` incluye `cumulative_accrued_cents`, `cumulative_cash_cents` y `cumulative_difference_cents`, cada uno con descripción y origen (Q2-65 A). M-11 de `C1-datos.md` no cambió desde la ronda 2 |
| V-17 | Extractos | script propio que normaliza `componentes-comunes.yaml#/` → `#/` | 17 YAML, todos ≤ 80,000 bytes (mayor `ajustes-c.yaml`, 79,094; `INDICE.md` l. 27 dice lo mismo). 212 de 212 operaciones, sin repetidas. Las 212 operaciones y las 304 copias de esquemas son **idénticas** a `openapi.yaml`. `ajustes-c.yaml` trae la regla 422 en sus 4 textos |
| V-18 | C1-15: marcas PROVISIONAL | script propio sobre todo `openapi.yaml`, y `grep` en `bloqueos/` | 6 `x-provisional: true`: `scope` → `C1-reparto-alcance`; `vat_included` (×3) y `missing_terms` → `S3-iva-lectura`; `JobUpdate.schedule` → `C1-hora-reportes-job`. Los 3 bloqueos existen y no tienen «Respuesta de Axel». `C1-datos.md` §13 (l. 1037–1044) suma `S1-reparto-edicion-gasto` (abierto) para `updateExpense` (ver N-1) |
| V-19 | «Solicitudes a C1» de `contratos-puerta2.md` | lectura de las descripciones crudas | Se siguen cumpliendo: `estimated_price_cents` («Ausente = sin precio: no suma en los totales»), `recommended_store_id`, `cheapest_store_id` (ausente si ninguna está completa), `difference_vs_cheapest_cents` (solo tiendas completas) e `identical_count` («Renglones idénticos (no piezas; S5 §7.2)») |
| V-20 | Alcance de la corrección r2 | `diff` de los respaldos `/tmp/*.bak-cc1r2` contra las fuentes actuales | `gen2_schemas.py`: solo `CREDIT_VAT_RULE` y `CREDIT_VAT_RULE_OP` (24 líneas). `verify_c1.py`: solo la comprobación nueva (50 líneas). `README.md`: 1 renglón. `gen3_ops.py`: sin cambio. Leí el código de la comprobación (l. 457–503): prueba contra la base dentro de `SAVEPOINT`, y un error que no sea «CHECK constraint failed» la hace fallar |

## Hallazgos

No hay altos ni medios.

### N-1 (baja) · `updateExpense` en `openapi.yaml` no lleva la marca que sí trae `C1-datos.md` §13
- **HECHO.** `C1-datos.md` l. 1044: «Edición de un gasto anterior después de “solo los nuevos” · `updateExpense` recalcula sus partes con la regla vigente (opción en uso, RN-03 1.4); PROVISIONAL · `S1-reparto-edicion-gasto` (abierto)».
- **HECHO.** En `openapi.yaml`, `updateExpense` solo dice «recalcula las partes». No tiene `x-provisional` ni cita `S1-reparto-edicion-gasto` (0 apariciones en `openapi.yaml`).
- **Por qué es bajo.** C1-15 admite la marca «`x-provisional: true` o la marca PROVISIONAL en `C1-datos.md`», y los dos textos no se contradicen. Pero quien implemente S6 desde el extracto no ve que la regla es provisional.
- **Ruta:** `backend/internal/api/c1gen/gen3_ops.py` (texto de `updateExpense`).
- **Corrección sugerida:** agregar a `updateExpense` «con la regla vigente (PROVISIONAL: bloqueo S1-reparto-edicion-gasto)», como hoy hace `updateCostShareRule` con `C1-reparto-alcance`.

### N-2 (baja) · las descripciones de un campo booleano usan los valores de la base
- **HECHO.** `vat_included` es `boolean` en el JSON, pero su descripción (de la `@col_note` de S1) dice «1 = sí … 0 = no … NULL = falta capturar» y «Con 0 y vat_on_interest_bp NULL». La regla 422 del mismo esquema usa `true`, `false` y `null`.
- **HECHO.** La descripción de `vat_on_interest_bp` (en `Credit`, `CreditCreate` y `CreditUpdate`) todavía empieza con «NULL = la tasa o la cotización ya incluye IVA (Kueski)». Desde 1.4 eso solo vale con `vat_included` = 1 (B-2 de la ronda 1).
- **Ruta:** `backend/internal/api/c1gen/gen2_schemas.py` (texto que toma de la nota de columna).
- **Corrección sugerida:** en el contrato, decir `true`, `false` y ausente o `null`, y quitar la frase de 1.3 de `vat_on_interest_bp`, o anteponerle «Hasta 1.3:».

### B-1 a B-5 de la ronda 2 (bajos) · siguen abiertos, sin cambio
- **B-1:** la tabla «Cambios» de `C1-datos.md` (l. 16–28) no tiene las filas «C1 v1.3 · corrección r1» ni «corrección r2», y el punto de IVA (l. 36) no menciona la regla 422. Ruta: `c1gen/c1_plantilla.md`.
- **B-2:** `C1-datos.md` l. 38 sigue con «Orden de los grupos: mercado, cadena y el grupo sin tienda al final». `openapi.yaml` pone el grupo sin tienda al final de **cada vertical**. Ruta: `c1gen/c1_plantilla.md`.
- **B-3:** `params_json` sigue con `additionalProperties: true` y no dice qué pasa si un cliente de v1.2 manda `send_time`, `advance_days`, `weekday` o `day_of_month`. Comprobé que la base no lo rechaza (solo `CHECK (json_valid(params_json))`, sin disparadores), así que no sale un 5xx: se guardaría sin efecto.
- **B-4 (fuera de C1):** el `Makefile` sigue con `TMP := $(TMP_ROOT)/verify-wf0` (l. 30).
- **B-5:** siguen los bajos de la ronda 1:
  - `docs/criterios/C1.md` sigue en 1.2 (C1-13 dice «C1 v1.2») y no registra la comprobación C1-04 del IVA;
  - `bloqueos/INDICE.md` no tiene `C1-reparto-alcance`, `C1-hora-reportes-job` ni `S3-iva-lectura` (`grep -c` da 0 en los tres);
  - `S3-iva-lectura` no anota lo que tomó C1;
  - `C1-puerta2.md` l. 84 sigue con «comprobaciones OK: 29». Hoy son 33.

## Estado de los hallazgos de la ronda 2

| Hallazgo r2 | Estado | Evidencia propia |
|---|---|---|
| M-1 (media), C1-04 y Q2-67 opción 1: el CHECK también rechaza un porcentaje con `vat_included` null o ausente | **CERRADO** | V-8 y V-9. Escenario 1 (alta sin `vat_included` con 1600): la base lo rechaza, y C1 dice 422 `validation_failed` con `fields[].field = vat_on_interest_bp`, que no guarda nada y que no deduce `false`. Escenario 2 (`{"vat_included": null}` sobre (false, 1600)): la base lo rechaza, y C1 da la misma respuesta, con el ejemplo literal. La comprobación nueva de `verify_c1.py` lo exige contra la base (V-2) |
| B-1 a B-5 (bajos) | Abiertos, sin cambio | Ver arriba. No venían en el encargo |

## Observaciones que no son hallazgo
- **HECHO.** `S1-puerta2.md` l. 29 copia un CHECK anterior: `(vat_included = 0 OR vat_on_interest_bp IS NULL)`. Con NULL ese CHECK aceptaría (NULL, 1600). La 00004 vigente usa `vat_included IS 0` y lo rechaza (V-8). C1 sigue a la migración, que es lo que manda. Es un reporte de S1 desactualizado: lo paso como solicitud fuera de alcance.
- **HECHO.** `Credit.missing_terms` lista `vat` con (0, sin porcentaje), con PROVISIONAL y el bloqueo abierto `S3-iva-lectura` (pregunta 2, «Sin recomendación»). Es la lectura en uso de S3 1.5 §5.3, donde ese caso responde `ErrMissingCapture`, y la marca cita el bloqueo abierto. No la cuento como una decisión de C1. Las rondas 1 y 2 la aceptaron igual.
