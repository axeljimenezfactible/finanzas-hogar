// Calendario de P-01 (U3 §5): pagos del calendario unificado, ingresos y eventos de despensa; quincenas K8; semana de lunes a domingo.
// < 768 px agenda de hoy a hoy + 30 (o + 60); ≥ 768 px mes con franjas de quincena; ≥ 1440 px con el panel del día seleccionado.
import { useState } from 'react'
import { ArrowDownLeft, ArrowUpRight, Calendar, CalendarClock, ChevronLeft, ChevronRight, CreditCard, FileChartColumn, PiggyBank, ShoppingCart, Store } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia } from '../../api/formatos'
import { Boton } from '../../componentes/Boton'
import { cn } from '../../componentes/utils'
import type { Esquemas } from '../../api/tipos'
import { AMPLIO, EnlaceAccion, EsqueletoFilas, Estados, Filtros, IconoInfo, Insignia, InsigniaPago, TABLET, diasAviso, estadoPago, karen, listaDe, nombreCredito, usarMedia, useInicio } from './comun'
import { diaDe, mesDe, quincenaDe, rejillaMes, sumarDias, sumarMeses, ultimoDia } from './fechas'
import { MESES_LARGOS, textosInicio as t } from './textos'
type Evento = Esquemas['CalendarEvent']
type Tipo = Exclude<Evento['event_type'], 'meal_plan' | 'savings_movement'>
const K = t.calendario
const TIPOS: Tipo[] = ['income', 'credit_payment', 'fixed_expense', 'upcoming_expense', 'card_statement', 'planned_grocery', 'market_visit', 'promotion', 'savings_goal', 'monthly_report']
const ICONO: Record<Tipo, typeof Calendar> = { income: ArrowDownLeft, credit_payment: CreditCard, fixed_expense: ArrowUpRight, upcoming_expense: ArrowUpRight, card_statement: CalendarClock, planned_grocery: ShoppingCart, market_visit: Store, promotion: Calendar, savings_goal: PiggyBank, monthly_report: FileChartColumn }
const COLOR: Partial<Record<Tipo, string>> = { income: 'text-role-income', fixed_expense: 'text-role-expense', upcoming_expense: 'text-role-expense', market_visit: 'text-store-market' }
const nums = (v: string[]) => (v.length ? v.map(Number) : undefined)

