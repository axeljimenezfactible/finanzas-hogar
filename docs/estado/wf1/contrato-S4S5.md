# contrato-S4S5 · reporte final (WF-1, 11-oct-2026, Opus)

## Hecho
- `S4-interfaz.md`: paquete `internal/recetas` (funciones puras, sin flotantes, un redondeo HALF_UP al final), esquema de `recetario.json` (60 recetas, llave `recetas`, esquema en `datos/semilla/recetario.schema.json`), escalado exacto con fracciones, plan semanal (4 tiempos, «fuera», Karen rechazada, domingo abre la semana siguiente K38, huecos K43), lista neteada (Σ escalado − existencia + reposición, paquetes, sobrantes), inventario híbrido (RN-15, estado K7), costo por receta, celda, día y semana con precios de S5 por `PriceSource`, kcal con fuente, `tabla.csv` y su cargador, IMC y meta Mifflin-St Jeor × actividad ± ajuste (K23, K36), fixtures `s4-*`, pruebas por criterio, tabla criterio → sección.
- `S5-interfaz.md`: paquete `internal/precios` y subpaquetes (`profeco`, `ocr`, `asistida`), 5 fuentes (3 compiten, 2 referencia: K22, K37), unidad base y sustitutos (RN-20, K7), frescura (K40, parámetros K44), precio vigente, comparativa en dos modos (compras registradas = canon; canasta por necesidad), mejor tienda por ingrediente y canasta, membresía de Sam's, mercados, precios viejos para S7, cargador PROFECO, OCR Tesseract 5.5.3 spa con confirmación humana, lectura asistida apagada, fixtures `s5-*`, tabla criterio → sección.
- `productos-formato.md`: formato de `productos.json` (grupos con canal, productos con id estable, unidad base, vertical, sinónimos, presentaciones con factor y `nombres_ticket`), 36 productos y 39 presentaciones de `data.json`, y la lista de los 39 renglones (61 piezas) con su nombre en el catálogo.
- Alineado con `oraculo/fixtures/FORMATO.md` 1.0 (apareció durante el trabajo; lo escribió `contrato-S2`): `fixture/1`, `salida/1`, `backend/cmd/motores s4|s5 dump --fixtures --out`, `compara.py`.
- 3 bloqueos nuevos con opción PROVISIONAL citada en los contratos.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/contratos/motores/S4-interfaz.md` (34.6 KB)
- `/home/deepsec/finanzas-hogar/docs/contratos/motores/S5-interfaz.md` (27.1 KB)
- `/home/deepsec/finanzas-hogar/docs/contratos/motores/productos-formato.md` (18.8 KB)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S4-reglas-lista.md`, `S5-reglas-comparativa.md`, `S5-profeco-almacen.md`
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/contrato-S4S5.md` (este reporte)

## Verificación
- Canon de la comparativa (python3 -I, `Decimal` y `Fraction` sobre `data.json`): 0 renglones con unitario × cantidad ≠ total; idénticos 277350 / 287400 → +10050 (**+100.50**); sustitutos 30000 / 35900; atún extra 2100; Bodega **309450**, Walmart **323300**, diferencia **13850** = 10050 + 5900 − 2100 (**+138.50**); piezas 55 / 61 / 60; renglones 33 / 3; ganadores 22 / 8 / 3 = `resultados_por_producto`. Unidad base: 2933 vs 2800, 12692 vs 16406, 3375 vs 2833 (mismo ganador que las notas).
- Frescura con el ticket del 10-oct-2026: 29 d (2026-11-08) peso 1; 30 d 60/61 = 0.9836; 31 d 59/61; 60 d 30/61 = 0.4918; 90 d (2027-01-08) 0; 91 d no compite.
- Salud: MSJ 1748.75 y 1330.25; × 1.2 = 2098.5 → 2099 y 1596.3 → 1596; −300 → 1799; +250 → 1846; RN-17 2555.5625 → 2556, con −500 → 2056; IMC 74/100/120 kg a 2 m = 18.5 / 25.0 / 30.0; 72.5 kg y 1.75 m → 237 (23.7). RN-15: 2838000 − 1351429 = 1486571. RN-16: martes siguiente al 10-oct = 2026-10-13. Costo ilustrativo §8.3: 4425 / 2213 / 8850.
- Catálogo (`catalogo.py` en el scratchpad): 36 productos, 39 presentaciones, 0 nombres de `data.json` sin presentación, 39 renglones de `cuentas` = **61 piezas** = **3,094.50**, piezas por presentación = `items_identicos[].cantidad` (atún 6 = 5 + 1).
- Documentos: bloques JSON válidos (3 de 3); 0 líneas PROVISIONAL sin bloqueo; S4-01…S4-15 y S5-01…S5-13 presentes; S4 §15 y S5 §18 = tabla criterio → sección; tamaños ≤ 40 KB.
- Entorno: `tesseract 5.5.3`, idiomas `osd spa`; `unrar` y `bsdtar` presentes. `datos/semilla/tickets/` no existe → S5-06 «sin medir».

## Supuestos
- S-1 «Fuera» excluye recetas en la misma celda (S4 §5.1). S-2 Celda con ingredientes de varias tiendas: sin `price_store_id` y fecha más vieja (S4 §5.3). S-3 Ingredientes opcionales fuera de lista, costo y kcal (S4 §5.4). S-4 `days_left` = días completos (piso) (S4 §7). S-5 Dato nutricional de la presentación antes que el genérico (S4 §9.3). S-6 retirado (FORMATO.md ya existe). S-7 Fecha de precio futura = error (S5 §3). S-8 Misma fecha: rige la última captura (S5 §6). S-9 `winners`/`ties` cuentan solo idénticos, como `data.json` (S5 §7.2). S-10 Día de visita: hoy cuenta (S5 §9). S-11 `needed_milli` hacia arriba a milésima entera (S4 §6.2). Ninguno cambia cifras canon.

## Bloqueos
- `docs/estado/bloqueos/S4-reglas-lista.md` (sobrante previsto, reposición por ritmo, clasificación del IMC).
- `docs/estado/bloqueos/S5-reglas-comparativa.md` (aplicación del peso, empates y canasta incompleta, vigencia de promociones, anualización y redondeo de la membresía).
- `docs/estado/bloqueos/S5-profeco-almacen.md` (tabla de registros crudos, liga con el catálogo, duplicados).
- Sigue abierto y no lo toco: `S1-parametros-avisos.md`.

## Pendientes
Ninguno de mi alcance. Para otros agentes: `datos/semilla/productos.json`, `recetario.json` y `recetario.schema.json`, `nutricion/tabla.csv` y `PROCEDENCIA.md`, `profeco/vinculos.csv`, `docs/estado/S5-terminos.md`, `docs/estado/S5-cobertura-profeco.md`, el código de `internal/recetas` e `internal/precios`, fixtures `oraculo/fixtures/s4|s5/`. Axel: fotografiar 5 tickets en `datos/semilla/tickets/`.

## Solicitudes fuera de alcance
1. `docs/contratos/motores/S3-interfaz.md` línea 24 y §9: usa `go run ./cmd/finanzas s3 dump --fixture`; FORMATO.md §4 fija `backend/cmd/motores <motor> dump --fixtures --out`. Alinear.
2. `oraculo/fixtures/FORMATO.md` §2: admitir `$ref` fuera de `_datos/` (S5 lee `../../../datos/semilla/data.json` sin copias, S5-01).
3. Esquema 1.4 y C1 (S1/S6): clave `app_setting ocr.timeout_seconds` (K44) para el tiempo máximo de Tesseract.
4. C1 `ProfecoImportResult`: sumar `read_count`, `rejected_count`, `filtered_out_count`, `duplicate_count` para que S5-05 muestre el conteo sin pérdidas.
5. `docs/estado/bloqueos/INDICE.md`: indexar los 3 bloqueos nuevos.
6. `docs/contratos/reglas-negocio.md` tras la puerta 2: escribir la regla de desempate (S5-03 la exige en RN) y la de clasificación del IMC (S4-07 la exige en RN-17).
7. Si `internal/dominio` publica un helper HALF_UP, unificar con el local de recetas y precios.

## Paquetes instalados
Ninguno (tesseract, unrar y bsdtar ya estaban).

## Propuestas para Axel
- `ocr.timeout_seconds` = 30 (K31/K44): Tesseract tarda unos segundos por foto en CPU; 30 s cubre fotos grandes sin colgar la petición; configurable.
- Bloqueos: recomiendo la opción A en S4-reglas-lista (las 3 preguntas), A en S5-reglas-comparativa 1 a 3 y 4b, B en 4a (frecuencia por historial real, sin inventar), y A en S5-profeco-almacen (no infla la base de producción).
- Grupos iniciales del catálogo (9, `productos-formato.md` §4), editables en Ajustes.

## Consumo aproximado
51 acciones; ~210 KB leídos; ~150 KB escritos (incluye una reescritura de 62 KB de S4 y S5 para alinear con FORMATO.md, que apareció a mitad del trabajo).
