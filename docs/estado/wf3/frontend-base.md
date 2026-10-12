# WF-3a · frontend-base · reporte (12-oct-2026)

## HECHO
- `frontend/` creado desde cero: Vite 8.3.4 + React 19.3 + TypeScript 5.9.3 + Tailwind 4.3.3 (solo escalas de C2) + Radix (`radix-ui`, patrón shadcn/ui con cva) + TanStack Table 9.2.8 + dnd-kit + ECharts 6.1 + vite-plugin-pwa 2.0 + lucide-react 1.55 + msw 3 + Playwright 1.60 + axe-core.
- Tokens: `npm run tokens` copia **tal cual** `docs/contratos/tokens.css` y `tailwind.theme.css` a `src/estilos/` (sha256 comprobado) e Inter `woff2` + licencia a `public/fonts/` con `@font-face` (DESIGN §3.1).
- PWA: manifest es-MX (`start_url /inicio`, `standalone`, colores de C2 resueltos de tokens.css), íconos 192/512/enmascarable autoalojados (ícono `wallet` del sprite de Lucide del repo, `rsvg-convert` local), service worker `generateSW` con caché de la app (21 entradas) y respaldo de navegación.
- `src/app/`: enrutador P-00..P-26 (rutas PROVISIONALES), carga de sesión con `/auth/me` (copia del teléfono sin red), P-00 Iniciar sesión (422 y 429 con hora), regreso a la página tras sesión vencida, tema sistema/claro/oscuro con interruptor y persistencia (sin destello: `public/tema-inicial.js`), registro de pantallas por `import.meta.glob`.
- `src/shell/`: barra lateral con las 9 secciones en orden fijo, ▸, líneas de crédito dinámicas, 14 secciones de Ajustes, colapsable (‹‹, Ctrl+\, panel flotante en el borde); barra superior con migas y pestañas; barra inferior móvil (Inicio · Comidas · «+» · Compras · Buscar); menú «/» y búsqueda (Ctrl+K, «/», comandos de U1 §5.2); hoja P-03 (Gasto con 6 categorías sugeridas N-050, «Más», «Más opciones», regla K10; Ingreso y Ticket llevan a su pantalla; Peso = P-22 con último peso N-241); franja «Sin conexión con el servidor de casa · n en cola» y panel para revisar, reintentar o descartar con confirmación.
- `src/api/`: tipos de `openapi.yaml` con openapi-typescript; cliente `pedir()` tipado por ruta y método, `Idempotency-Key` exigida por tipo donde C1 la declara, errores C1 §3, paginación por cursor, copia de lecturas para sin conexión, formatos de dinero (centavos ↔ `$1,234.56`, U+2212), fechas (`10-oct-2026`, `Vie 30 oct 2026`) y `captured_at` con desplazamiento.
- `src/cola/`: cola en IndexedDB solo para gastos y peso (K9): clave 16–128, `captured_at` del momento de captura, envío solo al volver la red, en orden, sin duplicar, Web Locks entre pestañas, falla → queda en la cola y avisa.
- `src/vistas-bd/`: tabla (TanStack 9: orden, filtro, agrupar), tablero (dnd-kit con teclado), calendario (lunes a domingo, K8), galería, lista ordenable, edición en el lugar con «Antes → Después».
- `src/componentes/`, `src/textos/es-MX.ts`, `src/mocks/` (216 operaciones), pruebas base y `docs/estado/wf3/CONVENCIONES-frontend.md` con pantalla de ejemplo completa.

## ARCHIVOS
Autor: `frontend/{package.json,package-lock.json,.gitignore,index.html,tsconfig.json,vite.config.ts,playwright.config.ts}`, `frontend/public/tema-inicial.js`, `frontend/scripts/{comun,tokens,iconos,generar-mocks,rutas-api,hosts-externos}.mjs`, `frontend/src/main.tsx`, `frontend/src/estilos/app.css`, `frontend/src/app/{rutas.ts,enrutador.tsx,tema.ts,sesion.tsx,IniciarSesion.tsx,registro.ts,App.tsx}`, `frontend/src/shell/{estado.ts,Layout.tsx,BarraLateral.tsx,Barras.tsx,MenuComandos.tsx,HojaRegistro.tsx,FranjaSinConexion.tsx}`, `frontend/src/api/{tipos,errores,formatos,conexion,cliente,lecturas,consulta}.ts`, `frontend/src/cola/{idb,cola}.ts`, `frontend/src/componentes/{utils.ts,Boton,Campo,Chip,TarjetaMetrica,Avisos,Estados,Hoja,Grafica}.tsx`, `frontend/src/vistas-bd/{tipos.ts,CeldaEditable,VistaTabla,VistaTablero,VistaCalendario,VistaGaleria,ListaOrdenable,VistaBD}.tsx`, `frontend/src/textos/es-MX.ts`, `frontend/src/mocks/{manejadores,navegador}.ts`, `frontend/pruebas/{apoyo/api-falsa,login.spec,shell.spec,cola.spec}.ts`, `docs/estado/wf3/{CONVENCIONES-frontend,frontend-base}.md`, `docs/estado/bloqueos/WF3-frontend-{rutas,servidor,k31}.md`.
Generados (no se editan; en `.gitignore`, los rehace `npm run preparar`): `src/api/esquema.gen.ts`, `src/mocks/{operaciones,sinteticos}.gen.json`, `src/estilos/{tokens.css,tailwind.theme.css,colores-manifest.json}`, `public/fonts/*`, `public/icons/*`, `public/favicon.svg`.

