# Bloqueo S5-reglas-comparativa
Contexto: al escribir el contrato del motor S5 (`docs/contratos/motores/S5-interfaz.md`, WF-1, agente `contrato-S4S5`, 11-oct-2026). K22/K40 fijan el **valor** del peso de un precio viejo, pero no cómo se aplica al elegir la mejor tienda; S5-03 pide una regla de desempate «escrita en `reglas-negocio.md`» que no existe; S5-10 pide vigencia de promociones y el esquema no tiene dónde guardarla; y la membresía de Sam's (K7, `C1-campos-calculados` pregunta 3) dejó sin decidir la frecuencia con que se anualiza el ahorro.
Qué falta (cuatro preguntas para Axel):
1. ¿Cómo «pesa menos» un precio viejo al elegir la tienda más barata?
2. ¿Qué pasa en un empate (por ingrediente y por canasta) y puede ganar una tienda a la que le faltan precios de la canasta?
3. ¿Cuánto dura vigente un precio de promoción (`price.is_promotional = 1`)?
4. (a) ¿Con cuántas compras al año se anualiza el ahorro de Sam's? (b) Si costo ÷ ahorro no es entero, ¿se redondea hacia arriba?
Opciones:
- 1. A (PROVISIONAL): ordenar por precio por unidad base ÷ peso; totales y pantallas con el precio real. B: el peso solo desempata (gana el más barato aunque sea viejo). C: un precio viejo solo gana si no hay precio fresco en otra tienda.
- 2. A (PROVISIONAL): empate = precios ajustados iguales como fracción; el renglón queda sin ganador y suma en «empates»; la tienda recomendada es la que ya tiene más renglones de la lista asignados, luego la de menos precios viejos, luego la de menor id; una tienda con renglones sin precio no puede ser «la más barata». B: en empate gana la tienda preferida que Axel marque en Ajustes (campo nuevo). C: en empate gana la más cercana al día de despensa (mercado en su día). Precisión de la corrección r3 (11-oct-2026, agente `correccion-contrato-S4S5-r3`, hallazgo M-r3-2): dentro de A, la tienda recomendada de un **renglón** ya no depende de la lista, porque `QuoteFor` no la recibe. `S5-interfaz.md` 1.3 §14.1 desempata con el menor precio del paquete, luego el orden de `stores` y luego el del catálogo (menor id en la API). «Más renglones, menos precios viejos, menor id» queda solo para ordenar tiendas que empatan en el total de la canasta (§7.4).
- 3. A (PROVISIONAL): vigente solo en su `observed_date`; RN-12 ya la excluye de la comparativa y solo se muestra. B: columna nueva `price.valid_to` (cambio de esquema y C1). C: ligar a la tabla `promotion` (hoy es de promociones de crédito).
- 4a. A: clave nueva en `app_setting` (p. ej. `price.membership_purchases_per_year`) que Axel captura en Tiendas y precios (K44). B (PROVISIONAL): las compras de despensa en cadenas registradas en los últimos 365 días. C: 52 (una lista por semana). D: 6 (bimestral, el ejemplo de RN-13).
- 4b. A (PROVISIONAL): hacia arriba (8.3 compras → 9). B: mostrar con un decimal (cambia el tipo entero de C1).
Impacto: ganador por ingrediente con precios viejos, tienda recomendada en empates, promociones mostradas y los campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases` cuando la división no es exacta. Ninguna cifra canon cambia: la canasta del 10-oct tiene precios del mismo día y Bodega gana sin empate; el ejemplo de S5-09 (1,000.00 ÷ 125.00 = 8) es exacto.
Cómo seguí: el contrato S5 usa las opciones marcadas PROVISIONAL con cita a este bloqueo (§6, §7.1, §7.4, §8 y, desde la versión 1.3, §14.1).

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-59 · (1) ¿Cómo «pesa menos» un precio viejo al elegir la tienda más barata?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «ordenar por precio por unidad base ÷ peso; totales y pantallas con el precio real.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: ganador por ingrediente con precios viejos, tienda recomendada en empates, promociones mostradas y los campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases` cuando la división no es exacta. Ninguna cifra canon cambia: la canasta del 10-oct tiene precios del mismo día y Bodega gana sin empate; el ejemplo de S5-09 (1,000.00 ÷ 125.00 = 8) es exacto.
- **Qué cambia:** `S5-interfaz.md` §6, §7.1, §7.4, §8 y §14.1 quitan la marca PROVISIONAL de esta parte; ninguna cifra canon cambia.
- **Estado:** RESUELTO

