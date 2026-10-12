# Bloqueo U4-reparto-preseleccion

Abierto el 11-oct-2026, noche, en la ronda P2c por el agente `correccion-U4-r1`, al corregir el hallazgo M-1 de `docs/estado/p2c/REVISION-U4-r1.md`. Se pregunta a Axel en la puerta 4.

Contexto: Axel eligió Q2-68 C (bloqueo `U4-reparto-retroactivo`, RESUELTO): al cambiar las partes de una regla de reparto reutilizable (por ejemplo «Conjunto 50/50»), P-26.12 pregunta el alcance en cada cambio, «solo los nuevos» o «también los ya capturados». La respuesta no dice qué opción aparece preseleccionada en el diálogo.

HECHO (11-oct-2026):
- Q2-68 C, literal: «Preguntar el alcance en cada cambio, como el presupuesto en K2 («solo los nuevos» o «también los ya capturados»).» No fija preselección.
- K2 sí la fija para el presupuesto: «de este mes en adelante» preseleccionado (`docs/plan/decisiones.md`). K46 es el modo de respuesta de la puerta 2 y no trata este punto.
- C1 v1.3 (`CostShareRuleUpdate.scope`, `onward` o `all`): «Obligatorio si viene parts […]. Sin valor por omisión: la pantalla pregunta el alcance en cada cambio (P-26.12).»
- `docs/design/U4-ajustes-b.md` 1.4 preseleccionaba «Solo a los gastos nuevos» como decisión del agente `U4-puerta2` (su supuesto 6, y la propuesta para Axel de `docs/estado/p2c/U4-puerta2.md`), sin marca ni bloqueo.

Qué falta (pregunta para Axel): en el diálogo de alcance de P-26.12, ¿qué opción aparece preseleccionada?

Opciones:
- A) «Solo a los gastos nuevos» preseleccionada. Consecuencia: con `Enter` o «Confirmar» sin elegir, los gastos ya capturados conservan sus partes y los estados de cuenta ya entregados a Karen no cambian (la razón de Q2-68 C). Es la que hoy describe U4.
- B) «También a los gastos ya capturados» preseleccionada. Consecuencia: con `Enter` o «Confirmar» sin elegir, se recalculan las partes de los gastos ya capturados con esta regla (cuáles, según el bloqueo abierto `C1-reparto-alcance`).
- C) Ninguna preseleccionada: «Confirmar» queda inactivo hasta que Axel elige una. Consecuencia: siempre hay una elección explícita y nada se aplica por inercia; un paso más en cada cambio.

PROPUESTA: el agente `U4-puerta2` propuso A (supuesto 6 de su versión). Este agente no recomienda: solo abre la pregunta.

Impacto: solo la preselección del diálogo y el comportamiento de `Enter` en P-26.12 (U4). No cambia C1: `scope` sigue obligatorio y sin valor por omisión con cualquier opción, porque la pantalla siempre manda el alcance elegido. Ninguna cifra canon cambia (la semilla no cambia ninguna regla de reparto).

Cómo seguí: `docs/design/U4-ajustes-b.md` 1.5 (P-26.12 y supuesto 6) marca **PROVISIONAL** la preselección con la opción en uso, A, citando este archivo.

- **Estado:** ABIERTO
