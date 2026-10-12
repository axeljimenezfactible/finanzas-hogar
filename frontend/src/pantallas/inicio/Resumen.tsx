// Resumen de P-01 (U3 §2): 6 tarjetas en orden fijo; cada una es un solo enlace a su pantalla destino y hace su propia petición.
import type { ReactNode } from 'react'
import { ArrowDownLeft, ArrowUpRight, Banknote, CalendarClock, CalendarX, PiggyBank, ShoppingCart, Utensils, Wallet } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia, porcentaje } from '../../api/formatos'
import { MarcaPersona, MarcaTienda } from '../../componentes/Chip'
import { cn } from '../../componentes/utils'
import { Aislado, CargadoA, EnlaceAccion, EsqueletoTarjeta, Estados, IconoAviso, IconoInfo, IconoOk, Insignia, PieEnlace, codigoTienda, diasAviso, entero, listaDe, nombreCredito, titulares, useInicio } from './comun'
import { diaDe, difDias, mesAbrev, rangoSemana, semanaISO, semanaVigente } from './fechas'
import { MESES_LARGOS, textosInicio as t } from './textos'

interface Vista { valor: string; rolValor?: 'advertencia' | 'ahorro'; iconoValor?: ReactNode; etiquetaValor?: string; subtitulo?: string; secundario?: ReactNode; extra?: ReactNode; insignias?: ReactNode; estado?: string; accion: { a: string; texto: string } }
const COLOR_VALOR = { advertencia: 'text-role-warning', ahorro: 'text-role-savings' }
interface PropsTarjeta<T> { n: number; titulo: string; icono: ReactNode; consulta: EstadoConsulta<T>; falla: string; accionBase: { a: string; texto: string }; vista: (d: T) => Vista }
const CAJA = 'flex h-full min-h-11 flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4 shadow-0 tablet:p-6'
const Titulo = ({ icono, titulo }: { icono: ReactNode; titulo: string }) => <span className="flex items-center gap-2 text-small font-medium text-text-secondary [&_svg]:size-5">{icono}{titulo}</span>

function Tarjeta<T>({ n, titulo, icono, consulta, falla, accionBase, vista }: PropsTarjeta<T>) {
  const listo = !consulta.cargando || consulta.datos !== undefined
  if (listo && !consulta.error && !consulta.nuncaCargada && consulta.datos !== undefined) {
    return <Aislado texto={falla} alReintentar={consulta.recargar}><Cargada n={n} titulo={titulo} icono={icono} desde={consulta.desdeCopia} v={() => vista(consulta.datos as T)} /></Aislado>
  }
  return (
    <section data-prueba={`tarjeta-${n}`} aria-label={titulo} className={CAJA}>
      <Titulo icono={icono} titulo={titulo} />
      <Estados consulta={consulta} falla={falla} esqueleto={<EsqueletoTarjeta />}>{() => null}</Estados>
      <EnlaceAccion a={accionBase.a} texto={accionBase.texto} />
    </section>
  )
}
function Cargada({ n, titulo, icono, desde, v: armar }: { n: number; titulo: string; icono: ReactNode; desde: number | null; v: () => Vista }) {
  const v = armar()
  const nombre = [titulo, v.valor, v.etiquetaValor, v.estado, v.accion.texto].filter(Boolean).join(' · ')
  return (
    <Enlace a={v.accion.a} data-prueba={`tarjeta-${n}`} aria-label={nombre} className={cn(CAJA, 'text-text-default no-underline hover:border-border-control-hover hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus')}>
      <Titulo icono={icono} titulo={titulo} />
      {v.subtitulo && <span className="text-caption text-text-secondary">{v.subtitulo}</span>}
      <span className="flex flex-wrap items-baseline gap-x-2">
        <span className={cn('flex items-center gap-1 whitespace-nowrap text-title-2 font-semibold tabular-nums [&_svg]:size-5', v.rolValor ? COLOR_VALOR[v.rolValor] : 'text-text-default')}>{v.iconoValor}{v.valor}</span>
        {v.etiquetaValor && <span className="text-small text-text-secondary">{v.etiquetaValor}</span>}
      </span>
      {v.secundario && <span className="flex flex-col gap-1 text-small text-text-secondary tabular-nums">{v.secundario}</span>}
      {v.extra && <span className="hidden flex-col gap-1 text-small text-text-secondary tabular-nums phone:flex">{v.extra}</span>}
      {v.insignias && <span className="flex flex-wrap gap-2">{v.insignias}</span>}
      <CargadoA desde={desde} />
      <PieEnlace texto={v.accion.texto} />
    </Enlace>
  )
}

