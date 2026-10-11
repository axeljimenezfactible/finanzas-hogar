# Bloqueo criterios-login-bloqueo
Contexto: S6-06 y S6-07: decisiones A fija «límite de intentos (5 en 15 min)» y «sesión larga por dispositivo (meses)».
Qué falta: (1) Duración del bloqueo tras 5 intentos fallidos (¿hasta que pasen los 15 min, otro plazo?) y si el contador es por usuario, por IP o ambos. (2) Duración exacta de la sesión larga (¿60, 90, 180 días?).
Opciones:
- (a) Ventana deslizante de 15 min por usuario+IP; rechazo mientras haya 5 fallos en la ventana; sesión de 90 días renovable con uso.
- (b) Bloqueo fijo de 15 min tras el 5.º fallo; sesión de 180 días.
- (c) Que S6 proponga valores en C1 y Axel los apruebe en la puerta 1.
Impacto: Comportamiento PROVISIONAL en C1 y S6; el umbral 5/15 sí es canon.
Cómo seguí: S6-06 y S6-07 prueban 5 en 15 min y una sesión de al menos 2 meses (plural de «meses»); el resto marcado PROVISIONAL.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción (a): ventana deslizante de 15 minutos por usuario e IP (el bloqueo se levanta cuando en los últimos 15 min hay menos de 5 fallos); sesión de 90 días renovable con cada uso. Dejan de ser PROVISIONALES.
- **Estado:** RESUELTO. Se aplica en: criterios S6-06 y S6-07, C1 v1.1, S6.
