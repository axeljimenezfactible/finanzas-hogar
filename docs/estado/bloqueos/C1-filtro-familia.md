# Bloqueo C1-filtro-familia
Contexto: Al definir los 6 filtros globales de C1 (plan §3.7: persona, familia, categoría, tienda, periodo, crédito) con una sola definición en `docs/contratos/openapi.yaml` (`components/parameters/Family`). U1 §0.2 define «Familia: Hogar o Karen», pero la necesidad N-040 (lista de gastos) ofrece «familia (Hogar, Karen, Conjunto)». S1 guarda las partes de dinero a nombre del titular (hogar) y de personas externas (Karen); la co-titular no lleva parte propia.
Qué falta: (1) ¿El filtro «familia» tiene dos valores (Hogar, Karen) o tres (Hogar, Karen, Conjunto)? (2) ¿«Hogar» incluye los registros conjuntos (todo lo que tiene parte del hogar) o solo los que son 100 % del hogar?
Opciones:
- A (la que tomé, PROVISIONAL). Tres valores `household`, `external`, `joint` con semántica incluyente: household = con parte del hogar (incluye conjuntos), external = con parte de Karen (incluye conjuntos), joint = con ambas partes. Consecuencia: «todo lo de Karen» es un solo filtro (sirve para su participación y su estado de cuenta).
- B. Tres valores exclusivos: household = solo hogar, external = solo Karen, joint = conjuntos. Consecuencia: ver «todo lo de Karen» pide dos valores a la vez (el parámetro tendría que aceptar lista).
- C. Dos valores (Hogar, Karen) incluyentes, sin «Conjunto». Consecuencia: N-040 pierde la opción «Conjunto» y U1 debe corregirse.
Impacto: Solo el significado del parámetro `family` en listExpenses, listCreditPurchases, listScheduledPayments, listCalendarEvents, getExpenseTotals, getExpenseSummary, getExpenseMonthlyTrend y queryView. Ninguna cifra canon depende de esto.
Cómo seguí: C1 v1.0 usa la opción A, marcada `x-provisional` en el parámetro y PROVISIONAL en `C1-datos.md` §6. Se registró la inconsistencia de U1 como solicitud fuera de alcance.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción A: tres valores incluyentes (Hogar, Karen, Conjunto). «Hogar» incluye lo conjunto; «Karen» incluye lo conjunto. Deja de ser PROVISIONAL.
- **Estado:** RESUELTO. Se aplica en: C1 v1.1 (parámetro `family`), U1 N-040.