export default function Calendario() {
  const c = useInicio(), hoy = c.hoy
  const ancho = usarMedia(TABLET), amplio = usarMedia(AMPLIO)
  const [mesV, setMesV] = useState(mesDe(hoy)), [dias, setDias] = useState(30), [sel, setSel] = useState(hoy)
  const [f, setF] = useState<Record<string, string[]>>({ tipos: [], credito: [], familia: [] })
  const period = ancho ? mesV : `${hoy}..${sumarDias(hoy, dias)}`
  const family = f.familia[0] as 'household' | 'external' | 'joint' | undefined
  const ev = usarConsulta('get', '/api/v1/calendar/events', { consulta: { period, event_type: f.tipos.length ? (f.tipos as Tipo[]) : TIPOS, credit_id: nums(f.credito), family, limit: 200 } })
  const pagos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { period, credit_id: nums(f.credito), family, status: 'pending', limit: 200 } })
  const aviso = diasAviso(c), idK = karen(c)?.id
  const porId = new Map(((pagos.datos?.items ?? []) as Esquemas['ScheduledPayment'][]).map((p) => [p.id, p]))
  const grupos = [
    { id: 'tipos', nombre: K.tipos, multiple: true, opciones: TIPOS.map((x) => ({ valor: x, texto: K.tiposEvento[x] })) },
    { id: 'credito', nombre: t.detalle.credito, opciones: listaDe(c.creditos).map((x) => ({ valor: String(x.id), texto: x.name })) },
    { id: 'familia', nombre: t.detalle.familia, opciones: (['household', 'external', 'joint'] as const).map((v) => ({ valor: v, texto: t.detalle.familias[v] })) },
  ]
  function Fila({ e, corto }: { e: Evento; corto?: boolean }) {
    const tipo = e.event_type as Tipo, Icono = ICONO[tipo] ?? Calendar
    const pago = tipo === 'credit_payment' ? estadoPago({ status: e.status === 'paid' || e.status === 'overdue' ? e.status : 'pending', due_date: e.date, shares: e.ref_table === 'scheduled_payment' ? porId.get(e.ref_id)?.shares : undefined }, hoy, aviso, idK) : undefined
    const monto = e.amount_cents === undefined ? '' : tipo === 'income' ? dinero(e.amount_cents, 'ingreso') : tipo === 'fixed_expense' || tipo === 'upcoming_expense' ? dinero(e.amount_cents, 'gasto') : dinero(e.amount_cents)
    const titulo = tipo === 'credit_payment' ? `${nombreCredito(c, e.credit_id)} · ${monto}` : [e.title, monto].filter(Boolean).join(' · ')
    const cuerpo = <>
      <span className={cn('flex items-center gap-2', corto && 'min-w-0')}><Icono aria-hidden className={cn('size-5 shrink-0', COLOR[tipo] ?? 'text-icon-default')} /><span className={cn('tabular-nums', corto ? 'truncate' : 'font-medium', COLOR[tipo])}>{titulo}</span></span>
      {!corto && tipo === 'credit_payment' && e.household_share_cents !== undefined && <span className="text-text-secondary tabular-nums">{t.comun.hogarYKaren(dinero(e.household_share_cents), e.external_share_cents ? dinero(e.external_share_cents) : undefined)}</span>}
      {!corto && tipo === 'promotion' && e.end_date && <span className="text-text-secondary">{K.delAl(fecha(e.date), fecha(e.end_date))}</span>}
      {!corto && pago && <InsigniaPago e={pago} />}
      {!corto && tipo === 'income' && e.status === 'expected' && <Insignia rol="info" icono={<IconoInfo />}>{t.pago.esperado}</Insignia>}
    </>
    const clase = cn('flex flex-col gap-1 text-small', !corto && 'border-t border-border-divider py-2')
    return (tipo === 'credit_payment' || tipo === 'card_statement') && e.credit_id !== undefined && !corto
      ? <li><Enlace a={rutaDe('P-18', { id: e.credit_id })} className={cn(clase, 'min-h-11 rounded-sm text-text-default no-underline hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-border-focus')}>{cuerpo}</Enlace></li>
      : <li className={clase}>{cuerpo}</li>
  }
  const filtros = <Filtros grupos={grupos} valores={f} alCambiar={(id, v) => setF((x) => ({ ...x, [id]: v }))} />
  return (
    <div className="flex flex-col gap-4">
      {filtros}
      <Estados consulta={ev} falla={K.falla} esqueleto={<EsqueletoFilas n={6} />}>{(d) => {
        const items = ([...d.items] as Evento[]).filter((e) => TIPOS.includes(e.event_type as Tipo)).sort((a, b) => a.date.localeCompare(b.date) || TIPOS.indexOf(a.event_type as Tipo) - TIPOS.indexOf(b.event_type as Tipo))
        const delDia = (dia: string) => items.filter((e) => e.date === dia)
        const vacio = !items.length && (listaDe(c.creditos).length === 0
          ? <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{K.sinProgramados}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={K.accionVacio} /></div>
          : <p className="text-body text-text-secondary">{K.vacio}</p>)
        if (!ancho) {
          const fechas = [...new Set(items.map((e) => e.date))]
          return <div data-prueba="agenda" className="flex flex-col gap-3">{vacio}
            {fechas.map((dia, i) => <section key={dia} aria-label={fechaConDia(dia)} className="flex flex-col">
              {(i === 0 || quincenaDe(fechas[i - 1]) !== quincenaDe(dia)) && <p className="mb-2 rounded-sm bg-surface-hover-subtle px-2 py-1 text-caption font-medium text-text-secondary">{diaDe(dia) <= 15 ? K.q1 : K.q2(ultimoDia(mesDe(dia)))}</p>}
              <h3 className="text-small font-semibold">{fechaConDia(dia)}</h3><ul>{delDia(dia).map((e, j) => <Fila key={j} e={e} />)}</ul></section>)}
            {dias < 60 && <Boton variante="secundario" className="w-max" onClick={() => setDias(60)}>{K.mas30}</Boton>}
          </div>
        }
        const u = ultimoDia(mesV), [a, m] = [Number(mesV.slice(0, 4)), Number(mesV.slice(5, 7))]
        return <div className="grid grid-cols-1 gap-4 wide:grid-cols-3">
          <div data-prueba="mes" className="flex flex-col gap-3 wide:col-span-2">
            <div className="flex flex-wrap items-center gap-2">
              <Boton tam="icono" variante="fantasma" aria-label={K.anterior} onClick={() => setMesV(sumarMeses(mesV, -1))}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
              <h2 className="min-w-40 text-title-3 font-semibold">{`${MESES_LARGOS[m - 1]} ${a}`}</h2>
              <Boton tam="icono" variante="fantasma" aria-label={K.siguiente} onClick={() => setMesV(sumarMeses(mesV, 1))}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
              <Boton variante="secundario" onClick={() => { setMesV(mesDe(hoy)); setSel(hoy) }}>{K.hoy}</Boton>
            </div>
            <div className="flex flex-wrap gap-2 text-caption text-text-secondary"><span className="rounded-sm border border-border-divider bg-surface-card px-2 py-1">{K.q1}</span><span className="rounded-sm bg-surface-hover-subtle px-2 py-1">{K.q2(u)}</span></div>
            {vacio}
            <table className="w-full table-fixed border-collapse text-small">
              <thead><tr>{K.semana.map((s) => <th key={s} scope="col" className="py-1 text-caption font-medium text-text-secondary">{s}</th>)}</tr></thead>
              <tbody>{rejillaMes(mesV).map((sem, i) => <tr key={i}>{sem.map((dia, j) => {
                if (!dia) return <td key={j} className="border border-border-divider" />
                const evs = delDia(dia), q2 = diaDe(dia) > 15
                return <td key={j} className={cn('h-24 border border-border-divider p-0 align-top', q2 && 'bg-surface-hover-subtle')}>
                  <button type="button" aria-pressed={sel === dia} aria-current={dia === hoy ? 'date' : undefined} aria-label={`${fechaConDia(dia)} · ${evs.length}`} onClick={() => setSel(dia)} className={cn('flex h-full min-h-11 w-full flex-col gap-1 p-1 text-left hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-border-focus', sel === dia && 'outline-2 outline-border-selected')}>
                    <span className={cn('tabular-nums', dia === hoy && 'font-bold text-text-link')}>{diaDe(dia)}</span>
                    {(diaDe(dia) === 1 || diaDe(dia) === 16) && <span className="truncate text-caption text-text-secondary">{diaDe(dia) === 1 ? t.detalle.q1 : t.detalle.q2}</span>}
                    <ul className="flex min-w-0 flex-col gap-1">{evs.slice(0, 3).map((e, k) => <Fila key={k} e={e} corto />)}</ul>
                    {evs.length > 3 && <span className="text-caption text-text-secondary">{K.masN(evs.length - 3)}</span>}
                  </button></td>
              })}</tr>)}</tbody>
            </table>
          </div>
          <section aria-label={K.diaSeleccionado} className={cn('flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4', !amplio && 'tablet:p-6')}>
            <h3 className="text-small font-semibold">{fechaConDia(sel)}</h3>
            {delDia(sel).length ? <ul>{delDia(sel).map((e, j) => <Fila key={j} e={e} />)}</ul> : <p className="text-small text-text-secondary">{K.sinEventosDia}</p>}
          </section>
        </div>
      }}</Estados>
    </div>
  )
}
