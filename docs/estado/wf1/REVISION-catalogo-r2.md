# Revisión del catálogo base de productos · ronda 2 de 3

Revisor: `revisor-catalogo-r2` (independiente; no escribió nada de lo revisado y no corrige). Fecha: 11-oct-2026. Artefacto: `datos/semilla/productos.json` (campo `version` = «1.0»; 160 productos, 186 presentaciones, 9 grupos) tras `correccion-catalogo-r1`, con `datos/semilla/productos-README.md` 1.1, el bloqueo `docs/estado/bloqueos/WF1-catalogo-huevo-por-kilo.md` y los reportes `docs/estado/wf1/catalogo-productos.md` y `correccion-catalogo-r1.md`. Se evalúa **solo** contra `docs/contratos/motores/productos-formato.md` 1.1 y `docs/contratos/motores/S4-interfaz.md`. Todas las verificaciones se corrieron con scripts propios en `datos/tmp/rev-cat-r2/` (borrados al terminar); no se usó el validador del autor ni el de la ronda 1.

## Veredicto
**APROBADO**: 0 hallazgos altos, 0 medios, 7 bajos (4 siguen abiertos desde la ronda 1 y 3 son nuevos). El M-1 de la ronda 1 está **cerrado**: limón, aguacate y tortilla de maíz pasaron a `g` con presentación «por kilo (granel)», y el huevo quedó en `pieza` con un bloqueo en el formato del protocolo §2. El validador propio da 0 errores, los 39 nombres de presentación y los 39 renglones de `data.json` empatan exactamente, y los 12 `gramos_por_unidad` coinciden con USDA (descarga propia).

## Criterios, estado y evidencia propia

