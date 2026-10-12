# Bloqueo S2-metas-vencida

Abierto el 11-oct-2026 en la ronda P2c por el agente `contratos-puerta2`, al fijar Q2-16 = A y Q2-18 = A en `S2-interfaz.md` 1.3 §3.13. Se decide en WF-3 o cuando Axel conteste por id.

Contexto: Q2-16 = A fija `months_left` = meses del mes de hoy al mes de `target_date`, los dos incluidos (0 si la fecha meta es de un mes anterior). Q2-18 = A da el estado `overdue` (vencida). Ninguna de las dos dice qué aportación necesita una meta vencida que no se cumplió (`months_left` 0 y restante mayor que 0). La fórmula del plan, base HALF_UP(restante ÷ n), no está definida con n = 0.

HECHO: el motor (`backend/internal/motor_fin/projection.go`) da `plan` = [{mes de hoy, restante}] y `needed_monthly_cents` = restante. El oráculo (`oraculo/fin_metas.py`) da `plan` vacío y `needed_monthly_cents` null. Ninguna fixture tiene una meta vencida (`s2-metas` trae una meta en curso, una a la mitad, una cumplida y un fondo).

Qué falta (pregunta para Axel): ¿qué aportación muestra la pantalla de Metas para una meta vencida que no se cumplió?

Opciones:
- A) Todo el restante en el mes de hoy: `plan` = [{mes de hoy, restante}] y `needed_monthly_cents` = restante (lo que hace el motor).
- B) Sin aportación calculada: `needed_monthly_cents` null y `plan` vacío, como una meta sin fecha (lo que hace el oráculo).
- C) `needed_monthly_cents` = restante y `plan` vacío.

Impacto: solo la tarjeta de una meta vencida en la pantalla de Metas. Ninguna cifra canon.

Cómo seguí: `S2-interfaz.md` 1.3 §3.13 marca PROVISIONAL la opción A, la del motor, y el oráculo se alinea a ella. Sin recomendación.

- **Estado:** ABIERTO
