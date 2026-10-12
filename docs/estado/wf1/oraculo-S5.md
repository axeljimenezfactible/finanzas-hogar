# oraculo-S5 · Oráculo independiente de S5 (precios y tiendas)

Agente `oraculo-S5` (Sonnet) · WF-1 · 11-oct-2026 · reanudación después del reinicio de las 19:00 (el intento anterior no dejó archivos; se partió de cero). Contrato leído en su versión vigente: `docs/contratos/motores/S5-interfaz.md` 1.4 (ronda 4).

## Hecho
- Oráculo en Python 3.14 con `-I` y solo biblioteca estándar (`json`, `csv`, `fractions`, `decimal`, `datetime`), escrito desde el contrato, sin leer `backend/internal/precios`. Las 12 funciones consultables de §15.2: `base_unit_price`, `freshness`, `current_price`, `quote_for`, `compare_carts`, `compare_basket`, `membership`, `next_visit`, `stale_prices_for_plan`, `profeco_import`, `ocr_parse_tsv`, `assisted_can_enable`.
- `Input` de §15.1 con `unknown_field`, `missing_field`, `invalid_date` y los errores de ingesta (`future_date`, `unknown_presentation`, `unknown_store`, `duplicate_price`); `report_file` arma los carritos `bodega` y `walmart` desde `data.json`.
- 27 fixtures `fixture/1` en `oraculo/fixtures/s5/` (los 13 casos de §15 más 14 de apoyo: frescura con parámetros K44, promoción, 4 empates/incompleta de canasta, membresía anual y por canasta, 3 variantes de lectura asistida, 4 errores de ingesta) y 4 archivos de `_datos/` (2 CSV de PROFECO, uno en UTF-8 con BOM y otro en Windows-1252, `vinculos` y un TSV de OCR).
- Los `esperado` salen de literales calculados a mano o de `data.json` leído directo en el generador; el generador no llama al oráculo (segunda derivación).
- `oraculo/salidas/s5.json` (formato `salida/1`), determinista (dos corridas, mismo sha256 `d73b7809…`).
- Verificación independiente `oraculo/pre_verify.py`: suma `data.json` renglón por renglón con `Decimal` y contrasta con la salida del oráculo.

## Archivos
- `oraculo/pre_comun.py`, `oraculo/pre_datos.py`, `oraculo/pre_calculo.py`, `oraculo/pre_profeco.py` (PROFECO y OCR), `oraculo/pre_run_fixtures.py` (ejecutor; el que nombra §15.3 como `s5_run_fixtures.py`), `oraculo/pre_genera_fixtures.py`, `oraculo/pre_verify.py`.
- `oraculo/fixtures/s5/*.json` (27), `oraculo/fixtures/s5/_datos/{profeco-muestra-a.csv,profeco-muestra-b.csv,profeco-vinculos.csv,ocr-muestra.tsv}`.
- `oraculo/salidas/s5.json`, `docs/estado/wf1/oraculo-S5.md`.

## Verificación
- `python3 -I oraculo/pre_genera_fixtures.py` -> `27 fixtures en …/oraculo/fixtures/s5`.
- `python3 -I oraculo/pre_run_fixtures.py` -> `27 fixtures, 27 coinciden, 1058 hojas de esperado, 0 diferencias` (la hoja incluye los largos `.#` de las listas).
- `python3 -I oraculo/compara.py s5 --sin-motor` -> `oráculo contra esperado: 747 de 747 hojas coinciden`, `provisionales: 225 de 225`, `model.json (ramas de s5): 0 de 0`, `diferencias: 0`, código 0. El motor Go de S5 aún no está comparado: falta `datos/tmp/motor-s5.json` (pendiente).
- Control negativo: con una copia de las fixtures en el scratchpad y un esperado alterado (2201 en lugar de 2200) el ejecutor imprime `DIF s5-cotizacion-recomendada canasta.stores[1].total_cents esperado=2201 oraculo=2200` y sale con 1.
- **Los dos valores canon desde `data.json`** (`python3 -I oraculo/pre_verify.py`, código 0, «TODO COINCIDE»):
  - Idénticos: Bodega 2,773.50 y Walmart 2,874.00 (Σ de 33 renglones, precio unitario x cantidad, 55 piezas). Sustitutos: Bodega 300.00 (66.00 + 99.00 + 135.00) y Walmart 359.00 (84.00 + 105.00 + 170.00). Sexta lata de atún de Bodega: 21.00.
  - Total Bodega = 2,773.50 + 300.00 + 21.00 = 3,094.50; total Walmart = 2,874.00 + 359.00 = 3,233.00.
  - **+138.50** (carrito contra carrito) = 3,233.00 - 3,094.50.
  - **+100.50** (mismos artículos, mismas cantidades) = 2,874.00 - 2,773.50.
  - Explicación: +100.50 (idénticos) + 59.00 (sustitutos: 359.00 - 300.00) - 21.00 (la lata que solo compró Bodega) = +138.50.
  - También coinciden K-04 a K-08 (61 / 60 / 55 piezas, 33 / 3 renglones) y ganadores 22 / 8 / 3.