export default function Resumen() {
  const c = useInicio()
  const hoy = c.hoy
  const comidas = usarConsulta('get', '/api/v1/dashboard/week-meals', { consulta: { period: semanaISO(hoy) } })
  const canasta = usarConsulta('get', '/api/v1/dashboard/best-basket', { consulta: { week: semanaVigente(hoy) } })
  const ahorro = usarConsulta('get', '/api/v1/dashboard/savings')
  const idGasto = canasta.datos?.basket === 'expense' ? canasta.datos.expense_id : undefined
  const gasto = usarConsulta('get', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: idGasto ?? 0 } }, idGasto !== undefined)
  const aviso = diasAviso(c)
  const tiendas = listaDe(c.tiendas)
  const personas = titulares(c)
  return (
    <div data-prueba="resumen" className="grid grid-cols-1 gap-4 tablet:grid-cols-2 desktop:grid-cols-3 desktop:gap-6">
      <Tarjeta n={1} titulo={t.t1.titulo} icono={<Wallet aria-hidden />} consulta={c.plan} falla={t.t1.falla} accionBase={{ a: rutaDe('P-06'), texto: t.t1.accion }} vista={(d) => {
        const conPlan = d.planned_cents !== undefined
        const rebasado = (d.remaining_cents ?? 0) < 0
        return {
          valor: !conPlan ? t.comun.sinDato : rebasado ? `${dinero(d.remaining_cents!)} ${t.t1.rebasado}` : dinero(d.remaining_cents ?? 0),
          etiquetaValor: conPlan && !rebasado ? t.t1.teQuedan : undefined, rolValor: rebasado ? 'advertencia' : undefined, iconoValor: rebasado ? <IconoAviso /> : undefined,
          secundario: conPlan ? <span>{t.t1.deDe(dinero(d.actual_cents), dinero(d.planned_cents!), d.used_pct !== undefined ? porcentaje(d.used_pct) : t.comun.sinDato)}</span>
            : <><span>{t.t1.sinPlan(MESES_LARGOS[Number(d.month.slice(5, 7)) - 1])}</span><span>{t.t1.llevas(dinero(d.actual_cents))}</span></>,
          extra: conPlan ? <span>{t.t1.faltan(d.days_left)}</span> : undefined,
          estado: d.status === 'warning' ? t.t1.vasEn(d.used_pct ?? 0) : d.status === 'exceeded' ? t.t1.rebasaste : undefined,
          insignias: d.status === 'warning' ? <Insignia rol="advertencia" icono={<IconoAviso />}>{t.t1.vasEn(d.used_pct ?? 0)}</Insignia> : d.status === 'exceeded' ? <Insignia rol="advertencia" icono={<IconoAviso />}>{t.t1.rebasaste}</Insignia> : undefined,
          accion: conPlan ? { a: rutaDe('P-06'), texto: t.t1.accion } : { a: rutaDe('P-26', { seccion: '02' }), texto: t.t1.accionVacio },
        }
      }} />
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
      <Tarjeta n={3} titulo={t.t3.titulo} icono={<CalendarClock aria-hidden />} consulta={c.proximo} falla={t.t3.falla} accionBase={{ a: rutaDe('P-17'), texto: t.t3.accion }} vista={(d) => {
        const np = d.next_payment, av = d.availability ?? []
        if (!np && !av.length) return { valor: t.t3.sinCreditos, accion: { a: rutaDe('P-26', { seccion: '05' }), texto: t.t3.accionVacio } }
        const linea = av.find((x) => x.credit_id === np?.credit_id) ?? av[0]
        const porVencer = np && !np.is_overdue && aviso.includes(np.days_left)
        const estadoTexto = np?.is_overdue ? t.pago.vencido : porVencer ? (np.days_left === 0 ? t.pago.venceHoy : np.days_left === 1 ? t.pago.venceManana : t.pago.venceEn(np.days_left)) : undefined
        const disponible = linea && (
          <span className="flex flex-col gap-1 border-t border-border-divider pt-2">
            <span className="tabular-nums text-text-default">{t.t3.disponibleHoy(dinero(linea.available_cents), linea.availability_reading)}</span>
            {!linea.availability_confirmed && <Insignia rol="info" icono={<IconoInfo />} etiqueta={t.t3.reposicion(nombreCredito(c, linea.credit_id))}>{t.t3.pendienteConfirmar}</Insignia>}
            {av.length > 1 && <span>{t.t3.masLineas(av.length - 1)}</span>}
          </span>
        )
        return {
          valor: np ? dinero(np.amount_cents) : t.t3.sinPagos,
          secundario: <>{np && <span>{`${fechaConDia(np.due_date)} · ${nombreCredito(c, np.credit_id)}`}</span>}{np && np.payment_count > 1 && <span>{t.t3.variosPagos(np.payment_count)}</span>}{disponible}</>,
          extra: np ? <><span>{t.comun.hogarYKaren(dinero(np.household_share_cents), np.external_share_cents > 0 ? dinero(np.external_share_cents) : undefined)}</span>
            {np.is_overdue ? <span>{t.t3.vencioHace(Math.max(1, difDias(np.due_date, hoy)))}</span> : !porVencer && <span>{t.t3.faltan(np.days_left)}</span>}</> : undefined,
          estado: estadoTexto,
          insignias: np?.is_overdue ? <Insignia rol="advertencia" icono={<CalendarX aria-hidden />}>{t.pago.vencido}</Insignia> : porVencer ? <Insignia rol="advertencia" icono={<CalendarClock aria-hidden />}>{estadoTexto}</Insignia> : undefined,
          accion: { a: rutaDe('P-18', { id: np?.credit_id ?? linea!.credit_id }), texto: t.t3.accion },
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
  )
}
