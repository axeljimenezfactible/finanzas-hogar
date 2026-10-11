# Revisión S1 · esquema 1.3 · ronda 1 de 3 (flujo P1c, 11-oct-2026)

Revisor independiente: `revisor-S1-r1`. Evalúa solo contra `docs/criterios/S1.md` (1.1) y `docs/plan/decisiones.md` K36, K37, K39, K40, K43 y K44. No corrige nada. Todas las verificaciones se reejecutaron con comandos y scripts propios en `datos/tmp/rev-s1-r1/` (borrados al terminar), sin tomar por bueno el reporte `docs/estado/p1c/S1-esquema-1.3.md`.

## Veredicto
**NO APROBADO.** 1 hallazgo alto, 1 medio y 5 bajos. Las pruebas técnicas pasan todas: migraciones, Down, goose, DDL, regresión Kueski, diccionario regenerado idéntico y cifras de RN-12, RN-17 y RN-27. Lo que impide aprobar: (H-01) los umbrales 80 % y 100 % que K44 nombra no son claves de `app_setting` sembradas; (H-02) la tabla de las 14 secciones de Ajustes (S1-13) no recoge 6 de los 12 parámetros del inventario.

## Verificaciones reejecutadas (a–j)
| # | Comprobación | Resultado | Evidencia (comando propio) |
|---|---|---|---|
| a | 00001+00002+00003 en `datos/tmp/rev-s1-r1/rev.db` | OK | `rev.py`: `integrity_check` = ok · `foreign_key_check` = 0 filas · 55 tablas, 0 sin STRICT · `health_profile` termina en `goal`, `goal_adjust_lose_kcal`, `goal_adjust_gain_kcal` |
| b | Down de 00003 = 1.2 | OK | base 1.2 contra 1.3 con Down: `sqlite_master` (type, name, tbl_name, sql) 278 = 278 objetos, idéntico. Datos: 55 de 55 tablas iguales; las únicas diferencias son `created_at`/`updated_at` porque son dos bases distintas. En una misma base con goose (`up-to 2` → volcado → `up` → `down-to 2` → volcado): `sqlite_master` (1295 líneas) y `.dump --data-only` (25 líneas) idénticos con `cmp` |
| c | goose v3.28.0 up / down-to 0 / up | OK | `up` → versión 3; `down-to 0` → 0 tablas; `up` → 55 tablas, `integrity_check` ok; `status` muestra las 3 migraciones aplicadas |
| d | DDL de 00001 y 00002 contra `git show HEAD` | OK | `git diff --quiet HEAD` y `git show HEAD:… \| cmp -`: los dos archivos son idénticos byte por byte, comentarios incluidos |
| e | `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` | OK | exit 0, **RESUMEN TODO OK** (26 líneas OK, 0 FALLA); KUESKI 389500 / 239528, 0 diferencias; DOWN13: 240 objetos y 130 filas iguales |
| f | `gen_dict.py` regenera idéntico | OK | `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev-s1-r1/dic.md` → `cmp` contra `docs/contratos/diccionario-datos.md` idéntico (182,480 bytes); «valores iniciales que no coinciden: []» |
| g | Valores de K44 contra §K | **Parcial (H-01)** | En la base 1.3, `app_setting` con alcance de la casa: `price.stale_after_days`=30, `price.zero_weight_after_days`=90, `session.max_age_days`=90, `auth.max_failed_attempts`=5, `auth.lockout_window_minutes`=15, `attachment.max_size_mb`=10, `comparator.timeout_seconds`=5, `health.goal_adjust_lose_kcal`=0, `health.goal_adjust_gain_kcal`=0 y `idempotency.retention_hours`=24: 10 de 10 iguales. `auth.password_min_length`=12 y `price.profeco_import_max_gb`=1 son de K6. **80 y 100: ninguna clave de `app_setting`**; solo existen en `notification.params_json` (`budget_threshold`). Ningún valor inventado: lo que no tiene decisión está listado «sin valor» (optimizer.* de K16, `date_window_days`, `send_time` y `advance_days`, que la respuesta de `C1-valores-k31` deja a S7, y cinco valores de K31) |
| h | Ninguna marca PROVISIONAL cita un bloqueo resuelto | OK | `prov.py` cruza 16 líneas con PROVISIONAL (backend/db, diccionario y reglas) contra los 44 archivos de bloqueos. Solo se cita `S1-parametros-avisos`, que está abierto. Los 8 bloqueos de la puerta 1b tienen «Respuesta de Axel» |
| i | RN-27 y RN-12 recalculados | OK | Decimal ROUND_HALF_UP sobre `model.json`: `periodos[0].efectivo.total` = 973.75 × 7 / 31 = **219.88** (diario 31.41); `periodos[4].efectivo.total` = 2,517.38 × 7 / 28 = **629.35** (diario 89.91); devengado 1,947.50 × 7 / 31 = 439.76. Con el diario redondeado saldrían 219.87, 629.37 y 439.74. Pesos: (90 − 30)/61 = **0.9836**; (90 − 60)/61 = **0.4918**; día 90 = 0. RN-17: 2,555.5625 → 2,556; con 500 kcal para bajar, 2,056 |
| j | Mínimos de seguridad fuera del esquema | OK | 00003 solo valida el tipo y el signo (entero JSON ≥ 0) y aplica `CHECK (≥ 0)` a los ajustes de kcal, que es la magnitud de K36. Prueba propia: `auth.max_failed_attempts`=0, `session.max_age_days`=0 y 100000 se aceptan; `-5` y `0.5` se rechazan. Los mínimos y máximos solo están en «Propuestas para Axel» del reporte |

