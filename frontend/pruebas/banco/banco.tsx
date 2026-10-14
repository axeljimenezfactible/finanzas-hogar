// Banco de pruebas de las piezas de la base (solo pruebas; nunca entra a dist/). Lo sirve un Vite de desarrollo que
// levanta pruebas/apoyo/banco.ts; cada caso se elige con ?caso=. Datos de ejemplo inventados, en centavos enteros.
import { StrictMode, useState } from 'react'
import { createRoot } from 'react-dom/client'
import '../../src/estilos/app.css'
import { ProveedorAvisos } from '../../src/componentes/Avisos'
import { ProveedorIconos } from '../../src/componentes/Icono'
import { Grafica } from '../../src/componentes/Grafica'
import { Boton } from '../../src/componentes/Boton'
import { Campo } from '../../src/componentes/Campo'
import { ChipSeleccionable } from '../../src/componentes/Chip'
import { BarreraErrores } from '../../src/componentes/BarreraErrores'
import { VistaBD } from '../../src/vistas-bd/VistaBD'
import { dinero, hoyISO } from '../../src/api/formatos'
import type { Cambio } from '../../src/vistas-bd/tipos'

interface GastoFijo { id: number; nombre: string; categoria: string; monto_cents: number; activo: boolean }
const FIJOS: GastoFijo[] = [
  { id: 1, nombre: 'Internet', categoria: 'servicios', monto_cents: 59_900, activo: true },
  { id: 2, nombre: 'Gimnasio', categoria: 'salud', monto_cents: 45_000, activo: true },
  { id: 3, nombre: 'Streaming', categoria: 'ocio', monto_cents: 29_900, activo: false },
]
const CATEGORIAS = [{ valor: 'servicios', texto: 'Servicios' }, { valor: 'salud', texto: 'Salud' }, { valor: 'ocio', texto: 'Ocio' }]
declare global { interface Window { __guardados: { ids: number[]; propiedad: string; valor: string }[]; __previas: number } }
window.__guardados = []
window.__previas = 0
const activosCents = (fs: GastoFijo[]) => fs.filter((f) => f.activo).reduce((s, f) => s + f.monto_cents, 0)

function CasoVistaBD() {
  const [filas, setFilas] = useState(FIJOS)
  return (
    <VistaBD def={{
      clave: 'banco-fijos', filas, idFila: (f) => f.id, vistas: ['tabla', 'galeria'], titulo: (f) => f.nombre,
      propiedades: [
        { id: 'nombre', nombre: 'Nombre', valor: (f) => f.nombre },
        { id: 'categoria', nombre: 'Categoría', valor: (f) => f.categoria, pintar: (f) => CATEGORIAS.find((c) => c.valor === f.categoria)?.texto },
        { id: 'monto_cents', nombre: 'Monto', valor: (f) => f.monto_cents, pintar: (f) => dinero(f.monto_cents), alinear: 'derecha' },
        { id: 'activo', nombre: 'Activo', valor: (f) => (f.activo ? 'Sí' : 'No') },
      ],
      editarVarias: {
        propiedades: [
          { id: 'activo', nombre: 'Activo', tipo: 'opcion', opciones: [{ valor: 'si', texto: 'Sí' }, { valor: 'no', texto: 'No' }] },
          { id: 'categoria', nombre: 'Categoría', tipo: 'opcion', opciones: CATEGORIAS },
        ],
        // Vista previa SUMADA: total mensual de gastos fijos activos antes y después de aplicar a las filas elegidas.
        vistaPrevia: async (sel, propiedad, valor): Promise<Cambio[]> => {
          window.__previas++
          if (propiedad !== 'activo') return [{ etiqueta: 'Registros que cambian de categoría', antes: '0', despues: String(sel.filter((f) => f.categoria !== valor).length) }]
          const ids = new Set(sel.map((f) => f.id))
          const despues = filas.map((f) => (ids.has(f.id) ? { ...f, activo: valor === 'si' } : f))
          return [{ etiqueta: 'Gasto fijo mensual', antes: dinero(activosCents(filas)), despues: dinero(activosCents(despues)) }]
        },
        guardar: async (sel, propiedad, valor) => {
          window.__guardados.push({ ids: sel.map((f) => f.id), propiedad, valor })
          const ids = new Set(sel.map((f) => f.id))
          setFilas((fs) => fs.map((f) => (!ids.has(f.id) ? f : propiedad === 'activo' ? { ...f, activo: valor === 'si' } : { ...f, categoria: valor })))
        },
      },
    }} />
  )
}
function CasoGrafica() {
  return <Grafica titulo="Gasto por categoría" resumen={`Despensa sumó ${dinero(400_000)} en oct-2026`} tipo="barras" categorias={['sep-2026', 'oct-2026']}
    series={[{ nombre: 'Despensa', valores: [0, 400_000] }, { nombre: 'Servicios', valores: [120_050, 98_000] }]} />
}
function CasoCalendario() {
  const mes = hoyISO().slice(0, 7)
  const filas = [{ id: 1, fecha: `${mes}-05`, titulo: 'Walmart · Sin intereses' }, { id: 2, fecha: `${mes}-05`, titulo: 'Bodega Aurrera · 3 meses' }, { id: 3, fecha: `${mes}-20`, titulo: 'Mercado · contado' }]
  return <VistaBD def={{ clave: 'banco-calendario', filas, idFila: (f) => f.id, vistas: ['calendario'], fecha: (f) => f.fecha, titulo: (f) => f.titulo, alAbrir: () => undefined,
    propiedades: [{ id: 'titulo', nombre: 'Título', valor: (f) => f.titulo }] }} />
}
function CasoControles() {
  return (
    <div className="flex flex-col gap-3">
      <Boton variante="primario">Guardar</Boton>
      <Campo etiqueta="Monto deshabilitado" disabled defaultValue="12.50" />
      <ChipSeleccionable seleccionado={false} disabled>Ocio</ChipSeleccionable>
    </div>
  )
}
function Explota(): never { throw new Error('forma inesperada') }
const CASOS: Record<string, () => React.ReactNode> = {
  vistabd: CasoVistaBD, grafica: CasoGrafica, calendario: CasoCalendario, controles: CasoControles,
  barrera: () => <BarreraErrores><Explota /></BarreraErrores>,
}
const Caso = CASOS[new URLSearchParams(location.search).get('caso') ?? ''] ?? (() => 'Caso desconocido')
createRoot(document.getElementById('raiz')!).render(<StrictMode><ProveedorIconos><ProveedorAvisos><main className="p-4"><Caso /></main></ProveedorAvisos></ProveedorIconos></StrictMode>)