- Criterios cubiertos por el oráculo (la cobertura de pruebas Go depende del motor):
  - S5-01 `s5-comparativa-canon`: K-01 a K-08 y ganadores.
  - S5-02 `s5-unidad-sustitutos`: 2933/2800, 12692/16406, 3375/2833; HALF_UP en mitades exactas (15497.5 -> 15498; 537.5 -> 538).
  - S5-03 `s5-mejor-tienda-canasta` (Bodega 3,094.50 es la más barata), empates, incompleta, `s5-cotizacion-recomendada`, `s5-canasta-regla-unica` (Walmart 9430 contra Bodega 9460).
  - S5-04 `s5-frescura-fronteras` (0.9836 / 0.4918 / 0, día 91 no compite, recta de 1/61 día por día de 30 a 90), `s5-frescura-parametros`, `s5-fuentes-compiten`.
  - S5-05 y S5-06 solo el parseo (`s5-profeco-muestra`, `s5-ocr-parseo`); la cobertura real y la medición con fotos no son del oráculo.
  - S5-07 `s5-lectura-asistida*`; S5-09 `s5-membresia-*`; S5-10 ingesta y `s5-promocion-vigencia`; S5-11 `s5-mercado-visita`; S5-12 `s5-precios-viejos`.
- Hallazgos de diferencia contra cifras canon: **ninguno** (`oraculo/model.json` no tiene rama de S5; contra `data.json` coincide todo).

