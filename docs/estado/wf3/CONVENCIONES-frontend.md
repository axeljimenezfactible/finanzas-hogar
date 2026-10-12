# Convenciones del frontend · Finanzas Hogar (WF-3)

Autor: agente `frontend-base` (12-oct-2026). Las siguen al pie de la letra los agentes de área. Lo que no esté aquí se pregunta (REGLA SUPREMA), no se inventa.

## 1. Estructura de `frontend/`

| Ruta | Qué es | Quién la toca |
|---|---|---|
| `package.json`, `vite.config.ts`, `tsconfig.json`, `playwright.config.ts`, `index.html`, `scripts/` | Configuración, PWA, scripts de verificación | frontend-base |
| `src/app/` | Enrutador (`rutas.ts`, `enrutador.tsx`), sesión (`sesion.tsx`, P-00 `IniciarSesion.tsx`), tema (`tema.ts`), registro de pantallas (`registro.ts`) | frontend-base |
| `src/shell/` | Barra lateral G-01, barra superior G-02, barra inferior G-05, menú «/» G-03/G-04, hoja P-03/P-22, franja «Sin conexión» G-07 y panel de la cola | frontend-base |
| `src/api/` | `esquema.gen.ts` (generado de `openapi.yaml`, no se edita), `tipos.ts`, `cliente.ts`, `consulta.ts`, `errores.ts`, `formatos.ts`, `conexion.ts`, `lecturas.ts` | frontend-base |
| `src/cola/` | Cola sin conexión en IndexedDB, **solo gastos y peso** (K9) | frontend-base |
| `src/vistas-bd/` | Tabla (TanStack Table 9), tablero (dnd-kit), calendario, galería, `VistaBD`, `CeldaEditable` («Antes → Después»), `ListaOrdenable` | frontend-base |
| `src/componentes/` | Botón, Campo, Chip de persona/tienda, TarjetaMetrica, Avisos, Estados (carga/vacío/error/franja/`ConDatos`), Hoja, Grafica (ECharts) | frontend-base |
| `src/textos/es-MX.ts` | Textos visibles de la base | frontend-base |
| `src/mocks/` | msw para las 216 operaciones de C1 (`npm run dev:mocks`) | frontend-base |
| `src/pantallas/<area>/` | Pantallas, textos, mocks y componentes propios del área | **el agente del área** |
| `pruebas/<area>/` | Pruebas Playwright del área | **el agente del área** |
| `pruebas/apoyo/`, `pruebas/banco/`, `pruebas/*.spec.ts` | API falsa, mediciones (`medidas.ts`), banco de pruebas de la base (Vite de desarrollo, nunca entra a `dist/`) y pruebas de la base | frontend-base |

Nunca se editan a mano: `src/estilos/tokens.css` y `src/estilos/tailwind.theme.css` (copias exactas de `docs/contratos/`, `npm run tokens`), `public/fonts/` (Inter de `docs/design/fuentes/inter/web/`), `src/api/esquema.gen.ts`, `src/mocks/*.gen.json`, `public/icons/`.

## 2. Comandos

```bash
cd /home/deepsec/finanzas-hogar/frontend
npm install                 # una vez
npm run dev                 # Vite en 127.0.0.1:5173 con proxy de /api y /health a FH_API (por omisión https://127.0.0.1:8443, PROVISIONAL)
npm run dev:mocks           # Vite con msw: todas las operaciones de C1 con ejemplos del contrato y datos sintéticos
npm run build               # preparar (tokens, íconos, tipos, mocks) + tsc --noEmit + vite build → dist/
npx tsc --noEmit            # tipos
npm run verificar:rutas     # C3-03: rutas que llama el frontend ⊆ openapi.yaml y ⊆ routes.go
npm run verificar:hosts     # U5-08 estático sobre dist/
FH_CHROMIUM=~/.cache/ms-playwright/chromium_headless_shell-1217/chrome-headless-shell-linux64/chrome-headless-shell npx playwright test
```

`FH_CHROMIUM` solo hace falta mientras el Chromium de `@playwright/test` 1.60.0 (revisión 1223) no esté descargado; con él, `npx playwright test` corre sin red.

## 3. Cómo crea un área sus pantallas

