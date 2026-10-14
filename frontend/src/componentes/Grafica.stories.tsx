// Variantes de la gráfica (DESIGN.md §7.7 y §8): barras, líneas y área, de 1 a 8 series, con marcador, trazo y patrón (K27).
import type { Meta, StoryObj } from '@storybook/react-vite'
import { Grafica, type Serie } from './Grafica'

const meses = ['may-2026', 'jun-2026', 'jul-2026', 'ago-2026', 'sep-2026', 'oct-2026']
const NOMBRES = ['Comida', 'Casa', 'Transporte', 'Salud', 'Ocio', 'Servicios', 'Ropa', 'Otros']
const serie = (i: number): Serie => ({ nombre: NOMBRES[i], valores: meses.map((_, j) => 150000 + ((i * 37 + j * 53) % 11) * 40000) })
const series = (n: number) => Array.from({ length: n }, (_, i) => serie(i))

const meta = {
  id: 'grafica',
  title: 'Componentes/Gráfica',
  component: Grafica,
  args: { titulo: 'Gasto por mes', resumen: 'El gasto subió 12 % respecto a septiembre', tipo: 'barras', categorias: meses, series: series(1) },
  argTypes: { tipo: { control: 'inline-radio', options: ['barras', 'lineas', 'area'] } },
  decorators: [(Historia) => <div className="max-w-3xl"><Historia /></div>],
} satisfies Meta<typeof Grafica>
export default meta
type Historia = StoryObj<typeof meta>

/** Editable con los controles de Storybook. */
export const Editable: Historia = {}
export const BarrasUnaSerie: Historia = { name: 'Barras, 1 serie' }
export const LineasTresSeries: Historia = { name: 'Líneas, 3 series', args: { tipo: 'lineas', series: series(3), titulo: 'Gasto por categoría' } }
export const AreaApilada: Historia = { name: 'Área apilada, 2 series', args: { tipo: 'area', series: series(2), titulo: 'Devengado contra efectivo' } }
export const OchoSeries: Historia = { name: 'Líneas, 8 series (máximo)', args: { tipo: 'lineas', series: series(8), titulo: 'Gasto por categoría' } }
export const BarrasOchoSeries: Historia = { name: 'Barras, 8 series con patrón', args: { tipo: 'barras', series: series(8), titulo: 'Gasto por categoría' } }
export const SinDatos: Historia = { name: 'Sin datos', args: { series: [] } }
