# Revisión del catálogo base de productos · ronda 1 de 3

Revisor: `revisor-catalogo-r1` (independiente; no escribió nada de lo revisado y no corrige). Fecha: 11-oct-2026. Artefacto: `datos/semilla/productos.json` v1.0 (160 productos, 186 presentaciones, 9 grupos) con `datos/semilla/productos-README.md` y el reporte `docs/estado/wf1/catalogo-productos.md`. Se evalúa **solo** contra `docs/contratos/motores/productos-formato.md` 1.1 y `docs/contratos/motores/S4-interfaz.md`. Todas las verificaciones se volvieron a correr con scripts propios en `datos/tmp/rev-cat-r1/` (borrados al terminar); no se usó el validador del autor.

## Veredicto
**NO APROBADO**: 0 hallazgos altos, **1 medio**, 4 bajos. El archivo cumple todas las reglas de validación del formato (0 errores en el validador propio), las 39 presentaciones y los 39 renglones de `data.json` empatan exactamente y los 15 `gramos_por_unidad` coinciden con USDA. Lo que impide aprobar es que tres alimentos básicos que en México se compran por kilo quedaron en `pieza` sin forma de recibir un precio por kilo, y eso se dejó como «pendiente» para el siguiente implementador en lugar de como bloqueo (M-1).

## Criterios, estado y evidencia propia