1. Crea `src/pantallas/<area>/indice.ts` (nombre del área en español, minúsculas, sin acentos: `gastos`, `presupuesto`, `comidas`, `compras`, `creditos`, `salud`, `reportes`, `ajustes`, `inicio`).
2. Exporta `pantallas`, un mapa `P-xx → componente perezoso`. El cascarón lo descubre solo (`import.meta.glob` en `src/app/registro.ts`): **nadie edita archivos de la base para registrar una pantalla**.
3. Cada componente recibe `PropsPantalla`: `params` (de la ruta, p. ej. `id`), `sub` (`resumen | detalle | tendencias | calendario`, ya resuelta con la de por defecto de U1) y `consulta` (`URLSearchParams`).
4. El cascarón ya pinta: migas, título (`h1` con el nombre de U1 §2.1), pestañas de subsección, franja «Sin conexión», barra lateral e inferior. La pantalla pinta solo su contenido (empieza con su barra de acciones o con `h2`).
5. Hojas propias del área: se abren con `?hoja=<nombre>` (el menú «/» ya manda `?hoja=ingreso`, `quincena`, `compra`, `pago`, `movimiento`, `precio`, `asignar`, `articulo`, `existencia`, `promocion`, `respaldo`, `generar`). Usa `<Hoja>` de `src/componentes/Hoja.tsx` y cierra con `navegar(conConsulta({ hoja: null }))`.
6. Rutas: `src/app/rutas.ts` (PROVISIONALES, bloqueo `WF3-frontend-rutas`). Para enlazar usa siempre `rutaDe('P-18', { id })` y `<Enlace a=…>`; nunca escribas la dirección a mano.
7. Textos del área en `src/pantallas/<area>/textos.ts` con el mismo estilo que `src/textos/es-MX.ts`: es-MX, sin inglés visible (U5-07), sin fechas ISO visibles.
8. Mocks con estado (opcional): `src/pantallas/<area>/mocks.ts` exporta `manejadores: HttpHandler[]` (msw); tienen prioridad sobre los genéricos.

## 4. Cliente del API (C1 v1.4)

- **Única puerta de red**: `pedir(metodo, '/api/v1/…', { ruta, consulta, cuerpo, idempotencia, senal })`. La ruta se escribe **literal, tal como está en el contrato** (con `{param}`); así la tipa `esquema.gen.ts` y la encuentra el script de C3-03. Prohibido `fetch` directo y rutas armadas con texto.
- Tipos: `Esquemas['Expense']`, `Respuesta<'/api/v1/expenses', 'get'>`, `Cuerpo<…>`, `Consulta<…>` en `src/api/tipos.ts`.
- **Idempotency-Key**: en los POST que la declaran el tipo la exige (`idempotencia: nuevaClave()`); un reintento de la misma captura reutiliza la misma clave.
- **Lecturas en componentes**: `const c = usarConsulta('get', '/api/v1/expenses', { consulta: { period: '2026-10', limit: 50 } })` y `<ConDatos consulta={c}>{(datos) => …}</ConDatos>`: pinta esqueleto, «Esto está tardando más de lo normal · Reintentar» a los 10 s, error con Reintentar, «Cargado a las hh:mm» con la copia del teléfono y «Sin conexión · Esta página aún no se ha cargado en este teléfono». Se recarga sola al volver la conexión. El cuarto argumento `activa=false` la pospone.
- **Paginación** (C1 §5): `limit` 50 por omisión y 200 máximo (K6); `next_cursor` null = fin. `todas('get', ruta, opciones)` junta todas las páginas (solo listas cortas).
- **Errores** (C1 §3): `ErrorApi` con `estado`, `codigo`, `message` (ya en español, se muestra tal cual), `campos` (`e.campo('amount_cents')`), `idPeticion`, `reintentarEn`. `ErrorSinConexion` cuando no responde el servidor de casa. 401 lleva a P-00 y, al entrar, regresa a la página donde estaba.
- **Sin conexión**: solo gastos y peso van a la cola (`src/cola`, ya cableado en P-03/P-22). Toda otra escritura sin conexión: `<Boton deshabilitadoPor={textos.conexion.requiere}>` usando `useConexion()`. Pantallas que no funcionan sin red (P-05, P-19, P-20): `<Franja rol="sinConexion">{textos.conexion.requiere}</Franja>` en lugar del contenido.
- **Datos que cambiaron**: tras guardar un gasto o un peso la hoja emite `fh:datos-cambiaron` (`detail.tipo`); Inicio, P-02 y P-21 lo escuchan y llaman a `recargar()`.

