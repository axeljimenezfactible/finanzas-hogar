# REVISIÓN frontend-core · ronda 4 · WF-3a · 12-oct-2026

Revisor: `revisor-frontend-core-r4`. Independiente: no escribió nada de lo revisado y no corrigió nada.

**Qué se revisó:** el trabajo de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras` en el árbol de trabajo actual (HEAD `a35451e` más las correcciones sin commit de `frontend-base-fix`, `front-gastos-presupuesto-fix` y `front-comidas-compras-fix`).

**Contra qué (solo esto):** U5-06 a U5-10 y U5-13 (parte estática y con mocks), C3-03, DESIGN.md, U1 §2.1/§2.2/§3 y U3.

**Cómo:** todas las cifras salen de comandos propios en `datos/tmp/rev-fcore-r4/` (ya borrado):
- build propio en `datos/tmp/rev-fcore-r4/dist` y servidor `vite preview` propio en el puerto 4894 (otro revisor corría a la vez sobre `frontend/dist`; así no se pisan);
- `rutas.py` propio contra `docs/contratos/extractos/*.yaml`;
- dos pruebas Playwright propias (`revision.spec.ts`: barrido de 29 vistas × 2 combinaciones, estados de U1, F-1, cola U5-13 y tema; `extra.spec.ts`: casillas, texto SVG de gráficas y P-03), que reutilizan solo los datos falsos de cada área;
- lista propia de 150 términos en inglés (incluidas las frases de ECharts, días, meses, `undefined`, `NaN`, `Invalid Date`, `AM`/`PM`), aplicada a `innerText` y a `aria-label`, `aria-description`, `title`, `alt` y `placeholder`;
- la medición de cifras tabulares incluye ahora los nodos `<text>` del SVG de ECharts (las rondas anteriores solo medían HTML).

## Veredicto

**NO APROBADO: 0 altos, 2 medios y 3 bajos nuevos.**

- Los dos medios de la ronda 3 (**N-04** y **N-05**) quedan **cerrados**, con medición propia. N-06 y N-03 (bajos) también quedan cerrados.
- **R4-01 (media, frontend-base):** en P-03, la pestaña «Ingreso» no registra el ingreso en la hoja y la forma de pago «A crédito» no lista las líneas (U1 §3.5: N-102 y N-053). Las piezas ya existen en el área; falta montarlas en `HojaRegistro.tsx`. Las pruebas que lo comprueban se saltan solas.
- **R4-02 (media, frontend-base):** las gráficas de `Grafica` no cumplen DESIGN §7.7/§8.5 en el dinero: eje con formato completo (`$25,000.00`) en vez del compacto, cifras del SVG sin `tabular-nums` (también en `GraficaInicio`) y montos sin etiqueta directa.
- Bajos nuevos: R4-03 (descripción de ECharts en inglés en `GraficaInicio`, oculta a lectores de pantalla), R4-04 (alto de los destinos de la barra inferior) y R4-05 (la excepción `PENDIENTES_INGLES` de la prueba).

## Tabla criterio → estado → evidencia propia

| Criterio | Estado | Evidencia (comando propio y salida resumida) |
|---|---|---|
| C3-03 (build y tsc) | Cumple | `npx tsc --noEmit -p .` → exit 0, sin salida (8 s). `npx vite build --outDir ../datos/tmp/rev-fcore-r4/dist` → exit 0, 0 errores, «precache 81 entries». |
| C3-03 (rutas) | Cumple | `rutas.py` propio: 170 rutas en los extractos; 165 literales `/api/v1/…` en `src` (sin `esquema.gen.ts` ni `mocks/`); **fuera del contrato = 0**; 364 llamadas método+ruta con **0 métodos no declarados**. Únicos `fetch(`: `api/cliente.ts:43` (puerta única) y `api/conexion.ts:23` (`/health`, está en `openapi.yaml:35`). Cruce con el servidor: `node scripts/rutas-api.mjs` → «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». |
| U5-08 | Cumple | Barrido propio (29 vistas en 412 táctil claro y en 1440 oscuro, más estados, F-1 y cola) con `page.on('request')`: **hosts = `localhost:4894`** únicamente. En `dist`, los `https?://` son solo literales: w3.org ×24, react.dev ×2, tailwindcss.com, github.com, bit.ly y scripts.sil.org. `hosts-externos.mjs` sobre mi `dist` → «hosts externos que se pedirían = 0». |
| U5-06 · objetivos táctiles | Cumple | 412 táctil, selector amplio (`button`, `a[href]`, `input`, `select`, `textarea`, `summary`, roles `button/link/tab/checkbox/switch/combobox/menuitem/option/radio` y `[tabindex]`): 702 controles en 29 vistas. Los únicos < 44 son 9 casillas de 20 × 20 (selección de filas en P-02 y P-15, «Marcar … como comprado» en P-13), y cada una vive en un `label` de **44 × 44** que recibe el toque (comprobado con `elementFromPoint` a ±20 px del centro). 0 enlaces chicos dentro de texto. |
| U5-06 · iconos | Cumple | 347 iconos Lucide en móvil y 548 en escritorio oscuro: **0 fuera de 20/24 px o de trazo 1.5 px dibujados**. Barra inferior: 5 iconos de **24 px** (ver N-05). |
| U5-06 · color y rol de la barra inferior (G-05) | Cumple | En `/inicio`, `/comidas` y `/compras` a 412 claro: fondo `rgb(250,250,251)` = `--ui-surface-raised`; píldora del destino activo **56 × 32** `rgb(240,241,255)` = `--ui-surface-selected`; icono activo `rgb(93,96,159)` = `--ui-icon-selected`; nombre activo `rgb(54,54,58)` peso 500 = `--ui-text-default`; reposo `rgb(102,103,110)` = `--ui-text-secondary`. Ver R4-04 (bajo) por el alto. |
| U5-06 · cifras tabulares | **No cumple** | HTML: 0 montos sin `tabular-nums`. **SVG de ECharts: 26 montos por combinación sin `tabular-nums`**: ejes de P-02 y P-06 Tendencias (`$0.00` … `$25,000.00`, `font-variant-numeric: normal`, 12 px) y de Inicio Tendencias (`$1 mil` … y las etiquetas «Efectivo acumulado $2,921.25», «Devengado acumulado $3,895.00») → **R4-02**. |
| U5-06 · etiquetas de dinero en texto y eje compacto | **No cumple** | `Grafica.tsx:66` usa `formato` completo en el eje (`$25,000.00`); DESIGN §7.7 y §8.5 piden el compacto solo en ejes (`$5 mil`). Las series de `Grafica` no tienen `label`: ningún monto en texto sobre la marca (§8.5 «Ningún monto vive solo en la descripción emergente»). `GraficaInicio` sí lo hace (`compacto`, `selectivos`) → **R4-02**. |
| U5-06 · colores y medidas | Cumple (bajo H-08 pendiente) | grep en `app api cola componentes estilos shell textos vistas-bd` y las 3 áreas: **0 hex**, 0 `rgb()`/`hsl()`/`oklch()` fuera de `tokens.css`, 0 usos de `dark:`. `style={{…}}`: 3, todos dinámicos (2 de dnd-kit y el ancho % de la barra de avance). Valores arbitrarios: `Layout.tsx:48 max-w-[85vw]` (H-08), `Hoja.tsx:16 max-h-[85dvh]`, zona segura en `Hoja.tsx:24` y `Barras.tsx:52`, y `transition-[color,background-color,border-color]` (lista de propiedades, no medida). |
| U5-06 · ≤ 6 tarjetas · U3 | Cumple | `/inicio` en 412 claro y 1440 oscuro: **6 tarjetas exactas**: «Mes contra plan», «Flujo de la quincena · $13,802.36», «Próximo pago · $973.75», «Comidas de la semana», «Mejor canasta · Bodega Aurrera $3,094.50» y «Ahorro y metas». |
| U5-06 · componentes de DESIGN §7 | Cumple (bajo H-09 pendiente) | Archivos que importan cada componente (inicio / gastos-presupuesto / comidas-compras): Boton 5/9/8, Campo 0/6/7, Chip 5/0/4, TarjetaMetrica 0/5/4, Grafica 1/5/2, VistaBD 0/6/7, Hoja 1/4/6, ConDatos 0/8/7, Vacio 0/7/6, Franja 0/3/6, useAvisos 0/8/7. `BarreraErrores` envuelve cada pantalla en `Layout`. |
| U5-07 · texto visible | Cumple | Lista propia de 150 términos sobre `innerText` y atributos expuestos de 29 vistas × 2 combinaciones más estados: **0 en inglés**. Fechas ISO visibles (`AAAA-MM`, `AAAA-MM-DD`, con hora): **0**. Sin red: «Sin conexión con el servidor de casa» y «2 en cola». |
| U5-07 · nombre accesible de las gráficas | Cumple (bajo R4-03) | `Grafica` en P-02 y P-06 Tendencias, 412 claro y 1440 oscuro: `role="img"`, `tabindex="0"`, `aria-label` = «Gasto del hogar por mes. Parte del hogar devengada por mes, las 6 categorías con más gasto y «Otras».» y «Plan contra real por mes. Plan total y real total por mes calendario.»; 0 frases de ECharts. En Inicio Tendencias, las 2 gráficas de `GraficaInicio` siguen con «This is a chart… the data for oct-2026 is 0, 389500…», pero en un `div aria-hidden="true"` sin `tabindex`: ni se pinta ni llega al lector de pantalla → R4-03. |
| U5-09 | Cumple | Prueba propia en línea desde P-02 (412 táctil): «Registrar gasto», «Transporte», «Guardar» = **3 toques** (teclear 97.40 no cuenta); 1 POST con `amount_cents` 9740, y «$97.40» visible en la tabla. Sin red (cola): 3 toques. e2e del grupo: «U5-09 toques… (sin conexión): 3», «(P-02, movil): 3», «desde Inicio (412 px): 3». |
| U5-10 | Cumple | En las 29 vistas, `colorScheme: dark` da `body` `rgb(34,35,37)` = `--ui-surface-page` oscuro (`#222325`), y claro `rgb(250,250,251)`. Interruptor (1440): sistema claro → `rgb(250,250,251)`; sistema oscuro → `rgb(34,35,37)` sin `data-theme`; clic 1 «Cambiar tema: Claro» → claro, `data-theme=light`; clic 2 «Oscuro» → `rgb(34,35,37)`, `dark`; clic 3 «Según el sistema» → sigue al sistema (oscuro y luego claro al cambiarlo). Es el token, no una inversión. |
| U5-13 pasos 1–3 y 6–7 (mocks) | Cumple | Variante propia: gasto $97.40 «Transporte» y peso 68.4 kg sin red. (1) «Sin conexión con el servidor de casa» y «2 en cola»; (2) tras recargar sin red, «2 en cola» y 0 registros; (7) red de vuelta con 500 en gastos: aviso «No se pudo enviar · Revisar», franja «No se pudo enviar · 2 en cola», 0 gastos y 0 pesos (la cola respeta el orden); «Revisar» → «Reintentar envío» con el servidor sano: **1 gasto y 1 peso**, contador desaparece; (3) **2 POST de gasto con 1 sola clave** de 36 caracteres; peso con su clave de 36; `captured_at` `2026-10-12T07:12:39-06:00` dentro de la ventana de captura y anterior al envío; cuerpo `amount_cents` 9740 y `weight_g` 68400; (6) reenvío con la misma clave y cuerpo → 200, **id 1**, sigue 1 gasto. |
| U1 §2.2 · estados | Cumple | **Carga** P-01 (respuestas con 2.5 s de retraso): 6 bloques `aria-busy`, sin `$0.00` provisional. **Error** P-01 (500): 6 tarjetas, cada una con su mensaje de U3 («No se pudo cargar el presupuesto del mes. Revisa que estés en la red de casa.») y **6 «Reintentar»**. **Error** 500 en P-02, P-04, P-06, P-08, P-10, P-11, P-13, P-14 y P-15: «Ocurrió un error en el servidor de casa.» con Reintentar, inglés 0 y **0 caídas** de la barrera. **Vacío**: «Sin gastos en oct-2026» (P-02), «Sin recetas que coincidan» (P-11) y «Sin productos» (P-15), con acción. Barrido: 0 «No se pudo mostrar esta página». |
| U1 §2.1 · §3 (P-03) | **No cumple** | Prueba propia: en la hoja, pestaña «Ingreso» → solo «El ingreso se captura en Presupuesto › Ingresos. · Ir a Ingresos», sin campo Monto. `HojaRegistro.tsx` no llama a `createIncome` ni a `listCredits` (grep = 0) → **R4-01**. |
| e2e del grupo | Cumple | `npx playwright test inicio/ gastos-presupuesto/ comidas-compras/ cola.spec.ts shell.spec.ts login.spec.ts medidas.spec.ts base.escritorio.spec.ts` (config propia sobre mi `dist`) → **108 passed, 7 skipped, 0 failed** (2.8 min). De los 7 saltados, 3 son `Ctrl+/` en móvil (atajo de teclado, correcto) y **4 son P-03 N-102 y N-053** (2 por proyecto), que se saltan solos porque la hoja muestra «Ir a Ingresos» / «Ir a compra a crédito» (R4-01). Pruebas propias: `revision.spec.ts` 10 passed. |
| `make verify-all` y `make verify-wf3-frontend` | Cumple | Ver «Resultado de make». |

## Hallazgos nuevos

### R4-01 · media · U1 §3.5 y §2.1 (P-03: pestaña Ingreso y forma de pago a crédito)

- **Ruta:** `frontend/src/shell/HojaRegistro.tsx:39` (pestaña Ingreso) y `:134` (forma de pago `credit`).
- **Dueño:** `frontend-base` (las piezas son de `front-gastos-presupuesto` y ya existen).

**Qué pasa:**
- Pestaña «Ingreso»: la hoja muestra «El ingreso se captura en Presupuesto › Ingresos.» y un botón «Ir a Ingresos». No hay campo Monto ni se registra nada (medido).
- Forma de pago «A crédito»: aviso «La compra a crédito se registra en Créditos con sus pagos.» y enlace «Ir a compra a crédito» a P-17; no se listan las líneas activas (sin `GET /api/v1/credits`).
- U1 §3.5 dice: «Las pestañas Ingreso y Peso de la misma hoja usan N-102 y N-240»; N-053: «Líneas de crédito activas (nombre y tipo), en orden alfabético. Elegir una revolvente o de tienda lleva a P-19 prellenada».
- `pantallas/gastos-presupuesto/hoja-registro.tsx` ya exporta `LineaCredito` (`:29`) y `FormIngreso` (`:66`). Es la «solicitud 1» de `front-gastos-presupuesto-fix.md:49`, que llegó después de `frontend-base-fix` y no se aplicó.
- Las pruebas `correccion-bis.escritorio.spec.ts:166` y `:177` se saltan con `test.skip(…isVisible())`: la verificación pasa en verde sin cubrirlo.

**Corrección sugerida:**
- En `HojaRegistro.tsx`, montar `<FormIngreso refMonto=… alTerminar=…/>` en la pestaña Ingreso (en lugar de `IrA`) y `<LineaCredito monto dia comercio alSalir/>` cuando `forma === 'credit'` (en lugar de la `Franja` con «Ir a compra a crédito»).
- Después, cambiar los `test.skip` condicionales de esas dos pruebas por aserciones, para que no vuelvan a saltarse.

### R4-02 · media · DESIGN §7.7 y §8.5 / U5-06 (cifras tabulares, etiquetas de dinero en texto)

- **Ruta:** `frontend/src/componentes/Grafica.tsx:66` (eje Y) y `:67-79` (series sin `label`); contenedor `:92`. Para las cifras tabulares, también `frontend/src/pantallas/inicio/GraficaInicio.tsx:73`.
- **Dueño:** `frontend-base` (y `front-inicio` para su contenedor, salvo que se resuelva con una regla común en `app.css`).
- **Afecta a:** P-02, P-06, P-07, P-08, P-09, P-10 y P-15 (Grafica) y P-01 Tendencias (cifras tabulares).

**Qué pasa (medido en P-02 y P-06 Tendencias, 412 y 1440):**

| Elemento | Medido | Lo que pide DESIGN |
|---|---|---|
| Eje Y de dinero | `$0.00`, `$1,000.00` … `$25,000.00` (formato completo con centavos) | §7.7: «dinero del eje compacto: `$0`, `$5 mil`, `$10 mil`»; §8.5: «el compacto (`$5 mil`) solo en ejes» |
| Cifras del SVG de ECharts | `<text style="font-size: 12px; font-family: Inter…">`, `font-variant-numeric: normal`: 5 en P-02 Tendencias, 6 en P-06 Tendencias y 15 en Inicio Tendencias, en las dos combinaciones | §7.7: ejes con `tabular-nums`; U5-06: «cifras tabulares en todo dinero» |
| Etiquetas de dinero | ninguna sobre las marcas; los montos solo en la descripción emergente, el eje y «Ver como tabla» | §8.5: «Ningún monto vive solo en la descripción emergente»: líneas con el último valor y los extremos; columnas con el valor si cabe (si no, en la tabla); con ≤ 4 series, etiqueta directa |

**Corrección sugerida:**
- Pasar a `api/formatos.ts` un `dineroCompacto` (el `compacto` de `GraficaInicio.tsx`, aritmética entera) y que `Grafica` lo use en `yAxis.axisLabel.formatter` cuando el formato sea dinero (prop opcional `formatoEje`; kilos y kcal quedan igual).
- Agregar `tabular-nums` al contenedor del gráfico (el `div ref` de `Grafica` y de `GraficaInicio`; el `<text>` lo hereda) o una regla común en `app.css`: `[_echarts_instance_] { font-variant-numeric: tabular-nums; }`.
- Etiquetado selectivo en las series de `Grafica` con `formato` completo (reutilizar `selectivos` de `GraficaInicio`): líneas con el último valor y los extremos; columnas con el valor solo con ≤ 4 series.
- Sumar a `dineroSinTabulares` (`pruebas/apoyo/medidas.ts:77`) los nodos `<text>` del SVG.

### R4-03 · baja · U5-07 (consistencia) / DESIGN §7.7

- **Ruta:** `frontend/src/pantallas/inicio/GraficaInicio.tsx:49`. **Dueño:** `front-inicio`.
- `aria: { enabled: true, decal: { show: false } }` deja que ECharts escriba su descripción automática en inglés con centavos crudos («This is a chart with type Line chart named Deuda viva total… the data for oct-2026 is 0, 389500…») en el `aria-label` del contenedor.
- Es baja porque el contenedor es `aria-hidden` y sin `tabindex`: no se pinta ni lo lee TalkBack. Pero es texto en inglés en el DOM y obliga a la excepción de R4-05.
- **Corrección:** `aria: { enabled: true, label: { enabled: false }, decal: { show: false } }` (el nombre accesible ya lo dan `figcaption` y la tabla `sr-only`), o `label.description` en español como en `Grafica`.

### R4-04 · baja · DESIGN §7.10 (alto de la barra inferior)

- **Ruta:** `frontend/src/shell/Barras.tsx:52`. **Dueño:** `frontend-base`.
- El `nav` mide `h-14` (56 px) con el borde superior y `pb-[env(safe-area-inset-bottom)]` dentro: cada destino mide **55 px** (medido 82 × 55) y, en un Android con zona segura, aún menos. DESIGN §7.10: «Alto `layout.bottombar.height` con la zona segura de Android… cada destino mide 56 px de alto».
- **Corrección:** quitar `h-14` del `nav` y dar `h-14` a cada destino, para que el borde y la zona segura se sumen al alto.

### R4-05 · baja · pruebas (U5-07)

- **Ruta:** `frontend/pruebas/medidas.spec.ts:36` (`PENDIENTES_INGLES`). **Dueño:** `frontend-base`.
- La prueba imprime como «PENDIENTE» y no falla con la descripción en inglés de Inicio Tendencias. Su comentario dice que la lista debe quedar vacía al corregirse, pero nada lo exige hoy.
- **Corrección:** al cerrar R4-03, vaciar la lista y quitar la excepción.

## Estado de los hallazgos anteriores

| Hallazgo | Estado | Comprobación propia |
|---|---|---|
| **N-04** (media) descripción de ECharts en inglés en `Grafica` | **Cerrado** | `Grafica.tsx:48` usa `aria.label.description` = «título. resumen»; el contenedor (`:92`) lleva `role="img"` y `tabindex="0"`. Medido en P-02 y P-06 Tendencias, 412 claro y 1440 oscuro: `aria-label` en español, sin centavos crudos, 0 frases de ECharts. El resto en `GraficaInicio` va como R4-03 (bajo). |
| **N-05** (media) barra inferior | **Cerrado** | Fondo = `--ui-surface-raised`, 5 iconos de 24 px, píldora 56 × 32 en `--ui-surface-selected` con icono `--ui-icon-selected`, nombre en `--ui-text-default` 500, reposo `--ui-text-secondary`; medido en los 3 destinos activos. `medidas.spec.ts` ya mide G-05. El alto va como R4-04 (bajo). |
| N-06 (baja) sin barrera de errores | **Cerrado** | `Layout` envuelve cada pantalla con `BarreraErrores`. Barrido y 9 pantallas con 500: 0 caídas y el cascarón sigue vivo. |
| N-03 (baja) «Buscar» y «Desplegar» `size-8` | **Cerrado** | `Barras.tsx:16` y `:26`: `size-9 coarse:size-11`. |
| H-07 (baja) API falsa con página para GET que no son listas | Sigue abierto | `api-falsa.ts:41` sin cambios. Ya no tira la aplicación gracias a la barrera. |
| H-08 (baja) `max-w-[85vw]` | Sigue abierto | `Layout.tsx:48`. |
| H-09 (baja) Inicio sin `TarjetaMetrica` ni `ConDatos` | Sigue abierto | 0 importaciones en `pantallas/inicio`. |

## Resultado de make

- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: exit 0** («Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»). Ya no falla `go test ./...` como en la ronda 3.
- **`make verify-wf3-frontend`: exit 0.** Corre `npm run verify`, que da:
  - build y tsc limpios;
  - «C3-03 frontend: 208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)»;
  - «U5-08 dist: hosts externos que se pedirían = 0»;
  - e2e de todas las áreas: **223 passed, 8 skipped** (4.9 min). Entre los saltados están las 4 de P-03 (R4-01).
- `git status` antes y después de los dos `make`: solo cambian los informes nuevos. Los archivos generados no cambiaron.

## Qué falta para aprobar

1. **R4-01:** montar `FormIngreso` y `LineaCredito` en `HojaRegistro.tsx` y activar sus pruebas.
2. **R4-02:** eje de dinero compacto, `tabular-nums` en el texto de las gráficas y etiquetas de dinero selectivas en `Grafica`.

R4-03, R4-04, R4-05, H-07, H-08 y H-09 son bajos y no bloquean.
