// Tendencias de P-01 (U3 §4): T-1 a T-4 en orden fijo, series en orden fijo con marcador, trazo o patrón y etiqueta directa.
import { useState, type ReactNode } from 'react'
import { CircleCheck } from 'lucide-react'
import { rutaDe } from '../../app/rutas'
import { usarConsulta } from '../../api/consulta'
import { dinero, mes } from '../../api/formatos'
import { ChipSeleccionable, MarcaPersona } from '../../componentes/Chip'
import { patron } from '../../componentes/Grafica'
import { abrirRegistro } from '../../shell/estado'
import { Boton } from '../../componentes/Boton'
import { EnlaceAccion, Estados, TABLET, entero, listaDe, titulares, usarMedia, useInicio } from './comun'
import { compacto, ejes, FiguraGrafica, selectivos, token } from './GraficaInicio'
import { diaDe, mesAbrev, mesDe, rangoSemana, sumarDias } from './fechas'
import { textosInicio as t } from './textos'
const G = t.graficas
const esqueleto = <div aria-hidden className="h-60 w-full rounded-md bg-skeleton-on-card tablet:h-72" />
const etiqueta = (formato: (v: number) => string) => ({ show: true, color: token('--ui-text-default'), fontSize: 12, formatter: (p: { value: number }) => formato(p.value) })
const conEtiquetas = (v: (number | null)[], formato: (v: number) => string, solo?: Set<number>) => { const s = solo ?? selectivos(v); return v.map((x, i) => (x === null ? null : s.has(i) ? { value: x, label: etiqueta(formato) } : x)) }
function VacioGrafica({ titulo, texto, accion }: { titulo: string; texto: ReactNode; accion?: ReactNode }) {
  return <section aria-label={titulo} className="flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6"><span className="text-small font-medium text-text-secondary">{titulo}</span><p className="flex items-center gap-2 text-body text-text-secondary">{texto}</p>{accion}</section>
}
function Chips<T extends string | number>({ etiqueta: e, opciones, valor, alCambiar, todas }: { etiqueta: string; opciones: { valor: T; texto: string; marca?: ReactNode }[]; valor: T[]; alCambiar: (v: T[]) => void; todas: string }) {
  return <div role="group" aria-label={e} className="flex flex-wrap gap-2">
    <ChipSeleccionable seleccionado={!valor.length} onClick={() => alCambiar([])}>{todas}</ChipSeleccionable>
    {opciones.map((o) => <ChipSeleccionable key={o.valor} seleccionado={valor.includes(o.valor)} marca={o.marca} onClick={() => alCambiar(valor.includes(o.valor) ? valor.filter((x) => x !== o.valor) : [...valor, o.valor])}>{o.texto}</ChipSeleccionable>)}
  </div>
}

