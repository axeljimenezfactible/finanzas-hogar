# Bloqueo S2-ahorro-quincenal

Abierto el 11-oct-2026 en la ronda P2c por el agente `contratos-puerta2`, al aplicar Q2-13 = B (bloqueo `S2-flujo-quincenal`, pregunta 2). Se decide en WF-3 o cuando Axel conteste por id.

Contexto: la métrica M-10 (`docs/contratos/C1-datos.md` §9) pide por quincena el «ahorro planificado». Q2-13 = B fija el reparto: «mitad y mitad con HALF_UP en la 1 y el resto en la 2» del ahorro mensual planificado (`planned_monthly_cents`). No fija sobre qué base se reparte ni qué hace el saldo con él. `S2-interfaz.md` 1.3 §3.11 deja `savings_planned_cents` en `null`, como ya lo hacen el motor y el oráculo.

Qué falta (pregunta para Axel):
1. ¿El reparto se hace meta por meta (cada `planned_monthly_cents` se parte y luego se suma) o sobre la suma del mes (`planned_this_month_cents` de §3.13)? Ejemplo: dos metas de $1.01 (101). Meta por meta, la quincena 1 lleva 51 + 51 = 102 y la 2 lleva 50 + 50 = 100. Sobre la suma (202), cada quincena lleva 101.
2. ¿Qué metas cuentan en una quincena futura?
3. ¿El ahorro planificado resta del saldo (`balance_cents` y `cumulative_balance_cents`) o solo se informa? M-10 lo lista junto al saldo sin decirlo.

Opciones:
- Pregunta 1: (A) sobre la suma del mes; (B) meta por meta.
- Pregunta 2: (A) todas las metas activas con `planned_monthly_cents`, el mismo conjunto de `planned_this_month_cents`; (B) además, solo desde el mes de su `start_date` hasta el mes de su `target_date`, y sin las cumplidas.
- Pregunta 3: (A) solo se informa; (B) resta del saldo como un egreso comprometido.

Impacto: solo `savings_planned_cents` y, con la opción B de la pregunta 3, el saldo de las quincenas futuras (M-02 y M-10). Ninguna cifra canon depende de esto: las del flujo son los pagos de despensa por quincena (R-13…R-31).

Cómo seguí: `savings_planned_cents` sigue en `null` y no entra a ningún saldo. El reparto de Q2-13 = B queda escrito como FIRME en `S2-interfaz.md` 1.3 §3.11 para cuando se fije la base. Sin recomendación.

- **Estado:** ABIERTO
