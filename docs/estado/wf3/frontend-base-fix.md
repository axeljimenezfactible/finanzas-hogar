# frontend-base-fix · WF-3a bis (12-oct-2026)

Agente: `frontend-base-fix` (Opus). Corrige los hallazgos de frontend-base abiertos en `REVISION-frontend-core-r3.md` (N-04, N-05, N-06, N-03), `REVISION-frontend-areas-r3.md` (H-03, H-07/H-08 parte base, H-10, R3-02) e `INTEGRACION-wf3a.md` §7 (I-01). No cambia cifras, motores, contrato ni decisiones de Axel. No se hizo `git add` ni `commit`.

## Hecho

| Hallazgo | Qué se hizo | Ruta |
|---|---|---|
| **N-04** (media, U5-07 / §7.7) | La descripción accesible ya no es la automática de ECharts: `aria: { enabled: true, label: { enabled: true, description: `${titulo}. ${resumen}` }, decal }`; el contenedor lleva `role="img"`. | `frontend/src/componentes/Grafica.tsx` |
| **N-05** (media, §7.10 / §5.2) | Barra inferior: fondo `bg-surface-raised`; iconos `size-6` (24 px, `ICONO_MD`); píldora de 56 × 32 (`h-8 w-14`, `size.bottombar.indicator-*` de C2) con `bg-surface-selected` y `text-icon-selected` en el destino con `aria-current="page"`; anillo de foco interior. «Compras» sale de `textos`. Regla «barra inferior → 24 px» sumada a `iconosFueraDeNorma`. | `frontend/src/shell/Barras.tsx`, `frontend/pruebas/apoyo/medidas.ts` |
| **H-03 (1)** (alta) | Calendario: por debajo de 768 px es **agenda por día** (DESIGN §7.5.4: «360: agenda por día…; 412: agenda…; 768 y 1024: mes») con encabezado fijo `Vie 30 oct 2026` y renglones de ancho completo y 44 px; desde 768, la cuadrícula del mes (renglón de 24 px mínimo con puntero fino, 44 con grueso). | `frontend/src/vistas-bd/VistaCalendario.tsx` |
| **H-03 (2)** | «Ver como tabla»: `inline-flex min-h-9 … coarse:min-h-11` (36 px fino = `size.control.fine`, 44 px grueso) con flecha `chevron-right` que gira al abrir. | `Grafica.tsx` |
| **H-03 (3)** | Descripción emergente con tokens releídos al cambiar el tema: fondo `--ui-surface-raised`, borde `--ui-border-divider`, texto `--ui-text-default`, `radius.md` y `shadow.1`, cifras tabulares; guía vertical de 1 px en `--ui-viz-baseline`; leyenda en `--ui-text-default` y apagada en `--ui-text-disabled`. | `Grafica.tsx` |
| **H-07 / H-08** (media, parte base) | Selección de filas en la vista tabla (casilla por fila «Seleccionar {fila}» y «Seleccionar todas» con estado mixto; área de toque 36 fino / 44 grueso), barra «n seleccionadas · Editar seleccionadas (Ctrl+/) · Quitar selección», y **`Ctrl+/`** con el foco en la vista abre el bloque «Editar n registros a la vez»: propiedad → valor → **vista previa sumada «Antes → Después»** (la da la pantalla) → `Enter` o «Aplicar a n» aplica a todas; `Esc` o «Cancelar» cierra sin cambiar ninguna y regresa el foco. Sin vista previa calculada no aplica («No se pudo calcular el efecto · Reintentar» si falla); sin conexión, «Aplicar» deshabilitado con motivo. API tipada `DefVista.editarVarias: { propiedades, vistaPrevia(filas, propiedad, valor), guardar(filas, propiedad, valor) }`, ejemplo y prueba en CONVENCIONES §7. El bloque «Antes → Después» quedó en un solo componente (`BloqueCambios`) que también usa `CeldaEditable`. | `frontend/src/vistas-bd/tipos.ts`, `VistaBD.tsx`, `VistaTabla.tsx`, `EdicionVarias.tsx` (nuevo), `CeldaEditable.tsx` |
| **H-10** (media, §7.0 / §4.5) | Botón: `transition-[color,background-color,border-color]` (el anillo de foco ya no tiene transición). Gráfica: `emphasis: { disabled: true }` por serie. `Campo` y `ChipSeleccionable` deshabilitados: `disabled:` relleno `surface-disabled`, texto `text-disabled`, borde punteado `border-control`, `cursor-not-allowed` y hover solo con `enabled:`. | `Boton.tsx`, `Grafica.tsx`, `Campo.tsx`, `Chip.tsx` |
| **N-06** (baja) | `BarreraErrores` alrededor de cada pantalla en `Layout` («No se pudo mostrar esta página · Reintentar»); se reinicia al cambiar de pantalla, subsección o consulta sin volver a montar la pantalla. | `frontend/src/componentes/BarreraErrores.tsx` (nuevo), `frontend/src/shell/Layout.tsx` |
| **N-03** (baja) | «Mostrar barra lateral» y «Buscar» de la barra superior: `size-9 coarse:size-11` (antes `size-8`). | `Barras.tsx` |
| **I-01** (baja) | `cola.spec.ts`: antes de que el servidor «se recupere» espera a que termine el intento automático tras recargar (POST contado y candado `fh-cola-envio` libre) y toca «Revisar» con reintento acotado (`toPass`). | `frontend/pruebas/cola.spec.ts` |
| R3-02 (baja, de paso) | `Celda` pinta «Sin valor» cuando el contenido es vacío: ninguna celda editable queda como botón sin nombre. | `VistaTabla.tsx` |
| Pruebas nuevas | Banco de pruebas de la base (Vite de desarrollo en 127.0.0.1, nunca entra a `dist/`) y `base.escritorio.spec.ts` (7 pruebas × 2 proyectos); medición G-05 en `medidas.spec.ts`; U5-07 (`textosEnIngles`: texto, `aria-label`, `title`, `alt`, incluida la descripción de ECharts) en todas las vistas de `medidas.spec.ts`; evento de `/ajustes/05` a 412 táctil en la prueba U5-06 de ajustes-a. | `frontend/pruebas/apoyo/banco.ts`, `frontend/pruebas/banco/`, `frontend/pruebas/base.escritorio.spec.ts`, `medidas.spec.ts`, `ajustes-a/ajustes.escritorio.spec.ts` |
| CONVENCIONES | §1 (pruebas de la base), §6 (Grafica: el `resumen` en español con dinero formateado; deshabilitado; BarreraErrores), §7 (calendario y edición en grupo con API, ejemplo con `planned-items` de C1 y prueba de referencia), §8 (mediciones y banco). | `docs/estado/wf3/CONVENCIONES-frontend.md` |

