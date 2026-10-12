# REVISIÓN frontend de áreas · ronda 4 · WF-3a (12-oct-2026)

Revisor: `revisor-frontend-areas-r4`. Soy independiente: no escribí nada de lo que reviso ni corregí nada. Trabajo revisado: front-creditos-salud-reportes, front-ajustes-a, front-ajustes-b y front-variantes, con las correcciones `frontend-base-fix`, `front-ajustes-a-fix`, `front-ajustes-b-fix` y `front-creditos-salud-reportes-fix`. Criterios: U5-06, U5-07 y U5-09 (`docs/criterios/U5.md`), C3-03 (`C3.md`), DESIGN.md, U1 y U4.

Todas las mediciones son mías y se hicieron en `datos/tmp/rev-fa-r4/`, que ya borré. Usé:
- un build propio (`vite build --outDir datos/tmp/rev-fa-r4/dist`), servido con `vite preview` en el puerto 4797;
- una configuración temporal de Playwright;
- sondas propias: medidas de H-03 y H-10 en claro y en oscuro, un barrido de 30 rutas × 2 temas × 2 proyectos y una sonda de error de servidor en 02;
- Storybook reconstruido y `contraste-variantes.mjs`.

No usé `frontend/dist` ni el puerto 4173.

## Veredicto: NO APROBADO · 0 altos y 2 medios

| Severidad | Abiertos | Ids |
|---|---|---|
| alta | 0 | — (H-03 cerrado) |
| media | 2 | R4-01 (nuevo, front-ajustes-a) y R4-02 (nuevo, front-ajustes-b): en Ajustes hay errores de guardado que no se muestran (U4 §0.3 y V-17) |
| baja | 1 | R3-01 (`ajv` sin declarar en `frontend/package.json`; sigue abierto) |

**HECHO.** Los cuatro hallazgos de la ronda 3 quedaron cerrados: H-03, H-07, H-08 y H-10. Los medí yo, y el detalle está en la tabla de abajo. También cerré los bajos R2-05, R2-06, R2-07 y R3-02.

Lo que impide aprobar es nuevo y no venía de la ronda 3. Al confirmar un monto en P-26.02, un `409 budget_overlap` del servidor termina en una promesa rechazada sin capturar (`pageerror`): la hoja se queda abierta y no muestra ningún mensaje, ni V-17 ni «No se pudo guardar». En P-26.13 hay dos botones con el mismo patrón. El defecto ya estaba en HEAD `a35451e`, donde `S02Presupuesto.tsx` tiene 0 `catch`.

## Tabla criterio → estado → evidencia

