// P-24 Reporte mensual, vista web (U1 §2.1, §2.2, §3.26: N-255..N-257; S7-11). Mismas cifras que el PDF: aquí no se
// calcula nada, solo se pinta lo que devuelve /monthly-report. Descargar el PDF o generarlo si no existe (N-248).
import { useState, type ReactNode } from 'react'
import type { PropsPantalla } from '../../app/registro'
import { conConsulta, navegar } from '../../app/enrutador'
import { usarConsulta } from '../../api/consulta'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { ErrorApi } from '../../api/errores'
import { dinero, fecha, fechaHora, mes, porcentaje } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { useAvisos } from '../../componentes/Avisos'
import { ConDatos, Esqueleto, Vacio } from '../../componentes/Estados'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { AccionesPdf, Datos, Rejilla, Seccion, Selector, itemsDe, mensajeError, usarPersonas } from './comun'
import { mesPorDefecto, nombrePdf } from './P23Centro'
import { ESTADO_PARTE, t } from './textos'

const m = t.mensual
const din = (n?: number | null) => (n == null ? t.comun.sinDato : dinero(n))
const kc = (n?: number | null) => (n == null ? t.comun.sinDato : t.comun.kcal(n))
const pct = (n?: number | null) => (n == null ? t.comun.sinDato : porcentaje(n))
type Reporte = Esquemas['MonthlyReport']
type Generado = Esquemas['GeneratedReport']
const esqueleto = <div className="flex flex-col gap-6">{[m.resumen, m.categorias, m.quincenas, m.creditos].map((s) => <div key={s} className="flex flex-col gap-2"><Esqueleto className="w-1/3" /><Esqueleto className="h-24" /></div>)}</div>

export default function P24Mensual({ consulta }: PropsPantalla) {
  const meses = usarConsulta('get', '/api/v1/monthly-report/months')
  const ms = meses.datos?.months ?? []
  const elegido = consulta.get('mes') ?? mesPorDefecto(ms)
  const rep = usarConsulta('get', '/api/v1/monthly-report', { consulta: { period: elegido } }, !!elegido)
  const sinDatos = rep.error instanceof ErrorApi && (rep.error.estado === 404 || rep.error.estado === 422)
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <ConDatos consulta={meses}>{() => !ms.length ? <Vacio texto={m.sinMeses} /> : <>
        <Selector className="tablet:w-64" etiqueta={m.cambiarMes} valor={elegido} alCambiar={(v) => navegar(conConsulta({ mes: v }), { reemplazar: true })} opciones={[...ms].sort((a, b) => b.month.localeCompare(a.month)).map((x) => ({ valor: x.month, texto: `${mes(x.month)}${x.is_partial ? ` (${t.comun.parcial})` : ''}` }))} />
        {sinDatos ? <Vacio texto={m.sinDatos(mes(elegido))} /> : <ConDatos consulta={rep} esqueleto={esqueleto}>{(r) => <Reporte r={r} />}</ConDatos>}
      </>}</ConDatos>
    </div>
  )
}

function Pdf({ r }: { r: Reporte }) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const [nuevo, setNuevo] = useState<Generado | null>(null)
  const [generando, setGenerando] = useState(false)
  const [clave] = useState(nuevaClave)
  const pdf = nuevo ?? r.pdf
  const generar = async () => {
    setGenerando(true)
    try { setNuevo(await pedir('post', '/api/v1/reports/monthly', { cuerpo: { month: r.month, send_telegram: false }, idempotencia: clave })); avisar({ texto: t.comun.pdfListo, rol: 'exito' }) } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) } finally { setGenerando(false) }
  }
  return <Seccion titulo={m.pdf}>
    {pdf ? <div className="flex flex-wrap items-center justify-between gap-3"><p className="text-small text-text-secondary">{t.comun.generadoEl(fechaHora(pdf.generated_at))}</p><AccionesPdf id={pdf.id} nombre={nombrePdf(pdf)} principal /></div>
      : <div className="flex flex-wrap items-center justify-between gap-3"><p className="text-small text-text-secondary">{m.sinPdf}</p><Boton variante="primario" deshabilitadoPor={generando ? t.comun.generando : enLinea ? undefined : textos.conexion.requiere} onClick={generar}>{generando ? t.comun.generando : t.comun.generarPdf}</Boton></div>}
  </Seccion>
}

