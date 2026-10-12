# Reporte · revisor-catalogo-r1 (WF-1, ronda 1 de 3)

## Hecho
- Revisión independiente de `datos/semilla/productos.json` v1.0 contra `docs/contratos/motores/productos-formato.md` 1.1 y `docs/contratos/motores/S4-interfaz.md`, con scripts propios (no se usó el validador del autor).
- Informe: `docs/estado/wf1/REVISION-catalogo-r1.md`. Veredicto **NO APROBADO**: 0 altos, 1 medio (M-1), 4 bajos (B-1 a B-4).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-catalogo-r1.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-catalogo-r1.md` (este reporte)
- Temporales borrados al terminar: `datos/tmp/rev-cat-r1/` (v1.py, v2.py, v3.py) y la descarga `datos/semilla/usda/raw/revisor-catalogo-r1/foods.json` (carpeta creada por este agente y eliminada completa).

## Verificación
- `python3 -I datos/tmp/rev-cat-r1/v1.py <raíz>` (campos, ids, patrón, unicidad, grupos contra §4, unidades, `gramos_por_unidad`, vertical contra canal, sinónimos, `presentacion_tipica`, `contenido_milli`, `granel`, códigos de barras; 39 nombres de `data.json`; 39 renglones de `cuentas` contra §5.2; precio unitario en centavos; tabla §5.1): `productos 160 presentaciones 186`, `nombres data.json 39 distintos 39`, `cuentas renglones 39 textos distintos 36 piezas 61 importe 3094.5`, `filas §5.2 leídas 39`, `filas §5.1 leídas 36`, **`ERRORES 0`** (la única nota es la sexta lata de atún, documentada en el formato §5).
- `python3 -I datos/tmp/rev-cat-r1/v2.py <raíz>` (listado completo, marcas, factor contra el tamaño del nombre, fuentes de gramos): 0 marcas fuera de `data.json`; 0 diferencias de factor en 64 presentaciones con tamaño explícito; 83 «por kilo (granel)» = 1000000; 0 alimentos en ml o pieza sin `gramos_por_unidad`.
- USDA FoodData Central (API oficial, una petición POST con los 15 `fdc_id` citados): las 15 porciones coinciden con `gramos_por_unidad`; densidades recalculadas con Decimal HALF_UP (0.913, 0.921, 1.078, 1.162, 1.031, 1.006) iguales al archivo.
- `python3 -I datos/tmp/rev-cat-r1/v3.py <raíz>`: 89 de 105 básicos presentes (84.8 %; el script contó «café» como faltante por empate exacto, corregido a mano: hay café soluble y café molido); piezas Walmart 60 = `articulos_walmart`.
- Criterios de fase: no hay archivo de criterios propio del catálogo; se cubren las reglas de validación de `productos-formato.md` §2–§5 y lo que S4 exige del catálogo (S4 §3.2 presentaciones del producto, §6.5/§11 canal, §9.2 gramos con fuente; S4-11 nombres en español; S4-12 ingredientes sin canal = 0). Tabla con evidencia en el informe.

## Supuestos
1. «Los 61 nombres de data.json» se leyó como: 39 nombres de presentación exactos + 39 renglones de `cuentas` (36 textos) que suman 61 piezas; se verificaron los tres conteos.
2. Los modos de compra por kilo de limón, aguacate, tortilla y huevo en México son conocimiento general del revisor (no hay dato del proyecto que los contradiga); sostienen M-1 junto con la propia nota del autor.

## Bloqueos
Ninguno abierto por el revisor (M-1 sugiere al corrector abrir `docs/estado/bloqueos/WF1-catalogo-pieza-por-kilo.md` si elige la opción B).

## Pendientes
- Ronda 2: verificar M-1 y los bajos tras la corrección.
- S4-11 (lectura de 20 nombres contra `nombre_fuente` en inglés) y el cruce de ids con `nutricion/tabla.csv` y `profeco/vinculos.csv` no se pudieron hacer: esos archivos aún no existen.

## Solicitudes fuera de alcance
Ninguna.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no hay valores operativos de K31 en este artefacto).

## Consumo aproximado
~19 acciones de herramienta · ~150 KB leídos (protocolo 6, decisiones §K 15, formato 19, reporte y README 10, S4-interfaz por partes ~20, criterios S4 18, salidas de scripts ~60) · ~25 KB escritos (scripts ~11, informe 12, reporte 4).
