// P-26.01 Cuenta y roles (U4 parte a): usuarios, Karen sin cuenta, contraseña propia, sesiones, este dispositivo y parámetros de acceso (K44).
import { useState } from 'react'
import { Download, LogOut } from 'lucide-react'
import { pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { fechaHora } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useSesion } from '../../app/sesion'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ConDatos } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { Aviso, Bloque, BloquePrevia, DialogoCritico, Fila, Parametro } from './comunes'
import { comoLista } from './comunes'
import { t } from './textos'
const c = t.cuenta
type ClaveSeg = keyof Esquemas['SecuritySettings']['parameters']

export default function S01Cuenta() {
  const usuarios = comoLista<Esquemas['UserAccount']>(usarConsulta('get', '/api/v1/users'))
  const personas = comoLista<Esquemas['Person']>(usarConsulta('get', '/api/v1/persons'))
  const sesiones = comoLista<Esquemas['DeviceSession']>(usarConsulta('get', '/api/v1/sessions', { consulta: { sort: '-last_seen_at', limit: 200 } }))
  const seguridad = usarConsulta('get', '/api/v1/settings/security')
  return (
    <div className="flex flex-col gap-4">
      <Bloque titulo={c.usuarios} id="b-usuarios">
        <ConDatos consulta={usuarios}>{(d) => (
          <ul className="flex flex-col">
            {d.items.map((u) => (
              <li key={u.id} className="flex flex-wrap items-center gap-x-4 gap-y-1 border-b border-border-divider py-2 last:border-b-0">
                <span className="min-w-32 font-medium">{u.display_name ?? u.username}</span><span className="text-small text-text-secondary">{c.usuario}: {u.username}</span>
                <span className="rounded-full bg-surface-selected px-2 text-caption">{c.coTitular}</span>
                <span className="text-small text-text-secondary">{c.ultimoAcceso}: {u.last_login_at ? fechaHora(u.last_login_at) : c.nunca}</span>
              </li>
            ))}
            <ConDatos consulta={personas}>{(p) => <>{p.items.filter((x) => x.participation === 'external' && !x.has_account).map((x) => (
              <li key={x.id} className="flex flex-wrap items-center gap-x-4 py-2"><span className="min-w-32 font-medium">{x.display_name}</span><span className="text-small text-text-secondary">{c.karen}</span></li>))}</>}</ConDatos>
          </ul>
        )}</ConDatos>
      </Bloque>
      <div className="flex flex-col gap-4 desktop:grid desktop:grid-cols-2">
        <Contrasena minimo={seguridad.datos?.parameters['auth.password_min_length'].value ?? 12} />
        <Bloque titulo={c.parametros} id="b-parametros">
          <ConDatos consulta={seguridad}>{(s) => {
            const p = s.parameters
            const guardar = (k: ClaveSeg) => async (n: number) => { await pedir('patch', '/api/v1/settings/security', { cuerpo: { parameters: { [k]: n } } }); seguridad.recargar() }
            const restablecer = (k: ClaveSeg) => async () => { await pedir('patch', '/api/v1/settings/security', { cuerpo: { reset_parameters: [k] } }); seguridad.recargar() }
            const v = (k: ClaveSeg) => p[k].value ?? p[k].initial_value ?? 0
            const props = (k: ClaveSeg) => ({ param: p[k], guardar: guardar(k), restablecer: restablecer(k) })
            return (<>
              <Parametro etiqueta={c.pIntentos} {...props('auth.max_failed_attempts')} vistaPrevia={(n) => ({ cambios: [{ etiqueta: c.reglaPrevia, antes: c.previaIntentos(v('auth.max_failed_attempts'), v('auth.lockout_window_minutes')), despues: c.previaIntentos(n, v('auth.lockout_window_minutes')) }] })} />
              <Parametro etiqueta={c.pVentana} {...props('auth.lockout_window_minutes')} vistaPrevia={(n) => ({ cambios: [{ etiqueta: c.reglaPrevia, antes: c.previaIntentos(v('auth.max_failed_attempts'), v('auth.lockout_window_minutes')), despues: c.previaIntentos(v('auth.max_failed_attempts'), n) }] })} />
              <Parametro etiqueta={c.pSesion} {...props('session.max_age_days')} vistaPrevia={(n) => ({ cambios: [{ etiqueta: c.reglaPrevia, antes: c.previaSesion(v('session.max_age_days')), despues: c.previaSesion(n) }] })} />
              <Parametro etiqueta={c.pLargo} {...props('auth.password_min_length')} />
            </>)
          }}</ConDatos>
        </Bloque>
      </div>
      <Bloque titulo={c.sesiones} id="b-sesiones">
        <ConDatos consulta={sesiones}>{(d) => <Sesiones filas={d.items} usuarios={usuarios.datos?.items ?? []} alCambiar={sesiones.recargar} />}</ConDatos>
      </Bloque>
      <EsteDispositivo />
    </div>
  )
}

