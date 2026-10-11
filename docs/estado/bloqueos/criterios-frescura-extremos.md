# Bloqueo criterios-frescura-extremos
Contexto: al aplicar K22 al criterio S5-04 (agente `criterios-S-1.1`, 11-oct-2026). K22 dice: «Peso 1 hasta el día 29; entre 30 y 90 días decae en línea recta de 1 a 0». Tomada al pie de la letra, la recta vale 1 en el día 30, igual que un precio fresco, aunque el precio ya está marcado como viejo.
Qué falta: ¿Cuánto pesa un precio de 30 días y cómo cambia día a día entre el 30 y el 90?
Opciones:
- (a) Lectura literal: la recta va de 1 en el día 30 a 0 en el día 90, peso = (90 − días) / 60 (día 30 = 1.0000; día 60 = 0.5000; día 89 = 0.0167). El precio de 30 días está marcado viejo pero todavía pesa igual que uno fresco. Es la lectura que usa el criterio mientras no haya respuesta.
- (b) La recta parte del 1 del día 29 y llega a 0 en el día 90, peso = (90 − días) / 61 (día 30 = 0.9836; día 60 = 0.4918). Todo precio viejo pesa estrictamente menos que uno fresco.
- (c) Otra pendiente o escalones; Axel los indica.
Impacto: solo el peso de los precios de 30 a 89 días en la comparativa de S5. Ninguna cifra canon cambia (la comparativa del 10-oct usa precios frescos). El peso de esos días queda PROVISIONAL.
Cómo seguí: S5-04 prueba lo que ambas lecturas comparten (peso 1 hasta el día 29, 0 en el día 90, 91 días no compite, peso que no sube y recta lineal) y deja el valor de los días 30 a 89 PROVISIONAL citando este bloqueo, con la lectura (a).

## Respuesta de Axel (puerta 1b, 11-oct-2026)
- **Decisión:** opción **(b)**: peso = (90 − días) / 61 para 30 ≤ días ≤ 90 (día 30 = 0.9836; día 60 = 0.4918; día 90 = 0); peso 1 hasta el día 29; después del día 90 no compite, solo se marca.
- **Estado:** RESUELTO (K40). Se aplica en: criterios S5-04 1.2, RN-12 1.3, S5 (WF-1).
