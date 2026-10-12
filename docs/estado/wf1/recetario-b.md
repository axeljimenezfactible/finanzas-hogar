# Reporte WF-1 · recetario-b (Sonnet) · 11-oct-2026

## Hecho
- 30 recetas de cocina casera mexicana para 2 adultos, comidas y cenas, escritas en `datos/semilla/recetario-b.json` con la estructura de `docs/contratos/motores/S4-interfaz.md` §3.1 (los campos de §3.2, ni uno más por receta).
- Todos los ingredientes usan ids de `datos/semilla/productos.json` (70 productos distintos), cantidades enteras en milésimas de la unidad base, `unidad` igual a la `unidad_base` del producto y la presentación comprable típica del catálogo (`presentacion_tipica`).
- Sin «al gusto»: sal, pimienta, orégano, comino y laurel llevan gramos exactos (la sal en gramos; el agua de cocción no es ingrediente del catálogo y se indica en los pasos).
- `kcal_por_porcion` = `null` y `fuentes_kcal` = `[]` en las 30 (los calcula el generador/motor; costos tampoco van).
- Las etiquetas (económica, rápida, sobrantes) viajan en la llave superior `etiquetas_por_receta` ({id: [etiquetas]}) porque §3.2 no define un campo de etiquetas por receta (ver Solicitudes).
- Los pasos son cadenas sin numerar; la numeración se hace al mostrar (§3.2).

## Archivos
- `/home/deepsec/finanzas-hogar/datos/semilla/recetario-b.json` (nuevo, 30 recetas, 107 KB con sangría de 2 espacios; sin sangría pesa casi la mitad)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/recetario-b.md` (este reporte)
- Scripts de apoyo (fuera del repositorio, en el scratchpad de la sesión): generador `build.py`, datos `part1.py`/`part2.py`, validador `validate.py`.

## Verificación
Comando: `python3 -I validate.py` (lee `productos.json` y `recetario-b.json`).
- Campos exactos por receta, id con patrón `^[a-z0-9]+(-[a-z0-9]+)*$` único, nombre único (3–80 caracteres), descripción ≤ 200, tiempos válidos sin repetidos, `porciones_base` = 2, minutos enteros ≥ 0, `origen` = `claude`, pasos no vacíos, `opcional` = false, un `producto_id` una vez por receta, producto y presentación existentes y de ese producto, unidad = unidad base, `cantidad_milli` entero > 0, `kcal_por_porcion` null y `fuentes_kcal` vacío.
- Salida: `recetas 30 errores 0`.
- Revisión cruzada pasos/ingredientes: todo ingrediente mencionado en los pasos (limón, cilantro, cebolla, ajo, crema, tortilla, etc.) está en la lista de la receta; 0 avisos.
- Conteo por tiempo de comida (una receta puede tener varios): comida (`lunch`) 27, cena (`dinner`) 28, desayuno (`breakfast`) 4, colación (`snack`) 0.
- Conteo por proteína principal: pollo 6, res 5, cerdo 2 (res + cerdo = 7 de carne no especializada), pescado 4 (tilapia 2, atún 2), huevo 3, legumbre 6, pasta 3, ensalada 1 (la de nopales con queso panela). Total 30.
- Etiquetas: sobrantes 11, económica 16, rápida 9.
- Masa por porción (suma de gramos de ingredientes / 2): mínimo 290 g (frijoles de la olla, en crudo) y máximo 998 g (caldo de res con verduras, con líquido de verduras y elote); sin valores absurdos.
- Criterio S4-01 (parcial, este lote): 30 de 30 válidas contra las reglas de §3.2; la unión a 60 y el esquema `recetario.schema.json` son del integrador/agente del recetario-a.

## Supuestos
- S-1: «Presentación comprable» = la `presentacion_tipica` de cada producto (granel para carne, verdura y tortilla; paquete donde el catálogo lo da). No cambia ninguna cifra; el motor puede reemplazarla por la que rija la cotización.
- S-2: Las recetas «de sobrantes» son platillos que se conservan y recalientan bien (guisos, caldos, tinga, frijoles, lentejas) o que aprovechan tortillas del día anterior (chilaquiles). El plan debe repetirlas; el motor solo descuenta sobrantes registrados de la misma receta (S4-interfaz §6.4). No se declara cruce entre recetas distintas.
- S-3: Lo que en el catálogo está en `g` pero en la cocina se cuenta por pieza (tortillas de maíz, limones, jitomates) lleva el peso aproximado y la medida casera en `texto`. Pesos aproximados de cocina, no datos de fuente.
- S-4: Las recetas con garbanzo piden remojo la noche anterior; el remojo no cuenta en `minutos_preparacion`, pero va como primer paso.
- S-5: El recetario de Claude no lleva opcionales (§5.4), así que adornos como cilantro o limón son ingredientes obligatorios con cantidad exacta.
- S-6: No hay recetas de colación; la tarea pedía comidas y cenas.

## Bloqueos
Ninguno. No se escribió archivo en `docs/estado/bloqueos/`.

## Pendientes
- Regenerar `kcal_por_porcion` y `fuentes_kcal` con el generador del recetario (cuando exista `datos/semilla/nutricion/tabla.csv` en modo `review`); hoy van `null` y `[]`, por lo que en modo `review` el cotejo daría `ErrKcalMismatch` hasta regenerar. En modo `production` no afecta.
- Unir con `recetario-a` en `recetario.json` (60 en total) y comprobar que no haya ids ni nombres repetidos entre ambos archivos (en este lote no pude verlo: `recetario-a.json` no existía al correr).
- Las 4 recetas de desayuno se quedan en 4; si se quiere más cobertura de desayuno y colación, debe venir del otro lote.

## Solicitudes fuera de alcance
1. `datos/semilla/recetario.schema.json` y `docs/contratos/motores/S4-interfaz.md` §3.2: decidir si las etiquetas (económica, rápida, sobrantes) son un campo de la receta (`etiquetas`: arreglo de cadenas). Motivo: la tarea las pide y §3.2 no las define; para no romper un esquema con `additionalProperties: false` van en la llave superior `etiquetas_por_receta`, que el integrador debe quitar o convertir al unir.
2. Integrador de S4: al unir en `recetario.json`, eliminar la llave `etiquetas_por_receta` de la raíz si el esquema no la admite, y volver a correr el validador de nombres contra los nombres en inglés de `tabla.csv` (S4-11; los nombres de este lote son es-MX).
3. Catálogo (`datos/semilla/productos.json`): ver «Productos faltantes».

## Productos faltantes
Productos básicos que habría usado y que no están en el catálogo; NO los usé, cambié la receta:
- chícharo (verde o en lata) para el arroz con pollo,
- mostaza y pepinillos para hamburguesas,
- pasta de mole,
- aceitunas y alcaparras para la veracruzana,
- rábano,
- cacahuate o nuez,
- ajonjolí,
- miel,
- pimentón o paprika,
- cubos de caldo de res (se usó solo el de pollo),
- crema de elote enlatada.
Ninguno hizo falta para cumplir las 30 recetas.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no hay valores K31 en esta tarea).

## Consumo aproximado
- Acciones de herramienta: 18.
- Leído: ≈ 150 KB (protocolo 6, S4-interfaz por extractos ≈ 12, criterios S4 18, decisiones D y K ≈ 22, catálogo compactado ≈ 15, más salidas de verificación).
- Escrito: ≈ 150 KB (107 de JSON generado, ≈ 36 de datos y scripts, ≈ 7 de reporte); en el límite del protocolo §4.
