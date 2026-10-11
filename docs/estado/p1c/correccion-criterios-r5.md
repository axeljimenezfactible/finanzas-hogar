# correccion-criterios-r5 · Ronda 5 de criterios 1.2 (flujo P1e) · 11-oct-2026

## Hecho
- **H-1 (alta) U4-12 e INDICE §6**: la excepción del bloqueo abierto `S1-parametros-avisos` cubre los cuatro parámetros `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1). Columna «clave» `notification.params_json` con PROVISIONAL solo sobre el lugar y cita del bloqueo; `weekday` y `day_of_month` con valor inicial FIRME (decisión F) y «restablecer»; `send_time` y `advance_days` sin valor inicial (S7, K31, puerta 4); sección propuesta por U4 con justificación en SUPUESTOS; cruce aparte contra el `@param_sin_valor` de 00003 y el `params_json` sembrado por 00001, fuera del cruce «clave de `app_setting`» hasta la «Respuesta de Axel». La lista mínima de la Verifica pasa de 14 a 16 filas (15 `weekday`, 16 `day_of_month`); no cambia ningún total de criterios. Enunciado, Verifica y las dos líneas de cambios (r1 y r4) de U4 actualizados; INDICE §6 (texto de la sección y fila del bloqueo), filas de §11 (r1/A-02 y r2/hallazgo 2) con la nota.
- **H-2 (media) S7-04 y S7-08 (4)**: el día se lee de `weekday` (7 = domingo, decisión F, FIRME; lugar PROVISIONAL según el bloqueo); `weekday` entra a las pruebas de configuración (cambio y restablecer) de S7-04 y S7-08.
- **H-3 (media) C1-12**: la Verifica excluye las propiedades con `x-app-setting` y nombra `auth.password_min_length` (`SecuritySettings`) como excepción; anotado en la lista de cambios de C1.md. Confirmado en `openapi.yaml:16809` que la propiedad lleva `x-app-setting`.
- **H-4 (baja) P3-04 y P3-03**: P3-04 cuenta llamadas de herramienta `mcp__claude-in-chrome__*` en el registro de herramientas del orquestador (journal del flujo WF-1 o la lista de herramientas de su reporte), no texto en prosa. P3-03 vuelve a `(password|contraseña|clave) *[:=] *[^ *]{6,}` y acepta la coincidencia que el revisor del paquete justifique por escrito (ej. «Paso clave: exporta…»).
- **H-5 (baja) S7-10**: «el valor inicial de la hora lo propone el agente…»; `day_of_month` FIRME.
- **H-6 (baja) U3**: §5 suma «decisiones K38 → U3-01» (cita literal de K38); U3-01 precisa que K38 aplica al destino de la tarjeta (4) (P-10 abre la semana según K38) y su contenido sigue U1 N-023 «Semana en curso». U3 no necesita K44: ningún criterio de U3 fija ni lee umbrales, días o anticipaciones operativas (U3-05 usa los límites de quincena de K8, que son calendario); los umbrales `alerts.*` los consumen S7-08 y C1-22/U4-12; si el diseño de la tarjeta (1) o (3) coloreara con ellos, leería `getAlertSettings`. Anotado en la lista de cambios de U3.md.
- Listas de cambios «1.2, ronda 5» agregadas en U4, S7, C1, U3 y P3; encabezado de versión de cada uno y del INDICE actualizado; INDICE §11 con 6 filas nuevas (H-1 a H-6) y la fila de conteo renombrada «rondas 4 y 5». Conteos sin cambio: 261 (201 / 60 / 0).

## Archivos
Modificados (parches con python `str.replace` con conteo exacto 1): `docs/criterios/U4.md`, `S7.md`, `C1.md`, `U3.md`, `P3-claude-design.md`, `INDICE.md`. Creado: este reporte. Estado real confirmado: `00003_puerta1b.sql` líneas 56-59 traen los cuatro `@param_sin_valor` (send_time, advance_days, weekday 7, day_of_month 1); `grep x-app-setting ajustes-a.yaml` da líneas 1058 y 1068 (`projection.repeat_cadence_enabled`).

## Verificación
- `bash docs/criterios/tools/verif_final.sh`: las 11 comprobaciones en lo esperado: 18 archivos, 261 criterios, 0 faltantes; refs 261/261; citas 306 verificadas / 0 sin fuente; N-xx 0; canon 178 / 0 diferencias; S1-09 SQL OK (54 tablas) y equivalencia 12 / 0; matriz 210 operaciones, 0 N-xxx sin operación; rutas `/tmp` 0; **«INDICE contra conteo propio: 0 diferencias»** (261: 201 / 60 / 0). Una primera corrida dio 1 cita sin fuente («Semana en curso» entre «» en la línea de K38 de U3 §5); corregido quitando las comillas angulares y repetida: 0.
- `make verify-criterios GOOSE=<scratchpad>/gobin/goose`: «las 11 comprobaciones dan lo esperado».
- `grep -n 'weekday\|day_of_month'`: aparecen en U4 (líneas 6, 9, 13, 74-75: U4-12), S7 (42-43 S7-04; 58-59 S7-08; 66 S7-10; 116 §5) e INDICE (§6 líneas 68 y 96; §11).
- `grep -n 'x-app-setting\|password_min_length' docs/criterios/C1.md`: líneas 6 (cambios) y 85 (Verifica de C1-12).
- Regex de P3-03 (`(password|contraseña|clave) *[:=] *[^ *]{6,}|bot[_-]?token…|api[_-]?key…|[0-9]{8,}:…`) con `grep -c -i -E`: 5 advertencias («**Contraseña:** nunca la pegues aquí», «No escribas tu contraseña…», etc.) = 0 coincidencias; 5 credenciales de muestra (password, contraseña, clave, bot_token, api-key) = 5 coincidencias. Archivos de muestra en el scratchpad, borrados.
- `grep -rn PROVISIONAL docs/criterios/`: los bloqueos citados son solo `S1-parametros-avisos` (14 menciones) y K31 (30); las demás líneas con «bloqueos/…» son referencias genéricas de patrón (U1-14, S3-20, C1-15) o avisos de que la marca se quitó (S2:6, S2:51, `criterios-promedios-redondeo`, RESUELTO con «Respuesta de Axel»). Bloqueos sin «Respuesta de Axel»: `S1-parametros-avisos` y `criterios-nxx`.

## Supuestos
- La lista mínima de U4-12 pasa a 16 filas para que `weekday` y `day_of_month` cuenten como parámetros con valor inicial; es una consecuencia mecánica de H-1 y no cambia el conteo de criterios. El «14 parámetros» de la línea histórica de r2 en U4 (línea 11) se deja como historial; el INDICE §11 lo anota.
- En P3-04 no se inventa el formato del journal de WF-1: el criterio dice «entradas de herramienta cuyo nombre empieza con `mcp__claude-in-chrome__`» y que el revisor indique la ruta que usó.
- U3: la línea de K38 en §5 cita K38 hasta «la que empieza el lunes» (el verificador de citas toma como cita el texto entre «» y rechazó las comillas anidadas).

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (sin «Respuesta de Axel»); no se tocó.

## Pendientes
- Cuando `S1-parametros-avisos` traiga «Respuesta de Axel», quitar las marcas PROVISIONAL y la excepción en C1-22, U4-12, S7-04, S7-08 y S7-10 (INDICE §6 lo dice).

## Solicitudes fuera de alcance
Ninguna.

## Propuestas para Axel
Ninguna nueva (la pregunta de lugar de los cuatro parámetros ya está en el bloqueo abierto).

## Consumo aproximado
17 acciones de herramienta; leído ≈ 120 KB (protocolo, bloqueo, REVISION r4, tramos de los seis criterios, salidas de verificación); escrito ≈ 38 KB (parches y reporte).
