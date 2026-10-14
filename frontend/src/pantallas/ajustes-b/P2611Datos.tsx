// P-26.11 Datos y respaldos (U4 parte b; decisión H, K31, K43, K44, K50, N-302..N-308, N-325).
// Hora del respaldo diario: updateJob (cron de 5 campos) PROVISIONAL hasta la puerta 4 (K31 punto 3).
import { textos } from '../../textos/es-MX'
import { SelectorArchivo } from '../../componentes/SelectorArchivo'
import { useEffect, useRef, useState } from 'react'
import { CircleAlert, CircleCheck, LoaderCircle } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { conConsulta, Enlace, navegar } from '../../app/enrutador'
import { pedir, nuevaClave } from '../../api/cliente'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { fechaHora } from '../../api/formatos'
import type { Cuerpo, Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { Hoja } from '../../componentes/Hoja'
import { Esqueleto, Vacio } from '../../componentes/Estados'
import { useAvisos } from '../../componentes/Avisos'
import { cn } from '../../componentes/utils'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { AntesDespues, Bloque, Cargar, DialogoCritico, FilaParametro, IndicadorGuardado, Interruptor, Marco, unir, useGuardado, useSinRed } from './comun'
import { esHora } from './calculos'
import { textosAjustesB as T } from './textos'

const t = T.d
type B = Esquemas['Backup']; type J = Esquemas['Job']; type AI = Esquemas['AiProvider']
const num1 = new Intl.NumberFormat('es-MX', { maximumFractionDigits: 1 })
export const tamano = (b?: number) => (b === undefined ? '—' : b >= 1024 ** 3 ? `${num1.format(b / 1024 ** 3)} GB` : `${num1.format(b / 1024 ** 2)} MB`)
const integridad = (b: B) => (!b.integrity_check ? '—' : b.integrity_check === 'ok' ? t.integro : b.integrity_check)
const guardarDatos = (c: Esquemas['DataSettingsUpdate']) => pedir('patch', '/api/v1/settings/data', { cuerpo: c })
const cronAHora = (s?: string) => { const m = /^(\d{1,2}) (\d{1,2}) \* \* \*$/.exec(s ?? ''); return m ? `${m[2].padStart(2, '0')}:${m[1].padStart(2, '0')}` : null }

export default function P2611Datos({ consulta }: PropsPantalla) {
  const respaldos = usarConsulta('get', '/api/v1/backups', { consulta: { limit: 50 } }) as unknown as EstadoConsulta<{ items: B[] }>
  const resto = unir<[{ items: J[] }, Esquemas['SeedStatus'], Esquemas['AiProviders'], Esquemas['SystemStatus'], Esquemas['DataSettings']]>([
    usarConsulta('get', '/api/v1/jobs', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/system/seed-status'),
    usarConsulta('get', '/api/v1/ai/providers'),
    usarConsulta('get', '/api/v1/system/status'),
    usarConsulta('get', '/api/v1/settings/data'),
  ])
  const enCurso = respaldos.datos?.items.some((b) => b.status === 'running')
  useEffect(() => { if (!enCurso) return; const i = setInterval(respaldos.recargar, 3000); return () => clearInterval(i) }, [enCurso, respaldos.recargar])
  const job = resto.datos?.[0].items.find((j) => j.code === 'backup_daily')
  return (
    <Marco id="11">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Respaldos c={respaldos} job={job} hoja={consulta.get('hoja')} />
      <Cargar c={resto} que={t.que} esqueleto={<Esqueleto className="h-40" />}>{([, seed, ia, sistema, datos]) => (
        <div className="grid grid-cols-1 gap-4 desktop:grid-cols-2">
          {job && <HoraRespaldo job={job} recargar={resto.recargar} />}
          <Importar />
          <Bloque titulo={t.semilla}>
            <ul className="flex flex-col gap-2">{seed.items.map((s) => (
              <li key={s.code} className="flex flex-wrap items-center justify-between gap-2 text-small"><span>{t.semillas[s.code]}</span>
                <span className="inline-flex items-center gap-1">{s.is_loaded ? <CircleCheck aria-hidden className="size-5 text-role-success" /> : <CircleAlert aria-hidden className="size-5 text-role-warning" />}{s.is_loaded ? t.cargada : t.pendiente} · {t.registros(s.row_count)}</span></li>
            ))}</ul>
          </Bloque>
          <Bloque titulo={t.ia} ayuda={t.iaAyuda}>{ia.providers.map((p) => <Proveedor key={p.provider} p={p} recargar={resto.recargar} />)}</Bloque>
          <Bloque titulo={t.servicio}>
            <dl className="grid grid-cols-2 gap-2 text-small">
              <dt className="text-text-secondary">{t.version}</dt><dd>{sistema.version}</dd>
              <dt className="text-text-secondary">{t.desde}</dt><dd className="tabular-nums">{fechaHora(sistema.started_at)}</dd>
              <dt className="text-text-secondary">{t.tamanoBase}</dt><dd className="tabular-nums">{tamano(sistema.database_size_bytes)}</dd>
            </dl>
          </Bloque>
          <Bloque titulo={t.parametros}>
            <FilaParametro etiqueta={t.adjunto} tipo="entero" unidad="MB" param={datos.parameters['attachment.max_size_mb']} mensaje={T.validar.mb} vistaPrevia={(v: number) => ({ cambios: [], nota: t.efectoAdjunto(v) })}
              guardar={async (v: number) => { await guardarDatos({ parameters: { 'attachment.max_size_mb': v } }); resto.recargar() }} restablecer={async () => { await guardarDatos({ reset_parameters: ['attachment.max_size_mb'] }); resto.recargar() }} />
            <FilaParametro etiqueta={t.retencion} tipo="entero" unidad="h" param={datos.parameters['idempotency.retention_hours']} mensaje={T.validar.horas} vistaPrevia={(v: number) => ({ cambios: [], nota: t.efectoRetencion(v) })}
              guardar={async (v: number) => { await guardarDatos({ parameters: { 'idempotency.retention_hours': v } }); resto.recargar() }} restablecer={async () => { await guardarDatos({ reset_parameters: ['idempotency.retention_hours'] }); resto.recargar() }} />
            <FilaParametro etiqueta={t.ocr} tipo="entero" unidad="s" param={datos.parameters['ocr.timeout_seconds']} mensaje={T.validar.segundos}
              vistaPrevia={(v: number) => ({ cambios: [{ etiqueta: t.efectoOcr, antes: `${datos.parameters['ocr.timeout_seconds'].value} s`, despues: `${v} s` }] })}
              guardar={async (v: number) => { await guardarDatos({ parameters: { 'ocr.timeout_seconds': v } }); resto.recargar() }} restablecer={async () => { await guardarDatos({ reset_parameters: ['ocr.timeout_seconds'] }); resto.recargar() }} />
          </Bloque>
        </div>
      )}</Cargar>
    </Marco>
  )
}

function Respaldos({ c, job, hoja }: { c: EstadoConsulta<{ items: B[] }>; job?: J; hoja: string | null }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const clave = useRef(nuevaClave())
  const [detalle, setDetalle] = useState<B | null>(null)
  const respaldar = () => g.ejecutar(async () => { await pedir('post', '/api/v1/backups', { idempotencia: clave.current }); clave.current = nuevaClave(); if (hoja) navegar(conConsulta({ hoja: null })); c.recargar() })
  const boton = <Boton variante="primario" tam="hoja" className="tablet:w-auto" deshabilitadoPor={sinRed} onClick={() => void respaldar()}>{t.respaldar}</Boton>
  return (
    <Bloque titulo={t.respaldos} accion={<span className="flex w-full flex-wrap items-center gap-2 tablet:w-auto">{boton}<IndicadorGuardado g={g} /></span>}>
      <Cargar c={c} que={t.que} esqueleto={<div className="flex flex-col gap-2">{[0, 1, 2].map((i) => <Esqueleto key={i} className="h-10" />)}</div>}>{(d) => {
        const ok = d.items.filter((b) => b.status === 'ok')
        const probada = d.items.filter((b) => b.restore_tested_at).sort((a, b) => b.restore_tested_at!.localeCompare(a.restore_tested_at!))[0]
        return (
          <>
            <p className="text-small tabular-nums text-text-secondary">{t.encabezado(ok[0] ? fechaHora(ok[0].finished_at ?? ok[0].started_at) : '—', job?.next_run_at ? fechaHora(job.next_run_at) : '—', probada ? fechaHora(probada.restore_tested_at!) : '—')}</p>
            <VistaBD def={{
              clave: 'P-26.11-respaldos', filas: [...d.items].sort((a, b) => b.started_at.localeCompare(a.started_at)), idFila: (b) => b.id, vistas: ['tabla'], titulo: (b) => fechaHora(b.started_at), alAbrir: setDetalle,
              vacio: <Vacio texto={t.sinRespaldos} accion={boton} />,
              propiedades: [
                { id: 'started_at', nombre: t.fecha, valor: (b) => b.started_at, pintar: (b) => fechaHora(b.started_at) },
                { id: 'kind', nombre: t.tipo, valor: (b) => t.tipos[b.kind] },
                { id: 'status', nombre: t.resultado, valor: (b) => t.estados[b.status], pintar: (b) => <Resultado b={b} /> },
                { id: 'size_bytes', nombre: t.tamano, valor: (b) => b.size_bytes ?? null, pintar: (b) => tamano(b.size_bytes), alinear: 'derecha' },
                { id: 'integrity_check', nombre: t.integridad, valor: integridad },
                { id: 'restore', nombre: t.restauracion, valor: (b) => b.restore_tested_at ?? null, pintar: (b) => (b.restore_tested_at ? `${fechaHora(b.restore_tested_at)} · ${b.restore_test_result === 'ok' ? t.integro : b.restore_test_result ?? ''}` : '—') },
                { id: 'file_path', nombre: t.ubicacion, valor: (b) => b.file_path ?? null, pintar: (b) => (b.file_path ? `${b.file_path} · ${(b.sha256 ?? '').slice(0, 12)}` : '—') },
              ],
            }} />
          </>
        )
      }}</Cargar>
      <Hoja abierta={hoja === 'respaldo'} alCambiar={(a) => !a && navegar(conConsulta({ hoja: null }))} titulo={t.respaldar} pie={<>{boton}<IndicadorGuardado g={g} /></>}><p className="text-small text-text-secondary">{t.proposito}</p></Hoja>
      {detalle && <Detalle b={detalle} cerrar={() => setDetalle(null)} />}
    </Bloque>
  )
}
function Resultado({ b }: { b: B }) {
  const Icono = b.status === 'ok' ? CircleCheck : b.status === 'failed' ? CircleAlert : LoaderCircle
  return <span className="inline-flex items-center gap-1"><Icono aria-hidden className={cn('size-5', b.status === 'ok' ? 'text-role-success' : b.status === 'failed' ? 'text-role-error' : 'text-role-info')} />{t.estados[b.status]}{b.status === 'failed' && b.error ? `: ${b.error.slice(0, 80)}` : ''}</span>
}
function Detalle({ b, cerrar }: { b: B; cerrar: () => void }) {
  const sinRed = useSinRed()
  const avisar = useAvisos()
  const [paso, setPaso] = useState<0 | 1 | 2>(0)
  const [texto, setTexto] = useState('')
  const [error, setError] = useState('')
  const f = fechaHora(b.started_at)
  const restaurar = async () => {
    try {
      const r = await pedir('post', '/api/v1/backups/{backup_id}/restore', { ruta: { backup_id: b.id }, cuerpo: { confirmation: 'RESTAURAR' } })
      if (r.integrity_check !== 'ok') { setError(t.noRestaurado); setPaso(0); return }
      avisar({ texto: t.restaurado, rol: 'exito' })
      setTimeout(() => window.location.reload(), 1500)
    } catch (e) { setError(e instanceof ErrorApi ? `${t.noRestaurado}. ${e.message}` : t.noRestaurado); setPaso(0) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={f} descripcion={`${t.tipos[b.kind]} · ${t.estados[b.status]}`}
      pie={b.status === 'ok' && <Boton variante="destructivo" tam="hoja" deshabilitadoPor={sinRed} onClick={() => setPaso(1)}>{t.restaurar}</Boton>}>
      <dl className="grid grid-cols-2 gap-2 text-small">
        <dt className="text-text-secondary">{t.tamano}</dt><dd className="tabular-nums">{tamano(b.size_bytes)}</dd>
        <dt className="text-text-secondary">{t.integridad}</dt><dd>{integridad(b)}</dd>
        <dt className="text-text-secondary">{t.ubicacion}</dt><dd className="break-all">{b.file_path ?? '—'} {b.sha256?.slice(0, 12)}</dd>
      </dl>
      {error && <p role="alert" className="mt-2 text-small text-role-error">{error}</p>}
      <DialogoCritico abierto={paso === 1} alCambiar={(a) => !a && setPaso(0)} titulo={t.restaurarTitulo} verbo={t.siguiente} alConfirmar={() => setPaso(2)}
        texto={t.restaurarTexto(f)}>
        <AntesDespues cambios={[]} nota={`${t.tamano}: ${tamano(b.size_bytes)} · ${t.integridad}: ${integridad(b)}`} />
      </DialogoCritico>
      <DialogoCritico abierto={paso === 2} alCambiar={(a) => { if (!a) { setPaso(0); setTexto('') } }} titulo={t.restaurarTitulo} verbo={t.restaurarTitulo} bloqueado={texto !== 'RESTAURAR'} alConfirmar={() => void restaurar()}
        texto={t.restaurarTexto(f)}>
        <label className="flex flex-col gap-1 text-small font-medium">{t.escribe}<input className={estilosControl} value={texto} onChange={(e) => setTexto(e.target.value)} autoComplete="off" /></label>
      </DialogoCritico>
    </Hoja>
  )
}
function HoraRespaldo({ job, recargar }: { job: J; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const actual = cronAHora(job.schedule)
  const [texto, setTexto] = useState(actual ?? '')
  const [error, setError] = useState('')
  const guardar = async () => {
    if (texto === (actual ?? '')) return
    if (!esHora(texto)) { setError(T.validar.hora); setTexto(actual ?? ''); return }
    setError('')
    const [h, m] = texto.split(':').map(Number)
    if (!(await g.ejecutar(async () => { await pedir('patch', '/api/v1/jobs/{job_id}', { ruta: { job_id: job.id }, cuerpo: { schedule: `${m} ${h} * * *` } }); recargar() }))) setTexto(actual ?? '')
  }
  return (
    <Bloque titulo={t.horaRespaldo} ayuda={t.horaRespaldoAyuda}>
      <div className="flex flex-wrap items-center gap-2">
        <label htmlFor="hora-respaldo" className="sr-only">{t.horaRespaldo}</label>
        <input id="hora-respaldo" className={cn(estilosControl, 'w-32 tabular-nums')} placeholder="hh:mm" inputMode="numeric" value={texto} disabled={!!sinRed || actual === null} aria-invalid={error ? true : undefined}
          onChange={(e) => setTexto(e.target.value)} onBlur={() => void guardar()} onKeyDown={(e) => { if (e.key === 'Enter') void guardar(); else if (e.key === 'Escape') setTexto(actual ?? '') }} />
        <span className="text-caption text-text-secondary">{`${T.r.proximo}: ${job.next_run_at ? fechaHora(job.next_run_at) : '—'}`}</span>
        <IndicadorGuardado g={g} />
      </div>
      {error && <p className="text-caption text-role-error">{error}</p>}
    </Bloque>
  )
}
function Importar() {
  const sinRed = useSinRed()
  const g = useGuardado()
  const [archivo, setArchivo] = useState<File | null>(null)
  const [previa, setPrevia] = useState<Esquemas['ImportPreview'] | null>(null)
  const [creadas, setCreadas] = useState<number | null>(null)
  const clave = useRef(nuevaClave())
  const forma = (f: File) => { const d = new FormData(); d.append('file', f); return d }
  const elegir = (f: File | null) => {
    setArchivo(f); setPrevia(null); setCreadas(null)
    if (f) void g.ejecutar(async () => setPrevia(await pedir('post', '/api/v1/imports/planned-items/preview', { cuerpo: forma(f) as unknown as Cuerpo<'/api/v1/imports/planned-items/preview', 'post'> })))
  }
  const importar = () => archivo && g.ejecutar(async () => {
    const r = await pedir('post', '/api/v1/imports/planned-items', { cuerpo: forma(archivo) as unknown as Cuerpo<'/api/v1/imports/planned-items', 'post'>, idempotencia: clave.current })
    clave.current = nuevaClave(); setCreadas(r.created_count); setPrevia(null)
  })
  return (
    <Bloque titulo={t.importar}>
      <SelectorArchivo etiqueta={t.archivo} boton={textos.archivo.elegirCsv} accept=".csv,text/csv" archivo={archivo} alElegir={elegir} deshabilitadoPor={sinRed || undefined} />
      <IndicadorGuardado g={g} />
      {previa && (
        <>
          <p className="text-small">{t.resumenImport(previa.valid_count, previa.error_count)}</p>
          <ul aria-label={t.seCrearan} className="flex max-h-64 flex-col gap-1 overflow-auto text-small">
            {previa.rows.map((r) => (
              <li key={r.line_no} className={cn('flex flex-wrap gap-2 rounded-sm p-1', r.errors.length > 0 && 'bg-role-error-surface')}>
                <span className="text-text-secondary">{`${t.renglon} ${r.line_no}`}</span>
                {r.planned_item && <span>{r.planned_item.description}</span>}
                {r.errors.map((e) => <span key={e.field} className="inline-flex items-center gap-1 text-role-error"><CircleAlert aria-hidden className="size-5" />{e.message}</span>)}
              </li>
            ))}
          </ul>
          <Boton variante="primario" deshabilitadoPor={sinRed ?? (previa.error_count > 0 ? t.corrige : undefined)} onClick={() => void importar()}>{t.importarBoton}</Boton>
          {previa.error_count > 0 && <p className="text-caption text-role-error">{t.corrige}</p>}
        </>
      )}
      {creadas !== null && <p role="status" className="text-small">{t.creadas(creadas)} · <Enlace a="/ajustes/02" className="text-text-link underline">{t.irPresupuesto}</Enlace></p>}
    </Bloque>
  )
}
function Proveedor({ p, recargar }: { p: AI; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const [dialogo, setDialogo] = useState(false)
  const poner = (v: boolean) => g.ejecutar(async () => { setDialogo(false); await pedir('patch', '/api/v1/ai/providers/{provider}', { ruta: { provider: p.provider }, cuerpo: { is_enabled: v } }); recargar() })
  const local = p.provider === 'local_model'
  return (
    <div className="flex flex-col gap-1 border-b border-border-divider py-2 last:border-b-0" data-prueba={`ia-${p.provider}`}>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="text-body font-medium">{t.proveedores[p.provider]}</span>
        <Interruptor ocultarEtiqueta etiqueta={`${T.n.activa} · ${t.proveedores[p.provider]}`} activo={p.is_enabled} deshabilitadoPor={sinRed} alCambiar={(v) => (v ? setDialogo(true) : void poner(false))} />
      </div>
      <p className="text-caption text-text-secondary">{p.allowed_uses.map((u) => t.usos[u]).join(' · ')}{p.is_enabled && ` · ${t.disponible}: ${p.is_available ? T.comun.si : T.comun.no}`}{!local && ` · ${t.llave}: ${p.is_configured ? t.configurada : t.sinConfigurar}`}</p>
      <IndicadorGuardado g={g} />
      <DialogoCritico abierto={dialogo} alCambiar={setDialogo} titulo={local ? t.verboLocal : t.verboClaude} texto={local ? t.activarLocal : t.activarClaude} verbo={local ? t.verboLocal : t.verboClaude} alConfirmar={() => void poner(true)} />
    </div>
  )
}
