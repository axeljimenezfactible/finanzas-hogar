# Bloqueo S5-carrito-comparativa
Contexto: corrección del contrato S4/S5 (WF-1, agente `correccion-contrato-S4S5-r1`, 11-oct-2026), hallazgo M1 de `docs/estado/wf1/REVISION-contrato-S4S5-r1.md`. El canon de la comparativa (Bodega Aurrera 3,094.50 contra Walmart 3,233.00, +138.50; idénticos +100.50; 61 / 60 / 55 piezas) es **carrito contra carrito**: cada tienda conserva sus propias cantidades (3 × 750 g contra 3 × 1 kg de arroz, 780 g contra 640 g de sustituto de crema, 4 contra 6 litros de leche, 6 contra 5 latas de atún). `data.json` dice «Precios del carrito ya con descuentos e IVA; sin envío» y solo la compra de Bodega se registró y se repartió (`cuentas` = 61 piezas = 3,094.50). En C1 v1.2, `getBasketComparison` con `basket = expense` recibe **una sola** `expense_id` (`compras-tiendas-a.yaml`) y no dice cómo se arma la canasta de las otras tiendas. Sembrar Walmart como un gasto de 3,233.00 sería un hecho no documentado.
Qué falta (preguntas para Axel):
1. Con `basket = expense`, ¿contra qué se compara la compra registrada en las otras tiendas?
2. ¿Con qué fuente se siembran los precios de la compra de Bodega del 10-oct-2026: `manual` (captura a mano sin ticket) o `ticket_manual` (hay ticket impreso)?
3. ¿La semilla de la app guarda los precios del carrito de Walmart del 10-oct-2026 en `price` con fuente `manual`, o solo existen como entrada del motor y de sus fixtures?
Opciones:
- Pregunta 1. **A)** Cada otra tienda se cotiza con sus precios vigentes que compiten (S5 §6) para los productos y cantidades de la compra base, eligiendo presentación como en S5 §7.3. No cambia C1. La API **no** reproduce +138.50 con la semilla (Walmart cotizado con las cantidades de Bodega da otro total); el canon K-01 a K-08 se reproduce en el motor (`CompareCarts`) y en la fixture `s5-comparativa-canon`. **B)** C1 agrega un segundo parámetro (`compare_expense_id`, o un carrito guardado con sus renglones) y la API compara compra contra compra o contra carrito: reproduce el canon por la API; cambia C1 (operación y quizá una tabla de carritos en el esquema). **C)** Las dos: A por omisión y B cuando se pasa el segundo parámetro.
- Pregunta 2. **A)** `manual`. **B)** `ticket_manual`. En los dos casos el canon es el mismo, porque compiten con el mismo peso y la misma frescura (K37).
- Pregunta 3. **A)** Solo entrada del motor y de la fixture: la app arranca sin precios de Walmart. **B)** También en `price` con fuente `manual` y fecha 2026-10-10: la comparativa de la lista de compras tiene precios de Walmart desde el primer día; son precios del 10-oct-2026 y se marcan viejos el 9-nov-2026 (K22).
Impacto: el motor S5 queda fijado (S5 §7.2 `CompareCarts`, con el canon exacto). Queda PROVISIONAL la parte de la API: cómo `getBasketComparison` con `basket = expense` arma las canastas de las otras tiendas (S5 §7.2 «API» y §14 `CompareExpense`), la fuente de la semilla de Bodega (S5 §3) y si los precios del carrito de Walmart se siembran en `price` (S5 §3).
Cómo seguí: pregunta 1 con la opción A (sin cambiar C1); pregunta 2 con la opción A (`manual`, porque no consta un ticket); pregunta 3 con la opción A (solo entrada del motor y de la fixture). Las tres llevan la marca PROVISIONAL citando este bloqueo en `docs/contratos/motores/S5-interfaz.md` 1.1. Ninguna cifra canon depende de estas respuestas.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-47 · (1) Al comparar una compra registrada, ¿contra qué se compara en las otras tiendas?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Cada otra tienda se cotiza con sus precios vigentes que compiten (S5 §6) para los productos y cantidades de la compra base, eligiendo presentación como en S5 §7.3. No cambia C1.»
- **Consecuencia:** No cambia C1. La API **no** reproduce +138.50 con la semilla (Walmart cotizado con las cantidades de Bodega da otro total); el canon K-01 a K-08 se reproduce en el motor (`CompareCarts`) y en la fixture `s5-comparativa-canon`.
- **Qué cambia:** `S5-interfaz.md` §7.2 «API» y §14 `CompareExpense`: quita la marca PROVISIONAL de esta parte; C1 no cambia; la API no reproduce +$138.50 pero el motor sí (`CompareCarts`, fixture `s5-comparativa-canon`).
- **Estado:** RESUELTO

### Q2-49 · (3) ¿Los precios del carrito de Walmart del 10-oct-2026 se siembran en la app?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Solo entrada del motor y de la fixture: la app arranca sin precios de Walmart.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: el motor S5 queda fijado (S5 §7.2 `CompareCarts`, con el canon exacto). Queda PROVISIONAL la parte de la API: cómo `getBasketComparison` con `basket = expense` arma las canastas de las otras tiendas (S5 §7.2 «API» y §14 `CompareExpense`), la fuente de la semilla de Bodega (S5 §3) y si los precios del carrito de Walmart se siembran en `price` (S5 §3).
- **Qué cambia:** `S5-interfaz.md` §3 (semilla): la app arranca sin precios de Walmart; quita la marca PROVISIONAL de esta parte.
- **Estado:** RESUELTO

## Puerta 2 (11-oct-2026): sin respuesta
Fuente: `docs/estado/puerta2/RESPUESTAS.md` («Las preguntas sin recomendación quedan sin respuesta en esta puerta: PROVISIONAL con la opción que ya está en uso, y pasan a WF-3») y hoja de respuestas de `PUERTA-2.md` §8. Axel puede mandar una excepción por id cuando lea el documento.

### Q2-48 · (2) ¿Con qué fuente se siembran los precios de la compra de Bodega del 10-oct-2026?
- PROVISIONAL: sigue en uso A (`manual`), «porque no consta un ticket» (cita de este archivo, «Cómo seguí»). Solo Axel sabe si hubo ticket impreso (PUERTA-2 §3.16 y §7). Pasa a WF-3.
- **Estado:** ABIERTO

**Estado del bloqueo (puerta 2): ABIERTO** — Q2-47 y Q2-49 están resueltas; queda abierta Q2-48 (fuente de la semilla de Bodega).
