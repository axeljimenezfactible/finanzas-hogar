// P-26.08 Reportes y frecuencias (U4 parte b; decisión F, K47, K50, Q2-66 B). Día y hora: claves reports.*
// con getReportSettings/updateReportSettings; «Activo» con updateNotification. Próximo envío al guardar:
// next_run_at de listJobs (PROVISIONAL, bloqueo C1-hora-reportes-job).
import { useRef, useState } from 'react'
import { TriangleAlert } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { conConsulta, Enlace, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { pedir, nuevaClave } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { fechaConDia, fechaHora, hoyISO, horaCorta, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ChipSeleccionable, MarcaPersona } from '../../componentes/Chip'
import { Hoja } from '../../componentes/Hoja'
import { Esqueleto } from '../../componentes/Estados'
import { useAvisos } from '../../componentes/Avisos'
import { estilosControl } from '../../componentes/Campo'
import { Bloque, Cargar, FilaParametro, IndicadorGuardado, Interruptor, Marco, unir, useEsEscritorio, useGuardado, useSinRed } from './comun'
import { ejemploMesCorto, proximoEnvio, textoEnvio, type Programa } from './calculos'
import { textosAjustesB as T } from './textos'

const t = T.r
type N = Esquemas['Notification']; type J = Esquemas['Job']; type R = Esquemas['NotificationRecipient']; type Pe = Esquemas['Person']
type Ajustes = Esquemas['ReportSettings']['parameters']
const CODIGOS = ['daily_summary', 'weekly_summary', 'monthly_report'] as const
const CLAVES = {
  daily_summary: { hora: 'reports.daily_send_time' as const, dia: null },
  weekly_summary: { hora: 'reports.weekly_send_time' as const, dia: 'reports.weekly_weekday' as const },
  monthly_report: { hora: 'reports.monthly_send_time' as const, dia: 'reports.monthly_day_of_month' as const },
}
const guardarAjuste = (parameters: Esquemas['ReportSettingsUpdate']['parameters'], reset?: Esquemas['ReportSettingsUpdate']['reset_parameters']) =>
  pedir('patch', '/api/v1/settings/reports', { cuerpo: reset ? { reset_parameters: reset } : { parameters } })
const programa = (code: string, dia: number): Programa => (code === 'daily_summary' ? { tipo: 'diario' } : code === 'weekly_summary' ? { tipo: 'semanal', dia } : { tipo: 'mensual', dia })
const envioDeJob = (j?: J) => (j?.next_run_at ? `${fechaConDia(j.next_run_at.slice(0, 10))}, ${horaCorta(new Date(j.next_run_at))}` : '—')
export const nombrePersona = (ps: Pe[], id: number) => ps.find((p) => p.id === id)?.display_name ?? ''

export default function P2608Reportes({ consulta }: PropsPantalla) {
  const datos = unir<[{ items: N[] }, { items: J[] }, Esquemas['ReportSettings'], { items: R[] }, { items: Pe[] }, Esquemas['ReportMonths']]>([
    usarConsulta('get', '/api/v1/notifications', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/jobs', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/settings/reports'),
    usarConsulta('get', '/api/v1/notification-recipients', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/persons', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/monthly-report/months'),
  ])
  const esqueleto = <div className="flex flex-col gap-3">{[0, 1, 2].map((i) => <Esqueleto key={i} className="h-16" />)}</div>
  return (
    <Marco id="08">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Cargar c={datos} que={t.que} esqueleto={esqueleto}>{([ns, js, aj, rs, ps, meses]) => {
        const filas = CODIGOS.map((code) => ns.items.find((n) => n.code === code))
        return (
          <>
            <Bloque titulo={t.titulo} ayuda={t.reglas}>
              <TablaReportes filas={filas} jobs={js.items} ajustes={aj.parameters} recips={rs.items} personas={ps.items} recargar={datos.recargar} />
            </Bloque>
            <BajoDemanda consulta={consulta} meses={meses.months} personas={ps.items} />
          </>
        )
      }}</Cargar>
    </Marco>
  )
}

function TablaReportes({ filas, jobs, ajustes, recips, personas, recargar }: { filas: (N | undefined)[]; jobs: J[]; ajustes: Ajustes; recips: R[]; personas: Pe[]; recargar: () => void }) {
  const escritorio = useEsEscritorio()
  const props = (n: N) => ({ n, job: jobs.find((j) => j.notification_id === n.id), ajustes, recips, personas, recargar })
  if (!escritorio) return <div className="flex flex-col gap-3">{filas.map((n, i) => (n ? <TarjetaReporte key={n.id} {...props(n)} /> : <Falta key={i} />))}</div>
  return (
    <table className="w-full text-left text-small">
      <thead className="text-caption text-text-secondary"><tr>{[t.reporte, t.activo, t.dia, t.hora, t.destinatarios, t.proximo, t.ultimo].map((c) => <th key={c} scope="col" className="border-b border-border-divider px-2 py-2 font-medium">{c}</th>)}</tr></thead>
      <tbody>{filas.map((n, i) => (n ? <FilaReporte key={n.id} {...props(n)} /> : <tr key={i}><td colSpan={7}><Falta /></td></tr>))}</tbody>
    </table>
  )
}
const Falta = () => <p role="alert" className="text-small text-role-error">{t.faltaReporte}</p>
type PropsFila = { n: N; job?: J; ajustes: Ajustes; recips: R[]; personas: Pe[]; recargar: () => void }

function useFila({ n, job, ajustes, recargar }: PropsFila) {
  const g = useGuardado()
  const claves = CLAVES[n.code as keyof typeof CLAVES]
  const dia = claves.dia ? (ajustes[claves.dia].value ?? 1) : 1
  const hora = ajustes[claves.hora].value ?? ''
  const previaDe = (d: number, h: string) => {
    const r = proximoEnvio(programa(n.code, d), h)
    const nota = !n.is_enabled ? t.apagado : n.code === 'monthly_report' && d >= 29 ? t.mesCorto(d, ejemploMesCorto(d)) : undefined
    return { cambios: [{ etiqueta: t.proximo, antes: envioDeJob(job), despues: textoEnvio(r) }], nota }
  }
  const activo = (v: boolean) => g.ejecutar(async () => { await pedir('patch', '/api/v1/notifications/{notification_id}', { ruta: { notification_id: n.id }, cuerpo: { is_enabled: v } }); recargar() })
  const editorDia = (compacto: boolean) => claves.dia && (
    <FilaParametro compacto={compacto} etiqueta={n.code === 'weekly_summary' ? t.diaSemana : t.diaMes} tipo={n.code === 'weekly_summary' ? 'diaSemana' : 'diaMes'} param={ajustes[claves.dia]}
      mensaje={n.code === 'weekly_summary' ? T.validar.diaSemana : T.validar.diaMes} vistaPrevia={(d: number) => previaDe(d, hora)}
      guardar={async (v: number) => { await guardarAjuste({ [claves.dia!]: v }); recargar() }} restablecer={async () => { await guardarAjuste(undefined, [claves.dia!]); recargar() }} />
  )
  const editorHora = (compacto: boolean) => (
    <FilaParametro compacto={compacto} etiqueta={t.hora} tipo="hora" param={ajustes[claves.hora]} mensaje={T.validar.hora} ayuda={t.pistas[n.code]} vistaPrevia={(h: string) => previaDe(dia, h)}
      guardar={async (v: string) => { await guardarAjuste({ [claves.hora]: v }); recargar() }} restablecer={async () => { await guardarAjuste(undefined, [claves.hora]); recargar() }} />
  )
  return { g, activo, editorDia, editorHora }
}

export function Destinatarios({ n, recips, personas, recargar }: Pick<PropsFila, 'n' | 'recips' | 'personas' | 'recargar'>) {
  const g = useGuardado()
  const sinRed = useSinRed()
  const elegidos = n.recipient_ids ?? []
  const conCuenta = personas.filter((p) => p.participation !== 'external')
  const cambiar = (id: number) => g.ejecutar(async () => {
    const ids = elegidos.includes(id) ? elegidos.filter((x) => x !== id) : [...elegidos, id]
    await pedir('put', '/api/v1/notifications/{notification_id}/recipients', { ruta: { notification_id: n.id }, cuerpo: { ids } })
    recargar()
  })
  return (
    <div className="flex flex-col gap-1">
      <div role="group" aria-label={`${t.destinatarios} · ${n.name}`} className="flex flex-wrap gap-2">
        {conCuenta.map((p) => {
          const r = recips.find((x) => x.person_id === p.id && x.is_registered !== false && x.verified_at)
          return (
            <ChipSeleccionable key={p.id} marca={<MarcaPersona codigo={p.code} />} seleccionado={!!r && elegidos.includes(r.id)} aria-pressed={!!r && elegidos.includes(r.id)}
              aria-disabled={!r || sinRed ? true : undefined} title={!r ? t.soloConfirmados : sinRed} onClick={() => r && !sinRed && void cambiar(r.id)}>{p.display_name}</ChipSeleccionable>
          )
        })}
      </div>
      {!elegidos.length && <span className="inline-flex items-center gap-1 self-start rounded-md bg-role-warning-surface px-2 py-1 text-caption text-text-default"><TriangleAlert aria-hidden className="size-5 text-role-warning" />{t.sinDestinatarios}</span>}
      <IndicadorGuardado g={g} />
    </div>
  )
}
function Ultimo({ job }: { job?: J }) {
  if (!job?.last_status) return <span className="text-text-secondary">{t.sinEnvios}</span>
  const r = t.resultado[job.last_status] ?? job.last_status
  return <span>{job.last_status === 'error' && job.last_error ? `${r}: ${job.last_error}` : r}{job.last_run_at && <span className="block text-caption text-text-secondary">{fechaHora(job.last_run_at)}</span>}</span>
}
function EnviarAhora({ n, recips, personas }: PropsFila) {
  const g = useGuardado()
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const confirmadas = personas.filter((p) => recips.some((r) => r.person_id === p.id && r.verified_at))
  return (
    <div className="flex flex-wrap items-center gap-2">
      {confirmadas.map((p) => (
        <Boton key={p.id} variante="fantasma" deshabilitadoPor={sinRed} onClick={() => void g.ejecutar(async () => { await pedir('post', '/api/v1/notifications/{notification_id}/test', { ruta: { notification_id: n.id }, cuerpo: { person_id: p.id } }); avisar({ texto: t.enviadoA(p.display_name), rol: 'exito' }) })}>
          {`${t.enviarAhora.replace('…', '')} ${p.display_name}`}
        </Boton>
      ))}
      <IndicadorGuardado g={g} />
    </div>
  )
}
function FilaReporte(p: PropsFila) {
  const f = useFila(p)
  const sinRed = useSinRed()
  return (
    <tr className="align-top" data-prueba={`reporte-${p.n.code}`}>
      <th scope="row" className="border-b border-border-divider px-2 py-3 font-medium"><span>{p.n.name}</span><span className="block text-caption font-normal text-text-secondary">{t.frecuencias[p.n.code]}</span></th>
      <td className="border-b border-border-divider px-2 py-3"><Interruptor ocultarEtiqueta etiqueta={`${t.activo} · ${p.n.name}`} activo={p.n.is_enabled} deshabilitadoPor={sinRed} alCambiar={(v) => void f.activo(v)} /><IndicadorGuardado g={f.g} /></td>
      <td className="border-b border-border-divider px-2 py-3">{f.editorDia(true) || '—'}</td>
      <td className="border-b border-border-divider px-2 py-3">{f.editorHora(true)}</td>
      <td className="border-b border-border-divider px-2 py-3"><Destinatarios {...p} /><EnviarAhora {...p} /></td>
      <td className="border-b border-border-divider px-2 py-3 tabular-nums">{envioDeJob(p.job)}{!p.n.is_enabled && <span className="block text-caption text-text-secondary">{t.apagado}</span>}</td>
      <td className="border-b border-border-divider px-2 py-3"><Ultimo job={p.job} /></td>
    </tr>
  )
}
function TarjetaReporte(p: PropsFila) {
  const f = useFila(p)
  const sinRed = useSinRed()
  const [abierta, setAbierta] = useState(false)
  return (
    <article className="flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-3" data-prueba={`reporte-${p.n.code}`}>
      <div className="flex items-center justify-between gap-2">
        <h4 className="text-body font-semibold">{p.n.name}</h4>
        <Interruptor ocultarEtiqueta etiqueta={`${t.activo} · ${p.n.name}`} activo={p.n.is_enabled} deshabilitadoPor={sinRed} alCambiar={(v) => void f.activo(v)} />
      </div>
      <p className="text-small tabular-nums text-text-secondary">{t.proximoEnvio(envioDeJob(p.job))}</p>
      {!p.n.is_enabled && <p className="text-caption text-text-secondary">{t.apagado}</p>}
      {!(p.n.recipient_ids ?? []).length && <span className="inline-flex items-center gap-1 self-start rounded-md bg-role-warning-surface px-2 py-1 text-caption"><TriangleAlert aria-hidden className="size-5 text-role-warning" />{t.sinDestinatarios}</span>}
      <IndicadorGuardado g={f.g} />
      <Boton tam="hoja" onClick={() => setAbierta(true)}>{t.editarFila}</Boton>
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={p.n.name} descripcion={t.frecuencias[p.n.code]}>
        <div className="flex flex-col gap-2">
          {f.editorDia(false)}
          {f.editorHora(false)}
          <p className="text-small font-medium">{t.destinatarios}</p>
          <Destinatarios {...p} />
          <EnviarAhora {...p} />
          <p className="text-small font-medium">{t.ultimo}</p>
          <Ultimo job={p.job} />
        </div>
      </Hoja>
    </article>
  )
}

function BajoDemanda({ consulta, meses, personas }: { consulta: URLSearchParams; meses: Esquemas['ReportMonths']['months']; personas: Pe[] }) {
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const hoja = consulta.get('hoja')
  const abrir = (h: string | null) => navegar(conConsulta({ hoja: h }))
  const ultimo = meses[meses.length - 1]?.month
  const [mesElegido, setMes] = useState(ultimo ?? '')
  const [telegram, setTelegram] = useState(false)
  const hoy = hoyISO()
  const [inicio, setInicio] = useState(`${hoy.slice(0, 7)}-01`)
  const [fin, setFin] = useState(hoy)
  const clave = useRef(nuevaClave())
  const g = useGuardado()
  const karen = personas.find((p) => p.participation === 'external')
  const fechasMal = fin < inicio
  const generar = () => g.ejecutar(async () => {
    await pedir('post', '/api/v1/reports/monthly', { cuerpo: { month: mesElegido, send_telegram: telegram }, idempotencia: clave.current })
    clave.current = nuevaClave(); abrir(null)
    avisar({ texto: t.enPreparacion, rol: 'exito', accion: { texto: t.verCentro, alHacer: () => navegar(rutaDe('P-23')) } })
  })
  const estado = () => karen && !fechasMal && g.ejecutar(async () => {
    await pedir('post', '/api/v1/reports/external-statements', { cuerpo: { person_id: karen.id, period_start_date: inicio, period_end_date: fin }, idempotencia: clave.current })
    clave.current = nuevaClave(); navegar(rutaDe('P-25'))
  })
  return (
    <Bloque titulo={t.bajoDemanda}>
      <div className="flex flex-col gap-2 tablet:flex-row tablet:flex-wrap">
        <Boton variante="primario" tam="hoja" className="tablet:w-auto" deshabilitadoPor={sinRed} onClick={() => abrir('generar')}>{t.generarMensual}</Boton>
        <Boton tam="hoja" className="tablet:w-auto" deshabilitadoPor={sinRed} onClick={() => abrir('estado-karen')}>{t.estadoKaren}</Boton>
        <Enlace a={rutaDe('P-24')} className="inline-flex h-12 items-center justify-center rounded-md px-4 text-body text-text-link underline-offset-2 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">{ultimo ? `${t.vistaWeb} · ${mes(ultimo)}` : t.vistaWeb}</Enlace>
      </div>
      <Hoja abierta={hoja === 'generar'} alCambiar={(a) => !a && abrir(null)} titulo={t.generarMensual}
        pie={<><Boton variante="primario" tam="hoja" deshabilitadoPor={sinRed} onClick={() => void generar()}>{t.generar}</Boton><IndicadorGuardado g={g} /></>}>
        <label className="flex flex-col gap-1 text-small font-medium">{t.mes}
          <select className={estilosControl} value={mesElegido} onChange={(e) => setMes(e.target.value)}>{meses.map((m) => <option key={m.month} value={m.month}>{mes(m.month)}</option>)}</select>
        </label>
        <label className="mt-3 flex min-h-11 items-center gap-2 text-small"><input type="checkbox" className="size-5 accent-action-primary-bg" checked={telegram} onChange={(e) => setTelegram(e.target.checked)} />{t.enviarTelegram}</label>
      </Hoja>
      <Hoja abierta={hoja === 'estado-karen'} alCambiar={(a) => !a && abrir(null)} titulo={t.estadoKaren}
        pie={<><Boton variante="primario" tam="hoja" deshabilitadoPor={sinRed ?? (fechasMal ? T.validar.fechas : undefined)} onClick={() => void estado()}>{t.generar}</Boton><IndicadorGuardado g={g} /></>}>
        <div className="flex flex-col gap-3">
          <Campo etiqueta={t.inicio} type="date" value={inicio} onChange={(e) => setInicio(e.target.value)} />
          <Campo etiqueta={t.fin} type="date" value={fin} onChange={(e) => setFin(e.target.value)} error={fechasMal ? T.validar.fechas : undefined} />
        </div>
      </Hoja>
    </Bloque>
  )
}
