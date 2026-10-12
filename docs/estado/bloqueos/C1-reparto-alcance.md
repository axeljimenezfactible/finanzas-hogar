# Bloqueo C1-reparto-alcance · fase C1 (v1.3, puerta 2)

Contexto: agente C1-puerta2 (P2c, 11-oct-2026). Axel eligió Q2-68 C (`U4-reparto-retroactivo`): al cambiar las partes de una regla de reparto reutilizable se pregunta el alcance, «solo los nuevos» o «también los ya capturados». C1 v1.3 lo expone como `CostShareRuleUpdate.scope` (`onward` o `all`) en `updateCostShareRule` (C1-datos §7.2). RN-03 1.4 dice que «qué gastos abarca "también los ya capturados" lo fijan C1 y el motor de reparto en WF-3». HECHO: en el esquema 1.4 un gasto puede tomar la regla por cuatro caminos: `expense.cost_share_rule_id`, `expense_item.cost_share_rule_id`, `category.default_cost_share_rule_id` (gasto con regla NULL) o la regla por defecto (categoría sin regla); `expense_share` guarda solo montos.

Qué falta: ¿qué gastos recalcula `scope = all`?

Opciones:
- A) Solo los que nombran la regla: `expense.cost_share_rule_id` o un renglón con `expense_item.cost_share_rule_id` igual a la regla. Consecuencia: un gasto que la toma por su categoría o por ser la de por defecto no cambia.
- B) Todos los que la usan, también por su categoría o por ser la de por defecto. Consecuencia: coincide con lo que se ve en pantalla («este gasto es Conjunto»), y cambia más estados de cuenta.
- C) B, limitado a un rango de fechas que se elige al confirmar (p. ej. «desde el mes en curso»). Consecuencia: cambia C1 (un campo más) y la pantalla P-26.12.

PROPUESTA del agente: B, porque «también los ya capturados» no distingue cómo llegó la regla al gasto. Sin decisión de Axel no se aplica.

Impacto: solo `scope = all`. `scope = onward` no toca nada guardado y no depende de esto. Ninguna cifra canon cambia (el valor inicial es 50/50 y no se cambia en la semilla). Lo usan C1 (texto de `CostShareRuleUpdate.scope` y §7.2), el motor de reparto (WF-3) y la vista previa de P-26.12 (U4).

Cómo seguí: C1 v1.3 marca PROVISIONAL solo el alcance de `all` (`CostShareRuleUpdate.scope`, C1-datos §7.2 y §13), citando este archivo.