| Criterio | Estado | Evidencia (mi corrida) |
|---|---|---|
| C3-03: build | CUMPLE | `npx tsc --noEmit`: rc 0. `npx vite build --outDir …/dist`: «✓ built», 81 entradas en precache. `tsc -p .storybook/tsconfig.json`: rc 0. |
| C3-03: rutas | CUMPLE | `npm run verificar:rutas`: «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». Cruce propio: las operaciones que U4 cita en P-26.01 a P-26.14 (140 citas, contra 215 operaciones de los extractos) tienen todas su ruta literal en `frontend/src`. |
| U5-08 (referencia) | CUMPLE | `hosts-externos.mjs` sobre mi dist: «hosts externos que se pedirían = 0». Solo aparecen literales que no son peticiones (tailwindcss.com, la licencia de Inter y bit.ly en workbox). |
| Pruebas Playwright de las 3 áreas | CUMPLE | `npx playwright test --config <tmp> ajustes-a ajustes-b creditos-salud-reportes`: **115 passed, 1 skipped, 0 failed** (2.2 min). La omitida es la medida táctil en escritorio. |
| U5-06: 44 px (H-03 punto 1) | CUMPLE | Sonda propia en /ajustes/05 a 412 táctil: agenda presente (`data-prueba="agenda"`) y evento «Walmart · Sin intereses» de **344×44**, en claro y en oscuro (fondos rgb(250,250,251) y rgb(34,35,37)). `controlesChicos` = [] entre 41 controles. A 1440 con puntero fino es la cuadrícula del mes y el evento mide 108×24. |
| U5-06: «Ver como tabla» (H-03 punto 2) | CUMPLE | En /creditos, /creditos/1 y /salud con `?sub=tendencias`, el `<summary>` mide **124×36** a 1440 y **124×44** a 412 táctil, en los dos temas. |
| U5-06: descripción emergente por token (H-03 punto 3) | CUMPLE | Paso el puntero por la gráfica y leo el estilo calculado. En oscuro: fondo rgb(54,54,58), borde rgb(77,78,83) y texto rgb(242,242,244), iguales a `--ui-surface-raised`, `--ui-border-divider` y `--ui-text-default` resueltos en la misma página. En claro: rgb(250,250,251), rgb(229,230,233) y rgb(54,54,58), también iguales. Radio 10 px (`--radius-md`), `tabular-nums` y la sombra de `--ui-shadow-1` (`none` en oscuro, como define el token). Mismo resultado en las 3 rutas y en los 2 proyectos. |
| U5-06: foco sin transición, series de C2 y deshabilitado (H-10) | CUMPLE | `storybook build` en tmp y `node scripts/contraste-variantes.mjs --salida <tmp>`: 37 historias, 148 páginas y 5 054 mediciones. Resultado: «FALLAS nuevas: 0», «hallazgos conocidos: 0 de 5 se reproducen» (FV-01 a FV-05 CERRADO) y «RESULTADO: OK». En la app, ningún botón de las 3 rutas tiene `outline` en `transition-property` (`Boton.tsx:7`: `transition-[color,background-color,border-color]`). |
| U5-06: barrido de las pantallas revisadas | CUMPLE | Barrido propio de /ajustes/01 a /ajustes/14, /ajustes y 15 rutas de P-17 a P-25 (con sus subvistas), en claro y en oscuro, a 412 táctil y a 1440: **120 de 120 páginas sin defectos**. En ninguna hubo controles bajo 44 px en táctil, controles visibles sin nombre, desborde horizontal, dinero sin `tabular-nums` ni franjas «No se pudo…». |
| U5-06: colores en el código | CUMPLE | `grep -E "#[0-9a-fA-F]{3,8}\b\|rgb\(\|hsl\(\|oklch\(\|dark:"` sobre las 3 carpetas de pantallas, `src/vistas-bd`, `src/componentes`, `.storybook` y las historias: 0 en el código. El único resultado es el texto de FV-03 en `.storybook/hallazgos.json`. Los tokens que usa `Grafica.tsx` existen en `tokens.css`. |
| U5-07: idioma, fechas y dinero | CUMPLE | En el mismo barrido de 120 páginas, sobre el `innerText`: 0 fechas ISO, 0 fechas con barras, 0 importes fuera de `^−?\$\d{1,3}(,\d{3})*\.\d{2}$` y `textosEnIngles` = [] (también cubre `aria-label`, `title` y `alt`). Con el bloque de `Ctrl+/` abierto, las pruebas de 02, 04, 05, 07, 09, 13 y 14 dan `textosEnIngles` = []. |
| U5-09 (referencia) | CUMPLE | Las pruebas de las áreas validan el tráfico contra C1 (`erroresTrafico` = [] en `afterEach`) y pasan. |
| U4 §0.2 punto 2 y §0.7: `Ctrl+/` con vista previa sumada (H-07, H-08) | CUMPLE | Las 7 pruebas pasan en los dos proyectos. Leí su código: marcan las casillas, ponen el foco en la vista, pulsan `Control+/`, comprueban el bloque `vista-previa` con cifras, `Escape` sin ningún PATCH y `Enter` con un PATCH por fila y su cuerpo exacto. **02**: «Gastos fijos de oct-2026 $699.00 → $0.00», las partes por persona y el aviso de P-07. **04**: «2 categorías pasan a Conjunto 50/50», el ejemplo de $199.00 y el gasto fijo que cambia. **05**: «vigentes hoy» y R5-3. **07**: tiempos aptos; si una receta se queda sin tiempo, no aplica. **09**: «Alertas activas 4 de 4 → 2 de 4». **13**: «Total reservado … $1,500.00 → $0.00». **14**: «1 pasa a Reponer…». El manejador está en `VistaBD.tsx:44` (con el foco dentro de la vista) y en `ajustes-b/varias.tsx` para las tarjetas de 09. `PlannedItemUpdate` de C1 sí tiene `is_active`, así que el cuerpo de 02 es válido. |
| U4 §0.2 «Antes → Después» en tableros (04 y 14) | CUMPLE | Las pruebas «P-26.04 tablero» y «P-26.14 tablero» pasan: vista previa al soltar, Esc sin PATCH y Enter con PATCH. |
| U4: claves `reports.*` y `alerts.*` | CUMPLE | Pasan «R2-06 · P-26.08» (`PATCH /api/v1/settings/reports` `{parameters:{'reports.daily_send_time':'22:15'}}`) y «R2-06 · P-26.09» (`PATCH /api/v1/settings/alerts` `{parameters:{'alerts.send_time':'08:30'}}`), con el tráfico validado contra C1. |
| U4 §0.3 y V-17: error al guardar | **NO CUMPLE** (medio) | R4-01 y R4-02, abajo. |
| U4: las 14 secciones | CUMPLE salvo R4-01 y R4-02 | Ver la tabla por sección. |
| U1: P-17 a P-25 | CUMPLE | Pasan las pruebas de P-17 a P-25 (13 por proyecto), entre ellas «P-17 N-193» (aviso A/B desde `getCredit`) y «P-18 N-208» (`listCalendarEvents` filtrado por línea), además de los flujos F-4, F-5 y F-6. `python3 -I scripts/c3-integracion.py` (solo lectura): «implementado sí 41/42»; el único «no» es P-03, que no es de estas áreas. Mi barrido de /creditos/1?sub=calendario: agenda en 412 y calendario en 1440, sin defectos. |
| front-variantes | CUMPLE | `tsc -p .storybook/tsconfig.json`: rc 0. Storybook: «build completed successfully». Contraste: «RESULTADO: OK». |
| `make verify-all` (referencia) | CUMPLE | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: rc 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |

## Tabla de las 14 secciones de Ajustes (U4)

Método:
- (a) Un script propio separa U4 por `## P-26.xx`, toma la primera columna de cada tabla de campos y busca cada etiqueta, o sus partes, en `src/pantallas/ajustes-*`, `src/textos`, `src/vistas-bd` y `src/componentes`. Las que no halló las revisé a mano con `grep`.
- (b) Cruzo las operaciones de C1 citadas en cada sección contra las rutas literales del código.
- (c) Uso el barrido y las pruebas.

| Sección | Campos U4 (hallados) | Operaciones sin ruta | Presentes | Faltan o fallan |
|---|---|---|---|---|
| P-26.01 Cuenta y roles | 13/13 (4 a mano: Karen, cerrar sesión, App instalada, certificado) | 0 de 11 | Sesiones, K44 con vista previa | — |
| P-26.02 Presupuesto | 20/20 (7 a mano). El texto V-17 está en `textos.ts:17`, pero **no se usa** | 0 de 12 | Alcance con «Confirmar», `Ctrl+/` con la vista previa 3 (h), cifras tabulares | **R4-01**: un 409 o un fallo al confirmar el monto no muestra nada |
| P-26.03 Ingresos | 7/7 | 0 de 4 | Todos | — |
| P-26.04 Categorías | 8/8 | 0 de 8 | Reparto con vista previa en celda, en tablero y en `Ctrl+/` | — |
| P-26.05 Créditos y condiciones | 49/49 (2 a mano: eliminar tasa o promoción) | 0 de 20 | Vistas previas 6, 7 y 10 (R2-05), N-284 con `getCredit`, `Ctrl+/` en promociones, agenda con eventos de 44 px | — |
| P-26.06 Tiendas y precios | 21/21 | 0 de 14 | Punto de equilibrio «Vista previa pendiente · PROVISIONAL» con su bloqueo citado (R2-07) | — |
| P-26.07 Recetas | 9/9 (2 a mano) | 0 de 14 | Alt+N, duplicar, `Ctrl+/` en tiempos aptos | — |
| P-26.08 Reportes y frecuencias | 14/14 (5 a mano: «Día» y «Hora» por reporte) | 0 de 14 | 5 claves `reports.*` con prueba de edición | — |
| P-26.09 Notificaciones | 24/24 (2 a mano) | 0 de 14 | 6 claves `alerts.*` con prueba, `Ctrl+/` en alertas | — |
| P-26.10 Apariencia | 2/2 | 0 de 0 | Tema al instante | — |
| P-26.11 Datos y respaldos | 20/20 (2 a mano: «Próximo: …» y usos de la IA en `P2611Datos.tsx:235`) | 0 de 13 | Respaldos, IA | — |
| P-26.12 Personas y reparto | 12/12 | 0 de 4 | Papeles, reglas, alcance | — |
| P-26.13 Salud y metas | 24/24 | 0 de 8 | Perfil, metas, `Ctrl+/` con el total reservado | **R4-02**: «Usar el de la casa» y el botón de la meta sugerida no muestran un fallo |
| P-26.14 Inventario | 9/9 | 0 de 4 | Modos, estado con texto, tablero y `Ctrl+/` con «{k} pasan a Reponer…» | — |

## Hallazgos abiertos

