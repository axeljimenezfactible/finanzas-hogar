// Gráficas de Tendencias (U3 §4, DESIGN §7.7 y §8): <figure> con <figcaption> (título + resumen), «Ver como tabla» que cambia
// la gráfica por la tabla (la tabla queda siempre para lectores de pantalla), ECharts 6 en SVG con colores leídos de los tokens.
import { useEffect, useRef, useState, type ReactNode } from 'react'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { AriaComponent, GridComponent, LegendComponent, MarkAreaComponent, MarkLineComponent, TooltipComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'
import { Boton } from '../../componentes/Boton'
import { cn } from '../../componentes/utils'
import { textosInicio as t } from './textos'
echarts.use([BarChart, LineChart, AriaComponent, GridComponent, LegendComponent, MarkAreaComponent, MarkLineComponent, TooltipComponent, SVGRenderer])

export const token = (v: string) => getComputedStyle(document.documentElement).getPropertyValue(v).trim()
/** «$12.3 mil»: formato compacto, solo en ejes (U3 §4, DESIGN §8.5). Vive en api/formatos.ts (`dineroCompacto`) para que
 *  Grafica y Tendencias usen el mismo. */
export { dineroCompacto as compacto } from '../../api/formatos'
/** Índices con etiqueta (DESIGN §8.5): el último valor y los extremos. Es el mismo de componentes/Grafica. */
export { selectivos } from '../../componentes/Grafica'
function usarTema() {
  const [n, setN] = useState(0)
  useEffect(() => {
    const mq = matchMedia('(prefers-color-scheme: dark)'), f = () => setN((x) => x + 1)
    const obs = new MutationObserver(f)
    obs.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme', 'class'] })
    mq.addEventListener('change', f)
    return () => { obs.disconnect(); mq.removeEventListener('change', f) }
  }, [])
  return n
}
export interface Tabla { columnas: string[]; filas: string[][]; total?: string[] }
interface Props { id: string; titulo: string; subtitulo: string; resumen: string; notas?: string[]; filtros?: ReactNode; opcion: () => echarts.EChartsCoreOption; tabla: Tabla; dependencias: unknown[] }
export function FiguraGrafica({ id, titulo, subtitulo, resumen, notas, filtros, opcion, tabla, dependencias }: Props) {
  const [comoTabla, setComoTabla] = useState(false)
  const ref = useRef<HTMLDivElement>(null)
  const tema = usarTema()
  useEffect(() => {
    if (comoTabla || !ref.current) return
    const g = echarts.init(ref.current, null, { renderer: 'svg' })
    const reducido = matchMedia('(prefers-reduced-motion: reduce)').matches
    g.setOption({ animation: !reducido, aria: { enabled: true, decal: { show: false } }, textStyle: { fontFamily: 'Inter, ui-sans-serif, system-ui', color: token('--ui-text-secondary') }, ...opcion() })
    const ro = new ResizeObserver(() => g.resize())
    ro.observe(ref.current)
    return () => { ro.disconnect(); g.dispose() }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [comoTabla, tema, ...dependencias])
  const tablaHtml = (
    <div className={cn('overflow-x-auto', !comoTabla && 'sr-only')}>
      <table className="w-full text-small tabular-nums">
        <caption className="sr-only">{titulo}</caption>
        <thead><tr>{tabla.columnas.map((c, i) => <th key={c} scope="col" className={cn('whitespace-nowrap px-2 py-1 font-medium text-text-secondary', i ? 'text-right' : 'text-left')}>{c}</th>)}</tr></thead>
        <tbody>{tabla.filas.map((f) => <tr key={f[0]} className="border-t border-border-divider">{f.map((v, i) => (i ? <td key={i} className="whitespace-nowrap px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="whitespace-nowrap px-2 py-1 text-left font-normal">{v}</th>))}</tr>)}</tbody>
        {tabla.total && <tfoot><tr className="border-t border-border-control font-medium">{tabla.total.map((v, i) => (i ? <td key={i} className="px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="px-2 py-1 text-left">{v}</th>))}</tr></tfoot>}
      </table>
    </div>
  )
  return (
    <figure data-prueba={`grafica-${id}`} className="flex min-w-0 flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6">
      <figcaption className="flex flex-col gap-1">
        <span className="text-small font-medium text-text-secondary">{titulo}</span>
        <span className="text-caption text-text-secondary">{subtitulo}</span>
        <span className="text-small text-text-default tabular-nums">{resumen}</span>
      </figcaption>
      {filtros}
      {!comoTabla && <div ref={ref} aria-hidden className="h-60 w-full tabular-nums tablet:h-72" />}
      {tablaHtml}
      {notas?.map((n) => <p key={n} className="text-caption text-text-secondary">{n}</p>)}
      <Boton variante="enlace" className="w-max" aria-pressed={comoTabla} onClick={() => setComoTabla((x) => !x)}>{comoTabla ? t.acciones.verComoGrafica : t.acciones.verComoTabla}</Boton>
    </figure>
  )
}
/** Ejes y rejilla comunes (DESIGN §8.1): un solo eje Y, sin retícula vertical, textos en los tokens de texto. */
export function ejes(categorias: string[], formatoY: (v: number) => string, ancho: boolean) {
  const sec = token('--ui-text-secondary')
  return {
    grid: { left: 8, right: ancho ? 120 : 16, top: 32, bottom: ancho ? 8 : 40, containLabel: true },
    legend: { [ancho ? 'right' : 'bottom']: 0, ...(ancho ? { top: 'middle', orient: 'vertical' } : {}), textStyle: { color: sec } },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: categorias, axisLine: { lineStyle: { color: token('--ui-viz-baseline') } }, axisTick: { show: false }, axisLabel: { color: sec, fontSize: 12 }, splitLine: { show: false } },
    yAxis: { type: 'value', splitLine: { lineStyle: { color: token('--ui-viz-grid') } }, axisLabel: { color: sec, fontSize: 12, formatter: formatoY } },
  }
}
