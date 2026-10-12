# Revisión de C1 v1.3 (puerta 2) · ronda 1 de 3

Revisor: `revisor-C1-r1`, independiente: no escribí nada de lo que reviso y no corregí nada. Flujo P2c, 11-oct-2026, noche (el reloj del equipo ya marcaba 12-oct 00:0x).

**Qué reviso.** C1 v1.3 del agente `C1-puerta2` (`docs/estado/p2c/C1-puerta2.md`):
- `docs/contratos/openapi.yaml`, `C1-datos.md`, `backend/internal/api/routes.go` y `docs/contratos/extractos/`, todos generados por `backend/internal/api/c1gen/`;
- los bloqueos nuevos `C1-reparto-alcance` y `C1-hora-reportes-job`.

**Contra qué.** Solo contra:
- `docs/criterios/C1.md` 1.2;
- `PUERTA-2.md` §2 y §3, más §4 para Q2-72 a Q2-77, porque §3.1 remite a ellas;
- `docs/estado/p2c/S1-puerta2.md` (claves);
- `docs/estado/p2c/contratos-puerta2.md` («Solicitudes a C1»);
- `RESPUESTAS.md`: modo «lo recomendado + excepciones», Q2-11 B y puerta 2 bis.

## Veredicto
**NO APROBADO: 0 altos, 3 medios y 5 bajos.**

Lo que está bien:
- El fondo de la puerta 2 está aplicado y las verificaciones pasan: las 10 claves de S1 con su rango, `params_json` limpio, versión 1.3.0, 212 operaciones y 168 rutas.
- No hay ninguna respuesta de Axel contradicha ni ninguna «sin recomendación» decidida.

Lo que falla:
- Dos propiedades salieron del generador con un guion bajo de más (H-1).
- M-11 no nombra los acumulados que pide la opción A aprobada de `U3-t2-acumulados` (H-2).
- La combinación nueva de IVA no dice qué responde la API (H-3).

