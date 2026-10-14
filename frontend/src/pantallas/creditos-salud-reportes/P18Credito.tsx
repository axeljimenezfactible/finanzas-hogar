// P-18 Crédito (U1 §2.1, §2.2, §3.20: N-200..N-210). Variante por tipo: revolvente (A/B, serie, tasas por plazo),
// tarjeta (estado del ciclo), préstamo (amortización) y tienda (abonos). Sin valores supuestos: «Falta capturar» (K18).
// Calendario de la línea (N-208 = N-192 filtrado por la línea): C1 listCalendarEvents con credit_id y los tipos de crédito.
import { useEffect, type ReactNode } from 'react'
import { Calculator, CircleCheck, Plus, Settings2 } from 'lucide-react'
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
import { itemsDe, Barra, Datos, FaltaCapturar, HojaPago, Seccion, agrupar, agruparEstado, enAgenda, rutaCondiciones, sumarDias, textoPartes, usarPersonas, type Pago, type Personas } from './comun'
import { FranjaLinea, conFaltantes } from './P17Lineas'
import { ESTADO_COMPRA, ESTADO_PARTE, EVENTO, EVENTO_LINEA, FRECUENCIA, TIPO_CREDITO, plazo, t } from './textos'

type Credito = Esquemas['Credit']
type Compra = Esquemas['CreditPurchase']
type Evento = Esquemas['CalendarEvent']
const c = t.creditos
/** Tipos de evento de una línea en el calendario (C1 listCalendarEvents: N-192 y N-208 = credit_payment, promotion y card_statement). */
const TIPOS_LINEA: ('credit_payment' | 'promotion' | 'card_statement')[] = ['credit_payment', 'promotion', 'card_statement']
/** Rango de N-208: del día 1 del mes en curso (o del pago sin marcar más viejo) a hoy + 60 días; sintaxis `desde..hasta` de C1 Period. */
export function periodoLinea(pagos: Pago[], hoy = hoyISO()) {
  const inicioMes = `${hoy.slice(0, 7)}-01`
  const viejo = pagos.filter((p) => p.status !== 'paid' && p.due_date < inicioMes).map((p) => p.due_date).sort()[0]
  return `${viejo ?? inicioMes}..${sumarDias(hoy, 60)}`
}
const claveEvento = (e: Evento) => `${e.event_type}-${e.ref_table}-${e.ref_id}`
const estadoEvento = (e: Evento) => (e.status === 'paid' || e.status === 'pending' || e.status === 'overdue' ? ESTADO_PARTE[e.status] : '')
/** Partes de un pago tal como las trae el evento (hogar y persona externa), si el pago programado no está en la lista cargada. */
function partesEvento(e: Evento, per: Personas) {
  const partes: [string, number | undefined][] = [[t.comun.hogar, e.household_share_cents], [per.externaPrincipal?.display_name ?? '', e.external_share_cents]]
  const est = estadoEvento(e)
  return [...partes.filter(([k, v]) => k && v).map(([k, v]) => `${k} ${dinero(v!)}`), ...(est ? [est.toLowerCase()] : [])].join(' · ')
}
const abrirEvento = (e: Evento, id: number) => (e.event_type === 'credit_payment' && e.ref_table === 'scheduled_payment' ? navegar(conConsulta({ hoja: 'pago', pago: String(e.ref_id) })) : navegar(rutaCondiciones(id)))
/** Agenda de la línea (N-208, móvil y escritorio): pagos y cortes de hoy a 60 días y los pagos sin marcar con fecha pasada, en el orden de C1. */
function AgendaLinea({ eventos, pagos, per }: { eventos: Evento[]; pagos: Pago[]; per: Personas }) {
  if (!eventos.length) return <p className="py-4 text-small text-text-secondary">{c.sinPagos}</p>
  return (
    <ul className="flex flex-col divide-y divide-border-divider rounded-lg border border-border-divider bg-surface-card" data-prueba="agenda-linea">
      {eventos.map((e) => {
        const pago = e.event_type === 'credit_payment' && e.ref_table === 'scheduled_payment' ? pagos.find((p) => p.id === e.ref_id) : undefined
        const detalle = e.event_type === 'credit_payment' ? (pago?.shares?.length ? agruparEstado(pago.shares, per) : partesEvento(e, per)) : e.event_type === 'card_statement' ? c.saldoCorte : EVENTO_LINEA[e.event_type as keyof typeof EVENTO_LINEA]
        const cuerpo = <>
          <span className="flex flex-wrap justify-between gap-2 text-body text-text-default"><span>{fechaConDia(e.date)} · {e.title}</span>{e.amount_cents != null && <span className="tabular-nums font-medium">{dinero(e.amount_cents)}</span>}</span>
          <span className="text-small tabular-nums text-text-secondary">{detalle}</span>
        </>
        return <li key={claveEvento(e)}>{e.event_type === 'credit_payment' && e.ref_table === 'scheduled_payment'
          ? <button type="button" data-prueba="pago-agenda" onClick={() => navegar(conConsulta({ hoja: 'pago', pago: String(e.ref_id) }))} className="flex min-h-11 w-full flex-col gap-1 px-4 py-3 text-left hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">{cuerpo}</button>
          : <div className="flex min-h-11 w-full flex-col gap-1 px-4 py-3">{cuerpo}</div>}</li>
      })}
    </ul>
  )
}
/** Pestaña «Calendario» de P-18 (N-208): agenda, promociones vigentes como franjas y, de tableta en adelante, el mes visible. */
function CalendarioLinea({ eventos, pagos, per, id }: { eventos: Evento[]; pagos: Pago[]; per: Personas; id: number }) {
  const hoy = hoyISO(), fin = sumarDias(hoy, 60)
  const agenda = eventos.filter((e) => e.event_type !== 'promotion' && e.date <= fin && (e.date >= hoy || (e.event_type === 'credit_payment' && e.status === 'overdue')))
  const promos = eventos.filter((e) => e.event_type === 'promotion' && (e.end_date ?? e.date) >= hoy)
  return <div className="flex flex-col gap-6">
    <Seccion titulo={c.agenda}><AgendaLinea eventos={agenda} pagos={pagos} per={per} /></Seccion>
    {promos.length > 0 && <Seccion titulo={c.promociones}>{promos.map((e) => <Franja key={claveEvento(e)} rol="info">{e.title} · {c.vigencia(fecha(e.date), fecha(e.end_date ?? e.date))}</Franja>)}</Seccion>}
    <div className="hidden tablet:block"><Seccion titulo={c.calendarioLinea}><VistaBD def={{ clave: 'P-18-calendario', filas: eventos, idFila: claveEvento, vistas: ['calendario', 'tabla'], fecha: (e) => e.date, titulo: (e) => (e.amount_cents != null ? `${e.title} · ${dinero(e.amount_cents)}` : e.title), alAbrir: (e) => abrirEvento(e, id),
      propiedades: [{ id: 'date', nombre: t.comun.fecha, valor: (e) => e.date, pintar: (e) => fechaConDia(e.date) }, { id: 'event_type', nombre: c.tipo, valor: (e) => EVENTO_LINEA[e.event_type as keyof typeof EVENTO_LINEA], grupo: (e) => EVENTO_LINEA[e.event_type as keyof typeof EVENTO_LINEA] }, { id: 'title', nombre: c.evento, valor: (e) => e.title }, { id: 'amount', nombre: t.comun.monto, valor: (e) => e.amount_cents ?? null, pintar: (e) => (e.amount_cents == null ? t.comun.sinDato : dinero(e.amount_cents)), alinear: 'derecha' }, { id: 'status', nombre: t.comun.estado, valor: (e) => estadoEvento(e) }] }} /></Seccion></div>
  </div>
}

