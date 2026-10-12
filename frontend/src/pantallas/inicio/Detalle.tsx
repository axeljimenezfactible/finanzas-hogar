// Detalle de P-01 (U3 §3): D-1 a D-7, cada panel con su petición, sus estados y los filtros que admite (§6).
import { CircleCheck } from 'lucide-react'
import { Enlace, conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia, porcentaje } from '../../api/formatos'
import { MarcaPersona, MarcaTienda } from '../../componentes/Chip'
import { ChipSeleccionable } from '../../componentes/Chip'
import { EnlaceAccion, EsqueletoFilas, Estados, Filtros, IconoAviso, IconoInfo, IconoOk, Insignia, InsigniaPago, Panel, codigoTienda, decimas, diasAviso, entero, estadoPago, karen, listaDe, nombreCredito, titulares, useInicio, type GrupoFiltro } from './comun'
import { mesDe, semanaISO, semanaVigente, sumarDias } from './fechas'
import { textosInicio as t } from './textos'
import type { Esquemas } from '../../api/tipos'
type Pago = Esquemas['ScheduledPayment']
import Karen from './Karen'
import type { PropsPantalla } from '../../app/registro'

const F = t.detalle
type Clave = 'periodo' | 'persona' | 'familia' | 'credito' | 'categoria' | 'tienda'
const NOMBRES: Record<Clave, string> = { periodo: F.periodo, persona: F.persona, familia: F.familia, credito: F.credito, categoria: F.categoria, tienda: F.tienda }
const lista = (v: string | null) => (v ? v.split(',').filter(Boolean) : [])
const nums = (v: string[]) => (v.length ? v.map(Number) : undefined)
const progreso = 'h-2 w-full appearance-none overflow-hidden rounded-full bg-viz-track [&::-webkit-progress-bar]:bg-viz-track [&::-webkit-progress-value]:bg-viz-meter [&::-moz-progress-bar]:bg-viz-meter'

