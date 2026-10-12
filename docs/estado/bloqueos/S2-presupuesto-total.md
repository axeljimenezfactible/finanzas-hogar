# Bloqueo S2-presupuesto-total
Contexto: interfaz del motor financiero, presupuesto del mes (`S2-interfaz.md` §3.8; RN-22; M-01). El esquema admite un renglón de presupuesto con `category_id` NULL («presupuesto total del mes») y C1 lo trata como una categoría más para los traslapes, pero M-01 define «plan = Σ `budget.amount_cents` del renglón vigente en el mes de las categorías activas».
Qué falta: (1) si existe un renglón del total, ¿el plan total es ese monto o la suma de las categorías? (2) Con un plan de 0.00 en una categoría, ¿qué porcentaje y qué banderas se muestran?
Opciones: (1) A) el renglón del total manda y la suma de categorías se muestra aparte; B) la suma de categorías manda (M-01) y el renglón del total solo se informa; C) el total es el mayor de los dos. (2) A) sin porcentaje y sin banderas (estado «sin plan»); B) cualquier gasto mayor que 0 enciende las dos banderas.
Impacto: `total.plan_cents` cuando existe el renglón del total y el caso de plan 0 quedan PROVISIONAL; las banderas por categoría con plan mayor que 0 (S2-11) no dependen de esto.
Cómo seguí: la interfaz aplica M-01 tal como está escrito (suma de categorías activas), informa el renglón del total en `declared_total_row_cents` y, con plan 0, no da porcentaje ni banderas; las fixtures `s2-alertas-80-100` y `s2-presupuesto-scope` no usan el renglón del total ni planes de 0.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-19 · (1) Si existe un renglón de «presupuesto total del mes», ¿el plan total es ese monto o la suma de las categorías?
- **Respuesta:** **B** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «B) la suma de categorías manda (M-01) y el renglón del total solo se informa»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `total.plan_cents` cuando existe el renglón del total y el caso de plan 0 quedan PROVISIONAL; las banderas por categoría con plan mayor que 0 (S2-11) no dependen de esto.
- **Qué cambia:** `total.plan_cents` (M-01 tal como está escrito); el renglón se informa en `declared_total_row_cents` (`S2-interfaz.md` §3.8): se quita la marca PROVISIONAL.
- **Estado:** RESUELTO

### Q2-20 · (2) Con plan $0.00 en una categoría, ¿qué porcentaje y banderas se muestran?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «A) sin porcentaje y sin banderas (estado «sin plan»)»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `total.plan_cents` cuando existe el renglón del total y el caso de plan 0 quedan PROVISIONAL; las banderas por categoría con plan mayor que 0 (S2-11) no dependen de esto.
- **Qué cambia:** Caso plan 0 de §3.8: sin porcentaje ni banderas, estado «sin plan»; la interfaz queda con esa regla FIRME.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
