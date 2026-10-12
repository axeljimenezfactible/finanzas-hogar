# revisor-frontend-areas-r4 · WF-3a (12-oct-2026)

Revisor independiente, ronda 4. No corregí nada. El informe completo está en `docs/estado/wf3/REVISION-frontend-areas-r4.md`.

## Hecho

- **Veredicto: NO APROBADO · 0 altos, 2 medios (R4-01 y R4-02, nuevos) y 1 bajo (R3-01).**
- Comprobé que los cuatro hallazgos de la ronda 3 quedaron **cerrados**: H-03 (alto), H-07, H-08 y H-10. También cerré R2-05, R2-06, R2-07 y R3-02.
- Encontré dos hallazgos nuevos:
  - **R4-01** (front-ajustes-a): en `S02Presupuesto.tsx:70-76` y `:93`, un 409 `budget_overlap` al confirmar un monto queda como promesa rechazada sin capturar. No se ve V-17 ni «No se pudo guardar». Lo medí con una sonda.
  - **R4-02** (front-ajustes-b): en `P2613Salud.tsx:112` y `:124` está el mismo patrón `void pedir(...).then(recargar)` sin `catch`. Lo vi leyendo el código.

## Archivos

- `docs/estado/wf3/REVISION-frontend-areas-r4.md` (informe)
- `docs/estado/wf3/revisor-frontend-areas-r4.md` (este reporte)
- Temporales en `datos/tmp/rev-fa-r4/`: build, configuraciones de Playwright, sondas, scripts y bitácoras. Ya están borrados.

## Verificación (cada criterio)

- **C3-03**:
  - `npx tsc --noEmit` rc 0; `vite build --outDir datos/tmp/rev-fa-r4/dist` «✓ built».
  - `npm run verificar:rutas`: «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0».
  - Las operaciones que cita U4 en las 14 secciones tienen todas su ruta literal en el código.
- **U5-08** (referencia): `hosts-externos.mjs` sobre mi dist da 0.
- **Playwright** de ajustes-a, ajustes-b y creditos-salud-reportes, en preview propio en el puerto 4797: **115 passed, 1 skipped, 0 failed**.
- **U5-06**:
  - Evento de /ajustes/05 a 412 táctil: 344×44 en claro y en oscuro, `controlesChicos` = [].
  - «Ver como tabla»: 36 px con puntero fino y 44 en táctil, en /creditos, /creditos/1 y /salud.
  - Descripción emergente: coincide con `--ui-surface-raised`, `--ui-border-divider` y `--ui-text-default` en los dos temas.
  - `contraste-variantes.mjs`: «0 de 5 se reproducen», «FALLAS nuevas: 0», «RESULTADO: OK».
  - Barrido de 120 páginas (30 rutas × 2 temas × 2 proyectos): 0 defectos.
  - `grep` de hex, `rgb(`, `hsl(`, `oklch(` y `dark:` en el código revisado: 0.
- **U5-07**: en el barrido, 0 fechas ISO, 0 fechas con barras, 0 importes mal formados y 0 textos en inglés; con el bloque de `Ctrl+/` abierto, `textosEnIngles` = [].
- **U5-09** (referencia): `erroresTrafico` = [] en las pruebas de las áreas.
- **U4**:
  - `Ctrl+/` con vista previa sumada en 02, 04, 05, 07, 09, 13 y 14: las 7 pruebas pasan, y su código pulsa `Control+/`, `Escape` y `Enter`.
  - Tableros de 04 y 14: pasan.
  - `reports.*` y `alerts.*`: hay pruebas que las editan y pasan.
  - 14 secciones: todos los campos de U4 están, algunos comprobados a mano.
  - Errores de guardado (§0.3 y V-17): **NO CUMPLE** (R4-01 y R4-02).
- **U1** P-17 a P-25: sus pruebas pasan, N-193 y N-208 incluidas. `c3-integracion.py` (solo lectura) da «implementado sí 41/42»; el que falta es P-03, que no es de estas áreas.
- **front-variantes**: `tsc -p .storybook/tsconfig.json` rc 0; Storybook se construye.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: rc 0, «Verificación ampliada OK».

## Supuestos

- HECHO: R4-02 viene de leer el código, no de una medición. La consecuencia del mismo patrón (rechazo sin capturar y ningún mensaje) sí la medí en R4-01.
- HECHO: el evento de /ajustes/05 a 1440 con puntero fino mide 108×24. Lo evalúo con el mínimo de 24 px para puntero fino; los 44 px son para táctil.
- PROPUESTA: Storybook se reconstruyó en `frontend/storybook-static`, que es una salida ignorada por git, porque el script lee esa ruta fija. Convendría que `contraste-variantes.mjs` aceptara `--estatico <dir>` para que un revisor no tenga que tocarla.

## Bloqueos

Ninguno nuevo.

## Pendientes

- R4-01 → front-ajustes-a; R4-02 → front-ajustes-b; R3-01 → frontend-base (dueño de `frontend/package.json`).
- PREGUNTA al orquestador o a Axel: ¿se encargan R4-01 y R4-02 con una revisión puntual después, o se aceptan como pendientes?

## Propuestas para Axel

- Ninguna de producto. Las correcciones de R4-01 y R4-02 no tocan cifras ni decisiones: solo muestran el mensaje que U4 ya define (V-17 y §0.3).

## Paquetes instalados

Ninguno.

## Consumo aproximado

~58 acciones · ~140 KB leídos · ~30 KB escritos (los dos documentos y los temporales ya borrados).
