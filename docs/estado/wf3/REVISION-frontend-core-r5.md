# REVISIÓN frontend-core · ronda 5 · WF-3a · 12-oct-2026

Revisor: `revisor-frontend-core-r5`. Es independiente: no escribió nada de lo revisado ni corrigió nada.

**Qué se revisó:** el trabajo de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras` en el árbol de trabajo actual. Es HEAD `a35451e` más las correcciones sin commit, entre ellas `correccion-frontend-core-r4`.

**Contra qué (solo esto):** U5-06 a U5-10 y U5-13 (la parte estática y con mocks), C3-03, DESIGN.md, U1 §2.1/§2.2/§3 y U3.

**Cómo:** todas las cifras salen de comandos propios, en `datos/tmp/rev-fcore-r5/` (ya borrado):
- `npm run build -- --outDir ../datos/tmp/rev-fcore-r5/dist`: un build completo (preparar, `tsc --noEmit` y vite) en mi propia carpeta, servido con `vite preview` en el puerto 4895. Otro revisor corría al mismo tiempo en el 4799.
- Una configuración de Playwright propia que corre la suite del grupo sobre mi `dist`.
- 3 pruebas propias:
  - `revision.spec.ts`: barrido de 36 vistas en 412 táctil claro y en 1440 oscuro, F-1, la cola U5-13, P-03, el tema y los estados de U1.
  - `revision2.spec.ts`: las Tendencias de P-06, P-07 y P-10, el orden de N-053 y los estados vacíos.
  - `tactil.spec.ts`: 20 vistas, con objetivos táctiles, iconos, foco y nombre accesible de las gráficas.
- Un barrido de fechas con formato híbrido.
- `rutas.py`, un script propio contra `docs/contratos/extractos/*.yaml` y `openapi.yaml`.
- `grep` estático de colores y medidas.
- `make verify-all`.

## Veredicto

**NO APROBADO: 0 altos y 1 medio. Además hay 4 bajos nuevos.**

- **R4-01** y **R4-02**, los dos medios de la ronda 4, quedan **cerrados**. Los medí yo mismo (ver «Estado de los hallazgos anteriores»).
- **R5-01 (medio, `front-comidas-compras`).** En P-10 Tendencias, el eje de semanas de las dos gráficas y su «Ver como tabla» muestran las fechas como «Lun 5 oct-2026». Ese formato híbrido no es ninguno de los que pide U5-07 (`10-oct-2026` / `Vie 30 oct 2026`). Además, se arma a mano en lugar de usar `fechaConDia()`, que es lo que pide CONVENCIONES-frontend §5.
- **Bajos nuevos:**
  - R5-02: «Guardar» no hace nada y no avisa cuando la forma de pago es «A crédito».
  - R5-03: N-053 también ofrece los préstamos a plazo y las tarjetas (PREGUNTA).
  - R5-04: la leyenda de las gráficas de Inicio usa `text-secondary`.
  - R5-05: una prueba e2e es intermitente.

## Tabla criterio → estado → evidencia propia

| Criterio | Estado | Evidencia (comando propio y salida resumida) |
|---|---|---|
| C3-03 · build y tsc | Cumple | `npm run build -- --outDir ../datos/tmp/rev-fcore-r5/dist --emptyOutDir` → exit 0, sin errores ni avisos de tipos, «precache 81 entries (2489.27 KiB)». Los archivos generados (`esquema.gen.ts`, tokens, iconos, mocks) no cambiaron: el md5 es igual y `git status` también. |
| C3-03 · rutas | Cumple | `rutas.py`: 172 rutas en los extractos, las mismas 172 que en `openapi.yaml`. Hay 208 llamadas método+ruta literales en `src` (sin contar `esquema.gen.ts` ni `mocks/`). **Fuera del contrato: 0. Métodos no declarados: 0.** Solo hay dos `fetch(`: `api/cliente.ts:43` y `api/conexion.ts:23` (`/health`). Cruce con el servidor (`node scripts/rutas-api.mjs`) → «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». P-03 llama ahora a `GET /api/v1/credits` y a `POST /api/v1/incomes`; las dos están en el contrato. |
| U5-08 | Cumple | El barrido de 36 vistas × 2 combinaciones, con F-1, la cola y P-03, registró `page.on('request')`. **Hosts: `localhost:4895`**, y ningún otro. `node scripts/hosts-externos.mjs` sobre mi `dist` → «hosts externos que se pedirían = 0». En `dist`, los `https?://` son solo literales: w3.org ×24, react.dev ×2, tailwindcss.com, github.com, bit.ly y scripts.sil.org. |
| U5-06 · objetivos táctiles | Cumple | 412 táctil, 20 vistas: **498 controles**. 0 controles < 44 px sin una zona de 44 × 44 que reciba el toque. En la hoja P-03 con N-053 y N-102 abiertos: 0 controles < 44. La suite (`medidas.spec.ts`) da 0 menores de 44 px en Inicio (83), Gastos y presupuesto (346), Comidas y compras (246) y estados abiertos (299). |
| U5-06 · iconos | Cumple | 222 iconos Lucide en 20 vistas. **0 fuera de 20/24 px o de 1.5 px de trazo dibujado.** Los de 20 px llevan `stroke-width` 1.8 con `absoluteStrokeWidth`: 1.8 × 20/24 = 1.5 px. |
| U5-06 · foco visible | Cumple | Tab en P-02 → «Abrir menú» con `outline: solid 2px rgb(117,121,194)`. |
| U5-06 · cifras tabulares | Cumple | Barrido de 36 vistas × 2, con HTML y `<text>` del SVG de ECharts: **0 montos sin `tabular-nums`**. En P-06 y P-10 Tendencias, 0 `<text>` sin tabular (`revision2.spec.ts`). |
| U5-06 · etiquetas de dinero en texto y eje compacto (§7.7, §8.5) | Cumple | Texto del SVG medido en cada gráfica: <br>• **P-02:** eje `$0 · $1 mil … $4 mil`; etiquetas `$4,000.00 · $1,500.00 · $0.00 · $1,200.00`. <br>• **P-06:** `$0 · $5 mil … $25 mil`; `$22,000.00 · $22,000.00 · $21,000.00 · $18,500.00`. <br>• **P-07:** `$0 · $3 mil … $18 mil`; `$18,000.00 · $17,000.00 · $1,000.00`. <br>• **P-08:** `$0 · $5 mil … $30 mil`; `$30,000.00 · $2,500.00`. <br>• **P-09:** `$0 · $10 mil … $50 mil`; `$45,000.00 · $1,500.00 · $3,000.00`. <br>• **P-10, costo:** `$0 · $200 … $1 mil`; `$980.00 · $68.76`. <br>• **P-01 Tendencias:** `$0 · $1 mil … $4 mil`; «Efectivo acumulado $2,921.25» y «Devengado acumulado $3,895.00». <br>El compacto solo sale en los ejes. Dinero con formato fuera de `$1,234.56` (o del compacto en ejes): **0**. |
| U5-06 · colores y medidas | Cumple (H-08 bajo sigue abierto) | `grep` en `app api cola componentes estilos shell textos vistas-bd` y en las 3 áreas, sin contar `tokens.css` ni `tailwind.theme.css`: **0 hex**, 0 `rgb()`/`hsl()`/`oklch()` y 0 `dark:`. Hay 3 `style={{…}}`, todos dinámicos: 2 de dnd-kit y el ancho en % de la barra de avance. Valores arbitrarios: `Hoja.tsx:16 max-h-[85dvh]`, zona segura en `Hoja.tsx:24` y `Layout.tsx:48 max-w-[85vw]` (H-08). El de `TarjetaMetrica.stories.tsx:52` es solo de Storybook. |
| U5-06 · ≤ 6 tarjetas · U3 | Cumple | `/inicio` a 412 claro y a 1440 oscuro: **6 tarjetas exactas** (`[data-prueba^="tarjeta-"]` dentro de `resumen`). |
| U5-06 · componentes de DESIGN §7 | Cumple | Archivos que importan cada componente (inicio / gastos-presupuesto / comidas-compras / shell): <br>• Boton 5/9/8/2 · Campo 0/6/7/2 · Chip 5/0/4/1 · TarjetaMetrica 0/5/4/0 · Grafica 2/5/2/0 <br>• VistaBD 0/6/7/0 · Hoja 1/4/6/2 · ConDatos 0/8/7/0 · Vacio 0/7/6/1 · Franja 0/3/6/2 · useAvisos 0/8/7/1 <br>`BarreraErrores` se usa en `Layout`. |
| U5-07 · texto visible | **No cumple (fechas)** | Lista propia de unos 170 términos en inglés (con días, meses, `undefined`, `NaN`, `Invalid Date`, `AM`/`PM` y frases de ECharts), aplicada a `innerText` y a `aria-label`, `title`, `alt`, `placeholder` y `aria-description`, fuera de `aria-hidden`: **0 en inglés** en 36 vistas × 2. **Fechas ISO visibles: 0.** Barrido de formatos híbridos con todas las `<details>` abiertas y el texto del SVG: **1 vista, P-10 Tendencias, con «Lun 5 oct-2026» y «Lun 12 oct-2026»** → R5-01. Sin red: «Sin conexión con el servidor de casa» y «2 en cola». |
| U5-07 · nombre accesible de las gráficas | Cumple (R4-03 bajo sigue abierto) | Las 4 `Grafica` medidas (P-02, P-06 y dos en P-10) tienen `role="img"`, `tabindex="0"` y `aria-label` en español. En `GraficaInicio` sigue la descripción automática en inglés, pero dentro de `aria-hidden`: la suite la imprime como «PENDIENTE U5-07» (R4-03). |
| U5-09 | Cumple | Prueba propia en línea desde P-02, a 412 táctil: «Registrar gasto», «Transporte» y «Guardar» = **3 toques**. 1 POST con `amount_cents` 9740 y una clave de 36 caracteres; «$97.40» aparece en la tabla P-02. Sin red: 3 toques. La suite da 3 sin conexión, 3 en P-02 (móvil y escritorio) y 3 desde Inicio. |
| U5-10 | Cumple | En las 36 vistas, el `body` es `rgb(250,250,251)` en claro (`#FAFAFB`, neutral 50) y `rgb(34,35,37)` en oscuro (`#222325`, neutral 950 = `ui.surface.page` oscuro). Interruptor a 1440, con 3 clics: <br>• sistema claro → claro, sin `data-theme`; sistema oscuro → oscuro, sin `data-theme`. <br>• clic 1 → claro (`light`); clic 2 → oscuro (`dark`); clic 3 → «Según el sistema»: sin `data-theme`, y sigue al sistema cuando este pasa a claro. <br>Es el token, no una inversión de colores. |
| U5-13 pasos 1–3 y 6–7 (mocks) | Cumple | Variante propia a 412 táctil: un gasto de $97.40 en «Transporte» y un peso de 68.4 kg, sin red. <br>• **Paso 1:** «Sin conexión con el servidor de casa» y «2 en cola». <br>• **Paso 2:** tras recargar sin red sigue «2 en cola», con 0 y 0 registros. <br>• **Paso 7:** la red vuelve con 500 en gastos. Aparece el aviso «No se pudo enviar» y la franja sigue en «2 en cola», con 0 gastos y 0 pesos. Con «Revisar» → «Reintentar» y el servidor sano: **1 gasto y 1 peso**, y el contador desaparece. <br>• **Pasos 3 y 6:** **2 POST de gasto con 1 sola clave** de 36 caracteres (el reintento no duplica). `captured_at` = `2026-10-12T07:52:07-06:00`, con desplazamiento, tomado 148 ms después de capturar y más de 2.5 s antes del envío. El peso tiene su propia clave de 36. |
| U1 §2.2 · estados | Cumple | **Carga** en P-01 (respuestas con 2.5 s de retraso): 6 bloques `aria-busy` y 0 «$0.00» provisionales. <br>**Error** (500 en todos los GET): P-01 muestra 6 «Reintentar», con el mensaje de U3 «No se pudo cargar el presupuesto del mes. Revisa que estés en la red de casa.». P-02, P-04, P-06, P-07, P-08, P-09, P-10, P-11, P-13, P-14, P-15 y P-16 muestran «Ocurrió un error en el servidor de casa.» con «Reintentar». 0 caídas de la barrera. <br>**Vacío:** «Sin gastos en oct-2026» con «Registrar gasto» (P-02), «Sin recetas que coincidan» (P-11) y «Sin productos» (P-15). <br>**Sin conexión:** ver U5-13. 0 errores de JavaScript en el barrido. |
| U1 §2.1 · §3.5 (P-03) | Cumple (R5-02 y R5-03 bajos) | Pestaña «Ingreso»: campos Monto, Concepto, Fecha y Nota; «Ir a Ingresos» = 0. `POST /api/v1/incomes` `{kind:'extra', description:'Venta de bici', amount_cents:120050, income_date:'2026-10-10'}` con una clave de 36 caracteres. «A crédito»: lista «Elige la línea · Ábaco tarjeta · Tarjeta bancaria · Línea de la tienda · Crédito de tienda · Zeta préstamo · Préstamo a plazo», en orden alfabético es-MX, sin la línea inactiva. «Ir a compra a crédito» = 0. Al elegir una línea va a `/creditos/{id}/compras/nueva?monto=185.50&fecha=2026-10-10`. |
| e2e del grupo | Cumple, con R5-05 bajo | `npx playwright test inicio/ gastos-presupuesto/ comidas-compras/ cola.spec.ts shell.spec.ts login.spec.ts medidas.spec.ts base.escritorio.spec.ts` con mi configuración y mi `dist` → **111 passed, 3 skipped, 1 failed** (2.9 min). <br>• Los 3 saltados son los de `Ctrl+/` en móvil (atajo de teclado, correcto). Ya no se salta ninguna de P-03. <br>• El fallo fue `gastos-presupuesto/correccion-bis.escritorio.spec.ts:79` (P-02 `Ctrl+/`, escritorio): «Expected 2, Received 0» parches tras Enter. Al reejecutarlo pasó 3/3 solo y 57/57 con `gastos-presupuesto/ --project=escritorio --repeat-each=3`. Es intermitente → R5-05. <br>• Mis pruebas propias: 8 + 1 + 1 + 1 passed (sin contar las reejecuciones), 0 fallos. |
| `make verify-all` | Cumple | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → exit 0: «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |

## Hallazgos nuevos

### R5-01 · media · U5-07 (fechas en formato) y CONVENCIONES-frontend §5

- **Ruta:** `frontend/src/pantallas/comidas-compras/P10Plan.tsx:245`.
- **Dueño:** `front-comidas-compras`.

**Qué pasa (medido a 412 táctil, con los datos falsos del área):**
- En P-10 Tendencias, «Costo semanal de comidas» y «kcal promedio por día contra la meta» rotulan las semanas como «Lun 5 oct-2026» y «Lun 12 oct-2026».
- Esas mismas etiquetas son los encabezados de fila de «Ver como tabla».
- El código las arma así: `diaCorto(w.week_start) + ' ' + textoMes(w.week_start.slice(0, 7))`.
- U5-07 solo admite `10-oct-2026`, o `Vie 30 oct 2026` con el día de la semana. CONVENCIONES §5 da `fechaConDia(iso)` → `Vie 30 oct 2026` para este caso.
- Es la única vista, de 32 barridas, con un formato híbrido.

**Corrección sugerida:**
- Usar `fechaConDia(w.week_start)` (→ «Lun 5 oct 2026») o `fecha(w.week_start)` (→ «5-oct-2026») de `src/api/formatos.ts` y quitar la concatenación.
- Agregar una aserción a `pruebas/comidas-compras/` sobre el texto de los ejes o de la tabla.

### R5-02 · baja · U1 §3.5 / DESIGN §7.2 (acción sin respuesta)

- **Ruta:** `frontend/src/shell/HojaRegistro.tsx:89` (`if (forma === 'credit') return`) y `:141`.
- **Dueño:** `frontend-base`.

**Qué pasa:** con «A crédito» elegido, «Guardar» sigue habilitado. Al tocarlo no pasa nada: la hoja sigue abierta y no hay mensaje (medido). La compra se termina con «Continuar en Compra a crédito». La propia `correccion-frontend-core-r4` lo dejó como PROPUESTA.

**Corrección sugerida (es decisión de diseño, va a Axel):** deshabilitar «Guardar» con el motivo visible («La compra a crédito se termina en Créditos»), con `deshabilitadoPor`, o esconderlo mientras esa forma esté elegida.

### R5-03 · baja / PREGUNTA · U1 N-053

- **Ruta:** `frontend/src/pantallas/gastos-presupuesto/hoja-registro.tsx:36` y `:54`.
- **Dueño:** `front-gastos-presupuesto`.

**Qué pasa:**
- N-053 dice: «Líneas de crédito activas (nombre y tipo)… Elegir una revolvente o de tienda lleva a P-19 prellenada».
- La pieza lista todas las activas, también `term_loan` (préstamo a plazo) y `bank_card`, y cualquiera lleva a P-19. Medido: «Zeta préstamo · Préstamo a plazo» → `/creditos/3/compras/nueva?monto=185.50…`.
- U1 no dice qué pasa con un préstamo a plazo ni con una tarjeta. No hay PROVISIONAL ni bloqueo.

**Corrección sugerida:** abrir un bloqueo (`WF3a-…-n053-tipos`) con estas opciones:
- (a) listar solo `revolving_line` y `store_credit`;
- (b) listar todas, pero que solo esas dos lleven a P-19;
- (c) dejarlo como está, si P-19 acepta los 4 tipos.

Mientras tanto, marcarlo PROVISIONAL.

### R5-04 · baja · DESIGN §7.7 (leyenda en `ui.text.default`)

- **Ruta:** `frontend/src/pantallas/inicio/GraficaInicio.tsx:76` (`legend.textStyle.color = --ui-text-secondary`).
- **Dueño:** `front-inicio`.

**Qué pasa:** DESIGN §7.7 dice «leyenda (con 2 o más series) con marca de forma y nombre en `ui.text.default`». `Grafica.tsx:75` sí lo cumple; las 3 gráficas de P-01 Tendencias usan el texto secundario.

**Corrección sugerida:** `textStyle: { color: token('--ui-text-default') }` en `ejes()`.

### R5-05 · baja · prueba intermitente

- **Ruta:** `frontend/pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts:79-97`.
- **Dueño:** `front-gastos-presupuesto`.

**Qué pasa:** en la corrida completa del grupo falló una vez: «Expected: 2 / Received: 0» en `expect.poll(() => parches.length)`, después de `Enter`. Pasó en las 6 reejecuciones siguientes. Al parecer, el `Enter` llega antes de que el bloque de edición en lote tenga el foco o el estado listos.

**Corrección sugerida:** antes de `Enter`, esperar a que el bloque `edicion-varias` tenga el foco (`toBeFocused` en su control) o pulsar el botón «Aplicar», y revisar si `EdicionVarias.tsx` escucha `Enter` en el contenedor correcto.

## Estado de los hallazgos anteriores

| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| **R4-01** (media) P-03: N-102 pestaña Ingreso y N-053 líneas | **Cerrado** | `HojaRegistro.tsx:42` monta `FormIngreso` y `:137` monta `LineaCredito`. Medido: <br>• la pestaña Ingreso tiene Monto, Concepto, Fecha y Nota, y registra `POST /api/v1/incomes` con clave; <br>• «A crédito» lista las líneas activas con su tipo, en orden alfabético, y sigue a P-19 prellenada; <br>• «Ir a Ingresos» = 0 e «Ir a compra a crédito» = 0. <br>Las pruebas `correccion-bis.escritorio.spec.ts:166` y `:179` ya no tienen `test.skip` y **pasan** en móvil y escritorio. Sin conexión, las dos piezas muestran «Requiere conexión» (código `hoja-registro.tsx:33` y `:73`). Ver R5-02 y R5-03, bajos nuevos. |
| **R4-02** (media) Gráficas: eje compacto, tabular en SVG y etiquetas de dinero | **Cerrado** | `Grafica.tsx:63` y `:89` ponen el compacto (`dineroCompacto`) solo en los ejes de dinero. `:121` y `GraficaInicio.tsx:64` usan `tabular-nums`, que heredan los `<text>`. `:95`–`:97` ponen las etiquetas selectivas con el formato completo y `hideOverlap`. Medido en P-01, P-02, P-06, P-07, P-08, P-09 y P-10: 0 montos del SVG sin tabular, y el compacto solo en los ejes (ver la tabla). P-15 usa la misma `Grafica` de dinero (`P15Productos.tsx:135`), pero con mis datos falsos no pinta gráfica: la revisé por el código. |
| R4-03 (baja) `aria.label` en inglés en `GraficaInicio` | Sigue abierto | La suite imprime «PENDIENTE U5-07 … This is a chart…», dentro de `aria-hidden`. |
| R4-04 (baja) alto de la barra inferior | Sigue abierto | No se encargó en la ronda 4. |
| R4-05 (baja) excepción `PENDIENTES_INGLES` | Sigue abierto | Va ligado a R4-03. |
| H-07, H-08, H-09 (bajas) | Siguen abiertos | H-08: `Layout.tsx:48 max-w-[85vw]`, medido con `grep`. |

## Preguntas para Axel

- **PREGUNTA (R5-03).** En la hoja «Registrar gasto», con la forma de pago «A crédito»: ¿se ofrecen también los préstamos a plazo y las tarjetas bancarias, o solo las líneas revolventes y de tienda, como dice N-053?
- **PREGUNTA (R5-02).** Con «A crédito» elegido, ¿«Guardar» se deshabilita con su motivo o se esconde?
