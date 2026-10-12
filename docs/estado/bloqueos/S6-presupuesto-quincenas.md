# Bloqueo S6-presupuesto-quincenas · plan de la quincena en getBudgetFortnights

Agente: api-presupuesto, 12-oct-2026. PROVISIONAL. C1 pide `planned_cents` «parte del plan de la quincena (motor S2)», pero el motor no reparte el presupuesto mensual por quincena. Hoy: Q1 = plan × 15 ÷ días del mes (ROUND_HALF_UP), Q2 = el resto (las dos suman el plan). Opciones: a) por días (actual); b) mitad y mitad; c) función nueva del motor S2. El real de la quincena es el devengado del hogar en la ventana (`motor_fin.Periods` con ventanas). Impacto: la gráfica de quincenas de P-06.
