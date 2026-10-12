# correccion-catalogo-r1 · WF-1 · 11-oct-2026

## Hecho
- **M-1 (media), opción A del revisor.** `limon`, `aguacate` y `tortilla-maiz` pasan de `pieza` a `g`: presentación «por kilo (granel)» = 1000000 con `granel: true`, `presentacion_tipica` apuntando a ella, y `gramos_por_unidad` / `fuente_gramos` en `null` (prohibido en `g`). La tortilla conserva un «paquete 500 g (aprox. 20 piezas)» = 500000 en lugar de los paquetes de 20 y 40 piezas. «2 limones» se expresa en `recipe_ingredient.texto`, con la cantidad de cálculo en gramos.
- Id de presentación nuevos: `limon.kilo-granel`, `aguacate.kilo-granel`, `tortilla-maiz.kilo-granel`, `tortilla-maiz.paquete-500g`. Desaparecen `limon.pieza`, `aguacate.pieza`, `tortilla-maiz.paquete-20`, `tortilla-maiz.paquete-40` (verificado: nada en `datos/semilla/` ni `docs/contratos/motores/` los cita; `nutricion/` está vacía y `profeco/` solo tiene los CSV de origen).
- **Huevo (relacionado en el hallazgo).** Queda en `pieza` (cartón 12/18, plato 30) porque pasarlo a gramos rompe «2 huevos» y su conteo exacto de piezas; el caso del precio por kilo en mercado es una decisión de modelo, escrita como bloqueo `docs/estado/bloqueos/WF1-catalogo-huevo-por-kilo.md` (4 opciones, impacto en S4 §8.2 y PROFECO), mencionado en el README. No se tocó `huevo` en `productos.json`.
- README a versión 1.1: conteos (unidad base g 140 · ml 11 · pieza 9; `gramos_por_unidad` en 12 productos), se quitó la advertencia del peso del limón y se reescribió la nota de limón/aguacate/tortilla/huevo.
- No se tocó nada más: los otros 157 productos y los 9 grupos son idénticos al respaldo. Las presentaciones siguen siendo 186 (limón 1→1, aguacate 1→1, tortilla 2→2).
- No se atendieron B-1 a B-4 (severidad baja; la tarea pidió solo alta y media). Siguen abiertos tal como los reporta `REVISION-catalogo-r1.md`.

## Archivos
- `/home/deepsec/finanzas-hogar/datos/semilla/productos.json` (3 productos modificados; mismo formato `indent=1`)
- `/home/deepsec/finanzas-hogar/datos/semilla/productos-README.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/WF1-catalogo-huevo-por-kilo.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-catalogo-r1.md` (este reporte)

## Verificación
Validador propio (`val.py` en el scratchpad de la sesión; reglas de `productos-formato.md` §1–§2, §3, §5.2), corrido sobre el archivo nuevo y sobre el anterior:
- Nuevo: 160 productos, 186 presentaciones, 9 grupos; unidad `g` 140 · `ml` 11 · `pieza` 9; 12 con `gramos_por_unidad`; **0 errores**: claves exactas, patrón/unicidad de ids, nombres únicos, categoría, unidad, `gramos_por_unidad` prohibido en `g`, `vertical` vs canal, sinónimos en minúsculas y sin compartir (0 compartidos), `presentacion_tipica` pertenece al producto, `contenido_milli` entero > 0, `codigo_barras`, 39 nombres de `data.json` (cada uno exactamente 1 vez), 39 renglones de `cuentas` con 61 piezas todos en algún `nombres_ticket`, 0 alimentos en `ml`/`pieza` sin `gramos_por_unidad`.
- Anterior (respaldo): mismos 0 errores, g 137 · pieza 12 · 15 con dato.
- Diferencia por producto contra el respaldo: cambiaron exactamente `limon`, `aguacate`, `tortilla-maiz`; ids de productos idénticos; `grupos` idénticos.
- Salida de la lectura del estado final: `limon g None [('limon.kilo-granel', 1000000, True)]`; `aguacate` igual; `tortilla-maiz g None [('…kilo-granel', 1000000, True), ('…paquete-500g', 500000, False)]`.
- Criterio 18 del revisor (muestra de 20): con el cambio, limón, aguacate y tortilla de maíz tienen unidad base y presentación sensata (20 de 20).
- Pruebas Go, oráculo y `compara.py`: no aplican, el artefacto es un JSON de siembra sin código; no se tocó ningún archivo de código.

## Supuestos
1. «Paquete 500 g (aprox. 20 piezas)» de tortilla es un contenido nominal típico sin marca (el reviewer sugirió «gramos nominales»); no es lectura de etiqueta. Lo ajusta Axel en Ajustes igual que los demás contenidos genéricos (README).
2. Con la unidad en gramos, la cantidad de cálculo de «2 limones» la fija el recetario en gramos (no se inventó un gramaje por pieza: el 67 g de USDA era de una lima; ya no se usa).
3. El huevo se dejó en `pieza` mientras Axel/S5 no decidan; el catálogo es válido con ambas respuestas.

## Bloqueos
- `docs/estado/bloqueos/WF1-catalogo-huevo-por-kilo.md` (nuevo): cómo registrar un precio por kilo de un producto con unidad `pieza` (huevo). No impide nada.

## Pendientes
- B-1 a B-4 de `REVISION-catalogo-r1.md` (básicos faltantes, 3 sinónimos que empatan con otro producto, id `tomate-bola`, confirmación de pesos por pieza USDA de bolillo y cubo de caldo) quedan sin atender por estar fuera de «alta y media».
- Los pesos del limón y el aguacate por pieza ya no están en el catálogo; si el recetario necesita convertir piezas a gramos, los tomará de su propia tabla (nutrición) y no de `productos.json`.

## Solicitudes fuera de alcance
1. `/home/deepsec/finanzas-hogar/docs/estado/wf1/catalogo-productos.md`, sección «Pendientes», primer renglón: sustituirlo por «Corregido en `correccion-catalogo-r1`: limón, aguacate y tortilla de maíz pasaron a `g` con presentación por kilo; el huevo queda en el bloqueo `WF1-catalogo-huevo-por-kilo`». Motivo: el pendiente ya no es cierto y el reporte original dice lo contrario del archivo.
2. El agente del recetario (60 recetas) debe expresar limón, aguacate y tortilla en gramos con `texto` («2 limones», «8 tortillas»). Archivo: `datos/semilla/recetario.json` (aún no existe).
3. El agente de `nutricion/tabla.csv` debe incluir `gramos_por_unidad` de esos tres productos solo si lo necesita su propia conversión; el catálogo ya no los lleva.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna sobre K31. Decisión pedida (puerta 2): opción del bloqueo `WF1-catalogo-huevo-por-kilo` (recomendada la 1: dejar el huevo en pieza y capturar por pieza o cartón).

## Consumo aproximado
Acciones: 15 · KB leídos ≈ 55 (protocolo §2/§4, informe de revisión, README, 4 productos y vecinos, trozos de formato y data.json) · KB escritos ≈ 20 (productos.json reescrito 106 KB por script, que no cuenta como lectura/escritura del agente por herramienta de edición; README, bloqueo y reporte ≈ 10 KB).
