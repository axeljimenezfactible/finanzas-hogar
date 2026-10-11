# REVISION-criterios-r4 · Criterios 1.2 (flujo P1d) · ronda 4 de 4 · 11-oct-2026

Revisor independiente `revisor-criterios-r4`. Evalúo solo contra `docs/criterios/INDICE.md` §1 y `docs/plan/decisiones.md` §K (K35 a K44). No escribí nada de lo revisado y no corregí nada. Reejecuté todas las comprobaciones con comandos y scripts propios en `datos/tmp/rev4c/` (borrados al terminar). No tomé como evidencia lo que dice `correccion-criterios-r4.md`.

## Veredicto
**NO APROBADO.** 1 hallazgo alto, 2 medios y 3 bajos.

La parte de los umbrales del hallazgo A está cerrada en todos los archivos. Las Verifica de C1-22 y U4-12 pasan contra la base migrada con goose y contra los extractos `ajustes-*.yaml`. B, C, E y F quedan cerrados, y D y G atendidos. Lo que no cierra:
- **Alto:** U4-12 limita la excepción a `send_time` y `advance_days`. El bloqueo abierto `S1-parametros-avisos` y la 00003 también incluyen `weekday` (7) y `day_of_month` (1), que viven en `notification.params_json`. Con el texto actual, un U4 fiel al esquema 1.3 no puede cumplir U4-12.
- **Medio:** S7-04 y S7-08 no piden que el día del resumen semanal y del aviso de plan sin armar se lea de `weekday`. S7-10 sí lo pide para `day_of_month`.
- **Medio:** la parte de C1-12 de B-01 (C1-r2) no se atendió. Su Verifica literal falla con el contrato correcto por `auth.password_min_length`.

