# REVISIÓN frontend de áreas · ronda 3 de 3 · WF-3a (12-oct-2026)

Revisor: `revisor-frontend-areas-r3` (independiente; no escribí lo revisado ni corregí nada). Trabajo revisado: front-creditos-salud-reportes, front-ajustes-a, front-ajustes-b y front-variantes, con la corrección `correccion-frontend-areas-r2`. Criterios: U5-06, U5-07 y U5-09 (`docs/criterios/U5.md`), C3-03 (`C3.md`), DESIGN.md, U1 y U4.

Todo lo de abajo lo medí yo, en `datos/tmp/rev-fa-r3/` (ya borrado): build propio, preview aislado en el puerto 4797, configuración temporal de Playwright, copias instrumentadas de las pruebas (graban el tráfico `/api/` y miden el estado final de cada página en claro y en oscuro), un validador ajv 2020 propio contra `docs/contratos/openapi.yaml` y Storybook construido en tmp. No usé `frontend/dist` ni el puerto 4173 para las pruebas de áreas.

## Veredicto: NO APROBADO · 1 alto y 3 medios

| Severidad | Abiertos | Ids |
|---|---|---|
| alta | 1 | H-03 (sigue abierto; frontend-base) |
| media | 3 | H-07 y H-08 (solo queda `Ctrl+/`; dependen de frontend-base), H-10 (frontend-base) |
| baja | 5 | R2-05, R2-06, R2-07 (siguen), R3-01 y R3-02 (nuevos) |

**HECHO.** Las áreas cerraron todo lo que dependía solo de ellas: R2-01, R2-02, H-09, R2-03 (en las pantallas), R2-04 y los dos tableros con vista previa (04 y 14). Lo que queda abierto está en archivos de frontend-base (`src/componentes/`, `src/vistas-bd/`), que no cambiaron desde la ronda 2 (fechas de modificación 03:47–04:13, antes de esa revisión), o necesita algo de la base primero (`Ctrl+/` requiere selección de filas en VistaBD).

## Tabla criterio → estado → evidencia

