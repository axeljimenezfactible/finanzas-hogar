// P-17 Todas las líneas (U1 §2.1, §2.2, §3.19: N-190..N-194). Galería de líneas, totales, aviso A/B, tendencias y calendario unificado.
import { useEffect, type ReactNode } from 'react'
import { CalendarDays, Calculator, Plus } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { Enlace, conConsulta, navegar } from '../../app/enrutador'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { itemsDe, Agenda, Barra, Datos, FaltaCapturar, HojaPago, Rejilla, Seccion, enAgenda, rutaCondiciones, usarPersonas } from './comun'
import { TIPO_CREDITO, plazo, t } from './textos'

type Credito = Esquemas['Credit']
const c = t.creditos
export const conFaltantes = (l: Credito) => (l.kind === 'bank_card' || l.kind === 'store_credit') && !!l.missing_terms?.length

export function FranjaAB({ lineas }: { lineas: Credito[] }) {
  const sin = lineas.filter((l) => l.kind === 'revolving_line' && !l.availability_confirmed)
  if (!sin.length) return null
  return <>{sin.map((l) => <Franja key={l.id} rol="advertencia" accion={<Enlace a={rutaCondiciones(l.id)} className="inline-flex min-h-11 items-center text-small text-text-link underline">{c.confirmarEnAjustes}</Enlace>}>{c.lecturaActiva(l.availability_reading)} · {l.name} · {c.avisoAB}</Franja>)}</>
}
function TarjetaLinea({ l }: { l: Credito }) {
  const disp = l.kind === 'revolving_line' ? [c.disponibleLectura(l.availability_reading), l.available_cents] : l.kind === 'bank_card' ? [c.disponible, l.available_cents] : [c.disponible, null]
  return (
    <article className="flex flex-col gap-2 rounded-lg border border-border-divider bg-surface-card p-4" data-prueba="tarjeta-linea">
      <Enlace a={rutaDe('P-18', { id: l.id })} className="inline-flex min-h-11 items-center text-title-3 font-semibold text-text-default hover:underline focus-visible:outline-2 focus-visible:outline-border-focus">{l.name}</Enlace>
      <p className="text-small text-text-secondary">{TIPO_CREDITO[l.kind]}{l.lender ? ` · ${l.lender}` : ''}{l.annual_rate_pct ? ` · ${c.tasa} ${c.tasaCapturada(l.annual_rate_pct)}` : ''}</p>
      <Datos className="border-0 p-0" filas={[
        [c.deudaViva, l.live_debt_cents == null ? t.comun.sinDato : dinero(l.live_debt_cents)],
        [disp[0] as string, disp[1] == null ? c.noAplica : dinero(disp[1] as number)],
        [c.proximoPago, l.next_payment_date ? `${dinero(l.next_payment_cents ?? 0)} · ${fecha(l.next_payment_date)}` : c.sinProximo],
        [l.kind === 'term_loan' ? c.montoOriginal : c.limite, l.credit_limit_cents == null ? t.comun.sinDato : dinero(l.credit_limit_cents)],
        [c.comprasVivas, String(l.live_purchase_count ?? 0)],
        ...(conFaltantes(l) ? [[c.pagoMinimo, <FaltaCapturar credito={l.id} que={l.missing_terms} />] as [string, ReactNode]] : []),
      ]} />
    </article>
  )
}
const porProximo = (a: Credito, b: Credito) => (a.next_payment_date ?? '9999').localeCompare(b.next_payment_date ?? '9999') || a.name.localeCompare(b.name, 'es')