| # | Criterio (fuente) | Estado | Evidencia del revisor |
|---|---|---|---|
| 1 | Claves de nivel superior y campos exactos de grupo, producto y presentación (formato §1, §2) | Cumple | `v2.py`: claves `version, fecha, grupos, productos`; conjuntos de campos exactos en 9 grupos, 160 productos y 186 presentaciones |
| 2 | Ids con el patrón y únicos; id de presentación = `<producto>.<slug>` (§2) | Cumple | 0 errores de patrón ni duplicados en 9 + 160 + 186 ids |
| 3 | Nombres únicos sin distinguir mayúsculas; genéricos en es-MX, nunca marca (§2) | Cumple | 0 duplicados con `casefold`; lectura de los 160 nombres: todos genéricos en español; el único «nombre con marca» del cruce automático es `lechuga-romana` (contiene «roma»), un falso positivo |
| 4 | Grupos iguales a §4 (id, nombre, canal), `orden` entero sin repetir | Cumple | Los 9 grupos coinciden uno a uno con §4 |
| 5 | `categoria` existe; `unidad_base` ∈ {g, ml, pieza} | Cumple | 0 errores; g 140 · ml 11 · pieza 9 |
| 6 | `gramos_por_unidad`: texto con ≤ 3 decimales, > 0, prohibido en `g`, con `fuente_gramos`; obligatorio en ml/pieza de comida (§2; S4 §9.2) | Cumple | 0 errores; 12 productos con dato; 0 alimentos en `ml` o `pieza` sin dato; limón, aguacate y tortilla de maíz con `gramos_por_unidad` y `fuente_gramos` en `null` |
| 7 | `gramos_por_unidad` reproducible desde la fuente citada | Cumple | API de USDA FoodData Central (12 `fdc_id`, descarga propia y borrada): aceite de oliva 216 g/taza → 216/236.588 = 0.913; soya 218 → 0.921; shoyu 255 → 1.078; inglesa 275 → 1.162; leche 244 → 1.031; vinagre 238 → 1.006 (Decimal, HALF_UP a 3 decimales); cubo 4 g; elote mediano «yields» 102 g; huevo grande 50 g; panecillo 3-1/2" 57 g; tortilla de harina 6" 32 g; tostada 12.3 g. 12 de 12 iguales |
| 8 | `vertical` coherente con el canal del grupo (§3) | Cumple | 0 errores; mercado 73 · cadena 67 · ambas 20 (todas en grupos de canal cadena) |
| 9 | Sinónimos en minúsculas, sin repetidos ni compartidos; ningún sinónimo igual al nombre de otro producto (§2) | Cumple con observación | 0 errores de forma; siguen 3 sinónimos que empatan semánticamente con otro producto (B-2) |
| 10 | `presentacion_tipica` pertenece al producto; ≥ 1 presentación (§2) | Cumple | 0 errores; la típica de limón, aguacate y tortilla es `*.kilo-granel` |
| 11 | `contenido_milli` entero > 0 y factor correcto (§2) | Cumple | Factor recalculado desde el tamaño escrito en el nombre en las 186 presentaciones: 1 diferencia aparente (`leche.santa-clara-4x1l-taza`, «4 pack 1 l»: mi expresión regular leyó 1 l; 4 × 1 l = 4000000 es correcto); 0 diferencias reales |
| 12 | `granel` booleano; `codigo_barras` nulo u 8–14 dígitos; `marca` texto o nulo; `nombres_ticket` lista de textos (§2) | Cumple | 0 errores; «por kilo» y «por pieza» con `granel: true`, empaques con `false` (única excepción: «Uva verde por pieza de 1 kg aprox», nombre literal de `data.json`); 186 códigos nulos |
| 13 | Los 39 nombres de presentación de `data.json` existen tal cual, una sola vez (§2, §5.1) | Cumple | 33 + 3 + 3 = 39 nombres distintos, cada uno exactamente 1 vez |
| 14 | Los 39 renglones de `cuentas` en algún `nombres_ticket` (§2, §5.2); 61 piezas | Cumple | 39 renglones, 36 textos distintos, cada uno en exactamente 1 presentación; 0 `nombres_ticket` fuera de `cuentas`; piezas 61 = `totales.articulos_bodega`; importe 3,094.50 = `totales.bodega`; precio unitario por renglón en centavos enteros = `bodega_unit` (o total ÷ cantidad en sustitutos) en 39 de 39; piezas por presentación = `cantidad` (atún 6 = 5 + la sexta lata); las 73 cadenas `producto`/`bodega`/`walmart` de todo `data.json` tienen presentación o alias |
| 15 | Tablas §5.1 (36 filas) y §5.2 (39 filas) del formato | Cumple | Leídas del formato por script: 0 diferencias en nombre, categoría, unidad, vertical, ids, contenido, cuenta, texto y piezas |
| 16 | Sin marcas inventadas (tarea d) | Cumple | 30 marcas en 38 presentaciones, todas de `data.json` y contenidas en su nombre; 0 marcas en las 147 presentaciones genéricas |
| 17 | Todo producto con grupo y canal (S4 §6.5, §11; S4-12) | Cumple | 160 de 160 |
| 18 | Muestra de 20 con unidad base y presentación sensatas (tarea c) | Cumple 20 de 20 | Ver «Muestra» |
| 19 | Cobertura de básicos de cocina casera mexicana (tarea e) | Cumple con faltantes | Lista propia de 157 básicos: 131 presentes (83.4 %); faltantes en B-1. No es la cobertura de K25 (esa es de código); el recetario todavía no existe |
| 20 | Nada cita los ids de presentación que se quitaron | Cumple | `grep` en todo el repositorio sin `.git`: `limon.pieza`, `aguacate.pieza`, `tortilla-maiz.paquete-20/40` solo aparecen como historia en el README y en `correccion-catalogo-r1.md`; ningún contrato de motores, fixture ni `profeco/FORMATO.md` cita limón, aguacate, tortilla de maíz ni huevo |