| Criterio | Estado | Evidencia (mi corrida) |
|---|---|---|
| C3-03: build | CUMPLE | `npx tsc --noEmit`: 0 errores. `npx vite build --outDir datos/tmp/rev-fa-r3/dist`: «✓ built», PWA con 80 entradas en precache. |
| C3-03: rutas | CUMPLE | `npm run verificar:rutas`: «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». En el tráfico grabado (883 respuestas `/api/` por esquema), 0 operaciones fuera de C1. |
| U5-08 (referencia) | CUMPLE | `node scripts/hosts-externos.mjs` sobre mi dist: solo literales que no son peticiones (w3.org, react.dev, tailwindcss.com, licencia de Inter, bit.ly en workbox). |
| Pruebas Playwright de las 3 áreas | CUMPLE | `npx playwright test ajustes-a ajustes-b creditos-salud-reportes`: **85 passed, 1 skipped** (la medida táctil, que solo corre en `movil`). Las copias instrumentadas: 85 passed en claro y 85 passed en oscuro. |
| Datos de prueba contra el esquema C1 | CUMPLE | Validador propio (ajv 2020, sin reutilizar `validar-c1.ts`) sobre los cuerpos literales que mandan las respuestas falsas de ajustes-a y ajustes-b (417 registros, 53 operaciones distintas): **0 fallas**. Las 12 escrituras validan con su código: `PATCH settings/security` 200, `PATCH budgets/{id}` 200, `PATCH credits/{id}` 200, `POST credit-purchases/preview` 200, `POST promotions` 201, `POST recipes/{id}/duplicate` 201, `POST recipes` 201, `PATCH categories/{id}` 200, `DELETE telegram-chat-candidates/{id}` 204, `PATCH settings/data` 200, `PATCH ai/providers/{provider}` 200, `PATCH inventory-items/{id}` 200. Las escrituras de créditos y reportes (`POST scheduled-payments/{id}/payments` 201, `POST reports/monthly` 202, `POST scenario-comparisons` 200) también validan en el tráfico grabado. |
| U5-06: objetivos táctiles de 44 px | **NO CUMPLE** (alto, frontend-base) | `controlesChicos` en el estado final de cada prueba en 412 táctil, en claro y en oscuro: un solo control bajo 44, el evento «Walmart · Sin intereses» de VistaCalendario, de **40×44** en /ajustes/05 (H-03). Las flechas «Semana anterior/siguiente» de /ajustes/06 ya no aparecen (R2-02 cerrado). |
| U5-06: «Ver como tabla» con puntero fino | **NO CUMPLE** (alto, frontend-base) | Sondeo propio en /creditos, /creditos/1 y /salud (`?sub=tendencias`) a 1440: el `<summary>` mide **20 px** en claro y en oscuro (44 en táctil). Mismo dato que FV-02. |
| U5-06: colores por token en la descripción emergente | **NO CUMPLE** (alto, frontend-base) | Paso el puntero por la gráfica y leo el estilo calculado del recuadro de ECharts: fondo `rgb(255,255,255)`, borde `rgb(183,185,190)` y texto `rgb(109,110,115)` **también en oscuro**, donde `--ui-surface-raised` es `#36363a` y `--ui-text-default` `#f2f2f4`. `Grafica.tsx:48` sigue con `tooltip: { trigger: 'axis', valueFormatter }`, sin colores. |
| U5-06: foco sin transición, series de C2 y deshabilitado (§7.0) | NO CUMPLE (medio, frontend-base) | `contraste-variantes.mjs` sobre Storybook construido en tmp: 37 historias, 148 páginas, 5 230 mediciones, «FALLAS nuevas: 0», y los 5 hallazgos conocidos se reproducen: FV-01 128, FV-02 12, FV-03 384, FV-04 52 y FV-05 16. En la app, `transition-property` del Botón incluye `outline-color` (`Boton.tsx:6`, `transition-colors`). |
| U5-06: cifras tabulares | CUMPLE | `dineroSinTabulares` en el estado final de todas las pruebas, en claro y en oscuro: 0 casos. `S02Presupuesto.tsx:99` envuelve el aviso del alcance en `tabular-nums` y su prueba lo mide con el diálogo abierto. |
| U5-06: hex en el código | CUMPLE | `grep -E "#[0-9a-fA-F]{3,8}\b\|rgb(\|hsl(\|oklch(\|dark:"` sobre las 3 carpetas de pantallas, `.storybook/*.ts*` y las historias: 0. |
| U5-07: idioma, fechas y dinero | CUMPLE | En el `innerText` final de cada prueba (85 en claro y 85 en oscuro): 0 fechas ISO, 0 fechas con barras, 0 importes fuera de `$1,234.56` (cada `$…` contra `^−?\$\d{1,3}(,\d{3})*\.\d{2}$`), 0 palabras de una lista de inglés de interfaz y 0 claves técnicas visibles (solo «finanzas.local», que es el nombre del servidor). |
| U5-09 (referencia, base) | CUMPLE | `cola.spec.ts` «captura sin conexión…»: «U5-09 toques para registrar un gasto (sin conexión): 3», 1 passed. |
| DESIGN §7 Botón, nombre visible o `aria-label` | CUMPLE en las pantallas | Botones, enlaces, `summary`, pestañas e interruptores visibles sin nombre en el estado final de todas las pruebas: 0 (antes 6 en /12, 4 en /13 y 6 en /14). La base sigue pintando `''` por omisión (R3-02, bajo). |
| U4: las 14 secciones | NO CUMPLE (medio) | Ver la tabla por sección. Solo falta `Ctrl+/` (H-07, H-08). |
| U4 §0.2 «Antes → Después» en tableros | CUMPLE | 04: al soltar «Despensa» en «Conjunto 50/50» queda en la columna nueva con la vista previa ($99.50), «Aplicar» con el foco, Esc la regresa sin PATCH, Enter guarda `PATCH /categories/1 {default_cost_share_rule_id: 2}`. 14: al soltar «Arroz» en «Ritmo de consumo» muestra «Estado: Agotado» y «Falta el ritmo», Esc sin PATCH, Enter guarda `tracking_mode: consumption_rate`. Las dos pruebas pasan en móvil y escritorio. |
| U4: claves `reports.*` y `alerts.*` | CUMPLE (con nota baja) | P2608Reportes.tsx:29-31 maneja las 5 claves `reports.*` y guarda con `PATCH /api/v1/settings/reports` (l. 34). P2609Notificaciones.tsx maneja las 6 `alerts.*` y guarda con `PATCH /api/v1/settings/alerts` (l. 27). En el tráfico grabado no hay ningún PATCH a esas dos rutas: ninguna prueba las edita (R2-06). |
| U1: P-17 a P-25 | CUMPLE | Pasan las 8 pruebas de P-17 a P-25 (en 412 táctil y en 1440) y los 4 flujos F-4, F-5 y F-6. N-225 ya está: P20Comparador.tsx:39 lee `GET /api/v1/settings/optimizer`, precarga colchón, máximo de compras vivas, fecha de deuda en cero y los 4 pesos con «De Ajustes», avisa «Solo para esta comparación» si cambian (l. 45), muestra los pesos normalizados (l. 98) y el tope por quincena de cada persona o «Sin tope» (l. 104). P-22 es de la base. |
| front-variantes (su alcance) | CUMPLE | `tsc -p .storybook/tsconfig.json`: 0 errores. Storybook build: «completed successfully». Contraste: «RESULTADO: OK (sin fallas nuevas)», 0 problemas de corrida, 0 hosts externos. Los 5 FV-xx son de componentes de la base. |
| `make verify-all` (referencia) | CUMPLE | `make verify-all GOOSE=…/goose`: salida 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |

