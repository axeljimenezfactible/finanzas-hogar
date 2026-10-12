# oraculo-S3 · Reporte final (WF-1, reanudación del 11-oct-2026, 19:00)

Agente: oraculo-S3 (Opus), «agente de regresión S3» de `S3-interfaz.md` 1.4 §10.0. Contrato vigente: ronda 4 (r4-P1). El intento anterior se interrumpió sin reporte; sus parciales se revisaron, se corrigieron contra la versión 1.4 y se completaron.

## Hecho
- `oraculo/cred_lib.py` (parcial revisado y corregido): fórmulas propias de §3–§7 y vistas de §9.3 con `fractions.Fraction`, calendario día por día. Correcciones contra la 1.4: (1) la tasa de una compra cotizada lleva los decimales de **su** `rate_decimals` (null → −1), nunca los de la tabla (§5.1, r3-N2); (2) `rate_display` = `""` con decimales −1; (3) `rate_decimals` < 0 → `invalid`; (4) listas `provisional` de §14.1: `S3-serie-sin-intereses` en toda compra de línea, MSI o tienda con interés 0 y n ≥ 2, `U4-iva-credito` en préstamo/`amortizar` con IVA nulo y en todo interés con fórmula o mínimo de tarjeta/tienda.
- `oraculo/cred_escenarios.py` (parcial revisado): búsqueda exhaustiva propia (candidatos crédito × fecha ± días × tienda × plazo/promoción, restricciones día por día, objetivo 2 con `math.isqrt`, normalización mín-máx en `Fraction`, 4 óptimos con desempates de §8.6 y recomendación). Se añadió: tiempo máximo real (revisa el reloj antes de cada candidato y al terminar; 0 s → `timeout`), funciones `std_milli` y `normaliza`, autoprueba con los ejemplos de §8.4 y §8.5, y el modo sin argumentos de §10.0 (avisa cada fixture faltante del optimizador y sale con 1).
- `oraculo/cred_run_fixtures.py` (parcial revisado): lee fixtures (`$ref`, `parche`), recorta con `esperado` como plantilla (FORMATO §4) y escribe `oraculo/salidas/s3.json` (`salida/1`). Se añadió la **segunda vía en float** (como `verify3.py`): series de pagos, partes, capital repuesto, series A/B, mínimos, sensibilidad e interés por monto, comparadas a 2 decimales con la salida exacta.
- `oraculo/cred_genera_fixtures.py` (nuevo): generador determinista de `_datos/kueski-semilla.json` (receta de «Datos compartidos») y de los **21 casos** de §10 que no son del comparador. Canon: hojas de `model.json` (y `data.json` para R-05/R-06) convertidas con FORMATO §3 leyendo solo esos archivos; sin canon: cifras de la interfaz/criterios/decisiones o cuentas a mano con `decimal.Decimal` (método distinto del oráculo) explicadas en `descripcion`. Nunca llama a `cred_lib` para un `esperado`.
- Fixtures generadas en `oraculo/fixtures/s3/`: kueski-despensa-1..4, feb-2027, series-a-b, series-desde-otra-fecha, reposicion-separada, deuda-viva-10-abr-2027, sensibilidad, interes-por-monto, tasa-plazo, promocion-0, tarjeta-sin-intereses, tarjeta-minimo (PROV), prestamo-frances (PROV, incluye una que obliga ajuste: pago fijo 156.42 y último 156.41), tienda-abonos, tienda-abonos-interes (PROV), unificado-sin-corte (PROV), unificado (PROV), unificado-duplicado.

## Archivos
- Modificados (parciales míos): `oraculo/cred_lib.py`, `oraculo/cred_escenarios.py`, `oraculo/cred_run_fixtures.py`.
- Nuevos: `oraculo/cred_genera_fixtures.py`, `oraculo/fixtures/s3/_datos/kueski-semilla.json`, `oraculo/fixtures/s3/s3-*.json` (21), `oraculo/salidas/s3.json`, este reporte.