## Verificaciones (todas reejecutadas por mí; scripts en `datos/tmp/rc1r1/`, ya borrados)
| Id | Verificación | Comando (resumen) | Resultado |
|---|---|---|---|
| a-1 | Generador idéntico y C1 | `make verify-c1` | rc 0. 21 archivos «IDÉNTICO» y «RESULTADO: los 21 archivos coinciden byte a byte con el generador». «comprobaciones OK: 29 · RESULTADO: TODO OK» |
| a-2 | `verify_c1.py` con mi propia base | `goose -dir backend/db/migrations sqlite3 datos/tmp/rc1r1/c1.db up` y después `uv run … verify_c1.py --db datos/tmp/rc1r1/c1.db` | goose: «successfully migrated database to version: 4». verify: rc 0, 29 líneas OK, 0 FALLA y «RESULTADO: TODO OK». Además: C1-01 «3.1.1 · rutas = 168 · operaciones = 212»; C1-13 «info.version: 1.3.0»; K44 «@param: 25 … fallan: [] · sembrados con ese valor: 25»; C1-15 «marcas PROVISIONAL en openapi.yaml: 6 · bloqueos abiertos citados: [C1-hora-reportes-job, C1-reparto-alcance, S3-iva-lectura] · sin bloqueo abierto: []»; EXT 212/212; S6-03 212/212; C1-23 «claves de puerta 2 con rango distinto del aprobado: [] · patrón acepta 00:00–23:59 y rechaza 24:00, 7:30, 12:60: True» |
| a-3 | OpenAPI válido | `make verify-openapi` (openapi-spec-validator 0.9.0) | «docs/contratos/openapi.yaml: OK» |
| a-4 | Verificación global | `make verify-all GOOSE=… TMP_ROOT=datos/tmp/rc1r1/tr` | rc 0 y «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». El primer intento, sin `TMP_ROOT`, dio rc 2 por «Error near line 1246: disk I/O error» en `verify-sql`: chocó con el `verify-all` que corría `rev-s5-p2c-r1` sobre la misma carpeta `datos/tmp/verify-wf0`. El suyo terminó con `EXIT=0`. No es defecto de C1 (ver Propuestas en `revisor-C1-r1.md`) |
| b-1 | Cada clave nueva de S1 en settings, con su rango | script `chk2.py`: busca `x-app-setting` en todos los esquemas | Las 10 claves aparecen. Las 5 `reports.*` están en `ReportSettings` (lectura, `SettingParameterTime` o `Integer`) y en `ReportSettingsUpdate`, con patrón HH:MM, día de la semana de 1 a 7 y día del mes de 1 a 31; las usan `getReportSettings` y `updateReportSettings`. `alerts.send_time` (patrón), `alerts.stale_prices_advance_days` (≥ 0) y `alerts.meal_plan_missing_weekday` (1 a 7) están en `AlertSettings`. `optimizer.date_window_days` (≥ 0) está en `OptimizerSettings` y `ocr.timeout_seconds` (≥ 0) en `DataSettings`. Rangos nuevos: `comparator.timeout_seconds` de 1 a 30 (Q2-74) y `price.stale_after_days` ≥ 1 (Q2-77). Todo coincide con la tabla de `S1-puerta2.md` l. 13 a 24 |
| b-2 | Valores sembrados | `sqlite3 c1.db "select setting_key, value_json from app_setting …"` | 21:30, 7, 19:00, 1, 08:00, 09:00, 1, 7, 7, 30, 5, 30 y 90: son Q2-72, Q2-73, Q2-75, Q2-76 y la decisión F. 30 claves en total |
| b-3 | `NotificationUpdate` sin `send_time` ni `advance_days` | `chk3.py` | `params_json` de `Notification` y de `NotificationUpdate`: `properties` vacío y sin `x-provisional`. `send_time` solo sale en el texto que explica dónde vive ahora. Lo que sembré: 7 avisos con `{}`, salvo `budget_threshold`, que queda con `{"scopes": [...]}` |
| c-1 | `vat_included` | `chk3.py` | Está en `Credit` (boolean, no obligatorio), `CreditCreate` (boolean) y `CreditUpdate` (boolean o null). `missing_terms` trae un enum de 5. Lo documentan C1-datos l. 36 y §13 l. 1043. Ver H-3, B-1, B-2 y B-3 |
| c-2 | Alcance del reparto | `chk3.py` | `CostShareRuleUpdate.scope`: enum `[onward, all]`, sin `default`, `x-entrada`, 422 si llega sin `parts`. Lo documenta C1-datos §7.2 (l. 313 a 325), con su tabla y un ejemplo sintético |
| c-3 | `assigned_store_id` | `chk4.py` | `ShoppingListItem.assigned_store_id` (calculado), `recommended_store_id` y `estimated_price_cents` con los textos de la solicitud. En `groups[]`: `store_id`, `visit_weekday` y `total_cents` con su texto, y orden mercado → cadena → sin tienda. **La fecha salió como `visit_date_`** (H-1) |
| c-4 | Acumulados de `AccrualVsCash` | `chk4.py` y `grep` | `months[].required` incluye `cumulative_accrued_cents`, `cumulative_cash_cents` y `cumulative_difference_cents`, con descripción «mismo origen». C1-datos solo los nombra en «Cambios» (l. 40). **M-11 (l. 372) no los nombra** (H-2) |
| d | Versión | `chk.py` y `grep` | `info.version` pasa de 1.2.0 (HEAD) a 1.3.0. `C1_VERSION = '1.3.0'` (`gen4_build.py` l. 6). C1-datos dice «C1 v1.3 … `info.version` 1.3.0» y `routes.go` dice «contrato C1 v1.3.0» |
| e | Conteos | `chk.py` (HEAD contra el árbol de trabajo) y `chk6.py` (extractos) | Rutas, operaciones y esquemas pasan de 167, 210 y 252 a 168, 212 y 255. Lo nuevo: la ruta `/api/v1/settings/reports`, las operaciones `getReportSettings` y `updateReportSettings`, y los esquemas `ReportSettings`, `ReportSettingsUpdate` y `SettingParameterTime`. No se quitó nada. 30 esquemas cambiados. Extractos: 212 operaciones, 0 faltantes, 0 sobrantes, 0 repetidas y 0 archivos de más de 80 000 bytes. Coincide con el reporte y con `extractos/INDICE.md` («212 operaciones en 168 rutas y 255 esquemas») |
| f | Marcas | `chk5.py` y `grep -c "Respuesta de Axel"` | Hay 6 `x-provisional`: `CostShareRuleUpdate.scope`, `Credit.vat_included`, `Credit.missing_terms`, `CreditCreate.vat_included`, `CreditUpdate.vat_included` y `JobUpdate.schedule`. Los bloqueos que citan (`C1-reparto-alcance`, `C1-hora-reportes-job`, `S3-iva-lectura`) y `S1-reparto-edicion-gasto` dan 0, o sea que siguen abiertos. Los 7 bloqueos ya respondidos (`S1-parametros-avisos`, `U4-iva-credito`, `U4-reparto-retroactivo`, `U3-t2-acumulados`, `U4-dia-mes-reporte`, `S5-reglas-comparativa` y `S4-asignacion-mercado`) dan 1, y ninguno aparece citado en `openapi.yaml` |
| g | Respuestas de Axel | lectura de PUERTA-2 §2 a §4 contra los esquemas | Q2-11 B, Q2-17 C (`progress_tenths_pct`), Q2-33 (b) (`Scenario.store_id`), Q2-39 a 42 A, Q2-02, 05 y 08 (orden de `stores[]`, diferencia solo en las completas, `identical_count` en renglones), Q2-65 A (campos), Q2-66 B, Q2-67 opción 1, Q2-68 C, Q2-72 a 77 y Q2-86 A: aplicadas. No hay decisiones sobre preguntas «sin recomendación»: Q2-94 no se agregó y Q2-07, 48, 50 a 52 y 56 a 58 no tocan C1 |