export default function Detalle({ consulta }: Pick<PropsPantalla, 'consulta'>) {
  const c = useInicio(), hoy = c.hoy
  const f = { periodo: consulta.get('periodo') ?? 'mes', persona: lista(consulta.get('persona')), familia: lista(consulta.get('familia')), credito: lista(consulta.get('credito')), categoria: lista(consulta.get('categoria')), tienda: lista(consulta.get('tienda')) }
  const period = f.periodo === 'mes' ? mesDe(hoy) : `${mesDe(hoy)}-${f.periodo}`
  const activos = (Object.keys(NOMBRES) as Clave[]).filter((k) => (k === 'periodo' ? f.periodo !== 'mes' : f[k].length > 0))
  const noUsa = (usa: Clave[]) => activos.filter((k) => !usa.includes(k)).map((k) => F.noUsa(NOMBRES[k]))
  const poner = (k: string, v: string[]) => navegar(conConsulta({ [k]: v.length ? v.join(',') : null }), { reemplazar: true })
  const grupos: GrupoFiltro[] = [
    { id: 'persona', nombre: F.persona, opciones: listaDe(c.personas).map((p) => ({ valor: String(p.id), texto: p.display_name, marca: <MarcaPersona codigo={p.code} /> })) },
    { id: 'familia', nombre: F.familia, opciones: (['household', 'external', 'joint'] as const).map((v) => ({ valor: v, texto: F.familias[v] })) },
    { id: 'credito', nombre: F.credito, opciones: listaDe(c.creditos).map((x) => ({ valor: String(x.id), texto: x.name })) },
    { id: 'categoria', nombre: F.categoria, opciones: listaDe(c.categorias).map((x) => ({ valor: String(x.id), texto: x.name })) },
    { id: 'tienda', nombre: F.tienda, opciones: listaDe(c.tiendas).map((x) => ({ valor: String(x.id), texto: x.name, marca: <MarcaTienda codigo={codigoTienda(x)} /> })) },
  ]
  // PROVISIONAL (bloqueo WF3-inicio-vacios): P-02 recibe los filtros con los nombres de C1.
  const aGastos = Object.entries({ period, person_id: f.persona.join(','), family: f.familia.join(','), category_id: f.categoria.join(','), store_id: f.tienda.join(','), credit_id: f.credito.join(',') }).filter(([, v]) => v).map(([k, v]) => `${k}=${encodeURIComponent(v)}`).join('&')
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2" role="group" aria-label={F.periodo}>
        {(['mes', 'q1', 'q2'] as const).map((p) => <ChipSeleccionable key={p} seleccionado={f.periodo === p} onClick={() => navegar(conConsulta({ periodo: p === 'mes' ? null : p }), { reemplazar: true })}>{F[p]}</ChipSeleccionable>)}
        <Filtros grupos={grupos} valores={{ persona: f.persona, familia: f.familia, credito: f.credito, categoria: f.categoria, tienda: f.tienda }} alCambiar={poner} />
      </div>
      <div className="grid grid-cols-1 gap-4 desktop:grid-cols-2 desktop:gap-6">
        <D1 period={period} categorias={f.categoria} aviso={noUsa(['periodo', 'categoria'])} />
        <D2 aviso={noUsa([])} />
        <D3 persona={f.persona} familia={f.familia} credito={f.credito} aviso={noUsa(['persona', 'familia', 'credito'])} />
        <D4 persona={f.persona} aviso={noUsa(['persona'])} />
        <D5 tiendas={f.tienda} aviso={noUsa(['tienda'])} />
        <D6 aviso={noUsa([])} />
        <Karen aviso={noUsa([])} />
      </div>
      <Enlace a={`${rutaDe('P-02')}?${aGastos}`} className="inline-flex min-h-11 items-center text-small font-medium text-text-link underline underline-offset-2">{F.verGastos}</Enlace>
    </div>
  )
}
function Avisos({ a }: { a: string[] }) { return <>{a.map((x) => <p key={x} className="text-caption text-text-secondary">{x}</p>)}</> }
function D1({ period, categorias, aviso }: { period: string; categorias: string[]; aviso: string[] }) {
  const c = useInicio()
    const q = usarConsulta('get', '/api/v1/budget/categories', { consulta: { period } })
    const nombre = (id: number) => listaDe(c.categorias).find((x) => x.id === id)?.name ?? t.comun.sinDato
    return <Panel etiqueta={t.d1.titulo} titulo={t.d1.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d1.falla} esqueleto={<EsqueletoFilas n={5} />}>{(d) => {
        const filas = d.rows.filter((r) => r.planned_cents !== undefined && (!categorias.length || categorias.includes(String(r.category_id)))).sort((a, b) => (b.used_pct ?? -1) - (a.used_pct ?? -1)).slice(0, 5)
        if (!filas.length) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d1.vacio}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '02' })} texto={t.d1.accionVacio} /></div>
        return <ul className="flex flex-col gap-3">{filas.map((r) => <li key={r.category_id} className="flex flex-col gap-1 text-small">
          <span className="flex flex-wrap items-center justify-between gap-2"><span className="font-medium">{nombre(r.category_id)}</span>{r.status === 'warning' ? <Insignia rol="advertencia" icono={<IconoAviso />}>{t.t1.vasEn(r.used_pct ?? 0)}</Insignia> : r.status === 'exceeded' ? <Insignia rol="advertencia" icono={<IconoAviso />}>{t.t1.rebasaste}</Insignia> : null}</span>
          <progress className={progreso} max={Math.max(r.planned_cents ?? 0, r.actual_cents, 1)} value={r.actual_cents} aria-label={nombre(r.category_id)} />
          <span className="tabular-nums text-text-secondary">{t.d1.fila(dinero(r.actual_cents), dinero(r.planned_cents ?? 0), r.used_pct !== undefined ? porcentaje(r.used_pct) : t.comun.sinDato)}</span>
        </li>)}</ul>
      }}</Estados><EnlaceAccion a={rutaDe('P-06')} texto={t.acciones.verTodo} /></Panel>
  }
function D2({ aviso }: { aviso: string[] }) {
  const hoy = useInicio().hoy
    const q = usarConsulta('get', '/api/v1/cash-flow/fortnights')
    return <Panel etiqueta={t.d2.titulo} titulo={t.d2.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d2.falla} esqueleto={<EsqueletoFilas n={7} />}>{(d) => {
        const r = d.rows.find((x) => x.start_date <= hoy && hoy <= x.end_date)
        if (!r || r.income_cents === 0) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d2.vacio}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={t.d2.accionVacio} /></div>
        const filas: [string, string][] = [[t.d2.ingresos, dinero(r.income_cents, 'ingreso') + (r.income_is_expected ? ` · ${t.d2.esperado}` : '')], [t.d2.pagosCredito, dinero(r.credit_payments_cents, 'gasto')], [t.d2.fijos, dinero(r.fixed_and_upcoming_cents, 'gasto')], [t.d2.variable, dinero(r.variable_estimate_cents, 'gasto')], [t.d2.ahorro, typeof r.planned_savings_cents === 'number' ? dinero(r.planned_savings_cents) : t.comun.sinDato], [t.d2.saldo, dinero(r.balance_cents)], [t.d2.acumulado, dinero(r.cumulative_balance_cents)]]
        return <dl className="grid grid-cols-2 gap-x-3 gap-y-2 text-small">{filas.map(([k, v]) => <div key={k} className="contents"><dt className="text-text-secondary">{k}</dt><dd className="text-right tabular-nums">{v}</dd></div>)}</dl>
      }}</Estados><EnlaceAccion a={rutaDe('P-07')} texto={t.acciones.verTodo} /></Panel>
  }
