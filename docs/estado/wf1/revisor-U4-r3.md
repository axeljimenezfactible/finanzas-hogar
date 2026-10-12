# revisor-U4-r3 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Ronda 3 de 3 de la revisión independiente de `docs/design/U4-ajustes-a.md` 1.2 y `docs/design/U4-ajustes-b.md` 1.2, solo contra `docs/criterios/U4.md` 1.2 y decisiones §A–§K.
- Veredicto: **NO APROBADO**, con 0 altos, 1 medio (R3-1) y 5 bajos (R3-2 a R3-6). Los hallazgos de la ronda 2 R2-1 (alta), R2-2 (media) y R2-3 (media) quedan **cerrados**.
- Informe con la tabla de criterios, la evidencia propia, los hallazgos y el estado de los anteriores: `docs/estado/wf1/REVISION-U4-r3.md`.

## Archivos
- Creados: `docs/estado/wf1/REVISION-U4-r3.md` y `docs/estado/wf1/revisor-U4-r3.md`.
- Temporales: `datos/tmp/revU4r3/` (ops.py, ops2.py, seed.py, chk.py, orac.py, tok.py, defs.py y ops.json), **borrados** al terminar.

## Verificación
- U4-01: `chk.py` → 14 secciones (7 + 7) con los nombres de la decisión G y 0 celdas vacías en 267 filas.
- U4-02: `ops2.py` contra `openapi.yaml` → 210 operaciones, 73 de Ajustes, **0 sin referencia en U4**. `ops.py` → 107 `operationId` citados, **0 que no existan**.
- U4-03: los ≈ 53 elementos están presentes. El interruptor K21 está en P-26.02 (baja R3-5).
- U4-04: `defs.py` → muestra al azar de 15 valores (semilla 20261011): todos tienen origen en §A–§K o son «sin valor»; 0 inventados.
- U4-05: `orac.py` → las series A y B son iguales a `credito.comparacion_series.esperado_A/B`; también coinciden 800.50, 973.75, las fechas de `despensas[0].pagos`, 99.50, 9-nov-2026, 8-ene-2027, 0.9836, 0.4918, IMC 23.7, metas de 2,556, 2,056 y 1,799, RN-15 1,486.571 y RN-18 50 %. Queda **R3-1 (media)**: la vista previa del gasto fijo solo muestra fechas, no el dinero que cambia.
- U4-06: `chk.py` → 13 atajos, 0 fuera de U1 §5.3; 8 comandos «/», todos en U1 §5.2. Las vistas coinciden con U1 §2.5 y §2.6, salvo inventario y recetario (baja R3-4).
- U4-07: muestra de 15 mensajes en es-MX, sin códigos.
- U4-08: 10 acciones críticas con §0.4. `tok.py` → 21 tokens existen en C2 y en DESIGN.
- U4-09: 0 hex y 0 `rgb(`; las 14 secciones tienen 412 px y 1440 px.
- U4-10: hay estados en las 14 secciones.
- U4-11: 0 fechas ISO; SUPUESTOS en las dos partes; 4 bloqueos `U4-*`; los ejemplos vienen del oráculo o de RN-15, RN-17 y RN-18.
- U4-12: `seed.py` y `grep @param` → las 15 claves de `app_setting` coinciden en clave, sección, valor inicial, decisión y «restablecer». Los 4 de `params_json` llevan el lugar PROVISIONAL (`S1-parametros-avisos` sin respuesta) y coinciden con 00001 (`weekday` 7, `day_of_month` 1).
- U4-13: 2 valores de la casa y 2 por persona, «mantener» sin campo; el ejemplo 2,098.50 − 300 = 1,798.50 → 1,799 está recalculado.

## Supuestos
- Clasifiqué R3-1 como media y no como alta: el gasto fijo sí tiene una vista previa definida (las próximas 3 fechas), pero no muestra las cifras que cambian, como pide la regla propia de U4 en §0.2, y quien implemente tendría que adivinar el contenido. Si el integrador lee U4-05 de forma literal («todo ajuste con efecto numérico»), sería alta.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`, sin «Respuesta de Axel»: no lo decido.

## Pendientes
- Para quien corrija U4 (fuera de mi alcance): R3-1 en `docs/design/U4-ajustes-a.md` P-26.02; R3-2 en P-26.11, P-26.06 y P-26.05; R3-3 en la parte a («Fronteras») y la parte b (línea 3 y P-26.13); R3-4 y R3-5 por solicitud a U1 o a los criterios; R3-6 lo resuelve el integrador.

## Solicitudes fuera de alcance
- `docs/design/U1-arquitectura.md` §2.5.a: alinear el tablero de `inventario_item` (Estado contra «por modo» de §2.3) y el de `receta` («por origen»). Motivo: U1 se contradice y U4 sigue a §2.3 (R3-4).
- `docs/criterios/U4.md` U4-03: decidir si el interruptor K21 se nombra en «Créditos y condiciones» o en «Presupuesto», como U1 §2.3 (R3-5).

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna nueva: no aplica a esta revisión, y K31 no me toca.

## Consumo aproximado
- Unas 41 acciones de herramienta, de un tope de 60.
- Unos 290 KB leídos, de un tope de 400: protocolo, decisiones, criterios U4, las dos partes completas, 5 bloqueos, RN-15 a RN-18, U1 §2.3, §2.5, §2.6, §5.2 y §5.3, y salidas de los scripts.
- Unos 30 KB escritos, de un tope de 150: informe, reporte y scripts temporales.
