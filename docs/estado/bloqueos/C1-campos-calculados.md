# Bloqueo C1-campos-calculados
Contexto: Al definir en C1 los campos calculados que entregan los motores S3, S4 y S5. La API los expone, pero su definición exacta no está en el plan ni en las decisiones; la deciden los motores en WF-1.
Qué falta: confirmar cinco definiciones.
1. Clasificación del IMC (`bmi_class`, N-230): ¿la de la OMS (bajo peso < 18.5, normal 18.5 a 24.9, sobrepeso 25 a 29.9, obesidad ≥ 30)?
2. Precio por unidad base (`base_unit_price_cents`, N-162, N-170, N-171): ¿centavos por 1 kg, 1 l o 1 pieza, con ROUND_HALF_UP? (El dinero nunca va en punto flotante y el precio por gramo tiene fracciones de centavo.)
3. Membresía de Sam's (N-164): ¿con qué frecuencia de compra se anualiza el ahorro y cómo se expresa el punto de equilibrio (compras al año necesarias para cubrir el costo anual)?
4. Rendimiento de compra por peso del comparador (`value_per_peso_milli`) y puntaje ponderado (`score_bp`, de 0 a 10000): dependen también del bloqueo existente `criterios-optimizador`.
5. Estado de inventario (suficiente, reponer, agotado; N-180): ¿reponer = existencia estimada ≤ `reorder_level_milli` y agotado = 0?
Opciones:
- A (la que tomé, PROVISIONAL). Las cinco definiciones tal como están escritas arriba. Consecuencia: los nombres y tipos de los campos quedan fijos; solo su cálculo puede cambiar.
- B. Que cada motor (S3, S4, S5) proponga su definición en `reglas-negocio.md` y Axel la apruebe en la puerta 2. Consecuencia: C1 v1.1 ajusta descripciones; los nombres de campo no cambian.
- C. Quitar de C1 los campos que no se aprueben. Consecuencia: U1 pierde esas cifras en pantalla.
Impacto: Campos `bmi_class`, `base_unit_price_cents`, `membership.*`, `value_per_peso_milli`, `score_bp` y `status` de inventario en `openapi.yaml`. Ninguna cifra canon del oráculo depende de esto.
Cómo seguí: Los campos van con `x-provisional` y la referencia a este bloqueo; las fórmulas FIRMES (IMC = peso ÷ estatura², promedios de kcal, sumas de dinero) no se marcaron.