## Comprobaciones pedidas
| # | Comprobación | Estado | Evidencia (reejecutada) |
|---|---|---|---|
| a | Estructura y conteo | Cumple | Script propio `chk.py`: 18 archivos con las secciones 1 a 5; 261 criterios (201 alta / 60 media / 0 baja); cada uno con 1 Verifica y 1 Evidencia; 0 fuera de §5 (con rangos expandidos); 0 repetidos. Las 18 filas de INDICE §2 y §10 (conteo y versión, U3 en 1.2) coinciden; totales §2 = 261 (201/60) y §10 = 261 (201/60/0). |
| b | Hallazgo A y Verifica de C1-22 y U4-12 | **Parcial (H-1, H-2)** | Umbrales y días antes del pago como `alerts.*` (00003, K44): C1-22 (`C1.md:14,121-122`), U4-12 (`U4.md:71-72`), S7-08 (`S7.md:54`), S7 §5 (`S7.md:112`), INDICE §6 (`:68`, `:96`) y §11 (`:174`, `:182`). `goose` v3.28.0 `up` 00001→00003 en base temporal (versión 3). Script `c122.py` (Verifica literal de C1-22): 15 claves `@param`, todas sembradas con su valor inicial y con cita `decisión [A-Z]` o `K[0-9]+` junto a la siembra; lectura con `value` e `initial_value` y escritura tipada en `ajustes-*.yaml`; los tres `alerts.*` en `getAlertSettings` y `updateAlertSettings`; 0 `x-provisional` sobre claves K44; 0 límites salvo `minimum: 0`; `send_time` y `advance_days` aparecen como `@param_sin_valor` con PROVISIONAL y bloqueo, en `updateNotification` con `x-provisional` y la cita, y sin clave en `app_setting`; `weekday` = 7 (`weekly_summary`, `meal_plan_missing`) y `day_of_month` = 1 en `params_json`; `thresholds_pct` y `days_before` en `params_json` = 0; `verify_c1.py --db` sobre esa base: 29 OK, «TODO OK». **C1-22: 0 fallas.** U4-12 (parte ejecutable; `docs/design/U4-ajustes.md` no existe hasta WF-1): las 12 claves de la lista mínima tienen valor y la sección que dice U4-12; 0 claves `@param` sin sección; **0 fallas**. Pero U4-12 no nombra `weekday` ni `day_of_month` (`grep -c "weekday\|day_of_month" U4.md` = 0): H-1. |
| c | B a G, I-02, B-01 de C1-r2, U3 con K38 | B y C cerrados; D a G atendidos; B-01 parcial (H-3) | P3-06 (`P3-claude-design.md:49-50`) y P3-04 (`:42`) corregidos. Regex de D probadas por mí: 6 advertencias dan 0; las credenciales de muestra dan 5 de 6 con P3-01/P3-12 (sin «clave», como antes) y 6 de 6 con P3-03. La propuesta de G existe: `S1-esquema-1.3.md:59` y `INTEGRACION.md:213`. I-02 en `S4.md:48`. C1-13: `openapi.yaml:7` `version: "1.2.0"`. C1-21: `grep -c operationId openapi.yaml` = **210** = total de `extractos/INDICE.md:31`. U3 cita K38 (`U3.md:6,19`). C1-12 no se tocó (H-3). |
| d | Cifras con Decimal | Cumple | 973.75 / 31 × 7 = **219.88** (con doble redondeo, 219.87); 2,517.38 / 28 × 7 = **629.35** (629.37); (90 − 30) / 61 = **0.9836**; (90 − 60) / 61 = **0.4918**; día 29 = 1; día 90 = 0; día 91 no compite. `canon.py`: 178 cifras, 0 diferencias; equivalencia del calendario: 12 cifras, 0 diferencias. |
| e | PROVISIONAL contra bloqueos resueltos | Cumple | Script propio `prov.py`: 44 bloqueos, 42 con «Respuesta de Axel»; abiertos `S1-parametros-avisos` y `criterios-nxx`. Las 11 líneas con PROVISIONAL que nombran un bloqueo resuelto son avisos de que la marca se quitó (S2:6, S2:51, S3:6, S4:9, S5:6, S5:36, U1:6, U2:6), Verifica que buscan que no haya marcas (C1:94, C1:126) o un falso positivo por la palabra INDICE (C1:137). **0 marcas activas que citen un bloqueo resuelto.** Script `umb.py`: ninguna oración pone los umbrales o los días antes del pago en `params_json` o con PROVISIONAL. |
| f | `verif_final.sh` y `make verify-criterios` | Cumple | `bash docs/criterios/tools/verif_final.sh` rc = 0: 261 (201/60/0); traza 0 faltantes; refs 261/261; citas 305/0; N-xx 0; canon 178/0; S1-09 SQL OK (54 tablas, K8) y equivalencia 12/0; matriz con 210 operaciones y 0 N-xxx sin operación; rutas `/tmp` 0; INDICE contra conteo: 0 diferencias. `GOOSE=… make verify-criterios` rc = 0: «las 11 comprobaciones dan lo esperado». |
| g | No se reabre §A–§K ni se inventan valores | Cumple, salvo H-1 | No hay valores inventados: el «1 MB = 1,048,576 bytes» de S6-10 se atribuye a S1 con su fuente y queda «a confirmar en la puerta 2»; el «resto de los días, semana en curso» de U3-01 sale de `U1-arquitectura.md:82`. H-1 obliga a decidir de hecho el lugar de `weekday` y `day_of_month` (opción B del bloqueo abierto) o a incumplir U4-12. |
| h | El diff solo toca lo listado | Cumple | `git diff HEAD --stat -- docs/criterios/`: 8 archivos (C1, INDICE, P3, S4, S6, S7, U3, U4); cada bloque cambiado corresponde a A, B a G, I-02, B-01 o U3/K38. `tools/` sin cambios. |

## Hallazgos

### H-1 (alta) · U4-12 deja fuera `weekday` y `day_of_month` de la excepción del bloqueo abierto
- **Ruta:** `docs/criterios/U4.md:6` y `:10` (listas de cambios) y `:71-72` (U4-12, enunciado y Verifica); `docs/criterios/INDICE.md:68`.
- **Qué pasa:** el bloqueo `S1-parametros-avisos` (ampliado por M-01 de S1-r2) y la 00003 (`@param_sin_valor`, líneas 58-59) tratan `weekday` = 7 (`weekly_summary`, `meal_plan_missing`) y `day_of_month` = 1 (`monthly_report`) como parámetros configurables en Ajustes. Viven en `notification.params_json`, con valor FIRME (decisión F) y lugar PROVISIONAL. C1-22 e INDICE §6 ya lo dicen. U4-12, en cambio, dice: «Cada parámetro es una clave de `app_setting`… **salvo `send_time` y `advance_days`**», y la lista de cambios (`U4.md:10`) agrega que «la columna “clave” admite `notification.params_json` **solo** para `send_time` y `advance_days`». Por la regla «los que inventaríe S1», U4 tiene que listar esos dos parámetros. Si les pone `params_json`, incumple U4-12; si les inventa una clave de `app_setting`, decide la opción B del bloqueo abierto. Además, U4-12 no pide la marca PROVISIONAL de su lugar, aunque INDICE §6 afirma que U4-12 cita el bloqueo sobre los cuatro. Es lo que falta para cerrar el hallazgo A, y contradice el esquema 1.3.
- **Corrección sugerida:** en U4-12 (enunciado, Verifica y las dos líneas de cambios), extender la excepción a `weekday` (7) y `day_of_month` (1): columna «clave» `notification.params_json` con PROVISIONAL (solo el lugar) que cite `S1-parametros-avisos`; valor inicial FIRME de la decisión F con «restablecer» (a diferencia de `send_time` y `advance_days`, sí tienen valor); sección que U4 proponga (p. ej. «Reportes y frecuencias»), justificada en SUPUESTOS. Cruzarlos aparte contra el `@param_sin_valor` de 00003 y el `params_json` sembrado por 00001, y excluirlos del cruce «clave de `app_setting`» mientras el bloqueo no tenga «Respuesta de Axel».

