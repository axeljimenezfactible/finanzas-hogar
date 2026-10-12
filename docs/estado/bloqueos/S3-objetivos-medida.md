# Bloqueo S3-objetivos-medida
Contexto: agente contrato-S3 (WF-1) al fijar el comparador (`S3-interfaz.md` §8). K16 y K17 deciden las 3 restricciones, la medida «desviación estándar de los pagos (y pico contra el tope)» y «unidades base de la canasta por peso gastado con los precios vigentes de cada tienda», pero dejan cuatro lecturas abiertas que cambian qué escenario gana. C1 v1.2 tampoco las fija (`Scenario` no trae `store_id` ni la desviación estándar; `ScenarioRequest` no trae canasta).
Qué falta:
  (1) Canasta y tienda (objetivo 4): ¿la tienda es una dimensión más del escenario (cada tienda con precios vigentes) y qué es «la canasta» de una compra planeada? (a) la lista de compras de la semana (P-13) con los precios vigentes de S5 en cada tienda, y lo financiado es lo que cuesta en esa tienda; (b) el monto de la petición es fijo y compra amount × unidades / costo de la canasta en cada tienda; (c) la tienda no es dimensión y el objetivo 4 solo distingue promociones (unidades fijas / total pagado).
  (2) Serie del objetivo 2: ¿desviación estándar de los pagos de la compra del escenario sola, o de los totales por quincena de todos los créditos del hogar (existentes + escenario) entre el primero y el último pago del escenario? ¿Y el pico contra el tope por persona es solo marca (S3-12: «se acepta y se marca») o entra a la medida?
  (3) Horizonte del colchón y del «máximo disponible»: ¿desde la fecha de compra del escenario hasta su último pago, o hasta el último pago de cualquier compra viva de esa línea?
  (4) Créditos sin límite (préstamo y crédito de tienda sin `credit_limit_cents`): ¿entran al comparador? Si entran, el colchón y el disponible no aplican.
Opciones (propuesta del agente marcada con *):
  (1) *(b): el monto de la petición es fijo (lo financiado es igual en todas las tiendas, como dice C1 `ScenarioRequest.amount_cents`); la tienda es dimensión solo cuando la petición no fija `store_id` ni `merchant_name`; S5 entrega por tienda las unidades base U y el costo C de la canasta con precios vigentes; rendimiento = (monto × U / C) / (monto + intereses). Alternativas (a) y (c). Consecuencia de (b): `Scenario` necesita `store_id` (solicitud a C1); la canasta la arma la API con S5.
  (2) *totales por quincena del hogar (calendario unificado) en la ventana del escenario; el pico contra el tope solo marca (`over_cap`) y desempata; alternativa: sumar el exceso sobre el tope a la medida.
  (3) *la vida de la compra del escenario (de su fecha a su último pago, inclusive): fuera de esa ventana el escenario no cambia el disponible, así que una baja previa o posterior no lo descarta; alternativa: hasta el último pago de cualquier compra viva de la línea.
  (4) *solo entran créditos con límite (línea revolvente, tarjeta y tienda con límite); los demás quedan fuera con motivo «sin límite»; alternativa: entran y el colchón no aplica.
Impacto: objetivos 2, 3 y 4, restricción de colchón, `S3-escenarios.md` y los 3 casos `s3-comparador-caso-1..3` (S3-12, S3-13, S3-14). Ninguna cifra canon depende de esto.
Cómo seguí: la interfaz implementa las propuestas * como PROVISIONAL citando este bloqueo, con la entrada de canasta como parámetro (`[]BasketQuote` de S5) para que cualquier opción quepa sin cambiar firmas; los 3 casos se diseñan para que los 4 óptimos sean distintos bajo la propuesta.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-33 · (1) ¿La tienda es una dimensión del escenario del comparador y qué es «la canasta»?
- **Respuesta:** **(b)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «*(b): el monto de la petición es fijo (lo financiado es igual en todas las tiendas, como dice C1 `ScenarioRequest.amount_cents`); la tienda es dimensión solo cuando la petición no fija `store_id` ni `merchant_name`; S5 entrega por tienda las unidades base U y el costo C de la canasta con precios vigentes; rendimiento = (monto × U / C) / (monto + intereses).»
- **Consecuencia:** `Scenario` necesita `store_id` (solicitud a C1); la canasta la arma la API con S5.
- **Qué cambia:** `Scenario` necesita `store_id` (solicitud a C1 v1.3); S5 entrega por tienda las unidades base U y el costo C de la canasta; objetivo 4 del comparador, `docs/estado/S3-escenarios.md` y los casos `s3-comparador-caso-1..3` (S3-12, S3-13, S3-14).
- **Estado:** RESUELTO

### Q2-34 · (2) ¿Sobre qué serie se mide la desviación de pagos y qué peso tiene el pico contra el tope?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «*totales por quincena del hogar (calendario unificado) en la ventana del escenario; el pico contra el tope solo marca (`over_cap`) y desempata»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: objetivos 2, 3 y 4, restricción de colchón, `S3-escenarios.md` y los 3 casos `s3-comparador-caso-1..3` (S3-12, S3-13, S3-14). Ninguna cifra canon depende de esto.
- **Qué cambia:** Objetivo 2 del comparador: `over_cap` solo marca y desempata; `S3-escenarios.md` y casos `s3-comparador-caso-1..3`.
- **Estado:** RESUELTO

### Q2-35 · (3) ¿Horizonte del colchón y del disponible?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «*la vida de la compra del escenario (de su fecha a su último pago, inclusive): fuera de esa ventana el escenario no cambia el disponible, así que una baja previa o posterior no lo descarta»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: objetivos 2, 3 y 4, restricción de colchón, `S3-escenarios.md` y los 3 casos `s3-comparador-caso-1..3` (S3-12, S3-13, S3-14). Ninguna cifra canon depende de esto.
- **Qué cambia:** Restricción de colchón y objetivo 3; `S3-escenarios.md` y casos del comparador.
- **Estado:** RESUELTO

### Q2-36 · (4) ¿Créditos sin límite entran al comparador?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «*solo entran créditos con límite (línea revolvente, tarjeta y tienda con límite); los demás quedan fuera con motivo «sin límite»»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: objetivos 2, 3 y 4, restricción de colchón, `S3-escenarios.md` y los 3 casos `s3-comparador-caso-1..3` (S3-12, S3-13, S3-14). Ninguna cifra canon depende de esto.
- **Qué cambia:** Los créditos sin límite quedan fuera con el motivo «sin límite»; `S3-escenarios.md` y casos del comparador.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
