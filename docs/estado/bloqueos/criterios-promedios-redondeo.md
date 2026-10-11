# Bloqueo criterios-promedios-redondeo
Contexto: al aplicar K20 (promedio diario = total del periodo ÷ días reales del periodo; semanal = diario × 7) al criterio S2-08 (agente `criterios-S-1.1`, 11-oct-2026), el semanal admite dos lecturas que difieren por centavos.
Qué falta: Para el promedio semanal, ¿se multiplica por 7 el diario **sin redondear** y se redondea una sola vez al final, o se multiplica el diario **ya redondeado a centavos**?
Opciones:
- (a) Diario exacto × 7 y un solo redondeo ROUND_HALF_UP al final. Octubre de 2026 (efectivo 973.75 ÷ 31 días): semanal 219.88. Febrero de 2027 (efectivo 2,517.38 ÷ 28): 629.35. Es la que usa el criterio mientras no haya respuesta.
- (b) Diario redondeado a centavos × 7. Octubre: 31.41 × 7 = 219.87. Febrero: 89.91 × 7 = 629.37. Lo que se ve en pantalla (diario) multiplicado por 7 siempre cuadra.
Impacto: solo el promedio semanal (hasta unos centavos por periodo). Ninguna cifra canon cambia: R-56 y R-57 son promedios mensuales y el diario es igual en ambas opciones.
Cómo seguí: S2-08 exige el mensual (R-56, R-57) y el diario exactos, y marca PROVISIONAL el semanal citando este bloqueo, con la opción (a) como lectura usada.

## Respuesta de Axel (puerta 1b, 11-oct-2026)
- **Decisión:** opción **(a)**: diario exacto × 7 y un solo redondeo ROUND_HALF_UP al final (octubre 2026: 219.88; febrero 2027: 629.35).
- **Estado:** RESUELTO (K39). Se aplica en: criterios S2-08 1.2, RN-27 1.3, S2 (WF-1).