## 5. Formatos (`src/api/formatos.ts`)

| Función | Entrada → pantalla |
|---|---|
| `dinero(c, signo?)` | `123456` → `$1,234.56`; `'gasto'` → `−$1,234.56` (U+2212); `'ingreso'` → `+$1,234.56`. Siempre con `tabular-nums` |
| `aCentavos(texto)` | `«185.5»` → `18550` (ROUND_HALF_UP; nunca flotantes) |
| `fecha(iso)` · `fechaConDia(iso)` · `mes('2026-10')` | `10-oct-2026` · `Vie 30 oct 2026` · `oct-2026` |
| `fechaHora(isoConDesplazamiento)` · `horaCorta(instante)` | `10-oct-2026 21:30` · `21:30` |
| `hoyISO()` · `ahoraConDesplazamiento()` | fecha de hoy en America/Mexico_City · `2026-10-12T09:30:00-06:00` (C1 DateTime) |
| `porcentaje(n)` · `kilos(g)` · `aGramos(texto)` | `80 %` · `72.5 kg` · `«72.5»` → `72500` |

Tasas: tal como se capturaron (K12); no se reformatean.

## 6. Componentes y estilo

- Solo clases de los tokens de C2 (`bg-surface-card`, `text-text-secondary`, `border-border-divider`, `bg-role-expense-surface`, `text-title-2`, `rounded-md`, `shadow-1`, puntos de quiebre `phone: tablet: desktop: wide:`, variante `coarse:` para 44 px). Tailwind trae **solo** esas escalas: un color o tamaño fuera de C2 no existe. Prohibido `dark:` y colores en línea: el tema lo resuelven las variables `--ui-*`.
- `Boton` (`variante`: primario, secundario, fantasma, destructivo, enlace; `tam`: normal, hoja 48 px, icono; `deshabilitadoPor` = motivo visible), `Campo` (etiqueta, `error`, `ayuda`), `ChipPersona`/`ChipTienda`/`ChipSeleccionable` (forma + nombre, orden fijo `ORDEN_PERSONAS`/`ORDEN_TIENDAS`, K41), `TarjetaMetrica` (≤ 6 por Resumen; `alAbrir` baja al detalle), `useAvisos()` (`{ texto, rol, accion: { texto: textos.app.deshacer, alHacer } }`), `Franja`, `Vacio`, `Cargando`, `ErrorCarga`, `ConDatos`, `Hoja`, `Grafica` (`tipo` barras/lineas/area; series en orden fijo con color `--ui-chart-N`, marcador, trazo y patrón K27; tabla alterna incluida).
- Íconos: `lucide-react` (empaquetado, autoalojado). Con `aria-hidden` si hay texto al lado.
- `Grafica`: la descripción accesible es `«{titulo}. {resumen}»` (se pasa a ECharts en `aria.label.description`), así que **el `resumen` debe ser una frase en español con el dinero ya formateado** (`dinero()`), nunca centavos. La descripción emergente, la guía vertical y la leyenda usan tokens (`--ui-surface-raised`, `--ui-border-divider`, `--ui-text-default`, `--ui-viz-baseline`) y se releen al cambiar el tema; las series no se recolorean al pasar el puntero (`emphasis.disabled`). Una pantalla que arme su propia gráfica con ECharts (p. ej. `GraficaInicio`) debe hacer lo mismo: `aria: { enabled: true, label: { description: `${titulo}. ${resumen}` } }`.
- Deshabilitado (§7.0): `Boton` con `deshabilitadoPor` (motivo visible), y `Campo`/`ChipSeleccionable` con `disabled` ya traen relleno `surface-disabled`, texto `text-disabled`, borde punteado y sin hover. El anillo de foco no tiene transición (la del Botón cubre solo color, fondo y borde).
- `BarreraErrores` (`src/componentes/BarreraErrores.tsx`): el cascarón ya envuelve cada pantalla (Layout); si una pantalla lanza al pintarse se ve «No se pudo mostrar esta página · Reintentar» y el resto de la aplicación sigue viva. No hace falta envolver nada más.

## 7. Vistas de base de datos