### Muestra de 20 (elegida por el revisor)
huevo (pieza, cartón 12 = 12000, plato 30 = 30000; el kilo de mercado va al bloqueo) ✓ · limón (g, por kilo = 1000000) ✓ · aguacate (g, por kilo) ✓ · tortillas de maíz (g, por kilo + paquete 500 g) ✓ · aceite de soya (ml, 946 ml) ✓ · vinagre blanco (ml, 1 l) ✓ · leche (ml, 4 × 1 l y 6 × 1 l) ✓ · jitomate saladet (g, por kilo) ✓ · tomate verde (g, por kilo) ✓ · chile serrano (g, por kilo) ✓ · cilantro (g, manojo 50 g + kilo) ✓ · plátano (g, por kilo) ✓ · elote en mazorca (pieza, por pieza, 102 g de grano) ✓ · bolillo (pieza, por pieza) ✓ · tostadas (pieza, paquete 30) ✓ · queso fresco (g, por kilo) ✓ · crema (g, bote 200 g, contenido nominal) ✓ · pechuga de pollo (g, por kilo) ✓ · frijol negro (g, bolsa 1 kg + granel) ✓ · caldo de pollo en cubos (pieza, 24 cubos; peso por cubo en B-4) ✓.

## Hallazgos

### B-1 (baja, sigue abierto desde r1) · Básicos de cocina casera que faltan
- **Ruta:** `datos/semilla/productos.json`.
- **Qué pasa:** de una lista propia de 157 básicos faltan 26. Los evidentes para guisos caseros: **chícharo**, **manteca de cerdo**, **clavo de olor**, **pimienta gorda**, **achiote**, **mole en pasta**, **rábano**, **mostaza**, **salsa picante embotellada**, **leche condensada**, **media crema** y **camarón**. Después: haba seca, tomillo, chile morita, ajonjolí, cacahuate, pasas, gelatina, miel, mermelada, galletas María, té de manzanilla, cereal, harina para hot cakes e hígado.
- **Corrección sugerida:** agregar al menos los 12 primeros con el mismo formato (ids nuevos; los existentes no cambian). El reporte del autor ya dice que el recetario puede ampliar el archivo; conviene que lo diga también el README.

### B-2 (baja, sigue abierto desde r1) · Tres sinónimos que empatan con otro producto
- **Ruta:** `datos/semilla/productos.json`: `frijol-bayo` («frijol canario»), `bolillo` («telera»), `petit-suisse` («queso fresco para niños»).
- **Corrección sugerida:** quitarlos o crear productos propios; un renglón de ticket «Frijol canario 1 kg» se compararía como frijol bayo.

### B-3 (baja, sigue abierto desde r1) · Id `tomate-bola` con nombre «Jitomate bola»
- **Ruta:** `datos/semilla/productos.json`, producto `tomate-bola`.
- **Corrección sugerida:** renombrarlo a `jitomate-bola` antes de que lo citen el recetario, `nutricion/tabla.csv` o `profeco/vinculos.csv` (hoy nadie lo cita fuera del catálogo).

### B-4 (baja, sigue abierto desde r1, ahora con 2 productos) · Pesos por pieza de USDA sin paso de confirmación
- **Ruta:** `datos/semilla/productos.json`: `bolillo` (57 g, panecillo kaiser) y `caldo-pollo-cubos` (4 g, cubo de EE. UU.). El limón salió de la lista al pasar a `g`.
- **Corrección sugerida:** incluirlos en la lista que Axel confirma en la puerta 2 o cotejarlos con báscula o etiqueta antes de sembrar.

### B-5 (baja, nuevo) · El campo `version` del archivo no cambió con la corrección
- **Ruta:** `datos/semilla/productos.json` (`"version": "1.0"`) contra `datos/semilla/productos-README.md` («Versión 1.1»).
- **Qué pasa:** el contenido cambió (3 productos, 4 ids de presentación nuevos y 4 quitados) y el README sube a 1.1, pero el JSON sigue diciendo 1.0; quien compare versiones no distingue la semilla corregida de la original.
- **Corrección sugerida:** `"version": "1.1"` en el JSON.