## Archivos

Modificados: `frontend/src/componentes/{Grafica,Boton,Campo,Chip}.tsx`, `frontend/src/shell/{Barras,Layout}.tsx`, `frontend/src/vistas-bd/{tipos.ts,VistaBD.tsx,VistaTabla.tsx,VistaCalendario.tsx,CeldaEditable.tsx}`, `frontend/src/textos/es-MX.ts`, `frontend/pruebas/{cola.spec.ts,medidas.spec.ts,apoyo/medidas.ts,ajustes-a/ajustes.escritorio.spec.ts}`, `docs/estado/wf3/CONVENCIONES-frontend.md`.
Nuevos: `frontend/src/vistas-bd/EdicionVarias.tsx`, `frontend/src/componentes/BarreraErrores.tsx`, `frontend/pruebas/apoyo/banco.ts`, `frontend/pruebas/banco/{index.html,banco.tsx}`, `frontend/pruebas/base.escritorio.spec.ts`, este reporte.

## Verificación (corrida por mí, 12-oct-2026)

- `npx tsc --noEmit` → 0 errores. `npm run build` → OK (dist sin el banco). `npm run verificar:rutas` → «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0». `npm run verificar:hosts` → «U5-08 dist: hosts externos que se pedirían = 0».
- `npx playwright test` (completo, 412 táctil + 1440): **163 passed, 2 skipped, 0 failed** (3.8 min). Los 2 omitidos son pruebas solo táctiles/solo escritorio.
- `cola.spec.ts` 3 corridas: 4/4, 4/4, 4/4; la prueba de I-01 (500) con `--repeat-each 10`: 10/10.
- **U5-07 / N-04**: aria-label de la gráfica = «Gasto por categoría. Despensa sumó $4,000.00 en oct-2026» en claro y oscuro, 412 y 1440; 0 textos en inglés. En `medidas.spec.ts`, Gastos y presupuesto (15 vistas, con P-02 y P-06 Tendencias) y Comidas y compras (10): inglés 0. Inicio Tendencias: 2 descripciones en inglés de `GraficaInicio.tsx` (archivo del área inicio, fuera de mi alcance): se imprimen como «PENDIENTE U5-07 (solicitud a front-inicio)», ver Solicitudes.
- **U5-06 / H-03 (1)**: evento «Walmart · Sin intereses» en el banco a 412 táctil: 378×44 (claro y oscuro), agenda visible y 0 controles < 44; en `/ajustes/05` real a 412 táctil: **344×44** y 0 controles < 44 en la página. A 1440: cuadrícula del mes.
- **H-03 (2)**: «Ver como tabla» 36 px con puntero fino (1440) y 44 px a 412 táctil.
- **H-03 (3)**: descripción emergente en oscuro: fondo rgb(54, 54, 58) = `--ui-surface-raised`, borde rgb(77, 78, 83) = `--ui-border-divider`, texto rgb(242, 242, 244) = `--ui-text-default`, radio 10 px; en claro rgb(250, 250, 251) / rgb(229, 230, 233) / rgb(54, 54, 58). Rellenos de las series iguales antes y después de pasar el puntero.
- **G-05 / N-05** (`medidas.spec.ts`): claro: fondo rgb(250, 250, 251) = raised; píldora «Inicio» **56×32** rgb(240, 241, 255) = selected; icono rgb(93, 96, 159) = icon.selected; iconos 24, 24, 24, 24, 24 px. Oscuro: fondo rgb(54, 54, 58); píldora 56×32 rgb(48, 50, 88); icono rgb(201, 207, 255).
- **H-07/H-08**: prueba «Ctrl+/» (1440): vista previa sumada «Gasto fijo mensual $1,049.00 → $0.00»; `Enter` sin valor no aplica; `Esc` → 0 escrituras y foco de regreso a la casilla; `Enter` → una escritura con las 2 filas; «Seleccionar todas» → 3; segunda propiedad (categoría) con «Aplicar a 3»; `Ctrl+/` sin selección solo avisa. Prueba táctil: «Editar seleccionadas» abre el bloque, «Aplicar» deshabilitado sin vista previa, «$1,348.00» sumado, 18 controles y 0 < 44.
- **H-10**: `npm run verificar:variantes` (Storybook + `contraste-variantes.mjs`): «hallazgos conocidos: 0 de 5 se reproducen» (FV-01, FV-02, FV-03, FV-04 y FV-05 CERRADO), «FALLAS nuevas: 0», «RESULTADO: OK». En el banco: `transition-property` del Botón = color, background-color, border-color; Campo y Chip deshabilitados: borde `dashed`, fondo = `--ui-surface-disabled`, texto = `--ui-text-disabled`, `not-allowed`.
- **N-06**: el banco monta una pantalla que lanza: alerta «No se pudo mostrar esta página» con «Reintentar».
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → salida 0 («Verificación ampliada OK»).