## VERIFICACIÓN
- `npm install`: 519 paquetes, sin errores (paquetes abajo).
- `npx tsc --noEmit`: 0 errores. `npm run build`: sin errores; `dist/` 1.2 MB; JS 346.1 KB (108.9 KB gzip), CSS 44.2 KB (8.8 KB gzip), precache 21 entradas (1,143.6 KiB, incluye las dos Inter de 344 y 379 KB).
- `python3 -m json.tool dist/manifest.webmanifest`: válido; name «Finanzas Hogar», lang es-MX, íconos 192, 512 y 512 enmascarable.
- `npx playwright test` (con `FH_CHROMIUM`): **10 de 10 pasan** (móvil 412×915 táctil y escritorio 1440×900).
- U5-08: Playwright registra solo `localhost` (externos 0). `npm run verificar:hosts`: hosts que se pedirían = 0; literales que no son peticiones: w3.org (xmlns), react.dev y bit.ly (mensajes de error), tailwindcss.com (comentario de licencia), github.com y scripts.sil.org (texto de la licencia OFL que acompaña a Inter).
- U5-09: «toques para registrar un gasto (sin conexión): 3» (K10: «+», categoría, Guardar; teclear no cuenta).
- U5-10: fondo distinto con `colorScheme` oscuro y claro; con el interruptor `data-theme="dark"` y el fondo oscuro del token.
- U5-13: sin conexión → franja «Sin conexión con el servidor de casa», «1 en cola», «2 en cola»; recarga sin red conserva «2 en cola»; al volver la red: gastos 1 y pesos 1 (antes 0 y 0), contador en 0; clave UUID (36 caracteres) y `captured_at` `2026-10-12T02:57:55-06:00`; reintento forzado misma clave y cuerpo → 200 sin duplicar; con falla 422 la captura sigue («No se pudo enviar», «1 en cola») y se avisa «No se pudo enviar · Revisar». Captura de pantalla en `datos/tmp` (temporal, borrada).
- U5-07 (base): todos los textos de la base en `src/textos/es-MX.ts`; axe sin violaciones WCAG 2.2 AA en P-00.
- C3-03: `npm run verificar:rutas` → 12 rutas invocadas; fuera de `openapi.yaml`: 0; sin registrar en `routes.go`: 0 (216 registradas).

## SUPUESTOS
- Direcciones de las pantallas (bloqueo `WF3-frontend-rutas`), entrega de `dist/` por el servidor (`WF3-frontend-servidor`) y valores operativos de sondeo, reintento y copias (`WF3-frontend-k31`): PROVISIONALES.
- `typescript` 5.9.3 y no 7.0: `openapi-typescript` 7.13 exige `typescript ^5.x`. `@playwright/test` 1.60.0 por ser la versión cuyo Chromium ya estaba en la máquina.
- Sidebar fija desde 768 px (U1 §1.3 solo fija < 768 y ≥ 1024); P-00 vive en `src/app/` porque la carga de sesión la necesita.
- Las pestañas Ingreso y Ticket de P-03 llevan a P-08 (`?hoja=ingreso`) y P-05: sus formularios son del área.

## BLOQUEOS
- `docs/estado/bloqueos/WF3-frontend-rutas.md`, `WF3-frontend-servidor.md`, `WF3-frontend-k31.md` (ninguno detiene a las áreas).

## PENDIENTES
- P-03 «Más opciones»: faltan reparto distinto (Hogar, Karen o Conjunto con %), tienda y adjunto (necesitan `listCostShareRules`, `listStores`, `uploadExpenseAttachment`); hoy: fecha, persona, forma de pago, descripción y nota.
- Atajos de U1 §5.3 más allá de Ctrl+K/P, «/», Alt+G, Ctrl+\ y Ctrl+Enter; hoja de atajos (`/atajos` deja `?hoja=atajos`).
- Pantallas P-01..P-26 (salvo P-00): las construyen las áreas; el cascarón muestra «Pantalla en construcción».
- Las capturas de U5 (4 por pantalla) y la matriz C3 las hace su agente cuando haya pantallas.

## PROPUESTAS PARA AXEL
- Aprobar las direcciones en español (opción A del bloqueo de rutas).
- Aprobar sondeo 15 s, reintento 30 s y copia de lecturas hasta cerrar sesión (puerta 4).
- Descargar el Chromium de Playwright 1.60 una vez (`npx playwright install chromium`) para no depender de `FH_CHROMIUM`.

## PAQUETES INSTALADOS
react 19.3.0, react-dom 19.3.0, radix-ui 1.8.0, class-variance-authority 0.7.1, clsx 2.1.1, tailwind-merge 3.7.0, @tanstack/react-table 9.2.8, @dnd-kit/core 6.3.1, @dnd-kit/sortable 10.0.0, @dnd-kit/utilities 3.2.2, echarts 6.1.0, lucide-react 1.55.0; dev: vite 8.3.4, @vitejs/plugin-react 6.1.2, typescript 5.9.3, tailwindcss 4.3.3, @tailwindcss/vite 4.3.3, vite-plugin-pwa 2.0.0, workbox-window 7.4.1, openapi-typescript 7.13.0, msw 3.0.3, @playwright/test 1.60.0, @axe-core/playwright 4.13.0, @types/react 19.3.0, @types/react-dom 19.3.0, @types/node 24.10.1, yaml 2.9.1. Sin `sudo` (mkcert no hizo falta en esta fase).

## CONSUMO APROXIMADO
~51 acciones de herramienta; ~175 KB leídos (documentos por secciones, extractos de C1 con python, tipos de TanStack 9); **~189 KB escritos a mano** (código 156 KB, configuración 8 KB, documentos 25 KB), **por encima del tope de 150 KB de protocolo §4 / K32**: el alcance (frontend completo de la base) no cupo en el tope; lo generado por scripts (esquema.gen.ts 743 KB, mocks, tokens, fuentes, íconos) no se cuenta.