## Hallazgos

### H-1 (media) · dos propiedades con un guion bajo de más (C1-14 y la solicitud de contratos)
- **Qué pasa.** `openapi.yaml` l. 13517 publica `ShoppingList.groups[].visit_date_` y l. 13718 publica `BasketComparison.cheapest_store_id_`. Lo mismo sale en `extractos/compras-tiendas-a.yaml` l. 797 y 998 y en `extractos/reportes.yaml` l. 802.
- **Con qué choca.** Todo lo demás usa los nombres sin guion bajo, `visit_date` y `cheapest_store_id`:
  - C1-datos l. 38 y 39;
  - el reporte `C1-puerta2.md` (puntos 5 y 8);
  - las «Solicitudes a C1» de `contratos-puerta2.md` l. 65 a 72;
  - las descripciones de los campos vecinos: `visit_weekday` dice «Día ISO de visit_date», el orden de `groups` dice «por visit_date ascendente» y `difference_vs_cheapest_cents` dice «total de cheapest_store_id».
- **Por qué importa.** C1-14 pide «propiedades en inglés corto igual que `nombre_codigo`». Quien implemente S6 o U3/U4 a partir de C1-datos leería un campo que la API no entrega. Ninguna verificación lo detecta: `verify_c1.py` y `build.py --check` pasan.
- **Causa.** `backend/internal/api/c1gen/gen2_schemas.py` l. 407 (`visit_date_=`) y l. 420 (`cheapest_store_id_=`) usan esos nombres como argumentos de `_ins`, y `_ins` los copia tal cual.
- **Corrección sugerida.**
  1. En `gen2_schemas.py`, cambiar los argumentos a `visit_date=` y `cheapest_store_id=`.
  2. Regenerar con `python3 -I backend/internal/api/c1gen/build.py`.
  3. Agregar a `verify_c1.py` una comprobación: 0 propiedades que terminen en `_`, y cada campo que nombra «C1 v1.3 (puerta 2)» de C1-datos existe en `openapi.yaml`.

