# Revisión S1 · esquema 1.3 · ronda 2 de 3 (flujo P1c, 11-oct-2026)
Revisor independiente: revisor-S1-r2. No escribí nada de lo revisado ni corregí nada. Solo evalué contra `docs/criterios/S1.md` y `docs/plan/decisiones.md` §K (K36, K37, K39, K40, K43 y K44). Reejecuté todas las verificaciones con mis propios scripts en `datos/tmp/rev2/`, que ya borré.

## Veredicto
**NO APROBADO.** Hay 0 hallazgos altos, 1 medio y 4 bajos.
- Los dos hallazgos de la ronda 1 quedaron **cerrados**: H-01 (alta, K44) y H-02 (media, S1-13).
- Pasan las 10 verificaciones técnicas, de la (a) a la (j).
- Lo que impide aprobar es M-01. El inventario de parámetros de K44 no menciona tres valores de la decisión F que siguen en `notification.params_json`: domingo (`weekday` 7, dos veces) y día 1 (`day_of_month` 1). No están en el inventario, ni en el bloqueo, ni llevan marca. Su caso es el mismo que el de `days_before` (decisión F), que la corrección sí pasó a `app_setting`. También es el mismo que el de `send_time`, que quedó en el bloqueo.

## Verificaciones reejecutadas (a–j)
| # | Verificación | Resultado | Evidencia (mis comandos) |
|---|---|---|---|
| a | 00001 + 00002 + 00003 en `datos/tmp/rev2/rev.db` | Cumple | Script propio con `sqlite3` de Python: `integrity_check` = ok; `foreign_key_check` = []; 55 tablas, 55 STRICT (con `pragma_table_list`); 152 índices y 75 disparadores |
| b | El Down de la 00003 deja la 1.2 | Cumple | Hice una base 1.2, la copié, apliqué Up 00003 (282 objetos) y luego Down. `sqlite_master` quedó idéntico (278 filas, con los autoíndices) y los datos de las 55 tablas también (21 filas sembradas). Con datos de 1.2 cambiados (`days_before` [3, 1, 0], umbrales [70, 95], `scopes` ["total"]), el Up siembra 70, 95 y `[3, 1, 0]` y el Down repone los valores y `updated_at`. Solo cambian los espacios del JSON de `scopes` (ver B-01) |
| c | goose up / down-to 0 / up | Cumple | `goose` v3.28.0 sobre `datos/tmp/rev2/goose.db`: códigos 0, 0 y 0; tablas 55, 0 y 55; `version` = 3; `integrity_check` ok; las claves `alerts.*` quedan en 80, 100 y `[1, 0]`, y las `health.*` en 0 y 0 |
| d | El DDL de 00001 y 00002 no cambia | Cumple | `git diff --quiet HEAD -- 00001_esquema.sql 00002_puerta1.sql` dio exit 0: son idénticos byte por byte, con comentarios y todo. `s1_verifica` DDL: True y True |
| e | `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` | Cumple | Exit 0, **RESUMEN TODO OK** con 28 líneas. KUESKI: suma de pagos 389500 y parte de Axel 239528, 0 diferencias. También pasan S1-02, S1-13 (15 de 15), S1-16, K44, AVISOS13, DDL, DOWN13 (130 filas), DOWN y GOOSE |
| f | `gen_dict.py` regenera el mismo diccionario | Cumple | `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev2/dic.md` y luego `cmp` contra `docs/contratos/diccionario-datos.md`: idéntico. Parámetros 15, sin valor 6, sin fila [], que no coinciden [] |
| g | Valores de K44 contra §K | Cumple | En la base 1.3 (alcance de la casa), los valores esperados por mí coinciden 15 de 15: 80, 100, 30, 90 (peso cero), 90 (sesión), 5, 15, 10 MB, 5 s, 0, 0, 24 h, 12 (K6), 1 GB (K6) y `[1, 0]` (F). No hay claves nuevas inesperadas: 1.3 tiene 20 claves, la 1.2 tenía 10 y 00003 siembra 10. No hay ninguna clave `optimizer.*` de K16, `send_time` ni `advance_days` sembrada: están listadas «sin valor» |
| h | Ninguna marca PROVISIONAL cita un bloqueo resuelto | Cumple | Crucé cada línea con PROVISIONAL de 00001–00003, de las herramientas, del diccionario y de las reglas contra todos los archivos de `docs/estado/bloqueos/`. Solo se cita `S1-parametros-avisos`, que no tiene «Respuesta de Axel». Las menciones de bloqueos resueltos en las reglas (línea 11, «Cambios en 1.3», y `S1-meta-kcal` en RN-17) son historia y no llevan marca |
| i | RN-27 y RN-12 recalculadas | Cumple | Saqué los valores de `model.json` con `Decimal` y ROUND_HALF_UP. `periodos[0]` va de 2026-10-10 a 2026-11-09: efectivo 973.75 × 7 / 31 = **219.88** (diario 31.41). `periodos[4]` va de 2027-02-10 a 2027-03-09: efectivo 2517.38 × 7 / 28 = **629.35** (diario 89.91). Con el diario redondeado darían 219.87 y 629.37, igual a lo que dice RN-27. Despensa de octubre: 439.76. Pesos: (90 − 30) / 61 = **0.9836**, (90 − 60) / 61 = **0.4918** y día 90 = 0 |
| j | Los mínimos de seguridad no están en el esquema | Cumple | El esquema acepta valores extremos: `auth.max_failed_attempts` 0 y 100000, `session.max_age_days` 100000, `auth.password_min_length` 0, `alerts.budget_warning_pct` 500, lista `[]` y `[400, 0, 0]`. Rechaza el tipo o el signo equivocado: −1, 5.0, "5", true, `[1, -1]`, `1` para la lista, `[1.5]` y `x`. Las propuestas (3 a 10, 5 a 60, 1 a 365, 0 a 1,000…) aparecen 0 veces en la 00003, las reglas y el diccionario: solo están en los reportes |

