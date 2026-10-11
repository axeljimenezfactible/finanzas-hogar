# Bloqueo C1-decimales-tasa
Contexto: pregunta sugerida por la revisión r2 (R2-06) y r3; sin archivo propio hasta la puerta 1.
Qué falta: ¿cuántos decimales de tasa se muestran en pantalla? Las tasas se guardan como fracción exacta (2 q = 53450/309450 = 17.2726 %; 4 q = 80050/309450 = 25.8685 %).
Opciones: como se capturó; 4 decimales siempre; 2 decimales siempre.
Impacto: solo presentación (`CreditTermRate.rate_pct`); los cálculos usan la fracción. También: el sufijo «Z» de `DateTime` solo se acepta en entradas (documentar).
Cómo seguí: C1 v1.0 muestra `rate_pct` con 4 decimales.

## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Como se capturó: 17.27 % y 25.8685 %. Se guarda el número de decimales capturado (`rate_decimals` o equivalente) y la pantalla lo respeta; los cálculos usan la fracción exacta. Documentar que «Z» solo se acepta en entradas.
- **Estado:** RESUELTO. Se aplica en: C1 v1.1 (`CreditTermRate`), S1 1.2 (columna de decimales capturados), U4.
