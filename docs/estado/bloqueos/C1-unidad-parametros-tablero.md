# Bloqueo C1-unidad-parametros-tablero · `unit` de las 5 claves de K54 es una descripción larga, no una unidad (U4 §0.5, U5-06)

Fecha: 2026-10-12 · Agente: `correccion-frontend-cierre-r1` (Opus, WF-3a cierre) · Origen: hallazgo M-1 de `docs/estado/wf3/REVISION-frontend-cierre-r1.md`. Nada de esto cambia cifras ni motores.

## Qué falta (HECHO)
- C1 (`SettingParameterInteger.unit`, extracto `docs/contratos/extractos/ajustes-e.yaml`) dice: «Unidad en español (intentos, minutos, días, MB…)». Es una palabra corta.
- El servidor llena `unit` con la tercera columna de la anotación `@param` de cada migración (`backend/internal/dominio/ajustes.go`, `parseParamDefs`: `Unit: f[2]`). En `backend/db/migrations/00005_tablero.sql:25-29` esa columna es una descripción larga, por ejemplo `% del ingreso mensual neto que se va en pagos de crédito del mes, frontera entre verde y amarillo del barómetro, de 1 a 100`.
- El componente común `Parametro` (`frontend/src/pantallas/ajustes-a/comunes.tsx`, U4 §0.5) pinta `«{valor} {unit}»` dos veces por fila: en «Valor inicial: …» y en el valor. Con el servidor real, P-26.02 «Parámetros del tablero» repetía unas 25 palabras dos veces por fila y decía «de 1 a 100», que contradice el rango dinámico de V-23 («del 1 al 39» para el fin del verde cuando el tope es 40).
- Las otras claves de K44 traen unidades cortas (`intentos`, `minutos`, `días`, `GB`, `kcal`…). Con el mismo patrón largo están `alerts.*`, `reports.*`, `optimizer.*` y `ocr.*` (por ejemplo `hora del centro de México, de 00:00 a 23:59`); quedan fuera de este bloqueo porque no las pinta `Parametro` de P-26.02 y no eran parte de la revisión.
- No se edita una migración aplicada (protocolo; decisión de esquema).

## PREGUNTA
¿De dónde sale la unidad corta de las 5 claves de K54 (`debt.capacity_warning_pct`, `debt.capacity_limit_pct`, `dashboard.payment_reminder_count`, `budget.signal_months`, `budget.signal_margin_pct`)?

| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. La API da una unidad corta para estas claves, **sin editar la 00005**. Por ejemplo, una migración nueva que solo vuelve a anotar las 5 claves con `@param` y unidad corta (`parseParamDefs` ordena los archivos y la última anotación gana), o una tabla de unidades en el dominio. | Cumple lo que C1 ya dice de `unit`; el frontend vuelve a usar `param.unit` y se quita la opción B | Medio: backend y verificaciones de K44 (`make verify-all` revisa el inventario de `@param`) | Bajo |
| B. P-26.02 usa unidades propias en `textos.ts` para las 5 claves (`%`, `%`, `pagos`, `meses`, `puntos`: la primera palabra de la descripción del servidor). **Aplicada como PROVISIONAL mientras se decide.** | Fila corta y sin «de 1 a 100»; nada cambia en la API | Bajo: un mapa en `frontend/src/pantallas/ajustes-a/textos.ts` y la prop `unidad` de `Parametro` | Medio: dos fuentes de la unidad; si la API cambia la clave o la unidad, la interfaz no se entera |
| C. `Parametro` muestra la unidad solo si es corta (por ejemplo, sin coma y de pocas palabras) | Arregla también las otras secciones con el mismo patrón | Bajo | Medio: con estas 5 claves no se vería ninguna unidad («Valor inicial: 30», sin «%»), y el umbral de «corta» sería una regla inventada |

PROPUESTA: A como definitiva (es lo que C1 ya promete) y B solo como puente; al aplicar A se quitan `unidades` de `textos.ts` y la prop `unidad` de las 5 filas de P-26.02.

## Qué se hizo mientras tanto (PROVISIONAL)
- Opción B en `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx` (`ParametrosTablero`), `textos.ts` (`presupuesto.unidades`) y `comunes.tsx` (prop opcional `unidad` de `Parametro`; sin ella, todo sigue usando `param.unit`).
- La prueba `frontend/pruebas/ajustes-a/parametros-tablero.escritorio.spec.ts` ya no inventa unidades: lee `unit` de las mismas anotaciones `@param` de `00005_tablero.sql` que lee el servidor, y comprueba que la descripción larga no aparece en P-26.02.

## Impacto
- Sin decidir, P-26.02 depende de un mapa local de unidades. Ninguna cifra ni regla cambia; V-23 sigue con los rangos de C1 v1.5 (bloqueo `C1-tablero`).
- Decide Axel en la puerta 4 o el orquestador.
