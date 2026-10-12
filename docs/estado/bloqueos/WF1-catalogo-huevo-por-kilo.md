# Bloqueo WF1-catalogo-huevo-por-kilo
Contexto: ronda `correccion-catalogo-r1` del catálogo (`datos/semilla/productos.json`). Limón, aguacate y tortilla de maíz se pasaron a `g` con presentación por kilo (hallazgo M-1). El huevo (`huevo`, unidad base `pieza`, 50 g por pieza, vertical `ambas`) tiene el mismo problema en menor grado: en mercado y tianguis también se vende por kilo, pero en cadena se vende por pieza (cartón 12, cartón 18, plato 30). Pasarlo a `g` quitaría la comparación natural por pieza y rompería el cálculo de huevos por receta («2 huevos»); dejarlo en `pieza` impide registrar un precio por kilo.
Qué falta: ¿Cómo debe registrarse un precio por kilo de un producto cuya unidad base es `pieza` (huevo)? Es decisión de modelo (S5/Ajustes), no del catálogo.
Opciones:
1. Dejar `huevo` en `pieza` y que Axel capture siempre el precio por pieza o por cartón (al comprar por kilo, divide). Consecuencia: cero cambios; el precio por kilo del mercado no se captura tal cual.
2. S5 convierte kilo a pieza con `gramos_por_unidad` (50 g): una presentación «por kilo» tendría `contenido_milli` en gramos y requeriría un campo nuevo de conversión. Consecuencia: cambia `productos-formato.md` §2 y el esquema; el huevo de 50 g es un promedio USDA, no pesa igual el huevo del mercado.
3. Pasar `huevo` a `g` con presentación por kilo y cartones en gramos nominales (cartón 12 ≈ 600 g). Consecuencia: costo por receta correcto con gramos, pero las recetas deben expresarse en gramos y los cartones pierden su conteo exacto de piezas.
4. Dos productos: `huevo` (pieza) y `huevo-granel` (g). Consecuencia: duplica el producto y su comparación de precios.
Impacto: S4 §8.2 (costo por receta) calcula con precio por pieza; un precio por kilo del mercado queda sin captura. PROFECO: el vínculo del huevo debe decidir si compara por pieza o por kilo. El catálogo sigue siendo válido; nada se bloquea.
Cómo seguí: dejé `huevo` sin cambios con sus tres presentaciones por pieza (la forma en que se compra en cadena, que es la fuente dominante del hogar) y lo documenté en el README. Ningún valor provisional entra en cifras ni esquema.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-71 · ¿Cómo se registra un precio por kilo del huevo, cuya unidad base es la pieza?
- **Respuesta:** **1** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Dejar `huevo` en `pieza` y que Axel capture siempre el precio por pieza o por cartón (al comprar por kilo, divide).»
- **Consecuencia:** cero cambios; el precio por kilo del mercado no se captura tal cual.
- **Qué cambia:** Ningún archivo cambia: el catálogo (`datos/semilla/productos.json`) conserva `huevo` en `pieza` con sus tres presentaciones; el README del catálogo ya lo documenta.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
