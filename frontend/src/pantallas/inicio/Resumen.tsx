// Tarjeta de métrica de P-01 (U3 §1.2) y «Más del hogar» (U3 2.0 §2.5, PROVISIONAL, bloqueo U3-tablero-mas-del-hogar): las tarjetas
// 2, 4, 5 y 6 de U3 1.x sin cambio de contenido. La 1 pasó a K-1 y la 3 a K-3 y K-4 (Tablero.tsx). Cada tarjeta es un solo enlace.
import type { ReactNode } from 'react'
import { ArrowDownLeft, ArrowUpRight, Banknote, PiggyBank, ShoppingCart, Utensils } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { dinero, fecha, porcentaje } from '../../api/formatos'
import { MarcaPersona, MarcaTienda } from '../../componentes/Chip'
import { cn } from '../../componentes/utils'
import { Aislado, CargadoA, EnlaceAccion, EsqueletoTarjeta, Estados, IconoAviso, IconoInfo, IconoOk, Insignia, PieEnlace, codigoTienda, entero, listaDe, titulares, useInicio } from './comun'
import { diaDe, mesAbrev, rangoSemana, semanaISO, semanaVigente } from './fechas'
import { textosInicio as t } from './textos'

export interface Vista { valor?: string; cuerpo?: ReactNode; rolValor?: 'advertencia' | 'ahorro' | 'ingreso' | 'deuda'; iconoValor?: ReactNode; etiquetaValor?: string; subtitulo?: string; secundario?: ReactNode; extra?: ReactNode; insignias?: ReactNode; estado?: string; accion: { a: string; texto: string } }
const COLOR_VALOR = { advertencia: 'text-role-warning', ahorro: 'text-role-savings', ingreso: 'text-role-income', deuda: 'text-role-debt' }
interface PropsTarjeta<T> { n: number | string; className?: string; titulo: string; icono: ReactNode; consulta: EstadoConsulta<T>; falla: string; accionBase: { a: string; texto: string }; vista: (d: T) => Vista }
export const CAJA = 'flex h-full min-h-11 flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4 shadow-0 tablet:p-6'
export const Titulo = ({ icono, titulo }: { icono: ReactNode; titulo: string }) => <span className="flex items-center gap-2 text-small font-medium text-text-secondary [&_svg]:size-5">{icono}{titulo}</span>

export function Tarjeta<T>({ n, className, titulo, icono, consulta, falla, accionBase, vista }: PropsTarjeta<T>) {
  const listo = !consulta.cargando || consulta.datos !== undefined
  if (listo && !consulta.error && !consulta.nuncaCargada && consulta.datos !== undefined) {
    return <Aislado texto={falla} alReintentar={consulta.recargar}><Cargada n={n} className={className} titulo={titulo} icono={icono} desde={consulta.desdeCopia} v={() => vista(consulta.datos as T)} /></Aislado>
  }
  return (
    <section data-prueba={`tarjeta-${n}`} aria-label={titulo} className={cn(CAJA, className)}>
      <Titulo icono={icono} titulo={titulo} />
      <Estados consulta={consulta} falla={falla} esqueleto={<EsqueletoTarjeta />}>{() => null}</Estados>
      <EnlaceAccion a={accionBase.a} texto={accionBase.texto} />
    </section>
  )
}
function Cargada({ n, className, titulo, icono, desde, v: armar }: { n: number | string; className?: string; titulo: string; icono: ReactNode; desde: number | null; v: () => Vista }) {
  const v = armar()
  const nombre = [titulo, v.valor, v.etiquetaValor, v.estado, v.accion.texto].filter(Boolean).join(' · ')
  return (
    <Enlace a={v.accion.a} data-prueba={`tarjeta-${n}`} aria-label={nombre} className={cn(CAJA, className, 'text-text-default no-underline hover:border-border-control-hover hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus')}>
      <Titulo icono={icono} titulo={titulo} />
      {v.subtitulo && <span className="text-caption text-text-secondary">{v.subtitulo}</span>}
      {v.valor !== undefined && <span className="flex flex-wrap items-baseline gap-x-2">
        <span className={cn('flex items-center gap-1 whitespace-nowrap text-title-2 font-semibold tabular-nums [&_svg]:size-5', v.rolValor ? COLOR_VALOR[v.rolValor] : 'text-text-default')}>{v.iconoValor}{v.valor}</span>
        {v.etiquetaValor && <span className="text-small text-text-secondary">{v.etiquetaValor}</span>}
      </span>}
      {v.cuerpo}
      {v.secundario && <span className="flex flex-col gap-1 text-small text-text-secondary tabular-nums">{v.secundario}</span>}
      {v.extra && <span className="hidden flex-col gap-1 text-small text-text-secondary tabular-nums phone:flex">{v.extra}</span>}
      {v.insignias && <span className="flex flex-wrap gap-2">{v.insignias}</span>}
      <CargadoA desde={desde} />
      <PieEnlace texto={v.accion.texto} />
    </Enlace>
  )
}