function DatosLinea({ l }: { l: Credito }) {
  const falta = (k: string) => (conFaltantes(l) && l.missing_terms?.includes(k as never) ? <FaltaCapturar credito={l.id} que={[k]} /> : null)
  const filas: [string, ReactNode][] = [[c.tipo, TIPO_CREDITO[l.kind]], [l.kind === 'term_loan' ? c.montoOriginal : c.limite, l.credit_limit_cents == null ? t.comun.sinDato : dinero(l.credit_limit_cents)], [c.frecuencia, FRECUENCIA[l.payment_frequency]],
    [c.tasa, falta('annual_rate') ?? (l.annual_rate_pct ? c.tasaCapturada(l.annual_rate_pct) : t.comun.sinDato)]]
  if (l.statement_day) filas.push([c.diaCorte, c.diaDelMes(l.statement_day)])
  if (l.payment_due_day) filas.push([c.diaPago, c.diaDelMes(l.payment_due_day)])
  if (l.annual_fee_cents != null) filas.push([c.comisionAnual, dinero(l.annual_fee_cents)])
  for (const k of ['rate_base', 'vat', 'min_payment_pct_bp', 'min_payment_fixed_cents']) { const f = falta(k); if (f) filas.push([c.faltan[k][0].toUpperCase() + c.faltan[k].slice(1), f]) }
  if (l.notes) filas.push([c.notas, l.notes])
  return <Seccion titulo={c.datosLinea}><Datos filas={filas} /></Seccion>
}
function ResumenRevolvente({ l }: { l: Credito }) {
  const ab = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-preview', { ruta: { credit_id: l.id } })
  const tasas = usarConsulta('get', '/api/v1/credits/{credit_id}/term-rates', { ruta: { credit_id: l.id } })
  return <>
    <ConDatos consulta={ab}>{(p) => { const act = p.availability_reading; const otra = act === 'A' ? 'B' : 'A'; const v = (x: 'A' | 'B') => (x === 'A' ? p.reading_a : p.reading_b)
      return <div className="grid grid-cols-1 gap-3 tablet:grid-cols-2">
        <TarjetaMetrica titulo={`${c.disponibleLectura(act)} (${c.activa})`} valor={dinero(v(act).available_today_cents)} rol="savings" nota={`${c.minimos(dinero(p.reading_a.horizon_min_cents), dinero(p.reading_b.horizon_min_cents))}`} />
        <TarjetaMetrica titulo={c.disponibleLectura(otra)} valor={dinero(v(otra).available_today_cents)} />
      </div> }}</ConDatos>
    <Seccion titulo={c.tasasPlazo}><ConDatos consulta={tasas}>{(d) => <Datos filas={[...itemsDe<Esquemas['CreditTermRate']>(d).sort((a, b) => a.term_count - b.term_count).map((r) => [plazo(r.term_count, l.payment_frequency), r.rate_pct ? c.tasaCapturada(r.rate_pct) : t.comun.sinDato] as [string, ReactNode]), [c.otrosPlazos, <Enlace a={rutaCondiciones(l.id)} className="inline-flex min-h-11 items-center text-text-link underline">{c.capturar}</Enlace>]]} />}</ConDatos></Seccion>
  </>
}
function EstadoTarjeta({ l }: { l: Credito }) {
  const e = usarConsulta('get', '/api/v1/credits/{credit_id}/card-status', { ruta: { credit_id: l.id } })
  const v = (n?: number | null, dep = false) => (dep && conFaltantes(l) ? <FaltaCapturar credito={l.id} que={l.missing_terms} /> : n == null ? t.comun.sinDato : dinero(n))
  return <Seccion titulo={c.estadoTarjeta}><ConDatos consulta={e}>{(s) => <>
    <Datos filas={[[c.fechaCorte, s.statement_date ? fecha(s.statement_date) : t.comun.sinDato], [c.fechaLimite, s.due_date ? fecha(s.due_date) : t.comun.sinDato], [c.saldoCorte, v(s.statement_balance_cents)], [c.pagoSinIntereses, v(s.no_interest_payment_cents)], [c.pagoMinimo, v(s.minimum_payment_cents, true)], [c.saldoActual, v(s.current_balance_cents)]]} />
    {s.installment_plans.length > 0 && <Seccion titulo={c.msi}><Datos filas={s.installment_plans.map((m) => [m.description, c.mensualidad(dinero(m.installment_cents), m.remaining_count)] as [string, ReactNode])} /></Seccion>}
  </>}</ConDatos></Seccion>
}
function TablaCompras({ compras, per, l }: { compras: Compra[]; per: Personas; l: Credito }) {
  return <VistaBD def={{ clave: 'P-18-compras', filas: compras, idFila: (x) => x.id, vistas: ['tabla', 'tablero'], titulo: (x) => x.description ?? fecha(x.purchase_date), alAbrir: (x) => navegar(rutaDe('P-19', { id: l.id, compra: x.id })),
    columnasTablero: { propiedad: 'status', columnas: [ESTADO_COMPRA.active, ESTADO_COMPRA.paid_off, ESTADO_COMPRA.cancelled] },
    vacio: <Vacio texto={c.sinCompras} />,
    propiedades: [
      { id: 'description', nombre: c.comercio, valor: (x) => x.description ?? '' },
      { id: 'purchase_date', nombre: t.comun.fecha, valor: (x) => x.purchase_date, pintar: (x) => fecha(x.purchase_date) },
      { id: 'principal', nombre: c.capital, valor: (x) => x.principal_cents, pintar: (x) => dinero(x.principal_cents), alinear: 'derecha' },
      { id: 'interest', nombre: c.intereses, valor: (x) => x.interest_cents, pintar: (x) => dinero(x.interest_cents), alinear: 'derecha' },
      { id: 'total', nombre: t.comun.total, valor: (x) => x.total_cents, pintar: (x) => dinero(x.total_cents), alinear: 'derecha' },
      { id: 'term', nombre: c.plazo, valor: (x) => x.term_count, pintar: (x) => plazo(x.term_count, x.payment_frequency) },
      { id: 'shares', nombre: c.partes, valor: (x) => textoPartes(agrupar(x.shares, (s) => s.total_cents, per)) },
      { id: 'pagos', nombre: c.pagos, valor: (x) => x.payments_made ?? 0, pintar: (x) => `${x.payments_made ?? 0} de ${x.payments_total ?? x.term_count}`, alinear: 'derecha' },
      { id: 'balance', nombre: c.saldo, valor: (x) => x.balance_cents ?? 0, pintar: (x) => dinero(x.balance_cents ?? 0), alinear: 'derecha' },
      { id: 'status', nombre: t.comun.estado, valor: (x) => ESTADO_COMPRA[x.status], grupo: (x) => ESTADO_COMPRA[x.status] },
    ] }} />
}
function TablaPagos({ pagos, titulo }: { pagos: Pago[]; titulo: string }) {
  return <Seccion titulo={titulo}><VistaBD def={{ clave: 'P-18-pagos', filas: pagos, idFila: (p) => p.id, vistas: ['tabla'], titulo: (p) => `${p.seq}`, alAbrir: (p) => navegar(conConsulta({ hoja: 'pago', pago: String(p.id) })), vacio: <Vacio texto={c.sinPagosPendientes} />,
    propiedades: [{ id: 'seq', nombre: c.numero, valor: (p) => p.seq, alinear: 'derecha' }, { id: 'due', nombre: t.comun.fecha, valor: (p) => p.due_date, pintar: (p) => fecha(p.due_date) }, { id: 'amount', nombre: t.comun.monto, valor: (p) => p.amount_cents, pintar: (p) => dinero(p.amount_cents), alinear: 'derecha' }, { id: 'principal', nombre: c.capital, valor: (p) => p.principal_cents, pintar: (p) => dinero(p.principal_cents), alinear: 'derecha' }, { id: 'interest', nombre: c.interesesCorto, valor: (p) => p.interest_cents, pintar: (p) => dinero(p.interest_cents), alinear: 'derecha' }, { id: 'saldo', nombre: c.saldo, valor: (p) => p.balance_after_cents ?? null, pintar: (p) => (p.balance_after_cents == null ? t.comun.sinDato : dinero(p.balance_after_cents)), alinear: 'derecha' }, { id: 'status', nombre: t.comun.estado, valor: (p) => ESTADO_PARTE[p.status ?? 'pending'] }] }} /></Seccion>
}