## Tabla de las 14 secciones de Ajustes (U4)

Método: un script cruza, por sección, (a) las etiquetas de las tablas de campos de U4 contra los textos de `src/pantallas/ajustes-*/` y `src/textos/`, y (b) las operaciones que nombra U4 contra las rutas literales en `frontend/src`. Las etiquetas que el script no halló las revisé a mano: todas están con otra forma (parametrizadas o compuestas). (c) Lo que mostraron las pruebas.

| Sección | Campos U4 hallados | Operaciones sin ruta en el código | Presentes | Faltan o fallan |
|---|---|---|---|---|
| P-26.01 Cuenta y roles | 17/17 (3 a mano: cerrar sesiones, parámetro) | `queryView` (U4 la da como alternativa: «`listSessions` (o `queryView`)») | Sesiones con sus alcances, K44 con vista previa | — |
| P-26.02 Presupuesto | 20/20 (5 a mano: alcance, «Solo este mes», traslape) | 0 | Alcance con «Confirmar», cifras tabulares | `Ctrl+/` (H-07) |
| P-26.03 Ingresos | 7/7 | 0 | Todos | — |
| P-26.04 Categorías | 8/8 | `queryView` (alternativa) | Reparto con vista previa en celda y en tablero, «Aplicar», Esc | `Ctrl+/` (H-07) |
| P-26.05 Créditos y condiciones | 53/53 (3 a mano: eliminar tasa o promoción, parámetro) | `queryView` (alternativa) | Vistas previas 6 y 7, alta de promoción, ←/→ | `Ctrl+/` (H-07). Evento del calendario de 40×44 (H-03, base). Vista previa 10 (R2-05, bajo) |
| P-26.06 Tiendas y precios | 24/24 | `queryView` (alternativa) | Calendario por día de visita con flechas de 44 px | Punto de equilibrio sin citar el bloqueo (R2-07, bajo) |
| P-26.07 Recetas | 10/10 (2 a mano) | `queryView` (alternativa) | Nueva receta (Alt+N) y duplicar | `Ctrl+/` (H-07) |
| P-26.08 Reportes y frecuencias | 14/14 (5 a mano: etiquetas «Día» y «Hora» por reporte) | 0 | 5 claves `reports.*` | Sin prueba que edite (R2-06, bajo) |
| P-26.09 Notificaciones | 24/24 (2 a mano) | 0 | 6 claves `alerts.*`, umbral por categoría, DELETE con 204 | `Ctrl+/` (H-08). Sin prueba que edite (R2-06, bajo) |
| P-26.10 Apariencia | 2/2 | 0 | Tema al instante | — |
| P-26.11 Datos y respaldos | 20/20 (2 a mano: próximo respaldo, usos permitidos) | 0 | Respaldos, parámetro con vista previa, IA con acción crítica | — |
| P-26.12 Personas y reparto | 12/12 | 0 | Papeles, reglas, V-16, «Sin valor» | — |
| P-26.13 Salud y metas | 24/24 | 0 | Perfil, meta sugerida, metas, «Sin valor» | `Ctrl+/` (H-08) |
| P-26.14 Inventario | 9/9 | 0 | Modos, estado con texto, días restantes, tablero con vista previa | `Ctrl+/` (H-08) |

## Hallazgos abiertos

