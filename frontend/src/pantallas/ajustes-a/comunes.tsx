// Piezas comunes de Ajustes (U4 §0): bloque, fila de propiedad, edición en el lugar con validación y vista previa
// «Antes → Después», parámetro K44 con «Restablecer», interruptor, diálogo de acción crítica y validadores V-xx.
import { useEffect, useId, useRef, useState, type ReactNode } from 'react'
import { AlertDialog } from 'radix-ui'
import { ArrowRight, CircleAlert, CircleCheck, ShieldAlert, TriangleAlert } from 'lucide-react'
import { useConexion } from '../../api/conexion'
import type { EstadoConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import type { Cambio } from '../../vistas-bd/tipos'
import { aFraccion } from './calculos'
import { t } from './textos'

/** `bloquea`: falta un dato para calcular la vista previa (p. ej. la compra de ejemplo de P-26.05); mientras tanto no guarda (U4 §0.2). */
export interface Previa { cambios?: Cambio[]; notas?: ReactNode[]; pendiente?: string; bloquea?: boolean }
export type Validar = (texto: string) => string | null
/** Las páginas de C1 se generan como Page & {items}; esto deja los items con su tipo (C1 §5). */
export type Pagina<T> = { items: T[]; next_cursor: string | null }
export const comoLista = <T,>(c: EstadoConsulta<unknown>) => c as unknown as EstadoConsulta<Pagina<T>>

// ——— Validadores (U4 §0.6)
export const vObligatorio: Validar = (s) => (s.trim() ? null : t.v.v01)
export const vNombre: Validar = (s) => (!s.trim() ? t.v.v01 : s.trim().length > 120 ? t.v.v09 : null)
export const vDinero = (mayorQueCero = false, opcional = false): Validar => (s) => {
  if (!s.trim()) return opcional ? null : t.v.v01
  const c = aCentavos(s)
  if (c === null || c < 0) return t.v.v03
  return mayorQueCero && c === 0 ? t.v.v04 : null
}
export const vEntero = (min = 0, max?: number, opcional = false): Validar => (s) => {
  if (!s.trim()) return opcional ? null : t.v.v01
  const ok = /^\d+$/.test(s.trim()), n = Number(s)
  if (max !== undefined) return ok && n >= min && n <= max ? null : t.v.v23(min, max)
  if (!ok || n < min) return min === 1 ? t.v.v24 : t.v.v02
  return null
}
export const vPlazo: Validar = (s) => (/^\d+$/.test(s.trim()) && Number(s) >= 1 ? null : t.v.v21)
export const vPct = (dec: 2 | 4, opcional = false): Validar => (s) => (!s.trim() ? (opcional ? null : t.v.v01) : aFraccion(s, dec) ? null : dec === 4 ? t.v.v05 : t.v.v06)
export const vFecha = (opcional = false): Validar => (s) => (!s ? (opcional ? null : t.v.v07) : /^\d{4}-\d{2}-\d{2}$/.test(s) ? null : t.v.v07)
export const vDia: Validar = (s) => (/^\d+$/.test(s) && +s >= 1 && +s <= 31 ? null : t.v.v10)
export const vOpcion = (opcional = false): Validar => (s) => (s || opcional ? null : t.v.v22)
/** Texto que se muestra en un campo de dinero al editarlo: 18550 → «185.50». */
export const textoDinero = (c?: number | null) => (c == null ? '' : `${Math.floor(c / 100)}.${String(c % 100).padStart(2, '0')}`)

// ——— Estructura
export function Bloque({ titulo, children, accion, id }: { titulo: string; children: ReactNode; accion?: ReactNode; id?: string }) {
  return (
    <section aria-labelledby={id} className="flex flex-col gap-3 rounded-lg border border-border-divider bg-surface-card p-4">
      <div className="flex flex-wrap items-center justify-between gap-2"><h3 id={id} className="text-title-3 text-text-default">{titulo}</h3>{accion}</div>
      {children}
    </section>
  )
}
export function Fila({ etiqueta, children, nota }: { etiqueta: ReactNode; children: ReactNode; nota?: ReactNode }) {
  return (
    <div className="flex flex-col gap-1 border-b border-border-divider py-2 last:border-b-0 desktop:grid desktop:grid-cols-3 desktop:items-start desktop:gap-4">
      <div className="pt-2 text-small text-text-secondary">{etiqueta}{nota && <div className="text-caption">{nota}</div>}</div>
      <div className="min-w-0 desktop:col-span-2">{children}</div>
    </div>
  )
}
export function Aviso({ rol, children }: { rol: 'advertencia' | 'error' | 'exito' | 'info'; children: ReactNode }) {
  const Icono = rol === 'error' ? CircleAlert : rol === 'exito' ? CircleCheck : TriangleAlert
  const color = { advertencia: 'bg-role-warning-surface text-role-warning', error: 'bg-role-error-surface text-role-error', exito: 'bg-role-success-surface text-role-success', info: 'bg-role-info-surface text-role-info' }[rol]
  return <p role={rol === 'error' ? 'alert' : 'status'} className={cn('flex items-start gap-2 rounded-md px-3 py-2 text-small', color)}><Icono aria-hidden className="size-5 shrink-0" /><span className="text-text-default">{children}</span></p>
}
/** U4 §0.3 y §0.6: texto de un guardado fallido. Gana el V-xx del código que se pase (409 `in_use`, `already_exists`…);
 *  un 4xx de validación trae su mensaje (igual que `Editable`); lo demás es «No se pudo guardar…», que se puede reintentar. */
export function textoFallo(x: unknown, porCodigo: Record<string, string> = {}): string {
  if (x instanceof ErrorApi) {
    if (porCodigo[x.codigo]) return porCodigo[x.codigo]
    if (x.estado >= 400 && x.estado < 500 && x.estado !== 401) return x.campos[0]?.message ?? x.message
  }
  return t.comun.noGuardo
}
/** U4 §0.3 «Error al guardar» dentro de una hoja o un diálogo: `circle-alert` y el texto en el rol de error, con «Reintentar»
 *  cuando el fallo es de conexión o del servidor. Lo capturado no se pierde: quien lo usa deja la hoja abierta. */
export function ErrorGuardado({ texto, alReintentar }: { texto: string; alReintentar?: () => void }) {
  return (
    <p role="alert" data-prueba="error-guardado" className="flex flex-wrap items-center gap-1 text-small text-role-error">
      <CircleAlert aria-hidden className="size-5 shrink-0" />{texto}
      {alReintentar && texto === t.comun.noGuardo && <> · <button type="button" className="min-h-11 text-text-link underline-offset-2 hover:underline desktop:min-h-0" onClick={alReintentar}>{textos.app.reintentar}</button></>}
    </p>
  )
}
export function BloquePrevia({ previa }: { previa: Previa }) {
  return (
    <div aria-live="polite" className="flex flex-col gap-1 rounded-md border border-border-divider bg-surface-raised p-2 text-small tabular-nums" data-prueba="vista-previa">
      {previa.cambios && previa.cambios.length > 0 && (
        <dl className="flex flex-col gap-1">
          <div className="flex gap-2 text-caption text-text-secondary"><span className="flex-1" /><span>{t.comun.antes}</span><ArrowRight aria-hidden className="size-5" /><span>{t.comun.despues}</span></div>
          {previa.cambios.map((c) => <div key={c.etiqueta} className="flex flex-wrap items-center gap-2"><dt className="flex-1 text-text-secondary">{c.etiqueta}</dt><dd>{c.antes}</dd><ArrowRight aria-label={t.comun.cambiaA} className="size-5" /><dd className="font-medium">{c.despues}</dd></div>)}
        </dl>
      )}
      {previa.notas?.map((n, i) => <p key={i} className="text-text-secondary">{n}</p>)}
      {previa.pendiente && <p className="text-role-warning"><span className="font-medium">{t.comun.previaPendiente}</span> · {previa.pendiente}</p>}
    </div>
  )
}
function Estado({ estado, alReintentar }: { estado: 'guardando' | 'guardado' | 'error' | null; alReintentar?: () => void }) {
  return (
    <span aria-live="polite" className="text-small">
      {estado === 'guardando' && <span className="text-text-secondary">{t.comun.guardando}</span>}
      {estado === 'guardado' && <span className="inline-flex items-center gap-1 text-role-success"><CircleCheck aria-hidden className="size-5" />{t.comun.guardado}</span>}
      {estado === 'error' && <span className="inline-flex flex-wrap items-center gap-1 text-role-error"><CircleAlert aria-hidden className="size-5" />{t.comun.noGuardo} · <button type="button" className="min-h-11 text-text-link underline-offset-2 hover:underline desktop:min-h-0" onClick={alReintentar}>{textos.app.reintentar}</button></span>}
    </span>
  )
}
function useEstadoGuardado() {
  const [estado, setEstado] = useState<'guardando' | 'guardado' | 'error' | null>(null)
  const reloj = useRef<number>(0)
  useEffect(() => () => window.clearTimeout(reloj.current), [])
  const ok = () => { setEstado('guardado'); window.clearTimeout(reloj.current); reloj.current = window.setTimeout(() => setEstado(null), 2000) }
  return { estado, setEstado, ok }
}

// ——— Edición en el lugar (U4 §0.2 y §0.3)
interface PropsEditable {
  etiqueta: string
  /** Texto con que abre el control (dinero «185.50», fecha ISO, «si»/«no» en interruptores). */
  valor: string
  /** Cómo se ve sin editar (dinero(), fecha()…). */
  mostrar: ReactNode
  tipo?: 'texto' | 'dinero' | 'numero' | 'fecha' | 'opcion' | 'interruptor' | 'area'
  opciones?: { valor: string; texto: string }[]
  validar?: Validar
  vistaPrevia?: (nuevo: string) => Previa | Promise<Previa>
  guardar: (nuevo: string) => Promise<void>
  soloLectura?: boolean
  /** Controles que alimentan la vista previa (compra de ejemplo de P-26.05); se muestran mientras se edita. */
  extra?: ReactNode
  /** Cuando cambia, la vista previa se vuelve a calcular con el valor tecleado. */
  dependencias?: string
}
export function Editable({ etiqueta, valor, mostrar, tipo = 'texto', opciones, validar, vistaPrevia, guardar, soloLectura, extra, dependencias }: PropsEditable) {
  const conexion = useConexion()
  const [editando, setEditando] = useState(false)
  const [texto, setTexto] = useState(valor)
  const [error, setError] = useState<string | null>(null)
  const [previa, setPrevia] = useState<Previa | 'error' | null>(null)
  const { estado, setEstado, ok } = useEstadoGuardado()
  const ultimo = useRef(valor)
  const ref = useRef<HTMLElement>(null)
  const idErr = useId()
  useEffect(() => { if (editando) ref.current?.focus() }, [editando])
  useEffect(() => { if (editando) calcular(texto) }, [dependencias])
  const calcular = (v: string) => {
    const e = validar?.(v) ?? null
    setError(e)
    if (!vistaPrevia || e || v === valor) return setPrevia(null)
    Promise.resolve().then(() => vistaPrevia(v)).then(setPrevia, () => setPrevia('error'))
  }
  const cancelar = () => { setEditando(false); setTexto(valor); setError(null); setPrevia(null) }
  const enviar = async (v: string) => {
    setEstado('guardando'); ultimo.current = v
    try { await guardar(v); ok() } catch (e) {
      if (e instanceof ErrorApi && e.estado >= 400 && e.estado < 500 && e.estado !== 401) { setEstado(null); setEditando(true); setTexto(v); setError(e.campos[0]?.message ?? e.message); return }
      setEstado('error')
    }
  }
  const confirmar = () => {
    if (texto === valor) return cancelar()
    const e = validar?.(texto) ?? null
    if (e) return setError(e)
    if (vistaPrevia && (previa === null || previa === 'error' || previa.bloquea)) return // la vista previa debe verse antes de guardar
    setEditando(false); setPrevia(null); void enviar(texto)
  }
  const teclas = (e: React.KeyboardEvent) => { if (e.key === 'Enter' && tipo !== 'area') { e.preventDefault(); confirmar() } if (e.key === 'Escape') { e.preventDefault(); cancelar() } }
  if (!editando) {
    const motivo = !conexion ? textos.conexion.requiere : undefined
    return (
      <div className="flex flex-col gap-1">
        <button type="button" disabled={soloLectura} aria-disabled={motivo ? true : undefined} title={motivo} aria-label={`${t.comun.editar(etiqueta)}: ${typeof mostrar === 'string' ? mostrar : texto || t.comun.sinValor}`}
          onClick={() => { if (motivo) return; setTexto(valor); setEditando(true); if (tipo === 'interruptor') calcular(valor) }}
          className="flex min-h-11 w-full items-center rounded-sm px-2 text-left text-body hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-disabled:cursor-not-allowed desktop:min-h-9">
          {tipo === 'interruptor' ? <MarcaInterruptor encendido={valor === 'si'} /> : mostrar}
        </button>
        {motivo && <span className="sr-only">{motivo}</span>}
        <Estado estado={estado} alReintentar={() => void enviar(ultimo.current)} />
      </div>
    )
  }
  const comunes = { ref: ref as never, 'aria-label': etiqueta, 'aria-invalid': error ? true : undefined, 'aria-describedby': error ? idErr : undefined, onKeyDown: teclas, onBlur: (e: React.FocusEvent) => { if (!e.currentTarget.parentElement?.contains(e.relatedTarget as Node)) confirmar() } }
  return (
    <div className="flex flex-col gap-1">
      {tipo === 'opcion' ? <select {...comunes} className={estilosControl} value={texto} onChange={(e) => { setTexto(e.target.value); calcular(e.target.value) }}>{opciones?.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}</select>
        : tipo === 'interruptor' ? <button {...comunes} type="button" role="switch" aria-checked={texto === 'si'} className="flex min-h-11 items-center rounded-sm px-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus desktop:min-h-9" onClick={() => { const v = texto === 'si' ? 'no' : 'si'; setTexto(v); calcular(v) }}><MarcaInterruptor encendido={texto === 'si'} /></button>
        : tipo === 'area' ? <textarea {...comunes} rows={3} className={cn(estilosControl, 'h-auto py-2')} value={texto} onChange={(e) => { setTexto(e.target.value); calcular(e.target.value) }} />
        : <input {...comunes} className={estilosControl} type={tipo === 'fecha' ? 'date' : 'text'} inputMode={tipo === 'dinero' || tipo === 'numero' ? 'decimal' : undefined} value={texto} onChange={(e) => { setTexto(e.target.value); calcular(e.target.value) }} />}
      {error && <p id={idErr} className="flex items-center gap-1 text-small text-role-error"><CircleAlert aria-hidden className="size-5" />{error}</p>}
      {previa === 'error' && <p role="status" className="text-small text-role-error">{t.comun.noCalculo} · <button type="button" className="min-h-11 text-text-link desktop:min-h-0" onMouseDown={(e) => e.preventDefault()} onClick={() => calcular(texto)}>{textos.app.reintentar}</button></p>}
      {extra && <div onKeyDown={teclas}>{extra}</div>}
      {previa && previa !== 'error' && <BloquePrevia previa={previa} />}
    </div>
  )
}
export function MarcaInterruptor({ encendido }: { encendido: boolean }) {
  return (
    <span className="inline-flex items-center gap-2">
      <span aria-hidden className={cn('relative inline-block h-6 w-10 rounded-full border border-border-control transition-colors', encendido ? 'bg-action-primary-bg' : 'bg-surface-disabled')}>
        <span className={cn('absolute top-0.5 size-4 rounded-full bg-surface-card shadow-1 transition-all', encendido ? 'left-5' : 'left-0.5')} />
      </span>
      <span className="text-small">{encendido ? t.comun.si : t.comun.no}</span>
    </span>
  )
}

/** Interruptor que guarda al tocar (sin vista previa); con `critico`, pide confirmar antes de encender (§0.4). */
export function Interruptor({ etiqueta, encendido, guardar, critico, bloqueadoPor }: { etiqueta: string; encendido: boolean; guardar: (v: boolean) => Promise<void>; critico?: { titulo: string; texto: string; verbo: string }; bloqueadoPor?: string }) {
  const conexion = useConexion()
  const { estado, setEstado, ok } = useEstadoGuardado()
  const [pregunta, setPregunta] = useState(false)
  const motivo = bloqueadoPor ?? (!conexion ? textos.conexion.requiere : undefined)
  const cambiar = async (v: boolean) => { setEstado('guardando'); try { await guardar(v); ok() } catch { setEstado('error') } }
  return (
    <div className="flex flex-col gap-1">
      <button type="button" role="switch" aria-checked={encendido} aria-label={etiqueta} aria-disabled={motivo ? true : undefined} title={motivo}
        onClick={() => { if (motivo) return; if (critico && !encendido) return setPregunta(true); void cambiar(!encendido) }}
        className="flex min-h-11 w-fit items-center rounded-sm px-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-disabled:cursor-not-allowed desktop:min-h-9">
        <MarcaInterruptor encendido={encendido} />
      </button>
      {motivo && <span className="text-small text-text-secondary">{motivo}</span>}
      <Estado estado={estado} alReintentar={() => void cambiar(!encendido)} />
      {critico && <DialogoCritico abierto={pregunta} alCambiar={setPregunta} titulo={critico.titulo} verbo={critico.verbo} alConfirmar={() => cambiar(true)}>{critico.texto}</DialogoCritico>}
    </div>
  )
}

/** Diálogo de acción crítica (U4 §0.4, DESIGN §7.12): «Cancelar» con el foco inicial; hoja inferior en móvil. */
export type AccionFallo = { texto: string; accion?: { texto: string; alHacer: () => Promise<void> | void } }
/** Si la confirmación falla, el diálogo queda abierto con el error de §0.3 y «Reintentar»; con `porCodigo`, un código de C1
 *  (p. ej. 409 `in_use`) muestra su V-xx y, si se da, una acción alterna (p. ej. «Desactivar») en lugar del verbo crítico. */
export function DialogoCritico({ abierto, alCambiar, titulo, verbo, children, alConfirmar, porCodigo }: { abierto: boolean; alCambiar: (a: boolean) => void; titulo: string; verbo: string; children: ReactNode; alConfirmar: () => Promise<void> | void; porCodigo?: Record<string, AccionFallo> }) {
  const [enviando, setEnviando] = useState(false)
  const [fallo, setFallo] = useState<(AccionFallo & { reintentar?: () => void }) | null>(null)
  useEffect(() => { if (!abierto) setFallo(null) }, [abierto])
  const correr = async (paso: () => Promise<void> | void) => {
    setEnviando(true); setFallo(null)
    try { await paso(); alCambiar(false) }
    catch (x) {
      const propio = x instanceof ErrorApi ? porCodigo?.[x.codigo] : undefined
      setFallo(propio ?? { texto: textoFallo(x), reintentar: () => void correr(paso) })
    } finally { setEnviando(false) }
  }
  const alterna = fallo?.accion
  return (
    <AlertDialog.Root open={abierto} onOpenChange={alCambiar}>
      <AlertDialog.Portal>
        <AlertDialog.Overlay className="fixed inset-0 z-40 bg-overlay-scrim" />
        <AlertDialog.Content className="fixed inset-x-0 bottom-0 z-50 flex flex-col gap-4 rounded-t-lg bg-surface-card p-4 shadow-2 desktop:inset-x-auto desktop:bottom-auto desktop:left-1/2 desktop:top-1/3 desktop:w-full desktop:max-w-lg desktop:-translate-x-1/2 desktop:rounded-lg">
          <AlertDialog.Title className="flex items-center gap-2 text-title-3 text-text-default"><ShieldAlert aria-hidden className="size-5 text-critical" />{titulo}</AlertDialog.Title>
          <AlertDialog.Description className="text-body text-text-default">{children}</AlertDialog.Description>
          {fallo && <ErrorGuardado texto={fallo.texto} alReintentar={fallo.reintentar} />}
          <div className="flex flex-col-reverse gap-2 desktop:flex-row desktop:justify-end">
            <AlertDialog.Cancel asChild><Boton variante="secundario" className="coarse:h-11">{textos.app.cancelar}</Boton></AlertDialog.Cancel>
            {alterna
              ? <Boton variante="primario" className="coarse:h-11" disabled={enviando} onClick={() => void correr(alterna.alHacer)}>{alterna.texto}</Boton>
              : <button type="button" disabled={enviando} onClick={() => void correr(alConfirmar)}
                  className="inline-flex h-11 items-center justify-center rounded-md bg-action-critical-bg px-4 text-small font-medium text-text-on-critical hover:bg-action-critical-bg-hover active:bg-action-critical-bg-pressed focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus desktop:h-9">{verbo}</button>}
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  )
}

// ——— Parámetro operativo K44 (U4 §0.5)
// `unidad` reemplaza a `param.unit` solo donde la API manda una descripción larga en lugar de una unidad (las 5 claves de K54:
// su `unit` es la tercera columna de @param de 00005_tablero.sql). PROVISIONAL, bloqueo docs/estado/bloqueos/C1-unidad-parametros-tablero.md.
type Param = Esquemas['SettingParameterInteger']
export function Parametro({ etiqueta, param, unidad, validar = vEntero(0), vistaPrevia, guardar, restablecer }: { etiqueta: string; param: Param; unidad?: string; validar?: Validar; vistaPrevia?: (n: number) => Previa; guardar: (n: number) => Promise<void>; restablecer: () => Promise<void> }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const v = param.value ?? param.initial_value
  const conUnidad = (n: number | null) => (n == null ? t.comun.sinValor : `${n} ${unidad ?? param.unit}`)
  return (
    <Fila etiqueta={etiqueta} nota={t.comun.valorInicial(conUnidad(param.initial_value))}>
      <div className="flex flex-wrap items-start gap-2" data-prueba={`parametro-${param.key}`}>
        <div className="min-w-40 flex-1">
          <Editable etiqueta={etiqueta} tipo="numero" valor={v == null ? '' : String(v)} mostrar={<span className="tabular-nums">{conUnidad(v)}</span>} validar={validar}
            vistaPrevia={vistaPrevia ? (s) => vistaPrevia(Number(s)) : undefined} guardar={(s) => guardar(Number(s))} />
        </div>
        {!param.is_initial && (
          <Boton variante="fantasma" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined}
            onClick={async function volver() {
              const previo = v
              try { await restablecer() } catch (x) { avisar({ texto: textoFallo(x), rol: 'error', accion: { texto: textos.app.reintentar, alHacer: volver } }); return }
              avisar({ texto: t.comun.restablecido, rol: 'exito', accion: previo == null ? undefined : { texto: textos.app.deshacer, alHacer: () => guardar(previo) } })
            }}>{t.comun.restablecer}</Boton>
        )}
      </div>
    </Fila>
  )
}
