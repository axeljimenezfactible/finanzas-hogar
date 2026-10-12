# Bloqueo WF3a-bis · semilla de desarrollo (C3-matriz punto 5, opción b)

12-oct-2026 · agente semilla-dev. Cada punto está marcado PROVISIONAL en el código y en `docs/estado/wf3/semilla-dev.md`. Nada de esto cambia cifras de los motores ni decisiones de Axel.

## 1. El gasto fijo (renta) y las cifras P1–P6 no caben juntos (PREGUNTA)
- **HECHO.** Con la renta de la fixture (`s2-proyeccion-csv`, partida 2: 8,000.00 mensual desde el 20-oct) `getAccrualVsCash` sube 800,000 centavos en devengado y en efectivo de **cada** mes (medido con el servidor: oct 194,750/97,375 pasa a 994,750/897,375). El motor S2 suma al mes las ocurrencias de gasto planificado (S2-interfaz §3.11, fórmulas de la proyección). El sueldo no mueve esas cifras.
- **Qué falta.** Las dos peticiones se contradicen: «partidas planificadas (gastos fijos)» y «cifras que sigan cuadrando con model.json (C3-09, P1–P6)».
- **Opciones.** (a) por omisión sin renta; con `--con-gasto-fijo` la agrega [elegida, PROVISIONAL]; (b) siempre con renta y el cotejo C3-09 de P1–P6 se corrige a «model.json + 8,000.00»; (c) renta ligada a un gasto capturado (no la cuenta dos veces, pero el gasto capturado suma igual).
- **Impacto.** Con (a), Inicio de hoy (12-oct) no cambia, porque la renta cae el 20-oct (quincena 16–31): solo la tarjeta de esa quincena mostraría un comprometido de 8,000.00. La prueba `TestDesarrolloConGastoFijoMueveLosMeses` deja escrito el efecto.

## 2. Candado «nunca en producción» (PROVISIONAL)
- Regla puesta: `finanzas semilla --desarrollo` se niega si `FINANZAS_ENTORNO=produccion`; si no, solo corre en una base cuyo archivo contiene «dev» (`datos/finanzas-dev.db`) o con `FINANZAS_SEMILLA_DEV=1`.
- **Falta** que WF-4 fije `FINANZAS_ENTORNO=produccion` en la unidad systemd y que `ops/runbook.md` lo diga (solicitud fuera de alcance en el reporte).
- Opciones: marca en `app_setting` (`entorno`), etiqueta de compilación (el binario de producción no incluye el overlay), o esta variable. Decide Axel/WF-4.

## 3. Presupuesto: solo Despensa (PREGUNTA)
- La única fixture con presupuesto de una categoría real es `s2-presupuesto-scope` (Despensa 2,000.00 desde oct-2026). No hay montos para Comida fuera ni Vivienda y no se inventaron. El total del mes es la suma de las categorías con plan (2,000.00).
- ¿Qué presupuesto quiere Axel para las demás categorías? Mientras tanto «Mes contra plan» muestra 119,764 de 200,000 (60 %, normal).

## 4. Hallazgo previo (no es del overlay): P3 a P6 de `getAccrualVsCash` no son las de model.json
- **HECHO.** Con la semilla real sola (sin overlay), oct y nov (P1, P2) cuadran con `model.json`; dic-2026 a mar-2027 no: devengado 503,474 / 0 / 503,474 / 0 contra 251,737 / 271,637 / 251,737 / 271,637, y efectivo de ene y mar 251,736 contra 271,636. Causa: `credit_purchase` 2 a 4 existen (installment_plan) pero sin gasto (`expense_id` NULL), así que su devengo no se reparte en 2 meses, y faltan los 2 papeles de contado de 199.00 (`semilla-oraculo.json` trae los gastos 2 a 6).
- **Impacto.** C3-09 «completo» (P1–P6 en pantalla) fallaría aunque el overlay no existiera, y `is_canonical` sale en true para esos meses.
- **Opciones.** (a) la semilla real agrega los gastos 2 a 6 de `semilla-oraculo.json` (futuros, proyectados) y los liga a sus compras; (b) un overlay de desarrollo hace lo mismo; (c) C3-09 compara solo P1 y P2 y deja P3 a P6 para cuando existan las compras. No lo toqué: es de la semilla real y cambia datos futuros.

## 5. Comidas de la semana: costo y kcal quedan en 0 (PREGUNTA)
- **HECHO.** Hay plan (6 recetas, lunes 12 a sábado 17, Axel y Ximena, `days_with_plan` = 6, `unplanned_slots` = 15). `cost_cents` = 0 y `avg_kcal_per_day` = 0 porque S4 no da costo a una receta con algún ingrediente sin precio (solo 36 de 160 productos tienen precio, los de la compra del 10-oct; la mejor receta tiene 4 de 9) y `nutrition_fact` está vacía (la semilla real no carga nutrición «propuesta», Q2-43 = A).
- **Opciones.** Que el agente externo (K53) cargue los precios que faltan; que Axel confirme filas de `nutricion/tabla.csv`; o aceptar el 0 en las capturas U5. No se inventó ningún precio ni caloría.

## 6. Ingresos: ninguno capturado
- La captura de la fixture `s2-liga-provisional` (1,195,000 el 14-oct) está en el futuro respecto de hoy (12-oct): no se cargó. Inicio muestra el sueldo de la quincena del 15 como esperado (`income_is_expected: true`, 12,000.00); el del 30 (12,500.00) sale de `second_half_amount_cents`.

## 7. Precios de Walmart en el overlay
- Los 36 precios de Walmart salen de `datos/semilla/data.json` (carrito comparado del 10-oct, fuente `manual`) y no están en la semilla real; con ellos «Mejor canasta» compara Bodega contra Walmart (Bodega 3,094.50, ahorro de 264.50 contra Walmart). ¿Se agregan también a la semilla real? Es dato real de Axel; hoy `price` queda vacía allí.