export default function Tendencias() {
  return <div className="grid grid-cols-1 gap-4 wide:grid-cols-2 wide:gap-6"><T1 /><T2 /><T3 /><T4 /></div>
}
function T1() {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const [cat, setCat] = useState<number[]>([])
  const q = usarConsulta('get', '/api/v1/budget/monthly-trend', { consulta: { category_id: cat.length ? cat : undefined } })
  const cats = listaDe(c.categorias)
  const nombre = cat.length ? cats.filter((x) => cat.includes(x.id)).map((x) => x.name).join(', ') : G.todas
  const filtros = <Chips etiqueta={t.detalle.categoria} todas={G.todas} valor={cat} alCambiar={setCat} opciones={cats.map((x) => ({ valor: x.id, texto: x.name }))} />
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const m = d.months
    if (!m.some((x) => x.planned_cents !== undefined)) return <VacioGrafica titulo={G.t1.titulo} texto={G.t1.vacio} accion={<><EnlaceAccion a={rutaDe('P-26', { seccion: '02' })} texto={G.t1.accionVacio} />{filtros}</>} />
    const u = m[m.length - 1], cs = m.map((x) => mes(x.month)), ult = new Set([m.length - 1])
    const sup = token('--ui-surface-card')
    return <FiguraGrafica id="t1" titulo={G.t1.titulo} subtitulo={G.ultimos12(nombre)} filtros={filtros} dependencias={[d, ancho]}
      resumen={u.planned_cents !== undefined ? G.t1.resumen(mes(u.month), dinero(u.actual_cents), dinero(u.planned_cents)) : G.t1.resumenSinPlan(mes(u.month), dinero(u.actual_cents))}
      tabla={{ columnas: [G.mes, G.t1.real, G.t1.plan], filas: m.map((x) => [mes(x.month), dinero(x.actual_cents), x.planned_cents !== undefined ? dinero(x.planned_cents) : t.comun.sinDato]) }}
      opcion={() => ({ ...ejes(cs, compacto, ancho), labelLayout: { hideOverlap: true }, series: [
        { name: G.t1.real, type: 'bar', barGap: '0%', data: conEtiquetas(m.map((x) => x.actual_cents), (v) => dinero(v), ult), label: { position: 'top' }, itemStyle: { color: token('--ui-chart-1'), borderColor: sup, borderWidth: 1 } },
        { name: G.t1.plan, type: 'bar', data: m.map((x) => x.planned_cents ?? null), itemStyle: { color: token('--ui-chart-2'), decal: patron(1, sup), borderColor: sup, borderWidth: 1 } },
      ] })} />
  }}</Estados>
}
function T2() {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const q = usarConsulta('get', '/api/v1/cash-flow/accrual-vs-cash')
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const m = d.months
    if (!m.length || m.every((x) => x.accrued_cents === 0 && x.cash_cents === 0)) return <VacioGrafica titulo={G.t2.titulo} texto={G.t2.vacio} accion={<Boton variante="enlace" className="w-max" onClick={() => abrirRegistro('gasto')}>{G.t2.accionVacio}</Boton>} />
    const u = m[m.length - 1], cs = m.map((x) => mes(x.month)), sup = token('--ui-surface-card')
    const hoyI = m.findIndex((x) => x.month === mesDe(c.hoy))
    const adelantado = m.filter((x) => x.cumulative_difference_cents < 0).map((x) => mes(x.month))
    const suma = (k: 'accrued_cents' | 'cash_cents' | 'accrued_household_cents' | 'accrued_external_cents' | 'cash_household_cents' | 'cash_external_cents') => dinero(m.reduce((a, x) => a + x[k], 0))
    const C = G.t2.columnas
    return <FiguraGrafica id="t2" titulo={G.t2.titulo} subtitulo={G.t2.subtitulo(mes(m[0].month), mes(u.month))} dependencias={[d, ancho]}
      resumen={G.t2.resumen(mes(u.month), m[6] ? mes(m[6].month) : t.comun.sinDato, dinero(u.cumulative_accrued_cents), dinero(u.cumulative_cash_cents))}
      notas={[d.projection_mode === 'with_repeat_cadence' ? G.t2.conRepetidas : G.t2.sinRepetidas, ...adelantado.map(G.t2.adelantadoEn)]}
      tabla={{ columnas: [G.mes, C.devengadoMes, C.efectivoMes, G.t2.diferencia, G.t2.devengado, G.t2.efectivo, C.devHogar, C.devKaren, C.efeHogar, C.efeKaren],
        filas: m.map((x) => [mes(x.month), dinero(x.accrued_cents), dinero(x.cash_cents), dinero(x.cumulative_difference_cents), dinero(x.cumulative_accrued_cents), dinero(x.cumulative_cash_cents), dinero(x.accrued_household_cents), dinero(x.accrued_external_cents), dinero(x.cash_household_cents), dinero(x.cash_external_cents)]),
        total: [G.total, suma('accrued_cents'), suma('cash_cents'), t.comun.sinDato, t.comun.sinDato, t.comun.sinDato, suma('accrued_household_cents'), suma('accrued_external_cents'), suma('cash_household_cents'), suma('cash_external_cents')] }}
      opcion={() => {
        // Áreas apiladas (DESIGN §8.3 y §8.5): los valores finales van en la leyenda, abajo y fuera del trazo (una etiqueta sobre
        // el relleno se encima con el área y, a la derecha, la leyenda tapaba su marca). La tabla sigue con todos los valores.
        const e = ejes(cs, compacto, ancho)
        const finales: Record<string, string> = { [G.t2.efectivo]: dinero(u.cumulative_cash_cents), [G.t2.diferencia]: dinero(Math.max(0, u.cumulative_difference_cents)) }
        return { ...e, grid: { ...e.grid, right: 16, bottom: ancho ? 40 : 64 }, legend: { bottom: 0, textStyle: e.legend.textStyle, formatter: (n: string) => `${n} ${finales[n] ?? ''}`.trim() }, series: [
        { name: G.t2.efectivo, type: 'line', stack: 'pila', symbol: 'circle', symbolSize: 8, lineStyle: { color: token('--ui-chart-1'), width: 2 }, itemStyle: { color: token('--ui-chart-1') }, areaStyle: { color: token('--ui-chart-1'), opacity: 1 },
          data: m.map((x) => x.cumulative_cash_cents),
          markLine: hoyI >= 0 ? { symbol: 'none', silent: true, lineStyle: { color: token('--ui-viz-baseline'), width: 1, type: 'solid' }, label: { formatter: G.hoy, color: token('--ui-text-secondary') }, data: [{ xAxis: cs[hoyI] }] } : undefined,
          markArea: m.length > 6 ? { silent: true, itemStyle: { color: token('--ui-viz-track'), opacity: 0.35 }, label: { formatter: G.proyeccion, color: token('--ui-text-secondary') }, data: [[{ xAxis: cs[6] }, { xAxis: cs[m.length - 1] }]] } : undefined },
        { name: G.t2.diferencia, type: 'line', stack: 'pila', symbol: 'rect', symbolSize: 8, lineStyle: { color: token('--ui-chart-2'), width: 2, type: [8, 4] }, itemStyle: { color: token('--ui-chart-2') }, areaStyle: { color: token('--ui-chart-2'), opacity: 1, decal: patron(2, sup) },
          data: m.map((x) => Math.max(0, x.cumulative_difference_cents)) },
      ] } }} />
  }}</Estados>
}
function T3() {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const [sel, setSel] = useState<number[]>([])
  const q = usarConsulta('get', '/api/v1/credits/debt-trend', { consulta: { credit_id: sel.length ? sel : undefined } })
  const creditos = listaDe(c.creditos), elegidos = creditos.filter((x) => sel.includes(x.id))
  const filtros = <Chips etiqueta={t.detalle.credito} todas={G.sinFiltro} valor={sel} alCambiar={setSel} opciones={creditos.map((x) => ({ valor: x.id, texto: x.name }))} />
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const m = d.months
    if (!m.length || m.every((x) => x.live_debt_cents === 0)) return <VacioGrafica titulo={G.t3.titulo} texto={<><CircleCheck aria-hidden className="size-5 text-role-success" />{G.t3.vacio}</>} accion={filtros} />
    const u = m[m.length - 1], cs = m.map((x) => mes(x.month))
    const series = elegidos.length ? elegidos.map((cr) => ({ nombre: cr.name, v: m.map((x) => x.credits.find((y) => y.credit_id === cr.id)?.live_debt_cents ?? null) })) : [{ nombre: G.t3.total, v: m.map((x) => x.live_debt_cents) }]
    const MARC = ['circle', 'rect', 'triangle', 'diamond'], TRAZO: (string | number[])[] = ['solid', [8, 4], [2, 4], [8, 4, 2, 4]]
    return <FiguraGrafica id="t3" titulo={G.t3.titulo} subtitulo={G.ultimos12(elegidos.length ? elegidos.map((x) => x.name).join(', ') : G.t3.total)} filtros={filtros} dependencias={[d, ancho, sel.join()]}
      resumen={elegidos.length ? G.t3.resumenFiltro(mes(u.month), series.map((s) => `${s.nombre} ${s.v[s.v.length - 1] !== null ? dinero(s.v[s.v.length - 1]!) : t.comun.sinDato}`).join(' · ')) : G.t3.resumen(mes(m[0].month), dinero(m[0].live_debt_cents), mes(u.month), dinero(u.live_debt_cents))}
      tabla={{ columnas: [G.mes, ...series.map((s) => s.nombre), G.t3.interesesPagados, G.t3.interesesPorPagar], filas: m.map((x, i) => [mes(x.month), ...series.map((s) => (s.v[i] !== null ? dinero(s.v[i]!) : t.comun.sinDato)), dinero(x.interest_paid_cents), dinero(x.interest_pending_cents)]) }}
      opcion={() => ({ ...ejes(cs, compacto, ancho), series: series.slice(0, 8).map((s, i) => ({ name: s.nombre, type: 'line', symbol: MARC[i % 4], symbolSize: 8, lineStyle: { color: token(`--ui-chart-${i + 1}`), width: 2, type: TRAZO[i % 4] }, itemStyle: { color: token(`--ui-chart-${i + 1}`) }, label: { position: 'top' }, data: conEtiquetas(s.v, (v) => dinero(v)) })) })} />
  }}</Estados>
}
function T4() {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const [sel, setSel] = useState<number[]>([])
  const q = usarConsulta('get', '/api/v1/meal-plan/trends', { consulta: { person_id: sel.length ? sel : undefined } })
  const pers = titulares(c), vistas = sel.length ? pers.filter((p) => sel.includes(p.id)) : pers
  const filtros = <Chips etiqueta={t.detalle.persona} todas={G.sinFiltro} valor={sel} alCambiar={setSel} opciones={pers.map((p) => ({ valor: p.id, texto: p.display_name, marca: <MarcaPersona codigo={p.code} /> }))} />
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const w = d.weeks
    if (!w.length) return <VacioGrafica titulo={G.t4.tituloA} texto={G.t4.vacio} accion={<EnlaceAccion a={rutaDe('P-10')} texto={G.t4.accionVacio} />} />
    const u = w[w.length - 1], cs = w.map((x) => `${diaDe(x.week_start)} ${mesAbrev(x.week_start)}`), rango = rangoSemana(u.week_start, sumarDias(u.week_start, 6))
    const kcal = (p: number, k: 'avg_kcal_per_day' | 'daily_kcal_target') => w.map((x) => x.persons.find((y) => y.person_id === p)?.[k] ?? null)
    const fmtK = (v: number) => entero(v)
    const ultimo = (p: number, k: 'avg_kcal_per_day' | 'daily_kcal_target') => { const v = u.persons.find((y) => y.person_id === p)?.[k]; return v !== undefined ? entero(v) : t.comun.sinDato }
    return <div className="flex flex-col gap-4">
      <FiguraGrafica id="t4a" titulo={G.t4.tituloA} subtitulo={G.t4.subtitulo} filtros={filtros} dependencias={[d, ancho]} resumen={G.t4.resumenA(rango, dinero(u.cost_cents))}
        tabla={{ columnas: [G.semana, G.t4.costo], filas: w.map((x, i) => [cs[i], dinero(x.cost_cents)]) }}
        opcion={() => ({ ...ejes(cs, compacto, ancho), series: [{ name: G.t4.costo, type: 'line', symbol: 'circle', symbolSize: 8, lineStyle: { color: token('--ui-chart-1'), width: 2 }, itemStyle: { color: token('--ui-chart-1') }, label: { position: 'top' }, data: conEtiquetas(w.map((x) => x.cost_cents), (v) => dinero(v)) }] })} />
      <FiguraGrafica id="t4b" titulo={G.t4.tituloB} subtitulo={G.t4.subtitulo} dependencias={[d, ancho, sel.join()]} resumen={G.t4.resumenB(rango, vistas.map((p) => G.t4.resumenPersona(p.display_name, ultimo(p.id, 'avg_kcal_per_day'), ultimo(p.id, 'daily_kcal_target'))).join('; '))}
        tabla={{ columnas: [G.semana, ...vistas.flatMap((p) => [p.display_name, G.t4.metaDe(p.display_name)])], filas: w.map((x, i) => [cs[i], ...vistas.flatMap((p) => { const y = x.persons.find((z) => z.person_id === p.id); return [y?.avg_kcal_per_day !== undefined ? entero(y.avg_kcal_per_day) : t.comun.sinDato, y?.daily_kcal_target !== undefined ? entero(y.daily_kcal_target) : t.comun.sinDato] })]) }}
        opcion={() => ({ ...ejes(cs, fmtK, ancho), series: vistas.flatMap((p, i) => { const color = token(`--ui-person-${p.code}`) || token(`--ui-chart-${i + 1}`)
          return [
            { name: p.display_name, type: 'line', symbol: i ? 'rect' : 'circle', symbolSize: 8, lineStyle: { color, width: 2, type: i ? [8, 4] : 'solid' }, itemStyle: { color }, endLabel: { show: true, formatter: p.display_name, color: token('--ui-text-default') }, label: { position: 'top' }, data: conEtiquetas(kcal(p.id, 'avg_kcal_per_day'), fmtK) },
            { name: G.t4.metaDe(p.display_name), type: 'line', symbol: 'none', lineStyle: { color, width: 1, type: [2, 4] }, itemStyle: { color }, endLabel: { show: true, formatter: G.t4.metaDe(p.display_name), color: token('--ui-text-secondary') }, data: kcal(p.id, 'daily_kcal_target') },
          ] }) })} />
    </div>
  }}</Estados>
}
