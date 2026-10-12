# revisor-frontend-areas-r5 · WF-3a (12-oct-2026)

Agente: `revisor-frontend-areas-r5`, revisor independiente en la ronda 5. Revisé front-creditos-salud-reportes, front-ajustes-a, front-ajustes-b y front-variantes contra U5-06, U5-07, U5-09, C3-03, DESIGN.md, U1 y U4. No corregí nada. No hice git add ni commit. El informe completo está en `docs/estado/wf3/REVISION-frontend-areas-r5.md`.

## Hecho
- **Veredicto: NO APROBADO**, con 0 altos, 1 medio y 2 bajos.
- **R4-01 y R4-02: CERRADOS.** Pasan las pruebas del área y mi sonda S-A, que probó el camino `POST` con 409 → V-17.
- **R5-01 (media, nuevo, dueño front-ajustes-a).** Es el defecto de R4-01 en otros puntos de ajustes-a: `DialogoCritico` (`comunes.tsx:228`) y seis handlers de alta o guardado en S02, S03 y S05 no tienen `catch`. Un fallo del servidor deja una promesa rechazada (`pageerror`) y no aparece ningún mensaje. Lo medí en 3 casos y en los 2 proyectos. Por la misma causa, V-18 («Eliminar promoción» con 409 `in_use`) no se muestra nunca.
- **H-13 (baja, reabierto, dueño front-creditos-salud-reportes).** Con 503 `unavailable`, P-20 muestra «tiempo máximo agotado» en lugar de «Comparador no disponible» (U1 N-226).
- **R3-01 (baja, dueño frontend-base).** Sigue abierto: `ajv` no está declarado en `package.json`.

## Archivos
- `docs/estado/wf3/REVISION-frontend-areas-r5.md` (informe)
- `docs/estado/wf3/revisor-frontend-areas-r5.md` (este reporte)
- Temporales en `datos/tmp/rev-fa-r5/` (build, configuración de Playwright, sondas y salidas). Ya los borré.

## Verificación (mis comandos, salidas resumidas)
- **C3-03 build**: `npx tsc --noEmit -p .` rc 0. `npx tsc --noEmit -p .storybook/tsconfig.json` rc 0. `npx vite build --outDir datos/tmp/rev-fa-r5/dist` da «✓ built» y «precache 81 entries».
- **C3-03 rutas**: `npm run verificar:rutas` da «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0». `python3 -I scripts/c3-integracion.py` (solo lectura) da «210 … fuera: 0 … sin registrar: 0» e «implementado sí 42/42».
- **U5-08 (referencia)**: `node scripts/hosts-externos.mjs <dist tmp>` da «hosts externos que se pedirían = 0».
- **Playwright de las 3 áreas** (`-c <tmp> ajustes-a ajustes-b creditos-salud-reportes`): 119 passed, 1 skipped, 0 failed.
- **U5-09**: `-g "U5-09|F-1 desde Inicio"`: 4 passed.
- **U5-06 hex**: 0 fuera de `src/estilos/tokens.css`. `rgb(`, `hsl(`, `oklch(` y `dark:`: 0.
- **U5-06 variantes**: `npm run build-storybook` y `contraste-variantes.mjs --salida <tmp>`: 37 historias, 148 páginas, «FALLAS nuevas: 0», «0 de 5 se reproducen», «RESULTADO: OK».
- **U5-06 y U5-07, barrido propio**: 128 páginas (32 rutas × 2 temas × 2 proyectos). Resultado: 0 `pageerror`, 0 desbordes, 0 controles < 44 px en táctil, 0 dinero sin `tabular-nums` y 0 textos en inglés. En 64 páginas a 1440: 0 fechas ISO visibles y 0 cifras sin formato.
- **U4, campos**: script propio sobre las tablas de las 14 secciones: 242 filas y 0 faltantes (207 exactas y 35 aproximadas).
- **U4 §0.2**: hay vista previa en las 13 secciones con efecto numérico, y las pruebas de `Ctrl+/` y del tablero pasan.
- **U4, claves `reports.*` y `alerts.*`**: las 11 claves de C1 son editables con `updateReportSettings` y `updateAlertSettings`, también con `reset_parameters`.
- **U4 §0.3**: NO CUMPLE (R5-01). Sondas S-C, S-D y S-F: `pageerror` 1 y 0 «No se pudo guardar» en los 2 proyectos.
- **U1, P-17 a P-25**: 21 de 22 textos y elementos de §3.19 a §3.27 están presentes. Falta «Comparador no disponible»: mi sonda S-G lo encontró 0 veces.

## Supuestos
- Clasifiqué R5-01 como **media**, igual que R4-01 y R4-02: es el mismo criterio de U4 §0.3 y se incumple de la misma forma.
- Mantengo H-13 en **baja**, como lo clasificó la ronda 1: el estado de error existe, no deja resultados parciales y tiene Reintentar. Solo falta el texto de U1.
- La franja «No se pudo enviar» de /reportes es contenido esperado (N-245, el estado de un envío fallido en los datos de prueba), no un defecto.

## Bloqueos
- Ninguno.

## Pendientes
- R5-01 para front-ajustes-a, con una revisión puntual después, si el orquestador lo decide (ver la PREGUNTA en el informe).
- H-13 para front-creditos-salud-reportes (baja).
- R3-01 para frontend-base (baja).

## Propuestas para Axel
- PROPUESTA: que la CONVENCIONES-frontend §4 diga explícitamente que todo `onClick` asíncrono que escribe pasa por un ayudante que captura el error, como `g.ejecutar` de ajustes-b o `Editable` de ajustes-a. También podría agregarse una prueba común que recorra los botones de Ajustes con el servidor en 500 y exija 0 `pageerror`. Así se cerraría este tipo de defecto de una vez, en lugar de caso por caso.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones: ~57.
- KB leídos: ~130.
- KB escritos: ~25, contando los dos reportes y los temporales ya borrados.