function Reporte({ r }: { r: Reporte }) {
  const per = usarPersonas()
  const cats = usarConsulta('get', '/api/v1/categories', { consulta: { limit: 200 } })
  const tiendas = usarConsulta('get', '/api/v1/stores', { consulta: { limit: 200 } })
  const cat = (id: number) => itemsDe<Esquemas['Category']>(cats.datos).find((c) => c.id === id)?.name ?? t.comun.sinDato
  const tienda = (id?: number) => itemsDe<Esquemas['Store']>(tiendas.datos).find((s) => s.id === id)?.name ?? t.comun.sinDato
  const b = r.budget
  const filasK: [string, ReactNode][] = []
  return (
    <div className="flex flex-col gap-8">
      <Pdf r={r} />
      <Seccion titulo={m.resumen}><Datos filas={[[m.ingreso, din(b.income_cents)], [m.plan, din(b.planned_spending_cents)], [m.real, din(b.actual_cents)], [m.ahorro, din(b.planned_savings_cents)], [m.restante, din(b.remaining_cents)], [m.sinAsignar, din(b.unassigned_cents)], ...(b.used_pct != null ? [[m.usado, pct(b.used_pct)] as [string, string]] : [])]} /></Seccion>
      <Seccion titulo={m.categorias}><VistaBD def={{ clave: 'P-24-categorias', filas: r.categories.rows, idFila: (x) => x.category_id, vistas: ['tabla'], titulo: (x) => cat(x.category_id),
        propiedades: [{ id: 'c', nombre: m.categoria, valor: (x) => cat(x.category_id) }, { id: 'p', nombre: m.plan, valor: (x) => x.planned_cents, pintar: (x) => din(x.planned_cents), alinear: 'derecha' }, { id: 'r', nombre: m.real, valor: (x) => x.actual_cents, pintar: (x) => din(x.actual_cents), alinear: 'derecha' }, { id: 'd', nombre: m.diferencia, valor: (x) => x.difference_cents, pintar: (x) => din(x.difference_cents), alinear: 'derecha' }, { id: 'u', nombre: m.usado, valor: (x) => x.used_pct ?? null, pintar: (x) => pct(x.used_pct), alinear: 'derecha' }, { id: 'n', nombre: m.gastos, valor: (x) => x.expense_count, alinear: 'derecha' }] }} /></Seccion>
      <Seccion titulo={m.quincenas}><Rejilla anchas>{r.fortnights.fortnights.map((q) => <Datos key={q.half} filas={[[m.quincena(q.half), `${fecha(q.start_date)} al ${fecha(q.end_date)}`], [m.plan, din(q.planned_cents)], [m.real, din(q.actual_cents)], [m.diaPago, fecha(q.pay_date)]]} />)}</Rejilla></Seccion>
      <Seccion titulo={m.creditos}><Datos filas={[[t.creditos.pagosMes, din(r.credits.month_total_cents)], [t.comun.hogar, din(r.credits.month_household_cents)], ['Karen', din(r.credits.month_external_cents)], [m.deudaCierre, din(r.live_debt_end_cents)], [t.creditos.interesesMes, din(r.credits.interest_paid_month_cents)]]} /></Seccion>
      <Seccion titulo={m.despensa}>{r.basket ? <Datos filas={[[m.masBarata, tienda(r.basket.cheapest_store_id)], ...r.basket.stores.map((s) => [tienda(s.store_id), `${din(s.total_cents)} · ${m.difBarata.toLowerCase()} ${din(s.difference_vs_cheapest_cents)}`] as [string, string])]} /> : <p className="text-small text-text-secondary">{m.sinCanasta}</p>}</Seccion>
      <Seccion titulo={m.comidas}><Datos filas={r.meals.weeks.map((w) => [`${m.semana} ${fecha(w.week_start)}`, `${din(w.cost_cents)} · ${w.persons.map((p) => `${per.nombre(p.person_id)} ${kc(p.avg_kcal_per_day)} / ${kc(p.daily_kcal_target)}`).join(' · ')}`])} /></Seccion>
      <Seccion titulo={m.karen}>{r.external_participation.map((p) => { filasK.length = 0
        filasK.push([m.despensaDevengada, din(p.accrued_groceries_cents)], [t.creditos.deudaViva, din(p.live_debt_cents)])
        for (const c of p.purchases) filasK.push([`${c.description} · ${fecha(c.purchase_date)}`, `${din(c.total_cents)} · ${c.payments_made} de ${c.payments_total} · ${t.creditos.saldo.toLowerCase()} ${din(c.balance_cents)}`])
        for (const s of p.payment_shares) filasK.push([`${t.pago.titulo} · ${fecha(s.due_date)}`, `${din(s.amount_cents)} · ${ESTADO_PARTE[s.status].toLowerCase()}`])
        return <div key={p.person_id} className="flex flex-col gap-2"><h3 className="text-small font-semibold">{per.nombre(p.person_id)}</h3><Datos filas={[...filasK]} /></div> })}</Seccion>
    </div>
  )
}
