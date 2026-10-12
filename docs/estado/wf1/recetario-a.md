# recetario-a · WF-1 (11-oct-2026)

Agente `recetario-a` (Sonnet). Alcance: `datos/semilla/recetario-a.json` (30 recetas, mitad del recetario de 60 que une el integrador de S4 en `recetario.json`).

## Hecho
- 30 recetas de cocina casera mexicana para 2 adultos, porciones base 2, origen `claude`, sin restricciones alimentarias.
- Ingredientes solo con ids de `datos/semilla/productos.json` (73 productos distintos), cantidades enteras en milésimas de la unidad base, `unidad` igual a la `unidad_base`, `presentacion_id` = `presentacion_tipica` del producto, `texto` casero y `opcional: false` en todos.
- Sin «al gusto»: sal, pimienta, comino y chiles secos llevan gramos exactos. El agua no es ingrediente del catálogo; los pasos la dan en tazas.
- `kcal_por_porcion: null` y `fuentes_kcal: []` (los llena el generador del recetario con la tabla nutricional; costos no van en el archivo).
- Tiempos con los códigos de S4-interfaz §3.2 (`breakfast`, `lunch`, `dinner`, `snack`).

## Archivos
- `/home/deepsec/finanzas-hogar/datos/semilla/recetario-a.json` (nuevo).
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/recetario-a.md` (este reporte).
- Scripts de generación y validación (fuera del repositorio, en el scratchpad de la sesión): `.../scratchpad/ra/{rec1.py,rec2.py,build.py,validate.py}`.

## Verificación
Comando: `python3 -I validate.py` (lee `productos.json` y `recetario-a.json`).
- `recetas 30 errores 0`: campos exactos, id con patrón y único, nombre 3 a 80 caracteres y único sin distinguir mayúsculas, descripción ≤ 200, tiempos válidos sin repetidos, porciones = 2, minutos enteros ≥ 0, origen `claude`, pasos no vacíos, ingredientes sin repetir producto, `producto_id` y `presentacion_id` existentes y del mismo producto, `unidad` = unidad base (g, ml, pieza), cantidad entera > 0, `opcional` falso, `texto` presente.
- Todo producto en `ml` o `pieza` usado tiene `gramos_por_unidad` (0 avisos de kcal incompleta).
- Conteo por tiempo (una receta puede tener varios): desayuno 10, comida 20, cena 17, colación 7.
- Proteína principal (receta cuenta en cada una que contenga): huevo 7, legumbre 7, res 4, cerdo/embutido 4, pollo 3, pescado/enlatado 2, sin carne ni legumbre ni huevo (verdura, lácteo o cereal) 9.
- Etiquetas: económica 20, rápida 14, sobrantes 7. Tiempo total máximo 135 min (frijoles de la olla, 120 min de cocción).
- Criterio S4-01: solo la parte de este lote (30 de 60); falta `recetario.schema.json` (no existe aún; no era de mi alcance) y las otras 30 recetas.

## Supuestos
1. El campo `etiquetas` (económica, rápida, sobrantes) lo pidió la tarea pero **no está** en S4-interfaz §3.2. Lo incluí; si el esquema final tiene `additionalProperties: false`, el integrador debe agregarlo al esquema o quitarlo al unir (ver Solicitudes).
2. «Sobrantes» = la receta indica cómo reaprovechar lo que sobra (guiso que rinde otra comida o relleno para el día siguiente). Las cantidades son siempre para 2 porciones; el sobrante entra al inventario por el mecanismo normal de S4 (paquetes completos), no con un campo propio.
3. `presentacion_id` = `presentacion_tipica` en todos los casos (la presentación por defecto del catálogo); el motor redondea a presentación al armar la lista.
4. Granel y piezas: cantidades en gramos para verdura y carne; huevo, bolillo, tostadas y cubos de caldo en piezas.
5. Los ingredientes de pocas unidades (1 cubo de caldo, 1 g de comino, 0.5 g de pimienta, 1 g de laurel, 1 g de chile de árbol) son estimaciones de cocina casera; se redondean a 1 gramo o más porque la menor unidad útil de la cocina es la pizca.
6. Las cantidades de lentejas y garbanzos son en crudo (secos); los garbanzos piden remojo de 8 horas, anotado en los pasos y fuera de `minutos_preparacion`.

## Bloqueos
Ninguno. No dependo del bloqueo S1-parametros-avisos.

## Pendientes
- `datos/semilla/recetario.schema.json` (esquema JSON): lo escribe el agente del recetario o el integrador de S4; hay que validar `recetario-a.json` contra él cuando exista.
- Las otras 30 recetas (`recetario-b.json`) y la unión en `datos/semilla/recetario.json` (60 exactas).
- Llenar `kcal_por_porcion` y `fuentes_kcal` con el generador (depende de `nutricion/tabla.csv`).
- Productos faltantes deseables, **no usados** (ninguna receta quedó sin ellos): miel, mermelada, pasas, nuez o cacahuate (avena y yogur), chícharo, granola, vainilla. Ninguno es imprescindible.
- Posible duplicado con `recetario-b`: el integrador debe cotejar nombres (por ejemplo, ensalada de atún, tinga, sopa de fideo) antes de unir, porque el nombre es único sin distinguir mayúsculas.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S4-interfaz.md` §3.2 (y `datos/semilla/recetario.schema.json` cuando exista): agregar el campo opcional `etiquetas` (arreglo de cadenas; valores `económica`, `rápida`, `sobrantes`) o decidir que se omite. Motivo: lo pide la tarea y sirve a la interfaz para filtrar.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no hay valores K31 en este alcance).

## Consumo aproximado
Acciones: 12. Leído: ≈ 85 KB (productos.json resumido por script ≈ 106 KB procesados sin imprimir íntegros; contrato S4 §3 y criterios por extractos; la lectura accidental del archivo completo de S4-interfaz fue de 56 KB). Escrito: ≈ 125 KB (datos 80 KB, scripts ≈ 40 KB, reporte ≈ 5 KB).