function D3({ persona, familia, credito, aviso }: { persona: string[]; familia: string[]; credito: string[]; aviso: string[] }) {
  const c = useInicio(), hoy = c.hoy
    const np = c.proximo.datos?.next_payment
    const filtros = { person_id: nums(persona), family: familia[0] as 'household' | 'external' | 'joint' | undefined, credit_id: nums(credito), sort: 'due_date', limit: 5 }
    const vencidos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { ...filtros, status: 'overdue', period: `${np?.due_date ?? hoy}..${sumarDias(hoy, -1)}` } }, np?.is_overdue === true)
    const proximos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { ...filtros, status: 'pending', period: `${hoy}..${sumarDias(hoy, 60)}` } })
    const aviso2 = diasAviso(c), idK = karen(c)?.id
    const av = c.proximo.datos?.availability ?? []
    const Fila = ({ p }: { p: Pago }) => {
      const k = p.shares?.find((s) => s.person_id === idK), h = (p.shares ?? []).filter((s) => s.person_id !== idK).reduce((a, s) => a + s.amount_cents, 0)
      return <li className="flex flex-col gap-1 border-t border-border-divider py-2 text-small tabular-nums">
        <span className="flex flex-wrap items-center justify-between gap-2"><span className="font-medium">{fechaConDia(p.due_date)} · {nombreCredito(c, p.credit_id)}</span><span>{dinero(p.amount_cents)}</span></span>
        <span className="flex flex-wrap items-center justify-between gap-2 text-text-secondary"><span>{t.comun.hogarYKaren(dinero(h), k && k.amount_cents > 0 ? dinero(k.amount_cents) : undefined)}</span><InsigniaPago e={estadoPago(p, hoy, aviso2, idK)} /></span>
      </li>
    }
    return <Panel etiqueta={t.d3.titulo} titulo={t.d3.titulo}><Avisos a={aviso} />
      {np?.is_overdue && <Estados consulta={vencidos} falla={t.d3.falla} esqueleto={<EsqueletoFilas n={2} />}>{(d) => <div><h3 className="text-small font-medium">{t.d3.vencidos}</h3><ul>{(d.items as Pago[]).map((p) => <Fila key={p.id} p={p} />)}</ul>{d.next_cursor && <p className="text-small">{t.d3.masVencidos}</p>}</div>}</Estados>}
      <Estados consulta={proximos} falla={t.d3.falla} esqueleto={<EsqueletoFilas n={5} />}>{(d) => {
        if (!d.items.length && !np?.is_overdue) return listaDe(c.creditos).length || av.length ? <p className="text-body text-text-secondary">{t.d3.sinPagos}</p> : <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d3.sinCreditos}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '05' })} texto={t.d3.accionVacio} /></div>
        return <div><h3 className="text-small font-medium">{t.d3.proximos}</h3><ul>{(d.items as Pago[]).map((p) => <Fila key={p.id} p={p} />)}</ul></div>
      }}</Estados>
      {av.length > 0 && <div className="flex flex-col gap-2"><h3 className="text-small font-medium">{t.d3.disponible}</h3>{av.map((a) => <Disponible key={a.credit_id} id={a.credit_id} />)}</div>}
      <EnlaceAccion a={rutaDe('P-17')} texto={t.acciones.verTodo} /></Panel>
  }
