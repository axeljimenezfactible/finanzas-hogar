# Revisión S1 · esquema 1.3 · ronda 3 de 3 (flujo P1c, 11-oct-2026) · agente revisor-S1-r3

Revisor independiente. No escribí nada de lo revisado y no corregí nada. Evalúo solo contra `docs/criterios/S1.md` (1.1) y `docs/plan/decisiones.md` §K (K36, K37, K39, K40, K43, K44). Reejecuté cada verificación con mis propios scripts en `datos/tmp/rev3/` (borrados al terminar). No di nada por bueno solo porque lo dijera un reporte.

## Veredicto
**APROBADO.** No quedan hallazgos de severidad alta ni media.
- M-01 de la ronda 2 está **cerrado**. `weekday` (7) y `day_of_month` (1) de la decisión F ya aparecen en estos lugares:
  - en el inventario de la 00003, como `@param_sin_valor`, en el rótulo «viven en otra tabla»;
  - en el diccionario generado: en la sección K44 y en la descripción de `notification.params_json`;
  - en RN-29: en el encabezado, en el índice y en el apartado «Viven en otra tabla»;
  - en el bloqueo abierto `S1-parametros-avisos`, que ahora pregunta por los cuatro valores y tiene el formato del protocolo §2.
- La marca PROVISIONAL cubre solo su lugar; su valor es FIRME (F).
- Quedan 6 hallazgos de severidad baja: 4 vienen de la ronda 2 y siguen abiertos, y 2 son nuevos.

## Verificaciones reejecutadas (a–j)
| # | Qué | Estado | Evidencia (mis comandos) |
|---|---|---|---|
| a | Aplicar 00001 + 00002 + 00003 en `datos/tmp/rev3/rev.db` | Cumple | `goose up`: rc 0, versión 3. `integrity_check` = ok y `foreign_key_check` = vacío. `pragma_table_list`: 56 tablas, y la única sin STRICT es `goose_db_version`, que es de goose. Las 55 del esquema son STRICT |
| b | El Down de la 00003 regresa a la 1.2 | Cumple | Base 1.2 (`up-to 2`), copia, `up` y `down` sobre la copia. `sqlite_master` (tipo, nombre, tabla y sql) es idéntico: 279 objetos. Los datos de las 55 tablas son iguales; solo cambia `sqlite_sequence` de `goose_db_version` (3 contra 4), que es contabilidad de goose y no de la migración. El `updated_at` de `notification` se conserva en el Up y en el Down. Con valores cambiados (umbrales [70, 95] y `days_before` [2, 1, 0] en 1.2), el Up siembra 70, 95 y `[2, 1, 0]`. El Down repone el JSON con el mismo significado, aunque no byte por byte (B-01). Con 75 en 1.3, el Down repone `[75, 100]` |
| c | goose up / down-to 0 / up | Cumple | Códigos [0, 0, 0]. Quedan 0 tablas tras `down-to 0` y 55 tras `up`, en versión 3. El esquema queda idéntico al de un `up` directo. Los datos también son iguales, salvo `created_at` y `updated_at`, que dependen de la hora en que se corre |
| d | DDL de 00001 y 00002 sin cambios | Cumple | Comparé con `git show HEAD:` y los dos archivos son idénticos byte por byte. Sin comentarios, el DDL también es igual |
| e | `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` | Cumple | exit 0, 29 líneas, **RESUMEN TODO OK**. KUESKI: pagos 389500 y parte de Axel 239528, con 0 diferencias. También pasan DOWN13, PROG13, GOOSE y DDL |
| f | `gen_dict.py` regenera el diccionario idéntico | Cumple | `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev3/dic.md` y `cmp` contra `docs/contratos/diccionario-datos.md`: idénticos. Resultado: 15 parámetros, 8 sin valor, 0 desajustes y 0 columnas sin descripción |
| g | Valores de K44 contra §K | Cumple | Leídos de la base migrada: alertas 80 y 100; precio viejo 30; peso cero 90; sesión 90; 5 intentos; ventana de 15; adjunto 10 (MB); comparador 5 (s); kcal 0 y 0; idempotencia 24 (h). Además: contraseña 12 y PROFECO 1 (K6), y días antes [1, 0] (F). Todos son de §A–§K. `send_time` y `advance_days` no tienen valor en K43, el bloqueo C1-valores-k31 ni la decisión F («noche»): están listados «sin valor», igual que optimizer.* (K16) y los cinco valores de K31. No hay ningún valor inventado |
| h | Ninguna PROVISIONAL cita un bloqueo resuelto | Cumple | Script propio sobre la migración, el diccionario y las reglas: 19 líneas con la marca. Las que citan un bloqueo citan solo `S1-parametros-avisos`, que no tiene «Respuesta de Axel». Las 3 sin cita son definiciones de la marca o la línea anterior de un bloque que sí cita. Los bloqueos resueltos (`S1-ajuste-objetivo-kcal`, `S1-fuentes-precio`, `criterios-promedios-redondeo` y `criterios-frescura-extremos`) solo aparecen en la nota «Cambios en 1.3» como quitados. 00001 y 00002 no tienen ninguna marca |
| i | Cifras recalculadas | Cumple | Con Decimal y HALF_UP desde `oraculo/model.json`: P1 efectivo 973.75 entre 31 días da 31.41 de diario y **219.88** de semanal (con el diario redondeado saldría 219.87). P5 efectivo 2,517.38 entre 28 días da 89.91 y **629.35** (629.345 sin redondear; la otra lectura da 629.37). Devengado de despensa de P1: 62.82 y 439.76. RN-12: (90 − 30)/61 = **0.9836** y (90 − 60)/61 = **0.4918**. Las fechas +29, +30, +60 y +90 desde el 10-oct-2026 caen en 8-nov, 9-nov, 9-dic y 8-ene. El ejemplo de RN-29, 60/91, da 0.6593. RN-17 da 2,556 y 2,056 kcal |
| j | Los mínimos de seguridad no están en el esquema | Cumple | «3 a 10», «5 a 60», «1 a 365», «365 días», «no menos de» y «mínimo razonable» aparecen 0 veces en la 00003, las reglas y el diccionario. El Up, sin comentarios, solo compara `< 0` y `>= 0` (tipo y signo). Los rangos están solo en la sección «Propuestas para Axel» de `S1-esquema-1.3.md` |

