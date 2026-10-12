// P-23 Centro de reportes (U1 §2.1, §2.2, §3.25: N-245..N-250). Reportes generados, próximos envíos y avisos enviados.
// Generar el mensual: hoja ?hoja=generar (también desde el menú «/»). El envío por Telegram lo hace el servidor.
import { useState } from 'react'
import { FileText, Globe, UserRound } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { Enlace, conConsulta, navegar } from '../../app/enrutador'
import { usarConsulta } from '../../api/consulta'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { fechaHora, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { useAvisos } from '../../componentes/Avisos'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { AccionesPdf, Barra, Datos, Seccion, Selector, itemsDe, mensajeError, periodoTexto } from './comun'
import { AVISO, ENTREGA, ENVIO, ESTADO_TRABAJO, TIPO_REPORTE, t } from './textos'

const r = t.reportes
type Reporte = Esquemas['GeneratedReport']
type Entrega = Esquemas['NotificationDelivery']
export const nombrePdf = (x: Reporte) => `${x.kind === 'monthly_pdf' ? 'reporte-mensual' : 'estado-de-cuenta-karen'}-${x.period_start_date.slice(0, 7)}.pdf`
/** Mes por defecto (N-256): el último cerrado; si no hay, el más reciente. */
export const mesPorDefecto = (ms: { month: string; is_partial: boolean }[]) => [...ms].sort((a, b) => b.month.localeCompare(a.month)).find((m) => !m.is_partial)?.month ?? ms[0]?.month ?? ''

export default function P23Centro({ consulta }: PropsPantalla) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const sin = enLinea ? undefined : textos.conexion.requiere
  const reportes = usarConsulta('get', '/api/v1/reports', { consulta: { limit: 200 } })
  const trabajos = usarConsulta('get', '/api/v1/jobs', { consulta: { limit: 200 } })
  const avisos = usarConsulta('get', '/api/v1/notification-deliveries', { consulta: { limit: 50 } })
  const reintentar = async (x: Reporte) => {
    try { await pedir('post', '/api/v1/reports/{report_id}/send', { ruta: { report_id: x.id }, idempotencia: nuevaClave() }); avisar({ texto: r.reenviando, rol: 'info' }); reportes.recargar() } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <Barra>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-25'))}><UserRound aria-hidden className="size-4" />{r.karen}</Boton>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-24'))}><Globe aria-hidden className="size-4" />{r.verWeb}</Boton>
        <Boton variante="primario" deshabilitadoPor={sin} onClick={() => navegar(conConsulta({ hoja: 'generar' }))}><FileText aria-hidden className="size-4" />{r.generarMensual}</Boton>
      </Barra>
      <Seccion titulo={r.proximos} accion={<Enlace a={rutaDe('P-26', { seccion: '08' })} className="inline-flex min-h-11 items-center text-small text-text-link underline">{r.cambiarFrecuencias}</Enlace>}>
        <ConDatos consulta={trabajos}>{(d) => { const js = itemsDe<Esquemas['Job']>(d).filter((j) => j.is_enabled && ['daily_summary', 'weekly_summary', 'monthly_report'].includes(j.code))
          return js.length ? <Datos filas={js.map((j) => [AVISO[j.code] ?? j.name, `${j.next_run_at ? fechaHora(j.next_run_at) : t.comun.sinDato}${j.last_status ? ` · ${r.ultimo.toLowerCase()}: ${ESTADO_TRABAJO[j.last_status].toLowerCase()}` : ''}`])} /> : <p className="text-small text-text-secondary">{r.sinProximos}</p> }}</ConDatos>
      </Seccion>
      <Seccion titulo={r.generados}>
        <ConDatos consulta={reportes}>{(d) => <VistaBD def={{ clave: 'P-23-reportes', filas: itemsDe<Reporte>(d), idFila: (x) => x.id, vistas: ['tabla', 'calendario'], fecha: (x) => x.generated_at.slice(0, 10), vacio: <Vacio texto={r.vacio} />,
          titulo: (x) => `${TIPO_REPORTE[x.kind]} · ${periodoTexto(x.period_start_date, x.period_end_date)}`,
          alAbrir: (x) => navegar(x.kind === 'monthly_pdf' ? `${rutaDe('P-24')}?mes=${x.period_start_date.slice(0, 7)}` : rutaDe('P-25')),
          propiedades: [
            { id: 'kind', nombre: r.tipo, valor: (x) => TIPO_REPORTE[x.kind], grupo: (x) => TIPO_REPORTE[x.kind] },
            { id: 'periodo', nombre: r.periodo, valor: (x) => x.period_start_date, pintar: (x) => periodoTexto(x.period_start_date, x.period_end_date) },
            { id: 'generado', nombre: r.generado, valor: (x) => x.generated_at, pintar: (x) => fechaHora(x.generated_at) },
            { id: 'envio', nombre: r.envio, valor: (x) => x.delivery_status ?? 'not_sent', pintar: (x) => x.delivery_status === 'failed'
              ? <span className="inline-flex flex-wrap items-center gap-2 text-role-error">{ENVIO.failed}<Boton variante="enlace" deshabilitadoPor={sin} onClick={(e) => { e.stopPropagation(); void reintentar(x) }}>{r.reintentarEnvio}</Boton></span>
              : x.delivery_status === 'sent' && x.sent_at ? `${ENVIO.sent} · ${fechaHora(x.sent_at)}` : ENVIO[x.delivery_status ?? 'not_sent'] },
            { id: 'pdf', nombre: t.comun.descargar, valor: () => '', pintar: (x) => <span onClick={(e) => e.stopPropagation()}><AccionesPdf id={x.id} nombre={nombrePdf(x)} /></span> },
          ] }} />}</ConDatos>
      </Seccion>
      <Seccion titulo={r.avisos}>
        <ConDatos consulta={avisos}>{(d) => <VistaBD def={{ clave: 'P-23-avisos', filas: itemsDe<Entrega>(d), idFila: (x) => x.id, vistas: ['tabla', 'calendario'], fecha: (x) => (x.sent_at ?? x.created_at).slice(0, 10), vacio: <Vacio texto={r.sinAvisos} />, titulo: (x) => AVISO[x.notification_code ?? ''] ?? r.avisos,
          propiedades: [
            { id: 'fecha', nombre: r.enviado, valor: (x) => x.sent_at ?? x.created_at, pintar: (x) => fechaHora(x.sent_at ?? x.created_at) },
            { id: 'tipo', nombre: r.tipo, valor: (x) => AVISO[x.notification_code ?? ''] ?? t.comun.sinDato, grupo: (x) => AVISO[x.notification_code ?? ''] ?? t.comun.sinDato },
            { id: 'resultado', nombre: r.resultado, valor: (x) => ENTREGA[x.status], grupo: (x) => ENTREGA[x.status] },
            { id: 'texto', nombre: r.texto, valor: (x) => x.message_text ?? '' },
          ] }} />}</ConDatos>
      </Seccion>
      <HojaGenerar abierta={consulta.get('hoja') === 'generar'} alTerminar={reportes.recargar} />
    </div>
  )
}

