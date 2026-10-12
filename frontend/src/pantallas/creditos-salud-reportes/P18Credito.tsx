// P-18 Crédito (U1 §2.1, §2.2, §3.20: N-200..N-210). Variante por tipo: revolvente (A/B, serie, tasas por plazo),
// tarjeta (estado del ciclo), préstamo (amortización) y tienda (abonos). Sin valores supuestos: «Falta capturar» (K18).
import { useEffect, type ReactNode } from 'react'
import { Calculator, CircleCheck, Plus, Settings2 } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { Enlace, conConsulta, navegar } from '../../app/enrutador'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { itemsDe, Agenda, Barra, Datos, FaltaCapturar, HojaPago, Seccion, agrupar, enAgenda, rutaCondiciones, textoPartes, usarPersonas, type Pago, type Personas } from './comun'
import { FranjaAB, conFaltantes } from './P17Lineas'
import { ESTADO_COMPRA, ESTADO_PARTE, EVENTO, FRECUENCIA, TIPO_CREDITO, plazo, t } from './textos'

type Credito = Esquemas['Credit']
type Compra = Esquemas['CreditPurchase']
const c = t.creditos

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
  const rev = linea.datos?.kind === 'revolving_line'
  const serie = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-series', { ruta: { credit_id: id } }, sub === 'tendencias' && rev)
  const trend = usarConsulta('get', '/api/v1/credits/debt-trend', undefined, sub === 'tendencias')
  const per = usarPersonas()
  useEffect(() => { if (consulta.get('hoja') === 'pago') pagos.recargar() }, [consulta.get('hoja')]) // eslint-disable-line react-hooks/exhaustive-deps
  const recargar = () => { pagos.recargar(); linea.recargar() }
  const siguiente = enAgenda(pagos.datos?.items ?? []).find((p) => p.status !== 'paid')
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <Barra>
        <Boton variante="secundario" onClick={() => navegar(rutaCondiciones(id))}><Settings2 aria-hidden className="size-4" />{c.editarCondiciones}</Boton>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-20') + `?credito=${id}`)}><Calculator aria-hidden className="size-4" />{c.simular}</Boton>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-19', { id, compra: 'nueva' }))}><Plus aria-hidden className="size-4" />{c.registrarCompraCorto}</Boton>
        <Boton variante="primario" deshabilitadoPor={siguiente ? undefined : c.sinPagosPendientes} onClick={() => siguiente && navegar(conConsulta({ hoja: 'pago', pago: String(siguiente.id) }))}><CircleCheck aria-hidden className="size-4" />{c.marcarPago}</Boton>
      </Barra>
      <ConDatos consulta={linea}>{(l) => {
        const nombre = () => l.name
        if (sub === 'calendario') return <ConDatos consulta={pagos}>{(pg) => <Seccion titulo={c.agenda}><Agenda pagos={enAgenda(pg.items)} per={per} nombreLinea={nombre} /></Seccion>}</ConDatos>
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
          <FranjaAB lineas={[l]} />
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
