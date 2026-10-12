# REVISIÓN frontend-core · ronda 3 de 3 · WF-3a · 12-oct-2026

Revisor: `revisor-frontend-core-r3`. Es independiente: no escribió nada de lo revisado y no corrigió nada.

**Qué se revisó:**
- El trabajo de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras`.
- La corrección `correccion-frontend-core-r2`.

**Contra qué se evaluó (solo esto):**
- U5-06 a U5-10 y U5-13 (la parte estática y con mocks).
- C3-03.
- DESIGN.md.
- U1 §2.1, §2.2 y §3.
- U3.

**Cómo se obtuvieron las cifras:**
- Todas salen de comandos propios, no de los reportes de los agentes.
- Las pruebas propias son de Playwright y están en `datos/tmp/rev-fcore-r3/` (`medir.ts`, `barrido.spec.ts`, `estados.spec.ts` y `tema.spec.ts`), con su propia configuración. Reutilizan solo los datos falsos de cada área.
- Las mediciones son propias, con un selector de controles más amplio que el del agente.
- Los temporales se borraron al terminar.

## Veredicto

**NO APROBADO: 0 altos, 2 medios y 1 bajo nuevos.**

**Hallazgos de la ronda 2:**
- N-01 (alto) y N-02 (medio) quedan **cerrados**. Se comprobó con medición propia.
- N-03 (bajo) sigue abierto. No estaba en el encargo.

**Los dos medios nuevos (N-04 y N-05):**
- Son de `frontend-base` y estaban desde la ronda 1.
- Las pruebas no los detectan por dos razones:
  - la lista de palabras en inglés de las rondas anteriores no incluía las frases que genera ECharts;
  - la regla de iconos de `medidas.spec.ts` no mira la barra inferior.
- Su corrección es local y no toca cifras.

## Tabla criterio → estado → evidencia propia

| Criterio | Estado | Evidencia (comando propio y salida resumida) |
|---|---|---|
| C3-03 (build y tsc) | Cumple | `npm run build`: exit 0, **0 `error TS`**, precache de 80 entradas. `npx tsc --noEmit -p .`: exit 0, sin salida. |
| C3-03 (rutas) | Cumple | Script propio `rutas.py` contra `docs/contratos/extractos/*.yaml` (172 paths): 163 rutas literales en `src`, **fuera del contrato = 0**. Los únicos `fetch` son los de `api/cliente.ts:43` (la puerta única) y `api/conexion.ts:23` (`/health`, que está en el contrato). `npm run verificar:rutas` → «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». |
| U5-08 | Cumple | Barrido propio de 31 vistas en 412 táctil claro y en 1440 fino oscuro, más los estados abiertos, registrando `page.on('request')`: **hosts = `localhost:4173`**. En `dist`, los `https?://` son solo literales que no se piden: w3.org ×24, react.dev ×2, tailwindcss.com, github.com, bit.ly y scripts.sil.org. `npm run verificar:hosts` → «hosts externos que se pedirían = 0». |
| U5-06 · objetivos táctiles | Cumple | A 412 px táctil, con un selector amplio (incluye `[role=link]`, `[role=combobox]`, `[tabindex]` y `[role=gridcell][tabindex]`): **0 controles menores de 44 × 44**. Se midieron 729 controles en 31 vistas y estos estados abiertos: F-1 Gasto; F-1 con «Más opciones», todas las categorías y «Crédito»; pestañas Ingreso, Peso y Ticket; menú «/»; panel lateral móvil; franja «2 en cola»; y panel de la cola con falla. «Más opciones» mide **123×44** y «Menos opciones» **141×44**. |
| U5-06 · iconos (tamaño y trazo) | Cumple | 381 iconos Lucide en móvil y 596 en escritorio oscuro, más los de los estados abiertos: **0 fuera de 20/24 px o de trazo 1.5 px dibujados**. Hay 0 `size-3`/`size-4` y 0 `strokeWidth` sueltos en el alcance. |
| U5-06 · iconos (regla de uso de §5.2) | **No cumple** | Controles de solo icono: **0 con icono de 20 px** (N-02 cerrado). En cambio, los 4 destinos con icono de la barra inferior (Inicio, Comidas, Compras y Buscar) miden **20 px**, y DESIGN §5.2 y §7.10 piden `size.icon.md` 24 → **N-05**. |
| U5-06 · color y rol de la barra inferior | **No cumple** | Fondo medido en claro: `rgb(242, 242, 244)` = `--ui-surface-card` `#f2f2f4`. DESIGN §7.10 pide `ui.surface.raised` (`#fafafb`). El destino activo no tiene la píldora de 56 × 32 en `ui.surface.selected` ni el icono en `ui.icon.selected` → **N-05**. |
| U5-06 · cifras tabulares | Cumple | Nodos de texto con `$\d` en 31 vistas × 2 combinaciones más los estados abiertos (incluido el panel de la cola, «Gasto $42.35»): **0 sin `tabular-nums`**. |
| U5-06 · colores y medidas | Cumple (con bajos pendientes) | grep en el alcance: **0 hex**, 0 `rgb()`/`hsl()`/`oklch()` y 0 `dark:`. Los `style={{…}}` son 3, todos dinámicos: 2 transformaciones de dnd-kit y el ancho en % de la barra de avance. `app.css` solo tiene el foco de 2 px (`border.width.focus` y `focus.offset`). Valores arbitrarios: `Layout.tsx:47 max-w-[85vw]` (H-08, bajo; DESIGN §4.3 fija 85 %). Además `Hoja.tsx:16 max-h-[85dvh]`, `Hoja.tsx:24 pb-[max(env(safe-area-inset-bottom),12px)]` y `Barras.tsx:48 pb-[env(safe-area-inset-bottom)]`: zona segura que DESIGN §4.3 pide, y 12 px = `space.3`. |
| U5-06 · ≤ 6 tarjetas · U3 | Cumple | `/inicio?sub=resumen` en móvil claro y en escritorio oscuro: **6 tarjetas exactas**. Son «Mes contra plan», «Flujo de la quincena», «Próximo pago», «Comidas de la semana», «Mejor canasta» y «Ahorro y metas». Con 500 en todas las lecturas, cada tarjeta muestra su error con «Reintentar» (**6 botones**), con el texto de U3 §2.1/§7 («No se pudo cargar el presupuesto del mes. Revisa que estés en la red de casa.»). |
| U5-06 · componentes de DESIGN §7 | Cumple (con bajo pendiente) | Archivos que importan cada componente, por área (inicio / gastos-presupuesto / comidas-compras): Boton 5/8/8, Campo 0/6/7, Chip 5/0/4, TarjetaMetrica 0/5/4, Grafica 0/5/2, VistaBD 0/6/7, CeldaEditable 0/2/0, Hoja 1/3/6, ConDatos 0/7/7, Vacio 0/7/6, Franja 0/2/6 y useAvisos 0/7/7. Inicio sigue sin `TarjetaMetrica` ni `ConDatos` (H-09, bajo). |
| U5-07 · texto pintado | Cumple | `innerText`, `placeholder`, `title` y `alt` de 31 vistas × 2 combinaciones más los estados abiertos, con una lista propia de 120 términos (incluidos `the`, `of`, `and`, `undefined`, `NaN`, `Invalid Date`, días y meses en inglés, `AM`/`PM`): **0 en el texto pintado**. Fechas ISO visibles (`AAAA-MM-DD` y `AAAA-MM`): **0**. Sin red: «Sin conexión con el servidor de casa» y «2 en cola». |
| U5-07 · nombre accesible de las gráficas | **No cumple** | En el `aria-label` del contenedor enfocable de cada `Grafica`, ECharts escribe su descripción automática en inglés, con centavos sin formato. Ejemplo: «This is a chart. It consists of 2 series count. The 0 series is a Bar chart representing Despensa. The data is as follows: the data for sep-2026 is 0, 400000…». Esa descripción sustituye al `aria-label` en español. Se midió en P-02 y P-06 Tendencias, en los dos proyectos → **N-04**. |
| U5-09 | Cumple | Las e2e del grupo dan 3 toques en todos los recorridos: «(sin conexión)», «(P-02, movil)», «(P-02, escritorio)», «desde Inicio (412 px)» y «(1440 px)». En la prueba propia, F-1 sin red con «Transporte» cuenta **3 toques** con la regla K10: «+», categoría y Guardar; teclear no cuenta. |
| U5-10 | Cumple | En las 31 vistas, `colorScheme: dark` da `body` `rgb(34, 35, 37)` = `#222325` = `color.neutral.950` (`ui.surface.page` oscuro), y `light` da `rgb(250, 250, 251)` = `#FAFAFB`. El interruptor se probó aparte (ver «Interruptor de tema» al final). |
| U5-13 pasos 1–3 y 6–7 (mocks) | Cumple | Variante propia, distinta de la del agente: 1 gasto ($42.35, «Transporte») y 1 peso (81.25 kg) capturados sin red. Tras recargar sin red: «2 en cola» y 0 registros en el servidor. Al volver la red con 500 en gastos: aviso «No se pudo enviar», franja «No se pudo enviar · 2 en cola» y 0 registros; la cola se detiene en la primera captura y respeta el orden de captura (U1 §2.4). Panel de la cola → «Reintentar envío» con el servidor recuperado: **1 gasto y 1 peso**, **2 POST de gasto con 1 sola clave** (36 caracteres) y **0 duplicados**. `captured_at` = `2026-10-12T04:39:12-06:00`, dentro de la ventana de captura y no la del envío. Cuerpo: `amount_cents` 4235 y `weight_g` 81250. La franja y el contador desaparecen al vaciarse la cola. |
| U1 §2.1 · §2.2 · §3 | Cumple (con un bajo nuevo) | **Vacíos** con API vacía: «Sin gastos en oct-2026 · Registrar gasto» (P-02), «Sin recetas que coincidan» (P-11) y «Sin productos» (P-15). **Error** con 500 en lecturas: P-02, P-06, P-08, P-10, P-13, P-14 y P-04 muestran «Ocurrió un error en el servidor de casa.» con **Reintentar** e inglés 0; P-01 lo hace por tarjeta. **Abiertos:** P-03 con sus pestañas y «Más opciones». Ver N-06 (bajo): sin una barrera de errores global, una respuesta con forma inesperada deja la aplicación en blanco. |
| e2e del grupo | Cumple | `npx playwright test pruebas/inicio pruebas/gastos-presupuesto pruebas/comidas-compras pruebas/cola.spec.ts pruebas/shell.spec.ts pruebas/login.spec.ts pruebas/medidas.spec.ts` → **64 passed** (1.9 min). Pruebas propias: barrido 6 passed; estados 3 passed, en dos corridas; tema 1 passed. |
| `make verify-wf3-frontend` | Cumple | Ver «Resultado de make». |

## Hallazgos nuevos

### N-04 · media · U5-07 (inglés y dinero sin formato) / DESIGN §7.7 (accesibilidad de gráficas)

- **Ruta:** `frontend/src/componentes/Grafica.tsx:44`. También afecta a `:76`.
- **Dueño:** `frontend-base`.
- **Afecta a:** toda pantalla con `Grafica`. Medido en P-02 y P-06 Tendencias. También la usan P-07, P-08, P-09, P-10 y P-15.

**Qué pasa:**
- `graf.setOption({ aria: { enabled: true, decal: { show: true } } })` activa la descripción automática de ECharts.
- ECharts escribe esa descripción en el `aria-label` del mismo `div` que la app etiquetó con `${titulo}. ${resumen}`, y la sustituye (`echarts/lib/visual/aria.js:206`).
- El resultado es un nombre accesible en inglés con importes en centavos crudos («the data for sep-2026 is 0, 400000»).
- Ese contenedor es enfocable (`tabIndex={0}`), así que TalkBack lo lee en voz alta.

**Qué incumple:**
- U5-07: texto en inglés y dinero fuera de `$1,234.56`.
- DESIGN §7.7: cada punto con nombre del tipo «Octubre, $8,240.10».

**Por qué es media:** es texto que se lee pero no se pinta. La tabla de severidad pone el texto en inglés como medio.

**Corrección sugerida:**
- Mantener el patrón de forma (decal, K27) y dar la descripción en español:
  `aria: { enabled: true, label: { description: \`${titulo}. ${resumen}\` }, decal: { show: true } }`.
  - Con `label.description`, ECharts usa ese texto tal cual (`aria.js:133-134`).
  - Otra opción es `label: { enabled: false }` y dejar el `aria-label` propio.
- Añadir a la lista de U5-07 de las pruebas frases como `This is a chart`, `series` y `the data`, y revisar `aria-label` y `aria-description` en las vistas Tendencias.

### N-05 · media · DESIGN §7.10 (barra inferior G-05) y §5.2 / U5-06 (iconos y tokens de rol)

- **Ruta:** `frontend/src/shell/Barras.tsx:42-56` (`BarraInferior`).
- **Dueño:** `frontend-base`.

**Qué pasa (medido a 412 px):**

| Elemento | Medido | Lo que pide DESIGN |
|---|---|---|
| Iconos de Inicio, Comidas, Compras y Buscar | `className="size-5"`, **20 px** | §7.10: «con icono `size.icon.md`»; §5.2: 24 px «solos en barras: barra superior, barra inferior» |
| Fondo de la barra | `bg-surface-card`: `rgb(242, 242, 244)` en claro | §7.10: «fondo `ui.surface.raised`», `#FAFAFB` en claro |
| Destino activo | Solo cambia el color y el peso del texto (`aria-[current=page]:text-text-default font-medium`); fondo detrás del icono: ninguno | §7.10, estado Activo: píldora de 56 × 32 px (`size.bottombar.indicator-width` × `-height`) en `ui.surface.selected` detrás del icono `ui.icon.selected` |

En oscuro `card` y `raised` valen lo mismo; la diferencia de fondo se ve en claro.

**Corrección sugerida:**
- Iconos con `className="size-6" size={ICONO_MD}`.
- Fondo `bg-surface-raised`.
- Envolver el icono en un `span` `inline-flex h-8 w-14 items-center justify-center rounded-full` con `bg-surface-selected text-icon-selected` cuando `aria-current="page"`.
  - Antes de usar `h-8 w-14`, confirmar en `tailwind.theme.css` que 32 px y 56 px existen como medidas de C2 (`size.bottombar.indicator-*`).
- Sumar a `iconosFueraDeNorma` la regla «barra inferior → 24 px».

### N-06 · baja · U1 §2.4 (estado de error; «nunca una pantalla en blanco»)

- **Ruta:** `frontend/src/app/App.tsx`, `frontend/src/shell/Layout.tsx`. Hoy solo hay barrera de errores en `pantallas/inicio/comun.tsx:52`.
- **Dueño:** `frontend-base`.

**Qué pasa:**
- No hay una barrera de errores (error boundary) alrededor del contenido de las pantallas.
- Con la API falsa de las pruebas, `/comidas` y `/compras` lanzan `TypeError` («Cannot read properties of undefined (reading 'filter' / 'flatMap')»). Esa API devuelve una página `{items,next_cursor}` para GET que no son listas (H-07).
- React desmonta toda la aplicación: queda `main` = 0 y desaparecen la barra lateral, la barra inferior y la franja.

**Por qué es baja:** el disparador es una respuesta fuera del contrato, que un servidor conforme a C1 no envía.

**Corrección sugerida:**
- Barrera de errores por pantalla en `Layout`, con `ErrorCarga` («No se pudo mostrar esta página · Reintentar»), para que el cascarón y la cola sigan vivos.
- Atender H-07 en `api-falsa.ts`.

## Estado de los hallazgos anteriores

| Hallazgo | Estado | Comprobación propia |
|---|---|---|
| **N-01** (alta) «Más opciones» de F-1 de 93×20 | **Cerrado** | `HojaRegistro.tsx:120` usa `<Boton variante="enlace" className="self-start">`. Medido: **123×44** («Más opciones») y **141×44** («Menos opciones»). F-1 con «Más opciones» abierto, todas las categorías y «Crédito»: 0 controles menores de 44 entre 21. `medidas.spec.ts` ya recorre los estados abiertos y pasa dentro de los 64 del grupo. |
| **N-02** (media) 14 controles de solo icono con icono de 20 px | **Cerrado** | Las 14 rutas citadas tienen `className="size-6" size={ICONO_MD}`: `Hoja.tsx:20`, `Barras.tsx:16`, `Calendario.tsx:80,82`, gastos-presupuesto `comun.tsx:38,40`, comidas-compras `comun.tsx:58,60,146`, `P10Plan.tsx:267,269` y `P12Receta.tsx:81,83,191`. Medición propia en 31 vistas × 2 y estados abiertos: **0 controles de solo icono con icono distinto de 24 px** y trazo de 1.5 px dibujados en todos. La regla quedó en `pruebas/apoyo/medidas.ts`. La barra inferior queda aparte en N-05. |
| N-03 (baja) «Buscar» `size-8` sin `coarse:size-11` | Sigue abierto | `Barras.tsx:26` sin cambios. Pasa lo mismo con «Desplegar» en `Barras.tsx:17` (`size-8`, por debajo de `size.control.fine` 36 en escritorio). No cae en las combinaciones de U5. |
| H-07, H-08, H-09, H-11 (bajas, ronda 1) | Siguen abiertos | H-07: la API falsa devuelve una página para GET que no son listas (dispara N-06). H-08: `Layout.tsx:47 max-w-[85vw]`. H-09: Inicio sin `TarjetaMetrica` ni `ConDatos`. H-11: no se revisó la prueba del agente; mi variante sí comprueba la clave única y `captured_at`. |

## Interruptor de tema (U5-10)

Prueba propia `tema.spec.ts`, a 1440 px, con el interruptor `data-prueba="interruptor-tema"` de la barra lateral:

| Paso | Etiqueta del interruptor | Fondo de `body` | `data-theme` |
|---|---|---|---|
| Al entrar, sistema claro | «Tema: Según el sistema» | `rgb(250, 250, 251)` | — |
| Sistema → oscuro | «Según el sistema» | `rgb(34, 35, 37)` | — |
| Sistema → claro | «Según el sistema» | `rgb(250, 250, 251)` | — |
| Clic 1 | «Claro» | `rgb(250, 250, 251)` | `light` |
| Clic 2 | «Oscuro» | `rgb(34, 35, 37)` | `dark` |
| Clic 3 | «Según el sistema» | `rgb(250, 250, 251)` | — |
| Clic 3 + sistema oscuro | «Según el sistema» | `rgb(34, 35, 37)` | — |

El oscuro es el token `#222325`, no una inversión de color.

## Resultado de make

**`make verify-wf3-frontend`: exit 0.** Corre `npm run verify`, que da:
- build y tsc limpios;
- «C3-03 frontend: 202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)»;
- «U5-08 dist: hosts externos que se pedirían = 0»;
- e2e de todas las áreas: **142 passed** (3.3 min).

**`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: exit 2** en `verify-wf1-cobertura` (`go test ./...`).
- WF-0, C1, S1, U1 y criterios pasan, y gofmt ya no falla como en la ronda 2.
- Fallan 3 pruebas de backend por la misma causa: la operación `getDashboardBestBasket` (`GET /api/v1/dashboard/best-basket`) no tiene handler en ningún área y responde 501.
  - `cmd/finanzas` `TestAreas_SinOperacionesEn501`: «1 de 216 operaciones sin registrar».
  - `internal/api` `TestContrato_Operaciones/getDashboardBestBasket` y `TestContrato_Rutas`: «sin implementar 1».
  - `internal/auth` `TestPermisos_MismasOperacionesAmbasCuentas`.
- **Está fuera del alcance de esta revisión**, porque es backend. El frontend llama a esa ruta, que sí está en el contrato y en `routes.go`, desde `pantallas/inicio/Resumen.tsx:57` (tarjeta «Mejor canasta») y `Detalle.tsx:142`. Contra el servidor real esa tarjeta mostraría su error por tarjeta. El orquestador asigna al dueño de backend.
- Lo dejo constar porque la verificación global de referencia está rota y porque C3-03 («implementado» = ruta registrada **y** atendida) no se cumple para esa operación del lado del servidor.

## Qué falta para aprobar

1. **N-04:** descripción accesible de `Grafica` en español y con dinero formateado.
2. **N-05:** barra inferior según DESIGN §7.10: iconos de 24 px, fondo `ui.surface.raised` y píldora del destino activo.

N-06, N-03 y los bajos de la ronda 1 son opcionales.