export default function P18Credito({ params, sub, consulta }: PropsPantalla) {
  const id = Number(params.id)
  const linea = usarConsulta('get', '/api/v1/credits/{credit_id}', { ruta: { credit_id: id } })
  const pagos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { credit_id: [id], limit: 200 } })
  const compras = usarConsulta('get', '/api/v1/credit-purchases', { consulta: { credit_id: [id], limit: 200 } }, sub === 'detalle')
  const eventos = usarConsulta('get', '/api/v1/calendar/events', { consulta: { period: periodoLinea(pagos.datos?.items ?? []), credit_id: [id], event_type: TIPOS_LINEA, limit: 200 } }, sub === 'calendario' && !pagos.cargando)
  const rev = linea.datos?.kind === 'revolving_line'
  const serie = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-series', { ruta: { credit_id: id } }, sub === 'tendencias' && rev)
  const trend = usarConsulta('get', '/api/v1/credits/debt-trend', undefined, sub === 'tendencias')
  const per = usarPersonas()
  useEffect(() => { if (consulta.get('hoja') === 'pago') pagos.recargar() }, [consulta.get('hoja')]) // eslint-disable-line react-hooks/exhaustive-deps
  const recargar = () => { pagos.recargar(); linea.recargar(); eventos.recargar() }
  const siguiente = enAgenda(pagos.datos?.items ?? []).find((p) => p.status !== 'paid')
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <Barra>
        <Boton variante="secundario" onClick={() => navegar(rutaCondiciones(id))}><Settings2 aria-hidden className="size-5" />{c.editarCondiciones}</Boton>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-20') + `?credito=${id}`)}><Calculator aria-hidden className="size-5" />{c.simular}</Boton>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-19', { id, compra: 'nueva' }))}><Plus aria-hidden className="size-5" />{c.registrarCompraCorto}</Boton>
        <Boton variante="primario" deshabilitadoPor={siguiente ? undefined : c.sinPagosPendientes} onClick={() => siguiente && navegar(conConsulta({ hoja: 'pago', pago: String(siguiente.id) }))}><CircleCheck aria-hidden className="size-5" />{c.marcarPago}</Boton>
      </Barra>
      <ConDatos consulta={linea}>{(l) => {
        // Mientras cargan los pagos (dan el inicio del rango), la consulta de eventos espera: se muestra «Cargando», no un error.
        if (sub === 'calendario') return <ConDatos consulta={{ ...eventos, cargando: eventos.cargando || pagos.cargando }}>{(ev) => <CalendarioLinea eventos={itemsDe<Evento>(ev)} pagos={pagos.datos?.items ?? []} per={per} id={id} />}</ConDatos>
        if (sub === 'tendencias') return <div className="flex flex-col gap-6">
          {rev && <ConDatos consulta={serie}>{(s) => <Seccion titulo={c.serie}>
            <p className="text-small tabular-nums text-text-secondary">{c.minimos(dinero(s.min_a_cents), dinero(s.min_b_cents))}</p>
            <Grafica titulo={c.serie} resumen={c.resumenSerie} tipo="lineas" categorias={s.events.map((e) => fecha(e.event_date))} series={[{ nombre: c.lectura('A'), valores: s.events.map((e) => e.available_a_cents) }, { nombre: c.lectura('B'), valores: s.events.map((e) => e.available_b_cents) }]} />
            <VistaBD def={{ clave: 'P-18-serie', filas: s.events.map((e, i) => ({ ...e, i })), idFila: (e) => e.i, vistas: ['tabla'], titulo: (e) => e.description, propiedades: [{ id: 'f', nombre: t.comun.fecha, valor: (e) => e.event_date, pintar: (e) => fecha(e.event_date) }, { id: 't', nombre: c.evento, valor: (e) => EVENTO[e.event_type] }, { id: 'd', nombre: c.comercio, valor: (e) => e.description }, { id: 'm', nombre: t.comun.monto, valor: (e) => e.amount_cents, pintar: (e) => dinero(e.amount_cents), alinear: 'derecha' }, { id: 'c', nombre: c.capitalRepuesto, valor: (e) => e.principal_restored_cents, pintar: (e) => dinero(e.principal_restored_cents), alinear: 'derecha' }, { id: 'n', nombre: c.interesNoRepuesto, valor: (e) => e.interest_not_restored_cents, pintar: (e) => dinero(e.interest_not_restored_cents), alinear: 'derecha' }, { id: 'a', nombre: c.lectura('A'), valor: (e) => e.available_a_cents, pintar: (e) => dinero(e.available_a_cents), alinear: 'derecha' }, { id: 'b', nombre: c.lectura('B'), valor: (e) => e.available_b_cents, pintar: (e) => dinero(e.available_b_cents), alinear: 'derecha' }] }} />
          </Seccion>}</ConDatos>}
          <ConDatos consulta={trend}>{(tr) => <Grafica titulo={c.tendenciaDeuda} resumen={c.resumenTendencia} tipo="lineas" categorias={tr.months.map((m) => mes(m.month))} series={[{ nombre: l.name, valores: tr.months.map((m) => m.credits.find((x) => x.credit_id === id)?.live_debt_cents ?? 0) }]} />}</ConDatos>
        </div>
        if (sub === 'detalle') return <div className="flex flex-col gap-6">
          {l.kind !== 'term_loan' && <Seccion titulo={c.comprasLinea}><ConDatos consulta={compras}>{(d) => <TablaCompras compras={d.items} per={per} l={l} />}</ConDatos></Seccion>}
          {(l.kind === 'term_loan' || l.kind === 'store_credit') && <ConDatos consulta={pagos}>{(pg) => <TablaPagos pagos={pg.items} titulo={l.kind === 'term_loan' ? c.amortizacion : c.abonos} />}</ConDatos>}
        </div>
        const prox = siguiente
        const partes = prox ? textoPartes(agrupar(prox.shares, (s) => s.amount_cents, per)) : undefined
        return <div className="flex flex-col gap-6">
          <p className="text-small text-text-secondary">{TIPO_CREDITO[l.kind]}{l.lender ? ` · ${l.lender}` : ''}</p>
          <FranjaLinea l={l} />
          <div className="grid grid-cols-1 gap-3 tablet:grid-cols-2 desktop:grid-cols-3">
            <TarjetaMetrica titulo={c.deudaViva} valor={l.live_debt_cents == null ? t.comun.sinDato : dinero(l.live_debt_cents)} rol="debt" />
            <TarjetaMetrica titulo={c.proximoPago} valor={prox ? dinero(prox.amount_cents) : c.sinProximo} nota={prox && `${fechaConDia(prox.due_date)}${partes ? ` · ${partes}` : ''}`} alAbrir={prox ? () => navegar(conConsulta({ sub: 'calendario' })) : undefined} />
            <TarjetaMetrica titulo={c.comprasVivas} valor={String(l.live_purchase_count ?? 0)} nota={!l.live_purchase_count ? (l.kind === 'revolving_line' ? c.sinCompras : c.sinPagosPendientes) : undefined} />
            {l.kind === 'bank_card' && <TarjetaMetrica titulo={c.disponible} valor={l.available_cents == null ? t.comun.sinDato : dinero(l.available_cents)} rol="savings" />}
            {conFaltantes(l) && <TarjetaMetrica titulo={c.pagoMinimo} valor="—"><FaltaCapturar credito={l.id} que={l.missing_terms} /></TarjetaMetrica>}
          </div>
          {rev && <ResumenRevolvente l={l} />}
          {l.kind === 'bank_card' && <EstadoTarjeta l={l} />}
          <DatosLinea l={l} />
        </div>
      }}</ConDatos>
      <HojaPago pagos={pagos.datos?.items ?? []} per={per} alCambiar={recargar} consulta={consulta} />
    </div>
  )
}
