# Reporte P1 · criterios-S-1.1 (11-oct-2026)
Agente del flujo «P1 · Aplicar respuestas de la puerta 1», modelo Sonnet. Alcance: `docs/criterios/S1.md` a `S8.md`, `WF4-operacion.md`, filas S y WF4 de `INDICE.md`, `docs/criterios/tools/`.

## Hecho
- S1 a S8 y WF4 pasan a **versión 1.1 (11-oct-2026, puerta 1)** con la lista «Cambios de la versión 1.1» antes de §1. No cambió el número de criterios ni ninguna severidad (134 en estos 9 archivos: 105 alta y 29 media).
- Respuestas de §K aplicadas: K22 (S1-03, S5-04), R2-07 ejemplos «ilustrativos» (S1-03), K19 (S2-07, con las 7 filas calculadas desde `model.json`), K20 (S2-08), K21 (S2-09), K18 (S3-08, S3-10), K16 (S3-12), K31 (S3-12 tiempo máximo; S6-09, S6-10, S7-04, S7-08, S7-10, S8-02, WF4-04, WF4-08), K17 y K13 (S3-13, S3-14), K15 (S3-05), K14 (S2-05), K23 y K7 (S4-07, S4-04), K25 (S4-15, S5-13, S6-02), K24 (S5-06), K7 (S5-02, S5-09), K28 (S6-01), K11 (S6-06, S6-07, S1-14), K29 (S6-14, WF4-03), K30 (S7-06, S7-07, WF4-01), K4 (S7-07, S7-09, WF4-05), K6 y K9 (S6-04, S6-05), K12 (S1-10), K6 y K7 sin marca (S1-16).
- Toda ruta `/tmp` pasó a `/home/deepsec/finanzas-hogar/datos/tmp/` con `mkdir -p` donde el comando escribe. S6-03 cita `build.py --check` y `verify_c1.py --db`.
- §5: las paráfrasis entre comillas de S2, S3, S5, S6, S7 y S8 son ahora citas textuales (o fragmentos con «…» en orden); `citas.py` no reporta ninguna cita sin fuente en S1–S8 ni WF4. Cada cambio de §K lleva su línea de trazabilidad.
- `docs/criterios/tools/` creada con los 8 scripts del anexo de `correccion-criterios-r1.md` (rutas ajustadas: `s109_sql.sh` ya no usa el scratchpad; `verif_final.sh` corre desde su carpeta) más `cuenta.py` (conteo por archivo y severidad). Todos corren.
- INDICE.md: filas S1–S8 y WF4 recalculadas desde los archivos (sin cambio de cifras), total de la tabla recalculado y nota de fecha en la línea 3.
- Dos bloqueos nuevos por ambigüedades que K no resuelve (ver Bloqueos).

## Archivos
Modificados: `docs/criterios/S1.md`, `S2.md`, `S3.md`, `S4.md`, `S5.md`, `S6.md`, `S7.md`, `S8.md`, `WF4-operacion.md`, `INDICE.md` (solo filas S/WF4, total y línea 3).
Creados: `docs/criterios/tools/{traza.py,refs.py,citas.py,canon.py,s109_sql.sh,s109_equiv.py,matriz_c102.py,verif_final.sh,cuenta.py}`, `docs/estado/bloqueos/criterios-promedios-redondeo.md`, `docs/estado/bloqueos/criterios-frescura-extremos.md`, este reporte.

## Verificación
- `bash docs/criterios/tools/verif_final.sh` (11-oct-2026 11:35):
  - `cuenta.py` sobre los 9 archivos: S1 17 (12/5), S2 16 (13/3), S3 20 (17/3), S4 15 (11/4), S5 13 (10/3), S6 16 (12/4), S7 14 (12/2), S8 11 (8/3), WF4 12 (10/2) = **134 (105 alta / 29 media / 0 baja)**; Verifica = Evidencia = criterios en los 9; 0 criterios sin línea en §5; incompletos: ninguno.
  - `traza.py` (18 archivos): secciones 1–5 completas, 0 sin traza. `refs.py`: 0 IDs citados que no existan.
  - `citas.py`: 253 citas verificadas, 5 sin fuente, **todas en C1, C2, U2 y U5** (no mías); 0 en S1–S8 y WF4.
  - `canon.py`: 178 comprobadas (R-01…R-85 en S2 y S3, K-01…K-08), 0 diferencias; sha256 de `model.json` 1a57066689dd…bfd4.
  - `s109_sql.sh`: Up OK, 54 tablas, 18 filas (6 meses + 12 quincenas), ventana 10→9 rechazada por CHECK, febrero paga el 28, integrity_check ok. `s109_equiv.py`: 12 de 12 cifras, 0 eventos entre los días 1 y 9. `matriz_c102.py` (con `uv run --with pyyaml`): 199 N-xxx, 197 operaciones, 0 sin operación ni Local.
  - Rutas `/tmp/` fuera de `datos/tmp/` en S1–S8 y WF4: 0.
- Cifras de S2-07 y S2-08 calculadas con Python desde `oraculo/model.json`: deuda viva en los 6 cierres y el 10-abr = 2,921.25 / 973.75 / 3,776.06 / 1,258.70 / 3,776.06 / 1,258.70 / 6,293.44 con saldo no devengado 1,947.50 / 0 / 2,517.37 / 0 / 2,517.37 / 0 / 5,034.74; diarios 31.41, 64.92, 89.91.
- Totales calculados para INDICE: filas S y WF4 = 134 (105 / 29), sin cambio. Total de la tabla = **248 (189 / 59)** (sin cambio). Hoy los archivos U/C suman más que sus filas (C1 21 = 19/2 y U1 16 = 10/6, por el otro agente); con esas cifras el total real sería **256 (196 / 60)**: el integrador reconcilia.

