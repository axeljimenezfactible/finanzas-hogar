# Catálogo base de productos (`datos/semilla/productos.json`)

Versión 1.1 · 11-oct-2026 · WF-1, agente `catalogo-productos`; ronda de corrección `correccion-catalogo-r1` (limón, aguacate y tortilla de maíz pasan a gramos). Formato: `docs/contratos/motores/productos-formato.md` (1.1). Solo siembra: Ajustes lo edita después (`createProduct`, `updateProduct`, `createPresentation`).

## Qué contiene
- **160 productos** y **186 presentaciones** en **9 grupos** (los de `productos-formato.md` §4, sin cambios).
- Los **36 productos de `data.json`** (39 presentaciones con su nombre exacto, 61 piezas en 39 renglones de `cuentas`, 36 textos de ticket distintos como `nombres_ticket`).
- **124 productos nuevos** para cocina casera mexicana: verduras, frutas, carne y pollo no especializados, huevo, lácteos, tortillas y pan, abarrotes, legumbres, pastas, chiles secos y especias.
- Sin precios, sin calorías y sin marcas inventadas: las marcas solo aparecen en las presentaciones que vienen de `data.json`. Los demás productos usan presentaciones genéricas («Azúcar estándar, bolsa 1 kg», «Cebolla blanca, por kilo (granel)»).

## Conteos
| Grupo (canal) | Productos |
|---|---|
| Frutas y verduras (mercado) | 59 |
| Carne y pollo no especializada (mercado) | 14 |
| Abarrotes (cadena) | 53 |
| Lácteos y huevo (cadena) | 15 |
| Salchichonería (cadena) | 4 |
| Panadería (cadena) | 5 |
| Limpieza del hogar (cadena) | 4 |
| Higiene personal (cadena) | 3 |
| Papel y desechables (cadena) | 3 |

Vertical: `mercado` 73 · `cadena` 67 · `ambas` 20. Unidad base: `g` 140 · `ml` 11 · `pieza` 9.

## Reglas que siguió
- `id` en minúsculas con guiones, estable (lo citarán `recetario.json`, `nutricion/tabla.csv` y `profeco/vinculos.csv`). Presentación: `<id del producto>.<slug>`.
- `contenido_milli` es el factor a la unidad base: 1 kg = 1000000 g; cartón de 12 huevos = 12000 piezas. `granel: true` en lo que se vende por kilo o por pieza suelta.
- `vertical = ambas` solo en productos de grupo `cadena` que también se consiguen en mercado o tianguis (huevo, quesos frescos, crema, bolillo, tortilla de maíz, frijol a granel, chiles secos, piloncillo, jamaica, maíz pozolero, chorizo, requesón). Verduras, fruta y carne quedan en `mercado`; una cadena con precio compite igual (`productos-formato.md` §3).
- `gramos_por_unidad` solo donde se usa en recetas con unidad `ml` o `pieza` (12 productos); cada uno con su fuente (`fdc_id` de USDA FoodData Central, SR Legacy 2018-04, CC0). En `ml` es la densidad (g/ml = gramos de una taza ÷ 236.588). Los productos de limpieza, higiene y papel no llevan dato (no entran en recetas).
- Sinónimos en minúsculas y en español de México; ningún sinónimo ni nombre se repite entre productos. Las variantes sin acento se incluyen donde el OCR o el teclado las suelen producir; la búsqueda debe normalizar acentos de todos modos.

## Qué falta confirmar (nada de esto es precio ni calorías)
- **Contenidos nominales sin marca.** Los tamaños de las presentaciones genéricas (bolsa 1 kg, lata 220 g, frasco 50 g, manojo aprox. 50 g…) son los típicos del mercado, no leídos de una etiqueta. Axel los ajusta en Ajustes la primera vez que compre cada uno. Los de `data.json` sí son exactos.
- **Pesos por pieza de USDA** que pueden diferir del producto mexicano: bolillo (57 g de panecillo tipo kaiser) y cubo de caldo (4 g de USDA contra la etiqueta de Knorr). Cada `fuente_gramos` lo dice.
- **Limón, aguacate y tortilla de maíz (corregido en 1.1).** Se compran por kilo en México (mercado, cadena, tortillería) y `contenido_milli` va en milésimas de la unidad base, así que pasaron de `pieza` a `g`: presentación «por kilo (granel)» = 1000000 g, sin `gramos_por_unidad` (prohibido en `g`). La receta que pida «2 limones» lo dice en `recipe_ingredient.texto`; su cantidad de cálculo va en gramos. La tortilla conserva además un «paquete 500 g (aprox. 20 piezas)», contenido nominal típico sin marca que Axel ajusta en Ajustes. Id de presentación nuevos: `limon.kilo-granel`, `aguacate.kilo-granel`, `tortilla-maiz.kilo-granel`, `tortilla-maiz.paquete-500g` (los anteriores `*.pieza`, `tortilla-maiz.paquete-20/40` ya no existen; nadie los citaba todavía).
- **Huevo** sigue en `pieza` (cartón 12, cartón 18, plato 30). Si el mercado lo cobra por kilo, ese precio no entra a ninguna presentación: está en el bloqueo `docs/estado/bloqueos/WF1-catalogo-huevo-por-kilo.md` (PROVISIONAL solo para ese caso; el producto y sus presentaciones son válidos tal cual).