## Supuestos (HECHO / PREGUNTA / PROPUESTA)

- HECHO: la agenda por día por debajo de 768 px sigue la tabla de DESIGN §7.5.4 (360 agenda; 412 agenda y mes compacto opcional; 768+ mes). No se hizo el mes compacto opcional en 412.
- HECHO: el borde de la descripción emergente usa `--ui-border-divider` (C2 no tiene `--ui-border-default`, el nombre que citaba la revisión; FV-03 pedía `border-divider`).
- HECHO: «Ver como tabla» mide 36 px fino (`size.control.fine`) en lugar de los 24 px mínimos que sugería la revisión.
- HECHO: `Ctrl+/` actúa con el foco dentro de la vista (no es global, para que dos tablas en una página no abran a la vez); la selección existe solo en la vista tabla y solo si la pantalla declara `editarVarias`.
- PROPUESTA: en táctil (sin atajos, U4 §0.7) el mismo bloque se abre con el botón visible «Editar seleccionadas». U4 no lo prohíbe ni lo pide; si Axel no lo quiere en Android, basta ocultarlo con `coarse:hidden`.
- PREGUNTA (para front-ajustes-a, no bloquea la base): U4 pide editar en grupo «activo» de los gastos fijos (P-26.02), pero `PlannedItemUpdate` de C1 no tiene un campo `active`; la pantalla debe decir con qué operación se hace o abrir su bloqueo. El ejemplo de CONVENCIONES usa solo `category_id`.

