# Corrección S1 · ronda 2 (flujo P1c, 11-oct-2026) · agente correccion-S1-r2

Corrige el único hallazgo de alta o media de `docs/estado/p1c/REVISION-S1-r2.md`: **M-01 (media)**, el inventario de K44 no mencionaba los valores de programación de la decisión F que siguen en `notification.params_json`. Los hallazgos B-01 a B-04 (baja) no se pidieron en esta ronda y no se tocaron. Ningún valor, columna, clave ni DDL cambia.

## Hecho
- Se eligió la primera corrección que sugiere el revisor, porque no decide nada que Axel no haya decidido. `weekday` = 7 (`weekly_summary` y `meal_plan_missing`) y `day_of_month` = 1 (`monthly_report`) tienen el valor de la decisión F (FIRME) y siguen en `notification.params_json`, donde los sembró la 00001. Dónde deben vivir es la misma pregunta que ya está abierta para `send_time`: `params_json` (K43) o `app_setting` (K44). Por eso entran en el bloqueo `S1-parametros-avisos` y quedan PROVISIONAL solo en su lugar. La otra opción del revisor (sembrarlos en `app_setting`) habría decidido el bloqueo sin Axel.
- `backend/db/migrations/00003_puerta1b.sql` (solo comentarios y anotaciones):
  - La cabecera dice que los días de programación de F siguen en `params_json`, con su lugar PROVISIONAL (corrección S1-r2, M-01).
  - El rótulo del bloque pasa a «sin valor sembrado o que viven en otra tabla».
  - Dos `@param_sin_valor` nuevas, `notification.params_json weekday (weekly_summary y meal_plan_missing)` y `notification.params_json day_of_month (monthly_report)`. Cada una lleva su valor de la decisión F, aclara que no se siembra en `app_setting` y es PROVISIONAL citando `S1-parametros-avisos`.
  - El `@col_note` de `notification.params_json` dice que `weekday` y `day_of_month` siguen ahí.
- `docs/estado/bloqueos/S1-parametros-avisos.md` (protocolo §2; el revisor lo pide en M-01):
  - El contexto explica la ampliación.
  - «Qué falta» pregunta por los cuatro valores: `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1).
  - Las opciones A, B y C dicen qué pasa con `weekday` y `day_of_month` en cada caso. Con B, una migración siembra tres claves citando F, las quita de `params_json` y su Down las repone. Con C, la programación por aviso se queda en `params_json` y la anticipación pasa a `app_setting`.
  - También se actualizaron Impacto y «Cómo seguí». El bloqueo sigue abierto: no tiene «Respuesta de Axel».
- `docs/contratos/reglas-negocio.md`:
  - El encabezado de RN-29 y su fila del índice nombran `weekday` y `day_of_month` (valores 7 y 1 de F) dentro de la marca PROVISIONAL del lugar.
  - Bullet nuevo «Viven en otra tabla», con el valor, el aviso, la sección de Ajustes según la tabla «Ajustes: 14 secciones» del diccionario y el bloqueo.
  - La línea «Esquema.» de RN-29 dice que `params_json` conserva `weekday` y `day_of_month`.
  - En «Cambios en 1.3» el bloqueo ya no dice «solo cubre `send_time` y `advance_days`».
- `backend/db/tools/gen_dict.py`: la nota «Cambios en 1.3» del diccionario añade «Corrección de la revisión S1-r2 de P1c … M-01».
- `docs/contratos/diccionario-datos.md`: **regenerado** con `gen_dict.py`, nunca editado a mano. Contra el diccionario anterior cambian 4 líneas: la nota de cambios, la descripción de `notification.params_json` y dos filas nuevas en «Existen sin valor sembrado o viven en otra tabla».
- `backend/db/tools/s1_verifica.py`: comprobación nueva **PROG13**. Revisa que en 1.3 `params_json` siga siendo `{"weekday": 7}`, `{"weekday": 7}` y `{"day_of_month": 1}`, y que haya exactamente una `@param_sin_valor` por valor, con el valor de F, PROVISIONAL y el bloqueo. También revisa que estén en la sección K44 del diccionario y en el encabezado y la lista de RN-29, que la pregunta del bloqueo nombre los cuatro valores y que ninguna clave de `app_setting` tenga `weekday` ni `day_of_month`.

## Archivos
- `backend/db/migrations/00003_puerta1b.sql` (modificado: comentarios y anotaciones; Up y Down sin cambios)
- `backend/db/tools/s1_verifica.py` (modificado: PROG13)
- `backend/db/tools/gen_dict.py` (modificado: nota de la corrección)
- `docs/contratos/diccionario-datos.md` (regenerado)
- `docs/contratos/reglas-negocio.md` (modificado: RN-29, índice, «Cambios en 1.3»)
- `docs/estado/bloqueos/S1-parametros-avisos.md` (ampliado)
- `docs/estado/p1c/correccion-S1-r2.md` (este reporte)
- 00001 y 00002: sin tocar, ni siquiera los comentarios.

## Verificación
- `GOOSE=…/gobin/goose python3 -I backend/db/tools/s1_verifica.py`: exit 0, 29 líneas, **RESUMEN TODO OK**. Pasan S1-01 a S1-17, KUESKI (pagos 389500, parte de Axel 239528, 0 diferencias), K12, NEG (46/46), NEG12 (26/26) y K44 (10 claves sembradas, 15 `@param` iguales a lo sembrado, sin marcas que citen bloqueos resueltos, cifras de RN-17, RN-12 y RN-27 presentes). También pasan AVISOS13, **PROG13** (nueva: `weekday` 1 y `day_of_month` 1 anotaciones; diccionario True/True; RN-29 True; bloqueo True; 0 claves de app_setting con esos nombres), DDL, DOWN13, DOWN y GOOSE (códigos [0, 0, 0], tablas [55, 0, 55], versión 3).
- S1-16: 19 líneas con PROVISIONAL (migración 6, diccionario 7, reglas 6). Todas citan solo el bloqueo abierto `S1-parametros-avisos` y 0 quedan sin bloqueo.
- Prueba negativa de PROG13: sobre una copia de la 00003 sin la corrección, `{weekday: 0, day_of_month: 0}` → FALLA (esperado); sobre la nueva, `{1, 1}` → pasa.
- `python3 -I backend/db/tools/gen_dict.py`: versión 1.3 · 15 parámetros · 8 sin valor (antes 6) · 0 valores iniciales distintos de lo sembrado · 55 tablas, 152 índices, 75 disparadores · 0 columnas sin descripción · 0 parámetros sin fila en «Ajustes». Regenerarlo otra vez a `datos/tmp` y compararlo con `cmp` da idéntico (idempotente).
- Diferencia exacta del diccionario: se regeneró el diccionario anterior desde una copia de la 00003 sin la corrección y se comparó con `diff`. Salen 4 líneas, todas las de M-01.
- `git diff --quiet HEAD -- 00001_esquema.sql 00002_puerta1.sql`: exit 0.
- Los mínimos propuestos (3 a 10, 5 a 60, 1 a 365) aparecen 0 veces en la 00003, las reglas y el diccionario.

## Supuestos
- Ninguno que cambie cifras, esquema, endpoints ni tokens. La sección de Ajustes citada para `weekday` y `day_of_month` («Reportes y frecuencias» y «Notificaciones») sale de la tabla «Ajustes: 14 secciones» que ya tenía el diccionario; no es nueva.

## Bloqueos
- `docs/estado/bloqueos/S1-parametros-avisos.md` (abierto, ampliado): ahora cubre dónde viven `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1).