### H-2 (media) · M-11 no nombra los acumulados (Q2-65 A y C1-09)
- **Qué pide la opción aprobada.** La opción A de `docs/estado/bloqueos/U3-t2-acumulados.md` (l. 10 y 25, aprobada en bloque) dice «C1-datos M-11 los nombra». También da las cifras para comprobarlos con cualquier reloj: nov-2026 2,921.25 y 3,895.00; mar-2027 13,103.78 y 14,362.48 (R-52 y R-51).
- **Qué hay.** La fila M-11 de `docs/contratos/C1-datos.md` (l. 372; viene de `c1gen/c1_plantilla.md` l. 268) sigue diciendo «por mes: devengado y efectivo (totales y partes); diferencia acumulada = …». No nombra `cumulative_accrued_cents` ni `cumulative_cash_cents`, ni el origen «desde el primer evento», ni un ejemplo.
- **Por qué importa.** C1-09 pide métricas con fórmula y ejemplo del oráculo. Sin la fila, quien dibuje T-2 en WF-3 solo tiene la nota de «Cambios» (l. 40).
- **Corrección sugerida.** En `c1_plantilla.md`, ampliar M-11:
  - la fórmula de los dos acumulados (al cierre del mes, desde el primer evento, mismo origen que la diferencia, y diferencia = devengado acumulado − efectivo acumulado);
  - un ejemplo con cifras de `oraculo/model.json`: oct-2026 1,947.50 y 973.75 (R-45) y una ventana que empiece en nov-2026, con las cifras del bloqueo comprobadas contra el oráculo, no copiadas a mano.

  Después, regenerar.

### H-3 (media) · no está escrito qué responde la API a «incluye IVA» con porcentaje (Q2-67 opción 1)
- **Qué pasa.** La 00004 agrega `CHECK (vat_included IN (0, 1) AND (vat_included = 0 OR vat_on_interest_bp IS NULL))` (`S1-puerta2.md` l. 29). C1 describe la regla en `vat_on_interest_bp` («con `vat_included` = 1 es NULL»). Pero ni `CreditCreate`, ni `CreditUpdate`, ni `createCredit` o `updateCredit` dicen qué responde la API a `{"vat_included": true, "vat_on_interest_bp": 1600}`.
- **El caso del Merge Patch.** Tampoco dicen qué pasa con `{"vat_included": true}` sobre un crédito que ya tiene 1600: ¿422 `validation_failed` con detalle en `vat_on_interest_bp`, o la API vacía el porcentaje?
- **Por qué importa.** Es lo que S6 tendría que adivinar. Si no lo mapea, la violación del CHECK puede salir como 5xx.
- **Referencia.** Para `price.zero_weight_after_days` ≥ `price.stale_after_days` y para `scope` sin `parts`, el mismo agente sí escribió el 422.
- **Corrección sugerida.** En `gen2_schemas.py` (bloque de `vat_included`, l. 312 en adelante), escribir la regla técnica en las descripciones de `CreditCreate` y `CreditUpdate` y en `createCredit` y `updateCredit`: «`vat_included` = true con `vat_on_interest_bp` (en el cuerpo o ya guardado) → 422 `validation_failed` con detalle por campo; la API no cambia el otro campo por su cuenta». Es la convención de C1 de no sustituir valores. Después, regenerar.

### B-1 (baja) · marca de más en `vat_included`
- **Qué pasa.** Las 3 marcas `x-provisional` de `Credit`, `CreditCreate` y `CreditUpdate.vat_included` cubren un campo FIRME (Q2-67 opción 1). Lo que está en duda es solo `missing_terms` con 0 y sin porcentaje, y `missing_terms` ya lleva su propia marca.
- **Corrección sugerida.** Dejar la marca solo en `missing_terms` y el texto PROVISIONAL en la descripción, para que C1-15 no lo lea como marca sobre algo aprobado.
- **El caso de `scope`.** Es distinto: no hay otra propiedad donde poner la marca, y la descripción ya acota qué está en duda.

