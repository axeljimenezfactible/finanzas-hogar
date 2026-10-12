# Corrección frontend de áreas · ronda 2 · WF-3a · 12-oct-2026

Agente: `correccion-frontend-areas-r2` (Opus). Encargo: los 9 hallazgos altos y medios de `REVISION-frontend-areas-r2.md`.
Alcance: solo archivos de front-creditos-salud-reportes, front-ajustes-a y front-ajustes-b que nombran los hallazgos. No toqué archivos de frontend-base (`src/componentes/`, `src/vistas-bd/`, `pruebas/apoyo/`, `package.json`): lo que les toca va abajo, en «Solicitudes a frontend-base», con la ruta y el cambio exacto. No reabrí decisiones, no cambié cifras y no hice `git add` ni `git commit`.

## Hecho

| # | Hallazgo | Estado | Qué cambió |
|---|---|---|---|
| 1 | R2-01 (alta) respuestas falsas de escritura fuera de C1 | **Corregido** | Cada escritura que hacen las pruebas tiene su respuesta con `satisfies Respuesta<ruta, método>` y el código de C1: 201 en `POST promotions` y `POST recipes/{id}/duplicate`, 204 en todo DELETE, y 200 con el esquema completo en `PATCH settings/security`, `budgets/{id}`, `credits/{id}`, `categories/{id}`, `settings/data`, `settings/reports`, `settings/alerts`, `ai/providers/{provider}` e `inventory-items/{id}`. Una escritura sin respuesta falsa ya no recibe `200 {}`: se marca como error. Validador nuevo `pruebas/ajustes-a/validar-c1.ts`: ajv 2020 sobre `docs/contratos/openapi.yaml`. Valida el cuerpo de cada petición y cada respuesta falsa con su código, también los GET. Un `afterEach` exige cero errores en cada prueba de ajustes-a y de ajustes-b. |
| 2 | R2-02 (alta) flechas de semana de 40×44 | **Corregido en el área** | `S06Tiendas.tsx:199,201`: `className="shrink-0"` en «Semana anterior» y «Semana siguiente». Una prueba nueva lo mide a 412 táctil: ancho y alto ≥ 44 y `controlesChicos` sin «Semana». Falta la parte de la base (`shrink-0` en `tam: 'icono'`), que va en solicitudes. |
| 3 | H-03 (alta) evento de VistaCalendario, «Ver como tabla» y descripción emergente | **Fuera de alcance (frontend-base)** | Todo está en `VistaCalendario.tsx` y `Grafica.tsx`. Ver solicitudes 1 a 3. |
| 4 | H-07 (media) `Ctrl+/` en 02, 04, 05 y 07 y arrastre del tablero de 04 | **En parte**: tablero hecho; `Ctrl+/` depende de la base | `S04Categorias.tsx`: al soltar una tarjeta en otro reparto (`columnasTablero.moverA`), la tarjeta queda en la columna nueva con la vista previa de U4 l. 236 (ejemplo de $199.00, gastos fijos sin reparto propio y aviso), con «Aplicar» y «Cancelar». El foco va a «Aplicar», así que Enter aplica con `updateCategory`. Esc la regresa sin guardar. Al soltar no se guarda nada. `Ctrl+/` necesita selección de filas en VistaBD y VistaTabla (base, solicitud 6). |
| 5 | H-08 (media) `Ctrl+/` en 09, 13 y 14 y arrastre del tablero de 14 | **En parte**: tablero hecho; `Ctrl+/` depende de la base | `P2614Inventario.tsx`: el tablero deja de ser de solo lectura. Al soltar en el otro modo se ve el estado, la existencia estimada, los días restantes y el consumo por semana, antes → después (la `previa` que ya existía), con «Aplicar» o Enter (`updateInventoryItem`) y Esc. `Ctrl+/`: igual que el punto 4. |
| 6 | H-09 resto (media) aviso del alcance sin `tabular-nums` | **Corregido** | `S02Presupuesto.tsx:99`: el texto del `Aviso` va en `<span className="tabular-nums">`. Una prueba nueva lo comprueba con `dineroSinTabulares`: 0 casos. |
| 7 | H-10 (media) FV-01, FV-04 y FV-05 | **Fuera de alcance (frontend-base)** | `Boton.tsx`, `Grafica.tsx`, `Campo.tsx` y `Chip.tsx`. Ver solicitudes 4 y 5. |
| 8 | R2-03 (media) celdas editables vacías sin nombre | **Corregido en el área** | En 12, 13 y 14, las propiedades editables que podían quedar vacías pintan «Sin valor» (`T.comun.sinValor`): `notes` en las tres y `rate_period_days` en 14. Una prueba nueva recorre 12, 13 y 14 y da 0 botones sin nombre. La base debería pintar lo mismo por omisión (solicitud 7). |
| 9 | R2-04 (media) P-20 sin los valores de Ajustes (N-225) | **Corregido** | `P20Comparador.tsx` lee `getOptimizerSettings`. Precarga el colchón, el máximo de compras vivas, la fecha de deuda en cero y los 4 pesos globales, con la ayuda «De Ajustes». Si un valor se cambia, la ayuda dice «Solo para esta comparación; el de Ajustes no cambia.». Muestra los pesos normalizados en % y el tope por quincena de cada persona (`person.fortnight_payment_cap_cents`, de `/persons`; «Sin tope» si está vacío). La comparación manda los valores mostrados. Si falta una global, el campo queda vacío, como antes. |