| Id | Sev. | Criterio | Ruta | Dueño | Qué pasa (medido por mí) | Corrección sugerida |
|---|---|---|---|---|---|---|
| H-03 | alta | U5-06: 44 px y colores por token en las 4 combinaciones | `frontend/src/vistas-bd/VistaCalendario.tsx:30`; `frontend/src/componentes/Grafica.tsx:48` y `:78` | frontend-base | (1) El evento del calendario mide 40×44 en /ajustes/05 con 412 táctil, en claro y en oscuro. (2) «Ver como tabla» mide 20 px con puntero fino. (3) La descripción emergente usa `#FFFFFF`, `#B7B9BE` y `#6D6E73` también en oscuro. Los tres archivos no cambiaron desde la ronda 2. | (1) `coarse:min-h-11 coarse:min-w-11 shrink-0` en el botón del evento, o una agenda en `phone:`. (2) `inline-flex min-h-6 items-center coarse:min-h-11` en el `<summary>`. (3) `tooltip: { backgroundColor, borderColor, textStyle: { color } }` leídos de `--ui-surface-raised`, `--ui-border-default` y `--ui-text-default`, y releídos al cambiar el tema. |
| H-07 | media | U4 §0.2 punto 2 y §0.7 (`Ctrl+/` con vista previa sumada) | `frontend/src/vistas-bd/VistaBD.tsx`, `VistaTabla.tsx`; `frontend/src/pantallas/ajustes-a/` (S02, S04, S05, S07) | frontend-base, después front-ajustes-a | `grep` de `Ctrl+/`, `key === '/'` con `ctrlKey` en `frontend/src`: 0 manejadores (el único `'/'` es el menú de Layout.tsx:31, que excluye `ctrlKey`). VistaBD no tiene selección de filas. El tablero de 04 ya está. | Selección de filas y `Ctrl+/` en VistaBD/VistaTabla con `editarVarias: { propiedades; vistaPrevia(filas, propiedad, valor): Promise<Cambio[]>; guardar }`, bloque sumado, Enter aplica, Esc cancela; luego conectarlo en 02 (categoría y activo de gastos fijos), 04 (tipo, reparto, activa), 05 (promociones) y 07 (tiempos aptos). |
| H-08 | media | U4 §0.2 punto 2 y §0.7 | `frontend/src/vistas-bd/VistaBD.tsx`; `frontend/src/pantallas/ajustes-b/` (P2609, P2613, P2614) | frontend-base, después front-ajustes-b | Igual que H-07 para 09 (activar o desactivar alertas), 13 (metas) y 14 (modo o nivel, «{k} pasan a Reponer…»). El tablero de 14 ya está. | La misma pieza de la base y conectarla en 09, 13 y 14. |
| H-10 | media | U5-06 / DESIGN §7.0 y §4.5: foco sin transición, series de C2, deshabilitado | `frontend/src/componentes/Boton.tsx:6`, `Grafica.tsx`, `Campo.tsx`, `Chip.tsx` | frontend-base | FV-01 (128), FV-04 (52) y FV-05 (16) se reproducen. `transition-property` del Botón incluye `outline-color`. | `transition-[background-color]` en lugar de `transition-colors`; `emphasis: { disabled: true }` (o sin recolorear) en Grafica; estilo `disabled` de §7.0 en Campo y ChipSeleccionable (`frontend/.storybook/hallazgos.json` trae el detalle). |
| R2-05 | baja | U4 P-26.05, vista previa 10 | `frontend/src/pantallas/ajustes-a/S05Creditos.tsx:158-161` | front-ajustes-a | Sin cambios: toma `credit-purchases?limit=1` sin orden y la frecuencia de la línea (`c.payment_frequency`). | Usar la compra del préstamo con su `payment_frequency` y primera fecha, con orden o filtro. |
| R2-06 | baja | Cobertura de pruebas de `reports.*` y `alerts.*` | `frontend/pruebas/ajustes-b/ajustes-b.escritorio.spec.ts` | front-ajustes-b | 0 PATCH a `settings/reports` y `settings/alerts` en el tráfico de las pruebas. | Una prueba por operación que edite una clave y compruebe el cuerpo, con respuesta falsa válida en C1. |
| R2-07 | baja | Regla «PROVISIONAL citando bloqueo» | `frontend/src/pantallas/ajustes-a/S06Tiendas.tsx:128` | front-ajustes-a | El punto de equilibrio sigue diciendo «Vista previa pendiente · PROVISIONAL» sin el bloqueo; el conteo de la l. 84 sí lo cita. | Agregar `(${s.bloqueo…})` con WF3-ajustes-a, como en la l. 84. |
| R3-01 | baja | Convenciones (§1 y §8): dependencias de pruebas declaradas y apoyo compartido en `pruebas/apoyo/` | `frontend/pruebas/ajustes-a/validar-c1.ts:7`; `frontend/package.json` | frontend-base (front-ajustes-a, front-ajustes-b) | `validar-c1.ts` importa `ajv`, que no está en `package.json` (llega como dependencia de otra dependencia). ajustes-b lo importa desde `../ajustes-a/` y también `tablero.ts`. Si cambia el árbol de npm, las pruebas de las dos áreas dejan de cargar. | Declarar `"ajv": "8.20.0"` en `devDependencies` y mover `validar-c1.ts` y `tablero.ts` a `pruebas/apoyo/` (solicitud 8 de la corrección r2). |
| R3-02 | baja | DESIGN §7 Botón (nombre visible), robustez | `frontend/src/vistas-bd/VistaTabla.tsx:11`; `frontend/src/vistas-bd/CeldaEditable.tsx:39` | frontend-base | Hoy hay 0 casos porque las áreas pintan «Sin valor» propiedad por propiedad, pero `Celda` sigue devolviendo `''` por omisión: cualquier propiedad editable nueva sin `pintar` vuelve a dar un botón sin nombre. | `textos.vistas.sinValor` cuando el contenido es `''` y `aria-label={`${prop.nombre}: ${texto}`}` en el botón de CeldaEditable. |

