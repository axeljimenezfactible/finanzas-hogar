# Bloqueo C1-huecos-u1-s1
Contexto: Al cruzar las necesidades de U1 con el esquema de S1 para escribir C1 aparecieron siete huecos (actualizado en la corrección r2 de C1, 11-oct-2026). El esquema 1.1 de S1 resolvió cuatro sin necesitar a Axel: descripción del gasto opcional, movimientos de inventario con motivo, tiempos aptos de una receta y monto por quincena. También agregó la tabla `idempotency_key`. C1 ya se regeneró contra ese esquema (`C1-datos.md` §14). Quedan tres preguntas de producto, las mismas del bloqueo `S1-huecos-negocio`: una sola respuesta de Axel resuelve los dos.
Qué falta:
1. Presupuesto (N-087, N-278): cuando Axel o Ximena cambian el monto mensual de una categoría en Ajustes o en Presupuesto, ¿el nuevo monto rige solo para ese mes o para ese mes y los siguientes hasta otro cambio?
2. Ticket por confirmar (N-075 «Guardar para después»): ¿se guardan las correcciones que la persona ya hizo a la lectura del OCR, o solo la lectura cruda?
3. Registro del chat de Telegram (N-299): ¿Ximena captura su chat_id a mano en Ajustes, o el sistema lo detecta cuando ella escribe al bot MultiX?
Opciones:
- Punto 1. A: rige para ese mes y los siguientes (renglón sin fin y cierre del anterior). B: rige solo para ese mes (renglón que empieza y termina en el mes). Consecuencia: C1 admite las dos con las mismas operaciones (`createBudget` y `updateBudget` con vigencia; un traslape responde 409 `budget_overlap`); cambia solo qué manda la pantalla.
- Punto 2. A: no se guardan las correcciones a medias, solo la lectura cruda (`attachment.ocr_text`); es lo que hace hoy `createTicketScan`. B: se guardan como borrador. Consecuencia: B pide una columna nueva en S1 y una operación para guardar el borrador en C1.
- Punto 3. A: captura manual del chat_id (`createNotificationRecipient`, dato solo de entrada que nunca se devuelve). B: el bot lo detecta y la persona lo confirma. Consecuencia: el esquema no cambia en ninguna; con B, S7 detecta el chat y C1 agrega la confirmación.
Impacto: lo que hace la API en N-087, N-278, N-075 y N-299 queda PROVISIONAL en `C1-datos.md` (§10, §13 y §14). Ninguna cifra canon depende de esto.
Cómo seguí: C1 usa las tablas del esquema 1.1 y deja PROVISIONAL solo estos tres puntos: el presupuesto acepta cualquier vigencia que no se traslape, «Guardar para después» guarda solo la lectura cruda y el chat se registra a mano. Ninguna necesidad queda sin operación. No inventé columnas ni valores.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Punto 1 (vigencia del presupuesto): configurable en cada cambio, «solo este mes» o «de este mes en adelante», con «de este mes en adelante» preseleccionado; la operación que cambia el monto recibe el alcance (`scope`). Punto 2 (ticket OCR «Guardar para después»): opción B, se guardan las correcciones como borrador (columna nueva en la migración 00002 y operación nueva en C1). Punto 3 (chat de Ximena): opción B, el bot lo detecta cuando ella escribe y ella confirma (S7 lo implementa; C1 suma la lectura de candidatos y la confirmación).
- **Estado:** RESUELTO. Se aplica en: C1 v1.1, S1 1.2 (migración 00002), U1 (N-087, N-275, N-278, N-075, N-299), S7 en WF-3.
