// Gráficas de Tendencias (U3 §4, DESIGN §7.7 y §8): <figure> con <figcaption> (título + resumen), «Ver como tabla» que cambia
// la gráfica por la tabla (la tabla queda siempre para lectores de pantalla), ECharts 6 en SVG con colores leídos de los tokens.
import { useEffect, useRef, useState, type ReactNode } from 'react'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { AriaComponent, GridComponent, LegendComponent, MarkAreaComponent, MarkLineComponent, TooltipComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'
import { Boton } from '../../componentes/Boton'
import { cn } from '../../componentes/utils'
import { Vacio } from '../../componentes/Estados'
import { textos } from '../../textos/es-MX'
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
type SerieOpcion = { type?: string; data?: unknown[]; label?: Record<string, unknown>; endLabel?: unknown; labelLayout?: unknown }
const seriesDe = (op: echarts.EChartsCoreOption) => ([] as SerieOpcion[]).concat((op.series ?? []) as SerieOpcion[])
const valorDe = (d: unknown): number | null => (typeof d === 'number' ? d : d && typeof d === 'object' && 'value' in d ? valorDe((d as { value: unknown }).value) : null)
/** Todas las series en 0 o sin valor: estado vacío en lugar de ejes falsos con «$0» (DESIGN §7.7). */
export const sinDatos = (op: echarts.EChartsCoreOption) => seriesDe(op).every((s) => (s.data ?? []).every((d) => !valorDe(d)))
/** Ancho mínimo de cada columna para su etiqueta directa («$1,197.64» a 12 px); si no cabe, el valor queda en la
 *  descripción emergente y en «Ver como tabla» (DESIGN §8.5). */
const ANCHO_ETIQUETA_COLUMNA = 72
/** Etiquetas sin encimarse (DESIGN §7.7): se recorren en vertical o se ocultan si chocan; con leyenda, los nombres al final
 *  de la línea pasan a la leyenda; las de columnas angostas se quitan. En las líneas, la etiqueta del primer punto se alinea
 *  a la izquierda y la del último a la derecha: así no invaden las etiquetas del eje Y ni salen por el borde (§8.5). */
function sinEncimar(op: echarts.EChartsCoreOption, ancho: number): echarts.EChartsCoreOption {
  const series = seriesDe(op)
  const eje = (Array.isArray(op.xAxis) ? op.xAxis[0] : op.xAxis) as { data?: unknown[] } | undefined
  const columnas = Math.max(1, series.filter((s) => s.type === 'bar').length)
  const caben = ancho / Math.max(1, eje?.data?.length ?? 1) / columnas >= ANCHO_ETIQUETA_COLUMNA
  const sinEtiqueta = (d: unknown) => (d && typeof d === 'object' && 'label' in d ? { ...(d as object), label: { ...(d as { label: object }).label, show: false } } : d)
  const conLabel = (d: unknown): d is { label: object } => !!d && typeof d === 'object' && 'label' in d
  const alineada = (d: unknown, i: number, n: number) => (n > 1 && (i === 0 || i === n - 1) && conLabel(d) ? { ...d, label: { ...d.label, align: i === 0 ? 'left' : 'right' } } : d)
  return {
    ...op,
    series: series.map((s) => ({
      ...s,
      labelLayout: s.labelLayout ?? { hideOverlap: true, moveOverlap: 'shiftY' },
      ...(s.endLabel && op.legend ? { endLabel: { show: false } } : {}),
      ...(s.type === 'bar' && !caben ? { label: { ...s.label, show: false }, data: (s.data ?? []).map(sinEtiqueta) } : {}),
      ...(s.type === 'line' ? { data: (s.data ?? []).map((d, i, todos) => alineada(d, i, todos.length)) } : {}),
    })),
  }
}
export interface Tabla { columnas: string[]; filas: string[][]; total?: string[] }
/** Envoltorio de la tabla alterna (U5-06; axe `scrollable-region-focusable`). Oculta: solo `sr-only`, sin desplazamiento
 *  (con `overflow-x-auto` quedaba una región desplazable sin foco). Visible: región desplazable con foco de teclado, anillo
 *  de foco del token y nombre = título de la gráfica. La usan FiguraGrafica y FiguraBarras (Tablero). */
export function ContenedorTabla({ comoTabla, titulo, children }: { comoTabla: boolean; titulo: string; children: ReactNode }) {
  if (!comoTabla) return <div className="sr-only">{children}</div>
  return <div role="region" aria-label={titulo} tabIndex={0} className="overflow-x-auto rounded-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">{children}</div>
}
interface Props { id: string; titulo: string; subtitulo: string; resumen: string; notas?: string[]; filtros?: ReactNode; opcion: () => echarts.EChartsCoreOption; tabla: Tabla; dependencias: unknown[] }
export function FiguraGrafica({ id, titulo, subtitulo, resumen, notas, filtros, opcion, tabla, dependencias }: Props) {
  const [comoTabla, setComoTabla] = useState(false)
  const ref = useRef<HTMLDivElement>(null)
  const tema = usarTema()
  useEffect(() => {
    if (comoTabla || !ref.current) return
    const g = echarts.init(ref.current, null, { renderer: 'svg' })
    const reducido = matchMedia('(prefers-reduced-motion: reduce)').matches
    g.setOption({ animation: !reducido, aria: { enabled: true, label: { enabled: true, description: `${titulo}. ${resumen}` }, decal: { show: false } }, textStyle: { fontFamily: 'Inter, ui-sans-serif, system-ui', color: token('--ui-text-secondary') }, ...sinEncimar(opcion(), ref.current.clientWidth) })
    const ro = new ResizeObserver(() => g.resize())
    ro.observe(ref.current)
    return () => { ro.disconnect(); g.dispose() }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [comoTabla, tema, ...dependencias])
  const vacia = sinDatos(opcion())
  const tablaHtml = (
    <ContenedorTabla comoTabla={comoTabla} titulo={titulo}>
      <table className="w-full text-small tabular-nums">
        <caption className="sr-only">{titulo}</caption>
        <thead><tr>{tabla.columnas.map((c, i) => <th key={c} scope="col" className={cn('whitespace-nowrap px-2 py-1 font-medium text-text-secondary', i ? 'text-right' : 'text-left')}>{c}</th>)}</tr></thead>
        <tbody>{tabla.filas.map((f) => <tr key={f[0]} className="border-t border-border-divider">{f.map((v, i) => (i ? <td key={i} className="whitespace-nowrap px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="whitespace-nowrap px-2 py-1 text-left font-normal">{v}</th>))}</tr>)}</tbody>
        {tabla.total && <tfoot><tr className="border-t border-border-control font-medium">{tabla.total.map((v, i) => (i ? <td key={i} className="px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="px-2 py-1 text-left">{v}</th>))}</tr></tfoot>}
      </table>
    </ContenedorTabla>
  )
  return (
    <figure data-prueba={`grafica-${id}`} className="flex min-w-0 flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6">
      <figcaption className="flex flex-col gap-1">
        <span className="text-small font-medium text-text-secondary">{titulo}</span>
        <span className="text-caption text-text-secondary">{subtitulo}</span>
        <span className="text-small text-text-default tabular-nums">{resumen}</span>
      </figcaption>
      {filtros}
      {vacia && <Vacio texto={textos.vacio.grafica} />}
      {!vacia && !comoTabla && <div ref={ref} aria-hidden className="h-60 w-full tabular-nums tablet:h-72" />}
      {!vacia && tablaHtml}
      {notas?.map((n) => <p key={n} className="text-caption text-text-secondary">{n}</p>)}
      {!vacia && <Boton variante="enlace" className="w-max" aria-pressed={comoTabla} onClick={() => setComoTabla((x) => !x)}>{comoTabla ? t.acciones.verComoGrafica : t.acciones.verComoTabla}</Boton>}
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
