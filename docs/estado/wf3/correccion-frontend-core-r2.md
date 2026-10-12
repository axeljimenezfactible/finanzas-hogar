# Corrección frontend-core · ronda 2 · WF-3a · 12-oct-2026

Agente: `correccion-frontend-core-r2` (Opus).
Encargo: corregir N-01 (alta) y N-02 (media) de `REVISION-frontend-core-r2.md`.
Alcance: solo los archivos de `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras` que nombran esos hallazgos.
No se reabrió ninguna decisión ni se tocó una cifra. No hubo `git add` ni `git commit`.

## Hecho

| Hallazgo | Estado | Qué se cambió |
|---|---|---|
| **N-01** (alta) U5-06 · DESIGN §4.6 | **Corregido** | Ver «N-01» abajo. |
| **N-02** (media) DESIGN §5.2 · U5-06 | **Corregido** | Ver «N-02» abajo. |

### N-01 · botón «Más opciones» de F-1

**Corrección en el código**
- `src/shell/HojaRegistro.tsx:120`: el `<button>` suelto pasa a `<Boton variante="enlace" className="self-start">`, el componente de DESIGN §7.2.
- Hereda `h-9 coarse:h-11`, el foco visible y los estados de enlace.
- Medido con una prueba propia temporal (ya borrada):

| Botón | 412 px táctil | 1440 px fino (`size.control.fine`) |
|---|---|---|
| «Más opciones» | **123×44** (antes 93×20) | 123×36 |
| «Menos opciones» | **141×44** (antes 111×20) | 141×36 |

**Ampliación de la prueba** (`pruebas/medidas.spec.ts`)
- Prueba nueva: «U5-06 a 412 px táctil · estados abiertos». Mide 11 estados con las mismas tres reglas (toque ≥ 44 × 44, iconos de §5.2 y dinero tabular):
  1. Hoja F-1, pestaña Gasto, en reposo.
  2. F-1 con «Más opciones» abierto, el chip «Más» (todas las categorías) y forma de pago «Crédito» (incluye «Ir a compra a crédito»).
  3. F-1, pestaña Ingreso.
  4. F-1, pestaña Peso.
  5. F-1, pestaña Ticket.
  6. Franja «Sin conexión · 1 en cola», con el aviso y su «Deshacer».
  7. Panel de la cola.
  8. Panel de la cola al confirmar un descarte.
  9. Menú «/» con la lista de comandos.
  10. Vista previa «Antes → Después» de la edición en el lugar (P-04, monto).
  11. Esa vista previa en «No se pudo calcular el efecto · Reintentar».
- Al medir el panel de la cola apareció un defecto: el importe «Gasto $99.00» se mostraba sin `tabular-nums` (DESIGN §3.4, U5-06).
  - Se corrigió en `src/shell/FranjaSinConexion.tsx:38`: se añadió `tabular-nums` al `span` del importe. Es un solo cambio de clase en el archivo de `frontend-base` donde vive el panel de la cola que nombra N-01.
  - Si el integrador lo considera fuera del alcance, se puede revertir, pero entonces la prueba ampliada falla en ese punto.

### N-02 · iconos de 24 px en controles de solo icono

**Corrección en el código**
- Los 14 iconos pasan de `className="size-5"` a `className="size-6" size={ICONO_MD}`. Así el trazo absoluto se calcula sobre 24 px y da 1.5 px dibujados.
- Donde faltaba, se añadió `import { ICONO_MD } from '…/componentes/Icono'`.

| Dueño | Ruta (líneas actuales) | Control |
|---|---|---|
| frontend-base | `src/componentes/Hoja.tsx:20` (import en :5) | «Cerrar» de toda hoja |
| frontend-base | `src/shell/Barras.tsx:16` | «Abrir menú» |
| front-inicio | `src/pantallas/inicio/Calendario.tsx:80,82` (import en :5) | «Mes anterior» y «Mes siguiente» |
| front-gastos-presupuesto | `src/pantallas/gastos-presupuesto/comun.tsx:38,40` (import en :5) | «Mes anterior» y «Mes siguiente» (P-02, P-06) |
| front-comidas-compras | `src/pantallas/comidas-compras/comun.tsx:58,60,146` (import en :6) | «Semana anterior», «Semana siguiente», «Más opciones» |
| front-comidas-compras | `src/pantallas/comidas-compras/P10Plan.tsx:267,269` (el import ya estaba) | «Mes anterior» y «Mes siguiente» |
| front-comidas-compras | `src/pantallas/comidas-compras/P12Receta.tsx:81,83,191` (import en :4) | «Menos porción», «Más porción», «Quitar ingrediente» |