## Criterios de `docs/criterios/S1.md`
| Criterio | Estado | Evidencia |
|---|---|---|
| S1-01 (alta) entidades | Cumple | La 00003 no quita ni renombra tablas; siguen siendo 55 |
| S1-02 (alta) migraciones en limpio | Cumple | Verificaciones (a) y (c) |
| S1-03 (alta) ejemplo numérico | Cumple | RN-27 con cuentas citadas que reproduje en (i). RN-12, RN-17 y RN-29 tienen ejemplos ilustrativos rotulados; (120 − 60) / 91 = 0.6593 es correcto. `s1_verifica` S1-03: 29 reglas |
| S1-04 (alta) sin huérfanas | Cumple | No hay tablas nuevas. Las columnas nuevas tienen `@used_new` (Salud, Ajustes, S4, C1 v1.2 y U4) |
| S1-05 (alta) dinero exacto | Cumple | No hay columnas de dinero nuevas. `goal_adjust_*` es INTEGER, en kcal |
| S1-06 (media) fechas | Cumple | La única coincidencia de `oct-2026` en la 00003 es el comentario de cabecera de la línea 3; ningún dato la usa |
| S1-07 / S1-08 / S1-09 / S1-10 / S1-11 / S1-12 | Cumple | La 00003 no toca esas tablas. `s1_verifica` pasa S1-07 a S1-12, KUESKI y la ida y vuelta (164 valores) |
| S1-13 (media) Ajustes | Cumple (H-02 cerrado) | Script propio: las 15 claves `@param` aparecen en la fila de su sección. Cuenta y roles tiene 4, Tiendas y precios 3, Créditos 1, Datos y respaldos 2, Salud 2 y Notificaciones 3 |
| S1-14 (alta) credenciales | Cumple | En la 00003, `password` solo aparece en el nombre de la clave `auth.password_min_length`. `grep "[0-9]\{8,\}:AA" backend/db/` da 0 |
| S1-15 (media) integridad | Cumple | La 00003 no tiene claves foráneas nuevas. `foreign_key_check` vacío |
| S1-16 (alta) sin inventados | Cumple, con la observación M-01 | (g) y (h): ningún valor inventado; 7 y 1 vienen de la decisión F. Que los parámetros de programación no estén en el inventario es inconsistencia (M-01), no un valor inventado |
| S1-17 (media) idioma | Cumple | Identificadores en inglés (`goal_adjust_lose_kcal`, `alerts.budget_warning_pct`) y textos en español de México |