Pruebas extra:
- **Disparadores de K44.** Rechazan -1, 5.0, `"5"` y `abc`, también en el alcance de una persona. Aceptan 3, 0 y '5' (texto JSON entero). La lista rechaza `[1, -1]`, `[1, "a"]` y `1`, y acepta `[]` y `[2, 1, 0]`.
- **`health_profile`.** Con una fila real, `goal_adjust_*` rechaza -1 (por CHECK) y `'x'` (por STRICT), y acepta 0, 500 y NULL.

## Criterios de `docs/criterios/S1.md`
| Criterio | Estado | Evidencia |
|---|---|---|
| S1-01 (alta) entidades | Cumple | s1_verifica S1-01: 31 de 31 entidades con su `nombre_codigo` |
| S1-02 (alta) migraciones | Cumple | (a) y (c): con goose v3.28.0, `ok` y ninguna fila |
| S1-03 (alta) ejemplos | Cumple | S1-03 OK, 29 reglas. Recalculé RN-12, RN-17, RN-27 y RN-29 en (i). El ejemplo de RN-27 cita rutas de model.json y da las cifras de K39. Ver la nota de redacción en B-06 |
| S1-04 (alta) sin huérfanas | Cumple | 55 de 55 tablas con «Usado por». Revisé a mano 5 tablas: `app_setting` (Ajustes, 32 menciones en U1), `health_profile` (pantalla Salud, 16), `weight_entry` (tendencia de peso, 51), `notification` (Ajustes «Notificaciones», 6) e `idempotency_key` (S6, POST de captura del flujo F-1). Las 5 tienen consumidor |
| S1-05 (alta) dinero | Cumple | Hay 0 columnas REAL, FLOAT o DOUBLE, y 37 columnas `_cents` de tipo INTEGER. La 00003 no agrega columnas de dinero |
| S1-06 (media) fechas | Cumple | S1-06 OK. En la 00003, la única «oct-2026» está en un comentario de la cabecera, no en datos |
| S1-07 (alta) personas | Cumple | S1-07 OK: rechaza la cuenta de Karen, una tercera cuenta y el ingreso de Ximena |
| S1-08 (alta) reparto | Cumple | S1-08 OK. Sin cambios desde 1.2 |
| S1-09 (alta) periodos | Cumple | S1-09 OK: 18 filas, límites de K8 y equivalencia con P1–P6 |
| S1-10 (alta) créditos | Cumple | S1-10 y K12 OK: 17.27 % y 25.8685 %, e intereses 103474 |
| S1-11 (media) comida y salud | Cumple | 12 de 12. El ajuste por objetivo vive en `health_profile` |
| S1-12 (alta) ida y vuelta | Cumple | 164 valores comparados, 0 diferencias; 389500 y 239528 |
| S1-13 (media) Ajustes | Cumple | 14 de 14 secciones, 46 referencias y 15 de 15 `@param` en la fila de su sección |
| S1-14 (alta) credenciales | Cumple | En la 00003, «password» solo aparece en la clave `auth.password_min_length`, y el patrón de token de Telegram da 0 |
| S1-15 (media) integridad | Cumple | S1-15 OK. Los disparadores nuevos validan el tipo de dato, no la integridad referencial |
| S1-16 (alta) sin inventados | Cumple | (g) y (h). Los valores de K6 y K7 no llevan marca, y M-01 está cerrado |
| S1-17 (media) idioma | Cumple | 658 identificadores en `[a-z0-9_]`. Las claves nuevas siguen el patrón `módulo.nombre` en inglés |