## Criterios
| Criterio | Estado | Evidencia |
|---|---|---|
| S1-01 (alta) entidades | Cumple | s1_verifica S1-01: 31 de 31 |
| S1-02 (alta) migraciones en limpio | Cumple | (a) y (c) |
| S1-03 (alta) ejemplo numérico | Cumple | S1-03: 29 reglas con ejemplo. RN-27 con cuentas citadas de `model.json` (219.88, 629.35, 439.76) y RN-29 con ejemplo ilustrativo rotulado. Observación H-03 (baja) |
| S1-04 (alta) sin tablas huérfanas | Cumple | 55 de 55 con «Usado por»; `health_profile` lleva `@used_new` de 1.3 |
| S1-05 (alta) dinero exacto | Cumple | 0 columnas REAL/FLOAT/DOUBLE |
| S1-06 (media) fechas y zona | Cumple | s1_verifica S1-06 |
| S1-07 (alta) personas y acceso | Cumple | s1_verifica S1-07 |
| S1-08 (alta) reparto | Cumple | s1_verifica S1-08 |
| S1-09 (alta) periodos | Cumple | s1_verifica S1-09: 18 filas |
| S1-10 (alta) créditos | Cumple | s1_verifica S1-10 y K12 |
| S1-11 (media) comida y salud | Cumple | 12 de 12 puntos |
| S1-12 (alta) ida y vuelta | Cumple | 0 diferencias; KUESKI 389500 / 239528 |
| S1-13 (media) Ajustes | **Parcial (H-02)** | La fila 13 suma las columnas de K36, pero las filas 1, 5, 6 y 11 no nombran 6 parámetros del inventario |
| S1-14 (alta) sesiones | Cumple | Fila 1: 5 intentos, 15 min y 90 días; sin secretos |
| S1-15 (media) integridad | Cumple | s1_verifica S1-15 |
| S1-16 (alta) sin inventados | Cumple | (g) sin valores inventados y (h); los supuestos 1 a 7 no cambian cifras. La sección de Ajustes y la exclusión de constantes y límites de formato van como propuesta para la puerta 2 en RN-29 |
| S1-17 (media) idioma | Cumple | 658 identificadores en `[a-z0-9_]` |
| K36 ajuste kcal | Cumple | Columnas por persona (NULL = valor de la casa, ≥ 0); claves 0 y 0; RN-17 FIRME con la fórmula y HALF_UP al final |
| K37 fuentes | Cumple | RN-12: compiten `ticket_manual`, `ticket_ocr` y `manual`; `site_assisted` y PROFECO solo como referencia |
| K39 promedio semanal | Cumple | (i) |
| K40 frescura | Cumple | (i); se siembra `price.zero_weight_after_days` = 90; la forma general con Dv y D0 da K40 con los valores iniciales |
| K43 valores de C1 | Cumple | 10 MB y 5 s sembrados; `send_time` y `advance_days` quedan sin valor porque S7 los propone (bloqueo `C1-valores-k31`, línea 6) |
| K44 parámetros configurables | **No cumple (H-01)** | Umbrales 80/100 % fuera de `app_setting` |