- Medido (prueba temporal) a 412 táctil y a 1440 fino: los iconos de «Mes anterior», «Mes siguiente», «Abrir menú» y «Cerrar» miden **24×24**.
- Siguen en 20 px:
  - las flechas «Desplegar» y «Plegar» de la barra lateral (§5.2: barra lateral = `size.icon.sm`);
  - los iconos que van junto a texto.

**Regla nueva en la prueba** (`pruebas/apoyo/medidas.ts`, `iconosFueraDeNorma`)
- Si el icono vive en un botón o enlace sin texto visible, debe medir 24 px.
  - Cuentan como control: `button`, `a[href]` y los roles `button`, `link`, `tab`, `menuitem` y `option`.
  - El texto `sr-only`, el oculto y el de los `svg` no cuenta como texto visible.
- Queda fuera la barra lateral (`nav[aria-label="Navegación principal"]`).
- Cada falla se informa con el control y la regla, por ejemplo: `solo icono → 24 px: button «Mes anterior»`.
- Comprobación de que la regla sí detecta el defecto (spec temporal, ya borrado):
  - en `/gastos` se devolvió a 20 px (trazo 1.8) el icono de «Mes anterior» y el de «Abrir menú»;
  - la regla marcó exactamente esos 2;
  - con el panel lateral móvil abierto (26 iconos, con las flechas de 20 px), no marcó nada más.

## Archivos

- `frontend/src/shell/HojaRegistro.tsx` (N-01)
- `frontend/src/shell/FranjaSinConexion.tsx` (`tabular-nums` del panel de la cola, encontrado por la prueba ampliada de N-01)
- `frontend/src/componentes/Hoja.tsx` (N-02)
- `frontend/src/shell/Barras.tsx` (N-02)
- `frontend/src/pantallas/inicio/Calendario.tsx` (N-02)
- `frontend/src/pantallas/gastos-presupuesto/comun.tsx` (N-02)
- `frontend/src/pantallas/comidas-compras/comun.tsx` (N-02)
- `frontend/src/pantallas/comidas-compras/P10Plan.tsx` (N-02)
- `frontend/src/pantallas/comidas-compras/P12Receta.tsx` (N-02)
- `frontend/pruebas/medidas.spec.ts` (N-01: estados abiertos; N-02: mensajes con la regla)
- `frontend/pruebas/apoyo/medidas.ts` (N-02: regla de solo icono)
- `docs/estado/wf3/correccion-frontend-core-r2.md` (este reporte)

## Verificación

- **C3-03 (build):** `npm run build` → exit 0, **0 `error TS`** (tsc dentro del build), precache de 80 entradas. `npx tsc --noEmit -p .` sin salida.
- **U5-06 (objetivos táctiles y estados abiertos):** `npx playwright test pruebas/medidas.spec.ts` (con `FH_CHROMIUM`) → **4 passed**:

  | Prueba | Controles | Menores de 44 px | Iconos | Fuera de norma | Dinero sin tabulares |
  |---|---|---|---|---|---|
  | Inicio (4 vistas) | 83 | 0 | 76 | 0 | 0 |
  | Gastos y presupuesto (15 vistas) | 341 | 0 | 169 | 0 | 0 |
  | Comidas y compras (10 vistas) | 239 | 0 | 100 | 0 | 0 |
  | **Estados abiertos (11)** | 283 | **0** | 184 | **0** | **0** |

  - «Fuera de norma» quiere decir: fuera de 20/24 px, trazo distinto de 1.5 px, o un botón de solo icono sin 24 px.
  - La prueba de la casilla de P-13 da 20×20 dentro de un área de 44×44.
  - Antes de añadir `tabular-nums` en el panel de la cola, la prueba de estados abiertos fallaba con «Panel de la cola: «Gasto $99.00»» ×2. Después pasa.