## Supuestos
- La etiqueta «1.1» ya existía (ronda 1 de WF-0, 10-oct); la tarea pide «1.1», así que la reemplacé por la de la puerta 1 y lo dejé escrito en la primera viñeta de cada lista de cambios.
- S2-09: «compromisos contratados» lo glosé como los pagos programados de los créditos ya tomados; el monto de las compras proyectadas con el interruptor encendido lo fija `reglas-negocio.md` (bloqueo `S2-proyeccion-interruptor` si falta fuente).
- S4-07: las constantes de Mifflin-St Jeor (10, 6.25, 5, +5 hombres, −161 mujeres) son las publicadas; los ejemplos usan personas ficticias. Los factores intermedios de actividad y la magnitud del ajuste por objetivo los escribe el agente de S4 con fuente (bloqueo `S4-meta-ajustes` si no hay).
- S5-09: «compras al año de equilibrio = costo anual ÷ ahorro de una compra con membresía» es mi lectura de K7; el ejemplo X = 1,000.00 y Y = 125.00 es de fixture, rotulado ilustrativo.
- S3-14: el caso máximo = mínimo de la normalización lo declara el agente de S3 (bloqueo `S3-normalizacion-empate` si no).
- Apliqué también K4, K6, K9, K12, K13, K14 y K15 en los criterios S que contradecían o ignoraban esas respuestas, aunque la lista de la tarea no las nombraba. Nada de eso cambia cifras ni cuenta criterios.
- S6-06: cómo se combinan usuario e IP queda a lo que escriba `C1-datos.md` v1.1 (K11 no lo precisa); la frontera exacta de los 15 min no se prueba.

## Bloqueos
- `docs/estado/bloqueos/criterios-promedios-redondeo.md`: promedio semanal con el diario sin redondear (219.88 en octubre) o ya redondeado (219.87). S2-08 marca el semanal PROVISIONAL con la opción (a).
- `docs/estado/bloqueos/criterios-frescura-extremos.md`: la recta de K22 vale 1 en el día 30 (lectura literal) o parte del 1 del día 29. S5-04 prueba lo común a ambas lecturas y deja los días 30–89 PROVISIONAL.
- Cerrados por K y ya reflejados en los criterios: `criterios-bind`, `criterios-cobertura`, `criterios-credito-formulas`, `criterios-frescura-peso`, `criterios-invariante-deuda`, `criterios-login-bloqueo`, `criterios-meta-kcal`, `criterios-ocr-tickets`, `criterios-optimizador`, `criterios-promedios`, `criterios-proyeccion-12m`, `criterios-token-servicio`, `criterios-umbrales-ausentes` (K31: PROVISIONAL con propuesta escrita hasta su puerta), `S6-version-go`, `S1-*` de los que tocan S.

## Pendientes
- Ninguno dentro del alcance. Para el integrador: ver solicitudes.

## Solicitudes fuera de alcance
- `docs/criterios/INDICE.md` §6 (tabla de bloqueos): marcar RESUELTOS los bloqueos de `criterios-*` y `S6-version-go` y agregar las filas de los dos bloqueos nuevos (S2-08 y S5-04); §9: registrar este flujo. Motivo: solo se me asignaron las filas S/WF4 y los totales.
- `docs/estado/bloqueos/INDICE.md`: agregar `criterios-promedios-redondeo` y `criterios-frescura-extremos` a la lista de preguntas para Axel.
- `docs/criterios/C1.md` (C1-15): quitar la exigencia de marca a los valores aprobados de K6 y K7, igual que S1-16 (otro agente).
- `INDICE.md` filas U, C, P3: reconciliar con C1 (21) y U1 (16) y el total real (256 si no hay más cambios).
- Los criterios C1, C2, U2 y U5 todavía traen 5 citas «…» que `citas.py` no halla en `decisiones.md` (paráfrasis): corregirlas con citas textuales.
- `plan.md` o `decisiones.md` (K22): aclarar el valor del día 30 (bloqueo `criterios-frescura-extremos`).
- Entregables de otros agentes que los criterios ahora exigen: `reglas-negocio.md` con ejemplos «ilustrativos» rotulados (S1), fórmula de promedios y de proyección con interruptor (S2), fórmula general de tarjeta y crédito de tienda (S3), meta de calorías con ajuste por objetivo (S4), `datos/semilla/tickets/` con las 5 fotos de Axel (S5, acción de Axel), `getfacl`/ruta de recorrido de `~/.openclaw/` para el usuario `finanzas` (WF4-01: la ACL de archivo necesita permiso de recorrido en las carpetas del camino).

## Paquetes instalados
Ninguno (`uv run --with pyyaml` descarga pyyaml en caché de uv para `matriz_c102.py`, sin instalar paquetes del sistema).

## Consumo aproximado
Acciones de herramienta: ~40. Leído: ~220 KB (protocolo, §K, 9 criterios, colas de 24 bloqueos, anexo de scripts, cabecera y tabla de INDICE, salidas de verificación). Escrito: ~95 KB (scripts de parche ~55 KB, scripts de `tools/` ~25 KB extraídos del anexo y nuevos, bloqueos y reporte ~15 KB). Temporales de `datos/tmp` borrados.
