// Piezas comunes de Ajustes 8 a 14 (U4 parte a §0): marco de sección, interruptor, guardado en el lugar,
// fila de parámetro operativo K44, diálogo de acción crítica y bloque «Antes → Después». Solo tokens C2.
import { ICONO_MD } from '../../componentes/Icono'
import { useCallback, useEffect, useRef, useState, useSyncExternalStore, type ReactNode } from 'react'
import { AlertDialog, Switch } from 'radix-ui'
import { CircleAlert, CircleCheck, RotateCcw, ShieldAlert } from 'lucide-react'
import { ErrorApi } from '../../api/errores'
import type { EstadoConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { ConDatos, Franja } from '../../componentes/Estados'
import { useAvisos } from '../../componentes/Avisos'
import { cn } from '../../componentes/utils'
import type { Cambio } from '../../vistas-bd/tipos'
import { textosAjustesB as t, DIAS_SEMANA } from './textos'
import { esHora } from './calculos'

const consultaEscritorio = '(min-width: 64rem)'
/** true desde el punto de quiebre «desktop» de C2 (64rem). */
export function useEsEscritorio() {
  return useSyncExternalStore((f) => { const m = window.matchMedia(consultaEscritorio); m.addEventListener('change', f); return () => m.removeEventListener('change', f) }, () => window.matchMedia(consultaEscritorio).matches)
}
/** Motivo para deshabilitar escrituras sin conexión (U4 §0.3). */
export const useSinRed = () => (useConexion() ? undefined : textos.conexion.requiere)

/** Contenido de una sección: el marco (columna de secciones, «‹ Ajustes» y título h2) lo pinta P26Ajustes del área ajustes-a. */
export function Marco({ id, children }: { id: string; children: ReactNode }) {
  return <div className="flex min-w-0 flex-col gap-6" data-prueba={`seccion-${id}`}>{children}</div>
}
export function Bloque({ titulo, children, accion, ayuda }: { titulo: string; children: ReactNode; accion?: ReactNode; ayuda?: ReactNode }) {
  return (
    <section className="flex flex-col gap-3 rounded-lg border border-border-divider bg-surface-card p-4">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h3 className="text-title-3 font-semibold text-text-default">{titulo}</h3>
        {accion}
      </div>
      {ayuda && <p className="text-small text-text-secondary">{ayuda}</p>}
      {children}
    </section>
  )
}
/** Une varias consultas en una sola (carga, error y recarga juntos). */
export function unir<T extends unknown[]>(cs: { [K in keyof T]: EstadoConsulta<T[K]> }): EstadoConsulta<T> {
  const l = cs as unknown as EstadoConsulta<unknown>[]
  const copias = l.map((c) => c.desdeCopia).filter((x): x is number => x !== null)
  return {
    datos: l.every((c) => c.datos !== undefined) ? (l.map((c) => c.datos) as T) : undefined,
    cargando: l.some((c) => c.cargando), error: l.find((c) => c.error)?.error ?? null,
    desdeCopia: copias.length ? Math.min(...copias) : null, nuncaCargada: l.some((c) => c.nuncaCargada),
    recargar: () => l.forEach((c) => c.recargar()),
  }
}
/** Consulta con su texto de error propio: «No se pudo cargar {qué} · Reintentar» (U4 §0.3). */
export function Cargar<T>({ c, que, children, esqueleto }: { c: EstadoConsulta<T>; que: string; children: (d: T) => ReactNode; esqueleto?: ReactNode }) {
  if (c.error && c.datos === undefined && !c.nuncaCargada && !c.cargando)
    return <Franja rol="error" alerta accion={<Boton variante="enlace" onClick={c.recargar}>{textos.app.reintentar}</Boton>}>{t.comun.noSeCargo(que)}</Franja>
  return <ConDatos consulta={c} esqueleto={esqueleto}>{children}</ConDatos>
}

export function Interruptor({ etiqueta, activo, alCambiar, deshabilitadoPor, ocultarEtiqueta }: { etiqueta: string; activo: boolean; alCambiar: (v: boolean) => void; deshabilitadoPor?: string; ocultarEtiqueta?: boolean }) {
  return (
    <label className="inline-flex min-h-11 min-w-11 cursor-pointer items-center gap-2 text-small text-text-default" title={deshabilitadoPor}>
      <Switch.Root checked={activo} onCheckedChange={(v) => !deshabilitadoPor && alCambiar(v)} aria-label={etiqueta} aria-disabled={deshabilitadoPor ? true : undefined}
        className="relative inline-flex h-6 w-11 shrink-0 items-center rounded-full border border-border-control bg-surface-pressed transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-disabled:cursor-not-allowed aria-disabled:opacity-60 data-[state=checked]:border-action-primary-bg data-[state=checked]:bg-action-primary-bg">
        <Switch.Thumb className="block size-4 translate-x-1 rounded-full bg-surface-card shadow-1 transition-transform data-[state=checked]:translate-x-6" />
      </Switch.Root>
      <span className={cn(ocultarEtiqueta && 'sr-only')}>{etiqueta}</span>
      <span className="text-text-secondary">{activo ? t.comun.encendido : t.comun.apagado}</span>
    </label>
  )
}

type EstadoG = 'inactivo' | 'guardando' | 'guardado' | 'error'
export function useGuardado() {
  const [estado, setEstado] = useState<EstadoG>('inactivo')
  const [mensaje, setMensaje] = useState('')
  const ultimo = useRef<(() => Promise<unknown>) | null>(null)
  const ejecutar = useCallback(async (fn: () => Promise<unknown>) => {
    ultimo.current = fn
    setEstado('guardando')
    try {
      await fn()
      setEstado('guardado')
      setTimeout(() => setEstado((e) => (e === 'guardado' ? 'inactivo' : e)), 2000)
      return true
    } catch (e) {
      setMensaje(e instanceof ErrorApi && e.estado === 422 ? e.message : t.comun.noSeGuardo)
      setEstado('error')
      return false
    }
  }, [])
  const reintentar = useCallback(() => { if (ultimo.current) void ejecutar(ultimo.current) }, [ejecutar])
  return { estado, mensaje, ejecutar, reintentar }
}
export function IndicadorGuardado({ g }: { g: ReturnType<typeof useGuardado> }) {
  return (
    <span aria-live="polite" className="inline-flex min-h-5 items-center gap-1 text-caption">
      {g.estado === 'guardando' && <span className="text-text-secondary">{t.comun.guardando}</span>}
      {g.estado === 'guardado' && <span className="inline-flex items-center gap-1 text-role-success"><CircleCheck aria-hidden className="size-5" />{t.comun.guardado}</span>}
      {g.estado === 'error' && (
        <span className="inline-flex flex-wrap items-center gap-1 text-role-error"><CircleAlert aria-hidden className="size-5" />{g.mensaje}
          {g.mensaje === t.comun.noSeGuardo && <Boton variante="enlace" onClick={g.reintentar}>{textos.app.reintentar}</Boton>}
        </span>
      )}
    </span>
  )
}

export function AntesDespues({ cambios, nota }: { cambios: Cambio[]; nota?: ReactNode }) {
  if (!cambios.length && !nota) return null
  return (
    <div role="status" aria-label={t.comun.vistaPrevia} className="flex flex-col gap-1 rounded-md border border-border-divider bg-surface-raised p-3 text-small text-text-default">
      {cambios.map((c) => (
        <p key={c.etiqueta} className="tabular-nums"><span className="text-text-secondary">{c.etiqueta}: </span>{c.antes} <span aria-hidden>→</span><span className="sr-only">{t.comun.despues}</span> <strong className="font-semibold">{c.despues}</strong></p>
      ))}
      {nota && <p className="text-text-secondary">{nota}</p>}
    </div>
  )
}

/** Parámetro operativo K44: valor vigente editable, «Valor inicial: …» en la misma fila y «Restablecer al valor inicial» si is_initial es falso. */
export type TipoParametro = 'entero' | 'hora' | 'lista' | 'diaSemana' | 'diaMes'
export interface Param<V> { key: string; value: V | null; initial_value: V | null; is_initial: boolean; unit?: string }
export function aTexto(tipo: TipoParametro, v: unknown): string {
  if (v === null || v === undefined) return ''
  if (tipo === 'lista') return (v as number[]).join(', ')
  return String(v)
}
export function mostrar(tipo: TipoParametro, v: unknown, unidad = ''): string {
  if (v === null || v === undefined) return t.comun.sinValor
  if (tipo === 'diaSemana') return DIAS_SEMANA[Number(v) - 1] ?? String(v)
  if (tipo === 'lista') return `[${(v as number[]).join(', ')}]`
  return unidad ? `${v} ${unidad}` : String(v)
}
export function interpretar(tipo: TipoParametro, s: string): number | string | number[] | null {
  const x = s.trim()
  if (tipo === 'hora') return esHora(x) ? x : null
  if (tipo === 'lista') {
    const l = x.split(/[,\s]+/).filter(Boolean)
    return l.length && l.every((d) => /^\d+$/.test(d)) ? l.map(Number) : null
  }
  if (!/^\d+$/.test(x)) return null
  const n = Number(x)
  if (tipo === 'diaSemana' && (n < 1 || n > 7)) return null
  if (tipo === 'diaMes' && (n < 1 || n > 31)) return null
  return n
}
export function FilaParametro<V extends number | string | number[]>({ etiqueta, param, tipo, mensaje, unidad, vistaPrevia, guardar, restablecer, ayuda, compacto }: {
  etiqueta: string; param: Param<V>; tipo: TipoParametro; mensaje: string; unidad?: string; ayuda?: ReactNode; compacto?: boolean
  vistaPrevia?: (nuevo: V) => { cambios: Cambio[]; nota?: ReactNode } | null
  guardar: (v: V) => Promise<unknown>; restablecer: () => Promise<unknown>
}) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const avisar = useAvisos()
  const [texto, setTexto] = useState(aTexto(tipo, param.value))
  const [error, setError] = useState('')
  useEffect(() => { setTexto(aTexto(tipo, param.value)) }, [param.value, tipo])
  const nuevo = interpretar(tipo, texto) as V | null
  const distinto = texto !== aTexto(tipo, param.value)
  const previa = distinto && nuevo !== null && vistaPrevia ? vistaPrevia(nuevo) : null
  const id = `param-${param.key}`
  const confirmar = async () => {
    if (!distinto) return
    if (nuevo === null) { setError(mensaje); setTexto(aTexto(tipo, param.value)); return }
    setError('')
    const ok = await g.ejecutar(() => guardar(nuevo))
    if (!ok) setTexto(aTexto(tipo, param.value))
  }
  const alRestablecer = async () => {
    const anterior = param.value
    const ok = await g.ejecutar(restablecer)
    if (ok && anterior !== null) avisar({ texto: t.comun.restablecido, rol: 'exito', accion: { texto: textos.app.deshacer, alHacer: () => { void guardar(anterior) } } })
  }
  const comun = {
    id, 'aria-invalid': error ? true : undefined, 'aria-describedby': `${id}-ayuda`, disabled: !!sinRed,
    onBlur: () => void confirmar(),
    onKeyDown: (e: React.KeyboardEvent) => { if (e.key === 'Enter') { e.preventDefault(); void confirmar() } else if (e.key === 'Escape') { setTexto(aTexto(tipo, param.value)); setError('') } },
  }
  return (
    <div className={cn('flex flex-col gap-1', !compacto && 'border-b border-border-divider py-3 last:border-b-0')} data-prueba={id}>
      <div className={cn('flex flex-col gap-2', !compacto && 'tablet:flex-row tablet:items-center')}>
        <label htmlFor={id} className={cn('text-small font-medium text-text-default', compacto ? 'sr-only' : 'tablet:w-64 tablet:shrink-0')}>{etiqueta}</label>
        <div className="flex flex-1 flex-wrap items-center gap-2">
          {tipo === 'diaSemana' ? (
            <select {...comun} value={texto} onChange={(e) => setTexto(e.target.value)} className={cn(estilosControl, 'w-40')}>
              {DIAS_SEMANA.map((d, i) => <option key={d} value={String(i + 1)}>{d}</option>)}
            </select>
          ) : (
            <input {...comun} value={texto} onChange={(e) => setTexto(e.target.value)} inputMode={tipo === 'hora' ? 'numeric' : tipo === 'lista' ? 'text' : 'numeric'} placeholder={tipo === 'hora' ? 'hh:mm' : undefined} className={cn(estilosControl, 'w-32 tabular-nums')} />
          )}
          {unidad && <span className="text-small text-text-secondary">{unidad}</span>}
          <span className="text-caption text-text-secondary">{t.comun.valorInicial(mostrar(tipo, param.initial_value, unidad))}</span>
          {!param.is_initial && <Boton variante="enlace" deshabilitadoPor={sinRed} onClick={() => void alRestablecer()}><RotateCcw aria-hidden className="size-5" />{t.comun.restablecer}</Boton>}
          <IndicadorGuardado g={g} />
        </div>
      </div>
      <p id={`${id}-ayuda`} className={cn('text-caption', error ? 'text-role-error' : 'text-text-secondary')}>
        {error ? <span className="inline-flex items-center gap-1"><CircleAlert aria-hidden className="size-5" />{error}</span> : ayuda}
      </p>
      {previa && <AntesDespues cambios={previa.cambios} nota={previa.nota} />}
    </div>
  )
}