| Decisión §K | Estado | Evidencia |
|---|---|---|
| K36 | Cumple | Columnas INTEGER NULL con `CHECK (≥ 0)` (NULL = valor de la casa); claves 0 y 0; RN-17 FIRME con `maintain` sin ajuste |
| K37 | Cumple | RN-12 FIRME: compiten `manual`, `ticket_manual` y `ticket_ocr`; `site_assisted` y PROFECO quedan solo de referencia |
| K39 | Cumple | RN-27 FIRME, con un solo HALF_UP; recalculado en (i) |
| K40 | Cumple | RN-12 (90 − días) / 61; `price.zero_weight_after_days` = 90 |
| K43 | Cumple | 10 MB y 5 s sembrados. `send_time` y `advance_days` no tienen valor: el bloqueo es real por el choque con K44 |
| K44 | Cumple en lo sembrado (H-01 cerrado); inventario incompleto (M-01) | 15 claves con su valor y una descripción que cita su decisión; una sola fuente para umbrales y días antes |

## Hallazgos
### M-01 (media) · El inventario de K44 no menciona los valores de programación de la decisión F que siguen en `params_json`
- **Ruta:** `backend/db/migrations/00003_puerta1b.sql:35-57` (bloque `@param` / `@param_sin_valor`), `docs/contratos/reglas-negocio.md` RN-29 (líneas 286-312), `docs/contratos/diccionario-datos.md` §«Parámetros operativos de app_setting (K44)» (líneas 2265-2296, generada) y `docs/estado/bloqueos/S1-parametros-avisos.md`.
- **Qué pasa:** en la base 1.3, `notification.params_json` todavía tiene `weekly_summary` = `{"weekday": 7}`, `meal_plan_missing` = `{"weekday": 7}` y `monthly_report` = `{"day_of_month": 1}`. Son valores numéricos operativos decididos en §F: «resumen semanal (Telegram, domingo)», «reporte mensual PDF (día 1)» y «plan de comidas sin armar (domingo)». El propio diccionario, en su línea 2257, los llama «parámetros ya decididos». K44 pide que el inventario cubra «los demás que inventaríe S1». Aun así, no aparecen en `@param`, ni en «Existen sin valor sembrado o viven en otra tabla», ni en «No son parámetros» de RN-29.
- **Por qué es inconsistente:**
  - La corrección S1-r1 pasó a `alerts.payment_due_days_before` el `days_before` [1, 0] de esa misma decisión F, que también estaba en `params_json`.
  - El bloqueo `S1-parametros-avisos` deja PROVISIONAL dónde vive `send_time`, que también es programación de un aviso.
  - El reporte `correccion-S1-r1.md` (Pendientes y Propuesta 3) dice que «falta decidir» si `weekday` y `day_of_month` son parámetros de K44, pero no abrió bloqueo ni puso marca (protocolo §2).
  - Así, los tres artefactos que dicen inventariar K44 no reflejan una pregunta que el agente ya identificó.
- **Corrección sugerida (sin cambiar valores):** primero, agregar `weekday` (7, en `weekly_summary` y `meal_plan_missing`) y `day_of_month` (1, en `monthly_report`) a la pregunta del bloqueo `S1-parametros-avisos`, porque es la misma que la de `send_time`: `params_json` (K43) o `app_setting` (K44). Segundo, listarlos con `@param_sin_valor` en «viven en otra tabla», con su valor de la decisión F y PROVISIONAL citando ese bloqueo (solo dónde viven). Después, regenerar el diccionario y reflejarlos en RN-29. Otra opción es tratarlos como `days_before`: claves `app_setting` sembradas con 7, 7 y 1 que citen la decisión F, con Down que las reponga. No hay que tocar 00001 ni 00002.