function Disponible({ id }: { id: number }) {
  const c = useInicio()
    const q = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-preview', { ruta: { credit_id: id } })
    return <Estados consulta={q} falla={t.d3.fallaDisponible} esqueleto={<EsqueletoFilas n={1} />}>{(d) => <div className="flex flex-col gap-1 text-small tabular-nums">
      <span className="font-medium">{nombreCredito(c, id)}</span>
      <span className="grid grid-cols-1 gap-2 phone:grid-cols-2">{(['A', 'B'] as const).map((l) => { const r = l === 'A' ? d.reading_a : d.reading_b
        return <span key={l} className="flex flex-col gap-1 rounded-sm border border-border-divider p-2"><span className="flex flex-wrap items-center gap-2">{t.comun.lectura(l)} {dinero(r.available_today_cents)}{d.availability_reading === l && <Insignia rol="info" icono={<CircleCheck aria-hidden />}>{t.comun.activa}</Insignia>}</span><span className="text-text-secondary">{t.d3.minimo(dinero(r.horizon_min_cents))}</span></span> })}</span>
      {!d.availability_confirmed && <Insignia rol="info" icono={<IconoInfo />}>{t.t3.reposicion(nombreCredito(c, id))}</Insignia>}
    </div>}</Estados>
  }
function D4({ persona, aviso }: { persona: string[]; aviso: string[] }) {
  const c = useInicio(), hoy = c.hoy
    const q = usarConsulta('get', '/api/v1/meal-plan/week', { consulta: { period: semanaISO(hoy), person_id: nums(persona) } })
    const tiempos = usarConsulta('get', '/api/v1/meal-slots', { consulta: { limit: 200 } })
    const deHoy = (q.datos?.cells ?? []).filter((x) => x.plan_date === hoy && (x.recipe_id !== undefined || x.is_eating_out))
    const recetas = usarConsulta('get', '/api/v1/recipes', { consulta: { limit: 200 } }, deHoy.some((x) => x.recipe_id !== undefined))
    const pers = titulares(c)
    return <Panel etiqueta={t.d4.titulo} titulo={t.d4.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d4.falla} esqueleto={<EsqueletoFilas n={4} />}>{() => {
        if (!deHoy.length) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d4.vacio}</p><EnlaceAccion a={rutaDe('P-10')} texto={t.d4.accionVacio} /></div>
        const orden = listaDe(tiempos) as Esquemas['MealSlot'][]
        const filas = [...deHoy].sort((a, b) => (orden.find((s) => s.id === a.meal_slot_id)?.sort_order ?? 9) - (orden.find((s) => s.id === b.meal_slot_id)?.sort_order ?? 9) || pers.findIndex((p) => p.id === a.person_id) - pers.findIndex((p) => p.id === b.person_id))
        return <ul>{filas.map((x) => { const p = pers.find((y) => y.id === x.person_id)
          return <li key={`${x.meal_slot_id}-${x.person_id}`} className="flex flex-wrap items-center justify-between gap-2 border-t border-border-divider py-2 text-small tabular-nums">
            <span className="flex items-center gap-2"><span className="font-medium">{orden.find((s) => s.id === x.meal_slot_id)?.name ?? t.comun.sinDato}</span>{p && <><MarcaPersona codigo={p.code} />{p.display_name}</>}</span>
            <span>{x.is_eating_out ? t.d4.fuera : (listaDe(recetas) as { id: number; name: string }[]).find((r) => r.id === x.recipe_id)?.name ?? t.comun.sinDato}</span>
            <span className="text-text-secondary">{x.kcal !== undefined ? t.d4.kcal(entero(x.kcal)) : t.comun.sinDato} · {x.cost_cents !== undefined ? dinero(x.cost_cents) : t.comun.sinDato}</span></li> })}</ul>
      }}</Estados><EnlaceAccion a={rutaDe('P-10')} texto={t.acciones.verTodo} /></Panel>
  }