| # | Criterio (fuente) | Estado | Evidencia del revisor |
|---|---|---|---|
| 1 | Claves de nivel superior y campos exactos de grupo, producto y presentación (formato §1, §2) | Cumple | `v1.py`: claves `version, fecha, grupos, productos`; conjuntos de campos exactos en 9 grupos, 160 productos y 186 presentaciones; 0 campos de más o de menos |
| 2 | Ids con el patrón y únicos; id de presentación = `<producto>.<slug>` (§2) | Cumple | 0 errores de patrón ni duplicados en 9 + 160 + 186 ids |
| 3 | Nombres únicos sin distinguir mayúsculas; nombre genérico es-MX, nunca marca (§2; S4-11) | Cumple | 0 duplicados con `casefold`; lectura de los 160 nombres: todos genéricos en español; las marcas solo están en presentaciones |
| 4 | Grupos iguales a §4 (id, nombre, canal) y `orden` entero | Cumple | Los 9 grupos coinciden uno a uno; orden 10…90 sin repetir |
| 5 | `categoria` existe; `unidad_base` ∈ {g, ml, pieza} | Cumple | 0 errores; unidades g 137 · ml 11 · pieza 12 |
| 6 | `gramos_por_unidad`: texto con ≤ 3 decimales, prohibido en `g`, con `fuente_gramos`; obligatorio en ml/pieza de comida (§2; S4 §9.2) | Cumple | 0 errores; 15 productos con dato, 0 alimentos en ml o pieza sin dato |
| 7 | Valores de `gramos_por_unidad` reproducibles desde la fuente citada | Cumple | Descarga propia de los 15 `fdc_id` (API de USDA FoodData Central): las 15 porciones coinciden (p. ej. huevo grande 50 g, tortilla de maíz 6" 26 g, tostada 12.3 g, aguacate sin cáscara 136 g, elote mediano «yields» 102 g); densidades recalculadas con Decimal y HALF_UP: 216/236.588 = 0.913, 218 → 0.921, 255 → 1.078, 275 → 1.162, 244 → 1.031, 238 → 1.006, iguales al archivo |
| 8 | `vertical` coherente con el canal del grupo (§3) | Cumple | 0 errores; mercado 73 · cadena 67 · ambas 20 (las 20 en grupos de cadena) |
| 9 | Sinónimos en minúsculas, sin repetidos (§2) | Cumple con observación | 0 errores de forma; 0 sinónimos compartidos entre productos; 3 sinónimos apuntan a otro producto (B-2) |
| 10 | `presentacion_tipica` pertenece al producto; ≥ 1 presentación (§2; S4 §3.2) | Cumple | 0 errores |
| 11 | `contenido_milli` entero > 0 y factor correcto (§2) | Cumple | 0 errores de tipo; factor contra el tamaño escrito en el nombre: 64 presentaciones con tamaño explícito fuera de `data.json`, 0 diferencias (p. ej. «cartón 12 piezas» = 12000, «lata 220 g» = 220000, «manojo (aprox. 50 g)» = 50000); las 83 «por kilo (granel)» valen 1000000 g |
| 12 | `granel` booleano; `codigo_barras` nulo o 8–14 dígitos; `marca` texto o nulo (§2) | Cumple | 0 errores; 186 códigos nulos (no se inventó ninguno) |
| 13 | Los 39 nombres de presentación de `data.json` existen tal cual, una sola vez (§2 validación, §5.1) | Cumple | 33 `items_identicos` + 3 `sustitutos.bodega` + 3 `sustitutos.walmart` = 39 nombres distintos, cada uno exactamente 1 vez |
| 14 | Los 39 renglones de `cuentas` en algún `nombres_ticket` (§2, §5.2); 61 piezas | Cumple | 39 renglones, 36 textos distintos, cada uno en exactamente 1 presentación; piezas 61 = `totales.articulos_bodega`; importe 3,094.50 = `totales.bodega`; las 39 filas de §5.2 (cuenta, texto, piezas, presentación) coinciden en orden; precio unitario de cada renglón = `bodega_unit` (o total ÷ cantidad en sustitutos) en 39 de 39, comparado en centavos enteros; piezas de Walmart 60 = `articulos_walmart` |
| 15 | Tabla §5.1 (36 productos, 39 presentaciones: nombre, categoría, unidad, vertical, ids, contenido) | Cumple | 36 filas leídas del formato por script; 0 diferencias |
| 16 | Sin marcas inventadas (tarea d) | Cumple | 30 marcas distintas, todas en presentaciones de `data.json` y contenidas en su nombre; 0 marcas en las 147 presentaciones genéricas |
| 17 | Todo producto con grupo y canal (S4 §6.5, §11; S4-12 «ingredientes sin canal = 0») | Cumple | 160 de 160 con `categoria` de un grupo con canal |
| 18 | Muestra de 20 con unidad base y presentación sensatas (tarea c) | Cumple 17 de 20 | Ver «Muestra»; limón, aguacate y tortilla de maíz no (M-1) |
| 19 | Cobertura de básicos de cocina casera mexicana (tarea e; formato §4) | Cumple con faltantes | Lista propia de 105 básicos: 89 presentes (84.8 %); faltantes en B-1. No es la cobertura de K25 (esa es de código); el recetario aún no existe, así que no se puede medir contra sus ingredientes |

### Muestra de 20 (elegida por el revisor)
huevo (pieza, cartón 12 = 12000, 50 g/pieza) ✓ · aceite de soya (ml, 946 ml, 0.921 g/ml) ✓ · aceite de oliva (ml, 500 ml) ✓ · jitomate saladet (g, por kilo) ✓ · cebolla blanca (g, por kilo) ✓ · ajo (g, por kilo) ✓ · cilantro (g, manojo 50 g típico + por kilo) ✓ · bolillo (pieza, por pieza) ✓ · leche (ml, 4 × 1 l) ✓ · arroz (g, 750 g) ✓ · frijol bayo (g, bolsa 1 kg + granel) ✓ · pechuga de pollo (g, por kilo) ✓ · carne molida de res (g, por kilo) ✓ · queso Oaxaca (g, por kilo + 400 g) ✓ · sal (g, 1 kg) ✓ · papel higiénico (pieza, 30 rollos = 30000) ✓ · elote en mazorca (pieza, 102 g de grano) ✓ · **limón (pieza, solo «por pieza»)** ✗ · **aguacate (pieza, solo «por pieza»)** ✗ · **tortillas de maíz (pieza, paquetes de 20 y 40)** ✗.

## Hallazgos

### M-1 (media) · Limón, aguacate y tortilla de maíz en `pieza` sin presentación por kilo; quedó como pendiente sin bloqueo
- **Ruta:** `datos/semilla/productos.json` (productos `limon`, `aguacate`, `tortilla-maiz`; relacionado: `huevo`, vertical `ambas`); `docs/estado/wf1/catalogo-productos.md` «Pendientes».
- **Qué pasa:** `contenido_milli` está en milésimas de la unidad base (formato §2), así que un producto en `pieza` no puede tener una presentación «por kilo». El limón y el aguacate se venden por kilo en cadenas y casi siempre en mercados; la tortilla de maíz se compra por kilo en la tortillería (el grupo es `ambas`); el huevo también se vende por kilo en el mercado. Con la base en `pieza`, un precio de $40/kg no entra a ninguna presentación, el costo por receta de S4 §8.2 (cantidad en unidad base × precio por unidad base) sale «sin precio» y la lista de §6.2 pide «paquetes» que no corresponden a cómo se compra. El autor lo vio y lo dejó como «Decisión del integrador/S5: añadir presentación en gramos con conversión por `gramos_por_unidad`, o capturar el precio por pieza»: una pregunta que el siguiente implementador tendría que adivinar, sin el archivo de bloqueo que pide el protocolo §2.
- **Corrección sugerida:** (A) dentro del formato actual y sin decisión nueva: pasar `limon`, `aguacate` y `tortilla-maiz` a `unidad_base = g` con presentación «por kilo (granel)» = 1000000 (la tortilla puede conservar los paquetes expresados en gramos nominales), quitar su `gramos_por_unidad` (prohibido en `g`) y dejar que la receta diga «2 limones» en `texto`; o (B) conservarlos en `pieza` y escribir `docs/estado/bloqueos/WF1-catalogo-pieza-por-kilo.md` con las opciones (conversión kilo↔pieza con `gramos_por_unidad` en S5, captura del precio por pieza, cambio de unidad base) y su impacto en S4 §8 y en los vínculos de PROFECO, marcando esos productos como PROVISIONAL en el README.

### B-1 (baja) · Básicos de cocina casera que faltan
- **Ruta:** `datos/semilla/productos.json`.
- **Qué pasa:** de una lista propia de 105 básicos faltan 16. Los más evidentes para guisos caseros: **chícharo**, **manteca de cerdo**, **clavo de olor**, **pimienta gorda**, **rábano**, **mole en pasta** y **mostaza**. Después: haba seca, tomillo, leche condensada, pasas, cacahuate, ajonjolí, salsa picante embotellada, gelatina y harina para hot cakes (el autor descartó a propósito rábano, pasas, cacahuate y ajonjolí).
- **Corrección sugerida:** agregar al menos los 7 primeros con el mismo formato (ids nuevos; los existentes no cambian), o dejar escrito en el README que el agente del recetario puede agregar productos a este archivo con estas reglas.

### B-2 (baja) · Tres sinónimos que empatan con otro producto
- **Ruta:** `datos/semilla/productos.json`: `frijol-bayo` («frijol canario»), `bolillo` («telera»), `petit-suisse` («queso fresco para niños»).
- **Qué pasa:** el frijol canario (peruano) es otra variedad, la telera es otro pan y «queso fresco para niños» choca en la búsqueda con `queso-fresco`. Como los sinónimos también sirven para empatar renglones de ticket y OCR (formato §2), un renglón «Frijol canario 1 kg» se compararía como frijol bayo.
- **Corrección sugerida:** quitar esos tres sinónimos (o crear productos propios si se quieren).

### B-3 (baja) · Id `tomate-bola` con nombre «Jitomate bola»
- **Ruta:** `datos/semilla/productos.json`, producto `tomate-bola` (su vecino es `jitomate-saladet`).
- **Qué pasa:** el id es estable y lo citarán el recetario, la tabla de nutrición y los vínculos de PROFECO; en México «tomate» a secas suele ser el verde (`tomate-verde`), así que el id se presta a confusión.
- **Corrección sugerida:** renombrarlo a `jitomate-bola` ahora, antes de que otro archivo lo cite.

### B-4 (baja) · Pesos por pieza de productos de EE. UU. sin un paso de confirmación
- **Ruta:** `datos/semilla/productos.json`: `limon` (67 g, lima persa de 2"), `bolillo` (57 g, panecillo kaiser), `caldo-pollo-cubos` (4 g, cubo de EE. UU.).
- **Qué pasa:** los valores son fieles a USDA (verificado) y el `fuente_gramos` dice «pendiente cotejar», pero `gramos_por_unidad` no tiene estado `propuesto`/`confirmado` como `tabla.csv` (S4 §9.1), así que entrarán a la base tal cual y moverán las kcal de las recetas que los usen.
- **Corrección sugerida:** incluir estos tres en la lista que Axel confirma en la puerta 2 (la misma del bloqueo `S4-confirmacion-empates`) o cotejarlos con la etiqueta o la báscula antes de sembrar.

## Notas sin severidad
- La tarea pedía «los 61 nombres»: en `data.json` hay 39 nombres de presentación y 36 textos distintos en 39 renglones de `cuentas`; 61 es el número de piezas. Se verificaron los tres conteos (criterios 13 y 14).
- El validador y el generador del autor solo existen en el scratchpad de la sesión; su solicitud 2 (`scripts/validar_productos.py`) haría reproducible la verificación.
- Sinónimos con marcas comunes («nescafé», «maizena», «royal», «queso filadelfia», «cloralex»): el formato solo prohíbe marcas en `.nombre`; ayudan a empatar tickets.

## Estado de rondas anteriores
Primera ronda: no hay hallazgos previos.