`<VistaBD def={…}>` con `DefVista<T>`: `clave` (Local), `filas`, `idFila`, `propiedades` (`id`, `nombre`, `valor`, `pintar`, `alinear: 'derecha'` para dinero, `grupo`, `editar`), `vistas` (`tabla | tablero | calendario | galeria`), `fecha` (calendario), `columnasTablero` (`propiedad`, `columnas`, `moverA`), `titulo`, `alAbrir`. La edición en el lugar con efecto numérico declara `editar.vistaPrevia` (devuelve `Cambio[]` para «Antes → Después»); sin vista previa calculada no guarda (U4 §0.2).

- **Calendario**: desde 768 px, mes en cuadrícula; por debajo, agenda por día (DESIGN §7.5.4) con renglones de 44 px. La pantalla no hace nada distinto: usa `vistas: [..., 'calendario']` y `fecha`.
- **Edición en grupo con `Ctrl+/`** (U4 §0.2 punto 2 y §0.7; U1 §5.3): con `def.editarVarias`, la vista tabla muestra una casilla por fila y «Seleccionar todas» (área de toque 36 px fino / 44 px grueso), «n seleccionadas · Editar seleccionadas (Ctrl+/) · Quitar selección», y `Ctrl+/` con el foco dentro de la vista abre el bloque «Editar n registros a la vez»: propiedad → valor nuevo → **vista previa sumada «Antes → Después»** (la calcula la pantalla) → `Enter` o «Aplicar a n» aplica a todas; `Esc` o «Cancelar» cierra sin cambiar ninguna y el foco regresa a donde estaba. Sin vista previa calculada (o si falla: «No se pudo calcular el efecto · Reintentar») no aplica; sin conexión, «Aplicar» queda deshabilitado con «Requiere conexión…». Al aplicar: aviso «Guardado en n registros» y se limpia la selección; si `guardar` lanza, se muestra el mensaje y el bloque sigue abierto.

```ts
// API (src/vistas-bd/tipos.ts)
interface PropiedadVarias { id: string; nombre: string; tipo: 'opcion' | 'texto' | 'dinero' | 'numero' | 'fecha'; opciones?: { valor: string; texto: string }[] }
interface EditarVarias<T> {
  propiedades: PropiedadVarias[]
  vistaPrevia: (filas: T[], propiedad: string, valor: string) => Promise<Cambio[]> // SUMADA sobre las filas
  guardar: (filas: T[], propiedad: string, valor: string) => Promise<void>        // todas o ninguna
}
```

```tsx
// Ejemplo (P-26.02, gastos fijos = planned-items de C1; solo la categoría, que existe en PlannedItemUpdate).
// `valor` llega como texto: el id de la opción o lo tecleado (dinero con aCentavos(), nunca flotantes). Cómo se
// edita «activo» lo decide la pantalla con C1 (PlannedItemUpdate no trae un campo `active`): no se inventa.
<VistaBD def={{ clave: 'P-26.02-fijos', filas, idFila: (f) => f.id, vistas: ['tabla'], titulo: (f) => f.description, propiedades,
  editarVarias: {
    propiedades: [{ id: 'category_id', nombre: t.categoria, tipo: 'opcion', opciones: categorias.map((c) => ({ valor: String(c.id), texto: c.name })) }],
    vistaPrevia: async (sel, _prop, valor) => [
      { etiqueta: t.cambianDeCategoria, antes: '0', despues: String(sel.filter((f) => String(f.category_id) !== valor).length) },
      { etiqueta: t.montoQueSeMueve, antes: dinero(0), despues: dinero(sel.filter((f) => String(f.category_id) !== valor).reduce((s, f) => s + f.amount_cents, 0)) },
    ],
    guardar: async (sel, _prop, valor) => {
      for (const f of sel) await pedir('patch', '/api/v1/planned-items/{planned_item_id}', { ruta: { planned_item_id: f.id }, cuerpo: { category_id: Number(valor) } })
      recargar()
    },
  } }} />
```

La prueba de referencia está en `pruebas/base.escritorio.spec.ts` («H-07/H-08 · Ctrl+/…»): marcar casillas, `page.keyboard.press('Control+/')`, `selectOption` del valor, comprobar el bloque `getByTestId('vista-previa')`, `Escape` sin escrituras y `Enter` con una escritura por fila. En la pantalla, la prueba equivalente cuenta los PATCH con `context.route`.

## 8. Pruebas

