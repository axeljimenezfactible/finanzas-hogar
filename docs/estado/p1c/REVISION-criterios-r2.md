# Revisión de criterios 1.2 · ronda 2 de 3 (flujo P1c, 11-oct-2026)

Revisor independiente `revisor-criterios-r2`. Evalúa `docs/criterios/` (diff contra HEAD `082235c`) solo contra `docs/criterios/INDICE.md` §1 y `docs/plan/decisiones.md` §K (K35–K44). No corrige nada. Todo se reejecutó con script propio (`datos/tmp/rev2c/audit.py`, ya borrado) y, como contraste, con `docs/criterios/tools/verif_final.sh`.

## Veredicto
**NO APROBADO.** Hallazgos: 0 altos, 4 medios y 6 bajos. Los 5 hallazgos de la ronda 1 (A-01, A-02, A-03, M-01 y M-02) quedan **cerrados**. Dos de los medios nuevos son de la misma clase que los anteriores: P3-02 sigue pidiendo su evidencia en un archivo que solo existe después de la puerta 3, y U4-12 pide un valor inicial que no está decidido.

## Tabla de verificaciones

| # | Verificación | Estado | Evidencia (reejecutada) |
|---|---|---|---|
| a | 18 archivos con las 5 secciones; cada criterio con «Verifica» y «Evidencia» y citado en §5 | CUMPLE | `audit.py`: 18 archivos con secciones [1,2,3,4,5]; 261 criterios, 0 sin Verifica ni Evidencia, 0 IDs repetidos y 0 IDs citados que no existan. WF4-07 aparece por el rango «WF4-01…WF4-12» (`traza.py`, con rangos expandidos: 0 sin línea en §5). |
| b | Cada K aplicada, con línea | CUMPLE | K35 P3:6,17,21,30,62,94,95; K36 S4:45, U4:67,71, C1:115; K37 S5:36,61; K38 U1:83; K39 S2:51; K40 S5:36,88; K41 U2:27,51; K42 U2:51; K43 C1:115,119, S3:66, S6:61, S7:35,51; K44 U4:67, C1:115, S3:66, S7:35,51,59. |
| c | 219.88, 629.35, 0.9836 y 0.4918 recalculados | CUMPLE | `oraculo/model.json` (sha256 `1a5706…bfd4`): efectivo 973.75 (P1, 31 días) y 2,517.38 (P5, 28 días). Con Decimal, el diario exacto × 7 y un solo HALF_UP da 219.88 y 629.35; con el diario redondeado da 219.87 y 629.37, que S2-08 rechaza. (90 − d) / 61 da 0.9836 (d = 30), 0.4918 (d = 60) y 0 (d = 90); 91 días no compite. |
| d | P3 1.2: nada del orquestador en el navegador; paquete completo; carpetas de salida | CUMPLE con hallazgos | P3-01, P3-04 y §3 prohíben que el orquestador abra claude.ai/design o use Chrome. P3-02 pide los 22 elementos de K35 con MANIFIESTO y sha256. P3-11 pide `salida/iter-N/` y `salida/final/`. Ver M-1 y M-4. |
| e | 0 PROVISIONAL que citen bloqueos resueltos | CUMPLE | `audit.py`: 10 líneas nombran un bloqueo resuelto junto a la palabra PROVISIONAL (S2:6,51; S3:6; S4:6; S5:6,36; U1:6; U2:6; C1:88,120). Todas dicen que la marca se quitó o comprueban que hay 0 marcas. Las marcas activas citan K31 (puerta 4) o el bloqueo abierto `S1-parametros-avisos`. |
| f | INDICE coherente con el conteo | CUMPLE con bajos | Conteo propio: 261 (201 / 60 / 0), igual a INDICE §1, §2 y §10, con versiones por archivo iguales a §10. `conteo.py` da 0 diferencias. §6: 25 filas (24 resueltas y 1 abierta) contra 43 archivos en `bloqueos/`; los dos sin «Respuesta de Axel» son `S1-parametros-avisos` (ABIERTO en §6) y `criterios-nxx` (resuelto con «Resolución»). Ver B-2 y B-3. |
| g | Ningún criterio reabre una decisión de §A–§K ni inventa valores | CUMPLE con hallazgos | No se reabre ninguna decisión. Los valores de prueba de S4-07 (300 y 250) están rotulados «fixture, no canon». Hallazgos: U4-12 exige el valor inicial de `send_time` y `advance_days` (M-2), y las Verifica de C1-22 y U4-12 no coinciden con dónde siembra el esquema 1.3 (M-3). |
| — | `verif_final.sh` (contraste) | exit 0 | 305 citas, 0 sin fuente; canon: 178 cifras, 0 diferencias; S1-09: 12 cifras, 0 diferencias; rutas `/tmp`: 0; conteo 261. El paso 9 sigue mostrando N-320…N-327 sin operación: es el B-01 de r1, en `C1-datos.md`, fuera de criterios. |

## Hallazgos