function HojaGenerar({ abierta, alTerminar }: { abierta: boolean; alTerminar: () => void }) {
  const avisar = useAvisos()
  const meses = usarConsulta('get', '/api/v1/monthly-report/months', undefined, abierta)
  const [elegido, setElegido] = useState('')
  const [telegram, setTelegram] = useState(false)
  const [clave, setClave] = useState(nuevaClave)
  const ms = meses.datos?.months ?? []
  const valor = elegido || mesPorDefecto(ms)
  const cerrar = () => navegar(conConsulta({ hoja: null }))
  const generar = async () => {
    if (!valor) return
    try { await pedir('post', '/api/v1/reports/monthly', { cuerpo: { month: valor, send_telegram: telegram }, idempotencia: clave }); setClave(nuevaClave()); avisar({ texto: r.generandoAviso, rol: 'info' }); cerrar(); alTerminar() } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  return (
    <Hoja abierta={abierta} alCambiar={(a) => !a && cerrar()} titulo={r.generarMensual} pie={<Boton variante="primario" tam="hoja" deshabilitadoPor={valor ? undefined : t.mensual.sinMeses} onClick={generar}>{r.generar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Selector etiqueta={t.comun.mes} valor={valor} alCambiar={setElegido} opciones={[...ms].sort((a, b) => b.month.localeCompare(a.month)).map((m) => ({ valor: m.month, texto: `${mes(m.month)}${m.is_partial ? ` (${t.comun.parcial})` : ''}` }))} />
        <label className="flex min-h-11 items-center gap-3 text-body"><input type="checkbox" className="size-5 accent-action-primary-bg" checked={telegram} onChange={(e) => setTelegram(e.target.checked)} />{r.enviarTelegram}</label>
      </div>
    </Hoja>
  )
}