- **U5-06 (iconos, N-02):** ver la regla y su comprobación en «N-02».
- **C3-03 (rutas), U5-08, U5-09, U5-13 y e2e de todas las áreas:** `make verify-wf3-frontend` → **exit 0**:
  - «C3-03 frontend: 202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)»;
  - «U5-08 dist: hosts externos que se pedirían = 0» y «U5-08 hosts pedidos: localhost; externos: 0»;
  - «U5-09 toques…»: 3 en cada recorrido;
  - «U5-13 registros: gastos 1, pesos 1 (antes 0 y 0); reintento forzado → 200»;
  - «U5-13 paso 7 (500)» y «(401)» con franja, «1 en cola» y aviso;
  - **142 passed** (3.3 min).
- **Verificación global:** `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → **exit 2** en `verify-go`.
  - Causa: «FALLA gofmt: internal/reportes/reportes_test.go, internal/reportes/vistas.go».
  - Es backend de S7 y está fuera del alcance de esta corrección; ningún archivo del frontend interviene.
  - Lo anterior a ese paso pasa: goose up/down/up y `go build`/`go vet`.

## Supuestos

- **HECHO:** «botón de solo icono» quiere decir un control sin texto visible. El nombre va en `aria-label` o en texto `sr-only`. Los iconos que van junto a texto siguen en 20 px (`size.icon.sm`, DESIGN §5.2: «junto a texto de 14 y 16 px»).
- **HECHO:** la barra lateral queda fuera de la regla de 24 px. Lo dice el hallazgo y lo dice §5.2 («barra lateral» en `size.icon.sm`).
- **HECHO:** el «Más opciones» de F-1 usa la variante `enlace` de `Boton` tal como está. Eso trae el `px-3` del tamaño normal (tailwind-merge) y el `font-medium`, igual que los otros 25 usos de esa variante. No se ajustó nada a mano.

## Bloqueos

Ninguno nuevo.

## Pendientes

- N-03 (baja, `Barras.tsx:26`, «Buscar» de 32×32 sin `coarse:size-11`) y los bajos de la ronda 1 (H-07, H-08, H-09, H-11): no estaban en este encargo.
- `make verify-all` sigue roto por gofmt en `backend/internal/reportes/` (dueño: S7).

## Propuestas para Axel

- **PROPUESTA (frontend-base):** añadir a `CONVENCIONES-frontend.md` §6 esta línea: «Botón de solo icono: icono con `className="size-6" size={ICONO_MD}` (DESIGN §5.2); junto a texto, 20 px por omisión». Hoy la convención solo dice «Con aria-hidden si hay texto al lado».
- **PROPUESTA (fuera de este grupo, para la revisión de áreas):** el mismo grep encontró, fuera de este grupo, botones de solo icono que no cumplen §5.2:
  - `src/pantallas/ajustes-a/S06Tiendas.tsx:199,201`: `size-6` sin `size={ICONO_MD}`. El trazo sale de 1.8 px dibujados.
  - `src/pantallas/creditos-salud-reportes/P21Salud.tsx:139`: `size-4`, 16 px.
  - `src/componentes/Boton.stories.tsx:70-72`: ejemplos de Storybook de solo icono con `size-5`.
  - `medidas.spec.ts` no recorre esas áreas.
- **PROPUESTA (frontend-base):** N-03: `Barras.tsx:26` `size-8` → `size-9 coarse:size-11`, igual que `Boton tam="icono"`.

## Paquetes instalados

Ninguno.

## Consumo aproximado

- Acciones: ~35.
- Leído: ~110 KB.
- Escrito: ~22 KB (código y pruebas ~14 KB; reporte ~8 KB).
- Temporales en `datos/tmp/corr-fcore-r2/` (copias de respaldo y bitácoras): se borraron al terminar. Los dos specs temporales de `pruebas/zz-tmp-*` también se borraron.