## Supuestos
El contrato no fija estas convenciones y el motor Go debe coincidir con ellas, o el integrador las corrige en el contrato. Ninguna cambia una cifra canon.
1. `compare_carts`: `stores[]` por total ascendente (empate: orden base y luego `others`); `winners[]` lista todas las tiendas en ese orden, también con `count` 0.
2. `rows[]` de `compare_carts`: primero los idénticos (en orden de aparición en la canasta base), luego los sustitutos (por producto, orden de aparición) y al final los «solo en una tienda» (por canasta y renglón). El excedente de una presentación idéntica (la sexta lata) sale como `only_one` aparte. Por eso `canon-filas` va en una consulta separada de `canon-totales`.
3. `winner_store` es `null` (no ausente) en empate y en `only_one`. Un renglón idéntico es una presentación presente en todas las canastas; un sustituto, un producto con presentaciones sin emparejar en dos o más.
4. `stale_count` de `compare_carts`: número de renglones agrupados del carrito si su `observed_date` tiene edad >= `stale_after_days`; si no, 0.
5. `compare_basket`: se listan todas las tiendas de `stores` (un mercado sin el grupo del producto cuenta ese renglón como faltante); orden por total ascendente, también las incompletas; `difference_vs_cheapest_cents` contra la más barata **completa** (puede ser negativa en una incompleta); `cheapest_store` es `null` si ninguna es completa; `piece_count` aquí es la suma de paquetes completos (no se compara en ninguna fixture); `recommended_presentation` nula en la lista deja `is_substitute = false`; sin cotización: `recommended_store`, `recommended_presentation` y `estimated_price_cents` nulos e `is_stale_price` falso.
6. `current_price`: `price` es el último regular de una fuente que compite, sin importar su edad (la edad la decide `quote_for`); `references[]` son todos los regulares de `site_assisted` o `profeco`, del más reciente al más antiguo; `promotion` es la promoción de fuente que compite con `observed_date` = hoy (PROVISIONAL, bloqueo `S5-reglas-comparativa`, pregunta 3). `quote_for` y `stale_prices_for_plan` solo consideran precios con `observed_date` <= la fecha de evaluación.
7. `freshness`: edad > 90 da peso 0/1, `weight_display` «0.0000», `competes` falso; fecha futura es `future_date`.
8. `membership`: comparadores = tiendas con `requires_membership` falso; un renglón sin precio en Sam's o sin comparador suma 0; compras al año = compras en tiendas que no son mercado con 0 <= edad < 365 días (hoy y los 364 anteriores); `savings_per_purchase_cents` se copia a `basket_savings_cents`; con costo ausente el veredicto es `missing_fee` aunque el ahorro sea 0.
9. `assisted_can_enable`: el sitio de `terms[].site` se liga a la tienda por **igualdad de `code`** (`walmart`, `bodega_aurrera`, `sams_club`): el contrato no dice cómo se relaciona una tienda con un sitio.
10. `next_visit`: una tienda sin `visit_weekday` es `invalid_argument`.
11. `profeco_import`: filtros por igualdad exacta tras quitar espacios (cadena = `profeco_name` de las tiendas `chain`; municipio = `municipality` de cualquier tienda); tienda del precio = sucursal (`branch`) cuyo `profeco_name` = `nombre_comercial`, si no la cadena; clave de liga = (producto, presentacion, marca) exactos; liga a una presentación inexistente va a `skipped`; orden de revisión de cada renglón: rechazo, filtro, liga, duplicado, carga; `read` = rechazados + filtrados + omitidos + cargados + duplicados; `row` = ordinal 1 de los registros de datos de cada archivo; `coverage[]` cuenta lo cargado y va ordenado por cadena y municipio; `rejections[]` solo trae rechazos de lectura; los `reason` son texto libre y no se comparan; fechas aceptadas: `AAAA-MM-DD`, `AAAA/MM/DD`, `DD/MM/AAAA`, con hora opcional. Provisional, bloqueo `S5-profeco-almacen`.
12. `ocr_parse_tsv`: se agrupan las palabras de nivel 5 con texto; cantidad = primer token numérico (con la unidad `kg`, `g`, `l`, `ml`, `pz`… pegada o separada, que se descarta); importe = último token con dos decimales y miles con coma; si el primer token es el propio importe no hay cantidad.

## Bloqueos
Ninguno nuevo. Las fixtures `provisional` citan `docs/estado/bloqueos/S5-reglas-comparativa.md` (promoción vigente, desempates, anualización y redondeo hacia arriba de la membresía) y `docs/estado/bloqueos/S5-profeco-almacen.md` (conteos y cobertura de PROFECO).

## Pendientes
- Correr `datos/tmp/motores s5 dump --fixtures oraculo/fixtures/s5 --out datos/tmp/motor-s5.json` y luego `python3 -I oraculo/compara.py s5` cuando el motor Go de S5 exista, y resolver las diferencias de las convenciones de Supuestos.
- Mapa de «sitio» a tienda (supuesto 9) y orden de `rows[]` (supuesto 2): confirmar o cambiar en el contrato.

## Solicitudes fuera de alcance
- `docs/contratos/motores/S5-interfaz.md` §15.2: escribir las convenciones 1 a 12 de Supuestos (orden de `stores[]`, `rows[]` y `winners[]`; `null` por «ausente»; `piece_count` de `compare_basket`; relación tienda-sitio en `assisted_can_enable`; formatos de fecha y filtros exactos de PROFECO; reglas del parseo de OCR).
- `docs/contratos/motores/S5-interfaz.md` §15.3 nombra `oraculo/s5_run_fixtures.py`; el ejecutor quedó como `oraculo/pre_run_fixtures.py` (nombre de la tarea). Pedir al integrador que ajuste la mención o renombre el archivo.
- `oraculo/compara.py` ya trae la rama `s5`; no se tocó.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: este agente no necesitó valores operativos nuevos de K31 (los umbrales 30/90 de frescura salen de `app_setting` y se prueban con otros valores en `s5-frescura-parametros`).

## Consumo aproximado
- Acciones de herramienta: 40 (menos de 60). Lectura: unos 175 KB (contrato S5 54 KB, criterios 19 KB, FORMATO 16 KB, protocolo 6 KB, formato del catálogo 10 KB, referencias del oráculo S4 25 KB, `data.json`). Escritura: unos 110 KB de código y documento a mano; las fixtures y `s5.json` los genera el script.
