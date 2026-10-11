# correccion-S1-r1 · Reporte final (flujo P1c, 11-oct-2026)
Corrige los hallazgos H-01 (alta) y H-02 (media) de `docs/estado/p1c/REVISION-S1-r1.md`. No reabre decisiones de §A–§K ni toca los hallazgos bajos (H-03 a H-07); H-04 queda resuelto de paso porque la línea de los umbrales salió del grupo «sin valor sembrado».

## Hecho
1. **H-01 · K44 aplicado a los avisos (opción B del bloqueo).** La 00003 siembra `alerts.budget_warning_pct` = 80 y `alerts.budget_limit_pct` = 100 (enteros, % del presupuesto) y `alerts.payment_due_days_before` = [1, 0] (lista de enteros, días antes del pago). Las tres van en la sección Notificaciones y su descripción cita la decisión F y K44.
   - Están en `@param` y en los disparadores. Los enteros entran en `trg_app_setting_param_int_ins` y `_upd`. Para la lista hay dos disparadores nuevos, `trg_app_setting_param_int_list_ins` y `_upd`: aceptan una lista JSON de enteros ≥ 0, solo validan el tipo y el signo, y no fijan mínimos ni cuántos elementos. El Down borra las tres claves en cualquier alcance y los dos disparadores.
   - Cada clave toma el valor que tenga el aviso al migrar. En una base nueva son 80, 100 y [1, 0]. Si el valor falta o no es válido, se usa el de la decisión F.
   - Una sola fuente: el Up quita `days_before` de `payment_due` (queda `{}`) y `thresholds_pct` de `budget_threshold` (queda `{"scopes": ["category", "total"]}`).
   - El Down repone los dos valores al principio de `params_json` con el valor vigente de `app_setting`. Con los valores iniciales el texto es el de 1.2 byte por byte. `updated_at` se conserva con un UPDATE en dos pasos, porque `trg_notification_updated_at` lo renueva si el UPDATE no lo trae.
   - Una anotación nueva, `-- @col_note`, documenta en el diccionario 1.3 que `notification.params_json` ya no lleva esos valores. No se cambia ningún comentario de 00001 ni de 00002: así `sqlite_master` de 1.2 no cambia.
2. **H-02 · S1-13.** `gen_dict.py` agrega a cada fila de «Ajustes: 14 secciones» las claves `@param` de su sección con su valor inicial. El sufijo generado se quita antes de volver a ponerlo, así que regenerar da el mismo archivo. Si una clave no encuentra su fila, el script sale con error. También corrige la fila 9: `params_json` queda con `scopes` y `weekday`, y los umbrales y los días antes van a `app_setting`. Añade la línea de datos de referencia de `notification` en la 00003 y actualiza «Cambios en 1.3» y la introducción del inventario. El diccionario está regenerado: 15 de 15 valores iniciales coinciden y las 6 filas llevan sus claves.
3. **`s1_verifica.py`:**
   - K44 exige las 15 anotaciones y las 3 claves `alerts.*` sembradas, y que `params_json` ya no lleve esos valores y conserve el `updated_at` de 1.2. Suma 6 negativas y 2 positivas.
   - S1-13 comprueba que cada clave `@param` aparece en la fila de su sección.
   - DOWN13 exige que `notification` vuelva byte por byte al estado de 1.2.
   - Prueba nueva AVISOS13: hace el traspaso en los dos sentidos con valores que no son los iniciales.
4. **`reglas-negocio.md`:**
   - RN-29: tabla con las tres claves `alerts.*`; definición con la lista de enteros, los disparadores nuevos y la regla de una sola fuente; línea de esquema; la marca PROVISIONAL queda solo para dónde viven `send_time` y `advance_days`.
   - RN-22 cita los umbrales configurables.
   - Se ajustaron «Cambios en 1.3» y el índice.
5. **Bloqueo `docs/estado/bloqueos/S1-parametros-avisos.md` reducido** a `send_time` y `advance_days`, que es el choque real entre K43 y K44 y no tienen valor. Las opciones se reescribieron, y la B es ahora la coherente con H-01.

## Archivos
- `backend/db/migrations/00003_puerta1b.sql` (modificado)
- `backend/db/tools/gen_dict.py` (modificado)
- `backend/db/tools/s1_verifica.py` (modificado)
- `docs/contratos/diccionario-datos.md` (regenerado con `gen_dict.py`)
- `docs/contratos/reglas-negocio.md` (modificado: RN-22, RN-29, índice y cambios 1.3)
- `docs/estado/bloqueos/S1-parametros-avisos.md` (reducido)
- `docs/estado/p1c/correccion-S1-r1.md` (este reporte)

## Verificación
- `GOOSE=…/gobin/goose python3 -I backend/db/tools/s1_verifica.py` → exit 0, **RESUMEN TODO OK** (27 pruebas). Puntos clave:
  - S1-02: 55 tablas STRICT, 152 índices, 75 disparadores, `integrity_check` ok, 0 FK rotas.
  - S1-13: 14 de 14 filas; claves `@param` en la fila de su sección: 15 de 15.
  - KUESKI 389500 / 239528, 0 diferencias. S1-12: 164 valores, 0 diferencias.
  - K44: 15 anotaciones iguales a lo sembrado; negativas 15 de 15 rechazadas; positivas 7 de 7 aceptadas; cifras RN-17 2556 y 2056, RN-12 0.9836, 0.4918 y 0, RN-27 219.88 y 629.35; `params_json` sin los valores y con el `updated_at` de 1.2 = True.
  - AVISOS13: un 1.2 con [75, 95] y [3, 1, 0] da `app_setting` 75, 95 y `[3, 1, 0]`. Tras poner 70 y [2,0] en 1.3, el Down da `{"thresholds_pct": [70, 95], …}` y `{"days_before": [2, 0]}`, con `updated_at` conservado.
  - DDL de 00001 y 00002 igual a git HEAD.
  - DOWN13: 240 objetos iguales a 1.2, `app_setting` igual (10 claves), 54 tablas y 130 filas iguales, `notification` incluida.
  - DOWN: Down de la 00002 igual a 1.1.
  - GOOSE: up / down-to 0 / up con códigos [0, 0, 0] y versión 3.
