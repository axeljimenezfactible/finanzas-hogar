# Bloqueo criterios-optimizador
Contexto: S3-13 y S3-14: decisiones C define 4 objetivos y «puntaje ponderado por defecto (pesos iguales normalizados)». No define cómo se pasan los 4 objetivos a una escala común.
Qué falta: (1) ¿Cómo se normaliza cada objetivo antes de ponderar (mín-máx entre los escenarios generados, relativo al mejor, rango fijo)? (2) ¿Qué mide «pago por quincena parejo y bajo tope»: desviación estándar de los pagos, diferencia máx-mín, pico contra el tope? (3) ¿Qué mide «comprar más por peso»: unidades base por peso gastado de la canasta, usando qué precios?
Opciones:
- (a) Mín-máx entre los escenarios generados (0 = peor, 1 = mejor); parejo = desviación estándar de pagos; peso = unidades base por $; cada definición con ejemplo.
- (b) Relativo al mejor escenario de cada objetivo (cociente).
- (c) Que el agente proponga las tres definiciones en `reglas-negocio.md` y Axel las apruebe en la puerta 2.
Impacto: La recomendación ponderada y la medida de los objetivos 2 y 4 quedan PROVISIONAL; los 4 óptimos por objetivo individual no dependen de la normalización.
Cómo seguí: S3-13 pide óptimos por objetivo independientes (4 distintos en 3 casos) y S3-14 pide la definición escrita; sin ella, PROVISIONAL.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción (a): normalización mín-máx entre los escenarios generados; «pago parejo y bajo tope» = desviación estándar de los pagos (y pico contra el tope cuando hay tope por persona); «comprar más por peso» = unidades base de la canasta por peso gastado, con los precios vigentes de cada tienda.
- **Estado:** RESUELTO. Se aplica en: criterios S3-13 y S3-14, S3.
