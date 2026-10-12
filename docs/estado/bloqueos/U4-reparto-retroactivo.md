# Bloqueo U4-reparto-retroactivo
Contexto: agente U4b-ajustes (WF-1, 11-oct-2026), Ajustes › Personas y reparto (P-26.12). El reparto por defecto del conjunto (50/50, decisión B) se cambia con `updateCostShareRule`, que **reemplaza** las partes de la regla (C1 v1.2, `CostShareRuleUpdate.parts`). En el esquema 1.3, `expense.cost_share_rule_id`, `planned_item.cost_share_rule_id` y los renglones de gasto apuntan a la regla y no guardan una copia de las partes, así que un cambio alcanzaría también a los gastos ya capturados y a los meses cerrados (y al estado de cuenta de Karen de esos meses).
Qué falta: cuando Axel cambia el porcentaje de una regla de reparto reutilizable (p. ej. «Conjunto 50/50»), ¿el cambio vale para los gastos ya capturados o solo para los nuevos?
Opciones:
- A. Vale para todo (como hoy en C1 y en el esquema). Consecuencia: la vista previa avisa que cambian los gastos ya capturados, incluidos meses anteriores y estados de cuenta ya entregados a Karen.
- B. Solo para los nuevos: el cambio crea una regla nueva (o una vigencia «desde») y los gastos anteriores conservan la suya. Consecuencia: S1 y C1 cambian (versión de regla o fecha de vigencia).
- C. Preguntar el alcance en cada cambio, como el presupuesto en K2 («solo los nuevos» o «también los ya capturados»). Consecuencia: cambio en C1 (parámetro de alcance) y en el motor de reparto.
Impacto: con el valor inicial (50/50) las cifras canon no cambian. Afecta la vista previa y el texto de confirmación de P-26.12 y, según la opción, S1, C1 y S2.
Cómo seguí: `docs/design/U4-ajustes-b.md` usa `updateCostShareRule` como está en C1 v1.2 y su vista previa muestra el ejemplo del nuevo reparto con la regla de K14 más el aviso del alcance, marcado PROVISIONAL citando este bloqueo. Nada entra en cifras ni en el esquema.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-68 · Si cambias el porcentaje de una regla de reparto reutilizable, ¿vale para los gastos ya capturados o solo para los nuevos?
- **Respuesta:** **C** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Preguntar el alcance en cada cambio, como el presupuesto en K2 («solo los nuevos» o «también los ya capturados»).»
- **Consecuencia:** cambio en C1 (parámetro de alcance) y en el motor de reparto.
- **Qué cambia:** C1 (`updateCostShareRule` con parámetro de alcance), motor de reparto (S2) y, según el alcance, S1; vista previa y texto de confirmación de P-26.12 en `docs/design/U4-ajustes-b.md` sin marca PROVISIONAL.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