function D5({ tiendas, aviso }: { tiendas: string[]; aviso: string[] }) {
  const c = useInicio(), hoy = c.hoy
    const week = semanaVigente(hoy)
    const canasta = usarConsulta('get', '/api/v1/dashboard/best-basket', { consulta: { week } })
    const b = canasta.datos
    const q = usarConsulta('get', '/api/v1/basket-comparison', { consulta: { week, basket: b?.basket, expense_id: b?.expense_id } }, b !== undefined && b.store_id !== undefined)
    const ts = listaDe(c.tiendas)
    const vacio = <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d5.vacio}</p><EnlaceAccion a={rutaDe('P-15')} texto={t.d5.accionVacio} /></div>
    return <Panel etiqueta={t.d5.titulo} titulo={t.d5.titulo}><Avisos a={aviso} />
      <Estados consulta={canasta} falla={t.d5.falla} esqueleto={<EsqueletoFilas n={5} />}>{(bb) => bb.store_id === undefined ? vacio :
        <Estados consulta={q} falla={t.d5.falla} esqueleto={<EsqueletoFilas n={5} />}>{(d) => {
          const ver = (id?: number) => !tiendas.length || tiendas.includes(String(id))
          const filas = d.items.filter((x) => x.difference_cents !== undefined && ver(x.winner_store_id)).sort((a, b2) => (b2.difference_cents ?? 0) - (a.difference_cents ?? 0)).slice(0, 5)
          const tienda = (id?: number) => ts.find((s) => s.id === id)
          return <div className="flex flex-col gap-3">
            <ul>{filas.map((x) => { const g = tienda(x.winner_store_id); const viejo = x.prices.some((p) => p.store_id === x.winner_store_id && p.is_stale)
              return <li key={x.product_id} className="flex flex-wrap items-center justify-between gap-2 border-t border-border-divider py-2 text-small tabular-nums">
                <Producto id={x.product_id} /><span className="flex items-center gap-1"><MarcaTienda codigo={codigoTienda(g)} />{g?.name ?? t.comun.sinDato}</span><span>{dinero(x.difference_cents ?? 0)}</span>
                {viejo && <Insignia rol="advertencia" icono={<IconoAviso />}>{t.d5.viejo}</Insignia>}</li> })}</ul>
            <div><h3 className="text-small font-medium">{t.d5.totales}</h3><ul>{[...d.stores].filter((s) => ver(s.store_id)).sort((a, b2) => orden(tienda(a.store_id)) - orden(tienda(b2.store_id))).map((s) => <li key={s.store_id} className="flex items-center justify-between gap-2 py-1 text-small tabular-nums"><span className="flex items-center gap-1"><MarcaTienda codigo={codigoTienda(tienda(s.store_id))} />{tienda(s.store_id)?.name ?? t.comun.sinDato}</span><span>{dinero(s.total_cents)}</span></li>)}</ul></div>
          </div>
        }}</Estados>}</Estados><EnlaceAccion a={rutaDe('P-14')} texto={t.acciones.verTodo} /></Panel>
  }
function D6({ aviso }: { aviso: string[] }) {
    const q = usarConsulta('get', '/api/v1/savings-goals', { consulta: { status: 'active', limit: 50 } })
    return <Panel etiqueta={t.d6.titulo} titulo={t.d6.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d6.falla} esqueleto={<EsqueletoFilas n={3} />}>{(d) => {
        if (!d.items.length) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d6.vacio}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '13' })} texto={t.d6.accionVacio} /></div>
        return <ul className="flex flex-col gap-3">{(d.items as Esquemas['SavingsGoal'][]).map((g) => { const est = g.status ? t.t6.estados[g.status] : undefined
          return <li key={g.id} className="flex flex-col gap-1 text-small tabular-nums">
            <span className="flex flex-wrap items-center justify-between gap-2"><span className="font-medium">{g.name}</span>{est && <Insignia rol={g.status === 'achieved' ? 'exito' : g.status === 'overdue' ? 'advertencia' : 'info'} icono={g.status === 'achieved' ? <IconoOk /> : g.status === 'overdue' ? <IconoAviso /> : <IconoInfo />}>{est}</Insignia>}</span>
            {g.target_cents !== undefined && <progress className={progreso + ' [&::-webkit-progress-value]:bg-role-savings! [&::-moz-progress-bar]:bg-role-savings!'} max={g.target_cents} value={Math.min(g.saved_cents ?? 0, g.target_cents)} aria-label={g.name} />}
            <span className="text-text-secondary">{t.d6.avance(g.progress_tenths_pct !== undefined ? decimas(g.progress_tenths_pct) : t.comun.sinDato)} · {t.d6.deDe(dinero(g.saved_cents ?? 0), g.target_cents !== undefined ? dinero(g.target_cents) : t.comun.sinDato)}</span>
            <span className="text-text-secondary">{g.required_monthly_cents !== undefined ? t.d6.aportacion(dinero(g.required_monthly_cents)) : null}{g.target_date ? ` · ${t.d6.fechaMeta(fecha(g.target_date))}` : ''}</span>
          </li> })}</ul>
      }}</Estados><EnlaceAccion a={rutaDe('P-09')} texto={t.acciones.verTodo} /></Panel>
  }
const ORDEN_T = ['walmart', 'bodega-aurrera', 'sams-club', 'market']
const orden = (s?: { name: string; kind: string }) => ORDEN_T.indexOf(codigoTienda(s))
function Producto({ id }: { id: number }) {
  const q = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: id } })
  return <span className="font-medium">{q.datos?.name ?? t.comun.sinDato}</span>
}
