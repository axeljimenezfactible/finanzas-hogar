# revisor-frontend-core-r3 · reporte · WF-3a · 12-oct-2026

Revisor independiente, ronda 3 de 3, del trabajo de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras`, y de la corrección `correccion-frontend-core-r2`.
No escribí ni corregí código. No hubo `git add` ni `git commit`. No se envió Telegram.
Informe completo: `docs/estado/wf3/REVISION-frontend-core-r3.md`.

## Hecho

- **Veredicto: NO APROBADO.** Hay 0 altos, 2 medios (N-04 y N-05) y 1 bajo nuevo (N-06).
- **Ronda 2:**
  - N-01 (alto) y N-02 (medio) quedan **cerrados**.
  - N-03 (bajo) y los bajos de la ronda 1 (H-07, H-08, H-09 y H-11) siguen abiertos. No bloquean.
- **Medios nuevos (los dos de `frontend-base`):**
  - **N-04:** el `aria-label` de toda `Grafica` sale en inglés y con centavos crudos. ECharts escribe ahí su descripción automática y sustituye la etiqueta en español (`componentes/Grafica.tsx:44`).
  - **N-05:** la barra inferior no sigue DESIGN §7.10 (`shell/Barras.tsx:42-56`):
    - iconos de 20 px en lugar de 24;
    - fondo `ui.surface.card` en lugar de `ui.surface.raised`;
    - el destino activo no lleva la píldora de 56 × 32 px.
- **Bajo nuevo:**
  - **N-06:** no hay barrera de errores global. Una respuesta con forma inesperada deja la aplicación en blanco.

## Archivos

- `docs/estado/wf3/REVISION-frontend-core-r3.md` (informe)
- `docs/estado/wf3/revisor-frontend-core-r3.md` (este reporte)
- Temporales en `datos/tmp/rev-fcore-r3/`, ya borrados: pruebas propias de Playwright con su configuración, `rutas.py`, `grep-estilo.sh` y bitácoras.

## Verificación

- **C3-03:**
  - `npm run build` → exit 0, 0 `error TS`.
  - `npx tsc --noEmit -p .` → exit 0.
  - `python3 -I rutas.py` (propio) → 163 rutas literales y **0 fuera** de los 172 paths de los extractos.
  - `npm run verificar:rutas` → «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- **U5-08:**
  - En el barrido propio solo hay un host: `localhost:4173`.
  - `npm run verificar:hosts` → «hosts externos que se pedirían = 0».
  - En `dist` solo hay URLs literales que no se piden.
- **U5-06:**
  - Toques: **0 controles menores de 44 × 44** entre 729 controles de 31 vistas y los estados abiertos. «Más opciones» mide 123×44 y «Menos opciones» 141×44.
  - Iconos: **0 fuera de 20/24 px o de trazo 1.5** (381 en móvil y 596 en escritorio oscuro).
  - Solo icono: **0 con 20 px**.
  - Barra inferior: **4 iconos de 20 px** → N-05.
  - Dinero sin `tabular-nums`: **0**.
  - Valores fuera de C2: 0 hex, 0 `rgb()`/`hsl()`/`oklch()` y 0 `dark:`.
  - Inicio: **6 tarjetas exactas**.
- **U5-07:**
  - Inglés en el texto pintado, `placeholder`, `title` y `alt`: **0**, con una lista de 120 términos.
  - Fechas ISO visibles: **0**.
  - En el `aria-label` de las gráficas: «This is a chart… the data for sep-2026 is 0, 400000» → N-04.
- **U5-09:**
  - Todas las e2e del grupo dan **3**.
  - Mi variante sin red da **3** con la regla K10.
- **U5-10:**
  - Oscuro: `body` `#222325` (`color.neutral.950`) en 31 vistas.
  - Claro: `#FAFAFB`.
  - El interruptor recorre Claro → Oscuro → Según el sistema, y «Según el sistema» sigue al sistema.
