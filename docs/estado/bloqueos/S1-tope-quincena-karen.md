# Bloqueo S1-tope-quincena-karen
Contexto: pregunta sugerida por la revisión r2 (R2-03); sin archivo propio hasta la puerta 1. `oraculo/model.json` dice «fuera del cap de 374.93» en dos notas de `papel`; S1 trató 374.93 como la parte de Karen por pago, no como un tope.
Qué falta: ¿el tope por quincena de Karen para el objetivo «pago parejo y bajo tope» es 374.93 o lo fija Axel?
Opciones: sin tope fijo (lo captura Axel por persona); 374.93; otro monto.
Impacto: `person.fortnight_payment_cap_cents`; solo alimenta el objetivo 2 del comparador, no es restricción dura.
Cómo seguí: columna vacía; ninguna cifra entró al esquema ni a C1.

## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Sin tope fijo: Axel lo captura por persona en Ajustes › Personas y reparto cuando quiera. 374.93 fue la parte de Karen por pago en la despensa 1, no un tope. La descripción de `person.fortnight_payment_cap_cents` se redacta neutral (R2-03).
- **Estado:** RESUELTO. Se aplica en: S1 1.2 (descripción de la columna), S3, U4.
