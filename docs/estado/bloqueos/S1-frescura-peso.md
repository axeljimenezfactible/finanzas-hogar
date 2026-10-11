# Bloqueo S1-frescura-peso
Contexto: Al formalizar la frescura de precios (reglas-negocio.md RN-12). La decisión E dice que un precio «se marca viejo a los 30 días y pesa menos en la comparativa». Es la misma pregunta del bloqueo `criterios-frescura-peso` (escrito en paralelo): una sola respuesta de Axel resuelve los dos.
Qué falta: (1) ¿Cuánto menos pesa un precio viejo en la comparativa de tiendas y canastas? (2) Confirmar que «a los 30 días» significa una edad de 30 días o más, por ejemplo que un precio del 10-oct-2026 ya es viejo el 9-nov-2026.
Opciones:
A. Peso fijo de 0.5 para precios viejos. Consecuencia: es simple; un precio viejo cuenta la mitad en promedios y en el puntaje.
B. Precio viejo solo de respaldo: se usa únicamente si no hay un precio fresco de ese producto en esa tienda, y se marca. Consecuencia: la comparativa nunca mezcla viejos con frescos.
C. Decaimiento lineal: peso 1 hasta los 29 días, que baja a 0 a los 90 días. Consecuencia: es más fino, pero más difícil de explicar en pantalla.
D. Sin peso, solo la marca visual. Consecuencia: contradice «pesa menos».
Impacto: S5 (mejor tienda por ingrediente y por canasta) y S7 (alerta de precios viejos). La clave `app_setting price.stale_weight` queda sin valor.
Cómo seguí: Dejé FIRME el umbral de 30 días (sembrado como `price.stale_after_days = 30`) con la lectura literal «30 días o más», y marqué el peso como PROVISIONAL, sin valor.
