// P-25 Estado de cuenta de Karen (U1 §2.1, §2.2, §3.27: N-260..N-262). Vista previa del periodo, generar PDF y compartir
// (Axel lo reenvía a mano desde la hoja de compartir de Android; el sistema no se lo manda a Karen, decisión A).
import { useState } from 'react'
import type { PropsPantalla } from '../../app/registro'
import { usarConsulta } from '../../api/consulta'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { dinero, fecha, fechaHora, hoyISO } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ChipSeleccionable } from '../../componentes/Chip'
import { useAvisos } from '../../componentes/Avisos'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { AccionesPdf, Datos, Seccion, itemsDe, mensajeError, periodoTexto, usarPersonas } from './comun'
import { nombrePdf } from './P23Centro'
import { ESTADO_PARTE, plazo, t } from './textos'

const k = t.karen
type Generado = Esquemas['GeneratedReport']

export default function P25Karen(_: PropsPantalla) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const per = usarPersonas()
  const karen = per.externaPrincipal
  const [modo, setModo] = useState<'mes' | 'rango'>('mes')
  const [mesSel, setMesSel] = useState(hoyISO().slice(0, 7))
  const [desde, setDesde] = useState('')
  const [hasta, setHasta] = useState('')
  const periodo = modo === 'rango' && desde && hasta ? `${desde}..${hasta}` : mesSel
  const est = usarConsulta('get', '/api/v1/persons/{person_id}/statement', { ruta: { person_id: karen?.id ?? 0 }, consulta: { period: periodo } }, !!karen)
  const anteriores = usarConsulta('get', '/api/v1/reports', { consulta: { kind: 'external_statement_pdf', limit: 50 } })
  const [nuevo, setNuevo] = useState<Generado | null>(null)
  const [generando, setGenerando] = useState(false)
  const generar = async (a: string, b: string) => {
    if (!karen) return
    setGenerando(true)
    try { setNuevo(await pedir('post', '/api/v1/reports/external-statements', { cuerpo: { person_id: karen.id, period_start_date: a, period_end_date: b }, idempotencia: nuevaClave() })); avisar({ texto: t.comun.pdfListo, rol: 'exito' }); anteriores.recargar() } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) } finally { setGenerando(false) }
  }
  if (per.lista.length && !karen) return <Vacio texto={k.sinExterna} />
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <Seccion titulo={k.periodo}>
        <div className="flex flex-wrap gap-2">{([['mes', k.mesEnCurso], ['rango', k.rango]] as const).map(([v, n]) => <ChipSeleccionable key={v} seleccionado={modo === v} onClick={() => setModo(v)}>{n}</ChipSeleccionable>)}</div>
        {modo === 'mes' ? <Campo className="tablet:w-64" etiqueta={t.comun.mes} type="month" value={mesSel} onChange={(e) => e.target.value && setMesSel(e.target.value)} />
          : <div className="grid grid-cols-2 gap-3 tablet:w-96"><Campo etiqueta={t.comun.desde} type="date" value={desde} onChange={(e) => setDesde(e.target.value)} /><Campo etiqueta={t.comun.hasta} type="date" value={hasta} onChange={(e) => setHasta(e.target.value)} /></div>}
      </Seccion>
      <Franja rol="info">{k.nota}</Franja>
      <ConDatos consulta={est}>{(e) => {
        const vacio = !e.purchases.length && !e.grocery_shares.length && !e.payment_shares.length
        const pdf = nuevo && nuevo.period_start_date === e.period_start_date && nuevo.period_end_date === e.period_end_date ? nuevo : null
        return <div className="flex flex-col gap-6">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <p className="text-body font-medium">{periodoTexto(e.period_start_date, e.period_end_date)}</p>
            {pdf ? <AccionesPdf id={pdf.id} nombre={nombrePdf(pdf)} principal /> : <Boton variante="primario" deshabilitadoPor={generando ? t.comun.generando : enLinea ? undefined : textos.conexion.requiere} onClick={() => generar(e.period_start_date, e.period_end_date)}>{generando ? t.comun.generando : t.comun.generarPdf}</Boton>}
          </div>
          <div className="grid grid-cols-1 gap-3 tablet:grid-cols-2">
            <Datos filas={[[k.pagado, dinero(e.paid_in_period_cents)]]} />
            <Datos filas={[[k.pendiente, <span className="font-semibold text-role-debt">{dinero(e.pending_total_cents)}</span>]]} />
          </div>
          {vacio ? <Vacio texto={k.vacio} /> : <>
            <Seccion titulo={k.compras}>{e.purchases.length ? <Datos filas={e.purchases.map((c) => { const s = c.shares?.find((x) => x.person_id === karen?.id)
              return [`${c.description ?? t.comun.sinDato} · ${fecha(c.purchase_date)}`, `${dinero(c.principal_cents)} + ${dinero(c.interest_cents)} = ${dinero(c.total_cents)} · ${plazo(c.term_count, c.payment_frequency)}${s?.fortnight_payment_cents != null ? ` · ${k.pagoQuincena.toLowerCase()} ${dinero(s.fortnight_payment_cents)}` : ''}`] })} /> : <p className="text-small text-text-secondary">{t.creditos.sinCompras}</p>}</Seccion>
            <Seccion titulo={k.despensas}>{e.grocery_shares.length ? <Datos filas={e.grocery_shares.map((g) => [fecha(g.purchase_date), `${dinero(g.principal_cents)} + ${dinero(g.interest_cents)} = ${dinero(g.total_cents)}`])} /> : <p className="text-small text-text-secondary">{t.comun.sinDato}</p>}</Seccion>
            <Seccion titulo={k.pagos}>{e.payment_shares.length ? <Datos filas={e.payment_shares.map((p, i) => [`${i + 1}. ${ESTADO_PARTE[p.status]}`, `${dinero(p.amount_cents)} · ${k.pagadoParte.toLowerCase()} ${dinero(p.paid_cents)}`])} /> : <p className="text-small text-text-secondary">{t.comun.sinDato}</p>}</Seccion>
          </>}
        </div>
      }}</ConDatos>
      <Seccion titulo={k.anteriores}>
        <ConDatos consulta={anteriores}>{(d) => { const rs = itemsDe<Generado>(d)
          return rs.length ? <ul className="flex flex-col divide-y divide-border-divider rounded-lg border border-border-divider bg-surface-card">{rs.map((x) => <li key={x.id} className="flex flex-wrap items-center justify-between gap-2 px-4 py-3"><span className="text-small"><span className="font-medium">{periodoTexto(x.period_start_date, x.period_end_date)}</span> · {fechaHora(x.generated_at)}</span><AccionesPdf id={x.id} nombre={nombrePdf(x)} /></li>)}</ul> : <p className="text-small text-text-secondary">{k.sinAnteriores}</p> }}</ConDatos>
      </Seccion>
    </div>
  )
}
