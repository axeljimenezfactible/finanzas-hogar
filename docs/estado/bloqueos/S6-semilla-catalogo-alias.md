# Bloqueo S6-semilla-catalogo-alias

Contexto: WF-3a, agente `semilla`. `datos/semilla/productos.json` trae por producto `sinonimos`, `presentacion_tipica` y `vertical`, y por presentación `nombres_ticket`. El esquema 1.4 no tiene columna para ninguno (`product` y `presentation` solo guardan nombre, grupo, unidad, gramos por unidad, marca, contenido, granel y código de barras).

Qué falta: ¿dónde viven los sinónimos y los nombres de ticket una vez sembrado el catálogo? El empate de renglones de ticket y la búsqueda (S5, OCR) los necesitan.

Opciones:
- A: el empate lee `productos.json` en cada arranque (los alias son de solo lectura y no se editan en Ajustes). Consecuencia: sin migración; lo que Axel agregue en Ajustes no tiene alias.
- B: migración 00005 con tablas `product_alias` y `presentation_ticket_name` y carga desde la semilla. Consecuencia: editables en Ajustes; cambia esquema, C1 y diccionario.
- C: no usarlos. Consecuencia: el OCR empata solo por nombre exacto.

Impacto: la semilla carga todo lo que el esquema permite (160 productos, 186 presentaciones, grupos, gramos por unidad). `nombres_ticket` SÍ se usa al sembrar la compra del 10-oct para ligar los 39 renglones de `data.json: cuentas` con su presentación. `presentacion_tipica` no se guarda: las recetas ya traen su `presentacion_id`.

Cómo seguí: opción A PROVISIONAL (nada en la base). Para B hace falta decisión de Axel y migración nueva (nunca editar una aplicada).

- **Estado:** ABIERTO
