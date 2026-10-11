# Bloqueo U1-captura-fuera-de-casa

Contexto: al definir el flujo F-1 (registrar gasto desde el teléfono en ≤ 3 toques, decisión G) y el estado «Sin conexión» de todas las pantallas (`docs/design/U1-arquitectura.md` §2.4 y G-07). El sistema solo se alcanza desde la red de casa (decisión H, pregunta 2 del plan). Es probable que muchos gastos se quieran registrar en la tienda o en la calle, donde la app instalada no llega al servidor.

Qué falta: ¿qué debe hacer la app instalada en el teléfono cuando Axel o Ximena quieren registrar un gasto (o un peso) fuera de casa, sin conexión con el servidor?

Opciones:
1. Sin captura fuera de casa. La app muestra «Sin conexión con el servidor de casa», deja ver lo último que cargó con su hora y deshabilita las capturas. Consecuencia: lo más simple; no agrega piezas; los gastos de la calle se capturan al volver a casa, con riesgo de olvidos.
2. Captura en cola en el teléfono. La hoja de gasto (y la de peso) funciona sin conexión, guarda lo capturado en el teléfono y lo envía al volver a la red de casa, con un aviso «3 gastos por enviar» y la opción de revisarlos antes. Consecuencia: el registro en ≤ 3 toques funciona también en la tienda; requiere almacenamiento local en la PWA, protección contra duplicados al enviar y una prueba en Android; ningún dato sale de la red de casa. Mientras no se envían, esos gastos no aparecen en las cifras.
3. Cola solo para gastos (como la opción 2, pero sin peso ni otros registros). Consecuencia: menos alcance y menos pruebas; cubre el caso más frecuente.

Impacto: el estado «Sin conexión» (G-07) y una variante de F-1 y F-6. El camino de 3 toques no cambia con ninguna opción. Con las opciones 2 o 3, C1 y S6 necesitan una forma de recibir varios gastos capturados antes (con su fecha y hora de captura) sin duplicarlos.

Cómo seguí: diseñé con la opción 1 marcada **PROVISIONAL** (§2.4 del documento). Ninguna cifra, esquema ni contrato depende de esto.