## Archivos (todos en `/home/deepsec/finanzas-hogar/frontend/`)
- Pantallas: `src/pantallas/ajustes-a/S02Presupuesto.tsx`, `S04Categorias.tsx`, `S06Tiendas.tsx`, `textos.ts`; `src/pantallas/ajustes-b/P2612Personas.tsx`, `P2613Salud.tsx`, `P2614Inventario.tsx`, `textos.ts`; `src/pantallas/creditos-salud-reportes/P20Comparador.tsx`, `textos.ts`.
- Pruebas: `pruebas/ajustes-a/ajustes.escritorio.spec.ts`, `pruebas/ajustes-a/validar-c1.ts` (nuevo), `pruebas/ajustes-a/tablero.ts` (nuevo: arrastre con teclado); `pruebas/ajustes-b/datos.ts`, `ajustes-b.escritorio.spec.ts`; `pruebas/creditos-salud-reportes/datos.ts` (`GET settings/optimizer`), `pantallas.escritorio.spec.ts` (P-20).
- Reporte: `docs/estado/wf3/correccion-frontend-areas-r2.md`.

## Verificación
Corrí todo aislado: build en `datos/tmp/corr-fa-r2/dist`, preview en el puerto 4799 y una configuración temporal de Playwright que extiende la del repositorio. No usé `frontend/dist` ni el puerto 4173, por si otro agente los estaba usando. Ya borré la carpeta temporal.
- C3-03 build: `npx tsc --noEmit` da 0 errores. `npx vite build --outDir …` da «✓ built» y la PWA generada.
- C3-03 rutas: `npm run verificar:rutas` da «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- U5-08: `node scripts/hosts-externos.mjs <dist>` solo encuentra literales que no son peticiones (licencia de Inter y bit.ly en workbox).
- Pruebas de las 3 áreas: `npx playwright test ajustes-a ajustes-b creditos-salud-reportes` da **85 passed y 1 skipped**. El skipped es la medida táctil, que solo corre en `movil`. Antes eran 78. Las nuevas: tablero de 04, tablero de 14, medidas U5-06 y nombres R2-03.
- Suite completa (todas las áreas y la base): 148 passed y 1 skipped. La primera corrida tuvo una falla intermitente del tablero de 04 en escritorio: el arrastre con teclado iba más rápido que el pintado de dnd-kit. `tablero.ts` ahora avanza flecha por flecha hasta que la región en vivo de dnd-kit anuncia la columna. Con `--repeat-each 4`, las pruebas de tablero dan 16 de 16.
- Datos de prueba contra C1 (R2-01): el `afterEach` con `erroresC1` da cero errores en las 37 pruebas de ajustes-a y ajustes-b, en móvil y en escritorio. Antes de corregir, el validador reproducía los fallos de la revisión: `PATCH settings/security` sin `parameters`, `DELETE telegram-chat-candidates` con 200 cuando C1 pide 204, `POST promotions` con 200 cuando C1 pide 201, y `PATCH ai/providers` sin `provider` ni `allowed_uses`.
- U5-06, 44 px (R2-02): a 412 táctil, «Semana anterior» y «Semana siguiente» miden ≥ 44 por lado (prueba «U5-06 en 412 táctil…»).
- U5-06, cifras tabulares (H-09): `dineroSinTabulares` da `[]` en /ajustes/02 con el diálogo del alcance abierto.
- DESIGN §7, Botón (R2-03): 0 botones visibles sin nombre en /ajustes/12, /13 y /14.
- U5-06, colores en las pantallas de las áreas: `grep` de hex, `rgb(`, `hsl(`, `oklch(` y `dark:` en las 3 carpetas da 0.
- U1 N-225 (R2-04): la prueba de P-20 comprueba el colchón precargado (500.00), el máximo de compras vivas (3) con «De Ajustes», el cambio que vale solo para esta comparación, los pesos al «25 %», el «Sin tope» y el cuerpo enviado (`min_available_buffer_cents: 50000`, `max_live_purchases: 2` y los pesos).
- No corrí `make verify-all` ni `make verify-wf3-frontend`, porque escriben `frontend/dist` y usan el 4173, que comparto con otros agentes. En su lugar corrí lo mismo aislado.