### M-1 (media) · P3-02 pide su evidencia en un archivo posterior a la puerta 3 (K35; INDICE §3)
- Ruta: `docs/criterios/P3-claude-design.md:28` (P3-02, «Evidencia: línea `P3-02` en `docs/estado/P3-iter-1.md`»).
- Qué pasa: INDICE §3 gobierna la puerta 3 con P3-01 y P3-02. K35 reduce esa puerta a «paquete listo y revisado». Pero `P3-iter-1.md` lo escribe el agente Sonnet después de que Axel avisa «iter-1 lista» (K35; P3-06 y P3-07). Por eso P3-02 no puede tener evidencia cuando se evalúa la puerta 3. P3-01 ya deja su evidencia en `docs/estado/wf1/`. Es la misma clase de defecto que M-02 de r1; esta línea no cambió desde la versión 1.1.
- Corrección sugerida: que la Evidencia de P3-02 sea una «línea `P3-02` en el reporte del agente `paquete-claude-design` y de su revisor en `docs/estado/wf1/`», con `ls -l` y `sha256sum -c`, como P3-01.

### M-2 (media) · U4-12 pide el valor inicial de `send_time` y `advance_days`, que no está decidido (K43, K31; coherencia con C1-22)
- Ruta: `docs/criterios/U4.md:67-68` (U4-12, enunciado y Verifica: «valor inicial … con los 11 parámetros … y 0 celdas vacías»).
- Qué pasa: K43 aprobó las claves y sus tipos, pero no sus valores. El bloqueo `C1-valores-k31` (punto 3) dice: «los valores por defecto los propone S7 con su tarea. Puerta 4». 00003 los declara `@param_sin_valor`. C1-22 lo dice explícitamente: «claves de K43 sin valor sembrado … su valor inicial lo propone S7-04/S7-08 con K31». U4-12 los cuenta entre los parámetros con «valor inicial», exige 0 celdas vacías y su marca PROVISIONAL cubre solo el lugar. Además U4 se construye en WF-1 y S7 en WF-3 (INDICE §2), así que la propuesta de S7 todavía no existe cuando se evalúa U4. El criterio empuja a inventar un valor o hace fallar la Verifica, y contradice C1-22.
- Corrección sugerida: en U4-12, para `send_time` y `advance_days`, pedir «sin valor inicial: lo propone S7 (K31), PROVISIONAL» y que esa celda no cuente como vacía. Precisar cuáles son los 11 parámetros, porque el recuento es ambiguo (kcal cuenta como 1 o 2; los días antes del pago no aparecen en la lista).

### M-3 (media) · Las Verifica de C1-22 y U4-12 buscan las claves solo en 00003, y el esquema 1.3 siembra 5 de ellas en 00001 y 00002 (K44)
- Rutas: `docs/criterios/C1.md:115-116` (C1-22: «valor inicial sembrado (migración 00003)», «lee las claves … sembradas en `00003*.sql` (o la migración vigente) … comentario que cita K-algo») y `docs/criterios/U4.md:68` (U4-12: «cruza la tabla con las claves de `app_setting` sembradas en `00003*.sql` (0 parámetros sin clave)»).
- Qué pasa: hay 5 claves de K44 que ya existían antes de 00003:
  - `auth.max_failed_attempts`, `auth.lockout_window_minutes` y `price.stale_after_days` se siembran en 00001, y su `description` cita «decisión A» y «decisión E», no un K-n;
  - `session.max_age_days` e `idempotency.retention_hours` se siembran en 00002.

  00003 solo las inventaría en sus líneas `@param`. Leída al pie de la letra, la Verifica de U4-12 marcaría como «sin clave» la sesión de 90 días, los 5 intentos en 15 min y el precio viejo a 30 días. La de C1-22 fallaría en el «comentario que cita K-algo» de tres claves, aunque el esquema esté bien.
- Corrección sugerida: «claves sembradas en 00001–00003 (o el inventario `@param` de 00003), con un comentario que cite su decisión de §A–§K»; en C1-22, quitar «(migración 00003)» o decir «00001–00003».

### M-4 (media) · Las búsquedas de «token» en P3 chocan con el contenido obligatorio (K35; INDICE §1)
- Rutas: `docs/criterios/P3-claude-design.md:22` (P3-01: `grep -rn -i "password\|contraseña\|token" …/LEEME.md` = 0), `:31` (P3-03: `grep … "contraseña\|password\|token\|clave" …/LEEME.md` = 0) y `:67` (P3-12: `grep … "password\|token\|api[_-]key" …/prompt.md` = 0).
- Qué pasa: el paquete de K35 incluye `C2-tokens.json` y `tokens.css`, y P3-02 exige que `prompt.md` incluya «los tokens de C2 en texto». Un `LEEME.md` que nombre los archivos que hay que subir, o una advertencia como «no escribas tu contraseña» (P3-03 pide advertir lo que no se hace), da coincidencias. Lo mismo un `prompt.md` con un apartado de tokens. Las Verifica de dos criterios altos de la puerta 3 fallarían con un paquete correcto. La de P3-01 es nueva en la versión 1.2; la de P3-12 viene de la 1.1.
- Corrección sugerida: buscar patrones de credencial y no la palabra suelta. Por ejemplo, `grep -n -i -E "password *[:=]|contraseña *[:=]|bot[_-]?token|api[_-]key|[0-9]{8,}:[A-Za-z0-9_-]{30,}"`. O bien excluir los nombres de archivo del paquete y la sección de advertencias.

