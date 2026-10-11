# Bloqueo criterios-ocr-tickets
Contexto: S5-06: plan.md §3.4 fila S5 pide «OCR ≥ 80 % renglones correctos en 5 tickets de prueba». No hay tickets en el repositorio.
Qué falta: (1) ¿De dónde salen los 5 tickets: reales que entrega Axel (tickets de Bodega Aurrera/Walmart de la compra del 10-oct u otros) o sintéticos generados? (2) ¿Qué es un «renglón correcto» (descripción + cantidad + importe exactos, solo importe, descripción con tolerancia)?
Opciones:
- (a) 5 tickets reales fotografiados por Axel (prueba más fiel; depende de él).
- (b) 5 tickets sintéticos renderizados que imitan papel térmico (el agente los puede generar, prueba menos fiel).
- (c) Mezcla: 3 reales y 2 sintéticos.
Impacto: La verificación del ≥ 80 % no puede declararse cumplida; con tickets sintéticos se declara «prueba limitada».
Cómo seguí: S5-06 exige fijar por escrito el origen y la definición antes de medir, con transcripción manual de los renglones esperados.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción (a): 5 tickets reales fotografiados por Axel (van en `datos/semilla/tickets/`). Renglón correcto = descripción legible + cantidad + importe exactos. S5 arranca sin las fotos y mide al final; si no están, el criterio queda «sin medir», no «cumplido».
- **Estado:** RESUELTO. Se aplica en: criterios S5-06, S5; acción de Axel: fotografiar 5 tickets.