### Q2-60 · (2) Empates, y si puede ganar una tienda a la que le faltan precios
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «empate = precios ajustados iguales como fracción; el renglón queda sin ganador y suma en «empates»; la tienda recomendada es la que ya tiene más renglones de la lista asignados, luego la de menos precios viejos, luego la de menor id; una tienda con renglones sin precio no puede ser «la más barata».»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: ganador por ingrediente con precios viejos, tienda recomendada en empates, promociones mostradas y los campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases` cuando la división no es exacta. Ninguna cifra canon cambia: la canasta del 10-oct tiene precios del mismo día y Bodega gana sin empate; el ejemplo de S5-09 (1,000.00 ÷ 125.00 = 8) es exacto.
- **Qué cambia:** `S5-interfaz.md` §6, §7.1, §7.4, §8 y §14.1 quitan la marca PROVISIONAL de esta parte; ninguna cifra canon cambia. Incluye la precisión de la corrección r3 (§14.1: en un renglón se desempata por el menor precio del paquete; «más renglones, menos precios viejos, menor id» solo para tiendas que empatan en el total, §7.4). RN de desempate en `reglas-negocio.md` (S5-03).
- **Estado:** RESUELTO

### Q2-61 · (3) ¿Cuánto dura vigente un precio de promoción?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «vigente solo en su `observed_date`; RN-12 ya la excluye de la comparativa y solo se muestra.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: ganador por ingrediente con precios viejos, tienda recomendada en empates, promociones mostradas y los campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases` cuando la división no es exacta. Ninguna cifra canon cambia: la canasta del 10-oct tiene precios del mismo día y Bodega gana sin empate; el ejemplo de S5-09 (1,000.00 ÷ 125.00 = 8) es exacto.
- **Qué cambia:** `S5-interfaz.md` §6, §7.1, §7.4, §8 y §14.1 quitan la marca PROVISIONAL de esta parte; ninguna cifra canon cambia. Sin columna nueva (esquema sin cambio).
- **Estado:** RESUELTO

### Q2-62 · (4a) ¿Con cuántas compras al año se anualiza el ahorro de la membresía de Sam's?
- **Respuesta:** **B** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «las compras de despensa en cadenas registradas en los últimos 365 días.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: ganador por ingrediente con precios viejos, tienda recomendada en empates, promociones mostradas y los campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases` cuando la división no es exacta. Ninguna cifra canon cambia: la canasta del 10-oct tiene precios del mismo día y Bodega gana sin empate; el ejemplo de S5-09 (1,000.00 ÷ 125.00 = 8) es exacto.
- **Qué cambia:** `S5-interfaz.md` §6, §7.1, §7.4, §8 y §14.1 quitan la marca PROVISIONAL de esta parte; ninguna cifra canon cambia. Campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases`. Nota: el archivo marca B como la opción vigente (PROVISIONAL) y PUERTA-2 la recomienda; Axel la aprueba.
- **Estado:** RESUELTO

### Q2-63 · (4b) Si costo ÷ ahorro no da entero, ¿se redondea hacia arriba?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «hacia arriba (8.3 compras → 9).»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: ganador por ingrediente con precios viejos, tienda recomendada en empates, promociones mostradas y los campos `membership.annual_savings_cents`, `verdict` y `break_even_purchases` cuando la división no es exacta. Ninguna cifra canon cambia: la canasta del 10-oct tiene precios del mismo día y Bodega gana sin empate; el ejemplo de S5-09 (1,000.00 ÷ 125.00 = 8) es exacto.
- **Qué cambia:** `S5-interfaz.md` §6, §7.1, §7.4, §8 y §14.1 quitan la marca PROVISIONAL de esta parte; ninguna cifra canon cambia.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