## Decisiones §K revisadas
| Decisión | Estado | Evidencia |
|---|---|---|
| K36 | Cumple | Columnas `goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` ≥ 0, donde NULL toma el valor de la casa. Claves `health.goal_adjust_*` = 0. RN-17 FIRME, con un solo HALF_UP al final |
| K37 | Cumple | RN-12: `manual` compite igual que el ticket, y `site_assisted` y `profeco` quedan solo como referencia |
| K39 | Cumple | RN-27 FIRME con un solo redondeo al final: 219.88 y 629.35 recalculados |
| K40 | Cumple | RN-12 (90 − días)/61, `price.zero_weight_after_days` = 90. Con 30 y 90, la generalización de RN-12 (supuesto 3 del agente) da exactamente K40 |
| K43 | Cumple | 10 MB y 5 s sembrados. `send_time` y `advance_days` existen sin valor, porque no hay uno decidido, y su lugar es PROVISIONAL por el bloqueo abierto |
| K44 | Cumple | 15 claves en el inventario con su valor decidido y su decisión citada; disparadores de tipo y signo sin mínimos; una sola fuente para umbrales y días antes. Los valores de la decisión F que siguen en `params_json` van marcados (M-01 cerrado). Ver B-04 |

## Hallazgos (todos de severidad baja)
- **B-01 (baja, sigue de r2) · El Down no devuelve el texto byte por byte si la casa cambió los valores.**
  - Dónde: `backend/db/migrations/00003_puerta1b.sql:157-162` (`json_remove` compacta el JSON) y `:170-182`.
  - Qué pasa: si `budget_threshold` tenía umbrales distintos de [80, 100], el Up deja `{"scopes":["category","total"]}` sin espacios. El Down repone `{"thresholds_pct": [70, 95], "scopes":["category","total"]}`, que tiene el mismo significado que el texto de 1.2 pero no los mismos bytes.
  - Corrección sugerida: normalizar el resto con `replace(json(...), ',', ', ')` y `': '` igual que ya se hace con `days_before`, o documentar que la igualdad es de significado cuando la casa cambió valores.