### B-01 (baja) · El Down de la 00003 no devuelve el texto de `params_json` byte por byte si la casa cambió otras claves del aviso
- **Ruta:** `backend/db/migrations/00003_puerta1b.sql:152-158` (Up, `json_remove`) y `:165-178` (Down).
- **Qué pasa:** con `budget_threshold` = `{"thresholds_pct": [70, 95], "scopes": ["total"]}`, el Up deja `{"scopes":["total"]}`, porque `json_remove` compacta el JSON. El Down repone `{"thresholds_pct": [70, 95], "scopes":["total"]}`. El JSON es equivalente pero el texto no es igual. Con los valores iniciales sí vuelve byte por byte, que es lo que afirma la cabecera.
- **Corrección sugerida:** documentar en la cabecera y en RN-29 que con valores personalizados el Down devuelve un JSON equivalente, o usar el Up con el mismo espaciado.

### B-02 (baja) · La descripción de `notification.params_json` en el diccionario empieza con datos que ya no aplican
- **Ruta:** `docs/contratos/diccionario-datos.md:1985`, generada con `@col_note` desde `backend/db/tools/gen_dict.py`.
- **Qué pasa:** la descripción empieza con «Parámetros del aviso (días antes 1 y 0, umbrales 80 y 100)» y después dice que desde 1.3 ya no están ahí. `gen_dict` añade la nota en lugar de reemplazar el paréntesis.
- **Corrección sugerida:** que `@col_note` permita reemplazar el texto de 1.2 en el diccionario, sin tocar el comentario de 00001.

### B-03 (baja) · La línea de versión de las reglas mezcla la 1.3 con la explicación de la 1.2
- **Ruta:** `docs/contratos/reglas-negocio.md:3`.
- **Qué pasa:** la línea «Versión: 1.3 …» todavía dice «Se numera 1.2 para ir a la par del esquema … este documento no tuvo 1.1».
- **Corrección sugerida:** decir que la 1.2 se numeró así para ir a la par (R2-04) o quitar esa frase de la línea de la 1.3.

### B-04 (baja) · `optimizer.weights` (decisión C) no está en el inventario de K44
- **Ruta:** `backend/db/migrations/00003_puerta1b.sql:35-49` y `docs/contratos/reglas-negocio.md` RN-29.
- **Qué pasa:** `optimizer.weights` (`{"lowest_interest": 1, …}`, «Iguales y se normalizan (decisión C)») es un valor numérico operativo que ya está en `app_setting` desde la 00001. Es configurable y está sembrado, así que K44 se cumple en lo mecánico. Pero no aparece en `@param` ni en RN-29. Las otras 4 claves de la 1.2 que tampoco aparecen son interruptores booleanos.
- **Corrección sugerida:** listarlo en el inventario, con tipo «objeto de enteros», decisión C y sección Créditos y condiciones, o explicar en RN-29 por qué queda fuera.

## Estado de los hallazgos de la ronda 1
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H-01 (alta) K44: umbrales 80/100 y `days_before` fuera de `app_setting` | **Cerrado** | Hay 3 claves `alerts.*` sembradas (80, 100 y `[1, 0]`) con una descripción que cita la decisión F y K44, en `@param` y en sus disparadores. `params_json` de 1.3 queda en `{}` y `{"scopes": ["category", "total"]}`. El Down repone los valores (verificaciones b y g). El bloqueo se redujo a `send_time` y `advance_days`, que es el choque real K43 / K44 y no tiene valor. RN-29 y RN-22 citan las claves |
| H-02 (media) S1-13: 6 parámetros fuera de la tabla de 14 secciones | **Cerrado** | Script propio: las 15 claves están en la fila de su sección: `auth.password_min_length` en Cuenta y roles, `comparator.timeout_seconds` en Créditos, `price.zero_weight_after_days` y `price.profeco_import_max_gb` en Tiendas, y `attachment.max_size_mb` e `idempotency.retention_hours` en Datos y respaldos. `gen_dict` sale con 1 si una clave no encuentra su fila |
| H-03 a H-07 (bajas) | Fuera del encargo de la corrección | No las reevalué. H-04 quedó resuelta de paso, según `correccion-S1-r1.md` |
