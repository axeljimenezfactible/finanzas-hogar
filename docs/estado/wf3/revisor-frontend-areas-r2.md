# revisor-frontend-areas-r2 · WF-3a (12-oct-2026)

Soy revisor independiente (ronda 2 de 3). No corregí nada ni hice git add o commit. El informe completo está en `docs/estado/wf3/REVISION-frontend-areas-r2.md`.

## Hecho
- Volví a correr por mi cuenta la build, tsc, el script de rutas, el de hosts, las pruebas Playwright de las 3 áreas y el script de contraste de variantes, todo aislado en `datos/tmp/rev-fa-r2/`.
- Instrumenté copias temporales de las 4 specs de las áreas para tres cosas:
  - grabar el tráfico `/api/` y validarlo con ajv contra `openapi.yaml`;
  - medir al final de cada prueba los objetivos de 44 px y el dinero sin `tabular-nums`;
  - revisar el texto visible (ISO, fechas con barras, inglés y formato de dinero), en claro y en oscuro.
- Crucé las 14 secciones de U4 (filas y operaciones) y P-17 a P-25 de U1 (elementos N) con el código.
- Veredicto: **NO APROBADO**, con 3 altos (R2-01, R2-02 y H-03 en parte) y 6 medios (H-07, H-08, H-09 en su resto, H-10, R2-03 y R2-04).

## Archivos
- `docs/estado/wf3/REVISION-frontend-areas-r2.md` (informe)
- `docs/estado/wf3/revisor-frontend-areas-r2.md` (este reporte)

## Verificación
- **C3-03**:
  - `npx tsc --noEmit` da 0 errores.
  - `npx vite build --outDir datos/tmp/rev-fa-r2/dist` da «✓ built» (80 entradas en precache).
  - `npm run verificar:rutas` da 202 invocadas, 0 fuera de openapi y 0 sin registrar (216 registradas).
  - Por mi cuenta: las operaciones de U4 tienen su ruta en el código.
- **Playwright** (`-c datos/tmp/rev-fa-r2/pw.config.ts`, puerto 4798): `ajustes-a ajustes-b creditos-salud-reportes` da 78 passed. Las copias instrumentadas dan 78 passed en claro y 78 en oscuro.
- **Esquema C1**: 775 respuestas grabadas y 795 validaciones correctas. Las respuestas falsas de escritura no validan en 7 operaciones (R2-01), y lo confirmé con ajv sobre los cuerpos literales del código.
- **U5-06**:
  - 832 controles medidos en 412 táctil, con 3 bajo 44 px (40×44): las flechas de semana de /ajustes/06 y el evento del calendario de /ajustes/05.
  - Dinero sin `tabular-nums`: 1 caso, en /ajustes/02.
  - 0 hex en el código de las áreas.
  - `contraste-variantes`: 0 fallas nuevas y 5 de 5 hallazgos FV conocidos reproducidos (FV-01 128, FV-02 12, FV-03 384, FV-04 52, FV-05 16).
- **U5-07**: 0 fechas ISO, 0 fechas con barras y 0 importes mal formados en 156 estados finales. «Mar» es español (falso positivo).
- **U5-09**, de referencia: 3 toques en móvil, en escritorio y desde Inicio.
- **U4**: tabla de las 14 secciones en el informe. Las claves `reports.*` (5) y `alerts.*` (6) coinciden con C1 y se guardan con `PATCH settings/reports` y `PATCH settings/alerts`.
- **U1**: P-17 a P-25 completas salvo N-225 en P-20 (R2-04).
- **No corrí** `make verify-all` ni `verify-wf3-frontend`: otro agente estaba usando `frontend/dist` y el puerto 4173.

## Supuestos
- Las respuestas que grabé como «__no_json__» (el cuerpo ya no estaba disponible después de navegar) no las conté como fallas por sí mismas. Las confirmé leyendo el cuerpo literal en el código.
- Califiqué R2-01 como alto con la misma regla que H-04 de la ronda 1: una prueba con datos que no validan contra el esquema.

## Bloqueos
Ninguno nuevo.

## Pendientes
Ninguno del revisor.

## Propuestas para Axel
- Pedir a frontend-base que cierre de una vez VistaBD (`Ctrl+/` con vista previa sumada) y VistaTablero (vista previa, «Aplicar» y Esc). De eso dependen H-07 y H-08 en 7 secciones.
- Agregar al `verify` del frontend una validación ajv del tráfico falso, como la que usé aquí, para que R2-01 no regrese.

## Paquetes instalados
Ninguno. Usé ajv y yaml, que ya estaban en `frontend/node_modules`.

## Consumo aproximado
Unas 57 acciones, unos 150 KB leídos y unos 35 KB escritos (scripts temporales, informe y reporte). Borré los temporales de `datos/tmp/rev-fa-r2/`.