## Hallazgos
### H-01 (alta) · K44: los umbrales 80 % y 100 % no son claves de `app_setting`
- **Ruta:** `backend/db/migrations/00003_puerta1b.sql:49` (`@param_sin_valor umbrales…`), `docs/contratos/reglas-negocio.md` RN-29 (filas `thresholds_pct` y `days_before`), `docs/estado/bloqueos/S1-parametros-avisos.md`, `backend/db/tools/s1_verifica.py:565` (la prueba K44 exige 0 claves `alerts.*`).
- **Qué pasa:** K44 nombra entre los parámetros los «umbrales de alertas 80 % y 100 %» y fija: «Esquema 1.3: cada parámetro es una clave de `app_setting` con valor inicial sembrado y comentario que cita su decisión; C1 v1.2 las expone por las operaciones de settings». La respuesta de Axel en `C1-valores-k31` repite el destino: «esquema 1.3 (claves en `app_setting`)». El esquema 1.3 no siembra esas claves; los valores siguen en `notification.params_json` y el agente abrió un bloqueo sobre algo que §K ya decidió para los umbrales. La comprobación (g) falla para 80 y 100. Lo mismo vale para `days_before` [1, 0] de `payment_due` («día anterior y el mismo día», §F), que es un valor numérico operativo decidido. La parte del bloqueo sobre `send_time` y `advance_days` sí es un choque real entre K43 y K44, y esos no tienen valor.
- **Corrección sugerida:** aplicar K44 en la 00003 (opción B del bloqueo). (1) Sembrar `alerts.budget_warning_pct` = 80 y `alerts.budget_limit_pct` = 100, enteros, sección Notificaciones, con descripción que cite §F y K44, y la anticipación del pago próximo (`days_before` [1, 0], de §F) como clave de `app_setting`. (2) Agregarlas a `@param`, a los disparadores (o a un disparador de lista de enteros) y al Down. (3) Para no dejar dos fuentes, quitar `thresholds_pct` y `days_before` de `notification.params_json` en el Up y reponerlos byte por byte en el Down. (4) Ajustar la prueba K44 de `s1_verifica.py`, RN-29 y el diccionario. (5) Pedir como solicitudes los cambios de C1 v1.2 (`NotificationUpdate` y settings) y de S7. (6) Reducir el bloqueo a `send_time` y `advance_days`. Otra salida: que el orquestador registre la opción A como excepción explícita a K44 aprobada por Axel; mientras no exista esa excepción, el esquema no es coherente con §K.

### H-02 (media) · S1-13: la tabla de las 14 secciones de Ajustes no recoge 6 parámetros del inventario
- **Ruta:** `docs/contratos/diccionario-datos.md:2307` y siguientes (generado por `backend/db/tools/gen_dict.py:220`, que solo parchea la fila 13).
- **Qué pasa:** el inventario K44 (diccionario y RN-29) asigna cada parámetro a una sección. La tabla del criterio S1-13 no los nombra: fila 1 «Cuenta y roles» sin `auth.password_min_length`; fila 5 «Créditos y condiciones» sin `comparator.timeout_seconds`; fila 6 «Tiendas y precios» sin `price.zero_weight_after_days` ni `price.profeco_import_max_gb`; fila 11 «Datos y respaldos» sin `attachment.max_size_mb` ni `idempotency.retention_hours`. Conteo propio: 0 apariciones en la tabla para esos 6, contra 1 para `stale_after_days`, `max_failed_attempts` y `session.max_age_days`. El diccionario queda inconsistente consigo mismo y U4, que lee la tabla, no los vería. `s1_verifica` S1-13 pasa porque solo comprueba que haya destinos.
- **Corrección sugerida:** que `gen_dict.py` agregue a cada fila de la tabla las claves `@param` de su sección, con valor inicial, de forma idempotente; regenerar. En `s1_verifica` S1-13, comprobar que cada clave `@param` aparece en la fila de su sección.