## Bloqueos

Ninguno nuevo.

## Pendientes

- Conectar `editarVarias` en Ajustes 02, 04, 05, 07 (front-ajustes-a) y 09, 13, 14 (front-ajustes-b), con su vista previa sumada de U4: es de los agentes de área (solicitudes 2 y 3).
- `GraficaInicio.tsx` del área inicio sigue con la descripción automática de ECharts (solicitud 1). Al corregirse, quitar la entrada de `PENDIENTES_INGLES` en `frontend/pruebas/medidas.spec.ts`.
- R3-01 (`ajv` sin declarar): requiere `frontend/package.json`, fuera de mi alcance (solicitud 4). Mover `validar-c1.ts` y `tablero.ts` a `pruebas/apoyo/` queda para cuando se declare `ajv`, para no chocar con los agentes de área.
- No hice: H-07/H-08/H-09/H-11 bajos de la ronda 1 de frontend-core (API falsa que devuelve página para GET que no son listas, `max-w-[85vw]`, Inicio sin `TarjetaMetrica`), ni la altura de la barra inferior sumando la zona segura.

## Solicitudes fuera de alcance

1. `frontend/src/pantallas/inicio/GraficaInicio.tsx:49` (front-inicio): cambiar `aria: { enabled: true, decal: { show: false } }` por `aria: { enabled: true, label: { enabled: true, description: `${titulo}. ${resumen}` }, decal: { show: false } }` y agregar `titulo` y `resumen` a las dependencias del efecto; además, como en `Grafica.tsx`, `tooltip.backgroundColor/borderColor/textStyle.color` de `--ui-surface-raised`, `--ui-border-divider`, `--ui-text-default` y `emphasis: { disabled: true }` en cada serie.
2. `frontend/src/pantallas/ajustes-a/` S02, S04, S05, S07 (front-ajustes-a): declarar `editarVarias` en sus `VistaBD` (02: categoría y activo de gastos fijos; 04: tipo, reparto y activa; 05: promociones; 07: tiempos aptos) con la vista previa sumada de cada sección de U4, y una prueba por sección como la de CONVENCIONES §7.
3. `frontend/src/pantallas/ajustes-b/` P2609, P2613, P2614 (front-ajustes-b): lo mismo para 09 (activar/desactivar alertas), 13 (metas) y 14 (modo o nivel, «{k} pasan a Reponer…»).
4. `frontend/package.json` (dueño de la configuración): agregar `"ajv": "8.20.0"` a `devDependencies` (R3-01).

## Paquetes instalados

Ninguno.

## Consumo aproximado

~58 acciones; ~150 KB leídos; ~60 KB escritos.