## Supuestos
- Celda vacía: usé el texto que ya tenía el área, «Sin valor» (`T.comun.sinValor`, igual que `textos.vistas.sinValor` de la base), y no el «Sin dato» que sugiere la revisión, para no tener dos textos para lo mismo.
- P-20: si la comparación manda la global tal cual, el resultado es el mismo que si no la manda (C1: un campo ausente usa la global). La holgura (`date_window_days`) no se precarga porque N-225 no la nombra. Los pesos normalizados se redondean a % entero con ROUND_HALF_UP en enteros.
- Tablero: la vista previa con «Aplicar» y Esc se hizo en las áreas con el `columnasTablero.moverA` que ya existía: `moverA` no guarda, solo deja el cambio pendiente. No hizo falta cambiar la base.
- El validador usa ajv 8.20.0, que ya está en `node_modules` porque lo trae `workbox-build`. No instalé nada. `package.json` es de la base (solicitud 8).

## Solicitudes a frontend-base (fuera de mi alcance)
1. H-03: en `src/vistas-bd/VistaCalendario.tsx:30`, el botón del evento con `coarse:min-h-11 coarse:min-w-11 shrink-0`, o una lista de agenda en `phone:`. Hoy mide 40×44 en /ajustes/05.
2. H-03 (FV-02): en `src/componentes/Grafica.tsx:78`, el `<summary>` de «Ver como tabla» con `inline-flex min-h-6 items-center coarse:min-h-11`. Hoy mide 20 px con puntero fino.
3. H-03 (FV-03): en `src/componentes/Grafica.tsx:48`, `tooltip: { backgroundColor, borderColor, textStyle: { color } }` leídos con `getComputedStyle(document.documentElement)` de `--ui-surface-raised`, `--ui-border-default` y `--ui-text-default`, y releídos al cambiar el tema.
4. H-10 (FV-01): en `src/componentes/Boton.tsx:6`, cambiar `transition-colors` por `transition-[background-color]`, para que el anillo de foco no tenga transición.
5. H-10 (FV-04 y FV-05): en `Grafica.tsx`, `emphasis: { disabled: true }` o un `emphasis` sin recolorear (solo colores de C2). En `Campo.tsx` y `Chip.tsx`, el estilo `disabled` de DESIGN §7.0, como está en `frontend/.storybook/hallazgos.json`.
6. H-07 y H-08: en `src/vistas-bd/VistaBD.tsx` y `VistaTabla.tsx`, selección de filas y `Ctrl+/` con `editarVarias: { propiedades: string[]; vistaPrevia(filas, propiedad, valor): Promise<Cambio[]>; guardar(filas, propiedad, valor) }`, con el bloque sumado antes de aplicar, Enter aplica y Esc cancela. Cuando exista, ajustes-a lo conecta en 02 (categoría y activo de gastos fijos), 04 (tipo, reparto y activa), 05 (promociones) y 07 (tiempos aptos), y ajustes-b en 09 (alertas), 13 (metas, «Activa») y 14. Opcional: `columnasTablero.vistaPrevia` en VistaTablero, para pasar ahí lo que hoy hacen 04 y 14.
7. R2-03: en `src/vistas-bd/VistaTabla.tsx:11`, cuando el contenido de `Celda` sea `''`, pintar `textos.vistas.sinValor`. En `CeldaEditable.tsx:39`, `aria-label={`${prop.nombre}: ${texto}`}`.
8. R2-01: pasar `pruebas/ajustes-a/validar-c1.ts` a `pruebas/apoyo/` para que lo usen todas las áreas, y declarar `"ajv": "8.20.0"` en `devDependencies` de `frontend/package.json`.
9. R2-02: `shrink-0` en `estilosBoton` para `tam: 'icono'` (`src/componentes/Boton.tsx`).

## Bloqueos
Ninguno nuevo. No hizo falta decidir ningún valor.

## Pendientes
- `Ctrl+/` en 02, 04, 05, 07, 09, 13 y 14: espera la solicitud 6.
- H-03 y H-10: esperan las solicitudes 1 a 5.
- Bajos de la revisión que no estaban en el encargo: R2-05 (vista previa 10 de 05), R2-06 (pruebas que editen las claves `reports.*` y `alerts.*`; sus respuestas falsas válidas en C1 ya están en `pruebas/ajustes-b/datos.ts`) y R2-07 (citar el bloqueo en el punto de equilibrio de 06).

## Propuestas para Axel
- Que el revisor de la ronda 3 corra el validador `erroresC1` también en las pruebas de las otras áreas cuando la base lo mueva a `pruebas/apoyo/`.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 52 acciones, ~150 KB leídos y ~40 KB escritos.
