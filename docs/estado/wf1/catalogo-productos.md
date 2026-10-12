# Reporte · catalogo-productos (WF-1)

## Hecho
- `datos/semilla/productos.json` v1.0: 9 grupos (`productos-formato.md` §4), 160 productos, 186 presentaciones.
- Los 36 productos de `data.json` tal como los fija `productos-formato.md` §5.1: 39 presentaciones con el texto exacto de `items_identicos[].producto` y `sustitutos[].bodega/walmart`; los 36 textos distintos de `cuentas` (39 renglones, 61 piezas) como `nombres_ticket` según §5.2.
- 124 productos nuevos para unas 60 recetas caseras (frutas y verduras 58, carne y pollo 14, lácteos y huevo 10, salchichonería 2, panadería 4, abarrotes 36).
- `gramos_por_unidad` con `fdc_id` de USDA (SR Legacy 2018-04, descargado a `datos/tmp`, borrado al terminar) en los 15 productos de comida con unidad `ml` o `pieza`: huevo, limón, aguacate, elote en mazorca, bolillo, tortilla de maíz, tortilla de harina, tostadas, caldo en cubos, aceite de oliva, aceite de soya, leche, vinagre, salsa de soya, salsa inglesa.
- `datos/semilla/productos-README.md`.

## Archivos
- `/home/deepsec/finanzas-hogar/datos/semilla/productos.json` (106 KB, nuevo)
- `/home/deepsec/finanzas-hogar/datos/semilla/productos-README.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/catalogo-productos.md` (este reporte)
- Generador y validador en el scratchpad de la sesión (no van al repositorio); el validador completo está abajo.

## Verificación
`python3 -I validar.py /home/deepsec/finanzas-hogar` (script al final de este archivo):
```
productos 160 presentaciones 186 grupos 9
data.json: presentaciones exactas 39 renglones de cuentas 39 piezas 61 textos de ticket únicos 36 cruces de precio unitario 39
categoria {'abarrotes': 53, 'carnes': 14, 'frutas-verduras': 59, 'higiene': 3, 'lacteos': 15, 'limpieza': 4, 'panaderia': 5, 'papel': 3, 'salchichoneria': 4}
vertical {'ambas': 20, 'cadena': 67, 'mercado': 73}
unidad {'g': 137, 'ml': 11, 'pieza': 12}
ERRORES 0
```
Comprobado: campos obligatorios; ids de producto y de presentación únicos y con el patrón; nombres únicos sin distinguir mayúsculas; `categoria` existe; unidad en {g, ml, pieza}; `contenido_milli` entero > 0; `gramos_por_unidad` con ≤ 3 decimales, > 0, ausente en `g` y con `fuente_gramos`; vertical coherente con el canal del grupo; `presentacion_tipica` pertenece al producto; sinónimos en minúsculas, sin repetidos ni compartidos entre productos; **los 39 nombres de presentación de `data.json` aparecen exactamente una vez**; los 36 textos de ticket (39 renglones, 61 piezas) aparecen en exactamente una presentación; **el precio unitario de cada renglón (importe ÷ cantidad) coincide con `bodega_unit`** del artículo que le corresponde (39 de 39).
No hay criterio de fase propio para este archivo en `docs/criterios/`; lo verifican S4-01 (`TestRecetarioValido`) y S5 cuando carguen el catálogo.

## Supuestos
1. «Nombre exacto como primer sinónimo» (tarea) contra `productos-formato.md` §2 (el texto exacto va en `presentacion.nombre` y en `nombres_ticket`; los sinónimos son genéricos en minúsculas): se siguió el formato, que «lo define todo». Los 61 artículos se cubren así: 39 nombres de presentación exactos + 36 textos de ticket exactos. Si se quiere el texto literal también entre los sinónimos, es un cambio mecánico.
2. Los contenidos de las presentaciones genéricas (sin marca) son típicos nominales, no leídos de etiqueta (README, «Qué falta confirmar»). No tocan cifras de `oraculo/model.json`.
3. `limpieza`, `higiene` y `papel` no llevan `gramos_por_unidad` (no entran en recetas); el formato lo exige solo para productos usados en recetas.
4. La búsqueda normaliza acentos y mayúsculas; por eso no se duplicaron todos los sinónimos sin acento.
5. Los pesos USDA son de alimentos estadounidenses; para limón, bolillo y cubo de caldo es solo la mejor fuente citable y queda dicho en `fuente_gramos`.
6. Sin código de barras (`codigo_barras: null`) en todo el catálogo: ninguno se conoce con certeza y no se inventan.

## Bloqueos
Ninguno. Ni S1-parametros-avisos ni otro bloqueo abierto tocan este archivo.

## Pendientes
- Limón y aguacate (y cualquier fruta que el mercado cobre por kilo y la receta cuente por pieza): `pieza` con presentación «por pieza» no representa un precio por kilo. Decisión del integrador/S5: añadir presentación en gramos con conversión por `gramos_por_unidad`, o capturar el precio por pieza. No bloquea la semilla.
- Los recetarios pueden pedir productos que no estén: se añaden con `createProduct` o ampliando este archivo (los ids existentes no cambian). Quedan fuera por marginales: ajo y cebolla en polvo, mermelada, galletas María, pasas, cacahuate, ajonjolí, miel y rábano (se descartaron para quedar en 160).
- Tabla `nutricion/tabla.csv` y `profeco/vinculos.csv`: otros agentes; usan estos `id`.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/productos-formato.md` §5.1/§2: aclarar que el nombre literal de `data.json` va en `presentacion.nombre` y `nombres_ticket`, no en `sinonimos` (contradice el texto de la tarea del flujo, supuesto 1). Motivo: evitar que un revisor lo lea como incumplimiento.
2. Opcional: copiar el validador de abajo a `scripts/validar_productos.py` y llamarlo desde la prueba Go del cargador. Motivo: reproducibilidad.

## Paquetes instalados
Ninguno. (`curl` y `python3 -I` ya existían; descarga de USDA SR Legacy en `datos/tmp/usda-sr-raw`, borrada.)

## Propuestas para Axel
Ninguna: no se usó ningún valor operativo de K31.

## Consumo aproximado
~22 acciones de herramienta · ~95 KB leídos (protocolo 6 KB, formato 19 KB, decisiones y S4 ~25 KB, `data.json` por extractos, USDA por búsqueda de texto) · ~135 KB escritos (JSON 106 KB, README 4 KB, reporte 6 KB, scripts de generación ~20 KB).

## Validador (`validar.py`, `python3 -I validar.py <raíz>`)
Verifica lo listado en «Verificación». Resumen de reglas en código: patrón `^[a-z0-9]+(-[a-z0-9]+)*$`; `unidad_base ∈ {g,ml,pieza}`; `gramos_por_unidad` con `^\d+(\.\d{1,3})?$` y > 0; `contenido_milli` entero > 0; `vertical` vs `canal`; 39 nombres de `items_identicos`/`sustitutos` con conteo exacto 1; 39 renglones de `cuentas` con su texto en una sola presentación; `abs(importe/cant − bodega_unit) ≤ 0.006`. (El script completo vive en el scratchpad de la sesión: `/tmp/claude-1000/-home-deepsec/ad9b2909-7e27-4a0a-9e87-9e6a5395adce/scratchpad/cat/validar.py`.)
