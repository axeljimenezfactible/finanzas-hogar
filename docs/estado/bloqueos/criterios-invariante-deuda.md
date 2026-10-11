# Bloqueo criterios-invariante-deuda
Contexto: Escribiendo S2.md (criterio S2-07). plan.md §3.4 fila S2 pide el invariante «devengado acumulado − efectivo acumulado = deuda viva». Lo comparé con `oraculo/model.json`.
Qué falta: ¿Qué significa «deuda viva» en ese invariante? Cálculo con `oraculo/model.json`, deuda viva = Σ totales de compras hasta la fecha − Σ pagos hasta la fecha: cierre P1 (9-nov-2026) deuda 2,921.25 vs diferencia 973.75; P2 973.75 vs 973.75 (cumple); P3 3,776.06 vs 1,258.69; P4 1,258.70 vs 1,258.70 (cumple); P5 3,776.06 vs 1,258.69; P6 (9-abr-2027) 1,258.70 vs 1,258.70 (cumple); y el 10-abr-2027, con la despensa 4 comprada ese día, 6,293.44 (`horizonte.deuda_viva_10_abr`) vs 1,258.70 (`horizonte.diferencia_final`). La forma literal solo se cumple cuando no queda parte de una compra sin devengar (P2, P4, P6). Se cumple siempre la identidad: deuda viva = (devengado acumulado − efectivo acumulado) + saldo no devengado de las compras (comprobado en los 6 cierres y en el 10-abr).
Opciones:
- (a) «Deuda viva devengada» = devengado acumulado − efectivo acumulado (la diferencia del oráculo, 1,258.70 al final): el invariante queda como identidad trivial y la cifra 6,293.44 se llama «deuda viva total» (consecuencia: dos nombres, ambos canon).
- (b) Invariante algebraico: deuda viva total = (devengado acumulado − efectivo acumulado) + saldo no devengado de las compras vivas (P1: 2,921.25 = 973.75 + 1,947.50; el 10-abr-2027: 6,293.44 = 1,258.70 + 5,034.74). Es verificable por prueba en los 6 cierres (consecuencia: reescribir el criterio del plan).
- (c) Cambiar el oráculo para devengar de otra forma (consecuencia: invalida cifras canon; no se recomienda).
Impacto: S2-07 queda PROVISIONAL (alta una vez resuelto). El motor y el revisor no pueden declarar cumplido el invariante literal del plan.
Cómo seguí: S2.md S2-07 exige verificar la opción (b) como identidad exacta, marcada PROVISIONAL, y que la diferencia acumulada coincida con `periodos[i].diferencia_acumulada`. Ninguna cifra cambió.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción (b): identidad con el saldo no devengado: deuda viva total = (devengado acumulado − efectivo acumulado) + saldo no devengado (10-abr-2027: 6,293.44 = 1,258.70 + 5,034.74). Se reescribe el criterio S2-07; ninguna cifra cambia.
- **Estado:** RESUELTO. Se aplica en: criterios S2-07, S2.