## Pendientes
- Hallazgos de baja de `REVISION-S1-r2.md` que esta ronda no pidió: B-01 (el Down no devuelve `params_json` byte por byte si la casa cambió otras claves), B-02 (la descripción de `notification.params_json` empieza con datos que ya no aplican), B-03 (la línea de versión de las reglas) y B-04 (`optimizer.weights` fuera del inventario).
- Cuando Axel responda `S1-parametros-avisos`: con A o C, quitar las marcas. Con B, escribir una migración que siembre tres claves (7, 7 y 1, citando F) y las quite de `params_json` con un Down que las reponga, igual que `days_before` en la 00003. Luego ajustar PROG13.

## Solicitudes fuera de alcance
- `docs/estado/bloqueos/S1-parametros-avisos.md`: se amplió la pregunta, que es la primera corrección que pide M-01. Lo permite el protocolo §2. El integrador debe saber que este archivo cambió fuera de la lista de alcance.

## Propuestas para Axel
- Para `S1-parametros-avisos`, S1 recomienda la **opción B** (todo en `app_setting`). Así K44 tiene una sola fuente por valor, como ya quedaron los umbrales y los días antes del pago próximo (corrección H-01), y Ajustes y C1 v1.2 los manejan con las mismas operaciones de settings. `weekday` 7, 7 y `day_of_month` 1 se moverían con su valor de la decisión F y un Down que los reponga. Si se quiere una hora distinta por aviso, la opción C es la alternativa.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: ~27.
- Leído: ~95 KB (protocolo 5.6; §K ~11; revisión ~7; bloqueo 2.9; 00003 18.5; RN-29 y cabecera de reglas ~10; partes de gen_dict ~9 y de s1_verifica ~14; grep del diccionario ~3; salidas de verificación ~8).
- Escrito: ~15 KB de contenido propio (parches a la 00003 ~2.5, a las reglas ~2.3, a s1_verifica ~2.6 y a gen_dict ~0.4; bloqueo 4.3; este reporte ~6). Aparte, el diccionario lo reescribió el generador.
