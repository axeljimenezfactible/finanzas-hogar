// Gráfica base (DESIGN §7.7 y §8.6, K27) con Apache ECharts 6 (renderizado SVG, módulos sueltos).
// Paleta categórica de C2 (--ui-chart-1..8) leída del tema activo; cada serie con marcador, trazo y patrón (decal).
// Accesibilidad (§7.7, U5-07): la descripción accesible es la propia «título. resumen» en español (aria.label.description),
// nunca la automática de ECharts (en inglés y con centavos crudos). La descripción emergente y la guía vertical usan
// tokens de C2 releídos al cambiar el tema; las series no se recolorean ni se atenúan al pasar el puntero (emphasis off).
// Dinero (§7.7 y §8.5): eje Y compacto («$5 mil») solo cuando la gráfica es de dinero; etiquetas directas con el formato
// completo y selectivas (líneas: último valor y extremos; área: valor final; columnas: el valor con ≤ 4 series, y si no
// cabe pasa a «Ver como tabla»); todo el texto del SVG hereda tabular-nums del contenedor.
import { useEffect, useRef, useState } from 'react'
import { ChevronRight } from 'lucide-react'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { AriaComponent, GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'
import { textos } from '../textos/es-MX'
import { dinero, dineroCompacto } from '../api/formatos'
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
/** Índices con etiqueta directa en una línea (DESIGN §8.5): el último valor y los extremos, nunca uno en cada punto. */
export function selectivos(v: (number | null)[]) {
  const idx = v.map((x, i) => [x, i] as const).filter(([x]) => x !== null) as [number, number][]
  if (!idx.length) return new Set<number>()
  const max = idx.reduce((a, b) => (b[0] > a[0] ? b : a)), min = idx.reduce((a, b) => (b[0] < a[0] ? b : a))
  return new Set([idx[idx.length - 1][1], max[1], min[1]])
}
/** Columnas con etiqueta directa solo hasta 4 series (§8.5: «si no cabe, en la tabla»). */
const MAX_SERIES_CON_ETIQUETA = 4
const DINERO = (v: number) => dinero(v)
export interface Serie { nombre: string; valores: number[] }
interface Props {
  titulo: string; resumen: string; tipo: 'barras' | 'lineas' | 'area'; categorias: string[]; series: Serie[]
  /** Formato completo de etiquetas, descripción emergente y tabla. Sin él, la gráfica es de dinero (`dinero`). */
  formato?: (v: number) => string
  /** Formato del eje Y. Por omisión: `dineroCompacto` si la gráfica es de dinero; si no, el mismo `formato` (kilos, kcal). */
  formatoEje?: (v: number) => string
}
export function Grafica({ titulo, resumen, tipo, categorias, series, formato: formatoDado, formatoEje: ejeDado }: Props) {
  const ref = useRef<HTMLDivElement>(null)
  const tema = usarTemaVivo()
  const visibles = series.slice(0, 8) // 8 series máximo; la novena va en «Otros» (la arma quien llama)
  const formato = formatoDado ?? DINERO
  const formatoEje = ejeDado ?? formatoDado ?? dineroCompacto
  useEffect(() => {
    if (!ref.current || !visibles.length) return
    const graf = echarts.init(ref.current, null, { renderer: 'svg' })
    const sup = css('--ui-surface-card'), sec = css('--ui-text-secondary')
    const columnasConEtiqueta = tipo === 'barras' && visibles.length <= MAX_SERIES_CON_ETIQUETA
    const conEtiqueta = (valores: number[]): Set<number> =>
      tipo === 'lineas' ? selectivos(valores) : tipo === 'area' ? new Set(valores.length ? [valores.length - 1] : []) : new Set()
    graf.setOption({
      aria: { enabled: true, label: { enabled: true, description: `${titulo}. ${resumen}` }, decal: { show: true } },
      textStyle: { fontFamily: 'Inter, ui-sans-serif, system-ui', color: sec },
      grid: { left: 8, right: 8, top: 24, bottom: 8, containLabel: true },
      legend: visibles.length > 1 ? { bottom: 0, textStyle: { color: css('--ui-text-default') }, inactiveColor: css('--ui-text-disabled') } : undefined,
      tooltip: {
        trigger: 'axis',
        className: 'fh-descripcion-grafica',
        valueFormatter: (v: unknown) => formato(Number(v)),
        backgroundColor: css('--ui-surface-raised'),
        borderColor: css('--ui-border-divider'),
        borderWidth: 1,
        padding: 8,
        textStyle: { color: css('--ui-text-default'), fontFamily: 'Inter, ui-sans-serif, system-ui', fontSize: 14 },
        extraCssText: `border-radius: ${css('--radius-md')}; box-shadow: var(--ui-shadow-1); font-variant-numeric: tabular-nums;`,
        axisPointer: { type: 'line', lineStyle: { color: css('--ui-viz-baseline'), width: 1, type: 'solid' } },
      },
      xAxis: { type: 'category', data: categorias, axisLine: { lineStyle: { color: css('--ui-viz-baseline') } }, axisLabel: { color: sec } },
      yAxis: { type: 'value', splitLine: { lineStyle: { color: css('--ui-viz-grid') } }, axisLabel: { color: sec, formatter: (v: number) => formatoEje(v) } },
      series: visibles.map((s, i) => {
        const marcados = conEtiqueta(s.valores)
        return {
          name: s.nombre,
          type: tipo === 'barras' ? 'bar' : 'line',
          data: s.valores.map((v, j) => (marcados.has(j) ? { value: v, label: { show: true } } : v)),
          label: { show: columnasConEtiqueta, position: 'top', color: css('--ui-text-default'), fontSize: 12, formatter: (p: { value: unknown }) => formato(Number(p.value)) },
          labelLayout: { hideOverlap: true },
          animation: false,
          color: css(`--ui-chart-${i + 1}`),
          symbol: MARCADORES[i],
          symbolSize: 8,
          lineStyle: { type: TRAZOS[i], width: 2 },
          areaStyle: tipo === 'area' ? { opacity: 0.2 } : undefined,
          stack: tipo === 'area' ? 'total' : undefined,
          itemStyle: { decal: patron(i, sup) },
          emphasis: { disabled: true },
        }
      }),
    })
    const ro = new ResizeObserver(() => graf.resize())
    ro.observe(ref.current)
    return () => { ro.disconnect(); graf.dispose() }
  }, [tema, tipo, categorias, visibles, formato, formatoEje, titulo, resumen])
  if (!visibles.length) return <Vacio texto={textos.vacio.grafica} />
  return (
    <figure className="flex flex-col gap-2 rounded-lg border border-border-divider bg-surface-card p-4">
      <figcaption>
        <span className="block text-small font-medium text-text-secondary">{titulo}</span>
        <span className="block text-small text-text-default">{resumen}</span>
      </figcaption>
      <div ref={ref} tabIndex={0} role="img" aria-label={`${titulo}. ${resumen}`} className="h-50 w-full tabular-nums tablet:h-70" />
      <details className="group text-small">
        <summary className="inline-flex min-h-9 cursor-pointer list-none items-center gap-1 text-text-link coarse:min-h-11 [&::-webkit-details-marker]:hidden">
          <ChevronRight aria-hidden className="size-5 group-open:rotate-90" />{textos.grafica.tablaAlterna}
        </summary>
        <table className="mt-2 w-full tabular-nums">
          <thead><tr><th className="text-left" scope="col">{titulo}</th>{visibles.map((s) => <th key={s.nombre} scope="col" className="text-right">{s.nombre}</th>)}</tr></thead>
          <tbody>{categorias.map((c, j) => <tr key={c}><th scope="row" className="text-left font-normal">{c}</th>{visibles.map((s) => <td key={s.nombre} className="text-right">{formato(s.valores[j] ?? 0)}</td>)}</tr>)}</tbody>
        </table>
      </details>
    </figure>
  )
}
