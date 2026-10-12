// Gráfica base (DESIGN §7.7 y §8.6, K27) con Apache ECharts 6 (renderizado SVG, módulos sueltos).
// Paleta categórica de C2 (--ui-chart-1..8) leída del tema activo; cada serie con marcador, trazo y patrón (decal).
import { useEffect, useRef, useState } from 'react'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { AriaComponent, GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'
import { textos } from '../textos/es-MX'
import { dinero } from '../api/formatos'
import { Vacio } from './Estados'
echarts.use([BarChart, LineChart, AriaComponent, GridComponent, LegendComponent, TooltipComponent, SVGRenderer])

export const MARCADORES = ['circle', 'rect', 'triangle', 'diamond', 'path://M0,0L10,0L5,10Z', 'path://M4,0H6V4H10V6H6V10H4V6H0V4H4Z', 'path://M1,0L5,4L9,0L10,1L6,5L10,9L9,10L5,6L1,10L0,9L4,5L0,1Z', 'emptyCircle']
export const TRAZOS: (number[] | 'solid')[] = ['solid', [8, 4], [2, 4], [8, 4, 2, 4], [4, 4], [12, 4], [8, 3, 2, 3, 2, 3], [2, 2]]
export function patron(i: number, superficie: string) {
  const rayas = (rot: number, paso: number) => ({ symbol: 'rect', symbolSize: 1, dashArrayX: [1, 0], dashArrayY: [1.5, paso - 1.5], rotation: rot, color: superficie })
  const puntos = (tam: number, paso: number) => ({ symbol: 'circle', symbolSize: tam / paso, dashArrayX: [[paso, 0]], dashArrayY: [paso, 0], color: superficie })
  return [undefined, rayas(Math.PI / 4, 6), rayas(-Math.PI / 4, 6), puntos(2, 6), { symbol: 'rect', symbolSize: 1, dashArrayX: [[1, 7]], dashArrayY: [1, 7], rotation: Math.PI / 4, color: superficie }, rayas(Math.PI / 4, 10), rayas(-Math.PI / 4, 10), puntos(3, 8)][i]
}
const css = (v: string) => getComputedStyle(document.documentElement).getPropertyValue(v).trim()
function usarTemaVivo() {
  const [n, setN] = useState(0)
  useEffect(() => {
    const mq = matchMedia('(prefers-color-scheme: dark)')
    const f = () => setN((x) => x + 1)
    const obs = new MutationObserver(f)
    obs.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
    mq.addEventListener('change', f)
    return () => { obs.disconnect(); mq.removeEventListener('change', f) }
  }, [])
  return n
}
export interface Serie { nombre: string; valores: number[] }
interface Props { titulo: string; resumen: string; tipo: 'barras' | 'lineas' | 'area'; categorias: string[]; series: Serie[]; formato?: (v: number) => string }
export function Grafica({ titulo, resumen, tipo, categorias, series, formato = (v) => dinero(v) }: Props) {
  const ref = useRef<HTMLDivElement>(null)
  const tema = usarTemaVivo()
  const visibles = series.slice(0, 8) // 8 series máximo; la novena va en «Otros» (la arma quien llama)
  useEffect(() => {
    if (!ref.current || !visibles.length) return
    const graf = echarts.init(ref.current, null, { renderer: 'svg' })
    const sup = css('--ui-surface-card'), sec = css('--ui-text-secondary')
    graf.setOption({
      aria: { enabled: true, decal: { show: true } },
      textStyle: { fontFamily: 'Inter, ui-sans-serif, system-ui', color: sec },
      grid: { left: 8, right: 8, top: 24, bottom: 8, containLabel: true },
      legend: visibles.length > 1 ? { bottom: 0, textStyle: { color: sec } } : undefined,
      tooltip: { trigger: 'axis', valueFormatter: (v: unknown) => formato(Number(v)) },
      xAxis: { type: 'category', data: categorias, axisLine: { lineStyle: { color: css('--ui-viz-baseline') } }, axisLabel: { color: sec } },
      yAxis: { type: 'value', splitLine: { lineStyle: { color: css('--ui-viz-grid') } }, axisLabel: { color: sec, formatter: (v: number) => formato(v) } },
      series: visibles.map((s, i) => ({
        name: s.nombre,
        type: tipo === 'barras' ? 'bar' : 'line',
        data: s.valores,
        animation: false,
        color: css(`--ui-chart-${i + 1}`),
        symbol: MARCADORES[i],
        symbolSize: 8,
        lineStyle: { type: TRAZOS[i], width: 2 },
        areaStyle: tipo === 'area' ? { opacity: 0.2 } : undefined,
        stack: tipo === 'area' ? 'total' : undefined,
        itemStyle: { decal: patron(i, sup) },
      })),
    })
    const ro = new ResizeObserver(() => graf.resize())
    ro.observe(ref.current)
    return () => { ro.disconnect(); graf.dispose() }
  }, [tema, tipo, categorias, visibles, formato])
  if (!visibles.length) return <Vacio texto={textos.vacio.grafica} />
  return (
    <figure className="flex flex-col gap-2 rounded-lg border border-border-divider bg-surface-card p-4">
      <figcaption>
        <span className="block text-small font-medium text-text-secondary">{titulo}</span>
        <span className="block text-small text-text-default">{resumen}</span>
      </figcaption>
      <div ref={ref} tabIndex={0} aria-label={`${titulo}. ${resumen}`} className="h-50 w-full tablet:h-70" />
      <details className="text-small">
        <summary className="cursor-pointer text-text-link coarse:py-3">{textos.grafica.tablaAlterna}</summary>
        <table className="mt-2 w-full tabular-nums">
          <thead><tr><th className="text-left" scope="col">{titulo}</th>{visibles.map((s) => <th key={s.nombre} scope="col" className="text-right">{s.nombre}</th>)}</tr></thead>
          <tbody>{categorias.map((c, j) => <tr key={c}><th scope="row" className="text-left font-normal">{c}</th>{visibles.map((s) => <td key={s.nombre} className="text-right">{formato(s.valores[j] ?? 0)}</td>)}</tr>)}</tbody>
        </table>
      </details>
    </figure>
  )
}