function Contrasena({ minimo }: { minimo: number }) {
  const conexion = useConexion()
  const [v, setV] = useState({ actual: '', nueva: '', repetir: '' })
  const [err, setErr] = useState<Partial<Record<'actual' | 'nueva' | 'repetir', string>>>({})
  const [hecho, setHecho] = useState(false)
  const [enviando, setEnviando] = useState(false)
  const enviar = async () => {
    const e: typeof err = {}
    if (!v.actual) e.actual = t.v.v01
    if (!v.nueva) e.nueva = t.v.v01; else if (v.nueva.length < minimo) e.nueva = t.v.v14(minimo)
    if (v.repetir !== v.nueva) e.repetir = t.v.v15
    setErr(e); setHecho(false)
    if (Object.keys(e).length) return
    setEnviando(true)
    try { await pedir('post', '/api/v1/auth/password', { cuerpo: { current_password: v.actual, new_password: v.nueva } }); setV({ actual: '', nueva: '', repetir: '' }); setHecho(true) }
    catch (x) { setErr(x instanceof ErrorApi ? (x.campo('new_password') ? { nueva: x.campo('new_password') } : { actual: x.campo('current_password') ?? x.message }) : { actual: t.comun.noGuardo }) }
    finally { setEnviando(false) }
  }
  return (
    <Bloque titulo={c.contrasena} id="b-contrasena">
      <form className="flex flex-col gap-3" onSubmit={(e) => { e.preventDefault(); void enviar() }} onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void enviar() } }}>
        <Campo etiqueta={c.actual} type="password" autoComplete="current-password" value={v.actual} error={err.actual} onChange={(e) => setV({ ...v, actual: e.target.value })} />
        <Campo etiqueta={c.nueva} type="password" autoComplete="new-password" value={v.nueva} error={err.nueva} ayuda={t.v.v14(minimo)} onChange={(e) => setV({ ...v, nueva: e.target.value })} />
        <Campo etiqueta={c.repetir} type="password" autoComplete="new-password" value={v.repetir} error={err.repetir} onChange={(e) => setV({ ...v, repetir: e.target.value })} />
        <Boton type="submit" variante="primario" className="w-fit" deshabilitadoPor={!conexion ? textos.conexion.requiere : enviando ? t.comun.guardando : undefined}>{c.cambiar}</Boton>
        {hecho && <Aviso rol="exito">{c.cambiada}</Aviso>}
      </form>
    </Bloque>
  )
}