- goose en una misma base, `up-to 2` → volcado → `up` → `down-to 2` → volcado (`iterdump` sin `goose_db_version`): `cmp` idéntico (1279 líneas). En 1.3, `budget_threshold` = `{"scopes": ["category", "total"]}`, `payment_due` = `{}` y las claves `alerts.*` = 80, 100 y `[1, 0]`.
- `python3 -I backend/db/tools/gen_dict.py` dos veces, la segunda a otro archivo: `cmp` idéntico. Parámetros 15, sin valor 6, valores que no coinciden [], parámetros sin fila [], columnas sin descripción [].
- Las marcas PROVISIONAL se cruzaron contra los nombres de todos los archivos de `docs/estado/bloqueos/`. En 00003, el diccionario y las reglas solo se cita `S1-parametros-avisos` (abierto). Bloqueos resueltos citados: 0.
- `git diff --stat` de 00001 y 00002: sin cambios.

## Supuestos
1. El nombre de las claves (`alerts.budget_warning_pct`, `alerts.budget_limit_pct` y `alerts.payment_due_days_before`) es el de la opción B del bloqueo y de la corrección sugerida. La sección Notificaciones es la PROPUESTA del orquestador en K44.
2. El Up trae a `app_setting` el valor que el aviso tenga en 1.2 y el Down devuelve a `params_json` el valor vigente en 1.3. Así no se pierde lo que la casa haya cambiado. Con los valores iniciales, todo vuelve byte por byte.
3. Si `thresholds_pct` tiene más de dos elementos, se toman el primero y el segundo. K44 y la decisión F definen exactamente dos umbrales.
4. `updated_at` de los avisos se conserva en la migración porque la casa no cambió nada. Ninguno de estos supuestos cambia cifras.

## Bloqueos
- `docs/estado/bloqueos/S1-parametros-avisos.md`, abierto y reducido a dónde viven `send_time` y `advance_days`.

## Pendientes
- Ninguno de H-01 ni de H-02. Los hallazgos bajos H-03 (texto de RN-27), H-05 (granularidad en GB y MB), H-06 (dos listas de `optimizer.*`) y H-07 (detección de prefijos en `s1_verifica`) no se tocaron porque están fuera del encargo.
- Observación: `params_json` todavía guarda valores de programación: `weekday` 7 en `weekly_summary` y `meal_plan_missing`, y `day_of_month` 1 en `monthly_report`. La revisión no los cita, así que no se tocaron. Ver Propuestas.

## Solicitudes fuera de alcance
1. Fuente de C1 (`backend/internal/api/c1gen/`) y regeneración de `docs/contratos/openapi.yaml`, `C1-datos.md` y los extractos, para C1 v1.2:
   - `NotificationUpdate` deja de documentar y aceptar `thresholds_pct` (en `budget_threshold`) y `days_before` (en `payment_due`).
   - Las operaciones de settings exponen `alerts.budget_warning_pct` y `alerts.budget_limit_pct` (entero ≥ 0) y `alerts.payment_due_days_before` (lista de enteros ≥ 0), con valor inicial y «restablecer».
   - Motivo: K44 y H-01. Ya no hay dos fuentes.
2. S7 (WF-1): los avisos de presupuesto y de pago próximo leen sus parámetros de `app_setting` `alerts.*`, no de `params_json`. Motivo: H-01.
3. U1 1.2 y U4: Ajustes «Notificaciones» muestra las tres claves con su valor inicial y «restablecer». Motivo: K44.
4. Orquestador: llevar a la puerta 2 la pregunta reducida del bloqueo `S1-parametros-avisos`.

## Propuestas para Axel
1. **Mínimos de los umbrales:** el aviso debe ser menor que el tope; si no, los dos llegan juntos. Se propone el aviso entre 1 y 99 % y el tope entre 100 y 200 %. Arriba de 200 % el aviso casi nunca llegaría. Se aprueba en la puerta 2 (K31, K44). Hoy el esquema no pone tope.
2. **Días antes del pago:** de 1 a 3 elementos, sin repetidos, cada uno de 0 a 14 días. 14 días son un periodo quincenal: más anticipación que eso se cruza con el pago anterior. La lista vacía no se permitiría: para apagar el aviso está `is_enabled`. Se aprueba en la puerta 2. Hoy el esquema solo exige enteros ≥ 0.
3. **Programación de reportes y avisos:** falta decidir si `weekday` (7) y `day_of_month` (1) en `params_json` cuentan como «parámetro operativo» de K44 y pasan a `app_setting`, como los umbrales, o siguen como programación de cada aviso. Conviene decidirlo junto con `send_time` y `advance_days`.

## Consumo aproximado
- 37 acciones de herramienta.
- Leído: unos 175 KB.
- Escrito a mano: unos 48 KB (scripts de parche y bloque SQL, bloqueo, progreso y reporte). Además, los scripts de parche reescriben archivos en disco y `gen_dict.py` regeneró el diccionario (186 KB).
- Temporales en `datos/tmp/cs1r1`, ya borrados.