| Id | Sev. | Criterio | Ruta | Dueño | Evidencia mía | Corrección sugerida |
|---|---|---|---|---|---|---|
| R4-01 | media | U4 P-26.02 (fila «Si se traslapa con otro monto» → V-17, `createBudget` responde conflicto) y U4 §0.3 «Error al guardar» | `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx:70-76` (`confirmar` sin `try/catch`) y `:93` (`onClick={() => void confirmar()}`) | front-ajustes-a | Sonda propia, en movil y en escritorio: edito el monto de Despensa, tecleo 4500, `Enter` y «Confirmar», con `POST`/`PATCH /api/v1/budgets*` respondiendo `409 {code:'budget_overlap', message:'Ya hay un monto…'}`. Resultado: `pageerror` con el mensaje del servidor (promesa rechazada sin capturar), la hoja sigue abierta y **ni el mensaje del servidor, ni V-17, ni «No se pudo guardar» aparecen**. `textos.ts:17` define `v17` y nadie lo usa (`grep v17` en `src`: solo la definición). En HEAD `a35451e` el archivo ya tenía 0 `catch`. | Capturar el error en `confirmar`. Si es `ErrorApi` con `codigo === 'budget_overlap'`, mostrar `p.v17` en la hoja (debajo del monto, con `circle-alert` y `text-role-error`) y dejarla abierta. Para cualquier otro error, mostrar el texto común de §0.3 «No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar» sin perder lo capturado. Agregar una prueba con la respuesta 409 de C1 que compruebe el mensaje y 0 `pageerror`. |
| R4-02 | media | U4 §0.3 «Error al guardar» (rige las 14 secciones) | `frontend/src/pantallas/ajustes-b/P2613Salud.tsx:112` y `:124` (`onClick={() => void pedir('patch', '/api/v1/health-profiles/{…}', …).then(recargar)}`) | front-ajustes-b | Lo vi leyendo el código, no lo medí. El patrón es el mismo de R4-01, cuya consecuencia sí medí: sin `catch`, un fallo del PATCH no muestra nada. En la misma pantalla, `:95` sí pasa por `g.ejecutar`. | Pasar los dos botones por el ayudante con estado de error que ya usa la sección (`g.ejecutar` o el de `comun.tsx:87`) y agregar una prueba con una respuesta 500 que compruebe «No se pudo guardar…». |
| R3-01 | baja | Higiene de pruebas | `frontend/package.json` | frontend-base (dueño de la configuración) | `"ajv"` no está en `devDependencies`. Solo existe como dependencia transitiva en `node_modules`. | Agregar `"ajv": "8.20.0"` a `devDependencies` (solicitud 4 de `frontend-base-fix`). |

## Estado de los hallazgos de la ronda 3

| Hallazgo | Sev. | Estado | Evidencia mía |
|---|---|---|---|
| H-03 (evento del calendario, «Ver como tabla», descripción emergente) | alta | **CERRADO** | Evento de 344×44 en agenda (412 táctil, claro y oscuro). Summary de 36 px con puntero fino y 44 en táctil. Descripción emergente igual a `--ui-surface-raised`, `--ui-border-divider` y `--ui-text-default` en los dos temas. |
| H-07 (`Ctrl+/` en 02, 04, 05 y 07) | media | **CERRADO** | 4 pruebas pasan en los 2 proyectos con `Control+/`, la vista previa sumada, Esc sin PATCH y Enter con un PATCH por fila. La selección de filas está en `VistaTabla` y el atajo en `VistaBD.tsx:44`. |
| H-08 (`Ctrl+/` en 09, 13 y 14) | media | **CERRADO** | 3 pruebas pasan. 09 usa `GrupoVarias` (`ajustes-b/varias.tsx`) con el mismo `EdicionVarias` de la base. |
| H-10 (FV-01, FV-04 y FV-05) | media | **CERRADO** | 0 de 5 hallazgos conocidos se reproducen y 0 fallas nuevas. Ningún botón tiene `outline` en `transition-property`. |
| R2-05, R2-06 y R2-07 | baja | **CERRADOS** | Las pruebas «R2-05», «R2-06 · P-26.08», «R2-06 · P-26.09» y «R2-07» pasan. |
| R3-01 (`ajv`) | baja | **ABIERTO** | Ver la tabla de hallazgos. |
| R3-02 (celda sin nombre) | baja | **CERRADO** | En el barrido de 120 páginas no hay ningún control visible sin nombre. |

## PREGUNTA para el orquestador o Axel

R4-01 y R4-02 son nuevos y no estaban en la ronda 3. Los dos se corrigen en pocas líneas dentro del alcance de cada área (front-ajustes-a y front-ajustes-b). ¿Se encargan a esas dos áreas con una revisión puntual después, o se aceptan como pendientes? No lo decido yo.

## Lo que no corrí

- `make verify-wf3-frontend`, porque escribe en `frontend/dist` y usa el puerto 4173. En su lugar corrí lo mismo de forma aislada: tsc, build, rutas, hosts, Playwright de las áreas y variantes.
- La suite completa de Playwright de todas las áreas: solo corrí las tres áreas que me tocan.
- Storybook se reconstruyó en `frontend/storybook-static` porque `contraste-variantes.mjs` lo lee de ahí fijo. Es una salida de build ignorada por git; ya existía antes y la reemplacé con el mismo contenido de fuente.