/** «Más del hogar» (U3 2.0 §2.5): título de sección y las tarjetas 2, 4, 5 y 6; 1 columna, 2 × 2 desde 768 y 4 en 1440. */
export function MasDelHogar({ className }: { className?: string }) {
  const c = useInicio()
  const hoy = c.hoy
  const comidas = usarConsulta('get', '/api/v1/dashboard/week-meals', { consulta: { period: semanaISO(hoy) } })
  const canasta = usarConsulta('get', '/api/v1/dashboard/best-basket', { consulta: { week: semanaVigente(hoy) } })
  const ahorro = usarConsulta('get', '/api/v1/dashboard/savings')
  const idGasto = canasta.datos?.basket === 'expense' ? canasta.datos.expense_id : undefined
  const gasto = usarConsulta('get', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: idGasto ?? 0 } }, idGasto !== undefined)
  const tiendas = listaDe(c.tiendas)
  const personas = titulares(c)
  return (
    <section data-prueba="mas-del-hogar" aria-labelledby="mas-del-hogar" className={cn('flex flex-col gap-4', className)}>
      <h2 id="mas-del-hogar" className="text-title-3 font-semibold text-text-default">{t.tablero.masDelHogar}</h2>
      <div className="grid grid-cols-1 gap-4 tablet:grid-cols-2 desktop:gap-6 wide:grid-cols-4">
      <Tarjeta n={2} titulo={t.t2.titulo} icono={<Banknote aria-hidden />} consulta={c.flujo} falla={t.t2.falla} accionBase={{ a: rutaDe('P-07'), texto: t.t2.accion }} vista={(d) => {
        // PROVISIONAL (bloqueo WF3-inicio-vacios): «sin quincena capturada ni esperada» = income_cents en 0.
        const sinIngreso = d.income_cents === 0
        const negativo = d.closing_balance_cents < 0
        const entra = <span className="inline-flex items-center gap-1"><ArrowDownLeft aria-hidden className="size-5 text-role-income" /><span className={sinIngreso ? '' : 'text-role-income'}>{sinIngreso ? t.comun.sinDato : dinero(d.income_cents, 'ingreso')}</span></span>
        const sale = (m: number) => <span className="inline-flex items-center gap-1"><ArrowUpRight aria-hidden className="size-5 text-role-expense" /><span className="text-role-expense">{dinero(m, 'gasto')}</span></span>
        return {
          subtitulo: t.t2.delAl(diaDe(d.start_date), diaDe(d.end_date), mesAbrev(d.end_date)),
          valor: negativo ? `${dinero(d.closing_balance_cents)} ${t.t2.faltante}` : dinero(d.closing_balance_cents), etiquetaValor: negativo ? undefined : t.t2.alCierre,
          rolValor: negativo ? 'advertencia' : undefined, iconoValor: negativo ? <IconoAviso /> : undefined,
          secundario: <><span className="flex flex-wrap items-center gap-1">{t.t2.entra} {entra} · {t.t2.sale} {sale(d.spent_cents + d.committed_cents)}</span>{sinIngreso && <span>{t.t2.faltaIngreso}</span>}</>,
          extra: <span className="flex flex-wrap items-center gap-1">{t.t2.yaSalio} {sale(d.spent_cents)} · {t.t2.porPagar} {sale(d.committed_cents)}</span>,
          estado: d.income_is_expected ? t.t2.esperado : undefined,
          insignias: d.income_is_expected && !sinIngreso ? <Insignia rol="info" icono={<IconoInfo />}>{t.t2.esperado}</Insignia> : undefined,
          accion: sinIngreso ? { a: rutaDe('P-26', { seccion: '03' }), texto: t.t2.accionVacio } : { a: rutaDe('P-07'), texto: t.t2.accion },
        }
      }} />
      <Tarjeta n={4} titulo={t.t4.titulo} icono={<Utensils aria-hidden />} consulta={comidas} falla={t.t4.falla} accionBase={{ a: rutaDe('P-10'), texto: t.t4.accion }} vista={(d) => {
        const vacia = d.unplanned_slots >= 21 && d.cost_cents === 0
        return {
          subtitulo: rangoSemana(d.week_start, d.week_end),
          valor: vacia ? t.comun.sinDato : dinero(d.cost_cents), etiquetaValor: vacia ? undefined : t.t4.costo,
          secundario: vacia ? <span>{t.t4.sinPlan}</span> : d.unplanned_slots === 0 ? <span className="inline-flex items-center gap-1"><IconoOk />{t.t4.completa}</span> : <span>{t.t4.sinArmar(d.unplanned_slots)}</span>,
          extra: vacia ? undefined : <>{personas.map((p) => { const x = d.persons.find((q) => q.person_id === p.id); if (!x) return null
            return <span key={p.id} className="inline-flex items-center gap-2"><MarcaPersona codigo={p.code} />{!x.days_with_plan || x.avg_kcal_per_day === undefined ? t.t4.personaSinDias(p.display_name) : t.t4.persona(p.display_name, entero(x.avg_kcal_per_day), x.daily_kcal_target !== undefined ? entero(x.daily_kcal_target) : null, x.days_with_plan)}</span> })}</>,
          accion: { a: rutaDe('P-10'), texto: vacia ? t.t4.accionVacio : t.t4.accion },
        }
      }} />
      <Tarjeta n={5} titulo={t.t5.titulo} icono={<ShoppingCart aria-hidden />} consulta={canasta} falla={t.t5.falla} accionBase={{ a: rutaDe('P-14'), texto: t.t5.accion }} vista={(d) => {
        if (d.store_id === undefined || d.total_cents === undefined) return { valor: t.comun.sinDato, secundario: <span>{t.t5.sinPrecios}</span>, accion: { a: rutaDe('P-15'), texto: t.t5.accionVacio } }
        const tienda = tiendas.find((s) => s.id === d.store_id), segunda = tiendas.find((s) => s.id === d.second_store_id)
        const dif = d.difference_vs_second_cents
        return {
          valor: `${tienda?.name ?? t.comun.sinDato} ${dinero(d.total_cents)}`, iconoValor: <MarcaTienda codigo={codigoTienda(tienda)} />,
          secundario: <>{dif !== undefined && <span>{segunda ? t.t5.menosQue(dinero(dif), segunda.name) : t.t5.menosQueSiguiente(dinero(dif))}</span>}
            {d.basket === 'expense' && <span>{t.t5.listaVacia(gasto.datos?.expense_date ? fecha(gasto.datos.expense_date) : t.comun.sinDato)}</span>}</>,
          extra: <span>{t.t5.cobertura(d.covered_count, d.item_count)}</span>,
          insignias: (d.stale_price_count > 0 || d.assumes_membership) ? <>{d.stale_price_count > 0 && <span className="hidden phone:inline-flex"><Insignia rol="advertencia" icono={<IconoAviso />}>{t.t5.viejos(d.stale_price_count)}</Insignia></span>}{d.assumes_membership && <Insignia rol="info" icono={<IconoInfo />}>{t.t5.membresia}</Insignia>}</> : undefined,
          accion: { a: rutaDe('P-14'), texto: t.t5.accion },
        }
      }} />
      <Tarjeta n={6} titulo={t.t6.titulo} icono={<PiggyBank aria-hidden />} consulta={ahorro} falla={t.t6.falla} accionBase={{ a: rutaDe('P-09'), texto: t.t6.accion }} vista={(d) => {
        // PROVISIONAL (bloqueo WF3-inicio-vacios): «sin fondo ni metas» = sin meta, sin ahorro y sin plan del mes.
        if (!d.next_goal && d.saved_total_cents === 0 && d.planned_month_cents === 0) return { valor: t.comun.sinDato, secundario: <span>{t.t6.sinAhorro}</span>, accion: { a: rutaDe('P-26', { seccion: '13' }), texto: t.t6.accionVacio } }
        const g = d.next_goal
        const est = g?.status ? t.t6.estados[g.status] : undefined
        return {
          valor: dinero(d.saved_total_cents), rolValor: 'ahorro', iconoValor: <PiggyBank aria-hidden />, etiquetaValor: t.t6.ahorrado,
          secundario: <span>{t.t6.esteMes(dinero(d.contributed_month_cents), dinero(d.planned_month_cents))}</span>,
          extra: g ? <span>{t.t6.meta(g.name, g.progress_pct !== undefined ? porcentaje(g.progress_pct) : t.comun.sinDato, g.target_date ? fecha(g.target_date) : t.comun.sinDato)}</span> : undefined,
          estado: est,
          insignias: g?.status ? <span className="hidden phone:inline-flex"><Insignia rol={g.status === 'achieved' ? 'exito' : g.status === 'overdue' ? 'advertencia' : 'info'} icono={g.status === 'achieved' ? <IconoOk /> : g.status === 'overdue' ? <IconoAviso /> : <IconoInfo />}>{est}</Insignia></span> : undefined,
          accion: { a: rutaDe('P-09'), texto: t.t6.accion },
        }
      }} />
      </div>
    </section>
  )
}