- **B-02 (baja, sigue de r2) · La descripción de `notification.params_json` empieza con datos que ya no aplican.**
  - Dónde: `docs/contratos/diccionario-datos.md:1985`, generado del comentario de la 00001 más el `@col_note`.
  - Qué pasa: la descripción todavía arranca con «(días antes 1 y 0, umbrales 80 y 100)».
  - Corrección sugerida: que `gen_dict.py` anteponga «Hasta 1.2:» a esa parte o que la sustituya con el `@col_note`, sin tocar la 00001.
- **B-03 (baja, sigue de r2) · La línea de versión de las reglas mezcla la 1.3 con la explicación de la 1.2.**
  - Dónde: `docs/contratos/reglas-negocio.md:3`.
  - Qué pasa: dice «Versión: 1.3 …» y luego «Se numera 1.2 para ir a la par…».
  - Corrección sugerida: «La 1.2 se numeró así para ir a la par del esquema…».
- **B-04 (baja, sigue de r2) · `optimizer.weights` no está en el inventario de K44.**
  - Dónde: `backend/db/migrations/00003_puerta1b.sql:36-61` y `docs/contratos/reglas-negocio.md` (tabla de RN-29). Está sembrada en la 00001 con «pesos iguales normalizados» (decisiones, comparador de escenarios) y `docs/contratos/diccionario-datos.md:2258` la lista.
  - Qué pasa: ya es una clave configurable con un valor decidido, así que cumple K44 en lo sustancial. Pero el inventario no la nombra: no está ni en `@param`, ni en `@param_sin_valor`, ni en «No son parámetros».
  - Corrección sugerida: agregarla al inventario como parámetro de tipo objeto, sin disparador de entero, o listarla con su motivo.
- **B-05 (baja, nuevo) · El reporte del agente S1 describe un bloqueo superado.**
  - Dónde: `docs/estado/p1c/S1-esquema-1.3.md:39`.
  - Qué pasa: su §Bloqueos todavía dice que los umbrales 80/100 % y los días antes viven en `notification.params_json` y que «seguí con la opción A (sin cambios)». La corrección S1-r1 los pasó a `alerts.*` de `app_setting` y la S1-r2 amplió el bloqueo a `weekday` y `day_of_month`. Quien lea solo ese reporte se lleva una idea vieja del estado.
  - Corrección sugerida: que el integrador lo consolide, o una nota al inicio del reporte que remita a `correccion-S1-r1.md` y `correccion-S1-r2.md`.
- **B-06 (baja, nuevo) · El ejemplo de RN-27 llama «periodo de octubre» a la ventana del oráculo.**
  - Dónde: `docs/contratos/reglas-negocio.md:278`.
  - Qué pasa: llama «periodo de octubre de 2026 (10-oct-2026 a 9-nov-2026, 31 días)» y «periodo de febrero de 2027 (10-feb-2027 a 9-mar-2027)» a las ventanas P1 y P5 del oráculo. Pero el periodo del sistema es el mes calendario (decisión B, S1-09). Las cifras son las mismas por la equivalencia que prueba S1-09, y 31 y 28 días coinciden, pero la redacción mezcla las dos nociones.
  - Corrección sugerida: «Mes de octubre de 2026 (31 días); su efectivo es igual al de `periodos[0]` (S1-09)» y lo mismo para febrero.

## Estado de los hallazgos de la ronda 2
| Hallazgo | Estado | Evidencia |
|---|---|---|
| M-01 (media), valores de programación de F fuera del inventario de K44 | **Cerrado** | `00003_puerta1b.sql:10-12, 58-59, 62`; RN-29 (`reglas-negocio.md:50, 286, 308, 313`); diccionario `:1985, 2294-2295` (regenerado, `cmp` idéntico); `docs/estado/bloqueos/S1-parametros-avisos.md` con «Qué falta» sobre los cuatro valores y las opciones A, B y C, sin «Respuesta de Axel». La base migrada mantiene `{"weekday": 7}` ×2 y `{"day_of_month": 1}`, y ninguna clave de `app_setting` los duplica. El agente no se decidió por la opción B, que habría resuelto el bloqueo sin Axel |
| B-01 a B-04 (baja) | Siguen abiertos | No se pidieron en la ronda de corrección. Ver arriba |
