# Bloqueo S6-ajustes-api · fase WF-3a (agente api-ajustes, 12-oct-2026)

Contexto: al implementar las 19 operaciones de Ajustes contra C1 v1.4 quedaron seis puntos que C1, S1, RN y U4 no fijan. Se implementó lo marcado como PROVISIONAL y se sigue con lo demás.

## 1. `PlannedItem.next_date` (x-calculado, opcional)
Qué falta: `motor_fin` calcula las ocurrencias de una partida (`occurrencesIn`, sin exportar) y los motores no se tocan. Recalcularlo en la API sería «recalcular a mano».
Opciones: A) exportar en `motor_fin` un envoltorio `NextOccurrence(p PlannedItem, from Date) (Date, bool)` sin cambiar cifras (lo aprueba Axel); B) dejar `next_date` siempre ausente; C) calcularlo en la API (descartada: regla de motores).
Cómo seguí: B (PROVISIONAL): `next_date` se omite. Impacto: P-26.02/P-26.03 no muestran «próxima» hasta A.

## 2. `updateCostShareRule` con `scope = all`
Qué falta: el bloqueo `C1-reparto-alcance` (qué gastos abarca) sigue ABIERTO y `motor_fin` no exporta el reparto (`splitRule`), que es el «motor de reparto» que pide C1-datos §7.2.
Opciones: A) exportar `SplitRule(amount Cents, r *ShareRule) ([]Share, error)` en `motor_fin` y, con la decisión de Axel sobre el alcance, recalcular `expense_share` en la misma transacción; B) responder 409 sin escribir nada hasta entonces.
Cómo seguí: B (PROVISIONAL): valida todo y responde 409 `conflict` «Por ahora el cambio no se puede aplicar a los gastos ya capturados. Elige «Solo los nuevos».». `onward` funciona completo.

## 3. `reorderCategories`: qué es «lista completa»
Opciones: A) todas las categorías, activas y archivadas; B) todas las activas, las archivadas opcionales.
Cómo seguí: B (PROVISIONAL): faltan activas, ids repetidos o inexistentes → 422 `validation_failed` en `ids`; `sort_order` queda 1..n en el orden recibido.

## 4. Orden por defecto sin texto en C1
`listPlannedItems`: `start_date`; `listCostShareRules`: `id` (Hogar, Karen, Conjunto en el orden de alta). PROVISIONAL.

## 5. Metas desactivadas
`motor_fin.SavingsGoals` solo informa metas activas (Q2-15). Una meta con `is_active = false` sale sin `saved_cents`, `progress_*`, `required_monthly_cents` ni `status`. Opciones: A) así; B) pedir al motor también las inactivas. PROVISIONAL A.

## 6. PATCH sin campos y valores fuera de enum
`{}` (minProperties 1) y valores fuera de la lista de C1 (`kind`, `frequency`, `scope`, `status`…) responden 400 `bad_request` (formato, C1 §3); las reglas de negocio, 422. Un requerido ausente en un alta, 422 `required` (CONVENCIONES §4). PROVISIONAL hasta que C1 lo diga por operación.

Impacto general: ninguna cifra canon cambia; solo afecta respuestas de Ajustes. Lo usan U4 (P-26.02, P-26.03, P-26.04, P-26.12, P-26.13) y frontend.

- **Estado:** ABIERTO
