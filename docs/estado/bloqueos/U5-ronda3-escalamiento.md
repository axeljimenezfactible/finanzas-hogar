# Bloqueo U5-ronda3-escalamiento

Contexto: WF-3b, `revisor-u5-r3` (12-oct-2026). U5-05 permite 3 rondas de revisión visual. La ronda 3 queda **devuelta** con 0 hallazgos altos y 2 medios (`docs/estado/wf3/REVISION-u5-r3.md`):
- M-7: barra lateral de escritorio sin los iconos de sección que piden DESIGN §7.9 y §5.4 (40 pantallas, Ec y Eo).
- M-8: U5-11, la matriz de `C3-integracion.md` no marca «capturado» con sus 4 rutas (41 filas). Es una actualización del documento, no de la aplicación.

Todos los hallazgos de la ronda 2 (1 alto y 4 medios) quedaron cerrados, y los criterios altos U5-01, 02, 03, 04, 07, 08, 09 y 13 pasan con comandos reejecutados.

Qué falta: decisión de Axel. U5-05 dice «si la ronda 3 no aprueba, se escala a Axel (no hay ronda 4)».

Opciones (PROPUESTA, no decidido):
| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. Autorizar una corrección acotada a M-7 y M-8 y una verificación puntual de esos 2 puntos | U5 cerrable sin abrir otra ronda completa | Bajo (iconos en un componente; columnas de C3) | Bajo |
| B. Aceptar U5 con M-7 y M-8 como pendientes con fecha | Se avanza a la puerta 4 | Ninguno ahora | La barra lateral queda distinta de DESIGN y C3 incompleto en la puerta |
| C. Pedir otra ronda completa | Revisión completa de nuevo | Medio (≈ 9 min de capturas y una revisión) | Contradice el límite de 3 rondas de U5-05 |

Recomendación del revisor: A.

Estado: ABIERTO (lo decide Axel).
