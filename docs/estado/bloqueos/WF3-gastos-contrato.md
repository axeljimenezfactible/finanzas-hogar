# WF3-gastos-contrato · Huecos de C1 v1.4 para P-02, P-04 y P-05

Abierto por `front-gastos-presupuesto` (WF-3a, 12-oct-2026). Estado: PREGUNTA para Axel. Lo afectado quedó PROVISIONAL en `frontend/src/pantallas/gastos-presupuesto/`.

| # | Qué falta | Opciones | Impacto hoy |
|---|---|---|---|
| 1 | `Expense` no trae el id de su adjunto (solo `has_attachment`), así que P-04 «Ver ticket» y las miniaturas de la galería de P-02 (N-046) no pueden pedir `downloadAttachment` | (a) agregar `attachments: [{id, kind, mime_type}]` a `Expense`; (b) operación nueva `listExpenseAttachments`; (c) dejarlo fuera de la v1 | «Ver ticket» se ve desactivado con su motivo; la galería muestra tarjetas sin foto |
| 2 | La edición en el lugar del monto cambia las partes (hogar y Karen), pero C1 no tiene una vista previa calculada del reparto (U4 §0.2) | (a) operación de vista previa (p. ej. `POST /expenses/{id}/preview`); (b) aceptar que la vista previa muestre solo el monto «Antes → Después» y que las partes se recalculen al guardar; (c) editar el monto solo en P-04 con una hoja | Hoy está (b): monto con «Antes → Después» y partes recalculadas por el servidor; el reparto no se edita en el lugar |
| 3 | `getExpenseTotals` y `getExpenseSummary` no aceptan `kind`, `payment_method` ni `has_attachment`, y N-041 pide los totales del mismo conjunto que N-040 | (a) agregar esos tres filtros a las dos operaciones; (b) ocultar la fila de totales cuando haya esos filtros | Con esos filtros, la fila de totales cuenta todo el periodo y la familia elegidos |
| 4 | P-05 «emparejado con catálogo» (N-072): falta definir qué operación busca productos desde la pantalla del ticket | (a) usar la búsqueda de productos del área Compras; (b) operación propia en Gastos | Los renglones muestran el texto leído y su cuenta; no hay buscador de producto |
