# Bloqueo U3-tablero-horizonte · cuántos meses ve cada gráfica del tablero (U3 2.0, §2.1)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54 («graficas de analisis… para irme previniendo»). Axel no fijó el horizonte. Mientras tanto rige la opción PROVISIONAL, que usa `docs/design/U3-dashboard.md` 2.0 (§2.1, G-1, G-6).

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **G-1 · ingreso, gasto y pagos por mes.** | (a) PROVISIONAL: los 11 meses anteriores y el mes en curso (lo real; el gasto futuro no existe). (b) 6 anteriores, el en curso y 5 siguientes, con el plan en lugar del gasto en los meses futuros. | Ventana de `getDashboardMonthlyOverview` y meses con datos. |
| 2 | **G-6 · tendencia del barómetro.** | (a) PROVISIONAL: los 6 meses anteriores, el en curso y los 5 siguientes. Los pagos futuros ya se conocen por el calendario de pagos y las quincenas esperadas por P-26.03, así que se ve antes un mes que se pasa del 40 %. (b) Los 12 meses anteriores. (c) El en curso y los 11 siguientes. | Ventana y número de puntos de proyección. |
| 3 | **Una o varias peticiones.** | (a) PROVISIONAL: una sola petición de 17 meses (`from` = en curso − 11, `to` = en curso + 5) y cada gráfica muestra su tramo. (b) Una petición por gráfica. | Rendimiento; no cambia cifras. |

G-4 (T-2: el mes en curso y los 11 siguientes) y G-5 (T-3: 12 meses) conservan las ventanas que ya aprobó la puerta 2.

**Qué cambia al contestar:** U3 2.0 §2.1, G-1, G-6, §6 y A-19; C1 v1.5 (límite de `from` y `to`).
