# Bloqueo S5-carritos-granel

Abierto el 11-oct-2026 en la ronda P2c por el agente `contratos-puerta2`, al aplicar Q2-03 = A (hallazgo H-r3-02 de `docs/estado/wf1/REVISION-S5-r3.md`). Se decide en WF-3 o cuando Axel conteste por id.

Contexto: `compare_carts` (`S5-interfaz.md` §7.2) suma los renglones de cada canasta. Con cantidades a granel (por ejemplo 0.5 kg), el importe de una línea, cantidad × precio, no es un número entero de centavos. El contrato no decía si ese importe se redondea por línea, como lo imprime un ticket, o si se suma exacto y se redondea una sola vez. El motor (`backend/internal/precios/compare.go`) redondea una sola vez; el oráculo (`_agrupa` en `oraculo/pre_calculo.py`) redondea cada línea y calcula el precio por unidad base desde el importe ya redondeado. El precio por unidad base siempre sale de la fracción exacta (S5-02: «aritmética exacta»); eso no está en duda.

Cifras del revisor (`REVISION-S5-r3.md`, H-r3-02; fixtures propias, ya borradas):
- `rv-carritos-granel`, dos renglones idénticos de 0.5 a 1.01 y 1.03: total de Bodega, motor 101 y oráculo 102; total de Walmart, 103 contra 104. `base_unit_price_cents`: 135, 137, 721 y 736 en el motor; 136, 139, 729 y 743 en el oráculo.
- `rv-carritos-granel-sust`, sustituto con 0.333 y 0.667: 2933 y 2800 en el motor; 2935 y 2801 en el oráculo.
Hoy no cambia el canon, porque todas sus cantidades son enteras, pero el ganador de un sustituto podría voltearse.

Qué falta (pregunta para Axel): en un carrito o una compra con cantidades a granel, ¿el importe de cada renglón se redondea como en el ticket, o se suma exacto?

Opciones:
- A) Exacto: el importe de cada línea es una fracción; el total de la tienda es HALF_UP de la suma exacta (un solo redondeo) y cada línea se muestra con HALF_UP de su importe (la suma de las líneas mostradas puede diferir del total en un centavo o dos). Es lo que hace el motor.
- B) Como el ticket: cada línea se redondea HALF_UP a centavos y el total es la suma de las líneas redondeadas; el precio por unidad base sigue saliendo de la fracción exacta. Según el informe, los totales del oráculo de hoy (102 y 104) son los de esta opción.
- C) Otra regla de redondeo por línea (por ejemplo, la que use cada tienda en su ticket), que habría que documentar con un ticket real.

Impacto: solo totales y diferencias de canastas con cantidades no enteras (1 o 2 centavos por renglón) y, en el peor caso, el ganador de un sustituto. Ninguna cifra canon.

Cómo seguí: `S5-interfaz.md` 1.5 §7.2 paso 8 escribe como PROVISIONAL la opción A, la del motor. Falta que el oráculo se alinee y agregar una fixture a granel (solicitud a los agentes de S5 en WF-3). Sin recomendación.

- **Estado:** ABIERTO
