# correccion-datos-r2 · Corrección de los hallazgos M-1 y M-2 de REVISION-datos-r2

Fecha: 11-oct-2026, noche. Alcance: `docs/design/maquetas/` (JSON y LEEME), `docs/design/tools/ds/generar_sinteticos.py` (el JSON es su salida) y `verificar_sinteticos.py`. Sin git add/commit.

## Hecho
- **M-1 (resumen del 20-oct).** `generar_sinteticos.py`: nuevas funciones `cn(f)` y `txt_dia(f, etiqueta)`; el conteo de gastos de los resúmenes diarios sale de `gastos` igual que el monto. El del 20-oct ahora dice «en 2 gastos» (monto $951.30 sin cambio); el del 19-oct sigue en 3 gastos y $645.00. Verificador: comprobación 34 recalcula monto y conteo de cada resumen diario.
- **M-2, camino (a): se sumaron las series al generador, con comprobaciones** (el hallazgo daba a elegir entre (a) y (b); se tomó (a) porque K45.1 manda cifras sintéticas y deja las vistas dibujables):
  1. `ingresos.tendencia` (P-08 N-103): 6 meses, mayo a octubre; octubre = datos de `ingresos`.
  2. `gastos_resumen.tendencia_por_categoria` (P-02 N-044, P-06 N-085): 6 meses × 13 categorías con plan, real y %; cada mes suma el plan y el real de `presupuesto.tendencia`; `seis_mayores` y `otras_*`.
  3. `productos[].historial_precios` (P-15 N-171): 105 series (una por par producto × presentación × tienda), 3 a 5 puntos desde mayo, último punto = precio vigente. Sin flotantes (el precio por unidad base se calcula; no se tocó B-4).
  4. `recetas[].uso` (P-12 N-141): `veces_8_semanas` y `proxima_fecha` coherentes con `plan_semana`.
  5. `ajustes.datos_semilla[].registros` (+ `registros_texto`, `registros_estado`) (P-26.11 N-305): conteos del mismo archivo (33 productos, 5 compras de Kueski, 10 registros por importar).
  - `indice_pantallas` apunta a las llaves nuevas (P-02/P-03, P-06, P-08, P-10/11/12, P-15). `meta.version` = 1.2 (alinea con el LEEME; era B-5, tocado de paso).
- **LEEME 1.2:** filas nuevas por llave, regla de «menos de 12 meses» ampliada a todas las series (6 meses; el historial de precios no se interpola), nota de textos que repiten cifras, 57 comprobaciones.
- **Verificador:** de 51 a 57 comprobaciones (34 a 39 son nuevas); `existe()` entiende rutas `lista[].llave`.

## Verificación
```
python3 -I docs/design/tools/ds/generar_sinteticos.py docs/design/maquetas/datos-sinteticos.json
(segunda generación a /tmp/regen.json) cmp                         -> REGEN-IDENTICA (determinista)
python3 -I docs/design/tools/ds/verificar_sinteticos.py             -> 57 OK, 0 FALLA, RESULTADO: OK
prueba negativa (JSON copiado con «en 3 gastos» el 20-oct)          -> FALLA 34: «resumen del 2026-10-20 dice ('$951.30', 3) y gastos da ('$951.30', 2)»
python3 -I docs/design/tools/ds/verificar_ds.py                     -> RESULTADO: OK · 445 tokens C2 ... 0 fallas
grep de cifras prohibidas P3-12 (JSON y LEEME)                      -> 0; sin correos (@) en el JSON
```
Las fechas nuevas (mayo a septiembre) no coinciden con `model.json` ni `data.json` (solo siguen las 15 de la regla de quincena); no caen en día 9, 10, 15, 28, 30 ni 31.

## Supuestos
- Ingresos de mayo a septiembre: quincenas = esperado de octubre × 0.955 a 0.974 (antes de `rige_desde`), extras en 2 de 5 meses (junio y julio); todo inventado.
- Categorías en meses previos: vivienda = plan; servicios entre 55 % y 66 % del plan; seguros en 0 (como en octubre); las variables reparten el real restante con ±12 % de ruido. Plan por categoría: las fijas igual que octubre, las variables proporcionales al plan de octubre.
- `seis_mayores` = las 6 de mayor plan (el hallazgo no fija el criterio de «mayores»; U1 dice «las 6 mayores»). Si Axel prefiere por gasto real, se cambia solo aquí.
- Historial: los puntos previos varían entre −6 % y +8 % por paso; la fuente es de ejemplo.
- `uso` de recetas: las 7 semanas previas son inventadas (0 a 5 veces); la semana actual sale del plan. La próxima fecha cuenta desde hoy inclusive.
- Los registros de Datos semilla son conteos de este archivo, no del seed real.

## Bloqueos
Ninguno.

## Pendientes
- Siguen abiertos los hallazgos bajos B-1 a B-5 (salvo la versión): campos sin cifras (N-060, N-128, N-274, N-294, N-317, N-011), horizonte del calendario, $9,000.00 en fechas coincidentes, `precio_unidad_base_cents` flotante en `precios`, `person_id` y `store_id` como texto y «Pendiente de importar» contra «Pendiente».
- La regla «menos de 12 meses» obliga a los autores a dibujar 6 meses vacíos en las vistas de 12 meses.

## Consumo aproximado
~25 acciones de herramienta; ~125 KB leídos (reporte de revisión, LEEME, generador y verificador por partes); ~70 KB escritos contando el JSON regenerado (252 KB de salida del generador, 42 KB más que antes), el código nuevo (~9 KB), el LEEME y este reporte.
