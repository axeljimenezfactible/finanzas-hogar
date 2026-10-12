# revisor-frontend-areas-r3 · WF-3a · 12-oct-2026

Agente: `revisor-frontend-areas-r3` (revisor independiente, ronda 3 de 3). No corregí nada. Informe completo: `docs/estado/wf3/REVISION-frontend-areas-r3.md`.

## Hecho
- Volví a correr cada verificación con mis propios comandos, en `datos/tmp/rev-fa-r3/`: build y tsc, rutas, hosts, Playwright de ajustes-a, ajustes-b y creditos-salud-reportes, copias instrumentadas de esas pruebas en claro y en oscuro (tráfico `/api/` y medidas del estado final), sondeo propio de las gráficas, validador ajv 2020 propio contra `openapi.yaml`, Storybook y contraste de variantes, cruce de U4 por sección y `make verify-all`.
- Comprobé los 9 hallazgos altos y medios de la ronda 2: 5 cerrados (R2-01, R2-02, H-09, R2-03 en pantallas, R2-04), 2 reabiertos en parte (H-07 y H-08: los tableros ya están, falta `Ctrl+/`) y 2 siguen abiertos (H-03 y H-10, de frontend-base).
- **Veredicto: NO APROBADO, con 1 alto y 3 medios.** Los cuatro están en archivos de frontend-base o dependen de ella. Hay además 5 bajos: R2-05, R2-06 y R2-07 siguen abiertos, y R3-01 y R3-02 son nuevos.

## Archivos
- `docs/estado/wf3/REVISION-frontend-areas-r3.md` (informe)
- `docs/estado/wf3/revisor-frontend-areas-r3.md` (este reporte)

## Verificación
- C3-03: `npx tsc --noEmit` da 0 errores. `npx vite build --outDir datos/tmp/rev-fa-r3/dist` da «✓ built» con 80 entradas en precache. `npm run verificar:rutas` da «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- Pruebas de las 3 áreas: `npx playwright test -c <tmp> ajustes-a ajustes-b creditos-salud-reportes` da 85 passed y 1 skipped. Las copias instrumentadas dan 85 passed en claro y 85 passed en oscuro.
- Datos de prueba contra C1: con mi validador, 0 fallas en 417 cuerpos literales de 53 operaciones, y las 12 escrituras de Ajustes validan con su código (201 y 204 donde C1 los pide). En el tráfico grabado, las escrituras de créditos y reportes también validan. Los únicos «fallos» del tráfico grabado son cuerpos que el navegador ya no dejó leer, porque la página había cambiado. Cambian de una corrida a otra, y los cubrí con los cuerpos literales.
- U5-06, 44 px: queda 1 control bajo 44, el evento del calendario de /ajustes/05 de 40×44 (H-03, base). Las flechas de semana de 06 ya cumplen.
- U5-06, «Ver como tabla»: mide 20 px con puntero fino (H-03, base).
- U5-06, descripción emergente: `rgb(255,255,255)`, `rgb(183,185,190)` y `rgb(109,110,115)` también en oscuro (H-03, base).
- U5-06, variantes: `contraste-variantes.mjs` da 37 historias, 148 páginas, 5 230 mediciones y «FALLAS nuevas: 0». Se reproducen los 5 conocidos: FV-01 128, FV-02 12, FV-03 384, FV-04 52 y FV-05 16 («RESULTADO: OK»).
- U5-06, cifras tabulares: 0 casos en claro y en oscuro.
- U5-06, hex en el código: 0 en las 3 carpetas de pantallas, en `.storybook` y en las historias.
- U5-07: en el texto visible, 0 fechas ISO, 0 fechas con barras, 0 importes fuera de formato y 0 palabras en inglés.
- U5-09: «toques para registrar un gasto (sin conexión): 3».
- DESIGN §7, Botón: 0 controles visibles sin nombre.
- U4: en las 14 secciones están todos los campos. La única operación sin ruta es `queryView`, que U4 da como alternativa. Los tableros de 04 y 14 tienen vista previa, «Aplicar», Esc y Enter, con sus pruebas. Hay 0 manejadores de `Ctrl+/` (H-07, H-08).
- U4, `reports.*` y `alerts.*`: P2608 maneja 5 claves con `PATCH /api/v1/settings/reports` y P2609 maneja 6 con `PATCH /api/v1/settings/alerts`. Ninguna prueba las edita (R2-06).
- U1, P-17 a P-25: N-225 está en P20Comparador.tsx:37-45, 98 y 104. Pasan las pruebas de P-17 a P-25 y de F-4 a F-6.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` termina con salida 0: «Verificación ampliada OK».

## Supuestos
- Medí el tamaño de «Ver como tabla» con puntero fino contra el mínimo de 24 px que cita FV-02 (WCAG 2.2, DESIGN §4.6 y §9.3), como la ronda 2.
- Tomé `queryView` como opcional, porque U4 lo escribe «`listX` (o `queryView` …)».
- Revisé a mano las etiquetas de U4 que el script no encontró. Todas están en el código, aunque compuestas o parametrizadas.

## Bloqueos
Ninguno nuevo.

## Pendientes
- H-03, H-10 y la parte de la base de `Ctrl+/` (H-07, H-08) esperan una corrección de frontend-base: solicitudes 1 a 7 de `correccion-frontend-areas-r2.md`, más R3-01 y R3-02. Después, las áreas tienen que conectar `Ctrl+/` en 02, 04, 05, 07, 09, 13 y 14.
- Bajos: R2-05, R2-06, R2-07, R3-01 y R3-02.

## Propuestas para Axel
- PROPUESTA: como esta es la última ronda, encargar una corrección corta a frontend-base con H-03, H-10, la selección de filas con `Ctrl+/` en VistaBD, R3-01 y R3-02. Después, una revisión puntual solo de esos puntos y de la conexión de `Ctrl+/` en las áreas. La alternativa es aceptar el trabajo con esos pendientes registrados. PREGUNTA: ¿cuál de las dos?

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 40 acciones, ~120 KB leídos y ~45 KB escritos (incluidos los temporales, ya borrados).