### H-2 (media) · S7-04 y S7-08 no leen `weekday` de su parámetro
- **Ruta:** `docs/criterios/S7.md:38-39` (S7-04) y `:54-55` (S7-08, alerta 4); `docs/criterios/INDICE.md:68`.
- **Qué pasa:** S7-10 ya lee `day_of_month` (= 1, decisión F) de su parámetro, con el lugar PROVISIONAL. S7-04 («el semanal el domingo») y S7-08 («(4) plan de comidas sin armar (domingo)») no piden lo mismo para `weekday` (= 7). Así, una implementación con el domingo fijo en el código pasaría S7-04 y S7-08, contra K44 y contra la 00003, que lo declara configurable. INDICE §6 dice que S7-04 y S7-08 citan el bloqueo sobre los cuatro parámetros, pero ningún criterio de S7 exige leer `weekday`.
- **Corrección sugerida:** en S7-04 y S7-08 (4), decir que el día se lee de `weekday` (7 = domingo, decisión F, FIRME; lugar PROVISIONAL según `S1-parametros-avisos`), y sumar `weekday` a las pruebas de configuración (cambio y restablecer).

### H-3 (media) · C1-12: la parte de B-01 de C1-r2 no se atendió y su Verifica literal falla
- **Ruta:** `docs/criterios/C1.md:81-82`.
- **Qué pasa:** `revisor-C1-r2.md:39` pide, dentro de B-01, actualizar C1-13, C1-21 y **C1-12 (excepción `x-app-setting`)**, y `INTEGRACION.md` §7 lo repite. La corrección solo hizo C1-13 y C1-21, y no justifica por qué dejó C1-12. Ejecuté la Verifica de C1-12 tal como está escrita sobre `openapi.yaml` (script `c112.py`): las operaciones públicas son `getHealth` y `login` (lo esperado); la búsqueda de `password|hash|token|secret` en los esquemas de respuesta encuentra `SecuritySettings.auth.password_min_length` (`getSecuritySettings`, `updateSecuritySettings`). Es un entero de configuración que K6 y K44 exigen exponer, pero el criterio solo admite «0 salvo la entrada de login». La Verifica falla con un contrato correcto.
- **Corrección sugerida:** en la Verifica de C1-12, excluir las propiedades con `x-app-setting` (claves de configuración, sin secretos) o nombrar `auth.password_min_length` como excepción, y anotarlo en la lista de cambios de C1.

### H-4 (baja) · P3-04 y P3-03: residuos de las regex
- **Ruta:** `docs/criterios/P3-claude-design.md:42` (P3-04) y `:38` (P3-03).
- **Qué pasa:** P3-04 corre `grep -c -E "mcp__claude-in-chrome__|navigate"` sobre `docs/estado/wf1/INTEGRACION.md`, que es un reporte en prosa y no un registro de llamadas. Si el integrador escribe «0 llamadas a `mcp__claude-in-chrome__`», el grep cuenta 1 con un proceso correcto. La vía alterna (declaración contrastada) evita el bloqueo, pero el propio criterio dice que «se cuentan acciones, no texto». Aparte, para no dar positivo con «Paso clave: exporta», P3-03 cambió `clave *[:=]` por «clave de acceso|clave secreta», y ahora una línea como «clave: Abc12345» ya no se detecta.
- **Corrección sugerida:** en P3-04, contar llamadas en el registro de herramientas del orquestador (las entradas de tipo herramienta con nombre `mcp__claude-in-chrome__*`), no texto de un reporte. En P3-03, volver a `clave *[:=] *[^ *]{6,}` y aceptar la coincidencia que el revisor del paquete justifique, como ya prevé el criterio.

### H-5 (baja) · S7-10: «su valor inicial» se puede leer sobre el día del mes
- **Ruta:** `docs/criterios/S7.md:62`.
- **Qué pasa:** «El día del mes del job (`day_of_month` = 1, decisión F, FIRME…) y su hora se leen de parámetros configurables… y su **valor inicial** lo propone el agente… PROVISIONAL hasta la puerta 4 (K31)». La frase admite la lectura de que el valor del día también es PROVISIONAL, cuando es FIRME.
- **Corrección sugerida:** escribir «el valor inicial de la hora lo propone el agente…».

