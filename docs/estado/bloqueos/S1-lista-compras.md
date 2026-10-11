# Bloqueo S1-lista-compras
Contexto: Corrección de S1 (ronda 1, hallazgo H-01 de `docs/estado/wf0/REVISION-r1.md`). La lista de compras de U1 necesita guardar marcas de comprado (N-151), artículos agregados a mano (N-152), la tienda o cantidad elegida por artículo (N-153) y «Usar esta tienda» (N-165), y el flujo F-3 (ruta A, paso 3) depende de esa última. Es la misma pregunta del bloqueo `C1-lista-compras`: una sola respuesta de Axel resuelve los dos.
Qué falta: ¿Las marcas de comprado y los artículos a mano se comparten entre los teléfonos de Axel y Ximena (se guardan en el servidor), viven solo en el dispositivo, o la lista es solo calculada?
Opciones:
- A. Compartidas en el servidor (tablas `shopping_list` y `shopping_list_item`). Consecuencia: los dos ven las mismas marcas y «Terminar compra» funciona desde cualquier teléfono; C1 agrega las 4 operaciones.
- B. Solo en el dispositivo. Consecuencia: las dos tablas salen del esquema; si Axel marca en su teléfono, Ximena no lo ve; U1 declara N-151 a N-153 y N-165 como «Local».
- C. Lista solo calculada, sin marcas ni artículos a mano. Consecuencia: las dos tablas salen del esquema y U1 quita N-151 a N-153 y N-165 (y el paso 3 de F-3 cambia).
Impacto: `shopping_list` y `shopping_list_item` (migración 00001), RN-15 de `reglas-negocio.md` y las operaciones de C1 para N-151, N-152, N-153 y N-165 quedan PROVISIONALES hasta la puerta 1. Ninguna cifra canon depende de esto.
Cómo seguí: Implementé la opción A marcada PROVISIONAL en la migración (`@entity` y `@desc` de las dos tablas), en el diccionario y en RN-15, para que la matriz U1 ↔ C1 pueda cerrarse. Si Axel elige B o C, se quitan las dos tablas y sus disparadores `updated_at` en la siguiente revisión del esquema.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción A: compartida en el servidor. Las dos tablas se quedan y dejan de ser PROVISIONALES.
- **Estado:** RESUELTO. Se aplica en: S1 1.2 (diccionario y reglas), C1 v1.1.
