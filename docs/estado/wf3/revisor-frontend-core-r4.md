# revisor-frontend-core-r4 · WF-3a · 12-oct-2026

Revisor independiente, ronda 4. No escribí ni corregí nada de lo revisado. Informe completo: `docs/estado/wf3/REVISION-frontend-core-r4.md`.

## Hecho

- **Veredicto: NO APROBADO (0 altos, 2 medios y 3 bajos nuevos).**
- Comprobé cerrados N-04 y N-05 (medios de la ronda 3), además de N-06 y N-03 (bajos). H-07, H-08 y H-09 (bajos) siguen abiertos.
- Hallazgos nuevos:
  - **R4-01 (media, frontend-base):** P-03 no registra el ingreso en la hoja ni lista las líneas a crédito (U1 §3.5, N-102 y N-053). `HojaRegistro.tsx:39` y `:134`. Sus pruebas se saltan solas.
  - **R4-02 (media, frontend-base):** en `Grafica` el eje de dinero no es compacto, las cifras del SVG no son tabulares (también en `GraficaInicio`) y no hay etiquetas de dinero en texto (DESIGN §7.7 y §8.5, U5-06). `Grafica.tsx:66`.
  - R4-03 (baja, front-inicio): descripción automática de ECharts en inglés, oculta, en `GraficaInicio.tsx:49`.
  - R4-04 (baja, frontend-base): los destinos de la barra inferior miden 55 px y la zona segura resta en vez de sumar (`Barras.tsx:52`).
  - R4-05 (baja, frontend-base): excepción `PENDIENTES_INGLES` en `medidas.spec.ts:36`.

## Archivos

- Escritos: `docs/estado/wf3/REVISION-frontend-core-r4.md` y este reporte.
- Temporales en `datos/tmp/rev-fcore-r4/` (build propio, `rutas.py`, configuraciones y pruebas Playwright propias, registros). Ya están borrados.
- No toqué código, pruebas, contratos ni Pista 3. No hice `git add` ni `commit`.

## Verificación (comandos propios; detalle en el informe)

- **C3-03:**
  - `npx tsc --noEmit -p .` → exit 0.
  - `npx vite build --outDir ../datos/tmp/rev-fcore-r4/dist` → exit 0, precache de 81 entradas.
  - `rutas.py` propio: 165 literales `/api/v1`; fuera del contrato 0; métodos no declarados 0 en 364 llamadas.
  - `rutas-api.mjs` → «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0».
- **U5-06:**
  - Controles: 702 en 29 vistas a 412 táctil. Los únicos menores de 44 son 9 casillas de 20 px, cada una dentro de un `label` de 44 × 44.
  - Iconos: 0 fuera de norma (347 en móvil, 548 en escritorio).
  - Colores y medidas: 0 hex, 0 `rgb()`/`hsl()` y 0 `dark:`.
  - Inicio: 6 tarjetas exactas.
  - Barra inferior con los tokens de §7.10.
  - **Cifras tabulares: 26 montos en SVG sin `tabular-nums` → R4-02.**
- **U5-07:**
  - Lista propia de 150 términos: 0 en inglés en el texto y en los atributos expuestos.
  - Fechas ISO visibles: 0.
  - `Grafica`: `aria-label` en español.
- **U5-08:** hosts = solo `localhost:4894`. `hosts-externos.mjs` → 0.
- **U5-09:** 3 toques en línea desde P-02 (prueba propia) y 3 sin red. Las e2e del grupo también dan 3.
- **U5-10:**
  - Fondo oscuro `rgb(34,35,37)` = `--ui-surface-page`; claro `rgb(250,250,251)`.
  - El tema sigue al sistema.
  - El interruptor recorre Claro → Oscuro → Según el sistema.
- **U5-13 pasos 1–3, 6 y 7 (mocks), variante propia:**
  - Paso 1: «Sin conexión…» y «2 en cola».
  - Paso 2: tras recargar, «2 en cola» y 0 registros en el servidor.
  - Paso 7: con 500, aviso «No se pudo enviar» y 0 registros.
  - Al reintentar: 1 gasto y 1 peso.
  - Paso 3: 2 POST con 1 sola clave de 36 caracteres; `captured_at` con desplazamiento y anterior al envío.
  - Paso 6: el reenvío con la misma clave responde 200 con el mismo id; 0 duplicados.
- **U1 §2.2:**
  - Carga en P-01: 6 bloques `aria-busy` y nunca `$0.00`.
  - Error en P-01: 6 tarjetas con «Reintentar».
  - Error 500 en 9 pantallas: «Reintentar» y 0 caídas.
  - Vacíos en P-02, P-11 y P-15.
- **U1 §3 (P-03):** no cumple → R4-01.
- **e2e del grupo** (config propia sobre mi `dist`): **108 passed, 7 skipped, 0 failed** (2.8 min). De los 7 saltados, 4 son de P-03 por R4-01. Pruebas propias: 10 passed.
- **`make verify-all GOOSE=…`:** exit 0 («Verificación ampliada OK»). No cambió el árbol de trabajo.
- **`make verify-wf3-frontend`:** exit 0. Da build y tsc limpios, rutas fuera 0 y hosts 0; e2e de todas las áreas: 223 passed y 8 skipped (4.9 min). El árbol de trabajo no cambió.

## Supuestos

- HECHO: R4-02 mide los `<text>` del SVG de ECharts. La ronda 3 solo medía nodos HTML; por eso antes daba 0.
- HECHO: una casilla de 20 px dentro de un `label` de 44 × 44 cuenta como objetivo táctil de 44 (WCAG 2.5.8 y DESIGN §4.6: área táctil). Lo comprobé con `elementFromPoint`.
- HECHO: la descripción en inglés de `GraficaInicio` va como baja porque su contenedor es `aria-hidden` y no se enfoca. Con el mismo criterio de la ronda 3 (lo que lee TalkBack), no llega al usuario.
- HECHO: compilé y serví mi propio `dist` (puerto 4894) para no pisar a otro revisor que corría a la vez. El build equivale a `npm run build` sin el paso `preparar`, que regenera archivos ya versionados. `make verify-wf3-frontend` sí corre el `npm run build` completo.

## Bloqueos

Ninguno nuevo.

## Pendientes

- Corregir R4-01 y R4-02 (frontend-base). R4-02 toca también `GraficaInicio.tsx` (front-inicio), salvo que se use la regla común de `app.css`.
- Bajos: R4-03 (front-inicio), R4-04 y R4-05 (frontend-base), y H-07, H-08 y H-09.

## Propuestas para Axel

- PROPUESTA: que las pruebas no usen `test.skip` condicionado a lo que pinta la pantalla para algo que pide U1. Así una integración pendiente sale en rojo y no como «saltada». Pasó con P-03 (R4-01).

## Paquetes instalados

Ninguno.

## Consumo aproximado

~56 acciones; ~170 KB leídos; ~45 KB escritos (pruebas temporales e informes).
