# Bloqueo S1-huecos-negocio
Contexto: Corrección de S1 (ronda 1, hallazgo H-02 de `docs/estado/wf0/REVISION-r1.md`). De los siete huecos S1 ↔ U1 del bloqueo `C1-huecos-u1-s1`, cuatro se corrigieron en el esquema sin necesitar a Axel: descripción del gasto opcional (punto 1), movimientos de inventario con motivo (punto 2), tiempos aptos de una receta (punto 5) y monto por quincena (punto 6); también se agregó la tabla `idempotency_key`. Quedan tres preguntas de producto (puntos 3, 4 y 7), que son las mismas de `C1-huecos-u1-s1`: una sola respuesta resuelve los dos bloqueos.
Qué falta:
1. Presupuesto (N-087, N-275, N-278): cuando Axel o Ximena cambian el monto mensual de una categoría en Ajustes o en Presupuesto, ¿el nuevo monto rige solo para ese mes o para ese mes y los siguientes hasta otro cambio?
2. Ticket por confirmar (N-075 «Guardar para después»): ¿se guardan las correcciones que la persona ya hizo a la lectura del OCR, o solo la lectura cruda?
3. Registro del chat de Telegram (N-299): ¿Ximena captura su chat_id a mano en Ajustes, o el sistema lo detecta cuando ella escribe al bot MultiX?
Opciones:
- Punto 1. A: rige para ese mes y los siguientes (renglón con `valid_to_date` NULL y cierre del anterior). B: rige solo para ese mes (renglón que empieza y termina en el mes). Consecuencia: el esquema admite las dos (vigencia `valid_from_date`/`valid_to_date` en `budget`); cambia solo qué escribe la API.
- Punto 2. A: no se guardan las correcciones a medias, solo la lectura cruda (`attachment.ocr_text`, `ocr_status = 'done'`). B: se guardan como borrador. Consecuencia: B pide una columna nueva (por ejemplo un borrador JSON en `attachment`).
- Punto 3. A: captura manual del chat_id. B: el bot lo detecta y la persona lo confirma. Consecuencia: el esquema no cambia en ninguna (`notification_recipient.chat_id`, `verified_at`); cambia el trabajo de S7.
Impacto: el comportamiento de N-087, N-275, N-278, N-075 y N-299 en la API (C1) queda PROVISIONAL. Ninguna cifra canon depende de esto.
Cómo seguí: `budget` tiene vigencia y admite las dos lecturas del punto 1 (marcado PROVISIONAL en la columna `valid_to_date` y en RN-22). No agregué columna para el punto 2 ni cambié nada para el punto 3. Nada provisional entró en cifras.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Punto 1: vigencia configurable por cambio («solo este mes» o «de este mes en adelante»; preselección «en adelante»); el esquema con `budget.valid_to_date` ya lo admite. Punto 2: opción B, borrador de correcciones del OCR (columna nueva, migración 00002). Punto 3: opción B, el bot detecta el chat y Ximena confirma (tabla o columna de candidatos pendientes, migración 00002, sin guardar más que el chat_id y el nombre que da Telegram).
- **Estado:** RESUELTO. Se aplica en: S1 1.2 (migración 00002, diccionario, reglas), C1 v1.1, S7.