### B-1 (baja) · El cambio de la versión 1.1 de S7 contradice a K43 sin aclararlo
- Ruta: `docs/criterios/S7.md:15` («S7-04, S7-08 y S7-10: contenido, horas, anticipaciones y definiciones … PROVISIONAL hasta entonces (K31)»).
- Corrección sugerida: anotar «(en la versión 1.2, la definición de «casillas sin armar» es FIRME por K43)». En S7-08, aclarar quién «cite el bloqueo abierto»: es la marca, no Axel.

### B-2 (baja) · La regla de estado de §6 no cubre `criterios-nxx`
- Rutas: `docs/criterios/INDICE.md:68` y `docs/criterios/C1.md:87` (C1-15: «bloqueo abierto = sin «Respuesta de Axel»»).
- Qué pasa: `criterios-nxx.md` está resuelto con una sección «Resolución» y sin pregunta para Axel. Con esa regla contaría como abierto. Hoy ninguna marca lo cita.
- Corrección sugerida: «o con «Resolución» cuando no hubo pregunta para Axel».

### B-3 (baja) · Filas de INDICE desactualizadas
- Rutas: `docs/criterios/INDICE.md:95` y la tabla de §11.
- Qué pasa: la fila de `C1-valores-k31` en §6 no lista S7-04 ni S7-08 (`send_time`, `advance_days` y «casillas sin armar») ni C1-22. La fila B-01 de §11 dice «§6 con los 24 bloqueos», cuando hoy son 25.
- Corrección sugerida: completar la fila de §6 y actualizar el número en §11.

### B-4 (baja) · INDICE §7 no tiene fila para §K
- Ruta: `docs/criterios/INDICE.md` §7 (líneas 98–111).
- Qué pasa: la tabla de trazabilidad con `decisiones.md` cubre de §A a §J. No tiene fila para §K (K1–K44), aunque casi todos los cambios de 1.1 y 1.2 vienen de ahí.
- Corrección sugerida: agregar una fila K con las frases cubiertas y los archivos.

### B-5 (baja) · Comandos nuevos con rutas relativas o elididas (INDICE §1, «Rutas y comandos: absolutos»)
- Rutas: P3-03 (`… ` elidido), P3-05 (`grep -c "DESIGN.md" LEEME.md`), P3-11 (`… LEEME.md`, que además no comprueba `salida/final/`), C1-22 (`grep -n "mínimo de seguridad\|minimum:" ` sin archivo) y U4-12 (`grep … ...`).
- Corrección sugerida: usar rutas absolutas bajo `/home/deepsec/finanzas-hogar/` y nombrar los archivos.

### B-6 (baja) · S7-08 (4) no dice qué semana se evalúa el domingo
- Ruta: `docs/criterios/S7.md:51`.
- Qué pasa: «la alerta se dispara si hay al menos una» es una lectura razonable de la decisión F y de K43, pero no está escrita en ellas. Tampoco se dice qué semana se revisa el domingo; K38 hace que el plan abra la semana siguiente.
- Corrección sugerida: decir que se evalúa la semana que empieza el lunes siguiente (coherente con K38) y presentar «al menos una» como lectura de la decisión F y K43.

## Estado de los hallazgos de la ronda 1

| Hallazgo r1 | Estado | Comprobación |
|---|---|---|
| A-01 (alta) · C1-22 fijaba claves de avisos en `app_setting` sin marca | **CERRADO** | C1-22 (`C1.md:115-116`) excluye los avisos, los deja en `notification.params_json` con la marca PROVISIONAL citando el bloqueo abierto `S1-parametros-avisos` y los comprueba aparte. Su Verifica ya no espera esas claves en 00003. Queda M-3, de otra causa. |
| A-02 (alta) · U4-12 contra K43 | **CERRADO** | U4-12 (`U4.md:67-68`): la columna «clave» admite `notification.params_json` con la marca y cita el bloqueo; el cruce con 00003 los excluye. Quedan M-2 y M-3, nuevos. |
| A-03 (alta) · S7-08 con «sin armar» PROVISIONAL; S7 sin configurables | **CERRADO** | S7 1.2: (4) es regla FIRME sin marca, con casos 21/1/«fuera»; §3 dice que la marca es un defecto; S7-04, S7-08 y S7-10 leen `send_time`, `advance_days`, los umbrales y los días antes del pago de su parámetro, con prueba de cambio y de «restablecer». |
| M-01 (media) · INDICE §6 sin `S1-parametros-avisos` | **CERRADO** | Fila ABIERTO en `INDICE.md:96`; la frase de estado dice 25 filas, 24 resueltas y 1 abierta, lo que coincide con `bloqueos/`. |
| M-02 (media) · P3-01 pedía el aviso de Axel | **CERRADO** | P3-01 es solo «paquete listo y revisado», con evidencia en `docs/estado/wf1/`; el aviso pasa a P3-06. Queda M-1, la misma clase en P3-02. |
