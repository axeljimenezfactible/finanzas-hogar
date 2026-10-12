# Reporte · revisor-catalogo-r2 (WF-1, ronda 2 de 3)

## Hecho
- Revisión independiente de `datos/semilla/productos.json` tras `correccion-catalogo-r1`, solo contra `productos-formato.md` 1.1 y `S4-interfaz.md`. Informe: `docs/estado/wf1/REVISION-catalogo-r2.md`.
- Veredicto: **APROBADO** (0 altos, 0 medios, 7 bajos). M-1 de r1 cerrado; B-1 a B-4 de r1 siguen abiertos (bajos); 3 bajos nuevos (B-5 versión del JSON, B-6 reporte original desactualizado, B-7 regla de granel que falta en S4 §6.2).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-catalogo-r2.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-catalogo-r2.md` (este reporte)
- Temporales en `/home/deepsec/finanzas-hogar/datos/tmp/rev-cat-r2/` (`v2.py`, `v3.py`, `cov.py`, `usda.py`, descarga de USDA en `usda-raw/`): borrados al terminar.

## Verificación
- `python3 -I datos/tmp/rev-cat-r2/v2.py <raíz>` (validador propio): productos 160 · presentaciones 186 · grupos 9 · unidad g 140 / ml 11 / pieza 9 · vertical cadena 67 / mercado 73 / ambas 20 · con gramos 12; nombres de `data.json` 39 (39 distintos), conteo ≠ 1: ninguno; renglones de `cuentas` 39, textos distintos 36, piezas 61, importe 3094.5 = totales; textos sin exactamente 1 presentación: ninguno; 73 cadenas de producto de todo `data.json` sin alias: ninguna; precio unitario por renglón contra `bodega_unit` en centavos: 0 diferencias; piezas por presentación = cantidad (atún 6 = 5 + 1); Walmart 60 = 60; **ERRORES 0**.
- `v3.py`: factor `contenido_milli` recalculado desde el nombre en 186 de 186 (1 diferencia aparente del propio regex, leche 4 × 1 l, correcta); tabla §5.1 36 filas 0 diferencias; tabla §5.2 39 filas 0 diferencias; 30 marcas en 38 presentaciones, 0 fuera de `data.json`; 9 productos en `pieza` revisados uno por uno; estado de B-2, B-3, B-4 y M-1.
- USDA FoodData Central (API, `DEMO_KEY`, 12 `fdc_id`): 12 de 12 porciones iguales al catálogo; densidades recalculadas con Decimal y HALF_UP (0.913, 0.921, 1.078, 1.162, 1.031, 1.006).
- `cov.py`: lista propia de 157 básicos, 131 presentes (83.4 %).
- `grep -rln` de los ids quitados en todo el repositorio sin `.git`: solo README (historia) y `correccion-catalogo-r1.md`; `grep` de limón, aguacate, tortilla y huevo en `docs/contratos/motores/`, `profeco/FORMATO.md` y fixtures: 0 líneas.
- Criterios de fase: no hay criterio de `docs/criterios/` propio del catálogo; se cubrió lo que el catálogo alimenta de S4-01 (presentación comprable y unidad base en cada producto: 160 de 160), S4-11 (nombres visibles en español: lectura de los 160) y S4-12 (0 productos sin canal).

## Supuestos
1. La descarga de USDA se hizo en una carpeta propia dentro de `datos/tmp/rev-cat-r2/usda-raw/` (sin scripts dentro) y no en `datos/semilla/usda/raw/`, porque el revisor no escribe en `datos/semilla/`; se borró al terminar.
2. «Los 61 nombres» de la tarea se leyó como los tres conteos de `data.json` (39 nombres de presentación, 36 textos en 39 renglones, 61 piezas); los tres cuadran.

## Bloqueos
Ninguno nuevo. Revisado: `docs/estado/bloqueos/WF1-catalogo-huevo-por-kilo.md` (cumple el formato del protocolo §2).

## Pendientes
Ninguno del revisor. Los 7 bajos quedan para la ronda de corrección si el orquestador la abre.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S4-interfaz.md` §6.2 pasos 5 y 6: definir qué se compra cuando la presentación que rige es a granel (`is_bulk = 1`), en lugar de paquetes enteros de 1 kg (B-7). Motivo: 86 presentaciones «por kilo» generarían compras y sobrantes de hasta casi 1 kg.
2. `docs/estado/wf1/catalogo-productos.md` «Pendientes» y «Verificación»: aplicar la solicitud 1 de `correccion-catalogo-r1.md` (B-6).
3. `docs/estado/bloqueos/INDICE.md`: sumar `WF1-catalogo-huevo-por-kilo` a la lista de la puerta 2.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna sobre K31.

## Consumo aproximado
Acciones de herramienta: 29 · KB leídos ≈ 180 (protocolo 6, decisiones §K 15, criterios S4 18, formato 19, S4-interfaz ≈ 25 por secciones, revisión r1 12, corrección 6, bloqueo 2, reporte 6, README por grep 3, salidas de scripts ≈ 25, listado del catálogo ≈ 20) · KB escritos ≈ 35 (informe 13, reporte 4, scripts 18).