### H-6 (baja) · U3: K38 sin traza en §5 y alcance ambiguo de la tarjeta (4)
- **Ruta:** `docs/criterios/U3.md:19` (U3-01, Verifica) y §5.
- **Qué pasa:** U3-01 cita K38, pero §5 no tiene la línea «decisiones K38 → U3-01». K38 nombra P-10, P-13 y P-14, y en U1 la fila N-023 de la tarjeta «Comidas de la semana» (`U1-arquitectura.md:464`) dice «Semana en curso». «La tarjeta (4) abre por defecto la semana que sigue K38» puede leerse como que la tarjeta lleva a P-10, que ya sigue K38 (lo que es coherente), o como que el contenido de la tarjeta muestra el domingo la semana siguiente (lo que choca con N-023). Además, `INTEGRACION.md` §7 anotó «`U3.md`: sin K38 ni K44»; la parte de K44 no se atendió ni se justificó.
- **Corrección sugerida:** sumar la línea de K38 a §5; precisar si K38 aplica al destino de la tarjeta (P-10) o a su contenido, en concordancia con U1 N-023; y decir en el reporte si U3 necesita K44 (p. ej. si la tarjeta (1) colorea con `alerts.budget_warning_pct` y `alerts.budget_limit_pct`) o por qué no.

## Estado de los hallazgos de la ronda 3
| Hallazgo r3 | Estado | Evidencia |
|---|---|---|
| A (alta) · parámetros de avisos | **Parcial** | Umbrales y días antes del pago: cerrados en C1-22, U4-12, S7-08, S7 §5 e INDICE §6/§11; la Verifica de C1-22 pasa contra la base y los extractos. Siguen abiertos `weekday` y `day_of_month` en U4-12 (H-1) y `weekday` en S7-04 y S7-08 (H-2). |
| B (media) · P3-06 | **Cerrado** | `P3-claude-design.md:49-50`: cuenta solo lo descargado, sin `siguiente-prompt.md`; asignación archivo o página → artboard en `P3-iter-N.md`; `P-xx` es solo sugerencia. |
| C (media) · P3-04 | **Cerrado**, con residuo bajo (H-4) | `:42`: ruta absoluta y conteo de herramientas de navegador, o declaración contrastada. |
| D (baja) · regex | **Atendido** (residuo en H-4) | Prueba propia: 6 advertencias dan 0, incluidas «**Contraseña:** nunca la pegues aquí» y «Paso clave: exporta…»; las credenciales se detectan. Hay una vía de justificación por escrito. |
| E (baja) · sección de `send_time`/`advance_days` | **Cerrado** | `U4.md:71`: «llevan la sección que U4 proponga… con su justificación en SUPUESTOS». |
| F (baja) · carpetas de salida | **Cerrado** | `P3-claude-design.md:70` incluye `salida/final/`; INDICE §3 (`:44`) pone P3-11 en la puerta 3. |
| G (baja) · 10 MB + 1 byte | **Atendido / justificado** | `S6.md:65`: el factor lo fija C1, con la propuesta de S1 citada (`S1-esquema-1.3.md:59`) y a confirmar en la puerta 2. |
| I-02 · `S4-meta-ajustes` | **Cerrado** | `S4.md:48`. |
| B-01 de C1-r2 | **Parcial** | C1-13 (1.2.0) y C1-21 (210 = 210) cerrados; C1-12 sin atender (H-3). |

## Verificación reejecutada (resumen)
- `bash docs/criterios/tools/verif_final.sh` y `GOOSE=… make verify-criterios`: rc = 0 los dos (detalle en la tabla).
- `goose -dir backend/db/migrations sqlite3 datos/tmp/rev4c/t.db up`: 00001, 00002 y 00003 OK, versión 3.
- `datos/tmp/rev4c/chk.py` (estructura, conteo y Decimal), `c122.py` (Verifica de C1-22 y la parte ejecutable de U4-12), `c112.py` (Verifica literal de C1-12), `prov.py` (PROVISIONAL contra bloqueos), `umb.py` (oraciones sobre umbrales), prueba de regex de P3 con textos de muestra propios y `verify_c1.py --db` sobre la base de goose (29 OK, TODO OK). Todo borrado al terminar.
- Nota: el agente `correccion-C1-r3` trabajaba en paralelo en `c1gen/` y `docs/contratos/`; mis corridas sobre el contrato son una foto de las 16:49 a 16:55 del 11-oct-2026.
