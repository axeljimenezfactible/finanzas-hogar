# Bloqueo S4-reglas-lista
Contexto: al escribir el contrato del motor S4 (`docs/contratos/motores/S4-interfaz.md`, WF-1, agente `contrato-S4S5`, 11-oct-2026). Tres reglas que C1 v1.2 expone como campos calculados no están definidas en decisiones §A–§K, en RN-15/RN-17 ni en los criterios S4.
Qué falta (tres preguntas para Axel):
1. **Sobrante previsto** (`MealPlanWeek.leftovers`, RN-15): ¿el plan cocina exactamente las porciones planeadas o la receta completa (múltiplos de sus 2 porciones base)?
2. **Reposición por ritmo** (`generateShoppingList`, motivo `replenishment`, RN-15 «más la reposición por ritmo»): ¿cuánto se agrega a la lista por un artículo llevado por ritmo de consumo y cuándo?
3. **Clasificación del IMC** (K7, S4-07 pide que RN-17 declare «la regla de redondeo antes de clasificar»; RN-17 no la declara): ¿se clasifica con el IMC exacto o con el que se muestra a un decimal?
Opciones:
- Pregunta 1. A (la que tomé, PROVISIONAL): se cocina lo planeado; el plan no genera sobrantes nuevos y `leftovers` muestra los sobrantes registrados que el plan no alcanza a consumir. B: se cocina la receta completa en múltiplos de las porciones base; la lista compra para el múltiplo y la diferencia es sobrante previsto (sube la lista). C: campo nuevo por casilla «porciones a cocinar» (cambio de esquema y C1).
- Pregunta 2. A (la que tomé, PROVISIONAL): si la existencia estimada al domingo de la semana es ≤ nivel de reposición, se agrega el consumo estimado de una semana (HALF_UP(ritmo × 7 / días del ritmo)); sin nivel no se repone. B: reponer hasta cubrir la semana más el nivel de reposición. C: agregar una presentación típica completa cuando toca reponer.
- Pregunta 3. A (la que tomé, PROVISIONAL): sobre `bmi_tenths`, el valor que se muestra (lo que la persona lee coincide con su categoría). B: sobre el valor exacto (18.4975 se muestra 18.5 pero es «bajo peso»).
Impacto: solo los campos `leftovers`, `needed_milli` de artículos por ritmo y `bmi_class` en la franja de ±0.05 de cada corte. Ninguna cifra canon del oráculo ni del reporte depende de esto.
Cómo seguí: el contrato S4 usa la opción A de cada pregunta marcada PROVISIONAL con cita a este bloqueo (§6.2 paso 3, §6.4 y §10.1); el resto del motor no depende de esto.