### H-03 (baja) · RN-27 llama «periodo de octubre» a la ventana del oráculo
- **Ruta:** `docs/contratos/reglas-negocio.md:278`.
- **Qué pasa:** el ejemplo dice «periodo de octubre de 2026 (10-oct-2026 a 9-nov-2026, 31 días)» y «periodo de febrero de 2027 (10-feb-2027 a 9-mar-2027)». La decisión B y S1-09 fijan que el periodo es el mes calendario. Las cifras no cambian: los días coinciden (31 y 28) y la equivalencia está probada en S1-09 (b).
- **Corrección sugerida:** «mes calendario de octubre de 2026 (31 días); su efectivo es `periodos[0].efectivo.total` por la equivalencia de S1-09», y lo mismo para febrero.

### H-04 (baja) · Encabezado engañoso en la 00003
- **Ruta:** `backend/db/migrations/00003_puerta1b.sql:42`.
- **Qué pasa:** «Parámetros que existen sin valor sembrado» agrupa también los umbrales y `days_before`, que sí tienen valor en `notification.params_json`. El diccionario lo dice bien: «Existen sin valor sembrado o viven en otra tabla».
- **Corrección sugerida:** usar el mismo texto que el diccionario (o dejar de agrupar los umbrales ahí si se aplica H-01).

### H-05 (baja) · Granularidad entera de unidades grandes
- **Ruta:** `backend/db/migrations/00003_puerta1b.sql:36` y el disparador de la línea 78.
- **Qué pasa:** `price.profeco_import_max_gb` es un entero en GB. Abajo de 1 GB solo se puede poner 0, así que no se puede bajar a 0.5 GB. Ni K6 ni K44 fijan que el valor sea entero.
- **Corrección sugerida:** guardar el límite en MB (1024) o llevar la unidad a las propuestas de la puerta 2. Lo mismo, de menor peso, para `attachment.max_size_mb`.

### H-06 (baja) · Dos listas de optimizer.* en el diccionario
- **Ruta:** `docs/contratos/diccionario-datos.md:2280` y `:2292`.
- **Qué pasa:** «Existen sin valor sembrado…», generada de la 00003, y «Claves de app_setting previstas sin valor», generada de la 00001, repiten las mismas cuatro claves `optimizer.*` desde dos fuentes. Hoy coinciden, pero pueden divergir.
- **Corrección sugerida:** que una lista remita a la otra, o generar las dos desde una sola fuente de anotaciones.

### H-07 (baja) · La detección de marcas con bloqueo resuelto en s1_verifica es parcial
- **Ruta:** `backend/db/tools/s1_verifica.py:562`.
- **Qué pasa:** la expresión solo reconoce los prefijos S1, C1, U1, U2, U2a, S6 y criterios. Una marca que cite un bloqueo resuelto S2-, S5-, P3- u otro pasaría. Mi cruce contra los 44 nombres de archivo encontró 0 casos hoy.
- **Corrección sugerida:** cruzar contra la lista de nombres de `docs/estado/bloqueos/*.md`.

## Estado de los anteriores
Ronda 1 de este flujo: no hay hallazgos previos de S1 en P1c. Los 5 bloqueos de la puerta 1b que tocan S1 (`S1-ajuste-objetivo-kcal`, `S1-fuentes-precio`, `criterios-promedios-redondeo`, `criterios-frescura-extremos` y `C1-valores-k31`) están aplicados y sus marcas se quitaron de migraciones, diccionario y reglas (h). Fuera del alcance de esta revisión y anotados por el agente como pendientes: `INTEGRACION.md` §7.2 (`interest_includes_vat`, K18) y §7.7 (renombrar `credit.rate_base`).