### B-2 (baja) · vocabulario de la base en `vat_included`
- **Qué pasa.** La descripción de `vat_included` habla como la base: «1 = sí … 0 = no … NULL = falta capturar». C1-datos l. 36 y la convención 6 (l. 78) dicen `true`, `false` y ausente.
- **Corrección sugerida.** Redactarla con los valores de JSON.

### B-3 (baja) · `S3-iva-lectura` no registra lo que tomó C1
- **Qué tomó C1.** Para la pregunta 2 (sin recomendación), la parte de C1 de la opción A: listar `vat` con `vat_included` = 0 y sin porcentaje. Va PROVISIONAL, y es coherente con el `ErrMissingCapture` que S3 1.5 tiene en uso.
- **Con qué choca.** RN-26 1.4 dice «solo mientras `vat_included` sea NULL», y el «Cómo seguí» del bloqueo dice «deja sin decidir si C1 lo avisa».
- **Pregunta 1.** La descripción de C1 no menciona la pregunta 1 (préstamo con NULL; en uso A, sin aviso).
- **Corrección sugerida.** Pedir al integrador una línea en «Cómo seguí» de `S3-iva-lectura` y en RN-26 que registre la lectura de C1 v1.3. En C1, citar también la pregunta 1 en la descripción de `vat_included`.

### B-4 (baja) · criterios desfasados (fuera del alcance de C1; es la solicitud 4 del implementador)
- **Qué pasa.** `docs/criterios/C1.md` sigue en 1.2:
  - C1-13 pide «C1 v1.2 (`info.version` 1.2.0)»;
  - C1-21 dice «hoy 210»;
  - C1-22 conserva la excepción PROVISIONAL de `send_time`, `advance_days`, `weekday` y `day_of_month` con `S1-parametros-avisos` abierto.
- **Por qué importa.** La regla de C1-22 («si el bloqueo trae Respuesta de Axel, se aplica sin excepción») cubre el fondo, pero la r2 y la r3 evaluarían contra un texto que contradice lo correcto.
- **Corrección sugerida.** El agente de criterios o el integrador pasa C1.md a 1.3: v1.3 y 1.3.0, 212 operaciones, C1-22 sin la excepción y con las claves de K47, y los rangos de Q2-74 y Q2-77.

### B-5 (baja) · `bloqueos/INDICE.md` no lista los dos bloqueos nuevos de C1
- **Qué pasa.** Faltan `C1-reparto-alcance` y `C1-hora-reportes-job`. `grep` da 0.
- **Corrección sugerida.** Es la solicitud 1 del implementador al integrador: agregarlos.

## Estado de los anteriores
- Es la primera ronda de C1 en P2c, así que no hay hallazgos anteriores de esta ronda que seguir.
- De rondas previas: los hallazgos de P1c y P1d (`docs/estado/p1c/REVISION-C1-r1.md` y `r2`) siguen cerrados. Las marcas de `C1-valores-k31` dan 0 (C1-23) y no queda ninguna marca que cite un bloqueo respondido.

## Fuera de esta revisión (no son hallazgos)
- **K53 (puerta 2 bis, se registra en P2d).** C1 todavía tiene la lectura asistida (`assisted_reading_enabled`, `PriceSources.sites`, el 409 de `updateStore`; 6 menciones de `S5-terminos.md` en `openapi.yaml`). Tampoco tiene la importación del agente externo. Lo declara el reporte del implementador y es trabajo de P2d.
- **Ubicación de `ocr.timeout_seconds`.** Su lugar en Ajustes («Datos y respaldos») es una PROPUESTA de S1 sin pregunta en la puerta. C1-datos §11.7 lo dice («sección: propuesta S1»). Cambiarlo solo movería la clave a otra operación de settings.