type Ses = Esquemas['DeviceSession']
function Sesiones({ filas, usuarios, alCambiar }: { filas: Ses[]; usuarios: Esquemas['UserAccount'][]; alCambiar: () => void }) {
  const conexion = useConexion()
  const [una, setUna] = useState<Ses | null>(null)
  const [alcance, setAlcance] = useState('')
  const [critico, setCritico] = useState(false)
  const nombre = (s: Ses) => usuarios.find((u) => u.id === s.user_account_id)?.display_name ?? s.username ?? ''
  const actual = filas.find((s) => s.is_current)
  const [scope, uid] = alcance.split(':')
  const afectadas = filas.filter((s) => (scope === 'all' ? true : !s.is_current && (scope === 'all_except_current' || String(s.user_account_id) === uid)))
  const quien = scope === 'user_except_current' ? (usuarios.find((u) => String(u.id) === uid)?.display_name ?? '') : c.todasPersonas
  const req = !conexion ? textos.conexion.requiere : undefined
  return (
    <div className="flex flex-col gap-3">
      <VistaBD def={{ clave: 'P-26.01-sesiones', filas, idFila: (s) => s.id, vistas: ['tabla'], titulo: (s) => s.device_label ?? c.sinNombre,
        propiedades: [
          { id: 'device_label', nombre: c.dispositivo, valor: (s) => s.device_label ?? '', pintar: (s) => <span>{s.device_label ?? c.sinNombre}{s.is_current && <span className="ml-2 rounded-full bg-surface-selected px-2 text-caption">{c.estaSesion}</span>}</span> },
          { id: 'user', nombre: c.usuario, valor: nombre, grupo: nombre },
          { id: 'created_at', nombre: c.inicio, valor: (s) => s.created_at, pintar: (s) => fechaHora(s.created_at) },
          { id: 'last_seen_at', nombre: c.ultimoUso, valor: (s) => s.last_seen_at, pintar: (s) => fechaHora(s.last_seen_at) },
          { id: 'expires_at', nombre: c.vence, valor: (s) => s.expires_at, pintar: (s) => fechaHora(s.expires_at) },
          { id: 'cerrar', nombre: c.cerrar, valor: () => '', pintar: (s) => s.is_current ? null : <Boton tam="normal" deshabilitadoPor={req} onClick={() => setUna(s)}>{c.cerrar}</Boton> },
        ] }} />
      <div className="flex flex-wrap items-end gap-2">
        <label className="flex min-w-60 flex-1 flex-col gap-1 text-small text-text-secondary">{c.alcance}
          <select className={estilosControl} value={alcance} onChange={(e) => setAlcance(e.target.value)}>
            <option value="">{t.v.v22}</option>
            {usuarios.map((u) => <option key={u.id} value={`user_except_current:${u.id}`}>{c.deUsuario(u.display_name ?? u.username)}</option>)}
            <option value="all_except_current">{c.salvoActual}</option><option value="all">{c.todas}</option>
          </select>
        </label>
        <Boton variante="secundario" deshabilitadoPor={req ?? (!alcance ? t.v.v22 : undefined)} onClick={() => setCritico(true)}>{c.cerrarVarias}</Boton>
      </div>
      <DialogoCritico abierto={critico} alCambiar={setCritico} titulo={c.cerrarVarias} verbo={c.verbo(afectadas.length)}
        alConfirmar={async () => { await pedir('post', '/api/v1/sessions/revoke', { cuerpo: { scope: scope as 'all', ...(uid ? { user_account_id: Number(uid) } : {}) } }); setAlcance(''); alCambiar() }}>
        {c.dialogo(afectadas.length, quien)}{scope === 'all' && actual ? ` ${c.salirPregunta}` : ''}
      </DialogoCritico>
      <Hoja abierta={!!una} alCambiar={(a) => !a && setUna(null)} titulo={c.cerrar} pie={<div className="flex justify-end gap-2"><Boton onClick={() => setUna(null)}>{textos.app.cancelar}</Boton><Boton variante="primario" onClick={async () => { if (una) await pedir('delete', '/api/v1/sessions/{session_id}', { ruta: { session_id: una.id } }); setUna(null); alCambiar() }}>{textos.app.confirmar}</Boton></div>}>
        <p>{una && c.cerrarUna(una.device_label ?? c.sinNombre)}</p>
      </Hoja>
    </div>
  )
}

function EsteDispositivo() {
  const { salir } = useSesion()
  const conexion = useConexion()
  const [pregunta, setPregunta] = useState(false)
  const instalada = typeof window !== 'undefined' && window.matchMedia?.('(display-mode: standalone)').matches
  const descargar = async () => {
    const b = (await pedir('get', '/api/v1/system/ca-certificate')) as unknown as Blob
    const a = document.createElement('a'); a.href = URL.createObjectURL(b instanceof Blob ? b : new Blob([String(b)])); a.download = 'finanzas-hogar-ca.pem'; a.click(); URL.revokeObjectURL(a.href)
  }
  return (
    <Bloque titulo={c.esteDispositivo} id="b-dispositivo">
      <Fila etiqueta={c.instalada}><span className="px-2">{instalada ? t.comun.si : t.comun.no}</span></Fila>
      <Fila etiqueta={c.certificado}><div className="flex flex-col gap-2"><Boton className={cn('w-fit')} deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void descargar()}><Download aria-hidden className="size-5" />{c.descargar}</Boton><p className="text-small text-text-secondary">{c.instrucciones}</p></div></Fila>
      <Fila etiqueta={c.salir}><Boton variante="secundario" className="w-fit" onClick={() => setPregunta(true)}><LogOut aria-hidden className="size-5" />{c.salir}</Boton></Fila>
      <Hoja abierta={pregunta} alCambiar={setPregunta} titulo={c.salir} pie={<div className="flex justify-end gap-2"><Boton onClick={() => setPregunta(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" onClick={() => void salir()}>{c.salir}</Boton></div>}>
        <BloquePrevia previa={{ notas: [c.salirPregunta] }} />
      </Hoja>
    </Bloque>
  )
}