## Verificación
- `python3 -I oraculo/cred_genera_fixtures.py` → «21 casos + semilla»; dos corridas dan bytes iguales (sha256 de las fixtures igual).
- `python3 -I oraculo/cred_run_fixtures.py` → «21 fixtures · 73 consultas» y «segunda vía (float, 2 decimales): 174 de 174 cifras coinciden»; dos corridas dan el mismo sha256 de `salidas/s3.json`.
- `python3 -I oraculo/compara.py s3` (código 0): oráculo contra esperado 893 de 893; **model.json (ramas de s3): 437 de 437**; data.json 9 de 9 (R-05, R-06); provisionales 695 de 695; diferencias 0; sin `SIN-DUEÑO`. Motor: «sin motor» (el núcleo S3 aún no existe).
- `python3 -I oraculo/cred_escenarios.py` (código 1, esperado por §10.0): autoprueba §8.4/§8.5 «bien» (0/100/200 → 81,650; 100/200/400 → 1, 2/3, 0; empate → 1); «falta oraculo/fixtures/s3/s3-comparador-…json (autor: optimizador S3)» ×6; «escenarios aceptados con violación = 0».
- Prueba de humo del comparador (en memoria, sin fixture): 2 líneas, 41 candidatos, 21 válidos, 20 descartados por `max_live_purchases`, 65 ms; `timeout_seconds` 0 → `timeout`; restricciones vacías → `missing_constraints` con las 4 en el orden de C1, `evaluated_count` 0, `provisional` [].
- `sha256sum oraculo/model.json oraculo/verify3.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` y `d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041`, sin cambio (S3-16). `grep reporte-compra-latex oraculo/cred_*.py` = 0; ningún `cred_*` lee `datos/tmp/motor-s3.json` ni `motor_cred`.
- Criterios cubiertos por el oráculo (su parte Python): S3-01 (canon de S3 y cruces con S2: 6,293.44, P1–P6, 13,103.78, vivas d3/d4), S3-02 (437/437), S3-03 (series A/B y lectura omitida → B con «pendiente confirmar»), S3-04 (16 fechas del canon y 28-feb / 29-feb), S3-05 (reposición solo al completar el pago, K15), S3-06 (sin tasa a 6–12 q; 1,034.74; decimales K12), S3-07 (4 fronteras), S3-08 (0.00 un día antes y el mismo día con y sin datos; fórmula y mínimo PROV; «falta capturar»), S3-09 (3 tablas francesas con propiedades y ajuste), S3-10 (abonos sin interés, «falta capturar», tienda con interés PROV), S3-11 (4 tipos, orden, suma por quincena, corte que sustituye, duplicado inválido), S3-14 (ejemplo de normalización), S3-16. S3-12/S3-13 quedan pendientes de las fixtures del optimizador.

## Supuestos
- La segunda vía toma las fechas de pago de `cred_lib.paydays_after` (no son dinero) y los valores 1,191.50 y 99.50 de la receta de §10; solo verifica, nunca entra en una salida.
- `normaliza` da 0 a un valor nulo cuando otros no lo son; ese caso es inalcanzable con la generación de §8.2 (con tienda en la petición todos comparten la misma canasta o ninguna; sin tienda, cada candidato sale de una canasta), así que no cambia cifras.

## Bloqueos
- Ninguno nuevo. Las fixtures PROVISIONALES citan los bloqueos abiertos que fija §14.1: `S3-formula-tarjeta`, `S3-convencion-prestamo`, `S3-origen-cargo-normal`, `S3-serie-sin-intereses`, `U4-iva-credito`. `S1-parametros-avisos` no toca a S3.

## Pendientes
- Fixtures del comparador (autor: optimizador S3, §10.0): `docs/estado/S3-escenarios.md`, `oraculo/fixtures/s3/_datos/comparador-caso-{1,2,3}.json`, `oraculo/fixtures/s3/s3-comparador-{caso-1,caso-2,caso-3,sin-restricciones,tiempo,carga}.json`. Cuando existan: `python3 -I oraculo/cred_escenarios.py` (debe dar código 0 y «aceptados con violación = 0»), `python3 -I oraculo/cred_run_fixtures.py` y `python3 -I oraculo/compara.py s3`; una diferencia con la tabla a mano es hallazgo con ambos valores (no se edita la fixture ajena).
- Motor Go S3 (núcleo): al existir `datos/tmp/motor-s3.json`, `compara.py s3` sin `--sin-motor` cierra «motor contra esperado» y «oráculo contra motor».

## Solicitudes fuera de alcance
- Ninguna nueva. Se mantiene la del contrato: que `oraculo/fixtures/FORMATO.md` admita una lista en `provisional` (hoy los bloqueos secundarios van en «otros bloqueos:» dentro de `que`).

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Sin valores K31 nuevos: el tiempo máximo del comparador ya está aprobado (5 s, K43) y las restricciones las captura Axel (K16). Dato informativo: el oráculo en Python evalúa ~40 candidatos en ~65 ms.

## Consumo aproximado
- 41 acciones de herramienta; ~300 KB leídos (S3-interfaz por secciones ~125 KB, criterios ~30 KB, FORMATO 16 KB, decisiones C y K 17 KB, parciales 53 KB, compara 7 KB, RN 12 KB, model.json/data.json en extracto 6 KB, salidas de corridas); ~70 KB escritos con herramientas (generador 42 KB, parches 15 KB, reportes 10 KB); los scripts generaron ~196 KB de fixtures y la salida del oráculo.
