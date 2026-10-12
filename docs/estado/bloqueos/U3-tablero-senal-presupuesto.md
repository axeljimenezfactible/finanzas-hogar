# Bloqueo U3-tablero-senal-presupuesto · regla «subir o recortar presupuesto» por categoría (U3 2.0, G-3)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54, propósito del sistema según Axel: «para ver si aumentamos presupuesto o lo recortamos para irme previniendo al respecto». El orquestador dio un ejemplo («tres meses seguidos arriba o abajo del plan por más de un umbral configurable»), pero Axel no fijó la regla. Mientras tanto rige la opción PROVISIONAL. Ya la usan `docs/design/U3-dashboard.md` 2.0 (G-3, §2.6) y la operación propuesta `getBudgetSignals`. La calcula el motor S2 y la interfaz solo la muestra.

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **Cuántos meses.** | (a) PROVISIONAL: **3 meses cerrados** antes del mes pedido (clave nueva `budget.signal_months`, valor inicial 3, de 2 a 12). El mes en curso no cuenta porque no ha cerrado. (b) 2 meses. (c) 6 meses. | Qué tan pronto aparece la sugerencia. |
| 2 | **Umbral.** | (a) PROVISIONAL: un solo margen de **10 %** para los dos lados (clave nueva `budget.signal_margin_pct`, valor inicial 10): «subir» si en los 3 meses el % usado pasó de 110; «recortar» si quedó abajo de 90. (b) Márgenes distintos para subir y recortar (p. ej. 10 y 20). (c) Comparar el promedio de los meses contra el plan, sin racha. | Cuántas categorías llevan sugerencia. |
| 3 | **Meses sin plan o sin gasto.** | (a) PROVISIONAL: un mes sin plan de esa categoría corta la racha (sin datos no es cero); un mes con plan y sin gasto cuenta como 0 % usado. (b) Se saltan sin cortar la racha. | Categorías nuevas o de temporada. |
| 4 | **Qué dice y dónde.** | (a) PROVISIONAL: insignia de rol información «Sugerencia: subir presupuesto» o «Sugerencia: recortar presupuesto» en G-3, sin monto sugerido. (b) Además un monto sugerido (p. ej. el promedio real de los 3 meses). (c) Además un aviso de S7 el día 1. | Texto, y en (b) y (c) una cifra o un aviso nuevos. |

**Qué cambia al contestar:** U3 2.0 G-3 y §2.6; C1 v1.5 (`BudgetSignals`); esquema 1.5 (las dos claves); motor S2 (función nueva de sugerencia, con su prueba).
