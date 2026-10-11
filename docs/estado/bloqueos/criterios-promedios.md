# Bloqueo criterios-promedios
Contexto: S2-08: plan.md §3.4 fila S2 pide «promedios diario/semanal/mensual». `model.json` solo tiene el promedio mensual (efectivo 2,183.96 = 13,103.78 / 6; devengado 2,393.75 = 14,362.48 / 6).
Qué falta: ¿Cómo se calculan el promedio diario y el semanal? Denominador: días reales del periodo (P1 = 31 días), 365/12 días por mes, 30 días por mes, o semanas de 7 días sobre el horizonte total (P1–P6 = 182 días).
Opciones:
- (a) Total del horizonte ÷ días reales del horizonte (10-oct-2026 a 9-abr-2027 = 182 días) y ÷ (días/7) para semanal.
- (b) Promedio mensual ÷ 30 y × 12 / 52 (convención simple).
- (c) Promedio mensual ÷ días del mes calendario correspondiente (varía por mes).
Impacto: Las cifras diarias y semanales no tienen valor canon; el criterio S2-08 las deja PROVISIONAL hasta que la fórmula esté escrita en `reglas-negocio.md` con ejemplo.
Cómo seguí: S2-08 exige que el motor y el oráculo coincidan con la fórmula escrita; el promedio mensual sí es canon (R-56, R-57).