- Archivos en `pruebas/<area>/*.spec.ts` (proyecto `movil` 412×915 táctil); los que terminan en `.escritorio.spec.ts` corren también en 1440×900.
- API falsa: `const api = new ApiFalsa(); await api.instalar(context)` (sesión, personas, categorías, gastos y peso con idempotencia; GET desconocidos → página vacía; demás → 501). Para rutas del área registra después tu `context.route(…)` y llama `route.fallback()` en lo que no atiendas (Playwright prueba primero la última ruta registrada).
- Cero hosts externos (U5-08): no se usan fuentes, íconos ni librerías remotas.
- Mediciones compartidas (`pruebas/apoyo/medidas.ts`): `controlesChicos` (44 × 44), `iconosFueraDeNorma` (20/24 px, trazo 1.5, 24 en solo icono y en la barra inferior), `dineroSinTabulares` y `textosEnIngles` (U5-07 en texto, `aria-label`, `title` y `alt`, incluida la descripción automática de ECharts «This is a chart…»).
- Banco de pruebas de la base (`pruebas/apoyo/banco.ts` + `pruebas/banco/`): `const banco = await iniciarBanco()` levanta un Vite de desarrollo en 127.0.0.1 (puerto 5181 o el siguiente libre) y `banco.url('vistabd' | 'grafica' | 'calendario' | 'controles' | 'barrera')` monta la pieza sola. Lo usan solo las pruebas de la base.

## 9. Pantalla de ejemplo completa (P-21 Salud, solo como guía)

```ts
// src/pantallas/salud/indice.ts
import { lazy } from 'react'
export const pantallas = { 'P-21': lazy(() => import('./P21Salud')) }
```

```tsx
// src/pantallas/salud/P21Salud.tsx
import { useEffect } from 'react'
import type { PropsPantalla } from '../../app/registro'
import { usarConsulta } from '../../api/consulta'
import { fecha, kilos } from '../../api/formatos'
import { ConDatos } from '../../componentes/Estados'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { Grafica } from '../../componentes/Grafica'
import { Boton } from '../../componentes/Boton'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { abrirRegistro } from '../../shell/estado'
import { textosSalud as t } from './textos'
export default function P21Salud({ sub }: PropsPantalla) {
  const pesos = usarConsulta('get', '/api/v1/weight-entries', { consulta: { sort: '-measured_date', limit: 50 } })
  useEffect(() => { const f = () => pesos.recargar(); window.addEventListener('fh:datos-cambiaron', f); return () => window.removeEventListener('fh:datos-cambiaron', f) }, [pesos.recargar])
  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end"><Boton variante="primario" onClick={() => abrirRegistro('peso')}>{t.registrarPeso}</Boton></div>
      <ConDatos consulta={pesos}>{(d) => {
        if (sub === 'resumen') return <TarjetaMetrica titulo={t.ultimoPeso} valor={d.items[0] ? kilos(d.items[0].weight_g) : '—'} nota={d.items[0] && fecha(d.items[0].measured_date)} />
        if (sub === 'tendencias') return <Grafica titulo={t.peso} resumen={t.resumenPeso} tipo="lineas" categorias={d.items.map((p) => fecha(p.measured_date))} series={[{ nombre: t.peso, valores: d.items.map((p) => p.weight_g) }]} formato={kilos} />
        return <VistaBD def={{ clave: 'P-21-pesos', filas: d.items, idFila: (p) => p.id, vistas: ['tabla', 'calendario'], fecha: (p) => p.measured_date, titulo: (p) => kilos(p.weight_g),
          propiedades: [{ id: 'measured_date', nombre: t.fecha, valor: (p) => p.measured_date, pintar: (p) => fecha(p.measured_date) }, { id: 'weight_g', nombre: t.peso, valor: (p) => p.weight_g, pintar: (p) => kilos(p.weight_g), alinear: 'derecha' }] }} />
      }}</ConDatos>
    </div>
  )
}
```

```ts
// pruebas/salud/salud.spec.ts
import { expect, test } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
test('P-21 muestra el último peso', async ({ page, context }) => {
  await new ApiFalsa().instalar(context)
  await context.route('**/api/v1/weight-entries?*', (r) => r.request().method() === 'GET' ? r.fulfill({ json: { items: [{ id: 1, person_id: 1, measured_date: '2026-10-05', weight_g: 84600, created_at: '2026-10-05T08:00:00-06:00', updated_at: '2026-10-05T08:00:00-06:00' }], next_cursor: null } }) : r.fallback())
  await page.goto('/salud')
  await expect(page.getByText('84.6 kg')).toBeVisible()
})
```