### B-6 (baja, nuevo) · El reporte original del catálogo contradice el archivo
- **Ruta:** `docs/estado/wf1/catalogo-productos.md`, «Pendientes» (primer renglón) y «Verificación» (conteos `g 137 · pieza 12`, 15 productos con gramos).
- **Qué pasa:** dice que limón y aguacate siguen en `pieza` como pendiente; el archivo ya los tiene en `g` (g 140 · pieza 9 · 12 con gramos). `correccion-catalogo-r1.md` lo pidió como solicitud 1 al integrador y no se ha aplicado.
- **Corrección sugerida:** aplicar esa solicitud (o una nota al inicio que remita a `correccion-catalogo-r1.md`).

### B-7 (baja, nuevo; del contrato, no del catálogo) · Presentaciones a granel y paquetes enteros en la lista
- **Ruta:** `docs/contratos/motores/S4-interfaz.md` §6.2 pasos 5 y 6, contra `productos-formato.md` §2 (`granel: true` = «se vende por kilo o por unidad base»).
- **Qué pasa:** S4 nunca usa `granel`: `packages_to_buy` = techo(`needed_milli / size_milli`). Con la corrección, una receta que pide 134 g de limón genera en la lista «1 × Limón, por kilo» y un sobrante previsto de 866 g que entra al inventario al confirmar la compra; pasa igual con las otras 83 presentaciones «por kilo» (86 en total: jitomate, cebolla, carne…). El costo por receta de §8.2 sí sale bien (fracción exacta). El catálogo cumple los dos contratos; la regla falta en S4.
- **Corrección sugerida (fuera del alcance del catálogo, solicitud al integrador / contrato S4S5):** en §6.2, si la presentación que rige tiene `is_bulk = 1`, comprar `needed_milli` redondeado a la fracción de venta que se decida (o tal cual) en lugar de paquetes enteros; si es decisión de Axel, abrir el bloqueo correspondiente.

## Notas sin severidad
- La tarea pide «los 61 nombres»: en `data.json` hay 39 nombres de presentación y 36 textos distintos en 39 renglones; 61 es el número de piezas. Se verificaron los tres conteos (criterios 13 y 14).
- `docs/estado/bloqueos/INDICE.md` no lista `WF1-catalogo-huevo-por-kilo` (tampoco los demás bloqueos de WF-1); el integrador debe sumarlo a la lista de la puerta 2.
- Los sinónimos con marcas comunes («nescafé», «maicena», «royal», «queso filadelfia», «cloralex», «danonino») no son marcas inventadas: el formato solo prohíbe marcas en `.nombre` y ayudan a empatar tickets.

## Estado de los hallazgos de la ronda 1
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| M-1 (media) limón, aguacate y tortilla de maíz en `pieza` sin precio por kilo; sin bloqueo | **Cerrado** | `limon`, `aguacate`: `g`, única presentación «por kilo (granel)» = 1000000, `granel: true`, típica esa misma, gramos `null`. `tortilla-maiz`: `g`, «por kilo (granel)» = 1000000 + «paquete 500 g (aprox. 20 piezas)» = 500000, `granel: false`. `huevo`: sigue en `pieza` con el bloqueo `WF1-catalogo-huevo-por-kilo.md` (Contexto, Qué falta, 4 opciones, Impacto, Cómo seguí: formato del protocolo §2) y la nota en el README. Con eso un precio por kilo entra a una presentación y S4 §8.2 calcula el costo |
| B-1 básicos faltantes | Sigue abierto (B-1) | Faltan 26 de 157 de mi lista |
| B-2 sinónimos que empatan con otro producto | Sigue abierto (B-2) | Los 3 sinónimos siguen |
| B-3 id `tomate-bola` | Sigue abierto (B-3) | `tomate-bola` existe; `jitomate-bola` no |
| B-4 pesos por pieza de USDA | Sigue abierto (B-4), 2 productos | `bolillo` 57 y `caldo-pollo-cubos` 4 con «pendiente cotejar» |