export default function P17Lineas({ sub, consulta }: PropsPantalla) {
  const lineas = usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } })
  const totales = usarConsulta('get', '/api/v1/credits/totals', undefined, sub === 'resumen')
  const pagos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { limit: 200 } }, sub === 'calendario' || consulta.get('hoja') === 'pago')
  const promos = usarConsulta('get', '/api/v1/promotions', { consulta: { valid_on: hoyISO() } }, sub === 'calendario')
  const trend = usarConsulta('get', '/api/v1/credits/debt-trend', undefined, sub === 'tendencias')
  const per = usarPersonas()
  useEffect(() => { if (consulta.get('hoja') === 'pago') pagos.recargar() }, [consulta.get('hoja')]) // eslint-disable-line react-hooks/exhaustive-deps
  const nombre = (id?: number) => itemsDe<Credito>(lineas.datos).find((l) => l.id === id)?.name ?? ''
  const primera = itemsDe<Credito>(lineas.datos)[0]?.id
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <Barra>
        <Boton variante="secundario" onClick={() => navegar(conConsulta({ sub: 'calendario' }))}><CalendarDays aria-hidden className="size-4" />{c.verCalendario}</Boton>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-20'))}><Calculator aria-hidden className="size-4" />{c.simular}</Boton>
        <Boton variante="primario" deshabilitadoPor={primera ? undefined : c.vacio} onClick={() => primera && navegar(rutaDe('P-19', { id: primera, compra: 'nueva' }))}><Plus aria-hidden className="size-4" />{c.registrarCompra}</Boton>
      </Barra>
      <ConDatos consulta={lineas}>{(d) => {
        const ls = [...d.items].sort(porProximo)
        if (!ls.length) return <Vacio texto={c.vacio} accion={<Enlace a={rutaCondiciones()} className="inline-flex min-h-11 items-center text-text-link underline">{c.irAjustes}</Enlace>} />
        if (sub === 'tendencias') return <ConDatos consulta={trend}>{(tr) => <div className="flex flex-col gap-6">
          <Grafica titulo={c.tendenciaDeuda} resumen={c.resumenTendencia} tipo="lineas" categorias={tr.months.map((m) => mes(m.month))} series={ls.map((l) => ({ nombre: l.name, valores: tr.months.map((m) => m.credits.find((x) => x.credit_id === l.id)?.live_debt_cents ?? 0) }))} />
          <Grafica titulo={c.interesesPorMes} resumen={c.resumenIntereses} tipo="barras" categorias={tr.months.map((m) => mes(m.month))} series={[{ nombre: c.pagados, valores: tr.months.map((m) => m.interest_paid_cents) }, { nombre: c.porPagar, valores: tr.months.map((m) => m.interest_pending_cents) }]} />
        </div>}</ConDatos>
        if (sub === 'calendario') return <ConDatos consulta={pagos}>{(pg) => <div className="flex flex-col gap-6">
          <FranjaAB lineas={ls} />
          <Seccion titulo={c.agenda}><Agenda pagos={enAgenda(pg.items)} per={per} nombreLinea={nombre} /></Seccion>
          <Seccion titulo={c.promociones}>{itemsDe<Esquemas['Promotion']>(promos.datos).filter((p) => p.is_active).map((p) => <Franja key={p.id} rol="info">{c.promo(p.merchant_name, plazo(p.term_count), dinero(p.min_purchase_cents), fecha(p.valid_from), fecha(p.valid_to))}</Franja>)}</Seccion>
          <div className="hidden tablet:block"><Seccion titulo={c.pagosProgramados}><VistaBD def={{ clave: 'P-17-pagos', filas: pg.items, idFila: (p) => p.id, vistas: ['calendario', 'tabla'], fecha: (p) => p.due_date, titulo: (p) => `${nombre(p.credit_id)} · ${dinero(p.amount_cents)}`, alAbrir: (p) => navegar(conConsulta({ hoja: 'pago', pago: String(p.id) })),
            propiedades: [{ id: 'due_date', nombre: t.comun.fecha, valor: (p) => p.due_date, pintar: (p) => fechaConDia(p.due_date) }, { id: 'linea', nombre: c.linea, valor: (p) => nombre(p.credit_id), grupo: (p) => nombre(p.credit_id) }, { id: 'amount', nombre: t.comun.monto, valor: (p) => p.amount_cents, pintar: (p) => dinero(p.amount_cents), alinear: 'derecha' }] }} /></Seccion></div>
        </div>}</ConDatos>
        if (sub === 'detalle') return <VistaBD def={{ clave: 'P-17-lineas', filas: ls, idFila: (l) => l.id, vistas: ['galeria', 'tabla'], titulo: (l) => l.name, alAbrir: (l) => navegar(rutaDe('P-18', { id: l.id })),
          propiedades: [{ id: 'kind', nombre: c.tipo, valor: (l) => TIPO_CREDITO[l.kind], grupo: (l) => TIPO_CREDITO[l.kind] }, { id: 'deuda', nombre: c.deudaViva, valor: (l) => l.live_debt_cents ?? 0, pintar: (l) => dinero(l.live_debt_cents ?? 0), alinear: 'derecha' }, { id: 'disp', nombre: c.disponible, valor: (l) => l.available_cents ?? null, pintar: (l) => (l.available_cents == null ? c.noAplica : dinero(l.available_cents)), alinear: 'derecha' }, { id: 'prox', nombre: c.proximoPago, valor: (l) => l.next_payment_date ?? '', pintar: (l) => (l.next_payment_date ? `${dinero(l.next_payment_cents ?? 0)} · ${fecha(l.next_payment_date)}` : c.sinProximo) }, { id: 'vivas', nombre: c.comprasVivas, valor: (l) => l.live_purchase_count ?? 0, alinear: 'derecha' }] }} />
        return <div className="flex flex-col gap-6">
          <ConDatos consulta={totales}>{(to) => <div className="grid grid-cols-1 gap-3 tablet:grid-cols-2 desktop:grid-cols-3">
            <TarjetaMetrica titulo={c.deudaVivaTotal} valor={dinero(to.live_debt_cents)} rol="debt" />
            <TarjetaMetrica titulo={c.pagosQuincena} valor={dinero(to.fortnight_total_cents)} nota={c.partesHogarKaren(dinero(to.fortnight_household_cents), dinero(to.fortnight_external_cents))} alAbrir={() => navegar(conConsulta({ sub: 'calendario' }))} />
            <TarjetaMetrica titulo={c.pagosMes} valor={dinero(to.month_total_cents)} nota={c.partesHogarKaren(dinero(to.month_household_cents), dinero(to.month_external_cents))} alAbrir={() => navegar(conConsulta({ sub: 'calendario' }))} />
            <TarjetaMetrica titulo={c.interesesMes} valor={dinero(to.interest_paid_month_cents)} />
            <TarjetaMetrica titulo={c.interesesAnio} valor={dinero(to.interest_paid_year_cents)} />
          </div>}</ConDatos>
          <FranjaAB lineas={ls} />
          <Seccion titulo={c.lineas}><Rejilla>{ls.map((l) => <TarjetaLinea key={l.id} l={l} />)}</Rejilla></Seccion>
        </div>
      }}</ConDatos>
      <HojaPago pagos={pagos.datos?.items ?? []} per={per} alCambiar={pagos.recargar} consulta={consulta} />
    </div>
  )
}
