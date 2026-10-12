# REVISIÓN frontend-core · ronda 2 de 3 · WF-3a · 12-oct-2026

Revisor: `revisor-frontend-core-r2` (independiente; no escribió nada de lo revisado y no corrigió nada).
Alcance: trabajo de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras`, más la corrección `correccion-frontend-core-r1`.
Se evaluó **solo** contra U5-06..U5-10 y U5-13 (parte estática y con mocks), C3-03, DESIGN.md, U1 §2.1, §2.2 y §3, y U3.
Todas las cifras de abajo salen de comandos propios. Los temporales quedaron en `datos/tmp/rev-fcore-r2/` y se borraron al terminar.

## Veredicto

**NO APROBADO**: 1 alto, 1 medio y 1 bajo nuevos.

Los 7 hallazgos de la ronda 1 que tocaba comprobar (H-01..H-06 y H-10) quedan **cerrados**.

Los dos hallazgos nuevos que bloquean son del mismo tipo que H-03 y H-02. La prueba `pruebas/medidas.spec.ts` no los detecta por dos razones:
- solo mide las pantallas en reposo, sin abrir la hoja F-1;
- solo comprueba que el icono mida 20 o 24 px, no qué tamaño le toca según DESIGN §5.2.

## Tabla criterio → estado → evidencia propia

| Criterio | Estado | Evidencia (comando propio y salida resumida) |
|---|---|---|
| C3-03 (build) | Cumple | `npm run build`: exit 0, **0 `error TS`** (tsc dentro del build). Precache de 80 entradas. |
| C3-03 (rutas) | Cumple | `python3 -I datos/tmp/rev-fcore-r2/rutas.py` (propio, contra `docs/contratos/extractos/*.yaml`): 170 paths del contrato, 162 rutas literales en `src`, **fuera = 0**. `make verify-wf3-frontend` → «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». |
| U5-08 | Cumple | Recorrido propio de Playwright por 30 vistas en móvil claro, 30 en escritorio oscuro y la hoja F-1, registrando `page.on('request')`: **hosts = `localhost:4173`** (1 host). En `dist`, los `https?://` son solo literales que no se piden (w3.org ×24, react.dev ×2, tailwindcss.com, github.com, bit.ly, sil.org). El script del agente da «hosts externos que se pedirían = 0». |
| U5-06 · objetivos táctiles | **No cumple** | Medición propia a 412 px táctil sobre 697 controles en 30 vistas, con un selector más amplio que el del agente (incluye `[tabindex]`, `[role=link]`, `[role=combobox]`, `[role=gridcell][tabindex]`). En reposo, **0 controles menores de 44 px**. Con la hoja F-1 abierta, el botón **«Más opciones» mide 93×20** y «Menos opciones» 111×20 → **N-01**. Las etiquetas `label[for]` de 20 px de alto no se cuentan: no son controles, el campo asociado mide ≥ 44. |
| U5-06 · iconos | Cumple en tamaño y trazo; **no cumple** la regla de uso de §5.2 | Medido con el `stroke-width` computado: 353 iconos Lucide en móvil y 562 en escritorio oscuro, **0 fuera de 20/24 px o de trazo 1.5 px dibujados**. Hay 0 `svg` que no sean Lucide fuera de gráficas. Hay 0 `size-3`/`size-4` en íconos y 0 `strokeWidth` sueltos. En cambio, hay **14 controles de solo icono con icono de 20 px**, cuando §5.2 pide 24 → **N-02**. |
| U5-06 · cifras tabulares | Cumple | Recorrido de los nodos de texto con `$\d` y de los `input` con `inputmode=decimal` o con valor de dinero, en 30 vistas × 2 combinaciones más F-1: **0 sin `tabular-nums`**. |
| U5-06 · colores y medidas | Cumple (con bajo pendiente) | grep en el alcance: **0 hex**, 0 `rgb()`/`hsl()`/`oklch()`, 0 `style={{…}}` con medidas. Queda 1 valor arbitrario: `shell/Layout.tsx:47 max-w-[85vw]` (H-08, bajo, de la ronda 1, fuera de esta corrección). |
| U5-06 · ≤ 6 tarjetas · U3 | Cumple | `/inicio?sub=resumen`, en móvil y escritorio: **6 tarjetas exactas** (`tarjeta-1`..`tarjeta-6`). |
| U5-06 · componentes de DESIGN §7 | Cumple (con bajo pendiente) | Los componentes de §7 existen en `componentes/` y `vistas-bd/` (Boton, Campo y CeldaEditable, Chip, VistaBD con tabla, tablero, calendario y galería, TarjetaMetrica, Grafica, Hoja, Estados, Avisos, BarraLateral, Barras, MenuComandos y el calendario de pagos). Las áreas los importan: TarjetaMetrica 9, Grafica 8, VistaBD 13, Chip 9, Campo 13, Hoja 10 y Vacio 13. Inicio sigue sin TarjetaMetrica (H-09, bajo, de la ronda 1). |
| U5-07 | Cumple | En `innerText`, `placeholder`, `aria-label` y `title` de 30 vistas × 2 combinaciones más F-1, con 37 términos de interfaz en inglés (la lista de U5-07 más Close, Edit, Next, Back, Filter, Sort, Retry, Select, Clear, Undo, Page, Rows y otros): **0**. Fechas ISO visibles: **0**. Sin red, la franja dice «Sin conexión con el servidor de casa» y el contador «2 en cola». |
| U5-09 | Cumple | e2e del grupo: «U5-09 toques… (sin conexión): 3», «(P-02, movil): 3», «(P-02, escritorio): 3», «desde Inicio (412 px): 3», «(1440 px): 3». La prueba cuenta solo `tap()`/`click()` sobre controles (los `fill` y `keyboard.type` no cuentan) y comprueba ≤ 3. El gasto aparece en la tabla de P-02, y en Inicio se vuelve a pedir `month-vs-plan`. |
| U5-10 | Cumple | Prueba propia: `colorScheme: dark` → `body` `rgb(34, 35, 37)` = `#222325` = `color.neutral.950` = `ui.surface.page.dark`, en las 30 vistas. `light` → `rgb(250, 250, 251)` = `#FAFAFB`. Interruptor: «Claro» → `#FAFAFB`, «Oscuro» → `#222325`, «Según el sistema» → sigue al sistema. |
| U5-13 pasos 1–3 y 6–7 (mocks) | Cumple | Variante propia, distinta de la del agente: 2 gastos capturados sin red. Tras recargar, «2 en cola». Al volver la red con 500: franja «No se pudo enviar · 2 en cola» y aviso «No se pudo enviar · Revisar». El panel muestra «No se pudo enviar: Ocurrió un error en el servidor de casa.». «Reintentar envío» con el servidor aún caído **vuelve a avisar**. Con el servidor recuperado: **2 registros, 0 duplicados** (4 POST con 2 claves únicas de 36 caracteres), `captured_at` con desplazamiento `-06:00` del momento de captura, y la franja desaparece. |
| U1 §2.1 · §2.2 · §3 | Cumple | Están en el código los 23 textos de vacío y error de P-00..P-16 del alcance: «No se pudo calcular», «Sin gastos en…», «Ningún gasto coincide · Quitar filtros», «Este gasto ya no existe», «No hay tickets por confirmar», «Leyendo el ticket…», «No se pudo leer», «Sin presupuesto definido», «Faltan datos para proyectar», «Sin ingreso esperado», «Aún no hay metas», «Semana sin armar», «Sin recetas que coincidan», «sin fuente», «Nada que comprar…», «Sin plan esta semana», «La canasta no tiene artículos», «Sin precios para comparar», «Sin productos», «Inventario vacío…», «Usuario o contraseña incorrectos», «Demasiados intentos» y «No se guardó». Hay esqueletos en 8 archivos. Las e2e de cada pantalla pasan. |
| e2e del grupo | Cumple | `npx playwright test pruebas/inicio pruebas/gastos-presupuesto pruebas/comidas-compras pruebas/cola.spec.ts pruebas/shell.spec.ts pruebas/login.spec.ts pruebas/medidas.spec.ts` → **63 passed** (1.8 min). Pruebas propias (`rev.spec.ts`, `rev2.spec.ts`): 6 passed. |
| H-10 · `make verify-wf3-frontend` | Cumple | Ver «Resultado de make» al final. |

## Hallazgos nuevos

### N-01 · alta · U5-06 / DESIGN §4.6 (`size.touch.min` 44 × 44)
- **Ruta:** `frontend/src/shell/HojaRegistro.tsx:120`.
- **Dueño:** `frontend-base`.
- **Qué pasa:** el botón «Más opciones» / «Menos opciones» de la hoja F-1 (registrar gasto, el flujo principal) es un `<button>` sin altura mínima. A 412 px táctil mide **93×20** («Menos opciones»: 111×20).
- **Por qué no lo detectó la prueba:** `medidas.spec.ts` nunca abre la hoja, así que no lo mide.
- **Corrección sugerida:**
  - Usar `<Boton variante="enlace">`, o añadir al botón `inline-flex items-center coarse:min-h-11`.
  - Ampliar `pruebas/medidas.spec.ts` para que mida también los estados abiertos: hoja F-1 con «Más opciones», panel de cola, menú «/» y hoja de vista previa.

### N-02 · media · DESIGN §5.2 (24 px en barras y botones de solo icono) / U5-06
14 controles de solo icono usan icono de 20 px. §5.2 asigna `size.icon.md` 24 px a «solos en barras: barra superior… botones de solo icono». La corrección de la ronda 1 dice que estos pasaron a `size-6`, pero quedan:

| Dueño | Ruta | Control |
|---|---|---|
| `frontend-base` | `componentes/Hoja.tsx:19` | «Cerrar» de toda hoja |
| `frontend-base` | `shell/Barras.tsx:16` | «Abrir menú» de la barra superior móvil |
| `front-inicio` | `pantallas/inicio/Calendario.tsx:79,81` | |
| `front-gastos-presupuesto` | `pantallas/gastos-presupuesto/comun.tsx:37,39` | «Mes anterior» y «Mes siguiente» de P-02 y P-06 |
| `front-comidas-compras` | `pantallas/comidas-compras/comun.tsx:57,59,145` | |
| `front-comidas-compras` | `P10Plan.tsx:267,269` | |
| `front-comidas-compras` | `P12Receta.tsx:80,82,190` | |

- **Medido:** en móvil, «Abrir menú», «Cerrar», «Mes anterior» y «Mes siguiente» miden 20 px.
- **Corrección sugerida:** `className="size-6" size={ICONO_MD}` en esos 14 iconos. Sumar a `iconosFueraDeNorma` la regla «botón o enlace sin texto visible → icono de 24 px».
- **Lo que no entra:** las flechas «Desplegar» y «Plegar» de las secciones de la barra lateral quedan en 20 px, porque §5.2 pone la barra lateral en `size.icon.sm`.

### N-03 · baja · DESIGN §4.6
- **Ruta:** `frontend/src/shell/Barras.tsx:26`.
- **Dueño:** `frontend-base`.
- **Qué pasa:** el botón «Buscar» de la barra superior (visible desde `tablet`) es `size-8` (32×32) y no tiene `coarse:size-11`.
  - En una tableta táctil queda por debajo de 44.
  - En escritorio queda por debajo de `size.control.fine` (36).
- **Por qué es baja:** no cae en las 4 combinaciones de U5 (412 táctil y 1440 fino).
- **Corrección sugerida:** `size-9 coarse:size-11`, igual que `Boton tam="icono"`.

## Estado de los hallazgos de la ronda 1

| Hallazgo | Estado | Comprobación propia |
|---|---|---|
| H-01 (alta) U5-13 paso 7: 5xx o 401 sin aviso | **Cerrado** | Ver la fila U5-13 de la tabla. `cola.ts` deja la captura `pendiente` con `error`. `useCola().retenidas` mantiene visible la franja y `App.tsx` avisa. Las e2e de 500 y 401 del agente pasan. |
| H-02 (alta) iconos con trazo 2 y de 16 px | **Cerrado** | `ProveedorIconos` (`LucideProvider` de lucide-react 1.55.0 con `absoluteStrokeWidth`) envuelve la app. 0 iconos fuera de 20/24 px o de trazo 1.5 en 915 medidos. La regla de uso residual queda aparte en N-02 (media). |
| H-03 (alta) celdas de 36–40 px y encabezado «Día» de 23 px | **Cerrado** | 0 controles menores de 44 en las 30 vistas en reposo, incluidas P-02, P-04, P-08, P-10, P-15 y P-16. |
| H-04 (alta) P-06 «Despensa» y «Transporte», P-04 «Todos los gastos» | **Cerrado** | Igual que H-03. El enlace de regreso de `Barras.tsx` lleva `min-h-11`. |
| H-05 (alta) casilla de P-13 y enlace de P-12 | **Cerrado** | La casilla vive en una etiqueta de 44×44 que la contiene (DESIGN §4.6: «aunque el dibujo sea menor»), y la prueba toca la esquina. El enlace de P-12 mide ≥ 44. |
| H-06 (media) dinero sin `tabular-nums` | **Cerrado** | 0 importes sin tabulares, incluidos los `input`. |
| H-10 (media) script `verify` | **Cerrado** | `"verify": "npm run verificar"` en `package.json`. `make verify-wf3-frontend` corre build, C3-03, U5-08, manifiesto y e2e. |
| H-07, H-08, H-09, H-11 (bajas) | Siguen abiertos (fuera del encargo de la corrección) | H-08: `Layout.tsx:47 max-w-[85vw]`. H-09: Inicio no importa `TarjetaMetrica` (`Resumen.tsx:17-20`). No bloquean. |

## Resultado de make

- `make verify-wf3-frontend` (exit 0) corre `npm run verify`, que da:
  - build y tsc limpios;
  - «C3-03 frontend: 202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)»;
  - «U5-08 dist: hosts externos que se pedirían = 0» y «hosts pedidos: localhost; externos: 0»;
  - e2e de todas las áreas: **133 passed** (3.1 min).
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: **exit 2** en `verify-go` por «FALLA gofmt: internal/reportes/reportes_test.go».
  - **Está fuera del alcance de esta revisión**: es backend de S7 y ningún agente del frontend lo toca.
  - Se deja constancia para el orquestador, porque la verificación global de referencia está rota.

## Qué falta para aprobar en la ronda 3

1. **N-01**: altura táctil del botón «Más opciones» de F-1, y que la prueba mida también los estados abiertos.
2. **N-02**: iconos de 24 px en los 14 controles de solo icono, y que la prueba compruebe la regla de uso.

N-03 y los bajos de la ronda 1 son opcionales.
