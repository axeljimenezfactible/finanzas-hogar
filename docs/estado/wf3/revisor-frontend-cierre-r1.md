# revisor-frontend-cierre-r1 · reporte (12-oct-2026)

Agente revisor independiente: evalúa y no corrige. El informe completo está en `docs/estado/wf3/REVISION-frontend-cierre-r1.md`.

## Hecho
- HECHO: revisé el diff de `frontend/src` (Barometro, comun, Tablero G-6, S02Presupuesto y textos) y de `frontend/pruebas-u5` (contraste-dom.js, flujos, pendientes y c3-09), la prueba nueva `parametros-tablero.escritorio.spec.ts` y el reporte `frontend-cierre.md` con su Continuación 1.
- HECHO: volví a correr todas las verificaciones que pide el encargo, de la (a) a la (e). Veredicto: **no aprobado**, con 0 altos, 2 medios (M-1 y M-2) y 4 bajos.

## Archivos
- Escritos: `docs/estado/wf3/REVISION-frontend-cierre-r1.md` y este reporte. No toqué ningún archivo del producto.

## Verificación
- (a) `npx tsc --noEmit` sale con 0. `npx vite build --outDir datos/tmp/rev-fc-r1/dist` da «precache 80 entries» y `diff -rq` contra `frontend/dist` da 0 diferencias. `npx playwright test pruebas/ajustes-a/ pruebas/inicio/` da 109 passed, 1 skipped.
- U4 P-26.02 / §0.2: las 5 claves, «Antes → Después», Esc y Enter, y Restablecer, en verde con datos falsos. Con el servidor real, las unidades salen largas y repetidas (M-1).
- U5-04: recaptura de P-01 Resumen sobre una copia limpia. Bajo umbral: android-claro 0 (579 pares), android-oscuro 0 (579), escritorio-claro 0 (610), escritorio-oscuro 0 (610). Detalle da 0 y P-26.02 da 0. `make verify-u5` da «bajo umbral = 2» porque la evidencia versionada no se ha regenerado (B-6).
- U5-06: el eje del barómetro mide 10 × 10 con anillo de 2 px. `scrollWidth` es 412/412 en P-01 Detalle Android y 1440/1440 en escritorio.
- U5-07: `textosEnIngles` da [] en P-01 Resumen, P-01 Detalle y P-26.02 en las 4 combinaciones. Las líneas agregadas tienen 0 hex, rgb, hsl u oklch.
- C3-03: rutas-api da «213 rutas invocadas; fuera de openapi.yaml: 0; sin registrar: 0». `c3-integracion.py` da «no llama: 1 ['listPeriods']; sin cobertura: 0», así que `updateDashboardSettings` sí se llama. rc 1 solo por la matriz del documento (B-6).
- C3-08, cada archivo en una copia limpia (8444): `flujos.spec.ts` da 7 passed (F-1 con 3 toques, F-2 con 4, F-3A con la condición (c), F-3B con 5, F-4 con 3, F-5 con 3, F-6 con 3). `pendientes.spec.ts` da 4 passed (F-3A3, la ruta A completa, en 3 toques). F-3A0 sale ✓ con su paso en FALLA (B-3).
- C3-09 sobre una copia limpia: «120 cifras leídas en pantalla · diferencias != 0: 0». 13 son del tablero y P-17 (≥ 10).

## Supuestos
- Corrí C3-09 sobre una copia limpia de la base de desarrollo (8444) en vez del 8443. Los datos son los mismos y así no se escribe nada en el servidor de desarrollo.
- Califiqué M-2 como medio aunque es un defecto previo: es una cifra visible distinta del oráculo y está en un archivo que esta ronda tocó.

## Bloqueos
- Ninguno nuevo de mi parte. PROPUESTA: abrir el bloqueo de unidades de M-1 (lo decide el orquestador).

## Pendientes
- Ninguno de mi alcance.

## Propuestas para Axel
- Decidir cómo se muestran las unidades de los parámetros K44 y K54 en Ajustes (M-1). Opciones: unidad corta desde la API, textos propios en la interfaz, o mostrar la unidad solo si es corta.

## Paquetes instalados
- Ninguno.

## Temporales
- Borrados: `datos/tmp/rev-fc-r1`, `datos/tmp/u5-flujos`, `datos/tmp/playwright-frontend` y `datos/tmp/playwright-u5-flujos`. Dejé `datos/tmp/rev-bc` y `datos/tmp/verify-wf3`, que son de otro agente. `git status`: 59 entradas, las mismas que al empezar más mis 2 informes.

## Consumo aproximado
- ≈ 52 acciones, ≈ 150 KB leídos, ≈ 16 KB escritos.
