# Bloqueo S1-reparto-edicion-gasto

Abierto el 11-oct-2026, noche, en la ronda P2c por el agente `correccion-S1-r1`, al corregir el hallazgo M-02 de `docs/estado/p2c/REVISION-S1-r1.md`. La pregunta la dejó el agente S1-puerta2 en «Pendientes» de `docs/estado/p2c/S1-puerta2.md` (punto 2), sin bloqueo. Sale de la respuesta de Axel a Q2-68 (opción C, `docs/estado/puerta2/RESPUESTAS.md` y `PUERTA-2.md`). Se decide en WF-3 o cuando Axel conteste.

Contexto: Q2-68 C pide que, al cambiar las partes de una regla de reparto reutilizable (por ejemplo «Conjunto 50/50»), la pantalla pregunte el alcance: «solo los nuevos» o «también los ya capturados» (RN-03, invariante 18 del diccionario). La respuesta no dice qué pasa si, después de un cambio con «solo los nuevos», alguien edita un gasto anterior.

HECHO:
- `expense_share` guarda por gasto y persona solo `amount_cents`, sin el porcentaje ni el monto fijo de la regla, y sin su versión (`backend/db/migrations/00001_esquema.sql`, tabla `expense_share`: «Se congela al guardar para que cambiar una regla no altere gastos pasados»).
- `cost_share_rule_part` guarda cada parte como `pct_bp` o `fixed_cents`. `updateCostShareRule` reemplaza las partes de la regla (C1 v1.2, `CostShareRuleUpdate.parts`), así que después del cambio la regla de un gasto anterior ya trae las partes nuevas.
- `updateExpense` (C1 v1.2, `docs/contratos/extractos/gastos.yaml`, operación `updateExpense`): «Edición en el lugar con guardado automático; recalcula las partes. `items` reemplaza los renglones.» Es decir, hoy recalcula con la regla vigente.
- El reparto de un gasto lo calcula la API al guardarlo (RN-03: los renglones se agrupan por regla, una parte con porcentaje vale HALF_UP(suma × `pct_bp` / 10000), una parte fija vale `fixed_cents` y la persona del resto recibe lo demás). S2 no reparte gastos registrados, porque recibe las partes congeladas (`docs/contratos/motores/S2-interfaz.md`, frontera con S3).

Qué falta (pregunta para Axel): después de cambiar una regla de reparto con «solo los nuevos», si se edita un gasto anterior (su importe o sus renglones), ¿con qué partes se recalcula?

Opciones (sin decidir; esta ronda no las recomienda):
- A) **Con la proporción de las partes congeladas del gasto.** Cada persona recibe HALF_UP(importe nuevo × su parte guardada ÷ importe anterior) y la persona del resto absorbe los centavos. Consecuencia: sin cambio de esquema. Una parte fija (`fixed_cents`, por ejemplo la parte fija de Karen en despensas futuras, 119150 en el diccionario) se escalaría como si fuera porcentaje, incluso sin cambio de regla. Si cambian los renglones, la proporción de todo el gasto no refleja reglas distintas por renglón. Además, `expense_share` no dice quién es la persona del resto.
- B) **Guardar cómo se calculó cada parte.** El porcentaje o el monto fijo de cada parte (o la versión de la regla) se guarda en `expense_share`, y el recálculo usa lo guardado. Consecuencia: es exacto también con renglones y partes fijas, pero pide una migración posterior a la 00004 (esquema 1.5, con su Down), un cambio en C1 y que el reparto de la API lo use.
- C) **Con la regla vigente.** Es lo que hace hoy `updateExpense` en C1 v1.2. Consecuencia: sin cambio de esquema ni de C1, pero el gasto editado pasa al reparto nuevo aunque se eligió «solo los nuevos». Si es de un mes cerrado, cambia el estado de cuenta de Karen de ese mes. La pantalla de edición tendría que avisarlo.

En A y B también falta decidir cómo se sabe que un gasto quedó fuera de un cambio: por ejemplo, comparando la fecha de sus partes con la del cambio de la regla, o con una marca.

Impacto:
- C1: `updateExpense` y quizá un campo o parámetro.
- El reparto de la API al guardar un gasto (WF-3).
- S1: con la opción B, una migración nueva.
- S7: estado de cuenta de Karen de meses ya entregados (opción C).
- U4/P-02: texto o aviso de la edición en el lugar.

Ninguna cifra canon cambia mientras no se cambie una regla con «solo los nuevos». La opción A sí cambia el resultado de editar un gasto cuya regla tiene una parte fija.

Cómo seguí: el esquema no cambia (1.4). RN-03 (`docs/contratos/reglas-negocio.md`) y la invariante 18 de `docs/contratos/diccionario-datos.md` marcan este caso PROVISIONAL, citando este bloqueo, con la opción que ya está en uso: la C, porque `updateExpense` de C1 v1.2 recalcula con la regla vigente. Nada entra en cifras.

- **Estado:** ABIERTO
