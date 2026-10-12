# correccion-frontend-core-r1 · WF-3a · 12-oct-2026

Agente: `correccion-frontend-core-r1` (Opus). Corrige los hallazgos **altos y medios** de `docs/estado/wf3/REVISION-frontend-core-r1.md` (H-01..H-06 y H-10) en los archivos de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras`. Los bajos (H-07, H-08, H-09, H-11) quedan fuera del encargo y siguen pendientes. No se tocaron otras áreas, `docs/design/`, migraciones, motores ni `docs/estado/wfp3/`. Sin `git add` ni `git commit`.

## Hecho

| Hallazgo | Estado | Qué se cambió |
|---|---|---|
| **H-01** (alta) U5-13 paso 7 · U1 §2.4 | **Corregido** | `cola.ts`: ante 5xx o 401 la captura **sigue `pendiente`** (se reintenta sola), pero se guarda con el motivo (`error`) y cuenta en `ResultadoEnvio.retenidas`, que emite `fh:cola-enviada`. `App.tsx` avisa «No se pudo enviar · Revisar» también en `retenidas`. `useCola()` expone `retenidas`; `FranjaSinConexion.tsx` sigue visible con conexión mientras haya capturas rechazadas **o retenidas**, con «No se pudo enviar · n en cola» y «Revisar». En el panel, la captura retenida muestra «No se pudo enviar: <motivo>». Sin red sigue igual: `break` sin aviso, porque la franja ya dice «Sin conexión». |
| **H-02** (alta) U5-06 · DESIGN §5.2 | **Corregido** | Nuevo `src/componentes/Icono.tsx` con `ProveedorIconos` (`LucideProvider` de lucide-react 1.55.0: `size` 20, `strokeWidth` 1.5 y `absoluteStrokeWidth`), que envuelve la app en `App.tsx`. El resultado es `stroke-width` 1.8 a 20 px y 1.5 a 24 px, los valores de `icon.stroke.user-units.sm/md`. Los íconos de 16 y 12 px (`size-4`, `size-3`) pasan a `size-5` (20 px). Los de solo icono (cerrar aviso, plegar o desplegar barra, buscar en la barra superior, mes anterior o siguiente del calendario, asa de arrastre, quitar renglón, eliminar) y los de estados vacíos pasan a `size-6` con `size={ICONO_MD}` (24 px), como pide §5.2. La insignia de Inicio (`[&_svg]:size-4`) pasa a `size-5`. El ícono aplastado de 11.4 px de ancho (alerta de categoría en `gastos-presupuesto/comun.tsx`) lleva ahora `shrink-0`. |
| **H-03** (alta) U5-06 · DESIGN §4.6 | **Corregido** | `CeldaEditable.tsx`: la celda lleva `coarse:min-h-11 coarse:min-w-11`, y el «Reintentar» de la vista previa `coarse:min-h-11`. `VistaTabla.tsx`: los encabezados ordenables `coarse:min-w-11` (ya tenían `coarse:h-11`) y el botón de título de fila `coarse:min-h-11 coarse:min-w-11`. La prueba nueva `pruebas/medidas.spec.ts` recorre **todos** los controles visibles a 412 px táctil. |
| **H-04** (alta) U5-06 · DESIGN §4.6 | **Corregido** (en archivos de la base) | Los botones «Despensa» y «Transporte» de P-06 Detalle son el botón de título de fila de `VistaTabla.tsx` y quedaron con el arreglo de H-03. El enlace «Todos los gastos» de P-04 (y «Recetario» de P-12) es el enlace de regreso móvil de `shell/Barras.tsx`, que ahora lleva `min-h-11 gap-1`. Ningún archivo de `pantallas/gastos-presupuesto` tenía estos controles. Nota de corrección añadida en `front-gastos-presupuesto.md`. |
| **H-05** (alta) U5-06 · DESIGN §4.6 | **Corregido** | El enlace «Recetario» de P-12 se corrigió en `Barras.tsx` (igual que H-04). La casilla de P-13 ya vivía en una etiqueta clicable de 44 × 44 (`label.size-11`), y el revisor midió el `input` de 20 × 20. Ahora lo comprueba una prueba (`medidas.spec.ts`, «U5-06 P-13»): dibujo 20 × 20, área (etiqueta) 44 × 44, y tocar la esquina del área, fuera del dibujo, marca el artículo (PATCH con `week=2026-W42`). En `front-comidas-compras.md` quedó tachada la afirmación sin medir, con la corrección. |
| **H-06** (media) U5-06 · DESIGN §3.4 | **Corregido** | `TarjetaMetrica.tsx`: la nota lleva `tabular-nums`, lo que cubre P-02 «$200.00 más que el periodo anterior» y P-09 «Planificado: $3,000.00». `P14Canasta.tsx`: la frase «…la más barata: $3,094.50» lleva `tabular-nums`. La medición nueva encontró además 3 importes de la galería de P-09 Detalle sin tabulares: `VistaGaleria.tsx` lleva `tabular-nums` en sus propiedades y los eventos de `VistaCalendario.tsx` también. |
| **H-10** (media) Convenciones backend-base ↔ frontend-base | **Corregido** | `package.json`: `"verify": "npm run verificar"`, así `make verify-wf3-frontend` corre build, tsc, C3-03, U5-08, manifiesto y e2e. `playwright.config.ts`: si falta `FH_CHROMIUM`, usa el `chrome-headless-shell` más nuevo de `~/.cache/ms-playwright` (hoy la revisión 1217). Sin eso, `npm run verify` llamado desde `make` fallaría al no encontrar la revisión 1223, porque no hay red para descargarla. |

Otros cambios que exigieron estos arreglos:
- `pruebas/apoyo/api-falsa.ts`: nueva opción `fallaServidorGastos`, que responde 500 al POST de gasto sin guardar nada.
- `pruebas/cola.spec.ts`: casos nuevos de 500 y 401.
- `pruebas/apoyo/medidas.ts`: mediciones de área de toque, íconos y dinero tabular.
- `Grafica.tsx`: el «Ver como tabla» (`summary`, que medía 346 × 20) lleva `coarse:py-3` y mide 44 px de alto.
- `VistaCalendario.tsx`: los eventos del mes (45 × 16) llevan `coarse:min-h-11 coarse:py-3`.

Los dos últimos los encontró la prueba nueva. Están en archivos de la base y entran en el criterio de H-03.

## Archivos

Nuevos:
- `frontend/src/componentes/Icono.tsx`
- `frontend/pruebas/apoyo/medidas.ts`
- `frontend/pruebas/medidas.spec.ts`
- `docs/estado/wf3/correccion-frontend-core-r1.md`

Modificados en la base (frontend-base):
- `frontend/src/cola/cola.ts`
- `frontend/src/shell/FranjaSinConexion.tsx`
- `frontend/src/app/App.tsx`
- `frontend/src/componentes/`: `Avisos.tsx`, `Campo.tsx`, `Chip.tsx`, `Estados.tsx`, `TarjetaMetrica.tsx`, `Grafica.tsx`
- `frontend/src/shell/`: `BarraLateral.tsx`, `Barras.tsx`, `MenuComandos.tsx`
- `frontend/src/vistas-bd/`: `CeldaEditable.tsx`, `ListaOrdenable.tsx`, `VistaCalendario.tsx`, `VistaTabla.tsx`, `VistaBD.tsx`, `VistaGaleria.tsx`
- `frontend/package.json`
- `frontend/playwright.config.ts`
- `frontend/pruebas/apoyo/api-falsa.ts`
- `frontend/pruebas/cola.spec.ts`

Modificados en las áreas (solo tamaño de íconos y `tabular-nums`):
- `frontend/src/pantallas/inicio/`: `Calendario.tsx`, `P01Inicio.tsx`, `Resumen.tsx`, `comun.tsx`
- `frontend/src/pantallas/gastos-presupuesto/`: `P02Gastos.tsx`, `P04Gasto.tsx`, `P05Ticket.tsx`, `P06Presupuesto.tsx`, `P07Flujo.tsx`, `P08Ingresos.tsx`, `P09Ahorro.tsx`, `comun.tsx`
- `frontend/src/pantallas/comidas-compras/`: `P10Plan.tsx`, `P11Recetario.tsx`, `P13Lista.tsx`, `P14Canasta.tsx`, `P15Productos.tsx`, `P16Inventario.tsx`

Notas de corrección en reportes:
- `docs/estado/wf3/front-comidas-compras.md` (afirmación de 44 px, H-05)
- `docs/estado/wf3/front-gastos-presupuesto.md` (solo 3 controles medidos, H-04)

## Verificación (comandos propios, salidas resumidas)

- `npx vite build`: exit 0, precache 80 entradas.
- `npx tsc --noEmit`: **0 errores en los archivos de este encargo**. Al cierre había 10 errores en `pruebas/creditos-salud-reportes/datos.ts`. Ese archivo es de otra área y otro agente lo modificó a las 03:52, durante esta corrección. A las 03:53 había además errores en `pruebas/ajustes-a` y `ajustes-b` que ya se resolvieron.
  - Mi `npm run build` de las 03:49 dio exit 0 con tsc limpio.
- **C3-03**: `npm run verificar:rutas` → «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- **U5-08**: `npm run verificar:hosts` → «hosts externos que se pedirían = 0» (solo literales que no son peticiones). Manifiesto JSON válido.
- **Playwright del grupo**: `npx playwright test pruebas/inicio pruebas/gastos-presupuesto pruebas/comidas-compras pruebas/cola.spec.ts pruebas/shell.spec.ts pruebas/login.spec.ts pruebas/medidas.spec.ts` → **63 passed** (1.9 min; móvil 412×915 táctil y escritorio 1440×900). Antes eran 57 y se añadieron 6.
- **U5-13 paso 7 (H-01)**:
  - 422: pasa como antes.
  - 500: «franja visible, «1 en cola», aviso «No se pudo enviar · Revisar»; registros 0». Tras recargar sigue «1 en cola» con «No se pudo enviar». Con el servidor recuperado, «Revisar» → «Reintentar envío» da 1 registro y la franja desaparece.
  - 401: aviso «No se pudo enviar · Revisar», P-00, y al volver a entrar la franja cuenta «1 en cola». «Reintentar envío» da 1 registro.
- **U5-13 pasos 1–3 y 6**: «2 en cola» sobrevive a recargar. Claves de 36 caracteres, `captured_at` con desplazamiento. El reintento forzado responde 200 sin duplicar.
- **U5-09**: «toques para registrar un gasto (sin conexión): 3».
- **U5-06 · iconos (H-02)** (`medidas.spec.ts`, DOM a 412 táctil):
  - Inicio: 76 íconos en 4 vistas, **0 fuera de 20/24 px o trazo 1.5**.
  - Gastos y presupuesto: 169 íconos en 15 vistas, **0 fuera**.
  - Comidas y compras: 100 íconos en 10 vistas, **0 fuera**.
  - Antes del arreglo la misma prueba contó 56 íconos de 16 px con trazo dibujado de 1.2, y uno de 11.4 × 20.
- **U5-06 · objetivos táctiles (H-03, H-04, H-05)**: **0 controles menores de 44 × 44** en todas las vistas medidas.
  - Inicio: 83 controles. Gastos y presupuesto: 341. Comidas y compras: 239.
  - P-13: dibujo 20×20, área de toque (etiqueta) 44×44. Tocar la esquina marca el artículo.
- **U5-06 · dinero tabular (H-06)**: **0 importes sin `tabular-nums`** en las 29 vistas.
- **H-10**: `grep '"verify"' frontend/package.json` → presente.
  - `make verify-wf3-frontend` → ver «Resultado de make» al final.
- **`make verify-all GOOSE=…`**: **falla en `verify-go`**, en gofmt de `backend/internal/reportes/reportes_test.go` y `vistas.go`.
  - Esos archivos son de otro agente: los modificó a las 03:32 y 03:37.
  - No vienen de este encargo, que no toca `backend/`. verify-all no incluye la interfaz.

## Supuestos

- **HECHO, no supuesto**: la regla de tamaños sale literal de DESIGN §5.2: 20 px junto a texto; 24 px en barras, en botones de solo icono y en estados vacíos (junto a títulos de 20 px o más). No se cambiaron los íconos que ya medían 20 px (`size-5`), aunque estuvieran en una barra: cumplen «20/24» y no estaban en el hallazgo.
- La etiqueta de 44 × 44 que envuelve la casilla cuenta como su área de toque: DESIGN §4.6 dice «aunque el dibujo sea menor». La prueba lo demuestra tocando la esquina.
- El dueño nominal de H-04 y H-05 (el enlace) era el área, pero el control vive en archivos de la base (`VistaTabla.tsx`, `Barras.tsx`). Se corrigió ahí, una sola vez para todas las pantallas.

## Bloqueos

Ninguno nuevo de este encargo. No se pidió ningún valor que no estuviera decidido.

## Pendientes

- Bajos de la revisión que no entraban en este encargo:
  - H-07: la API falsa responde con página a GET de cosas que no son listas, y P-16 pinta `NaN`.
  - H-08: medidas sueltas en `Layout.tsx` y `Hoja.tsx`.
  - H-09: Inicio reimplementa `TarjetaMetrica` y `ConDatos`.
  - H-11: `cola.spec.ts` no comprueba que `captured_at` sea anterior al envío ni que el reintento devuelva el mismo `id`.
- `make verify-wf3-frontend` completo depende de que `pruebas/creditos-salud-reportes/datos.ts` vuelva a compilar. Es trabajo en curso de otro agente.
- Las historias de Storybook (`*.stories.tsx` y `.storybook/preview.tsx`) aún usan `size-4` y no tienen `ProveedorIconos`. No son pantallas de la app y no entraban en el hallazgo. Se dejan para quien sea dueño de `.storybook`.

## Propuestas para Axel

- **PROPUESTA (H-01)**: el aviso de 5xx o 401 sale **una vez por captura**. Los reintentos automáticos de cada 30 s, que ya eran PROPUESTA de K31, no lo repiten mientras la falla siga. Sí lo repite el «Reintentar envío» manual. Así se evita un aviso cada 30 s con el servidor caído. La franja sigue contando todo el tiempo. Si prefieren otro ritmo, se cambia en una línea (`cola.ts`).
- En el calendario mensual a 412 px, los eventos ahora miden 44 px de alto con pantalla táctil (antes 16 px), lo que alarga la cuadrícula. Si Axel prefiere una agenda (lista) en móvil en lugar de la cuadrícula, es una decisión de diseño (U1 §3) que no se tomó aquí.

## Paquetes instalados

Ninguno.

## Consumo aproximado

~45 acciones · ~170 KB leídos · ~40 KB escritos.

## Resultado de make

`make verify-wf3-frontend`, que ahora sí corre `npm run verify` (03:57–04:00):
- build y tsc: exit 0. Para entonces el otro agente ya había resuelto los errores de tipos de creditos-salud-reportes.
- C3-03: «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- U5-08: «hosts externos que se pedirían = 0» y «hosts pedidos: localhost; externos: 0».
- e2e de todas las áreas: **122 passed, 1 failed**.
  - La que falla es `pruebas/ajustes-b/ajustes-b.escritorio.spec.ts` «P-26.09 Notificaciones», solo en escritorio.
  - El agente de ajustes-b estaba editando esa pantalla y su prueba durante la corrida: `P2609Notificaciones.tsx` cambió a las 03:56:45.
  - Al relanzarla aislada falló en otra aserción (`getByText('Xime')`), porque la prueba cambió entre una corrida y otra. Es trabajo en curso de otra área.
  - Ningún archivo de este encargo interviene en esa aserción.