- **U5-13 (pasos 1–3 y 6–7, con mocks):** variante propia con 1 gasto y 1 peso.
  - Recargar: «2 en cola».
  - 500: aviso y franja con «No se pudo enviar · 2 en cola».
  - Reintento: 1 gasto y 1 peso, con 2 POST que llevan **1 sola clave** (36 caracteres) y **0 duplicados**.
  - `captured_at` lleva `-06:00`, del momento de captura.
  - La franja desaparece al vaciarse la cola.
- **U1 §2.2 y U3:**
  - Vacíos de P-02, P-11 y P-15 con su texto.
  - Errores con Reintentar en P-02, P-04, P-06, P-08, P-10, P-13 y P-14.
  - P-01 muestra el error por tarjeta: 6 «Reintentar» con el texto de U3.
- **e2e del grupo:** `npx playwright test pruebas/inicio pruebas/gastos-presupuesto pruebas/comidas-compras pruebas/cola.spec.ts pruebas/shell.spec.ts pruebas/login.spec.ts pruebas/medidas.spec.ts` → **64 passed**.
- **`make verify-wf3-frontend`:** exit 0, **142 passed**.
- **`make verify-all GOOSE=…`:** **exit 2**, por backend y fuera del alcance de esta revisión.
  - `getDashboardBestBasket` no tiene handler y responde 501.
  - Por eso fallan `TestAreas_SinOperacionesEn501`, `TestContrato_Operaciones`, `TestContrato_Rutas` y `TestPermisos_MismasOperacionesAmbasCuentas`.
  - El frontend la llama desde la tarjeta «Mejor canasta» de Inicio.

## Supuestos

- **HECHO:** clasifiqué N-04 como **media** y no alta. La tabla de severidad pone «texto en inglés» como medio. Además, es texto que se lee en voz alta pero no se pinta.
- **HECHO:** clasifiqué N-05 como **media**. Es del mismo tipo que N-02 (regla de uso de §5.2) y suma tokens de rol de §7.10.
- **HECHO:** que la cola se detenga en la primera captura que falla, sin enviar las siguientes, cumple U1 §2.4 («en orden de captura»). No es un hallazgo.
- **HECHO:** N-06 es bajo. Solo se dispara con una respuesta fuera del contrato (la API falsa de H-07).

## Bloqueos

Ninguno nuevo.

## Pendientes (para la corrección; dueño `frontend-base`)

1. **N-04:** en `src/componentes/Grafica.tsx:44`, pasar `aria: { enabled: true, label: { description: \`${titulo}. ${resumen}\` }, decal: { show: true } }`. Ampliar la lista de inglés de las pruebas.
2. **N-05:** en `src/shell/Barras.tsx:42-56`:
   - iconos con `size-6` y `size={ICONO_MD}`;
   - fondo `bg-surface-raised`;
   - píldora del destino activo en `bg-surface-selected` con el icono en `text-icon-selected`;
   - regla «barra inferior → 24 px» en `pruebas/apoyo/medidas.ts`.
3. Opcionales:
   - N-06: barrera de errores por pantalla en `Layout`.
   - N-03: `Barras.tsx:17,26` `size-8` → `size-9 coarse:size-11`.
   - H-07, H-08, H-09 y H-11.

## Propuestas para Axel

- **PROPUESTA:** que el orquestador asigne al dueño de backend `getDashboardBestBasket` (501), porque rompe `make verify-all` y deja sin datos la tarjeta «Mejor canasta» de Inicio contra el servidor real.
- **PROPUESTA:** sumar a la lista de verificación de U5-07 de las pruebas el nombre accesible (`aria-label` y `aria-description`) de las gráficas. Hoy ninguna prueba lo mira.

## Paquetes instalados

Ninguno.

## Consumo aproximado

- Acciones: ~50.
- Leído: ~150 KB (reportes de la ronda 2, CONVENCIONES-frontend, secciones de DESIGN, U5, C3, U1 y U3, código y bitácoras).
- Escrito: ~45 KB (pruebas y scripts temporales ~22 KB; informe y reporte ~23 KB).
- `datos/tmp/playwright-frontend/` es la salida compartida de la configuración de Playwright del frontend, que también usan otros agentes, así que **no la borré**.
