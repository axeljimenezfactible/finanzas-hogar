# Bloqueo S6-semilla-despensas-futuras

Contexto: WF-3a, agente `semilla` (12-oct-2026). La tarea pide «la despensa 1 con su calendario de 4 quincenas tal como lo reproduce motor_cred», pero la prueba de cruce C3-09 pide una deuda viva al 10-abr-2027 de 6,293.44, y esa cifra solo existe si las despensas 3 y 4 (1,258.70 pendientes de la 3 + 5,034.74 de la 4) están en la base. `docs/plan/decisiones.md` (línea «Datos semilla») dice «crédito Kueski con 4 despensas y pagos (`oraculo/model.json`)» y `canontest.Seed` carga las cuatro.

Qué falta: ¿las despensas 2, 3 y 4 (compras futuras: 10-dic-2026, 10-feb-2027 y 10-abr-2027) se siembran como compras a crédito `credit_purchase` activas, sin gasto, o se dejan fuera y las produce la proyección a partir de la partida «Despensa» de `planned_item`?

Opciones:
- **A (la que apliqué, PROVISIONAL)**: 4 `credit_purchase` (la 1 enlazada a su `expense`; la 2, 3 y 4 sin gasto, con la nota «Despensa programada del oráculo») con sus 16 pagos y sus partes, todo calculado por `motor_cred` y cotejado contra `model.json`. Consecuencia: C3-09 sale sobre la base sembrada (A 10,000.00 / B 6,905.50, pago 973.75, deuda viva 6,293.44); el calendario de pagos y la línea de tiempo del disponible ya muestran las 4 despensas. Riesgo: son compras que aún no ocurren; si Axel cambia el plan (otro monto o fecha) debe editarlas, y si además captura una partida «Despensa» en `planned_item` (CSV) la proyección podría contarlas dos veces.
- **B**: sembrar solo la despensa 1; las demás las crea la proyección desde la partida programada. Consecuencia: C3-09 (6,293.44) y la línea de tiempo A/B solo se cumplen si el motor de proyección crea `credit_purchase` virtuales; hoy no está definido.
- **C**: A, pero con `status = 'cancelled'` hasta que se compren. Consecuencia: la deuda viva no las cuenta (6,293.44 no sale).

Impacto: solo estas 3 filas de `credit_purchase` (y sus pagos/partes). No cambia ninguna cifra del oráculo.

Cómo seguí: opción A, marcada PROVISIONAL en cada compra futura (`credit_purchase.notes`) y en este bloqueo. Para pasar a B basta borrar esas 3 compras (hace falta una operación de Ajustes o SQL).

- **Estado:** ABIERTO