/** Diálogo de acción crítica (U4 §0.4): shield-alert, verbo explícito, «Cancelar» con el foco inicial. */
export function DialogoCritico({ abierto, alCambiar, titulo, texto, verbo, alConfirmar, children, bloqueado }: {
  abierto: boolean; alCambiar: (a: boolean) => void; titulo: string; texto: ReactNode; verbo: string; alConfirmar: () => void; children?: ReactNode; bloqueado?: boolean
}) {
  return (
    <AlertDialog.Root open={abierto} onOpenChange={alCambiar}>
      <AlertDialog.Portal>
        <AlertDialog.Overlay className="fixed inset-0 z-40 bg-overlay-scrim" />
        <AlertDialog.Content className="fixed inset-x-0 bottom-0 z-50 flex flex-col gap-4 rounded-t-lg bg-surface-raised p-4 shadow-2 tablet:inset-x-auto tablet:bottom-auto tablet:left-1/2 tablet:top-1/2 tablet:w-full tablet:max-w-lg tablet:-translate-x-1/2 tablet:-translate-y-1/2 tablet:rounded-lg">
          <AlertDialog.Title className="flex items-center gap-2 text-title-3 font-semibold text-text-default"><ShieldAlert aria-hidden size={ICONO_MD} className="size-6 text-critical" />{titulo}</AlertDialog.Title>
          <AlertDialog.Description className="text-body text-text-default">{texto}</AlertDialog.Description>
          {children}
          <div className="flex flex-col-reverse gap-2 tablet:flex-row tablet:justify-end">
            <AlertDialog.Cancel asChild><Boton tam="normal" className="coarse:h-12">{t.comun.cancelar}</Boton></AlertDialog.Cancel>
            <Boton aria-disabled={bloqueado ? true : undefined} onClick={(e) => { if (bloqueado) { e.preventDefault(); return } alConfirmar() }} className="bg-action-critical-bg text-text-on-critical hover:bg-action-critical-bg-hover active:bg-action-critical-bg-pressed coarse:h-12">{verbo}</Boton>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  )
}