## Estado de los hallazgos de la ronda 2

| Id | Sev. | Estado | Evidencia mía |
|---|---|---|---|
| R2-01 (datos que no validan contra C1) | alta | **CERRADO** | Validador propio sobre los cuerpos literales: 0 fallas en 53 operaciones; `POST promotions` da 201 con un Promotion, `DELETE telegram-chat-candidates` da 204 sin cuerpo, `PATCH ai/providers` trae `provider` y `allowed_uses`, `PATCH settings/security`, `budgets` y `credits` validan. |
| R2-02 (flechas de semana de 40×44) | alta | **CERRADO** | Ninguna «Semana…» entre los controles bajo 44 en 412 táctil (claro y oscuro); `shrink-0` en S06Tiendas.tsx:199 y :201; su prueba pasa. |
| H-03 (evento, «Ver como tabla», descripción emergente) | alta | **ABIERTO** (frontend-base) | 40×44 en /ajustes/05, 20 px con puntero fino y `rgb(255,255,255)` en oscuro, medidos por mí. Archivos sin cambios. |
| H-07 (`Ctrl+/` en 02, 04, 05, 07 y tablero de 04) | media | **REABIERTO EN PARTE**: tablero de 04 cerrado; `Ctrl+/` abierto | Prueba del tablero de 04 pasa y la revisé (vista previa, «Aplicar» enfocado, Esc sin PATCH, Enter con PATCH). 0 manejadores de `Ctrl+/` en `frontend/src`. |
| H-08 (`Ctrl+/` en 09, 13, 14 y tablero de 14) | media | **REABIERTO EN PARTE**: tablero de 14 cerrado; `Ctrl+/` abierto | Prueba del tablero de 14 pasa (estado antes → después, Esc, Enter con `tracking_mode`). Sin `Ctrl+/`. |
| H-09 resto (aviso del alcance de 02) | media | **CERRADO** | 0 casos de `dineroSinTabulares` en todo el recorrido; `tabular-nums` en S02Presupuesto.tsx:99. |
| H-10 (FV-01, FV-04, FV-05) | media | **ABIERTO** (frontend-base) | 128, 52 y 16 mediciones; `outline-color` en la transición del Botón. |
| R2-03 (celdas editables sin nombre) | media | **CERRADO en las pantallas** (queda R3-02, bajo, en la base) | 0 controles sin nombre en todas las páginas medidas; prueba «R2-03» pasa en 12, 13 y 14. |
| R2-04 (P-20 sin valores de Ajustes, N-225) | media | **CERRADO** | P20Comparador.tsx:37-45, 98 y 104; la prueba de P-20 comprueba 500.00, 3, «De Ajustes», «Solo para esta comparación», «25 %», «Sin tope» y el cuerpo enviado. |
| R2-05, R2-06, R2-07 | baja | **ABIERTOS** | Ver la tabla de hallazgos. |

## PREGUNTA para el orquestador o Axel

Esta es la ronda 3 de 3. Lo que bloquea la aprobación (H-03, H-10 y la parte base de `Ctrl+/`) está en archivos de frontend-base, que ninguno de los cuatro agentes revisados puede tocar. ¿Se encarga a una corrección de frontend-base (solicitudes 1 a 7 de `correccion-frontend-areas-r2.md`, más R3-01 y R3-02) y después una revisión puntual, o se acepta con estos pendientes? No lo decido yo.

## Lo que no corrí
- `make verify-wf3-frontend`: escribe `frontend/dist` y usa el 4173. Corrí lo mismo aislado (tsc, build, rutas, hosts, Playwright y variantes). `make verify-all` sí, con salida 0.
